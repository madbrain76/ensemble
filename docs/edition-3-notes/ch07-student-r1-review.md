# Chapter 7 student — formal correction Round 1

## Disposition

Both assigned public-behavior corrections are implemented in a new clean student candidate. The **frozen Chapter 7 adapter passed 100/100, `passed:true`, exit 0**, and the **same exact CLI passed Chapter 1 100/100**, including both genuine zero-point delayed-header/body phases. The complete sequential module race/vet/build matrix, focused core/real-WebSocket/page regressions, and compile-valid mutation controls meet their expected outcomes.

This is **formal correction ROUND1**, not acceptance, tagging, or a replacement for independent whole-code review. Root alone accepts/tags. The unchanged-r4 reconciliation and every earlier failed grade remain historical and unmodified. The broad grader had already passed r4 after its accepted fixture repair; that score did not detect these two review findings. The new behavioral failures and controls, not score alone, are the correction proof.

## Authority and boundary

The complete `/tmp/cadence-ch07-correction-r1-task.txt` was read and hash-verified (`becb805c57e97abc2606be09f5ff1c2cf6c892be63a1fe8e0d9089f42120a88c`, 6907 bytes). Mandatory public prerequisites and current Chapter 7 were read in full in bounded ranges, with the relevant earlier explanation/refusal, identity and dispatch/wakeup passages recovered. `before.json` records public input identities and the clean required baseline.

Only the two assigned behaviors were changed. No future-chapter structures, private reviewer/author/parent material, hidden cases, supplied patches, or other workers' artifacts were consulted. No grader, manuscript, rubric, guidance, reference, acceptance tag or main-repository source was edited or staged. No children were spawned. No new browser, native/headed/audio/microphone, paid-provider or credential-reading campaign was performed. Pagehide/focus/voice exploration was explicitly not promoted into a correction.

A temporary evidence-path scope error is disclosed below rather than claiming perfect path compliance.

## Frozen identity and artifacts

- Repository: `/Users/bill/projects/ensemble-edition-3-student`
- Commit: `463996b765631d31923a5726036a34942a67cefd`
- Tree: `5e95dcdd716f7c25c96cc3ebcde13d000bb38670`
- Direct parent: `48e1d12d9c396b5988c0ce6eccdef67e581872ff` (the required r4)
- Worktree: **clean**, with all **120 tracked source-file hashes** verified after execution.
- Immutable `edition-3-ch06-r1`: tag object `84d0bce5f2525b406affa64fb1a8eb70711b1e06`, peeled commit `f595c787675e2a6343053ec42ea6c7e1e02a930e`; unchanged.
- All frozen binaries and the archive are under `.ch07-local/correction-r1/frozen/`, excluded from Git.

| Artifact | SHA256 |
|---|---|
| `consumer` | `5088a4a63efb7df413c6acdcf5d06c40ff08fda7dbe979d195a82ac223733d06` |
| `ensemble-cli` | `000a2997e5b8396feb2e0892f3b8160f48cce8adc34b44eeb940e63dd3873f11` |
| `ensemble-gui` | `c41d8115a6b22b8fb294ab11f4303926e5d485cac7b40025c39fa0e42c0eccd4` |
| `source.tar.gz` | `f2ba25cedc080abd295f35af9aeae3128dab0793d9fb1c6c71abef5c2750593a` |
| `streaming` | `bdd55dfcd26eb18b776f9a9b80a133a8a6604b8d0289498715ffd857d8f6c1ea` |
| `workflow` | `5fe5bcea5ae0a7cdff0bcbb0f3153aa75ed773477b9b76abbab1a015302175ee` |

`ch07-student-r1-source-manifest.json` pins every source file, archive, artifact, exact declarations, command receipts, grader and asset identity. The build-info logs record the compiled CLI/GUI module and VCS metadata. Source did not change after freezing; checks preceded/followed each sequential command. Both graders used the CLI hash shown above, verified again after both completed.

The public adapter builds optional module executables internally from the frozen module paths. Its temporary internal GUI binary hash is not available in the report; the standalone GUI hash is **not substituted** for that temporary binary. The source is pinned for both, and the standalone artifact is separately pinned.

## R1 — truthful execution projection

**Original failure:** Chapter 3 intentionally retains a `tool_called`/`tool_returned` pair for each explanation-refused proposal. The r4 GUI treated every `tool_called` as execution and therefore emitted/summarized refused proposals as `tool_dispatched`, both live and on replay. Removing the required called/returned pairs or filtering selected names/error strings would be incorrect.

**Narrow correction:**

1. `internal/common/common.go` adds optional data-only `ToolEvent.Dispatched`. It records admission, not successful completion.
2. `internal/agent/dispatch.go` stamps that fact only in the actual admitted `prepareCall` path, not the explanation-refusal branch. Actor decisions, job ownership, refusal feedback, identities and one-shot limits are unchanged.
3. `gui/hub.go` requires positive admission evidence: the durable bit, an actual retained-job record, or the real admission observation already seen for that identity. A proposal lacking that evidence renders as a final call artifact, **not** an execution lifecycle envelope. Correlated results/errors remain readable. Existing exact-once dispatch bookkeeping remains in place.

The real Agent fixture first returns an unexplained three-call proposal (file, process, control), then an explained admitted retry with reused provider call IDs. Before retry, actual file/process effects are absent, the provider receives all three refusal results, live/reconnected messages claim zero starts, and identity pairing remains intact. The retry performs the actual file/process work and admits the control; there are three core admissions, two real retained-job identities, three live dispatches and three replayed dispatches. Distinct durable display IDs separate the reused provider IDs.

The regression traverses local fake HTTP, a real public Agent/job owner, actual hub subscriptions, **real WebSocket frames**, serialized history load into a fresh Agent, and the controlled production page. No private actor state or injected GUI dispatch creates this proof. Complete raw wire/history captures stay local; `focused-safe-summary.json` and the pinned test/assertion logs are portable.

**Legacy qualification:** prior job-backed admission records retain execution proof. A legacy control with neither a job, durable admission bit nor current real observation cannot be truthfully classified as executed: its call and result remain artifacts, but replay does not invent a start. Newly admitted no-job controls have the durable bit and replay correctly. No existing history is rewritten or deleted.

## R2 — complete on-demand tool speech

**Original failure:** `tool_finished` appended visible output but left the artifact's readable text as its call summary. The actual speaker button therefore queued a caption while omitting the visible result.

`gui/web/gui/tool-readable.js` projects the readable call context/arguments and accumulated result or error as plain text. `gui/web/gui/app.js` keeps tool-call and result state together and refreshes the artifact's speaker text on both result arrival and later call-final updates. The short automatic dispatch summary remains separate and concise; receiving a tool result does **not** automatically narrate it. No native speech API, speech queue, parser, server result clipping, tool executor or renderer subsystem was redesigned.

The controlled test evaluates the **actual page module and its actual imported SpeechQueue**, delivers real-handler messages, and clicks the actual artifact speaker closure. It checks ordinary `read_file` text including literal markup, specialized `run_command` JSON with stdout/stderr and an error, full result inclusion, metadata exclusion, a late call-final update, replay/reconnect silence and on-demand recovery, empty/opaque silence, cancellation, stale callbacks, current error release, and pause/unpause intents. It also feeds actual core/hub/WebSocket refusal receipts through the page and verifies no automatic execution summaries while errors remain readable on demand.

The DOM, WebSocket client and synthesis boundary in that test are explicitly controlled test objects. These are **not native or browser/listening observations**. The test exercises production page/speaker logic rather than a pure-helper assertion or a receipt-only claim.

## Failure sensitivity, not merely green tests

- Both final regressions run unchanged against an isolated copy of **original r4 production**. They compile/evaluate and fail their intended assertions: false dispatches and omitted spoken results. `original-r4-r1-red.log` and `original-r4-r2-red.log` retain these failures. The R2 script stops at its first failed ordinary-result assertion; it is not claimed to have completed every later check on original code.
- A compile-valid R1 mutation claims every proposal was dispatched. It fails the refusal/admission assertions, including the wire path, while actual refused process effects remain absent.
- A compile-valid R2 mutation discards the stored result while leaving visible output rendering intact. The actual speaker assertion fails because the result is absent.
- Each mutation lives only in an excluded isolated copy, records before/broken/restored hashes, and passes after exact restoration. The working/frozen source is never mutated. These controls were repeated sequentially after freeze and are in `frozen-control-results.json` and their command/output receipts.

## Frozen sequential validation

`run-frozen-sequential.py`, per-command receipts and `sequential-completed.json` retain the exact cwd, closed environment, command, start/end, exit and output hashes. The completed sequence contains 26 commands: 24 expected zero exits and two **intentional semantic-failure** controls, followed by their passing restorations. There were no concurrent adapter runs or source edits during these phases.

| Phase | Outcome |
|---|---|
| Root `go test -race -count=1 ./...`, vet, CLI build | pass |
| GUI race tests, vet, GUI build | pass |
| Inherited consumer, workflow and streaming modules: race tests, vet, builds | pass, each separately |
| Focused real Agent/HTTP/WebSocket regression | pass |
| Controlled production-page/speaker regression | pass |
| Existing phrase/speech queue checks | pass |
| R1/R2 compile-valid mutations and exact restorations | intended failures; restored passes |
| Current repaired public Chapter 7 adapter | 100/100, all five checks pass, exit 0 |
| Chapter 1 on the same frozen CLI | 100/100, all nine checks pass, exit 0 |

The Chapter 1 delayed-header and delayed-body checks each completed their genuine **1m2s** phase, both `passed:true` despite zero points. No subprocess timeout, shortened fixture or forced exit was applied. The finite new focused unit-test timeout is not a deadline on those lifetime phases. Longer arbitrary request lifetimes still depend on source review; a 62-second observation is not an unbounded-runtime proof.

Exactly one final Chapter 7 and one final Chapter 1 adapter invocation ran in this correction. The same root-approved repaired grader executable is SHA256 `c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e`. The inherited consumer/workflow/streaming/GUI declarations and existing media map were reused unchanged; no new grader mode/model declaration or case assertion was invented.

Chapter 7 command (cwd `/Users/bill/projects/ensemble`):

```sh
/tmp/cadence-ch07-grader-pause-fixed -ch 7 -json -core /Users/bill/projects/ensemble-edition-3-student -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer -consumer-case /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json -workflow /Users/bill/projects/ensemble-edition-3-student/examples/ch06 -stream-exercise /Users/bill/projects/ensemble-edition-3-student/examples/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble-edition-3-student/gui -gui-model claude-sonnet-4-6 -media-models /Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-evidence/declared-media-models.json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli
```

Chapter 1 command, same cwd:

```sh
/tmp/cadence-ch07-grader-pause-fixed -ch 1 -json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli
```

The raw machine reports remain verbatim, including the Chapter 7 adapter's inherited Edition 2 title. The current repaired adapter's five checks cover observation types, replay, GUI log, pause-held dispatch and corrected Chapter 6 parity. A passing score is not whole-code acceptance.

## Native boundary and exact carryover

The correction brief states that prior completed native/heard-by-Bill evidence belongs to r4. It is **not relabeled** as a run of this candidate. `native-evidence-applicability.json` compares exact bytes against r4:

- `speech.js`, `queue.js`, `artifact-scroll.js`, `format.js`, `index.html` and `style.css` are byte-identical.
- The local-pause/release functions and the draft/control/pagehide/visibility/voice tail are unchanged.
- `app.js` changed in the tool-artifact projection, and `tool-readable.js` is new. Their relevant new result selection and queue interactions have controlled evidence, not new native speech-service, permission, listening, large-payload performance or whole-bundle browser evidence.

The earlier native mechanisms retain only the narrow applicability justified by those unchanged bytes. No new browser/native/audio/mic run, focus/activation manipulation, paid call, or native launcher/plan change occurred. Paid/live services were not run because they were not authorized.

## Retained failures, scope deviation and uncertainty

The development record is intentionally not scrubbed. `development-notes.txt` and the original logs retain: patch-anchor refusals/partial patch, an out-of-scope helper compile error and dependent missing-wire failure, a no-tests-selected command, the real-wire test reader's 127-byte sentinel bug, a failed metadata grep, and overbroad/misnamed preservation-helper checks. Setup, compile, missing-file, metadata and no-test results receive **no sensitivity or coverage credit**. Later successful logs use new paths rather than overwriting failed logs.

I initially created new evidence under the **wrong** path:

`/Users/bill/projects/ensemble/docs/edition-3-notes/ch07-correction-r1-evidence`

The authorized path is:

`/Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-r1-evidence`

Only files newly created by this correction turn were moved; no existing package was overwritten. Full raw wire/history captures were moved to excluded `.ch07-local/correction-r1/raw`. `scope-relocation-receipt.json` retains every exact old/new path and contemporaneous before/after byte count and SHA256. Historical command argv/cwd/path records remain unchanged and use that relocation map. The wrong directory is gone. `development-notes.txt` was subsequently appended to; the final check proves its original moved prefix and records its new hash. This is an acknowledged output-path scope deviation, not claimed perfect compliance.

The transient relocation tool command remains in the session. A separately saved byte-for-byte copy of its launcher script was not retained; the durable exact transfer/hashing receipt and this limitation are explicit. Parent was informed and directed continuation within the same bounded scope.

Final preservation verifies **252 original portable files, 22 original local artifacts, and 17 reconciliation files**, plus both original manifest hashes, unchanged. Old failed grades and the separate repaired-r4 pass remain intact. Private/raw captures, binaries and archives remain local-only; portable evidence is safe UTF-8.

Assistance was the public task's two findings, its prior accepted fixture-repair authority, and the parent's scope/progress clarification. No private review body, grader source/case body, supplied patch or acceptance decision was used as student implementation evidence.

## Final files and authority

- `ch07-student-r1-review.md` — this report.
- `ch07-student-r1-evidence/` — safe retained logs, command receipts, counts, controls, applicability and preservation.
- `ch07-student-r1-source-manifest.json` — complete source, artifact and invocation pins.
- `ch07-student-r1-evidence-manifest.json` — portable-file hashes and local-only pointers.

All launched jobs are finished; no children or execution jobs remain. The student worktree is clean at the identity above. No acceptance tag, push, publication, main commit/stage, live-GUI restart/disruption or installation was performed. Root owns the independent correction review and acceptance/tag decision.
