# Chapter 2 correction-round-one reviewer evidence

Disposition and source citations: `../ch02-code-review-r1.md`. Candidate: `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`. Private reviewer evidence; original report/evidence remain untouched.

## Identity and source

- `identity.json`: complete pinned source/test hashes, archive/report/binary identities, byte-identical independent build.
- `binary-build-identity.stdout`: candidate VCS revision and `vcs.modified=false`.
- `governing-hashes.json`: unchanged governing sources, checked against the prior complete-read record; current reviewer skill reread.
- `manifest-verification.json`: all 126 original and 135 round-one submitted manifest entries, including byte counts and SHA-256 hashes.
- `prior-review-integrity.json`: verification of all 63 prior review checksum entries, including the corrected initial report.
- `final-integrity.json`: final frozen-candidate and evidence checks; `SHA256SUMS` includes this new review and every new evidence file except the checksum file itself. Paths are relative to the authoritative repository root.

## Source / grader batch

`run-statuses.json` records exact sanitized argv/cwd, exits, elapsed times, watchdogs, and grader target hashes. Each command has `.stdout`/`.stderr`; grader `.json` files preserve full results. No credential variables are inherited; environment allowlist is PATH/HOME/TMPDIR plus `GOPROXY=off GOSUMDB=off`.

`$STUDENT=/Users/bill/projects/ensemble-edition-3-student`, `$REPO=/Users/bill/projects/ensemble`, `$EVIDENCE=$REPO/docs/edition-3-notes/ch02-code-review-r1-evidence`, `$NEW_CLI=$REPO/docs/edition-3-notes/ch02-student-r1-evidence/student-cli`; `$SCRATCH` is disposable reviewer storage.

```sh
# cwd: $STUDENT
 go build -o "$SCRATCH/student-cli" ./cmd/ch01
 go test -race -count=1 -timeout=180s -coverpkg=./... \
   -coverprofile="$EVIDENCE/reviewer-cover.out" ./...
 go test -race -count=1 -timeout=30s -v \
   -run '^TestChapter2(GeminiChronologicalNames|MeasuredEmptyReplies|OpenAIMultipleOpaqueParts)$' .
 go vet ./...
# cwd: $REPO
 go version -m "$NEW_CLI"
 go build -o "$SCRATCH/grade-edition2" ./cmd/grade-edition2
 go build -o "$SCRATCH/grade-edition3" ./cmd/grade-edition3
 "$SCRATCH/grade-edition2" -ch 2 -json "$NEW_CLI"
 "$SCRATCH/grade-edition3" -ch 1 -json "$NEW_CLI"
```

Source tests and Edition 3 grader separately ran their actual 62-second header/body probes. Outer watchdogs were 210 seconds for the source suite and 180 seconds for the Edition 3 grader; neither expired. Both graders target the same new binary, hash-checked before and after. The independently built binary is identical, not a separately assumed equivalent target.

`reviewer-cover.out` and `coverage-summary.json`: union duplicate source-span/statement-count blocks; any hit covers a block. Result 749/838 production statements = 89.379%, root module including CLI; excludes the separate GUI module, test source, external executable runs and historical live evidence. No numeric threshold claimed.

## Previous reviewer probe

`public-test.go.txt` is byte-identical to the prior reviewer's test, which previously exposed R1–R3. In scratch, copy it to `public_test.go` in a separate module requiring `ensemble v0.0.0` with `replace ensemble => $STUDENT`, then run:

```sh
go test -race -count=1 -timeout=30s -v ./...
```

`public-regressions-status.json` and stdout/stderr preserve the passing actual run. Public root/real children only, no internal imports or production mock interfaces.

## Three scoped fix reverts

`check-sensitivity.py` is an unchanged copy of the student's inspected round-one helper, with its output directory selected by its copied location. It copies tracked frozen source to a disposable temporary directory, reverts one correction at a time, executes its exact regression test, restores original bytes, reruns the test, and verifies source restoration.

```sh
# cwd: $REPO
python3 "$EVIDENCE/check-sensitivity.py"
```

`R1-sensitivity.txt`, `R2-sensitivity.txt`, `R3-sensitivity.txt` retain actual behavioral failures and restore passes. `sensitivity-summary.json` records selectors, exit 1 for each revert, exit 0 for each restore, and source hashes. `sensitivity-status.json` records the outer execution. These are three R1–R3 controls, not the initial review's five controls. Frozen source was never mutated.

## Affected offline and local CLI evidence

`check-offline.py` is an unchanged copy of the inspected student helper. `student-cli` is a hash-verified copy of the new pinned binary. The helper uses original captures from `../ch02-student-evidence` and old binary only for comparison; it does not invoke the live launcher or read credentials.

```sh
# cwd: $REPO
python3 "$EVIDENCE/check-offline.py"
```

`offline-checks-status.json`, stdout/stderr, per-CLI stdout/stderr, fixture/error logs, and `offline-checks-summary.json` preserve the actual execution. `offline-counts.json` derives and asserts the counts from its records:

- 3 new-binary dumps of original logs;
- 18 new-binary captured-log renders and 9 old-binary comparison renders;
- 6 affected keyless renders (R1, R3 exact/foreign, redaction across three vendors);
- 0 HTTP requests during offline work;
- 3 expected exit-1 local fake error-cost CLI calls; no provider calls;
- 39 CLI runs total: 36 exit 0, 3 expected exit 1.

Exactly three local fake HTTP paths are recorded. Original files remain hash-identical. These new-binary offline/local checks do NOT claim fresh live-provider validation: the six original paid text calls remain attributed to initial commit `74eb026…`, and the original launcher exit 1 remains preserved.

## Failures and limits

Expected negative results are retained for all three reverts and local CLI error cases. `inspection-note.txt` records a read-only inspection typo using nonexistent root `agent.go`, recovered as `internal/agent/agent.go`; it was not a candidate/test failure. No successful baseline failures or watchdog expirations occurred.

No unrelated five-control rerun, GUI runtime/build rerun, private-reference rereview, new paid call, credential read, installation, restart, push, commit, tag, manuscript/grader/source edit or author-phase action occurred. This round's verdict is scoped to the exact pinned correction candidate.

Verify checksums from `$REPO`:

```sh
shasum -a 256 -c docs/edition-3-notes/ch02-code-review-r1-evidence/SHA256SUMS
```
