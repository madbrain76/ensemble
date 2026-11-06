# Linux needs its own PTY launcher, not a pipes fallback

Branch `fix/linux-pty-launcher`, commit `2f6ba9ca`, from `fix/linux-path-case`.
Two product files and one test file.

## Symptom

`run_command` never ran. Each call produced:

```
tool_result (is_error): "Turn failed before dispatch; not executed."
turn_ended  reason: "This chapter endpoint uses the macOS system PTY launcher"
```

`tool_called` and `tool_started` shared a timestamp with `tool_returned`, which is
how you can tell the child never started: nothing was executed, so nothing could
fail. A batch of calls lost all of them, because the first throw abandons the
rest of the turn. The agent therefore looked idle while it was being refused
before dispatch, and it correctly reported the failures it had been handed.

## Cause

`Job.RunCommandAsync` (`agent/Jobs/Pty.cs:8`) had one implementation and it began
with `if (!OperatingSystem.IsMacOS()) throw new PlatformNotSupportedException`.
It is the book's endpoint: `book/chapter-04.md:493` says the PTY implementation
is macOS-specific and never falls back to pipes, and `:746` adds that other
operating systems need their own tested launcher.

That is a requirement, not a defect, so the fix is to supply the launcher the
book asks for rather than to relax it. Pipes were not used: `send_input` needs a
controlling terminal (`agent/Jobs/Job.cs:438` refuses a job whose process group is
zero) and `kill_job` needs a fresh group to signal.

util-linux `script` provides both, verified on this host before any code was
written: a real controlling terminal with `stty size` reporting the requested
geometry, and a child that is its own PID, PGID and SID.

## Three differences the launcher now selects per platform

| | macOS `/usr/bin/script` | util-linux `script` |
|---|---|---|
| spelling | `-q -F /dev/null <sh> -c <bootstrap> ensemble-pty <cmd>` | `-q -e -c <bootstrap> /dev/null` |
| command transit | positional `$1`, never re-quoted | environment, expanded by the bootstrap |
| exit status | propagated | only with `-e`; without it the launcher exits 0 for a failed command |

The command cannot ride as a positional parameter on Linux: `script` there
rejects trailing arguments with "unexpected number of arguments". Interpolating
it into the `-c` string instead is the obvious wrong answer, and a mutant proved
it: the bootstrap shell ate `$X`, so `X=hi; printf '%s' "$X"` reported `GOT[]`.
The environment keeps the command exactly as the model wrote it.

`-e` deserves its own test because nothing else noticed. Its failure is silent
and it is a lie with a number in it: a failed kernel build would be reported as
`exit_code: 0`, which is the opposite of what this book teaches about tool
results.

`Signal` was `[DllImport("libSystem.B.dylib")]`, so signalling a group would have
thrown `DllNotFoundException` on Linux. Each platform now names its own library
behind one chooser, which is the only caller. The existing kill test is what
covers it: it kills a `run_command` job and asserts a backgrounded grandchild
died, and that grandchild can only be reached through the group.

## Results

`agent/Tests` on `fix/linux-path-case` went from 341 passed / 4 failed to **343
passed / 2 failed**. `Chapter 3 tools and loop` and `Chapter 4 jobs and PTY` pass
on Linux for the first time, and that is the better evidence: the graded chapter
test asserts `test -t 0 && test -t 1`, `stty size` reporting `50 200`, and a
`send_input` round trip through `read`, so the book's own exercise now verifies
this launcher on this platform.

The remaining two, `Chapter 19 confined process boundary` and `... useful work`,
fail at `agent/Sandbox/Command.cs:35`, "Confined commands require the macOS
process sandbox". That is Seatbelt and is a separate platform facility; the
launcher is not involved, and the GUI does not request confinement, so ordinary
tool use is unaffected. Porting confinement to Linux means `bwrap` or Landlock and
its own threat review, which is a decision rather than a fix.

## Mutation checks

Three, each killed by a distinct assertion:

- drop `-e`: "a failing command must report its own exit code, not the
  launcher's: halfway".
- interpolate the command into `-c` rather than passing it through the
  environment: "the command was re-quoted or expanded too early: GOT[]". That
  assertion needed the two earlier checks neutralised first, because the same
  mangling broke the PTY case before reaching it. It can fail on its own, which
  is what the experiment established.
- deleting the group signal is not separately checked; the kill case would still
  pass on .NET's tree walk, so it pins the binding rather than the choice.

## End-to-end on the product path

Real launcher, copied data directory, real Firefox, prompt typed into the
composer, with the imported Go conversation of about three hundred messages, and
a scripted provider that asked for one `run_command`. The tool result the model
received:

```
/dev/pts/8madbrain
MARK os=Linux pid=23378

exit_code: 0
```

A terminal device, the operator's identity, a real child PID, and a reported
status. Provider saw two calls, the reply rendered, zero page errors. Two earlier
harness attempts were wrong rather than the product: single quotes around
`$(uname -s)` made the shell correctly print it literally, and `test -t 1` inside
`$( )` is always false because the subshell's stdout is captured.

## Not done

`agent/GuiTests` was re-run (clean) although no GUI code changed. Confinement on
Linux is untouched, as above. `fix/platform-refused-job-command`, which turns a
refused launcher into ordinary tool data instead of an application abort, stays a
separate branch; a refused platform is still an exceptional event and should stay
loud.
