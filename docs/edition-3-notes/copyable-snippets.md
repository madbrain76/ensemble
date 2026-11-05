# Ported feature: copyable text boxes inside replies

Branch `port/12-copyable-snippets`, based on `fix/linux-path-case`. Coordinator,
2026-11-03. Ported from `go-checkpoint` (`cedf90d7`).

A fenced block is content the operator was meant to take away, so every complete
fence becomes a copyable box: the literal payload in a read-only field labelled
"Copyable text", and one "Copy text" control that confirms with "Copied" and keeps
keyboard focus. The payload drops only the newline the fence itself contributes and
the collapsed-preview ellipsis, so a copy is what the model wrote.

Two defects the suite found: the copy control of a long block was being placed
inside its collapsed `<details>`, hiding the only way to copy it, so the block now
anchors on the collapsible wrapper and the control stays visible; and the payload
carried a stray trailing newline.

`copyable-snippets.cjs` drives real Firefox. Firefox grants no clipboard
permission, so both paths are observed rather than trusted: the Clipboard API call
is recorded with its exact argument, and the selection fallback is asserted to have
selected the whole literal payload — including the complete text of a block whose
preview was truncated.
