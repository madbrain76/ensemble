using System.Net;
using System.Net.Sockets;

namespace Ensemble;

// A local model server on a sleeping desktop costs the first turn a long
// timeout. When the operator supplies the host's MAC, a request that finds the
// endpoint closed sends a Wake-on-LAN magic packet and waits for the listener,
// instead of failing. The wake belongs to the transport path: it uses this
// request's configuration and cancellation, and it never claims model readiness.
public sealed record ProviderEndpoint(string Host, int Port)
{
    public string Address => Host + ":" + Port;
}

public sealed partial class Agent
{
    // A provider that must be woken keeps the Agent busy for as long as the host
    // takes to boot, which on a desktop can be the better part of two minutes.
    // The turn itself already reads as "waiting for the model", so without this
    // the operator sees an ordinary slow request and cannot tell a sleeping
    // machine from a stuck one.
    private int providerWakes;
    private string providerWakeDetail = "";
    internal partial void BeginProviderWake(string host);
    internal partial void EndProviderWake();
}

public sealed partial class Engine
{
    public const int WakePort = 9;
    public static readonly TimeSpan WakeProbeTimeout = TimeSpan.FromSeconds(1);
    // A GPU desktop has to boot an OS, load a driver and capture CUDA graphs before
    // the inference port exists, so ninety seconds routinely ends before the
    // listener does. The wait is patient because the operator is already waiting.
    public static readonly TimeSpan WakeBootTimeout = TimeSpan.FromMinutes(4);
    // A single datagram can be lost on the way to a card that is only barely
    // powered, and nothing else will ask again until the next prompt.
    public static readonly TimeSpan WakeRetransmit = TimeSpan.FromSeconds(30);

    // Parse the six-byte address the magic packet must carry. Only the usual
    // colon/dash/hex-string spellings are accepted, so a typo is a start-up
    // error rather than a wake that silently never arrives.
    internal static partial bool TryWakeMac(string? value, out byte[]? mac);

    // The host and TCP port that prove the provider machine is up. Absence of a
    // port follows the URL scheme, and a non-HTTP scheme is refused rather than
    // guessed at.
    internal static partial bool TryProviderEndpoint(string baseUrl,
        out ProviderEndpoint? endpoint);

    // The 102-byte magic packet: six broadcast bytes then the MAC sixteen times.
    internal static partial byte[] WakePacket(byte[] mac);

    // A TCP connect proves the provider listener; ICMP only proves the machine.
    // Neither proves inference works, which remains the HTTP request's job.
    internal static partial Task<bool> TcpReachableAsync(ProviderEndpoint endpoint,
        CancellationToken cancellation);
    internal static partial Task<bool> IcmpReachableAsync(string host,
        CancellationToken cancellation);

    // Send one magic packet toward the provider's subnet. Delivery to any of the
    // candidate destinations counts; none reachable is a transport failure.
    internal static partial bool SendWakePacket(string host, byte[] mac, int port,
        CancellationToken cancellation);

    // The datagram socket a wake is sent from. A sleeping host only hears a
    // broadcast, and a kernel refuses to deliver one unless the socket says so,
    // so this one option is the difference between a wake and a no-op. It is
    // separate from the send so a test can read it back without emitting traffic.
    internal static partial UdpClient WakeSocket();

    // Probe, wake, then poll for the listener. A configuration without a MAC does
    // nothing at all, so ordinary providers pay nothing for this feature.
    internal static partial Task WakeUnavailableProviderAsync(Agent parent,
        Configuration config, CancellationToken cancellation,
        TimeSpan? bootTimeout = null, TimeSpan? retransmit = null);

    // Send the packet without waiting for the listener, so a launcher can start a
    // suspended provider coming up while the operator is still typing rather than
    // only noticing the machine is down on the first request.
    public static partial Task WakeProviderAtStartupAsync(Agent parent,
        Configuration config, CancellationToken cancellation);
}
