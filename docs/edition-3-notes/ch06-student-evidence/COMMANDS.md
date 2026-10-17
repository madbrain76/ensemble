# Chapter 6 student evidence commands and qualifications

Student working directory: `/Users/bill/projects/ensemble-edition-3-student`.
Book grading directory: `/Users/bill/projects/ensemble`.
Candidate identities and binary/archive hashes are in `../ch06-student-source-manifest.json`.

## Required gates on the repaired candidate

Run sequentially, not as a concurrent build campaign:

```sh
go test ./... -count=1 -timeout 240s
go test -race ./... -count=1 -timeout 300s
go vet ./...
(cd examples/ch04-consumer && go test ./... -count=1 -timeout 90s && go vet ./...)
(cd examples/ch06 && go test ./... -count=1 -timeout 90s && go vet ./...)
(cd examples/ch07 && go test ./... -count=1 -timeout 90s && go vet ./...)
```

All binaries were built sequentially with `go build -trimpath`, respectively
`./cmd/ch01`, `examples/ch04-consumer`, `examples/ch06`, `examples/ch07`,
into `/tmp/e3-ch06-frozen-{core,consumer,workflow,shape}`.

From the book repository, all flags precede the explicit same core executable:

```sh
GOFLAGS=-mod=readonly go run ./cmd/grade-edition3 -ch 6 -json \
 -core /Users/bill/projects/ensemble-edition-3-student \
 -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer \
 -workflow /Users/bill/projects/ensemble-edition-3-student/examples/ch06 \
 -stream-exercise /Users/bill/projects/ensemble-edition-3-student/examples/ch07 \
 -stream-model claude-sonnet-4-6 \
 -media-models /Users/bill/projects/ensemble/docs/edition-3-notes/ch06-student-evidence/media-models.json \
 /tmp/e3-ch06-frozen-core
GOFLAGS=-mod=readonly go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch06-frozen-core
```

`final-*` receipts describe repaired candidate 800da21d215b9ca48de41516be4773962fe38ccd.
`candidate1/*`, original `grade-*` and `frozen-*` describe aa704ce00abe4203845e6dff86e4c96105fc7cef,
not the repaired bytes. The first Chapter 6 grade failed (70/100): no assistant
record from the live prompt path, and a default OpenAI media model outside the
student's declared inherited rows. The documented explicit media map selects
unchanged inherited rows; it does not expand their capabilities.

## Failure/assistance ledger supplement

The source tree's committed `evidence/ch06/attempts.md` retains prefreeze attempts.
This supplement records later evidence-tool and protocol corrections:

* The first local shape launcher assumed unrelated JSON field names. Actual shape
  output already used public `assistant/deltas/parts/first_delta_ms`. The harness
  expectation was corrected; its failure output is retained.
* The first live checklist was NOT authorized. Parent required real forwarding
  controls, fixed pins, allowlisted environment, explicit acceptance checks,
  genuine active hints, and stop-on-first-failure; the original launcher remains
  as `live_harness.py`, superseded by `live_harness_v2.py`.
* First v2 rehearsal sent an ordinary prompt rather than the inherited `/hint`
  command, correctly causing a queued turn and cap refusal. Second rehearsal
  admitted the hint but the simplistic fake provider discarded it. Both failed
  launcher versions and receipts are retained in `rehearsal-v2-attempt{1,2}`.
  The corrected fixture carries admitted hint instructions into its next visible
  response/history, and checks actual marker inclusion in later requests/answers.
* My first diagnosis of the failed Chapter 6 completion gate was false: I claimed
  the result handler suppressed assistant JSON in live mode. It already used
  `if !human`. Source inspection and the parent's independent check corrected
  this before any repair. The actual live `kind:prompt/hint` branch posted a
  fire-and-forget UserMessage, bypassing the private completion channel.
* Added `TestLiveProtocolReliableAssistant`: unchanged candidate-1 production
  compiled and failed in BOTH streamed and forced-plain modes with zero assistant
  records despite correct deltas/finals. The repair uses existing reliable submit
  for idle live prompts and explicit Hint for busy/explicit-hint admission. The
  final response handler and actor ownership are unchanged. The CLI suite passes;
  the regression also requires exactly one final answer and durable EOF capture.
* No inherited lifetime failure was blamed on code: an early SIGQUIT and a
  15-second reproducer timeout were student orchestration mistakes. Completed
  127–139 second suites and public 62-second header/body probes supersede those
  incomplete attempts without deleting them.
* `git diff --cached --check` at the first freeze flagged literal whitespace in
  retained unified-diff evidence, not Go source. No evidence was rewritten to
  conceal it. Source-only whitespace checks pass.

## Sensitivity scope

The committed prefreeze controls suppress live deltas, drop valid ID zero, and
replace arrival authorization with finalized-order authorization. Each compiled,
failed its intended public assertion, restored exact source hashes, and passed.
They are scoped to the implementation before the later two-file CLI admission
repair. That repair has its separate unchanged-production fail / corrected-pass
regression, not an alleged grader mutation. No grader source was inspected or edited.

## Launcher boundary

`live_harness_v2.py --rehearse` uses actual frozen binaries with the same local
forwarding boundary as live mode. It verifies source commit/tree/clean status,
manifest/archive/core/shape hashes before process launch and upstream forwarding;
uses a fixed child environment; validates path/model/stream/output cap; reserves
at most three requests under a lock; and stops the five-case series on its first
failed case. Local negatives test wrong pin/model/token/path, concurrent fourth+
requests at the real forwarding boundary, and a first provider failure.
All rehearsal credentials are literal local fixture values. Rehearsal never
enters the settings loader or contacts public provider hosts.

The proposed live command is exactly:

```sh
python3 /Users/bill/projects/ensemble/docs/edition-3-notes/ch06-student-evidence/live_harness_v2.py --approved-live
```

It must not run without separate parent approval of this exact pinned launcher.
The five cases are three human coding/hint exercises (one per declared vendor),
then the OpenAI shape stream/plain pair. Each has at most three requests and
1,024 output tokens/request; total ceiling 15 requests / 15,360 requested output
tokens. It has no retry/discovery path and stops on first failure. The child
process is stopped at 90 seconds; an in-progress upstream socket can require up
to its 60-second timeout to join. Raw stderr/api.log/history stay in private
mode-0700 temporary workspaces and are not portable evidence. Safe exports do
not contain keys/auth/private provider frames. Source and source-manifest bytes
must remain unchanged after approval.

### Approval qualification

Parent `main-28` approved the exact hash-pinned launcher once, conditional on all
final-byte gates passing. See `live-approval-main-28.md`. The wrong-pin negative
proves preflight throws; no-key-loader reachability is established by inspected
main control flow, not by a separately instrumented credential-read experiment.

### Authorized execution result

The one approved command ran after all final gates passed, then stopped with
exit 1 on Anthropic human `hint_in_answer=false`. Actual activity: three HTTP
200 upstream calls, no refusals; child exit 0; config read and exact STREAM_OK
write; genuine early stdout and next-request hint inclusion. Request 3/final
answer omitted the hint marker. Canonical numeric totals: 5,509 input / 244
output, zero cache counts. The four other live cases were NOT RUN. All safe
partial receipts are retained as `live-v2-*`, with `live-once.started` preventing
an accidental second launch. No further provider command ran or is authorized.

### Remaining-only verifier (not a rewrite of the failed live case)

Parent `main-29` required a new local-only rehearsal before any further live
permission. `live_harness_remaining.py --remaining --rehearse` uses only the
four originally unrun cases; Anthropic is absent from its route table and the
selector negative verifies rejection before any upstream server. The human hint
now explicitly asks for the marker in the VERY NEXT response before tools. The
forwarder observes visible provider text order only (not thinking), records a
marker-before-tool boolean/timestamps, and removes transient text suffixes before
export. PASS requires actual next-request marker delivery, next-response marker
before its first tool, and actual CLI stdout in addition to early streaming,
file effects, status/cap/exit checks. This is openly different from the original
final-answer persistence assertion, which remains FAILED and unchanged.

The new local rehearsal and actual forwarding negatives passed on unchanged
frozen production. The exact proposed remaining-only paid command is separately
gated on parent approval; it has no inherited authority from the first command:

```sh
python3 /Users/bill/projects/ensemble/docs/edition-3-notes/ch06-student-evidence/live_harness_remaining.py --remaining --approved-live
```

Its ceiling is four cases × three requests × 1,024 requested output tokens,
with the same watchdog/join bounds and stop-on-first-failure/no-retry policy.
Only the two remaining direct OpenAI/Gemini key fields would be loaded in process
after fixed-pin preflight. No Anthropic key or provider route is selectable.

### Separately authorized remaining execution result

Parent main-30 approved launcher SHA256
ad6d5603e637004567d06b60a392c6ca98e5e2398d7a0fba8993a9cd1503f0c7 once.
The exact `--remaining --approved-live` command ran once and exited 0. All four
cases passed: OpenAI human3 calls, Gemini human3, OpenAI shape stream2/plain2;
10 upstream requests, all HTTP200, no refusals/retries. Next-response-before-tool
hint compliance and actual read/write effects were demonstrated for both human
cases. Usage 8168 input/701 output, zero caches. Original Anthropic3-call failed
assertion remains intact. Combined actual requests13 and usage13677/945 are NOT
the requested output ceilings. All five exercise positions have now been tried;
original Anthropic final-answer compliance remains unsuccessful. Neither live
campaign establishes hidden thinking behavior or byte-identical stochastic replies.


### Metadata-only inventory correction after initial submission

Parent requested reclassification of `candidate-source.tar` and
`candidate1/source.tar` as LOCAL-ONLY, regardless of UTF-8 decodability, and
inclusion of the four existing candidate-1 frozen binaries from
`candidate1/preserved-artifacts.json` as LOCAL-ONLY. Paths and bytes remain
unchanged; no artifact was moved or deleted. Prior report/index are retained by
the parent at `/tmp/cadence-ch06-pre-inventory-correction` (untouched). Prior
report SHA256 is 7022fcc7331d04a7a6cd21f770aebe86e81e35234612384b82602b8261252429;
prior index SHA256 is a45b8c21ab677dc892622ab522dcf452e95c9d0b0ac9b1051bf622491238f0c4.

Updated only report/COMMANDS/index-builder/index metadata; regenerated the index
and checked hashes/counts against the previous indexed set and preserved map.
Counts: 474 total, 283 portable, 191 local-only (previously 470/285/185).
No implementation/commit changes, tests, builds, provider calls, credential reads
or optional evidence expansion. This is not a formal code-correction round.
All original live outcomes and paid receipts remain unchanged.
