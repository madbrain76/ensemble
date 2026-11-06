# GUI transport security: what was considered, built, and removed

Investigation of 2026-10-22, conducted by the `ensemble debugger` agent while
debugging local-model output delivery. Repository HEAD at the time of the
discussion was `0327827a88de31533dc6bbbb9da7f717a9e1db26`. Nothing described
here is committed: the digest-auth implementation was built, tested, verified
live, and then removed at Bill's instruction. It was never committed, so
**the code is gone** — this note is the only surviving record of the design.

Bill's request began narrowly — make the GUI reachable from the LAN rather
than loopback only — and became a survey of every way to authenticate a
browser to a private server. The narrow part is done and committed
(`--host`, `agent/gui/cmd/main.go`). The rest is this document.

## Environment as established during the discussion

These facts drove every conclusion. They were supplied by Bill, not inferred.

| Constraint | Consequence |
|---|---|
| GUI is a **web UI**, clients are browsers (Firefox; Android phone wanted) | Client is not our code; no custom handshake, no bespoke pin store |
| Client machines are **containers** Bill controls fully | No app-store distribution; no "just install an app" |
| **`GOPROXY=off`** — no new module downloads | No `lego`, no third-party crypto; stdlib and vendored modules only |
| Hostnames are **private, unpublished**; `.localdomain` | `.localdomain` is not a real TLD: **no CA may ever issue for it** |
| **No CA distribution to clients** | Rules out private roots deployed by hand |
| LAN contains devices Bill **does not trust**; others are **adversarial**; "no device is trusted" | On-path attacker assumed; cleartext disqualified |
| **Unmanaged switch, no VLANs** | Phone and container are on the same L2; **pfSense is not in that path** |
| No Apple hardware | The Darwin-only sandbox path is permanently unreachable |

Two of these invalidated recommendations I had already made, and are recorded
below as retractions rather than quietly dropped.

## Measurements, not recollection

Claims below were tested in this container. Where a claim could not be tested,
it is labelled as reasoning.

| Question | Result | Evidence |
|---|---|---|
| Can Go's `crypto/tls` serve **external** PSK identities? | **No** | Go 1.26.3; `go doc crypto/tls.Config` exposes PSK only for resumption; GOROOT grep finds no `PSKVerify`/`GetClientPSK`; sole PSK identifier is `alertUnknownPSKIdentity` |
| Can Firefox present a PSK identity? | **No client-facing path** | `libnss3.so` exposes only the resumption binder label `psk_id_hash`; `omni.ja` matches are the CSS colour `deepskyblue` and the `PskKeyExchangeModes` extension *type name* — neither is an API |
| Does the wire support PSK at all? | **Yes, TLS 1.2 only** | OpenSSL 3.5.5: `Cipher is PSK-AES128-CBC-SHA256`, `PSK identity: svc`. My first TLS 1.3 attempt reported `PSK identity: None` on a *resumed* session — PSK had never negotiated |
| Does the vendored Noise module do PSK? | **Yes** | `flynn/noise@v1.1.0` `state.go`: `MessagePatternPSK`, `psk` token, `willPsk`, `PresharedKeyPlacement`; placement folded into the protocol name, so a stripped-PSK downgrade fails cryptically |
| Is the GUI sandbox enforcing anything here? | **No** | `internal/sandbox/command.go:61` returns "unsupported on this platform" unless `runtime.GOOS == "darwin"`; `(deny network*)` sits at line 82, behind that gate |
| Are GUI tool sessions confined? | **No** | `agent/gui/cmd/main.go:144` sets `BuiltinTools = true`; `SafeMode`/`SandboxRoot` exist but default false/empty, and an empty root means full access |

Scope caveat: the Go and OpenSSL results are directly observed. "Browsers
cannot present a PSK identity" is a reasoned claim from the shipped build's
exported symbols plus the absence of any API or preference, not an exhaustive
wire probe.

## The negative result that ends the design space

A browser will honour a TLS server only if the certificate chains to a root it
already trusts, or the user accepts an exception. That is the complete list.
Combining it with the table above:

- **Public CA** — needs either inbound reachability (HTTP-01) or a publicly
  authoritative zone (DNS-01). Bill has neither, and `.localdomain` is not
  issuable under any circumstances.
- **Private CA** — requires placing a root in each client's trust store.
  Excluded.
- **TLS-PSK / Noise with a PSK** — excluded twice: browsers cannot present a
  PSK identity, and the code that would verify a PSK must itself be delivered
  over the hostile channel.
- **Application-layer crypto over cleartext (PAKE, SPAKE2+/OPAQUE, WebCrypto)**
  — excluded independently, for two reasons. `crypto.subtle` exists only in
  secure contexts, and `http://192.168.x.x` is not one. And even with WASM
  crypto, an on-path attacker rewrites the page before the handshake runs.
- **Network layer (WireGuard, SSH tunnel, TLS offload at pfSense)** — excluded
  by Bill. This was also factually wrong on my part twice over; see retractions.

The residual conclusion is not "hard" but **impossible**: an out-of-band trust
fact must reach the client by a channel the attacker does not control. Every
proposal that appears to avoid this has smuggled such a channel in and called
it cryptography. I checked for a fourth option and there is none.

## Retractions

Recording these because each was stated confidently and each was wrong.

1. **"WireGuard is overkill for a home LAN."** Wrong under the stated threat
   model. That judgement was conditioned on a *trusted* LAN; Bill then said
   devices are adversarial, which inverts it.
2. **"Put WireGuard on pfSense."** Useless here. On a flat L2 behind an
   unmanaged switch, phone→container traffic is switched at layer 2 and never
   reaches pfSense. I recommended a control that was not in the path.
3. **"20+ characters and rate limiting is appropriate for LAN use."** I agreed
   with this while assuming a trusted LAN. With an on-path attacker, digest
   authenticates only the *client*, so the attacker presents their own login
   form and relays to the real server — full impersonation, and every prompt
   and response readable in transit. The long secret is irrelevant because the
   attacker never needs to crack anything. Digest-over-cleartext is
   opportunistic-attacker protection only.
4. **"Enabling streaming would disable the thinking-only retry guard."** I
   justified leaving streaming off using `anthropic_output_retry_test.go` — a
   test for a surface Bill does not use. Bill runs `--vendor openai`. My
   probes had all been run with `--vendor anthropic`, and the model registry is
   keyed by name alone, so the flag moves both paths. The guard exists only on
   the Anthropic renderer; there is no equivalent on the OpenAI path at all.
5. **My own `git checkout` regression.** Backing out that option-A change
   deleted the `delta.reasoning_content` parsing from `openai_stream.go` and
   `openai.go` — the code producing Bill's thinking boxes on the interface he
   actually uses. I made the revert while still believing he was on Anthropic.

## What was built, and why it is gone

RFC 7616 digest authentication was implemented to the standard Bill chose
after reading the cost analysis:

- `SHA-256` only; MD5 never advertised, downgrade refused
- `qop="auth"`, server-generated single-use nonces with expiry, `nc`
  monotonicity plus a `cnonce` replay cache
- Replay state committed **only after** credentials verify, so a rejected guess
  cannot burn a challenge a legitimate client is mid-flight holding
- Constant-time comparison; unknown users indistinguishable from wrong
  passwords; client `realm`/`uri`/`method` validated so one signature cannot be
  replayed against another path
- **≥20-character secrets enforced at startup**, not warned about
- Per-source and per-name throttling with exponential backoff and
  `Retry-After`, cleared by a successful login
- Applied at `Server.Handler`, so the **`/ws` upgrade is authenticated before
  the handshake** — otherwise the socket is a bypass and the guard is theatre

Verified then: 86 tests passing, 1 environmental failure
(`TestExternalReuse`, which needs `dlclark/regexp2` under `GOPROXY=off`); a
known-answer vector computed **independently in Python's `hashlib`** so the
tests could not pass by agreeing with their own implementation; and a live
probe sequence — `401` with a well-formed challenge, `/ws` refused
pre-upgrade, valid digest admitted, replay refused, MD5 downgrade refused,
cross-path signature refused, then a guessing burst throttled
`401 401 401 429 429 429`.

Removed on instruction ("just rip out all authentication support", scoped to
the digest layer only; provider OAuth in `internal/auth` and gateway peer-token
auth were explicitly left intact).

## Crack cost, for the record

Measured rates matter more than adjectives. Per RFC 7616 §3.2 a digest guess is
roughly two hashes, `A1` cacheable across captures.

| Secret | Keyspace | Time, one GPU (~3 GH/s, SHA-256) |
|---|---|---|
| rockyou dictionary (14 M) | 1.4e7 | milliseconds |
| 6 chars `[a-z0-9]` | 2.2e9 | < 1 s |
| 8 chars `[a-z0-9]` | 2.8e12 | ~15 min |
| 4 random words | 3.6e15 | ~2 weeks |
| 12 random `[a-z0-9]` | 4.7e18 | ~50 years |
| 16 random `[a-z0-9]` | 3.4e25 | irrelevant |

Two aggravators: `A1 = H(user:realm:password)` is **unsalted**, so tables
precompute per `(user, realm)`; and the cheaper attack is **online** guessing
against the server's own challenges, which needs no capture at all and is
limited only by throttling. Algorithm choice is a factor of ~5–10; secret
entropy is orders of magnitude.

## Recommendation that survived

Transport crypto cannot be retrofitted onto a browser session over a hostile
segment, so the question that still has an answer is **what a stolen credential
gets you**. That is the boundary worth hardening, and it needs no PKI:

1. **Least privilege as the real boundary.** LAN-originated sessions should
   default to no shell, confined workspace, and approval gates on destructive
   tools. Today `BuiltinTools = true` and an empty `SandboxRoot` mean full
   access; the only network denial ever written is Darwin-only. A credential
   leak should yield a read-only assistant, not a shell on a host with a route
   into the house.
2. **Egress policy on the tool path, not the container.** The container
   legitimately needs LAN access; nothing says a model-requested shell command
   does. `unshare -n` per command, or a seccomp filter blocking `connect()` off
   loopback, needs no added capabilities in this container.
3. **Per-device credentials** — attributable, individually revocable and
   rate-limited — instead of one shared secret.
4. **Fail closed at startup**: refuse a non-loopback bind with no guard loaded.
   This behaviour was implemented and verified in the removed code and is worth
   keeping as a principle.

Order of precedence, stated plainly: encrypting the browser hop is worth less
than constraining what an Agent can do with a shell, because an attacker who
can reach a web search does not need to read the wire to make the Agent open
the connection for them.

## Spec references consulted

RFC 9110 §11 (auth framework) · **RFC 7616** (digest; §5 says it is not a
substitute for TLS) · RFC 7615 (cookie-managed digest) · RFC 7617 (basic +
charset) · RFC 6750 / **RFC 9700** (bearer; TLS required) · RFC 9449 DPoP ·
RFC 8705 mTLS · RFC 9421 message signatures (machine-to-machine only; browsers
cannot sign) · RFC 5802 / RFC 7804 SCRAM · **RFC 9381 SPAKE2+ / RFC 9382
OPAQUE** (the principled answer to password auth without TLS, and the reason a
PAKE was considered for the browser — killed by the secure-context and
attacker-writable-page arguments above).

## Follow-ups worth deciding

- Whether to re-add any authentication at all, given the conclusion that
  digest over cleartext is cosmetic against an on-path attacker.
- Whether to ship items 1–2 above, which protect against the realistic
  compromise and touch no transport.
- Whether the private hostnames can move to RFC 8375's `home.arpa`; it does
  not enable TLS, but `.localdomain` is invalid and `.local` collides with
  mDNS.
- Whether owning one cheap domain for ACME **DNS-01** is acceptable after all.
  It is the only route to browser-trusted TLS with zero client-side work —
  the CA reads a TXT record and never contacts the network, and the A record
  can stay split-horizon private. It was excluded, but it was excluded for
  "hostnames are private", which DNS-01 does not actually forbid.
