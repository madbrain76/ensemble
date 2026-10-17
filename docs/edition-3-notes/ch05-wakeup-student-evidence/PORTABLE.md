# Chapter 5 public wakeup amendment: portable evidence

This directory is text-only and student-authored, not a reference solution.
The adjacent `ch05-wakeup-student.md` is the bounded report;
`ch05-wakeup-student-source-manifest.json` pins the exact48-file final source,
archive and race binaries; `ch05-wakeup-student-evidence-manifest.json` hashes this
copied evidence, that report, and the source index. Nothing here was auto-staged.

Run `python3 verify-portable.py` from this copied directory after relocation.
It only reads/hashes files: no candidate execution, credential loading, network,
or provider call. Copy the report and both indexes alongside this directory.
Archive/binaries remain outside Git; their paths and hashes are in the indexes.
`source-snapshot/` contains the exact clean Git/worktree/archive text files, not
build outputs. `verify-final-source.py` documents the full local byte checks;
its old absolute paths must be adapted intentionally for a different machine.

## Primary evidence

- `input-pins.json`: fully read permitted public inputs and amended chapter pins.
- `implementation-notes.md`, `attempt-ledger.md`: plan, lifetimes/assistance and all
  unsuccessful attempts. Original/correction evidence was not rewritten.
- `compatibility-focused-r4.{json,log}`: ten corrected race repetitions plus R1–R3
  and inherited PTY compatibility, exit0.
- `sensitivity-compatible/{source-controls.json,commands.json,*.log}`: seven
  original6c controls compile and fail behaviorally, complete corrected restore
  passes with the same public tests. Invalid setup attempts stay uncredited.
- `candidate-r2/gates/commands.json`, individual logs/coverage: all13 final-byte
  root/GUI/consumer/workflow vet/race/test/build exits. CLI is root cmd/ch01.
- `candidate-r2/frozen/identity.json`, `final-*-buildinfo.txt`,
  `final-source-verification.json`: clean source/archive/race-binary identities.
- `local-wakeup-r5/summary.json` plus actual safe executable stdout/stderr,
  loopback HTTP bodies, real retained output/files and saved histories: initial
  and later3600s process wakes, ordinary-handler wake, new effects while old jobs
  live, no skipped effects, and EOF's honest late-handler join. All HTTP data is
  synthetic/local; no paid/real-provider payload, private conversation or opaque
  thinking data was collected. Earlier launcher failures remain separate.
- `candidate-r2/frozen/grade-commands-and-exits.json` and both grade JSON reports:
  required exact same CLI public regression gates. The historical grader does
  not test the new amendment; use the public wakeup tests and independent review.
- `prior-preservation.json`: original341+correction83 text sets (portable and
  local), all19 prior excluded artifacts, pinned reports/indexes and original
  46-file snapshot still match. No paid proof is attributed to the amendment.

Scratch sensitivity source trees, generated Python bytecode, binaries and tar
archives are not copied. Scratch provenance/restoration hashes and their actual
command outputs remain. The exact new source snapshot and tests allow replay;
original6c is available by its pinned Git commit and local archived control.

## Scope/credit

Final amendment commit f60ebf3aad2ece4251face9d95c106dd667abcd6 extends reviewed6c.
The formal correction count stays one; this is authorized new public work.
No child agent, paid call, credential lookup, GUI, install/restart/push, acceptance
tag or main commit occurred. The original12 paid calls remain initial00489-only.
Independent whole-candidate review/acceptance is still the parent's responsibility.
