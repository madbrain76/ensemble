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

Against the real endpoint: `stream: true` on the wire, `text/event-stream` back, and a
peer's turn rendering progressively. One reply did arrive in a single frame, which is
the endpoint's choice; that is what the negative observation is for, and nothing in the
client collapses chunks that really arrive.
