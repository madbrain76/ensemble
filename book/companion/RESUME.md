# Human-first companion — restart notes

## Task and authority

Bill requests autonomous completion of the entire human-first companion to *The Art of Building AI Coding Agents*, sourced from `book/edition-3`: preface, Chapters 0–27, and epilogue. Stories and technical knowledge in prose only: no TL;DRs, exercises, assignments, or bot operating instructions. Audience: an exceptionally capable, AI-coding-fluent future technical lead inheriting Bill's couple thousand hours of project experience. Light, enjoyable delivery must preserve the reasoning, mechanisms, failures, and tradeoffs. `book/companion/voice.md` governs the writing.

Use a fresh author and separate fresh independent reviewer for every chapter. Both read the complete companion voice guide, corresponding Edition 3 source, previous accepted companion chapter, and previous author carryover. Review the actual prose and source, resolve material findings, collect saved submissions and join workers before acceptance. Bill explicitly replaced the legacy four-large-document reading requirement with this companion-specific set; do not reinstate that blocker. This does not waive accuracy, privacy, or independent review. No student/grader pipeline or product build is required for prose.

Bill approved the preface as excellent: leave it unchanged. He permits the parent to author the epilogue; independently review it. Continue without routine approval loops, preserving unrelated concurrent work. Source Chapter 27 was initially a proposal but has since been updated by its owner to report implementation and local testing. Use the current source at authorship: preserve its bounded, source-attributed observations and distinguish them from deployment, student execution, live-provider tests, or verification by this writing team.

## Current ownership and state

- Preface: accepted after fresh independent review and explicitly approved by Bill; published.
- Chapter 0, including Bill’s requested learning/knowledge-preservation section: accepted after fresh independent review and pushed as `ea5025fb`.
- Body Chapters 1–27: managed sequentially by `companion-body-editor`. **Read `docs/companion-notes/body-progress.md` for the current accepted frontier**, then inspect actual files/reviews. At this update Chapters 1–7 are accepted and pushed (Chapter 7: `25d12781`). Chapter 8 has a saved draft and author note, published as a draft snapshot; fresh independent review is next. The application restart left no active children, so the body editor is being restarted at this frontier. Do not start a duplicate author while the editor is active.
- Parent owns Chapter 0, preface, epilogue, reader navigation, and this file. Editor owns Chapters 1–27 and their author/review notes plus body-progress.
- Epilogue and final reader navigation: not yet written; finish after the body so the closing reflects the actual book.
- Working notes are in `docs/companion-notes/`, outside the reader-facing manuscript. Historical notes describe their own stage, not necessarily the latest status.

## Chapter 0 requested addition

Bill asks why ordinary LLM use does not automatically preserve accumulated project expertise. Human experience consolidates, including during sleep; ordinary inference does not update the model's trained weights. Useful contextual adaptation, summaries, and external memory must be distinguished from additional training processes. Notes can preserve substantial experience, but only if recorded and supplied when needed. Avoid claims that models cannot learn or never match human detail; the sleep comparison is an analogy, not biological equivalence.

Human stewardship must preserve the reasons for decisions for future sessions and successors. An agentic codebook makes knowledge executable; this companion translates it for a lead who must understand tradeoffs and know when to depart from inherited decisions. Bill intends it as an example of that translation, not a verified world-first claim.

Fresh author `companion-author-00-learning` submitted and was joined. Changed only Chapter 0 plus `docs/companion-notes/chapter-00-learning-author.md`. Fresh reviewer `companion-reviewer-00-learning` accepted the revision with no material findings; saved review is `chapter-00-learning-review.md`. Both workers submitted and were joined. Parent inspected the actual diff and saved review and accepted the revision. Targeted author lint, whitespace and source-link checks passed; unchanged results were reused. Body work continues without restarting accepted chapters.

## Publication

Bill is reviewing from vacation and explicitly authorizes exact-path commits and pushes after each completed, reviewed unit. He also authorized pushing existing ahead commits. Repository: `/Users/bill/projects/ensemble`, branch `main`, upstream `origin/main`, GitHub `https://github.com/waywardgeek/ensemble`. Initial companion snapshot `32f50d36` was pushed; body editor publishes its accepted chapters.

Stage only owned exact paths; preserve all unrelated staged and unstaged work. Use `git commit --only -- <exact paths>` when appropriate. Never stage all, reset, force-push, or amend another agent's work. Disable pagers with `GIT_PAGER=cat PAGER=cat`. Coordinate parent/editor Git writes; an edit, commit, and successful push are distinct outcomes. No deployment, installation, or service restart is authorized by the publishing request.

## Resume and finish

Inspect actual manuscript, body-progress, scoped Git status, and active children before resuming. Do not duplicate an active writer or repeat accepted work without a concrete correctness issue. Let the body editor complete all 27 chapters, inspect accepted text and saved reviews in batches, and preserve the full fresh-pair sequence. Finish and independently review the epilogue, add a reader-facing contents page linking the preface, Chapters 0–27, and epilogue, check final structure and links, and publish remaining owned files. Collect final saved results and join all children before reporting completion. Do not end the autonomous run merely because one chapter or a checkpoint is complete.

## Goal reminder

At Bill’s request, the parent loaded the goals skill and enabled an eight-minute idle reminder using its Agent-owned `goal.md`. It prompts a state check, not duplicate work or cancellation. Completion still requires all reviews, final navigation and epilogue, successful publication, and child collection. Disable reminders at completion or a genuine human blocker. On another application restart, check whether the local configuration and children survived before restoring them.

## Latest reader direction

Keep the exact Edition 3 main chapter titles and use them in the contents page. Explain what each thing is and why it works that way, motivated where useful by source cold openings and “In Plain English” sections. No philosophical discussion of how agents ought to operate. Every sentence must carry knowledge or earn its place through cohesion, readability, or pleasure. The companion must be significantly shorter than the original book; compare actual lengths without inventing a fixed quota. Approved preface remains unchanged. Body editor is applying this to new work and assessing earlier chapters for concrete, bounded cuts rather than indiscriminate rewriting.

Current frontier at this update: Chapters 0–12 accepted and published; body-progress remains authoritative for later progress. Titles 01–11 were corrected in `7fd23663`; Chapter 12 published in `37736bab`. Parent title00/navigation/voice correction follows.
