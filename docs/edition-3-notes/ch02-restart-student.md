# Chapter 2 — restart student notes

Implemented from the accepted Chapter 1 workspace, retaining the star: Agent owns the session, History owns causal facts/replay, Engine owns disposable vendor projections and live transport accounting, and common owns shared contracts. The file CLI and separate minimal GUI use that core. No source changes followed the final build.

## Teaching feedback

The explicit exercise spellings and replay rules were valuable; typed events alone do not prevent a reducer from accidentally erasing failed attempts. The hardest distinctions were historical versus live usage, actual producer identity versus a routing alias, and an unknown usage operand versus measured zero. A small worked example with an invalid cache count would make the dependency rule especially concrete: unknown cache subtraction makes fresh input unknown, while independent output remains measured. Likewise, ordinary Gemini text carrying a signature must remain portable even when the signature cannot replay. Review findings, not reference code, helped expose those mistakes.

The inherited Chapter 1 tests required narrow, root-authorized assertion updates for Chapter 2 causal dialogue, independent measured buckets, offline credential-free construction, and four-bucket JSON. The missing-key guard before live HTTP remains tested. Initial tooling mistakes were a build run from the wrong directory and a script anchor mismatch; neither was claimed as successful validation. One focused Gemini fixture initially used valid JSON arguments while expecting a parser error; changing it to a genuinely invalid text field fixed that test without a production change.

## Evidence and limits

All raw logs and the sole frozen executable are under `/tmp/ch02-final-validation/`; executable: `frozen-cli`.

- Focused Round 2 regressions passed after that single fixture correction; Round 3 dependency-accounting tests passed all 18 cases. Earlier focused logs are `/tmp/ch02-round2-*` and `/tmp/ch02-round3-*`.
- Initial root race run exited 1 on stale Chapter 1 assertions; the authorized test-only update and one rerun exited 0. Root vet exited 0.
- Separate GUI module `go test ./...` and `go vet ./...` each exited 0. GUI has no authored tests; no browser exercise was run.
- Final CLI build exited 0. Initial Chapter 2 grade exited 1, parsed `passed:false`, 85/100; original `grade-ch02.json` and stderr are preserved. Root repaired two grader assumptions; the independent reviewer—not this student—ran corrected Chapter 2 once on the unchanged executable: 100/100, `passed:true`, exit 0 (`grade-ch02-corrected-reviewer.*`). No student Chapter 2 rerun.
- Remaining Chapter 1 grader ran on that same executable: 100/100, `passed:true`, exit 0 (`grade-ch01.*`). Both real 62-second header and body gates passed; observation waits did not cancel the run.

## Authorized live exercise

The unchanged frozen CLI ran `chat` for three turns per provider: invent a codename, answer an unrelated arithmetic question, recall the actual codename. A root-authorized in-memory credential proxy enforced fixed verified-TLS origins, one discovery GET/provider, native 1024 output caps, and three POSTs/provider. CLI environments received only dummy keys. The one run exited 0 in 20.05 seconds: exactly 3 discovery GETs and 9 paid POSTs, no retries/fallbacks. All CLI exits were 0. The proxy independently checked exact outgoing dialogue content at 1/3/5 messages and compared all four saved usage totals with native responses.

| Discovered model | Recalled codename | Input | Cache write | Cache read | Output |
| --- | --- | ---: | ---: | ---: | ---: |
| `claude-sonnet-5-5` | Quillfire-7 | 325 | 0 | 0 | 31 |
| `gpt-4.1-2025-04-14` | FrostPine | 222 | 0 | 0 | 9 |
| `gemini-3.8-flash` | VoxelDrift | 175 | 0 | 0 | 655 |

Raw sanitized output: `/tmp/ch02-final-validation/live/live.log`; per-provider chat output and causal logs are alongside it. One actual Anthropic log passed keyless `dump` and two byte-identical `render` invocations per vendor (all exit 0; no network). The owned proxy processes shut down and joined before exit.

Limits: these were text conversations, not live tools/media/signature or nonzero-cache tests. Chat's status display retains the inherited input/output-only line and `claude>` label for all vendors; four-bucket accounting was verified from native responses and saved facts, not claimed from that display. No browser, commits, or snapshot changes. Required assistance was the independent reviews, test-only rebaseline authorization, root/reviewer correction of grader assumptions, and bounded live credential-proxy authorization. This note records validation, not independent acceptance.
