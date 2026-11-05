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

    public static partial void UseDigestAuth(WebApplication host, string path,
        TextWriter? diagnostics)
    {
        var credentials = LoadDigestCredentials(path);
        var guard = new DigestGuard(credentials, RandomNumberGenerator.GetBytes(32), path,
            diagnostics, watchVerifier: true);
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

    public static partial int RunAuthCommand(string[] args, TextWriter output,
        TextWriter? diagnostics)
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
        // The listener reads <data directory>/gui-auth.json, so the subcommand looks in
        // the same place by default instead of the current directory. The last resort is
        // still the current directory, which is correct for an embedded host and wrong
        // for a launcher started with --data-dir: a password printed beside the file
        // nobody reads simply fails to work. Both facts are now stated, because the
        // silent version cost two debugging sessions.
        string environment = Environment.GetEnvironmentVariable("EN_DATA_DIR") ?? "";
        string resolved;
        string directory;
        if (!string.IsNullOrWhiteSpace(file))
        {
            resolved = "--file";
            directory = Path.GetDirectoryName(Path.GetFullPath(file!))
                ?? Directory.GetCurrentDirectory();
        }
        else if (!string.IsNullOrWhiteSpace(dataDirectory))
        {
            resolved = "--data-dir";
            directory = Path.GetFullPath(dataDirectory!);
        }
        else if (!string.IsNullOrWhiteSpace(environment))
        {
            resolved = "EN_DATA_DIR";
            directory = Path.GetFullPath(environment);
        }
        else
        {
            resolved = "the current directory";
            directory = Path.GetFullPath(Directory.GetCurrentDirectory());
        }
        string path = Path.Combine(directory, "gui-auth.json");
        string password = CreateDigestCredentials(path, username, replace);
        output.WriteLine(password);
        if (diagnostics is not null)
        {
            diagnostics.WriteLine("auth: verifier written to " + path + " (resolved from "
                + resolved + ")");
            if (resolved == "the current directory")
            {
                diagnostics.WriteLine("auth: no --data-dir or EN_DATA_DIR was given, so a "
                    + "listener started with one of those reads a different file; pass "
                    + "--data-dir to match it, or run the launcher's auth subcommand");
            }
            string stray = Path.Combine(Directory.GetCurrentDirectory(), "gui-auth.json");
            if (!string.Equals(Path.GetFullPath(stray), path, StringComparison.Ordinal)
                && File.Exists(stray))
            {
                diagnostics.WriteLine("auth: warning \"" + stray + "\" also exists and is not "
                    + "the file this listener reads; delete it to avoid an unused password");
            }
        }
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
        string method, string uri, Dictionary<string, string> fields) =>
        DigestRejection(credentials, method, uri, fields) is null;

    // Client-chosen text reaches the operator's log, so it is stripped of control
    // characters and shortened before quoting. A log line must not become a channel
    // for forged terminal output or an unbounded attacker string.
    private static string Quoted(string value)
    {
        string visible = new(value.Where(character => !char.IsControl(character)).ToArray());
        return "\"" + (visible.Length > 48 ? visible[..48] + "..." : visible) + "\"";
    }

    // Every rule the Digest answer must satisfy, in one place, reporting the first it
    // fails. The constant-time comparison stays last so no earlier check depends on it.
    // Naming the failing rule in the operator's log distinguishes a mistyped password
    // from a stale realm or a browser reusing a nonce, which a bare bool made
    // undiagnosable. The client still receives one uniform 401 for every reason, so
    // nothing here helps an attacker tell a good username from a bad one.
    internal static partial string? DigestRejection(GuiDigestCredentials credentials,
        string method, string uri, Dictionary<string, string> fields)
    {
        if (fields.ContainsKey("username*"))
        {
            return "the username* extension is not accepted";
        }
        foreach (string field in new[] { "nc", "cnonce", "nonce", "qop", "algorithm",
            "realm", "uri", "username", "response" })
        {
            if (!fields.ContainsKey(field))
            {
                return "the Authorization header omits " + field;
            }
        }
        string nc = fields["nc"], cnonce = fields["cnonce"], nonce = fields["nonce"],
            qop = fields["qop"], algorithm = fields["algorithm"], realm = fields["realm"],
            requested = fields["uri"], username = fields["username"];
        // The counter is hex, exactly eight digits, and starts at one: nc=0 would
        // let a client reuse the zero key in the replay table indefinitely.
        if (nc.Length != 8 || !IsHex(nc)
            || uint.Parse(nc, NumberStyles.HexNumber, CultureInfo.InvariantCulture) == 0)
        {
            return "the replay counter is not eight hex digits above zero";
        }
        if (qop != "auth")
        {
            return "qop must be auth, not " + Quoted(qop);
        }
        if (!string.Equals(algorithm, "SHA-256", StringComparison.OrdinalIgnoreCase))
        {
            return "algorithm must be SHA-256, not " + Quoted(algorithm);
        }
        if (realm != credentials.Realm)
        {
            return "realm " + Quoted(realm) + " does not match the verifier's "
                + Quoted(credentials.Realm);
        }
        if (requested != uri)
        {
            return "uri " + Quoted(requested) + " does not match the request target "
                + Quoted(uri);
        }
        if (username != credentials.Username)
        {
            return "username " + Quoted(username) + " is not the configured user "
                + Quoted(credentials.Username);
        }
        if (cnonce.Length == 0 || cnonce.Length > 256)
        {
            return "cnonce is empty or longer than 256 characters";
        }
        if (nonce.Length == 0)
        {
            return "nonce is empty";
        }
        if (fields.TryGetValue("userhash", out string? userhash) && userhash != "false")
        {
            return "userhash must be absent or false";
        }
        byte[] response;
        try
        {
            response = Convert.FromHexString(fields["response"]);
        }
        catch (FormatException)
        {
            return "response is not hexadecimal";
        }
        if (response.Length != 32)
        {
            return "response is not 32 bytes";
        }
        string requestHash = Convert.ToHexStringLower(SHA256.HashData(
            Encoding.UTF8.GetBytes(method + ":" + requested)));
        byte[] expected = SHA256.HashData(Encoding.UTF8.GetBytes(credentials.Ha1 + ":" + nonce
            + ":" + nc + ":" + cnonce + ":auth:" + requestHash));
        return CryptographicOperations.FixedTimeEquals(expected, response)
            ? null
            : "the password does not match the verifier";
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
        private GuiDigestCredentials credentials;
        private byte[] key;
        private byte[] verifierBytes;
        private readonly bool watched;
        private DateTimeOffset verifierSampled;
        private string? verifierFault;
        private readonly object gate = new();
        private readonly Dictionary<string, DateTimeOffset> replays = [];
        private readonly Dictionary<string, (int Count, DateTimeOffset Until, DateTimeOffset Updated)> failures = [];
        private readonly string? verifierPath;
        private readonly TextWriter? diagnostics;
        private bool verifierAnnounced;

        // watchVerifier is set only when the file is the guard's authority: UseDigestAuth
        // loaded the credentials from it, so a later reset must take effect. An embedding
        // host that supplies its own credentials names a path purely so refusals can say
        // where they came from, and that file must not be able to revoke them.
        internal DigestGuard(GuiDigestCredentials credentials, byte[] key,
            string? verifierPath = null, TextWriter? diagnostics = null,
            bool watchVerifier = false)
        {
            this.credentials = credentials;
            this.key = key;
            this.verifierPath = verifierPath;
            this.diagnostics = diagnostics;
            this.verifierBytes = Array.Empty<byte>();
            this.watched = verifierPath is not null && watchVerifier;
            if (this.watched)
            {
                // The startup load already succeeded, so a sample taken here is the
                // baseline a later reset has to differ from.
                try
                {
                    // watched implies the path, checked where the field is assigned.
                    this.verifierBytes = File.ReadAllBytes(verifierPath!);
                }
                catch (Exception)
                {
                    // A file that vanishes between load and sample is reported by the
                    // next refresh, which fails the request closed.
                }
            }
        }

        // One line per refused attempt, naming the rule that refused it. The verifier
        // path is stated on the first refusal because a password generated for one
        // file silently cannot unlock a listener reading another, and that mismatch is
        // otherwise invisible: both files look correct on disk.
        private void Explain(string source, string reason)
        {
            if (diagnostics is not null)
            {
                bool announce;
                lock (gate)
                {
                    announce = !verifierAnnounced;
                    verifierAnnounced = true;
                }
                // Written outside the gate: it also guards the replay and slowdown
                // tables, and no request should wait on somebody's terminal.
                diagnostics.WriteLine(announce
                    ? $"auth: rejected {reason} from {source}; verifier "
                        + (verifierPath is null ? "is held by the embedding host"
                            : "file " + verifierPath)
                    : $"auth: rejected {reason} from {source}");
            }
        }

        // A reset while this listener is running is otherwise refused forever: the
        // verifier is loaded once at startup, so a freshly generated password is
        // rejected as wrong by a file it was never written to, and the operator is
        // told their password does not match. The file is 128 bytes, so its contents
        // are sampled at most once a second and reloaded when they differ. Reloading
        // rotates the nonce key, which is what forces a browser holding the previous
        // password to take a fresh challenge instead of silently replaying the old
        // credential; the refusal Bill saw twice came from that replay. A verifier that
        // becomes unreadable fails closed, because replacing the file is also how a
        // credential is revoked, and a half-written file must not keep the old access.
        private void RefreshVerifier(DateTimeOffset now)
        {
            if (!watched)
            {
                return;
            }
            lock (gate)
            {
                if (verifierFault is null && now - verifierSampled < TimeSpan.FromSeconds(1))
                {
                    return;
                }
                verifierSampled = now;
            }
            string watchedPath = verifierPath!;
            byte[] sample;
            try
            {
                sample = File.ReadAllBytes(watchedPath);
            }
            catch (Exception failure) when (failure is IOException
                or UnauthorizedAccessException or NotSupportedException)
            {
                lock (gate)
                {
                    verifierFault = "could not be read: " + failure.Message;
                }
                return;
            }
            bool changed;
            lock (gate)
            {
                changed = verifierFault is not null
                    || !CryptographicOperations.FixedTimeEquals(verifierBytes, sample);
            }
            if (!changed)
            {
                return;
            }
            GuiDigestCredentials reloaded;
            try
            {
                reloaded = Gui.LoadDigestCredentials(watchedPath);
            }
            catch (Exception failure) when (failure is IOException or InvalidDataException
                or UnauthorizedAccessException or InvalidOperationException
                or NotSupportedException or FormatException)
            {
                lock (gate)
                {
                    verifierFault = "could not be loaded: " + failure.Message;
                }
                return;
            }
            lock (gate)
            {
                credentials = reloaded;
                key = RandomNumberGenerator.GetBytes(32);
                replays.Clear();
                verifierBytes = sample;
                verifierFault = null;
                verifierAnnounced = false;
            }
            diagnostics?.WriteLine("auth: verifier reloaded from " + verifierPath
                + "; clients holding the previous password will be challenged again");
        }

        private string? Fault()
        {
            lock (gate)
            {
                return verifierFault;
            }
        }

        // True admits the request; false has already answered it.
        internal async ValueTask<bool> AdmitAsync(HttpContext context)
        {
            DateTimeOffset now = DateTimeOffset.UtcNow;
            string source = context.Connection.RemoteIpAddress?.ToString() ?? "unknown";
            RefreshVerifier(now);
            string? fault = Fault();
            if (fault is not null)
            {
                Explain(source, "the verifier " + fault);
                await DenyAsync(context, 503, "GUI authentication is unavailable", 0);
                return false;
            }
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
                Explain(source, "more than one Authorization header was sent");
                Failed(source, now);
                await ChallengeAsync(context, now, false);
                return false;
            }
            // The Digest uri field must match the request target exactly, so the
            // same spelling is used for both comparison and the hash.
            string uri = (context.Request.Path.Value ?? "")
                + (context.Request.QueryString.HasValue ? context.Request.QueryString.Value! : "");
            var fields = ParseDigestFields(schemes[0]);
            if (fields is null)
            {
                Explain(source, "the Authorization header is not parseable Digest syntax");
                Failed(source, now);
                await ChallengeAsync(context, now, false);
                return false;
            }
            string? rejection = DigestRejection(credentials, context.Request.Method, uri, fields);
            if (rejection is not null)
            {
                Explain(source, rejection);
                Failed(source, now);
                await ChallengeAsync(context, now, false);
                return false;
            }
            DateTimeOffset? expiry = DigestNonceExpiry(key, fields["nonce"], now);
            if (expiry is null)
            {
                // The nonce is self-certifying, so this means the challenge did not come
                // from this listener: the usual cause is a restart between prompt and
                // retry, and a fresh challenge resolves it.
                Explain(source, "the nonce was not issued by this listener");
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
                Explain(source, "that nonce, cnonce and counter were already used");
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
