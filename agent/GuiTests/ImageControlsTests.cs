using System.Net;
using System.Net.Sockets;
using System.Net.WebSockets;
using System.Text;
using System.Text.Json;
using Ensemble;
using EnsembleGui;
using SkiaSharp;

static class ImageControlsTests
{
    static void Check(bool condition, string message)
    {
        if (!condition) throw new InvalidOperationException(message);
    }

    public static async Task RunAsync()
    {
        string work = Path.Combine(Path.GetTempPath(), "gui-images-" + Guid.NewGuid().ToString("N"));
        Directory.CreateDirectory(Path.Combine(work, "skills", "plain"));
        File.WriteAllText(Path.Combine(work, "skills", "plain", "SKILL.md"),
            "---\nname: plain\ndescription: Plain\ntype: primary\n---\nPlain\n");
        Console.WriteLine("Image controls workspace: " + work);
        var app = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Path.Combine(work, "fleet"), agentDiagnostics: true);
        var agent = new Agent(app, new Configuration(Vendor.OpenAI, "gpt-6-astra", "http://127.0.0.1:1", surface: Surface.Responses),
            name: "images", skillSpec: new SkillSpec(Path.Combine(work, "skills"), "plain"));
        using var deadline = new CancellationTokenSource(TimeSpan.FromSeconds(90));
        var token = deadline.Token;
        await agent.SkillManager()!.InitializeAsync(token);
        agent.StartDiagnostics();
        await agent.PostAsync(new AgentInput(InputKind.Pause));
        using var portProbe = new TcpListener(IPAddress.Loopback, 0);
        portProbe.Start();
        int port = ((IPEndPoint)portProbe.LocalEndpoint).Port;
        portProbe.Stop();
        string url = $"http://127.0.0.1:{port}";
        var builder = WebApplication.CreateBuilder();
        builder.Logging.ClearProviders();
        builder.WebHost.UseUrls(url);
        await using var host = builder.Build();
        using var stop = new CancellationTokenSource();
        string tracePath = Path.Combine(agent.DataDirectory(), "gui.log");
        Task serving = Gui.ServeAsync(host, agent, new GuiOptions(port,
            SettingsPath: Path.Combine(work, "settings.json")), stop.Token);
        using var socket = new ClientWebSocket();
        socket.Options.SetRequestHeader("Origin", url);
        async Task<JsonElement> Frame(string type)
        {
            byte[] buffer = new byte[8192];
            while (true)
            {
                using var bytes = new MemoryStream();
                WebSocketReceiveResult received;
                do
                {
                    received = await socket.ReceiveAsync(buffer, token);
                    Check(received.MessageType == WebSocketMessageType.Text, "Socket closed unexpectedly.");
                    bytes.Write(buffer, 0, received.Count);
                } while (!received.EndOfMessage);
                using var json = JsonDocument.Parse(bytes.ToArray());
                if (json.RootElement.GetProperty("type").GetString() == type) return json.RootElement.Clone();
            }
        }
        Task Send(object value) => socket.SendAsync(JsonSerializer.SerializeToUtf8Bytes(value), WebSocketMessageType.Text, true, token);
        int Images() => agent.History().Events().SelectMany(e => e.Parts).OfType<BlobPart>().Count();
        try
        {
            await host.StartAsync(token);
            await socket.ConnectAsync(new Uri($"ws://127.0.0.1:{port}/ws"), token);
            var replay = await Frame("replay");
            string target = replay.GetProperty("agent_id").GetString()!;
            using var bitmap = new SKBitmap(new SKImageInfo(1024, 1024, SKColorType.Rgba8888, SKAlphaType.Unpremul));
            new Random(71).NextBytes(bitmap.GetPixelSpan());
            for (int i = 3; i < bitmap.ByteCount; i += 4) bitmap.GetPixelSpan()[i] = 255;
            using var image = SKImage.FromBitmap(bitmap);
            using var encoded = image.Encode(SKEncodedImageFormat.Png, 100);
            string source = Convert.ToBase64String(encoded.ToArray());
            Check(source.Length > 1024 * 1024, "Fixture did not exceed ordinary control limit.");
            object Envelope(long sequence, string address, string data) => new
            {
                type = "image_prompt", sequence, agent_id = address, text = "WITH_IMAGE_TEXT",
                images = new[] { new { mime = "image/png", data } }
            };
            await Send(Envelope(1, target, source));
            var accepted = await Frame("image_accepted");
            Check(accepted.GetProperty("sequence").GetInt64() == 1 && Images() == 1, "Image was not admitted exactly once.");
            var dimensions = accepted.GetProperty("images")[0];
            Check(dimensions.GetProperty("original_width").GetInt32() == 1024
                && dimensions.GetProperty("width").GetInt32() == 1024
                && !dimensions.GetProperty("resized").GetBoolean(), "Acknowledgment omitted preparation metadata.");
            var retained = agent.History().Events().SelectMany(e => e.Parts).OfType<BlobPart>().Single();
            using var decoded = SKBitmap.Decode(Convert.FromBase64String(retained.Ref.Locator.Split(',')[1]));
            Check(decoded.GetPixel(51, 77) == bitmap.GetPixel(51, 77), "Actual retained pixels differ.");
            Check(agent.History().Events().SelectMany(e => e.Parts).OfType<TextPart>().Any(p => p.Text == "WITH_IMAGE_TEXT"),
                "Accompanying text was lost.");
            async Task Refuse(object value)
            {
                await Send(value);
                await Frame("error");
                Check(Images() == 1, "Refused input added an image to History.");
            }
            await Refuse(Envelope(1, target, source));
            await Refuse(Envelope(2, "wrong-agent", source));
            await Refuse(Envelope(3, target, ""));
            await Refuse(Envelope(4, target, "AAAA"));
            await Refuse(new
            {
                type = "image_prompt", sequence = 5, agent_id = target,
                images = new[] { new { mime = "image/png", data = source }, new { mime = "image/png", data = "AAAA" } }
            });
            await Refuse(new { type = "image_prompt", sequence = 6, images = new[] { new { mime = "image/png", data = source } } });
            await Refuse(new { type = "prompt", agent_id = target, images = new[] { new { mime = "image/png", data = source } } });
            await Refuse(new { type = "prompt", agent_id = target, text = new string('x', 1024 * 1024 + 1) });
            await Send(Envelope(7, target, source));
            await Frame("image_accepted");
            Check(Images() == 2, "Refusals disrupted subsequent valid admission.");
            // A document is rasterised in the listener, so every page enters History as an
            // ordinary image part: no new wire shape and no new capability concept. The
            // fixture is built here rather than embedded so its page count is auditable.
            string document = Document(3);
            int pagesBefore = Images();
            await Send(new
            {
                type = "image_prompt", sequence = 8, agent_id = target, text = "WITH_DOCUMENT",
                images = new[] { new { mime = "application/pdf", data = document } }
            });
            var paged = await Frame("image_accepted");
            Check(Images() == pagesBefore + 3, "A three-page document did not enter History as three pages.");
            Check(paged.GetProperty("images").GetArrayLength() == 3, "Acknowledgment did not report every page.");
            var pageParts = agent.History().Events().SelectMany(e => e.Parts).OfType<BlobPart>().TakeLast(3).ToList();
            Check(pageParts.Count == 3, "History does not hold three page parts.");
            Check(pageParts.All(p => p.Mime == "image/png"), "Pages were retained under a non-image MIME.");
            Check(pageParts.All(p => Convert.FromBase64String(p.Ref.Locator.Split(',')[1])[1] == (byte)'P'),
                "A retained page is not a PNG.");
            Check(pageParts.Select(p => p.Ref.Locator).Distinct().Count() == 3, "Distinct pages collapsed into one part.");
            Check(agent.History().Events().SelectMany(e => e.Parts).OfType<TextPart>().Any(p => p.Text == "WITH_DOCUMENT"),
                "Accompanying document text was lost.");
            // Each page costs a slot, so a document longer than the message allowance is
            // refused outright instead of being silently truncated: an answer about
            // "page 5" of a truncated document would otherwise look entirely plausible.
            await Send(new
            {
                type = "image_prompt", sequence = 9, agent_id = target,
                images = new[] { new { mime = "application/pdf", data = Document(5) } }
            });
            var ceiling = await Frame("error");
            Check(ceiling.GetProperty("image_sequence").GetInt64() == 9
                && ceiling.GetProperty("message").GetString()!.Contains("5 pages"),
                "Page ceiling was not correlated to its sequence or its count never reached the browser.");
            Check(Images() == pagesBefore + 3, "An over-long document was admitted anyway.");
            await Send(new
            {
                type = "image_prompt", sequence = 10, agent_id = target,
                images = new[] { new { mime = "application/pdf", data = Convert.ToBase64String("not a document"u8.ToArray()) } }
            });
            await Frame("error");
            Check(Images() == pagesBefore + 3, "A body claiming to be a document was admitted without a header check.");
            // Fill the remaining count allowance with tiny prepared images;
            // the next WebSocket image must be refused with its sequence intact.
            using var smallBitmap = new SKBitmap(4, 4);
            smallBitmap.Erase(SKColors.Blue);
            using var smallImage = SKImage.FromBitmap(smallBitmap);
            using var smallPng = smallImage.Encode(SKEncodedImageFormat.Png, 100);
            var small = agent.PrepareImage(smallPng.ToArray()).Image;
            // Pages consumed slots above, so the fill count is derived rather than
            // hard-coded: the budget that matters is the shared 32-image History total.
            // The deficit is captured once; re-reading it each turn would shrink as
            // fast as i grows and the allowance would never be filled.
            for (int i = 0, deficit = 32 - Images(); i < deficit; i++)
                await agent.PostAsync(new AgentInput(InputKind.Message, "fill", [small])).WaitAsync(token);
            int beforeRefusal = agent.History().Events().Length;
            await Send(Envelope(11, target, source));
            var overflow = await Frame("error");
            Check(overflow.GetProperty("image_sequence").GetInt64() == 11
                && overflow.GetProperty("message").GetString()!.Contains("micro_handoff"),
                "Overflow refusal was not correlated or actionable.");
            Check(Images() == 32 && agent.History().Events().Length == beforeRefusal,
                "Overflow admitted an image or caption.");
            // Inbound envelopes and outbound History frames both pass through
            // the same trace projection. No pixel payload belongs in this log.
            string trace = File.ReadAllText(tracePath);
            Check(trace.Contains("media_redacted") && !trace.Contains(source[..128]), "Trace retained source image bytes.");
            Check(!trace.Contains(retained.Ref.Locator.Split(',')[1][..128]), "Trace retained prepared image bytes.");
            Check(!trace.Contains(document[..100]), "Trace retained the document's source bytes.");
        }
        finally
        {
            socket.Abort();
            stop.Cancel();
            try { await serving.WaitAsync(TimeSpan.FromSeconds(15)); }
            finally
            {
                await host.StopAsync(CancellationToken.None);
                app.Destroy();
            }
        }
        Console.WriteLine("PASS GUI images: payloads, source preparation, explicit routing, duplicate refusal, control limits and trace privacy");
    }

    // A minimal but well-formed PDF: real xref offsets and a Helvetica text object, so
    // poppler rasterises ink rather than a blank page. Page n reads "VERMILION PAGE n".
    private static string Document(int pages)
    {
        var objects = new List<string>
        {
            "<< /Type /Catalog /Pages 2 0 R >>",
            "<< /Type /Pages /Kids [" + string.Join(" ",
                Enumerable.Range(0, pages).Select(index => (4 + index) + " 0 R")) + "] /Count " + pages + " >>",
            "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>"
        };
        for (int page = 0; page < pages; page++)
        {
            objects.Add("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 300 100] /Resources << /Font << /F1 3 0 R >> >>"
                + " /Contents " + (4 + pages + page) + " 0 R >>");
        }
        for (int page = 0; page < pages; page++)
        {
            string stream = "BT /F1 20 Tf 10 30 Td (VERMILION PAGE" + (page + 1) + ") Tj ET";
            objects.Add("<< /Length " + Encoding.UTF8.GetByteCount(stream) + " >>\nstream\n" + stream + "\nendstream");
        }
        var body = new StringBuilder("%PDF-1.4\n");
        var offsets = new List<int>();
        for (int index = 0; index < objects.Count; index++)
        {
            offsets.Add(Encoding.UTF8.GetByteCount(body.ToString()));
            body.Append((index + 1) + " 0 obj\n" + objects[index] + "\nendobj\n");
        }
        int crossReference = Encoding.UTF8.GetByteCount(body.ToString());
        body.Append("xref\n0 " + (objects.Count + 1) + "\n0000000000 65535 f \n");
        foreach (int offset in offsets)
        {
            body.Append(offset.ToString("D10") + " 00000 n \n");
        }
        body.Append("trailer\n<< /Size " + (objects.Count + 1) + " /Root 1 0 R >>\nstartxref\n" + crossReference
            + "\n%%EOF\n");
        return Convert.ToBase64String(Encoding.UTF8.GetBytes(body.ToString()));
    }
}
