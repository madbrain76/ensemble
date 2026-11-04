using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;

// The thinking model that spends its whole budget deliberating and never answers.
// Recovery is bounded, and it never rewrites an answer the operator already saw.
static class OutputLimitTests
{
    static void Check(bool value, string message)
    {
        if (!value) throw new Exception(message);
    }

    static readonly Provenance Origin = new(Vendor.OpenAI, "gpt-4o", Surface.ChatCompletions);

    static ParsedResponse Ended(params Part[] parts) => new([
        new Event(EventKind.ResponseEnded, actor: Actor.Agent, parts: parts)], []);

    // The provider must have said the budget ran out, or nothing is retried.
    static ParsedResponse Capped(params Part[] parts) => new([
        new Event(EventKind.ResponseEnded, actor: Actor.Agent, parts: parts)], [],
        OutputCapped: true);

    static JsonElement Thinking() => JsonSerializer.SerializeToElement(
        new { type = "thinking", thinking = "considering the request" });

    public static void Run()
    {
        // Only a deliberation-only completion counts as ending before an answer.
        Check(!Engine.EndsBeforeAnswer(Ended(new TextPart("the answer"))),
            "a visible answer was treated as truncated");
        Check(!Engine.EndsBeforeAnswer(Ended(new TextPart("  "), new TextPart("tail"))),
            "any visible prose should count as an answer");
        Check(Engine.EndsBeforeAnswer(Capped(new TextPart("   "))),
            "whitespace-only prose was accepted as an answer");
        Check(Engine.EndsBeforeAnswer(Capped(new OpaquePart(Origin, Thinking()))),
            "a thinking-only reply was not recognised as ending before an answer");
        Check(!Engine.EndsBeforeAnswer(Ended(new OpaquePart(Origin, Thinking()),
            new ToolCallPart("call-1", Origin, "read_file",
                JsonSerializer.SerializeToElement(new { path = "notes.md" })))),
            "a reply that asked for a tool call is not a silent truncation");
        Check(!Engine.EndsBeforeAnswer(new ParsedResponse([
            new Event(EventKind.ErrorOccurred, error: "HTTP exchange failed")], [])),
            "a transport failure was treated as an output limit");
        Check(!Engine.EndsBeforeAnswer(new ParsedResponse([
            new Event(EventKind.ResponseEnded, actor: Actor.Agent,
                parts: [new OpaquePart(Origin, Thinking())])], [])),
            "an uncapped thinking-only reply was called a truncation");

        // The cap doubles once and stops at the ceiling.
        Check(Engine.LargerOutputBudget("{\"model\":\"m\",\"max_completion_tokens\":1024}",
            out string doubled)
            && JsonNode.Parse(doubled)![
                "max_completion_tokens"]!.GetValue<int>() == 2048,
            "the chat cap was not doubled");
        Check(doubled.Contains("\"model\"") && doubled.Contains("\"m\""),
            "doubling dropped other request fields");
        Check(Engine.LargerOutputBudget("{\"max_tokens\":9000}", out string ceiling)
            && JsonNode.Parse(ceiling)!["max_tokens"]!.GetValue<int>() == 16384,
            "the doubled budget ignored the ceiling");
        Check(!Engine.LargerOutputBudget("{\"max_tokens\":16384}", out _),
            "a request already at the ceiling was raised again");
        Check(!Engine.LargerOutputBudget("{\"model\":\"m\"}", out _),
            "a request with no stated cap was \"doubled\"");
        Check(!Engine.LargerOutputBudget("not json", out _),
            "an unreadable body was treated as raisable");

        // The directive goes to whichever turn array this vendor uses, exactly once.
        Check(Engine.WithAnswerDirective("{\"model\":\"m\",\"messages\":[{\"role\":\"user\","
            + "\"content\":\"hi\"}]}", out string chat), "no directive for chat");
        var turns = JsonNode.Parse(chat)!["messages"]!.AsArray();
        Check(turns.Count == 2 && turns[1]!["role"]!.GetValue<string>() == "user"
            && turns[1]!["content"]!.GetValue<string>().Contains("visible text"),
            "the chat directive was not appended as one user turn");
        Check(Engine.WithAnswerDirective("{\"input\":[{\"role\":\"user\",\"content\":\"go\"}]}",
            out string responses) && JsonNode.Parse(responses)!["input"]!.AsArray().Count == 2,
            "no directive for the Responses shape");
        Check(Engine.WithAnswerDirective("{\"contents\":[{\"role\":\"user\",\"parts\":"
            + "[{\"text\":\"go\"}]}]}", out string gemini), "no directive for the Gemini shape");
        var added = JsonNode.Parse(gemini)!["contents"]!.AsArray()[1]!;
        Check(added["parts"]!.AsArray().Count == 1
            && added["parts"]![0]!["text"]!.GetValue<string>().Length > 0,
            "the Gemini directive did not use the parts shape");
        Check(!Engine.WithAnswerDirective("{\"model\":\"m\"}", out _),
            "a request with no turn array gained a directive anyway");

        // The ladder itself: bounded, and silent once anything was shown.
        Check(Engine.RecoveryBodies("{\"messages\":[],\"max_tokens\":1024}", true, false).Count == 0,
            "a published answer was scheduled for retry");
        Check(Engine.RecoveryBodies("{\"messages\":[],\"max_tokens\":1024}", false, true).Count == 0,
            "a cancelled turn was still retried");
        var ladder = Engine.RecoveryBodies("{\"messages\":[],\"max_tokens\":1024}", false, false);
        Check(ladder.Count == 2, $"the ladder offered {ladder.Count} attempts, expected 2");
        Check(JsonNode.Parse(ladder[0])!["max_tokens"]!.GetValue<int>() == 2048
            && JsonNode.Parse(ladder[0])!["messages"]!.AsArray().Count == 0,
            "the first rung changed more than the budget");
        var second = JsonNode.Parse(ladder[1])!;
        Check(second["max_tokens"]!.GetValue<int>() == 2048
            && second["messages"]!.AsArray().Count == 1,
            "the second rung did not build on the first");
        var single = Engine.RecoveryBodies("{\"messages\":[]}", false, false);
        Check(single.Count == 1 && JsonNode.Parse(single[0])!["messages"]!.AsArray().Count == 1,
            "an uncapped request should get only the directive rung");
    }
}
