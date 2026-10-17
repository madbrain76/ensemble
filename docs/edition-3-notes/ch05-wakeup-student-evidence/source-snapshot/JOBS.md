# Chapter 4 public jobs and extension surface

Create a root with `ensemble.New(logger)` and an Agent with `root.NewAgent`.
`Config.Workspace` fixes relative tool paths; `Config.DataDir` owns persistent
`cr/io/<handle>` output (defaults to Workspace). Each fresh shell starts in
Workspace; a per-call relative `cwd` resolves against it and appears in reports.
No process-wide `Chdir` is used. `root.Shutdown()` kills/reaps running process
groups and records shutdown events; call it **before** `Log.Dump`.

`Agent.RegisterTool(ensemble.ToolDefinition{...}, ensemble.ToolHandler)` adds a
tool without replacing builtins. Handlers receive `ensemble.ToolContext`, which
embeds their cancellation context and exposes the actual Agent and Job owners.
`ctx.Agent.Parent().Logger()` reaches the composition-root logger. Return output
normally; central dispatch persists it and enforces the same cursor/cap/limits
policy as builtins. Duplicate names and non-object schemas are rejected. The
independent `examples/ch04-consumer` module is a runnable public-only example.

Every ordinary invocation creates a retained record before `tool_called`.
`JobRecord` describes handle, running/done/killed status, output Ref, total byte
count, cwd override and exit code where applicable. Short completed output is
inline; running and truncated reports identify the handle and persistent output.
Waiting consumes **all** newly observed bytes, including omitted middle bytes.
Subsequent waits show only unseen output; use `read_file` on the locator to
retrieve retained content. A successful log load replays events, never live jobs.

`wait_for_job`, `send_input`, `kill_job`, and `tool_limits` are inline controls,
not jobs. Wait defaults are 3 seconds and 16 KiB. Delay and regex affect only
observation, never execution lifetime. Patterns match unseen output only.
`tool_limits` affects exactly the next invocation (including controls, unknown
tools, and failures), with explicit wait/run options winning. Its consumption is
noted in that result. `run_command` starts a fresh `/bin/sh` under a 200x50 PTY,
TERM=dumb, stripped provider credentials, merged output and normalized CRLF.
Input gets a trailing newline unless already present. No implicit timeout kills
work. Kill/shutdown terminate the process group and wait for reaping; concurrent
cleanup waits for the same completion. Go handlers may ignore cancellation and
keep running: killed status is terminal, waiters wake, and late output/results
are discarded. The result explicitly distinguishes this from stopping a goroutine.

Direct public `Agent.CallTool` is library-initiated dispatch. LLM-generated calls
continue through `Agent.Ask` and the inherited selected-model visible-explanation
gate **before** job creation or side effects. Refusal leaves only a transcript
pair and does not create a fake job.
