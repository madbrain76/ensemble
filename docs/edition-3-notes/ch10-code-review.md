# Chapter 10 whole-candidate code review — qualified / incomplete

## Verdict and root decision boundary

**NOT ACCEPTED as a completed whole-candidate review.** The required same-CLI gates and the assigned standard module gates are green, but this reviewer did not successfully execute the required independent high-risk controls or behavioral sensitivity. Two reviewer-owned probe attempts were invalid. **No candidate must-fix defect is proved by those failures, and no student code correction is justified from them.**

This is not another fixture review and is not a declaration that the candidate is defective. It is an incomplete first whole-candidate review. Root should not infer whole-code acceptance from the green gates or source-read volume, and should not infer a product defect from a stricter/invalid local fixture. The remaining independent review obligations need a valid review before acceptance. Operational exercises, adoption, publication, and any formal review/correction round remain root-owned and pending. The student remained idle and unchanged.

## Frozen subject and preservation

- Checkout: `/Users/bill/projects/ensemble-edition-3-student`.
- Candidate commit: `729bc919e734405307b4ac0060e65c40e33652e5`; tree: `025ef3744b42c8f9daec34e76661edc40583ca6a`.
- Accepted parent from Git/source manifest: `251ce88fb1677897dbcc6b33661096c5cb39cd8e`. No old parent-identity typo was used.
- Exactly **137 tracked source/support files**, with **19 changed paths** relative to the parent.
- Frozen CLI: `.ch10-local/bin/frozen-cli`, SHA256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33`.
- Corrected grader: `/tmp/cadence-ch10-grader-prompt-fixed`, SHA256 `394a51635fe8538653fe9e49025c9c3817612c05b67106effd2ecbd445ab9fb2`.
- Input manifest: `ch10-code-review-input.json`; source manifest: `ch10-student-source-manifest.json`.
- All 31 inputs, 33 protected files, 137 student files, frozen binaries, both existing archive identities, and seven prior inventories were reverified at each completed run tail. Existing main/student sources and prior evidence were not edited.

**Provenance qualification:** the owned copies were populated from the 137 hash-verified frozen checkout files, not extracted from the pinned source TAR as requested. Both archive identities were verified and left local/unchanged. The copied bytes were independently checked; final cleanup also proved original modes and exactly the original 137-file set. This establishes byte identity but does not erase the requested-procedure deviation. No new archive or binary is exported.

## Prerequisite and source reading

Unchanged governing/craft/historical prerequisites were carried only through exact hash comparisons against previous complete reads. Changed/new student reconciliation and parent-decision records were read. Current Chapter 10 and applicable Chapter 2/9 obligations governed the review; older author-first directions did not authorize changes or broader work.

All 19 actual changed files and the complete approximately 64 KiB frozen diff were read, including recovered truncation. Task example filenames were reconciled with actual Git paths. Adjacent public interfaces, Agent/jobs/skills ownership, all three native request render paths and stream/parsing boundaries, GUI HTTP/WebSocket transport, and browser projection code were inspected. No browser was opened. `changed-file-read-coverage.json` records exact source identities; `read-only-preflight.json` and each run-version `verification-before.json` preserve exact input/carry pins; the older complete-read ledgers remain in their verified prior packages.

The failed fixtures show that full reading was **not sufficient proof of accurate API use or behavioral completion**. This report does not turn source reading into a pass.

| Changed file fully read | Bytes | SHA256 |
|---|---:|---|
| `.gitignore` | 478 | `629492025dd79f316a4a0739825e2a516e7369c31e84c66bc5bf7066b4762971` |
| `CHAPTER10.md` | 3295 | `2cb90730395f04ddb4fb7debd178aebe4520e2bf4a58dd604a64bd2634437453` |
| `chapter10_test.go` | 13188 | `cedbcab8dfc057c92af0be67bf06c6a58bd7fc80a1109778bf153acfac2dc40a` |
| `chapter2_test.go` | 22205 | `6826740933a33b2ce175c98f1cf70556792a56a342eb219c45654b79a66242d6` |
| `cli_test.go` | 4490 | `6d5d0eaab26708e629fa56d058bf35ec1ccd4e4d0da74b0ad10308c25524ffc4` |
| `cmd/ch01/main.go` | 14951 | `0b5d00910bfb7c47508d7b09b1937962f5432867eb15cc37cbb5171d8a9e9323` |
| `cmd/ch01/main_test.go` | 11362 | `c4efa9be1713a8c5c219361ae826881f875adaf93ed8bfbf7b626487ac85bd52` |
| `ensemble.go` | 13456 | `b35a55819b47b1c4468d1e04d70005736f53982d887a990025502cebb4867802` |
| `gui/README.md` | 7473 | `147e3c40645dc4d400a73443dbad76493b0e38061ebfe227f4082c4c08d9e202` |
| `gui/chapter10_test.go` | 6284 | `231688834ee61014949cbe0cb568bdd64a0932565396a13039ce35319ca46999` |
| `gui/hub.go` | 10597 | `d1c28ea1c55821c0bc0d1688b7c85fd85452e6a96fa4b6221caf820753549d64` |
| `gui/main.go` | 8017 | `4b11e676333156bf795a97a3cd35dd76a6c83bf2a5b64d0cbfb1e42024535c34` |
| `internal/agent/actor.go` | 23886 | `3d22b923ec0d9b3256c7ee1e7746944ccaa4cbe8f028fd156df35c9644a78329` |
| `internal/agent/agent.go` | 9233 | `117180139d277fb7b5f5b5d033cce953c305e62b8ac7a8ceb1ba9f0ce2b24426` |
| `internal/agent/save.go` | 3487 | `447cb4705175e0efe23c10ec4bf8d782357df8eca1dd877484a6764f98e438e4` |
| `internal/common/actor.go` | 7649 | `2d91d34e91d5f35e907042fce0bc2121a7d2756552148d1b55ee6c923565fb02` |
| `internal/common/common.go` | 18747 | `071f7b7b5963e012378dd3923cdde031b3e2ae014906d975f0684772723645d6` |
| `internal/history/history.go` | 14802 | `03930e8dcd36dcf13b594939e3fa309b5b2c136977da1aebc73d905f2c7c9eaa` |
| `internal/history/save.go` | 4498 | `91d64c22516ad474ee32267c331aec33704bc6ee1587eb6f7cce4a56317724a3` |

## Required same-CLI gates

### Corrected Chapter 10 — carried, not rerun

The accepted fixture-review report and exact output/receipt pins were verified. On the same unchanged frozen CLI, the corrected grader returned **exit 0, passed true, 100/100**, exact title and seven expected checks. This task executed **zero** new Chapter 10 gates. It did not replay entry, prompt-fixture, or historical compatibility controls. Earlier 80/100 failure records and the measured LF negative/contrast remain intact; none was rewritten.

### Chapter 1 — one fresh execution, actual 62 seconds twice

`v2/01-fresh-ch1.*` and `v2/01-parsed-report.json` retain the actual fresh result:

- Command: corrected grader `-ch 1 -json` against the exact frozen CLI.
- PID **94354**, process exit **0**, no timeout, elapsed **124.19327 seconds**.
- Report title: **Chapter 1 — Edition 3 conversation and request lifetime**.
- **Passed true, 100/100**, actual nine check IDs: `protocol`, `wire`, `calls`, `replies`, `memory`, `growth`, `usage`, `generation_headers`, `generation_body`.
- The final two checks explicitly report a **1m2s headers delay** and **1m2s body delay**, with no timeout/cancel and the text recovered.
- Stdout SHA256: `f5f3b049294ce0f224cdf86853eee291fdcf937ce67daf30deedc7fa4b73019d`; 2,959 bytes. Stderr was empty.

My wrapper incorrectly expected 15 check rows and stopped **after this real PASS**. That wrapper failure is preserved in v2. Root `main-32` authorized carrying the actual nine-ID result and continuing only unrun controls. **No lifetime or Chapter 1 retry occurred.**

## Fresh bounded module gates

These ran sequentially on the owned byte-pinned copy, with `GOMAXPROCS=2`, `GOFLAGS=-mod=readonly -p=1`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, clean HOME/TMPDIR, explicit Go/module cache paths, and pinned Go/Delve resolution. No inherited provider credentials were admitted; no paid/provider endpoint or browser was used. Full argv, cwd, environment, PID, deadlines, elapsed times, actual exit codes, stdout/stderr and hashes are retained.

- Root: `go test -race -count=1 -timeout=240s ./...`, then `go vet ./...` — **PASS**.
- GUI module: separate `go test -race -count=1 -timeout=120s ./...`, then vet — **PASS**.
- `examples/ch04-consumer`: test/vet exit 0, but **no test files**; compile/static coverage only.
- `examples/ch06`: test/vet exit 0; actual workflow tests ran.
- `examples/ch07`: test/vet exit 0, but **no test files**; compile/static coverage only.

There was no full legacy suite or arbitrary coverage campaign. Ten prescribed module commands completed once; none was rerun after reviewer failures.

| Version / actual command label | Exit | Seconds | Classification |
|---|---:|---:|---|
| v2 / `01-fresh-ch1` | 0 | 124.19327 | PASS (actual check/test scope below) |
| v3 / `02-root-race` | 0 | 159.175948 | PASS (actual check/test scope below) |
| v3 / `03-root-vet` | 0 | 0.513532 | PASS (actual check/test scope below) |
| v3 / `04-gui-race` | 0 | 2.190186 | PASS (actual check/test scope below) |
| v3 / `05-gui-vet` | 0 | 0.250251 | PASS (actual check/test scope below) |
| v3 / `06-aux-1-test` | 0 | 0.36917 | PASS (actual check/test scope below) |
| v3 / `07-aux-1-vet` | 0 | 0.153402 | PASS (actual check/test scope below) |
| v3 / `08-aux-2-test` | 0 | 0.478388 | PASS (actual check/test scope below) |
| v3 / `09-aux-2-vet` | 0 | 0.145171 | PASS (actual check/test scope below) |
| v3 / `10-aux-3-test` | 0 | 0.142899 | PASS (actual check/test scope below) |
| v3 / `11-aux-3-vet` | 0 | 0.191147 | PASS (actual check/test scope below) |
| v3 / `12-independent-public` | 1 | 0.189225 | Reviewer compile failure; no assertions |
| v4 / `12-independent-public` | 1 | 0.677179 | Reviewer constructor failure; no persistence assertions |

The ten module-gate stdout summaries are retained; no-test modules are not advertised as tested behavior. Go vet's empty output is recorded as an exit-0 static check, not a runtime proof.

## Independent review controls: invalid attempts and unmet obligations

### First owned public probe: did not compile

`v3/12-independent-public.stderr.txt` records reviewer errors: invented `Redaction`/`SummaryReplacement` fields, wrong `Render` value-versus-pointer use, and treating `History.Append(Event) error` as a sequence result. No behavioral assertion executed. This was **not** a valid negative control, not a student failure, and not mutation sensitivity.

### Root-authorized mechanical correction: compiled, but constructors rejected the fixtures

Root `main-33` authorized one versioned mechanical correction only. The original sources and failed compile remain. `v4/*-probe-correction.diff.txt` records use of actual `Redact`/`Redaction` levels rather than an invented compaction API; JSONL `Dump`/`Load`; `Render(*Context)`; checked Append errors; `[]Part` ephemera; and required `turn_ended` outcomes.

The corrected command compiled and ran, but all relevant paths stopped in fixture construction:

- The standard/rich/atomic fixtures omitted a required positive `MaxTokens`: **“max tokens must be positive.”**
- The empty and whitespace primary cases attempted constructor configurations rejected with **“system prompt is required.”** They did not establish a persistence or recorded-metadata failure.
- **No replay-mode, rich-render, failed-load atomicity, regular-file save-atomicity, or frozen-prompt persistence assertion was reached.** `RUN` lines alone are not credited as those behaviors running.

`constructor-validation-source.json` records the actual validation locations and predecessor-window comparison. Rejecting an invalid constructor is not a violation of the persistence contract, and the task does not authorize imposing an arbitrary stronger local constructor policy. The final corrected run exited 1; the wrapper stopped because it lacked a green atomic-oracle baseline. No further correction/retry was made.

### Explicitly NOT RUN

1. A valid independent public checkpoint-only/null-context/zero-frontier/strict-tail replay and rich-render comparison.
2. A valid failed `ReadSave`/`Load` and existing-regular-file `Save` sentinel check.
3. The drafted real GUI WebSocket legacy-collision/reconnect/restart/new-generation check.
4. The planned one compile-valid early-truncation mutation on an owned `internal/agent/save.go` copy.
5. A restored positive run after that mutation.

**No behavioral mutant was applied.** All 137 copied original bytes remained/restored exact; this is preservation, not a completed mutation-restoration experiment. The required sensitivity obligation is unmet. Old fixture/entry negatives cannot substitute for it.

## Source observations and limits by public obligation

The applicable obligations remain those in public `books/edition-3/10-persistence.md`, inherited Chapter 2 replay/projection and Chapter 9 primary/skills semantics, and the actual public API. The following is an inspection summary, **not a list of proved product defects**:

| Area | Source reviewed / evidence | What is not established here |
|---|---|---|
| Four-field save document, typed parsing, frontier and staged installation | History/save code, public interfaces, submitted tests, carried full gate | Independent unusual-cache/strict-tail/public-load reproduction did not reach its body |
| Rich archived facts, summaries/redactions, ephemera, usage, nested references and vendor-bound opaque data | Full history/common changes, Chapter 2 tests, three renderer paths; fresh root race | No independent valid rich replay/render comparison or real-provider acceptance claim |
| Immutable/current configuration and primary instructions | Agent save/load path, Root setup, submitted tests, carried corrected gate and earlier measured prompt evidence | Empty/whitespace invalid constructors do not test older saved metadata or prove a policy defect |
| Error atomicity and file replacement | Temporary-file/write/sync/close/rename path and staged reader inspected; submitted root tests passed | Existing regular-file failure oracle never reached; no mutation sensitivity |
| Ownership, lifecycle, publication, CLI/GUI flush and absence-vs-malformed load | Full Agent/actor/public entry changes; fresh root/GUI race and same-CLI lifetime; carried Chapter 10 | No whole-code acceptance inferred solely from these gates |
| Legacy repeated display IDs and tool correlation | Full hub/presentation changes and adjacent browser consumer; real GUI regression suite green | Draft independent WebSocket collision control did not run; source inspection is not its substitute |
| Star topology and declared consumers | Root/common interface changes and actual module map; test/vet gates above | Two auxiliary modules have compile-only rather than behavioral test coverage |

I considered the zero-frontier/non-null-cache boundary for independent verification, but the fixture failed before a valid Agent existed. **This report does not promote that concern into a must-fix finding without a valid public-contract reproduction or complete supporting proof.** Likewise, no new compaction API or constructor behavior is required by this report.

## Prioritized blockers to completing this review (not student defects)

1. **Independent high-risk coverage and sensitivity are missing.** Root must obtain valid completion of the original required review obligations before using this attempt as a whole-code PASS.
2. **Reviewer fixture validity was not established.** The attempted assertions must be grounded in actual public constructor/History/Engine contracts; their failures must not be given to the student as a correction brief.
3. **Provenance method deviation:** checkout-copy rather than requested TAR extraction, despite exact file equality, must remain visible in the acceptance record.

There is no justified “smallest student code repair” from this attempt. A new root decision, not automatic continuation, is required for further review work. No private fixture, mutation implementation, or invented API expectation was sent to the student.

## Failure, assistance, and authorization ledger

`review-and-failure-ledger.json` records all actual commands, qualifications and unmet requirements. Retained evidence includes:

- Recovered read-only filename/revision/schema mistakes and truncation gaps; these were not runtime passes.
- A final-reconciliation inventory `entries` versus `files` schema mistake, corrected before controls.
- A wrong inferred example-module path that stopped a zero-command attempt; actual module paths were recovered from the pinned full argv.
- The post-PASS Chapter 1 row-count assertion and root-authorized no-replay continuation.
- Both invalid public-probe attempts, their actual exit-1 output, and the root-authorized one mechanical correction.
- A correction-writer filename/glob error before any corrected test ran.
- A cleanup guard failure: temporary-file hashes had been computed only in memory when the exact source-file-set check found two extra owned API logs. **That exhaustive per-temp-file hash list was not durably saved and is not claimed.** Previously sampled transient binary identities remain explicitly non-exhaustive. The two remaining logs were durably hashed before removal; their raw contents were never exported.

A first evidence-only finalizer also used the wrong process-group receipt key (`pgid` rather than `process_group`). Its failure, original script and one-field correction are retained; no inventory or runtime control was produced by that failed finalization.

No child agents, browser, paid provider, keys, public guidance edits, student edits, grader edits, shared cleanup, extra lifetime/gate replay, or automatic follow-on task was used.

## Cleanup, preservation, and inventory

All top-level jobs and samplers are finished. Both recorded owned TMPDIRs are empty. The two reviewer-added test files and two test-created API logs were removed only from the owned copy. Both retained local source copies now contain **exactly the original 137 files with matching bytes and modes**. No owned process remains. Shared caches were not cleaned; existing student/main source, archives, frozen binaries and old evidence remain unchanged.

`owned-cleanup.json`, `remaining-owned-log-identities.json`, and `cleanup-failure-ledger.json` preserve the cleanup facts and the lost-in-memory-hash qualification. Disk stayed above the required floor at command boundaries; final cleanup observed roughly **2.62 GB free**.

The portable payload is this new report plus the recursively enumerated new `ch10-code-review-evidence/` files. The inventory excludes only itself and records exact relative paths, byte sizes and SHA256. Every included file must decode as UTF-8 without NUL; no archive, executable, raw settings payload or file named `settings.json` is exported. Local source copies, archives, binaries and removed/raw runtime artifacts are local/hash-only as qualified. Root owns inspection, sealing/tail/join, acceptance and any formal round.
