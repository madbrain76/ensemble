# Chapter 2 candidate evidence index

Candidate: `74eb0265bd8bdbc7857d579704300b317b2365da`.

Read [`../ch02-student-review.md`](../ch02-student-review.md) for architecture, full command/cwd/status tables, assistance, development failures, and limitations.

## Final candidate

- `final-identity.txt`: source archive hash, CLI hash and Go embedded build identity.
- `final-source-checks.txt`: final root vet, diff check, clean status and matching hashes.
- `final-grade-e2-ch02.{json,stderr,exit}`: Chapter 2, exit 0 and passed:true.
- `final-grade-e3-ch01.{json,stderr,exit}`: Chapter 1 preservation, exit 0 and passed:true, including long response checks.
- `final-race.{txt,exit}`: full source race tests, exit 0.
- `final-cover.out`, `final-coverage-*`, `final-coverage.exit`: source-only 745/839 statements (88.8%). Separate GUI module and grader/live processes are not included.
- `final-gui-checks.txt`: isolated placeholder build/test, no GUI runtime claim.
- `final-mutations.txt`, `mutation-{ephemera,cache-write,opaque-isolation,anthropic-order,ref-stub}.txt`: five targeted behavioral failures followed by restored passes and matching source hashes.

## Live and offline capture

- `run-live-evidence.py`: original authorized bounded launcher. **Do not rerun merely to verify artifacts; it performs live calls.** It accessed only the three authorized key fields in memory, never stored keys, and made six generation calls total.
- `live-prompts.txt`: the two meaningful human-mode prompts used with each provider.
- `live-{anthropic,openai,gemini}.{stdout.txt,stderr.txt,jsonl}`: actual CLI transcripts, diagnostics, normalized saved event logs.
- `live-summary.json`: original summary, preserving its erroneous usage-regex boolean.
- `live-launch.{txt,exit}`: **exit 1**, preserving the evidence-only final assertion failure. All three CLI child exits were 0; no production failure is hidden.
- `verify-captures.py`: **offline-only** correction/verification over existing artifacts, no key access or additional live calls.
- `capture-verification.{json,txt,exit}`: corrected checks, **exit 0**; per-response and cumulative usage match saved facts; all log hashes unchanged.
- `offline-summary.json`: three exact dumps, nine byte-identical repeated cross-vendor render pairs (18 fresh processes), valid JSON, all exits 0, no API keys, local sentinel POST count 0.
- `offline-*`: actual dump/render stdout and stderr artifacts.
- `fake-redaction-*`: explicitly synthetic before/after JSONL and rendered requests, not live tools. Original fake result remains in the append-only after-log but not in any post-redaction request.

## Earlier failures / non-final development evidence

- `initial-grade-*`: initial Chapter 2 failure and initial Chapter 1 run, not final-candidate proof.
- `development-grade-*`: later Chapter 2 failure and development pass, not final-candidate proof.
- `development-*`, `precommit-*`: source/check/coverage development artifacts. Do not use these to claim final source coverage.
- `mutation-opaque-isolation-initial-no-tests.txt` and `mutation-*-precommit-harness-false-negative.txt`: preserved mutation harness mistakes.
- `public-openai-cache-write.txt`: independently fetched public documentation excerpt supporting optional cache-write normalization.

All outputs are sanitized student evidence. There are no raw credential stores, auth headers, or raw provider request/response archives. The executable `student-cli` is the exact final grader/live candidate. `gui-stub` is a separately built placeholder, not a live GUI application. The manifest hashes retained artifacts without changing them.
