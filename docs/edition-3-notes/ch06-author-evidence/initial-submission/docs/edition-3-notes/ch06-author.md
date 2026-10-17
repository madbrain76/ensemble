# Chapter 6 post-student author report

## Status and bounded scope

Initial Edition 3 post-student teaching draft, ready for the root-owned **one formal text review**. This is not code review, code acceptance, publication, or text acceptance. No formal text-review revision has been used by this author. The two refused pre-write edits and the local validator correction below are drafting/setup history, not formal author revision rounds.

Accepted student code remains `f595c787675e2a6343053ec42ea6c7e1e02a930e`, tree `6fb36168020a22b8993412812987338aa4b6e010`, immutable tag `edition-3-ch06-r1`. The acceptance record says one code correction round was used and code workers are closed. I did not inspect or modify student production source, run student tests, redo acceptance, or launch another audit.

Author writes are confined to:

- `book/edition-3/chapter-06.md`
- this report, `docs/edition-3-notes/ch06-author.md`
- `docs/edition-3-notes/ch06-author-evidence/`

Main HEAD was inspected as `b67386ab80cf2afe7985db5a32ab99ec2311f090` and was not changed. No children/reviewers were spawned. No installation, publication, push, staging, commit, tag, paid calls, key reads, GUI work, cleanup, or solution-source reads were performed. The parent supplies the later formal text review; this author has not launched it.

## Baseline and final bytes

Before editing, the current chapter matched the supplied `/tmp/cadence-ch06-poststudent-baseline/chapter-06.md` byte-for-byte:

- Baseline: **40,891 bytes**, SHA-256 `ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55`.
- Final: **42,312 bytes**, SHA-256 `864d806c63c019ede50810ae5618b6fa4545f7028ae3abcb060722458babb4e9`.

Portable copies are `ch06-author-evidence/baseline/chapter-06.md` and `ch06-author-evidence/final/chapter-06.md`. `chapter.diff` uses those same-basename paths; `chapter.git.diff` is the narrow repository diff. The net change is 1,421 bytes, not a general rewrite. The final `SHA256SUMS` inventories every author evidence payload plus this report and the chapter, excluding only the inventory's own self-referential hash. Its digest is supplied in the submission receipt.

## Why these three edits

1. **§6.1, grading identities.** The entry audit and grader review demonstrate that an executable path, core module root, inherited consumer/workflow, and new streaming exercise are different identities. Nested `cmd/ch01` architecture is valid. Replaced the old Edition 2 shortcut with a concise link to the current public Edition 3 route, named the separate inputs and same-CLI Chapter 1 gate, and explained why the distinction matters. The adjacent model sentence now requires an explicit `-stream-model` for Edition 3 while retaining `claude-sonnet-5` as the historical reference row. Older `ch07` discovery is identified as older behavior, not a substitute for explicit exercise input. Optional case/media inputs and the distinction between core observations and the exercise summary remain. The public brief contains the complete command and map schema; duplicating it here would add maintenance burden without teaching a new idea.

2. **§6.5, Gemini part ownership.** The original independent review and accepted R1 report identify a single demonstrated boundary/provenance defect. Added the reasoning beside the existing vendor-specific decoder explanation: array entries establish boundaries; a first entry across frames can continue compatible unsigned content; a trailing signature can finish unsigned text, thought, or tool content; attaching a signature closes its exact binding. The paragraph preserves standalone opaque-first material, signed text followed by unsigned or signed text, and signed thought followed by separate opaque material. Native exact-source replay remains the outcome. Existing §6.4 already teaches response-local IDs, valid zero IDs, parse reuse, and actor response identity; it was left unchanged rather than duplicated. This teaches the public contract's consequences, not a private fixture checklist or a new parser interface.

3. **§6.10, JSON completion.** The original student report records the fix from observation-driven JSON completion to reliable request replies. Existing §6.11 already says lossy observations are not reliable completion, and the chapter already distinguishes hint acknowledgments. Added only three sentences at the JSON `assistant` record: produce the answer from the reliable private request reply, not `Post` admission or a lossy queue, and leave busy hints as next-request context rather than completed answers. This connects existing principles at the implementation point without adding an API or duplicating the actor/reply chapter.

No change was needed to the opening, thesis, signature/API design, PartID rationale, three-clock explanation, streaming-disabled path, inherited job/lifetime obligations, historical mutant failures, CLI flush story, exercise measurements, live-run history, or Chapter 7 boundary. Bill's preface was not edited. No provider-compliance guarantee, new benchmark, provider client, controller, interface, or future-chapter obligation was added.

## Evidence and assistance limits

The accepted reports, not a new source review, supplied the demonstrated needs. The original independent reviewer explicitly found the public contract already sufficient; these are teaching clarifications, not retroactive requirements for acceptance.

All **13 paid exchanges belong only to original `800da21`**, not the corrected Gemini production bytes. The original Anthropic answer-marker assertion remains failed even though early streaming, next-request hint transport, and file effect succeeded. The later remaining-only 10-call run passes its own scope; it does not repair that assertion or prove the corrected Gemini implementation live. R1 used no paid calls. No new live proof, live streaming/plain equivalence, or speed benchmark is claimed here. The full existing reports and failure receipts were not rewritten, trimmed, or copied into a new apparent acceptance story. Their setup errors and reviewer GUI restart recovery remain evidence history, not added chapter lessons.

The parent provided this bounded assignment and a scope acknowledgment, and identified its portability note/acceptance/execution work as parent-owned. No other author/reviewer agent assisted. The public brief was already current; I did not edit it or the grader. Solution link targets were checked for existence only, not used as starting code.

## Actual read coverage

`ch06-author-evidence/read-inputs.json` records hashes and sizes for the files actually read in this session. These are reads, not inherited handoff memory:

- Full `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`, and `.agents/skills/edition-3-author/SKILL.md`; the coding skill's ending was explicitly recovered.
- Full mandatory `docs/agentic-codebook-workflow.md` (sequential ranges, with its final 19 lines explicitly recovered after confirming the actual 2,434-line count), `docs/edition-2-carryover.md`, `book/voice.md`, and `book/chapter-writing-procedure.md`. A combined voice/procedure output was truncated; each was then read separately in full. Historical author-first instructions were treated as history; current Edition 3 student-first governed.
- Full current `docs/edition-3-execution.md`, `docs/edition-3-readiness.md`, `docs/edition-3-carryover.md`, and public `docs/edition-3-student.md`.
- Full current `book/edition-3/chapter-06.md`, including bounded rereads of its actual line-wrapped text and recovery of the exercise section omitted from one inline reader result. Only this chapter was prepared.
- Full `ch06-student-review.md`, `ch06-student-review-r1.md`, `ch06-code-review.md`, `ch06-code-review-r1.md`, and `ch06-accepted.json`; acceptance was reread near freeze. Full `ch06-entry-audit.md` and `ch06-grader-review.md`, used only for the demonstrated grading identity gap, not to redo their work.
- `Makefile` and `cmd/lintprose/main.go` to discover the actual lint command and rules. No separate accepted student source inspection was necessary.

## Validation commands and results

Commands below run from the repository root. Full verbose lint output and explicit exit files are retained; no diagnostics were discarded for presentation.

| Command | Exit | Result / receipt |
|---|---:|---|
| `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go run ./cmd/lintprose -v docs/edition-3-notes/ch06-author-evidence/baseline/chapter-06.md` | 0 | `baseline-prose-lint.txt` and `.exit`; all hard checks pass |
| `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go run ./cmd/lintprose -v book/edition-3/chapter-06.md` | 0 | `final-prose-lint.txt` and `.exit`; all hard checks pass on final bytes |
| `python3 docs/edition-3-notes/ch06-author-evidence/validate_author.py` | 0 | `text-validation.txt` and `.exit`; baseline SHA, narrow change locations, final UTF-8/whitespace, seven local links, protected hashes |
| `git --no-pager diff --check -- book/edition-3/chapter-06.md` | 0 | `chapter-diff-check.txt` (empty diagnostics) and `.exit` |
| `git --no-pager diff -- book/edition-3/chapter-06.md` | 0 | Complete diff saved as `chapter.git.diff`; inspected for the three local changes |

Same basename matters: both lint inputs are named `chapter-06.md`. Baseline/final lint word counts are **5,288 / 5,515**, within the actual tool's 4,000–7,500 range. Neither has a length warning. Both retain soft warnings for negation forms (18/19) and person gap (4,860/5,087 words). Hard checks, including zero em-dashes and zero LLM crutches, pass. These soft warnings do not justify stripping unrelated teaching or artificially changing the chapter's voice. No invented quality score is assigned.

The validator resolves baseline links in the original chapter directory, not the relocated evidence directory. All six baseline local links and all seven final local links exist. The new brief link intentionally has no nonexistent heading fragment. No network link checks were made.

### Setup failures retained

- The initial guessed `scripts/lint_prose.py` path did not exist (`wc`/`cat` failed with exit 1). Inspecting `Makefile` and the repository located the actual `cmd/lintprose` tool before invoking it.
- The first Python exact-anchor edit refused with `AssertionError` (exit 1) before writing the chapter. Baseline/state receipts had already been saved. I recovered the exact wrapped chapter text and then used its actual anchors.
- The first text-validator run exited 1 because its own expected insertion range began at baseline line 440 rather than the actual boundary at 439. Full failure diagnostics and the original validator survive as `text-validation-attempt1.txt`, `.exit`, and `validate_author-attempt1.py`. Links, baseline SHA, and chapter whitespace passed in that run. This was not a chapter or code defect.
- A later exact replacement of the old discovery sentence also refused before writing (Python exit 1), because the actual sentence continued into inherited parity. The subsequent commands in that compound shell still ran: the unchanged validator repeated exit 1, lint passed, and diff-check passed; the shell's final exit 0 is not being represented as success of the refused edit. The precise sentence was then replaced once, preserving its remaining rationale. The validator's expected range was corrected locally. No manuscript broadening or test weakening was used to resolve it.
- Intermediate-draft lint passed and remains in `final-prose-lint-attempt1.txt` with its exit file. Its chapter bytes remain in `draft-before-route-discovery-clarification/chapter-06.md`; it is not the frozen final.

## Preservation and handoff

`protected-before.json` and `protected-comparison.json` cover the four root-owned tracked adapter/guidance edits, ten unrelated tracked publication edits, Bill's preface, selected acceptance/review inputs, and the parent-owned chapter execution log. All selected hashes match except `docs/edition-3-notes/ch06-execution.md`, which the parent owns and may append concurrently. I did not rewrite or revert it. The initial/final status receipts preserve the wider dirty-tree context; untracked historical artifacts were not cleaned. This is a selected-file scope check, not a claimed new full audit of all historical evidence.

The evidence directory contains only UTF-8 text, including Markdown snapshots, scripts, JSON receipts, diffs, and lint logs. No executable, archive, credential, or private provider trace was copied into it. The final freeze check (`python3 docs/edition-3-notes/ch06-author-evidence/freeze_author.py`, exit 0; `freeze-check.txt` and `.exit`) verifies the report/evidence text and generates the hash inventory. `shasum -a 256 -c docs/edition-3-notes/ch06-author-evidence/SHA256SUMS` verifies that inventory at submission (exit 0). At submission all author tool jobs are finished, no children exist, and no worker is pending. The next action belongs to root: review this exact frozen delta with its one formal text reviewer, then request at most the authorized author revision if needed.
