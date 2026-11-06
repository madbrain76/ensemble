using System.Diagnostics;
using Ensemble;
using System.Net;
using System.Net.Sockets;
using System.Text;


// Waking a sleeping local provider: the address the operator supplies, the
// packet that wakes it, and the rule that an awake host is never woken twice.
internal static class WakeOnLanTests
{
    static void Check(bool value, string message)
    {
        if (!value)
        {
            throw new Exception(message);
        }
    }

    static byte[] Mac() => [0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff];

    public static async Task Run(string temp)
    {
        // Accepted MAC spellings.
        foreach (string spelling in new[] { "aa:bb:cc:dd:ee:ff", "AA:BB:CC:DD:EE:FF",
            "aa-bb-cc-dd-ee-ff", "aabbccddeeff", "  aa:bb:cc:dd:ee:ff  " })
        {
            Check(Engine.TryWakeMac(spelling, out byte[]? parsed)
                && parsed!.SequenceEqual(Mac()), "rejected the MAC spelling " + spelling);
        }
        foreach (string bad in new[] { "", "   ", "aa:bb:cc:dd:ee", "aa:bb:cc:dd:ee:ff:00",
            "zz:bb:cc:dd:ee:ff", "aabbccddeef", "aa:bb:cc:dd:ee:gg", "10.0.0.5" })
        {
            Check(!Engine.TryWakeMac(bad, out byte[]? rejected) && rejected is null,
                "accepted the MAC " + bad);
        }

        // The magic packet is 6 broadcast bytes then the MAC sixteen times.
        byte[] packet = Engine.WakePacket(Mac());
        Check(packet.Length == 102, "magic packet is " + packet.Length + " bytes, not 102");
        Check(packet.AsSpan(0, 6).ToArray().All(value => value == 0xff),
            "the magic packet has no broadcast preamble");
        for (int copy = 0; copy < 16; copy++)
        {
            Check(packet.AsSpan(6 + copy * 6, 6).SequenceEqual(Mac()),
                "the MAC is missing at repetition " + copy);
        }

        // Endpoint derivation follows the URL, never a guess.
        foreach ((string url, string host, int port) in new[]
                 {
                     ("http://llm:8110/v1", "llm", 8110), ("https://api.example.com/v1", "api.example.com", 443),
                     ("http://127.0.0.1", "127.0.0.1", 80), ("https://loopback.test:8443", "loopback.test", 8443)
                 })
        {
            Check(Engine.TryProviderEndpoint(url, out ProviderEndpoint? found)
                && found!.Host == host && found.Port == port,
                "derived the wrong endpoint from " + url);
        }
        foreach (string bad in new[] { "", "not a url", "ftp://llm:8110", "/relative/path" })
        {
            Check(!Engine.TryProviderEndpoint(bad, out _), "accepted " + bad);
        }

        // Configuration validates the MAC at the boundary and keeps it across a
        // model switch, because the address belongs to the machine, not the model.
        var configured = new Configuration(Vendor.OpenAI, "gpt-4o", "http://127.0.0.1:8110/v1",
            wakeOnLanMac: "aa:bb:cc:dd:ee:ff");
        Check(configured.WakeOnLanMac == "aa:bb:cc:dd:ee:ff", "the configured MAC was lost");
        Check(configured.WithModel("gpt-4o-mini").WakeOnLanMac == "aa:bb:cc:dd:ee:ff",
            "a model switch dropped the wake address");
        try
        {
            _ = new Configuration(Vendor.OpenAI, "gpt-4o", "http://127.0.0.1:8110/v1",
                wakeOnLanMac: "not-a-mac");
            Check(false, "an invalid wake MAC was accepted");
        }
        catch (ArgumentException)
        {
        }

        // A listening port is already awake; a closed one is not.
        var listener = new TcpListener(IPAddress.Loopback, 0);
        listener.Start();
        int open = ((IPEndPoint)listener.LocalEndpoint).Port;
        Check(await Engine.TcpReachableAsync(new ProviderEndpoint("127.0.0.1", open),
                CancellationToken.None), "a listening port was reported unreachable");
        listener.Stop();
        Check(!await Engine.TcpReachableAsync(new ProviderEndpoint("127.0.0.1", FreeTcpPort()),
            CancellationToken.None), "a closed port was reported reachable");

        // A sleeping network card hears nothing but a broadcast, and the kernel
        // refuses to deliver one unless the socket says so. Read the option back
        // rather than infer it: without it the send below still appears to succeed
        // on unicast while no wake can ever arrive.
        using (UdpClient wakeSocket = Engine.WakeSocket())
        {
            Check(wakeSocket.Client.EnableBroadcast,
                "the wake socket cannot send broadcasts, so the magic packet never "
                    + "leaves this machine");
        }

        // The packet really leaves the process.
        using (var receiver = new UdpClient(new IPEndPoint(IPAddress.Loopback, 0)))
        {
            int wakePort = ((IPEndPoint)receiver.Client.LocalEndPoint!).Port;
            Check(Engine.SendWakePacket("127.0.0.1", Mac(), wakePort, CancellationToken.None),
                "the magic packet was not delivered");
            receiver.Client.ReceiveTimeout = 3000;
            IPEndPoint? sender = null;
            byte[] arrived = receiver.Receive(ref sender);
            Check(arrived.SequenceEqual(packet), "the delivered packet is not the magic packet");
        }

        // ICMP is a hint, not a requirement: it is only asserted where ping exists.
        if (PingAvailable())
        {
            Check(await Engine.IcmpReachableAsync("127.0.0.1", CancellationToken.None),
                "loopback was not reachable by ping");
            Check(!await Engine.IcmpReachableAsync("no-such-host.invalid", CancellationToken.None),
                "an unresolvable host answered ping");
        }

        // No MAC means no probing at all: the ordinary provider pays nothing.
        var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Path.Combine(temp, "wake-off"));
        var agent = new Agent(root, new Configuration(Vendor.OpenAI, "gpt-4o",
            "http://127.0.0.1:1/v1"), name: "wake-off");
        var started = Stopwatch.GetTimestamp();
        await Engine.WakeUnavailableProviderAsync(agent, agent.Configuration(),
            CancellationToken.None);
        Check(Stopwatch.GetElapsedTime(started) < TimeSpan.FromSeconds(1),
            "an unconfigured provider was probed anyway");

        // An awake endpoint is never woken: an open port answers before any probe,
        // so the operator's machine is not disturbed and nothing is logged.
        var awake = new TcpListener(IPAddress.Loopback, 0);
        awake.Start();
        int awakePort = ((IPEndPoint)awake.LocalEndpoint).Port;
        var awakeLog = new StringWriter();
        var awakeRoot = new Ensemble.Ensemble(awakeLog,
            dataDirectory: Path.Combine(temp, "wake-awake"));
        var awakeAgent = new Agent(awakeRoot, new Configuration(Vendor.OpenAI, "gpt-4o",
            $"http://127.0.0.1:{awakePort}/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"),
            name: "wake-awake");
        await Engine.WakeUnavailableProviderAsync(awakeAgent, awakeAgent.Configuration(),
            CancellationToken.None, TimeSpan.FromSeconds(2));
        awake.Stop();
        Check(!awakeLog.ToString().Contains("Wake-on-LAN"),
            "a provider with an open port was woken anyway");

        // A host that answers neither probe is woken, and the wait never claims the
        // model is ready: if the listener stays closed, that is reported honestly.
        // 192.0.2.1 is TEST-NET-1: unroutable by reservation, so both probes fail.
        var asleepLog = new StringWriter();
        var asleepRoot = new Ensemble.Ensemble(asleepLog,
            dataDirectory: Path.Combine(temp, "wake-asleep"));
        var asleepAgent = new Agent(asleepRoot, new Configuration(Vendor.OpenAI, "gpt-4o",
            "http://192.0.2.1:8110/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"),
            name: "wake-asleep");
        string failure = "";
        try
        {
            await Engine.WakeUnavailableProviderAsync(asleepAgent,
                asleepAgent.Configuration(), CancellationToken.None, TimeSpan.FromSeconds(2));
        }
        catch (InvalidOperationException error)
        {
            failure = error.Message;
        }
        string logged = asleepLog.ToString();
        bool attempted = logged.Contains("Wake-on-LAN", StringComparison.Ordinal);
        Check(attempted || failure.Contains("no local broadcast route", StringComparison.Ordinal),
            "a silent provider was neither woken nor reported unreachable: " + failure);
        if (attempted)
        {
            Check(failure.Contains("unreachable after Wake-on-LAN", StringComparison.Ordinal),
                "a sent wake did not report the listener staying closed: " + failure);
        }

        // A host that takes 90 seconds to boot must not look like a stalled model
        // request. The indicator is checked against the model branch on purpose:
        // both are true during a wake, and the wake is the more specific answer.
        var visibleRoot = new Ensemble.Ensemble(TextWriter.Null,
            dataDirectory: Path.Combine(temp, "wake-visible"));
        var visibleAgent = new Agent(visibleRoot, new Configuration(Vendor.OpenAI, "gpt-4o",
            "http://192.0.2.1:8110/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"), name: "wake-visible");
        visibleAgent.BeginEngineSend();
        visibleAgent.BeginProviderWake("llm");
        WorkActivity? waking = (await visibleAgent.WorkStatusAsync()).Work;
        Check(waking?.Phase == WorkPhase.WakingProvider && waking?.Detail == "llm",
            "a waking provider was not reported to the browser, or was hidden behind the "
                + "model branch: got '" + waking?.Phase + "/" + waking?.Detail + "'");
        visibleAgent.EndProviderWake();
        string modelDetail = (await visibleAgent.WorkStatusAsync()).Work?.Detail ?? "";
        Check(modelDetail.Contains("model", StringComparison.OrdinalIgnoreCase),
            "releasing the wake did not restore the model wait: got '" + modelDetail + "'");
        visibleAgent.EndEngineSend();

        // The release must happen on every exit, including the give-up path: an
        // indicator left set would claim a wake forever.
        var clearedLog = new StringWriter();
        var clearedAgent = new Agent(new Ensemble.Ensemble(clearedLog,
            dataDirectory: Path.Combine(temp, "wake-cleared")), new Configuration(Vendor.OpenAI,
            "gpt-4o", "http://192.0.2.1:8110/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"),
            name: "wake-cleared");
        try
        {
            await Engine.WakeUnavailableProviderAsync(clearedAgent, clearedAgent.Configuration(),
                CancellationToken.None, TimeSpan.FromSeconds(2));
        }
        catch (InvalidOperationException)
        {
            // Expected: the listener stays closed. What matters is what it left behind.
        }
        AgentWorkStatus settled = await clearedAgent.WorkStatusAsync();
        Check(settled.Work?.Phase != WorkPhase.WakingProvider,
            "the wake indicator outlived the wait: " + settled.Work?.Phase);
        // The positive half, observed on a real wake rather than a synthetic flag:
        // where this machine has a route, the packet is sent and the wait is long,
        // so the indicator must be visible for its duration. Where there is no
        // route the send is refused before any wait, which the send test above
        // already pins; asserting here anyway would only test the sandbox.
        var liveLog = new StringWriter();
        var liveAgent = new Agent(new Ensemble.Ensemble(liveLog,
            dataDirectory: Path.Combine(temp, "wake-live")), new Configuration(Vendor.OpenAI,
            "gpt-4o", "http://192.0.2.1:8110/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"),
            name: "wake-live");
        Task wake = Engine.WakeUnavailableProviderAsync(liveAgent, liveAgent.Configuration(),
            CancellationToken.None, TimeSpan.FromSeconds(6));
        string liveDetail = "";
        DateTimeOffset watchUntil = DateTimeOffset.UtcNow + TimeSpan.FromSeconds(4);
        while (DateTimeOffset.UtcNow < watchUntil && !wake.IsCompleted)
        {
            WorkActivity? live = (await liveAgent.WorkStatusAsync()).Work;
            liveDetail = live?.Phase + "/" + live?.Detail;
            if (live?.Phase == WorkPhase.WakingProvider)
            {
                break;
            }
            await Task.Delay(50);
        }
        try
        {
            await wake;
        }
        catch (InvalidOperationException)
        {
            // The listener never opens; the assertion below is about visibility.
        }
        bool sent = liveLog.ToString().Contains("Wake-on-LAN", StringComparison.Ordinal);
        Check(!sent || liveDetail.Contains("WakingProvider/192.0.2.1", StringComparison.Ordinal),
            "a wake was sent but never reported to the browser while it waited: got '"
                + liveDetail + "'");

        // The case a powered-down desktop actually presents: the network card still
        // answers a ping while the model server is not running. Treating any ping as
        // "awake" suppressed the wake here and the feature silently did nothing, so
        // this asserts the packet is sent even though the host is reachable by ICMP.
        if (PingAvailable())
        {
            var pingOnlyLog = new StringWriter();
            var pingOnlyRoot = new Ensemble.Ensemble(pingOnlyLog,
                dataDirectory: Path.Combine(temp, "wake-ping-only"));
            var pingOnlyAgent = new Agent(pingOnlyRoot, new Configuration(Vendor.OpenAI, "gpt-4o",
                "http://127.0.0.1:1/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"),
                name: "wake-ping-only");
            Check(await Engine.IcmpReachableAsync("127.0.0.1", CancellationToken.None),
                "loopback should answer ping for this case to mean anything");
            Check(!await Engine.TcpReachableAsync(new ProviderEndpoint("127.0.0.1", 1),
                CancellationToken.None), "a closed port must not read as a live provider");
            try
            {
                await Engine.WakeUnavailableProviderAsync(pingOnlyAgent,
                    pingOnlyAgent.Configuration(), CancellationToken.None, TimeSpan.FromSeconds(2));
            }
            catch (InvalidOperationException)
            {
                // Expected: nothing wakes on loopback, so the listener stays closed.
            }
            string pingOnlyText = pingOnlyLog.ToString();
            Check(pingOnlyText.Contains("Wake-on-LAN", StringComparison.Ordinal),
                "a host that only answers ping suppressed the wake instead of receiving it");
            Check(pingOnlyText.Contains("answers ping", StringComparison.Ordinal),
                "the report did not distinguish a live host with a dead service: " + pingOnlyText);
            pingOnlyRoot.Destroy();
        }

        // The launcher asks the provider to come up before the operator types, but
        // it must not wait for the boot: that would hold the browser off for minutes.
        var startupLog = new StringWriter();
        var startupAgent = new Agent(new Ensemble.Ensemble(startupLog,
            dataDirectory: Path.Combine(temp, "wake-startup")), new Configuration(Vendor.OpenAI,
            "gpt-4o", "http://127.0.0.1:1/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"),
            name: "wake-startup");
        var startupStarted = Stopwatch.GetTimestamp();
        await Engine.WakeProviderAtStartupAsync(startupAgent, startupAgent.Configuration(),
            CancellationToken.None);
        TimeSpan startupElapsed = Stopwatch.GetElapsedTime(startupStarted);
        Check(startupElapsed < TimeSpan.FromSeconds(10),
            "the launcher waited for a boot it was told not to await: " + startupElapsed.TotalSeconds);
        Check(startupLog.ToString().Contains("sent Wake-on-LAN packet", StringComparison.Ordinal),
            "starting the GUI did not ask a suspended provider to come up");

        // One datagram to a barely-powered card is not a reliable request. While the
        // listener stays closed the packet goes out again, so a lost first attempt
        // does not cost the operator the whole boot window.
        var repeatLog = new StringWriter();
        var repeatAgent = new Agent(new Ensemble.Ensemble(repeatLog,
            dataDirectory: Path.Combine(temp, "wake-repeat")), new Configuration(Vendor.OpenAI,
            "gpt-4o", "http://192.0.2.1:8110/v1", wakeOnLanMac: "aa:bb:cc:dd:ee:ff"),
            name: "wake-repeat");
        try
        {
            await Engine.WakeUnavailableProviderAsync(repeatAgent, repeatAgent.Configuration(),
                CancellationToken.None, TimeSpan.FromSeconds(3), TimeSpan.FromMilliseconds(500));
        }
        catch (InvalidOperationException)
        {
            // Expected: nothing answers TEST-NET-1.
        }
        int sends = repeatLog.ToString().Split("Wake-on-LAN packet",
            StringSplitOptions.None).Length - 1;
        Check(sends >= 2,
            "the magic packet was sent once and never repeated across a 3s wait: " + sends);

        // Cancellation must not be swallowed by the wake path.
        using var cancelled = new CancellationTokenSource(TimeSpan.FromMilliseconds(50));
        bool observed = false;
        try
        {
            await Engine.WakeUnavailableProviderAsync(asleepAgent,
                asleepAgent.Configuration(), cancelled.Token, TimeSpan.FromSeconds(30));
        }
        catch (OperationCanceledException)
        {
            observed = true;
        }
        Check(observed, "the wake path ignored cancellation");
        Check(!Engine.TryProviderEndpoint("ftp://llm:8110", out _), "a non-HTTP scheme resolved");
    }

    static int FreeTcpPort()
    {
        var probe = new TcpListener(IPAddress.Loopback, 0);
        probe.Start();
        int port = ((IPEndPoint)probe.LocalEndpoint).Port;
        probe.Stop();
        return port;
    }

    static bool PingAvailable()
    {
        try
        {
            using var process = Process.Start(new ProcessStartInfo(
                OperatingSystem.IsWindows() ? "ping" : "/bin/ping",
                OperatingSystem.IsWindows() ? "-n 0" : "-c 0")
            { RedirectStandardOutput = true, RedirectStandardError = true, UseShellExecute = false });
            return process is not null && process.WaitForExit(3000);
        }
        catch (Exception)
        {
            return false;
        }
    }
}
