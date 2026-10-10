# Edition 2: bounded enhancement pass

Authorized by Bill on October 10, 2026, after acceptance of original Chapters
1–9 at `db96c20` / `edition-2-ch09-r2`. Coordinator: Codex.

Bill: “Please proceed autonomously with the author and coder for chapters 1–9
(1–8 after deleting 5). Have a reviewer review the chapter and also the code for
each chapter. The goal is improving the chapter, not rewriting.” The preceding
agreement includes enhancing Chapter 0 with execution guidance and Chapter 1
with upfront architecture, using the carryover and workflow records.

## Boundaries and working agreement

- Enhance copies of the original chapters in `book/edition-2/`. Preserve their
  titles, stories, argument and useful explanations. Change only what recorded
  lessons, corrected architecture, clarified contracts or numbering require.
  Start from the original text, never the discarded Edition 2 draft.
- Author and prose reviewer read all of `book/voice.md` and
  `book/chapter-writing-procedure.md`. The old procedure's author-first sequence,
  fixed whole-book chapter count and expanded-contract machinery do not govern
  this pass. Bill's enhancement constraint takes precedence over stylistic
  cleanup that would rewrite sound original prose. Record inherited lint issues.
- Every coding/review handoff includes `.agents/skills/ensemble-coding/SKILL.md`.
  Read it and `docs/edition-2-carryover.md` in full. Use the workflow, chapter
  notes, student reviews, comparative reviews and accepted implementation.
- No new product features or speculative acceptance requirements. Author
  explains demonstrated behavior. Coder repairs documented grader defects and
  meaningful missing coverage; reviewer rejects implementation-shaped checks,
  bloat and weakened valid coverage. No mocks without Bill's approval.
- Preserve first-edition chapters/graders and all frozen source exports/results.
  Edition 2 graders must have an explicit separate entry point and preserve
  valid original checks. Use historical chapter IDs for evidence paths.
- Do not edit `agent/` or unrelated user work. The deleted old Edition 2 Chapter
  1 is intentionally replaced by an enhancement of original Chapter 1.
- Review one chapter before proceeding to the next. Author and coder may work
  together on that chapter. Default to three correction rounds after initial
  independent review. A critical scope question pauses the work for Bill.
- Chapter 10 student development and prose beyond the accepted boundary wait.
  No publication or push. Commit completed reviewed stages with explicit paths.

## Mapping and retained Chapter 5 lessons

| Edition 2 | First-edition source | Placement |
| --- | --- | --- |
| 0 | 0 | Preserve opening; add practical execution/review guidance. |
| 1 | 1 | Architecture and discipline before code, using current minimal objects. |
| 2 | 2 | Event history and request boundaries. |
| 3 | 3 | Tools; explain instance registration and owner access. |
| 4 | 4 plus relevant 5 | Jobs plus the existing external-consumer/custom-tool proof. |
| 5 | 6 | Actor, steering and reliable completion. |
| 6 | 7 | Streaming. |
| 7 | 8 | GUI, replay and Agent-wide pause. |
| 8 | 9 | Workbench and persistent settings. |

Original Chapter 5's refactor disappears; its external-consumer capability
remains. Folding that proof into Chapter 4 uses the accepted original Chapter 5
snapshot as the end-of-chapter reference, avoiding moving later job machinery
into the first tools exercise. Teach ownership and logging motivation earlier,
including the documented logger ambush where useful. Preserve original IDs/tags;
do not rename or modify frozen snapshots. Do not decide the rest of the book's
numbering in this pass.

## Per-chapter handoff and acceptance

Author records in `docs/edition-2-notes/enhancement-NN.md` the original source
revision, additions/corrections and why they belong, stories retained, relevant
notes incorporated or deferred, and any changed code examples/contracts. Keep
the record brief; do not paste internal workflow into the teaching prose.
Coder records exact grader commands, results, focused mutation evidence and
any corrected implementation separately from historical scores. Reuse unchanged
student live evidence with its provenance; new behavior requires its own checks.
Reviewer reads original and enhanced prose, relevant code, grader changes and
evidence. Assess teaching/enjoyment and technical accuracy, preservation of voice,
KISS, comment quality, test sensitivity and coverage. Report concrete findings,
then reassess corrections; do not silently fix the work being reviewed.
Coordinator checks the dispositions, mapping and source diff before acceptance.

| Edition 2 chapter | Author/coder | Independent review | Checkpoint |
| --- | --- | --- | --- |
| 0–1 | Preparing | Pending | — |
| 2 | Not started | Pending | — |
| 3 | Not started | Pending | — |
| 4 (includes old 5) | Not started | Pending | — |
| 5 (old 6) | Not started | Pending | — |
| 6 (old 7) | Not started | Pending | — |
| 7 (old 8) | Not started | Pending | — |
| 8 (old 9) | Not started | Pending | — |
