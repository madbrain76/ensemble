using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;

// Chat Completions replays a provider's own reasoning block as message-level
// fields, never as a content member: an untyped array element fails that
// endpoint's schema outright, which is what turned resumed chats into 422s.
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
        // Both shapes this codec really records: the envelope a provider's own
        // reasoning block becomes, and the legacy part recorded before the envelope
        // existed. A raw response fragment with no type is neither, and the
        // endpoint refuses it, so it is checked separately below.
        var returned = new OpaquePart(target, Json("""{"type":"chat_reasoning","reasoning_content":"PRIVATE RETURNED THOUGHT","text":"PRIVATE RETURNED THOUGHT"}"""));
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
        // Streaming is a request property, not a client preference: an endpoint that
        // is never asked streams nothing, so no fragment ever reached the page and an
        // artifact had no arrivals to measure a span from.
        Check(request["stream"]?.GetValue<bool>() == true,
            "a chat-completions turn did not ask the endpoint to stream: " + wire);
        Check(request["stream_options"]?["include_usage"]?.GetValue<bool>() == true,
            "a streamed reply was not asked to carry usage, so counts would be lost: " + wire);
        agent.Destroy();

        // The operator's opt-out is honoured in the same place it is expressed.
        var buffered = new Agent(app, new Configuration(Vendor.OpenAI, chat.Model, "http://127.0.0.1:1",
            stream: "thinking"));
        var quiet = buffered.Engine().Render(context);
        buffered.Destroy();
        Check(quiet["stream"]?.GetValue<bool>() == true,
            "a channel selection that still streams something stopped asking for a stream");

        Check(wire.Contains("PRIVATE RETURNED THOUGHT") && wire.Contains("PRIVATE RECORDED THOUGHT"),
            "provider reasoning was not replayed for the endpoint that returned it: " + wire);
        Check(wire.Contains("first answer") && wire.Contains("second answer"),
            "replay lost the delivered text beside the reasoning: " + wire);
        // Reasoning now leaves the content array entirely, so any member that does
        // arrive must carry a type this endpoint recognizes: an untyped element is
        // what turned every resumed conversation into a 422.
        // A raw response fragment is not an envelope this codec records, so it is
        // dropped rather than handed to the endpoint as content.
        var stray = new Agent(app, chat);
        var strayWire = stray.Engine().Render(new Context([
            new Entry(1, Actor.Agent, [new OpaquePart(target,
                Json("""{"reasoning_content":"PRIVATE STRAY FRAGMENT"}"""))])
        ])).ToJsonString();
        stray.Destroy();
        Check(!strayWire.Contains("PRIVATE STRAY FRAGMENT"),
            "an untyped provider block was handed to the endpoint as content: " + strayWire);
        foreach (var block in Parts(request))
        {
            string type = block["type"]?.GetValue<string>() ?? "";
            Check(type == "text" || type == "image_url",
                "content array carries a part this endpoint cannot parse: " + block.ToJsonString());
        }

        // A turn holding only reasoning is now a turn worth sending: its reasoning is
        // the provider's own continuation state, so omitting it would lose the thought.
        var assistants = request["messages"]!.AsArray()
            .Count(node => node!["role"]!.GetValue<string>() == "assistant");
        Check(assistants == 3,
            $"a turn holding only reasoning should replay, saw {assistants} assistant messages: {wire}");
        Check(request["messages"]!.AsArray()
            .Where(node => node!["role"]!.GetValue<string>() == "assistant")
            .All(node => node!["content"] is JsonValue),
            "an ordinary text turn should stay a string content without media: " + wire);

        // The messages dialect keeps its own thinking block, in its own shape, so
        // chat's message-level fields must not be pushed onto that surface either.
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
