# Chapter 10 retention amendment — completed

Bill requested an incremental lesson and implementation after diagnostic logs exhausted disk space. This is a separately authorized amendment to accepted Chapter 10, not a repeat of the original exercise or a retrospective failure of it.

## Accepted result

The student implemented the delta in an isolated worktree from accepted `729bc919e734405307b4ac0060e65c40e33652e5`. Accepted commit `b03c28fdb1586f16dfb53849f8c763576ebb39ae` has immutable tag `edition-3-ch10-r2`; `edition-3-ch10-r1` remains unchanged. The 18-file delta adds configurable whole-record diagnostic retention and explicit full versus dialogue-only saved-history modes. Diagnostics are distinct from canonical recovery history; a dialogue projection is not an event log suitable for full replay. Context and any necessary post-snapshot tail preserve restart state, not discarded historical steps.

The student and independent reviewer completed focused checks and compile-valid behavioral negatives followed by exact restoration. Root verified the sealed reports, source identities, outputs and durable final tails before joining and closing both code workers. There were zero formal code correction rounds. The student’s unintended development lifetime run, lost wrapper timing footers, and setup failures remain qualified in its report; the original grades, operational exercises and old binary evidence are not new results for this delta.

The post-student author recommended no further prose change. One fresh formal text review passed with no requested revision. Root verified its three-member seal, 21 inputs, 34 protected tracked files and final 165, joined the turn, and closed the author and reviewer. Lint was carried through unchanged chapter and lint-source hashes; it was not rerun for cosmetic freshness. Soft warnings remain disclosed.

## Evidence

- [Accepted amendment](accepted.json), with the earlier [code-only boundary](accepted-code.json) preserved.
- [Student report](../ch10-retention-student.md) and [parent verification](../ch10-retention-student-parent-verification.json).
- [Independent code completion](../ch10-retention-code-review-completion.md) and [parent verification](../ch10-retention-code-review-parent-verification.json). The initial incomplete review is retained, not relabeled.
- [Post-student author](post-student-author/report.md), [formal text review](text-review/report.md), and [text parent verification](text-review-parent-verification.json).
- [Chapter 14 pending obligation](../ch14-retention-entry-obligation.md).

## Limits and next construction step

This implements the student book exercise; it does **not** install configurable retention into the running Edition 2 orchestrator. The earlier binary rebuild preceded this feature. `micro_handoff` integration is deferred to Chapter 14, where that tool is introduced. No handoff, browser, provider, or power-loss guarantee is inferred from these focused retention checks.

The unfinished Chapter 11 source remains preserved in its original dirty worktree and verified backup. The same Chapter 11 student must integrate the accepted retention predecessor without losing MCP work, demonstrate the resulting construction and focused compatibility, then finish its previously unfinished gates. No Chapter 11 acceptance is implied here.

Main-repository closure includes only this manuscript amendment, its text evidence and the pending Chapter 14 note. Student source lives in its separate repository/tag; local executables, archives, raw captures and profiles are not portable Git evidence. Existing sealed failures are not rewritten. Concurrent Chapter 11 grading/guidance, runtime, publication and other work remains outside this commit. No push, publication build or live restart is part of closure.
