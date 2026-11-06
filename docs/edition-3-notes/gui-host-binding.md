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
  `Host` that names neither this machine nor the arrival address, both
  off-loopback and from loopback. Proven live: reinstating the loopback-only gate
  fails with "correct credentials were refused from an off-loopback client".
  `agent/Tests` stays at 341 passed / 4 pre-existing macOS PTY failures, so the
  shared origin rule regressed nothing.

This transport-trust rule was itself removed later, for the reason recorded in the
next section. Read it before treating any Host-name rule here as current.

## Removed: Host header verification

Bill reported the arrival-address rule as broken and useless and directed its
removal. He was right, and the transcript of his own session is the evidence:

```
./run.sh                                        -> refused Host "devbox.localdomain:8090"
./run.sh --allowed-host devbox.localdomain         -> refused Host "devbox.localdomain:8090"
./run.sh --allowed-host devbox.localdomain:8090    -> refused: not a bare name
```

The flag I invented to rescue the rule could not work: the advice line told him to
strip the port (`pass --allowed-host devbox.localdomain`) while the comparison still
ran against `request.Host.Host` including `:8090`. So the operator was told the
correct fix and the fix was then refused for being correct. A third attempt at
matching this rule is not worth the cost when the rule's own benefit is this thin:

- The listener it guarded is reachable only on a deliberately non-loopback bind,
  and this branch refuses to open such a socket without Digest credentials. An
  attacker who can produce a `Host` for it is already on the LAN and can dial the
  address directly.
- DNS rebinding was the threat the rule addressed. Under Digest that requires the
  victim to type their password into a prompt for the attacker's hostname, and the
  remaining same-origin protection is `MatchesOrigin`, which compares the page's
  `Origin` to the request's own origin and still refuses a cross-origin page.
- Every legitimate client is served by removing it, and no client I can name is
  protected by keeping it.

Removed outright rather than left dormant: `RequestOrigin`'s non-loopback Host
comparison, `MachineNames`/`LocalMachineNames`/`NamesThisMachine`,
`NamesAsPermitted`, `NormalizeAllowedHosts`, `RefusedOriginMessage`, the
`--allowed-host` and `EN_ALLOWED_HOST` inputs, and the `AllowedHosts` field threaded
through `GuiOptions` and four map functions. Net **-203 lines**. The loopback rule
is back to upstream's exactly as it was: `localhost`, `127.0.0.1`, `::1`, or the
`--public-origin` proxy. A non-loopback peer is admitted over http on the right
port and its `Host` becomes the request's origin, which is what the origin check
needs in order to work at all.

What still guards a LAN bind, in order: the bind itself refuses without credentials
(`GuiBindRefusalException`), every request requires a valid Digest answer, and
state-changing requests must carry a matching `Origin`.

`HostBindingTests` now asserts the opposite of what it did: an off-loopback client
naming itself `any-name.example` must be **served**, and an anonymous request with
that same `Host` must still get `401`, so removing the rule cannot become an
authentication bypass. Killing the fix (reinstating the arrival-address
comparison) fails that test with "a LAN client naming itself something other than
the arrival address was refused", so the assertion is live.

## Verified in real Firefox

Headless Firefox against the standalone `Examples/Gui` launcher, `--host 0.0.0.0`,
Digest credentials in a scratch data directory, and a local fake provider. No
`--wake-on-lan`, so nothing could reach the LAN. **12 of 12 checks pass.**

- `auth-init` through `run.sh` printed a password and named the file it wrote.
- An anonymous request carrying a chosen `Host` gets `401`: the credential rule
  survives the removal.
- Page loads over the LAN address **and** when the client dials a name that is
  neither the arrival address nor any name this host knows - the case that returned
  400 - with `#prompt`/`#send` present and scripts loaded.
- A prompt typed into the composer ran for real: send button armed after the
  WebSocket connected, one provider request, reply rendered. Zero page errors.
- Reinstating the arrival-address comparison turns the named case red in the
  browser (`NS_ERROR_NET_EMPTY_RESPONSE`), so this harness would have caught the
  bug Bill hit.

`agent/Tests` is unchanged at 341 passed / 4 pre-existing macOS PTY failures and
`agent/GuiTests` passes, so the shared origin rule regressed nothing else.

One limit worth naming: a browser on the listener's own machine typing
`http://dev:8090/` still gets refused, because upstream's loopback rule admits only
`localhost`/`127.0.0.1`/`::1`. That is upstream's pre-existing rule, not the feature
Bill asked me to delete, and the ready line prints the address to open. Restoring
local browsing by name would mean widening that upstream rule, so it is his call
rather than mine.

## Bill's verdict: removed, and verified served (2026-11-05)

Bill tried the rescue flag three ways -- `--allowed-host devbox.localdomain` ignored
because the incoming `Host` carries `:8090`, then `devbox.localdomain:8090` refused as
"not a bare name" -- and called the feature "broken and useless. remove it." He was
right that it could not be operated: the message told him the correct fix and then
rejected the fix for being correct. What remains is upstream's rule plus Digest
authentication, and `--allowed-host` is gone (the launcher now refuses it by name).

Verified over a genuine non-loopback socket (listener and client both on the machine's
LAN address, so `RemoteIpAddress` is not loopback), with a throwaway verifier rather
than his credential: anonymous over the LAN 401, anonymous dialed by name 401,
authenticated `Host: devbox.localdomain:PORT` **200**, authenticated `Host:
anything.example:PORT` **200**, authenticated by address 200, and a wrong password by
any name 401. Eight of eight. An earlier run of that matrix reported 429s because my
own failed attempts tripped the repeated-failure slowdown; the wrong-password case now
gets a fresh listener, and `qop="auth"` has to be unquoted before hashing.

The loopback limitation below still stands and is still his call: a browser on the
listener's own machine typing `http://dev:8090/` is refused, because upstream admits
only literal loopback names there. Removing that would let any website holding a
domain that resolves to 127.0.0.1 drive a credential-free local GUI, so it is kept
and named rather than deleted quietly.
