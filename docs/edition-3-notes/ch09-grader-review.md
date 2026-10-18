# Chapter 9 grader adapter — independent bounded review

## Recommendation and severity

**Review verdict: PASS within the bounded adapter scope. Recommend accepting this thin adapter candidate. No blocking or nonblocking adapter defect was found. Root retains acceptance, public-guide and student ownership.** This is not Chapter 9 student acceptance and does not reopen the entry audit, the erratum, or any legacy skill-lifetime campaign.

The frozen `/tmp/cadence-ch09-grader` actually routes Chapter 9 to the existing historical skills composition, keeps the Chapter 8 route working, and rejects the measured invalid inputs before target execution or build. One compiling dispatch mutation was detected from actual report behavior, despite its grader process returning success; an exactly restored source copy was rebuilt and positively exercised. No replacement checks, score changes, student API changes or production edits were needed.

Severity summary: critical/high/medium/low adapter findings **none** in this bounded review. Informational limitations and review-harness defects are recorded below rather than attributed to the candidate.

## Authority, prerequisite state and frozen identities

The complete authoritative task was read before work: `/tmp/cadence-ch09-grader-review-task.txt`, 5,369 bytes, SHA256 `f18c654c6897d1e9f041ca4aa7816d91509bb7a13428a4dd940c218811a8b3bf`. The complete frozen `ch09-grader-review-input.json`, current `AGENTS.md`, and reviewer skill were read. The prior full prerequisite reads were carried only after matching current hashes, with no compaction since those reads. The carried ledger includes the correctly measured 454-line Chapter 9 and 1,110-line Edition 2 carryover; over-EOF read endpoints are not file lengths. See `prerequisite-reads.json` and `pins-before.json`.

All frozen records matched before controls and after completion: 25 inputs, 126 accepted student source records and 29 protected records. The candidate executable remained 10,865,602 bytes, SHA256 `0438e74287c3fdf80b66dacef5ebb8fbc4fadd92d2d472ddf824a9579b73b73e`. Candidate `main.go` is pinned at `53c301ff0d0859622720357bc1b3d99f4848426173f6a182c378fda5111719d1`; the complete frozen diff was reviewed, not replaced with a reconstruction.

Repository refs remained:

- main HEAD: `5614faed78f10f12a940eef6573a3bfe50aa3f6c`;
- student HEAD and peeled `edition-3-ch08-r1`: `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`;
- student tree: `fa896caa404b87b08de61eb38ce079eb38138fbe`, with clean status.

All 15,287 initially tracked files remained byte-identical to the owned pre-control snapshot. No nonowned status entry changed. The original entry-audit inventory and all its listed contents remain sealed, as do the erratum inventory and its contents. Their inventory SHA256 values remain respectively `f6301805f129a1eedeb852baea7bd69dfcdda172278bab4d49304db4bb472bd3` and `614556ada99f3dfab5292e331f21c73031706d826cb26c7dfd34cf1b71e7e364`.

## Adapter inspection

The candidate changes only `cmd/grade-edition3/main.go` and its tests. It extends the supported set to 1, 4, 5, 6, 7, 8 and 9 and updates help, adds Chapter 9 to the existing required module/model/fixture validation gates, and adds the dispatch at `main.go:123-124`. The call is the existing `gradeedition2.Chapter9(bin, core, gui, guiModel, consumer, workflow, exercise, streamModel, example, models)` API, not a synthetic replacement.

In `main.go:63-103`, required module shape and nonblank model validation, consumer-case decoding and media-model decoding occur before `grade.Build` at line 105. The compiled guards below independently measured that ordering. The Chapter 8 branch is still separate. The changed candidate tests were read as source; their unit suite was not represented as a fresh run. The independently compiled controls are the fresh validation evidence. The `flags.Visit` presence checks and adjacent gate retain the Chapter 1/4/5/6 older-flag rejection rules, including explicit empty values; the final Chapter 1 fallback and Chapter 4/5/6/7 dispatches remain intact on inspection. Those older branches were not freshly rerun; the compiled Chapter 6 empty-GUI guard is the fresh older-flag example.

`internal/gradeedition2/ch09.go` retains eight skills checks and Chapter 8 parity. The existing builder passes through a supplied executable file rather than rebuilding it. The route therefore preserves the distinct executable target and explicitly supplied core identity; it does not substitute the executable's parent directory for the core source. Existing consumer/media readers retain their established schema checks and defaults. No helper implementation, check, public return type, score weight, or temporal contract was edited.

The report title remains the inherited `Edition 2, Chapter 9 — Skills`. That is accurately reported here, not relabeled as a new Edition 3 check suite. The thin adapter's task is routing and validation; the historical title is not a blocker in this scope.

## Real artifact composition and isolation

One unchanged 75-file copy of `solutions/ch09` supplied all real application modules at `/tmp/ch09-grader-review/source`. Its canonical file-map digest is `51b7ee4bb380ad29189c224e21a1f3134623733798bb7b3bb9f31e64da4694ec`; it matched the original before and after controls. Actual paths were:

| Role | Input or retained artifact |
|---|---|
| Explicit core-source identity | `/tmp/ch09-grader-review/source` |
| Independently built CLI target | `/tmp/ch09-grader-review/bin/cli` |
| GUI source | `source/gui` |
| External consumer | `source/ch05` |
| External workflow | `source/ch06` |
| External stream exercise | `source/ch07` |
| Stream and GUI models | Explicit `claude-sonnet-5` for both |
| Consumer case / media rows | Existing reader defaults, not new fixture protocols |

The CLI executable is outside the source directory and has SHA256 `42c73a32c3b4c5c5122ad538f91d72d03c4d8017789fb74a0966b24af8984ac0`. The observed GUI executable hash is `1f8e24c717ff8b13cd774cd70ecd9adcd244e24afc02e8815e315210b0011441`; consumer, workflow and stream binaries were also captured. `artifact-identities.json`, `generated-artifacts.json` and `passive-build-metadata.json` retain exact paths, hashes and passive `go version -m` observations. The observer recorded 41 transient executable observations without errors, but is a sampling observer, not an exhaustive process tracer; missing observations are not invented. Raw binary contents and full source maps remain in owned scratch, not portable evidence.

`environment.json` records the entire allowlisted child environment. HOME and TMPDIR were owned and initially empty; user/provider configuration and credentials were not inherited. Go proxy and checksum-network access were disabled, existing local caches were reused, `GOMAXPROCS=2` and `GOFLAGS=-p=1 -mod=readonly` bounded build resources. Existing helpers supply fake keys and local fake-provider endpoints. No browser, external provider, external search or child agent was used. Real controls used 300-second outer process-group deadlines; none fired. No internal timeout, assertion or scoring threshold was changed. Disk stayed above the 1 GiB floor; the final protected-state check recorded 3,389,177,856 free bytes.

## Fresh controls and outcomes

Every launch, argv, cwd, explicit environment override, PID/process group, exit status, duration and stdout/stderr byte/hash receipt is in `commands.json`. Raw output files are retained separately. There were 56 recorded child commands: four builds, two marker self-checks, 46 guard invocations and four real route invocations. There was no retry.

### Compiled validation and no-launch guards

The actual frozen candidate was used for all guards. A compiled target sentinel writes a marker and exits 73; its direct self-check produced that marker and status. A PATH-local `go` marker wrapper writes a separate build marker and exits 97; its direct self-check likewise succeeded. Self-check markers were moved to local retained evidence before guards. `marker-provenance.json` preserves their exact sources and marker contents.

The guard suite then tested both executable-sentinel and module-build targets for:

- help, unsupported chapter, unknown flag and an explicit GUI flag on an older chapter;
- each missing Chapter 9 required core/consumer/workflow/stream/GUI/model flag;
- whitespace-only stream and GUI models;
- absent core, nonmodule consumer and nonregular GUI `go.mod`;
- malformed, missing and schema-invalid consumer-case and media-model JSON.

It also tested no target and help without a target. All 46 invocations met their expected status: three help calls exited 0, and 43 rejection calls exited 2. Neither target nor build marker appeared in any case. This is compiled-process evidence, not merely a source assertion that validation precedes build. The marker mechanism was not used to synthesize passing application behavior. See `guard-results.json`, `guard-*.stdout/.stderr` and the self-check receipts.

### Real positive, regression, mutation and restoration

| Invocation | Actual grader exit and report | Independent route observation |
|---|---|---|
| Frozen `-ch 9`, real composition | Exit 0; 100/100; passed | Correct Chapter 9 title and all nine expected check IDs |
| Same frozen executable, `-ch 8`, same real composition | Exit 0; 100/100; passed | Four settings checks plus `ch7-parity`, correct Chapter 8 title |
| Compiling owned route mutant, `-ch 9`, same real composition | Exit 0; 100/100; passed **as Chapter 8** | **Rejected**: wrong title and all nine Chapter 9 IDs absent |
| Exactly restored owned route, rebuilt, `-ch 9`, same composition | Exit 0; 100/100; passed | Correct Chapter 9 title and all nine expected check IDs restored |

The Chapter 9 report's measured IDs are `initial-tools`, `system-prompt`, `ensemble-tools`, `load-skill`, `progressive-disclosure`, `depends-autoload`, `var-substitution`, `blocked-skill`, and `ch8-parity`. The eight skills checks earned 85/85 points and `ch8-parity` earned 15/15, not merely an unqualified 100-point report. Raw `frozen-ch09.stdout` and `route-restored.stdout` preserve their titles, points and earned values. The Chapter 8 regression preserves the original settings details and its explicit source/integration limitations.

The one mutation changes only `Chapter9(` to `Chapter8(` in the `-ch 9` dispatch in an owned 162-file adapter/dependency copy. No check, scoring or application source was mutated. It compiled successfully; it then actually ran and returned a Chapter 8 report containing `settings-roundtrip`, `settings-on-connect`, `settings-broadcast`, `settings-persist`, and `ch7-parity`. The detector compared the parsed report's title and ordered IDs to the expected Chapter 9 contract. Its calculated rejection status was 1. **The mutant grader itself exited 0; the calculated status is not a separate subprocess exit.** This distinguishes behavioral route detection from compilation failure or grep-only proof. See `route-mutation.diff`, `route-mutation.json`, `route-mutant.stdout` and `route-mutant-route-detection.json`.

Restoration matched all copied files, not just the mutated line. The complete adapter tree digest returned to `5af4ea8007250382168c52b7eecc63193149fb8fec0b23092f08bd7f69113e0a`, and `main.go` returned to the frozen candidate's SHA256. The rebuilt restored executable was positively exercised, not accepted solely from that hash. It was built from the owned copy and its binary hash differs from the frozen binary; both exact artifact identities are recorded, and no binary-equality claim is made.

## Failures, evidence caveats and limits

Two pre-control read commands failed on guessed historical filenames (`chapter9.go` and `solutions/ch09/consumer/go.mod`). The errors are preserved in `inspection-failure.txt`; actual filenames/module paths were then inspected. These were inspection mistakes, not failed candidate controls. Some combined helper listings were clipped; relevant portions were recovered separately, and no full-read claim is made for those combined listings. Required prerequisite reads were complete or explicitly eligible to carry by current unchanged hash.

The route detector's saved summaries have an unused `total: null` field because the actual report field is `max_score`. The raw JSON correctly records `score: 100` and `max_score: 100`; `normalized-observations.json` derives the correct field without a rerun. The typo did not affect title/ID detection or any test outcome. It and the calculated-vs-process status distinction remain visible.

This review does not prove every Chapter 9 source-level architecture obligation, arbitrary GUI behavior, live-provider behavior or new student readiness. Numeric inherited parity is not a substitute for those proofs. No public names, prompt conventions, return types, weights, legacy compatibility path, skill lifetime rule, or grader helper internals were redesigned. No second negative campaign was conducted.

## Completion and evidence boundary

All control processes were waited, the observer thread joined, and no recorded control process-group member remained. All passive metadata inspection children also ended. `process-check.json` and `pins-after.json` provide completion and protected-state receipts. Only the new report, its evidence directory and owned `/tmp/ch09-grader-review` scratch were written.

Portable evidence is UTF-8 text with no NUL bytes. `inventory.json` gives exact paths, sizes and hashes for the report and every portable evidence file except itself. `local-raw-artifacts.json` is a hash-only index of retained local binaries, raw source maps and NUL-separated status snapshots; those contents are not copied into portable evidence. Prior audit/erratum seals and existing final191/join2231 remain prior completed evidence, not new review results.

**Recommendation to root:** the frozen thin routing/validation adapter has passed this independent bounded review. Root may make the acceptance and guide/student-progression decision; this report makes no such ownership transition itself.
