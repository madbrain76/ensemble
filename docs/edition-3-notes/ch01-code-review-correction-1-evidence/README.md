# Chapter 1 correction round 1 — independent re-review evidence

Scope: close initial R1 and check regressions from the two-line test correction. No new hardening or author phase.

- Revised candidate: `1aad3975842d502b9144b7b9bddfd30ca7309db7`.
- Initially reviewed candidate: `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`.
- Unchanged book/grader: `7f3d069420112c2136bc95d1d551e12fad6964c1`.

## Reproducible reviewer checks

`verify.py` reads both exact Git revisions, checks the entire path set and file bytes, exports the revised candidate to a temporary archive, builds/vets/tests it, and exercises one first-block mutation in a separate disposable copy. It does not modify either repository's source. `verification.txt` and `verification.exit` record the runner outcome; per-command logs/exits retain actual results.

From the exact revised archive:

```sh
go version
go build -o "$SCRATCH/ch01" ./cmd/ch01
go vet ./...
go test -count=1 -v -run '^TestConversation$' .
go test -race -count=1 -v ./...
```

The focused command is also run in a disposable mutant copy after inserting `break` immediately after `text.WriteString(block.Text)`, then again after byte-for-byte restoration. `gofmt -l` is applied to all archived Go files without rewriting them. The full race suite includes the real sequential 62-second response-header and response-body waits; no short shell timeout terminates them.

## Artifact index

- `correction.patch`: complete exact-commit diff; only `ensemble_test.go:54–55` changed.
- `provenance.json`: both revisions, file hashes/identity, unchanged assertions, chapter/grader hashes, protected original/student artifact hashes, exact archive location, and final verification status.
- `scratch-path.txt`: temporary exact-commit archive used for this run.
- `go-version.txt`/`.exit`, `build.txt`/`.exit`, `vet.txt`/`.exit`: actual reviewer commands/results.
- `focused-control.txt`/`.exit`: revised fixture's passing control.
- `first-block-mutant.txt`/`.exit`: expected behavioral failure, exit 1, at `ensemble_test.go:71` with `all text blocks: "hello " <nil>`. This failure demonstrates retained sensitivity, not a candidate failure.
- `focused-restored.txt`/`.exit`: successful focused rerun after the disposable Engine file is restored byte-for-byte.
- `tests.txt`/`.exit`: independent full revised-candidate race suite.
- `format.txt`/`.exit`: command/cwd metadata followed by any gofmt filenames; no filenames means clean formatting.
- `verification.txt`/`.exit`: overall verification runner result; exit 0 requires all controls, expected mutant failure, suite and preservation checks to pass.
- `hygiene.txt`: final archive/repository/preserved-artifact checks.
- `student-evidence-validation.txt`: independently recomputed all 21 original-artifact hashes and both revision hashes for all nine tracked files recorded in the student's correction evidence; inspected student mutation, restoration, full tests and build metadata separately.

## Evidence attribution

These are new correction-round-1 checks. Earlier reviewer grader/live/mutation/source-review artifacts remain unchanged under `../ch01-code-review-evidence/`, associated with `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`. They are not renamed or presented as new `1aad397` runs. Production source and go.mod are byte-identical across those revisions; that verified identity is the reason the earlier production evidence remains applicable. VCS metadata or build paths can change binary bytes, so no old/new binary-identity claim is made.

No grader rerun, provider calls or credential access occurred in this correction round. The earlier successful executable-target grader and separate failed directory-target invocation keep their original attribution. No student source, manuscript, grader, historical reference, initial review/evidence or student artifacts were edited. Only the new correction review/evidence paths are written. The parent-reported dispatch hold is not treated as a test/code failure.
