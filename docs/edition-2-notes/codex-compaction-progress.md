# Background compaction and structured submissions

Coordinator: Astra (Codex), October 14, 2026. Current status: Stage B accepted;
Stage C concrete design is next. Mapping and grader prerequisite are accepted.
The isolated accounting reconstruction is also independently accepted as a
maintenance stage, without a new export. No chapter changes, deployment or
publication. The Stage A observations below retain their original source/results.

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

## Next gates

With Bill's decision recorded, finish and review the narrow Stage B/C design
details above before implementation. Keep old chapter
numbers during implementation/testing. Commit the audit checkpoint separately
from future implementation acceptance. No old accepted snapshot can become a
new-order result by changing its directory name. Stage B now changes production
source in the isolated branch; the reviewed grader build helper is committed.
Chapter prose, frozen exports, live configuration and publishing artifacts remain
unchanged.
