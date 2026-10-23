# Human-first companion — restart notes

## Purpose and agreed direction

Create a separate human-first companion to *The Art of Building AI Coding Agents*, using **Edition 3** as the source. The instructional book captures knowledge in a form bots can execute; this companion should explain that knowledge so a human can understand the field, take over, and exercise expert judgment.

## Audience and expertise transfer

The intended reader is extraordinarily smart, already very familiar with coding with AI, and a prospective technical lead who can move this work forward after Bill is no longer working on it. This is not an introductory AI-coding book. Explain project-specific concepts clearly without reteaching familiar basics.

Bill reports developing this expertise through at least a couple thousand hours building these tools with this assistant and other LLMs. The instructional book now encodes that expertise. The companion must accelerate a human's acquisition of the judgment behind it: recognizing failures, understanding design rationale and tradeoffs, and deciding when changed circumstances justify a different approach. Stories carry experience, not just entertainment. Never invent the missing history. Quick to read means efficient transfer of deep understanding, not an arbitrary short length.

## Reader-facing format

The reader-facing book is **stories and knowledge in prose**. No TL;DR sections, exercises, assignments, grading machinery, or bot-oriented operating instructions. This is not a mechanical abridgment: preserve the reasons, causal explanations, tradeoffs, failure modes, and technical distinctions that make the knowledge useful. Translate procedural material into understanding where it carries essential knowledge rather than simply discarding it.

The user explicitly wants a **fun, light, quick read that is still technically very heavy**. Light delivery, not light substance. Favor story-led explanations, concrete language, short readable paragraphs and sections, and natural humor. Do not force jokes, invent anecdotes, pad with folksy filler, or flatten technical nuance. Read and follow `voice.md`; the companion changes the teaching format, not the author's personality.

## Source and workspace boundaries

Repository: `/Users/bill/projects/ensemble`.

- Main source: `book/edition-3/` (chapter map: `chapters.md`). Inspect the live map before deciding the final companion scope; Edition 3 is under active development.
- Companion-specific voice: `book/companion/voice.md`, adapted from Edition 3 version 7 to this audience and prose-only format. This is the shared voice guide for companion authors and reviewers. Original `book/edition-3/voice.md` and `book/voice.md` remain unchanged; inherited factual/privacy guardrails remain applicable.
- Bill authorized the recommended blocker resolution: for this companion only, the full active-context reading set is `book/companion/voice.md`, the corresponding Edition 3 source chapter, the preceding accepted companion chapter (if any), and prior-author carryover notes. This replaces the four legacy full-document reading requirement. Repository preservation, privacy, factual accuracy, and independent-review rules remain in force. Do not load the old implementation/grading workflow merely to satisfy a superseded companion reading gate.
- New manuscript: `book/companion/chapter-NN.md`.
- Working author/reviewer notes: `docs/companion-notes/` (not reader-facing book content).
- This file is coordination material, not part of the published prose.

Other agents are actively changing Edition 3 and implementation code. Do not edit, reset, stage, or commit their work. Companion work is additive and isolated. Do not silently copy source publication-status claims, historical benchmarks, or promises into the companion as current facts.

## Required chapter workflow

The user explicitly requested **a fresh author and a fresh independent reviewer for every chapter**. Do not reuse one long-running writer or reviewer for the whole book.

For each chapter:

1. Give the fresh author the user brief above, current `voice.md`, the corresponding Edition 3 source, the preceding accepted companion chapter (none for Chapter 0), and the previous author's carryover notes.
2. Author writes a complete human-readable chapter and short carryover notes: narrative decisions, key knowledge retained, source references, uncertainties, and context needed by the next author.
3. Give a separate fresh reviewer the same voice/source/preceding-chapter context and author notes, plus the new draft. Review technical fidelity, missing essential knowledge, story accuracy, readability/fun, and voice continuity. Check that instructional apparatus has not leaked back in.
4. Resolve material review findings. Inspect actual changes rather than accepting an agent verdict alone. Use existing targeted prose checks where applicable; check tool usage before running them. No code grading or broad mutating book build just for prose.
5. Collect final saved results and join author/reviewer. Mark the chapter accepted only after review and corrections, then pass its prose and notes to the next chapter's fresh pair.

This author/reviewer sequence does not require student exercise or code-grader roles: those belong to the instructional edition, not this companion.

## Current state

- Companion direction and workflow agreed; this restart note is saved.
- Chapter 0 author was launched as `companion-author-00` in the originating session.
- Author ownership: `book/companion/chapter-00.md` and `docs/companion-notes/chapter-00-author.md` only.
- The author has been sent the additional fun/light/quick-but-technically-heavy requirement.
- `book/companion/voice.md` is written; audience/succession brief is saved here and in that guide.
- The initial author stopped and submitted its saved note without drafting after legacy full reads could not be retained. It has been joined. This is not a chapter acceptance or a claim of satisfying that old gate.
- Bill subsequently authorized the companion-specific reading set above. The blocker is resolved. A fresh Chapter 0 author, `companion-author-00-fresh`, is assigned the same owned files and the smaller required reading set.
- Bill describes himself as a developer, not a writer. Handle writing mechanics, continuity, and review without making him manage routine workflow. Ask him about meaning or genuine consequential ambiguity, not routine editorial choices.
- Chapter 0 draft is saved at `book/companion/chapter-00.md`, titled **Inheriting the Judgment**. Parent read the complete draft. Author carryover is saved at `docs/companion-notes/chapter-00-author.md`.
- Author reports targeted `go run ./cmd/lintprose -v book/companion/chapter-00.md` passed all rules, plus whitespace/newline and relative-source-link checks. These are prose checks, not product tests or editorial acceptance.
- **Chapter 0 is accepted.** Independent fresh reviewer `companion-reviewer-00` saved `docs/companion-notes/chapter-00-review.md` with no material findings or required revisions. Parent read the final chapter and actual review, agreed with the assessment, and confirmed the reviewed file hash matches the saved draft. Parent also checked manuscript/voice/restart whitespace and the source-link target.
- Both fresh author and reviewer submitted final saved results and were joined. No review corrections were needed. Author notes describing review as pending reflect the author-stage snapshot; this acceptance entry is the current state.
- No subsequent chapter has been started or accepted. Next: Chapter 1, with a fresh author and separate fresh reviewer, both reading companion voice, Edition 3 Chapter 1, accepted companion Chapter 0, and Chapter 0 carryover.

## Resume procedure

First inspect the actual files and scoped git status; this note may lag a worker's saved draft. If the originating agent session is available, retrieve the author status and final result instead of launching a duplicate writer. If it is not, inspect saved author notes and draft to establish what survived. Do not assume agent handles work in a new session.

The required-reading blocker is resolved by Bill's authorization above; do not reinstate the legacy reading loop. Chapter 0 is accepted; do not repeat its authorship or review without a concrete new issue. Continue with a fresh Chapter 1 author, followed by a fresh Chapter 1 reviewer, and record the actual accepted state here. Then proceed sequentially with fresh pairs, passing the previous accepted chapter and carryover notes each time. Keep these notes concise and current, separating drafts, reviewed chapters, and accepted chapters. Do not turn this into a separate audit-document project.

## Autonomous completion authorization

Bill explicitly requested autonomous completion of **all chapters, the preface, and the epilogue**, with the preface motivated by the original preface. Current Edition 3 map contains Chapters 0–27; cover all existing chapters, not the map's speculative future-topic list. Do not stop for routine approval between chapters.

Body managing editor `companion-body-editor` owns sequential Chapters 1–27 and their author/reviewer notes, using a fresh author and independent reviewer per chapter. Its saved frontier is `docs/companion-notes/body-progress.md`; inspect that plus actual drafts/reviews on restart. Parent retains preface, epilogue, navigation, and this restart file. Fresh `companion-preface-author` owns preface and its author note; fresh independent review follows. Epilogue comes after body completion and may be authored by parent, then independently reviewed. Preserve original source and concurrent work. Bill subsequently authorized committing only our files and pushing to GitHub after every achieved goal, and explicitly authorized pushing all existing commits as well. Never stage other agents' work or use broad `git add .`; use exact owned paths and scoped commits. No deployment or running-service restart is authorized.

## Preface and GitHub publication

Preface is saved and accepted after fresh independent review (`companion-preface-reviewer`), with no required changes. Parent read the draft and actual review; author and reviewer submitted and were joined. Review is `docs/companion-notes/preface-review.md`.

Bill is remote and asked for the current work to be committed now, then committed and pushed after every achieved goal. Existing ahead commits may also be pushed. Current branch/upstream: `main` / `origin/main`, GitHub repository `waywardgeek/ensemble`. Parent handles the initial companion snapshot (accepted Chapter 0/preface and clearly unaccepted Chapter 1 draft if present). Body editor may subsequently commit/push each accepted body chapter using only its exact owned manuscript/note paths, once parent confirms initial push done. All other agents' uncommitted/staged changes must remain untouched. A commit is not a successful push; record failures honestly and do not force-push or discard concurrent work.
