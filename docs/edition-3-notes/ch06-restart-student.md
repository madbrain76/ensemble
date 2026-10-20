# Chapter 6 restart — final local PASS; manual evidence with Gemini cap caveat

**HOLD for root/Bill's final disposition.** All 3 formal correction rounds and both explicit narrow exceptions were accepted by root/same reviewer. Source stayed unchanged during final gates/manuals. Latest local grades pass; manual execution reached its hard budget, but Gemini's hinted final was cap-limited and incomplete. Do **not** treat the runner's blanket success label as an unqualified semantic/manual pass. No further paid request or production change is authorized by this record.

Owned source: `solutions/edition-3/main`, based on accepted Chapter 5 `bd4e49387160d80d813584972e65502a2dce6dde`; this is the sole student note. Raw/build evidence remains under `/tmp`. No children, commits, snapshots, manuscript/grader changes, browser/audio/device use, future/reference/private-grader implementation reads, or credential discovery.

## Accepted work and carried checks

Engine owns shared plain/SSE parsing; Agent owns History/save-before-final/private replies; transient fences and lossless observations serve CLI/public shape consumers. Capability fallback/opt-out, sanitized trace, native part IDs/provenance/usage and accepted Chapter1–5 contracts retained.

- R1 fixed native boundaries/signatures/IDs/refusal, call-record-before-effects barrier, replaceable ephemera and Gemini error envelopes; README updated. Evidence `/tmp/e3-ch06-round1/logs/`.
- R2 forwards real state boundaries to JSON CLI. R3 makes tool-name/argument Chunk streams appendable and avoids duplicate human labels. Evidence `/tmp/e3-ch06-round{2,3}/logs/`.
- First Bill exception: shared OpenAI renderer uses documented `max_completion_tokens:1024`, not rejected max_tokens; no cap increase/removal/model framework. Other providers unchanged. Official reference: https://developers.openai.com/api/reference/resources/chat/subresources/completions/methods/create . Real local HTTP plain/SSE field/usage tests, affected race/vet and five builds passed; `/tmp/e3-ch06-token-exception/logs/`.
- Second Bill exception: local stable OpenAI wire-only complete-result projection. Call→hint→real-running-result remains in History; wire becomes call→result→hint. IDs/order/content preserved, no fabrication/rerun/status upgrade/global sort; incomplete/unrelated groups unchanged. Actual Agent/local HTTP and multi-group/opaque/portable-text regressions failed before and passed after; 10 distinct focused top-level races plus subcases, affected vet and five builds passed. `/tmp/e3-ch06-adjacency-exception/logs/`. Edit/fixture setup failures are separate from behavioral failures.

Carried, not rerun for receipts: main `go test -race ./...`, `go vet ./...`; separate `go test ./...` and `go vet ./...` in gui/ch06/ch07/examples/ch04-consumer all exited0 previously (`/tmp/e3-ch06-final/logs/`). GUI/ch07/consumer had no test files: compilation-only test gates. Focused exception checks cover changed paths. No redundant full-suite campaign.

## Fresh final gates — same adjacency CLI

Clean offline sequential wrapper: GOMAXPROCS=2, GOFLAGS=`-mod=readonly -p=1`, GOTOOLCHAIN=local, GOPROXY/GOSUMDB=off, explicit `/Users/bill/go/pkg/mod`, existing Delve PATH, owned HOME/TMP/cache; no inherited real keys or LLM_BASE_URL. Disk stayed >1GiB; no shared cleanup.

Exactly one final build:
`go build -o /tmp/e3-ch06-final-adjacency/frozen-cli ./cmd/ch01` —0.

From repository root:
```sh
go run ./cmd/grade-edition3 -ch 6 -json -core /Users/bill/projects/ensemble/solutions/edition-3/main -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 -stream-model claude-sonnet-4-6 -media-models /tmp/e3-ch06-final/media-models.json /tmp/e3-ch06-final-adjacency/frozen-cli
go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch06-final-adjacency/frozen-cli
```
Map: `{"anthropic":"claude-sonnet-4-6","openai":"gpt-5.4","gemini":"gemini-3.8-flash"}`.

Both **process0, parsed passed:true,100/100** on the same binary. Chapter1 took **125s**, genuinely completing62s headers and62s body with full answer/usage, totals511input/65output. No short outer timeout. Logs: `/tmp/e3-ch06-final-adjacency/logs/` (`frozen-cli-build.log`, `ch06-grade{.json,.stderr.log,-status.log}`, `ch01-grade{.json,.stderr.log,-status.log}`).

Earlier binaries/results preserved: `/tmp/e3-ch06-final/` failed20 (boundary/launcher), `final-r2/` failed85 (tool-name chunks), `final-r3/` both grades passed, `final-exception/` both passed before adjacency repair. No older binary was substituted.

## Final manual execution and limits

Root explicitly authorized ≤8 additional POSTs after6 used, hard **14 cumulative including failures**, no retries/fallback/discovery/Claude repetition. Stages: OpenAI human≤3, Gemini human≤3, GPT streamed/plain shape≤1 each. All latest8 were sent, all HTTP200, all child exits0; runner exit0. **Budget exhausted:14/14; no further calls.** Each active case <120s; latest campaign ~26s, below600s.

Latest invocation, once in clean env:
`/usr/bin/python3 /tmp/e3-ch06-manual-adjacency/control/run.py --live-approved` —0.
Only approved runner-version changes: pins/output directory and cumulative accounting6→14. Earlier runners/receipts immutable. Setup anchor mismatch wrote nothing and made no request. Shape built once from held source; same frozen CLI used for both human sessions. No proxy body rewrite, redirect, artificial response hold or retry.

Fixed verified-HTTPS endpoints/system trust; exact directOpenAIAPIKey/directGeminiAPIKey fields used only in forwarding-process memory/headers. No Claude lookup in replacement runners. Dummy-key clean children, controls/HOME/TMP outside tool cwd, explicit marker-path exception, owned child/proxy cleanup completed. Native per-request caps remained1024: OpenAI max_completion_tokens; Gemini maxOutputTokens; carried Claude max_tokens. No image/browser/audio/device execution: Gemini's SVG was merely returned text, not rendered.

### Fresh human observations

| Case | POSTs / elapsed | Actual producer | Measured input/output | Outcome |
|---|---:|---|---:|---|
| OpenAI, requested gpt-5.4 |3 /6.55s|gpt-5.4-2026-03-05 on all replies|1799 /248|Complete transport and human exercise; no400 after adjacency repair |
| Gemini, requested/known target gemini-3.8-flash |3 /16.46s|gemini-3.8-flash on all replies|3092 /1712|Delivery/hint/job evidence obtained; hinted final incomplete at native cap |

Cache write/read0/0; no missing usage buckets. Requested/known target and actual returned producer kept distinct; no manufactured OpenAI provenance.

OpenAI first real stdout **1.543s < upstream EOF2.682s**; marker waiting2.714s, hint acknowledged/released2.730s. Next request valid roles `system,user,assistant(tool_calls),tool(real running result),user(hint)`, HTTP200. Final summary adopted the blue-lantern direction and honestly qualified its running observation; recall likewise did not invent observed completion. Native finishes: tool_calls,stop,stop.

Gemini first real stdout **4.462s < upstream EOF4.476s**; marker waiting4.496s, hint acknowledged/released4.512s. Hint is present in the next request and influenced the blue-lantern output. Both sessions preserve a stable paired job handle1, initial running observation, later saved terminal done/exit0/26retained bytes; each real `result.txt` contains `FILE_JOB_OK` plus newline. One job each, no rerun or fabricated completion. Automated routing/no-duplicate checks passed; final observations are not substituted for independent terminal job facts.

**Gemini caveat / review limitation:** its second native finish reason is **MAX_TOKENS**, with957 thought +63 visible tokens (1020 output), yielding a137-character partial SVG rather than a complete factual summary. Native visible text exactly equals saved/delivered text: no demonstrated parser loss. First/third replies end STOP; recall accurately describes the earlier running observation. The runner's permissive blue-content/exit checks labeled this `chat_pass` and did not classify the native cap stop as failure, so recall and shape requests ran before offline inspection identified the incomplete final. This limitation is explicitly reported; no unqualified Gemini complete-final/manual PASS claim, no retry, cap increase, patch or additional call followed. Root/Bill must decide whether the cap-limited evidence satisfies the manual requirement.

### Shape pair — each generation checked against itself

GPT-5.4, actual producer gpt-5.4-2026-03-05. One POST each, child0, native finish stop each, each usage51input/42output/cache0/0.

- Streamed:41 published text callbacks (2 empty),1 final part,50 published content bytes; first delta1344ms, final elapsed1590ms.
- Plain:1 callback,1 final part,50 bytes; first delta1042ms, final elapsed1045ms.
- Each native response's aggregate exactly matches **its own** saved/private final. Native prefixes all match that final; streamed has nonempty strict prefixes. Published nonempty byte-length sequence matches native content sequence. Trace stage `delta_publish` distinguishes consumer publication from separate `parse_delta` telemetry.
- Independent generations happened to return identical numbers1–20, but that equality is **not** used as parser-parity evidence. The example reports counts, not persisted callback text; offline proof uses captured native content, saved final and published trace lengths. Earlier offline filename/count/stage assumptions were analysis errors, not product failures, and were corrected without app/provider reruns.

Raw latest evidence: `/tmp/e3-ch06-manual-adjacency/control/` (`campaign.log`,`status.log`, case request/response bodies, stdout/stderr, History). Workdirs contain API traces and actual effects. `offline-review.log` records native finishes/text preservation; `shape-offline-review.log` records corrected response-local checks. Superseded count/stage assumptions retained separately. No new harness/smoke/diagnostic execution campaign.

## Preserved prior live evidence and failures

- `/tmp/e3-ch06-manual-r3/control/`: Claude3POSTs, child0,13.18s, actualclaude-sonnet-4-6; firststdout1.320s<EOF3.809s; hint/paired real job and no duplicate final;3555input/522output, cache0/0. Carried without repetition because OpenAI-only exceptions leave its path unchanged. Initial GPT400 rejected max_tokens; no usage/producer claimed for error.
- `/tmp/e3-ch06-manual-exception/control/chat-openai/`: first GPT request200 with actualgpt-5.4-2026-03-05 and433input/114output; second400 on hint/tool-result adjacency. Original causal History, real effects and exact324-byte error/request preserved. No automatic patch/retry then.
- Total provider POSTs14 includes those2 HTTP400s; failed calls supplied no measured usage, so no invented zero/billing claim. Successful measured usage is recorded per case above.

## Final pins and HOLD

Current CLI `/tmp/e3-ch06-final-adjacency/frozen-cli`: `dc1f9aaeddb20e830829e726fed741b84035192bc8cf4b4d096198a36a27a5b0`.
Shape `/tmp/e3-ch06-manual-adjacency/control/ch07-shape`: `0a1eab7a90e1042796c34cf0ac0dd889911e182b618ad374dc5e2e7f52e3f3ee`.
Pins rechecked unchanged after execution. Prior runner hashes remain `d7b5eef16e4a2dc0d0a379a45880a29170792063550d1d3f5d04a0cf9750f895` and `5255664b58f091350f2f6b80d502631134d4946ff40bffa4758bfe0bfa2271cd`; previous binaries/failure outputs retained.

**Saved final + HOLD:** local gates pass, repaired OpenAI live adjacency succeeds, and response-local shape evidence passes. Gemini's capped incomplete hinted final and the runner's late recognition remain explicit acceptance caveats. No overall acceptance claim, no further source edit or paid execution. Any new production work or renewed budget requires Bill/root authorization.

## Final closure — root/reviewer accepted scoped Chapter 6 evidence

Root and the same reviewer accept the scoped evidence **with the Gemini limitation**, superseding the pending-disposition wording above. No further work is authorized.

- Current CLI: `/tmp/e3-ch06-final-adjacency/frozen-cli`, SHA-256 `dc1f9aaeddb20e830829e726fed741b84035192bc8cf4b4d096198a36a27a5b0`. Shape: `/tmp/e3-ch06-manual-adjacency/control/ch07-shape`, SHA-256 `0a1eab7a90e1042796c34cf0ac0dd889911e182b618ad374dc5e2e7f52e3f3ee`.
- Fresh grades on that same CLI: Chapter6 and Chapter1 each100/100, parsed passed:true, process0. Chapter1 took125s with genuine62s headers+62s body and complete usage. Earlier full suites/focused races/vet carried, not repeated.
- Final live total **14POSTs including2 preserved HTTP400 failures**. Fresh OpenAI human, Gemini human and GPT shape pair used8HTTP200 requests; current OpenAI adjacency works. Claude's prior successful3-request session is carried unchanged, not repeated. Native1024 caps stayed intact. Own-response shape aggregate/prefix evidence is accepted; independent generation equality is not used as proof.
- Gemini's hinted reply ended MAX_TOKENS: **all137 native characters were delivered and saved, but the natural-language summary was incomplete**. This is not demonstrated parser loss. The runner failed to flag that cap stop and continued through recall/shapes before offline review; its blanket pass label is qualified accordingly. No retry followed.
- Owned children/proxies cleaned up; controls stayed outside tool cwd, credentials confined to authorized forwarding memory/verified headers, dummy child keys. Prior binaries/runners/400s/Claude evidence preserved. No additional source edits, calls, analysis or receipts after closure authorization.

**Saved final + HOLD.** Scoped Chapter6 evidence accepted with the explicit incomplete-Gemini-answer qualification; no claim of full natural-language completion and no further execution budget.
