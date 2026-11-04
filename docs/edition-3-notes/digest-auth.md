# Ported feature: SHA-256 Digest authentication for the LAN GUI

Ported from `go-checkpoint` (items `15084c1e`, `89ef8e11`) to the Edition 3 C#
tree by the coordinator on 2026-11-03, at Bill's request. Branch
`port/17-digest-auth`, based on `fix/linux-path-case` because upstream `main`
alone does not build on a case-sensitive filesystem.

## Why

The GUI binds loopback and is otherwise open to any process on the host, and
`--public-origin` admits one HTTPS origin behind a proxy without authenticating
it. Bill's requirement was to stop a rogue LAN program that merely connects and
starts issuing controls. Transport encryption stayed out of scope by his
decision, so the threat model is an on-path attacker who cannot manipulate
traffic: Digest with SHA-256 keeps the password off the wire and off disk, and
slows repeated guessing.

## Shape

- `agent/gui/Common/DigestAuth.cs` declares the vocabulary: the
  `GuiDigestCredentials` record (version, username, realm, HA1 only), the
  lifetime and table limits, and the partial members.
- `agent/gui/Server/DigestAuth.cs` implements them: strict verifier loading,
  creation/reset, the RFC 7616 field parser, the SHA-256 answer check, the
  self-certifying HMAC nonce, the per-address slowdown, and the middleware.
- `GuiOptions.AuthPath` carries a decided path. `ServeAsync` only mounts it;
  `Examples/Gui` decides policy from the data directory, because the reusable
  GUI may be embedded by a caller whose Agent has no data directory at all.
- `auth-init` and `auth-reset` are offline subcommands of the GUI entry point.
  They compose no Agent, open no port and read no provider credential.

Presence is the switch: an existing `<data-dir>/gui-auth.json` turns
authentication on, and its absence leaves a loopback listener exactly as before.
A non-loopback bind will later be refused without credentials by the host-binding
port (feature 18).

## Findings during porting

- The first parser version trimmed whitespace only after a comma, not at the top
  of each iteration, so it rejected every real header with more than one field.
  Caught by `DigestAuthTests.Parser`, before any browser was involved.
- The port initially dropped the Go rule that the nonce counter is nonzero;
  `nc=00000000` was accepted, which would have pinned one replay-table key.
- Reading `Agent.DataDirectory()` unconditionally from `ServeAsync` broke six
  existing GUI/Agent tests that construct an Agent without a data directory.
  The policy moved to the composition root, which is also where the bind-address
  decision belongs.
- `agent/Tests` reports 339 passed / 6 failed both with and without this change:
  four are the `agent/Jobs/Pty.cs` macOS gate, and the other two drift between
  `Chapter 11 automatic collection`, `Chapter 21 fixed child change waits` and
  `auth login stale registration` across runs. They are pre-existing timing
  sensitivity in this environment, not regressions from this port.
- The credential mode check is Unix-only, matching where the concept exists;
  `.NET` exposes no portable uid read, so ownership is not re-verified beyond
  the 0600 requirement.
