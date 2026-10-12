# New Edition 2 chapter: the temporary Codex endpoint

Coordinator brief, October 11, 2026. This records Bill's assignment and evidence;
it is not a finished chapter, reviewed design or implemented feature.

## Bill's direction

Edition 1 ends at Chapter 27. New chapters belong only to Edition 2. Apply this
change only to `solutions/edition-2/main/`; do not edit `agent/`, first-edition
solutions or frozen Edition 2 exports. Under the existing three merges, Gateway
is planned as Edition 2 Chapter 24; this first wholly new chapter follows as 25.
Gateway's Edition 2 implementation and review remain outstanding.

The cold open should explain why the transition to OpenAI ran into an economic
blocker: our public OAuth route answered successfully but reported no reuse of
the context we kept resending. Bill reports filing a bug with OpenAI. Using
Codex's endpoint is intended as a temporary way to continue while caching on the
official OAuth route is broken in this workflow. The intended destination stays
the official route once a matched probe establishes that it works.

## Evidence for the opening

Use the [probe record](../codex-endpoint-probe.md) and its raw sanitized receipts.
The original three-call Codex test reached 4,864 cached tokens out of 4,998.
The fresh comparison used identical request bodies across the two credential/
endpoint pairs. Both completed public OAuth calls reported zero cached tokens;
both completed Codex controls reported 4,864. The third request on each route
hit overload, and a bounded public retry also hit overload. Keep those failures
separate from completed cache measurements. No claim of a provider-wide diagnosis
or a third successful public request is supported.

For GPT-5.6 Sol, the checked ordinary/cache-read rates are $4/$0.40 per million
input tokens. That tenfold rate difference motivates the chapter. It is not a
measured tenfold total invoice: output and cold input still cost separately,
and the experiment did not observe subscription billing. Bill's report of the
financial impact should remain attributed to him.

The story's concrete sequence is available without inventing an incident:
successful answers concealed zero reported caching; an existing Codex sign-in
reached a different Responses route; the matched request reused most input
there; the public retest still failed to show reuse. Explain why a correct
answer was insufficient and why the workaround must remain replaceable.

## Implementation and authorship boundaries

No implementation or new chapter manuscript is complete. The existing book
team is paused after the Chapter 16 draft. Before later authorship, the author
and reviewer must hold all four governing documents in active context and reload
all four after every compaction. The new chapter must follow the implemented,
independently reviewed change, not invent a replacement authentication design.

Retain separate public OAuth and Codex credential/route identities, current
explicit billing/fallback rules and existing ownership. The successful probe
establishes inference and cache reuse only; production token renewal, tool and
reasoning replay, usage reporting and return to the public route remain work
to implement or verify. Do not turn these items into a speculative framework.
All coding work uses the selected GPT-6 Astra/high; the matched diagnostic model
remains GPT-5.6 Sol because it is the model in Bill's proposed test.
