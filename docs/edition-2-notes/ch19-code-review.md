# Original Chapter19 / Edition2 Chapter18 — independent review

## October11 public-contract preparation only

Reviewer `enhancement_reviewer_resumed` read all666 lines of original
`book/chapter-19.md`, revision `e1605916114230570b348ab5961c193fa1cebf80`,
SHA-256 `41fe6b803bd837c0ec3f4ed818cf02187c588c63b1ebee4713484d9fa404161c`,
and the complete [coordinator preparation and official contracts](ch19.md).
Mandatory `.agents/skills/ensemble-coding/SKILL.md` was reread in full; current
carryover/workflow apply. No implementation, old answer, grader, private auth
profile or current main source was inspected. No provider request was made.
Original18 remains active and its formal review takes priority.

The chapter's scope is credential separation, actual PKCE/OIDC identity and scope
validation, serialized rotating-token refresh and persistence, model/route-aware
Responses transport, separate incremental reasoning summaries, and visible billing
policy. Existing owner interfaces and the optional GUI boundary still apply.
The paragraph saying a renderer calls GetBearerToken conflicts with the surrounding
no-renderer-auth seam; the governing ownership and credential-isolation rules settle
this without a new architecture question. Credentials must stay out of requests'
model-visible content, snapshots and diagnostics. Passing login is not inference
permission, and HTTP200 is not a completed successful Responses turn.

The coordinator's sourced current contract resolves the short versus full quota
error-code discrepancy. HTTP status alone cannot authorize credential destruction
or metered fallback. The explicit “Yours” permission to defer chatgpt_then_api
remains meaningful despite the general three-mode description and billing rubric.
If the student implements fallback, explicit enablement and announcement before
the first paid request are binding; deferral is not permission for silent fallback.
Multiple-account schema support does not require a multi-account UI. Provider-managed
compaction mentioned in the narrative is not part of this exercise.

The remaining consequential route question is already documented: Chapter19's
unconditional explicit caching and later reported plan restrictions cannot both
be treated as universal facts. Use actual supported-route evidence when reached;
do not silently discard required behavior or call fake cache reads a live grant.
The same applies to model eligibility and streamed summaries: actual nontrivial
reasoning/tool use supplies evidence; one trivial response without summaries does
not settle capability. Current official notes are reused, not independently
re-researched or claimed as an authenticated account result here.

The public source shows login/status/ask examples but lacks the complete grader
invocation, fake authorization/discovery/JWKS/Responses endpoint overrides, isolated
profile selection and observation envelope for summary/billing events. Coordinator
may recover these missing interface facts when required, with exact source and
supplied contract recorded. Do not disclose old implementation choices or infer
new functionality from a private harness. Suggested helper/type names remain
adaptable to the existing codebase.

Live sign-in may require Bill's direct participation when19 arrives. Existing
permission to use selected API keys does not authorize copying or rotating another
running application's OAuth profile. A faithful local OAuth/JWKS fake is valid
protocol evidence but cannot be labeled live authentication. This preparation
does not initiate sign-in, waive live use, or stop unrelated Chapter18 work.

No formal review round or acceptance begins here. Wait for the authorized chapter,
frozen submission and complete student review before comparative assessment.
Original22/24 stops, stopped authorship and eventual direct Edition2 ch18 export
remain unchanged. Required handoff: `.agents/skills/ensemble-coding/SKILL.md`.

## October11 narrow grader-contract preflight (implementation still underway)

Independently read the coordinator's newly supplied external contract, unchanged
checks/runner and relevant launch/fake-provider sections. No old answer or
unfinished student source was inspected; no grader, provider or auth flow ran.
Everything below is a **latent source observation**, not a student failure,
score, acceptance or formal correction round.

The external invocation described in ch19.md matches the harness: login uses
`auth login --print-url`, isolated HOME/EN_HOME and issuer override; the same GUI
binary later receives a private absolute `--save`, provider configuration, JSONL
stdin and WebSocket summary/text/turn observations. Both scenarios have an API
key but only logged-in ones require the minted OAuth bearer. This tests default
credential selection, not permission to override an explicitly selected billing
mode or silently fall back. A course model is an input fixture, not a live model.
The skill fixture indeed omits `type`; established original10 validation can mask
inference startup. Login might execute earlier, so do not predict every scenario
or score from this omission before the actual run.

A consequential historical conflict is now explicit in the unchanged runner:
`ch19CheckCacheBreakpoints` rejects **any** plan-route cache options/markers;
`ch19CheckCacheBreakpointsMetered` requires explicit mode and an actual developer
content breakpoint on the API-key route. That disagrees with original19's
unconditional explicit-cache wording and matches the later historical route
restriction already noted in preparation. The harness comment attributes it to
its own live observation; this is not current account-specific evidence. No
implementation policy or exception follows merely from reading the fake.
Retain the actual supported-route evidence requirement; do not remove metered
caching globally to pass a plan restriction.

Scope precision correction: fake OAuth token issuance at lines953–972 returns
openid/profile/email/offline_access plus chatgpt.tokens.use.direct, without
resource.invoke. I initially described this too strongly as a possible required
login rejection. Coordinator re-opened the current official sign-in page: its
request table includes resource.invoke, but the explicit granted-scope validation
checks chatgpt.tokens.use.direct. Therefore the fake's omission is only a scope-
shape observation, **not a confirmed defect or waiver**. Request and granted-scope
requirements must not be conflated. The coordinator records the independently
verified official wording in student-safe ch19.md; this review did not repeat
provider research or access credentials.

For later review, the seven labels overstate several predicates:

- Credential15 checks minted bearer versus API key and some secret exposure;
  it cannot prove common-interface ownership, all durable/log surfaces or that
  rendering remains independent of credential acquisition.
- OAuth15 checks the positive PKCE/issued-client exchange and a bad-signing-key
  scenario. The latter passes whenever no inference occurs, including unrelated
  login/startup failure. It needs a successful positive control; state/nonce/
  audience/expiry/granted-scope rejection and sign-out need their own valid
  current-scope evidence, not inference from that score.
- Refresh15 uses a90-second token (inside the source's five-minute window) and
  checks refresh occurrence, no observed token reuse and two turn completions.
  It does not prove concurrent serialization, atomic persisted replacement,
  terminal-error behavior or preservation on transient failure.
- Responses15 checks actual route/request fields, requested encrypted reasoning,
  successful turn flags and a correlated call/output pair. Merely requesting
  encrypted output does not prove replaying its received bytes; an error tool
  output can satisfy the pair predicate. Preserve actual parser/tool/replay proof.
- Summary15 counts separate kinds and ordering before the first text delta.
  A buffered batch emitted at completion can preserve that ordering; this alone
  does not prove incremental delivery while inference is still active.
- Billing15's two rejection scenarios use `usage_limit_reached`, not the prepared
  current full quota code. More seriously, a launch fatal is accepted as a stop,
  so zero actual inference can pass. It does not exercise permitted announced
  fallback, persisted billing policy or credential retention. Deferring optional
  fallback remains allowed by the public exercise; if implemented, the existing
  consent/announcement requirements still bind.

These are masked-obligation risks for the later source/test/live review, not
instructions to expand scope or rewrite the grader. Current public security,
transport and billing behaviors remain required even if a fixture prevents them
from running. No extra retry policy, profile feature or test framework is implied.

Exact inspected file provenance (file revision; SHA-256):

- ch19_checks.go: `3bfb49407e060eea45e8665fa1c969818195eb31`;
  `5be31eb750e63ccf2aeab13e8948099b8fb60a10d37a9b2684494801982b80b3`.
- ch19_run.go: `d68fef808f0d6d50daf682afec976cb401e1b586`;
  `f0b667d1669a9002e606caf0204d8d68b1f11d7e0701ba80139c1e376466fdbd`.
- ch19_harness.go: `70e9e3f4d72848c21143c1a3874a165debff32f7`;
  `d5acd0fdb3db968a7242fb96d0a45afbbb96c8756dde8c2a1c3a181674cff37e`.
- ch19_fakeoauth.go: `6012fc717f0cf33785ef5093b1373f17ba0a1e3e`;
  `26291e3641daae5fcf196e981c2738b7143bb5b9879af913f5683b24c183d218`.
- ch19_fakeresponses.go: `d68fef808f0d6d50daf682afec976cb401e1b586`;
  `0de734804694fc31c336dfea663f0c9afe323109476ded6784fc5a2cb130c5e2`.

Keep grader-internal coverage analysis in this reviewer-only record. Only narrow
public-contract corrections go through the coordinator to the student. Mandatory
handoff `.agents/skills/ensemble-coding/SKILL.md`; preserve Edition1 work and the
original22/24 stops. Await frozen submission for formal comparative review.

## Original grader observation — before final submission

October11 independent evidence check. The unchanged original command completed
on `06cbf20000cf3592ac2966a1041ece12b33dbc04`, main tree
`c550fa3c13faec59268cc8d7486a6ed8f590d68f`: **15/100 FAIL**, grader exit1,
make exit2. Read [raw output](../../solutions/edition-2/evidence/ch19/original-grade.txt)
and [provenance](../../solutions/edition-2/evidence/ch19/original-grade-provenance.json),
independently matched the recorded source tree and raw SHA-256
`bd940a590b17495f5789b0a6e39fa9c825e7e7997d5151ca83a3c914f92e1a6a`.
Coordinator retained these at d44532b. No repeat was run.

Six failed checks (85points) report the same actual GUI startup rejection:
`skill "base" has invalid type ""`. Five first failures name scenario A and
refresh names B. This confirms the preflight's fixture mismatch, not a discovered
Responses, caching, refresh or reasoning implementation failure. From the
unchanged launch order, reaching those GUI errors means their fake login command
returned successfully; it does **not** establish all OAuth validation claims,
because the later positive predicates were skipped. It is not real plan login.
The first-failure aggregation does not separately expose every later scenario's
observations, so do not infer them from repeated labels.

The awarded billing15/15 supplies no positive billing proof. The unchanged C
predicate explicitly accepts a launch fatal and requires neither a successful
inference request nor a quota refusal. The common untyped-skill launch prevents
inference; the aggregate PASS is therefore consistent with—and explained by—the
preflight's vacuous path. No per-scenario positive C trace appears in this output.
Likewise the bad-key scenario's absence of inference cannot alone separate
correct signature rejection from unrelated startup failure. Preserve the actual
15/100 rather than calling this an85-point implementation failure or a demonstrated
billing success.

Cache-route restrictions, default credential selection, refresh error classes,
summary incrementality, replay and real tool outcomes remain **masked obligations**
for final independent review. No cache rejection or private-save-path problem was
observed by this run. The resource.invoke request/grant distinction remains the
corrected observation above, not an accepted fixture defect. The narrow missing-
type explanation is independently substantiated, but chapter acceptance and any
final exception decision await frozen implementation, full student response and
real supplemental evidence. No old answer or changing product was inspected in
this check; zero formal correction rounds. Required handoff:
`.agents/skills/ensemble-coding/SKILL.md`.
