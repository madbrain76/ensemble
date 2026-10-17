# Edition 3 Chapter 4 — finite entry/grader audit

## Disposition

**ENTRY ONLY. Alignment is required before student launch.** No Chapter 4 student was launched or accepted by this audit. Do not treat the historical Chapter 4 reference or a jobs score as Edition 3 acceptance.

The smallest valid path preserves the accepted root-library plus `cmd/ch01` layout, separates executable identity from core-module identity, and composes the existing Edition 2 Chapter 4 report with the existing Edition 3 Chapter 1 lifetime report against the **same explicit executable**. Require independent coverage of the inherited selected-model policy and the Chapter 4 concurrency/cursor/limits contract; do not invent CLI configuration requirements or a new grader framework.

## Scope and provenance

- Main HEAD verified: `bac07e644facb25ba3e5bffb2961674763fed4c4`.
- Accepted predecessor HEAD and peeled `edition-3-ch03-r1`: `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`. The tag object is distinct from its commit; compare the peeled commit. Student tree clean.
- Main tracked diff was empty at entry; pre-existing untracked files were preserved. Later concurrent changes to `cmd/grade-edition3/main.go` and new `main_test.go` were observed and left untouched. This is a **baseline HEAD snapshot audit**, not a review of those new working-tree changes. No manuscript, grader, application, student, original solution, or reference edits. Only this report and `ch04-entry-evidence/` were written in the live repository.
- Exact Git exports of both revisions were made in a temporary directory (recorded in `ch04-entry-evidence/scratch.txt`). All builds, mutations, fake-provider executions, consumer runs, and tests used those exports. No credentials or live calls, installs, service restarts, or pushes. Provider/credential environment variables were removed for controls; `GOPROXY=off`, `GOTOOLCHAIN=local`.
- Required reading completed: full current AGENTS, coding/reviewer skills, Edition 3 execution/readiness/carryover, workflow, Edition 2 carryover, voice, chapter-writing procedure, and full Edition 3 Chapter 4. Oversized/truncated procedure/chapter reads were recovered separately; workflow was read through its actual line 2434. Current Edition 3 student-first procedure overrides historic author-first phases/paths; explicit current model restrictions remain. Chapter 3 entry and completion receipts were consulted rather than rerunning their unchanged sensitivity matrix.
- Toolchain: Go 1.25.4 darwin/arm64. Source hashes and observed statuses: `ch04-entry-evidence/provenance.json`. Per-command argv/cwd/exit/elapsed/output hash: `runs.json`. Binary hashes: `binary-identities.json`.

## Actual composition and supported invocation

At audited HEAD `bac07e6`, `cmd/grade-edition3/main.go` supports **only `-ch 1`**. `-ch 4` is rejected, not an available Chapter 4 mode.

`cmd/grade-edition2/main.go` uses flags `-ch 4`, `-consumer`, and optional `-consumer-case`; chapter names are not positional subcommands. A missing target silently selects a historical solution, which is invalid student evidence.

The actual wrapper (`internal/gradeedition2/ch04.go`) runs historical `grade.Ch4Run/Ch4Evaluate`, **replaces `ch3parity` using the corrected Edition 2 `Chapter3` result**, and appends the external-consumer checks. Edition 2 Chapter 3 composes Chapter 2, which composes its Chapter 1 checks. Thus current Chapter 2/3 historical parity fixes are inherited. It does **not** invoke `gradeedition3.Chapter1`, nor independently establish the Edition 3 selected-model narration/refusal/retry policy. The Chapter 3 accepted receipt records a paired grade and independent policy evidence; that evidence establishes the predecessor, not immunity from Chapter 4 regression.

### Root-library mismatch

`internal/grade/harness.go:Build` accepts an existing executable. For a directory it tries `./cmd/` and otherwise builds `.`; it does not discover named commands such as `cmd/ch01`. On the accepted student, `go build -o <scratch>/student-root .` exits 0 but creates an **ar archive**, not a CLI. Directory grading then exits 2 with a launch permission error. This is a build/launch boundary, not a Chapter 4 behavior failure.

Build the existing accepted command explicitly, without converting the public root into `package main`:

```sh
cd /absolute/path/to/candidate
go build -o /absolute/path/to/candidate-cli ./cmd/ch01
```

The Chapter 4 CLI has a second obstacle: it passes its positional `path` both to `grade.Build(path)` and to `gradeedition2.Chapter4(bin, path, consumer, example)`. `consumerTopology` treats its `core` argument as a directory containing `go.mod`. Consequently a valid explicit executable runs the jobs/inherited checks but fails source topology because it is treated as the module directory. A consumer `-consumer-case` JSON file does not fix this identity mismatch.

There is a supported **Go API** composition already present:

```go
example, err := gradeedition2.ReadConsumerCase("") // existing default shout fixture
report, err := gradeedition2.Chapter4(binaryPath, coreModuleDir, consumerModuleDir, example)
```

A thin driver inside the scratch main module called this unchanged API and wrote the existing `grade.Report` JSON. Its exact source is `ch04-entry-evidence/driver.go.txt`. This is not a new grader/framework or a claimed existing public CLI flag. The audited CLI has no equivalent separate core-source argument for this layout; concurrent adapter edits are outside this entry review.

**Smallest proposed adapter:** retain existing positional directory/executable semantics and add an explicit core-module override for Chapter 4, defaulting to the positional path. Feed it only to the already-separated source/topology parameter, never infer the candidate core from a consumer replacement. Validate directory and explicit-binary routes, and wrong-core rejection. This audit makes no adapter edit. An entry recipe must also name the real `./cmd/ch01` build target and both required grader invocations; illustrative `go build .` commands in Chapter 4 do not justify changing the accepted library architecture.

## Existing coverage: useful, but not complete

Source anchors: `internal/grade/ch04_harness.go`, `internal/grade/ch04_checks.go`, `internal/gradeedition2/ch04_consumer.go`, `internal/gradeedition2/ch04_topology.go`.

The jobs/inherited portion has 100 points; consumer checks add 60, for **160**, not a standalone 100-point acceptance verdict. Check JSON `passed` as well as score/exit. Existing checks cover:

- corrected Edition 2 Chapter 3 parity and its inherited Chapter 2/1 checks;
- three distinct, increasing job handles before execution for `read_file`, `list_directory`, and `run_command`; completed status, output locator, byte count and exact small-result persistence; `wait_for_job` not itself becoming a job;
- applied/recorded per-call cwd, its absence on the next call, and refusal of a missing directory;
- early running status, finite wait completing, a subsequent finished wait returning without error or long delay, and final output/exit status persisted;
- pattern wake on an interactive prompt and a deliberately delayed echo that catches stale prompt wake-ups; real `send_input` and clean exit;
- actual Delve interaction under a PTY (missing Delve is a failure, not a skip);
- serial explicit kill followed by a long wait returning promptly with killed rather than normal-completion status, reason logging, and helper PID dead after the session; analogous shutdown event and helper PID death;
- full large output persisted for both command and `read_file`, bounded inline output with locator and total byte count, plus explicit 2048-byte command cap (with reporting headroom);
- one-shot delay and pattern overrides, default restoration, explicit-delay precedence, and consumption notes on job and inline `tool_limits` paths;
- implementation import edges through the required `internal/common` hub, an external module's public import and actual resolved candidate core identity, declared custom-tool description/object schema, exact successful correlated tool text, and a parseable timestamped diagnostic.

Concrete limits of that coverage:

1. **Edition 3 request lifetime is absent.** Chapter 4 reaches only the old Chapter 1 checks. The reference still constructs `http.Client{Timeout: 60 * time.Second}` (`solutions/ch04/internal/engine/engine.go:36`). Score alone cannot preserve the inherited no-whole-request/read-deadline obligation.
2. **Pending output cap is untested.** `tool_limits` fixture checks delays/patterns and explicit delay precedence, but never a pending `max_output_bytes` applied to a file/custom tool. The paired mutation below tests this exact omission. Invalid-number/regex failures and consumption by failed operations also are not exercised there.
3. **Cursor evidence is partial.** The delayed echo detects stale prompt matching on that path, but the finished wait only checks non-error/short duration, not absence of replay. No exact two-wait empty-tail assertion, discarded-middle non-replay case, or file/custom-tool cursor-boundary case. Large-output byte bounds and locator checks do not establish exact head/tail or sub-cap/long-line semantics.
4. **Concurrency/kill evidence is partial.** The serial post-kill wait does not prove notification of an already-blocked waiter. Helper PID death is checked, but no spawned descendant distinguishes group kill from killing just the helper. No late file-result overwrite check, blocked-file containment case, malformed-regexp cleanup/no-leak case, or abandoned-waiter case is composed into this grader.
5. **PTY/base-directory evidence is partial.** Delve is a strong actual-terminal exercise, and cwd override/nonpersistence/missing-directory checks exist. They do not establish exact `TERM=dumb`, 50x200 sizing, newline normalization, a shell `cd` not leaking into a later fresh command, or embedded per-Agent base-directory independence. The documented command-string interface must remain; no raw-argv API is required.
6. **External callbacks are exercised, not fully supervised.** The consumer test checks declaration, exact correlated result, and timestamped diagnostic, but not a job envelope, saved full output/metadata, running handle, shared cap/cursor behavior, blocked callback, builtins retention, or public shutdown/late-completion behavior. Timestamped output does not prove the running tool follows the actual parent chain to the root logger. Inspect registration validation, setup-versus-in-flight mutation and per-Agent isolation through focused tests/source review.
7. **Selected-model policy regression remains a separate obligation.** Existing Chapter 3 parity does not replace its independently proved selected/unselected model behavior, no-side-effect refusal, history/result pairing, bounded correction, and source-API contract. Do not introduce a universal narration gate or undocumented environment switch to make a grader happy.
8. **Source topology checks have bounded scope.** They inspect core package import edges and the consumer dependency graph, including canonical resolved core identity and positive public imports. They expressly defer real composition-root ownership, parent interfaces, common-package behavior, and absence of mutable application globals to source review. Credit the implemented graph checks rather than calling them a direct-requirement-only scan.

The historical reference already contains narrowly targeted public/job tests in `solutions/ch04/ensemble/jobs_test.go` (file cap/cursor/redaction; wait override; PTY cwd; blocked waiter/process group; consumed pattern; blocked file; shutdown). Their existence is not a grader pass or student proof, but it demonstrates that these gaps can be addressed with focused fixtures/tests instead of a new framework.

## Finite controls

Exact command receipts and outputs are adjacent under `ch04-entry-evidence/`. All controls are against scratch exports, not the student working tree.

| Control | Actual outcome |
| --- | --- |
| Build exact HEAD grade-edition2, grade-edition3, corrected thin API driver | exit 0 |
| Build historical Chapter 4 reference at its real module-root main target | exit 0 |
| Build accepted predecessor explicitly with `./cmd/ch01` | exit 0; executable identified in hash manifest |
| Build accepted predecessor with `go build -o <file> .` | exit 0, but `file` reports **current ar archive**, not executable |
| Audited `grade-edition3 -ch 4 <reference-bin>` | exit 2, unsupported chapter |
| `grade-edition2 -ch 4 -json <student-module-dir>` | exit 2 at process launch (permission denied); no student Chapter 4 behavior verdict |
| `grade-edition2 -ch 4 -json -consumer <reference/ch05> <reference-bin>` | exit 1, **135/160**, `passed:false`; only `star-topology` fails because executable path is not the core module directory |
| Existing `Chapter4(reference-bin, reference-core, reference-consumer, default-case)` API | exit 0, **160/160**, `passed:true` (8.373 s) |
| Existing E3 `-ch 1 -json` on that same reference binary, corrected 180 s outer bound | exit 1, **100/100 but `passed:false`**; both `generation_headers` and `generation_body` fail, each receiving zero complete answers/usage lines (120.042 s combined) |
| Existing `TestJobFileCapCursorAndRedaction`, original reference, `-count=1` | exit 0; real matching test, not a no-tests run |
| Build isolated pending-cap mutant | exit 0 |
| Same separated Chapter 4 API on mutant binary + its own core/consumer | exit 0, **160/160**, `passed:true` (8.690 s): measured coverage gap |
| Same focused existing test on mutant, `-count=1` | exit 1 at `jobs_test.go:54`, `cap/notice`: expected 128-byte pending cap is silently replaced by default 16 KiB |

The one mutation inserts `base.MaxBytes = defaults().MaxBytes` in `internal/jobs/jobs.go:Consume` after pending settings are consumed and before explicit arguments are applied. Delay, pattern, consumption notice and explicit output-cap override remain intact. This isolates pending-cap semantics rather than deleting all truncation. The mutant consumer resolves its own mutated core, not the original reference. Before/after source hashes and patch: `mutation.json`, `pending-cap.patch`. This is a **non-detecting grader negative control paired with a detecting contract test**, not evidence that the mutant is acceptable.

Audit setup failures are preserved separately: initial driver wrong-module compile error; initial 120-second outer lifetime timeout (SIGKILL, inconclusive); initial mutation-anchor assertion. The corrected remaining stages ran once, without repeating already-recorded controls. `audit-setup-notes.md` also records missing `rg` and the observed owned Git pager interruption. None is relabeled as an application failure. The final runner completed and all behavioral outcomes above are finite.

### Executable identities (SHA-256)

- Edition 2 grader: `3e24f12519fdaf875dfc86738880e016ce546a8de2500962563a239fffc7ac74`
- Edition 3 grader: `e7057c765859dd4719d66c9ff0d5329c436418933b451f9e20813c7e29697cab`
- Thin API driver: `8a9da8e7bd02a128d1ded511319ead15a833d1d289f6de26144fd47a3a027369`
- Historical reference: `c4754399d1427f33534fb1c9afc32942c03244521e50a7d0991e104c7cc69eb7`
- Accepted predecessor command: `24e0e4093bc7ff0ceffab2f2a863277649207f9d41e49d2bee7a34c704608bf1`
- Pending-cap mutant: `b3dfeaeb1302f641965181132af4949493f9f82f513ca2de20b2ddb6f2ad28fe`

## Smallest valid acceptance gate and next step

1. Before launch, approve/fix the executable-versus-core invocation seam and document the supported explicit target recipe. Do not force a root `main`, change the accepted public package, or silently grade a historical default.
2. At candidate freeze, record source/core/consumer revisions and binary hashes; require the existing Chapter 4 jobs + corrected inherited parity + external consumer report **and** Edition 3 Chapter 1 lifetime report to pass against that identical binary. Two zero-point lifetime checks mean a numeric score alone is insufficient.
3. Require focused contract evidence and independent source review for the measured pending-cap hole and unproven cursor, blocked waiter/process-group/late-result, PTY/base-directory, and registered-tool supervision behavior. Reuse the existing harness/API and smallest targeted tests. Keep inherited selected-model policy checks from Chapter 3 and source review for arbitrary larger request deadlines. Do not broaden into speculative redesign or demand undocumented helper signatures/configuration.
4. Once entry alignment is accepted, launch a fresh Chapter 4 student from exactly `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, under the current student-first Edition 3 procedure. Code acceptance, model-authorized live evidence, independent frozen review, and later author/text phases are future work; none occurred here. Current explicit model restrictions continue to apply. This audit needed no live service or credentials.
