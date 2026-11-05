using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;

namespace Ensemble;

public sealed partial class Engine
{
    // Deliberately trivial: the probe asks only whether the server accepts the
    // parameter value, so it needs no content and the smallest completion.
    private const string EffortProbePrompt = "ok";

    public partial async Task ValidateReasoningEffortAsync(CancellationToken cancellation,
        TimeSpan? probeTimeout)
    {
        Configuration config = parent!.Configuration();
        string effort = config.ReasoningEffort.Trim();
        if (effort.Length == 0)
        {
            return;
        }
        // Only this surface's request shape and listing document are verified here.
        // Probing another vendor with a body it does not define would report
        // "unverified" for the wrong reason, so the value passes through instead.
        if (config.Vendor != Vendor.OpenAI || config.Billing != BillingMode.ApiKey)
        {
            parent!.Parent().Log($"reasoning effort \"{effort}\" sent unverified: this check "
                + "covers the OpenAI API-key surfaces only");
            return;
        }
        (string? modelList, string? refusal, string? silence) = await ProbeAsync(config, effort,
            cancellation, probeTimeout ?? EffortProbeTimeout);
        if (refusal is not null)
        {
            // The server rejected this exact parameter with its own words. Silently
            // continuing would advertise a capability the provider denied.
            throw new InvalidOperationException(
                $"provider refused reasoning effort \"{effort}\": {refusal}");
        }
        if (modelList is not null)
        {
            parent!.Parent().Log($"reasoning effort \"{effort}\" accepted by the provider");
        }
        else if (refusal is null)
        {
            // Nothing denied the value, so the launch proceeds. It proceeds
            // honestly too: the two silences mean different things to an operator
            // diagnosing a launch, so the log says which one happened.
            string detail = silence ?? "the endpoint answered the request without naming "
                + "the effort levels it supports";
            parent!.Parent().Log($"reasoning effort \"{effort}\" sent unverified: {detail}");
        }
    }

    // Returns the listing document that confirmed the value, or the server's own
    // refusal message. Both null means the endpoint gave no usable answer; silence
    // then explains which kind, because a busy provider and an uninformative one
    // need different follow-ups.
    private async Task<(string? Advertised, string? Refusal, string? Silence)> ProbeAsync(
        Configuration config, string effort, CancellationToken cancellation, TimeSpan budget)
    {
        using var timeout = CancellationTokenSource.CreateLinkedTokenSource(cancellation);
        timeout.CancelAfter(budget);
        try
        {
            string? advertised = await AdvertiseAsync(config, effort, timeout.Token);
            // An advertisement is only a hint: at least one server publishes a strict
            // subset of what it accepts, so a missing name must still get to answer.
            if (advertised is not null)
            {
                return (advertised, null, null);
            }
            return (null, await AcceptsAsync(config, effort, timeout.Token), null);
        }
        catch (OperationCanceledException) when (!cancellation.IsCancellationRequested)
        {
            // Our own deadline, not a shutdown. A provider that is busy has not
            // refused anything, so this must read exactly like an endpoint that
            // answered nothing. Aborting here turned a slow LAN server into a GUI
            // that never started, with an OperationCanceledException as its only
            // trace.
            return (null, null, $"the endpoint did not answer within "
                + $"{budget.TotalSeconds:0.##} seconds");
        }
    }

    private async Task<string?> AdvertiseAsync(Configuration config, string effort,
        CancellationToken cancellation)
    {
        try
        {
            using var request = new HttpRequestMessage(HttpMethod.Get, ModelsUrl(config));
            Authorize(request, config);
            using var response = await http.SendAsync(request, cancellation);
            if (!response.IsSuccessStatusCode)
            {
                return null;
            }
            using var document = JsonDocument.Parse(await response.Content.ReadAsStringAsync(cancellation));
            var entries = new List<(string Id, string[] Efforts)>();
            foreach (string container in new[] { "data", "models" })
            {
                // OpenAI-compatible servers publish entries under either key; reading
                // both keeps one shape from silently masking the other.
                if (!document.RootElement.TryGetProperty(container, out var list)
                    || list.ValueKind != JsonValueKind.Array)
                {
                    continue;
                }
                foreach (var entry in list.EnumerateArray())
                {
                    if (entry.ValueKind != JsonValueKind.Object
                        || !entry.TryGetProperty("id", out var id)
                        || id.ValueKind != JsonValueKind.String)
                    {
                        continue;
                    }
                    string[] levels = entry.TryGetProperty("supported_reasoning_efforts",
                        out var supported) && supported.ValueKind == JsonValueKind.Array
                        ? supported.EnumerateArray().Where(value =>
                            value.ValueKind == JsonValueKind.String).Select(value => value.GetString()!)
                          .ToArray()
                        : [];
                    entries.Add((id.GetString()!, levels));
                }
            }
            var match = entries.FirstOrDefault(candidate =>
                string.Equals(candidate.Id, config.Model, StringComparison.Ordinal));
            if (match.Efforts.Length == 0)
            {
                return null;
            }
            return match.Efforts.Any(candidate =>
                string.Equals(candidate, effort, StringComparison.Ordinal)) ? match.Id : null;
        }
        catch (Exception) when (cancellation.IsCancellationRequested)
        {
            throw new OperationCanceledException(cancellation);
        }
        catch (Exception)
        {
            // A listing the client cannot read says nothing about effort levels.
            return null;
        }
    }

    private async Task<string?> AcceptsAsync(Configuration config, string effort,
        CancellationToken cancellation)
    {
        // The smallest request that still carries the field under test. Responses
        // nests the level inside a reasoning object; chat-completions takes it at
        // the top level, matching each surface's documented request shape.
        Dictionary<string, object> body = config.Surface == Surface.Responses
            ? new Dictionary<string, object>
            {
                ["model"] = config.Model,
                ["input"] = new object[] { new Dictionary<string, object>
                    { ["role"] = "user", ["content"] = EffortProbePrompt } },
                ["reasoning"] = new Dictionary<string, object> { ["effort"] = effort },
                ["max_output_tokens"] = 1,
                ["stream"] = false
            }
            : new Dictionary<string, object>
            {
                ["model"] = config.Model,
                ["messages"] = new object[] { new Dictionary<string, object>
                    { ["role"] = "user", ["content"] = EffortProbePrompt } },
                ["reasoning_effort"] = effort,
                ["max_tokens"] = 1,
                ["stream"] = false
            };
        try
        {
            using var request = new HttpRequestMessage(HttpMethod.Post,
                config.Surface == Surface.Responses ? ResponsesUrl(config) : ChatUrl(config));
            Authorize(request, config);
            request.Content = new StringContent(JsonSerializer.Serialize(body),
                Encoding.UTF8, "application/json");
            using var response = await http.SendAsync(request, cancellation);
            string text = await response.Content.ReadAsStringAsync(cancellation);
            if (response.IsSuccessStatusCode)
            {
                return null; // Accepted: nothing to report.
            }
            string message = RefusalMessage(text);
            // Only an explicit rejection of this specific parameter aborts the
            // launch; a login page, a rate limit or an unrelated 4xx is not evidence.
            return message.Contains("reasoning_effort", StringComparison.OrdinalIgnoreCase)
                ? message : null;
        }
        catch (Exception) when (cancellation.IsCancellationRequested)
        {
            throw new OperationCanceledException(cancellation);
        }
        catch (Exception)
        {
            return null; // Unreachable is not the same as refused.
        }
    }

    private static string RefusalMessage(string text)
    {
        try
        {
            using var document = JsonDocument.Parse(text);
            if (document.RootElement.TryGetProperty("error", out var error))
            {
                if (error.TryGetProperty("message", out var message)
                    && message.ValueKind == JsonValueKind.String)
                {
                    return message.GetString() ?? "";
                }
                if (error.ValueKind == JsonValueKind.String)
                {
                    return error.GetString() ?? "";
                }
            }
        }
        catch (JsonException)
        {
        }
        return text.Length <= 300 ? text : text[..300];
    }

    private static void Authorize(HttpRequestMessage request, Configuration config)
    {
        if (!string.IsNullOrEmpty(config.ApiKey))
        {
            request.Headers.TryAddWithoutValidation("Authorization", "Bearer " + config.ApiKey);
        }
    }

    private static string ModelsUrl(Configuration config)
        => Trim(config) + "/models";

    private static string ChatUrl(Configuration config)
        => Trim(config) + "/chat/completions";

    private static string ResponsesUrl(Configuration config)
        => Trim(config) + "/responses";

    // Accept either a service base or the complete documented endpoint, the same
    // latitude the request path allows.
    private static string Trim(Configuration config)
    {
        string url = config.BaseUrl.TrimEnd('/');
        foreach (string suffix in new[] { "/chat/completions", "/responses", "/v1" })
        {
            if (url.EndsWith(suffix, StringComparison.OrdinalIgnoreCase))
            {
                url = url[..^suffix.Length];
            }
        }
        return url.EndsWith("/v1", StringComparison.OrdinalIgnoreCase) ? url : url + "/v1";
    }
}
