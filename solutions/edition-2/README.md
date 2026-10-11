# Edition 2 solutions

`ch01/` through `ch08/` match the enhanced book's chapter numbers. Each contains
an accepted cumulative solution. `main/` remains the authoritative working
implementation; extend it for future student work rather than editing snapshots.

The original student run has resumed beyond the enhanced book. Its new accepted
original-numbered exports live under [history](history/#resumed-original-student-run),
starting with original Chapter 10; `main/` advances after each independent review.
These new student checkpoints do not claim corresponding manuscript enhancements.

Bill authorized the solution renumbering on October 10, 2026. The code is
unchanged: the old Chapter 5 solution supplies the combined Chapter 4, including
jobs and the separately built custom-tool consumer. Later snapshots shift down
one. The intermediate old Chapter 4 is retained in the historical archive.

| Book chapter | Solution | Original student chapter | Historical checkpoint |
| --- | --- | --- | --- |
| 1 | [ch01](ch01/) | 1 | `edition-2-ch01-r1` |
| 2 | [ch02](ch02/) | 2 | `edition-2-ch02-r3` |
| 3 | [ch03](ch03/) | 3 | `edition-2-ch03-r3` |
| 4: jobs and external consumer | [ch04](ch04/) | 5 | `edition-2-ch05-r2` |
| 5: actor and steering | [ch05](ch05/) | 6 | `edition-2-ch06-r3` |
| 6: streaming | [ch06](ch06/) | 7 | `edition-2-ch07-r4` |
| 7: GUI and replay | [ch07](ch07/) | 8 | `edition-2-ch08-r3` |
| 8: workbench and settings | [ch08](ch08/) | 9 | `edition-2-ch09-r2` |

From the repository root, use the Edition 2 grader. Its default solution path
now matches the requested chapter:

```sh
go run ./cmd/grade-edition2 -ch 4
go run ./cmd/grade-edition2 -ch 8
```

An explicit submission directory and the existing consumer/workflow/GUI flags
remain available. Chapter 4 combines the original jobs checks (100 points) with
external-consumer, topology and diagnostic checks (60). Later graders carry the
corrected earlier checks forward. See the
[enhancement record](../../docs/edition-2-enhancement.md) for reviewed results,
exceptions and the distinction between original and revised scores.

## Historical numbering

[history/](history/) preserves the six original `ch04`–`ch09` snapshots byte for
byte. Original Chapters 1–3 remain unchanged at their existing paths. The
[mapping manifest](history/mapping.json) records source commit, Git tree IDs,
file counts and content digests. Existing tags retain the original paths.

Earlier student and enhancement reports use the original numbering and describe
commands that were actually run then. For example, their old
`solutions/edition-2/ch05` now means [history/ch05](history/ch05/), whose current
book-numbered counterpart is [ch04](ch04/). Do not run an old command against the
new same-numbered directory without translating that path. The
[previous student-run README](history/student-run-readme.md) and
[chapter records](../../docs/edition-2-notes/) preserve original results.

Evidence directories (`evidence/chNN`), nested application names and their
interfaces retain historical numbering. In particular, the consumer remains
`ch05/`, the three-Agent workflow `ch06/`, and the streaming example `ch07/`
inside each applicable snapshot. Their relative module dependencies still point
to that snapshot's core. Snapshot-local READMEs and source comments likewise
retain their original labels. This outer mapping identifies their place in the
new book; no runtime or historical artifact was rewritten to change numbering.

Keep the archived and reader snapshots immutable. Preserve old tags and actual
results when a future correction is reviewed and checkpointed.
