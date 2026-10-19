# Chapter 1 restart — independent code review

## Verdict and comparison

**Recommend acceptance of the held `solutions/edition-3/main` candidate.** One formal correction round closed the only material source finding; no unresolved correctness, ownership, or material-bloat finding remains. Root retains final acceptance/snapshot authority.

Reviewed the current public Chapter 1 contract and governing guidance, actual student production source, original `solutions/ch01` source, then changed areas and final test source/results. The original remained read-only and was not supplied to the student.

The student preserves the mandatory star: the public root imports and constructs both implementation spokes; Agent and Engine depend on shared interfaces rather than each other. Root owns its logger/Agents, Agent owns configuration/history, and Engine owns transport/measured usage. Stored parent interfaces identify the actual owners and carry configuration/diagnostic access. Shared vocabulary is not a state registry. The synchronous core is reused by both CLI modes.

Relative to the original, the modest extra code buys useful ownership comments/accessors, validation, sanitized transport errors, redirect refusal, and caller cancellation—not a framework or speculative feature. The original is compact and clear but its `internal/engine/engine.go:34` installs a **60-second total HTTP timeout**, forbidden by the current chapter. The student's timeout-free client correctly avoids copying that defect. No reward is given for fewer lines at the expense of current requirements.

One HTTP exchange per question, full retained history, ordered text-block concatenation, provider-derived cumulative usage, and commit-only-after-success history are supported by source and tests. Measured usage survives an unusable-text response; missing accounting is not silently treated as zero. Default JSON-lines output and opt-in chat remain distinct. Comments explain ownership, failure boundaries, and request lifetime without excessive machinery.

## Correction closure and reviewer retraction

Round 1 replaced default `bufio.Scanner` with `Reader.ReadString` (`cmd/ch01/main.go:46–62`), removing the undocumented approximately 64 KiB physical-line ceiling. The reader handles CRLF, an unterminated final line, normal EOF, and genuine read errors. The small nil-root guard (`ensemble.go:46–49`) is covered too.

My first early handoff contained inaccurate source citations and claims about constructor placement, a request counter, lazy configuration, and an error-JSON requirement. Those claims were explicitly retracted after rereading the absolute paths; they are **not findings or acceptance evidence**. Root constructs both children; no application request counter exists; configuration validation is eager; the chapter does not require an error envelope. The scanner finding was subsequently reproduced, not merely inferred.

## Actual validation and qualifications

I read the following receipts under `/tmp/ensemble-ch01-restart.bj87Gl`; I did not rerun tests/builds/grading or access credentials/providers.

- `large-line-before.log`: genuine well-formed 70 KiB input failed with zero HTTP calls. `large-line-after.log`: regression passed after correction.
- `race.log`: root and CLI tests passed. Reviewed tests exercise real local HTTP, parent/config/logger isolation, history copies and replay, failure accounting, cancellation during headers/body, redirects, sanitized errors, and CLI framing. Student/root report successful `go vet ./...` and final build exits; `vet.log` and `final-build.log` themselves are empty.
- `grade.json`: **100/100, passed=true**, seven scored checks and both zero-point real 62-second header/body gates passed. `grade-status.log`: exit 0, elapsed 125 seconds. These are existing results, not reviewer executions.
- `focused-tests-peer-hang.log` preserves an earlier cancellation fake-peer teardown hang. The corrected peer drains the request and has a cleanup release; `cancellation-fixed.log` passes. This was test setup, not a demonstrated product cancellation failure.
- Final student note and `live-run.log`/`live-chat.stdout`/`live-chat.stderr`: one discovery GET selected `claude-sonnet-5-5`; three Messages POSTs returned HTTP 200. Replies recalled Cedar/version 7 and sum 12. Usage 92/23 + 141/80 + 212/45 matches CLI totals **445/148**. CLI/proxy exited 0 in 4.63 seconds. Root verified the proxy and unchanged graded binary/native 1024 output budget; I read the sanitized receipts, not credentials. The earlier root-supplied wrong credential field failed before networking and is preserved as a setup error, not erased or charged as a model request.

The final student report records successful build/vet exits and scratch mutation work; neither was independently rerun here. Root reports production held after the final build, with only README/test comments/report changes afterward.

## Limits

Acceptance is for the chapter's synchronous, non-streaming conversation, not concurrent calls or later chapters. The three-turn live exercise is a bounded smoke test, not a general model-quality guarantee. The 62-second gates alone cannot exclude longer hidden deadlines; direct source review supplies that distinction. No further provider requests, mutation campaign, or review reports were needed.

## Root acceptance

Accepted after direct source/reference comparison and inspection of the reported local and live results. One code-correction round; no material bloat requiring re-execution. The nine source/documentation files are copied unchanged to `solutions/edition-3/ch01`; `main` is the evolving workspace. Text closure remains pending.
