using System.Net;

namespace EnsembleGui;

// Private LAN protection for the GUI. Digest with SHA-256 keeps the password off
// the wire and off disk: the file holds only HA1 = SHA-256(username:realm:password),
// so a leaked file still needs an offline crack before it can log anyone in, and
// the server never compares a password directly. Authentication is optional and
// presence-driven: an existing credential file turns it on, its absence leaves a
// loopback listener open exactly as before.
public sealed record GuiDigestCredentials(int Version, string Username, string Realm,
    string Ha1)
{
    // The realm binds a verifier to this application's challenge, so a captured
    // response from another Digest realm cannot be replayed here.
    public const string RealmName = "Ensemble";
    public const int CurrentVersion = 1;
}

public static partial class Gui
{
    // Nonces are self-certifying (issued time plus an HMAC over it), so the server
    // keeps no nonce table: only entries that were actually accepted are remembered,
    // and those only to reject replays.
    public const int DigestNonceLifetimeSeconds = 300;
    public const int DigestReplayLimit = 4096;
    public const int DigestFailureLimit = 1024;
    public const int DigestCredentialsLimit = 4096;

    // Read the private verifier with the checks that make it trustworthy: regular
    // file, owner-only mode, bounded size, exactly the four documented fields.
    public static partial GuiDigestCredentials LoadDigestCredentials(string path);

    // Create or reset the verifier and return the generated password once. It is
    // never stored and never written to a log: a reset is the only recovery.
    public static partial string CreateDigestCredentials(string path, string username,
        bool replace);

    // Decide whether a listener must authenticate. An explicit path is required to
    // exist; otherwise a present file enables auth and a loopback listener may run
    // without it. Any other bind address is refused until credentials exist.
    public static partial string? DigestPath(string? dataDirectory, string? requested,
        bool requestedExplicitly, bool loopback);

    // Guard every route, including the WebSocket upgrade, before any handler runs.
    // Diagnostics receives one short sentence per rejected attempt so an operator can
    // tell a wrong password from a drifted realm; it stays null for embedded hosts and
    // fixtures, which have no terminal to read.
    public static partial void UseDigestAuth(WebApplication host, string path,
        TextWriter? diagnostics = null);

    // "auth-init" and "auth-reset" entry point. Returns the process exit code and
    // writes only the one-time password to output. Diagnostics, when given, receives
    // the verifier path so the command that prints a password also says which file it
    // belongs to: a password generated for one file can never unlock another.
    public static partial int RunAuthCommand(string[] args, TextWriter output,
        TextWriter? diagnostics = null);

    // Reject a credential path that collides with application state, so a
    // misconfigured --gui-auth-file cannot truncate a save, log or settings file.
    public static partial void CheckDigestPath(string path, IEnumerable<string> reserved);

    // Parse the Authorization header into its Digest fields. Rejects malformed
    // syntax, duplicates, embedded control characters and oversized headers rather
    // than guessing at a credential the browser did not clearly send.
    internal static partial Dictionary<string, string>? ParseDigestFields(string? header);

    // Verify response = SHA-256(HA1:nonce:nc:cnonce:qop:SHA-256(method:uri)).
    internal static partial bool DigestResponseMatches(GuiDigestCredentials credentials,
        string method, string uri, Dictionary<string, string> fields);

    // The check behind DigestResponseMatches, naming the first failed condition or
    // null when the answer is correct. One implementation serves both, so the reason
    // an attempt was refused can never drift from the rule that refused it.
    internal static partial string? DigestRejection(GuiDigestCredentials credentials,
        string method, string uri, Dictionary<string, string> fields);

    // Issue and validate the self-certifying nonce.
    internal static partial string DigestNonce(byte[] key, DateTimeOffset issued);
    internal static partial DateTimeOffset? DigestNonceExpiry(byte[] key, string nonce,
        DateTimeOffset now);
}
