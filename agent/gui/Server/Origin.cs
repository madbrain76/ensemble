using System.Net;

namespace EnsembleGui;

public static partial class Gui
{
    private static partial string? ValidatePublicOrigin(string? origin)
    {
        if (origin is null) return null;
        if (!Uri.TryCreate(origin, UriKind.Absolute, out var uri)
            || uri.Scheme != "https" || uri.HostNameType != UriHostNameType.Dns
            || uri.IsLoopback || uri.UserInfo.Length != 0 || uri.AbsolutePath != "/"
            || uri.Query.Length != 0 || uri.Fragment.Length != 0
            || origin.Contains('\\') || origin.Any(char.IsWhiteSpace))
        {
            throw new ArgumentException("Public origin must be one HTTPS DNS origin, without credentials, path, query or fragment.");
        }
        return uri.GetLeftPart(UriPartial.Authority);
    }

    private static partial string? RequestOrigin(HttpContext context, int port,
        string? publicOrigin, IReadOnlyList<string> allowedHosts)
    {
        var request = context.Request;
        if (context.Connection.RemoteIpAddress is not { } remote) return null;
        if (IPAddress.IsLoopback(remote))
        {
            // Both native local access and the authenticated proxy arrive over
            // loopback. Never use X-Forwarded-Host/Proto to widen this boundary.
            if (request.Scheme == "http" && request.Host.Port == port
                && (request.Host.Host is "localhost" or "127.0.0.1" or "::1" or "[::1]"
                    || NamesThisMachine(request.Host.Host)
                    || NamesAsPermitted(request.Host.Host, allowedHosts)))
            {
                return $"http://{request.Host}";
            }
            if (publicOrigin is null) return null;
            var configured = new Uri(publicOrigin, UriKind.Absolute);
            if (!string.Equals(request.Host.Host, configured.Host, StringComparison.OrdinalIgnoreCase)
                || (request.Host.Port ?? 443) != configured.Port) return null;
            return publicOrigin;
        }
        // A client off-loopback is only served by a listener deliberately bound to a
        // LAN address, which the host refuses to open without Digest credentials. The
        // Host header must then name this machine: either the address the connection
        // actually arrived on, or a name the local operating system knows for itself.
        // LocalIpAddress comes from the transport and the name list comes from the
        // host, never from the request, so this still refuses a rebinding attack,
        // where a hostile page resolves its own hostname to this machine's address
        // and sends that hostname in Host. A genuine client can only name this
        // machine by a name its own resolver already maps here.
        if (request.Scheme != "http" || request.Host.Port != port) return null;
        if (context.Connection.LocalIpAddress is not { } local || IPAddress.IsLoopback(local))
        {
            return null; // No transport-named destination means no claim to trust.
        }
        string named = local.ToString();
        string bracketed = local.AddressFamily == System.Net.Sockets.AddressFamily.InterNetworkV6
            ? "[" + named + "]" : named;
        return string.Equals(request.Host.Host, named, StringComparison.OrdinalIgnoreCase)
            || string.Equals(request.Host.Host, bracketed, StringComparison.OrdinalIgnoreCase)
            || NamesThisMachine(request.Host.Host)
            || NamesAsPermitted(request.Host.Host, allowedHosts)
            ? $"http://{request.Host}"
            : null;
    }

    private static partial bool MatchesOrigin(HttpRequest request, string expected, bool required)
    {
        var origins = request.Headers.Origin;
        if (origins.Count == 0) return !required;
        return origins.Count == 1 && string.Equals(origins[0], expected, StringComparison.Ordinal);
    }

    // The names the local operating system knows for this machine, read once at
    // startup. A dial-in client is allowed to name any of them because none of them
    // is chosen by the client: they come from the host, so a rebinding page cannot
    // borrow one without the user first typing that name themselves.
    private static readonly string[] MachineNames = LocalMachineNames();

    private static string[] LocalMachineNames()
    {
        var names = new List<string>();
        void Add(string? name)
        {
            if (name is { Length: > 0 } value)
            {
                string trimmed = value.TrimEnd('.').ToLowerInvariant();
                if (!names.Contains(trimmed)) names.Add(trimmed);
            }
        }
        try
        {
            Add(Dns.GetHostName());
            Add(Environment.MachineName);
            var identity = System.Net.NetworkInformation.IPGlobalProperties.GetIPGlobalProperties();
            Add(identity.HostName);
            Add(identity.DomainName is { Length: > 0 } domain
                ? identity.HostName + "." + domain : null);
        }
        catch (Exception error) when (error is PlatformNotSupportedException
            or NotSupportedException or InvalidOperationException)
        {
            // No local identity is available, so only the arrival address is trusted.
        }
        return [.. names];
    }

    private static bool NamesThisMachine(string host)
        => Array.IndexOf(MachineNames, host.TrimEnd('.').ToLowerInvariant()) >= 0;

    // Operator-declared names only. The incoming Host is never resolved: a name the
    // client picks could be made to answer with this address by its own DNS, which is
    // the rebinding attack the Host check exists to refuse.
    private static bool NamesAsPermitted(string host, IReadOnlyList<string> allowedHosts)
    {
        string named = host.TrimEnd('.').ToLowerInvariant();
        foreach (string allowed in allowedHosts)
        {
            // Entries are already normalized at launch, so this is one comparison
            // per declared name and allocates nothing on the request path.
            if (string.Equals(allowed, named, StringComparison.Ordinal)) return true;
        }
        return false;
    }

    public static partial IReadOnlyList<string> NormalizeAllowedHosts(IEnumerable<string>? hosts)
    {
        if (hosts is null) return [];
        var names = new List<string>();
        foreach (string entry in hosts)
        {
            string name = (entry ?? "").Trim().TrimEnd('.').ToLowerInvariant();
            if (name.Length == 0) continue;
            if (name.Any(char.IsWhiteSpace) || name.Contains("://") || name.Contains('/')
                || name.Contains(':') || name.Contains('@') || name.Contains('*'))
            {
                throw new ArgumentException(
                    "Allowed host must be a bare name such as devbox.localdomain, with no"
                    + " scheme, port, path, credential or wildcard.");
            }
            if (!names.Contains(name)) names.Add(name);
        }
        return [.. names];
    }

    // One short sentence for the operator's stream. The client still gets the same
    // uniform refusal, so a scanner learns nothing about which rule failed, but the
    // person who just saw a bare 400 in a browser is told what to type instead.
    private static string RefusedOriginMessage(HttpContext context, int port)
    {
        string host = context.Request.Host.HasValue ? context.Request.Host.Value : "(no Host header)";
        string? arrival = context.Connection.LocalIpAddress is { } local
            ? local.AddressFamily == System.Net.Sockets.AddressFamily.InterNetworkV6
                ? "[" + local + "]" : local.ToString()
            : null;
        string suggestion = arrival is null
            ? $"the listener's own address on port {port}"
            : $"http://{arrival}:{context.Request.Host.Port ?? port}/";
        return "gui: refused Host \"" + host + "\"; it names neither the arrival"
            + " address, this machine, nor a declared --allowed-host. Open " + suggestion
            + ", or pass --allowed-host " + host.Split(':')[0] + " if this name should be served.";
    }
}
