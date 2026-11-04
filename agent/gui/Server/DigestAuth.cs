using System.Buffers.Binary;
using System.Globalization;
using System.Net;
using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace EnsembleGui;

public static partial class Gui
{
    // Owner-only verifier, read before the listener accepts any request.
    public static partial GuiDigestCredentials LoadDigestCredentials(string path)
    {
        Info(path);
        byte[] data = ReadBounded(path);
        string text;
        try
        {
            text = new UTF8Encoding(false, true).GetString(data);
        }
        catch (ArgumentException)
        {
            throw new InvalidDataException("GUI auth file is not valid UTF-8");
        }
        var invalid = new InvalidDataException("invalid GUI auth credential document");
        int version = 0;
        string username = "", realm = "", ha1 = "";
        var seen = new HashSet<string>(StringComparer.Ordinal);
        try
        {
            using var document = JsonDocument.Parse(text,
                new JsonDocumentOptions { AllowTrailingCommas = false, CommentHandling = JsonCommentHandling.Disallow });
            if (document.RootElement.ValueKind != JsonValueKind.Object)
            {
                throw invalid;
            }
            foreach (var property in document.RootElement.EnumerateObject())
            {
                if (!seen.Add(property.Name))
                {
                    throw invalid; // A duplicate key would silently shadow a verifier field.
                }
                switch (property.Name)
                {
                    case "version":
                        if (property.Value.ValueKind != JsonValueKind.Number
                            || !property.Value.TryGetInt32(out version))
                        {
                            throw invalid;
                        }
                        break;
                    case "username":
                        username = ReadString(property.Value) ?? throw invalid;
                        break;
                    case "realm":
                        realm = ReadString(property.Value) ?? throw invalid;
                        break;
                    case "ha1":
                        ha1 = ReadString(property.Value) ?? throw invalid;
                        break;
                    default:
                        throw invalid; // Unknown fields are refused, so a typo cannot disable a check.
                }
            }
        }
        catch (JsonException)
        {
            throw invalid;
        }
        if (seen.Count != 4 || version != GuiDigestCredentials.CurrentVersion
            || realm != GuiDigestCredentials.RealmName || !ValidUsername(username))
        {
            throw invalid;
        }
        byte[] digest;
        try
        {
            digest = Convert.FromHexString(ha1);
        }
        catch (FormatException)
        {
            throw invalid;
        }
        if (digest.Length != 32 || ha1 != ha1.ToLowerInvariant())
        {
            throw invalid;
        }
        return new GuiDigestCredentials(version, username, realm, ha1);
    }

    public static partial string CreateDigestCredentials(string path, string username,
        bool replace)
    {
        if (!ValidUsername(username))
        {
            throw new ArgumentException("GUI auth username must be 1-256 ASCII letters, digits, dots, underscores or hyphens");
        }
        DateTime beforeWrite = default;
        long beforeLength = 0;
        if (replace)
        {
            var info = Info(path);
            LoadDigestCredentials(path); // Only a usable verifier may be replaced.
            beforeWrite = info.LastWriteTimeUtc;
            beforeLength = info.Length;
        }
        else if (File.Exists(path))
        {
            throw new IOException("GUI auth file already exists");
        }
        byte[] secret = RandomNumberGenerator.GetBytes(16);
        string password = Base64Url(secret);
        string ha1 = Convert.ToHexStringLower(SHA256.HashData(Encoding.UTF8.GetBytes(
            username + ":" + GuiDigestCredentials.RealmName + ":" + password)));
        var credentials = new GuiDigestCredentials(GuiDigestCredentials.CurrentVersion,
            username, GuiDigestCredentials.RealmName, ha1);
        string directory = Path.GetDirectoryName(Path.GetFullPath(path)) ?? ".";
        if (OperatingSystem.IsWindows())
        {
            Directory.CreateDirectory(directory);
        }
        else
        {
            Directory.CreateDirectory(directory, UnixFileMode.UserRead | UnixFileMode.UserWrite
                | UnixFileMode.UserExecute);
        }
        string temporary = Path.Combine(directory, $".gui-auth-{Guid.NewGuid():N}.tmp");
        string payload = JsonSerializer.Serialize(new
        {
            version = credentials.Version,
            username = credentials.Username,
            realm = credentials.Realm,
            ha1 = credentials.Ha1
        }) + "\n";
        try
        {
            using (var file = new FileStream(temporary, FileMode.CreateNew, FileAccess.Write,
                       FileShare.None))
            {
                if (!OperatingSystem.IsWindows())
                {
                    File.SetUnixFileMode(temporary, UnixFileMode.UserRead | UnixFileMode.UserWrite);
                }
                file.Write(Encoding.UTF8.GetBytes(payload));
                file.Flush(true);
            }
            if (replace)
            {
                var current = new FileInfo(path);
                if (!current.Exists || current.LastWriteTimeUtc != beforeWrite
                    || current.Length != beforeLength)
                {
                    throw new IOException("GUI auth file changed during reset");
                }
                File.Move(temporary, path, true);
            }
            else
            {
                File.Move(temporary, path, false); // Refuses to clobber a concurrent initializer.
            }
        }
        finally
        {
            if (File.Exists(temporary))
            {
                File.Delete(temporary);
            }
        }
        return password;
    }

    public static partial string? DigestPath(string? dataDirectory, string? requested,
        bool requestedExplicitly, bool loopback)
    {
        if (!requestedExplicitly && string.IsNullOrWhiteSpace(dataDirectory))
        {
            throw new ArgumentException("A conventional GUI auth location needs a data directory");
        }
        string path = requestedExplicitly
            ? (string.IsNullOrWhiteSpace(requested)
                ? throw new ArgumentException("--gui-auth-file requires a nonempty path")
                : Path.GetFullPath(requested!))
            : Path.Combine(dataDirectory!, "gui-auth.json");
        if (File.Exists(path))
        {
            return path;
        }
        if (!requestedExplicitly && loopback)
        {
            return null; // Local use stays open; presence of the file is the opt-in.
        }
        throw new GuiBindRefusalException(
            "GUI authentication required; run auth-init --file \"" + path + "\" before starting this listener");
    }

    public static partial void UseDigestAuth(WebApplication host, string path)
    {
        var credentials = LoadDigestCredentials(path);
        var guard = new DigestGuard(credentials, RandomNumberGenerator.GetBytes(32));
        host.Use(async (context, next) =>
        {
            if (await guard.AdmitAsync(context))
            {
                await next();
            }
        });
    }

    internal static partial IPAddress ResolveBindAddress(string host)
    {
        string text = (host ?? "").Trim();
        if (text.Length == 0)
        {
            throw new ArgumentException("GUI host must be an address, not empty");
        }
        if (text.Equals("localhost", StringComparison.OrdinalIgnoreCase))
        {
            return IPAddress.Loopback;
        }
        if (text.IndexOfAny([':', '/', '%']) < 0 && text.Any(character =>
                character is >= 'a' and <= 'z' or >= 'A' and <= 'Z'))
        {
            // A hostname would make the listener's reach depend on a resolver and
            // its hosts file, so only literal addresses may widen it.
            throw new ArgumentException("GUI host must be an IP address or localhost, not a name");
        }
        if (!IPAddress.TryParse(text, out IPAddress? address) || address.AddressFamily
            is not System.Net.Sockets.AddressFamily.InterNetwork
            and not System.Net.Sockets.AddressFamily.InterNetworkV6)
        {
            throw new ArgumentException("GUI host must be a valid IPv4 or IPv6 address");
        }
        return address;
    }

    public static partial bool IsLoopbackHost(string host) =>
        IPAddress.IsLoopback(ResolveBindAddress(host));

    public static partial int RunAuthCommand(string[] args, TextWriter output)
    {
        bool replace = args[0] == "auth-reset";
        string? file = null;
        string username = "ensemble";
        string? dataDirectory = null;
        for (int index = 1; index < args.Length; index++)
        {
            string flag = args[index];
            string? Value()
            {
                if (index + 1 >= args.Length)
                {
                    throw new ArgumentException(flag + " requires a value");
                }
                return args[++index];
            }
            switch (flag)
            {
                case "--file":
                    file = Value();
                    break;
                case "--username":
                    username = Value()!;
                    break;
                case "--data-dir":
                    dataDirectory = Value();
                    break;
                default:
                    throw new ArgumentException(
                        "usage: auth-init|auth-reset [--file path] [--username ensemble] [--data-dir path]");
            }
        }
        string path = string.IsNullOrWhiteSpace(file)
            ? Path.Combine(Path.GetFullPath(dataDirectory ?? Directory.GetCurrentDirectory()),
                "gui-auth.json")
            : Path.GetFullPath(file!);
        string password = CreateDigestCredentials(path, username, replace);
        output.WriteLine(password);
        return 0;
    }

    public static partial void CheckDigestPath(string path, IEnumerable<string> reserved)
    {
        string target = Final(path);
        foreach (string other in reserved.Where(candidate => !string.IsNullOrWhiteSpace(candidate)))
        {
            if (string.Equals(Final(other), target, StringComparison.Ordinal))
            {
                throw new InvalidOperationException(
                    "GUI auth file must not use an application log, configuration or persistence path");
            }
        }
    }

    internal static partial Dictionary<string, string>? ParseDigestFields(string? header)
    {
        if (header is null || header.Length > 8192)
        {
            return null;
        }
        int space = header.IndexOf(' ');
        if (space < 0 || !string.Equals(header[..space].Trim(), "Digest",
                StringComparison.OrdinalIgnoreCase))
        {
            return null;
        }
        var fields = new Dictionary<string, string>(StringComparer.Ordinal);
        string rest = header[(space + 1)..].Trim();
        while (rest.Length != 0)
        {
            // Separators are followed by optional whitespace, so each field starts
            // by skipping it; otherwise the second field would parse as a name of
            // zero length and the whole header would be refused.
            rest = rest.TrimStart();
            string name = TakeToken(rest, out rest);
            name = name.ToLowerInvariant();
            if (name.Length == 0)
            {
                return null;
            }
            rest = rest.TrimStart();
            if (!rest.StartsWith("="))
            {
                return null;
            }
            rest = rest[1..].TrimStart();
            string value;
            if (rest.StartsWith("\""))
            {
                rest = rest[1..];
                var decoded = new StringBuilder();
                bool closed = false;
                while (rest.Length != 0)
                {
                    char character = rest[0];
                    rest = rest[1..];
                    if (character == '"')
                    {
                        closed = true;
                        break;
                    }
                    if (character == '\\')
                    {
                        if (rest.Length == 0)
                        {
                            return null;
                        }
                        character = rest[0];
                        rest = rest[1..];
                    }
                    if (character < ' ' || character == (char)127)
                    {
                        return null;
                    }
                    decoded.Append(character);
                }
                if (!closed)
                {
                    return null;
                }
                value = decoded.ToString();
            }
            else
            {
                value = TakeToken(rest, out rest);
                if (value.Length == 0)
                {
                    return null;
                }
            }
            if (fields.ContainsKey(name))
            {
                return null;
            }
            fields[name] = value;
            rest = rest.TrimStart();
            if (rest.Length == 0)
            {
                continue;
            }
            if (rest[0] != ',')
            {
                return null;
            }
            rest = rest[1..];
            if (rest.Trim().Length == 0)
            {
                return null;
            }
        }
        return fields.Count > 0 ? fields : null;
    }

    internal static partial bool DigestResponseMatches(GuiDigestCredentials credentials,
        string method, string uri, Dictionary<string, string> fields)
    {
        if (fields.ContainsKey("username*") || !fields.TryGetValue("nc", out string? nc)
            || !fields.TryGetValue("cnonce", out string? cnonce)
            || !fields.TryGetValue("nonce", out string? nonce)
            || !fields.TryGetValue("qop", out string? qop)
            || !fields.TryGetValue("algorithm", out string? algorithm)
            || !fields.TryGetValue("realm", out string? realm)
            || !fields.TryGetValue("uri", out string? requested)
            || !fields.TryGetValue("username", out string? username))
        {
            return false;
        }
        // The counter is hex, exactly eight digits, and starts at one: nc=0 would
        // let a client reuse the zero key in the replay table indefinitely.
        if (nc.Length != 8 || !IsHex(nc)
            || uint.Parse(nc, NumberStyles.HexNumber, CultureInfo.InvariantCulture) == 0
            || qop != "auth"
            || !string.Equals(algorithm, "SHA-256", StringComparison.OrdinalIgnoreCase)
            || realm != credentials.Realm || requested != uri || username != credentials.Username
            || cnonce.Length == 0 || cnonce.Length > 256 || nonce.Length == 0
            || (fields.TryGetValue("userhash", out string? userhash) && userhash != "false"))
        {
            return false;
        }
        byte[] response;
        try
        {
            response = Convert.FromHexString(fields.GetValueOrDefault("response", ""));
        }
        catch (FormatException)
        {
            return false;
        }
        if (response.Length != 32)
        {
            return false;
        }
        string requestHash = Convert.ToHexStringLower(SHA256.HashData(
            Encoding.UTF8.GetBytes(method + ":" + requested)));
        byte[] expected = SHA256.HashData(Encoding.UTF8.GetBytes(credentials.Ha1 + ":" + nonce
            + ":" + nc + ":" + cnonce + ":auth:" + requestHash));
        return CryptographicOperations.FixedTimeEquals(expected, response);
    }

    internal static partial string DigestNonce(byte[] key, DateTimeOffset issued)
    {
        byte[] data = new byte[56];
        BinaryPrimitives.WriteInt64BigEndian(data.AsSpan(0, 8), issued.ToUnixTimeSeconds());
        RandomNumberGenerator.Fill(data.AsSpan(8, 16));
        HMACSHA256.HashData(key, data.AsSpan(0, 24)).CopyTo(data.AsSpan(24, 32));
        return Base64Url(data);
    }

    internal static partial DateTimeOffset? DigestNonceExpiry(byte[] key, string nonce,
        DateTimeOffset now)
    {
        byte[] data;
        try
        {
            data = FromBase64Url(nonce);
        }
        catch (FormatException)
        {
            return null;
        }
        if (data.Length != 56
            || !CryptographicOperations.FixedTimeEquals(HMACSHA256.HashData(key,
                data.AsSpan(0, 24)), data.AsSpan(24, 32)))
        {
            return null;
        }
        DateTimeOffset issued = DateTimeOffset.FromUnixTimeSeconds(
            BinaryPrimitives.ReadInt64BigEndian(data.AsSpan(0, 8)));
        if (issued > now.AddSeconds(1))
        {
            return null; // A future-dated nonce cannot have come from this server.
        }
        return issued.AddSeconds(DigestNonceLifetimeSeconds);
    }

    // Internal so the regression suite can drive the request decision directly:
    // the listener's nonce key is deliberately unreachable from outside.
    internal sealed class DigestGuard
    {
        private readonly GuiDigestCredentials credentials;
        private readonly byte[] key;
        private readonly object gate = new();
        private readonly Dictionary<string, DateTimeOffset> replays = [];
        private readonly Dictionary<string, (int Count, DateTimeOffset Until, DateTimeOffset Updated)> failures = [];

        internal DigestGuard(GuiDigestCredentials credentials, byte[] key)
        {
            this.credentials = credentials;
            this.key = key;
        }

        // True admits the request; false has already answered it.
        // True admits the request; false has already answered it.
        internal async ValueTask<bool> AdmitAsync(HttpContext context)
        {
            DateTimeOffset now = DateTimeOffset.UtcNow;
            string source = context.Connection.RemoteIpAddress?.ToString() ?? "unknown";
            bool blocked = false;
            int retry = 0;
            lock (gate)
            {
                if (failures.TryGetValue(source, out var recorded) && now < recorded.Until)
                {
                    blocked = true;
                    retry = Math.Max(1, (int)(recorded.Until - now).TotalSeconds + 1);
                }
            }
            if (blocked)
            {
                await DenyAsync(context, 429, "Too many authentication attempts", retry);
                return false;
            }
            var schemes = context.Request.Headers["Authorization"];
            if (schemes.Count == 0)
            {
                await ChallengeAsync(context, now, false);
                return false;
            }
            if (schemes.Count != 1)
            {
                Failed(source, now);
                await ChallengeAsync(context, now, false);
                return false;
            }
            // The Digest uri field must match the request target exactly, so the
            // same spelling is used for both comparison and the hash.
            string uri = (context.Request.Path.Value ?? "")
                + (context.Request.QueryString.HasValue ? context.Request.QueryString.Value! : "");
            var fields = ParseDigestFields(schemes[0]);
            if (fields is null
                || !DigestResponseMatches(credentials, context.Request.Method, uri, fields))
            {
                Failed(source, now);
                await ChallengeAsync(context, now, false);
                return false;
            }
            DateTimeOffset? expiry = DigestNonceExpiry(key, fields["nonce"], now);
            if (expiry is null)
            {
                Failed(source, now);
                await ChallengeAsync(context, now, false);
                return false;
            }
            if (now >= expiry.Value)
            {
                // A stale nonce is the browser reusing an old one, not a guess: it
                // costs no failure credit, only a fresh challenge.
                await ChallengeAsync(context, now, true);
                return false;
            }
            string replayKey = fields["nonce"] + ":" + fields["cnonce"] + ":"
                + uint.Parse(fields["nc"], NumberStyles.HexNumber, CultureInfo.InvariantCulture)
                    .ToString("x", CultureInfo.InvariantCulture);
            bool replay = false;
            bool busy = false;
            lock (gate)
            {
                foreach (string previous in replays.Where(entry => now >= entry.Value)
                             .Select(entry => entry.Key).ToArray())
                {
                    replays.Remove(previous);
                }
                replay = replays.ContainsKey(replayKey);
                if (!replay && replays.Count < DigestReplayLimit)
                {
                    replays[replayKey] = expiry.Value;
                    failures.Remove(source); // One success clears this address's slowdown.
                }
                else if (!replay)
                {
                    busy = true;
                }
            }
            if (replay)
            {
                Failed(source, now);
                await ChallengeAsync(context, now, false);
                return false;
            }
            if (busy)
            {
                await DenyAsync(context, 503, "Authentication temporarily busy", 1);
                return false;
            }
            return true;
        }

        private async ValueTask ChallengeAsync(HttpContext context, DateTimeOffset now,
            bool stale)
        {
            string nonce = DigestNonce(key, now);
            context.Response.Headers.WWWAuthenticate = "Digest realm=\"" + credentials.Realm
                + "\", algorithm=SHA-256, qop=\"auth\", charset=UTF-8, nonce=\"" + nonce + "\""
                + (stale ? ", stale=true" : "");
            await DenyAsync(context, 401, "Unauthorized", 0);
        }

        // Repeated failures from one address back off exponentially, capped, so a
        // host on the LAN cannot grind through guesses at line speed.
        private void Failed(string source, DateTimeOffset now)
        {
            lock (gate)
            {
                foreach (string address in failures.Where(entry =>
                             now - entry.Value.Updated >= TimeSpan.FromMinutes(10))
                             .Select(entry => entry.Key).ToArray())
                {
                    failures.Remove(address);
                }
                int count = failures.TryGetValue(source, out var previous)
                    ? Math.Min(previous.Count + 1, 8) : 1;
                if (count == 1 && failures.Count >= DigestFailureLimit)
                {
                    string oldest = failures.OrderBy(entry => entry.Value.Updated)
                        .Select(entry => entry.Key).First();
                    failures.Remove(oldest);
                }
                DateTimeOffset until = count >= 3
                    ? now.AddSeconds(1L << Math.Min(count - 3, 5)) : DateTimeOffset.MinValue;
                failures[source] = (count, until, now);
            }
        }

        private static async Task DenyAsync(HttpContext context, int status, string message,
            int retryAfter)
        {
            if (retryAfter > 0)
            {
                context.Response.Headers.RetryAfter = retryAfter.ToString(CultureInfo.InvariantCulture);
            }
            context.Response.Headers.CacheControl = "no-store";
            context.Response.StatusCode = status;
            await context.Response.Body.WriteAsync(Encoding.UTF8.GetBytes(message + "\n"));
        }
    }

    private static FileInfo Info(string path)
    {
        var info = new FileInfo(path);
        if (!info.Exists || (info.Attributes & FileAttributes.ReparsePoint) != 0
            || !info.Attributes.HasFlag(FileAttributes.Normal))
        {
            throw new InvalidDataException("GUI auth file must be a regular file, without symlinks");
        }
        // Ownership and hard-link count are checked by the owning user's umask and
        // the 0600 requirement below; .NET exposes no portable uid read on Linux.
        // Windows has no Unix mode bits and an NTFS ACL is outside this file's
        // scope, so the owner-only requirement is enforced where the concept exists.
        if (!OperatingSystem.IsWindows() && File.GetUnixFileMode(path)
            != (UnixFileMode.UserRead | UnixFileMode.UserWrite))
        {
            throw new InvalidDataException("GUI auth file must have mode 0600");
        }
        if (info.Length > DigestCredentialsLimit)
        {
            throw new InvalidDataException("GUI auth file is too large");
        }
        return info;
    }

    private static byte[] ReadBounded(string path)
    {
        FileInfo before = Info(path);
        using var file = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.Read);
        Info(path); // Re-check after opening: the path may have been swapped.
        if (file.Length != before.Length || file.Length > DigestCredentialsLimit)
        {
            throw new InvalidDataException("GUI auth file changed while opening");
        }
        byte[] data = new byte[file.Length];
        file.ReadExactly(data);
        return data;
    }

    private static string? ReadString(JsonElement element) =>
        element.ValueKind == JsonValueKind.String ? element.GetString() : null;

    private static bool ValidUsername(string username) =>
        username.Length is > 0 and <= 256 && username.All(character =>
            char.IsAsciiLetterOrDigit(character) || character is '.' or '_' or '-');

    private static bool IsHex(string value) => value.All(character =>
        character is (>= '0' and <= '9') or (>= 'a' and <= 'f') or (>= 'A' and <= 'F'));

    private static string TakeToken(string text, out string rest)
    {
        int index = 0;
        while (index < text.Length && (char.IsAsciiLetterOrDigit(text[index])
               || "!#$%&'*+-.^_`|~".Contains(text[index])))
        {
            index++;
        }
        rest = text[index..];
        return text[..index];
    }

    private static string Base64Url(byte[] data) => Convert.ToBase64String(data)
        .TrimEnd('=').Replace('+', '-').Replace('/', '_');

    private static byte[] FromBase64Url(string value) => Convert.FromBase64String(
        value.Replace('-', '+').Replace('_', '/')
        .PadRight(value.Length + (4 - value.Length % 4) % 4, '='));

    private static string Final(string path)
    {
        string full = Path.GetFullPath(path);
        try
        {
            return File.ResolveLinkTarget(full, returnFinalTarget: true) is { } target
                ? target.FullName : full;
        }
        catch (IOException)
        {
            return full; // Nothing to resolve for a path the caller has not created yet.
        }
    }
}
