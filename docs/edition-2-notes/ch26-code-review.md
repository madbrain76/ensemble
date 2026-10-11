# Original26 / Edition2ch23: independent reviewer preparation

Reviewer `reviewer_astra`, GPT-6 Astra/high, October 11, 2026.
**Preparation only; awaiting first-green implementation. Zero formal rounds.**
Mandatory `.agents/skills/ensemble-coding/SKILL.md` was read in full, with the
current AGENTS, governing updates and student-safe [ch26 notes](ch26.md).
No original26 implementation answer or student implementation was inspected;
no grader or provider call was run.

Read all466 lines of `book/chapter-26.md` at
`b451d2f14d8170b9075027a0a6606badfb442dad`; SHA-256 remains
`8db9779e9099d971e8fb8fd5c8f531e110d6a94c0bf234ed836f029aa159c56e`.
The coordinator's assignment-boundary completion check satisfies Bill's release
condition; historical source scores are not new student results. Accepted
predecessor is original25 source `e00d5bc`, main tree
`30fb05d154b9808eb88622aaff9c1424cd7e8af1`, direct combined Edition2ch22.
The existing narrow public contracts in ch26 notes remain sourced to grader
`9bab19dc914f04f16c5f7f5fdcc6c8345cf0b3c7`; verified file SHA-256
`ec6e7f98b4576cbb2e3bd834632cee8106ecc9152bf0edaa0d9115ad28dbc287`.
No further grader inspection is needed to restate those facts.

The independent review will preserve these manuscript obligations:

- Activity transitions through the existing Observer, with useful tool/detail
  information, and lifecycle distinct from thinking/tool/waiting/idle phases.
- Real named-peer create/restart/conversation preservation and deletion that
  refuses active children, with the registry owned by the existing Ensemble.
  Spawned children remain ephemeral; persistent spawning from the old brief is
  not part of the assignment.
- Each new management operation resolves only the caller's actual children;
  absent and foreign targets must be indistinguishable. Accepted duplicate names
  across ancestors and trees remain valid, so browser selection must resolve the
  actual node rather than assume a globally unique display ID.
- Change waits cover the specified state/exit/message/pattern/stuck conditions;
  cursor reads preserve actual new output; synchronous questions and asynchronous
  hints use the real actor; submissions have their required durable fallback.
- Section26.4 distinguishes orderly shutdown from immediate kill. Both interrupt
  modes end the current turn while preserving the Agent; Bill's prior ruling
  preserves other jobs and their real results. These lifetimes need actual checks,
  not labels on one equivalent shutdown path.
- The real CLI and optional GUI use the accepted public composition root. New
  fleet GUI clients enable original25's factory protocol for root and descendants.
  Initial tree plus ongoing changes, nested activity, conversation click selection,
  named create controls (name/model/skill), deletion and visible rejection remain
  required despite the older coder review's narrower historical deferrals.

Constructor names, React, and a particular incremental-versus-full-snapshot
implementation do not override Bill's existing architecture or KISS. Either
recorded tree-message spelling is accepted; visible behavior and actual owner
relationships remain the test. Do not recreate a flat root or introduce callback
bags to match the source narrative. Source confidence claims about compile-time
wiring, fleet-scale security or independent mutations are not assumed facts.

After first green: read the student review and new code, then matching original
source/tests; compare both editions against their own predecessor, including
active code, useful comments, exported contracts, tests and support files.
Final disposition waits for actual unchanged grading, focused mutation outcomes,
real browser/live evidence and the complete student review. No exception is
preaccepted. Preserve concurrent CodeRhapsody work and exact-owned-path commits.

## Stage1 preflight, October 11 (not formal review)

Implemented stage `55cbef230e675ed0d42547b7ad21e53c5b8a554f`, tree
`967ed4ff7d9379c37ce3eeadac2800803431e442`, covers activity and management.
Read the student stage notes and all current changed files; matching original
management code is now permitted comparison context. Registry and GUI stages
are explicitly unfinished and are not treated as final omissions. No old answer
code was relayed to the student.

A real defect is independently reproduced: `WaitChange([watched])` reports
`watched`/`message` when only the unselected child `other` asks a question.
Both child requests are parked, so no watched-child activity/state/output changed.
The implementation applies a manager-wide message sequence to every selected ID;
that is not the requested set-scoped monitoring (§26.4). The race-enabled probe
fails cleanly in0.075s. Controls prove no-progress/no-question waits time out and
a selected-child question wakes with its own identity (both pass,1.458s).
[Failure and controls](../../solutions/edition-2/evidence/ch26/reviewer-wait-scope.txt)
and adjacent reviewer probe sources preserve the setup and an initial HTTP-fake
cleanup timeout. This finding was sent through the coordinator as behavior and
rationale, without reference-answer code. The coder owns its correction.

Two small documentation inaccuracies were also relayed for final cleanup:
`FleetTools` still describes seven declarations after adding eight management
tools; `ChildStatus.State` omits the newly returned `idle` state. These are
accuracy corrections to exported contracts, not new capability requests.
Zero formal rounds used; final disposition remains pending full implementation.

## Stage2 preflight, October 11 (not formal review)

Reviewed exact `33c4e3315b2210bff946aa9d94bb27f7dda0c340`, main tree
`918fd9c466bdc55d7d4b245d78b0d846b043cc96`, through an isolated archive.
Moving GUI browser files were excluded. Read the student stage notes, complete
registry and tests, startup/fleet/test delta, and matching first-edition registry
and composition implementation. Mandatory `.agents/skills/ensemble-coding/SKILL.md`
and carryover were reread after context compaction.

The exact three preserved reviewer wait probes now pass under race (1.807s):
an unselected child's question no longer wakes the watched child, a selected
child's question does wake it with the correct identity, and the no-question
control times out. Per-sender revisions replace the manager-wide attribution.
The prior failure remains preserved; this is a verified correction, not a
reclassification of the old result. [Exact rerun](../../solutions/edition-2/evidence/ch26/reviewer-stage2-wait.txt).
The two exported-contract inaccuracies identified in stage1 are corrected too.

The shared `Start` path installs the recursive factory and human sink before
constructing/restoring root and named peers. Both actual executables call it;
offline render/dump retain the non-loading path. Named peers stay separate from
ephemeral children, while tree keys follow actual parent paths so repeated local
names remain selectable. The registry uses the required four identity fields,
restores conversations through normal Agent loading, and writes via temporary
file/rename. Named deletion checks children, shuts down the owner, then removes
registry identity and files. Its existing real-runtime tests independently pass
under race (1.169s), covering restart/provider-history preservation, repeated-name
path identity, delete/restart absence and an actually parked child's rejection.
[Focused rerun](../../solutions/edition-2/evidence/ch26/reviewer-stage2-registry.txt).
This does not claim every concurrency interleaving or final browser integration.

The modified GUI budget fake still exercises the built executable, persisted
settings, a two-batch turn bound, next-turn eighteen-batch behavior, restart and
zero/default behavior. Submitting as batch18 replaces the former eighteenth
shell call while retaining the actual dispatched-count/error assertions; it
adapts the external response to the enabled fleet completion protocol without
waiving the inherited budget contract.

Compared with the matching old registry/composition, this stage extends the
existing ownership root and Observer rather than adding callback wiring or
copying sibling services. Useful comments explain path identity, persistence
and lifecycle boundaries; the new public types/members have contract comments.
The change is proportional to the registry/selection requirements. Final whole-
chapter own-predecessor counts, management edge coverage, mutations, original
grading, browser/live evidence and final student review remain pending. No new
blocking finding in this stage, no acceptance or grader exception, zero formal
correction rounds.

## Full first-green preflight (not formal disposition)

Exact source `974f67be89272ea60fa17b460cefde0e478a137b`, main tree
`c34e436798cead5d2dc9d8899887da7b6eeda485`. Read the complete student stage3
notes, new GUI/view/browser code and tests, management/registry changes, public
contracts and README. Matching original activity, management, registry, GUI and
composition changes/tests were read after implementation, never supplied to the
student. Original GUI has no conversation-selection implementation and creates
by name/task rather than model/skill; those omissions do not waive this student's
manuscript obligations. The new implementation actually routes per-connection
selection and replay by owned Agent, reuses the Observer, and traverses trees in
a coalesced worker. Atomic child-admission closure repairs the earlier deletion
check boundary. No callback-bag or optional-module boundary regression was found.

**Blocking pre-submission finding:** repeated `select_agent` can deadlock the
GUI after its slow browser disconnects. The new selection branch sends a replay
batch to a bounded queue while holding the shared observer mutex. Once full,
a failed/exited writer leaves the socket receiver blocked on that send, retaining
the mutex needed by observations and Close. A real WebSocket peer, one64KiB
history item and320 selections reproduce it; no provider or internal mock is
involved. A one-selection positive control passes. The final race run fails
the intended bounded shutdown assertion in4.141s, with actual goroutine stacks
at `gui.go:332` (channel send) and `gui.go:465` (Close mutex). The first600ms
admission attempt passed; allowing3s exposed the full queue. Both are preserved,
not relabeled. [Probe, control and stacks](../../solutions/edition-2/evidence/ch26/reviewer-selection.txt)
and adjacent exact sources record the finding sent through the coordinator.
No production code was changed; correction and its mutation belong to the coder.

Original grading actually reports0/100 FAIL: six checks cannot start inference
without a model; the GUI check supplies `--port` to the headless executable and
gets usage/refused connection. Those are the observed fixture boundaries, not
fleet behavioral failures, but they also provide no proof of the masked valid
obligations. Full exception assessment awaits the final student explanation and
supplemental evidence; no additional unobserved envelope failure is claimed.

[Provisional own-predecessor counts](../../solutions/edition-2/evidence/ch26/reviewer-first-green-counts.json)
separate languages, tests and scripts. Student production Go adds1007 physical,
860 active,99 comment-only and48 blank lines; original adds1067/838/107/122.
Student Go tests add412/380/18/14 and browser tests68/66/1/1; original Go tests
add864/669/93/102. Student production JS adds34/32/2/0, HTML7/7/0/0, CSS9/7/1/1.
Original JS adds98/79/5/14, HTML changes-3/-4/+1/0, CSS adds59/53/0/6; its
removed418-line browser-support script is separate, not a production saving.
The original snapshot also removes `mcpws/mcpws.go` (109/87/13/9); its net Go
total includes that deletion, which is not an equivalent new management feature.
These are net tracked-source totals, excluding evidence/generated artifacts;
comment classification is lexical and inline code stays active. Refresh at final
freeze. Lower test totals alone are not an improvement, nor are the original
source's confirmation-only interrupt checks equivalent to observing a retained
real tool result and a later successful task. The current explanatory comments
cover ownership and critical lifetime ordering, but the new production comment
ratio is below the teaching target; assess clarity rather than pad to a quota.

Final checks still depend on correction of the demonstrated deadlock, targeted
mutation outcomes, actual live/browser provenance and the full student review.
No formal correction round consumed; no chapter acceptance or export.
