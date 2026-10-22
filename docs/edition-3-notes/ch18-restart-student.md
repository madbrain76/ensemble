# Chapter 18 restart — round0 compiling HOLD

Implemented directly from accepted main/prerequisite/public material; no separate design HOLD. No children, private reference/reviewer material, live services, credentials, commits, snapshots, installs, or shared-process/device operations. Existing `.ch15-work` / `.ch16-work` scratch left untouched. The requested `.agents/skills/coding/SKILL.md` was absent; read `.agents/skills/ensemble-coding/SKILL.md` as directed by the Edition 3 role. Read the governing public documents and recovered clipped chapter sections. Parent assistance: task boundaries and a question about sharing appropriate contracts through common; no implementation supplied.

## Owners and files

- `internal/common/{mcp,common,actor,cache,environment}.go`, `ensemble.go`: URL/auth-reference contract and shared URL validation; agent model-control method and typed observation; selected/last-producing model fields on existing typed EngineStatus; child-environment filtering of MCP credential references. No copied usage counters or new owner/backpointer framework.
- `internal/mcp/{http.go,http_test.go,manager.go,peer.go,transport.go}`: Streamable HTTP below existing peer/correlation/registration path; JSON/SSE, session/protocol headers, environment bearer reference, redirect refusal, sanitized HTTP errors, context-aware request cancellation, and owned connection closure. Existing stdio retained. No client-wide HTTP timeout.
- `internal/skills/parse.go`, `skills/{ensemble,web-search}/SKILL.md`: URL declaration parsing/validation; genuinely loadable Firecrawl-backed web skill with independent `FIRECRAWL_API_KEY` reference; status in primary skill.
- `internal/tool/{status.go,tool.go,command.go}`: `agent_status` reaches ToolCall.Parent → Agent.Engine → Engine.Status; local command environment uses shared credential filtering.
- `internal/engine/{engine,accounting}.go`: selected model and last producer added to existing per-producer accounting; existing known/partial/unknown cost semantics retained.
- `internal/agent/{model.go,actor.go}`: idle-boundary, same-provider model change via urgent actor control; validates before changing model fields only; preserves history and prior accounting; reuses OpenAI normalization.
- `internal/cli/{chat,cli}.go`, `gui/{server,hub,cmd/main}.go`, `gui/web/gui/{app.js,index.html}`: CLI `/status`, `/model`; GUI model/status controls and live typed meter; optional embedded loopback `--mcp-listen`.
- `mcp.go`, `mcp_listener.go`: dedicated and embedded listener compositions share an explicitly owned, cancellable/joined listener handle. Existing identified-observer unsubscribe ownership preserved.
- `chapter18_test.go`: real public tool/actor/web-manifest tests with owned loopback fakes; `chapter03_tools_test.go`: default registry expectation includes the new tool without broadening explicit allowlists.

## Actual validation

All Go commands used `/private/tmp/e3-ch18-restart-student/offline`, which executes `env -i` with:

```
PATH=/usr/local/go/bin:/Users/bill/go/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin
HOME=/private/tmp/e3-ch18-restart-student/home
TMPDIR=/private/tmp/e3-ch18-restart-student/tmp
GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local
GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod
GOCACHE=/Users/bill/Library/Caches/go-build
ANTHROPIC_API_KEY=dummy OPENAI_API_KEY=dummy GEMINI_API_KEY=dummy FIRECRAWL_API_KEY=dummy
```

No inherited account settings or generic `LLM_BASE_URL`. Commands below ran sequentially from `solutions/edition-3/main` unless noted; prefix each with the wrapper above.

- `go test -vet=off ./internal/mcp -run '^TestHTTPChapter18' -count=1 -timeout=45s` — exit 0.
- `go test -vet=off . ./internal/mcp -run '^(TestChapter18|TestHTTPChapter18|TestChapter3Round1Registry)' -count=1 -timeout=60s` — exit 1: new actor method omitted required mailbox urgency argument; MCP package passed. Corrected to existing urgent-control convention.
- `go test -vet=off . -run '^(TestChapter18|TestChapter3Round1Registry)' -count=1 -timeout=60s` — first exit 1: extracting listener removed `net` import still required by `MCPPipe`; restored. Next exit 1: test compared history after invoking the real status tool, whose activity legitimately records history. Corrected test to compare immediately after model change. Final invocation exit 0.
- `go test -vet=off . -run '^TestChapter18StatusAndModel$' -count=1 -timeout=60s` — exit 0 after narrowing model-config mutation to owned fields.
- `go test -vet=off . ./internal/mcp -run '^(TestChapter18PublicWebSkill|TestHTTPChapter18)' -count=1 -timeout=60s` — exit 0 after configurable credential-filter integration.
- `go test -vet=off . -run '^TestChapter18PublicWebSkill$' -count=1 -timeout=60s` — exit 0 with custom `CH18_WEB_SECRET` reference, proving filtering is not just the shipped-name blacklist.

Final compilation-only command, from each of `.`, `gui`, `ch06`, `ch07`, `examples/ch04-consumer`, `examples/ch13-speech`:

```
go test -c -vet=off -o /private/tmp/e3-ch18-restart-student/compile/<label>/ ./...
```

Labels are the module paths with `/` and `.` replaced by `_`. All six exited 0. No compiled test binaries executed. Exact logs and result table: `/private/tmp/e3-ch18-restart-student/compile/`. Final free disk observation: 21 GiB.

Setup/edit issues distinct from product failures: several guarded edit anchors refused mismatches and were corrected by rereading current source; one early gofmt command used repository-root cwd and failed to find the module-relative file (exit 2), then was rerun correctly. Formatting performed on changed Go files. `git --no-pager diff --check -- solutions/edition-3/main` passed at source-review stage; later narrow environment/model/test edits were gofmt'd.

## Deferred and HOLD

No full suite, race, vet, numerical grade, lifetime/stress test, GUI/browser/native/device exercise, live provider/web/MCP probe, or operational/account discovery. Focused tests establish only the exercised paths; compilation is not implementation acceptance. Root source comparison and retained independent review precede broader authorized gates. Source editing/testing stopped after six successful compilation checks. Steady round0 compiling HOLD.

## R1/10 — corrected public blockers; compiling HOLD

Assistance: coordinator supplied three public review blockers and bounded regression requirements; no private review/reference read and no child workers. Existing owners retained; no redesign.

Production changes: `internal/mcp/http.go` and `transport.go` resolve the connection's bearer reference once, log a nonsecret missing-variable/anonymous-attempt diagnostic through the existing writer, and actually POST without Authorization when absent. Non-2xx errors retain bounded useful remote reason text. HTTP error redaction covers literal/JSON-escaped bearer text and decoded RPC error/tool-isError structures for JSON and SSE, preserving IDs, successful content, tool content shape, and numeric precision. No client-wide timeout, new correlation map, or stdio semantic change. `gui/web/gui/app.js` now shares one current-selected-model projection for label, editor, and tree title; controls/state updates no longer restore stale replay configuration. Old producing-model identity remains in the live status tooltip.

Regressions: `internal/mcp/http_test.go` covers real anonymous headers + named diagnostic, bounded reflected-token 401/403 including long JSON, decoded JSON/SSE RPC and nested tool-error content. `internal/cli/chapter18_test.go` uses existing `runChat` with actual `/model` and `/status` commands, producer/cost retention, and invalid cross-provider selection. `gui/chapter18_test.go` + `.cjs` send actual WS `set_model`, `agent_status`, then `pause`; captured frames feed production JS through the existing Node fixture and assert current labels/editor/tree plus old producing-model tooltip after controls/state updates.

All following Go commands used the existing clean offline wrapper, sequentially, from main except GUI commands from main/gui. No broad gates or lifetime tests run.

- `go test -vet=off ./internal/mcp ./internal/cli -run '^(TestHTTPChapter18R1|TestChapter18Human)' -count=1 -timeout=60s` — exit 1 before fixes: expected HTTP failures (no anonymous request/diagnostic, missing remote reason, JSON/SSE RPC/tool credential leaks); CLI passed. Receipt `r1-root-before.log`.
- `go test -vet=off . -run '^TestChapter18WSModelProjection$' -count=1 -timeout=60s` — initial test-setup compile failure (nonexistent ensemble.State alias); narrowed decoder to the actual paused wire field. Rerun exit 1 before production fix: production JS label/tree still showed `claude-sonnet-4-6` instead of selected `new-private-model`. Receipt `r1-gui-before.log`.
- `go test -vet=off . ./internal/mcp ./internal/cli -run '^(TestChapter18|TestHTTPChapter18)' -count=1 -timeout=60s` — root/public accounting+web and CLI passed; exit 1 from anonymous test's nil result destination after POST began succeeding. Fixed only that fixture destination. Receipt `r1-root-after.log`.
- `go test -vet=off ./internal/mcp -run '^TestHTTPChapter18' -count=1 -timeout=60s` — exit 0 after fixture correction (`r1-mcp-after.log`); exit 0 again after adding the concrete long escaped-JSON error regression/fallback (`r1-mcp-final.log`).
- `go test -vet=off . -run '^(TestChapter18WSModelProjection|TestChapter15R1RealUsageWireAndConsumer)$' -count=1 -timeout=60s` — exit 0 (`r1-gui-after.log`).
- Compilation only, from main and main/gui: `go test -c -vet=off -o /private/tmp/e3-ch18-restart-student/r1-compile/<label>/ ./...` — both exit 0, no binaries executed. Labels `_` and `gui`; result table/logs in `r1-compile/`.
- `git --no-pager diff --check -- solutions/edition-3/main` — exit 0. Final `df -h .`: 21 GiB free.

All R1 receipts above are under `/private/tmp/e3-ch18-restart-student/`. Guarded edit anchors stopped on spacing/exact-expression mismatches; recovered by rereading the actual source, with no behavior test run on the partial edits. The four unchanged auxiliary-module compile receipts are carried from round0, not rerun. Stdio/cancellation ownership was not changed; broader/lifetime validation remains deferred. Root/public and GUI focused passing results were reused where inputs were unchanged. Source edits/testing stopped after successful affected compilation. R1 compiling HOLD is not acceptance; coordinator and retained reviewer own the next gate.

## R2/10 — bounded undecodable-error redaction; compiling HOLD

Coordinator supplied one public defect: bounded non-2xx JSON can end mid-string and retain alternate Unicode-escaped bearer text. No private material or workers used. Only production change is the fallback in `internal/mcp/http.go`: after literal/canonical redaction, an undecodable diagnostic with remaining escape syntax is conservatively omitted. The existing HTTP status/reason fallback remains available. Bounded reads, fully decoded errors, successful responses, JSON/SSE correlation, stdio, and GUI/CLI code are unchanged. This intentionally does not promise faithful display of arbitrary malformed/truncated bodies.

Added `TestHTTPChapter18R2TruncatedUnicodeBearer` in `internal/mcp/http_test.go`: actual loopback 403 returns complete valid JSON, with bearer `secret` represented using single-backslash Unicode escapes and 6,000 filler bytes before closing JSON. The test checks the real request's Authorization, recoverable-token absence, bounded output, and safe status/reason retention.

Commands used the same clean offline wrapper, sequentially:
- From main: `go test -vet=off ./internal/mcp -run '^TestHTTPChapter18R2TruncatedUnicodeBearer$' -count=1 -timeout=60s` — exit 1 BEFORE production correction, demonstrating recoverable `\u0073\u0065\u0063\u0072\u0065\u0074` in the surfaced HTTP diagnostic. Saved `/private/tmp/e3-ch18-restart-student/r2-http-before.log`.
- From main: `go test -vet=off ./internal/mcp -run '^TestHTTPChapter18' -count=1 -timeout=60s` — exit 0 AFTER correction; includes existing useful 401/403, long canonical JSON, decoded JSON/SSE RPC/tool-error safeguards. Saved `r2-http-after.log` in the same directory. Edit script emitted a harmless Python invalid-escape-literal warning; source was inspected and gofmt succeeded, with no Go setup/build failure.
- From main and main/gui: `go test -c -vet=off -o /private/tmp/e3-ch18-restart-student/r2-compile/<label>/ ./...` — both exit 0, labels `_` and `gui`. No binaries executed; logs/result table saved under `r2-compile/`.
- `git --no-pager diff --check -- solutions/edition-3/main` — exit 0; final disk observation 21 GiB free.

Unchanged R1 GUI/CLI/source evidence and round0 auxiliary-module compilation receipts carried forward, not rerun. No broad/full/race/vet/grade/lifetime/live/provider/browser/credential/child-worker activity. Source editing/testing stopped after affected compilation; stable R2 compiling HOLD for coordinator/reviewer, not acceptance.

## Test-only repair after stopped root-race gate — remains R2/10

Root authorized fixture/expectation repair after its reviewer classification; no production round or production edits. Did not access the private harness or root gate log. Updated only `chapter03_loop_test.go` (14-tool call count, named `agent_status` membership, and `want == 14` guard), `chapter03_schema_test.go` (14 projected/default declarations), and `internal/engine/cache_test.go` (three invalid nil-parent literals replaced by `New(parseAgent{config: ...})` with checked errors). All provider/schema/order/no-mutation and cost/unknown/bucket/alias assertions preserved; no production nil guard.

From main, same clean offline wrapper, focused race command only:
`go test -race -vet=off . ./internal/engine -run '^(TestChapter3LoopOrderedBatchProviders|TestChapter3GeminiParameterProjection|TestChapter15CostsKeepCapturedIdentity|TestChapter15WeightedRatesAndUnknownPrices)$' -count=1 -timeout=60s` — exit 0; both packages passed. Receipt `/private/tmp/e3-ch18-restart-student/r2-test-repair-race.log`.

Inspected/saved exact three-file `git --no-pager diff` as `/private/tmp/e3-ch18-restart-student/r2-test-repair.diff`; `git --no-pager diff --check` for those files passed. Disk 21 GiB free. No further suites/gates or repeated compile-only checks; held production and existing compilation evidence unchanged. Test edits/validation stopped; stable R2 HOLD for root to inspect and resume its stopped gate.

## Standalone DOM fixture repair — remains R2/10

Root/retained reviewer classified missing page `$` in `gui/mcp_dom_test.cjs` as an isolated test fixture mismatch; production/full-page evidence carried unchanged. Test-only edit supplies the two page-control listener targets and `$` helper in the VM, separate from snapshot nodes. Still evaluates actual production source from the historical MCP marker through EOF; every snapshot/state/UTF-8/elision/reconnect/JSON-RPC assertion remains unchanged. No production relocation or edits.

From main/gui, same clean offline wrapper: `node mcp_dom_test.cjs` — exit 0, all 3 tests passed. Receipt `/private/tmp/e3-ch18-restart-student/r2-dom-fixture.log`; inspected exact diff saved as `r2-dom-fixture.diff` alongside it. `git --no-pager diff --check -- solutions/edition-3/main/gui/mcp_dom_test.cjs` — exit 0. Disk observed 21 GiB free before the test. No Go gates/builds, other standalone scripts, grades, or operational probes rerun. Test editing/validation stopped; stable R2 HOLD for root to continue remaining root-owned gates.

## R3/10 — unknown-only unavailable-tool refusal; compiling HOLD

Root supplied public Book18 lines445–449 conformance finding; read that requirement and actual dispatch owner, no private grade/review access. Sole production edit: `internal/tool/dispatch.go` shared unavailable-entry branch now reports `unknown tool %q` instead of `unknown or disabled tool %q`. No grant/dispatch/model logic changes. Updated two inherited wording assertions in `chapter03_loop_test.go`. Added `TestChapter18UnavailableToolsAreUnknown` in `chapter18_test.go`: same name called while nonexistent and then registered-but-ungranted; exactly equal error responses, unknown-only language, no model declaration, and atomic zero-handler-execution assertion through the real Agent registry.

Same clean offline wrapper, sequential commands from main unless stated:
- `go test -race -vet=off . -run '^TestChapter18UnavailableToolsAreUnknown$' -count=1 -timeout=60s` — exit 1 BEFORE fix, reproducing availability-disclosing wording. Receipt `/private/tmp/e3-ch18-restart-student/r3-unknown-before.log`.
- `go test -race -vet=off . -run '^(TestChapter18UnavailableToolsAreUnknown|TestChapter3Round1Registry|TestChapter3LoopOrderedBatchProviders|TestChapter3LoopEmptyRegistryProviders|TestMCPExplicitOwningGrants)$' -count=1 -timeout=60s` — exit 0 AFTER fix; receipt `r3-unknown-after.log` alongside it.
- From main and main/gui: `go test -c -vet=off -o /private/tmp/e3-ch18-restart-student/r3-compile/<label>/ ./...` — both exit 0; labels `_` and `gui`, logs/results in `r3-compile/`. No test binaries executed.
- Targeted `git --no-pager diff --check` passed; inspected production diff is exactly the one-line message change. Final disk21GiB free.

No test setup failures, full-suite/grade/browser/live-provider operations, child workers, or unrelated source changes. Root-owned GUI fixture task untouched. Earlier unrelated gate evidence carried; all broader gates remain root-owned. Source editing/testing stopped after successful affected compilation; stable R3 compiling HOLD, not acceptance.

## R4/10 — legal Anthropic tool-cache marker; compiling HOLD

Root supplied public Anthropic tool-search prompt-caching rule and current owner classification: deferred tools cannot carry cache_control. No private/live artifacts read and no live verification. Sole production change is `internal/engine/render.go`: scan actual Anthropic catalog backward and mark the last non-deferred declaration, or omit the tool marker if none exists. No catalog reorder/removal, defer/grant/reference/history changes, or global cache disablement; other marker placement unchanged.

Tests: extended existing `ch11_projection_test.go` real Agent→loopback provider fixture to check post-load request legality/last-eligible marker/budget, unsupported fallback, and identical full ordered catalog/schema/defer flags before/after disclosure. Existing real load pair/reference-only projection/grants/history/no-mutation assertions retained. Added `TestChapter18ToolCacheBreakpointEligibility` in `internal/engine/cache_test.go` for trailing deferred, eligible final, all-deferred, and empty boundaries, preserving catalog data/order and system/history marker budget. Reused relevant existing content/reference/marker/fallback tests, not a new harness.

Same clean offline wrapper, sequential commands from main:
- `go test -vet=off . ./internal/engine -run '^(TestChapter11AutomaticAnthropicProjection|TestChapter18ToolCacheBreakpointEligibility)$' -count=1 -timeout=60s` — exit 1 BEFORE production fix. Actual second HTTP request showed a deferred `agent_status` declaration carrying cache_control; trailing-deferred and all-deferred renderer cases also failed. Separate test mistake required markers on unsupported fallback; corrected to require zero, matching existing behavior. Receipt `/private/tmp/e3-ch18-restart-student/r4-cache-before.log`.
- `go test -vet=off . ./internal/engine -run '^(TestChapter11AutomaticAnthropicProjection|TestChapter18ToolCacheBreakpointEligibility|TestChapter15LegalMarkersPreserveContent|TestSkillToolReferencesAreContentNotInstructions|TestChapter09Round2OrdinaryResults|TestChapter15R1MarkerPositions|TestChapter15R1ActualHandoffAndFallbackPositions)$' -count=1 -timeout=60s` — exit 0, both packages passed. Receipt `r4-cache-after.log` alongside it.
- From main and main/gui: `go test -c -vet=off -o /private/tmp/e3-ch18-restart-student/r4-compile/<label>/ ./...` — both exit 0; labels `_` and `gui`, result table/logs in `r4-compile/`; no test binaries executed.
- Targeted `git --no-pager diff --check` passed; inspected production diff limited to eligibility scan/comment. Final disk20GiB free.

Connection-reset interruption was infrastructure of unproven cause: recovered the already completed before-fix job, did not repeat it or restart work. No Go setup failure; product failure and fixture expectation error separated above. No broad gates/grade/live/provider keys/browser/children/commit. Earlier unrelated evidence carried. Source editing/testing stopped after affected compilation; stable R4 compiling HOLD, not acceptance.
