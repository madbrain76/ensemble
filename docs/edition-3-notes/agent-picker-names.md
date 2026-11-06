# Agent picker showed wire addresses instead of names

Found while verifying the Go history import in Firefox: the header dropdown listed
`ensemble%20debugger`, `PDF%20designer`, `Yolink%20local`. Every imported peer has a
space in its name, so all six displayed as addresses.

`showAgents()` built its label from `row.key`, the percent-encoded address the wire
carries, while the tree and the composer placeholder already use `row.id`, the name:

    const label = row.key === "" ? "Original unnamed conversation" : row.key ?? row.id;

The same text is read by `captureRequested()` for the capture consent dialog, so the
percent-decoded-looking address reached a prompt the user has to answer. One line,
`row.id || row.key`, fixes both.

This is upstream code (`7938082b`, "Finished ch21."), not one of the ported features.
It survived because every test fixture used single-word names, so no assertion ever
distinguished key from id. `agent/GuiTests/agent-picker-names.cjs` now serves rows
whose keys and ids differ -- including an encoded space and a parent reference -- and
asserts the option text, the selected text the consent dialog reads, and the composer
placeholder all show the name and none contains `%20`. Reverting the line fails the
first assertion, so the test is live.

## Verification

Firefox (headless, self-contained fixture): 7/7. Merged into `local/daily` and
observed on Bill's own restored history, where the picker reads
`ensemble debugger — idle` and clicking the row loads that peer's 3624-card
transcript, distinct from the primary's 1389.
