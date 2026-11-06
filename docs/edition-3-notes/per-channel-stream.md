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

## Correction 2026-11-06: the flag was rejected by its own launcher

Nothing in the suite starts the GUI, so two launcher defects survived on this
branch. `--stream` was listed both among the valueless switches and among the flags
that take a value; the valueless branch runs first, so `--stream all` recorded a
switch and then tried to parse `all` as a flag — "Unknown or incomplete GUI flag:
all". The configuration path was worse: `Selected("--stream", "EN_STREAM")` returned
an empty string when the flag was absent, `Configuration` forwarded it, and
`ParseStream` refuses an empty selection, so the GUI could not start at all without
`--stream`. An absent flag now keeps the `Configuration` default, and only a value
the operator typed is forwarded, so `--stream ""` is still refused — in the
validator's own words rather than as a silent "all".

Measured on the built launcher: `all`, `text,tools` and `thinking` serve; `bogus` is
refused as an unknown channel; `""` is refused by name; no flag serves.
