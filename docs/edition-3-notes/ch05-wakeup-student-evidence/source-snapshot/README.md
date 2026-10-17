# Chapter 2 — one log, three vendors

Built on accepted Chapter 1 (`1aad3975842d502b9144b7b9bddfd30ca7309db7`,
annotated tag `edition-3-ch01-r1`). The root remains the `ensemble` library;
the existing CLI stays at `./cmd/ch01` to preserve its public launch path.
No SDK, streaming, tools execution, media loader, GUI transport, retry framework,
or blanket generation deadline is introduced.

## Build and use

```sh
go build -o /absolute/path/to/student-cli ./cmd/ch01
# Export credentials securely, never as command arguments.
# LLM_VENDOR: anthropic (default), openai, gemini
# LLM_MODEL, LLM_BASE_URL, LLM_API_KEY
# Optional: LLM_RESOLVED_MODEL (exact replay identity).
# CLI system is "You are a helpful assistant."; max output is 1024 tokens.
# CH02_LOG defaults to ./ch02.jsonl.
/absolute/path/to/student-cli chat  # human conversation, EOF exits
/absolute/path/to/student-cli       # JSONL protocol
```

No-argument protocol accepts one `{"user":"..."}` or
`{"ephemeral":"..."}` object per line. The latter queues instruction text for
one request and emits `{"ack":"ephemeral"}`. Successful sends emit
`{"assistant":"..."}`. EOF emits cumulative measured
`{"usage":{"input":0,"cache_write":0,"cache_read":0,"output":0}}` with real
counts. Diagnostics are stderr-only. Sessions save their normalized log on both
success and failure. Anthropic's predecessor `ANTHROPIC_*` environment settings
remain fallbacks. `LLM_RESOLVED_MODEL` changes replay identity, not routing;
a returned model ID takes precedence when parsing a response.

```sh
# Fresh-process offline use; no key required and no network call is made.
/absolute/path/to/student-cli dump > /tmp/saved.jsonl
LLM_VENDOR=openai LLM_MODEL=gpt-4.1-nano \
  /absolute/path/to/student-cli render /tmp/saved.jsonl
```

`render LOG` reads a JSONL file, rebuilds Context, and writes exactly one JSON
request. Repeat it to compare bytes or choose another supported vendor. Config's
stable system prompt is outside the log. The CLI uses the same fixed prompt
for live and offline calls; library clients supply their own fixed Config.System.
Logs begin with `{"log_version":1}`; headerless current logs are also accepted.
Unknown events/parts/actors/vendors/surfaces, obsolete blob `path`, malformed
references, and unsupported media rendering fail loudly. RefPath/RefURI/RefHandle
survive dump/load, but their bytes are never opened/fetched by this chapter.

## Owners and public seam

Clients import only `ensemble`. Construct `ensemble.New(logger)`, then
`root.NewAgent(ensemble.Config{...})`. Call `agent.Ask(ctx,text)` sequentially.
Offline clients may omit provider credentials. Live-send validation happens
before HTTP. The caller owns cancellation; the HTTP transport bounds only
connection establishment, not response headers, response body, or generation.

Ensemble retains Agents and its logger. Agent retains immutable Config, Engine,
and History. Engine owns HTTP and measured live usage. History is the sole owner
of facts, validation, load/dump, and the pure Context reducer. The composition
root assembles concrete implementations; spokes depend only on shared interfaces,
not sibling implementations. Each child exposes its real Parent. Public aliases
keep clients out of `internal/`.

- `agent.Log().Append/Events/Load/Dump/Context` operate on the real owned History.
  Append and snapshots deep-copy mutable content. Failed loads are transactional.
- `agent.Engine().Render(context, provenance)` is pure, deterministic wire JSON.
- `agent.Engine().Parse(bytes, provenance)` returns normalized response facts;
  it does not append or change usage. Agent records returned facts after Send.
- `agent.History()` is the retained Chapter 1 committed-text compatibility view.
  Full facts (including failures) are available only through Log, not this view.
- `Engine.Usage()` tracks this engine's measured live responses (including an
  unusable but measured answer). `Context.Usage` reduces captured completed
  responses; an error-only fact is not a fabricated zero-usage response.

Context contains turn state, dialogue, pending ephemera and cumulative captured
usage; it never contains provider wire names. `request_sent` consumes ephemera.
Every event's effect is explicit, including identity transitions. Redactions
transform only projected spans; Ref-bearing result stubs retain exact references
in Context but render only their inert stub text. Original log facts remain.

Renderers use typed wire structs. Tool results keep provider-specific authorship;
Anthropic merges user material and stably moves **all** tool results ahead of
text. Gemini resolves function names from preceding calls. Missing wire call IDs
are generated from event sequence/part position where needed and never stored as
provider facts. Opaque blocks and per-call signatures replay only to the exact
vendor/model/surface. An empty text value is a wire `""`, never `null`.

## Tests and limits

```sh
go vet ./...
go test -race -count=1 ./...
go test -coverpkg=./... -coverprofile=/tmp/ch02.cover ./...
python3 scripts/check-mutations.py /absolute/evidence/directory
(cd gui && go test ./... && go build -o /tmp/ch02-gui-stub .)
```

Tests enter the public root or the actual command's `run` path and use focused
HTTP protocol fixtures, not mocks or alternate application construction. Two
inherited lifetime checks intentionally wait 62 seconds each. Mutation checks
restore exact source bytes and require behavioral failures plus restored passes.
The optional `gui/` module is only a public-library client stub.

Live evidence demonstrates ordinary text conversations and captured usage on
Anthropic Messages, OpenAI ChatCompletions and Gemini generateContent. Complex
reasoning/signature/tool/ref/redaction shapes are fake-protocol/source evidence,
not claims of live tool execution, paid cache activation, media support or GUI
integration. Grader processes do not contribute to source coverage.

## Chapter 3 — owned, synchronous tools

Built from accepted Chapter 2 `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`.
The root public library and `./cmd/ch01` CLI remain the entry points. The CLI
now enables `read_file`, `write_file`, `edit_file`, `list_directory`,
`search_files`, and `run_command`. A plain library `Config` remains tool-free;
set `Config.EnableTools: true` to install these six Agent-owned operations.
`Agent.Tools()` is a detached registry snapshot; each Tool exposes its Parent,
Definition (name/description/object schema), and context-aware Execute method.
Agent, not the tools or renderer, owns model dispatch and paired History facts.

### Selected-model explanation policy

Set `Config.RequireToolExplanation: true` for the particular configured model
whose tool batches require preceding visible explanation. The CLI equivalent is
`LLM_REQUIRE_TOOL_EXPLANATION=true`; explicitly set `false` or unset it to clear
that choice. The default is false. No model-name convention or hidden selection
is used. Every CLI model receives a stable system instruction to explain tool
use; only this explicit policy selection enforces it in application code.

The gate requires nonblank ordinary text **before the first tool call in that
response**. It does not judge prose semantics. Opaque reasoning, whitespace,
previous-turn text, and text following a call do not qualify. Preflight occurs
before any tool runs. A refused batch produces one paired error result per call
and preserves the full assistant response. Those errors tell the model to retry
with preceding text. A corrected retry executes the real tools. Both refused and
executed rounds consume the same maximum of ten model responses per Ask; no
separate retry allowance or automatic network retry is added. The API returns the
final no-call response, not intermediate explanation text. Full explanations and
results remain available in Log. Missing native IDs receive Agent execution
correlation IDs tied to response sequence/part index; native IDs are unchanged.

### Tool contracts and limits

- Argument decoding requires an object, correct types, required fields and no
  unknown fields/nulls. Tool errors are results, not silent claimed success.
- `read_file`: one-based inclusive start_line/end_line (omitted or zero means first/EOF),
  then max_bytes (omitted/zero: 32768) with explicit truncation. Negative and
  out-of-file ranges fail. `list_directory`: sorted immediate names,
  directories suffixed `/`. `search_files`: recursive Go-regexp path:line:text; context_lines defaults to zero,
  context uses path-line-text with merged windows and -- group separators. Binary files containing NUL and symlinks are not traversed.
- `write_file` creates parent directories; existing files refuse without
  overwrite:true, except append:true preserves bytes and adds content. Append
  takes precedence if both flags are true. Empty content is valid. `edit_file` requires a nonempty exact unique
  anchor; absent, duplicate, and overlapping anchors fail without writing.
- `run_command` uses fresh `/bin/bash -c` in the process working directory.
  Results include separate stdout, stderr and exit_code. Nonzero shell exits are ordinary data, not tool errors.
  Each stream retains 64 KiB with an explicit truncation marker while draining
  the child. Shell cwd/exports do not persist between calls. Ordinary environment
  is retained except supported credential names LLM_API_KEY, ANTHROPIC_API_KEY,
  OPENAI_API_KEY, GEMINI_API_KEY. This is credential hygiene, **not a sandbox**.
  Tools have local user permissions and can access paths outside cwd. Use only
  disposable projects and do not run untrusted commands with accessible secrets.
- Caller context controls cancellation. Commands have no arbitrary short work
  deadline; cancellation kills the shell process group, with a one-second pipe
  drain bound. Provider generation lifetime remains the Chapter 1 policy.
- Read output is capped after selecting lines, but reads still load the file
  into memory. Search output is unbounded; no transactional edit,
  OS sandbox, approvals, parallelism, GUI, job system or streaming is claimed.

Chapter 3 tests explicitly set and clear the public policy through both Config
and CLI environment, exercise actual HTTP wires, real disk/shell effects and
feedback, and cover strict args, cancellation, environment canaries, pairing,
retry and round limits. Live evidence and source/binary identity are maintained
in the separate student report, not encoded as tests or model-name switches.

Signed ordinary Gemini text remains ordinary text: only `thought:true` marks
text as hidden. A text part can retain its original `thoughtSignature` JSON
value in `Part.Opaque`, bound by `Part.From` to the exact vendor, surface, and
model. Gemini same-source rendering attaches that signature to the single
visible text part; other targets keep the text and omit the signature. True
hidden parts remain opaque and never qualify as a tool explanation. History
dump/load preserves these facts without rendering or reinterpreting old logs.
