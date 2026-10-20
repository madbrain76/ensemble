# Chapter 10 restart — Round 1 HOLD

Formal correction **1 of 10** complete; source held for root and independent review. Initial implementation was from accepted Chapter 9 at assigned launch HEAD `06f6b9e2`. Initial compile-only HOLD used no correction round. This round repaired only the six released public-contract gaps. Focused validation and compilation are not chapter acceptance.

Mandatory public guidance/chapter reads and permitted Chapter 9 source were used; initially truncated reads were recovered. No private reviewer/grader source, original answer, future chapter, children, commits, snapshots, installation, restart, shared cleanup, browser, or external provider was used.

## Ownership and Round 1 changes

- **History** (`internal/common/persistence.go`, `internal/history/persistence.go`): removed the unsolicited top-level version field; retained strict decoding/anchor checks. Save creates missing parent directories only after encoding/validation and recovery-candidate checks. Dialogue-only projection includes visible human hints with original sequence/role, preserving ordering and repeated-save deduplication. Backup/replay algorithms otherwise unchanged.
- **Diagnostics** (`internal/diagnostics/log.go`, root façade `diagnostics.go`): Write only appends complete records; serialized `Retain` independently applies the current limit. Reopening on append avoids stale offsets after replacement. Zero remains unlimited; negative retention is rejected.
- **CLI** (`internal/cli/cli.go`): ordinary flag parsing accepts `--save PATH`; render/dump remain offline and do not touch the selected conversation archive. Retention occurs after valid recovery or successful fresh save, and on normal exit only after successful shutdown/save. Failed load/save does not trigger destructive retention.
- **GUI lifecycle** (`gui/terminal.go`, `gui/cmd/main.go`): terminal input uses existing reliable `SubmitLive` result channels, keeps accepting live hints, and drains admitted terminal work at EOF before application shutdown joins owned jobs and saves. Unbounded line reads remain. Diagnostic retention uses the same safe startup/normal-exit boundaries.
- **GUI projection** (`gui/identity.go`): retained lifecycle events bind to the existing context-prefix call identity despite the prefix lacking TurnID. Results and late job links reuse that association, including the display-only nested job link. Canonical/native IDs are unchanged.
- **Focused tests**: new `ch10_round1_test.go`, `internal/cli/ch10_round1_test.go`, `gui/ch10_round1_test.go`, `gui/cmd/ch10_round1_test.go`; updated `internal/history/persistence_test.go` and `diagnostics_test.go`.

The existing compact History/Agent/GUI ownership remains; root is wiring. Initial Chapter 10 changes outside these Round 1 paths were not redesigned. This note is the only owned document.

## Environment and outputs

`out=/tmp/ensemble-ch10-round1.lLZ73v`; HOME=`$out/home`; TMPDIR/TMP/TEMP=`$out/tmp`.

```
GOMAXPROCS=2
GOFLAGS='-mod=readonly -p=1'
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
GOMODCACHE=/Users/bill/go/pkg/mod
GOCACHE=/Users/bill/Library/Caches/go-build
PATH=/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin:/usr/sbin:/sbin
```

Unset `LLM_BASE_URL LLM_API_KEY LLM_KEY_FILE ANTHROPIC_API_KEY OPENAI_API_KEY GEMINI_API_KEY GOOGLE_API_KEY`. No real key reads. HTTP tests used loopback fixtures and explicit fake keys only. Runs were sequential. Last disk check: **4,697,916 KiB available**, above the 1 GiB floor. Logs and build outputs are private temporary files, not repository artifacts.

## Before-fix evidence and setup failures

Commands below use `go test -vet=off -count=1` from main unless stated.

- `-run '^TestCh10R1' . ./internal/history ./internal/cli` → exit 1, `before-main.log`: actual failures for unsolicited version, omitted human hints, missing default Agent directories, and append-time destructive retention. CLI fixture initially had the wrong Run signature (not behavioral evidence).
- Corrected CLI fixture: `-run '^TestCh10R1' ./internal/cli` → exit 1, `before-cli.log`: `--save PATH` rejected by actual CLI parsing.
- From gui: `-run '^TestCh10R1' . ./cmd` → exit 1, `before-gui.log`: actual GUI-main subprocess with closed stdin lost the delayed last prompt/effect/final saved answer. The GUI package initially had a fixture field compile error; separate prefix fixture attempts then lacked provenance/used a wrong release tag (`before-prefix.log`, `before-prefix-valid.log`), not product evidence.
- Valid prefix fixture, from gui: `-run '^TestCh10R1SnapshotPrefixTailWS$' .` → exit 1, `before-prefix-valid2.log`: prefix ID 7 split into fallback 9007199254740991 and result lost its prefix association.
- During repair, `after-gui-1.log` records a completion-channel API wiring compile error (fixed to match existing CLI usage). Expanded fixtures had a typed-vendor compile error, incorrect job wire tag/card expectation (`after-gui-expanded-1.log`, `after-gui-expanded-2.log`, `after-prefix-expanded.log`). `final-offline-cli.log` records a missing dump-input fixture file; providing that required history file fixed setup. These are not credited as behavioral failures of production. No dependency setup failure or unexpected outside-scope product failure remains.

## Passing focused validation

Exact final commands (same environment):

| Cwd | Command after `go test -vet=off -count=1` | Exit | Log under `$out` |
| --- | --- | --- | --- |
| main | `-run '^(TestCh10R1\|TestArchive\|TestDialogueArchive\|TestBadArchive\|TestDiagnosticRetention)' . ./internal/history ./internal/cli` | 0 | `final-main.log` |
| main/gui | `-run '^(TestCh10R1\|TestLegacyDisplayIDs\|TestCh7R1Terminal\|TestCh7R2TerminalInterrupt\|TestSettingsPatchValidation\|TestSettingsPersistenceAndFailedSave)' . ./cmd` | 0 | `final-gui.log` |
| main | `-run '^TestCh10R1CLIReadOnlyModes$' ./internal/cli` | 0 | `final-offline-cli-valid.log` |

Table pipes are Markdown-escaped; the actual regex alternation characters are `|` (no backslashes).

Coverage includes archive backup/replay/refusal, repeated dialogue-only hint saves, two fresh library Agents with isolated defaults, actual CLI override/read-only modes, quiet retention, failed load/exit-save preservation, post-retention append offset, zero/negative settings, and legacy IDs. The GUI-main subprocess uses real stdin and delayed loopback HTTP responses, executes a local write_file effect, and verifies the final answer and completed saved history before exit. The History/Agent/hub/WebSocket regression covers snapshot+tail, context-only+live tail, full-log carry, late jobs, new display IDs, reconnect, and save/reload without rewriting native IDs. The withdrawn general full-log-double-seeding hypothesis was not reinstated.

## Compilation and HOLD

From each explicit module `.`, `gui`, `ch06`, `ch07`, `examples/ch04-consumer`, with `name=$(printf '%s' "$module" | tr '/.' '__')`:

```
go build -o "$out/compile-$name/build/" ./...
go test -vet=off -c -o "$out/compile-$name/test/" ./...
```

All five builds and test-source compilations **exit 0**. Also **exit 0**: main `go build -o "$out/ch01" ./cmd/ch01`; main/gui `go build -o "$out/gui-command" ./cmd`. Full command/output log: `$out/final-compile.log`. `-c` compiles test sources without executing suites. Scoped `git diff --check` passed. Earlier initial compile-only evidence remains at `/tmp/ensemble-ch10-restart-compile.KAGLW0`.

Source is held after these checks. No broad suite, race/vet campaign, grader, browser/native UI exercise, real provider validation, or mutation campaign ran. Remaining review is source/reference correctness and broader inherited behavior not covered by these focused cases; compilation does not establish acceptance or universal filesystem durability. Backup recovery still does not promise newer bytes lost with a corrupt primary.

## Full-validation release — STOP at first gate

Root and the same reviewer closed the six R1 findings for validation, not acceptance. Production and fixtures were held unchanged. New owned output directory: `/tmp/ensemble-ch10-validation.LbEDm2` (HOME=`home`, TMPDIR/TMP/TEMP=`tmp` beneath it). Preflight available disk: 4,655,988 KiB; `/tmp/e3-ch10-final/frozen-cli` did not exist.

From `/Users/bill/projects/ensemble/solutions/edition-3/main`, ran exactly:

```
env -i HOME=/tmp/ensemble-ch10-validation.LbEDm2/home TMPDIR=/tmp/ensemble-ch10-validation.LbEDm2/tmp TMP=/tmp/ensemble-ch10-validation.LbEDm2/tmp TEMP=/tmp/ensemble-ch10-validation.LbEDm2/tmp PATH=/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin:/usr/sbin:/sbin GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod GOCACHE=/Users/bill/Library/Caches/go-build go test -race ./...
```

**Exit 1.** Full output preserved at `/tmp/ensemble-ch10-validation.LbEDm2/main-race.log`. Reported failures include inherited Chapter 4/5 job tests rejected by save validation (`missing tool call ID`), tests rejected for overlapping recovery paths, CLI expectations affected by loading `save.json`, and a nil-pointer panic at `internal/agent/persistence.go:19` called from shutdown in `actor.go:391`. These observations are not yet a diagnosis or a race finding. No automatic production/fixture patch or cleanup was attempted; any test-created runtime files are left intact pending direction.

Stopped at this first gate. Root vet, GUI race/vet, separate-module gates, frozen binary build/SHA, model map, Chapter 10 grader, and Chapter 1 grader were **not run**. No operational proposal is submitted because the required gates have not passed. Awaiting explicit triage/correction authority; source remains HOLD.

### Authorized bounded read-only triage

Parent authorized inspection of named failing tests and their owner call paths only; no edits, reruns, cleanup, or later gates. Classification is based on actual composition, not test age:

1. **Production contract defect: valid direct dispatch cannot be saved.** `internal/tool/dispatch.go:21–22` implements public Dispatch using a Part with name/args and no native CallID. Await posts ToolAdmission (same file:78 onward); `internal/agent/actor.go:323–345,639–688` admits it and records called/released/returned events without inventing a CallID. New `internal/history/persistence.go:375–378` rejects every Tool event with an empty CallID, so default shutdown Save rejects History produced by this valid public API. `chapter04_jobs_test.go:18–37` and `c4Dispatch` in `chapter04_round1_test.go` exercise that path. This is not grounds to suppress shutdown errors or dismiss the tests.
2. **Production default-path isolation defect, also causing history contamination.** Root `ensemble.go:110–113` forms `DataDir/agent-N/save.json` before Agent normalization. Agent normalizes BaseDir but leaves the generated relative SavePath process-relative (`internal/agent/agent.go:107–127`). By contrast, jobs resolve relative DataDir beneath BaseDir (`internal/jobs/jobs.go:117–123`). Distinct `c4Agent` roots use distinct temporary BaseDirs but the same literal `DataDir:"cr"`, each starting at agent-1, so they load/save the same process-relative archive. The unexpected `cwd-override` event in `TestChapter5SupervisionLogSchema` matches a different test's explicit Invoke ID (`chapter04_round1_test.go:395–398`). The isolation defect is established from source; the exact pre-run/in-run origin of every contaminated event was not traced.
3. **Fixture contract mismatch: copying an effective recovery path into a second live Agent.** `chapter04_round1_test.go:550–556` and `chapter04_round2_test.go:216–221` copy `a.Config()`, change only provider fields, then call `app.NewAgent`. Config now includes the first Agent's concrete SavePath. Root's distinct-primary/backup guard correctly rejects that shared path. These fixtures need an intentionally distinct recovery destination (or new default), not removal of the guard.
4. **CLI fixture isolation mismatch.** `cmd/ch01/main_test.go` invokes run repeatedly without a temporary cwd or distinct save override and requires empty diagnostics/exact fresh turn counts (lines 87–90,151–152). Current default persistence deliberately uses `save.json` and reports recovered boundaries (`internal/cli/cli.go:56–57,106–107`). Earlier test invocations therefore become later invocations' history. The observed recovery diagnostics and accumulated turn counts are expected for that shared destination; preserve the behavioral assertions but provide isolated state when correction is authorized.
5. **Internal fixture composition mismatch behind the panic.** `internal/agent/actor_test.go:92–103` manually constructs and starts Agent with fake Engine/History/Tools/Jobs but never SetSkills. Default SavePath is now populated; shutdown calls Save (`actor.go:391`), which dereferences missing skills at `persistence.go:19`. The Engine/History fixtures also embed nil interfaces without implementing EffectiveSystemPrompt/SaveFile. Public `Ensemble.NewAgent` attaches actual Engine, History, Tools, Skills and loads skills before Start. Thus this crash establishes an incomplete internal-owner fixture, not a demonstrated panic through valid public root composition; no nil-guard production patch was made.

No data-race report was identified in the failed output, but the failed/panicking race gate is not a race pass. Read-only inspection included a few nonexistent filename guesses; they changed nothing and are not product failures. Source/fixtures remain HOLD, output remains preserved, and explicit correction authority is required before any repair or rerun. No operational proposal or grader result is available.

## Correction Round 2 of 10 — HOLD

Assistance/scope: root and retained reviewer supplied the stopped-gate classification and explicitly authorized two production repairs plus the listed fixture adaptations. No children, private review/grader inspection, broad validation, or operational work. Source is held after focused tests and compilation; this is not chapter acceptance or release of the interrupted full-validation campaign.

**Production (two files):** `internal/history/persistence.go` permits an absent native CallID only for turn-zero Tool events, as produced by real public Dispatch; model-turn events still require a CallID, and tag/payload/argument/sequence validation remains. Nothing invents historical IDs or suppresses save errors. `ensemble.go` derives default library recovery paths from BaseDir-relative DataDir, resolves the derived path absolutely, and retains per-Agent directories. Explicit SavePath semantics and overlap refusal are unchanged; fixtures continue to use relative `DataDir:"cr"`.

**Approved fixture adaptations:** `chapter04_round1_test.go` separates intentionally broken job storage from valid scratch recovery storage and gives the second Agent a distinct SavePath. `chapter04_round2_test.go` clears the copied effective SavePath so the second Agent gets its own default. `cmd/ch01/main_test.go` gives each intended fresh session an `ENSEMBLE_SAVE` scratch override, retaining its input/error/output assertions. `internal/agent/actor_test.go` explicitly excludes persistence from its handcrafted actor-ordering seam; stale-reply/job-fact assertions remain, and no production nil-skip was added.

**Focused regressions:** new `ch10_round2_test.go` exercises real Dispatch (ordinary write, control tool, failed read), shutdown/save/reopen, actual effect, supplied native-ID correlation, and exact canonical event equality after restore. It also checks two roots with distinct absolute/relative BaseDirs and relative DataDir, independent save/reload, no launch-directory archive, explicit paths and overlap refusal. `internal/history/persistence_test.go` adds rejection checks for missing model-turn call ID, non-object args, mismatched event payload/tag, and invalid sequence anchor. Existing backup/full/dialogue/refusal checks were retained.

### Commands and results

All commands sequential, clean/offline. `out=/tmp/ensemble-ch10-round2.LyoGdf`; same environment as the stopped gate, but run with `env -i`, HOME=`$out/home`, TMPDIR/TMP/TEMP=`$out/tmp`. PATH `/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin:/usr/sbin:/sbin`; GOMAXPROCS=2, GOFLAGS=`-mod=readonly -p=1`, GOTOOLCHAIN=local, GOPROXY=off, GOSUMDB=off, GOMODCACHE=`/Users/bill/go/pkg/mod`, GOCACHE=`/Users/bill/Library/Caches/go-build`. No inherited keys or generic LLM_BASE_URL. Last pre-compile disk check: **4,519,288 KiB** available. Previous stopped evidence/runtime directory `/tmp/ensemble-ch10-validation.LbEDm2` was left intact; no shared cleanup.

From main, each following command uses `go test -vet=off -count=1`:

- Before production repair: `-run '^TestCh10R2' .` → **1**, `$out/before-valid.log`: valid direct dispatch failed shutdown save; derived default escaped BaseDir. Existing stopped-race evidence was reused rather than repeating the campaign.
- `-run '^TestCh10R[12]' .` → **0**, `$out/final-root.log` (includes final absolute/relative BaseDir check).
- `-run '^(TestCh10R[12]|TestArchive|TestDialogueArchive|TestBadArchive)' ./internal/history` → **0**, `$out/after-history.log`.
- `-run '^TestQueuedStaleResponseIsNotPublishedInNewTurn$' ./internal/agent` → **0**, `$out/after-actor.log`.
- `-run '^(TestLargeJSONLine|TestInputLines|TestChatAndConfiguration)$' ./cmd/ch01` → **0**, `$out/after-cli-final.log`.
- `-run '^(TestChapter4InitialJobSmoke|TestChapter4Round1(KilledMetadata|OutputWakeups|ConsumedLimits|FileStreaming|SetupFailureCorrelation)|TestChapter4Round2(TerminalSettingsAndParity|LineInputAndTerminalCursor|StagedCRLFAndUnseenPattern|AgentWireAndKilledObservation)|TestChapter5(WakeInitialAndLaterObservation|ExplicitGoKillReturnsBeforeHandler|ShutdownJoinsKilledFactDelivery|SupervisionLogSchema)|TestCh7R2ControlPrepareHasNoEarlyEffect)$' .` → **0**, `$out/after-inherited.log`. This selects the inherited failures named in the stopped gate, not a full historical suite.

Setup failures are retained, not counted as product evidence: `$out/before.log` used the wrong Dispatch return type; `after-root.log`/`after-root-valid.log` had a fixture redeclaration compile error (including a failed text-edit attempt); `after-root-fixed.log` exposed an incorrect new assertion that Invoke retains caller display IDs. Existing actor admission explicitly allocates display IDs, so the fixture now checks recorded-link consistency and exact save/reopen equality while preserving the supplied native ID. `after-cli.log`/`after-cli-fixed.log` record a missing filepath import and failed import-edit attempt. All corrected; no unexpected outside-scope product failure was repaired.

From each of the five explicit modules `.`, `gui`, `ch06`, `ch07`, `examples/ch04-consumer`, with `name=$(printf '%s' "$module" | tr '/.' '__')`, ran `go build -o "$out/compile-$name/build/" ./...` then `go test -vet=off -c -o "$out/compile-$name/test/" ./...`: **all exit 0**. Exact expanded commands/output: `$out/compile.log`. These are compilation checks only; ch07 and consumer report no test files. Scoped `git diff --check` passed.

No full race/vet/grader/provider/browser/key work, commits, snapshots or cleanup. Full validation remains paused for root/retained-reviewer inspection of the actual R2 delta. No frozen CLI/SHA, model map, grader result, or operational proposal was produced.

## Resumed full validation after R2 review — STOP

Root reported retained-reviewer R2 PASS for validation, with R1 still closed. Held production and fixtures were not edited. New outputs: `/tmp/ensemble-ch10-validation-r2.xG4Y2Q`; the first failed campaign remains untouched. Disk preflight: **4,309,760 KiB** available; `/tmp/e3-ch10-final/frozen-cli` absent.

From `/Users/bill/projects/ensemble/solutions/edition-3/main`, executed exactly:

```
env -i HOME=/tmp/ensemble-ch10-validation-r2.xG4Y2Q/home TMPDIR=/tmp/ensemble-ch10-validation-r2.xG4Y2Q/tmp TMP=/tmp/ensemble-ch10-validation-r2.xG4Y2Q/tmp TEMP=/tmp/ensemble-ch10-validation-r2.xG4Y2Q/tmp PATH=/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin:/usr/sbin:/sbin GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod GOCACHE=/Users/bill/Library/Caches/go-build go test -race ./...
```

**Exit 1**, preserved output `/tmp/ensemble-ch10-validation-r2.xG4Y2Q/main-race.log`. The command's default vet phase rejected `ch10_round2_test.go:30:41`: Fatalf `%s` was given a `PartList`. Root package therefore did not run. This is a new fixture diagnostic-format mistake missed by compile-only/focused `-vet=off` checks, not an observed runtime/race defect. Other package results printed by that single command do not make the gate pass.

Stopped immediately after this command: no automatic patch, rerun, cleanup, standalone vet, GUI/separate-module gate, frozen CLI/SHA, map, or grader. Assistance requested from parent for an explicit bounded fixture-format correction; source remains HOLD and full validation remains incomplete. No operational proposal is warranted yet.

## Full-validation continuation — local gates PASS, Chapter 10 FAIL / HOLD

Parent verified the previous diagnostic-only error and authorized **only** `ch10_round2_test.go` Fatalf PartList formatting `%s`→`%v`; no assertion/production change and no additional correction round. Prior focused runs had `-vet=off`, explaining why they missed it. Preserved `main-race.log`; resumed with new output. No other source/fixture edit during validation.

Same clean `env -i` prefix and owned directory `/tmp/ensemble-ch10-validation-r2.xG4Y2Q` as immediately above. All following commands completed sequentially with **exit 0**:

- main: `go test -race ./...` → `main-race-format-fixed.log`; then `go vet ./...` → `main-vet.log`.
- main/gui: `go test -race ./...` → `gui-race.log`; then `go vet ./...` → `gui-vet.log`.
- main/ch06: `go test ./...` → `ch06-test.log`; `go vet ./...` → `ch06-vet.log`.
- main/ch07: `go test ./...` → `ch07-test.log`; `go vet ./...` → `ch07-vet.log`. **No test files: compile-only plus vet**, not behavioral testing.
- main/examples/ch04-consumer: `go test ./...` → `examples_ch04-consumer-test.log`; `go vet ./...` → `examples_ch04-consumer-vet.log`. **No test files: compile-only plus vet**.

These are the released full gates, not carried-over focused results. Go reused valid unchanged package cache results where reported, as parent instructed. Before build/grader, disk available **4,260,704 KiB** and frozen CLI absent. From main, one `go build -o /tmp/e3-ch10-final/frozen-cli ./cmd/ch01` succeeded (exit 0), output `frozen-build.log`. `shasum -a 256 /tmp/e3-ch10-final/frozen-cli` recorded `35b93169bec3b4487404101840e9f184d2bc475ef3befe01697c5cd00f826f25` in `frozen-cli.sha256.txt`. No replacement/rebuild after freezing.

Owned map `/tmp/ensemble-ch10-validation-r2.xG4Y2Q/ownedmap.json` contains exactly `{"anthropic":"claude-sonnet-4-6","openai":"gpt-5.4","gemini":"gemini-3.8-flash"}`. From repository root `/Users/bill/projects/ensemble`, with the same clean environment, executed:

```
go run ./cmd/grade-edition3 -ch 10 -json -core /Users/bill/projects/ensemble/solutions/edition-3/main -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui -gui-model claude-sonnet-4-6 -media-models /tmp/ensemble-ch10-validation-r2.xG4Y2Q/ownedmap.json /tmp/e3-ch10-final/frozen-cli
```

**Process exit 1 and JSON `passed:false`.** Full stdout: `grade-ch10.json`; stderr: `grade-ch10.stderr.log` (both under the owned validation directory). Passed checks: save shape, replay equals snapshot, tail applied once, snapshot without log. Failed checks/report details:

- `default-load`: core/GUI report historical model/vendor/prompt replaced runtime configuration.
- `bad-save-refused`: core/GUI rejected malformed saves without requests and preserved bytes, but report says diagnostic `named=false`.
- `ch9-parity`: nested streaming delta/final mismatch and missing final text; GUI log timestamp/direction/JSON format and direction failures; inherited replay-is-live exits on `validating save: context disagrees with canonical boundary`.

These are black-box observations, not yet diagnosed. No private grader inspection, automatic patch, cleanup or retry. **Chapter 1 did not run** because Chapter 10 failed the required JSON/process condition. Frozen binary, map, failed outputs and earlier campaigns remain intact. Parent asked for next authority; source remains HOLD. No operational proposal, browser/native/model exercise, key access or discovery occurred.

### Authorized read-only triage of Chapter 10 failure

Parent authorized bounded public-code/public-contract inspection only. No edits, tests, retries, private grader inspection, frozen-binary replacement or cleanup followed the grader failure.

- **Corrupt-save naming: supported production diagnostic gap.** Chapter 10 rule 2 requires naming candidate failures. `internal/history/persistence.go:LoadFile` reports generic primary/backup roles plus decoder errors, without explicitly wrapping each with the selected filename. Decoder/schema errors themselves do not carry a path. Preserve refusal/no-install/no-overwrite behavior; a future correction should add path context, not weaken validation.
- **GUI trace format: supported production contract mismatch.** Chapter 7 §trace example/body explicitly says RFC3339 timestamp, `>`/`<`, then one compact JSON payload; the whole line is not a JSON object. `gui/server.go:wireTrace` instead writes a JSON object containing time/direction/payload-string. Retention can remain writer-owned; its current envelope violates inherited wire-log format.
- **Replay boundary: supported representational comparison defect, not permission to remove strict comparison.** `internal/history/history.go:Apply` can leave a non-nil empty `Pending` slice after the final tool result. `internal/common/history.go:Context.Pending` has `json:"pending,omitempty"`, so the snapshot round-trip omits it and decodes nil. `internal/history/persistence.go:sameContext` normalizes only empty Ephemera, then calls reflect.DeepEqual. A complete canonical prefix ending during that state can therefore disagree solely on nil versus empty Pending. This is a concrete source-level explanation compatible with the replay-is-live boundary report, but no reproduction was run during read-only triage. RawMessage byte comparison is also representation-sensitive; do not claim it caused this report without a focused case.
- **Config finding remains unresolved; do not infer restoration from grader wording.** Root NewAgent validates current config, constructs current owners, then Agent.Load installs History/context and seeds IDs. Agent.Load does not assign SaveFile.Config back to Agent config. Engine.Render and EffectiveSystemPrompt read current Agent Config and current skills snapshot. Inspected GUI settings apply only current settings policies, not saved archive model/vendor/prompt. Saved metadata and effective prompt composition still need a focused captured-request comparison under a correction release; no historical-config restoration assignment was found in the inspected path.
- **Streaming delta/final report remains unresolved.** Read the current GUI projection/identity and immutable-checkpoint sync paths; the identity map explicitly hands a partial ID to its final. Static inspection did not establish the report's cause. No general full-log double-seeding claim is revived. A future bounded real-path regression is needed before changing identity or streaming machinery.

State: local released gates passed; Chapter 10 black-box failed; Chapter 1 remains unrun. Frozen binary and all report/evidence files remain preserved. **Blocked HOLD pending root/retained-reviewer classification and explicit correction authority.**

## Correction Round 3 of 10 — focused PASS / compiling HOLD

Root + retained reviewer released exactly three production areas and retracted any blanket final-ID reallocation claim. Current prompt composition and literal error-filename grader alignment remain root-owned; no production prompt/config/path-message change was made. No private grader inspection or children.

**Production changes (existing owners only):**
- `gui/server.go`: restores Chapter 7 trace lines: RFC3339Nano timestamp, direction, compact JSON. One serialized complete-record write remains; malformed inbound bytes are represented as a readable JSON string rather than extra lines. Safe-boundary retention is unchanged.
- `internal/history/persistence.go`: comparison serializes temporary Context values and decodes JSON values with `UseNumber`, normalizing required empty collections and existing `omitempty` fields. Object key order/RawMessage whitespace and empty Pending allocation no longer create false disagreements; large integer distinctions are preserved. Canonical History/context bytes are not rewritten for comparison; strict validation remains in place.
- `gui/hub.go`: exact streaming cause reproduced: a coalesced subscriber wake can see the final after the actor has removed its last partial. IDs were already stable. Hub now emits the unseen ordinary-text suffix from the authoritative final before that final under the same ID. No identity-map redesign or canonical/provider-ID modification.

**Regression files:** new root and GUI `ch10_round3_test.go`; additions in `internal/history/persistence_test.go`. Actual local-HTTP tool dispatch → JSONL-imported complete/prefix histories → save/reopen → equal rendered requests and canonical event data; content/args/provenance/anchor/sequence negatives; semantic raw-JSON and large-integer negative. Actual WS trace traffic, compaction, both directions, retained whole lines and append offset. Actual restored historical Agent + two fresh SSE responses + hub/WS, with subscriber delayed across the last chunk: concatenated deltas equal final under the same agent/part key; distinct new IDs and untouched historical canonical events checked.

**Environment/output:** `/tmp/ensemble-ch10-round3.NI8aQa` (`out` below), new owned HOME/TMP, no inherited keys or generic LLM_BASE_URL. Every Go command used:
```
env -i HOME="$out/home" TMPDIR="$out/tmp" TMP="$out/tmp" TEMP="$out/tmp" PATH=/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin:/usr/sbin:/sbin GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod GOCACHE=/Users/bill/Library/Caches/go-build
```
Before production changes (each test command **exit 1**):
- main: `go test -vet=off -count=1 -run '^TestCh10R3' . ./internal/history` → `before-main.log`: actual tool-result-prefix Save rejected its valid boundary; semantic JSON/empty-collection comparison failed. Completed-turn case passed.
- gui: `go test -vet=off -count=1 -run '^TestCh10R3' .` → `before-gui.log`: trace JSON envelope rejected; both fresh responses had stable IDs (5 then 8), but deltas `head-N-` lacked final suffix `tail`.

After repairs (**exit 0**):
- main `go test -vet=off -count=1 -run '^TestCh10R[123]' .` → `after-root.log`.
- main `go test -vet=off -count=1 -run '^(TestCh10R|TestArchive|TestDialogueArchive|TestBadArchive)' ./internal/history` → `after-history.log`.
- gui `go test -vet=off -count=1 -run '^(TestCh10R|TestLegacyDisplayIDs|TestCh7R4EmptyStartLiveAndReplay|TestCh7R1Terminal|TestCh7R2TerminalInterrupt)' .` → `after-gui-fixed-ack.log`. Includes the earlier legacy/collision/context-prefix-tail carry cases and existing empty-start/terminal checks.

Setup/fixture failures, not product evidence: one refused hub text anchor caused no hub edit; fixed against the actual block. Trace fixture initially mistook queued startup settings for an acknowledgement, and `subscribe` has no response (`after-gui.log`, exit 1). One attempted edit used the wrong cwd and did not apply; unchanged fixture rerun failed (`after-gui-ack.log`, exit 1). Corrected fixture consumes startup frames and uses acknowledged `update_settings` in a scratch settings file. No unrelated product failure occurred.

**Compile-only:** disk before compilation **4,047,868 KiB** free. Sequentially for `module=. gui ch06 ch07 examples/ch04-consumer`, from each `$base/$module` (`base=/Users/bill/projects/ensemble/solutions/edition-3/main`), with `name=$(printf '%s' "$module" | tr '/.' '__')`:
```
go build -o "$out/compile-$name/build/" ./...
go test -vet=off -c -o "$out/compile-$name/tests/" ./...
```
All ten commands **exit 0**. Summary `compile.log`; per-module `compile-$name/build.log` and `test-compile.log`; outputs under those owned directories. Test compilation is not execution. `git diff --check` clean.

**HOLD:** no source edits after successful compilation. No full race/vet/grader run in R3, no cleanup, no prompt/config changes, no operational work. Earlier failure campaigns and `/tmp/e3-ch10-final/frozen-cli` remain untouched; no new frozen binary. Await root/retained-reviewer inspection and explicit validation release; these focused/compile results are not chapter acceptance.

## R3 validation release — all authorized gates PASS; operational HOLD

Root and retained reviewer accepted R3 for validation and authorized only changed-owner gates plus grading. Source/fixtures stayed held; no edits or retries. Fresh output directory: `/tmp/ensemble-ch10-validation-r3.5VP9vE` (`out`). Disk checks: **3,655,208 KiB** initially, **3,630,660 KiB** before freeze, **3,588,768 KiB** before Chapter 1 (all above 1 GiB).

All Go commands used clean `env -i` with `HOME=$out/home`, `TMPDIR/TMP/TEMP=$out/tmp`, `PATH=/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin:/usr/sbin:/sbin`, `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, `GOMODCACHE=/Users/bill/go/pkg/mod`, `GOCACHE=/Users/bill/Library/Caches/go-build`. No inherited real keys or generic LLM_BASE_URL.

**Sequential changed-owner gates, all exit 0:**
- main `go test -race ./...` → `main-race.log`; `go vet ./...` → `main-vet.log`.
- main/gui `go test -race ./...` → `gui-race.log`; `go vet ./...` → `gui-vet.log`.
- As explicitly directed, carry unchanged ch06 tests/vet and ch07/consumer compile-only (no tests)/vet from `/tmp/ensemble-ch10-validation-r2.xG4Y2Q/{ch06,ch07,examples_ch04-consumer}-{test,vet}.log`. No separate reruns of those modules.

**New frozen CLI:** from main, `go build -o /tmp/e3-ch10-final/frozen-cli-r3 ./cmd/ch01` → `frozen-build.log`, **exit 0**, built once after verifying the name absent. `shasum -a 256 /tmp/e3-ch10-final/frozen-cli-r3` → `frozen-cli-r3.sha256.txt`:
`e148230ed0edbbf1d921fd0fcbe3680ff40cb5f41a9d293ed653c9d08610d980`.
Old `/tmp/e3-ch10-final/frozen-cli` and prior failure campaigns untouched.

From repository root `/Users/bill/projects/ensemble`, same clean environment, executed exactly:
```
go run ./cmd/grade-edition3 -ch 10 -json -core /Users/bill/projects/ensemble/solutions/edition-3/main -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui -gui-model claude-sonnet-4-6 -media-models /tmp/ensemble-ch10-validation-r2.xG4Y2Q/ownedmap.json /tmp/e3-ch10-final/frozen-cli-r3
```
Reused map contains anthropic=`claude-sonnet-4-6`, openai=`gpt-5.4`, gemini=`gemini-3.8-flash`. **Process exit 0; parsed JSON `passed:true`; 100/100, all checks pass including inherited parity.** Stdout `grade-ch10.json`; stderr `grade-ch10.stderr.log`. JSON parsed with Python stdlib and required `d.get("passed") is True`; did not rely on title/numeric chapter.

Only after that pass, from repository root with the same environment and **same frozen executable**:
```
go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch10-final/frozen-cli-r3
```
**Process exit 0; parsed JSON `passed:true`; 100/100.** Stdout `grade-ch1.json`, stderr `grade-ch1.stderr.log`, elapsed/exit `grade-ch1-timing.log`. Actual elapsed **125 seconds**. Report explicitly confirms complete answer and usage after **1m2s delayed headers** and separately **1m2s delayed body**; normal five-round usage totals input=511/output=65 match. No outer cancellation deadline was imposed (allowance exceeds 360s); observation waits were not deadlines. These observations do not claim universal arbitrary-duration behavior.

**Short operational proposal — requires separate authorization, not prepared/executed:**
1. Approve the specific owned scratch paths, restart targets, model/credential delivery and bounded paid allowance. Use an approved real model for a tool-bearing CLI file-read turn through EOF/save, then relaunch and recall without resupplying the fact; verify loading itself makes no provider request.
2. In an owned GUI/server/browser profile with a local fake provider and mic/read-aloud off, repeat EOF/save/relaunch, inspect old cards/controls and fresh noncolliding IDs. This is the browser evidence not supplied by backend tests.
3. Use independent owned copies for filename override, snapshot-only/full-log reconstruction, missing/corrupt primary backup recovery, both-invalid refusal and failed replacement preservation; check reported boundary, bytes and the next successful save without extra paid requests.

No operational scripts/harness preparations, browser/native/provider/key/discovery actions, cleanup, children, commits or manuscript/grader edits occurred. Only the assigned public exercise wording was reread for this proposal. **HOLD pending operational authorization; validation success is not a claim that the operational exercise has been completed.**

## Bounded operations — local PASS; live stopped at proxy setup guard

Owned scratch: `/tmp/e3-ch10-operations-wy8YEO` (macOS resolves `/private/tmp/...`); scripts/raws only under `control/`, outside all tool cwds. Source and fixtures unchanged. Consulted only the explicitly allowed prior scaffold files for launch/CDP/TLS mechanics, not their cases/budgets/assertions. Inspected current flags/readiness, DOM selectors/settings format, request encoding and tool-selection env.

Built exactly one GUI from main/gui with clean offline Go environment:
`go build -o /tmp/e3-ch10-operations-wy8YEO/control/gui ./cmd` — **exit 0** (`control/gui-build.log`). GUI SHA256 `31cf49210ecddff70af46ae2f412351715349fb03d425504a7bd81e211a1a91b`.
Required CLI unchanged: `/tmp/e3-ch10-final/frozen-cli-r3`, SHA256 `e148230ed0edbbf1d921fd0fcbe3680ff40cb5f41a9d293ed653c9d08610d980`. Pins and source-content aggregate were checked before/after local execution and before any credential read. Disk remained above 1 GiB. No real credentials entered a child environment, argv, file or log.

**Exact operational commands** (fresh owned HOME/TMP and minimal PATH `/usr/local/go/bin:/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin`):
- `env -i HOME=$B/home TMPDIR=$B/tmp PATH=$PATH /usr/bin/python3 -m py_compile $B/control/local.py` and `node --check $B/control/browser.mjs`: **0**.
- Same clean prefix `/usr/bin/python3 $B/control/local.py > $B/control/local-run.log 2>&1`: **0**. Controller launched one owned muted Chrome/profile/dynamic CDP port and GUI on dynamic port 60631, never 8088; one disposable GUI restart. Actual launch commands/PIDs/exits are in `control/local.jsonl`.
- `python3 -m py_compile $B/control/live.py`: **0** under the same clean prefix.
- Same clean prefix `/usr/bin/python3 $B/control/live.py > $B/control/live-run.log 2>&1`: **1**, stopped, no retry.

**Local observed PASS:** `control/local.pass`, `local-run.log`, `local.jsonl`, `browser.jsonl`, `frames.json` and `fake-{1..6}.{request,response}.json`.
- Active duration **4.224s**, six fake POSTs total: four for two actual GUI tool-bearing turns; two for actual frozen-CLI rendering captures. Native 1024 cap unchanged. No additional requests for load/replay/recovery/refusal/failure variants.
- First real `read_file` returned the owned fact file; paired IDs **2/3**. Browser observed those same IDs and old full-result control targeting the same long result after restart. Second real `write_file` produced `effect.txt` and distinct IDs **7/8**; old controls still addressed old results.
- Both GUI stdin EOF exits **0**, saved boundaries **12** and **24**; these operational EOFs followed completed browser turns. Existing active-terminal EOF receipt remains the separate evidence for draining an in-flight terminal turn.
- Actual CDP clicks, DOM, WebSocket and geometry inspected: viewport 1300×756; tools pane 360×577; old summary had visible positive bounds. Error DOM empty, mic=false, read-aloud=false, Chrome `--mute-audio`; no speech or visual paid call. **No pixel-level visual assessment claimed.**
- Explicit save override created no accidental default. Independent snapshot-only and full-log-only restores produced equal actual outbound request JSON. Missing/corrupt primary recovered the valid backup at boundary 12, reported backup recovery, successfully saved next primary, and preserved backup bytes.
- Both-invalid case: expected **exit 1**, bytes unchanged, zero HTTP. Nonempty `.bak` directory blocked atomic backup replacement: expected **exit 1**, valid primary and blocker unchanged, zero HTTP. These are expected negative-case outcomes, not an unexpected product failure.
- Owned GUI PIDs 74984/75005 exit 0, browser driver 74998 exit 0, Chrome 74985 exit 0; variant exits recorded. Proxy 60630, GUI 60631 and CDP 60634 confirmed closed. Only the exited owned profile was subsequently removed with a path-checked `shutil.rmtree`; raw/control evidence retained. No shared cleanup or hostwide zero-effects claim.

**Live stopped, not passed:** `control/live.jsonl`, `live-run.log`, `live-cli-1.stderr`, `live-models.json`; no `live.pass`.
- After source/binary checks, proxy read `directClaudeAPIKey` from the authorized settings file once, in memory only. One verified-TLS `GET https://api.anthropic.com/v1/models?limit=100` returned **200** and confirmed exact `claude-sonnet-4-6`.
- Single first CLI process started without an HTTP request on load. On prompting, proxy guard failed with `model/native cap/plain request guard` **before atomic POST reservation or forwarding**. **Upstream POSTs reserved/completed: 0/0; models GETs: 1; key reads: 1.** No fallback, redirect, retry or body rewrite.
- Runner mistakenly checked `body.get('stream') is False`. Current `internal/engine/render.go` uses `Stream bool json:"stream,omitempty"`; owned frozen-CLI local render capture has correct model/1024 and no `stream` member. This supports a runner guard/setup mistake, not a demonstrated product defect. The rejected live body was not captured before the guard, so there is no direct byte receipt for that body.
- Configured caps were max six upstream POSTs, 180s provider window and 240s absolute execution, verified TLS fixed endpoint, unchanged body forwarding, atomic reservation, watchdog/socket shutdown, read_file-only child tool selection, dummy child key. Stopped elapsed **1.162s**; no completed paid turn/save/recall claimed.
- Owned live CLI PID 75155 exited **1**; proxy port 60689 confirmed closed; key reference discarded during cleanup. Failure/raw local evidence preserved. No automatic runner fix, rerun, additional key read or additional models GET. Remaining allowance is not authorization to continue.

**HOLD / question to root:** whether to authorize a guard-only runner correction and explicitly address the already-consumed key-read/models-GET authority before resuming this same bounded case. Production remains held and unchanged; operational completion is not claimed.

### Authorized controller-only continuation — STOP / incomplete live exercise

Root independently inspected the v1 CLI API diagnostic `control/live-cli-1-api.log` and confirmed the correct model/1024/read_file request with absent `stream`. Root explicitly authorized one narrow controller-only continuation: version/preserve original files, absent-or-literal-false guard only (not null/0/truthiness), fresh cwd/save/output paths, reuse the successful model GET, one additional in-memory key read, max six **total** upstream POSTs, fresh 180s provider/240s absolute limits. No production correction or local/browser/fake/build rerun was authorized or performed.

Versioned runner: `$B/control/live-v2.py`, SHA256 `e11b1df0104f054984259449c1fbc7c75d972c1b0c17b457a03d07ee5fbf98dd`; fresh outputs `$B/control/live-v2/`, work `$B/live-v2-work`, save `control/live-v2/live-save.json`. Original runner/logs/failed save remained untouched. Syntax check under the same clean environment **exit 0**. Exact execution:
```
env -i HOME=/tmp/e3-ch10-operations-wy8YEO/home TMPDIR=/tmp/e3-ch10-operations-wy8YEO/tmp PATH=/usr/local/go/bin:/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin /usr/bin/python3 /tmp/e3-ch10-operations-wy8YEO/control/live-v2.py > /tmp/e3-ch10-operations-wy8YEO/control/live-v2/live-run.log 2>&1
```
**Exit 1** at watchdog, elapsed **180.234s including cleanup**. Full counters/events in `control/live-v2/live.jsonl`; original failure still preserved. Source/binary checks passed before the newly authorized in-memory key read. Model confirmation reused; no new discovery.

Observed partial real-model evidence:
- First CLI load made no HTTP. Two atomically reserved, unchanged native-body POSTs, each `max_tokens:1024`, each **HTTP 200**; both completed by about 3.27s after run start. No redirect/retry/substitution/fallback or body rewrite.
- Sonnet requested genuine `read_file` on `fact.txt`, native call ID `toolu_012gsZi5qbahnZR4zVpwUScG`. Next request included the actual successful tool result, `The sealed jar holds seven amber washers.\n`, paired to the same native ID. Second response stopped with `end_turn` and a correct complete answer, also present in `live-cli-1.stdout`.
- Provider usage: request 1 input **759**, output **64**; request 2 input **845**, output **23**; total input **1604**, output **87**, cache creation/read **0**. Raw requests/responses are `live-{1,2}.{request,response}.json` in the v2 output directory; no credential headers were captured.
- **Controller fault, not provider timeout:** v2 waited for both `assistant` and `usage` stdout records before EOF. Current `internal/cli/cli.go:140–205` emits each assistant immediately, but aggregate `usage` only after stdin EOF. The runner therefore withheld EOF despite the complete answer until its 180s watchdog. No runner edit/retry was made after this failure.
- Watchdog stopped owned CLI PID **75224**, exit **-15**; proxy port **60697** confirmed closed. No second process/recall. Actual archive remains fresh startup boundary **0**, empty log: **no successful paid EOF save/relaunch/recall is claimed**.

Cumulative Chapter 10 live allowance used: **2/6 upstream POSTs**, **1 models GET total**, **2 explicitly authorized in-memory key reads total** (v1 and v2). V2 alone GETs=0, key reads=1, reserved/completed POSTs=2/2. No further calls or credential reads. Remaining capacity is not a reason to continue. New owned processes/connections cleaned; prior evidence untouched; no additional profile cleanup. Local operational PASS stands, but the required real-model two-process exercise is incomplete. **Blocked HOLD for root/retained-reviewer review; source and fixtures remain held.**

### Script-only v3 prepared — NOT executed; awaiting root release

Root confirmed v2's controller deadlock (assistant before EOF, aggregate usage after EOF) and authorized only the narrow script correction and review artifacts. No production finding or production/fixture edit.

- Exact script: `/tmp/e3-ch10-operations-wy8YEO/control/live-v3.py`.
- Exact unified v2→v3 diff: `/tmp/e3-ch10-operations-wy8YEO/control/live-v3/v2-to-v3.diff`.
- Script SHA256: `abeef3bf0c943116189ec7ab401a6683a06df3c84072fe5b9047e75100e3dc4e` (`control/live-v3/runner.sha256.txt`).
- Syntax command: `env -i HOME=/tmp/e3-ch10-operations-wy8YEO/home TMPDIR=/tmp/e3-ch10-operations-wy8YEO/tmp PATH=/usr/local/go/bin:/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin /usr/bin/python3 -m py_compile /tmp/e3-ch10-operations-wy8YEO/control/live-v3.py` — **exit 0**; `control/live-v3/syntax.log`.

Diff is limited to fresh v3 output/cwd/save paths; atomic per-attempt POST cap **4** and matching final accounting; wait for the reliable assistant record alone; verify semantic answer; close stdin; wait for **exit 0**, then require usage and a nonempty saved canonical history. Before process two, require saved canonical `read_file` calls and matching successful native-ID results. Same corrected handshake applies to process two. No reconstruction from prior captures; v1/v2 scripts, outputs and failed saves preserved.

Unchanged: CLI/GUI/source pins, models-GET receipt reuse (**zero new GETs**), verified TLS/fixed endpoint, absent-or-literal-false stream guard, native 1024 cap, unchanged body forwarding, read_file/path scope, dummy child credentials, 180s provider window/240s absolute limit, stop/no-retry behavior. Intended future paths: `control/live-v3/`, `live-v3-work/`, and `control/live-v3/live-save.json`.

**Script-ready HOLD only.** Controller not executed; no operational process, external call, credential read, local/browser rerun, build, or validation gate. Cumulative live counters remain POST **2/6**, models GET **1**, authorized key reads **2**. A future key read/execution is expressly **not authorized yet**. Root must review and release before any continuation.

### V3 explicitly released and executed once — bounded live case PASS; prior failures retained

Root reviewed the exact script/diff and authorized execution of SHA256 `abeef3bf0c943116189ec7ab401a6683a06df3c84072fe5b9047e75100e3dc4e`, including one new proxy-only memory key read after pin checks. Script pin verified; disk **3,591,300 KiB** available. No controller edits, local/browser/fake/gate/build/research reruns or further discovery.

Exact command (source/binary checks also run inside the pinned controller):
```
env -i HOME=/tmp/e3-ch10-operations-wy8YEO/home TMPDIR=/tmp/e3-ch10-operations-wy8YEO/tmp PATH=/usr/local/go/bin:/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin /usr/bin/python3 /tmp/e3-ch10-operations-wy8YEO/control/live-v3.py > /tmp/e3-ch10-operations-wy8YEO/control/live-v3/live-run.log 2>&1
```
**Controller exit 0**. Case PASS at **4.615s**, provider window **4.599s**, full execution including cleanup **5.029s**. New reserved/completed upstream POSTs **3/3**, all HTTP 200, within the remaining cap 4; no new GET, one newly authorized key read. Fixed verified-TLS endpoint, unchanged native 1024 bodies, dummy child credentials, atomic reservation, no retries/redirects/fallback/body rewriting. Source and frozen CLI/GUI pins checked unchanged after success.

**Actual saved/reopened evidence** under `control/live-v3/`:
- `live-cli-1.stdout`: complete correct answer before EOF, then usage input **1604**, output **87** after EOF. CLI PID **75328**, exit **0**. Actual save `live-first-complete.json`, boundary **12**.
- Canonical `read_file` at sequence **5**, native ID `toolu_0158pHAHfFzHcUjSqM9XmG2E`, args `{"path":"fact.txt"}`; successful paired result at sequence **8**. Both verified in the actual saved log before process two launched.
- `live-cli-2.stderr`: loaded the v3 save through sequence **12**, backup=false. No request during startup. `live-3.request.json` contains the restored native call and matching tool result; no file/tool call during recall and the question did not resupply the fact.
- `live-cli-2.stdout`: correct recall, then post-EOF usage input **897**, output **19**. CLI PID **75329**, exit **0**. Final `live-save.json` boundary **17**; first saved canonical prefix unchanged.
- `live.pass`, `live.jsonl`, `live-run.log`, `live-{1,2,3}.{request,response}.json` preserve exact receipts. V3 provider totals input **2501**, output **106**, cache creation/read **0**.
- Both owned CLI processes exited 0 and owned proxy **60720** was confirmed closed. No further processes, profile cleanup or calls. Fresh v3 work/save paths only; v1/v2 artifacts preserved.

**All live attempts remain accounted for (not an overall-green rewrite):**
| Attempt | Outcome | Upstream POSTs | New models GETs | Authorized key reads |
|---|---|---:|---:|---:|
| v1 | exit 1, overly strict omitted-stream guard; no forwarded POST | 0 | 1 | 1 |
| v2 | exit 1, controller waited for EOF-only usage and hit watchdog; real read/answer but no successful EOF save/recall | 2 | 0 | 1 |
| v3 | exit 0, complete two-process read → EOF/save → reload → recall → EOF/save | 3 | 0 | 1 |

Cumulative live usage: **5/6 POSTs**, **1 models GET**, **3 explicitly authorized memory-only key reads**; provider input **4105**, output **193**, cache creation/read **0** across v2+v3. The one unused POST is unused. Prior local operational PASS remains qualified by DOM/geometry/WS evidence, not pixel assessment. No hostwide zero-effects claim, no new production correction round, no source/fixture changes or commits. **HOLD for root/retained-reviewer inspection; chapter acceptance remains theirs.**
