# Chapter 6 whole-candidate review evidence

Report: `../ch06-code-review.md`. Verdict: **required changes**, one cohesive Gemini signed/opaque boundary defect. Read R1 before treating the 100-point grader as whole-contract acceptance.

## Start here

- `final-integrity.json`, `final-integrity-check-r2.{meta.json,stdout,stderr}` — final authority/source/artifact verification, protected audits, unchanged production, restored controls.
- `commands.json` — every recorded command with cwd, exact argv, exit, duration, environment policy, output hashes, disk observations. Each label has complete `.stdout` / `.stderr`; recorder stdout previews are not the full evidence.
- `input-verification.json`, `source-hashes.json`, `prerequisite-hashes.json`, `production-source-review.json`, `candidate-diff.numstat` — provenance and source-read coverage.
- `SETUP-AND-FAILURES.md` — all uncredited setup/fixture/compile failures and their qualified successors.
- `SHA256SUMS`, `freeze.json` — exact frozen evidence set and report hash.

## Qualified finding and positive gates

- `independent-public-final-r2.*` + `public_checks_test.go` + `batch_interrupt_test.go`: four Gemini failures, seven other top-level tests pass; public module vet-r2 passes.
- `root-race.*`, `root-vet.*`, `ch04-consumer-{race,vet}.*`, `ch06-{race,vet}.*`, `ch07-{race,vet}.*`.
- `grade-ch06.*`: 100/100, pass, exit 0 on exact frozen core.
- `grade-ch01.*`: 100/100, pass, exit 0 on the same core, both real 62-second controls.
- `base-sequence.py`, `base-sequence-results.json`, `base-sequence.*` — sequential orchestration.

## Actual executable and inherited actor behavior

- `shape_fixture_original.py`, `shape_fixture.py`, `public-shape-reproduction.*`, `shape-fixture-summary.json`, `shape-fixture-{stream,plain}.*` — exact frozen external shape binary, pure loopback fixture, no live provider.
- `cli-gated-control.py`, `actual-cli-gated-hint-and-stricter-assertion-r2.*`, `gated-cli.json`, `gated-cli.{stdout,stderr,api.log}` — gated real CLI timing/admission, next-request-only hint, trace sanitization; intentional stricter answer assertion failure remains.
- `prepare-wakeup.py`, `review_stream_fixture_test.go`, `wakeup-selection.json`, `streaming-ch05-wakeup-and-interruption-r2.*` — ten selected original assertions under SSE, no production changes. Initial wrong-model wrapper/receipts are retained separately.

## Sensitivity and restoration

- `mutation-controls.py`, `mutation-results.json`, `ignore-disable.diff`, `wrong-final-id.diff`.
- `sensitivity-baseline*`, `ignore-disable-*`, `wrong-final-id-*`: changed-byte compilation, targeted failure, full-tree restoration, restored compilation/pass.
- `public_checks_mutation_control_source.go.txt` is the exact test source used for these compiled controls, before the later added multi-call test.

## Owned scratch and resource accounting

- `scratch-path.txt`; `owned-scratch-manifest.json` indexes every retained owned file, including source, review modules, mutated/restored test binaries, fixture outputs, and runtime byproducts.
- `source-runtime-extras.json`, `stream-wakeup-runtime-extras.json`, `mutation-source-*-runtime-extras.json` distinguish runtime-created files from the 98 pinned source files. They are not candidate changes and were not cleaned away.
- `resource-summary.json`: sequential gates, bounded Go concurrency, minimum recorded free space, no timeout wrapper/cache cleanup.
- `main-status-initial.txt` / `main-status-final.txt` where present record observed unrelated main changes; source/authority pins, not unrelated publication work, define this candidate.

No private provider body, settings file, or credential is copied here. API traces in this evidence set are reviewer-created loopback fixtures with redacted synthetic key sentinels. Existing private artifacts were hash-verified only. No jobs or children remain pending.
