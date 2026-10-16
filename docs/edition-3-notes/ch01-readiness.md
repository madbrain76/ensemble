# Edition 3 Chapter 1 readiness

## Scope and source identity

Readiness checker, **not a student run**. Initial source HEAD:
`2be995ebc12c6a7efb21cee3b9740a758489adda`.
Inspection and checks performed 2026-10-16 (UTC). No commit made.

Read the complete root `AGENTS.md`,
`.agents/skills/ensemble-coding/SKILL.md`,
`docs/edition-3-execution.md`, `docs/edition-3-carryover.md`,
`book/edition-3/chapter-01.md`, and the inherited
`docs/edition-2-carryover.md`. Read the original Chapter 1 runner, evaluator,
script, fake provider, Edition 2 overlay, command, and relevant tests.
The skill was reread completely before implementation.

The initial authorization permitted only this note and isolated build outputs.
Subsequent explicit parent authorization permitted the bounded addition of:

- `cmd/grade-edition3/main.go`
- `internal/gradeedition3/ch01.go`
- `internal/gradeedition3/ch01_test.go`
- this readiness note

No original or Edition 2 grader, historical solution, application, or manuscript
was edited by this checker. No historical score was changed. Historical solution
execution below is a grader **control**, not Edition 3 student evidence. No
browser/MCP connection, installation, restart, student launch, or commit occurred.
Other concurrent preparation changes appeared in repository status; those are not
this checker's changes and were not reverted.

Source revision evidence:

- Chapter 1 Edition 2 overlay's last source commit:
  `02efcd48324beac4c10b506194cbeb1de55ea034`.
- Original evaluator/script/fake's last source commit:
  `515e2d883ef0994db652a031eb13e6974ba7b829`.
- Edition 3 Chapter 1 source's last source commit:
  `2be995ebc12c6a7efb21cee3b9740a758489adda`.
- SHA-256 `book/edition-3/chapter-01.md`:
  `76e2ed60924f9eb0e1c7790be2843a5c1d96ae5a0a200e3394cf6f80133e40bf`.
- SHA-256 original `internal/grade/harness.go`:
  `43f93da9b67504fc9cb69a88572fdc3ebd382463a2d5e31c6703fd422c01fb7c`.
- SHA-256 original `internal/grade/checks.go`:
  `271efc2f59265ca07e9cd6d0a55d5db522db540330005949167cddb1d10c5ae4`.
- SHA-256 original `internal/fakeanthropic/fake.go`:
  `4f0326f1359774b75529b798ee4ccaa860f903d6f6ff666ee6cf82d8240fdfcd`.
- SHA-256 inherited `internal/gradeedition2/ch01.go`:
  `01385980eca5813b405acf9e871265435f0f57c5fbea8a87f0866964ad7277a2`.
- SHA-256 historical `solutions/ch01/internal/engine/engine.go`:
  `1dbbf5119bb160f5702aa29d59935af90ae7b82edceae5c2c069cb8c076379b6`.

The new overlay is an uncommitted working-tree addition, not part of the original
HEAD. Final validation/source status is recorded below.

## Authoritative invocation and initial harness result

**Working directory: `/Users/bill/projects/ensemble` (repository root).**
Root `go.mod` is module `github.com/waywardgeek/ensemble`, Go 1.25.
The available toolchain reported `go version go1.25.4 darwin/arm64`.
The nested `agent/` module is not the grader's working directory.

Inherited command, verified from source and executed:

```sh
cd /Users/bill/projects/ensemble
go run ./cmd/grade-edition2 -ch 1 -json /submission
```

In this environment `/submission` is absent. The literal invocation produced
`grader: stat /submission: no such file or directory`, `exit status 2`;
`go run` returned shell status 1. This is a missing-target failure, not a grader
build failure. The compiled grader's own corresponding status is 2.

The inherited CLI defaults to `solutions/ch01` when the path is omitted
(`cmd/grade-edition2/main.go:45-48`). Therefore **an explicit fresh path is
mandatory for Edition 3**. Do not interpret an omitted-path historical pass as a
fresh Chapter 1 run. The new Edition 3 CLI rejects an omitted path and only
implements Chapter 1:

```sh
cd /Users/bill/projects/ensemble
go run ./cmd/grade-edition3 -ch 1 -json /submission
```

No book command was patched. The readiness/execution caller should deliberately
select this new overlay for the timeout gate while retaining the inherited
command as historical evidence.

Exact offline validation commands (local toolchain; no dependency downloads):

```sh
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
tmp=$(mktemp -d /tmp/ensemble-ch01-readiness.XXXXXX)
go build -o "$tmp/grade-edition2" ./cmd/grade-edition2
"$tmp/grade-edition2" -h
go test ./internal/grade ./internal/gradeedition2 ./internal/fakeanthropic \
  -run '^(TestBuild|TestChapter1VersionAndPreservedChecks)$' -count=1 -v
rm -rf "$tmp"
go run ./cmd/grade-edition2 -ch 1 -json "$PWD/solutions/ch01"
go test -race ./internal/grade ./internal/gradeedition2 ./internal/fakeanthropic \
  -run '^(TestBuild|TestChapter1VersionAndPreservedChecks)$' -count=1
```

Results: build/help passed; focused normal and race tests passed; historical
control returned **100/100, `passed:true`, shell status 0**. All seven inherited
checks passed: `protocol`, `wire`, `calls`, `replies`, `memory`, `growth`,
`usage`. `TestBuild` passed its five layout/error cases (`root`,
`command_container`, `preferred_command`, `command_library`,
`broken_preferred_command`). `TestChapter1VersionAndPreservedChecks` verifies
baseline acceptance and the captured-version mutation. The fake provider package
has no standalone test files. Temporary build/control directories were cleaned
by their owning commands.

## Contract consistency and precise coverage gaps

The existing fake checks model and API-key overrides exactly; these are **not**
missing checks (`internal/fakeanthropic/fake.go:203-228`). It also checks required
headers, base-URL traffic, request shape, a positive max-token budget, and non-streaming.
The Edition 2 overlay adds the exact `anthropic-version: 2023-06-01` check
(`internal/gradeedition2/ch01.go:11-30`). The five-turn script/evaluator cover
full growing history, nonempty system text, concatenating multiple text blocks,
ordinary input/output usage sums, response shape, and clean stdout/EOF exit.

| Finding | Exact evidence | Readiness disposition |
| --- | --- | --- |
| **Blanket generation deadline accepted by inherited grader** | Edition 3 Chapter 1:365-366 explicitly prohibits a fixed wall-clock timeout over the whole HTTP request/body. Historical `solutions/ch01/internal/engine/engine.go:34` constructs `http.Client{Timeout: 60 * time.Second}` yet the exact inherited CLI above returns 100/100. | Bounded Edition 3 overlay added; two transport-phase controls described below. Historical graders/results remain unchanged. |
| Test watchdog is not application-timeout coverage | `internal/grade/harness.go:21-22` has 20-second round / 15-second exit watchdogs; fake responses are prompt. The harness never waits beyond the historical client deadline. | Preserve these historical watchdogs. New lifetime probes have their own longer **grader** watchdog; it is not an application deadline recommendation. |
| Thinking configuration/visible-text distinction untested | Chapter 1:420-428 explains ignoring non-text blocks and explicitly says the grader does **not** check that filtering step; the Appendix advises avoiding explicit thinking at this stage. Fake request types and ordinary replies do not exercise thinking. | Record as ungraded/source-review territory, not a new hidden scoring gate. Do not claim the timeout overlay tests it. |
| Cache-counter separation and missing usage members untested | Chapter 1:565-567 requires plain input/output sums, excluding cache counters, and zero for missing members. Fake usage in `internal/fakeanthropic/fake.go` has only input/output; scripted replies supply them. | Future targeted fake cases could cover absent counters and nonzero cache counters. Not necessary to introduce a new framework for the present timeout fix. |
| Provider errors/malformed response handling not systematically checked | Chapter 1:524-531 and 569-571 specify a useful machine-readable failure. The five-turn script exercises successful replies; invalid-wire failures are not equivalent to deliberate error-response controls. | Source review remains required; future narrowly scoped error controls are justified. No retries required or added. |
| CLI compatibility flags not exercised | Chapter 1:357, 527, 533-537 specify optional `--stdio` and accepting/ignoring `--loop`; inherited `grade.Start` launches with no arguments. | A future invocation probe is reasonable; this readiness patch is limited to the demonstrated timeout carryover. |
| Architecture/ownership/author review is not automatic | Chapter 1:113-177 and 438-503 require transport/session/I/O ownership reasoning; checks operate on observable traffic and output. | Keep source/architecture/author review as execution-loop gates. A 100 score alone is not acceptance. |

A finite delay cannot prove the absence of *every* arbitrary deadline (for example
five minutes). The explicit remaining review is to inspect the student's client
and request contexts for `http.Client.Timeout`, `ResponseHeaderTimeout`,
`context.WithTimeout`/`WithDeadline`, body-read timers, or equivalent blanket
limits. Legitimate connect/TLS limits are different and are not prohibited by
these probes. No exhaustive timeout detector or broad future-chapter framework
was introduced.

## Edition 3 overlay and focused validation

The new command reuses `grade.Build`. `gradeedition3.Chapter1` first runs
`grade.Run` and evaluates the unchanged `gradeedition2.Chapter1` checks. It then
adds required zero-point checks `generation_headers` and `generation_body`.
The total remains 100; **inspect `passed`, not score alone**. A submission can
retain 100 historical points and fail a lifetime gate, yielding CLI status 1.

Each extra probe launches a fresh real process in a temporary directory against
a local HTTP fake. One delays response headers for **62 seconds**. The other
flushes HTTP 200 headers immediately, then delays the response body for 62
seconds. Each requires exactly one request, the eventual assistant answer,
matching usage (7 input/3 output), successful EOF exit, and no unexpected protocol
lines. The watchdog is delay + inherited round grace + inherited exit grace
(97 seconds per phase), only for terminating stuck grading processes. A canceled
client releases the fake handler. No live provider is used by these checks.

Focused controls copy the historical source into `t.TempDir`, change only the
known timeout expression, compile a real executable, and run the same probe:
zero timeout succeeds through both 250ms delays; a 25ms blanket timeout fails
both. Assertions reject unrelated zero-request or watchdog failures as mutation
success. A separate test pins the production delay strictly above 60 seconds.
Historical source is never rewritten.

Exact commands:

```sh
gofmt -w cmd/grade-edition3/main.go internal/gradeedition3/ch01.go \
  internal/gradeedition3/ch01_test.go
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off \
  go vet ./cmd/grade-edition3 ./internal/gradeedition3
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off \
  go test -race ./cmd/grade-edition3 ./internal/gradeedition3 \
    ./internal/grade ./internal/gradeedition2 ./internal/fakeanthropic \
    -run '^(TestGenerationLifetimeControls|TestHistoricalTimeoutIsExceeded|TestBuild|TestChapter1VersionAndPreservedChecks)$' \
    -count=1 -v
```

First attempt: vet and new-package build failed because the new caller omitted
`NewTitledReport`'s third (stderr) argument. Corrected the new call to pass
`res.Stderr`; reran exactly these commands. Final focused result: **vet status 0;
race-enabled tests status 0**, all four short lifetime control cases passed and
both inherited test groups passed. No historical source fix was needed.

Full-duration CLI controls also completed successfully. Executed control setup
and CLI commands (report-printing/assertion helper omitted here):

```sh
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
tmp=$(mktemp -d /tmp/ensemble-e3-ch01-controls.XXXXXX)
trap 'rm -rf "$tmp"' EXIT
cp -R solutions/ch01 "$tmp/positive"
python3 - "$tmp/positive/internal/engine/engine.go" <<'PY'
from pathlib import Path
import sys
p=Path(sys.argv[1]); s=p.read_text()
a='&http.Client{Timeout: 60 * time.Second}'
assert s.count(a)==1
s=s.replace(a,'&http.Client{}').replace('\t"time"\n','')
p.write_text(s)
PY
go build -o "$tmp/grade-edition3" ./cmd/grade-edition3
"$tmp/grade-edition3"
"$tmp/grade-edition3" -ch 2 /submission
"$tmp/grade-edition3" -ch 1 -json /submission
(go run ./cmd/grade-edition3 -ch 1 -json "$PWD/solutions/ch01" \
  >"$tmp/historical-negative.json" 2>"$tmp/historical-negative.stderr"; \
  echo "$?" >"$tmp/historical-negative.exit") & negative_pid=$!
(go run ./cmd/grade-edition3 -ch 1 -json "$tmp/positive" \
  >"$tmp/corrected-positive.json" 2>"$tmp/corrected-positive.stderr"; \
  echo "$?" >"$tmp/corrected-positive.exit") & positive_pid=$!
wait "$negative_pid"; wait "$positive_pid"
```

The actual script disables shell `errexit` around the deliberately failing
invocations, restores it before assertions, and cleans the temp directory on
exit. No source copy was placed in the repository.

| Control | Inherited checks | Header/body lifetime gates | Final result |
| --- | --- | --- | --- |
| Unchanged historical 60-second client | 7/7 passed, 100/100 | both failed; exactly one API request each, no complete reply/usage | `passed:false`, `go run` exit 1 |
| Temporary corrected `http.Client{}` | 7/7 passed, 100/100 | both passed after 62-second delays | `passed:true`, `go run` exit 0 |
| Compiled CLI with no path | not run | not run | usage, exit 2 |
| Compiled CLI with `-ch 2` | not run | not run | usage, exit 2 |
| Compiled CLI with missing `/submission` | not run | not run | missing-path error, exit 2 |

The assertion helper verified exactly nine checks, both supplemental IDs/order,
all seven inherited passes, score/max score 100, each expected lifetime verdict,
and each expected command exit status; emitted
`FULL_DURATION_CONTROL_ASSERTIONS=PASS`. The historical negative failed because
its 60-second deadline discarded the delayed answer, not because of a build
failure, absent traffic, or grader watchdog.

Additional stability and preservation checks:

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test ./internal/gradeedition3 -count=5
gofmt -l cmd/grade-edition3/main.go internal/gradeedition3/ch01.go \
  internal/gradeedition3/ch01_test.go
git --no-pager diff --exit-code -- internal/grade internal/gradeedition2 \
  internal/fakeanthropic cmd/grade-edition2 solutions/ch01 book/edition-3/chapter-01.md
```

Results: five repetitions passed (status 0); gofmt listed no files; protected-path
diff status 0. Source hashes were remeasured and the note's transcriptions checked
against actual tool output. Full repository/later-chapter suites were **not run**;
this bounded check avoids unrelated provider/GUI fixtures and later-chapter work.
No browser or MCP session was started.

New source SHA-256 values after validation:

- `cmd/grade-edition3/main.go`:
  `08ca96d522cafde0d07d6843afe9df4982b002f478b7ba3ea9882b4974ca9ffc`
- `internal/gradeedition3/ch01.go`:
  `a1716c6b11b0531acad0a1c815450a979d2ae382e343162f379285c24d9cf97e`
- `internal/gradeedition3/ch01_test.go`:
  `3d229c5b85c2939ecd2af2969028d02e7f3d0076193371b2c56cb0f6b2b66d2d`

Final source check at 2026-10-16T17:15:25Z still reported HEAD
`2be995ebc12c6a7efb21cee3b9740a758489adda`; the new command/package/note are
untracked additions, with no commit or historical-path diff.

## Anthropic credential and live-access evidence (sanitized)

Codex access was not treated as Anthropic access. Initial checks printed only
presence/nonemptiness of known variables:

- `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `ANTHROPIC_BASE_URL`,
  `ANTHROPIC_MODEL`: absent/empty.
- Inspected loader's generic aliases `LLM_API_KEY`, `LLM_BASE_URL`, `LLM_MODEL`,
  `LLM_VENDOR`: absent/empty. `agent/ensemble/config.go:1-58` documents their
  fallback behavior; this is not evidence of a populated environment.
- Known configured reference `~/.cr/settings.json`: exists, regular file,
  readable. No broad credential search was performed.
- After explicit authorization for schema-only inspection, JSON parsing of that
  one known file succeeded. `directClaudeAPIKey` is a nonempty string. Only
  schema/type/nonempty metadata was printed; no scalar credential value was
  displayed. No credential was written into this repository or command text.

Initially live access was **not run**, pending explicit parent permission. Parent
then explicitly approved one GET to the intended Anthropic service and, if a
permitted supported model was found, one tiny inference (no 5.5, no Codex).
Using Python `urllib.request`, a redirect-denying opener, and the one known local
key field, the exact HTTP operations were:

```text
GET https://api.anthropic.com/v1/models
  x-api-key: [resolved in memory from authorized directClaudeAPIKey]
  anthropic-version: 2023-06-01

POST https://api.anthropic.com/v1/messages
  x-api-key: [same in-memory key; never printed]
  anthropic-version: 2023-06-01
  content-type: application/json
  {"model":"claude-sonnet-5","max_tokens":16,
   "messages":[{"role":"user","content":"Reply OK."}]}
```

Only status, safe model identifiers, usage, stop reason, and reply-match booleans
were emitted. No request/response headers, raw errors, or credential values were
printed. Requests did not follow redirects or retry. No raw authentication file
was displayed. No arbitrary overall generation timeout was passed to the opener.

Observed outcomes:

- GET: **HTTP 200**, model list included `claude-sonnet-5`, no further page.
- POST: **HTTP 200**, resolved model `claude-sonnet-5` (not 5.5).
- Nonempty text: true; reply equals `OK`: true; stop reason `end_turn`.
- Usage: input 12, output 4, cache creation 0, cache read 0.

This proves Anthropic credential-route and minimal inference usability at check
time. The shell environment remains unpopulated: a later authorized live student
check must resolve/inject the known credential securely into its own process,
not assume the grader or student's shell already has it. The local fake grader
does not require this live credential. No additional live requests were made.

## Start decision / remaining work

- Original harness builds and runs Chapter 1 from the root module.
- New bounded overlay builds; focused, repeated, and full-duration
  positive/negative controls pass. It can actually run Chapter 1 against an
  explicitly supplied fresh directory or executable.
- Anthropic access is verified separately from Codex, including a tiny permitted
  `claude-sonnet-5` inference.
- `/submission` is absent: no actual fresh student target was graded or launched.
- Parent must complete independent review/acceptance of the overlay, route the
  execution caller to the Edition 3 command, provision the fresh target, and
  authorize the fresh student run. This checker has not performed that launch.
- Later chapters were not audited by this readiness check. Remaining Chapter 1
  review-only gaps above are explicitly not claimed as automatic coverage.
