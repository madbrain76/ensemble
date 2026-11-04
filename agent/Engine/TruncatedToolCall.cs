using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

// A call whose arguments were cut off by the output budget is not a call. Partial
// JSON must never be dispatched, so the whole attempt is repeated at a larger cap.
public sealed partial class Engine
{
    internal const int OutputBudgetCeiling = 16384;
    internal static readonly string[] OutputCapFields =
        ["max_tokens", "max_completion_tokens", "max_output_tokens"];

    // True when the provider said the budget ran out and the reply was reaching for
    // a tool call at that moment.
    internal static partial bool ToolCallCapped(Configuration config, string completedJson,
        IEnumerable<Event> events);

    internal static partial bool LargerOutputBudget(string body, out string retry);

    internal static partial bool SaysLength(Configuration config, string completedJson);
}
