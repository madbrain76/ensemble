# Edition 2 chapter mapping

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
| 26 | 23 | `solutions/edition-2/ch23` | Conditional on completed original Chapter 26 |

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

Current checkpoint: original Chapter 23 / Edition 2 Chapter 21 is independently
accepted. Direct `ch21` preserves tree `9dc3e9e36a854ede4cc2cbef22dc389850cedc0c`
under `edition-2-ch23-r1`. Combined original 21 + 22 remains direct `ch20` at
`edition-2-ch22-r1`; its earlier original-21-only snapshot remains in its prior
tag. Original24 is now independently accepted at `d747be7`, source `0430df4`,
main tree `8555e3cb9ab5200f0a1c52691007a7fc4a1c6797`, intermediate checkpoint
`edition-2-ch24-r1`. Original25 follows before the combined `ch22` export.

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
