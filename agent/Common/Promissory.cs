namespace Ensemble;

// A weak model sometimes ends its turn by announcing the work it is about to do
// ("Now I'll write the regression test:") and never emits the call. The engine
// treats text without a call as the final answer, so the turn simply dies mid-task
// with no error. One bounded correction recovers almost every such turn.
public sealed partial class Agent
{
    // Two per turn: a model that keeps narrating after being told twice is telling
    // us something, and an unbounded nudge would loop forever.
    internal const int PromissoryLimit = 2;

    internal const string PromissoryCorrection = "Your last response announced work you were "
        + "about to do but made no tool call. Continue the task now: issue the tool call you "
        + "described, or give the user the complete final answer.";

    // A reply with neither prose nor a call shows the user nothing at all, so it is
    // not an answer even though the provider reported success. Ask once for the
    // missing half rather than close the turn blank.
    internal const string SilentCorrection = "Your last response was empty: it contained no "
        + "text and no tool call, so the user sees nothing. Either issue the tool call the "
        + "task needs, or answer in visible text.";

    // After the correction budget is spent, say why the turn ended instead of
    // leaving a silence that reads as a hang.
    internal const string SilentStopNotice = "The model returned an empty reply: no text and "
        + "no tool call. Nothing was done, and further correction did not help.";

    // True when an answer-only reply reads like a promise of work rather than the
    // work itself. Deliberately narrow: an ordinary final answer, however colon-
    // terminated, must not be nudged.
    internal static partial bool PromisesAction(string? answer);

    // The most recent visible agent prose, or empty when that reply showed none.
    internal partial string LastAgentAnswer();
}
