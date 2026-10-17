# Edition 3 Chapter 5 entry audit

Date: 2026-10-17 UTC. **Read-only entry audit; not a Chapter 5 student acceptance.** No Chapter 5 student implementation existed at the audited predecessor. No grader, student, manuscript, or historical solution source was changed in either original repository. Only this report and its adjacent evidence directory were added by this audit. Private scratch drivers/probes and binaries stayed outside both repositories. No child delegates, paid/provider calls, credentials, installation, push, restart, or live GUI were used.

## Frozen identities and governing reads

- Main: `686c1ef027c32b69d5c44230004b42d59045989a`.
- Student: `ec4a350f564e820d398f9125fd73a9b3b27129cd`, clean before the audit; accepted predecessor `edition-3-ch04-r1` peels to this commit. Annotated tag object: `5af671222370c6ebff1f369e93a7c3cf9d6e2ec2` (not the commit).
- Go: `go1.25.4 darwin/arm64`; scratch snapshots were made using `git archive HEAD`, not worktree copies. Archive builds therefore do not embed Git VCS revision metadata: archive origin, commands, and hashes supply that provenance.
- Fully read current `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`, `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, `docs/edition-3-execution.md`, `docs/edition-3-student.md`, `docs/edition-3-readiness.md`, `docs/edition-3-carryover.md`, `.agents/skills/edition-3-code-reviewer/SKILL.md`, and `book/edition-3/chapter-05.md`. Combined-output truncations were recovered with smaller reads. `source-hashes.json` records the exact document bytes, line counts, and SHA-256 values, plus inspected grader and student source.
- Current Edition 3 student-first procedure supersedes old reference-first/literary workflow instructions. The historical Chapter 5 solution was used only as a private audit control, never as a student substitute or repair handoff. Existing unrelated main-tree untracked work was left alone.
- Accepted Chapter 4 completion record is `docs/edition-3-notes/ch04-completion.md`; its accepted source/API obligations remain inherited.

## Bottom line

1. **There is no current `grade-edition3 -ch 5`.** The executable accepts only 1 and 4; the measured Chapter 5 request exited 2 before grading.
2. **The manuscript's Edition 2 shorthand is not a valid complete invocation for this student's layout.** `grade-edition2 -ch 5` reuses its sole positional argument as executable/build target and core source root. The accepted student's root is an importable `ensemble` library; its CLI is `cmd/ch01`; its Chapter 4 consumer is `examples/ch04-consumer`, not the legacy discovery location. Root, nested CLI directory, and prebuilt executable each fail a different part of that implicit contract. Explicit `-consumer` does not separate the CLI from core.
3. **The existing Chapter 5 package API already has the necessary separation.** A scratch-only driver calling it against the correct CLI and core reached all checks faithfully. The predecessor earned **35/100, false**, from Chapter 4 parity, inherited replay, and clean common. This is the expected missing-Chapter-5 baseline, not a Chapter 4 regression.
4. **Corrected Chapter 5 green alone is not Edition 3 acceptance.** A private historical positive earned 100/100, true, yet failed both actual 62-second Edition 3 lifetime probes with score 100/100, **false**. The existing Chapter 5 suite also does not establish the selected explanation policy or Chapter 5's default 200 dispatch batches and final result-consuming request.
5. Smallest recommendation: expose the already-separated Chapter 5 API parameters through a small explicit entry adapter (or use an audit-only adapter until authorized), preserve existing checks, and separately retain the inherited Edition 3 lifetime/public-policy gates plus focused public tests/source review for the chapter's named gaps. **No framework, student-layout rewrite, speculative feature, or change to an older grader is necessary.**

## Current contract (not inferred from historical implementation)

`book/edition-3/chapter-05.md` is authoritative for new work:

- One actor/conversation owner per Agent. `Post` is a nonblocking mailbox door, including under burst input; live prompt/hint is attached to the active turn or starts work when idle. Interrupt ends the active turn, not the actor, and retains real results from continuing jobs.
- Blocking Ask queues a distinct request/turn with its own buffered, reliable reply. Concurrent blocking requests are not hints. Request acknowledgement must not rely on bounded/droppable observations, shared last-text state, or competing Wait subscribers.
- Observer snapshots preserve per-Agent identity, state changes and complete parts; progress is nonblocking and may drop. Agent/Job ownership, full result metadata, and per-turn correlation must survive the actor migration. Public state/metadata access and configuration, initial input, and turn IDs must be real public seams, not private-only shortcuts.
- `Ref`/feature vocabulary remains in common; lookup/rendering policy does not. Supported local images embed actual bytes; unsupported media explicitly names model and type; unknown models have no permissive media default. Ordinary text does not require a media capability row.
- Canonical-log prefix rendering must equal actual captured provider bytes. Mid-HTTP hints cannot be consumed by the wrong request, lost, or duplicated. The chapter requires a regression around that interleaving.
- Three distinguishable author/editor/reviewer Agents in a separate-module real workflow use the public API; manuscript mutations/results must flow through the intended handoffs. Markers alone do not prove real tool effects or downstream visibility.
- The chapter explicitly requires a **default budget of 200 tool-dispatch batches**, counts refusals against the budget, permits the final request that consumes the last allowed batch's results, and rejects any subsequently requested over-budget batch before effects. An override, if exposed, must be documented as intentional rather than a hidden compatibility fallback. These are current requirements, not an optional compatibility request to retain the predecessor's ten-round cap.

Inherited obligations remain independent:

- Chapter 1: let a legitimate generation finish; no arbitrary whole-request/response-body deadline, preserve complete answer and usage. Test real delayed headers and real delayed body separately; source review still needed for deadlines longer than the measured delay and non-default paths.
- Chapter 3: selected model policy requires nonblank ordinary assistant text **before the first call in the same response**. Reject the whole unexplained batch before any handler effect; preserve the response and a matching error result for each call; explained retry shares the normal budget. Prior text, late text, whitespace, hidden thinking and reasoning summaries do not qualify; disabled policy remains different. Signed ordinary Gemini text remains visible with exact part/source binding; hidden thinking remains hidden. Neither a global stricter default nor silent refusal is an acceptable substitute.
- Chapter 4: public importable root, separate real consumer using that root, pure shared vocabulary and first-party leaf allowance, no spoke-to-spoke dependency, actual parent/ownership boundaries. Jobs retain full files and structured metadata, callback delay is not a process deadline, output is capped/cursor-correct without losing full bytes, completed jobs can be waited on, blocked waiters wake on output/completion/kill, kill is process-group termination and retained `killed`, and a late completion cannot resurrect it. Preserve limits semantics and consumer effects/diagnostic provenance. A bare topology score is not a source-architecture verdict.

## Existing entry points and identity boundaries

Sources: `cmd/grade-edition3/main.go`, `cmd/grade-edition2/main.go`, `internal/grade/harness.go`, `internal/gradeedition2/ch05.go`, `ch04_topology.go`, `ch04_consumer.go`, `ch05_workflow.go`.

- `grade.Build(TARGET)` builds a directory using `go build -o TEMP .`; it does not discover a nested CLI. For a non-directory it accepts the path as the executable. Building a library root can succeed yet yield a non-executable archive: this predecessor's root subsequently failed with `permission denied`.
- Edition 3 Chapter 4 already separates positional CLI target from `-core`, `-consumer` and `-consumer-case`.
- Edition 2 Chapter 5 calls `Chapter5(bin, target, consumer, workflow, example, models)`. There is no `-core` flag there. Passing `cmd/ch01` or a binary repairs execution but makes source checks target the wrong object.
- Consumer/workflow module discovery is a convenience for legacy locations, not a binding to this student's layout. Use explicit paths. Consumer defaults assume the legacy word-count example; supply the accepted `examples/ch04-consumer/consumer-case.json` instead. The absent Chapter 5 workflow was supplied explicitly as an absent path to record a real missing-feature failure, never silently substituted.
- `publicDependencyTarget` inspects dependencies of the actual executable target, requires distinct module declarations and a positive public-core import, rejects direct private imports, and checks the resolved module directory against the canonical core directory. A private mixed-core/apps control failed precisely this directory check; it does not accept a matching module name with a stale `replace` target.
- **The CLI binary itself is not authenticated to `core`.** A separate control using the historical binary but accepted student core/consumer passed Edition 3 Chapter 4, 160/160. This is a measured provenance boundary, not evidence that the student used the historical binary. Build the candidate CLI from the frozen candidate and record its hash/command; source and consumer checks cannot prove that relation themselves.

### Smallest valid explicit invocation today

There is no one-shot existing command that supplies all necessary Chapter 5 identities for this layout. The smallest measured adapter called the unchanged API:

```text
grade.Build(TARGET)
gradeedition2.ReadConsumerCase(CASE)
gradeedition2.ReadMediaModels(MODELS)
gradeedition2.Chapter5(bin, CORE, CONSUMER, WORKFLOW, example, models)
```

It emitted the unchanged JSON report and exited 1 when `Passed` was false, 2 on entry/build errors. Built inside the frozen scratch grader module (to respect Go `internal` visibility), it added no replacement checks. Its interface and measured predecessor command were:

```sh
# cwd: <AUDIT_SCRATCH>/student
go build -o <AUDIT_SCRATCH>/student-cli ./cmd/ch01
# cwd: <AUDIT_SCRATCH>/grader
go build -o <AUDIT_SCRATCH>/ch05-audit ./cmd/ch05-audit
<AUDIT_SCRATCH>/ch05-audit \
  <AUDIT_SCRATCH>/student-cli \
  <AUDIT_SCRATCH>/student \
  <AUDIT_SCRATCH>/student/examples/ch04-consumer \
  <AUDIT_SCRATCH>/student/examples/ch05-workflow \
  <AUDIT_SCRATCH>/student/examples/ch04-consumer/consumer-case.json ''
```

The empty last argument selects `ReadMediaModels`' existing defaults; it is not a wildcard or guessed capability table. For a future candidate, supply a grader-only mapping to that candidate's explicitly implemented model table if necessary. The future workflow path must be the student's actual separate module, not the nonexistent audit placeholder. The adapter exists only in disposable scratch and is **not** a newly available repository command.

Independently invoke the current supported gates with the same freshly built CLI:

```sh
<AUDIT_SCRATCH>/grade3 -ch 1 -json <AUDIT_SCRATCH>/student-cli
<AUDIT_SCRATCH>/grade3 -ch 4 -json \
  -core <AUDIT_SCRATCH>/student \
  -consumer <AUDIT_SCRATCH>/student/examples/ch04-consumer \
  -consumer-case <AUDIT_SCRATCH>/student/examples/ch04-consumer/consumer-case.json \
  <AUDIT_SCRATCH>/student-cli
```

Allow at least the full two 62-second delays plus grader/build/exit grace. These actual controls had a 240-second outer watchdog. Tool callbacks merely returned handles; their early return did not terminate observation.

## Measured controls

Exact argv, cwd, elapsed seconds, exit status, watchdog state and full sanitized stdout/stderr are adjacent in `ch05-entry-evidence/NAME.{meta.json,stdout,stderr}`. No score is inferred from exit status alone.

| Control / evidence stem | Seconds | Result |
|---|---:|---|
| `e3-ch5-unsupported` | 0.006 | exit 2: only Chapters 1 and 4 supported |
| `e2-ch5-root` | 0.638 | exit 2: root archive cannot execute (`permission denied`) |
| `e2-ch5-cli-directory` | 8.698 | exit 1, 15/100: nested CLI runs and inherited replay passes, but it is the wrong core for topology/common checks; workflow/media features absent |
| `e2-ch5-binary` | 8.421 | exit 1, 15/100: inherited replay passes, but binary is not a source directory; workflow/media features absent |
| `predecessor-api-ch5` | 8.646 | exit 1, 35/100: corrected Chapter 4 parity, inherited replay and hub pass; absent workflow and image/refusal contract fail |
| `e3-ch4` | 8.797 | exit 0, 160/160, true; accepted predecessor source, CLI and external consumer |
| `e3-lifetime` | 124.055 | exit 0, 100/100, true; actual 62-second headers and body phases each pass with complete answer/usage |
| `public-policy-independent` | 0.221 | exit 0; public selected-policy whole-batch/no-effects/retry across three adapters and presence/order/disabled controls pass |
| `cli-policy-independent` | 0.199 | exit 0; true/false/default public CLI policy-selection controls pass |
| `predecessor-batch-boundary` | 0.326 | exit 1; 9 batches → 10 requests, 9 results, final `done`; 10 or 200 batches → 10 requests, 10 results, empty answer and `tool round limit reached (10 rounds)` |
| `reference-e2-ch5` | 12.305 | private historical positive only: exit 0, 100/100, true |
| `reference-e3-lifetime-retry` | 120.056 | private historical negative: exit 1, score 100/100 but false; both real 62-second probes lose complete answer/usage under historical 60-second timeout; no outer watchdog expiry |
| `mixed-identity-private` | 8.600 | exit 1, 35/100: historical apps resolving historical core rejected against supplied student core |
| `mixed-binary-ch4-private` | 8.547 | exit 0, 160/160: executable/source binding not enforced by this source/consumer gate |
| `predecessor-tests-retry` | 125.652 | exit 0: full accepted student `go test -count=1 ./...`, including actual 62-second header/body tests, with a 300-second watchdog |
| `reference-public-tests` | 2.822 | private historical suite exit 0; not a student result |
| `grader-lifetime-fixture-retry` | 1.351 | exit 0: Chapter 1 short-delay positive/negative unit controls and command tests after required historical fixture restored to scratch |

The batch probe used the real accepted public `ensemble.New`/`NewAgent`/`Ask`, ordinary registered `read_file`, and a loopback provider returning N tool batches then final text. It did not replace Agent/renderer/job logic. The private test was added only to the archived scratch student, then removed; only its observed result is exported. This is an **entry baseline**, not an accusation that the accepted Chapter 4 work should already satisfy Chapter 5's larger budget. No 200-batch student implementation is available to test.

### Retained failures and incomplete attempts

- `build-reference`: exit 1, 0.006s; auditor initially guessed nonexistent `./cmd`. Corrected root build `build-reference-root`: exit 0, 0.289s. The first `reference-e3-lifetime` invocation then exited 2 because its binary did not exist. Neither was counted as a killed timeout control. The retry above is the real measured timeout failure.
- `predecessor-tests`: `go test ./...`, 120-second watchdog, killed at 120.009s (exit -9), incomplete. The suite itself includes sequential 62-second header/body tests; this allowance was insufficient. It is not a student regression. The 300-second `-count=1` retry completed successfully in 125.652s, exit 0; both the incomplete first run and completed retry are retained.
- `grader-tests`: broad package test attempt exited 1 at 0.760s because the initial minimal grader archive omitted `solutions/*`, which package tests use as fixtures. Its many fixture/build failures, including later chapters, were retained; no GUI was built or launched successfully and these are not Chapter 5 behavioral findings.
- `grader-self-contained-tests`: 0.213s, exit 1: even the narrowed Edition 3 lifetime unit test needs historical Chapter 1 source. This wrong self-contained assumption is retained. Copying only tracked `solutions/ch01` to scratch yielded the successful narrow retry above. The full historical grader package suite was not retried.
- `public-policy`: the earlier same public test selection also passed, 0.233s. The independently recorded retry explicitly uses `-count=1`; it is the table's noncached evidence.

## What the current grader actually covers—and does not

Source trace, with measured results distinguished above:

| Existing check | Actual exercised boundary | Remaining current-chapter obligation |
|---|---|---|
| `ch4-parity` (10) | Runs corrected Chapter 4's Chapter 3 parity, real consumer, tool diagnostic, topology and job checks; all must pass | Does not re-add Edition 3 Chapter 1 lifetime checks or selected explanation policy; source architecture and broader public job regressions remain necessary |
| `not-deaf` (25) | Builds actual workflow module, local provider requests `think`, sends delayed hint, requires matching call-ID job report and hint observed first | Does not prove arbitrary mailbox burst behavior, a single owner, reliable private replies, all interrupt orderings, or that a progress drop cannot lose a request reply |
| `replay-is-live` (15) | Captures three actual HTTP request bodies across two turns/tool result and renders saved prefixes before matching `request_sent`; byte comparison except surrounding whitespace | No actual mid-HTTP hint injection in this probe; source-required interleaving regression remains necessary |
| `observer-fires` (15) | Parsed workflow observations include both state and part kinds | Does not fully establish snapshot isolation, nonblocking/droppable subscribers, Wait cancellation or private acknowledgement reliability |
| `multi-agent-observations` (15) | Distinct observed author/editor/reviewer labels, ordered expected completion markers | Cannot establish real manuscript-tool effects, stored edit handoff to reviewer, saturated-stream replies, or concurrency from marker output alone |
| `loud-refusal` (10) | Local image bytes via each adapter; unsupported audio and unknown-model image errors naming model/type | Not every media/ref kind or capability-row combination; ordinary text without media row and all MIME cases are not comprehensively tested |
| `hub-clean` (10) | `go list` topology of common allows standard library/first-party leaves and rejects spokes/external modules | Cannot tell whether behavior was hidden in shared vocabulary or whether root composition/parent interfaces are genuine |

The actual point weights total 100. No Chapter 5 checker counts 200 dispatch batches, verifies the final result-consuming request, or checks selected policy before side effects. The private historical positive proves that passing this suite does not establish Edition 3 lifetime. Public predecessor tests demonstrate the inherited policy positively; they must be preserved and rerun after actor migration, not presumed covered by `ch4-parity`.

Additional **source-only conclusions**, not claimed as executed candidate tests: preserve one owner and reliable per-request replies under concurrent Ask/Post and saturated/dropped observer events; validate required public state/config/turn identity; test pending hints at a real in-flight HTTP boundary; test the concrete workflow's manuscript tools and downstream edit visibility; keep chosen/default batch budget, refused-batch accounting and no-partial-over-budget-effects explicit. These are named manuscript/inherited obligations, not a proposed generalized framework or a demand for implementation choices beyond the chapter.

## Smallest recommendation for the next authorized step

1. Resolve only the entry seam: add a narrow Edition 3 Chapter 5 adapter that accepts separate CLI target, core, consumer/case, workflow and media-model mapping and delegates to the existing corrected Chapter 5 function, or retain the private explicit adapter for the next audit. Use the existing Chapter 4 separation as the pattern. Do not claim the command already supports this, change Edition 2's historical meaning, or require the student to turn its public root into `package main`.
2. At freeze, build the candidate's actual CLI from its exact clean revision. Record binary SHA-256, source revision, build command and resolved consumer/workflow module directories. Never reuse a historical or stale binary merely because it passes.
3. Grade the real student with explicit paths. Keep Edition 3 `-ch 1` (inspect **Passed**, not score), supported `-ch 4`, and public Chapter 3 explanation tests separate and visible; then run the Chapter 5 checks. Preserve failures and distinguish setup failures from contract failures.
4. Require the chapter's focused public integration/source evidence for reliable request replies, saturated observers, in-flight hints, real manuscript handoffs, default 200 batches/final result-consuming request and inherited job/public API behavior. No speculative subsystem or teaching repair fixture is needed.
5. This report/evidence may guide acceptance boundaries. Private historical source, audit probe code and repair fixtures must not be handed to the future student. Review the exact frozen student candidate independently when available; the present audit does not pre-approve it.

## Evidence and artifact hashes

Adjacent `ch05-entry-evidence/` is portable text only. `<AUDIT_SCRATCH>`, `<MAIN_REPO>`, `<STUDENT_REPO>`, and `<USER_HOME>`, and `<SYSTEM_TMP>` replace machine-specific roots. Local fake-provider diagnostics may include loopback ports and fake fixture constants; no ambient credentials were passed. Build/test environments inherited only PATH/TMPDIR where present and explicitly retained existing Go caches, with scratch HOME, `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOWORK=off`, `LLM_VENDOR=anthropic`. No downloading/installing occurred.

Key executable SHA-256:

```text
student-cli  c2ac49db3fb60d860b7663507fb64bf773acba573078140468ab38dd9f9d8999
grade3       03608eadbccf8a816a68ce8d253d9320cdfd9fcf637c3f7373b5e2914dca5190
grade2       3a23dfa6276a2adb9506d54d447347148a514d63dc6610ed508031738dcff40f
ch05-audit   0277946b6e4c128b1986072fb0d5511ecccc4697c5be5655259e64c54ea06217
```

`artifact-hashes.json` includes all private artifact hashes/build metadata, including the historical binary and adapter source digest, without distributing those artifacts. `source-hashes.json` binds source/doc inputs. `provenance.json` binds revisions, tag and environment. `SHA256SUMS` covers all exported evidence files except itself. These hashes identify the measured artifacts; path-dependent Go builds are not promised to reproduce identical binary bytes elsewhere.

## Final reconciliation

Resumed after a transport `no route to host` interruption from durable results; no completed behavioral controls were rerun. That interruption was not a test failure. The draft misreported the replay-positive scores: raw JSON establishes **35/100** for the explicit API predecessor baseline and **15/100** for each legacy CLI-directory/binary entry. The report now agrees with those unchanged raw records. Inherited replay is a measured positive, not an absent predecessor feature.

The properly budgeted full student-suite retry passed, all scheduled controls completed, failed/incomplete attempts remain exported, and the evidence checksum manifest verified. Final original-tree checks found both tracked trees unchanged, student HEAD still at the accepted predecessor, and student status clean. Main's unrelated untracked files were preserved. This is a completed entry audit with an entry-adapter recommendation, **not Chapter 5 acceptance**.
