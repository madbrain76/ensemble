# Chapter 8 author pass

Author: Codex, delegated bounded author. No independent text review commissioned.
Assignment: `/tmp/cadence-ch08-author-task.txt`, 6771 bytes, SHA-256
`861f523b8e20f9848373634d95b255cbff3c0ee9830fb00402c2f3c301cbf27e`.
Root baseline: `/tmp/cadence-ch08-text-phase-baseline/baseline.json`.

## Result and exact manuscript

One reader-facing block changed: the TL;DR grading recipe, immediately before
the check/points table. All prose before and after that block is byte-identical
to the assigned chapter baseline. Root owns the single formal text review;
this report does not claim prose acceptance.

| Chapter | Bytes | SHA-256 |
| --- | ---: | --- |
| Before | 26124 | `8eb2d52071ae59ca095bdad3431ba72903419fbec1eae4bf0779ec99600ba14d` |
| After | 27335 | `fb16fce247ad7ae955113ee9e1793681e17beee14b754b0b5505fa4f06c3c662` |

The [exact diff](ch08-author-evidence/chapter-08.diff) and same-basename
[before](ch08-author-evidence/baseline/chapter-08.md) /
[after](ch08-author-evidence/final/chapter-08.md) copies are saved.

### Change and rationale

Replace the stale Edition 2 single-core command with the current
[public Edition 3 recipe](../edition-3-student.md): an already-built core CLI,
five explicit module paths, declared stream/GUI model selections, optional
consumer-case/media-map paths before the executable, and Chapter 1 on that
same frozen CLI. Both reports require `passed=true` and both commands exit zero;
the two real 62-second lifetime phases are named. The chapter does not require
a particular student's CLI package layout, and no grader or student build was run.

The inherited-media clarification is earned by the
[initial student report](ch08-student-review.md): the unmapped initial run
scored 80/100 before explicit declared-model mappings restored the inherited
checks. The new paragraph states the three JSON keys and nonempty declared
model IDs without copying local IDs, absolute run paths, or a provider matrix.
It expressly limits the map to grader-input selection, not provider capability
discovery. Existing launch, restart, dependency-boundary, and parity explanations
remain in the block.

### Deliberate no-change decisions

The [code review](ch08-code-review.md) and
[closure](ch08-code-review-r1.md), corroborated by the
[revision student report](ch08-student-r1-review.md), identify one formal code
correction: the inherited GUI test fixture, not production behavior. It had to
exercise current page initialization/Actions routing and settings entrypoints
while preserving its earlier assertions. Production remained identical to the
initial `f97a6850468a9cdc6875507f2821ccf6c397654d` freeze. The accepted student
is `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`, tag `edition-3-ch08-r1`.
These are report-attributed identities, not a new inspection of student source.

No second fixture paragraph was added. Existing browser-click-path, settings
round-trip, and inherited-behavior teaching already explains the relevant lesson;
an implementation-specific fixture story would duplicate it. Likewise the budget
owner/next-turn boundary, default 200 and zero-unbounded convention, sparse patch,
save-before-publish ordering, native fractional speech rates, and visual-judgment
sections remain untouched. Historical chapter failures remain historical, not
claims about this student's implementation. No runtime-policy change or memory
migration was introduced.

## Reads, pins, and process qualifications

The [coverage record](ch08-author-evidence/read-coverage.json) enumerates all
11 required current prerequisites and all 10 bounded evidence inputs, with
actual prerequisite line counts and identities. The
[input pins](ch08-author-evidence/input-pins.json) preserve all 23 root-listed
inputs. The current linter was read in full; the extra text-reviewer skill was
identity-checked only, not represented as an author full read. No raw-evidence
sweep, historical implementation inspection, or current student-source inspection
was performed. Current Edition 3 phase rules took precedence over historical
workflow instructions.

**There were two assignment-process deviations.** The workflow read originally
stopped at line 2421; its final 13 lines were recovered during the closing audit,
after the initial edit. They concern original Chapter 17 and change no Chapter 8
instruction. Thus the strict full-prerequisites-before-edit requirement was not
completely met. Truncated combined voice/procedure outputs were otherwise
recovered with separate reads.

The initial preservation loop also hashed every root-listed protected tracked
file, including excluded publication files. This was an out-of-scope read even
though no publication content was displayed/interpreted, edited, or rebuilt.
Earlier broad receipts are retained, but their publication results are withdrawn
as preservation assurance. Root owns moving publication scope. Final checks
exclude those publication paths. See the complete
[failure ledger](ch08-author-evidence/failures.md), including the failed
read-only lookup of nonexistent `scripts/lint_docs.py` before locating the real
linter. These failures are not concealed as successful setup.

## Mechanical checks and limits

Both same-basename copies were checked sequentially from
`/Users/bill/projects/ensemble` using the same current linter:

```sh
go run ./cmd/lintprose /Users/bill/projects/ensemble/docs/edition-3-notes/ch08-author-evidence/baseline/chapter-08.md
go run ./cmd/lintprose /Users/bill/projects/ensemble/docs/edition-3-notes/ch08-author-evidence/final/chapter-08.md
git diff --check
```

Linter SHA-256 (8366 bytes):
`a0321d4ed35983ca909d56a456f9984a5bae4a782bb75814990656c4e41a5ea6`.
All three commands exit 0. [Baseline receipts](ch08-author-evidence/baseline-checks.json)
and [final checks](ch08-author-evidence/final-checks.json) give actual argv, cwd,
exit status, and log names. Both lint runs retain the soft negation warning
(11, allowance 6) and person-gap warning after line 11 (3044 before, 3156 after;
allowance 1200). Short closers remain 14; paragraph count changes 67 to 69.
No prose was altered merely to reduce these warnings. Lint is mechanical support,
not proof of narrative quality or factual correctness.

Final checks verify: all 22 non-chapter input pins unchanged; 21 non-publication
protected tracked pins unchanged; original and final chapter-copy identities;
unchanged prefix/suffix around the single edit; five existing relative chapter
links resolving and unchanged; and only Chapter 8 added to the root's tracked
status delta. Status comparison is not a publication preservation assertion.
Existing dirty guidance, linter, grader, agent, and publication work was neither
reverted nor claimed as this author's work. No old report or seal was rewritten.
The reproducible [check script](ch08-author-evidence/verify.py) writes only this
new evidence directory.

## Evidence attribution and exclusions

The [grader seal](ch08-grader-accepted.json) and
[final acceptance](ch08-accepted.json) supersede earlier pending statuses;
accepted implementation was not reopened. The
[scroll continuation](ch08-operational-scroll-continuation.md) records the
packaging worker's stream `INTERNAL_ERROR` after successful runs: root sealed
saved facts, not a fictional final worker submission. Local/native 1.75,
1.2345 (native float representation), and 5 clamped to 3 with a fresh fourth
cancel were measured; no new human listening report is asserted.

The accepted two-request Claude GUI run is distinct from the one-request
image-capable review of two captured screenshots in the
[visual report](ch08-visual-review-evidence/visual-review.md) and
[root decision](ch08-visual-review-evidence/root-decision.json). Neither this
author nor root/student is claimed to have personally inspected PNGs. No next
paid-request body contents, retained job-file pins, provider failure, extra live
runs, or new visual/audio observation is invented. These qualifications stay in
the report, not a new chapter audit ledger.

No children, review commission, external assistance request, paid/provider call,
key access, browser/audio exercise, student build/test, grader run, commit/tag,
cleanup, or shared-cache deletion occurred. All command jobs returned; none was
started in the background. Bounded accepted reports supplied the assistance.

The self-excluding [inventory](ch08-author-evidence/inventory.json) seals this
small UTF-8/no-NUL evidence set, with the report and live chapter explicitly
identified as external entries. The final seal check is saved alongside it.
Root may now inspect the final saved manuscript and commission its one formal
text review, with the process deviations above visible.
