# Retained setup, inspection and packaging incidents

These are not behavior detections. They are recorded separately from the two successfully executed contract checks. Exact executed test commands, cwd/environment, timings, exits, stdout and stderr are in `logs/`. Nothing was fixed in production source or the grader.

## Reads and guessed paths

- Several initial inspection commands guessed absent GUI support/evidence filenames (including a `gui/test/pause.mjs` path and root/student-local review artifact locations). The actual manifest and directory inventory were then used. No absent file was counted as a student defect.
- The first private-copy setup guessed a source-manifest path inside the student tree. The authoritative manifest is main `docs/edition-3-notes/ch07-source-manifest.json`. It uses `source_files`, not `files`; failed setup assertions/KeyError were corrected before tests ran. All 117 source hashes then matched.
- A guessed `chapter1-r4.stdout.json` path was absent. The actual carried report is `chapter1-r4-report.json`, with its separate exact invocation receipt.
- Large read outputs were checked for cap omissions. The carryover 801–1000 span was recovered, finally as complete 801–900 and 901–1000 reads. The student-review overlap and final actor lines were also recovered. No omitted middle was silently treated as read.

## Executable check setup

- `page-contract-checks` exited 1 with `TypeError: doc.createDocumentFragment is not a function`. The inert DOM was incomplete. This is not a detection. The first harness source and all stdout/stderr/metadata are preserved.
- Added only inert `createDocumentFragment`/`createTextNode` methods and reran the same three cases as `page-contract-checks-dom-stub-fixed`. Production modules were unchanged. It completed and observed the required speaker defect, qualified page-hide behavior, and passing stale-cause guard.
- The first full-module gate plan guessed separate `examples/ch04`/`examples/ch05` directories. Validation raised `missing requested module examples/ch04` before executing any gate. The actual modules are root, `gui`, `examples/ch04-consumer`, `examples/ch06`, `examples/ch07`. The corrected finite sequence completed successfully.
- Selected core/race tests ran before the later mandatory full-root race gate; that overlapping coverage is explicit. The full root gate took 139.72 seconds and was not cancelled. No fresh Chapter 1 CLI long-gate campaign was run; its exact frozen-r4 receipt is carried.

## Preservation/packaging assumptions

- An inspection guessed `initial-identities.json`/`initial-verification.json` rather than the actual saved snapshot names. The scratch directory inventory resolved this; no identity conclusion relied on those failed reads.
- An initial preservation assertion manually mistyped the tree as `17f79bc4d384c2d856b2830b27880235ecb073724`. The recorded and actual tree both are `0f5f7fbe4273e91e1be3aa13d0715a99f671eaf9`; the script now uses the review-input value.
- The next assertion compared annotated tag object `edition-3-ch06-r1` directly with the predecessor commit. It was corrected to `edition-3-ch06-r1^{commit}`. The accepted predecessor is unchanged.
- Exact scratch-set assertion expected only 117 original files plus one reviewer test. Full Go tests also generated `api.log`, `cmd/ch01/api.log`, and `cr/io/*`. Their names/hashes were inventoried separately. All 117 original files remained byte-identical. No logs were removed or their payloads copied into portable evidence.
- The first generated-file inventory console output was capped; the complete durable `private-generated-files.json` is retained. Only hashes/metadata, not raw logs, are portable.
- Pre-seal cross-check caught draft transcription errors in four binary/source hashes, an incorrect helper name (`toolEvent`), and the carried Chapter 1 score/duration/metadata wording. These were replaced from the measured review-input manifest and exact receipt before sealing. The correct carried Chapter 1 result is 100/100, passed true, exit 0, 124.154179 seconds; the receipt pins the CLI and grader-source hashes. No raw evidence or score was altered.

## Classification, not a hidden correction requirement

An early milestone listed three candidate findings. Root's later source/public-wording check qualified page-hide as non-blocking; see `contract-classification.md`. The final verdict requires only truthful dispatch projection and readable tool-result speech. The raw page test still records its original stronger expectation, with this explicit qualification.

## Final state

All executed commands/tests finished. Final preservation verifies input pins, 117 frozen files, clean student, peeled predecessor, thirteen runtime sources, twelve publication paths, four exact prior bundles, 33 native-index files and 31 discrepancy-index files. No child, browser/native/live-runtime action, paid provider request, credential access or shared cleanup occurred.
