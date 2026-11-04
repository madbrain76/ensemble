namespace Ensemble;

public sealed partial class Engine
{
    internal static partial StreamChannel ParseStream(string? value)
    {
        string trimmed = (value ?? "").Trim();
        if (trimmed.Length == 0)
        {
            throw new ArgumentException(
                "empty streaming selection: use all or a list of text, thinking, tools");
        }
        StreamChannel bits = StreamChannel.None;
        foreach (string entry in trimmed.Split(','))
        {
            string channel = entry.Trim().ToLowerInvariant();
            if (channel == "all")
            {
                return StreamChannel.All;
            }
            bits |= channel switch
            {
                "text" => StreamChannel.Text,
                "thinking" => StreamChannel.Thinking,
                "tools" => StreamChannel.Tools,
                _ => throw new ArgumentException(
                    $"invalid streaming channel \"{entry}\": use all or a list of text, thinking, tools")
            };
        }
        return bits;
    }

    internal static partial string StreamName(StreamChannel channels) => channels switch
    {
        StreamChannel.All => "all",
        StreamChannel.None => "all",
        _ => string.Join(',',
            (channels.HasFlag(StreamChannel.Text) ? new[] { "text" } : Array.Empty<string>())
            .Concat(channels.HasFlag(StreamChannel.Thinking) ? new[] { "thinking" } : Array.Empty<string>())
            .Concat(channels.HasFlag(StreamChannel.Tools) ? new[] { "tools" } : Array.Empty<string>()))
    };

    internal static partial bool StreamSelected(StreamChannel selected, ChunkKind kind) =>
        kind switch
        {
            ChunkKind.Reply or ChunkKind.ReasoningSummary =>
                selected.HasFlag(StreamChannel.Text),
            ChunkKind.Thinking => selected.HasFlag(StreamChannel.Thinking),
            ChunkKind.ToolCall => selected.HasFlag(StreamChannel.Tools),
            _ => true
        };
}
