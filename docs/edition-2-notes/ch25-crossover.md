# Edition 2 Chapter 25: Codex-to-Ensemble crossover

Bill assigns Chapter 25 to the live crossover from Codex to Ensemble, Edition 2.
This is his first look at this GUI. Capture fixes discovered during actual use;
this note is not a chapter draft or an acceptance claim. Chapter 16 remains paused.
All implementation changes here belong in `solutions/edition-2/main/`; preserve
frozen chapter exports and Edition 1.

## First fix: approximate 5x speech speed

Request: raise the slider beyond 2x, accounting for Eloquence's nonlinear native
rate. Bill reports rate 2.0 is about 2x, 2.6 about 3.5x, and 2.9 about 5x; rate
3.0 is already about 6x. These are user observations, not new timing measurements.

Changes in `main/gui`:
- Slider and saved-settings clamp now allow 0.5–5; label says approximate speed.
- Rates through 2 are unchanged. Above 2, linearly interpolate the reported
  (displayed multiplier, native rate) anchors (2,2), (3.5,2.6), (5,2.9).
- Calibration applies to the default browser voice; it is tailored to this
  Eloquence setup, not a universal speech-engine timing guarantee.
- Real settings-wire/disk tests now cover a high fractional setting (4.2),
  multi-tab broadcast, partial patch preservation, restart, and a 5x clamp.
- Browser speech test includes a named subtest exercising the actual slider,
  settings broadcast, reload, and native utterance rates at seven settings.
  Only the native speech boundary is recorded; no audible timing is inferred.

Validation:
- `go test -race -run '^TestSettingsWireAndDisk$' -count=3 .`: PASS.
- Browser subtest `5x slider, Eloquence calibration and reload`: PASS.
- Deliberately replacing calibration with direct native rate makes that subtest
  fail at 3x (expected 2.4, observed 3); correct implementation restored.
- Full GUI Go suite fails `TestSavedBudgetEntry: application did not start`;
  also observed on the run started before source edits. Not repaired here.
- Full browser speech suite fails later at `whole final is spoken`, recording
  four identical artifacts instead of one. Reproduced with the original HEAD
  speech test against the current implementation; not caused by inserting the
  new speed subtest. Its underlying cause remains uninvestigated.
- Independent focused review (`review-speed`): no actionable findings. Reviewer
  did not rerun main actor's suites; explicitly limited timing claims to the
  recorded native-rate boundary.

Deployment: do not interrupt the active crossover session. Stage a replacement
GUI as `/Users/bill/ensemble-edition2-codex/bin/ensemble-gui.next` for a deliberate
restart; the currently running process and browser still use the old version.

## Live GUI fix: browser observations become on-demand tools

The active checkout now uses `agent/` (the earlier `solutions/edition-2/main/`
paths above describe the earlier checkpoint). This entry records a live runtime
change, not a manuscript draft, grader result, or chapter acceptance claim.

Bill noticed hearing "GUI snapshot" and "tts queue" repeatedly while Astra
worked. Inspection showed that `gui-debug` registered both browser tools with
`ephemeral: "round"`. Before every model request the runtime invoked them,
recorded tool events, and attached their latest output as ephemeral context.
Those ordinary tool events also reached the GUI's speech announcement path.
Astra saw one current snapshot and one current queue payload, not a growing
stack of snapshots: replacement of ephemeral context was working as intended.
The unwanted behavior was automatic collection and its audible/Actions traffic.

Bill chose on-demand inspection rather than merely silencing automatic calls:
- `agent/gui/web/gui/mcp.js`: remove the two ephemeral registration fields;
  preserve the `gui_snapshot` and `tts_queue` implementations and schemas.
- `agent/skills/gui-debug/SKILL.md`: instruct explicit calls for inspection and
  fresh snapshots after actions, without claiming current state from old output.
- `agent/gui/README.md`: document on-demand behavior.
- Leave the shared MCP/ephemeral mechanism unchanged for other tools.

Regression coverage uses the shipped embedded browser module in a real headless
browser (`components.test.mjs`, served by the external-consumer reuse test).
It checks all four browser tools remain discoverable, none requests ephemeral
collection, snapshot calls see changed visible text, and one explicit queue
call reads and returns the queue once. Discovery itself does not read the queue.
Initial `go test -race -run '^TestExternalReuse$' -count=1 -v .` in `agent/gui`
passed. Mutation and wider validation results will be added below.

Deployment status at this entry: source and skill updated; no GUI binary built,
installed, or restarted for this change yet. The running browser/server still
advertise the old automatic registrations until the updated assets are loaded
and the MCP connection is recreated. GUI visual polish remains unfinished and
is not part of this observation-tool fix.

Validation/deployment update for the on-demand fix:
- Deliberately restoring `ephemeral: "round"` on `gui_snapshot` makes the real
  embedded browser test fail with `gui_snapshot must be on-demand`. Restored
  the correct source in a `finally` block; the external reuse/race test passed
  again afterwards (including public events, reusable assets and relay).
- Core `go test -race -run '^TestMCP' -count=1 ./ensemble`: PASS; shared ephemeral
  collection and normal MCP calls remain covered by the existing tests.
- Full GUI `go test -race -count=1 ./...`: NOT GREEN. It failed only the reported
  `TestSavedBudgetEntry: application did not start` check in this run, matching
  the earlier recorded startup failure. Its cause was not investigated here.
- Tests above unset inherited provider selection/billing/reasoning environment
  variables to avoid treating live Codex settings as fake-provider credentials.
- Built `agent/gui/cmd` successfully and installed the replacement at the actual
  launcher target `/Users/bill/ensemble-edition2-codex/bin/ensemble-gui`.
  Prior executable retained as `bin/ensemble-gui.before-on-demand` in that same
  external runtime directory. Nothing is left merely staged as `.next`.
- Did not restart the running server or reload the human's tab. User must restart
  Ensemble with the usual launcher, then reload the browser (hard refresh if
  needed). Live confirmation that auto-attachments cease is still pending.

Live verification after deployment: restarting the server alone left the open
browser tab advertising the old ephemeral registrations. HTTP inspection showed
that the restarted server's `/mcp.js` matched the updated source exactly. After
Bill hard-refreshed the tab, the next user message and subsequent tool round no
longer included automatic snapshot/queue attachments. Both tools became directly
callable; explicit `gui_snapshot` returned the current three-pane GUI (with its
normal large-history truncation), and explicit `tts_queue` returned pending
speech. Explicit calls still generate ordinary tool-name speech announcements,
as intended. No further build or server restart was required after refreshing.

## Live configuration: request detailed reasoning summaries

Bill reported seeing only thinking titles. Inspection of the current runtime's
`session/api.log` showed `reasoning: {effort: "high", summary: "auto"}` in the
main Codex requests. Readable reasoning summaries returned by the endpoint were
short, single-line bold headings (roughly 19–32 characters in inspected parts),
or absent. The GUI renders returned readable summary text inside its Thinking
expander; no detailed paragraphs were present in those inspected responses.
Encrypted reasoning alongside summaries is opaque replay state, not display text.

Changed the external launcher `/Users/bill/ensemble-edition2-codex/start` to export
`OPENAI_REASONING_SUMMARY=detailed`, keeping `OPENAI_REASONING_EFFORT=high`.
Existing ConfigFromEnv and Responses request rendering already support this;
no production Go change or binary rebuild is necessary. Shell syntax check passed.

Added `TestCodexDetailedSummaryFromEnvironment` in
`agent/ensemble/codex_test.go`: a fake Codex endpoint asserts that real request
rendering carries both environment-selected reasoning settings. The initial test
setup had an invalid Config field; corrected it before validation. The test passed
three times under the race detector. The broader `^Test(Codex|Responses)` tests
also passed under the race detector with inherited provider selection and
reasoning/billing environment variables cleared.

Live verification is pending: the running process retains its previous environment.
Restart with the usual launcher, then inspect a new request in api.log for
`summary: "detailed"` and separately check returned readable summary lengths/content.
This change requests provider-designated detailed summaries; it neither guarantees
provider verbosity nor exposes encrypted/private reasoning.
