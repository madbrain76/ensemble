# Edition 3 working manuscript

This is a working copy of the current book, not a newly published edition.
Read the chapter sources in `chapters.md` order. The original `book/` tree and
Edition 2 evidence remain unchanged.

- [Execution agreement](../../docs/edition-3-execution.md): execute first;
  author clarifications follow student completion and passing grading/code
  acceptance, then one text-review pass and one author revision.
- [Carryover](../../docs/edition-3-carryover.md): observed crossover corrections,
  reading coverage, and remaining execution gates.
- [Baseline manifest](../../docs/edition-3-baseline.json): exact copied files.

`solutions/edition-3/` is the copied application baseline, not evidence of a
fresh student build. Student predecessor trees and grader paths must be selected
explicitly before execution. Existing grader commands in the chapters name the
retained Edition 2 harness; they do not claim the changed contracts are tested.

The copied assembled Markdown, PDF, EPUB and EPUB resource file remain Edition 2
artifacts. They have not been rebuilt from these edited sources. Root build
commands target the published `book/` tree; do not use them to overwrite it when
building this working edition. An isolated Edition 3 publication build is still
to be prepared. Bill's preface and attributed historical accounts are preserved.

The [baseline publication guide](../README.md) retains the EPUB navigation,
print typography, pagination and release checks; reuse those checks with
isolated output paths when building Edition 3.
