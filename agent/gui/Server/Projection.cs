using Ensemble;
using System.Collections.Immutable;
using System.Text.Json.Nodes;

namespace EnsembleGui;

public static partial class Gui
{
    public static partial string PartID(string agentID, long responseID, int partID)
    {
        return $"{agentID}:r{responseID}:p{partID}";
    }

    public static partial JsonObject PartFrame(string agentID, string id, Part part, Actor role,
        long sequence, DateTimeOffset? recordedAt, PartSpan? span)
    {
        var wire = JsonCodec.WritePart(part);
        wire["role"] = role == Actor.Agent ? "assistant" : role.ToString().ToLowerInvariant();
        // History retains the signed native value for compatible replay. The
        // browser needs only displayable text, never encrypted replay material.
        // A Responses summary describes reasoning; it is not the private trace.
        if (part is OpaquePart opaque)
        {
            var value = opaque.Data;
            bool summary = opaque.From.Surface == Surface.Responses;
            wire = new JsonObject
            {
                ["type"] = summary ? "reasoning_summary" : "thinking",
                ["role"] = role == Actor.Agent ? "assistant" : role.ToString().ToLowerInvariant()
            };
            string text = "";
            if (value.ValueKind == System.Text.Json.JsonValueKind.Object)
            {
                if (summary && value.TryGetProperty("summary", out var blocks)
                    && blocks.ValueKind == System.Text.Json.JsonValueKind.Array)
                {
                    // The parser validates this shape. Keep the projection safe
                    // for historical values too, without exposing unknown fields.
                    text = string.Concat(blocks.EnumerateArray()
                        .Where(block => block.ValueKind == System.Text.Json.JsonValueKind.Object
                            && block.TryGetProperty("type", out var kind)
                            && kind.ValueKind == System.Text.Json.JsonValueKind.String
                            && kind.GetString() == "summary_text"
                            && block.TryGetProperty("text", out var content)
                            && content.ValueKind == System.Text.Json.JsonValueKind.String)
                        .Select(block => block.GetProperty("text").GetString()));
                }
                foreach (var key in summary ? Array.Empty<string>() : new[] { "thinking", "text", "reasoning_content" })
                {
                    if (value.TryGetProperty(key, out var field) && field.ValueKind == System.Text.Json.JsonValueKind.String)
                    {
                        text = field.GetString() ?? "";
                        break;
                    }
                }
            }
            wire["text"] = text;
        }
        if (span is not null && span.Started > 0)
        {
            // A recorded span survives reload intact. A part that never streamed
            // carries neither key, so the view shows a clock without a duration
            // rather than a production time of zero.
            wire["started"] = JsonValue.Create(span.Started);
            wire["elapsed"] = JsonValue.Create(span.Elapsed);
        }
        return new JsonObject
        {
            ["type"] = "part_final", ["agent_id"] = agentID, ["part_id"] = id,
            ["seq"] = sequence, ["part"] = wire, ["recorded_at"] = recordedAt?.ToString("O")
        };
    }

    private static partial ImmutableArray<Event> DisplayFacts(HistorySnapshot snapshot)
    {
        // Context, not the diagnostic tail, owns the finalized transcript. This
        // works equally after snapshot-only load and after a later redaction.
        // The transcript is authoritative for content only. Each reconstructed
        // fact takes its wall time from the durable event that recorded it, so a
        // reload never stamps history with page-load time.
        var facts = snapshot.Events.Where(fact => fact.Kind is not
            (EventKind.MessageReceived or EventKind.HintReceived or EventKind.ResponseEnded or EventKind.ToolReturned))
            .ToList();
        foreach (var entry in snapshot.Context.Dialogue)
        {
            if (entry.Parts.Any(part => part is SkillPart))
            {
                continue;
            }
            if (entry.Parts is [ToolResultPart result])
            {
                facts.Add(new Event(EventKind.ToolReturned, seq: entry.Seq, time: entry.Time, parts: result.Parts,
                    callID: result.CallID, isError: result.IsError, job: result.Job,
                    toolReferences: result.ToolReferences, observationDuration: entry.ObservationDuration, workProduct: result.WorkProduct,
                    directoryListing: result.DirectoryListing, fileEvidence: result.FileEvidence, searchEvidence: result.SearchEvidence));
            }
            else if (entry.Actor == Actor.Agent)
            {
                // The dialogue entry keeps content; the durable event keeps the
                // per-part production spans, which only the streaming path observes.
                facts.Add(new Event(EventKind.ResponseEnded, seq: entry.Seq, time: entry.Time, parts: entry.Parts,
                    responseID: entry.ResponseID, partIDs: entry.PartIDs,
                    partSpans: snapshot.Events.FirstOrDefault(fact => fact.Seq == entry.Seq)?.PartSpans));
            }
            else
            {
                facts.Add(new Event(EventKind.MessageReceived, seq: entry.Seq, time: entry.Time, actor: entry.Actor, parts: entry.Parts));
            }
        }
        return facts.OrderBy(fact => fact.Seq).ToImmutableArray();
    }

    public static partial Projection Project(HistorySnapshot snapshot, string agentID)
    {
        var frames = ImmutableArray.CreateBuilder<JsonObject>();
        var ids = ImmutableHashSet.CreateBuilder<string>();
        var calls = new Dictionary<string, string>(StringComparer.Ordinal);
        // A Job's continuing view is keyed by the tool that owns it, so the name
        // must survive a reload even when the model's own call part has left the
        // working Context. The durable dispatch fact is the authority for it.
        var names = new Dictionary<string, string>(StringComparer.Ordinal);
        string state = "idle";
        string lastAnswer = "";
        bool paused = false;
        bool shutdown = false;
        var context = new Context();
        var display = DisplayFacts(snapshot);
        foreach (var fact in display)
        {
            if (fact.Kind is EventKind.ToolCalled or EventKind.ToolStarted
                && fact.CallID.Length > 0 && fact.Name.Length > 0)
            {
                names[fact.CallID] = fact.Name;
            }
        }
        foreach (var fact in display)
        {
            // DisplayFacts combines the authoritative retained transcript with
            // diagnostics, not a complete replayable History. Reduce only turn
            // phase inputs here: memory transactions and context cuts may refer
            // to source that the snapshot has already retired.
            if (fact.Kind is EventKind.RuntimeConfigured or EventKind.MessageReceived
                or EventKind.RequestSent or EventKind.ResponseEnded or EventKind.ToolReturned
                or EventKind.TurnInterrupted or EventKind.TurnEnded or EventKind.ErrorOccurred)
            {
                context = History.Reduce(context, fact);
            }
            if (fact.Kind == EventKind.RuntimeConfigured)
            {
                paused = false;
                shutdown = false;
            }
            if (fact.Kind == EventKind.AgentPaused)
            {
                paused = true;
            }
            else if (fact.Kind == EventKind.AgentResumed)
            {
                paused = false;
            }
            if (fact.Kind == EventKind.AgentShutdown)
            {
                shutdown = true;
            }
            string next = shutdown ? "shutdown"
                : paused ? "paused" : context.Phase.ToString().ToLowerInvariant();
            if (next != state)
            {
                frames.Add(new JsonObject
                {
                    ["type"] = "state_changed", ["agent_id"] = agentID, ["seq"] = fact.Seq,
                    ["from"] = state, ["to"] = next
                });
                state = next;
            }
            if (fact.Kind is EventKind.MessageReceived or EventKind.HintReceived or EventKind.ResponseEnded)
            {
                if (fact.Kind == EventKind.ResponseEnded)
                {
                    lastAnswer = string.Concat(fact.Parts.OfType<TextPart>().Select(part => part.Text));
                }
                for (int i = 0; i < fact.Parts.Length; i++)
                {
                    // Historical logs have no parser IDs. Their event sequence
                    // and offset give a stable, disjoint replay-only namespace.
                    string id = fact.ResponseID == 0
                        ? $"{agentID}:e{fact.Seq}:p{i}"
                        : PartID(agentID, fact.ResponseID, fact.PartIDs[i]);
                    var role = fact.Kind == EventKind.ResponseEnded ? Actor.Agent : fact.Actor ?? Actor.Human;
                    PartSpan? span = null;
                    if (fact.Kind == EventKind.ResponseEnded && !fact.PartSpans.IsDefaultOrEmpty
                        && i < fact.PartSpans.Length)
                    {
                        span = fact.PartSpans[i];
                    }
                    frames.Add(PartFrame(agentID, id, fact.Parts[i], role, fact.Seq, fact.Time, span));
                    ids.Add(id);
                    if (fact.Parts[i] is ToolCallPart call)
                    {
                        calls[call.CallID] = id;
                    }
                }
            }
            else if (fact.Kind == EventKind.ToolStarted)
            {
                // Dispatch diagnostics must not recreate model calls removed
                // from the authoritative working Context by a recorded cut.
                if (!calls.TryGetValue(fact.CallID, out string? id))
                {
                    continue;
                }
                frames.Add(new JsonObject
                {
                    ["type"] = "tool_dispatched", ["agent_id"] = agentID, ["seq"] = fact.Seq,
                    ["recorded_at"] = fact.Time?.ToString("O"),
                    ["part_id"] = id, ["call_id"] = fact.CallID, ["name"] = fact.Name,
                    ["input"] = JsonNode.Parse(fact.Args!.Value.GetRawText()),
                    ["job"] = fact.Job is null ? null : JsonNode.Parse(JsonCodec.WriteEvent(fact))?["tool"]?["job"]?.DeepClone()
                });
            }
            else if (fact.Kind is EventKind.ToolReturned or EventKind.ToolObserved)
            {
                string id = $"{agentID}:e{fact.Seq}:result";
                frames.Add(new JsonObject
                {
                    ["type"] = "tool_finished", ["agent_id"] = agentID, ["seq"] = fact.Seq,
                    ["observation_ms"] = fact.ObservationDuration?.TotalMilliseconds,
                    ["work_product"] = fact.WorkProduct is null ? null : new JsonObject
                    {
                        ["id"] = fact.WorkProduct.ID, ["name"] = fact.WorkProduct.Name,
                        ["offered_bytes"] = fact.WorkProduct.OfferedBytes
                    },
                    ["recorded_at"] = fact.Time?.ToString("O"),
                    ["part_id"] = id, ["tool_part_id"] = calls.GetValueOrDefault(fact.CallID),
                    ["call_id"] = fact.CallID, ["result"] = JsonCodec.Text(fact.Parts),
                    ["name"] = names.GetValueOrDefault(fact.CallID),
                    ["search_evidence"] = fact.SearchEvidence is null ? null
                        : JsonNode.Parse(JsonCodec.WriteEvent(fact))?["tool"]?["search_evidence"]?.DeepClone(),
                    ["file_evidence"] = fact.FileEvidence is null ? null
                        : JsonNode.Parse(JsonCodec.WriteEvent(fact))?["tool"]?["file_evidence"]?.DeepClone(),
                    ["directory_listing"] = fact.DirectoryListing is null ? null
                        : JsonNode.Parse(JsonCodec.WriteEvent(fact))?["tool"]?["directory_listing"]?.DeepClone(),
                    ["parts"] = new JsonArray(fact.Parts.Select(p => (JsonNode)JsonCodec.WritePart(p)).ToArray()),
                    ["is_error"] = fact.IsError,
                    ["job"] = fact.Job is null ? null : JsonNode.Parse(JsonCodec.WriteEvent(fact))?["tool"]?["job"]?.DeepClone()
                });
                ids.Add(id);
            }
            else if (fact.Kind is EventKind.ErrorOccurred or EventKind.TurnInterrupted or EventKind.JobKilled or EventKind.AgentShutdown)
            {
                frames.Add(new JsonObject
                {
                    ["type"] = "lifecycle", ["agent_id"] = agentID, ["seq"] = fact.Seq,
                    ["kind"] = JsonNode.Parse(JsonCodec.WriteEvent(fact))!["type"]!.GetValue<string>(), ["message"] = fact.Error.Length == 0 ? fact.Reason : fact.Error
                });
            }
            else if (fact.Kind == EventKind.TurnEnded)
            {
                frames.Add(new JsonObject
                {
                    ["type"] = "turn_ended", ["agent_id"] = agentID, ["seq"] = fact.Seq,
                    ["reason"] = fact.Reason,
                    ["text"] = fact.Reason == "completed" ? lastAnswer : "",
                    ["error"] = fact.Reason == "completed" ? "" : fact.Reason
                });
                calls.Clear();
                lastAnswer = "";
            }
        }
        return new Projection(frames.ToImmutable(), state, ids.ToImmutable());
    }
}
