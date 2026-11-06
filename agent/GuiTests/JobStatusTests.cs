using System.Net;
using System.Net.Sockets;
using System.Net.WebSockets;
using System.Text.Json;
using Ensemble;
using EnsembleGui;

// External public-API consumer. A real Job runs a real process while the GUI
// reports it; nothing here is inferred from a stub.
static class JobStatusTests
{
    static void Check(bool condition, string message)
    {
        if (!condition) throw new InvalidOperationException(message);
    }

    // The tree is only queued when it actually changed, so a waiter must keep
    // reading frames until the value it wants is published.
    static async Task<JsonElement> RowWhen(ClientWebSocket socket, string name, int jobs,
        string? state, CancellationToken token)
    {
        byte[] buffer = new byte[65536];
        while (true)
        {
            using var bytes = new MemoryStream();
            WebSocketReceiveResult result;
            do
            {
                result = await socket.ReceiveAsync(buffer, token);
                Check(result.MessageType == WebSocketMessageType.Text, "Unexpected socket close");
                bytes.Write(buffer, 0, result.Count);
            } while (!result.EndOfMessage);
            using var document = JsonDocument.Parse(bytes.ToArray());
            var root = document.RootElement;
            if (root.GetProperty("type").GetString() != "agent_tree")
            {
                continue;
            }
            foreach (var row in root.GetProperty("agents").EnumerateArray())
            {
                if (row.GetProperty("id").GetString() == name && row.GetProperty("jobs").GetInt32() == jobs
                    && (state is null || row.GetProperty("state").GetString() == state))
                {
                    return row.Clone();
                }
            }
        }
    }

    public static async Task RunAsync()
    {
        string directory = Directory.CreateTempSubdirectory("ensemble-job-status-gui-").FullName;
        string skills = Path.Combine(directory, "skills");
        Directory.CreateDirectory(Path.Combine(skills, "shell"));
        File.WriteAllText(Path.Combine(skills, "shell", "SKILL.md"),
            "---\nname: shell\ntype: primary\ndescription: Job status GUI fixture\n"
            + "tools: run_command wait_for_job\n---\n$TOOLS");
        using var deadline = new CancellationTokenSource(TimeSpan.FromSeconds(90));
        var token = deadline.Token;
        var app = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: directory);
        var config = new Configuration(Vendor.OpenAI, "offline", baseUrl: "http://127.0.0.1:1");
        var agent = new Agent(app, config, name: "shell-root", skillSpec: new SkillSpec(skills, "shell"));

        using var listener = new TcpListener(IPAddress.Loopback, 0);
        listener.Start();
        int port = ((IPEndPoint)listener.LocalEndpoint).Port;
        listener.Stop();
        string origin = $"http://127.0.0.1:{port}";
        using var stop = new CancellationTokenSource();
        Task serving = Gui.RunAsync(agent, new GuiOptions(port,
            SettingsPath: Path.Combine(directory, "appearance.json"),
            TracePath: Path.Combine(directory, "gui.log")), stop.Token);
        try
        {
            using var http = new HttpClient();
            while (true)
            {
                if (serving.IsCompleted) await serving;
                try
                {
                    if ((await http.GetAsync(origin, token)).IsSuccessStatusCode) break;
                }
                catch (HttpRequestException) { }
                await Task.Delay(20, token);
            }
            using var socket = new ClientWebSocket();
            socket.Options.SetRequestHeader("Origin", origin);
            await socket.ConnectAsync(new Uri($"ws://127.0.0.1:{port}/ws"), token);

            // Before any command exists, the tree reports no running work.
            await RowWhen(socket, "shell-root", 0, null, token);

            var call = new ToolCallPart("run-slow", config.Target(), "run_command",
                JsonSerializer.SerializeToElement(new { command = "sleep 8 && echo JOBS_DONE" }));
            ToolResultPart report = await agent.DispatchAsync(call, token);
            string text = JsonCodec.Text(report.Parts);
            // The report is a live Job, not completion: this is precisely the
            // window in which the old header claimed idle while work continued.
            Check(text.Contains("status: running"), "run_command did not hand back a running Job: " + text);

            // The exact window Bill saw: the report came back, the turn is over and
            // the actor is honestly idle, while the process it started still runs.
            var running = await RowWhen(socket, "shell-root", 1, "idle", token);
            Check(running.GetProperty("jobs").GetInt32() == 1, "Running Job was not counted");

            var finished = await RowWhen(socket, "shell-root", 0, "idle", token);
            Check(finished.GetProperty("id").GetString() == "shell-root",
                "The count returned on a different row");
        }
        finally
        {
            stop.Cancel();
            try { await serving.WaitAsync(TimeSpan.FromSeconds(10)); } catch { }
            agent.Destroy();
            app.Destroy();
        }
        Console.WriteLine("PASS job status: running Job counted through the tree, then released when the process exits");
    }
}
