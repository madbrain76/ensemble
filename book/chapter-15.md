# Chapter 15: Forgetting on Purpose

Memory is the most important feature I never knew I needed.

In the spring of 2026, I asked CodeRhapsody to interview another
agent we had built, a family assistant bot, looking for ideas to
improve the code that agent relied on. Both agents had memory: a
persistent record of past work that survived across sessions. The
instant improvement was impossible to overstate. Suggestions I had
never considered, each building on the context of the last.

I took the exact same code and methodology to work the following
week and tried to get my agents there to self-improve. It was a
miserable failure. I blamed Gemini, because Claude was doing great
at home. Then I discovered that my memory code had a fatal bug that
broke it entirely at work. None of my Gemini agents had memory at
all.

I fixed the bug and reran my experiments. In a rare reversal, Gemini
outperformed Claude. One agent invented procedural memory, what
this book calls skills, and implored me to build it. I was too
ignorant to realize the importance of what it had proposed. Other
agents fixed broken tools and motivated improvements I never would
have thought of.

That week I wrote an internal paper for my colleagues: "All AI
Coding Agents Must Have Memory NOW!" It was one of my most widely
read papers.

The difference between failure and breakthrough was not the model,
the prompt, or the temperature. It was memory.

Without persistent memory, a long-running agent can meet a stranger every
morning: the new session has the code and lacks the reasons for it. Chapter 14
removes old tool traffic. This chapter compresses completed conversation into
memories that survive a restart. What those memories preserve still needs checking.
A shorter conversation can lose the reason for a decision, or retain that reason
and fail to use it.

## TL;DR

The agent gets memory that the model reads on every request. Five
bands hold text between the frozen prefix and the conversation,
ordered from most stable to most volatile. A
checkpoint triggers compression: the oldest finished work becomes a
memory, and the raw conversation it replaced leaves the context.
When a band fills, its oldest contents graduate into the denser band
above. Every compression is a recorded event carrying the compressed
bytes, so replay reproduces the bands without asking a model. Forcing
is by tool removal: a nearly full context gets `micro_handoff` as its
only tool.

Shared vocabulary belongs in `internal/common`; Memory owns corpus policy and
files through its required Agent parent. The Agent actor commits facts and worker
results, History projects them, and Engine performs the bounded compressor call.
The optional GUI continues to use public interfaces.

```go
// Band orders memory from stable identity to recent completed work.
type Band uint8

const (
    BandSoul Band = iota + 1 // SOUL.md, maintained by the human
    BandMemory              // MEMORY.md, terminal band
    Band64x
    Band8x
    BandSession
)

// MemoryFileID is independent of any event log's sequence numbers.
type MemoryFileID struct {
    Date string `json:"date"`
    Num  int    `json:"num"`
}

// BandPopulatedData records actual replacement or disk-read bytes.
type BandPopulatedData struct {
    Band   Band         `json:"band"`
    File   MemoryFileID `json:"file"`
    Text   string       `json:"text"`
    Thru   MemoryFileID `json:"thru,omitzero"`
    Source string       `json:"source,omitempty"`
}

// BandDepopulatedData retires a prefix, or an entire disabled band.
type BandDepopulatedData struct {
    Band   Band          `json:"band"`
    Thru   *MemoryFileID `json:"thru,omitempty"`
    Reason string        `json:"reason"`
}

// CompactorLaunchedData correlates an attempt with its eventual population.
type CompactorLaunchedData struct {
    Band Band         `json:"band"`
    Thru MemoryFileID `json:"thru"`
}
```

Add `KindSoul`, `KindMemory`, `Kind64x`, `Kind8x`, `KindSession`, and the event
names `band_populated`, `band_depopulated`, `compactor_launched`. Each memory
entry has `Seq: 0`, its own kind and file identity, and text rendered as data.
The log event still has its ordinary sequence number.

The application's memory settings live at the owning Agent's settings path:

```json
{"memory": {
  "conversation": {"budget": 3000},
  "session": {"budget": 1500},
  "8x": {"budget": 1500},
  "64x": {"budget": 1500}
}}
```

Budgets and watermarks are bytes. Zero selects defaults; optional
`high_watermark` and `low_watermark` default to twice and once the budget.
`disabled: true` turns a band off. `ContextWindow` in the model table is tokens.

1. **`micro_handoff` is the sole trigger.** After stripping tool
   calls and results (chapter 14's mechanism), measure the remaining
   `KindDialogue` text. If it exceeds the conversation budget,
   select one or more whole completed checkpoint segments, oldest first,
   including their handoff notes. Compress after the successful tool batch,
   before the next main request; leave the unfinished tail and independent
   skills/declarations alone.
2. **The compressor gets one tool, `submit`, and writes in first
   person.** It receives the stripped conversation text and submits
   `{"memory":"compressed text"}`. Require exactly one valid, nonempty
   submission from this finite direct call. No tool loop or recursive retry.
3. **Two events record the compression.** A `redacted` event with
   level `redact_summary` and an empty replacement removes the selected
   dialogue and its handoff. A session `band_populated` adds the
   compressed memory as its own entry. Hold the cut until its matching
   population arrives. No text appears twice.
4. **Band entries sit at `Seq: 0` and are ordered canonically** by
   `(band, date, num)`, ahead of the conversation. Switching bands
   off and restoring them in any order must produce the same context.
5. **Graduation is oldest-first and watermark-based.** When a band
   crosses its high watermark, compress its oldest files until the source
   reaches its low watermark, taking at least two files. Emit
   `band_depopulated` then `band_populated` carrying the replacement bytes.
   Measure the destination and continue if it too crossed its watermark:
   session to 8x, 8x to 64x, 64x to MEMORY.md. Keep retired files outside
   the active scan pattern. SOUL never graduates; MEMORY is terminal.
6. **A band will not graduate into a disabled band.** Refuse and
   surface a warning while preserving the source. Disabling either source
   or destination during compression also prevents that reply from committing.
7. **Replay applies recorded events.** A log with compaction events
   rebuilds the context with zero vendor calls. The bytes and cuts are in
   the events; rendering neither recomputes them nor reads today's files.
8. **An abandoned compaction is dropped.** Emit `compactor_launched`
   when a compressor starts. A launch at startup with no matching population
   carrying the same band and `Thru` is abandoned: the old source stays
   live, and the next checkpoint measures again. Failed or malformed
   submissions likewise leave the source available.
9. **Forcing.** Warn at 90% of the model's `ContextWindow`; at 95%,
   declare and admit only `micro_handoff`. Compare total input from the
   latest completed main request, counting cached input exactly once.
   Exclude output and auxiliary compressor usage. After the handoff,
   restore normal tools. Tool restriction still permits text or refusal.
10. **Memory is data.** Band entries are in the conversation, not in
    the frozen system prompt. Changing a memory must not rewrite instructions.
11. **Fresh start and ordinary restart.** Scan enabled bands at running
    startup, before the first user request. Record each file's current bytes;
    an edited file replaces the old text at the same identity. Offline replay
    still reproduces the recorded version without consulting disk.
12. **Disable and enable.** Switching a band off records depopulation;
    its entries leave context but its files remain. Switching it back on
    rescans disk and records population per file. A human edit must reach
    the next request. A settings acknowledgment follows effective application,
    including while idle.

The deterministic forcing row remains `claude-ch16-course`, with a 4,096-token
window and `StubsToolResults` true; its original-numbered name is a fixture ID.
Unknown window sizes do not justify guessed thresholds. The compressor prompt,
internal corpus layout and settings appearance are yours; retain actual ownership
and the public submit contract.

```sh
go run ./cmd/grade-edition2 -ch 15
```

The default core is `solutions/ch15`; the existing `-gui` option selects
an alternate GUI. Section 15.11 maps the original eleven weights to revised
checks and their limits.

## §15.1 In Plain Words

Chapter 14 gives tool traffic a removal policy. This chapter gives completed
conversation a compression policy. The difference is what the replacement
tries to preserve: facts, decisions, reasons, failed approaches and unfinished work.

**Removal accepts the loss.** A tool result that becomes a stub
loses its content. The file can be re-read, the command re-run, so
the loss is recoverable. Removing conversation from the working set can hide irreplaceable meaning.
The retained event log may still contain it, but ordinary file tools cannot
reconstruct why the model rejected an approach from the source file alone. Removal is the
right tool for recoverable bytes. It is the wrong tool for
irreplaceable ones.

**Five bands of memory sit between the prefix and the
conversation.** They are ordered from what changes least to what
changes most: identity, curated knowledge, two compression tiers,
and the session band where the agent's own notes live. The ordering
matches the cache: stable bytes sit closest to the prefix, where
earlier bytes have a better chance of retaining their prefix match. Actual cache
grants still require their own observations.

**A checkpoint triggers compression.** When the agent calls
`micro_handoff`, chapter 14's mechanism strips every tool call and
tool result. This chapter measures what remains. If the conversation
exceeds its budget, the oldest finished work goes to a compressor
that returns a memory in the agent's own voice. The memory enters
the context as a band entry, and the conversation it replaced is
gone.

**When a band fills, its oldest memories graduate upward.** Session
memories fold into the 8x band, which folds into the 64x band,
which folds into MEMORY.md. Each graduation is two events: one
retiring entries from the source, one adding the compressed version
to the destination. The events carry the bytes.

**Every compression is a recorded event.** This is the invariant
that links chapters 14 and 15, stated here because it has a new
consequence. Chapter 14 records the cut selected under the settings effective then.
Changing those settings tomorrow must not change yesterday's replay. This
chapter also records a model's compressed text: nondeterministic, costly,
and potentially different on another call.
The event carries the bytes, and replay applies them without asking
a model. A compression that fires again returns a different summary.
A recorded event returns the same one.

Both chapters require recorded decisions. Recomputing even a deterministic
redaction during rendering would use today's setting and could change the
request being reconstructed. The compressor adds another reason to record the
answer: reproducing a past window must not buy a new summary and call it history.

## §15.2 One Door

CodeRhapsody, the agent that co-wrote this book, has two tools that
end a context. `save_memory` triggers the compression cascade.
`handoff_task` writes a structured document and starts a fresh
instance. Both clear the conversation. Only one feeds the cascade.
Bill first recorded the gap on 2026-04-18. As of this chapter's
design session, it had not been fixed.

The repair gives the cascade one trigger, reducing the number of clearing
paths that must agree. `micro_handoff` already strips tool calls
and results (chapter 14, rule 5). Measuring the conversation after
the strip is one function call at one call site. Neither `save_memory` nor `handoff_task` exists. One door in, one door out.

The acting model writes the checkpoint note because it has the current task
in front of it. That is the chance to state the reasons that later compression
must preserve. It is not a promise to capture hidden thinking. Tool clearing
removes calls and results; it does not categorically erase retained opaque
provider parts. Opaque replay material and a provider's private computation are
also different things. The compressor receives selected readable text, including
the note, and cannot reconstruct reasons that text never expressed.

The Agent runs this work after the successful tool batch and
before the next main request. The ordinary handoff result acknowledges the
checkpoint; the actor then asks its Memory child for a finite selection and
starts an Engine worker. That worker owns the HTTP call, while the actor remains
available for settings and interruption. Its returned text comes back through
the actor's ordered input path. A late reply from an interrupted turn is dropped.
The next main request waits for the compression outcome and any resulting
cascade; this is not a second independent conversation loop.

A checkpoint that sees the conversation under budget skips conversation
compression; a full memory band can still graduate. A checkpoint that sees the
conversation over budget selects the oldest whole segments
between checkpoints. Never half a segment. The grain of memory is
the grain of checkpointing, which is why the tool description asks
the model to checkpoint at completed sub-tasks, never at a token
threshold: a checkpoint taken mid-task cuts the task in half, and a
memory spanning half of two tasks is worse than a larger memory
covering one whole one.

## §15.3 The Five Bands

Five bands, five kinds. Each band has its own `EntryKind`, and
Chapter 14's rule applies to all of them: anything that must outlive
tool clearing is its own kind, and the kind determines which verb
removes it. The tool-clearing path removes tool parts from `KindDialogue`
entries and leaves the other kinds alone. A band entry is structurally safe.

| Band | File on disk | Kind | Cascades to |
|---|---|---|---|
| SOUL | `SOUL.md` | `KindSoul` | never |
| MEMORY | `MEMORY.md` | `KindMemory` | (terminal) |
| 64x | `bucket-1/*.md` | `Kind64x` | MEMORY.md |
| 8x | `bucket-0/*.md` | `Kind8x` | 64x |
| Session | `memory/*.md` | `KindSession` | 8x |

These paths are relative to the owning Agent's data directory. The supplied
application selects `cr`, so its session files are `cr/memory/*.md`, separate
from another Agent's corpus. Graduated source files receive `.md.retired`; they
remain on disk but no longer match the active `*.md` scan. Their identities still
reserve the next same-day session number after a fresh start.

The order is not cosmetic. The most stable band sits closest to the
prefix, where its unchanged bytes can remain eligible for prefix reuse. Identity changes
less often than curated knowledge, which changes less often than a
compression that ran today. Two independent arguments, one about
cache economics and one about reading order, give the same layout.

Band entries sit at `Seq: 0` rather than at the `Seq` of the event
that landed them. `Seq` exists so a redaction span can match an entry,
and no span ever reaches into a band. Giving a band entry the
landing event's `Seq` would put the order of restoration into the
context. Restoration order must not matter: switching every band off
and restoring them out of order must produce the same context.
Canonical placement by `(band, date, num)` makes the property true
by construction.

Memory file identifiers are dates and sequence numbers, not
event-log `Seq` values. A compressed file in `bucket-0/` is named
by the file range it spans:
`2026-09-20-1_2026-09-26-3.md`. The naming says "this file covers
the first memory of September 20 through the third memory of
September 26," and it says so without any event log in hand. The
memory corpus outlives any one log, and a range identifier tied to a
log-specific `Seq` would be the wrong coupling.

## §15.4 Two Events

Two event types handle all band changes, following chapter 14's
precedent of one event type and an enum field, rather than one event
per band.

`BandPopulated` adds one memory to a band. It carries the text
verbatim, because nothing downstream is allowed to look it up.
The first design specified file references. Its review found that neither the
pure reducer nor the renderer owned a place to resolve them. The file could
exist on disk while the model received no memory at all. The boundary is simple: Memory reads the file, the actor records its bytes,
and History applies those bytes without filesystem access.

A text part is sufficient. The media machinery from earlier chapters does not
make a mutable filename into a replayable memory. Tomorrow's edited file and
yesterday's delivered text are different facts, and the event must say which
text was delivered.

`BandDepopulated` retires entries from a band. For graduation, it
names a `Thru` file: everything at or before `Thru` is removed. For
disabling, it names nothing: the reducer wipes every entry of that
band's kind, computed at apply time. The asymmetry with
`BandPopulated` is the point. Removing need only say what goes.
Adding must say what arrives.

A graduation is one of each: a `BandDepopulated` retiring the
oldest source entries, then a `BandPopulated` landing the compressed
version in the destination band. Disabling is one `BandDepopulated`
alone. Re-enabling is one `BandPopulated` per file found on a
directory rescan. `CompactorLaunched` is the third event, and it
exists for one purpose: abandonment detection. It records that a
compressor started, so that a restart can see whether it finished.
The event is observable and never replayed as work.

## §15.5 The Compressor

The compressor is a direct LLM call: one prompt in, one response
back, no tools except `submit`. The prompt carries the stripped
conversation text and asks for a compressed memory at a target
ratio. Bounded input, bounded output. The accepted response is the memory; malformed or empty submissions leave the
source live and produce a diagnostic.

The call offers exactly one tool, `submit`, and the compressor
submits its result through it. The exact argument is `{"memory":"..."}`. That single-tool signature is what
identifies a compressor request on the wire.
The grader uses it to separate compressor traffic from turn traffic
without inspecting the student's function names or directory layout.

The compressor writes in the first person, as the agent. "I found
six dead settings fields" reads as memory. "The agent found six dead
settings fields" reads as a report about someone else. The distant
past is exactly the text a compressor writes, and if it reads as a
briefing about a stranger, the compaction introduced the
discontinuity the design exists to remove.

This is not a sub-agent. A sub-agent would need inherited settings,
a skill, tools, and a way to report back. The compressor needs none
of those: bounded input, one response, one tool. For this
chapter, the direct call is the right shape.

## §15.6 The Graduation Cascade

When a band crosses its high watermark, its oldest contents fold
into the band above. The high watermark defaults to twice the
budget; the low watermark defaults to the budget itself. The gap
between them is what keeps the cascade from firing on every
checkpoint: a band that just graduated is at its low watermark, and
it has a full budget's worth of room before the next graduation.

The cascade runs bottom-up. After a `BandPopulated{BandSession}`
lands, measure the session band. If it exceeds its high watermark,
select the oldest session files, enough to bring the band to its
low watermark, and compress them into a new file. On completion:
`BandDepopulated{BandSession}` retires the source entries,
`BandPopulated{Band8x}` lands the compressed version. Then measure
the 8x band. If it crossed its high watermark, compress its oldest
files into the 64x band. The chain continues to MEMORY.md.

One bug in this design waited through two drafts before the coder
caught it. The cascade originally required a full batch of eight
files before it would fold. A band over budget holding five large
memories would wait forever for a sixth that could only arrive by
growing further past its budget. The band that most needed folding
was the one that could never fold. That historical repair folded the available batch rather than waiting for eight.
This exercise selects the oldest files until the remaining source
reaches its low watermark, with a minimum of two. One oversized file still waits
for another; a band name is not a required count of eight files.

A failed compressor must leave its source available. The launch event records
the destination band and `Thru` identity. After valid text returns, Memory writes
the replacement file before the actor records retirement and population. History
holds that proposed retirement until a population with the matching band and
`Thru` arrives. A process interrupted between those two events therefore replays
the old source rather than a hole where a memory was supposed to be.

An unmatched launch at startup is abandoned. Replay does not retry it. The next
checkpoint measures again, and can propose a new fold if the band is still over
its watermark. Failed HTTP, malformed submit, interruption or a disabled band
can all prevent commitment. Old files remain retained outside the active scan
pattern after successful graduation. Recovery here covers process crashes, under
Chapter 14's persistence rules; it does not promise power-loss durability.

## §15.7 Forcing a Handoff

A nearly full context leaves little room to finish the current task or write
its checkpoint. Three approaches compete.

A stern prompt begging the agent to checkpoint relies on compliance.
The instruction competes with whatever the agent is doing. The framework
writing the handoff itself cannot supply reasons the model has never expressed.
The third approach narrows its executable choices: remove every tool except
`micro_handoff`. Warn at 90% of the known model window; at 95%, apply the restriction.

The numerator matters. Bill corrected the original chapter's word “cumulative”:
use total input tokens from the **latest completed main request**, including cache
buckets exactly once. Lifetime usage only rises. If it controlled forcing, one
large past turn could keep a newly compacted conversation restricted forever.
The main request measures occupancy; cumulative Engine usage measures work spent.

For the adapters in this snapshot, the raw fields normalize as follows:

| Provider response | Total input for this request |
| --- | --- |
| Anthropic | `input_tokens + cache_creation_input_tokens + cache_read_input_tokens` |
| OpenAI Chat Completions | `prompt_tokens`; cached tokens in `prompt_tokens_details` are already included |
| Gemini | `promptTokenCount`; `cachedContentTokenCount` is already included |

The normalized internal `Input`, `CacheWrite` and `CacheRead` buckets are disjoint,
so their sum recovers that total. Do not add an already-included raw cached subset
twice. Output, including visible or billed reasoning output, is not input occupancy.
The auxiliary compressor still contributes to Engine's cumulative usage, but its
small prompt must not overwrite the main conversation's last-input observation.

The course row makes the distinction easy to test: a 4,096-token window, one
main response above 95%, then a later main response of 512 input tokens. The
session total remains large; the later request is small. A compressor response
between them is a separate request with a separate purpose. Unknown windows stay
unknown rather than inheriting a guessed limit from another model.

The Agent records its warning as transient context and its declaration changes
through Chapter 14's existing path. The same restriction applies to dispatch,
so guessing a hidden tool name does not bypass it. After the handoff, normal tools
return; later checks use the next completed main request. This still permits a
text-only reply or a refusal. Removing tools cannot guarantee a handoff, a useful
note, or preservation of private thinking. The 95% threshold reserves headroom;
it is not a proof that every next response fits.

On a non-inline model, narrowing tools changes the top-level declaration prefix.
The historical design tested an Opus 5 model's recall of a chosen sentence after
a tool result was stubbed and after CodeRhapsody rewrote its system prompt on a
skill load. Bill compared the reply with an external record he could see. That
was stronger than accepting the model's assurance that it remembered. The notes
also acknowledge that the tool-array change was not isolated. Those observations
concern recalled text in those runs; they cannot establish how private reasoning
survives arbitrary prefix changes, nor certify today's model behavior.

## §15.8 Switching Bands Off and On

A band can be switched off with `"disabled": true` in its settings.
Disabling emits `BandDepopulated{Reason: "disabled"}`, which retires
every entry of that band's kind. The files on disk are untouched.
Re-enabling re-scans the directory and emits one `BandPopulated` per
file found.

The re-scan is the mechanism. A file edited on disk while its band was off is reflected on re-enable, because
the files are the memory, not the log's copy. The reference built
this property because a specific bug demanded it: the original
check matched on file identity alone, so after a restart the log
replayed what a file said when it was first read, and the disk copy
was decorative. Editing a memory to correct something the agent
believed changed nothing. The check now compares text as well as
identity. The property was found by accident: a test asserting that out-of-order restore produces
identical context found a bug that was not the property.

An ordinary restart must honor the same disk authority. Offline restoration
first reproduces the recorded text without reading files. Starting the Agent
then scans every enabled band and records current file contents, replacing the
same identity when its text changed. A startup test that first disables the band
can miss this: its entries are already gone, so even an identity-only comparison
appears to work. The second-edition implementation initially made that mistake too.
Editing an enabled file and inspecting the next request exposed it.

Canonical front placement also determines the cost of toggling. Re-enabling
MEMORY puts it before newer bands and conversation, not at the tail. Disabling
or re-enabling can therefore change a large request suffix while leaving the
frozen instruction prefix intact. Enabling memory is not always a cheap append; the actual
cache cost depends on the changed request and the provider's grants.

One interaction requires a guard. A band whose upward neighbor is
disabled cannot graduate. Without the guard, a session band over its
high watermark would attempt to fold into an 8x band that is
switched off, and the graduation would either drop the compressed
text or let the source band grow without bound. Neither is
acceptable. The Agent retains the source and logs a configuration warning. Disabling the
source or destination while a compressor is in flight also prevents that reply
from committing into the disabled band.

The settings acknowledgment created another trap in the student implementation.
The GUI held its observer mutex while waiting for the Agent to apply a memory
patch. During streaming, that actor was trying to publish an observation through
the same GUI mutex. The settings handler waited for the actor; the actor waited
for the settings handler. Even shutdown could not finish. Idle toggling passed,
which was why the first tests missed it.

The correction keeps settings transactions ordered without holding the progress
lock across the actor wait. A small separate settings mutex serializes patches
and subscriptions; the observer mutex protects only brief snapshot and publication
operations. A regression gates a real streamed delta, waits until the settings
file has been written, then releases the observation. That sequence reproduces
the old cycle and verifies the corrected acknowledgment and shutdown. Making the
acknowledgment arrive early would hide the deadlock by lying about application.

## §15.9 Memory Is Data

Band entries arrive in the conversation, as data, not in the system
prompt as instructions. The distinction is chapter 14's: `EntryKind`
records the origin of an entry, and the kind determines which
channel it speaks on.

There is a security-shaped reason too. When the agent instructs a
sub-agent, instructions flow down into a fresh context and output
returns as data. A handoff flows the opposite direction: forward to
the same agent, across a boundary that removes source context from its
working window. The tool results that would reveal where a piece of text
came from are stripped by the same operation that carried the text
forward. Compaction is the laundry. A memory arriving as
instruction (in the system prompt) would wear the agent's own
authority, while tracing its origin would require returning to the retained log.

A second reason to keep memory out of the system prompt: the system
prompt is part of the frozen prefix. Adding a memory to it means a
cache miss on every byte that follows, which is the entire
conversation. A memory in the dialogue leaves the frozen instruction untouched. It still
changes the request at its canonical band position, so cache effects require
actual request and provider observations.

The three settings categories this creates are worth naming.
`ContextTarget` is a policy setting: it steers a procedure that
records its own decisions, and a change applies to the next cut
without disturbing recorded ones. `Model` and `SystemPrompt` are
prefix settings: changing them needs a full refresh. Memory
enable/disable is a content setting: it changes what is in the
context, and without an event the log cannot reconstruct it. Each
category has a different cost, and only the third demands a new
event type.

Fresh starts use the same event. If memory files exist on disk and
their band is enabled, the Agent records `BandPopulated` for each at
startup, before the first user message. No special-casing: the
event is the same one a compressor emits, with a different `Source`
field for observability. An agent starting with a `SOUL.md` comes
up with its identity already in context. An agent starting with
`MEMORY.md` comes up with its curated knowledge. The mechanism is
the event; the file is where a human writes the content.

## §15.10 Measure It

The first edition reports a synthetic reference session from
`internal/llm/measure_test.go`: 40 turns, a checkpoint every 4
turns, 400 bytes of dialogue per turn, a conversation budget of
3,000 bytes, and a session band budget of 1,500 bytes.

| | Memory on | Memory off |
|---|---|---|
| Conversation | 1,600 bytes | 16,000 bytes |
| Session band | 27 bytes | n/a |
| 8x band | 22 bytes | n/a |
| **Total context** | **1,649 bytes** | **16,000 bytes** |
| Memories written | 4 | 0 |
| Cascade launches | 1 | 0 |

Cascade launches per checkpoint: 0.10. The watermark gap works:
nine of ten checkpoints find the band under budget and do nothing.

These are the reference test's reported structural numbers. The conversation stays at 1,600
bytes against 16,000 unbounded. Four memories were written. One fold
fired. The ladder from chapter 14 and the memory bands from this
chapter are additive: the ladder removes tool bytes, the bands
compress dialogue, and the combination holds the total context to a
fraction of its unmanaged size.

The band sizes are not a compression measurement, and printing them
as one would be wrong. The compressor in this test is a fake that
returns a short fixed string, so 27 bytes and 22 bytes are artifacts
of the fake's reply length, not of any real summarization. The 2x,
8x, and 64x compression targets in the design notes are inherited
from the shipped system; the retained study does not measure whether they are right.
Verifying them requires a real model, and nothing in this table
should be cited as evidence for them.

The live sessions below measure actual text reduction and expose what the short
fixed summaries could not test.

## §15.11 The Exercise, Graded

The grader reads what its fake vendor receives and what the agent
leaves on disk. It builds the optional GUI, applies memory settings through its
public settings channel, and waits for completed turns before checking requests
or closing stdin. A settings file in the wrong directory can leave a perfectly
working compressor idle throughout a test.

**micro-handoff-compresses** (15) scripts enough dialogue and
checkpoints to exceed the conversation budget. After the
compression, the next request contains a session-band entry with
the returned text, and the raw conversation it replaced is gone. The compressor
request must actually contain the selected work.

**planted-fact-survives** (10) plants a distinctive token in early
dialogue, pushes it through a checkpoint and compression. The token
must reach the compressor; the fake's returned memory must then reach the next
main request. This tests delivery through compression, not a real model's judgment
about which facts to preserve.

**graduation-fires-oldest-first** (12) scripts enough session
memories to exceed the session band's high watermark. An 8x entry
appears in the next request; the graduated oldest entries are gone,
and newer session work remains. The fixture distinguishes taking the oldest
files from taking the newest.

**disable-enable-idempotent** (12) switches the session band off,
checks that its entries leave the next request, then switches it
back on. The entries that reappear must match the entries before
disabling, by count and content. A second pass edits a memory file
on disk between disable and re-enable, and the edit must be
reflected. An ordinary enabled-band restart also rereads the edit. Checks inspect
actual memory order in the request, rather than sorting markers in the test.

**disabled-neighbor-refused** (8) disables the 8x band while the
session band is over its high watermark. No fold reaches the compressor,
and the source memories remain present exactly once.

**forced-handoff** (12) uses `claude-ch16-course`, a model with a
4,096-token context window. The grader drives the agent until
context approaches 90%, checks for a warning, then past 95%, and
checks that only `micro_handoff` remains in the tool list. Mixed cache buckets
count once. After the handoff and a smaller completed main response, normal tools
must be available even though lifetime usage has grown. A separate phase checks
that large auxiliary compressor usage does not restrict small main requests.

**replay-needs-no-llm** (10) starts from recorded memory populations and compares
a saved snapshot with full-log reconstruction. Startup alone makes zero HTTP
requests; identical subsequent prompts produce byte-identical vendor requests
without a compressor call. Running startup also rescans disk. Pure offline
restoration of populations and retirements, without disk reads, remains a separate
source regression.

**abandon-on-restart** (8) first populates a real source band, then exercises
a failed fold and a genuinely pending fold interrupted by process kill. Restart
retains the source without immediately retrying compression; a later checkpoint
can try again. An empty band with zero attempts cannot pass as failure recovery.

**memory-is-data** (5) checks that the returned memory appears in the
conversation rather than the system prompt. The fake supplies first-person text;
review of the compressor instructions establishes that the real model is asked
to write that way. The numeric check does not grade resistance to malicious text.

**fresh-start-populates** (5) places memory files on disk before
the agent starts. The first request sent to the model contains
their content.

**ch14-parity** (3) runs the chapter 14 grader against the same
solution and checks it still scores full marks.

These eleven weights total 100. Low-watermark arithmetic, the required disabled-
neighbor diagnostic, all-provider normalization, whole-segment reducer cases,
native browser settings and the GUI/actor interleaving retain their focused source
and browser checks. The original grader scored **20/100 FAIL** on an
intermediate student build: its settings missed the owning Agent's directory,
and its awarded recovery points did not require a fold to occur. The
[student record](../docs/edition-2-notes/ch16.md) preserves that result and the
independent review; revised grading is a separate measurement.

## §15.12 Taking It for a Spin

Use a scratch copy of your Chapter 15 core. Build its optional GUI, whose terminal
and browser share the same Agent:

```sh
(cd gui && go build -o /tmp/ensemble-ch15-gui ./cmd)
mkdir -p cr skills/base
```

Write `skills/base/SKILL.md`:

```yaml
---
name: base
type: primary
description: Analyze supplied cases and retain useful checkpoint memories
tools: read_file write_file list_directory
---
Analyze the supplied case. Preserve facts, decisions, reasons, failed approaches,
user constraints and open questions in checkpoint notes. Do only the current task.
```

The existing lifecycle and context tools include `micro_handoff`. Keeping file
tools available lets the exercise distinguish obeying a no-tools instruction from
having no tools to call. Write `cr/settings.json` with deliberately small budgets:

```json
{
  "tts_enabled": false,
  "max_tool_rounds": 3,
  "memory": {
    "conversation": {"budget": 600},
    "session": {"budget": 500, "high_watermark": 1000, "low_watermark": 500}
  }
}
```

With `LLM_API_KEY` set in the environment, launch from that core directory:

```sh
LLM_VENDOR=anthropic LLM_MODEL=claude-sonnet-5 \
EN_PRIMARY_SKILL=base EN_SKILLS_DIR="$PWD/skills" \
/tmp/ensemble-ch15-gui --port 8088 --gui-dir "$PWD/gui/web/gui"
```

Use the terminal or open `http://localhost:8088`. Send each of the following as
one message, waiting for the turn to finish before sending the next. Input while
busy is a hint. The cases are the actual inputs from the recorded memory trials;
the explicit skill and supported Sonnet-5 selection make a reproducible setup for
this edition, not a claim of a new live result. This snapshot has no declared
context-window size for Sonnet-5, so the run exercises compression rather than
the forcing threshold. Use the course row with the deterministic grader for that.

The first completed analysis supplies an archive decision and its reason:

```text
This is a bounded memory exercise in an empty scratch workspace. Do not use filesystem or command tools. Analyze the supplied case, explain the decision briefly, then call micro_handoff exactly once with the important facts and reasons. After the tool returns, finish this turn without another tool call. Case one: Our project is Cedar Lantern, release code CL-721. We rejected automatic deletion of imported notes because the support team needs the originals when customers report mistakes. The replacement design keeps source notes in an archive and stores the accepted summary separately. Morgan owns the archive review, scheduled for October 18. Acceptance requires preserving a source identifier with each summary. The unresolved question is whether attachments belong in the same archive or a separate one; do not claim that has been decided. The default retention is thirty days. A failed summarization must leave the original note available. The initial approach coupled cleanup to request arrival, so an unrelated request could delete unreviewed work. We moved cleanup to an explicit completed-review boundary instead. Explain why the latter boundary matters, and preserve these facts in your checkpoint.
```

The second supplies a different failure and another unresolved question:

```text
This is a bounded memory exercise in an empty scratch workspace. Do not use filesystem or command tools. Analyze the supplied case, explain the decision briefly, then call micro_handoff exactly once with the important facts and reasons. After the tool returns, finish this turn without another tool call. Case two: The deployment rehearsal uses Harbor Finch, diagnostic tag HF-934. Priya owns the rehearsal scheduled for October 21. We rejected retrying a database migration automatically after a timeout: the remote operation may already have committed. The chosen procedure reads the migration ledger before deciding whether any retry is safe. The operator must retain the exact migration identifier and measured error, not a guess about success. We have not settled whether the ledger should live alongside application data or in a separate control database. The rollback package must remain available for seven days after deployment. Last week a log line said the worker stopped, but the child process was still alive and held the port open. The rehearsal now checks the child lifetime directly. Explain how this evidence changes the decision and checkpoint the completed analysis, preserving the unresolved question as unresolved.
```

Then ask for recall without tools:

```text
Without using tools, recall both project names, diagnostic codes, owners, dates, rejected approaches and reasons. Explicitly identify the two unresolved questions. Do not invent a resolved answer.
```

Inspect `event.log`, `api.log` and the files under `cr/`. After a handoff,
the compressor should receive a whole completed segment, including its note.
A successful response lands a recorded cut and a session memory. When two session
memories cross the watermark, their oldest work graduates into `cr/bucket-0/`.
Model output lengths vary; the two turns need not cross the same thresholds on
every run. Inspect the actual calls and files instead of inferring a fold from a
short answer. The request layout after graduation is:

```text
[frozen instruction and tool prefix]
[enabled SOUL and MEMORY, if present]
[64x memories, then 8x memories, then session memories]
[remaining conversation]
```

The restart is the interesting part. End input with Ctrl-D after the completed
turn, then run the same command again. Restoration rebuilds the recorded bands
without a compressor call; running startup rereads enabled files before the next
user request. The model receives the resulting memories when you ask it to continue.

Switch a populated band off in Settings. Its entries leave the context. Edit an
active memory file on disk: correct a name, add a detail, remove a paragraph.
Switch the band back on. The entries return, and the edited file carries its new
text. Try the ordinary restart too, leaving the band enabled while editing its
file. The files are the memory. The log's copies record what was delivered, so
reconstructing yesterday's request and reading today's corrected memory can both
be right.

### What the live runs remembered

Five recorded sessions on October 11 each finished three turns, exited normally,
and made three compressor calls: two session memories and one graduation to 8x.
They also show why watching the files appear is only the first test.

The first Haiku 5.5 run preserved the decisions to reject automatic deletion and
automatic retries. Graduation discarded the reasons. Asked why those approaches
were rejected, the model answered, “No reason was given.” The reasons were in the
user messages and in both handoffs; they were absent from the graduated memory.
A smaller request had thrown away exactly what would stop the next agent from
repeating the mistake.

The compressor prompt was corrected to preserve decisions **and why**, failed
approaches, uncertainty and still-applicable user constraints. Its size ratio
became a guide: use more space rather than discard the reason. The bounded repeat
kept those facts and reasons, but also carried forward a model-invented question
about retention expiring before review. The final answer reported that question
alongside the two the user had actually left open.

| Recorded run | Preservation and task result |
| --- | --- |
| Haiku 5.5, initial | Reasons lost in graduation; final answer says none were given and misreads the supplied diagnostic identifiers. |
| Haiku 5.5, corrected prompt | Reasons retained; an invented retention-versus-review question enters memory and final recall. |
| GPT-4o mini, initial | Reasons and facts retained; two file writes violate the user's instruction. One written artifact turns an open attachment question into a decision. |
| GPT-4o mini, corrected prompt | No extra tools, but final recall omits Morgan and substitutes October 11 for the October 18 review date. Both correct facts were in the final request. |
| Gemini 3.8 Flash | Recalls the supplied facts, reasons and two open questions; still calls `list_directory` despite the no-extra-tools instruction. |

The Haiku identifiers describe historical runs. Bill subsequently excluded 5.5
models from current use because of his experience with ephemeral-context
compatibility. These transcripts establish the failures shown, without diagnosing
that compatibility issue or predicting how the supported Sonnet-5 will perform.

The corrected OpenAI run is the harder failure to detect. Its graduated file
says, “Morgan will lead the archive review on October 18.” That same text appears
in the actual last provider request, followed by an assistant message repeating
the owner and date. The final answer still says “Owner: Not specified” and gives
the memory's creation date. Inspecting only the file would certify preservation
and miss failed recall. Inspecting only the answer might blame compression for
facts it had retained.

Measured source-to-memory UTF-8 byte ratios across these calls range from 1.358
to 3.678. The 8x name describes a tier's intended compression, not an achieved
ratio. Preserving reasons used more bytes. The [run records and actual memories](../docs/edition-2-notes/ch16.md#live-mechanics-versus-preservation-and-task-quality)
retain all five outcomes.

Test the three steps separately: did the framework replace and replay the right
work, did the memory preserve what mattered, and did the model use it correctly
while obeying the task? A memory file that passes the first test can still fail
the other two. Bill's workplace experiment began with memory that did not run at
all. Once it runs, the next question is what the agent actually learned.
