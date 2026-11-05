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
