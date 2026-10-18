# Chapter 8 — fixture-only formal R1 re-review

**Verdict: PASS — R1 resolved; the initial whole-code review’s sole blocker is cleared.**

This closes the narrow fixture-compatibility finding against the exact frozen correction. It does not claim browser/native/live acceptance, create an acceptance tag, alter the original review, or turn carried grades into fresh runs. No additional mandatory correction was found in this one-file review.

## Identity and reviewed scope

- Task SHA-256: `c0f94a707152430d1a6d450df4c86f2d5e589bfba5940c70e4c8029f60f95584`.
- Authoritative input `ch08-code-review-r1-input.json`: `a2566bede526ee5836c97a6e987eb1200689dca1c9ec9b0c9dade8d677747be7`.
- Candidate: `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`; tree `fa896caa404b87b08de61eb38ce079eb38138fbe`; parent `f97a6850468a9cdc6875507f2821ccf6c397654d`.
- The only changed tracked path is `gui/test/correction-r1-page.mjs`, SHA-256 `8adfb4dbd88e883d50be7556ce4b72f137771c1ef140a36c4332ef0004762ae0`. All other 125 source members match the initial reviewed parent. Production, Go tests, module files, grader inputs and APIs are unchanged.
- Required `AGENTS.md`, coding/reviewer skills, and execution instructions were reloaded fully. Other prerequisite reads carry by their exact initial-review hashes, retaining prior full-versus-section qualifications; see `prerequisite-qualification.json`.
- Read the complete student R1 report, changed fixture/diff, source/evidence manifests, wire provenance, relevant production app/workbench/SpeechQueue/tool projection and the unchanged Go producer. This is a fixture-only re-review, not another entry audit or whole-production replay.

## R1 closure rationale

The initial failure and its original seal remain unchanged. The submitted correction now repairs the actual boundary rather than skipping its assertions:

| Location in `gui/test/correction-r1-page.mjs` | Review evidence |
|---|---|
| 10–22, 24–44 | Minimal event/style/media/scroll boundaries cover current production imports. IDs/tags come from actual `index.html`; missing/duplicate IDs, module escape, server-message rejection and unexpected reconnect timers fail explicitly. The actual `app.js` and its modules execute; no alternate production implementation is substituted. |
| 45–55 | A full accepted initial settings snapshot enables the real control. Cancellation dispatches `change`, checks the sparse `update_settings` message, verifies draft rollback/no premature cancellation, and then verifies cancellation on the accepted snapshot. This replaces the obsolete direct `.onchange` assumption. |
| 47, 65–75 | Artifact lookup spans Transcript and Actions but asserts exactly one matching display identity. Tool location is explicitly asserted to be Actions. Visible call/result/error, arguments, opaque exclusion, replay silence/manual readout, pause/unpause and stale-callback protection remain active. |
| 78–86 | Empty/opaque artifacts remain silent. The optional real-wire branch checks current Actions placement, refused batches, on-demand errors, replay/reconnect silence and six retained result artifacts; it is not silently bypassed. |

Fresh standalone and real-wire runs both pass. The workbench and queue fixtures also pass. The correction maintains the original behavioral assertions and strengthens current entrypoint/identity checks; it does not weaken a failing production behavior into a green fixture.

## Independent fresh controls

Node `v24.20.0`; exact executable hash, scrubbed environment, cwd/HOME, argv, output hashes, deadlines, timing and exit statuses are in the indexed receipts. Commands are sequential, with 30-second deadlines and owned process groups. No credentials are inherited.

| Control | Exit | Time | Actual outcome |
|---|---:|---:|---|
| `current-page` | 0 | 0.132s | PASS |
| `current-page-real-wire` | 0 | 0.075s | PASS |
| `current-workbench` | 0 | 0.040s | PASS |
| `current-queues` | 0 | 0.040s | PASS |
| `guard-old-parse` | 0 | 0.040s | PASS |
| `guard-old-behavior` | 0 | 0.040s | PASS |
| `guard-negative-parse` | 0 | 0.040s | PASS |
| `guard-negative-behavior` | 1 | 0.040s | Expected behavioral failure |
| `guard-restored-parse` | 0 | 0.040s | PASS |
| `guard-restored-behavior` | 0 | 0.040s | PASS |

### Independent executable negative and restoration

The owned-copy mutation removes only the production `SpeechQueue` completion guard that rejects a canceled/stale utterance. The submitted fixture is unchanged in all stages. This is independent of the student’s carried tool-result-projection mutant.

- Old: parse **0**, behavior **0**.
- Negative: parse **0**, behavior **1**, specifically `AssertionError: 'unpause' !== 'pause'` at the stale-callback assertion (fixture line 75). This is executable wrong behavior, not a syntax/import failure.
- Restored: parse **0**, behavior **0**; original/restored `speech.js` SHA-256 `818693251c3bea40bea32e0c10cc8c8dc31a908055c2320dc564546d14cc5ad6`; negative `935b2e5adf6eb8abec8d5b93e08817ac1d18640c98ecb23564cd812b08797344`.

All copied production files restore to their original hashes. `mutation-results.json`, `speech-guard-mutation.diff` and old/negative/restored receipts preserve the distinction.

## Real-wire provenance and limits

The fresh optional branch consumes the pinned local-only `.ch08-local/r1-wire.json` (78,386 bytes; SHA-256 `99a7078ba9a64843b42d40823d1ae0df86fb5b9b15855ed1f5937380d66f635c`). Its producer is unchanged `gui/correction_r1_test.go` (SHA-256 `36f9ad1a5a85e36e13c1a23e64462d0b6e918a0bfb0d52377091a192c57964d8`), whose carried receipt records the real local Agent/hub/WebSocket path. The producer was not rerun here.

The fresh consumer exercises the named refusal/replay/WebSocket streams through the actual production page handlers. Current settings in the retained replay/WebSocket streams have TTS enabled, so silence assertions are not explained away by a disabled preference. This remains a controlled DOM/socket/synthesis boundary with captured real local wire: **not browser rendering, native speech, permission, microphone, audio/listening or paid-provider evidence**. Raw wire is not copied into the portable bundle.

## Carried evidence; not rerun

`carried-controls.json` records exact hashes for the initial five-module race/vet receipts and both full grade receipts. All six runtime identities and nine retained historical binaries verify unchanged; their original commit/VCS metadata remains original, not relabeled as newly built R1 artifacts. The prior whole-code review’s 132 sealed members are intact.

- Initial root/GUI/consumer/workflow/stream race/vet remain **carried passes**.
- Initial Chapter 8 **100/100** and the single Chapter 1 lifetime **100/100** remain **carried passes** on the originally pinned same CLI, supported by the unchanged production inputs. Fresh R1 grade/lifetime runs: **zero**.
- Initial real-process budget/settings and fractional-rate evidence remains qualified as before. The original failed page fixture remains a historical failure; the new repaired fixture’s pass does not rewrite it.
- Student producer and result-projection sensitivity receipts are carried student evidence, not fresh reviewer runs.
- No browser/native, paid/live provider, secret/key/config, student/child agent, broad Go suite, entry matrix, commit/tag, source/grader/guide/old-report edit, or shared cleanup was performed. No new API or framework recommendation is warranted.

## Failed preparations and preserved qualifications

Read/setup mistakes are preserved, not presented as candidate failures: the initial diff entered a pager because `--no-pager` was omitted; only its uniquely verified owned pager was terminated, the read job ended, and the full diff was recovered with paging disabled. A guessed adjacent filename and a manifest-key assumption caused read-only lookup errors. Two original preparation scripts stopped before controls on a guessed raw-wire path and an incorrect envelope-shape assumption. The corrected preparation derives the path/hash from provenance and validates the actual named streams. `attempt-notes.json`, `preparation-failure.json`, `preparation-envelope-failure.json` and all three scripts retain the sequence. There was no unexpected runtime-control failure or repeat-to-green run.

Both source TARs remain local-only/excluded and unchanged; no TAR, binary payload, raw private wire, bytecode or symlink is exported in this review bundle.

## Preservation, inventory and next boundary

Before/after verification: **20 inputs, 126 source members, 31 protected files, 44 student R1 indexed entries, 132 prior review members, six runtime identities and nine retained binaries** match. Candidate and tags remain unchanged/clean, the one-file diff is unchanged, and the owned copy is restored. All **10** command receipts are terminal, and the earlier pager/read group is also ended.

`ch08-code-review-r1-evidence/inventory.json` is the self-excluded portable hash/size index. `local-artifacts.json` separately records retained owned scratch metadata; no local payload is implicitly called portable. `final-verification.json` and `final-seal.json` are the preservation/process/scope receipts. Only the new R1 report/evidence paths and owned scratch were written.

**Next boundary:** coordinator may record R1/code-review closure using this exact candidate and seal. Subsequent rendered/native/live acceptance remains separate; this review neither performs it nor adds routine permission or architecture prerequisites.
