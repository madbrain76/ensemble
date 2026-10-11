# Edition 2 Chapter 11 enhancement: MCP

Author: `edition2_author`, October 11, 2026. Original Chapter 12 becomes
Edition 2 Chapter 11, matching frozen `solutions/edition-2/ch11`.
This is an authorized enhancement of the original, with revised grader work
owned by `edition2_grader_coder` and independent review by `reviewer_astra`.

## Reading and source

The author context compacted at the start of this chapter. Before any new
writing, all four required documents were reloaded in full, without omitted
ranges: `docs/agentic-codebook-workflow.md` (1825 lines),
`docs/edition-2-carryover.md` (871), `book/voice.md` (632), and
`book/chapter-writing-procedure.md` (599). They remain together in active
context. The full `.agents/skills/ensemble-coding/SKILL.md` was also reread,
with its truncated final range explicitly recovered, plus current AGENTS,
enhancement record and chapter map. No further compaction preceded this draft.

Read original `book/chapter-12.md` in full, file revision
`ec41c6e8cf0800e266f7a39ebadf255d94379082`, SHA-256
`a754bd50ca7f3c7a2da868ff5acbdb98b4d5eddaa05349e7890c4a3dbf6cd4b9`;
the complete [student record](ch12.md) and [comparative review](ch12-code-review.md);
the accepted direct snapshot's common MCP interfaces, manager, client,
transports, Agent invocation/collection, browser transport/tools, completion
and server-lifecycle tests, README and public construction/dispatch paths.
Read accepted Chapter 10 for presentation continuity. Historical raw live,
browser and slow-peer evidence was inspected alongside the source and reports.
No main-tree future capability or first-edition runtime was edited.

## Outline and preservation

Through-line: the reader wants to add an external capability without rebuilding
the agent, while preserving trustworthy execution and fresh context. The
original has no documented personal incident to reconstruct. Its compiled-tool
opener, transport progression and sharp explanation of silent truncation and
unreadable control states carry the teaching. Those remain, with enough context
to explain why each boundary matters. No anecdote, dialogue or result was invented.

Order remains extension argument → transport → correlation → MCP client →
automatic tools → bridge → reverse calls → WebSocket/browser → skill lifecycle
→ exercise, followed by the actual student demonstration and corrections.
The small original chapter earns additions for the actual owners, missing wire
and launch facts, and observed timing failure; there is no reduction target.

Section facts that earn their place: one transport interface preserves optional
GUI reuse; a reverse call can deadlock its own reader; a callback can return
before execution ends; a 3.4-second snapshot was replaced by a three-second
running notice; discovery does not grant every offered capability; reverse
admission exposes enabled Agent tools; a textual speech queue cannot prove
audibility; stopping a peer differs from retaining its manual; a marker in call
arguments cannot establish the returned result.

## Disposition of the student and reviewer findings

- Corrected the original's automatic admission claim: per-server skill grants
  precede Agent enabled-tool policy, and hidden automatic names are refused by
  ordinary model/reverse dispatch. Explained reverse trust separately from the
  offered-tool filter, without promising a sandbox.
- Replaced copied sibling services and callback-bag sketches with the real
  owner chain and common/public interfaces. Kept the simpler combined correlation
  client, with out-of-order reply and reader-liveness reasons. Internal type and
  method spelling is explanatory, not a new grader requirement.
- Replaced the original bridge's `context.Background()` with an excerpt using
  the actual job context. Distinguished callback deadline, explicit kill, turn
  interruption and shutdown. Retained the slow automatic-completion defect and
  narrow correction, including delayed output retention without stale attachment.
- Published newline framing, integer correlation, dedicated `--mcp-pipe` behavior
  and historical `2024-11-05`. Reverse calls, `ephemeral` and `$/cancelRequest`
  are explicitly course contracts; no current MCP conformance claim is made.
- Kept the four browser tools and the original state/truncation explanation.
  Clarified the 4096-byte UTF-8 cap, actual per-connection delivery, one active
  browser target and speech-observation limit. Source/real-browser evidence is
  separate from a faithful JSON-RPC peer.
- Made frontmatter valid and explicit about descriptions, types and tool grants.
  Explained actual subprocess startup/reaping, workspace-relative paths, literal
  arguments and credential separation. External unload is immediate; compiled
  tools/instructions retain the earlier pending policy. No compaction was added.
- Added a small reproducible marker peer/primary skill and real CLI launch.
  They reproduce the historical pattern, without claiming identical original
  generated bytes. Retained the empty Anthropic tool-argument failure, exit-zero
  limitation, one bounded correction repeat and missing initial response capture.
- Preserved the argument-forwarding assertion survivor as a teaching example:
  inspect correlated result blocks rather than searching the whole request.
- Distinguished browser applied/effective font state from disk persistence.
  Fifteen actual Chrome checks used a provider fake; three live-provider CLI
  sessions used a local MCP peer. Neither is relabeled as the other.
- Preserved original full interrupted/unscored exit143 and interim-only80/90,
  with recursive parity unrun and vacuous behavior credits disclosed. Original
  source review accepted one correction round; it remains distinct from this
  manuscript/grader enhancement gate.

Deferred/declined: no resources/prompts protocol, full YAML processor, compaction,
sandbox, generalized codec framework, future GUI automation behavior, new public
mode, new runtime tests or paid repetition. Existing focused evidence covers
valid source obligations beyond the revised numeric grader's fixtures.

## Validation and remaining gate

Direct `go run ./cmd/lintprose book/edition-2/chapter-11.md` passes all hard
checks. Its lower-length and person-gap warnings do not justify padding with
invented stories. The original text and all frozen runtime exports remain intact.

Agreed grader contract retains weights15/15/15/20/15/10 plus10 corrected
Chapter10 parity. Dedicated raw pipe covers handshake/reverse, real CLI plus a
skill-declared stdio peer covers model declarations/results/automatic context,
and the actual optional GUI covers its tunnel. Existing grader-only launch
options remain; no new product interface is required.

Exact manuscript Go, Python, frontmatter and command validation, final revised
grader results and independent combined acceptance are pending. Do not call
this draft accepted or advance Chapter12. Final bounded results will be recorded
here with links to the [coder report](enhancement-11-code.md) and independent
review. Mandatory technical handoff: `.agents/skills/ensemble-coding/SKILL.md`.

Initial independent prose review of `685f0bfe` finds no blocker. During sample
reconciliation, the coordinator and coder flagged the inherited generic
`./ensemble --mcp-pipe` path beside the later explicit `/tmp/ensemble-ch11`
build; the accepted core also has an `ensemble/` library directory. The exercise
now states the core-module build and uses `/tmp/ensemble-ch11 --mcp-pipe`
consistently. This is executable-location clarity, not a claim of a failed
protocol or changed runtime requirement. Exact execution remains pending.
