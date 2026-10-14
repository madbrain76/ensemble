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

## Stage B/C independent review resumed

2026-10-14, reviewer `compaction_runtime_audit`, checkpoint `7a362ec2`.
The coordinator reports Bill approved the recommended mapping, staged accounting
and reviewed replacement exports. Those Stage A decision gates are resolved;
the original audit above retains its then-current questions. Accounting will
preserve the existing owner's session-total promise by attributing measured
auxiliary usage exactly once, while keeping latest main-request usage separate.
Stage C must retain the actual worker/model/tier identity and avoid re-importing
historical spending on restart.

After compaction the reviewer reloaded the entire mandatory
`.agents/skills/ensemble-coding/SKILL.md`, workflow, carryover, voice and chapter
procedure, recovered truncated read ranges, and checked the clean worktree at
the checkpoint. Current Chapters 22 and 23 were also read in full for the
submission and fleet contracts. The original full Chapters 14/15 reading is
recorded above. No implementation, test or acceptance result is claimed by this
preparation. The Stage B coder is separate from this reviewer; its concrete
schema/API/storage/ownership design is awaiting review before implementation.

### Stage B design acceptance

The independent reviewer accepts the design in
[`codex-structured-submissions.md`](codex-structured-submissions.md), including
the coder's subsequent `SubmissionData` refinement, for implementation. This is
a design gate, not formal implementation acceptance; zero formal code-review
rounds have occurred. The coordinator and coder were notified directly.

The accepted shape preserves `Spawn`, `Submit(string)`, legacy `Result` fields
and `{"result":"..."}` durable rows. Custom contracts cover the direct submit
argument object, use an owned raw-JSON `Payload`, and store an explicit version-1
schema/payload envelope. `SubmissionData` retains default/custom identity even
when a custom payload happens to look like a legacy result. Legacy string retrieval
reports a clear error for custom results. The existing string-valued `Ask` and
`ask_agent` completion can render custom JSON as text; structured callers keep
the raw payload. Durable fallback never marks a newly spawned task complete.

The pinned genuine validator is `github.com/santhosh-tekuri/jsonschema/v6`
`v6.0.3`, using an explicit bounded Draft 2020-12 subset. All references and
unsupported keywords are rejected; a refusing loader provides a second boundary
against resource fetching. Local dependency-source inspection confirmed the
draft/loader/resource APIs, exact-number parser and standard Go regexp default.
No online research or dependency mutation was performed by the reviewer.
The proposed size/depth/node/numeric ceilings are explicit admission limits,
not truncation or summarization targets. Their implementation and tests still
need review; declaring a limit alone does not prove bounded evaluation.

The Fleet owns the immutable compiled contract. Construction validates/copies
before side effects, Spawn validates before a custom factory, and the returned
child must retain the requested contract. Every submit entry point uses that
contract before persistence and completion; copies also protect outgoing Spec,
declaration, status and retrieval buffers. Ordinary grandchildren reset the
parent's custom schema.

Two narrow clarifications accompany acceptance: absent-schema submission retains
the old acceptance of extra keys and normalizes to its result text/legacy row;
custom submission retains the full object. The keyword admission walk visits
schema-bearing positions only, not object keys inside enum, const, default or
examples data. These clarifications prevent compatibility loss and accidental
rejection of valid data without adding another validator or result subsystem.

The coder's planned real-runtime/provider-fake checks cover corrected invalid
submissions, pending joins, public API enforcement, defensive copying, hierarchy,
ownership, durable retrieval/reload, adapter declarations and GUI replay. Review
will assess the actual assertions and targeted mutation outcomes. No new tests
have been run by this reviewer at this design boundary. Stage C's quiescence,
restricted construction, background scheduling and accounting remain separate.

### Stage B pre-submission source review

The coordinator resumed production review while the coder completed its report,
with the full governing context still active. The worktree had advanced to
`a8eec84a`; implementation was uncommitted. This review does not consume a formal
correction round. Findings were sent directly to coder and coordinator before
the consolidated submission:

- GUI result assertions initially searched for marker words instead of decoding
  the actual submitted object. They could not establish the claimed nested-data
  preservation. The coder is strengthening both live delivery and replay checks.
- Separate schema/input ceilings did not control their multiplied validation
  work. Dependency source collects errors across all array elements and retains
  each element's missing required names; formatting those errors can amplify a
  small input. No large or small executable amplification probe was run before
  this finding. The agreed correction adds an explicit conservative combined
  byte budget before validation, without replacing schema semantics or claiming
  a CPU deadline.
- Incoming byte caps preceded JSON's HTML escaping. An admitted annotation or
  payload containing many `<` characters could be serialized above the same cap
  enforced during durable rereading. The accepted shape would then reject its
  own stored record. This was initially source-identified, not reproduced; the
  coder is adding canonical-byte refusal checks and small restored-defect cases.
- Moving factory checks into the public `CreateChild` path is appropriate. Its
  inherited wrong-owner branch must not shut down an unrelated returned Agent.
  In particular, returning the caller itself from a factory invoked by Spawn
  would re-enter the manager lock during shutdown. Ownership rejection and
  disposal of a newly created, correctly owned wrong-contract child are distinct.
  The coder received this narrow finding and a real-Agent survival test suggestion.

A host disk-full event temporarily paused builds/tests. The reviewer had no
active process and ran no builds during that pause; the coordinator cleared
regenerable caches and subsequently released bounded checks. Source analysis
and later executable evidence retain their distinct statuses.

### Stage B formal implementation review: accepted

October 14, 2026. Independent reviewer `compaction_runtime_audit` accepts Stage B
at initial formal review, with zero formal correction rounds. All pre-submission
findings above remain attributed assistance. After the final context compaction,
the reviewer again read the entire `.agents/skills/ensemble-coding/SKILL.md` and
all four governing/craft documents before this decision, then reread the frozen
production diff, schema/storage implementation and consolidated coder report.
HEAD is `a92aaad9`; the Stage B source is still uncommitted for the coordinator
to checkpoint. The sorted changed/new `agent/` path-and-SHA256 manifest hashes to
`516f2d8205da583393b0fc2563f1032ddaebf9d295433985db0fd4a5f08f82ce`.
All 214 files in the coordinator's grader source manifest still match disk.

No Stage B blocking finding remains. The final public `CreateChild` boundary
also refuses nil or foreign parents before invoking the factory, matching the
existing required-parent construction rule. The real-Agent tests cover this
entry point directly and through Spawn. A factory cannot discard the promised
schema; wrong ownership does not authorize shutting down a live foreign Agent
or the caller. This closes the additional pre-submission parent-boundary finding
without adding lifecycle policy outside this stage.

The implementation follows the accepted ownership plan: common holds data and
interfaces; Fleet owns the immutable compiled contract, validation and durable
submission behavior; the composition root validates before construction effects.
The genuine validator owns schema semantics, with explicit admission restrictions
and no external resource resolution. Exact-number decoding, raw and canonical
byte checks, and the combined byte budget address the identified precision,
reload and repeated-work risks. The combined budget is deliberately conservative
and does not establish a measured CPU deadline; no amplification stress run was
performed or needed for this acceptance.

Legacy string submission and rows remain compatible. Custom JSON survives
declaration, tool dispatch, status, join, typed retrieval and durable rereading
without double encoding. The explicit envelope preserves contract identity even
when custom data resembles a legacy result. Validation and successful writing
precede completion publication; historical fallback does not complete a new task.
Spec, declaration and returned payload bytes are copied. The existing text-valued
Ask/display API renders compact custom JSON, while structured readers retain its
object. GUI production code is unchanged; the new test checks full nested data
in both actual subscribed live frames and restored replay, with no inference on
reload.

Independent reviewer executions:

- `go test -race ./internal/fleet -run TestSubmission -count=1`: pass, 1.080s.
- `go test -race ./ensemble -run TestStructuredSubmission -count=1`: pass,
  1.429s after the HTTP fixture correction. The earlier failure was the fixture's
  missing Anthropic usage, and remains a failed setup rather than an ownership
  defect. After the final parent guard,
  `go test -race ./ensemble -run 'TestStructuredSubmission(Construction|FactoryOwnership)$' -count=1`
  passes, 1.036s.
- GUI `go test -race . -run TestStructuredSubmissionGUIReplay -count=1`:
  pass, 1.079s with the exact live/replay assertion.

The consolidated coder evidence supplies passing final focused Stage B race
checks, full GUI race (11.856s), both module vets and full ensemble race
(91.778s). That full ensemble run began before the final required-parent guard;
the final focused race run covers the added guard. Full core short/race is not
claimed green: its CLI human-sink case reproduces on baseline. The coder's
standalone browser glob was an invalid suite invocation, and its separate speech
persona failure reproduces on exact `a8eec84a`; neither is silently passed or
repaired in this stage. ENOSPC link failures and the initial missing GUI
subscription are retained as setup failures. No live provider, deployment or
publication result is claimed.

The coordinator's advertised Edition 2 compatibility commands using current
chapter numbers and explicit `./agent` score Chapter 22 200/200 and Chapter 23
100/100, both exit 0. The reviewer inspected their JSON receipts and verified
the complete source manifest against disk. These are the revised compatibility
graders; historical original scores are unchanged. Receipts are under
`/tmp/ensemble-compaction-baseline/grader-{22,23}-stageb*`, with
`stageb-agent-source.json` binding the source.

Mutation metadata and failure excerpts were inspected in
`/tmp/ensemble-structured-submissions-mutations/`. They cover the reported 21
distinct defects over 24 executions, including repeated location checks and the
factory recheck. Canonical schema/payload and combined-budget omissions each
admit the prohibited value; restored foreign-agent shutdown fails with
`agent stopped`; removing parent admission invokes the factory; double encoding
fails the actual nested GUI result comparison. Raw-byte corruption is caught on
decode. The earlier factory run lacked a valid baseline until its peer supplied
usage, so the repaired baseline/recheck supplies that claim. Declaration aliasing
also causes cleanup serialization errors, and several schema-semantic assertions
overlap. These are sensitivity observations, not 24 independent proofs or proof
of every upstream keyword. The tests use real runtime objects with small external
HTTP peers, without mocks or a new test framework.

The independent size count agrees with the coder. Across 13 changed/new production
Go files, physical/active/comment-only lines move from 2212/1701/332 to
2596/2007/388, a net 384/306/56. Three new test files contain 764 physical,
712 active and 33 comment-only lines. Active counts exclude blanks and
comment-only lines; moved submission code is counted on both sides, avoiding a
false wholly-new-code claim. Module dependency files are separate. This compares
the maintenance change against its own `a8eec84a` predecessor; there is no
first-edition structured-submission implementation to claim as its counterpart.

The growth implements the requested contract and storage boundary rather than a
new scheduler, result browser or schema framework. Useful comments explain
ownership, copies, refusal versus semantic validation and the write/completion
ordering. Exported API additions document their contracts. Comment-only lines
are about 16% of nonblank production lines in this changed-file set, below the
rough 20% teaching target but with the important new reasoning present; padding
would not improve the code. Tests are longer because they exercise actual
provider adapters, factory ownership, persistence and GUI delivery, not because
they reconstruct the implementation behind mocks. No scope or KISS issue remains.

Stage C still requires its concrete design gate. Stage B's submitted state is
not worker quiescence: retain reliable turn completion through the actor tail
and save before cleanup, and remove owned root/GUI references when that lifetime
ends. The earlier disk/event findings now have coordinator-run controlled
saved-state boundary/restart reproductions on exact `a8eec84a`: uncommitted output
is reinjected before population, and retired sources return before Archive.
No-write/completed-Archive controls pass. These are not forced process kills;
the probe and logs are in `/tmp/ensemble-compaction-window-probe/`, probe SHA256
`ef869c4e0945e84c763c837dfd81f7a9f1ca85941b4aa4139d1418b5c5fae2af`.
Acceptance of structured submissions does not claim those memory/lifecycle
problems are fixed, authorize chapter moves now, or satisfy later live gates.

### Stage C concrete design: accepted for implementation

October 14, 2026. The independent reviewer accepts the revised
[`codex-background-compaction.md`](codex-background-compaction.md) for Stage C
implementation. Stage B is checkpointed at `a0acfb5e`; subsequent Cadence note
merges and accounting notes do not change that runtime baseline. The entire
mandatory `.agents/skills/ensemble-coding/SKILL.md` and four governing/craft
documents remain loaded from the final B review, with no intervening compaction.
The reviewer read both complete versions of the concrete C design and refreshed
the actual Memory, actor completion, Engine usage, factory and admission paths.
This is a design acceptance only; no C implementation, formal correction round,
runtime test or recovery success is claimed.

The proposed one owner-lifetime attempt reuses the existing actor, Fleet child,
reliable request Result and mailbox. Ordinary turn completion/interruption leaves
it useful; explicit kill, owner shutdown and relevant policy changes cancel or
invalidate it. Commit waits for a safe parent boundary and exact selected-source,
destination and policy-generation checks. Unrelated appended dialogue does not
invalidate the result. A known 95%-full main context may wait for that outstanding
attempt; failure produces an explicit error rather than another full-context send.
This is a stated exception to ordinary background progress, not a model timeout.

Restricted purpose is enforced in construction, declaration and dispatch, with
no copied primary/skills/MCP, recall, memory recursion or arbitrary tools. It grants
only the needed submit operation to the worker. Source inspection confirms
`BuiltinTools=false` currently selects default registration rather than an inherited
security ceiling; keeping the parent's declarations and FleetEnabled unchanged
honors that contract. Actual SafeMode, web and sandbox ceilings remain binding.
The revised design accurately distinguishes normalized input and checked public
construction from an arbitrary trusted custom factory's own pre-return effects.
The root cannot retrospectively prevent those effects by rejecting its result.

The finite policy is three Engine.Send invocations, including correction/resend
paths, with existing bounded authentication/HTTP retries explicitly separate.
Prose-only completion fails and the purpose cannot acquire another task to reset
its budget. The reviewed provider output policy is deliberate and worker-specific;
its live acceptance and semantic quality remain Stage F obligations. Following
review feedback, adapters must detect actual length/incomplete terminal reasons
before dispatch even when a submit's JSON is syntactically valid, while retaining
measured usage. Malformed JSON alone would not test that distinct failure.

The single hidden candidate plus correlated journal receipt is justified by the
two reproduced disk/event windows. Memory owns file staging, promotion, source
archive and recovery; History atomically projects the matched cut/population and
retains the receipt through snapshots. A latest terminal outcome recognizes the
final cleanup window, and no new attempt starts before staging removal succeeds.
Recovery precedes startup/re-enable population and recall. After journal commit,
an unfinished disk finalization is a recoverable committed state, never abandonment.
The clarified hold blocks main rendering and recall as well as corpus scans until
finalization succeeds or the turn receives an explicit error. Human-edit conflicts
preserve copies and receipt, rather than overwriting them. The design claims
process-crash ordering, not power-loss durability or transactions against external
editors. The implementation must still prove each described boundary.

Completion, cancellation and accounting share a real child drain and save attempt
before reference removal. The revised accounting proposal uses the active attempt's
identity and accounted bit instead of an unbounded Engine task-ID deduplication set.
Engine remains owner of measured model/tier buckets and adds final auxiliary usage
once, including failed/canceled work, without replacing latest-main usage or replaying
historical spending into a fresh session. Exact-child shared dismissal prevents
manual kill races from removing a replacement. Root and GUI references must go
with the Fleet entry; retained files provide the audit trail.

The reversed-range issue now has a separate coordinator baseline reproduction:
`ranges_test.go` under `/tmp/ensemble-compaction-window-probe/`, SHA256
`a513c5c1fde6fa7bbd26ddc7ea0924394117f0d9264b41a421784be4056b27dd`.
Malformed coverage retires nothing and immediately selects the same source; valid
filenames pass the control. C's range admission and exact retirement-prefix checks
address that demonstrated progress obligation, without claiming it caused Cadence's
earlier conversation-compressor failure.

Validation must hold the external request outside capture locks. The coordinator's
separate real-HTTP probe reproduces the shared fakevendor Route lock blocking main
requests and Requests(); it is Stage D fixture evidence, not a compaction result.
The existing Chapter 23 reverse-proxy barrier is a narrow later grader option.
C's independent concurrent HTTP peer can exercise actual main progress now. Tests
remain real runtime plus external fakes, with targeted mutations and genuine
filesystem failures; no new persistence mock or generic task framework is needed.
No unresolved design blocker remains. Production complexity, concrete ownership,
boundary tests and required compatibility results still face independent code
review before Stage C acceptance.
