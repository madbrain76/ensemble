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

## A refusal that could not be diagnosed

Reported after the port was in daily use: a browser prompt rejected the correct
username and the operator could not tell why. Two separate files were the cause —
`auth-reset` had been run without `--data-dir`, so it wrote
`~/projects/ensemble/gui-auth.json` while the listener read
`~/.local/share/ensemble-csharp/gui-auth.json`. The printed password was genuine
and simply belonged to a file nothing consulted. Nothing on either side said so.

- Every rejection was one bool. `DigestResponseMatches` returned `false` for a
  wrong password, a wrong username, a drifted realm, a malformed counter and a
  browser that reused a nonce identically, and `AdmitAsync` recorded only the
  address for the slowdown. The uniform 401 is correct as a *wire* behaviour; as
  operator feedback it left nothing to act on.
- `DigestRejection` is now the single implementation and names the first rule the
  answer failed; `DigestResponseMatches` is `DigestRejection(...) is null`, so the
  reason can never drift from the rule that produced it. The constant-time
  comparison stays last.
- The client still receives one uniform 401 whatever the reason. Naming the failed
  rule is a terminal-only privilege: it distinguishes a good username from a bad
  one for whoever can read the log, which on a LAN is the operator, not a scanner.
- Reasons name fields and fixed guidance; client-chosen text is passed through
  `Quoted`, which strips control characters and truncates, so a hostile header
  cannot forge terminal output or bloat the log. No password, digest or verifier
  material is ever written.
- `DigestGuard` takes the verifier path and prints it on the first refusal only.
  One statement of the path is enough to expose a split-brain file, and repeating
  it on every failed attempt would turn a brute-force attempt into log spam.
  `Explain` decides inside the gate and writes outside it: the same lock guards the
  replay and slowdown tables, and no request should wait on a terminal.
- `auth-init`/`auth-reset` now follow `$EN_DATA_DIR` like the launcher, print the
  verifier path they wrote, and warn when some other `gui-auth.json` sits in the
  current directory.
  - CORRECTION (November 5): an earlier version of this note claimed a forgotten
    `--data-dir` "can no longer write a second gui-auth.json that nothing reads".
    That was false and cost another session. The subcommand still falls back to the
    current directory when neither `--data-dir` nor `EN_DATA_DIR` is set, and
    `run.sh` exported the variable only into the process it launches — so a bare
    `./ensemble auth-reset` in a login shell wrote `~/projects/ensemble/gui-auth.json`
    while the listener read `~/.local/share/ensemble-csharp/gui-auth.json`. Both files
    were well formed, both used username and realm `ensemble`, and the printed password
    simply did not work. The stray check was also one-directional: it looked for a file
    in the current directory only when writing somewhere else, so the exact failure that
    happened was the one case it never reported. Now the report names the source
    (`--file`, `--data-dir`, `EN_DATA_DIR` or the current directory) and calls the last
    one out as a guess, and `run.sh` forwards the subcommands so the launcher's data
    directory is used by construction. `agent/start` was deliberately left alone: it
    documents that it keeps the caller's working directory.
- Tests: `RejectionReasons` pins a distinct reason per rule, agrees the bool and
  the reason cannot disagree, and checks sanitization of client text;
  `GuardDiagnostics` drives the real guard and asserts the two reasons, the single
  path statement, and that no password appears. Both proven live by mutating the
  reason strings and the path sentence.

## A reset could not take effect until the listener restarted (November 5)

The verifier was loaded once by `UseDigestAuth` and held in a `readonly` field, so
`auth-reset` against a running GUI was refused forever and the operator was told
`the password does not match the verifier` — accurate about the arithmetic, useless
about the situation, because the process answering was holding a copy of a file the
reset had never touched. This is the same class of failure as the split-brain file
above, and it produces the identical message, which is why it was worth fixing rather
than documenting.

- `DigestGuard` samples the 128-byte verifier at most once a second and reloads when
  the bytes differ. Contents, not `mtime`, because an atomic replace can land inside
  the filesystem's timestamp granularity and the file is small enough to read.
- Reloading rotates the nonce key and clears the replay table. A self-certifying nonce
  minted before the swap can no longer verify, which is what forces a browser that
  cached the previous credential to take a fresh challenge instead of silently
  replaying a dead password. Rotating the key, rather than only swapping credentials,
  is the part that makes `auth-reset` reliable from the browser.
- A verifier that disappears or stops parsing answers 503 and admits nothing. The
  alternative — keep serving on the last known good copy — would let a half-written
  or deliberately removed file leave old access alive, and replacing the file is also
  how a credential is revoked.
- Watching is opt-in (`watchVerifier`). `UseDigestAuth` sets it because the file is its
  authority. Inferring it from `verifierPath` was wrong and the suite said so:
  `GuardDiagnostics` builds a guard with in-memory credentials and names a path only so
  refusals can cite it, and that unrelated file then revoked the guard's own access.
  An embedding host supplies its own credentials; a file it never loaded must not be
  able to lock it out.
- Verified end to end through the real launcher and `run.sh` forwarding, no restart:
  anonymous 401, original password 200, `auth-reset`, retired password 401, new password
  **200**, with `auth: verifier reloaded from <path>; clients holding the previous
  password will be challenged again` in the operator log. `ReloadAndFault` covers the
  same ground plus the 503 and the recovery; each assertion is killed by its own
  mutation (retired password, stale nonce, ignored fault, EN_DATA_DIR source, the
  fallback sentence).

The one-second sample window means a delete can be unnoticed briefly; it is bounded,
single-file, and one read per second against a 128-byte file, so it was preferred to
a `FileSystemWatcher` and its own failure modes. The 503 body carries no
`WWW-Authenticate`, which is correct — nothing can be answered while the verifier is
unavailable — and the next successful request re-challenges normally.
