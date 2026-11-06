namespace Ensemble;

// A degraded endpoint can answer with the same character or phrase thousands of
// times over. Without a guard that stream becomes the reply: it is recorded into
// history, billed, and ends the turn, while the client waits for a model that is
// no longer saying anything. Watching is per part, and tool arguments are exempt
// because a generated file may legitimately hold long runs of separators.
public sealed partial class Engine
{
    internal const int RunawayRunes = 500;
    internal const int RunawayPhrase = 400;
    internal const int RunawayPeriod = 40;

    // True when the tail of a part is one short phrase repeating itself. A single
    // character is left to the run rule, which uses the higher bound so a drawn
    // rule or banner is tolerated.
    internal static partial bool LooksCollapsed(string tail);
}
