# One-field pause-fixture repair — frozen independent evidence

**Verdict: PASS for the fixture repair only. NOT whole-student acceptance.** See `../ch07-pause-fixture-review.md`.

## Evidence map

- `task.txt`, `instruction-pins.json`, `read-ledger.json`: authorization, complete fresh reads/recovery versus prior complete byte-identical historical reads.
- `production.diff`, `root-candidate.json`, `root-source-change.json`, `root-binary.json`: exact one-field production repair and frozen pins. Other `root-*` files are explicitly CARRIED root race/vet/build receipts, not fresh independent tests.
- `full-route-argv.json`, `run.py`, `*.run.json`, matching stdout/stderr: actual independent clean-env commands, CWD, timing, exit, output hashes. No long-gate/race repetition or watchdog.
- `measured-results.json`: explicit PAUSE_REPORT scores/Passed extracted from framed test output. Setup0/20 is NOT detection; builds/tests have no invented grader score.
- `build-observer.py`, `go-invocations.jsonl`, `actual-full-gui.json`, `go-build-metadata.txt`, `binary-identity-inspection.txt`, `local-artifacts.json`: captured real full-grader GUI, compiled mutant/restored provenance and distinct build IDs/hashes. Exact source restoration is NOT claimed as binary equality.
- `source-input-hashes.json`, `mutation-source-pins.json`, `restored-source-hashes.json`, `pause-bypass.diff`:117file original, one actual private admission-gate mutation, exact117file restoration. No private application source payload packaged.
- `harness-input-hashes.json`, `capture-only.diff`, `pause_review_test.go`: copied repaired fixture/evaluator; capture/workspace-retention/wrapper only. No schedule/assertion/score changes.
- `pause-*-argv.json`, `observed-facts.json`: same pause protocol before/mutated/restored; one/three/one paused dispatches, false/true/false actual second-effect file, successful results and other-tab unpause. Server duplicate peer deliveries are not extra jobs.
- `preservation.json`, `package-check.py`: protected source, binaries,31file package, prior76/71/47bundles, clean r4 and twelve publication files; all checks pass.
- `incidents.txt`, `LOCAL_ONLY.md`, `main-start-status.txt`: all uncredited setup/inspection/packaging failures, ownership and raw-local-only limitations.

## Result boundary

Corrected full route: exit0,100/100,Passed=true. Real compiled gate bypass: exit1,0/20,Passed=false, actual GUI_SECOND effect while paused. Exact restored source/application: exit0,20/20,Passed=true. Historical GUI unit check passes. No inferred full score for pause-only controls. No native/browser/paid/provider campaign or whole-candidate approval. No pending jobs/tests/questions/children.

## Checksum rule

From repository root: `shasum -a 256 -c docs/edition-3-notes/ch07-pause-fixture-review-evidence/SHA256SUMS`.

Manifest covers external report plus every evidence file except itself, including this INDEX. Final submission supplies manifest's nonrecursive self-hash. Exact directory set below includes manifest. External report: `docs/edition-3-notes/ch07-pause-fixture-review.md`.

## Exact directory set

```text
INDEX.md
LOCAL_ONLY.md
SHA256SUMS
actual-full-gui.json
binary-identity-inspection.txt
build-observer.py
build-pause-bypass.run.json
build-pause-bypass.stderr.txt
build-pause-bypass.stdout.txt
build-pause-wrapper.run.json
build-pause-wrapper.stderr.txt
build-pause-wrapper.stdout.txt
build-restored.run.json
build-restored.stderr.txt
build-restored.stdout.txt
capture-only.diff
corrected-full-r4.run.json
corrected-full-r4.stderr.txt
corrected-full-r4.stdout.txt
full-route-argv.json
go-build-metadata.txt
go-invocations.jsonl
harness-input-hashes.json
historical-gui.run.json
historical-gui.stderr.txt
historical-gui.stdout.txt
incidents.txt
instruction-pins.json
local-artifacts.json
main-start-status.txt
measured-results.json
mutation-source-pins.json
observed-facts.json
package-check.py
pause-bypass-argv.json
pause-bypass.diff
pause-bypass.run.json
pause-bypass.stderr.txt
pause-bypass.stdout.txt
pause-positive-argv.json
pause-positive-retry-argv.json
pause-positive-retry.run.json
pause-positive-retry.stderr.txt
pause-positive-retry.stdout.txt
pause-positive.run.json
pause-positive.stderr.txt
pause-positive.stdout.txt
pause-restored-argv.json
pause-restored.run.json
pause-restored.stderr.txt
pause-restored.stdout.txt
pause_review_test.go
preservation.json
production.diff
read-ledger.json
restored-source-hashes.json
root-binary.json
root-build.stderr.txt
root-build.stdout.txt
root-candidate.json
root-commands.json
root-fixture.diff
root-historical-gui-race.stderr.txt
root-historical-gui-race.stdout.txt
root-source-change.json
root-vet.stderr.txt
root-vet.stdout.txt
run.py
source-input-hashes.json
task.txt
```
