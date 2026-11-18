# Ported feature: the GUI browser tier runs in the operator's Firefox

Branch `port/21-firefox-harness`, based on `fix/linux-path-case`. Coordinator,
2026-11-03, at Bill's request.

## What was wrong

Every one of the 34 `agent/GuiTests/*.cjs` suites launched Chromium; 20 required
Playwright out of the frozen `edition-2/agent/gui/node_modules`, 10 hard-coded
`/Applications/Google Chrome.app/...` or `channel: 'chrome'`, and one spawned
`/Users/bill/.local/share/ensemble-dotnet10/dotnet`. On Bill's Linux box, which
has no Chrome and reviews this GUI in Firefox, the whole tier could not run, so
"verified in a browser" silently meant "verified on the author's Mac".

## Shape

- `browser-launch.cjs` is the single launch point: Playwright Firefox, headless,
  resolving the module from `PLAYWRIGHT_MODULE`, the vendored copy beside the GUI,
  or `playwright`, and `FIREFOX_BINARY` for the executable. It drops Chrome-only
  switches instead of forwarding them, and performs no download: the operator
  installs the browser once.
- `run-browser-tests.cjs` gives every suite its own disposable fixture host in a
  private temp directory, runs the tier, and exits nonzero on any failure. One
  shared host cannot work: each suite replays the mounted Agent's durable
  history, so an earlier suite's observations become a later suite's unexpected
  transcript. Suites whose real driver supplies its own fixture, arguments and
  environment name them in `PLAN`; a suite the platform genuinely cannot drive is
  recorded in `SKIP` with the reason, never counted as a pass.

## Result on Linux + Firefox

After merging upstream `ac3ac1a2`, 67 pass, 0 fail, 4 skipped of 71, on
2026-11-17. Every skip is a capability the platform does not have, not an
outcome: `page-capture-native.cjs`, `page-capture-native-ui.cjs` and
`capture-model-native.cjs` drive the macOS system PTY and launch-services
capture path, and `advanced-artifact-figures.cjs` is the human-authorized
book-figure capture that overwrites the approved PNGs in
`book/figures/advanced-artifacts`. Bill runs that one; the tier must not.

Upstream rewrote `configured-models.cjs` and `peer-choices.cjs` for its new
funding controls, and both still launched through the frozen `edition-2`
Playwright copy, so neither had ever run off macOS: each stopped at its first
assertion and everything behind that stop was unmeasured. What follows is what
they turned out to need.

Upstream added nine browser suites since the previous merge. Six arrived still
requiring Playwright out of the frozen `edition-2` tree, so on Linux they could
not start at all, and four more needed the fixture their driver supplies rather
than the crossover host. They are listed below, because a suite that never ran
is worse evidence than one that fails.

## What Firefox needed that Chrome was handed for free

Each entry is a harness accommodation for an API or a fixture the suite was
written against. None of them weakens an assertion, and none is a product change.

- **The Capture Handle is Chromium-only.** `agent/gui/wwwroot/page-capture.js`
  refuses to offer the camera unless `getDisplayMedia`,
  `navigator.mediaDevices.setCaptureHandleConfig` and
  `MediaStreamTrack.prototype.getCaptureHandle` all exist, and the last two are
  absent in Firefox. `installCaptureHandleShim` in `browser-launch.cjs` supplies
  exactly the two probed members. The capture suites already fabricate
  `getDisplayMedia` and a per-track `getCaptureHandle` of their own, so the
  identity check still runs against a handle the suite controls, and a wrong tab
  or origin still refuses. `page-capture.cjs` is a pure unit suite over an
  injected `MediaDevices`: all 10 of its scenarios now run here.
- **`voice-server-page.cjs` was given the wrong host.** It needs the separately
  funded transcription payer, the model provider and `/fixture/status`, all of
  which belong to `VoiceLeaseTests --browser-host`. Under `--crossover-host` the
  browser correctly kept its hold button disabled forever, because that fixture
  advertises no transcription funding. It is now in the runner's `PLAN`, which is
  what that file's own header prescribes for a suite whose real driver supplies
  its fixture, rather than reporting a missing argument as a Firefox failure.
- **The state word was matched against the whole element.** `speech-pause.cjs`
  and `working-context.cjs` required `#state` to read exactly `idle` or
  `paused`. The element reads `"<state> · <context>"`, and `job-status.cjs:61`
  pins that a live Job outranks the turn word, so a fixture holding a real
  300-second Job answers `"1 job running"` — which is the accurate display and
  was the assertion's blind spot, not the product's. They now assert the state
  each connection has acknowledged, and `working-context.cjs` still matches the
  working-context half as text. `header-status.cjs` and `job-status.cjs` remain
  the owners of how the word is rendered.
- **Firefox ignores `clipboardData` in `ClipboardEventInit`.** It constructs the
  event with an empty transfer, so `image-browser.cjs` now attaches its
  fabricated `DataTransfer` as an own property of the event it dispatches.
- **A hover that survives navigation.** `header-models.cjs` hovers Create Agent,
  and Firefox keeps the pointer there across the next `goto`, so the following
  iteration measured the `:hover` background as the resting one and read
  `rgb(38, 62, 100)` where Chrome reported a transparent `rgba(0, 0, 0, 0)`. The
  loop parks the pointer before each load; the resting-transparent and
  hover-painted pair are both still asserted.
- **Preferences that moved to the browser profile.** Upstream `784e841d` made
  the font size and the voice shortcut local to the browser, so
  `crossover-browser.cjs` and `voice-server-page.cjs` read
  `ensemble.device.v1` out of `localStorage` instead of waiting for an Agent
  `settings_updated` frame that no longer carries them.
- **The fake provider now answers the media probe honestly.**
  `agent/GuiTests/CrossoverFixture.cs` decodes the probe's base64 PNG, inflates
  its IDAT and reports the argmax pixel, so an unlisted model earns images in
  the tier the way it does against a real endpoint instead of being refused by a
  fixture that never looked at the bytes.
- **Contrast on transparent artwork.** `crossover-browser.cjs` measured the
  declared `color` of `#send`, which is an `<img class="send-icon">` on a
  transparent button: that declared colour never paints a pixel, so the check
  proved nothing. It now draws the artwork into an `OffscreenCanvas`, applies the
  `brightness()` factor the stylesheet actually declares, and takes the best
  contrast against the ancestor background. Both themes clear the same 4.5 floor,
  at 17.40:1.
- **Upstream's new suites still pointed at `edition-2`.** `history-window.cjs`,
  `layout-browser.cjs`, `error-cards-browser.cjs`, `recall-browser.cjs`,
  `settings-autosave-browser.cjs` and `workflow-job-sessions.cjs` required
  Playwright from the frozen tree and launched `channel: 'chrome'`. They now go
  through `launchBrowser`, which drops the Chrome-only selection; their
  assertions are untouched.
- **Suites whose driver also supplies a fixture.** `history-window.cjs`,
  `error-cards-browser.cjs`, `recall-browser.cjs` and
  `workflow-job-sessions.cjs` read their origin from a fixture directory argument
  and stop that fixture by writing beside it, so they belong in `PLAN` against
  `--artifact-renderer-fixture`, the mode that serves the shipped page from a
  private directory. The crossover host cannot stand in: these suites inject
  frames into the page rather than drive a mounted Agent, and without the
  argument they died on an undefined path, which is not a Firefox failure.
- **A sandbox that must supply what `render()` now asks for.** `render()` refuses
  a card the page does not own and asks the Job session where a result belongs, so
  the two suites that evaluate it outside the running app declare `jobSessions`,
  and `message-playback.cjs` registers the cards it fabricates.
  `error-cards.cjs` and `settings-autosave-browser.cjs` evaluate a slice of
  `app.js`, so each now supplies or includes the helpers that slice calls
  (`recordTurnTiming`, `updateAgentActions`, `siblings`). These are the
  dependencies the shipped page already provides, not softened assertions.
- **Labels that carry live state.** `capture-consent-browser.cjs` compared the
  consent dialog's requester with the selected picker option, and both end in the
  Agent's state, which the suite's own typing had just changed. It now compares
  the Agent each label names, which is what the check was about.
- **A physical size measured at fractional device pixels.**
  `layout-browser.cjs` required Create Agent to be at least 44 CSS pixels and
  Firefox reported 43.99998 for a rule that says `44px`. The assertion allows half
  a pixel at each end; an icon-sized target against a pane-wide one is still the
  distinction being tested.
- **A recorded-observations group that is not unique on the page.** Upstream moved
  a Job's recorded observations inside its session node, so
  `job-sessions.cjs` opens that group by this suite's own handle, the same rule
  its stubbed route already follows.
- **A refused write now keeps the operator's edit.** Upstream's autosave separates
  an unsent settings draft from the acknowledged policy, so a save the Agent
  refuses leaves the typed number in the field instead of rewriting it.
  `working-context.cjs` asserted the old rewrite; it now asserts both halves of
  the new behaviour — the draft is kept, and the policy the Agent reports is the
  one it was before the refusal, read back through a reload that discards drafts.
  This checks the durable policy more directly than watching a form snap back.
- **A picker rebuilt between two round trips.** `crossover-browser.cjs` added an
  `Invalid empty model` option and then selected it in a separate call, but
  `showUsage()` calls `picker.replaceChildren()` on every usage frame, so the
  injected option was sometimes gone before the selection arrived and the refusal
  it tests never happened. Option, selection and the page's own `change` event now
  go over in one evaluate, which keeps the refusal and drops the race.
- **A download button whose Job the display window chose.** Upstream keeps the
  newest 100 events per pane, so the oldest card no longer reaches the browser: the
  fixture reports `handle` 1, while the only surviving "Download complete output"
  control asked for `handle` 12 and was served that Job's 273 bytes, not Job 1's
  68,014. The check compared one Job's hash against another's bytes. It now compares
  the saved file with what the endpoint serves for the Job the page itself asked
  about, and requires that request to be a guarded `/jobs/output` for this Agent.
  The preceding check still pins `fixture.outputHash` for Job 1, so nothing about
  the guard's authority is now taken on trust.

## Product fix found on the way

A large text size in a narrow window made the Create Agent control unreachable.
Its heading row is a single flex line whose controls keep a fixed physical size
while its text scales, so at 640px and 28px the row overflowed the Agents column
and the control landed under the chat pane, visible but unable to receive a
click. `flex-wrap: wrap` on `.agents-heading` keeps the controls inside their
column, and `header-models.cjs` now measures a reachable control in all twelve
theme, width and font combinations.

`agent/Jobs/Job.cs` recorded only an allowlist of command failures. A
`PlatformNotSupportedException` escaped unrecorded, faulted the worker task, and
was then rethrown by teardown as "Job shutdown could not record or close every
resource", which aborted `Ensemble.Destroy` for the entire application. A command
the platform refuses is now recorded like any other tool failure.

That change is its own branch, `fix/platform-refused-job-command`, and this
branch sits directly on it because the crossover fixture cannot start on Linux
without it. It is not a Firefox change and should be reviewed, and if accepted
merged, as a separate pull request.

## What the funding suite found after the merge to `ac3ac1a2`

`configured-models.cjs` asserted a label upstream itself had renamed, so it
stopped before reaching any of its funding flow. Once that stopped masking the
rest, two more assumptions proved specific to the author's machine:

- It expected no exact-route options beside an unresolved default. The frame the
  same fixture serves carries two qualified choices, and
  `agent/GuiTests/ConfiguredModelTests.cs` pins exactly that count and requires
  them to stay displayable for recovery, so the picker must show both. The
  browser check now asserts two and says why.
- It expected the funding chooser to survive the removal of one of two bindings.
  With one binding left there is nothing to disambiguate, and the panel renders
  no chooser, which is the documented behaviour of the sole-binding case in
  `agent/gui/wwwroot/settings.js`. The suite now selects by option text when a
  chooser exists and proceeds when it legitimately does not.

Neither change weakens a check: the funding POST, the refusal rollback, the
unavailable reload and the peer recovery are all still exercised, and the
suite now runs its whole length here for the first time.

## Why `renderer-trials.cjs` needed a precondition wait

Run inside the full tier on a loaded machine this suite failed at
`agent/GuiTests/renderer-trials.cjs:37` with "element is not visible" after 59
retries; run on its own it passed all nine checks. The cause is a click that
precedes its own listener: `agent/gui/wwwroot/settings.js` attaches the open
handler to `#settings-open` when the module runs, so under CPU starvation the
button is clickable before anything is listening, the dialog never opens, and
the `Speech & Display` tab inside it stays hidden forever. Playwright's own
visibility wait cannot rescue that, because the precondition that failed was the
previous click.

`panel()` now re-asserts the precondition instead of assuming it: it clicks,
waits two seconds for `#settings-dialog.open`, and retries up to five times
before clicking the tab. The check is unchanged in what it proves -- the real
renderer controls behind the real settings dialog -- and the retry only removes
the ordering assumption. `renderer-trials.cjs` passes standalone and in the full
tier, which is what took the result from 66/1/4 to 67/0/4.
