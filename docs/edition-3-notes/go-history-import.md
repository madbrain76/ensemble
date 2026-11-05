# Importing history written by the Go implementation

Bill's conversations live on disk as the Go version left them: the repository root
was the Go data directory, so `<root>/save.json` is the primary Agent and
`<root>/agents/<name>/save.json` are the six peers named in `<root>/ensemble.json`
(2026-11-05: 2,515 + 13,131 events, 4.6 million characters of text). The C# GUI
showed none of it. `import-go` (`agent/Application/GoImport.cs`, dispatched from
`agent/Examples/Gui/Program.cs:8`) closes that gap without asking anyone to move
files by hand, and never writes to the source tree.

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

This is the one place where an import changes the identity of a record, and it is
recoverable: the Go files are never modified, so original sequences can be read from
them. `--dry-run` prints the plan without writing.

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

`agent/Tests/GoImportTests.cs` writes fixtures in Go's own shapes — a model call as a
`tool_call` part inside a response, each response behind its own request, a result as
text parts on a `tool_returned` — because an invented shape would test the fixture
rather than the import. It asserts the drops and translations are reported by name,
the source files are byte-identical afterwards, the recall keeps text and hashes, a
tool result keeps its call id, the redaction range moved 12..14 → 10..12, the history
is an unbroken run, the Go-redacted reply stays out of the render, and the surviving
conversation renders through `Engine().Render(History().View())`. `Refusals` covers
unreadable saves, `--dry-run`, and that a damaged peer registry is never replaced.
Suite: 343 passed, 4 failed — the pre-existing `agent/Jobs/Pty.cs` macOS gates.
Mutating the renumbering, the range translation and the recall join each turn one
specific assertion red, so none of them is decorative.

Real data: all seven Agents imported, `StartPersistence` restored every one, and
`Engine().Render` produced a provider request for each (396 KB for the primary,
715 KB for `ensemble debugger`). Two limits to state plainly. The `io/` gap above
means some assistant and tool bodies are absent; and a peer's model comes from the
registry while the primary's comes from the launcher's flags, because the system
prompt and model are live composition here, not saved history — Go recorded them for
reproducibility, and this build reads them from the launcher and its skills.
