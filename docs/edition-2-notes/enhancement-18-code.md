# Edition 2 Chapter 18 grader: coherent stage, review pending

Coder: GPT-6 Astra/high. Original Chapter 19, accepted direct `ch18`, legacy
`solutions/ch19`. Read the complete original, student/comparative records,
mandatory `.agents/skills/ensemble-coding/SKILL.md`, carryover and workflow.
Original graders, runtime exports and historical results remain unchanged.

The separate command is `go run ./cmd/grade-edition2 -ch 18 -json`; the legacy
comparison adds `-gui ./solutions/ch19` and positional `./solutions/ch19`.
It builds the actual optional GUI command, uses a typed primary and isolated
profile/workspace, and preserves all seven original weights: 15/15/15/15/10/15/15.
The small explicitly Edition 2 adapter reuses the original signed OAuth and
Responses HTTP peers. New behavioral predicates and launch code live separately
in `internal/gradeedition2/ch18*.go`.

## Current results and fixture corrections

Evidence is under `solutions/edition-2/evidence/ch19/enhancement18/`.
Actual full results `accepted-green.json` and `legacy-green.json` are both
100/100, exit zero, in 2.676 s and 1.794 s respectively. `stage-source.json`
binds those bytes. `green-source.json` binds the subsequent adapter publication
synchronization; focused real-command race testing passes in 2.557 s. Vet passes.
Exact examples, targeted sensitivity and final independent review remain pending.

Retained initial results are accepted 100/100 and legacy 55/100. The legacy
correctly refused missing-grant and quota scenarios but deliberately exited 1;
that completed refusal is now distinguished from successful sessions, which must
exit zero. Actual token issuance/request admission, observed turn error and no
metered credential switch remain required. A signal, early exit or arbitrary
nonzero status alone cannot establish refusal.

The original fake's terminal envelope omitted `model`. Legacy provenance-aware
replay correctly excluded opaque material with unknown model provenance. Actual
retained terminal responses identify their model: `model-fixture-provenance.json`
records sanitized fields from original live log lines 14 and 191. The forwarding
peer fills only an absent terminal model from the actual request; exact encrypted
replay remains required. The original fake is preserved. A failed initial
extraction used SSE framing against per-event JSON logging and found no events;
the corrected extraction reads the actual event/data representation.

The first corrected accepted attempt scored 30/100 because its callback server
closed after receiving the code but before finishing the browser response. The
source confirms this ordering. Login now waits for the bounded actual command
and checks real token exchanges instead of treating browser EOF as failed login.
That attempt remains `accepted-stage.json`; it is not a transport/parser defect.

The first race check exposed the original OAuth peer assigning its server pointer
after starting its handler goroutine. The adapter publishes that completed
assignment through the existing mutex that handlers acquire before reading it,
before releasing the URL to the child. Both original warning output and corrected
race result remain. This is fixture synchronization, not a runtime change.

## Observed claims and boundaries

Positive API-key and subscription requests use their actual configured/minted
bearers. Secret canaries are absent from observed output, request bodies and
non-profile files. Actual PKCE/issued-client exchange, forged-signature and nonce
rejection, inference-grant refusal, and two proactive rotating renewals are
observed. Each renewed request is checked against its corresponding minted token,
not the final token retrospectively. Profile security, concurrent/cross-process
single flight, atomic replacement, logout revocation and transient recovery retain
accepted source tests and the earlier 39-mutation evidence; this numeric fixture
is not a general OAuth/security proof.

Responses checks include actual successful `/etc/hosts` bytes in the correlated
function result and an encrypted reasoning item in the next request. The plan
peer retains absent Content-Type and empty terminal output. The summary fixture
withholds the remaining response until the GUI receives the first summary text;
valid chunk splitting is allowed. Complete separate summary/text deltas include
tool-call-adjacent content, not a summary conditioned on completed tool output.
Native speech and listening retain their original separate evidence.

API-key caching requires the actual explicit marker on developer constitution
content plus explicit mode and the four-marker budget. Plan requests omit the
paired mechanism under the accepted exception; historical provider rejection
established only the options field, not an independent marker rejection. Both
admission and midstream quota scenarios require real authenticated requests and
reported refusal without switching to the available metered key.

Original grading remains 15/100 on `06cbf20`, including its vacuous billing award.
The real plan task completed in 14.887 s, but its summary-dependent launcher
killed the idle CLI at 150.03 s, without graceful shutdown/save. No public summary
or cache grant was observed. Those historical results, source versions and
limitations are not replaced by these local protocol-peer runs. No paid request
or production sign-in occurred during enhancement.
