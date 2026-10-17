# Edition 3 Chapter 4 student report

## Scope and identity

- Authorized launch instructions: `/tmp/cadence-ch04-student-task.txt`, SHA256 `f90d785c5b3ba37182e3e0170566d621fc0f79ef40fb7d125fd9343b049eeac8` (verified).
- Public chapter: `/Users/bill/projects/ensemble/book/edition-3/chapter-04.md`, SHA256 `016c4e88af4bcb657aad6733e5f65aa09509b0397e84a37059ca245c588eaa55` (verified).
- Student repository: `/Users/bill/projects/ensemble-edition-3-student`; accepted clean predecessor `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, tag `edition-3-ch03-r1`.
- Required public GUIDE, author/student role, coding rules and whole Chapter 4 read before implementation. Ownership/lifecycle plan sent before edits; parent approved.
- Frozen staged source tree: `840b5b3bab4daf7490b5555451bb6506b273d7b9`. Final commit: `062c146cafa9b105745f6b44921160198bbaee32` (tree match and clean status verified).
- Tested/live core binary: `/tmp/ch04-core`, SHA256 `b0ed87f60783977c2024547690b4865ff6269da0f34a51c484f39901eb003a57`.
- Independent consumer binary: `/tmp/ch04-consumer`, SHA256 `d6e4074774199f5cc2476b5fd648bf14cb60d6325e6c94d09158fd53adaf8dd2`.
- Both binaries were built from this exact frozen dirty/staged tree before commit. Live metadata therefore records the predecessor HEAD plus the frozen source tree, not a false claim that a later commit had already existed. Final commit has the identical tree. No production source changes during final gates/live exercises.
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
10. Parent approved exactly one bounded Anthropic follow-up with four prompts, unchanged binary/model/round cap: real 8-second Go test with failure→repair→PASS; real HTTP service across turns; real Delve commands; still-running server immediately before EOF with OS process observation and shutdown event. The chapter requires Delve, not GDB; no additional provider campaign. Follow-up exited0 in77.856s, 62 neutral JSONL records. All four turns completed; exact observations are detailed below.

## Final sequential validation — completed

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

Final results from exact retained receipts (all exit0), in strict sequence:

|Gate|Result|UTC start → end|
|---|---|---|
|Full suite|PASS; ensemble125.212s|2026-10-17 00:36:27 → 00:38:32|
|Race suite|PASS; ensemble126.458s|00:38:32 → 00:40:39|
|Public Chapter4|passed=true;160/160; every check passed|00:40:39 → 00:40:48|
|Public Chapter1|passed=true;100/100; both62s header/body lifetime checks passed|00:40:48 → 00:42:52|

Receipts: `raw/ch04-final-sequence.log`, `raw/ch04-full-final.log`, `raw/ch04-race-final.log`, `raw/ch04-grade4-final.log`, `raw/ch04-grade1-final.log`. They completed before transport recovery; they were not rerun merely because the runtime restarted. Supplemental focused suite, `go vet ./...`, builds and `git diff --check` passed on the frozen source. Public reports explicitly note that the62s measurements alone cannot exclude every longer arbitrary deadline; source review supplies that boundary.

Focused tests cover ordinary file success/failure, full persisted output, capped head/tail/cursor behavior, limits consumed by controls and failures, invalid options before shell side effects, registered handlers/retained builtins, owner chain, blocked waiters, unkillable Go handlers and late result suppression, PTY input/TERM/size/CRLF/exit, per-call cwd and fresh-shell environment, unseen-only patterns, process-group descendant kill, actual running shutdown and shutdown log reason, history load without live jobs, independent roots/snapshot isolation, persistent-file collision protection and concurrent process-launch/shutdown.

## Live entry, identity and observations

Only approved credentials were injected through process environment from `~/.cr/settings.json`; no key values, copied secret env file, settings content, raw provider wire dump or private opaque signed blobs are in the report/evidence. Neutral logs sanitize opaque fields. Retained launcher scripts describe configuration without secret values. Each run has separate workspace/data directories and exact prompt, requested model, endpoint, native surface, binary digest, source tree, exit code and elapsed time metadata.

|Run|Requested model / observed neutral provenance|Outcome|
|---|---|---|
|Anthropic initial core|claude-sonnet-4-6 / messages, same identity|exit1 after10 rounds; 43 JSONL records; first turn only, no claimed completed conversation|
|OpenAI core|gpt-4.1-2025-04-14 / chat_completions, same identity|exit0; two turns; 48 records; actual file repair, process handle4 killed, no replay, token continuity; final printf refused, not executed|
|Gemini core|gemini-3.8-flash / generate_content, same identity|exit0; two turns;44 records; repair, process handle4 kill/no replay, continuity, successful limits-consuming printf; no downgrade|
|Anthropic external consumer|claude-sonnet-4-6 / messages, same identity|exit0;8 records; actual shout returns HELLO FRAMEWORK; timestamped `shout called` diagnostic from root logger|
|Anthropic bounded follow-up|claude-sonnet-4-6 / messages, same identity|exit0; four turns;62 records;77.856s; actual8s test failure→repair→PASS, live HTTP server, Delve, shutdown at EOF|

The event provenance uses the accepted neutral model identity, not an independent server-attestation claim. No mismatched model routing identifier was observed in the retained logs. All live results are attributed to the listed binary bytes and frozen source tree. Metadata correctly retains the predecessor HEAD at execution time; the final commit below records that same tree after validation.

### Observed useful workload, Delve and EOF evidence

`anthropic-followup-tool-evidence.json` extracts actual returns; `anthropic-core-followup.events.jsonl` is the sanitized source of record:

- Actual `go test -v .` returned running handle1. A later wait captured `TestGreeting (8.00s)` failure, `want hello, got "helo\n"`, and exit1. After read/edit, another actual `go test` returned running handle6; a later wait captured PASS after8.00s and exit0. The fixture is intentionally an eight-second useful asynchronous test, **not** a measured ninety-second test.
- A real Python HTTP server on localhost59154 returned running handle2 with `Serving HTTP`. Separate commands curled it across turns: initial `helo`, repaired `hello`, and final `hello`, each exit0. The same server stayed alive while the conversation and actual Go tests/debugger continued. A live sleep is not credited as server proof.
- Real `/Users/bill/go/bin/dlv exec ./debug-target` returned running handle8. `send_input` sent newline-separated `break main.main`, `continue`, `next`, `print value`, `quit`. The first consumed result contains `Breakpoint 1 set`; the next wait on **the same handle8** contains the actual `[Breakpoint 1]` stop at main.go:3, stepping to main.go:5, `(dlv) 41`, and `exit_code: 0`. These are actual debugger observations, not model prose. Source JSONL lines42,46,50 preserve launch/input/wait.
- The initial postprocessor's `real_breakpoint_continue:false` is preserved in `anthropic-followup-checks.json`. Its predicate incorrectly required the set and hit strings in one result, overlooking normal cursor consumption across successive results. `anthropic-followup-checks-reconciled.json` records the source-log SHA, same handle, exact original returns/line numbers and separate successful set/hit/step/value/exit predicates. No tool output was changed and no extra provider call was made to obtain this correction.
- Final tool observation was `wait_for_job(handle2,delay0.01)` reporting server2 **running** and its three real HTTP requests. No intervening tool call follows. EOF shutdown produced matching `job_killed` handle2, reason `shutdown` in the saved neutral log. Independent OS observation recorded real Python PID83531 alive until0.264340s before CLI exit and absent afterward (`anthropic-core-followup.processes.json`). The model's statement that the server “shuts down naturally” is imprecise: the evidence proves root-owned shutdown, not spontaneous server completion.
- Token MAPLE continuity and repaired file state were exercised across the four turns. The original failed Anthropic run and empty late monitor remain separately retained and are not substituted for this successful follow-up.

Source ownership/concurrency/lifetime claims rest on source review, real focused tests, race suite and public gates—not provider prose. The current chapter requires Delve, not GDB. The real Delve session fulfills that interactive debugger exercise; there is no host-conditional GDB requirement.

## Interruption, recovery, assistance and provenance

After the completed work above, the human reported hotspot transport loss; a runtime restart removed the live child registry before submission. This is continuation of the same student attempt, not a new attempt or correction round. Original source, binaries and evidence survived. Parent reported the completed receipts; on recovery I independently checked their contents, actual statuses, source tree and binary hashes. The unfinished report is preserved verbatim as `raw/report-before-recovery.md`. No paid request, suite or public grader was repeated solely for recovery, and no production source was changed after the frozen final runs.

Assistance disclosed: parent supplied the authorized public launch task, approved the initial ownership/lifecycle plan, corrected my misreading of `-core` and the exact provider models/credential source **before live entry**, approved the saved public live checklist and one bounded four-prompt follow-up, and identified that sleep alone would not demonstrate useful server interaction. Parent supplied recovery locations and noted the debugger predicate discrepancy; I reconciled it against my original actual job returns. No helper agents were launched. No private reference solution, private review, grader internals or future chapter was used. No borrowed completion result is credited as my own observed test; all reported results point to the retained receipts or outputs.

`source-manifest.json` records all staged paths, Git blob identities and SHA256s. `local-artifacts-manifest.json` distinguishes temporary binaries/fixtures/persisted outputs from portable evidence. `redaction-verification.json` records a process-memory-only scan against the three authorized configured credential values and a recursive opaque-field check: no findings. No credential values are printed or written by the verifier. The original paid-run launchers and prompts are retained as operational evidence; they are not an instruction to rerun the campaign. The source tree is unchanged by recovery/report/evidence updates.

## Candidate status

All required final local/public gates passed and the bounded live evidence is complete within the stated limits. Final candidate: `062c146cafa9b105745f6b44921160198bbaee32` (`Implement Chapter 4 supervised jobs and public tool registration`), parent `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, tree `840b5b3bab4daf7490b5555451bb6506b273d7b9`. Commit contains exactly the previously staged/tested source; no source correction occurred after the final runs. `git status --porcelain=v1` is empty; both unstaged and staged diffs are empty. `candidate-identity.json` records the verified identity and unchanged binary hashes; `evidence-manifest.json` hashes this report and every retained evidence file (excluding itself). This is a student candidate, **not independent acceptance**. Parent/root owns acceptance, any accepted tag, and the later author phase. No further source changes or live calls are requested.
