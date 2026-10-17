# Edition 3 Chapter 4 — focused C1 re-review, code correction Round 1

## Verdict: ACCEPT

**C1 is resolved on the frozen Round 1 candidate. No further code correction is requested.** The original external-consumer counterexample now passes. Reverting only the production fix in scratch reproduces both that failure and all three omitted-field regressions; restoring the exact corrected source restores passing results. Full/race/vet and both public graders independently pass on the corrected candidate, including the same executable’s real 62-second header/body phases.

This closes the initial review’s one required code correction. It is not another correction round, an adapter re-review, a manuscript/text-review approval, or an acceptance tag. The parent owns acceptance/tagging and the subsequent author procedure.

## 1. Scope and exact identities

- Corrected student commit: `ec4a350f564e820d398f9125fd73a9b3b27129cd`.
- Tree: `932a5fb7305526276a1c98b53fac9b5a92e55d16`.
- Parent/initial candidate: `062c146cafa9b105745f6b44921160198bbaee32`.
- Student worktree clean before and after review.
- Main HEAD remains `bac07e644facb25ba3e5bffb2961674763fed4c4`; intentional adapter/guidance/evidence and unrelated work preserved.
- `/tmp/ch04-r1-core`: SHA256 `7342ee1d5e91791c23bc3ce32a73db51a6ed157a8d02ce07c2d3296fdf028cd3`.
- `/tmp/ch04-r1-consumer`: SHA256 `7e032ea5e2692838e2a2a612acc929611e644105809d6b2bc6c6a7856fac95f7`.
- `/tmp/ch04-r1-candidate-source.tar.gz`: SHA256 `e66358c5ecfc1a63701bd387da1cae4147b8ca95cd4886c21f34e0c7452add8c`.
- Student R1 report: `docs/edition-3-notes/ch04-student-review-r1.md`, SHA256 `32ffb2bd1b815d220529ec5254129217a141794322b9cb2871dc6315630f8c5a`.

Both executable build records identify the corrected commit and `vcs.modified=false`. I verified their hashes before/after execution and at closure. The archive’s complete file set and every file’s bytes equal the committed tree: 35 files. I did **not** rebuild the two release artifacts to claim byte-identical rebuild reproducibility. Test compilation from the verified archive and execution of the pinned artifacts are separate, explicitly identified controls.

The diff is exactly two files: an eight-line production diff in `internal/agent/dispatch.go` (seven insertions, one deletion), plus the new 117-line `chapter4_limits_test.go`. No other implementation, job/process/engine/history code, public surface, consumer module, or CLI changed.

Verified 50 R1 evidence/report records, 35 source records, 16 local artifacts, and all 63 original evidence records. The original evidence remains immutable. Both separate scratch source trees match the corrected candidate at closure, including restoration of the sensitivity tree.

All governing prerequisites and accepted adapter hashes match the complete readings from the initial review; no changed prerequisite was skipped. The expanded public execution chronology and the complete R1 report were read. Initial whole-contract review is carried from `ch04-code-review.md`, SHA256 `6e89613a28ea33ea61f31d900adca83825ad793e829a57869a40f136a85cd89a`, not replaced by a new speculative review phase.

## 2. Why C1 is closed

Chapter 4’s single-use rule (`book/edition-3/chapter-04.md:389,397–406`) explicitly consumes a pending setting even when the next call is another `tool_limits`. The next setting must not inherit omitted fields from that spent setting.

The existing consumption/annotation logic in `internal/agent/dispatch.go:61–78` remains unchanged. At `:164–172`, the correction creates the next pending setting using:

```go
next, err := limits(raw, defaults())
```

and stores `next`, rather than storing the present call’s merged observation options. This directly separates the consumed setting from the newly installed setting. It reuses the existing validation helper and defaults; no new framework, scheduler, API, config, or target convention was introduced. Normal failure handling and loud consumption annotations remain in place.

The new tests cover:

- **Omitted cap:** consume a 1,024-byte setting with a delay-only replacement, then inspect the default-sized truncated read and complete 40,960-byte persisted output.
- **Omitted delay:** consume a 1 ms setting with a cap-only replacement, then wait for an actual delayed registered job to finish rather than returning prematurely.
- **Omitted pattern:** consume `MATCH` with a cap-only replacement; a real job writes `MATCH` before delayed completion, so a leaked pattern produces an early return.
- Explicit call-level cap precedence over the new pending cap; invalid replacement consumption/reset; separate-Agent isolation; and consumption/reset around a failing registered job.

These inspect actual reports, jobs, persisted bytes, and registration/dispatch rather than private pending-state fields. Their timing cases are finite 150 ms fixtures, not a new product latency requirement or proof of every scheduler interleaving.

## 3. Independent counterexample and mutation sensitivity

I reused the exact hashed external public-API probe from the initial review, in an independent module importing only public `ensemble` and replacing it with the frozen scratch source. The original sequence remains unchanged:

1. `tool_limits({max_output_bytes:1024})`;
2. `tool_limits({ai_callback_delay:0.3})`;
3. `read_file` of a 40 KB file.

It now passes, along with the existing ordinary reset/explicit-override/failure/isolation control. This is not a newly weakened test.

| Independent control | Exit | Observation |
|---|---:|---|
| Public counterexample + surrounding controls, corrected source, race | 0 | Pass |
| New chained-limits regressions, corrected source, race | 0 | Pass |
| New regressions, production-only revert | 1 | All three omitted-field subtests fail behaviorally |
| Original public counterexample, production-only revert | 1 | Old 1,024-byte cap survives; report is 1,080 bytes |
| New regressions, exact restored source, race | 0 | Pass |
| Public counterexample + surrounding controls, exact restored source, race | 0 | Pass |

The mutant uses the exact initial candidate’s `internal/agent/dispatch.go`; the corrected new tests remain unchanged. The failures are **not compilation failures**. The omitted-cap failure reports `report bytes=1080`; delay and pattern failures return running jobs with leaked `1ms` and `MATCH` settings. Both scratch trees and the live candidate were verified unchanged/restored afterward.

Receipts: `sensitivity-runs.json`, `sensitivity-inputs.json`, `regressions-reverted.log`, `public-reverted.log`, and the corrected/restored logs. The reused probe is stored as `.go.txt` documentation, not as an accidental Go package in the main repository.

One reviewer harness error is preserved: after recording the three genuine regression failures, the first runner expected different assertion wording and stopped. Its `finally` restored the source, verified by hash. The corrected runner matched the actual messages and resumed only unfinished controls; completed controls were not needlessly rerun. See `reviewer-setup-failure.txt` and `run-sensitivity-attempt1.py.txt`.

## 4. Independently repeated baseline gates

All five commands finished naturally within adequate outer observation bounds; no shortened long test or interrupted command is counted as a pass. Exact argv/cwd/duration/exit/output hashes are in `baseline-runs.json`.

| Control | Result | Elapsed |
|---|---|---:|
| Verified archive: `go test -count=1 ./...` | exit 0 | 125.785 s |
| Verified archive: `go test -race -count=1 ./...` | exit 0 | 127.149 s |
| Verified archive: `go vet ./...` | exit 0 | 0.165 s |
| Corrected pinned executable: public E3 Chapter 4 | exit 0, `passed:true`, 160/160 | 8.549 s |
| **Same corrected executable:** public E3 Chapter 1 | exit 0, `passed:true`, 100/100 | 124.115 s |

Both public invocations ran from `/Users/bill/projects/ensemble`:

```sh
go run ./cmd/grade-edition3 -ch 4 -json \
  -core /Users/bill/projects/ensemble-edition-3-student \
  -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer \
  -consumer-case /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json \
  /tmp/ch04-r1-core

go run ./cmd/grade-edition3 -ch 1 -json /tmp/ch04-r1-core
```

The new Chapter 1 report explicitly records complete answer/usage after **1m2s headers delay** and **1m2s body delay**. Its zero-point lifetime checks and final `passed` value were inspected, not inferred from the 100 points. These are fresh R1 measurements, not reused initial-candidate results. The Chapter 4 route builds its temporary consumer from the pinned module; that is distinct from the separately pinned consumer artifact audited below.

The inherited whole-contract/source assessment remains applicable: the correction changes only installation of new pending limits. The full/race suite and public gates were repeated, but unchanged process-group, late-result, PTY, cursor, logger, projection, explanation-policy, and shutdown code was not subjected to a new optional redesign or exhaustive fault-injection campaign. Prior independent contract controls and their stated limits remain identified in the initial report.

## 5. Audited retained local CLI and consumer boundary

I read the complete loopback fixture and independently checked the saved receipt claims against the actual neutral JSONL events, stdout/stderr, and persistent output files. `audit-retained-local.py` and `retained-local-audit.json` record this verification; they do not launch an application or provider session.

For **each** exact R1 binary, the retained fixture has four local HTTP requests and nine correlated actual tool calls/returns. It exercises three batches: partial-settings reset, invalid replacement consumption, and explicit call override. The consumer route uses the public registered `shout` handler for ordinary jobs, not just a build/import smoke.

| Saved run | Default-reset returned bytes | After invalid replacement | Explicit override | Full job files |
|---|---:|---:|---:|---|
| Core CLI | 16,441 | 16,320 | 8,247 | 40,960 / 40,960 / 40,974 bytes |
| Registered consumer | 16,437 | 16,320 | 8,247 | 40,960 / 40,960 / 40,974 bytes |

I verified consumption annotations, the invalid control’s error classification, explicit 8,192-byte call-level cap over the pending 2,048-byte setting, done metadata, exact output prefixes, full persisted bytes/hashes, and two timestamped consumer diagnostics. The fixture’s local placeholder model/key and loopback endpoint are not paid/vendor evidence or a credential-store read. These are **audited retained local runs**, not falsely labeled independently rerun standalone sessions. Fresh independent public-API tests and public grader executions are separately recorded above.

The student’s earlier setup failures remain visible and excluded from passing claims: initial test compilation/API-shape errors; the missing-system test fixture; the wrong historical grader path after completed full/race/vet, which exited 127; and the first JSONL parser treating the version header as an event. The corrected parser explicitly validates and skips that header. The later correct public routes ran on unchanged pinned bytes and completed successfully. I inspected the raw failures and corrected fixture/results rather than copying the report’s summary alone.

## 6. Carried live evidence and limits

No paid calls, credential/settings reads, live GUI, or unrelated process control occurred in this re-review. Independent commands stripped provider/secret environment variables and disabled dependency network fetching. No new vendor session is claimed on R1.

The original real-vendor evidence remains tied to commit `062c146…` and the original artifact `/tmp/ch04-core` SHA256 `b0ed87f60783977c2024547690b4865ff6269da0f34a51c484f39901eb003a57`, not the corrected artifact. Its initial Anthropic failure, successful bounded follow-up, OpenAI refusal-only result, Gemini/consumer successes, Delve consumed-return reconciliation, saved server shutdown observation, and disclosed assistance/recovery are preserved with the exact limits stated in the initial review. The eight-second test remains eight seconds, not a 90-second measurement; model prose is not upgraded into process/side-effect evidence.

The narrow settings-state fix is accepted using source review, independent public-API and mutation controls, fresh full/race/vet/public graders, and audited local CLI/consumer execution. No hard coverage threshold, new live-call requirement, or speculative API obligation is invented. No byte-identical rebuild, new platform campaign, or universal race/lifetime proof is claimed.

## 7. Preservation and handoff

Candidate source, student report/evidence, manuscript, guides, accepted adapter, historical graders/references, and previous review reports were not edited. The only production mutation was in the reviewer’s independent scratch sensitivity tree and was restored byte-for-byte. No child workers were used. No staging, commit, tag, install, restart, or push occurred.

All pinned artifacts and 11 command-output hashes verify at closure; both scratch source trees match the corrected candidate, and the live worktree remains clean. Current prerequisites/adapter and original 63 evidence records are unchanged. This review’s own artifacts are text-only under `ch04-code-review-r1-evidence/`; binaries and executable tests are outside the main repository.

**Recommendation: accept the Round 1 code candidate and close C1.** No additional code correction or paid call is requested by this review. The parent owns final verification, saved-tail/join closure, acceptance/tagging, and transition to the existing author procedure. The concrete teaching feedback from the initial review is now embodied in the chained-partial-settings regressions; no additional manuscript feedback or text-review verdict is introduced here.
