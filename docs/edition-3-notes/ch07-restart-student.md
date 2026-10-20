# Chapter 7 — R4 accepted; code gates PASS; operational v4 HOLD

**4 of 10 production correction rounds used. Source unchanged and HELD.** Root/SAME reviewer accepted R4 and the repaired grader. Latest authorized operational v4 stopped on missing native START observation; no automatic retry or source/helper repair followed.

## Valid code/grade gates

Existing `/tmp/e3-ch07-final/frozen-cli` passed corrected Ch7 and FIRST same-CLI Ch1: **100/100, `passed:true`, exit0** both. Ch1 took125s and explicitly completed answer/usage after genuine62s headers +62s body, without a shortened/killing outer deadline. Reports: `grade-ch07-v3.{json,stderr.log}`, `grade-ch01-v3.{json,stderr.log}` in that directory. CLI SHA-256 unchanged before/between/after: `d61fb25e94bb9e3a1a2b71373623b36a905806687d4ca5effa224e894edbbd9b`. Original35 and all failure/pass logs retained. Root/module results carry; GUI race/vet freshR4. No grader internals read.

One fresh operational GUI build succeeded at `/tmp/e3-ch07-operational-20261020/gui`; all later attempts symlinked the same binary, with no rebuild or production change. R4 remains only bounded existing-writer Pong handling and empty chunk serialization. `/tmp/ch07-r4-{focused,gui-race,gui-vet}.log` pass; pre-fix and earlier R3/fixture-adaptation evidence remains preserved.

## Latest actual run: operational v4

Root prospectively authorized ONE v4 setup continuation after v3 was saved/joined. `/tmp/e3-ch07-operational-20261020-v4/run-v4.py` and `operate-v4.mjs` retained the scenario/caps/flags, removed background rAF waiting, and used actual enabled **`HTMLElement.click()` normal DOM activation**. Rect/hit-test/visibility were diagnostic only. This is **programmatic DOM interaction, not physical mouse or OS focus proof**; no direct production callback or synthesizer handler was invoked/replaced to manufacture results.

`operation-v4.log`: GUI54724, Chrome54725, owned fresh profile, dynamic loopback CDP56116; ownership/listeners verified before CDP; exactly two owned pages. Overall run approximately22s, **exit1 at the first unexpected failure**.

Actual local positives:
- Early streamed partial reached the real GUI; partial reconnect retained content/display identity and did not speak.
- Cross-tab Pause and **actual Continue delivery** passed. `work/gui.log` records explicit pause then unpause at01:25:21.610. Immutable `local.jsonl` records `local-edit` requested, dispatched, returned successfully; the real built-in file job changed owned `work/exercise-note.txt` to `status: local verified`.
- Held response2 received the public terminal hint; request3 actually contained that hint. Three finite local requests completed through real Agent/jobs.
- Completed reconnect retained card identities/DOM uniqueness and remained silent. Legitimate same-card metadata updates were not rejected.
- Two safe **page-only** screenshots captured, not viewed: `compact.png`85989B and `expanded.png`92741B. Normal1280×900 viewport. Capture/geometry is not visual inspection; **no delegated visual assessment yet**.

Native attempt and stop:
- One short Speak request queued/claimed in the actual application, with actual Pause. A passive additive native listener observed **END** for “Native queue check. A short harmless local fixture.” at elapsed4315.9ms. Real wire then recorded `speech_done/status:done` and unpause.
- The initial **START was not observed** by the post-construction/queryObjects polling listener. Controller stopped after its20s observation allowance: `native initial START not observed`. This is an observer timing gap, **not proof that native/app speech failed**. No retries, synthetic native events or callback replacements occurred.
- Cancellation/recovery and later START+END were **not attempted/passed**. One short phrase consumed; two remain under the original three-phrase ceiling. Chrome had `--mute-audio`; human hearing, physical silence and default-voice locality are not asserted.

Cleanup: GUI/Chrome **exit0**; `ps` for54724/54725 and `lsof`18147/18148/56116 empty. `key_read:false`, `fake:3`, `paid:0`. No listeners displaced/shared cleanup. Owned profiles/raw evidence retained.

## Unfinished and remaining limits

Not reached: intentional terminal interrupt/EOF phase, native cancel/recovery, credential read, model GET discovery, delegated screenshot assessment, or live paid GUI flow. There have been **0 key reads,0 discovery calls,0 paid POSTs**. Existing screenshots have reached the total ceiling of2 and can be reused, not recaptured. Five local fake requests total across actual v2/v3/v4 attempts (1+1+3); original finite ceiling8. One native phrase completed; original ceiling3 phrases/90s remains. Any continuation needs root authorization and must carry this consumed evidence/budget, not silently reset it. Paid ceiling remains4 native1024 POSTs: one delegated image assessment plus<=3 actual unchanged GUI requests, active GUI provider phase<=90s, only after local prerequisites pass.

Earlier failures retained: original `operation.log` refused prelaunch after preparation consumed its early build-timed allowance; root authorized invocation-timed v2. v2 mixed hint-before-dispatch with tool-effect expectation and had no delivered Continue. Root authorized split-scenario v3; its rAF wait on a background page timed out before Continue delivery. v4 resolved that delivery and scenario issue. These were separately authorized setup continuations, **not production rounds or native/provider passes**.

Initial helper read accidentally included old assertion context; disclosed immediately. Root classified a minor nonsecret scope deviation and explicitly permitted helper-context reading as needed. No old assertions/spec models/counters/permissions adopted, adjacent old raw data/other solutions not read, real settings/key file never opened.

Clean env/owned HOME+TMP/known cache and Delve; GUI/Chrome/Node explicit clean env, no inherited real keys/LLM_BASE_URL. No MIC/OS focus/volume/settings/install/live8088 action, children, commits/staging, manuscript/grader edits or private/reference source reads. **STOP/HOLD: preserve partial real evidence without claiming complete native, visual or paid validation.**


## FINAL — authorized operational continuation completed; source HELD

**This section supersedes the v4 unfinished status and remaining-budget statements above. Historical failures and STOP records are preserved, not rewritten. Production is unchanged; 4 of 10 correction rounds used.** Root/SAME reviewer inspected the actual subsequent scripts/logs.

- **v5 local/native completion:** `/tmp/e3-ch07-operational-20261020-v5/operation-v5.log` records passive listeners attached before the unchanged real native speak call, with trusted START → Cancel audio → trusted `interrupted`, followed by trusted START+END and cleared pause/queue. Application/synthesizer callbacks and handlers were not replaced. The remaining real terminal `/interrupt` and EOF-detachment checks passed. Together with carried v4 evidence, this completes the local pause/Continue, real file job/hint, partial/completed silent reconnect, native cancellation/recovery and terminal scope.
- **Delegated visual assessment:** v5 reused the two existing safe page PNGs, read the explicitly authorized credential only in proxy memory, and made the one GET that listed `claude-sonnet-4-6`. Its one native-1024 screenshot POST completed (`end_turn`, 3137 input/767 output tokens). `visual-response.json` contains the model’s strengths and concerns. The automatic issues STOP remains an actual historical STOP. Root/SAME reviewer subsequently reconciled those concerns as **no actual production blocker** and required no source edit. This is a **DELEGATED assessment**, not direct pixel viewing by this agent/root.
- **v6 real GUI PASS:** under separate prospective authorization for one additional process-memory key read, `/tmp/e3-ch07-operational-20261020-v6/operation-v6.log` records **exit0**, three actual unchanged GUI POSTs to `/v1/messages`, all HTTP200, requested/returned `claude-sonnet-4-6`, native cap1024 each, **9.285s active phase**. No discovery, fake, native or screenshot exercise was repeated. Actual `read_file` then `edit_file` calls had successful paired results (call/result seq5/8 and12/15). The owned file exactly matched the requested status replacement with other content preserved. Validation used the **actual final response seq18, saved turn-ended seq19 and rendered final card display ID8**, not intermediate text. GUI usage was 3283 input/397 output tokens across the three responses. Raw requests, incremental SSE responses, history, file and GUI logs remain in that owned directory.
- **Final totals:** four paid POSTs (one delegated visual + three actual GUI), two separately authorized memory-only credential reads, one provider-model GET, seven local fake requests, three short native phrases, two page-only screenshots. No paid rerun for a cosmetic answer; no secret was passed through argv, logs, files or child environment. Verified HTTPS and pre-upstream reservations were used. No further budget or credential authority is implied.
- **Cleanup and qualifications:** v5’s owned primary PIDs/listeners were verified gone. v6 recorded GUI55316 and Chrome55317 **exit0**; no separate post-exit listener check is saved for v6, so listener clearance is **not independently claimed**. Interaction was programmatic DOM activation/submission, not physical input or OS-focus proof. Chrome used the mute flag; there is no human-hearing, physical-silence or default-voice-locality guarantee. Screenshots/geometry were not treated as direct viewing. No additional work was run to refresh these observations.

The prior full code gates and same-frozen-CLI Ch7/Ch1 passes remain valid. All authorized operational jobs have completed. **Saved FINAL/HOLD for root inspection/join; no source change, rebuild, further test, provider/browser action or cleanup sweep is pending or authorized by this note.**
