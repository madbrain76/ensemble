# Initial independent Chapter 2 review evidence

Review: `../ch02-code-review.md`. Frozen student candidate: `74eb0265bd8bdbc7857d579704300b317b2365da`. This directory is private reviewer evidence, not a student answer source.

## Identity / integrity

- `governing-hashes.json`: unchanged governing sources checked against authoritative `5351583208988f8370937cd0494727645efa8b46`, following prior complete reads and current reviewer-skill reread.
- `identity.json`, `binary-build-identity.stdout`, `final-integrity.json`: commit/archive/binary hashes, byte-identical independent build, clean frozen candidate and unchanged chapter/report.
- `retained-manifest-verification.json`: independent verification of all 126 submitted manifest entries, bytes and hashes.
- `source-review-metrics.json`: descriptive private source comparison; not a quota or an acceptance gate.
- `SHA256SUMS`: report/evidence checksums; paths are relative to the authoritative repository root. Verify there with `shasum -a 256 -c docs/edition-3-notes/ch02-code-review-evidence/SHA256SUMS`.

## Source and real-duration grader batch

`run-statuses.json` records exact sanitized argv/cwd, time, exit, watchdog, and target hashes; `<name>.stdout/.stderr` preserve each command's output. `grade-e2-ch02.json` and `grade-e3-ch01.json` preserve complete reports. Read `passed`, exits, and failed checks—not only points.

Environment allowlist: `PATH`, `HOME`, `TMPDIR`; Go fetching disabled via `GOPROXY=off GOSUMDB=off`; no provider key variables or credential reads. `$STUDENT` is `/Users/bill/projects/ensemble-edition-3-student`; `$REPO` is `/Users/bill/projects/ensemble`; `$SCRATCH` is a reviewer-only disposable directory. Commands:

```sh
# cwd: $STUDENT
 go test -race -count=1 -timeout=180s -coverpkg=./... \
   -coverprofile="$REPO/docs/edition-3-notes/ch02-code-review-evidence/reviewer-cover.out" ./...
 go vet ./...
 go build -o "$SCRATCH/student-cli" ./cmd/ch01
# cwd: $STUDENT/gui
 go test -count=1 ./...
 go build -o "$SCRATCH/gui-stub" .
# cwd: $REPO
 go version -m "$REPO/docs/edition-3-notes/ch02-student-evidence/student-cli"
 go build -o "$SCRATCH/grade-edition2" ./cmd/grade-edition2
 go build -o "$SCRATCH/grade-edition3" ./cmd/grade-edition3
 "$SCRATCH/grade-edition2" -ch 2 -json "$REPO/docs/edition-3-notes/ch02-student-evidence/student-cli"
 "$SCRATCH/grade-edition3" -ch 1 -json "$REPO/docs/edition-3-notes/ch02-student-evidence/student-cli"
```

The source suite ran real 62-second header/body tests; the Edition 3 grader separately ran its actual 62-second phases. The source process watchdog was 210 seconds and Edition 3 grader watchdog 180 seconds; neither expired. No substituted short test controls.

`reviewer-cover.out` and `coverage-summary.json`: union duplicate instrumentation blocks by source range and statement count; hits from any package test count. Result 745/839 = 88.796%; root module including command, excluding nested GUI/test source/external process coverage.

## External public-root regressions

`public-test.go.txt` contains the exact final Go tests. Copy it to a scratch `public_test.go`, create a scratch module requiring `ensemble v0.0.0` with `replace ensemble => $STUDENT`, and run:

```sh
go test -race -count=1 -timeout=30s -v ./...
```

No internal imports or mock production interfaces. Uses public root/real agents/logs/engines, with local HTTP servers where needed. `public-smoke-status.json` and `.stdout/.stderr` retain the failing R1/R2/R3 cases and passing construction/isolation case. `initial-public-*` preserve the reviewer's earlier compilation failure caused by calling root construction methods through the narrow parent interface; that is not a candidate defect.

## Mutation sensitivity

Candidate `scripts/check-mutations.py` was executed in a scratch extraction of the exact candidate `git archive`, with output directed to `mutation-rechecks/`:

```sh
# cwd: $SCRATCH/mutation-source
python3 scripts/check-mutations.py "$REPO/docs/edition-3-notes/ch02-code-review-evidence/mutation-rechecks"
```

Each of five controls fails when mutated and passes when restored. `mutation-rechecks-status.json`, `.stdout/.stderr`, and individual records preserve both outcomes. `final-integrity.json` verifies every frozen tracked file equals the restored scratch counterpart. No candidate write occurred.

## Retained live evidence and independent offline checks

No reviewer live calls or credential reads. The original student's three live CLI processes exited 0, but the live launcher exited 1 at its incorrect final regex assertion; those original artifacts remain unchanged.

1. Copy the retained evidence directory to scratch and run its `python3 verify-captures.py` there. `offline-verifier-status.json`, `.stdout/.stderr`, and `independent-capture-verification.json` retain this successful offline-only rerun.
2. From `$REPO`, run the independent helper with the submitted bundle and this evidence destination:

```sh
python3 docs/edition-3-notes/ch02-code-review-evidence/offline-recheck.py \
  docs/edition-3-notes/ch02-student-evidence \
  docs/edition-3-notes/ch02-code-review-evidence
```

The helper verifies manifest hashes/sizes, executes 3 fresh dumps and 18 keyless renders of retained captures against a loopback POST sentinel, and compares results to retained bytes and hashes. It uses the actual `render LOG` and `CH02_LOG`/`dump` interfaces. `independent-offline-recheck.json` records output hashes, per-process exits/times, equality results, and zero POSTs; `.stdout/.stderr/.exit` preserve the helper result. Its first failed version and TypeError are retained under `initial-offline-*`; it failed before rerendering because the reviewer assumed the wrong summary shape. The corrected helper neither invokes the live launcher nor reads credential settings.

## Boundaries

Only this evidence directory and `../ch02-code-review.md` were written in the repository. No manuscript, grader, student, or reference source was modified. No GUI launch, install, restart, push, commit, new feature, live-provider repetition, or security-proof claim. The report requires three corrections and authorizes no acceptance or author/text phase.
