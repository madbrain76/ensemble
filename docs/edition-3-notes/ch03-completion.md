# Edition 3 Chapter 3 — complete

Chapter 3 completed student-first execution, one code-correction round, independent code acceptance, and one formal text review followed by one author revision. Root verified the two requested text corrections without commissioning another review. Chapter 4 has not started.

## Exact construction checkpoint

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Accepted predecessor: `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`, `edition-3-ch02-r1`.
- Initial Chapter 3 candidate: `0eb6fdffe62e662b0ab96982b01287afba19ff42`.
- Accepted Chapter 3 source: `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, immutable annotated tag **`edition-3-ch03-r1`** peeled to that commit. The student tree was clean at closure.
- Accepted CLI SHA-256: `db9a31df81f88cace5edcdde6af50dfa35ffe25e96fbc7d7585e51dd6e264dde`; accepted source archive SHA-256: `8efe096945bfb35cda09316c30a2956302f0275406a795de4e6d283a366558d3`.
- Main repository preparation parent: `f7530a6643ca00bd2f59a8da60021c567cd7d174`. The copied final implementation remains reference-only; student source does not replace `agent/` or `solutions/edition-3`.

See [machine-readable acceptance](ch03-accepted.json), [initial candidate pin](ch03-candidate.json), [corrected candidate pin](ch03-candidate-r1.json), and the chronological [execution record](ch03-execution.md).

## Entry and input separation

The [entry audit](ch03-entry-audit.md) retained the existing Edition 2 Chapter 3 grader and the Edition 3 Chapter 1 lifetime overlay on the **same explicitly built CLI**, rather than inventing an Edition 3 Chapter 3 mode or a hidden policy-activation flag. Both reports must pass and exit successfully. The older score does not establish selected-model explanation enforcement: public-API tests and independent source/sensitivity review are mandatory. The real negative control passed the old grader even with its explanation after the calls; a no-enabled-tools control failed the tool checks.

The fresh student received current Chapter 3, public instructions, its accepted predecessor and recorded operational guidance, not private grader source, historical solutions, reference implementation or reviewer patches. This is instruction/context separation, not an OS sandbox. Its ownership plan was approved before implementation. Assistance and self-corrections are recorded: six tools; separate stdout/stderr; credential removal while retaining ordinary build environment; no invented fixed shell deadline; and an explicit policy opt-in without provider/model-name heuristics.

The recorded Gemini 3.0-or-newer constraint was enforced with **`gemini-3.8-flash`**. Root discovered that Chapter 2's earlier 2.5 choice conflicted with that constraint and recorded the mistake instead of silently carrying it forward. Minimal access readiness is not student live-use evidence. No deprecation date is asserted here.

## Code acceptance and the useful failure

Both initial grader reports passed, but the independent [initial code review](ch03-code-review.md) required correction of one real defect: ordinary Gemini text carrying a thought signature was incorrectly hidden as opaque thought. This caused blank human answers, false explanation-policy refusal and lost portable text. Provider failure, token exhaustion and a need to downgrade the model were not established causes.

The student corrected parsing/rendering using the existing part metadata rather than adding a new architecture. Ordinary text remains visible and portable; signatures remain bound to their exact part and source; genuine hidden thinking remains hidden. [Round-one code review](ch03-code-review-r1.md) independently accepted the exact revised commit with no remaining blocking finding. It verified the source/archive identities, tested the exact pinned CLI with both graders including real 62-second header/body lifetime gates, ran the full race suite and vet in pristine archive copies, and repeated all eleven exact mutation patches with assertion failure followed by restoration success. There is no numeric coverage prerequisite invented for this chapter.

The [student's initial report](ch03-student-review.md) and [correction report](ch03-student-review-r1.md) preserve failures and recovery: the first local Chapter 3 grader scored 41/100 before self-correction; the wrong fixture pointer made one policy test invalid until corrected; the original signed-text display problem was not hidden behind later success. These were distinguished from the one formal correction round. A report paragraph initially mixed student and reviewer mutation categories; its wording was corrected without changing findings or raw receipts.

## Live and human-interface evidence limits

- Initial-candidate Anthropic and OpenAI real tool/repair/chat evidence remains attributed to that candidate, not relabeled as fresh accepted-binary calls. Initial Gemini 3.8 used all six tools and repaired the project but displayed empty terminal answers; that was an acceptance blocker, not a successful chat exception.
- After the source fix and local reruns, one finite Gemini 3.8 actual `chat` session used the **exact accepted binary**. It completed two turns with nonblank answers, performed the real repair (local test exit 1 to 0), wrote the follow-up marker and preserved conversational context. All six tools succeeded at least once. There were eleven calls/results and one recoverable search error, not eleven uniformly successful results.
- The launcher checked and hashed stdout in process but did not retain raw stdout. Root and reviewer verified the retained history/artifact hashes and ordinary terminal text, and inspected the repaired project. They did not claim to reread absent stdout or hidden provider content. These observations complement, rather than erase, that evidence limit.
- Root also used the actual human `chat` interface with a local provider fixture: an unexplained two-call batch wrote no files; the explained retry wrote a marker; a follow-up retained context. This is synthetic interface/policy proof, not live-provider proof.
- Bounded observation windows lived in launchers, not blanket application generation deadlines. One original Gemini no-result run and the same-model retry remain separate observations. No fallback model or speculative token-cap change was used.

The evidence directories retain source identities, grader JSON/status, failures, tests, exact mutations and restored passes. Executables and disposable live projects remain local artifacts at their recorded paths; they are not embedded in this text-evidence commit. A fresh checkout must rebuild/materialize such artifacts before comparing their recorded hashes.

## Book improvement: one review, one revision

The author improved only Chapter 3: a short signed-ordinary-text failure lesson, building the actual CLI and running both public graders, and consistent live/real-work build commands. Typed-schema field inventory and additional policy machinery were deliberately rejected as unnecessary. No earlier chapter, runtime, grader or reference solution was edited.

The sole [independent text review](ch03-text-review.md) returned **REVISE**, requesting exactly two corrections: express signature binding without universally prescribing the accepted student's `Opaque`/`From` fields, and fix the final root-library build command. The [author's one revision](ch03-author-review.md) made precisely those changes. Root verified that reversing those two edits reproduces the exact reviewed draft, preserved the original author report as an unchanged prefix, and closed both findings. This is not a second reviewer PASS or a second formal pass.

Final chapter SHA-256: `a9cc08a3cc86d34f612226c582bfb130a9589299d721d5c9b154c6f6a32442c1`. Final author report: `6691ceb649dad280f77a2756f2c2e4006ce1017d29ea7bcad0b734203bee189e`. Independent review: `ab761baae75e1ec633c298f6b2503234cc199429d2091e5b3ae80c963c7e50ef`.

Final chapter change is +175 whitespace-delimited words, four hunks, with 12 balanced fenced blocks and the same two resolving local links. [Prose lint](ch03-prose-lint.json) passes all hard checks for baseline, draft and final revision; existing negation/person-gap warning categories are retained as diagnostics, not padding or deletion targets. Root verified all 231 copied reference files and 47 original manuscript files against the baseline, and all 11 pinned grader-source entries. Unchanged prefaces and other chapters remain protected. External URL/anchor validity and rendered publications were not comprehensively checked; no PDF/EPUB build occurred.

## Lifecycle and next boundary

Student final turn 549, code reviewer 451, text reviewer 387 and author final revision 341 were verified durably saved before cleanup. All workers were joined and closed; the live child registry was empty. Code acceptance, text correction closure and repository commit are distinct recorded facts. No push, installation, runtime patch, restart, browser reload or deployment occurred in Chapter 3 execution.

Next: read the current Chapter 4 contract and align its grader plus inherited obligations at entry, then launch a fresh student from `edition-3-ch03-r1`. Do not use the reference solution as starting code, assume a nonexistent Edition 3 chapter mode, or begin another author-led preparation loop.

## Scoped commit verification

The final checkpoint scope is 257 named text files: Chapter 3, the current execution/student guidance, and its acceptance records/evidence. All intended staged bytes and 22 local file links were verified. The credential hook remains enabled. No executables or unrelated files are included.

Git's 209 whitespace warnings are retained evidence formatting, not silently repaired data: 183 diff-context/help indentation warnings, 25 trailing-space/tab warnings (including diff blank lines, build metadata and two four-space lines in the initial patch-failure receipt), and one final blank line in the audit diff. All were classified against actual bytes. Author prose outside its four-backtick diff fences and all other staged files pass whitespace checks. Reviewed hashes and exact mutation patches remain intact.
