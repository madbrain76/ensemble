# Chapter 7 student harness reconciliation — unchanged frozen r4

## Result and strict scope

The **one authorized repaired Chapter 7 adapter run passed: 100/100, `passed:true`, process exit 0**. The original Chapter 7 results remain failed historical results; none were replaced, renamed, reinterpreted as passing, or edited. This was harness reconciliation, **not a student source change, new candidate, formal code correction round, or whole-code acceptance**.

Root's supplied follow-up reports independent acceptance of a one-field fixture repair: the prior pause fixture used a legacy argument name inconsistent with the public Chapter 4 contract, while timing, assertions and scores remain unchanged. The student used that expressly authorized repaired executable, without reading private review/case bodies or changing the grader. Root's explanation is attributed assistance, not a new student investigation.

Authority: `/tmp/cadence-ch07-post-fixture-student-followup.txt`, 4575 bytes, SHA256 `e7258e2fe819c6294ba77cc73435b03cdd8f5f46b79705f871297815f090d1a0`, read fully. Grader `/tmp/cadence-ch07-grader-pause-fixed`, SHA256 `c7fa9373e019a9656f485b852ad9be13a3dc31d8139766b7d853b802194e295e`, verified before launch and after the run.

## Exact unchanged identity

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`
- Commit: `48e1d12d9c396b5988c0ce6eccdef67e581872ff`
- Tree: `0f5f7fbe4273e91e1be3aa13d0715a99f671eaf9`
- Worktree: clean before and after; all 117 source hashes verified.
- Accepted predecessor and immutable tag's peeled commit: `f595c787675e2a6343053ec42ea6c7e1e02a930e`
- Core CLI: `/Users/bill/projects/ensemble-edition-3-student/.ch07-local/frozen-r4/ensemble-cli`
- Core CLI SHA256: `4beda0aa805195229c5987a025a2b794b4fe3e77a6526ff194927282fddd85cb`
- Preserved standalone GUI SHA256: `39ebcd3e42a57d8972b6278fb8cc66021e38afd41afa2cd230b93499e6b13a81`

The adapter receives the existing CLI plus unchanged optional module roots and builds modules within its own public adapter workflow. No standalone student candidate build was performed. As before, its temporary internally built GUI binary hash is not retained by the student; the frozen module source is pinned. The standalone GUI hash is not falsely substituted for that temporary binary.

## Actual repaired invocation and outcome

CWD: `/Users/bill/projects/ensemble`

```sh
/tmp/cadence-ch07-grader-pause-fixed -ch 7 -json -core /Users/bill/projects/ensemble-edition-3-student -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer -consumer-case /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json -workflow /Users/bill/projects/ensemble-edition-3-student/examples/ch06 -stream-exercise /Users/bill/projects/ensemble-edition-3-student/examples/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble-edition-3-student/gui -gui-model claude-sonnet-4-6 -media-models /Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-evidence/declared-media-models.json /Users/bill/projects/ensemble-edition-3-student/.ch07-local/frozen-r4/ensemble-cli
```

The exact closed process environment is retained in `repaired-chapter7-invocation.json`; no inherited LLM/provider settings or secrets were used. The existing media map and inherited consumer/workflow/streaming/GUI declarations were reused without modification. Start `2026-10-17T21:53:39.138812+00:00`, end `2026-10-17T21:53:53.816881+00:00`, exit `0`.

| Public check | Outcome |
|---|---|
| WebSocket observations | Pass, 25/25; 71 messages and all six required types |
| Replay | Pass, 20/20; 11 messages; report retains its active-partial/handoff/window proof qualification |
| GUI traffic log | Pass, 15/15; 85 lines, both directions/timestamps/JSON |
| Pause gate | Pass, 20/20; all three tools dispatched/finished and turn ended after unpause |
| Corrected Chapter 6 parity | Pass, 20/20 |

Original machine report title is retained verbatim (including its inherited Edition 2 title); no report fields were rewritten. Report SHA256 `053997da4084ca3044b2d55a9e0057be5fa0d02c07e4293da00d81d953630428`; stderr SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

## Historical package preservation and qualified carryover

Before and after execution, all **252 original portable UTF-8 file hashes**, **22 indexed local artifact hashes**, and **117 source hashes** were verified against the existing sealed manifests. The original manifest files themselves also retain their exact hashes:

- Source manifest SHA256 `8e4bb14356fb45dee8a7ebc914f439d1c38a46fd7772184ef1506546a9681b6a`
- Evidence manifest SHA256 `f8bf301d5d041ae94cc8ccca5aed7c25c8e2a245c1121648860a5b857e10ee48`

Original r4 Chapter 7 report SHA256 `103f5543a667643f6ca93edd24a5857fe1567b48d5b47f81682989bc145b9f17` remains **80/100, passed:false**; earlier failures and intermediate candidates remain intact. New evidence lives only in the separately authorized report and reconciliation-evidence directory, with its own index. Historical and new provenance are intentionally not merged by editing the old seal.

The actual existing same-byte r4 Chapter 1 report (SHA256 `f5f3b049294ce0f224cdf86853eee291fdcf937ce67daf30deedc7fa4b73019d`) remains 100/100 with both genuine 62-second zero-point lifetime phases passing. Its exact CLI hash matches this run. **No new Chapter 1 run occurred.** The original final matrix and r1-only browser observations likewise remain historical, with their original limitations; **no new full matrix, standalone browser, native, headed, audible, microphone, paid, key-access, production-edit or candidate-build campaign occurred**. No children were spawned. No acceptance/tag/publication action was taken.

## Native plan/launcher preparation — no execution

`native-exercise-plan.md`, `native-policy.json`, `native-expected-request.json`, `native-owned-probe-proposed.mjs` and `native-proposal-receipt.json` prepare the finite future scenario and exact command. Launcher SHA256 `2d582d6fb4caa5704e07a44a6d98c2799cd8e7b2bdb9f03e8880f84decfddcd6`; policy SHA256 `2bc352bd307417d7392125b1b9c5f993ac81443ba2e705bae0f8e445981da9dc`. Only syntax checking was performed, including after small preparation-only hardening. The launcher was never imported/executed and no native/browser preflight was run.

**Known blocker:** frozen production does not select or expose the resolved utterance voice. A local-default advertisement cannot prove actual service selection. The new launcher has an intentional pre-launch hard stop with no flag override. Root must resolve this without changing production speech behavior, replacing native callbacks, or silently permitting remote service before any new pinned launcher can execute. Runtime ownership/request/cleanup controls are proposed, not observed. The exact served-asset ordering distinction and production visibility-change empty-utterance caveat are explicit.

The prepared finite scenario uses **“Harbor lights shine softly.”** and **“Quiet gardens rest.”**: queue both, observe first native start and pause, genuinely cancel the active/queued work via the real page control, observe native idle/pause release, requeue the second phrase and require its native start **and end**, then pause release. Errors/unavailable conditions stop the attempt; no retry/fallback. Budget 55 seconds plus 10 seconds owned cleanup. No microphone/recognition/focus manipulation. Headed/audible output needs separately explicit later permission and human listening remains separately consented evidence.

All final-r4 application/browser usability, native voice locality/callback recovery, subjective hearing, recognition, paid/live and acceptance gates remain unfulfilled absent their separately required observations/permissions. The repaired automated pass alone does not waive them. Earlier incomplete speech evidence and the earlier unmuted-probe authorization gap are preserved.

## Preparation failures and limits

During read-only preparation, two guessed engine paths did not exist; the tracked request builder was then read directly and its `stream,omitempty` behavior verified. These were inspection errors, not source changes or execution attempts. A policy-shape check was also prematurely described as an iteration bug; direct reading confirmed the existing `Object.entries` loop was already correct, so that diagnosis was retracted and no iteration code changed. The initial sealing check then hit a regex false positive on the filename `native-policy.json`; this was corrected and retained without executing the launcher. Proposed provider/ownership/voice/browser controls remain untested except for JavaScript parsing. The current known voice blocker means this is **not an execution-ready approved harness**.

## Files and completion

`ch07-student-harness-reconciliation-evidence/index.json` seals every newly produced safe text artifact and this report, excluding the index itself to avoid recursive self-hashing. `before-identity.json`, `after-grade-identity.json` and final preservation receipt qualify the historical package independently. All finite authorized execution/preparation is complete; no execution jobs or children remain. Parent owns any later review, consent, permission, execution decision and acceptance.
