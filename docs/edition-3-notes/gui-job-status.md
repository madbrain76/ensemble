# Header status: idle actor, running Job

## Symptom

`local/daily` showed `Connected idle` in the top-left while the agent was working on a
long build, and gave no signal when the work actually ended.

## What the label actually is

`index.html` renders two spans: `#connection` (socket) and `#state` (actor turn state).
`#state` is `acknowledgedState` (`app.js`), updated only by `agent_state`/`state_changed`
frames, which the server derives from `History`'s `Phase`:
`idle | inputpending | inflight | toolspending` (`agent/History/History.cs:286-361`).

That is a faithful picture of the **actor**. It was never a picture of the **work**.

## Mechanism, measured before changing anything

Firefox against the real `Examples/Gui` launcher, a fake provider that returns a
`run_command` call for `sleep 30`, headers sampled every 150 ms:

```
t+0.2s  inflight        model thinking
t+2.1s  toolspending    tool being observed
t+5.2s  inflight        model thinking again
t+9.2s  idle            turn ended  <-- build still running
```

The tool result the model received is the reason:

```
job 1, status: running, bytes: 0, truncated: false, output_path: .../agents/<name>/jobs/1
```

`run_command` returns a **running report** once its observation delay (`JobLimits.Delay = 3`)
expires, and the Job keeps executing under `Jobs` (`ToolKind.RunCommand` →
`Job.RunCommandAsync`; the model collects it later with `wait_for_job`). The turn ends, the
actor is genuinely idle, and the label is literally true about the only thing it watches.
A 30-minute kernel build would have read `idle` for 30 minutes with no end marker.

Two false leads were eliminated by measurement rather than reading:

- **Frame loss.** `Queue` uses `TryWrite` and aborts the socket when full
  (`agent/gui/Server/Host.cs`), so frames are never silently dropped; the socket stayed
  open and every expected `state_changed` arrived, captured at the browser.
- **The sidebar row lagging** was an artifact of my own probe reading the wrong row on a
  200 ms poll, not the header.

## Fix

`RefreshTree` already visits every Agent every 200 ms and is the only thing that sees live
activity, so it now publishes a running-Job count and the header follows it:

- server: `["jobs"] = RunningJobs(node.Agent)` counting `Jobs().Children()` whose
  `Snapshot().Status == "running"`. The count is read through the Agent's own `Jobs`
  child; no browser counter is invented, per the file's standing rule.
- `Agent.Jobs()` throws `ObjectDisposedException` for an Agent that never ran a command,
  so the count is guarded locally: that is **zero**, not a row to drop from the tree.
- client: `1 job running` / `3 jobs running` when no turn is in flight; `inflight`,
  `toolspending` and `inputpending` stay more specific than a count, and `paused` does not
  hide work that is still running.
- `handle()` now repaints the header on `agent_tree`, which is what carries the count.

After the change, same harness:

```
t+9.2s  "1 job running"   turn ended, process still live
t+32.2s idle              the process exited
```

## Tests

- `agent/GuiTests/JobStatusTests.cs`: real Agent, real `run_command`, real GUI socket.
  It waits for the specific `jobs == 1 && state == "idle"` combination — the exact window
  that was wrong — then for `jobs == 0`. Replacing the count with a constant makes it hang
  until its deadline, so it is not decorative.
  An earlier version asserted `state == "idle"` on the first frame that showed one job and
  raced: the tool is still in flight when the Job is created. Waiting for the combination
  is what made it deterministic.
- `agent/GuiTests/job-status.cjs`: Firefox renders the header for 0/1/3 jobs, another
  Agent's jobs (must not leak), a missing row, paused, inflight, toolspending, and the
  `agent_tree` repaint. Two mutants, two distinct failures: neutralising the label
  expression fails the count cases, removing the repaint fails the frame case.

## Limits and follow-ups

- The count is per-Agent and covers `Jobs` children only. Work owned elsewhere — a child
  Agent spawned with `spawn_subagent`, an MCP call in flight — is still not represented.
- Nothing wakes the model when a Job finishes; collecting the result remains the model's
  `wait_for_job`. This change reports the state, it does not schedule anything.
- `Agent.Jobs()` throwing `ObjectDisposedException` for "never had any" conflates an absent
  subsystem with a destroyed one. The GUI has to know that to avoid dropping rows; a
  dedicated query on Agent would be the cleaner seam if this grows.
- The vocabulary is still internal jargon in the header (`inflight`, `toolspending`).
  Humanising it is a separate display decision, deliberately not mixed into this fix.
