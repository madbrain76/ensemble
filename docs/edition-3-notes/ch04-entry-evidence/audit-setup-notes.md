# Audit setup events (not candidate failures)

- `rg` was unavailable; ordinary grep/find were used.
- A read-only Git diff opened pager PID 78582 under owned Git PID 78581; only that observed pager was terminated. Later Git commands use `--no-pager`.
- Initial scratch driver used the wrong import module (`ensemble/...`); build exited 1. See `build-driver.log`. Corrected driver uses `github.com/waywardgeek/ensemble`, `ReadConsumerCase`, and the existing `grade.Report` API.
- Initial lifetime outer timeout was 120 seconds, insufficient for the two sequential 62-second probes. `reference-lifetime.log`/runs.json retain the timeout and SIGKILL. Only that inconclusive control was rerun with 180 seconds.
- Initial cap mutation stage stopped at a source-anchor assertion before mutating anything. Corrected anchor is the actual `limits(raw, base)` call.
- Previously recorded controls were not repeated; the corrected runner resumes from records. The script is the final corrected recipe, while runs.json retains earlier failures.
- During the audit, concurrent edits appeared in live cmd/grade-edition3/main.go and an untracked main_test.go. Neither was modified or tested here. All grader executions used the exact bac07e6 Git export created before those edits.
