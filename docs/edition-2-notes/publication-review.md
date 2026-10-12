# Publication preparation: independent review

## October 12: relocation and frontmatter

Reviewer: `reviewer_astra`, GPT-6 Astra/high. **Accept the byte-preserving
relocation and submitted frontmatter. Publication readiness remains pending**
current tooling, runnable paths/links, the coordinator's new epilogue and assembled
artifact inspection. This is not an external publication or completed crossover.

After compaction, reread the entire mandatory
`.agents/skills/ensemble-coding/SKILL.md` (373 lines), workflow (1–2371),
carryover (1–1046), voice (1–632) and chapter procedure (1–599), in bounded,
untruncated reads. All four documents were active together during this review.

### Relocation

Reviewed the complete coordinator script and relocation plan before application.
An isolated Git control found its combined `diff-tree -rz` argument invalid
(exit129), which would have interrupted index reconciliation after publication
of the commit. The coordinator changed it to separate `-r`, `-z` arguments with
`--no-renames`, and added a second empty-index guard before publication. The
corrected isolated control returns both old and new paths. No application-source
change or new test framework was needed.

Independently verified promotion
`564e316f3f8aeac2ec26fe8edfbd8cdf42aa350b` against source
`c35aab7380e75563d268d049ba5ae3e6c0564012` and the saved relocation state:

- All 9,020 committed entries equal the prepared tree.
- All 8,386 moved original mode/blob pairs equal their committed destinations.
- All 8,386 moved working-file hashes and 443 protected local-file hashes match
  at their new paths, including unfinished edits rather than just committed bytes.
- The eleven original dirty tracked files remain dirty in `agent/` and
  `edition-1/book/chapters.md`; they were not included as accepted changes.

The archive keeps the original grading module and its dependencies. The revised
grader retains shared support at the root; the old entry point and original
tests are archived. This avoids duplicating grader behavior merely to rename
the default edition. Directory renames also carry ignored contents; those are
not claimed as members of the explicit 443-file hash audit. The coordinator
reports resolving an optional-lock collision with one index-info update;
the independent checks above establish the resulting committed and working
bytes, not an uninterrupted execution claim.

### Frontmatter

Read all five promoted files: `book/preface.md`, `front-matter.md`,
`front-matter-print.md`, `note-from-bill.md` and `chapter-00.md`. Compared the
original preface, both title pages and Bill note in `edition-1/book/`, and the
Chapter0 diff for author `4663e67c` plus narrow lint correction `a77bf084`.

The preface retains the engineer's bet, deletion of the proof of concept, BASIC
book's first-page teaching format and the motivation for an agent tailored to
its reader. Unsupported financial interpretations, universal-first claims and
guaranteed self-improvement are removed rather than silently asserted as results
of this edition. The failed author-first attempt and student-first correction
explain the changed method without turning the preface into a review ledger.

Authorship names Bill, first-edition CodeRhapsody and second-edition Codex
accurately. Chapter24's combined-role exception is disclosed. The Bill note is
consistent with his explicit publication direction, without invented dialogue.
Gateway implementation/manuscript completion, the bounded Astra coding example,
the begun crossover and reserved Chapter25 are kept distinct from producing
Edition3. No new paid run, cost measurement or general reliability claim is made.

Chapter0's edits accurately qualify numeric grading, parity and mutation claims
and update the actual model/crossover status. Its inherited comparison rhetoric
and figures remain explicitly dated original material; this review does not
recertify them as current comparative measurements. The narrow lint edits do
not alter their factual content. Final current-path reconciliation and artifact
rendering are still separate gates below; in particular, the old `../../docs/`
Chapter0 link needs the coordinated relocation pass.

No blocking frontmatter finding remains. No reviewer changes to runtime,
originals, frozen exports, user work or the process on port8088 were made.

## Remaining publication checks

Current grader/build/export defaults, original archive execution, updated links
and example paths, the final epilogue, and Markdown/EPUB/print/cover-PDF output
remain to be reviewed after their owners submit them. Earlier chapter acceptance
and historical receipts retain their original sources and results.
