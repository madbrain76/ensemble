using Ensemble;
using System.Collections.Immutable;
using System.Diagnostics;
using System.Text;
using System.Text.Json;

namespace EnsembleGui;

public static partial class Gui
{
    private static partial async Task<(AgentInput Input, ImmutableArray<PreparedImage> Images)> PrepareImageInput(
        Agent selected, JsonElement command, CancellationToken cancellation)
    {
        string text = command.TryGetProperty("text", out var textValue) ? textValue.GetString() ?? "" : "";
        if (Encoding.UTF8.GetByteCount(text) > 32 * 1024)
            throw new InvalidDataException("Image message text exceeds 32 KiB.");
        var images = command.GetProperty("images");
        if (images.ValueKind != JsonValueKind.Array || images.GetArrayLength() is < 1 or > Agent.MaximumInputImages)
            throw new InvalidDataException("An image message requires one to four images.");
        var sources = new List<string>();
        var parts = new List<Part>();
        long total = 0;
        foreach (var image in images.EnumerateArray())
        {
            string mime = image.GetProperty("mime").GetString() ?? "";
            if (mime is not ("image/png" or "image/jpeg" or "image/webp" or "image/gif" or "application/pdf"))
                throw new InvalidDataException("Unsupported source attachment format.");
            string data = image.GetProperty("data").GetString() ?? "";
            if (data.Length == 0 || data.Length % 4 != 0)
                throw new InvalidDataException("Image payload must be unspaced base64.");
            long bytes = (long)data.Length / 4 * 3;
            if (data.EndsWith("==", StringComparison.Ordinal))
            {
                bytes -= 2;
            }
            else if (data.EndsWith('='))
            {
                bytes--;
            }
            total += bytes;
            if (total > Agent.MaximumImageSourceBytes)
                throw new InvalidDataException("Combined image sources exceed 64 MiB.");
            if (mime != "application/pdf")
            {
                sources.Add(data);
                parts.Add(new BlobPart(mime, new BlobRef(RefKind.Uri, $"data:{mime};base64,{data}")));
                continue;
            }
            // Each page costs the same as an attached image, so a document competes for
            // the very allowance an image would use rather than for an extra one.
            foreach (string page in PdfPages(data,
                         Agent.MaximumInputImages - sources.Count))
            {
                long pageBytes = (long)page.Length / 4 * 3;
                total += pageBytes;
                if (total > Agent.MaximumImageSourceBytes)
                    throw new InvalidDataException("Attached document pages exceed 64 MiB.");
                sources.Add(page);
                parts.Add(new BlobPart("image/png", new BlobRef(RefKind.Uri, $"data:image/png;base64,{page}")));
            }
        }
        // A model the curated table does not name proves itself here, once, on
        // this async boundary rather than inside rendering. Refusing bytes the
        // endpoint would have accepted is the failure this prevents.
        await selected.Engine().ProbeImagesAsync(cancellation);
        // Refuse provider capabilities before raster decoding. The same
        // renderer validator remains authoritative for admitted image parts.
        selected.Engine().ValidateMedia(parts, Actor.Human);
        var prepared = ImmutableArray.CreateBuilder<PreparedImage>();
        int outputBytes = 0;
        foreach (string source in sources)
        {
            cancellation.ThrowIfCancellationRequested();
            if (source.Any(char.IsWhiteSpace))
                throw new InvalidDataException("Image payload must be unspaced base64.");
            byte[] bytes = Convert.FromBase64String(source);
            var image = selected.PrepareImage(bytes);
            outputBytes += image.OutputBytes;
            if (outputBytes > Agent.MaximumInputImageBytes)
                throw new InvalidDataException("Prepared images exceed the combined 16 MiB input limit.");
            cancellation.ThrowIfCancellationRequested();
            prepared.Add(image);
        }
        var imagesForInput = prepared.Select(image => image.Image).ToImmutableArray();
        return (new AgentInput(InputKind.Message, text, imagesForInput), prepared.ToImmutable());
    }

    // One document, one rasterisation. pdfinfo settles the page count first so an
    // over-long document is refused with a usable message instead of being silently
    // truncated, which would otherwise make an answer about "page 3 of 40" look correct.
    private static IEnumerable<string> PdfPages(string source, int allowance)
    {
        if (allowance < 1)
            throw new ImageBudgetExceededException("A document's pages each count as an attached image, and this message"
                + " has no free image slot. Remove an image or send the page separately.");
        byte[] bytes;
        try
        {
            bytes = Convert.FromBase64String(source);
        }
        catch (FormatException)
        {
            throw new InvalidDataException("Document payload must be unspaced base64.");
        }
        if (bytes.Length < 5 || Encoding.ASCII.GetString(bytes, 0, 5) != "%PDF-")
            throw new InvalidDataException("Document is not a PDF.");
        string directory = Path.Combine(Path.GetTempPath(), "ensemble-pdf-" + Guid.NewGuid().ToString("N"));
        Directory.CreateDirectory(directory);
        try
        {
            string document = Path.Combine(directory, "document.pdf");
            File.WriteAllBytes(document, bytes);
            int pages = PdfPageCount(RunPdfTool("pdfinfo", [document], directory));
            if (pages > allowance)
                throw new ImageBudgetExceededException($"Document has {pages} pages but only {allowance} image slot"
                    + (allowance == 1 ? "" : "s") + " remain in this message. Print a page range, or send the pages separately.");
            string prefix = Path.Combine(directory, "page");
            RunPdfTool("pdftoppm", ["-png", "-r", "150", document, prefix], directory);
            var pagesRendered = Directory.GetFiles(directory, "page*.png").OrderBy(name => name, StringComparer.Ordinal)
                .Select(page => Convert.ToBase64String(File.ReadAllBytes(page))).ToList();
            if (pagesRendered.Count != pages)
                throw new InvalidDataException($"Document reported {pages} pages but {pagesRendered.Count} rendered.");
            return pagesRendered;
        }
        catch (Exception error) when (error is not InvalidDataException && error is not ImageBudgetExceededException)
        {
            // A missing or failing poppler says nothing about the document, but it must
            // not surface as an unhandled crash in the request path either.
            throw new InvalidDataException("Document could not be rasterised: " + error.Message);
        }
        finally
        {
            try
            {
                Directory.Delete(directory, true);
            }
            catch (IOException)
            {
                // A leaked temporary directory is preferable to failing the attachment.
            }
            catch (UnauthorizedAccessException)
            {
            }
        }
    }

    private static int PdfPageCount(string information)
    {
        foreach (string line in information.Split('\n'))
        {
            string trimmed = line.Trim();
            const string label = "Pages:";
            if (trimmed.StartsWith(label, StringComparison.OrdinalIgnoreCase))
            {
                // The prefix is dropped before digits are taken: the value follows
                // a run of padding spaces, so a scan that only skips leading
                // whitespace from the start of the line finds none and reads "".
                string digits = new(trimmed[label.Length..].TrimStart().TakeWhile(char.IsDigit).ToArray());
                if (int.TryParse(digits, out int count) && count >= 1)
                {
                    return count;
                }
            }
        }
        throw new InvalidDataException("Document page count could not be read.");
    }

    // Conversion runs in the listener, not in a tool sandbox: this is the transport
    // preparing a request body, not the model choosing to run a command.
    private static string RunPdfTool(string program, string[] arguments, string directory)
    {
        var start = new ProcessStartInfo(program)
        {
            WorkingDirectory = directory,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
        };
        foreach (string argument in arguments)
        {
            start.ArgumentList.Add(argument);
        }
        using var process = Process.Start(start) ?? throw new InvalidOperationException(program + " could not start");
        // Both pipes are drained concurrently. Reading them one after the other
        // deadlocks once a child writes enough to the unread one to fill its buffer.
        Task<string> errors = process.StandardError.ReadToEndAsync();
        string output = process.StandardOutput.ReadToEnd();
        string written = errors.GetAwaiter().GetResult();
        if (!process.WaitForExit(30_000))
        {
            try
            {
                process.Kill(entireProcessTree: true);
            }
            catch (InvalidOperationException)
            {
            }
            throw new TimeoutException(program + " timed out");
        }
        if (process.ExitCode != 0)
            throw new InvalidOperationException(program + " exited with " + process.ExitCode + ": " + written.Trim());
        return output;
    }

    private static partial string TraceProjection(string json)
    {
        using var document = JsonDocument.Parse(json);
        bool ContainsImage(JsonElement value)
        {
            if (value.ValueKind == JsonValueKind.Object)
            {
                foreach (var field in value.EnumerateObject())
                {
                    bool imageMime = field.Name is "mime" or "mimeType" && field.Value.ValueKind == JsonValueKind.String
                        && (field.Value.GetString() ?? "").StartsWith("image/", StringComparison.OrdinalIgnoreCase);
                    // MCP payload privacy must not depend on valid MIME/type
                    // metadata: malformed image replies are traced before refusal.
                    bool opaqueData = field.Name == "data" && field.Value.ValueKind == JsonValueKind.String;
                    if (field.Name is "images" or "audio" || imageMime || opaqueData)
                    {
                        return true;
                    }
                    if (ContainsImage(field.Value))
                    {
                        return true;
                    }
                }
            }
            else if (value.ValueKind == JsonValueKind.Array)
            {
                foreach (var item in value.EnumerateArray())
                {
                    if (ContainsImage(item))
                    {
                        return true;
                    }
                }
            }
            else if (value.ValueKind == JsonValueKind.String
                && (value.GetString() ?? "").Contains("data:image/", StringComparison.OrdinalIgnoreCase))
            {
                return true;
            }
            return false;
        }
        if (!ContainsImage(document.RootElement)) return json;
        string type = document.RootElement.TryGetProperty("type", out var value) && value.ValueKind == JsonValueKind.String
            ? value.GetString() ?? "" : "";
        // Do not persist any of the frame, including nested History/tool media.
        return JsonSerializer.Serialize(new { type, media_redacted = true, utf8_bytes = Encoding.UTF8.GetByteCount(json) });
    }
}
