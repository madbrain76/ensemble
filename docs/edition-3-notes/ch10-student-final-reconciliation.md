# Chapter 10 final local reconciliation

**Status: both assigned final local grades passed on the unchanged frozen candidate. This is not whole-candidate acceptance; independent whole-code review and root verification/join remain pending. Browser and paid-provider observation remain NOT RUN.**

## Exact unchanged candidate and correction to the old predecessor text

Student repository: `/Users/bill/projects/ensemble-edition-3-student`.

| Identity | Verified value |
|---|---|
| Source commit | `729bc919e734405307b4ac0060e65c40e33652e5` |
| Source tree | `025ef3744b42c8f9daec34e76661edc40583ca6a` |
| Actual parent and `edition-3-ch09-r1` commit | `251ce88fb1677897dbcc6b33661096c5cb39cd8e` |
| Exact source membership | 137 tracked files, each size/SHA256 verified |
| Frozen CLI | 9,805,778 bytes; SHA256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33` |
| Frozen GUI | 10,645,906 bytes; SHA256 `fb506081f98c09d2e6da628edbb6cb922f067f9db24983ae856d2d1222dcbf53` |
| Existing exact source TAR | 952,320 bytes; SHA256 `8a4fc7d17b4371a06d3db91287115c1f2ac6027fa932d8ba89a2141018b50742` |
| Approved corrected black-box grader | 10,934,098 bytes; SHA256 `394a51635fe8538653fe9e49025c9c3817612c05b67106effd2ecbd445ab9fb2` |

The earlier review/submission transcribed the predecessor incorrectly as `251ce8848b2e64080f65ec23384447a2bd0fcd65c`. The actual Git parent and accepted tag both resolve to the value above. This new text corrects the error; the old report, source manifest, seal and other historical bytes were **not rewritten**.

No source edit, rebuild, new freeze, commit or tag occurred in this resumed task. Both grades used the same already-built CLI. Source, all 137 tracked files, both binaries, old archive and prior sealed artifacts were verified before and after the commands. Source remains clean. Original grader retained; only the separately approved corrected executable was invoked.

## Fixture-only correction and preserved history

Root's public authorization reports that the independently reviewed default-load fixture generated a prompt with a terminal newline while its unchanged exact expectation did not. The fixture input now matches that expectation. No student production defect was established by this mismatch, no trimming behavior was requested, and no student source correction was authorized or made. This report does not derive that diagnosis from private review or grader internals; it records root's supplied public outcome and the new black-box result.

The original development **65/100, exit 1** and frozen **80/100, exit 1** results remain failed historical attempts. They are not relabeled. The old blocked seal remains SHA256 `37eecf6fa8424a471ba80c98af75478f6353d7b9b1b4fc3e1fa3d52af19b95b7`; its three reports and 143 evidence members were reverified unchanged. The earlier expected-negative vendor diagnostic remains an exit-1 failed assertion; its numeric-vendor hypothesis was disproved, and no passing diagnostic was manufactured from it.

The existing real TAR under the **old main evidence directory** is retained unchanged and is **local-only, excluded from Git staging and from this new portable text inventory**. It is not text. Direct Git checks show it is untracked (`ls-files --error-unmatch` exit 1) and is not covered by an ignore rule (`check-ignore` exit 1). Thus this is an explicit submission/staging exclusion, not a claim that Git mechanically ignores it. No ignore rule or other main file was changed, and no new TAR or binary was copied into main. Hash/path references to local binaries/archives are retained instead. Root must continue excluding the old TAR when staging evidence.

## Actual remaining gate sequence and results

Only one corrected Chapter 10 invocation and the first Chapter 1 invocation for this candidate were run. Exact argv, cwd, complete clean environment, timings, actual subprocess exits, source/artifact before/after pins, stdout/stderr bytes and SHA256 are in the new receipt files.

### Corrected Chapter 10

- Executable: `/tmp/cadence-ch10-grader-prompt-fixed`.
- Every argument after the executable is identical to the prior frozen Chapter 10 invocation: same module/case/media declarations, same absolute frozen CLI.
- Actual grader exit **0**, `passed: true`, **100/100**, **18.987 seconds**.
- Report identity: `Edition 2, Chapter 10 — Persistence`. This is the grader's literal inherited title; Chapter 10 is identified by the title and `-ch 10` invocation, not an invented numeric JSON field.
- All seven checks passed: save shape; default load; snapshot/replay equivalence; tail once; context-only continuation; malformed refusal; corrected Chapter 9 parity.

The local receipt runner made one metadata error: it expected a numeric `chapter` field that this JSON report does not contain. It preserved actual grader exit 0 and `passed:true`, then returned outer exit 1 before launching Chapter 1. That original receipt and runner are retained unchanged. `chapter10-identity-reconciliation.json` explicitly records the mistaken wrapper assumption and verifies the actual title-based Chapter 10 identity. **Chapter 10 was not rerun.** Only the local receipt helper for the still-unrun Chapter 1 invocation was corrected to recognize the title. No student production source or assertion about grader success was weakened.

### First Chapter 1 on this same CLI

Exact argv:

```text
/tmp/cadence-ch10-grader-prompt-fixed -ch 1 -json /Users/bill/projects/ensemble-edition-3-student/.ch10-local/bin/frozen-cli
```

- Actual grader and runner exit **0**, `passed:true`, **100/100**, **124.142 seconds**.
- Identity: `Chapter 1 — Edition 3 conversation and request lifetime`.
- Protocol: **5 rounds, clean exit 0**.
- Wire/calls: **5 requests**, all well formed; exactly **one API call per round**.
- Replies matched the fake server's returned text.
- Request 4 contained **7 messages**, including the prior assistant marker; history grew **1 → 9 messages**, with earlier turns preserved byte-for-byte.
- Usage: cumulative **input=641, output=65**. Input per round: **82 → 111 → 128 → 149 → 171**; output per round from retained stderr: **20, 8, 10, 13, 14**.
- `generation_headers`: complete answer and usage after an actual **1m2s header delay**.
- `generation_body`: complete answer and usage after an actual **1m2s body delay**.
- No short outer deadline, cancellation, shortened timing, persistence-suppression wrapper or replayed predecessor evidence. The grader itself notes that longer arbitrary deadlines still require source review; these two phases do not prove unbounded lifetime behavior.

## Carried evidence, not fresh reruns

`carried-receipts.json` pins the original command index and all 31 indexed own receipts, distinguishing failures, unpinned/prefreeze history and source-pinned observations. The unchanged source/artifact identities permit carrying these prior results; they were not rerun here:

- Root tests and race passed (134.106s and 158.871s), and vet passed.
- GUI, consumer, workflow and stream test/race/vet commands passed. Consumer/stream packages had no test files; this is not inflated into runtime coverage.
- The compile-valid checkpoint-seeding mutation failed on the intended snapshot assertion; exact 137-file restoration then passed.
- Actual fake-HTTP CLI/GUI stdin/WebSocket process control v4 passed. Prior v1/v2/v3 harness failures remain failures. This is not browser evidence.

All earlier failures and disclosures remain in the old sealed artifacts: early compile/test failures; the 136.674s source-changing pre-freeze suite explicitly excluded from final-source evidence; incomplete early environment metadata; the disk-floor interruption; and both byte-preserving, root-authorized receipt/note relocations after initial writes outside permitted output paths.

The earlier recursive grep accidentally entered ignored `.ch09-local` predecessor-copy directories and printed test-line matches. That exposure was disclosed; root authorized continuation. No such path was reopened in this resumed task. Only explicitly owned prior evidence/source identities were inspected. This report preserves the boundary disclosure rather than asserting it never happened.

## Environment, ownership, assistance and remaining scope

Both new grader commands used a clean, explicit environment: `GOMAXPROCS=2`, `GOFLAGS=-mod=readonly -p=1`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`; required PATH including installed Go and `/Users/bill/go/bin` Delve; owned `.ch10-local/final-reconciliation/home` and `tmp`; explicit nonsecret existing Go cache/module-cache paths. The exact environment is recorded, with no inherited provider-key variables or generic `LLM_BASE_URL`. Free space was checked against the 1-GiB floor before each gate. No shared cleanup occurred.

Both actual grader subprocesses exited and their separately owned process groups were absent after waiting. No separately launched student fake service or child agent remains. These are scoped process observations, not an account-wide absence/security proof. The black-box grader owned its local fake fixtures; its private internals were not inspected.

The 5,519-byte final task was read in full and pinned at SHA256 `997d2bae0943465db4aa9937e99433b4fb241a2abdd6de4f67d82629a734e136`. Root supplied the public fixture-only outcome and exact resumed scope. No private review/audit/coordinator manifests were opened; root retains responsibility for the 33 private protection checks. Formal student production correction rounds used: **zero**. No source change was made to obtain these successes.

Browser/native/audio/MIC, paid/live-provider/key access, runtime-agent observation, children, installs, pushes and publication were not performed under this resumed authority. Whole-code independent review and any later browser/provider gates remain pending separate root authorization. Root must verify/join this submission before deciding broader acceptance.

## New package

`ch10-student-final-reconciliation-evidence/inventory.json` lists this report and every file in the new evidence directory except itself, with exact size/SHA256 and validated UTF-8/no-NUL content. It references, but does not copy or relabel, the old sealed artifacts and local binary/TAR identities. The inventory's own hash is supplied with the handoff. All old sealed bytes remain unchanged.
