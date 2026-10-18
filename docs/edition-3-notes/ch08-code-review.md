# Chapter 8 — INITIAL whole-code review

**Verdict: NEEDS CORRECTION — one blocking test/integration-maintenance finding (R1).**

The reviewed production implementation passed the fresh module race/vet checks, Chapter 8 publication grade, once-only Chapter 1 lifetime grade, and independent real-process/page sensitivities. I found no production-code blocker in the exercised settings, budget, speech-rate or ownership seams. However, a checked-in inherited production-page integration fixture no longer runs against this candidate: it crashes during app initialization, before its assertions. The owned diagnostic proves this is repairable test scaffolding/locator drift, **not evidence of a browser defect**. Do not call the submitted whole candidate code-accepted until the checked-in fixture is corrected and the affected checks pass.

This is an initial review, not an acceptance tag or correction implementation. Student remains submitted and unchanged; formal correction rounds used remain 0. No student agent, child agent, browser/native UI, paid provider, key/config, public grader, guide, manuscript or old-report edit was performed. No Chapter 7/entry/adapter campaign was run.

## 1. Frozen identity and prerequisites

- Task SHA-256: `79343abe8bfea0b8c2efa5a1a4c149d5b1f5e3a531c74daf7150f2e76f2d980b`.
- Authoritative input: `ch08-code-review-input.json`, SHA-256 `267c288144090f392df6d1eed0223cc32cfa451386c5be40636124d191804306`.
- Candidate: `f97a6850468a9cdc6875507f2821ccf6c397654d`, tree `b3c4f005d5f43a7f20fc6227c65e1f763b93fac9`; clean before/after. Review covered the whole production candidate, not only its 19 changed paths.
- Predecessor: `463996b765631d31923a5726036a34942a67cefd`, tree `5e95dcdd716f7c25c96cc3ebcde13d000bb38670`; tag `edition-3-ch07-r1` / object `b66ba6520e0569c65259905a93057f3d6d04e8de` remains unchanged.
- Current public execution/student/carryover/readiness/Chapter 8 instructions and the complete student report/manifests were read. The relevant inherited Chapters 1–6 public contract sections were read; Chapter 7, governing workflow/carryover/craft and reviewer-role readings carry only with verified unchanged pins. See `prerequisite-carry-forward.json`, `source-review-coverage.json`, and `input-verification-before.json` for the exact boundary.
- All 27 pinned inputs, 126 tracked source records, 31 protected existing files, 126 student indexed entries, and 630 retained local artifacts verified. The source archive has 126 files plus 26 directories; its bytes/members match the source, without extraction or repacking.

### Runtime artifact identities

The five current submitted executables and frozen grader were hash/build-info verified. The root CLI below is the **same executable** used for both fresh publication gates. The Chapter 8 command uses explicit core/consumer/workflow/stream/GUI paths, declared `claude-sonnet-4-6` models and the pinned `media-models.json`; exact argv is in `chapter8.receipt.json`. Existing earlier/unqualified binaries were not silently substituted.

| Artifact | SHA-256 | recorded VCS revision / modified |
|---|---|---|
| `consumer-ch08` | `9bbcc2a86713aac4daefc81e034e70ec9c7e518ba8c424bacfedf13a17a59612` | f97a6850468a9cdc6875507f2821ccf6c397654d / true |
| `gui-ch08` | `6b513f3b10c5906bf66c6acf89e645903c95b0a397b711d93ae067b5b35ec0ae` | f97a6850468a9cdc6875507f2821ccf6c397654d / false |
| `root-ch08` | `160946a79ca8ebdcf17acfaf315b9710be87cd7294387eb36ac95d0062491945` | f97a6850468a9cdc6875507f2821ccf6c397654d / false |
| `stream-ch08` | `fed84b61e5aa033f857061824263f2a012aba85321276e65f8206c392f07c934` | f97a6850468a9cdc6875507f2821ccf6c397654d / true |
| `workflow-ch08` | `f04f9aa00066277168011deaf1c6d5c08706a7ad57a86c18bb225fcce5a5b3b9` | f97a6850468a9cdc6875507f2821ccf6c397654d / true |
| `cadence-ch08-grader` | `529187f45b272aeb025670651fa9157ea1f7e1262170ae5593a517afa07a0645` | 2cebe5ae76b47eabfd1f1cb7e8bef95c98b7c20e / true |

Recorded `vcs.modified` metadata is preserved, not rewritten into a pristine-build claim. Source membership, retained hashes, fresh source race/vet builds, and owned-copy production builds provide distinct evidence. `artifact-verification.json` retains complete `go version -m` output for all six identities.

## 2. R1 — restore the checked-in inherited production-page fixture

**Priority: P2; blocking whole-candidate code-review closure.**

**Location:** student `gui/test/correction-r1-page.mjs:22–35` (old element/window shim, artifact locator and cancellation entrypoint), compared with `gui/web/gui/workbench.js:3–7` and current `gui/web/gui/app.js` workbench initialization, Actions routing and accepted-settings handling.

**Reproduction:** the exact submitted file, with no edits:

```text
node --experimental-vm-modules /Users/bill/projects/ensemble-edition-3-student/gui/test/correction-r1-page.mjs
```

`page-inherited-integration.receipt.json` records exit **1**; stderr records `TypeError: window.matchMedia is not a function` during `setupWorkbench`. The old shim also omits the new workbench controls/root-style boundary, searches only Transcript for artifacts now intentionally routed to Actions, and invokes an obsolete read-aloud `.onchange` hook rather than current accepted-settings handling. These latter incompatibilities are source findings; the unchanged failing command reaches only the first error.

**Impact:** this retained regression entrypoint does not currently execute its pause/steering, authoritative artifact identity, replay/TTS dedupe, manual readout, or cancellation-generation assertions. The student's green focused workbench/queue fixtures and public grader do not make this separate checked-in fixture green. Its first crash is an actual test-regression gap, not grounds to claim native layout/audio failure or to redesign the production API.

**Independent diagnostic:** one owned copy supplies the current inert DOM/style/media boundaries, locates the same artifacts across their current panes, and exercises cancellation via an accepted settings snapshot. Original behavioral assertions remain. Against identical production app bytes it exits **0** (`page-adapted-diagnostic.*`, `page-fixture-diagnostic.diff`, `page-diagnostic-identity.json`). This was not a student fix and does not override the original failure.

**Minimal correction:** update this existing controlled-page fixture for the current Chapter 8 app, retaining its inherited behavioral assertions and correctly covering its current entrypoints/locators (including any optional wire-backed branch). Run that submitted fixture and the affected page/settings checks. No new grading framework, public API, product architecture change, or malformed-input campaign is warranted. Unchanged expensive gate evidence may carry with verified pins; this review does **not** request another full lifetime run.

## 3. Whole-candidate production review

| Boundary | Review/result and limits |
|---|---|
| Public core / common vocabulary / external clients | Reviewed `ensemble.go`, common definitions, CLI and the three independent example modules. Public-client module seams remain intact; no new reference-internal dependency requirement or prescribed API spelling. Fresh race/vet and Chapter 8 inherited parity passed. |
| Engine versus Agent ownership | Reviewed model/render/request/parse/SSE/stream code and Agent admission/dispatch/history. Engine remains the single-request transport/provider boundary; orchestration, history and tool rounds remain Agent-owned. No settings copied into provider routing or model-name heuristic introduced. |
| Budget ownership | `internal/agent/agent.go:203–224` owns the thread-safe getter/update seam; `actor.go:278–282` captures the effective default-or-positive budget at admission, and `actor.go:405` applies that captured value. Fresh source/real-process checks prove current-turn isolation, effective next turns, positive 18 beyond the old small-loop failure mode, and zero reset exercised through 20 fake-provider tools. A bounded test does not prove an infinite budget. |
| Settings persistence / atomic publication | `gui/settings.go:57–67,106–158`: full typed snapshots and validated partial merges; temporary-file save/rename succeeds before Agent policy/memory publication and broadcast. Negative budget and font/TTS ranges follow the candidate's declared clamp/default policy, not an invented reviewer range. Independent process control rejects a wrong-typed sibling patch and forces a real final-save rename failure; neither partially publishes or changes effective policy. |
| Startup / terminal / sockets / replay | Reviewed GUI entrypoint, terminal loop, socket framing, hub observation and replay ownership. Actual process tests use terminal input and WebSocket turns, concurrent subscribers and restart; no local in-memory replacement server was used. |
| Prompt and visible-context isolation | Independent fake-provider request assertions hold the system/base prompt hash constant and require every user text item to be one of the explicit test turns. Settings/preferences never become user-turn text in those real requests. Source keeps the setting in its owner rather than modifying the immutable base prompt. |
| History / jobs / process lifetime | Reviewed history, jobs, tools/process ownership and CLI interruption/EOF paths. Existing Agent-versus-job lifetimes and process cleanup remain inherited boundaries. The required same-CLI Chapter 1 full lifetime gate passed freshly once, not as cached evidence. |
| Theme / font / panes / resizing | Reviewed actual HTML/CSS/app/workbench and artifact scrolling. Controls apply accepted full snapshots; system theme listens through the media boundary; transcript/Actions and independent scrolling/layout roles are explicit. Controlled workbench assertions passed. This is source/synthetic-DOM evidence, not pixel, native drag or browser persistence proof. |
| Read Aloud / fractional speed | Actual production `SpeechQueue` reads accepted continuous speed at each utterance admission; existing active speech is not arbitrarily rerated. Independent fractional/manual/cancel-generation checks passed and a rounding mutant failed. Extended rate mapping is explicitly approximate; no native acoustic-rate proof or fixed preset policy is claimed. |
| Maintainability / comments / duplication | Production changes stay at settings, budget admission, workbench and speech seams rather than duplicating orchestration in the GUI. The substantive maintenance defect is R1's stale adjacent integration fixture. No additional mandatory production-code change found in this finite review. |

## 4. Fresh controls and real outcomes

All commands used scrubbed credentials, isolated cwd/HOME/TMPDIR, offline dependency resolution, single-worker builds, deadlines and owned process groups. Receipt argv, cwd, environment reference, stdout/stderr hashes, timing, exit and process termination are retained. No browser/native/live/paid endpoint was used.

| Control | Exit | Elapsed | Outcome |
|---|---:|---:|---|
| `root-race` | 0 | 153.389s | PASS |
| `root-vet` | 0 | 0.711s | PASS |
| `gui-race` | 0 | 1.577s | PASS |
| `gui-vet` | 0 | 0.180s | PASS |
| `consumer-race` | 0 | 0.125s | PASS |
| `consumer-vet` | 0 | 0.125s | PASS |
| `workflow-race` | 0 | 1.332s | PASS |
| `workflow-vet` | 0 | 0.179s | PASS |
| `stream-race` | 0 | 0.132s | PASS |
| `stream-vet` | 0 | 0.127s | PASS |
| `page-workbench` | 0 | 0.040s | PASS |
| `page-queues` | 0 | 0.038s | PASS |
| `page-inherited-integration` | 1 | 0.040s | FAIL — R1 |
| `chapter8` | 0 | 18.248s | PASS |
| `chapter1-lifetime-once` | 0 | 124.101s | PASS |
| `page-adapted-diagnostic` | 0 | 0.040s | PASS |

Both publication reports are **100/100, `Passed=true`**, against the same pinned root CLI. Chapter 1 full lifetime execution count is **one**. The module checks are `go test -race -count=1 ./...` and `go vet ./...` for the root, GUI, consumer, workflow and streaming roots. Their actual coverage includes modules with no tests; a successful no-test module result is build/vet evidence, not invented assertions.

### Independently strengthened real-process control and sensitivities

The real-process helper reuses the submitted fixture's reviewed process/WebSocket/fake-provider transport machinery, then adds independent wrong-type atomicity, actual filesystem save failure, immutable system-prompt and no-settings-as-user-text assertions. It executes owned binaries built from the actual copied production GUI/core. `real-process-helper.diff` exposes every helper change.

| Owned mutation | Old | Negative | Restored | Meaning |
|---|---|---|---|---|
| Active Agent turn rereads live budget instead of captured admission budget | compile 0 / behavior 0 | compile 0 / behavior 1 | compile 0 / behavior 0 | Negative fails `live update rewrote active turn`; it is executable wrong behavior, not a compiler failure. |
| Production SpeechQueue rounds the fractional rate at admission | parse 0 / behavior 0 | parse 0 / behavior 1 | parse 0 / behavior 0 | Negative fails the fractional assertion; synthetic synthesis only, not native speech. |

Both copies restore to their original source hashes. Budget old/restored executable hashes are identical: `0d86e089c4e4158a8598d08b08fe6bc29f3f3b78142a36caa37808de5b3b0f9e`. Complete source/binary/parse/exit identities are in `budget-mutation-results.json` and `speech-mutation-results.json`. There were exactly two production-code mutations; the third copied-file operation was the separately qualified fixture diagnostic, not a product mutant or candidate edit.

## 5. Carried, limited and deliberately not run

- Student local fixtures, initial failed process setup attempts, unmapped inherited-model Chapter 8 failure, later mapped pass, and expected student sensitivity failures remain carried and qualified in the unchanged student report/package. Fresh reviewer results above are not attributed to the student and do not erase its failures.
- Prior adapter/reconciliation sealed bundles remain verified unchanged (168 and 54 members respectively). This turn did not repeat entry controls or issue new adapter/API-preparation recommendations.
- Chapter 7 standalone publication, legacy entry matrices, another Chapter 1 full run, new malformed/write-failure campaigns, native GUI/browser automation, paid/live providers, secret/key/config inspection and authoring are **not run**. The one real save-failure probe belongs to the finite settings atomicity control above, not a new campaign.
- Actual browser/native layout, OS-theme behavior, real drag/focus persistence, microphone and speech audibility/rate, and any later live provider acceptance remain separate acceptance stages. Their absence is a scope boundary, not a failed test or a demand for routine fresh permission.
- The retained 860,160-byte `ch08-student-evidence/source.tar` (`a0783f71ee167f54ace212d11a57a75ab0ada8871592a3ac5f186a9cf6e463fe`) is the already-qualified local-only TAR mistakenly placed in portable evidence. It is preserved/excluded from commit and the old seal stays unchanged. It was not extracted, repacked, copied into this bundle or described as runnable.
- Preparation/read-only mistakes are preserved in `review-attempt-notes.json` and the original scripts/receipts: absent guessed paths; selecting nine retained build records instead of the requested six runtime identities; initially assuming a nonexistent report flag (no controls started); and a wrong summary receipt-key name. The corrected grader plan consumes JSON stdout. No failed actual control was silently retried into a green original receipt. The sole unexpected runtime failure remains R1.

## 6. Inventory, preservation and next action

Authoritative portable index: **`ch08-code-review-evidence/inventory.json`**, self-excluded, with exact report/evidence paths, byte counts, SHA-256 and classifications. The final seal verifies UTF-8/no-NUL/no-binary/symlink boundaries. `local-artifacts.json` separately inventories 349 retained owned scratch files (31,994,328 bytes), including binaries/tool effects; those bytes are not portable deliverables.

`final-verification.json` confirms all pinned preservation counts above, six unchanged runtime identities, the unchanged accepted predecessor tag, restored owned source copy, clean student, 28 terminal command receipts and **zero remaining owned process groups**. No shared/unrelated path was cleaned. Free space at that check: 6,186,065,920 bytes. `final-seal.json` records the closing check.

**Next action:** return R1 as a narrow student correction, then review the resulting exact diff and updated fixture evidence. Preserve the accepted production evidence where identities remain valid. This report is not code acceptance, browser acceptance or authority to edit the submitted files; no acceptance tag, commit or push was made.
