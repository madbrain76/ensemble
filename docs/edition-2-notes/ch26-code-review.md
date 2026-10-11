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
