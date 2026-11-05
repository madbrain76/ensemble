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
        string? publicOrigin)
    {
        var request = context.Request;
        if (context.Connection.RemoteIpAddress is not { } remote) return null;
        if (IPAddress.IsLoopback(remote))
        {
            // Both native local access and the authenticated proxy arrive over
            // loopback. Never use X-Forwarded-Host/Proto to widen this boundary.
            if (request.Scheme == "http" && request.Host.Port == port
                && (request.Host.Host is "localhost" or "127.0.0.1" or "::1" or "[::1]"))
            {
                return $"http://{request.Host}";
            }
            if (publicOrigin is null) return null;
            var configured = new Uri(publicOrigin, UriKind.Absolute);
            if (!string.Equals(request.Host.Host, configured.Host, StringComparison.OrdinalIgnoreCase)
                || (request.Host.Port ?? 443) != configured.Port) return null;
            return publicOrigin;
        }
        // A client off-loopback is served only by a listener deliberately bound to a
        // LAN address, which the host refuses to open without Digest credentials. The
        // Host header then names the origin the request carries, exactly as a browser
        // sets it: it is compared against the page's own Origin header by MatchesOrigin,
        // so a cross-origin page still fails. Verifying this name against the arrival
        // address was tried and removed; see docs/edition-3-notes/gui-host-binding.md.
        if (request.Scheme != "http" || request.Host.Port != port) return null;
        return $"http://{request.Host}";
    }

    private static partial bool MatchesOrigin(HttpRequest request, string expected, bool required)
    {
        var origins = request.Headers.Origin;
        if (origins.Count == 0) return !required;
        return origins.Count == 1 && string.Equals(origins[0], expected, StringComparison.Ordinal);
    }
}
