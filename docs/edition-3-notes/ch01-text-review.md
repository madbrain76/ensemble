# Edition 3 Chapter 1 — independent text review

## Verdict

**Accepted as written for this bounded text-review pass. Zero required findings; no author revision requested.** This is the one formal independent TEXT review after code acceptance and initial author clarification. It does not reopen code acceptance, authorize another review loop, or authorize the next chapter.

The clarification fixes a demonstrated invocation problem without changing the Chapter 1 implementation contract. It preserves the useful directory example, states its supported layouts, and gives a supported executable-target alternative for a nested CLI. No speculative API sketch, duplicate state/usage exposition, new feature, or architecture requirement is needed.

## Frozen inputs and scope

- Book/grader baseline HEAD: `7f3d069420112c2136bc95d1d551e12fad6964c1`.
- Full current `book/edition-3/chapter-01.md`: SHA256 `d12183d54265fdd849293946878b0b8d7131830de2aa702928419795ecb34cd7`.
- Full [initial author record](ch01-author-review.md): SHA256 `cea3102119af1a22d060984e27eb8cdf8ae4641dfa3a6c15b3b2c2776268c82f`.
- Accepted student context, supplied by the handoff and recorded in the accepted correction review: `1aad3975842d502b9144b7b9bddfd30ca7309db7`, tag `edition-3-ch01-r1`. No new student source review was performed.

Both frozen document hashes matched before review and again after all substantive review checks. The exact chapter diff against the baseline is one hunk, **19 added / 3 removed lines**, whole-Markdown whitespace-split words **5,464 → 5,555 (+91)**. Everything before “Grade yourself” and after the edited passage is byte-identical to baseline. The author's separate +74 prose-word count is recorded in the author report; I did not rerun that linter/counting method.

Prerequisite reading was completed in the preceding readiness-only turn: `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`, `.agents/skills/edition-3-text-reviewer/SKILL.md`, `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, `docs/edition-3-execution.md`, and `docs/edition-3-carryover.md`. The current Edition 3 guide and text-reviewer overlay were reread for this pass; current instructions supersede historical phase/path guidance.

## Assessment and evidence

### Technical accuracy and command plausibility

The changed passage is `book/edition-3/chapter-01.md:631–657`.

- **Directory layouts (`:631–638`):** `internal/grade/harness.go:102–128` starts with the root build target, inspects only direct `./cmd/`, and selects it when its package name is `main`. It does not recursively discover `cmd/ch01`. The qualification is accurate and does not prescribe a root-main architecture.
- **Executable alternative (`:640–652`):** the example builds `./cmd/ch01` from the solution's Go module root, uses an absolute output path, explicitly returns to the book repository root, and passes that executable as the positional target. `cmd/grade-edition3/main.go:20–28` accepts `-ch`, `-json`, and exactly one positional target; `internal/grade/harness.go:88–95` resolves a file target to an absolute path rather than rebuilding it. The instruction to substitute actual paths and the CLI package makes the placeholders clear. The archive explanation accurately identifies why building a library root is not equivalent to building its CLI.
- **No contradictory gate:** the existing seven conversation checks, two lifetime gates, and requirement for `passed: true` remain unchanged. The edit changes invocation guidance, not scoring, implementation scope, or accepted behavior.

These are **source-inspection conclusions**, not a new executed grader pass. A scoped baseline diff confirmed that the two inspected grader source files are unchanged from the supplied book/grader revision.

### Demonstrated usefulness, voice, and KISS

Read in full: [student review](ch01-student-review.md), [student correction](ch01-student-correction-1.md), [initial code review](ch01-code-review-r1.md), and [accepted correction review](ch01-code-review-correction-1.md), as well as the chapter and author record.

The student's first proposed teaching clarification is supported by retained negative and positive evidence, not an imagined learner problem:

- `ch01-student-evidence/grader-initial.stderr` and `.exit` record the directory-target `permission denied` failure, child status 2 / Go wrapper exit 1.
- `ch01-code-review-evidence/directory-grader.stdout`, `.stderr`, and `.exit` retain the independent reproduction of that failure.
- Both evidence directories' `README.md`, `grader.json`, `grader.stderr`, and `grader.exit` document the explicit executable-target invocation and its successful result: exit 0, `passed: true`, 100/100, all seven scored checks and both lifetime gates true.

Those passing runs belong to the initial student candidate `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`. The accepted correction review establishes the later test-only correction's production identity. I inspected the saved invocation evidence; I did not rerun it or relabel it as fresh verification at the accepted commit.

The new prose is short, direct, and placed beside the command that caused the actual difficulty. “You do not need to move your main package” prevents the build clarification from becoming an architecture mandate. The chapter's narrative and teaching explanations remain intact. No repair-reference story, internal orchestration mistake, malformed-thinking-fixture requirement, or speculative future feature was added to the manuscript. The author's declined state table, facade sketch, and duplicate usage/anecdote labels do not leave a demonstrated material gap in the full chapter.

### Checks actually run in this text review

Read-only checks against the frozen chapter:

- SHA256 checks before review and after substantive checks: both frozen hashes match.
- Exact baseline diff and independent word/hunk/line counts: results above.
- `bash -n` on each of the two grading shell examples: exit 0, no diagnostics. This checks shell syntax, not execution of placeholder paths.
- Fence scan: all **24** fence markers pair correctly.
- Local Markdown link check: `../../solutions/evidence/ch01/live.txt` resolves. No link was added by the edit.
- Trailing-whitespace scan: none.
- `git --no-pager diff --check 7f3d069420112c2136bc95d1d551e12fad6964c1 -- book/edition-3/chapter-01.md`: exit 0, no diagnostics.

## Limits and handoff

No new build, grader/lifetime/race run, live-provider exercise, credential access, external-link request, or publication/render build was performed. The author's build, EOF, parser/path, and prose-linter checks remain attributed to the author record, not to this reviewer. Source inspection, saved execution evidence, and the checks above are sufficient for this narrow invocation clarification; they are not a fresh code-quality verdict.

Only this new report was written. No manuscript, author record, student source, grader, application, preface/frontmatter, or other chapter was edited. No delegation, runtime changes/restart, staging, commit, or broad Git operation was performed. Unrelated runtime work was not investigated.

**Final disposition: zero findings; retain the frozen chapter as written.**
