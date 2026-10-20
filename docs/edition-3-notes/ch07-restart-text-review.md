# Chapter 7 restart — formal text review

## Verdict: PASS — unchanged after report correction

No material chapter revision conditions. No author response is requested. This is the corrected report of the same single formal review, not a second formal review; root acceptance remains the root's decision.

Reviewed current `AGENTS.md`, coding and text-reviewer skills, `docs/edition-3-execution.md`, full `book/voice.md` and chapter-writing procedure, full current `book/edition-3/chapter-07.md`, the restart author note, and the actual student/code-review restart notes, including their final acceptance sections. Inspected `git --no-pager diff 039e977b -- book/edition-3/chapter-07.md`: six localized clarifications, 19 insertions and 3 deletions. Historical workflow/carryover archives were waived for this bounded phase. The immutable author first-draft copy was left untouched.

## Teaching and contract

The six changes are useful and consistent with the accepted contract:

- **Empty chunk key:** an empty `part_delta` must still carry `"chunk":""`; omitting the key changes the public schema. This is a schema-preservation clarification, not an identity-key or speech-queue instruction.
- **Admission, save, release:** saving and publishing the dispatch fact share the same admission boundary as releasing the handler. Cancellation after publication may stop admitted work, but must not make the published dispatch describe a reservation that was never released. This sharpens the existing distinction between a request and an actual dispatch without prescribing a new framework.
- **WebSocket Pong:** each Ping payload is echoed in a Pong control frame through the writer. An application-level JSON message named `pong` is not that control reply. This preserves the sole-writer lesson.
- **Card identity:** dispatch status and job metadata can change while `(agent_id, part_id)` continues to identify the same logical part. Reusing that identity for a different part is a collision. The addition does not prescribe a fatal-collision policy.
- **Replay count:** the window remains the last hour OR the last 100 renderable History events, whichever reaches farther back. An event with several rendered parts counts once. Replay batches are measured by those entries, not socket-envelope or unique-card counts; reset frames and rebuilt cards are not substituted for the History-event unit.
- **Native speech observation:** attach start, end and error listeners before playback begins so a fast start is not missed. This is observation-timing advice, not a new product requirement or a claim of hearing.

On the corrected reading, no material contradiction among the chapter's contract, explanation and accepted behavior is identified. No optional polish or invented requirement is promoted to an acceptance condition.

## Voice and evidence

The opening, historical account, human motivation, ownership rationale, examples, exercises and conclusion remain intact. The additions explain concrete boundaries in the chapter's existing technical voice; they do not replace its narrative with an operational report.

The final restart records support acceptance of code and bounded operations after four of ten authorized correction rounds. Their qualifications remain controlling: actual native callbacks are not human-hearing evidence; screenshots received delegated model assessment, not root/reviewer pixel viewing; the live file read preceded the edit, with no post-edit reread; earlier STOPs remain preserved. The chapter neither claims otherwise nor needs to narrate every current-run failure. No code review or browser campaign was reopened.

Carried the author's actual recorded baseline/final `lintprose` results: both commands exited 0. Both retain soft warnings for **20 negation forms**; the person gap increases from **2,334 words after L414 to 2,428 words after L420**. Short closers remain 35, against 105 baseline and 107 final paragraphs. The other displayed checks are `ok`. These soft warnings do not justify invented personal experience or rewriting the retained technical rationale. No lint was rerun.

## Erratum to the first saved report

At root's request, preserved the first report at `/tmp/ch07-restart-text-review-first-note.md`, then reread the actual current chapter diff against `039e977b` and the author note's full displayed baseline/final lint outputs.

The first report incorrectly described the empty-chunk edit as a `(turn_id, text_id)`/speech-queue instruction, reversed the replay counting requirement by substituting emitted frames/cards for renderable History events, and invented a fatal cross-node/cross-kind collision policy. It also incorrectly reported nine voice warnings and eight style warnings as unchanged. Those statements are withdrawn and replaced above. The unsupported claim that the author's recorded checks included a clean `git diff --check` is also withdrawn. These were review-report errors, not chapter defects.

**The PASS verdict is unchanged after this accuracy reconciliation:** the actual six edits clarify existing boundaries and preserve the chapter's voice without adding requirements. No author revision is needed to repair the reviewer's mistakes.

Only this project review note was edited, plus the requested preservation copy in `/tmp`. No chapter/source/grader edits, tests, builds, lint execution, providers, browsers, credentials, children or commits were used for the review or its report correction.
