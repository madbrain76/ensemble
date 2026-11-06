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
        int cap = chat["max_completion_tokens"]!.GetValue<int>();
        // A reasoning model pays for deliberation out of this same budget. At 1024
        // the measured response was finish_reason "length" with the tool call cut off
        // mid-arguments, which the GUI can only present as an agent that never acts.
        Check(cap == 4096, $"chat completions must budget room to think and then act, got {cap}");
        // Recovery doubles this when a reply is cut off, so the base must stay low
        // enough that escalation is still available. 4096 leaves four doublings.
        Check(cap <= 8192, $"base budget {cap} leaves the output-limit ladder no headroom");

        // The same rule was not measured on the messages dialect, which has always
        // required its own value; widening it here would be an unmeasured change.
        var messages = Request(Vendor.Anthropic);
        Check(messages["max_tokens"]!.GetValue<int>() == 1024,
            "the messages dialect's cap was changed without evidence");
    }
}
