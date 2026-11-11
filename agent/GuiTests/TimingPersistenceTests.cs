using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;
using EnsembleGui;
using Application = Ensemble.Ensemble;

// Artifact clocks must come from the recorded event, never from the moment a
// browser happened to open. This drives a real streamed turn and follows the
// measured span through History, the save codec and the replay projection.
static class TimingPersistenceTests
{
    static void Check(bool condition, string message)
    {
        if (!condition) throw new InvalidOperationException(message);
    }

    static async Task<string> ReadRequest(Stream stream, CancellationToken token)
    {
        var buffer = new byte[8192];
        var text = new StringBuilder();
        int prior = 0;
        while (true)
        {
            int read = await stream.ReadAsync(buffer, token);
            if (read == 0) break;
            text.Append(Encoding.UTF8.GetString(buffer, 0, read));
            string value = text.ToString();
            int body = value.IndexOf("\r\n\r\n", StringComparison.Ordinal);
            if (body >= 0)
            {
                int length = value[..body].Split("\r\n")
                    .Where(line => line.StartsWith("Content-Length:", StringComparison.OrdinalIgnoreCase))
                    .Select(line => int.Parse(line.Split(':')[1].Trim())).First();
                int have = Encoding.UTF8.GetByteCount(value[(body + 4)..]);
                if (have >= length) return value[(body + 4)..];
            }
            prior += read;
        }
        throw new InvalidOperationException("Request ended early after " + prior + " bytes");
    }

    // One Anthropic SSE reply split across two writes so the fragments really do
    // arrive apart; a single-shot JSON reply shares `body` and has no fragments.
    static async Task<(Application app, Agent agent, string dir)> Streamed(string marker,
        bool streaming, CancellationToken token)
    {
        string directory = Directory.CreateTempSubdirectory("ensemble-timing-").FullName;
        var listener = new TcpListener(IPAddress.Loopback, 0);
        listener.Start();
        int port = ((IPEndPoint)listener.LocalEndpoint).Port;
        string first = "data: {\"type\":\"message_start\",\"message\":{\"model\":\"claude-3-haiku-20240307\",\"usage\":{\"input_tokens\":3}}}\n\n"
            + "data: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\",\"text\":\"\"}}\n\n"
            + "data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"opening \"}}\n\n";
        string tail = "data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"" + marker + "\"}}\n\n"
            + "data: {\"type\":\"content_block_stop\",\"index\":0}\n\n"
            + "data: {\"type\":\"message_delta\",\"usage\":{\"output_tokens\":2}}\n\n"
            + "data: {\"type\":\"message_stop\"}\n\n";
        string body = "{\"type\":\"message\",\"role\":\"assistant\",\"content\":[{\"type\":\"text\",\"text\":\"plain " + marker + "\"}],"
            + "\"usage\":{\"input_tokens\":3,\"output_tokens\":2},\"stop_reason\":\"end_turn\"}";
        _ = Task.Run(async () =>
        {
            using var client = await listener.AcceptTcpClientAsync(token);
            using var stream = client.GetStream();
            await ReadRequest(stream, token);
            string payload = streaming ? first + tail : body;
            int size = Encoding.UTF8.GetByteCount(payload);
            string head = "HTTP/1.1 200 OK\r\nContent-Type: " + (streaming ? "text/event-stream" : "application/json")
                + "\r\nContent-Length: " + size + "\r\nConnection: close\r\n\r\n";
            if (streaming)
            {
                await stream.WriteAsync(Encoding.UTF8.GetBytes(head + first), token);
                await stream.FlushAsync(token);
                await Task.Delay(120, token);
                await stream.WriteAsync(Encoding.UTF8.GetBytes(tail), token);
            }
            else
            {
                await stream.WriteAsync(Encoding.UTF8.GetBytes(head + payload), token);
            }
            await stream.FlushAsync(token);
        });
        var app = new Application(TextWriter.Null, dataDirectory: directory);
        var agent = new Agent(app, new Configuration(Vendor.Anthropic, "claude-3-haiku-20240307",
            "http://127.0.0.1:" + port, "timing-fixture-key"), name: "timed-root");
        return (app, agent, directory);
    }

    static JsonObject? ReplyFrame(Projection projection, string marker)
    {
        foreach (var frame in projection.Frames)
        {
            if (frame["type"]?.GetValue<string>() != "part_final") continue;
            if (frame["part"]?["text"]?.GetValue<string>()?.Contains(marker) == true) return frame;
        }
        return null;
    }

    public static async Task RunAsync()
    {
        using var deadline = new CancellationTokenSource(TimeSpan.FromSeconds(90));
        var token = deadline.Token;

        // 1. A streamed reply records its own span, survives the codec, and reaches
        //    the browser with the recorded instant rather than a page-load clock.
        var (app, agent, dir) = await Streamed("STREAM-MARKER", streaming: true, token);
        try
        {
            Check(await agent.SendAsync("say it", token).WaitAsync(token) == "opening STREAM-MARKER",
                "Streamed reply did not complete");
            var ended = agent.History().Events().Last(fact => fact.Kind == EventKind.ResponseEnded);
            Check(!ended.PartSpans.IsDefaultOrEmpty,
                "the streamed reply recorded no production span");
            Check(ended.PartSpans[0].Started > 0 && ended.PartSpans[0].Elapsed >= 80,
                $"production span is not a real measurement: {ended.PartSpans[0]}");

            var restored = JsonCodec.ReadEvent(JsonCodec.WriteEvent(ended));
            Check(restored.PartSpans.Length == ended.PartSpans.Length
                && restored.PartSpans[0].Started == ended.PartSpans[0].Started
                && restored.PartSpans[0].Elapsed == ended.PartSpans[0].Elapsed,
                "production span did not survive the save codec");
            Check(restored.Time == ended.Time, "recorded instant did not survive the save codec");

            var frame = ReplyFrame(Gui.Project(agent.History().Snapshot(), "timed-root"), "STREAM-MARKER");
            Check(frame is not null, "reply never reached the replay projection");
            Check(frame!["recorded_at"]?.GetValue<string>() == ended.Time?.ToString("O"),
                $"frame instant is not the recorded event time: {frame!["recorded_at"]}");
            Check(frame["part"]?["started"]?.GetValue<long>() == ended.PartSpans[0].Started
                && frame["part"]?["elapsed"]?.GetValue<long>() == ended.PartSpans[0].Elapsed,
                $"replay lost the measured span: {frame["part"]}");
        }
        finally { app.Destroy(); Directory.Delete(dir, recursive: true); }

        // 2. A reply that never streamed has no fragments to time, so it must show
        //    a clock alone. Claiming a zero-length production would be invented data.
        var (plainApp, plain, plainDir) = await Streamed("PLAIN-MARKER", streaming: false, token);
        try
        {
            await plain.SendAsync("say it", token).WaitAsync(token);
            var ended = plain.History().Events().Last(fact => fact.Kind == EventKind.ResponseEnded);
            Check(ended.PartSpans.IsDefaultOrEmpty,
                "a reply with no observed fragments claimed a measured span");
            var frame = ReplyFrame(Gui.Project(plain.History().Snapshot(), "timed-root"), "PLAIN-MARKER");
            Check(frame is not null && frame["recorded_at"] is not null,
                "a non-streamed reply lost its recorded instant too");
            Check(frame!["part"] is not null && frame["part"]!["started"] is null
                && frame["part"]!["elapsed"] is null,
                $"unmeasured production was published as a number: {frame["part"]}");
        }
        finally { plainApp.Destroy(); Directory.Delete(plainDir, recursive: true); }
        Console.WriteLine("PASS timing persistence: streamed span recorded, saved, and replayed with its original clock; unstreamed reply shows a clock alone");
    }
}
