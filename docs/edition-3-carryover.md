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
| 06–07, 10, streamed artifacts and save/load | Part identity must address the same artifact through deltas, final output, replay and restart. The diagnostic reproduced a new response merging into an old card although archive messages remained distinct. | Clarify the existing identity obligation and a restart scenario. Do not prescribe a new ID data structure. At this baseline `Load` restores History without restoring the runtime part counter; the copied code is not claimed fixed. Legacy colliding saved IDs require an explicit decision, not destructive archive rewriting. |
| 08, settings and speech | UI speed and native voice rate are not the same scale. GUI range, saved validation and the consuming browser must agree; installation, server restart and browser reload are different actions. | Replace the stale universal 2× ceiling with the established approximate 0.5–5× setting; retain <=2× behavior and observed round-trip evidence. Explain higher-rate mapping as approximate rather than calibrated speech timing. Do not rebrand old 2× test receipts as 5× verification. |
| 11, tools and context examples | Browser snapshots and speech-queue inspection were too noisy as automatic observations. General ephemeral/automatic MCP context remains useful. | Keep `gui_snapshot` and `tts_queue` callable on demand. Use a generic environment observation to teach automatic context; do not delete the mechanism. Agent-facing instructions must explicitly request/confirm observations. |
| 12, virtual-user interaction | The chapter must not depend on the automatic browser observation removed above. | Replace automatic/next-round observation promises with explicit observation and confirmation. Preserve accessible/semantic interaction and the real evidence requirements. |
| 14, checkpoint projection | Old encrypted reasoning survived handoffs and dominated wire input. | Already integrated in current teaching. Retain model/surface replay ownership and completed-checkpoint filtering; don't add a second explanation or claim full archived history is erased. |
| 20, delegation | A submission can precede the child's final saved tail; main is not a submitted worker. | Already integrated. Preserve explicit delegated-only submission and quiescent join/shutdown handling. The Chapter 20 reference to Chapter 21 concerns GUI integration and is correct; leave it and historical labels intact. |
| 22–23, memory and recall | Compaction is a visible restricted child, with durable replacement before retirement and safe parent settlement. Memory preserves first-person reasons, failures, constraints and uncertainty. | Already integrated. Keep the current source-selection, structured submission, lifecycle and recall contracts. This actor's first Conversation→Session live test committed a replacement and cleaned up without resetting; that does not verify every band, fault case or Edition 3 grader. |
| 25, late integration | File configuration, provider policy, deployment distinctions and visible explanation were actual user-driven corrections after the initial construction. | Replace the placeholder with a concise account linked to the real crossover notes. Keep complete application configuration at this late boundary rather than make early settings depend on later providers/memory/fleet. Explain visible text before calls, not opaque reasoning, and preserve file/legacy test interfaces. Do not fabricate a Chapter 25 executable contract or passing grader. |
| All working-copy links/build labels | The copy is one directory deeper and includes old generated publications. | Repair relative links, identify Edition 3 as a working draft, and prevent old PDF/EPUB/assembled Markdown from masquerading as new builds. Preserve source preface and historical evidence. |

## Execution gates and unresolved work

1. The existing chapter graders and frozen evidence remain Edition 2 artifacts. Before execution, select the student worktree/predecessor and grader invocation explicitly. Copying the final answer does not establish construction boundaries. Do not silently run a student against the full answer tree.
2. Carryover changes that affect observable behavior need actual grader/test alignment and fresh execution evidence at their chapter. This editorial pass neither patches the original graders nor certifies sensitivity. Use existing public interfaces and negative controls; do not invent requirements to make a grader more elaborate.
3. Chapter 25 is an existing crossover record, not currently a fully specified independent construction exercise. Its full execution/absorption treatment remains an explicit decision before reaching that stage; do not hide that gap behind a new speculative author assignment.
4. The runtime role skills are not yet installed/enabled. Package the agreed procedure without duplicating it, keeping student context separate from historical answers and author-only repair records. Avoid a new orchestration framework.
5. Rebuild Edition 3 publications only with paths isolated from the published book. A source-only light edit does not update the copied assembled Markdown, PDF or EPUB.
6. No push, deployment, installation or live restart is part of this preparation. No reviewer or student has been launched.

## Result tracking

The baseline manifest remains immutable evidence of the copy. The [preparation check](edition-3-preparation.json) records 278 unchanged original files, 231 matching copied application files, 102 resolving local link paths, and no added trailing whitespace. Twenty-seven manuscript/support files changed, many only for deeper relative paths. Chapter/front-and-back narrative grew by 574 words (0.3557%); the total includes the short Chapter 25 crossover account. Bill's preface and note are byte-identical. Historical publication artifacts remain unreconstructed.

This is a completed reading and source-edit pass, not successful chapter execution. No graders, student, reviewer, publication build, push, deployment or live restart ran. Role-skill packaging, a clearly separate fresh student location, and grader alignment remain before the corresponding execution stages.
