# Chapter 10 student handoff — BLOCKED, required gate failed

## Exact candidate

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Source-only commit: `729bc919e734405307b4ac0060e65c40e33652e5`.
- Tree: `025ef3744b42c8f9daec34e76661edc40583ca6a`.
- Parent/accepted Chapter 9 commit: `251ce8848b2e64080f65ec23384447a2bd0fcd65c`.
- Exact 137-file Git archive SHA256: `8a4fc7d17b4371a06d3db91287115c1f2ac6027fa932d8ba89a2141018b50742`.
- Frozen CLI: `.ch10-local/bin/frozen-cli`, SHA256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33`.
- Frozen GUI: `.ch10-local/bin/frozen-gui`, SHA256 `fb506081f98c09d2e6da628edbb6cb922f067f9db24983ae856d2d1222dcbf53`.

Source, exact archive membership, per-file hashes, binary sizes/build info and final cleanliness are recorded in `ch10-student-source-manifest.json`. Normal Git hooks were not bypassed. This commit is in the **student** repository, not the main manuscript repository. No tag was made. No tracked source changed after the freeze; all final receipt source/binary before/after pins are retained. This is a failed-gate handoff, not acceptance or implementation-complete certification.

## Implementation ownership

`internal/history/save.go` owns strict envelope validation and staged snapshot/log installation. `internal/history/history.go` owns a checkpoint-seeded context reducer and sequence frontier; it retains the old JSONL inspection path. `internal/common/common.go` declares the four-field save vocabulary and explicit tool declaration fields. `internal/agent/save.go` owns actor-serialized explicit Save/Load, missing-only startup, same-directory atomic replacement and presentation-ID high-water seeding. Existing Agent construction and shutdown invoke the lifecycle. Runtime Config and executable tools are not installed from archive metadata.

Both clients select working-directory `save.json`, with `--save` selecting the same load/write path. Independent unspecified library Agent directories are unique; explicit DataDir remains caller-owned. Offline CLI JSONL inspection uses isolated temporary persistence. The CLI retains its inherited protocol and persists pending ephemera at EOF. GUI stdin and WebSocket use the same Agent; EOF closes admission, drains and shuts down. GUI history projection includes context-only turns and deterministic presentation-only repairs for colliding archived IDs. `Ephemeral` is an explicit durable mailbox input; existing Hint semantics remain unchanged.

The source tests cover snapshot/full-log/tail/context-only behavior, strict atomic refusal, rich neutral parts and same/cross-vendor rendering, empty primary instructions, independent Agent directories, caller settings/tools, client lifecycle and GUI replay identity. These are observed local tests, not proof that the required external gate passed.

## Observed final-source evidence

All run receipts are under `ch10-student-evidence/`, including exact argv, cwd, clean/offline environment, duration, exit status, output hashes and source/binary before/after pins.

- Root test: exit 0, 134.106 seconds. Root race: exit 0, 158.871 seconds. Root vet: exit 0.
- GUI, Chapter 4 consumer, Chapter 6 workflow and Chapter 7 stream modules: test/race/vet exit 0. Consumer and stream report **no test files**; they are not represented as runtime test coverage.
- Compile-valid sensitivity: an isolated copy disables checkpoint seeding; `TestChapter10SnapshotReplayTailAndTrim` fails on missing snapshot context (exit 1), not a compilation error. Restoring the exact file and verifying all 137 source-file hashes produces exit 0. The real frozen worktree was never mutated.
- Actual CLI/GUI process control v4: exit 0. Local fake HTTP and real WebSocket exercise CLI default/override restart and malformed-file preservation; GUI stdin/WebSocket shared conversation, one-use/pending ephemera, EOF drain, replay/reconnect and distinct continued part IDs. All owned children exited. This is **not browser observation** and uses no live provider.
- Required Chapter 10 grader: **80/100, passed=false, exit 1**, 17.800 seconds. Save shape, replay equivalence, tail-once, context-only continuation, malformed refusal and corrected Chapter 9 parity passed. `default-load` failed for both CLI/core and GUI, with the message `historical model/vendor/prompt replaced runtime configuration`.
- Chapter 1 on this candidate: **NOT RUN**. The sequential script stopped on required Chapter 10 failure; root explicitly directed no Chapter 1 run at this boundary. Its required 62-second header/body phases therefore have no final grader result. Root-suite duration is not offered as a substitute.

## Unresolved diagnostic, not a speculative fix

Root authorized bounded diagnosis using public Chapter 10 text, invocation help and own code/tests. The student's hypothesis that Vendor wrote integers was incorrect: the inherited Vendor marshaler already uses the public string spelling. No source was changed in response.

`vendor-disk-negative-v1` remains an **exit-1 failed expected-negative assertion**, not a passing test. Its actual local process outputs show the initial fake request succeeded, and a restart after historical config was set to openai/historical-model/historical instruction produced a second normal fake response and a save recording current anthropic/model/prompt values. The diagnostic's first request assertions executed; the positive-branch assertions for the second request did not execute. The second subprocess exit status and full request body were not separately persisted, so the report does not upgrade them into a fresh passing assertion. The compound expected-refusal assertion failed. Existing retained observations disprove the proposed numeric-vendor explanation but do not explain or override the required grader failure. Root directed sealing this exact candidate for independent diagnosis, with no validation weakening or rerun just to get a positive label.

## Failures, deviations and assistance retained

- Initial development Chapter 10 grade: **65/100, exit 1**, preserved unmodified in `ch10-grader-ch10.json` and companion files. Its candidate binary remains pinned; its uncommitted source tree was not captured exactly. Its invocation omitted the explicit media-model map. Later corrections addressed client-owned job paths and runtime-system handling; this earlier result is not relabeled as final evidence.
- Early compilation/test failures and their available logs remain in `development-log-index.json` and `development-logs/`. Missing original timing/environment data is not invented.
- The 136.674-second pre-freeze failed root suite ran while source edits occurred. It is explicitly **not final-source evidence**. Later targeted checks and all frozen gates are separate receipts.
- Early commands omitted some assigned clean/offline/resource settings. Subsequent pinned gates use the recorded clean environment, `GOMAXPROCS=2`, `GOFLAGS=-mod=readonly -p=1`, local toolchain, GOPROXY/GOSUMDB off, and the required PATH including Delve. The bounded process controls use only dummy fixture credentials and provider-specific localhost endpoints, not generic endpoint pollution.
- A recursive grep unintentionally entered ignored `.ch09-local` predecessor-copy directories and emitted test-line matches. It was disclosed immediately; root authorized continuation. No further inspection of those copies or private solutions was performed. Later search used explicitly tracked paths.
- Four early grader receipts and two notes were written outside the four permitted main output paths. Root authorized only their exact byte-preserving relocation. Both relocation manifests retain original paths, before/after sizes/hashes and ownership. No shared cleanup was performed.
- The disk guard stopped a targeted run before launch below 1 GiB. All owned jobs were finished at pause. Root later relayed Bill's explicit permission to continue after freeing space. Resource-interruption evidence is retained; no completed gates were replayed as part of resume.
- Freeze-runner setup had a shell quoting failure before builds; preserved separately. Process-control v1/v2/v3 failed on harness assumptions (`asof` versus `as_of`, CLI ephemeral envelope, `id` versus `part_id`). Their source versions and receipts remain unchanged; only ignored harness code was corrected before v4.
- Root supplied scope/sequence, path/resource authorization and this stop/seal decision, not implementation patches. No children, private grader source/fixtures, historical/final solutions, runtime agent, browser, live provider, credential file, manuscript edit, main-repository commit or tag was used. Formal independent code-correction rounds used: **zero**. Root's retained independent reviewer has not yet reviewed this handoff.

## Remaining work

The required Chapter 10 default-load/runtime-config failure is unresolved. Root must diagnose against this frozen failed candidate. Chapter 1 true+exit0 with unshortened lifetime phases remains outstanding. Browser/provider and runtime-agent observation remain deferred and require separate root authorization. No additional student behavioral diagnostics are running or authorized at this handoff.
