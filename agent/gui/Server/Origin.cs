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

    private static partial string? RequestOrigin(HttpContext context, int port, string? publicOrigin)
    {
        var request = context.Request;
        if (context.Connection.RemoteIpAddress is not { } remote) return null;
        if (IPAddress.IsLoopback(remote))
        {
            // Both native local access and the authenticated proxy arrive over
            // loopback. Never use X-Forwarded-Host/Proto to widen this boundary.
            if (request.Scheme == "http" && request.Host.Port == port
                && request.Host.Host is "localhost" or "127.0.0.1" or "::1" or "[::1]")
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
        // Host header must then name the address the connection actually arrived on.
        // LocalIpAddress comes from the transport rather than the request, so this
        // still refuses a rebinding attack, where a hostile page resolves its own
        // hostname to this machine's address and sends that hostname in Host, while a
        // genuine LAN client names the address it dialed without any cooperation.
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
            ? $"http://{request.Host}"
            : null;
    }

    private static partial bool MatchesOrigin(HttpRequest request, string expected, bool required)
    {
        var origins = request.Headers.Origin;
        if (origins.Count == 0) return !required;
        return origins.Count == 1 && string.Equals(origins[0], expected, StringComparison.Ordinal);
    }
}
