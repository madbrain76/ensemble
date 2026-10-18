# Revision-only process and failed inspections

The original author's failed inspections and false evidence descriptions remain in its untouched report/evidence. This note records this sole revision's separate process, without relabeling original failures as revision checks.

## Full reading

Read the complete task and formal review before edits, then the original report, current chapter, public student brief, current author/repository guidance, accepted fixture/live/native records, code reviews and sanitized primary receipts used for attribution. Reloaded the workflow in lines 1–840 / 841–1680 / 1681–2520, carryover in 1–620 / 621–1241, and the full voice/procedure documents. Each output cap exceeded its range bytes. Exact paths came from the original ledger after the failed path guess below. The ledger pins these source files but does not pretend all historical student evidence was newly executed or independently re-reviewed.

## Revision setup failures and retractions

1. A guessed `book/llm-development-workflow.md` path did not exist. Used the original ledger's actual `docs/agentic-codebook-workflow.md` instead and read it completely. No source write resulted.
2. The first snapshot attempt used incorrectly copied expected review hashes. Its assertion failed before protected inputs were read or the chapter changed; only the empty new evidence directory existed. Re-read the task and complete formal review; the real hashes matched the supplied authoritative values. No review drift was found or inferred.
3. The next attempt assumed the saved exclusion metadata used a `paths` key and stopped with `KeyError`. The actual metadata has a `files` mapping. Only its twelve keys were subsequently used to exclude publication paths before hashing. No publication file was opened, hashed or restored and no chapter write occurred during these failures.
4. Provisional visible commentary incorrectly called the fixture field `yield_time`. The actual receipt says `ai_callback_delay`; this was explicitly corrected before the report. The final report uses only `delay -> ai_callback_delay`.
5. Provisional visible commentary also mentioned a SpeechViewer utterance and real microphone capture. Those claims are unsupported and withdrawn. The actual native receipts record headed Chrome, native callbacks, a separate qualified human quote and microphone **never requested**. No SpeechViewer or microphone proof is asserted in the final report.

## Actual revision and checks

After successful snapshots, one fail-fast chapter write made exactly the review's three local changes: paired gate recipe, engine-to-core attribution, and deletion of the unsupported opaque-description sentence. All anchors were unique and the before chapter matched the reviewed hash. No manual GUI change or style cleanup was made.

Ran `go run ./cmd/lintprose` on baseline/chapter-07.md and final/chapter-07.md. Both exited 0 and their advisory warnings are retained verbatim. Ran `git --no-pager diff --check` scoped to the chapter, shell syntax parsing only, exact content/heading/link/manual-GUI checks and preservation/index checks. Grading/build commands in the chapter were not executed. All Git invocations used `--no-pager`; no pager was opened during this revision. No children, asynchronous jobs, application tests, provider/browser/native runs, publication tasks, commits or tags were created. Every tool command returned a completed exit status before handoff.
