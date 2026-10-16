# Edition 3 Chapter 2 — student implementation review

## Candidate (submitted for independent review, not self-accepted)

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Exact candidate commit: **`74eb0265bd8bdbc7857d579704300b317b2365da`** (`Implement Chapter 2 event history and three-provider context seam`).
- Accepted predecessor/tag peel remains **`1aad3975842d502b9144b7b9bddfd30ca7309db7`**, `edition-3-ch01-r1^{commit}`. Built on that source, not an empty tree or copied reference. No new accepted tag, reset, clean, push, or Chapter 3 work.
- Authoritative input revision supplied: `5351583208988f8370937cd0494727645efa8b46` at `/Users/bill/projects/ensemble`.
- Source archive SHA-256 (`git archive --format=tar HEAD`, student **root** cwd): **`bff37519226b7aca26b273128cdc4e922e0c458c35234f751139ffd7178b3ad5`**.
- Actual CLI: `/Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-evidence/student-cli`.
- CLI SHA-256: **`1f6494edd754535ff5d1f2ca0f2663ccfa4ea236a5655dd8894b1940dc70f7e4`**.
- Build command, student root cwd: `go build -o /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-evidence/student-cli ./cmd/ch01`.
- Build identity: Go 1.25.4, darwin/arm64, `vcs.revision=74eb0265bd8bdbc7857d579704300b317b2365da`, **`vcs.modified=false`**. Full identity in `final-identity.txt`.
- Worktree clean after final testing; hashes rechecked in `final-source-checks.txt`.

Evidence paths below are relative to [`ch02-student-evidence/`](ch02-student-evidence/). Evidence/report files are deliberately **not staged or committed in the authoritative repository**.

## Inputs and assistance

Read the authoritative `AGENTS.md`, ensemble-coding skill, Edition 3 coder skill/shared procedure, `docs/edition-3-student.md`, entire actual `book/edition-3/chapter-02.md`, and accepted predecessor source/tests. Recovered truncated chapter reads. The public student brief superseded references to historical solutions/carryover; no solutions, future chapters, grader internals, private provider receipts, repair records, or other agents' data were consulted. No sub-agents were used.

Explicit assistance received:

1. The assignment's launch clarification: the accepted root is a library; build **`./cmd/ch01`**, not `.`. Execute both graders from the authoritative root against the exact same absolute binary. Edition 3 CLI supports Chapter 1 only.
2. The assignment's credential authorization and time-bound economical model guidance: `directClaudeAPIKey`, `directOpenAIAPIKey`, `directGeminiAPIKey` from `~/.cr/settings.json`; suggested callable IDs `claude-haiku-4-5-20251001`, `gpt-4.1-nano`, `gemini-2.5-flash-lite`. The actual live calls below, not a private receipt, establish this candidate's observed access.
3. One parent question after a public grader failure: which public compatible ChatCompletions field describes cache-write input? The parent supplied the current OpenAI Chat Completions documentation field `usage.prompt_tokens_details.cache_write_tokens` and clarified that it is optional. I independently fetched the public documentation and retained a sanitized field excerpt/URL in `public-openai-cache-write.txt`. No grader internals were read. The normalized ordinary-input count subtracts both reported cache categories and clamps at zero.

All model/key hints are recorded as assistance, not permanent API guarantees. Keys were selected only into the launcher process/child environment. No credentials, full settings store, authentication headers, or raw provider request/response archives were printed, written, or committed. The saved logs contain the student's normalized chapter events and visible conversation content.

## Architecture and Chapter 1 delta

- Preserved the library root, public `Ensemble` entry point, root-owned agent construction, explicit parent identity, private implementation fields, existing command package, and interface-only implementation cross-dependencies. Import listing confirmed no implementation spoke imports another spoke.
- Root assembly attaches one History and one Engine to each Agent. Agent owns orchestration/capture; History owns the append-only event slice and disposable projection; Engine owns rendering, parsing, transport, and live usage totals. No alternative production test constructor or test hooks were added.
- Replaced string-pair storage as the source of truth with typed, versioned JSONL facts. Shared vocabulary includes closed Vendor/Surface/RefKind values, explicit model/surface provenance, text/tool-call/tool-result/blob/opaque/redacted parts, events, and turn state. The legacy `History()` text view remains available for Chapter 1 committed-dialogue expectations.
- History clones nested DTO data at public boundaries, validates before append/load, assigns sequence identities, accepts an optional supported version header, rejects legacy/zero-kind blob references, and loads transactionally. Stable Config system text is outside the log. Timestamps are metadata, never ordering authority.
- Context rebuild handles identity transitions, message authorship, pending tools, all four redaction levels, Ref-retaining stubs, and one-request system ephemera. Render consumes snapshots without mutating History or cumulative usage. Redaction changes only the projection; original facts remain dumpable.
- Typed request structs and full-response parsers implement Anthropic Messages, OpenAI **ChatCompletions**, and Gemini generateContent. Returned provider model identity takes precedence; optional `LLM_RESOLVED_MODEL` supplies an exact replay identity when the routing name differs. Model-bound material is replayed only on exact vendor/model/surface equality. Gemini signatures stay attached to their call and absent native IDs remain absent in normalized facts; missing outbound IDs use deterministic sequence/part positions where required.
- Anthropic grouping places tool results before text across a whole merged user message. OpenAI uses assistant calls plus `role:tool` results. Gemini uses model calls plus user functionResponse parts. Unsupported blobs are refused loudly; no loader, inline-media fallback, tools execution/registration, streaming, configuration framework, CLI ID API, or native Responses requirement was introduced.
- CLI contract is the actual chapter contract: `chat`, `render LOG`, `dump`, or no args for JSONL stdio. `LLM_VENDOR`, `LLM_BASE_URL`, `LLM_API_KEY`, `LLM_MODEL`, and `CH02_LOG` are honored; historical Anthropic variables remain compatible. Ephemeral directives get an acknowledgement, no model call, and no usage. Sessions persist on success and failure. Offline commands need no key.
- No overall HTTP/client generation deadline or response-header timeout was added. Dial/TLS establishment remains bounded, caller cancellation propagates, and delayed body/header checks passed beyond a minute each.
- Added a separate `gui/` module placeholder that consumes the public root and builds independently. It is not a GUI deployment or live app integration and has no GUI dependencies or automated GUI behavior claims.

## Final validation — exact candidate only

All four main final commands started after committing and rebuilding the exact candidate; source did not change while they ran. The two graders used the same binary/hash above.

| Check | Exact command / cwd | Result | Evidence |
|---|---|---|---|
| Chapter 2 grader | authoritative root: `go run ./cmd/grade-edition2 -ch 2 -json /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-evidence/student-cli` | **exit 0, passed:true, 100/100** | `final-grade-e2-ch02.json`, `.stderr`, `.exit` |
| Edition 3 Chapter 1 preservation | authoritative root: `go run ./cmd/grade-edition3 -ch 1 -json /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-evidence/student-cli` | **exit 0, passed:true, 100/100**, including 62-second header and body delays | `final-grade-e3-ch01.json`, `.stderr`, `.exit` |
| Full source race tests | student root: `go test -race -count=1 ./...` | **exit 0**, root 125.685s, CLI 1.039s | `final-race.txt`, `.exit` |
| Source coverage | student root: `go test -count=1 -coverpkg=./... -coverprofile=/Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-evidence/final-cover.out ./...` | **exit 0**, root 124.698s; merged **88.8%** | `final-coverage-tests.txt`, `final-coverage.exit`, `final-cover.out`, `final-coverage-functions.txt`, `final-coverage-summary.json` |
| Static/style | student root: `go vet ./...`; `git diff --check` | both **exit 0**, clean worktree | `final-source-checks.txt` |
| Isolated GUI placeholder | student `gui/`: `go test ./...`; `go build -o /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-evidence/gui-stub .` | both **exit 0**; no GUI-module tests | `final-gui-checks.txt` |

Coverage denominator: **745 / 839 statements**, 630 unique instrumented blocks, all root-module packages selected by `./...`, including the shipped CLI code. The separate GUI module, test-source lines, graders, and external/live CLI executions are **not** counted. Direct CLI-runner tests use real stdin/stdout files and construct the real public root. Package-local 0% lines for implementation packages with no colocated tests are not the merged result; public-root tests exercise those implementations via `-coverpkg=./...`. The assigned rules require explicit coverage/mutation evidence but specify no numeric percentage threshold. Uncovered branches remain; no 100% claim is made.

### Targeted mutation checks

Command (student root): `python3 scripts/check-mutations.py /Users/bill/projects/ensemble/docs/edition-3-notes/ch02-student-evidence`.

Each final mutation ran a selected real-root test with `go test . -count=1 -run ...`, produced a behavioral failure/nonzero exit, restored the exact source bytes, then produced a passing exit. Source hashes before/after match in each record. `final-mutations.txt` reports all five PASS; individual output is in `mutation-*.txt`.

1. Fail to consume one-request ephemera → state/projection test fails.
2. Omit OpenAI cache-write subtraction/accounting → normalized usage test fails.
3. Replay Anthropic opaque material across a model mismatch → provenance isolation test fails.
4. Put text before results in an Anthropic merged user message → merged-order test fails.
5. Drop original Ref from a redacted result stub → reduced Ref-retention test fails.

This complements (not replaces) the real-CLI grader checks, which cannot expose exact reduced Ref objects.

## Actual live CLI demonstrations and saved replay

Exactly **six live generation calls**: one actual `student-cli chat` process per provider, two user turns per process, no automatic retries, max 1024 output tokens per call, EOF exits. No tiny setup probes or alternate application constructors are presented as evidence. `live-prompts.txt` contains the input: explain event-log versus context-view with a concrete chat example, then explain why rendering the saved log offline needs no further model inference. Each second answer refers to the first explanation.

| Provider/surface | Routing model → returned model | CLI exit | Captured input / cache-write / cache-read / output |
|---|---|---:|---:|
| Anthropic Messages | `claude-haiku-4-5-20251001` → same | 0 | **279 / 0 / 0 / 197** |
| OpenAI ChatCompletions | `gpt-4.1-nano` → `gpt-4.1-nano-2025-04-14` | 0 | **244 / 0 / 0 / 147** |
| Gemini generateContent | `gemini-2.5-flash-lite` → same | 0 | **192 / 0 / 0 / 120** |

For each provider, `live-<vendor>.stdout.txt` and `.stderr.txt` are the actual human-mode transcript/usage diagnostics; `live-<vendor>.jsonl` is the actual saved normalized log, containing two messages, two request facts, and two response facts. Capture verification checks both per-response diagnostics and cumulative totals against those response facts. Log hashes and returned identities are in `live-summary.json` and `capture-verification.json`.

**Important preserved evidence failure:** the first evidence launcher itself exited **1** after all three CLI sessions and all offline commands succeeded. Its final verification regex assumed the wrong capitalization/token order for the human usage line. Therefore `live-summary.json` retains `transcript_usage_matches_log:false` and `live-launch.txt` retains the assertion failure; these were not rewritten to look like a clean first run. The separate `verify-captures.py` reads existing sanitized artifacts only, uses the actual documented output order, and exits **0** (`capture-verification.txt`/`.exit`). It made **zero** additional live calls and did not access credentials. The candidate source/binary never changed for this evidence correction.

Fresh-process, keyless replay:

- Three `dump` commands reproduce each saved log byte-for-byte.
- Each live log is rendered twice into **each of the three target vendors**: nine repeat pairs / eighteen independent render processes, all exit 0, valid JSON and byte-identical within each pair.
- All offline render targets pointed to a local HTTP sentinel; its POST count was **0**. No API keys were supplied to those processes.
- Exact commands/cwd, target names, exits, and output hashes are in `offline-summary.json`; bodies and empty stderr files are `offline-<source>-to-<target>-{1,2}.*`.
- Additional separately labeled **fake** before/after redaction facts were rendered by the same binary for all three vendors. `PRIVATE-FAKE-RESULT` appears before redaction, disappears afterward, and remains in the original append-only after-log. Evidence: `fake-redaction-*`, indexed in `capture-verification.json`. This is not a live tool demonstration.

## Failures preserved and development corrections

Development results are **not final-candidate validation**. They used earlier uncommitted working versions; only `final-*` and the identified live/replay artifacts above claim the exact committed candidate.

- `initial-grade-e2-ch02.*`: exit 1, passed:false, 10/100. I initially invented flag-based commands/event spellings rather than correctly implementing the chapter's exact subcommands/schema. I reread the public contract and removed that wrong surface. The accompanying initial Chapter 1 run passed but is not relabeled as final.
- `development-grade-e2-ch02.*`: exit 1, passed:false, 60/100. Remaining public failures concerned cache-write normalization, Gemini per-call signature replay, and returned provenance. Corrected those in implementation, using the public cache-write clarification above. `development-grade-e2-ch02-fixed.*` then passed; final runs were still performed anew.
- Early compile/test mistakes included a wrong module import (`chapter1` instead of accepted `ensemble`), old constructor/wiring signatures, inherited fake request decoding that assumed string rather than block content, an old constructor credential expectation incompatible with offline construction, and test helper/signature/name assumptions. Inherited behavior assertions were retained or moved to the real send boundary as appropriate; no production mock hook was added. `development-source-tests.txt` retains the early focused failure. Subsequent temporary compile errors in added tests (nonexistent helper names) were fixed before commit.
- The first mutation attempts refused ambiguous/stale anchors without editing source. An opaque selector initially matched no tests; its output is preserved in `mutation-opaque-isolation-initial-no-tests.txt`. Later the runner mistakenly rejected real behavioral failures because unrelated CLI packages reported no matching tests; `mutation-*-precommit-harness-false-negative.txt` retains those records. The corrected root-package selectors produced all five final fail/restore passes.
- The initial GUI build asked for module metadata synchronization; local `go mod tidy` fixed the replacement-module Go version without adding third-party dependencies. Final GUI checks passed.
- Git's PTY pager interfered with an early diff-stat display; later Git inspection explicitly used `GIT_PAGER=cat`. No repository reset/clean was used.
- A final bookkeeping vet/hash command accidentally used GUI cwd; retained as `development-source-checks-gui-cwd.txt`. Repeated from explicit student-root cwd in `final-source-checks.txt`; root source hash matches the build identity record.
- The live evidence launcher regex failure is described above and remains an actual nonzero exit, not concealed behind later verification.

## Limits and feedback for the existing teaching

- Live evidence covers ordinary text conversation and measured zero cache categories. Nonzero cache-write/read accounting, opaque reasoning/signatures, synthetic tool calls/results, redaction categories, and Ref edge cases were verified with source tests/fake-protocol grader boundaries, **not claimed as live tools or nonzero live cache observations**.
- No tools execute; no media bytes are opened/downloaded; no streaming, compaction/eviction, future chapters, GUI runtime tests, live app connection, deployment, restart, or billing-cost estimation was performed.
- The GUI module is a buildable placeholder only. Root source coverage excludes that separate module and does not imply every branch is covered.
- Live availability is an observation for these sessions, not a guarantee those model aliases remain callable. Returned OpenAI identity visibly differs from the routing alias, demonstrating why exact provenance and routing must remain distinct.
- Chapter 1's Engine accounting includes measured unusable responses; the saved normalized response facts are successful parses, while failed captures record an error event rather than archiving raw provider payloads. The live usage comparisons here involve successful replies only.
- Student feedback: the distinction between the explanatory pseudo-code and the exact CLI/JSONL contract is easy to miss and costly when skimmed. My first failure was an implementation-reading error, not a reason to change the contract or grader. A student checklist should emphasize `render LOG`/`dump`, the optional header, `agent` actor spelling, failure capture, and generic `LLM_*` configuration.
- The library-root build clarification was essential and should accompany the chapter's `go build ... .` example for accepted predecessors. Current OpenAI optional cache-write documentation and Gemini signature-value placement were genuine provider-boundary details worth pointing to with current public URLs.
- The mutation exercise was useful: it exposed both a wrong selector and a false-negative harness condition. Ref-retention needs a source-level assertion because vendor wire output intentionally hides the locator after redaction.
- The exact candidate is ready for the parent's independent revision-specific code review. No acceptance is declared and no wait-for-review loop or Chapter 3 work is started.
