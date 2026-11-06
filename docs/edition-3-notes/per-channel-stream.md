# Ported feature: per-channel --stream selection

Branch `port/8-per-channel-stream`, based on `fix/linux-path-case`. Coordinator,
2026-11-03. Ported from `go-checkpoint` (`a6e23122`).

`--stream` (or `EN_STREAM`) accepts `all` or a comma list of `text`, `thinking`,
`tools`, and is validated at construction: an empty selection or an unknown
channel is refused with the accepted values in the message, because a silent empty
selection would look like a stalled model. Unselected channels simply do not stream
live — the settled response still delivers every part, so the selection changes
watching, never the answer. `StreamName` round-trips the choice so a model switch or
a security narrowing cannot quietly reset it to `all`.
