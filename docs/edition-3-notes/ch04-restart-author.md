# Chapter 04 restart author pass

Student-first feedback and accepted Chapter 04 source support three small clarifications, not a broader rewrite. Most findings were implementation omissions of requirements already explicit in the chapter.

- §4.3: the cut-down dispatch sketch now cautions that failed output-file preparation must not start the handler or bypass exact-call-ID error pairing and the remaining response calls.
- §4.4: notifying a sleeping waiter is distinguished from returning an observation. A later no-pattern wait returns on unseen output; the initial no-pattern observation preserves the fast full-result path. Pattern matching remains limited to unseen bytes.
- §4.4: observing `killed` does not by itself fail a wait; genuine job/observation errors remain distinct.

The chapter diff is 17 added / 4 removed lines. History opening, TL;DR, rationale, examples, exercise requirements, and already-explicit PTY/newline/geometry/CRLF/footer material are preserved. No implementation, grader, other chapter, or guidance files changed.

Current inputs were read; root explicitly waived the historical workflow/carryover reads for this lean pass. Factual source checks were limited to accepted `solutions/edition-3/ch04` job observation, dispatch, and Agent call processing. No runtime tests, builds apart from the lint command, provider calls, children, snapshots, or commits were run.

## Lint

Ran current `go run ./cmd/lintprose <absolute chapter path>` once before and once after, offline with the existing Go cache. Both exited 0; both stderr files are empty. Same-basename chapter copies and raw paired outputs remain under:

- `/tmp/e3-ch04-restart-author.vfF5SX/baseline/{chapter-04.md,prose-lint.stdout,prose-lint.stderr}`
- `/tmp/e3-ch04-restart-author.vfF5SX/final/{chapter-04.md,prose-lint.stdout,prose-lint.stderr}`

Pre-existing soft warnings are unchanged: 29 negation forms (budget 17) and a 5,482-word person gap (budget 1,200); only its line reference moved from 375 to 388. No density rewrite was warranted. Chapter diff whitespace check passed. Available disk space remained above 11 GiB; no cleanup occurred.

This is an author pass, not formal review. Root owns the fresh text review and any single subsequent author-response request.
