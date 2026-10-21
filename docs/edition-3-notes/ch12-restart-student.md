# Chapter 12 restart — compiling HOLD

Source: `solutions/edition-3/main`. Initial candidate is held for coordinator source comparison and retained independent review; no behavioral validation or operational release is claimed.

## Ownership and delta

- Existing CLI adapter owns shared launch selection (`internal/cli/skills.go`), used by CLI and optional GUI: `--gui-debug`, `--skills-dir`, explicit environment selection, and launch-directory `skills/ensemble/SKILL.md` defaults. Plain launches without that default remain usable.
- `Config.StartupSkills` is copied with agent configuration. Existing skill admission loads these after the primary with `initial` provenance, before SYSTEM rendering; normal offered/loadable admission remains required. Existing dynamic SYSTEM and pending-unload paths were not replaced.
- Existing MCP owner has only the bundled four-tool `browser-debug` WebSocket expected catalog (`internal/mcp/browser.go`). This resolves startup-before-GUI-attachment without pretending a peer is connected. Real attachment still initializes and validates schemas; explicit pipe/socket overrides still discover their peer. Tool dispatch waits on the existing connection readiness and rechecks active ownership. GUI tunnel selection respects the effective endpoint override.
- Browser `draftChanged` records trusted input cause rather than inferring human typing from any nonempty value; human whitespace counts. Existing click/input control mutation and input/change dispatch, explicit snapshot/queue tools, native speech, Continue, DOM rendering, and external pending-unload teardown remain their existing owners.
- Focused Go regressions cover launch selection and startup instructions/admission/provenance. Existing GUI fixtures now identify simulated human input explicitly; a synthetic-draft/speech regression was added. The browser socket fixture now expects startup admission with disconnected call failure and the real four-tool schema contract.

Assistance: coordinator clarified/approved the narrow known-browser startup admission boundary and required compile-only HOLD. No child agents were used. Required public guidance/chapter and actual main source were read; no original/future/private snapshots or graders were inspected. No unresolved public-requirement ambiguity is claimed; expected browser schemas must stay synchronized with `gui/web/gui/app.js`.

## Compilation

All commands ran sequentially with an empty inherited environment, private HOME/TMP under `/tmp/ch12-restart-compile.sisGbe`, Go/Delve PATH, `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, `GOENV=off`, and the specified existing module/build caches. No real keys or generic `LLM_BASE_URL` were inherited. Every compile checked available disk against 1 GiB first. Final available space: 1,150,416 KiB; no cleanup/deletion was performed.

Final successful commands (all exit 0; outputs under the private compile directory):

- Main: `go build -o "$work/ch01" ./cmd/ch01`.
- Main: `go test -vet=off -c -o "$work/core.test" .`.
- Main: `go test -vet=off -c -o "$work/cli.test" ./internal/cli`.
- Main: `go test -vet=off -c -o "$work/mcp.test" ./internal/mcp`.
- `examples/ch04-consumer`: `go build -o "$work/consumer" .`.
- `ch06`: `go build -o "$work/workflow" .`.
- `ch07`: `go build -o "$work/stream" .`.
- `gui`: `go build -o "$work/gui" ./cmd`.
- `gui`: `go test -vet=off -c -o "$work/gui.test" .`.
- `gui`: `go test -vet=off -c -o "$work/gui-cli.test" ./cmd`.

Two initial core test-compilation attempts exited 1 because my new test used incorrect fixture/API names, then a nonexistent `Allowed` method. These new-test mistakes were corrected to the existing `c9Fixture`, `Ensemble.NewAgent`, and capability snapshot interfaces. Failure logs are retained as `core-test-compile.log` and `core-test-compile-retry.log`; final success is `core-test-compile-final.log`. Full command/exit sequence: `/tmp/ch12-restart-compile.sisGbe/commands.log`. Static `git diff --check -- solutions/edition-3/main` also exited 0.

Limits: no test binaries or JavaScript fixtures executed; no behavioral suites, races, vet, grader, browser, native speech, provider/discovery/network/keys, installation/restart, live 8088, microphone/device, deployment/publication, commit/tag, or snapshot actions. Compilation does not establish runtime correctness. Source is frozen pending root review/release.

## Production correction ROUND 1 — resource HOLD

Only the two requested production gaps changed:

- `internal/mcp/manager.go` retains real catalog metadata during attachment and rejects nonempty `Ephemeral` for the known four ordinary WebSocket browser tools, even with matching schemas. Effective non-WebSocket discovery and ordinary attachment retain their existing paths.
- `gui/hub.go` derives `Wire.BrowserMCP` eligibility through the existing effective-endpoint helper; `gui/web/gui/app.js` uses that boolean for auto-open. `gui/projection.go` carries the boolean only: endpoint details are not added and canonical skill declarations are unchanged.

Focused real-path tests added in `gui/chapter12_round1_test.go`: actual MCP handshakes (all four names, both round/turn modes and ordinary controls), actual serialized `skillFrame` into the full production app via the controlled DOM fixture, real local stdio override discovery/call, normal WebSocket preconnect and duplicate-frame suppression. Existing DOM fixture was exported for reuse, given missing current-app DOM stubs, and given an isolated `chapter12` typing entry. The existing same-socket unload/reload fixture now uses dynamic loading again; startup/preconnect is covered by the new tests.

Commands used the same clean offline Go environment and owned HOME/TMP as above, adding `/opt/homebrew/bin` for local Node. Disk was checked before executions; commands were sequential. Logs remain in `/tmp/ch12-restart-compile.sisGbe`:

- Before fix: `go test -vet=off -count=1 -timeout=40s -o "$work/gui.test" -run '^TestChapter12(BrowserAutomaticIdentity|SkillFrameBrowserEligibility)$' .` exited 1 (`round1-before.log`). The metadata defect reproduced for all eight automatic cases; ordinary cases passed. Frame cases were setup failures, not behavioral evidence: initial preparation used wrong relative paths; the inherited DOM fixture lacked `addEventListener`, and the new stdio catalog lacked descriptions.
- After fixture setup corrections, before production fix: `go test -vet=off -count=1 -timeout=30s -o "$work/gui.test" -run '^TestChapter12SkillFrameBrowserEligibility$' .` exited 1 (`round1-frame-before.log`). Effective stdio emitted one unwanted tunnel open; actual stdio call and normal WebSocket control succeeded.
- After production fix: both new regression groups passed, but the combined run exited 1 because the existing reload fixture expected a new event from idempotent startup loading (`round1-after.log`). An attempted fixture connection wait used a nonexistent `mcp_connected` event and also timed out (`round1-after-final.log`). Neither timeout establishes a production failure. The fixture was corrected to its ordinary dynamic-load/unload/reload path; no additional production change was made.
- Final: `go test -vet=off -count=1 -timeout=40s -o "$work/gui.test" -run '^(TestChapter12(BrowserAutomaticIdentity|SkillFrameBrowserEligibility)|TestMCPBrowserReloadSameSocket)$' .` exited 0 (`round1-gui-final.log`). This compiles the current shared core and GUI and executes only the named local fixtures.
- `node gui/client_round1_test.cjs chapter12` from main exited 0 (`round1-typing.log`): synthetic draft does not acquire typing pause; human whitespace does; clearing human input releases it. No browser/audio was launched.

Latest observed free space: 1,051,936 KiB, only about 3.3 MiB above the strict 1 GiB floor. No deletion/cleanup occurred. Further executable builds would require more temporary space than that margin. Coordinator directed resource HOLD, no further compiler/test processes, and read-only review. **Not run:** five-module post-round compilation, focused startup/client tests, full legacy DOM suite, or any broad/race/vet/grader/native/provider/operational validation. Initial candidate module compilation remains recorded above but is not claimed as a post-round rebuild. Resource pause is not a product failure or a new correction round. Source frozen for root and retained reviewer; no children or other assistance used beyond coordinator direction.

## Post-restart held-source validation

Bill manually removed the old `/tmp/ch12-restart-compile.sisGbe` tree during restart; this was not an implementation failure. Coordinator and retained reviewer recovered the original R1 behavioral negatives and focused passes from retained tool output. Those results are carried forward, not recreated for log freshness. R1 source repairs were accepted for validation; no source edits were made in this phase.

New single owned workspace: `/tmp/e3-ch12-validation.Z0OXma` (logs, private HOME/TMP, one `bin/ch01`). Sequential clean offline environment used the specified shared Go caches, `GOMAXPROCS=2`, readonly/p1 flags, local toolchain, proxy/sumdb off, and no inherited credentials or generic base URL. Disk was checked before each execution; initial available space was 24,750,932 KiB.

All authorized local gates exited 0: main `go test -race ./...`, `go vet ./...`; GUI the same two commands; `ch06`, `ch07`, and `examples/ch04-consumer` each `go test ./...` and `go vet ./...`. Main includes the previously unrun Chapter 12 startup/client regressions. `ch07` and consumer report no test files: their test commands establish compilation only. Logs are `root-{race,vet}.log`, `gui-{race,vet}.log`, `workflow-{test,vet}.log`, `stream-{test,vet}.log`, `consumer-{test,vet}.log`; `commands.log` records commands/exits. No redundant standalone build/test-c pass or broad legacy Node run was added.

One frozen CLI built with `go build -o /tmp/e3-ch12-validation.Z0OXma/bin/ch01 ./cmd/ch01` (exit 0; `cli-build.log`). SHA-256: `016ae5794ffce6505dc005008c49ce1ef436822a648220e4c26f3647575ba63e` (`cli.sha256`). Next setup check, `command -v grade-edition3`, exited 1 on the clean PATH. No grader invocation occurred; stopped and asked coordinator for the approved existing executable path rather than searching private locations or rebuilding the grader. No production/fixture changes, children, keys, external providers, native/browser/device, installation/publication, commits/tags, or cleanup actions.

Coordinator clarified the approved grader entry point is `go run ./cmd/grade-edition3` from `/Users/bill/projects/ensemble`; its build was authorized, without inspecting grader source/fixtures. Supplied public map saved as `/tmp/e3-ch12-validation.Z0OXma/media-models.json`: anthropic `claude-sonnet-4-6`, openai `gpt-5.4`, gemini `gemini-3.8-flash`. No discovery or capability edits occurred.

Executed once, in the same clean offline environment:

```sh
go run ./cmd/grade-edition3 -ch 12 -json \
  -core /Users/bill/projects/ensemble/solutions/edition-3/main \
  -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer \
  -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 \
  -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 \
  -stream-model claude-sonnet-4-6 \
  -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui \
  -gui-model claude-sonnet-4-6 \
  -media-models /tmp/e3-ch12-validation.Z0OXma/media-models.json \
  /tmp/e3-ch12-validation.Z0OXma/bin/ch01
```

**Actual exit 1; report `passed: false`, score 90/100.** All six direct Chapter 12 checks passed. Inherited `ch11-parity` failed with these exact details:

- `ws-tunnel: activate declared browser skill: WebSocket observation deadline`
- `ch10-parity: ch9-parity: ch8-parity: ch7-parity: websocket-streams: GUI public dependency: consumer imports private package ensemble/internal/cli; observed wire checks passed before dependency assessment: true`

Preserved output: `chapter12.json`, `chapter12.stderr.log`, `chapter12.exit` in the validation workspace; full invocation also in `commands.log`. Stopped and reported to coordinator immediately; no inspection of private grader code, diagnosis-by-fixture, source edits, retries, or speculative corrections. **Chapter 1 same-binary gate and operational proposal were not run/provided because their Chapter 12 success prerequisite failed.** Source remains held. Further scope requires coordinator review.

## Chapter 12 correction round 2 — focused/compile HOLD

Root and the retained reviewer classified both inherited failures in the first 90/false grade as production defects and authorized this bounded round. No grader or accepted-snapshot edits; no children. Production delta:

- `gui/cmd/main.go` now uses a small local `configureSkills` in `gui/cmd/skills.go`, using only public `ensemble` configuration types. Defaults remain client-owned; explicit directory/environment/default-primary/debug startup semantics are preserved without exporting CLI internals.
- `internal/mcp/manager.go` discovers an attached peer's actual catalog. The known browser catalog is only a disconnected fallback; ordinary skill admission still rejects unknown cold required names. Real attachment validates the schema **and automatic identity of admitted tool definitions**, not an unconditional four-tool catalog for every peer named `browser-debug`. This preserves R1's cold ordinary-mode checks while admitting explicitly granted generic ordinary/hidden automatic tools. `gui_snapshot` and `tts_queue` remain on-demand; other generic automatic tools are not prohibited. `internal/mcp/browser.go` documents this narrower ownership. No new connection/protocol/lifecycle owner or framework.

Focused regressions added in `gui/cmd/chapter12_test.go`, `gui/chapter12_round2_test.go`, and `internal/mcp/browser_test.go`. The real GUI entry point runs with EOF and temporary skills/save/settings, never a model turn. The generic regression uses an actual Agent + GUI MCP tunnel: unknown cold names refuse; the attached peer's distinct schema activates; ordinary actions execute; round/turn tools remain hidden from guessed model calls and gather only when granted; pending unload disables external capability while retaining instructions. Existing R1 metadata/stdio-override/normal-browser and same-WebSocket reload assertions were retained.

One new workspace: `/tmp/e3-ch12-r2.PUv3WB`; clean offline environment and shared caches as previously specified, owned HOME/TMP. Disk checked before executions, always far above 1 GiB (last pre-build check 24,484,728 KiB). Prior `/tmp/e3-ch12-validation.Z0OXma` frozen CLI and failed grade report were not rebuilt or edited.

Results (all tests below used `-vet=off -count=1`, no race/full suite):

- Before production repairs: GUI `go test -run '^TestChapter12(GUIPublicImports|GUIStartupSelection|GenericBrowserTunnel)$' . ./cmd` exited 1 (`before.log`). Behavioral negatives: actual generic tunnel rejected with `MCP browser schema mismatch: gui_snapshot`; GUI import check found `ensemble/internal/cli`. Startup cases passed except one test expecting the wrong missing-directory diagnostic wording; corrected assertion, then `go test -run '^TestChapter12GUIStartupSelection$' ./cmd` passed before production edits (`startup-before.log`).
- After repairs: GUI `go test -run '^(TestChapter12.*|TestMCPBrowserReloadSameSocket)$' . ./cmd` exited 0 (`gui-after.log`): new import/startup/generic-tunnel tests, known cold catalog/automatic-metadata refusal, effective stdio override, normal browser opening, and same-socket reload all passed.
- Main `go test -run '^Test(StartupSkill|ClientSkillSelection|ClientExplicitPrimaryWins)' . ./internal/cli` exited 0 (`startup-client-after.log`). The new `TestBrowserObservationModes` also passed (`core-after.log`). An earlier overly narrow selector in that log matched no main/CLI tests; it was not counted as behavioral evidence, and the corrected command above ran them.
- Five sequential module compiles exited 0: main and GUI each `go build ./...`; `ch06`, `ch07`, and consumer each `go build -o /tmp/e3-ch12-r2.PUv3WB/bin/{workflow,stream,consumer} .` in their respective module directories. Logs: `main-build.log`, `gui-build.log`, `workflow-build.log`, `stream-build.log`, `consumer-build.log`. `commands.log` retains labels/exits and the final commands.

Setup/edit issues, not product evidence: the first test-launch command used an incorrect workspace suffix and never started tests; an exact-anchor MCP edit refused before writing MCP source, then was reapplied against the current block. Tool reports retain those failures. No full suites, race, vet, grader retry, browser/native/device/provider/key/discovery actions, operations, commits/tags, cleanup, or snapshot creation in this round. Source frozen after focused success and five-module compilation; awaiting root/retained-reviewer inspection and next release.

## Held R2 validation — passed; operations proposal only

Root and the retained reviewer accepted R2 for validation, not overall acceptance. No source edits in this phase. One workspace: `/tmp/e3-ch12-r2-validation.20Rg53`, with owned HOME/TMP and the required clean offline environment/shared caches. Sequential main and GUI `go test -race ./...` and `go vet ./...` all exited 0 (`root-race.log`, `root-vet.log`, `gui-race.log`, `gui-vet.log`). Unchanged ch06 actual test/vet and ch07/consumer compile-only test/vet results are carried from `/tmp/e3-ch12-validation.Z0OXma`, together with R2's five module compiles; they were not repeated for receipts. Disk checks remained above 23 GiB.

Built exactly one new frozen CLI with `go build -o /tmp/e3-ch12-r2-validation.20Rg53/bin/ch01 ./cmd/ch01` (exit 0). SHA-256 `d72313023079544f5a31e11de978150506db1c32d0a1daba48890c1a1c92f52f` (`cli.sha256`, `cli-build.log`). The prior frozen CLI and 90/false report remain untouched.

From the book repository root, with `B=/Users/bill/projects/ensemble/solutions/edition-3/main` and `W=/tmp/e3-ch12-r2-validation.20Rg53`, executed once each:

```sh
go run ./cmd/grade-edition3 -ch 12 -json -core "$B" -consumer "$B/examples/ch04-consumer" -workflow "$B/ch06" -stream-exercise "$B/ch07" -stream-model claude-sonnet-4-6 -gui "$B/gui" -gui-model claude-sonnet-4-6 -media-models "$W/media-models.json" "$W/bin/ch01"
go run ./cmd/grade-edition3 -ch 1 -json "$W/bin/ch01"
```

The owned public media map retains anthropic `claude-sonnet-4-6`, openai `gpt-5.4`, gemini `gemini-3.8-flash`; no discovery or capability changes. Chapter 12: **actual exit 0, parsed passed:true, 100/100**, including inherited Chapter 11 parity. Only then ran the first Chapter 1 gate on the same hash-verified CLI: **actual exit 0, parsed passed:true, 100/100**, both genuine 62-second delayed-header/body checks passed with complete answer and usage; observed elapsed time 124 seconds. No outer timeout was imposed (allowance exceeded 360 seconds); observation waits did not cancel execution. Reports/exits/stderr are `chapter12.{json,exit}`, `chapter12.stderr.log`, `chapter1.{json,exit}`, `chapter1.stderr.log`, plus `chapter1.timing.log`; `commands.log` records commands/exits. No new failures, retries, grader-source inspection, source edits, children, credentials/external providers, browser/native/device actions, installation, publication, cleanup, commits, or snapshots.

**Short operational proposal — not yet authorized or executed:**

1. In an owned isolated loopback GUI session (not live 8088), launch the held client with `--gui-debug --skills-dir` pointing to its actual skills and runtime saves/settings outside source. Confirm startup skill admission and browser attachment; do not add unsolicited DOM/TTS context.
2. Through the existing Agent/MCP path, explicitly request `gui_snapshot` and `tts_queue`, then `gui_input` on the draft and `gui_click` on a harmless visible control. Request another snapshot to confirm actual effects. Observe that programmatic input dispatches control updates without claiming human typing/pause; contrast a brief human whitespace input and explicit Continue. No native/device or paid-provider action is implied.
3. Unload the skill and confirm external calls are blocked/peer closed while pending instructions remain; optionally reload in the same page to complete that bounded lifecycle check. Keep observations concise and do not expand into prior-chapter/provider campaigns.

HOLD for root inspection and a separate operational release. Grader/local success is not a claim of real-browser or overall chapter acceptance.

## Operational preparation only — ready for root inspection, UNEXECUTED

Source and both validated CLI/report artifacts remain held. Built one GUI from held `gui/cmd` with the established clean offline environment/shared cache: `go build -o /tmp/e3-ch12-operations.AOMZI3/bin/gui ./cmd`, exit 0 (`gui-build.log` in that workspace). GUI SHA-256: `1832186dc65e4f2e0c7ee67b7cd9e3f9591ce07eab7985e6ac9574467b649347`.

One controller: `/tmp/e3-ch12-operations.AOMZI3/controller.mjs` (262 lines). Final SHA-256: `7edd6c2bab30daa884975d76df4d32062fe0912c33fc5ccf8cd9f1eecf52b8b2`. Node `--check` exit 0 (`syntax.log`); embedded artifact hashes verified read-only. No controller invocation, GUI launch, Chrome, native speech, provider/key read, models GET, or upstream request occurred. Preparation-only path/anchor misses and guessed protocol fields were corrected against current source; they are not operational failure evidence. No suites/graders repeated, children delegated, or production edits made.

**Local plan:** one owned dynamic-loopback GUI/proxy and private Chrome profile, actual `--gui-debug --skills-dir .../main/skills --gui-dir .../main/gui/web/gui --port 0 --settings OWNED --save OWNED`; current Agent/jobs/MCP path with exactly eight fake POSTs, ten-minute absolute cap. Explicit snapshot/queue, actual sidebar click and draft input, correlated fresh snapshot plus effective DOM/input+change observations. Automatic read-aloud remains false. Speak one short final card on demand through production speech callbacks; require native START→END and actual `tts_queue` capture, then automatic pause release with the programmatic draft retained. Contrast trusted CDP whitespace DOM input (`isTrusted=true`, not a claim of a physical person typing) and real Continue release. One phrase only; no retry/second phrase planned. Default voice locality unknown; no hearing claim. Unload must remove the next request's browser declarations, close the external tunnel, and retain pending instructions. No optional reload/layout/prior-chapter campaign.

**Live plan:** separate invocation, accepted only after matching local PASS. One `claude-sonnet-4-6` task, sequential snapshot→actual Settings→font20→fresh snapshot/effective card+draft font→unload→final/post-unload request, at most six TOTAL POSTs. Read-aloud false and Chrome muted; no native speech. At most 60s setup, 180s active, 240s absolute. After future execution approval, pins and owned PID/command/profile/listeners precede attachment/key/upstream; one exact `/Users/bill/.cr/settings.json` `directClaudeAPIKey` read and one `GET https://api.anthropic.com/v1/models`. Key stays in controller memory/HTTPS header; GUI has a dummy key/minimal environment. Exact `/v1/messages`, native `max_tokens:1024`, original body bytes, atomic request slots, no redirects/retries/fallbacks. TLS uses the pinned Node runtime's explicit root bundle and actual authorized handshake. Upstream timeout uses remaining case budget, not a blanket 3s. Inherited menu remains visible, but irrelevant returned tool effects and uncorrelated automatic browser observations are refused before forwarding the response. Stdin stays open until final→EOF→exit/save/full-usage-ledger verification. Stop-first-failure preserves partial logs/save/profile; cleanup targets only owned process groups/listeners, with profile retained for inspection.

Exact commands below are **NOT authorization to run**; root must separately release local, then live only after genuine local PASS:

```sh
env -i HOME=/tmp/e3-ch12-operations.AOMZI3/home TMPDIR=/tmp/e3-ch12-operations.AOMZI3/tmp PATH=/usr/bin:/bin:/usr/sbin:/sbin LANG=en_US.UTF-8 /opt/homebrew/Cellar/node@24/24.20.0/bin/node /tmp/e3-ch12-operations.AOMZI3/controller.mjs local --approved-local 7edd6c2bab30daa884975d76df4d32062fe0912c33fc5ccf8cd9f1eecf52b8b2

env -i HOME=/tmp/e3-ch12-operations.AOMZI3/home TMPDIR=/tmp/e3-ch12-operations.AOMZI3/tmp PATH=/usr/bin:/bin:/usr/sbin:/sbin LANG=en_US.UTF-8 /opt/homebrew/Cellar/node@24/24.20.0/bin/node /tmp/e3-ch12-operations.AOMZI3/controller.mjs live --approved-live 7edd6c2bab30daa884975d76df4d32062fe0912c33fc5ccf8cd9f1eecf52b8b2
```

Controller pins include the held GUI and CLI (`d723130…92f52f`), actual Chrome and Node executables, the three served assets and two current skill files, and Node's TLS roots. Static compilation/syntax/pin checks do not establish operational success; both cases await review and future execution approval. Last free-space check: 24,253,708 KiB, above the 1 GiB floor. Assistance: root's bounded public operational brief; no private grader/reviewer inputs or sub-agents.

### Versioned controller correction — still UNEXECUTED

At root's explicit static review request, created `/tmp/e3-ch12-operations.AOMZI3/controller-v2.mjs`; SHA-256 `bc6d0cd70ff1232c1917fd97cc6400213c69959c338a5c73a6538e3075fadc78`. Original `controller.mjs` preserved byte-for-byte at `7edd6c2bab30daa884975d76df4d32062fe0912c33fc5ccf8cd9f1eecf52b8b2`.

Only three corrections (four diff hunks): remove unsupported settings `version`; require the current primary prefix `Use the enabled tools to inspect facts before claiming results.` while retaining full gui-debug/manual and unsolicited-observation guards; both native waits now promptly reject actual `error/failed/unsupported/disabled/cancelled` states. Genuine START→END remains required. No unload-ordering change, new case, production/binary edit, or execution. Current settings/manual/speech owners were read to verify these anchors.

`node --check /tmp/e3-ch12-operations.AOMZI3/controller-v2.mjs` in clean environment exited 0; output `/tmp/e3-ch12-operations.AOMZI3/syntax-v2.log`. Original/new hashes and the bounded diff were checked; last disk check 24,218,032 KiB. This is syntax/static evidence only. All operational approval gates and budgets remain unchanged. Updated commands, **unexecuted and requiring separate root releases**:

```sh
env -i HOME=/tmp/e3-ch12-operations.AOMZI3/home TMPDIR=/tmp/e3-ch12-operations.AOMZI3/tmp PATH=/usr/bin:/bin:/usr/sbin:/sbin LANG=en_US.UTF-8 /opt/homebrew/Cellar/node@24/24.20.0/bin/node /tmp/e3-ch12-operations.AOMZI3/controller-v2.mjs local --approved-local bc6d0cd70ff1232c1917fd97cc6400213c69959c338a5c73a6538e3075fadc78

env -i HOME=/tmp/e3-ch12-operations.AOMZI3/home TMPDIR=/tmp/e3-ch12-operations.AOMZI3/tmp PATH=/usr/bin:/bin:/usr/sbin:/sbin LANG=en_US.UTF-8 /opt/homebrew/Cellar/node@24/24.20.0/bin/node /tmp/e3-ch12-operations.AOMZI3/controller-v2.mjs live --approved-live bc6d0cd70ff1232c1917fd97cc6400213c69959c338a5c73a6538e3075fadc78
```

### One authorized local invocation — STOPPED before product exercise

Ran the exact v2 `local --approved-local` command above once in the saved clean environment. Actual exit **1**, `pass:false`; stdout/stderr: `/tmp/e3-ch12-operations.AOMZI3/local-controller.log`, exit file `local-controller.exit`, structured result `/tmp/e3-ch12-operations.AOMZI3/local/result.json`, events `local/events.log`.

Failure: `owned command mismatch: /private/tmp/e3-ch12-operations.AOMZI3/controller-v2.mjs`. The command used the approved `/tmp/...` spelling; the controller self-path check expected `/private/tmp/...`. It stopped at its own proxy-ownership check before GUI/Chrome startup. Result: **0 fake POSTs, 0 models GETs, 0 usage, nativeStartEnd:false**, elapsed 9ms. Thus this is a controller startup/ownership failure, not behavioral evidence for or against the product. No native utterance, key read, external provider, or live case occurred.

Events record one owned proxy at 127.0.0.1:64632 (controller PID 88640), then stop and owned-only cleanup; the controller process exited. Partial settings/events/result and the unstarted private profile directory are retained. No retry, controller/source/binary edit, extra case, or shared cleanup. Pre-run free space: 24,213,208 KiB. Existing controller/GUI/CLI identities are recorded in the result. Returned HOLD for root verification/classification; live remains unauthorized and unrun.

### One authorized canonical-path v3 local invocation — partial evidence, overall FAIL

Assistance: root made only the mechanical v3 controller copy (new `local-v3`/`live-v3` case directories and matching local prerequisite), supplied SHA `285ab30bad9751563c26a450b1628069f827d6a1f34882c91ffcab6c19d442ed`, and authorized canonical `/private/tmp/.../controller-v3.mjs` argv. Student did not edit controller or source. V2 failed-local files remain untouched.

Executed once with the prior exact clean environment/Node: `/private/tmp/e3-ch12-operations.AOMZI3/controller-v3.mjs local --approved-local 285ab30bad9751563c26a450b1628069f827d6a1f34882c91ffcab6c19d442ed`. Actual exit **1**, `pass:false`, failure **`deadline: trusted whitespace typing pause`**, elapsed 22,461ms. Six fake POSTs, zero models GETs, recorded fake usage input66/output18/cache0. No retry, live/key/external-provider action, or controller/source/build change.

Completed partial path: owned dynamic GUI/Chrome/proxy identity checks; all four explicit browser tools; fresh snapshot showing the collapsed sidebar and actual programmatic draft; programmatic cause/native automatic-release gates; one on-demand 19-character phrase. Actual `tts_queue` (`events.log` lines26–27) records speech_id1 `speaking` at1792544527300, `chunk-ended`/`done` at1792544529207 (~1.907s), empty current/pending/request queues and native speaking/pending/paused=false. This is browser callback/queue evidence, **not hearing or local-voice proof** (`locality:unknown`). Result records nativeStartEnd:true but overall PASS remains false.

The subsequent trusted-whitespace pause condition did not establish before its deadline; cause remains unclassified, not asserted as a production defect. Continue, unload/capability withdrawal, requests7–8, and the successful final→EOF/full-usage validation were not reached. Cleanup records Chrome PID88676 and GUI PID88669 stopped, then owned-only cleanup; profile and partial save/output retained. No shared cleanup.

Evidence: `/tmp/e3-ch12-operations.AOMZI3/local-v3-controller.log`, `local-v3-controller.exit`, and `local-v3/{result.json,events.log,gui.stdout.log,gui.stderr.log,chrome.stderr.log,save.json}`. Owned listeners were proxy64635, GUI64636, Chrome64637; none8088. Pre-run free space24,177,996KiB. Returned stopped/HOLD for root verification; no live execution or automatic rerun follows.

### Remaining-only preparation — UNEXECUTED

Root/reviewer classification carried accurately: v3's trace establishes that typing did pause; the five-condition wait's unmet component was not captured. No product correction round3, source drift, or whitespace-regex defect is claimed. V3 remains overall FAIL. Its actual four-tool/programmatic-draft/native START→END/automatic-release observations are carried separately, not rerun. Root's withdrawn diagnosis is not an implementation premise.

Prepared `/tmp/e3-ch12-operations.AOMZI3/controller-v4-remaining.mjs` (211 lines), SHA-256 **`078ba14c966159c172bbbe1ddcab16e66a449f75b7c40e5a3bb2dde935cf8431`**. Exact diff: `controller-v4-remaining.diff` in the same workspace. Original/v2/v3 and their results preserved; v3 hash rechecked unchanged (`285ab30…d442ed`). Reused reviewed pins, fake SSE, CDP transport, owned process/listener checks, EOF/save/usage validation and cleanup. No production or held-binary changes.

New case only: fresh `remaining-v4` home/tmp/profile/settings/save, muted Chrome and read-aloud false. Fresh blank draft→real CDP DOM focus and `Input.insertText(' ')`→trusted whitespace pause→actual Continue preserving that draft. No select-all shortcut, forged event, handler/callback substitution, native utterance, or browser-tool replay. Records exact active-element/focus, selection start/end/direction, draft/value/length, trusted input/change records and client state before/after each action; failed waits log last view and each typing predicate, and every stop includes timestamped last observation.

Then real Agent/job `unload_skill` and final consumption with **at most two fake POSTs**, next-request browser-menu absence/pending manual/closed tunnel, normal final→EOF→exit/save/full usage. Three-minute budget with ten seconds reserved for reviewed cleanup; first failure stops and retains partial records. No live/HTTPS/key/models-GET path exists in this file. Fresh result explicitly sets nativeExercised:false/nativeStartEnd:false and links, but does not promote, the prior v3 evidence.

Preparation validation only: clean-environment `node --check .../controller-v4-remaining.mjs` exit0; log `syntax-v4-remaining.log`. No controller invocation, GUI/Chrome/native/provider/key/discovery execution or repeated cases. Last free-space check23,874,284KiB. Ready for root read; command below is **UNEXECUTED, not execution authorization**:

```sh
env -i HOME=/tmp/e3-ch12-operations.AOMZI3/home TMPDIR=/tmp/e3-ch12-operations.AOMZI3/tmp PATH=/usr/bin:/bin:/usr/sbin:/sbin LANG=en_US.UTF-8 /opt/homebrew/Cellar/node@24/24.20.0/bin/node /private/tmp/e3-ch12-operations.AOMZI3/controller-v4-remaining.mjs remaining --approved-remaining 078ba14c966159c172bbbe1ddcab16e66a449f75b7c40e5a3bb2dde935cf8431
```

Existing live code is unchanged and remains blocked by its original genuine-local-PASS requirement. Any later live release must explicitly gate on **both** root-verified carried v3 tool/programmatic/native receipts **and** a genuine fresh remaining-only PASS with matching held artifacts; it must not pretend native work was repeated or make fresh nativeStartEnd true. No live gate adaptation or paid release is made here. Assistance: root's public failure classification and remaining-only scope; no private review/grader input or children.

### One authorized remaining-v4 invocation — partial evidence, overall FAIL

Ran the exact canonical v4-remaining command above once. Actual exit **1**, `pass:false`, failure **`browser capability menu mismatch`** on the second inbound fake POST; elapsed1,565ms. Two fake POST slots consumed (v3 six + remaining two = eight total), but only the first response completed: recorded fake usage input11/output3/cache0. Zero models GETs, key reads, external-provider/live actions, native utterances or MIC activity. No retry/diagnosis campaign, controller/source/build edit, or replay of completed cases.

Per-predicate evidence is now retained: `remaining-v4/events.log` line11 has actual trusted `input`/`insertText`, focused TEXTAREA `#input`, draft=`" "`, length1, selection1..1, typing/inputOccupied/serverPaused=true; all five typing predicates true. Line14 records Continue with the same draft/length/selection, typing/inputOccupied/serverPaused=false and continue_epoch1. Physical-human typing is not claimed; this is trusted CDP DOM input through production handling. Native speech history empty and native speaking/pending/paused=false in those observations.

The next capability-menu guard failed before the second request was logged/served. Cause remains unclassified; no product-defect claim or harness repair is made here. Unload withdrawal/pending-manual/closed-peer verification, final consumption, and successful normal EOF/save/full-usage validation were **not established**. Partial save/output retained. Lines19/23 preserve stopped/failed-wait last views and predicate values. Lines20–22 record owned ChromePID89043 and GUIPID89036 stopped and owned-only cleanup; profile retained, no shared cleanup.

Evidence under `/tmp/e3-ch12-operations.AOMZI3/`: `remaining-v4-controller.log`, `remaining-v4-controller.exit`, and `remaining-v4/{result.json,events.log,gui.stdout.log,gui.stderr.log,chrome.stderr.log,save.json}`. Runtime pins/recorded identities: controller `078ba14c966159c172bbbe1ddcab16e66a449f75b7c40e5a3bb2dde935cf8431`; held GUI `1832186dc65e4f2e0c7ee67b7cd9e3f9591ce07eab7985e6ac9574467b649347`; held CLI `d72313023079544f5a31e11de978150506db1c32d0a1daba48890c1a1c92f52f`; source/assets/skills unchanged by this work. Pre-run free space23,878,360KiB.

Carried v3 tool/programmatic/native observations remain explicitly separate; its overall FAIL is not promoted. Fresh result correctly reports nativeExercised:false/nativeStartEnd:false and links prior v3 receipts. Returned saved result/HOLD for root verification; live remains unauthorized and unexecuted.

### Production correction ROUND3/10 — focused checks + compile HOLD

Narrow production delta only in `main/internal/engine/engine.go` and `main/internal/tool/skills.go`: supported Anthropic projection omits **disabled MCP-sourced** declarations, retaining ordinary dormant compiled catalog entries and automatic exclusions. Capability snapshot/admission, teardown, pending instructions/grants, reference selection, real IDs and canonical History are unchanged. Unload description/result now distinguish retained instructions/compiled grants from immediate external-grant removal and closure of peers with no remaining grants. Shared grants are not described as forcibly closed. This corrects a declaration leak/misleading wording, **not a successful callable-capability leak**.

Added `main/chapter12_unload_test.go`: real Agent, local Anthropic HTTP fixture and actual MCP transport, genuine model load/unload IDs. External-only and mixed compiled-grant variants verify next-request absence of all four external declarations, ordinary dormant/initial catalog entries, pending manuals/compiled grants, real paired status and canonical-history stability, ordinary reference retention/status fallback, actual peer close, automatic exclusion/cessation and external refusal. Existing fake-peer same-socket reload/refusal check also passes; no Chrome/native operational case is replayed.

Fresh owned workspace `/tmp/e3-ch12-r3.X9bcOh`; all tests/builds sequential in clean offline env (GOMAXPROCS=2, readonly/-p=1, local toolchain, GOPROXY/GOSUMDB off, existing shared GOMODCACHE/GOCACHE, owned HOME/TMP, no inherited keys/base URL). Disk checked before each execution; last pre-build free23,772,232KiB.

Actual commands/results (each log and numeric `.exit` in that workspace):
- Before production edits, root `go test -vet=off -count=1 -run '^TestChapter12UnloadAnthropicCatalog$' .` — **exit1**, `before.log`. Compiled regression reproduced all four disabled declarations with deferred=true in both variants, plus old inaccurate status/description. Admission/closure/history/reference checks did not fail. This is new compiled reproduction evidence; it is **not** a claim that the absent raw second body from remaining-v4 was captured.
- Root `go test -vet=off -count=1 -run '^(TestChapter12UnloadAnthropicCatalog|TestChapter11AutomaticToolsStayOutOfModelProjection)$' .` — **exit0**, `focused-core.log`. **Coverage correction from root:** only `TestChapter12UnloadAnthropicCatalog` matched. `TestChapter11AutomaticToolsStayOutOfModelProjection` does not exist; the actual existing test is `TestChapter11AutomaticAnthropicProjection`. The exact command/exit are retained, but this focused run is not credited with executing that existing test. The released full root race suite will cover it without a separate replay.
- GUI `go test -vet=off -count=1 -run '^TestMCPBrowserReloadSameSocket$' .` — **exit0**, `focused-mcp-reload.log`; existing local fake peer, not a browser/native session.
- Root and GUI `go build ./...` — **exit0** each, `build-main.log`, `build-gui.log`.
- `ch06`, `ch07`, `examples/ch04-consumer`: each `go build -o /dev/null .` — **exit0**, `build-workflow.log`, `build-stream.log`, `build-consumer.log`. Compile-only, not exercise execution.

Two wording edit anchors safely refused before the exact owner was reread; `tool.go` was unchanged. Final production diff is the projection predicate/comment and two strings in the skill-tool owner. No transport/lifecycle framework or controller/grader weakening. Assistance: root and retained reviewer's public classification/scope; no children/private notes or grader input.

All earlier failures/controllers/results and R2 frozen CLI/GUI remain untouched; those binaries do **not** contain R3. No new operational binary, full suite/race/vet/grader, GET/key/external provider/browser/native execution, cleanup, source commit or snapshot. V3 four-tool/programmatic/native facts and v4 trusted-whitespace/Continue facts carry with their overall FAIL qualifications, not replay or promoted acceptance. Source frozen at coherent focused-check/compile HOLD for root review.

### R3 validation release — gates PASS, operations still pending

Corrected the prior R3 focused coverage receipt in place while preserving its exact command/exit: the nonexistent selector did not execute an existing test. Fresh root **`go test -race ./...` exit0** now covers the real `TestChapter11AutomaticAnthropicProjection` together with the new regression; **`go vet ./...` exit0**. Logs/exits: `/tmp/e3-ch12-r3-validation.Nmod9X/root-{race,vet}.{log,exit}`. Sequential clean offline/shared-cache environment as above, owned HOME/TMP and disk checks. GUI/application suite results are carried only as root authorized from verified unchanged adapter/module inputs, R3 focused GUI reload and current five-module builds; those suites were not repeated.

Built exactly one new frozen CLI and GUI, both exit0; build logs/exits and `.sha256` records in the new workspace. No GUI process:
- CLI `/tmp/e3-ch12-r3-validation.Nmod9X/bin/ch01` — **`4421938b217c257b6fa67966a085cacfdef67a5c757a5bdc257fcf4afdb414a1`** (`go build -o .../bin/ch01 ./cmd/ch01` from main).
- GUI `/tmp/e3-ch12-r3-validation.Nmod9X/bin/gui` — **`767910b6290dda7c12d3c0f697f4b56b3acc50d9c32ff637cd63d54368744fdd`** (`go build -o .../bin/gui ./cmd` from gui).
Old R2 binaries, failure reports and all controller/operational records untouched. No production edits during validation.

From book root, with `B=/Users/bill/projects/ensemble/solutions/edition-3/main`, `W=/tmp/e3-ch12-r3-validation.Nmod9X`, actual commands:
```sh
go run ./cmd/grade-edition3 -ch 12 -json -core "$B" -consumer "$B/examples/ch04-consumer" -workflow "$B/ch06" -stream-exercise "$B/ch07" -stream-model claude-sonnet-4-6 -gui "$B/gui" -gui-model claude-sonnet-4-6 -media-models "$W/media-models.json" "$W/bin/ch01"
go run ./cmd/grade-edition3 -ch 1 -json "$W/bin/ch01"
```
Unchanged public media map copied as plain JSON (no discovery). **Chapter12 actual exit0, parsed passed:true, 100/100 including corrected Chapter11 parity.** Only afterward, first R3 Chapter1 on the **same hash-verified CLI**: **actual exit0, parsed passed:true, 100/100**. Genuine62-second delayed headers and62-second delayed body each delivered complete answer+usage; elapsed124s. No outer timeout/cancellation (allowance unbounded, exceeding360s). Ordinary full usage input511/output65 matched fixture totals. Raw expanded commands, stdout JSON, stderr and exits retained as `chapter{12,1}.{command,json,stderr.log,exit}`; Chapter1 timing in `chapter1.timing.log`. Last pre-command free23,663,100KiB.

No grader/private-source inspection or grader/guidance/source change, controller preparation, real browser/native/MIC, key read/models discovery, live/external provider, cleanup, child delegation, commit or snapshot. V3 tool/programmatic/native and v4 trusted-whitespace/Continue/teardown facts carry with their failed-run qualifications; none was replayed or promoted to overall PASS. **HOLD:** current gates pass, but successful R3 real-GUI post-unload declaration/consumption/EOF closure and separately authorized live exercise remain operational boundaries. No execution or preparation beyond this release.

### R3 live-only completion — PREPARATION ONLY / HOLD

Prepared `/tmp/e3-ch12-operations.AOMZI3/controller-v5-live.mjs` (269 lines), exact v3→v5 diff `controller-v5-live.diff`. SHA256 **`b3bdaa10dd6526170cefaa71ba774187af59030e54cb51c1a0d2ebe41288cd00`**. `node --check` **exit0**, `syntax-v5-live.log`/`.exit`; read-only hashing verified all17 embedded artifact/receipt pins. Original v3 remains `285ab30bad9751563c26a450b1628069f827d6a1f34882c91ffcab6c19d442ed`; all earlier controllers, failed results and R2 binaries preserved. No source/binary/gate change or repeated suite/grade.

One live-only path; local/native/fake execution removed. New fresh case will be `live-r3-v5`, with canonical script argv. Pins NEW R3 CLI4421938b…/GUI767910b… plus unchanged assets/manuals, Node, Chrome and Node's explicit TLS-root bundle. Replaces impossible whole-local-PASS gate with pinned exact v3/v4 records and semantic checks of real paired four-tool results, programmatic effect, actual tts_queue START→END/drain and subsequent Agent progress, and trusted whitespace/Continue preserving draft; both prior overall statuses are required to remain **false**. Also requires actual current R3 report passed:true, exit0 receipts and matching binary/hash/command identities. These checks carry older-GUI observations, **not whole-GUI/native proof on R3**. Current result will explicitly say nativeExercised:false/nativeStartEnd:false, with separately qualified carried evidence.

Prepared task: explicit snapshot → actual Settings click using an observed visible selector → another snapshot for newly visible controls → gui_input font20 → fresh correlated snapshot plus actual control/CSS/card font20 → actual unload → following request with all four browser declarations absent and pending manual/peer teardown → real returned final text rendered, then EOF→exit0/save/full usage. No arbitrary final answer phrase. Existing paired/no-unsolicited-observation guards retained; inherited menu may remain, but unrelated effects are refused before forwarding model tool calls to the Agent. Sequential tool rounds; **at most8 total POST slots** (usual sequence7), native1024/stream/body bytes unchanged. One exact model GET and one exact memory-only directClaudeAPIKey read are code paths reserved solely for a later execution release. Fixed api.anthropic.com TLS/SNI/verified pinned interpreter bundle, no redirects/retries/fallback, no blanket short upstream timeout. <=180s active/60s setup/240s absolute, with10s reserved for owned cleanup. Owned muted Chrome, read-aloud false, no native/MIC/OS-focus/volume/shared8088; minimal dummy-key GUI env. Known ownership/CDP/cleanup reused; last observed view retained on failure. Private owned profile retained as evidence; no shared cleanup.

**NOT EXECUTED — only after separate root inspection and execution approval:**
```sh
env -i HOME=/tmp/e3-ch12-operations.AOMZI3/home TMPDIR=/tmp/e3-ch12-operations.AOMZI3/tmp PATH=/usr/bin:/bin:/usr/sbin:/sbin LANG=en_US.UTF-8 /opt/homebrew/Cellar/node@24/24.20.0/bin/node /private/tmp/e3-ch12-operations.AOMZI3/controller-v5-live.mjs live --approved-live b3bdaa10dd6526170cefaa71ba774187af59030e54cb51c1a0d2ebe41288cd00
```

Preparation consisted only of reading the reviewed controller/current snapshot/header owners and owned public receipts, writing this version/diff, read-only hashes and syntax check. No controller import/case execution, live directory, GUI/Chrome/native/MIC, key read, GET/discovery, external provider or paid request. No children, new binaries, production changes, cleanup, commit/snapshot or separate harness project. Root supplied the combined-completion scope and qualified-carry rule. **HOLD for complete-script review; syntax/ready is not an operational PASS.**

### Authorized ONE v5 live execution — actual PASS / final HOLD

Executed the exact canonical clean-env command above once, with unchanged controller SHA **b3bdaa10dd6526170cefaa71ba774187af59030e54cb51c1a0d2ebe41288cd00** and held R3 GUI **767910b…** / CLI **4421938b…**. Outer `live-v5-controller.log`/`.exit` in `/tmp/e3-ch12-operations.AOMZI3`: **exit0**. Actual case `/tmp/e3-ch12-operations.AOMZI3/live-r3-v5/result.json`: **pass:true**, failure:null, elapsed **19,085ms** including owned cleanup. Runtime pins/ownership/prerequisite checks passed. No controller/source/binary changes or retries.

Actual authorized upstream use: **one models GET + seven total POSTs**, all HTTP200; exact `claude-sonnet-4-6`, every request native `max_tokens:1024`, streaming true, original bytes forwarded unchanged. Approved single memory-only directClaudeAPIKey path used; no additional key/network access by the controller. Verified TLS connection recorded authorized and reused. Within60s setup/180s active/240s absolute and8-POST ceiling. No native utterance/MIC/OS-focus/volume operation, broader tests or local-case replay; owned Chrome muted and read-aloud false.

Read-only receipt inspection independently confirmed seven actual request/response records and **six real tool/result ID pairs** in order: snapshot → observed-selector Settings click → snapshot of newly visible controls → font20 input → fresh snapshot → unload → seventh response end_turn. Fresh snapshot at the actual selected control records value20; simultaneous actual control, shared setting, CSS variable, draft and conversation-card computed fonts are20px. Owned `settings.json` persists font_size20/tts_enabled:false. No raw model/browser payload is copied into this note; full evidence remains private in `events.log`, `save.json` and owned settings/profile files.

Actual unload status distinguishes pending instructions/compiled grants from immediate external grants/removal of unused peers. Following request contains **none of the four browser declarations** and retains the **full pinned392-character gui-debug manual**. Final browser view: MCP ready:false, gui-debug enabled:true/pending_unload:true, font20, idle/unpaused, MIC/read-aloud off, speech history empty. Actual returned final text rendered; then normal EOF→GUI **exit0**→save. Saved context idle; **seven response_ended records, all usage complete**, real tool_dispatched/job_report records, no automatic browser records. Controller/API aggregate, normal-exit and saved usage agree: **input38,979 / output702 / cache_read0 / cache_write0** (total39,681). Cleanup records owned Chrome stopped, GUI already normally exited, proxy/controller closed; private owned profile retained, no shared cleanup.

Two initial read-only extraction labels guessed a CSS ID and obsolete manual phrase and therefore missed those fields; corrected against the **actual saved selector** and **full pinned manual body**, both confirmed. This was receipt extraction only, not a product/controller failure or rerun. Outer execution exit and all private partial/full evidence remain preserved.

Fresh result explicitly nativeExercised:false/nativeStartEnd:false. V3 overallFAIL and v4 overallFAIL remain unchanged; their older-GUI four-tool/programmatic/native callback/automatic-release and trusted-whitespace/Continue facts are carried separately, not promoted to whole-GUI/native R3 proof. Native callback evidence still makes no human-hearing or voice-locality claim. No secret/private payload in Git; only this concise note changed. No children, source/grader/guidance/manuscript changes, commits/snapshots, extra preparation or operations. **HOLD for root verification of the actual live result and qualified combined evidence.**
