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
