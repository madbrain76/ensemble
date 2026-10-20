# Chapter 5 restart author pass

Status: **HOLD** for the root-owned formal text review.

Read the full current `AGENTS.md`, coding and Edition 3 author skills,
`docs/edition-3-execution.md`, `book/voice.md`, chapter-writing procedure,
Chapter 5, and actual restart student/code-review notes. Recovered clipped
output in bounded reads. The root's historical-archive waiver governed this
pass. Inspected relevant accepted mailbox/job/shutdown code, library/CLI
layout and `cmd/grade-edition3/main.go`; did not change them.

## Bounded changes

- Corrected the manual build target from `.` to `./cmd/ch01`, with allowance
  for the reader's own command path. The accepted module root is an
  importable library, not the CLI.
- Distinguished reliable mailbox admission from buffered reply delivery.
- Made explicit that Go-only Kill returns without joining its handler,
  whereas Shutdown joins workers and terminal-fact delivery before saving.
  Handler exit alone does not establish that an admitted kill was recorded.

Existing teaching already covered most restart findings. These clarifications
are not a claim that missed explicit requirements were manuscript ambiguity.
The grader block already describes Chapter 5's core/consumer/workflow/media
inputs followed by Chapter 1 on the same CLI; it remains unchanged. Opening,
history and personality remain intact. No new API, framework, acceptance gate,
provider recipe or harness requirement was added. No Gemini cause, model
superiority or runtime-hint claim was added.

## Checks and limits

Baseline and final copies both used basename `chapter-05.md` under
`/tmp/e3-ch05-author.QwUkSz/{baseline,final}/`. Both runs of
`go run ./cmd/lintprose -v <copy>` exited 0. Hard failures: none. Both retained
38 negation forms and a 4,217-word person gap as soft warnings; short closers
remained 35 and pointer closers 1. No style cleanup was undertaken for these
pre-existing warnings. The final bounded diff was inspected.

Only Chapter 5 and this note were edited. No code tests, graders, provider,
browser/native calls, publication build, children or commit were run. Accepted
code/operations remain the root's evidence, not a new author verification.
One formal text review remains with root; at most one author response is
reserved, with no rereview.
