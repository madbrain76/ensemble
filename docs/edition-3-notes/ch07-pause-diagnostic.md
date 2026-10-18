# Chapter 7 r4 — bounded pause-discrepancy diagnostic

## Verdict

**Reconciled: the measured failure is a harness/input-contract mismatch, not demonstrated failure of r4’s pause gate.** The existing pause fixture passes `delay: 10` to `run_command`; current Chapter 4 specifies `ai_callback_delay`. Frozen r4 rejects the unknown `delay` field. Its three admitted tool jobs fail rapidly, producing error results before the inbound pause command. Consequently the intended first blocking command never provides the scheduling boundary assumed by the grader’s count assertion.

In a private copy, changing **only that JSON field name** to `ai_callback_delay`, with unchanged command values, actor/GUI, timing, two-socket schedule, barrier and evaluator, yields **20/20, `passed=true`**. Real first-command completion happens while paused; second/third starts and successful effects follow explicit unpause from the other socket.

This is not whole-student acceptance, universal pause proof, a production grader patch, or permission to rewrite old results. The unchanged full Chapter 7 grade remains **80/100, `passed=false`, exit 1**. No student compatibility alias or private grader accommodation is justified by this diagnostic.

## Authority, reads and identities

Task `/tmp/cadence-ch07-r4-pause-diagnostic-task.txt`, including final package addendum, verified SHA256 `5f4c4c9d26b0515f6be711130db62bfa33d350be94fc6ffac3ae57212a107a45`. Entry and thin-adapter verdicts remain accepted and their evidence bundles unchanged. Current instructions and Chapter 7 were read; the inherited Chapter 4 public argument spelling was checked directly. `read-ledger.json` distinguishes fresh full reads, overlap recovery, and prior complete historical-role reads verified byte-identical, rather than calling summaries full reads.

Frozen candidate source: commit **48e1d12d9c396b5988c0ce6eccdef67e581872ff**, 117 files. All 117 archive-extracted files and corresponding student files matched the supplied manifest before/after diagnostic work. Student remained clean at that commit. Discrepancy index SHA256 `101e30a07599d680bd4dc6325186e10b16be499cf2d4712c868245770a3f4003`; all 31 indexed payloads verified (165,053 bytes), unchanged at final check.

Local-only pinned artifacts:

| Artifact | SHA256 |
| --- | --- |
| Frozen source archive | `7beef80a0af5282c793954d738f1df9da51c510051f68bc2ab5454278f304624` |
| Frozen CLI target | `4beda0aa805195229c5987a025a2b794b4fe3e77a6526ff194927282fddd85cb` |
| Student frozen GUI | `39ebcd3e42a57d8972b6278fb8cc66021e38afd41afa2cd230b93499e6b13a81` |
| Accepted explicit grader executable | `b72d12ccf29c9d56eeea0109710bb8930923085f9fd68cd80dcf211bf3b10826` |
| Independently rebuilt GUI used for both pause controls | `22dc29b6f31d56208727a624f978de6dced2c3bc37a6690ba94f803ff236c029` |

### GUI provenance qualification

The full grader builds the `-gui` module itself, selecting `./cmd/`; it does not launch the supplied frozen GUI. Exact commands/cwd are retained in `go-invocations.jsonl`. A capture wrapper initially failed to copy the temporary output because macOS supplied canonical `/private/tmp` rather than `/tmp` for cwd. Therefore **no binary hash is claimed for that already-removed original full-grade GUI output**.

A subsequent build used the same recorded `go build -o … ./cmd/` target from the verified source. Both pause controls used that one rebuilt binary. `go version -m` distinguishes the student frozen GUI’s `-trimpath=true`, clean git revision 48e1d12d metadata from the extracted-source rebuild’s no-git/no-trimpath metadata. They have matching verified source and relevant module/toolchain identity, not matching binary bytes. The saved original student full-grade failure and the new same-source full-grade reproduction agree; the student’s frozen-GUI effects probes are inspected carried evidence, not relabeled fresh controls here.

## Finite measurements

Every measured command has cwd, explicit clean environment, duration, exit, output hashes and stdout/stderr in the evidence. There were no outer test deadlines, paid calls, credential values, browser/native actions, children, repeated 62-second gates or race campaign. Only existing local fake transports were used. Tests were sequential.

| Run | Seconds | Exit | Score | Passed | Meaning |
| --- | ---: | ---: | ---: | --- | --- |
| Unchanged explicit full r4 route | 15.609 | 1 | 80/100 | false | Same three-start pause error; all other Chapter 7 checks pass |
| Compile capture-only pause harness | 1.320 | 0 | — | — | Setup succeeds |
| First capture launch | 0.010 | 1 | 0/20 | false | **Setup failure only:** copied GUI lacked execute mode |
| Capture-only retry, unchanged fixture/assertions | 0.250 | 1 | 0/20 | false | Same three-start error, captured causal facts |
| Compile explanation-text variant | 0.737 | 0 | — | — | **Never run**; hypothesis disproved by capture |
| Compile public-argument-only variant | 0.704 | 0 | — | — | Setup succeeds |
| Public-argument-only pause control | 0.283 | 0 | 20/20 | true | Existing evaluator passes with real command effects |

The same-source GUI rebuild also succeeds, exit 0, 0.440 seconds. No new full-grade 100/100 result is asserted. The last control invokes only the existing 20-point pause check, not the rest of Chapter 7.

### Exact causal facts: original fixture

The control’s fresh Agent has two subscriptions before the prompt. Both initial replay brackets are empty: this is not a historical replay contaminating the count. Socket A’s selected events contain exactly three distinct call IDs/display IDs, each followed by an error result:

- `gui_pause_1`, part 2;
- `gui_pause_2`, part 3;
- `gui_pause_3`, part 4.

All return `is_error=true` and `error: invalid arguments: json: unknown field "delay"`. The durable History records three actual admitted job handles, each ending with a 53-byte error report and no subprocess exit code. This is not proof that three intended shell commands executed; they did not produce the blocking/marker effects.

The selected server log records all three dispatch/finish pairs between **21:26:58.327584Z and 21:26:58.328065Z**, before the inbound pause at **21:26:58.328070Z**. Those timestamps include deliveries to both sockets. The evaluator counts only socket A; duplicated server-log deliveries must not be counted as additional jobs. The second-tool marker remains absent after the run. The fake provider’s final success wording is scripted and is not evidence of successful tool execution.

The pause/pong barrier is real and remains unchanged. It orders pause processing after the command, but cannot undo the three already completed invalid invocations or make the missing blocked subprocess exist. The fixture’s inference from queued dispatch messages to “starts allowed while paused” is invalid for this input.

### Exact causal facts: one-field control

The sole behavioral delta in `public-argument-only.diff` is:

```diff
- "delay": 10
+ "ai_callback_delay": 10
```

No explanation text remains in this executed variant. No new API, actor change, sleep, ordering change or weakened assertion is involved.

Captured order:

1. First dispatch at **21:27:40.786524Z**.
2. Inbound pause at **21:27:40.786594Z**; existing pong barrier completes before release-file creation.
3. First real command finishes at **21:27:40.789781Z**, `GUI_FIRST`, `is_error=false`, exit 0, while paused.
4. Existing driver’s paused count and absence-of-second-marker assertions pass, including its unchanged 200 ms observation interval.
5. Other socket sends unpause at **21:27:40.994436Z**.
6. Second dispatch at **21:27:40.994828Z**, then `GUI_SECOND`, exit 0; actual `second` file exists after execution.
7. Third dispatch at **21:27:41.011356Z**, then `GUI_THIRD`, exit 0; turn finishes.

Thus the valid schedule provides more than a passing count: the active tool completes during pause, later effects are held until the explicit other-client unpause, and successful results/marker corroborate execution. The original stale/error messages are not silently filtered to obtain this pass.

## Source and public-contract reconciliation

- Current `book/edition-3/chapter-04.md:383–418` defines `ai_callback_delay` as a JSON number of seconds, ordinary `run_command` argument, and observation bound rather than cancellation. `delay` is not the public spelling.
- Existing `internal/gradeedition2/ch07_driver.go`, `driveGUIPause`, builds its three command argument objects with `delay`; its existing pause/pong barrier and first-release/second-marker checks were preserved.
- Frozen `internal/agent/dispatch.go` recognizes the public limit names; `internal/tools/tools.go` declares public `run_command` schema and uses `DisallowUnknownFields` when decoding tool arguments.
- Frozen actor admits jobs and preserves truthful error results. GUI durable tool facts project to dispatch/result messages. For this failure, capture shows real job admissions with argument-error outcomes **before pause**, not evidence of a new post-pause dispatch or a fabricated replay.
- Frozen GUI/socket source applies pause through the actor-owned gate; neither the GUI nor actor was modified for any measurement.

An early hypothesis was that a tools-only response might be refused for missing preceding explanation. **Measured results disprove it as the cause here.** The explanation-only private variant compiled, but was never executed; it was reverted before the final single-field control. Its diff and build receipt remain as an explicitly unused diagnostic branch. No explanation-policy or projection defect is inferred from this failed hypothesis.

## Student probes and package comparison

The supplied `pause-boundary-probe.mjs`, protocol probe and fixture were inspected, and fresh/warm r3/r4 recordings parsed, not rerun. They use valid `ai_callback_delay`, real filesystem effects and replay-aware live identity counts. The r4 fresh receipt records one live start with first effect present/second absent before unpause, then two live starts and second effect present afterward. Its warm receipt has two historical starts plus one live start before unpause; naïvely counting all three would be wrong there. The warm/fresh recordings contain 105/54 events respectively and name the frozen GUI SHA above.

That is a valid distinction inside those probes, but **the new grader failure is not that warm-replay issue**: its fresh replay is empty and its three original jobs fail their arguments before pause. The probes also use different schedules (including a later second connection and additional pause/disconnect/unpause actions), so their pass was not treated as a substitute for the grader’s exact schedule. The new two-socket valid-argument control supplies the narrower bridge.

## Failures and qualifications retained

`incidents.txt` preserves these uncredited setup/inspection failures: guessed report filename; invocation JSON `argv` instead of actual `command`; local-declaration anchor assumption despite named return values; `/tmp` versus `/private/tmp` capture mismatch; copied-binary execute-mode loss; guessed protocol filename. No failed setup is counted as sensitivity detection. The unused explanation variant and corrected hypothesis are explicit above. Bounded output clipping was recovered where needed; raw provider payloads and full private traffic remain local only.

The full reproduced report and original 80/false receipt are not changed. This diagnostic makes no claim about complete replay-window/270-event handling, every active replay interleaving, all reconnect/stale-cause cases, DOM/speech/native evidence, model policy, live provider behavior, later Chapter 8 workspace or Chapter 10 persistence, or complete Chapter 7 acceptance.

## Recommendation and handoff

Root should separately decide the minimal historical pause-fixture argument repair, preserving the real blocking command, ordered pause barrier, other-client unpause, actual effect checks and failing controls. Independent sensitivity validation and any subsequent corrected full-grade acceptance remain root work, not accomplished here. Do not relax r4’s public schema to accommodate a private historical misspelling, weaken the pause assertions, hide error events, or turn a warm-count filter into a claimed fix for this fresh failure.

At final check: all twelve current publication paths match `/tmp/cadence-ch07-publication-preservation-after-body11.json`; source/archive/package and prior accepted evidence checksums match; main/student/grader/manuscript/guidance/reference files were not edited. Only this report, `ch07-pause-diagnostic-evidence/`, and owned `/tmp/cadence-ch07-r4-pause-independent` were written. No staging, commit, tag, push, cleanup, service action or student test was performed by this auditor.

**No pending jobs or tests.** Evidence INDEX and SHA256SUMS define the frozen exact portable set. No binaries, archives, private raw payloads or keys are included.
