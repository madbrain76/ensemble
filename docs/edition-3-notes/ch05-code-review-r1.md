# Chapter 5 — independent code correction round 1

**Verdict: R1–R3 closed under pre-amendment contract; new authorized job-wakeup amendment not evaluated.**

This is closure of the three findings in `ch05-code-review.md`, not a second entry audit, adapter review, new-feature review, acceptance tag, or unilateral chapter acceptance. The root orchestrator owns acceptance and phase progression. This is code correction round 1 of the existing three-round limit. No further correction is requested for R1–R3.

## Exact scope and identities

| Item | Verified identity |
|---|---|
| Main HEAD | `b9304c266ed50c0fd9234d5eaca80b486ca287b3` |
| Corrected student HEAD | `6c4561cbd6cbae080d952da34a5fcd588a76920f` |
| Corrected tree | `8ba148fa32b5211b0f148e9da137544523374720` |
| Original Chapter 5 production used for negative controls | `00489b05b65cfe5209f26cb655667df44b3847e2` |
| Pre-amendment adjudication chapter SHA-256 | `5d3db35e29da8a7fa8b343c01c88b3040339d75c347e845f434f124f3b17b0f2` |
| Corrected frozen source archive SHA-256 | `56535661fcfdbc5121384353aa02475bec791ac04fad4d57154856216dc81bb9` |
| Exact corrected CLI SHA-256 | `9bd612504b9856fd4b64e0b3b386de965eae6403b0c60c00cd14cd3973bc6d61` |
| Frozen workflow SHA-256 | `acfab0d19f435529f551331f065098b6341a0d57384dce17146401733723fdf4` |
| Student correction report SHA-256 | `9b6dca6bee16bf58eb739578f2974e2cecfe44ebc8cfe167ad85e95692662966` |

The clean candidate's 47 indexed source files were independently checked against Git, worktree, and the exact extracted source archive. All 83 portable student evidence files and all five excluded local artifacts were independently hash-verified. The complete 114-line student report and both manifests were read. Full required reviewer prerequisites were reloaded, with truncated ranges recovered.

The adjudication text is the saved `/tmp/cadence-ch05-before-job-wakeup.md`, not the moving current chapter. The current Chapter 5 amendment was read and recorded (observed hash `0aa8a2060025b59c538a8cb6f0a389428ae2b6673382fb7433f1496f16f0c913`). Bill's newly authorized uniform-job / foreground-observation wakeup / same-ID resumption requirement is **not** a retroactive finding against this candidate. This report does not evaluate or accept that new feature. The current guidance and main working changes were left untouched.

Production changes are confined to `internal/agent/actor.go` and `internal/jobs/jobs.go`; the only other candidate change is the new 417-line `chapter5_correction_test.go`. Engine, History, renderer, providers, CLI, public module interface and workflow production files are unchanged. The original reference solutions were not used. Negative controls use the student's own exact original Chapter 5 production archive, whose 46 indexed source files were separately verified.

## R1 — same-turn pending hints: closed

**Pinned contract:** the turn may finish when no further calls **or pending hints** remain (line 566); request snapshots and live hints remain separate (lines 560–606).

The corrected `onModel` checks pending ephemera in its no-tool-call branch and sends the next request on the existing active turn. It does not rewrite the frozen in-flight request, start a fresh user turn, or reset the dispatched-batch counter. Existing save-before-successful-completion and request-sent hint consumption remain in the owner path.

Independent external-package tests use a real loopback HTTP handler held across public `Post(Hint)`, not a mocked Agent or transport method:

- Corrected baseline and restored baseline: first frozen request contains no hint; the original turn performs request 2, consumes the hint, returns the post-hint answer, and has zero pending ephemera. An unrelated later prompt produces request 3 without replaying the consumed hint. The successful answer is already in saved history.
- Exact original production: exit 1. The original turn returns the pre-hint answer after only one request with one pending ephemera item; the unrelated next prompt consumes that hint. The no-hint positive still passes.
- No-hint positive: one request, correct answer, no ephemera.
- Adjacent budget cases pass on corrected, original and restored source: configured two dispatched batches plus the third result-consuming request; rejection of an over-budget batch with exactly two effects and a paired not-executed result, followed by a usable later turn; default 200 dispatched batches plus request 201 consuming their results.

The public Config retains zero as the default selector, so a logged `effective_limit=0` field in the private probe is not a zero dispatch budget: the actual measured default is **200 effects / 201 requests**. This review does not revive the old refusal-count wording or impose the old ten-round cap.

## R2 — canceled queued callers and private completion: closed

**Pinned contract:** `Ask` waits for its own reply, caller cancellation, or actor shutdown; successful private completion remains buffered and durable (lines 600–617).

The corrected `Ask` normalizes nil context, retains the request's private reply channel, observes caller cancellation/shutdown, and checks that channel again before returning a cancellation/shutdown error. An already-buffered durable result is not deliberately discarded. Cancellation does not require the actor to send to an unbuffered abandoned waiter or cancel unrelated work. The actor still owns eventual completion.

Independent results:

- Already-canceled caller behind a held unrelated HTTP request: corrected return in 8.167 microseconds; restored return in 12.042 microseconds. These are observations, not new timing requirements. The unrelated first HTTP remained held.
- Exact original production: exit 1; the already-canceled caller remains blocked for the full 1.001522459-second observation and returns only after the unrelated request is released.
- Cancel-while-queued test: corrected and restored callers observe the real context deadline while the first HTTP remains held; original production remains blocked through the additional 500 ms observation. The canceled queued prompt never reaches HTTP, and the active request completes correctly after release.
- Two later submitted prompts retain distinct buffered results while neither private channel is consumed. `Drain` completes, both answers are already saved, and the exact replies remain readable after caller cancellation and root shutdown. A new post-shutdown `Ask(nil, ...)` reports `ErrStopped` rather than stealing an earlier success.
- The new student public tests additionally use real context-channel observation to establish admission/cancellation ordering, concurrent caller isolation, and persisted successful replies. They run in the pristine race suite.

## R3 — process cleanup before ordinary-handler joins: closed

**Pinned contract:** shutdown stops managed processes, joins execution/report workers, preserves true ordinary Go outcomes, and saves (lines 616–617 and 727–734).

`Table.Shutdown` now closes admission, snapshots jobs, requests managed-process termination in the first pass, and only then joins cleanup/execution for all jobs in a second pass. Ordinary handlers still return their true result; their status is not fabricated as killed. No new public job API, scheduler, retry framework, or unrelated lifecycle mechanism was introduced. Existing actor/root shutdown ordering and late-report persistence were inspected adjacent to the diff.

The strengthened independent test launches four real held Go jobs and a real PTY shell with a live `sleep` child. It checks both PIDs initially exist and share a process group. Before releasing the Go handlers, it checks process completion and eventual `ESRCH` for both shell and child. Shutdown must still be waiting for the real Go handlers. After release, every true Go return, output file and saved terminal `job_report` is checked, along with final process disappearance.

- Corrected baseline: all four layouts clean up the process group before Go release; shutdown then joins and retains all true results/reports.
- Exact original production: exit 1; **three of four layouts** leave both shell and child alive while cleanup is blocked behind the ordinary handlers. **One favorable layout passes** and is retained; old map iteration is not claimed to fail deterministically on every run. All layouts are cleanly released/joined afterward.
- Restored corrected source: all four layouts pass again.

The student's own retained R3 sensitivity likewise records a favorable original ordering and later original failures. Its report did not erase that favorable observation. This review establishes stop/reap-before-held-Go-join behavior for these real local process groups; it does not claim exhaustive scheduler proof.

## New measured checks

All commands ran with credential-free environment selection, fresh scratch HOME, offline dependency resolution, and explicit cwd. Expensive builds ran sequentially. `run.py` records command, cwd, wall time, child exit, watchdog and raw stdout/stderr. Short tool callbacks were observation handles, not test deadlines.

| New independent check | Exit | Wall seconds | Result |
|---|---:|---:|---|
| Exact frozen CLI Chapter 5 | 0 | 159.089 | `passed:true`, 100/100, all seven checks |
| Same exact CLI Chapter 1 | 0 | 127.247 | `passed:true`, 100/100; both real 62-second header/body delays complete |
| Corrected external public probes, guarded retry | 0 | 6.467 | All six top-level tests, including budget subcases, pass |
| Original-source compile-only gate, retry | 0 | 1.526 | Successful compile; no behavioral credit by itself |
| Exact original R1 subset, retry | 1 | 5.160 | Genuine pending-hint failure; no-hint/budget positives pass |
| Exact original R2 subset, retry | 1 | 1.885 | Both real queued cancellation cases fail |
| Exact original R3 subset, retry | 1 | 2.712 | Three blocked layouts, one favorable layout |
| Restored corrected external public probes, retry | 0 | 6.384 | All pass again |
| Pristine corrected root `go test -race -count=1 -timeout 5m ./...` | 0 | 137.708 | Includes inherited tests and all new correction tests |
| Pristine corrected root `go vet ./...` | 0 | 0.237 | Pass |
| Separate consumer race test / vet | 0 / 0 | 0.178 / 0.128 | Module builds; no test files, not an additional behavior claim |
| Separate workflow race test / vet | 0 / 0 | 1.698 / 0.125 | Workflow tests and vet pass |

The lifetime JSON proves its configured delayed-header/body cases; it does not claim an extra cancellation test. Caller/actor cancellation evidence is separately described above and in the race suite.

From the main repository root, the exact gate commands were:

```sh
S=/Users/bill/projects/ensemble-edition-3-student
F=/tmp/cadence-ch05-correction-r1-evidence
go run ./cmd/grade-edition3 -ch 5 -json \
  -core "$S" -consumer "$S/examples/ch04-consumer" \
  -workflow "$S/examples/ch06" -media-models "$F/media-models.json" \
  "$F/frozen/ensemble"
go run ./cmd/grade-edition3 -ch 1 -json "$F/frozen/ensemble"
```

Both commands received a 300-second watchdog and finished normally. The Chapter 5 grader builds its source-selected consumer/workflow; that is not a claim that its workflow run executed the separately pinned workflow binary. The workflow binary and archive were verified for identity. Pristine suite/probe builds are separate source-build evidence, not asserted byte-identical replacements for the frozen race-built CLI.

All external probe runs used `go test -race -count=1 -timeout 60s -v`, with the exact selected test patterns in the retained metadata. Before each credited corrected/original/restored selection, a successful `go list -m -json ensemble` record verifies the resolved source directory. Neither production tree was mutated: only the private external module's replacement path changed.

## Failures and limits retained

- Initial reviewer manifest helper used the nonexistent key `student_head` instead of actual `source_commit`: setup-only `KeyError`, before any gate. Its raw receipt and corrected verification are retained.
- The first six external probe/compile attempts each exited 1 immediately because the private module lacked its indirect PTY requirement (`go: updates to go.mod needed`). **No tests executed; none is negative-control credit.** The first queue also checked the logging launcher's return rather than its recorded child exit, so it continued through those setup failures. A single bounded retry corrected only private module setup, reads actual recorded child exits, verifies source resolution, and preserves all first attempts. Source-path spelling ambiguity was also removed rather than assumed harmless.
- The student's retained first public fixture expected the wrong report representation; its first CLI capture expected a physical line without leading whitespace and was terminated with exit -15. These are explicitly uncredited setup attempts. Their later local interface receipts are carried student evidence, not fresh reviewer CLI-interaction runs.
- No new provider call, credential load, GUI launch, deployment, install, restart, cache cleanup, child delegate, source repair, commit or tag occurred. The prior whole-candidate review's live observations belong to the earlier exact binary; this round does not relabel them as live verification of the new binary.
- No watchdog expired in the 23 completed command records. Source-resolution commands and setup failures remain separate from behavior tests. Disk observations ranged from about 3.3 GiB to a final 3.1 GiB free; expensive builds were sequential and no shared cache was deleted.
- This bounded round did not repeat the accepted adapter campaign or every historical control. It reviewed the changed paths and adjacent ownership, hint snapshot/consumption, acknowledgment, budget, shutdown/report/save behavior, with paired public gates and pristine regression suites.

## Evidence and preservation

`ch05-code-review-r1-evidence/` contains sanitized portable raw command records, the original failed-helper receipt, complete command index, input/probe hashes, source-resolution receipts, final preservation verification and an integrity manifest. Private probe code, original/corrected source archives, binaries and process artifacts remain outside Git; no repair fixture is supplied as a student input. The student and all previous reports/evidence remain untouched. Current grader adapter hashes match the accepted adapter; inherited grader source is unchanged.

Root acceptance remains separate. **R1–R3 closed under pre-amendment contract; new authorized job-wakeup amendment not evaluated.**
