# Chat Completions must budget room to think and then act

Branch `fix/chat-completions-output-budget`, commit `dbf98faf`, from
`fix/linux-path-case`. Upstream defect: `agent/Examples/Gui` on `main` sends
`max_completion_tokens: 1024` (`agent/Engine/Render.cs:601` at `ba12137a`).

## Symptom

The agent narrated a plan and then stopped. No tool ran, so the operator's
summary row read `Tools (0)` and the question was "why are you idling?". The
agent's own account was accurate: "three calls failed before dispatch".

## Cause

Deliberation and the answer share one output budget. With `reasoning_effort` at a
high level this model writes 1,500 to 6,800 characters of reasoning before it
emits anything, which is more than 1024 tokens. The provider said so plainly --
`finish_reason: "length"` with no usable tool call -- and the engine treated the
truncated narration as a completed turn.

The reference implementation on this surface uses 4096
(`go-checkpoint:agent/internal/engine/openai.go:186`, `maxTokens := 4096`, omitted
entirely for compaction). The messages dialect's 1024 was carried across correctly
and stays: that value was never measured against this failure, so widening it
without evidence would be an unmeasured change.

## Measured on the live endpoint

The request the agent had actually logged, replayed unchanged except for the
budget, `max_completion_tokens` pinned:

| cap | finish_reason | reasoning chars | tool calls |
|---|---|---|---|
| 1024 | `length` | 1,558 | 1 -- truncated mid-arguments |
| **4096** | `tool_calls` | 6,811 | **2** |
| 8192 | `tool_calls` | 2,478 | 2 |

1024 did not merely cut the answer short; it produced a **partially formed tool
call**, which is the worst case, since the arguments are JSON that never ended.
The companion gap is handled separately by `port/5-truncated-tool-call`, which
refuses to dispatch a call that the budget truncated.

## Product path

Real launcher, copy of the live data directory, real Firefox, prompt typed into
the composer, scripted provider that only offers a call when the budget is large
enough. Requests carried 4096, the tool executed, and the model received:

```
E2E-CAP-OK
Linux

exit_code: 0
```

Zero page errors. The counterfactual with the constant restored to 1024 and the
same provider: no tool call, no reply.

## Test

`OutputBudgetTests` pins 4096 for this surface, pins 1024 for the messages
dialect so the widening cannot spread silently, and keeps the base low enough that
the doubling ladder still has steps available. Reinstating 1024 fails it with
"chat completions must budget room to think and then act, got 1024".

## Related, not this branch

A separate runaway was found while measuring this: the promissory correction
re-armed its own budget and issued 1,372 requests in 20 seconds. That is fixed on
`port/2-promissory-nudge` and documented in `promissory-loop-bound.md`, which also
records a transport-retry loop left open on purpose.
