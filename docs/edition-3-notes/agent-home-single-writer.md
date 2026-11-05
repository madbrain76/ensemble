# One writer per actor home

## Why

Bill's launcher began failing with `GUI failed: InvalidDataException. Check
private diagnostics and workspace/data permissions; no private error text is
printed here.` while every log in the actor home was empty. The saved state
explained it: `agents/root/save.json.journal.jsonl` held the sequences
`[15, 16, 17, 16]` and `events.jsonl` held a torn 6 KB record with two JSON
objects concatenated.

Two `ensemble` processes had been alive over one data directory at once -- the
first was the silent server that looked hung because the launcher printed no
address, and the second was the relaunch that replaced it. Each numbered its
events from the sequence it had restored, so both appended a record numbered 16
to the same journal, and the file stopped being contiguous. `ReadJournal`
refuses a non-contiguous journal, so every later launch failed: an actor home
made permanently unstartable by the application's own writes.

## Decision

Take exclusive writer authority over the directory holding an actor's save and
recovery journal, from `StartPersistence` until `Destroy` closes it, and refuse
the second launcher naming the contested directory. `WriterConflictException`
carries fixed, printable guidance, so it is read as a sentence rather than a
type name.

The lock is an ordinary advisory handle, so the operating system releases it if
the holder is killed: a crashed launcher cannot strand the directory behind a
stale lock file. Verified by three consecutive start/stop cycles that left the
journal contiguous.

The refusal was proven between two real processes, which is the failure being
prevented:

```
GUI failed: Another Ensemble is already using actor home .../agents/root;
close it before starting this one.
```

## Scope limits, stated plainly

.NET enforces `FileShare.None` with per-process advisory locks on Unix, so two
`Agent` objects inside one process do contend, and
`Actor home single writer` depends on that. What this guard does not claim:

- It does not stop a writer that ignores the lock file, such as a tool editing
  the journal directly.
- Sibling actors in one fleet keep separate directories and are unaffected.
- Prevention only. State damaged before this change is still damaged, so the
  matching recovery is `fix/journal-recovery-diagnostics`.

## Defect found while testing

`agent/Agent/Storage.cs` cannot reach the sandbox through `Parent()`: an
`Agent`'s parent is the `Ensemble` container, which has no `Sandbox()`. The
agent owns it directly. Caught by the compiler on the first build.
