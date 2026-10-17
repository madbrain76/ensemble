# Independent whole-candidate review evidence

Verdict and findings are in ../ch05-code-review.md. This package contains 17 raw command records (stdout/stderr/metadata), source/input identity, safe-event-derived facts, scratch sensitivity provenance, and final verification. It does not contain binaries, opaque provider history, private probe implementation, student repairs, or historical reference code.

Portable aliases: `<MAIN>` = main repository; `<STUDENT>` = exact clean candidate; `<ORIGINAL_EVIDENCE>` = original frozen candidate evidence directory; `<REVIEW_SCRATCH>` = private review workspace; `<SYSTEM_TEMP>` = OS temporary directory. Original local paths remain in the private launcher records.

`budget-mutant` from the first cycle is UNcredited because the mutation anchor assertion failed before editing; its passing exit is preserved. Only `budget-*-r2` is the successful pass/fail/restored-pass sensitivity cycle. `mixed-shutdown-probe` is a setup rejection; `mixed-shutdown-probe-retry` is the actual behavioral failure. The two other behavior failures are in `public-boundary-probes`.

Retained student live evidence was read and hash-verified, not rerun. Safe-event facts exclude opaque history. No provider credentials were loaded by the reviewer. No original failure was repaired or replaced.

Verify SHA256SUMS from this directory. Manifest excludes itself.
