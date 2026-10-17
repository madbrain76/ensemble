# Chapter 5 public wakeup amendment — finite student candidate

## Status and scope

This is **new work under the authorized public amendment**, not a retroactive
failure of the preceding candidate and not a second formal code-correction round.
The existing correction count remains one; R1–R3 are retained and tested. Parent
review/acceptance/tagging is still required. No acceptance tag, main-repository
commit, push, install, GUI session, credential lookup, or paid provider call was
performed for this amendment.

Final source: **`f60ebf3aad2ece4251face9d95c106dd667abcd6`**, extending the clean,
reviewed predecessor `6c4561cbd6cbae080d952da34a5fcd588a76920f`. There are 48 tracked
source files. Seven production files and one added test file differ from that
predecessor. All final-byte local gates, actual executable regressions, and both same-binary
public grades passed; exact provenance and limits are recorded below.

Permitted public inputs were fully read and pinned before edits: main commit
`4f1701a76f57ec57d362d4f1a30264f30b65ddf8`, AGENTS, exact coding skill path,
Edition 3 coder role, public student brief and shared execution guide, and both
current chapters in full. The student brief replaces answer-bearing historical
carryover for this role; no historical solution, private audit/probe, grader
implementation, final reference, parent memory or other agent file was consulted.

- Chapter 4 SHA256: `030c9a7ceb6c18f56596e840bb9df404fa418ea5f22c77bc7712b7ec1e6d8c55`.
- Chapter 5 SHA256: `4f4c2595ef65a08511d5cc68553f853feaf854e4cb0fe7d4d9094d4ae7f5a344`.
- Full mandatory input hashes: `input-pins.json` in the new evidence directory.

## Ownership, lifetimes, and ordering

The parent approved the bounded ownership plan (main-18) before source changes.
No old architecture approval was requested and no new scheduler, subagent system,
GUI, streaming, OS abstraction, messaging framework, or special join tool exists.
An ordinary registered held Go handler supplies the join-like boundary in tests.

1. **Actor ownership.** One mailbox still owns the active request, original model
   call batch/index, pending one-shot limits, history, admitted message and batch
   disposition. Sender, recipient, request ID and turn ID are copied scalar data,
   not dispatch/routing policy. The actual message parts and supplied metadata
   reduce into model context; the hint acknowledgement is not delivery proof.
2. **Observation ownership.** Each prepared call has its own cancellation-cause
   context and observation identity. A busy `Hint` or `UserMessage` is durably
   appended before canceling only that observation. It neither retires the turn
   nor cancels the already-sent HTTP request, job execution context, retained
   output, or other continuing jobs. Caller/turn cancellation keeps its distinct
   ordinary cancellation semantics.
3. **Job/report ownership.** `Job.Wait` retains its one locked status/cursor/report
   path. A message wake uses that path to report the actual running/terminal state,
   handle, bytes and real error, rather than an invented failure. Its timer is
   stopped on every return. The retained job and its execution worker outlive the
   observation. A failed retained control result preserves the nonempty report
   instead of replacing it with just an error string.
4. **Batch ordering.** The actor accepts the current report once (turn plus
   observation identity), keeps completed real outcomes, then pairs the remaining
   undispatched calls with exact-ID not-executed errors before sending a new model
   request. No skipped handler, filesystem effect or new job is created; a skipped
   setter does not arm limits. Literal next-call consumption is shared with the
   existing dispatcher, including synthetic undispatched errors. Original assistant
   calls and matching signed content remain intact. There is no automatic replay.
5. **Dispatch races.** If the actor has actually dispatched the next call before
   admitting the hint, that call is a started call and returns its real observed
   state. A setter may already have completed synchronously before its report is
   handled; it is reported truthfully, and its literal next skipped call consumes
   the pending setting. The model may issue a genuinely new call after waking
   while the older job is still alive.
6. **Available output compatibility.** Parent main-19 clarified the one ambiguity
   found by the inherited PTY gate: `wait_for_job` enables a small existing
   `WaitOptions.ReturnAvailable` data option for its initial locked snapshot.
   Already-unseen bytes return immediately; subsequently arriving nonmatching
   bytes notify/recheck but do not bypass the pattern/delay/completion rules.
   Initial calls, `send_input`, and `kill_job` keep their inherited callback rules.
   All these observations remain message-wakeable. No new control tool or method
   was introduced.
7. **Shutdown.** Existing stop-admission, all-managed-process stop sweep, honest
   ordinary-handler/report-worker join, and save-before-success boundaries remain
   unchanged. The original R3 shutdown regression still exercises them. Actual
   workflow EOF proof below also requires the real outstanding Go-handler output.

### Changed files and requirements

| Source | Narrow responsibility |
|---|---|
| `internal/agent/actor.go` | Actual message admission; per-observation cause/identity; batch cutoff, truthful report acceptance and same-turn continuation; literal next-error disposition |
| `internal/agent/dispatch.go` | Shared one-shot consumption; preserve nonempty failed control reports; enable initial available-output behavior only for `wait_for_job` |
| `internal/common/actor.go` | Optional application-supplied `UserMessage` correlation data |
| `internal/common/common.go` | Durable `MessageEvent` correlation data |
| `internal/common/jobs.go` | Observation-wake cause vocabulary and small observation-data option; no shared behavior implementation |
| `internal/history/history.go` | Preserve actual message and supplied correlation in reduced request context |
| `internal/jobs/jobs.go` | Locked truthful wake report, cursor/cap preservation, initial-available check and timer cleanup |
| `chapter5_wakeup_test.go` | Public API/real HTTP/real process/ordinary registered-handler regressions and race controls |

## Regression and sensitivity evidence

No production function mocks are used. Public `Agent`/`Job` methods, actual
loopback provider HTTP, real managed shell commands, and ordinary registered Go
handlers are used. New fields/options are set with JSON in predecessor controls
so the original candidate **compiles**: losing behavior, not losing a struct
field, must cause the failure.

| New public regression | Direct evidence |
|---|---|
| `WakeupInitialMixedBatch` | Held ordinary handler at 3600 seconds; addressed question plus metadata in the actual next HTTP; completed result preserved; current handle/running state; skipped effect/setter paired without execution; matching signature retained; new work while old job lives; same-ID terminal output/cap/full retained file/repeat cursor; saved history |
| `WakeupSetterDispatchBoundary` | Setter that actually dispatched reported truthfully; literal next skipped call consumes it; skipped later setter never arms; no effect/job/limit leakage |
| `WakeupLaterProcessWait` | Real process; explicit hint wakes later 3600-second `wait_for_job`; actual next HTTP contains message plus same-ID running result; same process later produces real final output |
| `WakeupTimerAndAvailableOutput` | Initially available unseen bytes return immediately; initially empty later wait rechecks nonmatching then matching arriving bytes; cursor excludes old output; terminal immediate result |
| `WakeupShortTimerCannotReleaseLaterWait` | Explicit initial hint then addressed question; later 3600-second observation survives the earlier short timer's former deadline; no stale/duplicate tool result or canceled job |
| `WakeupCompletionRaces` | Message-first, completion-first, simultaneous schedules; actual question survives, each call has one result and real completion is retained |
| `WakeupFailedRetainedReport` | Real failed ordinary job's retained output/handle/error survives later control observation |
| `WakeupAddressedProcessInitialAndLater` | Real process questions at both initial and later 3600-second observations; exact supplied metadata in each actual next HTTP; once-only message consumption; same job remains active and eventually returns real terminal output |

Final focused command (exact receipt: `compatibility-focused-r4.json`):

```sh
go test -race -count=10 -timeout 120s \
  -run 'TestChapter5Wakeup|TestChapter5Correction|TestChapter4PTYInputCWDPatternAndShutdown' -v .
```

Exit **0**, elapsed **8.628s**, including the closed R1 text-only continuation and
no-hint positive, both R2 cancellation cases and durable result preference, and
all eight R3 process/Go-handler shutdown layouts. The short test deadlines are
failure bounds, not substitutions for configured 3600-second observations.

Final sensitivity campaign: `run-sensitivity-compatible.py` and
`sensitivity-compatible/{commands.json,source-controls.json,*.log}`. It creates a
scratch tree from exact predecessor `6c4561c`, verifies its original files, adds
the same final public test bytes, and runs seven independent original-production
controls. **All seven compile and fail behaviorally** (initial mixed batch,
setter/batch disposition, later process wait, available-output wait, short timer /
later question, failed retained report, addressed initial/later process). It then
restores the complete corrected production diff from `6c4561c` and the same tests:
**exit 0**. The candidate worktree is never reverted. Completion-race controls and
all preserved correction tests are included in that restored passing run; no
claim is made that every arbitrary race schedule must fail on the predecessor.

## Attempts, assistance and limits of credit

All failures remain in the new evidence tree; see `attempt-ledger.md` and saved
logs. In particular:

- Wrong retained-output path in the first mixed fixture, model-ID mismatch in the
  signed fixture, wrong process-test types and metadata expectations, and failed
  scripted edit anchors are disclosed as fixture/setup failures, **not valid
  sensitivity detections**.
- The first frozen amendment candidate `ccf82d8a5101a4a79a4f3d0979db04450f163282`
  passed focused checks but failed the full inherited PTY input test after
  140.054s. Its overly broad unseen-output readiness changed `send_input` behavior.
  This is a genuine intermediate implementation compatibility failure, not a
  transport/environment failure. Its source archive, full log, and exit remain.
- The parent answered the concrete `wait_for_job` versus `send_input` ambiguity
  before the small observation-data option was implemented. This assistance and
  main-18 plan approval are disclosed; the parent did not author the implementation.
- An owned Git pager paused a freeze command. Only that identified pager was
  stopped; the original command continued, without a duplicate commit/build.
- Builds are sequential, disk checks are recorded, and no shared cache/tmp cleanup
  occurred. Free space stayed above the launcher stop threshold. No compiler,
  disk, timeout harness, or route-mock failure is credited as feature sensitivity.

The historical public grader does **not** cover this new wakeup contract. Its
passing scores, once recorded below, are regression gates, not a substitute for
these public regressions and independent review. Local HTTP fixtures are not paid
provider/model validation. The original **12 paid calls belong only to initial
candidate `00489b05`**, not `6c4561c` or this amendment; they were neither rerun nor
relabeled. No future messaging/join/MCP framework is implemented.

## Final-byte gates, executable proof and public grades

All **13 final-byte gates exited 0**, with `GOWORK=off` as recorded in each
`environment_overrides`, and no concurrent Go builds. The root
race/coverage run completed in **141.032s**. The operational CLI is `cmd/ch01` in
the root module; the actual additional modules are GUI, consumer and workflow
(no separate CLI marker module exists in this accepted source).

| Gate | Exact command | Exit | Elapsed seconds |
|---|---|---:|---:|
| root-vet | `go vet ./...` | 0 | 0.239 |
| root-race | `go test -race -count=1 -timeout 5m -coverpkg=./... -coverprofile=/tmp/cadence-ch05-wakeup-evidence/candidate-r2/gates/coverage.out ./...` | 0 | 141.032 |
| root-build | `go build -race ./...` | 0 | 0.526 |
| cli-build | `go build -race -o /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble ./cmd/ch01` | 0 | 0.289 |
| gui-vet | `go vet ./...` | 0 | 0.123 |
| gui-race | `go test -race -count=1 -timeout 3m ./...` | 0 | 0.104 |
| gui-build | `go build -race -o /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/gui .` | 0 | 0.256 |
| consumer-vet | `go vet ./...` | 0 | 0.096 |
| consumer-race | `go test -race -count=1 -timeout 3m ./...` | 0 | 0.111 |
| consumer-build | `go build -race -o /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/consumer .` | 0 | 0.272 |
| workflow-vet | `go vet ./...` | 0 | 0.110 |
| workflow-race | `go test -race -count=1 -timeout 3m ./...` | 0 | 1.359 |
| workflow-build | `go build -race -o /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/workflow .` | 0 | 0.295 |

Coverage (informational; no added numerical acceptance floor): `total:						(statements)		87.1%`.

### Actual local executable wakeup proof

`python3 /tmp/cadence-ch05-wakeup-evidence/local-wakeup.py /tmp/cadence-ch05-wakeup-evidence/local-wakeup-r5`
exited **0**, and both actual candidate processes exited **0**. Full synthetic
HTTP requests, safe actual stdout/stderr, process receipts, retained outputs,
written manuscript/review, and saved histories are retained under
`local-wakeup-r5/`; `summary.json` carries the binary hashes.

- CLI: initial 3600-second observation woke in **0.002017s**; later same-job
  `wait_for_job` at 3600 seconds woke in **0.001879s**. Five real HTTP exchanges
  against the local fake show actual hint plus correct paired results. Job 1
  stays live; new job 2 writes the answer while it runs; skipped write has no
  file/job. A fresh wait on job 1 yields `ACTUAL_PROCESS_FINAL`, and EOF saves one
  correlated result per call. Exactly two retained jobs exist: controls and
  undispatched siblings create none.
- Workflow: the ordinary registered `think` observation at 3600 seconds woke in
  **0.003806s**. Seven local HTTP exchanges show actual hint, running job 1 and
  exact skipped draft error. Real author/editor/reviewer writes finish while the
  old handler is still alive. EOF then honestly joins the four-second handler,
  retains `ACTUAL_ORDINARY_THOUGHT`, and saves its terminal `job_report`. Actual
  draft/revision/review effects are required, not merely pipeline markers.
- Earlier executable attempts are preserved: r1/r2 used the wrong state matcher,
  r3 used an invented workflow input kind (the CLI scenario passed), and r4
  expected `job_reported` instead of the actual saved `job_report` (both processes
  exited 0). None of those launcher exits is relabeled as a passing run; r5 is
  the clean final script execution. A mistaken commentary suspicion about the
  `think` parameter was corrected by inspection—it already used `thought`.

### Public grades

Both exact-binary public routes **passed:true, 100/100, exit 0**. The chain itself
exited 0. Chapter 5 elapsed **159.223s**; Chapter 1 elapsed **127.251s**.
Chapter 1 explicitly completed both real **1m2s delayed-header** and **1m2s
delayed-body** phases with complete answer and usage. No short subprocess deadline
or simulated wait substituted for those phases. Commands ran from main
`/Users/bill/projects/ensemble` at pinned `4f1701a7…`:

```sh
go run ./cmd/grade-edition3 -ch 5 -json -core /Users/bill/projects/ensemble-edition-3-student -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer -workflow /Users/bill/projects/ensemble-edition-3-student/examples/ch06 -media-models /tmp/cadence-ch05-wakeup-evidence/candidate-r2/media-models.json /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble
go run ./cmd/grade-edition3 -ch 1 -json /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble
```

`candidate-r2/frozen/grade-commands-and-exits.json` binds both command exits,
elapsed times, grading-main HEAD, final source commit, and the **same** CLI SHA256
`e51138885a359eaf801258a52bc21f28ca41691683256647f8c9911be410b233`.
The map is the copied public fake-grade model map (SHA256
`cc75c8f40899ce7dd6046f26c413c7afff518425bb86c72b853306966de9796b`).
Raw public JSON reports and stderr are preserved. These historical checks are
regression gates, **not** evidence that the newly amended wakeup contract was
privately or automatically graded. That still requires the supplied public
regressions and independent whole-source review.

### Exact final-byte provenance

- Source commit: `f60ebf3aad2ece4251face9d95c106dd667abcd6`; tree: `3de5d0d24dd9f593ec017fceccfbff1b209c5d57`.
- Exact `git archive` SHA256: `cc58a900e2f5e2040df75455c9b61d8d24f9b29afa65b627bf1d91a4b63d6e46`.
- Race cli: `/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble`; SHA256 `e51138885a359eaf801258a52bc21f28ca41691683256647f8c9911be410b233`.
- Race workflow: `/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/workflow`; SHA256 `c724d5018c45152efb47079226d2fce825cc23bdc3b327714dc5bb0d0338caf4`.
- `final-source-verification.json`: clean worktree/commit/tree, all **48** Git
  blobs, worktree files, exact archive members and new text snapshot agree.
  Both binaries report `-race=true`, exact `vcs.revision`, `vcs.modified=false`.
- All these are new-candidate proofs, not reuse of original binaries or source
  identities. Binary files and all source archives remain outside Git.

## Preservation and portable handoff

`prior-preservation.json` verifies unchanged: all **341 original portable and
341 original local text files**, **14 original excluded artifacts**, all **83
correction portable and 83 correction local text files**, **5 correction excluded
artifacts**, the pinned original/correction reports and indexes, and the original
**46-file source snapshot**. None was rewritten. The new evidence has its own
names and the original main source snapshot remains at the initial candidate.

Only these new main paths are owned for this amendment:

- `docs/edition-3-notes/ch05-wakeup-student.md` (this report)
- `docs/edition-3-notes/ch05-wakeup-student-evidence/` (portable UTF-8 text only)
- `docs/edition-3-notes/ch05-wakeup-student-source-manifest.json`
- `docs/edition-3-notes/ch05-wakeup-student-evidence-manifest.json`

Full local artifacts are under `/tmp/cadence-ch05-wakeup-evidence/`. Exact source
archive and race binaries stay outside Git, with paths/bytes/hashes indexed. The
new source text copy, if included in portable evidence, is explicitly student
source, not a reference solution, and is not automatically staged. The final
portable verifier checks relative evidence/report/source-manifest hashes after
relocation; source archive membership and final binary/VCS stamps are verified
before packaging. No acceptance/tag is implied by this handoff.

## Reader friction and bounded follow-up notes

The amended distinction between waking an observation and interrupting a turn is
clear and testable. The only contract ambiguity encountered during implementation
was the scope of initially available output on the shared wait primitive versus
`send_input`; the parent resolved it narrowly, and both the inherited PTY test and
new later-output pattern control are retained. A possible future wording tweak
would state explicitly that subsequent output wakes/rechecks the observation but
need not return unless its pattern, completion, incoming-message or silence
boundary is reached. No manuscript/guidance change was made here. The local
launcher naming/schema errors were my fixture mistakes, not defects attributed
to the chapter or excuses for the original failed receipts.
