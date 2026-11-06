using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;

// Chat Completions has no field for replaying a provider's own reasoning block,
// and an untyped member of a content array fails that endpoint's schema outright.
internal static class ReasoningReplayTests
{
    private static void Check(bool condition, string message)
    {
        if (!condition)
        {
            throw new InvalidOperationException(message);
        }
    }

    private static JsonElement Json(string value) => JsonSerializer.Deserialize<JsonElement>(value);

    private static IEnumerable<JsonObject> Parts(JsonObject request)
    {
        foreach (var node in request["messages"]!.AsArray())
        {
            if (node!["content"] is JsonArray blocks)
            {
                foreach (var block in blocks)
                {
                    yield return (JsonObject)block!;
                }
            }
        }
    }

    public static void Run()
    {
        var app = new Ensemble.Ensemble(TextWriter.Null);
        var chat = new Configuration(Vendor.OpenAI, "reasoning-replay-model", "http://127.0.0.1:1");
        var target = chat.Target();
        // Both shapes this program has really written: the field name the endpoint
        // returned, and the part shape this codec records when it reads that field.
        var returned = new OpaquePart(target, Json("""{"reasoning_content":"PRIVATE RETURNED THOUGHT"}"""));
        var recorded = new OpaquePart(target, Json("""{"type":"reasoning_content","text":"PRIVATE RECORDED THOUGHT"}"""));
        var context = new Context([
            new Entry(1, Actor.Human, [new TextPart("question")]),
            new Entry(2, Actor.Agent, [new TextPart("first answer"), returned]),
            new Entry(3, Actor.Human, [new TextPart("next")]),
            new Entry(4, Actor.Agent, [new TextPart("second answer"), recorded]),
            new Entry(5, Actor.Human, [new TextPart("last")]),
            new Entry(6, Actor.Agent, [returned])
        ]);
        var agent = new Agent(app, chat);
        var request = agent.Engine().Render(context);
        var wire = request.ToJsonString();
        agent.Destroy();

        Check(!wire.Contains("PRIVATE RETURNED THOUGHT") && !wire.Contains("PRIVATE RECORDED THOUGHT"),
            "provider-only reasoning was replayed to the endpoint: " + wire);
        Check(wire.Contains("first answer") && wire.Contains("second answer"),
            "dropping reasoning also dropped the delivered text: " + wire);
        // An array member without a type is what turned every resumed conversation
        // into a 422, so the whole content array is checked, not just reasoning.
        foreach (var block in Parts(request))
        {
            string type = block["type"]?.GetValue<string>() ?? "";
            Check(type == "text" || type == "image_url",
                "content array carries a part this endpoint cannot parse: " + block.ToJsonString());
        }

        var assistants = request["messages"]!.AsArray()
            .Count(node => node!["role"]!.GetValue<string>() == "assistant");
        Check(assistants == 2,
            $"a turn holding only reasoning should be omitted, saw {assistants} assistant messages: {wire}");
        Check(request["messages"]!.AsArray()
            .Where(node => node!["role"]!.GetValue<string>() == "assistant")
            .All(node => node!["content"] is JsonValue),
            "an ordinary text turn should stay a string content without media: " + wire);

        // The other dialects do accept their own opaque blocks, so the removal must
        // stop at this surface rather than quietly lose thinking replay everywhere.
        var claude = new Configuration(Vendor.Anthropic, "reasoning-replay-model", "http://127.0.0.1:1");
        var thinking = new OpaquePart(claude.Target(),
            Json("""{"type":"thinking","thinking":"KEEP THIS THOUGHT","signature":"sig-1"}"""));
        var messages = new Agent(app, claude);
        var kept = messages.Engine().Render(new Context([
            new Entry(1, Actor.Human, [new TextPart("question")]),
            new Entry(2, Actor.Agent, [thinking, new TextPart("answer")])
        ]));
        messages.Destroy();
        Check(kept.ToJsonString().Contains("KEEP THIS THOUGHT"),
            "the messages surface lost its own thinking replay: " + kept.ToJsonString());
    }
}
