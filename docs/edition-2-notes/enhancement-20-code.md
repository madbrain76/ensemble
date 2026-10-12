# Edition 2 Chapter 20: technical enhancement in progress

Coder: GPT-6 Astra/high. Scope agreed with coordinator and independent reviewer:
original Chapters 21 + 22, against the frozen combined `solutions/edition-2/ch20`.
The entire `.agents/skills/ensemble-coding/SKILL.md`, carryover and workflow were
reloaded after compaction, with both originals, student/comparative records and
original grader implementations read. No runtime or frozen source changes.

The combined rubric is 155 points: the original seven web checks retain 100;
build, status, accounting and resolved tool/common imports add 5/20/15/15.
The old ten-point web parity row is represented directly by the seven checks.
The old parent and composition shape detectors become mandatory independent
source review, not automatic credit for identifiers or guessed structs. Their
actual ownership and executable construction obligations remain binding.

Planned checks use actual CLI web dispatch and GUI turn-boundary model changes,
typed skills, correlated full results, extra advertised but ungranted tools,
HTTP JSON/SSE/session/auth, refusal after unload, real error propagation and
successful EOF. Normal tool timeout, blocked connection cancellation, redirect
and credential-redaction controls retain unchanged source/mutation evidence.

Status uses supported Sonnet-5 → Opus 4.6 → Sonnet-5. The accepted frozen table
prices only Sonnet among those rows; known-price work is checked numerically and
unpriced work must remain disclosed after returning to Sonnet. Legacy prices are
explicit fixture inputs, not silently inferred expectations. Historical accepted
two-priced-model evidence retains its actual old identities; no new 5.5 run or
model-table edit is authorized. All four usage buckets enter the cache denominator.

Original web 0/100 and architecture 20/100 remain unchanged historical results.
The first revised baseline was 35/155. Its HTTP fixture omitted `transport: url`,
so actual `load_skill` returned `invalid MCP server declaration`; the downstream
web failures were masked by that setup error. A retained-dialogue assertion also
used the recall-specific extractor, which intentionally excludes ordinary human
text. Both fixture mistakes are corrected. Inspection also replaced a flat-key
meter comparison with the actual nested cumulative counts before/after selection.
The initial report remains `accepted-initial.json`; it is not a runtime diagnosis.

First green full commands:

```
go run ./cmd/grade-edition2 -ch 20 -json solutions/edition-2/ch20
go run ./cmd/grade-edition2 -ch 20 -gui solutions/ch22 -status-prices 3,3.75,0.3,15 -status-second-prices 5,6.25,0.5,25 -json solutions/ch22
go test -race ./internal/gradeedition2 -run '^TestChapter20Combined$' -count=1
```

Accepted: **155/155**, exit 0 (`accepted-revised.json`). Legacy: **135/155**,
grade exit 1 (`legacy-initial.json`): its status tool omits cache-write tokens
from the hit-rate denominator, independently confirmed in its formatter. Both
web sessions and all cost assertions pass. Focused race passes in 2.054 seconds.
These fresh sessions use distinct temporary directories and ephemeral ports;
no user listener or account profile is accessed. Evidence is under original
`solutions/edition-2/evidence/ch21/enhancement20/` for the combined chapter.

Focused sensitivity, exact examples and independent final acceptance remain
pending. No later chapter is released by this green stage.
