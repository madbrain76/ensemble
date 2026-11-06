# Ported feature: Wake-on-LAN for a sleeping local provider

Ported from `go-checkpoint` (`7717894c`, `feature/wake-on-lan`) on 2026-11-03 by
the coordinator, at Bill's request. Branch `port/22-wake-on-lan`, based on
`fix/linux-path-case`; independent of the authentication and host-binding ports.

## Why

Bill runs a model server on a desktop that powers down. Before this port the
first request after a sleep paid a connect timeout and failed, so every session
started by hand-waking the box.

## Shape

- `Configuration.WakeOnLanMac` is validated at construction, so a typo stops the
  launch instead of turning every request into a wake that never arrives. It is
  a local transport hint: never provider input, never History content, and
  preserved across `WithModel` because the address belongs to the machine.
- `Engine.WakeUnavailableProviderAsync` runs on the transport path with the
  request's own cancellation. Empty MAC returns immediately, so ordinary
  providers pay nothing. An open TCP port or an ICMP reply means "already awake":
  no packet, no wait, nothing logged.
- Only when both probes fail is one magic packet sent, then the listener is
  polled for up to 90 seconds. HTTP remains the readiness test: this path never
  claims a model is loadable.
- `agent/Engine/WakeOnLan.cs` builds the 102-byte packet, derives the directed
  broadcast from the interface that would route to the provider, and also tries
  the address itself, since unicast wake is common and it keeps the path testable.

## Findings during porting

- The first probe helpers used `Task.Wait`, which surfaced a refused connect as
  `AggregateException` instead of "not reachable". Both probes are now async with
  a linked cancellation token, so a refused port, a probe timeout and a genuine
  caller cancellation stay three distinct outcomes.
- The first wake test pointed at a dead port on `127.0.0.1`, which is the *awake*
  case: loopback always answers ping, and skipping a wake then is correct. It now
  uses TEST-NET-1 (`192.0.2.1`), unroutable by reservation, so both probes fail
  and the wake path is actually exercised.
- Deviation from the Go branch: the GUI-side MAC entry field was not ported,
  because the Edition 3 GUI protocol has no configuration channel for it yet. The
  operator selects it with `--wake-on-lan aa:bb:cc:dd:ee:ff` or
  `EN_WAKE_ON_LAN`. Recorded here so the gap is not mistaken for an omission.

## Defect found when the flag was actually used

The first revision registered `--wake-on-lan` in both launchers and then never
read it: no `Configuration` in `Examples/Gui` or `Cli` passed `wakeOnLanMac`, so
the flag parsed successfully, was silently discarded, and a sleeping provider
still cost the first turn its full connect timeout. `Configuration` carries the
field and `Engine` reads it, so every test stayed green while the feature could
not be switched on from the command line — exactly the failure mode a flag that
is accepted but ignored produces.

Both launchers now thread the value in, with `EN_WAKE_ON_LAN` as the environment
form, matching how the other route options are supplied. `Configuration`
validates the six-byte address at construction, so a typo stops the launch. The
GUI additionally prints the validator's own sentence for that case: its failure
handler reports unknown exceptions by type name only, and the guidance text is
now a single `Configuration.WakeOnLanGuidance` constant that the validator throws
and the launcher echoes, so the two cannot drift.

Verified against the built binaries, without touching a real LAN: the CLI refuses
`EN_WAKE_ON_LAN=zz:zz` at launch and serves a valid address; the GUI prints
`Wake-on-LAN MAC must be six hex bytes, such as aa:bb:cc:dd:ee:ff` for a malformed
value and starts normally for a valid one. A wake is deliberately not attempted
against `127.0.0.1`, because a refused port there with a pingable host means the
machine is already up — proving the packet path needs an actually powered-down
host, which is Bill's `llm` and remains unverified here.

## The wake never left the machine (found 2026-11-06, by Bill)

Bill reported that after shutting `llm` down the server simply never came up. It
could not have: `SendWakePacket` opened a plain `UdpClient`, and the kernel refuses
to deliver a broadcast on a socket that did not ask for `SO_BROADCAST`. Measured
against his configured provider, with the destination computed from the routing
interface `10.0.4.30/255.255.248.0`:

```
SEND 10.0.7.255:9 -> SocketException: Permission denied   <- the only packet a powered-off card hears
SEND 10.0.5.20:9  -> wrote 102                             <- unicast, dropped at ARP
SHIPPED SendWakePacket(llm) -> True                              <- "sent Wake-on-LAN packet"
```

So the one destination that can wake a sleeping host was refused, the destination
that cannot was counted as delivery, and the operator was told a packet had been
sent. `EnableBroadcast = true` fixes it; the same call then reports `wrote 102` for
both destinations and `/proc/net/snmp` `OutDatagrams` rises by two, so the packets
genuinely leave the host. Go never showed this because its `net` package sets
`SO_BROADCAST` for a broadcast destination itself; the port inherited the semantics
without the option.

**Why the suite was green.** Two independent holes. The delivery assertion sent to
`127.0.0.1`, where `DirectedBroadcast` yields the loopback broadcast — before the
accounting change the refused broadcast was ignored and the successful loopback
unicast counted, so the check passed with the feature dead. And the wake-flow
assertion accepted `attempted || failure contains "no local broadcast route"`, so a
run that emitted nothing at all satisfied it. `SendWakePacket` now reports a
delivered wake only when the broadcast it needed was delivered — a successful
unicast alone is never a wake, since a powered-off card has no live ARP entry and
the kernel drops the frame with no error visible to the caller. That single change
made the existing loopback assertion live: without the socket option it now returns
false and the test fails. `Engine.WakeSocket()` exists so a test can read the option
back without putting a magic packet on the wire.

## Nothing said the provider was being woken (Bill's UI expectation)

Bill expected the packet to go out when he typed into an agent, and it does: the
wake sits on the transport path inside `SendCoreAsync`. What he could not see was
any of it. The wait runs up to `WakeBootTimeout`, ninety seconds, during which the
tree reported `Waiting for the model response` — because `engineSends != 0` and
`ReadWorkActivity` asks that before anything else. A booting desktop and a hung
request rendered identically, which is how a feature that sends nothing at all got
read as a slow server.

`Agent.BeginProviderWake`/`EndProviderWake` now publish a `WorkActivity` in the
existing live-status seam, checked ahead of the model branch and released from a
`finally`, so the row reads `Waking llm` for as long as the machine is coming up.
Two mutants, each caught by its own assertion: dropping the `BeginProviderWake` call
fails the observation of a real wake in flight (`got ''`), and moving the check
behind the model branch fails with `got 'Waiting for the model response'`, the exact
symptom reported. The positive assertion is conditional on a packet having been
sent, so a sandbox with no route cannot fail it; on this machine the route exists
and the mutant proves the branch executes rather than being skipped.

Still unverified here: waking Bill's actual hardware. Both fixes are measured at the
socket and at `/proc/net/snmp`, which is the last point this side of the network
card, and a real wake needs `llm` powered down.
