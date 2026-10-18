# Chapter 8 student review — frozen local-only candidate

## Status and identity

**Implemented and locally tested; not accepted. Browser/native/manual, real-provider/model and independent acceptance evidence remains open.** No acceptance tag, main-repository commit, push, install or publication was made.

- Student root: `/Users/bill/projects/ensemble-edition-3-student`.
- Accepted base: `463996b765631d31923a5726036a34942a67cefd`, tree `5e95dcdd716f7c25c96cc3ebcde13d000bb38670`, untouched `edition-3-ch07-r1`.
- Local freeze: `f97a6850468a9cdc6875507f2821ccf6c397654d`, tree `b3c4f005d5f43a7f20fc6227c65e1f763b93fac9`; clean tracked/untracked source status. 126 tracked files, with exact file hashes and complete archive entry membership in the source manifest.
- Frozen CLI: `.ch08-local/bin/root-ch08`, SHA-256 `160946a79ca8ebdcf17acfaf315b9710be87cd7294387eb36ac95d0062491945`. Both final grader receipts hash these same bytes immediately before and after execution. The GUI and three example executables are also indexed locally.
- Public materials revision: `2cebe5ae76b47eabfd1f1cb7e8bef95c98b7c20e`; exact chapter, seven prerequisite documents, assignment and supplied grader hashes are in `ch08-student-source-manifest.json`. The public grader was invoked as a black box, never inspected.

The seven mandatory public prerequisites were read in full during the preserved plan turn, including recovery of truncated reads. Relevant earlier Chapter 5 budget/dispatch and media-schema sections, Chapter 7 GUI/replay/pause sections, and the accepted student source supplied the remaining context. The unchanged `ch08-student-plan.md` records those reads and the initial disk stop. Root subsequently verified restored disk space and approved this bounded local implementation. No historical solutions, private grader code, coordinator decisions, other-worker outputs or credentials were read; no child agents/reviewers were used.

## What changed and why

`internal/agent` owns an atomic next-turn tool-batch policy, exposed by the small public `ToolBatchLimit`/`SetToolBatchLimit` interface. Each human turn captures it once. Positive N means N batches; zero selects the named default of 200. Launch `Config()` remains immutable; Engine remains the single-request owner. GUI settings are not core preferences.

`gui/settings.go` owns the five-field value, shared sparse-patch/disk validator, and replacement-file persistence at the process-relative `cr/settings.json`. It preserves explicit false/zero, rejects invalid types/null, normalizes unknown themes, clamps agreed ranges and retains fractional speeds. Under the hub lock it validates, saves, applies the Agent budget and then publishes a full snapshot. A failed save leaves both accepted settings and effective policy unchanged and sends `settings_error`. Startup loads before either input view is exposed; invalid/unreadable files fail visibly. Rename-based replacement is not advertised as power-loss durability.

The existing app, queues and renderers now sit in independently scrolling sidebar/chat/actions regions. The sidebar has Agents/Settings tabs, keyboard/pointer dividers and a recursive single-Agent tree driven by real states. Fifteen semantic color tokens support dark/light/system themes. Four actual controls send sparse updates and repaint only accepted snapshots; rejected or unacknowledged drafts do not become effective. Font size reaches the document, Read aloud reaches the inherited queue consumer, and every newly created utterance gets the accepted speed. Ordinary fractional rates through 2 are exact; requests above 2 map monotonically into native 2–3 and remain explicitly approximate. Existing stdin/WebSocket views, replay, FIFO reveal, speech queue, pause/hint/interrupt and provider-wire ownership were retained.

## Executed local evidence

All final commands used a whitelisted environment without provider credentials, sequentially, with exact argv/cwd/start/elapsed/exit/output receipts. No short wrapper deadline was imposed on real lifetime probes.

| Gate | Observed result |
|---|---|
| Core full tests / race / vet / build | Exit 0; full tests 129.766s, race 151.201s. |
| GUI tests / race / vet / build | Exit 0; tests 0.263s, race 1.739s. |
| Consumer, workflow, streaming modules | Test, race, vet and build commands exit 0. These example modules report **no test files**, not invented case coverage; actual application behavior is exercised by inherited grading. |
| Focused owner tests and Node checks | Pass: captured 2/18/zero-default budget, malformed/ranged/false/zero disk/patch settings, two-peer broadcast, restart and save failure; queue ordering/cancel/stale callbacks, fractional/extended next-utterance rates, pane following, accepted control effects, theme media changes, tabs/dividers/tree, HTML IDs and asset syntax. JSON targeted receipts retain named Go test/subtest results. DOM/speech boundaries are inert fixtures, not a browser or native engine. |
| Built GUI process fixture | Exit 0, 0.652s on frozen GUI: terminal turn retained budget 2 during a live WebSocket update; saved 18 survived process restart and allowed 18 GUI-requested real file-tool batches; zero allowed 20 terminal batches and survived reconnect/restart. Two real WebSocket peers and stdin used the same Agent with a loopback streaming HTTP fixture. |
| Final Chapter 8 black-box grader | `passed=true`, 100/100, exit 0, 14.405s. All settings checks and inherited Chapter 7/6/5/4 parity passed. |
| Following Chapter 1 on identical CLI | `passed=true`, 100/100, exit 0, 124.058s. Separate required zero-point `generation_headers` and `generation_body` checks both passed after actual 62-second delays. The earlier Chapter 1 run also passed in 124.053s and is retained. |

Stream/GUI model: explicitly declared `claude-sonnet-4-6`. The public grader-only media mapping selects the existing `claude-sonnet-4-6`, `gpt-4.1-2025-04-14` and `gemini-3.8-flash` capability rows. This is local fixture input, not provider discovery or a claim about current live availability.

Three named faults in disposable source copies were each caught: reading the live policy during an active turn changed its allowance; mutating settings before a failed save leaked the candidate; forcing utterance rate to 1 defeated the real queue's fractional-rate assertion. Mutations, nonzero receipts and restored unmutated passes are retained. The frozen source was never mutated for sensitivity testing.

## Failures, assistance and limitations

Two initial process-fixture attempts failed: guessed launcher/protocol assumptions, then a 17-byte WebSocket handshake key. Inspection of the actual launcher/socket contract corrected the fixture; attempt three and the frozen rerun passed. The first Chapter 8 grade returned 80/100, `passed=false`, exit 1: its default OpenAI image model `gpt-5.4-mini` was not a declared predecessor capability row. The documented `-media-models` input resolved this without changing source, binary or grader. The failed report remains intact; score alone was not treated as success.

Source/diff inspection also caught early unmatched browser-edit anchors before freeze. One Git pager was stopped; subsequent commands used `--no-pager`. A preliminary lightweight Node check overlapped the first core run; the complete frozen gate chain was then sequential. Unqualified example builds left three generated executables, which were moved to the ignored local artifact directory and preserved. These disclosures and raw failure diagnostics are retained, not hidden behind successful retries.

Root supplied disk recovery and scope/continuation authorization, not implementation or private answers. Storage started this turn at 6.2 GiB and remains 5.98 GiB free at packaging; no shared cache or accepted evidence was deleted. Development logs, raw loopback fixture HTTP/event logs and binaries stay local; the portable index records their hashes. Historical Chapter 7 evidence remains prior-candidate evidence and is not relabeled as a Chapter 8 browser/model run.

Useful teaching clarifications: declare inherited media rows even when using a newer chapter grader; distinguish launch configuration, next-turn policy and an already captured allowance; distinguish accepted snapshots from a control draft; and do not count controlled speech setters or DOM fixtures as audible/visual proof. The observed rate transform is a policy, not a native-duration guarantee.

## Minimal separately authorized continuation

1. One isolated headed browser session with two tabs against this frozen GUI and an owned local fake provider; one app restart. Inspect all three panes at normal/narrow widths, light/dark/system contrast, keyboard tabs/dividers, actual fractional settings/reload/reconnect, tool routing and independent scroll following. Exercise inherited replay and prompt/hint/pause/cancel/unpause/interrupt behavior and capture readable screenshots/console/transport evidence.
2. At most six short native utterances using baseline and requested 1.25, 1.75 and 2.5 rates, with direct native-event/utterance inspection and a human listening observation. Verify next-utterance application, cancellation/replay behavior and label extended rates approximate. **No microphone exercise is proposed or authorized.**
3. If approved separately, one useful real-provider coding/tool turn on an explicitly approved declared model, at most three generation requests and one tool batch, no retries/discovery/crossover campaign. Credentials only through the authorized loader; retain sanitized request/usage evidence. This does not waive any inherited live-model obligation the coordinator requires beyond this minimal case.

None of those browser, native, microphone, credential, discovery or paid-provider actions has occurred in this turn. Independent review/acceptance must still follow. The portable inventory includes this review and both external manifests; it explicitly names its own self-hash exclusion. No root-owned sealed freeze record was created.
