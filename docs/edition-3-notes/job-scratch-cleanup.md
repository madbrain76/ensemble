# Job scratch reclamation (November 5)

## Symptom

`/tmp` held 27,500 entries, 23,672 of them `ensemble-jobs-<guid>` directories, all
owned by the operator, most empty and the rest holding one or two job output files.
Nothing ever removed them.

## Mechanism

`Jobs` builds its output directory in its constructor (`agent/Jobs/Jobs.cs:12`) as
`Path.GetTempPath()/ensemble-jobs-<guid>` when the owning Agent supplied neither an
output directory nor a sandbox root, and creates it immediately. `OutputDirectory` is
`readonly` (`agent/Common/Jobs.cs:28`) and the next `Jobs` draws a fresh GUID, so once
the owning Agent is gone **no object in the program can name that directory again**:
it is unreachable by construction, not merely unused. `Jobs.Destroy()` joined children
and removed the reverse membership but never touched the directory, so every no-home
Agent ever constructed left one behind forever.

`Agent.ConstructOwners` (`agent/Agent/Agent.cs:40`) defaults the directory to
`<home>/jobs` only when the Agent has a home, so a named Agent keeps its retained
output under its data directory and is not involved. That is why the GUI root agent
(`Examples/Gui/Program.cs:167`, named `root`) leaks nothing: measured, a real launch
created zero `/tmp` entries and put its jobs at `<data-dir>/agents/root/jobs`. The
leak came from no-home Agents: the `Examples/*` programs and the test suites, which
construct roughly four hundred of them. One `agent/Tests` run leaked 392 directories.

## Fix and its boundary

`Jobs` now records whether the constructor invented the name
(`temporaryOutput = outputDirectory is null && files.Root().Length == 0`) and removes
that tree at the end of `Destroy()`, after children are joined.

Two cases are deliberately **not** removed, and each has its own case in
`JobScratchTests`:

- a directory the caller named, which is retained output somebody chose;
- a directory under a confined Agent's root, which `Jobs` documents as retained output
  for inspection and which lives where the operator pointed it.

Deleting either would destroy evidence to satisfy a tidiness argument. Removal failures
are swallowed: this is operating-system scratch, its own sweeper may be racing us, and
failing an Agent shutdown because `/tmp` is busy would be worse than leaving a
directory. The path is refused if it is a reparse point.

## Verification

- `Chapter 19 temporary Job scratch is reclaimed at destroy` fails when the delete is
  commented out, so the case is live rather than decorative.
- One `agent/Tests` run: 392 directories leaked before the change, **0 after**.
  One `agent/GuiTests` run: **0**. One real GUI launch: **0**.
- `agent/Tests`: 344 passed, 4 failed — the known `agent/Jobs/Pty.cs` macOS gates,
  unchanged from the base branch.

## Separate finding, same directory

`run.sh` created `mktemp -d /tmp/ensemble-profile.XXXXXX` on **every launch** and
exported `EN_CPU_PROFILE` into it. No C# code reads that variable; it came from the Go
launcher, where GUI CPU profiling is feature 16 and still unported. 36 empty
directories accumulated without the file ever being opened. The block is gone, and the
variable is now passed through only if the operator sets it.
