# Edition 3 Chapter 3 — finite entry audit

## Disposition

**Entry alignment still needed; no Chapter 3 student implementation launched, reviewed for acceptance, or accepted by this audit.** The existing Edition 2 Chapter 3 grader is a useful baseline, not a complete current Edition 3 gate. Its positive/negative controls work for several historical requirements, but a late-narration control still earns 100/100. There is no `grade-edition3 -ch 3` implementation. The smallest valid entry path is the existing Chapter 3 executable grader **plus** the existing Edition 3 Chapter 1 lifetime gate against the identical binary, supplemented by an explicit selected-model dispatch/source-API verification gate and independent source review. Root owns approval, launch, and later author timing.

Reference HEAD: `f7530a6643ca00bd2f59a8da60021c567cd7d174`.
Student predecessor: `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`, clean, with `edition-3-ch02-r1` peeled to that commit. No student source was changed. The audit rebuilt that predecessor only to verify inherited gates, **not** to grade it as a Chapter 3 candidate.

Fully read the required AGENTS, coding/reviewer skills, current Edition 3 execution/readiness/carryover, entire workflow and Edition 2 carryover, voice, chapter procedure, and actual Edition 3 Chapter 3. Recovered the truncated chapter-procedure output separately. Current Edition 3 procedure supersedes historical paths and phases. Also inspected relevant Chapter 1/2 public sections and grader implementations. Private reference/harness knowledge in this note is reviewer-side; do not send its fixtures, control patches, or reference code to the student.

Only this note and `ch03-entry-evidence/audit-*` were written by this worker in the repository. Other untracked files were left alone. No source edits, commit, deployment, restart, push, child launch, paid request, or executable artifact in the repository. The parent separately supplied the Gemini readiness receipt described below.

## Public-contract basis

### Chapter 3 dispatch boundary

`book/edition-3/chapter-03.md` §3.2, especially lines 162–180 and 212–228, requires:

- Preserve ordered assistant parts and call identity. When no provider ID exists, synthesize stable identity from captured event/part position; it must survive logging, execution, results, and replay.
- For the **selected model whose policy requires it**, ordinary visible assistant text must be nonblank and precede the **first** tool call in the **same response**. One explanation can cover a batch. Earlier dialogue, tool output, hidden thinking, readable reasoning summaries, and text after the first call do not qualify. The runtime checks presence/order, not prose quality; there is no heading or minimum-length rule.
- Reject the **entire** unexplained batch **before any handler**. Preserve the assistant response; return one matching error result per call, stating nothing ran and asking for an explained retry. Consume the existing round budget; do not add an unbounded repair loop.
- Test an observable side effect and the explained retry. The rejection must have no effect; the retry must really execute. Do not make this a universal policy for other models.
- Actual **stream-arrival** order belongs to Chapter 6. This audit does not demand it in Chapter 3.

The same chapter requires sequential tools in requested order, result pairing across vendors, all six concrete tools, their published defaults/error semantics, useful failure results rather than crashes, and preserved earlier chapters. §3.7 step 4 also requires a real tool probe per vendor and a small repair with the resulting agent. These live obligations are later candidate evidence, not facts established by a fake grader or this entry audit.

### Inherited obligations

- Chapter 2: deterministic event/history reduction and replay; no fabricated unsupported IDs/model identity; exact vendor/model/surface gating for model-bound opaque/signature material; Gemini function-response naming; tool results before text in merged Anthropic user messages; measured four-way usage; redaction preservation, transactional load, ephemera, and the public `dump`/`render` CLI. Native OpenAI Responses is **not** retroactively required: Chapter 2's executable scope remains Chat Completions.
- Chapter 1 §§1.3–1.5: real public-library ownership and parent chain; required parents; configuration/usage ownership; bounded connection/TLS setup but **no blanket model-generation or response-header deadline**; caller cancellation; failure must not leave a completed/unmatched conversational exchange; missing usage is not zero; measured usage survives an unusable/empty answer; safe diagnostics and no concealed paid retry.
- Accepted Chapter 2 review evidence remains historical preservation context, not fresh tests of Chapter 3. Its ID-less-pair and exact-provenance reasoning defects were corrected, and the new chapter must not regress them. Many of these properties require source/API review: historical grader success alone cannot prove them.

### Supported CLI layout

The public root library with `cmd/ch01` is valid. Chapter 1 §1.7 explicitly supports building `./cmd/ch01` to an absolute executable path and grading that path. `internal/grade/harness.go:Build` accepts an existing executable directly; for a directory, it queries `./cmd/` and builds that target if its package is `main`, otherwise it builds root `.`. It does **not** discover named subcommands or arbitrary nested command layouts. Existing `TestBuild` covers the old directory cases and rejects a library-only directory. This is not a demand to turn the public root into `package main` or rename the CLI.

Chapter 3's illustrative root-build examples do not revoke the executable escape hatch. The entry recipe must name the real target, not assume `go build .` works for this accepted architecture.

## Existing grader coverage and gaps

### Real entry points

- `cmd/grade-edition3/main.go` accepts Chapter 1 only. A direct invocation of the built grader with `-ch 3` exited **2**, printing `usage: grade-edition3 -ch 1 [-json] submission-directory-or-executable`. No Chapter 3 run occurred.
- Chapter 3 currently means `cmd/grade-edition2 -ch 3`, implemented by `internal/gradeedition2/ch03.go`, wrapping `grade.Ch3Run` / `grade.Ch3Evaluate` and Edition 2 `Chapter2`.
- Edition 2 Chapter 3 replaces the inherited `ch2parity` verdict with the Edition 2 Chapter 2 result; it includes the corrected version header, result naming/grouping, capture identity, redaction/ref scope, and existing replay fixtures. It does **not** call `gradeedition3.Chapter1` and so does not inherit the long header/body deadline gate.

### What the existing Chapter 3 checks actually inspect

The ten check IDs remain `ch2parity`, `toolsdecl`, `toolloop`, `multiblock`, `readtools`, `mutatetools`, `writeguard`, `runcommand`, `toolerror`, `editcontract` (100 points total). Source locations: `internal/grade/ch03_checks.go`, scenarios and log/wire decoding in `internal/grade/ch03_harness.go`, file extensions in `internal/gradeedition2/ch03.go`.

They test tool declarations in all three dialects; a two-call loop and a parallel-call assistant turn; exact recorded text and calls in selected happy-path turns; call/result ID matching; first-position Anthropic result projection; basic log ordering; the five file handlers and shell; overwrite refusal; all-occurrence/absent/empty edit semantics; missing files, unknown tool and malformed args; follow-up response after errors; cumulative usage. Edition 2 adds real `search_files.context_lines` and `write_file.append` effects.

**Useful limits, not a complete specification:** the log-order check verifies response → called → returned, not a visible-text-before-first-call policy. Text checks count parts and compare expected strings; they do not inspect part ordering. Error checks accept ordinary dispatch errors and prove a later request, not corrective policy rejection. The harness has no selected-model dispatch activation contract, no whole-batch narration refusal fixture, no no-side-effect assertion for refusal, no policy-specific result/history pairing or retry-bound test, and no unaffected-model control. Happy-path provider IDs are populated; ID-less Chapter 3 execution pairing is not covered by those cases. The 20-second process cap can stop a hang but does not prove the specified round budget.

The existing fake has `Reply.Text`, `Thinking`, and `Tools`. Anthropic/Gemini bodies emit thinking, then text, then calls (`internal/fakevendor/fake.go`); it cannot express arbitrary interleaved/late text through that simple Reply structure. Chat Completions separates content from tool calls; do not confuse JSON object-key order with response-part order. A focused ordered-response/source-API fixture is sufficient for Chapter 3; it need not create a streaming harness.

The Edition 2 search/append extension's first reply is **tool-only**, with no text (`internal/gradeedition2/ch03.go:79–83`). If a future selected-model test marks that fake model as requiring narration, this happy-path fixture conflicts with the new requirement. Conversely, leaving all existing fake models unselected means none tests the new gate. This is a **conditional fixture/activation gap**, not proof that every conforming selected-model implementation must fail the old grader. Do not solve it by enforcing the policy universally or inventing a private environment switch.

## Actual finite runs

Receipts: `ch03-entry-evidence/audit-focused-tests.txt`, `audit-results.json`, `audit-controls.diff.txt`, `audit-manifest.json`. JSON reports retain verdict/check details but intentionally omit captured stderr/request bodies. Source and executable SHA-256 values are in the manifest. Temporary binaries and control source were kept outside the repository.

### Existing focused harness tests — all passed

From reference root:

```sh
go test ./internal/grade ./internal/fakevendor ./internal/gradeedition2 ./internal/gradeedition3 \
  -run '^(TestBuild|TestRoutesAllThreeDialectsByPath|TestCycleWrapsAndUniquifiesIDs|TestChapter1VersionAndPreservedChecks|TestChapter2AcceptedSolution|TestChapter3AcceptedSolution|TestGenerationLifetimeControls|TestHistoricalTimeoutIsExceeded)$' \
  -count=1 -v
```

Exit **0**. Observed coverage:

- `TestBuild`: root, command container, preferred command, library preferred command, and expected library-only failure.
- Fake provider: three-dialect routing and cycling/unique call IDs.
- Chapter 1 header correction: dated correct version accepted; missing/empty/undated versions rejected; preserved old verdicts/score shape.
- Historical Chapter 2 and Chapter 3 solution controls: **100/100**. These are older references, not Edition 3 correctness proofs.
- Chapter 1 short lifetime controls: good implementation accepted with header/body delay; total-HTTP-timeout and header-timeout defects detected. `TestHistoricalTimeoutIsExceeded` confirms the production delay exceeds 60 seconds; it does not itself sleep 62 seconds.

### Accepted predecessor, same explicit binary

Built from clean student root with `go build -o "$scratch/student-cli" ./cmd/ch01`. Built both graders once from the reference tree. The same `student-cli` was used for:

| Actual invocation | Result | Meaning |
|---|---|---|
| `$scratch/grade-e3 -ch 3 -json $scratch/student-cli` | exit **2**, usage | Confirms unsupported CLI chapter, not student failure |
| `$scratch/grade-e2 -ch 2 -json $scratch/student-cli` | exit **0**, **100/100**, `passed:true` | Fresh inherited Chapter 2 baseline only |
| `$scratch/grade-e3 -ch 1 -json $scratch/student-cli` | exit **0**, **100/100**, `passed:true` | Fresh inherited Chapter 1 baseline, including **62-second delayed headers and 62-second delayed body**, both passing |

No Chapter 3 grader was run against the Chapter 2 predecessor: missing tools there would be expected and would not answer this entry question.

### Temporary historical Chapter 3 controls

Each scratch copy changed only `solutions/ch03/internal/agent/agent.go`, was built to an explicit executable, and ran the existing Edition 2 Chapter 3 grader. Diffs are retained as text, not executable fixtures. These are diagnostic controls, **not** proposed student implementation or complete policy implementations.

| Single boundary perturbation | Actual outcome | Interpretation |
|---|---|---|
| Remove ordinary text from tool-bearing recorded responses | exit **1**, **70/100**; `toolloop` and `multiblock` fail | Existing checks detect missing happy-path text; tool handlers still pass. Does not test rejection before effects. |
| Replace that text with whitespace | exit **1**, **70/100**; same two failures, exact-text mismatch | Existing expected-string checks are sensitive here; this is not a model-policy whitespace-rejection test. |
| Move text after tool-call parts before recording/dispatch, retaining text/calls/behavior | exit **0**, **100/100**, all ten checks pass | Concrete **ordering blind spot**. The harness accepts the mutated ordered response; it does not establish an explained-before-tools boundary. Existing fake models are not declared policy-selected, so this is sensitivity evidence, not an Edition 3 candidate verdict. |
| Add a diagnostic pre-dispatch rejection for tool responses without preceding nonblank text, return matching errors, keep existing loop budget | exit **1**, **91/100**; `readtools` and `mutatetools` fail on search-context/append errors | Demonstrates the tool-only extension fixture conflict **if** its model is subject to narration. Control is intentionally unconditional, not a valid selected-model implementation. Most other old checks remain green. |

No simulated result is presented as a run. Lookup misses for guessed filenames (`ch02_replay.go`, `ch01_lifetime.go`, a nonexistent `ch02-entry-gate.md`) were inspection errors only; actual implementations were subsequently located in the files cited above. They were not failed candidate tests.

## Smallest valid alignment, approval required

### 1. Prefer composition over a new grader framework

For a future frozen Chapter 3 candidate, the minimal existing executable recipe is:

```sh
# Candidate repository; target is the actual accepted nested CLI:
go build -o /absolute/scratch/ch03-agent ./cmd/ch01

# Reference repository; both must use that identical executable/hash:
go run ./cmd/grade-edition2 -ch 3 -json /absolute/scratch/ch03-agent
go run ./cmd/grade-edition3 -ch 1 -json /absolute/scratch/ch03-agent
```

Require exit 0 **and** `passed:true` for both. Chapter 3 already calls Edition 2 Chapter 2; a third routine Chapter 2 invocation adds duplicate work, not lifetime protection. Keep the separate source/API gate below; do not call these two commands full current-contract coverage by themselves. The commands above are a proposed candidate recipe, **not** a claim that they have been run on a Chapter 3 candidate.

### 2. Add a narrow, public-contract dispatch gate

Before launch, root should approve an implementer-neutral external executable/source-API verification plan that selects a genuinely policy-requiring model using the candidate's public configuration/capability seam. The chapter describes **selected-model policy**, not a public `REQUIRE_NARRATION` environment flag or a particular internal package/table. No such flag should be invented solely for grading. If a CLI-only check needs an activation identity not yet settled publicly, make that mapping explicit and obtain approval first; do not encode an undisclosed model-name trick in private fixtures.

Use temporary files/counters as observable effects, and actual emitted facts/requests to verify:

1. Selected model, missing or whitespace text: **zero** handlers/effects; every call gets exactly one matching policy-error result and original assistant parts remain intact. Use at least two calls so partial-batch execution cannot hide.
2. Prior-turn text, tool output, hidden thinking/summary-only, and text after the first call all remain insufficient. A late explanation cannot rescue an already-started batch.
3. One ordinary nonblank text part before two calls permits **both**, in order. No minimum length/heading, no semantic prose grader.
4. An explained corrective retry really executes once and preserves original rejection facts, ID/name pairing, results-first provider projection, usage, and continued conversation.
5. Repeated unexplained responses terminate through the **existing** round budget. Count requests/results; do not treat a killed hung process as success. A normal provider-error path must not acquire paid auto-retries.
6. A nonselected model remains unaffected. Include an ID-less call pair and exact-provenance replay preservation in source/API review, reusing Chapter 2 evidence where applicable.

Feasible sensitivity controls for that future gate: skip the guard; check any text anywhere; use `TrimSpace` incorrectly/omit it; count thinking as narration; inspect only the first call; dispatch before checking; suppress one rejection result; lose call identity; erase rejected assistant facts; reset the round counter on repair; enforce policy on every model; and reject every tool call. A correct implementation must pass both the no-effect negative case and actual successful explained execution. **These new controls were not implemented or run in this audit.** Existing scratch runs above cover only their stated perturbations.

Source/API controls avoid prescribing private constructor names or a new configuration framework. If root instead wants a repeatable Edition 3 CLI wrapper, the smallest source change is an Edition 3 Chapter 3 overlay composing the existing Chapter 3 and Chapter 1 work and a mandatory policy check—not rewriting old graders. A local ordered-response fixture can test late text without modifying every historical fake. Keep original graders callable and historical scores meaningful. Any correction to old happy-path narration fixtures should be narrowly scoped, independently reviewed, and sensitivity-tested. **No source change is authorized or made by this note; ask root before implementing any of these alternatives.**

### 3. Retain independent review and live boundaries

Review ownership/API/configuration boundaries, setup-vs-generation protections and caller cancellation, failed-exchange/history and measured-cost behavior, exact opaque/signature replay gating, ID-less call pairing, credential exclusion from tool environments, and meaningful tool errors. The existing report does not cover every concrete file-tool edge/default either. Do not impose a new numeric coverage threshold, root-main layout, native Responses route, or future streaming/job/caching/memory features.

Candidate live evidence must use the actual CLI/agent, per-vendor tool calls, recorded usage and tool/repair outcomes. Fake success, historical solution success, and provider discovery are different evidence classes. Do not require re-running the whole book as this entry gate.

## Model-policy reconciliation

- `docs/edition-2-carryover.md:436–445` records the user's restriction: Gemini demonstrations use **3.0 Flash or newer**, specifically requesting **3.8 Flash** for hint behavior; never silently substitute an older generation.
- Current Chapter 3 §3.6 (`:541–558`) requires discovered identifiers and verified tool support. It labels the listed `gemini-3.8-flash`/other defaults as **historical observations**, not fixed binary defaults. Discovery does not override the user's minimum-generation restriction.
- Accepted Chapter 2 evidence really used `gemini-2.5-flash-lite` at the original pre-correction candidate, with that attribution retained in `ch02-student-review-r1.md`. This is a policy discrepancy to disclose, **not** permission to reuse it or a reason to rewrite its evidence or reopen the accepted predecessor silently. Correction-round evidence did not perform new paid calls.
- Parent supplied `ch03-entry-evidence/gemini-38-readiness.json`: 2026-10-16T22:11:49.765473+00:00, `gemini-3.8-flash`, generateContent, HTTP **200**, returned model `gemini-3.8-flash`, reply **OK**. It reports 6 prompt, 1 candidate and 67 thought tokens, 74 total. I read the sanitized receipt; **this worker made no live request**. This establishes only minimal access to the requested model, **not** tool support, CLI integration, repair success, hint behavior, or a selected-model narration capability mapping.
- Future smoke tests should discover/confirm the current allowed model and surface, record actual identity/results, and not silently fall back to 2.5. If an access restriction prevents an authorized selection, report it and ask root; do not invent a new assignment requirement or generalize readiness from a different model. Parent retains provider-choice/launch authority.

## Evidence limits and completion

Fresh results are limited to the named focused tests, the predecessor's two inherited gates, unsupported Chapter 3 CLI invocation, and four historical scratch controls. Not run: a Chapter 3 student candidate, new policy gate/controls, full-book grading, source race/coverage suite, candidate live tool/repair exercises, real streaming-order checks, GUI/product integration, author phase, deployment/restart. No conclusion here is code acceptance.

All launched finite jobs completed. No children were launched. Reference tracked diff remained empty and the student worktree remained clean at the recorded predecessor after the runs. Existing untracked artifacts and the parent's concurrent execution/readiness notes are not audit source changes. Root can now choose and authorize the minimal dispatch-gate alignment, independently review any source change if selected, and decide student launch timing.
