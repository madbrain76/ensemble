using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

public sealed partial class Engine
{
    internal static partial bool EndsBeforeAnswer(ParsedResponse parsed)
    {
        // The provider must have said the budget ran out. A model that answered
        // empty, or was cut off by a filter, is not this condition, and retrying it
        // would double the cost of every such exchange.
        if (!parsed.OutputCapped)
        {
            return false;
        }
        Event? ended = null;
        foreach (var candidate in parsed.Events)
        {
            if (candidate.Kind == EventKind.ResponseEnded)
            {
                ended = candidate;
            }
        }
        if (ended is null)
        {
            return false; // No response at all is a transport failure, not a cap.
        }
        foreach (var part in ended.Parts)
        {
            switch (part)
            {
                case TextPart text when !string.IsNullOrWhiteSpace(text.Text):
                case ToolCallPart:
                    return false;
                case TextPart:
                    break;
            }
        }
        return true;
    }

    internal static partial List<string> RecoveryBodies(string body, bool published,
        bool cancelled)
    {
        var attempts = new List<string>();
        // Retrying is only safe while nothing was published. A partially delivered
        // answer is the operator's already, and a cancellation must not spend more
        // provider requests on work nobody is waiting for.
        if (published || cancelled)
        {
            return attempts;
        }
        string current = body;
        if (LargerOutputBudget(body, out string doubled))
        {
            attempts.Add(doubled);
            current = doubled;
        }
        if (WithAnswerDirective(current, out string directed))
        {
            attempts.Add(directed);
        }
        return attempts;
    }

    internal static partial bool LargerOutputBudget(string body, out string retry)
    {
        retry = body;
        JsonNode? root = TryParse(body);
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

    internal static partial bool WithAnswerDirective(string body, out string retry)
    {
        retry = body;
        JsonNode? root = TryParse(body);
        if (root is not JsonObject request)
        {
            return false;
        }
        foreach (string field in new[] { "messages", "input", "contents" })
        {
            if (request[field] is not JsonArray turns)
            {
                continue;
            }
            // Gemini takes parts, the other two take a plain content string.
            turns.Add(field == "contents"
                ? new JsonObject { ["role"] = "user", ["parts"] = new JsonArray
                    { new JsonObject { ["text"] = OutputDirective } } }
                : new JsonObject { ["role"] = "user", ["content"] = OutputDirective });
            retry = request.ToJsonString();
            return true;
        }
        return false;
    }

    private static JsonNode? TryParse(string body)
    {
        try
        {
            return JsonNode.Parse(body);
        }
        catch (JsonException)
        {
            return null; // A body this path cannot read is left exactly as it was.
        }
    }
}
