# Chapter 5 text-review evidence

This directory belongs only to the single normal post-student prose review in
`../ch05-text-review.md`. It contains no new source tests or grading campaign.

## Primary receipts

- `frozen-pins.json`, `input-hashes.json`: exact reviewed input identities.
- `author-manifest-verification.json`: all 23 original author entries verified.
- `baseline/chapter-05.md`, `candidate/chapter-05.md`, `chapter.diff`: exact reviewed
  snapshots and author diff; the original chapter is not edited by the reviewer.
- `preservation.json`: complete insertion text/ranges and preservation checks.
- `commands.json`: argv, absolute cwd, environment overrides, exit, duration and
  filenames for every command in the successful validation run.
- `lint-{baseline,candidate}.{stdout,stderr,exit}.txt`: full independent lint runs,
  same `chapter-05.md` basename; both exit 0. All lint output except the first-line
  path exactly matches the original author runs.
- `chapter-whitespace.*`, `local-links.json`, `scope-validation.json`, full
  `status-*`, `tracked-*`, and `staged-*`: bounded whitespace, link and scope proof.
- `validation.{stdout,stderr,exit}.txt`: successful runner's complete streams and
  exit. `validate.py` is the exact successful runner, not implementation source.
- `author-receipts/`: original full lint output/exits, scope/final-check failures,
  corrected checks, reconciliation, and parent-confirmed execution append.
- `historical-pty/`: manifest-verified intermediate root test/command log and
  final public grade/command receipt, copied without running either test/grade.
- `final-checks.json`, `seal.py`: final report/link/scope and index-sealing check.
- `SHA256SUMS`: evidence-directory-relative paths; includes the report as
  `../ch05-text-review.md`, excludes only itself. From this directory,
  `shasum -a 256 -c SHA256SUMS` verifies the sealed files without executing code.
  Rerunning either writing validation script changes evidence and requires a
  deliberate reseal; do not confuse a modified bundle with this submitted pin.

## Unsuccessful reviewer operations (no manuscript findings)

Before validation, a combined exploratory `cat` assumed these nonexistent
receipt names under the author directory: `baseline-prose-lint.stdout`,
`baseline-prose-lint.stderr`, `final-prose-lint.stdout`,
`final-prose-lint.stderr`, `command-receipts.json`, `scope-and-preservation.json`,
`validation.json`, `failed-scope-receipt.json`, and `reconciliation.md`.
They returned `No such file or directory`. The actual manifest names were then
read in full. No check was credited to those failed reads.

Exploratory `find`/`cat`/`grep` also tried absent portable
`ch05-wakeup-student-evidence/gates-first/` and absent
`frozen/grade-ch05.json` / `frozen/grade-commands-and-exits.json`.
A conditional check of `/tmp/cadence-ch05-wakeup-implementation/gates-first`
found no directory. Inspection of the actual portable manifest/layout located
`gates/root-race.txt` and `gates/commands.json` (intermediate failing candidate)
and `candidate-r2/frozen/grade-*.json` (accepted candidate). These are the saved
raw receipts used, not invented missing logs. They agree with the complete
student and review reports. The failed lookup details here are a narrative
record; the full streams retained below are for the executable validator runs.

Three validator launches failed before a successful fourth run:

1. `failed-validation-attempt/`: manifest entries incorrectly resolved from the
   repository root. FileNotFoundError; exit 1. Fixed to the author evidence root.
2. `failed-validation-attempt-2/`: all 23 manifest entries verified, but an assumed
   `candidate/chapter-05.md` author snapshot path failed. Exit 1. Fixed to the
   manifest's `final/chapter-05.md`.
3. `failed-validation-attempt-3/`: pins, preservation and receipt-copy checks
   reached lint, but `go run ./cmd/prose-lint` named a nonexistent directory.
   Lint exit 1 and runner assertion exit 1. The actual saved author command is
   `go run ./cmd/lintprose -v`; corrected run uses it with offline Go settings.

Each failed directory retains its exact runner and complete runner
stdout/stderr/exit, command receipt, and then-existing top-level validation
files. Some earlier checks succeeded within failed runs; only the completed
successful run is used as overall validation credit. No failed attempt changed
manuscript, code, graders, guides, Git index, or existing evidence.

## Limits

Hash scope is the named input set plus all five previously modified tracked
files. Before/after status captures cover all 1,735 pre-existing untracked file
paths, not byte hashes of all unrelated untracked content. Link validation checks
explicit local inline Markdown file destinations, not fragment accuracy or
remote URLs. Governing records and full supplied reviews, not implementation
source, are the basis for contract fidelity. Historical grader output is not
exhaustive wake-up proof. No provider call, installation, code test, grader run,
commit, tag, push, restart or child agent was used.
