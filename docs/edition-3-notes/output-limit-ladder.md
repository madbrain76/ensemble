# Ported feature: recover an output limit instead of failing the turn

Branch `port/1-output-limit-ladder`, based on `fix/linux-path-case`. Coordinator,
2026-11-03. Ported from `go-checkpoint` (`3429b42c`).

A reasoning model can spend the entire output budget deliberating and stop with no
visible answer. The engine now asks again, once at a doubled cap and then with one
appended user turn demanding the final answer, and reports an error only if both
fail. Retries are gated three ways: the provider must have returned its own length
verdict (`OutputCapped`, distinct from `IncompleteOutput`, which also fires on
filters and a missing finish reason), nothing may have been published, and the turn
must not be cancelled. A partially delivered answer is never rewritten.

The first version keyed only on "no visible text", and it broke
`Chapter 14 opt-in keep/stub actor and snapshot replay`: that fixture's scripted
empty reply was being retried until its request budget ran out. Requiring the
provider's length reason fixed it and is the honest gate — a model that chose to
answer empty should not be billed twice.
