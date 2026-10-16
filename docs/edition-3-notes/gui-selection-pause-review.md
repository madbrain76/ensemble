# GUI selection pause review — initial finding and final recheck

## Final verdict

**Accept the corrected revisions recorded in “Final bounded recheck” below. R1 is resolved.** No remaining blocking finding was identified in this bounded correction review. The eight added Chapter 21 lines remain acceptable without edits. The initial finding, evidence, and revision boundary below are retained as history, not as an unresolved rejection of the corrected files.

## Initial verdict (superseded by final recheck)

**Changes requested: R1 blocked acceptance.** The ordinary single-selection path and its new regression test were sound, but rapid repeated selections could release the pending guard too early. The initial review covered the revisions in the first table below; acceptance applies only to the corrected hashes in the final recheck.

## Scope and verified revision

Repository: `/Users/bill/projects/ensemble`.

Verified HEAD: `930dc42a65ea81560577388cd73140bbb017e47c`. This differs from the supplied transcription. The review concerns working-tree changes over that HEAD.

SHA-256 of reviewed files:

| File | SHA-256 |
| --- | --- |
| `agent/gui/web/gui/app.js` | `4c2e15b1d76b1b45aadf1678ef69b585d4dbeabff2cc4b79231adf5582bca036` |
| `agent/gui/selection-pause.test.mjs` (latest teardown revision independently run) | `bb5b0fc61f915b64481402a8c931ce30c5b12c95e3db22a8728935a1433925d2` |
| `book/edition-3/chapter-21.md` | `2e0692d73d3f6379c8829239f575a7385c6123a37721f87b82da424956136e00` |

Read all required current prerequisites: `AGENTS.md`, the ensemble-coding skill, workflow, Edition 2 carryover, voice, chapter-writing procedure, Edition 3 execution procedure, both applicable Edition 3 reviewer roles, and the whole Chapter 21. Inspected the full app, new test, and surrounding fleet/speech/selection tests and server selection/pause routing. Review scope remains the GUI-only fix, its test, and the eight added chapter lines. No source, reference, grader, live settings, or runtime-observation files were changed by this reviewer. This note is the only repository file written.

## R1 — repeated target acknowledgments prematurely clear the pending guard

Locations: `agent/gui/web/gui/app.js`, `selectAgent` (around lines 138–147), and `agent_selected` handling (around lines 370–376).

`selectAgent` permits another selection while one is pending. The guard contains only the requested key. Two rapid requests for the same key therefore cannot be distinguished: the first `agent_selected` with that key clears the guard even though the second selection remains outstanding.

Independent reproduction used an in-memory variant of the new test with its real isolated GUI process and WebSocket proxy; no repository test or source file was edited:

1. Hold real selection acknowledgments at the proxy.
2. Click `peer`, wait for its first acknowledgment to be held, then click `peer` again.
3. Wait until both real `agent_selected(peer)` acknowledgments are held.
4. Deliver only the first acknowledgment.
5. Enqueue speech before delivering the second acknowledgment, using the same application-owned speech object and controlled native boundary used by the proposed test.

Captured outgoing command types:

```text
["unpause", "select_agent", "unpause", "select_agent", "pause"]
```

The assertion that no `pause` is emitted while the second selection remains pending failed. This is a real command through the real connection, not a mock of the selection or pause implementation.

Consequences established by code inspection: the server processes each connection's commands in order and routes pause commands to that connection's currently selected Agent. The second `unpause` can therefore already affect `peer`. Speech allowed by the premature guard release can send a new `pause`; the later acknowledgment clears `localPaused` and cancels speech silently, so it does not release that new pause. This breaks the incoming-pause boundary and can leave a tool gate paused without a remaining local cause.

Smallest correction direction: allow only one pending selection at a time (including repeated clicks), restore selection availability at its completion/error boundary, and add the repeated-click case under the existing held-ack fixture. No server protocol, global resume, shared speech settings, or new framework is warranted. The parent accepted R1 and will make the bounded correction.

A follow-up probe attempted to inspect nonexistent public `Speech` state properties and timed out; it supplies **no additional evidence** and is not used for this finding. The first probe's captured wire sequence and failed assertion are the valid reproduction. The follow-up fixture exited; no live GUI was signaled. A review commentary statement that the source had changed was mistaken: rehashing confirmed the same app revision. The actual variable is `pendingSelection`.

## Evidence and invariant checks

- Independently ran `node --test selection-pause.test.mjs` from the separate `agent/gui` module. Latest hashed test passed: 1 test, 0 failures, approximately 1.90 seconds total. Earlier teardown revisions also passed independently; they are not substituted for the latest result.
- The test starts an isolated real GUI with a provider fake and uses real paused `write_file` dispatch. It verifies the file does not exist while paused, then appears without another prompt after selection releases the outgoing gate.
- A second real connection keeps the incoming Agent paused. Both speech-caused and typing-caused outgoing pauses are exercised. The test checks the same-connection `unpause` then `select_agent` sequence, delayed acknowledgment, suppression of queued outgoing speech, and a stale native completion callback.
- The single-selection implementation resets local causes before cancellation, explicitly sends release before selection on the same socket, and suppresses pacing/native enqueue while pending. `Speech.cancel()` invalidates the active utterance before native cancellation; the existing callback identity check rejects a stale callback. Acknowledgment cleanup is deliberately silent.
- Static inspection: selection errors clear the pending guard; reconnect clears it and sends `subscribe`, not a global release. Those paths and interleaved unrelated errors were not independently browser-driven in this pass. A key-only guard is insufficient for rapid requests, as R1 demonstrates.
- `git diff --check` passed.
- The parent is running the original-assets negative control and broader fleet/speech/GUI race checks. This report does not claim their results as independently observed.

The native synthesizer is a controlled fake. These results establish deterministic client ordering and real server/tool-gate behavior, **not live native-speech proof**. The latest test uses `SIGKILL` only to dispose of its own temporary GUI fixture. That cleanup does not establish graceful shutdown behavior and does not expand this review into runtime shutdown.

## Text assessment and KISS

The paragraph at the replay/selection contract is load-bearing. It states the exact ordering, the otherwise stranded pending-tool consequence, the late-callback/queued-output boundary, and an observable test for outgoing progress while preserving the incoming pause. Removing it would remove a concrete selection requirement rather than merely shorten commentary.

Chapter 20 introduces headless subagents; Chapter 21 introduces the GUI fleet and connection-local selection. Placing the requirement here avoids a forward dependency and stays within Bill's explicit authorization for the GUI fix and a quick explanatory requirement. It does not add shared speech preferences, a server pause redesign, a global resume mechanism, or another independent chapter feature. No prose revision requested.

The patch is otherwise small and reuses existing speech cancellation, pacing, and selection paths. Resolve R1 with bounded client-side request handling and a regression assertion; do not redesign the pause model. Acceptance remains withheld pending that correction and a recheck of its actual hashes.

## Handoff revision boundary

At the final post-write hash check, the parent's correction had arrived: app SHA-256 `57f24a5255ef9835a353e8fd4c47e35d96d121191bcf5b1148338fa576f5226f`, test SHA-256 `10f8d0de41206494398c25e674f30100f19e36f7e7643dc3dfd1896bec02fa60`; the chapter hash was unchanged. Those are observed current worktree revisions, **not reviewed revisions in this initial pass**. The finding and passing-test evidence above remain tied to the earlier table. `git diff --check` remained clean. A bounded recheck is required before accepting the correction.

## Final bounded recheck — R1 resolved, accepted

The parent froze the correction and requested only the R1 recheck. No further editorial pass or source changes were made by this reviewer.

Verified HEAD remains `930dc42a65ea81560577388cd73140bbb017e47c`.

Final reviewed SHA-256 values, independently verified before and after the test:

| File | SHA-256 |
| --- | --- |
| `agent/gui/web/gui/app.js` | `57f24a5255ef9835a353e8fd4c47e35d96d121191bcf5b1148338fa576f5226f` |
| `agent/gui/selection-pause.test.mjs` | `10f8d0de41206494398c25e674f30100f19e36f7e7643dc3dfd1896bec02fa60` |
| `book/edition-3/chapter-21.md` (unchanged) | `2e0692d73d3f6379c8829239f575a7385c6123a37721f87b82da424956136e00` |

Correction inspection:

- `selectAgent` now returns immediately if a selection is already pending. It sets the guard synchronously before cancellation or either socket command, so a queued duplicate handler cannot emit another release/selection pair.
- The tree is immediately rendered with selection buttons disabled. This makes the temporary restriction visible rather than silently accepting a second click. The handler guard remains necessary and is present.
- The matching acknowledgment clears the guard. Existing server `selectConversation` sends `agent_selected`, then `agent_tree`, in one ordered replay batch; that tree render restores selection buttons. The error branch clears the guard and explicitly rerenders. Reconnect still clears pending state and uses the existing subscribe/replay path rather than sending a global resume.
- The revised regression holds the real acknowledgment, checks that the peer button is disabled, explicitly dispatches another click, and verifies that the wire sequence remains exactly `unpause`, `select_agent`. Both speech and typing iterations still test the real outgoing `write_file` gate, preservation of the independently paused peer, queued speech suppression, and old native callbacks.
- This is a small client-side serialization fix using the existing state and render path. It closes the reproduced R1 sequence without protocol changes, new shared settings, server pause redesign, or an additional framework.

Independent verification:

```text
cd agent/gui
node --test selection-pause.test.mjs
# 1 test, 1 pass, 0 failures
# test duration 1986.667459 ms; total 2147.562584 ms
```

`git diff --check` also passed. The clean result and guard inspection resolve R1. The new test does not independently drive every error/reconnect interleaving; the initial coverage limits remain explicit. Controlled native speech remains deterministic boundary evidence, not live speech proof.

Parent-reported complementary evidence (not independently rerun in this recheck): final new test passed in approximately 2.08 seconds; clean original-JS negative control failed in approximately 2.50 seconds at actual `["select_agent"]` versus required `["unpause", "select_agent"]`, without timeout; existing fleet/speech tests and the full GUI Go race suite repassed the latest app. The earlier negative-control teardown timeout remains a failed attempt, not sensitivity evidence. Temporary-fixture force-kill cleanup remains separate from graceful runtime-shutdown claims.

**Final disposition:** accept these corrected hashes for the authorized GUI-only fix and quick Chapter 21 requirement. R1 is closed; no remaining blocking finding was identified in this bounded recheck. The unchanged paragraph remains load-bearing and within authorization, and the correction remains KISS. No installation, live-service restart, live-browser inspection, settings changes, or unrelated-file edits were performed.
