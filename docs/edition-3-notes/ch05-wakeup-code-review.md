# Edition 3 Chapter 5 — whole-candidate job-wakeup code review

## Verdict and phase boundary

**ACCEPTED — no blocking findings in the frozen amended candidate.**

This is the independent whole-candidate review of Bill's authorized public
job-wakeup amendment, not another entry/adapter audit, a prose re-review, or a
second formal code-correction round. The accepted pre-amendment R1–R3 closure
remains intact. Formal corrections used remain **one of three**. Root retains
acceptance, staging, commit and tag authority; this reviewer created no tag and
made no production, manuscript, guidance, grader, application or reference edit.

The verdict rests on source/ownership review, exact-binary public regression
gates, independent actual-boundary positives/original-production negatives/
restored positives, the complete pristine race/vet/module checks, and inspection
of the retained real-executable local evidence. A score alone does not cover the
amendment. No optional polish or new framework is made a condition of acceptance.

## Exact identities and input verification

| Input | Verified identity |
|---|---|
| Main committed amendment | `4f1701a76f57ec57d362d4f1a30264f30b65ddf8` |
| Student, clean | `f60ebf3aad2ece4251face9d95c106dd667abcd6` |
| Student tree | `3de5d0d24dd9f593ec017fceccfbff1b209c5d57` |
| Accepted predecessor / original-production controls | `6c4561cbd6cbae080d952da34a5fcd588a76920f` |
| Chapter 4 | `030c9a7ceb6c18f56596e840bb9df404fa418ea5f22c77bc7712b7ec1e6d8c55` |
| Chapter 5 | `4f4c2595ef65a08511d5cc68553f853feaf854e4cb0fe7d4d9094d4ae7f5a344` |
| Exact source archive | `cc58a900e2f5e2040df75455c9b61d8d24f9b29afa65b627bf1d91a4b63d6e46` |
| Frozen race CLI | `e51138885a359eaf801258a52bc21f28ca41691683256647f8c9911be410b233` |
| Frozen race workflow | `c724d5018c45152efb47079226d2fce825cc23bdc3b327714dc5bb0d0338caf4` |
| Parent candidate receipt | `6f5ae577f0523f698532a423c9b34b59ff7dc1935854f67d34cdd95853961c29` |
| Student report | `6e4188f8656a7ffab9409956bbd0dbec88623d1efa5aad9dfb65e58cb418f680` |
| Student source manifest | `d6388c761707d26ad6629a811381632e11a974fc13d330a741fddabe0e5b3748` |
| Student evidence manifest | `48226502a9a2aa0324e87815f39b77b8633abbf838de18c7accaeb559456a89f` |

All 48 source files were independently verified against Git/worktree and the
exact extracted archive; all 315 indexed portable text files and 11 excluded
local artifacts verified. The complete 290-line report, attempt ledger,
assistance chronology, test source, sensitivity receipts and relevant captured
local requests/outputs/effects were read, not just the summary. Actual `go version
-m` reads of both frozen executables show race builds, the exact candidate VCS
revision, and `vcs.modified=false`.

All current mandatory reviewer documents were fully reloaded: AGENTS, coding
skill, reviewer role, governing workflow and Edition 2 carryover, voice/writing
procedure, current Edition 3 execution/student/readiness/carryover, and complete
amended Chapters 4 and 5. Truncated combined reads were recovered. The chapter
hashes above, not a private remembered requirement, define this review. The
seven-file production diff plus new `chapter5_wakeup_test.go` was inspected against
the previously reviewed source/architecture. No historical answer implementation
was substituted for student source. Private control code remains outside Git and
is not a future student input.

## Source/contract assessment

### Actor, message and observation ownership

Chapter 4 §§4.3–4.6 and Chapter 5's **Job reports are observations** and **Turn
completion and shutdown** sections distinguish waking an observation from
canceling a turn, HTTP request or job. `internal/agent/actor.go` keeps admission,
original batch/index, current observation identity, one-shot state and durable
history on the actor. The message fact is appended before `wakeObservation`
(lines 364–369 and the input branches around 558–585). Per-observation cancellation
uses the distinct common wake cause; execution and retained output remain owned
by the existing job. Observation identity checks reject stale reports.

The added fields in `internal/common/actor.go` and `common.go` are optional scalar
application correlation, not a routing system. `internal/history/history.go`
projects actual content plus supplied sender/recipient/request/turn metadata
(lines 364 onward). The independent control sees those exact values once in the
next actual local HTTP request; it does not infer delivery from an observer ack.
Held-HTTP/frozen-request behavior and same-turn text-only continuation still pass
the previously accepted independent R1 control.

### Common jobs, mixed batches and literal-next settings

Chapter 5's job/batch requirements (particularly lines 603–665) require truthful
state for dispatched work, preservation of completed results, exact paired errors
for every undispatched sibling, no skipped effect/job/replay, unmodified assistant
calls/signatures, and fresh work while the old job remains alive.

The existing dispatcher creates ordinary jobs through the common table; controls
retain their existing no-new-job path. The actor's wake-tail disposition runs
only after accepting the actual dispatched report. `notExecuted` and
`consumeLimits` in `internal/agent/dispatch.go` share actor-owned one-shot
accounting; a skipped setter cannot install settings. The actual-dispatch branch
is not relabeled as an undispatched call. The full suite and retained sensitivity
receipts cover the setter-dispatch and completion races, and the independent
mixed case checks completed output, matching signed content, all expected call-ID
pairs, skipped handlers/jobs, a new call while the held handler is live, and a
second addressed-question wake of the same handle.

A source question about over-budget rejection and pending settings was tested
rather than reported as a presumed defect. `corrected-literal-limit-r3` passed:
one rejected batch causes no effect and no settings leakage into the next actual
call. No finding is raised. Dispatch batches remain the budget unit; the old
entry-audit refusal-count shorthand remains withdrawn.

### Output readiness, cursor, timer and process compatibility

`internal/jobs/jobs.go:287–340` uses the same locked status/output/cursor/report
path for completion, pattern, elapsed observation and message wake. `timer.Stop`
is deferred. `ReturnAvailable` affects the initial snapshot only; only
`wait_for_job` enables it (`internal/agent/dispatch.go:227`). New bytes notify and
recheck without unconditionally completing an existing pattern observation.
This is a small data option on the existing mechanism, not another scheduler.

The inherited PTY regression independently passes. A separate real PTY control
proves new nonmatching output does **not** prematurely finish `send_input`, an
actual hint does release its 3600-second observation with truthful running output,
and the same process later exits **7** with its real terminal bytes. A separate
ordinary-job control proves already-available output returns while work is still
running; later matching output releases a blocked pattern observation without
replaying old cursor bytes; and legal `kill_job` plus one-shot settings kills and
joins a real managed process without a control job. Instantaneous setters/kill
completion are not misrepresented as long observations. The retained short-timer,
failed-result, simultaneous-completion and dispatch-boundary tests are included
in the independently run full race suite.

### Inherited whole-candidate obligations

The Engine remains a one-HTTP-request executor; the Agent still owns the model
loop/history and the public parent interfaces remain the ownership path. No
fleet, MCP, hidden join primitive, callback messaging framework, streaming or GUI
feature was added. The eight-file change set does not alter the public CLI,
consumer, workflow, renderer/media tables or provider adapters.

The accepted independent R1–R3 probes were hash-checked and rerun against this
candidate. They pass: no-hint positive; a held-request hint consumed by the same
turn without mutating the frozen request or repeating next turn; already-canceled
and cancel-while-queued callers; independent buffered durable replies across
shutdown; and four real process-group cleanup/reaping layouts before releasing
four held Go handlers, with true results/reports saved. Budget controls execute
200 batches plus the final consuming request, configured boundaries, paired
no-effect over-budget errors and later usability.

The pristine root race suite also exercises the selected public explanation
policy, actual job/PTY/cursor/report behavior, public actor/history paths, media
and provider rendering/replay, CLI protocol and workflow behavior. The separate
consumer/workflow modules pass; GUI module checks were noninteractive compile/
package checks, not a live GUI session. Unchanged parts are reviewed in their
verified prior context plus these regressions, not claimed to be exhaustively
proven by the new two-scenario probe alone.

## Independent measurements

Every command's exact argv, cwd, child exit, elapsed time, watchdog and full
stdout/stderr is indexed in `ch05-wakeup-code-review-evidence/`. There were no
watchdog expirations. Builds/tests were sequential. Free disk stayed roughly
2.5–2.6 GiB; no shared cache or unrelated artifact was deleted.

Both required public gates ran from main against the **same pinned shipped CLI**:

```sh
go run ./cmd/grade-edition3 -ch 5 -json -core /Users/bill/projects/ensemble-edition-3-student -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer -workflow /Users/bill/projects/ensemble-edition-3-student/examples/ch06 -media-models /tmp/cadence-ch05-wakeup-evidence/candidate-r2/media-models.json /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble
go run ./cmd/grade-edition3 -ch 1 -json /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble
```

Both return `passed:true`, **100/100**, exit **0**. The Chapter 1 report contains
both actual **1m2s** delayed-header and delayed-body phases with complete answer
and usage. Observation callbacks were not treated as deadlines. These inherited
graders do not test the entire wakeup amendment. Chapter 5 builds its public
consumer/workflow from their explicit source directories; it does not directly
execute the frozen workflow path. Suites/probes use separately extracted source;
no byte-identical rebuild claim is made.

| Recorded command/control | Exit | Seconds |
|---|---:|---:|
| `grade-ch05` | 0 | 160.028 |
| `grade-ch01` | 0 | 127.273 |
| `root-race` | 0 | 138.815 |
| `root-vet` | 0 | 0.336 |
| `consumer-race` | 0 | 0.175 |
| `consumer-vet` | 0 | 0.124 |
| `workflow-race` | 0 | 1.765 |
| `workflow-vet` | 0 | 0.124 |
| `gui-race` | 0 | 0.126 |
| `gui-vet` | 0 | 0.122 |
| `inherited-r1-r3-current` | 0 | 6.916 |

Pristine suite commands are `go test -race -count=1 -timeout 5m ./...` for root,
`go vet ./...`, and the recorded module variants. Exact timeout arguments, cwd and
launcher bounds are authoritative in the JSON receipts.

### Valid original-production sensitivity (not compiled mutants)

The external probe module resolves to a separately extracted exact `6c4561c`
archive for negatives, and back to exact candidate source for restoration.
`go list -m -json ensemble` verifies each resolution. Original source compiles
successfully; no candidate or predecessor production source is patched. These
controls use real public Agent/Job methods, ordinary held Go handlers, actual PTY
processes, and a fake external HTTP boundary—not mocked production functions.

| Recorded control | Exit | Seconds |
|---|---:|---:|
| `corrected-wake-r3` | 0 | 1.484 |
| `original-compile-r3` | 0 | 1.370 |
| `original-mixed-r3` | 1 | 1.114 |
| `original-input-r3` | 1 | 1.204 |
| `restored-wake-r3` | 0 | 1.412 |
| `corrected-literal-limit-r3` | 0 | 1.316 |
| `available-corrected-r2` | 0 | 4.515 |
| `available-original-r2` | 1 | 3.853 |
| `available-restored-r2` | 0 | 4.446 |

- **Mixed join:** corrected source makes three requests, executes exactly one
  fresh effect and one held join, carries both actual messages/correlation and
  paired results, retains the original signed block, keeps handle 2 alive and
  later returns its capped 2011-byte retained output without replay. Original
  source makes no next request while the join is held; after deliberate release
  it runs the skipped tail (three effects). Restored source passes.
- **Real `send_input`:** corrected source preserves a running handle and actual
  nonmatching bytes on hint wake, then retains actual exit 7. Original source
  remains blocked until deliberate process release and then reports done/exit 7
  instead of the required still-running wake. This is a behavioral timing/state
  failure with real completion recorded, not a harness timeout or killed test.
- **Initially available output:** corrected source returns it before ordinary
  work completes and passes later pattern/cursor/kill checks. Original source
  cannot return the available bytes while work is held; after deliberate release
  it returns the truthful terminal output. Restored source passes. Short measured
  assertion windows establish the missing prompt response; no 3600-second test
  timeout is passed off as a sensitivity result.

## Retained attempts and assistance — credit boundaries

The complete retained student attempt ledger was inspected. The first feature
freeze `ccf82d8` genuinely failed inherited PTY behavior because every arriving
byte could finish `send_input` before exit 7. That failure remains a failure.
Main-18 approved the ownership plan; main-19 resolved the documented available-
output compatibility distinction after that failure. The source implementation
was not supplied by the parent. This is disclosed new-feature implementation
history, not an extra formal correction round.

The student's seven original-production behavioral failures/restored pass and
10 focused race repetitions were inspected as **retained student evidence**;
they were not falsely counted as newly run by this reviewer. The reviewer ran
independent original-production scenarios and the complete pristine suite.
Student fixture/type/path/model-ID, launcher matcher/schema, pager and other
setup mistakes remain uncredited. Coverage 87.1% is descriptive, not a gate.

This review also retains its own failures:

1. Initial external fixture used nonexistent `NewEnsemble`/`APIURL` names and the
   wrong root type. The no-suffix probe commands failed compilation; none is a
   valid behavioral negative.
2. The `-r2` retry compiled but omitted required `Config.System`. Runtime fixture
   construction failed before the tested path; its nonzero results are also
   uncredited. The no-test original compile succeeds but proves no behavior.
3. The first available-output cycle used unsupported per-call
   `ai_callback_delay` on `kill_job`. Candidate correctly rejected that input;
   the ensuing missing-reap assertion is not a production finding. The final
   `-r2` adjacent cycle uses the actual documented `handle` schema and preceding
   `tool_limits`, and passes on corrected/restored source. The original available-
   output failure happened before that invalid step, but the final valid cycle
   is the credited comparison.
4. Inspection-path misses and an aborted private orchestration assertion are
   setup/inspection issues, not student failures. A final evidence-exporter schema
   error (string digest versus object digest) was retained and corrected before
   the successful preservation/hash verification. No failed file is overwritten.

The final valid runs use the real API/schema, preserve all earlier receipts and
verify module resolution. No production behavior was mocked and no repair code
was supplied to the student.

## Retained executable proof and limits

The final retained `local-wakeup-r5` proof was read directly: five CLI HTTP
requests and seven workflow HTTP requests, safe stdout, actual retained job
output, paired call dispositions and actual manuscript/review files. The CLI
keeps the original process live while fresh work runs and later captures
`ACTUAL_PROCESS_FINAL`; the workflow executes real author/editor/reviewer effects
while its ordinary thought handler remains live, then EOF joins it and retains
`ACTUAL_ORDINARY_THOUGHT`. Failed r1–r4 launcher receipts remain failed; r5 is the
clean recorded execution. These are **inspected retained exact-binary local
proofs**, not a newly rerun executable campaign.

Both frozen binary stamps were independently read. The original twelve paid
calls belong only to initial `00489b05`, not this amendment or corrected `6c4561c`.
There were no new provider calls, credentials, live GUI, deployment, streaming,
or future fleet/MCP validation. Local HTTP fixtures are not provider acceptance.
Finite race schedules and positive suites do not establish correctness under
every imaginable schedule; they cover the specified ownership/race boundaries
without adding a speculative framework requirement.

All 341 original and 83 correction portable files, all 19 prior excluded local
artifacts, and the accepted R1–R3 report were reverified unchanged. Final source,
chapter, archive, binary and grader checks are indexed with the new evidence.
Private binaries, extracted source and probe code stay outside Git. Only this
report and its evidence directory are this review's repository outputs.

**Final recommendation: ACCEPTED for this exact whole amended candidate.** No
blocking code correction is requested; root owns subsequent acceptance and tags.
