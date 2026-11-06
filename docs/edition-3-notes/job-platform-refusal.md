# Fix: a command this platform cannot launch is ordinary tool data

Branch `fix/platform-refused-job-command`, based on `fix/linux-path-case`. Coordinator, 2026-11-03.

`Job.Start` caught a fixed list of exception types and let anything else escape unrecorded. `agent/Jobs/Pty.cs:10` gates the PTY launcher on macOS, so on Linux `Job.RunCommandAsync` throws `PlatformNotSupportedException`, the worker task faulted with no recorded result, and teardown rethrew it as "Job shutdown could not record or close every resource" — which then aborted `Ensemble.Destroy` for the whole application. Any consumer that started such a job could not shut down cleanly.

The refusal now goes through the same path as any other tool failure: `AppendOutput("Tool failed: ...")` and a recorded result, so the model can read why the command did not run. A platform that cannot launch a command is a fact about the request, not a defect to propagate.

Effect on the suite, stated precisely because the count is unchanged: the four Linux-expected failures stay four, but they are no longer teardown aborts. Before this change `Chapter 3 tools and loop`, `Chapter 4 jobs and PTY` and both `Chapter 19 confined process` cases failed with `Agent cleanup failed (Job shutdown could not record or close every resource ... macOS system PTY launcher)`. After it they fail on their own real assertions — `nonzero exit is ordinary data`, `real PTY with dimensions`, `Confined environment not allowed...`, `Could not locate toolchain cache...` — which is the honest Linux gap, and `Ensemble.Destroy` now completes. Verified by full `agent/Tests` runs on the base and this branch: 341 passed and 4 failed in both, with only the failure identities changed.
