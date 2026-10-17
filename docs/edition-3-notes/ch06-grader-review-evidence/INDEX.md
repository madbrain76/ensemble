# Independent Chapter 6 adapter review evidence

Scope: exact frozen root adapter, not student acceptance or entry re-audit. Report: `../ch06-grader-review.md`.

## Authoritative receipts

- `commands.json`: full recorded command metadata (argv, cwd, UTC times, exits, durations, environment policy, disk, stdout/stderr hashes).
- Each `<label>.meta.json`, `<label>.stdout`, `<label>.stderr`: unabridged command receipts. Intentional negative exits are not hidden.
- `results.json`: seven completed grader outcomes on their actual score scales.
- `positive.*`: exact candidate, distinct binary/core/consumer/workflow/stream-exercise identities, known historical model: 100/100, exit 0.
- `wrong-stream-model.*`: unlisted but syntactically valid model reaches grading: 35/100, exit 1.
- `wrong-stream-exercise.*`: regular-go.mod source-less wrong module reaches dependency/exercise proof: 90/100, exit 1.
- `custom-case-media.*`: valid-but-wrong case and media map both visibly forwarded: 80/100, exit 1.
- `legacy-ch4.*`, `legacy-ch5.*`: 160/160 and 100/100 respectively, exits 0.
- `legacy-ch1.*`: accepted frozen f60 CLI, 100/100, exit 0, 127.147 seconds. Both actual 62-second header/body delay checks finished naturally.
- `measurement-sequence.*`: entire sequential seven-grader sequence, exit 0 with every intended exit verified.
- `unit.*`, `race.*`, `vet.*`, `whitespace.*`, `build-sentinel.*`: all exit 0.
- `gates-r2.*`, `runtime-gates.json`: completed 61-case runtime gate matrix and direct compiled-sentinel positive self-check. Full per-case argv/stdout/stderr/exits/durations retained. Marker remained absent for all candidate gate cases.
- `final-integrity.*`: full final source/pin/evidence preservation check, exit 0.

## Identities and source review

- `candidate-pins.json`: root's frozen candidate provenance, independently verified before/after.
- `source-binary-pins.json`: candidate, historical core, accepted f60 CLI, compiled sentinel, and source-manifest hashes.
- `initial-tracked.diff`: exact root-supplied two-file diff, unchanged by reviewer.
- `source-comparison.json`: exact legacy suffix equality after removing the Chapter 6 branch, line counts, and diff size.
- `prerequisite-verification.json`: retained full prerequisite reads, unchanged hashes, frozen entry report/evidence validation.
- `final-verification.json`: final candidate/input/source/student/FINAL/entry preservation, output hashes, and minimum disk headroom.
- `root-*`: copied root receipts, initial pre-write guard failure, original/refined tests, source pins, build record, final race/vet/whitespace evidence. These are root provenance, not falsely fresh reviewer runs.
- `root-test-refinement.diff`: independently retained initial/final root test comparison.

## Reproduction and setup history

- `run.py`: no-deadline sequential command recorder; credentials excluded from inherited environment.
- `measure.py`: exact seven-grader command sequence. Expected negative exits are explicit.
- `gates.py`: exact compiled-sentinel runtime matrix.
- `sentinel-source.go.txt`: compiled sentinel source; writes only its owned marker and exits 99.
- `custom-consumer-case.json`, `custom-media-models.json`: actually used forwarding-negative controls.
- `explicit-default-*.json`: prepared but unrun configurations; not credited as separate runtime evidence.
- `finalize.py`: final read-only assertion/summary generator; no new behavioral test.
- `scratch-path.txt`: `/tmp/ensemble-ch06-adapter-review-c6v5ft04`; only owned scratch is frozen. Root's candidate executable and entry scratch were not modified.
- `setup-notes.md`: every retained review setup/observation mistake and its recovery, separate from candidate behavior.
- `gates.*`, `initial-gates.py`, `initial-runtime-gates.json`: first harness attempt stopped on its own diagnostic assertion; no grading started, no candidate launch, not credited as a completed matrix.
- `initial-measure.py`: unrun draft lacking `-json`, fixed before any grading.

## Freeze

- `REPORT-SHA256`: binds the sibling review report.
- `SHA256SUMS`: SHA-256 of every evidence file except itself.
- `freeze-status.json`: exact file count and completed-job state.

All behavioral/build/test/observation jobs completed. No children or pending jobs. No application/student/manuscript/grader edits by reviewer, no student launch or solution supplied, and no original entry artifact changed. Verdict is PASS for this exact thin adapter only, with all entry-audit source/integration limitations carried forward.
