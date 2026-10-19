# Chapter 10 entry follow-up — environment reconciliation

## Result and scope

**The one root-authorized, environment-corrected full exported Chapter 10 API run passed: 100/100, `passed: true`, all seven expected checks green, helper exit 0.** This supplies the full-composition positive control missing from the original entry audit. It does not implement or authorize an Edition 3 Chapter 10 route, grade the accepted student as a Chapter 10 submission, establish browser behavior, or close Chapter 10.

Authority: `/tmp/cadence-ch10-entry-followup-task.txt`; its exact size/hash is in `ch10-entry-followup-evidence/retained-identities.json`. This was a narrow continuation of the original audit, not another broad investigation. No sentinel, focused sensitivity, lifetime, browser, provider-service, or other behavioral control was repeated. No source, adapter, student, manuscript, grader, guide, or original sealed evidence was edited. No child agent, install, network module fetch, real credential, private journal, port 8088 inspection, or shared cleanup was used.

The original report and evidence remain immutable, including the 85/100 full report and both `no tests to run` mistakes. This follow-up must be read alongside them; it does not retroactively relabel those outcomes.

## Preservation and carried-forward prerequisites

Before execution, all 78 originally inventoried files and the closed original path set were reverified. Original inventory: 16,859 bytes, SHA-256 `1000f80deb8d88f6efd9c00a2c1c4d4b242b46df36b00c8f2409d33acc3dd26e`. Original report: 23,227 bytes, SHA-256 `8fee621cd8d8da979425450e12de930e8e912cff68333befc4adc6a749db2baa`.

The original full prerequisite reading carries forward because every one of the 18 input hashes still matched the governing manifest. The 29 protected tracked files, 80 historical reference files, and 132 accepted predecessor files also matched before and after this run. The accepted predecessor remains clean at commit `251ce88fb1677897dbcc6b33661096c5cb39cd8e`, tree `a2add41c306d567c05d98ae760911bc13c18b75a`. Main/student HEAD, tree, staged-diff digest, and unstaged-diff digest were unchanged across the follow-up. No unrelated/private diff contents were exported.

`verification-before.json` and `verification-after.json` contain exact source/preservation identities and the original seal's file pins. The retained full-API executable, its original sealed source, and the exact public consumer-case file were checked before invocation. **The helper was not rebuilt.**

## Exact deliberate environment correction

The original allowlist is preserved verbatim as `environment-v1-original.json`; the new version is `environment-v2-corrected.json`. A key-set assertion permits exactly these three changes:

| Variable | Original | Corrected |
|---|---|---|
| `PATH` | `/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin` | `/Users/bill/go/bin:/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin` |
| `GOMAXPROCS` | absent | `2` |
| `GOFLAGS` | `-mod=readonly` | `-mod=readonly -p=1` |

All other values are identical, including the original owned HOME, TMPDIR, and GOPATH, explicit existing caches, offline `GOPROXY=off`/`GOSUMDB=off`, local toolchain, `GOWORK=off`, and loopback proxy bypass. No shell credential environment was merged in. Fake-provider settings added internally by the unchanged grader retain their existing dummy/loopback behavior.

Installed Delve was pinned, not installed or searched for: `/Users/bill/go/bin/dlv`, 20,038,402 bytes, SHA-256 `10dbb1ce8e6d2cb77ae4d30c7ca1db0302b83529cd6cd39453f009a7653c17fe`. In the **same corrected allowlist**, explicit command resolution returned that exact Delve path and `/usr/local/go/bin/go`. `dlv version` exited 0 and reported version `1.27.2`, build `360e7b2181d3da54115aaee4f569954bca95cfe3`.

The original full run's missing-Delve failure is therefore classified as an audit setup limitation. The follow-up changed three variables together as authorized, so this is not a one-variable causal experiment about performance or resource limits. It is a correctly provisioned full positive control, with the original failure retained.

## One actual full-composition run

Retained executable: `/tmp/cadence-ch10-entry-audit-tr4jypr8/bin/full-api`, 10,882,018 bytes, SHA-256 `ebb48e170df896688142f58ffc4dadebcfb9f744fd9370f4d899d93ecc8c2655`.

Exact argv:

```text
/tmp/cadence-ch10-entry-audit-tr4jypr8/bin/full-api
/Users/bill/projects/ensemble/solutions/ch10
/Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json
```

Cwd: `/Users/bill/projects/ensemble`. PID: `83502`. Start: Unix `1792364221.234524`. Duration: **17.361 seconds**. Deadline: **600 seconds**, not an observation timeout. Exit: **0**. No timeout or cancellation. There was **one** full API attempt, no tuning retry.

The exact existing helper selects the same historical core, GUI `solutions/ch10/gui`, consumer `solutions/ch10/ch05`, workflow `solutions/ch10/ch06`, stream `solutions/ch10/ch07`, public case, and model declarations as before: `claude-sonnet-4-6` for GUI/streaming; media map Anthropic `claude-sonnet-4-6`, OpenAI `gpt-5.4-mini`, Gemini `gemini-3.8-flash`. The public case remains 111 bytes, SHA-256 `39bae58b556760c55e120f890564a39e3a2b94011d50ba14df989c06c1215071`. The accepted case is a declaration input, not a launch or grade of the accepted student implementation.

The helper's zero exit is insufficient by itself because it also exits zero for a generated failed report. The runner parsed stdout and required all of: exact title `Edition 2, Chapter 10 — Persistence`, `passed` literally true, score and max score both 100, the seven exact expected IDs in order, every check passed, helper exit 0, and no timeout. All requirements were met:

| Check | Earned |
|---|---:|
| `save-shape` | 15/15 |
| `default-load` | 20/20 |
| `replay-equals-snapshot` | 20/20 |
| `tail-applied-once` | 15/15 |
| `log-not-needed` | 10/10 |
| `bad-save-refused` | 5/5 |
| `ch9-parity` | 15/15 |
| **Total** | **100/100** |

The prior chapter parity check's actual title is `All corrected Chapter 9 checks pass`. Thus the full inherited composition is green under the corrected environment; the original six persistence checks remain green as well. This is an unchanged historical grader/reference positive control using real client processes and loopback fake providers, not production inference or a live GUI/browser acceptance claim.

## Standalone factual reconciliation of the original report

1. **Accepted example paths.** The original report invented `examples/ch05-workflow` and `examples/ch06-stream`. Those paths do not exist in the verified accepted predecessor. The actual public example paths are **`examples/ch04-consumer`**, **`examples/ch06` (workflow)**, and **`examples/ch07` (streaming)**. These are distinct from the historical reference's `ch05`, `ch06`, and `ch07` modules. Future public-interface planning must use the actual accepted declarations and module layout, not the incorrect prose.

2. **Already-colliding legacy archives are required, not optional.** The original phrase “Legacy duplicate-ID replay needs an explicit strategy **if supported**” is retracted as too weak. Current `book/edition-3/chapter-10.md:79–81` requires that “Already-colliding archived IDs must also produce distinct presentations, deterministically across replay and restart, without rewriting archive facts. Preventing new collisions alone does not repair that case.” Candidate acceptance therefore requires recovery/rendering of already-colliding legacy archives while preserving durable-history facts and provider bytes. Avoiding future allocator collisions alone is insufficient. This is a public chapter/source/integration acceptance obligation, not permission to rewrite history or provider request data and not a proposal to invent extra historical score conditions.

3. **Original resource deviation.** The original audit allowlist omitted task-requested `GOMAXPROCS=2` and build `-p=1`; it used only `GOFLAGS=-mod=readonly`. That was an audit execution deviation, now explicitly disclosed rather than erased. The follow-up supplies both limits for its one full run. It does not retroactively claim that earlier controls used them and does not repeat focused tests merely to change their resource environment. Original valid sensitivity evidence remains the disclosed exact-name green → red → restored-green sequence; the two empty selectors remain invalid.

4. **Source observation versus live behavior.** The historical/current-runtime source observations in the original package remain source-only where identified as such. The separately sampled old-reference save artifacts remain limited to what their bytes demonstrate. Neither becomes a live GUI/browser finding through this follow-up. No new identity investigation was performed.

5. **No speculative historical score changes.** Any original suggestion to add contract/identity “safeguards” must not be interpreted as authorization to alter historical scoring. Existing public candidate acceptance still owns source architecture, integration, deterministic presentation recovery, and GUI behavior. This follow-up preserves the existing historical score/check definitions exactly.

`failure-ledger.json` preserves these qualifications and identifies the immutable original failed and empty-selection receipts. Nothing in the new 100/100 result removes those facts.

## Evidence, ownership, and completion

New outputs are confined to this report and `docs/edition-3-notes/ch10-entry-followup-evidence/`. `run-followup.py` is the exact one-shot driver. Receipts include command argv, cwd, complete explicit parent environment reference, PID, timing, deadline, exit, timeout flag, executable identity, and stdout/stderr. `parsed-full-report.json` and `run-result.json` provide explicit semantic acceptance, not inferred success. The command-resolution and version preflights have their own receipts.

Four transient `grade.Build` executable identities and observed process command lines were sampled at 100 ms. These observations are **non-exhaustive**; they are not a claim that every short-lived child/binary/environment was independently captured. No binary is exported. Original retained scratch is preserved; no helper was rebuilt. The unchanged grader built and cleaned its normal per-run submission binaries. All commanded jobs were waited, and no owned process remained after the full run. No shared cache or unrelated file was removed.

The 1 GiB disk floor was checked before every top-level invocation and at final verification; recorded free space after the full invocation was 2,734,972,928 bytes. All work was sequential. The follow-up produced no file named `settings.json`, copied no private journal/diff body, and made no implementation change.

The new `inventory.json` inventories this report and every other new evidence file, by repository-relative path, byte count, and SHA-256. It excludes only itself under its explicit non-recursive self-inventory rule. All portable files must be UTF-8, NUL-free regular files; symlinks, binaries, and `settings.json` are excluded. The original inventory is only referenced and hash-pinned, never resealed.

**Conclusion:** the historical full positive-control gap is now closed by one environment-corrected 100/100 run. The factual corrections above supersede the identified prose mistakes without mutating the original package. Adapter edits, student launch/progression, historical-score expansion, and further work remain unauthorized by this follow-up.
