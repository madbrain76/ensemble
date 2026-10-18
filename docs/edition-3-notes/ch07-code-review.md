# Edition 3 Chapter 7 — initial whole-code review

## Verdict and authority

**CHANGES_REQUIRED: two public-contract defects require correction before whole-candidate acceptance.**

This is the initial independent whole-code review of frozen student commit `48e1d12d9c396b5988c0ce6eccdef67e581872ff`, tree `0f5f7fbe4273e91e1be3aa13d0715a99f671eaf9`. It is not another entry, adapter, pause-diagnostic or fixture-repair review. Those four accepted packages remain unchanged. Formal correction rounds remain **0/3** at submission; root alone issues any correction packet, accepts, commits or tags.

Required changes:

1. **R1 — High: the GUI falsely reports dispatch for a batch the core truthfully refused before admission.** This occurs live and on replay.
2. **R2 — Medium: a tool-result card's speaker omits the readable tool result.**

A measured page-hide/unpause behavior is retained as **C1, a non-blocking student-task deviation/contract-clarity note**, not a third required correction. Root explicitly qualified that finding after checking the student-visible task and the less explicit chapter wording. No new code requirement is inferred from the private entry audit.

## Frozen inputs and scope

- Main HEAD: `23c75ede7944d9a912ae74d3f683354d163f0c86`.
- Student: the commit/tree above; clean at final verification.
- Immutable accepted predecessor: `edition-3-ch06-r1^{commit}` = `f595c787675e2a6343053ec42ea6c7e1e02a930e`.
- Current Chapter 7 SHA-256: `68d02d8e1170c2912676b8f162847f8fac25c87216514ed06adbd5f0ea011b47`.
- Repaired pause driver SHA-256: `bcdbcd50397d40de60c9fa871cd9d588403efeee34162aa01830cf6c99913b09`.
- Corrected grader SHA-256: `c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e`.
- Frozen CLI SHA-256: `4beda0aa805195229c5987a025a2b794b4fe3e77a6526ff194927282fddd85cb`.
- Frozen GUI SHA-256: `39ebcd3e42a57d8972b6278fb8cc66021e38afd41afa2cd230b93499e6b13a81`.
- Initial review-input manifest: `/tmp/cadence-ch07-code-review-input-pins.json`; its exact copied bytes and all referenced hashes are in evidence.
- New owned scratch: `/tmp/cadence-ch07-whole-code-review-independent`.

The private source copy was verified against all **117** source-manifest hashes before execution and again afterward. No original source file was changed. One reviewer-only GUI test file was added to that copy. Tests generated local logs/job reports, inventoried separately by hash and not packaged as payloads. No production source, main/student guidance, manuscript, grader or runtime file was edited. No live server, current browser, native speech/microphone, paid provider, credential store, service restart or deployment was used. No child agent was launched and no cleanup was performed.

## Reads and whole-code coverage

The current AGENTS, coding skill, reviewer overlay, execution/readiness/student/carryover documents and complete Chapter 7 were reloaded. The historical workflow, Edition 2 carryover, voice and writing procedure were also reloaded in full. Cap-clipped sections were recovered; the carryover 801–1000 span was finally reread as two complete halves. Evidence records actual file hashes, lengths and ranges, rather than substituting a summary for these reads.

Reviewed the complete **26-file delta** from the accepted Chapter 6 predecessor and relevant surrounding actor, dispatch, history, public configuration/interface and GUI code. This includes all new GUI assets and support tests, not just the grade-facing functions. The README, submitted review/checklist, exact source/evidence manifests, discrepancy reconciliation and accepted native/listening records were inspected. The final two actor lines and clipped review-report overlap were recovered. Necessary inherited clauses concerning explanation-before-dispatch, truthful paired undispatched results, hint observation wake/job continuation and public callback-delay naming were checked.

This is not a claim to have freshly audited every unchanged implementation line in all 117 files or re-proved the whole book. The accepted predecessor and prior independent audits are identified as carried evidence.

## R1 — false dispatch projection for a refused batch

**Severity: High. Required correction.**

### Source and contract

`gui/hub.go` translates every durable `tool_called` record into `tool_dispatched` in `capture` and retains that projection for replay (notably the `tool_called` branch around lines 89–94). That is not sufficient evidence of actual admission.

The inherited actor deliberately appends a truthful `tool_called`/error-result pair when `RequireToolExplanation` refuses a response without prior explanatory text: `internal/agent/actor.go:605–617`. It then retries the model without dispatching that batch. This pairing is not a claim that a job ran. Real admission has its separate `ensemble.ToolDispatched` observation in `internal/agent/dispatch.go`.

Chapter 7 requires actual content and state rather than IDs or invented lifecycle observations, truthful tool dispatch/result events, matching job/call identities, and a pause gate measured between actual starts. The Chapter 5/6 explanation/undispatched obligations remain in force. A truthful refused pair must remain visible without becoming a fabricated start.

### Fresh measured evidence

A private test compiled successfully against the unchanged public root and actual GUI hub. It used a local fake provider and the same valid `run_command` in two cases, with an actual filesystem effect as the independent admission observation. Streaming was disabled for this narrow explanation/refusal comparison, not replaced with a browser mock.

| Case | Real core `ToolDispatched` | Actual command effect | Live GUI `tool_dispatched` | Replay GUI `tool_dispatched` |
| --- | ---: | --- | ---: | ---: |
| Explained positive | 1 | present | 1 | 1 |
| Refused before admission | 0 | absent | **1** | **1** |

The refused case also retains the honest error result: the batch was refused for lacking explanation and **no tool in it ran**. Thus the failure is in the GUI projection, not evidence that the core ran a refused command. The explained positive validates that the observation/effect path can detect a real admission.

Compilation exited 0. The compiled test then exited 1 on the intended projection assertion, not on startup, compilation, timeout or disk failure. Exact command, clean environment, hashes, timing, stdout and stderr are retained under `logs/compile-refusal-projection.*` and `logs/refusal-projection.*`; `reviewer_contract_test.go` is reviewer-only evidence.

### Required outcome

Keep the requested call and truthful refused result, with stable identities, but do not label it as dispatched unless admission actually occurred. Preserve correct actual-dispatch observations and their replay de-duplication. Cover both live and replay behavior. No new diagnostic API, arbitrary framework or private fixture interface is required.

## R2 — tool-result speaker omits its readable result

**Severity: Medium. Required correction.**

### Source and contract

In `gui/web/gui/app.js:29–32`, the card speaker enqueues `a.readable`. The `tool_dispatched` branch sets that field to a short tool/argument summary. The `tool_finished` branch around line 65 updates the result node but does not update the card's speech text. Consequently the displayed result is not what its speaker reads.

The public requirement is explicit: Chapter 7's artifact contract gives artifacts a speaker, §7.10 says its click reads the **readable contents**, and §7.12 requires a clickable speaker for **every artifact** (`chapter-07.md` around lines 154, 422 and 553). Readable text must be distinguished from opaque data, but the measured missing content is ordinary text, not a signature or opaque blob.

### Fresh measured evidence

A non-browser VM executed the unchanged production page, speech queue, phrase queue, formatter and scroll modules. Its DOM/WebSocket/native-speech boundaries were inert fixtures; this was not a new browser/native exercise.

The page received a real-shaped `tool_dispatched` followed by `tool_finished` for the same call, with ordinary readable result `DISTINCT_TOOL_RESULT_41`. The result node contained that text. Clicking the actual card speaker called the production speech path with:

```
read file: note.txt
```

The utterance omitted `DISTINCT_TOOL_RESULT_41`. The completed harness exited 1 and reported `passed=false` for the tool-result speech check. This observes the actual text delivered to speech synthesis, not whether a human heard it; it does not negate the accepted native listening positive for other text.

The first attempt stopped at a missing inert DOM method, `createDocumentFragment`. That setup failure is retained and not counted as a detection. Adding only the missing inert DOM/text-node methods allowed the same cases to run against unchanged application source. See `logs/page-contract-checks-dom-stub-fixed.*` and the exact two harness revisions.

### Required outcome

A tool-result artifact's speaker must include its readable result (including readable error text where applicable), rather than only the tool name/argument summary. Preserve the safe/opaque distinction and production queue/error/cancellation behavior. Do not satisfy this by removing the speaker or by claiming a native listening test of an unrelated text artifact covers it.

## C1 — page-hide release: measured, but not a blocking code correction

The same unchanged page harness measured this sequence:

1. A nonempty draft's input event sends `pause`.
2. `pagehide` runs while the socket is open.
3. The handler at `app.js:77` cancels speech, clears local causes/draft, invokes `localPause`, then closes the socket.
4. **`unpause` is queued before close.**

This is a real generated-control observation in the non-browser page harness. It is not a fresh native measurement of a remote gate after browser teardown, and it does not establish that every transport delivers the queued frame before disconnect.

The exact student-visible assignment `/tmp/cadence-ch07-student-task.txt`, SHA-256 `8d34f80dfce2dfc506a2dfd1b3f3a18129279227846a738914fd85619d882d2f`, line 9, included `disconnectnotauto-release/reconnectstale causesnotreassert`. Therefore this is not merely an invented private-audit preference.

However, the chapter itself explicitly describes server-held state, automatic local-cause reversals and stale reconnect behavior without unambiguously resolving this page-hide cleanup transition. Root checked the exact task and directed that, under Bill's prohibition on added requirements, this be recorded as a **student-task deviation / public-contract clarity note, not a third blocking correction solely on that sentence**. That disposition is retained in `contract-classification.md`. No public text or source was changed here. Root owns any clarification and subsequent correction scope.

A separate positive in the same harness also matters: after another tab unpauses, ending speech while an unchanged draft remains **does not reassert pause**. The earlier cause-removal hypothesis was disproved and is not a finding.

## Fresh checks and exact limits

All commands were sequential. There was no outer watchdog on the required full root race run, which was allowed to complete its genuine long waits. The later full root gate subsumes the earlier selected core checks; that overlap is recorded rather than hidden.

| Fresh command/check | Result | Wall time |
| --- | --- | ---: |
| Selected changed core/inherited tests, `go test -race -count=1` | exit 0 | 2.37 s |
| Pristine GUI module, `go test -race -count=1 ./...` | exit 0 | 1.42 s |
| Submitted production queue tests in Node | exit 0 | 0.03 s |
| Refusal projection compile | exit 0 | 0.96 s |
| Compiled explained/refused comparison | explained passes; refusal projection fails; exit 1 | 0.01 s |
| Page checks, initial inert DOM | setup exception; exit 1; not a detection | 0.03 s |
| Page checks after stub repair | speech failure; qualified page-hide observation; stale-cause positive; exit 1 | 0.03 s |
| Full root `go test -race -count=1 ./...` | exit 0 | **139.72 s** |
| Root `go vet ./...` | exit 0 | 0.28 s |
| GUI `go vet ./...` | exit 0 | 0.17 s |
| `examples/ch04-consumer`, race/vet | exit 0 / 0 | 0.12 / 0.10 s |
| `examples/ch06`, race/vet | exit 0 / 0 | 1.68 / 0.11 s |
| `examples/ch07`, race/vet | exit 0 / 0 | 0.12 / 0.10 s |

The actual module layout is preserved: library root, nested `cmd/ch01`, separate importing GUI, external consumer, workflow and stream examples. `examples/ch07` remains the stream-shape example, not a replacement GUI. The consumer and stream-shape modules report no test files, so those commands establish build/vet coverage, not new behavior cases. The workflow module has its own passing test. No module was flattened.

Go test/vet commands do not emit grader score/Passed fields; the raw receipts correctly retain those as null. The page harness emits its own `passed=false`, not a normalized Chapter 7 grade. These results must not be described as a fresh 100/100 whole-candidate verdict.

## Carried evidence, not new claims

- The accepted corrected Chapter 7 route's **100/100, `passed=true`** and real compiled pause-bypass/restoration remain valid carried evidence for the same source pins. This review did not repeat that campaign. The historical student's pre-repair **80/100, `passed=false`** remains preserved, not rewritten.
- The exact frozen CLI's Chapter 1 receipt was inspected: `chapter1-r4-invocation.json` identifies CLI `4beda0…`, original Chapter 1 plus Edition 3 lifetime companion, the recorded offline environment, exit 0, 124.154179 seconds. The report says **100/100, `passed=true`**. This is carried r4 evidence, not a fresh long-gate invocation here. Its invocation and report are pinned in the review inputs; the invocation records the exact CLI and grader-source hashes.
- Root's accepted native/browser/listening package remains a positive for the actual cases it exercised. No browser/native/microphone action or paid comparison was repeated. That positive does not establish the untested tool-result speaker content or refused-dispatch projection.
- Passing parser/queue/race/replay tests do not prove every native rendering, listening or causal live/replay path. Conversely, the two narrow failures do not invalidate all existing positive evidence or establish a core pause-gate defect.

## Required corrections, acceptance limits and disposition

Only R1 and R2 are required code corrections in this initial report. They are behavioral outcomes of the current public contract, not stylistic preferences, hidden interfaces, future Chapter 8 workspace/tab features, disk-restore Chapter 10 requirements or a new coverage floor. Preserve existing signed/provider identities, stable display identity, 270-event replay, observer/history ownership, partial/live handoff, hint wake semantics, truthful undispatched pairing and Agent-wide pause while correcting them.

The code's overall separation is sound: the executable remains a consumer of the public core; the GUI is an optional importing module; durable history remains core-owned; provider-local and display identifiers are separate; the core gate and ordinary job continuation retain their established boundaries. No mandatory architecture rewrite is proposed.

No claim of whole-candidate correctness, exhaustive race freedom, native speech quality, complete DOM security or every event interleaving is made. The new page checks do not render a real browser. The refused-call test exercises actual public core behavior and hub projection, but does not open a browser or assert a native UX result. Root alone decides the correction packet and acceptance after reviewing this report and evidence.

## Failures and provenance discipline

`incidents.md` records all setup/inspection/packaging failures. These include guessed absent evidence paths, wrong manifest key/path assumptions, the incomplete inert DOM fixture, the initially wrong separate-module names (stopped before execution), and preservation-script assumptions corrected using measured facts: a manually mistyped tree, an annotated tag compared without peeling, and generated scratch logs mistaken for extra authored source.

None counts as a behavioral detection. The final preservation check passes. Failed command receipts and the first page harness revision are retained. No failed attempt was discarded, no report score was repaired by hand, and the root's page-hide qualification is explicit.

## Preservation and frozen handoff

Final hash-only verification establishes:

- all 117 frozen files unchanged in student and private copy; student clean at the exact commit/tree;
- immutable Chapter 6 tag still peels to the accepted predecessor;
- all thirteen protected runtime source paths unchanged, including `agent/gui/README.md`;
- all twelve publication paths match the supplied after-body11 baseline;
- all review-input pins unchanged;
- exact prior audit bundles unchanged: entry 76 files, adapter 71, diagnostic 47, fixture review 70;
- all 33 files declared in the native review index and all 31 discrepancy-package files verify;
- only the reviewer test is added authored source in scratch; generated raw logs/reports remain local and are not portable payloads.

Portable outputs are this report and `ch07-code-review-evidence/`. `INDEX.md` declares the exact set, and `SHA256SUMS` covers the report and every other evidence file except the manifest itself. Evidence includes actual commands/cwd/clean environment/timings/exits, safe stdout/stderr, test source, source/hash inventories, read ledger, preservation receipt, carried-evidence pins and root's classification. No binary, archive, raw provider/job payload, key or native recording is included.

**At submission: all commands and tests are complete; no children or pending jobs remain. Verdict: CHANGES_REQUIRED for R1 and R2 only. Correction rounds remain 0/3 until root acts.**
