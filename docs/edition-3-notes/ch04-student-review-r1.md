# Edition 3 Chapter 4 — formal code correction round 1

## Scope and candidate identity

Parent opened **formal code correction round 1 of at most 3**, with exactly one required correction, C1. The supplied public failure: `tool_limits({max_output_bytes:1024})`, then `tool_limits({ai_callback_delay:0.3})`, then a 40KB `read_file` inherited the spent1024 cap. Current Chapter4 lines389,397–406 require consumption by the very next call, including another `tool_limits`; the new setting must start from defaults plus only its new explicit fields.

Candidate repository: `/Users/bill/projects/ensemble-edition-3-student`.

- Corrected commit: `ec4a350f564e820d398f9125fd73a9b3b27129cd` — `Reset chained tool limits to fresh defaults`.
- Corrected tree: `932a5fb7305526276a1c98b53fac9b5a92e55d16`.
- Parent: initial student candidate `062c146cafa9b105745f6b44921160198bbaee32`; its tree remains `840b5b3bab4daf7490b5555451bb6506b273d7b9`.
- Corrected core: `/tmp/ch04-r1-core`, SHA256 `7342ee1d5e91791c23bc3ce32a73db51a6ed157a8d02ce07c2d3296fdf028cd3`.
- Corrected independent consumer: `/tmp/ch04-r1-consumer`, SHA256 `7e032ea5e2692838e2a2a612acc929611e644105809d6b2bc6c6a7856fac95f7`.
- Committed-source archive: `/tmp/ch04-r1-candidate-source.tar.gz`, SHA256 `e66358c5ecfc1a63701bd387da1cae4147b8ca95cd4886c21f34e0c7452add8c`.
- `candidate.patch` records the complete two-file diff. Binaries were built **after committing**, from the clean corrected commit. No source change after pinning is planned or permitted without a new explanation.

All round evidence is separate under `docs/edition-3-notes/ch04-student-r1-evidence/`. Initial report and evidence remain immutable. At round entry all63 entries of the initial evidence/report manifest verified. The preserved initial report hash is `010b71e4809b46545680ad24824279267894f9300241b965348bc230f747efd9`, and initial manifest hash is `530e7e8d3958b379e5075ff7ab159a1ab1c7f5213c3fd5e4af55667e3b86e0e9`. See `initial-evidence-verification.json`.

## Required reads, public assistance and bounded change

Reread the expressly authorized launch task, current public AGENTS, student/coder/execution guidance, complete coding rules, cited current chapter contract, and the current candidate's public implementation/tests. Parent supplied C1's public reproduction and contract, not repair code, and approved a narrow ownership/lifecycle-preserving plan. No private reviewer report/tests/evidence, grader implementation, historical solution, future chapter, or other worker history was read. No helper agent was launched. No manuscript, guide, grader, API, model, budget, provider policy, process lifetime, or job architecture changed.

The existing dispatch still consumes and clears the previous pending setting at the present call boundary, records its consumption annotation, applies ordinary per-call overrides, and preserves failed-call consumption. The **only production change** is in `internal/agent/dispatch.go`'s `tool_limits` control: calculate `next` with `limits(raw, defaults())`, propagate a validation error, and store `&next` rather than storing the present call's merged options. Ownership, explanation gate, wait/input/kill, cancellation/reaping, cursors and persisted output paths are unchanged.

## Own regressions and meaningful sensitivity proof

`chapter4_limits_test.go` adds public-Agent tests, reusing the accepted public test helpers. No internal pending-state inspection is used:

1. Omitted cap: old1024 cap, new delay0.3, actual40,960-byte file. New default-sized truncation is observed and the complete persisted file is checked.
2. Omitted delay: old1ms delay, new2048 cap, real registered handler that finishes after150ms. It must finish within the default observation instead of prematurely returning running.
3. Omitted pattern: old`MATCH` pattern, new2048 cap, registered handler writes `MATCH` then finishes after150ms. The spent pattern must not cause an early return.
4. Explicit call-level cap override after chained settings, invalid second-control consumption, per-Agent isolation, registered failed-call consumption, and absence of sticky state on the following operation.

A separate scratch tree was made from initial commit062c146 plus these **same new tests**, leaving only the production fix reverted. `raw/sensitivity-reverted.log` records actual compiled behavioral failures, exit1:

- omitted cap: **1080bytes**, matching the public failure;
- omitted delay: old1ms annotation and `job 1 running` before completion;
- omitted pattern: old`MATCH` annotation and `job 1 running` after the matching output.

No compiler failure is counted as mutation/sensitivity evidence. Copying the corrected `dispatch.go` into that same scratch tree restored all test passes, exit0 (`raw/sensitivity-restored.log`). The candidate-focused run also passes (`raw/focused-fixed.log`). `raw/sensitivity-status.log` preserves both actual exit codes; the scratch location is retained in `sensitivity-scratch-path.txt`.

### Intermediate failures retained honestly

The first scripted production anchor did not match and therefore made no production change. The first test draft used guessed API names and failed compilation. Its raw output is retained as `raw/initial-anchor-and-compile-failure.log` and `raw/focused-initial.log`. After correcting the test APIs from actual source, a second-Agent fixture omitted the required system prompt and failed while the three omitted-field regressions passed (`raw/focused-corrected.log`). The fixture was corrected before pinning. These are real failed preparation attempts, not successful regression detection or evidence against production. No long-running test was interrupted this round.

The first gate script used a nonexistent grader executable path after full/race/vet had passed. It stopped with exit127 before grading; `raw/final-sequence.log` and `raw/chapter4.log` retain that failure. The authorized public `go run ./cmd/grade-edition3` invocation was recovered from the permitted launch instructions and resumed separately with the same binaries. Initial archive snapshots were also placed in the evidence directory; this workflow mistake was corrected before freezing by moving them outside the main repository, unchanged, to `/tmp/ch04-r1-candidate-{source,binaries}.tar.gz`. `archive-relocation-receipt.json` records the paths. No production source or initial evidence changed.

## Final sequential gates on the corrected commit

`run-final-gates.sh` runs full, race, vet, public Chapter4, and public Chapter1 in that order, preserving exact commands, UTC boundaries and exit statuses in `raw/final-sequence.log`. It stops on any nonzero exit and checks `passed:true` plus every public check. Both grader invocations use the same `/tmp/ch04-r1-core` bytes. Chapter4's core source is the corrected module directory; consumer source and case are the separate `examples/ch04-consumer` module and its `consumer-case.json`.


| Gate | UTC start | UTC finish | Result |
|---|---|---|---|
| full | 2026-10-17T01:26:25Z | 2026-10-17T01:28:31Z | exit0; PASS; package ensemble125.534s |
| race | 2026-10-17T01:28:31Z | 2026-10-17T01:31:19Z | exit0; PASS; package ensemble126.807s |
| vet | 2026-10-17T01:31:19Z | 2026-10-17T01:31:19Z | exit0; PASS |
| chapter4 | 2026-10-17T01:31:19Z | 2026-10-17T01:31:19Z | exit127; Invocation failure only; no grading (nonexistent executable) |
| public-chapter4 | 2026-10-17T01:32:39Z | 2026-10-17T01:32:48Z | exit0; PASS; passed:true;160/160, every check passed |
| public-chapter1 | 2026-10-17T01:32:48Z | 2026-10-17T01:34:52Z | exit0; PASS; passed:true;100/100, every check passed |

The successful gates remained sequential. The initial gate-script exit127 stopped grading before any grader ran; `run-public-graders.sh` then used the exact authorized public entry point, in sequence, without rerunning the already-passed suites or changing the candidate. `raw/public-chapter4.json` and `raw/public-chapter1.json` each contain `passed:true`; their process exits are0 in `raw/public-graders-sequence.log`. The adapter's reused output title says Edition2; the actual invocation was the public Edition3 `grade-edition3` entry with `-ch 4` or `-ch 1`, as recorded.

Chapter1's `generation_headers` and `generation_body` checks both explicitly report complete answer **and usage** after an actual **1m2s** delay. No long test was interrupted, skipped or shortened. These observed62s probes do not alone exclude every longer arbitrary cutoff; independent source review remains necessary.

## Targeted local binary verification

`local-boundary-check.py` is a local-only HTTP boundary fixture. It runs the actual corrected chat CLI and actual independent consumer executables with explanation policy enabled, using loopback Anthropic-shaped responses and an explicit local-fixture model identifier. It does not replace Agent dispatch, tools, jobs, file reads, shell execution, registered consumer handlers, or root logging.

Executed checks: chained partial cap settings, invalid replacement consumption, explicit per-call cap override; exact returned tool results and annotations; actual40KB persisted outputs; done job metadata; independent consumer registration retaining builtins and logging through the real root logger. The consumer's large synthetic `shout` input is deliberate local output-cap stress, not evidence that a live model emitted that call within its real token budget. Ordinary `run_command` uses the actual existing shell/PTY runner, but this is not a rerun of the original interactive debugger/server campaign.


| Actual executable | Process exit / assertions | Three observed result sizes | Ordinary job statuses |
|---|---|---|---|
| Corrected core | exit0 / all passed | 16441, 16320, 8247 bytes | done, done, done |
| Corrected consumer | exit0 / all passed | 16437, 16320, 8247 bytes | done, done, done |

The three batches respectively cover default cap reset after partial replacement, failed replacement consumption with default cap restored, and explicit8192 call-level cap overriding a new2048 setting. Each final fixture session made four loopback requests and finished with `LOCAL C1 CHECKS PASSED`. All three actual ordinary outputs were retained on disk, at least40,960bytes each. Consumer `shout` executed twice and emitted two real root-logger diagnostics; its builtins remained available for the shell case. Per-application receipts, output, neutral logs and binary/source identities are under `local-core/` and `local-consumer/`. Their temporary workspaces and job outputs are outside the main repository and indexed separately using the actual receipt paths. Python allocated them under the host temporary directory, not necessarily `/tmp`. The first incomplete attempt lacks a receipt; its directory was identified by this round's unique own fixture prefix in that host temp directory.

The first local attempt reached the core's final response but its evidence parser failed on the JSONL version header (which has no `type`). The failed script, traceback and core stdout/stderr/log are retained as `raw/local-boundary-attempt1.py`, `raw/local-boundary-attempt1.log`, and `local-core-attempt1/`. I initially misdiagnosed this as field capitalization; inspecting actual records showed lowercase event fields following `{"log_version":1}`. Only the fixture was corrected to validate and skip that header. The successful final local checks use unchanged source and binaries. No successful consumer run is claimed for the incomplete first attempt. An artifact-index query also incorrectly assumed every shell-job record included `cwd`; that field was omitted in this case and the query failed before writing its manifest. Its raw output is retained as `raw/artifact-index-cwd-assumption-failure.log`; indexing was then completed from actual receipt paths and the unique own temporary prefix.

These are **local synthetic provider-boundary checks**, not live vendor identity evidence, actual token/cost observations, or new vendor acceptance. No production behavior was mocked below the external HTTP response boundary.

## Live-evidence attribution and exclusions

**No fresh paid calls are authorized, requested or made in this correction round.** The configured credential store is not needed or accessed for this local correction. The local fixture child environment uses only an unmistakable non-credential placeholder and a loopback URL; no inherited provider credentials are passed.

All previously recorded actual Anthropic, OpenAI, Gemini3.8, consumer/root-logger, Delve, useful Go-test/server, and EOF/OS-process observations remain attributed to initial commit062c146's tree and original binary hashes. They are retained in the immutable initial report/evidence; **none is claimed as newly run on ec4a350**. The current chapter requires Delve, not GDB. No new vendor campaign, debugger campaign, timeout, model change, budget change, fallback or installation occurred. Independent source rereview plus the corrected local checks are required to assess this small change against the existing historical live observations.

## Final status

All required gates and the targeted local checks are complete as recorded above. The corrected source remains exactly commit`ec4a350f564e820d398f9125fd73a9b3b27129cd`, tree`932a5fb7305526276a1c98b53fac9b5a92e55d16`, with a clean working tree. No source edits or rebuilds occurred after pinning. `source-manifest.json` covers all35 committed files; the source snapshot was verified byte-for-byte against them. `binary-manifest.json` covers the exact two binaries and their compressed snapshot outside the main repository. `local-artifacts.json` distinguishes the binary/source archives and disposable outputs from portable evidence. `evidence-manifest.json` inventories this report and the round's portable evidence; raw failures are not rewritten. Initial report/evidence hashes are reverified unchanged.

Teaching feedback: explicitly separating the present call's consumed observation options from a newly created pending setting is the key C1 invariant. Independent omitted-cap, omitted-delay and omitted-pattern cases expose the same state-retention defect from different public behaviors; successful explicit-field cases alone would miss it.

This is a student correction candidate only, **not acceptance**. Parent owns focused independent source rereview, acceptance and any tag. No acceptance tag is created here.
