# Edition 3 Chapter 5 — student candidate report

## Status and provenance

**Candidate, not accepted.** Parent owns review, acceptance, and tags. No acceptance tag was created. No paid provider call has been made at the time of this initial report section. Later results are appended below.

- Launch task: `/tmp/cadence-ch05-student-task.txt`, SHA256 `b5e56c262865f7d7af4a2067bb864788371a7b90c7f0d00b87c1a65c51db8646`.
- Current chapter: `book/edition-3/chapter-05.md`, SHA256 `5d3db35e29da8a7fa8b343c01c88b3040339d75c347e845f434f124f3b17b0f2`.
- Exact predecessor: `ec4a350f564e820d398f9125fd73a9b3b27129cd`, peeled `edition-3-ch04-r1`.
- Frozen student source: `/Users/bill/projects/ensemble-edition-3-student`, clean commit **`00489b05b65cfe5209f26cb655667df44b3847e2`**.
- Main/public-grader revision at resumed validation: `b9304c266ed50c0fd9234d5eaca80b486ca287b3`. Parent confirmed its advance from the launch HEAD was permitted; task and chapter hashes were unchanged. No private grader/audit/reviewer implementation, historical answer, final reference, or parent memory was read.
- Ordinary source snapshot: `/Users/bill/projects/ensemble/.agents/chapters/code/ch05-student`, verified against every frozen source hash.
- Source archive: `/tmp/cadence-ch05-evidence-r1/frozen/source.tar`, SHA256 `336b26e12c16c268a1634dd690d4fc9be075aeed4b1f323ff36c42a1436f70df`.
- Race-built CLI: `/tmp/cadence-ch05-evidence-r1/frozen/ensemble`, SHA256 `258233df606e71d50aa0103c5865fcd2e383efc471c3d3f8e71fe59f07609863`.
- Race-built workflow: `/tmp/cadence-ch05-evidence-r1/frozen/workflow`, SHA256 `c5c2fee6aa8c463853c1c6210f7b2a699be8bf6806e550ba0b28d0aa44507114`.
- Compiler: `go version go1.25.4 darwin/arm64`. Both frozen binaries record the candidate revision and `vcs.modified=false`; see `frozen/build-info.txt`.
- Complete identity: `/tmp/cadence-ch05-evidence-r1/frozen/identity.json`, copied to `docs/edition-3-notes/ch05-student-manifest.json`. Earlier pre-commit binaries and manifests remain intact and are not described as these frozen binaries.

## Implementation and ownership

| Requirement | Files / owner |
|---|---|
| Shared typed input, reliable result, observation, media metadata vocabulary | `internal/common/actor.go`, `internal/common/common.go`; public aliases in `ensemble.go` |
| One Agent mailbox, copied admission, serialized history commits and transitions | `internal/agent/actor.go`, `agent.go` |
| Generation-tagged model/tool work, pending hints, independent Ask requests, atomic configured history save before reply | `internal/agent/actor.go` |
| Job launch versus asynchronous observation; one-shot limits and call/result pairing | `internal/agent/dispatch.go` |
| Hint, turn-end, late job-report reduction | `internal/history/history.go` |
| Admission closure, ordinary Go-handler joining, managed process stop/join, retained output | `internal/jobs/jobs.go`, `internal/common/jobs.go`, root shutdown in `ensemble.go` |
| Explicit model rows and native image/audio/video/PDF request branches | `internal/engine/models.go`, `render.go`, `engine.go` |
| Long-lived reader/select loop; strict live JSONL commands, human hints/interruption, EOF drain | `cmd/ch01/main.go` |
| Application-owned author/editor/reviewer roles, tools, manuscript and reliable handoffs | standalone public-importing `examples/ch06` module |
| Focused actor/media/shutdown tests and workflow test | `chapter5_test.go`, `examples/ch06/main_test.go` |
| Caller-facing boundaries | `ACTORS.md`, `examples/ch06/README.md` |

The actor is the only conversation-transition owner. Provider rendering/loading/HTTP and tool observation run off-loop. Workers report immutable generation-tagged outcomes. Interrupt pairs current/unstarted calls and ends the current request, but detaches observation rather than killing the actual retained job. A new turn can complete while that job is unfinished; later job reports remain durable but cannot resume the interrupted turn. Root shutdown closes every Agent's admission, stops/joins managed processes, waits for ordinary handlers, drains reports, then saves final logs.

`Post` is admission, not completion. UserMessage and Hint are classified against current actor state: idle/interrupted starts a prompt; active work receives one-request pending context. `Submit`/`SubmitMessage`/`Ask` have independent buffered reliable replies. Observer callbacks execute after reduction on the owner and must not block; `Wait`/`WaitAny` queues are bounded and lossy. Workflow handoffs use Ask replies, not observations. Whole-part deltas are supplied; token streaming is not claimed.

Batch limits count admitted dispatched batches, default 200. The next over-budget tool batch receives synthetic paired failures. Explanation-policy refusals are separately bounded and do not consume executed-batch budget. The live evidence forwarder therefore enforces its own small hard request ceiling rather than treating the tool cap as a billing guarantee.

Media loading clones request context before reading stable local refs; bytes do not enter persisted history. Redacted parts are skipped before file access. Explicit rows reject unknown or unsupported combinations. Anthropic supports image/PDF; OpenAI's selected row supports image; Gemini's selected row supports image/audio/video/PDF. Local files are bounded at 20 MiB. URI/opaque-handle fetching is not implemented. Only an authorized positive image run will be claimed as live; audio/video/PDF are fake-render/HTTP evidence.

Library callers may leave HistoryPath empty for memory-only use. Both executables configure persistence. Arbitrary ordinary Go handlers cannot be forcibly stopped and can delay shutdown indefinitely; that limitation is explicit rather than hidden behind a synthetic successful cancellation.

The GUI remains a stub. Only its standalone dependency metadata changed so it can build the inherited PTY dependency. No GUI installation or live session was attempted.

## Local validation and sensitivity

Fresh sequential gate receipts: `/tmp/cadence-ch05-evidence-r1/resume-r2/gate-exits.txt`. All exited 0:

- Root: `go vet ./...`; `go test -race -count=1 -timeout 4m -coverpkg=./... -coverprofile=/tmp/cadence-ch05-evidence-r1/resume-r2/coverage.out ./...`; `go build -race ./...`.
- Each of `gui`, `examples/ch04-consumer`, `examples/ch06`: `GOWORK=off go vet ./...`; `GOWORK=off go test -race -count=1 -timeout 90s ./...`; `GOWORK=off go build -race -o <receipt-recorded-output> .`.
- Frozen builds: root `go build -race -o /tmp/cadence-ch05-evidence-r1/frozen/ensemble ./cmd/ch01`; workflow `GOWORK=off go build -race -o /tmp/cadence-ch05-evidence-r1/frozen/workflow .`.

Root statement coverage is **85.5%**, including internal packages (`resume-r2/coverage-summary.txt`). No source changed between the full gate compilation and the candidate commit. Frozen executable rebuilds additionally establish clean VCS metadata.

| Acceptance property | Deterministic evidence | Valid mutation rejection |
|---|---|---|
| Held-model hints remain responsive and feed the next request exactly once | `TestChapter5Hints…` | dropping hint reduction |
| Per-request replies survive dropped observations and follow configured persistence | hint/reliable test with 140 observations | skipping request save |
| Interruption returns promptly; new turn completes before old handler release | `TestChapter5Interrupt…` held channels | killing ordinary handlers during shutdown |
| Late actual job facts survive without resuming old turn | same test; saved final history and retained output | dropping late job reports |
| Configurable exact tool batch boundary | `TestChapter5Configured…` | bypassing batch-cap branch |
| Positive image/audio/video/PDF request branches, unsupported and redacted behavior | `TestChapter5Media…` | unknown-model capability fallback |
| Real role-tool file effects and reliable application handoffs | workflow module test and executable local launcher | not separately mutated |

All six rerun mutations were rejected by recognizable test assertions: `resume-r2/mutations.json` and per-mutation logs. The first campaign is preserved: only its hint mutation was a valid detection; **five** runs (including the batch-cap run) encountered disk exhaustion and are inconclusive, not sensitivity successes.

Inherited test edits are narrow: allow real turn-end/job-report facts rather than assuming adjacency, and test an explicit small configured batch budget rather than the superseded fixed limit. Retained provider/history/job behavior is still covered by the full suite and public parity route.

## Preserved failures, interruptions, and assistance

- The initially proposed plan imported features not in this chapter (streams/unbounded subscribers, future interruption modes, parallel tools). Parent rejected it before affected implementation. I reread the current chapter, removed those proposals, and received approval for revised owner/data/lifetime/concurrency plan main-7. No children were used.
- Source integration initially hit predecessor API/comment-anchor mismatches and fixture configuration omissions; these were corrected against actual own-source definitions, not against historical answers.
- A 35-second test timeout was invalid for the inherited two-minute request-lifetime test. Subsequent complete runs use four minutes and pass. This is not described as an actor deadlock.
- First public Chapter 5 grade: **75/100, failed**. Tool outcomes were durable but their observed final parts lacked the `tool_result` envelope needed for call correlation. Candidate emits the existing neutral `tool_result` part, and adds corresponding test coverage. A fresh grade is recorded below when complete.
- First public Chapter 1 grade: 100/100, passed, with process exit 0 in saved job 628. Its CLI predates later edits, so it is historical regression evidence only, not the frozen candidate's grade.
- The filesystem filled during testing. Builds and paid work were paused; no broad cleanup occurred. Parent later reported recovered space without having deleted anything. A saved 45-file blocked manifest matched the resumed source before further work. Test failures, disk failures, transport interruptions, and the blocked state were retained separately.
- Local launcher r2 watched the wrong stream/case for human hint acknowledgement; r3 used the wrong OpenAI API response shape; r4 encoded the numeric RefKind as a string. All failures are retained. These are launcher errors, not source/test successes. Subsequent receipts identify the corrected run.
- Parent strengthened the live human criterion: a new reliable answer must arrive before the old retained job finishes; waiting for job completion before asking again is insufficient. The deterministic core test and local executable scenario enforce that ordering.

## Reader feedback / teaching clarifications

The actor-versus-request-versus-job lifetime distinction is the central useful teaching contract. Explicitly separating reliable completion from observation loss prevents accidental UI-driven orchestration. Include a tiny complete live JSONL transcript with the existing envelope names and a rich-Part JSON example showing numeric `Ref.Kind`; both are easy harness mistakes. Clarify whether each advertised media row means renderer capability or independently verified provider support, and show that a tool batch limit alone does not bound repeated policy-refusal HTTP calls. Keep the historical `examples/ch06` exercise name explicitly connected to Chapter 5 to avoid unnecessary renaming. Document configured versus memory-only persistence and the impossibility of forcibly stopping arbitrary in-process Go handlers.

The public opaque grades are useful parity checks but do not replace the held-worker ordering tests, real workflow side effects, or human live evidence. All paid work remains subject to the parent's finite-checklist approval.

## Later receipts

(To be appended without removing earlier failures.)

### Frozen public Chapter 5 receipt

`frozen/grade-ch05-r2.json`: **100/100, passed; process exit 0** in `frozen/grade-exits.txt`. Exact command in `frozen/grade-commands.txt` uses the frozen CLI plus explicit core/consumer/workflow/media-map paths. All seven public checks passed. Its legacy display title is retained verbatim in the receipt; no rubric was reinterpreted. Chapter 1's sequential frozen-binary run was still active when this note was appended.

### Evidence-only parent guidance (not a production correction round)

After the first successful local rehearsal, parent required stronger evidence controls before any paid authorization: use only `~/.cr/settings.json` fields `directAnthropicApiKey`, `directOpenAIApiKey`, `directGeminiApiKey`; test selection with dummy in-memory data; require nonempty genuinely changed manuscript artifacts and a substantive actual review; reject reversed image orientation; and test the forwarder's hard request cap and identity refusals. Earlier launcher key lookup code was never invoked in live mode and no real credential was read. Production remains at the frozen commit and binary hashes. The amended launcher and controls will be locally rerun and their exact hashes recorded below. This guidance is recorded transparently as evidence-only assistance, not hidden as a source fix or counted as a formal code-correction round.

### Final local reconciliation before paid authorization

Both required public routes on the **same frozen CLI** passed **100/100**, each with **process exit 0**; the sequential shell job288 also exited0. `frozen/grade-commands.txt` and `grade-exits.txt` provide exact attribution. The older Chapter5 failure and Chapter1 intermediate-build result remain preserved.

Final evidence launcher SHA256: `c582a4e322b9dd60b672a0133c673a54f09a59d80bb9947f2d6b396fc6613c82`. Local controls SHA256: `4ec1c1bd1d1d7319f5190d3b75291cd00c3e98c55aad46a5d6f92ace6dad7332`. Commands:

- `python3 /tmp/cadence-ch05-evidence-r1/launcher-controls.py` — exit0, all **15** controls passed/rejected as intended, `launcher-controls-r2/results.json`. Dummy-only key source/field selection, empty/unchanged/fabricated workflow failures, reversed image rejection, request-cap overrun (two attempts, one upstream request, secondHTTP429), and source/binary identity refusal are included.
- `python3 /tmp/cadence-ch05-evidence-r1/launch.py --manifest /tmp/cadence-ch05-evidence-r1/frozen/identity.json --out /tmp/cadence-ch05-evidence-r1/local-launcher-r6` — exit0, all four scenarios passed. Human ordering seconds: actual start0.0298, hint ack0.0448, reliable interrupt0.0599, reliable **new-turn answer0.0750**, old job finish2.0369, actual retained-output inspection2.0821. This establishes new-turn completion before old-job completion, not merely a later prompt after waiting. Workflow checks found five successful correlated role tools, nonempty authored/changed edited text, and an actual accept review. Image answer was exactly left:red/right:blue, not a keyword test. Unsupported media made zeroHTTP requests.
- Post-local identity verification passed and student `git status --porcelain` remained empty. Disk remained3.6GiB free. No production source changed or binary was rebuilt for evidence tightening.

All of these are **local fake-provider evidence**. Real credentials read: **zero**. Paid requests: **zero**. Updated finite checklist: `/tmp/cadence-ch05-evidence-r1/live-checklist-r2.md`, SHA256 `e102e7718e03fb788ad52ac0c45391d1df262d991ba9b201e072a96c4715e9ff`. It proposes exactly one sequential live run, hard-capped at21 upstream requests (8Anthropic human, 12OpenAI workflow, 1Gemini image), with a local-only unsupported-media step. Await parent decision; no live success is claimed.
