# Chapter 10 — final bounded operational handoff

**GUI exercise PASS; authorized paid CLI recovery PASS. Root operational acceptance, tag, and prose closure remain pending.** No additional execution is authorized or needed.

## Frozen candidate and authority

Source `729bc919e734405307b4ac0060e65c40e33652e5`, tree `025ef3744b42c8f9daec34e76661edc40583ca6a`, tracked status clean. CLI SHA256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33`; GUI SHA256 `fb506081f98c09d2e6da628edbb6cb922f067f9db24983ae856d2d1222dcbf53`. Production source/assets/binaries were not changed or rebuilt. No completed suite/grade was repeated and no tag or main commit was made.

Prospective root instructions authorized the local modes and the separate bounded browser corrections. The successful quiet-state browser receipt was then carried by exact receipt, original seal/policy/manifest, source, and unchanged flow-file pins. After the initial paid failure, root explicitly authorized one paid-only recovery using the installed public CA and at most three new reservations. Full authority/diffs and exact argv/environment/timings are in `authority-and-static-checks.json`, `paid-cli-tls-r1.receipt.json`, and the sealed bundle. No private root authorization/protection manifest was read.

## A — owned browser exercise

`browser-quiet-r1`: PASS in 2.747s, complete owned cleanup, zero credential opens. One headed owned Chrome/profile/page, one GUI restart, dynamic loopback ports, three fake-provider requests. DOM card IDs were `1,2,4,3` before EOF and after restore, then `1,2,4,5,6,3` after the new final. Read-aloud and microphone were inactive. Recorded speech events were only silent cancels—not utterances. Typed read/result cards, old/new IDs, replay and stale-card assertions passed.

The pre-EOF and restored PNG hashes are identical (`2bf0119066a5e4e64dd3f7ab0d9b36b46599266fc66ff0deb2cf6c95cafdca7e`); the new-final PNG differs (`eaa0fabad54df5806eddeb1c34ae5a9f1a67ad9005a8b8644b5bea3a485fa3fc`). **PNG files were captured and hashed, not human-viewed or heard.** No human visual/audio inspection is claimed. Exact carried receipt, DOM-file and local-only PNG pins are in `final-operational-facts.json`; these were not rerun or recopied. Prior browser failures and the narrow observed-favicon classification remain preserved, not relabeled.

## B — paid attempts and recovery

The initial attempt remains **FAIL**, `proxy-transport-failure`, 0.780s, one reservation and one authorized credential open. Saved counters are `forwarded=false`, `connect_calls=0`; no provider response/usage or underlying exception type was retained. These counters do **not** prove no attempted DNS/TCP/TLS handshake. A local-proxy CLI 502 is not a provider HTTP status. Successful recovery does not retroactively establish the old cause.

The single TLS recovery is **PASS** in 5.660s. Both CLI launches exited 0, no signals were needed, their recorded groups were absent, and all owned resources cleaned up. The same work/save path was used across EOF/relaunch; startup request counts were checked before admitting each prompt. Actual `read_file` call/result provenance and a nonempty resumed answer were verified from the real CLI/provider flow, not fabricated history. Archive `as_of` advanced from **9 to 13**, with the retained prefix checked. The first-save hash is a contemporaneous phase observation; that same archive path was intentionally updated by the follow-up. Final archive and history checks are pinned in `final-operational-facts.json`.

Endpoint: `POST https://api.anthropic.com/v1/messages`; model `claude-sonnet-4-6`; native `max_tokens=1024`, validated rather than clamped. All three incoming body hashes equal their forwarded-body hashes:

| Slot | HTTP | Header seconds | Input tokens | Output tokens | Request body SHA256 |
|---|---:|---:|---:|---:|---|
| 1 | 200 | 2.056644 | 1772 | 102 | `e5e126f92c66073a0e33b1634879f3ac4c9f4c50af045ef1601b2c1158dd26ba` |
| 2 | 200 | 1.102918 | 1895 | 20 | `65c561bed8947c3845de41f5329e1dcd96413a82ecd1dec0b2c3fe91331a1c5f` |
| 3 | 200 | 0.976344 | 1949 | 17 | `05a832c4525a10fdc72d2cd464fcb438d1baeb7f92c86640a1d903ce545ba5a2` |

Provider-returned recovery totals: **5616 input / 139 output**, cache creation/read input fields **0 / 0**. These are observations for the three successful recovery responses, not a cost estimate or aggregate billing claim for the unresolved first attempt. Response message IDs, tool-call correlations, stop reasons, response/body hashes, replies and per-call usage are in `final-operational-facts.json`. HTTP request-ID headers were not retained.

## TLS, budget and credential boundary

Verified public CA: `/etc/ssl/cert.pem`, 333483 bytes, SHA256 `9dae8d76e55cb08991f2b672d58999ea15560d910759c16b544f843bdffbb994`. Static context loading observed 128 CA certificates; certificate verification and hostname checking remained enabled. The CA pin is revalidated before credential/upstream entry. Only transport exception **type**, not message/args/headers/key values, was added; no transport exception occurred in the recovery. The original missing exception detail is not reconstructed.

The first reserved attempt counts conservatively: **1 prior + 3 recovery reservations = 4**. This does not assert four completed provider requests: exactly three recovery HTTP 200 pairs were observed. The recovery cap was three total across both launches, with no retry, alternate model/provider, warmup or diagnostic request. Two authorized key-file opens were recorded across the two paid attempts—one per proxy process, cached across that attempt’s launches. Child argv/environment used only a nonsecret proxy placeholder. This is Python-proxy/owned-process evidence, not an OS-wide key-open or packet-capture audit.

## Preservation and limits

All old reports/seals/failures remain unchanged. Earlier NOTEXECUTED/blocked reports describe their original timestamps. Local rehearsal receipts are referenced, not rerun. The v1 CLI argument error, prior browser failures, initial paid failure, and static/metadata mistakes (including the malformed CA comparison) remain recorded. No production correction was made.

Raw API traces, archives, profiles, screenshots, binaries and the earlier real TAR remain local-only/excluded from Git, referenced by hash, not copied here. `final-preservation.json` verifies the frozen candidate, 144 prior sealed members and 254 older immutable entries. `inventory.json` is the exact current UTF-8/no-NUL text inventory, including this report and excluding itself. All owned jobs have finished. Root alone decides final acceptance/tag/text.
