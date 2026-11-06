# The promissory bound had a hole that made it unbounded

Follow-up on `port/2-promissory-nudge`, commit `6295a7f4`.

## What was observed

Through the real launcher, with the imported Go conversation and a provider that
replies with promising text and no tool call:

| | provider requests | median gap |
|---|---|---|
| before | 1,372 in 20 s (second run: 793 in 10 s) | 11 ms |
| after | 5, all within 0.1 s | then silence |

`events.jsonl` for that run: 240 `message_received` all from the **system** actor
(the correction text), 413 `request_sent`, 174 `turn_ended`. Instrumenting every
`requests.Enqueue` site showed only the nudge site firing, and the queue depth
climbing monotonically -- `q=92, 93, ... 99` -- with exactly two corrections at
each depth. The queue was filling faster than it drained. Nothing ever executed,
which is what the operator experiences as an agent that idles.

## Why

```csharp
if (promissoryReplies < PromissoryLimit && PromisesAction(LastAgentAnswer()))
{
    promissoryReplies++;
    Record(...PromissoryCorrection...);
    requests.Enqueue(new Mail(MailKind.Submit));   // one more request
    return;
}
promissoryReplies = 0;                             // re-arms on give-up
Finish();
```

The give-up path cleared the counter it had just exhausted. Since each correction
enqueues another request, and `LastAgentAnswer()` is unchanged by a refusal, every
queued turn received a fresh allowance of two. Self-amplifying: the bound was a
per-turn allowance on a queue that manufactures turns.

The constant's own comment already said the right thing -- "an unbounded nudge
would loop forever" -- so this was a placement error, not a disagreement about
intent.

## Fix

Clear the count where the model **acts** (the branch with tool calls), and where a
**human** starts a turn. Do not clear it when the bound is refused: that refusal is
the statement that this reply is not going to act. Releasing on human input is what
keeps the guard honest for the next task rather than muting the agent for the rest
of the session.

## Test, and what it does not do

`PromissoryLoopTests` drives one exchange against a provider that always promises,
asserts the request count stays within `1 + PromissoryLimit`, asserts a second human
turn still gets corrected, and turns a hang into a failure instead of wedging the
suite.

It must be read honestly: with the defect reinstated, the test still passes.
Neither a bare exchange nor one with a hint arriving mid-turn reproduces the
runaway, which needs the GUI's queue behaviour. So this test guards the property
and would not have found the bug. The measurement above is the evidence, taken on
the product path with the real launcher and a restored conversation.

## Still open, deliberately untouched

Killing the provider produced a second hot loop: 51 consecutive
`request_sent -> error_occurred -> turn_ended "HTTP exchange failed
(HttpRequestException)"` at the same ~11 ms cadence, with no backoff and no bound.
Transport retry policy is a design decision (how many attempts, what delay, whether
an operator can cancel), not a line placement, so it is left for a decision rather
than fixed by guesswork. It matters for the real setup: a provider that restarts, as
the GPU server does, would be hammered at 90 requests a second by an agent that
looks like it is merely waiting.
