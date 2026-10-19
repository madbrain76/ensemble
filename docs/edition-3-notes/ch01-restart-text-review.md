# Chapter 1 restart — formal text review

**PASS. No required revision.** This is the single formal text review.

Read the governing instructions and author/reviewer guidance, the full current
`book/edition-3/chapter-01.md`, its actual git diff, and both restart notes.
The diff adds only 26 lines after “Grader mode is the default”; nothing is
removed. Consulted only the permitted accepted `ch01` CLI and README for
implementation behavior, without reopening code acceptance.

- The addition addresses the student's reported 70 KiB well-formed-line failure:
  default `bufio.Scanner` imposes roughly a 64 KiB limit, whereas the illustrated
  `bufio.Reader.ReadString` loop accommodates longer physical lines. The baseline
  chapter contained no Scanner example. This is new explanation, not correction
  of an old manuscript example; the author note correctly distinguishes that
  from the mistaken attribution in the root orchestrator's assignment.
- The EOF ordering is accurate: process a nonempty final unterminated line once,
  then stop on empty EOF. Non-EOF read errors return before terminal protocol
  usage; the surrounding prose correctly requires stderr and nonzero exit.
  Successful protocol EOF emits cumulative usage, distinct from chat's per-turn
  usage. This agrees with the permitted CLI and README.
- The fragment is explicitly a loop sketch inside an error-returning CLI
  function, not a standalone program. Its standard-library `bufio`, `io`, and
  `errors` qualifiers are understandable; line processing and the enclosing
  input/error plumbing are intentionally left to the reader.
- Whole-chapter reading found no contradiction introduced by the addition.
  Existing voice, history, preface references, public contract, star ownership
  and KISS boundaries remain intact. The sketch introduces no library API or
  hidden acceptance requirement; equivalent contract-satisfying implementations
  need not copy its private names or exact structure. No broader rewrite,
  features, exercises, benchmarks or process diary are needed.

Per `ch01-restart-author.md`, baseline and final normal prose lint both exited
0 with all hard ceilings met. The same soft warning remained: person gap after
L110, 1611 words versus 1200. These are attributed author results, not new runs.
No tests, builds, lint replay, provider calls or browser work were performed;
no source or manuscript was edited. Only this review note was written.
