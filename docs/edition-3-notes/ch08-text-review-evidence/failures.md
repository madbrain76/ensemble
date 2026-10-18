# Review setup failures and recoveries

This is the same single formal review, not a retry commission. No chapter, implementation, grader, acceptance file, publication asset, or other existing document was edited. No child was created and no runtime exercise was performed.

## Oversized reads

The combined carryover 631–1110 + both full voice guides + full chapter procedure emitted 153,255 bytes against a 54,000-byte cap. It was not counted as a complete displayed read. All affected documents were recovered in smaller complete ranges. The carryover 631–949 recovery emitted 28,252 bytes plus its heading against a 28,000-byte cap; an overlap was initially recovered, then both 631–790 and 791–949 were explicitly read in full. A later combined carryover recovery plus AGENTS emitted 52,742 bytes against 42,000; AGENTS and the affected carryover sections were read separately afterward. `read-coverage.json` selects the successful full reads, not these truncated emissions. `reads.jsonl` records successful file-read operations and hashes; it must not by itself be mistaken for proof that oversized tool output was displayed completely.

## Nonexistent guessed receipt paths

These `read.py` calls failed with FileNotFoundError and did not read or mutate a file:

- `docs/edition-3-notes/ch08-acceptance-gates.md`
- `docs/edition-3-notes/ch08-author-evidence/initial-chapter.diff`
- `docs/edition-3-notes/ch08-author-evidence/final-chapter.diff`
- `docs/edition-3-notes/ch08-author-evidence/checks-output.txt`
- `docs/edition-3-notes/ch08-author-evidence/verification.json`
- `docs/edition-3-notes/ch08-author-evidence/finalization.json`
- `docs/edition-3-notes/ch08-author-evidence/commands.jsonl`

A schema-inspection Python block also failed opening nonexistent author `MANIFEST.json`; the actual seal is `inventory.json`. An attempted empty-log read opened existing baseline and final diff-check logs (both empty), then failed on nonexistent `final-report-diff-check.log`. The subsequent loop entry `initial-final-diff-check.log` was not reached. These were reviewer retrieval errors, not missing required author evidence. The actual inventory names were listed/read, and the required actual records were recovered. Some combined shell calls continued after individual failures; a shell's aggregate exit status is not used to erase them.

## Parent lookup / missing search utility

The media-map validator was outside the two initially permitted source files. Permission was requested, not presumed. An incorrect `wait_for_job` handle 72 returned `Unknown job handle: 72`; a repeated parent question was then issued. The successful wait on handle 2158 returned parent message `main-12`, explicitly authorizing only ReadMediaModels and directly necessary validation source. It also confirmed the earlier permission was sent as `main-11`. Permission is saved verbatim in `narrow-source-permission.json`.

`rg -n '^func ReadMediaModels|^func .*MediaModels' internal/gradeedition2` failed because `rg` is unavailable. Standard `grep -n '^func ReadMediaModels\|^func .*MediaModels' internal/gradeedition2/*.go` located the function without opening unrelated source contents. The subsequent source read was `internal/gradeedition2/ch05_media.go` lines 1–60; the relevant validator is lines 18–41. The bounded window incidentally exposed the opening of the next function (44–60). That opening was not assessed, no later fixture body was read, and no grader was executed. The whole-file identity was recorded without treating unshown lines as read.

## Whitespace command syntax

The saved `check.py` used `git diff --check -- book/edition-3/chapter-08.md`, omitting mandatory `--no-pager`. It returned zero but was a procedural mistake. The corrected exact command was run: `git --no-pager diff --check -- book/edition-3/chapter-08.md`; it returned zero with empty output, retained in `chapter-diff-check-corrected.log` and `commands.jsonl`. The original script/log remain unchanged as evidence, not a recommended rerun recipe.

## Interpretation of receipts

`read.py`, `coverage.py`, `check.py`, and `finalize.py` retain the executed helper source. The reads ledger records exact helper argv, file identity, line range and emitted-output hash. The command ledger records subprocess argv, return status and retained output identity. This file records retrieval/tool failures not captured by the successful-read helper. These records are not claimed to be a complete raw transcript of every tool wrapper or every user-visible commentary.

The author separately disclosed its workflow-tail read-order failure and earlier excluded-publication hashing. Those remain author-process deviations; this reviewer did not repeat those publication reads/hashes or certify them as compliant.

All required current input/evidence reads and the permitted static checks were complete before the first report/verdict. No outstanding text finding remains.

## Finalizer pre-seal assertion

The first finalizer invocation stopped before writing the inventory: its report check assumed `git --no-pager diff --no-index --check -- /dev/null REPORT` must return zero. It returned 1 with zero output bytes, the no-index files-differ status. The command receipt and empty output were retained. The finalizer was corrected to accept 0/1 only with empty diagnostics, and also explicitly reject trailing spaces/tabs in the report. No sealed member existed at this point; no inventory was rewritten. The corrected finalizer is the only retained script version, with this entry preserving the failed assertion and recovery.
