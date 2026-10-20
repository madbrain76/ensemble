# Chapter 10 restarted author pass — HOLD

Post-student text-only pass. Root supplied acceptance of code and scoped operations after 3/10 production rounds; this pass does not re-establish that acceptance. Root owns the sole text review, with at most one author response.

## Scope and exact diff

Read current `AGENTS.md`, `docs/edition-3-execution.md`, `.agents/skills/edition-3-author/SKILL.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, the full current chapter, current carryover, both concise restart student/code-review notes, snapshot-recovery followup, and Chapter 14 retention-entry obligation. Historical workflow/Edition 2 campaign archives were waived for this narrow phase. Read relevant accepted source paths to verify the clarifications; no source edits.

Saved the current chapter outside Git before editing, including root's existing entry-backup amendment:
`/tmp/ch10-author-before.dexX0n/chapter-10.md`.

Author-only chapter diff: four hunks, **22 added / 3 removed lines** (664 → 683 lines):
- Baseline line 314: public direct `Tools().Dispatch` may record no native call ID outside a model turn; do not invent IDs or weaken model-turn correlation.
- Baseline line 364: replay/snapshot comparison concerns decoded meaning, not JSON layout; normalize empty in-memory collections and preserve exact numbers without relaxing required snapshot fields, facts, sequences or anchors.
- Baseline line 487: derived relative `DataDir` follows owning `BaseDir`; an explicit `SavePath` remains caller-selected.
- Baseline line 508: EOF ends admission, not completion; reliable completion precedes shutdown, and CLI usage output precedes deferred saving rather than acknowledging it.

Evidence mapping: student findings/round 3 semantic-comparison repair; round 2 direct-dispatch repair; round 1 EOF repair; round 2 owning-directory repair; final manual usage-timing correction. Verified against `solutions/edition-3/main/internal/history/persistence.go`, `internal/tool/dispatch.go`, `ensemble.go`, `internal/cli/cli.go`, and GUI terminal/command lifecycle source. Existing diagnostic safe-boundary and card-continuity explanations were already adequate and were not expanded into a repair diary.

Root's backup amendment, Bill's voice/reasons/history, grader 100/seven-check explanation, current-configuration rule and backup cases remain unchanged from the saved baseline. No new first-person experience or stronger operational claims were added.

## Commands and results

Read the actual interface in `cmd/lintprose/main.go`. From `/Users/bill/projects/ensemble`, same-basename lint commands:

```sh
env GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOFLAGS='-mod=readonly -p=1' GOMAXPROCS=2 go run ./cmd/lintprose /tmp/ch10-author-before.dexX0n/chapter-10.md
env GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOFLAGS='-mod=readonly -p=1' GOMAXPROCS=2 go run ./cmd/lintprose book/edition-3/chapter-10.md
diff -u /tmp/ch10-author-before.dexX0n/chapter-10.md book/edition-3/chapter-10.md
git --no-pager diff --check -- book/edition-3/chapter-10.md docs/edition-3-notes/ch10-restart-author.md
```

Both baseline and final lint exited **0**. One intermediate post-edit lint also exited 0; final rerun followed the required-field wording clarification. Soft warnings remain: negation forms **15 → 15** (word-count allowance 9 → 10), person gap **2825 → 2928**. Other categories remain OK. Warnings are advisory, not a mandate to cut useful explanations. Author-only `diff -u` exited **1**, meaning the four expected text hunks differ. Final scoped whitespace check exited **0**.

Source navigation initially looked for nonexistent `internal/tools`; those lookup failures were corrected to `internal/tool/dispatch.go`, not treated as product defects. `df -k .` remained above the 1 GiB floor (last observed 3,517,800 KiB available); no cleanup.

## Limits

Only the chapter and this note were written in the repository. The temporary chapter baseline is the sole outside-Git author artifact. No code/tests, grader run, provider request, browser, research, children, commit, deployment or publication build. Historical passing receipts were not rerun or relabeled as new passes. Existing manual evidence remains local rendered DOM, not pixel inspection, plus one Anthropic read/EOF/save/restart-recall exercise. Earlier controller failures remain in their existing notes. Snapshot backup recovery is neither a journal nor a universal power-loss guarantee.

**HOLD** for root's actual-diff text review. No text-review verdict claimed.
