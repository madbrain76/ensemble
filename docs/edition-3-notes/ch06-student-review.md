# Edition 3 Chapter 6 — student implementation review

## Identity and status

This is a student implementation, not independent acceptance or an acceptance tag.
Accepted predecessor: `f60ebf3aad2ece4251face9d95c106dd667abcd6`, annotated
`edition-3-ch05-r1`, tree `3de5d0d24dd9f593ec017fceccfbff1b209c5d57`.
Assignment SHA256 `0d202210609993531539dcd4d11edf9aad7ebecf8893ea3f22cb91702c383a1a`;
public chapter SHA256 `ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55`.
Book revision and supplied adapter hash, final commit/tree/archive and all four
actual executable hashes are recorded in `ch06-student-source-manifest.json`.
The repaired source is clean at `800da21d215b9ca48de41516be4773962fe38ccd`.
The first candidate, failed grade, and genuine passing lifetime evidence remain
preserved separately under `ch06-student-evidence/candidate1`.

## What I implemented

* Kept the existing Engine, actor, History, root job service, and CLI composition.
  `Engine.Parse(*http.Response, StreamCallbacks) error` is the single Parse seam.
  SendPrepared transfers the live response and closes it after parsing; plain JSON
  and provider SSE use the same final semantic decoder and normalized usage rules.
* Added SSE framing and provider accumulation for Anthropic, OpenAI, and Gemini.
  Stable IDs are allocated on first part arrival (including valid zero), reset
  per response, and carried through worker completion to model finals. Provider
  opaque replay/signatures are retained in canonical finals, not partial history.
* Deltas cross a turn/request-fenced actor mailbox. Complete response state and
  usage remain reliable worker facts, saved before model-final publication.
  Explanation authorization uses visible text before the first tool in actual
  arrival order AND finalized order; thinking does not authorize tool dispatch.
* Configured-key redaction precedes frame forwarding, root trace output is
  serialized, and callback values are copied. Request cancellation uses the
  caller/HTTP lifetime without an arbitrary generation deadline.
* A CLI display owner queues observations off the actor and joins producers before
  draining at EOF/error/interrupt. Human reply text goes to stdout exactly once;
  optional reasoning/tool styling stays on terminal stderr. Hint routing uses
  origin metadata, not text prefix. `EN_DISABLE_STREAMING` is the core opt-out.
* Added the separate public-library `examples/ch07` module, leaving the inherited
  Chapter 5 collaboration module intact. It counts model-response finals and
  delta kinds and measures first-delta latency; `CH07_NO_STREAM` selects plain mode.
* Preserved reliable JSON completion alongside streaming observations. The first
  candidate missed this for live `kind:prompt`: a fire-and-forget admission path
  bypassed the private reply. The narrow repair now uses the existing reliable
  submit path for an idle live prompt and explicit Hint for busy/explicit hints.

## Evidence and current limitations

Detailed commands, failed attempts, assistance and sensitivity scope are in
`ch06-student-evidence/COMMANDS.md` and committed `evidence/ch06/attempts.md`.
The final same-byte test and grading results and any separately authorized live
results are appended below only after their jobs have ended.

Coverage includes delayed public Parse/Ask callbacks, a real gated CLI flush of
legitimate `hint:early` text before HTTP continuation, provider canonical parity,
zero/reset IDs, fragmented arguments, malformed/EOF/read failures, usage known
zero, cancellation/recovery, trace redaction, callback mutation isolation, tool
arrival authorization, private replay, opt-outs and slow display drain. Inherited
root suites exercise jobs, interruption, budgets, reliable replies and lifetime.
The three retained mutation controls test specific behavioral assertions; they
are not evidence that every parser/CLI path was independently audited.

No paid-provider behavior is inferred from local fixtures or the deterministic
scores. Raw provider logs (if any are later authorized) remain private/local;
portable evidence reports safe observations and honest failed/unavailable cases.

## Assistance, mistakes and teaching clarification

I read only the allowed public prerequisites and my accepted/current student
source. I used no child agents, private parent memory, answer/reference source,
grader source, or future chapter implementation. Parent approved the ownership
plan and later challenged the live harness and an inaccurate diagnosis; those
interventions improved the boundary controls rather than supplying solution code.

I interrupted one intentionally long inherited test and used an invalid short
reproduction timeout; both are recorded as my orchestration mistakes. Later
complete suites and real delayed-header/body gates are the valid lifetime evidence.
Several mechanical edit/compile and harness assumptions failed and are retained.
The initial Chapter 6 grade correctly exposed missing reliable assistant output.
My first claimed cause was wrong and was corrected by direct source inspection
and the parent's check before the narrow repair. The default media-model grade
failure is preserved; the final documented map selects unchanged inherited rows.

The most useful lesson was that streaming progress, canonical final content and
reliable request completion are three different obligations. A UI can display
perfect deltas yet still lose the answer contract by admitting a request through
a fire-and-forget path. Likewise, final provider ordering cannot prove a visible
explanation arrived before a tool call. Local gates which hold the server open
and inject actual behavioral defects made those distinctions observable.

## Completed gates on final repaired bytes

All below ran on clean `800da21d215b9ca48de41516be4773962fe38ccd`; stdout,
stderr and separate exit receipts are indexed under `final-*`.

| Gate | Actual result |
| --- | --- |
| Root `go test ./... -count=1 -timeout 240s` | exit 0; root package 127.764 s |
| Root race `go test -race ./... -count=1 -timeout 300s` | exit 0; root package 139.342 s |
| Root `go vet ./...` | exit 0 |
| `ch04-consumer`, inherited `ch06`, new `ch07` standalone test + vet | all exit 0 |
| External final shape with local delayed SSE and plain JSON | both exit 0; same `shape one two` answer |
| Edition 3 Chapter 6 route, explicit stream/media models | exit 0; `passed:true`, 100/100 |
| Edition 3 Chapter 1 route, same executable | exit 0; `passed:true`, 100/100, both required zero-point lifetime checks true |

The final local shape stream measured two text deltas, one finalized model part,
and first delta 1.331 ms vs 158.195 ms total; plain mode measured one final-part
delta, the same answer and one final part. These are local fixture measurements,
not provider latency claims. The public Chapter 6 gate saw 36 text deltas across
two parts, 54 thinking deltas excluded from reply text, and 14 tool fragments;
stream/plain comparison had the same 55-byte answer and five final part kinds,
with 104 versus five deltas. The report preserves its inherited title; the actual
command was the authorized `cmd/grade-edition3 -ch 6` route.

The Chapter 1 grade completed both actual one-minute-two-second header and body
waits and received complete answer/usage. It also passed protocol, wire, call
count, reply, memory, append-only history and usage checks. These completed
same-byte results—not the earlier aborted attempts—are the lifetime evidence.

The repaired-byte local forwarding rehearsal passed all three real human CLI
exercises and the OpenAI shape pair, using only fake upstream servers and fixture
credentials. Negative controls exercised wrong pin/model/token/path, synchronized
fourth-plus reservations, and stop-on-first-provider-failure. Parent `main-28`
then approved one finite live execution conditional on these completed gates.

## One approved live attempt — partial coverage, failed assertion

After all final-byte gates passed, I ran the exact launcher authorized by parent
`main-28` once. It stopped after its FIRST case and returned exit 1. This is not
a passing live campaign and I did not retry, substitute a model or continue.

Actual provider activity was **three Anthropic claude-sonnet-4-6 requests, all
HTTP 200**, with no refusals. The human CLI process exited 0, read the fixture
`config.json`, and wrote `output.txt` with exact bytes `STREAM_OK`. Canonical
usage totaled **5,509 input / 244 output tokens**, zero cache-read/write; each
request's configured 1,024 output cap is a limit, not reported usage.

First substantive stdout was measured at 1,134.636 ms, before the first forwarded
response ended at 1,862.455 ms. The hint was sent at 1,134.638 ms and acknowledged
at 1,134.714 ms while forwarding was active. Its marker was present in actual
request 2, proving admission to the next prepared context, but absent from
request 3 and the model answer. Therefore the launcher's `hint_in_answer` check
was false even though early streaming, next-request hint delivery and actual tool
file effects were true. I preserve that failure, rather than reinterpret the
stricter answer requirement as success or infer its cause from hidden reasoning.
The accepted hint design is transient; I have not changed that ownership/lifetime
contract to make a paid observation pass.

**Not run at the original campaign stop:** OpenAI human, Gemini human, OpenAI shape streaming, and OpenAI
shape plain. Their successful local fixture rehearsal is separately labelled and
is not live evidence. The stop-on-first-failure rule left these four cases
uncovered. No further paid calls are authorized or claimed.

Safe request/status/timing/usage/file-effect records and stdout are under
`ch06-student-evidence/live-v2-*` and `live-observations-safe.json`. Raw provider
frames, stderr and canonical history are preserved only in the private temporary
workspace, listed local-only by the artifact manifest. Credentials were loaded
only by the approved launcher after fixed-pin checks and were not supplied to
children, tools, command arguments, fixtures or portable evidence.

At this original stop point, local/public gates passed but the first paid exercise
was partial/failed. The separate remaining-only section below records the later
explicitly authorized completion of the four unrun positions. The original failed
assertion is retained; this report does not declare Chapter 6 accepted.

The exact paired-grade cwd/argv, `GOFLAGS=-mod=readonly` override, selected
media-map bytes/hash, same core binary hash, and each stdout/stderr/exit hash are
recorded in `ch06-student-evidence/paired-grade-command-receipt.json`. It records
completed tool job 1703 after execution; no grading rerun was performed to create
that provenance file. The human-readable exact command block is in `COMMANDS.md`.

## Separately authorized remaining live exercise — four cases completed

The preceding statement that four cases were unrun describes the stop point of
the ORIGINAL campaign, not the final coverage state. Parent `main-29` required a
new, explicitly remaining-only local rehearsal. Parent `main-30` then approved
its exact hash-pinned command once after the actual frozen binaries and forwarding
boundary passed that rehearsal and its no-Anthropic/stop-first-failure controls.
No production bytes changed and no Anthropic repeat was made.

The new human hint explicitly requested its marker in the **very next response,
before tools**. This intentionally tests the accepted transient hint's next-context
behavior rather than persistence through another tool cycle. The verifier requires
marker delivery in request 2, visible marker before that response's first tool
frame, and actual CLI stdout. It excludes private thinking from the text check.
The original final-answer assertion remains failed; it was not relabelled or
rerun under the changed verifier.

The remaining-only launcher completed with **exit 0 and all four cases passing**:

| Case | Actual upstream calls/statuses | Observed result |
| --- | --- | --- |
| OpenAI gpt-4.1-2025-04-14 human | 3, all HTTP 200 | Config read and exact STREAM_OK write; active hint in next request, next visible response before tool, and actual stdout |
| Gemini gemini-3.8-flash human | 3, all HTTP 200 | Same coding/file-effect and next-response hint checks |
| OpenAI shape streaming | 2, both HTTP 200 | 56 text + 7 tool deltas, 0 thinking; 3 model finals; first delta 1084.638 ms, total 2321.767 ms |
| OpenAI shape CH07_NO_STREAM | 2, both HTTP 200 | 2 text + 1 final tool delta, 0 thinking; 3 model finals; first delta 1457.792 ms, total 2196.851 ms |

OpenAI human first substantive stdout was 1865.963 ms, before first-response end
2449.601 ms; hint acknowledgement was 1866.157 ms. Gemini first content was
3150.497 ms, before first-response end 3319.763 ms; acknowledgement 3150.624 ms.
Both actually executed the read/write tools. There were no refusals or retries.
The shape replies both reported `stream-study` and setting `7`; punctuation and
usage differed, as expected for separate live generations. I do **not** claim
byte-identical live output or a causal provider speed benchmark. Canonical byte
parity is established by the controlled local fixtures, not these stochastic runs.
The shape prints its summary at completion; its `first_delta_ms`, not summary
stdout arrival, measures the observer's first chunk.

Remaining-run canonical usage was 8,168 input / 701 output tokens, zero cache
counts. Across BOTH approved commands, actual totals were **13 provider requests,
all HTTP 200; 13,677 input / 945 output tokens**, with zero cache counts. The
first command still exited 1 on its stricter Anthropic assertion, while the
remaining command exited 0. No provider case is now unattempted, but the original
Anthropic final-answer compliance check is explicitly unsuccessful. Actual live
thinking-delta coverage is not claimed: the shape observed zero thinking deltas;
provider thinking/replay behavior is covered by the labelled local fixtures.

Safe final receipts are `remaining-live-*`, including exact counts, statuses,
timing, observed file effects and summary values. The source and executable pins
are unchanged. All authorized provider/fixture/test/build jobs have ended. Raw
private artifacts remain local-only and are hash-indexed; no additional paid
calls, acceptance tag, push, install or service restart was performed.


## Bounded inventory packaging correction

After the initial submission, the parent requested a metadata-only correction,
not a formal code-correction round. The parent-held prior report/index boundary
remains at `/tmp/cadence-ch06-pre-inventory-correction`, untouched. Prior report
SHA256: `7022fcc7331d04a7a6cd21f770aebe86e81e35234612384b82602b8261252429`;
prior index SHA256: `a45b8c21ab677dc892622ab522dcf452e95c9d0b0ac9b1051bf622491238f0c4`.

The evidence copies `candidate-source.tar` and `candidate1/source.tar` are now
**local-only**, regardless of whether their bytes can decode as UTF-8. Their
paths, sizes and hashes did not change. The four already-retained candidate-1
frozen binaries named in `candidate1/preserved-artifacts.json` are now explicitly
included in the main inventory as local-only, with their existing hashes verified.
Nothing was deleted or moved. The inventory changes from 470 entries (285
portable / 185 local-only) to **474 entries (283 portable / 191 local-only)**.

Only this disclosure, COMMANDS, the inventory builder and generated index changed.
Production source, commit/tree, archives, binaries, failed original live assertion
and all paid receipts are unchanged. No tests, builds, provider calls, credential
reads or optional evidence expansion were performed. The original live failure
and later separate successes retain their distinct outcomes. No jobs are pending;
context is retained for the forthcoming independent review findings.
