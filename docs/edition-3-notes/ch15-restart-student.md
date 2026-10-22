# Chapter 15 restart — student HOLD

**Stable compiling HOLD; formal production corrections 0/10.** Source is held in `solutions/edition-3/main`; no snapshot, commit, tag, installation, deployment, or restart was made. Root and the retained reviewer should inspect this source before expensive gates.

## Ownership and changes

Engine owns the send-boundary cache child, exactly two retained requests, detached structural diagnostics, and current-process incurred accounting. Common contains value contracts and explicit capability/price lookups; root wires the child. GUI reads Engine snapshots, not History totals. History retains measured token facts (including `-1` unknown), never dollars or lens state. There is no cache service, alternate client, sibling dependency, or auxiliary request path added.

Changed production paths:
- `ensemble.go`; `internal/common/{common,history,cache}.go`.
- `internal/engine/{engine,parse,render,tools,cache,accounting,breakpoints}.go`.
- `internal/history/{history,persistence}.go`.
- `gui/{hub,projection}.go`; `gui/web/gui/{app.js,index.html,style.css}`.

The lens sees the exact actor-frozen bytes immediately before HTTP dispatch, not previews/replays. Anthropic Messages markers use at most two baseline and two legal conversation blocks; manuals remain SYSTEM and handoff/ephemera retain their accepted projection. Only explicitly supported Sonnet 4.6 and Opus 4.6 get markers/prices. Default-five-minute cache and token rates were verified by source-fetch of the chapter-linked `https://docs.anthropic.com/en/docs/build-with-claude/prompt-caching`; an unverified Haiku 3 row was removed. Other model/API prices remain honestly unknown. Costs accumulate under captured API/served-model identity rather than the active model's price. GUI has the live four-number meter and explanatory tooltip, including partial/unknown costs and local-prefix-versus-provider-hit distinctions.

New tests: `chapter15_cache_test.go`, `internal/engine/cache_test.go`, `gui/ch15_meter_test.{go,cjs}`. Updated affected expectations/fixtures in `chapter02_round{2,3}_test.go`, `chapter06_{adjacency,token_limit}_test.go`, `chapter14_instruction_test.go`, `ensemble_test.go`, and `internal/engine/sse_test.go`. These retain earlier replay/error/stream/manual assertions while accepting unknown cache measurements and legal SYSTEM block encoding.

## Actual verification

Commands below used `.ch15-work/offline` (absolute path when inside another module). It runs `env -i` with owned `.ch15-work/{home,tmp}`, PATH `/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin`, dummy Anthropic/OpenAI/Gemini/Google keys, and no generic LLM_BASE_URL or inherited credentials. Required settings: `GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod GOCACHE=/Users/bill/Library/Caches/go-build`. Go execution was sequential.

Final focused passes:
```
go test -vet=off . -run 'TestChapter15|TestChapter14|TestRound2LiveAccountingIsNotReplay|TestRound2KnownChargesSurviveResponseErrors|TestRound3UsageDependencies|TestCompletionTokenLimitHTTP|TestOpenAIBusyHintToolAdjacencyHTTP|TestConversationAndOwners|TestResponseErrorsAndUsage' -count=1
go test -vet=off ./internal/engine -run 'TestChapter15|TestChapter14ProviderContext|TestPlainAndStreamFinalParity|TestInterruptedToolNotFinal' -count=1
# cwd gui:
go test -vet=off . -run 'TestChapter15|TestChapter14' -count=1
# cwd main:
node gui/ch15_meter_test.cjs
node gui/ch14_context_test.cjs
node --check gui/web/gui/app.js
```

Explicit final compilation: in each of `.`, `gui`, `ch06`, `ch07`, `examples/ch04-consumer`, `examples/ch13-speech`, ran `go test -vet=off -run '^$' ./...` through the same wrapper. **All six passed**, with no tests selected. Free space was checked before each module; final available space was 17,983,004 KiB, safely above the 1 GiB floor. `git --no-pager diff --check -- .` passed.

Earlier incremental runs used `go test -run '^$' ./internal/engine ./internal/common`, owner selection above, root `go test . -run '^TestChapter15' -count=1`, GUI `go test . -run '^TestChapter15' -count=1`, and progressively expanded subsets of the final root selection. Those initial commands lacked `-vet=off`, so Go's default package vet may have run; no explicit vet gate was run. Final commands explicitly disabled it.

Failures/corrections: test fixtures initially used nonexistent Config/Agent methods, the wrong opaque-part construction, exact floating equality, wrong GUI Agent ID, a nonexistent loop variable, and `Prices` instead of `LookupPrices`. The initial send test wrongly expected a dispatch-time stream rewrite; digest/length proved the baseline sends its frozen request unchanged. Review caught and fixed the missing actual HTML meter node and fresh bucket `null` serialization. Earlier tests exposed superseded zero assumptions and a string-only SYSTEM decoder; these were updated, and the negative-usage error wording preserved. No unresolved focused-test or compilation failure remains; no dependency/setup blocker occurred.

## Limits and deferred work

No private/reference implementation, other sibling snapshot, grader implementation, or root/reviewer memory was inspected. No children or external assistance were used; no consequential public ambiguity required a question. Only public provider documentation was fetched, never inference/discovery/account endpoints.

Local HTTP fixtures and production JS/DOM projection checks are not observed browser rendering, real provider caching, or acceptance. No live GUI/8088, browser, audio/MIC/device activity, paid/model call, full suite, race/lifetime campaign, grade15, or same-held-source CLI Chapter 1 gate was run. Root's original comparison, retained-reviewer review, explicit-GUI grade15 and same-held-source CLI Chapter 1, and any broader validation remain deferred. Source/tests are stable at this HOLD.

---

## R1 — stable compiling HOLD (production correction 1/10)

This section supersedes the initial implementation claims above without rewriting that history. The existing Engine/cache/accounting owners and transport remain; no restart, framework, shadow client, or parallel accounting path was introduced.

**R1 changed paths:** `internal/common/{actor,cache}.go`; `internal/engine/{breakpoints,cache,cache_wire,engine,parse,render,sse}.go`; `gui/{cmd/main,hub,projection}.go`; `gui/web/gui/{app.js,index.html}`. Tests: new `internal/engine/ch15_r1_test.go`, `chapter15_r1_test.go`, `gui/ch15_r1_{wire_test.go,meter_test.cjs}`; revised `internal/engine/cache_test.go` and `gui/ch15_meter_test.{go,cjs}`.

- Exact GPT-5.6 / Chat Completions policy: message-text markers only, request `prompt_cache_options:{mode:"explicit"}` only after a marker lands; no tool-definition markers or provider-wide opt-in. Unsupported rows/APIs retain implicit fallback. GPT-5.6 prices remain unknown; no price rows were invented.
- Boundaries are baseline SYSTEM (and Anthropic last tool), one last eligible retained-history endpoint, and the actual typed surviving handoff at the ephemera front. The handoff is split before flattening/coalescing; transient tails and protected native blocks are not substitute targets. Manuals remain SYSTEM. Explicit-cache routes retain stable block shapes; Anthropic markers are inserted without re-encoding the marked block.
- The Engine child retains two actual sends and writes private readable `cache-current.txt` / `cache-prior.txt` under Agent DataDir; GUI defaults DataDir to `cr`. New processes write honest no-send placeholders, not restored request claims. Raw-byte spans strip only recognized marker members at legal positions; other bytes, key order, whitespace, literals and tool data survive. Tools/system/messages, non-prompt options, first byte difference and prefix estimates are separate; route/model changes reset prediction eligibility. `Observe` returns a request-local token; response reporting uses that token without holding the lens lock over HTTP.
- Qualified OpenAI `prompt_tokens_details.cache_write_input_tokens` is parsed as exact int64. Input is total minus read minus write only when all operands are valid/consistent. Missing/invalid dependencies remain unknown while independent facts survive; final cumulative SSE usage is not added to earlier snapshots.
- The GUI now emits/consumes flat `usage` frames (`input/output/cache_read/cache_write`, `total_*`, `hit_rate/cost_usd/priced`) from Engine snapshots. Visible values are session cost, latest all-input, latest cache-read, session hit rate. Tooltip retains totals, served-model buckets, unknown/partial costs and request-local prediction-versus-measurement. No GUI accumulation counters were added.

### R1 actual checks and failures

All R1 Go commands used `.ch15-work/offline` (absolute wrapper path in other modules), now with **no key environment variables**, no inherited environment or LLM_BASE_URL, the same required offline/cache/Go settings recorded above, and owned HOME/TMP. The wrapper checks the 1 GiB free floor before every command. Final free space: 17,690,892 KiB.

Pre-fix owner tests reproduced missing OpenAI opt-in, the arbitrary second-history boundary, canonicalization hiding byte changes, and forced-zero qualified writes. After fixing its fixture APIs, the local-HTTP test reproduced missing diagnostic files. Receipts remain in `.ch15-work/r1-prefixed-{tests,files,meter}.txt` and the tool transcript. The pre-fix JS attempt had fixture/setup errors, not a valid behavioral-failure receipt.

During R1, tests caught and corrected marker-induced key reordering, missing measured-response correlation, and fixture/edit mistakes (wrong import/handler, enum/string/API assumptions, source anchors, temporarily removed adjacent helpers and pointer construction). The old path-only lens assertion was updated to the inspectable byte-difference contract. No unresolved focused-test or compilation failure remains.

Final successful focused commands (through the wrapper, cwd main unless noted):
```
go test -vet=off ./internal/engine -run '^(TestChapter15|TestPlainAndStreamFinalParity|TestInterruptedToolNotFinal)' -count=1
go test -vet=off . -run '^(TestChapter15|TestChapter14|TestRound3UsageDependencies|TestRound3SavedUsageIsNotRecharged|TestCh2R2KeepsMeasuredBucketsOnIncompleteUsage|TestCh6UsageAdjacentToToolChain|TestCh6TokenLimitPreservesUsageBeforeExposingError|TestConversationAndOwners|TestResponseErrorsKeepMeasuredUsage)' -count=1
go test -vet=off . -run '^TestChapter15R1' -count=1
# cwd gui:
go test -vet=off . -run '^TestChapter15' -count=1
# standalone production-JS checks also passed:
node gui/ch15_r1_meter_test.cjs
node gui/ch15_meter_test.cjs
```
The broad root regex above selected Chapter 14/15 and existing matching tests; regex alternatives that have no matching test are not additional coverage claims. The final R1-only run followed the last test-only pretty-send refinement. GUI tests include a real local WebSocket and pass actual initial/live/partial/reconnect frames to the production socket/DOM JavaScript. Distinct latest versus total values, unknown served-model pricing, restart zeroes, raw pretty HTTP bodies, marker-looking literals, protected thinking, cumulative SSE precision, and delayed response correlation were exercised.

Final sequential compile command:
```
base="$PWD"
for module in . gui ch06 ch07 examples/ch04-consumer examples/ch13-speech; do
  (cd "$module" && "$base/.ch15-work/offline" go test -vet=off -run '^$' ./...) || exit
done
```
**All six modules passed**, no tests selected. `git --no-pager diff --check -- .` passed. Earlier initial-HOLD evidence for unchanged paths was retained rather than re-running a suite.

### Public inputs, scope and deferred acceptance

Read-only documentation fetch: the chapter-linked OpenAI caching guide, `https://developers.openai.com/api/docs/guides/prompt-caching`, confirmed the explicit marker mode/budget. Python TLS validation failed; it was not bypassed. System curl fetched the guide successfully; the linked Chat Completions reference returned 404. No discovery, account, inference, credential access, or paid call occurred. No private review/reference input or assistance was used.

Scope slip: one `git status --short` / `git diff --check` invocation omitted the main pathspec and returned unrelated repository pathnames. No outside-main diff text was returned or used; no private path listed by that output was subsequently opened or used as source. No outside-main source edits were made (only this permitted note was updated). Subsequent Git inspection was bounded with `-- .`.

No full suite, race, vet gate, grade/lifetime campaign, browser/native/audio/device activity, live GUI/controller, children, installation, cleanup, commit or snapshot was performed in R1. Local fixtures and DOM tests do not prove real provider cache grants or visual GUI acceptance. Root's whole-source inspection and the **same retained reviewer's delta verification** remain required before broader validation, explicit-GUI grade15, or same-held-source CLI Chapter 1. Source/tests are stable at this R1 HOLD.

---

## R2 — stable compiling HOLD (production correction 2/10)

**Scope/paths:** only `internal/engine/cache.go`, `internal/engine/render.go`, new `chapter15_r2_test.go`, this note, and owned temporary receipts changed in R2. Existing ownership, transport, accounting, GUI, History, and raw-span helpers remain intact.

1. Diagnostics now print actual marker-free raw **TOOLS / SYSTEM / MESSAGES** spans in provider cache order, followed by separate **NONPROMPT OPTIONS**, before the retained exact original-wire and marker-free whole-request views. Absent sections are labeled; spans are not misleadingly presented as independent JSON. No JSON roundtrip or reporting framework was added.
2. OpenAI's existing disposable tool-result adjacency projection now runs before selecting the final history marker and before appending handoff/transient ephemera. The request sends that already-projected message list. The later interleaved hint therefore receives the history marker, not a result moved ahead of it; canonical facts/provenance are unchanged.

**Failure receipts:** `.ch15-work/r2-prefixed-tests.txt`, `r2-prefixed-confirmed.txt`, and `r2-prefixed-isolated.txt` retain the pre-fix runs (exit 1). Initial fixture corrections supplied tool-call provenance, excluded the existing SYSTEM declaration from causal-dialogue assertions, and supplied explicit resolved-model evidence for native replay. The isolated run then failed on exactly the two production findings: missing raw TOOLS body and marker on the moved tool result instead of the later hint; unsupported-model fallback passed. The first edit script applied the diagnostic correction but stopped at a mismatched renderer anchor; `.ch15-work/r2-postfixed-tests.txt` records the still-failing marker check (exit 1). After moving the actual initializer call, `.ch15-work/r2-final-focused.txt` records both R2 tests passing (exit 0).

The new real-HTTP/files test deliberately serializes messages before system/tools, checks exact raw named bodies/options in cache order, preserves the original transmitted bytes, rotates the prior report, and verifies one changed section plus inspectable first difference. The Render/Agent-HTTP test checks complete call→hint→real-result canonical facts versus adjacent result→later-hint wire order, native reasoning/argument precision, separate handoff, unmarked transient tail, implicit fallback, and unchanged History/events/source context.

**Actual successful commands**, cwd main, through the unchanged `.ch15-work/offline` wrapper (cleared/keyless environment, required sequential offline Go settings/caches, owned HOME/TMP, 1 GiB floor checked before every invocation):
```
go test -vet=off . -run '^TestChapter15R2' -count=1
go test -vet=off ./internal/engine -run '^(TestChapter15R1MarkerPositions|TestChapter15R1MovingMarkerStablePrefix|TestChapter15R1ActualHandoffAndFallbackPositions|TestChapter15LegalMarkersPreserveContent)$' -count=1
go test -vet=off . -run '^(TestOpenAIBusyHintToolAdjacencyHTTP|TestChapter15R1DiagnosticFilesActualSend)$' -count=1
```
The interrupted inspection/check call did not dispatch; completed R2 tests were not repeated after the status interruption. The pending affected checks above passed. Bounded `git --no-pager diff --check -- .` passed.

Final sequential compile-only pass (all six passed; no tests selected):
```
base="$PWD"
for module in . gui ch06 ch07 examples/ch04-consumer examples/ch13-speech; do
  (cd "$module" && "$base/.ch15-work/offline" go test -vet=off -run '^$' ./...) || exit
done
```

**Receipt clarification for retained R1 history:** job 3547 was mixed, not all green: Engine fixture compilation failed while root/GUI passed. Job 3549 subsequently passed corrected Engine/root checks; job 3550 passed compilation. R2 runs explicitly preserve/check failed subcommand statuses; final pending-check/compile job 3633 exited 0.

Unchanged R1 GUI/accounting/replay evidence was carried forward, not rerun as a behavior campaign. No other production corrections, external network/provider/key access, browser/native activity, children, full suite/race/vet/grade/lifetime gates, manuscript/grader/guide edits, installation, commit, snapshot, or cleanup occurred. Local fixture results do not establish provider cache grants or acceptance. **Source/tests are held; same-reviewer R2 delta verification and root approval remain required before validation.**

---

## R3 — stable compiling HOLD (production correction 3/10)

**Narrow delta:** `gui/projection.go` now emits latest disjoint `input/output/cache_read/cache_write` at top level and `total` as the existing Engine session `Usage` value. The four `total_*` aliases were removed. `gui/web/gui/app.js` reads tooltip latest counts from the public frame and session counts from `frame.total`; auxiliary meter data remains for model/cost/cache diagnostics. No Engine/accounting owner, price/model policy, or four-visible-value calculation changed.

Focused fixtures changed: `gui/ch15_r1_wire_test.go`, `gui/ch15_r1_meter_test.cjs`, `gui/ch15_meter_test.cjs`. The real WebSocket test independently requires a non-null `total` with exactly four disjoint fields and rejects flat aliases on initial/live/reconnect frames. It still passes actual emitted frames to production JavaScript. Consumer checks prove auxiliary `meter.last/session` cannot override public counts, and cover unknown counts, partial price, initial zeroes and reset. Existing latest-versus-session and reconnect assertions remain.

**Before-fix failures:** through `.ch15-work/offline`, cwd `gui`:
```
go test -vet=off . -run '^TestChapter15R1RealUsageWireAndConsumer$' -count=1
node ch15_r1_meter_test.cjs
```
Both exited 1 for the intended defects: real usage frame lacked nested `total`; tooltip displayed deliberately differing auxiliary counts. Preserved in `.ch15-work/r3-prefix-wire.txt` and `r3-prefix-consumer.txt`. No fixture/setup or subsequent implementation failures occurred in R3.

**After-fix checks, all passed**, same wrapper/cwd:
```
go test -vet=off . -run '^(TestChapter15R1RealUsageWireAndConsumer|TestChapter15LiveMeterAndReplay)$' -count=1
node ch15_r1_meter_test.cjs
node ch15_meter_test.cjs
go test -vet=off -run '^$' ./...
```
The last command compiled GUI and GUI/cmd with no tests selected. Bounded `git --no-pager diff --check -- gui` passed. Wrapper retained cleared/keyless environment, required offline Go/cache settings, owned HOME/TMP, and the 1 GiB pre-command free-space check. Local HTTP/WS fixtures only; no external network/provider use.

Unchanged R2 Engine/root checks, R1 accounting/restart evidence, and the R2 six-module compilation are carried forward, not rerun. No grader/root fake-provider inspection or alignment, unsupported alias/price additions, framework/second owner, broader suite/race/vet/grade/lifetime/operational/browser/native activity, children, installation, commit, snapshot or cleanup. Only the listed GUI source/tests, this allowed note and owned temporary receipts changed. **Steady HOLD: same retained reviewer must verify the R3 delta before broader validation.**

---

## Released validation — STOPPED at first gate; final HOLD

Held source was not edited. The existing cleared/keyless offline wrapper was verified; required Go/cache settings and owned HOME/TMP were retained. Pre-gate free space was 17,357,556 KiB (wrapper enforces the 1 GiB floor).

**Actual first command**, cwd `/Users/bill/projects/ensemble/solutions/edition-3/main`:
```
.ch15-work/offline go test -race ./...
```
**FAILED, process exit 1.** Full command/cwd, unabridged output and exit status are preserved at:
`main/.ch15-work/validation-r3/root-race.log`.

First reported failure: `TestChapter4Round1SetupFailureCorrelation/mkdir`, `chapter04_round1_test.go:505`, opening its temporary `cr` path failed with **“is a directory.”** The same single invocation also reported:
- `TestChapter12UnloadAnthropicCatalog` (both variants): cannot unmarshal an array into the SYSTEM string field.
- `TestFailureDoesNotCommitHistory`: history/measured-usage assertion.
- `TestTransportErrors/unreachable=false`: assertion with returned usage `{-1,-1,-1,-1}`.
- CLI `TestLargeJSONLine` and `TestInputLines` (ordinary/final-without-newline): output assertions with cache-write/read `-1`.

These are observed failures, **not yet classified** as product defects versus stale tests or fixture setup. No race-detector warning appears in the captured output, but the gate failed; passing packages within it do not make the command successful. The command was allowed to finish normally; tool observation waits did not cancel it.

Stopped without retry, repair, source/test inspection for diagnosis, or later gate. **Root vet, all GUI/separate-module validation gates, frozen binary builds/hashes, Chapter 15 grader, and Chapter 1 grader were NOT RUN.** No new binary or acceptance claim was produced. Earlier focused/compile receipts remain unchanged. No production/test edits, children, installations, cleanup, external network/provider execution, browser/native activity, commits or snapshots occurred.

Production correction count remains **3/10**. Final HOLD awaits root's failure classification and explicit next release; code/operational acceptance remains root-owned.

---

## R4 — stable compiling HOLD (production correction 4/10)

**Production ownership/delta:** only `internal/engine/parse.go` and `internal/engine/sse.go`. In a valid usage object, omitted optional cache counts are again **zero** under inherited Chapter 2; supplied invalid/null/negative fields remain unknown/errors. Missing required counts and missing/null/invalid whole usage remain unknown/errors. Invalid OpenAI detail containers cannot fall through to optional zeroes; only dependent uncached input is tainted, preserving independent valid counts. Qualified current write fields and unsupported fallback remain unchanged; no aliases/prices/accounting owner added. SSE retains supplied whole usage (including invalid/null), merges valid partial objects, and shares the final parser; it no longer fabricates an empty usage object for complete absence or silently discards invalid supplied usage.

**Permitted fixture alignments:** `chapter04_round1_test.go` now obstructs `cr/io` after both constructors, preserving cache diagnostics and all effects/admission/pairing/limits/job assertions. `chapter12_unload_test.go` decodes actual SYSTEM string or text blocks; all integration assertions remain. `ensemble_test.go` changes only unreadable-HTTP200 transport usage to wholly unknown (unreachable still zero); sanitized error/no-answer/no-completed-history checks remain. `TestFailureDoesNotCommitHistory` and all CLI successful output/large-input/EOF/read-error assertions were not changed.

Tests previously edited under the blanket-unknown assumption were restored in `chapter02_round2_test.go`, `chapter02_round3_test.go`, `chapter06_adjacency_test.go`, `chapter06_token_limit_test.go`, the corresponding usage assertions in `ensemble_test.go`, and Engine `cache_test.go`, `ch15_r1_test.go`, `sse_test.go`. `chapter15_cache_test.go` now uses an actually invalid supplied cache count for its partial/unknown persisted-fact case. New `internal/engine/ch15_r4_usage_test.go` covers absent/invalid/required/whole usage plain–SSE parity, independent counts, qualified writes, invalid-final versus valid-later cumulative usage, and unsupported fallback. No other production paths changed.

**Failure/assistance:** first focused command failed compiling my new fixture (nonexistent `common.DefaultSurface`; wrong `parse` signature), before behavioral tests ran. Stopped and asked root. Root verified/classified these as new-test setup errors and authorized only fixture API corrections, with no extra production round. Original `.ch15-work/r4-engine-focused.log` is preserved (exit 1). The corrected `.ch15-work/r4-engine-focused-v2.log` passed (exit 0). Its regex contained two non-matching older-test names; verified names were subsequently selected in `r4-engine-focused-v3.log`, which passed. No unexpected behavioral failures occurred.

**Actual focused commands**, main cwd via unchanged `.ch15-work/offline`:
```
go test -vet=off ./internal/engine -run '^(TestChapter15R4.*|TestChapter15MeasuredUsage|TestChapter15R1QualifiedWriteAndCumulativeParity|TestPlainAndStreamFinalParity|TestInterruptedToolNotFinal)$' -count=1
# Above: first fixture compile failure, then setup-corrected PASS (R4 matrix + matching SSE tests).
go test -vet=off ./internal/engine -run '^(TestChapter15Usage|TestChapter15R1QualifiedOpenAIWrites|TestChapter15R1WriteDependenciesAndSSE|TestChapter15CostsKeepCapturedIdentity|TestChapter15WeightedRatesAndUnknownPrices)$' -count=1
go test -vet=off . -run '^(TestChapter4Round1SetupFailureCorrelation|TestChapter12UnloadAnthropicCatalog|TestFailureDoesNotCommitHistory|TestTransportErrors|TestConversationAndOwners|TestResponseErrorsAndUsage|TestRound2KnownChargesSurviveResponseErrors|TestRound3UsageDependencies|TestOpenAIBusyHintToolAdjacencyHTTP|TestCompletionTokenLimitHTTP|TestChapter15SentLensAndRestartIsolation)$' -count=1
go test -vet=off ./cmd/ch01 -run '^(TestLargeJSONLine|TestInputLines)$' -count=1
```
Last three all **PASS**; full command/output/exit receipts: `.ch15-work/r4-engine-focused-v3.log`, `r4-root-focused.log`, `r4-cli-focused.log`. `git --no-pager diff --check -- .` passed.

**Sequential compile-only:** root, GUI, ch06, ch07, examples/ch04-consumer, examples/ch13-speech each ran `.ch15-work/offline go test -vet=off -run '^$' ./...` from its module (absolute wrapper path for nested modules). All passed; receipts `.ch15-work/r4-compile-*.log`. No tests selected; no-test modules ch07/ch04-consumer are compilation evidence only. Required offline Go/cache settings, owned HOME/TMP, cleared/keyless environment and 1 GiB pre-command floor retained.

The original `.ch15-work/validation-r3/root-race.log` is unchanged. Full validation/race/vet/grades, frozen builds and later acceptance remain **NOT retried**. No browser/native/provider/external-network operations, grader/guide/manuscript inspection/edits, children, installation, commits/snapshots or shared cleanup. No unrelated production repair or price/model adjustment. **Source/tests held for the same retained reviewer's R4 delta inspection before any resumed validation.**

---

## Renewed R4 validation — local PASS; Chapter 15 launch failure; final HOLD

No source/test edits during renewed validation. Original `validation-r3` failure and all focused receipts remain unchanged. New receipt directory: `main/.ch15-work/validation-r4/`. Required cleared/keyless offline wrapper, existing caches, owned HOME/TMP and pre-command 1 GiB floor retained (initial free space 17,131,836 KiB).

**Sequential local gates, all exit 0**, via `.ch15-work/offline` from each module:
- Root: `go test -race ./...`, then `go vet ./...` — `root-race.log`, `root-vet.log`.
- GUI: `go test -race ./...`, then `go vet ./...` — `gui-race.log`, `gui-vet.log`.
- ch06, ch07, examples/ch04-consumer, examples/ch13-speech: each `go test ./...`, then `go vet ./...` — corresponding `ch06-*`, `ch07-*`, `ch04-consumer-*`, `ch13-speech-*` logs.
- ch07 and ch04-consumer reported **no test files**: test-command results are compilation evidence only, with separate successful vet. No separate redundant focused-JS/owner runs were added.

**Exactly one build of each frozen binary**, from the same held R4 source; both exit 0 (`build-cli.log`, `build-gui.log`):
```
# root cwd
go build -o /Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/validation-r4/frozen-cli ./cmd/ch01
# GUI cwd
go build -o /Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/validation-r4/frozen-gui ./cmd
```
Both used the offline wrapper. SHA-256:
- CLI: `876b309993c65098dabd4086887ac4e20561066c11671f41719bdede4669f8a5`
- GUI: `c42a6f1d9b35d6425edf2c88e42d9b633e0e32af7a983361011cd9faab3eac2d`

**Chapter 15 black-box gate**, book-root cwd `/Users/bill/projects/ensemble`, via the same wrapper:
```
go run ./cmd/grade-edition3 -ch 15 -json -cache-prices 3,3.75,0.3,15 /Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/validation-r4/frozen-gui
```
**FAILED: process exit 1; JSON reports `passed: false`, score 0/100.** All seven check reports share this launch failure:
```
GUI exited before ready: exit status 1: ... LLM_STREAM_CAPABILITY must be sse or plain
```
Thus the reported score does not establish failures in the cache behaviors themselves: the GUI did not reach ready. Cause is not classified or investigated. Full outputs/status/command preserved in `grade15.stdout.json`, `grade15.stderr.log`, and `grade15-run.log`. No grader source was inspected; only authorized black-box execution and command-location presence were used. Default artifact paths were used.

Stopped immediately after this failed gate. **No retry, diagnostic scope expansion, repair, binary rebuild, or Chapter 1 campaign.** The frozen CLI remains unused by the grader; genuine 62-second header/body lifetime behavior has not yet been exercised. Observation waits did not cancel any task. No additional production correction (still **4/10**), operations/browser/native/MIC/provider/discovery activity, children, installation, commits/snapshots or cleanup.

**Final HOLD** with source and both frozen binaries retained. Await root classification/release; renewed validation is not chapter or operational acceptance.

---

## Authorized fixture-only Chapter 15 retry — FAILED; final HOLD (still 4/10)

Root classified the prior invalid stream-capability launch literal as fixture-only and authorized one new black-box run. No student source/test changes, rebuilds, hash regeneration, local-gate reruns, or private fixture inspection occurred; the already verified frozen GUI/CLI were retained.

From book root `/Users/bill/projects/ensemble`, through the unchanged offline wrapper:
```
go run ./cmd/grade-edition3 -ch 15 -json -cache-prices 3,3.75,0.3,15 /Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/validation-r4/frozen-gui
```
**Process exit 1; JSON reports `passed: false`, score 0/100.** This time the GUI announced `http://127.0.0.1:56577`, but all seven check reports ended with **“WebSocket observation deadline”** and the GUI launch line. The grader completed itself; tool observation waiting did not cancel it. No attribution to product versus fixture has been made, and the cache-check behaviors were not established by these deadline reports.

New full receipts, preserving all earlier attempts:
- `.ch15-work/validation-r4/grade15-v2.stdout.json`
- `.ch15-work/validation-r4/grade15-v2.stderr.log`
- `.ch15-work/validation-r4/grade15-v2-run.log` (exact command/cwd/process exit)

Required offline/cache/owned HOME-TMP environment and pre-command 1 GiB floor remained in effect. Only authorized local fake-HTTP/GUI grader execution occurred. **Stopped at this failed gate: no further retry, diagnosis, repair, rebuild, cleanup or Chapter 1 campaign.** Existing local PASS receipts and frozen hashes remain valid/unchanged. Correction count remains **4/10**. Final HOLD awaits root classification and explicit release; no chapter/operational acceptance is claimed.

---

## Authorized Chapter 15 v3 — 75/100, FAILED; final HOLD (4/10)

After root classified/corrected the prior replay-observer fixture issue, ran exactly one authorized black-box attempt on the unchanged, already frozen GUI. No private grader inspection, source/test edit, rebuild, local-suite rerun or operations preparation.

Book-root cwd `/Users/bill/projects/ensemble`, unchanged offline wrapper:
```
go run ./cmd/grade-edition3 -ch 15 -json -cache-prices 3,3.75,0.3,15 /Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/validation-r4/frozen-gui
```
**Process exit 1; JSON `passed: false`, score 75/100.** Cache lens, stable prefix, request breakpoints, cache-read rate and all-vendor usage checks passed. Two reported failures:
- Cost: `response 2 cost <nil> priced=false, want 0.006045000000 at supplied rates [3 3.75 0.3 15]`.
- Session scope: `new session counts/cost include old work`.

These are grader observations, not independently diagnosed causes. Full new receipts under `.ch15-work/validation-r4/`: `grade15-v3.stdout.json`, `grade15-v3.stderr.log`, `grade15-v3-run.log`. All prior failures/receipts and both frozen binaries/hashes remain unchanged. Required offline settings, owned HOME/TMP and pre-command 1 GiB floor retained; authorized local fake-HTTP grader execution only.

Stopped at this failure with no diagnosis, fix, additional retry or cleanup. **Chapter 1 has still not run** because the Chapter 15 pass condition was not met. Source correction count remains **4/10**. Saved final HOLD awaits root classification/release; no chapter or operational acceptance claimed.

---

## Authorized Chapter 15 v4 + FIRST Chapter 1 — PASS; final saved HOLD (4/10)

Ran exactly the newly released Chapter 15 attempt on the existing frozen GUI after root's independently reviewed fixture-only producer-identity correction. No private fixture inspection, source/test edits, rebuilds, local-suite repeats or hash regeneration. Same held R4 source and previously verified frozen artifacts throughout.

Book-root cwd `/Users/bill/projects/ensemble`, through the unchanged offline wrapper:
```
go run ./cmd/grade-edition3 -ch 15 -json -cache-prices 3,3.75,0.3,15 /Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/validation-r4/frozen-gui
```
**PASS: process exit 0, independently parsed `passed: true`, 100/100; all seven checks passed.** Full new receipts: `.ch15-work/validation-r4/grade15-v4.stdout.json`, `grade15-v4.stderr.log`, `grade15-v4-run.log` (including JSON-check exit 0).

Only after that pass, ran the **first and only Chapter 1 campaign** on the original same-source frozen CLI:
```
go run ./cmd/grade-edition3 -ch 1 -json /Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/validation-r4/frozen-cli
```
**PASS: process exit 0, independently parsed `passed: true`, 100/100.** Five conversation rounds/requests, byte-preserved growing history, and cumulative usage input=511/output=65 passed. Both lifetime checks explicitly passed with **complete answer and usage after 1m2s headers delay and after 1m2s body delay**. Recorded process elapsed time: **124 seconds**. No outer deadline was imposed (allowance exceeded 360 seconds), and observation intervals did not cancel the job. Longer arbitrary deadlines remain a source-review matter, as the grader notes.

Full receipts: `.ch15-work/validation-r4/grade1-v1.stdout.json`, `grade1-v1.stderr.log`, `grade1-v1-run.log` (command/cwd, start/end timestamps, elapsed time, process and JSON-check exits).

All earlier failures remain preserved. Previously successful local gates and frozen hashes carry unchanged. Offline/cache/owned HOME-TMP/keyless environment and wrapper's 1 GiB pre-command floor remained in effect; only authorized local fake-HTTP grading ran. No production correction beyond **4/10**, external provider use, browser/native/MIC activity, children, installation, cleanup, commit/snapshot or operational preparation.

**FINAL HOLD:** all released validation gates have now passed, but this is not chapter or operational acceptance. Source and both frozen artifacts are retained unchanged for root's final disposition; operations remain closed.

---

## Chapter 15 operational proposal ONLY — HOLD, no execution authorized

Read current public §15.7 and current execution/student guidance; verified held GUI flags and provider-specific environment names read-only. Code/local validation is accepted at **4/10**, not the overall chapter. Source, frozen binaries and assets remain unchanged. No scripts or processes were created/launched for this proposal.

**Ownership/launch, proposed only.** Use `MAIN=/Users/bill/projects/ensemble/solutions/edition-3/main`, existing `$MAIN/.ch15-work/validation-r4/frozen-gui`, and its current `$MAIN/gui/web/gui` assets. Fresh owned directories `$MAIN/.ch15-work/ch15-ops/rehearsal` and `.../paid` keep fake and paid histories separate. Proposed loopback ports: relay **18765**, GUI **18766**, CDP **18767**—never 8088. Availability/ownership must be checked during separately approved preparation; collision means HOLD, not killing another listener.

From the chosen owned working directory, an ordinary launcher would keep GUI stdin open and run this existing binary (cleared environment, owned HOME/TMP/PATH; no inherited keys or generic `LLM_BASE_URL`):
```
env -i PATH=/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin \
HOME="$WORK/home" TMPDIR="$WORK/tmp" GOMAXPROCS=2 \
LLM_VENDOR=anthropic LLM_MODEL=claude-sonnet-4-6 \
LLM_STREAM_CAPABILITY=plain ANTHROPIC_API_KEY=local-dummy \
ANTHROPIC_BASE_URL=http://127.0.0.1:18765 \
EN_PRIMARY_SKILL=base EN_SKILLS_DIR="$WORK/skills" \
"$MAIN/.ch15-work/validation-r4/frozen-gui" --port 18766 \
  --gui-dir "$MAIN/gui/web/gui" --save "$WORK/save.json" \
  --settings "$WORK/cr/settings.json" --speaking=false
```
Use one owned Chrome/profile at a time, headless real rendering to avoid OS focus: proposed executable `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`, flags `--headless=new --mute-audio --no-first-run --no-default-browser-check --disable-background-networking --disable-component-update --disable-sync --disable-extensions --use-mock-keychain --remote-debugging-address=127.0.0.1 --remote-debugging-port=18767 --user-data-dir="$WORK/chrome-profile" --window-size=1440,1000`, opening only `http://127.0.0.1:18766`. Executable/tool attachment is not yet verified. Existing browser/CDP tools must bind explicitly to that owned target; if they cannot, HOLD rather than build another browser/controller stack. No MIC, native audio, OS focus or volume operations.

**Two approvals, not a campaign.** Preparation approval would permit a small owned launcher/relay and **one keyless local rehearsal**: frozen GUI/current assets, one fake turn, actual DOM/meter observation, minimal budget/model/native-cap rejection checks, then EOF and joins. No discovery or credential access in rehearsal. Its concrete implementation and guard review are the next separately approved phase, not delivered here. Paid approval would follow successful preparation/rehearsal, not automatically inherit from it.

**Proposed paid envelope.** Sonnet **4.6**, not the historical Sonnet-5 example. Exactly one approved direct Claude API-key read into the eventual relay's memory only; root must designate the precise allowed source. No search/fallback/login; no real key in children, environment, argv, files or logs. One bounded `GET https://api.anthropic.com/v1/models?limit=100` to confirm the chosen model; no pagination/discovery campaign or automatic model substitution. Forward only native `/v1/messages`, verbatim bodies, verified model and documented `max_tokens: 1024`; no body rewrite to manufacture a cap or cache hit. At most **four total attempted paid POSTs**, one armed per manual turn, with the same budget retained across GUI restart. No retries/redirect fallback/automatic tool continuation. Tool use, an extra POST, HTTP/guard/product failure or unverifiable identity/rate means stop and HOLD.

Proposed bounds: setup **8 minutes**, paid active window **15 minutes**, absolute owned run **25 minutes**; discovery **30 seconds**, each complete paid request **150 seconds**, each shutdown/join **30 seconds**. Observation expiry alone is not cancellation; these separately approved run limits are the stop boundaries. Maintain at least **1 GiB** free. No restarting the exercise to obtain a warm cache grant.

**Exact public workload.** `$WORK/skills/base/SKILL.md` contains this header/instruction text, followed immediately by exactly 700 numbered lines (i=1 through 700, one newline each), with no extra instructions:
```
---
name: base
description: Answer questions about the fixed archive reference.
tools: list_directory read_file
---
Answer the user's questions from this reference without using tools. Keep each answer to one short sentence.
Project codeword: COPPER-FERRY-18.
Rule 37: preserve the original source and record the reason for each change.
```
Each following line is exactly `Archive record {i}: color copper; preserve the original source and record the reason for each change.`
Submit these four public questions, sequentially, waiting for the real answer/usage before the next:
1. `Without tools, what is the project codeword?`
2. `Without tools, what does Rule 37 require?`
3. `Without tools, name the color mentioned in the archive records.`
4. `Without tools, repeat the project codeword in at most five words.`
Assess actual tool-free answers semantically: correct codeword, both Rule-37 obligations, copper, then correct short codeword answer. Do not demand incidental exact phrasing.

**Observe, restart, retain.** At initial zero and after each completed turn, inspect actual production DOM and screenshots—not injected counters—for the four values: session cost, latest all-input (disjoint buckets summed once), latest cache-read, session hit rate. Read the actual tooltip/title for last/total buckets, models, prices and prediction versus granted usage. Visually inspect captured screenshots of the visible meter; native title-tooltip rendering may not appear in headless captures, so do not claim visual tooltip evidence if only its DOM title is readable. Compare to real response usage, Engine usage frames and the sent-request lens, including marker placement and stable/shared prefix. Check the documented/approved Sonnet-4.6 standard per-million rates **3 / 3.75 / 0.3 / 15**, not rates inferred from discovery or invented aliases. Valid zero is zero; true unknowns remain unknown. A positive warm grant or arbitrary hit threshold is not a pass requirement.

After turn four: wait idle, capture evidence, send GUI stdin **EOF**, join it, and perform **exactly one GUI restart in the same paid directory** with the same save/settings/assets. Keep the relay budget and Chrome ownership intact. Before any new model request, inspect restored conversation and a new zero session meter; verify POST count stayed four. No fifth question. Then EOF/join again and stop only owned relay/Chrome PIDs/listeners; no shared cleanup or `killall`. Retain raw request/response bodies (never auth headers/secrets), saves, settings, usage/lens observations and screenshots locally under the owned directory, **not Git**; only explicitly approved owned-profile cleanup later.

**Needed next permission:** preparation/rehearsal only, including the named owned paths/listeners/browser target and minimal guard implementation. A subsequent explicit paid/key/discovery authorization is still required. No operations have begun; saved proposal HOLD awaits root review.

---

## One authorized keyless rehearsal — STOP / HOLD, no retry

Prepared one owned launcher using installed Node 24 built-in HTTP/WebSocket/CDP support, no dependency install or product changes:
- Script: `/private/tmp/ensemble-ch15-rehearsal.plRSMV/launcher.mjs`
- SHA-256: `206f7bc64f73f057ce0abf1f3846ed2a69dd1ed872cb8bec4722ed0831904002`
- Exact command: `/Users/bill/projects/ensemble/solutions/edition-3/main/.ch15-work/offline node /private/tmp/ensemble-ch15-rehearsal.plRSMV/launcher.mjs --rehearsal`

Used the **exact current public header including `type: primary`**, 700 records and four public question constants, without the proposal's extra Rule-37 line/rewording. Actual reference is `work/skills/base/SKILL.md` beneath that owned directory. Paid mode is hard-closed (`PAID_ENABLED=false`); real key/discovery code is unreachable in this invocation. No settings/key source was read and no external provider request occurred.

Startup short receipt recorded:
- Go-source/module hash (sorted relative `.go`, `go.mod`, `go.sum` paths + NUL + bytes; `.git`/`.ch15-work` excluded): `596f7a3f06665aae1c470f519cd076790b42a697cd3b9fbbc3df98da8d45d24e`
- Existing frozen GUI: `c42a6f1d9b35d6425edf2c88e42d9b633e0e32af7a983361011cd9faab3eac2d`
- Current GUI assets: `f18fd3d93fc64c3463cdfb44665905664a1900a578c4c5f0fe15bc605d2e2fd3`

Minimal local guard checks rejected wrong model, wrong native cap and over-budget requests before forwarding. Unused ports were checked; GUI PID 18986/executable/listener `127.0.0.1:18766` and Chrome PID 18994/executable/owned profile/listener `127.0.0.1:18767` were verified **before CDP attachment**. Relay/fake listeners were owned by launcher PID 18981 on loopback 18765/18768. Explicit cleared child environments/owned HOME-TMP/dummy GUI key and muted headless Chrome were used; no MIC, OS focus/volume or native-audio operation.

**Precise result: exit 1, `GUI DOM readiness deadline`, elapsed 121 seconds.** The attach/navigation sequence returned, but the GUI readiness predicate did not succeed before the 120-second setup deadline. No request/response, DOM receipt or screenshot files were produced; **no turn was sent and no image was captured or viewed**. Zero paid POSTs, provider discovery and key reads. This is not a successful rehearsal. Attach's target-fetch/open awaits still lack individual time bounds; no correction or further investigation was attempted pending root inspection.

The launcher performed owned cleanup: GUI stdin EOF, GUI and Chrome joined exit 0, listeners released. At root's stop instruction, the cleanup confirmation found those PIDs already exited, no group remnants and no listeners on the four owned ports; no forced signal was needed. Profile/save/settings/partial logs remain outside Git. No shared cleanup, retry, variant, source edit, rebuild, grade or further operational action.

Full local receipts: `run.log`, `receipt.log`, `gui.log`, `chrome.log`, `forced-stop.log` in `/private/tmp/ensemble-ch15-rehearsal.plRSMV/`. The last confirmed stage and ownership are preserved there; absent request/DOM/image receipts are explicitly listed in `forced-stop.log`.

**Saved HOLD for root's direct exact-script/secret-boundary/last-stage inspection. Paid/key/discovery remain CLOSED. Code/local validation remains accepted at 4/10, not overall chapter acceptance.**

### Final STOP/HOLD clarification — root owns controller next

Root confirms the receipt: **121 seconds, exit 1, GUI DOM readiness deadline; GUI/Chrome exit 0 and owned listeners released.** Cause remains **unknown**; an attach failure is **not proven**. Zero model POSTs/key reads remain the recorded relay boundary result. Root reports **Chrome-updater activity in the Chrome log**; no host-wide absence-of-side-effects claim is made.

Root has taken temporary ownership of the operational controller. Student will not inspect further, retry, edit the controller, or perform additional operational work without a new task. Failed script/receipts are preserved; production and frozen artifacts remain held. **STOP/HOLD.**
