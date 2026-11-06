# Ported feature: per-turn Total / Model / Tools summary

Branch `port/15-turn-timing`, based on `fix/linux-path-case`. Coordinator,
2026-11-03. Ported from `go-checkpoint` (`84cbd438`, `770249af`).

The page accumulates wall-clock spans from the frames it already receives:
`prompt` submission opens the turn, `inflight` is model wait, `toolspending` is
tool wait, and the first non-busy state closes it and prints
`Total 11.9s · Model 10.9s · Tools (1) 1.0s` under the exchange.

Two deliberate omissions. There is **no "Other" column**: from the page every span
is model or tool wait, so the bucket would sit at `0.0s` forever and invite the
question of what it measures — exactly the exchange Bill had during the Go work.
The real phase split (admission queue, context copy, request render) needs
server-side instrumentation and is deferred rather than faked. There is also no
durable record: storing it would extend the graded `Event` contract for a
visual-only reading.

Two defects surfaced while testing: durations rounded at ten seconds, which is
precisely where a quick answer is distinguished from a slow one (now one decimal
below a minute, `Xm YYs` beyond), and `accountTurn` did not advance its segment
start, so the closing state charged the final tool span twice.
`agent/GuiTests/turn-timing.cjs` runs the shipped state machine in a vm and pins
both, plus the tool-free turn and the no-summary-before-completion case.

## Correction 2026-11-06: the summary published zeros for a live turn

As on the artifact note, the state machine was tested in isolation and the frame
path was not. Two defects reached the operator. `recordTurnTiming` runs at the top
of `handle`, before the per-Agent filter further down, so a sibling Agent's state
frame drove this Agent's accounting. And the first non-active state closed the turn
unconditionally: the prompt is admitted and the actor has not reached the send yet,
so a 200 ms poll could catch an idle frame and publish `Total 0.0s · Model 0.0s ·
Tools (0) 0.0s` for a turn that had not started — a row of confident zeros, the same
class of lie as a tool result claiming to have run.

Accounting now follows only frames for the selected Agent, and an idle state closes
the turn only once it has been inflight or has delivered live output. Firefox on the
real launcher: `Total 1.5s · Model 1.5s · Tools (0) 0.0s` for a 1.4 s streamed
answer. Covered by `agent/GuiTests/timing-live.cjs`, whose two turn cases fail
independently when either guard is removed.
