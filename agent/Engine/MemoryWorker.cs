using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

public sealed partial class Engine
{
    private static partial bool IncompleteOutput(Configuration config, string body)
    {
        if (config.Surface == Surface.Responses)
        {
            // The Responses parser already requires terminal completed success.
            return false;
        }
        var root = JsonNode.Parse(body)!.AsObject();
        return config.Vendor switch
        {
            Vendor.OpenAI => root["choices"]?[0]?["finish_reason"]?.GetValue<string>()
                is not ("stop" or "tool_calls"),
            Vendor.Anthropic => root["stop_reason"]?.GetValue<string>()
                is not ("end_turn" or "tool_use"),
            Vendor.Gemini => root["candidates"]?[0]?["finishReason"]?.GetValue<string>() != "STOP",
            _ => true
        };
    }

    private static partial bool OutputCapped(Configuration config, string body)
    {
        JsonNode? node;
        try
        {
            node = JsonNode.Parse(body);
        }
        catch (JsonException)
        {
            return false;
        }
        if (node is not JsonObject root)
        {
            return false;
        }
        return config.Vendor switch
        {
            // The names differ per vendor; only these say "the budget ran out".
            Vendor.OpenAI => root["choices"]?[0]?["finish_reason"]?.GetValue<string>()
                is "length" or "max_tokens" || root["stop_reason"]?.GetValue<string>() == "max_tokens",
            Vendor.Anthropic => root["stop_reason"]?.GetValue<string>() == "max_tokens",
            Vendor.Gemini => root["candidates"]?[0]?["finishReason"]?.GetValue<string>() == "MAX_TOKENS",
            _ => false
        };
    }
}
