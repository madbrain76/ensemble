# Chapter 6 — independent post-student text review

## Verdict: REQUIRED REVISION

One minimal prose correction is required: make the **Taking it for a spin** build recipe select the reader's actual CLI package/executable rather than unconditionally building the core module root. No code change, new API, provider campaign, additional feature, broad rewrite, or cosmetic revision is requested. The three bounded authored clarifications are otherwise satisfactory.

This is the single formal Edition 3 Chapter 6 post-student text review, not an entry audit or code re-review. Root alone accepts and closes. The author had used zero formal text revisions at review entry; root may now request the one permitted revision and verify the exact condition below without commissioning another formal review.

## Reviewed identity and evidence

- Frozen main HEAD: `b67386ab80cf2afe7985db5a32ab99ec2311f090`.
- Reviewed chapter: `book/edition-3/chapter-06.md`, **42,312 bytes**, SHA-256 `864d806c63c019ede50810ae5618b6fa4545f7028ae3abcb060722458babb4e9`.
- Baseline: **40,891 bytes**, SHA-256 `ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55`. The author's baseline, the root `/tmp/cadence-ch06-poststudent-baseline/chapter-06.md` copy, and this review's baseline snapshot match exactly.
- Author report: `docs/edition-3-notes/ch06-author.md`, SHA-256 `94af4bedf32dba09c01bea6a01117a90b9f1ffea2aa41171e3d06ded0763f2d7`.
- Author portable diff: SHA-256 `9eebf86b7aad8d0be503969803d82f455117715591677239702be09672324973`. Author manifest: SHA-256 `bdae2167a55ff3382b3dedcd2ac6e1070e7c518f84c87fd106e04fe1c61da71b`.
- Accepted student: `f595c787675e2a6343053ec42ea6c7e1e02a930e`, tree `6fb36168020a22b8993412812987338aa4b6e010`; annotated `edition-3-ch06-r1` object `84d0bce5f2525b406affa64fb1a8eb70711b1e06`, resolving to that commit. The student worktree remained clean. This is the already accepted one-code-round result, not a new acceptance by this reviewer.

Evidence directory: `docs/edition-3-notes/ch06-text-review-evidence/` (called **E** below). `E/INDEX.md` describes the exact hash inventory. It includes this external report and every evidence payload except the checksum file itself. Paths in `SHA256SUMS` are repository-root-relative. `E/pins.json`, `E/initial-state.json`, and `E/author-manifest-verification.json` establish the input identities; all **31** author manifest entries matched, with the exact expected set and safe UTF-8 text.

## R1 — make the manual CLI build recipe match the permitted layout

**Location:** reviewed chapter lines **757–764**, under **Taking it for a spin**; contrast lines **183–189** in **The exercise**. Exact excerpts are in `E/chapter-review-excerpts.txt`.

The new exercise guidance correctly says a CLI under `cmd/ch01` is valid and its directory or frozen binary is not the core module root. The later runnable recipe still says:

> Build your own CLI from the core module root

and supplies:

```sh
go build -o /tmp/ch6agent .
```

Running that literal command from a library core root does not build the nested CLI. The accepted layout is precisely such a case: the core is package `ensemble`, while `cmd/ch01/main.go` is package `main`. Thus the next `/tmp/ch6agent chat` instruction is not supplied with the required CLI by the preceding recipe. This is a concrete executable-versus-source identity mismatch, not a preference for a directory convention. It already existed in the baseline, but the bounded identity clarification leaves this second reader-facing invocation incomplete.

**Minimal correction / closure condition:** tell the reader to select the actual CLI build target when building from the core module root, and make the example runnable for a nested CLI (for example `go build -o /tmp/ch6agent ./cmd/ch01` for the accepted layout), while making clear that `.` is appropriate only when the CLI main package actually lives at the module root. An equivalent layout-neutral instruction using the actual selected CLI executable is acceptable. Preserve the subsequent `/tmp/ch6agent chat` flow. Do **not** require every implementation to adopt `cmd/ch01`, change accepted student source, change graders, or expand the exercise.

**Reason:** the reader should not receive a correct grading-path explanation and then a copyable build command that selects a different package. The public recipe and declared package boundary suffice to establish the issue; no build or behavioral rerun is necessary.

I raised this consequential ambiguity with root before finalizing. Root confirmed the accepted package declarations and agreed that this is the one concrete completeness issue in the bounded invocation clarification, not optional polish or a new contract. The full question and reply are retained in `E/root-clarification.json`. The review finding is limited to the runnable recipe above; it does not reopen code acceptance.

## Assessment of the authored changes

### Exercise/grading identities — satisfactory apart from R1's later recipe

Lines 183–200 now link directly to the current Edition 3 student grading route, distinguish CLI executable/build target, core module, Chapter 4 consumer, Chapter 5 workflow, and separate `ch07` shape exercise, and retain the Chapter 1 lifetime gate on the same CLI. The paragraph explains why this matters: executable identity is not source-module identity. It also distinguishes a declared `-stream-model` grader input from provider discovery and deliberately unlisted plain-fallback fixtures from invented capability rows.

The author properly keeps the older grader's shape discovery as history while requiring the explicit Edition 3 identity. The entry audit and grader-review reports support this distinction; no adapter or grading API revision is implied. The adjacent explanation still distinguishes incremental core observations from the external application's aggregate report. Those are different checks, not interchangeable success receipts.

### Gemini boundaries and provenance — satisfactory

Lines 445–453 clarify the already published parity obligation exposed by the original independent code review: explicit `parts` array entries delimit parts; a first entry of a later frame can continue compatible unsigned content; a trailing signature can finish an unsigned text, thought, or tool part; an attached signature closes that binding. Standalone opaque/signature-only material and signed/unsigned neighboring parts remain distinct, preserving native payload for exact-source replay.

The explanation gives the reason, not merely a checklist: appending after a signature changes its bound content, and replacing it loses the original binding. The original/R1 student reports and original/R1 independent reports support these distinctions, including legitimate trailing-signature controls. The prose does not prescribe a private fixture sentinel, implementation helper, provider-intent theory, new parser API, or universal real-provider guarantee. The existing `Parse` seam, response-local IDs, plain/SSE parity, opaque-data privacy, and provenance-gated replay obligations remain intact.

### Reliable JSON completion — satisfactory

Lines 667–670 distinguish three different facts: request admission, best-effort observations, and reliable completion. The final `assistant` record comes from the reliable private request reply; a `Post` acknowledgment or lossy observer queue is not a substitute. A busy hint supplies context to the next provider request, and its acknowledgment is not an additional completed answer.

This directly addresses the student's documented original CLI completion failure without widening pending-context lifetime or promising that a model will echo a hint in a later answer. The surrounding JSON envelope, part-delta/final correlation, tool-result separation, and distinct core/shape flags remain unchanged.

## Clarity, correctness, voice, and preserved teaching

The chapter still opens with the stalled-cursor problem and explains streaming as an ownership and commitment problem, rather than reducing it to token display. Its historical sequence and useful failures remain. The three clocks, ordered finalized response, worker/actor/display separation, source-bound private replay, honest usage, framing details, and eventual GUI requirements continue to provide reasons for the mechanisms.

The additions are specific to observed student/reviewer friction and fit the existing explanatory voice. The signature paragraph's small list of boundary shapes is justified by its binding explanation; it is not grounds for a stylistic rewrite. The reliable-completion addition is concise and teaches the distinction the student actually encountered. No new lead, motivational gloss, generic pedagogy section, API inventory, or speculative feature is necessary.

The diff is restricted to the exercise identity paragraphs, the Gemini paragraph, and the reliable-completion paragraph: **32 added / 14 removed lines, net +1,421 bytes**. The opening, history, surrounding examples, and downstream GUI rationale are unchanged. Bill's preface is unchanged. I found no other consequential clarity/correctness issue requiring a revision. Existing soft prose warnings are not a reason to gut useful explanation or force a chapter rewrite.

## Independent text checks

All checks below are text/identity checks, not renewed code acceptance.

| Check | Result and evidence |
|---|---|
| Actual lint command | Located and read `cmd/lintprose/main.go`; no guessed lint path used for execution |
| Baseline lint | `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go run ./cmd/lintprose -v docs/edition-3-notes/ch06-text-review-evidence/baseline/chapter-06.md`; exit **0** |
| Final lint | Same command with `final/chapter-06.md`; exit **0**; run sequentially after baseline |
| Lint comparability | Identical `chapter-06.md` basenames; full diagnostics match the author's diagnostics apart from input-path header |
| Lint counts | Baseline **5,288** words; final **5,515**. Negation warnings **18 → 19**; person-gap warnings **4,860 → 5,087**. Hard checks pass; warnings remain disclosed |
| Local links | All **7** local Markdown destinations resolve; no local fragment targets need separate anchor validation. `E/local-links.json` |
| Narrow diff | Reviewer portable and Git diffs match the author's exactly. `E/chapter.diff`, `E/chapter.git.diff` |
| Whitespace | `git --no-pager diff --check -- book/edition-3/chapter-06.md` exit **0**; no added trailing whitespace issue, final newline retained |
| Pins and author inventory | Main/student/tag/tree identities and frozen chapter/report/diff/manifest hashes match; author exact 31-entry set passes |
| Independent validator | `E/validate_review.py`, first attempt exit **0**; full receipt `E/text-validation-attempt1.txt` and `.exit` |
| Concurrent preservation | All **10** declared unrelated tracked publication files match root's saved observed hashes. `E/preservation.json` |
| Selected protected inputs | **22 of 23** author's selected initial input hashes match; the only difference is the permitted root-owned `ch06-execution.md` append, not author drift. Four root-owned adapter/guidance inputs match. `E/selected-protected-check.json` |

Full lint stdout/stderr and exits are retained as `E/baseline-prose-lint.*` and `E/final-prose-lint.*`; commands are in `E/lint-commands.txt`. Lint success does not erase the invocation finding.

## Evidence qualifications and limitations

I read the full current project instructions, coding/reviewer/author role overlays, Edition 3 execution/readiness/student/carryover documents, applicable historical workflow and Edition 2 carryover, voice and chapter-writing procedure. The current bounded student-first procedure governs over historical author-first/repeated-review instructions. I read the entire current chapter and baseline, full author report and both diffs, complete author evidence inventory and safe supporting command/validation diagnostics; the full original/R1 student and independent code-review reports; accepted receipt; entry and adapter review reports; and artifact-portability note. `E/read-inputs.json` records current safe input hashes. This was not a whole-book preparation pass.

Evidence is kept at its actual scope:

- All **13 paid exchanges** belong to the old `800da21d215b9ca48de41516be4773962fe38ccd` candidate, not corrected `f595c787...` Gemini bytes. There is no fresh R1 Gemini live proof here.
- The original **three Anthropic exchanges** failed the stronger final-answer marker assertion despite observed early text, hint transport and file effect. That failed assertion stays failed. The remaining **ten exchanges/cases** passed a different next-response scope, not a retroactive repair of the original result.
- Local plain/SSE comparisons and timing are fixture observations, not live equivalence, benchmarks, universal provider claims, or evidence of provider intent.
- The independent reviewer's GUI-restart recovery reconciled existing receipts and hashes; it was not a new behavioral rerun. This text reviewer did not rerun the code gates or controls accepted by root and the independent code reviewer.
- Old artifact indices retain `LOCALONLY` exclusions. The five preserved excluded files are explained by `ch06-artifact-portability.md`; this review does not relabel those old indices as rewritten, completely portable success. No private provider payload or key/settings contents were read or copied.

Reviewer setup limitations are recorded in `E/SETUP-AND-LIMITATIONS.md`. Six read-only probes used guessed author metadata filenames and returned “No such file or directory”; the actual inventory was then used. These are reviewer setup errors, not failed author artifacts or chapter checks. Capped rule and explanatory displays were recovered completely. The author's wrong lint-path attempt, two refused prewrite anchors, and off-by-one validator failure remain in the original evidence and are not counted as a formal author revision. No failed lint run or failed independent text-validation assertion occurred in this review.

No student/source/grader/manuscript/guidance edits, paid calls, behavioral tests, publication builds, installation, restart, cleanup, staging, commits, tags, or pushes were performed. Only this report and its evidence directory were written. No children were launched. All review command jobs and the clarification request finished before submission; no tests, provider work, or questions remain pending. The final report and inventory are frozen for root to inspect, join, and decide the one bounded author revision/closure.
