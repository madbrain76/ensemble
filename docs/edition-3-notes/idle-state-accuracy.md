# An idle agent that was not idle, and a reply that said nothing

Both defects came from Bill's "Yolink local" session and share one root: the
operator's picture of the agent was not the truth. Commit `44e49d50`.

## The empty reply that ended a turn

The provider reported success with no text and no tool call. Nothing claimed
that case:

- `PromisesAction("")` is false, correctly -- the classifier decides whether a
  sentence promises work, and "" promises nothing. `PromissoryTests` pins it as
  an ordinary answer, so the fix belongs at the decision site, not the filter.
- `EndsBeforeAnswer` (`agent/Engine/OutputLimit.Decisions.cs:12`) excludes it on
  purpose: "A model that answered empty ... is not this condition, and retrying
  it would double the cost of every such exchange."

So the turn closed and the page showed nothing at all. That is the stop that
looks like idling.

A silent reply is now corrected inside the existing promissory budget, with its
own message that names the missing half rather than accusing the model of a
promise it never made. When the budget is spent and the model still says
nothing, a durable `ErrorOccurred` records why the turn ended: "Check private
diagnostics" pointing at empty logs is not an explanation. A turn that recovers
is not reported at all -- only a turn that never produced anything.

`SilentReplyTests` drives a real agent against a loopback provider and pins the
request count (`1 + PromissoryLimit`, so a runaway fails), the correction
message, and the notice. Deleting `|| silent` from the condition fails it with
"1 requests for an always-empty model"; deleting the notice fails it with "the
turn closed silently".

## The queue that was invisible

`ReadWorkActivity` consulted wakes, transport, memory pressure, workers,
controls and startup replies, and then returned `Idle`. It never looked at
`requests`, the queue of submissions the actor has admitted but not started. An
operator who believed that idle started a second agent on the same project, and
the two collided.

`WorkPhase` gained `Queued` -- admitted, not started, not free -- rather than
borrowing `Thinking`, which would have double-counted in
`ActivityModelCycle` and mislabelled a span in which no model is being asked.
Only `Submit` mails enter that queue, so a non-empty queue is always real work.
The header now defers to the server's live phase whenever the turn-level word is
not in flight, and shows the detail sentence, which is written for exactly this
display. `paused` still surfaces running jobs, as before.

`ActivityQueuedSubmissionIsNeverIdle` settles an agent first, so a later non-idle
answer can only come from pending work, and pauses it so the submission is
deterministically queued rather than racing a dequeue.

## What the timing row was actually doing

Bill also quoted `Total 0.0s · Model 0.0s · Tools (0) 0.0s`. Measured on the
current build, the row is honest: `Total 0.2s · Model 0.1s · Tools (1) 0.1s` for
a turn that ran one command, with `tool_dispatched` observed on the wire.

The `Tools (0)` he saw was correct at the time, and for a reason worth
remembering: that agent's `run_command` returned `Tool is not enabled:
run_command`, because a GUI agent has no tools until a skill grants them, and a
refused call is never dispatched. The zeros for model and total time were the
pre-fix build from the earlier defect.

One of my own checks was vacuous while finding this: it asserted the page
contained `SUMMARY-MARKER-73`, which the tool card also prints as the echoed
arguments, so it passed while the command had not run. The check now uses a
marker the shell computes (`echo SUM-$((6*7))RESULT`) and requires both the
computed output and the echoed formula.

## Fixture note

`OnceOnlyResultStubbing` scripted exactly five requests, which quietly encoded
the assumption that an empty reply ends an exchange. It now serves while the
scenario asks, and tolerates a client hangup on read or write, because the actor
may abandon a request when a turn ends. Its per-round keep/stub checks still run
on every round, and the awaited turn results are what guarantee the requests
happened -- a missing request hangs until the deadline. The `rounds >= 4` guard
covers only the rounds whose bodies the loop interprets by number.

## Verification

`agent/Tests`: 368 passed, 2 failed -- the two macOS Seatbelt confined-command
cases, unchanged. `agent/GuiTests` clean. `job-status.cjs` 21 checks including
four for pending phases, each failing when the header rule is removed. Real
headless Firefox on the real launcher: the header holds `inflight` across a
five-second provider hold and becomes `idle` only after the answer arrives, with
no page errors.

## Limits

- The browser proves the two states a user can act on. The queued phase itself is
  proven in C#, because a second submission during a live turn arrives as a hint
  observation rather than a queued submit, so it cannot be forced from the page.
- A live tool phase shows the server's detail sentence, which is English and
  fixed; humanising those words further is still Bill's call.
- A false *busy* is still possible in the other direction: a client that dies
  without closing its socket leaves the last published phase standing until the
  next poll. The 200 ms tree poll re-reads the server, so it self-corrects, but a
  wedged actor would keep claiming work.
