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
