# Chapter 7 thin-adapter review — frozen evidence index

Verdict: PASS for the exact frozen CLI adapter only. Report: `../ch07-grader-review.md`.

## Meaning and boundaries

- `candidate.json`, `input-pins.json`, `instruction-pins.json`: frozen candidate and inspected input identity.
- `read-ledger.json`: full prerequisite/contract/diff reads, including truncation and endpoint recovery.
- `*.run.json`: exact command/cwd/timing/exit/environment controls, stdout/stderr hashes, and score/Passed only when a real report exists. Matching stdout/stderr retained, including initial `legacy-ch4` usage failure.
- `validation-matrix.json`: 128 frozen-binary invocations (64 cases, two observed targets), expected exits and no build/launch. `positional-guards.json`: three additional guarded invocations.
- `results.json`, `legacy-results.json`: real reports; passing unlisted GUI-model fallback is not mislabeled a negative detection.
- `source-review.json`, `main.diff`, `main-test.diff`: complete candidate diff and exact legacy build/dispatch/output/exit-suffix preservation proof.
- `sentinel.go`, `run.py`, `validate.py`, `measure.py`, `legacy.py`: independent harness sources. `measure.py` preserves the first legacy4 setup mistake; `legacy.py` is the corrected finite continuation.
- `alternate-gui-go-mod.diff`: only private fixture binding change. No application/model behavior was mutated.
- `local-artifacts.json`, `historical-positive-source-hashes.json`, `input-map-hashes.json`: local-only private input identities, not fixture payloads.
- `carried-ch1-receipt-pins.json`: prior expensive Chapter1 evidence identity, explicitly not fresh or a run on current f595.
- `preservation.json`: final candidate/source/student/12-publication-path checks plus unchanged prior entry bundle and66file historical fixture.
- `incidents.txt`, `readme-local-only.txt`: failure qualifications, privacy and ownership boundaries.

No binary/archive/pyc/raw private traffic/full historical solution/key payloads are portable. Everything in the set below is UTF-8 text. These are reviewer-only entry/adapter materials, not student inputs or student acceptance. The prior entry evidence set is preserved separately and unchanged.

## Checksum rule

From repository root: `shasum -a 256 -c docs/edition-3-notes/ch07-grader-review-evidence/SHA256SUMS`.

The manifest covers the report plus every evidence file except itself, including this INDEX. Its nonrecursive self-hash is supplied in the final handoff. The exact set below includes the manifest itself. Report outside this directory: `docs/edition-3-notes/ch07-grader-review.md`.

## Exact evidence-directory set

```text
INDEX.md
SHA256SUMS
alternate-gui-build.run.json
alternate-gui-build.stderr.txt
alternate-gui-build.stdout.txt
alternate-gui-go-mod.diff
candidate.json
carried-ch1-receipt-pins.json
custom-case-media.run.json
custom-case-media.stderr.txt
custom-case-media.stdout.txt
historical-positive-source-hashes.json
incidents.txt
input-map-hashes.json
input-pins.json
instruction-pins.json
legacy-ch4-corrected.run.json
legacy-ch4-corrected.stderr.txt
legacy-ch4-corrected.stdout.txt
legacy-ch4.run.json
legacy-ch4.stderr.txt
legacy-ch4.stdout.txt
legacy-ch5-corrected.run.json
legacy-ch5-corrected.stderr.txt
legacy-ch5-corrected.stdout.txt
legacy-ch6-corrected.run.json
legacy-ch6-corrected.stderr.txt
legacy-ch6-corrected.stdout.txt
legacy-results.json
legacy.py
local-artifacts.json
main-test.diff
main.diff
measure.py
nondefault-gui-model.run.json
nondefault-gui-model.stderr.txt
nondefault-gui-model.stdout.txt
positional-guards.json
positive.run.json
positive.stderr.txt
positive.stdout.txt
preservation.json
race.run.json
race.stderr.txt
race.stdout.txt
read-ledger.json
readme-local-only.txt
results.json
run.py
sentinel-build.run.json
sentinel-build.stderr.txt
sentinel-build.stdout.txt
sentinel-positive.run.json
sentinel-positive.stderr.txt
sentinel-positive.stdout.txt
sentinel.go
source-review.json
validate.py
validation-matrix.json
vet.run.json
vet.stderr.txt
vet.stdout.txt
whitespace.run.json
whitespace.stderr.txt
whitespace.stdout.txt
wrong-gui-root.run.json
wrong-gui-root.stderr.txt
wrong-gui-root.stdout.txt
wrong-stream-model.run.json
wrong-stream-model.stderr.txt
wrong-stream-model.stdout.txt
```
