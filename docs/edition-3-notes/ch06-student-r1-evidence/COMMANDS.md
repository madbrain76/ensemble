# R1 command and evidence guide

Source cwd: `/Users/bill/projects/ensemble-edition-3-student`.
The original candidate800da21 and all original artifacts are retained unchanged.
Only the Gemini accumulator and a new public test were committed for R1.

## Reproduced old-code failure and correction

`go test -run '^$' ./...` first exposed a test-only field-name mistake; the
corrected compile passed. `go test -count=1 -timeout=60s -run
'^TestChapter6GeminiBoundary' .` then exposed a missing positive fixture token
limit. Both setup failures remain separately labelled in attempts.md.

`original-boundary-corrected-inputs.json` pins original production and the exact
new tests for the meaningful failing command above. The unchanged final test bytes
are reconciled in `control-reconciliation.json`. After the scoped correction,
that same command passed (`corrected-boundary-applied.*`). A refused edit followed
by a shell rerun left `corrected-boundary.*` as another ORIGINAL-code failure,
not a passing corrected run. No receipt was overwritten.

Focused command: `go test -count=1 -timeout=90s -run
'^TestChapter6|^TestLiveProtocol|^TestDisplay' ./...`.

## Freeze and final gates

Explicit stage paths: `internal/engine/stream.go`,
`chapter6_gemini_boundary_test.go`. Source commit and local-only artifact pins are
in the immutable R1 source manifest. `build-commands.json` records exact sequential
build cwd/argv/exits. `git archive HEAD` created the source archive exclusively in
the new correction temp directory. It is LOCAL-ONLY.

`run-final-gates.py` owns the one sequential final-byte gate sequence.
`final-planned-commands.json` records plans; `final-command-receipts.json` records
ONLY actual completed commands, cwd, argv, bounded timeout, elapsed duration,
exit/stdout/stderr hashes, source identity and exact core binary pin. The runner
preserves the current parent environment PATH only plus HOME and explicitly
chosen Go/temp settings; no provider key environment is forwarded. GOPROXY and
GOSUMDB are off, GOTOOLCHAIN local, GOFLAGS -mod=readonly. No install/download,
provider call or credentials read is part of this sequence.

Root test/race limits240/300seconds accommodate both genuine62-second phases.
The runner does not shorten any inner lifetime check. Module test/vet commands are
run for ch04-consumer, ch06 workflow and ch07 shape. Final focused tests retain the
actual gated human CLI check. `shape_fixture.py` runs the newly frozen shape
against local Gemini SSE and JSON, with signed/opaque/unsigned boundaries and a
150ms pause, restricted child environment, explicit stream-off comparison and
known-zero usage. Raw local runtime state stays beneath the new temp directory.

## Exact paired grade route

Both commands run from `/Users/bill/projects/ensemble`, using the same
`/tmp/e3-ch06-correction-r1/core` and all flags before that executable.
The actual structured argv/exits/hashes are in final-command-receipts.json.
The Chapter6 flags supply source core, consumer, workflow and stream-exercise
module directories, stream model `claude-sonnet-4-6`, and the new R1
`media-models.json`. That92-byte file is an unchanged copy of the previously
selected declared rows, not a new capability claim. Chapter1 uses the same core
binary. Grade scores never override boundary tests or independent review.

No original paid receipt is rerun, edited, reclassified as passing, or silently
carried forward to changed Gemini production. Final report/index will distinguish
completed gates from any pending/failed command. No acceptance tag is created.

## Completed final boundary

All13 planned gates ran to completion, exit0. Both paired grades returned
passed:true and100/100. `paired-grade-command-receipt.json` is derived from the
actual runner receipts, not a reconstructed plan. Both inner62-second Chapter1
phases completed naturally. `final-reconciliation.json` verifies clean source
and every frozen artifact pin. `original-boundary-preservation.json` verifies
all474 original indexed artifacts unchanged. No R1 paid call/key read took place.
Archives and binaries are LOCAL-ONLY in the R1 inventory; no UTF-8 exception.
