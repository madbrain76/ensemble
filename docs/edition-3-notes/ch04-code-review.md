# Edition 3 Chapter 4 — initial independent student code review

## Verdict: REQUIRED CORRECTIONS

**One demonstrated correction: C1, pending settings leak across a consuming `tool_limits` call.** This is the initial code review, not an adapter re-review, student correction round, student acceptance, or manuscript/text-review pass. No student code was repaired by this reviewer. No author work is authorized by this verdict; the parent owns the next correction handoff and acceptance.

The public Chapter 4 and Chapter 1 gates both pass on the exact pinned executable. Full/race/vet also pass. Those results do not cover C1, which is independently reproduced through an external consumer of the submitted public API.

## 1. Frozen inputs and independence

- Main HEAD: `bac07e644facb25ba3e5bffb2961674763fed4c4`, with the intentional accepted adapter/guidance/evidence work preserved.
- Student commit: `062c146cafa9b105745f6b44921160198bbaee32`.
- Student tree: `840b5b3bab4daf7490b5555451bb6506b273d7b9`.
- Accepted parent: `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`.
- Student tree clean before and after; all 34 source-manifest entries verified. Scratch was a Git archive of this commit, and its source files also remained byte-identical.
- Core `/tmp/ch04-core`: SHA256 `b0ed87f60783977c2024547690b4865ff6269da0f34a51c484f39901eb003a57`.
- Retained live consumer `/tmp/ch04-consumer`: SHA256 `d6e4074774199f5cc2476b5fd648bf14cb60d6325e6c94d09158fd53adaf8dd2`.
- Student public report: SHA256 `010b71e4809b46545680ad24824279267894f9300241b965348bc230f747efd9`.
- All 63 retained evidence-manifest records and 48 local-artifact records verified, including the corrected report. Both binaries were unchanged before/after the public gates and at final verification.
- Accepted adapter `main.go`: `539ee3e6cc7313171ab57d268803ddd4c2d49be33fb028c38fa43b3b31c12834`; `main_test.go`: `c94e75cc895bfcf5e2625b523ff5571955c565ad8ef9f261b1f7ae3f7637fa2f`.
- Environment: Go 1.25.4, darwin/arm64; dependency network fetching disabled for independent commands. Provider/credential environment variables were stripped from these processes.

The binary build metadata honestly embeds predecessor revision `28be8bf…` with `vcs.modified=true`. It does **not** embed the final commit. The public report/manifests identify the staged tree used for those builds. This review verified and executed the exact retained core artifact; it does not claim a final-commit live request or byte-identical rebuild. The Chapter 4 grader builds its own temporary consumer from the pinned source; that is distinct from the retained live `/tmp/ch04-consumer` artifact.

All current required prerequisites were reloaded completely, recovering omitted portions: AGENTS, coding/reviewer skills, workflow, Edition 2 carryover, voice, writing procedure, current Edition 3 execution/student guides, readiness/carryover, and all of Chapter 4. Their current hashes are in `input-identities.json` and were rechecked unchanged at closure. Current student-first/three-correction/one-text-review procedure supersedes the historical author-first descriptions.

No private student conversation or parent runtime history was read. Evidence sources were the frozen source, public reports/receipts, current contract, and this reviewer’s earlier entry/adapter reports. No credentials/settings files were read, no paid calls or GUI sessions were made, and no unrelated process was killed. No candidate, manuscript, guide, historical grader/reference, or other report was edited. Main-repository writes are limited to this report and its evidence directory; binaries and test consumers remain in scratch.

## 2. Required correction C1 — rearming leaks consumed pending fields

**Contract:** `book/edition-3/chapter-04.md:389` defines `tool_limits` arguments as fields to set for the next call. Lines 397–406 make the pending setting single-use, consumed by the very next call, explicitly including **another `tool_limits`**. Defaults are 3 seconds/no pattern/16 KiB at lines 369–376. A consumed setting cannot silently become the next call’s pending setting again merely because the second control omits that field.

**Implementation:** `internal/agent/dispatch.go:61–78` consumes/clears `a.pending` but initializes `opts` from the consumed value and overlays the current arguments. Then `:164–166` handles `tool_limits` by storing that same merged `opts` back into `a.pending`. Thus omitted old fields survive into a second later call. The consumption annotation is present and truthful about which call consumed the old setting, but the reinstalled setting contradicts the reset rule.

**Independent public-API reproduction, on the unchanged candidate:**

1. Create an ordinary public Agent with tools enabled and write a 40 KB file containing distinguishable head/tail text.
2. `CallTool("tool_limits", {"max_output_bytes":1024})`.
3. `CallTool("tool_limits", {"ai_callback_delay":0.3})`.
4. `CallTool("read_file", {"path":file})`.

The second return says:

```
[tool_limits consumed by this tool_limits call: ai_callback_delay 3s,
 ai_callback_pattern none, max_output_bytes 1024]
tool_limits set: the next call, whichever tool that is, will report the consumption.
```

The following file result is only **1,080 bytes**, instead of receiving the default 16 KiB payload allowance. The external test fails:

```
old 1024-byte pending cap survived consumption by tool_limits;
following read=1080 bytes, want default 16KiB payload
```

Receipts: `pending-rearm-reproduction.log`, `public-api-runs.json`; test `TestPublicPendingRearm` in `public_api_test.go.txt` (and the exact initial successful-probe source retained as `public-api-initial-success-source.go.txt`). The command’s recorded cwd is an independent module importing only public `ensemble`, with a local replacement to the frozen scratch source. No internal-package import or implementation patch was used.

**Sensitivity/restored control:** supplying `"max_output_bytes":16384` explicitly on the second `tool_limits` makes the same following read pass (`pending-explicit-reset-control.log`). Ordinary one-shot reset, explicit `run_command` override, consumption by an unknown failing call, and separate-Agent pending isolation also pass in `positive-public-controls.log`. This is an input-level sensitivity control, not a claim that candidate source was repaired or that a source-fix mutant was tested.

**Required result:** separate the consumed observation settings for the present call from the new pending settings established by this `tool_limits` call. Omitted fields in the new pending setting must start from defaults, not the old consumed setting. Preserve the loud consumption annotation, explicit-field behavior, failure consumption, and per-Agent isolation. Add a focused regression covering chained partial `tool_limits` calls (including omitted delay/pattern/cap fields as appropriate). No new scheduler, CLI, config, or target mechanism is requested.

To reproduce from the retained scratch directory:

```sh
cd "$(cat /Users/bill/projects/ensemble/docs/edition-3-notes/ch04-code-review-evidence/scratch.txt)/public-consumer-r2"
GOPROXY=off GOTOOLCHAIN=local GOFLAGS=-mod=readonly \
  go test -count=1 -v -run '^TestPublicPendingRearm$' .
```

Expected on this initial candidate: exit 1. The explicit-default control exits 0. Scratch lifetime is not guaranteed; the public test source and exact command receipts are portable evidence for recreating the external module.

## 3. Independently measured automation

All five baseline commands finished naturally, with no timeout or interrupted test counted as a result. `baseline-runs.json` stores argv, cwd, elapsed time, exit code, report values, and output hash.

| Control | Result |
|---|---|
| Scratch `go test -count=1 ./...` | exit 0, 125.404 s |
| Scratch `go test -race -count=1 ./...` | exit 0, 126.793 s |
| Scratch `go vet ./...` | exit 0, 0.169 s |
| Exact pinned core, new E3 Chapter 4 route | exit 0, `passed:true`, 160/160, 8.604 s |
| **Same exact core**, E3 Chapter 1 route | exit 0, `passed:true`, 100/100, 124.118 s |

The Chapter 1 report separately confirms complete answer/usage after **1m2s response-header delay** and **1m2s response-body delay**. These are zero-point required checks; the verdict was read, not inferred from the score. Source review also found no newly introduced arbitrary provider-generation deadline. The source review is not a claim that finite tests prove every possible longer duration.

Actual public invocations, both from `/Users/bill/projects/ensemble`:

```sh
go run ./cmd/grade-edition3 -ch 4 -json \
  -core /Users/bill/projects/ensemble-edition-3-student \
  -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer \
  -consumer-case /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json \
  /tmp/ch04-core

go run ./cmd/grade-edition3 -ch 1 -json /tmp/ch04-core
```

These are new-adapter results on the actual frozen student artifact, not relabeled historical-wrapper/reference controls. The historical reference’s lifetime failure and root-library launch failures remain preserved in the earlier entry evidence; they were not needlessly repeated.

### Independent external public-API controls

The successful positive batch ran under the race detector and exercised:

- pending cap/head-tail output, ordinary default reset, explicit run-command cap override, failed unknown-call consumption, and per-Agent pending isolation;
- external registered handler running as a job, root logger reachability, called/returned job metadata, explicit kill, wake of an existing long waiter, and a matching `job_killed` terminal event;
- an intentionally uncancellable plugin released **after kill**, with terminal status remaining killed and no late bytes accepted into persistent output;
- separate-root handle/state isolation and no replay of already-consumed completed output;
- fresh command cwd/environment, actual terminal configuration, and an old callback pattern not waking on already-consumed bytes while later new bytes remain visible.

The initial canonical-mode large-input fixture passed but did **not** block its writer; its log honestly records `blocked=false`. A subsequent noncanonical/no-echo fixture did establish `blocked=true`, and kill interrupted it in 0.803 ms under race; the small-input comparison passed as well (`pty-noncanonical.log`). No rescue signal was required. Terminal name/dimension observations are not new mandated configuration values.

C1 is the only completed probe that demonstrated a candidate defect. A suspicion about a mutex being held during PTY input was dismissed after recovering the actual source (`Input` unlocks before writing) and testing it. A concern about lost tool-error classification was likewise dismissed: `Job.report` returns the stored failure on inline and framed paths, preserving dispatch `IsError`.

## 4. Whole-contract/source assessment and coverage limits

The root remains a library with a separate CLI. `internal/common` owns the shared contracts; the root composes Agent, history, engine, tools, and jobs. Feature leaves depend through common parent interfaces rather than constructing sibling owners or introducing a new parallel framework. Chapter 4’s public topology/consumer check passes against the actual source. The independent consumer imports only the public root package.

The central Agent dispatch creates registered-handler/ordinary jobs, installs called/returned job metadata, and routes the four supervision controls without creating jobs for them. Explanation/selected-model refusal remains before execution; inherited projection/ordering/error tests pass. Existing engine/projection behavior was not replaced with a jobs-specific transport. The provider request-lifetime checks above remain green.

Jobs own synchronized terminal state, an output file, cursor, cancellation, process attachment, and notification. Source and completed tests cover one-way kill state, late-result suppression, current-waiter notification, atomic attachment/cancellation handling, process-group stop/reap, and shutdown. The process path uses a real PTY, fresh shell, base directory/current-call override, fixed terminal settings, newline input, and explicit execution cancellation rather than imposing the observation delay as a lifetime.

The submitted `chapter4_test.go` tests and the full/race run add coverage for process attachment, descendants/shutdown races, caps/persistence, concurrent waits, and invalid/control consumption. Independent probes add the external-plugin/late-result/waiter/cursor paths above and reveal the chained-limits hole. No hard coverage percentage is invented. Race-test success is not proof of all possible schedules. A Go handler that ignores cancellation cannot be forcibly terminated; the contract requires terminal-state/late-result protection, which was exercised, not a fictional goroutine kill.

The historical Chapter 4 fixture’s 160 points and this finite probe set do not establish every conceivable scheduler, filesystem-failure, resource-exhaustion, operating-system, or duration boundary. No additional platform run, new 90-second background workload, comprehensive fault injection, or new paid provider session was performed. None is silently claimed. The contract, not a speculative future scheduler/API, defines the correction above.

## 5. Retained live evidence: actual credit and limits

I verified the manifest bytes and inspected sanitized neutral tool facts, actual returns, saved stdout/stderr, and the saved follow-up process observation. `retained-live-tool-facts.json` preserves source filenames/hashes/line numbers without copying opaque provider blocks.

- Original Anthropic ten-round failure remains a failure/limit result, not a successful campaign. The original empty late process monitor provides no OS proof.
- OpenAI’s two-turn run preserves conversation and **policy refusal only**. The prose saying it used `printf` is not credited as execution; the returned tool result is a refusal with no side effect.
- Gemini’s two-turn run has actual tool side-effect/readback evidence. The external consumer has actual registered `shout` invocation, correlated result, and timestamped root-owned diagnostic—not merely an import/build.
- The authorized single four-turn Anthropic follow-up completed on the unchanged artifact/tree. Its actual jobs show a real approximately **8-second** Go test failure, repaired source, and a subsequent passing test; a real HTTP server and requests across turns; and a real Delve session. This is not a measured 90-second test.
- Delve set/hit/continue are separate consumed returns for the same job 8: source JSONL lines **42, 46, 50**. Requiring all three in one return was a faulty verifier predicate, not a code defect. Current Chapter 4 requires Delve, not GDB.
- Final server handle 2 was observed running, then received matching `job_killed` reason `shutdown`. Saved OS evidence records PID 83531 last alive 0.264340 s before CLI exit and absent afterward. I reviewed that receipt, not today’s PID table. The model’s “shuts down naturally” wording is not credited; the evidence shows root-owned EOF shutdown.

No paid calls were rerun after restart or for this review. The public execution/report record discloses initial plan approval, source-versus-binary clarification, correction of model/key-file proposals before credential read/call, the useful-server checklist pointer, one bounded follow-up approval, recovery receipt locations/debugger predicate correction, and report-only UTC/GDB wording correction. This assistance is not misrepresented as unaided student discovery or a prior formal code-correction round. The early SIGQUIT-interrupted student long test is neither a pass nor a candidate failure; only its later completed runs and this review’s completed runs are credited.

## 6. Teaching feedback and review setup failures

Demonstrated teaching feedback is narrow: the chapter’s “very next call, including another `tool_limits`” rule is useful and sufficiently explicit, but deserves a chained-partial-settings regression. A loud consumption annotation alone does not prove reset semantics. Preserve consumed-return boundaries when checking interactive debugger progress; model prose is not evidence of side effects or process lifetime. The source/executable distinction and permitted single bounded follow-up were material assistance and are appropriately disclosed. This is feedback for the later author handoff, not a manuscript approval or requested rewrite.

Reviewer setup failures are retained and excluded from candidate verdicts:

- An initial read-only Git invocation reached its pager; only that owned pager was closed. Subsequent Git operations used no-pager controls.
- A few report/log paths were initially guessed incorrectly and then resolved through the pin/manifest; no missing file was treated as missing candidate behavior.
- First external-probe attempt used wrong public API signatures/types and failed compilation. Second reached config validation but omitted the required system prompt. Both sets of logs/receipts and source snapshots are retained as `setup-*` and `setup2-*`; neither is a behavioral failure.
- The final local-artifact verifier first assumed an object rather than a list. The corrected read-only verifier matched all 48 records.
- The successful initial public probe source and later added noncanonical-input source are retained separately so source/receipt attribution is not silently rewritten.

No candidate-source mutation or repair was made. The C1 failing and explicit-reset passing controls are input perturbations. Earlier entry/adapter mutation evidence is not relabeled as sensitivity of this student implementation.

## 7. Handoff

Return C1 as the first code-correction handoff. Request the focused chained-settings regression, a new frozen candidate/report with exact identities, and the required gates on the corrected candidate under the existing procedure. No additional paid call is requested by this review. No acceptance tag, commit, stage, or author-phase transition was performed.

Evidence directory: `docs/edition-3-notes/ch04-code-review-evidence/`. Start with `input-identities.json`, `baseline-runs.json`, `public-api-runs.json`, `pending-rearm-reproduction.log`, `pty-noncanonical.json`, `retained-live-tool-facts.json`, and `final-identities.json`. The directory manifest hashes all retained review artifacts. The parent owns saved-tail verification, join closure, and final acceptance.
