# Edition 3 Chapter 5: independent frozen-adapter review

## Verdict

**ACCEPT the exact frozen adapter identified below.** No blocking adapter finding and no correction request. This review establishes that the thin Edition 3 command exposes the existing Chapter 5 API faithfully, preserves Chapter 1/4 routing, rejects invalid invocation before target launch, and is sensitive to real forwarding errors. It does **not** accept a Chapter 5 student implementation or convert the historical score into complete current-contract coverage. Student launch remains the parent's decision after accepting this review.

This is a new bounded adapter review, not a repeated entry audit. The preceding corrected entry report is `ch05-entry-audit.md`, SHA-256 `9a5e43deaee5614207be97841361fd4e0ec085855a3f485a45cefdd5825662e9`; its 88 covered evidence files were already verified. The earlier transport interruption and completed shutdown did not invalidate that submission and are not test failures.

## Frozen identity and scope

Main HEAD: `686c1ef027c32b69d5c44230004b42d59045989a`.

| Candidate input | SHA-256 |
|---|---|
| `cmd/grade-edition3/main.go` | `a7a9e4b4e499b15a4887e9c10be4bae1ee1241d76b5ce301befdb0189113be2c` |
| `cmd/grade-edition3/main_test.go` | `c25ca42eadf65702153b0f28746ed54e8f8c46609106962c3b2c90fbc766b6c6` |
| Compiled `/tmp/cadence-ch05-grader` | `ffd18a1b7b4b179c7a5fd84108b322346da9dfa8d1f405f8a4916cf78a3f8250` |

Accepted student predecessor: `ec4a350f564e820d398f9125fd73a9b3b27129cd`, tag `edition-3-ch04-r1`. Student source and CLI were controls, not new Chapter 5 work.

The exact diff was inspected and saved as `ch05-grader-review-evidence/candidate.diff`. Only the command and its test file differ from main HEAD. No inherited grader packages, student source, chapter, or governing document changed. Prior required-read and source hashes were checked rather than blindly rereading the governing corpus; the only changed file in that manifest was the expected command adapter. The entire tracked historical Chapter 5 scratch source was byte-compared with main HEAD, the retained student source inputs matched their recorded hashes, and both reused CLI binaries and the credential-free launcher matched the prior artifact hashes.

Every new compiled route used a scratch copy whose bytes equal the frozen `/tmp/cadence-ch05-grader` hash. The independent command tests used an isolated source copy containing the exact candidate command/tests and unchanged inherited packages. All binaries, launchers, private case/model inputs, source copies and sentinel files stayed outside Git. The review changed no adapter, grader, student, book or guide file.

## Source review and actual API

The added branch is deliberately narrow:

- Chapter selection accepts 1, 4 or 5; exactly one nonempty target is still required.
- Explicit flag visitation rejects `workflow` and `media-models` on Chapters 1 and 4 even when their values are empty. Existing Chapter 1 rejection of Chapter 4-only flags remains.
- Chapters 4/5 validate explicit core and consumer module directories; Chapter 5 additionally requires an explicit workflow module directory, each with a regular `go.mod`.
- `ReadConsumerCase` and, for Chapter 5, `ReadMediaModels` run before `grade.Build`. Invalid files/contents therefore cannot launch the target.
- The unchanged `grade.Build` receives only the positional target. The Chapter 5 call is exactly `gradeedition2.Chapter5(bin, *core, *consumer, *workflow, example, models)`.
- Chapter 4 retains its prior `Chapter4(bin, *core, *consumer, example)` call. Chapter 1 retains `gradeedition3.Chapter1(bin)` and its real lifetime probes.
- Existing cleanup, JSON/text reporting, and exit policy remain: help 0, invalid/setup 2, failed report 1, passing report 0. No historical defaults or implicit source inference were introduced.

The schemas were read from the actual unchanged implementation, not guessed:

- Consumer-case JSON has `name`, object `arguments`, nonempty `result`, and nonempty `diagnostic`; the reader returns `gradeedition2.ConsumerCase`. Empty path uses the existing shout demonstration.
- Media JSON is `map[string]string`, with exactly the keys `anthropic`, `openai`, `gemini`, each with a nonblank identifier. Empty path uses the existing declared-capability fixture rows. The reader validates the mapping, not live model availability or media capability; the runtime check establishes support/refusal.

Smallest complete current invocation for this layout, with optional files omitted when their defaults actually match the submission:

```sh
grade-edition3 -ch 5 -json \
  -core "$CORE_MODULE" -consumer "$CONSUMER_MODULE" \
  -workflow "$WORKFLOW_MODULE" \
  "$CORE_MODULE/cmd/ch01"
```

A prebuilt executable can replace the final argument, but must be independently tied to the same frozen source. To choose explicit grader inputs, place `-consumer-case "$CASE_JSON" -media-models "$MEDIA_JSON"` before the positional target. The adapter now separates build/executable identity from core source identity; it does not alter `grade.Build` to turn a library root into a runnable CLI.

## Newly executed controls

Raw stdout/stderr and exact argv/cwd, elapsed seconds, exit code and outer-watchdog status are saved per named run. Paths are sanitized using `<REVIEW_SCRATCH>`, `<ENTRY_SCRATCH>`, `<MAIN_REPO>`, `<STUDENT_REPO>`, `<USER_HOME>` and `<SYSTEM_TMP>`. Durations below are wall-clock measurements, not callback intervals.

| Run | Seconds | Result |
|---|---:|---|
| `candidate-command-tests` | 3.427 | `go test -race -count=1 ./cmd/grade-edition3`: exit 0 |
| `candidate-ch1-lifetime` | 124.048 | exit 0, 100/100, Passed=true; both actual 62-second header/body probes pass |
| `candidate-ch4-regression` | 8.768 | exit 0, 160/160, Passed=true on accepted predecessor, explicit core/consumer/case |
| `candidate-ch5-positive-retry` | 12.205 | exit 0, 100/100, Passed=true on the private historical executable/core/consumer/workflow, explicit non-default case and media mapping |
| `candidate-ch5-predecessor` | 9.018 | exit 1, 35/100, Passed=false; corrected Chapter 4 parity, inherited replay and common cleanliness pass; workflow/media fail as expected |
| `sensitivity-media-map-retry` | 12.063 | exit 1, 90/100, Passed=false; only `loud-refusal` fails for the selected unsupported image-model row |
| `sensitivity-consumer-case-retry` | 12.003 | exit 1, 90/100, Passed=false; only Chapter 4 parity fails on the deliberately wrong complete tool-result expectation |
| `sensitivity-core-identity-retry` | 8.716 | exit 1, 35/100, Passed=false; consumer and workflow resolved-module checks reject the substituted core directory |

The Chapter 1 run had a 240-second outer watchdog and completed normally. All other grading controls had a 90-second outer watchdog. No new watchdog expired. No short tool callback was treated as a deadline.

### Real positive and forwarding sensitivity

The positive used the verified historical source's actual external module paths (`solutions/ch05/ch05` consumer and `solutions/ch05/ch06` workflow), not mocks replacing their behavior. The existing grader exercised HTTP-backed tool dispatch, returned consumer results, workflow observations, replay, actual image bytes and refusals. The case selected a different text/result from the default demonstration. The media map selected other already-declared image-capable rows for all three providers. The 100-point positive is a private adapter-routing control, never a substitute for student work.

The three sensitivity controls changed **valid inputs**, not the grader or target implementation:

1. A structurally valid map selected an unsupported Anthropic model. Runtime failed specifically with `model "unknown-review-model" does not support image media`; other Chapter 5 checks passed. This shows the map reached media rendering instead of being dropped for defaults.
2. A valid case changed only the expected result. Runtime observed the genuine correlated text `INDEPENDENT ADAPTER REVIEW` and rejected the deliberate different expectation. This shows the case reached the consumer check instead of being ignored.
3. The executable and real historical apps stayed fixed while `-core` selected the verified predecessor. The consumer and workflow each resolved their core to the historical directory, not the supplied predecessor, and failed the appropriate checks. This establishes forwarding of the explicit core into both identity checks, rather than inference from the executable.

These are genuine post-validation grading failures, not compile failures. Workflow forwarding is additionally observed through the predecessor run: a source-less scratch module was explicitly supplied because the accepted predecessor has no workflow. The grader examined that exact path and reported `no Go files`; it did not substitute the ordinary CLI or a historical workflow. That placeholder is solely a negative baseline input, not student implementation or a workflow-positive claim.

### Compiled invalid-input/no-launch matrix

All **31** sentinel cases passed, with no target marker created:

- 26 invalid invocations returned 2: no/empty/extra target, unsupported chapter, unknown flag, arguments after the positional target, missing core/consumer/workflow, empty/missing/nonmodule workflow, bad/missing/incomplete media mapping, bad/missing/incomplete consumer case, and both empty and nonempty new flags on each of Chapters 1 and 4.
- Five help forms returned 0: `-h`, `-help`, and `-ch 1|4|5 -h target`.

`guard-summary.json` records expected/actual exit and the sentinel-marker observation for every case. Per-case streams and exact invocations are retained. This tests the real frozen executable, not only calls to `run(args)` inside unit tests. The candidate's source-level tests also passed independently under `-race`.

## Retained mistakes, incomplete attempts and carried evidence

No new behavioral run was incomplete. Four first attempts used my wrongly guessed historical `examples/consumer` and `examples/workflow` paths: `candidate-ch5-positive`, `sensitivity-media-map`, `sensitivity-consumer-case`, `sensitivity-core-identity`. Each returned 2 in 0.005 seconds at module-directory validation. All four records remain. Their corrected `-retry` runs are the real controls above. The failed path-listing inspection was the same reviewer setup error, not an adapter defect. The negative path attempts are not credited as forwarding sensitivity.

The parent's two earlier compile failures are copied unchanged except path sanitization: wrong guessed `run` signature/import, then nonexistent `gradeedition2.MediaModels`. They are retained as setup history, not negative grading controls. Parent race/vet/build/diff-check results and candidate record are retained separately; this review independently reran the command race tests and exercised the exact frozen executable.

Carried, not newly rerun: the entry audit's historical lifetime failure, public explanation-policy tests, full predecessor suite retry, batch-limit control, executable/source provenance control and broader source-only coverage analysis. Their recorded failures/incomplete attempts remain in `ch05-entry-evidence/`. This review verified the reused inputs and cited that evidence rather than blindly repeating it.

## Acceptance limits and next boundary

- This adapter deliberately reuses unchanged historical grader checks. A Chapter 5 100/100 result does not establish the current default **200 dispatch batches**, the consuming request after the last allowed dispatch, reliable private Ask replies under saturation/concurrency, actor ownership, pending-hint HTTP interleavings, or actual manuscript effects/handoffs. Those entry-audit coverage gaps remain; this review creates no new framework or speculative requirements.
- Chapter 1 lifetime, Chapter 3 public explanation policy, and Chapter 4 jobs/public-consumer obligations remain inherited obligations. Routing Chapter 5 does not silently waive them. This review newly establishes the lifetime and Chapter 4 routes on the same candidate executable; public-policy and full-suite positives are explicitly carried evidence.
- Core/application resolved-directory checks do not authenticate a prebuilt target executable to those sources. Freeze and record source plus binary identity for future student acceptance.
- No student implementation has been frozen or reviewed here. A future independent student acceptance needs that exact candidate and separate contract evidence, not the private historical control or this adapter verdict.
- No credentials, paid calls, live GUI, children, installs, restarts or pushes were used. Loopback fake providers and credential-free bounded launchers were used throughout. Original unrelated work was preserved.

The requested narrow entry adapter is sufficient and accepted at these hashes. No adapter correction or optional polish is requested. The review ends here; there is no pending test or indefinite wait.
