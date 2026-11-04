using System.Net;

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

public sealed partial class Engine
{
    public const int WakePort = 9;
    public static readonly TimeSpan WakeProbeTimeout = TimeSpan.FromSeconds(1);
    public static readonly TimeSpan WakeBootTimeout = TimeSpan.FromSeconds(90);

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

    // Probe, wake, then poll for the listener. A configuration without a MAC does
    // nothing at all, so ordinary providers pay nothing for this feature.
    internal static partial Task WakeUnavailableProviderAsync(Agent parent,
        Configuration config, CancellationToken cancellation,
        TimeSpan? bootTimeout = null);
}
