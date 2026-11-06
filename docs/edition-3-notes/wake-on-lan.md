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

**The first version was not visible either.** Publishing a `WorkActivity` is not the
same as showing one: the agent tree renders `phase/tool/detail` only into the row's
`title` attribute, so the wake appeared on hover and nowhere else. Firefox measured
`#state=root busy` with the new activity set — the assertion passed at the data layer
while the operator still saw nothing. `WorkPhase.WakingProvider` is now a phase in
its own right instead of `Waiting` carrying prose, because the header has to branch
on it and string matching is how a reworded detail would silently blank it, and
`Examples/Gui` shows `Waking <host>` in `#state` ahead of both the state word and the
job count. Measured in Firefox against the real launcher with a TEST-NET provider and
a dummy MAC: `#state=Waking 192.0.2.1`, no page errors, and the ten job-status
assertions still pass unchanged.

## Ping was suppressing the wake (found 2026-11-06, Bill's second report)

After the broadcast fix Bill powered `llm` down again and still got nothing, with
only `Connection unavailable; work may still be running.` in the browser and no
wake line anywhere in the log. The remaining defect was the decision itself:

```csharp
if (await TcpReachableAsync(...) || await IcmpReachableAsync(...)) return; // "already awake"
```

A sleeping machine's network card commonly answers a ping -- that is the normal
state of a host waiting for a magic packet, not evidence that it is up. Any NIC or
switch that replies to ICMP while the machine is suspended therefore suppressed the
wake completely, silently, and the operator had no way to distinguish "no wake
needed" from "wake never sent", because the branch logged nothing at all. The Go
reference carries the same rule, so this was faithful to a wrong design.

Only the provider's own port may excuse a wake now. ICMP still runs, but for the
report rather than the decision, because `the host answers ping but the port is
closed` and `the endpoint is silent` are different faults even though both need the
same packet. Every path now says which one it took, including the listening case, so
this question is answerable from the log instead of by experiment.

The regression case is `127.0.0.1` with a closed port: loopback always answers ping,
so it reproduces a powered-down desktop exactly, and it is the first version of this
feature that a test actually exercised -- the earlier loopback case used an *open*
port, which is genuinely the awake case. Reinstating the ping branch fails it with
`a host that only answers ping suppressed the wake instead of receiving it`.

## The sequence Bill actually runs

Suspend llm, start `run.sh`, type into an agent. Nothing in the code supported that
order: the packet was sent only from inside a turn, so step 2 did nothing at all and
step 3 began the wake then, making the first prompt wait out the whole boot. The GUI
now asks the provider to come up while composing the root Agent, before it binds, and
deliberately does not wait -- a suspended desktop must not keep the browser from
opening. The first turn still performs the patient wait itself.

Two constants were wrong for a GPU box. `WakeBootTimeout` was ninety seconds, but the
inference port exists only after the OS boots, the driver loads and CUDA graphs are
captured, which is routinely longer; a wait that expires while the machine is still
coming up reports the provider as unreachable when the truth is that it is slow. It is
now four minutes. And one datagram was the entire request: a broadcast lost on the way
to a barely-powered card left the operator staring at a wait nothing would end early,
until the next prompt happened to ask again. The packet is now repeated every thirty
seconds while the listener stays closed, each repeat logged, so `sent once at boot and
never heard from` is no longer a possible outcome.

`WakeProviderAtStartupAsync` is public because `Examples/Gui` is a separate assembly;
its failure is caught and printed rather than allowed to stop the launch, since a
provider that cannot be woken should still leave a readable transcript. Covered by two
cases: startup returns in well under the boot window and logs the send, and a three
second wait against an unreachable endpoint records more than one send. Removing the
retransmission fails the second with `the magic packet was sent once and never repeated
across a 3s wait: 1`.
