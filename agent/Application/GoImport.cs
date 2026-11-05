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
    public static partial int RunGoImportCommand(string[] args, TextWriter output, TextWriter error)
    {
        try
        {
            return Import(args, output, error);
        }
        catch (Exception fault)
        {
            // Operator commands report; they do not hand the operator a stack trace.
            error.WriteLine("import-go: " + fault.Message);
            return 1;
        }
    }

    private static int Import(string[] args, TextWriter output, TextWriter error)
    {
        int position = -1;
        string Next(string flag)
        {
            position++;
            if (position >= args.Length) throw new ArgumentException($"import-go: {flag} needs a value");
            return args[position];
        }

        string source = "";
        string destination = "";
        string primaryName = "root";
        bool replace = false;
        bool dryRun = false;
        List<string> only = [];
        while (position < args.Length - 1)
        {
            position++;
            switch (args[position])
            {
                case "--from": source = Next("--from"); break;
                case "--data-dir": destination = Next("--data-dir"); break;
                case "--agent": only.Add(Next("--agent")); break;
                case "--as": primaryName = Next("--as"); break;
                case "--replace": replace = true; break;
                case "--dry-run": dryRun = true; break;
                default:
                    error.WriteLine($"import-go: unknown flag {args[position]}");
                    return 2;
            }
        }

        string Env(string name) => Environment.GetEnvironmentVariable(name) ?? "";
        string origin;
        if (destination.Length == 0)
        {
            destination = Env("EN_DATA_DIR");
            origin = destination.Length == 0 ? "the current directory" : "EN_DATA_DIR";
        }
        else
        {
            origin = "--data-dir";
        }
        if (destination.Length == 0)
        {
            destination = Directory.GetCurrentDirectory();
        }
        if (source.Length == 0)
        {
            error.WriteLine("import-go: --from <directory> is required; it is the directory the Go version kept its state in.");
            return 2;
        }
        output.WriteLine($"import-go: destination {destination} (from {origin})");
        if (!Directory.Exists(source))
        {
            error.WriteLine($"import-go: no such source directory: {source}");
            return 2;
        }
        if (Path.GetFullPath(source) == Path.GetFullPath(destination))
        {
            error.WriteLine("import-go: the source and destination are the same directory.");
            return 2;
        }

        // Validate names and resolve homes through the same rule the running
        // application uses, so an import can never install a home the Agent could
        // not later open.
        var container = new Ensemble(TextWriter.Null, dataDirectory: destination);

        var planned = new List<Adapted>();
        var problems = new List<string>();

        string primarySource = Path.Combine(source, "save.json");
        if (File.Exists(primarySource))
        {
            Plan(primaryName, primarySource);
        }

        string registryPath = Path.Combine(source, "ensemble.json");
        if (File.Exists(registryPath))
        {
            JsonArray registry;
            try
            {
                using var document = JsonDocument.Parse(File.ReadAllBytes(registryPath));
                if (document.RootElement.ValueKind != JsonValueKind.Array) throw new FormatException();
                registry = JsonNode.Parse(document.RootElement.GetRawText())!.AsArray();
            }
            catch (Exception fault) when (fault is JsonException or FormatException)
            {
                problems.Add($"{registryPath}: the peer registry must be a JSON array of name, model, skill and created.");
                registry = [];
            }
            foreach (var row in registry)
            {
                if (row is not JsonObject entry)
                {
                    problems.Add($"{registryPath}: a peer registry entry must be an object.");
                    continue;
                }
                string? name = entry["name"]?.GetValueKind() == JsonValueKind.String ? entry["name"]!.GetValue<string>() : null;
                if (name is null || name.Length == 0)
                {
                    problems.Add($"{registryPath}: a peer registry entry has no name.");
                    continue;
                }
                if (only.Count > 0 && !only.Contains(name, StringComparer.Ordinal)) continue;
                Plan(name, Path.Combine(source, "agents", name, "save.json"));
            }
        }
        else if (only.Count > 0)
        {
            problems.Add($"{registryPath}: no peer registry to resolve the requested --agent names against.");
        }

        if (only.Count > 0)
        {
            foreach (string wanted in only)
            {
                if (!planned.Any(item => item.Name == wanted) && !problems.Any(text => text.Contains(wanted, StringComparison.Ordinal)))
                {
                    problems.Add($"--agent {wanted}: no such peer in {registryPath}.");
                }
            }
        }

        void Plan(string name, string sourcePath)
        {
            string? home;
            try
            {
                home = container.AgentDirectory(name);
            }
            catch (ArgumentException fault)
            {
                problems.Add($"{name}: {fault.Message}");
                return;
            }
            if (home is null)
            {
                problems.Add($"{name}: the destination cannot host a named Agent.");
                return;
            }
            if (!File.Exists(sourcePath))
            {
                problems.Add($"{name}: no saved history at {sourcePath}.");
                return;
            }
            string destinationPath = Path.Combine(home, "save.json");
            if (File.Exists(destinationPath) && !replace)
            {
                problems.Add($"{name}: {destinationPath} already holds history; pass --replace to overwrite it (a copy is kept).");
                return;
            }
            JsonNode root;
            try
            {
                root = JsonNode.Parse(File.ReadAllText(sourcePath))!;
            }
            catch (Exception fault) when (fault is JsonException or IOException or UnauthorizedAccessException)
            {
                problems.Add($"{name}: {sourcePath} could not be read as JSON: {fault.Message}");
                return;
            }
            try
            {
                planned.Add(Adapt(name, sourcePath, destinationPath, root.AsObject()));
            }
            catch (Exception fault) when (fault is FormatException or InvalidDataException)
            {
                problems.Add($"{name}: {sourcePath} is still unreadable after adaptation: "
                    + (fault.InnerException?.Message is { Length: >0 } inner ? inner + " / " : "") + fault.Message);
            }
        }

        if (problems.Count > 0)
        {
            // Nothing is written until the whole fleet validates, so a half-imported
            // peer set cannot be mistaken for the restored one.
            foreach (string problem in problems) error.WriteLine("import-go: " + problem);
            return 1;
        }
        if (planned.Count == 0)
        {
            error.WriteLine($"import-go: {source} holds no Go history this command recognises.");
            return 1;
        }

        foreach (var item in planned)
        {
            output.WriteLine($"imported \"{item.Name}\": {item.Events} events, {item.Characters} characters of text"
                + (item.Dropped.Count == 0 ? "" : ", dropped " + string.Join(", ", item.Dropped.OrderBy(entry => entry.Key).Select(entry => $"{entry.Value} {entry.Key}")))
                + (item.Translated.Count == 0 ? "" : ", translated " + string.Join(", ", item.Translated.OrderBy(entry => entry.Key).Select(entry => $"{entry.Value} {entry.Key}")))
                + (item.NoText == 0 ? "" : $", {item.NoText} text parts carried no body in the Go files"));
        }
        if (dryRun)
        {
            output.WriteLine($"import-go: dry run; nothing was written. {planned.Count} Agent(s) ready to install in {destination}.");
            return 0;
        }

        foreach (var item in planned)
        {
            Directory.CreateDirectory(Path.GetDirectoryName(item.DestinationPath)!);
            if (replace && File.Exists(item.DestinationPath))
            {
                string kept = $"{item.DestinationPath}.pre-import-{DateTime.UtcNow:yyyyMMdd-HHmmss}";
                File.Copy(item.DestinationPath, kept, overwrite: true);
                output.WriteLine($"imported \"{item.Name}\": previous history kept at {kept}");
            }
            string temporary = item.DestinationPath + ".importing";
            File.WriteAllText(temporary, item.Json);
            File.Move(temporary, item.DestinationPath, overwrite: true);
        }

        MergeRegistry(registryPath);
        output.WriteLine($"import-go: installed {planned.Count} Agent(s). Start Ensemble with --data-dir {destination} to resume them.");
        return 0;

        void MergeRegistry(string sourceRegistry)
        {
            if (!File.Exists(sourceRegistry)) return;
            string destinationRegistry = Path.Combine(destination, "ensemble.json");
            var rows = new JsonArray();
            var names = new HashSet<string>(StringComparer.Ordinal);
            if (File.Exists(destinationRegistry))
            {
                // A damaged registry is never replaced by an import.
                using var existing = JsonDocument.Parse(File.ReadAllBytes(destinationRegistry));
                if (existing.RootElement.ValueKind != JsonValueKind.Array)
                {
                    error.WriteLine($"import-go: {destinationRegistry} is not a JSON array; peers were not merged.");
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
            using var sourceRegistryDocument = JsonDocument.Parse(File.ReadAllBytes(sourceRegistry));
            foreach (var row in sourceRegistryDocument.RootElement.EnumerateArray())
            {
                if (row.ValueKind != JsonValueKind.Object
                    || !row.TryGetProperty("name", out var key) || key.ValueKind != JsonValueKind.String)
                {
                    continue;
                }
                if (names.Add(key.GetString()!)) rows.Add(JsonNode.Parse(row.GetRawText()));
            }
            string temporary = destinationRegistry + ".importing";
            File.WriteAllText(temporary, rows.ToJsonString());
            File.Move(temporary, destinationRegistry, overwrite: true);
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
            string candidate = root.ToJsonString(new JsonSerializerOptions { WriteIndented = true });
            SaveDocument document;
            try
            {
                document = JsonCodec.ReadSave(candidate);
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
            if (dropped.Count > 0)
            {
                // A removed record leaves a hole, and replay requires an unbroken
                // run of sequences. Close it before the document is offered.
                CloseGaps(root, dropped);
                candidate = root.ToJsonString(new JsonSerializerOptions { WriteIndented = true });
                document = JsonCodec.ReadSave(candidate);
            }
            // Write exactly the text that the reader above accepted.
            int characters = document.Events.Sum(event_ => event_.Parts.Sum(part => part is TextPart body ? body.Text.Length : 0));
            int silent = document.Events.Sum(event_ => event_.Parts.Sum(part => part is TextPart body && body.Text.Length == 0 ? 1 : 0));
            return new Adapted(name, destinationPath, candidate, document.Events.Length, characters, silent, dropped, translated);
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
                // Ranks are one-based positions in the renumbered history.
                range["from"] = low + 1;
                range["to"] = high + 1;
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
