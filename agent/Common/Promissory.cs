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

    // True when an answer-only reply reads like a promise of work rather than the
    // work itself. Deliberately narrow: an ordinary final answer, however colon-
    // terminated, must not be nudged.
    internal static partial bool PromisesAction(string? answer);

    // The most recent visible agent prose, or empty when that reply showed none.
    internal partial string LastAgentAnswer();
}
