using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;
using Application = Ensemble.Ensemble;

// The Go implementation kept the same nested event envelopes but a different
// snapshot layout, a wider vocabulary, and omitempty text parts. These fixtures are
// written in Go's shape on purpose: the point is that foreign bytes this build never
// wrote become a history it can open.
internal static class GoImportTests
{
    private static void Check(bool value, string message)
    {
        if (!value) throw new Exception(message);
    }

    private static string GoSave(long asOf, params string[] events)
    {
        return new JsonObject
        {
            ["config"] = new JsonObject
            {
                ["model"] = "qwen-test", ["vendor"] = "openai", ["system_prompt"] = "go prompt",
                ["tools"] = new JsonArray(new JsonObject
                {
                    ["name"] = "list_directory", ["description"] = "List a directory",
                    ["schema"] = new JsonObject { ["type"] = "object" }
                })
            },
            // Go's snapshot fields; none of them exist in this codec's reader.
            ["as_of"] = asOf,
            ["context"] = new JsonObject
            {
                ["dialogue"] = new JsonArray(), ["disabled_bands"] = new JsonArray(false, false, false, false, false, false),
                ["ephemera"] = new JsonArray(), ["sent_through"] = 0, ["startup_tools"] = new JsonArray(),
                ["tools"] = new JsonArray(), ["turn"] = 0
            },
            ["log"] = JsonNode.Parse("[" + string.Join(",", events) + "]")
        }.ToJsonString(new JsonSerializerOptions { WriteIndented = true });
    }

    // Written the way Go actually wrote them: the model's call is a tool_call part
    // inside a response, each response follows its own request, and a result rides on
    // a tool_returned with plain text parts. History pairs a result with the call that
    // precedes it, so an invented shape here would test the fixture, not the import.
    // The two droppable events sit before the redaction on purpose: its range has to
    // move, and the conversation the test reads back has to stay visible.
    private static readonly string[] Primary =
    [
        """{"seq":1,"type":"message_received","time":"2026-10-25T04:00:00Z","message":{"actor":"human","parts":[{"type":"text","text":"design a slip-in PDF cover"}]}}""",
        """{"seq":2,"type":"request_sent","time":"2026-10-25T04:00:01Z","request":{"to":{"vendor":"openai","model":"qwen-test","surface":"messages"},"startup_tools":[]}}""",
        """{"seq":3,"type":"response_ended","time":"2026-10-25T04:00:02Z","response":{"parts":[{"artifact_id":1,"type":"opaque","from":{"vendor":"openai","model":"qwen-test","surface":"messages"},"data":{"type":"thinking","thinking":"plan the cover"}}]}}""",
        """{"seq":4,"type":"request_sent","time":"2026-10-25T04:00:03Z","request":{"to":{"vendor":"openai","model":"qwen-test","surface":"messages"}}}""",
        // Go's omitempty wrote no key at all for a legitimate empty reply body.
        """{"seq":5,"type":"response_ended","time":"2026-10-25T04:00:04Z","response":{"parts":[{"artifact_id":2,"type":"text"}]}}""",
        """{"seq":6,"type":"request_sent","time":"2026-10-25T04:00:05Z","request":{"to":{"vendor":"openai","model":"qwen-test","surface":"messages"}}}""",
        """{"seq":7,"type":"response_ended","time":"2026-10-25T04:00:06Z","response":{"parts":[{"artifact_id":3,"type":"tool_call","from":{"vendor":"openai","model":"qwen-test","surface":"messages"},"call_id":"call-1","name":"list_directory","args":{"path":"."}}]}}""",
        """{"seq":8,"type":"tool_returned","time":"2026-10-25T04:00:07Z","tool":{"name":"list_directory","call_id":"call-1","parts":[{"type":"text","text":"manuals/ blender.pdf"}],"is_error":false}}""",
        """{"seq":9,"type":"turn_ended","time":"2026-10-25T04:00:08Z"}""",
        // Dropped: Go recorded a job's lifecycle with no call behind it, and a tool
        // set change this build reads from the saved configuration instead.
        """{"seq":10,"type":"job_finished","time":"2026-10-25T04:00:09Z","tool":{"call_id":"","job":{"handle":268,"status":"done","output":{"kind":3,"locator":"io/268"},"bytes":244}}}""",
        """{"seq":11,"type":"tools_changed","time":"2026-10-25T04:00:10Z","tools":{"added":["run_command"],"removed":[]}}""",
        // A second exchange that the migration's redaction below covers, so its
        // range has to move by the two retired records without losing its meaning.
        """{"seq":12,"type":"message_received","time":"2026-10-25T04:00:11Z","message":{"actor":"human","parts":[{"type":"text","text":"what else is in that folder?"}]}}""",
        """{"seq":13,"type":"request_sent","time":"2026-10-25T04:00:12Z","request":{"to":{"vendor":"openai","model":"qwen-test","surface":"messages"}}}""",
        """{"seq":14,"type":"response_ended","time":"2026-10-25T04:00:13Z","response":{"parts":[{"artifact_id":4,"type":"text","text":"the folder holds manuals"}]}}""",
        """{"seq":15,"type":"redacted","time":"2026-10-25T04:00:14Z","redact":{"from":12,"to":14,"level":"redact_dialogue","reason":"compacted by go","replacement":[{"type":"text","text":"a second exchange, summarised"}]}}""",
        """{"seq":16,"type":"recall_attached","time":"2026-10-25T04:00:15Z","recall_attached":{"parts":[{"type":"text","text":"[Auto-recalled memories] blender notes"}],"hashes":["3b38c1c92b6e1add6ee7cdd0bf5d7e1a152b60be27b9bf03e4704a3851b8bf99"]}}""",
        // A peer call's result: Go's job_report, which this build records as an
        // observed tool outcome because no assistant call stands behind it.
        """{"seq":17,"type":"job_report","time":"2026-10-25T04:00:16Z","tool":{"name":"peer_tool","call_id":"call-2","parts":[{"type":"text","text":"peer answered"}],"job":{"handle":9,"status":"done","output":{"kind":3,"locator":"io/9"},"bytes":12}}}"""
    ];

    // The launcher calls this with no operator involvement, so the tests drive it the
    // same way: point the working directory at a Go tree and ask for a restore.
    private static string? Restore(string destination)
    {
        var diagnostics = new StringWriter();
        string? summary = Application.RestoreGoHistory(destination, diagnostics);
        RestoreLog = diagnostics.ToString();
        return summary;
    }

    public static string RestoreLog { get; private set; } = "";

    private static void GoEnvironment(string? source, Action body)
    {
        string original = Directory.GetCurrentDirectory();
        string? previous = Environment.GetEnvironmentVariable("EN_GO_DATA_DIR");
        try
        {
            Environment.SetEnvironmentVariable("EN_GO_DATA_DIR", source);
            body();
        }
        finally
        {
            Environment.SetEnvironmentVariable("EN_GO_DATA_DIR", previous);
            Directory.SetCurrentDirectory(original);
        }
    }

    private static void WriteGoTree(string source, bool brokenPeer = false)
    {
        string peerHome = Path.Combine(source, "agents", "PDF designer");
        Directory.CreateDirectory(peerHome);
        File.WriteAllText(Path.Combine(source, "save.json"), GoSave(17, Primary));
        File.WriteAllText(Path.Combine(source, "ensemble.json"),
            """[{"name":"PDF designer","model":"qwen-test","skill":"","created":"2026-10-25T04:55:26.331352951Z"}]""");
        File.WriteAllText(Path.Combine(peerHome, "save.json"), brokenPeer
            ? "{ \"context\": {\"turn\": 1}, \"as_of\": 1, \"config\": {}, \"log\": [{\"seq\":1,\"type\":\"request_sent\",\"time\":\"" + "2026-10-25T05:00:00Z\"}]}"
            : GoSave(1, """{"seq":1,"type":"message_received","time":"2026-10-25T05:00:00Z","message":{"actor":"human","parts":[{"type":"text","text":"peer question"}]}}"""));
    }

    public static void Run(string temp)
    {
        string source = Path.Combine(temp, "go-auto");
        string destination = Path.Combine(temp, "live-auto");
        Directory.CreateDirectory(destination);
        WriteGoTree(source);
        byte[] before = File.ReadAllBytes(Path.Combine(source, "save.json"));

        GoEnvironment(null, () =>
        {
            // No environment hint: discovery falls back to the working directory,
            // which is exactly how the Go version found its own state.
            Directory.SetCurrentDirectory(source);
            string? summary = Restore(destination);
            Check(summary is not null && summary.Contains("2 Agents") && summary.Contains(source),
                "starting beside a Go tree did not restore it: " + summary);
            Check(RestoreLog.Contains("skipped 1 job_finished") && RestoreLog.Contains("tools_changed"),
                "dropped event types were not named in the log: " + RestoreLog);
            Check(RestoreLog.Contains("recall_attached to recall_added")
                && RestoreLog.Contains("job_report to tool_observed"),
                "a translation was not logged: " + RestoreLog);
            Check(RestoreLog.Contains("had no body"), "an absent part body was not reported: " + RestoreLog);
        });

        Check(before.AsSpan().SequenceEqual(File.ReadAllBytes(Path.Combine(source, "save.json"))),
            "discovery modified the Go source files");

        var imported = JsonCodec.ReadSave(File.ReadAllText(Path.Combine(destination, "agents", "root", "save.json")));
        Check(imported.Events.Length == 15, $"expected the fifteen readable events, got {imported.Events.Length}");
        Check(imported.Events.Any(e => e.Kind == EventKind.MessageReceived
            && e.Parts.Length == 1 && e.Parts[0] is TextPart ask && ask.Text.Contains("slip-in PDF cover")),
            "the user's own words did not survive");
        Check(imported.Events.Any(e => e.Kind == EventKind.ToolReturned
            && e.Parts.Any(part => part is TextPart result && result.Text.Contains("blender.pdf"))),
            "tool output did not survive");
        var recall = imported.Events.Single(e => e.Kind == EventKind.RecallAdded);
        Check(recall.Parts.Length == 1 && recall.Parts[0] is RecallPart attached
            && attached.Attachment.Text.Contains("blender notes"),
            "the recalled memories lost their text");
        var redaction = imported.Events.Single(e => e.Kind == EventKind.Redacted);
        Check(redaction.RedactFrom == 10 && redaction.RedactTo == 12,
            $"the migrated redaction range drifted: {redaction.RedactFrom}..{redaction.RedactTo}");
        Check(imported.Events[^1].Seq == imported.Events.Length && imported.Sequence == imported.Events.Length,
            "the migrated history is not an unbroken run from one");
        Check(imported.Context is null, "Go's snapshot was carried over instead of rebuilt");

        var registry = JsonNode.Parse(File.ReadAllText(Path.Combine(destination, "ensemble.json")))!.AsArray();
        Check(registry.Count(row => row?["name"]?.GetValue<string>() == "PDF designer") == 1,
            "the restored peer never reached the registry");

        string original = Directory.GetCurrentDirectory();
        try
        {
            Directory.SetCurrentDirectory(destination);
            var app = new Application(TextWriter.Null, dataDirectory: destination);
            try
            {
                var agent = new Agent(app, new Configuration(Vendor.OpenAI, "qwen-test", "http://127.0.0.1:1"), name: "root");
                try
                {
                    bool restored;
                    try
                    {
                        restored = agent.StartPersistence(new PersistenceOptions(Path.Combine("agents", "root", "save.json")));
                    }
                    catch (Exception fault)
                    {
                        throw new Exception($"restore failed: {fault.Message}", fault);
                    }
                    Check(restored, "an imported Agent could not restore its history");
                    string wire = agent.Engine().Render(agent.History().View()).ToJsonString();
                    Check(wire.Contains("design a slip-in PDF cover") && wire.Contains("blender.pdf"),
                        "the resumed conversation is missing its substance");
                    Check(wire.Contains("call-1"), "the imported call identity was lost");
                    Check(!wire.Contains("the folder holds manuals"),
                        "a record the Go history redacted with dialogue scope came back into view");
                }
                finally { agent.Destroy(); }
            }
            finally { app.Destroy(); }
        }
        finally { Directory.SetCurrentDirectory(original); }

        // An Agent that already has history here keeps it, byte for byte: startup
        // convenience must never be able to replace work with a copy of something
        // else. The source is a healthy Go peer, so surviving this is a refusal
        // rather than an accident of the source being unreadable.
        string keptPath = Path.Combine(destination, "agents", "PDF designer", "save.json");
        string keptText = File.ReadAllText(keptPath).Replace("peer question", "an answer written here later");
        Check(keptText.Contains("an answer written here later"), "the fixture did not diverge from its source");
        File.WriteAllText(keptPath, keptText);
        // Remove the marker so this really is a re-scan: the refusal to overwrite has
        // to come from the install path, not from discovery having given up already.
        File.Delete(Path.Combine(destination, ".go-history-restored"));
        GoEnvironment(source, () =>
        {
            Restore(destination);
            Check(File.ReadAllText(keptPath) == keptText,
                "discovery replaced an Agent's existing history");
        });

        // Steady state: the second launch finds nothing to do and says nothing.
        GoEnvironment(source, () =>
        {
            string? again = Restore(destination);
            Check(again is null, "a second launch re-restored history it had already installed: " + again);
        });
    }

    public static void Adoption(string temp)
    {
        string source = Path.Combine(temp, "go-adopt");
        WriteGoTree(source);

        // A data directory that has never held a conversation adopts the Go primary
        // as its own, which is what "as it was before" means for the main Agent.
        string quiet = Path.Combine(temp, "live-quiet");
        Directory.CreateDirectory(Path.Combine(quiet, "agents", "root"));
        File.WriteAllText(Path.Combine(quiet, "agents", "root", "save.json"),
            new JsonObject
            {
                ["context"] = null, ["as_of"] = 1,
                ["config"] = JsonNode.Parse("""{"model":"qwen-test","vendor":"openai","system_prompt":"here","tools":[]}"""),
                ["log"] = JsonNode.Parse("""[{"seq":1,"type":"agent_shutdown","time":"2026-11-05T00:00:00Z","reason":"shutdown"}]""")
            }.ToJsonString());
        GoEnvironment(source, () =>
        {
            Check(Restore(quiet) is not null, "a fresh data directory did not adopt the Go primary");
            Check(File.Exists(Path.Combine(quiet, "agents", "root", "save.json")), "the adopted primary is not in place");
            Check(!File.Exists(Path.Combine(quiet, "agents", "Go archive", "save.json")),
                "a quiet home was filed as a peer instead of becoming the primary");
        });

        // A home that already talks keeps its work and receives the Go primary as a
        // named peer: a convenience at startup must never bury a live conversation.
        string busy = Path.Combine(temp, "live-busy");
        Directory.CreateDirectory(Path.Combine(busy, "agents", "root"));
        File.WriteAllText(Path.Combine(busy, "agents", "root", "save.json"),
            new JsonObject
            {
                ["context"] = null, ["as_of"] = 1,
                ["config"] = JsonNode.Parse("""{"model":"qwen-test","vendor":"openai","system_prompt":"here","tools":[]}"""),
                ["log"] = JsonNode.Parse("""[{"seq":1,"type":"message_received","time":"2026-11-05T00:00:00Z","message":{"actor":"human","parts":[{"type":"text","text":"work already underway here"}]}}]""")
            }.ToJsonString());
        GoEnvironment(source, () =>
        {
            Check(Restore(busy) is not null, "a busy data directory restored nothing");
            string kept = File.ReadAllText(Path.Combine(busy, "agents", "root", "save.json"));
            Check(kept.Contains("work already underway here"), "the restore overwrote the live primary conversation");
            Check(File.Exists(Path.Combine(busy, "agents", "Go archive", "save.json")),
                "the Go primary did not arrive as a peer");
        });
    }

    public static void Graceful(string temp)
    {
        // One unreadable peer must not cost the operator the rest of the fleet, and a
        // directory that is not Go state at all must leave the launch untouched.
        string source = Path.Combine(temp, "go-partial");
        WriteGoTree(source, brokenPeer: true);
        string destination = Path.Combine(temp, "live-partial");
        Directory.CreateDirectory(destination);
        GoEnvironment(source, () =>
        {
            string? summary = Restore(destination);
            Check(summary is not null && summary.Contains("1 Agent "),
                "one damaged peer stopped the whole restore: " + summary);
            Check(RestoreLog.Contains("PDF designer"), "the damaged peer was not reported: " + RestoreLog);
            Check(File.Exists(Path.Combine(destination, "agents", "root", "save.json")),
                "the healthy primary was withheld behind a damaged peer");
        });

        string plain = Path.Combine(temp, "not-go");
        Directory.CreateDirectory(plain);
        // Valid for this codec, but written here rather than by Go: a permissive
        // detector would happily adopt it as the primary and the operator would see
        // their own data directory arrive as a stranger in the tree.
        File.WriteAllText(Path.Combine(plain, "save.json"),
            @"{""context"":null,""as_of"":1,""config"":{""model"":""here"",""vendor"":""openai"",""system_prompt"":""here"",""tools"":[]},"
            + @"""log"":[{""seq"":1,""type"":""message_received"",""time"":""2026-11-05T00:00:00Z"",""message"":{""actor"":""human"",""parts"":[{""type"":""text"",""text"":""written here, not by Go""}]}}]}");
        string elsewhere = Path.Combine(temp, "live-plain");
        Directory.CreateDirectory(elsewhere);
        GoEnvironment(plain, () =>
        {
            // Without this the fallback would find the repository the suite runs in,
            // which really is a Go tree, and the check would prove nothing.
            Directory.SetCurrentDirectory(plain);
            Check(Restore(elsewhere) is null, "an ordinary data directory was treated as Go state");
            Check(!File.Exists(Path.Combine(elsewhere, "agents", "root", "save.json")),
                "a non-Go directory produced an Agent");
        });

        // Deleting an imported Agent stays deleted: the marker is what keeps a
        // start-up convenience from resurrecting what the operator removed.
        GoEnvironment(source, () =>
        {
            string marker = Path.Combine(destination, ".go-history-restored");
            Check(File.Exists(marker), "no marker recorded the restore");
            File.Delete(Path.Combine(destination, "agents", "root", "save.json"));
            Check(Restore(destination) is null, "a removed Agent was resurrected on the next launch");
            Check(!File.Exists(Path.Combine(destination, "agents", "root", "save.json")),
                "the marker did not suppress a re-restore");
        });

        // A damaged registry is never replaced by an import.
        string registry = Path.Combine(destination, "ensemble.json");
        File.WriteAllText(registry, "{ not json");
        File.WriteAllText(Path.Combine(destination, ".go-history-restored"), "again");
        GoEnvironment(source, () =>
        {
            Restore(destination);
            Check(File.ReadAllText(registry) == "{ not json", "discovery overwrote a damaged peer registry");
        });
    }
}
