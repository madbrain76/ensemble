# Chapter 4 — restart student notes

Implemented from the accepted Chapter 3 predecessor. I did not read a reference solution, private review, or grader implementation. Root and the same independent reviewer supplied public-contract findings in two correction rounds; no architecture restart was needed.

## Implementation and teaching feedback

The compact design worked: Ensemble owns the shared integer handle table; each Agent owns its registry, limits, base directory, and Jobs collection. Ordinary built-ins and registered handlers share dispatch, disk output, and lifecycle. The separate `examples/ch04-consumer` module exercises the importing API and the real ToolCall → Agent → Ensemble logger chain. Chapter 5 mailbox/actor behavior was not added.

The difficult distinction was **execution lifetime versus observation**. A wait must wake on new unseen output without cancelling execution; a fast initial call still needs complete small-result parity. Separating initial observation from subsequent waits made that distinction explicit without adding a scheduler. Related requirements worth keeping together in the teaching text are:

- A killed job is immutable, but observing it is not itself a tool error. Exit-code setters, late writers, reports, and serialized kill events must agree.
- One-shot limits are consumed by literally the next call, including errors and another setter. The consumed-settings diagnostic must also reach the recorded tool result.
- PTY line submission, `TERM=dumb`, 50×200 geometry, and CRLF normalization before disk/cursor/regex accounting form one contract. A pending trailing CR must survive split reads without corrupting UTF-8 or lone CR bytes.
- Failed output setup must still produce exact-ID tool-call/result pairs and allow the rest of the model’s batch to be recorded, without invoking unowned effects.

Review also caught missing cwd metadata, buffered list/search output, and incomplete fast-command/disk parity. These were implementation omissions, not reasons for a plugin or job framework.

## Checks, failures, and assistance

Final sequential root race tests and vet passed. GUI and consumer test/vet/build commands passed, but those modules have **no test files**; this is compile/static evidence, not behavioral GUI testing. Focused tests exercised real files, PTYs, processes, asynchronous writers, late handlers, cursor wakeups, split normalization, and Agent/local-HTTP correlation.

The first Chapter 4 grading attempt failed **110/160**. After the second reviewed correction round, the new frozen CLI passed Chapter 4 **160/160** and Chapter 1 **100/100**, both with parsed `passed: true` and process exit 0. Both genuine 62-second header/body checks completed naturally.

Authorized inherited-fixture changes covered six→ten tools, full read output on disk with bounded observations, and merged PTY output instead of separated shell JSON. A later truncated-read fixture incorrectly assumed no report header; my first replacement also misstated the marker literal. Both failed race logs were retained before the corrected fixture passed. Explicit read bounds/errors, exact retained bytes, shell environment isolation, normal/nonzero exits, and launch failures remained asserted.

The local manual session proved natural **90.017-second** completion after an early callback and same-handle wait. Actual Agent tools drove Delve through breakpoint, continue, step, value `42`, and exit 0. Its first attempt failed because I guessed a `/tmp` source path; querying Delve’s actual `/private/tmp` mapping fixed only the disposable fixture. A real CLI-started loopback server and descendant were alive before EOF and gone afterward, including the process group; the persisted shutdown event had no invented exit code.

## Bounded live result and limits

The approved CLI cases used Anthropic `claude-sonnet-5-5`, OpenAI Chat Completions `gpt-4.1-2025-04-14`, and Gemini generateContent `gemini-3.8-flash`, with explicit resolved targets. Each completed one harmless `printf` job with paired IDs, exact retained output, exit 0, and measured usage. The importing consumer completed one `shout` call and logged through its public parent chain. These model identities are observations for this run, not permanent availability claims.

There were **9 actual upstream requests**, including an initially rejected OpenAI response. Its quoted printf contained literal LF rather than backslash-n; the first proxy guard was too strict. Work stopped, and root explicitly authorized accepting only those two equivalent spellings and continuing with the original request counts. Anthropic was not repeated. All upstream statuses were 200; the successful child processes exited 0. Native output limits stayed exactly 1024 tokens/request, with three-request case ceilings and a twelve-request cumulative ceiling. Keys remained only in proxy memory/HTTPS headers; verified TLS, isolated dummy-key child environments, and unchanged frozen binaries were used. No USD cost is claimed without pricing.

Raw successes and failures remain locally under `/tmp/e3-ch04-student/logs` and `/tmp/e3-ch04-student/live-control`; disposable job outputs are under the owned temporary work directories. No broad provider campaign, browser/audio/focus exercise, or additional production correction was performed. All owned activity is closed; source and binaries are held for final result review.
