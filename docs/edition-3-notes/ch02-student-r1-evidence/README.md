# Chapter 2 correction round 1 evidence

Candidate: **89cc874b3c32e65ee43f44a72ad9ca48c9dd3323**.

Read [the round-one review](../ch02-student-review-r1.md) for exact command/cwd/status tables, R1–R3 implementation details, assistance, development failures, and limits. The original review and evidence directory are unchanged.

## Identity and final checks

- `final-identity.txt`, `final-clean-recheck.txt`: clean candidate, source archive hash, binary hash and embedded revision (`vcs.modified=false`).
- `student-cli`: exact new executable used by both final graders and local CLI evidence checks.
- `final-grade-e2-ch02.{json,stderr,exit}` and `final-grade-e3-ch01.{json,stderr,exit}`: both exit 0 and passed:true. Edition 3 includes both actual 62-second response phases.
- `final-race.{txt,exit}`: full race tests, exit 0.
- `final-coverage.{txt,exit}`, `final-cover.out`, `final-coverage-functions.txt`, `coverage-denominator.json`: 749/838 source statements, 89.4%; no grader/live-process coverage claim.
- `final-source-checks.txt`: vet, diff check and package import listing.
- `import-smoke/`, `import-smoke.{txt,exit}`: external module consumes the real public root and checks ownership, exit 0.
- `gui-checks.txt`, `gui-stub`: isolated placeholder test/build, not a GUI runtime claim.

## New regressions and sensitivity

- `regressions-before-fixes.{txt,exit}`: all three new regressions fail before the fixes.
- `regressions-after-fixes.{txt,exit}`, `final-focused-regressions.txt`: passing fixed-source runs.
- `check-sensitivity.py`, `sensitivity-run.{txt,exit}`, `sensitivity-summary.json`, `R{1,2,3}-sensitivity.txt`: each fix independently reverted **only in scratch**, behavioral failure followed by restored pass. Candidate source hashes remain unchanged.

## New local-only CLI evidence

- `check-offline.py`, `offline-checks-run.{txt,exit}`, `offline-checks-summary.json`: final local-only run, exit 0. No credential reads, no new paid calls.
- `dump-original-*`: fresh new-binary dumps equal the three original saved live logs.
- `replay-*-new{1,2}.*`: nine deterministic cross-vendor render pairs under the new binary.
- `replay-*-old-binary-now.*`: newly executed **offline** old-binary baselines, not new or re-dated live evidence. All match the corresponding new renders.
- `fake-R1-*`, `R1-sequential-idless.*`: synthetic sequential call/result name check through actual CLI.
- `fake-R3-*`, `R3-*`: synthetic two-field opaque replay and other-model exclusion through actual CLI.
- `redaction-*`: unchanged original fake redaction fixture rendered with the new CLI.
- `fake-R2-empty-*.jsonl`, `R2-empty-*`: three local fake HTTP responses with measured 11/2 usage. CLI failures are expected, measured diagnostics survive, and logs contain error facts rather than fabricated completed responses.
- `original-live-attribution.json`: original six calls remain dated to 2026-10-16 21:01:06–21:01:15 UTC and commit **74eb0265bd8bdbc7857d579704300b317b2365da**. This correction binary has not been newly live-validated.

The final offline replay stage made zero HTTP requests. The subsequent fake failure-capture stage made exactly three local fake requests. No raw provider payloads, auth headers, or credentials are archived.

## Preserved failed attempts

- `edit-anchor-failure-regressions.*`: an edit anchor failed before writing source; regressions remained failing.
- `first-attempt-*`: original local harness/script incorrectly supplied `--protocol`; CLI rejected it. These records are not final pass evidence.
- `second-attempt-*`: retry stopped because the GUI build left its executable untracked. That generated executable was moved here as `gui-stub`; source remained unchanged, then the final run passed.

The old live launcher regex failure remains in the **original** evidence directory and was not overwritten or reclassified. The new manifest hashes retained round-one files and source; the original manifest is separately verified.
