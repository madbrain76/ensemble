using System.Diagnostics;
using System.Net;
using System.Net.NetworkInformation;
using System.Net.Sockets;

namespace Ensemble;

public sealed partial class Engine
{
    internal static partial bool TryWakeMac(string? value, out byte[]? mac)
    {
        mac = null;
        string text = (value ?? "").Trim().Replace('-', ':');
        if (text.Length == 0)
        {
            return false;
        }
        string[] groups = text.Contains(':') ? text.Split(':') : [text];
        if (groups.Length == 1)
        {
            if (groups[0].Length != 12)
            {
                return false;
            }
            groups = [.. Enumerable.Range(0, 6).Select(index => groups[0].Substring(index * 2, 2))];
        }
        if (groups.Length != 6)
        {
            return false;
        }
        byte[] parsed = new byte[6];
        for (int index = 0; index < 6; index++)
        {
            if (groups[index].Length != 2 || !byte.TryParse(groups[index],
                    System.Globalization.NumberStyles.HexNumber,
                    System.Globalization.CultureInfo.InvariantCulture, out parsed[index]))
            {
                return false;
            }
        }
        mac = parsed;
        return true;
    }

    internal static partial bool TryProviderEndpoint(string baseUrl,
        out ProviderEndpoint? endpoint)
    {
        endpoint = null;
        if (!Uri.TryCreate(baseUrl, UriKind.Absolute, out Uri? uri)
            || (uri.Scheme != Uri.UriSchemeHttp && uri.Scheme != Uri.UriSchemeHttps)
            || string.IsNullOrWhiteSpace(uri.Host))
        {
            return false;
        }
        int port = uri.Port > 0 ? uri.Port
            : uri.Scheme == Uri.UriSchemeHttps ? 443 : 80;
        endpoint = new ProviderEndpoint(uri.Host, port);
        return true;
    }

    internal static partial byte[] WakePacket(byte[] mac)
    {
        byte[] packet = new byte[6 + 16 * mac.Length];
        packet.AsSpan(0, 6).Fill(0xff);
        for (int index = 6; index < packet.Length; index += mac.Length)
        {
            mac.CopyTo(packet.AsSpan(index));
        }
        return packet;
    }

    internal static partial async Task<bool> TcpReachableAsync(ProviderEndpoint endpoint,
        CancellationToken cancellation)
    {
        // One second is generous for a LAN and short enough that an asleep host is
        // not mistaken for a slow one. A refused connect is an answer, not a fault.
        using var deadline = CancellationTokenSource.CreateLinkedTokenSource(cancellation);
        deadline.CancelAfter(WakeProbeTimeout);
        try
        {
            using var probe = new TcpClient();
            await probe.ConnectAsync(endpoint.Host, endpoint.Port, deadline.Token);
            return probe.Connected;
        }
        catch (SocketException)
        {
            return false;
        }
        catch (ArgumentException)
        {
            return false;
        }
        catch (OperationCanceledException) when (!cancellation.IsCancellationRequested)
        {
            return false; // Our own probe deadline elapsed.
        }
    }

    internal static partial async Task<bool> IcmpReachableAsync(string host,
        CancellationToken cancellation)
    {
        try
        {
            // Missing ping, blocked ICMP and timeout all count as an inconclusive
            // probe, never as proof the host is asleep. TCP is tried first so a live
            // provider never waits on a ping.
            ProcessStartInfo start = OperatingSystem.IsWindows()
                ? new ProcessStartInfo("ping", "-n 1 -w 1000 " + Quote(host))
                : new ProcessStartInfo("ping", OperatingSystem.IsMacOS()
                    ? "-c 1 -W 1000 " + Quote(host) : "-c 1 -W 1 " + Quote(host));
            start.RedirectStandardOutput = true;
            start.RedirectStandardError = true;
            start.UseShellExecute = false;
            start.CreateNoWindow = true;
            using var process = Process.Start(start);
            if (process is null)
            {
                return false;
            }
            using var deadline = CancellationTokenSource.CreateLinkedTokenSource(cancellation);
            deadline.CancelAfter(TimeSpan.FromSeconds(2));
            try
            {
                await process.WaitForExitAsync(deadline.Token);
            }
            catch (OperationCanceledException) when (!cancellation.IsCancellationRequested)
            {
                try
                {
                    process.Kill(entireProcessTree: true);
                }
                catch (InvalidOperationException)
                {
                    // It exited between the check and the kill.
                }
                return false;
            }
            return process.ExitCode == 0;
        }
        catch (Exception) when (cancellation.IsCancellationRequested)
        {
            throw new OperationCanceledException(cancellation);
        }
        catch (Exception)
        {
            return false;
        }
    }

    internal static partial UdpClient WakeSocket()
    {
        UdpClient client = new(AddressFamily.InterNetwork);
        // Without this the kernel rejects every broadcast with EACCES, the directed
        // broadcast is silently skipped, and the unicast that follows cannot reach a
        // powered-off card. The feature then reports success and wakes nothing.
        client.EnableBroadcast = true;
        return client;
    }

    internal static partial bool SendWakePacket(string host, byte[] mac, int port,
        CancellationToken cancellation)
    {
        IPAddress[] addresses;
        try
        {
            addresses = Dns.GetHostAddresses(host);
        }
        catch (Exception) when (!cancellation.IsCancellationRequested)
        {
            return false;
        }
        var destinations = new List<(IPEndPoint Point, bool Broadcast)>();
        foreach (IPAddress address in addresses.Where(value =>
                     value.AddressFamily == AddressFamily.InterNetwork))
        {
            // A sleeping machine answers nothing, so the packet goes to the subnet's
            // directed broadcast; the address itself is also tried, because many
            // switches and NICs accept unicast wake, and it keeps the path testable.
            IPAddress? broadcast = DirectedBroadcast(address);
            foreach (IPAddress candidate in new[] { broadcast, address }
                         .Where(value => value is not null).Select(value => value!)
                         .Distinct())
            {
                destinations.Add((new IPEndPoint(candidate, port),
                    candidate.Equals(broadcast)));
            }
        }
        bool unicastDelivered = false;
        bool broadcastDelivered = false;
        bool broadcastTried = false;
        using UdpClient client = WakeSocket();
        foreach ((IPEndPoint destination, bool isBroadcast) in destinations)
        {
            cancellation.ThrowIfCancellationRequested();
            broadcastTried |= isBroadcast;
            try
            {
                if (client.Send(WakePacket(mac), destination) > 0)
                {
                    if (isBroadcast)
                    {
                        broadcastDelivered = true;
                    }
                    else
                    {
                        unicastDelivered = true;
                    }
                }
            }
            catch (SocketException)
            {
                // This route is unusable; another interface may still reach the host.
            }
        }
        // Only the broadcast can reach a powered-off network card, so a successful
        // unicast alone is never a delivered wake: an asleep host has no live ARP
        // entry, and the kernel drops the frame without the caller learning. Where a
        // broadcast route existed, its failure is the answer.
        return broadcastTried ? broadcastDelivered : unicastDelivered;
    }

    internal static partial async Task WakeUnavailableProviderAsync(Agent parent,
        Configuration config, CancellationToken cancellation, TimeSpan? bootTimeout)
    {
        if (config.WakeOnLanMac.Length == 0)
        {
            return;
        }
        cancellation.ThrowIfCancellationRequested();
        if (!TryProviderEndpoint(config.BaseUrl, out ProviderEndpoint? endpoint))
        {
            throw new InvalidOperationException(
                "Wake-on-LAN requires an absolute http(s) provider base URL");
        }
        // Only the provider's own port proves the service is listening. A sleeping
        // machine's network card commonly still answers a ping -- that is precisely
        // the state this feature exists for -- so letting ICMP suppress the wake
        // made the feature do nothing on a real powered-down host, while every test
        // stayed green because loopback never reproduces it. Ping is kept for the
        // operator's sentence: "host up, service down" and "host silent" are
        // different faults even though both need the same packet.
        if (await TcpReachableAsync(endpoint!, cancellation))
        {
            parent.Parent().Log("provider " + endpoint!.Address + " is listening; no wake needed");
            return;
        }
        bool answersPing = await IcmpReachableAsync(endpoint!.Host, cancellation);
        if (!TryWakeMac(config.WakeOnLanMac, out byte[]? mac))
        {
            throw new InvalidOperationException("invalid Wake-on-LAN MAC");
        }
        // From here the wait belongs to a machine rather than a model, and a
        // booting desktop can hold the turn for the whole boot timeout. Reported
        // before the send so a refused broadcast is also covered by the release.
        parent.BeginProviderWake(endpoint!.Host);
        try
        {
            if (!SendWakePacket(endpoint!.Host, mac!, WakePort, cancellation))
            {
                throw new InvalidOperationException("no local broadcast route to provider");
            }
            parent.Parent().Log(answersPing
                ? "provider port is closed but the host answers ping; sent Wake-on-LAN packet"
                : "provider endpoint is silent; sent Wake-on-LAN packet");
            // Only the listener is awaited here. Loading and testing a model remains
            // the HTTP request that follows, so a booting host is never declared ready.
            DateTimeOffset deadline = DateTimeOffset.UtcNow + (bootTimeout ?? WakeBootTimeout);
            while (DateTimeOffset.UtcNow < deadline)
            {
                cancellation.ThrowIfCancellationRequested();
                if (await TcpReachableAsync(endpoint!, cancellation))
                {
                    return;
                }
                await Task.Delay(WakeProbeTimeout, cancellation);
            }
            throw new InvalidOperationException(
                "provider remained unreachable after Wake-on-LAN");
        }
        finally
        {
            parent.EndProviderWake();
        }
    }

    private static IPAddress? DirectedBroadcast(IPAddress destination)
    {
        foreach (NetworkInterface adapter in NetworkInterface.GetAllNetworkInterfaces())
        {
            if (adapter.OperationalStatus != OperationalStatus.Up)
            {
                continue;
            }
            foreach (UnicastIPAddressInformation unicast in adapter.GetIPProperties()
                         .UnicastAddresses)
            {
                if (unicast.Address.AddressFamily != AddressFamily.InterNetwork
                    || unicast.IPv4Mask is null || !SameSubnet(unicast.Address,
                        unicast.IPv4Mask, destination))
                {
                    continue;
                }
                byte[] local = unicast.Address.GetAddressBytes();
                byte[] mask = unicast.IPv4Mask.GetAddressBytes();
                return new IPAddress([.. local.Zip(mask,
                    (value, bits) => (byte)(value | ~bits))]);
            }
        }
        return null;
    }

    private static bool SameSubnet(IPAddress local, IPAddress mask, IPAddress other)
    {
        byte[] left = local.GetAddressBytes(), right = other.GetAddressBytes(),
            bits = mask.GetAddressBytes();
        for (int index = 0; index < 4; index++)
        {
            if ((left[index] & bits[index]) != (right[index] & bits[index]))
            {
                return false;
            }
        }
        return true;
    }

    private static string Quote(string host) => "\"" + host.Replace("\"", "") + "\"";
}
