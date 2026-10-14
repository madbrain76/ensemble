# Background compaction and structured submissions

Coordinator: Astra (Codex), October 14, 2026. Current status: Stage B accepted;
Stage C supporting slices C1–C9 are checkpointed; final actor integration is
under independent and coordinator review. The coordinator's factory-refusal
reproduction passes after correction; a related GUI removal edge is being fixed.
Accounting, authentication and historical crossover reconstruction are accepted.
The reviewed combined web/status reconstruction plan is released for implementation.
No chapter changes, exports, deployment or publication. Earlier observations
retain their actual source/results.

Bill authorized proceeding after reading Cadence's
[handoff](codex-compaction-subagents-plan.md). He separately required full
Chapters 14 and 15 and both governing documents before surgery. Those full
reads are complete, including recovery of truncated output. The coordinator
also read AGENTS.md, the entire `.agents/skills/ensemble-coding/SKILL.md`,
deployment HANDOFF.md, voice and chapter procedure, chapter map, enhancement
record, publication record and crossover notes. Reload full governing/craft
documents after compaction before formal authorship/review.
After coordinator compaction, all four governing/craft documents and the entire
coding skill were reloaded in full before assessing the completed audit notes.
The coordinator repeated that full reload after the next compaction before
accepting the accounting reconstruction plan; the truncated workflow range was
read separately to recover the missing text.

## Isolation and ownership

- Baseline: `56ec71cbde8477b6959ad2fbc65f9237a0fa8346` on main.
- Worktree: `/Users/bill/projects/ensemble-compaction`, branch
  `codex/background-compaction`. Main checkout remains untouched.
- At capture there were no tracked modifications. Existing untracked publication,
  historical, runtime and evidence paths remain untouched. Full status is in
  `/tmp/ensemble-compaction-baseline/initial-status.txt`.
- Cadence owns its data-directory migration and memory repository. No edits to
  deployment, credentials, launcher, runtime settings, binary or shared crossover
  notes. Reconcile any overlapping application changes before integration.
- Codex owns new progress/audit notes and the eventual submission/compaction
  changes in this separate worktree. Do not stage unrelated files or push.

Later main reached `f910e8f5`; its two new commits change only Cadence's migration
note, with no application-code overlap. The coordinator read the full updated
note: Cadence records Bill applied the move and verified runtime output under
`/Users/bill/projects/cadence-memories`, preserving the explicit save path and
all six disabled bands. These are Cadence's observations, not a Codex live test.
Do not rerun migration or assume the original deployment data path remains active.
Reconcile those note-only commits at a clean checkpoint; preserve both workstreams.

## Review assignments

- `compaction_dependency_audit`: independent chapter/fixture/parity/snapshot
  dependency audit, read-only except its own audit note.
- `compaction_runtime_audit`: independent runtime ownership, schema and lifecycle
  audit, read-only except its own audit note.
- Both are GPT-6 Astra/high, fresh contexts, with explicit full-reading and
  `.agents/skills/ensemble-coding/SKILL.md` handoffs. Their initial Stage A was
  read-only; the separately assigned implementation roles are listed below.
  This is maintenance of the existing runtime, not a new blind student evaluation.

## Stage A observations

Chapter 15 explicitly requires a finite synchronous direct compressor and says
it is not a sub-agent. Background children deliberately replace that scheduling
contract; they must preserve its source-retention and memory semantics. Current
full governing records contain historical versions/pauses; their current header
and Bill's new direction govern this maintenance.

Sequential baseline commands use an environment stripped of inherited `LLM_`,
`OPENAI_`, `ANTHROPIC_`, `GEMINI_`, `EN_` and `CODEX_` variables, only in the test
process. Raw output and command/timing records are in
`/tmp/ensemble-compaction-baseline/`. No live provider requests are authorized
by these fake-provider checks.

- `go run ./cmd/grade-edition2 -ch 14 -json ./agent`: 100/100, exited 0.
- Same for 15 and 16: each 100/100, exited 0.
- Same for 22 and 23: build failed before grading. `grade.Build` selects
  `agent/cmd/`, which now has subcommands but no Go source of its own. No runtime
  behavior was assessed by those invocations. Preserve this result; investigate
  a narrow explicit-entry-point fixture repair separately.
- Supplying a CLI built with `go build -o /tmp/ensemble-compaction-baseline/ensemble-cli .`
  from `agent/` to the unchanged `gradeedition2.Chapter22` and `Chapter23`
  functions gives 200/200 and 100/100 respectively. A temporary root-module
  driver invoked these functions with explicit core/GUI paths, then was removed.
  These diagnostic passes bypass only build selection; the advertised CLI
  invocations above still fail. No grader source was changed.
- Core `go test -short -race ./...` exited 1; root package reproduced
  `TestFleetCLIHumanReplyAndEscalation/true: human sink was unreachable` before
  any edits. Other reported packages passed. No all-green suite claim.
- Core and GUI `go vet ./...` each passed.
- GUI `go test -race ./...` initially failed because this new worktree lacked
  Playwright. After `npm ci --ignore-scripts --no-audit --no-fund` in its GUI
  directory, the same suite passed (11.916 seconds). No tracked dependency
  files changed.
- Root `go test ./internal/gradeedition2 ./internal/grade` exited 1 (21.281
  seconds): `TestChapter13Speech` cannot import Playwright from the frozen
  `solutions/ch13/gui` fixture. Installing current GUI dependencies does not
  install that separate snapshot's dependencies. This setup failure is retained;
  no speech-runtime regression is established by it.
- The separately run response-lifetime/cancellation baseline now passes:
  `go test -race ./ensemble -run 'TestModelResponseBeyondOneMinute|TestBlockedModelResponseCancellation' -count=1`
  from `agent/`, exit0 in63.110 seconds. Both real61-second cases were exercised.
  New compaction/schema tests and live work have not run. The baseline is not
  acceptance of the proposed implementation.

## Accepted prerequisite: grader entry-point selection

Coordinator changed `internal/grade/harness.go` to prefer `cmd/` only when Go's
package metadata identifies `main`. A directory containing subcommands or a
library no longer hides the root executable. A selected command's compilation
error still fails; it does not fall back to a healthy root. No grading predicate
or application behavior changed.

Five real temporary-module tests cover root, nested command container, command
preference, library fallback and broken preferred command. Before the fix, the
container and library scenarios failed; all pass afterward. Independent review
accepted with zero correction rounds. Its scratch mutations (always-root and
restored directory-only selection) fail the intended checks, with the correct
control passing; exact hashes and comparison are in the dependency audit.

Actual advertised explicit-source commands now pass: Chapter22 200/200 and
Chapter23 100/100, exits0 in2.489 and1.125 seconds. Raw reports use
`grader-{22,23}-buildfix` under the existing temporary evidence directory.
After installing locked dependencies with `npm ci --ignore-scripts --no-audit
--no-fund` in the isolated frozen ch13 GUI fixture (only ignored node_modules),
`go test ./internal/gradeedition2 ./internal/grade` passes in30.429 seconds.
No frozen tracked source changed. Earlier failures remain recorded above.

## Reviewed dependency result and approved mapping

Both independent audit notes are complete:
[chapter dependencies](codex-compaction-dependency-audit.md) and
[runtime boundaries](codex-compaction-runtime-audit.md). Their findings are
source review, separate from the coordinator's executed baseline above.

Recommended mapping: current17–23 become15–21; current15/16 become22/23;
0–14 and24–26 keep their numbers. The accounting chapter currently requires
compressor/judge usage and memory cache boundaries before the proposed new
memory chapters. Stage those obligations at new22/23 while keeping ordinary
accounting and handoff boundaries in new15/18. Nothing is dropped. Memory's
existing explicit Chapter14 parity stays Chapter14 parity.

The cumulative references are the material consequence: old17 already includes
memory/recall. Renaming it new15 cannot demonstrate the new sequence. Recommend
reviewed staged reconstruction from accepted14 through new23, then affected
carry-forward/validation for24 and26. This is maintenance using reviewed source,
not a fresh blind student run. Replace book-facing exports only after validation;
preserve existing tags, source/manifests/results and Edition1. Historical grading
can use the immutable old checkout and its own grader without a new dual-mode
framework. Chapter25 remains a placeholder.

Bill approved this mapping, staged accounting and reviewed replacement of current
exports on October14: "Yes, approved." This resolves the handoff's mapping and
historical-export decision. Apply the move only after the behavior/review gates;
the approval is not acceptance of unbuilt intermediate solutions.

Stage B is assigned to `structured_submission_coder` (fresh GPT-6 Astra/high),
with `compaction_runtime_audit` as independent design/code reviewer. The coder
must present concrete schema/API/persistence design before implementation.
Coordinator owns the narrow grader build-selection repair, independently reviewed
by `compaction_dependency_audit`. Work remains in the isolated worktree.
The Stage B design is now accepted by the coordinator and independent runtime
reviewer: see [the coder's contract](codex-structured-submissions.md). Implementation
is released. The design uses a pinned genuine validator with an explicitly
admitted Draft2020-12 subset and no references; admission bounds reject rather
than truncate input. This is design acceptance, not implementation acceptance.

Separate preparation branch `codex/reordered-foundations` at
`/Users/bill/projects/ensemble-reordered` starts from `a8eec84a`.
`reordered_foundations_coder` (fresh GPT-6 Astra/high) owns its staged `agent/`
and note, with the dependency reviewer checking it independently. Its first scope
is accepted14 plus accounting only, with future memory/recall obligations deferred.
This isolated preparation does not replace runtime main or publish new snapshots.
The coordinator and dependency reviewer accepted its concrete source plan before
implementation: accepted14 plus the selected16-to17 accounting delta, with
memory/recall integrations deferred to their new introduction points. Existing
supported-model and response-lifetime corrections carry forward explicitly;
affected separate example modules join the core/GUI validation. Implementation
was released in that branch. The resulting `288fc6e8` (agent tree
`2b671b45dc5f7b760f2faaff2a77c83e6c7c09c8`) is now independently accepted at
initial review, zero correction rounds. Its old14/17 graders each score100/100;
affected checks and five targeted defects pass their intended observations.
The original concurrent EOF timing failure remains unexplained, alongside passing
exact14 and isolated staged controls. No export, renumber or auth stage is accepted.

Validation encountered host disk exhaustion at142 MiB free. The coordinator
paused new builds across both coders and the reviewer, verified no Go build/test
process remained, and ran `go clean -cache` against the11 GiB regenerable build
cache. The command completed successfully and restored11 GiB free. Source,
module downloads and evidence were preserved. Link/build failures during ENOSPC
are retained as setup failures; affected checks require clean reruns.

## Stage B acceptance

Independent runtime review accepts the frozen structured-submission implementation
at initial formal review, zero correction rounds. The coordinator read the full
disposition and accepts the stage. Common holds vocabulary, Fleet owns the compiled
contract and durable results, and the composition root refuses invalid construction
before factory effects. Legacy APIs/rows remain compatible; custom results retain
JSON data through public retrieval and GUI live/replay paths. The accepted bounded
dialect and admission ceilings are in the coder's contract.

Reviewer feedback corrected canonical-size/reload inconsistency, multiplicative
validation-work exposure, unsafe shutdown on a wrong-owner factory result, missing
required-parent admission and an incomplete GUI assertion. These were genuine
pre-submission assistance, not independent student discoveries. No Stage C cleanup
or memory behavior is claimed by this acceptance.

Final focused races, full GUI race and both vets pass. Full ensemble race passes
in91.778s, with the later parent guard separately covered by focused races. The
full core suite retains its reproduced CLI human-sink failure. An unnecessary
browser glob invocation lacked required harness environment; its separate speech
assertion also fails on exact `a8eec84a`. Failed fixture/setup attempts and the
mutation evidence's coupled/repeated claims remain in the coder/reviewer notes.

Coordinator ran advertised explicit-source old22/23 commands sequentially on the
frozen stage:200/200 in2.356s and100/100 in1.072s, exits0. All214 source-manifest
files remained unchanged during the runs and matched at independent review.
The reviewer verifies21 distinct targeted defects across24 executions with stated
limits. Changed-file active production grows306 lines; the reviewed explanation
and runtime coverage justify that scope. No new result browser or scheduler.

## Runtime decisions to carry into the implementation plan

Coordinator preparation reproduced both source-identified disk/event windows on
an exact `a8eec84a:agent/` archive. The temporary external module uses real
Ensemble, Memory and History with an HTTP provider fake. It saves at the selected
boundary, creates a fresh Agent, restores that saved history and starts an ordinary
turn. This is a controlled restart-boundary probe, not a forced process kill.

- After replacement `Memory.Write`, before cut/population, restart injects the
  uncommitted candidate alongside the retained conversation. The source-retention
  control passes, then the candidate-absence assertion fails.
- After paired retirement/population, before `Memory.Archive`, the pre-restart
  context has retired the source, but restart reinjects those old source files.
- Controls omit the first replacement write and complete the second archive;
  both pass on the same unchanged application source.

`go test -mod=mod -v -count=1 ./...` in the temporary probe module: baseline
exit1, two intended failures (0.010s); controls exit0, two passes (0.011s).
The original probe was restored after controls. Raw files are under
`/tmp/ensemble-compaction-window-probe/`; probe SHA-256
`ef869c4e0945e84c763c837dfd81f7a9f1ca85941b4aa4139d1418b5c5fae2af`.
An initial run against the in-progress Stage B tree reproduced the same failures;
the immutable-archive repeat establishes that they precede Stage B. These results
support the narrow Stage C recovery work, not a claim that recovery is fixed.

A separate exact-baseline range probe creates two source files whose parsed Thru
precedes File. Scan accepts them; a completed write/paired retirement/population/
archive retires no context entries, and Graduation immediately selects the same
source again. `TestReversedRangesCannotRepeatFold` fails for that intended reason
(0.008s, exit1). Replacing only the filenames with valid single-file identities
passes (0.008s, exit0). The original probe was restored; SHA-256
`a513c5c1fde6fa7bbd26ddc7ea0924394117f0d9264b41a421784be4056b27dd`, file
`ranges_test.go` and logs in the same temporary directory. This establishes the
separate no-progress defect before repair; it does not explain Cadence's observed
Conversation compressor failure.

- Preserve the current accounting promise: compaction cost belongs in owner
  session totals exactly once, attributed to the actual model/tier, while main
  last-request counts remain separate. A separately retained child receipt is
  necessary but does not replace that existing promise. Review the concrete
  auxiliary import boundary before implementation; never count it again through
  a fleet aggregate or replay it as new spending.
- Ordinary parent-turn completion and interruption leave useful compaction alive,
  following the handoff's ordinary-child preference. Explicit worker kill,
  owner shutdown and relevant memory disable cancel/invalidate it. This changes
  the old interrupted-direct-compressor expectation deliberately.
- A submission is a candidate. Acceptance waits for reliable child completion,
  including save, and a safe parent boundary with selected-source/settings
  revalidation. Automatic cleanup removes Fleet, root and GUI references and
  preserves durable outcomes. No inference on replay.
- Retain default string submission compatibility and add an explicit raw JSON
  custom path; a string convenience call must obey the same custom validator.
  The accepted Stage B contract specifies the schema subset, pinned validator and
  durable envelope. Bounded worker request policy and Memory-owned crash recovery
  remain Stage C design work before that implementation.
- The audit initially identified disk/event crash windows from source; the
  controlled restart probes above now reproduce them before implementation.

## October 14: next-stage preparation

The coordinator reloaded all four governing/craft documents and the complete
coding skill after compaction. Primary maintenance HEAD is `0a732ed4`, which
merges Cadence's documentation-only migration commits through `f910e8f5`.
Cadence reports the applied data directory as
`/Users/bill/projects/cadence-memories`; its explicit conversation save remains
at the deployment's `session/save.json`. This is Cadence's recorded live
verification, not a new coordinator deployment test. No application overlap or
live changes were introduced by the merge.

The independent dependency reviewer accepted authentication source plan
`185a1e31969eb94c2cbec54a9437455988ce1fba`. The coordinator read the full plan
and disposition and released that bounded reconstruction from the accepted
accounting tree, with implementation review still required. It retains the
single Send path, resolved-credential sanitization, once-only successful renewal
accounting and the no-total-response-deadline correction. Recall/judge integration
remains deferred. No export, renumbering or new live use is implied.

Coordinator source inspection found two Stage D fixture issues at this HEAD:

- `internal/gradeedition2/ch15_driver.go:memoryCompressor` checks only that the
  sole tool is named `submit`. The handoff's claim that it also inspects the
  `memory` property is inaccurate. Revised classification must distinguish an
  ordinary submit-only child from the compaction worker's actual contract.
- `internal/fakevendor/fake.go:handle` holds its server-wide mutex while calling
  `Options.Route`. A held memory route therefore blocks main requests and
  request inspection too. Releasing only the memory router's own mutex cannot
  establish background progress. Audit other route consumers before a narrow
  shared-fake correction, or use an appropriate external fake; preserve archived
  Edition 1 graders.

The coordinator then reproduced that fake-server serialization through real HTTP
in `/tmp/ensemble-fakevendor-concurrency-probe`. On the unchanged helper,
`go test -mod=mod -race -v -count=1 ./...` fails both intended assertions:
the unrelated main request and `Requests()` remain blocked while the routed
auxiliary response is held (exit1,0.612s). Releasing only that gate before the
main request gives the passing control (exit0,1.017s). Original test restored,
SHA-256 `13511e65951acc094fc04f8e72faa6a4e18e060dcf469b3ff96bd4ba592ffaaa`;
raw `baseline.log` and `control.log` remain beside it. This proves a fixture
limitation, not a new runtime compaction result. No fake implementation changed.

Further source inspection found an existing narrow option:
`ch23_driver.go:managementGUI` holds selected HTTP requests in a reverse-proxy
boundary outside the same capture mutex. Stage D should consider that established
pattern before changing shared fake concurrency semantics for every route caller.
This is a design option, not an implemented correction or extra requirement.

The current main-turn completion barrier also does not establish background
worker completion. Stage D needs observable separate outcomes and causal
barriers while retaining all eleven existing weighted obligations. No grader,
scoring or chapter change has been made or approved by this preparation.

## Stage C design accepted; implementation released

Coordinator and independent runtime reviewer accept the revised complete
[`codex-background-compaction.md`](codex-background-compaction.md). The design
uses one owner-lifetime attempt, restricted real children, reliable completion
and save, selected-source validation, Memory-owned hidden staging and journal
receipts, and bounded child/reference cleanup. Three Engine.Send invocations
are distinct from existing bounded authentication retries. Ordinary parent
completion/interruption preserves useful work; explicit kill/shutdown and
relevant policy changes cancel or invalidate it. The stated near-window wait
is limited to an outstanding attempt at the existing 95% main-input threshold.

Pre-implementation review clarified bounded attempt-owned accounting instead
of an accumulating deduplication set, default tool registration versus actual
safety ceilings, and the trusted custom-factory callback boundary. A valid JSON
submit accompanied by a provider length/incomplete terminal reason still fails
before dispatch. Committed-but-unfinalized disk work blocks recall, corpus scans
and main rendering until recovery; it cannot be mislabeled abandonment. The
runtime audit records the independent disposition and required evidence.

The coder is released to implement coherent tested slices and checkpoint them,
with independent implementation review afterward. No runtime correctness or
semantic memory quality is established by design acceptance. Grader/chapter
changes remain Stage D, and renumbering/export/live work retain their later gates.

Parallel authentication reconstruction is accepted at `970676abeebce8f156d5cdc7ee0dd47a40708136`,
Agent tree `21b15bc6cd1c24a2ee3e03873a87a9c080744402`, initial independent review
with zero correction rounds. The reviewer verified all reported checks and source
counts, old18/17/14 each100/100, both61-second regressions and four intended
mutation failures; its independent frozen-source focused race run passes1.236s.
The coordinator read and accepted the full disposition, then released a source
plan only for current19/new17, the still-ungraded historical crossover. Historical
live evidence keeps its source identity; no paid repeat, export or renumbering.

Main advanced to documentation-only `67060b6e` while this work continued. Its new
`cadence-session-directory-migration.md`, read in full by the coordinator, describes
a prepared and fixture-tested move of application session files into
`cadence-memories/session/`; it explicitly says not applied. It adds no application
code conflict. Merge that note at a clean checkpoint, and re-establish actual
deployment paths before any later authorized live operation. Do not treat the
proposed session location as active or run Cadence's migration.

## October 14: coordinator review and reconstruction resumed

After answering status questions, the coordinator resumed active work at Bill's
"Please proceed." Earlier replies had ended coordination while runtime delegates
continued; that delayed the crossover source-plan gate. This repeats the workflow's
handoff lesson and is coordinator delay, not a technical blocker.

Before formal review, reloaded all four complete governing/craft documents and
the entire coding skill, then reread all current Chapters14/15 and the complete
handoff. Root inspected the historical crossover's complete three-path delta,
plan484c7234 and independent dispositionfa29b7ee. The bounded current19/new17
implementation is released with its stated two mutations and affected checks.
It remains ungraded; no paid task, new export or whole-tree merge is authorized.

Runtime slices C1–C9 are checkpointed throughf111f1b3: range admission,
restricted finite workers, provider sizing/terminal handling, selected-source
validation, correlated journal receipts, staged filesystem recovery, auxiliary
usage buckets, child completion/removal, and GUI dismissal. Each has scoped
independent review; none alone establishes the complete background actor.
The coder's C10 working tree is separate from this exact review baseline.
Root's overall review is underway, not accepted, and uses the immutable source.

### Reproduced reverse-purpose factory refusal defect

Root's real-runtime external-boundary test on exact f111f1b3 calls
StartCompaction with a trusted constructor-only custom factory. The factory
returns a real owned child with Compaction=false and LogPath pointing to a
preexisting sentinel file. CreateChild correctly refuses the purpose mismatch,
but its Shutdown writes the empty log header over that file; ForgetCompaction
then retains the rejected ordinary child because its purpose is false.
The same redirected-path factory retaining Compaction=true is a passing control:
refusal preserves the log and removes that managed child. No inference or mock
was used, and arbitrary callback side effects are outside this finding.

Actual command: `go test -race -run '^TestRootFactoryPurposeRefusal$' -count=1 ./ensemble`
in an isolated f111f1b3 Agent archive, with provider/runtime variables stripped
and a fresh EN_HOME. Exit1,0.022s package time; the ordinary-return case fails
both exact-byte preservation and root-membership checks. Raw output:
`/tmp/ensemble-root-review-factory-refusal.txt`; scratch location is recorded in
`/tmp/ensemble-root-review-path.txt`, with test `agent/ensemble/root_factory_refusal_test.go`.
The archive extraction emitted a Python tarfile deprecation warning; compilation
and test execution succeeded before the intended behavioral failures.

Coder and runtime reviewer received the unchanged probe. The existing C8 purpose
regression covers the opposite ordinary-request/managed-return case and missed
this direction. Prior scoped acceptance remains history. Correction must preserve
ordinary valid lifecycle while rejecting the returned child without writing into
its rejected paths or retaining its new registry entry. Recheck the unchanged
probe and appropriate ordinary/managed controls before accepting the correction.
This is overall review of accepted support, not a claim C10 has been submitted.

## October 14: integration review and next reconstruction handoff

After the next compaction, the coordinator again reloaded the complete skill,
workflow, carryover, voice and chapter procedure, then all Chapters 14/15 and the
handoff before resuming formal review. The initial combined tool output was
truncated; separate complete reads recovered it. Summaries were not substituted.

The unchanged root factory-refusal test body now passes on the C10 correction:
`go test -race -run '^TestCoordinatorFactoryPurposeRefusal$' -count=1 ./ensemble`,
exit0,1.030s. Root overlaid current changed/new Agent files onto its exact f111f1b3
scratch, binding file hashes in `/tmp/ensemble-root-review-fixed-source.json` and
checking those source bytes did not change during capture. The first attempt
failed to compile because the coder had retained the same test name; root renamed
only its independent function/comment, preserving the body and assertions.
That setup output remains `/tmp/ensemble-root-review-factory-refusal-fixed-setup.txt`;
the actual passing result is `/tmp/ensemble-root-review-factory-refusal-fixed.txt`.

Root inspected the explicit construction-discard path and exact-identity removal,
then the complete C10 production diff. The earlier source-only completion-count
concern is corrected: attempt identity and first-completion admission precede
worker decrement and auxiliary import. No duplicate-delivery runtime reproduction
is claimed. Independent runtime review separately identified the GUI consequence
of removing a rejected ordinary child; the small membership correction and actual
refusal test remain pending. Full Stage C acceptance is withheld until that final
source and the independent review are complete. The coder's reported full core
run retains its reproduced CLI human-sink failure; the Ensemble package and full
GUI pass. Exact commands, timings and mutation evidence live in the coder/reviewer
notes rather than being relabeled as coordinator executions.

Historical crossover reconstruction `078cccbd` is accepted after independent `003ae75d`,
Agent tree `f94f414deb977d9f106b8d749a57faa934987049`. The next plan `68793cad` is accepted
by coordinator and independent reviewer `076048cd`: current20 → new18 combined
web/status, including only the necessary earlier GUI settings-lock correction.
Root read the complete plan, exact `05b43a0` lock diff, 23-path historical delta
inventory and predecessor locking sites. The new synchronous model acknowledgment
requires that lock ordering now; memory itself remains deferred. Supported fixture
and unknown-price limits stay explicit. Implementation is released in the isolated
reconstruction branch, with focused causal regression, grader serialization and
independent implementation review before advancement. No whole-tree merge or
new accepted export is authorized at this boundary.

## Next gates

Stage B is accepted and Stage C's reviewed design is checkpointed at `01dfa4ce`.
Complete staged implementation and independent review before releasing the
chapter/grader enhancement. Keep old chapter numbers during implementation and
testing. Continue each reconstructed cumulative stage only after its independent
acceptance. No old accepted snapshot can become a new-order result by changing
its directory name. Chapter prose, frozen exports, live configuration and
publishing artifacts remain unchanged.
