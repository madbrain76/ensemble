# Chapter 4 restart — formal text review

**PASS. No required revisions.** This is the single formal text review of the held author pass after root accepted the 42-file Chapter 4 snapshot and two code-correction rounds. Root alone owns acceptance and commit; no second review or author response is needed.

## Scope and checks

Read the full current `AGENTS.md`, coding/text-reviewer/author roles, shared execution procedure, `book/voice.md`, `book/chapter-writing-procedure.md`, current Chapter 4, and the restart student/code-review/author notes. Reviewed the actual chapter-only diff: **+17/−4 lines**, confined to §§4.3–4.4. The two historical workflow reads were explicitly waived. No plan gate was reopened.

- **Failed output setup:** the new paragraph correctly preserves `tool_called`/`tool_returned` correlation and prohibits running the handler without its output owner. Accepted `solutions/edition-3/ch04/internal/jobs/jobs.go` (`Start`) stops before handler launch on setup failure; `internal/tool/dispatch.go` (`dispatch`) returns the correlated error rather than dropping the call. The prose teaches the observable obligation, not a new storage or orchestration design.
- **Wake versus return:** the revision correctly distinguishes notification/rechecking from returning an observation. Accepted `InitialReport` and `Wait` share the locked readiness loop but differ on no-pattern output readiness: later waits return for unseen bytes; initial reports retain the completion/delay window. Pattern readiness remains restricted to unseen bytes, with completion or delay also able to return. The permissive “may keep” wording preserves fast-call behavior without imposing an extra delay requirement.
- **Killed observations:** the added sentences correctly distinguish truthful killed status/reason from failure of the observing tool. Accepted `Kill` clears completion metadata without setting failure merely because of the kill; `observeLocked` still carries genuine failure state, and lookup/wait errors remain errors. The chapter does not promise successful completion or an invented exit code.

## Voice, preservation, and scope

The small additions stay in the chapter's concrete teaching voice. Their rationale is explicit: no unowned effects, no sleeping through useful later output, no needless conversion of a fast call into a running handle, and no conflation of cancellation with observation failure. The lock/channel explanation, PTY and cursor rules, usage/API examples, first-person experience, dated experiments, historical results, validation teaching, and closing motivation for Chapter 5 remain intact. Nothing imports Chapter 5 behavior or adds tools, dependencies, release gates, or a new architecture. The accepted implementation omissions remain implementation omissions; the edit does not recast the correction rounds as broad manuscript faults.

Inspected the author's saved same-basename baseline/final lint outputs under `/tmp/e3-ch04-restart-author.vfF5SX/`: unchanged soft warnings (29 negation forms and the 5482-word person gap), empty stderr, and author-reported exit 0 for both runs using current `cmd/lintprose`. No lint rerun or warning-cleanup expansion was needed.

This review ran no builds, tests, graders, runtime/live/browser activity, or children and did not access keys. Only this review note was written; the manuscript, accepted code, other notes, and unrelated staged memory files were left untouched. No manifest, seal, or audit package was created.
