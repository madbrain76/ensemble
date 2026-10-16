# Independent Chapter 1 initial-review evidence

Review: [../ch01-code-review-r1.md](../ch01-code-review-r1.md).

Candidate `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`, empty predecessor `69859da3328c3f710f8143203b99f934fed778a1`. Book/grader `7f3d069420112c2136bc95d1d551e12fad6964c1`. These are reviewer-produced artifacts, distinct from `../ch01-student-evidence/`.

## Exact source and commands

`provenance.json` records pinned revisions, SHA256s, the built executable hash, clean student status, and byte-for-byte verification of all archived tracked files against the candidate. `scratch-path.txt` identifies the disposable exact-commit archive used in this session; source was exported with `git archive d2baff64055ac33702a91c6d6ee5a4246e84c0c9`, not copied from an unpinned moving tree.

From the archive root:

```sh
go build -o "$SCRATCH/ch01" ./cmd/ch01
go vet ./...
go test -race -count=1 -v ./...
```

From `/Users/bill/projects/ensemble`:

```sh
go run ./cmd/grade-edition3 -ch 1 -json "$SCRATCH/ch01"
go test -race -count=1 -v ./internal/gradeedition3 ./internal/gradeedition2 \
  -run 'TestGenerationLifetimeControls|TestHistoricalTimeoutIsExceeded|TestChapter1'
```

`$SCRATCH` was the absolute path in `scratch-path.txt`; the grader received the absolute **binary** positional target, not a default directory. The test selector unintentionally also matched existing Chapter 10–19 tests; those incidental passes are not Chapter 1 evidence. No 90-second shell cutoff was applied. Tool callback/observation intervals did not cancel running work.

## Artifacts

- `build.txt`, `build.exit`: actual exact-archive build output/exit. The successful initial build was repeated during final hygiene to persist a distinct exit file.
- `vet.txt`, `vet.exit`: independent final vet rerun. Initial vet was also part of the successful combined test command.
- `tests.txt`, `tests.exit`: Go version and complete full race suite; exit 0, 125.261 s, header/body lifetime subtests 62.00/62.01 s.
- `format.txt`: `gofmt -l` on every archived Go file; empty means no formatting differences.
- `grader.json`, `grader.stderr`, `grader.exit`: unchanged grader; exit 0, 100/100, `passed: true`, all seven scored and both required lifetime checks true. Empty stderr expected.
- `grader-tests.txt`, `grader-tests.exit`: grader regression/control tests; exit 0.
- `directory-grader.stdout`, `directory-grader.stderr`, `directory-grader.exit`: separately reproduced exact directory invocation `go run ./cmd/grade-edition3 -ch 1 -json /Users/bill/projects/ensemble-edition-3-student` from main root. Failed before checks with permission denied, go-run exit 1 / grader status 2. Preserved as failed invocation/documentation evidence, not waived or called a pass.
- `mutations.py`, `mutations.txt`, `mutations.exit`: nine independent mutations of final-candidate disposable copies, exact replacements and test commands, all killed by behavioral failures; runner exit 0 means all expected failures occurred. No mutant entered either repository. Baseline/archive source was reverified afterward. The real one-minute deadline mutant failed after 62 s.
- `mutations-initial-partial.txt`, `mutations-initial-partial.exit`: first reviewer runner's six completed probes followed by a scratch-copy setup failure. `shutil.ignore_patterns('ch01')` mistakenly skipped the `cmd/ch01` directory as well as the root binary; corrected to exclude only the root binary. The runner raised `FileNotFoundError` for its missing temporary `cmd/ch01/main.go`, exit 1. This is a reviewer-harness failure, not a student defect, and not the final mutation result.
- `live.py`, `live.txt`, `live.exit`: independent reviewer real-provider session, not a fake/recorded replay. Successful fresh model discovery, selected `claude-sonnet-5`, three-turn PTY conversation remembering Cobalt Wren, cumulative input/output 430/187, EOF exit 0. Script reads only the authorized `directClaudeAPIKey` in memory, passes it to the intended child environment, and sanitizes before writing. Its 300-second reviewer safety watchdog is not a production generation deadline.
- `line-counts.json`: per-file/total active, comment-only, blank and physical Go lines for candidate production/tests, promoted Edition 2 Chapter 1 and archived original Chapter 1. Classification is leading `//` versus blank versus active; no block comments occur in these files. README and go.mod are counted separately in the review.
- `hygiene.txt`: final revision/source/credential-presence checks. No credential value is recorded.

## Interpretation limits

Passing tests/grader/live use do not settle source quality or fake fidelity: the formal review contains one required fixture correction. Nine mutants prove those nine sensitivities, not exhaustive correctness. The short grader unit controls are distinct from the actual candidate's real 62-second probes. Student prompt-experiment evidence at `60bf58a7728a96e27c2e13eea1afbf744889cb6e` was inspected, not independently repeated; it is not part of this reviewer's live run. The final candidate's implementation delta and evidence provenance are discussed in the review.

Scripts use disposable local paths and are preserved for audit, not installed tools. Reproduction requires a fresh exact archive and updating `scratch-path.txt`; live reproduction additionally requires the explicitly authorized credential and incurs provider use. Neither student source nor grader/reference/manuscript was edited. Only the permitted review/evidence paths were written in the main repository.
