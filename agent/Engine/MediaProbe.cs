using System.Text.Json.Nodes;

namespace Ensemble;

// The curated table cannot name a local or newly released endpoint, so a model
// absent from it earns image capability from the endpoint itself. A solid color
// is the cheapest proof that image bytes actually arrived: a server that
// silently drops them has no way to name the color, so it fails the probe
// instead of authorizing hallucinated answers. Only the delivered answer counts.
// Reasoning that guesses a color is not evidence the image was received.
public sealed partial class Engine
{
    // One-pixel solid RGB PNGs. A random pick keeps a server that memorized one
    // answer for one model from passing every probe.
    private static readonly (string Color, string Png)[] probeColors =
    [
        ("red", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGL6z8AACAAA//8DCQECWLbVUAAAAABJRU5ErkJggg=="),
        ("green", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGJi+M8ACAAA//8CCgECX2t3gQAAAABJRU5ErkJggg=="),
        ("blue", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGJiYPgPCAAA//8BCwEC1pVDqAAAAABJRU5ErkJggg==")
    ];

    // The marker lets offline traces separate a capability probe from
    // conversation without changing what a capable model must answer.
    private const string ProbePrompt = "ensemble-media-probe Which single color fills this entire"
        + " image? Answer with exactly one word.";

    // Each result, failure included, is retained per route so neither a
    // text-only endpoint nor an unreachable server is charged again for an
    // answer that cannot change. A negative expires so one transient timeout
    // never becomes a permanent refusal; a positive needs no refresh.
    private static readonly object probeGate = new();
    private static readonly Dictionary<string, (bool Proven, DateTimeOffset At)> probedImages = new(StringComparer.Ordinal);
    private static readonly TimeSpan probeNegativeLifetime = TimeSpan.FromMinutes(10);
    private static readonly TimeSpan probeTimeout = TimeSpan.FromSeconds(30);

    // Surface is part of the route: the same endpoint and model reached over
    // two dialects takes different request shapes, so proof on one is not
    // proof on the other.
    private static string ProbeKey(Configuration config)
        => config.Vendor + "\0" + config.Surface + "\0" + config.BaseUrl.TrimEnd('/') + "\0" + config.Model;

    internal static partial bool ProvenImages(Configuration config)
    {
        if (LookupModel(config.Vendor, config.Model, out var features))
        {
            return (features.Media & Media.Image) != 0;
        }
        lock (probeGate)
        {
            return probedImages.TryGetValue(ProbeKey(config), out var record) && record.Proven;
        }
    }

    // Ask the selected endpoint to prove it receives image bytes. Table models
    // are never probed: the curated row stays authoritative in both directions.
    public partial async Task<bool> ProbeImagesAsync(CancellationToken cancellation)
    {
        Configuration config;
        try
        {
            config = Parent().Configuration();
        }
        catch (ObjectDisposedException)
        {
            return false;
        }
        if (LookupModel(config.Vendor, config.Model, out var features))
        {
            return (features.Media & Media.Image) != 0;
        }
        string key = ProbeKey(config);
        lock (probeGate)
        {
            if (probedImages.TryGetValue(key, out var record)
                && (record.Proven || DateTimeOffset.UtcNow - record.At < probeNegativeLifetime))
            {
                return record.Proven;
            }
        }
        bool proven;
        try
        {
            // Two different colors, asked at once. A single probe is only a
            // one-in-three test of an endpoint that never sees the bytes and
            // answers from a fixed word, and a positive is retained. Requiring
            // both colors, in parallel so the first image send still waits for
            // only one round trip, makes a memorized answer fail.
            var asked = TwoProbeColors();
            var answers = await Task.WhenAll(ProbeOnceAsync(config, asked[0], cancellation),
                ProbeOnceAsync(config, asked[1], cancellation)).ConfigureAwait(false);
            proven = answers[0] && answers[1];
        }
        catch (OperationCanceledException) when (cancellation.IsCancellationRequested)
        {
            // The caller's own cancellation says nothing about the model, so it
            // must not be cached as a negative.
            throw;
        }
        catch (ObjectDisposedException)
        {
            // Our own teardown is not evidence about the endpoint, so it must
            // not become the route's cached answer.
            return false;
        }
        catch (Exception)
        {
            // Any transport, status, authentication, parsing or answer failure
            // leaves images unsupported. Fail-closed is the whole contract.
            proven = false;
        }
        lock (probeGate)
        {
            probedImages[key] = (proven, DateTimeOffset.UtcNow);
        }
        if (proven)
        {
            try
            {
                Parent().Log("Media probe: model " + config.Model + " accepts image input.");
            }
            catch (ObjectDisposedException)
            {
                // A proven route stays proven even when its Agent is gone; the
                // log line is not worth an unobserved background fault.
            }
        }
        return proven;
    }

    // Two distinct colors in an unpredictable order. Distinct is the property
    // that matters: one fixed answer cannot satisfy both questions.
    private static (string Color, string Png)[] TwoProbeColors()
    {
        int first = Random.Shared.Next(probeColors.Length);
        int second = (first + 1 + Random.Shared.Next(probeColors.Length - 1)) % probeColors.Length;
        return [probeColors[first], probeColors[second]];
    }

    private async Task<bool> ProbeOnceAsync(Configuration config, (string Color, string Png) color,
        CancellationToken cancellation)
    {
        using var bound = CancellationTokenSource.CreateLinkedTokenSource(cancellation);
        bound.CancelAfter(probeTimeout);
        // The request is fully described by the frozen configuration value it
        // was started with, so it takes no Engine-send guard. A capability
        // probe must never be the reason a turn or a model selection waits.
        // SendCore owns authentication, same-source plan retry, parsing and
        // auxiliary accounting exactly as the recall judge does; the probe adds
        // no parallel transport, and an unlisted model renders non-streaming.
        var parsed = await SendCoreAsync(config, ProbeRequest(config, color.Png).ToJsonString(),
            bound.Token, onChunk: null, auxiliary: true).ConfigureAwait(false);
        string answer = string.Join(" ", parsed.Events
            .SelectMany(entry => entry.Parts)
            .OfType<TextPart>()
            .Select(part => part.Text)).ToLowerInvariant();
        // Naming the asked color is necessary but not sufficient: a reply that
        // also names the other two is a hedge, not an observation.
        return answer.Contains(color.Color, StringComparison.Ordinal)
            && probeColors.Where(entry => entry.Color != color.Color)
                .All(entry => !answer.Contains(entry.Color, StringComparison.Ordinal));
    }

    // Hand-built request shapes: the renderers would apply the very capability
    // rule that this probe exists to question.
    private static JsonObject ProbeRequest(Configuration config, string png)
    {
        string data = "data:image/png;base64," + png;
        if (config.Vendor == Vendor.Anthropic)
        {
            return new JsonObject
            {
                ["model"] = config.Model,
                ["max_tokens"] = 512,
                ["messages"] = new JsonArray(new JsonObject
                {
                    ["role"] = "user",
                    ["content"] = new JsonArray(
                        new JsonObject
                        {
                            ["type"] = "image",
                            ["source"] = new JsonObject
                            {
                                ["type"] = "base64", ["media_type"] = "image/png", ["data"] = png
                            }
                        },
                        new JsonObject { ["type"] = "text", ["text"] = ProbePrompt })
                })
            };
        }
        if (config.Vendor == Vendor.Gemini)
        {
            return new JsonObject
            {
                ["contents"] = new JsonArray(new JsonObject
                {
                    ["role"] = "user",
                    ["parts"] = new JsonArray(
                        new JsonObject
                        {
                            ["inline_data"] = new JsonObject { ["mime_type"] = "image/png", ["data"] = png }
                        },
                        new JsonObject { ["text"] = ProbePrompt })
                }),
                ["generationConfig"] = new JsonObject { ["maxOutputTokens"] = 512 }
            };
        }
        if (config.Surface == Surface.Responses)
        {
            return new JsonObject
            {
                ["model"] = config.Model,
                ["max_output_tokens"] = 512,
                ["input"] = new JsonArray(new JsonObject
                {
                    ["role"] = "user",
                    ["content"] = new JsonArray(
                        new JsonObject { ["type"] = "input_text", ["text"] = ProbePrompt },
                        new JsonObject { ["type"] = "input_image", ["image_url"] = data })
                })
            };
        }
        return new JsonObject
        {
            ["model"] = config.Model,
            ["max_tokens"] = 512,
            ["messages"] = new JsonArray(new JsonObject
            {
                ["role"] = "user",
                ["content"] = new JsonArray(
                    new JsonObject { ["type"] = "text", ["text"] = ProbePrompt },
                    new JsonObject { ["type"] = "image_url", ["image_url"] = new JsonObject { ["url"] = data } })
            })
        };
    }
}
