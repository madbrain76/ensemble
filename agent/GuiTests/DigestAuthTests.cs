using System.Net;
using Microsoft.AspNetCore.Http;
using System.Net.Sockets;
using System.Security.Cryptography;
using System.Text;
using Ensemble;
using EnsembleGui;

// LAN authentication: the verifier file, the Digest challenge/answer rules, the
// self-certifying nonce, the per-address slowdown, and the live middleware that
// guards assets and the WebSocket upgrade alike.
static class DigestAuthTests
{
    static void Check(bool value, string message)
    {
        if (!value)
        {
            throw new InvalidOperationException(message);
        }
    }

    static void Refuses(Action action, string message)
    {
        try
        {
            action();
        }
        catch (Exception)
        {
            return;
        }
        throw new InvalidOperationException(message);
    }

    static string Hex(byte[] data) => Convert.ToHexString(data).ToLowerInvariant();

    static string ToBase64Url(byte[] data) => Convert.ToBase64String(data).TrimEnd('=')
        .Replace('+', '-').Replace('/', '_');

    static byte[] FromBase64Url(string value) => Convert.FromBase64String(
        value.Replace('-', '+').Replace('_', '/')
        .PadRight(value.Length + (4 - value.Length % 4) % 4, '='));

    static string Answer(string username, string realm, string password, string method,
        string uri, string nonce, string nc = "00000001", string cnonce = "cnonce")
    {
        string ha1 = Hex(SHA256.HashData(Encoding.UTF8.GetBytes(username + ":" + realm + ":" + password)));
        string requestHash = Hex(SHA256.HashData(Encoding.UTF8.GetBytes(method + ":" + uri)));
        string response = Hex(SHA256.HashData(Encoding.UTF8.GetBytes(
            ha1 + ":" + nonce + ":" + nc + ":" + cnonce + ":auth:" + requestHash)));
        return "Digest realm=\"" + realm + "\", nonce=\"" + nonce + "\", uri=\"" + uri
            + "\", qop=\"auth\", algorithm=SHA-256, username=\"" + username + "\", nc=" + nc
            + ", cnonce=\"" + cnonce + "\", response=\"" + response + "\"";
    }

    static Dictionary<string, string> Fields(string username = "ensemble",
        string realm = "Ensemble", string method = "GET", string uri = "/x",
        string password = "s3cret", string nonce = "nonce", string nc = "00000001",
        string cnonce = "cc", string qop = "auth", string algorithm = "SHA-256")
    {
        string ha1 = Hex(SHA256.HashData(Encoding.UTF8.GetBytes(username + ":" + realm + ":" + password)));
        string requestHash = Hex(SHA256.HashData(Encoding.UTF8.GetBytes(method + ":" + uri)));
        return new Dictionary<string, string>(StringComparer.Ordinal)
        {
            ["realm"] = realm, ["nonce"] = nonce, ["uri"] = uri, ["qop"] = qop,
            ["algorithm"] = algorithm, ["username"] = username, ["nc"] = nc,
            ["cnonce"] = cnonce,
            ["response"] = Hex(SHA256.HashData(Encoding.UTF8.GetBytes(
                ha1 + ":" + nonce + ":" + nc + ":" + cnonce + ":" + qop + ":" + requestHash)))
        };
    }

    public static async Task RunAsync()
    {
        string work = Path.Combine(Path.GetTempPath(), "gui-digest-" + Guid.NewGuid().ToString("N"));
        Directory.CreateDirectory(work);
        try
        {
            Credentials(work);
            Parser();
            Verification();
            RejectionReasons();
            GuardDiagnostics();
            Nonce();
            Staleness();
            PathPolicy(work);
            CommandPaths(work);
            await Listener(work);
            await ReloadAndFault(work);
        }
        finally
        {
            try
            {
                Directory.Delete(work, true);
            }
            catch (IOException)
            {
            }
        }
    }

    static void Credentials(string work)
    {
        string file = Path.Combine(work, "gui-auth.json");
        string password = Gui.CreateDigestCredentials(file, "ensemble", false);
        Check(password.Length >= 20, "generated password is too short to be a secret");
        Check(!File.ReadAllText(file).Contains(password), "the verifier file stored the password");
        var loaded = Gui.LoadDigestCredentials(file);
        Check(loaded.Username == "ensemble" && loaded.Realm == "Ensemble"
            && loaded.Version == 1, "verifier identity changed");
        Check(Convert.FromHexString(loaded.Ha1).Length == 32, "HA1 is not a SHA-256 digest");
        Check(loaded.Ha1 == loaded.Ha1.ToLowerInvariant(), "HA1 is not lowercase");
        if (!OperatingSystem.IsWindows())
        {
            Check((int)File.GetUnixFileMode(file) == 0b110_000_000,
                "verifier file is not mode 0600");
        }
        Refuses(() => Gui.CreateDigestCredentials(file, "ensemble", false),
            "an existing verifier was replaced without reset");
        Check(Gui.CreateDigestCredentials(file, "ensemble", true) != password,
            "auth-reset reused the previous password");
        foreach (string bad in new[] { "", " ", "has space", "ha#sh", new string('x', 257) })
        {
            string copy = Path.Combine(work, "user-" + Math.Abs(bad.GetHashCode()) + ".json");
            Refuses(() => Gui.CreateDigestCredentials(copy, bad, false),
                "an invalid username was accepted");
        }
        var current = Gui.LoadDigestCredentials(file); // Reset issued a new digest.
        foreach (string tamper in new[] { "mode", "extra", "duplicate", "realm", "version",
            "upper", "short", "garbage" })
        {
            string copy = Path.Combine(work, "tamper-" + tamper + ".json");
            File.Copy(file, copy);
            switch (tamper)
            {
                case "mode":
                    if (OperatingSystem.IsWindows())
                    {
                        continue; // Mode bits do not exist to tamper with on Windows.
                    }
                    File.SetUnixFileMode(copy, UnixFileMode.UserRead | UnixFileMode.UserWrite
                        | UnixFileMode.GroupRead);
                    break;
                case "extra":
                    File.WriteAllText(copy, File.ReadAllText(copy).TrimEnd()[..^1]
                        + ",\"extra\":1}\n");
                    break;
                case "duplicate":
                    File.WriteAllText(copy, File.ReadAllText(copy).TrimEnd()[..^1]
                        + ",\"username\":\"ensemble\"}\n");
                    break;
                case "realm":
                    File.WriteAllText(copy, File.ReadAllText(copy)
                        .Replace("\"realm\":\"Ensemble\"", "\"realm\":\"Elsewhere\""));
                    break;
                case "version":
                    File.WriteAllText(copy, File.ReadAllText(copy)
                        .Replace("\"version\":1", "\"version\":2"));
                    break;
                case "upper":
                    File.WriteAllText(copy, File.ReadAllText(copy).ToUpperInvariant()
                        .Replace("\"VERSION\":2", "\"version\":1"));
                    break;
                case "short":
                    File.WriteAllText(copy, File.ReadAllText(copy).Replace(current.Ha1,
                        current.Ha1[..62]));
                    break;
                case "garbage":
                    File.AppendAllText(copy, "trailing garbage\n");
                    if (!OperatingSystem.IsWindows())
                    {
                        File.SetUnixFileMode(copy, UnixFileMode.UserRead | UnixFileMode.UserWrite);
                    }
                    break;
            }
            Refuses(() => Gui.LoadDigestCredentials(copy),
                "a tampered verifier file was accepted (" + tamper + ")");
        }
    }

    static void Parser()
    {
        string good = Answer("ensemble", "Ensemble", "pw", "GET", "/x", "nonce-value");
        var fields = Gui.ParseDigestFields(good);
        Check(fields is not null && fields["realm"] == "Ensemble" && fields["nc"] == "00000001"
            && fields["username"] == "ensemble",
            "a valid Authorization header was rejected: [" + good + "] -> "
            + (fields is null ? "null" : string.Join(",", fields.Keys.OrderBy(k => k))));
        Check(Gui.ParseDigestFields(null) is null, "a missing header parsed");
        Check(Gui.ParseDigestFields("Basic realm=\"Ensemble\"") is null, "another scheme parsed");
        Check(Gui.ParseDigestFields("Digest") is null, "a scheme-only header parsed");
        Check(Gui.ParseDigestFields(good + ", realm=\"Other\"") is null,
            "a duplicate field was accepted");
        Check(Gui.ParseDigestFields(good + ", username=\"control\u0001character\"") is null,
            "a control character was accepted");
        Check(Gui.ParseDigestFields(good + ",") is null, "a trailing comma was accepted");
        Check(Gui.ParseDigestFields(good[..^1]) is null, "an unclosed quote parsed");
        Check(Gui.ParseDigestFields("Digest x=" + new string('a', 9000)) is null,
            "an oversized header was accepted");
        // Unquoted values are RFC 7215 tokens, which exclude "/": a uri must be
        // quoted, while these token-safe fields need not be.
        var unquoted = Gui.ParseDigestFields("Digest realm=Ensemble, qop=auth, algorithm=SHA-256,"
            + " nonce=n, uri=\"/\", cnonce=c, nc=00000001, username=ensemble, response="
            + new string('0', 64));
        Check(unquoted is not null && unquoted["realm"] == "Ensemble",
            "an unquoted value did not parse");
    }

    static void Verification()
    {
        var credentials = new GuiDigestCredentials(1, "ensemble", "Ensemble",
            Hex(SHA256.HashData(Encoding.UTF8.GetBytes("ensemble:Ensemble:s3cret"))));
        Check(Gui.DigestResponseMatches(credentials, "GET", "/x", Fields()),
            "a correct answer was rejected");
        Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", Fields(password: "guess")),
            "a wrong password was accepted");
        Check(!Gui.DigestResponseMatches(credentials, "POST", "/x", Fields()),
            "the method was not covered by the digest");
        Check(!Gui.DigestResponseMatches(credentials, "GET", "/y", Fields()),
            "the uri was not covered by the digest");
        Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", Fields(realm: "Elsewhere")),
            "another realm was accepted");
        Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", Fields(username: "root")),
            "another username was accepted");
        Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", Fields(qop: "auth-int")),
            "another qop was accepted");
        Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", Fields(algorithm: "MD5")),
            "MD5 was accepted");
        foreach (string nc in new[] { "0000001", "0000000g", "00000000", "0000000000" })
        {
            Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", Fields(nc: nc)),
                "nc " + nc + " was accepted");
        }
        Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", Fields(cnonce: "")),
            "an empty cnonce was accepted");
        foreach (string extra in new[] { "userhash", "username*" })
        {
            var mutated = Fields();
            mutated[extra] = extra == "userhash" ? "true" : "escapes";
            Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", mutated),
                extra + " was accepted");
        }
        foreach (string response in new[] { "", "zz", new string('0', 63), new string('0', 65) })
        {
            var mutated = Fields();
            mutated["response"] = response;
            Check(!Gui.DigestResponseMatches(credentials, "GET", "/x", mutated),
                "a malformed response digest was accepted");
        }
    }

    // Naming the rule a request failed is what makes a LAN authentication failure
    // debuggable, so each rule reports itself rather than sharing one blunt "no".
    static void RejectionReasons()
    {
        var credentials = new GuiDigestCredentials(1, "ensemble", "Ensemble",
            Hex(SHA256.HashData(Encoding.UTF8.GetBytes("ensemble:Ensemble:s3cret"))));
        Check(Gui.DigestRejection(credentials, "GET", "/x", Fields()) is null,
            "a correct answer reported a reason");
        Check(Gui.DigestRejection(credentials, "GET", "/x", Fields(password: "guess"))
            == "the password does not match the verifier", "a wrong password was not named");
        Check(Gui.DigestRejection(credentials, "GET", "/x", Fields(username: "root"))!
            .StartsWith("username \"root\" is not the configured user", StringComparison.Ordinal),
            "a wrong username was not named");
        Check(Gui.DigestRejection(credentials, "GET", "/x", Fields(realm: "Elsewhere"))!
            .StartsWith("realm ", StringComparison.Ordinal), "a wrong realm was not named");
        Check(Gui.DigestRejection(credentials, "GET", "/y", Fields())!
            .StartsWith("uri ", StringComparison.Ordinal), "a mismatched uri was not named");
        // The method is inside the hash, not a separate field, so it surfaces as a
        // digest mismatch rather than a named field.
        Check(Gui.DigestRejection(credentials, "POST", "/x", Fields())
            == "the password does not match the verifier",
            "a changed method did not surface as a digest mismatch");
        Check(Gui.DigestRejection(credentials, "GET", "/x", Fields(nc: "0000000g"))!
            .Contains("replay counter", StringComparison.Ordinal), "a bad nc was not named");
        Check(Gui.DigestRejection(credentials, "GET", "/x", Fields(qop: "auth-int"))!
            .StartsWith("qop must be auth", StringComparison.Ordinal), "a bad qop was not named");
        Check(Gui.DigestRejection(credentials, "GET", "/x", Fields(algorithm: "MD5"))!
            .StartsWith("algorithm must be SHA-256", StringComparison.Ordinal),
            "a bad algorithm was not named");
        var missing = Fields();
        missing.Remove("nonce");
        Check(Gui.DigestRejection(credentials, "GET", "/x", missing)
            == "the Authorization header omits nonce", "a missing field was not named");
        var extension = Fields();
        extension["username*"] = "escapes";
        Check(Gui.DigestRejection(credentials, "GET", "/x", extension)
            == "the username* extension is not accepted", "username* was not named");
        // The bool view and the stated reason share one implementation and must agree.
        foreach (var probe in new[] { Fields(), Fields(password: "guess"),
            Fields(username: "root"), Fields(qop: "auth-int") })
        {
            Check(Gui.DigestResponseMatches(credentials, "GET", "/x", probe)
                == (Gui.DigestRejection(credentials, "GET", "/x", probe) is null),
                "the bool answer and the stated reason disagree");
        }
        // Client-chosen text reaches the operator's log, so it is stripped of control
        // characters and bounded: no forged terminal output and no unbounded string.
        var sneaky = Fields(username: "a\u001b]0;pwned\u0007" + new string('x', 200));
        string reason = Gui.DigestRejection(credentials, "GET", "/x", sneaky)!;
        Check(!reason.Any(character => char.IsControl(character)) && reason.Length < 160,
            "client-chosen text reached the log unsanitized");
    }

    // A terminal that says only "Unauthorized" cannot tell a wrong password from a
    // verifier file the listener never read, which is exactly the failure that cost
    // an operator an afternoon. The reason and the file must both be said.
    static void GuardDiagnostics()
    {
        var credentials = new GuiDigestCredentials(1, "ensemble", "Ensemble",
            Hex(SHA256.HashData(Encoding.UTF8.GetBytes("ensemble:Ensemble:s3cret"))));
        byte[] key = RandomNumberGenerator.GetBytes(32);
        var log = new StringWriter();
        var guard = new Gui.DigestGuard(credentials, key, "/var/lib/ensemble/gui-auth.json", log);
        DateTimeOffset now = DateTimeOffset.UtcNow;
        void Refuse(string authorization)
        {
            var context = new DefaultHttpContext();
            context.Connection.RemoteIpAddress = IPAddress.Loopback;
            context.Request.Method = "GET";
            context.Request.Path = "/x";
            context.Request.Headers["Authorization"] = authorization;
            Check(!guard.AdmitAsync(context).GetAwaiter().GetResult(),
                "the guard admitted a request it should refuse");
            Check(context.Response.StatusCode == 401,
                "a refusal was not the uniform 401 challenge");
        }
        Refuse(Answer("ensemble", "Ensemble", "not-the-secret", "GET", "/x",
            Gui.DigestNonce(key, now)));
        Refuse(Answer("root", "Ensemble", "s3cret", "GET", "/x", Gui.DigestNonce(key, now)));
        string text = log.ToString();
        Check(text.Contains("the password does not match the verifier"),
            "a wrong password was not explained to the operator");
        Check(text.Contains("username \"root\" is not the configured user"),
            "a wrong username was not explained to the operator");
        Check(text.Contains("/var/lib/ensemble/gui-auth.json"),
            "the verifier file was never named");
        Check(text.Split("verifier file").Length - 1 == 1,
            "the verifier path was repeated instead of stated once");
        Check(text.Split("auth: rejected").Length - 1 == 2,
            "each refusal was not logged exactly once");
        // Secrets stay out of the log: neither password, and no digest material.
        Check(!text.Contains("s3cret") && !text.Contains("not-the-secret"),
            "a password reached the operator log");
        // Silent hosts, such as fixtures, write nothing at all.
        var quiet = new Gui.DigestGuard(credentials, key);
        var anonymous = new DefaultHttpContext();
        anonymous.Connection.RemoteIpAddress = IPAddress.Loopback;
        anonymous.Request.Method = "GET";
        anonymous.Request.Path = "/x";
        Check(!quiet.AdmitAsync(anonymous).GetAwaiter().GetResult(),
            "a request that should be refused was admitted");
    }

    static void Nonce()
    {
        byte[] key = RandomNumberGenerator.GetBytes(32);
        byte[] other = RandomNumberGenerator.GetBytes(32);
        DateTimeOffset now = DateTimeOffset.UtcNow;
        string nonce = Gui.DigestNonce(key, now);
        // The issued time is stored in whole seconds, so the lifetime is truncated
        // to the second rather than matching a fractional clock exactly.
        var expiry = Gui.DigestNonceExpiry(key, nonce, now);
        Check(expiry is not null && Math.Abs((expiry.Value - now).TotalSeconds - 300) <= 1,
            "a fresh nonce did not report its five-minute lifetime");
        Check(Gui.DigestNonceExpiry(key, nonce, now.AddSeconds(298))!.Value > now,
            "a still-valid nonce was refused");
        // Reporting expiry is the helper's job; refusing an expired nonce belongs to
        // the request path, which answers it with a stale challenge instead.
        var late = Gui.DigestNonceExpiry(key, nonce, now.AddSeconds(302));
        Check(late is not null && late.Value <= now.AddSeconds(302),
            "a spent nonce did not report a passed expiry");
        Check(Gui.DigestNonceExpiry(key, nonce, now.AddSeconds(-10)) is null,
            "a future-dated nonce was accepted");
        Check(Gui.DigestNonceExpiry(other, nonce, now) is null,
            "another server's nonce was accepted");
        byte[] raw = FromBase64Url(nonce);
        raw[40] ^= 0x01;
        Check(Gui.DigestNonceExpiry(key, ToBase64Url(raw), now) is null,
            "a forged nonce was accepted");
        Check(Gui.DigestNonceExpiry(key, "not-base64!!", now) is null,
            "an undecodable nonce was accepted");
        Check(Gui.DigestNonceExpiry(key, ToBase64Url(raw[..55]), now) is null,
            "a short nonce was accepted");
    }

    static void Staleness()
    {
        var credentials = new GuiDigestCredentials(1, "ensemble", "Ensemble",
            Hex(SHA256.HashData(Encoding.UTF8.GetBytes("ensemble:Ensemble:s3cret"))));
        byte[] key = RandomNumberGenerator.GetBytes(32);
        var guard = new Gui.DigestGuard(credentials, key);
        DateTimeOffset now = DateTimeOffset.UtcNow;
        async Task<DefaultHttpContext> Send(string? authorization)
        {
            var context = new DefaultHttpContext();
            context.Connection.RemoteIpAddress = IPAddress.Loopback;
            context.Request.Method = "GET";
            context.Request.Path = "/x";
            if (authorization is not null)
            {
                context.Request.Headers["Authorization"] = authorization;
            }
            Check(!await guard.AdmitAsync(context), "the guard admitted a request it should refuse");
            return context;
        }
        // A spent nonce gets a fresh challenge marked stale.
        var stale = Send(Answer("ensemble", "Ensemble", "s3cret", "GET", "/x",
            Gui.DigestNonce(key, now.AddSeconds(-301)))).GetAwaiter().GetResult();
        Check(stale.Response.StatusCode == 401
            && stale.Response.Headers.WWWAuthenticate.ToString().Contains("stale=true"),
            "a spent nonce was not answered with a stale challenge");
        // Reusing one accepted answer is a replay, and unlike staleness it is a
        // failure: it must not be admitted twice.
        string live = Gui.DigestNonce(key, now);
        string answer = Answer("ensemble", "Ensemble", "s3cret", "GET", "/x", live,
            cnonce: "once");
        var first = new DefaultHttpContext();
        first.Connection.RemoteIpAddress = IPAddress.Loopback;
        first.Request.Method = "GET";
        first.Request.Path = "/x";
        first.Request.Headers["Authorization"] = answer;
        Check(guard.AdmitAsync(first).GetAwaiter().GetResult(), "a correct answer was refused");
        var second = Send(answer).GetAwaiter().GetResult();
        Check(second.Response.StatusCode == 401
            && !second.Response.Headers.WWWAuthenticate.ToString().Contains("stale=true"),
            "a replay was answered like a stale nonce");
        // Staleness alone must not build toward the slowdown: three spent nonces
        // still leave the address able to authenticate immediately.
        var fresh = new DefaultHttpContext();
        fresh.Connection.RemoteIpAddress = IPAddress.Loopback;
        fresh.Request.Method = "GET";
        fresh.Request.Path = "/x";
        fresh.Request.Headers["Authorization"] = Answer("ensemble", "Ensemble", "s3cret",
            "GET", "/x", Gui.DigestNonce(key, now), cnonce: "later");
        Check(guard.AdmitAsync(fresh).GetAwaiter().GetResult(),
            "stale nonces were counted as failures and slowed the address down");
    }

    // The verifier is read once when the listener starts, so a reset that lands while
    // the listener runs used to be refused forever and reported as a wrong password.
    // This drives the real listener: replace the file, then prove the new password is
    // accepted without a restart, that a credential minted before the swap no longer
    // works, and that a verifier which disappears denies access instead of trusting
    // the last known good copy.
    static async Task ReloadAndFault(string work)
    {
        string file = Path.Combine(work, "reloaded-auth.json");
        string first = Gui.CreateDigestCredentials(file, "ensemble", false);
        var root = new Ensemble.Ensemble(TextWriter.Null,
            dataDirectory: Path.Combine(work, "reload-fleet"));
        var agent = new Agent(root, new Configuration(Vendor.OpenAI, "offline",
            "http://127.0.0.1:1"), name: "reloading");
        var probe = new TcpListener(IPAddress.Loopback, 0);
        probe.Start();
        int port = ((IPEndPoint)probe.LocalEndpoint).Port;
        probe.Stop();
        string url = "http://127.0.0.1:" + port;
        var builder = WebApplication.CreateBuilder();
        builder.Logging.ClearProviders();
        builder.WebHost.UseUrls(url);
        await using var host = builder.Build();
        var report = new StringWriter();
        using var stop = new CancellationTokenSource(TimeSpan.FromSeconds(60));
        Task serving = Gui.ServeAsync(host, agent, new GuiOptions(port,
            SettingsPath: Path.Combine(work, "reload-settings.json"), AuthPath: file,
            DiagnosticsOutput: report), stop.Token);
        if (serving.IsCompleted)
        {
            await serving;
        }
        await host.StartAsync();
        using var http = new HttpClient(new HttpClientHandler { UseProxy = false });
        async Task<(HttpStatusCode Status, string Challenge)> Ask(
            string? authorization = null)
        {
            using var request = new HttpRequestMessage(HttpMethod.Get, url + "/");
            if (authorization is not null)
            {
                request.Headers.TryAddWithoutValidation("Authorization", authorization);
            }
            using var reply = await http.SendAsync(request,
                HttpCompletionOption.ResponseHeadersRead, CancellationToken.None);
            return (reply.StatusCode, reply.Headers.WwwAuthenticate.FirstOrDefault()?.ToString() ?? "");
        }
        string Nonce(string challenge)
        {
            int start = challenge.IndexOf("nonce=\"", StringComparison.Ordinal) + 7;
            int end = challenge.IndexOf('"', start);
            Check(start > 6 && end > start, "the challenge carried no nonce: " + challenge);
            return challenge[start..end];
        }
        var opening = await Ask();
        Check(opening.Status == HttpStatusCode.Unauthorized, "the reload test started open");
        string staleNonce = Nonce(opening.Challenge);
        var known = await Ask(Answer("ensemble", "Ensemble", first, "GET", "/", staleNonce));
        Check(known.Status == HttpStatusCode.OK, "the original password was refused");
        string second = Gui.CreateDigestCredentials(file, "ensemble", true);
        await Task.Delay(1300);
        var beforeSwap = await Ask();
        var stale = await Ask(Answer("ensemble", "Ensemble", second, "GET", "/",
            staleNonce));
        Check(stale.Status == HttpStatusCode.Unauthorized,
            "a nonce issued before the verifier swap was still honoured");
        var old = await Ask(Answer("ensemble", "Ensemble", first, "GET", "/",
            Nonce(beforeSwap.Challenge)));
        Check(old.Status == HttpStatusCode.Unauthorized,
            "the retired password still opened the listener");
        var next = await Ask();
        var renewed = await Ask(Answer("ensemble", "Ensemble", second, "GET", "/",
            Nonce(next.Challenge)));
        Check(renewed.Status == HttpStatusCode.OK,
            "the reset password was refused without restarting the listener: "
                + renewed.Status + " " + renewed.Challenge);
        File.Delete(file);
        await Task.Delay(1300);
        // A faulted verifier answers before any Digest parsing, so no challenge is
        // offered: the point is that even a well-formed answer cannot get through.
        var missing = await Ask(Answer("ensemble", "Ensemble", second, "GET", "/",
            staleNonce));
        Check(missing.Status == HttpStatusCode.ServiceUnavailable,
            "a deleted verifier left the listener serving on its last known good copy: "
                + missing.Status);
        // Recreated with init, not reset: a reset refuses to invent a verifier it
        // cannot first read, which is the same fail-closed rule as the guard above.
        string third = Gui.CreateDigestCredentials(file, "ensemble", false);
        await Task.Delay(1300);
        var restored = await Ask(Answer("ensemble", "Ensemble", third, "GET", "/",
            Nonce((await Ask()).Challenge)));
        Check(restored.Status == HttpStatusCode.OK,
            "the listener stayed locked out after the verifier came back: "
                + restored.Status);
        string log = report.ToString();
        Check(log.Contains("verifier reloaded from " + file),
            "the verifier swap was silent: " + log);
        stop.Cancel();
        await host.StopAsync();
        try
        {
            await serving;
        }
        catch (OperationCanceledException)
        {
        }
    }

    // A password printed beside the wrong file simply does not work, and the old report
    // said only where it had written, never how it chose that place. The source is now
    // named, and a fallback to the current directory is called out because that is the
    // case where a listener started with --data-dir reads somewhere else entirely.
    static void CommandPaths(string work)
    {
        string original = Directory.GetCurrentDirectory();
        string fromFlag = Path.Combine(work, "cmd-flag");
        string fromEnv = Path.Combine(work, "cmd-env");
        string fallback = Path.Combine(work, "cmd-cwd");
        Directory.CreateDirectory(fromFlag);
        Directory.CreateDirectory(fromEnv);
        Directory.CreateDirectory(fallback);
        string saved = Environment.GetEnvironmentVariable("EN_DATA_DIR") ?? "";
        try
        {
            var stdout = new StringWriter();
            var log = new StringWriter();
            Gui.RunAuthCommand(["auth-init", "--data-dir", fromFlag], stdout, log);
            Check(File.Exists(Path.Combine(fromFlag, "gui-auth.json")),
                "--data-dir did not place the verifier beside the data directory");
            Check(log.ToString().Contains("resolved from --data-dir"),
                "the report did not name --data-dir as the source");
            string printed = stdout.ToString();
            Check(printed.Trim().Length > 8 && !printed.Contains("auth:")
                && printed.Split('\n').Length == 2,
                "stdout must carry the password alone, for $(ensemble auth-init): " + printed);

            Environment.SetEnvironmentVariable("EN_DATA_DIR", fromEnv);
            log = new StringWriter();
            Gui.RunAuthCommand(["auth-init"], TextWriter.Null, log);
            Check(File.Exists(Path.Combine(fromEnv, "gui-auth.json")),
                "EN_DATA_DIR was ignored");
            Check(log.ToString().Contains("resolved from EN_DATA_DIR"),
                "the report did not name EN_DATA_DIR as the source");

            Environment.SetEnvironmentVariable("EN_DATA_DIR", null);
            Directory.SetCurrentDirectory(fallback);
            log = new StringWriter();
            Gui.RunAuthCommand(["auth-init"], TextWriter.Null, log);
            Check(File.Exists(Path.Combine(fallback, "gui-auth.json")),
                "the current directory fallback did not write the verifier");
            string text = log.ToString();
            Check(text.Contains("resolved from the current directory")
                && text.Contains("reads a different file"),
                "the silent current-directory fallback was not called out: " + text);

            // The other direction: writing to the real data directory while an unread
            // verifier sits beside the operator must still be reported.
            log = new StringWriter();
            Gui.RunAuthCommand(["auth-reset", "--data-dir", fromFlag], TextWriter.Null, log);
            Check(log.ToString().Contains(Path.Combine(fallback, "gui-auth.json"))
                && log.ToString().Contains("also exists"),
                "an unused verifier in the current directory was not warned about: " + log);
        }
        finally
        {
            Directory.SetCurrentDirectory(original);
            Environment.SetEnvironmentVariable("EN_DATA_DIR",
                string.IsNullOrEmpty(saved) ? null : saved);
        }
    }


    static void PathPolicy(string work)
    {
        string data = Path.Combine(work, "data");
        Directory.CreateDirectory(data);
        Check(Gui.DigestPath(data, null, false, true) is null,
            "a loopback listener without credentials required authentication");
        Refuses(() => Gui.DigestPath(data, null, false, false),
            "a non-loopback listener started without credentials");
        string file = Path.Combine(data, "gui-auth.json");
        Gui.CreateDigestCredentials(file, "ensemble", false);
        Check(Gui.DigestPath(data, null, false, true) == file,
            "an existing verifier did not enable authentication");
        Refuses(() => Gui.DigestPath(data, Path.Combine(data, "missing.json"), true, true),
            "an explicit missing verifier was ignored");
        Refuses(() => Gui.DigestPath(data, "   ", true, true),
            "a blank explicit verifier path was accepted");
        string save = Path.Combine(work, "save.json");
        File.WriteAllText(save, "{}");
        Refuses(() => Gui.CheckDigestPath(save, [save]),
            "the verifier file was allowed to reuse application state");
        Gui.CheckDigestPath(file, [save, ""]);
    }

    static async Task Listener(string work)
    {
        string file = Path.Combine(work, "live-auth.json");
        string password = Gui.CreateDigestCredentials(file, "ensemble", false);
        var root = new Ensemble.Ensemble(TextWriter.Null,
            dataDirectory: Path.Combine(work, "fleet"));
        var agent = new Agent(root, new Configuration(Vendor.OpenAI, "offline",
            "http://127.0.0.1:1"), name: "guarded");
        var probe = new TcpListener(IPAddress.Loopback, 0);
        probe.Start();
        int port = ((IPEndPoint)probe.LocalEndpoint).Port;
        probe.Stop();
        string url = "http://127.0.0.1:" + port;
        var builder = WebApplication.CreateBuilder();
        builder.Logging.ClearProviders();
        builder.WebHost.UseUrls(url);
        await using var host = builder.Build();
        using var stop = new CancellationTokenSource(TimeSpan.FromSeconds(45));
        Task serving = Gui.ServeAsync(host, agent, new GuiOptions(port,
            SettingsPath: Path.Combine(work, "live-settings.json"), AuthPath: file), stop.Token);
        if (serving.IsCompleted)
        {
            await serving;
        }
        await host.StartAsync();
        using var http = new HttpClient(new HttpClientHandler { UseProxy = false });
        async Task<(HttpStatusCode Status, string Challenge, int? Retry)> Ask(
            string? authorization = null)
        {
            using var request = new HttpRequestMessage(HttpMethod.Get, url + "/");
            if (authorization is not null)
            {
                request.Headers.TryAddWithoutValidation("Authorization", authorization);
            }
            using var reply = await http.SendAsync(request, HttpCompletionOption.ResponseHeadersRead,
                CancellationToken.None);
            int? retry = reply.Headers.TryGetValues("Retry-After", out var values)
                ? int.Parse(values.Single()) : null;
            return (reply.StatusCode,
                reply.Headers.WwwAuthenticate.FirstOrDefault()?.ToString() ?? "", retry);
        }
        string Nonce(string challenge)
        {
            int start = challenge.IndexOf("nonce=\"", StringComparison.Ordinal) + 7;
            int end = challenge.IndexOf('"', start);
            Check(start > 6 && end > start, "the challenge carried no nonce: " + challenge);
            return challenge[start..end];
        }
        var anonymous = await Ask();
        Check(anonymous.Status == HttpStatusCode.Unauthorized,
            "an unauthenticated page request was served");
        Check(anonymous.Challenge.StartsWith("Digest realm=\"Ensemble\"")
            && anonymous.Challenge.Contains("algorithm=SHA-256")
            && anonymous.Challenge.Contains("qop=\"auth\"")
            && anonymous.Challenge.Contains("charset=UTF-8"),
            "the challenge is not SHA-256 Digest auth: " + anonymous.Challenge);
        var guessed = await Ask(Answer("ensemble", "Ensemble", "guessed", "GET", "/",
            Nonce(anonymous.Challenge)));
        Check(guessed.Status == HttpStatusCode.Unauthorized, "a guessed password was accepted");
        var fresh = await Ask();
        var accepted = await Ask(Answer("ensemble", "Ensemble", password, "GET", "/",
            Nonce(fresh.Challenge)));
        Check(accepted.Status == HttpStatusCode.OK,
            "the correct credentials were refused: " + accepted.Status);
        var replay = await Ask(Answer("ensemble", "Ensemble", password, "GET", "/",
            Nonce(fresh.Challenge)));
        Check(replay.Status == HttpStatusCode.Unauthorized, "a replayed answer was accepted");
        // Keep guessing until the address is slowed down. A blocked address is also
        // refused before any challenge is issued, which is itself the proof we want.
        bool slowed = false;
        for (int attempt = 0; attempt < 12 && !slowed; attempt++)
        {
            var challenge = await Ask();
            if (challenge.Status == HttpStatusCode.TooManyRequests)
            {
                Check(challenge.Retry is > 0, "the slowdown carried no Retry-After");
                slowed = true;
                break;
            }
            Check(challenge.Status == HttpStatusCode.Unauthorized
                && challenge.Challenge.Contains("nonce="),
                "an anonymous request returned " + challenge.Status);
            var slow = await Ask(Answer("ensemble", "Ensemble", "guessed", "GET", "/",
                Nonce(challenge.Challenge)));
            if (slow.Status == HttpStatusCode.ServiceUnavailable)
            {
                throw new InvalidOperationException("the failure table reported busy, not rate limited");
            }
            if (slow.Status == HttpStatusCode.TooManyRequests)
            {
                Check(slow.Retry is > 0, "the slowdown carried no Retry-After");
                slowed = true;
                break;
            }
            Check(slow.Status == HttpStatusCode.Unauthorized,
                "a guessed password was accepted during the slowdown");
        }
        Check(slowed, "repeated failures never triggered the slowdown");
        // Whether the address is still in slowdown or not, an unauthenticated upgrade
        // must never reach a handler: 401 or 429 both prove the middleware intercepted.
        var websocket = await Ask();
        Check(websocket.Status is HttpStatusCode.Unauthorized
            or HttpStatusCode.TooManyRequests,
            "an unauthenticated page/upgrade request was served: " + websocket.Status);
        stop.Cancel();
        await host.StopAsync();
        try
        {
            await serving;
        }
        catch (OperationCanceledException)
        {
        }
    }
}
