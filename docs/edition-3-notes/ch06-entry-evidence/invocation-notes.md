# Retained invocation errors (not grading failures)

1. The first historical CLI comparison attempted nonexistent `-chapter` and `-core` options. Exact argv, stderr and exit 2 are retained in `positive-old-cli-source.*` and `positive-old-cli-binary.*`. The actual CLI uses `-ch 6` and a positional build target. Corrected comparisons are separate `*-valid.*` files; failed attempts were not overwritten or counted as grades.
2. A read-only shell batch ended with `cmp docs/edition-3-notes/ch06-entry-evidence/api/main.go <(gofmt docs/edition-3-notes/ch06-entry-evidence/api-source-initial.txt)` under `/bin/sh`. It exited 2 with `/bin/sh: -c: line 31: syntax error near unexpected token '('`. This occurred after saving the initial source text and did not modify any grader/application/student. The portable Python/subprocess comparison in `verify-harness-format.*` then completed with exit 0. That comparison establishes exact gofmt equivalence, not a new binary build.
3. An informational disk-summary read used nonexistent metadata field names `disk_free_before` / `disk_free_after` and printed `None`; authoritative fields are `free_bytes_before` / `free_bytes_after`. The final minimum is calculated from the actual fields in `freeze-status.json`.

No implementation was repaired to hide a failed reference. Both real mutation failures and both exact restorations are retained separately.
