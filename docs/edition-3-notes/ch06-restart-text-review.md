# Chapter 6 restart — single formal text review

**Verdict: REQUIRED_REVISION — one bounded attribution condition.**

Read the current prerequisites, full Chapter 6, actual 25-addition/3-deletion diff, author note, current student and consolidated independent-review records, and relevant accepted `solutions/edition-3/ch06` source. This is the single final text review, not a reopened implementation or operations review.

## Findings

The four additions are useful and accurate: native output-limit fields precede stream parsing; the protocol CLI must forward the existing `in_flight` boundary before response deltas; tool display chunks include each arriving name fragment without relying on supplemental metadata or repeating cumulative names; and OpenAI result adjacency belongs to the disposable wire projection, not rewritten chronological History. The accepted renderer, parser/plain-delta path and CLI support these explanations. The opening, terminal scene and teaching sequence remain intact. No expansion into process history or additional machinery is warranted.

**Required correction: distinguish retained reference API illustrations from the accepted restart's literal API.** The TL;DR says the interfaces follow the “accepted Ensemble implementation” and calls `Render(*Context)` the accepted signature. Its `StreamCallbacks`/`Parser` and `Stream` bitmask examples, carried into §§6.3–6.6, describe the older reference rather than the restart snapshot. Current `internal/common/common.go:101–141` instead exposes `Render(Context)`, `Parse(*http.Response, ParseCallbacks) ([]Event, error)` and string-valued `StreamCapability`; `config.go` uses `ModelFeatures.Streaming`. Existing historical/example qualifications are useful, but do not resolve those literal accepted-implementation assertions. In §6.6, saying a single boolean must lose text streaming or invent tool-argument chunks, and that a bitmask is the minimum required structure, also overstates the representation requirement: the accepted SSE/plain capability plus native parser preserves complete Gemini tool objects without manufacturing fragments.

In one localized author response, make the attribution consistent: identify the retained API/capability examples as the reference design rather than mandatory literal restart signatures or representation. Do not replace the examples wholesale or require code changes. Preserve the actual behavioral contract: parent-owned configuration/services, one parsing entry point across delivery modes, response-local part identity and authoritative finals, supported versus plain fallback behavior, and the CLI boundary/tool-chunk requirements. This is a prose-accuracy correction, not permission to weaken those requirements.

## Evidence limits and closure

The existing baseline/final lint outputs both report **0 errors**, with the same two warnings: 19 negation forms and a person-gap increase from 5,109 to 5,319 words. Neither justifies voice surgery; lint was not rerun.

The capped Gemini result was a fully delivered **137-character partial answer**, not complete prose. Qualified operational acceptance rests on the recorded stream, hint and job evidence, not typing cadence or guaranteed answer completeness. Claude's unchanged path remains carried evidence. Shape comparisons concern each response's own aggregation, not byte equality between independent generations. Historical examples remain historical; correction rounds and approved exceptions need not become chapter features.

No manuscript, source, grader, earlier-chapter note, tests, builds or live providers were changed or run. After at most one author response, root should directly verify the attribution condition and close it; no second formal review or new campaign is requested.

## Root condition verification

Verified the author's sole response against the reviewed draft and accepted source: API and capability illustrations are now attributed to the retained reference; actual restart signatures are stated; the bitmask is no longer claimed to be universally necessary. The ownership, single parsing path, plain/stream behavior, response identity, authoritative finals, CLI boundary/name payload and wire-only adjacency requirements remain intact. Revised lint exits 0 with the disclosed soft warnings. The required condition is closed; this is direct root verification, not a second formal text review or a new PASS verdict.
