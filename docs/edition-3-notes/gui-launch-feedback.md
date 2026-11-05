# GUI launch feedback: what a silent terminal actually meant

## The report

`run.sh` printed nothing at all and was killed after about three minutes. Nothing
was listening, so it looked like a hung process, and a busy LLM server was the
plausible suspect. Neither diagnosis was right, and both were invisible.

## What it was

Two defects, both about reporting rather than about the provider.

**The launcher never announced itself.** `Examples/Gui` starts Kestrel and then
waits. It printed no address and no port, so a GUI serving correctly was
indistinguishable from one that had died. Reading the saved event log settled it:
the launch ran three minutes and twenty one seconds with no turn in flight, which
is somebody waiting for output that was never going to arrive.

**A busy provider failed the launch outright.** The reasoning-effort check waits
on the provider before the socket is bound. Its own ten second deadline escaped
as an `OperationCanceledException`, which the launcher treats as a launch
mistake, so a slow but living server stopped the GUI from starting at all. Fixed
on `port/9-reasoning-effort-check`; see its note.

## The fix

`GuiOptions` takes an optional `ReadyOutput`. `RunAsync` writes to it only after
`host.StartAsync` has returned, so no address is ever advertised for a socket that
did not bind, and the standalone launcher is the only host that supplies one:
embedded mounts and fixtures stay as quiet as before. A wildcard listener is
expanded to concrete addresses, because a URL ending in `0.0.0.0` names no host a
browser can open, and an authenticated listener says authentication is required,
since a bare 401 otherwise reads as a broken page.

Separately, the launcher prints `Starting: model ... at ...` before any network
work, so the provider wait has a label instead of looking like a hang. Offline
`--render` and `--dump` keep stdout for their own output.

## What was measured

- Slow provider (25 second answer): `Starting:` immediately, `Serving:` at the
  ten second probe budget, `GET /` then returns 200. Before the fix this case
  never bound a socket.
- Fast provider: `Serving:` in under a second.
- LAN bind with credentials: `Serving http://10.0.4.30:8090` plus
  `Digest authentication is required for every request.`, and an unauthenticated
  `GET /` returns 401.
- `agent/GuiTests --host-binding-tests` covers all three shapes, including that a
  refused bind announces nothing. Verified live by breaking one expected URL and
  watching the check fail with the real output.

## Still deferred on this branch

Merging every feature into the daily launcher hit conflicts that are not
mechanical, so these stay out until resolved properly: `port/3` (engine send
path), `port/5` and `port/1` both define `OutputBudgetCeiling`, `OutputCapFields`
and `LargerOutputBudget`, so the two share cap machinery and cannot be unioned;
`port/8` touches `Configuration` in five files at once; `port/11` conflicts in
`app.js`. A blind union produced invalid C#, which is why the branch was reset
and rebuilt from the last verified commit instead.
