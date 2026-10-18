# Edition 3 Chapter 7 — formal round-1 corrected-code re-review

## Verdict and authority

**PASS for this frozen corrected-candidate whole/adjacent code re-review. R1 and R2 are closed. No new must-fix existing-contract code defect was found in the bounded checks.**

This is the same independent reviewer continuing after the initial review, not a fresh entry, adapter, or pause-fixture campaign. Exactly the two assigned required findings were assessed. **C1 remains non-blocking, unchanged, and outside the correction requirement.** It is not a hidden acceptance condition.

The student remains submitted. Formal correction usage remains **1/3**, not another round. Root alone decides acceptance/tagging after inspecting and joining this report. This PASS does **not** turn offline models, an inert DOM/speech boundary, or carried r4 listening into a new native/current-bundle or real-provider exercise claim.

## Frozen identity

- Candidate: `463996b765631d31923a5726036a34942a67cefd`.
- Tree: `5e95dcdd716f7c25c96cc3ebcde13d000bb38670`; clean student worktree.
- Parent/original negative control: `48e1d12d9c396b5988c0ce6eccdef67e581872ff`.
- Accepted Chapter 6 tag still peels to `f595c787675e2a6343053ec42ea6c7e1e02a930e`.
- Main HEAD remains `23c75ede7944d9a912ae74d3f683354d163f0c86`.
- **120 frozen source/support members**, not “120 production files”; **138 submitted portable files** verified. The seven changed paths remain within the two corrections.
- Chapter 7 contract SHA-256: `68d02d8e1170c2912676b8f162847f8fac25c87216514ed06adbd5f0ea011b47`.
- Corrected source archive SHA-256: `f2ba25cedc080abd295f35af9aeae3128dab0793d9fb1c6c71abef5c2750593a`.
- Frozen corrected CLI SHA-256: `000a2997e5b8396feb2e0892f3b8160f48cce8adc34b44eeb940e63dd3873f11`.
- Frozen corrected GUI SHA-256: `c41d8115a6b22b8fb294ab11f4303926e5d485cac7b40025c39fa0e42c0eccd4`.
- Accepted grader executable SHA-256: `c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e`.
- Accepted pause-driver source SHA-256: `bcdbcd50397d40de60c9fa871cd9d588403efeee34162aa01830cf6c99913b09`.

All 31 review-input pins match. Paths, byte lengths, the remaining importing binaries, old r4 artifacts, command inputs, and local-only control binaries are inventoried in the evidence. The grader builds its GUI from the verified private source module; a test binary or source rebuild is **not** represented as byte-identical to the student's frozen GUI artifact.

## Read and review scope

The task explicitly permits retained complete prerequisite reads when unchanged. All twelve governing/prerequisite/chapter hashes matched the prior complete reads, including `AGENTS.md`, both skills, current Edition 3 execution/readiness/student/carryover documents, historical workflow/carryover, voice/writing procedure, and full Chapter 7. They are recorded as **reused full reads**, not newly reread files. No new execution-procedure entries were present at that comparison.

Fresh inspection covered the complete correction diff and new helper/tests, the affected admission/observer/history projection and page paths, the student's round-1 report and manifests, root's required-finding decision, native applicability and exact grader records. Every declared source and portable member was hash-verified. Relevant public API definitions were inspected while building the independent test; incorrect reviewer guesses about those APIs are retained as setup failures.

Adjacent review covered ordinary/custom/process/control admissions, real errors versus refusals, paired history, replay identities, old-history interpretation, readable/opaque projection, compact automatic narration, final/result ordering, cancellation/error symmetry, the library-root/importing-module boundary, and the actual unchanged pause/control tail. It did not expand into later-chapter tabs/trees/workspaces or add a disk-restore requirement.

## R1 — truthful tool admission: closed

### Source assessment

`internal/common/common.go` adds an optional durable `ToolEvent.Dispatched` fact. `internal/agent/dispatch.go` sets it on the actual admission path, including admitted no-job controls. Refusal/pairing paths do not acquire that fact merely because a `tool_called` record is retained.

`gui/hub.go` now projects a dispatch only with affirmative admission evidence: the explicit fact, job-backed legacy evidence, or an actual observed dispatch. An ambiguous job-less historical `tool_called` remains a readable call artifact and retains its result; it is not promoted into a start. This is narrower than filtering every nil-job call, which would incorrectly hide current admitted controls.

### Independent real boundary

The private test compiles the unchanged production GUI/core together with one reviewer test file. It uses the actual public Agent/jobs APIs, a loopback fake Messages server, the production WebSocket upgrade/handler and real masked client frames. It does not replace the controller with a fake page model or add a candidate diagnostic API.

The same valid six calls are first returned without required explanation and then in an explained batch, reusing provider call IDs:

- ordinary file read and file write;
- a real process that writes a marker;
- a registered custom tool that really executes and errors;
- admitted `tool_limits` without a job;
- admitted `wait_for_job` for an unknown handle, producing a no-job control error.

A barrier at the second local model request permits inspection before retry. Corrected results:

- Refused batch: **0 actual core admissions, 0 custom invocations, no file/process effect, 0 dispatch frames live/replay/fresh reconnect**, and six truthful error results.
- Explained batch: **6 actual admissions**, one custom invocation, both file and process effects, and six projected admissions. The custom/job error and no-job control error remain genuine admitted errors, not suppressed calls.
- Across both batches: **12 called/returned pairs**, eight expected error results (six refusals plus two admitted errors), and twelve distinct call/result display identities despite reused provider IDs.
- A fresh completed-history socket and a serialized/rehydrated history each replay six actual admissions, not the six refusals. Stable per-call result identities agree across those views.
- Removing only the new admission bit from a private history copy conservatively retains four unambiguous job-backed admissions and all paired results; unknown old no-job controls are not invented as starts. This is schema/projection compatibility evidence, not Chapter 10 persistence acceptance.

### Sensitivity and restoration

The exact same reviewer test compiles against original r4. It fails behaviorally: six refused dispatches appear in each relevant refused view; completed/rehydrated/legacy projections inflate to twelve. The real core still admits only the six explained calls. This is **not** a compiler, server-start, timeout, or disk failure.

The corrected binary/source is then selected again and passes the unchanged assertions. All original 117 r4 members and all corrected 120 members were verified unchanged; there was no mutation of either real worktree. “Restoration” means returning to the hash-verified corrected path, not claiming a nonexistent production rollback.

## R2 — complete readable tool-artifact speech: closed

`gui/web/gui/tool-readable.js` selects only the readable call context/arguments and accumulated result/error. `app.js` keeps those values together and refreshes the speaker text on result arrival and later call-final updates. The short automatic dispatch summary remains separate. A completed result does not trigger an unsolicited full-result narration.

The independent page check executes the actual unchanged page modules in a Node VM with inert DOM/WebSocket/speech boundaries. It feeds the real protocol captures above and selected ordering/content cases; it is not a replacement page implementation and is not native-browser evidence.

All nine final corrected cases pass:

1. Long readable result/arguments survive result arrival and a later call-final update. Automatic summary remains 22 characters; on-demand speech includes the complete call/result (3,990 characters in the fixed fixture).
2. The corresponding error case includes complete error content and its error label (3,995 characters), not only the compact summary.
3. Cancellation and stale callbacks do not release a newer utterance's pause; actual current completion does.
4. The analogous `onerror` path has the same generation/gate behavior.
5. Empty/signed/opaque content is not spoken.
6. Actual refused-wire frames produce no automatic execution summaries, and all six error artifacts speak their complete readable result/context.
7–9. Fresh reconnect, rehydrated history and the conservative legacy view each preserve complete readable speech for all twelve result artifacts.

The unchanged final check also runs on original r4. Six relevant cases fail for the intended R1/R2 reasons: full result/error speech is still the 22-character summary, result speakers omit their payloads, and refused calls generate execution summaries. The three unchanged opaque/queue-generation guards still pass. Returning to corrected source and fresh corrected wire passes all nine cases again.

The final result check compares the **complete wire result** plus the appropriate error/result label. It does not demand identical punctuation to the DOM's concatenated caption. A preliminary overly strict caption assertion and the inert DOM's initially incomplete `textContent` behavior are retained as reviewer errors, not candidate failures.

## Fresh final gates

All gates ran sequentially in owned scratch with a clean, explicitly recorded environment, local fake transports, no credentials/paid requests, and no outer deadline that could cancel the genuine long waits. The full root race suite was allowed to finish. No new browser, microphone or native speech run occurred.

| Gate | Exit | Wall time | Recorded outcome |
|---|---:|---:|---|
| `root-race` | 0 | 140.489 s | PASS (test/vet exit; no grader score) |
| `root-vet` | 0 | 0.909 s | PASS (test/vet exit; no grader score) |
| `gui-race` | 0 | 1.912 s | PASS (test/vet exit; no grader score) |
| `gui-vet` | 0 | 0.168 s | PASS (test/vet exit; no grader score) |
| `consumer-race` | 0 | 0.128 s | PASS (test/vet exit; no grader score) |
| `consumer-vet` | 0 | 0.107 s | PASS (test/vet exit; no grader score) |
| `workflow-race` | 0 | 1.381 s | PASS (test/vet exit; no grader score) |
| `workflow-vet` | 0 | 0.112 s | PASS (test/vet exit; no grader score) |
| `stream-race` | 0 | 0.122 s | PASS (test/vet exit; no grader score) |
| `stream-vet` | 0 | 0.105 s | PASS (test/vet exit; no grader score) |
| `chapter7-fresh` | 0 | 14.573 s | 100/100; `passed=true` |
| `chapter1-fresh` | 0 | 124.066 s | 100/100; `passed=true` |

The root remains a library with nested `cmd/ch01`; the separate importing GUI and `examples/ch04-consumer`, `examples/ch06`, and `examples/ch07` modules were not flattened. Race/vet covered the current submitted tests **before** adding the reviewer-only test file to the private copies.

### Exact final grader routes

Both use the same exact corrected frozen CLI above, not the old r4 CLI. Cwd is `/Users/bill/projects/ensemble`.

```sh
/tmp/cadence-ch07-grader-pause-fixed -ch 7 -json -core /tmp/cadence-ch07-code-r1-rereview-independent/source -consumer /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch04-consumer -consumer-case /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch04-consumer/consumer-case.json -workflow /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch06 -stream-exercise /tmp/cadence-ch07-code-r1-rereview-independent/source/examples/ch07 -stream-model claude-sonnet-4-6 -gui /tmp/cadence-ch07-code-r1-rereview-independent/source/gui -gui-model claude-sonnet-4-6 -media-models /Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-evidence/declared-media-models.json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli

/tmp/cadence-ch07-grader-pause-fixed -ch 1 -json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/correction-r1/frozen/ensemble-cli
```

Chapter 7 freshly reports **100/100 and `passed=true`**. This includes the already accepted corrected pause fixture; no old entry/adapter/pause-sensitivity campaign was repeated.

Chapter 1 freshly reports **100/100 and `passed=true`** in **124.066300 seconds**. Its `generation_headers` and `generation_body` checks each report receipt of the complete answer/usage after a **1m2s** delay. No separate per-phase sub-second elapsed measurements are emitted by that JSON report, and none are invented here. Both genuine 62-second gates completed. This is a fresh same-CLI result, not a carried previous-commit score.

## C1, native applicability, and real-model boundary

C1 stays exactly as root classified it: a measured non-blocking contract-clarity note, not a third required finding. The control tail and its hash `f760a9928e290234cbb87fe6438c862787fa322cf525a7457491c3bfe048d378` are unchanged. No pagehide, visibility, focus, microphone, activation or launcher correction was requested or made.

The earlier r4 native exercise remains a real historical positive. The human receipt records Bill's exact statement, `I heard "Quiet Gerdens Rest"`, with its existing qualification. That is not erased, but it is also not a fresh hearing of this corrected bundle or of a corrected full tool-result utterance.

`speech.js`, `queue.js`, `artifact-scroll.js`, `format.js`, `index.html`, and `style.css` are byte-identical to r4. Their prior mechanism evidence has that narrow applicability. `app.js` changed and `tool-readable.js` is new. Their current projection/order behavior has real-wire plus actual-module controlled evidence, **not new native speech-service, permissions, long-utterance audio/performance, or whole-bundle browser evidence**.

All new model traffic in this review was local fixture traffic. No current-candidate real-provider streamed-GUI demonstration is established here. Historical book links are not new student proof, and a 100/100 offline grade is not such a demonstration. If a real-model operational demonstration remains outstanding under the existing exercise, it remains a root-level exercise-evidence item; this review neither invents permission to run it nor treats it as silently completed. No new code blocker was inferred merely from an unauthorized run not being performed. Root retains final exercise/acceptance authority.

## Failures, assistance, and what was not counted

Every completed command has its actual argv, cwd, clean environment, wall time, exit, and stdout/stderr. Grader score and `passed` are recorded separately; compiler/unit/vet commands do not receive invented grader scores.

Retained reviewer mistakes include: source-manifest shape and recorder-argv assumptions; querying an evidence-read handle instead of the long-gate handle (no gate was cancelled or duplicated); private-test package/public-API/history-method/system-prompt mistakes; the reviewer WebSocket extended-length parser bug; a wrong page asset directory; inert DOM text aggregation; the overly strict DOM-caption expectation; summary-field/stdout-suffix/publication-baseline filename assumptions. Their original stderr/source variants and packaging ledger distinguish them from the two intended original-r4 behavior failures.

The final shared assertions pass corrected, fail original r4 for the intended reasons after successful compilation, and pass corrected again. There was no grader change, score accommodation, hidden candidate hook, supplied student patch, or source change to manufacture those outcomes. No children or outside implementation assistance were used.

## Preservation and handoff

`preservation.json` verifies:

- all 31 current pins; clean commit/tree/parent and immutable predecessor tag;
- all 120 candidate members in the worktree and private source, and 138 declared corrected portable files;
- thirteen unrelated runtime source paths and twelve publication paths against the sealed protected baseline;
- five exact prior audit bundles, their manifests and covered reports;
- 252 old student-evidence members and 33 indexed native members;
- unchanged native mechanism/control bytes and the main HEAD.

The private source has one reviewer test plus test-generated `api.log`/`cr` artifacts. Those extras are inventoried by hash only, not confused with student changes or copied as raw payloads. No shared or manual cleanup was performed. Last observed free disk was approximately 1.5 GiB; no resource failure was credited as a detection.

Portable evidence contains only safe UTF-8 receipts, reviewer sources/variants, summaries and hash inventories. Raw wire/provider/tool payloads, binaries, archives and generated private traces remain local. All owned test/command jobs are finished. No experiment, child, or question is pending. No source/manuscript/grader/guidance/reference/runtime edit, install, service operation, commit, tag, stage or push was performed.

Report: `docs/edition-3-notes/ch07-code-review-r1.md`.
Evidence: `docs/edition-3-notes/ch07-code-review-r1-evidence/` (`INDEX.md` and `SHA256SUMS` define the exact sealed set and cover this report; the checksum manifest excludes itself).

**Final disposition: PASS for the assigned corrected-code review; R1 closed, R2 closed, C1 non-blocking. Root acceptance/tagging and any remaining operational exercise decision are still pending root's independent inspection.**
