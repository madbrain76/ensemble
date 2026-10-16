# Edition 3 Chapter 2 — complete

Chapter 2 completed the student-first loop: entry gate → fresh implementation and live use → independent code review → one bounded correction round → code acceptance → student-informed author clarification → one independent text review. No required findings or exceptions remain. No Chapter 3 student or implementation has started.

## Accepted identities

- Main repository at entry/review: `5351583208988f8370937cd0494727645efa8b46`.
- Student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Accepted predecessor: `1aad3975842d502b9144b7b9bddfd30ca7309db7`, immutable `edition-3-ch01-r1`.
- Initial Chapter 2 candidate: `74eb0265bd8bdbc7857d579704300b317b2365da`.
- Accepted Chapter 2: **`89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`**, immutable **`edition-3-ch02-r1`**. Clean source tree; one of three allowed code-correction rounds used.
- Reviewed/final chapter SHA-256: `04b112390fac40d3dc12c2b58533573faffc717b997dffc9ce94494c97a487da`.
- One formal text-review pass, zero author revisions. [Acceptance identity](ch02-accepted.json) pins source/evidence and lifecycle facts; [execution record](ch02-execution.md) preserves chronology, errors, assistance and decisions.

## Entry gate and student separation

The [entry audit](ch02-entry-audit.md) accepted composition of existing graders, not a new framework or a Chapter 2 mode in the Edition 3 CLI. Both commands run from the book repository root against the **same explicitly built executable**:

```sh
go run ./cmd/grade-edition2 -ch 2 -json /absolute/path/to/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/student-cli
```

Require `passed: true` and successful exit from both, not score alone. The unchanged historical Chapter 2 control scored 100 while failing the actual delayed-header/body lifetime gates; a scratch control removing only the blanket timeout passed both graders. These are setup controls, not student acceptance. Historical solutions and graders were not edited.

The fresh student received the current chapter, accepted Chapter 1 predecessor, role/public guidance, exact grader invocation and scoped provider-access instructions—not reference solutions, grader internals or private repair answers. The saved child history started fresh with no inherited memory; this is instructed/context separation, not an OS sandbox claim.

[Provider readiness](ch02-provider-readiness.md) verified actual minimal calls for the three taught interfaces. Parent initially misstated Chapter 2 as requiring native OpenAI Responses. Exact-file inspection corrected that: Chapter 2 uses Chat Completions. The unnecessary Responses setup call remains recorded separately; an additional minimal Chat Completions call closed the actual interface check. A report's wrong repository revision was corrected against saved tool output. Neither mistake became a student requirement.

The student's cache-write question was answered from the current official Chat Completions reference, which explicitly documents `usage.prompt_tokens_details.cache_write_tokens`; the assistance and primary-source check are retained. It was not answered with grader implementation code or fixture values. The separate `gui/` module is an import/construction smoke stub, not an early GUI implementation.

## Implementation and acceptance evidence

[Initial student report](ch02-student-review.md), [initial independent review](ch02-code-review.md), [correction report](ch02-student-review-r1.md), and [final code acceptance](ch02-code-review-r1.md) retain exact artifacts, commands, outcomes and limits.

Both final grader reports passed, each 100/100, on one executable; the inherited lifetime gate actually exercised both 62-second phases. Independent review rebuilt a byte-identical accepted binary, reran both graders and the full race/coverage suite, repeated the former public regression failures, and individually reverted all three fixes in scratch to prove regression sensitivity. Final statement coverage was 749/838 (89.4% displayed). The current governing files impose no numeric threshold; the auditor's earlier numeric-threshold claim was withdrawn.

The first review found three contract violations missed by the passing score: chronological association for an ID-less Gemini result, retaining measured usage on an unusable answer, and preserving multiple supported OpenAI opaque fields. Round 1 fixed those without a redesign. The independent re-review closed R1–R3 with no further findings. A mistaken mutation-summary paragraph in the initial review was corrected against its actual logs; its findings and underlying test outcomes were unchanged.

Six real CLI generation calls—two conversational turns per vendor—belong to initial candidate `74eb026`. Parent directly read their stdout/stderr; meaningful follow-up and cumulative usage are visible. **The corrected candidate was not live-retested.** Its deltas were verified with affected saved-log replay, keyless/local CLI/error-cost cases, source tests, mutation controls and independent source review. No fake/synthetic signature or media case is presented as live-provider proof. Original source, failure artifacts and live provenance remain preserved rather than relabeled as new-binary observations.

Five local test executables are excluded from Git; [artifact exclusions](ch02-artifact-exclusions.json) records paths, sizes and hashes. Source/build metadata and textual receipts remain versioned. Historical manifests describe the artifacts present at acceptance; a fresh clone must rebuild/materialize those executables before checking every recorded binary path.

## Bounded book improvement

[Author report](ch02-author-review.md) and the [verbatim independent text review](ch02-text-review.md) establish the accepted three-hunk edit:

- Locate the optional Chat Completions cache fields without changing normalization or historical evidence.
- Build the actual nested CLI, retaining the root-main alternative, and run both existing graders against the same explicit executable.
- Make the live-build example consistent with that layout without changing its exercise.

Net +154 split words; 22→23 balanced fenced blocks. No new architecture, features, defaults, media/streaming/tool exercise, future Responses requirement or test-inventory lecture. The existing text already taught the obligations behind R1–R3, so those student testing failures did not justify new requirements. Original historical failures, rationale, voice and evidence qualifications remain intact; Bill's preface is unchanged.

The author delegated the single independent text review to a fresh read-only child. Parent had separately prepared another reviewer for prerequisite reading only; once that overlap was discovered, parent verified the completed nested review rather than launch a second pass. The nested reviewer passed the exact frozen chapter. The author's final report appends that verdict; its review-time input hash is separately retained. Saved final turns were verified before joining/closing all workers, including the nested reviewer and the unused readiness-only worker.

## Closure and next step

The shared execution status and public student brief now agree with the accepted chapter and paired invocation. Checkpoint scope is explicit: this chapter, those two guidance files, and the new named Chapter 2 records/evidence only. Unrelated files remain untouched. The working application, historical solutions/graders, original manuscript and immutable final reference were not edited. No installation, restart, publication build, deployment or push occurred.

**Next is Chapter 3 entry readiness**, not an assumed grader pass: reload its current contract and governing instructions, align/verify its grader and inherited obligations, then launch a fresh student from `edition-3-ch02-r1` with permitted inputs only. No speculative author-first pass or all-grader rewrite is authorized.

### Final verification qualification

The exact checkpoint list contains 501 named files (mostly retained evidence), excluding five local binaries. Protected checks verify all 231 reference and 47 original manuscript files; 34 local closure links resolve. The author report was normalized only by removing one extra EOF blank line after review, with original/final hashes retained in the acceptance receipt. Thirteen whitespace warnings in raw build metadata, CLI prompts and the public documentation excerpt are preserved and classified in the execution record rather than silently changing evidence. Chapter/current guidance/closure prose pass independently; credential checking remains enabled. External URLs and anchors were not comprehensively validated, and no rendered publication was built.

The orchestrator ran the required repository prose linter on both the baseline and accepted Chapter 2 (same basename), preserving the [full receipt](ch02-prose-lint.json). Both exited 0 with all hard checks passing. Prose words were 9,793→9,932; existing length/person-gap warnings remain, and negation forms 19→21 add a soft density warning. The reviewed clarifications use those contrasts to prevent demonstrated mistakes; no extra author revision or cosmetic rewrite was opened. This final check added one explicit receipt to the earlier 500-file scope.
