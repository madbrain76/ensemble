# Chapter 3 live plan — written before provider discovery/live calls

Only `/tmp/ensemble-ch03-evidence/student-cli` (rebuilt from this candidate) will
act. Each provider gets a separate `/tmp/ensemble-ch03-live-<vendor>-...` cwd.
The disposable Go repair project contains a deliberately broken Add function,
a small failing test and a README. No tools run in main/student repositories.
Policy is explicitly selected by `LLM_REQUIRE_TOOL_EXPLANATION=true` for these
configured models; this is not a claim that every model universally requires it.
Current provider model discovery precedes selection; Gemini is `gemini-3.8-flash`
or stops to ask parent (no silent alternative). OpenAI uses ChatCompletions only.

| Interactive feature | Planned live step | Independent proof |
| --- | --- | --- |
| model/tool capability | discover current models, then real tool-bearing Ask | successful provider HTTP/model identity, parsed tool calls/results |
| orientation/list_directory | Ask 1: list current disposable directory | saved actual directory inventory and history tool/result summary |
| read_file and search_files context | Ask 1: read README, search Add in Go files with context | requested args names and actual returned content hashes/counts |
| run_command and error feedback | Ask 1: run `go test ./...` on broken project | result exit nonzero, next provider request contains feedback |
| edit_file repair | Ask 1: fix subtraction to addition via edit_file | independently inspect Add source bytes and run Go tests |
| write_file | Ask 1: create a new repair note | actual new file, content hash |
| preserve follow-up conversation | Ask 2 in same CLI process: refer to prior repair, test again, update note with append | two successful assistant replies, chronological history and paired results |
| preceding visible explanation | both turns under selected public policy | inspect each response: nonblank text before first call, or paired refusal then corrected retry |
| credential subprocess hygiene | benign-canary actual-shell unit test before live | no supported dummy key variable reaches shell; ordinary setting remains |
| cancellation/time boundary | HTTP/provider and subprocess fixtures | caller cancellation kills process group; no invented short shell-work deadline |

Launcher reads only the three explicitly authorized JSON fields in process
memory, never prints keys or credential-store content, and passes the selected
key in the CLI's environment (never argv or prompts). Tool subprocesses remove
all four supported provider-key variables. No whole store is copied. Stored
live evidence is sanitized metadata/results, not credentials or raw provider
bodies. Each provider has at most one initial repair attempt and one bounded
retry if needed; process deadline is 240 seconds per session. Provider/network
failures, refused or missing tool use, and limit hits are reported honestly.

## Final human-chat exercise (parent approved before execution)

After the protocol sessions, run exactly one additional `student-cli chat`
session per discovered selected provider/model, using the final frozen binary
and fresh disposable projects. Plain-text human input exercises the actual
prompt/reply loop with the same repair and follow-up feature mapping above.
The launcher watchdog is 120 seconds per session, outside application/model
HTTP lifetime. A timeout is incomplete, never success; no automatic retries.
The Gemini retry already completed actual repair and follow-up but returned
empty final answer strings (opaque-only terminal responses); do not present
that as a complete human-facing reply. Normalized History does not retain a
stop/finish reason, so a token-limit explanation is unproven speculation.
