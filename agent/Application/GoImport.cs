using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

public sealed partial class Ensemble
{
    // The Go implementation wrote the same nested event envelopes this codec reads,
    // so history travels as data rather than through a per-vendor adapter. Four
    // differences still have to be bridged, each handled where it is understood
    // rather than by loosening the reader for foreign input:
    //
    //   1. Go's compacted snapshot ("context") laid out different fields
    //      (disabled_bands, sent_through, turn, ephemera) than this codec's
    //      (memories, phase, compression). The event log carries the same facts, so
    //      the snapshot is dropped and History rebuilds its own by replay.
    //   2. Go recorded job_finished, job_report and tools_changed. There is no
    //      counterpart here, so they are dropped and counted, never silently.
    //   3. Go omitted "reason" on turn_ended and put the interrupt cause under
    //      "error"; this codec requires reason for that event group.
    //   4. Go's text part omitted "text" whenever it was empty, because its field
    //      is omitempty. An absent key is a legitimate empty reply here, so it is
    //      written as one.
    //
    // An older Go build also spilled part bodies into <home>/io/<artifact id>. That
    // directory is deliberately not read: nothing in the build that wrote these
    // saves consults it, no saved part carries a locator for it, and its ids were
    // reused across restarts, so attaching a file to a part would be a guess that
    // could put another turn's words into an Agent's mouth. Such parts import as
    // empty and are counted in the report.
    // This runs at every launch, silently, and does nothing once the history is
    // in place: the Go version kept its state in whichever directory it was started
    // from, so the conversations Bill had been resuming for weeks were simply not
    // where the C# launcher looked. Discovering them there is the equivalent of
    // starting the old launcher in the old directory, with none of the ceremony.
    //
    // Restoring is deliberately conservative. The source is only ever read; a save
    // this build cannot load is left where it is and reported; an Agent whose home
    // already holds history is never overwritten; and a failure here can never stop
    // the GUI, because losing a conversation is bad but losing the application is
    // worse. A marker records that the restore ran, so deleting one imported Agent
    // does not resurrect it under the same name on the next start.
    private const string Marker = ".go-history-restored";
    private const string ArchiveName = "Go archive";

    public static partial string? RestoreGoHistory(string destination, TextWriter diagnostics)
    {
        try
        {
            return Restore(destination, diagnostics);
        }
        catch (Exception fault)
        {
            // Discovery is a convenience at the edge of composition. Whatever it
            // could not work out belongs in the log, never in the operator's way.
            diagnostics.WriteLine($"go history: not restored ({fault.GetType().Name}: {fault.Message})");
            return null;
        }
    }

    private static string? Restore(string destination, TextWriter diagnostics)
    {
        string wanted = Path.GetFullPath(destination);
        if (File.Exists(Path.Combine(wanted, Marker))) return null;

        string? source = FindGoTree(wanted, diagnostics);
        if (source is null) return null;

        // A container that owns no Agent: it is here only to resolve and validate
        // homes through the same rule the running application uses.
        var container = new Ensemble(TextWriter.Null, dataDirectory: destination);
        var installed = new List<(string Name, int Events, int Characters, int NoText,
            Dictionary<string, int> Dropped, Dictionary<string, int> Translated)>();
        var notes = new List<string>();

        // The primary Agent kept its save at the data directory itself; the Go GUI
        // opened it as the Agent it was serving. It becomes this build's primary
        // home when that home holds nothing worth keeping, and a named peer when it
        // holds a conversation, so a restore can never bury work already done here.
        string primarySource = Path.Combine(source, "save.json");
        string? primaryTarget = File.Exists(primarySource) ? PrimaryTarget(container, destination) : null;
        if (primaryTarget is not null)
        {
            // The home's directory names the Agent; the file inside it is always
            // save.json, so the name has to come from the directory above it.
            bool adopted = string.Equals(Path.GetFileName(Path.GetDirectoryName(primaryTarget)),
                "root", StringComparison.Ordinal);
            Plan(container, source, primarySource, adopted ? "root" : ArchiveName,
                installed, notes, diagnostics, adoptQuietHome: adopted);
        }

        string registryPath = Path.Combine(source, "ensemble.json");
        if (File.Exists(registryPath))
        {
            try
            {
                using var document = JsonDocument.Parse(File.ReadAllBytes(registryPath));
                if (document.RootElement.ValueKind == JsonValueKind.Array)
                {
                    foreach (var row in document.RootElement.EnumerateArray())
                    {
                        if (row.ValueKind != JsonValueKind.Object
                            || !row.TryGetProperty("name", out var key) || key.ValueKind != JsonValueKind.String)
                        {
                            continue;
                        }
                        string name = key.GetString()!;
                        if (name.Length == 0) continue;
                        Plan(container, source, Path.Combine(source, "agents", name, "save.json"), name,
                            installed, notes, diagnostics);
                    }
                }
                else
                {
                    notes.Add("the peer registry is not a JSON array; only the primary conversation was restored");
                }
            }
            catch (Exception fault) when (fault is JsonException or IOException or UnauthorizedAccessException)
            {
                notes.Add($"the peer registry could not be read ({fault.GetType().Name})");
            }
        }

        if (installed.Count == 0)
        {
            // Nothing to install is a normal steady state, not a failure: report it
            // once in the log and leave the terminal alone.
            File.WriteAllText(Path.Combine(wanted, Marker),
                $"restored at {DateTime.UtcNow:O}; source {source}; nothing new to install\n");
            foreach (string note in notes) diagnostics.WriteLine("go history: " + note);
            return null;
        }

        foreach (var item in installed)
        {
            diagnostics.WriteLine($"go history: restored \"{item.Name}\" from {source} - {item.Events} events, "
                + $"{item.Characters} characters of text"
                + (item.NoText == 0 ? ""
                    : $", {item.NoText} {(item.NoText == 1 ? "part" : "parts")} had no body in the Go files")
                + (item.Translated.Count == 0 ? ""
                    : ", translated " + string.Join(", ", item.Translated.OrderBy(entry => entry.Key)
                        .Select(entry => $"{entry.Value} {entry.Key}")))
                + (item.Dropped.Count == 0 ? ""
                    : ", skipped " + string.Join(", ", item.Dropped.OrderBy(entry => entry.Key)
                        .Select(entry => $"{entry.Value} {entry.Key}"))));
        }
        MergeRegistry(source, destination, diagnostics, notes);
        foreach (string note in notes) diagnostics.WriteLine("go history: " + note);
        File.WriteAllText(Path.Combine(wanted, Marker),
            $"restored at {DateTime.UtcNow:O}; source {source}; {installed.Count} Agent(s): "
            + string.Join(", ", installed.Select(item => item.Name)) + "\n");
        int events = installed.Sum(item => item.Events);
        return installed.Count == 1
            ? $"restored 1 Agent from the Go history in {source} ({events} events)"
            : $"restored {installed.Count} Agents from the Go history in {source} ({events} events)";
    }

    // Where the Go primary conversation should land. A home with no conversation of
    // its own adopts it; one with a conversation becomes a peer instead.
    private static string? PrimaryTarget(Ensemble container, string destination)
    {
        string? home = container.AgentDirectory("root");
        if (home is null) return null;
        string path = Path.Combine(home, "save.json");
        if (!File.Exists(path)) return path;
        try
        {
            var existing = JsonCodec.ReadSave(File.ReadAllText(path));
            bool talks = existing.Context.Dialogue.Any(entry => entry.Parts.Any(part =>
                (part is TextPart body && body.Text.Length > 0) || part is ToolCallPart));
            // Nothing said here yet: the Go primary is the conversation this operator
            // was resuming, so it becomes the primary rather than a stranger in the
            // tree. A home that has said something keeps it, and arrives as a peer.
            return talks ? Archive(container) : path;
        }
        catch (Exception fault) when (fault is FormatException or InvalidDataException
            or IOException or UnauthorizedAccessException or JsonException)
        {
            // An unreadable existing home is not a reason to discard it, nor to
            // overwrite it: the Go primary arrives as a peer and the log says so.
            return Archive(container);
        }
    }

    private static string? Archive(Ensemble container)
    {
        string? home = container.AgentDirectory(ArchiveName);
        return home is null ? null : Path.Combine(home, "save.json");
    }

    // Validate before installing, and never install beside a history that already
    // exists: the two would otherwise be indistinguishable in the tree.
    private static void Plan(Ensemble container, string source, string sourcePath, string name,
        List<(string Name, int Events, int Characters, int NoText,
            Dictionary<string, int> Dropped, Dictionary<string, int> Translated)> installed,
        List<string> notes, TextWriter diagnostics, bool adoptQuietHome = false)
    {
        string? home = container.AgentDirectory(name);
        if (home is null)
        {
            notes.Add($"{name}: this data directory cannot host a named Agent");
            return;
        }
        string destinationPath = Path.Combine(home, "save.json");
        bool exists = File.Exists(destinationPath);
        // An Agent that already has history here is left alone, unless this is the
        // Go primary arriving in a home that has never said anything.
        if (exists && !adoptQuietHome) return;
        if (!File.Exists(sourcePath))
        {
            notes.Add($"{name}: no saved history at {sourcePath}");
            return;
        }
        JsonNode? root;
        try
        {
            root = JsonNode.Parse(File.ReadAllText(sourcePath));
        }
        catch (Exception fault) when (fault is JsonException or IOException or UnauthorizedAccessException)
        {
            notes.Add($"{name}: {sourcePath} could not be read as JSON");
            return;
        }
        if (root is null)
        {
            notes.Add($"{name}: {sourcePath} is empty");
            return;
        }
        try
        {
            Adapted adapted = Adapt(name, sourcePath, destinationPath, root.AsObject());
            Directory.CreateDirectory(Path.GetDirectoryName(adapted.DestinationPath)!);
            if (exists)
            {
                // Adoption is still a replacement of bytes that were here first. The
                // quiet save is kept beside the new one rather than discarded, so a
                // restore can be undone by whoever notices it should not have run.
                string kept = adapted.DestinationPath + ".before-go-restore";
                File.Copy(adapted.DestinationPath, kept, overwrite: true);
                diagnostics.WriteLine($"go history: the quiet home at {adapted.DestinationPath} was kept at {kept}");
            }
            string temporary = adapted.DestinationPath + ".restoring";
            File.WriteAllText(temporary, adapted.Json);
            File.Move(temporary, adapted.DestinationPath, overwrite: exists);
            installed.Add((name, adapted.Events, adapted.Characters, adapted.NoText, adapted.Dropped,
                adapted.Translated));
        }
        catch (Exception fault) when (fault is FormatException or InvalidDataException
            or IOException or UnauthorizedAccessException)
        {
            notes.Add($"{name}: its Go history is not readable by this build ({fault.Message})");
        }
    }

    // The Go registry is the list of peers; merging keeps whatever this data
    // directory already had, and a damaged registry is left alone rather than
    // replaced by an import that would look like it had always been there.
    private static void MergeRegistry(string source, string destination, TextWriter diagnostics, List<string> notes)
    {
        string sourceRegistry = Path.Combine(source, "ensemble.json");
        if (!File.Exists(sourceRegistry)) return;
        string destinationRegistry = Path.Combine(destination, "ensemble.json");
        var rows = new JsonArray();
        var names = new HashSet<string>(StringComparer.Ordinal);
        if (File.Exists(destinationRegistry))
        {
            try
            {
                using var existing = JsonDocument.Parse(File.ReadAllBytes(destinationRegistry));
                if (existing.RootElement.ValueKind != JsonValueKind.Array)
                {
                    notes.Add($"{destinationRegistry} is not a JSON array; peers were not merged");
                    return;
                }
                foreach (var row in existing.RootElement.EnumerateArray())
                {
                    rows.Add(JsonNode.Parse(row.GetRawText()));
                    if (row.ValueKind == JsonValueKind.Object && row.TryGetProperty("name", out var key)
                        && key.ValueKind == JsonValueKind.String)
                    {
                        names.Add(key.GetString()!);
                    }
                }
            }
            catch (Exception fault) when (fault is JsonException or IOException or UnauthorizedAccessException)
            {
                notes.Add($"{destinationRegistry} could not be read; peers were not merged");
                return;
            }
        }
        try
        {
            using var incoming = JsonDocument.Parse(File.ReadAllBytes(sourceRegistry));
            if (incoming.RootElement.ValueKind != JsonValueKind.Array) return;
            foreach (var row in incoming.RootElement.EnumerateArray())
            {
                if (row.ValueKind != JsonValueKind.Object
                    || !row.TryGetProperty("name", out var key) || key.ValueKind != JsonValueKind.String)
                {
                    continue;
                }
                if (names.Add(key.GetString()!)) rows.Add(JsonNode.Parse(row.GetRawText()));
            }
        }
        catch (Exception fault) when (fault is JsonException or IOException or UnauthorizedAccessException)
        {
            notes.Add("the Go peer registry could not be read; existing peers were kept");
            return;
        }
        string temporary = destinationRegistry + ".merging";
        File.WriteAllText(temporary, rows.ToJsonString());
        File.Move(temporary, destinationRegistry, overwrite: true);
    }

    // The places the Go version was ever started from. The working directory first,
    // because that is where it kept its state when it was the daily driver; then the
    // XDG location a later Go build used. Nothing here searches the filesystem.
    private static string? FindGoTree(string wantedDestination, TextWriter diagnostics)
    {
        var candidates = new List<string>();
        string configured = Environment.GetEnvironmentVariable("EN_GO_DATA_DIR") ?? "";
        if (configured.Length > 0) candidates.Add(configured);
        candidates.Add(Directory.GetCurrentDirectory());
        string home = Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
        if (home.Length > 0) candidates.Add(Path.Combine(home, ".config", "ensemble"));

        foreach (string candidate in candidates.Distinct(StringComparer.Ordinal))
        {
            if (string.Equals(Path.GetFullPath(candidate), wantedDestination, StringComparison.Ordinal)) continue;
            if (LooksLikeGoState(candidate)) return candidate;
            diagnostics.WriteLine($"go history: {candidate} holds no Go history");
        }
        return null;
    }

    // Go's save documents carry snapshot fields this codec never writes, which makes
    // them recognisable without probing a directory that merely happens to contain a
    // save.json. A registry of peers counts on its own.
    private static bool LooksLikeGoState(string directory)
    {
        try
        {
            string registry = Path.Combine(directory, "ensemble.json");
            if (File.Exists(registry))
            {
                using var document = JsonDocument.Parse(File.ReadAllBytes(registry));
                if (document.RootElement.ValueKind == JsonValueKind.Array
                    && document.RootElement.GetArrayLength() > 0) return true;
            }
            string save = Path.Combine(directory, "save.json");
            if (!File.Exists(save)) return false;
            using var head = JsonDocument.Parse(File.ReadAllText(save));
            if (head.RootElement.ValueKind != JsonValueKind.Object
                || !head.RootElement.TryGetProperty("log", out var log)
                || log.ValueKind != JsonValueKind.Array)
            {
                return false;
            }
            if (!head.RootElement.TryGetProperty("context", out var context)
                || context.ValueKind != JsonValueKind.Object)
            {
                return false;
            }
            foreach (string key in new[] { "disabled_bands", "sent_through", "startup_tools", "ephemera", "turn" })
            {
                if (context.TryGetProperty(key, out _)) return true;
            }
            return false;
        }
        catch (Exception fault) when (fault is JsonException or IOException or UnauthorizedAccessException)
        {
            return false;
        }
    }

    // Go event type -> the event this build uses for the same fact. Only
    // job_report is here. job_finished looks like a tool result and was mapped to
    // one at first, but every one of Bill's 668 job_finished events carries an
    // empty call_id: Go recorded a job's lifecycle, not the outcome of a model's
    // call, and this codec requires a real call ID for a tool result. Inventing an
    // identity to satisfy the reader would attach job output to a call that never
    // asked for it, so job_finished stays what the reader says it is, and the
    // adaptation loop reports it as dropped.
    private static readonly Dictionary<string, string> Renamed = new(StringComparer.Ordinal)
    {
        ["job_report"] = "tool_observed"
    };

    // One imported Agent, validated by the real reader before it is offered for install.
    private sealed record Adapted(string Name, string DestinationPath, string Json,
        int Events, int Characters, int NoText, Dictionary<string, int> Dropped, Dictionary<string, int> Translated);

    // Rewrite one Go save document into what this codec reads, then prove it reads.
    // Rather than freeze a list of Go's vocabulary, the loop asks the real reader
    // what it cannot parse and retires exactly that, by name, until the document
    // loads. Anything unknown is therefore reported instead of silently skipped, and
    // a malformed event still fails the import instead of being dropped with it.
    private static Adapted Adapt(string name, string sourcePath, string destinationPath, JsonObject root)
    {
        // Rebuilt by replay; this build has no reader for Go's snapshot layout.
        root["context"] = null;
        var dropped = new Dictionary<string, int>(StringComparer.Ordinal);
        var translated = new Dictionary<string, int>(StringComparer.Ordinal);
        // Normalise before asking the reader: absent text keys and missing turn
        // reasons are Go's omitempty, not damage this import should report as such.
        FillEmptyText(root, translated);
        for (int attempt = 0; ; attempt++)
        {
            if (attempt > 20) throw new InvalidDataException("too many unreadable event types to adapt");
            if (dropped.Count > 0)
            {
                // Retiring a record leaves a hole, and this reader demands an unbroken
                // run. Closing has to happen before the next read, or the recovery
                // trips over the very gap it just made. Renumbering an already
                // unbroken log leaves it, and every redaction range, untouched.
                CloseGaps(root, dropped);
            }
            string candidate = root.ToJsonString(new JsonSerializerOptions { WriteIndented = true });
            try
            {
                _ = JsonCodec.ReadSave(candidate);
            }
            catch (FormatException fault) when (fault.Message.StartsWith("Unknown EventKind:", StringComparison.Ordinal))
            {
                string kind = fault.Message["Unknown EventKind:".Length..].Trim();
                int removed = 0;
                if (root["log"] is JsonArray log)
                {
                    foreach (var entry in log.ToArray())
                    {
                        if (entry is JsonObject item && item["type"]?.GetValueKind() == JsonValueKind.String
                            && item["type"]!.GetValue<string>() == kind)
                        {
                            removed++;
                            log.Remove(item);
                        }
                    }
                }
                if (removed == 0) throw;
                dropped[kind] = dropped.GetValueOrDefault(kind) + removed;
                continue;
            }
            // Write exactly the text that the reader above accepted.
            var written = CountWritten(root);
            return new Adapted(name, destinationPath, candidate, written.Events, written.Characters,
                written.Silent, dropped, translated);
        }
    }

    // The reader folds the log into Context and retains no event collection, so the
    // import reports the dialogue it actually wrote by reading its own document back.
    // Each event kind carries its parts under its own envelope -- message, response,
    // tool -- so the walk is structural rather than tied to one layout. Text is the
    // only part an operator reads as speaking, and a text part that arrived with no
    // body is Go's omitempty, which the operator is told to expect.
    private static (int Events, int Characters, int Silent) CountWritten(JsonObject root)
    {
        int events = 0;
        int characters = 0;
        int silent = 0;
        if (root["log"] is JsonArray log)
        {
            foreach (var entry in log)
            {
                if (entry is not JsonObject item)
                {
                    continue;
                }
                events++;
                CountText(item, ref characters, ref silent);
            }
        }
        return (events, characters, silent);
    }

    private static void CountText(JsonNode node, ref int characters, ref int silent)
    {
        if (node is JsonArray entries)
        {
            foreach (var entry in entries)
            {
                if (entry is not null) CountText(entry, ref characters, ref silent);
            }
            return;
        }
        if (node is not JsonObject item)
        {
            return;
        }
        if (item["type"]?.GetValueKind() == JsonValueKind.String
            && item["type"]!.GetValue<string>() == "text")
        {
            string body = item["text"]?.GetValueKind() == JsonValueKind.String
                ? item["text"]!.GetValue<string>() : "";
            characters += body.Length;
            if (body.Length == 0)
            {
                silent++;
            }
        }
        foreach (var child in item)
        {
            if (child.Value is not null) CountText(child.Value, ref characters, ref silent);
        }
    }

    // Go's omitempty dropped the key for legitimate empty text; an absent key and an
    // empty reply are the same fact here. Bodies that an older Go build spilled into
    // <home>/io/ are not consulted, because its ids were reused across restarts and
    // nothing in the build that wrote these saves reads that directory: guessing
    // would put another turn's words into an Agent's mouth.
    private static void FillEmptyText(JsonObject root, Dictionary<string, int> translated)
    {
        if (root["log"] is not JsonArray log) return;
        foreach (var entry in log)
        {
            if (entry is not JsonObject item) continue;
            string type = item["type"]?.GetValueKind() == JsonValueKind.String ? item["type"]!.GetValue<string>() : "";
            if ((type is "turn_ended" or "turn_interrupted") && item["reason"]?.GetValueKind() != JsonValueKind.String)
            {
                JsonNode? cause = item["error"]?["message"];
                item["reason"] = cause?.GetValueKind() == JsonValueKind.String ? cause.GetValue<string>() : "";
            }
            // Go had three lifecycle events this codec names differently, because the
            // two builds describe the same facts from different owners. The mapping
            // follows this codec's own rules: a peer call's result is an audit outcome
            // with no assistant call behind it (EventKind.ToolObserved), exactly as
            // agent/Agent/Actor.cs chooses that kind for a host-invoked tool.
            // Both are renamed only when the payload already fits what this reader
            // requires; anything else is left for the reader to refuse or drop, so a
            // partial match never becomes an invented history event.
            if (Renamed.TryGetValue(type, out string? outcome)
                && item["tool"] is JsonObject report
                && report["call_id"]?.GetValueKind() == JsonValueKind.String
                && report["call_id"]!.GetValue<string>().Length > 0)
            {
                item["type"] = outcome;
                translated[$"{type} to {outcome}"] = translated.GetValueOrDefault($"{type} to {outcome}") + 1;
            }

            // Go nested a recall under its type name and stored the attached body as
            // an ordinary text part plus a parallel hash list. This codec carries the
            // same facts as one recall part holding both, so the two are joined only
            // when the Go shape is exactly that: one text part and a list of hash
            // strings. Anything else is left untouched, where the reader reports it as
            // an unknown kind instead of this import inventing a shape for it.
            if (type == "recall_attached" && item["recall_attached"] is JsonObject recall
                && recall["parts"] is JsonArray { Count: 1 } recalled
                && recalled[0] is JsonObject attached
                && attached["type"]?.GetValueKind() == JsonValueKind.String
                && attached["type"]!.GetValue<string>() == "text"
                && attached["text"]?.GetValueKind() == JsonValueKind.String
                && recall["hashes"] is JsonArray hashes
                && hashes.All(hash => hash?.GetValueKind() == JsonValueKind.String))
            {
                item["parts"] = new JsonArray(new JsonObject
                {
                    ["type"] = "recall",
                    ["text"] = attached["text"]!.DeepClone(),
                    ["hashes"] = hashes.DeepClone()
                });
                item.Remove("recall_attached");
                item["type"] = "recall_added";
                translated["recall_attached to recall_added"] = translated.GetValueOrDefault("recall_attached to recall_added") + 1;
            }
            foreach (string holder in new[] { "response", "tool", "message" })
            {
                if (item[holder]?["parts"] is JsonArray parts) FillPartsEmpty(parts);
            }
            if (item["parts"] is JsonArray direct) FillPartsEmpty(direct);
        }
    }

    // A dropped Go event leaves a gap, and History replays a save only when its
    // sequences run unbroken from one (agent/History/Persistence.cs:18). Survivors
    // are therefore renumbered in place, and the one other fact in a save that
    // names a sequence is a redaction range: translating those too is what keeps a
    // migrated redaction covering the records the Go history said it covered rather
    // than silently redacting different ones. A redaction whose range no longer
    // reaches anything surviving is retired with the records it described.
    private static void CloseGaps(JsonObject root, Dictionary<string, int> dropped)
    {
        if (root["log"] is not JsonArray log) return;
        var survivors = new List<long>();
        foreach (var entry in log)
        {
            if (entry is JsonObject item && item["seq"]?.GetValueKind() == JsonValueKind.Number)
            {
                survivors.Add(item["seq"]!.GetValue<long>());
            }
        }

        // The surviving records that a range actually covered are those whose
        // original sequences lie inside it. Snapping each end to the nearest
        // neighbour on the wrong side would widen the range onto records the Go
        // history never redacted, so both bounds are taken from within.
        static int FirstAtOrAfter(List<long> survivors, long sequence)
        {
            int low = 0;
            int high = survivors.Count;
            while (low < high)
            {
                int middle = (low + high) / 2;
                if (survivors[middle] < sequence) low = middle + 1;
                else high = middle;
            }
            return low;
        }

        static int LastBefore(List<long> survivors, long sequence)
        {
            int low = 0;
            int high = survivors.Count;
            while (low < high)
            {
                int middle = (low + high) / 2;
                if (survivors[middle] < sequence) low = middle + 1;
                else high = middle;
            }
            return low - 1;
        }

        var retiring = new List<JsonNode?>();
        foreach (var entry in log)
        {
            if (entry is not JsonObject item) continue;
            if (item["redact"] is JsonObject range
                && range["from"]?.GetValueKind() == JsonValueKind.Number
                && range["to"]?.GetValueKind() == JsonValueKind.Number)
            {
                int low = FirstAtOrAfter(survivors, range["from"]!.GetValue<long>());
                int high = LastBefore(survivors, range["to"]!.GetValue<long>() + 1);
                if (low > high || low >= survivors.Count)
                {
                    retiring.Add(entry);
                    continue;
                }
                // Ranks are one-based positions in the renumbered history. They are
                // written as the same width the Go sequences were read with, because a
                // second pass over an already-closed log reads these values back.
                range["from"] = (long)(low + 1);
                range["to"] = (long)(high + 1);
            }
        }
        foreach (var entry in retiring)
        {
            log.Remove(entry);
            string key = "redacted (covered nothing that survived)";
            dropped[key] = dropped.GetValueOrDefault(key) + 1;
        }

        long assigned = 0;
        foreach (var entry in log)
        {
            if (entry is JsonObject) entry["seq"] = ++assigned;
        }
        root["as_of"] = assigned;
    }

    private static void FillPartsEmpty(JsonArray parts)
    {
        foreach (var part in parts)
        {
            if (part is not JsonObject item) continue;
            if (item["type"]?.GetValueKind() == JsonValueKind.String && item["type"]!.GetValue<string>() == "text"
                && item["text"]?.GetValueKind() != JsonValueKind.String)
            {
                item["text"] = "";
            }
            if (item["parts"] is JsonArray nested) FillPartsEmpty(nested);
        }
    }
}
