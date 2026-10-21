# Chapter 14 restart — student compiling HOLD

Implementation is held in `solutions/edition-3/main`; this is **not acceptance**. Initial whole-source review blocked four requirements. Production correction **ROUND 1/10** is now at a compiling HOLD for the same reviewer's delta check; source is steady. Earlier implementation evidence below remains valid where unchanged.

## Round 1 correction

- **One-shot round-trip decision:** `internal/common/history.go`, `internal/history/{history,persistence}.go`, `internal/agent/context.go` record the previous presented boundary. Only the newly consumed window is eligible for the course round-trip rule; previously kept results remain reclaimable by ordinary byte bands. Real Agent/local-HTTP coverage exercises keep, text-only/later turns below the bands, a new unkept batch, then band reclamation.
- **GUI draft ownership:** `gui/web/gui/app.js` keeps a target-specific draft and pending convergence separate from accepted settings. New typing invalidates the earlier pending submission. Remote snapshots, refusal and disconnection preserve the draft; an accepted normalized value updates an unchanged submitted draft. Status distinguishes unsaved text from shared settings. Tests invoke production input/socket handlers and feed actual settings-WebSocket frames through them (two local clients; no native browser).
- **Live paired references:** History attaches `AdmissionSeq` only to the latest genuine projected assistant call, copying its slice rather than changing the canonical response. `internal/engine/engine.go` uses that live provenance and a later successful result, never `History.EventsRange`. Startup/manual/unpaired fallback and provider legality remain. A guarded History fails any archive-range access during actual Engine renders, including after save and handoff; foreign provenance gets ordinary declarations rather than native references. Replay/non-mutation is checked.
- **Active-only membership:** `internal/agent/actor.go` deletes inactive `liveCalls` entries. Running eligibility is unchanged; the existing job owner/killed-state/late-result code is untouched. The focused membership regression and existing saved-admission/kill-boundary fixture pass; no wider job campaign was run.

Changed tests: `chapter14_correction_test.go`, `internal/{agent,history}/context_test.go`, `gui/ch14_context_test.go`, `gui/ch14_context_test.cjs`. The JS tests reuse the unchanged `gui/client_round1_test.cjs` VM fixture. The old unit expectation that kept results should be reclassified below the bands was narrowly corrected, with new public-path coverage of legitimate later reclamation.

All Round 1 logs are under `/tmp/e3-ch14-student/`. Genuine pre-fix failures: `round1-prefix-root.log` (keep/references), `round1-prefix-js.log` (three draft-loss cases), and the agent section of `round1-prefix-go-2.log` (inactive membership). The early Go commands also had new-fixture construction errors; these are not product verdicts. Later fixture corrections removed a nonexistent ToolData field, used the protocol's actual `current_settings` frame (the wrong name caused a timeout), and stopped applying a no-assistant-call helper to legitimate foreign replay. Their failed output remains in `round1-focused-go-final.log` and `round1-focused-gui.log`. Failed patch anchors stopped before subsequent owners; completed changes were inspected and continued, not restarted.

Final focused commands, each exit **0**, under the cleared offline environment below with `/opt/homebrew/bin` prepended to PATH for the existing Node binary:

```
# main; round1-focused-go-2.log
go test . ./internal/agent ./internal/history -run "TestChapter14|TestChapter09EnablingEventAndRendering|TestChapter09Round1DirectLoadExposure|TestChapter12UnloadAnthropicCatalog|TestCh10R3SemanticContextJSON|TestCh7R1SavedAdmissionFrontier" -count=1
# gui; round1-focused-gui-2.log (includes production JS plus real WS frames)
go test . -run "TestChapter14ContextTargetSettings|TestChapter14TargetDraftRealSettingsWS" -count=1 -v
```

The unchanged six-module driver/commands below ran sequentially again: exit **0**, `/tmp/e3-ch14-student/round1-compile-six.log`. No test bodies were selected for compilation. Disk checks stayed above 1 GiB. Scoped `git --no-pager diff --check -- .` passed. No broad/race/grade/lifetime/native/browser/live-provider checks, external operations, children, commits or acceptance snapshots. The same reviewer delta check and any later full validation remain unreleased. Assistance for this correction was the public four-item review brief plus current own source; no private/reference code was read.

## Initial HOLD ownership and implementation

- `internal/common` carries projected Context, survivor/declaration facts, redaction selectors, and the small actor/registry seams. `internal/history` owns total semantic reduction, held survivor ordering, recorded redaction replay, snapshot validation, and the `.journal` sibling. Known semantic junk remains canonical and is skipped/diagnosed; unknown/envelope failures still reject recovery.
- `internal/agent/context.go` owns the actor-local serialized-byte ladder and model-course round-trip policy. Context target uses zero as default (400,000 bytes), clamps positive values to 20,000, and is independent of tool-round budgets. Complete already-presented batches are eligible; unseen, incomplete, live-job and current keep-protected batches are excluded. Handoff follows its ordinary receipt, clears tool/private bytes after the complete batch, preserves visible dialogue and instruction survivors, then saves before asking the existing managed diagnostic writer to retain records.
- History keeps the current projection rather than reconstructing it for requests. Snapshot publication, journal append and rotation share its lock. Journal retention uses the validated backup boundary. Full raw archives remain in the existing snapshot format; explicit export/display and full saves can read the archived prefix. The in-memory canonical suffix is pruned after snapshot publication. No archive migration/new persistence architecture.
- `internal/skills` still owns discovery, activation and dispatch grants. Facts enter through the Agent mailbox; declaration deltas are computed under the registry lock without a callback into that lock. Compiled pending unload withdraws its sole grant immediately; independent grants and instruction survivors remain. Root composes restoration and freezes startup declarations.
- `internal/engine/context.go` assembles current SYSTEM and provider lanes. Only the two course rows are synthetic: course Sonnet gets inline SYSTEM skills, course Opus also gets startup catalog plus actual addition/removal blocks and flags. Real providers retain supported SYSTEM/declaration lanes. `internal/tool/context.go` exposes the three local context tools.
- Root configuration and GUI settings wire `context_target` for existing and new Agents. Existing diagnostics, jobs, cancellation, streaming, speech and MCP owners were extended/reused rather than replaced.

Focused new tests live in `chapter14_context_test.go`, `internal/{agent,engine,history}/context_test.go`, and `gui/ch14_context_test.go`. Inherited expectations changed narrowly: default tool count; sole compiled pending-unload grants; recognized semantic archive junk; backup recovery now includes its retained journal tail. Chapter 9's independently granted dependency replay remains checked unchanged in substance.

## Executed checks and limits

All Go commands used a cleared environment (`env -i`) with:

```
PATH=/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin
HOME=/tmp/e3-ch14-student/home TMPDIR=/tmp/e3-ch14-student/tmp
GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local
GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod
GOCACHE=/Users/bill/Library/Caches/go-build
```

Final retained focused commands (root cwd unless specified):

```
go test ./internal/history ./internal/agent ./internal/engine . -run "TestChapter14|TestChapter09EnablingEventAndRendering|TestChapter12UnloadAnthropicCatalog|TestCh10R2DirectToolArchiveValidation|TestArchiveBackupRecovery|TestBadArchiveDoesNotInstallOrOverwrite|TestCh10R3SemanticContextJSON|TestContractSeamQueuedResponseAfterCancelIsDiagnosticOnly" -count=1
```

`focused-final-root-2.log`: all three owner packages passed; aggregate exit **1** solely because my revised Chapter 9 assertion incorrectly ignored its still-active dependency grant. That assertion was restored with an explanation. Subsequent root command:

```
go test . -run "TestChapter14|TestChapter09EnablingEventAndRendering|TestChapter12UnloadAnthropicCatalog" -count=1
```

Exit **0**, `focused-root-final.log`. Unchanged actor/engine results above remain valid.

After the final History recovery diagnostics/continuity changes:

```
go test ./internal/history -run "TestChapter14|TestArchiveReplayTailAndNullContext|TestCh10R2DirectToolArchiveValidation|TestArchiveBackupRecovery|TestBadArchiveDoesNotInstallOrOverwrite|TestCh10R3SemanticContextJSON" -count=1
```

Exit **0**, `focused-history-hold.log`. GUI cwd:

```
go test . -run "^TestChapter14ContextTargetSettings$" -count=1
```

Exit **0**, `focused-gui.log`.

Final compile-only commands, sequentially, all exit **0**:

| Cwd relative to main | Command |
| --- | --- |
| `.` | `go test . ./internal/common ./internal/agent ./internal/engine ./internal/history ./internal/skills ./internal/tool ./internal/diagnostics ./cmd/ch01 ./cmd/ch02 -run '^$'` |
| `gui` | `go test . ./cmd -run '^$'` |
| `ch06` | `go test . -run '^$'` |
| `ch07` | `go test . -run '^$'` |
| `examples/ch04-consumer` | `go test . -run '^$'` |
| `examples/ch13-speech` | `go test . -run '^$'` |

Exact driver: `/tmp/e3-ch14-student/compile-modules.sh`; final output: `/tmp/e3-ch14-student/compile-hold-final.log`. These selected **no test bodies**; common, diagnostics, cmd/ch02, ch07 and ch04-consumer report no test files. Disk was checked before each compile and remained above 1 GiB.

Earlier failures remain in `/tmp/e3-ch14-student/`: initial missing Go PATH, incorrect field/anchor assumptions, new-fixture provenance/type errors, initial policy/recovery integration failures, and a nonexistent `cmd/ensemble` compile target. `compile-modules.log` retains that setup failure; the actual command paths were then inspected and corrected. `focused-*.log`, `compile-owners*.log`, `owners-build-*.log` retain iteration output. No failures were hidden by weakening unrelated coverage. A few failed patch scripts wrote nothing; one unchanged focused rerun followed such a script and is not counted as new evidence.

Limits/gates: no broad suite, race, grader, lifetime gate, root/reference comparison, browser/native/device/provider operations, installation, restart, or publication. Tests use local in-memory/filesystem/scripted HTTP fixtures, not real integrations. Managed diagnostic retention integration, long-running compaction under load, full GUI behavior and wider cancellation/jobs/speech regressions remain for authorized review/full validation. Root must inspect source/reference and the retained reviewer independently review before releasing those gates.

Assistance: current public task and entry-aligned Chapter 14, permitted public instructions, own accepted Chapter 13 source/baseline (including a permitted baseline type check). No Chapter 14 reference solution, private grader/review material, old-run implementation, child worker, or external assistance was used. No commit, snapshot or self-acceptance.

## Round 2 — latest stable HOLD

R2/10 addresses only the remaining negative-target draft blocker. Production change: `gui/web/gui/app.js` assigns a normalized pending-acceptance value only for nonnegative targets. A negative edit still reaches the existing server validator, but cannot be mistaken for an accepted 20,000-byte target when an unrelated broadcast arrives before refusal. No protocol/framework or other production owner changed.

Regressions extend `gui/ch14_context_test.cjs` and `gui/ch14_context_test.go`: accepted 20000 → type/submit -1 → unrelated shared 20000 → refusal, with the draft retained throughout. The actual settings owner supplies the queued unrelated/refusal frames to the production JS handlers. Existing valid normalization/ack, newer typing, disconnect and shared-state assertions remain.

Using the same cleared offline environment as R1, GUI cwd:
- `/opt/homebrew/bin/node ch14_context_test.cjs`: pre-fix exit **1**, `round2-prefix-js.log`.
- `go test . -run "^TestChapter14TargetDraftRealSettingsWS$" -count=1 -v`: pre-fix exit **1**, `round2-prefix-gui.log`. Both new production-JS paths genuinely lost `-1` to `20000`; no setup failures.
- `go test . -run "TestChapter14ContextTargetSettings|TestChapter14TargetDraftRealSettingsWS" -count=1 -v`: post-fix exit **0**, `round2-focused-gui.log`; all five JS cases and both selected Go tests passed. This compiles the changed GUI test package. No separate rebuild was needed for the JS-only production change; prior R1 Go/six-module evidence carries unchanged.

All logs above are in `/tmp/e3-ch14-student/`. Disk checks remained above 1 GiB (last 19,663,504 KiB free). Scoped whitespace check passed. No root/full/race/grade/lifetime/browser/provider checks or operations were launched. Source is steady for the retained reviewer's R2 delta check; this is not acceptance. Assistance was the public R2 brief and current own source only.

## Authorized local validation — STOPPED, held production

Only two small test files were added: `internal/history/ch14_boundaries_test.go` and `chapter14_boundaries_test.go`. No production edits or existing-test expectation changes were made; production correction count remains two.

The actual History test races 48 appends with 12 saves, verifies exact restart, then makes the primary pathname absent while preserving its bytes and verifies validated backup + intervening/journal-only tail recovery and continued append. The real Agent/local-HTTP integration uses the managed diagnostic writer, checks successful handoff snapshot/retention, deliberately blocks snapshot replacement while leaving the journal writable, verifies no retention after that failed save, and verifies continued complete diagnostic writes in both cases. Existing backup/torn-tail and Chapter 10 diagnostic tests were reused.

Commands used the same cleared offline environment/PATH/HOME/TMP and cache settings recorded above. From main:

```
go test ./internal/history -run "TestChapter14ConcurrentAppendSaveMissingPrimary|TestChapter14SnapshotBackupJournal|TestArchiveBackupRecovery" -count=1 -failfast -v
# exit 0; validation-boundaries-history.log

go test . -run 'TestChapter14HandoffSaveDiagnosticBoundary|TestDiagnosticRetention|TestCh10R1DiagnosticAppendDoesNotTrim|TestCh10R1DiagnosticRetentionBoundaries' -count=1 -failfast -v
# exit 0; validation-boundaries-agent.log

/usr/bin/time -p ... go test -race ./...
# exit 1; validation-main-race.log; real 616.34s, user 13.87s, sys 1.70s
```

All logs are retained under `/tmp/e3-ch14-student/`. The race command started 2026-10-21 16:29:22 UTC and exited 16:39:38 UTC. No short observation cancellation or retry occurred. Disk before that phase was 19,636,660 KiB free, above the 1 GiB floor.

**Blocking result:** root `TestCh7R3ActualAdmissionSave/inspect_saved_boundary` hit Go's 10-minute timeout. The retained dump shows its admission-save callback waiting at `chapter07_round3_test.go:57` via `History.AppendSaved` (`history.go:157`), while the inspecting test waits in `History.EventsRange` (`history.go:49`). Earlier root failures include canonical reopen/argument-prefix assertions (`TestCh10R2DirectDispatchArchive`, `TestCh10R3ToolArgumentSaveAndRenderPrefix`), fresh-root isolation, ephemera causal-dialogue assertions on all three providers, and Chapter 3 registry/declaration assertions. Some declaration failures explicitly expect ten tools rather than the current thirteen; this does not classify or excuse the other failures. `internal/cli.TestCh10R1FailedExitSaveKeepsDiagnostics` also failed its save-refusal assertion with JSON/archive-prefix errors. Full details remain in the unabridged race log; no blanket product-versus-fixture verdict or race-clean claim is made.

Stopped after that command's failure. **Not run:** main vet, GUI race/vet/remaining standalone JS suites, ch06/ch07/consumer/speech test+vet phases, frozen CLI/GUI builds/hashes, Chapter 14 black-box grader, or Chapter 1 grader. No grader source was read; no grader was executed. No production autofix, test weakening, repeat campaign, operations, children, commits, acceptance snapshot or cleanup. Test additions are formatted and the scoped whitespace check passed. Source is held pending root triage and a new release; no operational proposal is made because validation did not pass.

## Production R3/10 — compiling HOLD for same-reviewer delta

Production changes are limited to `internal/history/history.go`, `internal/history/persistence.go`, and `internal/skills/skills.go`. History now has a writer reservation separate from its committed-state mutex. Append, staged admission, save/rotation and load/repair remain serialized; `AppendSaved` releases only the committed-state mutex around its external callback, so readers see the OLD frontier while later writers cannot overtake staging. Nothing publishes before successful save/journal admission. The actual Agent dump callback does not acquire its snapshot `saveMu`. Skills construction now uses the same validated catalogue insertion helper without live declaration publication; normal registration still publishes after insertion. Root's existing initial-load/restore establishes effective grants before use. No persistence framework/format or public API was added.

Only the authorized fixture expectations changed: initial empty `tools_changed` plus strict root/path/recovery/CWD isolation (`ch10_round2_test.go`); all THREE causal entries and unchanged ephemera/replay checks (`chapter02_round2_test.go`); exact 13-name menus, schemas, empty selection and isolation (`chapter03_loop_test.go`, `chapter03_schema_test.go`, `chapter03_tools_test.go`); actual archive-prefix refusal wording while preserving exact 200-old + 10-fresh diagnostic bytes and both corrupt files (`internal/cli/ch10_round1_test.go`). The admission barrier and exact canonical reopen/argument-prefix tests were NOT weakened.

Added `chapter14_reopen_test.go`: default/selected/empty/recovered-skill construction must preserve exact canonical events, full projection and effective menus; subsequent live registration must emit exactly its genuine addition. Strengthened `internal/history/ch14_boundaries_test.go`: concurrent snapshot saves now overlap alternating plain and real staged-file appends, with OLD-frontier reads inside callbacks, durable staged-log parsing, restart and missing-primary backup-tail checks.

Commands used the same cleared offline environment and sequential limits as above, with `-count=1 -timeout=60s -v` on focused Go tests (the unchanged barrier still has its own checks; timeout was not increased). Logs below are in `/tmp/e3-ch14-student/`:

```
go test . -run 'TestCh7R3ActualAdmissionSave|TestCh10R2DirectDispatchArchive|TestCh10R3ToolArgumentSaveAndRenderPrefix|TestChapter14ReopenPreservesFactsAndEffectiveMenus' -count=1 -timeout=60s -v
# round3-admission-reopen.log: aggregate exit 1 solely from NEW-test setup mistakes
# (exclusive EventsRange boundary and wrong fixture primary name). All unchanged
# admission modes and exact reopen/argument-prefix assertions PASS on this source.

go test ./internal/history -run 'TestChapter14ConcurrentAppendSaveMissingPrimary|TestChapter14SnapshotBackupJournal|TestArchiveBackupRecovery|TestArchiveReplayTailAndNullContext|TestBadArchiveDoesNotInstallOrOverwrite|TestDialogueArchiveRetainsRecoveryContext' -count=1 -timeout=60s -v
# exit 0; round3-history-boundaries.log

go test . -run 'TestChapter14ReopenPreservesFactsAndEffectiveMenus|TestCh10R2BaseRelativeRecoveryIsolation|TestRound2EphemeraTailFailureAndReplay|TestChapter3LoopOrderedBatchProviders|TestChapter3LoopEmptyRegistryProviders|TestChapter3GeminiParameterProjection|TestChapter3Round1Registry|TestChapter09ActivationAndIsolation|TestChapter09EnvironmentAndStartupRenderer' -count=1 -timeout=60s -v
# exit 0; round3-startup-menus-fixtures.log (new test setup corrected)

go test . ./internal/cli -run 'TestChapter14HandoffSaveDiagnosticBoundary|TestCh10R1CLIOverrideAndRetention|TestCh10R1FailedLoadKeepsDiagnostics|TestCh10R1FailedExitSaveKeepsDiagnostics' -count=1 -timeout=60s -v
# exit 0; round3-agent-cli-retention.log

# From gui:
go test . -run 'TestCh10R1RetentionSettings|TestChapter14ContextTargetSettings|TestChapter14TargetDraftRealSettingsWS' -count=1 -timeout=60s -v
# exit 0; round3-gui-settings.log (production JS/actual settings frames included)
```

The existing explicit six-module compile sequence was reused with round-specific output names (`/tmp/e3-ch14-student/round3-compile-modules.sh`); aggregate and all module exits 0 in `round3-compile-modules.log` and `round3-compile-*.log`. Every invocation used `go test <explicit packages> -run '^$'`: root's recorded owner/command package list; gui `. ./cmd`; ch06, ch07, consumer and speech driver each `.`. These are compile-only results, including packages/modules with no tests. Disk before focused checks was 19,527,764 KiB free, and the script enforced the 1 GiB floor before each module. Scoped whitespace check exit 0.

The genuine stopped-validation failure remains unaltered in `validation-main-race.log`; no ten-minute baseline rerun was needed. Only the two new-test setup mistakes above occurred in this round; they are not sensitivity evidence. Prior passing checks for unchanged work were reused. No full race/vet campaign, graders, lifetime gate, browser/native/provider operation, production autofix beyond these two repairs, children, commits, acceptance snapshots or shared cleanup. No frozen binaries were built. Assistance: root's public R3 brief and allowed current source/tests, not private analysis/reference/grader code. Production is steady; renewed validation remains gated on root and the same retained reviewer's delta check.

## R3 renewed local validation — STOPPED at GUI race (held bytes)

No production or fixture edits were made. The inspected R3 focused results were reused. Fresh evidence is under `/tmp/e3-ch14-student/r3-validation/`; previous `validation-main-race.log` and every R3 correction log remain unchanged. The small `run-gate.sh` records command/CWD/disk/start/end/exit plus separate stdout/stderr, propagates failure, and refuses to overwrite a phase. Environment remained cleared/offline with GOMAXPROCS=2, GOFLAGS='-mod=readonly -p=1', GOTOOLCHAIN=local, GOPROXY/GOSUMDB=off, prescribed shared caches and existing Go/Delve/Node PATH. HOME/TMP were fresh owned `r3-validation/home` and `r3-validation/tmp`; no inherited key or LLM_BASE_URL.

Sequential actual commands/results:

| CWD | Command | Exit | UTC interval | Wall time |
|---|---|---:|---|---:|
| main | `go test -race ./...` | 0 | 2026-10-21 17:12:38–17:13:17 | 38.80s |
| main | `go vet ./...` | 0 | 17:13:17–17:13:18 | 0.77s |
| gui | `go test -race ./...` | 1 | 17:15:03–17:15:24 | 20.49s |

Full records: `main-race.{status,stdout,stderr}`, `main-vet.{status,stdout,stderr}`, and `gui-race.{status,stdout,stderr}` in that directory. The lowest phase-start disk reading was 19,387,224 KiB, above the 1 GiB floor. Commands finished naturally; no observation expiry was treated as failure and no test was retried.

GUI failures, unabridged in `gui-race.stdout`:
- `TestCh10R1SnapshotPrefixTailWS/snapshot_tail` and `/context_only`: `cannot recover history: primary: invalid context checkpoint` (fixture backup absent), `ch10_round1_test.go:89`.
- `/full_carry`: `canonical prefix ID rewritten`, `ch10_round1_test.go:216`.
- `TestCh7R1TerminalUnboundedInput`: `unbounded terminal input lost`, `round1_test.go:442`.
- `TestSettingsPersistenceAndFailedSave`: `settings leaked into history`, `settings_test.go:83`.

No source triage, fixture alteration, retry or blanket product/fixture classification was performed after this unexpected failure. The same GUI command completed `ensemble/gui/cmd` successfully, but the GUI gate as a whole failed. GUI vet was blocked by shell failure propagation. No standalone JS suite was subsequently launched (Chapter 14 JS is invoked by the GUI Go suite); no ch06/ch07/consumer/speech test+vet phases, frozen CLI/GUI builds/hashes, Chapter 14 grader or Chapter 1 grader ran. No private grader source was inspected. No operations/proposal, children, installs, commits, acceptance snapshots or cleanup. Production correction count remains 3/10. Source is steady and validation is stopped for root triage.

## GUI fixture-only repair — HOLD for root inspection (production still R3, 3/10)

Only `gui/ch10_round1_test.go`, `gui/round1_test.go`, and `gui/settings_test.go` changed. Exact diff: `/tmp/e3-ch14-student/gui-fixture-repair/fixture.diff` (9,270 bytes). No production edits; inspected R3 main race/vet results carry unchanged, and `r3-validation/gui-race.*` remains intact.

- Snapshot fixture now captures the actual prefix `AsOf`, selects the three canonical tail facts by sequence, and writes matching checkpoint/context data to a fresh owned destination separate from the seed's journal/backup. Context-only also bounds its derived display archive by `SourceSeq`. Checks require no journal/backup side input, the intended loaded path/boundary/mode and exact supplied canonical events. Immediately before the context-only live append, the canonical log must still be empty at the prefix boundary and initial WS replay must contain no result. The canonical ID check finds exactly the `native-call` tool call and requires original DisplayID 7. Existing live calls/results/jobs links, counts, new/stable IDs, reconnect, save and exact-produced-archive reopen checks are unchanged.
- Terminal fixture captures the canonical frontier before Terminal, then waits for a newly admitted HUMAN text message containing the exact entire trimmed >1 MiB payload. Startup facts no longer satisfy the wait. Pause/offline behavior and the existing wait remain; no provider completion is required.
- Settings fixture captures the complete canonical event baseline before operations and requires exact equality afterward. Existing broadcast, persistence, live policy, failed-save/no-broadcast, old-file and reconnect assertions remain unchanged.

From the GUI module, with cleared prescribed offline environment/shared caches and fresh owned `gui-fixture-repair/home` and `tmp`, ran ONLY:

```
go test . -run '^(TestCh10R1SnapshotPrefixTailWS|TestCh7R1TerminalUnboundedInput|TestSettingsPersistenceAndFailedSave)$' -count=1 -failfast -v
```

Evidence in `/tmp/e3-ch14-student/gui-fixture-repair/`:
- `selected.{status,stdout,stderr}`: exit 1, 0.65s, fixture compilation only. I initially used nonexistent `DialogueRecord.Seq`; no test bodies ran. Corrected solely to the actual `SourceSeq` field after reading its definition. This is setup evidence, not product/sensitivity failure.
- `selected-2.{status,stdout,stderr}`: exit 0, 1.88s wall (2026-10-21 17:34:07–17:34:09 UTC), package test time 0.584s. All three snapshot variants, full terminal payload, and settings test PASS.
- Scoped `git diff --check` exit 0. Lowest pre-run free disk 19,301,292 KiB, above floor.

No other test selection, full suite, race/vet campaign, build command, grader, browser/native/provider operation, keys/discovery, children, commit/snapshot or cleanup. No unexpected product failure occurred. No validation continuation is claimed or started: fixtures and held production are ready for root inspection first.

## Validation resumed after accepted GUI fixtures — STOPPED at client JS

Source and fixtures stayed held; production remains R3 / 3 rounds. Main race/vet were carried from `r3-validation/` and not repeated. Fresh logs are `/tmp/e3-ch14-student/r3-validation-2/`, using fresh owned HOME/TMP, the same cleared offline environment/shared caches, and first-failure shell propagation. Previous failed GUI race and fixture-repair evidence remain unchanged.

Actual sequential GUI-CWD commands:

| Phase | Command | Exit | UTC interval (2026-10-21) | Wall |
|---|---|---:|---|---:|
| gui-race | `go test -race ./...` | 0 | 17:42:49–17:43:11 | 21.68s |
| gui-vet | `go vet ./...` | 0 | 17:43:11–17:43:11 | 0.50s |
| js-ch08 | `node ch08_round1_test.cjs` | 0 | 17:44:02–17:44:03 | 0.03s |
| js-client | `node client_round1_test.cjs` | 1 | 17:44:03–17:44:03 | 0.02s |

Each phase has exact command/CWD/disk/timestamps/exit in `.status`, full `.stdout`, and full `.stderr` with timing. GUI race passed both `ensemble/gui` and `ensemble/gui/cmd`. The controlled-JS selection excluded duplicate Chapter 14 execution (already invoked by GUI Go); the normal standalone client entry executes checks not run merely by importing its fixture.

First unexpected failure: `AssertionError [ERR_ASSERTION]: 1 !== 2`, actual 1 / expected 2, at `gui/client_round1_test.cjs:35:50`, Node v24.20.0. Full stack is `js-client.stderr`; stdout is empty. This is an executed assertion failure, not a setup failure. No retry, source triage, production/fixture edit or product-versus-fixture classification was attempted.

Shell propagation prevented MCP DOM and Chapter 13 correction JS suites from starting. No remaining ch06/ch07/consumer/speech test/vet gates, frozen CLI/GUI builds or hashes, Chapter 14 grader, or Chapter 1 grader ran. Lowest phase-start free disk was 19,240,416 KiB. All started commands finished naturally. No private grader inspection, browser/native/live-provider/discovery/key operations, children, commits, acceptance snapshots or cleanup. No operational proposal: validation is stopped for root triage, not accepted.

## Invocation-corrected continuation — STOPPED at Chapter 14 GUI startup

No source or fixture edits; production remains R3 / 3 rounds. Carried main/GUI race+vet and fresh Chapter 8 JS results without rerunning them. New logs are `/tmp/e3-ch14-student/r3-validation-3/`, with exact commands/CWD/disk/UTC start/end/exit in each `.status`, raw `.stdout`, and `.stderr` including timing. Clean offline environment/shared caches retained; fresh owned HOME/TMP; minimum phase-start free disk 19,149,228 KiB.

Authorized corrected JS invocations (GUI cwd), both exit 0 / 0.03s:
- `node client_round1_test.cjs chapter12` → `js-client-ch12.*`. This credits only the later trusted/synthetic typing assertions. The entire legacy Chapter 7 block was skipped; adversarial DOM, silent replay, long queue and direct late-delta assertions are NOT newly proved by this mode. The failed default `r3-validation-2/js-client.*` remains failed and preserved.
- `node ch13_correction_test.cjs` → `js-ch13.*` (its actual controlled cases passed).

Remaining module commands ran sequentially as `go test ./...` followed by `go vet ./...`:

| module | test/vet exits | wall seconds (test/vet) | qualification |
|---|---|---|---|
| ch06 | 0 / 0 | 1.17 / 0.13 | actual tests (package 0.252s) |
| ch07 | 0 / 0 | 0.11 / 0.11 | no test files: compile-only test command |
| examples/ch04-consumer | 0 / 0 | 0.10 / 0.11 | no test files: compile-only test command |
| examples/ch13-speech | 0 / 0 | 0.88 / 0.13 | actual tests (package 0.614s), not native audio |

Logs: `ch06-{test,vet}.*`, `ch07-{test,vet}.*`, `consumer-{test,vet}.*`, `speech-{test,vet}.*`. No additional standalone JS suite was added beyond the latest two-command release.

Built exactly once each (both exit 0; 0.80s / 0.37s): main `go build -o /tmp/e3-ch14-student/r3-validation-3/frozen-cli ./cmd/ch01`; GUI `go build -o /tmp/e3-ch14-student/r3-validation-3/frozen-gui ./cmd`. Logs `build-cli.*`, `build-gui.*`. Hashes in `frozen.sha256`, checked before and after the failed grade without rebuilding:

```
676fa756a1e9beaa2e45be719f8f003b7443fbd01990465a336b311123850a0e  frozen-cli
f806ff1ab39966f92aaa7e226602c19df532eb4edf8ca2ac18422cb048704456  frozen-gui
```

From book root `/Users/bill/projects/ensemble`, ran exactly:

```
go run ./cmd/grade-edition3 -ch 14 -json /tmp/e3-ch14-student/r3-validation-3/frozen-gui
```

**Exit 1**, 1.56s (2026-10-21 17:56:28–17:56:29 UTC). Raw JSON: `grade14.stdout`; stderr/timing: `grade14.stderr`; command/exit/times: `grade14.status`. Machine-parsed `grade14.parsed` records `passed:false`, score 0/100, 8 checks. The output reports GUI exit before readiness because `/tmp/e3-ch14-student/r3-validation-3/tmp/context-grade-484661564/gui/index.html` is missing. Sessions reporting request counts made zero model requests; dependent checks did not run their sessions. This is a failed grader launch as executed, not evidence that context-policy semantics were exercised and failed. No private grader or temporary grader-fixture source was inspected.

Stopped at this first unexpected failure. No retry, asset staging, production/fixture change or rebuild. Chapter 1 was NOT run, so genuine 62s-header + 62s-body/full-usage behavior remains unverified by that gate. Both frozen binaries are held unchanged for root triage. No browser/native/live-provider/key/discovery operation, children, commits, acceptance snapshots or cleanup. No operational proposal because the released validation did not all pass.

## One renewed Chapter 14 grade — STOPPED, 85/100 (same frozen bytes)

Root authorized one renewed GUI-direct grade after its independently reviewed inert-static-fixture setup correction. No private grader source or temporary fixture files were inspected. No student source/fixture edits, local-gate repeats or rebuilds; production remains R3 / 3 rounds. Prior 0/100 startup failure remains intact in `r3-validation-3/`.

Fresh evidence: `/tmp/e3-ch14-student/r3-validation-4/`. From book root ran exactly:

```
go run ./cmd/grade-edition3 -ch 14 -json /tmp/e3-ch14-student/r3-validation-3/frozen-gui
```

Exit **1**, 3.12s, 2026-10-21 18:09:30–18:09:33 UTC. Raw JSON `grade14.stdout`, stderr/timing `grade14.stderr`, exact command/CWD/disk/timestamps/exit `grade14.status`, machine-parsed summary `grade14.parsed`. Parsed `passed:false`, score 85/100; six checks passed. Failures are preserved without root-cause guesses:

- `frozen-prefix`: `request 2: skill body occurs 0 times in SYSTEM, want 1`.
- `total-reducer`: fixture preparation stopped at `lstat /tmp/e3-ch14-student/r3-validation-4/tmp/context-grade-182472604/C/skills: no such file or directory`; this check did not exercise the reducer.

Both existing frozen hashes matched before and after (`frozen-before.check`, `frozen-after.check`); binary paths/hashes remain those recorded above in `r3-validation-3/frozen.sha256`. Clean prescribed offline environment, fresh owned HOME/TMP and shared caches; pre-grade free disk 19,073,280 KiB. The grader completed naturally. Process0 AND parsedpass were not achieved, so FIRST Chapter 1 was NOT run and its genuine 62s+62s/full-usage gate remains unverified.

Stopped without retry, autofix, source/fixture change, asset staging, rebuild, private-source inspection, browser/native/live-provider/key/discovery activity, children, commits, snapshots or cleanup. No operational proposal. Stable HOLD pending root triage of the exact black-box failures.

## Production R4/10 — compiling HOLD for same-reviewer delta

Narrow Engine-only correction: `internal/engine/context.go` now appends every active/pending surviving manual to SYSTEM for every model and excludes `KindSkill` from wire dialogue without changing History. `internal/engine/render.go` reserves inline SYSTEM dialogue handling for declaration entries only. `InlineTools` still selects frozen startup declarations plus chronological deltas; the native paired-reference payload branch is unchanged. No History, registry, replay, journal, handoff or grant-owner production changes.

Focused tests:
- `internal/engine/context_test.go` now inspects native SYSTEM separately from dialogue on both actual course rows, Anthropic, OpenAI and Gemini; whole-payload occurrence alone no longer suffices.
- New `chapter14_instruction_test.go` uses real Agents/registry/History, local HTTP and local file tools. Both course rows and real Anthropic cover load → following request, complete two-call/result batch ordering, stable primary/manual SYSTEM exactly once with no manual in dialogue, explicit tool clearing, pending unload/grant withdrawal, handoff, canonical replay and snapshot/journal restart. Synthetic startup declarations remain unchanged with ordered add/remove groups. Genuine Anthropic reference-only results, ordinary fallback after clearing, canonical History immutability and post-restart requests are checked.

Exact R4 delta: `/tmp/e3-ch14-student/round4/change.diff` (against the directly inspected pre-edit files, not private/reference inputs). Source/fixtures outside these four files were not edited. Whitespace check passed.

Executed from main with the prescribed clean offline environment/shared caches and owned `round4/home`, `round4/tmp`:

```
go test ./internal/engine . -run '^(TestChapter14ProviderContext|TestChapter14InstructionSystemLifecycle)$' -count=1 -timeout=60s -v
```

`round4/prefix.log`: exit 1, 2.93s. Genuine pre-fix failures only on Opus course: late manual count 0 in SYSTEM and present in dialogue, also across clearing/unload/handoff/replay/restart. Other provider/lifecycle cases passed. No setup failure in this sensitivity run.

```
go test ./internal/engine . -run '^(TestChapter14ProviderContext|TestChapter14InstructionSystemLifecycle|TestChapter14RenderUsesLiveLoadPair|TestChapter09DynamicRenderAndProviderDeclaration|TestChapter09EnablingEventAndRendering|TestChapter09SameResponseLoadDoesNotAuthorizeSibling)$' -count=1 -timeout=60s -v
```

`round4/focused.log`: exit 0, 3.20s. Actual reported tests: provider-field control, lifecycle matrix, guarded live pairing/provenance and existing enabling/rendering control. Two optional selector names matched no tests and are not credited; the actual existing controls were then run explicitly:

```
go test ./internal/engine . -run '^(TestSkillToolReferencesAreContentNotInstructions|TestChapter09Round2OrdinaryResults|TestChapter09Round1DirectLoadExposure|TestChapter09EnvironmentAndStartupRenderer)$' -count=1 -timeout=60s -v
```

`round4/reference-controls.log`: exit 0, 0.86s. Includes native reference-only content, successful ordinary fallback/failed load/missing result/ordinary tool, Dispatch/Invoke/manual direct load and startup renderer.

Explicit six-module compile driver `/tmp/e3-ch14-student/round4/compile-modules.sh` exited 0; `round4/compile.log` and `compile-{root,gui,ch06,ch07,ch04-consumer,ch13-speech}.log` contain exact package commands/exits. Every command used `go test <explicit packages> -run '^$'`; **no test bodies selected**. No-test packages/modules are compile-only. No full suite/race/vet/grade/lifetime or frozen-binary rebuild.

Setup qualifications: one patch attempt stopped on a comment-anchor mismatch before source writes/tests; one delta-record generation attempt stopped on a test-block boundary mismatch without source changes. Both were corrected from the actual small source sections, not counted as product failures or sensitivity evidence. Lowest observed disk 18,834,908 KiB.

`round4/frozen-older.check` verifies the two R3 frozen hashes still match; they are OLDER-candidate evidence, not R4 binaries. Original 0/100 and 85/100 false grade outputs remain unchanged. Root's separate total-reducer fixture-layout repair was not inspected or touched. FIRST Chapter 1 remains unrun. Assistance was the public R4 brief/current public chapter and own source/public fixtures only; no private grader/review/original inputs, children, manuscript/guidance edits, commits, snapshots or cleanup. HOLD for the same reviewer and root release before renewed validation; no chapter acceptance or operations.

## R4 validation completed — saved HOLD, root acceptance/operations pending

No source or fixture edits during this release. Fresh evidence is under `/tmp/e3-ch14-student/r4-validation/`, using owned `home`/`tmp`, the prescribed cleared offline environment/shared caches, and no inherited keys or `LLM_BASE_URL`. Every gate checked the 1 GiB floor; lowest observed free space was 18,682,788 KiB. Commands ran sequentially, with first-failure propagation. All newly released gates passed.

Fresh affected-owner gates (each `.status` records exact cwd/command/start/end/exit; `.stdout` and `.stderr` retain full output/time):

| Cwd | Command | Exit | Wall seconds |
|---|---|---:|---:|
| main | `go test -race ./...` | 0 | 40.90 |
| main | `go vet ./...` | 0 | 0.62 |
| main/gui | `go test -race ./...` | 0 | 20.68 |
| main/gui | `go vet ./...` | 0 | 0.26 |

Carried per root authorization, not rerun as separate campaigns: accepted Chapter 8 JS; `node client_round1_test.cjs chapter12` ONLY the later trusted/synthetic typing block; `node ch13_correction_test.cjs`; ch06 and ch13-speech actual tests plus vet; ch07 and ch04-consumer no-test compile-only commands plus vet. Their exact earlier logs remain in `r3-validation-2/` and `r3-validation-3/`. GUI Go tests invoked their normal JS integrations as part of this fresh GUI run. The old default-client failure remains failed: its skipped legacy/adversarial-DOM/silent-replay/long-queue/direct-late-delta assertions are NOT newly proved. R4 six-module compile-only evidence is retained in `round4/`, not rerun here.

Built ONCE each after affected gates, no rebuild between grades:
- Main: `go build -o /tmp/e3-ch14-student/r4-validation/frozen-cli ./cmd/ch01` — exit 0, 0.64s.
- GUI: `go build -o /tmp/e3-ch14-student/r4-validation/frozen-gui ./cmd` — exit 0, 0.33s.
- CLI SHA-256: `05663974c2785fca2d5a3b0602c22be42d838a244a3d994de2504ba636e17a58`.
- GUI SHA-256: `fddf82b3d230a244d88eb536c81e3f042488b7f1a5c4d29d9819319eaedb8f22`.
- `frozen.sha256`, `frozen-before-ch1.check`, and `frozen-final.check` record identity; both hashes match after grading. R3 binaries and all earlier failures remain preserved as older evidence.

Black-box commands from `/Users/bill/projects/ensemble`, without private-source inspection:

```
go run ./cmd/grade-edition3 -ch 14 -json /tmp/e3-ch14-student/r4-validation/frozen-gui
go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch14-student/r4-validation/frozen-cli
```

- Chapter 14 ONCE: exit 0; parsed `passed:true`, 100/100; all eight checks passed, including frozen-prefix and total-reducer. Wall 3.17s. `grade14.stdout` is the complete JSON, with separate stderr/status and `grade14.parsed`.
- FIRST Chapter 1 ONCE, released only after that pass: exit 0; parsed `passed:true`, 100/100. Wall 124.89s, 18:53:42–18:55:47 UTC. JSON explicitly confirms complete answer/usage after separate genuine **1m2s headers** and **1m2s body** delays. Five conversation rounds/requests; cumulative usage input 511/output 65 matches the fixture. No outer cancellation deadline was imposed (allowance exceeded 360s); no short-observation cancellation or retry. Full JSON/stderr/status are `grade1.*`.

No provider/key/discovery, browser/native/device, child-worker, deployment/restart, commit/snapshot or cleanup operation was performed. No old failed suite/default-JS/grade result was erased or relabelled. A note-write call initially omitted its path and was rejected without writing; it was corrected without rerunning any gate. This is validation evidence and HOLD, not self-acceptance of the chapter or authorization for operations.

### Conditional proposal only — public §14.16

After explicit approval of one supported real-model selection, credential channel, spend/stop limits and an owned local launch, use one stdin-driven GUI session with actual GUI assets and a scratch workspace containing the public base/manual skills and six ordinary record files. Use `context_target=20000`, `max_tool_rounds=4`, TTS off, and the real-model context flags left false (not synthetic course aliases). Load/read record 1, then read/report records 2–6 in separate human turns; request exactly one handoff, then unload and report first/last codes plus instruction survival. Inspect provider-bound SYSTEM/declarations, whole result presentation, recorded cuts/narration, handoff and withdrawn grants—not just the final answer. Close normally via EOF, retain save/journal/request log together, and perform only an explicitly approved restart check. Stop on repeated checkpoints or the agreed cost/turn boundary; no provider matrix or native/browser campaign. This is a proposal, not a script or permission to access keys or execute it.

## Sonnet 4.6 preparation — HOLD after one fake rehearsal setup failure

Production/source, frozen R4 CLI/GUI and actual GUI assets remain unchanged. Prepared only owned `/tmp/e3-ch14-student/sonnet46-preparation/` content; no real credential access, provider/discovery request, browser/native/MIC/audio operation, child agent, rebuild or shared cleanup.

New single-purpose stdlib Python launcher/proxy: `launcher.py`, SHA-256 `3fa4e75f060548584817fca5cdbbf637bc34e184b15e0024ed0e79bd40b6f4c2`. No suitable prior Python launch scaffold was reused; no private/reference/runtime inputs were consulted. `pins.sha256` is a plain hash list covering the frozen GUI, launcher, copied full assets and immutable public fixture inputs. `assets/{index.html,app.js,style.css}` match the actual held source files byte-for-byte, not inert grader files. Inputs reproduce public base/manual and six ordinary records plus the three settings. Prepared seven stdin turns, exact Sonnet 4.6/model/native-1024 checks, unchanged-byte forwarding, atomic irrevocable reservations, fail-closed stop latch, a fixed 24-POST live cap, 420s active / 480s absolute (five seconds reserved for owned cleanup) / 60s setup, and normal EOF/join on success. Unsupported methods/routes, framing, transport/provider errors, out-of-fixture tools/reads and repeated/out-of-order handoffs stop forwarding. Future live mode is unexecuted: one certificate-validated `/v1/models` GET, one direct `~/.cr/settings.json` / `directClaudeAPIKey` read inside proxy memory, no pagination/fallback/redirect/retry or key in argv/env/child/files/logs. These live-path safety properties are prepared code for root review, NOT live-tested claims.

Exact controlled commands (prescribed cleared offline environment, owned check HOME/TMP and existing caches):

```
python3 launcher.py --mode rehearse --run-name must-not-launch --pins wrong-pins.sha256
python3 launcher.py --mode cap-check --run-name cap-check
python3 launcher.py --mode rehearse --run-name rehearsal-01
```

- Wrong-pin control: expected exit 1, exact frozen-GUI hash mismatch. No run directory, listener or GUI was created. `wrong-pin.{stdout,stderr,status}` retained.
- Reduced-cap local HTTP control: exit 0; one actual fake upstream POST, then statuses 503/503 with the stop latch closed. One local fake catalog GET, zero real key reads. `cap-check/receipt.json`, raw request/fake-received/response files and trace retained. This exercised cap=1 through the real proxy path, not 24 requests and not a mutation campaign.
- ONE real frozen-GUI fake rehearsal: exit 1, wall 60.16s, setup deadline at elapsed 60.0s. The GUI started with `--port 0`, full owned `--gui-dir`, explicit `--skills-dir`, `--save`, `--settings`, `--log-mode full`, owned cwd/HOME/TMP, dummy key, explicit loopback `ANTHROPIC_BASE_URL`, no generic inherited `LLM_BASE_URL`, and nonstreaming. It announced its loopback listener and wrote startup `state.save`; however launcher readiness never succeeded. Zero stdin turns, zero reserved POSTs, one LOCAL fake catalog GET, zero real key reads. No external model/discovery request occurred. No retry or production repair was attempted.

Qualification: this is a launcher-readiness/setup failure, not evidence of a Chapter 14 product defect. The readiness predicate invokes `lsof -Fn` without `-nP` yet accepts only literal `n127.0.0.1:`; its output was not retained, so that representation mismatch is a concrete inspection concern rather than a claimed observed cause. The script was not changed/re-pinned after the failed rehearsal. Seven-turn completion, full provider-bound reads, SYSTEM/manual survival, ordered cuts/handoff/unload, normal EOF/save and final operational receipt remain UNPROVED in this rehearsal.

Evidence: `rehearsal-01.{stdout,stderr,status}`, `rehearsal-01/{receipt.json,trace.jsonl,gui.stdout,gui.stderr,work/state.save,work/state.save.journal}`, and `post-failure-observation.json`. Startup save is AsOf 2 with only `skill_loaded`/`tools_changed`, zero request/human facts. Owned GUI PID 10537 no longer exists; the launcher finished its owned listener/process closure. `pins-after.check` and `frozen-after.check` verify unchanged pins/frozen hashes. No files or shared processes were manually cleaned. Minimum observed free disk 18,516,284 KiB.

HOLD for root inspection. No paid/live/key/discovery release, further rehearsal, restart or operations are assumed. Earlier R4 validation/grades and qualifications remain intact; this preparation failure is separate and does not consume a production correction round.

## Preparation controller v2 — local fake rehearsal PASS, live remains closed

Controller-only change; production/fixtures/candidate binaries remain held, no production round. Preserved v1, its pins and failed `rehearsal-01` intact. New owned files under `/tmp/e3-ch14-student/sonnet46-preparation/`: `launcher-v2.py`, `pins-v2.sha256`, `v1-to-v2.diff`, `rehearsal-02.command`. V2 SHA-256: `0d442acaeb1a4f78e6e16c061151a3ae34ccf6645366c72f32ae91c525d1dbbc`.

Exact delta: readiness uses `lsof -nP`, retains last return/stdout (4096-byte bound), verifies the owned Popen PID alive before/after lookup and refuses listener 8088. Response validation uses global + current-batch handoff count, permitting ordered `micro_handoff` then `unload_skill` companions while retaining reversed/repeated/wrong-turn refusals and final canonical order checks. Only default pin filename changed otherwise. No input, cap, timeout, TLS/CA, production or final-grant-check changes; base correctly has no read_file grant. V1 historical timeout cause remains unproved because its lsof stdout was not retained.

Syntax/exact-pin and seven small pure batch controls passed: ordered companions, separate batches; refusal of reversed, repeated local/global, wrong turn and unload-before-handoff; rejected batches did not commit count. `v2-pure-controls.{stdout,stderr,status}` retained. No listeners/processes/GETs/POSTs/key reads in these controls. Unchanged wrong-pin/cap controls were carried, not repeated.

ONE GUI fake rehearsal, with command recorded before launch in `rehearsal-02.command` (same clean offline environment/shared caches/owned HOME/TMP):

```
python3 /private/tmp/e3-ch14-student/sonnet46-preparation/launcher-v2.py --mode rehearse --run-name rehearsal-02 --pins /private/tmp/e3-ch14-student/sonnet46-preparation/pins-v2.sha256
```

Exit 0; wall 1.20s, receipt elapsed 1.099s. Owned GUI PID 12076/listener `127.0.0.1:54589`, lsof return0/numeric output retained. Seven stdin turns ended canonically at seq23/36/49/62/75/88/112; normal stdin EOF/join, GUI exit0. Six complete 8051–8054-byte read results reached provider-bound requests3/5/7/9/11/13. Six L2 result cuts at seq25/38/51/64/77/90; L3 paired-tool cut91; one handoff99 before unload105. Final snapshot/context/log all AsOf112. Manual remains pending with body; final effective tools are load_skill/unload_skill/compact_context/micro_handoff/keep_tool_results, no read_file grant/declaration.

All 16 requests used exact Sonnet4.6/native1024. Primary SYSTEM count1 throughout; manual SYSTEM count0 initially then1 for requests2–16, manual dialogue count0 throughout. Request bytes matched local fake upstream captures. One LOCAL fake catalog GET, 16 LOCAL fake POSTs, zero real-key reads/external requests; no browser/native/audio/discovery/paid operation. This proves the local driver/owner interaction with scripted fake responses, NOT real-model judgment or live-key/TLS/provider success. Same-batch companions were pure-controlled; this GUI rehearsal used separate handoff/unload batches.

Evidence: `rehearsal-02/{receipt.json,inspection.json,trace.jsonl,request-*.json,fake-received-*.json,response-*.json,gui.stdout,gui.stderr,work/state.save,work/state.save.journal}`, top-level `rehearsal-02.{stdout,stderr,status}`. A later ad-hoc read-only cross-check used an incorrect primary sentence and failed; source/input/rehearsal were not changed or rerun. Reading the actual pinned input confirmed the launcher’s original marker was correct; `inspection.json` records actual field counts with that marker. No product failure inferred from that analysis error.

`pins-v2-after.check`, `pins-v1-preserved.check`, `frozen-v2-after.check` pass; owned PID no longer exists. Minimum observed free disk18,449,972KiB. No shared cleanup, broader validation, rebuild or grade. Stable preparation HOLD for root inspection and a separate finite live release; no credential/provider action is authorized by this result.

## Finite live-01 — STOPPED / HOLD, no remaining authorization

Executed the ONE root-released live attempt with unchanged approved v2/pins/R4 GUI/assets/inputs. Preflight pins matched launcher SHA `0d442acaeb1a4f78e6e16c061151a3ae34ccf6645366c72f32ae91c525d1dbbc`; disk18,414,036KiB. Exact interpreter/command (full cleared environment in `sonnet46-preparation/live-01.command`):

```
/opt/homebrew/opt/python@3.13/bin/python3.13 /private/tmp/e3-ch14-student/sonnet46-preparation/launcher-v2.py --mode live --run-name live-01 --pins /private/tmp/e3-ch14-student/sonnet46-preparation/pins-v2.sha256
```

UTC19:52:41–19:53:17, exit1; wall35.18s, receipt elapsed35.019s. Stop latch: `unexpected/repeated handoff` on response15 during human turn7. No retry, reserve extension, fallback, independent diagnostic/provider run, body rewrite, rebuild or source/fixture/controller change. All unused budget expires; live-01 is NOT a successful seven-turn exercise or normal-EOF result.

Actual real observations, distinct from prior fake/control evidence:
- One memory-only key open inside the approved proxy and one certificate-validated models GET confirmed exact `claude-sonnet-4-6` before POST. No other key read, pagination, fallback or discovery. Secret never supplied to child/argv/env/logs or read by the assistant.
- 15 irrevocable POST reservations, all exact model/native1024; 14 responses completed/forwarded. Exchange reached response validation for #15, where the repeat-handoff guard closed forwarding and terminated only the owned GUI. Upstream accounting is 1 GET200 +15 POST200 **inferred from the pinned exchange path** (it rejects non-200 before confirmation/response validation); separate per-exchange numeric status receipts were not recorded. Local rejection attempts503; delivery to the terminated GUI is not confirmed.
- Retained responses1–14 report input35,983/output1,018/cache-creation0/cache-read0; all four fields are present in all14. Response15 was rejected before raw-response storage, so its usage is unknown. These totals are LOWER BOUNDS, not complete billable totals; no credential re-read/API query was made to fill the gap.
- Six human turns completed canonically at23/36/49/62/75/88. Seven inputs were sent. Six complete ordinary file results (8051–8054bytes) reached requests3/5/7/9/11/13. Primary SYSTEM count1 across all15 requests; manual count0 then1, and manual dialogue count0 throughout.
- Actual L2 cuts25/38/51/64/77/90 and paired-tool cut91 preceded one recorded handoff99, then unload102. Saved manual is pending with its body hash intact; saved effective tools exclude read_file. Snapshot/context AsOf105, journal reaches107 with turn7 interrupted. Final requested answer/normal EOF success are not claimed. No root-cause/product classification or automatic repair follows this bounded stop.

Owned GUI PID12189/listener127.0.0.1:54645; PID no longer exists. Launcher exited after owned listener/process closure. No shared cleanup. Pins and both frozen R4 hashes still match (`live-01-after-pins.check`, `live-01-after-frozen.check`). Minimum final free disk18,395,340KiB.

Raw traces remain LOCAL under `/tmp/e3-ch14-student/sonnet46-preparation/live-01/`; no provider-body export. References/hashes: `local-artifact-references.json`, request/retained-response hashes in `safe-summary.json`. Also retained: `receipt.json`, `trace.jsonl`, `models-response.json`, raw requests/retained responses, GUI stdout/stderr, owned save/journal; outer `live-01.{command,status,stdout,stderr}`. Safe summary preserves numeric usage, qualifications, structural canonical facts and cleanup observation only. All older failures/control receipts remain intact.

HOLD for root decision. Production rounds remain4/10; no further live branch, real credential access, network diagnostics, browser/native action or operational continuation is authorized by this result.
