# Chapter 3 student candidate review

## Status and identity

**Student implementation submitted for independent review, not accepted. Gemini
human-answer acceptance is incomplete.** Both required final graders and the
final race/coverage tests pass, but a concrete inherited Gemini signed-text
classification defect was discovered in actual live use and remains open.
No acceptance tag was created and nothing was pushed.

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Candidate: `0eb6fdffe62e662b0ab96982b01287afba19ff42` (clean).
- Accepted parent: `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`, also the peeled
  `edition-3-ch02-r1` target.
- Main repository revision used for grader execution:
  `f7530a6643ca00bd2f59a8da60021c567cd7d174`.
- Absolute executable, **outside Git**:
  `/tmp/ensemble-ch03-evidence/student-cli`.
  SHA-256 `e2ce182d3220ffa05c3584b9c92ef84536925ffbbe0f88f5bccba9ac29b2aa2c`.
- Exact committed source archive, **outside Git**:
  `/tmp/ensemble-ch03-evidence/ch03-student-source.tar.gz`.
  SHA-256 `5cc91e15b10c11f95ff2b5ff89a1ca66bd3874dc68c79e9c8b53e7517ef1e3e8`.
- Toolchain: Go 1.25.4, darwin/arm64.

The final binary was compiled from the final working-tree bytes and those bytes
were then committed unchanged. Consequently its embedded Go VCS metadata reports
the predecessor plus `modified=true`, not the subsequently created candidate
commit. `identity.json` and `final-binary-buildinfo.txt` make that distinction
explicit; this is not a binary rebuilt after the receipt was recorded.

All evidence below is under `ch03-student-evidence/` beside this report. Read
`identity.json`, `commands.json`, and the actual grader `passed` fields and
`.status` files rather than inferring success from score alone.

## Inputs, boundaries, ownership and assistance

Read the required AGENTS, coding skill, coder overlay, public student brief,
Chapter 3, and the additionally required shared execution procedure; bounded
rereads recovered truncated material. Inspected only the accepted student source
and my own new source/evidence. Did not inspect main-repository `agent/`,
solutions, grader source/tests, future chapters, private prior student/reviewer
records, or parent context. Public grader commands were invoked without reading
implementation. No subagents were launched. No GUI, runtime installation,
restart, push, manuscript edit, or OS-sandbox claim was made.

Before implementation I sent the ownership/configuration/test plan and waited
for parent approval. The parent specifically corrected the plan to all **six**
tools, distinct command stdout/stderr/exit status, ordinary environment
preservation except supported provider credentials, and caller-owned lifetime
without an arbitrary command deadline. Later assistance authorized one finite
human-chat session per provider after protocol live use, and required the final
Gemini issue to remain open rather than silently changing caps/models/retrying.

Ownership stays `Ensemble -> Agent -> Engine/History/Tool`. The root composes
optional Agent-owned `internal/tools` instances. Agent owns whole-batch policy,
dispatch, response/result recording, and the ten-request round budget. Tool
instances own synchronous filesystem/command operations. Engine owns typed
provider schemas and wire rendering; History remains authoritative. No new
framework, streaming/job system, GUI, sandbox subsystem, public call bypass,
or model-name convention was added. The existing root public library and
`./cmd/ch01` CLI layout remain intact.

Public selection is explicit: `Config.EnableTools` installs the six built-ins;
ordinary library construction leaves the registry empty. The real CLI enables
them. `Config.RequireToolExplanation` selects the explanation requirement for
that Agent's selected model. CLI `LLM_REQUIRE_TOOL_EXPLANATION=true` selects it;
false or unset clears it, and invalid booleans fail. This is not a universal
policy, hardcoded vendor rule, hidden grader mode, or private model-name rule.
The README documents both public seams and the six argument contracts.

## Implemented behavior and direct contract proof

`chapter3_test.go` exercises the public package with real tools and faithful
HTTP boundary fixtures. `cmd/ch01/main_test.go` explicitly selects/clears the
CLI policy. Important tests are:

- `TestChapter3SelectedPolicyWholeBatchAndRetry`: all three provider wires;
  opaque reasoning is preserved but does not authorize tools. An unexplained
  multi-call batch causes no file/command side effect; every attempted call gets
  a paired error result. The next request carries feedback and preserved model
  response. A genuinely explained retry then executes actual side effects,
  followed by a terminal answer. This proves more than a mocked dispatch count.
- `TestChapter3PolicySelectionAndPrecedingText`: public bool both enabled and
  disabled; whitespace and text after the first call do not authorize a batch.
- `TestChapter3RoundBudgetAndErrors`: repeated explanation refusals and ordinary
  unknown tools consume the existing ten-round budget, without a separate retry
  allowance. Missing-argument errors are ordinary feedback.
- `TestChapter3RegistryAndSchemas` and
  `TestChapter3NativeResultIDsAndEmptyRegistryFeedback`: six authoritative
  declarations, no declarations for empty registry, provider-native pairing,
  Gemini IDs, feedback with no executable tools, and Gemini's supported schema
  subset. Absent IDs receive stable Agent execution IDs; supplied IDs survive.
- File tests cover sorted directory discovery; one-based bounded reads and
  byte-cap markers; missing/zero defaults; guarded write/overwrite/append;
  exact/overlapping ambiguous/no-match edits and empty replacement; regex
  matches, context-window merging/separators, and actual unchanged bytes after
  refusal.
- Command tests use disposable CWDs, separate streams/status, normal nonzero
  exit as data, fresh shell state, caller cancellation, ordinary inherited
  variables, and benign credential canaries. Provider credentials do not reach
  tool subprocess environments. This is environment filtering, not filesystem
  isolation or an OS sandbox.

The response text requirement is about visible normalized `text` before the
first call, not opaque thinking or text from an earlier response. Whole-batch
preflight happens before *any* operation. Complete model parts are stored before
paired results; retry traffic comes from History's actual projection.

## Verification and sensitivity

From the student root, ordinary tests and race tests passed on behaviorally
final source. After comment-only documentation corrections, the exact committed
bytes passed:

```
go test -race ./... -count=1 -coverpkg=./... -coverprofile=/tmp/ensemble-ch03-evidence/final-coverage.out
```

Exit 0 (`final-race-coverage.txt`). The intentional two 62-second lifetime phases
were allowed to finish. `go vet ./...` also passed. The initial coverage command
without `-coverpkg` misleadingly reported uninstrumented internal packages as
zero; that evidence is retained and was superseded, not presented as meaningful
coverage of those owners.

The combined final profile reports **90.6% statements**, with Tool `Execute`
93.3%, command `run` 96.0%, Agent `Ask` 89.1%, and provider render functions
94.6–100%. These numbers are diagnostics, not a numeric acceptance gate.
Malformed-input/I/O/cancellation and inherited provider branches remain partly
uncovered. More importantly, the Gemini signed final-text semantic case was
*not* covered despite high coverage; the live defect demonstrates that limit.

Both following commands ran from `/Users/bill/projects/ensemble`, on the same
final absolute binary identified above:

```
go run ./cmd/grade-edition2 -ch 3 -json /tmp/ensemble-ch03-evidence/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /tmp/ensemble-ch03-evidence/student-cli
```

Both report **`passed:true`, exit 0** (100/100). Chapter 3 includes Chapter 2;
the separate Edition 3 Chapter 1 run preserves its lifetime gates. No unsupported
Edition 3 `-ch 3` command was used. Grading alone is not explanation-gate proof.

`mutations-final/index.json` records **17 independently failing mutations**, each
restored to exact original bytes with its targeted test passing afterward; none
was a compile-only failure. The final sweep ran in
`/tmp/ensemble-ch03-mutations-final`, not the candidate repository. Mutations
covered gate bypass; opaque/blank/late-text acceptance; skipped dispatch;
corrupted result IDs; extra round; missing declarations; unsupported Gemini
schema field; lost Gemini result ID; overwrite, read cap, search context and
ambiguous-edit safety; credential filtering; nonzero command classification;
and CLI policy clearing. Each mutation has separate mutated/restored output.
Earlier invalid anchors/test names and an insensitive first harness attempt are
retained under `mutations/`; they are not counted as passing sensitivity proof.
The initial mutation harness restored in-place source before I corrected the
workflow to scratch-only; final evidence uses the isolated sweep.

## Actual CLI/live exercise

`live-plan.md` was written **before** provider calls, mapping features to steps;
it was updated before the parent-approved human sessions. Current model-list
requests succeeded (`live-discovery.json`), then the actual CLI exercised:

| Provider/model | Protocol + real `chat` sessions | Verified artifact outcome | Human final answers |
|---|---|---|---|
| Anthropic `claude-sonnet-4-6` | all six tools; two human turns each | unchanged repair test goes 1→0; follow-up marker persists | both nonblank |
| OpenAI **ChatCompletions** `gpt-4.1-2025-04-14` | all six tools; two human turns each | unchanged repair test goes 1→0; follow-up marker persists | both nonblank |
| Gemini `gemini-3.8-flash` | initial HTTP400; one schema-repaired protocol retry; one `chat` session | retry/chat use all six tools, tests 1→0, marker persists | **both empty in each successful tool session; incomplete** |

Every session used a distinct disposable project identified in its `live-*.json`.
The task repaired actual `Add` arithmetic without changing its test, wrote a
repair note, then used a second human turn to read/append the note and rerun the
test. Directory discovery, file reads/search/edit/write, real shell execution,
tool feedback and follow-up conversation are evidenced in normalized History.
All observed call IDs are paired with returns. Command results show the original
failure and subsequent successes `[1,0,0]`. Independent launcher tests verified
disk state, rather than accepting model claims. The OpenAI protocol session also
contains one **real live policy refusal followed by an explained successful
retry**; see its history summary. No automatic retry count is hidden.

The protocol launcher had a 240-second observation watchdog; human `chat` had
120 seconds per provider, outside the application. Completed sessions did not
hit those limits. Human sessions passed plain-text input through the real
`chat` loop, not the grading JSON-lines interface. No live GUI was used.

Only the three authorized credential fields were loaded by local launchers into
process memory. No credentials were passed as command arguments, included in
model prompts, printed, committed, or copied to evidence. Own-source/evidence
secret scanning emitted only a boolean (`secret-scan.json`: no matches).
Raw normalized logs remain in the identified disposable projects; committed
sanitized summaries contain field names, IDs, lengths, hashes, booleans and
outcomes, not prompts, raw arguments, model/hidden text, or provider keys.

Live binaries are explicitly **not falsely equated** with the final hash:
protocol Anthropic/OpenAI and Gemini's first failed attempt used hash
`60bbd2a138efb2c10d868b074c2aba544c0e8a18e01dd57b2c6e6ace2b079164`;
Gemini's schema-repaired protocol retry used
`f1f112397474f881027af6639cbb4135417786e8e52114db1bd5ef1c852afb71`;
all three final human sessions used
`588c6069397e10a0f272b8416e9b78a9279461264916b6d2033f833782cbe9f9`.
The latter differs from final source only by Go documentation comments.
`human-live-to-final-comments.patch` retains a comparable reconstructed delta;
`live-build-attribution.json` explicitly says this is reconstruction, not a
contemporaneous archive. Earlier protocol runs precede the evidenced Gemini
schema repair and are not credited as final-Gemini-interface proof.

## Preserved failures and open Gemini defect

Failures were recorded before repairs, including initial scripted anchor/build
errors, a mistaken OpenAI part-order test expectation, an incorrect CLI config
function test name, the first Chapter 3 grader contract failures, the second
remaining edit-file diagnostic failure, invalid mutation harness attempts,
and Gemini's initial live HTTP400. See `initial-patch-failure.txt`,
`second-patch-failure.txt`, `first-contract-test-failure.txt`,
`cli-test-name-failure.txt`, `grader-edition2-ch3-attempt1.{json,status}`,
`grader-edition2-ch3-attempt2.{json,status}`, and `gemini-schema-diagnostic.json`.
The initial tool-contract failures were my omissions of explicit §3.3 details,
not evidence that the chapter omitted those requirements. They were repaired
from the public table, not grader internals.

Gemini initially rejected `additionalProperties` inside `parameters`. A narrow
actual-rendered-request diagnostic established the error; the correction uses a
typed Gemini-supported schema subset while keeping local argument validation
strict. The first failure is retained as `live-gemini-attempt1.json`.

**Open defect:** terminal Gemini signed non-thought text is suppressed locally.
`gemini-terminal-metadata.json` records four terminal blocks (protocol and chat,
sequences 31/50): each contains raw keys `text` and `thoughtSignature`, nonblank
text (874/591/763/570 bytes), and **no `thought` flag**. No hidden text was emitted.
The unchanged inherited parser at `internal/engine/parse.go:240-243` treats
`block.Thought || len(block.Signature)>0 || block.Text==nil` as opaque. Agent
`internal/agent/agent.go:145-159` concatenates only normalized text and returns
`"", nil` if there are no calls. Protocol `cmd/ch01/main.go:171-178` emits an
empty assistant string; human loop `:203-214` prints a blank line and ultimately
returns success. The actual CLI exits were 0, **not errors**. Tool work and the
follow-up really happened, but this is not complete human-answer acceptance.

Normalized History does not retain provider finish/stop reasons, so none is
claimed; no token-budget explanation is supported by this evidence. The
signed-text classification path, however, is directly evidenced in existing
stored parts. Per parent direction, the candidate stays frozen for independent
review rather than silently changing inherited parsing, caps, or retrying live
requests. Acceptance must resolve this issue; passing graders do not erase it.

## Reader feedback (student experience, not an author rewrite)

- **§3.3 tools/table:** sufficient and concrete. I initially failed to implement
  existing-file protection, max-byte read behavior, grep context, and normal
  nonzero status correctly. Rereading the table resolved these; this is coder
  error, not a request to add requirements or expand the chapter.
- **§3.6 provider declarations/Gemini:** live Gemini rejects the registry's
  `additionalProperties` in its `parameters` Schema dialect. A concise reminder
  that equivalent contracts may need provider-supported schema representation
  would have prevented my direct schema reuse. The measured interface is
  `gemini-3.8-flash`; I make no claim about every Gemini endpoint/version.
- **§3.5 and §3.8 explanation gate:** the whole-batch preflight, visible-before-
  first-call ordering and same-budget retry contract were actionable. Public
  explicit selection/clearing tests are materially better than model-name
  guesses. The real OpenAI refusal/retry reinforced this; no new feature needed.
- **§3.8 live-test checklist:** retain actual disk/test verification and human
  follow-up, not merely exit status or tool calls. This exposed a signed-text
  parser issue that the grader and high coverage missed. Any editorial response
  should wait for independent code acceptance and the signed-text semantics
  review; I do not propose a new API or manuscript rewrite here.

Independent acceptance, any bounded correction, and later author/text review
belong to the parent. This is a finite student candidate with a precisely
identified open defect, not an acceptance declaration.
