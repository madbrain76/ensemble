# Chapter 10 restart — formal text review

**Verdict: PASS.** No required revision and no optional polish requested. This is the single post-student text review; root owns closure.

## Scope

Reviewed the full current chapter and the actual author-only diff against `/tmp/ch10-author-before.dexX0n/chapter-10.md`: four hunks, **22 added / 3 removed lines**. Read current `AGENTS.md`, execution instructions, author/text-reviewer roles, `book/voice.md`, `book/chapter-writing-procedure.md`, and the restart student, code-review and author records. Applied root's narrow-phase waiver of obsolete workflow/archive reads. Root's independently accepted entry-backup amendment is present in the baseline and unchanged by these hunks.

Root supplied code/scoped-operations acceptance after **3/10** correction rounds and the accepted 120-file `ch10 == main` snapshot; this prose review does not re-establish that acceptance. The corrected author note now accurately maps R1 to EOF, R2 to direct dispatch and directory ownership, and R3 to semantic comparison.

## Reasons

- **Direct dispatch (§10.3):** The addition distinguishes legitimate public `Tools().Dispatch` events outside a model turn from model-supplied calls. `internal/tool/dispatch.go` preserves the supplied call ID, while `internal/history/persistence.go:eventValid` permits its absence for turn-zero tool events and still requires it for model-turn events. The prose neither invents historical IDs nor weakens correlation.
- **Decoded comparison (§10.4):** `internal/history/persistence.go:sameContext` compares decoded JSON using `UseNumber` and normalizes empty in-memory collections. The paragraph explains why formatting and nil/empty differences must not reject equivalent context, and why rounded numbers can conceal a changed argument. It expressly retains required snapshot fields, facts, sequence ordering and anchors; it does not turn representation tolerance into permissive loading.
- **Directory ownership (§10.6):** `ensemble.go:NewAgent` resolves relative `DataDir` beneath the owner's `BaseDir` when deriving the default path and leaves an explicit `SavePath` caller-selected. The addition states that distinction without prescribing extra configuration machinery.
- **EOF and usage (§10.6):** `gui/terminal.go` drains private `SubmitLive` completions before terminal completion; `gui/cmd/main.go` then shuts down. `internal/cli/cli.go` emits usage before deferred shutdown/save. The addition correctly distinguishes admission, reliable completion and successful persistence, explaining why observer silence or a usage record cannot acknowledge saving.

The additions are load-bearing, student-grounded clarifications in the chapter's direct technical voice. Existing examples, first-person history and reasons remain intact. There is no fabricated experience, new feature requirement, repair diary or stronger operational claim. Backup recovery remains bounded snapshot recovery, not an early journal or universal power-loss guarantee.

## Limits

Only this review note was written. No chapter/source edits, tests, builds, relint, browser/provider activity, research, children or commits. Existing lint results remain reported evidence: baseline/final exit 0, with retained soft warnings of 15 negations and person gap 2825 → 2928; no fresh lint pass is claimed. Historical receipts remain historical. Restart operations establish local rendered-DOM behavior, not pixel inspection or native listening, and one live Anthropic read/EOF/save/restart-recall exercise, not broader provider coverage. Earlier controller failures remain in their existing records rather than entering the chapter.
