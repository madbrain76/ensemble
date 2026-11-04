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
- `run-browser-tests.cjs` starts one disposable `GuiTests.dll --crossover-host`
  in a private temp directory, runs the suites, and reports the tier with a
  nonzero exit on failure. A suite that dials the crossover host is reported as
  skipped when that fixture cannot start, because this build's shell launcher is
  the macOS system PTY.

## Result on Linux + Firefox

12 pass, 3 fail, 19 skipped of 34. Passing: browser-theme, capture-consent,
composer-layout, image-picker, message-colors, message-playback, model-controls,
voice-browser, voice-control, voice-hold, voice-session, voice-shortcut.
Two of the three remaining failures are pre-existing Linux defects the tier had
never been able to reach: `header-models.cjs` expects a transparent computed
background where Firefox reports `rgb(38, 62, 100)`, a real cross-browser styling
question, and `page-capture.cjs` drives Chrome's tab-capture API. `voice-funding.cjs`
still resolves a build artifact the author's layout provided. None is a Firefox
regression in the product; all three are recorded rather than papered over.

## Product fix found on the way

`agent/Jobs/Job.cs` recorded only an allowlist of command failures. A
`PlatformNotSupportedException` escaped unrecorded, faulted the worker task, and
was then rethrown by teardown as "Job shutdown could not record or close every
resource", which aborted `Ensemble.Destroy` for the entire application. A command
the platform refuses is now recorded like any other tool failure.

That change is its own branch, `fix/platform-refused-job-command`, and this
branch sits directly on it because the crossover fixture cannot start on Linux
without it. It is not a Firefox change and should be reviewed, and if accepted
merged, as a separate pull request.
