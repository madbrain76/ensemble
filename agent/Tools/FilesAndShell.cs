using System.Diagnostics;
using System.Text;
using System.Text.Json;
using System.Text.RegularExpressions;

namespace Ensemble;

public sealed partial class Tool
{
    private string ReadFile(JsonElement args)
    {
        var start = Integer(args, "start_line", 1);
        var end = Integer(args, "end_line", int.MaxValue);
        var maximum = Integer(args, "max_bytes", int.MaxValue);
        if (start == 0)
        {
            start = 1;
        }
        if (end == 0)
        {
            end = int.MaxValue;
        }
        if (maximum == 0)
        {
            maximum = int.MaxValue;
        }
        if (maximum < 0)
        {
            throw new ArgumentException("max_bytes must be nonnegative.");
        }
        if (start < 1 || end < start)
        {
            throw new ArgumentException("Line numbers must be positive, with end_line >= start_line.");
        }
        using var input = Parent().Sandbox().OpenFile(String(args, "path"), FileMode.Open, FileAccess.Read);
        using var reader = new StreamReader(input);
        var text = reader.ReadToEnd();
        // Include each line's original terminator. Splitting on '\n' and joining
        // would silently normalize CRLF or add a newline to an unterminated file.
        var result = new StringBuilder();
        int offset = 0;
        int line = 1;
        while (offset < text.Length && line <= end)
        {
            var newline = text.IndexOf('\n', offset);
            var limit = newline < 0 ? text.Length : newline + 1;
            if (line >= start)
            {
                result.Append(text, offset, limit - offset);
            }
            offset = limit;
            line++;
        }
        var bytes = Encoding.UTF8.GetBytes(result.ToString());
        if (bytes.Length <= maximum)
        {
            return result.ToString();
        }
        // Back up to a UTF-8 boundary: a byte cap must not manufacture a broken
        // character. The report is separate from the capped file payload.
        var boundary = maximum;
        while (boundary > 0 && (bytes[boundary] & 0xc0) == 0x80)
        {
            boundary--;
        }
        return Encoding.UTF8.GetString(bytes, 0, boundary) + "\n[truncated]\n";
    }

    // A work product the user asked for becomes downloadable: the tool only
    // publishes a descriptor, and the browser route re-opens the file through the
    // same sandbox, so a name here is never a path the client may choose.
    private string OfferDownload(JsonElement args)
    {
        var path = String(args, "path", "");
        if (path.Trim().Length == 0)
        {
            throw new ArgumentException("Missing argument: path");
        }
        var sandbox = Parent().Sandbox();
        using (var probe = sandbox.OpenRegularRead(path))
        {
            if (probe.Length == 0)
            {
                throw new IOException("Work product is empty: " + path);
            }
            string name = Path.GetFileName(path.TrimEnd('/', '\\'));
            if (name.Length == 0)
            {
                throw new ArgumentException("Work product must be a file path: " + path);
            }
            long bytes = probe.Length;
            return JsonSerializer.Serialize(new
            {
                download = new
                {
                    name,
                    path,
                    mime = DownloadMime(name),
                    size = bytes
                }
            });
        }
    }

    // Only types a browser accepts as a file download without sniffing. An
    // executable or HTML deliverable stays a workspace file: the operator can
    // still read it, and the GUI never becomes a distribution point.
    public static partial string DownloadMime(string name) => Path.GetExtension(name).ToLowerInvariant() switch
    {
        ".pdf" => "application/pdf",
        ".txt" or ".log" or ".md" => "text/plain; charset=utf-8",
        ".csv" => "text/csv",
        ".json" => "application/json",
        ".zip" => "application/zip",
        ".tar" => "application/x-tar",
        ".gz" => "application/gzip",
        ".png" => "image/png",
        ".jpg" or ".jpeg" => "image/jpeg",
        ".gif" => "image/gif",
        ".webp" => "image/webp",
        ".svg" => "image/svg+xml",
        ".mp3" => "audio/mpeg",
        ".wav" => "audio/wav",
        ".ogg" => "audio/ogg",
        ".mp4" => "video/mp4",
        ".webm" => "video/webm",
        _ => throw new ArgumentException(
            "Work product type is not offered for download: " + Path.GetExtension(name))
    };

    private string ListDirectory(JsonElement args)
    {
        var directory = String(args, "path", ".");
        var files = Parent().Sandbox();
        var result = new StringBuilder();
        foreach (string name in files.ReadDirectory(directory))
        {
            result.Append(name);
            try
            {
                if (files.IsDirectory(Path.Combine(directory, name)))
                {
                    result.Append('/');
                }
            }
            catch (IOException)
            {
                // A vanished entry or an inaccessible/outward symlink remains
                // a visible name. Decoration must not bypass the boundary.
            }
            result.Append('\n');
        }
        return result.ToString();
    }

    private IEnumerable<string> SearchPaths(string path, HashSet<string> visited)
    {
        var files = Parent().Sandbox();
        if (!files.IsDirectory(path))
        {
            yield return path;
            yield break;
        }
        // Canonical identity prevents inward links from recursively revisiting
        // a directory. This is only a cycle check: each enumeration and open
        // still enforces the rooted boundary independently.
        if (!visited.Add(files.Resolve(path)))
        {
            yield break;
        }
        foreach (string name in files.ReadDirectory(path))
        {
            string entry = Path.Combine(path, name);
            bool directory;
            try
            {
                directory = files.IsDirectory(entry);
            }
            catch (Exception error) when (error is IOException or UnauthorizedAccessException)
            {
                // A recursive scan may encounter an outward link or a vanished
                // entry. Skip it without opening outside or losing other matches.
                // The explicitly requested root was checked above and still
                // refuses; this catch applies only to discovered children.
                continue;
            }
            if (directory)
            {
                // Preserve full-access search's existing no-directory-links
                // behavior. Host metadata is consulted only in full access;
                // confined traversal never uses path-based host operations.
                if (files.Root().Length == 0 && (File.GetAttributes(entry) & FileAttributes.ReparsePoint) != 0)
                {
                    continue;
                }
                foreach (string file in SearchPaths(entry, visited))
                {
                    yield return file;
                }
            }
            else
            {
                yield return entry;
            }
        }
    }

    private string SearchFiles(JsonElement args)
    {
        var radius = Integer(args, "context_lines", 0);
        if (radius < 0)
        {
            throw new ArgumentException("context_lines must be nonnegative.");
        }
        var pattern = new Regex(String(args, "pattern"));
        var result = new StringBuilder();
        var visited = new HashSet<string>(StringComparer.Ordinal);
        string path = String(args, "path");
        bool recursive = Parent().Sandbox().IsDirectory(path);
        foreach (var file in SearchPaths(path, visited))
        {
            var lineList = new List<string>();
            try
            {
                using var input = Parent().Sandbox().OpenFile(file, FileMode.Open, FileAccess.Read);
                using var reader = new StreamReader(input);
                while (reader.ReadLine() is { } line)
                {
                    lineList.Add(line);
                }
            }
            catch (Exception error) when (recursive && (error is IOException or UnauthorizedAccessException))
            {
                // Recursive discovery is best effort for inaccessible and
                // non-file entries; a direct file request remains a refusal.
                continue;
            }
            var lines = lineList.ToArray();
            var matches = lines.Select(line => pattern.IsMatch(line)).ToArray();
            var windows = new List<(int Start, int End)>();
            for (int i = 0; i < lines.Length; i++)
            {
                if (!matches[i])
                {
                    continue;
                }
                // Use long arithmetic so a legal large radius cannot overflow.
                var start = (int)Math.Max(0L, (long)i - radius);
                var end = (int)Math.Min(lines.Length - 1L, (long)i + radius);
                if (windows.Count > 0 && start <= windows[^1].End + 1)
                {
                    windows[^1] = (windows[^1].Start, end);
                }
                else
                {
                    windows.Add((start, end));
                }
            }
            for (int w = 0; w < windows.Count; w++)
            {
                if (w > 0 && radius > 0)
                {
                    result.Append("--\n");
                }
                for (int i = windows[w].Start; i <= windows[w].End; i++)
                {
                    var separator = matches[i] ? ':' : '-';
                    result.Append(file).Append(separator).Append(i + 1)
                        .Append(separator).Append(lines[i]).Append('\n');
                }
            }
        }
        return result.ToString();
    }

    private string WriteFile(JsonElement args)
    {
        var path = String(args, "path");
        var append = Boolean(args, "append");
        var overwrite = Boolean(args, "overwrite");
        var files = Parent().Sandbox();
        var directory = Path.GetDirectoryName(path);
        files.CreateDirectory(string.IsNullOrEmpty(directory) ? "." : directory);
        // CreateNew enforces no-overwrite at the actual open. A separate
        // existence/size probe would race, and must not reopen a host path.
        var mode = append ? FileMode.Append : overwrite ? FileMode.Create : FileMode.CreateNew;
        using var stream = files.OpenFile(path, mode, FileAccess.Write);
        using var writer = new StreamWriter(stream, new UTF8Encoding(false));
        writer.Write(String(args, "content"));
        return "File written.";
    }

    private string EditFile(JsonElement args)
    {
        var path = String(args, "path");
        var oldText = String(args, "old_text");
        var replacement = String(args, "new_text");
        if (oldText.Length == 0)
        {
            throw new ArgumentException("old_text is empty for " + path + "; read the file and choose a nonempty unique anchor.");
        }
        // Read and replace through one handle. A swapped directory entry must
        // not redirect the write after the anchor was checked in another file.
        using var stream = Parent().Sandbox().OpenFile(path, FileMode.Open, FileAccess.ReadWrite);
        string text;
        using (var reader = new StreamReader(stream, leaveOpen: true))
        {
            text = reader.ReadToEnd();
        }
        var first = text.IndexOf(oldText, StringComparison.Ordinal);
        if (first < 0)
        {
            throw new ArgumentException("old_text has no match in " + path + "; read the file and choose an exact anchor.");
        }
        // Advance by one, not by anchor length: "aa" occurs twice in "aaa".
        if (text.IndexOf(oldText, first + 1, StringComparison.Ordinal) >= 0)
        {
            throw new ArgumentException("old_text has multiple matches in " + path + "; read the file and narrow the anchor.");
        }
        var edited = text[..first] + replacement + text[(first + oldText.Length)..];
        stream.Position = 0;
        stream.SetLength(0);
        using var writer = new StreamWriter(stream, new UTF8Encoding(false));
        writer.Write(edited);
        return "File edited.";
    }

}
