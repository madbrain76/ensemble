using System.Net;
using System.Net.NetworkInformation;
using System.Net.Sockets;
using System.Security.Cryptography;
using System.Text;
using Ensemble;
using EnsembleGui;

// Where the listener may bind. Loopback stays the default and needs no
// credential; anything beyond it is a deliberate LAN exposure and therefore
// requires Digest authentication.
static class HostBindingTests
{
    static void Check(bool value, string message)
    {
        if (!value)
        {
            throw new InvalidOperationException(message);
        }
    }

    static int FreePort()
    {
        var listener = new TcpListener(IPAddress.Any, 0);
        listener.Start();
        int port = ((IPEndPoint)listener.LocalEndpoint).Port;
        listener.Stop();
        return port;
    }

    static Agent LocalAgent(string work, string name)
    {
        var root = new Ensemble.Ensemble(TextWriter.Null,
            dataDirectory: Path.Combine(work, "fleet-" + name));
        return new Agent(root, new Configuration(Vendor.OpenAI, "offline",
            "http://127.0.0.1:1"), name: name);
    }

    static string? LanAddress()
    {
        return NetworkInterface.GetAllNetworkInterfaces()
            .Where(adapter => adapter.OperationalStatus == OperationalStatus.Up
                && adapter.NetworkInterfaceType != NetworkInterfaceType.Loopback)
            .SelectMany(adapter => adapter.GetIPProperties().UnicastAddresses)
            .Select(candidate => candidate.Address)
            .FirstOrDefault(address => address.AddressFamily == AddressFamily.InterNetwork
                && !IPAddress.IsLoopback(address))?.ToString();
    }

    static string Hex(byte[] data) => Convert.ToHexString(data).ToLowerInvariant();

    public static async Task RunAsync()
    {
        string work = Path.Combine(Path.GetTempPath(), "gui-host-" + Guid.NewGuid().ToString("N"));
        Directory.CreateDirectory(work);
        try
        {
            // Accepted spellings.
            Check(IPAddress.IsLoopback(Gui.ResolveBindAddress("127.0.0.1"))
                && IPAddress.IsLoopback(Gui.ResolveBindAddress("localhost"))
                && IPAddress.IsLoopback(Gui.ResolveBindAddress("LOCALHOST"))
                && IPAddress.IsLoopback(Gui.ResolveBindAddress("::1")),
                "a loopback spelling was not recognized");
            Check(Gui.ResolveBindAddress("0.0.0.0").Equals(IPAddress.Any),
                "the unspecified address did not resolve");
            Check(Gui.ResolveBindAddress("192.168.7.19").ToString() == "192.168.7.19",
                "a literal address did not resolve");
            Check(Gui.IsLoopbackHost("localhost") && !Gui.IsLoopbackHost("0.0.0.0"),
                "the loopback decision is wrong");
            // Refused spellings: a name, a port, a URL, or nonsense must not bind.
            foreach (string bad in new[] { "", "   ", "example.com", "localhost:8090",
                "127.0.0.1:8090", "http://127.0.0.1", "127.0.0.256", "not-an-address", "::gg" })
            {
                bool refused = false;
                try
                {
                    Gui.ResolveBindAddress(bad);
                }
                catch (ArgumentException)
                {
                    refused = true;
                }
                Check(refused, "accepted the bind address \"" + bad + "\"");
            }

            // Beyond loopback without a credential, the listener must refuse to start.
            bool started = false;
            string refusal = "";
            try
            {
                using var deadline = new CancellationTokenSource(TimeSpan.FromSeconds(20));
                await Gui.RunAsync(LocalAgent(work, "lan-open"),
                    new GuiOptions(FreePort(), Host: "0.0.0.0",
                        SettingsPath: Path.Combine(work, "open-settings.json")),
                    deadline.Token);
                started = true;
            }
            catch (InvalidOperationException failure)
            {
                refusal = failure.Message;
            }
            Check(!started, "an unauthenticated listener bound beyond loopback");
            Check(refusal.Contains("authentication", StringComparison.OrdinalIgnoreCase),
                "the LAN bind was refused for an unrelated reason: " + refusal);
            // Loopback without a credential still starts: LAN protection is opt-in.
            int openPort = FreePort();
            using (var loopbackDeadline = new CancellationTokenSource(TimeSpan.FromSeconds(20)))
            {
                Task local = Gui.RunAsync(LocalAgent(work, "local-open"),
                    new GuiOptions(openPort, SettingsPath: Path.Combine(work, "local-settings.json")),
                    loopbackDeadline.Token);
                await Task.WhenAny(local, Task.Delay(3000));
                await RequireRunning(local, "the default loopback listener");
                loopbackDeadline.Cancel();
                try
                {
                    await local;
                }
                catch (OperationCanceledException)
                {
                }
            }

            // With credentials, a LAN bind serves the page and still authenticates.
            string file = Path.Combine(work, "gui-auth.json");
            string password = Gui.CreateDigestCredentials(file, "ensemble", false);
            int lanPort = FreePort();
            string? lan = LanAddress();
            using var lanDeadline = new CancellationTokenSource(TimeSpan.FromSeconds(30));
            Task serving = Gui.RunAsync(LocalAgent(work, "lan-guarded"),
                new GuiOptions(lanPort, Host: "0.0.0.0", AuthPath: file,
                    SettingsPath: Path.Combine(work, "lan-settings.json")), lanDeadline.Token);
            await Task.WhenAny(serving, Task.Delay(3000));
            await RequireRunning(serving, "the authenticated LAN listener");
            using (var http = new HttpClient(new HttpClientHandler { UseProxy = false }))
            {
                var targets = new List<string> { "127.0.0.1" };
                if (lan is not null)
                {
                    targets.Add(lan); // Proves the socket is genuinely reachable off-loopback.
                }
                foreach (string target in targets)
                {
                    using var anonymous = await http.GetAsync($"http://{target}:{lanPort}/");
                    Check(anonymous.StatusCode == HttpStatusCode.Unauthorized,
                        "the LAN bind served the page without credentials from " + target);
                    Check(anonymous.Headers.WwwAuthenticate.Any(value =>
                        value.Scheme == "Digest" && (value.Parameter ?? "").Contains("SHA-256")),
                        "the LAN bind did not issue a SHA-256 Digest challenge");
                }
                string? parameter;
                using (var unauthenticated = await http.GetAsync($"http://127.0.0.1:{lanPort}/"))
                {
                    parameter = unauthenticated.Headers.WwwAuthenticate.First().Parameter;
                }
                Check(parameter is not null, "the LAN challenge carried no parameters");
                string challenge = parameter!;
                int begin = challenge.IndexOf("nonce=\"", StringComparison.Ordinal) + 7;
                string nonce = challenge[begin..challenge.IndexOf('"', begin)];
                string ha1 = Hex(SHA256.HashData(Encoding.UTF8.GetBytes("ensemble:Ensemble:" + password)));
                string uri = "/";
                string requestHash = Hex(SHA256.HashData(Encoding.UTF8.GetBytes("GET:" + uri)));
                string cnonce = Guid.NewGuid().ToString("N");
                string response = Hex(SHA256.HashData(Encoding.UTF8.GetBytes(
                    ha1 + ":" + nonce + ":00000001:" + cnonce + ":auth:" + requestHash)));
                using var request = new HttpRequestMessage(HttpMethod.Get, $"http://127.0.0.1:{lanPort}/");
                request.Headers.TryAddWithoutValidation("Authorization",
                    "Digest realm=\"Ensemble\", nonce=\"" + nonce + "\", uri=\"" + uri
                    + "\", qop=\"auth\", algorithm=SHA-256, username=\"ensemble\", nc=00000001"
                    + ", cnonce=\"" + cnonce + "\", response=\"" + response + "\"");
                using var admitted = await http.SendAsync(request);
                Check(admitted.StatusCode == HttpStatusCode.OK,
                    "correct credentials were refused on a LAN bind: " + admitted.StatusCode);
            }
            lanDeadline.Cancel();
            try
            {
                await serving;
            }
            catch (OperationCanceledException)
            {
            }
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

    // A listener that started is still serving; anything that finished early failed
    // to start, and its real reason must surface instead of a timeout.
    static async Task RequireRunning(Task serving, string label)
    {
        if (!serving.IsCompleted)
        {
            return;
        }
        string reason;
        try
        {
            await serving;
            reason = "returned without an error";
        }
        catch (Exception failure)
        {
            reason = failure.GetType().Name + ": " + failure.Message;
        }
        throw new InvalidOperationException(label + " stopped before serving: " + reason);
    }
}
