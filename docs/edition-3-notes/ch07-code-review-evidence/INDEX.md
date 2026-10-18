# Chapter 7 initial whole-code review evidence

**Verdict: CHANGES_REQUIRED for two findings only.**

- R1: an explanation-refused call has zero actual core admissions/effect, but is labeled dispatched by the GUI live and on replay.
- R2: a tool-result card's speaker omits its readable result.
- C1 page-hide release is a measured non-blocking student-task/contract-clarity note under root's explicit classification; it is not a third required correction.

Report: `../ch07-code-review.md`. Student remains frozen at `48e1d12d9c396b5988c0ce6eccdef67e581872ff`, tree `0f5f7fbe4273e91e1be3aa13d0715a99f671eaf9`. Formal correction rounds remain 0/3; root alone accepts/tags.

## How to inspect

1. Read the complete report, then `contract-classification.md` and `incidents.md`.
2. `logs/` contains exact commands, cwd, clean environment, timings, exits and safe stdout/stderr. Null score/Passed fields for Go test/vet are not grader grades.
3. `page-contract-checks*.mjs` are non-browser production-module checks; the first source/setup failure is retained. They do not claim native listening/rendering.
4. `reviewer_contract_test.go` is one private added test of actual public core behavior and GUI projection, with an explained positive and refused negative. Compilation succeeds; the intended projection assertion fails.
5. `module-gate-plan.json` and `run-summary.json` record the completed finite package checks, including the 139.72-second full root race gate. No new Chapter 1 CLI long-gate run occurred; `carried-evidence.json` identifies the exact r4 100/100 positive and native/corrected-Ch7 positives as carried.
6. `read-ledger.json` and `changed-source-read-ledger.json` distinguish full actual reads, recovered caps, surrounding source inspection and unchanged-source carryover.
7. `final-preservation.json` verifies input pins, 117 frozen files, clean student, predecessor, 13 runtime paths, 12 publication files, four exact prior bundles, and 33 native/31 discrepancy indexed files.
8. `private-generated-files.json` and `local-only-artifact-pins.json` are hash inventories only. The generated raw logs, binaries, source archive and native/provider payloads are not included.

No production edits, browser/native/live-runtime action, paid provider call, key access, children or cleanup occurred. No commands/tests remain pending. The private scratch path is `/tmp/cadence-ch07-whole-code-review-independent`; local paths are provenance, not portable binary payloads.

## Exact set

This directory contains exactly **70 files**, including this index and `SHA256SUMS`. The manifest covers the report plus every evidence file except itself (70 entries). The manifest's own hash is supplied in the final submission; there is no circular self-hash.

- `INDEX.md`
- `SHA256SUMS`
- `carried-evidence.json`
- `changed-source-read-ledger.json`
- `contract-classification.md`
- `final-preservation.json`
- `incidents.md`
- `local-only-artifact-pins.json`
- `logs/compile-refusal-projection.run.json`
- `logs/compile-refusal-projection.stderr.txt`
- `logs/compile-refusal-projection.stdout.txt`
- `logs/gui-race.run.json`
- `logs/gui-race.stderr.txt`
- `logs/gui-race.stdout.txt`
- `logs/module-examples-ch04-consumer-race.run.json`
- `logs/module-examples-ch04-consumer-race.stderr.txt`
- `logs/module-examples-ch04-consumer-race.stdout.txt`
- `logs/module-examples-ch04-consumer-vet.run.json`
- `logs/module-examples-ch04-consumer-vet.stderr.txt`
- `logs/module-examples-ch04-consumer-vet.stdout.txt`
- `logs/module-examples-ch06-race.run.json`
- `logs/module-examples-ch06-race.stderr.txt`
- `logs/module-examples-ch06-race.stdout.txt`
- `logs/module-examples-ch06-vet.run.json`
- `logs/module-examples-ch06-vet.stderr.txt`
- `logs/module-examples-ch06-vet.stdout.txt`
- `logs/module-examples-ch07-race.run.json`
- `logs/module-examples-ch07-race.stderr.txt`
- `logs/module-examples-ch07-race.stdout.txt`
- `logs/module-examples-ch07-vet.run.json`
- `logs/module-examples-ch07-vet.stderr.txt`
- `logs/module-examples-ch07-vet.stdout.txt`
- `logs/module-gui-vet.run.json`
- `logs/module-gui-vet.stderr.txt`
- `logs/module-gui-vet.stdout.txt`
- `logs/module-root-race.run.json`
- `logs/module-root-race.stderr.txt`
- `logs/module-root-race.stdout.txt`
- `logs/module-root-vet.run.json`
- `logs/module-root-vet.stderr.txt`
- `logs/module-root-vet.stdout.txt`
- `logs/page-contract-checks-dom-stub-fixed.run.json`
- `logs/page-contract-checks-dom-stub-fixed.stderr.txt`
- `logs/page-contract-checks-dom-stub-fixed.stdout.txt`
- `logs/page-contract-checks.run.json`
- `logs/page-contract-checks.stderr.txt`
- `logs/page-contract-checks.stdout.txt`
- `logs/production-queue-tests.run.json`
- `logs/production-queue-tests.stderr.txt`
- `logs/production-queue-tests.stdout.txt`
- `logs/refusal-projection.run.json`
- `logs/refusal-projection.stderr.txt`
- `logs/refusal-projection.stdout.txt`
- `logs/selected-core-race.run.json`
- `logs/selected-core-race.stderr.txt`
- `logs/selected-core-race.stdout.txt`
- `module-gate-plan.json`
- `page-contract-checks-first-attempt.mjs`
- `page-contract-checks.mjs`
- `preserve.py`
- `private-generated-files.json`
- `private-source-verification.json`
- `protected-runtime-source-hashes.json`
- `read-input-pins.json`
- `read-ledger.json`
- `review-input.json`
- `reviewer_contract_test.go`
- `run-summary.json`
- `run.py`
- `task.txt`
