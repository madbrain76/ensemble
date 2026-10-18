# Chapter 9 — bounded operational exercise

**Both actual exercises passed.** Exactly one browser-local attempt then one paid Claude CLI attempt ran. No rerun/retry, new browser, paid GUI, provider matrix, production change, rebuild/test/grader, commit/tag or child agent. No owned process/job remains. **An outer browser receipt-wrapper failure is preserved, not relabelled green.**

## Actual results and limits

| Exercise | Actual result | Requests/status | Native cap | Sealed elapsed / wrapper |
|---|---|---|---|---|
| Muted owned Chrome + frozen GUI, local fake upstream | Controller exit 0, `ok:true` | 6 / six HTTP 200 | 4096 | 6.053s / 6.192s; outer wrapper exit **1** (lookup defect below) |
| Frozen CLI, real `claude-sonnet-4-6` | Controller and wrapper exit 0, `ok:true` | 5 / five HTTP 200 | **1024**, unchanged | 15.183s / 15.305s |

Browser authority was one normal muted profile, ≤600 seconds/8 fake requests. Paid authority was one session, ≤8 total requests, ≤60 seconds setup/180 active/240 absolute; measured setup was 0.054s. No automatic next attempt, redirect, fallback, token/body change or retry occurred. HTTP/request counters are actual proxy reservations, including all turns. Only Anthropic POST `/v1/messages` with the pinned model was forwarded in the paid case.

The graph `base → code-tools → helper` activated through real calls. Both branches loaded `code-tools`, wrote **`note.txt = SKILLS_OK`**, read that exact value, marked unload pending and answered verified content. Primary system hash remained identical across requests; next declarations retained activated tools. Browser-local additionally wrote `after-unload.txt = SKILLS_OK` through the still-enabled dynamic tool. Paid made no extra post-unload tool call. Durable call/result IDs, file hashes, pending result, declarations and system hashes are in each exercise’s `sanitized-effects.json` / `proxy-requests.json`.

## GUI, privacy and usage

The real page, production WebSocket and frozen assets were used. Typed messages correlate five dispatch/result pairs, completed turn and idle/unpaused state with five DOM action cards. Helper instructions and `CH09_CALLBACK_VALUE` interpolation are visible; transcript answers the verified file value. DOM facts show three pane bounds, no body overflow, connected/idle, read-aloud false. Informational completed-job notices were not treated as errors. Two settled screenshots were captured, but **not personally visually inspected**; they remain local/hash-only at `/Users/bill/projects/ensemble-edition-3-student/.ch09-local/operational-preparation/execution-runs/ch09-root-execution-01/browser-local/initial.png` and `skills-result.png`. No speech/MIC/audio/focus, restart, second-page or paid GUI credit is claimed.

The read-only Python audit counted **0 credential-file opens in browser-local and 1 in paid-cli**. The latter is the sealed proxy’s approved `~/.cr/settings.json` access after pin/reservation; values were not logged. Python-observed browser destinations were loopback only; paid had five TLS connections to the resolved Anthropic endpoint. This does **not** assert all Chrome/host background networking or writes were absent. Child HOME/environment were isolated, with dummy keys, and Chrome was muted. Raw API/auth/event/CDP traces and the profile are local/hash-only, not portable evidence.

Actual CLI usage was **8,358 input / 729 output tokens**, cache write/read both 0. Five retained provider response usage objects sum to those exact values; per-response outputs were 127, 195, 76, 131 and 200. No price lookup or dollar estimate was made. An optional post-hoc decoder did **not** establish byte equality between reconstructed API-log JSON and raw proxy fingerprints; that negative extraction is preserved and uncredited, not an HTTP/provider failure. Local upstream body-hash equality is established; paid body preservation is supported by the unchanged sealed forwarder and actual request fingerprints, not a claimed independent upstream byte receipt.

## Failures, corrections and transport recovery

1. Root’s original command/approval/run-root did not fit the frozen controller interface. Prelaunch stopped before artifact/key/upstream. Root supplied exact v2 task/approval and confirmed no attempt consumed. Original inputs remain in evidence. A separate read-only asset-path lookup failure was corrected before launch, with its failure retained.
2. The actual browser exercise passed, but new outer `run-once.py` searched `browser-local/browser/*browser-exit.json`; the real receipt is `browser-local/owned-browser-exit.json`. Original wrapper, outer exit 1 and `gate_passed:false` remain unchanged. Root independently confirmed successful actual cleanup and explicitly authorized receipt-only reconciliation (main-18). New paid-only wrapper/diff checks the direct receipt and actual prerequisite without changing sealed functions/assertions/timeouts. There is no global all-green relabel.
3. **Error 707 was an inference/transport interruption after browser completion**, not a browser or paid result. Existing job2843 was collected; durable results were reconciled; no browser replay occurred. The unused paid attempt then ran once under still-valid v2 authority.
4. The earlier unapproved three-provider/paid-GUI preparation was removed after root’s scope correction before execution. Prior bundles remain untouched. Earlier sealed report’s `dependencies` wording is inaccurate: root verified the original failing fixture field was **`depends-on`**. This qualification does not rewrite that report.
5. `CH09_CALLBACK_VALUE` is the accepted interpolation marker; **`SKILLS_OK` is file content**. Its name is not extra proof of public `RegisterVar` callback execution and is not a new product requirement.

## Cleanup, identity and exact evidence

Chrome exit0/listener closed; GUI and paid CLI exit0/process exited; both proxies closed with zero active upstreams and zero retries. All recorded owned PIDs are absent. Source remains `251ce88fb1677897dbcc6b33661096c5cb39cd8e`, tree `a2add41c306d567c05d98ae760911bc13c18b75a`; 132 tracked source files, all five binaries, served assets, sealed scripts/policy/fixtures, prior public-read hashes and prior bundles verified unchanged. Root’s separately attributed attestation reports **33/33 protected tracked inputs unchanged, 0 drift at 2026-10-18T21:54:41.327408+00:00**; student did not read the private baseline. This is time-qualified, not a universal preservation claim.

The exact commands, cwd, sanitized environment, preflight pins, timing, exit status, read-only audit, original failures, v2 authority and root reconciliation are retained. Preparation seal remains `070e03c3703a629fe3b29cd1adadc9292563a6bcf07ce1839aaefe910fb43ab1`; policy remains `b8d8525398503564367003eaa3c34dc73dab68be971238c86377752e213a0e79`. Wrapper/collector scripts are new receipt-only artifacts, not production patches. No private audit/grader input was read; assistance was limited to the exact public tasks, approved scaffolding/preparation and root’s explicit operational corrections/attestation.

`ch09-operational-exercise-evidence/index.json` is the exact self-excluded UTF-8/no-NUL portable inventory and includes this report. `execution/local-artifact-index.json` hashes the owned raw traces/profile/PNGs without copying them. `execution/final-preservation.json` binds cleanup and unchanged identities; `execution/qualifications.json` preserves each limitation. No genuine blocker remains; root’s final inspection/acceptance is still its own decision.
