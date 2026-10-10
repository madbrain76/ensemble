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
| 0–1 | Complete | Accepted after one prose correction | `02efcd4`, `edition-2-enhanced-ch01-r1` |
| 2 | Complete | Accepted at initial formal review | `ec8c901`, `edition-2-enhanced-ch02-r1` |
| 3 | Complete | Accepted at initial formal review | `601791a`, `edition-2-enhanced-ch03-r1` |
| 4 (includes old 5) | Complete | Accepted at initial review with recorded source-test exception | `edition-2-enhanced-ch04-r1` |
| 5 (old 6) | Not started | Pending | — |
| 6 (old 7) | Not started | Pending | — |
| 7 (old 8) | Not started | Pending | — |
| 8 (old 9) | Not started | Pending | — |

### Chapters 0–1 acceptance

Author `enhancement_author`, coder `enhancement_coder`, independent prose/code
reviewer `enhancement_reviewer`. Combined review: [enhancement-01-review.md](edition-2-notes/enhancement-01-review.md).
Preserved the original stories and section progression; added execution guidance
and architecture before code. Corrected one `go run` exit-status explanation.
The separate Edition 2 grader retains all seven original checks and fixes the
API-version-header survivor. Accepted and legacy source both score 100/100;
wrong-version source scores original 100 versus revised 85, with only wire failing.
Original code, graders and student results are unchanged. Historical Chapter 0
lint issues and frozen Chapter 1 documentation gaps are recorded, not called passes.

### Coordination handoff, October 10

Incoming coordinator thread `01a12710-86d5-7dc0-af0f-1afc272964cc` reported Bill's
fresh-session instruction to continue and wake the author/coder/reviewer. This
coordinator has stopped starting chapter work and asked active delegates to
checkpoint safely for takeover. Last accepted enhancement remains `02efcd4` /
`edition-2-enhanced-ch01-r1`; later first-edition Chapter 25 commits belong to
Bill's other work and must be preserved.

Existing delegates in this collaboration tree:

- `/root/enhancement_author`: Chapter 2 prose and source/disposition notes ready;
  independent initial prose review found no blockers. Completed read-only Chapter
  3 preparation but wrote no Chapter 3 file. Chapter 2 files are
  `book/edition-2/chapter-02.md` and `docs/edition-2-notes/enhancement-02.md`.
- `/root/enhancement_coder`: Chapter 2 grader implementation/validation in
  progress in `cmd/grade-edition2/main.go` and new
  `internal/gradeedition2/ch02{,_ref,_test}.go`. Requested safe handoff and
  `enhancement-02-code.md` with exact results/pending checks. No Chapter 3 work.
- `/root/enhancement_reviewer`: Chapter 2 prose initially clear; code review and
  validation pending. Requested `enhancement-02-review.md` checkpoint. No final
  Chapter 2 acceptance or completed correction-round count is claimed.

Pending Chapter 2 work at takeover:

1. Coder complete baseline and applicable legacy results, focused mutations,
   and execution of the author's exact nine-event public sample. Initial
   supplement failed because its fixture omitted tool-call provenance; this was
   a grader setup failure, not a student defect. The coordinator also required
   documented timestamps and complete input fields. Revised fixture is on disk.
2. Reviewer finding: redaction check accepted a nonempty result object whose
   actual text could be empty. Require observable nonempty stub content and a
   targeted mutation. Keep the claimed dump preservation limited to what the
   check observes; exact reduced Ref preservation stays in source tests/review.
3. Check useful comments on new grader fields, including anonymous JSON decode
   structs, under Bill's exported-symbol rule. This does not reopen frozen
   historical documentation gaps.
4. Final independent combined prose/code decision, coordinator source-diff
   check, explicit-path commit/tag; only then release Chapter 3 drafting/coding.

Approved Chapter 2 grader plan: reuse original checks; remove premature URI
rendering (4 points), assign those points to corrected Ref/redaction (8), keep
100 total; strengthen result names, captured identity/signatures, interleaved
Anthropic results and corrected Chapter 1 parity. No new public API for grading,
no original grader edits, no product features. Legacy defects found by valid
new checks must be reported, not hidden to force a pass.

Cross-thread reply tools failed in this outgoing session (local MCP transport
failure; send requires approval while policy reports never). Canonical delegate
IDs above are known; the collaboration API does not expose underlying thread
UUIDs. This file supplies the durable handoff even if the reply cannot be sent.

The incoming coordinator recovered the old child thread IDs: author
`01a126fc-b64c-7ae0-8ecc-9d26c242fed6`, coder
`01a126fc-dd24-7e72-b3cb-f053dabb48e8`, reviewer
`01a126fc-fbb2-7a71-b0df-07b05171608d`. Those v2 children cannot be resumed from
the incoming tree, so the incoming coordinator will create replacement roles
with this handoff and unchanged review history. The outgoing coordinator issued
interrupts to all three: author and reviewer had completed; coder was already
interrupted. No outgoing delegate remains authorized to write. Reviewer saved
`docs/edition-2-notes/enhancement-02-review.md`; coder's final validation report
was not received. Replacement coder must re-establish actual state and check
for any interrupted subprocess before continuing; do not assume pending tests
passed. Safe for takeover; no new chapter work started here.

Incoming coordinator verified the outgoing stop and restarted the three roles
as `enhancement_author_resumed`, `enhancement_coder_resumed` and
`enhancement_reviewer_resumed`, with full-access command execution. They read
the current artifacts and retain the preceding review history. Author recovery
confirmed Chapter 2 prose/evidence alignment; final grader validation and the
independent combined decision remain pending.

### Chapter 2 acceptance after recovery

The resumed team completed Chapter 2 at initial formal combined review, with
zero correction rounds. See [author notes](edition-2-notes/enhancement-02.md),
[grader evidence](edition-2-notes/enhancement-02-code.md) and
[independent review](edition-2-notes/enhancement-02-review.md). The coordinator
checked the source-relative prose diff, dispositions and untouched original
grader/frozen-source paths. Both student and legacy implementations score
100/100 under the revised grader; the student's original 92/100 remains FAIL.
Nine targeted source defects fail the intended new checks. The reviewer also
independently verified the empty-stub mutation and exact published sample.
Affected tests/vet pass. The broad root test is explicitly incomplete following
unrelated script-package setup failures; it is not reported as passing.
No product code or original source changed. Chapter 3 may now begin.

### Chapter 3 acceptance

Accepted at initial formal combined review, zero correction rounds. The original
tool-count and editing stories remain; additions clarify arguments, ownership,
credential handling and actual live use. See [author notes](edition-2-notes/enhancement-03.md),
[grader evidence](edition-2-notes/enhancement-03-code.md) and
[review](edition-2-notes/enhancement-03-review.md). Revised student grading is
100/100; legacy grading is 96/100 because its search tool omits the context
required by the original prose. Original results remain student 90/100 FAIL
and legacy 100/100 PASS. Seven targeted defects and two score-accounting cases
verify the changes; affected tests/vet and own-binary fake demos pass.
The Chapter 2 capture fake now terminates after a tool result so Chapter 3 can
use it without looping; both chapters were revalidated. The coordinator checked
the source diff and unchanged original/frozen paths. Chapter 4 may now begin,
including the original Chapter 5 external-consumer proof.

### Chapter 4 grader composition agreed before implementation

Keep the original job checks at 100 points, with corrected inherited parity.
Add 60 points for external consumer execution (25, including separate build),
star imports (25), and an emitted tool diagnostic (10). Grader-only consumer
directory/case inputs allow an arbitrary custom tool without prescribing public
API spellings or adding product commands. The external dependency must resolve
to the submitted core. New results use 160 total and preserve historical Chapter
4 and Chapter 5 scores separately.

The reviewer inspected original Chapter 5's unexported global-variable grep
check: name-based exemptions can hide mutable state, while other declarations
can falsely fail. Its flawed logger AST check likewise misses actual multi-hop
parent interfaces. Correct these claims by retaining mandatory source review
of ownership, actual composition, common-package behavior and no mutable global
application state; do not report automatic proof from names or emitted logs.
The architecture remains binding. This scope decision does not authorize a
new semantic-analysis framework or removal of valid job/consumer coverage.

### Chapter 4 acceptance

Accepted at initial formal combined review, zero correction rounds, with one
recorded inherited source-test exception. See [author notes](edition-2-notes/enhancement-04.md),
[coder evidence](edition-2-notes/enhancement-04-code.md) and
[independent review](edition-2-notes/enhancement-04-review.md). The chapter retains
the original jobs stories and incorporates the external consumer/logger lesson.
Revised results are accepted 160/160 and legacy 90/160; original results remain
unchanged. Thirteen targeted changed-behavior defects are detected. Affected
tests/vet and published examples pass. Two immediate descendant-PID existence
assertions failed in the unchanged source tests; bounded diagnostic runs captured
no failing process state. The reviewer accepted that stronger assertion as a
test limitation using source and independent group-kill evidence, without claiming
a zombie was observed or waiving group kill. No runtime or frozen export changed.
The coordinator checked the prose diff, grader, evidence and disposition. Chapter
5 (original 6) may now begin after this checkpoint.
