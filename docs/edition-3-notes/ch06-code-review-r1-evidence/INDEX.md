# Chapter 6 R1 independent re-review evidence index

Verdict and qualifications: `../ch06-code-review-r1.md`.
This is the continuation of the same independent reviewer and exact R1 task, not an entry/adapter restart. The final report and `SHA256SUMS` are the durable handoff. No paid calls or real key/settings reads were made.

## Authority and preservation

- `integrity.json`, `prerequisite-hashes.json`, `recovery-prerequisite-continuity.json`: exact authority pins, complete prerequisite continuity, all 310 R1 and 474 original artifact entries, previous review exact-set/index checks.
- `source-hashes.json`: 99 corrected source-file hashes. The 98 original extracted files were independently checked against their Git tree by the verifier; those completed checks are recorded in `integrity.json` and the verification receipts.
- `input-verification-r3.*`: successful initial verification. The two preceding failed attempts and verifier snapshots remain uncredited.
- `final-integrity-r2.*`: successful pre-interruption final check.
- `recovery-final-integrity.*`, `final-state.json`, `final-verifier.*`: post-recovery hash-only check; no behavioral replay.
- `main-status-final.txt`: final main state snapshot; concurrent publication edits are not candidate changes.
- `original-runtime-byproducts.json` / `source-runtime-byproducts.json`: separately retained local test output files, not additional candidate source.
- `core-build-info.txt`, `shape-build-info.txt`: exact frozen executables identify the corrected clean commit.

## Exact commands and fresh source/grade gates

`commands.json` inventories every recorded command. Each label has `LABEL.meta.json`, `LABEL.stdout`, `LABEL.stderr`, with exact argv/cwd, allowlisted environment, timestamps/duration, exit, disk observations and output hashes. No short timeout substituted for the actual 62-second phases.

- `base-sequence.py`, `base-planned.json`, `base-results.json`, `base-sequence.*`: sequential orchestration, all completed before interruption.
- `root-race.*`, `root-vet.*`: fresh full root race and vet.
- `ch04-consumer-{race,vet}.*`, `ch06-{race,vet}.*`, `ch07-{race,vet}.*`: separate module compile/race/vet gates; no-test-file output is not unit coverage.
- `grade-ch06.*`, `grade-ch01.*`: both 100/100, exit 0, exact same frozen R1 core. Chapter 1 includes genuine completed 62-second header and body phases.

## Independent causal and canonical proof

- `public_checks_test.go`, `batch_interrupt_test.go`: byte-identical qualified probes from the original whole review. Their source hashes are checked, not silently adjusted.
- `gemini_r1_independent_test.go`: additional public-only explicit canonical values and 13 boundary/continuation layouts, reuse/IDs/finals/deltas/usage, native replay and model/vendor/surface exclusion.
- `prepare-controls.py`, `public-module.go.mod.txt`: owned scratch preparation, no student/source-tree edit.
- `control-sequence.py`, `control-results.json`, `control-sequence.*`: corrected positive → compiled original implementation negative → byte-restored positive, same assertions.
- `original-restoration-control.diff`: the one production file temporarily restored to original bytes in scratch.
- `independent-corrected-{compile,all}.*`, `independent-vet.*`, `independent-original-{compile,boundaries}.*`, `independent-restored-{compile,all}.*`: full receipts. The original negative is exit 1; all compile/positive/vet controls exit 0. Restored and corrected test executables have the same SHA256.
- Binaries and their exact paths/hashes are in `control-results.json`; files remain in the owned scratch directory identified by `scratch-path.txt`.

## Exact-binary local integrations

- `shape_fixture.py`, `shape-fixture-input.py.txt`, `shape-fixture-path-only.diff`, `fixture-provenance.json`: inspected R1 local Gemini fixture, same assertions, owned output-root adaptation only.
- `actual-gemini-shape.*`, `shape-fixture-summary.json` and `shape-*` outputs: exact frozen R1 shape; streaming/plain both exit 0, same answer and three final parts, known-zero usage, early streaming delta.
- `cli-gated-control.py`, `cli-gated-control-input.py.txt`, `cli-pin-only.diff`: retained gated CLI assertions, exact R1 core pin and new owned output path only.
- `actual-cli-gated-and-stricter-assertion-r2.*`, `gated-cli.json`, associated `gated-cli.*` outputs: candidate exit 0; all public checks pass; the stricter `hint_in_answer` assertion remains false and checker exit 1. No model-compliance pass is inferred.
- `cli-pin-recovery.json`, `cli-gated-control-initial-pin.py.txt`, `cli-pin-only.initial.diff`, first CLI receipt: incorrect copied hash stopped before launch; preserved uncredited.
- `integration-sequence.*`, `integration-results.json`: the initial sequence failure remains explicit. The aggregate was reconciled after recovery using completed shape and qualified CLI receipts, not a rerun.

## Recovery and audit bookkeeping

- `recovery-state.json`: interruption boundary, same reviewer/context, 28 completed pre-interruption receipts, absent report/index at resume, no active relevant process, no replay or recovery-backup import.
- `SETUP-AND-FAILURES.md`: complete credit boundaries for student and reviewer mistakes, meaningful original-code negatives, stricter assertion failure and recovery.
- `verify*.py*`, `finalize*.py*`, `write-report*.py*`: retained working/pre-fix reviewer verification/report scripts. Failed setup variants are not candidate implementations.
- `resources.json`: minimum recorded disk space, sequential execution, no cleanup, unknown restart-era disk increase.
- `freeze.json`, `SHA256SUMS`: final report hash, evidence exact-set counts and per-file SHA256 manifest. `SHA256SUMS` excludes only its own recursive self-hash.

All candidate archives/executables and private existing artifacts remain at their pinned original paths and are not replaced by text exports. This directory contains safe local fixture material and artifact metadata, not raw private provider payloads or credentials. Owned scratch runtime output is retained separately from pinned source. The entire new evidence set is frozen read-only at handoff; copy it into a new owned directory before any future authorized replay. All jobs are complete.
