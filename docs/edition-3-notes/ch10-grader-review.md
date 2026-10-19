# Chapter 10 thin-adapter independent review — qualified, not accepted

## Disposition

**NOT ACCEPTED / incomplete control evidence.** No concrete defect in the frozen two-file adapter was established by source inspection or the completed Chapter 10 run. This is not an approval: the required Chapter 9 regression failed, the external Build-marker self-check was omitted, and the required copied-adapter route-sensitivity/restoration sequence did not run.

Root disposition (`main-20`): seal this attempt with all failures and unrun controls explicit; do not rerun, mutate, or repair source/harness. Root will inspect the failed Chapter 9 composition before deciding whether to authorize a bounded follow-up. The omitted Build-marker check and incorrect unexecuted mutation plan are **reviewer setup/scope failures, not student defects**.

No Chapter 10 student launch, public-guidance update, historical score change, or later work is authorized by this review. Prior entry/follow-up acceptance and both original sealed packages remain unchanged.

## Frozen subject, prerequisites, and scope

Task `/tmp/cadence-ch10-grader-review-task.txt`: 5,783 bytes, SHA256 `c1f4e6bdb7925df02e597adf083acd50f01d6ca2afbf4a5560efc9db50af58b1`. The governing private review input is `docs/edition-3-notes/ch10-grader-review-input.json`; its exact hash is retained in `task-and-manifest.json`.

Only the following candidate files were under review; the reviewer did not edit either:

| File | Bytes | SHA256 |
|---|---:|---|
| `cmd/grade-edition3/main.go` | 6,371 | `24fd6236dc3d885221c4f35ff1460d0036bd19aebb3626f658a01b1af743a024` |
| `cmd/grade-edition3/main_test.go` | 11,509 | `f69c2b3b91426b2223c05b23e61a51eede8f13461f3cebb811ee1e01e6f0af51` |

The exact candidate diff matched the frozen root diff. Historical `internal/gradeedition2/ch10.go` remains 20,693 bytes / `c7a63c5dc97a859f8d96f48f508edb0daad597dadbfb9a17dcd657c9240b1932`; no historical grader/reference/application mutation occurred.

Before and after controls, exact verification covered all **24 review inputs, 29 protected tracked files, 132 predecessor files, and 80 historical reference files**, plus both prior closed inventories and every listed file. Existing governing/historical full reads were carried only after exact hash comparison; the current reviewer role, frozen source/tests, exact diff, candidate/root gate receipts, acceptance record, and relevant adjacent validation/build/report routines were read. Truncated combined displays were not relied upon as complete; the full two candidate files and needed parser sections were read separately. The task was reread when reconciling the reviewer's mistaken control plan.

The accepted predecessor remains commit `251ce88fb1677897dbcc6b33661096c5cb39cd8e`, tree `a2add41c306d567c05d98ae760911bc13c18b75a`, clean. Its examples are `examples/ch04-consumer`, `examples/ch06`, and `examples/ch07`. No student behavior was graded or student source changed; the pinned consumer case was read as an existing declaration.

## Source review findings and limits

The actual diff adds 10 to the existing supported-chapter, grouped validation, help, and diagnostic lists; adds a dispatch branch calling unchanged `gradeedition2.Chapter10`; and extends the shared GUI-chapter validation tests. It does not replace Chapter 9 dispatch, change historical score weights, or change grader timing. Root's earlier prewrite assertions assumed a different source structure; those assertions failed before writing and are not evidence of source drift. Their reconciliation was reviewed as such.

The frozen source establishes the following:

- `flag.Visit` tracks presence independently of values, so unsupported earlier-chapter options remain rejectable even when explicitly empty.
- The supported chapter/arity guard rejects 11, missing targets, empty targets, and extra targets; no historical source default is added.
- Required core, consumer, workflow, stream, and GUI paths must have a regular `go.mod`; required stream/GUI models are checked with `TrimSpace`.
- `ReadConsumerCase` and `ReadMediaModels` execute before `grade.Build`. Their existing parsers reject malformed/missing supplied files and incomplete required data; media declarations require exactly the three supported vendor keys with nonblank identifiers.
- Help returns before target build/dispatch. Chapter 10 uses the same explicit executable/source/GUI/consumer/workflow/stream identities as the existing grouped routes.
- Setup/parser/build/report-write errors return 2; a well-formed failed grade report returns 1; a passing report returns 0. No changed unconditional-success return was found.
- Earlier 1/4/5/6/7/8 dispatch branches were inspected, not given fresh full/lifetime campaigns. Source review is not a substitute for unperformed runtime controls.

No independently established candidate code defect is alleged. Conversely, source ordering alone does **not** satisfy the task's independent no-Build instrumentation requirement, and a baseline identity assertion alone does not prove sensitivity to the wrong Chapter 10 route.

## Actual environment and artifacts

Every recorded command has exact argv, cwd, executable identity, PID, timing, deadline, exit, output filenames, and environment reference in `command-receipts.json`. Every output and receipt is hashed in the portable inventory. The clean environment is recorded in `environment.json`, not reconstructed from the invoking shell. It uses owned HOME/TMPDIR/GOPATH, the existing explicit Go caches, `/Users/bill/go/bin` before Go/OS paths, `GOMAXPROCS=2`, `GOFLAGS=-mod=readonly -p=1`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, and `GOWORK=off`. No inherited generic/provider endpoint or credential is included; real clients use the unchanged grader's local fakes and dummy keys.

Same-allowlist `command -v` resolves `/Users/bill/go/bin/dlv` and `/usr/local/go/bin/go`; `dlv version` succeeds. Delve matches 20,038,402 bytes / SHA256 `10dbb1ce8e6d2cb77ae4d30c7ca1db0302b83529cd6cd39453f009a7653c17fe`.

The reviewer built a binary once from the exact frozen source. It is byte-identical to the retained root candidate `/tmp/cadence-ch10-grader`: **10,934,098 bytes**, SHA256 **`e1262d3aff21be10d08426bf8fca58fc9e9178a114f3556dec5fa8860cd67fa5`**. The recorded controls invoked this reviewer-owned copy. This is an observed equality, not a reproducible-build claim. Root gates were already available; the reviewer additionally ran the exact package race test/vet/build, not a full-suite or lifetime campaign. This unnecessary additional local work is disclosed rather than represented as required by the finite full-route assignment.

Both full commands explicitly selected historical `solutions/ch10` as target and core, its separate `gui`, `ch05`, `ch06`, and `ch07` modules, and the same pinned accepted consumer case. The recorded three-vendor map is Anthropic `claude-sonnet-4-6`, OpenAI `gpt-5.4-mini`, Gemini `gemini-3.8-flash`; the declared stream and GUI model is `claude-sonnet-4-6`. Chapter 9 used these same compatible historical artifacts to isolate routing. No library-root/student identity was substituted for an executable.

## Completed controls and actual results

| Control | Actual result | Qualification |
|---|---|---|
| Same-allowlist resolution / Delve version | Exit 0 / 0 | Pinned installed tools; no install/discovery |
| Exact adapter package `go test -race -count=1` | Exit 0, 2.303 s | Package output reports tests passed |
| Exact adapter package `go vet` | Exit 0, 0.259 s | No diagnostics |
| Reviewer adapter build | Exit 0, 0.863 s | Binary exactly matches root candidate |
| Sentinel build and direct invocation | Build 0; target exit 99 | Real compiled marker created with exact expected content; marker then removed only in owned scratch |
| 80 external invalid/help cases | All expected exits, no launch marker | **Launch-refusal evidence only; no fake-Go Build marker** |
| Chapter 10 full adapter | **Exit 0; true; 100/100**, 16.704 s | Exact title and seven IDs checked; 85 persistence + 15 parity |
| Chapter 9 route regression | **Exit 1; false; 85/100**, 15.553 s | Genuine unexpected baseline failure; stopped here |
| Required copied-adapter route mutation | **Not run** | No source copy/mutation/rebuild/detector-failure evidence |
| Required exact restoration/rebuild/pass | **Not run** | No mutation occurred to restore |

The 80 sentinel cases were fixed before controls: missing/empty/whitespace required values; absent/nonmodule/nonregular `go.mod` paths; malformed/missing/invalid optional case/media JSON; no/empty/extra targets; unsupported chapters/unknown flag; applicable older-chapter options including empty values; and help including Chapter 10. `sentinel-matrix-plan.json`, each command receipt/output, direct positive marker receipt, and `sentinel-results.json` preserve the exact matrix. The direct target marker proves the launch sentinel was functional. **There was no fake-go wrapper and no positive Build-marker self-check.** Absence of the target marker does not prove that `grade.Build` was never called, especially when it receives an existing executable. The matrix therefore does not complete assigned control 3.

The successful Chapter 10 report is `07-ch10-positive.stdout.txt` with parsed copy. It has exact title `Edition 2, Chapter 10 — Persistence` and ordered IDs `save-shape`, `default-load`, `replay-equals-snapshot`, `tail-applied-once`, `log-not-needed`, `bad-save-refused`, `ch9-parity`, all passing, score/max score 100/100. Exit zero alone was not accepted as sufficient.

### Genuine Chapter 9 failure, not explained away

`08-ch9-positive.stdout.txt` is a complete JSON report titled `Edition 2, Chapter 9 — Skills`; stderr is empty. All eight skill checks pass (85 points). The `ch8-parity` check earns 0/15 and says:

> ch7-parity: ch6-parity: ch5-parity: ch4-parity: ch3parity: turning tool calls into jobs broke Chapter 3: ch2parity (0/10)

The adapter returns **1** consistently with the report's `passed: false`. It was not killed and did not hit the 600-second deadline. The immediately preceding Chapter 10 composition passed its Chapter 9 parity, but that does not erase this independent failing baseline. The nested report does not expose the lower Chapter 2 assertion or enough evidence to identify its cause. Calling it a confirmed flake, environment fault, adapter regression, student defect, or source change would exceed the evidence. No lower-grader probe, retry, tuning, or timing change was made.

## Reviewer errors, omitted controls, and stop behavior

1. Early read-only verification used the older manifest key `accepted_predecessor` instead of `predecessor`; an adjacent lookup used unavailable `rg` and an incorrect `internal/grade/build.go` filename. Correct names were inspected and verification completed. These attempts were not controls.
2. The first prepared runner used guessed top-level `candidate_binary` and base-HEAD fields. It stopped before preflight or any control with `KeyError: candidate_binary`. Its exact source and initial environment/workspace metadata are retained; correction used the actual nested `candidate.source`/`candidate.binary` map and before/after repository identities. No candidate file changed. The original failed attempt is not hidden by the corrected runner.
3. The external matrix omitted the explicitly required fake-Go Build-marker self-check. Its passing result is narrower than assigned.
4. The runner's planned final sequence mistakenly copied the old entry audit's historical effective-prompt mutation. The task instead requires a compile-valid **copied-adapter Chapter10-to-Chapter9 route mutation**, including an identity detector that fails even if the wrong route scores 100 and exits 0, followed by exact restoration/rebuild/pass. The erroneous branch was **never reached**: the Chapter 9 baseline failed first. No historical file, adapter copy, or student file was mutated, and no sensitivity/restoration pass is claimed. The recorded plan/driver are preserved as failed-review evidence and retired from execution, not future authorization.
5. On the genuine Chapter 9 failure, the runner stopped and reverified all source/old-seal identities. Root was asked for disposition and instructed the reviewer to seal without rerun, mutation, or repair. No additional runtime control followed that direction.

The generated completion receipt says `owned_copy_restored: false`; its precise meaning is **the copy/mutation was never started**, not that a modified copy was left behind. `final-reconciliation.json` makes that distinction explicit. Root assistance/disposition is preserved in `root-disposition.json`.

## Isolation, preservation, and evidence boundary

The active scratch root is recorded in `workspace.json`; the failed pre-control attempt's separate root is in `initial-workspace.json`. Both are reviewer-owned. The environment directs temporary grader builds/launch directories beneath owned TMPDIR, and HOME is separate from user configuration. The unchanged grader's launch helpers allocate fresh directories; persistence phases use their owned temporary save paths. No live GUI service, port 8088, private user settings, microphone, or provider account was accessed. No settings filename or binary is exported in the portable package. No copied reference was created because execution stopped before that step.

Every recorded subprocess was waited to termination. Final process inspection found no process referencing the active owned scratch. No shared cache cleanup, existing-source cleanup, commit, tag, deployment, child agent, or publication occurred. Disk-free checks stayed above the 1 GiB floor at recorded boundaries (about 2.66 GB at completion). No claim of continuous disk monitoring is made.

Eight transient executable identities were observed at 100 ms intervals. Their local hashes and observed process commands are preserved; sampling is explicitly **non-exhaustive**. No binary or archive was exported, and no extra run was made for sampling completeness. No wire/request-body trace was exported; the actual returned grader reports establish the stated pass/fail observations, not broader GUI/architecture acceptance.

Before/after receipts confirm the exact candidate, all manifest inputs/protected files, predecessor/reference files, prior inventory closures, and main/student HEAD/tree/staged/unstaged tracked-difference identities stayed unchanged. The predecessor remained clean. Original user changes were hashed for preservation, not exported as unrelated diff bodies. Already-colliding legacy artifact recovery remains a mandatory later candidate integration/source/GUI obligation; this review adds no historical score check and claims no such runtime coverage.

## Portable package and root handoff

`ch10-grader-review-evidence/inventory.json` lists this report and every regular portable evidence file with repository-relative path, byte count, and full SHA256, excluding only the inventory itself. The final submission carries the inventory's external hash. All portable files are verified UTF-8/NUL-free; no symlinks, archives, binaries, or settings files are included.

**Root action required:** inspect the failed Chapter 9 composition and this review's incomplete-control record. Do not treat passing local gates, the 100/100 Chapter 10 result, or the launch-marker matrix as complete adapter acceptance. Any further investigation/control correction requires a new explicit bounded authorization. No automatic continuation is planned.
