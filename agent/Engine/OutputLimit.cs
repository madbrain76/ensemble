using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

public sealed partial class Engine
{
    // The thinking model that spent 4,096 tokens deliberating and then stopped has
    // not failed to answer; it ran out of room to answer in. Asking again at a
    // larger cap often suffices, and if it does not, saying so explicitly stops the
    // model from deliberating again. The budget stops doubling at 16384.
    internal const int OutputBudgetCeiling = 16384;
    internal static readonly string[] OutputCapFields =
        ["max_tokens", "max_completion_tokens", "max_output_tokens"];

    private const string OutputDirective = "Your previous attempt used the entire output budget on "
        + "private deliberation and ended without a visible answer. Do not deliberate further: reply "
        + "now with the complete final answer to the last request in visible text.";

    // A response that carried only deliberation: no visible prose and no tool call.
    // A truncated answer that did publish text is deliberately excluded, because
    // rewriting an answer the operator already watched would silently replace it.
    internal static partial bool EndsBeforeAnswer(ParsedResponse parsed);

    // The ordered attempts a truncated exchange deserves, or none. An empty list is
    // the ordinary path: no extra request, no cost, no delay.
    internal static partial List<string> RecoveryBodies(string body, bool published,
        bool cancelled);

    internal static partial bool LargerOutputBudget(string body, out string retry);

    internal static partial bool WithAnswerDirective(string body, out string retry);
}
