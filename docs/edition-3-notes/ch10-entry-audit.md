# Chapter 10 entry audit — persistence boundary, not implementation

## Verdict

**The existing focused persistence oracle is substantive and mutation-sensitive. Chapter 10 is still unsupported by the Edition 3 CLI. No adapter, student, chapter, grader, or historical solution was changed. This audit does not authorize implementation or Chapter 10 closure.**

The single full exported `gradeedition2.Chapter10` run returned **85/100, failed**: all six persistence checks passed (85/85); the inherited `ch9-parity` check failed (0/15) because the audit's clean environment omitted the installed Delve executable from PATH/GOPATH. This is an **audit environment/setup limitation**, not a persistence failure or proof that the historical reference fails corrected Chapter 9 behavior. The full composition was not repeated. A future separately authorized adapter-preparation task must establish a correctly provisioned full-composition positive control; this report cannot supply that missing green result.

The exact existing `TestChapter10Persistence` then ran against unchanged owned copies of the grader and historical Chapter 10 solution: **baseline pass → one compiling effective-config mutant fail → exact restoration pass**. Both actual CLI and GUI command clients are exercised by that test; only the provider network is fake. The mutation failed `save-shape` and `default-load` for both clients on saved `config.system_prompt` versus actual wire configuration. No timeouts or compilation failures were used as mutation evidence.

A non-exhaustive sample of real full-run save files also demonstrates an important boundary: event sequences continue correctly, but the historical reference stores repeated `artifact_id: 1` across restarted replies. Text replay/request equivalence does **not** prove restored artifact display identity. No browser or live GUI was used or claimed.

## Authority, prerequisites, and read-only scope

Task: `/tmp/cadence-ch10-entry-task.txt`, 7,549 bytes, SHA-256 `8db6000f036e00013183a9e0ddd61b475cf66b11f21c38d3d7509c6188f173f7`.
Governing manifest: `docs/edition-3-notes/ch10-entry-manifest.json` (manifest-pinned entry scope). Main checkout: `/Users/bill/projects/ensemble`. Accepted student: `/Users/bill/projects/ensemble-edition-3-student`.

Before the first control I fully read the manifest and bounded task, AGENTS, ensemble-coding skill, Edition 3 code-reviewer skill, all four current Edition 3 policy documents, `book/edition-3/chapter-10.md`, `docs/edition-3-notes/ch10-execution.md`, the complete workflow, inherited carryover, voice guide, chapter-writing procedure, and artifact identity investigation. The workflow and carryover were read in successive line ranges, not searched as substitutes. Truncated output was reread in smaller ranges; wrong inferred paths were corrected to actual manifest/source paths. An initial attempted hash verifier used an incorrect manifest key and failed; the successful keyed-map verifier is the evidence, not that attempt.

All **18 manifest input files**, **29 protected tracked files**, and **132 accepted-predecessor files** matched pinned size and SHA-256 before controls (`input-verification.json`). The prerequisite rows include actual byte and line counts. The accepted tree is the clean Chapter 9 predecessor at commit `251ce88fb1677897dbcc6b33661096c5cb39cd8e`, tree `a2add41c306d567c05d98ae760911bc13c18b75a`; it was inspected, not graded as an unfinished Chapter 10 submission. Historical directions to start old chapters are subordinate to current Edition 3 policy and this finite task.

Also read: existing Chapter 10 grader/test, Edition 3 adapter and tests, historical CLI routing, Chapter 9→8→7→6→5→4 composition, shared `grade.Build`, GUI launcher, reference persistence/actor/history code, accepted public API and CLI lifecycle, module/example declarations, and identity allocation. No private score expectations were placed in student-visible prose. The historical book Chapter 10 was additionally read during an early path-discovery mistake; it was not treated as current authority.

No child agents, paid inference, browsers, browser setup, real credentials, copied private agent history, shared cleanup, commits, student edits, or existing-file edits were performed. Runtime source mutation was confined to a separately owned copy. The allowed report/evidence directory and a named local helper directory were the only main-tree additions. Global Go caches were reused, never cleaned; network module lookup was disabled.

## Evidence and execution boundary

Portable evidence is under `docs/edition-3-notes/ch10-entry-audit-evidence/`. `inventory.json` is the final closed inventory: report plus every other regular evidence file, exact relative path/bytes/SHA-256, excluding itself under the stated non-recursive rule. It inventories no symlink or executable binary. Absolute scratch/source locations in receipts describe this run, not portable prerequisites. Compiled binaries are represented by hashes and receipts; scratch is not part of the portable evidence.

Key files:

- `input-verification.json`, `reference-source-pins.json`, final read-only verification: input, original-source, and predecessor identity.
- `launch-environment.json`, `command-receipts.json`: exact allowlisted parent environment, argv, cwd, deadlines, PID, elapsed time, exit, log path, and invoked executable identity.
- `00-go-version.stdout.txt`: `go1.25.4 darwin/arm64`.
- `full-api.go.txt`: exact exported-API helper; no grader assertion changes.
- `run-controls.py`, `run-focused-correction.py`: exact sequential drivers, including the initial selector mistake, not cleaned-up fictitious scripts.
- `04/05` logs and `routing-marker-result.json`: compiled unsupported-route results and positive marker evidence.
- `07-full-chapter10.stdout.txt`: actual complete JSON score; stderr identifies the CLI build source/path.
- `14/15/16` focused logs, `mutation.patch`, `mutation-source-identity.json`, `restoration-result.json`, `focused-grader-copy-identity.json`: valid green/red/green sensitivity sequence and exact restoration.
- `explicit-binaries.json`: SHA-256/size of adapter, helper, marker, test executable, and successfully built mutant CLI/GUI. `observed-transient-binaries.json` contains four sampled `grade.Build` submissions; it explicitly does **not** claim exhaustive child binary coverage.
- `observed-processes.txt`: sampled child command lines, not an exhaustive process trace. `observed-save-index.json`, nine `observed-saves/*` files, and `save-summary.json`: real fake-provider save artifacts sampled while the unchanged grader owned its lifecycle. These are not a complete phase capture or raw HTTP transcript.
- Pre/post git identities/status and diff digests: existing unrelated dirty work is verified without exporting its 11.7 MB binary diff. No source-ownership claim rests on an absent diff export.

Every top-level control was sequential. Builds had 180-second deadlines, full API 600 seconds, and focused test runs a 150-second Go deadline inside a 170-second process-group deadline. No deadline fired. The runner handles timeout by terminating only its own process group, then waiting. Launch env was rebuilt rather than copied from the shell: owned HOME/TMPDIR/GOPATH, fixed system/Go PATH, explicit caches, local toolchain, `-mod=readonly`, `GOWORK=off`, `GOPROXY=off`, `GOSUMDB=off`, and loopback proxy bypass. Application launches add only grader-declared dummy provider keys, loopback fake endpoints, declared models/skill selections, and ordinary local paths. No real provider credentials were inherited.

The clean environment had a consequential omission: `/Users/bill/go/bin/dlv` exists, but PATH excluded that directory and GOPATH deliberately pointed at owned scratch. The full score's sole failing leaf was `debugger: dlv was not found on PATH or in $(go env GOPATH)/bin`. This result is retained as a red substrate/setup control, not relabeled parity success.

## Actual finite controls

| Control | Actual result | What it establishes / does not |
|---|---|---|
| Compile `./cmd/grade-edition3` | Exit 0 | Actual current adapter binary, not `go run` ambiguity. |
| Compile marker; invoke it directly with its output path | Exit 0; `DIRECT_POSITIVE_MARKER` written | Fixture really executes. Marker bytes/hash unchanged during refusal controls. This marker is not exhaustive process instrumentation. |
| Compiled adapter `-ch 10 -json`, no target, cwd marker module | Exit 2, usage only | Unsupported/no-target route refuses, rather than selecting an implicit historical solution. |
| Compiled adapter `-ch 10 -json MARKER_DIR` | Exit 2, usage only | Explicit marker target also refused. Source `cmd/grade-edition3/main.go:58` rejects chapter 10 before `grade.Build`. No claim of a usable Chapter 10 route. |
| Exported `gradeedition2.Chapter10` once | Exit 0 with failed JSON report, 85/100; about 15.94 s | All six real persistence checks pass; inherited parity is red only with the recorded missing-Delve detail. CLI exit 0 here means report generated, **not grade passed**. |
| Initial baseline/mutant test selectors | Both exit 0 with `warning: no tests to run` | **Invalid evidence.** Neither exercised persistence; retained and explicitly excluded. |
| Corrected exact baseline | `TestChapter10Persistence` runs; PASS, about 0.66 s | Focused real CLI/GUI persistence positive control. |
| Mutant CLI and GUI explicit builds | Both exit 0 | Mutation preserves compilation of both real clients. |
| Corrected exact mutant | Exact test runs; FAIL, exit 1, about 0.61 s | Causal failures for `save-shape` and `default-load`: effective saved system prompt differs from wire, in core and GUI. |
| Exact restore and corrected exact rerun | Original file hash restored; PASS, about 0.61 s | Red follows mutation, not permanently damaged fixture or grader changes. |

The initial runner stopped rather than silently accepting unexpected mutant success. Inspection found the wrong selector `^TestPersistenceLifecycleThroughCLIAndGUI$`; the unchanged actual declaration is `^TestChapter10Persistence$`. The correction script required its `=== RUN` line, preserved the invalid logs, restored the copied file before baseline, and ran the intended sequence exactly once. Neither routing nor the full composition was repeated. All jobs were awaited to completion.

Full API arguments were concrete, not defaults: core `solutions/ch10`, GUI `solutions/ch10/gui`, consumer `solutions/ch10/ch05`, workflow `solutions/ch10/ch06`, streaming exercise `solutions/ch10/ch07`; public shout consumer case copied from the accepted declaration; `claude-sonnet-4-6` GUI/stream model; explicit media map Anthropic `claude-sonnet-4-6`, OpenAI `gpt-5.4-mini`, Gemini `gemini-3.8-flash`. The reference module/registry supports those names. Persistence itself deliberately runs its own Anthropic `fake-model` phases and current-runtime OpenAI `gpt-5-course` phase. External declarations are inputs, not guessed opaque helper internals.

## Current grader surface and composition

`internal/gradeedition2/ch10.go:21` exports `Chapter10`. It first calls corrected `Chapter9`, which composes corrected Chapter 8, 7, 6, 5, and 4, then runs `persistenceChecks` and appends a 15-point `ch9-parity` summary. Nested failure details preserve the lower failing leaf. `cmd/grade-edition2` exposes the historical route; `cmd/grade-edition3` has supported routes 1 and 4–9 only. Tests of a lower exported API do not create an Edition 3 CLI surface.

`persistenceChecks` runs the CLI, builds the supplied GUI directory with `grade.Build`, then runs that command with `-serve-only -listen 127.0.0.1:0`. Results are merged by check ID: failure in either client fails the shared check, with client-prefixed diagnostics. This exercises the actual GUI process's persistence startup/EOF/shutdown path, **not a browser or rendered DOM**. `grade.Build` resolves a supplied file/binary/package and builds directory targets into owned temporary `course-grade-*` directories; omitted targets/incorrect module assumptions must not be hidden in a new adapter.

`TestChapter10Persistence` builds `../../solutions/ch10`, calls precisely `persistenceChecks`, and fails every false check. It deliberately does **not** include `ch9-parity`; its green result therefore cannot erase the full composition's red result.

The historical reference public package is `ensemble/ensemble`, its CLI is the root module, and auxiliaries are `ch05`, `ch06`, `ch07`, `gui`. The accepted student exposes `ensemble` at module root, CLI `cmd/ch01`, separate GUI `gui/cmd`, and public examples `examples/ch04-consumer`, `examples/ch05-workflow`, `examples/ch06-stream`. The accepted public root owns Agents; Agents own History, Engine, and tool/lifecycle state. Root `Shutdown` stops admission, joins Jobs, shuts down Agents, and closes observations; public request replies, not lossy observations, remain the completion authority. Existing lower topology checks use resolved import/module graphs, but source review explicitly still owns proof of root/parent construction and absence of mutable globals. No source-layout or public-API mismatch was dynamically proved by this full run; those are distinct preparation obligations, not invented failure explanations.

## Persistence contract map: assertions, status, and limits

| Contract / points | Existing exact boundary | Evidence and limitation |
|---|---|---|
| `save-shape` / 15 | Public top-level `config`, `as_of`, `context`, `log`; non-null snapshot; nonempty, strictly increasing positive log Seq; last Seq equals anchor; model/vendor/effective prompt/tool declarations compared with actual captured request; no dummy credential/endpoint leak | Full run green for both clients. Mutation directly proves effective prompt sensitivity. Saved config is compared to wire, not merely supplied raw config. This does not prove all secrets or every malformed semantic schema are rejected. |
| `default-load` / 20 | Fresh cwd first startup; one prompt per process with EOF drain; successive no-flag restarts retain FIRST/SECOND/THIRD markers; original saved log event JSON remains unchanged and first new Seq is exactly old anchor+1 | Green. These are real independent process starts, not a reused in-memory mock. Exactly-once word checks plus wire comparison are stronger than an LLM saying it remembers. |
| Explicit save override (inside default-load) | Seed only `chosen.json`; `--save PATH` loads it, retains marker conversation, changes its bytes on exit, and does not also create cwd `save.json` | Green through actual clients. Does not establish many-Agent file isolation or all path-error/atomic-write cases. |
| Current config versus old metadata (inside default-load) | Restore Anthropic history into current OpenAI `gpt-5-course`, primary `CURRENT_RUNTIME_PROMPT`, tools exactly `load_skill`, `read_file`, `unload_skill`; old conversation remains, current config wins; save compared again against wire | Green baseline; red under mutant. Demonstrates boundary, not broad three-vendor rich replay equivalence. |
| `replay-equals-snapshot` / 20 | Set only top-level `context:null`; send identical NEXT prompt; actual next HTTP body must equal full-snapshot baseline **byte-for-byte**; resulting numbering checked | Green. Saved log event equality uses decoded JSON; request equality is raw bytes. No claim that save formatting itself is fixed. |
| `tail-applied-once` / 15 | Old `context` plus old `as_of`, newer retained log; require real nonempty tail; next request equals full snapshot and words appear once; post-save Seq continues after newer anchor | Green. Distinguishes normal restore from replaying all old facts twice. |
| `log-not-needed` / 10 | Keep full context and anchor, replace log with `[]`; next HTTP request equals baseline; post-exit first new Seq exactly old anchor+1 | Green. Captured snapshot-only resave has just Seq 10,11,12 with `as_of:12`, proving this historical implementation does not restart numbering at `len(log)+1`. |
| `bad-save-refused` / 5 | Separate syntax-invalid `{"config":` and JSON `null`; real nonzero process exit (not deadline/spawn failure), stderr names save file, no fake HTTP request, exact original bytes unchanged | Green. Does not enumerate all negative semantic version/anchor/event combinations. Both invalid inputs were sampled unchanged. |
| `ch9-parity` / 15 | Full corrected prior-chapter composition, independent of persistence | Red due clean-env missing Delve. No claim that inherited behavioral checks or external-consumer topology all passed. |

Relevant source: `ch10.go:91–156` process runner/save parser, `:164–194` sequence checks, `:196–350` restarts and splice variants, `:377–413` save shape/effective wire metadata, `:448–473` override, `:476–514` current runtime, `:519–548` malformed refusal. Check weights/titles are `internal/grade/ch11_checks.go`. The existing focused test checks all six results for both clients; there is no exact numeric return score in its PASS output.

## Identity revalidation: three different concepts

1. **Durable event sequence/anchor.** Historical `solutions/ch10/internal/history/history.go` assigns `h.asOf+1`; restore stages a separate history, retains saved events, applies only events above a snapshot anchor, and seeds the next sequence above the maximum retained/anchored position. Full-run samples have Seq 1..9 or 1..12, and snapshot-only 10..12. The accepted Chapter 9 History presently continues from its last retained event; Chapter 10 must introduce anchor-safe behavior for empty/trimmed logs rather than copy a `len(events)` assumption.
2. **Parser/progress correlation.** A runtime/response-local parser part ID correlates deltas with finalization. It is not by itself a durable history order or display address. Historical `runtime.displayID` maps parser-local IDs into a runtime counter, but the runtime starts at zero and `Load` restores History without restoring that counter.
3. **Artifact display identity after restore.** The historical note targeted the root `agent/` application on a prior snapshot. Current source inspection still finds a zero-initialized runtime counter there; this audit did not rerun that note's browser experiment or inspect its old private diagnostics. Separately, **current `solutions/ch10` executions** in this audit produced repeated assistant `artifact_id:1` in normal, snapshot, log-only, tail, override, and current-config retained histories (`save-summary.json`). These actual saved records support a present historical-reference allocator concern. They do not prove what any browser currently renders.

The accepted student differs: `internal/agent/actor.go:621–639` seeds `displayID` from restored message/response/tool event identities before its owner loop; deltas and finals share `displayFor`. Therefore the old investigation is **not** evidence that the accepted student has the same current defect. Chapter 10's new snapshot-only/log-trimmed path must nevertheless preserve allocator authority beyond retained event scanning. Preserve separate response-local correlation and stable producer-owned display identity. Legacy duplicate-ID replay needs an explicit strategy if supported; inventing new IDs independently for each chunk or clearing history hides, rather than solves, the problem.

No new browser/identity probe was launched. This revalidation used current source and the already-running real persistence control's artifacts, within the bounded task. Text-only oracle success is explicitly insufficient for replayed card identity, live delta/final routing, reconnect ordering, cross-Agent identity, mixed-kind artifacts, interrupted partials, or legacy-collision migration.

## Remaining source-review and GUI obligations

Before claiming an Edition 3 Chapter 10 pass, source review still must establish:

- Storage location is chosen by application/root policy, with correct per-Agent ownership and isolation. History owns durable facts/projection; Engine owns provider translation; GUI merely adapts public APIs. No global mutable registry, cwd-dependent core singleton, or external callback bypasses the owner chain.
- Load happens before work admission; every invalid load is fatal without deferred overwrite or provider side effects. Encoding/write/close/rename errors propagate; shutdown drains admitted work and joins producers before save. A clean EOF happy path does not prove interruption, active tools, retained jobs, disk failure, or concurrent shutdown safety.
- Replay uses the same reducer as live append, is deterministic, preserves existing facts/neutral parts and usage, and never resends a request or executes historical tools. Stable sequence high-water survives context-only saves, malformed anchors, trimming, redaction, and future multi-Agent isolation.
- Effective metadata describes actual prompt/model/vendor/tools yet never overrides current runtime policy or serializes credentials/endpoints. Rich provider/tool data, opaque provenance, and current skills/visibility need their own review/tests; the minimal text fixtures cannot establish them all.
- Public external consumers remain viable with the accepted root API, root-owned shutdown, current example declarations, and the real GUI command. No requirement to adopt historical internal filenames or package layout should be smuggled into the adapter.
- Display IDs remain stable and collision-free across restored snapshots and live streaming, including partials and tool events. A future narrowly authorized real GUI/browser acceptance check must observe old replayed content unchanged while a new partial streams and finalizes into its own artifact. No such visual acceptance happened here.

## Smallest next task / completion limit

A reasonable **separately authorized** next task is bounded adapter/grader preparation, not student implementation: expose Chapter 10 only after validating mandatory explicit target/core/GUI/auxiliary arguments; retain pre-build rejection and no historical defaults; preserve the six substantive oracle properties and corrected prior-chapter composition; use the accepted public layout rather than force historical internal structure. First preflight Delve under the same allowlisted environment and establish the one missing full-composition positive control in that task. Keep the focused effective-config mutation regression and add only contract-justified source/identity safeguards, not score tricks.

The historical reference is a useful persistence control, not a canonical identity implementation. Do not copy its runtime-local artifact allocator into the student. Do not drop prior-chapter parity to turn 85/100 green. The two empty test selections must remain classified as invalid; the final green/red/green sequence is the valid focused evidence. No chapter route, student progression, live operational closure, or paid/browser work is authorized by this report.

Read-only and process completion facts are sealed in `final-verification.json`; the inventory defines the exact deliverable set. All reported control jobs finished, and no unfinished child agent exists.
