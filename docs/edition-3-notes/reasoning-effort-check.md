# Ported feature: ask the provider before honouring a reasoning effort

Ported from `go-checkpoint` (`928d036e`) on 2026-11-03 by the coordinator, at
Bill's request. Branch `port/9-reasoning-effort-check`, based on
`fix/linux-path-case`; independent of the other ports.

## Why

A requested effort that the endpoint will not honour is invisible: the operator
believes they are getting extended thinking while the server ignores or rejects
the field. The table row only records what this application knows.

## Shape

- The flag keeps upstream's name, `--reasoning-effort`, as Bill specified.
- `Engine.ValidateReasoningEffortAsync` runs once at composition, before the GUI
  serves, and only when an effort is actually configured. An empty request sends
  no traffic, so ordinary launches and every offline fixture stay probe-free.
- The model listing's `supported_reasoning_efforts` is read as a hint only: at
  least one server publishes a strict subset of what it accepts, so a value the
  listing omits still gets to answer for itself with a one-token probe.
- Only an explicit refusal naming `reasoning_effort` aborts the launch, and it
  aborts with the provider's own message. A 404, an offline host, a credit error
  or an unreadable listing is not evidence of refusal and launches with a note.
- Both OpenAI request shapes are probed in their own form: `reasoning.effort` for
  Responses, top-level `reasoning_effort` for chat-completions. Other vendors are
  passed through unverified rather than probed with a body they do not define.

## What the curated table still decides

Bill settled the open question this section used to defer: the only thing worth
establishing about an effort level is whether the endpoint **accepts** it, not
whether levels differ. So the table keeps its authority where it actually has
knowledge, and stays out of the way where it has none:

- A listed model with no verified reasoning (`gpt-4o-mini`) is a real negative and
  still refuses an effort. `Tests/ResponsesTests.cs` pins that.
- A listed model still refuses a level outside its row (`gpt-6-astra` and
  `not-a-level`) and still refuses any level on a non-Responses surface.
- A model the table does not name at all is **silence, not refusal**, so the
  requested value is passed through and the endpoint answers for it. This is the
  same rule `port/4-dynamic-media-probe` applies to image capability, and it is
  what makes `--reasoning-effort` usable at all on Bill's local Qwen model.
- Only the effort is delegated. A reasoning *summary* is a different provider
  field with nothing probing it, so an unlisted model still cannot ask for one,
  and no default is invented: an unrequested effort sends no traffic.

## Two defects found while landing that rule

**The field never reached the wire on chat-completions.** `RenderResponses` emits
`reasoning.effort`, but `RenderOpenAI` had no equivalent, so a chat-completions
request silently dropped the level. The provider check would have asked about a
field that the real turn never sent, which is precisely the invisible failure this
feature exists to prevent. `agent/Engine/Render.cs` now emits top-level
`reasoning_effort` for that dialect, matching the shape `AcceptsAsync` probes.

**The GUI hid the refusal it promised to surface.** This note previously claimed a
refusal "stops with the provider's own message". It did not: `Examples/Gui`
catches everything and prints only type names, deliberately, because unknown
exception text can carry credentials or paths. The launch stopped, and the
operator saw `GUI failed: InvalidOperationException. Check private diagnostics...`
and nothing else. That policy is correct and was not touched; instead the
validation call now has its own narrow `catch` that prints the provider's sentence
before rethrowing, because for this one failure the text *is* the answer. The CLI
needed no such handling — its shared handler already prints
`InvalidOperationException.Message` verbatim — and the first draft of this fix
duplicated it there, which was removed.

The CLI also never ran the check at all while accepting `--reasoning-effort` and
`OPENAI_REASONING_EFFORT`; it now calls it for the commands that actually infer,
after `render` and `dump` have returned and skipping `--mcp-pipe`.

## Verification

`agent/Tests/ReasoningEffortTests.cs` drives a fake OpenAI-compatible listener
and asserts: no effort means no requests; the table still refuses an unlisted
model and an unlisted level; an advertised level costs only the listing; an
unadvertised level still gets its probe; an explicit refusal surfaces the
provider's message; an unrelated 4xx, a 404 listing and an offline host all
launch; and a configured level actually reaches the rendered request body.
Full suite on this branch: 342 passed, 4 failed — only the four
`agent/Jobs/Pty.cs` macOS gates. The earlier revision of this note recorded 341
passed and 5 failed on this branch, counting one drifting Chapter 11/19/21
timing case that passed on this run; that earlier figure is its own record and
was not remeasured here.

The acceptance rule was then driven for real, not just through fakes in the
suite. In headless Firefox against the actual `Examples/Gui` server with an
unlisted model and `--reasoning-effort ultra`, the provider saw `ultra` on both
the probe and the ordinary turn and the reply rendered with no page errors — so
the value the check asked about is the value the turn sends. Against a provider
that rejects that value, the same GUI exits non-zero and prints
`provider refused reasoning effort "ultra": Unsupported reasoning_effort:
'ultra'`. The CLI against the same refusing provider exits non-zero with the
provider's sentence exactly once.
