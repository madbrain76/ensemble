# Recover an actor from a rejected journal

## The failure

`History.LoadSave` called `ReadJournal` on the recovery journal *before* the
loop that tries the primary save and then its backup. A journal that
`ReadJournal` rejects therefore threw out of the load entirely: no backup was
attempted, nothing was logged about the offending file, and the operator saw
only `InvalidDataException. Check private diagnostics ...` against empty logs.
One damaged recovery tail made the actor permanently unstartable even though
`save.json` and `save.json.bak` were both valid.

## Decisions

**Say which file and why.** The rejection is logged with its path and the inner
reason, and separately when nothing can be continued from. This is the same
shape as the save-candidate logging added earlier, which had left the journal as
the one silent path.

**Salvage the usable prefix rather than discard the file.** Records before the
first gap are contiguous and durable, so `ReadJournal(path, strict: false)`
returns them and recovery continues at that point. `StartJournal` and
`TrimJournal` keep the strict reader: they rewrite the journal in place, and a
partial read there could publish a file that quietly omits records. Keeping one
strict and one lenient reader, chosen by the caller's exposure, is smaller than
making either reader conditional on who called it.

**Preserve the rejected bytes under a new name.** Continuing requires the
damaged file to stop being the active journal, because `StartJournal` rewrites
the active file and would destroy the only copy of what it could not read. It is
moved whole to `<journal>.rejected-<utc>` first. A failed move aborts the load
rather than falling through to that rewrite.

**Refuse when the journal is the only record.** With no usable snapshot there is
nothing to continue from, so the load rethrows and leaves every file exactly as
found. Verified by `Chapter 14 journal refusal without snapshot`.

Meantime behaviour, on a copy of the real damaged directory:

```
persistence: .../save.json.journal.jsonl rejected (InvalidDataException: Journal
sequence is not contiguous); preserved as .../rejected-20261105-023326;
continued from the saved snapshot at sequence 17 with 3 journal record(s) before
the damage
```

The preserved file held all four original records; the rebuilt journal came out
`[17, 18, 19, 20]`.

## What recovery costs, honestly

Records after the last save and after the first gap are not applied. They stay
readable in the preserved file, but they are not in the actor. In the observed
case the salvage recovered everything contiguous and the loss was nil; a longer
unsaved tail would be a real loss, and no client-side fix can invent those
records. The alternative was refusing to start forever.

## Deliberately not changed

`ensemble.log` is opened with `FileMode.Create`, so a relaunch erases the
previous session's host diagnostics, and `DiagnosticTests.cs:181` pins that
("Host restart appended diagnostics"). It is why Bill saw an empty log after
relaunching: the surviving log belongs to the newest launch. Truncating per
session is a tested contract with a bounded file (a 64 MiB cap plus a reset
marker), so it was left alone rather than quietly flipped to append. Read the
diagnostics **before** launching again.

## Test notes

`JournalRecoveryTests` builds damaged state through the real `History` API and
codec, not hand-written JSON. `Owner` is the local agent helper; naming it
`Actor` collides with the `Actor` enum used for `Actor.Human` in the same file.
Assertions were mutation-checked: discarding instead of salvaging fails the
survival check, and removing the quarantine fails the active-journal check.
