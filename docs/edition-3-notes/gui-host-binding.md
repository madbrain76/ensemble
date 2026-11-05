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

## A bound listener that could not be used (found in daily use)

`--host 10.0.4.30` bound the socket, issued its Digest challenge, and then
answered a correctly authenticated LAN request with **400**. The GUI was
reachable in exactly the way that makes a lock look working and is useless:
unauthenticated requests were refused, authenticated ones were rejected too.

- `RequestOrigin` began with `if (remote is null || !IPAddress.IsLoopback(remote))
  return null;`. Upstream's model had two trusted clients — local access, and an
  HTTPS reverse proxy that reaches the listener over loopback and is named by
  `--public-origin`. A plain-HTTP LAN client is neither, so it fell outside the
  gate. The bind port opened the socket without ever widening the origin rule to
  match, and Digest auth runs before the route group, so the rejection surfaced
  as a 400 after a successful login rather than as an obvious refusal.
- `HostBindingTests` did contain "a LAN bind serves the page and still
  authenticates", and it passed. It checked off-loopback only for the anonymous
  401, then built its Digest answer from a challenge fetched at
  `http://127.0.0.1:{port}/` and sent it to `http://127.0.0.1:{port}/`. The
  authenticated leg never left loopback. Asserting the refusal without asserting
  the admission is how a listener that cannot serve its own purpose looked green;
  the earlier claim in this note that the test "confirms a non-loopback interface
  answers" was therefore overstated.
- The gate now trusts the transport rather than a new configuration input: for a
  non-loopback remote the `Host` header must name `Connection.LocalIpAddress`,
  the address the connection actually arrived on, with a matching port and plain
  http. No parameter was threaded through the four `RequestOrigin` call sites,
  and nothing client-supplied decides the comparison.
- The rebinding property survives: a hostile page whose name resolves to this
  machine sends its own hostname, which cannot equal the arrival address, so it
  is still refused. `--publicOrigin` keeps its proxy path unchanged, loopback
  behaviour is byte-for-byte the same, and a loopback-bound listener can never
  take the new branch because its local address is loopback.
- Tests now authenticate over the LAN address itself, and assert refusal for a
  `Host` that is not the arrival address both off-loopback and from loopback.
  Proven live: reinstating the loopback-only gate fails with "correct credentials
  were refused from an off-loopback client". `agent/Tests` stays at 341 passed /
  4 pre-existing macOS PTY failures, so the shared origin rule regressed nothing.
