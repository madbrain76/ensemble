# Stage C: restricted background compaction design

Coder: `/root/structured_submission_coder`, 2026-10-14. Design accepted by the
coordinator and independent runtime reviewer on October 14; implementation is
released in coherent slices and remains subject to separate review. Stage B is
independently accepted at `a0acfb5e`; design began at `e8ddced4` and is checkpointed
after documentation-only `cd798d36`. This is maintenance with prior-answer
exposure, not a blind student run. Mandatory skill:
[`.agents/skills/ensemble-coding/SKILL.md`](../../.agents/skills/ensemble-coding/SKILL.md).
The full skill and carryover were reloaded after compaction; this design also
uses the complete workflow, current chapters 14–16 and 22–23, complete handoff,
runtime audit and progress record. No chapter, grader, export or deployment edits.

## Ownership and public seams

Keep the actor, Fleet, Memory, History and Engine as their existing owners.
Shared immutable records/interfaces belong in common; implementation spokes
continue importing common rather than each other. No new scheduler or store.

- A creation-only `AgentSpec.Compaction` purpose selects the restricted worker.
  The ordinary `Spawn` path resets it, as it resets `SubmitSchema`. Root validates
  and normalizes the specification before invoking the custom factory, and the
  normal constructor enforces it before skills, primary prompts, MCP or recall can
  initialize. A public custom factory must return an actually restricted child
  with the requested owners/schema/purpose; the checked constructor refuses a
  mismatch. Safety ceilings still clamp normally. A custom Go factory is trusted
  application code: the root can normalize its input and verify its returned child,
  but cannot prevent arbitrary side effects performed inside that callback before
  return. Test the normal constructor's pre-effect restriction and reject a
  nonconforming returned child; do not claim the callback itself is sandboxed.
- `Fleet.StartCompaction(id, task)` constructs one fresh child and returns its
  actual Agent and reliable `Result` channel. It shares ordinary child creation,
  registration, paths and observation, but does not discard Result or reload a
  previous child task. A random path-safe task ID also names the child. Source,
  policy and commit metadata remain with the parent, never model-authored.
- Root supplies its ordinary default factory when no custom factory is set;
  this does not set the public `FleetEnabled` flag. `EnableFleet` continues to
  enable public delegation and preserves an already configured custom factory.
  Internal construction must work without granting delegation to headless or
  `BuiltinTools=false` parents. `BuiltinTools` controls default registration, not
  an inherited security ceiling: register only submit in this special child while
  leaving the parent's declarations and FleetEnabled unchanged. Preserve actual
  AgentSpec safety ceilings (SafeMode, web and sandbox narrowing). If those ceilings
  prevent required child construction or trace persistence, fail explicitly and
  retain source; do not widen them to make internal delegation work.
- Memory extends its existing selection value with selected-entry identities,
  relevant policy and file/destination fingerprints. Its narrow new operations
  are `ValidateSelection`, `Stage`, `Finish` and `Recover`; they own corpus I/O.
  History projects correlated facts. Agent owns when facts may be appended.
- Fleet dismissal stops/saves an exact child identity and removes that identity
  from its map and root registry. Root issues a removal observation; GUI handles
  it before creating a view, drops the retained view, and switches selected tabs
  to the owner with the ordinary replay path. Directory allocation uses a root
  monotonic counter, rather than the now-shrinking registry length.

## Restricted construction and finite execution

The child inherits provider/model/route and safety settings, not the parent's
role or conversation. Before any startup side effect, clear primary/startup
skills and variables, bypass skill/MCP initialization, disable recall and all
memory bands, and bypass context forcing. Give it only the normal submit tool,
with Stage B's exact whole-argument schema:

```json
{"type":"object","properties":{"memory":{"type":"string"}},"required":["memory"],"additionalProperties":false}
```

Both declaration and actual execution admission accept only `submit`, including
guessed names and later registration attempts. These restrictions are immutable
for this purpose; settings cannot reopen tools/recall/memory. Internal persistence
still uses the normal sandbox and does not become a model filesystem capability.

The dedicated prompt asks for the **source owner's first-person memory**, retaining
decisions and reasons, failed approaches, uncertainty, constraints and preferences.
The task message contains delimited source data with role labels and handoff notes;
quoted instructions remain data. There is no main coding/no-submit prompt, cloned
dialogue or opaque model state. Visible explanations precede submit where current
model policy requires them; they stay outside the memory string.

One attempt permits at most **three Engine.Send invocations and three tool batches**:
the first request plus two opportunities to correct validation/visible-explanation
errors. Count invocations immediately before actual sends, including ephemera-induced
resends, across the worker's entire lifetime. Do not use ordinary missing-submit
idle wakes. A completed prose-only response fails immediately. More than one call
in a response, or a call other than submit, rejects that batch with paired errors;
no first-call partial success. A successful submit finishes after its tool-return
tail. The purpose admits one task only; messages cannot start fresh turns or reset
its budget. An invocation can include the existing bounded authentication refresh
and HTTP retry; three invocations is not a claim of three physical HTTP exchanges.
Keep that separate retry bound and test that it cannot restart the worker budget.
Failures do not automatically start another attempt on the same source.
There is no total model deadline; explicit cancellation reaches the provider.

For worker requests only, omit the optional Chat `max_completion_tokens` and
Gemini `maxOutputTokens` caps, preserving existing Responses route behavior.
Their APIs explicitly make these fields optional; Gemini's default varies by model.
Anthropic requires `max_tokens`: request **16,384** for workers, a practical allowance
over twice the observed 7,317-token failed compression, rather than the current
1,024/4,096 chat allowance. This is an execution budget, not a claimed universal
model maximum or a conversion from file bytes. Sonnet-5's documented output capacity
exceeds it. A provider/model rejecting that allowance fails visibly with source
retained; there is no silent model substitution, cap reduction or beta feature.
Ordinary requests keep their existing settings. For workers, adapters must preserve
and check actual terminal reasons in both streamed and nonstreamed responses:
Anthropic `max_tokens`, Chat `length`, Gemini `MAX_TOKENS`, and Responses incomplete
output are failures even when the returned submit arguments are syntactically valid.
Propagate the failure before any submit tool dispatch; still retain final measured
usage. Stage B's payload bounds also apply. No byte truncation forces a ratio.
Test a memory over the old caps, malformed JSON, and a valid submit paired with a
length stop; live quality/sizing evidence remains Stage F.

Sources checked October 14: [Anthropic Messages](https://platform.claude.com/docs/en/api/messages/create),
[Sonnet-5](https://platform.claude.com/docs/en/models/sonnet-5/overview),
[Chat Completions](https://developers.openai.com/api/reference/resources/chat/subresources/completions/methods/create),
[Gemini GenerationConfig](https://ai.google.dev/api/generate-content#generationconfig).
These establish parameter contracts, not new model support; Bill's excluded 5.5
models remain excluded. Existing synthetic course fixtures are not live evidence.

## Parent lifetime, scheduling and acceptance

The parent runtime holds one private attempt slot, independent of its ordinary
turn: ID, immutable selection, generation, child, cancellation state, and eventual
completion. Launch uses a background owner context, not `active.ctx`. Record launch
correlation, source boundaries and child trace path. Main rendering continues after
checkpoint cleanup while the child runs. Repeated handoffs cannot overlap this slot.

Ordinary turn finish and parent interruption preserve the child. Explicit child
kill, owner shutdown or a relevant memory policy change invalidate the attempt and
cancel/drain it. Compare source/destination band settings at the actor's existing
settings transition and increment their generation on every relevant change;
disable/re-enable cannot restore an earlier generation. Unrelated UI settings and
new dialogue do not invalidate it. Worker replies carry immutable task identity;
late replies cannot adopt a newer slot. Shutdown drains the managed completion
before the parent's final save/exit, retaining accounting and failure evidence.

Keep existing 90% warning and 95% forcing arithmetic based solely on last completed
main input, with cache buckets counted once. If an attempt is outstanding and that
last main request reached 95% of a known context window, defer the next main render
until this attempt resolves; show waiting activity and keep interruption/settings
usable. Successful compaction resumes it. If that held attempt fails, end that main
turn with the explicit compaction/headroom error rather than immediately sending the
same full context. Unknown context windows retain the ordinary provider error path.
This is the documented near-window exception to background progress, not a timer.

Completion becomes a candidate only after the child's reliable Result, full
shutdown/drain, and required Save succeed. `Status.Submitted`/Join alone is insufficient.
The completion worker copies the validated payload and final usage, then posts one
result through the parent's existing mailbox. Commit only when idle or at next-render
with the tool batch complete, never while an already-sent response is in flight.

At that boundary, Memory revalidates exactly the selected source and destination:

- Conversation selection uses whole eligible completed checkpoint segments and
  exact selected entry/part identity; independent skills/declarations, unfinished
  tail, and newly appended dialogue are outside the hash and outside the cut.
- Graduation validates selected entries against current file bytes, not only IDs
  or sizes. Hash file bytes and coverage metadata; same-name human edits invalidate.
  Validate the destination's prior bytes/absence too, including existing MEMORY.md.
- Scan rejects malformed or reversed coverage ranges. Graduation also verifies its
  retirement prefix equals the selected set and makes progress; overlapping ranges
  cannot silently retire an unselected file or repeatedly select retired-nothing work.
- Recheck relevant settings/generation. Reject empty/whitespace memory and stale
  candidates explicitly. Existing filenames are not overwritten on collision.

After a successful commit, policy may select the next graduation, one child at a
time. Do not restart failed work without a later checkpoint. A graduation must
retire its selected nonempty source set before another cascade step starts.

## Durable replacement and restart protocol

Retain paired cut/retirement and population facts. The identified baseline defects
require a hidden durable candidate and a durable commit receipt, not a new database.
Memory owns one atomic `.compaction/pending.json` record per parent, outside all
active/recall scans. It includes version, task ID, selected file paths/hashes,
destination precondition, candidate bytes/hash and population identity. All paths
derive from the owner and use its Sandbox; the model supplies none of them.

Extend the existing PendingCompression projection with TaskID and a committed
population receipt. Legacy uncorrelated events keep their current replay behavior.
For new tasks, the matching population atomically applies the held cut/retirement
and retains the receipt until disk finalization. A correlated terminal outcome
clears it and keeps only the latest task outcome in Context, so an orphan staging
record can be recognized after a crash. Do not launch another attempt until this
record has been removed successfully. This bounded receipt survives snapshots
and log retention; it is not an ever-growing job ledger. Outcome facts do not change
ordinary turn state (unlike `error_occurred`).

Before launch, establish the parent's existing snapshot/journal if the embedding
client has not called Load. Make successful History.Checkpoint establish its journal
path under the same lock; Agent.Save then supplies this boundary without replaying
or replacing live history. A persistence failure refuses launch. This also makes
plain public NewAgent/Ask memory use durable without demanding hidden client setup.

| Durable boundary | Normal action | Restart before the next action |
| --- | --- | --- |
| Launch journaled | Run fresh worker; source remains active | Mark abandoned; retain source; never respawn |
| Hidden candidate installed | Revalidate selection and append cut/retirement, then matching population | Without matching population, discard hidden candidate and abandon; source stays active |
| Cut/retirement alone journaled | Source is still in projection until matching population | Discard held cut and candidate; source stays active |
| Population journaled | History commits the pair; Memory promotes staged bytes into active destination | Recover committed receipt before any startup population or recall; complete disk work without a model |
| Destination installed | Rename selected originals to `.retired` | Verify destination hash, finish remaining expected source renames idempotently |
| All source renames complete | Journal terminal committed outcome, then remove hidden record | Receipt or latest terminal outcome recognizes installed/retired state; finish only missing steps |
| Terminal outcome journaled | Remove hidden record; release slot | Matching last outcome permits hidden-record cleanup without reapplying history |

`Memory.Recover` runs after History recovery and before **all** startup corpus
population/recall, including re-enable paths. It uses recorded bytes/decisions and
never reconstructs the replacement with an LLM. A committed receipt with missing
or mismatched staging data is an explicit recovery error; do not scan conflicting
active files into context or silently start a new attempt. A cleared staging file
with an already terminal outcome is normal. Old pending attempts retain legacy
abandonment behavior; no claim is made to repair every old ambiguous disk artifact.

Stage before appending the pair; on a pre-commit persistence failure, source remains
unchanged. After journal commit, finalization errors leave the receipt and candidate
available for retry/restart, block further compaction, corpus scans, recall and main
rendering, and report the error. Startup and re-enable must pass recovery before
any of those operations. If finalization cannot finish, return the explicit error
to the waiting/new main turn; remain interruptible and retain the recovery hold
until a successful recovery. Never send a mixture of committed projection and
unfinalized active corpus. Do not call a committed replacement an abandoned attempt. Idempotent
Finish accepts already-installed matching destination and already-retired matching
source, but never overwrites a human edit or an unrelated `.retired` file. Recheck
file preconditions immediately before each mutation. An external edit racing these
steps can require explicit conflict repair: preserve both copies and the receipt,
report failure, and stop automatic finalization. Portable filesystem operations do
not provide a transaction against arbitrary concurrent external editors. This
continues the repository's process-crash, not power-loss/fsync, durability contract.

## Completion, dismissal and usage

Capture the reliable Result, then call the existing child Shutdown and Save before
acceptance/removal. Shutdown joins actual actor/job/provider workers; even canceled
children can account usage after an interrupt Result. Do not cancel from the submit
callback. Save failures prevent candidate acceptance, retain trace files, and become
parent outcome errors. Fleet cleanup uses exact identity and a shared completion
for concurrent manual kill/automatic dismissal, so it cannot delete a replacement
child with the same name or save/remove a still-running child.

Remove Fleet/root/GUI references only after drain and required save attempt, even
on failure. Preserve journal, submit row and trace path on disk. No late observation
may recreate a dismissed GUI view. Parent outcome records completion, cancellation,
invalidation, persistence and cleanup errors; repeated dismissal is harmless.

Engine exposes a copied breakdown of its existing usage buckets (vendor, model,
large-price tier, disjoint usage) and an auxiliary addition operation. The actor's
single active attempt owns an `accounted` bit: accept completion only for that exact
attempt, add its final buckets once, then set the bit before releasing the slot.
The shared child completion handles racing kill/dismissal; duplicate mailbox
delivery cannot import twice, and a late result cannot match a subsequent attempt.
Engine retains no task-ID deduplication set. Import after actual drain, for successes
and failures, into parent session totals and corresponding model/tier buckets.
Do not alter Engine.Last,
History.LastUsage, parent response facts or latest-main context occupancy. Do not
also add child totals through Fleet aggregation. The outcome retains the usage
receipt for audit; replay does not import historical spend into fresh session
totals. Pricing remains owned by Engine, including unknown-price behavior.

The existing exported `Engine.Compress` method remains deliberately source-compatible
for external callers, with a deprecation comment explaining that it is a direct
legacy operation. The normal actor path never calls it or falls back to it.

## Implementation slices and evidence

First port root's real Memory/History restart-window and reversed-range probes into
focused source regressions, retaining their exact baseline failures and positive
controls. They model saved boundaries, not a forced process kill. Then implement
coherent reviewable pieces: restricted construction/finite execution; Memory staged
recovery and correlated projection; owner-lifetime scheduling/completion/accounting;
registry/GUI dismissal. The full protocol above governs every piece.

Use real Ensemble/actor/Engine/Fleet/Memory/History and a small concurrent HTTP peer.
Identify workers by the actual sole submit declaration with a `memory` property.
Root found the old grader classifier too broad and its fake server globally locked
around Route, so a held compressor also blocks main requests there. C tests must
release their request-capture lock before waiting at a provider barrier; graders
remain Stage D work. No mocks, live calls or production GUI connections.

Focused claims: main completes and accepts another turn while worker held; exact
role/source/schema and visible explanation; default/custom/headless construction;
no startup skill/MCP side effects and forged tool admission; finite requests despite
ephemera, prose-only and invalid submit responses; unchanged selected source with
unrelated dialogue; file/destination edits and off/on invalidation; safe tool-batch
commit; all restart boundaries and journal/write failures; no truncated/empty commit;
final submit tail and save failure; interrupt versus kill/owner shutdown; exact usage
model/tier once with unchanged main Last; repeated registry/job/GUI removal including
selected-child fallback; ordinary grandchildren and recall corpus remain compatible.

Use deterministic external barriers and real filesystem refusal cases rather than
timing guesses or mock persistence. Apply focused restored-defect mutations per
distinct claim, recording coupled failures honestly. After coherent focused checks,
run core/GUI tests, race, vet and required explicit-source old-number graders with
root; retain known baseline CLI human-sink/browser-fixture failures separately.
Strip only LLM_/OPENAI_/ANTHROPIC_/GEMINI_/EN_/CODEX_ in test subprocesses. Do not
glob browser scripts or repeatedly rerun unchanged broad suites. Semantic memory
quality, live provider acceptance and deployment remain unclaimed until Stage F.

## C1 implementation evidence: coverage admission

On design checkpoint `01dfa4ce`, source regressions first reproduced all three
root probes: uncommitted active replacement after restart, committed retirement
before archive after restart, and reversed-range no-progress reselection. All
three failed for their intended reasons in 0.015s. The initial filename
`compaction_windows_test.go` was accidentally GOOS-filtered on macOS, causing
undefined test helpers; this was setup failure, corrected before the measured run.
Raw outputs: `/tmp/ensemble-compaction-stagec-baseline{,-corrected}.txt`.

The first narrow runtime slice rejects reversed/malformed filename coverage and
checks graduation's restored entries for invalid, reversed or overlapping ranges.
It protects the prefix retirement rule without changing watermarks or selecting
different valid work. Focused race tests include a valid ranged-file scan, exact
prefix retirement/progress and inherited MemoryWaterfall; the final focused race
run passed in 1.142s; after review strengthened retained-file identity assertions,
the affected race checks passed in 1.143s. Targeted vet passed.
Restoring scan admission admits both reversed and extra-segment filenames and
fails the intended checks; restoring graduation admission fails both reversed and
overlapping restored-entry checks while the valid control passes. These are two
focused mutations, not four independent mutation proofs. A third mutation reverses
the retirement comparison: the surviving count stays one, but the strengthened
check rejects source 1 surviving instead of source 3. Scratch copies removed.
Raw outputs: `/tmp/ensemble-compaction-stagec-ranges{,-final}.txt` and
`/tmp/ensemble-compaction-stagec-mutant-{scan,graduation,prefix}.txt`; reviewed final
race output is `/tmp/ensemble-compaction-stagec-ranges-reviewed.txt`.

The two legacy Write/event/Archive restart probes are still failing and remain
outside this coverage-only checkpoint. New Stage/Finish recovery tests will test
the replacement protocol explicitly; they must not relabel these old-API failures
as unchanged passing regressions. No background scheduling or recovery acceptance
is claimed by C1. Independent source review found no production blocker and asked
for useful test comments and exact retained/selected identities; both were added
before checkpoint. The reviewer accepted C1 for checkpoint, with no production
changes requested; this is not full Stage C acceptance or a formal correction
round. Final slice disposition is recorded in the runtime audit.

## C2 implementation evidence: worker construction and finite execution

Work after `c8ec3d4b` adds creation-only Compaction purpose, normalizes its startup
specification before custom factories and in the normal constructor, replaces the
main role, disables memory/recall/forcing, and registers/admit only submit. Default
child construction no longer requires public delegation; EnableFleet preserves an
existing custom factory. Ordinary Spawn and the public AgentSpec.Child helper reset purpose as well as schema.
Ordinary Spawn still refuses when public delegation is disabled; the internal
default factory does not silently grant that public capability. The actor
admits one task, at most three Send invocations/three batches, rejects entire extra
or forged-call batches, and treats prose-only completion as failure without idle wakes.

Focused real-runtime tests cover default/custom construction with and without
builtins, normalized factory input, primary-startup refusal avoidance, unchanged
parent tools, actual safety ceilings, factory purpose mismatch, valid submission,
invalid schema correction exhaustion, extra/forged calls, a held prose response
with a pending hint, second-task refusal, and ordinary descendant defaults. This
slice's provider is a concurrent Anthropic HTTP fake with faithful usage/call IDs;
provider-specific output sizing and terminal-stop handling remain the next slice.
Tests also run the Stage B structured-submission cases. Final focused race passed
in 1.519s; affected ensemble/agent/fleet vet passed. Inherited Fleet, Named,
Management and sandbox child/safe-mode compatibility race checks passed in 3.228s
(`/tmp/ensemble-compaction-stagec-worker-compatibility.txt`). Initial tests incorrectly
expected forbidden RegisterTool to succeed so later dispatch could reject it;
registration already enforces admission. That fixture error is retained in
`/tmp/ensemble-compaction-stagec-worker-initial.txt`; corrected and race outputs
use the same prefix with `fixture-corrected`, `race-final`, and `vet` suffixes.

Twelve focused scratch mutations failed their intended assertions: skipped startup
normalization, inherited main role, unavailable default factory, overwritten custom
factory, unchecked returned purpose, reopened tool registration, four-send budget,
second task admission, successful prose-only response, accepting the first submit
of a multi-call batch, inherited main forcing, and inherited worker purpose in an
ordinary descendant. The startup mutant's nonexistent-skill failure demonstrates
unexpected startup loading, not an actual MCP connection measurement. Logs are
`/tmp/ensemble-compaction-stagec-worker-mutant-<case>.txt`; scratch copies removed.

C2 does not yet launch these workers from normal memory policy: StartCompaction,
provider sizing/terminal reasons, staged replacement/recovery, owner-lifetime
scheduling, drained completion/accounting and registry/GUI dismissal remain required
Stage C work. The legacy recovery tests still fail as recorded above. Independent
C2 slice review found a missing purpose reset in the public Child helper and weak
schema substring assertions. The helper now resets purpose; a real derived child
proves ordinary construction. Both declarations and every actual input_schema are
compared structurally with the exact memory contract. Two additional focused
mutations remove that reset and allow extra submit properties; both fail the
intended assertions. Raw outputs use the worker-mutant prefix with public-child
and exact-schema suffixes. A C2-only scratch run restores HEAD Engine and omits
C3/recovery tests: final focused race passed in 1.141s. The independent equivalent
run including all Stage B tests passed in 1.526s. Scratch copies removed. These
results exclude in-progress C3 provider changes. No full Stage C acceptance is
claimed. The independent reviewer accepted C2 for scoped checkpoint with no
production blocker and no formal Stage C correction round; its audit records the
size comparison, exact source binding and remaining exclusions.

## C3 implementation evidence: provider allowance and terminal failure

After C2 checkpoint `ce52c166`, worker requests use the accepted Anthropic16,384
allowance and omit optional Chat/Gemini caps; ordinary rendering is unchanged.
All three SSE reconstruction paths preserve actual stop reasons. Semantic parsing
rejects worker length stops before submission dispatch even when arguments parse.
Responses already rejected incomplete output; it now retains complete measured
usage on terminal failure and on malformed arguments in a completed envelope.
Classic semantic failures likewise retain validated measurements. No successful
response fact is fabricated for that accounting, and absent usage is not zero cost.
This slice does not recover missing final measurements from an interrupted stream.

The real-runtime HTTP fixtures cover four routes × streamed/nonstreamed × success/
length-stop, plus Chat and Responses malformed-argument envelopes in both delivery
modes. Success preserves the exact roughly52KB memory string; synthetic9000 output
tokens are provider fixture measurements, not a tokenizer or live quality result.
Failure checks require no submission/tool dispatch and exact measured usage; the
main Agent's usage/occupancy remains unchanged before later managed import. Initial
compile setup used a nonexistent History.Events helper and was corrected to the
public EventsAfter API; the original setup failure remains in
`/tmp/ensemble-compaction-stagec-provider-initial.txt`. Corrected16-case tests
passed0.103s. The final20-case focused race passed1.919s. Inherited streaming,
visible-explanation and engine boundary race checks passed (engine1.019s,
ensemble1.366s); targeted vet passed. Raw files use the same provider prefix with
`corrected`, `race`, `compatibility`, and `vet` suffixes.

Thirteen focused mutations each failed their intended assertion: each of three
worker caps restored; each of three semantic length refusals removed; each of
three streamed terminal reasons discarded; classic failed usage dropped; Responses
terminal failed usage dropped; ordinary Anthropic requests given the worker cap;
and completed malformed Responses usage dropped. The last mutation covers both
delivery modes through their shared semantic path, not two independent proofs.
Logs: `/tmp/ensemble-compaction-stagec-provider-mutant-<case>.txt`; scratch removed.
Existing Responses incomplete-output rejection and unchanged schema validation
reuse their inherited tests; this slice does not claim a new independent mutation
for each inherited assertion. Final combined C2/C3 race passed2.163s. Independent C3-only source reconstruction
(accepted C2 archive plus exactly the C3 diff/test, excluding C4 and recovery)
passed C3+C2 race2.030s. The reviewer accepted C3 for scoped checkpoint with no
blocker and zero formal Stage C rounds; its audit retains the exact source binding.
No parent
lifetime, disk recovery, once-only auxiliary import or live capacity claim is made.

## C4 implementation evidence: selected-source acceptance

Memory selections now capture a digest of the exact selected entries and task
boundaries, effective source/destination policy, selected source-file bytes and
prior destination bytes/absence. ValidateSelection rechecks these values. It does
not hash the growing history: new dialogue, skills and unselected later files
remain independent. Actual source bytes must match the selected context already
at launch, including existing curated MEMORY.md. Nonterminal destination collisions
are refused. Policy generations across off/on transitions remain actor work;
equal current settings alone cannot prove generation continuity.

Real Memory/History/corpus tests cover unchanged selection, unrelated growth,
selected projected-content edits, same-name/same-length source-file edits,
unselected-file edits, changed policy, destination creation, existing destination
collision, and curated-file edits before/after selection. The C1 valid-prefix
fixture now writes actual source files, since selection fingerprints them; invalid
coverage still refuses before filesystem selection. Existing range and waterfall
controls pass. Initial setup misspelled common.SkillData as SkillLoadedData and did
not compile; corrected focused tests passed0.033s, final focused race1.167s and
memory/ensemble vet passed. Outputs:
`/tmp/ensemble-compaction-stagec-selection-{initial,corrected,race,final,vet}.txt`.

Five focused mutations fail intended assertions: ignore selected entry bytes;
ignore relevant policy; permit an existing destination; accept curated file/context
disagreement; and omit selected file-byte validation. The last mutation removes
both file/context comparison and captured-file comparison, explicitly a coupled
restoration of the old unchecked path, not two independent proofs. Raw outputs:
`/tmp/ensemble-compaction-stagec-selection-mutant-{entries,policy,collision,curated,source-bytes}.txt`.
Scratch copies removed. Independent reconstruction from accepted C3 plus the five
C4 paths passed selection/C1/waterfall race1.170s. The reviewer accepted C4 for
scoped checkpoint after a comment correction clarifying that Sources holds
fingerprints, not recoverable file bytes; no behavior correction or formal Stage C
round. No staged replacement, restart-window repair or actor generation/lifetime
claim is made. The runtime audit retains exact source binding and size comparison.

## C5 implementation evidence: correlated projection and journal receipt

After C4 checkpoint `ffa0bb68`, History accepts task-correlated launch, cut/retirement,
replacement and terminal outcome facts. A matching population applies the pair and
retains its committed receipt until a matching committed outcome; abandonment
cannot clear that committed state. An unrelated disk population cannot satisfy a
managed cut. New launches cannot replace an outstanding managed receipt, including
legacy launches. The latest terminal outcome survives snapshots for subsequent
staging-file cleanup; repeated identical-state completion is harmless and does not
change the main turn. Legacy uncorrelated event behavior remains compatible.
Successful Checkpoint now establishes its journal path under the same lock, so a
fresh embedding client calling Save before Load journals subsequent launch/commit
facts. This is a projection/journal slice, not completed disk recovery.

Real public History/Agent tests cover deferred source removal, exact task/file/range/
band correlation, before-cut population refusal, premature success, duplicate
launch/replacement, committed-versus-abandoned state, terminal cleanup identity,
unchanged main-turn state, and actual snapshot+journal recovery without model work.
Graduation retains the exact unselected third file. Initial new tests passed0.013s;
affected memory/context/save compatibility plus new receipt race passed3.062s;
history/ensemble vet passed. Outputs use
`/tmp/ensemble-compaction-stagec-receipt-{initial,race,vet}.txt`.

Eleven focused restored defects fail intended checks: ignored task/file/range/band
identity, wrong terminal commit state, discarded committed receipt, disk population
satisfying a managed pair, early cut, uncleared pending receipt after terminal
outcome, legacy launch overwriting a managed receipt, and missing journal activation
on fresh Save. Logs: `/tmp/ensemble-compaction-stagec-receipt-mutant-<case>.txt`.
Scratch copies removed. Final focused receipt race passed1.042s. Independent
reconstruction from accepted C4 plus the seven C5 paths passed receipt/C4/waterfall
race1.173s; the reviewer accepted C5 for scoped checkpoint with no findings and
zero formal Stage C rounds. Its audit retains exact source binding. The two original unsafe
Write/event/Archive restart probes remain failures until Memory staging/finalization
and runtime recovery gates are implemented; these receipt tests do not relabel them.

## C6 implementation evidence: staged disk recovery and corpus holds

After C5 checkpoint `c00af62b`, Memory owns the hidden versioned staging record,
exact launch/selection/candidate binding, idempotent destination promotion and
selected-source archiving. Finish retains the record; Recover returns a terminal
fact for the actor to journal before a second call removes staging. Committed
receipts refuse missing/corrupt/mismatched records and conflicting human edits.
Startup and memory-settings synchronization recover before scanning; a failed
committed receipt also blocks Memory selection/scan, recall search and rendering
of that Context. Recall construction now defers indexing until its first query,
so NewAgent cannot scan files before its embedding client loads History. This is
the recovery slice; normal actor compaction still uses its legacy path pending
the next scheduling/lifetime slice.

The original direct-Write probes remain recorded failures of the old API. Their
source is preserved at `/tmp/ensemble-compaction-legacy-recovery-before-staging.go`
and the original raw baseline described under C1. The new regressions explicitly
use Stage/Finish and real History snapshot/journal recovery; they do not claim
unchanged legacy probes now pass. Saved-boundary cases cover hidden candidate,
held cut, committed pair, promoted destination, one source archived, all sources
archived and terminal outcome before staging cleanup. These are constructed real
saved states followed by a fresh Agent Load/Ask, not forced process kills.

Eleven damaged-record/file cases prove recovery errors preserve the committed
receipt, prevent requests, block scan/selection/render/recall/re-enable, and permit
retry after an explicit repair. Empty/stale output and a real not-a-directory
staging failure preserve active source. A real unreadable directory proves recall
construction defers indexing while the subsequent query actually reaches that
corpus; this permission test passed, rather than skipping, on this host. Initial
new tests passed0.046s; final focused race passed1.243s. Affected memory/recall/
compaction compatibility race passed ensemble18.578s; internal memory/history had
no package tests and internal recall/engine selected no tests in that command.
Targeted vet passed. Raw outputs use
`/tmp/ensemble-compaction-stagec-staging-{compile,initial,race,final,compatibility,vet}.txt`.

Nineteen focused scratch mutations failed the intended checks: putting the hidden
record in active memory; skipping startup recovery; skipping archiving; ignoring
destination, source or archive conflict; rejecting already-retired matching source;
ignoring selection digest, candidate digest, committed population identity, record
version or unknown fields; bypassing scan, render or recall hold; eager constructor
recall; deleting staging before terminal outcome; accepting empty memory; and
accepting a stale Stage task ID. The active-candidate mutation exposes serialized
record bytes to Scan, sufficient to demonstrate visibility rather than semantic
memory quality. Skip-recovery fails on the required successful restart because the
remaining scan hold refuses access. Logs use
`/tmp/ensemble-compaction-stagec-staging-mutant-<case>.txt`; each scratch tree was
removed. The stale-ID mutation initially matched two source sites and was refused
by the scratch edit assertion before building; a scoped replacement then failed
the intended Stage refusal check. That setup assertion is not a mutation kill.
Malformed JSON uses the standard decoder; write refusal uses the actual filesystem;
inherited journal ordering/selection checks retain C4/C5 evidence rather than
claiming new independent mutations for unchanged behavior.

No model was called live. No power-loss guarantee, transaction against arbitrary
external editors, full worker lifetime/accounting, or complete Stage C acceptance
is claimed. Independent C6 review accepted this scoped slice with no production
correction and zero formal Stage C rounds. Its exact predecessor-plus-eight-path
reconstruction passed the six focused test functions under race in1.262s, with
no skipped startup permission test; the runtime audit retains the source binding.

Independent review requested the distinct existing-MEMORY.md promotion branch.
Added a real staged commit/restart with exact curated replacement and a human
edit refusal preserving bytes/receipt. The final focused race including both
passed1.258s (`staging-reviewed.txt`). Two additional mutations fail the intended
curated checks: skipping an existing destination update leaves old bytes; ignoring
its prior fingerprint overwrites the human edit (`staging-mutant-curated-{replacement,conflict}.txt`).
These bring the C6 total to21 intended mutant failures, with the setup assertion
above kept separate. No production correction was required for these added cases.

## C7 implementation evidence: auxiliary accounting seam

After accepted C6 `59c71d90`, Engine exposes value-copied UsageBuckets and
AddAuxiliaryUsage. The buckets retain vendor, configured model and the original
request's large-price tier; aggregation cannot reconstruct that tier afterward.
Both ordinary accounting and imports use the same locked totals helper. Imports
leave Engine's latest-main measurement and History unchanged, with no task-ID
ledger. The actor's forthcoming active-attempt bit remains responsible for calling
once after drain; this seam deliberately does not claim duplicate-delivery safety.

A real main actor and restricted child use an HTTP peer. The child corrects one
schema-invalid response, crossing the existing small/large price threshold; both
responses include cache-write/read counts. Cases cover successful completion,
measured length failure and an unpriced supported model. After actual child
Shutdown/Save, a manual import preserves exact original buckets, totals and cost,
leaves latest-main/History occupancy untouched, and does not reimport spend after
parent Save/Load. Receipt mutation cannot affect either Engine. The returned type
contains only value fields; that copy property also follows directly from its
representation, without adding artificial pointer state for a mutation test.

Initial test execution failed one overly strict whole-meter equality assertion:
map iteration changes floating-point addition order. The corrected assertion
compares exact counts/identity and uses a small tolerance only for derived money.
The original failure is retained in
`/tmp/ensemble-compaction-stagec-accounting-initial.txt`. Affected provider/cache/
status accounting race passed2.057s and targeted vet passed. After strengthening
exact imported-bucket assertions, final focused race passed1.086s. Seven targeted
mutations fail intended assertions: omit import, lose original tier, substitute
parent model or vendor, overwrite latest-main usage, omit copied input/cache counts,
and add each imported bucket twice. Logs use the accounting prefix with `race`,
`vet`, `final` and `mutant-<case>` suffixes; scratch copies removed. Duplicate
addition inside this single operation is distinct from future duplicate mailbox
handling. Existing reload/session and provider usage semantics retain inherited
coverage; no new live pricing or semantic memory claim is made.

Independent C7 review accepted the scoped slice with no production finding and
zero formal Stage C rounds. A slice-wide copy assertion replaces the earlier
map-order-dependent first-element check; its focused race passed1.084s. The
independent exact-predecessor-plus-five-path accounting/C3 race passed1.989s;
the runtime audit retains the source binding. Managed launch, drained completion/dismissal,
owner-lifetime scheduling, generation/headroom and once-only actor import remain
required before full Stage C acceptance.
