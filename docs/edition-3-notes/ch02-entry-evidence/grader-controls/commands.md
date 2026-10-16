# Full-duration Chapter 2 composition controls

Private author/auditor evidence; do not expose historical reference details to a future student. This follow-up performs only the previously missing real-duration controls. No student or tracked reference/grader/chapter source changes, no extra mutations, no coverage expansion, no live providers or credentials.

## Inputs and environment

- Repository HEAD: `5351583208988f8370937cd0494727645efa8b46`.
- `$REPO`: grader/book repository root; local absolute paths are sanitized in reports/statuses.
- `$SCRATCH`: fresh disposable directory, created with `mktemp -d /tmp/e3-ch02-full-controls.XXXXXX`.
- Each build/grader subprocess inherits only `PATH`, `HOME`, `TMPDIR`, `GOPROXY=off`, `GOSUMDB=off`; no credential variables are inherited or recorded.
- `source-hashes.json` records SHA-256 for every historical/scratch source file. Exactly `internal/engine/engine.go` differs: remove only `Timeout: 60 * time.Second` from its `http.Client` literal and the now-unused `time` import. All other bytes/files remain unchanged. No formatting sweep or broader repair.

## Commands executed

Each command below runs separately with the allowlisted environment above. Builds have a 90-second outer watchdog. Chapter 2 has a 90-second outer watchdog; each Edition 3 Chapter 1 invocation has a 180-second outer watchdog. These do not change the grader's real `generationDelay = 62 * time.Second` or either header/body probe.

```sh
# Prepare positive-source by copying $REPO/solutions/ch02 to $SCRATCH/positive-source,
# then apply only the two removals described above and verify the source manifest.

(cd "$REPO/solutions/ch02" && go build -o "$SCRATCH/historical-ch02" .)
(cd "$SCRATCH/positive-source" && go build -o "$SCRATCH/positive-ch02" .)
(cd "$REPO" && go build -o "$SCRATCH/grade-edition2" ./cmd/grade-edition2)
(cd "$REPO" && go build -o "$SCRATCH/grade-edition3" ./cmd/grade-edition3)

# Run sequentially from $REPO; preserve the negative exit code rather than aborting.
"$SCRATCH/grade-edition3" -ch 1 -json "$SCRATCH/historical-ch02"
"$SCRATCH/grade-edition2" -ch 2 -json "$SCRATCH/positive-ch02"
"$SCRATCH/grade-edition3" -ch 1 -json "$SCRATCH/positive-ch02"
```

No rebuild occurs between the two positive grader invocations. `binary-hashes.json` identifies the explicit executables, and `statuses.json` records target hashes before/after each grader, elapsed monotonic seconds, UTC start timestamps, exact sanitized argv/cwd, exit status, watchdog status, report score/Passed and failed check IDs. Each command has `.stdout`/`.stderr` files; grader runs also have parsed `.report.json`. `integrity.json` is written only after all three runs finish and source manifests are rechecked. A nonzero negative-control exit is expected evidence, not a failed setup to repair. A build/watchdog failure is a stop condition, not permission to expand scope.

The receipt in `../../ch02-entry-audit.md` records actual outcomes after completion. These reference controls test the two-command composition; they are not acceptance or coverage evidence for an unlaunched Chapter 2 student implementation.
