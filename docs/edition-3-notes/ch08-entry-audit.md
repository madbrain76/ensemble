# Chapter 8 finite entry audit

## Decision

**Not ready for student launch.** This completes only the authorized read-only entry audit. There were no source/grader repairs, no Chapter 8 student run, no browser work, no provider-key/paid calls, no children, and no inherited Chapter 7/6/5/4 parity execution. All owned jobs ended before completion reporting.

### Findings, highest impact first

1. **Blocking — no Edition 3 Chapter 8 publication entry.** `cmd/grade-edition3/main.go:24-34,58-60,112-122` advertises/accepts only Chapters 1/4/5/6/7 and has no Chapter 8 dispatch. Built help exited 0; `-ch 8` with an executable sentinel exited 2 without launching it. The chapter opening still prints historical `grade-edition2`/`solutions/ch08` commands; those cannot stand in for a current publication/student entry. This audit deliberately does not fix either.
2. **High — persistence/echo is not execution-policy acceptance.** The contract at `chapter-08.md:109-124,318-330` requires real saved budget use, turn-start capture, subsequent-turn changes and restart. `internal/gradeedition2/ch08.go:81-237` only changes `max_tool_rounds` in JSON and compares snapshots. It sends no tool workload. The pinned student has root module `ensemble`, optional module `ensemble-gui`, immutable Config, public default 200 and construction-time negative rejection (`internal/common/common.go:285-305`, `internal/common/actor.go:9-10`, `internal/agent/agent.go:86`); mutable next-turn settings are not already exposed through its Agent interface. Translate that boundary deliberately from accepted Chapter 7, without forcing historical method names or importing a reference implementation. Historical `TestSavedBudgetEntry` is source-only evidence in this audit.
3. **High — required full initial snapshot has a proven blind spot.** `ch08.go:143-159` requires four preference fields but omits budget completeness; `grade/ch09_checks.go:65-85` accepts any nonempty initial snapshot. The isolated mutant that omitted only initial `max_tool_rounds` passed **all four settings checks**. Later explicit positive/zero budget patches are checked correctly; the result is narrowly an initial-snapshot hole, not a claim that false/zero are wholly untested.
4. **High — failed acceptance/publication paths are not graded.** Chapter 8 requires decode/write failures to preserve old state and avoid success broadcasts (`247-310`); existing driver exercises successful sparse updates and restart only. Historical solution code/tests contain relevant validation, write/rename and active-budget evidence, but those tests are not invoked by this grader and were not run here. These remain source/integration gates, not assertions of a current student defect.
5. **Required non-automated gate — UI and inherited ownership are outside the settings score.** Pane resizing, multiple artifact tabs/focus, accessibility, system theme, actual speech/font consumers and visual quality need future browser/source evidence. Chapter 7's public-core/optional-client and no-stall/identity rules carry forward. The settings control did not prove them and no browser was used.

## Verified identity and preservation

Private manifest: `docs/edition-3-notes/ch08-entry-manifest.json` (not edited or published as student input). Main HEAD `2cebe5ae76b47eabfd1f1cb7e8bef95c98b7c20e`; student HEAD/peeled `edition-3-ch07-r1` `463996b765631d31923a5726036a34942a67cefd`; student tree `5e95dcdd716f7c25c96cc3ebcde13d000bb38670`. The student remained clean. Before/after receipts verify 19 manifest inputs, 120 student files and 25 protected tracked changes, plus a separate 465-file baseline covering grader/adapter/historical Chapter 8/Edition 3 solution sources. The only deliberate exception in the pinned input group at final check is the authorized append to `docs/edition-3-carryover.md`; its original bytes remain a prefix. No other input or protected file changed. Final receipts and complete inventory accompany this report.

Full prerequisites read before controls: root AGENTS; both required skills; Edition 3 execution/readiness/student/carryover; full workflow; full Edition 2 carryover; voice; chapter-writing procedure; Chapter 8; private manifest. The accepted Chapter 7 chapter contract was then read in full for carry-forward. Source inspections and missing-path corrections are recorded separately from executed tests.

## Historical implementation assessment

Reusable source-level driver behavior: build the actual optional GUI module (using `./cmd` when present); start a fake local HTTP provider; isolate settings working directory; two subscribers; full-snapshot comparison against sparse patches; false/zero presence checks; fresh subscriber; real process stop/restart in the same workdir; bounded readiness/read/shutdown. The existing `TestChapter8Settings` confirms this narrow historical control still passes.

Important limits: `settingsChecks` is unexported, uses the historical GUI launch convention and expects `/ws`/`subscribe` plus static `web/gui` path; none licenses a new public fixture API. It slices four old Chapter 9 checks and populates their original-scenario result inputs alongside its stronger comparisons. Thus the old “theme survived restart” detail understates the current all-expected-field comparison; the report does not infer theme-only persistence coverage from that detail.

`Chapter8` itself adds `ch7-parity` and calls `Chapter7` (`ch08.go:18-30`). Chapter7 calls Chapter6 and its own GUI wire/replay/trace/pause/dependency checks; Chapter6 calls Chapter5; Chapter5 calls Chapter4; Chapter4 retains earlier executable checks and external-consumer checks (`ch07.go:18-30`; `ch06.go`, `ch05.go`, `ch04.go`). Publication Chapter 7 passes separate executable, core, GUI, consumer, workflow, stream exercise, declared models and case mappings. A future Chapter 8 wrapper must preserve those public identities and failure-before-launch rules, not reuse Edition 2's source inference or a solution default. None of this inherited chain was run in this audit.

## Actual controls

| Run | Result |
| --- | --- |
| One serial focused package command, 60s package deadline | exit 0, both packages passed; publication invalid-input/help/prelaunch tests + historical settings test only |
| Publication binary build / owned settings driver build | both exit 0 |
| Publication help with sentinel | exit 0; Chapter 8 absent; no sentinel launch |
| Publication Chapter 8 with sentinel | exit 2; usage only; no sentinel launch |
| Initial-budget-omission mutant, settings-only | exit 0; 4/4 checks pass — demonstrated coverage gap |
| No-persistence mutant, settings-only | exit 1; restart check fails, first three pass — sensitivity confirmed |

Only two mutants ran, once each. Both used temporary fixtures, fake/local HTTP, no inherited credentials, and bounded owned process groups. No control timed out or left a live process group. Disk remained above 300 MiB. Exact logs, argv, sources, exits and cleanup receipts: [controls](ch08-entry-controls/README.md). Detailed contract ownership: [matrix](ch08-entry-matrix.md).

## Recommendation and inventory

Ask Bill to authorize only the next bounded publication/settings readiness work described in [next steps](ch08-entry-next-steps.md), with a fresh reviewer/authorization before any student launch. Do not edit Chapter 8 prose or source as an incidental follow-up to this report. Do not interpret historical reference success as student acceptance.

Outputs: this audit; `ch08-entry-matrix.md`; `ch08-entry-next-steps.md`; `ch08-entry-controls/` (README, before/after/final-state receipts, exact logs and exit files, reproducer scripts, small fixture sources, hashed inventory); one appended section in `docs/edition-3-carryover.md`. No readiness edit was needed. The private manifest is an existing input, not an audit output. See `ch08-entry-controls/inventory.json` for every control file and `final-state.json` for exact permitted worktree additions and cleanup.
