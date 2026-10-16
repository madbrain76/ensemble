# Edition 3 Chapter 2 — initial independent code review

**Disposition: required revisions; not accepted.** This reviews frozen candidate `74eb0265bd8bdbc7857d579704300b317b2365da`, not a moving draft or another entry audit. **Student correction rounds used: 0 of at most 3.** No author/text phase or Chapter 3 work is authorized by this review. Private reference comparisons and reviewer fixtures remain reviewer-only; relay the scoped findings, not private answer material.

## Identity and authority

- Authoritative repository: `5351583208988f8370937cd0494727645efa8b46`; no tracked changes.
- Student: `/Users/bill/projects/ensemble-edition-3-student`, clean at the candidate above; accepted predecessor/tag peel `edition-3-ch01-r1^{commit}` remains `1aad3975842d502b9144b7b9bddfd30ca7309db7`.
- Actual `book/edition-3/chapter-02.md` SHA-256: `1a67ee9e3b6ae5f0f776a30f58908c6d53b53be2d8136ef44661728a677dc2d2`.
- Submitted review SHA-256: `8817b507bb67c69a726d4e7d04769143e9e7905b36603d9ec86ceb2444140416`.
- Candidate `git archive --format=tar` SHA-256: `bff37519226b7aca26b273128cdc4e922e0c458c35234f751139ffd7178b3ad5`.
- Submitted CLI SHA-256: `1f6494edd754535ff5d1f2ca0f2663ccfa4ea236a5655dd8894b1940dc70f7e4`. Independent `go build -o "$SCRATCH/student-cli" ./cmd/ch01` is byte-identical. `go version -m` records candidate revision and `vcs.modified=false`.

Re-read the current reviewer skill and verified all twelve governing/chapter files against the unchanged authoritative revision and prior complete entry-audit reads. Recovered truncated review/source sections rather than relying on summaries. See `ch02-code-review-evidence/governing-hashes.json`, `identity.json`, and `final-integrity.json`. Current Edition 3 rules supersede historical paths/phases. OpenAI remains **Chat Completions**, not native Responses.

## Required corrections

All three are independently reproduced through an external Go module importing only public `ensemble`, with the real root/agent/engine/log. The frozen candidate was not patched. Exact tests and failing output are `ch02-code-review-evidence/public-test.go.txt` and `public-smoke.stdout`; the final test invocation exits **1**. Public construction/ownership/isolation passes in the same run.

### R1 — Gemini result names can be rebound by a later call

**Source:** candidate `internal/engine/render.go:432–449` builds one name map from the entire dialogue; `:416–424` uses it for every result. All ID-less calls share the empty-string key.

**Observed:** two ordinary sequential call/result pairs, with different names and omitted provider IDs, render both results with the *second* call's name. The earlier `first_tool` result becomes `second_tool`. This is not an ambiguous parallel-call fixture: each result immediately follows its own call. Missing provider IDs are an explicitly supported input shape.

**Basis:** current Chapter 2 `:1003–1008` says to find the matching call **earlier in the dialogue**; the candidate README `:81–86` also promises preceding-call resolution. A later call must not change an already-paired earlier result.

**Smallest required correction:** preserve chronological call/result identity when rendering Gemini and add a public-root regression for sequential ID-less calls with distinct names. Do not add a function-name field to Context, a tool executor, or a generalized matching framework. Ambiguous/unsupported input must not be silently assigned a different call's name; this finding requires no new broad pairing API.

### R2 — measured usage disappears on unusable OpenAI/Gemini responses

**Source:** `internal/engine/parse.go:129–130` and `:207–208` return before filling `ResponseEvent.Usage`; OpenAI's assignment is at `:170–171`, and Gemini's is likewise after candidate parsing. `internal/engine/engine.go` accumulates the returned usage before returning the parse error, so an early return supplies zero.

**Observed:** real HTTP test servers return valid JSON with empty `choices` / `candidates` and explicitly measured input/output `11/2`. `Ask` correctly fails, but `Engine.Usage()` is `0/0` for both vendors rather than `11/2`. No completed text turn should be committed.

**Basis:** inherited Chapter 1 `:349–353` explicitly counts measured usage despite an unusable answer; `:372–374` and `:613–614` require measured cumulative accounting. The candidate README `:71–73` makes the same promise. This does **not** require fabricating a completed response in the event log: `Context.Usage` and live engine accounting have deliberately different domains.

**Smallest required correction:** retain independently reported usage on these failure paths while preserving the error and non-commit behavior; add cross-vendor failure/accounting regressions through the real public root. Do not estimate missing usage, add retries, or turn error-only log facts into successful responses.

### R3 — exact-provenance OpenAI opaque parts overwrite one another

**Source:** `internal/engine/parse.go:159–168` preserves `reasoning_content` and `reasoning` as separate opaque parts. `internal/engine/render.go:295–305` assigns *both* destination fields for each part, including the absent/nil field.

**Observed:** parsing a Chat Completions message containing both already-supported fields and rendering its captured response to the exact same provenance retains `reasoning` but silently drops `reasoning_content`. The test follows Parse → Log.Append → Context → Render, rather than constructing private engine state.

**Basis:** current Chapter 2 `:939–947` says to carry opaque material and hand it back to its issuing model; `:1148–1150` requires no part content loss through the log/context seam. This is preservation of fields the candidate itself already accepts, not a demand to implement new OpenAI fields or native Responses.

**Smallest required correction:** keep one preserved opaque field from clearing another and add a regression for multiple same-provenance opaque parts. Preserve the existing vendor/model/surface exclusion behavior. No new wire surface or generic opaque-state framework is needed.

## Independent execution and evidence

All reviewer provider traffic used local servers; the process environment was allowlisted to `PATH`, `HOME`, `TMPDIR` with Go dependency fetching disabled. No credentials were read. Exact argv/cwd, elapsed durations, exits, target hashes, and grader check reports are retained in `ch02-code-review-evidence/`.

| Check | Actual result |
|---|---|
| Candidate root `go test -race -count=1 -timeout=180s -coverpkg=./... -coverprofile=… ./...` | exit 0; 127.272 s; includes candidate's real 62-second header and body tests |
| Candidate root `go vet ./...` | exit 0 |
| Separate required `gui/` module `go test -count=1 ./...` and explicit stub build | both exit 0; no live GUI launched |
| Explicit independent `./cmd/ch01` build | exit 0; same SHA-256 as submitted CLI |
| Root-built Edition 2 grader, `-ch 2 -json <exact submitted executable>` | exit 0, score 100/100, `passed:true`; 0.418 s |
| Root-built Edition 3 grader, `-ch 1 -json <same executable>` | exit 0, score 100/100, `passed:true`; 124.055 s; both actual 62-second probes pass |
| External public-import tests, `go test -race -count=1 -timeout=30s -v ./...` | exit 1: R1, R2 (two vendor subtests), R3; construction/isolation passes |
| Existing five mutation cases repeated in an exact scratch candidate archive | all five fail under mutation and pass after restore; harness exit 0; final scratch files equal all 17 frozen tracked files |
| Copied retained `verify-captures.py` | exit 0; offline only, zero new live calls |
| Independent retained-artifact and fresh keyless replay check | exit 0; 126 manifest hashes/sizes match; 3 fresh dumps and 18 fresh render processes match retained bytes; zero sentinel POSTs |

Both graders used the **same verified binary**, not a default historical target, a different scratch fix, or a short lifetime-test substitute. The source and grader lifetime checks both actually ran. Passing grades do not cover the three newly reproduced source failures.

### Coverage and sensitivity

Independent union of unique instrumented source blocks gives **745/839 statements = 88.796% (88.8%)**, matching the submitted claim. It includes all root-module production packages and `cmd/ch01`, excludes test-source statements and the nested GUI module, and does **not** count external CLI/grader/live processes as instrumented source coverage. Package-local 0% reports for packages exercised from external tests are not the merged result. See `reviewer-cover.out` and `coverage-summary.json`.

**Correction to the entry audit:** its claim of an inherited numeric source-coverage threshold was unsupported. The unchanged current governing text requires meaningful tests, honest coverage, and targeted mutation evidence; it sets no numeric source-coverage minimum. The student is correct on this point. The roughly 20% *teaching-comment* target in the coding skill is a different, non-padding target, not coverage. This review imposes no invented percentage gate.

The repeated five mutation cases were **`ephemera`, `cache-write`, `opaque-isolation`, `anthropic-order`, and `ref-stub`**, exactly as recorded in `ch02-code-review-evidence/mutation-rechecks.stdout`. Their failing and restored outputs are preserved individually. They do not establish exhaustive sensitivity: the normal suite and both graders still passed the paths with R1–R3. The reviewer regressions demonstrate those gaps without altering production code.

**Report-only correction before re-review:** the earlier version of this paragraph incorrectly named omitted-required-usage, lost-resolved-provenance, and unknown-event probes instead of the actual executed cases above. This corrects the report description only; no tests were rerun, original evidence was not changed, and the initial R1–R3 findings and required-revisions verdict remain unchanged. No current student-worktree state was inspected or reasserted. The prior report SHA-256 was `34c5bdcf79500fb9c44e6008351e5ed96bf7ec7642681cf2c6202fc21ed692bd`.

### Retained live evidence: accepted as evidence, not independently repeated live

Inspected the bounded launcher and verifier without opening credential settings. The retained record has **six actual live generation calls**, two human-interface turns per vendor: Anthropic Messages, OpenAI Chat Completions, Gemini `generateContent`. The three live CLI processes exited 0; captured requests/responses and measured totals correspond to those sessions. The report/manifest and executable identity are consistent.

The original **launcher exit 1 remains a failure**: its final regex assumed the wrong human usage-line form. `live-launch.txt`, `live-launch.exit`, and `live-summary.json` retain that assertion failure/false flag. The separate offline verifier compares the saved response usage and cumulative displayed totals with the correct form; the reviewer reran it on a scratch copy, not on the retained originals. This repairs the evidence interpretation, not the original exit status, and made no new model calls or source/binary changes.

Independently verified saved artifact hashes, recomputed dump equality, and reran all nine source/destination render pairs twice under a keyless environment and loopback POST sentinel. All match. Live-provider behavior was **not independently repeated**; retaining and auditing this bounded evidence is sufficient here, given reproducible source defects and no need to spend more on ordinary text calls. The six calls do not prove live signed tool-call replay, live blob handling, every provider error, or security. Relevant replay/signature duties rely on source/local fixtures and remain subject to the corrections above.

## Architecture, inherited duties, security, and simplicity

Reviewed all candidate production Go files, relevant tests, the predecessor delta, current graders, retained scripts, and the matching private first-edition Chapter 2 implementation. The old answer is a structure/teaching comparison, not an oracle or a student input.

- **Public architecture:** importable `ensemble` root is the sole external construction path; concrete descendants are in `internal/`, with meaningful narrow parent contracts. Public external construction/isolation smoke passes. `cmd/ch01` retains CLI/stdio handling, not provider/history internals. Required `gui/` is a separate module importing the public root and is only a stub. No production mocks, plugin/factory framework, config framework, or new dependency was introduced.
- **Ref/reduction:** typed kind/value validation; unsupported/legacy Ref forms refuse; append/load validate before changing history; load is transactional; nested data is copied across public boundaries. Stub reduction retains exact Ref facts while request rendering emits only inert stub text. Retained events are not rewritten by redaction. See `internal/common/common.go`, `internal/history/history.go`, `chapter2_test.go:123–254`, and the Ref mutation.
- **Replay/state/ephemera:** history owns sequence/order; timestamps are metadata. Event transitions are explicit, request facts consume pending ephemera, redactions affect projection, and usage is reduced from completed response facts. CLI capture/replay and failure persistence run through the real root (`chapter2_test.go:59–122,414–497`; `cmd/ch01/main_test.go`). Engine live accounting still needs R2.
- **IDs/provenance/opaque/order:** tests protect deterministic distinct synthesized IDs in the covered Anthropic shape, preserve issued IDs, and do not write synthesized wire IDs into facts. Exact vendor/model/surface gating and per-call Gemini signature preservation are covered by source fixtures (`chapter2_test.go:255–387`). Anthropic grouping occurs after merging and is mutation-sensitive. Those checks are useful but not universal: Gemini chronological lookup needs R1; multi-part OpenAI opaque replay needs R3. No new generalized collision or pairing contract is inferred.
- **Isolation/lifetime:** public boundaries copy histories, parts, raw JSON, and state maps; agent ownership/configuration is instance-local. Caller cancellation propagates through headers/body. The engine has no blanket generation deadline or response-header cutoff; only connection establishment is bounded. Source cancellation tests and both real-duration gate phases pass.
- **Security review gate:** inspected URL restrictions, credential handling, remote-error redaction, local/offline boundaries, capture paths, and opaque forwarding. No demonstrated credential leak or unintended provider call in reviewed paths. API keys are request headers, not log facts; untrusted remote error bodies are not copied into diagnostics; offline rendering makes zero observed POSTs. This is a source-and-test review, **not a security proof**, exhaustive fuzzing, or a live GUI audit.
- **KISS/teaching:** material growth is traceable to typed facts, validation/copy boundaries, three adapters, and actual capture/replay—not tools, media loaders, streaming, or future chapters. No architecture rewrite is requested. A simple count gives 305 full-line comments among 2,064 nonblank root-production lines (14.8%, versus 14.9% in the private comparison); this is descriptive, not a quota. Important ownership/lifetime/reduction/order reasons are documented. Some wire-field comments restate types and could be improved opportunistically; that is optional pedagogy polish, not a fourth blocker or permission to add padding.

## Limits, preserved failures, and bounded next step

Reviewer test-setup failures are separate from candidate defects: the first external smoke incorrectly called construction methods through the deliberately narrow parent interface (compile failure preserved under `initial-public-*`); the first offline helper misread the summary shape and used an unverified CLI form (its actual pre-command TypeError is preserved under `initial-offline-*`). Corrected reviewer tests are retained alongside them. A stale job-handle lookup and a few wrong-path inspection reads were tooling mistakes, not test runs or evidence of candidate failures; no results were counted from them.

No new live calls, GUI session, installs, restarts, private credential reads, pushes, commits, grader changes, manuscript edits, reference edits, or candidate edits. No unseen error branches or all-provider security guarantee is claimed. All reviewer scratch mutations are restored; frozen source/evidence identities remain unchanged.

**Next step:** return only R1–R3 and their public-contract bases to the student for **one bounded correction pass**. Require narrow real-root regression evidence and a newly frozen source/report/binary identity. Then rerun source/race/coverage, relevant sensitivity controls, both explicit graders on that same new binary including the real 62-second phases, and the affected offline replay checks. Do not rebuild graders, add native Responses, generalize deadlines/matching, expand tool/media scope, or start author/text work. Repeating unrelated live text calls is not automatically required for these local corrections; assess whether the actual delta invalidates the retained live evidence.
