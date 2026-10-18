# Frozen r4 pause diagnostic — exact evidence index

Verdict: harness/input-contract mismatch established. Original fullgrade80/100 false preserved; NOT whole-student acceptance. Report: `../ch07-pause-diagnostic.md`.

## Evidence map

- `task.txt`: full frozen task/addendum SHA5f4c4c9d26b0515f6be711130db62bfa33d350be94fc6ffac3ae57212a107a45.
- `initial-verification.json`, `source-verified.json`:31file package (165053payload bytes),117file source, archive/binary pins; no source payload.
- `read-ledger.json`: fresh full reads/overlap recovery versus prior complete historical reads verified byte-identical.
- `reproduction-argv.json`, `*.run.json`, paired stdout/stderr: exact commands/cwd/env/timing/exit/output. Permission setup failure retained separately from retry.
- `measured-results.json`: actual extracted PAUSE_REPORT scores/Passed; build nulls not passes, permission failure not detection.
- `go-invocations.jsonl`, `build-observer.py`, `gui-build-metadata.txt`, `local-artifacts.json`: actual ./cmd build; canonical-path capture miss; same-source rebuild and distinct binary provenance. Original removed GUI output hash not claimed.
- `capture-only.diff`, `pause_diagnostic_test.go`: private capture/workspace-retention/wrapper, unchanged evaluator.
- `public-argument-only.diff`: sole executed behavioral delta, delay→ai_callback_delay.
- `explanation-only.diff`, build-pause-explained receipts: compiled-but-UNRUN disproved hypothesis.
- `observed-facts.json`: selected safe events/results/jobs/pause ordering and marker facts. Full raw private data stays local; duplicated server rows are peer deliveries, not additional jobs.
- `student-controls-inspected.json`: carried fresh/warm/protocol evidence inspected, not rerun.
- `harness-input-hashes.json`, `preservation.json`:166original internal files unchanged, student clean/source117match,31file package/12publication/prior entry+adapter preservation.
- `run.py`, `incidents.txt`, `LOCAL_ONLY.md`: runner, all uncredited failures including corrected draft byte-total assertion, ownership/privacy limits.

All tests completed sequentially. No pending jobs, paid/live-provider/native/browser/liveGUI work, source/grader edit,62second/race campaign, children or shared/manual cleanup. No binaries, archives, raw private payloads or keys in bundle.

## Checksum rule

From repository root: `shasum -a 256 -c docs/edition-3-notes/ch07-pause-diagnostic-evidence/SHA256SUMS`.

Manifest covers report and every evidence file except itself, including INDEX. Final submission supplies nonrecursive manifest self-hash. Exact directory set includes manifest. External report: `docs/edition-3-notes/ch07-pause-diagnostic.md`.

## Exact directory set

```text
INDEX.md
LOCAL_ONLY.md
SHA256SUMS
build-observer.py
build-pause-capture.run.json
build-pause-capture.stderr.txt
build-pause-capture.stdout.txt
build-pause-explained.run.json
build-pause-explained.stderr.txt
build-pause-explained.stdout.txt
build-pause-public-argument.run.json
build-pause-public-argument.stderr.txt
build-pause-public-argument.stdout.txt
capture-only.diff
capture-original-retry.run.json
capture-original-retry.stderr.txt
capture-original-retry.stdout.txt
capture-original.run.json
capture-original.stderr.txt
capture-original.stdout.txt
capture-public-argument.run.json
capture-public-argument.stderr.txt
capture-public-argument.stdout.txt
explanation-only.diff
go-invocations.jsonl
gui-build-metadata.txt
harness-input-hashes.json
incidents.txt
initial-verification.json
local-artifacts.json
measured-results.json
observed-facts.json
pause_diagnostic_test.go
preservation.json
public-argument-only.diff
read-ledger.json
rebuild-identical-gui.run.json
rebuild-identical-gui.stderr.txt
rebuild-identical-gui.stdout.txt
reproduction-argv.json
run.py
source-verified.json
student-controls-inspected.json
task.txt
unchanged-r4.run.json
unchanged-r4.stderr.txt
unchanged-r4.stdout.txt
```
