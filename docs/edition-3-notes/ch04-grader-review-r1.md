# Chapter 4 command adapter R1 — focused independent re-review

## Decision: ACCEPT

R1 is resolved. The frozen adapter is accepted for the stated **CLI adapter scope only**. No remaining adapter finding. This is not Chapter 4 student acceptance, a manuscript review, or closure of the entry audit's historical-grader coverage gaps.

### Exact reviewed inputs

- `cmd/grade-edition3/main.go`: SHA-256 `539ee3e6cc7313171ab57d268803ddd4c2d49be33fb028c38fa43b3b31c12834`.
- `cmd/grade-edition3/main_test.go`: SHA-256 `c94e75cc895bfcf5e2625b523ff5571955c565ad8ef9f261b1f7ae3f7637fa2f`.
- Parent manifest inspected: `docs/edition-3-notes/ch04-grader-evidence/candidate-r1.json`.
- Base and final main HEAD: `bac07e644facb25ba3e5bffb2961674763fed4c4`.
- Independently built **new R1 binary**: SHA-256 `5234579fcb72b4faf6426b9757fa44c2630d09227762c76e4ae9b12dbcd39647`.
- Toolchain: `go version go1.25.4 darwin/arm64`.

## Delta and source review

Independently removing only the three-line `if err == flag.ErrHelp { return 0 }` branch reproduces the previous `main.go` SHA-256 exactly (`d08616eee707b4b91363f528abafa85f9b2268af964a0c6d903efbae95c291a4`). Removing only the appended `TestHelpSucceedsWithoutLaunchingTarget` reproduces the previous `main_test.go` SHA-256 exactly (`f505aebe9a8ef9bea5495c6ee723eea0880b1ae3735e4b36adadff2cceb71a5f`). The new test exercises four help forms, including target sentinels under Chapters 1 and 4.

The branch handles the standard parser's help sentinel with successful status 0, before build/dispatch. Actual parsing errors remain status 2. Local flag state, explicit Chapter 4 source/consumer validation, optional case handling, build cleanup, Chapter 1/4 routing, report scoring, and all non-help behavior are byte-for-byte unchanged. This is the smallest appropriate fix; it adds no framework, API, application requirement, or implicit target.

All 157 previously recorded historical grading/provider source hashes were independently reverified unchanged. A fresh scratch export was extracted from the hash-verified prior Git archive and overlaid with only the two new frozen files. No live source was modified.

## Fresh bounded R1 controls

| Control | Actual result |
| --- | --- |
| `go test -race -count=1 -v ./cmd/grade-edition3` | exit 0 (4.035 s): existing 18 invalid-input subcases and new four-form help test pass |
| `go vet ./cmd/grade-edition3` | exit 0 (1.120 s) |
| Build new R1 adapter | exit 0 (0.601 s) |
| Compiled `-h`, `-help`, `-ch 1 -h sentinel`, `-ch 4 -h sentinel` | all exit 0, display help, no sentinel launch |
| Compiled unknown flag and Chapter 4 missing target | both exit 2, no sentinel launch |
| Compiled Chapter 1 with `-core=`, `-consumer=`, `-consumer-case=` | all exit 2, no sentinel launch |
| Scratch mutation removing only help branch while retaining new test | targeted test exits 1; **all four** cases report status 2 instead of 0 (0.287 s) |
| Restore exact frozen source, rerun targeted help test | exit 0 (0.276 s); restored source hash matches frozen R1 |

All 14 command receipts are finite and their output hashes verify. No unexpected failure, timeout, or launch occurred. The mutation failure is intentional, establishing that the new test detects the original regression. Scripts, source diffs, exact argv/cwd/status/timing/output hashes and completion records are in `ch04-grader-review-r1-evidence/`.

## Carried full-route evidence — NOT rerun on the R1 binary

The prior independent review remains at `ch04-grader-review.md`, SHA-256 `5a589d5f83692e556fc2b0010f668165def7fce65b6b8f1cebc8a3ae27cf0b1b`. Its applicable output hashes were reverified and copied as references into `carried-route-evidence.json`. These are explicitly measurements of **prior adapter binary** SHA-256 `3f7af0cd6df371069ceae55240ca6a670d1a102cb8cf1e434fd0e185bb228a92`, not the newly built R1 binary:

- Explicit executable/default consumer case: **160/160, passed=true**.
- Supported directory target/non-default consumer case: **160/160, passed=true**.
- Wrong core with the same module name: **135/160**, only topology fails.
- Valid case with deliberately wrong expected result: **135/160**, only external-consumer fails, proving case forwarding.
- Private core-forwarding mutation: **135/160**, only topology fails.
- Chapter 1 on the **same immutable reference executable** used for Chapter 4: **100/100, passed=true**, including both actual **62-second** delayed-header/body phases; elapsed 124.046 seconds.

That paired reference binary's SHA-256 remains `c7a958a2df3e208f32644ca72f48be08fb79880f44eaa603f4a1894cf78c4b21`. It was a scratch-only copy corrected solely for the known legacy 60-second HTTP timeout; original reference failure and student root-library archive/launch evidence remain preserved in the prior review. No new reference change or full-route/lifetime rerun was performed in this focused turn.

Carrying these controls is justified by the exact reverse-delta proof: only help-error handling and its test changed; all successfully parsed Chapter 1/4 routes and historical grading code are identical. Claiming fresh R1 execution of the 160/100-point or 62-second phases would be incorrect.

## Scope and preservation

The adapter now supports explicit executable/core/consumer identity separation through the versioned E3 command while reusing the existing Chapter 4 API and 160-point contract. Edition 3 lifetime remains a separate `-ch 1` invocation on the same future submission binary. The entry audit's selected-model, pending-cap, cursor/concurrency/process-group, registered-tool lifecycle, and ownership review obligations remain applicable to the future student.

Final frozen input hashes match. Accepted student HEAD and peeled `edition-3-ch03-r1` remain `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, working tree clean. No student launch, student/manuscript/historical-grader/original-reference edit, paid/live call, credentials, installation, push, restart, or concurrent-work removal. Only this report, its small evidence directory, and disposable scratch files were written. The reviewer is finished, not waiting for a student candidate.
