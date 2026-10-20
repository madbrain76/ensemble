# Chapter 9 restart — student candidate HOLD

Implemented in `solutions/edition-3/main`; held at the first coherent compiling candidate for root and independent source review. Initial implementation, not a correction round (0 of 10 used). No source edits while review runs.

## Official Anthropic research

Retrieved the current public Markdown documentation over HTTPS; no provider requests or credentials were used.

- https://platform.claude.com/docs/en/agents-and-tools/tool-use/tool-search-tool (retrieved as the same URL plus `.md`), especially **How tool search works**, **Deferred tool loading**, **Custom tool search implementation**, and **Model compatibility**. The page states: “You still send every tool's full definition in the `tools` array on every request, including the deferred ones.” Its custom implementation returns a standard `tool_result` containing `{"type":"tool_reference","tool_name":"discovered_tool_name"}`. Each reference requires a corresponding top-level definition, normally marked `defer_loading: true`. At least one tool remains non-deferred. This is supported content-block disclosure, not an instruction channel. The current compatibility table explicitly excludes Opus 4.1 and earlier; the adapter conservatively recognizes the documented model families and dated identifiers, with ordinary declaration fallback for unknown models. This is documentation-based compatibility, not observed provider behavior.
- https://platform.claude.com/docs/en/api/messages (retrieved with `.md`). The input-message documentation says to use the top-level `system` parameter and that “there is no `"system"` role for input messages in the Messages API.” No supported incremental instruction channel was identified in these interfaces. Loaded instructions therefore join system text; ordinary lifecycle-tool output is never promoted into instructions.
- https://platform.claude.com/docs/en/build-with-claude/context-editing (retrieved with `.md`), **Server-side strategies**, **Tool result clearing**, and **Thinking block clearing**. Documented beta strategies clear tool-result content and thinking; `clear_tool_inputs` can additionally clear tool-use parameters. These are not documented per-skill instruction/declaration revocation operations. No context-editing or compaction feature is enabled here. Pending unload retains instructions, declarations/references, and dispatch capability.

The deferred catalog is sent on every Anthropic request because the documented API requires it. Dormant catalog entries are not yet model-visible declarations and are not dispatch-enabled. Successful tool-driven activation attaches references at the original enabling call's result. Initial/manual tools remain non-deferred. System text changes can invalidate caching; no cache-hit claim is made. OpenAI, Gemini, and unsupported Anthropic models receive the current system instructions and ordinary enabled declarations.

## Ownership and implementation

- `internal/skills` owns parsing/discovery, graph validation, activation, rendered bodies, offered/loadable filtering, variables, pending state, and concrete tool activation/provenance. Common contains contracts; the root injects owners. No sibling imports or shared state bag.
- The existing tool owner supplies built-in handlers and retains jobs/dispatch. Declaration rendering and dispatch consult the same capability snapshot. The concrete record is named `ToolSpec` to avoid collision with the existing internal `Tool` actor constant.
- Full dependency/tool validation precedes exposure. Activation publication is mutex-protected; application variable renderers and state callbacks run outside owner locks. Startup construction also releases the root lock before invoking renderers.
- Initial provenance remains EventSeq 0. The Agent supplies the actual recorded `tool_called` sequence before release. Shared/duplicate grants retain first-enable provenance. Explicit manual API/GUI loading has a separate manual origin and does not invent a tool event.
- Known variables are expanded and retained per activation; unknown variables remain literal. Pending unload does not remove capabilities or rewrite instruction bodies.
- Optional `SkillSpec` preserves the non-skill registration path. CLI/GUI use the existing `ConfigFromEnv` entry path with `EN_SKILLS_DIR`, `EN_PRIMARY_SKILL`, and captured `EN_CUSTOM_VAR`; config maps are copied per Agent. GUI exposes discovered/enabled/loadable/pending state and manual controls, including reconnect state.
- New focused tests cover metadata/lists, graph failures, direct-load boundaries, actual writer dispatch, provenance, Agent isolation, callbacks, variables, pending unload, provider rendering, environment inputs, and GUI state. They have only been compiled.

## Actual checks and limits

Passed, sequential module checks: `go build ./...` and `go test -c -o <owned module output directory>/ ./...` in root/core, `gui`, `ch06`, `ch07`, and `examples/ch04-consumer`. Modules/packages without tests reported that fact. Also passed explicit root `go build -o <owned path> ./cmd/ch01`, GUI `go build -o <owned path> ./cmd`, and `node --check gui/web/gui/app.js`.

Builds used a clean environment with owned HOME/TMP/cache/output directories, `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, and `GOMODCACHE=/Users/bill/go/pkg/mod`; no inherited provider keys or `LLM_BASE_URL`. Final observed free disk: about 6.7 GiB.

Initial compile failures were corrected before HOLD: the Tool name collision/rename fallout and new-test API/helper name mismatches. A power/connection-reset interruption occurred before compilation; Bill authorized resumption in the same context. Saved work was reconciled, not treated as a product failure. Earlier documentation temporary paths were unavailable after recovery, so the two official pages were retrieved again and the current model table reconciled. No behavioral result is claimed.

No behavioral/full/race/grader/provider/browser/native/MIC campaigns, paid calls, key reads, children, commits/tags, installation/restart, shared cleanup, or manuscript/grader/guidance edits. Inherited tests remain unchanged; obsolete frozen-prompt grader alignment remains root-owned. Assistance: task constraints and resume authorization from Bill/root; no delegated implementation or review. Source correctness, runtime behavior, GUI interaction, provider acceptance, and model-list evolution remain for review/authorized later verification.

## Formal correction Round 1 of 10 — complete, HOLD

Only the three released findings were addressed; no architecture/API restart or withdrawn loadable-formula change.

- Restored registration checks in the existing skills registry owner: nonblank description and explicit `type: object` in supplied parameter schemas. Omitted parameters retain the existing object-schema default. Both public registration paths reject invalid input before publication; focused tests check unchanged declarations/provider catalog, unavailable dispatch, and successful later valid registration.
- Anthropic rendering now defers an enabled tool only when that request actually emits its reference at a successful, real paired `load_skill` result. Direct Dispatch/Invoke activation and contexts lacking the enabling pair use ordinary declarations. Dormant tools stay deferred; startup/manual paths and real paired model disclosure remain intact. This is request-local projection, not altered History or provenance, and does not repair generic orphan-result behavior.
- Corrected malformed-metadata fixtures to test required boundaries/name/type, and checked minimal metadata plus duplicate discovered identity. Parser production behavior was not tightened. Replaced the manually appended provenance fixture with real Agent admission: local-HTTP model-authored loads are distinct from public Dispatch/Invoke. Tests verify the actual recorded first sequence, replay/absent-pair fallback, shared initial provenance, pending/reload stability, manual fallback, and failed-load non-publication.

Actual receipts (existing local files; no new evidence bundle):

- `/tmp/ensemble-ch09-restart-student/round1-before.log`: initial test-compilation setup failure (unused imports after replacing the old fixture), corrected before behavioral testing.
- `/tmp/ensemble-ch09-restart-student/round1-before-behavior.log`: before-fix registration acceptance/dispatch and direct-load deferral failures. Its model-path result was limited by a test decoder that did not yet accept inherited text-only string messages; that decoder was corrected, not product code.
- `/tmp/ensemble-ch09-restart-student/round1-before-model.log`: corrected decoder reproduces missing-pair model-context fallback failure before the fix.
- `/tmp/ensemble-ch09-restart-student/round1-before-parser.log`: original parser test fails on its unsupported missing-description requirement.
- `/tmp/ensemble-ch09-restart-student/round1-after.log`: passing focused root regressions, explicit parser tests, and `TestSkillToolReferencesAreContentNotInstructions` after production fixes.
- `/tmp/ensemble-ch09-restart-student/round1-final-focused.log`: passing final root `^TestChapter09(Round1|EnablingEventAndRendering)` selection, including manual fallback and duplicate discovered identity. Unchanged passing parser/renderer checks were reused.
- `/tmp/ensemble-ch09-restart-student/round1-compile.log`: passed sequential `go build ./...` and `go test -c` in root, GUI, ch06, ch07, and examples/ch04-consumer, plus explicit CLI `./cmd/ch01` and GUI `./cmd` builds. Compiled binaries are under `/tmp/ensemble-ch09-restart-student/build/round1/`.

Same clean offline Go environment and owned HOME/TMP/cache as initial compilation; local HTTP only with fixture credentials. Last observed free disk before compilation was about 6.3 GiB. No full/race/grader/browser/native/real-provider campaign, private grader inspection, children, source commits, or manuscript changes. No new documentation research. Root/independent findings were the only correction assistance. All commands finished; source remains HELD for review. One formal correction round used.

## Authorized full validation — stopped at first failure

Ran root `go test -race ./...` in the same clean offline environment. Exit 1: `TestCh6R1CallRecordBarrier/invoke/append/registered` panicked with a nil pointer at `internal/tool/dispatch.go:20` (`registry.RegisterTool`), called via `Register` from `chapter06_round1_test.go:101`. The unaltered full output is `/tmp/e3-ch09-final/root-test.log`. This is observed test failure evidence; root classification is pending.

The sequential gate command stopped immediately. Root vet, GUI gates, separate-module gates, final frozen CLI build, and graders were not run. No production/fixture changes or failure investigation patches were made. Initial observed free disk was about 6.2 GiB. Source remains HELD; root has been asked to classify before any further action.

## Fixture-only adaptation and validation retry

Root and the same reviewer classified the first panic as fixture-only and authorized only `chapter06_round1_test.go` composition repair. Added the actual skills owner from `registry.Catalog()`, attached it, and checked `LoadInitial()` before `Start`, matching the production root. No production fallback or test assertion changes. Focused `go test . -run '^TestCh6R1CallRecordBarrier$' -count=1 -v` passed once, all 36 cases: `/tmp/e3-ch09-final/barrier-fixture-v2.log`. This is not another production correction round.

Resumed root `go test -race ./...` with a versioned log. Exit 1 at a new panic: `TestCh7R1SavedAdmissionFrontier/failure=false/kill-before-release=false`, `internal/tool/dispatch.go:20` via `Register`, called from `chapter07_round1_test.go:109`. Exact output: `/tmp/e3-ch09-final/root-test-v2.log`. Original `/tmp/e3-ch09-final/root-test.log` preserved. Stopped immediately for root classification; no Chapter 7 fixture investigation/edit, no root vet or later gates, no frozen CLI or grades. Last observed free disk about 6.0 GiB. Production and fixtures remain held after the one explicitly authorized composition adaptation.

## Validation completed after authorized fixture composition fixes — HOLD

Root classified the Chapter 7 panic as the same fixture-only omission and authorized one bounded tracked-test composition scan. Added only actual Skills-owner construction/attachment/checked `LoadInitial` before `Start` in `chapter07_round1_test.go`, `chapter07_round2_test.go`, and `chapter07_round3_test.go`. No assertions or production code changed; the Agent seam fixture was not the same omission and was left untouched. Focused saved-frontier, before-execution-boundary, and actual-admission-save tests passed: `/tmp/e3-ch09-final/ch7-fixtures-v3.log`. Reused the Chapter 6 barrier's 36-case PASS. These adaptations do not consume another production correction round.

Completed sequential clean offline gates:

- Root: `go test -race ./...`, then `go vet ./...` — PASS (`/tmp/e3-ch09-final/root-test-v3.log`, `root-vet-v3.log`).
- GUI module: the same race/vet commands — PASS (`/tmp/e3-ch09-final/gui-test-v3.log`, `gui-vet-v3.log`).
- Separate ch06, ch07, examples/ch04-consumer modules: `go test ./...`, then `go vet ./...` — PASS (`/tmp/e3-ch09-final/ch06-test-v3.log`, `ch06-vet-v3.log`, `ch07-test-v3.log`, `ch07-vet-v3.log`, `consumer-test-v3.log`, `consumer-vet-v3.log`). Ch07 and consumer reported no test files: compile/vet coverage only. Earlier failure logs are preserved.

Built once from main: `go build -o /tmp/e3-ch09-final/frozen-cli ./cmd/ch01`. Build/identity receipt: `/tmp/e3-ch09-final/frozen-cli-build.log`. Size 9,556,706 bytes; SHA-256 `f7699809a0303e26047a77bd298981af6aad453874c22078e96b35f02eb17eb7`. The same hash was checked before Chapter 1; no rebuild between grades.

From `/Users/bill/projects/ensemble`, invoked only the grader command, without source inspection:

```
go run ./cmd/grade-edition3 -ch 9 -json -core /Users/bill/projects/ensemble/solutions/edition-3/main -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui -gui-model claude-sonnet-4-6 -media-models /tmp/e3-ch09-final/media-models.json /tmp/e3-ch09-final/frozen-cli
go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch09-final/frozen-cli
```

Owned media map declares strings: anthropic=`claude-sonnet-4-6`, openai=`gpt-5.4`, gemini=`gemini-3.8-flash`. Chapter 9 returned exit 0 and `passed: true` (100/100); only then Chapter 1 ran, also exit 0 and `passed: true` (100/100). Raw outputs: `/tmp/e3-ch09-final/ch09-grade.json`, `ch09-grade.stderr`, `ch01-grade.json`, `ch01-grade.stderr`. Chapter 1 explicitly reports complete answer and usage after both real `1m2s` header and body delays. No outer timeout was imposed and no observation interval cancelled execution. Latest observed free disk before Chapter 1 was about 5.9 GiB.

All validation used the prescribed clean environment, owned HOME/TMP/cache, offline module settings, and no real provider keys or generic LLM_BASE_URL. No grader source inspection/edits, production patches, children, commits/tags, or manuscript changes. All commands finished.

Proposed remaining manual scope (not executed; requires root operational authorization):
1. One temporary Chapter 9 CLI/GUI walkthrough covering offered skill load, dependency activation, variable rendering, an actual file write, and pending unload retaining effective capabilities; exercise GUI Load/Unload controls and displayed state. Reuse the fixture rather than replay Chapter 8 speech/layout work.
2. One bounded real Anthropic supported-model scenario verifying the model-authored load -> paired tool_reference -> newly enabled tool path and ensuing real file operation, with loaded instructions in system context. No vendor campaign; root must authorize provider/key access and any needed GUI launch/browser use first.

Production and fixtures HELD. Await root operational authorization before any manual execution, harness rehearsal, provider/key/discovery, browser, or native work. One production correction round used.

## Operational attempt stopped at startup — awaiting root classification

Prepared only the newly authorized operation under `/tmp/e3-ch09-operations-ZHspdm`, reading the two explicitly permitted nonsecret helper files and current source. Built the one permitted GUI binary (`go build -o /tmp/e3-ch09-operations-ZHspdm/control/gui ./cmd` from main/gui, prescribed clean offline environment); SHA-256 `ae7bebd31a0571122ba149aa61c6c27a7d221ba45911301020e18f385dd17a15`. Frozen CLI hash still matched `f7699809a0303e26047a77bd298981af6aad453874c22078e96b35f02eb17eb7`; operation-start source aggregate hash was `eb4c77b3074857fdcce0c1dbd5f1a55d4a4e97e4c76e3ccbea3b5132992adfdb`. Source and fixtures unchanged. Scripts passed syntax checks only; no rehearsal.

Single launch command: clean minimal-env `python3 /tmp/e3-ch09-operations-ZHspdm/control/run.py`. It stopped immediately: GUI exit 2, stderr `flag provided but not defined: -addr`. Usage reports `-port` (default 8088), `-gui-dir`, and `-settings`, not `-addr`. No automatic retry or patch. Controller cleanup also rejected its canonical `/private/tmp/...` profile path because its guard expected a `/tmp/...` prefix. GUI PID 66775 was confirmed exited, owned proxy port 58244 had no listener, and the still-empty owned profile was removed with `rmdir`; Chrome had never launched. Approximately 5.8 GiB free.

Exact saved evidence and code: `/tmp/e3-ch09-operations-ZHspdm/control/run.py`, `operate.mjs`, `controller.stdout`, `controller.stderr`, `gui.stderr`, `operation.jsonl`. Local walkthrough did NOT pass. Counts: fake requests 0/4, paid POSTs 0/6, discovery GETs 0/1, key reads 0/1, browser pages 0, native actions 0, screenshots 0. No live phase or provider access occurred. Root classification requested before any controller correction/retry; production/fixtures remain HELD.

## Authorized operational continuation — local PASS; live stopped on HTTP 400

Root classified the first attempt as setup-only and authorized one versioned controller correction/new owned attempt. Preserved `/tmp/e3-ch09-operations-ZHspdm` unchanged; reused its single built GUI and the same frozen CLI (both hashes rechecked). Fresh attempt: `/tmp/e3-ch09-operations-v2-HNhkeT`. Controller changes only: actual `-port` chosen free/non8088, absolute `-gui-dir /Users/bill/projects/ensemble/solutions/edition-3/main/gui/web/gui`, owned settings, PID/listener/HTTP readiness instead of guessed stdout, and exact resolved scratch/profile cleanup containment. No source/fixture edits or rebuilds. Invocation: clean minimal-env `python3 /tmp/e3-ch09-operations-v2-HNhkeT/control/run.py`.

Local GUI walkthrough PASSED in about 5.37 seconds, 4/4 fake model requests. Verified owned GUI PID 66852/listener 58254 and muted private-profile Chrome PID 66861/CDP 58256 before attachment; proxy port 58253. One page. Used actual public Load and Mark pending unload controls: code-tools available -> enabled with helper dependency -> pending unload/still active. Actual jobs: write note.txt, then another write AFTER pending state, then read updated file. Actual result/final correlation and rendered DOM were checked; system contained code-tools/helper instructions and `SKILLS_OK`. Final local file `/tmp/e3-ch09-operations-v2-HNhkeT/local-work/note.txt` contains `SKILLS_OK\nAfter pending unload.\n`. This is controlled CDP/DOM plus real Agent/job/file evidence, not human/direct visual assessment. No screenshot or native action. Local raw evidence: `control/page.jsonl`, `frames.json`, `local-1..4.request.json`, `local-1..4.response.sse`, `local-history.jsonl`, `browser.stdout` under that attempt.

Only after local PASS, cleanup and unchanged CLI/source identity, performed the one authorized exact key read and verified-TLS discovery GET. Fresh `/v1/models` listed `claude-sonnet-4-6` (`control/models.json`). Live POST #1 succeeded, returning model-authored `load_skill` for code-tools, call ID `toolu_014HPG4Y5HAJYEpXJYt4nE4r`. POST #2 was rejected upstream with HTTP 400:

> `invalid_request_error`: `Tool definitions/code execution functions cannot be mixed with other content`

Exact raw error: `/tmp/e3-ch09-operations-v2-HNhkeT/control/live-2.response.sse` (contains the JSON error body despite the stream-oriented filename). Provider request ID: `req_011CgDjUuZ8wyoMknK7YYWm1`. `live-2.request.json` contains the real load call and matching successful tool_result, whose content has the load status text block plus `tool_reference` for write_file; catalog write_file remains deferred, read_file/startup controls are normal. System contains rendered code-tools/helper instructions and `SKILLS_OK`. Requests retained `model=claude-sonnet-4-6`, `max_tokens=1024`, `stream=true`; raw request bytes and successful response stream were forwarded unchanged. The controller's fail-stop response surfaced as HTTP 502 to the CLI; the actual upstream HTTP 400 body is preserved separately. No automatic retry/fallback or source patch. Live note.txt was not created; live write/read/unload/final behavior is NOT established.

Total operational counts across both attempts: fake 4/4, paid POSTs 2/6, discovery GET 1/1, exact key read 1/1, native actions 0, screenshots 0. POST #1 usage: input 1004, output 93, cache-create/read 0; rejected POST #2 returned no usage. First live POST reservation to rejection was about 2.39 seconds; live phase through cleanup about 3.10 seconds. No credentials in argv, files, logs, or child environments; owned children used dummy key/minimal environment. Source remained held.

Cleanup verified GUI/Chrome/CLI exits, no remaining profile-owned Chrome helpers, removed only owned profile, and closed ports 58253/58254/58256. Approximately 5.8 GiB free. Full current evidence/scripts: `/tmp/e3-ch09-operations-v2-HNhkeT/control/{run.py,operate.mjs,controller.stdout,controller.stderr,operation.jsonl,gui.stderr,cli.stdout,cli.stderr,live-1.request.json,live-1.response.sse,live-2.request.json,live-2.response.sse,live-history.jsonl}`. Raw requests remain private local. Stopped and requested root classification; no further key read, discovery, provider use, browser run, research, or patch authorized by this result. HOLD.

## Formal correction Round 2 of 10 — refs-only wire projection; HOLD

Implemented only the classified production correction in `main/internal/engine/render.go`: a successful, actually paired, reference-bearing Anthropic `load_skill` result now contains only documented `tool_reference` blocks. Removed the status-text initializer from that wire-only branch and retained a short explanation. No handler, History, activation/provenance, catalog, system-instruction, cache, or other-provider changes.

Strengthened existing `chapter09_round1_test.go` actual-Agent/local-HTTP checks for next-request/replay references-only content, real pairing/success, canonical handler status, unchanged serialized History/context, provenance, pending/reload, direct Invoke/Dispatch/manual fallback, and failed loads. Extended `internal/engine/ch09_skills_test.go` for preserved canonical status and non-Anthropic output plus ordinary/no-reference/failed/missing Anthropic results. No new harness/framework.

Before-fix behavior reproduced: root actual model-load/replay assertions failed on mixed `text` content (`/tmp/e3-ch09-round2/root-before-behavior.log`); engine assertion failed on canonical status appearing in reference-bearing wire content (`engine-before.log`). Unchanged ordinary/failed/direct/manual cases passed before the fix. One distinct setup failure preceded execution: prior owned TMP path was absent (`root-before.log`); created fresh Round 2 owned HOME/TMP/cache and reran without source workaround.

After-fix focused commands passed sequentially, using the prescribed clean offline environment (`GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, local toolchain, proxy/sumdb off, existing module cache, owned HOME/TMP/cache, no provider keys/base URL):
- `go test . -run '^(TestChapter09EnablingEventAndRendering|TestChapter09Round1DirectLoadExposure)$' -count=1 -v` — `/tmp/e3-ch09-round2/root-after.log`.
- `go test ./internal/engine -run '^(TestSkillToolReferencesAreContentNotInstructions|TestChapter09Round2OrdinaryResults)$' -count=1 -v` — `/tmp/e3-ch09-round2/engine-after.log`.
- Root `go build ./...`; then sequential `go test -c -o /tmp/e3-ch09-round2/build/<module>/ ./...` for root, GUI, ch06, ch07, consumer — all compile PASS, `/tmp/e3-ch09-round2/compile.log`. No-test outputs are compile-only. No new retained CLI/GUI executable; existing frozen operational binaries were not rebuilt and do not include this correction. About 5.7 GiB free.

Production/tests now HELD for root + retained review. No full suites/race/vet/grades/browser/provider/key/discovery or operational retry in Round 2. Earlier four-fake local GUI PASS carries only for unchanged behavior; the two-paid live attempt remains FAILED and this correction is not yet provider-validated. Root owns any grader adjustment and further validation/live authorization.

## Round 2 held-source validation — PASS; HOLD for separate live authorization

After root + retained reviewer accepted the narrow correction for validation, ran fresh root `go test -race ./...` and `go vet ./...` sequentially in the prescribed clean offline environment. Both PASS. Receipts: `/tmp/e3-ch09-round2/root-race-validation-v1.log` and `root-vet-validation-v1.log`. Source/tests unchanged during validation. Carried, without rerunning, the accepted unchanged GUI/module gates and four-fake local GUI walkthrough.

Built exactly one new retained CLI: root `go build -o /tmp/e3-ch09-round2/frozen-cli ./cmd/ch01`. SHA-256 `1471a4c5eb7a61e66aef47748c91fa036b1b854eebb4544a69129e21e1e343a4`; `/tmp/e3-ch09-round2/frozen-cli-build-v1.log`. Old `/tmp/e3-ch09-final/frozen-cli` preserved, hash still `f7699809a0303e26047a77bd298981af6aad453874c22078e96b35f02eb17eb7`. No GUI build.

From book root, black-box invocation only:
`go run ./cmd/grade-edition3 -ch 9 -json -core /Users/bill/projects/ensemble/solutions/edition-3/main -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui -gui-model claude-sonnet-4-6 -media-models /tmp/e3-ch09-final/media-models.json /tmp/e3-ch09-round2/frozen-cli`.
PASS: process exit 0 and parsed JSON `passed:true` (100/100); no numeric chapter/title assertion. Same model map: anthropic `claude-sonnet-4-6`, openai `gpt-5.4`, gemini `gemini-3.8-flash`. Receipts: `/tmp/e3-ch09-round2/ch09-grade-v1.json` and `ch09-grade-v1.stderr`.

Only after that PASS, `go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch09-round2/frozen-cli`, same hash verified. PASS: process exit 0 and parsed JSON `passed:true` (100/100). Report explicitly confirms complete answer AND usage after real 1m2s headers delay and 1m2s body delay. No outer deadline/observation cancellation. Receipts: `/tmp/e3-ch09-round2/ch01-grade-v1.json` and `ch01-grade-v1.stderr`.

No grader source inspection/edit, production patch, browser/native action, new provider request, key read, discovery, operational retry, or new harness/proposal stage. Earlier live two-POST failure stays failed; real-provider acceptance of the new projection remains unverified. Source, fixtures, and operations HELD pending root inspection and explicit live-continuation authorization.

## Authorized corrected live-only attempt — PASS; HOLD for root review

One bounded run, no rehearsal or additional analysis afterward. Adapted the existing owned v2 controller by removing all GUI/browser/fake/native/screenshot/discovery execution paths. Fresh owned attempt: `/tmp/e3-ch09-live-r2-SpjB1F`, same base/code-tools/helper fixture and `EN_CUSTOM_VAR=SKILLS_OK`. Script syntax check passed, then clean minimal-env `python3 /tmp/e3-ch09-live-r2-SpjB1F/control/run.py` executed once (process exit 0). No source/fixture product edits or builds.

Pinned CLI `/tmp/e3-ch09-round2/frozen-cli`, SHA-256 `1471a4c5eb7a61e66aef47748c91fa036b1b854eebb4544a69129e21e1e343a4`; held-source aggregate `6db5998e1473f1178893cd6927fa580e87d8a3821cf623aa8a4ac3ac8d4fa725`, verified before key access and after the exercise. Reused already-verified discovery `/tmp/e3-ch09-operations-v2-HNhkeT/control/models.json` (SHA-256 `02a9188989a823990a206aff15d9764786d298d249cd7b4f355ab834361fcd9e`), which lists `claude-sonnet-4-6`; no new GET. Exactly one newly authorized memory-only directClaudeAPIKey read; key not placed in argv/child env/files/logs. Children used dummy key/owned HOME/TMP; no generic LLM_BASE_URL.

Five new POSTs, all HTTP 200: actual model `load_skill(code-tools)` -> `write_file(note.txt)` -> `read_file(note.txt)` -> `unload_skill(code-tools)` -> final `end_turn`. No forced batching or fixed final phrase. Successful load call ID `toolu_01FRSwsmtsFZxe7J4mhS6WzG`; the NEXT request's paired result was checked to contain ONLY the `write_file` tool_reference, with the valid deferred catalog and primary + rendered code-tools/helper system instructions. Canonical status remained text and the actual first enabling tool_called sequence/return correlation were checked and recorded as `canonical_enable_verified` in `control/operation.jsonl`; no fabricated events. Every model call had a successful actual paired canonical result. Actual `/tmp/e3-ch09-live-r2-SpjB1F/live-work/note.txt` contains exactly `SKILLS_OK\n` (10 bytes). Pending unload and final completion were verified; final reported write/read success and still-active pending state. Future compaction behavior was not exercised. Source/binary pins remained unchanged.

New usage by POST (input/output): 1004/73, 1364/98, 1483/70, 1570/74, 1672/175; total 7093 input + 490 output; cache-create/read zero throughout. New POST count 5/6, cumulative Chapter 9 POST count 7/8 including the preserved earlier two-POST failed attempt. No retry/fallback or request/success-stream body rewriting. New discovery GETs 0; new key reads 1 (cumulative authorized reads 2); no browser/native/fake/screenshot actions. Setup ~0.020s; entire run through cleanup 10.61s, within 60s setup/180s active/240s absolute ceilings. Unused cap was not consumed.

Cleanup verified CLI PID 69266 exit 0 and owned proxy port 58813 closed; controller PID 69265 completed. No GUI/profile existed in this attempt; no shared cleanup. Exact script/raw evidence: `/tmp/e3-ch09-live-r2-SpjB1F/control/run.py`, `controller.stdout`, `controller.stderr`, `operation.jsonl`, `live-1..5.request.json`, `live-1..5.response.sse`, `live-history.jsonl`, `cli.stdout`, `cli.stderr`. Raw request/response evidence remains private local. Prior failure evidence is unchanged; accepted local four-fake GUI/post-pending-write evidence carries without replay. Final acceptance remains root's decision. All source and operations HELD.

### Final evidence clarification for acceptance review

Corrected live attempt: **5 new HTTP 200 POSTs**, **7 cumulative paid POSTs** including the prior failed two-POST attempt. The earlier **4-fake local GUI PASS carries**, without replay. Across operations: **2 total authorized key reads**, **1 total discovery GET** (none in the corrected run). Every new Messages request retained native **max_tokens=1024**; no retries, fallback, or request-body rewriting.

Actual sequence: model load -> write -> read -> pending unload -> end_turn. The successful load's canonical **tool_called Seq 5 / tool_returned Seq 8** retained the handler's text status; its next Anthropic wire result contained **only tool_reference content**, paired to that real call. Actual note.txt contained exactly `SKILLS_OK\n` (10 bytes). New-run usage was **7093 input / 490 output** tokens; reported cache-create/read counts were zero. This is **not a cache-hit claim**. The earlier successful response in the failed attempt reported 1004/93; its rejected response supplied no usage.

CLI exit 0, owned proxy listener closed, controller exit 0. Prior GUI startup/cleanup-guard setup failures, fixture-composition race failures, and the upstream HTTP 400 live failure remain preserved in the earlier records; none is relabeled as a pass. No direct human visual or hearing claim; no native/speech action. Future compaction behavior remains unexercised. Source and operations remain HELD; final acceptance/prose belongs to root after saved-result join and retained review.
