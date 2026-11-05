namespace EnsembleGui;

// Value-like launch options, not an owner or another settings store. Relative
// state paths are resolved by the borrowed Agent's Sandbox. Null asset override
// selects assembly resources: deployment does not require a source checkout.
// RunAsync owns the listener on Host; ServeAsync requires the caller to bind
// the same address and Port. BasePath is a literal route prefix, not a route template.
// Account controls follow it, but registration reserves /auth/callback at root.
// The trusted GUI can select existing addresses in the borrowed Agent's Ensemble.
// Selection is per socket and does not transfer ownership. The optional agent_id
// page/socket query chooses a fixed browser relay target for that connection;
// select_agent changes only its conversation view. The caller must keep the
// Ensemble alive until ServeAsync has drained; peers may close independently.
// Host is the bind address: "127.0.0.1" or "localhost" stays on this machine,
// while a specific address or the unspecified address admits the LAN. Only
// literal addresses are accepted, so a hostname cannot silently move the
// listener through DNS. Binding beyond loopback requires Digest credentials.
// PublicOrigin optionally admits one HTTPS origin behind an authenticated local
// reverse proxy. It is not authentication; the proxy must gate every route,
// including WebSockets. Forwarded headers do not grant host/origin authority.
// Trusted peer POST controls use this mounted Agent as the explicit defaults
// source and cannot delete it. Pending creation does not block socket traffic.
// AuthPath names the private Digest credential file. Null defers to the
// conventional <data directory>/gui-auth.json: an existing file turns LAN
// authentication on, and its absence leaves the listener open as before.
// Null TracePath uses Agent-owned diagnostics. An explicit path is a fresh,
// operation-owned aggregate trace for embedding/tests, not the normal host mode.
// DiagnosticsOutput is the operator's error stream, distinct from any projection or
// trace file: it carries short, actionable sentences such as why an authentication
// attempt was refused. Embedded hosts and fixtures leave it null and stay silent.
// AllowedHosts names the extra values a client may put in Host, such as the DNS
// name a router gives this box. They are never looked up: resolving a name the
// client chose is exactly the DNS rebinding hole this check exists to close.
public sealed record GuiOptions(int Port = 8090, string? AssetsDirectory = null,
    string SettingsPath = "settings.json", string? TracePath = null, string BasePath = "",
    string? PublicOrigin = null, string? AuthPath = null, string Host = "127.0.0.1",
    TextWriter? DiagnosticsOutput = null, IReadOnlyList<string>? AllowedHosts = null);

// A refused bind address is a launch decision carrying fixed operator guidance,
// not provider, credential or filesystem text. The host reports unknown
// exceptions by type name only, so this one type is what lets the operator read
// what to do next instead of being told to go look at diagnostics.
public sealed class GuiBindRefusalException : InvalidOperationException
{
    public GuiBindRefusalException(string message) : base(message)
    {
    }
}
