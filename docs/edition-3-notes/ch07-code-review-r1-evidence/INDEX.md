# Chapter 7 formal round-1 independent re-review evidence

Verdict: PASS for corrected code/whole-adjacent re-review. R1/R2 closed; C1 remains non-blocking. Root alone accepts/tags. Student source: `463996b765631d31923a5726036a34942a67cefd`, tree `5e95dcdd716f7c25c96cc3ebcde13d000bb38670`; formal correction usage 1/3.

## Reading order

1. `../ch07-code-review-r1.md` — bounded verdict, exact candidate, measured controls/gates and limitations.
2. `input.json`, `initial-verification.json`, `source-manifest.json`, `evidence-manifest.json`, `read-ledger.json` — frozen inputs and honest retained/fresh read scope.
3. `preservation.json`, `control-source-verification.json` — unchanged source, runtime/publication, prior bundles and corrected/original comparison.
4. `wire-facts.json`, final page `*.facts.json`, corresponding `.run.json`/`.stdout.txt`/`.stderr.txt` — genuine behavior outcomes.
5. `reviewer_wire_r1_test.go`, `reviewer_page_r1.mjs` — unchanged final reviewer assertions on both source versions. Earlier `reviewer_*.txt` snapshots and `INCIDENTS.md` retain setup/assertion mistakes; they are not extra detections.
6. `gate-plan.json`, `gates-summary.json`, `chapter1-lifetime-checks.json`, the twelve gate receipts — sequential fresh race/vet, Chapter7 and real Chapter1 results. Score and passed are distinct from process exit.
7. `QUALIFICATIONS.md` and `local-only-artifact-pins.json` — no fresh native/paid/browser claim; local-only binaries/raw traces are not packaged.

## Scope and safety

All files here are UTF-8 text without NUL bytes. No executable, archive, pyc, real credential, raw wire/provider trace, or generated private cr/api payload is copied. Constant local-fixture authentication and WebSocket nonce literals in reviewer source are not real keys. Hash-only inventories preserve local provenance without supplying payloads.

The original-r4 controls compiled and failed behaviorally; corrected controls passed and passed again on restoration. Every setup/API/parser/DOM/packaging failure is separately retained. Returning to the unchanged corrected path is the restoration; no candidate production member was mutated. All jobs finished; no children or pending experiment.

`SHA256SUMS` covers the report plus every listed evidence file except itself. The manifest's own hash is supplied in the final submission. Exact evidence set: **128 files**, below. No additional files are part of this seal.

## Exact file set

- `INCIDENTS.md`
- `INDEX.md`
- `QUALIFICATIONS.md`
- `SHA256SUMS`
- `all-run-summary.json`
- `chapter1-fresh.run.json`
- `chapter1-fresh.stderr.txt`
- `chapter1-fresh.stdout.txt`
- `chapter1-lifetime-checks.json`
- `chapter1-lifetime-facts.json`
- `chapter7-fresh.run.json`
- `chapter7-fresh.stderr.txt`
- `chapter7-fresh.stdout.txt`
- `compile-wire-corrected-api-fixed.run.json`
- `compile-wire-corrected-api-fixed.stderr.txt`
- `compile-wire-corrected-api-fixed.stdout.txt`
- `compile-wire-corrected-final.run.json`
- `compile-wire-corrected-final.stderr.txt`
- `compile-wire-corrected-final.stdout.txt`
- `compile-wire-corrected-history-fixed.run.json`
- `compile-wire-corrected-history-fixed.stderr.txt`
- `compile-wire-corrected-history-fixed.stdout.txt`
- `compile-wire-corrected-package-fixed.run.json`
- `compile-wire-corrected-package-fixed.stderr.txt`
- `compile-wire-corrected-package-fixed.stdout.txt`
- `compile-wire-corrected-parser-fixed.run.json`
- `compile-wire-corrected-parser-fixed.stderr.txt`
- `compile-wire-corrected-parser-fixed.stdout.txt`
- `compile-wire-corrected.run.json`
- `compile-wire-corrected.stderr.txt`
- `compile-wire-corrected.stdout.txt`
- `compile-wire-original-r4.run.json`
- `compile-wire-original-r4.stderr.txt`
- `compile-wire-original-r4.stdout.txt`
- `consumer-race.run.json`
- `consumer-race.stderr.txt`
- `consumer-race.stdout.txt`
- `consumer-vet.run.json`
- `consumer-vet.stderr.txt`
- `consumer-vet.stdout.txt`
- `control-source-verification.json`
- `evidence-manifest.json`
- `gate-plan.json`
- `gates-summary.json`
- `gui-race.run.json`
- `gui-race.stderr.txt`
- `gui-race.stdout.txt`
- `gui-vet.run.json`
- `gui-vet.stderr.txt`
- `gui-vet.stdout.txt`
- `initial-verification.json`
- `input.json`
- `local-only-artifact-pins.json`
- `original-r4-source-pins.json`
- `page-corrected-assets-fixed.run.json`
- `page-corrected-assets-fixed.stderr.txt`
- `page-corrected-assets-fixed.stdout.txt`
- `page-corrected-dom-fixed.run.json`
- `page-corrected-dom-fixed.stderr.txt`
- `page-corrected-dom-fixed.stdout.txt`
- `page-corrected-restored.facts.json`
- `page-corrected-restored.run.json`
- `page-corrected-restored.stderr.txt`
- `page-corrected-restored.stdout.txt`
- `page-corrected-wire-payload-fixed.facts.json`
- `page-corrected-wire-payload-fixed.run.json`
- `page-corrected-wire-payload-fixed.stderr.txt`
- `page-corrected-wire-payload-fixed.stdout.txt`
- `page-corrected.run.json`
- `page-corrected.stderr.txt`
- `page-corrected.stdout.txt`
- `page-original-r4.facts.json`
- `page-original-r4.run.json`
- `page-original-r4.stderr.txt`
- `page-original-r4.stdout.txt`
- `prerequisite-comparison.json`
- `preservation.json`
- `preserve.py`
- `private-extra-file-inventory.json`
- `read-ledger.json`
- `reviewer_page_r1.dom-caption-expectation.txt`
- `reviewer_page_r1.inert-text-error.txt`
- `reviewer_page_r1.mjs`
- `reviewer_wire_r1_test.first-compile.txt`
- `reviewer_wire_r1_test.frame-parser-error.txt`
- `reviewer_wire_r1_test.go`
- `reviewer_wire_r1_test.history-api-error.txt`
- `reviewer_wire_r1_test.missing-system.txt`
- `reviewer_wire_r1_test.public-api-errors.txt`
- `root-race.run.json`
- `root-race.stderr.txt`
- `root-race.stdout.txt`
- `root-vet.run.json`
- `root-vet.stderr.txt`
- `root-vet.stdout.txt`
- `run-plan.py`
- `run.py`
- `runtime-before.json`
- `source-manifest.json`
- `stream-race.run.json`
- `stream-race.stderr.txt`
- `stream-race.stdout.txt`
- `stream-vet.run.json`
- `stream-vet.stderr.txt`
- `stream-vet.stdout.txt`
- `task.txt`
- `wire-corrected-parser-fixed.run.json`
- `wire-corrected-parser-fixed.stderr.txt`
- `wire-corrected-parser-fixed.stdout.txt`
- `wire-corrected-restored.run.json`
- `wire-corrected-restored.stderr.txt`
- `wire-corrected-restored.stdout.txt`
- `wire-corrected-valid-fixture.run.json`
- `wire-corrected-valid-fixture.stderr.txt`
- `wire-corrected-valid-fixture.stdout.txt`
- `wire-corrected.run.json`
- `wire-corrected.stderr.txt`
- `wire-corrected.stdout.txt`
- `wire-facts.json`
- `wire-original-r4.run.json`
- `wire-original-r4.stderr.txt`
- `wire-original-r4.stdout.txt`
- `workflow-race.run.json`
- `workflow-race.stderr.txt`
- `workflow-race.stdout.txt`
- `workflow-vet.run.json`
- `workflow-vet.stderr.txt`
- `workflow-vet.stdout.txt`
