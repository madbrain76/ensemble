using Ensemble;

// A stream that stops saying anything new must not become the answer.
static class RunawayGuardTests
{
    static void Check(bool value, string message)
    {
        if (!value) throw new Exception(message);
    }

    static bool Feed(Engine.RunawayGuard guard, string text, ChunkKind kind = ChunkKind.Reply,
        int partId = 1)
        => guard.Observe(partId, kind, text);

    public static void Run()
    {
        var bangs = new Engine.RunawayGuard();
        Check(!Feed(bangs, new string('!', 300)), "a short run of one character tripped the guard");
        Check(Feed(new Engine.RunawayGuard(), new string('!', 520)),
            "500+ identical characters were not recognised as collapse");
        Check(!Feed(new Engine.RunawayGuard(), new string(' ', 900)),
            "whitespace padding, as in a wide aligned table, tripped the guard");
        Check(!Feed(new Engine.RunawayGuard(), new string('-', 480)),
            "a drawn rule just under the bound was rejected");
        var phrase = new Engine.RunawayGuard();
        Check(Feed(phrase, string.Concat(Enumerable.Repeat("The user", 80))),
            "a repeated phrase was not recognised as collapse");
        Check(!Feed(new Engine.RunawayGuard(),
            string.Concat(Enumerable.Repeat("The user", 12)) + " reviewed the report and approved it."),
            "ordinary prose that repeats a phrase a little was rejected");
        Check(!Feed(new Engine.RunawayGuard(), string.Concat(Enumerable.Repeat(
            "| RS-10-01    | restore      | a recovery point stored on NTFS |\n", 30))),
            "a genuinely repetitive markdown table was rejected");
        // Tool arguments carry long runs legitimately.
        Check(!Feed(new Engine.RunawayGuard(), new string('x', 4000), ChunkKind.ToolCall),
            "tool arguments were treated as a collapsed reply");
        // Each part is judged on its own.
        var mixed = new Engine.RunawayGuard();
        Check(Feed(mixed, new string('#', 600), partId: 7), "the flooding part was not caught");
        Check(!Feed(mixed, "plain prose in another part", partId: 8),
            "a clean part was condemned by another part's flood");
        // Chunks arrive in pieces; the run must survive the boundaries.
        var split = new Engine.RunawayGuard();
        bool tripped = false;
        for (int index = 0; index < 60; index++)
        {
            tripped |= Feed(split, new string('?', 12));
        }
        Check(tripped, "a flood delivered in small chunks was not detected");
    }
}
