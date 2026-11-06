namespace Ensemble;

/// <summary>Which delivery channels may stream live to the caller.</summary>
[Flags]
public enum StreamChannel
{
    None = 0,
    Text = 1,
    Thinking = 2,
    Tools = 4,
    All = Text | Thinking | Tools
}

public sealed partial class Engine
{
    // Accept "all", or a comma-separated list of text, thinking and tools. The keys
    // mirror the flags so one channel can be requested at a time. An empty selection
    // is refused rather than silently meaning "watch nothing", which would look like
    // a stalled model.
    internal static partial StreamChannel ParseStream(string? value);

    // The single decision the streaming path asks about.
    internal static partial bool StreamSelected(StreamChannel selected, ChunkKind kind);

    // Round-trip a selection back to its accepted spelling, so a derived
    // configuration keeps the operator's choice instead of falling back to "all".
    internal static partial string StreamName(StreamChannel channels);
}
