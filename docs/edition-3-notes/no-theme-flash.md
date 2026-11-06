# Ported feature: paint the confirmed theme before the stylesheet

Branch `port/19-no-theme-flash`, based on `fix/linux-path-case`. Coordinator,
2026-11-03. Ported from `go-checkpoint` (`d3a0e26f`).

`index.html` declared `data-bs-theme="dark"` and every script is deferred, so a
light-profile first frame was dark and then flipped. A small non-deferred inline
script now resolves the stored preference — including `system` against the
browser's own media query — before the stylesheet is requested. Only a stored,
valid value is trusted; anything else keeps the markup default until `theme.js`
decides, because guessing wrong is its own flash.

`theme-no-flash.cjs` holds the Bootstrap stylesheet back with request interception
and reads the attribute inside that gap: a stored light profile is already `light`
while fewer than two stylesheets have applied, a stored dark profile stays dark,
`system` follows the browser, and a nonsense stored value is not trusted.

## 2026-11-05: "system" was not honoured, and the caution above was the cause

Bill reported that a light desktop came up dark. The paragraph above reasoned that
an absent or unrecognised stored value should keep the markup default, because
guessing wrong is itself a flash. That reasoning was wrong in a specific way: the
value is not unknown to the browser, it is the operating system's own preference,
and `theme.js` always resolves missing and invalid to `system`. Leaving the markup
dark therefore did not avoid a guess; it guaranteed a dark first frame for every
OS-light visitor and then always corrected it.

Three separate defaults had to move, and only the first was obvious:

- `Examples/Gui` served `Preferences.Theme` from a record default of `"dark"`
  (`agent/gui/Common/Settings.cs:11`). `app.js` passes that value to
  `browserTheme.seed`, which *persists* it into `localStorage` when the profile has
  no stored preference. So the server's untouchable default became the browser's own
  durable choice, and the page could never return to following the OS. The Go tree
  pinned the opposite: `agent/gui/settings_test.go` asserts `theme == "system"`. The
  port silently regressed it, and no C# test read that field, which is why nothing
  caught it.
- `index.html` preselected `<option value="dark" selected>`, so the control
  contradicted the resolved default for the moment before `theme.js` runs.
- The early paint returned without deciding when nothing was stored.

`io/` is not involved here. Its per-agent files are named by artifact id and were
written by an older Go build; no code in the build that wrote these saves reads that
directory, and no saved part carries a locator, so the C# has no reason to consult it.

`ThemeDefaultTests` (`agent/GuiTests/ThemeDefaultTests.cs`, mode `--theme-default-tests`)
now pins the record default, the value `Gui.ReadPreferences` serves for a fresh GUI,
that an explicit stored `dark` still wins, and that a nonsense stored value falls
back. `theme-no-flash.cjs` gained the first-visit case it never probed — nothing
stored, OS light, must not paint dark — plus the matching OS-dark case, and its
nonsense-stored expectation now follows the OS. Reinstating the returned-without-
deciding early paint turns the unrecognised-value assertion red, and hard-coding
`Theme` back to `"dark"` turns the C# checks red, so both layers are live.

## Second pass: the stored browser copy, not the default, was the real lock (2026-11-05)

Bill relaunched after the default fix and was still dark. The default was only half
of it. `theme.js` seeded the server value through `select()`, which writes
`localStorage`, so while `Preferences.Theme` was `"dark"` every profile that ever
opened the page **recorded dark as its own choice**. A stored value outranks the
operating system permanently, and `seed()` returned early once anything was stored,
so no later server-side change could reach those profiles. Fixing the default could
not help anyone who had already loaded the page -- which is everyone.

The browser owns appearance: the client never patches `theme`, and `Settings.cs`
already labels the field legacy. So `seed()` now applies the value for the frame
without persisting it, and only a real click or select is recorded. The key moves to
`ensemble.theme.v2`; a stored v1 cannot be told apart from a deliberate pick, so it
is dropped rather than trusted and removed on sight. The cost is one click for
anyone who did genuinely choose dark since the bug shipped.

Both theme tests had been *pinning* the defect: `browser-theme.cjs` and
`theme-page.cjs` asserted the seeded value lands in localStorage, so the poisoning
had green tests behind it. They now assert the opposite. `theme-no-flash.cjs` gained
Bill's exact profile -- v1 dark, OS light, nothing stored under v2 -- which must
paint light and lose v1; the ten cases pass in headless Firefox, and reverting the
early paint turns the new case red, so it is live rather than decorative.
