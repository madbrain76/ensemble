# Chapter 1 restart — author draft

Addressed the student’s demonstrated 70 KiB input failure in the Exercise, immediately after “Grader mode is the default.” The current manuscript had no Scanner snippet to replace, so the change adds one compact `bufio.Reader.ReadString` loop sketch and its explanation. It identifies Scanner’s implicit roughly 64 KiB limit, processes a final unterminated line once, and returns on non-EOF read errors before cumulative protocol usage. Errors go to stderr with a nonzero exit; terminal usage remains successful-EOF output. The accepted `ch01` CLI and README informed the example.

No clarification elsewhere was needed: this was a CLI line-reading issue, not an ownership or lifetime gap. All existing teaching, history, requirements and other passages remain unchanged. No source changes or runtime tests were made.

Normal prose lint (`go run ./cmd/lintprose`) ran on the unchanged baseline at `/tmp/ensemble-ch01-author.6afYy6/chapter-01.md` and on final `book/edition-3/chapter-01.md`, preserving the same basename. Both exited 0 with all hard ceilings met. The sole warning was unchanged: person gap after L110, 1611 words versus 1200. Chapter whitespace check passed. Draft ready for the one scheduled formal text review.
