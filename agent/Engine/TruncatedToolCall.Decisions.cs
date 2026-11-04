using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

public sealed partial class Engine
{
    internal static partial bool ToolCallCapped(Configuration config, string completedJson,
        IEnumerable<Event> events)
    {
        if (completedJson.Length == 0 || !SaysLength(config, completedJson))
        {
            return false;
        }
        foreach (var fact in events)
        {
            if (fact.Kind != EventKind.ResponseEnded)
            {
                continue;
            }
            foreach (var part in fact.Parts)
            {
                if (part is ToolCallPart)
                {
                    return true;
                }
            }
        }
        return false;
    }

    internal static partial bool SaysLength(Configuration config, string completedJson)
    {
        JsonNode? node;
        try
        {
            node = JsonNode.Parse(completedJson);
        }
        catch (JsonException)
        {
            return false;
        }
        if (node is not JsonObject root)
        {
            return false;
        }
        // Only these vendor values say the budget ran out. Anything else, including
        // a missing reason, must not spend a retry.
        return config.Vendor switch
        {
            Vendor.OpenAI => root["choices"]?[0]?["finish_reason"]?.GetValue<string>()
                is "length" or "max_tokens" || root["stop_reason"]?.GetValue<string>() == "max_tokens",
            Vendor.Anthropic => root["stop_reason"]?.GetValue<string>() == "max_tokens",
            Vendor.Gemini => root["candidates"]?[0]?["finishReason"]?.GetValue<string>() == "MAX_TOKENS",
            _ => false
        };
    }

    internal static partial bool LargerOutputBudget(string body, out string retry)
    {
        retry = body;
        JsonNode? root;
        try
        {
            root = JsonNode.Parse(body);
        }
        catch (JsonException)
        {
            return false;
        }
        if (root is not JsonObject request)
        {
            return false;
        }
        foreach (string field in OutputCapFields)
        {
            if (request[field] is not JsonNode value
                || value.GetValueKind() != JsonValueKind.Number)
            {
                continue;
            }
            int cap;
            try
            {
                cap = value.GetValue<int>();
            }
            catch (Exception notANumber) when (notANumber is InvalidOperationException
                or FormatException or OverflowException)
            {
                continue;
            }
            if (cap <= 0 || cap >= OutputBudgetCeiling)
            {
                continue;
            }
            request[field] = Math.Min(cap * 2, OutputBudgetCeiling);
            retry = request.ToJsonString();
            return true;
        }
        return false;
    }
}
