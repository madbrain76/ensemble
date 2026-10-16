# Chapter 3 code-correction Round 1 re-review

## Verdict: PASS — R1 closed; code-review gate passes

Revised candidate: **`28be8bffdbc1356eaaa50a97369e59d0e0556dc3`**, parent `0eb6fdffe62e662b0ab96982b01287afba19ff42`, in `/Users/bill/projects/ensemble-edition-3-student`.

The sole blocking finding from the initial review is corrected. Ordinary signed Gemini text is now visible, can satisfy the selected-model explanation rule when it precedes the first tool call, and survives cross-model conversation projection as ordinary text. Its signature remains attached to the original part and is replayed only to the exact issuing vendor/surface/model. Genuine hidden thought remains hidden and cannot authorize a tool batch. Independent regression tests, actual pinned-CLI controls, scratch reverts/restorations, the full race suite, both graders, and the bounded live-evidence audit support this conclusion.

**No blocking introduced defect found. No further correction is requested by this review.** This is closure of R1 against the frozen revision, not a new whole-chapter requirements exercise. The initial whole-chapter review otherwise stands. Root owns final acceptance, tagging, and the author loop; this reviewer does not perform those actions.

## Frozen inputs and scope

- Reference HEAD is unchanged at `f7530a6643ca00bd2f59a8da60021c567cd7d174`. Retained governing/reference hashes from the entry audit match. The earlier full governing reads remain valid.
- Verified all **21 source hashes** and **56 new evidence hashes** in `ch03-candidate-r1.json`, and all **146 original evidence hashes** in `ch03-candidate.json`. The revised archive contains exactly the 21 tracked files with matching contents. Student worktree is clean at the pinned revision.
- Read the full `ch03-student-review-r1.md`, recovering the small section omitted by combined output. Read the complete 415-line `chapter3_r1_test.go`, all 11 exact mutation patches/index, and the production/documentation diff.
- The corrected initial review remains SHA-256 `a3e3264b8e2862c05b5fbe66289e6f9b80d43435e29886ea9615cf6b491bdf00`. Its original evidence manifest's pre-transcription-correction report hash is qualified by the saved correction receipt, not silently rewritten.
- Current explicit CLI: `/tmp/ensemble-ch03-r1/student-cli`, SHA-256 **`db9a31df81f88cace5edcdde6af50dfa35ffe25e96fbc7d7585e51dd6e264dde`**. Archive SHA-256: **`8efe096945bfb35cda09316c30a2956302f0275406a795de4e6d283a366558d3`**. `go version -m` independently reports the revised commit and `vcs.modified=false`. Unlike the initial revision, there is no live/comment-only final delta in this R1 attribution.
- Production behavior changes are confined to Gemini ordinary-text parsing and rendering. Other edits are field comments, README, and the new test file. No Agent, History, tool, CLI, option, timeout, retry-budget or model-selection behavior was changed.
- Tests/mutations ran in separate temporary archive copies, not the student tree. Both graders and the independent CLI matrix used the exact manifest executable. Go tests compile their own test/scratch executables; this is not a claim that those have byte-identical build metadata to the committed binary.
- Only this report and `ch03-code-review-r1-evidence/` were written in the reference repository. No student contact, source/grader/manuscript edits, installs, subagents, commits, pushes, runtime restarts, or reviewer-paid provider calls. No input or scratch evidence cleanup was performed.

## Public-contract closure and source review

The basis remains Chapter 3's ordinary synchronous final-answer flow, same-response visible nonwhitespace explanation **before the first call** for explicitly selected policy, and inherited Chapter 2's portable public text versus exact-source opaque replay. This does not introduce a universal explanation rule, model-name inference, a new flag, a generic empty-answer policy, historical-log migration, or Chapter 6 arrival-order requirements.

### Parser: `internal/engine/parse.go:240–249`

The old condition classified any signed ordinary text as opaque. The revised condition treats a part as opaque when it is genuinely `thought: true` or lacks text. Ordinary text with an absent/false thought flag becomes `Type: "text"`. When it has a signature, the existing per-part `From` and `Opaque` fields retain exact provenance and signature separately from its visible text.

This is the relevant semantic distinction: signature presence is replay metadata, not proof that the text is hidden. It fixes both the empty terminal-answer symptom and the false refusal of signed ordinary pre-tool prose without changing the Agent's rule. It also restores ordinary text's portability. Hidden parts still take the opaque branch; tool calls retain their separate pre-existing path.

### Renderer: `internal/engine/render.go:442–449`

The new explicit Gemini text branch emits **one** text part. It attaches the retained signature only when `p.From == to`, comparing the full provenance value rather than model name alone. Rendering elsewhere retains the ordinary text without its signature; other vendor renderers already project ordinary text without opaque metadata. Same-source hidden material and signed tool-call behavior remain in their existing branches.

No duplicate public prose is introduced. Dump/load retains the per-part signature/provenance association through the existing History representation. The implementation does not broaden model equivalence, expose raw opaque JSON/reasoning, discard all signatures, or add a parallel history owner. Comments accurately describe this use of the existing fields; no new serialized field or public seam was needed.

### Required boundaries checked

- Absent and explicitly false thought flags: signed ordinary pre-tool text authorizes actual dispatch; signed final text reaches the caller/CLI.
- True thought flag: even nonblank signed hidden text cannot authorize a batch; refusal has no write effect, preserves pairing, and a later ordinary explanation can dispatch.
- Exact vendor/surface/model replay preserves the original visible-part signature association, including after dump/load. Changed model, changed vendor, changed surface and cross-vendor projection do not expose the signature/hidden material. Public text remains portable.
- Native signed tool-call behavior remains covered. Usage still accumulates measured facts; the signed-text fixture checks four measured input/output tokens over its two requests.
- Fallback IDs remain distinct within and across batches, with one matching result per call and dump/load/render pairing. The new coverage closes the initial review's test gap without modifying fallback-ID production code.
- Existing policy true/false/unset, invalid-configuration failure, whole-batch no-effect refusal, late/blank/opaque/prior-prose rejection, real explained retry and the shared ten-round budget remain intact. No opaque signature is printed as the answer.

## Independent runs: actual results

All owned jobs finished. Logs, JSON outcomes, commands, mutation transcripts and hashes are saved in `ch03-code-review-r1-evidence/`.

| Run | Measured outcome |
|---|---|
| `go test -race -count=1 ./...` on pristine archive | **Exit 0**. Root package 126.634 seconds; CLI package and internal packages pass/compile. Includes inherited generation-lifetime tests and the new R1 tests. |
| `go vet ./...` on the same pristine copy | **Exit 0**, no diagnostics. |
| `go run ./cmd/grade-edition2 -ch 3 -json /tmp/ensemble-ch03-r1/student-cli` | **Exit 0, passed true, 100/100**; all ten checks, including the 13-check Chapter 2 parity gate. |
| `go run ./cmd/grade-edition3 -ch 1 -json /tmp/ensemble-ch03-r1/student-cli` | **Exit 0, passed true, 100/100**; all nine checks. Real **62-second header and 62-second body** phases completed with full answer/usage. |
| Initial independent public-API fixture, unchanged | **Exit 0**. Signed ordinary text (absent/false thought), signed explanation and answer, hidden/provenance boundaries, and ID-less-call pairing now pass. This is the same saved fixture that exposed R1 initially. |
| Independent actual-CLI local matrix | **Exit 0**; ten policy/budget cases and two signed-terminal cases pass. Same current pinned binary, synthetic providers only. |
| All 11 exact student scratch patches, independently applied | Every mutated run **exit 1 with named test assertion failures**, no compile failures; every restored run **exit 0**. Whole production reversion also runs the independent reviewer fixtures. |

The ordinary race suite remained pristine while the API fixtures and mutations used separate copies. Supported provider-key variables were removed from grading/race/mutation process environments; local fixtures supply synthetic credentials. No test requests were relabeled as live-provider evidence.

### Actual-CLI matrix detail

Reused the initial review's saved driver, changing only its scratch path, revision manifest, and Gemini expectation from the known empty-answer failure to two visible answers. The exact driver is retained as text.

- Policy true with missing, whitespace-only, opaque-only, late or prior-turn prose refuses all six calls before disk effects, records six matching error results, then accepts an ordinary explained retry that actually executes and returns an answer. At the retry request boundary the protected seed and forbidden write/command markers demonstrate no refused-batch mutation.
- Proper explanation succeeds directly. Policy false and unset permit ordinary tool dispatch for the same arbitrary model string. Invalid policy exits nonzero without a provider request.
- Repeated invalid batches stop at exactly ten requests and sixty paired refusals, with no write/command effect.
- Gemini signed terminal text is visible **twice** through both `--protocol` and `chat`, with exit 0, two requests, and no synthetic signature appearing in stdout. These actual-binary cases failed in the initial review and now pass.

The new student R1 actual-CLI test additionally exercises signed pre-tool text, native call signatures, same-source replay, a successful command, and a signed final answer. Its source was reviewed, it passes in the full suite/restoration runs, and parser/render reversions cause its behavioral assertions to fail.

### Independent sensitivity of the bounded correction

The reviewer read, hash-checked and applied each saved patch independently, checked exactly which source files changed, restored those files after the mutation, and verified every candidate source file matched again. These are actual independent executions of the **11 R1 patches**, unlike the narrower reconstruction used for the original revision's 17-control evidence audit.

| Exact mutation | Boundary detecting it |
|---|---|
| `classification-reversion` | Visible ordinary text / actual CLI answer |
| `signature-capture-removed` | Exact signature retention |
| `provenance-capture-removed` | Original-part exact provenance |
| `render-reversion` | Same-source signature replay, including CLI wire |
| `exact-provenance-weakened` | Vendor/surface isolation, not only model equality |
| `signature-cross-model-leak` | Signature must not follow portable public text across models |
| `duplicate-signed-prose` | One visible part, no duplicate answer text |
| `hidden-text-made-visible` | True thought cannot become answer/explanation or cause a write |
| `fallback-ids-collide-between-batches` | Cross-batch uniqueness and result pairing |
| `fallback-ids-collide-within-batch` | Within-batch uniqueness and result pairing |
| `whole-R1-production-reversion` | New student tests **and unchanged independent reviewer regressions** fail; restoration passes |

Fallback-ID patches are coverage controls against unchanged Agent code, not evidence that this correction changed it. No surviving mutation or unexplained baseline failure remains in this R1 run. Passing sensitivity is meaningful evidence of these boundaries, not a claim of exhaustive proof.

## Bounded live-evidence audit

The one authorized Gemini session was conducted by the student launcher, not by this reviewer. Its plan, sanitized receipt, pre/post test transcripts and artifact pins were read. The receipt identifies the current exact-commit binary, Gemini **`gemini-3.8-flash`**, GenerateContent, `LLM_REQUIRE_TOOL_EXPLANATION=true`, unchanged 1024 output-token cap and ten-round budget. It records **25.371 seconds**, exit 0, no watchdog expiration, and no retry/cap/model change. No further paid session was requested or made.

Independent checks against the retained local artifacts:

- All six History/project artifact hashes match the receipt. `add_test.go`, `go.mod` and README remain unchanged from their initial hashes; `add.go` now contains addition; the marker exists. The retained independent test receipts show the original `Add(7,5)=2` failure (**exit 1**) and post-repair success (**exit 0**). This review did not rerun the live project tests or network session.
- History has **11 calls and 11 corresponding results**, with ordered ID/name matches. All six tool names have successful returns. There is **one recoverable `search_files` error**, followed by successful search; this is not eleven uniformly successful results and is not an explanation-policy refusal. Tool errors remain usable feedback.
- Terminal History events **33 and 42** each contain one ordinary signed text part with exact response provenance and no calls. Independently extracted ordinary text is nonblank and has **1015** and **746** characters respectively, matching the receipt. Both texts were actually read: the first coherently describes inspection/repair/test verification; the second describes the marker, repair and later test result.
- The same harmless recall marker is present in both ordinary terminal answers and absent from all five retained project files. This independent check supports conversation preservation. The launcher's intended original prompt-token validation remains its recorded observation; the reviewer did not inspect other private History text.
- Only ordinary text from those two terminal events was extracted for reading. No opaque content, hidden reasoning, signature values or credentials were printed or copied into the review artifacts. The durable review evidence records lengths/hashes and booleans, not the live answers themselves.

**Raw stdout and stderr were hashed but not retained.** They were therefore **not reread or independently hash-verified from original bytes** by this reviewer. The launcher recorded checking both complete ordinary answers on actual stdout in process. That contemporaneous observation, retained normalized answers, source/binary identity and the independent actual-CLI regression are complementary evidence; none is falsely presented as a new review-time live observation. `live-audit.json` explicitly separates launcher observations from independently read artifacts.

The independent verifier initially asserted that every live tool result was successful and stopped. Inspection of metadata exposed the recoverable search error; the verifier was corrected to the actual recorded behavior and rerun without any live request. That was a review-script assumption error, not a product failure or a hidden live retry.

Prior Anthropic/OpenAI live evidence remains attributable to the initial revision. It was not rerun or relabeled as R1 evidence. The narrow Gemini-only production change, preserved prior evidence, independent three-vendor graded regressions and this single R1 Gemini session are the bounded closure evidence. The parent's additional local-provider smoke is separate orchestration evidence, not a paid/live proof and not needed to substitute for these reviewer runs.

## Limits, handoff and disposition

- Scores alone remain insufficient: the independent signed-text, exact-source, hidden-text, CLI and sensitivity checks are what close the initial known gap.
- No full-book run, Chapter 6 streaming-order test, provider discovery campaign, new historical-log migration behavior, or generalized empty-provider-answer policy was required or attempted.
- The original input evidence, initial review/correction receipt and R1 candidate remain frozen. Scratch files and review receipts are retained; the final handoff text is saved before any cleanup. Final checks verify source/binary hashes and repository cleanliness again.
- **R1 is closed on `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`; this code-review gate passes.** Root may now apply its acceptance/tagging and author-loop procedure. This review does not itself tag, launch an author, or expand student scope.
