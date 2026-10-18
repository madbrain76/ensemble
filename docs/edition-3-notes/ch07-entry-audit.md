# Edition 3 Chapter 7 — independent entry audit

## Verdict

**Public behavioral contract: ready. Student launch through the present Edition 3 CLI: not ready until the narrow invocation gap is resolved.** The existing Chapter 7 grader API works without grader changes. No student Chapter 7 attempt exists, and this is not candidate acceptance, author review, or permission for paid calls.

The smallest measured successful route is a private command calling the existing explicit API, with a separately built core executable, its module root, the optional GUI module/model, and inherited modules/models. The unchanged historical positive passed **100/100, `passed=true`, exit 0**. The supported historical Edition 2 CLI also passed that matching flat fixture. Current Edition 3 rejects `-ch 7`; Edition 2 lacks a separate `-core` parameter and couples its positional build target to core-module discovery. Do not flatten the accepted library-root / `cmd/ch01` architecture to satisfy that older CLI.

Required entry work belongs to the coordinator: provide and independently accept a thin, explicit-input launch route (or equivalent public invocation) preserving executable/module separation and existing checks. No new grader algorithm, diagnostic protocol, framework, hidden acceptance floor, manuscript rewrite, or student implementation fix is established as necessary by this audit. Any adapter's own implementation/testing is outside this frozen report.

## Scope, inputs, and authority

Audit performed 2026-10-17 as a fresh independent reviewer. Fully read current `AGENTS.md`, coding/reviewer skills, all four Edition 3 procedure/readiness/student/carryover documents, historical workflow and Edition 2 carryover, voice and chapter-writing procedure, **all current Chapters 4–7**, both Chapter 7 grader files, and both CLI entrypoints. `ch07-entry-evidence/read-inventory.json` pins exact files and lengths; `read-ledger.txt` records actual ranges and truncation recoveries, including corrected endpoint annotations. Source summaries did not substitute for those reads. Edition 3 student-first rules override historical author-first ordering.

Pins, reverified after all tests:

- Main HEAD: `23c75ede7944d9a912ae74d3f683354d163f0c86`.
- Student: `/Users/bill/projects/ensemble-edition-3-student`, clean HEAD `f595c787675e2a6343053ec42ea6c7e1e02a930e`, tree `6fb36168020a22b8993412812987338aa4b6e010`.
- Immutable predecessor tag: `edition-3-ch06-r1`; tag object `84d0bce5f2525b406affa64fb1a8eb70711b1e06`, peeled commit equals student HEAD.
- Chapter 7: 692 lines, SHA-256 `68d02d8e1170c2912676b8f162847f8fac25c87216514ed06adbd5f0ea011b47`.
- Initial root pins: `/tmp/cadence-ch07-entry-initial-pins.json`, copied as evidence `initial-pins.json`.
- At the completed pre-handoff invariant snapshot, the ten concurrent publication edits checked against `/tmp/cadence-ch06-concurrent-publication-state.json` had **all ten exact hashes unchanged**. Every source/guidance/grader file in the initial pin snapshot also matched. See `final-invariants.json`. After that check, root announced a separately user-authorized Update2 typography-only rebuild (12pt → 11pt). Publication hashes may legitimately change under that work; they are not audit/student drift. This audit did not reread or modify publication outputs after that notice.

No student/reference/manuscript/grader/guidance/source edits; no stage, commit, tag, push, provider request, credential read, browser navigation, live GUI action, service restart, or child agent. Historical `solutions/edition-3` was not student input and was not used as the positive. No source diff against original-book code is mislabeled as a student change. The predecessor was inspected only for layout and clean pins, not graded as a Chapter 7 submission; missing or stub GUI behavior there is not a student failure.

## Actual public contract and inherited obligations

### Chapter 7 scope

One Agent displayed through an **optional separate importing GUI application**, leaving the headless core usable. Browser code presents state and sends controls; the Agent owns history, control, and dispatch. The real GUI imports the submitted public core module. A separate `examples/ch07` is not automatically a GUI: the accepted student's existing one imports the framework as an external **streaming** consumer and implements its Observer / protocol loop. Agent tabs/tree/workspace expansion is Chapter 8, not this assignment; “another tab” in the pause requirement means another browser connection to the same Agent.

`book/edition-3/chapter-07.md` specifies:

- Six observation/message types with real content, correlated tool calls/results, truthful turn termination, and state `from`/`to`, not empty kind-only signals. Preserve Chapter 5's existing tool-start/tool-result facts rather than introducing a second ownership route.
- Actor-owned history and stable conversation display identities through deltas, authoritative finals, distinct model responses, and **in-process** history replay. Repeated response-local vendor IDs cannot collide in a retained conversation. Disk save/load identity enters Chapter 10; it is not a new Chapter 7 requirement.
- Atomic capture/installation of completed history, current streaming partials, and the live subscription. No replay/live loss, duplicates, stale completed partial, ordering inversion, or premature live presentation. The replay selector is fixed at capture: last hour **OR** last 100 renderable events, whichever includes more. A selected 270-event replay must be delivered, not silently truncated to a live queue's size. Transient progress may coalesce; irreplaceable replay/finalized content cannot be dropped under that excuse.
- Shared Agent-wide pause gates future tool starts, including later calls in the batch and following starts; already running tools continue and report completion. Explicit unpause from **any** browser connection releases the gate. No persistent per-tab veto, disconnect auto-release, reconnect resurrection, or stale typing/speech cause reassertion. Pause is not job cancellation or Agent interrupt.
- One reusable Artifact presentation, safe DOM construction for untrusted text, optional vetted Markdown/ANSI rendering, readable speech projection excluding opaque/private reasoning payloads, and correct tool summaries/results. Raw text must not become unsafe `innerHTML`.
- Actual speech-queue lifecycle: autoplay and on-demand behavior, no tool-result narration, long-part chunking, finals without deltas, no repeated streamed final, `onerror` symmetric with `onend`, cancellation-safe late callbacks, and effective pause causes from typing/queued/active speech rather than a checkbox that falsely claims the entire cause model.
- Actual `gui.log` inbound/outbound WebSocket traffic with timestamp, direction, and JSON; no reconstructed transcript pretending to be traffic.

The chapter's historical links are explanations and old receipts, not fresh student proof. Native-page behavior and **listening** are different evidence obligations.

### Inherited Chapters 4, 5, and 6

- Chapter 4: public framework/tool consumer boundary, real tools and PTY job ownership, bounded model-facing reports with complete output preserved, continuation and truthful job state, cancellation/lifetime distinctions, and HTTP continuation beyond the old short deadline. Preserve the accepted modular architecture.
- Chapter 5: Agent-owned control loop/history, mailbox ordering and reliable paired replies, hints wake an in-flight tool wait with an **OBSERVATION** while the underlying job continues, and later output remains reachable. Interrupt/shutdown/explicit job kill remain distinct. Undispatched advertised calls receive truthful correlated paired handling, not fabricated successful execution. The GUI must not become an alternate turn controller.
- Chapter 6: typed delta/final identity, authoritative finals, tool argument bytes, readable-versus-opaque separation, current model routing, cancellation/error truth, public streaming consumer, and unaffected CLI/framework behavior. Existing accepted provider and long-form evidence remains bound to its original predecessor; this audit does not relabel it as fresh Chapter 7 proof.

## Existing route inspection and measured invocation

The actual API is:

```go
Chapter7(bin, core, gui, guiModel, consumer, workflow, exercise, streamModel string,
    example ConsumerCase, media map[string]string)
```

`internal/gradeedition2/ch07.go` passes independent core/source and inherited inputs to historical Chapter 6 parity, then builds the explicit GUI module. `buildGUI` checks the GUI's public dependency resolves to the provided core, builds `./cmd` if present (otherwise root), and locates `web/gui`. The driver uses ephemeral loopback ports, local fake vendor responses, disposable workspaces, and the existing `--port` / `--gui-dir` interface. No extra diagnostic endpoint or candidate protocol is required.

`internal/grade/harness.go:Build` builds the target directory itself; it does not discover a nested command. `cmd/grade-edition2/main.go` derives `core := target`, or `filepath.Dir(target)` for a file, and has no independent `-core`. Thus its measured successful flat historical route does not establish compatibility with the student's library root and nested executable. We did not run a deliberately incompatible predecessor GUI or flatten its source to manufacture a result. The direct API already accepts the independent values needed by that architecture. `cmd/grade-edition3/main.go` supports Chapters 1, 4, 5, and 6 only and already has the relevant core/build separation, but no Chapter 7 branch or GUI flags.

### Private measurement inputs

Owned workspace: `/tmp/cadence-ch07-entry-auditor` (macOS canonical path `/private/tmp/cadence-ch07-entry-auditor`). The harness copied current `go.mod`, `go.sum`, `internal/`, and the two existing CLI directories privately. The positive was an exact private copy of `solutions/ch07`, **66 files**, with its own `gui`, `ch05`, `ch06`, and `ch07` modules. Fixture source is auditor-only and must not be handed to the student. Initial grader and historical fixture per-file hashes are retained; local-only executable hashes are in `local-only-artifacts.json`.

Private additions were only an explicit API command and a test exposing existing package-private `guiChecks`; portable source is retained as `.go.txt`. No existing grading implementation or assertion was changed. The direct command reads the existing case/media formats with existing helpers, emits `Report.WriteJSON`, and exits nonzero when `Passed` is false.

Build commands, from their corresponding module roots:

```sh
go build -o /tmp/cadence-ch07-entry-auditor/bin/entry-api ./cmd/entry-api
go build -o /tmp/cadence-ch07-entry-auditor/bin/positive-core .
go build -o /tmp/cadence-ch07-entry-auditor/bin/grade-edition2 ./cmd/grade-edition2
go build -o /tmp/cadence-ch07-entry-auditor/bin/grade-edition3 ./cmd/grade-edition3
go test -c -o /tmp/cadence-ch07-entry-auditor/bin/gui-checks ./internal/gradeedition2
```

All five builds exited 0. Each `*.run.json` records exact cwd, argv, times, exit, and environment overrides; stdout/stderr are separate safe UTF-8 files. `run.py` preserved the subprocess exit in the JSON even when its enclosing collection command continued to gather later runs. `TMPDIR` was owned scratch, `GOPROXY=off`, `GOTOOLCHAIN=local`; ambient credential-bearing variable names were removed. Local fake-server traffic is not paid provider evidence.

Successful explicit call, cwd `/tmp/cadence-ch07-entry-auditor/harness`:

```sh
/tmp/cadence-ch07-entry-auditor/bin/entry-api \
  -bin /tmp/cadence-ch07-entry-auditor/bin/positive-core \
  -core /tmp/cadence-ch07-entry-auditor/positive \
  -gui /tmp/cadence-ch07-entry-auditor/positive/gui \
  -gui-model claude-sonnet-4-6 \
  -consumer /tmp/cadence-ch07-entry-auditor/positive/ch05 \
  -workflow /tmp/cadence-ch07-entry-auditor/positive/ch06 \
  -stream-exercise /tmp/cadence-ch07-entry-auditor/positive/ch07 \
  -stream-model claude-sonnet-4-6 \
  -consumer-case /tmp/cadence-ch07-entry-auditor/consumer-case.json \
  -media-models /tmp/cadence-ch07-entry-auditor/media-models.json
```

The `shout` case and model map are pinned as their own safe JSON files. Historical Edition 2 CLI positive used these same explicit flags, replacing `entry-api -bin ... -core ...` with `grade-edition2 -ch 7 -json` and a **positional** positive module root at the end. Exact non-abbreviated argv is retained.

## Finite results and sensitivity

| Run | Exit | Score | `passed` | Meaning |
| --- | ---: | ---: | --- | --- |
| `api-positive` | 0 | 100/100 | true | Existing explicit API, complete Chapter 7 report; 13.83 s |
| `cli2-positive-corrected` | 0 | 100/100 | true | Supported historical flat-fixture CLI route |
| `cli3-unsupported-corrected` | 2 | not emitted | not emitted | `-ch 7` unsupported; entry gap, not student/grader behavior failure |
| `cli2-no-core-corrected` | 2 | not emitted | not emitted | `flag provided but not defined: -core`; independent-input gap |
| `replay-mutant-build` | 0 | n/a | n/a | Mutated GUI compiled |
| `replay-mutant` | 1 | 60/80 | false | Only `event-replay` failed: actual final content/identity lost |
| `replay-restored` | 0 | 80/80 | true | Exact original source restored; all four GUI checks pass |
| `pause-mutant-build` | 0 | n/a | n/a | Mutated actor compiled through GUI build |
| `pause-mutant` | 1 | 60/80 | false | Only `pause-holds-tools` failed: 3 starts before other-tab unpause |
| `pause-restored` | 0 | 80/80 | true | Exact original restored; all four GUI checks pass |

There were exactly **two behavioral sensitivity controls**. Replay changed only replay final `part_id`, leaving live output and content intact. Pause changed only the actor's dispatch-gate condition so it no longer waited. Minimal diffs and original/mutant hashes are retained. They are reviewer controls, not student implementation instructions. Both are compile-success **intended behavior detections**, not build/watchdog/disk/startup failures. The four-check reports have max 80 because they deliberately omit the 20-point inherited parity wrapper; they are not separate 100-point Chapter 7 grades. Both controls restored byte-for-byte; `restoration.json` verifies the entire positive file inventory, not just the changed files.

### Failures and qualifications retained

The initial three CLI attempts mistakenly used nonexistent `-solution`; each exited 2 at argument parsing and emitted no score/Passed. They are retained as `cli3-unsupported`, `cli2-no-core`, and `cli2-positive`, **not** counted as detections. Corrected attempts used the positional target. An inspection tried unavailable `rg`, and another initially looked for `examples/ch07` in main rather than the student workspace; standard tools and the correct student path recovered those reads. A packaging script initially treated the publication snapshot's map as a list and raised `TypeError`; the corrected script verified all pins. A milestone-handle query returned `unknown job 345`; no test was canceled or inferred successful from that. See `audit-incidents.txt`.

Tests were sequential. The API's own historical inherited parity completed in seconds and did **not** run the separate Edition 3 Chapter 1 `lifetimeCompanion` with its genuine 62-second HTTP request (located in `internal/gradeedition3/ch01_lifetime.go`). That inherited obligation remains; this is neither its fresh proof nor a duplicate Chapter 6 campaign. No genuine long gate was interrupted. No native-browser/speech test or paid campaign was run. Disk was initially about 2.5 GiB free and about 2.0 GiB free at completion; no disk/build/watchdog blocker, shared-cache cleanup, or authority to delete unrelated files was used. Only existing grader-managed disposable-workspace teardown occurred.

## Precise coverage limits for later public evidence

A score of 100 establishes the current finite assertions, not comprehensive acceptance. The driver proves actual local WS prompt/subscription, required content-bearing observations and tool output, completed-turn replay identity/content, valid actual traffic logs, and real shared pause/release behavior for its three-tool batch. It has a bounded quiet collection after a completed turn; this is **not** a causal active-replay barrier proof. It checks current state fields are populated, not every possible state transition. The inherited parity score is the existing aggregate, not fresh comprehensive Edition 3 lifetime acceptance.

Later student/reviewer evidence should use the public interfaces and real submitted page, without new hidden APIs:

1. **Owner and dependency inspection:** optional GUI imports the exact core; history and dispatch remain actor-owned; public calls/mailbox ordering, slow-observer handling, and shutdown preserve ownership. No controller hidden in browser callbacks or GUI server.
2. **Identity and active replay:** deliberate repeated vendor-local IDs across responses, in-process reconstruction, active streaming reconnect with causally fixed boundaries, finalized/partial/live ordering without loss/duplication, replay queue delivery under backpressure, and fixed last-hour OR 100-event selection including a 270-event case. No disk-restoration requirement added here.
3. **Complete pause/control matrix:** next batch and following starts blocked while active tools continue, explicit other-connection unpause releases everyone, disconnect does not release, reconnect and stale local typing/speech causes do not reassert, and hint/interrupt still remain reachable while waiting at the gate. Preserve truthful job continuation and correlated undispatched results.
4. **Actual DOM and speech queue:** real emitted page/assets, untrusted tool text safely rendered, readable/opaque projection, chunking/autoplay/on-demand/tool-result silence, final-only and delta-plus-final cases, typing causes, error completion, and canceled-generation late callbacks. Stubbed queue callbacks are deterministic fault evidence, not native speech proof.
5. **Native browser and listening:** demonstrate real page loading/render/control and actual browser speech separately; listening must assess intelligibility and omissions. Use a later isolated candidate instance only. Never navigate/refresh or restart the owner's live GUI at port 8088/current browser. Record any unavailable native evidence honestly.
6. **Inherited/public use:** accepted CLI, public consumer/workflow/stream modules, model/media routing, true job/lifetime behavior, and required source-bound provider/human-interface evidence remain independently scoped. Historical links or a fixture pass cannot stand in for fresh candidate receipts. Any paid work requires the later authorized checklist and frozen inputs.

No full browser/code acceptance campaign, extra mutation suite, or whole-book preparation is proposed at this entry stage.

## Handoff and evidence boundary

Coordinator should inspect and join this report, then own any thin adapter and its independent acceptance before launching the fresh student. Publish only the legitimate reader-facing invocation/contract clarification needed by that student, **not** the private positive, grader internals, control diffs, or historical answer code. No student implementation feedback is warranted yet.

Portable evidence: `docs/edition-3-notes/ch07-entry-evidence/`, with `INDEX.md`, `SHA256SUMS`, exact file set, source/diff control text, safe command reports, timings/exits, and source/binary hash inventories (binaries remain local only). Checksum scope explicitly excludes the manifest's own impossible self-hash. Local scratch is retained without cleanup. All owned tests/builds are done; no children or pending test jobs remain. This report is the finite completed entry audit, not a Chapter 7 validation checkpoint.
