# Chapter 8 entry controls — finite audit only

No student launch, inherited parity, browser, paid/provider-key run, source/grader edit, child, or background monitoring loop occurred. These are local entry-audit controls, not Chapter 8 acceptance.

## Execution environment and ownership

Owned scratch: `/tmp/ch08-entry-audit.aTPJZB`. Commands ran serially. The environment wrapper uses `env -i`, an owned HOME/TMPDIR, `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOFLAGS=-p=1 -mod=readonly`, and `GOMAXPROCS=2`. Existing Go module/build caches were reused, never deleted. No real credential was inherited by a control. GUI launches use the existing local fake provider and its literal fake key. No provider conversation is requested by the settings-only driver. Settings files live only in the launcher's temporary workspace, reused across its one restart then removed.

Disk checks remained above the 300 MiB stop boundary; the final pre-cleanup measurement was 997,088 KiB available. The owned scratch occupied 20,608 KiB before binary removal. There was no resource failure or retry.

## Exactly one package-test command

From repository root, through `clean-env`:

```sh
go test -p 1 -count=1 -timeout=60s \
  -run '^(TestInvalidInputsFailBeforeGrading|TestHelpSucceedsWithoutLaunchingTarget|TestChapter7ValidationBeforeLaunch|TestChapter8Settings)$' \
  ./cmd/grade-edition3 ./internal/gradeedition2
```

Exit **0**. Output: publication package **0.010s**, inherited package **0.828s**. This selected the publication negative-validation/help tests and the existing `TestChapter8Settings` only. That historical settings test built `solutions/ch08/gui`, used fake HTTP and an isolated workspace, connected two subscribers, applied sparse updates including false/zero, and restarted. It did **not** call `Chapter8`, `Chapter7`, or inherited paid parity. It is historical positive wire evidence, not an Edition 3 student result. Package build overhead exceeded these per-package timings; the owned job was explicitly joined before further controls.

## Standalone owned builds and runs

Two serial `go build` commands (not additional `go test` commands): publication command to the scratch binary and the scratch driver's own module. Both exited **0**; build logs are empty. The scratch driver imports the existing grader package and uses an audit-only `go:linkname` bridge to its unexported `settingsChecks`. This is deliberately not a proposed publication API. No source/overlay/test file was written into a grader package. Exact driver and mutant sources are retained in `fixtures/`; `prepare.py` reconstructs them using the pinned repository's `go.sum`.

`run-controls.py` records argv, PID/process group, actual exit, timeout and elapsed time in `control-results.json`. Every process group had disappeared when its parent exited; the watchdog did not need to kill anything. Each run had a 30-second external limit. The underlying GUI driver also bounds launch, reads, shutdown, and process lifetime.

| Control | Actual outcome | Interpretation |
| --- | --- | --- |
| publication `-h SENTINEL` | exit 0; no sentinel marker | Help lists only 1/4/5/6/7; no Chapter 8 claim. |
| publication `-ch 8 SENTINEL` | exit 2; no sentinel marker | Unsupported chapter is rejected before target build/launch. This is the entry blocker, not a failed student run. |
| `initial-budget-omitted` | exit 0; all four settings checks pass | Deliberately invalid initial `current_settings` omits `max_tool_rounds`; later patches/persistence are valid. A demonstrated false negative in required full initial snapshot coverage. |
| `no-persistence` | exit 1; first three pass, `settings-persist` fails | Updates broadcast correctly but are never written. Actual restart receives defaults; the saved log reports `theme: want present system, got dark (present=true)`. A demonstrated true negative for restart persistence. |

The mutants implement only the narrow settings wire required for these controls. They are **not** conforming Chapter 7/8 applications and their scores do not imply parity, ownership, budget execution, or browser correctness. Exactly these two mutant cases ran, once each. No active-turn budget, malformed-input or write-failure mutant was run.

## Evidence interpretation

`identity-before.json` verifies main HEAD, student HEAD/tree, the peeled annotated Chapter 7 tag, 19 pinned inputs, 120 student files, and 25 protected tracked worktree files. `identity-after.json` supplies the final counterpart and the 465-file protected-source scratch baseline comparison. The manifest's `student_tag` stores the resolved commit. The receipt compares `edition-3-ch07-r1^{}` to that same pinned student commit. The initial exploratory unpeeled-tag comparison was corrected before controls; the resolved identity matches.

No solution-level `TestSettingsWireAndDisk` or `TestSavedBudgetEntry` was executed. Those tests were inspected as historical source evidence only. Missing-path exploratory reads (`ensemble/`, `gui/server.go`, `internal/gradeedition3/ch07.go`) were corrected to the actual root core module, `gui/settings.go`, and `gradeedition2.Chapter7` dispatch; they did not run or alter anything.

See `final-state.json` for final process/disk/temp cleanup receipts and the allowed-worktree status comparison. `inventory.json` enumerates these durable controls with byte sizes and SHA-256 hashes (excluding itself).
