# Edition 3 Chapter 4 student report

## Scope and identity

- Authorized launch instructions: `/tmp/cadence-ch04-student-task.txt`, SHA256 `f90d785c5b3ba37182e3e0170566d621fc0f79ef40fb7d125fd9343b049eeac8` (verified).
- Public chapter: `/Users/bill/projects/ensemble/book/edition-3/chapter-04.md`, SHA256 `016c4e88af4bcb657aad6733e5f65aa09509b0397e84a37059ca245c588eaa55` (verified).
- Student repository: `/Users/bill/projects/ensemble-edition-3-student`; accepted clean predecessor `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, tag `edition-3-ch03-r1`.
- Required public GUIDE, author/student role, coding rules and whole Chapter 4 read before implementation. Ownership/lifecycle plan sent before edits; parent approved.
- Frozen staged source tree: `840b5b3bab4daf7490b5555451bb6506b273d7b9`. Final commit: **pending**.
- Tested/live core binary: `/tmp/ch04-core`, SHA256 `b0ed87f60783977c2024547690b4865ff6269da0f34a51c484f39901eb003a57`.
- Independent consumer binary: `/tmp/ch04-consumer`, SHA256 `d6e4074774199f5cc2476b5fd648bf14cb60d6325e6c94d09158fd53adaf8dd2`.
- Both binaries were built from this exact frozen dirty/staged tree before commit. Live metadata therefore records the predecessor HEAD plus the frozen source tree, not a false claim that a later commit had already existed. Final commit must have the identical tree. No production source changes during final gates/live exercises.
- Evidence directory: `docs/edition-3-notes/ch04-student-evidence/` in the main book repository. Manuscript, grader internals, historical answers, private reviews and future chapters were not read or modified. No application generation deadline or provider downgrade added.

## Implementation and ownership

Root `ensemble.Ensemble` constructs one `internal/jobs.Table`; each job stores its actual Agent owner. Consumers, tools, jobs and History communicate through public aliases/common contracts; feature leaves do not import one another. No process-global registry or callback side table.

- `internal/agent/dispatch.go`: centralized ordinary dispatch, job creation before tool_called, inline supervision, ephemeral one-call limits, explicit option precedence/validation, unknown-tool failures, metadata and the unchanged explanation-policy boundary. Refusals produce paired feedback but no job or side effect; pending limits are consumed even by refused calls.
- `internal/jobs/jobs.go`: synchronized retained jobs, exclusive persistent `cr/io/<handle>` files (never truncate an earlier run), disk-backed complete output, shared consumable cursor, unseen-only regex matching, bounded head/tail reports, terminal killed state, waiter wakeup and late-result suppression. Process launch/attachment is atomic with kill; concurrent cleanup waits for reaping. Go cancellation is explicitly not a claim to stop a handler goroutine.
- `internal/tools/process.go`: fresh `/bin/sh` per command through a real 200x50 PTY, TERM=dumb, original non-secret environment, provider credential removal, merged output, CRLF normalization, exit_code trailer, newline input, process-group kill and reaping. Observation delay is not an execution timeout. Workspace-relative per-call cwd does not mutate process cwd or later shells.
- `internal/tools/tools.go`: original six tools retained, four supervision definitions added; ordinary local operations use job persistence/caps. Relative file/search paths use the Agent workspace; persisted output locators resolve under data dir.
- `internal/history/history.go`: locking plus independent snapshots for concurrent kill events, structured job metadata and job_killed events; loading history does not reconstruct jobs.
- Root public API: RegisterTool, ToolContext, ToolHandler, Job/JobTable/JobRecord/WaitOptions; root Shutdown. CLI calls Shutdown before Dump, including error paths, so EOF cleanup is persisted.
- `examples/ch04-consumer/`: independently buildable module with local replacement, public-only import, additive shout registration, actual ToolContext→Agent→root logger diagnostic, explicit provider/workspace/data/gate mapping, shutdown-before-dump. Documented `consumer-case.json` and README.
- Added `chapter4_test.go`, `JOBS.md`, PTY dependency/checksums and narrow runtime-output ignores. Updated only Chapter 3 assertions for intentionally replaced six-tool/blocking stdout/stderr behavior and workspace capture timing; provider/explanation tests remain.

## Chronology, raw failures, and corrections

Raw logs are retained under `ch04-student-evidence/raw/`; none is replaced with a cleaned success transcript.

1. Initial scripted edit failed on a nonmatching anchor (no partial Agent write); retained `initial-edit-anchor-failure.txt`. Narrowed against actual source. First compile exposed an overbroad field replacement in AttachTools and missing environment helper; next compile exposed invented event constants not present in the accepted vocabulary. Exact failures: `ch04-first-test.log`, `ch04-second-test.log`, `initial-gofmt-integration-failure.txt`.
2. I prematurely SIGQUIT'd a test process after roughly 30 seconds believing it stalled. Stack proved it was the inherited delayed-header lifetime test, not a job deadlock. Induced failure/stack retained in `ch04-third-test.log`. Subsequent full/race runs were allowed to finish naturally (about 125/126 seconds each).
3. Focused failures included replaced six-tool/JSON stdout-stderr expectations, refusal wording regression, mistaken test API name Lookup vs Get, omitted System/MaxTokens test config, and real non-strict tool_limits decoding. Corrected those narrowly; all intermediate focused logs retained. Added real file/PTY/process and public registration tests, not fake output or shell substitution.
4. My first Chapter 4 invocation incorrectly passed the executable to `-core`. The launch actually required the module directory; parent corrected my misreading. No grading occurred (`ch04-grade4-first.log`). I also tried unsupported Edition3 `-ch 3` for diagnostics; adapter returned usage and graded nothing (`ch04-grade3-diagnostic.log`). No grader internals were read.
5. Corrected public Chapter 4 run scored 150/160: all jobs/consumer checks passed, Chapter3 toolerror parity failed. Restored the requested tool name in unknown-tool feedback; next public run scored 160/160. Same executable Chapter1 scored100/100 including both 62-second lifetime checks. These are intermediate results, not final acceptance.
6. Added final lifecycle refinements: serialize process startup with shutdown; wait for concurrent cleanup; avoid truncating existing persistent output; workspace-correct search; direct registered operations through jobs. Added actual concurrent launch/shutdown and persistence tests. Focused final and vet passed. Earlier full/race also passed but predated these refinements.
7. I initially proposed inaccurate live model IDs and a nonexistent copied-env setup in a parent checklist. Parent corrected it before any credential read or live call. No copied credential file was created or read. Reread/verified full launch and saved corrected `live-checklist.md`, then got parent approval. Some intermediate suites/graders overlapped; final required gate sequence is strictly full→race→Chapter4→Chapter1.
8. Initial Anthropic live prompt was overfull and exhausted the inherited ten-tool-round cap, exit1. It repaired the file and exercised capped persistence and real refusal/retry, but did not finish the planned PTY/debugger or two-turn continuity. Original stdout/stderr/events/metadata retained separately. An OS monitor started after that failed run and observed nothing; `anthropic-eof-process-observation.json` is NOT shutdown proof.
9. OpenAI/Gemini and consumer runs completed. OpenAI's final prose falsely said a continuity printf executed; its neutral log proves the command was refused for missing visible explanation, consumed pending limits, and created no job. This report credits the refusal/consumption only, not a successful command. A short evidence parser first failed on the JSONL version header; corrected parser skips non-event rows (no source behavior changed).
10. Parent approved exactly one bounded Anthropic follow-up with four prompts, unchanged binary/model/round cap: real 8-second Go test with failure→repair→PASS; real HTTP service across turns; real Delve commands; still-running server immediately before EOF with OS process observation and shutdown event. No GDB on this macOS host; no additional provider campaign. **Follow-up result pending.**

## Final sequential validation (pending completion)

Frozen source/binaries above; exact command order recorded with start/end UTC and exit statuses in `ch04-final-sequence.log`:

```sh
# Student repo
go test -count=1 ./...
go test -race -count=1 ./...
# Book repo
go run ./cmd/grade-edition3 -ch 4 \
  -core /Users/bill/projects/ensemble-edition-3-student \
  -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer \
  -consumer-case /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json \
  -json /tmp/ch04-core
go run ./cmd/grade-edition3 -ch 1 -json /tmp/ch04-core
```

Current known final result: full suite passed in125.212s; remaining sequence pending. Final report will replace this status only after all processes exit. Supplemental focused suite, `go vet ./...`, builds and `git diff --check` passed on the frozen source.

Focused tests cover ordinary file success/failure, full persisted output, capped head/tail/cursor behavior, limits consumed by controls and failures, invalid options before shell side effects, registered handlers/retained builtins, owner chain, blocked waiters, unkillable Go handlers and late result suppression, PTY input/TERM/size/CRLF/exit, per-call cwd and fresh-shell environment, unseen-only patterns, process-group descendant kill, actual running shutdown and shutdown log reason, history load without live jobs, independent roots/snapshot isolation, persistent-file collision protection and concurrent process-launch/shutdown.

## Live entry, identity and observations

Only approved credentials were injected through process environment from `~/.cr/settings.json`; no key values, copied secret env file, settings content, raw provider wire dump or private opaque signed blobs are in the report/evidence. Neutral logs sanitize opaque fields. Retained launcher scripts describe configuration without secret values. Each run has separate workspace/data directories and exact prompt, requested model, endpoint, native surface, binary digest, source tree, exit code and elapsed time metadata.

|Run|Requested model / observed neutral provenance|Outcome|
|---|---|---|
|Anthropic initial core|claude-sonnet-4-6 / messages, same identity|exit1 after10 rounds; 43 JSONL records; first turn only, no claimed completed conversation|
|OpenAI core|gpt-4.1-2025-04-14 / chat_completions, same identity|exit0; two turns; 48 records; actual file repair, process handle4 killed, no replay, token continuity; final printf refused, not executed|
|Gemini core|gemini-3.8-flash / generate_content, same identity|exit0; two turns;44 records; repair, process handle4 kill/no replay, continuity, successful limits-consuming printf; no downgrade|
|Anthropic external consumer|claude-sonnet-4-6 / messages, same identity|exit0;8 records; actual shout returns HELLO FRAMEWORK; timestamped `shout called` diagnostic from root logger|
|Anthropic bounded follow-up|claude-sonnet-4-6 / messages|pending|

The event provenance uses the accepted neutral model identity, not an independent server-attestation claim. No mismatched model routing identifier was observed in the retained logs. All live results are attributed to the listed binary bytes; staged source commit identity will be finalized after gates.

Actual useful live asynchronous workload and EOF/Delve evidence: **pending follow-up**. The eight-second fixture is intentionally shorter live async-use evidence, not a measured ninety-second test. Source ownership/concurrency/lifetime claims rest on source review, real focused tests, race suite and public gates, not provider prose. A live sleep is not credited as development-server proof.

## Candidate and remaining status

No final acceptance claimed. Final gates, bounded follow-up, evidence credential/opaque scan, exact commit/tree identity and clean status still to be finalized. No further source changes are planned. Parent/root remains orchestrator and reviewer.
