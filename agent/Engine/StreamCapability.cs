namespace Ensemble;

public sealed partial class Engine
{
    // Whether this route may be asked to stream its reply. A curated row answers
    // first, in both directions: a model documented to buffer is never asked, and a
    // documented streamer is never made to prove it again. An unlisted model is
    // nobody's knowledge but the endpoint's, so it is asked once and the answer is
    // remembered; an operator's explicit no-stream always wins over both.
    private static readonly object streamGate = new();
    private static readonly Dictionary<string, (bool Streams, DateTimeOffset At)> observedStreams = new(StringComparer.Ordinal);
    private static readonly TimeSpan streamNegativeLifetime = TimeSpan.FromMinutes(10);

    internal static bool MayStream(Configuration config)
    {
        if (config.DisableStreaming)
        {
            return false;
        }
        if (LookupModel(config.Vendor, config.Model, out var features))
        {
            return features.Stream;
        }
        lock (streamGate)
        {
            if (observedStreams.TryGetValue(StreamKey(config), out var record)
                && (!record.Streams && DateTimeOffset.UtcNow - record.At < streamNegativeLifetime))
            {
                return false;
            }
        }
        return true;
    }

    // What the endpoint actually did with the last request that asked. Only an
    // unlisted model is recorded: a curated row is application knowledge and no
    // single exchange may overrule it.
    internal static void RecordStream(Configuration config, bool streamed)
    {
        if (config.DisableStreaming || LookupModel(config.Vendor, config.Model, out _))
        {
            return;
        }
        lock (streamGate)
        {
            observedStreams[StreamKey(config)] = (streamed, DateTimeOffset.UtcNow);
        }
    }

    // Surface belongs to the key: the same model over two dialects takes different
    // request shapes, so one answer is not evidence about the other. So does the
    // offering: two bindings at one address are distinct routes, and what one of
    // them did with the flag says nothing about what the other will do.
    private static string StreamKey(Configuration config)
        => config.Vendor + "\0" + config.Surface + "\0" + config.BaseUrl.TrimEnd('/') + "\0"
            + config.Model + "\0" + config.Offering?.EndpointID;
}
