# GUI selection releases the outgoing pause

Bill authorized a GUI-only fix and a short book requirement after Chapter 1 completion. Base revision: `930dc42a65ea81560577388cd73140bbb017e47c`. Chapter 2 was not launched.

## Defect and bounded change

Previously, clicking another Agent sent `select_agent` before canceling outgoing speech. The server changed the connection's selected Agent before returning `agent_selected`; that acknowledgment then canceled speech, whose unaddressed `unpause` affected the incoming Agent. The outgoing Agent could remain behind a tool gate even after its model response was complete. Typing causes also needed to be cleared at this view boundary.

The browser now cancels the outgoing speech and clears local speech/typing causes, explicitly sends `unpause`, then sends `select_agent` over the same WebSocket. Native speech and automatic pacing are suppressed while selection is pending; transcript recording is retained. Acknowledgment cleanup is silent, so it cannot release the incoming Agent's pause. Existing active-utterance identity checks reject late native callbacks.

Independent review demonstrated that two rapid selections of the same target could produce two acknowledgments and release a key-only guard too early. The bounded correction permits one outstanding selection, disables selection buttons during that interval, rejects an already queued click, and restores the controls at completion/error. No server protocol or pause-gate change was made.

Final source identities:

- `agent/gui/web/gui/app.js`: SHA256 `57f24a5255ef9835a353e8fd4c47e35d96d121191bcf5b1148338fa576f5226f`.
- `agent/gui/selection-pause.test.mjs`: SHA256 `10f8d0de41206494398c25e674f30100f19e36f7e7643dc3dfd1896bec02fa60`.
- `book/edition-3/chapter-21.md`: SHA256 `2e0692d73d3f6379c8829239f575a7385c6123a37721f87b82da424956136e00`.

## Executed evidence

From the separate `agent/gui` module:

- `node --test selection-pause.test.mjs`: [passed](gui-selection-pause-evidence/fixed-positive.txt), 1 test, 0 failures, about2.08seconds. An isolated real GUI/server and fake provider produce an actual `write_file` call; the file is absent behind the outgoing pause and appears after switching, without another prompt. A second real connection maintains the incoming Agent's pause. Both speech and typing cases are tested, together with ordered commands, delayed acknowledgment, an already queued repeated click, queued outgoing speech, and late native callbacks.
- The same final test with original HEAD JavaScript supplied through `GUI_ASSETS`: [failed](gui-selection-pause-evidence/original-negative.txt) in about2.50seconds at `['select_agent']` versus required `['unpause','select_agent']`. This is the clean sensitivity control, not a timeout.
- `GUI_BINARY=/tmp/cadence-selection-pause-evidence/ensemble-gui node fleet.test.mjs`: [passed](gui-selection-pause-evidence/fleet.txt) against a fresh `go build -o ... ./cmd` binary.
- `node --test speech.test.mjs`: [passed](gui-selection-pause-evidence/speech.txt), 2 reported tests, including the existing speech/pacing and speed/reload checks.
- `go test -race ./...`: [passed](gui-selection-pause-evidence/gui-go-race.txt) for the GUI module and its command packages.
- Both changed JavaScript files passed `node --check`; tracked diff whitespace checks passed. The short paragraph adds75 whitespace-separated words; fenced-example count is unchanged.

These are isolated deterministic browser tests. A controlled native synthesizer is not live speech proof. Provider fakes do not establish live-provider behavior. No test browser attached to the live GUI/MCP connection.

## Failed attempts and limits

Initial fixture attempts used the wrong build target/environment convention and an incorrect pause-label assertion; neither reached the intended acceptance check. The fixture was corrected to use actual `./cmd`, `LLM_*`, existing skill flags, installed Chrome and the UI's literal `paused` state.

The first original-assets negative attempt [timed out during teardown](gui-selection-pause-evidence/initial-negative-teardown-timeout.txt). Stdin EOF, followed by SIGTERM to the verified temporary GUI process, did not complete shutdown while the original bug stranded the tool. That attempt is **not** the clean sensitivity evidence. The known isolated process was disposed after inspection; the test now deterministically disposes only its own temporary GUI with SIGKILL in `finally` and awaits exit. It does not impose a generation deadline, signal a live process, or prove graceful shutdown. Server-shutdown behavior remains outside this patch.

The [independent review](gui-selection-pause-review.md) preserves its original repeated-selection finding, correction verification and evidence limits. It also identifies an invalid exploratory probe separately rather than relying on it.

## Teaching and preservation boundary

Chapter20 introduces headless subagents and explicitly defers multi-agent GUI selection. Chapter21 introduces the tree, connection-local selection and replay boundary; its eight-line requirement therefore belongs there. It requires release-before-switch and preservation of incoming pause under late callbacks/queued output, with speech and typing acceptance cases. It does not add a new architecture or shared-preference/global-resume requirement.

This is Bill-authorized teaching from a demonstrated runtime defect, not a speculative author loop or a claimed Chapter21 student pass. The existing Chapter21 entry gate must enforce it when that chapter is executed. The accepted Chapter1 source/tag and text, historical book, copied final reference, and graders remain unchanged.

No shared speech-preference change, global-resume change, application installation, live restart, browser reload, publication rebuild or push is included. The isolated test binary is not a staged production replacement. Source acceptance and observed live deployment are separate facts.

## Final acceptance

Independent review accepted the final three source hashes above, with R1 closed and no remaining blocking finding. The unchanged paragraph was accepted in its single text pass. Review record SHA256: `55f19907e40c575ce773acce1d47451c17a14c5d953fa0b0224bfc6e6de9c069`. The reviewer was joined; its final saved `turn_ended` event197 was verified before closure. Installation/live behavior remain unclaimed.

Checkpoint whitespace qualification: the full staged check reports two space-only blank lines (17 and22) in the verbatim `original-negative.txt` assertion output. They were retained to preserve that evidence. The production source, test, chapter and both prose records pass the separate whitespace check.
