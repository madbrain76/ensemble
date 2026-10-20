# Chapter 5 restart — independent text review

**Verdict: PASS. No required revision or author response.**

One formal review of the current Chapter 5 and actual author diff, under the current governing instructions and root’s bounded historical-archive waiver. The full required instruction/craft/chapter reads and concise student, code-review and author notes were completed; clipped ranges were recovered. Accepted source was consulted only for the changed contracts and layout.

- **Reliable completion:** §5.1 correctly distinguishes admission to the reliable mailbox from delivery into a private one-result reply buffer. The buffer keeps an abandoned caller from blocking the actor; it cannot rescue a request silently dropped at admission. This repairs an incomplete explanation without prescribing a new API or queue implementation. The accepted snapshot’s `internal/agent/actor.go` supplies the matching lossless admission and buffered reply paths.
- **Kill versus Shutdown:** §5.5 preserves immediate containment of explicit Go-only Kill without pretending Go can terminate the handler. Shutdown separately joins real work and terminal-fact delivery before saving. `internal/jobs/jobs.go` and the actor’s killed-report acknowledgment support that distinction, including the external-Kill reporting race. No requirement to wait for a Go handler during explicit Kill was introduced.
- **Build layout:** the live CLI command now targets `./cmd/ch01` from the core module root. The accepted snapshot has `package ensemble` at that root and `package main` in that command directory. The separate workflow module’s build instructions remain appropriate; no library flattening is requested.

These are load-bearing, local clarifications. The opening, explanation, ownership rationale, existing grader block and historical demonstrations are unchanged. Most accepted student findings concerned requirements already explicit in the chapter; turning them into more warnings or new requirements would add noise. No provider-specific hint/memory behavior or cause for Gemini’s empty response was invented. The restart notes’ qualified test/live results and narrow acceptance exception remain qualifications, not upgraded claims in the chapter.

Inspected both saved lint reports in `/tmp/e3-ch05-author.QwUkSz/`: both report zero errors and the same pre-existing warning set, with later line numbers shifted by the additions. No lint rerun, tests, builds, providers or new coverage campaign were performed.

**HOLD:** this completes the single formal text-review pass. Only this note was written. No second review is requested; chapter acceptance and any subsequent release remain with the parent/root.
