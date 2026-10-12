# Original Chapter 27 / Edition 2 Chapter 26: independent code review

Reviewer: `reviewer_astra`, GPT-6 Astra/high. October 12, 2026.
**Preparation only. No student implementation reviewed or accepted; zero formal rounds.**

## Reading and source boundary

Reread all469 lines of original [Chapter27](../../book/chapter-27.md), SHA-256
`5463f060485fc838858ebe44115c70ac062582ba9f3179281df95ddcce9d974d`, and all
current [student-safe preparation/contract notes](ch27.md). Earlier complete
reading of the historical brief, coder review and working notes remains recorded
there. The final chapter requires actual Ensemble actor integration and a
three-process model interview; the historical external subject is not that proof.
No original answer code was read before the student's implementation.

Read all original grader files at revision
`558b29fb098a1f79a60794f467b75762177fbeee`:

| File | Lines | SHA-256 |
| --- | ---: | --- |
| `internal/grade/ch27_harness.go` |315|`cab0928523a1ee489980d65ac571f125d10b9718987996036ae99f2ae0834689`|
| `internal/grade/ch27_probe.go` |553|`1d1d1bf9cbad8c389f1aa553a1c3bf48a0028ddc4f637e39cd95b68b73d8e4ec`|
| `internal/grade/ch27_checks.go` |90|`eeec8a960e14677be504d2399c791c4505b3927f9763ddf11be5125f549aee07`|

The complete workflow/carryover/voice/procedure remain active together from the
full reload recorded at the Chapter23 gate; no intervening compaction. Entire
`.agents/skills/ensemble-coding/SKILL.md` and AGENTS instructions remain active.
The student extends committed main tree `ada66c1572f1e5d22efb6ac59f23ab2b61dd03ac`
in the coordinator's isolated worktree. Shared user main/GUI/fleet/skills changes
and frozen exports are outside this reviewer's mutation scope.

## Grader preflight

The seven original rows total100. The coordinator's public launch/config/API
erratum matches the actual probe calls. No additional input contract or old
implementation recipe is needed at this stage. Run the original grader unchanged
on the frozen student submission; no score or exception is preaccepted.

The external module aliases the submitted source to the first-edition module
name (harness73; probe import). This is a **potential** fixture mismatch, not an
inevitable compile failure: a self-contained public package may still compile
through Go's replacement. Diagnose the actual result and imports rather than
renaming the student's module or installing a compatibility facade in advance.
The dependency detector at harness51–54 only searches `/agent/internal`; it
cannot detect the student's `ensemble/internal/...` paths. Any awarded C6 credit
therefore still needs the actual resolved graph and real outside consumer.

Several labels exceed what the original probe observes. Its provenance check
requires only a header prefix and contained body, not the exact authenticated
identity/header fields. Its delivery handler returns immediately; later queue
inspection after sleeps does not establish that acknowledgment waits for handler
return. Reconnect exercises the running server, not restoration of dynamic state
after a server restart. The probe does not run the interviewer against Ensemble
or inspect actual model input/output. These are coverage limits, not reasons to
invent more grader points, edit the original harness or waive chapter obligations.
The eventual source/tests/live record must establish those existing promises.

## Independent review focus after a coherent freeze

Follow actual owners and construction in gateway server, public client and each
binary. The gateway owns routing/authorization/delivery, never an Agent or model
request. The chapter's public Handler callbacks are the required external API;
they do not authorize internal sibling-service callback bags. Inspect necessary
transport/library ownership under the skill without imposing a preferred package
design. Public client reuse must avoid private runtime imports in the real graph.

Verify authenticated sender provenance independently of group authorization, a
successful authorized route beside refusal, subset registration beside over-grant
refusal, and immutable configured versus persisted runtime policy. Inspect real
handler-return acknowledgment, offline replay/resume and dynamic-state restart,
with observations tied to actual effects and source versions. At-least-once is
not exactly-once; periodic saves do not promise zero loss after abrupt termination.

The actual gateway, Ensemble and interviewer must run as separate processes.
Follow received envelopes into the existing actor and responses back through the
public client, with real model interview evidence and normal process completion.
An external stand-in or a plausible report file cannot establish that integration.
Keep security claims scoped to tested boundaries; primary Noise explanation
corrections already in ch27 notes do not prescribe new cryptography.

After implementation, compare the student's complete review, production/comments/
tests and own-predecessor delta with the matching original answer. Reuse valid
unchanged evidence and require meaningful targeted defects for new claims; inspect
whole failure outputs and cleanup. No broad or paid freshness run is requested by
this preparation. Gateway authorship remains gated on implementation acceptance;
Chapter25 remains reserved.


Student design question, before implementation: whether the server may consume
public `gateway` wire/Noise types, or must split a second public protocol package.
Reviewer response: no split is required solely to change the package diagram.
The chapter explicitly requires the public boundary to avoid private runtime
imports. Server use of its protocol/crypto API is not by itself sibling-service
wiring; actual behavior placement and ownership still require source review.
Keep server policy/queues with their owner rather than putting them in Client.
This is a bounded architecture interpretation, not acceptance of unwritten code
or a supplied implementation recipe. Coordinator informed. Mandatory
`.agents/skills/ensemble-coding/SKILL.md` included in the student handoff.
