# Ported feature: abandon a degenerate reply instead of answering with it

Branch `port/3-runaway-guard`, based on `fix/linux-path-case`. Coordinator, 2026-11-03. Ported from `go-checkpoint` (`12c8ccb5`, `2fd11b97`).

A degraded endpoint can stream one character or one short phrase thousands of times. The attempt is now cancelled as soon as 500 identical non-whitespace characters arrive, or 400 characters turn out to be one phrase of at most 40 characters repeated; the bytes are discarded, never recorded, and one retry is allowed before the fault is reported.

Whitespace is never evidence, so a wide aligned table survives; a uniform single-character cycle is left to the higher run bound, so a drawn rule survives; and tool arguments are exempt, so a generated file of separators is unaffected. Observed and user cancellation stay distinguishable, because the guard cancels only its own linked token.
