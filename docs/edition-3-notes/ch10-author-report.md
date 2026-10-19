# Chapter 10 bounded author report

**Status: author pass complete; parent-owned independent text review pending.**
Code and operations were already accepted. This pass changes only Chapter 10
prose and its new author receipts; it does not declare text acceptance.

## What changed

- Replaced the obsolete Edition 2 launch with the current Edition 3 route:
  five module roots, declared streaming/GUI models, optional inherited fixture
  inputs, and an explicit CLI executable. Added the same-binary Chapter 1 check,
  its real delayed phases, and the `passed: true` requirement.
- Distinguished an importable library root from an executable package. The
  concrete accepted CLI example builds `./cmd/ch01`; root `.` is conditional on
  an actual `main` package.
- Corrected the manual exercise around the accepted CLI's model environment,
  flags-before-subcommand parsing, human chat versus live JSONL envelope,
  scratch working directory, real tool turn, EOF, same-path relaunch and override.
- Explained why provider tool correlations and stable browser card addresses
  are separate identities, including already-colliding legacy IDs. Added the
  owned GUI restart/card-control exercise without prescribing a new public API.
- Preserved all numbered technical sections, persistence invariants, explanatory
  code, grader check table, historical evidence links and source-test caveats.
  Marked the older effective-prompt story as Edition 2 history. Added the fresh
  student's frozen failure, corrected testing path, and bounded operational
  results with transport, fake-provider and screenshot-review limitations.

The JSON companion gives section-by-section coverage. Nothing in Chapters 1–9,
application modules, grader implementation, existing evidence or acceptance
status was edited.

## Read and source coverage

All 18 mandatory files were read in full before the first write. Long or omitted
combined displays were recovered in consecutive smaller reads; no compaction
was used. `ch10-author-evidence/required-reads.json` records paths, byte/line
counts and SHA-256 digests, including the task digest. Additional narrow reads
covered the accepted CLI source/tests and README/CHAPTER10, public grading
instructions and prose-linter rules. The source pins are in
`protected-before.json`; `protected-after.json` verifies all 22 scoped paths
unchanged. Existing settings and credentials were not inspected or exported.
Historical chapter links were checked for existence only, without opening their
solution or reviewer contents.

## Text-only checks

`go run ./cmd/lintprose -v book/edition-3/chapter-10.md` ran before and after the
edit from the book repository, with a minimal recorded environment, offline
module lookup, read-only module mode, local toolchain and bounded parallelism.
Both baseline and final exits are **0**, with **no hard failures**. Negation
count stays **6**; the person-gap warning improves from **2,175** to **1,297**
words (threshold 1,200). That remaining inherited soft warning is disclosed,
not suppressed. Intermediate results are retained separately from final results.

`git diff --check -- book/edition-3/chapter-10.md` exits **0**. All four shell
fences pass `/bin/sh -n`; fences balance and all five relative links exist.
These checks establish text/syntax consistency only. The author did not run the
manual recipe, grader, application, browser or provider. Command semantics were
checked against the accepted CLI source and current public student brief.

The failed first edit anchor caused no write; baseline hash equality was checked.
A chapter-only diff opened a pager; only its identified child pager was stopped,
and the owned job was joined at exit zero. No children were created. All jobs
are finished. Disk stayed above the 1 GiB floor.

## Evidence boundary and handoff

The chapter preserves the original failed 80/100 report and attributes the final
100/100 and same-binary Chapter 1 success to the accepted corrected evidence.
The live tool/restart report covers one declared model; the original transport
failure remains unresolved in cause. The owned GUI exercise used a fake provider
and DOM assertions. Its screenshots were not a human visual/audio review, and
this author pass adds no such review.

`seal.json` pins the final chapter, both reports and every text evidence file
(excluding the seal itself). The protected-path check is deliberately scoped;
it is not an audit of private evidence or the whole filesystem. Git status before
and after records the unrelated work preserved in place. Root owns the **one
independent text review** and any subsequent text-acceptance decision.
