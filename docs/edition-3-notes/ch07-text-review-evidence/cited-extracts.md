# Exact source extracts

Safe report/receipt excerpts only, not new executions. Sources are hashed in read-ledger.json.

## book/edition-3/chapter-07.md

```text
164: ```
165: local_pause_cause = tts_queued_or_speaking OR user_typing
166: ```
167: 
168: The Agent checks its gate before starting each tool; the server calls the
169: public Agent's `Pause` or `Unpause` on receiving those commands. Pause is
170: Agent-wide. An explicit unpause in any tab unpauses it for everyone; unchanged
171: local typing or speech in another tab cannot veto that decision. Rules:
172: 
173: 1. TTS queues an utterance: send `pause`. Set the flag when the
174:    utterance is queued, not when audio begins.
175: 2. TTS queue empties: send `unpause`, unless the user is typing.
```

```text
206: Build an optional public-library GUI module: one Agent, one HTTP server on
207: `--port` (default 8088), serving its static files through `--gui-dir`. The
208: accepted example lives in `gui/`, with `cmd/` and `web/gui/` inside it; historical
209: chapter IDs in evidence remain `ch08`. The binary accepts prompts
210: over WebSocket, streams observations to all connected clients,
211: supports reconnection, gates tool calls on pause, and logs to gui.log.
212: 
213: Configuration uses public flags, including `--port`, `--gui-dir`, `--vendor`,
214: `--model`, `--base-url`, `--workspace`, `--event-log` and `--no-stream`.
215: Explicit flags override the inherited provider environment defaults. Keep
216: credentials in the intended process environment rather than command arguments.
217: The original flags-only sentence conflicted with its own environment-configured
218: grader; it did not justify breaking the earlier key setup.
219: 
220: For Edition 3, build the headless CLI from its command directory first. For a
221: nested `core/cmd/ensemble` layout, run:
222: 
223: ```sh
224: (cd /path/to/core/cmd/ensemble && go build -o /tmp/ch7-cli .)
225: ```
226: 
227: Then, from this book repository's root, pass that executable separately from
228: the source modules. Replace the paths and model names with your declarations:
229: 
230: ```sh
231: go run ./cmd/grade-edition3 -ch 7 \
232:   -core /path/to/core \
233:   -consumer /path/to/consumer \
234:   -consumer-case /path/to/consumer-case.json \
235:   -workflow /path/to/workflow \
236:   -media-models /path/to/media-models.json \
237:   -stream-exercise /path/to/stream-exercise \
238:   -stream-model DECLARED_STREAM_MODEL \
239:   -gui /path/to/gui \
240:   -gui-model DECLARED_GUI_MODEL \
241:   /tmp/ch7-cli
242: ```
243: 
244: The positional target is the executable, not the core directory or the GUI
245: server. The shallow directory build does not find a nested CLI command. Keep
246: the declared media map and all inherited module arguments; omitting them does
247: not isolate the GUI checks. An early student run also supplied a blanket
248: `LLM_BASE_URL` override, sending fixture requests to the wrong endpoint. Leave
249: that override out of the grading environment so the adapter can route them.
250: 
251: The grader builds the GUI module separately from the explicit `-gui DIR`.
252: It builds `./cmd` when present, otherwise the module root, and launches `--port`
253: plus `--gui-dir` pointing to that module's `web/gui`.
254: HTTP `/` must become available and WebSocket `/ws` accepts the messages above.
255: Its faithful provider peer supplies the inherited `LLM_VENDOR`, `LLM_MODEL`,
256: `LLM_API_KEY` and `LLM_BASE_URL` environment. `-gui-model MODEL` selects the
257: declared Anthropic fixture row; the historical default `claude-sonnet-5` is a
258: local test-table entry, not a model-discovery claim. Earlier consumer, workflow,
259: media and streaming flags still apply to core parity.
260: 
261: The optional module must resolve its public core dependency to the submitted
262: core. This boundary check accompanies `websocket-streams`; source review still
263: establishes actual ownership and shared-Agent use. Active-partial replay,
264: handoff ordering, the fixed window, paused interrupt and browser behavior also
```

```text
326: `ToolDispatched` fires before execution: call id, tool name, raw
327: input. `ToolFinished` carries the result of waiting: call id, result, error
328: flag. Chapter 4's distinction still holds. A timeout can return a handle while
329: the actual job continues, and its later output must remain visible even if
330: the human has interrupted the turn. `TurnEnded` announces the end of the turn
331: to displays; it does not replace the reliable completion used by `Ask`.
332: 
333: A requested call is not yet an execution. The Edition 3 student GUI announced
334: starts for a batch the engine had refused, both live and on replay.
335: `ToolDispatched` belongs only to calls admitted to run. Keep the requested calls
336: and refusal results visible, but do not infer an execution from their presence
337: in history. A result can explain why a tool never ran.
```

```text
373: 
374: ## 7.4 The pause gate
375: 
376: Pause is a live Agent control, not another conversation message waiting for
377: reduction. The Agent owns its gate. Common declares the public interface;
378: the behavior belongs with Agent, rather than in a shared utility package.
379: The GUI calls that Agent interface. It does not receive an Engine sibling,
380: become the Agent's parent, or create a second gate.
```

```text
568: The correction projects readable content explicitly. Empty text stays empty.
569: Readable Anthropic thinking and Gemini thought text can become prose; signature
570: and redacted bodies remain opaque replay material. Preserve those bytes in
571: history without pretending they are something to read aloud.
572: 
573: The focused browser test reused the retained Gemini part, plus readable and
574: opaque controls from an external peer. It checked the card, automatic speech
575: input and speaker-button input separately. Restoring the fallback failed all
576: three checks for each empty or opaque case, while readable controls still
577: produced native start/end events. No new model generation was needed to diagnose
578: what was already visible in the saved page.
```

```text
580: ---
581: 
582: ## 7.10 The voice
583: 
584: TTS auto-speaks the stream. Thinking and chat in full. Tool
585: dispatches abbreviated: "read file: main.go, line 42." Tool results
586: silent. A speaker icon on every Artifact card speaks the accessible
587: text on click.
588: 
589: The Edition 3 student page could display a long tool result yet speak only its
590: short generated caption. On-demand speech must use the full readable result or
591: error available to the page, not the abbreviated label used for automatic
592: speech. Opaque media still needs only its stable spoken description.
```

## docs/edition-3-student.md

```text
107: For Chapter 7, keep the accepted core CLI and inherited modules separate from the new GUI module. Build the actual CLI target, then run from the book repository root:
108: 
109: ```sh
110: go run ./cmd/grade-edition3 -ch 7 -json \
111:   -core /absolute/path/to/core-module \
112:   -consumer /absolute/path/to/consumer-module \
113:   -workflow /absolute/path/to/workflow-module \
114:   -stream-exercise /absolute/path/to/stream-exercise-module \
115:   -stream-model YOUR_DECLARED_STREAMING_MODEL \
116:   -gui /absolute/path/to/gui-module \
117:   -gui-model YOUR_DECLARED_GUI_MODEL \
118:   /absolute/path/to/student-cli
119: go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/SAME/student-cli
120: ```
121: 
122: All five module roots and both model identifiers are explicit; flags precede the executable. Optional `-consumer-case` and `-media-models` paths are forwarded unchanged to the inherited checks. Use models declared by your implementation, not historical defaults or invented feature rows. A passing WebSocket grade does not substitute for the chapter's actual browser, native-speech, replay, and pause exercises. Run isolated servers and browser contexts without disrupting the orchestrator's live GUI. The same frozen core CLI must pass Chapter 7 and Chapter 1, including the real delayed-response phases.
123: 
124: The Edition 3 command requires an explicit target and currently supports Chapters 1, 4, 5, 6 and 7. Its Chapter 1 mode preserves the inherited conversation checks and adds delayed-header/body lifetime gates. Check the report's `passed` field and exit status, not score alone: required zero-point checks can fail a 100/100 score. Those real delayed probes take about two minutes, so a 90-second suite timeout is invalid. Omitting the path from the older Edition 2 command silently targets an old solution and is not a valid student result. A grader result does not replace required live exercises, test-sensitivity evidence or independent acceptance. Credentials remain in the authorized provider loader/test process, never prompts, command arguments, evidence or commits.
125: 
126: Submit one concise chapter review identifying the chapter/code/grader revisions, what you built, actual tests and live observations, failures, ambiguity, assistance and useful teaching clarifications. You may explain a suspected grader defect; you may not silently redefine success or edit the grader. The orchestrator supplies any missing public contract without old answer code. Follow the shared bounded review loop.
```

## docs/edition-3-notes/ch07-author-review.md

```text
18: 
19: 1. **Exercise: explicit Edition 3 invocation.** Replace the historical Edition 2 positional-core shorthand with a nested-command build example and the public `grade-edition3 -ch 7` recipe. The binary is the positional target; `-core`, consumer/case, workflow/media map, streaming/model, and GUI/model declarations remain separate. The measured broad `LLM_BASE_URL` override and missing media declaration are explained as invocation errors, not tool-stream failures. The GUI build explanation is retained, with its `-gui` input made explicitly required for this adapter. Existing historical fixture-model qualification remains untouched.
20: 2. **Observation section: requested calls are not execution.** The initial independent review's R1 showed that a refused mixed batch appeared to start in both live observations and replay. The addition states the semantic boundary: only admitted calls get `ToolDispatched`; requested calls and refusal results remain visible without implying execution. It does not prescribe the correction's storage fields, migration mechanism, or admission-marker implementation.
21: 3. **Speech section: complete readable tool results.** The initial independent review's R2 showed a page that displayed a full tool result but spoke only its generated caption. The addition distinguishes on-demand full readable result/error speech from automatic abbreviated labels, retaining the stable-description exception for opaque media. It does not require speech to re-read artifacts automatically or introduce a new speech queue policy.
22: 
23: These are local clarifications of existing contracts, not new architecture. The radio/receiver and mixer explanations, diagrams, named lifecycle steps, turn-completion separation, pause/resume/interrupt semantics, history/replay/parent-link examples, tests and historical failures, exercises, pedagogy, and deferred-work boundary remain. There is no cleanup/pruning pass and no deleted explanatory paragraph. The only removed prose is the superseded positional-core/default-GUI invocation wording; its GUI build meaning is retained at the same location.
24: 
25: ## Evidence acceptance and limits
26: 
27: - [Initial student report](ch07-student-review.md): preserves the early broad endpoint override, missing media-model map, pause-fixture failures, actual peer-EOF reset defect and correction, plus the eventual coherent student bundle. These attempts are not recast as successes.
28: - [Harness reconciliation](ch07-student-harness-reconciliation.md) and [pause-fixture acceptance](ch07-pause-fixture-accepted.json): the original pause predicate had both protocol-completeness and step-boundary mismatches. The repaired provider-faithful predicate used the intended completed-call boundary; passing repaired fixture evidence does not erase original failures or substitute for production-provider/native evidence. No author fixture change was made.
29: - [Initial independent review](ch07-code-review.md) and [root decision](ch07-code-review-root-decision.json): R1 and R2 required one correction. C1, page-hide speech interruption, was a non-blocking clarification comment; this author does not turn it into a new hidden-tab cancellation requirement or rewrite the existing queue policy.
30: - [Correction report](ch07-student-r1-review.md), [independent correction review](ch07-code-review-r1.md), and [root correction decision](ch07-code-r1-root-decision.json): support acceptance of the truthful-dispatch and readable-result fixes and explain the student's focused local checks. Those implementation details are not prescribed in the chapter.
31: - [Accepted bundle](ch07-accepted.json): distinguishes accepted student code/grader/live evidence from the still-pending text phase. The author does not confer code acceptance or produce a new grader result.
32: - [Real-provider live exercise](ch07-live-exercise.md): its original launcher artifact remains failed. The unexpected *different-ID* tool call invalidated a naive no-tools shortcut, but root's correlated observations/logs established that the held completed calls did not dispatch, interruption cancelled the paused turn, and the same agent resumed with those held calls still absent. This is narrower than a clean launcher pass and narrower than a claim that the provider made no tool request. That distinction is preserved here; the chapter's existing call-ID/replay lesson needed no new audit paragraph.
33: - [Native exercise](ch07-native-exercise.md): the earlier provider/native Electron run includes human-heard speech and actual interaction, with retained failed attempts. The later corrected-bundle controlled Electron checks are useful for layout, artifact controls and routing but do not constitute a fresh paid-provider/native speech run for the corrected code. Root's acceptance and these limits are retained, not inflated into fresh author testing.
```

## docs/edition-3-notes/ch07-student-harness-reconciliation.md

```text
1: # Chapter 7 student harness reconciliation — unchanged frozen r4
2: 
3: ## Result and strict scope
4: 
5: The **one authorized repaired Chapter 7 adapter run passed: 100/100, `passed:true`, process exit 0**. The original Chapter 7 results remain failed historical results; none were replaced, renamed, reinterpreted as passing, or edited. This was harness reconciliation, **not a student source change, new candidate, formal code correction round, or whole-code acceptance**.
6: 
7: Root's supplied follow-up reports independent acceptance of a one-field fixture repair: the prior pause fixture used a legacy argument name inconsistent with the public Chapter 4 contract, while timing, assertions and scores remain unchanged. The student used that expressly authorized repaired executable, without reading private review/case bodies or changing the grader. Root's explanation is attributed assistance, not a new student investigation.
8: 
9: Authority: `/tmp/cadence-ch07-post-fixture-student-followup.txt`, 4575 bytes, SHA256 `e7258e2fe819c6294ba77cc73435b03cdd8f5f46b79705f871297815f090d1a0`, read fully. Grader `/tmp/cadence-ch07-grader-pause-fixed`, SHA256 `c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e`, verified before launch and after the run.
```

## docs/edition-3-notes/ch07-pause-fixture-accepted.json

```text
1: {
2:   "checked_at": "2026-10-17T21:52:07.149932+00:00",
3:   "manifest_entries": 70,
4:   "evidence_files": 70,
5:   "all_files_bytes_with_report": 169844,
6:   "hashes_and_exact_set": true,
7:   "single_field_diff": true,
8:   "final_turn": 900,
9:   "joined": true,
10:   "report_sha256": "24610b9a6e84de6786189eeb8c272dab90b207738d03d85db80f23f42d6ec000",
11:   "manifest_sha256": "4411a4b6eda7c3c086973e7d56bcfe002c6843dd62abe7eafa8b50a715296127",
12:   "status": "accepted_fixture_repair_only",
13:   "scope": "driveGUIPause argument delay -> ai_callback_delay; assertions, timing, scores unchanged",
14:   "student_commit": "48e1d12d9c396b5988c0ce6eccdef67e581872ff",
15:   "student_accepted": false,
16:   "formal_student_correction_rounds_used": 0,
17:   "corrected_grader": {
18:     "path": "/tmp/cadence-ch07-grader-pause-fixed",
19:     "sha256": "c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e"
20:   },
21:   "full_route": {
22:     "score": 100,
23:     "passed": true,
24:     "exit": 0
25:   },
26:   "sensitivity": "compiled real tool-admission bypass produces actual GUI_SECOND effect before unpause, score0/20 exit1; restored117file source passes20/20 exit0",
27:   "native_paid_permission": false
28: }
```

## docs/edition-3-notes/ch07-live-exercise.md

```text
1: # Chapter 7 live GUI exercise — root acceptance and launcher false negative
2: 
3: Bill explicitly approved one bounded Claude Sonnet 4.6 GUI exercise: maximum four upstream requests, the frozen GUI's existing4096 output tokens/request, isolated muted Chrome, no automatic paid retry. Actual source remained463996b765631d31923a5726036a34942a67cefd. The pinned launcher-root-r1.mjs and unchanged policy were verified before any credential access/process launch; root read the key only through the approved process-memory path. No install/restart/shared browser/microphone action occurred.
4: 
5: ## Actual run
6: 
7: `live-approved-01` ran2026-10-18T04:34:33Z,8.627seconds,launcher exit1. Both requests used `claude-sonnet-4-6`/stream=true/max_tokens4096 and returnedHTTP200. First response requested read_file; second contained the actual result and completed the answer. Recorded usage totals3553input/179output tokens. No third request, fallback, retry or discovery. Both owned GUI/Chrome processes exited0; credentials were not exported to children or evidence.
8: 
9: Actual page observations show Connected/in_flight and a partial Assistant artifact before completion; final Connected/idle DOM shows visible explanation, read_file card and exact file contents, and final answer. The result correctly explains projectLantern Orchard, three reviewers, blue folder. Those facts are not encoded as an arbitrary exact answer phrase requirement.
10: 
11: ## Failed launcher assertion retained; independently reconciled evidence
12: 
13: The launcher reported passed=false solely because actualReadResult stayedfalse. Its event predicate wrongly required `tool_finished.name === 'read_file'`. The actual production event has no name field; it correlates by stable call_id and part_id. Saved gui.log establishes exactly one read_file dispatch for exercise-note.txt, exactly one nonerror completion with the SAME call/display IDs, and byte-identical result/fixture/job output. The final DOM displays that result. The second authenticated provider request receipt contains that actual tool result, and the final answer correctly uses all three facts. Root saved the exact correlated safe events and full semantic reconciliation, without rerunning or editing student code. The original exit1 and passed=false remain FAILED launcher outcomes; this is not a retroactive green report.
14: 
15: Root accepts the existing real-model/live-GUI exercise requirement on that combined observed evidence. This is not universal provider/GUI proof or a fresh native speech test. Prior Bill-heard audio remains qualified earlier evidence; screenshots were captured but root does not claim an unperformed visual-image review. The plan's prose called the fixture chapter7-note.txt; actual pinned policy, prompt, launch and result consistently use exercise-note.txt. That prose typo is a plan-record issue, not an application mismatch.
16: 
17: ## Evidence and portability
18: 
19: `ch07-live-exercise-evidence/` holds actual invocation/exits, original stdout/stderr, sanitized receipt, extracted correlated tool events, root reconciliation and local-artifact hash metadata. Its text index excludes itself and the two PNG screenshot copies; screenshots and raw runtime/API captures are local-only pinned artifacts, not automatically chapter-commit contents. No private settings, credentials, raw authenticated headers or browser profile are publication material.
20: 
21: The failed original launcher predicate was not changed after the paid run, and no paid permission remains to retry automatically. Future implementation changes would need their own qualified evidence rather than inheriting this run as if it used different production bytes. Code acceptance, tag and later author/text closure remain separate recorded decisions.
```

## docs/edition-3-notes/ch07-live-exercise-evidence/live-approved-01-correlated-tool-events.json

```text
1: [
2:   {
3:     "at": "2026-10-18T04:34:37.21104Z",
4:     "message": {
5:       "call_id": "toolu_013yPgvEruEBrGiQTHcM5GT4",
6:       "input": {
7:         "path": "exercise-note.txt"
8:       },
9:       "name": "read_file",
10:       "part_id": 3,
11:       "type": "tool_dispatched"
12:     }
13:   },
14:   {
15:     "at": "2026-10-18T04:34:37.212059Z",
16:     "message": {
17:       "call_id": "toolu_013yPgvEruEBrGiQTHcM5GT4",
18:       "is_error": false,
19:       "part_id": 3,
20:       "result": "Project: Lantern Orchard\nReviewers needed: 3\nFolder: blue\n",
21:       "type": "tool_finished"
22:     }
23:   },
24:   {
25:     "at": "2026-10-18T04:34:39.36653Z",
26:     "message": {
27:       "error": "",
28:       "text": "Here's a brief summary of what the file contains:\n\n- **Project:** *Lantern Orchard* \u2014 this is the project the note describes.\n- **Reviewers needed:** *3* \u2014 three reviewers are required for this project.\n- **Folder to use:** *blue* \u2014 the designated working folder is named \"blue\".",
29:       "type": "turn_ended"
30:     }
31:   }
32: ]
```

## docs/edition-3-notes/ch07-live-exercise-evidence/live-approved-01-root-reconciliation.json

```text
1: {
2:   "checked": "2026-10-18T04:36:03.789070+00:00",
3:   "launcher_exit": 1,
4:   "launcher_passed": false,
5:   "launcher_failed_assertion": "actualReadResult required a name on tool_finished which has only stable correlation IDs",
6:   "same_call_and_display_id": true,
7:   "result_matches_fixture_and_job_output_exactly": true,
8:   "actual_partial_DOM": true,
9:   "actual_final_DOM_result": true,
10:   "next_authenticated_request_contains_actual_tool_result": true,
11:   "root_semantic_review": "Final answer correctly identifies Lantern Orchard, three reviewers, and blue folder; no exact answer phrase imposed.",
12:   "live_exercise_root_accepted": true,
13:   "student_production_change": false,
14:   "paid_retry": false,
15:   "requests": 2,
16:   "input_tokens": 3553,
17:   "output_tokens": 179,
18:   "owned_processes_exit_zero": true,
19:   "qualification": "Launcher failure retained; evidence reconciliation is not a retroactive launcher PASS. Screenshots captured, not claimed visually reviewed here. Default native listening belongs to earlier authorized r4 evidence."
20: }
```

## docs/edition-3-notes/ch07-native-exercise.md

```text
1: # Chapter 7 native browser exercise — root execution
2: 
3: Bill approved the audible exercise and subsequently authorized autonomous retries. This exercises unchanged student commit `48e1d12d9c396b5988c0ce6eccdef67e581872ff`; it is not whole-code acceptance or paid-provider proof.
4: 
5: ## Result
6: 
7: The third launcher attempt exited 0. In a disposable headed Chrome profile, the actual production speech queue emitted native start for “Harbor lights shine softly.”, cancellation through the page's Interrupt control, and native start then end for “Quiet gardens rest.”. At the final observation speech was neither speaking, pending nor paused and the GUI speech-pause indicator was released. No native error was used as completion. Human listening remains a separate pending observation. Default voice service locality is unknown and was explicitly consented to; no production voice selection or callback replacement was added.
8: 
9: All seven served assets matched pinned student bytes. One fixed local-provider request was handled; no paid request, credential read or microphone operation occurred. GUI/Chrome PID and listener ownership were checked before attachment. Both owned child processes exited 0 and the local fixture closed. No shared browser or live GUI was attached or restarted. Runtime receipt duration was 2.974 seconds; the enclosing command was 12.963 seconds including harness-process teardown. No latency benchmark is claimed.
10: 
11: ## Retained failures and scope correction
12: 
13: Attempt 1 exited 1 before GUI/Chrome/speech: a redundant check of the harness's own process compared canonical Node executable spelling with a symlink invocation. Attempt 2 exited 1 before speech because the harness's blanket DNS block also blocked loopback, triggering a failed navigation/retry. Both receipts are preserved. Removing the redundant self-command check and exempting loopback corrected the harness only. No student source changed. Root also guessed the wrong speech-asset path and then a nonexistent policy root key during report preparation; both reads failed before report writes, and actual manifest paths were used afterward.
14: 
15: The earlier requirement to prove default-voice locality was root-imposed test policy, not the chapter contract. Bill specifically reminded root to avoid invented requirements that obstruct efficient execution of the existing book. Root withdrew that restriction rather than require a student feature. Existing process isolation and fixed local-provider boundaries remain; this is not an OS-network-isolation proof.
16: 
17: ## Evidence
18: 
19: `native-result-summary.json` records the exact event sequence. Each attempted launcher and receipt remains separate; `attempt-3-receipt.json` contains full event and process records and `attempt-3-command.json` the actual command, environment and exit. The original student proposal/index remain unchanged. `SHA256SUMS` covers this report and the evidence directory except the checksum file itself.
```

## docs/edition-3-notes/ch07-native-exercise-evidence/human-listening.json

```text
1: {
2:   "recorded_at": "2026-10-17T23:58:22.793384+00:00",
3:   "source": "Bill direct observation",
4:   "exact_quote": "I heard \"Quiet Gerdens Rest\"",
5:   "test": "attempt-4",
6:   "intended_phrase": "Quiet gardens rest.",
7:   "qualification": "Human confirms audible phrase; not a transcription-accuracy metric or a claim that first canceled phrase was heard",
8:   "native_receipt_sha256": "8d6f3182454659469c5b8d67981d38be68c2af6dbe7e294aded5301bfa5a236f"
9: }
```

## docs/edition-3-notes/ch07-student-r1-review.md

```text
35: 
36: `ch07-student-r1-source-manifest.json` pins every source file, archive, artifact, exact declarations, command receipts, grader and asset identity. The build-info logs record the compiled CLI/GUI module and VCS metadata. Source did not change after freezing; checks preceded/followed each sequential command. Both graders used the CLI hash shown above, verified again after both completed.
37: 
38: The public adapter builds optional module executables internally from the frozen module paths. Its temporary internal GUI binary hash is not available in the report; the standalone GUI hash is **not substituted** for that temporary binary. The source is pinned for both, and the standalone artifact is separately pinned.
39: 
40: ## R1 — truthful execution projection
41: 
42: **Original failure:** Chapter 3 intentionally retains a `tool_called`/`tool_returned` pair for each explanation-refused proposal. The r4 GUI treated every `tool_called` as execution and therefore emitted/summarized refused proposals as `tool_dispatched`, both live and on replay. Removing the required called/returned pairs or filtering selected names/error strings would be incorrect.
43: 
44: **Narrow correction:**
45: 
46: 1. `internal/common/common.go` adds optional data-only `ToolEvent.Dispatched`. It records admission, not successful completion.
47: 2. `internal/agent/dispatch.go` stamps that fact only in the actual admitted `prepareCall` path, not the explanation-refusal branch. Actor decisions, job ownership, refusal feedback, identities and one-shot limits are unchanged.
48: 3. `gui/hub.go` requires positive admission evidence: the durable bit, an actual retained-job record, or the real admission observation already seen for that identity. A proposal lacking that evidence renders as a final call artifact, **not** an execution lifecycle envelope. Correlated results/errors remain readable. Existing exact-once dispatch bookkeeping remains in place.
49: 
50: The real Agent fixture first returns an unexplained three-call proposal (file, process, control), then an explained admitted retry with reused provider call IDs. Before retry, actual file/process effects are absent, the provider receives all three refusal results, live/reconnected messages claim zero starts, and identity pairing remains intact. The retry performs the actual file/process work and admits the control; there are three core admissions, two real retained-job identities, three live dispatches and three replayed dispatches. Distinct durable display IDs separate the reused provider IDs.
51: 
52: The regression traverses local fake HTTP, a real public Agent/job owner, actual hub subscriptions, **real WebSocket frames**, serialized history load into a fresh Agent, and the controlled production page. No private actor state or injected GUI dispatch creates this proof. Complete raw wire/history captures stay local; `focused-safe-summary.json` and the pinned test/assertion logs are portable.
53: 
54: **Legacy qualification:** prior job-backed admission records retain execution proof. A legacy control with neither a job, durable admission bit nor current real observation cannot be truthfully classified as executed: its call and result remain artifacts, but replay does not invent a start. Newly admitted no-job controls have the durable bit and replay correctly. No existing history is rewritten or deleted.
55: 
56: ## R2 — complete on-demand tool speech
57: 
58: **Original failure:** `tool_finished` appended visible output but left the artifact's readable text as its call summary. The actual speaker button therefore queued a caption while omitting the visible result.
59: 
60: `gui/web/gui/tool-readable.js` projects the readable call context/arguments and accumulated result or error as plain text. `gui/web/gui/app.js` keeps tool-call and result state together and refreshes the artifact's speaker text on both result arrival and later call-final updates. The short automatic dispatch summary remains separate and concise; receiving a tool result does **not** automatically narrate it. No native speech API, speech queue, parser, server result clipping, tool executor or renderer subsystem was redesigned.
61: 
62: The controlled test evaluates the **actual page module and its actual imported SpeechQueue**, delivers real-handler messages, and clicks the actual artifact speaker closure. It checks ordinary `read_file` text including literal markup, specialized `run_command` JSON with stdout/stderr and an error, full result inclusion, metadata exclusion, a late call-final update, replay/reconnect silence and on-demand recovery, empty/opaque silence, cancellation, stale callbacks, current error release, and pause/unpause intents. It also feeds actual core/hub/WebSocket refusal receipts through the page and verifies no automatic execution summaries while errors remain readable on demand.
63: 
64: The DOM, WebSocket client and synthesis boundary in that test are explicitly controlled test objects. These are **not native or browser/listening observations**. The test exercises production page/speaker logic rather than a pure-helper assertion or a receipt-only claim.
65: 
66: ## Failure sensitivity, not merely green tests
67: 
68: - Both final regressions run unchanged against an isolated copy of **original r4 production**. They compile/evaluate and fail their intended assertions: false dispatches and omitted spoken results. `original-r4-r1-red.log` and `original-r4-r2-red.log` retain these failures. The R2 script stops at its first failed ordinary-result assertion; it is not claimed to have completed every later check on original code.
69: - A compile-valid R1 mutation claims every proposal was dispatched. It fails the refusal/admission assertions, including the wire path, while actual refused process effects remain absent.
70: - A compile-valid R2 mutation discards the stored result while leaving visible output rendering intact. The actual speaker assertion fails because the result is absent.
71: - Each mutation lives only in an excluded isolated copy, records before/broken/restored hashes, and passes after exact restoration. The working/frozen source is never mutated. These controls were repeated sequentially after freeze and are in `frozen-control-results.json` and their command/output receipts.
72: 
73: ## Frozen sequential validation
74: 
75: `run-frozen-sequential.py`, per-command receipts and `sequential-completed.json` retain the exact cwd, closed environment, command, start/end, exit and output hashes. The completed sequence contains 26 commands: 24 expected zero exits and two **intentional semantic-failure** controls, followed by their passing restorations. There were no concurrent adapter runs or source edits during these phases.
```

```text
90: 
91: Exactly one final Chapter 7 and one final Chapter 1 adapter invocation ran in this correction. The same root-approved repaired grader executable is SHA256 `c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e`. The inherited consumer/workflow/streaming/GUI declarations and existing media map were reused unchanged; no new grader mode/model declaration or case assertion was invented.
92: 
93: Chapter 7 command (cwd `/Users/bill/projects/ensemble`):
94: 
95: ```sh
96: /tmp/cadence-ch07-grader-pause-fixed -ch 7 -json -core /Users/bill/projects/ensemble-edition-3-student -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer -consumer-case /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json -workflow /Users/bill/projects/ensemble-edition-3-student/examples/ch06 -stream-exercise /Users/bill/projects/ensemble-edition-3-student/examples/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble-edition-3-student/gui -gui-model claude-sonnet-4-6 -media-models /Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-evidence/declared-media-models.json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli
97: ```
98: 
99: Chapter 1 command, same cwd:
100: 
101: ```sh
102: /tmp/cadence-ch07-grader-pause-fixed -ch 1 -json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli
103: ```
104: 
105: The raw machine reports remain verbatim, including the Chapter 7 adapter's inherited Edition 2 title. The current repaired adapter's five checks cover observation types, replay, GUI log, pause-held dispatch and corrected Chapter 6 parity. A passing score is not whole-code acceptance.
106: 
107: ## Native boundary and exact carryover
108: 
109: The correction brief states that prior completed native/heard-by-Bill evidence belongs to r4. It is **not relabeled** as a run of this candidate. `native-evidence-applicability.json` compares exact bytes against r4:
110: 
```

## docs/edition-3-notes/ch07-code-review-r1.md

```text
36: ## R1 — truthful tool admission: closed
37: 
38: ### Source assessment
39: 
40: `internal/common/common.go` adds an optional durable `ToolEvent.Dispatched` fact. `internal/agent/dispatch.go` sets it on the actual admission path, including admitted no-job controls. Refusal/pairing paths do not acquire that fact merely because a `tool_called` record is retained.
41: 
42: `gui/hub.go` now projects a dispatch only with affirmative admission evidence: the explicit fact, job-backed legacy evidence, or an actual observed dispatch. An ambiguous job-less historical `tool_called` remains a readable call artifact and retains its result; it is not promoted into a start. This is narrower than filtering every nil-job call, which would incorrectly hide current admitted controls.
43: 
44: ### Independent real boundary
45: 
```

```text
69: 
70: ## R2 — complete readable tool-artifact speech: closed
71: 
72: `gui/web/gui/tool-readable.js` selects only the readable call context/arguments and accumulated result/error. `app.js` keeps those values together and refreshes the speaker text on result arrival and later call-final updates. The short automatic dispatch summary remains separate. A completed result does not trigger an unsolicited full-result narration.
73: 
74: The independent page check executes the actual unchanged page modules in a Node VM with inert DOM/WebSocket/speech boundaries. It feeds the real protocol captures above and selected ordering/content cases; it is not a replacement page implementation and is not native-browser evidence.
75: 
76: All nine final corrected cases pass:
77: 
78: 1. Long readable result/arguments survive result arrival and a later call-final update. Automatic summary remains 22 characters; on-demand speech includes the complete call/result (3,990 characters in the fixed fixture).
79: 2. The corresponding error case includes complete error content and its error label (3,995 characters), not only the compact summary.
80: 3. Cancellation and stale callbacks do not release a newer utterance's pause; actual current completion does.
81: 4. The analogous `onerror` path has the same generation/gate behavior.
82: 5. Empty/signed/opaque content is not spoken.
83: 6. Actual refused-wire frames produce no automatic execution summaries, and all six error artifacts speak their complete readable result/context.
84: 7–9. Fresh reconnect, rehydrated history and the conservative legacy view each preserve complete readable speech for all twelve result artifacts.
85: 
86: The unchanged final check also runs on original r4. Six relevant cases fail for the intended R1/R2 reasons: full result/error speech is still the 22-character summary, result speakers omit their payloads, and refused calls generate execution summaries. The three unchanged opaque/queue-generation guards still pass. Returning to corrected source and fresh corrected wire passes all nine cases again.
87: 
88: The final result check compares the **complete wire result** plus the appropriate error/result label. It does not demand identical punctuation to the DOM's concatenated caption. A preliminary overly strict caption assertion and the inert DOM's initially incomplete `textContent` behavior are retained as reviewer errors, not candidate failures.
89: 
90: ## Fresh final gates
91: 
92: All gates ran sequentially in owned scratch with a clean, explicitly recorded environment, local fake transports, no credentials/paid requests, and no outer deadline that could cancel the genuine long waits. The full root race suite was allowed to finish. No new browser, microphone or native speech run occurred.
93: 
94: | Gate | Exit | Wall time | Recorded outcome |
95: |---|---:|---:|---|
96: | `root-race` | 0 | 140.489 s | PASS (test/vet exit; no grader score) |
97: | `root-vet` | 0 | 0.909 s | PASS (test/vet exit; no grader score) |
98: | `gui-race` | 0 | 1.912 s | PASS (test/vet exit; no grader score) |
99: | `gui-vet` | 0 | 0.168 s | PASS (test/vet exit; no grader score) |
100: | `consumer-race` | 0 | 0.128 s | PASS (test/vet exit; no grader score) |
101: | `consumer-vet` | 0 | 0.107 s | PASS (test/vet exit; no grader score) |
102: | `workflow-race` | 0 | 1.381 s | PASS (test/vet exit; no grader score) |
103: | `workflow-vet` | 0 | 0.112 s | PASS (test/vet exit; no grader score) |
104: | `stream-race` | 0 | 0.122 s | PASS (test/vet exit; no grader score) |
105: | `stream-vet` | 0 | 0.105 s | PASS (test/vet exit; no grader score) |
106: | `chapter7-fresh` | 0 | 14.573 s | 100/100; `passed=true` |
107: | `chapter1-fresh` | 0 | 124.066 s | 100/100; `passed=true` |
108: 
109: The root remains a library with nested `cmd/ch01`; the separate importing GUI and `examples/ch04-consumer`, `examples/ch06`, and `examples/ch07` modules were not flattened. Race/vet covered the current submitted tests **before** adding the reviewer-only test file to the private copies.
110: 
111: ### Exact final grader routes
112: 
113: Both use the same exact corrected frozen CLI above, not the old r4 CLI. Cwd is `/Users/bill/projects/ensemble`.
114: 
115: ```sh
116: /tmp/cadence-ch07-grader-pause-fixed -ch 7 -json -core /tmp/cadence-ch07-code-r1-rereview-independent/source -consumer /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch04-consumer -consumer-case /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch04-consumer/consumer-case.json -workflow /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch06 -stream-exercise /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch07 -stream-model claude-sonnet-4-6 -gui /tmp/cadence-ch07-code-r1-rereview-independent/source/gui -gui-model claude-sonnet-4-6 -media-models /Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-evidence/declared-media-models.json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli
117: 
118: /tmp/cadence-ch07-grader-pause-fixed -ch 1 -json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli
119: ```
120: 
121: Chapter 7 freshly reports **100/100 and `passed=true`**. This includes the already accepted corrected pause fixture; no old entry/adapter/pause-sensitivity campaign was repeated.
122: 
123: Chapter 1 freshly reports **100/100 and `passed=true`** in **124.066300 seconds**. Its `generation_headers` and `generation_body` checks each report receipt of the complete answer/usage after a **1m2s** delay. No separate per-phase sub-second elapsed measurements are emitted by that JSON report, and none are invented here. Both genuine 62-second gates completed. This is a fresh same-CLI result, not a carried previous-commit score.
124: 
125: ## C1, native applicability, and real-model boundary
```

## docs/edition-3-notes/ch07-code-review.md

```text
94: A tool-result artifact's speaker must include its readable result (including readable error text where applicable), rather than only the tool name/argument summary. Preserve the safe/opaque distinction and production queue/error/cancellation behavior. Do not satisfy this by removing the speaker or by claiming a native listening test of an unrelated text artifact covers it.
95: 
96: ## C1 — page-hide release: measured, but not a blocking code correction
97: 
98: The same unchanged page harness measured this sequence:
99: 
100: 1. A nonempty draft's input event sends `pause`.
101: 2. `pagehide` runs while the socket is open.
102: 3. The handler at `app.js:77` cancels speech, clears local causes/draft, invokes `localPause`, then closes the socket.
103: 4. **`unpause` is queued before close.**
104: 
105: This is a real generated-control observation in the non-browser page harness. It is not a fresh native measurement of a remote gate after browser teardown, and it does not establish that every transport delivers the queued frame before disconnect.
106: 
107: The exact student-visible assignment `/tmp/cadence-ch07-student-task.txt`, SHA-256 `8d34f80dfce2dfc506a2dfd1b3f3a18129279227846a738914fd85619d882d2f`, line 9, included `disconnectnotauto-release/reconnectstale causesnotreassert`. Therefore this is not merely an invented private-audit preference.
108: 
109: However, the chapter itself explicitly describes server-held state, automatic local-cause reversals and stale reconnect behavior without unambiguously resolving this page-hide cleanup transition. Root checked the exact task and directed that, under Bill's prohibition on added requirements, this be recorded as a **student-task deviation / public-contract clarity note, not a third blocking correction solely on that sentence**. That disposition is retained in `contract-classification.md`. No public text or source was changed here. Root owns any clarification and subsequent correction scope.
110: 
```
