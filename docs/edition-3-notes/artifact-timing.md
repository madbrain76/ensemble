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
