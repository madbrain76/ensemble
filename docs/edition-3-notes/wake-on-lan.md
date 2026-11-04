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
