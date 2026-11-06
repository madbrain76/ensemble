using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

// A call whose arguments were cut off by the output budget is not a call. Partial
// JSON must never be dispatched, so the whole attempt is repeated at a larger cap.
public sealed partial class Engine
{
    // The cap fields, ceiling and LargerOutputBudget used to grow a budget live with
    // the output-limit ladder in OutputLimit.cs. This surface needed the same
    // doubling independently, so the two copies were compared and one was kept:
    // a truncated call and a truncated answer must escalate identically or the
    // recovery depends on which kind of reply happened to be cut off.

    // True when the provider said the budget ran out and the reply was reaching for
    // a tool call at that moment.
    internal static partial bool ToolCallCapped(Configuration config, string completedJson,
        IEnumerable<Event> events);

    internal static partial bool SaysLength(Configuration config, string completedJson);
}
