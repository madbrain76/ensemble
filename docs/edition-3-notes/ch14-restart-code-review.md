# Chapter 14 restart — independent code review

## Entry review (read-only; before student)

**Disposition:** align the public Chapter 14 contract and existing grader before
starting the student. This is an entry recommendation, not implementation or
grader acceptance. No student exists and no production round has begun.

Reviewed the complete current `AGENTS.md`, coding and reviewer skills, shared
readiness, execution, student and carryover documents, applicable workflow,
all 976 lines of Chapter 14, Chapter 9's current instruction/disclosure policy,
Chapter 10's recovery/diagnostic contract, and both retention obligation notes.
The lean/current directions supersede historical audit/mutation/author-first
procedures. Inspected actual accepted owners and the existing Chapter 14 checks,
helpers, tests and command routing. No tests, builds, reference grades, browser,
provider, mutation campaign or network operation were run.

Entry facts: HEAD is `f3194f1c34cac9fecd48acb92bda6ec26d796aa6`; accepted main has
165 tracked files. The protected paper draft and unrelated/untracked material
were not changed. Only this entry note was written.

### Minimum public alignment

1. **Preserve accepted Chapter 9, not the obsolete frozen-SYSTEM policy.**
   Chapter 14 lines 35–36, Rule 1, the overview, §14.3 Law 1 and the
   frozen-prefix figure freeze SYSTEM and place runtime skill bodies outside
   it; §14.4 also retains the old tool-result-instruction narrative. Rule 4
   already distinguishes skill entries from acknowledgment results. Chapter 9
   instead explicitly appends loaded instruction
   bodies to the next SYSTEM, and successful supported Anthropic discovery
   uses paired reference-only tool results. Actual owners are
   `internal/engine/engine.go:Render/EffectiveSystemPrompt`,
   `internal/engine/render.go:system`, `internal/engine/tools.go` and
   `internal/tool/skills.go` under accepted main. State that skill/manual entries
   survive ordinary result clearing and continue contributing to dynamic SYSTEM;
   they are not ordinary tool-result instructions. Any stable-prefix statement
   must be limited to the unchanged base/startup material and the applicable
   fixture declaration behavior, not all SYSTEM bytes across a skill load.
   Keep provider-specific references paired to their actual enabling call and
   catalog definition; do not mix instruction/status text into reference-only
   results. Preserve ordinary meaningful results for non-reference paths.

2. **Keep `InlineTools` explicitly synthetic.** The two `*-course` rows in
   Rule 6 are local deterministic fixtures, not new Anthropic API capabilities.
   The mid-conversation addition/removal vocabulary and placement exercise can
   remain on those rows. Do not infer it for real models or let it replace
   Chapter 9's live Anthropic paired-reference convention. Replace the stale
   named live-model recommendation with supported current capability guidance.
   `internal/common/actor.go:LookupModel` is the existing exact-row table;
   `internal/engine/tools.go:supportsToolReferences` owns the live reference path.

3. **State the intentional unload advance.** Chapter 14 Rule 4 legitimately
   advances compiled-tool unload: remove that skill's grants immediately while
   its instruction body survives until the context checkpoint. This differs
   from Chapter 9's pending compiled grants; it is not grounds to freeze the old
   behavior. Preserve shared grants and existing immediate MCP teardown. Actual
   owner: `internal/skills/skills.go` (`UnloadSkill`, enabled-tool calculation,
   instructions/snapshot). Existing `chapter09_round1_test.go` includes the
   assertion that pending unload leaves rendering unchanged; that expectation
   will need a narrow Chapter 14 update, not removal of unrelated Chapter 9
   disclosure/provenance coverage.

4. **Replace event-count retention with the backup-anchored recovery contract.**
   Rule 10, §14.12–13 and settings prose currently introduce `log_retention`
   (zero means full log) and allow loss of old history. Remove that
   canonical-event cap. At a safe boundary publish valid primary and backup
   recovery state before trimming the independent journal. Keep the complete
   journal tail strictly after the **validated backup's** anchor, including
   records between backup and newer primary. If a usable backup boundary is not
   established, do not discard required records. Recovery must select/validate
   primary or backup, apply that candidate's embedded tail and remaining journal
   records exactly once in sequence, ignoring only already-covered overlap.
   Failed validation/publication/replacement must leave required recovery
   material available. Serialize retention with the actual event writer, not
   merely with other saves. Ordinary process-kill recovery is required;
   power-loss durability is not newly required.

   Full archival history and dialogue-only archives remain Chapter 10 policies,
   distinct from the recovery journal. Diagnostic retention also stays separate:
   `diagnostic_max_records`, default 200, zero unlimited, negative invalid, with
   safe startup/shutdown application. After a **successful** `micro_handoff`
   recovery save, invoke the same existing diagnostic retention operation;
   do not run a second ad-hoc trimmer or trim before failed recovery publication.

5. **Correct runnable public examples and invocation.** Chapter 14
   line 158 and its exercise/run instructions retain Edition 2 invocation
   or original `solutions/ch14` assumptions. Publish the
   Edition 3 chapter-14 invocation with explicit core and GUI paths. The manual
   primary fixture must explicitly grant `load_skill` and `unload_skill` along
   with the shown context/read tools: accepted `skills` computes grants from
   metadata, not a blanket implicit lifecycle allowlist. Do not promise these
   tools remain available under a primary that omits them.

### Existing persistence is the implementation starting point

- `internal/common/persistence.go` already defines `SaveFile`, `SavedConfig`,
  archive modes and `Recovery`. `internal/history/persistence.go:SaveFile`
  captures a boundary, encodes/validates it, rotates only a validated prior
  primary into `.bak`, and replaces primary through a closed temporary file.
  `LoadFile` independently validates fallback and publishes only a complete
  candidate. `decodeSave` handles snapshot plus embedded tail, complete-log
  null-context replay, dialogue-only constraints and semantic context equality.
  Extend these owners; no parallel snapshot framework is necessary.
- `internal/agent/persistence.go:Save` already serializes saves with `saveMu` and
  records effective prompt/capability metadata without restoring runtime secrets
  or authority. `Recovery.AsOf` currently reports the effective loaded boundary,
  including an embedded tail, not necessarily the file's raw snapshot anchor.
  Retention must derive its cutoff from validated recovery data, not assume these
  two values are interchangeable.
- `internal/agent/actor.go:saveDump` currently atomically rewrites the independent
  `LogPath` JSONL from History; `internal/history/history.go:AppendSaved` stages
  durable tool admission before publication. `internal/agent/agent.go` currently
  loads that JSONL only when no snapshot was selected. These are concrete owners
  for the Chapter 14 tail/replay/writer-serialization change. Merely retaining a
  suffix in `SaveFile.log`, or placing a mutex around `Save` alone, does not
  implement independent-journal crash recovery.
- `internal/diagnostics/log.go:Log.Retain` already shares the writer mutex,
  atomically replaces whole complete records, and has subsequent writes reopen
  the destination path. Reuse it. Public alias/wiring is `diagnostics.go`;
  callers are `internal/cli/cli.go` and `gui/cmd/main.go`. Configuration owners
  are `internal/common/common.go`, `internal/agent/persistence.go`, and
  `gui/settings.go` (plus its browser controls). Keep `log_mode` and diagnostic
  policy intact while adding context-target controls.
- New context events/entries and replay belong with existing
  `internal/common/history.go`, `internal/history/history.go` and archive
  validation; context actions/pending-call clearing belong in the existing
  Agent actor/tool path. Renderer placement remains Engine work. These are
  Chapter 14 implementation obligations, not pre-student source edits.

### Minimum grader alignment (keep 8 checks / 100 points)

Retain the existing IDs and weights from `internal/grade/ch15_checks.go`:
`skill-survives-the-ladder` 15, `micro-handoff-shape` 15,
`keep-or-stub` 10, `ladder-is-recorded` 15, `frozen-prefix` 10,
`replay-equals-snapshot` 15, `crash-recovery` 15, `total-reducer` 5. Keep thresholds, request timing,
pairing, exactly-once effects and independent side-effect checks; do not replace
these with prose-only or declaration-only evidence.

- `internal/gradeedition2/ch14.go:contextSkillOnce` (used by `contextLadder`) currently rejects the
  manual marker in SYSTEM. Replace that obsolete predicate with persistent
  instruction authority in dynamic SYSTEM across load, ladder and pending
  unload; retain immediate removal of the manual's compiled grant.
  `contextPrefix` must not compare all SYSTEM bytes across skill changes;
  retain unchanged base material, fixture-only inline placement, effective
  declarations/removals and the real MCP input side effect. Existing Chapter 9
  `skillSystem`, `skillResult`, `skillDeclarations` and their contract tests show
  accepted SYSTEM representations and paired-reference semantics—reuse rather
  than invent another live-vendor interpretation.
- `internal/gradeedition2/ch14_driver.go:contextCreateSkills` needs explicit
  primary lifecycle grants. Its reuse of Chapter 12's `guiDebugPeer` is sound,
  and `contextPrefix` already checks the correlated, successful
  `INPUT_RESULT_ONLY:INPUT_ARGUMENT_ONLY` result. Preserve that real execution
  evidence, rather than counting echoed call arguments. Remove the inherited
  expectation of automatic browser snapshot/TTS ephemera:
  accepted Chapter 12 uses explicit model calls and paired results. Synthetic
  ephemera placement can remain a renderer fixture without reintroducing
  unsolicited browser work.
- `contextReplay` compares snapshot/null-context replay, and `contextCrash`
  covers process-kill recovery; neither establishes backup-anchored retention.
  Extend these existing checks' fixtures to exercise validated backup fallback plus a
  complete independent journal tail/overlap exactly once after safe trimming;
  preserve byte comparison and no repeated effects. Keep the full/null-context
  replay case on a sufficient full archive, not a deliberately trimmed recovery
  suffix. Focused failure/retention tests belong beside the existing persistence
  tests; there is no need for a new campaign or scoring framework.
- `contextTotal` must continue proving that an unappliable old event does
  not prevent a later valid event from replaying. Its payload-less `skill_loaded`
  fixture follows Chapter 14's explicit example but intersects Chapter 10's
  strict tagged-payload validation (see clarification below).
- `ch14_wire.go` / `ch14_test.go:TestChapter14InlineGroups` already preserve valid
  split inline groups and reject intervening user content; retain that boundary.
  The integration test still builds original `solutions/ch14/gui`; it is not an
  Edition 3 candidate test. Keep original/reference use distinct from explicit
  candidate paths.
- `cmd/grade-edition3/main.go` and `internal/gradeedition3/ch01.go` currently stop
  at 13. Add the ordinary thin Chapter 14 branch, routing to corrected existing
  checks. `ch14_driver.go:contextLaunch` currently assumes a GUI-shaped binary
  (`--port`, `--gui-dir`, settings/WebSocket helpers). Separate core and GUI
  invocation explicitly; do not pass GUI flags to the core or infer sibling
  paths. Chapter 12's `guiDebugChecks` already demonstrates explicit separately
  built paths. Keep the eight-check score rather than adding an unrelated new
  parity score. No obsolete baseline/reference grade is needed before alignment.

### Consequential clarification before student

Chapter 14 Rule 8/§14.10 permits malformed/unappliable payloads to be
skipped, and the existing grader treats a payload-less `skill_loaded` as a parseable
but unappliable event, whereas accepted Chapter 10 validation requires exactly
one correctly tagged payload and rejects malformed archives before mutation
(`history.go:validate`, `persistence.go:decodeSave/eventValid`). State the intended
boundary explicitly: Chapter 14's specified recoverable old/unappliable cases
must skip with diagnostics and allow later valid events, without relaxing
whole-file JSON/sequence/candidate validation or silently accepting arbitrary
corruption. The public example and `contextTotal` fixture must agree with
that boundary. This is a narrow contract reconciliation, not permission to drop
Chapter 10 validation or discard the total-reducer check.

No other scope expansion is needed for entry. The remaining ladder, handoff,
model-row, settings, replay and effect behavior is student implementation work
under the existing public eight-check contract. Retain this reviewer context for
later candidate review; this note does not initiate that work.

## Held entry-alignment diff review — HOLD (before student)

Reviewed the complete held manuscript/adapter/grader diff and new
`ch14_contract_test.go`, with the relevant accepted GUI save lifecycle,
registry, persistence and grader helpers. No student exists; this is not
production acceptance. No tests or grades were run in this review. The existing
`/tmp/ch14-entry.LuPUZ1/{adapter.log,contract.log}` focused PASS results are
reused as adapter/predicate/orchestration evidence, not student behavior proof.
Available disk remains about 20 GiB.

### One blocking fixture finding

**`internal/gradeedition2/ch14.go:contextTotal` can pass without applying the
edited primary.** `contextChecks` clones the entire successful crash-recovery
directory into T. `contextTotal` then nulls/replaces only `save.json`, inserts
its two bad familiar events and checks that later dialogue survives. It leaves
`save.json.bak` and the independent journal untouched. With the newly required
complete backup-anchored suffix, refusing the edited primary (for example,
retaining Chapter 10's strict missing-known-payload behavior) and falling back
to the valid backup plus journal can reconstruct every prompt/reply asserted by
this check. No assertion distinguishes that fallback from successful tolerant
reduction of the edited archive. The new crash control tests do not cover this
false-positive path.

Minimal repair: make this existing full-log/null-context T fixture recover from
its edited canonical archive alone, without alternate backup/journal recovery
sources in this cloned scenario, or add equally direct evidence that the edited
primary was the candidate replayed. Preserve its five points and current public
behavior. A small focused control should show that refusing the edited primary
cannot pass by obtaining the expected dialogue from untouched recovery files.
This is not a request for a mutation campaign or another recovery framework.

### Alignment otherwise sound

- The manuscript correctly separates stable primary instructions from dynamic
  SYSTEM skill projection, and retains supported Anthropic paired reference-only
  results. Synthetic inline rows do not become a live API claim. Recorded manual
  survivors remain through unload and handoff; immediate compiled-grant removal
  is explicitly the Chapter 14 advance over Chapter 9.
- Known-event semantic tolerance is now explicit without weakening whole JSON,
  version, anchor, ordering or unknown-type refusal. Independent journal,
  backup-boundary suffix, selected-candidate embedded-tail overlap, writer
  serialization, validated rotation and separate diagnostic retention are stated
  coherently. No parallel persistence framework is demanded. Read Rule 10's
  failed-save preservation as preservation of **required recovery records**,
  not an extra all-files-byte-identical transaction: Chapter 10 can already
  safely publish a validated backup before primary replacement fails.
- Adapter 14's direct positional GUI binary/module route matches `Chapter14`'s
  actual driver and avoids irrelevant parity flags. Separately freezing/grading
  the core CLI under Chapter 1 is sufficient; the earlier suggestion of extra
  core/GUI adapter flags is superseded. No-target guards and the existing
  `grade.Build` path retain explicit selection.
- SYSTEM/manual marker predicates accept the intended projection and reject
  duplicated dialogue bodies. Removing unsolicited snapshot/TTS expectations
  matches the on-demand Chapter 12 peer; the successful correlated `gui_input`
  result, inline ordering, additions and removals remain checked.
- `contextCrash` forks before recovery mutates files, tests intact/missing/corrupt
  primary independently, compares requests and counts every prior turn once.
  Its prior-backup setup is supported by the accepted implementation:
  `gui/cmd/main.go` saves the fresh empty Agent at startup, and
  `internal/agent/actor.go` saves at normal shutdown. Thus the first clean session
  rotates a genuinely earlier backup before the second session is killed.
  Backup coverage must still mean its anchor **plus embedded tail** on recovery;
  raw file inequality alone is not a general proof of distinct effective
  coverage. The present lifecycle gives this fixture the needed earlier state.
- Eight IDs, 100 points and driver timing constants are unchanged. The new unit
  tests are appropriately small controls, not claims of implementation success.

Correction to the initial entry note: accepted
`internal/skills/skills.go:activate` seeds `load_skill`/`unload_skill` during
initial activation. Explicit fixture grants are clear and harmless, but their
previous absence was not by itself an accepted-registry dispatch defect.

**Disposition:** hold entry acceptance for the isolated `contextTotal` issue;
otherwise the held alignment is a suitably small, coherent route to Chapter 14.
No manuscript, grader, production or protected paper files were changed by this
review. Await root direction; no student review or additional work is initiated.

## Entry completion — PASS (isolated total-reducer repair)

The prior HOLD is resolved. Reviewed only the requested repair in
`contextChecks`, `contextCopyReducerInput`, `contextTotal` and
`TestChapter14TotalCannotFallback`; unchanged prior review remains applicable.
T is now a fresh directory populated only with the skills tree and complete
primary archive. Backup/journal fallback cannot supply the dialogue if that
edited archive is refused. The existing copy helper creates the needed directory;
copy failures and final archive encode/write failures stop before launching the
runner. This is a small, appropriate fixture isolation, not a new recovery API.

The focused control verifies the destination inventory despite alternate
recovery files in its source, sees the actual null-context/planted-known-event
archive at invocation, and confirms strict rejection fails `total-reducer`.
Reused the executed PASS in
`/tmp/ch14-entry.LuPUZ1/total-isolation-v3.log`; no rerun was needed. The earlier
exact-anchor edit refusal/no-tests selector run is not test evidence, and the
v2 missing-parent failure was test setup, not useful behavioral sensitivity.
Their logs remain preserved; v3's explicit parent creation fixes that setup.

**Final verdict: ENTRY PASS.** The single blocking finding is closed without
changing check IDs, points or public requirements. This accepts entry alignment
only, not student production or implementation behavior. No source, manuscript,
grader or protected-paper edits by this reviewer; only this note was appended.
No further work initiated.


## Root entry decision

Accepted the saved, joined entry PASS and personally inspected the affected public contract, adapter, predicates and recovery owners. Eight IDs/100 points remain. The actual GUI-direct API needs no unused core/parity flags; Chapter 1 remains a separate frozen CLI gate. Explicit lifecycle grants clarify the fixture, not a claim that accepted startup activation previously failed. Unit positives/negatives and total-input isolation are grader controls, not production evidence. The first missing-selector/parent failures and initial punctuation lint failures remain in `/tmp/ch14-entry.LuPUZ1`; executed adapter/contract/isolated-primary checks and final punctuation-corrected lint passed. Soft prose warnings are retained. Original history, accepted source and concurrent paper are untouched. No Chapter 14 student grade, browser or provider operation has run. Public command examples added from the reviewed route; entry acceptance does not accept a student implementation.

## Whole held-student review — BLOCKED, before production round 1

Reviewed the fresh held `solutions/edition-3/main` against accepted `ch13`, not
merely the accepted entry changes. This is a source/architecture verdict at
0/10 production rounds, **not PASS-for-validation**.

### Bounded findings

1. **Per-round stubbing revisits and destroys a previously kept batch.**
   `internal/agent/context.go:198–232` (`stubPresentedResults`) scans every
   complete batch through `SentThrough`; its only keep exclusion requires
   `KeepAt == SentThrough`. `prepareContext` calls this on every request.
   A large result can be kept on the consuming response, survive the next
   request, then be stubbed after that request's boundary advances—even if
   its response was text-only and both byte bands remain below their triggers.
   Rule 6 makes this decision about the just-consumed batch; this is not a
   license to reclassify all older results on each request. Preserve the
   once-per-presented-batch decision separately from the later byte-band
   ladder. The same transient exclusion appears in `compactBands:174`, but
   this finding does not demand permanent immunity from the ordinary ladder.
   This is a concrete source trace, not an executed behavioral reproduction.
   My earlier live message's descriptive name `stubSeenResults` was inaccurate;
   the current function is `stubPresentedResults`.

2. **GUI settings broadcasts clobber unsent Context Target edits.**
   `gui/web/gui/app.js:52–54` rewrites every preference input from `settings`;
   `applySettings:64–66` always invokes it. There is no dirty-edit/focus
   protection for the new numeric input. A user can be typing an uncommitted
   context target when another tab's settings snapshot arrives and replaces
   those bytes. The common onchange path also resets controls before server
   acknowledgement (`73–78`). The public GUI requirement explicitly preserves
   unsent edits; the Go settings test exercises persistence/live targets, not
   this browser state. Source-established overwrite path; no browser was run.

3. **Normal paired-reference rendering repeatedly hydrates the full archive.**
   `internal/engine/engine.go:47–78` calls `History.EventsRange` separately for
   each dynamically enabled tool, before checking whether its enabling result
   survived compaction. `internal/history/history.go:49–63,619–642` reads and
   unmarshals the complete save under the History mutex for every such range.
   This turns a one-event correlation lookup in ordinary rendering into
   repeated lifetime-archive allocation/I/O and blocks the actual append owner.
   It remains expensive after handoff removes the enabling result. Material
   avoidable complexity for a chapter about bounded live context: use bounded
   enabling-reference/projection state, or otherwise avoid full archive reads
   on this normal rendering path. Full archives remain legitimate for explicit
   export/recovery. Current bytes use `json.Unmarshal`, **not** `decodeSave` or
   reducer replay; no timing or memory magnitude is claimed without a probe.

4. **The new live-job cache accumulates completed identities indefinitely.**
   `internal/agent/actor.go:216–226` writes `liveCalls[CallID] = false` on
   completion/kill; there is no deletion/reset anywhere. Its sole query
   (`context.go:279`) only needs active membership. Thus each lifetime job
   leaves another otherwise useless key in resident state. Remove inactive
   membership rather than introducing a second retained job history. This is
   a small new-state correction, not a request to redesign inherited jobs.

### Architecture and inspected scope

The basic split is appropriate and should survive correction: actor-owned
policy resolves numeric redaction boundaries; History owns reduction and the
serialized actual journal/snapshot writer; Engine projects provider requests;
skills own grants; the managed diagnostic writer keeps its independent policy.
Reviewed the complete changed-owner surface: root construction and persistence
integration; Agent mailbox/turn/result/job/context paths; common facts/context
and exact model rows; History append/reduction, journal, save/load/validation;
skill lifecycle/restoration/declaration deltas; context-tool registration;
Engine SYSTEM/inline/paired-reference rendering; GUI Go settings and JS/HTML;
and the added focused tests. Adjacent cancellation, pending receipts, replay,
EOF/save and tool admission were checked where these changes cross them.

In particular, the inspected writer serializes append/rotation under the same
History lock, validates backup candidates, retains the suffix after the backup
anchor, and considers embedded-tail coverage during replay. Handoff records
follow the acknowledgement and wait for outstanding receipts before survivor
flush; diagnostic `Retain` is reached only after a successful recovery save.
Known payload semantics are checked before reduction while envelope/type and
replay-boundary validation remain separate. The synthetic inline rows are
explicit, dynamic SYSTEM/manual survivors remain separate from dialogue, and
compiled grant reconciliation excludes pending-unload owners while checking
other active grants. These are source observations, not a declaration that all
failure, restart or provider cases have been demonstrated.

Read original `solutions/ch14` policy/history only as a simplicity comparison
(round-local presentation is useful); did not adopt its obsolete frozen-SYSTEM
or recovery policies. Current guidance, student handoff, Chapter 14 and the
accepted Chapter 9/10 contracts/retention notes governed the review; unchanged
entry prerequisite reads carried forward.

### Evidence and limits

Read the focused test source and actual saved logs under
`/tmp/e3-ch14-student`. `focused-final-root-2.log` is a **failed aggregate**
(`TestChapter09EnablingEventAndRendering/writer: unloaded compiled declaration
remains`) despite its three passing internal packages. The later scoped
`focused-root-final.log`, `focused-history-hold.log` and `focused-gui.log` are
passing evidence only for their selected inputs. `compile-hold-final.log`
shows six modules compile with `-run ^$`; it is not executed acceptance testing.
No grade, broad suite, race, lifetime, browser, provider, network or behavioral
probe was run by this reviewer. Failure injection, GUI observation and complete
runtime acceptance remain unestablished. No production/grader/manuscript or
protected-paper edits; only this review note was appended. Root owns correction
scope and any later release to validation.

### Root decision: production R1 released

Root read the saved whole-review verdict and joined the retained reviewer before releasing one bounded round (1/10). Retain the owner-local implementation; correct the four findings above with focused actual-path regressions and explicit-module compilation, then HOLD. No broad/race/grade/lifetime/browser/provider operation is released. Earlier focused passes and compilation are not acceptance.

Root independently inspected the changed ownership paths and original actor/journal/tool helpers as a simplicity comparison, rejecting the original's obsolete recovery policy. Two tentative descriptions were wrong and withdrawn after exact current-byte reads: the current keep mechanism is `KeepAt`/`KeptBatch`, not a `protectedBatch` helper; archived event lookup performs raw unmarshal, not reducer replay. A tentative log-only skill-fact error concern was also withdrawn: current `RecordFact` errors are returned. No source-drift or concurrent-writer allegation follows. The supported findings concern round-local eligibility, actual settings overwrite, full-archive hydration on normal rendering, and resident inactive-call entries. Root's wrong temporary log-path lookup was corrected to `/tmp/e3-ch14-student`; it does not imply lost tests or cleanup.

## R1 held-candidate delta review — BLOCKED on one GUI ordering case

Read the updated student handoff, current correction source/tests, and executed
`/tmp/e3-ch14-student/round1-focused-go-2.log`, `round1-focused-gui-2.log`, and
`round1-compile-six.log`. Retain the earlier whole-source architecture review;
no new execution or probe was needed for this delta review.

Three findings are adequately corrected for release to validation:

- `stubPresentedResults` now considers only batches in
  `(PreviousSentThrough, SentThrough]`. The reducer advances these request
  cursors and saves/replays them. The real actor test carries a kept result
  across text-only turns and a new unkept batch, then lowers T to show that
  ordinary bands can still remove it. This is one-shot presentation eligibility,
  not permanent keep immunity.
- `Engine.Render` obtains enabling-call identity from live assistant parts'
  derived `AdmissionSeq`, never a full-archive range read. History binds that
  sequence on a matching real `ToolCalled` without mutating the canonical
  response. The matching vendor/surface call and surviving successful result
  are still required; missing/foreign/handoff-cleared pairs fall back to
  ordinary declarations. Read the save/replay test, archive-read trap and real
  load/render/handoff/foreign test. The new state is small and appropriately
  owned by the live projection.
- Completed and killed job IDs are deleted from `liveCalls`; the focused
  active-membership test covers both transitions. No extra job framework.

The production JS now preserves ordinary unsent edits, handles newer typing,
separates accepted shared settings from drafts, and clears pending submission
on refusal/disconnect. Its real-settings-frame/production-handler test is useful,
but one source-established ordering defect remains:

**Negative drafts can be falsely acknowledged by unrelated settings frames.**
`gui/web/gui/app.js:86–94` admits a negative safe integer to the existing server
refusal path, but records its pending value as `Math.max(20000, value)`.
`applySettings:69–70` treats *any* shared frame with that target as acceptance.
With accepted target 20000, type and submit `-1`; if a previously queued or
another tab's unrelated settings update arrives before `settings_error`, the
20000 frame clears the draft and replaces `-1`, although no target change was
accepted. The later refusal cannot restore the lost edit. The current refusal
test sends `settings_error` immediately and therefore does not cover this
ordering. Do not normalize an invalid negative submission into a valid pending
acknowledgement; preserve the editable draft on this refusal path. This is a
bounded correction to the same unsaved-edit obligation, not a new settings API.
No browser/probe reproduction is claimed.

**Verdict: BLOCKED only on that remaining draft case; not yet PASS FOR
VALIDATION.** The other corrections do not justify a redesign or further
campaign. The three supplied logs pass their selected checks; six-module
compile-only remains compile evidence. Earlier pre-fix and setup failures stay
qualified in the student note, not counted as acceptance/sensitivity. Node
production-handler checks are not connected-browser observation. No new tests,
builds, browser/provider/network operations or source edits by this reviewer;
only this existing review note was appended. Chapter acceptance remains root's
later decision after authorized validation.

### Root decision: narrow R2 released

Root read the R1 source/tests and actual pre-fix/final logs, then read and joined the retained review verdict. The three non-GUI owner corrections are adequate for validation, not chapter acceptance. R2/10 addresses only the remaining negative-draft/unrelated-settings ordering case above, with production-handler regression and relevant focused checks. No broad suite, grade, browser, provider, or repeated module campaign is released; unchanged R1 evidence carries. Root does not patch student production.

## R2 narrow completion — PASS FOR VALIDATION

The final R1 blocker is closed. Read the R2 student note, production
`gui/web/gui/app.js:86–96`, and the extended JS/real-settings-WS tests. The only
production correction guards pending-acceptance normalization with `value >= 0`,
with a useful explanation. Negative drafts still reach the existing server
validator, but cannot be falsely acknowledged as 20000 by an unrelated frame.
Draft preservation, legitimate nonnegative normalization, server refusal and
shared accepted settings retain their existing owners. This is the adequate
small correction; no new protocol or settings framework is needed.

Read `/tmp/e3-ch14-student/round2-prefix-js.log` and
`round2-prefix-gui.log`: they demonstrate the exact pre-fix `-1` → `20000`
draft loss, not setup failures. `round2-focused-gui.log` records both selected
Go tests and all five production-JS cases passing, including unrelated-20000
before refusal. These are existing executed results, not a fresh reviewer run.

**Verdict: PASS FOR VALIDATION.** Carry forward the whole-candidate architecture
review and the three corrections closed in R1; no source blocker remains from
this review. Prior unchanged Go/compile evidence and its limitations carry
forward. This is permission for root's authorized validation decision, **not**
operational readiness or chapter acceptance. No new execution, builds, probes,
browser/provider work, source edits or further tasks by this reviewer; only
this existing review note was appended.

## First validation stop — bounded failure classification

Read the complete failure messages and timeout stacks in the 1,647-line
`/tmp/e3-ch14-student/validation-main-race.log`, the implicated current and
predecessor owners/tests, and both new boundary-test sources and passing logs.
No rerun, probe, build, patch or student communication. The aggregate failed
assertions and then hit the Go 600-second timeout (616.34 seconds wall clock).
It contains no reported data race; this is **not** a race-clean suite result.
Downstream validation did not run.

### Production blockers — do not weaken these inherited assertions

1. **Unpublished admission blocks committed-history readers.**
   `TestCh7R3ActualAdmissionSave/inspect_saved_boundary` deliberately pauses the
   real save callback after its staged dump, then inspects `EventsRange` to prove
   the dispatch is not published and the handler has not run. The timeout stack
   puts the callback at `chapter07_round3_test.go:57` and History
   `AppendSaved:157`; the test's history read waits for the same `h.mu`.
   Current `internal/history/history.go:144–163` holds that state mutex across
   the external save callback. The predecessor releases it around the staged
   save and publishes only afterward. Chapter 7's durable-before-publication/
   effect boundary and readable committed frontier remain required alongside
   Chapter 14's actual writer serialization. Preserve serialization of writes,
   saves and retention without holding the committed-state reader mutex over
   the save callback. Increasing the timeout or releasing the test latch before
   the committed-state assertion would hide the product regression.

2. **Construction publishes false tool removals into recovered history.**
   `TestCh10R2DirectDispatchArchive` logs 16 original facts and 17 after reopen.
   Offline comparison of those logged JSON arrays shows the original 16 are
   unchanged; the extra fact removes all 13 tools. `skills.New` registers the
   catalogue through `RegisterTool`, which now calls
   `publishDeclarationsLocked`. On a recovered `ToolsInitialized` context this
   runs before `LoadInitial`/`RestoreContext` installs effective grants
   (`ensemble.go:147–169`). Thus incomplete construction is incorrectly treated
   as a live capability change. Preserve canonical IDs/arguments/anchors and
   restore real grants before publishing live deltas; do not filter the false
   removal out of the archive equality assertion. The completed and
   tool-result-prefix variants of `TestCh10R3ToolArgumentSaveAndRenderPrefix`
   fail at the same reopen comparison and traverse this same faulty path;
   no evidence here establishes argument corruption as a separate cause.

### Precisely justified fixture alignment

- `TestCh10R2BaseRelativeRecoveryIsolation` requires zero initial events, then
  exactly one event after adding a human message. Chapter 14 legitimately
  records the initial `tools_changed` fact even for an explicitly empty
  registry. Allow/assert that local initial capability fact while preserving
  the path, empty recovery-path, own human-message, no foreign history and
  no process-CWD `cr` checks. This explains the *fresh-root* assertion only;
  it does not excuse the false removal on reopen above.
- `TestRound2EphemeraTailFailureAndReplay` (`chapter02_round2_test.go:173`)
  requires total `Context.Dialogue` length three and the third entry's first
  text part to be `"new question"`. The three causal entries are the old human
  message `"old"`, the assistant's `"old reply"`, and the new human question.
  Initial capability history now adds a `KindTools` survivor. Update the
  assertion to identify/preserve all three causal `KindDialogue` entries and
  their payloads, retaining the separate ephemera-clearing/idle assertion,
  foreign-provenance suppression and replay/request equality. The logged
  “error removed causal dialogue” is not proof of deletion: it comes from the
  total-length/third-entry-payload condition. This corrects the earlier numeric
  description only; classification and verdict are unchanged.
- `TestChapter3LoopOrderedBatchProviders`, `TestChapter3GeminiParameterProjection`
  and the final cross-Agent check in `TestChapter3Round1Registry` still expect
  10 defaults. The observed 13 are exactly the prior ten plus Chapter 14's
  `compact_context`, `micro_handoff`, `keep_tool_results`. Update those exact
  default-set/count expectations, retaining provider schema conversion,
  explicit-empty/selected tools, identity and isolation checks. The registry
  test's earlier expectation is already 13; its later 10 is inconsistent.
- `TestCh10R1FailedExitSaveKeepsDiagnostics` deliberately corrupts primary and
  backup after startup, then insists on the phrase “neither recovery candidate”.
  It now receives an archive-prefix JSON decode error first because full
  canonical history is read from the archived prefix. That is a legitimate
  save refusal, not evidence of diagnostic truncation. Assert refusal without
  that internal wording dependency and retain the exact 200 old + 10 fresh
  diagnostic-line assertion. The failed run stops before that byte assertion;
  it therefore does **not** establish that this particular retention scenario
  passed. CLI source still gates final `Retain` on successful shutdown/save.

The new focused history and Agent boundary logs pass their explicit concurrent
save/recovery, backup, and successful/failed handoff-diagnostic scenarios. Reuse
those results for unchanged inputs; they neither exercise the blocked staged
admission reader nor validate construction against recovered capabilities.
**Status: validation stopped; production corrections required for the two
blockers, with bounded fixture alignment above.** Earlier PASS FOR VALIDATION
was a release to execute gates, not acceptance or authority to bypass this stop.

### Root validation-stop decision and R3 release

Root read the failed assertions, actual staged-save lock path, startup registry ordering, and full retained classification. A structural comparison of the saved archive diagnostics found all original 16 facts unchanged plus one false all-13-tool removal. The ten-minute timeout is not cured by extending the deadline or releasing the barrier before its committed-view assertion. R3/10 is released only for those two production corrections and the specifically justified inherited fixture alignments; focused checks and affected-module compilation precede another HOLD. No full-suite retry, grades or operations are released. The review note's initial two-entry shorthand was corrected to the actual three causal dialogue entries before this decision. The failed aggregate remains failed; the new boundary passes retain their limited scope.

## R3 held-candidate delta — PASS FOR VALIDATION

Refreshed the current role/coding/shared guidance and read the full R3 student
section. Reviewed the three changed production owners, actual Agent admission
and snapshot callbacks, new reopen test, strengthened History boundary test,
authorized fixture diffs and all six named R3 evidence logs. Earlier whole-source
review and closed R1/R2 findings carry forward; no new execution or probe.

Both validation-stop production blockers are closed:

- **Committed-reader boundary:** History consistently takes its writer
  reservation before the committed-state mutex for append, staged append,
  snapshot save/rotation and load/repair. `AppendSaved` prepares an immutable
  staged suffix under the state mutex, releases that mutex for the external
  callback, and retains the writer reservation until journal admission and
  publication finish. Callback/serialization failure leaves the committed
  frontier unchanged; queued writers cannot overtake staging. The actual Agent
  callback is `saveDump`, not snapshot `Save`: it writes the provided dump and
  does not acquire Agent `saveMu`, call back into a History writer, or create a
  reverse lock order. Read-only History access remains available during staging.
  The unchanged Chapter 7 barrier and its serializer/duplicate-save-failure
  cases pass. The strengthened real-file staged/plain append test also overlaps
  snapshot saves, reads the old frontier inside callbacks and verifies staged
  bytes, primary recovery and missing-primary backup-tail continuation.
- **Recovered capability construction:** `skills.New` uses the same validated
  insertion helper without publishing declarations. Normal `RegisterTool`
  still locks, inserts and publishes the genuine live delta. Root initial
  activation/restoration remains responsible for effective grants before use.
  Default, selected, empty and recovered-skill reopen cases now preserve exact
  canonical facts, projection and effective menus; a later registration emits
  precisely its one addition. The original direct-dispatch archive and
  completed/tool-result-prefix argument comparisons remain intact and pass.

The fixture changes are narrowly justified: fresh-root isolation permits exactly
its empty initial capability fact while retaining path/recovery/CWD and own-message
checks; the causal assertion checks **all three** entries, actors and payloads,
with ephemera/replay checks retained; default menus require all 13 names and keep
schema, explicit selection and isolation checks; failed-exit-save accepts the
archive-prefix refusal and still requires exact 200-old + 10-fresh diagnostic
bytes, now also requiring both corrupted save files to remain untouched.
No inherited admission barrier or canonical/argument equality assertion was
relaxed. The owner changes are clear and proportionate—two mutex responsibilities
and one insertion helper, not a new persistence or registry framework.

Evidence qualification: `round3-admission-reopen.log` remains aggregate **exit 1**.
Its unchanged regression cases pass; its new test initially used an exclusive
range boundary incorrectly and the wrong fixture primary name. Those setup
failures are not sensitivity evidence. The corrected new test passes in
`round3-startup-menus-fixtures.log`. History boundaries, Agent/CLI diagnostic
retention and GUI settings logs pass their selected cases. Six-module compile
logs are compile-only evidence. The original `validation-main-race.log` remains
failed/timed out; R3 focused checks do not constitute a full race/vet pass.

**Verdict: PASS FOR VALIDATION.** No remaining blocker from this delta review.
Root may decide the renewed validation sequence; this is neither chapter nor
operational acceptance. No full suites, grades, binaries, browser/provider work,
source edits or further tasks were initiated by this reviewer. Only this existing
review note was appended.

### Root R3 validation release

Root read the held writer/callback and startup-publication paths, exact fixture diffs, new reopen/concurrent-writer tests and all named focused logs, then collected the retained reviewer's saved R3 PASS FOR VALIDATION. The committed-reader fix preserves writer serialization without holding the state lock through the external save callback; constructor population no longer emits partial catalogue facts. Existing admission and canonical replay assertions remain intact. The first new-test aggregate remains failed, with corrected setup demonstrated separately. Three production rounds are used. Root released a fresh sequential validation from root race/vet through GUI/application modules, then one frozen CLI/GUI pair, Chapter 14, and only on its parsed/process success the real Chapter 1 lifetime gate. No automatic fixes, grade retries or operations are authorized. This is validation release, not chapter acceptance.

## R3 GUI validation stop — fixture classification only

Read `r3-validation/gui-race.{stdout,stderr,status}`, the implicated complete
GUI tests, checkpoint/journal loading and relevant startup, terminal and settings
owners. Root main race and vet are fresh passing gates (38.80s / 0.77s); GUI race
exited 1 in 20.49s on the listed assertions. No data-race warning is reported in
that GUI output, but the GUI suite did not pass. No subsequent gate, rerun,
probe, build or source edit by this reviewer.

**No production blocker is established by these failures.** The required narrow
fixture repairs are:

1. **`TestCh10R1SnapshotPrefixTailWS`: derive the real snapshot boundary.**
   The seed now contains the legitimate initial `tools_changed` fact, followed
   by the assistant call. Consequently `prefix.AsOf` is 2, not 1. Writing
   `s.AsOf = 1` beside that Context is invalid; `contextValid` correctly rejects
   it. Select the tail by actual event sequences strictly after the captured
   prefix boundary, not `Log[1:]`, and use that boundary consistently. Do not
   relax checkpoint validation or turn all variants into full-log replay.
   For `full_carry`, `Dialogue[0]` is the typed initial tool survivor, not the
   assistant call. Locate the exact canonical `native-call` tool-call part and
   retain its required `DisplayID == 7` assertion; do not accept an arbitrary
   nonzero/new GUI ID instead.

   **Isolate the recovery inputs as part of this same fixture repair.** The seed
   also leaves `<save>.journal`. Rewriting only primary to `context_only` still
   lets `LoadFile` restore the original tail from that journal *before* the test
   explicitly appends it. Construct the selected archive in a fresh owned save
   location (or otherwise isolate the owned backup/journal side inputs). Assert
   the intended context-only boundary before adding its real tail. This keeps
   context-only and embedded-tail paths meaningful rather than silently testing
   independent-journal fallback. Keep every existing live-frame call/result/job
   association, stable/new-ID assertion, reconnect check, save and exact-produced-
   archive reopen check. In this failed run the first two modes never loaded;
   `full_carry` stopped at its stale canonical-index assertion before final
   save/reopen. Those downstream portions are not passing evidence yet.

2. **`TestCh7R1TerminalUnboundedInput`: wait for the actual new message.**
   Startup already gives `View().Through > 0` and `Events[0]` is capability
   history. Thus the wait can succeed before the asynchronous terminal reader
   admits any input, and the subsequent first-event assertion necessarily reads
   the wrong fact. Capture the pre-input frontier and wait for the newly admitted
   human message containing the exact entire trimmed long input; select that
   message by type/identity, not array position. Preserve the >scanner-limit
   payload and exact bytes. `Terminal` still uses `bufio.Reader.ReadString`,
   submits through `SubmitLive`, and the actor records the human message before
   paused provider dispatch. Admission can be observed without requiring an
   offline provider completion or reducing the payload.

3. **`TestSettingsPersistenceAndFailedSave`: compare the unchanged baseline.**
   A fresh Agent now legitimately has startup capability history. The final
   `len(Events()) == 0` assertion does not identify a settings leak. Capture and
   compare exact canonical history before and after the settings operations;
   do not merely allow extra facts or replace the assertion with a nonzero
   count. Keep accepted broadcast/persistence/live-policy checks, failed-save
   refusal, no rejected broadcast, unchanged previous settings file and reconnect
   settings. The owner saves settings before changing policy/broadcasting, and
   no history append is present on that path.

**Status: stopped at GUI validation; bounded fixture-only alignment is justified,
not a new production correction round.** The new recovery-isolation requirement
above preserves existing mode coverage; it adds no production API or chapter
feature. No claim is made that corrected fixtures or downstream gates pass until
root authorizes and observes them. Prior R3 source verdict and unchanged evidence
remain scoped as recorded.

### Root GUI fixture acceptance and validation continuation

Root read the exact three-file fixture diff and selected-test logs. The revised snapshot fixtures derive actual anchors, isolate seed recovery side inputs, bound the context-only display archive and assert the intended loaded boundary before appending the tail. Exact native-call canonical ID, WebSocket/reconnect/reopen checks remain. Terminal waits for the complete new human message beyond its baseline; settings compares exact unchanged canonical events. The first selection failed compilation on the nonexistent DialogueRecord.Seq field and ran no test bodies; corrected SourceSeq selection passed all three snapshot modes and terminal/settings cases. No production changed; rounds remain three. Root accepted the fixture repair and released continuation from GUI race/vet through remaining modules and conditional grades, carrying the unchanged fresh root race/vet passes. First unexpected failure still stops execution; no automatic fixes or operations are authorized.

## R3 client-script stop — obsolete default mode; scoped resume recommended

Source/log-only classification. Read the actual default failure, full
`client_round1_test.cjs`, current Chapter 8/13 controls, relevant page rendering,
replay/delta branches and actual GUI regression owners. Fresh GUI race/vet and
`ch08_round1_test.cjs` pass in `r3-validation-2`; default client invocation exits
1 at line 35 (`conversation` has one card, expected two). No new execution,
fixture/source edits or automatic production round.

**No current product defect is demonstrated.** The legacy first block expects
both thinking and tool cards in `conversation`. Accepted Chapter 8 workbench
rendering sends tool calls/results to `tool-cards`; one thinking card in Chat is
correct. Direct comparison with accepted `ch13/gui/web/gui/app.js` confirms the
card/replay/partial-final paths are unchanged; Chapter 14 changes there are
context-target settings/draft handling. The legacy mock also uses `parent`
rather than `parentElement` and appends existing nodes without real DOM move
semantics, so repairing only the first count is not a sound plan for treating
the entire old block as a current workbench gate.

**Recommended bounded resume, for root authorization:** use the existing
`node client_round1_test.cjs chapter12` mode and current
`node ch13_correction_test.cjs`, retaining the already-valid current Chapter 8
and actual GUI results rather than rerunning them. The `chapter12` guard is
already in the script; importing its fixture from the Chapter 13 file does not
execute the old block. No broad Chapter 7 fixture repair/import campaign or
production change is justified by this stop.

Coverage qualification is important:

- The selector skips the **whole** legacy Chapter 7 block: safe-text/edit-diff
  DOM assertions, compact details, silent replay, several old speech/Continue
  sequences, long native-queue content and direct partial/final DOM fencing.
  Only its later trusted-versus-synthetic typing/whitespace assertions still
  execute in `chapter12` mode; this is not an equivalent pass of the old block.
- Fresh Chapter 8 controls exercise current workbench layout/follow behavior,
  current Continue authority, cancellation/stale callbacks and Read Aloud/HUMAN
  holds. The Chapter 13 file supplies current disconnected/whitespace submission,
  speech-vs-typing release, duplicate native callbacks, segmented disclosure,
  target outcomes and later-human-edit protections. It does **not** directly
  replace the skipped adversarial DOM and stale-delta assertion pair.
- Existing actual GUI tests retain live/replay wire identity, partial/final
  content and job/result matching, accepted speech content and real WS Continue
  protocol. Their fresh suite pass is transport/server evidence, not a browser
  DOM or native-hearing claim. Source still uses `textContent` through `element`
  for user/provider/edit-diff text; `apply` ignores a delta after a nonpartial
  final for the same card key. These unchanged production branches support
  carrying forward source confidence without inventing a new executable gate
  for this unrelated Chapter 14 settings delta. The skipped exact DOM assertions
  remain unexecuted in the proposed scoped invocation.

The unavailable historical Chapter 13 `client-round1.log` is **not evidence** of
which command previously passed. This recommendation rests on current source,
existing explicit mode selection and observed current gates—not a reconstructed
historical receipt. Root must observe any resumed selected commands before
crediting them; no chapter or operational acceptance follows from this note.

### Root JS invocation decision

Fresh GUI race/vet and Chapter 8 controlled JS passed. The default inherited client script failed its old single-pane count assertion; root and retained reviewer classified the current source and unchanged renderer, without executing probes. Root released only its existing `chapter12` mode and current Chapter 13 correction suite, then the outstanding modules and conditional grades. This skips the entire legacy Chapter 7 block, including exact adversarial DOM and direct late-delta assertions; those are not newly credited. Source inspection and actual GUI wire checks are not equivalent browser/native evidence. The unavailable historical command receipt was not reconstructed. No production or fixture edit is warranted by this invocation mismatch; rounds remain three, prior passed gates carry, and the failed default run is preserved.

## First Chapter 14 grade stop — grader asset setup defect

Source/log classification only. Read `r3-validation-3/grade14` stdout, stderr,
status and parsed result; current `contextChecks`/launcher, neighboring grader
asset conventions, public Chapter 14 invocation contract and actual GUI server
asset validation. The saved report is **0/100, false, exit 1** with eight checks,
but the independent scenarios all stop before readiness on missing
`<owned-temp>/gui/index.html`; recorded request counts are zero. Dependent checks
then report their prerequisite sessions did not complete. This is not evidence
of failure of their retention/recovery predicates.

**Confirmed grader setup defect, not a student production defect.**
`contextChecks` creates an empty directory and supplies it as `--gui-dir`.
`NewServerWithSettings` correctly requires existing non-directory entries for
`index.html`, `app.js` and `style.css` before listening. Public Chapter 14 says
the grader supplies GUI assets. Both an explicit frozen GUI binary and a GUI
module built by the adapter receive this same bad directory; changing target
form or rebuilding the held GUI cannot fix it. Production asset validation
must remain intact.

**Recommended smallest repair:** populate that already-owned fixture directory
with a minimal inert `index.html`, `app.js` and `style.css`, checking directory
creation and each write and returning a setup/Fatal error on failure. Keep it
local to the existing Chapter 14 fixture setup, with an explicit comment that
these are WS-only static fixtures, **not** browser, page behavior or shipped-
asset quality evidence. The launcher only connects to `/ws`, reads settings,
sends stdin prompts and observes protocol/recorded effects; it never retrieves
or executes any of these page files. Static fixture content therefore cannot
satisfy or bypass an assertion in the eight checks.

This follows the existing `--gui-dir` server/launcher contract without adding
an adapter framework or new argument. Chapters 9/10/12 use actual module asset
directories because those routes possess a GUI source-directory identity;
Chapter 14 deliberately accepts a standalone frozen executable. Inferring a
source/assets location from that binary would be incorrect. An explicit actual-
assets argument could work but is unnecessary surface for these WS-only checks.
Do not import browser/parity checks or pretend inert files validate real assets.

Preserve the actual frozen GUI executable, startup flags, settings path, all
eight IDs/predicates/weights, request/effect assertions, crash forks and deadlines.
No student source change, rebuild or production correction round is justified.
The original failed grade receipt stays failed; Chapter 1 and later operations
remain withheld until root authorizes continuation. No probe, grade, patch,
child, browser/provider work or new execution was performed by this reviewer;
only this existing note was appended. Root owns repair and any renewed grade.

## Static-asset fixture repair — PASS

Reviewed the current replacement in `contextChecks`, appended
`contextStaticFixture`, its focused test, and saved
`/tmp/ch14-static-fixture.uTPsRD/test.log` / `test.exit` (executed PASS, exit 0).
No new execution by this reviewer.

The helper creates the three required nonempty inert files in the existing
owned directory, propagates mkdir/write errors, and labels the files accurately
as launch prerequisites rather than UI/shipped-asset evidence. `contextChecks`
returns a setup/Fatal error before launching if preparation fails. The test
checks exactly three readable nonempty files, refusal of a non-directory parent
and propagation of an asset-write failure. This is the recommended small fix;
no new adapter or asset-location assumption is introduced.

**PASS for this grader fixture delta.** The launcher, frozen student owner,
flags and eight behavioral checks remain unchanged. The focused helper test is
not a passing Chapter 14 grade or browser test. Keep the first 0/100 false/exit-1
receipt; no student rebuild or production round is needed. Root may authorize
the one renewed grade and conditional Chapter 1 as planned. No grade, rebuild,
probe or further task was initiated here; only this existing note was appended.

### Root static-fixture acceptance and grade release

All remaining local gates passed with their recorded compile-only and JS-selection limits; root inspected exact statuses, outputs and frozen hashes. The first Chapter 14 grade remained 0/false/exit1: the grader's empty static directory prevented startup and yielded no model requests. Root missed this prerequisite at entry. Root added only checked inert fixture files, tested setup success and directory/write refusals, and obtained the retained reviewer's saved PASS for that delta. These WebSocket checks do not validate shipped browser assets. Root released one new grade on the unchanged frozen GUI, with Chapter 1 on the unchanged frozen CLI only after parsed/process success. No rebuild or production round, and no operations are authorized.

## Renewed Chapter 14 grade — one fixture defect, one product blocker

Read the current `r3-validation-4/grade14.stdout`, actual session/fixture owners,
Engine projection and disclosure paths, model-row metadata, projection test and
public rules 1/2/4. Source/log classification only; no probe, rerun, patch or
production round initiated. The saved outcome remains **85/100, false, exit 1**.
Six checks pass, including ladder, handoff, keep, snapshot/replay and crash
recovery; Chapter 1 remains withheld. Two failures classify differently:

1. **`total-reducer`: confirmed grader setup defect.** All sessions receive
   the same outer `tmp/skills` through the launch closure. Session C does not
   own `C/skills`. `contextCopyReducerInput` wrongly tries to copy that nonexistent
   subtree, so the semantic-reduction scenario never launches. Its unit fixture
   copied the same incorrect ownership assumption and therefore missed this.
   Small repair: ensure owned T exists and copy **only** C's complete `save.json`
   there, with errors propagated; retain the existing edited/null-context input.
   Shared skills are already supplied by the unchanged real launch closure.
   Update the unit setup to the actual sibling/shared-skills layout, requiring
   only the primary in T and still proving planted known events reach the runner
   and strict refusal fails. No whole-directory copy, backup or journal fallback
   is permissible: fallback would reopen the earlier false-positive loophole.
   This failure says nothing yet about the product's total reducer.

2. **Inline-model instructions: demonstrated product/public-contract blocker.**
   Prefix request 2 follows a successful manual load and has zero manual bodies
   in SYSTEM, where the public rule requires exactly one. The source explains
   it: `internal/engine/context.go:projectPrompt` appends active skill bodies to
   the base only when `!features.InlineTools || skill.ActiveSeq == 0`; for inline
   mode it instead retains runtime `KindSkill` as a dialogue System entry.
   `render.go:system` specifically excludes that inline Skill entry, and the
   Anthropic dialogue renderer carries it outside the top-level SYSTEM field.
   This is not a missing registry instruction: `Skills.Snapshot` contains active
   bodies, but Engine intentionally uses recorded `ActiveSkills` when present.

   Rules 1/4 require recorded surviving skill bodies in SYSTEM on subsequent
   requests, including replay, tool clearing, unload-pending survival and handoff.
   Synthetic `InlineTools` changes **declaration placement only**: startup tools
   stay fixed and `KindTools` deltas follow the complete correlated turn. It
   does not exempt skill instructions from dynamic SYSTEM or permit a duplicate
   dialogue instruction body. Correct the projection at its current owner;
   retain that distinction, stable primary, startup declarations and inline
   delta ordering. Do not weaken the grader's manual-SYSTEM assertion or treat
   the synthetic row as a new live API. Keep the separate real-vendor Anthropic
   full catalogue/genuine paired reference-only path intact, without manual prose
   or acknowledgement mixed into a native reference result.

   Existing `internal/engine/context_test.go` includes both course rows but only
   requires each body once **anywhere** in serialized output. That assertion
   allows the wrong lane and is not contrary evidence; a public placement check
   must distinguish SYSTEM from dialogue while retaining its declaration checks.
   The prefix predicate stops at the placement failure, so its later MCP/inline
   checks are not credited by this run.

**Status: BLOCKED for the public instruction-placement defect, plus bounded
root-owned fixture alignment.** Existing local gates/frozen artifacts remain
what they were; neither the six passing checks nor this diagnosis authorizes a
new round or acceptance. Root must inspect and decide correction/resume. The
private grader is not student input; no findings were sent to the student.

## Primary-only reducer fixture repair — PASS

Reviewed only the current root-owned `contextCopyReducerInput` delta, updated
`TestChapter14TotalCannotFallback`, and
`/tmp/ch14-primary-fixture.rRI0vu/test.{log,exit}` (executed PASS, exit 0).
The helper checks creation of owned T and copies only `save.json`, propagating
read/write errors; its comment accurately leaves shared skills to the launch
closure and excludes recovery side inputs. The test now creates sibling source
and shared-skills directories, requires exactly one primary file in T despite
backup/journal files in source, and retains the planted known-event/null-context
runner assertion and strict-refusal negative result. This closes the wrong-
ownership fixture defect without reopening fallback or weakening the predicate.

**PASS for this fixture delta only.** The earlier 85/100 false/exit-1 report
still records a root-owned setup mistake: total reduction was not exercised.
This unit result is not a product reducer pass or renewed grade. Student R4 and
its public instruction-placement correction are outside this review. No new
execution, source edits, private student input or further task was initiated;
only this existing review note was appended. Await root's held-R4 request.

### Root renewed-grade triage and R4 release

Root read the 85/false/exit1 report, shared skills launch closure, primary-copy helper and actual Engine projectPrompt. Six check passes are scoped; prefix checks after the early instruction failure are not credited, and total reduction never launched. The retained reviewer independently confirmed the split. Root releases production round 4 of ten: every surviving dynamic manual belongs exactly once in SYSTEM, including synthetic inline-tool rows; declaration lanes, recorded survivors, native paired reference-only results and canonical history remain intact. Focused field-aware tests and compilation only precede a new held-source review. Root owns the separate wrong-skills-layout fixture defect: copying only the primary into the isolated directory, shared launch skills unchanged. Its existing no-fallback test passes on the corrected sibling layout at /tmp/ch14-primary-fixture.rRI0vu; root read and accepted the retained reviewer's saved fixture-only PASS, then joined it. No grade retry, Chapter 1, or operational work is authorized at this boundary.

## Production R4/10 delta — PASS FOR VALIDATION

Reviewed the full R4 handoff, saved exact delta, actual two changed Engine
production files, field-sensitive context test and new real-Agent lifecycle
test. Carried the unchanged architecture/ownership review; did not repeat entry
review or the separately accepted primary-only grader fix. No execution by this
reviewer, and no changes except this appended note.

**The public instruction-placement blocker is closed in source.**
`projectPrompt` now adds every recorded active/pending survivor body to SYSTEM
for all models, and removes `KindSkill` only from projected wire dialogue—not
History. The two `render.go` exceptions now reserve inline treatment for
`KindTools` alone. This is the smallest correction at the existing owner, not
new state or a provider-specific workaround. Stable primary instructions,
synthetic frozen startup declarations and chronological declaration deltas
remain separate from dynamic skill instructions. The genuine Anthropic
catalogue/admission-pair/reference-only branch is unchanged.

Evidence inspected:

- `round4/prefix.log` demonstrates the actual pre-fix Opus-course failure in
  SYSTEM-versus-dialogue placement, both directly and through load, clearing,
  pending unload, handoff, replay and restart. Other covered rows passed.
  This is meaningful sensitivity evidence, unlike the separately reported patch
  anchor/delta-record setup mistakes, which changed no source and get no credit.
- `focused.log` passes the field-sensitive provider matrix, real Agent/local
  HTTP lifecycle, live-pair provenance and enabling/rendering controls actually
  listed in that log. Tests require primary/manual exactly once in native SYSTEM
  and absent from dialogue; inspect complete two-call/result ordering and inline
  addition/removal groups; preserve startup declarations; retain the manual
  after clearing, grant withdrawal and handoff; compare canonical replay and
  snapshot/journal restart requests; and assert rendering leaves canonical
  History bytes unchanged. The lifecycle uses real registry and local file
  dispatch, rather than only manually constructing a rendered request.
- Real Anthropic lifecycle assertions require genuinely paired reference-only
  content and ordinary fallback after clearing. `reference-controls.log` also
  passes the explicit native-reference, ordinary/failed/missing-result fallback,
  direct Dispatch/Invoke/manual load and startup renderer controls. The two
  unmatched optional names in the first selector are **not** counted as tests;
  the explicit later run supplies the stated existing controls.
- `compile.log` and all six module logs show compile-only success with
  `-run '^$'`. They select no test bodies and do not establish broad runtime
  acceptance, race/vet results, a passing grade or browser/provider operation.

**PASS FOR VALIDATION for R4's held source.** No new blocker or material avoidable
complexity found in this bounded delta. Root owns subsequent validation and
artifact builds. The old frozen R3 binaries and original 0/100 and 85/100 false
results remain historical and are not R4 evidence; the first Chapter 1 grade is
still unrun. No grade, probe, rebuild, provider/browser task, student/private
feedback or automatic further correction was initiated here.

### Root R4 validation release

Read the held R4 production delta, field-aware tests and saved sensitivity/focused/compile results; retained reviewer PASS FOR VALIDATION was saved and joined. Accepted the compact Engine-only correction for renewed validation, not chapter acceptance. Released sequential root and GUI race/vet for the changed request path, fresh CLI/GUI builds, then one Chapter 14 grade and conditionally the still-unrun Chapter 1 genuine lifetime gate. Unchanged JS/application-module results carry with their existing limits; no repetitive campaign. Any unexpected failure stops before repair or downstream work. No operational execution is authorized; only a short proposal may follow all passing gates. Older R3 binaries and failed 0/85 grades stay historical.

## Final R4 code/validation verdict — PASS; operations pending

Completed the results-only review against the already reviewed R4 source.
Read the current student validation qualifications, exact command/exit statuses,
full main/GUI race and vet outputs, build receipts, frozen hash/check receipts,
and complete Chapter 14 and FIRST Chapter 1 JSON/stderr in
`/tmp/e3-ch14-student/r4-validation/`. Resumed after the connection interruption
without restarting source review or executing any check.

- Main race/vet: exit 0, 40.90s / 0.62s; GUI race/vet: exit 0, 20.68s / 0.26s.
  Package cache/no-test labels remain as reported. The two one-time builds
  succeeded; recorded hash checks passed before Chapter 1 and after grading.
- Frozen CLI SHA-256:
  `05663974c2785fca2d5a3b0602c22be42d838a244a3d994de2504ba636e17a58`.
  Frozen GUI SHA-256:
  `fddf82b3d230a244d88eb536c81e3f042488b7f1a5c4d29d9819319eaedb8f22`.
  These are the inspected receipt identities; root separately verified the
  actual artifacts. They are not the older R3 binaries.
- Chapter 14: **100/100, passed true, process exit 0**, all eight checks,
  3.17s. Both instruction-prefix and isolated total-reducer checks now pass.
- FIRST Chapter 1: **100/100, passed true, process exit 0**, 124.89s.
  Full JSON records five rounds/requests, preserved conversation, cumulative
  input/output usage **511/65**, and complete answer/usage after separate genuine
  **62-second header and 62-second body delays**. This is actual lifetime
  evidence, not a short observation or an interrupted attempt; it does not by
  itself prove arbitrarily longer lifetimes.

**CODE/VALIDATION PASS.** No remaining code blocker from the retained source
reviews or these results. Unchanged JS and ancillary module results are carried
under root's release, not relabelled fresh executions; compile-only/no-test
modules stay compile-only. The selected `client_round1_test.cjs chapter12` mode
is not a pass of the old default script or its skipped adversarial-DOM/late-delta
block. Fresh GUI tests' ordinary integrations retain their actual scope.

All earlier failed suites, setup failures, default-client failure, original
0/100 and 85/100 grades, and their classifications remain historical evidence.
This verdict is not overall operational/chapter closure: external provider,
credentials/spend, browser/native/device, restart/deployment and any operational
exercise remain subject to root's separate authorization and observation.
Local preparation and the operational proposal were not audited here. No new
execution, probe, build, grade, source edit, child or operational task occurred;
only this existing review note was appended. Root owns acceptance and next steps.

### Root code acceptance; operational preparation only

Read and joined the retained final code/validation PASS after its same-context connection recovery. Independently inspected commands, complete grade JSON, gate output, actual frozen hashes and lifetime results. Accept Chapter 14 code and local validation after four production rounds; this is not overall chapter acceptance. Released one small seven-turn, stdin-driven GUI operational preparation/local fake rehearsal only, with source and binaries held. No provider/key/discovery/browser/native execution is released. Scope is the public ladder exercise, not a provider matrix or new harness framework; root will inspect the exact launcher before any separately bounded real-model release.

### Root operational preparation review

The first local frozen-GUI rehearsal stopped at its 60-second readiness deadline before any stdin turn or POST. No real key/discovery/provider activity occurred. Root read the complete 559-line launcher and receipt: lsof requested named addresses while matching literal numeric loopback; the missing original lsof output prevents claiming that as the proven timeout cause. Released only a versioned numeric-listener predicate repair and one local rehearsal. Also corrected a controller restriction that disallowed a legitimate ordered handoff/unload same-batch companion; final canonical ordering remains required. Root withdrew an initial shared-read_file-grant concern after reading the actual prepared/public base (no grant) and manual; no skill input change is needed. System Python has 191 CA roots with certificate and hostname verification enabled, so no TLS modification was warranted. Code remains accepted and held; real execution remains closed pending the saved local result and root release.

### Root one-attempt real-model release

Read the complete v1 launcher, narrow v2 diff, actual pins and rehearsal receipt. V2 fake rehearsal passed all seven turns with sixteen local POSTs, no real key reads, full results, recorded result/tool cuts, SYSTEM/manual survival, handoff/unload and normal EOF. The initial readiness failure remains failed and causally qualified. Pure checks cover ordered handoff/unload companions and forbidden sequences; root's guessed standalone check filename did not exist (checks were inline), not missing execution.

Released only pinned launcher-v2.py (SHA256 0d442acaeb1a4f78e6e16c061151a3ae34ccf6645366c72f32ae91c525d1dbbc), live-01: one discovery GET confirming Sonnet 4.6, one approved in-memory key read, at most 24 POSTs with unchanged native 1,024 cap, existing 60-second setup/420-active/480-absolute bounds, no retries/fallback/redirect/body changes. Owned stdin GUI session, no browser/native/audio/MIC/live8088. First failure stops; unused slots expire. This is execution authority, not operational acceptance.

## Bounded real Sonnet observation — partial failed task; no demonstrated product defect

Read-only classification of `sonnet46-preparation/live-01` receipt, safe summary,
trace, saved response 14/request 15, snapshot/journal and the approved
`launcher-v2.py` guard/exchange/storage order, against current §14.16. Inspected
existing JSON only; no replay, product execution, request, browser interaction,
key access or new campaign. No raw capture or credential material copied here.
The prior CODE/VALIDATION PASS remains unchanged.

**Keep the live receipt failed:** `passed:false`, stop reason
`unexpected/repeated handoff`. Fifteen POSTs were reserved; fourteen responses
were saved and forwarded. Response 15 reached the structural guard after the
exchange returned, but that guard precedes response storage/forwarding, so its
raw response and usage are unavailable. The HTTP-200 conclusion follows the
exchange path's status check; it is not a separately logged per-exchange status.
The trace/receipt reports one model GET and one key read. All fifteen saved
requests retain the pinned real model and native `max_tokens:1024` body.
Known usage for responses 1–14 is **35,983 input / 1,018 output**, cache writes
and reads zero; response 15 is **unknown**, not zero. These are incomplete totals,
not a final bill. The paid authorization is expired; no continuation is implied.

### What this real session establishes

- Six completed human turns and all six complete file-read observations; the
  retained causal replies carry their inspection codes and statuses. Recorded
  ladder activity includes six result cuts and one paired-tool cut. Real Sonnet's
  two context flags stayed false; this is ladder observation, not evidence of
  synthetic inline declarations or automatic result-stubbing behavior.
- The primary remains in SYSTEM; the manual occurs there exactly once on each
  saved request after load and does not migrate into dialogue. Narration and
  ordinary causal conversation survive the cuts.
- One handoff is recorded at sequence 99, unload at 102, and removal of the
  manual's sole `read_file` grant at 103. Snapshot/context through 105 retains
  the pending-unload manual. Final request 15 retains its SYSTEM instructions,
  lacks the `read_file` declaration, and carries one handoff note as user text.
- That final request has no tool traffic, including the current handoff/unload
  batch's receipts, while retaining ordinary conversation and assistant narration.
  This is the public handoff rule, **not** evidence that the implementation lost
  a required acknowledgement or should restore tool traffic. The saved note's
  next step still says to unload although unload actually completed in the same
  batch; the application preserves the model-authored note rather than rewriting
  it. Neither this mismatch nor the surviving narration establishes why the next
  response tripped the guard.

### What it does not establish

The seventh turn did not complete its requested final answer/confirmation.
Journal 106 records the final request and 107 the interrupted turn; the snapshot
stops at 105. There is **no joint successful seventh turn plus normal EOF save**
in this live attempt, and no live restart/recovery replay was performed here.
A successful fake seven-turn/EOF rehearsal cannot supply those missing real-run
observations. This session also does not independently prove all retention
failure/backup/diagnostic boundaries already covered by code and local controls.
Native/browser behavior was neither performed nor required for this exercise.

**Classification:** unmet live-model task/guard outcome, with substantial scoped
real-model evidence and **no demonstrated product/public-contract blocker**.
The launcher reason identifies the guarded additional handoff, not provider
intent or its cause; response 15 cannot be reconstructed from the saved inputs.
§14.16 calls for observing actual reads, retained narration, cuts and handoff and
explicitly discusses unsuccessful model outcomes. It does not promise model
compliance or require a green controller/final phrase as a new product invariant.
Do not change the failed receipt, invent such a requirement, prescribe a product
patch from this observation, or spend again. Root decides overall acceptance
from the established facts and these explicit limits; no operational acceptance
or additional task was initiated by this review.

### Root Chapter 14 acceptance

Accepted code and the required scoped real-model observation after four production rounds. Joined and read the retained evidence classification; independently inspected the final response/request and saved facts. The live receipt remains FAILED: the second handoff request was blocked, the seventh answer and joint normal-EOF success were not achieved, and response 15/its usage were not retained. This is not a demonstrated product violation of the specified tool-traffic-clearing policy, nor grounds for a new promise or another paid run. Real evidence established six full reads, retained narration and SYSTEM manuals, result/tool cuts, handoff followed by unload, and pending manual survival with capability withdrawal. Known usage covers fourteen responses only (35,983 input / 1,018 output, cache zero); fifteen POSTs plus one GET/key open, with HTTP 200 inferred from the inspected exchange path rather than separately recorded status fields. All live authority is exhausted.

Code/grades, qualified local rehearsal and the real observation satisfy this chapter's requirements; root accepts their scope without relabeling the failed model task. No browser/native/hearing or live restart/recovery is claimed. Proceed to accepted source snapshot and minimal fresh prose pass; one text review and at most one author response remain before closure.

### Author briefing correction

The fresh author caught root's wrong Chapter 1 usage shorthand: the actual R4 grade receipt reports **511 input / 65 output**, not 641/65. Root reread the receipt and confirmed the correction; this changes no grade result or lifetime evidence. Root's guessed source paths for a streaming flag lookup did not exist and establish nothing; the author is checking the real configuration/GUI entry point to remove the obsolete `--no-stream` recipe without adding a feature.

### Root prose and closure decision

Read the saved author HOLD/note and exact 28-addition/8-deletion author delta; one fresh formal text reviewer returned PASS with zero revisions. Saved finals were read and joined; both text contexts are closed. The first 161 chapter lines remain byte-identical to the author baseline, preserving Bill's opening. Baseline/final scoped lint passed with disclosed soft warnings, not universal prose proof. Accepted snapshot copies all 181 source/support files from main with exact bytes/modes, excluding runtime artifacts. Retargeted only the Chapter 14 integration test's snapshot path; `/tmp/ch14-closure-compile.vF5RxR/test.log` is a successful compile-only check with no tests selected, not a replayed grade. Closure will use explicit owned text paths and the normal credential hook; concurrent paper and unrelated work stay excluded. Chapter 15 is next and has not started.
