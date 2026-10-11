# Chapter 14: The World's Best Context Engineering, Before Breakfast

This may be the most important chapter in the book.

If you've ever started a new chat because your agent seemed confused,
well, that's yesterday's context engineering. I did it myself, several
times, on the day I designed this chapter. As of September 2026, what
follows is tomorrow's. The book is evergreen, so I hope it always has
modern advice.

The design had been in my head for a couple of months. On a Tuesday
morning I typed it out for about twenty minutes while CodeRhapsody
turned it into a design record and checked it against the code. Then I
took a shower. The first commit landed at 7:42, before breakfast, and
frankly it felt good to get it out.

This is how long-running actors should work: what the actor knows,
organized by section at several levels of detail, with the most
critical parts kept longest, instead of some lame summary of the first
half of the context.

A long session can fill with work the agent has already finished: files read,
commands run, decisions made before lunch. A fresh chat lets a human choose what
to carry over. This chapter gives the running Agent a smaller intervention:
remove recoverable tool traffic while preserving its conversation, instructions
and checkpoint notes. Each cut is recorded, so replay knows what left the window.

The title preserves the ambition of the original design. Its design record
explicitly distinguished the best design its author knew from a demonstrated
best system. The useful test is what the mechanisms preserve, what they remove,
and whether the agent can still do its work afterward.

## TL;DR

The tool portion of context gains a budget; dialogue and instruction survivors
can still grow. Add recorded cuts, per-model keep/stub behavior, a checkpoint
tool, instruction and declaration survivors, and process-crash journaling. Extend
your Chapter 13 solution, preserving serial dispatch and the existing owners.

The shared entry vocabulary gains four kinds. These definitions belong in
`internal/common`; reduction belongs in History, not on common-package structs:

```go
// EntryKind identifies what removes an entry. Survivors never contain tool parts.
type EntryKind uint8

const (
    KindDialogue EntryKind = iota + 1
    KindHandoff
    KindSkill
    KindTools
)

// MicroHandoffData records the checkpoint text that survives its tool call.
type MicroHandoffData struct {
    Text string `json:"text"`
}
```

Each context entry carries its kind alongside sequence, actor and parts. Reuse
the existing skill body and redaction payloads. Add `micro_handoff` and
`tools_changed` facts; the latter carries full `added` declarations and `removed`
names. A `tool_decls` part holds that delta without becoming a call or result.

1. **Keep the prefix stable.** The startup instruction stays frozen. Inline-capable
   models retain the startup tool array too. Other models fold recorded declaration
   deltas into the top-level array when tools change. Do not promise cache grants.
2. **Record actual declaration changes.** Loading, unloading and external MCP
   changes produce `ToolsChanged` entries. Inline declarations must follow the
   complete user/tool-result/ephemeral turn, before its assistant continuation or
   at the request end. The log's chronology stays unchanged.
3. **Give survivors their own kind.** Skill bodies, handoff notes and declaration
   deltas contain no tool-call/result parts. A load result acknowledges success;
   its rendered body is recorded separately and replay never rereads the skill file.
4. **Unload executable capability now.** Remove grants no other active skill
   supplies and stop the skill's servers; retain its instruction entry. This
   chapter does not implement the later instruction-removal verb.
5. **Checkpoint after the whole batch.** `micro_handoff({"text":"..."})` produces
   its call, ordinary result, then recorded handoff. Hold its effect until all
   results from that model reply arrive, including a companion call still waiting
   in the serial dispatcher. Remove all tool parts and append the note once.
6. **Keep the previous batch.** On a `StubsToolResults` row, results larger than
   T/100 are presented whole once, then stubbed unless the next model response
   calls `keep_tool_results({})`. That keeps the preceding batch, never its new
   companions; its own acknowledgement is not stubbed. Record each cut.
7. **Use the ladder for every model.** Whole tool traffic has a T/8 budget; older
   calls and stubs have T/16. When a band exceeds twice its budget, cut an old
   prefix back toward its budget. Work on completed batches already presented;
   retain call/result pairing. Store numeric sequence bounds, not a policy to
   recompute during replay.
8. **Skip only unappliable known events.** Validate before changing context,
   diagnose and skip a malformed familiar fact, then continue with later facts.
   Refuse invalid save JSON, invalid ordering and unknown event vocabulary.
9. **Journal before publication.** Append facts beside the selected save before
   publishing them. Recover snapshot plus anchored tail and journal. Checkpoint
   on normal shutdown: retain the prior snapshot, install the new one, then reset
   the covered journal. Process crashes are in scope; fsync and power loss are not.
10. **Expose the two settings.** `context_target` is bytes: zero selects 400,000;
    positive values clamp to at least 20,000. `log_retention` counts saved events,
    zero keeps all. Use the owning Agent's settings path and preserve unsent browser
    edits when another settings snapshot arrives.

The model-feature matrix is part of the exercise:

| Row | `InlineTools` | `StubsToolResults` |
| --- | --- | --- |
| `claude-opus-5-course` | true | true |
| `claude-sonnet-5-course` | false | false |
| All existing real-model rows | false | false |

The course rows are deterministic fixtures, not provider model recommendations.
The original real-Opus narrative used a different feature configuration; it does
not override this matrix. No live inline-tool or per-round stubbing success is
claimed for the accepted student solution.

Yours: the journal's internal implementation, diagnostic wording and settings
appearance. Keep the existing save/restore contract and actual owners; no new
memory bands, compressor, canonical disk format or generic policy framework.

From the repository root:

```sh
go run ./cmd/grade-edition2 -ch 14
```

The default core is `solutions/edition-2/ch14`; a positional core directory and
existing `-gui` option select another implementation. Section 14.14 gives the
actual launch, settings and eight-check contract.

## §14.1 In Plain Words

The context window is not storage. The log is storage: Chapter 2 made
it append-only, and it captures the facts the agent saw or produced. Saved-log retention can later discard old records. The
window is a working set rendered from that log, and its retained content is rendered again on each request. Cache accounting determines what is charged; bytes alone do not measure that bill. Every byte has to answer three
questions: why is it here, who put it here, and what removes it. A byte
that cannot answer the third question stays forever, and a window full
of such bytes can crowd out the work a person wanted the agent to remember.

Tool traffic can dominate the window. Chapter 2 reports the retained first-edition measurements put tool results at about 42 percent of history by volume and call arguments at another 30 percent. Those are measurements of those sessions, not a universal workload ratio. Much of that traffic is needed to make the immediate decision. The agent reads a
file, decides, edits; after that, the file's
contents are recoverable from disk and the decision lives in what the
agent said about it. Hence the chapter's rule: keep the words, let the
bytes go. Dialogue stays. Tool results go first, then the calls that
produced them, oldest first.

Some things must never go with them. A loaded skill's instructions and
the note an agent writes to its future self at a checkpoint look like
tool traffic today, because they arrive through tools. If they stay
tool traffic, the first cleanup deletes the manual and keeps the tools
it explains. Anything that must survive becomes its own kind of
entry, removed only by its own verb. Survivors are safe by
construction.

Every cut records exact sequence bounds; replay applies those bounds rather than rerunning the policy. That is what keeps Chapter 10's promise:
replaying the log reproduces the same window, byte for byte, even after
the settings that chose the cuts have changed.

Two laws decide where cuts may happen. The front of the request, the
system prompt and the startup tools, is stable so it remains eligible for prefix reuse. For a prefix cache, an earlier change invalidates a longer reusable suffix than a later change. Actual cache grants still depend on the provider and request configuration. So cuts come in steps, not a trickle: one
larger cut now and then invalidates that region less often than a small one every turn.

The redaction levels come from Chapter 2, the frozen prompt from Chapter 9, and the snapshot and replay from Chapter 10. This chapter adds the survivor vocabulary and the policy that decides when to use those existing mechanisms.

## §14.2 Why Is This Byte Here?

The first-edition agent had built redaction machinery before its policy existed. Its `RedactData` event has four levels. `RedactResult` turns a tool
result into a stub and keeps the call. `RedactTool` removes both and
keeps the reasoning around them. `RedactDialogue` and `RedactSummary`
reach into the conversation itself. The original account found no production caller for those levels yet. The dialogue level even says where the
missing piece lives:

```go
case RedactDialogue:
        // Prose and reasoning go. Survivors are defined by the compaction
        // policy, which is a later chapter's problem.
```

This chapter supplies the tool-traffic policy. Dialogue removal remains a later exercise. Context engineering, as the book uses the term,
is managing every byte in the context data structure as well as
today's models allow. The management is a policy: for each byte, a
reason it is in the window, a record of what put it there, and a verb
that takes it out.

Walk a request from a preceding agent after an hour of work and ask each
byte the three questions. The system prompt answers all three: the
agent's constitution, written at startup, removed by nothing. A user
prompt answers them. A 9,000-byte file read forty minutes ago answers
the first question with "the agent needed it once", the second with
`read_file`, and the third with silence. Nothing removes it. It rides
along on every request until the session ends or a human copies the
good parts into a fresh chat.

The fix is to give that old tool traffic a removal policy. This chapter controls those bytes; it does not yet bound the dialogue around them.

## §14.3 Two Laws

The original design reasons about a prefix cache: an earlier change prevents
reuse of more content than a change at the tail. This is a layout rule, not an
invoice formula. A stable prefix alone neither requests caching nor proves that
a provider granted it; the historical live run at the end reported zero cache
reads and writes.

**Law 1: the prefix is frozen.** The instruction stays frozen, and an inline-capable model retains byte-identical startup tool declarations throughout the session. The non-inline fallback below is the deliberate exception for tool changes. Chapter
9 already made the system prompt a constitution and put dynamic
skills in the dialog. It left one leak: a skill that brings tools
changed the tool list, and the tool list is the first thing in the
request. Loading one skill could invalidate reuse from the tool array onward.

The chapter closes the leak with an event. A skill load, a skill
unload, or an MCP server connecting mid-session emits `ToolsChanged`,
a delta of declarations added and names removed. The reducer turns it
into a `Tools` entry at the tail of the context, where it preserves the earlier declaration prefix. On a model whose features row sets
`InlineTools`, the renderer sends that entry as a mid-conversation
system message:

```json
{"role": "system",
 "content": [{"type": "tool_addition",
              "tool": {"type": "tool_definition",
                       "definition": {"name": "gui_click",
                                      "description": "Click a matching element",
                                      "input_schema": {"type": "object",
                                                       "properties": {"selector": {"type": "string"}},
                                                       "required": ["selector"]}}}}]}
```

The accepted course renderer sends the historical header
`inline-tools-2026-09-15`. A removal uses a `tool_removal` block containing
`{"type":"tool_reference","name":"gui_click"}` under `tool`. These are the
course wire shapes verified for this implementation. The October 2026 student
checked the provider's documentation then; this chapter does not assert that a
dated beta is a permanent or currently enabled live capability.
A model without the feature gets the fallback: the renderer folds each delta into the effective declarations and sends that array at the top level. That changes the prefix; no measured cache charge is inferred from the change alone.

The declaration's position matters as much as its bytes. During the student run,
a new MCP declaration followed a tool result, then automatic ephemera added
another user message after it. The roles became:

```text
user, assistant, user, system, user
```

The initial fake accepted that request, and prefix and replay comparisons passed.
A faithful placement check rejected it: the contentful system message belongs
after the complete user turn. The corrected renderer holds declaration blocks
until the assistant boundary or request end, after tool results and automatic
context:

```text
user, assistant, user, system
```

This moves no recorded fact and reclassifies no ephemera as instructions. Only
the disposable wire view changes. The same strict real-Agent/MCP probe then
passes; restoring the old renderer produces the expected placement failure.
Equal replay bytes are useful only after the request itself is valid.

The prefix does change on purpose once in a while. Changing defining instructions requires constructing a refreshed Agent. Changing `context_target` does not rewrite those instructions or the cuts already recorded. Do not add a new refresh command for this exercise.

**Law 2: the cost of changing a byte grows with its distance from the
end.** Under the prefix model, everything following a changed position loses that match. Distance measures the extent of invalidation, not tokens or dollars. Stubbing the result that arrived on the last round trip changes content near the tail.
Stubbing a result a hundred thousand bytes back changes the reusable prefix much earlier.

The original design record preserves a useful failed prediction. The design predicted that
stubbing every tool result one round trip after it arrived would wreck
the cache, since the request changes on every trip. The first session
run that way, on CodeRhapsody with Opus 5.5 on 2026-09-22, reported a
cumulative cache hit rate of 77 percent, after a cold first request
that missed on over 100,000 tokens. The prediction had the law
backwards: a just-finished result sits at the tail, where less later content loses its prefix match. The cumulative figure also hides the steady state, because that
one cold miss is averaged into every later request, which is why a
display of cache rate should show the last request beside the total.

That report concerns one CodeRhapsody session, not a new measurement on the student solution.

Law 2 has a second consequence. A trickle of small cuts deep in the
window breaks that prefix match on every turn; one larger cut now and then
breaks it once. The agent cuts in steps.

## §14.4 What Survives Is Never a Tool Call

Before this chapter, a dynamically loaded skill returns its manual inside the loading tool's result. The first edition illustrated that arrangement with:

```go
result += "## Instructions\n\n" + body
```

The arrangement becomes a bug the moment cleanup can discard that result. Tool results are the first bytes
to go. The first time the ladder stubs old results, the manual becomes
a one-line stub and the tools it explains stay declared. The agent
keeps calling `gui_click` with no memory of the rules for calling it.
A checkpoint note passed as a tool-call argument has the same defect
one level later, because `RedactTool` deletes calls.

The cure is a naming rule, printed in the code as the comment on
`EntryKind`: anything that must outlive tool clearing is its own kind.

| Kind | Created by | Holds |
|---|---|---|
| `KindDialogue` | prompts, hints, model output | text, reasoning, tool calls and results |
| `KindHandoff` | `MicroHandoff` | the actor's checkpoint note |
| `KindSkill` | `SkillLoaded` | a skill's body |
| `KindTools` | `ToolsChanged` | a declaration delta |

Tool clearing only ever touches tool parts, and tool parts only ever
live in `Dialogue` entries. The other three kinds hold none, so tool-part clearing cannot remove their contents. Survivors are safe by construction,
without a list of exceptions for the ladder to consult. `load_skill`
still returns a result, a one-line acknowledgement that the skill
loaded; the body arrives through its own recorded `skill_loaded` fact. The first edition already carried that body in `SkillData`; this implementation makes its separate survival explicit here.

One timing problem remains even though dispatch is serial. The model can return several calls together; after the first runs, others from the same reply are still outstanding. A `micro_handoff`
issued beside a `read_file` clears every tool part in the context. If
it landed immediately, the `read_file` result would arrive afterwards
with its call already gone, and a result without its call is a request
the vendor refuses. So the reducer holds survivors until the batch
completes:

```go
// Holding survivors preserves call/result pairing even under serial dispatch:
// another call from this same model reply may still be waiting to start.
func survivor(c *common.Context, entry common.Entry) {
    if !toolsComplete(c.Dialogue) {
        c.Held = append(c.Held, entry)
        return
    }
    land(c, entry)
}
```

When the last result of the batch arrives, the held entries land in
order. The handoff branch of `land` clears tool parts before appending its survivor. Every removed call now has its result, including the companion that had not started when the note was written.

Unloading changes executable capability immediately. The registry drops grants
that no other active skill supplies and stops the skill's external servers.
The next recorded declaration delta makes the model's view agree with dispatch.
Its instruction entry remains: removal of old instruction text is a different
operation, deferred beyond this chapter. Lifecycle and checkpoint tools stay
available, and a shared grant survives while another active skill owns it.

The prior chapters retained compiled grants on pending unload. This exercise
changes that policy deliberately. Keeping an old manual does not make its removed
tool callable; guessed names must still fail dispatch.

## §14.5 The Layout

With four kinds and two laws, the request has one shape:

```
[ frozen prefix   ]  system prompt; startup tools on inline-capable rows
[ memory region   ]  empty until the agent has memories
[ context         ]  dialogue and survivors, in Seq order
                     ^ oldest                  newest ^
                     cheap to keep            cheap to change
```

The order runs from least volatile to most, which is also the order of
meaning: who the agent is, then what it knows, then what it is doing.
Two independent arguments, one about cache economics and one about
reading order, give the same layout, which is a good sign that the
layout is right.

Survivors land in `Seq` order among the dialogue, at the tail where
they were created. Moving them up next to the prefix would read better
and would edit bytes far from the end. Reordering is a verb of its
own, and this exercise preserves sequence placement except for its explicitly recorded cuts.

The kinds also preserve why an entry belongs in context. The system prompt
and a skill body are instruction: the agent is meant to follow them. A
tool result is data: the agent reads it and decides. A dialogue entry
is conversation. Keeping those apart is why `Kind` records the origin
of an entry instead of guessing it from the text, because text that
looks like instructions and arrived as a file's contents is still a
file's contents.

## §14.6 Visible Reasoning Is the Storage Format of the Self

After the ladder runs, what does the agent remember about a file it
read an hour ago? Exactly what it said about the file. The result is a
stub, the call may be gone, and opaque provider material has its own replay restrictions. This handoff removes tool parts; it does not claim to erase every opaque part. The words the agent wrote in the dialogue are the only record that
survives every level up to `RedactDialogue`. Chapter 2 wrote the rule
into the definition of `RedactTool`: remove calls and results
entirely, keep visible reasoning.

So an agent that narrates keeps its past, and an agent that works in
silence loses it. Compare two turns that make the same edit:

```
(silent)   read_file config.go   edit_file config.go

(narrated) The port is set in config.go; reading it.   read_file config.go
           Port is 8092 at line 40, hardcoded. Moving it to settings.
           edit_file config.go
```

Once both tool bands clear that turn, the silent version may retain no explanation of the edit. The narrated version still says which port, where it was,
and why it moved. Asking for one sentence of intent before each tool call gives the retained dialogue something useful to preserve. It remains a prompt policy whose quality needs observation, not a mechanical guarantee that every narration is correct.

Narration has a failure mode of its own. An agent rereading its old
sentences can trust a figure it wrote down over the file that has
since changed, and a figure written from memory carries the same
confidence as one copied from the source. The defence is in what the
narration records: where a fact lives, the path and line and command,
so the agent's reflex is to re-read the source rather than quote
itself.

## §14.7 The Ladder

The ladder keeps tool bytes in two bands measured back from the tail,
with one number, the target size T, setting both:

```
    stub threshold = T/100    results band = T/8    calls band = T/16
```

The Agent reads T from its owned configuration. Its policy measures serialized tool parts, excluding ordinary dialogue and survivor bodies. The newest
T/8 bytes of tool traffic stay whole. Older results become stubs by
`RedactResult`, and their calls stay, for the next T/16 bytes. Tool
bytes older than both bands go entirely by `RedactTool`, leaving the
dialogue around them. At the default T of 400,000 bytes, without converting it into an assumed token count, the stub threshold is 4,000 bytes, the
results band is 50,000 and the calls band 25,000.

Law 2 sets the rhythm. A band is cut when it reaches twice its budget,
and one event cuts it back to its budget, so cuts arrive as occasional
steps. The two doubled budgets sum to 3T/8. That arithmetic describes the eligible tool bands, not a hard maximum for the whole request. A newly presented large result must survive once, whole batches constrain cut boundaries, and dialogue and survivors have no budget here.

Three rules keep the ladder honest.

- **Every cut is an event.** The Agent runs policy before rendering and records `redacted` facts through its existing append path. History alone changes the projection; Engine renders it. Its `To` field is a number, a `Seq` in the log.
  The watermark is computed once, at the moment of the cut, and never
  again. Replay applies the recorded number, so lowering T tomorrow
  cannot reach back and cut yesterday's session differently.
- **Nothing is cut before the model has seen it.** A result goes out
  whole at least once, however large. The ladder only reaches entries
  already sent.
- **No call loses its result.** The calls-band cut snaps back to a
  point where every call removed takes its result with it.

The ladder applies to every model, because it needs no judgement from
the model. It is also blunt: a 40,000-byte file read that the agent
will never look at again stays whole until the band passes it.

## §14.8 Keep or Stub

A capable model can do better than the ladder, because it knows which
results it is finished with. On a model whose features row sets
`StubsToolResults`, every tool result above the stub threshold becomes
a stub in the request after the one that carried it, unless the
actor's next message calls `keep_tool_results`. The keep takes no
arguments and keeps the whole batch.

The keep governs the batch that has already arrived, the one the
actor is looking at. Called alongside other tools, it keeps the
previous batch, never the one it rides in, since those results do not
exist yet when the actor decides. Its own result is never stubbed. The
call survives either way, so the agent can see what it asked for and
re-run it. The stub is an ordinary `RedactResult` event, so replay
reproduces it without asking the model again.

Law 2 favors this placement. The stub replaces bytes that arrived one round
trip ago, at the tail, which is the cheapest place in the window to
change anything. The 77 percent in §14.3 was measured under this rule:
CodeRhapsody, the agent that co-wrote this book, runs it with its own
thresholds.

The feature is a column in the model table because judging which evidence to
retain is a model behavior, not a transport capability. The original account
rated particular Opus and Sonnet versions differently. Those judgments are
historical observations, not authorization to enable the flags on this exercise's
real rows. The course Opus row exercises the mechanism; every existing real row
keeps both new flags false.

> **Open risk.** Per-round-trip stubbing may interfere with the
> model's thinking: a model reasoning across several round trips has
> its evidence removed between steps. Nobody has measured it. The
> experiment that would settle it is one task, one model, stubbing on
> and off, with the outcomes compared by quality rather than by bytes.

## §14.9 From compress_context to micro_handoff

The ladder removes bytes the agent is done with. It gives the agent no
way to say what it needs in order to keep going. A checkpoint does,
and the book's own lineage tried two designs for it before this one.

CodeRhapsody's first design was a tool called `compress_context`: the
model chose a range of messages and replaced them with its own summary.
Bill described how that went on 2026-09-13:

> "compress_context was the old system that you (Claude) did pretty
> well, picking a range of messages to summarize, but freaking Gemini
> almost always deleted 80% of messages, starting with message 1, with
> a terrible summary, lobotomizing the LLM, so we switched to handoffs
> instead."

The failure has a shape. The cut ran by position, from message 1, and
message 1 is where the user stated the goal. The summary was written
by the model that was about to lose the originals, so nothing could
check it, and summarized dialogue is in no file the agent can re-read.
Bill's account describes destructive choices in those trials. It motivates limiting discretion to recoverable traffic without making a universal claim about that model family.

The handoffs that replaced it, `handoff_task`, wrote a structured
document and started a fresh instance from it with an empty window.
That design assumes continuity lives in a note to a stranger. It also
added a second path for writing the agent's durable state, and second
paths drift: CodeRhapsody's handoff path never triggered the memory
cascade that its `save_memory` path did, and the gap sat in its notes
as a known bug for months. Ensemble builds neither tool.

`micro_handoff` keeps the same instance and the same dialogue. The
actor writes a note, and the reducer removes every tool call and tool
result in the context and appends the note as a `Handoff` entry. Two
lessons from the older designs are built in. The model's discretion
covers only what is recoverable: files and tool operations can often be revisited, and the dialogue, which cannot, stays. The quality of that judgement still needs observation on each model. Unlike automatic stubbing, the handoff tool is available across the exercise's model rows.

In the log it is three records: the tool call, an ordinary result
that acknowledges it, and a `MicroHandoff` event carrying the text.
The event is what the reducer acts on. The result is an
acknowledgement for the same reason `load_skill`'s is: a tool result
is tool traffic, and the note must outlive tool traffic. Replay sees
the event and makes the same cut.

What goes in the note decides whether the checkpoint works. The
version CodeRhapsody uses asks for fields, and its tool description is
blunt about which one matters. Of `tried_and_failed`, it says: "Highest
value per byte in the document: it is the only field that prevents
repeating a mistake, and errors are the part of a record most likely
to be discarded as noise." It also says when to call the tool: "at a
completed micro-goal, never at a token threshold." A checkpoint taken
mid-task drops the tool result that held the task's state. For this exercise, the handoff clears tool traffic while retaining conversation. A useful note still records conclusions and failed approaches; it cannot rely on a future model inferring them from discarded results.

## §14.10 Compaction Is Described, Never Performed

The policy does not edit context directly. It records a cut; the History reducer applies it, and the retained log preserves the decision. With the full required log retained, prior windows can be reconstructed, including ones the current settings would cut differently. Log retention can remove the evidence needed to reconstruct older windows; a current snapshot does not restore that missing past.
The grader's `ladder-is-recorded` check does exactly that: it replays a
session under a changed T and expects the same requests.

A retained log can contain a known event the current reducer cannot apply: a redaction naming an entry a later
handoff removed, a payload a newer version malformed. The previous loader refused such a known invalid fact. The accepted History
owner now has a small replay boundary:

```go
func (h *history) replay(e common.Event) {
    if err := validateFor(&h.context, e); err != nil {
        h.parent.Ensemble().Logf("restore: skipped event %d: %v", e.Seq, err)
        return
    }
    apply(&h.context, e)
}
```

The skip is clean because validation checks an event before it mutates
anything, so a rejected event leaves the context exactly as it was.
Totality covers known event types. An event type the code has never
heard of is still refused at load, loudly, because a log from a newer
agent is a different problem than a bad record in a familiar one.

The original design also considered provider-side context editing. Its dated
API documentation described a server editing the input while the client kept its
full history. That can coexist with a local log, but it moves the decision about
what the model reads outside this reducer. Here, the chapter needs a recorded
keep decision and reproducible local projection, so the existing event path owns
the cut. This choice does not require a second editing API or a new provider
feature in the exercise.

## §14.11 Keep the Words, Let the Bytes Go

A stub removes bytes from the window. It removes nothing from the log.
The `ToolReturned` event that carried the result stays in the log
until log retention truncates it, and in the common case the bytes
were never unique anyway: the file is still on disk, and the operation can often be run again. That does not promise identical output from a changed file or a command with side effects.

The stub carries no address. The reason is a bug in the agent that
co-wrote this book. CodeRhapsody's stubs cite the output file of the
command they replace, a path like `cr/io/26`. When this chapter was
designed, its handle counter started again at 1 on every launch, so
after a restart `cr/io/26` named a different command's output. The
agent that followed the stub read a file, got plausible output, and
had no way to tell it was the wrong one. Chapter 12's rule applies:
wrong data is worse than no data. The original account reports a subsequent counter correction. Ensemble's stubs say the bytes are gone and
leave recovery to the tools that produced them.

## §14.12 Crash-Safe Persistence

Chapter 10 saved at shutdown. A process that dies mid-session loses
everything since the last save, including work whose results have already reached the model.

The current Agent appends each accepted event to a journal beside the save file,
`<save>.journal`, one JSON line per event, as the event happens.
Recovery is Chapter 10's load with one more source: `Recover` reads
the snapshot, applies the tail after its `as_of` anchor, then applies
the journal. A normal shutdown writes a new snapshot and resets the
journal. The grader's `crash-recovery` check runs a clean session,
then runs a second for two turns and kills it with SIGKILL, restarts
it, and checks that journaled work survives the restart. This does not certify every possible crash point or recover a malformed partial write by guessing.

Two orderings are load-bearing, and neither is commutative. The new
snapshot is written before the log is truncated: in the other order, a
crash between the two steps leaves neither the old events nor the new
snapshot. And the log is never truncated past the snapshot's anchor,
because the anchor is where replay starts. `log_retention` is a saved-history preference; the anchor is a correctness boundary, and the
anchor wins. The accepted History owner installs its checkpoint in that order and keeps the previous save at `<save>.bak`. The Agent supplies the selected path; the GUI does not own a second persistence loop.

Two things are out of scope. The journal is never fsynced, so power
loss can take the last few events; "crash" means a process crash.
And truncation deletes old `ToolReturned` events, which are the only
on-disk copy of a redacted text result. The snapshot can still restore the current window, but retention may make older windows and raw results unavailable. Zero retention keeps the full saved log; it does not promise bounded disk growth.

## §14.13 One Knob

The tool budgets run off one setting. The Context Management tab in the
settings panel writes `context_target`: T in bytes, 0 for the default
of 400,000, clamped to at least 20,000. The original budget comment explains its floor: below it, "the stub threshold would fall under the size of a
directory listing". The tab's other field, `log_retention`, is the
number of events the saved log keeps, 0 for all.

The GUI's file is `DataDir/settings.json`; the supplied application chooses
`cr`, so use `cr/settings.json` from its launch directory. A file named
`settings.json` beside the executable is not interchangeable. The browser's
public settings update reaches the owning Agent before subsequent policy work.

A user changes the tool-traffic target, not a ceiling on every kind of context.
Every threshold scales with T. `log_retention` separately controls saved history;
changing either setting cannot rewrite already recorded cuts.

The settings browser also exposed an editing race: an acknowledgement for one
field arrived while another number was typed but not yet committed by blur.
Replacing every field from the snapshot erased that unsent edit. Keep effective
settings updates while protecting the person's pending input. The corrected
two-tab test catches the restored defect; an echoed setting alone would not.


The retained first-edition coder review reports a synthetic reference session: the fake
vendor scripted four prompts, each followed by fifteen `read_file`
calls on 8,000-byte files, 64 requests per run.

| run | last request | total sent | cuts |
|---|---|---|---|
| no ladder (T = 10^12) | 542,600 | 17,470,874 | none |
| ladder only (T = 400,000, Sonnet row) | 143,152 | 6,172,202 (35%) | 7 `RedactResult` |
| ladder and stubs (T = 400,000, Opus row) | 62,456 | 2,334,066 (13%) | 59 stubs |

That review reports the largest request in the ladder run as 146,378 bytes in total,
below 150,000 bytes, the sum of the two doubled tool-band budgets. As §14.7 explains, that sum is not a universal upper bound. The
prefix, 2,991 bytes, was byte-identical across all 64 requests in all
three runs. The stubs run never triggered the ladder, because stubbing
kept every band under its budget.

Bytes are the easy half. The table says nothing about cache rates on a
real vendor, nothing about whether the answers got better or worse,
and nothing about the open risk in §14.8. The following historical live account supplies a different observation; it does not turn these byte totals into cache savings or a quality comparison.

## §14.14 The Exercise, Graded

Extend the current core and optional GUI. The Agent chooses policy and follows
its configuration; History validates, journals and reduces; Engine translates the
result into vendor requests. Keep those existing ownership paths rather than
injecting a settings callback or moving runtime behavior into common.

The grader drives the actual GUI executable through newline-delimited stdin
prompts, using the existing envelope:

```json
{"kind":"prompt","text":"Load manual and inspect the prepared files."}
```

It supplies provider environment, the two course model rows, complete skills and
GUI assets. Public settings updates establish the effective target before traffic;
the GUI persists them at its own Agent path. Normal EOF drains work, snapshots and
exits successfully; the crash scenario observes completed work before SIGKILL.
These are the existing human/client boundaries, not a new grading protocol.

| Check | Points | Observable obligation |
| --- | ---: | --- |
| skill-survives-the-ladder | 15 | The manual remains after both tool bands cut old paired traffic |
| micro-handoff-shape | 15 | Call/result/event, one note, no orphan including a same-batch companion |
| keep-or-stub | 10 | Previous-batch keep, next-batch distinction and non-stubbing model control |
| ladder-is-recorded | 15 | Recorded numeric cuts survive restart under changed T |
| frozen-prefix | 10 | Stable inline prefix, actual skill/MCP declarations and valid placement |
| replay-equals-snapshot | 15 | Snapshot and sufficient log replay produce the same next vendor bytes |
| crash-recovery | 15 | Completed journaled work survives an actual killed process |
| total-reducer | 5 | Bad known facts are skipped; later good facts still apply |

The original rubric has eight checks totaling 100; there is no additional parity
row. Preserve earlier behavior through its affected regressions. Diagnostics, validation before mutation, disk failure
ordering, backups, actual GUI edits and the stricter automatic-context placement
also have source/integration proof; a check title alone does not establish all
clauses in a paragraph.

The original student score remains **45/100 FAIL**. Three failures came from a
fixture writing the wrong settings path, so the intended small target never
reached the Agent. Another requested undefined `think` rather than supplying a
real capability. Those masked later assertions; they did not prove replay broken
or certify its remaining obligations. The revised fixtures use the effective
settings path and existing `list_directory`, while supplemental real-Agent checks
establish the valid behavior. No compatibility tool is added to the solution.

The first-edition audit found two other useful holes. Short file paths did not
fill the calls band, so a supposed two-band test reached only result stubbing.
Bad events at the very end let a loader stop early without losing any later good
content. Exercise both watermarks and put valid work after a rejected event.
The student's inline-ordering failure supplies the same lesson at another
boundary: a fake that accepts invalid provider requests can certify replay of the
same mistake.

## §14.15 A Historical Run, and What It Established

The first edition records a September 23, 2026 reference run with
`claude-opus-5` and both new features enabled. Its request table and journal
excerpts are retained here as that chapter's reported evidence, not a repeat on
the accepted student binary. This configuration differs from the explicit course
exercise, which leaves existing real rows false. The historical run used:

```json
{"context_target":20000}
```

Twenty thousand is the floor the settings clamp allows. It puts the
stub threshold at 200 bytes, the results band at 2,500 and the calls
band at 1,250. The workspace held seven files copied from
`agent/internal/common`, 60,865 bytes in all, three times the target.
The prompt:

```text
Read each .go file in src/ one at a time with read_file (budget.go,
context.go, save.go, journal.go, model.go, event.go, part.go). Then
answer in three sentences: how does this agent keep its context from
growing without bound?
```

The recorded account reports eleven requests and ten tool calls: a directory
listing, seven reads, two reads of files already read, and the
answer. It prints nine `redacted` events from the journal, all of one shape:

```json
{"seq": 12, "type": "redacted", "time": "2026-09-23T14:44:45.308867Z", "redact": {"from": 6, "to": 6, "level": "redact_result", "reason": "round trip: not kept"}}
```

Ten results, nine stubs. The tenth result was still in its one full
request when the turn ended. Its table of reported input tokens, request by request:

| Request | Carried in full        | Input tokens |
|--------:|------------------------|-------------:|
| 1       | the prompt             | 3,787        |
| 2       | directory listing      | 3,955        |
| 3       | budget.go, 1,683 B     | 4,560        |
| 4       | context.go, 18,855 B   | 10,889       |
| 5       | save.go, 6,282 B       | 6,943        |
| 6       | journal.go, 4,522 B    | 6,347        |
| 7       | model.go, 7,347 B      | 8,475        |
| 8       | event.go, 10,742 B     | 10,109       |
| 9       | part.go, 11,434 B      | 10,655       |
| 10      | budget.go, again       | 7,112        |
| 11      | context.go, again      | 13,522       |

Each file is a tooth: it rises for one request and falls on the next.
The floor under the teeth rises too, because every call, every stub
and every sentence the model wrote stays. The first request, before
any file, was already 3,787 tokens of system prompt and tool
declarations. The target caps neither of those. It sets three numbers,
and all three govern tool bytes.

The ladder never fired. It never cuts a result that has not yet gone
out whole once, and the stub rule had already removed every older
result, so no request held anything the ladder was allowed to cut.
Ten calls with short paths never reached the calls band's trigger.
On a stubbing model under this workload the ladder is the backstop.
On a row without automatic stubbing, the ladder does the tool-byte cleanup.

The interesting part is the dialogue. The account quotes the model's narration from its journal:

```text
Noting: context.go — Context has bounded fields; redaction/summary/handoff clear tool parts. Next file.

model.go: a per-model feature table; the context-relevant column is `StubsToolResults` — "every tool result above the stub threshold becomes a stub in the request after the one that carried it, unless the model's next message calls keep_tool_results" (on for opus-5/opus-5-course; off for sonnet-5, where "the ladder" alone applies). Next file.

part.go: content is `[]Part`, and `RedactedPart` is the *result* of a redaction — its stub is synthesized deterministically by the reducer (no stored growth) and it carries forward the superseded part's `Ref`, so dropped bytes stay fetchable. Now re-reading the two files whose contents got stubbed out of my own context.

budget.go confirmed: one knob T, with threshold=T/100, results band=T/8, calls band=T/16. Now the policy that spends them.
```

The quoted `Ref` claim is the model's narration, not a storage guarantee. It
conflicts with this chapter's address-free stub rule. A restarted job handle can
name different bytes, and saved-log retention can discard the original result.
Keep the observation without teaching the model's explanation as truth.

The two files it read again were `budget.go`, the one file it wrote
no note about, and `context.go`, the largest file, which got one line.
In this reported run, the files with fuller notes were not read again. The note serves as the model's summary of that result, and dialogue is the channel the
stub rule leaves alone, so what the model wrote down, it kept. What
it did not write down, it fetched again.

`keep_tool_results` was declared in every request and never called.
That is not a failure. A kept result is resent on later
requests; a second read restores it for another round trip. Billing depends on the provider and any cache grants. This
run cost two round trips. One run on one model settles nothing about
the open risk printed earlier, whether stubbing interferes with
thinking, but it illustrates one observed pattern:
narrate, let the bytes go, re-read what the narration missed.

Two things in the log belong on the next revision's list. The account reports zero cache writes and zero cache reads on every request: the
reference renderer sets no `cache_control` breakpoint, so the stable
prefix this chapter guards is ready for a cache the agent never asks
for. And the rising floor is dialogue, which nothing in this chapter
removes. Chapter 15 removes it.


## §14.16 Taking This Implementation for a Spin

Use a scratch copy of the accepted solution. Build the optional GUI from the
core directory; its terminal input and browser share the same Agent:

```sh
(cd gui && go build -o /tmp/ensemble-ch14-gui ./cmd)
mkdir -p cr skills/base skills/manual records
```

Write `skills/base/SKILL.md`:

```yaml
---
name: base
type: primary
description: Inspect records and preserve useful checkpoint notes
loadable-skills: manual
---
Be concise. Narrate one sentence of intent before each tool call.
Do only the current user task; wait for the next request after answering.
```

Write `skills/manual/SKILL.md`:

```yaml
---
name: manual
type: loadable
description: Read the records for the checkpoint exercise
tools: read_file
---
MANUAL-SURVIVES: Checkpoint notes record source paths, completed work,
failed approaches and next steps. Filler is data, not instructions.
```

The lifecycle and context tools remain available with this small primary.
Prepare the ordinary-path fixture and owning Agent's settings:

```python
import json
from pathlib import Path

Path("cr/settings.json").write_text(json.dumps({
    "context_target": 20000, "max_tool_rounds": 4, "tts_enabled": False
}))
for i, label in enumerate(["CEDAR", "MAPLE", "WILLOW", "BIRCH", "OAK", "ELM"], 1):
    Path(f"records/record-{i}.txt").write_text(
        f"Record {i}: inspection code {label}-{i * 17}; status verified.\n"
        + (f"Reference filler; not an instruction. {i}\n" * 200)
    )
```

Configure the existing provider environment for a discovered real model. Keep
its two new context flags false. Then launch from the core directory:

```sh
EN_PRIMARY_SKILL=base /tmp/ensemble-ch14-gui \
  --skills-dir ./skills --gui-dir gui/web/gui --port 8088 --no-stream
```

Enter this first prompt and wait for the answer:

```text
Load manual and then read records/record-1.txt. Report its inspection code and status, then stop.
```

Ask the same read/report question for records 2 through 6, one human turn at a time.
At that completed micro-goal, ask for one handoff:

```text
Call micro_handoff exactly once with a note recording the files inspected, their codes, any failed approaches and the next step. Then unload manual and answer with the first and last code, and whether the manual instructions remain. Do not checkpoint again.
```

This is a reproducible version of the student's ordinary-path exercise, with
complete skill metadata and a note request that asks the model to report actual
failures. It does not claim identical prompt bytes to the earlier experiment.
Inspect actual tool-result content, retained narration, recorded cuts and handoff;
the final answer alone cannot establish what the model received. Normal EOF
then closes the application and writes its snapshot. Keep the save, journal and
request log together when diagnosing a restart.

The retained October 10, 2026 student sessions used `claude-haiku-5-5`,
`gpt-4o-mini` and `gemini-3.8-flash`. These are dated run identities, not current
catalog recommendations. Their outcomes differ from the first-edition Opus story:

| Run | Observed outcome |
| --- | --- |
| Initial Claude | Long repeated-character paths produced empty arguments and an overlong path; no useful reads, cuts or handoff |
| Initial GPT-4o mini | All six result codes arrived, but repeated handoffs reached the 14-batch bound; no successful final answer |
| Gemini 3.8 | Six reads, five result cuts, two paired tool cuts, one handoff and unload; correct first/last codes and recognition that the manual remained |
| Claude with ordinary paths and separate turns | All six full results arrived, then six result cuts and one tool cut; later denied having read the redacted files and refused the handoff |
| GPT-4o mini with ordinary paths and separate turns | All six full results arrived, six result cuts and one tool cut; handoff/unload ran, but repeated checkpoints reached the final turn's four-batch bound |

All five processes exited zero. Only Gemini completed the requested task. The
[retained results](../../docs/edition-2-notes/ch15.md#live-use-and-its-limits)
distinguish successful mechanisms from unsuccessful model behavior; no further
paid run was spent to make the table look better.

Claude's corrected-path run makes the risk concrete. It first received the codes
in actual provider-bound result blocks. After those results were redacted, it said:

> That file was redacted, so I never read it.

The transcript disproves the claim. The ladder had removed already-presented
bytes; it had not prevented the read. A smaller request is easy to measure. Whether
the model interprets that smaller working set correctly needs the task and the
actual received results beside it. The recorded cuts make that diagnosis possible;
they do not make the model's recollection true.
