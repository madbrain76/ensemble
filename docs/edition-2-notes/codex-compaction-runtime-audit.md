# Compaction runtime design audit

Read-only audit by `compaction_runtime_audit`, an independent design reviewer,
2026-10-14. Baseline: `56ec71cbde8477b6959ad2fbc65f9237a0fa8346` in
`/Users/bill/projects/ensemble-compaction`. The working tree was clean when
inspected. This note is the only file written by this reviewer. No implementation,
test, grader, provider call, credential access, deployment read/change, or commit
was performed. Test names below describe inspected coverage, not new results.

Read in full: `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`,
`docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`,
`book/chapter-writing-procedure.md`, and current Chapters 14 and 15. Read the
360-line handoff plan from the main checkout. This is the runtime portion of
Stage A; the coordinator and separate dependency reviewer own the full chapter
mapping, remaining preparation, and sequential baseline results. It is not a
formal prose review or acceptance of the whole handoff.

## Assessment

The plan is feasible through the existing composition root, Agent actor, Fleet,
Memory, History and Engine. It requires real changes at their existing boundaries;
wrapping `Engine.Compress` in a goroutine or calling `EnableFleet` everywhere does
not implement the requested contract. There is no need for a second scheduler,
event bus, mutable conversation builder, memory store, or retry framework.

Keep current chapter IDs while implementing. Chapter 15 explicitly requires a
finite direct call before the next main request (TL;DR rules 1–2 and §15.2/§15.5).
The carryover also teaches that choice. Background actor compaction deliberately
replaces that scheduling/construction contract while retaining selection,
first-person memory, replay, source preservation and ownership.

## 1. Construction and capability boundaries

Verified source:

- `agent/ensemble/fleet.go:26–57`: `FleetEnabled()` means `factory != nil`;
  `CreateChild` fails without a factory. `NewChildAgent` verifies the actual
  parent and calls `AgentSpec.Clamp` before construction.
- `agent/ensemble/registry.go:85–94`: `EnableFleet` installs a recursive factory
  using the parent's Config and replaces any prior factory. `Start` calls it.
- `agent/internal/fleet/fleet.go:72–102`: `Spawn` copies the whole parent Spec,
  sets owner-relative child paths, calls the application factory, verifies
  back-pointers, loads/saves, inserts the child, and uses its normal actor.
- `agent/ensemble/ensemble.go:274–320`: `BuiltinTools=false` still installs all
  fleet tools if fleet is enabled. Skills can install lifecycle tools, load
  primary/startup instructions and launch MCP peers. Recall defaults on.
- `agent/internal/agent/agent.go:217–242,274–281`: admission and declarations
  both use `allowedTool`/`toolEnabled`. SafeMode removes shell/supervision tools,
  not file reading, all management tools, memory or recall.

Smallest clear change: separate default child-construction availability from
public delegation-tool enablement. Preserve `SetAgentFactory` callers' existing
fleet behavior and honor a configured custom factory. Do not call `EnableFleet`
merely to make internal compaction work. A headless `NewAgent` client must retain
its existing public tool surface, including `BuiltinTools=false` clients.

Add a creation-only managed-compaction purpose alongside the optional submission
contract in the shared child specification; exact field names are the coder's
choice. Normalize that purpose in `NewChildAgent` before skill/MCP side effects.
It means fresh history, a dedicated role, no primary/startup skills or variables,
no recall, all memory bands disabled, no recursive compaction, and only `submit`
admitted by declaration **and execution**. The worker's internal journal/result
writes still use its inherited Sandbox; no filesystem tool is granted to the model.
Keep SafeMode, web and sandbox-root ceilings; if the child data path cannot be
constructed under them, fail visibly and retain source.

Do not copy the parent's submission schema or managed purpose into ordinary
grandchildren. Existing Spawn copies the entire Spec, so inheritance requires an
explicit reset. Ordinary children still restore same-ID histories; compaction
attempts need fresh, uniquely correlated child identities, including after
restart. Use the owning Agent's configured DataDir, never a deployment path.

`internal/agent/context.go:19–26` runs forcing for every Agent. A restricted worker
must also suppress checkpoint forcing, otherwise its sole submit capability can
be hidden by `forceHandoff` after a large response. Disabling memory settings alone
does not disable that code. Preserve model-specific visible-explanation enforcement.

## 2. Submission contract and persistence

All current submission boundaries are string-only: `common.ChildStatus.Result`,
`SubAgentManager.Submit(string)`, `Submission(string) (string,error)`, the manager's
`result`, and `submit.md` JSON lines containing `{"result":"..."}`. Join encodes
ChildStatus values; `get_submission` encodes another string result. Sources:
`internal/common/fleet.go`, `internal/fleet/fleet.go:330–356`,
`internal/fleet/management.go:47–70`, and `internal/tools/management.go:80–89`.

Recommended compatible boundary:

1. Keep old Spawn/Submit/string retrieval entry points as wrappers. Add an explicit
   options/raw-JSON path rather than changing every caller's return type to `any`.
2. `submit_schema` describes the complete direct argument object of `submit`.
   Absence retains the existing `{"result":"text"}` behavior. Custom compaction
   uses exactly `{"memory":"..."}` without an extra wrapper.
3. Compile/validate and defensively copy the schema before factory invocation.
   Declare that canonical contract and validate raw arguments before disk output,
   manager completion or join notification. Old `Submit(string)` must traverse the
   same validation boundary; it cannot bypass a custom schema.
4. Preserve custom payloads as owned `json.RawMessage` or an equivalent JSON value
   in status, joins, retrieval and durable records. An additional explicit payload
   field can preserve the old `Result string` surface. Do not put JSON into that
   string and call it structured support. Copy incoming and outgoing byte slices.
5. Make custom durable records distinguishable from legacy lines and retain enough
   contract identity to interpret them on reload. The exact record envelope is a
   Stage B design choice. Continue reading old `submit.md` records. A historic
   fallback submission is not proof that a newly started task has completed.

No genuine general schema validator was found in the current core/GUI dependency
lists or implementation. `RegisterTool` and MCP discovery use `json.Valid`.
`internal/tools/tools.go:77–93` checks required presence/null and Go decoding;
it ignores extra properties and does not evaluate nested schema semantics.
Metadata is insufficient. Introduce one shared validation boundary in a responsible
implementation package, reached through shared contracts, rather than separate
submit/compressor validators or behavior in `common`.

Before coding, freeze the supported dialect/keywords, root-object requirement,
null/type/integer/array/nested/enum semantics, limits and reference policy. A
well-defined supported dialect with explicit rejection of unsupported constructs
is acceptable; a silent partial validator is not. Prefer a reviewed validator
dependency over a handwritten imitation of general JSON Schema. Disable external
reference resolution; reference support and complexity limits must be tested at
spawn as well as submit. No particular dependency was researched or selected in
this source-only audit.

The renderers currently pass raw schemas through: `anthropic.go:122`,
`openai.go:142`, `gemini.go:133`, `responses.go:99` (Responses sets `strict:false`).
Wire-schema tests must cover all supported routes, including the exact memory
schema, without relying on provider enforcement. General structured submission
does not require converting every ordinary tool result into structured data.

## 3. Submission is earlier than quiescence

`fleet.Submit` writes/closes `submit.md`, then sets `submitted` and signals the
parent from inside the running tool handler. `Join` returns on that flag.
`agent/dispatch.go:13–23` waits for the tool batch to drain, records `turn_ended`,
then finishes. `actor.go:305–340` marks Fleet/activity/TurnEnded **before** its
final log save; only the reliable request Result contains that save error.
`Spawn` currently discards the Result value and retains only its notification.
Status can still say submitted after a failed request save.

Use the already-created reliable request completion for managed tasks; retain its
error. Never accept a status transition or lossy observation as the completion
barrier. Successful completion must include the submit tail and save, then stop
and join the child and checkpoint it before dismissing references. `Flush` waits
for turns and saves, not all background jobs; `Shutdown` joins workers/jobs;
`Save` checkpoints after shutdown. Existing `GracefulShutdown` supplies much of
this ordering, but alone cannot recover a prior Result error that Spawn discarded.

Make managed cleanup idempotent against explicit kill and owner/application
shutdown, using retained child identity rather than a reused string ID. Record
cleanup/save errors rather than treating removal as success. The parent needs a
terminal outcome and child trace location after dismissal.

There are two concrete retention leaks to close for repeated workers:

- `ensemble.go:322–324` appends every constructed Agent to `e.agents`; child
  shutdown removes only the manager map (`fleet.go:196–211`). There is no matching
  application removal. Also audit construction failures after root registration.
  The default DataDir uses `len(e.agents)+1`, so removal must not introduce
  accidental identity reuse.
- `gui/gui.go:98–105` calls `view(a)` for every observation. `gui/tree.go:21–27`
  retains the Agent and copied events forever, even if never selected. Tree refresh
  does not prune views. Dismissal must retire these references and handle browsers
  selecting that child, through the existing observer/tree machinery.

Do not add an archive browser or a new general registry for this task. Preserve
trace files and parent outcomes while removing runtime ownership references.

## 4. Parent lifetime and finite worker attempts

Current compression is tied to `compressed.turn`, uses the active request's
context, and drops completion when that pointer is no longer active
(`internal/agent/memory.go:15–20,76–100`). `checkpoint` sends the next main request
only after compression/cascade completion. A background task must instead live
on the runtime/service with one slot for an in-flight or commit-pending attempt.
Workers return immutable messages to the existing owner mailbox.

Recommended lifecycle: selected → launched → candidate/failed/canceled → committed
or invalidated → cleanup complete. These are task facts, not a parallel Agent
execution loop. The actor admits completion only at idle or before a new render
with no incomplete main tool batch. A candidate arriving during an HTTP request
waits; do not mutate context so that the in-flight request appears to have contained
it. Unrelated appended dialogue does not invalidate selected completed work.

Recommended cancellation contract: normal parent-turn completion **and parent-turn
interrupt** leave a useful compaction child alive, matching ordinary background
children. Explicit worker cancellation/kill, owner shutdown, and relevant memory
disable invalidate the task and stop its wasted work. The plan §7 already prefers
ordinary background-child semantics; adopting it appears within the authorized
handoff, not a new Bill approval gate. Record this deliberate replacement of
Chapter 15 and `TestMemoryInterruptedCompressor`, rather than quietly deleting a
failing old expectation. If Bill explicitly intended interrupt to cancel all
owner compaction, that contrary instruction would take precedence.

Ordinary child Idle supplies three extra corrective wakes (`fleet.go:369–385`),
and MaxToolRounds bounds tool batches, not all model-only replies. Managed workers
need an explicit finite attempt policy and no automatic idle wakes or immediate
same-source relaunch after failure. A small existing round budget can allow a
bounded correction after invalid submit; document actual maximum requests as well
as dispatched batches. A text-only finish without submit is a failed compaction.
No new wall-clock model deadline is needed. A later explicit checkpoint may retry.

## 5. Selection, generations and safe acceptance

Keep selection in Memory. `memory.go:190–235` selects oldest whole completed
checkpoint segments, including handoff notes and role labels, excluding independent
skills/declarations and the unfinished tail. Tool parts are removed by the existing
handoff path; source building keeps readable text. Graduation (`:239–295`) uses
canonical files and watermarks, at least two files, and includes existing MEMORY.md
when producing the terminal band. Keep those semantics.

Capture task identity, policy revision, selected boundaries/identities and exact
selected bytes when choosing work. At commit, verify those selected entries/files
and relevant source/destination settings. Do not rerun selection against the entire
now-growing history or hash all history. A monotonic policy revision must remember
disable/re-enable even if final settings equal the old values; comparing final
booleans alone permits stale results. Do not tie it to unrelated GUI preferences.

Disk revalidation matters even when context has not rescanned a human edit.
Fingerprint the selected graduation files and the existing terminal MEMORY.md
input. Protect a newly allocated output identity from collision with a file created
while the worker ran. A terminal summary must not overwrite a human's newer curated
memory. Reject whitespace-only memory separately from schema validation. Validate
at the owner commit boundary, not merely at child submit.

`memory.Scan` parses range endpoints but does not reject `Thru < File`;
`Graduation` trusts the last selected entry's Thru, while History retirement compares
that endpoint with each entry's File. This is a source-supported possible no-progress
fold. Establish it with a focused failing test before changing it. It is not evidence
for the reported live Conversation compressor failure. A fold must retire selected
work rather than schedule the same source indefinitely.

## 6. Disk/event ordering needs crash-window coverage

`Memory.Write` (`memory.go:299–310`) installs the replacement at its active `.md`
path before the actor appends cut/retirement and population. Archive follows those
events (`agent/memory.go:105–135`). History holds proposed removals until matching
band/Thru population (`history/reduce.go:82–105`), preserving source in pure replay.
That pairing should remain.

Two separate disk/replay windows are not solved by the existing pair:

- Crash after the active replacement file is written but before population:
  restart Scan can inject that uncommitted output alongside retained source.
- Crash after population but before Archive: restart Scan can reintroduce old
  source files whose entries the completed event pair retired.

The current abandonment test restores a launch/cut with no matching population;
the held-request test interrupts before output is written. Neither establishes
these disk windows. Add narrow Memory-owned staging/recovery for its own attempts,
with enough task metadata to reconcile completed versus abandoned work before
running startup scans. Preserve replacement bytes before retiring source and keep
staged/uncommitted outputs out of ordinary scans. Do not build a general transaction
framework, and do not claim paired history facts alone make corpus files atomic.

Extend launch/outcome correlation with a task identity; current `(band,Thru)` can
repeat on abandoned retries. Record failure/cancellation/invalidation/abandonment
and cleanup without changing the dialogue turn state. `error_occurred` currently
sets History's turn idle, so it is unsuitable as an unqualified background-outcome
event while the main request is running. Replay must apply only recorded decisions,
clear abandoned pending work and never spawn or infer a new model result.

## 7. Roles, sizing, and accounting

Current `Engine.Compress` constructs fresh dialogue but calls the same Engine's
Render, which uses the owning main SystemPrompt (`engine/compress.go:14–22`,
`engine.go:405–425`). A dedicated child fixes that role-source conflict only if
construction also removes the inherited primary. The ordinary child instructions
mention tools a restricted child lacks, so supply a coherent dedicated purpose.
Ask for the **source owner's** first-person memory, with delimited attributed data,
decisions and reasons, failed approaches, uncertainty and constraints. Retain the
ratio as guidance and exclude visible submit narration from the memory payload.

Sizing is not fully settled by using the normal Engine: Anthropic currently caps
output at 1,024 or 4,096 tokens (`anthropic.go:133–136`); Gemini at 4,096
(`gemini.go:142`); Responses has no general configured output cap here, and the
subscription capability table omits `max_output_tokens`. Do not equate byte budgets
with tokens or add a tiny universal cap. The coder should justify supported-route
sizing and test truncation/failure retains source; this audit selects no new cap.

The direct compressor currently contributes once to the parent's Engine totals
through `send/account`, without updating Engine.last or History.LastUsage.
`agent/memory.go:142` forces from History.LastUsage only. A real child naturally
accounts to its own Engine. Preserve that owner and the main-only forcing numerator;
never append child `response_ended` events to the parent's dialogue to transfer cost.

Freeze the reporting choice: parent-only meter plus separately retained worker
usage is the smallest change; if the existing owner total must still include
compaction, use an explicit auxiliary usage record/import once at completion,
including vendor/model/tier and failed attempts with measured usage. Do not add both
that import and the same child again in a fleet aggregate. Do not overwrite the
main last-request meter or restore historic usage as new process spending. Parent
outcomes should retain worker usage before its runtime is dismissed. The coordinator
owns the accounting/chapter-placement decision; this audit does not silently change
the meaning of the existing session total.

## 8. Tests and scoped file ownership

Useful existing real-runtime coverage to retain/adapt:

- `ensemble/memory_test.go`: three-provider compression, forcing, disk toggles,
  waterfall, replay abandonment, failure retention, whole segments, enabled restart
  edits and interrupted compressor. Its HTTP fake holds a global mutex through the
  response callback and uses fixed request ordinals; a held compressor can therefore
  block the fake's main request independently of runtime behavior. Use per-request
  barriers and concurrent-safe classification for new background-progress checks.
- `ensemble/fleet_test.go`: hierarchy/questions, main-vs-child role, lifecycle and
  same-ID restore, bounded idle correction. `ensemble/management_test.go` explicitly
  Flushes after Join, confirming those are different boundaries.
- `ensemble/actor_test.go`: reliable completions, save failure, interruption and
  shutdown. `ensemble/visible_reasoning_test.go`: real actor refusal, stream order,
  bounded correction. `gui/selection_test.go`: submitted turn and replay state.

New checks need real actor/engine/fleet with only the provider faked: valid/invalid
structured submission including old direct API attempts; schema/result buffer
mutation; declaration and forged-call rejection; main progress behind a child
barrier; completion after parent idle/interrupt; explicit cancellation; selected
source and terminal-destination edits; disable/re-enable; no overlap/relaunch;
save failure after submit; crash points around disk/events/archive; cleanup races;
repeated registry/GUI retention; and accounting exactly once. Test the intended
failure with targeted mutations, preserving the existing source-preservation checks.

Expected owned source areas, subject to coder's smallest implementation:

| Boundary | Existing paths |
| --- | --- |
| Shared submission/task vocabulary | `agent/internal/common/{skills,fleet,memory,types}.go` |
| Factory, construction, detach | `agent/ensemble/{fleet,ensemble,registry,spec}.go` |
| Submit tools and fleet lifecycle | `agent/internal/tools/{fleet,management}.go`, `agent/internal/fleet/{fleet,management}.go` |
| Parent scheduling/admission | `agent/internal/agent/{actor,memory,dispatch,context,agent}.go` |
| Corpus and replay | `agent/internal/memory/memory.go`, `agent/internal/history/{history,reduce}.go` |
| Role/sizing/accounting only as needed | `agent/internal/engine/{compress,engine,usage,anthropic,gemini,responses}.go` |
| Visible dismissal and structured display | `agent/gui/{gui,tree}.go` and affected renderer/tests |

`ensemble.go`, `registry.go`, and `common/skills.go` are announced shared hotspots
with Cadence. The audit also establishes a concrete need to coordinate GUI view
cleanup. No evidence here requires editing live configuration/startup files or
Cadence's memory repository. Do not claim exclusive ownership from this table.

Stage A can proceed to the coordinator's contract/mapping decision. Before Stage B/C
implementation, record schema semantics, structured persistence shape, finite worker
attempt bounds, task acceptance/recovery protocol and accounting attribution.
The model/factory/lifecycle gaps above are engineering work under the handoff;
the final chapter move and historical snapshot policy remain Bill's stated gates.
