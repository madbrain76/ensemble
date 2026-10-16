# Edition 3 Chapter 1 — initial post-student author review

## Scope and decision

Code acceptance was settled before this pass: student commit `1aad3975842d502b9144b7b9bddfd30ca7309db7`, also resolved by immutable tag `edition-3-ch01-r1`, with R1 closed and no exception. Book/grader HEAD is `7f3d069420112c2136bc95d1d551e12fad6964c1`. Complete current governing documents, chapter and Chapter 1 records were read; oversized reads were recovered in smaller sections. The current Edition 3 execution guide governs this bounded pass.

**Accepted feedback:** [student review](ch01-student-review.md), suggestion 1, supported by the [initial review](ch01-code-review-r1.md), [accepted correction review](ch01-code-review-correction-1.md) and [execution record](ch01-execution.md). A valid public-root library with `cmd/ch01` hit a permission error because directory building produced an archive. The existing grader accepts an executable. Only “Grade yourself” changes: retain and qualify the useful directory example; show a nested CLI build from the student's module root followed by an explicit positional executable target from the book repository root. No root-main layout, `-impl` flag, architecture change or grader change is prescribed.

**Not adopted:**

- Suggestion 2, a state table: §1.4 already explains success-only history commits, measured usage on empty text, and missing accounting versus zero. The student explicitly calls the prose clear; their initial decoder omission does not establish another teaching gap.
- Suggestion 3, a public-root API sketch: the ownership diagram and existing exercise already state the boundary. A concrete facade would unnecessarily prescribe a valid implementation choice.
- Suggestion 4, usage/anecdote labels: §1.5 already distinguishes per-round from cumulative counts; the live exercise requires accumulated usage. §1.6 identifies the preserved transcript and explicitly qualifies the prompt anecdote. Preserve these rather than duplicate them.
- The malformed thinking fixture was test fidelity, corrected in [the student correction](ch01-student-correction-1.md), not a missing manuscript requirement. No thinking support or new malformed-input obligation is added. The parent's wrong credential reference and pause/status delay remain orchestration issues, not chapter requirements.

## Exact change

One hunk, 19 added / 3 removed lines. Whitespace-split whole-Markdown words: **5,464 → 5,555 (+91)**. Repository prose-linter words (excluding headings/code/tables): **4,664 → 4,738 (+74)**. Exercise, implementation and all other chapter prose remain unchanged.

Chapter SHA256 before: `ab016baaff3f6a9d83ecffcbfc3772cd692c9d33c06d9fa1c52ad9a8708d25df`.

Chapter SHA256 after: `d12183d54265fdd849293946878b0b8d7131830de2aa702928419795ecb34cd7`.

```diff
diff --git a/book/edition-3/chapter-01.md b/book/edition-3/chapter-01.md
index f3af60ac..79157d29 100644
--- a/book/edition-3/chapter-01.md
+++ b/book/edition-3/chapter-01.md
@@ -628,15 +628,31 @@ paying in Chapter 3, and that is where it gets graded.
 
 ### Grade yourself
 
-Free, as often as you like:
+Free, as often as you like. Run the grader from this book repository's root.
+For a solution with a main package at its root or directly in `cmd/`, pass the
+solution directory:
 
 ```bash
 go run ./cmd/grade-edition3 -ch 1 path/to/your/solution
 # Add -json before the solution path for machine output.
 ```
 
-Run this from the repository root. Edition 3 retains the seven conversation
-checks and adds two required, zero-point lifetime gates: real 62-second delays
+Directory builds do not discover nested commands such as `cmd/ch01`. For that
+layout, build the CLI from your solution's Go module root, then pass the
+executable instead:
+
+```bash
+cd /absolute/path/to/your/solution
+go build -o /absolute/path/to/ch01 ./cmd/ch01
+cd /absolute/path/to/this/book/repository
+go run ./cmd/grade-edition3 -ch 1 /absolute/path/to/ch01
+```
+
+Substitute your actual paths and CLI package. Building a library root produces
+an archive, not a runnable CLI; you do not need to move your main package.
+
+Edition 3 retains the seven conversation checks and adds two required,
+zero-point lifetime gates: real 62-second delays
 before response headers and during the response body, taking about two minutes
 in total. Require `passed: true`; a score of 100 alone is not a pass.
 
```

## Validation and limits

- **Source inspection:** `cmd/grade-edition3/main.go` accepts `-ch`, `-json` and one explicit positional target. `internal/grade/harness.go:Build` returns an absolute path for a file target; for directories it checks only a root or direct `cmd/*.go` main package, not `cmd/ch01`. Inspection is not an executed grader pass.
- **Executed:** exported the accepted commit with `git archive` to disposable scratch; from that module root, `go build -o "$SCRATCH/ch01" ./cmd/ch01` exited 0. `file` identified a Mach-O executable. Empty-stdin execution with dummy configuration and a loopback base URL exited 0 and printed `{"usage":{"input":0,"output":0}}`; no request was made. Scratch was removed.
- **Executed from book repository root:** `go run ./cmd/grade-edition3 -h` exited 0. `go run ./cmd/grade-edition3 -ch 1 -json "$SCRATCH/missing-ch01"` reached the positional-target stat check and failed as expected: missing file, child status 2 / Go wrapper exit 1. This is a parser/path probe, not passing grading. The full documented invocation against the existing executable is supported by inspection and the retained passing executable-target reports, not a new full-suite run here.
- **Honest author-check failure:** the first scratch validation asserted provider usage keys (`input_tokens`/`output_tokens`) instead of the documented CLI keys (`input`/`output`). Build and EOF execution had succeeded; only that assertion failed. Correcting the validation assertion made the rerun pass. The final record whitespace check also initially flagged literal Git diff context lines; the corrected check exempts that exact-diff block rather than changing it. No student or manuscript change resulted.
- **Prose/mechanics:** repository `lintprose -v` passed both original and edited Chapter 1, exit 0, with the same existing 1,611-word person-gap warning after line 110. `git diff --check` passed. All 24 chapter fence markers pair correctly; the sole local Markdown link resolves; no trailing whitespace. No external-link requests or publication build.
- **Protected scope:** the before/after fingerprint of all 10,160 tracked and nonignored untracked files outside the two allowed paths matches: `4d075ff0c32d2a96c7121dcf8e1494156cc9a455e705403a766ac4556fc2cd45` (sorted path plus file SHA256 manifest). Index unchanged; accepted student tree remains clean at the same commit/tag. Existing unrelated artifacts were preserved. Grader/reference source, preface/frontmatter, other chapters and application/student code were not edited.

No full grading/race/lifetime rerun, provider call, credential access, runtime change/restart, staging, commit, push, publication or delegation. Prior grading/live results retain their original attribution to the initial candidate and are applicable through the accepted test-only correction's verified production identity; this pass does not relabel them as fresh executions.

This is the **initial author pass only**, ready for one independent text-review pass. At most one later author revision may be assigned; no extra loop or next-chapter work was started.
