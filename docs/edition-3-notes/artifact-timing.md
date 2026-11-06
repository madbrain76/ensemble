# Ported feature: per-artifact timestamp and production duration

Branch `port/14-artifact-timestamps`, based on `fix/linux-path-case`. Coordinator,
2026-11-03, at Bill's request. Ported from `go-checkpoint` (`043e3c87`, `0dbb1cdd`,
`a6018d74`, `902b1497`).

Each output box now opens with `HH:MM:SS · 3.0s`: when that artifact began and how
long the model took to produce it. Arrival is measured in the page as chunks land
(`card.first`/`card.last` on the card), and a frame carrying a server `time` wins,
so a replay keeps its original moment rather than the reload time.

Measured in the browser, deliberately: putting it in `Event` would extend a
durable History contract that chapters and graders already pin, for a fact the
reader only ever wants visually. A part with no arrival shows nothing rather than
a misleading `0.0s`, and a part still streaming shows `· …` instead of a duration
that would immediately be wrong. The timing line is written on finalization only —
rebuilding it on every delta would mean a DOM node per chunk on a long answer.

`agent/GuiTests/artifact-timing.cjs` extracts the shipped formatter from
`app.js` and pins: no arrival shows nothing, streaming shows the ellipsis, a
zero-span final never claims `0.0s`, sub-ten-second values keep one decimal,
longer values round to whole seconds, and a replayed part keeps the server clock.

## Correction 2026-11-06: the line never appeared on a live turn

Everything above describes what the formatter produces, which is how this note came
to read as a feature that worked. It did not: `renderTiming` was called only from
the `part_final` branch, and `part_final` is emitted only by the replay projection,
so during a running turn no clock line was drawn at all — it appeared only after a
reload. Two further breaks sat behind that one: `render()` rebuilds a card with
`replaceChildren()`, wiping the line just prepended, so a reply delivered in a
single chunk never showed a duration even after the turn closed; and nothing in the
live stream marks a part finished, so the duration would have stayed an ellipsis
forever. The frame handler now renders on `part_delta`, `render()` keeps the timing
line across a rebuild, and `closeTurn()` finalises the cards that are still open,
because at that moment production genuinely has ended.

`agent/GuiTests/timing-live.cjs` covers the path the formatter suites could not,
driving the shipped `handle`: an in-flight part shows the clock with an honest
ellipsis, the same part carries a measured duration once the turn closes. Each
fails on its own when the corresponding call is removed. Measured in Firefox on the
real launcher: `03:51:03 · 1.4s` and `03:51:03 · 0.7s` on a streamed reply.

## Replaying a large history froze the page

Selecting `ensemble debugger` (3,624 cards, 3.36M characters) blocked the main
thread for **57.98 seconds** in a single task, which is what made Firefox offer to
stop the page. `render()` read `pane.scrollTop`, `pane.scrollHeight` and
`pane.clientHeight` before appending each card, so every insertion forced a fresh
layout of an ever-growing subtree: quadratic, measured at ~15 ms per card, and all
of it inside `render` (57.0 s self time out of the 57.98 s task; every callee
together accounted for under 0.4 s).

Cards now accumulate in a detached fragment per pane while `replaying` is set, and
the scroll position is decided once. That alone took the task from 58.0 s to 1.87 s.
The remaining second was attaching all 3,624 nodes at once, so only the most recent
150 are attached before the pane is usable and the rest arrive in animation-frame
slices that grow the transcript outward from what the operator is reading, each
slice re-pinning the scroll position so the view does not drag backward. A card
re-rendered by a live frame while its slice is still pending is dropped from the
backfill and attached in place, and a newer Agent selection bumps a generation
counter so an in-flight backfill cannot interleave two transcripts.

Real history, same click: **0.92 s** worst task, down from 57.98 s. Synthetic 1,200
parts in `agent/GuiTests/replay-performance.cjs`: 127 ms synchronous replay, 147 ms
worst task, complete in arrival order with no duplicates.

Two things worth keeping from this. The order of a backfill slice was wrong until
the test checked it — `pending.pop()` yields newest-first, and prepending that into
a fragment reverses the chunk, so history would have read scrambled. And the
first version of the perf test passed against the unfixed code, because an interval
watchdog records no samples at all while the thread is blocked: `Math.max(...[])`
is zero, so "no task longer than 1000 ms" was satisfied by the worst case. The test
now requires evidence that the watchdog ran before it is allowed to report it.
