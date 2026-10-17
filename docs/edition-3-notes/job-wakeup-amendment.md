# Uniform jobs and prompt foreground wake-up: authorized amendment

This is a teaching-contract amendment to Chapters 4 and 5, explicitly requested during Chapter 5 execution. It is not a runtime installation, a passing implementation result, or a rewrite of historical Chapter 4 acceptance.

Chapter 4 clarifies the existing ordinary-job path, common output and callback limits, literal-next-call `tool_limits`, and controls that operate on existing jobs without allocating another job. Chapter 5 requires hints and addressed questions to release foreground observation promptly while work continues. A completed call retains its real result; a started job retains its real state and handle; an undispatched batch-tail call gets its own interrupted/not-executed result, no effects or new job, and no automatic replay. New calls may run while the original job continues, allowing the parent to answer a student before resuming the same join job. Limits cannot leak from discarded calls, and the original assistant call records remain intact.

The chapter includes Bill's observed parent/student waiting cycle and the reported `message_fleet` attempt, with motive and incident-trace limits explicit. The exercise requires real public-interface local regressions for initial and later waits, message delivery, retained output, timer/completion races, and the undispatched batch tail. Historical grader scores and paid demonstrations do not establish these new behaviors.

## Review and validation

- [Author report](job-wakeup-amendment-author.md), [exact final diff](job-wakeup-amendment-evidence/amendment.diff), and [queued-tail clarification](job-wakeup-amendment-evidence/queued-tail-clarification.diff). The direct user clarification arrived before the formal verdict; the superseded draft is preserved.
- [Independent text review](job-wakeup-amendment-review.md): PASS, one formal review, no reviewer-requested author revision. Author final turn 203 and reviewer final turn 308 were verified saved, joined, and both text workers closed.
- Same-basename baseline/final prose lint: exit 0, zero HARD findings; soft warnings and increased word counts are disclosed, not represented as unchanged. The chapter diffs are whitespace-clean.
- Parent verified all 28 authored inventory entries, final reviewer pins and report/index hashes, 231 frozen reference files and 47 original manuscript files unchanged, and all eight local chapter links resolving. No new chapter links were introduced.
- [Validation and artifact hashes](job-wakeup-amendment-validation.json) record the exact scope. Limited credential-marker screening and the normal staged hook are separate checks, not comprehensive security proof.

The existing student's correction candidate `6c4561cbd6cbae080d952da34a5fcd588a76920f` has independently closed its three earlier defects under the earlier contract. It is the starting point for the new feature, not proof of it. No Chapter 5 acceptance tag or new paid calls are authorized by this text amendment. Original source/evidence and live-call attribution remain intact. The separate student implementation and independent code acceptance must follow; nothing was pushed, published, installed, or restarted by this amendment.
