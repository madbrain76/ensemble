# Finding history the Go implementation left behind

Bill's conversations live on disk as the Go version left them: the repository root
was the Go data directory, so `<root>/save.json` is the primary Agent and
`<root>/agents/<name>/save.json` are the six peers named in `<root>/ensemble.json`
(2026-11-05: 2,515 + 13,131 events, 4.6 million characters of text). The C# GUI
showed none of it, because it was told to look somewhere else.

There is **no import command**. `Examples/Gui` calls `Ensemble.RestoreGoHistory`
(`agent/Application/GoImport.cs`) after resolving its data directory and before
composing any Agent, and the operator sees one line when something arrived:

```
Restored: restored 7 Agents from the Go history in ~/projects/ensemble (15547 events)
```

Discovery looks in the places the Go version was ever started from — `EN_GO_DATA_DIR`
if set, then the working directory, then `~/.config/ensemble` — and never searches the
filesystem. It is deliberately conservative, because it runs unattended at every
launch:

- The source is only ever read. Verified by hashing the Go files before and after.
- An Agent whose home here already holds history is **never** replaced, so a re-scan
  cannot overwrite work done in the C# version.
- A directory that merely happens to contain a `save.json` is not Go state: Go's
  snapshot keys (`turn`, `disabled_bands`, `sent_through`, `ephemera`,
  `startup_tools`) identify it. A valid C#-shaped save in the candidate is left alone.
- A failure never stops the GUI. One damaged peer is logged and the healthy ones still
  arrive; any other fault logs one line and serves the fleet it has. Losing a
  conversation is bad, losing the application is worse.
- A marker, `<data-dir>/.go-history-restored`, records that the restore ran. Deleting
  one imported Agent therefore stays deleted instead of resurrecting on next start;
  deleting the marker is the way to ask for a re-scan after adding a new Go peer.

**Where the Go primary lands.** If this data directory's primary Agent has never said
anything (no text and no tool call), the Go primary is adopted as the primary — that is
what "the same conversations as before" means for the main view. If it has said
something, the Go primary arrives as a peer named `Go archive` and the live conversation
is untouched. An unreadable existing home is also treated as busy rather than discarded.

## What already matched, and what did not

The nested event envelopes are the same design, which is the reason this is a
migration rather than a translator: `message`, `request`, `response` and `tool`
carry the same fields, and the peer registry is the identical four-string-field
array, so `ReadPeerLaunches` accepts Go's file unchanged.

Four things did not line up, and each is fixed where it is understood rather than by
loosening the reader for foreign input:

1. **Go's compacted snapshot.** `context` held `disabled_bands`, `sent_through`,
   `turn` and `ephemera`; this codec has `memories`, `phase` and `compression`. The
   event log carries the same facts, so the snapshot is set to null and History
   rebuilds its own by replay. `JsonCodec.ReadSave` already accepts a null context.
2. **Go's wider vocabulary.** `tools_changed` has no counterpart here (the tool set
   comes from the saved configuration), and `job_finished` records a job's lifecycle
   rather than a call's outcome. The adaptation loop does not hard-code this: it asks
   the real reader what it cannot parse and retires exactly that, **by name**, so an
   unexpected type is reported instead of quietly skipped, while a genuinely damaged
   event still refuses the import.
3. **Turn reason.** Go omitted `reason` on `turn_ended` and put the interrupt cause
   under `error`. This codec requires `reason` for that event group
   (`agent/Serialization/JsonCodec.cs:428`).
4. **Empty text.** Go's `Text.Text` is `omitempty`, so a legitimate empty reply
   carried no key at all (`agent/internal/common/types.go:139`: "including legitimate
   empty text"). Absent and empty are the same fact, so the key is written as empty.

Two renames are faithful translations, both checked against the data first:

- `recall_attached` → `recall_added`, joining Go's nested text part and its parallel
  hash list into the single recall part this codec requires. Applied only when the
  Go shape is exactly one text part plus string hashes; anything else is left for the
  reader to report.
- `job_report` → `tool_observed`, because a peer call has no assistant call behind
  it, which is precisely how `agent/Agent/Actor.cs:863` chooses that kind.

`job_finished` is **not** mapped, and the reason is worth keeping. It looks like a
tool result, and was first mapped to one; the import then refused with "Tool event
needs call ID", because all 668 of Bill's `job_finished` events carry an empty
`call_id`. Go recorded the job, not the answer to a call. Inventing an identity would
attach job output to a call that never asked for it, so it stays a reported drop.

## Sequence numbers move, and their references move with them

Replay requires an unbroken run from one (`agent/History/Persistence.cs:18`), so
retiring events leaves holes that make the whole save unreadable. Survivors are
renumbered, and the only other fact in a save that names a sequence is a redaction
range: both bounds are taken from **inside** the original range, because snapping to
the nearest neighbour on the wrong side would silently redact records the Go history
never covered. A redaction that no longer reaches anything surviving is retired with
the records it described.

This is the one place where restoring changes the identity of a record, and it is
recoverable: the Go files are never modified, so an original sequence can still be
looked up there.

## `io/` is deliberately not read

Each Go agent home holds an `io/` directory of bodies named by artifact id. An
earlier reading of the C# port treated those as the missing text of the 417 keyless
parts. They are not safe to use:

- Nothing in the build that wrote these saves reads or writes `io/`: no Go source
  path touches it, and no saved part carries a `ref`, so nothing points at it.
  (Job payloads do carry `tool.job.output.locator` addresses into it, which this
  codec preserves as `JobInfo.OutputPath` — that is separate from part bodies.)
- Its ids were reused across restarts: in `ensemble debugger`, artifact ids run
  10..48 twice, and `io/400` is dated a month before `io/8`. One file cannot serve
  two parts.

So parts whose body is absent import as empty and are counted in the report
("314 text parts carried no body in the Go files"). Blank cards in a resumed
conversation are a gap; another turn's words in an Agent's mouth would be a lie.

## Verification

`agent/Tests/GoImportTests.cs` drives the launcher's own entry point and writes its
fixtures in Go's shapes — a model call as a `tool_call` part inside a response, each
response behind its own request, a result as text parts on a `tool_returned` — because
an invented shape would test the fixture rather than the restore. Three cases:

- **found and resumed**: discovery from the working directory with no environment hint,
  both Agents installed, drops and translations named in the log, source bytes
  unchanged, registry merged, recall text and hashes intact, call id intact, redaction
  range moved 12..14 → 10..12, an unbroken run from one, the Go-redacted reply absent
  from the render, and the surviving conversation rendered by `Engine().Render`.
- **adoption**: a primary that has never spoken adopts the Go primary; one with a live
  conversation keeps it and receives the Go primary as `Go archive`.
- **stays out of the way**: one damaged peer does not withhold the healthy ones; a
  valid non-Go save is not adopted; the marker keeps a deleted Agent deleted; an
  existing save survives a forced re-scan; a damaged registry is left alone.

Suite: 344 passed, 4 failed — the pre-existing `agent/Jobs/Pty.cs` macOS gates. Five
mutants each turn one distinct assertion red: adopting into a busy primary, ignoring
the marker, dropping the working-directory fallback, accepting any `save.json` as Go
state, and overwriting an Agent's existing history. Three of those five survived an
earlier version of the test, which is why the fixtures are valid saves rather than
arbitrary bytes: a fixture the restore could not read anyway proves nothing about
refusal.

Real data: all seven Agents restored, `StartPersistence` opened every one, and
`Engine().Render` produced a provider request for each (396 KB for the primary,
715 KB for `ensemble debugger`). In Firefox, starting the launcher from the repository
root with no extra flags brought back the primary conversation and all six peers, and
selecting a peer replaced the transcript with its own cards.

Two limits worth repeating, because both look like bugs from the inside. Some parts
arrive with no text, for the `io/` reason above — the launch log counts them per Agent.
And the system prompt and model are live composition here, read from the launcher and
the peer registry, not from a saved history: Go recorded them for reproducibility, so
importing them would silently change which model an old conversation now talks to.
