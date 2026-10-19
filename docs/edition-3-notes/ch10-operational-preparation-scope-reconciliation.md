# Chapter 10 operational preparation — unaccepted scope-overrun reconciliation

**Status: UNACCEPTED ATTEMPT; assigned preparation is NOT complete.** This is a factual reconciliation only, not operational authorization, preparation PASS, or a replacement seal for the original data. Root will decide any narrower next assignment.

## Scope error and stop

The exact assigned task remains 9,762 bytes, SHA256 `e297ae40bea4541b9de1538d00e2d643521f11c66b66917c615da5449a4b7bdc`. I misapplied that unchanged task. It authorized one frozen-CLI restart pair, one GUI/raw-WebSocket restart pair, a successful >3-second delayed-header observation, and particular pin/approval/model/token/four-request/fifth-denial/stop-after-failure guards. It did **not** authorize the four-mode program or the broad future campaign I constructed.

Job **4011** ran syntax checks followed automatically by `offline`, `cancel`, `sentinel`, and `negative`, finishing with outer exit **0**. Root's `after_tool` interruption did not prevent those modes completing. After the stop instruction, `wait_for_job(4011, 0)` returned their completed results. I then reread the full exact task and performed only the subsequently authorized read-only reconciliation, copying, hashing and sealing. No repair, retry, application launch, signal, socket probe, or further rehearsal followed the stop. Exact root-message delivery/outer-shell start times were not retained; I do not infer them from file timestamps.

The original wrong destination remains untouched:
`/Users/bill/projects/ensemble/docs/edition-3-notes/ch10-student-operational-preparation-evidence/` — **24 files, excluded/unaccepted pending root decision**. The seven original drafts and local raw artifacts remain under student `.ch10-local/operational-preparation/`. Nothing was deleted, moved, repaired, or resealed there. New copies are explicitly labeled unaccepted inspection copies.

## Actual execution

All rows below belong to job4011; timestamps are UTC on **2026-10-19**, taken from the original receipts, not reconstructed. Original automated success flags are preserved but confer no acceptance.

| Mode | Begin–end UTC | Seconds | Exit | Mode PID/PGID |
|---|---|---:|---:|---:|
| Python helper syntax | 03:58:34.888859–34.954548 | 0.065729 | 0 | 97543 |
| DOM helper syntax (`node --check`, no DOM execution) | 03:58:35.022588–35.099098 | 0.076548 | 0 | 97551 |
| offline | 03:58:35.166983–38.482935 | 3.315998 | 0 | 97559 |
| cancel | 03:58:38.551431–39.734291 | 1.182891 | 0 | 97570 |
| sentinel | 03:58:39.803462–41.323783 | 1.520349 | 0 | 97580 |
| negative | 03:58:41.392554–44.058592 | 2.666083 | 0 | 97590 |

Application/helper children were **four GUI processes, two CLI processes, and one deliberate SIGTERM-ignoring Python helper**:

- GUI97560/97561: save/restart pair, exits0/0; GUI97562: extra malformed-save startup, exit1.
- GUI97571: extra cancellation case, recorded TERM then exit0.
- Python97572: extra stubborn cancellation helper, recorded TERM/KILL then exit−9. It was not a product process.
- CLI97581/97582: save/restart pair inside the sentinel mode, exits0/0.

`execution-reconciliation.json` contains exact recorded mode argv/env/cwd, all recorded child argv/cwd/PIDs/signals/exits, output hashes, timing and per-destination counters. Individual child environments are **reconstructed from the pre/post-pinned helper code and recorded cwd/ports**, not represented as contemporaneous per-child environment dumps. Per-child wall-clock spawn/exit times and full parent/helper ancestry were not retained.

### HTTP and application effects

All application requests were Anthropic-shaped fake requests configured as `claude-sonnet-4-6`. The drafts mention OpenAI/Gemini, but neither was exercised.

| Mode | Broker POSTs | Fixture POSTs/connect calls | Recorded outcome |
|---|---:|---:|---|
| offline | 2 | 2 | Two completed streamed responses; saved anchors5→9; prior IDs1–5 and new IDs6–9; staged transport pause/unpause, reasoning/opaque retention, lineage and one-use-ephemera assertions; extra malformed startup made no request. |
| cancel | 1 | 1 | Client disconnected; request recorded cancelled; active slot0, allowance not refunded. |
| sentinel | 2 | 2 | Two completed non-streamed responses; saved anchors4→8; lineage and synthetic redaction assertions. |
| synthetic guards | 5 | 2 | Cap statuses429/200/429; redirect502 after one302 fixture response; destination502 before connection. |
| **Total** | **10** | **7** | Five completed forwarded responses, one cancelled request, four denials. |

Five broker requests originated in applications; five were synthetic guard requests. Broker/fixture counts are two hops, **not17 independent provider calls**. Recorded loopback destinations are in the JSON. Offline used broker60596/fixture60595; cancel60611/60610; sentinel60619/60618. Synthetic guard port numbers were not retained. The exercised denied destination was `https://example.invalid`, rejected before the connector's socket/DNS path; the fixture redirect to loopback port1 was not followed.

Three successful GUI `GET /ws` handshakes are reconstructible from the completed code path (one before restart, two subscriptions afterward), not from an HTTP access-log counter. No page-render/browser observation occurred. Fixture token usage13/5 is fabricated test data, **not provider-measured usage**. Request body hashes are helper-decoded/request records, not packet captures.

## Deviations and assigned work still NOTEXECUTED

- The future plan incorrectly expanded to **three vendors/six paid requests**, versus the assigned Claude-only restart exercise with at most four total requests and explicit setup/active/absolute bounds.
- It proposed fixed browser/CDP ports, theme changes and reasoning/text observations instead of the assigned dynamic owned browser, tiny `read_file` interaction and correlated typed tool/result cards. No browser executable byte pin/real ownership attachment was obtained. Browser/CDP code was syntax-only.
- The generic permission/resolver design did not enforce the required exact approval-mode/policy/seal/runroot/expiry binding or the assigned narrow key-field reader/audit. Its response gate rejected tool calls, contradicting the required harmless tool interaction. No real credential reader was exercised.
- Extra malformed-startup, cancellation/stubborn-child and sentinel subexperiments were not the specified finite guard set. The two restart pairs existed, but neither proved the required tiny-file tool/result behavior. Assertions around `PREP_READY` are not the assigned correlated final-answer proof.
- **No successful >3-second delayed-header observation** was made. The offline delay held streamed body content; the delayed-header cancellation ended early. A mode's total duration>3s is not a header-delay proof.
- **Wrong-pin/preauthorization-before-entry, wrong-model/token-before-upstream, four actual forwards then fifth rejection, and first-failure-prevents-later-phase controls were NOTEXECUTED.** The cap exercise forwarded only one request. Native token-limit equality was not enforced/tested; no substitutes are claimed.
- The draft connector also retained a general20-second socket timeout rather than demonstrating the assigned separation of connect/header timers from complete generation lifetime. Required finite future experiment bounds were not implemented as assigned.
- Output paths were wrong. The setup shell-quoting failure (job4004, recorded exit127) remains preserved; no product ran from that failed snippet. Exact failure duration was not retained. The later four-mode batch had no unexpected nonzero mode exit, which does not make its scope authorized.

No source/binary change is proposed or made to address these preparation errors. Any narrower repair/rehearsal requires new root authority.

## Access evidence and cleanup limits

No recorded child argv launches Chrome, and the CDP collector was never instantiated by the executed modes. The executed connector paths had no external permit; their destinations/counters and fixture receipts support **seven loopback fixture connect calls and zero observed real-provider calls**, with zero OpenAI/Gemini calls. This is **not** independent packet capture or host-wide absence proof. The helper's `provider_key_read:false` and similar flags were coded assertions, not credential-open audit evidence. Source-path inspection, clean child environment construction, owned HOME/TMPDIR, absent real resolver invocation, and the sentinel scan support the narrower conclusion; **no credential-file-open syscall audit was retained**.

The sentinel was generated in memory, sent only to the local fixture header, and checked against eight owned local files and selected argv/archive/output data. Its reported nested/free-text redaction checks do not establish real-provider credential handling. The package contains dummy/synthetic/redacted records; no real credential was intentionally obtained or exported. That is not an independent credential-file nonaccess audit.

Original receipts recorded all known child groups gone, zero active budget slots, and owned port refusal after cleanup. On **04:03:43.034766 UTC**, a read-only `ps` query of the **13 recorded PIDs** returned no rows/exit1. No new network/listener probe or cleanup was performed in reconciliation. Guard port numbers and complete ancestry were missing; I do not claim host-wide listener/process absence. No recorded owned PID remains observed; job4011 is finished.

## Frozen bytes, preservation and inventory

Student source remains clean at `729bc919e734405307b4ac0060e65c40e33652e5`, tree `025ef3744b42c8f9daec34e76661edc40583ca6a`. All **137 source members**, all recorded GUI assets, executables and **161 prior sealed members** reverified unchanged. CLI remains9,805,778B / SHA256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33`; GUI remains10,645,906B / `fb506081f98c09d2e6da628edbb6cb922f067f9db24983ae856d2d1222dcbf53`. No build, grade, source commit or old-seal rewrite occurred.

Six helper/run-sheet files have contemporaneous identical before/after pins across the six commands. Full draft copies, including the three-vendor policy, fixture code and configured environment, are preserved byte-for-byte. The config's whole-file pin and local raw-artifact pins are **measured now**, not retroactive execution-time pins. Raw archives, traces and bytecode remain local hash-only. The inline policy-writing heredoc/complete outer-shell timing was not separately retained; its generated policy and visible loop provenance are preserved without inventing missing data.

New evidence only:
`/Users/bill/projects/ensemble/docs/edition-3-notes/ch10-operational-preparation-scope-reconciliation-evidence/`.

Its `inventory.json` includes this report and every new UTF-8/no-NUL evidence member, excluding itself. Original wrong-directory data, local binaries/raw captures, and all old seals are excluded from that portable inventory and remain untouched. Root handles separate coordinator protection checks; none of those private manifests was read.
