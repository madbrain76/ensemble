# Edition 2 Chapter 14 grader — stage, review pending

Coder: GPT-6 Astra/high. Mandatory full reading:
`.agents/skills/ensemble-coding/SKILL.md`, carryover and workflow, original
Chapter 15 at `994bf218` (SHA-256
`6c6d953e754f73b8bafd45a2f6009cab5dce9d7557c008e79fcc4084f27c2e79`),
and complete original15 student/comparative records. Accepted direct solution
is `solutions/edition-2/ch14`, tree `0c751dbd3985dee1a80284461e31346471ccc984`;
legacy is `solutions/ch15`. Neither source nor original grader is edited.

The separately versioned command is `go run ./cmd/grade-edition2 -ch 14 -json`.
It builds the optional GUI; `-gui` selects a different existing GUI command.
Legacy invocation adds `-gui solutions/ch15 solutions/ch15`. The original eight
weights are retained: 15/15/10/15/10/15/15/5, totaling100 without parity.

The runner applies `context_target` through public `update_settings` and checks
the complete `settings_changed` snapshot before prompting. Ladder uses20000;
the restart experiment explicitly applies20000 versus400000. Keep/no-stub uses
400000, preventing the ladder from confounding per-round stubbing. An actual
`turn_ended` observation replaces the original400ms delay before EOF or SIGKILL;
normal EOF now requires a successful process exit. The first local runner forgot
`subscribe`; `initial-test.txt` retains its zero-request setup failure. Adding
that existing public message made `subscribed-test.txt` pass in0.971s.

The manual grants existing `list_directory`, replacing undefined `think`.
Original scenario assertions remain. New checks require successful correlated
file delivery before later absence proves a cut, a real same-batch read alongside
`micro_handoff`, and its ordinary dispatch/result/separate-note chronology.
The prefix scenario now loads, invokes and unloads the external MCP peer, checks
actual declaration/removal block kinds, changing output-only automatic content,
correlated argument-dependent results and strict inline-system placement after
the whole user turn. It reuses Chapter12's small faithful subprocess peer.

Actual full revised results: accepted100/100 in0.814s, legacy100/100 in0.842s,
both exit0. Raw command receipts and reports are in
`solutions/edition-2/evidence/ch15/enhancement14/`. Source binding is
`stage-source.json`; only documentation comments changed after those runs.
Formatting, focused test and affected vet pass. These are first-green results,
not final mutation or independent-review acceptance.

The original45/100 and partial live Anthropic/OpenAI task failures remain in
`ch15.md` and `ch15-code-review.md`. The course-only model feature matrix governs
these deterministic inline/stubbing checks; no real provider cache, inline or
stubbing success is inferred. Existing source tests retain richer reduction,
journal/checkpoint-failure, feature-table and GUI/native proof. No paid rerun is
needed for unchanged evidence. Exact author examples and targeted sensitivity
checks remain pending before final submission.
