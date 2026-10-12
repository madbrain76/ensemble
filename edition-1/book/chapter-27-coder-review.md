# Chapter 27 — Gateway: Coder Review

## Must-Fix

None. The chapter text matches the verified code.

## Enrichment (optional, author's call)

1. The `NopHandler` embedding pattern is a standard Go idiom worth a sentence
   calling out that it exists for forward compatibility, which the chapter
   already does. No change needed.

2. The MCP tool surface (8 tools: `gateway_declare_contactable`,
   `gateway_grant_membership`, `gateway_revoke_membership`,
   `gateway_send_message`, `gateway_list_agents`, `gateway_create_group`,
   `gateway_edit_group`, `gateway_delete_group`) could be enumerated if
   the author wants to show the full admin surface. Currently covered by
   "The administration surface is exposed as MCP tools."

## Do NOT Add

- Do not name the Noise-KK implementation source as "OpenADP" — Bill's
  ruling supersedes the design doc. The code was copied from
  `~/projects/gateway/pkg/noise/noise_kk.go`, byte-identical.

- Do not describe the keys as "ed25519." They are X25519/Curve25519
  Diffie-Hellman keys. `KeyPrefix = "x25519:"` in `gateway/keys.go:22`.
  The `ed25519:` label is explicitly called a misnomer in the code comment
  at `keys.go:17`.

- Do not mention bill's employer, its frameworks, or any coworker by name.

- Do not discuss Ensemble actor-loop integration — deliberately deferred
  per the TL;DR's "What NOT to build" section.

## Facts: Verified vs Unverified

### Verified (from code, this session)

| Claim | Source |
|-------|--------|
| Keys are X25519, labeled `x25519:` | `gateway/keys.go:22` (`KeyPrefix = "x25519:"`) |
| Provenance header format `[GATEWAY source=<channel> user=<user> trust=<trust>]` | `gateway/envelope.go:103` (`FormatProvenance`) |
| Three trust levels: owner, trusted, external | `gateway/envelope.go:27-34` |
| Handler interface: 4 methods (OnDeliver, OnConnected, OnDisconnected, OnError) | `gateway/handler.go:18-33` |
| NopHandler exists for embedding | `gateway/handler.go:37-44` |
| Gateway binary: 116 lines | `cmd/gateway/main.go` (verified by `wc -l`) |
| Total non-test lines: 3,058 | `find ... ! -name '*_test.go' \| xargs wc -l` |
| Total test lines: 1,326 | `find ... -name '*_test.go' \| xargs wc -l` |
| Total: 4,384 lines across 20 Go files | sum of above |
| Persist every 30 seconds | `cmd/gateway/main.go:110` (`time.NewTicker(30 * time.Second)`) |
| 10 frame types, 5 each direction | `gateway/protocol.go` (hello/ack/response/ping/mcp_request + welcome/deliver/pong/error/mcp_response) |
| Interviewer: 4 fixed questions via `defaultQuestions()` function | `cmd/interviewer/main.go:38-44` |
| `go list -deps ./gateway` shows zero `agent/internal` imports | verified via command |
| Client API: Connect, Send, Respond, RegisterChild, DeclareContactable, Close, Done | `gateway/client.go` public methods |
| ParsePublicKey accepts any prefix | `gateway/keys.go:43-47` |

### Unverified (from handoff, not re-measured this session)

| Claim | Source |
|-------|--------|
| `make grade27` = 100/100 | handoff memory, previous session |
| Mutation audit 7/7 | handoff memory, previous session |
| Full cross-chapter sweep clean | handoff memory, previous session |

The previous coder instance (Opus 5) was interrupted during `make grade27`
by the model switch. The score is reported from the handoff, not from a
measurement in this session. The author should verify before citing it.

## Decisions Made During Coding

1. **Noise-KK copied, not adapted.** Bill's ruling: copy from
   `~/projects/gateway/pkg/noise/noise_kk.go` (166 lines), byte-identical
   to the copy under `coderhapsody/gateway/`. The design doc's "adapt from
   OpenADP" was overruled. Bill: "It's boilerplate code, not at all novel,
   but easy to get wrong."

2. **X25519 labeling.** The code explicitly labels keys `x25519:` and
   documents that `ed25519:` is a misnomer from upstream habit. `ParsePublicKey`
   accepts any label prefix but enforces 32-byte key length.

3. **Trust stamped at delivery, not acceptance.** Provenance is added when
   the gateway delivers, so a redelivered message cannot accumulate two
   headers.

4. **Queue before deliver.** The server queues messages to disk before
   attempting delivery, so a message written to a socket the recipient
   never read is not lost.

5. **Reply authorized by origin, not contactability.** The sender started
   the exchange, so reply trust is computed via Authorize but always
   delivered (falling back to external trust).

6. **Runtime groups in a separate namespace.** Config groups are immutable
   at runtime because a runtime edit would silently revert on restart.

7. **Conversation ID owned by the sender.** The interviewer creates the
   conversation ID, because only it knows what constitutes one conversation.

8. **Package-level vars converted to functions.** `defaultQuestions` in
   the interviewer and `CipherSuite` in noise.go were both converted from
   package-level `var` to functions after the ch5 grader caught them.
   Third occurrence of this pattern in the book.

## Unverified Claims in the Chapter

- "it hasn't failed since I started it in the spring of 2026" — Bill's
  claim from his story. No independent verification available. Faithfully
  reproduced per §7.8 (our own claims stay in print with results attached).

- The Cloudflare TLS termination scenario assumes a standard reverse-proxy
  deployment. The description is architecturally correct for any TLS-
  terminating proxy, not just Cloudflare.
