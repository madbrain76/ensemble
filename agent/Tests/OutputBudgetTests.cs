using System.Text.Json.Nodes;
using Ensemble;

// The output budget is the difference between an agent that acts and one that
// appears idle, so it is pinned here rather than left to a literal in a renderer.
internal static class OutputBudgetTests
{
    private static void Check(bool condition, string message)
    {
        if (!condition)
        {
            throw new InvalidOperationException(message);
        }
    }

    private static JsonObject Request(Vendor vendor)
    {
        var app = new Ensemble.Ensemble(TextWriter.Null);
        try
        {
            var agent = new Agent(app, new Configuration(vendor, "output-budget-model", "http://127.0.0.1:1"));
            var request = agent.Engine().Render(new Context([
                new Entry(1, Actor.Human, [new TextPart("question")])
            ]));
            agent.Destroy();
            return request;
        }
        finally { app.Destroy(); }
    }

    public static void Run()
    {
        var chat = Request(Vendor.OpenAI);
        long cap = chat["max_completion_tokens"]!.GetValue<long>();
        // A reasoning model pays for deliberation out of this same budget. The
        // measured failure was finish_reason "length" with the tool call cut off
        // mid-arguments, which the GUI can only present as an agent that never
        // acts, so both dialects now ask a provider for the whole ceiling.
        Check(cap == Engine.OutputBudgetCeiling,
            $"chat completions must budget room to think and then act, got {cap}");

        // The messages dialect shares that ceiling. A model configured below it
        // still doubles, so the ladder survives as a per-model decision rather
        // than a default that pays for headroom the operator did not ask for.
        var messages = Request(Vendor.Anthropic);
        Check(messages["max_tokens"]!.GetValue<long>() == Engine.OutputBudgetCeiling,
            "the messages dialect's cap was changed without evidence");
    }
}
