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

    public static void Run(string temp)
    {
        string source = Path.Combine(temp, "go");
        string destination = Path.Combine(temp, "live");
        string peerHome = Path.Combine(source, "agents", "PDF designer");
        Directory.CreateDirectory(peerHome);
        Directory.CreateDirectory(destination);
        File.WriteAllText(Path.Combine(source, "save.json"), GoSave(17, Primary));
        File.WriteAllText(Path.Combine(source, "ensemble.json"),
            """[{"name":"PDF designer","model":"qwen-test","skill":"","created":"2026-10-25T04:55:26.331352951Z"}]""");
        File.WriteAllText(Path.Combine(peerHome, "save.json"), GoSave(1,
            """{"seq":1,"type":"message_received","time":"2026-10-25T05:00:00Z","message":{"actor":"human","parts":[{"type":"text","text":"peer question"}]}}"""));
        byte[] before = File.ReadAllBytes(Path.Combine(source, "save.json"));

        var output = new StringWriter();
        var errors = new StringWriter();
        int code = Application.RunGoImportCommand(
            ["--from", source, "--data-dir", destination], output, errors);
        Check(code == 0, "import refused a valid Go tree: " + errors + output);
        string report = output.ToString();
        Check(report.Contains("dropped 1 job_finished") && report.Contains("tools_changed"),
            "drops were not reported by name: " + report);
        Check(report.Contains("1 recall_attached to recall_added")
            && report.Contains("1 job_report to tool_observed"),
            "a translation was not reported: " + report);
        Check(report.Contains("1 text parts carried no body"), "an absent part body was not reported: " + report);

        // Migration must be read-only on the tree it comes from.
        Check(before.AsSpan().SequenceEqual(File.ReadAllBytes(Path.Combine(source, "save.json"))),
            "import-go modified the Go source files");

        var imported = JsonCodec.ReadSave(File.ReadAllText(Path.Combine(destination, "agents", "root", "save.json")));
        Check(imported.Events.Length == 15, $"expected the fifteen readable events, got {imported.Events.Length}");
        Check(imported.Events.Any(e => e.Kind == EventKind.ToolObserved
            && e.Job is JobInfo report && report.Handle == 9
            && report.OutputPath.EndsWith("io/9", StringComparison.Ordinal)),
            "the peer job report lost its outcome or its output locator");
        Check(imported.Events.Any(e => e.Kind == EventKind.MessageReceived
            && e.Parts.Length == 1 && e.Parts[0] is TextPart ask && ask.Text.Contains("slip-in PDF cover")),
            "the user's own words did not survive");
        Check(imported.Events.Any(e => e.Kind == EventKind.ToolReturned
            && e.Parts.Any(part => part is TextPart result && result.Text.Contains("blender.pdf"))),
            "tool output did not survive");
        Check(imported.Events.Any(e => e.Kind == EventKind.TurnEnded && e.Reason is not null),
            "turn end still has no reason");
        var recall = imported.Events.Single(e => e.Kind == EventKind.RecallAdded);
        Check(recall.Parts.Length == 1 && recall.Parts[0] is RecallPart attached
            && attached.Attachment.Text.Contains("blender notes")
            && attached.Attachment.Hashes.SequenceEqual(["3b38c1c92b6e1add6ee7cdd0bf5d7e1a152b60be27b9bf03e4704a3851b8bf99"]),
            "the recalled memories lost their text or their hashes");
        // Renumbering is only honest if a redaction still names the records it
        // covered: the range was written against Go's sequences, two of which the
        // migration retired.
        var redaction = imported.Events.Single(e => e.Kind == EventKind.Redacted);
        Check(redaction.RedactFrom == 10 && redaction.RedactTo == 12
            && redaction.Level == Redaction.Dialogue,
            $"the migrated redaction range drifted: {redaction.RedactFrom}..{redaction.RedactTo}, wanted 10..12");
        Check(imported.Events[^1].Seq == imported.Events.Length && imported.Sequence == imported.Events.Length,
            "the migrated history is not an unbroken run from one");
        Check(imported.Context is null, "Go's snapshot was carried over instead of rebuilt");

        // The real test of an import is that an Agent opens it and can continue.
        string original = Directory.GetCurrentDirectory();
        try
        {
            Directory.SetCurrentDirectory(destination);
            var trace = new StringWriter();
            var app = new Application(trace, dataDirectory: destination);
            try
            {
                var agent = new Agent(app, new Configuration(Vendor.OpenAI, "qwen-test", "http://127.0.0.1:1"),
                    name: "root");
                try
                {
                    bool restored;
                    try
                    {
                        restored = agent.StartPersistence(new PersistenceOptions(Path.Combine("agents", "root", "save.json")));
                    }
                    catch (Exception fault)
                    {
                        throw new Exception($"restore failed: {fault.Message} / {trace}", fault);
                    }
                    Check(restored, "an imported Agent could not restore its history: " + trace);
                    string wire = agent.Engine().Render(agent.History().View()).ToJsonString();
                    Check(wire.Contains("design a slip-in PDF cover") && wire.Contains("blender.pdf"),
                        "the resumed conversation is missing its substance");
                    // The system prompt belongs to the live configuration, not to a
                    // saved history: Go recorded it for reproducibility, and this
                    // build composes it from the launcher and its skills instead.
                    // The call identity is what ties a result to the model's request
                    // for it; losing it would leave a tool answer belonging to nobody.
                    Check(wire.Contains("call-1"), "the imported call identity was lost");
                    Check(!wire.Contains("the folder holds manuals"),
                        "a record the Go history redacted with dialogue scope came back into view");
                }
                finally
                {
                    agent.Destroy();
                }
            }
            finally
            {
                app.Destroy();
            }
        }
        finally
        {
            Directory.SetCurrentDirectory(original);
        }

        var again = Application.RunGoImportCommand(["--from", source, "--data-dir", destination], new StringWriter(), errors);
        Check(again == 1 && errors.ToString().Contains("--replace"),
            "a second import overwrote or hid the existing history: " + errors);
    }

    public static void Refusals(string temp)
    {
        string source = Path.Combine(temp, "go-refuse");
        string destination = Path.Combine(temp, "live-refuse");
        Directory.CreateDirectory(source);
        Directory.CreateDirectory(destination);
        var output = new StringWriter();
        var errors = new StringWriter();

        // A history this build cannot read must not replace one it could.
        // A request without its target is damage, not an unfamiliar event type:
        // it must refuse the import rather than be quietly retired with the rest.
        File.WriteAllText(Path.Combine(source, "save.json"), GoSave(1,
            """{"seq":1,"type":"request_sent","time":"2026-10-25T04:00:00Z"}"""));
        Check(Application.RunGoImportCommand(["--from", source, "--data-dir", destination], output, errors) == 1,
            "an unreadable save was accepted");
        Check(!File.Exists(Path.Combine(destination, "agents", "root", "save.json")),
            "a refused import still wrote a save file");

        // A dry run proves the plan without touching the destination.
        File.WriteAllText(Path.Combine(source, "save.json"), GoSave(1, Primary[0]));
        Check(Application.RunGoImportCommand(["--from", source, "--data-dir", destination, "--dry-run"], output, new StringWriter()) == 0,
            "a valid dry run failed");
        Check(!Directory.Exists(Path.Combine(destination, "agents")), "a dry run wrote history");

        // A damaged destination registry is never replaced by an import.
        string registry = Path.Combine(destination, "ensemble.json");
        File.WriteAllText(registry, "{ not json");
        Check(Application.RunGoImportCommand(["--from", source, "--data-dir", destination], output, errors) == 0,
            "an import failed only because the registry needed repair");
        Check(File.ReadAllText(registry) == "{ not json", "import-go overwrote a damaged peer registry");
    }
}
