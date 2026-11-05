# Ported feature: resume a turn that promised a tool call it never emitted

Branch `port/2-promissory-nudge`, based on `fix/linux-path-case`. Coordinator, 2026-11-03. Ported from `go-checkpoint` (`bc9f8afa`, `d995c6ef`, `00a8169d`).

Text with no tool call is the final answer as far as the engine is concerned, so a model that ends its turn announcing the work it is about to do simply stops, mid-task, with no error. The actor now recognises that shape and appends one system correction telling it to issue the call or answer fully, bounded to two per turn so a model that keeps narrating ends the turn instead of looping.

The predicate is deliberately narrow. A trailing colon alone is not a promise, because `Fix order (cheapest risk-reduction first):` followed by the list is a complete answer; a colon counts only when the same reply also names work in progress. First-person futures (`I'll`, `let me`, `I should`) match anywhere, and an action gerund counts as an opener even after connector words or a parenthetical question, which is what hid the real stall `Now inspecting the panics (are they invariant checks or control flow?), ...`.

`agent/Tests/PromissoryTests.cs` pins the seven stall texts verbatim from Bill's sessions plus seven ordinary answers that must not be nudged, including `Tail:` and a markdown summary list.
