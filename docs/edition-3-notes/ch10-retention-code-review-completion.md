# Chapter 10 retention amendment — independent delta-review completion

## Verdict

**PASS for this bounded amendment review; recommend acceptance of the frozen retention delta.** No blocking product defect was found in the completed source assessment or focused execution. This is not whole-Chapter-10 recertification and does not reopen the original lifetime, grade, or whole-suite acceptance.

This completes the same review after explicit parent authorization to recover reviewer setup. It does **not** rewrite the earlier INCOMPLETE result or conceal its failures. Student correction rounds: **zero**. No candidate, manuscript, or paused Chapter 11 edits; no children.

## Identity and custody

- Candidate: `b03c28fdb1586f16dfb53849f8c763576ebb39ae`.
- Tree: `b2f86ba9fd48ff9a923a8a924a9612549aa762c0`, obtained with `git --no-pager rev-parse HEAD^{tree}` and compared with the task pin. The private input has no `tree` or `candidate_tree` field.
- Accepted parent: `729bc919e734405307b4ac0060e65c40e33652e5`.
- Scope: 18 changed files; source-identical 142-file owned copy at `/tmp/ch10-retention-independent-review-20261019/source`.
- Final candidate status is clean. Frozen and owned-copy bytes/modes match all 142 pins; owned copy has no extra files. Governing task, public inputs, and protected tracked files remain pinned.
- Original report SHA-256: `81d71fe3c27d3d87f15f39655cb530428661fcd23d9072a76c38cb136e19c881`.
- Original eight-member inventory SHA-256: `cb941b59e0cda24408f60fbec5069e98466f69482321d38bc6cc0f05bba70108`.
- Student inventory SHA-256: `1f48c96a14b7d7825671314f9d2e318e63223c1490c88669b40d6593c4fc56b3`; all 34 members reverified. These are carried evidence, not new execution claims.

`preservation.json` records final checks. The earlier report/evidence remain byte-for-byte unchanged. Historical transcripts were not recopied. Unchanged prerequisites were not reread wholesale; their governing pins and prior reading ledger were retained, while unfinished source assessment and relevant public amendment clauses were completed. Earlier clipped source windows were recovered.

## Source assessment against the public amendment

All 18 changed files and necessary adjacent ownership/fixture code were assessed; the file list and citations are in `source-assessment.json`.

1. **Settings and runtime policy.** `internal/common/retention.go:5–23` supplies `200`/`full` defaults, preserves explicit zero as unlimited, and rejects negative limits and unknown modes. `diagnostics.go:165–188` decodes the shared settings file onto defaults and rejects null retention fields. `gui/settings.go` stages and validates patches, persists before publishing/live mutation, and carries values through reconnect/restart. CLI and GUI construction read the same `cr/settings.json`. Runtime retention is separate from provider configuration and descriptive archive metadata.
2. **Diagnostic ownership and boundaries.** `diagnostics.go:15–163` serializes append and replacement under one owner lock. New writes frame whole logical records; recognized legacy starts retain continuation lines, and unknown legacy bytes are not guessed apart. Retention scans record offsets, copies the newest complete suffix to a sibling, syncs, and renames before switching descriptors. Zero does not trim. `internal/agent/save.go:34–50` establishes successful recovery or a fresh save before startup trimming. `internal/agent/actor.go:737–745` saves after admitted work has finished and only then trims on orderly exit. Save/load/trim errors surface; diagnostic trimming does not target canonical History storage.
3. **Archive representation and recovery.** `internal/history/save.go:15–38` observes context, frontier, and available facts under one History lock. `dialogue_only` writes the current snapshot with already-folded canonical tail empty; it does not pretend that derived text is replayable `Event` data. `projectDialogue` (`208–245`) carries prior derived records, retains visible human/assistant top-level text including ordinary signed text, omits private signatures/opaque/tool/system payloads from the derived records, and sorts by source sequence. Mixed-content events are not discarded wholesale.
4. **Validation, mode changes, and invariants.** `ReadSave` (`59–179`) stages validation before installing state: reduced archives require a snapshot; derived role/provenance/order and canonical tail are checked; snapshotless reduced/derived/incomplete history is refused. The reducer applies only events newer than the snapshot anchor. Runtime config and tool handlers are not loaded from descriptive declarations. Switching back to full retains available originals and prior derived dialogue without fabricating discarded originals. Repeated reduced saves retain earlier dialogue while folding a legitimate canonical tail.
5. **Test quality and ownership.** Selected tests exercise real Agent, History, CLI `run`, GUI hub/startup, and filesystem boundaries using existing offline helpers. Assertions inspect saved files, unchanged state on refusal, whole records, exact projection, frontier continuation, and settings persistence/failure. Helpers were inspected but not changed. No pass depends on a reviewer-added production seam.

**Findings:** none requiring student correction. The public contract is retention, not removal of private replay data from the canonical recovery snapshot. The documented single diagnostic-owner assumption is not an external multi-process writer guarantee.

## New actual execution

Exactly **four** Go invocations ran, serially, under Go `go1.25.4 darwin/arm64`. Exact argv, cwd, full scrubbed environment, timestamps, elapsed time, disk readings, source pre/post verification, actual `RUN` names, top-level `PASS`/`FAIL` names, and child exit status are in the corresponding JSON receipts; raw merged output includes subcase passes.

Environment: private HOME/TMPDIR; no inherited provider credentials/settings; `GOMAXPROCS=2`, `GOFLAGS=-mod=readonly -p=1`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, `GOWORK=off`, `GOENV=off`; installed Go/module/build caches, with no cache cleanup. Every invocation used `-race -v -count=1 -timeout=90s` and an anchored named selection. A 1 GiB disk floor was enforced; final free space was about 15.4 GB. Fixtures may use local loopback and dummy credentials; no live provider or browser was used.

| Receipt | Actual selection/result | Go exit | Wall seconds |
|---|---|---:|---:|
| `root` | Root plus `./cmd/ch01`: 10 named cases PASS; 25 actual RUN entries including 15 subcases | 0 | 4.740 |
| `gui` | GUI module: 5 named cases PASS; 15 actual RUN entries including 10 subcases | 0 | 1.545 |
| `negative` | Only `TestRetentionDialogueArchiveAndModeSwitch`: intended assertion FAIL on compile-valid mutant | 1 | 0.406 |
| `restored` | Only `TestRetentionDialogueArchiveAndModeSwitch`: restored candidate PASS | 0 | 1.358 |

Root/CLI named cases:

- `TestRetentionDiagnosticRecordsAndConcurrentAppends`
- `TestRetentionDiagnosticReplacementFailure`
- `TestRetentionStartupAndExitBoundaries`
- `TestRetentionFailedLoadSaveAndTrimBoundaries`
- `TestRetentionPolicyDefaultsZeroAndValidation`
- `TestRetentionDialogueArchiveAndModeSwitch`
- `TestRetentionMalformedArchiveAtomicRefusal`
- `TestRetentionLegacyWholeRecords`
- `TestRetentionDiagnosticsNeverTrimReplayFacts`
- `TestRetentionCLISettingsAndRestart`

GUI named cases:

- `TestRetentionSettingsUpdateReconnectRestart`
- `TestRetentionGUIActualStartupExitAndZeroRestart`
- `TestSettingsPatchAndLoad`
- `TestSettingsCommitBroadcastRestartAndFailure`
- `TestSettingsConcurrentPatchSubscribe`

### One independent negative and exact restoration

After both race baselines passed, the reviewer changed exactly one owned-copy line in `internal/history/save.go`:

```diff
-for _, d := range prior {
+for _, d := range prior[:0] {
```

This deliberately drops previously derived dialogue while preserving compilation and current-tail processing. It complements the student's diagnostic-count negative by challenging recovery/mode-switch preservation. The unchanged test failed at `retention_test.go:379`, **`folded tail lost`**: only source sequences 7 and 8 remained instead of four derived records. It was an assertion failure, not a build/setup failure, panic, timeout, or race report. `negative.diff`, `negative.txt`, and `negative.json` preserve it. The receipt runner exits zero when it successfully observes this expected negative; **the actual Go subprocess exit was 1**.

The original file bytes were immediately restored, then all 142 files/modes and absence of extra files were verified. Only the focused restored positive ran; it passed under `-race`. No alternative mutant, retry, helper edit, or extra suite followed.

The reviewer independently selected execution and mutation but reused unchanged student-authored tests/helpers. No independently authored test suite is claimed.

## Reviewer failures and mechanical recovery

- The original INCOMPLETE seal preserves the two nonexistent manifest-key failures. Parent authorized setup recovery without charging a student correction round. `setup-correction.diff` is an explanatory schema-use excerpt, not a product patch; the new runner uses actual `candidate`/`base`/`candidate_files` fields and the explicit task tree pin.
- Completion `plan.json` initially contained a guessed projection anchor and assertion string. Its pre-execution anchor check returned zero. No source was edited and no test launched. After inspecting the real helper/assertion, `plan-v2.json` corrected them; both plans, `plan-correction.diff`, and `preexecution-anchor-check.txt` remain sealed.
- The first final custody script resolved student-inventory members from the wrong directory and exited 1 before writing its result. `custody-path-failure.txt` preserves the failure. Inspecting the actual inventory showed paths relative to `docs/edition-3-notes`; `custody-path-correction.diff` records the mechanical path correction. Final verification passed. No test or fixture was repaired or rerun for this error.

There was **no genuine candidate assertion failure**. The sole failing Go command was the declared negative, followed by exact restoration and the permitted focused positive.

## Limits and seal

No full suite, full race sweep, vet, grade, original Chapter 10 lifetime, live provider, browser, or Chapter 11 run occurred. No exhaustive process-crash/filesystem durability or external same-path writer guarantee is claimed. Long-run behavior remains carried earlier evidence rather than new measurements. No product changes or student correction round were needed.

All owned command jobs finished; no children were created. The owned `/tmp` copy and review directories remain; no shared-cache or unrelated-worktree cleanup occurred. Only the authorized completion report/evidence and owned local resources were written.

New portable seal: `ch10-retention-code-review-completion-evidence/inventory.json`. Paths are relative to `docs/edition-3-notes`; the report and every completion evidence member are included with byte length, mode, and SHA-256; only the inventory itself is excluded. All members are UTF-8 and contain no NUL bytes. Parent verification and final acceptance remain separate from this reviewer recommendation.
