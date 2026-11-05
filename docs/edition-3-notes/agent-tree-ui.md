# Ported feature: Agent tree arrangement, one action set, guarded deletion, hotkeys

Branch `port/10-agent-tree-ui`, based on `fix/linux-path-case`. Coordinator, 2026-11-03. Ported from `go-checkpoint` (`f8eb24a4` + hotkeys).

- **Drag to reorder siblings only.** A row accepts a drop only from its own parent; a child dragged toward another parent is refused, because arranging is this browser's business and the hierarchy is the server's. The arrangement is stored under `ensemble.agentOrder.v1` and survives reload.
- **One action bar for the selected Agent** (`#agent-actions`) replaces the copy of the delete button on every row. It targets whichever row is selected and goes inert for the mounted root.
- **Confirmation before deletion.** The shared control reuses the existing `data-delete-agent` path, so the same `window.confirm` names the Agent; declining sends nothing.
- **Alt+Shift+Arrow walks the tree.** Plain arrows stay with the rows, `Alt+Arrow` remains Firefox history and `Ctrl+Shift+Arrow` word selection, and an editable target keeps its keys. The walk continues from the focused row so a second press advances a second row even while a selection change is still in flight.

Three defects surfaced while testing: my children lookup inverted the parent test so nested rows vanished; the shortcut focused a node captured before `selectAgent` re-rendered the tree, which silently did nothing; and the test itself reused stale node handles for the same reason.

`crossover-browser.cjs` drove the per-row trash buttons that this change replaces, so its deletion steps now select the row and use the shared control, and it asserts the root still has no trash affordance. **That suite needs the macOS shell fixture and could not be run here** — it parses, and its expectations were rewritten to the new contract, but it is unverified on this machine.
