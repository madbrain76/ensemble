using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;

// A tool call cut off by the output budget must be retried, never dispatched.
static class TruncatedToolCallTests
{
    static void Check(bool value, string message)
    {
        if (!value) throw new Exception(message);
    }

    static readonly Provenance Origin = new(Vendor.OpenAI, "gpt-4o", Surface.ChatCompletions);

    static IEnumerable<Event> With(params Part[] parts) =>
        new[] { new Event(EventKind.ResponseEnded, actor: Actor.Agent, parts: parts) };

    static Part Call() => new ToolCallPart("call-1", Origin, "write_file",
        JsonSerializer.SerializeToElement(new { path = "notes.md" }));

    public static void Run()
    {
        var openai = new Configuration(Vendor.OpenAI, "gpt-4o", "http://127.0.0.1:1");
        var anthropic = new Configuration(Vendor.Anthropic, "claude-3-haiku-20240307",
            "http://127.0.0.1:1", surface: Surface.Messages);
        var gemini = new Configuration(Vendor.Gemini, "gemini-2.5-flash",
            "http://127.0.0.1:1", surface: Surface.GenerateContent);

        Check(Engine.SaysLength(openai, "{\"choices\":[{\"finish_reason\":\"length\"}]}"),
            "the OpenAI length verdict was missed");
        Check(Engine.SaysLength(openai, "{\"stop_reason\":\"max_tokens\"}"),
            "the Responses-style max_tokens verdict was missed");
        Check(!Engine.SaysLength(openai, "{\"choices\":[{\"finish_reason\":\"stop\"}]}"),
            "a finished reply was treated as capped");
        Check(!Engine.SaysLength(openai, "{\"choices\":[{\"finish_reason\":\"content_filter\"}]}"),
            "a filtered reply was treated as capped");
        Check(!Engine.SaysLength(openai, "{}"), "a missing finish reason was treated as capped");
        Check(!Engine.SaysLength(openai, "not json"), "unreadable JSON was treated as capped");
        Check(Engine.SaysLength(anthropic, "{\"stop_reason\":\"max_tokens\"}"),
            "the Anthropic max_tokens verdict was missed");
        Check(Engine.SaysLength(gemini,
            "{\"candidates\":[{\"finishReason\":\"MAX_TOKENS\"}]}"),
            "the Gemini MAX_TOKENS verdict was missed");

        Check(Engine.ToolCallCapped(openai, "{\"choices\":[{\"finish_reason\":\"length\"}]}",
            With(Call())), "a capped tool call was not recognised");
        Check(!Engine.ToolCallCapped(openai, "{\"choices\":[{\"finish_reason\":\"length\"}]}",
            With(new TextPart("prose only"))),
            "a capped prose reply was treated as a truncated call");
        Check(!Engine.ToolCallCapped(openai, "{\"choices\":[{\"finish_reason\":\"stop\"}]}",
            With(Call())), "a complete tool call was treated as truncated");
        Check(!Engine.ToolCallCapped(openai, "", With(Call())),
            "a missing body was treated as truncated");

        Check(Engine.LargerOutputBudget(
            "{\"model\":\"m\",\"max_completion_tokens\":1024}", out string wider)
            && JsonNode.Parse(wider)!["max_completion_tokens"]!.GetValue<int>() == 2048
            && wider.Contains("\"model\""), "the budget was not doubled in place");
        Check(Engine.LargerOutputBudget("{\"max_output_tokens\":12000}", out string ceiling)
            && JsonNode.Parse(ceiling)!["max_output_tokens"]!.GetValue<int>() == 16384,
            "the doubling ignored the ceiling");
        Check(!Engine.LargerOutputBudget("{\"max_tokens\":16384}", out _),
            "a request at the ceiling was raised again");
        Check(!Engine.LargerOutputBudget("{\"model\":\"m\"}", out _),
            "a request with no cap was doubled");
        Check(!Engine.LargerOutputBudget("nope", out _), "unreadable JSON was rewritten");
    }
}
