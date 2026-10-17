# Edition 3 Chapter 5 — student candidate report

> **Final status:** clean finite candidate `00489b05b65cfe5209f26cb655667df44b3847e2`, tree `095e8f8a2559275d94415bcb18bef5b5540dc2b9`; both frozen public grades100/100 (exit0), full source gates passed, three-provider live feature evidence complete with12total requests. The original human launcher failure remains failed and is independently reconciled; the mocked route control is uncredited. Parent review/acceptance/tag is pending. Earlier “pending/no paid” entries below are chronological records, not the final status.

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

After the first successful local rehearsal, parent required stronger evidence controls before any paid authorization: use only the approved `~/.cr/settings.json` fields and test selection with dummy in-memory data (I initially assumed incorrect field spellings; correction recorded below); require nonempty genuinely changed manuscript artifacts and a substantive actual review; reject reversed image orientation; and test the forwarder's hard request cap and identity refusals. Earlier launcher key lookup code was never invoked in live mode and no real credential was read. Production remains at the frozen commit and binary hashes. The amended launcher and controls will be locally rerun and their exact hashes recorded below. This guidance is recorded transparently as evidence-only assistance, not hidden as a source fix or counted as a formal code-correction round.

### Final local reconciliation before paid authorization

Both required public routes on the **same frozen CLI** passed **100/100**, each with **process exit 0**; the sequential shell job288 also exited0. `frozen/grade-commands.txt` and `grade-exits.txt` provide exact attribution. The older Chapter5 failure and Chapter1 intermediate-build result remain preserved.

Final evidence launcher SHA256: `c582a4e322b9dd60b672a0133c673a54f09a59d80bb9947f2d6b396fc6613c82`. Local controls SHA256: `4ec1c1bd1d1d7319f5190d3b75291cd00c3e98c55aad46a5d6f92ace6dad7332`. Commands:

- `python3 /tmp/cadence-ch05-evidence-r1/launcher-controls.py` — exit0, all **15** controls passed/rejected as intended, `launcher-controls-r2/results.json`. Dummy-only key source/field selection, empty/unchanged/fabricated workflow failures, reversed image rejection, request-cap overrun (two attempts, one upstream request, secondHTTP429), and source/binary identity refusal are included.
- `python3 /tmp/cadence-ch05-evidence-r1/launch.py --manifest /tmp/cadence-ch05-evidence-r1/frozen/identity.json --out /tmp/cadence-ch05-evidence-r1/local-launcher-r6` — exit0, all four scenarios passed. Human ordering seconds: actual start0.0298, hint ack0.0448, reliable interrupt0.0599, reliable **new-turn answer0.0750**, old job finish2.0369, actual retained-output inspection2.0821. This establishes new-turn completion before old-job completion, not merely a later prompt after waiting. Workflow checks found five successful correlated role tools, nonempty authored/changed edited text, and an actual accept review. Image answer was exactly left:red/right:blue, not a keyword test. Unsupported media made zeroHTTP requests.
- Post-local identity verification passed and student `git status --porcelain` remained empty. Disk remained3.6GiB free. No production source changed or binary was rebuilt for evidence tightening.

All of these are **local fake-provider evidence**. Real credentials read: **zero**. Paid requests: **zero**. Updated finite checklist: `/tmp/cadence-ch05-evidence-r1/live-checklist-r2.md`, SHA256 `e102e7718e03fb788ad52ac0c45391d1df262d991ba9b201e072a96c4715e9ff`. It proposes exactly one sequential live run, hard-capped at21 upstream requests (8Anthropic human, 12OpenAI workflow, 1Gemini image), with a local-only unsupported-media step. Await parent decision; no live success is claimed.

### Credential-field setup correction caught before paid work

Parent explicitly withheld authorization after catching my incorrect field-name assumptions. Those spellings were my mistake, **not the parent-supplied names**. The exact approved fields are `directClaudeAPIKey`, `directOpenAIAPIKey`, and `directGeminiAPIKey` in `~/.cr/settings.json`. No real credential had been read and no paid call had occurred. Previous scripts/report are preserved under `/tmp/cadence-ch05-evidence-r1/key-selector-r3/`.

Only the selector mapping and matching dummy fixture spellings changed. `key-selector-r3/launch.py.diff` proves one changed mapping line and byte-for-byte equality of all scenario, semantic-check, identity, and cap code. The corrected selector passed three positive dummy vendor tests and twelve negative tests (missing, old spelling, empty, non-string), with the loader's `read_text` replaced in memory and its exact approved path asserted. Command `python3 /tmp/cadence-ch05-evidence-r1/key-selector-controls-r3.py` exited0. No production rebuild or full grade/scenario rerun was performed: unchanged controls/rehearsal receipts carry forward honestly for unchanged code only. Corrected selector evidence is new, not retroactively attributed to the earlier run.

Current launcher SHA256 `4516948d1dc809555bb0cad9c0c055c450f85c9cdd793c574ec4e36d617a34c2`. Corrected finite checklist is `live-checklist-r3.md`; exact checklist/control hashes and scope proof are in `key-selector-r3/result.json`. Frozen production identity reverified. **Paid calls still not authorized; zero real settings reads, zero paid requests.** Awaiting the parent's explicit bounded decision.

### Authorized live attempt 1 — preserved launcher failure, no retry

Parent main-5 approved exactly one run of checklist r3 / launcher `4516948…`, with hard caps8+12+1 and stop on first failure. Executed the exact command into `/tmp/cadence-ch05-evidence-r1/live-r1`; shell job338 and launcher exited2. The run stopped before workflow, image, or the final local unsupported step. **Actual forwarded requests: four Anthropic, all HTTP200; zero OpenAI and zero Gemini.** Human CLI process itself exited0 after28.207s. Request counters and all safe partial receipts are retained.

Precise diagnosis: the final human guard assertion searched all assistant stdout for the string `KEY_ENV_LEAK`. The assistant's actual final prose said **“No `KEY_ENV_LEAK` was printed, confirming `LLM_API_KEY` was not set in the environment.”** That benign sentence triggered the broad substring check. The real `run_command` argument contains the expected environment guard. Its retained `cr/io/1` contains only `CH05_JOB_STARTED`, `CH05_JOB_FINISHED`, and `exit_code: 0`, with **no guard sentinel**. The actual side-effect file contains FINISHED. Thus this is an evidence-check false positive, not evidence of a credential leak or a production failure. No credential values were printed or included in evidence.

The preserved actual facts independently show the requested human behavior: hint was visibly acknowledged while the job was running; current turn ended interrupted at `2026-10-17T05:51:09.380481Z`; the new user turn ended at `05:51:10.421531Z`; the old retained job's done/exit0 report arrived at `05:51:29.380193Z`; actual `wait_for_job` output was inspected and the final turn ended at `05:51:34.122969Z`. The new answer preceded actual old-job completion by about19seconds. Exact inputs/times, safe stdout/stderr/events, exit receipt and actual job files survive under `live-r1/human`. No synthetic completion is used as evidence of the real job outcome.

The original `live-r1/summary.json` remains **passed:false** and is not rewritten as success. Its expected `human/result.json` was not written because the final assertion failed. Remaining required live coverage is the OpenAI three-role workflow and the Gemini image path. No retry, further provider call, model substitution, source edit, or rebuild has occurred; parent is being asked for a bounded next decision.

### Guard reconciliation and remaining-only local route

Parent main-6 accepted the false-positive diagnosis, not the original launcher PASS, and approved only an evidence correction plus local validation. Actual human facts/hashes/times are saved in `live-r1-human-observed-facts.json`: new turn precedes actual old-job terminal report by18.958662seconds; actual guarded process exit0 and real inspection retained. No human inference was repeated.

Current launcher `68066167187e3618beba4d9c8a49b59ee5f2ba105286dc3ee23235960c4cbfad` checks the sentinel in actual `cr/io/1`, and adds an explicit `--scenario remaining` path that cannot select human. Prior launcher/diff remain under `guard-route-r4`. Credential selection, caps, models and remaining scenario bodies are byte-for-byte unchanged. One local parser-insertion error was caught before execution and preserved (`first-control.stderr`); it was not a provider/source failure.

`python3 /tmp/cadence-ch05-evidence-r1/guard-route-controls-r4.py` exited0 with five controls: benign assistant prose accepted, real malicious process sentinel rejected, saved human artifacts/hashes reconciled without inference, no-human dispatch enforced, hard cap verified (second local attempt429; one upstream). `python3 /tmp/cadence-ch05-evidence-r1/launch.py --manifest /tmp/cadence-ch05-evidence-r1/frozen/identity.json --scenario remaining --out /tmp/cadence-ch05-evidence-r1/local-remaining-r4` exited0 and passed only workflow/image/unsupported fake scenarios, recording `human_run:false`. These are fresh bounded evidence controls, not a rerun of production grades or the entire earlier campaign.

Remaining-only proposed checklist: `live-checklist-remaining-r4.md`. It requires new explicit approval for one workflow(max12)+image(max1)+local unsupported command; no additional human call. Existing actual paid count remains4Anthropic/all200,0OpenAI,0Gemini. Production source/archive/binaries remain frozen.

### Remaining-only authorization, recovery, and evidence-credit correction

Parent main-7 approved exactly one remaining-only command with launcher `68066167187e3618beba4d9c8a49b59ee5f2ba105286dc3ee23235960c4cbfad`, checklist `db2f64394c847e6b558fa039690037ef5e613a2f3d48cedaef335a9f10f195d6`, and unchanged frozen identity. Before execution, a coding-agent inference/transport request failed (“Responses inference did not complete”); this was not a student source/test or provider-scenario failure. On resume, I verified that `live-remaining-r1`, its stdout/stderr/exit/authorization receipts, and any matching launch process were absent; hashes matched. Only then was the exact authorized command started once (owned job401). No human or previously started scenario was blindly repeated.

**Correction to credit, not to the raw control receipt:** `guard-route-controls-r4.py` used `unittest.mock` to replace scenario functions in its no-human dispatch check. Bill's no-mocks rule was not waived. That mocked route control is **uncredited toward acceptance**, retained and disclosed as an attempt. The actual completed `local-remaining-r4` rehearsal using the frozen executables—six fake workflow HTTP requests, one fake image request, unsupported zero HTTP—plus inspected dispatch code is the route evidence. The direct real-output guard controls, saved-artifact/hash reconciliation, and real local HTTP forwarding-cap control retain their separate value. The earlier raw receipt lists five controls, but it must not be represented as five credited acceptance controls. Credential-selector dummy data checks likewise are not a substitute for actual live credential loading; their mocked file-loader check is limited to setup evidence, not production acceptance.

No production source, archive, binary, model, or approved launcher hash changed during recovery. The remaining authorized ceiling is12OpenAI+1Gemini, local unsupported0, stop first failure, no retry/substitution. Existing human count remains4Anthropic HTTP200; original failed summary remains unchanged.

## Final candidate outcome

**Finite candidate ready for parent review, not accepted/tagged.** Production remains clean at `00489b05b65cfe5209f26cb655667df44b3847e2`, with the previously recorded source archive and both frozen race binaries unchanged. Post-live identity verification passed and `git status --porcelain` was empty. No source change, rebuild, model substitution, speculative retry, or repeated human inference occurred during live evidence collection.

The remaining-only authorized command (job401) exited0. `live-remaining-r1/summary.json` records `passed:true`, `live:true`, `scenario:remaining`, `human_run:false`.

| Feature | Actual evidence | Result / limits |
|---|---|---|
| Human actor responsiveness, hint, interruption, new turn while old work continues | Anthropic `claude-sonnet-4-6`; original human CLI exit0; four requests, allHTTP200; actual job output/side effect and saved event times | New reliable turn finished18.958662s before old actual job report; real retained output inspected. Original overall launcher exit2/failed summary remains preserved due diagnosed guard-prose false positive; independent reconciliation is credited, not an invented original PASS. |
| Three-role application workflow | OpenAI `gpt-4.1-2025-04-14`; process exit0 in15.187s; seven requests, allHTTP200 | Actual five correlated role tools succeeded; author78words, editor89words and genuinely changed text; actual nonempty accept review; final manuscript571characters and review313characters, with contents matched to real tool calls/results. |
| Positive local image | Gemini `gemini-3.8-flash`; process exit0 in9.221s; one request, HTTP200 | Exact structured answer left:red/right:blue from supplied PNG; stable reference and no image body in persisted history. Fixture SHA256 `057d5c091556b8846631c6c1e00c276ee1033a6d5b2c93cc078a38429e62605e`. |
| Unsupported media | Local fake only, same frozen CLI with Anthropic selected row and video | Expected executable exit1, explicit refusal, **zero HTTP requests**. This is boundary evidence, not a paid support claim. |
| Retained provider/CLI parity and Chapter5 public acceptance surface | Same frozen CLI, public opaque routes | Chapter5 and Chapter1 each100/100 passed; both process exits0. |
| Concurrency regression sensitivity | Full root race suite, standalone module gates, six valid mutations | All source gates passed;85.5% root statement coverage. Five earlier disk-blocked mutations remain inconclusive. Mocked route control is uncredited; actual local executable rehearsal is credited. |

**Total actual upstream requests:12 =4Anthropic+7OpenAI+1Gemini, allHTTP200.** No extra request was made after either authorized invocation. Hard budgets were respected. Actual billing USD is not inferred from request count. Raw safe stdout/stderr, exact stdin/times, filtered events, process exits, actual role/job artifacts, counts, and hashes are preserved. No real credential values, raw provider response bodies, thinking content, or signed opaque payloads are included in the portable evidence.

Remaining limits are explicit: whole-part deltas, not token streaming; observers must not block and bounded queues may lose observations; ordinary nonreturning Go handlers can delay shutdown; library persistence is configured/optional while the executables configure it; only selected explicit model capability rows are supported; audio/video/PDF positive evidence is local renderer/HTTP testing, not live provider evidence; URI fetching is unsupported; GUI remains a build-tested stub. Parent owns code review, acceptance and tagging. No acceptance claim is made here.

Portable evidence/index paths and hashes are recorded in the companion `ch05-student-evidence-manifest.json`. Earlier failed/interrupted/inconclusive receipts and evidence-only corrections remain present; later passing checks do not overwrite them.

### Portable handoff locations

- Main report: `/Users/bill/projects/ensemble/docs/edition-3-notes/ch05-student-review.md`.
- Main sanitized text evidence/scripts/receipts: `/Users/bill/projects/ensemble/docs/edition-3-notes/ch05-student-evidence/`.
- Main exact evidence index: `/Users/bill/projects/ensemble/docs/edition-3-notes/ch05-student-evidence-manifest.json`. Run the included `verify-portable.py` to verify every copied text hash and this report; it performs no provider call and reads no credentials.
- Existing source manifest: `/Users/bill/projects/ensemble/docs/edition-3-notes/ch05-student-manifest.json`; existing ordinary source snapshot `.agents/chapters/code/ch05-student` is left untouched pending parent scope review.
- Authoritative source archive and race binaries remain outside Git under `/tmp/cadence-ch05-evidence-r1/frozen/`, exactly identified by `final-candidate-state.json` and frozen identity SHA256 `950d5c3fb732ad49ff132bbb240e3e0405628c623313cc28da23baa458c6ac19`. Archive SHA256 `336b26e12c16c268a1634dd690d4fc9be075aeed4b1f323ff36c42a1436f70df`; CLI `258233df606e71d50aa0103c5865fcd2e383efc471c3d3f8e71fe59f07609863`; workflow `c5c2fee6aa8c463853c1c6210f7b2a699be8bf6806e550ba0b28d0aa44507114`.
- Binaries, archives, PNGs and interpreter caches are not copied into the main text evidence. Their local paths/hashes and exclusion reasons are indexed. The exact PNG generator is retained in the tested launcher; fixture hash is in both actual image and final receipts. No optional archive under `.agents/chapters/code/ch05-student-evidence.tar.gz` was created.

The portable receipts retain original absolute invocation paths for historical accuracy, while the evidence index uses relative paths so copied text can be verified after relocation. Do not execute historical paid-launch commands from the bundle: each was authorized for one invocation only, now consumed. Any new provider call requires new explicit approval.
