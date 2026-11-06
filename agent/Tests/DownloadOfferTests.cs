using System.Text.Json;
using Ensemble;

// Offering a finished work product: the descriptor the tool publishes, and the
// files it must refuse to publish.
static class DownloadOfferTests
{
    static void Check(bool value, string message)
    {
        if (!value) throw new Exception(message);
    }

    public static async Task Run(string temp)
    {
        var root = new Ensemble.Ensemble(TextWriter.Null,
            dataDirectory: Path.Combine(temp, "downloads"));
        var agent = new Agent(root, new Configuration(Vendor.OpenAI, "gpt-4o",
            "http://127.0.0.1:1"), name: "downloads");
        agent.InstallTools();
        Directory.SetCurrentDirectory(agent.DataDirectory());
        const string payload = "%PDF-1.4 fake but nonempty";
        File.WriteAllText("report.pdf", payload);
        File.WriteAllText("notes.html", "<script>alert(1)</script>");
        File.WriteAllText("empty.txt", "");
        Directory.CreateDirectory("build");
        var origin = agent.Configuration().Target();

        async Task<(string Text, bool Error)> Offer(string path)
        {
            var result = await agent.DispatchAsync(new ToolCallPart("call-" + path, origin,
                "offer_download", JsonSerializer.SerializeToElement(new { path })),
                CancellationToken.None);
            return (string.Concat(result.Parts.OfType<TextPart>().Select(part => part.Text)),
                result.IsError);
        }

        var good = await Offer("report.pdf");
        Check(!good.Error, "a real work product was refused: " + good.Text);
        using (var document = JsonDocument.Parse(good.Text))
        {
            var offer = document.RootElement.GetProperty("download");
            Check(offer.GetProperty("name").GetString() == "report.pdf"
                && offer.GetProperty("mime").GetString() == "application/pdf"
                && offer.GetProperty("size").GetInt64() == payload.Length,
                "the descriptor is wrong: " + good.Text);
        }
        Check((await Offer("missing.pdf")).Error, "a missing file was offered");
        Check((await Offer("empty.txt")).Error, "an empty file was offered");
        Check((await Offer("build")).Error, "a directory was offered");
        Check((await Offer("notes.html")).Error, "HTML was offered for download");
        Check((await Offer("")).Error, "an empty path was accepted");
        // The whitelist is shared with the browser route, so one list decides both.
        Check(Tool.DownloadMime("archive.tar.gz") == "application/gzip", "gzip was mislabelled");
        Check(Tool.DownloadMime("plot.svg") == "image/svg+xml", "svg was mislabelled");
        try
        {
            _ = Tool.DownloadMime("installer.sh");
            Check(false, "a shell script passed the download whitelist");
        }
        catch (ArgumentException)
        {
        }
        agent.Destroy();
        root.Destroy();
    }
}
