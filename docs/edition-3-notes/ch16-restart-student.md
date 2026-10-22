# Chapter 16 lean restart — R3 compiling HOLD

Source frozen after production correction R3 of 10 in `solutions/edition-3/main`, extending accepted Chapter 15 HEAD `eb068a3bf78d8aeca9fe67aba8d7b32001b9fd26`. No commit, tag, installation, or restart.

## R1 — four authorized corrections

Ownership remains unchanged; no restart or framework. Changes are confined to:

1. `auth_command.go`: shared CLI/GUI `auth login --print-url` support and `EN_BROWSER_CMD` handling. The URL stays printable; the documented `EN_BROWSER_CMD=true` invocation suppresses browser launch. The URL is passed as a shell positional argument, not evaluated as shell text. Auth storage, validation, and issued-client ownership are unchanged.
2. `internal/engine/responses_stream.go` / `responses.go`: retain completed native-indexed `output_item.done` items when terminal output is empty/omitted; nonempty terminal output remains authoritative. Existing delta identities and the same parts parser are reused. Failed/incomplete/EOF/unaccounted/malformed terminal responses cannot commit or dispatch those items. Summaries, raw continuation fields, source provenance, and terminal usage are retained.
3. `internal/common/actor.go` / `auth.go`: exact supported Responses rows (`gpt-ch19-course`, `gpt-6.1-sol`, `gpt-5.6-sol`) own their reasoning defaults, summary support, and ordinary explanation policy in the existing model contract. Explicit supported `none` choices are honored (summary `none` omits that field); unknown/API-incompatible capability claims stay conservative. Summary-only turns dispatch for these rows without pretending the summary is ordinary explanation; retained ordinary gates still work.
4. Responses replay rejects duplicate calls/results and orphan results rather than overwriting/reusing them. Unique late results retain exact IDs and pairing; missing results retain the honest unknown-outcome placeholder. Rendering does not modify the supplied canonical dialogue. Existing final-response duplicate-call validation was not changed.

### R1 focused commands and outcomes

All execution was sequential/offline, with scrubbed environment, existing caches, owned `.ch16-work/home` and `.ch16-work/tmp`, dummy credentials, `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, `GOMODCACHE=/Users/bill/go/pkg/mod`, and `GOCACHE=/Users/bill/Library/Caches/go-build`. Every R1 test command used `-vet=off`.

- Before correction, `go test -vet=off ./internal/engine -run '^TestResponsesR1' -count=1` reproduced completed-item loss, capability/explanation/summary-none failures, and ambiguous replay acceptance. The same new cases pass after correction.
- Built owned executables using root `go build -o "$ROOT/.ch16-work/r1-bin/cli" ./cmd/ch01` and GUI-module `go build -o "$ROOT/.ch16-work/r1-bin/gui" ./cmd`. With `CH16_TEST_CLI`/`CH16_TEST_GUI` set to those paths, `go test -vet=off . -run '^TestCh16R1HeadlessCommands$' -count=1` reproduced both binaries' flag rejection; corrected binaries pass real delegated login against a signed local issuer, including issued-client reuse. The baseline suppression-only invocation (`-run '^TestCh16R1HeadlessCommands$/suppressed'`) also reproduced launch despite `EN_BROWSER_CMD=true`; final cases pass. A harmless owned browser sentinel, never a real browser, detects attempted launches.
- `go test -vet=off . -run '^TestCh16R1CompletedToolAdmission$' -count=1`: passed actual Engine/Agent/local HTTP/SSE admission, retained ordinary gates, terminal usage, and no dispatch/history tool facts after failed/incomplete/EOF/partial arguments.
- Combined affected run, with both owned-binary variables set: `go test -vet=off ./internal/auth ./internal/engine . -run 'Test(Auth|Responses|Ch16)' -count=1` — passed.
- Strengthened exact-ID replay check: `go test -vet=off ./internal/engine -run '^TestResponsesR1ReplayPairing$' -count=1` — passed. GUI-module `go test -vet=off . -run '^TestCh16ResponsesPublicWire$' -count=1` — passed.
- Explicit module compilation, sequentially in root, gui, ch06, ch07, examples/ch04-consumer, examples/ch13-speech: `go test -vet=off ./... -run '^$'` — all six passed. This is compilation-only evidence.
- Final `git diff --check -- .` clean; approximately 17.4 GiB available.

R1 setup failures were a test `AuthStatus.Subject` field mismatch and incorrect fixture host/scope names; initial production edit compile failures were a missed map declaration/raw-item reference and missing `os` import. These were corrected and are not claimed as defect sensitivity. The flag and suppression defects were separately reproduced through actual binaries after fixture repair. Tests live in `internal/engine/responses_r1_test.go`, `chapter16_r1_agent_test.go`, and `chapter16_r1_command_test.go`.

No children, live provider/browser activity, full suite, race, vet, grading, lifetime campaign, installation, cleanup, commit, or restart in R1. Retained source review and broader gates remain deferred. Local synthetic identity/HTTP checks are not proof of live interoperability.

## Round 0 ownership and held implementation

Ensemble owns one auth service; Engine resolves credentials at its HTTP boundary. Auth is neither History nor a tool. Existing Agent admission/History persistence remain authoritative; GUI receives only public projections.

- New `internal/auth/{auth,store,login,identity}.go`, public `auth.go`, and shared `auth_command.go`: protected installation/issued-client identity and per-account grants, process-shared writer locking, atomic rotation, PKCE loopback login, signed identity validation, inference-scope gating, persisted billing selection, status/models/logout, and bounded same-source renewal. Both CLI entry points delegate to that owner.
- `internal/common/auth.go`, configuration/construction and persistence owners: public `gpt-ch19-course` capability contract, explicit surface/stream controls, configurable reasoning, conservative unknowns, non-secret saved policy. Plan streaming contradictions fail explicitly. Automatic paid fallback is deliberately deferred/rejected.
- Engine Responses renderer/parser and existing transport/cache owners: public Responses route, plan namespaces and no plan cache markers, API-only explicit cache markers, streaming/text/tools/usage, portable public summaries, source-bound opaque continuation, and redacted diagnostics. Shared SSE framing is reused.
- Existing History/Agent part handling and GUI projection/browser presentation: summaries remain distinct from text; continuation becomes a non-secret placeholder; billing mode is visible without entering conversation facts.

## Round 0 focused offline results

Commands used isolated HOME/TMP under `.ch16-work`, scrubbed environment, dummy credentials, offline module settings, and local ephemeral HTTP peers only.

- `go test ./internal/auth -run '^TestAuth' -count=1`: passed signed synthetic login/claims, protected storage, actual cross-process refresh rotation, invalidation, missing inference grants, and transient-refresh preservation.
- Final affected checks: `go test -vet=off ./internal/engine . -run 'Test(Responses|Ch16)' -count=1`, followed by `go test -vet=off . -run '^TestCh16' -count=1`: parser checks passed; final root checks passed after correcting the affected fixture. Coverage includes real Agent/tool/save/replay, plan/API transport pairing, renewal, shared auth commands, override/default routing, usage and diagnostic privacy.
- GUI: `go test . -run '^TestCh16ResponsesPublicWire$' -count=1`: passed real loopback GUI replay projection of summary/answer/usage/billing with no private continuation/token bytes.
- Compilation only, sequentially in `.`, `gui`, `ch06`, `ch07`, `examples/ch04-consumer`, and `examples/ch13-speech`: `go test -vet=off ./... -run '^$'` — all six passed. This is compilation evidence, not full behavior coverage.
- `git diff --check -- .`: clean. Final disk availability approximately 17.8 GiB.

## Round 0 failures, limitations, assistance

Intermediate failures were retained in command results. They included mechanical adapter/test mismatches, fixture permissions/visibility/scope mistakes, and real integration gaps subsequently corrected (Surface persistence, `off` handling, resolved replay identity, and default Responses endpoint). Those setup/compile failures are not claimed as negative-sensitivity tests. Initial speculative auth/model assumptions were replaced with the chapter's issued-client and synthetic capability contract before HOLD.

No child assistance, external documentation fetch, live OAuth/account/model call, browser/native/speech run, full suite, race, standalone vet, grade, or lifetime campaign. Early focused `go test` commands did not explicitly suppress Go's implicit vet pass; final affected checks and all six compile-only commands used `-vet=off`. An untracked-name enumeration accidentally listed old `.ch15-work` filenames; none of their contents were opened. No original sibling solution/private grader/reviewer content was read.

Live interoperability, broader regression/security validation, and expensive gates remain deferred to held-source review. Local synthetic OAuth is not proof of real-provider login. Raw temporary profiles, saves, settings and wire data remain outside tracked source in owned temporary areas; no settings fixture was added to source.

## R1 test-only trace-capture race repair — stable HOLD

Authorized narrow repair only in `chapter16_test.go`, for `TestCh16ResponsesAgentToolRecovery`; no production source changes and no new production correction round (still R1, 1/10). The test already called Shutdown before inspecting diagnostics; that ordering was not the defect. Its caller-owned capture now locks both Write and String, and takes one snapshot. A nonblocking receipt signal plus a bounded two-second select waits for the continuation-bearing tool-result request diagnostic before Shutdown; there are no sleeps or production queue/drain changes. Existing credential/continuation absence and archive/replay assertions remain intact.

From `solutions/edition-3/main`, the sole test command was `go test -race -vet=off . -run '^TestCh16ResponsesAgentToolRecovery$' -count=1`, under the prior scrubbed offline/cache environment with `HOME=$PWD/.ch16-work/home` and `TMPDIR=$PWD/.ch16-work/tmp` (non-symlink owned paths). It passed, exit 0, on the first post-repair run. Receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r1-test-trace-race-repair.log`. Disk remained above the 1 GiB floor. The coordinator's prior race finding was reused; no broad gate, auth-storage weakening, runtime amendment, or concurrent-work edits were attempted. Root owns inspection and the broad rerun with canonical environment paths.


## R2 — incremental framing and explicit Responses cache envelope

Stable compiling HOLD, correction 2/10. Existing owners only: `internal/engine/engine.go` and `responses_stream.go` sniff a small initial SSE field prefix when a successful Responses response lacks a useful SSE media type, then use the existing incremental parser. No whole-body read precedes those deltas. Ordinary JSON, readable JSON refusal, malformed/error handling, terminal/usage gating, R1 completed-item fallback and nonempty terminal authority remain intact; other vendor dispatch is unchanged.

`internal/engine/breakpoints.go` now puts the exact supported API-route constitution in one developer-role input content block with its legal explicit marker and paired top-level `prompt_cache_options.mode="explicit"`. Plan and unsupported routes omit both. The existing `cache_wire.go` lens strips that top-level option alongside legal content markers, preserving opaque provider fields. No capability names/prefixes/prices were added. Tests assert exact constitution/nonmarker contents, unchanged canonical context, legal retained-history boundaries, the marker budget, and existing marker accounting.

New focused tests are in `chapter16_r2_test.go` and `internal/engine/responses_r2_test.go`; the existing `chapter16_auth_test.go` real-auth HTTP comparison now checks actual API/plan/near-name unsupported envelopes. The local HTTP framing peer withholds the stream remainder until the first callback, bounded by a timeout (not a sleep). Controls cover absent/plain/misleading headers, complete JSON, malformed JSON/SSE, readable refusal, and HTTP/provider errors.

All commands used the existing scrubbed `env -i` offline/cache environment recorded above, `GOMAXPROCS=2`, readonly/-p=1 flags, and non-symlink owned `.ch16-work/home` and `.ch16-work/tmp`. No real credentials or inherited LLM_BASE_URL.

- Before correction: `go test -vet=off . -run '^(TestCh16R2ResponseFraming|TestCh16PlanBillingTransport)$' -count=1` — exit 1, reproducing all three delayed-callback/header failures and the system-role/missing-options cache envelope. Receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r2-baseline.log`.
- The same command after correction — exit 0. Receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r2-focused.log`.
- `go test -vet=off ./internal/engine -run '^TestResponsesR2CachePreservesContent$' -count=1` — initially exit 1 because the existing lens retained the newly added top-level mode; corrected by its Responses-specific span. Concrete failure receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r2-lens-failure.log`.
- Final affected check: `go test -vet=off ./internal/engine . -run '^(TestResponsesR1CompletedItems|TestResponsesTerminalAndPrivateParts|TestResponsesR2CachePreservesContent|TestCh16R1CompletedToolAdmission|TestCh16R2ResponseFraming|TestCh16PlanBillingTransport)$' -count=1` — exit 0, both packages passed. Receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r2-affected.log`.
- Sequential root and GUI compilation only: `go test -vet=off ./... -run '^$'` — both exit 0. Receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r2-compilation.log`. No public interface changed; inherited-module results were carried rather than repeated. Compilation is not behavior coverage.

Setup mistakes are separate: an unavailable public `Response` type name and an unnecessary nonexistent `Instructions` fixture field caused compile failures (`r2-before.log`, `r2-lens-before.log`, `r2-lens-baseline.log` in the same owned receipt directory). The first refusal control incorrectly expected a parser error; the established readable-refusal behavior was preserved and its assertion corrected. These are not claimed as defect sensitivity. No mutation campaign occurred.

Final whitespace check clean; approximately 22.4 GiB available. No full suite, race, vet, grading, live-provider/browser activity, children, runtime/book changes, commit, installation, or restart. R1's other areas and the test-only synchronized diagnostic capture remain carried unchanged. Root and retained reviewer own renewed broader gates.


### R2 validation follow-up — test-only assertion repair (still 2/10)

Stable HOLD. Only `chapter04_round1_test.go` and `ensemble_test.go` changed; no production code or separate handoff work changed. The shutdown test still observes real late output, checks its exact on-disk content, requires Shutdown to join, and forbids killed events. Its terminal done/exit-99/empty-reason assertion now uses fresh `job.Data()` after the successful join, rather than a legitimately earlier output-wake report.

The cancellation test separately diagnoses the expected `context.Canceled` result and empty completed-text History. It joins the actual worker with a bounded channel wait, then requires exactly one stale worker error receipt in canonical History. Transport cancellation requires zero usage; only the specific canceled HTTP-200 body read/decode receipts (and only the body-phase case) require all-unknown usage. Unexpected receipts fail. This is not a blanket zero/unknown alternative, and canonical History is not required to be empty. No sleeps or production synchronization changes.

Commands used the same scrubbed offline environment, existing Go caches, and canonical owned HOME/TMP described above; disk remained above 1 GiB:
- `go test -race -vet=off . -run '^(TestChapter4Round1KilledMetadata|TestCancellationWhileWaiting)$' -count=3` — exit 0 (2.562s). Exact receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r2-test-assertion-repair-race.log`.
- Root compilation only: `go test -vet=off ./... -run '^$'` — exit 0. Exact receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r2-test-assertion-repair-compile.log`.

Both checks passed on the first repair attempt; whitespace check clean. The retained failed root gate `/private/tmp/ch16-r2-validation.yZZ4zC` was not read, altered, or removed. No broad gates, vet, browser/provider traffic, credentials, children, commits, or new production round. Root owns inspection and resumption of its stopped gate; production binaries can remain carried.


## R3 — public-summary checkpoint catch-up (3/10)

Stable compiling HOLD. Only production change: `gui/hub.go` applies its existing final unseen-suffix catch-up to agent `reasoning_summary` as well as ordinary `text`, and emits the final's actual kind under its existing stable ID. Prefix comparison, ordering before final, partial deletion, cursor/revision idempotence, bounded queues, and reconnect projection are unchanged. No parser, auth, Engine, canonical History, or hidden/opaque reasoning changes.

`gui/chapter16_r3_test.go` drives the real Hub/client queues with controlled immutable checkpoints and canonical Agent History. Both entirely unseen and partially seen summary finals must deliver exactly the missing summary bytes before their final, preserve the adjacent tool and separate ordinary text, retain IDs, emit no duplicate suffix/final on repeated checkpoints, and reconnect with finals rather than stale chunks. It checks early prefix delivery, private continuation absence, and unchanged canonical events. Existing public WS projection and coherent ordinary-text handoff tests also pass.

Commands ran from `main/gui`, sequentially in the same scrubbed offline environment/caches and owned canonical HOME/TMP; disk was above 1 GiB (about 22.2 GiB at entry):
- Before fix: `go test -race -vet=off . -run '^TestCh16R3SummaryCheckpointCatchup$' -count=1` — exit 1 for missing summary bytes in both cases. Clean baseline receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r3-gui-baseline-corrected.log`. The first receipt, `r3-gui-baseline.log` in that directory, additionally exposed an incomplete fixture: opaque Responses provenance was missing, so its existing fallback placeholder was `thinking`. The fixture was corrected; production opaque handling was not changed or classified as defective.
- Affected checks: `go test -race -vet=off . -run '^(TestCh16R3SummaryCheckpointCatchup|TestCh16ResponsesPublicWire|TestCh7R1CoherentFinalHandoff)$' -count=1` — exit 0. Receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r3-gui-affected-race.log`.
- GUI compilation only: `go test -vet=off ./... -run '^$'` — exit 0 for GUI and GUI/cmd. Receipt: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch16-work/r3-gui-compilation.log`.

Whitespace check clean; no checks pending. No broad gates, vet, grading, native/browser/provider activity, children, book/guide/runtime edits, commits, installation, or restart. Unchanged root production/artifacts and evidence remain carried for coordinator review; compilation is not behavior coverage.
