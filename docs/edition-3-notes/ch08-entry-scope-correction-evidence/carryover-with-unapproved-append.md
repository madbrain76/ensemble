# Edition 3 carryover: execute the existing book, then clarify

Prepared by Cadence, 2026-10-16, before the Edition 3 light edit. This is a preparation record, not a new feature backlog or proof that Edition 3 has passed any exercise. Bill's current execution instructions are in [edition-3-execution.md](edition-3-execution.md).

## Baseline and reading

The application and book were copied from the tracked working tree at `4309d568` to `solutions/edition-3/` and `book/edition-3/`; [edition-3-baseline.json](edition-3-baseline.json) records exact source revision, file hashes and byte counts. The copied application is an answer/reference baseline, not a fresh student's starting state. No original source, original manuscript, frozen solution or historical evidence is changed by this pass.

I read the complete source manuscript in order: Chapters 00–26, both front-matter variants, the preface, epilogue and Bill's note. The initial large early batch was truncated; Chapter 02 and then front matter plus Chapters 00, 01, 03 and 04 were reread completely in bounded outputs. All remaining chapters were read in full batches. I also read the core authoring/workflow/carryover/voice documents, the archived supporting guides that remain present, edition planning/publication records, the Knowledge Patch workflow, and the Palimpsest introduction. This reading does not substitute for loading the required documents when a future author or reviewer starts work.

Crossover inputs: [Chapter 25 notes](edition-2-notes/ch25-crossover.md), the manuscript's current Chapter 25, and [artifact identity investigation](edition-2-notes/artifact-identity-investigation.md). Existing implementation was checked for total HTTP timeout removal, visible-explanation gating, and restoration of artifact allocation state. Older observations are dated evidence, not automatically current defects.

## Rules for this pass

- Do not redesign chapters or anticipate imaginary student problems. Incorporate observed crossover lessons and fix concrete contradictions/path errors.
- Preserve Bill's preface and attributed historical accounts, including failed runs and accepted exceptions. Do not promote old evidence into Edition 3 receipts.
- Every sentence must earn its place. Useful mechanisms, reasons, demonstrated clarifications and grounded narrative are load bearing; redundant warnings and words for their own sake are not.
- KISS applies to code and prose. No new abstractions, elaborate task system, invented grader or speculative prerequisite.
- No reviewer in this preparatory pass. During execution, the author acts only after student completion and grading/code acceptance, using the student's review; then one text-review pass and one author revision.

## Carryover by existing chapter boundary

| Location | What the evidence teaches | Bounded preparation action / execution check |
|---|---|---|
| 00, local preparation guidance | The first run failed when editorial imagination expanded requirements before student experience. | Identify the Edition 3 working copy and keep the post-student author loop explicit in procedural guidance. Do not rewrite the historical failure story or add another requirement framework. |
| 01, transport | A whole-response 60-second HTTP deadline killed legitimate requests; polling callbacks were not deadlines. | Clarify caller-owned cancellation and absence of a blanket generation deadline. Retain connect/TLS protections. During execution use the recorded delayed-header/midstream and cancellation scenarios; document missing grader coverage instead of claiming it exists. |
| 06–07, 10, streamed artifacts and save/load | Part identity must address the same artifact through deltas, final output, replay and restart. The diagnostic reproduced a new response merging into an old card although archive messages remained distinct. | Clarify the existing identity obligation and a restart scenario. Do not prescribe a new ID data structure. At this baseline `Load` restores History without restoring the runtime part counter; the copied code is not claimed fixed. Legacy colliding saved IDs require deterministic, non-destructive presentation that preserves archive facts; Chapter 10 owns that acceptance case without prescribing a new identity API. |
| 08, settings and speech | UI speed and native voice rate are not the same scale. GUI range, saved validation and the consuming browser must agree; installation, server restart and browser reload are different actions. | Replace the stale universal 2× ceiling with the established approximate 0.5–5× setting; retain <=2× behavior and observed round-trip evidence. Explain higher-rate mapping as approximate rather than calibrated speech timing. Do not rebrand old 2× test receipts as 5× verification. |
| 11, tools and context examples | Browser snapshots and speech-queue inspection were too noisy as automatic observations. General ephemeral/automatic MCP context remains useful. | Keep `gui_snapshot` and `tts_queue` callable on demand. Use a generic environment observation to teach automatic context; do not delete the mechanism. Agent-facing instructions must explicitly request/confirm observations. |
| 12, virtual-user interaction | The chapter must not depend on the automatic browser observation removed above. | Replace automatic/next-round observation promises with explicit observation and confirmation. Preserve accessible/semantic interaction and the real evidence requirements. |
| 14, checkpoint projection | Old encrypted reasoning survived handoffs and dominated wire input. | Already integrated in current teaching. Retain model/surface replay ownership and completed-checkpoint filtering; don't add a second explanation or claim full archived history is erased. |
| 20, delegation | A submission can precede the child's final saved tail; main is not a submitted worker. | Already integrated. Preserve explicit delegated-only submission and quiescent join/shutdown handling. The Chapter 20 reference to Chapter 21 concerns GUI integration and is correct; leave it and historical labels intact. |
| 22–23, memory and recall | Compaction is a visible restricted child, with durable replacement before retirement and safe parent settlement. Memory preserves first-person reasons, failures, constraints and uncertainty. | Already integrated. Keep the current source-selection, structured submission, lifecycle and recall contracts. This actor's first Conversation→Session live test committed a replacement and cleaned up without resetting; that does not verify every band, fault case or Edition 3 grader. |
| 25, late integration | File configuration, provider policy, deployment distinctions and visible explanation were actual user-driven corrections after the initial construction. | Retain the historical account and real evidence links. Construction ownership is now Chapters 3/6, 16, 20, 21 and the memory/recall schema extensions in 22/23, as detailed below. Chapter 25 has no independent exercise or code checkpoint; Chapter 26 extends accepted Chapter 24 source. |
| All working-copy links/build labels | The copy is one directory deeper and includes old generated publications. | Repair relative links, identify Edition 3 as a working draft, and prevent old PDF/EPUB/assembled Markdown from masquerading as new builds. Preserve source preface and historical evidence. |

## Execution gates and unresolved work

1. The existing chapter graders and frozen evidence remain Edition 2 artifacts. Before execution, select the student worktree/predecessor and grader invocation explicitly. Copying the final answer does not establish construction boundaries. Do not silently run a student against the full answer tree.
2. Carryover changes that affect observable behavior need actual grader/test alignment and fresh execution evidence at their chapter. This editorial pass neither patches the original graders nor certifies sensitivity. Use existing public interfaces and negative controls; do not invent requirements to make a grader more elaborate.
3. Chapter 25's treatment is resolved: historical account only, no independent exercise. The construction owners and Chapter 26 predecessor are recorded below and in the current execution guide.
4. The five thin role files now exist for explicit complete reads, with one shared procedure and a separate public student brief. Runtime registration is not required; preserve the student/reference input boundary.
5. Rebuild Edition 3 publications only with paths isolated from the published book. A source-only light edit does not update the copied assembled Markdown, PDF or EPUB.
6. No push, deployment, installation or live restart is part of this preparation. Bounded preparation authors/reviewers were subsequently authorized; no student execution or normal post-student author loop has begun.

## Result tracking

The baseline manifest remains immutable evidence of the copy. The [preparation check](edition-3-preparation.json) records 278 unchanged original files, 231 matching copied application files, 102 resolving local link paths, and no added trailing whitespace. Twenty-seven manuscript/support files changed, many only for deeper relative paths. Chapter/front-and-back narrative grew by 574 words (0.3557%); the total includes the short Chapter 25 crossover account. Bill's preface and note are byte-identical. Historical publication artifacts remain unreconstructed.

This initial pass completed reading and source edits, not chapter execution. It ran no graders, student, reviewer, publication build, push, deployment or live restart. The subsequent authorized closure below records role packaging, the empty student repository and Chapter 1 grader readiness separately.

## Pre-execution closure — 2026-10-16

Bill subsequently authorized closing the known preparation issues before student execution. This supersedes the earlier unresolved Chapter 25 placement, not the student-first chapter loop. See [readiness](edition-3-readiness.md) for actual preparation results and the [crossover audit](edition-3-notes/crossover-closure-audit.md) for source-backed decisions and withdrawn preliminary suggestions.

The remaining crossover contracts have explicit owners:

| Owner | Contract and acceptance responsibility |
| --- | --- |
| 1 | Keep the small environment-configured loop. Use the Edition 3 grader's inherited conversation checks plus delayed-header/body lifetime gates. A 100 score with `passed:false` is failure. Architecture and arbitrary longer deadline avoidance still need source review; check the chapter's documented error, usage and compatibility behavior rather than assume every obligation is automatically graded. |
| 3 and 6 | Enforce visible explanation for the model capability that requires it, before the first tool call in the same batch. Chapter 3 owns side-effect-free rejection and paired history/results within the existing iteration bound; Chapter 6 owns actual stream-arrival order, not post-hoc reordered content. No fixed heading or minimum length. |
| 7 and 10 | Preserve artifact identity across replay/restart. Chapter 10 acceptance includes the established FIRST→FIRSTSECOND→SECOND collision case and deterministic, non-destructive presentation of legacy collisions. Archive facts must remain intact. This closes the behavioral requirement, not the baseline implementation defect; no new public identity API is prescribed. |
| 16 | Provider/model capabilities, native versus compatible routing, and verified reasoning-summary configuration. Preserve the existing default; explicit `detailed` is supported, not a new mandatory default. No speculative pricing or new core-only configuration loader. |
| 20 | Root/child construction shares the existing path and preserves the unmodified base prompt; each child computes its effective model policy. Delegated submission and final-save lifecycle obligations remain. No new CLI file-mode surface is invented here. |
| 21 | Existing `ApplicationConfig`/`NewConfigured`, strict versioned file mode, credential references, private atomic writes, and separate active-versus-pending provider display. Explicit file mode ignores inherited provider environment and allows only the documented launch overrides; no-config legacy behavior remains. Compose only fields whose features already exist. Chapter 8 browser preferences remain separate. |
| 22 and 23 | Extend the same schema for memory and then recall when their types and policies exist. Chapter 22 is not the final dependency boundary because recall arrives in Chapter 23. Live memory application does not make provider changes live. |
| 25 and 26 | Chapter 25 is historical only, with no independent implementation or grader pass. Chapter 26 extends the exact accepted Chapter 24 student source; do not import an old crossover solution or invent a Chapter 25 construction step. |

At each affected chapter entry, align its grader with these existing public contracts, independently review any grader change and preserve sensitivity evidence before launching the student. Historical scores are not evidence for newly clarified checks. This is a scheduled acceptance gate, not a deferred architectural decision or a claim that later checks have already passed. Any materially new requirement discovered later must be discussed rather than smuggled into this mapping.

The Chapter 1 timeout gate is being closed now because it blocks the first exercise. Later source defects are work for the corresponding honest student implementation, not a reason to repair or copy the final reference before the run. Publications remain separate outputs; no copied PDF/EPUB was rebuilt by this closure.


## Chapter 8 finite entry audit — 2026-10-18

The read-only entry audit is complete, not a student run or readiness approval. See [audit](edition-3-notes/ch08-entry-audit.md), [contract matrix](edition-3-notes/ch08-entry-matrix.md), [bounded next-step request](edition-3-notes/ch08-entry-next-steps.md), and [actual local controls](edition-3-notes/ch08-entry-controls/README.md). The Edition 3 publication command still rejects Chapter 8 before launch. Existing settings controls cover sparse patch/broadcast/subscription/restart, but an initial snapshot missing `max_tool_rounds` passes; a no-persistence mutant correctly fails restart. Saved budget execution/next-turn policy, failed persistence/decode atomicity, and actual browser consumers remain separate acceptance obligations.

Only one focused package-test command, two temporary binary builds, two publication no-launch controls, and two settings-only mutants ran; all owned jobs ended. No inherited parity, source/grader edit, student change/run, child, browser, paid/provider-key call or reference repair occurred. Main/student pins and protected files were verified; unrelated changes were preserved. The accepted student remains `edition-3-ch07-r1`; request separate bounded publication/settings-readiness authorization before any Chapter 8 student launch.
