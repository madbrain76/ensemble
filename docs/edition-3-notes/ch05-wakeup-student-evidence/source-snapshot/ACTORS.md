# Chapter 5 actors

`root.NewAgent` starts one mailbox owner. `Post(UserMessage{...})` and `Post(Hint{...})`
are copied admissions, not acknowledgements: at an idle/interrupted boundary they
start a prompt; during a turn they add one-request context without replacing that
turn. `Post(Interrupt{})` detaches model/tool observation and pairs unfinished calls.
It never kills a retained job. Explicit `kill_job` remains a different operation.

`Submit`, `SubmitMessage`, and `Ask` queue independent requests. Each has one buffered
reliable reply, sent after its turn-end fact and configured history save. `Drain`
is a reliable EOF barrier for previously admitted work, not an observation predicate.
`Config.HistoryPath` enables atomic JSONL replacement; both executables set it.
Library callers may deliberately keep memory-only history by leaving it empty.

`Observe` callbacks run on the actor after reduction: they must not block, call
`Ask`, or perform synchronous shutdown. `Wait` and root `WaitAny` consume bounded,
lossy observation queues (64 and 128 entries); these are UI conveniences, never
completion authorities. Observers receive detached Part payloads. Returned job
facts and tool-result parts are distinct from the assistant's signed response.

HTTP/rendering and tool observation run in workers. Tagged late model results do
not commit or dispatch after interruption; retained jobs keep real late output.
Root shutdown closes all admission first, cancels HTTP/observation, kills and
joins managed processes, joins actual ordinary Go handlers, drains terminal job
reports, and saves. A handler that never returns can prevent shutdown; Go does
not provide safe forced termination of an arbitrary handler.

`ToolBatchLimit=0` selects `DefaultToolBatchLimit` (200); a positive limit counts
dispatched batches. The last permitted batch may produce one final answer. A
subsequent requested batch is paired with limit errors without dispatch. Repeated
explanation-policy refusals have an independent 200-refusal safety bound.
Applications accept `LLM_TOOL_BATCH_LIMIT` as a nonnegative override.

## CLI

Legacy JSONL `{"user":"..."}` requests remain sequential and produce reliable
`{"assistant":"..."}` replies plus cumulative usage. The live protocol accepts
`{"kind":"prompt","text":"..."}`, `{"kind":"hint","text":"..."}`, and
`{"kind":"interrupt"}`; prompt/hint meaning is classified by the actor's state.
It emits `state_changed`, `part_delta`, and `part_final` observations. Explicit
legacy `{"hint":"..."}` and `{"interrupt":true}` are also accepted. `{"parts":[...]}`
submits rich neutral parts. Malformed/unknown commands fail, never become prompts.
`chat` provides `/hint TEXT`, `/interrupt`, and Ctrl-C. EOF drains admitted work.

## Media

`LookupModel` has explicit rows, with no unknown-model fallback. Images are supported
on the selected Anthropic/OpenAI/Gemini models; Anthropic also supports PDF and
Gemini supports audio/video/PDF. `blob` parts need a `RefPath` locator and MIME;
optional `sha256:` digest mismatches fail. Bytes are loaded off the actor into a
copied provider request. Redacted parts are skipped before path access. URI and
opaque handle fetching are not implemented. No bytes are added to persisted history.

## Application workflow

`examples/ch06` is a separate public-importing module, not a framework graph. It
owns author/editor/reviewer prompts, role tools, a synchronized `manuscript.txt`,
`review.json`, tagged observations and `CH06_LOG` (plus per-role history files).
Run `(cd examples/ch06 && GOWORK=off go run .)` with the same provider environment
as the CLI. Input is `{"kind":"prompt","text":"..."}`, role-targeted
`{"kind":"hint","agent":"author","text":"..."}`, or
`{"kind":"interrupt","agent":"author"}`. Handoffs wait for each reliable Ask
completion, never observer timing. Each successful handoff emits `turn_ended`
and `author_done`/`editor_done`/`reviewer_done` application markers.
