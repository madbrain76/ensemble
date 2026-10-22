# Chapter 18 restart — retained independent code review

## ENTRY PASS — bounded alignment, 2026-10-22

The root's stable entry delta is suitable for assigning Chapter 18. This is **not** student-code, whole-chapter, browser, provider, or live-OAuth acceptance. No student implementation has been reviewed yet. Chapter 17 is the closed historical field report; the source predecessor remains accepted Chapter 16, not a Chapter 17 snapshot.

Scope: source-only independent review under the current lean ten-round rules. Read the current governing instructions, coding/reviewer/orchestration roles, execution/readiness, student brief, carryover, and full Chapter 18, recovering clipped ranges. Inspected the actual Edition 3 command/tests; Chapter 18 grader, driver, peer and relevant build/environment helpers; accepted CLI/GUI/configuration, owner interfaces, render/accounting and MCP paths. The original `solutions/ch18` was a simplicity comparison only, never executed as a positive control or supplied as a student answer. Reviewer executed no commands, network requests, browser/provider operations, installs or commits and edited only this note.

## Independent findings and their disposition

1. **Adapter defect, resolved:** historical `Chapter18(guiBin, core, ...)` builds the source root as an executable. Accepted main is a library with nested `cmd/ch01`. Edition 3 now supplies distinct CLI, GUI and core identities through `Chapter18Targets` (`internal/gradeedition2/ch18.go:33`, called by `cmd/grade-edition3/main.go:216`); required tariffs and unrelated flags are validated. The historical caller remains available. The same eleven checks and 155 points remain; no fake root executable or layout change is required.
2. **Fixture launch defects, resolved:** status GUI startup previously supplied no checked assets and inherited ambient process configuration. The driver now creates checked inert WS assets and reuses `invoiceEnv`. `webStatusEnv` (`internal/gradeedition2/ch18_driver.go:248–265`) removes `CH07_AGENT_LOG`, `CH02_LOG`, `E3_TRACE` and the inherited fixture token `EDITION2_WEB_TOKEN`. The first three are actual accepted CLI consumers at `solutions/edition-3/main/internal/cli/cli.go:61–63,88`; web launch then supplies its owned token. Plain text on GUI stdin is clear and supported. Importantly, accepted CLI `internal/cli/protocol.go` already accepts both historical JSON wrappers; those were **not** a demonstrated CLI defect.
3. **Controlled model/price mismatch, resolved:** accepted `solutions/edition-3/main/internal/common/cache.go:88–98` defines `LookupPrices`, which prices Sonnet 4.6 and Opus 4.6 on Anthropic Messages, not historical Sonnet-5 packets. The new E3 path uses explicit Sonnet 4.6 producer identity and, only when a second tariff is supplied, explicit Opus 4.6. Without that tariff the middle control remains genuinely unpriced. `statusLaunchModels` constructs `fakevendor.Reply.JSONBody` packets via `statusReply` for tool calls and `invoiceAnthropicUsage` for final text. Both put the controlled producer in the JSON `model` field, important because accepted `internal/engine/accounting.go:40–47` prefers returned producer identity for `LookupPrices`. Rates are independently supplied, not read from student code. These are documented controlled fixtures, not aliases, new capability rows, live availability claims, or student pricing requirements.
4. **Reference-name mismatch, resolved:** accepted tools receive `common.ToolCall.Parent()` and use `Agent.Engine()`, `Engine.Parent()` and `Engine.Status()`; there is no accepted `ToolContext`/`UsageMeter` naming obligation. The chapter now explicitly permits equivalent accessors and the actual singular `internal/tool` registry. The resolved-import gate accepts singular/plural registry paths while retaining the common-only application dependency boundary and failure when neither registry exists. Its tests include forbidden imports and the missing-registry case. No compatibility facade or speculative public probe was added.
5. **Manual invocation and authority defects, resolved:** manuscript now builds CLI `./cmd/ch01` and separate GUI `./cmd`, uses Chapter 18 artifact names and actual `--gui-dir`, and omits lone `EN_PRIMARY_SKILL` (which `ConfigFromEnv` rejects without `EN_SKILLS_DIR`, before flags are applied). Explicit `--skills-dir` retains the accepted default primary. An unused nonzero owned GUI port is required. Historical Astra/Sol and subscription examples are qualified; Bill's live-OAuth deferral remains a deferral, with no fresh entitlement or credential authority inferred. The student brief agrees and its supported-chapter list includes 18.

The initial review was HOLD. During alignment, review caught and closed the remaining singular-registry gate, hard-coded accessor wording, nonexistent `cmd/ensemble`, nonexistent GUI `--assets`, and incomplete primary-skill environment example. The root's proposed extra concrete status probe was withdrawn after inspection of the actual grader; it is not part of the accepted delta. My early cross-provider suggestion was also withdrawn: the accepted table does not price `gpt-4o`, and the final same-vendor alignment avoids that assumption.

## Source-reference erratum — root-requested correction

The first submitted version contained four factual source-reference errors, identified by the root. They were reviewer errors, not alternate API names; findings 1–3 above are corrected after a fresh source-only inspection:

- `Chapter18WithArtifacts` was wrong. The actual exported adapter is `Chapter18Targets`, verified at its definition and Edition 3 call site.
- `E3_TOOL_TRACE` / `CH4_LOG` / `CH5_LOG` were wrong. The wrapper excludes `CH07_AGENT_LOG`, `CH02_LOG`, `E3_TRACE` and `EDITION2_WEB_TOKEN`. The first three match the actual accepted CLI environment reads; `invoiceEnv` separately removes provider/settings prefixes.
- `fakevendor.Script.Models` was wrong. `ch18_driver.go:91–95,270–276` uses `fakevendor.Reply.JSONBody`; the final-text helper is `invoiceAnthropicUsage` in `ch15.go:146–154`. `fakevendor/fake.go:300–304` writes that exact body. Both packet builders explicitly encode the producer and usage. The existing `TestChapter18StatusFixturePackets` exercises this mechanism, not a Script field.
- `internal/common/pricing.go` was wrong. The actual accepted table is `solutions/edition-3/main/internal/common/cache.go`, function `LookupPrices`; accepted `internal/engine/accounting.go` calls it using returned producer identity first.

**Reassessment:** ENTRY PASS remains justified for bounded entry alignment: the actual source implements the separate-artifact route, removes the actual inherited CLI log/trace consumers, and supplies exact controlled producer/usage packets compatible with the real price lookup. The incorrect references undermine the accuracy of the original review record and are explicitly withdrawn; they do not reveal a new entry product defect. No gates were rerun, no new execution evidence is claimed, and only this note changed. Student launch remains blocked until the root reads this corrected note and JOINs. This does not expand PASS into student or operational acceptance.

## Public contract retained for whole-student review

- Status must read live owners: tool-call context → Agent → Engine → Ensemble and Engine status snapshot. Accept equivalent accessors, not placeholder state, a second counter/accounting system, reflection conveniences, or copied private ownership.
- Preserve the common-only tool boundary, all actual composition roots (CLI, GUI, tests and public programmatic use), real model selection, per-response producer attribution, session totals/cache fields/cost, legitimate unknown prices, reset/lifecycle behavior and selected-model versus producing-model distinction.
- MCP stdio and Streamable HTTP, shared registration/filtering/transport ownership, dynamic skill loading/unloading, finite web/status behavior, human CLI/GUI paths and required focused tests remain substantive obligations. Passing numerical checks cannot establish all of these.
- Prefer existing owners and helpers and the simplest adequate implementation. Reference API spelling, directory names, historical fixture models, and historical positive-test source are not extra acceptance requirements. Fixtures must not redefine the student contract.

## Evidence inspected and limits

Root-run evidence is saved under `/private/tmp/ch18-entry-checks.HgxI4R`:

- `adapter.log` / `adapter.exit`: full Edition 3 adapter suite PASS, exit 0, including Chapter 18 validation/routing and older chapter guards. The valid-target sentinel intentionally exits 71 and produces a failed eleven-check report: this proves artifact routing/report shape, **not** a passing implementation.
- `fixture.log` / `fixture.exit`: environment and explicit status-packet controls PASS, exit 0, for Sonnet, Opus and unpriced rows.
- `graph.log` / `graph.exit`: `TestChapter18ImportBoundary` PASS, exit 0; inspected the test's real allowed/forbidden/missing-registry cases.
- Same-basename `lint-baseline-v2` and `lint-final-v3`: hard checks PASS, exit 0. Warnings remain (negation count 24 → 29; person-gap 1677 → 1764); this is not warning-free prose. The final removal of a redundant launch environment assignment does not change those prose metrics. Initial wrong `cmd/prose-lint` setup failure remains separate from product defects; the root records its other intermediate edit/setup failures in the execution entry.

I inspected source and saved results; I did not independently rerun them. Root reports all 216 accepted predecessor source/support files and modes still match main; this reviewer did not perform a command-level tree comparison. No original-solution positive grade, new student grade, complete Chapter 18 end-to-end run, Chapter 1 gate, real web fetch, browser/native observation, model request or account login is claimed by this entry review. No such execution authority transfers from this PASS.

**Disposition:** ENTRY PASS. Retain this same independent reviewer for the coherent compiling student HOLD, all justified corrections, and eventual whole-student acceptance review. Root must inspect this saved result and JOIN; implementation/operational release remains root-owned.


### Root entry acceptance

I read both submitted finals, the full corrected note, actual adapter,
fixture and import-boundary changes, and saved focused receipts. I
verified the corrected symbols against actual CLI consumers and patch
bytes. ENTRY is accepted after JOIN; no student source or live behavior is
accepted. I retained this reviewer and released one student directly to
implementation and a compiling HOLD, with 0/10 correction rounds.

## Round0 whole-code review — HOLD, 0/10 corrections

This section supersedes the historical entry-only statement that no student implementation has been reviewed. Reviewed the held main implementation against current Chapter 18 and relevant accepted16 owners, including the student-listed new `http.go`, `http_test.go`, `status.go`, `model.go`, `mcp_listener.go`, and `chapter18_test.go`, rather than limiting review to tracked-file diffs. Original Chapter 18 HTTP source was a non-executed simplicity/behavior reference only. No executions, production edits, gates, accounts, browser/network operations, or new workers in this review.

### Concrete public blockers

1. **Missing auth reference incorrectly prevents an anonymous request.** `internal/mcp/http.go`, `SendContext`, reads the declared `AuthEnv` and returns an error before POST when its value is empty. Current Chapter 18 Authentication explicitly requires anonymous connection plus a nonsecret diagnostic naming the missing variable. The constructor is `newHTTPTransport`; it does not itself reject the missing value. Correct the existing request/diagnostic path; do not invent credentials, fallback secrets, provider policy, or a new authentication subsystem. The current `http_test.go` does not exercise missing auth. My interim message incorrectly attributed refusal to construction and claimed a missing-env test existed; both details are withdrawn.
2. **Useful transport diagnostics are discarded, and decoded error paths lack token redaction.** The same HTTP sender replaces a non-2xx response with only `MCP HTTP status <code>` rather than surfacing the remote reason required by the public transport-error-is-reported check. Retaining a useful bounded reason must redact the configured bearer value. Independently, HTTP JSON/SSE is passed into the shared peer without credential redaction; decoded `rpcError.Message` and tool `isError` text can reach the caller through the existing manager path. A remote error echoing the bearer value therefore has no redaction boundary. Fix diagnostic preservation and secret removal in the existing transport/error path, not by indiscriminately hiding every failure. Existing status-error tests only establish rejection, not useful safe diagnostics.

3. **GUI selected-model views revert to stale replay configuration after a live switch.** In `gui/web/gui/app.js`, the usage handler (line 168) updates `a.usage`/`a.meter`, not `a.config`. `renderMeter()` sets the header model from `usage.meter.selected_model`, but `controls()` (line 250) immediately overwrites the same header with `a.config.model`; `renderAgentTree` also derives its title from that stale config. Later state updates call `controls()` (lines 182–185). The selected model must stay truthful after state updates without requiring reload. Root independently confirmed this narrow stale-view defect. Correct the existing GUI view/state path, not the accounting owner. The public settings paragraph does not establish a new global model-storage requirement; no such persistence requirement is imposed here.

**Disposition: HOLD, not PASS FOR VALIDATION.** The three bounded correction areas are missing-auth anonymous behavior, useful redacted MCP diagnostics, and truthful live GUI model views. They are sufficient to withhold broader gates. Root owns the correction assignment and round count; this remains the initial 0/10 review.

### Ownership, simplicity, and lifecycle assessment

- **Status/accounting:** the actual tool follows `ToolCall.Parent()` → Agent → `Engine().Status()` and reaches the existing Ensemble owner; no second token ledger or fake status owner was added. Existing per-producer accounting and `LookupPrices` remain the source of costs. Selected model and last-producing model are separate; model switches do not reprice old work. The known subtotal with explicit partial/unknown status is honest, including per-model buckets; it is not a claim that unpriced work costs zero. GUI/CLI consume that snapshot. Mandatory source review remains necessary beyond numeric grader checks.
- **Model controls:** CLI `/model` and JSON model commands, plus the GUI model command, reach the Agent's urgent control path. The idle-boundary mutation changes owned model fields, validates same-provider selection, preserves history/accounting, and leaves provider credentials/configuration intact. The browser control is session/Agent-local, not a rewrite of global persisted GUI settings; the chapter does not justify requiring a new global model preference schema. Engine request preparation includes current model identity. The stale header observation above is distinct from successful model selection.
- **Transport:** HTTP is selected below the existing peer/manager registration and filtering path. Source handles JSON and SSE replies, response IDs, session/protocol headers, redirect refusal, context-bound HTTP work and owned cancellation/close. No client-wide timeout was added. Existing stdio still uses the accepted stream/process ownership: child pipes and process group close, and a joined process waiter. The peer retains cancellation-aware writer admission and the existing non-context transport close-to-unblock behavior. These source properties are not a race/stress or protocol-conformance execution result.
- **Skills/credentials:** primary exposes status; web-search is a separately loadable URL MCP skill using its own auth environment reference. Shared parsing/manager ownership and existing unload/disconnect lifecycle remain in use. `ProcessEnvironment(agents...)` removes configured auth-reference names from local command/MCP child environments as well as the established provider names; it is not limited to the shipped Firecrawl variable. This does not replace the missing response redaction identified above.
- **Compositions/observer lifetime:** CLI, GUI and public API reach the existing owners. Dedicated `--mcp-pipe`/`--mcp-socket` branches remain separate from normal UI startup. The new `StartMCPListener` handle/`--mcp-listen` is not needed merely to implement the outgoing URL client or status tool; assessed as the requested sibling-server composition, it factors the existing server rather than adding a second Agent or protocol owner, and callers defer its owned close/join. No new architecture or mandatory removal is prescribed. Identified observer queues retain their existing close/unsubscribe ownership (CLI defers queue close; GUI closes observation subscriptions through its server lifecycle); no parallel status subscription/counter service was added. These conclusions do not certify unexecuted lifetime behavior across every composition.

### Evidence qualifications and next boundary

The student's durable round0 note reports successful focused tests and six compilation-only module checks under `/private/tmp/e3-ch18-restart-student/compile`. It also preserves earlier real compile/test failures and fixes. I inspected the added test source, not a fresh run; the remote HTTP test's non-2xx assertion does not prove diagnostic preservation/redaction, and the public tests do not resolve the blockers above. No independent execution or full receipt audit is claimed here. Root has the saved execution evidence and is reviewing in parallel. A single resumed infrastructure Responses failure did not authorize repeating tests or source work.

After the minimal corrections, retain this reviewer for the corrected coherent HOLD and source/delta reassessment. No numerical grade, full/race/vet gate, browser/native observation, real provider/web probe, live OAuth, or operational acceptance follows from compilation or this review. Root must read this saved report and JOIN before issuing the correction assignment or lifting HOLD.

### Root round0 decision and R1 release

I read the held changed/new files, actual CLI/GUI/MCP composition and observer
paths, original compact transport/ownership reference, and the complete retained
review. The original was a simplicity reference, not a test oracle. I agree with
the three public blockers and retain the existing owner split. The retained
reviewer final was saved/read and JOINed before R1 release.

I independently read student receipts io/1098,1102,1104,1105,1106,1111,1117,1118
under the saved student context. The successful focused outputs and six-module
compile table are evidence, not full gates. Missing post priority and net import
were setup compile failures; the asynchronous History assertion mismatch and
its corrected receipt remain preserved. The student's later Responses inference
interruption happened after its note/checks: same context resumed to inspect
saved state and submit, without replaying checks. The reviewer assessed source;
its note appropriately does not claim independent execution.

Released R1/10 to the same student: logged anonymous fallback when a declared
credential reference is unset; useful bounded HTTP/RPC/tool diagnostics with
bearer redaction; and non-stale live GUI model views. Required focused actual
HTTP/control-path regression coverage and affected root/GUI compilation only.
No broad gates, grading, operations, new architecture, or private fixture
requirements were released. Student remains subject to existing public-source,
offline, owned-temp and disk-floor boundaries. Same reviewer retained.

## R1 delta review — HOLD (1/10), one remaining blocker

Scope: only the three round0 blockers, their changed production paths (`internal/mcp/http.go`, `transport.go`, `gui/web/gui/app.js`), added/updated focused tests, and the named saved receipts. Prior whole-code ownership/simplicity conclusions and ENTRY erratum carry forward. No commands, tests, grades, broad gates, network/browser/provider operations, or production edits were performed by this reviewer.

**Closed: anonymous attempt and diagnostic.** `dial` now passes the existing diagnostic writer into `newHTTPTransport`; the manager supplies the owning Ensemble logger's writer. An absent declared credential logs only its environment-variable name and the anonymous-attempt message, leaves the bearer empty, and no longer prevents POST. `TestHTTPChapter18R1Anonymous` exercises that real dial/HTTP path, the absent Authorization header, and the named diagnostic.

**Closed: truthful live GUI model views.** `selectedModel(a)` uses the live meter identity before replay-config fallback. Header, model editor and tree title share it, including `controls()` after state updates; the last-producing identity remains separate in the meter tooltip. The new Go GUI test sends actual WS `set_model`, `agent_status`, and `pause` commands; its captured frames enter `socket.onmessage` in production `app.js` loaded by the existing controlled Node fixture. Assertions cover the labels/editor/tree after control/state refresh and retain the old producing-model tooltip. This is meaningful controlled consumer evidence, not a real-browser observation. The CLI test likewise exercises actual human `/model` and `/status` paths and producer/cost retention without demanding global model persistence.

**Remaining public blocker: bounded HTTP diagnostics can retain a Unicode-escaped bearer.** In `exchange`, a non-2xx body is read only through `4096 + 6*len(token)` bytes. If that prefix is incomplete JSON, `redactHTTPError` falls back to replacing only the literal token and its canonical `json.Marshal(token)` spelling. It does not cover other valid JSON spellings. Concrete source-derived case: token `secret`, HTTP 403 with a complete response body beginning `{"error":"quota denied \u0073\u0065\u0063\u0072\u0065\u0074` (single JSON backslashes on the wire), followed by more than 5000 filler characters and the closing string/object. The bounded read loses the closing JSON; decode fails; neither fallback replacement matches the Unicode escapes. The reported reason therefore retains the recoverable escaped credential. A short complete JSON body is correctly decoded and redacted. Root independently confirmed this analysis. **This case was not executed**, and the passing long-JSON test only uses canonical `json.Marshal` escaping, so it does not settle this boundary.

Keep the correction narrow: safely handle or conservatively omit unsafe undecodable fragments while retaining safe status/useful reason where available. Do not require faithful display of arbitrary malformed bodies, a new authentication/response framework, or changes to successful payloads or correlation. The existing decoded JSON/SSE error path is otherwise well-scoped: only error/result-isError fields are rewritten, top-level IDs remain raw, `UseNumber` preserves numeric values, and native content structure remains intact. No new ownership or stdio lifecycle regression was identified in this delta.

**Evidence inspected:** all named `r1-*.log` receipts under `/private/tmp/e3-ch18-restart-student/` and `r1-compile/results.txt`. The before logs demonstrate the original HTTP failures and stale production-JS view. `r1-root-after.log` has passing root/CLI packages but an MCP test-destination `json: Unmarshal(nil)` setup failure; subsequent `r1-mcp-after.log` and `r1-mcp-final.log` pass. `r1-gui-after.log` passes; compilation table records root and GUI exit 0. These are saved student executions, not fresh reviewer executions; compilation is not broad validation. Unchanged round0 conclusions/results are carried rather than rerun.

**Verdict: HOLD, not PASS FOR VALIDATION.** Only the same public redaction requirement remains open. Root owns the R2 assignment after reading/JOIN; anonymous/GUI work and unrelated architecture need not be reopened. All broader gates and operational acceptance remain deferred.

### Root R1 decision / R2 release
I read the complete R1 source delta, HTTP/CLI/actual-WS-to-production-JS tests,
saved before/after logs and affected compile table. Anonymous fallback and live
model-view corrections are closed at source/focused-evidence level. I agree
with the retained reviewer’s concrete incomplete-JSON redaction gap; this was
source-derived, not a newly executed negative control. Reviewer final saved,
read and JOINed. Released only that bounded diagnostic repair as R2/10 to the
same student, with an actual HTTP regression and focused HTTP checks plus
affected compilation. Safe conservative omission of undecodable fragments is
allowed; faithful arbitrary malformed-body display/new parsing machinery is
not required. Unchanged GUI/CLI and ownership evidence carries. No broader
gates or operations released.

## R2 delta review — PASS FOR VALIDATION (2/10)

Reviewed only the changed `redactHTTPError` fallback in `internal/mcp/http.go`, its existing non-2xx caller, new `TestHTTPChapter18R2TruncatedUnicodeBearer`, student R2 note, and named saved receipts. No executions, new exploration, production edits, or reopening of closed GUI/CLI/ownership work.

The remaining R1 blocker is closed. After removing literal and canonical JSON token spellings, the undecodable-body fallback now omits a fragment if any backslash escape syntax remains. The caller retains the HTTP code and standard status reason (`403: Forbidden` for the reviewed case). This deliberately conservative behavior is acceptable at the bounded/undecodable boundary: it avoids exposing recoverable alternate escaped credentials without promising faithful display of arbitrary malformed fragments. Safe useful diagnostics and fully decoded JSON redaction remain on their existing paths. The small fallback guard adds no authentication/response framework and does not change successful payloads, correlation IDs/native values, HTTP read limits, stdio, or GUI/CLI behavior.

The actual-HTTP regression sends complete valid 403 JSON with the Unicode-escaped token followed by 6,000 filler bytes, causing the production bounded read to truncate it. It verifies the request's bearer, absence of the recoverable token, bounded output, and retained safe status/reason. Inspected `/private/tmp/e3-ch18-restart-student/r2-http-before.log`: real pre-fix leak failure. Inspected `r2-http-after.log`: passing saved `TestHTTPChapter18` run, including prior cases as specified by the recorded command. Inspected `r2-compile/results.txt`: root and GUI exit 0, compilation only. These are student-run receipts, not reviewer executions or broad gates.

**Verdict: PASS FOR VALIDATION.** All three initial blockers are now closed across R1/R2; prior ENTRY corrections and whole-code conclusions carry forward. Root must read/JOIN and owns any next validation release. This verdict is not numerical grading, full/race/vet validation, real-browser/native evidence, live/provider approval, operational acceptance, or final chapter acceptance. Reviewer remains retained and HOLD pending authorized validation findings.

### Root R2 verification and local validation release
I inspected the minimal fallback guard, actual-loopback regression, complete R2
student note, real before-failure/after-pass receipts and root/GUI compile-only
table. I agree with conservative omission at this undecodable boundary rather
than new parsing machinery. The retained review final was saved, fully read and
JOINed. R2/10 is PASS FOR VALIDATION, not chapter acceptance.

Released sequential offline local gates in
`/private/tmp/ch18-validation.l8oKlb`: fresh held CLI/GUI artifacts, root/GUI
race and vet, then auxiliary-module tests/vet. Optional inherited command test
receives actual `CH16_TEST_CLI`/`CH16_TEST_GUI` variables. Canonical owned
HOME/TMP, clean dummy environment, existing caches and disk headroom apply.
Each step saves its own log/exit; first failure stops later gates. Grading and
live operations have not run. Primary configuration and live GUI are untouched.

## First root validation stop — fixture classification, source-only

Inspected `/private/tmp/ch18-validation.l8oKlb/root-race.log` and the exact affected tests/construction paths. Root-race exited 1; preserve this failed receipt. The log reports declaration assertions and a nil-owner panic, not a race-detector report. It is still a failed gate, not evidence that the root race suite passed. Builds passed as reported by root; subsequent vet/GUI/auxiliary/grade stages were not run. No reviewer executions or edits to production/tests.

- **Chapter 3 declaration expectations are stale test data.** `chapter03_loop_test.go:244` expects 13 default declarations; the logged actual set contains the same 13 plus required `agent_status`. `chapter03_schema_test.go:73` similarly hardcodes 13 for both common and Gemini declarations. Minimal repair: update the intended exact default set/count to 14, including `agent_status`, while keeping ordered effects, call/result identity, schema projection, strictness and no-mutation assertions. Importantly, `c3CheckDeclarations` gates its named-default checks on `want == 13`; merely changing the call count would silently bypass them. Update that default-set condition and named list coherently rather than weakening the assertion or deriving all expectations from the implementation under test.
- **Chapter 15 panic is an invalid private fixture, not a supported ownerless Engine.** Current `internal/engine/engine.go:33–37` explicitly rejects a nil Agent parent; accepted16 already had that requirement. The public construction in `ensemble.go:137–147` creates the Agent, calls `engine.New(a)`, and wires its cache/owner. The test instead directly constructs `&engine{last: common.UnknownUsage()}` and bypasses the invariant. `Status()` now legitimately obtains the selected model from `e.parent.Config()` (`accounting.go:81`). There is no source basis for adding production nil-parent fallback solely to support this fixture.
- **Minimal accounting fixture repair:** use the existing `parseAgent` config-bearing test owner (`sse_test.go:14–21`) with the existing constructor, or provide that owner in the narrowly scoped accounting fixture. Preserve all captured-producing-identity, price, partial/unknown cost, bucket-copy isolation and weighted-cache assertions. The immediately adjacent `TestChapter15WeightedRatesAndUnknownPrices` has the same invalid construction for both `e` and `unknown`; it was masked by the earlier panic and should receive the same setup-only correction. This does not authorize changing expected accounting results or broad fixture refactoring.

**Recommendation:** narrow test-only repair, no product patch established by these failures. Prior source PASS FOR VALIDATION is unchanged, but validation remains STOPPED until root authorizes the repaired gate attempt. No fresh gate, numerical acceptance, or claim about unrun stages follows from this classification.

### Root test-only repair verification
I inspected the exact three-test-file diff and focused race receipt. The named
14-tool set/guard and provider/schema/no-mutation checks remain explicit;
accounting tests use the existing valid owner/constructor without changing
cost, unknown, bucket-copy or weighted-rate expectations. No production file
changed. Same student saved HOLD/JOIN; production remains R2/10. Resumed the
root gate as `root-race-v2` and remaining original sequence, reusing unchanged
held CLI/GUI artifacts and retaining failed `root-race` intact. No grade yet.

## Standalone MCP DOM script stop — fixture boundary classification

Inspected `/private/tmp/ch18-validation.l8oKlb/js-mcp_dom_test.log`, `gui/mcp_dom_test.cjs`, and the relevant `app.js` helper/bindings. All three scenarios fail during `fixture()` evaluation with `ReferenceError: $ is not defined`, before their behavior assertions. Preserve this failed receipt; it is not a passing DOM check.

The exact current fixture evaluates `source.slice(source.indexOf('// Historical browser MCP tools:'))` through EOF. That tail now includes the model/status button registrations (`app.js:442–443`), but its VM context does not supply `$`. The real page defines `$ = id => document.getElementById(id)` at `app.js:2`, outside the excerpt. Thus this particular failure is the controlled fixture's incomplete evaluation boundary/environment, not an established production defect.

Minimal fixture-only correction: either bound the excerpt to the actual browser-MCP implementation, or supply the minimal page-control lookup/listener stubs needed by the unrelated registrations. Continue executing the actual production browser-MCP functions; preserve all snapshot/selector/state assertions, UTF-8/Markdown elision accounting assertions, and JSON-RPC/reconnect readiness assertions unchanged. Do not remove those scenarios, catch-and-ignore evaluation errors, or rewrite a substitute peer implementation.

Prior source PASS FOR VALIDATION remains unchanged. Root reports root/GUI race/vet and auxiliary gates passing after the verified test-only repair, and prior `client_round1`, Chapter12 and Chapter08 scripts passing; this note does not rerun them. The earlier nonexistent `gui_cards` filename is a separate preserved command-setup failure, not a product defect. Remaining scripts and grading are still stopped. No reviewer executions, code edits, or broader exploration.

## First Chapter 18 grade stop — two distinct causes

Inspected `/private/tmp/ch18-validation.l8oKlb/grade18-v1.log`, the relevant `ch18.go`/`ch18_driver.go` checks, actual dispatch and GUI/Agent model-control paths, and current public chapter. Preserve **FAIL 105/155, exit 1, parsed false**. No Chapter 1 grade, new grade, test, or code edit was performed by this reviewer. Root reports the preceding root/GUI race/vet, auxiliary and repaired standalone-JS gates passing; that does not make this grade pass.

1. **Status scenario: grader adapter mismatch, not evidence of a model-selection/accounting defect.** `ch18_driver.go:187–190` sends `update_settings` with `{model: ...}` and waits for `settings_changed`. Actual `gui/settings.go` deliberately has no model preference: `patchSettings` rejects it as an unknown field; `server.go:191–195` sends `settings_error`. Therefore the driver's `map[]`/deadline diagnostic is explained by waiting for an event that this command cannot produce. The real session operation is `set_model` with `text` (and the target `agent_id`): `server.go:207–212` calls `Agent.SetModel`, then `Hub.refreshStatus` (`hub.go:464–469`) publishes a `usage` snapshot from the owning Engine. `internal/agent/model.go` serializes that change through the actor. Current public chapter lines 62–65 and 989–994 explicitly separate live session selection from saved preferences. Minimal adapter repair: use the actual session control and confirm its target/model usage event; surface relevant errors instead of waiting only for a settings success. Do not add a persisted model field or rewrite student preferences semantics to accommodate this stale driver. Preserve all later checks of actual provider model, prior history, counts, producing-price attribution and unknown cost. Those checks did not finish in this failed run.

2. **Unknown/disabled refusal: genuine narrow public wording mismatch, not failed grant enforcement.** `internal/tool/dispatch.go:175–176` correctly refuses an absent selected definition as `unknown or disabled tool %q`; the gate's declaration/dispatch checks do not report actual execution of the hidden/unloaded tools. However the current chapter is explicit: lines 445–449 require an unknown-tool refusal rather than “disabled,” with a stated information-disclosure rationale. `webUnknown` (`ch18.go:235–239`) accepts several equivalent unknown categories but rejects “disabled”; it therefore matches that published distinction. This cannot be classified honestly as a pure grader false negative merely because the inherited generic message also contains “unknown.” The smallest conforming production correction is unknown-only wording for this nonexistent/ungranted-definition path, preserving the real failure result and zero execution. If root instead chooses to accept the inherited generic category, that requires explicit public-contract alignment, not silently weakening the grader. No new policy machinery or exact punctuation requirement is warranted.

**Disposition:** validation remains STOPPED. One adapter correction and one narrowly scoped production/public-conformance issue are established; neither justifies changing ownership, accounting, live preferences semantics, or grant behavior. The prior source PASS FOR VALIDATION is historical, not final acceptance and not an override of this newly identified wording issue. Root owns the correction decision and sequential release; reviewer retained/HOLD.

### Root grade-stop correction and R3 verification
I initially overclassified both grade failures as fixture assumptions. Rereading
public lines445–449 confirmed unknown-only wording is explicit; I withdrew that
claim and kept `webUnknown` unchanged. Real grants refused calls; this was not
successful unauthorized execution. The model timeout is a separate root driver
mismatch, now repaired only on the Edition3 route with actual live `set_model`
and target-correlated `usage` acknowledgment. Historical settings callers,
3-second deadline, all11 checks/155points and substantive packet/history/cost
assertions remain unchanged.

R3 student HOLD/JOIN: one production refusal-line change, actual absent versus
registered-ungranted regression, unchanged grant/dispatch semantics. I read the
actual regression and before/after/compile receipts; first guessed test filename
was a lookup failure only. Root loopback model-control fixture passed in6.018s,
exit0, including wrongAgent/wrongModel/refusal and legacy branch controls. These
are separate from a passing chapter grade. Same reviewer reviewing both deltas.

Carried local evidence: root-race-v2/root-vet, GUIrace/vet, auxiliary tests/vet;
ch06 and consumer no-test modules are compilation-only. Standalone clientRound1
chapter12, ch08, repaired MCPDOM, ch13 and ch15 meter checks passed. Original
missing `gui_cards_test.cjs` invocation was a root setup failure; original DOM
ReferenceError and failed root gate remain preserved. Source-backed fixture
repairs do not consume production rounds. No installation or live operation.

## R3 stable delta — source PASS FOR VALIDATION; root model-control fixture accepted

**Student source: PASS FOR VALIDATION.** Read the actual shared unavailable-entry branch in `internal/tool/dispatch.go`: it now returns `unknown tool %q`; grant selection and dispatch are unchanged. `TestChapter18UnavailableToolsAreUnknown` calls the same name first absent and then registered-but-ungranted through the real Agent registry. It retains exact response equality, error status for both, unknown-only wording, no declaration, and atomic zero-handler-execution assertions. The two inherited `chapter03_loop_test.go` wording changes preserve their order/ID, provider, filesystem-effect and error-result assertions. No additional production correction is indicated by this delta.

Inspected saved evidence: `r3-unknown-before.log` reproduces the old availability-disclosing response; `r3-unknown-after.log` passes the focused race run (`ok ensemble 1.720s`), whose command and included regressions are recorded in the student R3 note. `r3-compile/results.txt` records root and GUI exit 0, compile-only. These are saved student-run results, not reviewer executions or a fresh full-suite claim.

**Root fixture: accepted for the next controlled grade.** Read propagation through `Chapter18Targets` → `chapter18Targets` → `statusChecks` → `statusLaunchModels`: live control is true only for the explicit restarted-artifact entry; historical `Chapter18` and `statusLaunch` retain false. `selectStatusModel` keeps the legacy settings branch and uses the prior actual usage snapshot's Agent ID for live `set_model` (`text`/`agent_id`). It requires matching Agent and `meter.selected_model`, surfaces matching/unscoped errors, and retains the three-second deadline. It does not synthesize usage or add persisted preferences.

Read the complete `TestChapter18ModelControl` in `ch18_contract_test.go`: a real local WS peer exercises matching usage after unrelated/stale frames, wrong-Agent-only and old-model-only deadline failures, refusal, and legacy success, with exact outbound-command assertions. Inspected `/private/tmp/ch18-validation.l8oKlb/model-control-fixture.log` and `.exit`: package PASS in 6.018s, exit 0. No new fixture framework is necessary.

The unknown-category predicate, grant checks, 11-check/155-point rubric, and substantive later provider-model, conversation-history, prior-count/cost preservation, status/count/rate and known/unknown accounting assertions remain in force. Matching a usage acknowledgment does not replace those later checks. The grade18-v1 FAIL 105/155 is preserved and is not retroactively changed to PASS.

**Release boundary:** source and fixture review permit root's stated next step—rebuild the CLI/GUI and perform one controlled Chapter 18 grade. Earlier root/GUI race/vet, auxiliary and standalone-JS evidence is carried for unchanged behavior; this review does not propose or claim a full-gate replay. No executions, production/grader edits, Chapter 1 grade, operational authorization, or final chapter acceptance by this reviewer. Retained reviewer submits/HOLD for root read/JOIN and subsequent evidence.

## Chapter 18 grade v2 stop — status representation mismatch

Inspected `/private/tmp/ch18-validation.l8oKlb/grade18-v2.log` and the actual handler, Engine snapshot, correlated-result extractor and status checks. Preserve **FAIL 120/155**, with no Chapter 1 grade. The failures report missing model/count/hit-rate prose and complete numeric cost; they do not demonstrate that the corresponding structured values are absent or wrong. The transport/grants/model-control corrections remain carried; no new reviewer executions or code edits.

**Concrete source path (corrected):** `internal/tool/status.go` defines `statusTool`; its handler returns `json.NewEncoder(call).Encode(call.Parent().Engine().Status())`, writing the snapshot directly to the `common.ToolCall`. `internal/common/cache.go:72–82` defines `EngineStatus` with `selected_model`, `last`, `session`, `session_hit_rate`, and `cost` (`usd` and `status`), among other fields. `internal/engine/accounting.go:71–87` builds that disposable view from the actual Engine buckets/counts and owning Agent's selected model. This is not an absent handler or a separate status meter.

The actual grader path read here is `skillResult` (`ch09.go:177`) → `anthropicToolText` (`ch03.go:118`) → `statusReport`/`statusCost` (`ch18.go`). `statusToolText` and `statusToolReplyText` were not present in the current grader source searched; do not cite them as the implemented path. Extraction preserves the correlated tool-result text/error flag, but the subsequent checker recognizes only historical prose lines (`Model:`, `Last:`, `Session total:`, a percent line and `cost: $...`). A JSON-encoded snapshot does not match these regexes, explaining the missing-field diagnostics without implying missing accounting.

**Ownership: narrow root/grader adapter correction.** Current restarted TL;DR and ownership contract require the semantic fields and real owner chain, not that historical line serialization. The later historical discussion of five labels is not a mandate to replace the accepted structured snapshot. Support that actual structured result alongside the retained historical prose form. Compare selected model, both four-counter groups, session hit-rate fraction against the independently calculated rate, and `cost.usd`/`cost.status` against independently priced or unknown/partial expectations. Required fields/types must not silently default to zero; malformed/missing values must fail. In particular, cost completeness must come from the cost status, not an unrelated occurrence of “unknown” elsewhere in the JSON. Keep actual tool-call correlation, failed-result detection, provider model/history/prior-total preservation and the complete rubric unchanged. No new production formatter, counter, prices or aliases are justified by this failure.

**Limit:** source establishes the representation mismatch, not that all captured numeric values will pass a corrected comparison. Validation remains stopped until root's bounded adapter/control review and sequential release; grade v2 stays failed evidence. Prior source PASS FOR VALIDATION remains, with no newly established student product defect. Reviewer retained/HOLD.

## Structured-status reader — accepted for controlled validation

**Verdict: root fixture delta ACCEPTED.** Read `statusJSON` and both new structured branches in `ch18.go`, together with the full `TestChapter18StructuredStatus` in `ch18_contract_test.go`. The reader directly interprets the correlated snapshot; it does not manufacture prose or derive expected answers from student code. Selected model must match; every required last/session counter must be a JSON number equal to its expected disjoint count, without a tolerance. Hit rate must be finite and in [0,1], uses input+cache-read+cache-write as denominator, and retains the historical 0.15-percentage-point tolerance. Cost status explicitly distinguishes known from partial/unknown; known and partial amounts (and any supplied unknown subtotal) must be finite, nonnegative and agree with the independently computed price within the existing tolerance. An unknown cost may honestly omit a number. Unrelated diagnostic strings no longer decide pricing completeness. Malformed/overflow JSON fails rather than falling back to prose. The historical regex path is retained.

The focused controls assert known/partial/unknown success, unrelated-unknown-string independence, missing/wrong model, missing/fractional/string counters, absent rate/wrong denominator, missing cost, null/string/wrong amount, contradictory known/partial status, malformed/overflow/empty objects, and the historical prose representation. Inspected `/private/tmp/ch18-validation.l8oKlb/structured-status-reader.log` and `.exit`: package PASS, 0.009s, exit 0. These are saved root-run results; no reviewer execution occurred.

The surrounding 11-check/155-point flow still obtains the actual correlated non-error result and checks actual provider models, declared status capability, conversation history, prior totals/cost preservation and the independent measured values. No student round, production change, binary rebuild or gate replay is caused by this reader correction. Failed grades v1 and v2 remain preserved. Root's stated one next controlled grade using unchanged R3 binaries is consistent with this verdict after read/JOIN; this is not an advance passing-grade or chapter-acceptance claim.

**Explicit source-reference erratum:** the preceding v2 note originally described a nonexistent `c.Buffer()` write. Re-reading the actual 20-line `internal/tool/status.go` confirms `json.NewEncoder(call).Encode(call.Parent().Engine().Status())`; the paragraph above has been corrected. The returned type is `common.EngineStatus`, not a hypothetical `UsageMeter`. This was my inaccurate source description; it does not change the verified structured-result mismatch or this reader verdict. No other production API is being asserted by that shorthand. Reviewer retained/HOLD.

## Final LOCAL-code disposition — ACCEPTED; live verification pending

**LOCAL CODE ACCEPTED on the held R3 source and artifacts.** This combines the completed whole-code/public-contract/ownership/KISS review, comparison with accepted Chapter 16 and the original reference, the reviewed three production correction rounds, accepted root-owned fixture corrections, and the carried local validation evidence. No outstanding local-code blocker is established by the reviewed record. This is not overall Chapter 18 closure.

Read the actual final receipts under `/private/tmp/ch18-validation.l8oKlb/`:

- `grade18-v3.log` and `.exit`: **155/155, passed true, exit 0**, all 11 checks passing. This includes actual URL transport/grants/dispatch/errors, status counts/hit rate, turn-boundary model selection and mixed-model accounting. The source-reviewed owner chain remains necessary beyond the resolved-import check; its passing score does not replace that review.
- `grade1-v1.log`, `.exit`, and `.elapsed`: **100/100, passed true, exit 0; 124.541 seconds elapsed**. The report records five well-formed requests/rounds, append-only history, exact cumulative usage **511 input / 65 output**, and complete answer/usage after separate genuine **1m2s delayed-header and 1m2s delayed-body** cases. Those finite waits are observed evidence, not a claim about every possible duration.

Root/GUI race/vet, auxiliary and standalone-JS results are carried with their documented fixture-only repairs, together with the focused R3 assertion/race and compile evidence. No unchanged tests were replayed by this reviewer. Earlier grade v1 **105/155** and v2 **120/155** failures, the earlier test/setup failures, and their explicit classifications remain preserved; they are not retroactively rewritten as passes. Root fixture alignment is separate from the three student production correction rounds.

**LIVE NOT YET RUN in this final local-validation phase.** Bill now explicitly authorizes live tests/API-key use; root will conduct the finite live exercise separately. Local acceptance does not supply live results, model availability, successful authentication, OAuth/subscription entitlement, installation, runtime rollout, or overall chapter closure, and the historical OAuth exercise is not borrowed as evidence. No reviewer executions, source edits or operational actions occurred in this disposition. Same reviewer/context retained for the live evidence and closure; submit/HOLD for root read/JOIN.

### Root failed live setup and explicit corrected attempt

Bill explicitly authorized live tests/API keys. `/private/tmp/ch18-live.JtNrDD` failed before any MCP call: root reconstructed an older primary and omitted inherited `loadable-skills`, causing the correct offered-skill refusal. Root also used older recalled prompts rather than the current reproduction. The current book already declares loadable-skills and specifies the stable Go release/archive task, status/unload, then real model switch/status; root withdrew its initial attribution to the book. No manuscript or production fix was made. Four Sonnet4.6 POSTs and one discovery GET, zero MCP POSTs; no web evidence, model switch, later turns, or normal EOF. Actual usage 6 fresh input/4763 writes/3918 reads/844 output; independent $0.03171465 equals process meter. SIGTERM was owned cleanup. An exact-anchor edit failed after writing a corrected controller but before any checks/keyread/live launch; root retained that directory and finished only missing setup. The separate `/private/tmp/ch18-live-corrected.lwY8Lj` attempt uses the current primary and exact three public prompts, changing only the explicitly allowed supported API-key model pair. Same pinned R3 CLI, finite limits, no automatic retry. Root explicitly released this corrected attempt under Bill’s current authorization; preserve the first failure.

## Live exercise stop — inherited Anthropic marker-legality defect

**Concrete production blocker; HOLD pending bounded correction.** Read the corrected attempt's `receipt.json`, `response-2.json`, MCP status records and `primary.md` under `/private/tmp/ch18-live-corrected.lwY8Lj`, and the relevant renderer in current main and accepted Chapter 16. The first Sonnet 4.6 provider request returned HTTP 200 and loaded the web skill. Actual anonymous MCP initialize/list returned HTTP 200 SSE, with initialized notification HTTP 202. The second provider request returned HTTP 400: `Tool 'firecrawl_search' cannot have both defer_loading=true and cache_control set. Tools with defer_loading cannot use prompt caching.` This is not evidence of an MCP initialization/session/ownership failure.

**Actual source and inheritance:** `internal/engine/tools.go`'s `anthropicCatalog` sets `DeferLoading: !t.Enabled || disclosed[t.Name]`. After building that catalog, `internal/engine/render.go:276–287` unconditionally assigns `tools[len(tools)-1].CacheControl = ephemeralCache()` whenever caching is enabled and tools are nonempty. Accepted `solutions/edition-3/ch16` contains the same relevant code. `breakpoints.go` owns the cache-control type/helper, but the defective final-tool assignment is in `render.go`; no `applyAnthropicCache` helper was found in the current inspected engine source. The inherited renderer can therefore mark a deferred last tool, exactly the combination rejected during this Chapter 18 flow.

Root separately verified the provider's primary documentation, “Prompt caching” in <https://platform.claude.com/docs/en/agents-and-tools/tool-use/tool-search-tool>: deferred tools cannot also carry `cache_control`; put the breakpoint on a non-deferred tool. I did not fetch that document or make another paid request. Together with the saved provider rejection and actual source, this is a public provider-legality correction, not a private fixture requirement or a new model/price policy.

**Minimal ownership/scope:** correct the Engine's tool-marker eligibility to use the last eligible non-deferred tool, or no tool marker when none is eligible. Preserve the full catalog, defer flags, paired tool references, content and other legal cache markers; do not remove progressive disclosure or change grants to make the request pass. The accepted MCP transport and real accounting/ownership conclusions remain carried. Prior local passes remain historical evidence, but this newly observed production blocker prevents overall acceptance until corrected and appropriately validated. R4 is not released by this note.

**Limits/correction of attribution:** two provider POSTs were recorded, with known usage only for the first; search/scrape counts are zero, no model switch occurred, and cleanup ended by SIGTERM rather than normal EOF. This is not a completed live exercise or a verified final cost. Both the current book and corrected saved primary include `loadable-skills: web-search`; the earlier root-reconstructed primary omitted it, so blame of the book for that earlier attempt is withdrawn. No reviewer executions, tests, paid retry, implementation changes or additional exploration. Same reviewer retained/HOLD for root read/JOIN and bounded next delta.

### Root R4 release and continuation

After reading and JOINing the live-blocker classification, root released production R4/10 to the same student: choose a legal non-deferred tool breakpoint or omit that slot; preserve catalogue, defer flags, paired references, grants, native content, other markers and History. Focused actual-request regressions plus affected root/GUI compile only, no external tests or broad gates. The corrected failed attempt has $0.00578025 of independently computed known first-response cost; rejected response has no usage, so total billed usage is not established. Both attempt directories remain preserved. Later student model transport ended with connection reset; root resumed the same context from saved work after Bill reported the laptop reopened, without claiming that report proves the request-failure cause. Bill renewed autonomous continuation. No tunnel settings/processes were changed; local8088 was observedHTTP200 during read-only diagnosis.

## R4 delta review — PASS FOR VALIDATION

**Source verdict: PASS FOR VALIDATION, limited to the live-discovered marker defect.** Read the actual `internal/engine/render.go` change, the full affected `TestChapter11AutomaticAnthropicProjection`, `TestChapter18ToolCacheBreakpointEligibility`, the existing declaration/marker helpers, the student's R4 record, and the named saved receipts. The backward scan uses the actual rendered Anthropic declarations and places at most one tool breakpoint on the last declaration whose `DeferLoading` is false; no eligible declaration means no tool marker. This directly satisfies the provider rule and prevents the rejected deferred-tool/`cache_control` combination. It does not reorder/remove tools, change schemas/defer flags, disable caching globally, alter grants/references, or change the existing system/history markers. The short inline scan and explanatory comment are proportionate; no new cache framework or ownership layer is introduced.

The real-Agent/local-HTTP control observes the request after an actual `load_skill`, rejects any deferred declaration carrying a marker, checks the last eligible position and marker budget, and compares the ordered catalog/schema/defer flags before and after disclosure. Existing startup/dynamic provenance, real automatic-handler counts, reference-only paired-result projection, ordinary unsupported fallback, canonical load call/result, replay and history/context non-mutation assertions remain. The renderer boundary cases cover a deferred tail, eligible final entry, all-deferred and empty catalogs, including exact remaining system/history marker budget and catalog immutability.

**Saved evidence inspected:** `/private/tmp/e3-ch18-restart-student/r4-cache-before.log` genuinely fails the actual post-load HTTP request and the trailing/all-deferred renderer cases. Its separate unsupported-model failures were an erroneous new test expectation: current final test correctly requires zero markers for that unchanged unsupported fallback, rather than weakening the supported-model checks. `r4-cache-after.log` records both root and engine packages passing (0.395s/0.007s) for the focused commands listed in the student note. `r4-compile/results.txt` records root and GUI compile-only exit 0; no executed-suite claim follows from compilation. The preserved infrastructure interruption is not relabeled a product/test failure or grounds to repeat completed work.

All previously reviewed, unchanged Chapter 18 transport/accounting/model-control/ownership/KISS areas and evidence remain carried. This closes the source-level correction request, not the live exercise: the earlier provider rejection and other failed receipts remain preserved. Only root may now release the appropriate affected local validation and a separately explicit bounded live attempt; no broad gate, grade, network/credential action, production edit or reviewer execution occurred here. No advance claim of post-R4 local acceptance or successful live search/scrape/model switch/normal EOF. Same reviewer retained; submit/HOLD for root read/JOIN.

### Root R4 local validation and bounded live release

Root read the complete held R4 code/tests and saved receipts, agreed with retained PASS, and JOINed. Fresh root race/vet passed, with correctly supplied fresh command-test artifacts; unchanged GUI/auxiliary/JS evidence carried. `/private/tmp/ch18-validation.l8oKlb/grade18-v4.json` passed155/155/true/exit0; `grade1-v2.json` passed100/100/true/exit0 on fresh same-source CLI, 124.49s including both62-second phases, usage511/65. Artifacts in `bin-r4`, hashes `r4-binaries.sha256`. The separate `/private/tmp/ch18-live-r4.FWfiti` controller copies the corrected exact-book primary/three prompts, pins R4 CLI, retains all budgets/TLS/no-redirect/no-auto-retry/memory-only key controls, adds a pre-upstream illegal-marker guard, and removes an unnecessarily restrictive literal release-query requirement while keeping Go-only search. Syntax/keyless guard checks passed with zero keyreads/network/childlaunches; current primary/prompts verified. Root explicitly releases one R4 attempt under Bill’s existing Chapter18 live/API-key authorization: one modelavailability GET, <=12 POSTs at1024 nativeoutput, one search/scrape plus bounded MCPsetup, exact Sonnet4.6→Opus4.6 switch, ownedprocesses only. This is release, not success; previous failures remain.

## Final R4 code and required CLI-operations disposition

**CODE ACCEPTED; required bounded CLI live exercise COMPLETED.** No remaining blocker is established by the reviewed R4 source and evidence. Carry the completed whole-code/public-contract/ownership/original-reference/KISS review and the four production correction rounds. This disposition does not perform or claim author work, snapshot creation, commit, installation or overall chapter closure.

**Final local evidence read:** under `/private/tmp/ch18-validation.l8oKlb/`, `r4-root-race.log/.exit` shows all root packages passing with exit 0; `r4-root-vet.log/.exit` has no diagnostics and exit 0. `r4-binaries.sha256` records the freshly built `bin-r4/ch18` and `bin-r4/ch18-gui` identities. Actual final grade filenames are `grade18-v4.json` and `grade1-v2.json` (not `.log`); both `.exit` files are 0. Chapter 18 passes all 11 checks, **155/155, true**. Chapter 1 passes **100/100, true**, with exact **511/65** usage, append-only history, normal exit and both genuine **62-second header/body** waits. Root reports 124.49s elapsed; the inspected final stderr files are empty and no separate final elapsed file was present, so that total timing is root-reported, not independently read here. Unchanged GUI/auxiliary/standalone-JS broad evidence is carried with its documented fixture repairs, not silently relabeled as fresh R4 runs.

**Live evidence read:** `/private/tmp/ch18-live-r4.FWfiti/receipt.json` records **8 provider POSTs, all HTTP 200**: six Sonnet 4.6 and two Opus 4.6, with matching producer identities. Exactly five MCP POSTs cover initialize, initialized notification, list, one search and one scrape; no duplicate search/scrape is recorded. CLI exits **0, no signal**. Root separately reports one availability GET. Bill's explicit live/API-key authorization belongs to this bounded attempt; no historical OAuth authorization/result is borrowed.

The inspected `verification.txt` records root's independent checks of all five CLI status snapshots and both actual paired `agent_status` results against the appropriate response prefix, including last/session counts, cache-hit fraction and per-model pricing. It shows selected Opus changing at the human switch while the six-response totals remain unchanged, then accumulating the two Opus responses. Final totals are **14 fresh input / 39,375 cache-write / 36,687 cache-read / 1,895 output**, with the known-tariff model estimate **$0.23440875**. The saved verification reports all eight tool-marker placements legal, no external web declarations after unload, all status result pairs found and normal EOF/exit 0. These are inspected saved root verification results, not a reviewer replay or an independent read of every raw provider/tool body.

Root's saved content verification supports the fetched Go 1.27.3 release, required archive filename/citation, disclosed truncation and the extra **64-character** checksum match. Claims about omitted regions are not thereby proven. “Known” tariff arithmetic is not invoice proof and excludes Firecrawl billing and earlier attempts. The keyless live MCP session supplies no live bearer-auth claim; controlled redaction/auth tests remain the relevant carried evidence. No optional live GUI/browser/native-launch/install result is claimed.

**Failures preserved:** the earlier root-reconstructed-primary setup failure used four provider POSTs, no MCP, with a reported known-model estimate $0.03171465; it was not a book failure. The next attempt hit the real HTTP 400 marker defect after two provider POSTs, only one known-usage response ($0.00578025), no successful search/scrape/switch and SIGTERM cleanup. Neither estimate proves complete billing for those failed attempts. Earlier local grade/test/setup failures remain intact with their separate causes and repairs; none becomes a retrospective pass. No new reviewer execution, credential access, source edit or operational action. Submit for root read/JOIN; author/snapshot/commit still pending.
