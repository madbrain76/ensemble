# Original student snapshots

These are the byte-for-byte accepted exports of original Chapters 4–9 before
the October 10, 2026 solution renumbering. Original Chapters 1–3 remain at
`../ch01`–`../ch03`. The current book-numbered solutions are one level above;
see [the complete mapping](../README.md).

| Original export | Current Edition 2 solution |
| --- | --- |
| [ch04](ch04/) | Intermediate jobs-only snapshot; retained here |
| [ch05](ch05/) | [Chapter 4](../ch04/), including external consumer |
| [ch06](ch06/) | [Chapter 5](../ch05/) |
| [ch07](ch07/) | [Chapter 6](../ch06/) |
| [ch08](ch08/) | [Chapter 7](../ch07/) |
| [ch09](ch09/) | [Chapter 8](../ch08/) |

[mapping.json](mapping.json) binds every export to its original Git tree in
`c2df11b5f86c2abbcc14a827e59d263f766ee8e6`, with file counts and sorted
relative-path/SHA-256 manifest digests. Existing checkpoint tags are unchanged.
Both content and executable modes are preserved. Nested application directories
and relative module replacements retain their original structure.

[student-run-readme.md](student-run-readme.md) is the prior outer README copied
unchanged. Its relative paths and current-status language describe its original
location and date. Use the mapping above to locate those sources now.
[Student reviews](../../../docs/edition-2-notes/) and [run evidence](../evidence/)
retain their original chapter numbers, scores and command transcripts.

Do not edit these historical exports. This archive records the original student
run and is distinct from the retired author-first attempt.

## Resumed original student run

New accepted snapshots keep original chapter numbers here while the completed
book-facing Chapters 1–8 remain frozen. They extend the authoritative `../main/`.

| Original chapter | Snapshot and source manifest | Checkpoint |
| --- | --- | --- |
| 10: Skills | [ch10](ch10/), [manifest](ch10-manifest.json), 75 files | `edition-2-ch10-r2` |
| 11: Persistence | [ch11](ch11/), [manifest](ch11-manifest.json), 80 files | `edition-2-ch11-r1` |
| 12: MCP | [ch12](ch12/), [manifest](ch12-manifest.json), 91 files | `edition-2-ch12-r1` |
| 13: The Agent Sees Itself | [ch13](ch13/), [manifest](ch13-manifest.json), 93 files | `edition-2-ch13-r1` |

Chapter 10 is accepted with independently recorded exceptions. The full original
run was interrupted/unscored; focused unchanged skills checks earned 65/85, with
two undefined `think` name checks failing. See the [review](../../../docs/edition-2-notes/ch10-code-review.md).
This table supplements the immutable renumbering manifest above.

Chapter 11 is accepted after one correction round with inherited exceptions.
Its full original run remains interrupted/unscored; focused original checks
earned 70/85 on the initial submission, with final correction evidence and the
masked save-shape predicates documented in the [review](../../../docs/edition-2-notes/ch11-code-review.md).

Chapter 12 is accepted after one correction round with recorded grader exceptions.
The full original run remains interrupted/unscored; unchanged interim current
checks earned 80/90, with the fixed-path WebSocket check failing and recursive
parity unrun. Actual browser and runtime evidence, including the corrected slow
automatic-context case, is assessed in the [review](../../../docs/edition-2-notes/ch12-code-review.md).

Chapter 13 is accepted after one correction round with missing-type fixture and
structural parity exceptions. The complete original grade remains 0/100 FAIL.
Actual browser/startup evidence and the native-speech pause correction support
acceptance; both OpenAI live task failures remain recorded separately from their
capability observations. See the [review](../../../docs/edition-2-notes/ch13-code-review.md).
