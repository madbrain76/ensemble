# Chapter 7 native browser exercise — root execution

Bill approved the audible exercise and subsequently authorized autonomous retries. This exercises unchanged student commit `48e1d12d9c396b5988c0ce6eccdef67e581872ff`; it is not whole-code acceptance or paid-provider proof.

## Result

The third launcher attempt exited 0. In a disposable headed Chrome profile, the actual production speech queue emitted native start for “Harbor lights shine softly.”, cancellation through the page's Interrupt control, and native start then end for “Quiet gardens rest.”. At the final observation speech was neither speaking, pending nor paused and the GUI speech-pause indicator was released. No native error was used as completion. Human listening remains a separate pending observation. Default voice service locality is unknown and was explicitly consented to; no production voice selection or callback replacement was added.

All seven served assets matched pinned student bytes. One fixed local-provider request was handled; no paid request, credential read or microphone operation occurred. GUI/Chrome PID and listener ownership were checked before attachment. Both owned child processes exited 0 and the local fixture closed. No shared browser or live GUI was attached or restarted. Runtime receipt duration was 2.974 seconds; the enclosing command was 12.963 seconds including harness-process teardown. No latency benchmark is claimed.

## Retained failures and scope correction

Attempt 1 exited 1 before GUI/Chrome/speech: a redundant check of the harness's own process compared canonical Node executable spelling with a symlink invocation. Attempt 2 exited 1 before speech because the harness's blanket DNS block also blocked loopback, triggering a failed navigation/retry. Both receipts are preserved. Removing the redundant self-command check and exempting loopback corrected the harness only. No student source changed. Root also guessed the wrong speech-asset path and then a nonexistent policy root key during report preparation; both reads failed before report writes, and actual manifest paths were used afterward.

The earlier requirement to prove default-voice locality was root-imposed test policy, not the chapter contract. Bill specifically reminded root to avoid invented requirements that obstruct efficient execution of the existing book. Root withdrew that restriction rather than require a student feature. Existing process isolation and fixed local-provider boundaries remain; this is not an OS-network-isolation proof.

## Evidence

`native-result-summary.json` records the exact event sequence. Each attempted launcher and receipt remains separate; `attempt-3-receipt.json` contains full event and process records and `attempt-3-command.json` the actual command, environment and exit. The original student proposal/index remain unchanged. `SHA256SUMS` covers this report and the evidence directory except the checksum file itself.
