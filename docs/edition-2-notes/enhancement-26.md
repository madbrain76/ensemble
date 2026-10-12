# Edition 2 Chapter 26: Gateway author record

Author: `edition2_author`, October 12, 2026. Coherent manuscript submitted for
independent prose and revised-grader review; **not yet accepted**. Own paths are
`book/edition-2/chapter-26.md` and this record. Mandatory technical handoff:
`.agents/skills/ensemble-coding/SKILL.md`.

## Reading and source identity

A compaction occurred during initial preparation. Before authorship, the author
reloaded all four COMPLETE documents together in active context, with these
untruncated reads:

- Workflow1–350,351–750,751–1150,1151–1550,1551–1950,1951–2325.
- Carryover1–350,351–700,701–1000.
- Voice1–320,321–632.
- Chapter procedure1–310,311–599.

The earlier failed whole-voice read was not counted. No further compaction
preceded this draft. The full373-line mandatory coding skill, current AGENTS
instructions and updated enhancement/map headers were read. Every later
compaction still requires another full four-document reload before writing.

Read all469 lines of original `book/chapter-27.md`, source revision
`6452456e5549be5566bfc920138b82f53a715756`, SHA-256
`5463f060485fc838858ebe44115c70ac062582ba9f3179281df95ddcce9d974d`;
all363 lines of `ch27.md` and282 of `ch27-code-review.md`; full original coder
review and working notes. Accepted source is the exact197-file direct `ch26`
export, tree `6d3e3d2e54a4def135f9d5c23ce3242d7d4c0313`, integrated at
`36dcaaa0` after independent acceptance `0afc45e8`. Read its full public
Gateway README, protocol/client/Noise implementations, server/admin and all
three command implementations, plus actual live report/result and independent
live-artifact check. Chapter23's accepted opening/TL;DR supplied continuity.

## Source-relative plan and preservation

Stake: the operator needs agents to receive authenticated work across process
boundaries without parking their tool handlers or confusing network receipts
with completed work. It resolves in the actual three-process Ensemble interview.

Preserve the Deephold game and Puffin care motivation, with the reported uptime
bounded to the original account's date. Preserve the original seven-section
progression, router-versus-Agent distinction, key-label explanation, readable
provenance example, immutable-config-group argument, small Handler interface,
offline delivery reasoning and fixed-question interviewer. The original's
line counts and sweeping mutation/security conclusions do not teach the current
implementation and are retained in their historical sources instead.

Necessary additions are the missing public contract, actual actor integration,
and complete reproducible launch. They serve the original exercise rather than
adding a transport, application framework, scheduler, dashboard or coding tool.
No word-count target applies; hard lint and the source-relative reading matter
separately. Review bookkeeping stays here rather than becoming chapter narrative.

## Material feedback dispositions

- **Welcome ambiguity:** Bill explicitly asks whether the revised chapter still
  leaves the interface ambiguous. TL;DR and §26.4 now publish `QueueDepth int`,
  JSON `queue_depth`, and its outstanding-after-resume meaning. The public
  ClientConfig, Handler, message/child fields, key/Noise helper signatures and
  callable operations are explicit. MCP names/argument keys, hello/welcome/ack
  wire shapes and complete server configuration are also published.
- **Original0/100:** retain the actual unchanged original grade on `15e0386b`
  as a probe compile failure at undocumented `QueuedCount`. The module alias
  did compile. No runtime failure or passing original row was established.
  The new chapter teaches the corrected public contract without burdening its
  reader with the old private spelling. Revised legacy comparison may explicitly
  select that historical spelling; frozen runtime remains unchanged.
- **Noise:** correct the first ephemeral public key to cleartext and explain
  own-private/peer-public DH, directional session keys and identity-bound
  prologue. Keep the public library rather than importing the old private helper.
  Correct reverse-proxy TLS legs: origin TLS is possible, but the proxy still
  terminates the client connection. Avoid universal endpoint secrecy claims.
- **Provenance:** source fields come from authenticated agent state and the
  strongest accepted group, not arbitrary asserted platform identity or an
  unrelated admin group. A readable header supplies evidence to the model,
  not mechanical instruction compliance. Puffin's platform example motivates
  the distinction without claiming a shipped Discord adapter.
- **Public callbacks:** sequential Handler callbacks cannot perform blocking
  RPCs on their own receive goroutine. Actual adapter uses existing Agent.Submit,
  a readiness publication for reply workers and explicit shutdown joins. Owner
  diagnostics remain reachable; no callback service bag or core change is taught.
- **Delivery:** Send admits to the queue; handler return triggers ack; an actor
  admission is not durable completion. Explain explicit reconnect, recipient
  sequence identity, plaintext queue files,30-second/final individual snapshots
  and their nontransactional/abrupt-stop limits. The actual ack-test correction
  teaches why an earlier channel notification did not establish receipt order.
- **Historical interview:** final manuscript obligation governs the older
  stand-in/actor-deferral notes. Actual second-edition three-process run is
  retained, with model opinions and coding capability kept distinct.

## Primary-source checks

Read the official [Noise specification](https://noiseprotocol.org/noise.html),
revision34 dated2018-07-11, on October12. Its handshake state machine explicitly
writes ephemeral public keys in cleartext and computes DH from local private
and remote public values. The KK pattern/security sections support the narrow
explanation; no novel cryptography or complete-security claim is introduced.
Read [Cloudflare SSL/TLS concepts](https://developers.cloudflare.com/ssl/concepts/):
edge and origin certificates/connection legs are separate, and encryption modes
control the origin leg. These primary links appear beside the teaching.

Read the OpenAI Docs skill for the named model guidance and searched then opened
its official [GPT-6 Astra model page](https://developers.openai.com/api/docs/models/gpt-6-astra).
The identifier and `high` effort are supported; this does not establish account
eligibility or convert the optional Chapter24 credential route into an official
public endpoint promise. Reader examples use the already accepted launch, not a
new provider feature or unsupported Anthropic model.

## Actual evidence preserved

Read `evidence/ch27/live-result.json`, `live-interview.md` and
`reviewer-live-check.json` directly. On `ab6d0516`, the actual Gateway, Ensemble
adapter and public interviewer completed four questions in11.739 seconds, with
three distinct process identities and all exits0. Saved incoming messages,
request/completion facts and report text agree. Actual model is GPT-6 Astra;
high effort and Codex route retain launch/source provenance, not an independent
full provider-wire capture. Inputs90/161/232/305, outputs42/42/45/41, all cache
buckets0. The chapter quotes the actual token-use answer and paraphrases the
proposed context-finding tool as opinions from this no-tools interview.

Initial event counting looked for nonexistent `response_received`; four actual
`response_ended` records corrected the capture without another paid run. Final
schema/null-state/diagnostic and ack-test ordering corrections retain their
later source and affected deterministic/race evidence. The original25 ultimately
detected defects include invalid/coupled/mislabeled initial attempts; they are
not25 independent security proofs. No new paid demonstration is authorized or
needed for this enhancement.

## Coherent-stage checks and pending gate

Direct `go run ./cmd/lintprose book/edition-2/chapter-26.md` passes every hard
check at4384 prose words. Soft negation/person-gap warnings were read; technical
API/lifetime detail is retained without invented anecdotes or padding. Exact
examples and both numeric rubrics are with the coder; full independent prose
review is with `reviewer_astra`. Proposed revised command uses `-ch26` and seven
original weights15/10/15/15/15/10/20, total100 with no parity. Numeric public
client/socket claims remain distinct from actual actor/live/owner/lifetime proof.

Initial coherent manuscript SHA-256:
`55f7c95a2391c85e45efd6e8e2651d270e96fe1220c3ccd3c24b4731c09e7ab7`.
Final grader/source/sample evidence and independent disposition remain pending.
No original/frozen/runtime/crossover files were edited. Publication migration is
explicitly deferred until this gate; current examples retain current real paths.
