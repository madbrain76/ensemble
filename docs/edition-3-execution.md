# Edition 3: execute, then improve

Bill's instructions of October 16, 2026 govern this run. This is execution and improvement of the existing book, not new-chapter design. They supersede contrary historical phase instructions and the coordinator's proposal to revise each chapter before its student attempt.

## Preparation

1. Copy current application code to `solutions/edition-3/` and current book to `book/edition-3/`. Completed from source commit `4309d5688992f64067f4f75e1cef6bd2e9e25e6a`; `edition-3-baseline.json` records exact copied bytes. The code copy is a reference baseline, not evidence of fresh student completion. Copied publication artifacts remain Edition 2 artifacts until deliberately rebuilt.
2. Cadence reads the entire current book and crossover notes and writes a sourced Edition 3 carryover. Record reading coverage; recover truncated output. Do not infer future requirements from editorial imagination.
3. Cadence makes a light edit pass incorporating the carryover. No reviewer for this preparatory pass. Preserve Bill's preface verbatim. Keep original book, source, graders, frozen answers and historical evidence unchanged.
4. Once prepared, execute chapter by chapter under the existing execution safeguards, with the role/timing corrections below. Student working location and chapter checkpoints must be explicit before launch; the full copied answer is not the student's predecessor.

## Chapter loop and roles

- **Orchestrator:** prepares permitted inputs, supervises workers, answers narrow contract questions from authoritative sources, retains actual evidence and assistance, and respects scope/architecture stop gates. Skills hold procedure; handoffs hold position, reasons, blockers and next action. No new goals subsystem is required.
- **Coder/student:** executes the existing chapter against its accepted predecessor, runs grading and required tests/live exercises, and writes the chapter review. KISS: implement only the actual requirement, with the simplest adequate design. No speculative abstractions, future features or machinery to satisfy invented promises. Preserve explanations of non-obvious decisions. Do not inspect old answers or grader internals. Only the agent doing the exercise supplies experience-based clarification requests; the author does not simulate a student's difficulties.
- **Code reviewer:** independently reviews the result and student's chapter review under the existing code-review rules, including KISS, architecture, test sensitivity, scope and bounded correction rounds. Findings about code do not license new chapter requirements. Preserve failures and explicitly recorded exceptions; never relabel them as passing tests.
- **Author:** begins only after the student has finished and grading/code acceptance is settled. Improve the existing chapter solely from the student's chapter review and verified clarifications. Do not invent requirements, enlarge the exercise, redesign the chapter, add anticipated student difficulties or expand text to demonstrate contribution. A chapter needing no clarification may remain unchanged.
- **Text reviewer:** makes one pass on that bounded edit: accuracy, faithfulness to the student's findings, no new requirements, and load-bearing prose. Not a second designer or a source of speculative exercise scope.
- **Author:** makes one revision in response. Do not start an unbounded prose-review loop. Escalate any material unresolved issue rather than silently add rounds or requirements.

## Load-bearing writing

Every sentence must do useful work: teach a mechanism, explain a reason, remove a demonstrated ambiguity, state a necessary contract, or support the book's grounded narrative. Words for their own sake, repeated caveats, review bookkeeping in the chapter, speculative completeness and arbitrary length targets are not goals. Keep useful stories, voice and rationale; KISS is not a demand to erase explanations. Apply this standard to both author and text reviewer, and simplicity to the coder and code reviewer.

Deletion or simplification must not introduce ambiguity, weaken a requirement, remove a necessary condition or obscure the reason for a design. For a proposed cut, identify what the passage currently contributes and where that meaning remains. If preserving the meaning is uncertain, leave the passage alone. Repetition between the TL;DR, explanation and exercise is not automatically padding.

After the initial carryover pass, Bill authorized one read-only pre-execution pruning pilot on Chapter 11. This is not the post-student author loop or permission for new requirements, a whole-book rewrite or automatic rollout. The reviewer proposes exact low-risk cuts and records why they preserve meaning; the orchestrator checks them in context. Bill's opening is excluded.

## Existing safeguards retained

Use the governing workflow and mandatory coding skill for architecture, honest fresh-student inputs, public-API tests with faithful boundary fakes, mutation evidence, required live exercises, independent code review and bounded correction rounds. Full-document reading requirements remain; handoff summaries are not substitutes. Historical per-edition path/phase instructions are not current authorization. An exception is recorded as an exception, not as passing grading; surface conflicts with Bill's current acceptance instruction before relying on an old waiver.

No push, publication, live installation, restart or destructive repository operation is authorized by preparation. Preserve concurrent changes. No Edition 3 students or authors have been launched. The single authorized Chapter 11 pruning reviewer is a preparation pilot, not chapter execution or acceptance.

## Current position

Baseline copy, full-book/crossover reading, sourced carryover and light manuscript edit are complete. Mechanical verification is recorded separately in `edition-3-preparation.json`. The Chapter 11 pruning pilot recommended zero cuts and initially left the chapter unchanged; its findings remain in `edition-3-notes/ch11-pruning-review.md`. Bill subsequently authorized trimming the chapter and clarifying on-demand browser observations. That edit and its checks are recorded in `edition-3-notes/ch11-trim.md`. No grader changes or Edition 3 execution yet. The five runtime skills have not been installed or enabled; this record captures the agreed procedure, not a claim of skill availability. Next: package the role instructions, establish the fresh student/predecessor location and grader invocation, then execute Chapter 1. Do not hand the copied final implementation to the student.

## High-growth chapter review follow-up

The comparative review and bounded edits to Chapters 6, 7, 8, 9, 10, 16, 22 and 26 are complete; Chapter 16 was left unchanged. Results, exact revisions and checks are recorded in `edition-3-notes/growth-review.md` and its linked JSON. Bill explicitly added a qualitative GUI design goal: attempt beauty, strong contrast and intentional tabs without a prescribed template or automated beauty score. All reviewers are closed. This does not change the remaining execution prerequisites or claim grader/live-GUI acceptance.
