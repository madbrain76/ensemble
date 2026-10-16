# Chapter 3 student correction round 1

## Status and bounded scope

Student candidate for independent acceptance, **not an acceptance declaration**. R1 is the sole production correction requested after review of initial candidate `0eb6fdffe62e662b0ab96982b01287afba19ff42`. The initial commit, report (`ch03-student-review.md`), evidence (`ch03-student-evidence/`), source archive, and executable remain intact. No author-stage work, acceptance tag, push, grader/manuscript/main-application edits, subagents, or additional production fixes were made.

Revised clean commit: **`28be8bffdbc1356eaaa50a97369e59d0e0556dc3`**.

The correction makes Gemini signed **ordinary** text visible without treating a signature as evidence of hidden reasoning. It preserves exact-source signature replay and true-hidden isolation. A targeted fallback-ID collision regression was added without changing the already-correct ID implementation.

## Exact identity

All source, tests, comments, and README edits were completed **before committing**, building, and pinning the artifact, and before the paid live session and both graders. No post-pin source changes occurred.

| Item | Identity |
|---|---|
| Student cwd | `/Users/bill/projects/ensemble-edition-3-student` |
| Commit | `28be8bffdbc1356eaaa50a97369e59d0e0556dc3` |
| Parent | `0eb6fdffe62e662b0ab96982b01287afba19ff42` |
| Git tree | `8ebc664456b22591a2b0028be751d33d7aeddd64` |
| Executable, outside Git | `/tmp/ensemble-ch03-r1/student-cli` |
| Executable SHA256 | `db9a31df81f88cace5edcdde6af50dfa35ffe25e96fbc7d7585e51dd6e264dde` |
| Exact committed source archive, outside Git | `/tmp/ensemble-ch03-r1/source.tar.gz` |
| Archive SHA256 | `8efe096945bfb35cda09316c30a2956302f0275406a795de4e6d283a366558d3` |
| Build command | `go build -o /tmp/ensemble-ch03-r1/student-cli ./cmd/ch01` |
| Toolchain | Go 1.25.4, darwin/arm64 |
| Grader cwd | `/Users/bill/projects/ensemble` |
| Grader repository revision | `f7530a6643ca00bd2f59a8da60021c567cd7d174` |

Evidence directory: `docs/edition-3-notes/ch03-student-r1-evidence/` in the main repository. `identity.json` and `buildinfo.txt` record the pin. The executable embeds the revised commit and `vcs.modified=false`. The actual CLI built by the final local regression test has the **same SHA256** as the pinned live/grader binary, recorded in `full-race-coverage.txt`. This round does not need a reconstructed source/comment delta.

## Implementation and preserved boundaries

Only these production paths changed:

- `internal/engine/parse.go`, Gemini normalization around lines 240–248: `thought:true` still creates an opaque part. Absent/false thought flags with ordinary text now create a `text` part. Its signature JSON value is attached in the existing `Part.Opaque` field, with exact `Part.From` provenance.
- `internal/engine/render.go`, `geminiParts` text branch around lines 442–450: render the public text **once**, with its original attached signature only when the complete vendor/surface/model provenance matches. Cross-model Gemini keeps text but omits that signature; other supported vendor surfaces render the public text without Gemini metadata.
- `internal/common/common.go`: comments document the existing fields' signed-text use; no new field or framework.
- README: documents visible/hidden distinction and exact-source replay.

No changes to Agent, History, tool dispatch, request budgets, CLI/public configuration, generation caps, provider model names, endpoint selection, usage arithmetic, tool subprocess environment, or cancellation policy. The selected-model explanation policy remains opt-in via the existing public configuration; no universal gate, model-name exception, new flags, empty-answer policy, or streaming was added. True-hidden raw blocks retain their prior exact-source replay/isolation behavior. No opaque JSON or hidden reasoning was dumped into reports.

## Failure-first regression and sensitivity evidence

`baseline-r1-regressions.txt` preserves the pre-fix run: signed ordinary text with both absent and false flags caused blank successful final answers, falsely refused explained tool calls, and loss of public prose on cross-target rendering. The first fixture iteration also had two harness issues, subsequently corrected: a nonunique prose counter could match the system instruction, and a provenance-isolation probe touched unsigned user text. These are disclosed rather than counted as useful product failures. The final corrected tests independently fail under isolated production reversions, so the conclusion does not rest on those initial fixture mistakes.

`first-correction-tests.txt` records the focused corrected pass before pinning. The final committed full suite repeats all four new public tests:

- `TestR1SignedTextPublicAPI`: absent **and** false thought flags; real write dispatch after preceding signed visible explanation; successful feedback on the next HTTP request; nonblank signed final answer; exact usage; paired nonempty IDs; normalized signature/provenance; dump/load equality; unchanged History after rendering; same-source text/signature replay without duplication or malformed combined text/function parts; cross-model and both other vendor surfaces keep prose and exclude signatures/hidden content; exact vendor and surface equality independently exercised.
- `TestR1TrueThoughtCannotExplain`: actual `thought:true` text with signature cannot satisfy the gate; real disk remains unchanged and refusal feedback reaches the next request.
- `TestR1FallbackIDsUniqueAcrossBatches`: four missing-native-ID calls across two batches must have unique nonempty IDs, matching result order, successful execution, and stable dump/load projection. This is coverage, not a second production correction.
- `TestR1ActualCLISignedText`: real built `chat` subprocess and local HTTP fixture, for absent/false flags. The signed explanation must actually write a file; both human final answers must occur exactly once on stdout; follow-up request retains the first answer and its signature; saved/reloaded History retains both public answers. No real credentials are supplied to these fixtures.

### Eleven isolated mutations

All mutations ran in `/tmp/ensemble-ch03-r1/mutation-tree`, **not** the student working tree. Each mutation produced a meaningful assertion failure (exit 1, no compile failure), and its individually restored test selection passed (exit 0):

1. Revert ordinary/hidden classification.
2. Remove signature capture.
3. Remove provenance capture.
4. Revert signed-text rendering.
5. Weaken equality to model-only.
6. Leak the signature to another model.
7. Duplicate signed public prose.
8. Make genuine hidden text visible.
9. Collide fallback IDs between batches.
10. Collide fallback IDs within one batch.
11. Revert the complete R1 parser/renderer production change, testing API and actual CLI failures together.

`mutations/index.json` identifies each test, scratch cwd, exit statuses, exact patch SHA256, and source commit. Each mutation has its unified `.patch`, `-failed.txt`, and `-restored.txt`. These include the old classification's actual CLI blank-answer and gate failures, not just internal shape comparisons.

## Final local and grader verification

Commands and explicit outcomes are recorded in `commands.json`. All final commands below use the revised committed source; both graders use the **same pinned executable**.

| Command | Cwd | Outcome |
|---|---|---|
| `go vet ./...` | Student root | Exit 0, `final-vet.status` |
| `go test -race ./... -count=1 -coverpkg=./... -coverprofile=/Users/bill/projects/ensemble/docs/edition-3-notes/ch03-student-r1-evidence/coverage.out -v` | Student root | Exit 0, `full-race-coverage.txt` |
| `go run ./cmd/grade-edition2 -ch 3 -json /tmp/ensemble-ch03-r1/student-cli` | Main root | **passed:true, exit 0**, 100/100; `grader-edition2-ch3.{json,status}` |
| `go run ./cmd/grade-edition3 -ch 1 -json /tmp/ensemble-ch03-r1/student-cli` | Main root | **passed:true, exit 0**, 100/100; `grader-edition3-ch1.{json,status}` |

The full suite retains the prior whole-batch no-effect policy, pairing, real explained retry, ten-round budget, explicit nonselected/cleared policy, usage, and History tests. Both deliberate approximately 62-second test phases completed; both Chapter 1 grader delayed-header/body gates also completed. No lifetime run was killed as stuck.

Combined statement coverage is **90.5%**, not an acceptance threshold. `coverage.out` and `coverage-functions.txt` disclose the remaining uncovered branches (Gemini parser 86.2%, Gemini part renderer 86.5%). Focused assertions and sensitive reversions establish the R1 behavior; coverage percentage and grader scores alone do not.

## Single authorized actual Gemini human-chat run

`live-plan.md` was written before the paid call. This was **one** two-turn actual `chat` invocation, no automatic retry, no model/options/cap changes. Provider: Gemini; model: **`gemini-3.8-flash`**; default GenerateContent endpoint; default output cap 1024; selected explanation policy enabled; existing ten-round budget. The external watchdog was 240 seconds solely for observation.

- Exact argv: `/tmp/ensemble-ch03-r1/student-cli chat`.
- Exact disposable cwd: **`/tmp/ensemble-ch03-r1-chat-gemini-70_5qji4`**.
- Pinned commit and binary: exactly those in the identity table above.
- Duration: **25.371 seconds**; CLI exit **0**; watchdog did not expire.
- **Both human answers nonblank**, 1015 and 746 visible characters. Each terminal response contains one signed ordinary-text part, and each complete normalized visible answer was present on actual stdout. The second answer includes the harmless first-turn recall token.
- **11 actual calls and 11 results**, all IDs paired, no explanation refusals. All six required tools returned successfully: list_directory, read_file, search_files, edit_file, write_file, run_command.
- The Agent repaired the actual subtraction bug to addition in `add.go`, created `followup-marker.txt`, read it on follow-up, and executed real tests. Independently run `go test ./...` changed from exit **1** before repair to exit **0** afterward.
- Test, module, and README hashes remained unchanged. The recall token is absent from project sources/marker, so the follow-up evidence is not merely rereading a token stored in the project.
- `live-gemini-chat.json` contains sanitized response-part metadata, tool names/IDs/outcome metadata, terminal-answer nonblank/stdout/recall checks, source/binary identity, and artifact paths/hashes. `live-gemini-pre-test.txt` and `live-gemini-post-test.txt` preserve independent test receipts.
- The local normalized History is `/tmp/ensemble-ch03-r1-chat-gemini-70_5qji4/session.jsonl`; its hash is recorded in the live evidence. It is not copied into the report or printed. Raw stdout/stderr are represented by hashes; no opaque provider content or reasoning is emitted.

The launcher loaded only the authorized Gemini credential for client use in process memory, never into arguments, prompts, report files, or tool commands. Tool subprocess credential stripping remains covered by the full benign-canary tests. Actual tool work was confined by instruction to the disposable cwd; this is not an OS sandbox claim.

No paid Anthropic/OpenAI calls were repeated in R1. Their prior successful human/protocol runs are **prior-revision evidence**, retained with their original identities in the initial report. Current local cross-vendor fixtures and both graders passed, but those are not described as new paid live verification.

## Assistance, limits, and completion

The parent supplied the public R1 behavior and independent reproduction, identified the fallback-ID coverage gap, and authorized exactly one finite Gemini human session. No reviewer/private test implementation or historical answer was read or copied; no production patch or architecture was supplied. No subagents were launched.

This correction is not a history migration: previously stored misclassified events remain original facts. New response normalization and subsequent dump/load/replay are corrected and tested. Other retained scope limits are unchanged from the initial report: no OS sandbox, approval system, streaming, universal explanation policy, arbitrary short command deadline, or guarantee about every future provider response. Actual provider stop reasons remain outside this bounded correction. Successful observed Gemini answers replace the initial **observed live limitation for this new candidate**, not the record of the old failure.

All launched verification, mutation, grader, and live jobs completed. No further provider call or retry was made. The source and pinned binary remain unchanged; the student tree is clean. Independent acceptance remains the parent's responsibility.
