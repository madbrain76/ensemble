# Chapter 7 entry evidence — frozen portable index

This is reviewer-only entry-audit evidence, not a student answer bundle or candidate acceptance.

- Report: `../ch07-entry-audit.md`.
- `*.run.json`: exact argv/cwd/environment controls/timing/actual subprocess exit. Matching stdout/stderr files retain successes and failures, including initial wrong CLI syntax.
- `runs-inventory.json`: exact captured run metadata in one file.
- `read-inventory.json`, `read-ledger.txt`: complete prerequisite/contract reads and truncation recovery.
- `entry-api-main.go.txt`, `existing-gui-checks-test.go.txt`, `run.py`: private invocation harness using unchanged existing checks.
- `control.py`, `replay.diff`, `pause.diff`, `*.hashes.json`, `restoration.json`: two compiled behavior controls and exact restorations. These are not student implementation guidance.
- `consumer-case.json`, `media-models.json`: explicit public-form input maps used by the historical positive.
- `historical-positive-initial-hashes.json`, `harness-initial-hashes.json`: source input SHA-256 inventories. No full historical solution payload included.
- `initial-pins.json`, `final-invariants.json`: initial and pre-handoff source/student/publication snapshot. Later root-authorized typography publication changes are outside these temporal claims.
- `local-only-artifacts.json`: private scratch binary locations/hashes; binaries, archives, raw traffic traces, and historical fixture source remain local-only.
- `audit-incidents.txt`: setup/inspection/packaging failures and qualifications, not behavior detections.

All portable files are UTF-8 text. No keys, binary/archive payloads, raw provider traffic, or raw GUI traffic logs are included. Existing grader reports include only safe aggregate assertions. Local scratch remains `/tmp/cadence-ch07-entry-auditor`; it is not a student input.

## Integrity and exact set

Run from repository root: `sha256sum -c docs/edition-3-notes/ch07-entry-evidence/SHA256SUMS` (or `shasum -a 256 -c ...`). Manifest covers this INDEX, every other evidence file except itself, and the report. Manifest's own SHA is supplied in the final handoff; there is no recursive self-hash. This listed set includes the manifest itself.

Report outside directory: `docs/edition-3-notes/ch07-entry-audit.md`.

```text
INDEX.md
SHA256SUMS
api-positive.run.json
api-positive.stderr.txt
api-positive.stdout.txt
audit-incidents.txt
build-api.run.json
build-api.stderr.txt
build-api.stdout.txt
build-cli2.run.json
build-cli2.stderr.txt
build-cli2.stdout.txt
build-cli3.run.json
build-cli3.stderr.txt
build-cli3.stdout.txt
build-gui-checks.run.json
build-gui-checks.stderr.txt
build-gui-checks.stdout.txt
build-positive.run.json
build-positive.stderr.txt
build-positive.stdout.txt
cli2-no-core-corrected.run.json
cli2-no-core-corrected.stderr.txt
cli2-no-core-corrected.stdout.txt
cli2-no-core.run.json
cli2-no-core.stderr.txt
cli2-no-core.stdout.txt
cli2-positive-corrected.run.json
cli2-positive-corrected.stderr.txt
cli2-positive-corrected.stdout.txt
cli2-positive.run.json
cli2-positive.stderr.txt
cli2-positive.stdout.txt
cli3-unsupported-corrected.run.json
cli3-unsupported-corrected.stderr.txt
cli3-unsupported-corrected.stdout.txt
cli3-unsupported.run.json
cli3-unsupported.stderr.txt
cli3-unsupported.stdout.txt
consumer-case.json
control.py
entry-api-main.go.txt
existing-gui-checks-test.go.txt
final-invariants.json
harness-initial-hashes.json
historical-positive-initial-hashes.json
initial-pins.json
local-only-artifacts.json
media-models.json
pause-mutant-build.run.json
pause-mutant-build.stderr.txt
pause-mutant-build.stdout.txt
pause-mutant.run.json
pause-mutant.stderr.txt
pause-mutant.stdout.txt
pause-restored.run.json
pause-restored.stderr.txt
pause-restored.stdout.txt
pause.diff
pause.hashes.json
read-inventory.json
read-ledger.txt
replay-mutant-build.run.json
replay-mutant-build.stderr.txt
replay-mutant-build.stdout.txt
replay-mutant.run.json
replay-mutant.stderr.txt
replay-mutant.stdout.txt
replay-restored.run.json
replay-restored.stderr.txt
replay-restored.stdout.txt
replay.diff
replay.hashes.json
restoration.json
run.py
runs-inventory.json
```
