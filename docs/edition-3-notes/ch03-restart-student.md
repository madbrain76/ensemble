# Chapter 3 restart — final student report

**Authorized validation is complete; source and the final CLI are held for root acceptance.** Three formal production rounds were used. The last Gemini completion was harness-only: no source/test/binary change or grade replay.

## Implementation and corrections

Extended the accepted Chapter 2 owners with a per-Agent registry joining six declarations and handlers, Engine-owned wire declarations, and synchronous ordered dispatch. Empty registries stay empty. Explanation/budget refusals pair results before effects; batch 16 gets its final consumption request. Tests use real files/processes and actual Agent/Engine local-HTTP exchanges.

- **Round 1:** joined registry entries; removed four API-key variables from actual shell children; range-first reads with 32,768-byte default/truncation marker; distinct stdout/stderr/exit status; cross-file search separators; size-bearing overwrite refusals.
- **New user requirement, not a formal round:** edit diagnostics count overlaps, show at most eight one-based start lines plus omitted-location counts, and guide rereading. Refusals preserve bytes; unique replacement/deletion work.
- **Round 2:** the first executable Chapter 3 grade scored 95 because edit errors lacked file/anchor identification. Adding the supplied path fixed it without anchor excerpts.
- **Round 3:** a separately authorized diagnostic established Gemini's rejection of `additionalProperties`. Engine removes only that top-level keyword from fresh Gemini parameter maps, preserving shared schemas, other provider rendering, and strict local validation.

Other preserved failures: initial compile mistakes; a new test's wrong `enabled`/`disabled` expectation; an inherited assertion counting OpenAI declarations as calls (narrowed to actual message calls); and an isolated-cache grader setup failure (resolved using the approved offline cache). Two guessed focused names matched nothing; the final full race run covered the actual suite.

## Final local gates

Temporary root **`/tmp/e3-ch03-early.jR0AOc`**; frozen CLI **`round3-final/ch03-final`** beneath it.

- Root `go test -race ./...` and `go vet ./...`: exit 0.
- GUI test/vet: exit 0; **tests are compile-only, with no test files**.
- From repository root, `go run ./cmd/grade-edition2 -ch 3 -json <final-cli>`: **100/100, passed:true, exit 0**, including all 13 carried Chapter 2 checks.
- `go run ./cmd/grade-edition3 -ch 1 -json <same-final-cli>`: **100/100, passed:true, exit 0**. Genuine 62-second header plus 62-second body delays; 124 seconds elapsed under a 360-second observation window, without cancellation.

Exact logs: `round3-final/{root-race,root-vet,gui-test,gui-vet,final-build,grade-ch3,grade-ch1}.log`. Clean sequential environments used the required Go limits, owned HOME/TMP, and approved offline cache. No gates were replayed after the harness-only identity correction.

## Live results and qualifications

| Session | POSTs | Measured input/output tokens | Observed result |
| --- | ---: | ---: | --- |
| Claude Sonnet 5.5 original work | 10 | 33692 / 2022 | All six tools; real repair, unchanged test, test exits 1→0; cap blocked distinct recall. |
| Claude separate clean read/recall | 4 | 6929 / 204 | Three completed turns; actual third-turn recall. |
| OpenAI GPT-4.1-2025-04-14 | 6 | 4498 / 156 | Repair, unchanged test, test exit 0, paired calls, third-turn recall. |
| Earlier Gemini attempts | 5 | 1559 / 209 on successful responses only | Schema failure/diagnostic, then unknown-affinity signature rejection; error usage unmeasured. |
| Final resolved-identity Gemini 3.8 Flash | 6 | 6730 / 463 | All HTTP 200; CLI exit 0; three completed turns; four ordered pairs; actual `helo`→`hello` repair, unchanged test, test exit 0, third-turn recall without tools. |

**Total 31 POSTs, within the final authorized ceiling of 33.** Final Gemini elapsed 13.68 seconds under eight-POST/180-active/240-absolute bounds. Exact-route/model/token guards, atomic reservations, TLS verification, native 1024-token limits, and no redirects/retries applied. Real keys stayed in proxy memory/HTTPS headers, never CLI/tool environments or logs. All proxies/children exited.

The initial Claude workspace exposed its prompt file, and its launcher's recall flag was a false positive. That failed attempt remains qualified; later clean recall is a separate session. Claude/OpenAI evidence is from the round-2 binary, carried only by verified unchanged provider paths—not fresh final-binary proof.

The first Gemini 400 body was unavailable; a separately authorized diagnostic established the schema defect. After its fix, an unknown-affinity continuation lacked required thought signatures. Root and reviewer independently confirmed that known identity adds only the original signature keys structurally. **Both my launcher and the original root authorization omitted resolved producer identity.** The final authorized harness supplied `LLM_RESOLVED_MODEL=gemini-3.8-flash`, the observed exact producer on the same vendor/surface/route. The unchanged frozen CLI then succeeded. This was not a fourth production round or automatic retry.

All failed evidence remains under `live/`, `live-v2/`, `gemini-diagnostic/`, and `live-gemini-final/`. Retained error/request files are mode 0600. Final evidence: `live-gemini-resolved/run.log` and its `gemini/control/{session.jsonl,chat.stdout,chat.stderr,result.txt}`. Clean-session controls/logs/HOME/TMP stayed outside workspaces; final Gemini's workspace contained only its three exercise files.

## Author lessons and assistance

Teach Gemini's JSON-Schema subset and known producer identity explicitly: common schemas are not universally accepted, and conservative opaque replay must not be misdiagnosed when the harness omits identity. Keep controls outside tool workspaces; verify the actual third response and actual file/test effects. Passing local graders is not a substitute for provider-strict live evidence.

Root and one independent reviewer supplied public findings and bounded authorizations. I did not inspect reference solutions, grader source, runtime-agent source, or sibling answers. The existing owners accommodated the changes without a new framework.
