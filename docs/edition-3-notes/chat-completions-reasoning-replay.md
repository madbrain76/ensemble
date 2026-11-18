# Chat Completions must not replay provider reasoning

Branch `fix/chat-completions-reasoning-replay`, commit `a4340b79`, from
`fix/linux-path-case`. Product code only plus one test file; nothing else moved.

## Symptom

A prompt in a restored conversation produced no reply. The endpoint answered
before generating anything:

```
POST /v1/chat/completions  422 Unprocessable Content
  body.messages[2].content.str  Input should be a valid string
```

Any thread with recorded reasoning failed on its first turn, so the conversation
looked dead rather than degraded.

## Cause

`RenderOpenAI` built the assistant `content` array from every provenance-matched
`OpaquePart` verbatim. On this surface `content` is a typed part array, so a
member without `type` fails the endpoint's schema. Two shapes reach history:

| shape | origin | endpoint verdict |
|---|---|---|
| `{"reasoning_content": "..."}` | written by the Go engine, so imported history carries it | `422` |
| `{"type":"reasoning_content","text":"..."}` | `agent/Engine/Parse.cs:170` when this codec reads the field back | `400 Unsupported content part type` |

Probe against the live endpoint, one small request per shape: string content
`200`; untyped part `422`; `type: reasoning_content` `400`; message-level
`reasoning_content` `200`; `type: thinking` `400`. Only the message-level field
is a vendor extension, and a stricter OpenAI-shaped server could reject it as an
unknown field, so nothing is invented to replace the dropped blocks.

## Fix

Drop opaque parts in `RenderOpenAI`, both from the emitted array and from the
decision to keep an entry at all. A turn holding only reasoning is then omitted
exactly as a bare turn is. This matches the reference implementation on this
surface, whose `case "opaque":` is a no-op carrying the same comment
(`go-checkpoint:agent/internal/engine/openai.go:119`), and it leaves the
delivered text alone.

`agent/Engine/Render.cs:351` (Messages) and `:721` (GenerateContent) keep
emitting their own blocks: thinking with its signature, and Gemini thought
signatures, both round trip on the surfaces that accept them. The defect was
this one renderer, not opaque replay in general.

## Test

`agent/Tests/ReasoningReplayTests.cs`, registered in `agent/Tests/Program.cs`.
Both recorded shapes are replayed into a chat-completions request and the whole
content array is checked, not just reasoning: every member must be `text` or
`image_url`. Delivered text must survive, a reasoning-only entry must be omitted,
an ordinary text turn must stay a string content, and the Messages surface must
still carry its thinking block, which is what stops a future fix from deleting
replay everywhere.

Reinstating the deleted line fails the case, and the failure prints the exact
poison body, so the assertion is live rather than decorative.

## Verification

- `agent/Tests`: 342 passed, 4 failed -- only the known `agent/Jobs/Pty.cs`
  macOS gates on this base. No change in count.
- Real Firefox against the real launcher on a copy of the live data directory:
  the imported 587-card conversation loaded, the composer
  armed, the prompt was sent, one request reached a capturing provider and the
  reply rendered with zero page errors.
- A/B against the live endpoint. `BEFORE` is the request his agent logged in
  `api.log`; `AFTER` is the request the fixed renderer produced from the same
  restored history. Both with `max_completion_tokens` pinned to 1, since a
  schema failure happens before generation.

| request | untyped parts | mentions reasoning | verdict |
|---|---|---|---|
| before (logged) | 219 | yes | `422` |
| after (fixed) | 0 | no | `200` accepted |

Dropping the replay also cut that request from 784,870 to 246,479 bytes. That is
not the point of the change, but recorded reasoning was two thirds of every
resumed prompt, so prefill gets materially cheaper as a side effect.

## Not done

Not rebuilt into the live launcher while his GUI held the data directory; the
process must be closed before the dll is replaced. `agent/GuiTests` was not
re-run because no GUI, JS, CSS or wire-format contract changed -- only an
outgoing request body. Say so if you want it run anyway.
