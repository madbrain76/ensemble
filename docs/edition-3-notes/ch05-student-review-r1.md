# Edition 3 Chapter 5 — formal code-correction round 1

Final status: clean corrected candidate frozen; all local gates and actual executable regressions passed; Chapter 5 and Chapter 1 each passed 100/100 with exit 0 on the same exact corrected CLI. This is a submission for independent whole-candidate re-review, not acceptance or a tag.

## Frozen scope and provenance

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Clean correction commit: `6c4561cbd6cbae080d952da34a5fcd588a76920f`.
- Tree: `8ba148fa32b5211b0f148e9da137544523374720`.
- Unchanged parent/original candidate: `00489b05b65cfe5209f26cb655667df44b3847e2`.
- Only three source files changed: `internal/agent/actor.go`, `internal/jobs/jobs.go`, and new public-boundary `chapter5_correction_test.go`.
- Source archive: `/tmp/cadence-ch05-correction-r1-evidence/frozen/source.tar`, SHA256 `56535661fcfdbc5121384353aa02475bec791ac04fad4d57154856216dc81bb9` (47 tracked files).
- Corrected CLI: `/tmp/cadence-ch05-correction-r1-evidence/frozen/ensemble`, SHA256 `9bd612504b9856fd4b64e0b3b386de965eae6403b0c60c00cd14cd3973bc6d61`.
- Corrected workflow: `/tmp/cadence-ch05-correction-r1-evidence/frozen/workflow`, SHA256 `acfab0d19f435529f551331f065098b6341a0d57384dce17146401733723fdf4`.
- `frozen/identity.json` records every source-file hash and binary/archive identity. Both `*-build-info.txt` receipts identify the correction revision with `vcs.modified=false`.

All newly named relative receipt paths in this report are beneath `/tmp/cadence-ch05-correction-r1-evidence/`, with portable text copies in `docs/edition-3-notes/ch05-student-r1-evidence/`. Binaries and tar archives stay outside Git. New source/evidence indexes are `ch05-student-r1-source-manifest.json` and `ch05-student-r1-evidence-manifest.json` beside this report.

No original report, evidence index, source snapshot, failed human summary, live receipt, commit, or acceptance tag was rewritten. No children, paid/provider calls, manuscript/grader/guidance edits, or unrelated feature changes were used in this correction.

## Findings, bounded fixes, and ownership/lifetime reasoning

### R1 — acknowledged hint stranded after a text-only response

The actor still appends the actual model response. For a valid nonempty text-only response, it now checks actor-owned pending ephemera before declaring the turn successful. Pending context triggers the existing `sendModel` boundary in that **same turn**. That boundary snapshots immutable request context and consumes ephemera once through `request_sent`; it does not mutate the HTTP request already sent. No dispatched-tool-batch count is incremented for this text-only continuation, and no new counter or policy requirement was introduced.

The public regression holds real HTTP, posts a hint, waits for its actual observation acknowledgement, and only then releases a normal text-only response. It verifies the first wire payload cannot contain the hint; the reliable request cannot finish before its continuation; the next payload contains the hint exactly once; completion is saved; pending ephemera is empty; and an unrelated subsequent prompt cannot inherit it. A separate ordinary no-hint positive requires exactly one request.

Actual executable regressions repeat the held-HTTP seam through the human CLI and through the author stage of the public-import workflow. Both use a configured dispatched-batch limit of one. The workflow genuinely runs its role tools and completes the author/editor/reviewer file handoffs; the extra hint continuation does not spend the one dispatched batch.

### R2 — cancelled synchronous Ask waiting behind unrelated work

`Ask` remains an adapter over `Submit`, not another conversation loop. It normalizes nil context, first checks its own private buffered result, then waits for that result, caller cancellation, or actor shutdown. If cancellation/shutdown wins, it checks its own buffer once more before returning cancellation/stopped. A result already buffered therefore takes precedence; an abandoned caller cannot block the actor's eventual completion. The actor's existing save-before-success and per-request one-element reply channel remain unchanged.

Public tests hold a distinct active HTTP request and exercise both already-cancelled and cancelled-while-queued Ask callers. The caller returns cancellation without releasing or cancelling the active request. Later concurrent callers receive their own replies; the cancelled prompt never reaches HTTP; successful later replies appear in saved history. The queued-case test instruments a real context's `Done` access while retaining its real cancellation channel; it does not replace Agent behavior. Separate real-boundary positives show an unread private result does not block a drain barrier, is saved before completion, survives subsequent caller cancellation and shutdown, and is not confused with a new stopped request.

### R3 — map-order-dependent shutdown postponing process stop

The retained-job table now performs two sweeps. First it uses the existing managed `Kill` path for all managed entries and retains cleanup errors. That path still performs its own process-group cleanup/reaping; it is not replaced with a status-only update. Only after the managed-stop sweep does shutdown wait on every job's cleanup and ordinary-handler execution. Started ordinary Go handlers are not force-aborted or given invented outcomes. Existing report-worker joining and final actor saves remain in their existing owners.

Each of eight public test layouts starts four real held Go handlers and a ready real PTY process. While every Go handler remains held, the test requires the process's cleanup channel to close **and** its OS PID to cease to exist. Root shutdown must still be waiting for the real Go handlers. Releasing them must produce their actual retained output, `done` records, and final saved `job_report` facts. No map implementation, process cleanup function, actor, or tool handler is mocked.

## Regression sensitivity: restored original versus corrected source

`run-sensitivity.py` extracts the exact original archive into an owned disposable directory, verifies all 46 original file hashes, and adds only the new public test file. Original production is never substituted into the candidate worktree. An original compile-only race run must pass before any failure can count. Environmental, build, panic, timeout, and race-detector failures are excluded from detection credit.

| Control | Original production result | Corrected result |
|---|---|---|
| Compile with new public tests | exit 0 | full final-byte gates exit 0 |
| R1 held HTTP / acknowledged hint / text-only response | exit 1; exact stranded-hint assertion once | focused race and frozen full suite pass; no-hint positive also passes |
| R2 cancelled Ask behind held request | exit 1; exact blocked-caller assertion twice | both cancellation cases, private-reply isolation, and buffered/durable positives pass |
| R3 held Go handlers plus real process | exit 1; actual-process-alive assertion in 7 of 8 layouts | all 8 layouts pass on corrected source |

The one successful original R3 layout is a legitimate favorable map ordering, not a discarded result. All eight layout results remain in `sensitivity/original-R3.txt`. The runner did not force map order. See `sensitivity/source-proof.json`, `results.json`, `commands.json`, original logs, and `corrected-focused.txt`. The added public test hash is `abd3055d5c12b9f9f3eeb2fd9b5651192195c7c1a3d24080b8dba55ef4159ef4`.

These are three restored-original regression controls, not a claim of a newly repeated broad mutation campaign. Earlier mutation and live receipts remain attributed to the original candidate.

## Final-byte local gates and command receipts

`gates/commands.json` records every command, working directory, explicit `GOWORK=off`, source commit, exit, elapsed time, and pre-step free space. All 13 commands exited 0:

- Root: `go vet ./...`; `go test -race -count=1 -timeout 5m -coverpkg=./... -coverprofile=/tmp/cadence-ch05-correction-r1-evidence/gates/coverage.out ./...`; `go build -race ./...`.
- CLI: `go build -race -o /tmp/cadence-ch05-correction-r1-evidence/frozen/ensemble ./cmd/ch01`.
- Independently from `gui`, `examples/ch04-consumer`, and `examples/ch06`: `go vet ./...`; `go test -race -count=1 -timeout 3m ./...`; `go build -race -o <recorded outside-Git binary> .`.

The full root inherited/race run took **140.4 seconds**. Combined statement coverage is **86.2%**, from `gates/coverage.out` and `coverage-summary.txt`. GUI and Chapter 4 consumer remain compile/test module gates, not a claim of a live GUI session. No long test was replaced by a short observation window.

### Actual executable boundary evidence

`local-interface.py` has no live-provider route or credential loader: subprocess environments contain a synthetic fixture key and a loopback base URL. The corrected source and both exact executable hashes are checked before launch.

`local-interface-r2/summary.json` reports exit 0 for both real programs:

- CLI: actual hint acknowledgement while first HTTP is held; 3 local requests; same-turn continuation; once-only hint consumption; no leak to the unrelated request.
- Workflow: actual author hint acknowledgement while first HTTP is held; 7 local requests; ordered reliable role handoffs; successful `write_draft`, `word_count`, `read_draft`, `edit_draft`, and `review`; genuinely edited manuscript and nonempty review artifact. A one-batch-per-role limit remains sufficient.

Raw local stdout/stderr, input timing, request bodies, actual neutral history, and workflow files are retained. These are **fake-provider** receipts through actual executable and tool boundaries, not live inference or mocked Agent behavior.

## Preserved mistakes and assistance

- The first focused test run failed because the new R3 fixture tried to JSON-decode the existing human-readable job report. It was changed to obtain the handle from the public recorded `ToolEvent.Job`. `focused-first.txt` is preserved and earns no sensitivity credit. R1/R2 passed in that run. `focused-corrected.txt`, the final sensitivity corrected pass, and the frozen full suite establish the corrected fixture's result.
- The first local executable capture required exactly `Hint: ` plus the marker. The CLI correctly preserved an additional leading space in hint text. `local-interface/cli/stderr.txt` contains the actual acknowledgement; the capture timed out and terminated its child with exit -15. The script version and failed receipts are preserved. Only the capture matcher was relaxed to recognize the actual acknowledgement. The source and binaries did not change. `local-interface-r2/` is the successful rerun.
- Parent supplied only the three public findings/contracts, approved the bounded plan before edits, and required preserving the favorable original shutdown layout and all failures. No private audit/probe source, historical solution, or grader implementation was read.
- The earlier approval-plan fields, live controls, and original human false-positive reconciliation are not relabeled as corrected-source proof. This round did not invoke them or any provider.

## Limits and pending independent review

No paid calls were authorized or made in this round. The 12 original provider requests and their artifacts remain evidence of **00489b05**, not fresh live evidence of **6c4561cb**. Current evidence is full frozen-source local testing, real OS/process cleanup, actual CLI/workflow local HTTP boundaries, and the public routes recorded below. The noncancellable Go-handler shutdown wait is deliberate and honest; this patch does not promise to force-stop an arbitrary Go handler that never returns. Ordinary observation loss, whole-part notifications, provider/media policy, existing dispatched-batch default of 200, and the GUI stub are unchanged. The policy-refusal guard was not redesigned or claimed to be the dispatched-batch counter.

Parent retains whole-candidate review, acceptance, and tagging authority. No acceptance tag was created.

## Completed same-binary public routes

Both commands ran sequentially from `/Users/bill/projects/ensemble`, after the clean correction commit and final binaries were frozen:

```sh
go run ./cmd/grade-edition3 -ch 5 -json -core /Users/bill/projects/ensemble-edition-3-student -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer -workflow /Users/bill/projects/ensemble-edition-3-student/examples/ch06 -media-models /tmp/cadence-ch05-correction-r1-evidence/media-models.json /tmp/cadence-ch05-correction-r1-evidence/frozen/ensemble
go run ./cmd/grade-edition3 -ch 1 -json /tmp/cadence-ch05-correction-r1-evidence/frozen/ensemble
```

| Route | Process exit | Public result | Measured wall time |
|---|---:|---:|---:|
| Chapter 5, explicit core/consumer/workflow/media map | 0 | passed, 100/100 | 159.073 s |
| Chapter 1, same exact CLI bytes | 0 | passed, 100/100 | 127.260 s |

The Chapter 1 JSON independently reports complete answer/usage after **1m2s headers delay** and after **1m2s body delay**. These phases actually ran; observation windows were never used as deadlines. `frozen/grade-commands-and-exits.json` binds both complete commands and exits to the same CLI SHA256 `9bd612504b9856fd4b64e0b3b386de965eae6403b0c60c00cd14cd3973bc6d61`, corrected source commit, workflow hash, media-map hash and recorded main HEAD. The raw public JSON and stderr remain alongside it. No grader implementation or private probe source was inspected.

`coordinator-exits.json` preserves the top-level exits, including both fixture/capture failures, and points to the detailed per-command ledgers. Initial fixture attempts were intermediate, unfrozen test/capture setup and are not promoted into final-byte or sensitivity successes. The final controls, full gates, local executable receipts and public grades have exact frozen attribution.

## Portable handoff and preservation

Only this new report, the new `ch05-student-r1-evidence/` text directory, and new exact source/evidence manifests are written under the permitted main note paths. The original 341-file portable set, original report/index, initial-candidate local artifacts and source snapshot are verified against their existing hashes rather than rewritten. Any generated-cache exclusions are called out by the index, not relabeled as source/evidence changes. Binaries and archives stay at their recorded outside-Git paths. The new text verifier checks the copied files and report without executing code under test or accessing credentials. The manifest records privacy scan results, preserved-source verification, all exclusions and exact hashes.

No additional provider calls, source edits, optional features, acceptance tags or shared-cache cleanup followed these checks. Independent re-review remains with the parent.
