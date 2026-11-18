# Chat Completions never asked the endpoint to stream (November 8)

Bill: "streaming also doesn't seem to work properly in the GUI". He was right, and
the cause was not the page. `request["stream"] = true` existed in exactly one place
(`agent/Engine/Render.cs:145`), guarded by `config.Billing == BillingMode.Codex`.
Ordinary Chat Completions traffic therefore carried no `stream` field at all, so a
conformant endpoint returned one buffered body: nothing appeared until the answer was
finished. A logging proxy in front of his real endpoint settled it — response
`content-type: application/json`, one chunk, and zero `part_delta` frames reaching the
page.

The reference names the cost of this in its own comment
(`go-checkpoint:agent/internal/engine/media.go:175`): asking an unlisted model to
stream, when it can, is what avoids "time to first readable text" and is what gives
"artifacts ... a measurable span". That last clause is why this surfaced next to the
timestamp bug: a buffered answer has no fragment arrivals to time.

## The gate is capability, not preference

`Configuration.Streams` is the operator's per-channel *display* selection; it is not
evidence that an endpoint streams, so using it as the request gate was wrong and the
Chapter 6 test said so ("Unknown model assumed streaming support"). `ModelFeatures.Stream`
already existed as the declared permission, unused for this dialect.

`Engine.MayStream` (`agent/Engine/StreamCapability.cs`) now decides in this order:
`DisableStreaming` wins outright; a curated row answers in both directions and is never
overruled by one exchange; an unlisted model is asked once and `Engine.RecordStream`
remembers what the endpoint did with the request, so a route observed buffering is not
asked again for ten minutes and a route observed streaming keeps being asked. Recording
happens in `SendCoreAsync` from the response media type, only when we actually asked,
so an auxiliary call that never requested a stream cannot blacklist a route.

This changes prose the book pinned ("Unknown models use plain delivery",
`book/chapter-06.md`). Bill's standing decision for unlisted models is that the
provider answers whether it accepts a thing, which is how images and reasoning effort
already work, so the chapter now states the rule that runs. `internal/grade` holds no
check on the old wording.

## Verification

`StreamingActorTests` covers the ladder: unlisted asked, observed-buffered not asked
again, observed-streaming asked again, and `DisableStreaming` outranking an earned
capability. `ReasoningReplayTests` pins `stream: true` plus
`stream_options.include_usage` on the ordinary request. Forcing the gate off turns both
red, separately. A caller-supplied body stays untouched on the wire: the flag is read by
substring, because `JsonDocument.Parse` there threw on the deliberately malformed input
that test requires be forwarded verbatim, and the timeout it produced was the bug, not
the test's.

`agent/GuiTests/streaming-render.cjs` is new: it drives the shipped frame handler the
way an endpoint does and requires the page be observable mid-answer — five distinct
growing texts (`"the"`, `"the page"`, ...), appended rather than rewritten, and a
thinking card that grows in its own box. Before this fix a buffered reply could only
ever produce one state. `agent/Tests` 368 passed / 2 failed, the same two macOS Seatbelt
confined cases; GuiTests clean; ten browser suites green.

Against the real endpoint: `stream: true` on the wire and `text/event-stream` back.
One reply did arrive in a single frame, which is the endpoint's choice; that is what the
negative observation is for, and nothing in the client collapses chunks that really
arrive.

## CORRECTION, same evening: enabling streaming broke chat, and my verification missed it

The paragraph above checked that a turn **rendered progressively** and never that it
**completed**. That is precisely the hole this endpoint fell into. The operator's inference box runs a freetoken
fork that appends its own informational frames, and `AccumulateOpenAI` required
`choices` on every frame, so `{"freetoken_completion_tokens_delta":1}` aborted the
exchange: a minute of prefill, then `Invalid provider stream`. One request after the
change, one failure -- effectively all chat, not intermittent.

The defect was not the streaming gate. Buffered answers go through `ParseResponse`,
which reads the keys it wants and ignores the rest, so these frames were always
harmless until streaming existed for an unlisted model. Two of the three dialects
already tolerated the equivalent bookkeeping -- Anthropic returns on unknown frame
types ("ping and other non-content protocol bookkeeping") and Gemini returns when
`candidates` is absent or empty -- so OpenAI's strictness was an inconsistency inside
one file, not a policy.

Fixed in `AccumulateOpenAI`: an absent or explicitly null `choices` is bookkeeping and
is ignored, mirroring both siblings. A member of the **wrong type** still fails, so a
stream that promises content and withholds it remains an error rather than becoming
silently empty. `usage` merging happens before the guard, so billing survives.

The refusal also said only `FormatException`, discarding `Missing choices` -- the two
words that named the fix. `agent/Engine/Streaming.cs` now reports a `FormatException`'s
message, which is always a fixed literal this file writes, while other exception types
keep only their type name because their text is runtime-supplied. Same rule as the
launch diagnostics: our own strings are safe to show, arbitrary exception text is not.

Coverage that did not exist before: a streamed answer through a telemetry frame and a
null `choices`; a refusal that names its rule; and a wrong-typed member that must stay
an error. Installed the old strict behaviour and the Chapter 6 case fails with
`Sequence contains no elements` -- 367 passed, 3 failed -- so the assertions are live
rather than decorative. Restoring the fix returns 368/2, the two macOS Seatbelt
confined cases.

**The lesson this note should carry:** verifying that output *appears* is not verifying
that a turn *finishes*. A streamed reply renders right up to the frame that destroys
it, so a progressive-rendering assertion can be fully green while the answer is
discarded. Any streaming claim has to assert the settled result.

## Two follow-ups, both surfaced by the Firefox tier

The recall judge builds its own request body, and `AddStreamingOptions` gated on
`config.Features.Stream` — the curated row — instead of `MayStream`. An unlisted
model therefore never streamed for recall, and because the table names no local
model, that was every model Bill actually runs. The documented capability rule
now governs both request builders, and the duplicate `MayStream` branch that had
grown in `Render` is gone, so there is one answer to the question.

`StreamKey` also had to grow the offering. Two bindings of one vendor over one
base URL are distinct routes, and what one of them did with the flag says nothing
about the other; sharing a key let an endpoint observed buffering poison its
sibling's negative record. The negative cache itself is a tested contract
(`StreamingActorTests` refuses to ask a route twice within its window), not a
defect — only its identity had been too coarse.
