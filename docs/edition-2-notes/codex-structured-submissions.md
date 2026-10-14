# Stage B: structured child submissions

October 14, 2026. Maintenance coder, GPT-6 Astra/high. Design accepted by coordinator and independent runtime reviewer before implementation.
Formal implementation review remains pending.
Baseline `7a362ec2`, current chapter numbers retained. Mandatory reading includes
`.agents/skills/ensemble-coding/SKILL.md`, AGENTS, the full workflow/carryover,
Chapters14/15/22/23 and all three compaction audit/progress notes. This is
maintenance with prior-answer exposure, not a blind student exercise.

## Proposed boundary

- Add `AgentSpec.SubmitSchema json.RawMessage` and `SpawnWithSchema(id, task,
  model string, schema json.RawMessage)`. Keep `Spawn` as the absent-schema
  wrapper. A child's default derivation explicitly clears its parent's schema,
  including the public `AgentSpec.Child` helper.
- Keep `Submit(string)`; it encodes `{"result": text}` and calls the same
  `SubmitJSON(json.RawMessage)` boundary as the model tool. Custom schemas govern
  the entire direct argument object. All accepted data remains JSON, never JSON
  embedded in a string. `ChildStatus.Result` stays the default text field;
  `ChildStatus.Payload` carries custom JSON only. `SubmissionData(id)` returns a typed `Submission` containing default `Result`
  or custom `Payload`. Legacy `Submission(id)` returns default text and reports
  a clear use-SubmissionData error for custom results.
- The Fleet implementation owns a compiled immutable submission contract. The
  composition root compiles it before constructor side effects and passes it to
  `fleet.New`; Spawn also validates/copies before invoking any custom factory.
  Spec ingress/egress, declaration and result snapshots own their byte slices.
  The factory result must retain the requested schema; a custom factory cannot
  silently discard or replace that contract. No new callback/service layer.
- `submit` obtains its declaration from that compiled contract. Runtime validation
  happens before durable write, state change or join notification. Rejection leaves
  the child unsubmitted, allowing its existing bounded lifecycle to correct it.
  This stage does not change submission-versus-quiescence ordering (Stage C).

## Validator and bounded dialect

Use `github.com/santhosh-tekuri/jsonschema/v6` pinned to `v6.0.3`, with Draft
2020-12 selected explicitly and a loader that refuses all external resources.
The library handles metaschema checking and instance semantics. A small recursive
admission check rejects unsupported keywords; it is not a second validator.
The root must explicitly be `type: object`. Nested boolean schemas are allowed.

Supported keywords: `$schema` (only the 2020-12 URI), `type`, `properties`,
`required`, `additionalProperties`, `items`, `enum`, `const`, `minimum`,
`maximum`, `exclusiveMinimum`, `exclusiveMaximum`, `multipleOf`, `minLength`,
`maxLength`, `pattern`, `minItems`, `maxItems`, `uniqueItems`, `minProperties`,
`maxProperties`, and annotations `title`, `description`, `default`, `examples`.
Nested objects/arrays and nullable type arrays use normal draft semantics;
integer accepts mathematically integral numbers, not merely integer spelling.
All other keywords, references (including local), IDs, dynamic references,
combinators, formats/content and custom vocabularies are explicitly rejected.
This avoids recursive/reference expansion and silently unenforced assertions.
Go's standard RE2-style regexp engine remains the library default.

Limits proposed: schema64 KiB, submission1 MiB, JSON nesting64, schema nodes2048,
instance nodes65536; numeric lexemes128 bytes and exponent magnitude1000. These
are admission ceilings against recursion/allocation and enormous rational-number
construction, not memory summarization targets. Rejection is explicit and never
truncates data. They leave ordinary nested results and the memory schema ample
room. No arbitrary model deadline or token cap is added.

Dependency research: [upstream library](https://github.com/santhosh-tekuri/jsonschema)
and [versioned API](https://pkg.go.dev/github.com/santhosh-tekuri/jsonschema/v6@v6.0.3),
plus downloaded v6.0.3 source. `DefaultDraft`, `UseLoader`, `AddResource` and
`Validate` provide the boundary; `UnmarshalJSON` retains exact `json.Number`
values. Upstream documents standard-library regexp as the default. Downloading
the module did not change the repository dependency files.

## Storage and presentation

Default submissions keep existing `submit.md` JSONL rows `{"result":"..."}`.
Custom rows use an explicit envelope `{"version":1,"submit_schema":{...},
"payload":{...}}`. Each row is self-describing for durable retrieval even when
a same-ID child is respawned with a different contract. Read legacy and custom
rows without promoting an old result into completion of the new task; historical
custom rows are checked against their recorded contract when read. Reject unknown
envelope versions and malformed rows. Ownership lookup precedes every disk read.

Join/status use Payload as raw JSON; `get_submission` returns `agent_id` plus
legacy `result` or custom `payload`. Existing tool-result text carries a JSON
document, as before, but its inner payload is an object. The GUI already renders
tool arguments and result text safely; verify actual selected conversation/replay
and nested JSON rather than introduce a new result browser. Agent startup retains
its caller-supplied immutable schema; saved history preserves old declaration and
tool arguments, and `submit.md` supplies historical result retrieval. Ephemeral
children still are not automatically resurrected.

## Verification

Use real Agent/Engine/Fleet/public constructors with small HTTP provider fakes,
including invalid-then-correct submit, pending Join, direct string API rejection,
default child/grandchild reset, foreign/missing retrieval and durable reload.
Exercise exact memory declarations across Anthropic, Chat Completions, Responses
and Gemini, plus nested values/arrays/null/integer/enum/additional-properties,
invalid schemas/keywords/refs and admission limits. Mutate caller and returned
schema/result buffers and verify retained declarations/results. Check a same-ID
new task cannot become submitted from an old file. Test GUI selection/replay with
a custom result. Run focused races, full core/GUI short/race and vet once coherent;
coordinator owns advertised grader execution/build repair. No live inference.
Targeted scratch mutations will cover the meaningful validation, ownership,
copying, inheritance and persistence claims. Record actual failures/limits here.

Accepted clarifications: default submissions continue accepting extra keys but
normalize to the legacy result string/row; only custom contracts preserve the
full object. Admission checks descend only into schema-bearing positions, never
object values inside enum/const/default/examples. Existing Ask/ask_agent and
TurnEnded text display use compact JSON for custom results, while typed status
and SubmissionData retain raw JSON data. No empty success or new result browser.
