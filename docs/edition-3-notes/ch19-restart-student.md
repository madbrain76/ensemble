# Chapter 19 restarted student — R10 compiling/tested HOLD

**Production is HELD after bounded R10.** No edits/checks are running; await the same direct/retained review. The historical ten-round stopping limit is lifted, without changing review standards or this turn's authority.

## R10 delta and focused evidence

- Only production change: `internal/sandbox/sandbox.go`, in the existing Profile owner. Added metadata-only literals `/var`, `/private/var/select`, and `/private/var/select/developer_dir`, using the existing literal accumulator/deduplication. A brief comment explains lexical alias traversal and inspection of the developer-directory link into the already-readable Library runtime tree.
- Added `internal/sandbox/ch19_r10_test.go`. It requires the three metadata-only literals and rejects selection subtrees/data grants, global metadata/root subtree grants, and unrelated select/sh or gitattributes paths. It also checks the retained literal-root read, network denial, and immediate xcrun_db cache boundary. R9 and other owners are unchanged.
- No resolver, namespace/chroot, installer, host-selection, tempfile/environment, or warning-cleanup changes. The native observation supplied by root establishes sufficiency of the bounded combination, not every literal's necessity; this student made no native correctness claim.

Verified existing reports (commands and outcomes directly read):

| Report path | Command/outcome |
|---|---|
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/3016` | `go test -vet=off ./internal/sandbox -run '^TestCh19R10DeveloperSelectionMetadata$' -count=1 -timeout=30s`, exit 1: before test reported all three missing metadata literals. Pure profile-shape evidence, not a native failure repro. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/3018` | New R10 profile test plus existing `TestCh19R9LiteralRootRead`, exit 0. Both are pure string/profile checks; no PrepareCommand or system-tool exercise. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/3019` | Root then GUI `go test -vet=off ./... -run '^$' -count=1`, both exits 0. Compilation only, including the nested CLI and separate GUI module. |

All Go work used the established sequential env-i canonical owned HOME/TMP, readonly/p1, GOMAXPROCS=2, local toolchain, offline proxy/sum settings and existing caches. Scoped diff check was clean; no setup failures or unresolved blocker. No student native sandbox/PrepareCommand/Python/Git workload, full suite/race/vet/grade, provider/browser/install/restart, private fixture/reviewer/reference/root-receipt access, children, snapshot or commit. Root owns subsequent native release and validation. No extra pre-failure or sensitivity campaign ran.

---

## Retained R9 record

**Production is HELD after bounded R9.** No edits/checks are running; await the same direct/retained review. The latest brief lifts the historical ten-round cap, but does not widen this turn's authority. Architecture remains trusted core + rooted builtin filesystem operations + OS-confined commands with explicit runtime exceptions; no namespace/chroot or new temporary-file framework.

## R9 delta

- `internal/sandbox/sandbox.go`: Profile adds only `(allow file-read* (literal "/"))`. Existing ancestor metadata, bounded immediate physical xcrun_db cache grant, system/device references, environment/Git handling, network denial and descendant confinement are retained. No `/private/var/select/sh` grant or broader root/temp subtree grant was added.
- Sandbox validation/operation errors retain the supplied pathname while keeping filesystem causes (`errors.Is` etc.). Canonical containment and rooted operations are unchanged. Typed PathError/LinkError fields are reattributed at the owner, not globally string-scrubbed.
- `internal/tool/tool.go` / `files.go`: file tools pass supplied names into Sandbox rather than expanding them first. Write checks the full target before parent creation; recursive search keeps relative traversal/display names. Errors from opened-handle reads/writes retain the requested filename. ReadDir sorting and successful relative/absolute read/write/edit/list/search behavior are preserved. No read_file schema or safe-mode text changes.
- `internal/tool/dispatch.go` / `command.go`: cwd resolver refusals occur before job reservation, so they create no handle/output locator, publish no expanded cwd, and dispatch no process. Admitted jobs retain existing expanded cwd/output/handle reporting. Confined preparation still checks the original cwd before launch; ordinary unconfined cwd expansion remains at its existing boundary. Established unconfined missing-cwd startup failures remain truthful job failures, not relabeled as resolver refusals.

## Focused results and verified reports

Reports below were directly read and verified. Go work was sequential, env-i with canonical owned HOME/TMP, dummy/no real credentials, readonly/p1, GOMAXPROCS=2, local toolchain, offline proxy/sum settings and existing caches. Commands/outcomes are saved in the reports.

| Existing report path | Meaning |
|---|---|
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2878` | `go test -vet=off . ./internal/sandbox -run '^TestCh19R9' -count=1 -timeout=30s`, exit 1. Root fixture signatures failed compilation (setup only); Sandbox tests actually reproduced missing literal-root read, missing supplied spelling, and lost filesystem-error spelling. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2879` | Corrected root fixture, `go test -vet=off . -run '^TestCh19R9' -count=1 -timeout=30s`, exit 1. Real file-tool failures lost spelling; bad cwd reserved expanded cwd/output/handle context. These commands were aborted at reservation, never released. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2888` | First post-edit check, exit 1 at compilation: unused sort import exposed omitted ReadDir sorting. Sorting was restored; not behavioral evidence. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2889` | Post-edit focused run, exit 1: real missed edit_file target expansion still exposed synthesized names; Sandbox checks passed. Fixed that caller. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2906` | Final root/Sandbox focused run: TestCh19R9 family plus TestCh19PreparationPolicy, exit 0. Covers relative traversal, symlink/outside and supplied-absolute refusals, cause preservation, allowed relative/absolute file operations and sorted/relative display, rejected-cwd envelopes/no reservation/no dispatch, and admitted reservation metadata without starting confined commands. Literal-profile assertion rejects broad grants. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2910` | Existing TestChapter3Round1ShellEnvironmentAndResults, TestChapter4Round1Cwd, TestChapter4Round2TerminalSettingsAndParity, and TestPTYNormalizationReadBoundaries, exit 0. Existing unconfined shell key/output/exit, cwd/context, PTY and streaming assertions retained. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2912` | Root then GUI `go test -vet=off ./... -run '^$' -count=1`, both exits 0. Compilation only, including nested CLI and separate GUI module. |

Two edit scripts stopped on unmatched anchors before production writes: existing reports `io/2880` and `io/2882` under the same directory show AssertionError. Trailing inspection commands caused their outer exit 0; those are not successful edits/checks. No known unresolved blocker remains. Scoped diff check was clean; about 17 GiB free.

Only existing authorized unconfined local shell/PTY fixtures and the existing unexecuted preparation helper (including its trusted getconf lookup) ran. New confined command cases were inspected/aborted without launch. No native Seatbelt/host-escape, full suite/race/vet/grade, live provider/OAuth/browser/native GUI/install/restart, private grader/reference access, child, snapshot or commit. Profile/source assertions and local file checks are not universal profile or OS-process-confinement proof. Root owns native validation; its earlier shell startup observation and select/sh warning are not grounds for additional grants.

---

## Retained R8 record

R8 was held at this point. Its source/test record is retained below; current R9 status is above.

## Test-only HOLD: Chapter 3 diagnostic fixture (not R9)

R8 production remains HELD and unchanged. Only `TestChapter3EditDiagnostics` in `chapter03_tools_test.go` and this note were edited for the authorized fixture repair.

- Validate the exact leading `file %q:` identifier with `strings.CutPrefix`, then check the untouched diagnostic body for anchors/excerpts and required count/location/truncation/guidance. No-mutation, failure, overlap and all existing expected-value assertions remain intact.
- Every fixture now uses deterministic `aa-file`, reproducing the path/anchor collision independently of random TMPDIR suffixes. The same checker must reject injected body anchors, fixture excerpts, and a repeated full file prefix in the body; it cannot silently delete arbitrary path occurrences.
- Verified saved report: `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2771`. Command: `go test -vet=off . -run '^TestChapter3EditDiagnostics$' -count=1 -timeout=30s`, using the required isolated/offline environment and existing caches. PASS, exit 0. Report includes the actual command and outcome.
- Root's retained failure supplies before evidence; no new before run, full suite, broad compilation, race/vet or other gate was run. No production edits, snapshot/commit/child or unrelated repair. Await tiny test-repair review before root continuation.

## R8 delta and coverage

- Production change is confined to `internal/skills/parse.go`. The explicitly application-selected corpus now uses host ReadDir/Stat/ReadFile, as Chapter 19's separate trusted-input contract requires. The optional Sandbox parser argument is removed. Directory selection remains BaseDir-relative when relative; parsing, validation, source metadata and duplicate-name handling are unchanged. A comment explains the trust boundary and that later loads use retained declarations by name.
- Reviewed adjacent discovery, activation and restoration paths: they retain declared-name/primary/offered/dependency admission and do not turn later `load_skill` calls into filesystem reads. No activation, SandboxRoot, file/state/media/job/process I/O, MCP execution or auth behavior was changed. Corpus files are not copied or relocated.
- `chapter19_r8_corpus_test.go` uses real Agents/registries with an outside corpus, both absolute and BaseDir-relative. It verifies primary and offered manual/dependency activation, real admitted `load_skill` provenance for a declared name, refusal of unoffered/dependency/primary misuse and path-shaped requests (including a valid unselected document), disabled-web/dependency refusal and retained credential-name inventory. Skill names deliberately differ from folder names.
- Paired file-tool controls succeed on an inside file and refuse absolute/relative reads of the same outside primary document whose trusted instructions loaded. Root/corpus configuration and source paths remain unchanged. External duplicate/malformed metadata and directory-name-as-primary are refused; focused existing parser list/validation checks also pass.

## Verified R8 receipts

Root's preserved native startup finding supplies the before evidence. No additional pre-failure campaign or access to private grader/fixture/review/runtime files was performed. These existing reports were directly read and verified; Go work used the required sequential env-i, canonical owned HOME/TMP, offline/local toolchain and existing caches.

| Existing report path | Command/outcome |
|---|---|
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2742` | `go test -vet=off . ./internal/skills -run '^(TestCh19R8TrustedCorpusAndConfinedFiles\|TestCh19R8ExternalCorpusValidation\|TestFlexibleListsAndLiteralVariables\|TestMalformedMetadata)$' -count=1 -timeout=30s`; both packages pass, exit 0. The report contains the exact executable regex without Markdown table escaping. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2745` | Root then GUI `go test -vet=off ./... -run '^$' -count=1`; both exits 0. Compilation only, including nested CLI and separate GUI command package. |

No R8 setup/test failure or unresolved blocker is known. `git diff --check -- .` was clean; about 18 GiB remained free. The tests use local files/metadata and builtin tools only; no shell, provider or MCP connection was executed. No full suites/race/vet/grades, actual Seatbelt/host escape, native/browser/OAuth/install/restart, commits or snapshot changes were run. Local paired file checks and compilation do not establish native startup/grade success or OS-process confinement; those remain root-owned validation.

---

## Retained R7 record

R7 was held at this point. Its source/test record is retained below; current R8 status is above.

## R7 delta and coverage

- `startup_file.go` adds public synchronous `Ensemble.ReadStartupFile`. Core privately reuses the existing unstarted Agent/Sandbox construction, closes that Sandbox before every return, and exposes only bytes/existence, resolved file path, and selected DataDir/SavePath choices. It installs no Agent/transport/capabilities and exposes no bootstrap lifetime or callback. Missing files remain distinct from read/confinement/close errors; JSON parsing is not core's responsibility.
- `gui/server.go` no longer imports `ensemble/internal/agent`. Startup reads through the public seam, parses with existing `patchSettings`, copies only selected DataDir/SavePath into the original Config, applies saved safe/web policy, then calls ordinary NewAgent. The final store binds to the final Agent using the same parsed snapshot, without rereading.
- `gui/settings.go` shares the existing launch-policy defaults between preload and attached stores; schema, validation, explicit persisted precedence and standalone defaults remain GUI-owned. No unrelated production/fixture changes or weakening of the public-import test.
- New helper tests cover BaseDir-relative reads, inside/outside owned symlinks and refusals, missing versus opaque malformed bytes, explicit caller-relative SavePath, derived/empty-DataDir path conventions, no Agent registration/ID consumption, and Sandbox descriptor closure on success/missing/refused/failed-construction paths.
- New GUI tests use a test-only existing Auth seam to change an owned settings file between preload and normal construction. They verify one parsed snapshot, original auth selection (not wholesale bootstrap normalization), final-owner binding, relative settings resolution, and actual explicit caller SavePath placement. Empty/malformed/type-invalid settings are rejected before Agent installation. Existing R2/R3 startup/recovery policy, rooted settings and one-DataDir tests pass.

## Verified R7 receipts

Root supplied the valid before finding: GUI's unchanged public-import test rejected the private core import during its R6 GUI-race run. I did not rerun that before case or inspect root validation/runtime files. These student reports were directly read after completion; all Go work used the required sequential env-i, owned canonical HOME/TMP, offline/local toolchain and existing caches.

| Existing report path | Command/outcome |
|---|---|
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2687` | Initial focused run, aggregate exit 1: new root/GUI test fixtures incorrectly used nonexistent Config.UseStub and did not compile. This is setup failure, not behavioral evidence. In the same run, `ensemble/gui/cmd` ran the unchanged `TestChapter12GUIPublicImports` under the named selector and passed. That valid public-import result was reused after test-only corrections. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2688` | Root `go test -vet=off . -run '^TestCh19R7StartupFile' -count=1 -timeout=30s`, then GUI package `go test -vet=off .` selecting the existing R2 startup/default, R3 one-DataDir, rooted settings/policy, and new R7 snapshot/malformed cases (exact regex in report), `-count=1 -timeout=60s`; both exits 0. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2690` | Root then GUI `go test -vet=off ./... -run '^$' -count=1`; both exits 0. Compilation only, including nested CLI and GUI command package. |

`git diff --check -- .` was clean; about 18 GiB remained free. No unresolved R7 source/test blocker is known. No full suites/race/vet/grades, actual Seatbelt/host escape, real provider/OAuth, GUI/browser/native, install/restart, commits or snapshot changes were run. Local file/fake-peer tests and compilation are not live application or OS-confinement proof. Root's full validation remains stopped pending its review/release.

---

## Retained R6 record

R6 was held at this point. Its source/test record is retained below; current R7 status is above.

## R6 delta and focused coverage

- `internal/common/skill.go` adds the narrow `SkillRegistry.MCPAuthEnvNames()` query. `internal/skills/skills.go` implements it by reading already-known parsed declarations under the existing registry lock and returning sorted, unique nonempty environment names. It does not return credential values, cache a second registry, or grant/advertise capabilities.
- `internal/common/environment.go` uses that query instead of capability-filtered ListSkills. Provider/Config MCP exclusions, ordinary unconfined build settings, and strict confined handling are otherwise unchanged. ListSkills/Snapshot filtering, activation, and accepted earlier owners were not changed.
- `chapter19_r6_environment_test.go` constructs a real registry with a custom harmless auth variable declared only by web-search. It directly inspects ProcessEnvironment in initial-disabled, enabled-but-unloaded, and disabled-again states. It also asserts ordinary GOFLAGS/custom build settings survive, web visibility follows policy, and no unloaded web tool is enabled or exposed. It never activates the skill or invokes shell/MCP/network/provider paths.
- The existing R5 environment fake implements the new names-only query; its prior credential/build/confinement assertions remain. No production or test setup failure occurred in R6.

## Verified R6 receipts

These existing reports were directly read after completion; they contain the traced commands, outcomes and exits. All Go work used the required sequential env-i, canonical owned HOME/TMP, offline/local toolchain flags and existing caches. No checks were repeated for documentation.

| Existing report path | Command/outcome |
|---|---|
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2638` | `go test -vet=off . -run '^TestCh19R6HiddenSkillCredentialEnvironment$' -count=1 -timeout=30s`; exit 1. Genuine before: ProcessEnvironment included the dummy skill auth variable in both disabled states. The enabled-but-unloaded control did not fail. No shell or transport leakage was exercised. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2640` | `go test -vet=off . ./internal/common -run '^(TestCh19R6HiddenSkillCredentialEnvironment|TestCh19R5ProcessEnvironmentContexts)$' -count=1 -timeout=30s`; exit 0 in both packages. Real-registry and affected environment checks passed. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2641` | Sequential root, GUI, ch06, ch07, examples/ch04-consumer and examples/ch13-speech `go test -vet=off ./... -run '^$' -count=1`; all six exits 0. Compilation only, covering consumers of the registry interface change. |

`git diff --check -- .` was clean. No unresolved R6 source/test blocker is known. Full suites/race/vet/grades, actual Seatbelt/host escape, provider/OAuth, GUI/browser/native, install/restart, commits and snapshot changes remain unreleased and were not run. These local helper and compilation results do not establish live transport or OS-confinement behavior. Root should review the held delta and decide continuation of validation.

---

## Retained R5 record

R5 was held at this point. Its source/test record is retained below; current R6 status is above.

## R5 production delta

- `internal/agent/agent.go`: explicit SavePath is made absolute against caller cwd before Sandbox validation. Relative DataDir still derives from BaseDir. Outside caller-selected paths are refused under confinement rather than silently relocated. Existing composition-root derived paths and persistence ownership/overlap logic remain intact.
- `internal/common/environment.go`: restored the accepted Chapter 18 unconfined contract—retain ordinary environment/build settings, exclude supported provider credentials and Config/declared-skill MCP AuthEnv variables, force TERM=dumb. A nonempty SandboxRoot uses selected locale/timezone values and fixed system PATH; no arbitrary build settings, startup hooks, inherited HOME/TMP/Git overrides, or credential variables pass through that branch. No generalized secret registry.
- `internal/sandbox/sandbox.go`: root explicitly approved completing the previously missing Git preparation within R5. PrepareCommand now generates a fixed identity-only Git file using existing rooted CreateTemp/Rename operations, atomically replaces it, and sets rooted GIT_CONFIG_GLOBAL, GIT_CONFIG_SYSTEM=/dev/null, and GIT_CONFIG_NOSYSTEM=1. Existing rooted HOME/TMP and confinement preparation remain. No Git command or Seatbelt process was launched.
- Accepted credential/FIFO/actor/startup owners were not changed.

## Test-only alignments and focused coverage

- Chapter 10 assertions now compare the resolved DataDir and canonical caller-selected explicit SavePath; root isolation, initial events, reopen, cwd, same-path recovery, and primary/backup overlap assertions remain.
- Chapter 4 file streaming and killed observation use absolute Output.Locator directly (join BaseDir only for relative locators). Exact bytes/order/context/footer/killed assertions remain. Chapter 3's original shell key/output/build-setting/exit assertions were not altered.
- Chapter 18's intended enabled public-web fixture explicitly enables web search; its fake-peer auth, unloaded/enabled dispatch and unload assertions remain. Human model/status parses and asserts each safe/web/root prefix before JSON while retaining model/accounting/session checks.
- Disk history tests now use a real unstarted internal Agent with its parent-linked rooted Sandbox; diagnosticParent delegates Sandbox while retaining its logger override. Memory-only New(nil) tests remain. No production nil-owner filesystem fallback. Persistence, invalid candidates, backups, journal, concurrency, recovery and diagnostic assertions remain; malformed-JSON output stays inside the fixture root so a path refusal cannot mask that assertion.
- Added focused explicit-relative SavePath/inside-vs-outside confinement and unconfined/confined environment checks. Preparation assertions cover fixed PATH, rooted HOME/TMP, generated identity-only Git content, and exclusion of inherited Git includes/helpers/rewrites/count injections. Preparation inspects files/environment without starting the prepared target.

## Verified R5 receipts

All paths below were directly read after completion. Each contains the exact traced Go command(s), outcomes and exits. All Go commands used the required sequential env-i, canonical owned HOME/TMP, offline toolchain/module settings and existing caches. No evidence bundle or test rerun was made for documentation.

| Existing report path | Command/outcome meaning |
|---|---|
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2567` | Initial test-construction/edit-anchor errors; Go build failed on nonexistent Agent.Load and an embedded SkillInfo field literal. Exit 1; no behavioral before credit. Corrected against actual interfaces. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2569` | `go test -vet=off . ./internal/common -run '^(TestCh19R5|TestChapter3Round1ShellEnvironmentAndResults)' -count=1 -timeout=45s`; exit 1. Genuine before: original shell lost build settings, explicit relative saves relocated in both contexts/outside path silently moved inside, and confined environment inherited PATH. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2574` | Traced named owner tests and nine selected disk-history tests. Root/common/CLI cases and selected history tests passed. Overall exit 1 solely for missing confined Git variables; this is retained meaningful before evidence for the explicitly approved Git addition. The selector's literal `TestCh10R2` did NOT select its two longer test names; those were run separately below. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2583` | `go test -vet=off . -run '^(TestCh10R2BaseRelativeRecoveryIsolation|TestCh10R2ExplicitRecoveryPathAndOverlap)$' -count=1 -timeout=45s`; exit 0. Fills the selector gap, not a rerun for receipts. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2585` | First Git insertion failed compilation because Sandbox has no WriteFile method. Exit 1; no behavioral result. Corrected using its actual rooted temp/rename API. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2587` | `go test -vet=off ./internal/sandbox -run '^TestCh19PreparationPolicy$' -count=1 -timeout=30s`; exit 0. Corrected Git/file/environment preparation checks passed; no prepared command launch. |
| `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2588` | Sequential root and GUI `go test -vet=off ./... -run '^$' -count=1`; both exits 0. Compilation only, including importable root, nested CLI, and separate GUI module. |

`git diff --check -- .` was clean. About 19 GiB remained free at the last disk check. Earlier passing R5 results were reused where subsequent changes did not affect them.

## R5 limits / next step

All released focused checks have passing coverage across the reports above; there is no unresolved R5 source/test blocker known here. Existing local shell/job fixtures and loopback peers ran only under the released dummy/owned environment. No full suite, race, vet, grader, actual Seatbelt/host-escape scenario, real stdio MCP, Git command, live provider/OAuth, GUI/browser/native operation, install/restart, commit, or snapshot change occurred. Preparation and compilation are not OS-confinement proof. Root should review the held delta and decide continuation of its isolated full validation.

---

## Retained R4 record

R4 was held at this point. Its source/test record is retained below; the current R5 status is above.

## R4 delta

- `internal/engine/credentials.go` now normalizes the exchange's exact literal/standard-JSON spellings once: longest-first with deterministic ties and duplicates removed. That same list drives the existing incremental matcher and a standard-library `strings.Replacer` shared by complete Parts/JSON/event response fields and delta Name. Replacement is one-pass, so inserted `[REDACTED]` markers are not processed as further tokens. No global `common.Redact` change, heuristic detector, registry, or custom replacement framework.
- `internal/engine/engine.go` passes the shared replacer to final publication. Error fields retain their existing additional exact diagnostic values using the same normalized one-pass builder before callbacks; the old second event-message scrub was removed. Renewal-error text also uses that helper. Actor/FIFO logic and the approved conditional-delay tradeoff are unchanged. No GUI/startup/DataDir changes.
- Added `chapter19_r4_credentials_test.go`; adjusted the existing raw-field helper regression to the shared replacer API. Fake Auth renewal checks overlapping current/renewed tokens in plain and SSE responses, text deltas versus OnEvent and returned content, complete/delta tool names, CallID/nonsecret argument preservation, and a dummy old value occurring inside the marker (to reject sequential marker reprocessing).

## Verified existing reports

Reports below were read directly and contain the stated commands/outcomes. All Go commands were sequential under the required `env -i` owned canonical HOME/TMP and offline settings: GOMAXPROCS=2, GOFLAGS='-mod=readonly -p=1', GOTOOLCHAIN=local, existing caches, GOPROXY/GOSUMDB off.

1. `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2465` — `go test -vet=off . -run '^TestCh19R4' -count=1 -timeout=30s`, exit 1. Genuine before behavior in all four subcases: final/OnEvent text and delta Name/complete tool fields kept the renewed suffix while streamed chunks replaced the whole value; the marker-overlap case also showed marker reprocessing. This is exact-content inconsistency, NOT demonstrated full-token disclosure or a live-provider finding. No fixture/setup failure.
2. `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2471` — `go test -vet=off . ./internal/engine -run '^TestCh19R[234]' -count=1 -timeout=30s`, exit 0 for root and engine. Includes new overlap cases and affected current/renewed, JSON-byte preservation, ordinary/split-secret/interleaved streaming, actual actor arrival-policy, and terminal/error/cancellation regressions. No post-fix setup failures.
3. `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2474` — `go test -vet=off ./... -run '^$'` in root, exit 0; importable root and nested CLI compile. Compilation only, no tests run. No unaffected GUI campaign or other-module reruns; earlier valid unchanged receipts remain below, not represented as fresh R4 results.

`git diff --check -- .` clean; last disk check ~19 GiB free. Receipt-reading/note-writing reports are not themselves test evidence. No new evidence bundle, commit, cleanup, or unrelated edit.

## Limits / next step

Stable HOLD for same-reviewer inspection; no unresolved setup blocker. Dummy Auth and loopback/in-process fixtures only. No full suites, race/vet/grades, Seatbelt/sandbox-exec, host probes, real stdio process, browser/native GUI, provider/live OAuth, install/restart, or unrelated process interaction. Local consistency tests and compilation are not OS-confinement or live-provider proof. R3's explicitly approved conditional buffering remains unchanged and documented in the retained record.

---

## Retained R3 record

R3 was held at this point; the current R4 status is above. Same coder/workspace and retained reviewer. This historical record does not release broader gates.

## R3 delta and agreed tradeoff

- `internal/engine/credentials.go`: replaced independent per-part publication with a request-local conditional FIFO. Undecided bytes are matched within their original part, then resolved into their original delta slots/metadata. Ready fragments drain in original cross-part arrival order. A split credential is replaced at its first byte's original position; later consumed fragments retain their metadata. No actor or explanation-policy code changed, no final-part-order inference, no extra actor signal or credential registry.
- The parent explicitly approved this necessary tradeoff: ordinary unambiguous streaming drains immediately, but later payload waits behind an unresolved exact-token-prefix ambiguity. If that earlier part never resumes, later payload necessarily waits until parse completion, including error/cancellation completion. The queue is request-owned and flushed by the existing Exchange completion path; no blanket response buffering, arbitrary size/failure policy, or background worker was added.
- `gui/server.go`: carries bootstrap's selected normalized DataDir into ordinary Agent construction. If both original DataDir and SavePath were empty, it also carries the bootstrap save.json path so supplying the selected DataDir cannot accidentally switch to root's explicit-DataDir/per-Agent SavePath convention. Explicit paths, existing rooted parser/saved-policy precedence, and Sandbox-handle closure remain intact.
- Tests added: `chapter19_r3_stream_test.go` and `gui/ch19_datadir_test.go`. Fake HTTP/SSE and an in-process counter-only tool exercise actual Engine → actor → dispatch paths without process/host effects. Before evidence demonstrates both policy inversions, not merely a display-order discrepancy.

## Exact verified receipts / commands

All Go work used the authorized sequential `env -i` dummy HOME/TMP/PATH and offline Go environment: GOMAXPROCS=2, GOFLAGS='-mod=readonly -p=1', GOTOOLCHAIN=local, existing caches, GOPROXY/GOSUMDB off. The following existing on-disk reports were read directly; each contains its actual command lines and outcomes. The separate receipt-reading/note-writing reports are not test evidence.

1. **Before failures:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2428` — root then GUI: `go test -vet=off . -run '^TestCh19R3' -count=1 -timeout=30s`; both exit 1. Delta order/metadata changed in both directions. Real actor test wrongly refused text-first (`effects=0`, want 1) and wrongly admitted tool-first (`effects=1`, want 0). GUI helper allocated two temporary data roots. These were product failures, not fixture/setup errors.
2. **Focused after + affected R2 regressions:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2431` — root: `go test -vet=off . ./internal/engine -run '^TestCh19R[23]' -count=1 -timeout=30s`; GUI: `go test -vet=off . -run '^TestCh19R[23]' -count=1 -timeout=30s`; all exit 0. Covers original cross-part positions/metadata, both actual actor-policy outcomes, interleaved split-secret non-disclosure in both directions, final tool names/arguments, ordinary early streaming, current/renewed echoes, raw JSON/empty-delta preservation, non-resuming-prefix terminal/error/cancellation flush, one selected DataDir with unchanged default SavePath, and affected rooted saved-policy/default/recovery behavior. No post-fix test setup failures.
3. **Affected-module compilation only:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2434` — sequential root then GUI: `go test -vet=off ./... -run '^$'`; both exit 0, including nested CLI and GUI command. No tests run in this gate. Unchanged other-module compile results remain the R2 receipts below; they were not rerun for documentation or represented as fresh R3 results.

`git diff --check -- .` clean. Last disk check ~19 GiB free. No commits, cleanup, new evidence bundles, or unrelated edits. No new blocker or unresolved implementation ambiguity.

## Limits / next step

Stable HOLD for same-reviewer inspection. Tests cover local Chat Completions SSE and fake Auth/transports, not every provider dialect or a live service. No full suites, race/vet/grades, sandbox-exec/Seatbelt, host escape probes, real stdio process, browser/native GUI, provider/live OAuth, install/restart, or unrelated process interaction. Conditional buffering is documented above rather than hidden behind a claim of unconditional incremental delivery. Compilation and local fixture success are not OS-confinement proof. No broader gate is released.

---

## Retained R2 record

R2 was held at this point; the current R3 status is above. Same workspace/coder, accepted Chapter 18 base, and retained reviewer. This historical record does not release broader validation.

## R2 changes

- **Persisted policy before startup:** `gui/server.go` adds the launcher-specific `NewServerFromConfig`. A short-lived, unstarted internal Agent supplies the existing parent-linked Sandbox to the existing settings parser. It has no installed History/skills/MCP/Engine or root identity; its Sandbox is closed on success and failure. Parsed safe/web values enter Config before ordinary `root.NewAgent`, so recovery and startup capabilities see that policy. The same parsed settings/store snapshot is then attached to the real Agent and Hub. `gui/cmd/main.go` uses this path for GUI launch; protocol mode retains ordinary Agent construction. Asset validation and rooted settings path checks remain. No optional Config, second parser, public lifecycle hook, registry, or GUI default root was added.
- **Exact transport echoes before public publication:** `internal/engine/engine.go` retains current/renewed resolved tokens only for the exchange. Final events are sanitized before `OnEvent` and return; `OnDelta` is filtered incrementally. `internal/engine/credentials.go` edits known string values and their standard JSON spellings, not key-shaped patterns or unrelated environment values. Raw JSON fields retain unrelated bytes (including number spellings and whitespace); only strings containing a known value are replaced. Delta tails hold only a possible credential suffix per part and flush at completion, retaining per-part append order; ordinary streaming continues before response completion. Renewal-error events and request traces use the known values too. Existing frame diagnostic protection remains.

## Focused coverage and actual receipts

All Go commands were sequential under the authorized `env -i` dummy HOME/TMP/PATH and offline Go settings (GOMAXPROCS=2, GOFLAGS='-mod=readonly -p=1', local toolchain, existing caches, GOPROXY/GOSUMDB off). No real credentials were used. Existing on-disk reports below were directly read and verified; final reports contain the command lines and exits.

- **Meaningful before:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2329`. Root then GUI `go test -vet=off . -run '^TestCh19R2' -count=1 -timeout=30s`; both exit 1. Valid fake responses leaked current/renewed tokens through returned parts, OnEvent, and deltas, including a token split across SSE chunks. GUI startup made one web handshake before saved web-off applied; disabled recovered/requested web was not refused under saved policy. The same output also has an incorrect test expectation about existing argument compaction: that expectation was corrected, and is NOT credited as a production defect. The safe-initial history assertion was strengthened later; this before report does not establish a separate safe-initial failure.
- **Intermediate compile/edit failure:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2336`, exit 1. An exact-anchor edit assertion stopped; two incorrect Part field references failed compilation. No behavioral result from that run. Corrected without widening scope.
- **First post-fix focused pass:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2338`, root/GUI targeted tests both exit 0. Later edge additions are covered by the final report below.
- **Final focused checks:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2341`, exit 0. Root command: `go test -vet=off . ./internal/engine -run '^TestCh19R2' -count=1 -timeout=30s`; GUI command: `go test -vet=off . -run '^TestCh19R2' -count=1 -timeout=30s`. Covers plain/current and renewal response publication, OnEvent/OnDelta/final events, split-token SSE with ordinary first-chunk delivery before completion, nonsecret tool arguments, unrelated environment-marker preservation in public events, renewal-error echoes, exact escaped JSON strings with unchanged numeric/whitespace bytes, and empty delta notifications. GUI covers zero startup handshakes under saved web-off, stricter recovery refusal, saved safe mode before initial tool declarations, missing/legacy deliberate launch policy, and outside-symlink settings refusal before Agent installation.
- **Six-module compile attempt:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2342`, aggregate exit 1. Sequential `go test -vet=off ./... -run '^$'`: root, ch06, ch07, examples/ch04-consumer, examples/ch13-speech all exit 0. GUI alone failed on a duplicate `err` declaration in its launcher. This is not a six-module passing receipt.
- **GUI compile correction:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2348`, exit 0, `go test -vet=off ./... -run '^$'` in GUI. Both GUI library and command compile after removing the duplicate declaration. The other five unchanged module passes from 2342 are reused; they were not rerun. Together these are passing compilation coverage for all six modules, not behavioral test runs.

`git diff --check -- .` is clean. Approximately 19 GiB free at the last disk check. No independent evidence bundle, commits, or cleanup was created. The receipt-verification/note-write tool reports are not themselves test receipts.

## Limits / next step

No unresolved setup blocker. Await direct/retained source review of the held R2 delta. Tests use dummy AuthService values, in-process/loopback fake peers, and owned temporary roots. Streaming fixture coverage is Chat Completions SSE, not live-provider compatibility or every provider stream dialect. No full suites, race/vet/grades, Seatbelt launch, real stdio process, host escape probe, provider/live OAuth, browser/native GUI, install/restart, or unrelated process interaction occurred. Compilation/local tests are not OS-confinement proof. No further work is authorized by this HOLD.

---

## Retained R1 record

Base: accepted Chapter 18 HEAD `08cef6aaf0c1472f75589e7016d2f52599ee9d0d`.
Workspace: `solutions/edition-3/main`. R1 was held at this point; the current R2 status is above. This retained record is not acceptance or OS-confinement proof. No new architecture, AgentSpec, fleet/spawning, or credential store.

### R1 delta (retained)

1. **GUI policy and controls:** `gui/settings.go` loads through the attached Agent Sandbox. Missing/legacy fields inherit deliberate launch SafeMode/EnableWebSearch values; explicit persisted fields and user patches remain authoritative. Both fields have typed patch handling. `gui/web/gui/{index.html,app.js}` contain actual controls and roundtrip bindings. `gui/hub.go` adds a settings-update mutex: policy setters run outside Hub.mu because synchronous capability callbacks re-enter that lock. Settings updates remain serialized.
2. **Retained job output:** `internal/jobs/jobs.go` observation, pattern scanning, and completion reopen output through the owning Sandbox. Tests replace the output with an outside symlink and verify refusal on all three actual methods.
3. **Profile/process boundary:** `internal/sandbox` grants metadata only on literal necessary ancestors, not globally; the immediate physical `xcrun_db*` cache regex remains. PrepareCommand confines independently authorized processes without deciding arbitrary-command permission. `internal/tool/command.go` explicitly performs ValidateCommand before preparation. Safe arbitrary commands remain forbidden; authorized skill MCP uses the existing prepared dial path.
4. **Web transitions:** discovery retains the web skill declaration. Config/Snapshot/load/dependency policy still withholds it when disabled. `Skills.RefreshPolicy` uses existing pending-unload state, publishes filtered declarations, reconciles MCP connections, and notifies existing subscribers. SetSecurityPolicy calls it after releasing configMu. Re-enable exposes the declared skill for explicit reload, rather than silently dialing a remote service. Loopback test verifies initial-disabled enable/load, successful tool call, disable/refusal, and a new handshake on re-enable/reload.
5. **Handle ownership:** deferred construction guards close Sandbox on every failure after acquisition in both `internal/agent.New` and root `Ensemble.NewAgent`; successful construction transfers ownership to shutdown. No finalizer-dependent cleanup added.
6. **Credential diagnostics:** Engine sanitizes frame callbacks/traces using the actual AuthService credential token, including renewal tokens, before diagnostic publication. Error messages are also sanitized. Existing MCP escaped-JSON redaction moved into the small shared `common.RedactDiagnostic` helper; MCP reuses it. CredentialValues additionally includes Config API keys and explicitly named Config MCP auth environment values. No resolved token is stored in Config or a new registry.
7. Added bounded tests in `chapter19_r1_test.go`, `chapter19_state_test.go`, `internal/common/ch19_policy_test.go`, `internal/jobs/ch19_output_test.go`, `internal/sandbox/ch19_test.go`, and `gui/ch19_policy_test.go`. Existing MCP loopback credential-error tests were reused.

R0 owners remain: public Config Child(dataDir)/Clamp(parent), Agent→Sandbox, rooted file/history/journal/cache/media/diagnostics, application defaults/flags, and shared safe/web policy. Importable root, nested cmd/ch01, and separate GUI module remain intact.

## Commands, outcomes, exact existing receipts

All Go work was sequential/offline, using this environment (no inherited real keys or endpoint):

```sh
env -i HOME=/private/tmp/e3-ch19-restart-student/home \
 TMPDIR=/private/tmp/e3-ch19-restart-student/tmp \
 PATH=/usr/local/go/bin:/Users/bill/go/bin:/usr/bin:/bin \
 GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local \
 GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod \
 GOCACHE=/Users/bill/Library/Caches/go-build <command>
```

The following are existing tool report files containing output and exit codes, not new evidence bundles. Receipt-only verification after HOLD directly reread these files; no code changed and no tests were rerun. **Canonical map:** 2150 = pre-fix local failures; 2181 = pre-fix credential failure; 2188 = mixed first post-fix run; 2193 = corrected profile/GUI pass; 2220 = additional run with descriptor-test setup failure; 2221 = corrected descriptor pass; 2223 = six-module compilation. Disregard any alternate mapping: `io/2159` and `io/2182` do not exist, `io/2186` contains a source read, `io/2224` is only a receipt-location inventory, and `io/2226` is only the note-writing result. None of those is a test/compile receipt. The commands below identify the original invocations; the linked report files retain their output and exits.

- **Pre-fix local regressions, exit 1:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2150`. `go test -vet=off ./internal/jobs ./internal/sandbox -run '^TestCh19' -count=1`, followed sequentially by GUI `go test -vet=off . -run '^TestCh19' -count=1`. Reproduced all three job-output disclosures, rejected safe authorized command preparation, missing/legacy GUI policy resets, and GUI outside-symlink settings read.
- **Pre-fix resolved-token leak, exit 1:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2181`. Root `go test -vet=off . -run '^TestCh19ResolvedCredentialDiagnostics$' -count=1`; valid fake loopback response leaked the profile-derived escaped token in the frame diagnostic.
- **First post-fix focused group, overall exit 1:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2188`. Root `go test -vet=off . ./internal/jobs ./internal/sandbox -run '^TestCh19' -count=1`, then GUI equivalent. Root credential/web-transition and all job tests passed. Profile test caught a missed metadata edit. GUI hung on the synchronous Hub callback lock cycle. Only owned test PID 36415 was sent SIGQUIT; its stack is retained in this report. Both product issues were corrected within R1.
- **Corrected profile and GUI, exit 0:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2193`. `go test -vet=off ./internal/sandbox -run '^TestCh19' -count=1 -timeout=30s`, then GUI `go test -vet=off . -run '^TestCh19' -count=1 -timeout=30s`. Both passed.
- **Additional bounded cases, overall exit 1 for test setup:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2220`. `go test -vet=off . ./internal/common ./internal/mcp -run '^(TestCh19(RootedStateAndMedia|FailedConstructionClosesSandbox|ChildClamp)|TestHTTPChapter18R1RemoteErrors|TestHTTPChapter18R2TruncatedUnicodeBearer)$' -count=1 -timeout=45s`. Config child/clamp and reused MCP credential tests passed. The rooted-state/media case was selected and not listed as failing; this nonverbose report is **not a standalone passing root-package receipt**. The root package failed with `lstat /dev/fd/3: bad file descriptor` in the descriptor test's os.ReadDir setup; this is **not** leak evidence.
- **Corrected descriptor test, exit 0:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2221`. Root `go test -vet=off . -run '^TestCh19FailedConstructionClosesSandbox$' -count=1 -timeout=30s`. Name-only descriptor enumeration verified repeated internal/root construction failures do not retain handles; successful shutdown closes the owned root.
- **Six-module compilation, exit 0:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2223`. Sequentially in `.`, `gui`, `ch06`, `ch07`, `examples/ch04-consumer`, `examples/ch13-speech`: `go test -vet=off ./... -run '^$'`. Each module exit marker is 0; all are compilation only.
- **Read-only final inventory/whitespace check, exit 0:** `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2222`, including `git --no-pager diff --check -- .`.

Other setup mistakes during R1: checked edit anchors stopped before applying later portions; early new tests used incorrect helper/type/API names and one response fixture lacked required usage data. These were corrected and are not negative-test credit. Root reported one undispatched tool-call infrastructure interruption; resumed the same files and reused completed checks without duplicating the worker.

R0 compile receipts remain `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2036` (root) and `/Users/bill/projects/cadence-memories/children/e3-restart-ch19-student/io/2039` (other modules).

## Limits / next step

- Focused tests used owned temp roots, dummy credentials, and loopback/in-process fixtures only. Command preparation was inspected without launching its target; its existing trusted getconf cache-location lookup ran. No sandbox-exec, arbitrary shell/job execution, stdio MCP child launch, OS-confinement/escape probe, live provider/OAuth/account, external network, GUI/browser/native action, install/restart, full suite, race/vet/grader, commit/push, or unrelated cleanup.
- Visible GUI controls were source-reviewed and backend-tested, not browser-tested. The profile's syntax/behavior, real toolchain startup, and actual confined MCP launch remain unverified until released. Compilation and local path checks do not prove Seatbelt enforcement.
- The descriptor check observes current-process descriptor names only; it does not read descriptor contents. Only the owned hung test was signaled; root GUI and unrelated processes were untouched.
- Last disk check remained about 20 GiB free. No genuine blocker or new consequential ambiguity remains for this bounded R1. Review the frozen delta before further production changes or broader gates.

---

## Final completion and chapter-feedback addendum

Production R10 remains frozen. Root reports that local validation and a bounded live CLI exercise completed; final retained review is being obtained. No further implementation or validation was performed by this student for this addendum.

Public observations supplied by root: useful file work, Python, and Git succeeded; escaped-file, symlink, outside-write, subprocess, and network attempts were refused. Synthetic-token filtering worked; narrower safe-mode file behavior and absence of execution declarations were observed; the CLI ended normally at EOF. Exact actual write/read results matched. Root controller budget and arbitrary-newline assertion failures are not student product defects.

Teaching feedback worth retaining:
- Keep the distinct trust boundaries explicit: application-selected skill files are trusted input; ordinary file/state/media/job operations stay rooted; command descendants use the OS profile with narrow runtime exceptions. This is not a chroot or a claim that the trusted core has no host access.
- Make application policy effective before recovery or capability activation, while respecting the public GUI/core import boundary. Keep hidden capability metadata available for credential exclusion without re-enabling the capability.
- Streaming credential protection must preserve exact content and arrival evidence, not just final text. Diagnostics likewise need to distinguish required caller-supplied file identification from prohibited content echoes.
- Apple runtime alias traversal is a useful concrete lesson: narrowly justified metadata literals can be necessary even when the final runtime tree is already readable. Warnings alone do not justify additional grants, and a sufficient combination is not proof that every element is necessary.
- Separate product defects from fixture/controller assumptions, and distinguish compilation, local helper tests, bounded live observations, and broader security claims.

These are bounded observations, not proof of universal security, all possible escape resistance, or every runtime/tool configuration. No OAuth or browser success is claimed. Root's live authority, correction-limit exceptions, final review, and any later release decisions remain root-owned; they require no student action. This addendum is the only final change, in the authorized student note.
