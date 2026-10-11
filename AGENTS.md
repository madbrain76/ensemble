# Ensemble student restart

## October 11: stopped before original Chapter 22

Original21 / Edition2 ch20 is independently accepted with recorded exceptions.
The exact accepted export is `solutions/edition-2/ch20/`; main tree is
`f351d66a7bc456650e05ac6656a7db59bbbd725a`. Coder and reviewer are stopped.
Wait for Bill to decide the placement of original22 material before further
implementation or authorship. Do not merge, skip, renumber or advance to23.
The original24 and author-phase stops remain. Use GPT-6 Astra/high for subsequent
authorized coding and code review. Preserve CodeRhapsody’s concurrent work.

## October 11: GPT-6 Astra for coding

Bill approved **GPT-6 Astra (`gpt-6-astra`) at high reasoning** and requested it
for all coding tasks. Use that model for new coding and code-review delegates
and for Ensemble's live coding work. This replaces the earlier GPT-5 Astra
destination and resolves the original20 model-choice stop. Resume the crossover
verification and then the authorized run through original21; the original22/24
stops remain. Preserve historical results and frozen snapshots. A model setting
alone does not demonstrate the crossover; retain actual coding/task evidence.

## October 11: preserve CodeRhapsody’s concurrent Edition 1 work

Bill explicitly reiterates that CodeRhapsody is writing new Edition 1 Chapters
25 and 26. Preserve its commits and working changes. Do not undo those commits,
reset or rewrite shared history, or restore/stage unrelated files. Edition 2
agents commit only their exact owned source, notes and evidence paths, using an
isolated index while other writers are active. Reconcile only those owned index
entries afterward; leave all other work untouched.
Keep the isolated index’s base and the commit parent identical. Publish only if
HEAD still equals that base, using an atomic expected-old-reference check; if a
concurrent writer advanced it, rebuild the owned-path commit on the new HEAD.
Do not attach an older whole-tree index to a newer parent: that can undo unrelated
changes even without a reset. Preserve configured commit hooks where present.

## October 11: Auto-Recall decisions approved

Bill agrees with both original Chapter 17 / Edition 2 Chapter 16 decisions:
apply source quotas to final mechanical selection when the judge is absent or
fails, preserving a successful judge's choices (§17.3); search all memory files
and skip content already in actual context, without excluding the newest two.
The Chapter 17 decision stop is lifted. Continue the authorized student/reviewer
run, preserving original sources and the original Chapter 22/24 stops.

## October 11: prior-request input correction

Bill resolves original Chapter 16 / Edition 2 Chapter 15 forcing: use total input
tokens from the prior request, including its cache buckets exactly once, not
lifetime cumulative usage. The book wording is erroneous. Coordinator/reviewer
may inspect the original solution for this narrow check; provide the student
only verified behavior and provider field contracts, not answer code. Get
Anthropic and OpenAI working first, then Gemini. Preserve original sources and
record the correction in the chapter notes. The Chapter 16 pause is lifted;
all later review and original Chapter 22/24 decision gates remain.

## October 10: corrected solution layout

Bill requires accepted solutions directly in `solutions/edition-2/chNN/`, using
Edition 2 numbering, with no symlinks. Original Chapter 5 is absorbed into
Chapter 4; original Chapters 10–13 are Edition 2 `ch09`–`ch12`. Only the retired
original `history/ch05/` remains. The coder extends `solutions/edition-2/main/`;
original Chapters 14–21 are accepted at top-level Edition 2 `ch13/`–`ch20/`.
Stop before original22 for Bill’s placement decision. The bounded
original20 coding crossover succeeded and was independently accepted.
Never put new accepted solutions in history. Original exercise,
grader, note, evidence and existing tag numbers remain unchanged; state both
numbers in handoffs. This explicitly supersedes the coordinator's history-only
export policy. Preserve snapshot contents and existing Git tags when moving.

## October 10: student run resumed

After completion of the early enhancement and solution renumbering, Bill asked
the coordinator to release the coder and reviewer again. Resume the original
student exercises at Chapter 10 in `solutions/edition-2/main/`, proceeding through
Chapter 21 with independent review and checkpoints at each chapter. The author
remains stopped. Use fresh student contexts without old-answer exposure; the
enhancement coder's reference-answer access does not carry into this run.
The stops before original Chapters 22 and 24 below remain binding. Work promptly
without expanding exercises or repeating unchanged evidence unnecessarily;
required grading, live use, comparative review and correction limits still apply.
Original Chapter 25 is now available and has received full carryover reading;
see `docs/edition-2-notes/ch25.md`. Earlier forthcoming references are historical.
Its capabilities remain beyond the decision gates, not current implementation.

## October 10: later refactor decision gates

Bill requires stops at original first-edition Chapters 22 and 24. Complete the
preceding chapter and its review, then stop before implementing or drafting
either refactor chapter. At 22, Bill will decide whether and how to merge its
material into original Chapter 21 or 23. At 24, stop again for his decision on
placement. The coordinator pauses active delegates and waits; no autonomous
merge, skip, renumbering or advance past these gates. Carry forward the lessons
already approved for earlier exercises, without adding later capabilities early.
These future gates do not themselves release the completed early enhancement
pass into further chapter work.

## October 10: bounded author–coder enhancement pass

Bill explicitly authorized author and coder to enhance original Chapters 1–9
as second-edition Chapters 1–8, dropping original Chapter 5 as a standalone
chapter, with independent review of both prose and code for each chapter.
This supersedes the author-phase timing prohibition below for this bounded
pass only. Include the agreed Chapter 0 execution guidance. Enhance the original
text; do not rewrite it or revive the retired author-first attempt. Preserve
first-edition text, graders, frozen student exports and actual results.
Follow `docs/edition-2-enhancement.md` for mapping, assignments and gates.
The student remains accepted through original Chapter 9; do not advance Chapter
10 during this pass. No future features or expanded exercises are authorized.

Bill's explicit crossover goal: the coder's new Ensemble must run GPT-6 Astra
at high reasoning, and Edition 3 of the book is to be built using that Ensemble
agent rather than Codex. Keep this destination in view without implementing
later features early. See the workflow's crossover goal and verification notes.

Bill cannot use Codex for continuous real-time collaboration. Do not rely on him
monitoring progress or interrupting a mistaken direction. Work autonomously
within agreed scope. For a critical question, stop the work, pause active
delegates, ask plainly with the relevant context, and wait for Bill's answer.
Do not ask asynchronously and continue working, or treat silence as agreement.
Close real-time collaboration in Ensemble is part of the crossover goal.

Bill authorizes the coordinator to resolve missing exercise information by
inspecting the grader when public instructions are insufficient. Supply only
the missing input/output contract, not solution code or inferred extra features.
Record the source revision and exactly what was supplied to the student in
`docs/edition-2-notes/chNN.md`. The student still does not inspect grader internals
or old answers. Actual scope, behavior or architecture decisions remain stop points.

Codex is the primary creator of Edition 2. Bill continues writing and publishing
Edition 1 with CodeRhapsody; its sub-agent chapter is forthcoming. Follow newly
available chapters as part of that evolving source, recording which revision
was used. Do not invent the missing chapter's requirements or claim its exercise
complete before it exists and has been implemented and reviewed.

Human enjoyment and teaching quality are explicit goals alongside correct code.
In the later author phase, both author and prose reviewer read `book/voice.md`
and `book/chapter-writing-procedure.md`, retained for editing after the coder is done:
preserve documented stories and personality, explain the reasons behind the
mechanisms, and judge whether the chapter is engaging and understandable.
An executable specification or passing prose lint alone is not a finished textbook.

Bill retired the author-driven second-edition attempt. Its history is preserved
at `edition-2-attempt-1-stopped`; do not resume its agents or workflow. See
`docs/edition-2-attempt-1.md` for the decision.

Keep the learning record in `docs/agentic-codebook-workflow.md` current. Record
substantive discoveries, failed approaches, corrections and unresolved questions
with evidence. Any agent—coder, reviewer or orchestrator—may add chapter notes
for the future author in `docs/edition-2-notes/chNN.md`, creating a file when
there is a finding to record. Attribute contributions and link evidence or the
student review; preserve disagreements and corrections. Notes do not change the
current exercise or authorize extra features. See the workflow's chapter-note
guidance for details.
Required behavior and architecture are binding, but suggested designs are
provisional. The coder may choose a simpler working design and explain why;
consequential requirement changes must go back to Bill.
The student need not pass every grader test if the grader is flawed: it may
explain the failed check and supporting evidence in its chapter review instead.
Bill authorizes a specific exception: when the coder documents a grader defect
or why an exercise requirement is unreasonable and the reviewer independently
agrees, the reviewer may pass the chapter with a recorded exception and advance the coder without
asking again. Record the actual grader/live results and what was waived for the
future author; do not describe a failing grader as passing. This does not waive
Bill's explicit architecture rules or authorize extra features.

Before coding or reviewing code, and after compaction, read the entire skill:
`.agents/skills/ensemble-coding/SKILL.md`. Every coder/reviewer handoff must include
that exact path. It replaces the deleted second-edition workflow skill.

No mocks without Bill's explicit approval. Use grader-provided fakes where
suitable or small additional fakes. If a coder believes a mock is necessary,
stop and explain why fakes are insufficient, then wait for Bill's decision.
The reviewer's exercise/grader exception authority does not waive this rule.

Before the student run, a reviewer reads the entire available current edition
and completes `docs/edition-2-carryover.md`; the coordinator checks the sourced
lessons and recorded reading coverage. This is preparation, not chapter
authorship or new feature design.
Bill reviewed the carryover and `docs/agentic-codebook-workflow.md` with Codex
and explicitly authorized the fresh student run on October 9: "Please proceed
autonomously." The preparation gate is satisfied. The workflow is the edition-level procedure; the old
`book/chapter-writing-procedure.md` does not reinstate the retired workflow.
The student reads that carryover alongside the mandatory skill and follows
first-edition `book/chapter-NN.md`, in order from the unchanged Chapter 1.
Bill's architecture rules apply from the beginning. Bill has deleted the
second-edition Chapter 1 draft too; use the original chapter and carryover.
There is no author agent. Use the original grader;
record suspected grader defects and teaching gaps for future work rather than
changing chapters or graders during the student attempt.

Only after the complete new Ensemble implementation succeeds and its comparative
code reviews establish the improvements does the author phase begin. The author
and coder then work together chapter by chapter from the original book, all
agents' notes, student reviews and the reviewed implementation. Bill authorizes
improvements to both chapters and graders in this phase: the author revises the
book, the coder improves and validates graders, and each feeds findings back to
the other. Review prose and grader changes, preserve valid coverage, and keep
original graders and student results as history alongside the first edition. Do not
draft ahead of the coder or turn editorial work into new feature requirements.

The restart began from empty in `solutions/edition-2/main/`. Each subsequent
chapter extends its accepted predecessor in that authoritative working tree;
do not restart it or edit frozen exports. The student
must not read or copy `agent/`, first-edition answers, discarded second-edition
answers, or their Git blobs to solve the exercise. Frozen exports are not working
trees. Do not edit the first edition. Never reset history or move existing tags.

After the student's implementation, original grading and live exercise, a code
reviewer reads the new code, matching first-edition solution and the coder's
review for the future author. Both coder and reviewer use KISS as the yardstick.
The skill defines the required comparison of active code, comments and tests.
Reject feature creep, unjustified bloat and complexity added for elegance;
passing the grader is not sufficient. Return findings
and rationale to the coder, not answer code. Do not add features or substitute
reviewer preferences for the original exercise. The coder may revise both code
and its chapter review, rerun affected checks and return to the reviewer.
Default to three revision rounds after the initial review, unless Bill sets
another limit. If still unaccepted, stop progression, record that the student
could not pass within the limit, and bring Bill the evidence to resolve the
chapter, grader or implementation problem. Do not silently start another loop.

Feature creep and code bloat violate the task. Ask Bill about consequential
ambiguity; coordinator preferences are not user requirements. In particular,
pause is Agent-wide: an explicit unpause in any tab unpauses it for everyone.
