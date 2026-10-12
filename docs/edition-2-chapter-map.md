# Edition 2 chapter mapping

## Edition 1 endpoint and new Edition 2 material, October 11

Bill freezes Edition 1 at Chapter 27. The three existing chapter merges remain,
so original Chapter 27 (Gateway) maps to Edition 2 Chapter 24. Its Edition 2
implementation and enhancement are not yet complete; the accepted export table
below still ends at Chapter 23. Do not create an accepted export by renumbering
unfinished work.

The first wholly new Edition 2 chapter follows as Chapter 25: the temporary
Codex-endpoint workaround for the public OAuth caching problem. Subsequent new
chapters are Edition 2 only, with no corresponding first-edition Chapter 28.
Bill's current instruction is to verify the public endpoint again before using
the result in the cold open. Preserve Edition 1 and its solution unchanged.

Bill authorized the coordinator on October 11, 2026 to choose both remaining
merges and proceed autonomously. Exactly three first-edition standalone chapters
are absorbed: 5, 22 and 24. This supersedes the earlier placement stops.

| First edition | Second edition | Direct accepted solution | Placement |
| --- | --- | --- | --- |
| 0–3 | 0–3 | `ch01`–`ch03` (no Chapter 0 implementation) | Unchanged |
| 4 + 5 | 4 | `solutions/edition-2/ch04` | Jobs and public framework consumer |
| 6–20 | 5–19 | `solutions/edition-2/ch05`–`ch19` | Existing mapping |
| 21 + 22 | 20 | `solutions/edition-2/ch20` | Web access, status and accounting |
| 23 | 21 | `solutions/edition-2/ch21` | The Sandbox |
| 24 + 25 | 22 | `solutions/edition-2/ch22` | The Ensemble, including reusable GUI boundary |
| 26 | 23 | `solutions/edition-2/ch23` | Fleet management and GUI; accepted |

With first-edition Chapters 1–26 complete, Edition 2 has Chapters 1–23. Chapter
0 remains introductory in both. Future chapters, if separately authorized,
retain the three-chapter offset. Do not create standalone Edition 2 refactor
chapters.

## Why these merges

Original Chapter 22 completes original Chapter 21 before sandboxing: the agent
gains status access through its existing owners, accurate per-model accounting
and a practical model-switch exercise. Most architectural repair lessons already
govern the student from the beginning; do not recreate their old defects.
Preserve the status feature and structural/runtime verification, including real
executable composition and observer-disconnect behavior. Keep the sandbox
chapter focused on confinement rather than inserting an unrelated architecture
detour.

Original Chapter 24 belongs with original Chapter 25 because a reusable frontend
boundary and public event vocabulary support application assembly around the
agent fleet. Preserve embedded deployment with disk override, component-only
reuse, public MCP relay, external consumers, transitive dependency controls and
unchanged wire behavior. The GUI is already a separate optional module; do not
copy the source's callback bundle or reintroduce coupling just to remove it.
Original Chapter 25 remains the terminal fleet exercise; multi-agent GUI
features wait for original Chapter 26. The later prose should use coherent
sections within the combined chapter, with the right teaching length rather than
a target word reduction.

## Execution and preservation

Current checkpoint: original26 is independently accepted at `0d9e823` after
one correction round, source `450f7a9`, main tree
`27d7fcf9b8b31275068fa705d6969cfe3eafedc1`. Its exact direct `ch23` export is
checkpointed as `edition-2-ch26-r1`. Direct chapters1–23 are accepted and match
the table. Original26 was complete at original25 review, satisfying the condition.
Original24's intermediate tag, the combined24+25 `ch22`, and prior tags remain.
Original results and accepted exceptions are preserved; later prose/grader
enhancement remains separate from these completed implementations.

Use original exercises and unchanged graders in original order: 22, 23, 24, 25.
Each original exercise still receives its own student review, actual grader/live
results, independent comparative review and bounded correction loop. Notes,
evidence and original-numbered tags keep original numbers. Only the book-facing
chapter/solution numbering follows the table above.

Extend `solutions/edition-2/main` throughout. After original Chapter 22
acceptance, replace the direct `ch20` export with the newly reviewed cumulative
original 21 + 22 source and a new manifest/tag. The prior original-21-only
export and manifest remain exactly preserved in immutable tag
`edition-2-ch21-r1`; no existing tag is moved and no new history directory or
symlink is created. This specific replacement is authorized by Bill's requested
merge. Original Chapter 23 acceptance creates direct `ch21`. Review and
checkpoint original Chapter 24 as an intermediate stage in Git, then original
Chapter 25 acceptance creates the combined direct `ch22`. Never publish an
incomplete combined exercise as an accepted final chapter. Original Chapter 26
acceptance, if reached, creates `ch23`.

At original Chapter 25 review completion, inspect original Chapter 26's actual
manuscript, grader, reference snapshot and completion records. A snapshot alone
does not establish that the chapter is complete. If completed, read the full
source, update carryover coverage and proceed autonomously with original Chapter
26→Edition 2 Chapter 23. If still ongoing, checkpoint through `ch22` and report
that condition without inventing requirements. Preserve CodeRhapsody's commits,
index entries and working files at every step.

GPT-6 Astra/high remains mandatory for coding and code review. The complete
implementation precedes the later author–coder manuscript/grader enhancement
phase; this map directs that phase without starting an author ahead of the
coder. All inherited architecture, no-mocks, evidence and safe shared-commit
rules remain.
