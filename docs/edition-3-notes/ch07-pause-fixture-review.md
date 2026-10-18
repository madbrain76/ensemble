# Chapter 7 one-field pause-fixture repair — independent review

## Verdict and boundary

**PASS — accept the one-field pause-fixture repair. No further repair is required by this review.** The corrected full Chapter 7 route passes on frozen r4 inputs. The repaired evaluator detects a compiled, real tool-admission bypass; exact source restoration passes again.

This accepts **only the grader fixture change**, not the whole student candidate, native/browser presentation, speech behavior, publication, or a tag. Root must inspect this report, exact evidence index/checksums and finalized submission and join the reviewer before acting. Entry, adapter and discrepancy reviews remain separate accepted records.

The production change is exactly:

```diff
-            "command": cmd, "delay": 10,
+            "command": cmd, "ai_callback_delay": 10,
```

Exact comparison to `/tmp/cadence-ch07-pause-fixture-original.go` confirms replacing that single occurrence produces the pinned current driver. No evaluator, score, schedule, timeout, assertion, production application, student source, chapter, guide, API or other grader code was changed by this reviewer.

Chapter 4's public argument table requires `ai_callback_delay`. Chapter 7 requires the current tool to continue, subsequent batch admissions to pause, and another tab's explicit unpause to release the shared gate. The repair lets the fixture exercise that contract. It does not require an undocumented `delay` alias. The original 80/100 failure and diagnostic remain preserved.

## Read coverage and inherited evidence

Reloaded completely: current `AGENTS.md`, coding skill, reviewer skill, execution procedure, readiness, student brief, Edition 3 carryover, full Chapter 7, authorized task, previous sealed diagnostic report and INDEX, current diff and actual root fixture receipts. Recovered the capped AGENTS middle explicitly; skill and Chapter 7 reads used bounded sections. Read the applicable Chapter 4 callback-delay passage and actual actor/driver paths.

Historical workflow, Edition 2 carryover, voice and writing procedure are byte-identical to previously completed full reads in this retained reviewer context. Those are carried full reads, not newly reread documents or substitutes for missing material. The read ledger distinguishes them. This was not an entry redo.

Root's focused historical race and vet receipts were inspected/pinned, **not relabeled as fresh independent race/vet runs**. The reviewer freshly ran the small non-race historical GUI regression below. No 62-second Chapter 1/HTTP or broad Chapter 6/race campaign was launched.

## Frozen inputs and actual GUI provenance

| Input | SHA-256 / identity |
|---|---|
| Repaired `internal/gradeedition2/ch07_driver.go` | `bcdbcd50397d40de60c9fa871cd9d588403efeee34162aa01830cf6c99913b09` |
| Unchanged `internal/gradeedition2/ch07.go` | `2dc39413f88b1fd28059c84b52e9b48661be6438696f921fabe2ddc08d2ddb2d` |
| Accepted Edition 3 adapter main | `1041fe571947f1a17d7e04160a363869c874e9264363c799e0b7d42afa9f06fe` |
| Accepted adapter tests | `55ef77d6191213a9cc11c30dacbe27e94d201c47e5909dbce1c73ac4266525da` |
| Corrected grader `/tmp/cadence-ch07-grader-pause-fixed` | `c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e` |
| Preserved original grader `/tmp/cadence-ch07-grader` | `b72d12ccf29c9d56eeea0109710bb8930923085f9fd68cd80dcf211bf3b10826` |
| Student HEAD, clean | `48e1d12d9c396b5988c0ce6eccdef67e581872ff` |
| Exact frozen CLI | `4beda0aa805195229c5987a025a2b794b4fe3e77a6526ff194927282fddd85cb` |
| Frozen archive | `7beef80a0af5282c793954d738f1df9da51c510051f68bc2ab5454278f304624` |
| Separately frozen student GUI — **not** full-grader target | `39ebcd3e42a57d8972b6278fb8cc66021e38afd41afa2cd230b93499e6b13a81` |

Owned scratch: `/tmp/cadence-ch07-pause-fixture-review-independent`. Source was copied from the previously verified private extraction and all 117 hashes checked. The full route used those source/module identities and the exact frozen CLI. The real grader selected `go build -o <owned temporary output> ./cmd/` in `source/gui`. A transparent Go wrapper delegated to real Go, matched canonical paths, retained the actual built executable with its execute mode, and recorded its hash. It did not substitute a diagnostic GUI.

- Actual full-grade GUI: `991bab3e1cd23dfb84e4e5adcde4cd2bea1379bf53c8b5a6911833b8178c5e5c`.
- Compiled bypass GUI: `1cdd02e6e062d71bfbf831b0a80a01414556c538397b74b2b75298b85887b286`.
- Compiled restored GUI: `d0aee65c932c49d14da792eba641ebeea990d848ca9db96ef94698807d8e5496`.

The restored source is **exactly identical across all 117 files** to initial source. The restored binary is **not byte-identical** to the full-grade binary: an optional equality assertion failed and was retained. Go/module/platform/signature metadata agree, but build IDs differ; this review does not invent a cause. Binaries and commands have separate pins. The restored executable was actually tested after successful compilation; restoration does not rest on a byte-identity claim.

## Exact route and execution environment

CWD: `/Users/bill/projects/ensemble`. Expanded argv is in `full-route-argv.json`. Here `$W` means the owned scratch and `$S` the pinned student worktree:

```sh
/tmp/cadence-ch07-grader-pause-fixed -ch 7 -json \
  -core "$W/source" \
  -consumer "$W/source/examples/ch04-consumer" \
  -consumer-case "$W/source/examples/ch04-consumer/consumer-case.json" \
  -workflow "$W/source/examples/ch06" \
  -stream-exercise "$W/source/examples/ch07" \
  -stream-model claude-sonnet-4-6 \
  -gui "$W/source/gui" -gui-model claude-sonnet-4-6 \
  -media-models /Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-evidence/declared-media-models.json \
  "$S/.ch07-local/frozen-r4/ensemble-cli"
```

The runner constructs an explicit clean environment: wrapper PATH, HOME, LANG, owned TMPDIR, local toolchain, GOPATH, `GOPROXY=off`, `GOSUMDB=off`, and capture phase label. No inherited LLM endpoint/key is passed. Existing grader loopback fake transports remain the only provider fixtures. `run.py`, receipts and `go-invocations.jsonl` record actual argv/CWD/environment/time/exit/output hashes. No watchdog cancellation or process deadline. Behavioral tests were sequential; preparation/compilation did not mutate inputs of a running test.

## Fresh results, including failures

| Run | Exit | Score / max | Passed | Wall time | Meaning |
|---|---:|---:|---|---:|---|
| Corrected full r4 route | 0 | 100/100 | true | 15.268 s | Actual production route, including existing inherited parity; not fresh long-gate proof |
| Build private pause wrapper | 0 | — | — | 1.284 s | Existing repaired driver/evaluator plus capture |
| First pause wrapper attempt | 1 | 0/20 | false | 5.029 s | **Setup failure**, wrong asset directory; not detection |
| Corrected baseline pause wrapper | 0 | 20/20 | true | 0.282 s | Actual full-grade GUI, correct `gui/web/gui` assets |
| Compile gate-bypass application | 0 | — | — | 0.792 s | Successful application compilation before negative |
| Real gate-bypass pause check | 1 | 0/20 | false | 0.243 s | Intended behavior detection with real premature effects |
| Compile exact restored application | 0 | — | — | 0.299 s | All 117 source hashes restored before compilation |
| Restored pause check | 0 | 20/20 | true | 0.278 s | Unchanged repaired evaluator passes again |
| Fresh historical GUI regression | 0 | — | — | 2.190 s | `go test -count=1 ./internal/gradeedition2 -run '^TestChapter7GUI$' -v`; PASS, test duration 1.414 s |

The pause executable prints an actual `PAUSE_REPORT` from the existing `grade.Ch8Evaluate` pause check. Raw run metadata has null score/Passed when stdout includes Go test framing; `measured-results.json` extracts the explicit report, not assumptions. Builds/unit tests have no grader score/Passed. No full-grade score is inferred for the pause-only mutant/restoration.

## Real sensitivity and causal evidence

Only the private source copy was mutated:

```diff
- if !a.admitStart("tool") {
+ if false && !a.admitStart("tool") {
```

This is the actual `dispatchNext` admission guard, not a fabricated grader message, fake error, compile failure, or assertion change. Pause-state handling, transport, jobs, results and evaluator stayed intact. Full source comparison showed exactly `internal/agent/actor.go` changed. `pause-bypass.diff` records it.

The private copied grader adds only capture/workspace retention and a test wrapper calling existing `driveGUIPause` and evaluator. It does not change provider reply, arguments, socket order, pause barrier, fixed wait, command release, assertions, score or unpause behavior. `capture-only.diff` and wrapper source expose that boundary.

At the existing paused checkpoint:

| Observation | Baseline | Bypass | Restored |
|---|---:|---:|---:|
| Dispatch messages on inspected connection | 1 | 3 | 1 |
| Non-error completed results | 1 | 3 | 1 |
| Second-tool effect file exists | false | **true: `GUI_SECOND`** | false |
| Provider requests | 1 | 1 | 1 |

Baseline/restoration subsequently show three successful completed tools, the effect file, and a second provider request **after the other connection's unpause**, ending the turn. Mutant fails before unpause with three successful results and the real second command's filesystem effect. This distinguishes a gate regression from the previous unknown-argument/error-result failure and replay counting. Relevant tool messages are not replay messages; server rows delivered to two peers are not extra jobs.

The evaluator reports its first failure at the three-start assertion. Captured marker/results independently establish a real premature effect, not just an event-label anomaly. `observed-facts.json` retains selected timestamps, identities, result hashes/known markers, stage state and effect facts. Raw GUI/transport captures stay local-only.

## Retained uncredited incidents

All are recorded in `incidents.txt`:

1. Three guessed root receipt filenames absent; actual receipts subsequently read.
2. Initial wrapper passed module root as asset path: startup failed. Guessed `static` path then failed before retry launch, followed by missing-argv/JSON errors. Actual `web/gui` path passed.
3. Optional executable equality assertion failed after exact source restoration and successful build. Build metadata/IDs inspected; restored behavior passed. No equality claim.
4. A guessed builder-source filename absent during metadata inspection.
5. Packaging treated prefixed GUI log lines as JSONL and stopped with `Extra data`; corrected to actual timestamp/direction-prefix format.
6. Guessed publication-baseline pathname absent; exact task pathname then used.
7. Uppercase index glob missed actual lowercase `r4-discrepancy-package-index.json`; exact record and 31 hashes then passed.
8. First report-write tool request refused for missing visible preface; no file was written by that attempt.

None altered production inputs, scores or assertions. Successful and unsuccessful test receipts/stdout/stderr remain in evidence. These issues are not counted as detections.

## Preservation, limits and next step

Final preservation verifies: student clean at r4; all 117 student/restored scratch source hashes match; all 166 production internal hashes match repaired starting pins; accepted adapter/corrected grader pins match; frozen CLI/GUI/archive match; all 31 discrepancy payloads/index intact; previous entry 76-file, adapter 71-file and diagnostic 47-file bundles exact including reports/manifests; twelve publication paths match `/tmp/cadence-ch07-publication-preservation-after-body11.json`; root fixture receipts pinned/unmodified.

Disk monitored: 2.2 GiB initially, 2.0 GiB during packaging. No resource blocker/cleanup. No student/main application/grader/guidance/manuscript edits, staging, commits, tags, children, paid requests, key access, browser/native/microphone work, live8088 interaction, deployment, service restart, installation or shared/manual cleanup.

This establishes that the spelling repair restores the intended valid-input schedule and retains meaningful pause sensitivity. It does not establish every pause interleaving, stale-local/reconnect cause, following-start class, active replay causality, DOM safety, native audio, typing behavior, identity durability or whole-candidate correctness. These remain separate public Chapter 7 obligations. No paid/native permission or student acceptance follows automatically.

**Next step:** root may accept this one-field repair after full artifact inspection and join, then resume its separately authorized workflow. No additional change requested. No pending jobs, tests, questions or children. `ch07-pause-fixture-review-evidence/INDEX.md` gives the frozen exact set/checksum rule.
