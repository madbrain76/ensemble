# Ported feature: configurable GUI listen host

Ported from `go-checkpoint` (`8e970375`, `feature/gui-host`) on 2026-11-03 by the
coordinator, at Bill's request. Branch `port/18-gui-host`, stacked on
`port/17-digest-auth` and therefore on `fix/linux-path-case`.

## Why

`Host.cs` hardcoded `server.Listen(IPAddress.Loopback, options.Port)`, so the GUI
could not be reached from another machine at all, and `--public-origin` only
covers a proxied deployment. Bill wants direct LAN access protected by the
Digest port rather than a proxy he does not run.

## Shape

- `GuiOptions.Host` defaults to `127.0.0.1`, preserving today's behavior.
- `ResolveBindAddress` accepts `localhost`, an IPv4 or an IPv6 literal, and
  refuses empty strings, hostnames, `host:port` spellings, URLs and malformed
  addresses. Names are refused deliberately: the listener's reach must not depend
  on a resolver or the hosts file.
- `RunAsync` refuses any bind beyond loopback when `AuthPath` is empty. That
  invariant is checked where the bind happens, so an embedded caller cannot skip
  the composition root and still open a LAN socket.
- `Examples/Gui` adds `--host`, and feeds the loopback decision into the
  credential-path policy, so `--host 0.0.0.0` without `auth-init` fails at start.

## Verification

`agent/GuiTests/HostBindingTests.cs` covers accepted and refused spellings, the
unauthenticated LAN refusal (asserting the *reason* is authentication, not any
`InvalidOperationException`), that the default loopback listener still starts
without credentials, and that a credentialed `0.0.0.0` bind answers a real
non-loopback interface address while still demanding a valid Digest answer.

A listener that has started is still running, so "did it start" is asserted by
the task not having completed, with its real exception surfaced if it did. The
first draft of this test inverted that and failed for the right reason.

## Defect found while launching it for real

`--host` beyond loopback refuses to start without Digest credentials, and the
refusal was unreadable: `agent/gui/Server/DigestAuth.cs:192` threw a plain
`InvalidOperationException`, and `Examples/Gui` deliberately reports unknown
exceptions by type name only, because exception text can carry credentials, URLs
or private paths. Binding a LAN listener therefore failed as "GUI failed:
InvalidOperationException. Check private diagnostics..." with no hint that
`auth-init` was the missing step. `Host.cs`'s own later gate had the same problem.

`GuiBindRefusalException` now marks the two throws whose text is fixed operator
guidance rather than provider or filesystem data, and `ReportFailure` prints that
one type's message from any depth of the cause chain, because composition joins
the serving and fleet-restore tasks and the refusal usually arrives wrapped. Subclassing `InvalidOperationException` keeps every existing
`is InvalidOperationException` expectation and test intact. Reproduced and verified
against the real binary: `--host 10.0.4.30` without credentials now prints
`GUI authentication required; run auth-init --file "<data dir>/gui-auth.json"
before starting this listener`, while `--host 127.0.0.1` still serves unauthenticated.

The gate belongs to this branch rather than `port/17-digest-auth` because without
a configurable bind address nothing could ever leave loopback, so the unreadable
message was unreachable before now.
