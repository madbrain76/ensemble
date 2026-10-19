# Chapter 10 current-prompt fixture correction — independent bounded review

## Verdict and boundary

**PASS for the assigned fixture-correction review.** The source change is exactly the removal of one terminal newline from the synthetic current primary-skill body. It does not weaken or otherwise change the persistence assertions, score composition, deadlines, adapter, public guidance, or frozen student.

The two authorized runtime controls each executed once and passed:

1. The corrected full Chapter 10 gate on the unchanged frozen student CLI returned **exit 0, `passed: true`, 100/100**, with the exact historical Chapter 10 title, seven IDs, and original weights.
2. The existing focused historical CLI/GUI persistence test actually ran and passed, **exit 0**. The task initially named a nonexistent test; root corrected the selector before any runtime control.

This is not an independent whole-code review or Chapter 10/student acceptance. **Chapter 1 remains NOT RUN.** No student release, public-guide update, publication, commit, or next-chapter authorization is implied. Root owns acceptance of the fixture and any subsequent code review or chapter gates. The original blocked 80/100 result and every prior seal remain unchanged.

## Reviewed identities and carried prerequisites

The bounded task is `/tmp/cadence-ch10-prompt-fixture-review-task.txt`, 5783 bytes, SHA256 `e006773a2a3bd710b84df0935f9bd47aea5a081f4a02f3d9ad7ff7258fd95a33`. The governing review input is `docs/edition-3-notes/ch10-prompt-fixture-review-input.json`.

Before controls, the verifier checked all **31 pinned inputs**, **33 protected tracked files**, all **137 frozen student files**, both frozen student executables, the source archive in both pinned local locations, and all **80 historical Chapter 10 reference files**. It verified the blocked package's **146 members** and five earlier exact inventories containing **78, 22, 212, 86, and 59 members**. The repeated candidate map is identical to the previously fully read frozen source manifest. The student remained clean throughout.

Previously complete governing, historical, public-contract, student-source, and diagnosis reads were carried only on matching hashes. The changed checker and new root decision, fixture-candidate/check receipts, and diagnosis parent verification were read. The actual focused test was read directly. The root's compile-only, adapter race/vet/build, and whitespace command outputs were verified by their pinned identities and exit codes; they were **not replayed**.

Exact source, executable, and candidate identities are retained in `v2/verification-before.json`, `v2/verification-after.json`, and the final receipt. In particular:

- Checker before: 20693 bytes, SHA256 `c7a63c5dc97a859f8d96f48f508edb0daad597dadbfb9a17dcd657c9240b1932`.
- Checker after: 20691 bytes, SHA256 `c3f5bd4a560e67922a2410139ab60c8ffbbe491d22bca3eff94c1ff20b025ece`.
- Corrected grader: `/tmp/cadence-ch10-grader-prompt-fixed`, SHA256 `394a51635fe8538653fe9e49025c9c3817612c05b67106effd2ecbd445ab9fb2`.
- Frozen student CLI: SHA256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33`.

No executable or source archive is exported in the portable evidence.

## Exact source review

The complete before/after byte comparison proves precisely one replacement in `checkPersistenceCurrentConfig` in `internal/gradeedition2/ch10.go`:

```diff
- ... CURRENT_RUNTIME_PROMPT\n"), 0600)
+ ... CURRENT_RUNTIME_PROMPT"), 0600)
```

Two source characters—the `\n` escape—were removed, removing **one runtime byte, hexadecimal `0a`**, from the primary file's final body. The complete current source equals the pinned previous source after this single replacement. The exact Git diff also equals the root's pinned candidate diff: 780 bytes, SHA256 `5b4c4cba86a5634dd1a8d4b0147371dbe122028acf8e72fe72f30c7408a24cac`.

Everything else in the checker remains byte-identical, including:

- the exact model check for `gpt-5-course` and exact system-text comparison with `CURRENT_RUNTIME_PROMPT`;
- the active tool declarations and the expected three tool names;
- the `CURRENT_WORD` request and restored-history checks;
- the recorded-config/current-wire comparison;
- source-level failure behavior and error text;
- the six persistence checks, their weights, the corrected Chapter 9 parity composition, and all timing limits.

Protected adapter source/tests, Chapter 10 prompt, public guidance, and student source/binaries were not changed by this review. There is no new trimming, fallback comparison, broadened accepted string, ignored assertion, point redistribution, or changed route.

### Why the fixture change is supported

The sealed diagnosis established, using the actual frozen clients, that the old primary body produced system text `CURRENT_RUNTIME_PROMPT\n`—23 UTF-8 bytes—while the exact assertion expected `CURRENT_RUNTIME_PROMPT`—22 bytes. Both current models were correct, history loaded correctly, and rewritten config metadata matched the actual current request. There was no measured historical-config restoration in those probes. The one authorized contrast removed just terminal byte `0a` and made the exact predicate true.

This review carried and rehashed those measurements instead of replaying them. Its retained negative-control record contains the two false LF-terminated predicates and the true no-LF contrast, tied to the unchanged diagnosis inventory and observation hashes. The new fixture now supplies the exact effective primary body that its unchanged assertion expects.

The root's decision is a fixture-alignment policy, not a new public requirement that all applications trim primary instructions. This review does not independently broaden that policy. The diagnosis used an explicitly handcrafted snapshot-only seed, not a byte-for-byte replay of the original historical grade's save; that limitation and its setup/environment qualifications remain in the sealed diagnosis. The new full gate supplies the directly requested corrected-grader result.

## Authorized runtime control 1 — corrected complete Chapter 10

The argv is the original blocked gate's full argv with **only argv[0] changed** from `/tmp/cadence-ch10-grader` to `/tmp/cadence-ch10-grader-prompt-fixed`. It retains:

- the same frozen CLI target;
- the same student core and GUI directories;
- the declared `examples/ch04-consumer` consumer and case;
- the declared `examples/ch06` workflow and pinned media-model file;
- the declared `examples/ch07` stream exercise;
- the same `claude-sonnet-4-6` stream and GUI model declarations;
- Chapter 10 and JSON-report flags.

Complete argv, cwd, environment, PID, deadline, duration, executable identity, and stdout/stderr identities are in `v2/command-receipts.json`. Core/GUI builds performed internally by the unchanged grader are normal parts of the assigned full gate, not additional review controls or student source edits.

Observed process **PID 92644**, deadline **600 seconds**, duration **18.974507 seconds**, **exit 0**, no timeout, empty stderr. The entire JSON report was parsed; exit zero alone was not accepted.

Report title: **`Edition 2, Chapter 10 — Persistence`**.

| Check ID | Earned / original weight | Passed |
| --- | ---: | --- |
| `save-shape` | 15 / 15 | true |
| `default-load` | 20 / 20 | true |
| `replay-equals-snapshot` | 20 / 20 | true |
| `tail-applied-once` | 15 / 15 | true |
| `log-not-needed` | 10 / 10 | true |
| `bad-save-refused` | 5 / 5 | true |
| `ch9-parity` | 15 / 15 | true |

**Total: 100/100; report `passed: true`.** Every expected ID appeared in the expected order with unchanged weight and full earned points. The original report remains separately preserved as **80/100, `passed: false`, exit 1**, with `default-load` false. It was not overwritten or relabeled.

The full gate exercises its normal CLI and GUI paths and composed lower checks. This result is evidence for the existing gate, not a claim of browser testing, untested save validation, whole-code completeness, or future chapter behavior.

## Authorized runtime control 2 — existing historical compatibility

The task text requested `TestChapter10CLISavesAndLoads`, but the pinned 27-line `internal/gradeedition2/ch10_test.go` declares only **`TestChapter10Persistence`**. Its unchanged implementation resolves `solutions/ch10`, calls `grade.Build`, invokes `persistenceChecks` for the real CLI and GUI, and fails on any returned check failure.

Root authorization **main-31** corrected the selector before execution and required actual run/pass markers. No empty selector was run or counted. Only after the full student gate passed, the reviewer ran exactly:

```text
go test -count=1 -run ^TestChapter10Persistence$ -v -timeout=180s ./internal/gradeedition2
```

Observed process **PID 93066**, outer deadline **240 seconds**, duration **4.091818 seconds**, **exit 0**, no timeout, empty stderr. The saved output is:

```text
=== RUN   TestChapter10Persistence
--- PASS: TestChapter10Persistence (1.51s)
PASS
ok      github.com/waywardgeek/ensemble/internal/gradeedition2     1.516s
```

The assertion checked the actual `RUN` and `PASS` markers, package `PASS`, exit zero, no timeout, and absence of `no tests to run`. This is one existing focused test, not an internal suite, race run, all-module campaign, legacy full grade, or fresh sensitivity/identity/lifetime experiment. The historical source and test were not copied, patched, or changed.

## Environment, boundaries, and cleanup

Both controls ran sequentially from the main checkout with an allowlisted environment containing every required field:

```text
GOMAXPROCS=2
GOFLAGS=-mod=readonly -p=1
GOTOOLCHAIN=local
GOPROXY=off
GOSUMDB=off
```

PATH includes the installed Go and Delve locations. Their local identities are recorded. HOME and TMPDIR are new and reviewer-owned. Explicit existing Go build/module caches are used without deletion or shared cleanup. No provider credential, provider endpoint, primary-skill selector, or skill directory was inherited into the control processes. The unchanged harness supplies its own loopback fakes and dummy keys to client launches. No browser, paid provider, key access, child agent, or standalone Chapter 1 command was used. Disk checks remained above the 1 GiB floor.

Both top-level jobs were waited to exit. The non-exhaustive transient-executable sampler stopped and joined; six observed temporary executables are recorded by local identity only. Sampling is not claimed to enumerate every internal child executable.

At runtime completion, no owned process remained. Six `ensemble-agent-*` directories remained **inside this review's owned TMPDIR**. After recording ten local-only file identities, the reviewer removed exactly those six owned directories. The original completion receipt is retained, and the separate cleanup receipt records the before/after list. No shared cache, source directory, old evidence, or pre-existing runtime state was cleaned. Final verification confirms the owned TMPDIR is empty and no owned process remains.

## Failures and scope reconciliation

1. **Root task selector error, resolved before controls.** The task named a nonexistent historical test. Root main-31 authorized the actual pinned test and forbade an empty-selector pass. No test was run under the incorrect name.
2. **Reviewer preservation-verifier setup failure.** The initial runner treated `candidate.tracked_files` as a map instead of a list. It stopped before creating a workspace or launching either runtime control. Its source, traceback, failure/completion records, and empty executable-observation record remain intact.
3. **Related read-only schema lookup error.** A lookup guessed `seal.entries`; the actual schema has `reports` plus `evidence`, and the archive uses `local_path`/`evidence_path`. These were corrected only in a versioned verifier. The complete driver recovery diff is retained. A read-only preflight then verified all preserved inputs and old seals before either runtime command.
4. **No runtime retries or hidden assertion weakening.** The corrected full gate ran once; the exact historical test ran once. Both succeeded. No extra negative replay, whole-code correction, Chapter 1 run, or cleanup of shared paths followed.

These setup errors are evidence-processing failures, not student or grader defects. They are not silently rewritten as an uninterrupted first attempt.

## Evidence handoff

The new evidence directory is `docs/edition-3-notes/ch10-prompt-fixture-review-evidence/`; the closed inventory includes this external report and every portable evidence file, excluding only the inventory itself. It retains:

- the initial zero-control failure and original driver;
- the versioned recovery diff and successful read-only preflight;
- prerequisite, selector, source-diff, old-negative-carry, and runtime plans;
- complete per-command receipts and stdout/stderr;
- the parsed Chapter 10 report, exact report verdict, and real-test verdict;
- source/candidate/archive/binary/reference/declaration verification before, after, and at finalization;
- non-exhaustive local executable identities, owned-cleanup facts, output identities, and final receipt.

All portable members are regular UTF-8/NUL-free files with path, byte count, and SHA256. No binary, archive, settings-file payload, raw API dump, or pre-existing package is exported or modified. Earlier reports/seals and root receipts remain separate and immutable.

**Handoff:** the evidence supports root acceptance of this narrowly scoped fixture correction and records a new green corrected Chapter 10 gate on the unchanged candidate. It does not itself authorize student continuation or supply the still-separate independent code review, Chapter 1 gate, or chapter closure decision. No automatic follow-on work is scheduled.
