# Chapter 10: Persistence

Every agent in this book so far is a goldfish. It reasons, calls tools,
streams into its GUI and loads skills on demand, and the moment the
process exits it forgets the user's name, the project, the afternoon of
work and every decision made along the way. The next start meets a
stranger. This chapter ends that. Nothing else in the book changes more
about what the agent is: a program that runs becomes one that
remembers.

## TL;DR

The agent saves itself on exit and loads itself on start. Neither needs
a flag. The save file is one JSON object: the configuration that shaped
the wire, the Chapter 2 context as a snapshot, the Seq that snapshot was taken
at, and the retained history. Full history is the default; an explicit
dialogue-only mode keeps a derived conversation record instead of the full
canonical archive. A load installs the snapshot and replays only the
events after its anchor. A snapshot with no log is a complete save. So
is a complete log with no snapshot.

```go
// SaveFile joins the conversation snapshot to the available recorded history.
type SaveFile struct {
    // Config records request settings without restoring them later.
    Config SaveConfig `json:"config"`
    // AsOf is the last event folded into Context, not the log length.
    AsOf uint64 `json:"as_of"`
    // Context is the Chapter 2 projection; null requests a full rebuild.
    Context *Context `json:"context"`
    // Log contains canonical events: full-mode history or the recovery tail.
    Log []Event `json:"log"`
    // LogMode is "full" (also the meaning of an omitted value) or "dialogue_only".
    LogMode string `json:"log_mode,omitempty"`
    // DialogueLog is a derived text archive, never input to the event reducer.
    DialogueLog []DialogueRecord `json:"dialogue_log,omitempty"`
}

// DialogueRecord projects visible dialogue without claiming to be an Event.
type DialogueRecord struct {
    SourceSeq uint64 `json:"source_seq"` // provenance, not a newly allocated Seq
    Role string `json:"role"`           // user or assistant
    Text string `json:"text"`           // visible dialogue only
}

// SaveConfig records what shaped the wire. Never the API key.
type SaveConfig struct {
    // Model identifies the configured request model.
    Model string `json:"model"`
    // Vendor uses its existing disk spelling: anthropic, gemini or openai.
    Vendor Vendor `json:"vendor"`
    // SystemPrompt is the effective instruction used by the renderer.
    SystemPrompt string `json:"system_prompt"`
    // Tools contains declarations, never executable handlers.
    Tools []ToolDeclaration `json:"tools"`
}

// ToolDeclaration is the existing tool vocabulary with explicit disk names.
type ToolDeclaration struct {
    // Name identifies the callable tool.
    Name string `json:"name"`
    // Description tells the model when to use it.
    Description string `json:"description"`
    // Schema describes its JSON arguments.
    Schema json.RawMessage `json:"schema"`
}
```

The original four fields follow the accepted implementation's common vocabulary;
`log_mode` and `dialogue_log` extend that format for the retention exercise.
Preserve the JSON fields without requiring these exact internal type names.
Ordinary full-mode saves may omit both new fields.

1. **Default location.** The single-Agent commands select `save.json` in their
   working directory. The library defaults to `save.json` under the owning
   Agent's data directory. `--save PATH` names a different
   file. The flag changes where, never whether: the same path is loaded
   at start and written at exit.
2. **Load at start.** Validate the selected primary completely. If it is
   missing or invalid, validate its `.bak` sibling with the same rules
   (`save.json.bak` by default). Start fresh only when both are absent.
   If neither existing candidate validates, exit non-zero, name the failures,
   install no Agent and leave both files untouched. Report backup recovery
   and its saved boundary rather than implying that newer lost events survived.
3. **Snapshot plus tail.** If `context` is non-null, install it, then
   apply in order every log event with `seq > as_of`. If `context` is
   null, require a complete canonical log and apply it to a fresh context.
   A dialogue-only save requires its snapshot. Events at or below
   `as_of` are already inside the snapshot; applying one twice is a bug.
4. **The log is not needed.** `"log": []` with a non-null `context` is a
   complete save. The vendor sees the context, never the log.
5. **Numbering continues.** The first new event gets the Seq one past
   the larger of `as_of` and the last log event's Seq. GUI artifact identities
   must also survive restart: replay an old response, then stream a new one
   and verify that the new text cannot append to or replace the old card.
   Already-colliding archived IDs must also produce distinct presentations,
   deterministically across replay and restart, without rewriting archive
   facts. Preventing new collisions alone does not repair that case.
6. **Save at exit.** When stdin closes, finish admitted work, shut down the
   Agent and join its managed jobs, then write the file and exit. `as_of`
   is the Seq of the last event folded into the saved context. The
   saved log is the loaded log plus every event created since, oldest
   first, Seq strictly increasing. `log_mode` defaults to `full`, retaining
   all available canonical history. `dialogue_only` retains a derived dialogue
   archive instead, with the snapshot and complete recovery tail kept
   separately (§10.5).
7. **Config is a record, not a restore.** On load the running agent's
   own model, vendor, prompt and tools win. The context is
   vendor-independent (ch2); the save must not pin a model.
8. **Replay is deterministic.** For a save whose log is complete, loading it
   as written and loading it with `context` set to null must produce byte-identical vendor
   requests for the same next prompt under the same runtime configuration.
   A trimmed log (rule 4) cannot rebuild the missing history, so the rule cannot
   apply to it, including a later save descended from that trimmed input.
   Reject and diagnose attempted full-history replay when the originals are
   unavailable; a dialogue projection is not an executable event stream.
9. **Diagnostic retention is separate.** `diagnostic_max_records` defaults
   to 200 complete records per managed diagnostic log at each retention
   boundary. Zero means unlimited; negative values are invalid. This policy
   never trims the canonical saved log, recovery journals or Agent state.
10. **Trim at safe boundaries.** Apply diagnostic retention at startup after
    successful recovery, and after a successful normal-exit save. Establish
    a valid recovery save before destructive cleanup; failed load or save
    leaves diagnostics intact. Coordinate trimming with the owning writer
    (§10.5). Chapter 14 adds the successful `micro_handoff` save boundary.

11. **Keep a validated backup.** Fully encode and validate the new save before
    changing either file. Before atomically replacing a valid primary, retain
    that validated primary as its `.bak` sibling. Never rotate corrupt primary
    bytes over a good backup. After backup recovery, preserve that backup until
    the new primary is installed. Serialize saves for the same Agent; a failed
    encode, validation or write must leave a valid recovery candidate intact.
    Use temporary files and rename for replacement, without claiming universal
    power-loss durability or introducing Chapter 14's journal early.

Yours: indentation, a `verify` subcommand for debugging, and the wording of
load and recovery notices. A backup restores its own boundary, not newer bytes
lost with a damaged primary.

**Exercise.** Extend your Chapter 9 agent with this save/load lifecycle in
both CLI and GUI. Keep one public Agent implementation behind them.

Build the executable package for your module layout. A library root
can be imported by the grader without being a runnable command. The fresh
student implementation keeps its CLI in `cmd/ch01`, for example:

```sh
cd /absolute/path/to/your/core-module
go build -o /tmp/ensemble-ch10 ./cmd/ch01
```

Use your own CLI package path. Use `go build -o /tmp/ensemble-ch10 .` only if
that directory actually contains your command's `main` package.

From the **book repository root**, run the current Edition 3 command with all
five module roots and both declared model IDs. Replace the placeholders with
your implementation's paths and capability rows:

```sh
go run ./cmd/grade-edition3 -ch 10 -json \
  -core /absolute/path/to/core-module \
  -consumer /absolute/path/to/consumer-module \
  -consumer-case /absolute/path/to/consumer-case.json \
  -workflow /absolute/path/to/workflow-module \
  -media-models /absolute/path/to/media-map.json \
  -stream-exercise /absolute/path/to/stream-exercise-module \
  -stream-model YOUR_DECLARED_STREAMING_MODEL \
  -gui /absolute/path/to/gui-module \
  -gui-model YOUR_DECLARED_GUI_MODEL \
  /tmp/ensemble-ch10

go run ./cmd/grade-edition3 -ch 1 -json /tmp/ensemble-ch10
```

The consumer-case and media-model map paths are optional inputs inherited from
the earlier chapters; remove those two flags if you do not supply them. The
media map selects your declared `anthropic`, `openai`, and `gemini` model rows.
It does not establish live-provider capability. All flags precede the explicit
CLI executable; the core module and executable are separate inputs.

Both reports must have `passed: true` and exit zero on the **same frozen
binary**. Chapter 1 still runs the real 62-second delayed-header and delayed-body
phases. Allow both to finish: a 90-second outer timeout cannot cover them, and
a 100/100 score can still fail a required zero-point check.

The GUI stays a separate module; the grader uses the appropriate command for
each phase rather than applying headless checks to a GUI server.

| Check | Points | Proves |
|---|---|---|
| save-shape | 15 | recorded effective config and declaration fields, correct anchor and increasing log sequences |
| default-load | 20 | restart retains conversation; default and override use the same path for load and save; runtime configuration wins |
| replay-equals-snapshot | 20 | rule 8 |
| tail-applied-once | 15 | turn-2 snapshot spliced onto the turn-3 log yields the same next request as the turn-3 save |
| log-not-needed | 10 | rules 4 and 5 |
| bad-save-refused | 5 | validated backup fallback; refusal without destructive writes when neither candidate is valid |
| ch9-parity | 15 | Chapter 9 still passes through the appropriate clients |

The persistence phases exercise both actual commands through their existing
stdin input and EOF lifecycle. In the GUI command, redirected prompt lines
use the same `{"kind":"prompt","text":"remember this"}` envelope as the
headless protocol; ordinary terminal text still works. Closing that stdin now
flushes admitted work and ends the application, rather than merely detaching
the terminal view. Browser prompts continue to use the same Agent.

Keep source tests for independent Agent directories, rich neutral context
and intentionally empty primary instructions. The numeric grader's chosen
fixtures cannot prove every ownership or renderer property.

## §10.1 In Plain Words

The event log is the truth. The context is what the truth means right
now. Chapter 2 split them on purpose: events are appended and never
edited, and the context is whatever the reducer makes of them. The
split has carried streaming, jobs, artifacts and skills. Persistence now
tests that relationship across a process boundary.

The context must be a pure function of the events.
No clock, no map iteration order, no field set by the engine behind the
reducer's back. The earlier replay checks depend on that rule too.
An agent that runs start to finish in one process can violate it
forever and never notice, because the live context is the only copy
anyone ever looks at.

Saving creates a second copy. Once a snapshot sits on disk next to the
log that produced it, the claim becomes falsifiable: rebuild from the
log and compare the rendered requests under the same configuration. A difference
exposes something the restore, reducer or renderer failed to preserve.
Persistence is the feature a user sees. Verification is the reason it belongs
this early in the book.

## §10.2 The Save File

The recovery fields remain `config`, `as_of`, `context`, `log`. Optional
`log_mode` and `dialogue_log` identify a derived archive (§10.5); old four-field
saves mean full mode. Use the existing Context and Event types for recovery.
The dialogue projection is not a second state model and never drives replay.

`Context` is the snapshot: exactly what the renderer reads. `Log` holds
available canonical events, not the derived dialogue archive. `AsOf` joins
the snapshot to those events. It is the Seq of the last event
already folded into the snapshot, and it is the field everything else
in the chapter hangs on. Without it, a loader holding a snapshot and a
log cannot tell which events the snapshot already contains, so it has
two choices and both are wrong: apply the whole log and duplicate every
turn, or apply none of it and lose whatever came after the snapshot.

`Config` records vendor, model, system prompt, and tool declarations.
It is a record only. The running agent's own configuration
wins on load, so a save made with one model resumes under whatever
model the user starts with today. A save that pinned its model would
turn every model retirement into a pile of unloadable files. API keys
and base URLs are excluded from this metadata, so the save is portable across
machines and endpoints without carrying transport credentials.

There is a small distinction between configured and effective instructions.
The earlier Edition 2 implementation's first save code copied
`Config.SystemPrompt`. An ordinary Agent without a skill left that field empty,
while Engine supplied a nonempty default
to the actual request. The conversation resumed correctly and the metadata lied.
The grader's supplied primary skill hid the problem because that path had
already filled in the configured text.

That correction asks Engine for the effective instruction through its existing
owner interface. Saving and all three renderers now use the same decision.
An intentionally empty primary body must remain empty; blindly filling every
empty string with the default would fix one case by breaking the other.
The regression compares saved `system_prompt` with the captured HTTP field,
then checks an empty primary separately. Restoring either defect fails its
corresponding assertion. No new paid run was needed to verify this correction.

The save does not reconstruct a running job, a handler closure or the skills
registry's loaded/pending state. The new Agent supplies its current tools and
primary skill; its restored conversation may contain earlier tool results and
manuals. Serialized declarations describe what was available when saved. They
do not install executable capability on restart.

## §10.3 Snapshot Plus Tail

Loading has one loop. This excerpt belongs in the history implementation,
where `validate` and `apply` are the existing event validator and reducer;
`SaveFile` itself stays shared vocabulary in common:

```go
func restoreContext(sf common.SaveFile) (*common.Context, error) {
    ctx, after := sf.Context, sf.AsOf
    if ctx == nil {
        ctx, after = &common.Context{Turn: "idle"}, 0
    }
    for _, e := range sf.Log {
        if err := validate(e); err != nil {
            return nil, fmt.Errorf("restore: event %d: %w", e.Seq, err)
        }
        if e.Seq <= after {
            continue
        }
        apply(ctx, e)
    }
    return ctx, nil
}
```

With a snapshot, install it and apply only the tail: events whose Seq
is strictly greater than `AsOf`. Without one, start from an empty
context and apply everything. The same loop serves both, which leaves
one replay path to get right instead of two.

Decode into a fresh save value and validate its configuration, snapshot state
and increasing log sequence numbers before installing the restored history.
The excerpt shows the fold, not that surrounding validation and installation.
That validation must also reject a dialogue-only save without its snapshot
and any attempted rebuild without complete canonical history. Never feed
`dialogue_log` records into this loop.
An invalid later event must not leave the live Agent half-restored. The accepted
loader stages the candidate context and log, then replaces the owned state only
after all checks succeed.

Validation must accept the events the public API legitimately produces.
`Tools().Dispatch` can record tool events outside a model turn without a native
call ID. Preserve that absence rather than inventing historical IDs; model-turn
tool calls still require their supplied correlation IDs.

The tail looks unnecessary. A normal save from a complete history has
`AsOf` equal to the last Seq in its log, so the tail is empty and a
loader that skips it passes every test built from its own output. The
tail matters the moment a snapshot and a log come from different
moments: a snapshot kept from an earlier save, with a newer log written
after it. That is the shape a crash leaves behind, once the log is
written event by event and the snapshot only now and then. The grader
builds that shape on purpose, splicing an old snapshot onto a newer
log to distinguish a correct loader from one that ignores its anchor.
This chapter still saves on normal exit. Journaling events during execution
and recovering from a process crash arrive in Chapter 14; the tail fixture
does not claim that machinery exists yet.

The `<=` is the other half. An event at or below the anchor is already
inside the snapshot. Applying it again gives the conversation a
duplicate turn, and the model answers a question it has already
answered.

Take the snapshot and its anchor together. If one is copied before an append
and the other afterward, the file can claim to contain a fact its context has
never seen. History already owns those values; its existing lock can protect
the pair without a new persistence coordinator. Load before submitting work;
save after shutdown has joined the jobs that can still append final facts.

## §10.4 Replay Equals Snapshot

The chapter's central claim fits in one sentence. For a save whose log
is complete, loading it as written and loading it with `context` set
to null must produce byte-identical vendor requests for the next
prompt under the same runtime configuration.

The first load takes the snapshot path. The second takes the rebuild
path, replaying every event from nothing. If the reducer is a pure
function of the events, the two paths land on the same context and the
renderer turns that context into the same bytes. If they differ, inspect the
restoration and reduction first, then the renderer's ordering and configuration.
Both paths must use the same model, declarations and settings for the comparison
to isolate persistence.

The grader compares vendor requests.
A request is the only thing the model ever sees, and it is a format
every student's agent already produces, so the grader can compare two
of them without knowing anything about how a particular solution
stores its context. Two contexts that differ in some field the renderer
never reads are the same conversation. Two requests that differ by one
byte are not.

Byte-identical requests do not require byte-identical save files. Indentation
and object-key formatting on disk are yours; the renderer already determines
the bytes compared here. When validating a replay against its snapshot, compare
decoded values: formatting inside a JSON tool argument does not change its
meaning. Normalize nil and empty in-memory `Dialogue` or `Ephemera` collections
for that comparison. Preserve exact numbers; rounding a large integer could
conceal a changed argument. This tolerance for representation does not relax
required snapshot fields, snapshot facts, sequence ordering or anchor checks.
A canonical disk format, checksums or an origin certificate would solve a
different problem and add nothing to this test.

The precondition is real. A save with a trimmed log (§10.5) lacks the
complete history needed for a rebuild. An empty log would leave an empty
agent; a retained suffix cannot reconstruct the missing prefix. Refuse that
rebuild rather than accepting either result as recovered history.
A run that starts fresh and retains its log can write a complete history.
After loading a deliberately trimmed save, however, the agent can retain only
the events it was given and the ones it creates next. A new save cannot recover
discarded facts. Keep the snapshot in that case; do not call nulling it a
replay-equivalence test.

## §10.5 The LLM Never Sees the Log

The renderer reads the context. The context holds dialogue entries,
the entries hold parts, and at no point does the renderer consult the
log. A save with `"log": []` and a non-null context is therefore
complete: it loads, it resumes, and the next request carries the full
conversation.

This permits trimming already-folded log events without changing restored
behavior. It also creates one small trap. Numbering must continue after the
loaded events, and a save with an empty log knows its anchor and
nothing else. The next Seq is one past the larger of `AsOf` and the
last event in the log; neither alone is enough.

For example, a snapshot anchored at 7 with an empty log must assign 8 to the
next event. The resulting save has anchor 8 and one retained event, not an
eight-event log. Sequence identity and array length measure different things.

The browser has a related continuity problem. A provider's tool-call ID and a
transcript card's address serve different owners. Reusing a parser's turn-local
number as a global card address can make an old result select a new card after
restart. The fresh student implementation had to preserve old card addresses
while allocating new ones beyond the restored conversation.

That includes loading a legacy archive whose parser-local IDs already collide.
Repairing the display must leave provider payloads and call/result correlations
intact. Renumbering those facts to make the DOM convenient would change the
conversation being replayed. Old cards must keep their addresses and new cards
must get unused ones; that constraint does not prescribe a new public identity API.

### Diagnostic retention and saved-history modes

Recovery history and diagnostic output have different jobs. The canonical
save preserves conversation; API and GUI traces explain operation. Diagnostic
files can exhaust a disk without adding recoverable conversation. Bill reports
that his Gemini CLI fork also ran out of disk space; its storage internals and
whether it retained context snapshots remain unverified here.

Expose `log_mode` and `diagnostic_max_records` through the
existing ownership and configuration interfaces. The Agent owns saved-history
policy; the diagnostic writer owns cleanup. Use Chapter 8's settings persistence,
not provider secrets, browser display preferences or extra `SaveConfig` fields.
These names do not prescribe a configuration-file location or nesting. Omitted
values select defaults; an explicit diagnostic zero must survive updates and
restart. Reject negative diagnostic limits and unknown modes. The save's
`log_mode` describes its representation; runtime settings govern the next save.

In `full` mode, preserve all available canonical history. Loading a previously
trimmed save cannot restore its missing prefix.

In `dialogue_only` mode, the compacted archive is `dialogue_log`: ordered,
derived user/assistant text records with source sequence provenance. Retain
visible dialogue, excluding tool calls/results, opaque reasoning or private
replay payloads, and API/GUI traces. For a mixed assistant event, project only
its dialogue text; do not disguise the filtered record as an unchanged Event
or fabricate executed events. Preserve available prior dialogue records when
saving again. Arbitrary last-N truncation is not this dialogue projection.
Switching back to full mode cannot recover discarded originals.

The snapshot remains authoritative and may retain current state needed for
recovery, including tool or provider material. Dialogue-only mode requires a
non-null `context`; it does not compact that context. Keep its canonical
post-`as_of` recovery tail complete in `log`, separate from `dialogue_log`,
until a newer snapshot folds those events. Chapter 14's journal likewise
stays separate and must not be filtered. Apply only canonical events above
the anchor, once; preserve their Seq values and continue allocation past the
larger of the anchor and last event. Build the replacement save before
replacing the old one; a failed save must not discard the prior recovery
record or trim live history.

An available checkpoint supports diagnosis of the state represented there;
its complete tail supports replay from that boundary. Dialogue helps explain
the conversation, not reconstruct discarded steps, context windows or raw tool
results. `as_of` identifies the last event folded into the snapshot, not an
oracle for missing facts. Reject attempted full-history replay when originals
are unavailable, rather than treating filtered dialogue or a bounded suffix
as a complete canonical stream. Complete-log replay equivalence stays unchanged.

Diagnostic retention counts complete records, not lines or bytes. Keep the
newest configured number in order; zero disables trimming. At startup, finish
load/recovery before cleanup; if no valid recovery save exists, successfully
establish one first. At normal exit, finish admitted work and save successfully
before trimming diagnostics. Load or save failure leaves diagnostics intact.
These are retention boundaries, not a continuous disk-size cap.

Serialize cleanup with the writer's appends, retaining whole records and
resetting its offset or reopening the replacement before resuming writes.
External truncation can leave a non-append descriptor at an old offset and
create a sparse gap; unlinking or renaming its pathname does not close the
old inode or guarantee reclaimed space. Report cleanup failure without
destroying the prior diagnostic file or invalidating the recovery save.
Never clean up `save.json`, recovery journals, memories or worker state as
diagnostics; shared rebuildable caches require separate authorization.

Chapter 14 must add the same diagnostic-retention operation after a successful
`micro_handoff` save, preserving its snapshot/journal save-then-trim ordering,
and test that integration. Align that chapter's existing `log_retention` with
these modes. This chapter does not require compaction machinery or a
`micro_handoff` tool.

## §10.6 Default Load, and Refusing a Bad File

The first design had a `--load` flag. Bill's review of it was one
line: "Let's load by default without a flag." The agent now loads
`./save.json` at startup if it exists and writes it when stdin closes.
`--save PATH` names a different file for both directions and leaves
loading on.
That is the command's choice, so a project directory remembers its own
conversation. The library itself resolves its default under the owning Agent's
data directory; two Agents with different directories must not share one
process-wide save. A relative `DataDir` follows the owning `BaseDir`, rather
than the launch directory. Resolve that derived default without rebasing an
explicit `SavePath`, whose location remains the caller's choice. In the
current GUI command, preferences remain at `cr/settings.json`; the save does
not need to move beside them or take over their persistence.

Loading by default makes one failure mode dangerous. If `save.json`
cannot be validated, starting fresh looks friendlier and is worse: the fresh
session saves on exit and overwrites the history the user came back for.
Try `save.json.bak` instead, including when the primary is missing. Validate
that candidate independently before installing any state. Only two absent
files mean a fresh conversation; if neither existing candidate validates,
refuse startup and preserve both files. A refusal costs one confusing startup.
A silent reset costs the conversation.

A recovered backup is an older valid snapshot, not evidence that nothing was
lost. Report its anchor. Keep Context, `as_of` and Log from the selected
candidate coherent; do not splice in an arbitrary tail from a corrupt primary.
This chapter's log lives inside the save, so newer bytes lost with that file
cannot be reconstructed. Chapter 14 adds a separate journal retaining complete
events after the **backup** anchor. That later mechanism is not part of this
exercise.

EOF ends input admission, but accepted work may still be running. Wait for its
reliable completion before shutdown; a quiet observer stream does not prove
that the final turn has finished. The CLI emits final usage before its deferred
shutdown saves. Receiving that usage record is not a save acknowledgment:
wait for successful process exit before inspecting or reopening the file.

Install save-on-exit only after load succeeds. An unconditional deferred save
placed above `Load` would overwrite the bad file while returning the right
error. Check the bytes after refusal and verify that no HTTP request was sent;
an error message alone leaves the destructive half of that bug invisible.

For a successful exit, encode and validate the complete new save before
changing primary or backup. If the existing primary validates, install its
bytes as the backup using a temporary file and rename before replacing the
primary the same way. If recovery selected the backup because the primary was
missing or corrupt, leave the good backup alone while installing the new
primary. Never copy corruption over the recovery record.

Keep this sequence with the existing save owner, serialized for one Agent.
Encoding, close and replacement failures propagate to the caller; failure must
leave a valid recovery candidate intact and must not trigger diagnostic
trimming. Atomic replacement prevents partial JSON installation, not every
filesystem or power-loss failure. A database, checksum framework or write-ahead
journal is unnecessary for this normal-exit exercise.

## §10.7 Taking It for a Spin

Use a fresh scratch directory so you cannot accidentally load an earlier
conversation. The following launch matches the fresh implementation's
`cmd/ch01` command built above. It chooses the model through `LLM_MODEL` and
puts `--save` **before** the `chat` subcommand. Other implementations may expose
a different launcher; use their documented model selection and save-path options.

Configure credentials through your normal authorized provider setup. Keep keys
out of commands and prompts. This example uses the declared
`claude-sonnet-4-6` row; choose a supported model for your implementation before
beginning this bounded exercise.

```sh
work=$(mktemp -d "${TMPDIR:-/tmp}/ensemble-ch10.XXXXXX")
cd "$work"
printf 'hello from disk\n' > note.txt
LLM_VENDOR=anthropic LLM_MODEL=claude-sonnet-4-6 \
  /tmp/ensemble-ch10 --save "$work/save.json" chat
```

Ask the agent to read `note.txt` and summarize it. Verify an actual `read_file`
call and result in the transcript; the answer alone could have guessed the file
contents.
Wait for the completed answer, then send EOF with Ctrl-D on an empty line.
Wait for the process to exit successfully and check that `save.json` exists.

Restart from that same directory with the same model and save path:

```sh
LLM_VENDOR=anthropic LLM_MODEL=claude-sonnet-4-6 \
  /tmp/ensemble-ch10 --save "$work/save.json" chat
```

Before sending a prompt, verify that loading the archive alone has not made a
provider request. Ask, "What did I ask you to read?" Check the restored answer,
then send EOF again and confirm another successful exit. If driving the
accepted CLI over JSONL instead, omit `chat` and send a live-protocol record
such as `{"kind":"prompt","text":"What did I ask you to read?"}`; close stdin
and wait for completion. Do not mix that envelope with the human chat input.

The explicit path above selects the ordinary filename. In a second fresh
directory, use `--save chosen.json` on both starts to check the override: it
should resume `chosen.json` without also creating `save.json`. Keep the complete
save from the first run for the independent variants below.

Repeat the tool-bearing turn through an owned GUI instance and browser profile,
using a separate scratch directory and an unused loopback port. Keep read-aloud
and microphone inactive. Send the first prompt through the GUI's JSONL stdin,
let its tool turn finish, and record the browser's tool/result cards and IDs.
Close stdin, wait for the server to exit, and confirm that `save.json` appears.

Restart from that same GUI working directory and reopen the browser view. The
restored conversation and old card addresses should be present before a new
prompt, with no provider request caused merely by loading the archive or
requesting state. Make a new tool-bearing turn, then check that its new card IDs
do not collide with the old ones and that old controls still select old results.
Close the owned server and browser after the check. A backend replay test alone
cannot establish these browser properties.

Then break copies on purpose. Preserve the complete first save, and make each
variant independently from it. Delete `"log"` down to `[]` and the agent
still remembers, because the log was never on the rendering path. In another
copy, set `"context"` to `null` and it still remembers, because the full log
rebuilds the same context. Do not make that second edit to the already trimmed
copy. With a valid backup retained from the second save, truncate the primary
and verify recovery from that backup; repeat independently with a missing
primary. Confirm the recovery notice, boundary and next save, and that corrupt
primary bytes never replace the backup. With both candidates invalid, verify
non-zero exit, zero provider requests and unchanged bytes. Exercise a failed
replacement locally too: a valid recovery copy must survive. These local
failure checks need no additional paid model requests.

The first-edition reference also shipped a `verify` subcommand that rebuilt
from the log and printed `MATCH` or `MISMATCH` against the saved
snapshot. It remains an optional debugging aid outside the contract; the
accepted Edition 2 solution does not need it. The grader never trusts an agent's
opinion of its own determinism; it compares the requests.

Historical [Claude](../../solutions/evidence/ch11/anthropic-live.json),
[OpenAI](../../solutions/evidence/ch11/openai-live.json) and
[Gemini](../../solutions/evidence/ch11/gemini-live.json) sessions each received:

> Remember the secret word cobalt-juniper. Reply briefly; do not use tools.

A second process for each provider was asked for the word without supplying it
again, and each answered `cobalt-juniper`. Claude also resumed independently
from snapshot-only and log-only variants, with the same answer. All response
processes exited normally without a timeout. The truncated-save process exited
1 before any HTTP request and left the damaged bytes unchanged. The receipts
retain the actual inputs, replies and captured conversation payloads. Their
model choices are historical, not an account inventory for your next run.

Remembering the word demonstrates useful restoration, but it cannot establish
byte identity. A model might answer correctly despite duplicated history or
a dropped tool result. The deterministic source checks use all three real
renderers with recorded tool identities, opaque material, a redaction and pending
ephemera. Snapshot, full replay, older-snapshot-plus-tail and snapshot-only
variants produce the same rendered request. Appending afterward verifies the
next sequence and the next save's anchor as separate obligations.

The original full grader was interrupted and unscored during inherited parity.
Its separately invoked Persistence phases scored **70/85** on the initial
submission. The shape check stopped at the undefined `think` declaration, before
checking later sequence and anchor predicates; those predicates cannot be called
original-grader passes. Independent final-source inspection and focused tests
verified the valid shape obligations after the effective-prompt correction.
The [original review](../../docs/edition-2-notes/ch11-code-review.md) preserves
that distinction and the live receipts' earlier source provenance. The
[Edition 2 grader report](../../docs/edition-2-notes/enhancement-10-code.md)
records revised results separately.

The fresh Edition 3 run encountered a different problem. Its first frozen
Chapter 10 grade was **80/100**: both default-load checks failed, while snapshot,
replay-tail and refusal checks passed. The student preserved the failed report
and stopped rather than changing persistence code to fit an uncertain diagnosis.

The default-load fixture had supplied `CURRENT_RUNTIME_PROMPT` with a terminal
newline, while its unchanged exact expectation had none. Correcting the fixture
input resolved the mismatch; no student production defect was established and
no prompt trimming was requested. The production source stayed unchanged.
Chapter 10 then passed **100/100**, and the same frozen CLI passed Chapter 1,
including its delayed-header and delayed-body phases. Restoring conversation and
replacing runtime configuration remain separate operations: the running agent's
configuration still wins on load.

The operational exercise then took a real tool turn through EOF and relaunch:
read a small file, save, and answer from the restored conversation. A first paid
attempt failed at the transport boundary; a bounded TLS recovery completed both
launches. The original failure's cause remained unproven. Separately, an owned
browser exercise used a local fake provider to check old and new card IDs across
a GUI restart. Its DOM assertions passed; captured screenshots were not a human
visual or audio review. These observations support specific restart claims. The exercise covered one
live model and the asserted browser behavior; other providers and unobserved
interactions remain outside that evidence.
