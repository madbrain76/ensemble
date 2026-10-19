# Chapter 1 restart — student review

Implemented current Chapter 1 from scratch in `solutions/edition-3/main`; ready for independent review. Acceptance and the `ch01` snapshot remain root's job.

## Architecture

`ensemble.go` is the public composition root, owning its logger/Agents and wiring private implementations. `internal/common` contains shared interfaces/values, not state. Agent owns configuration/completed history; Engine owns raw HTTP/measured usage and obtains configuration/diagnostics through its actual parents. No sibling imports, tools, retries, provider framework, or later-chapter machinery. Calls within an Agent are synchronous.

Failed exchanges do not commit history. Valid measured usage is counted even for unusable text; absent usage is not zero. State accessors return copies. HTTP has no blanket model-work timeout: standard connection/TLS protections remain and caller context controls header/body waiting. Redirects are refused; diagnostics exclude remote bodies/URLs/secrets. CLI defaults to JSON lines, with opt-in chat. The corrected line reader has no arbitrary size cap and handles final unterminated lines, EOF, and read errors.

## Commands and outcomes

Outputs/executable are retained only in `/tmp/ensemble-ch01-restart.bj87Gl` (`$T` below). Commands ran in the implementation module unless noted, with clean fake-credential environments. Disk checks at build boundaries showed about 14 GiB available; no shared cache was deleted.

- Initial `go build -o "$T/ch01" ./cmd/ch01`: exit 0 (`build.log`). Then held unchanged for root's requested early comparison.
- Round1: `go test ./cmd/ch01 -run '^TestLargeJSONLine$' -count=1` reproduced the 70 KiB well-formed-line failure: `reading input failed`, zero HTTP calls, exit 1 (`large-line-before.log`). Replacing Scanner with `ReadString` made that same test pass (`large-line-after.log`).
- `go test ./... -count=1 -v`: most cases passed, but my cancellation peer blocked `httptest.Server.Close` because it had not consumed the request body. Client cancellation had returned. I captured SIGQUIT from only that disposable test process; exit 1/full stack retained in `focused-tests-peer-hang.log`. After draining the peer body and adding a cleanup release, `go test . -run '^TestCancellationWhileWaiting$' -count=1 -v` passed (`cancellation-fixed.log`).
- `go test . -run '^(TestValidationBeforeHTTP|TestTransportErrors|TestRedirectIsNotFollowed)$' -count=1 -v`: passed (`focused-remaining.log`). Tests exercise real local HTTP, ownership/config/logger isolation, history, usage/errors, cancellation, and CLI line/chat behavior.
- `python3 "$T/mutation-check.py"`: exit 0; 16 representative mutations in disposable copies of my own source were rejected by their intended tests (`mutations.log`). Candidate source was not mutated.
- Stable gates: `go test -race ./...`, `go vet ./...`, and `go build -o "$T/ch01" ./cmd/ch01`: all exit 0 (`race.log`, `vet.log`, `final-build.log`). Only README/test comments and this note changed afterward, not executable behavior.
- From the book repository root, exactly one final `go run ./cmd/grade-edition3 -ch 1 -json /tmp/ensemble-ch01-restart.bj87Gl/ch01`: exit 0, independently parsed `passed: true`, 100/100. Real 62-second header and 62-second body phases both passed; elapsed 125 seconds, observation interval 360 seconds, no outer deadline (`grade.json`, `grade.stderr`, `grade-status.log`). No grader internals were inspected or successful grade repeated.

## Bounded live exercise

Root authorized one discovery GET plus three Messages POSTs through a small temporary Python proxy, using the same graded binary in `chat` mode. Only the proxy read the exact authorized credential into memory; verified HTTPS used `/etc/ssl/cert.pem`. CLI received a dummy key, loopback URL, discovered model, clean environment, and owned temporary HOME. Bodies were unchanged; slots reserved before forwarding; no redirects/retries/pagination. The 180-second exercise deadline did not change production HTTP behavior.

Command: `env -i HOME="$HOME" PATH=/usr/bin:/bin:/usr/sbin:/sbin /usr/bin/python3 "$T/live-proxy.py"`. First setup exited 1 before networking because root specified absent field `ClaudeAPIKey`: zero discovery/messages, no CLI/listener launched (`live-setup-wrong-field.log`). Root acknowledged its error and authorized exact field `directClaudeAPIKey`; I did not search for alternatives. The corrected setup used the remaining original budget:

- Discovery returned selected model `claude-sonnet-5-5`. One GET and three POSTs all returned HTTP 200.
- Turn 1 supplied Cedar/version 7. Reply: “Got it: your project is named Cedar, and it's at version 7.” Usage 92 input / 23 output.
- Turn 2 asked 7 plus 5 while retaining context. Reply: “7 + 5 = 12, so for Cedar, going 5 versions past version 7 would land on version 12.” Usage 141 / 80.
- Turn 3 asked for original facts and the sum. Reply: “You told me the project is named **Cedar** and its version is **7**. The sum from our previous exchange was 7 + 5 = 12.” Usage 212 / 45. Semantic recall passed.
- CLI cumulative totals: 92/23 → 233/103 → 445/148. EOF/CLI exit 0, proxy exit 0, elapsed 4.63 seconds. Owned CLI reaped/listener closed. No further provider requests.

Sanitized chat is retained in `live-chat.stdout` / `live-chat.stderr`; counts/status/usage/body hashes/exits in `live-run.log`. No raw API captures or credentials were saved.

## Author feedback and assistance

Specific feedback: the chapter's default Scanner example silently imposes roughly a 64 KiB physical-line limit. Prefer `bufio.Reader.ReadString` with final-line/EOF explanation unless an explicit input-size limit is part of the contract. No ownership/lifetime prose change suggested.

Assistance: root coordination and public round1 feedback from root/independent review identified the scanner limit; root supplied/corrected live credential handling. I used no children, browser, reference implementation, grader internals, retired solution, or later chapters. Source/tests are new work from allowed public instructions. No installs, restarts, pushes, commits, nested repositories, or student-created snapshots.
