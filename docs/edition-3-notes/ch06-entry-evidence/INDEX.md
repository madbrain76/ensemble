# Chapter 6 entry evidence index

Private reviewer evidence; not student answer material. Main report: `../ch06-entry-audit.md`.

## Read first

- `results.json`: every recorded command exit/duration and all scored reports' score/pass/failed-check IDs.
- `commands.txt`: exact argv, working directory, controlled environment policy/overrides, timestamps, exits, elapsed seconds, disk measurements, and stdout/stderr hashes.
- `*.meta.json`, `*.stdout`, `*.stderr`: authoritative unabridged per-command evidence. Never substitute a truncated tool tail for these files.
- `inputs.json`: hashes/byte/line counts of full prerequisites and inspected grading inputs.
- `git-identities.txt`, `initial-status.txt`, `final-verification.json`: immutable pins, clean tracked main/student, unchanged source/manuscript/grader hashes, protected FINAL tree, and intentionally absent Chapter 6 exercise.
- `historical-source-sha256.json`, `student-source-sha256.json`: source byte manifests.
- `binary-sha256.json`, `final-binary-sha256.json`, `buildinfo.txt`: initial/final binary identities and Go build metadata.
- `SHA256SUMS`: final checksum manifest of every evidence file except itself.
- `REPORT-SHA256`: SHA-256 binding the sibling audit report.

## Grading runs

- `positive-api.*`: pinned historical positive via explicit-identity existing API, 100/100, pass, exit 0.
- `positive-old-cli-source-valid.*`: corrected old CLI source-root target, 100/100, pass, exit 0.
- `positive-old-cli-binary-valid.*`: identical frozen core binary but conflated old CLI source identity, 70/100, fail, exit 1.
- `ch05-nonsubmission-baseline.*`: accepted predecessor and genuinely absent Ch6 module, 20/100, fail, exit 1; inherited parity passes. NOT a Ch6 student submission.
- `mutant-disable.*`: 90/100, fail, exit 1; real disable/request flag defect.
- `mutant-correlation.*`: 50/100, fail, exit 1; real delta/final ID defect.
- `restored-disable.*`, `restored-correlation.*`: restored source/rebuilt binary controls, both 100/100, pass, exit 0.
- `positive-old-cli-source.*`, `positive-old-cli-binary.*` without `-valid`: wrong-flag invocation attempts, exit 2; no report/score. Kept rather than overwritten.

All builds have `build-*.{meta.json,stdout,stderr}` and exited 0. All grader invocations were sequential. There is no wall-clock kill timeout in the runner.

## Reproduction/provenance artifacts

- `run.py`: sequential command recorder; environment allowlist excludes ambient credentials; only local fake keys supplied by existing grading code.
- `api/main.go`: formatted evidence-only caller of `internal/gradeedition2.Chapter6`; not a grader edit.
- `api-source-initial.txt`: exact original source used to compile the retained API binary. `verify-harness-format.*` confirms the formatted saved source is gofmt-identical.
- `archive.*`: exact pinned historical Git export command and exit.
- `provenance.py`, `provenance.*`: initial identities/manifests and source pins.
- `mutations.py`, `mutations.json`, `disable.diff`, `correlation.diff`, `mutant-*-source-sha256.json`: two exact scratch defects and changed source bytes.
- `restore.py`, `restore-source.*`, `restoration.json`: before/after tree equality and retained mutant binary hashes.
- `finalize.py`: final read-only identity/source verification and summary generator (no behavioral tests).
- `student-consumer-case.json`, `student-media-models.json`: accepted baseline configuration copied without changing the student.
- `scratch-path.txt`: `/tmp/ensemble-ch06-entry-7_7zbtqs`, containing retained original/mutant/restored core binaries, harness, old CLI, and restored source copies.
- `invocation-notes.md`: retained tooling mistakes, explicitly separated from grading outcomes.
- `freeze-status.json`: final verification status, evidence count and minimum recorded disk availability.

To reproduce a recorded command in a fresh owned scratch directory, use its exact `.meta.json` argv/cwd/environment policy, remapping only the recorded scratch prefix consistently. Do not rerun scripts into this frozen evidence directory or mutate an accepted student/solution tree. Internal grader-built short-lived consumer/workflow/exercise executable hashes were not intercepted; their source identities/build commands are recorded, not falsely represented as captured binary hashes.

This entry audit did not implement a root adapter, authorize a student, or discharge the remaining public source/integration obligations. No pending audit jobs or children remain.
