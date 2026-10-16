# Edition 3 crossover integration

Preparation-only prose integration against repository HEAD `2be995ebc12c6a7efb21cee3b9740a758489adda`.
This records the bounded authorized edits, including corrections to the earlier
unaccepted handoff. It is not a formal review, student execution result or
replacement for the separate crossover audit.

## Changed files and placement

All chapter paths below are under `book/edition-3/`:

- `chapter-01.md`: new Edition 3 grader command with positional implementation
  path; seven inherited checks plus required real 62-second header/body gates
  (about two minutes), `passed` rather than score alone, and historical status
  of old commands. Removed the added four-line timeout paragraph; the
  preexisting transport-timeout paragraph is unchanged.
- `chapter-03.md` and `chapter-06.md`: selected-model explanation-before-dispatch
  enforcement and actual stream-arrival ordering, including rejection before
  tool effects and bounded corrective retry.
- `chapter-10.md`: deterministic distinct presentation across restart for
  already-colliding archived IDs while preserving archive facts. No ID API or
  allocator strategy is prescribed.
- `chapter-16.md` and `chapter-20.md`: selected-model policy with reasoning
  summary default `auto`; root versus worker instruction composition and
  inheritance of resolved running configuration, not pending saved values.
- `chapter-21.md`: real GUI `--config` ownership, strict versioned
  `ApplicationConfig`, external credential references, private atomic writes,
  environment isolation, legacy launch compatibility, and visible active versus
  saved-next-launch settings with correlated saves. No new CLI loader scope.
- `chapter-22.md` and `chapter-23.md`: extend the same schema with memory then
  recall when their types exist; durable/live application before acknowledgment
  while preserving unrelated and pending provider values.
- `chapter-25.md` and `chapter-26.md`: historical-only crossover interlude,
  without a construction exercise or code checkpoint; Gateway starts from the
  exact accepted Chapter 24 commit. Chapter 26 changes only this boundary sentence.

This record is the only other file written for this integration. No runtime,
solution baseline, original grader, role/student guide, publication, commit or
preface edits were made by this worker.

## Checks and limitations

- `git --no-pager diff --check` passed for all eleven chapters.
- All **116** fenced blocks balance. Compared with HEAD, the only fenced-example
  change is Chapter 1's `grade-edition2` to `grade-edition3` command.
- All **19** local file-link occurrences in those chapters resolve. External
  URLs, fragment anchors and rendered publications were not checked.
- Against `docs/edition-3-baseline.json`, all **278** original source hashes and
  **231** copied application hashes match. Edition 3's Bill preface, note and
  essay also match their recorded baseline hashes.
- Protected tracked runtime, copied application and original grader paths have
  no diff from HEAD. Chapter 26's remaining bytes, including mixed line endings,
  are unchanged after removing the one inserted sentence.
- Inspected the new grader CLI, lifetime checks and report semantics; did not
  run its two-minute gates, live providers, browser tests or student execution.
  No implementation acceptance or independent review is claimed.
- The worker's initial **Chapter 1 draft** incorrectly used `-impl`; inspection
  of `cmd/grade-edition3/main.go` corrected it to a positional path. The parent's
  student brief was already positional; the worker's earlier contrary message
  was mistaken. An exact-anchor edit stopped safely before its unmatched
  replacement and was retried with verified anchors. A Git pager interruption
  during checks was closed and the check rerun with paging disabled.

## Final chapter SHA256

The record's own hash is supplied in the handoff, avoiding a self-referential
hash inside this file.

```text
ab016baaff3f6a9d83ecffcbfc3772cd692c9d33c06d9fa1c52ad9a8708d25df  book/edition-3/chapter-01.md
2366f74bdb873af57fe95e0c7c0b040187365a8db7efc09f5e63ce1f3c59c801  book/edition-3/chapter-03.md
ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55  book/edition-3/chapter-06.md
abda8ed3a74c2ee2c64bc2208ea104e0547810d6e06e1d5e772c0bb4ee84b3ed  book/edition-3/chapter-10.md
7e20c9116cc5d1f0d7514ca61a8501c342d1ac49b044f4ba283b28e1be0e1f28  book/edition-3/chapter-16.md
291038c32369f9e519afea9cb30d43023c17f82fd8945578410b00fc0f21fe2e  book/edition-3/chapter-20.md
356d762d00d8a5276f12c3f5c1b357e7c8298e6cb3f7ee5a4bf3220bfeae07cf  book/edition-3/chapter-21.md
d471365f586d0c083082933ea73bdc2f5cf0cc089d75885f71409a85c470d54e  book/edition-3/chapter-22.md
d0d155215945e7e943287ee49f2a90954b51b97e3198367e73528a30f26d4f43  book/edition-3/chapter-23.md
db4581c30727babd43cbbd0c46c9573d2d6db0599ebd293c84d6120f2dea6e78  book/edition-3/chapter-25.md
1a1d074520793c01ffd95735a95d00e8e2992f1821ec9a81714ed9c1bf4b2cad  book/edition-3/chapter-26.md
```
