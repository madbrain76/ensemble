# Edition 3 Chapter 5 — initial independent whole-candidate code review

## Verdict: REVISE

Candidate `00489b05b65cfe5209f26cb655667df44b3847e2` is not ready for a Chapter 5 acceptance tag. Both required frozen-CLI graders pass, and substantial actor, history, jobs, media, and workflow work is real. Three public-contract failures remain: a pending hint is stranded by a text-only response; a canceled queued `Ask` does not return until unrelated work finishes; and shutdown can join an ordinary handler before stopping another managed process.

This is one initial whole-chapter independent review, not an author pass, adapter re-review, or correction round. The root owns correction assignment, final acceptance, and tags. No student production, tests, manuscript, grader, guidance, or prior reports were edited. No child delegates, provider calls, credential access, installations, live GUI sessions, deployments, commits, or cleanup were used. Scratch mutation and probes stayed outside Git. No reference solution was used.

## Frozen identity and governing inputs

- Main HEAD: `b9304c266ed50c0fd9234d5eaca80b486ca287b3`.
- Student HEAD: `00489b05b65cfe5209f26cb655667df44b3847e2`, clean.
- Student tree: `095e8f8a2559275d94415bcb18bef5b5540dc2b9`.
- Accepted predecessor: `ec4a350f564e820d398f9125fd73a9b3b27129cd`, `edition-3-ch04-r1`.
- Current Chapter 5 SHA-256: `5d3db35e29da8a7fa8b343c01c88b3040339d75c347e845f434f124f3b17b0f2`.
- Candidate receipt: `ba35402e06c212d792da9278c70488a315478b1e61cbd1aec356bcbbd86a85cc`.
- Frozen CLI: `258233df606e71d50aa0103c5865fcd2e383efc471c3d3f8e71fe59f07609863`.
- Frozen workflow: `c5c2fee6aa8c463853c1c6210f7b2a699be8bf6806e550ba0b28d0aa44507114`.
- Source archive: `336b26e12c16c268a1634dd690d4fc9be075aeed4b1f323ff36c42a1436f70df`.
- Student report: `2c47fb072bb6a93d17fa6738c21cdbd0f306c88e32947b0539cb961c596c9020`.
- Student manifest: `950d5c3fb732ad49ff132bbb240e3e0405628c623313cc28da23baa458c6ac19`.
- Student evidence manifest: `7180c6e8ec7d1cc63d53d33e5e70a5c292e8aa4a3ca05c8121597c8f6d47bff0`.

Independently verified all 46 source-file hashes, archive/binary hashes, and all 341 portable evidence entries (847,782 bytes). The 14 excluded local artifact entries remain identified separately, not silently treated as portable evidence. The exact frozen CLI, not a rebuild, was passed to both graders. Suites and public probes used the pristine extracted archive; source-built test executables are not claimed byte-identical to the frozen race-built CLI. The workflow grader builds the supplied workflow module; retained live evidence supplies the separate frozen-workflow provenance.

Fully reloaded current AGENTS, ensemble-coding skill, workflow, Edition 2 carryover, voice, chapter-writing procedure, Edition 3 execution/student/readiness/carryover, reviewer skill, and the complete pinned Chapter 5. Truncated combined reads were recovered by ranges. Also read the full student report and Chapter 5 execution chronology. Current student-first instructions supersede historical answer-first workflow. Main's later Chapter 2 amendment is teaching context, not a retroactive Chapter 5 condition or a claim that this candidate fixed it. The accepted four adapter/guidance edits and unrelated main-tree material were preserved.

## Required findings

### R1 — pending hint is not consumed before text-only turn completion

**Required; actor turn correctness.** Chapter 5 §5.3 says a turn can finish when no further calls **or pending hints** remain (lines 280–283); §§5.3/5.7 require a live hint to become input to the next HTTP request, without changing the frozen request.

`internal/agent/actor.go:443–449` unconditionally ends a successful text-only response. It never considers the ephemera admitted while that request was held. The hint itself is correctly queued and does not mutate the first request, but no consuming request is issued in that turn.

Independent public-API probe, using a real loopback HTTP server:

1. Hold the first request after receiving its actual bytes.
2. Post a hint and wait until the actor's observer acknowledges that exact hint.
3. Release a normal nonempty text-only response, with no tool call.
4. Inspect the private reply and public History state; then submit an unrelated later prompt.

Observed on pristine candidate source: original reply is `pre-hint answer`; only **one** HTTP request occurred; **one pending ephemera entry** remained; the first frozen request correctly lacked the hint. Only the unrelated later prompt caused request two to carry the previous hint. The normal no-hint text-completion positive passed. The existing `TestHeldEphemeraReplay` and hint test force a tool call on the held response, so they do not cover this branch.

Required outcome: a pending admitted hint must be consumed by the appropriate continuation before that turn reports successful completion; it must not leak into an unrelated next prompt. Preserve request immutability, once-only consumption, and dispatched-batch limits. No particular implementation shape is prescribed.

### R2 — canceled queued Ask blocks behind another request

**Required; public caller cancellation.** Chapter 5 “Ask, rebuilt,” lines 587–615, explicitly says `Ask` waits for its own reply, caller cancellation, or shutdown. The one-element reply buffer permits the actor to finish after a caller stops waiting.

`internal/agent/actor.go:121–124` instead performs an unconditional receive from `Submit`. A queued request's context is checked only when `start` reaches it; therefore cancellation cannot release that caller while an earlier request is held.

Independent public-API probe held one actual HTTP request, called `Ask` with an already-canceled context for a distinct queued request, and observed that it remained blocked for **1.0025 seconds** while the first HTTP request was deliberately still held. After releasing the unrelated request, the queued caller finally returned `context canceled`. This is not an HTTP timeout, dropped observer, shared reply, or slow provider conjecture: the controlled first request was the only impediment.

Required outcome: cancellation must release the affected synchronous caller independently of earlier queued work, without stealing another reply, blocking actor completion, or losing an already completed durable result. Retain private buffered replies, save-before-success, and concurrent request isolation.

### R3 — shutdown may wait on a Go handler before stopping a managed process

**Required; shutdown lifecycle ordering.** Chapter 5 lines 616–620 and 730–735 require stopping managed processes, joining workers, preserving actual outcomes of ordinary non-cancellable Go handlers, and saving. Waiting honestly for a Go handler must not leave other controllable processes alive until that handler returns.

`internal/jobs/jobs.go:72–96` snapshots jobs from a map and, for each job, calls its cancel function (if any) and immediately waits for `execDone`. If the first entry is an ordinary handler, that wait prevents the loop from reaching a later process's cancel function. Because ordinary Go work can be non-cancellable, the process can outlive shutdown initiation indefinitely for a reason unrelated to that process.

Independent public-path probe registered a real ordinary handler held on a releasable channel, started it through `CallTool`, and started a real PTY `run_command` process that printed its readiness before sleeping. Root shutdown ran concurrently. In **3 of 4** iterations the managed process was still `running` after 500 ms while the unrelated handler remained held. Releasing that handler immediately allowed shutdown to kill the process and join. The fourth iteration reached the process first and stopped it; this order sensitivity matches the map iteration source. All handlers and processes were released/joined, with no watchdog expiry or leftover child work.

Required outcome: initiating shutdown must stop all managed processes without waiting for unrelated non-cancellable handlers first, while still joining those handlers/report workers and recording their real results. This does not require fake cancellation of ordinary Go functions or any new job framework.

## Independent measured controls

Exact command arrays, working directories, watchdogs, exit codes, output, and elapsed times are retained in `ch05-code-review-evidence/results/`. Watchdogs are distinct from tool observation callbacks. Neither required gate expired.

From the main root, using the accepted adapter:

```sh
go run ./cmd/grade-edition3 -ch 5 -json \
  -core "$STUDENT" \
  -consumer "$STUDENT/examples/ch04-consumer" \
  -workflow "$STUDENT/examples/ch06" \
  -media-models "$ORIGINAL_EVIDENCE/media-models.json" \
  "$ORIGINAL_EVIDENCE/frozen/ensemble"
go run ./cmd/grade-edition3 -ch 1 -json \
  "$ORIGINAL_EVIDENCE/frozen/ensemble"
```

`STUDENT` is the verified clean candidate checkout. `ORIGINAL_EVIDENCE` is the retained candidate evidence root. These are portable aliases, not different builds or sources.

| New control | Exit | Seconds | Result |
|---|---:|---:|---|
| Frozen CLI Chapter 5 | 0 | 159.359 | `passed:true`, 100/100; all seven checks pass |
| Same frozen CLI Chapter 1 | 0 | 127.263 | `passed:true`, 100/100; configured 62 s headers delay and 62 s body delay both complete |
| Pristine root `go test -race -count=1 ./...` | 0 | 138.122 | Root and CLI packages pass, including inherited real lifetime tests |
| Pristine root `go vet ./...` | 0 | 0.183 | Pass |
| Separate ch04 consumer race test / vet | 0 / 0 | 0.123 / 0.130 | No consumer test files; module checks compile/vet, not invented behavioral coverage |
| Separate ch06 race test / vet | 0 / 0 | 1.357 / 0.133 | Real workflow tests pass |
| Public boundary probes | 1 | 1.557 | Plain-text positive passes; R1 and R2 fail behaviorally |
| Mixed shutdown, corrected setup | 1 | 1.899 | R3 reproduced in 3/4 real mixed-job runs |
| Budget baseline, corrected sensitivity cycle | 0 | 5.998 | Positive N, default 200, over-budget pairing and later-turn usability pass |
| Budget one-extra-batch regression | 1 | 5.101 | Real extra effect: 3 effects / 4 requests instead of 2 effects / 3 requests; no paired refusal |
| Budget restored source | 0 | 6.182 | Same public assertions pass again; source hash restored |

The budget probe establishes **dispatched batches**, not an obsolete ten-round limit and not a requirement that policy refusals consume the same counter. With default configuration, it measured 200 real tool effects and the allowed 201st, result-consuming request. With explicit limit 2, two effects plus request three for the final answer succeeded. An over-budget third tool-call response produced zero additional effect and one paired `not executed` error; the subsequent independent turn succeeded. The negative altered only the boundary comparison in the disposable source. This is a behavioral negative, not a compile failure.

The candidate separately bounds explanation-refusal retries. No finding asks that those refusals be counted identically to dispatched batches. This corrects the shorthand in the earlier private entry audit; current chapter text governs.

## Whole-contract review and inherited coverage

- **Ownership and topology:** public aliases/interfaces in `ensemble.go` and the common hub feed an Agent-owned queue and History writer. Engine sends/parses one request rather than owning a sibling turn loop. Cross-package interfaces remain parent-facing; consumer and workflow modules import the public package. This is source review plus the source-root topology and actual module builds, not ownership inferred from a grade label.
- **Admission, requests, observers:** mailbox admission is mutex-linearized and lossless; private replies are buffered and separate from bounded lossy observations. The actor owns interpretation and sends snapshots to workers. Concurrent-submit, stalled-observer, state/config, metadata, and save-failure tests pass under race detection. R2 qualifies the synchronous adapter's cancellation claim; R1 qualifies live-hint completion.
- **History/provider parity:** request recording precedes held HTTP, event order drives replay, tool calls/results remain paired, provider-private opaque material stays typed and out of ordinary display. Actual captured-request/replay tests and all three provider envelope checks pass. The public hint counterexample does not invalidate frozen-request immutability, but does invalidate universal hint-continuation correctness.
- **Jobs, interruption, shutdown:** callbacks remain reports rather than deadlines; retained handles, PTY input, cursor/pattern behavior, one-shot limits, and honest ordinary-handler results are exercised through the public API. Interrupt/new-turn and stale-result tests pass. Exact frozen live evidence additionally shows a new answer before the old process's terminal report. R3 is the uncovered mixed ordinary/process shutdown case; the individual-job shutdown positives remain real.
- **Explanation policy:** existing root tests exercise false/true `Config.RequireToolExplanation` plus selected `LLM_REQUIRE_TOOL_EXPLANATION` values through public paths, including real side effects, refusals, ordering, and later valid explanation. These passed in the independent pristine race suite. No monkeypatch policy shim was introduced.
- **Media:** explicit model rows, unknown-false capability behavior, text without a declared row, actual `RefPath` bytes, unsupported-model pre-HTTP refusal, redactions and provider render/parse behavior are present in source and covered by the relevant local tests/grader. Retained live image evidence is specifically Gemini PNG, not live acceptance of every modality/provider combination.
- **Actual workflow:** separate `examples/ch06` has three sequential roles with required tools `write_draft`, `word_count`, `read_draft`, `edit_draft`, `review`, durable manuscript/review effects, and reliable public `Ask` handoffs. Source, new module tests, grader, and retained safe events/artifacts establish more than names or log markers. R2 still affects that common synchronous API when a caller cancels while queued.
- **Scope restraint:** no future streaming engine, GUI feature, grouped parallel runner, universal schema, coverage percentage, or abstraction framework was imposed. The existing GUI module was inspected as inherited source; no live GUI ran. Earlier adapter controls were not repeated or relabeled as new candidate controls.

## Retained live and assistance evidence: read, not rerun

Verified the indexed evidence and read actual safe human output, actual guarded job output, request budgets, local reconciliation, workflow manuscript/review artifacts, image output, unsupported refusal, and local remaining-only rehearsal. Read only `.safe` event logs, not opaque provider history; derived event-only facts and their hashes are saved in this review evidence.

The frozen candidate's retained live runs show 4 Anthropic `claude-sonnet-4-6` calls, 7 OpenAI `gpt-4.1-2025-04-14` Chat Completions calls, and 1 Gemini `gemini-3.8-flash` PNG call: all twelve HTTP 200, without retries or model substitution. Human CLI exited 0. Its original launcher summary exited 2 because an overly broad sentinel search matched the assistant's benign denial. That failed summary is preserved; the authorized reconciliation was local, not a paid rerun. Actual guarded process output was retained and contains no leaked sentinel. Safe events place the new answer **18.958662 seconds before** the old job's terminal report. A subsequent real wait returned the completed job output.

Workflow evidence records author `write_draft`/`word_count`, editor `read_draft`/`edit_draft`, reviewer `review`, 78 to 89 manuscript words, a changed manuscript, and a stored review with final count 89. The image answer identifies left red/right blue. Unsupported media exited 1 with zero provider requests. These are carried candidate demonstrations, not new review executions; they do not cover all live modalities, token streaming, or a native GUI.

The execution chronology retains: the pre-implementation rejected broad plan; approved narrower plan; parent's wrong tool-name guidance corrected after student challenge; an implementation-stage observed 75-point envelope failure; disk-pressure mutation attempts that were inconclusive; later six actual behavioral mutations; two transport interruptions; credential-field corrections before calls; strengthened evidence checks without hidden production repairs; and corrected tree/archive shorthand. Earlier-binary Chapter 1 passes are not used as proof of this candidate. No formal code correction round had occurred before this review.

The `unittest.mock` route-control script remains explicitly **uncredited**. Real frozen-executable fake-HTTP rehearsal and dispatch inspection support the remaining-only route. Dummy loader setup checks are not production acceptance. None of these saved scripts was executed against paid providers by this reviewer.

## Failures, setup misses, and limitations retained

- R1/R2 public failures and R3 mixed-job failure are actual pristine-source behavior.
- First mixed-job probe exited 1 during setup because its scratch Config omitted required `MaxTokens`. It is not shutdown evidence. The corrected retry is independently identified.
- First budget mutation used the wrong source anchor (`turn.batches` instead of `t.batches`); its assertion failed before changing source, and the shell continued. Consequently its “mutant” test passed on baseline bytes. That attempt is **uncredited**, not a successful sensitivity check. Its original outputs and explanation remain. The guarded second cycle verified distinct mutant bytes, behavioral failure, exact restoration, and pass.
- Several initial source/evidence path guesses produced read errors; the actual manifest paths were then used. These inspection misses are not student test failures. No failed attempt was relabeled as a candidate regression or deleted.
- No transport interruption was counted as a test failure. No new timed control expired. No paid-call retry was performed.
- Free data volume stayed approximately 3.3–3.4 GiB. Builds were bounded, shared caches were not deleted, and unrelated artifacts were preserved.
- This is finite testing and source review, not exhaustive schedule/model coverage. The package-wide race pass does not prove every concurrent lifecycle ordering; R3 demonstrates why.

## Smallest next step and re-review boundary

Root should assign one bounded student correction round for R1–R3 against the existing public contract, with focused regressions that fail on this candidate and pass on the corrected candidate. Preserve the successful architecture, batch semantics, policy, jobs, history, media and workflow behavior; no grader edits or framework are requested. Do not send private probe source or repair code as student input.

After a new exact clean candidate is frozen, re-review these findings, adjacent actor/shutdown paths, unchanged-input identity, real sensitivity, the required Chapter 5 and same-binary lifetime gates, and inherited regressions. Root decides any live re-demonstration need under separate authority; this review authorizes none. No Chapter 5 tag or overall acceptance should be issued for the current candidate solely because its two grader scores are 100.
