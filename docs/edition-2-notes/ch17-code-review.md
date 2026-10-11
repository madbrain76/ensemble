# Original Chapter17 / Edition2 Chapter16 — independent review

## October11 public-contract preparation — no acceptance

Reviewer `enhancement_reviewer_resumed` read all 552 lines of original
`book/chapter-17.md`, file revision `c4770e9ba9475fcb6622df7349d0214c88ba3af9`,
SHA-256 `9b77e32240d10743da3abf0dfd22e585cebcc382990f0517f68d9d006a02f893`,
and current [student-safe preparation](ch17.md). Mandatory skill, carryover
and workflow have been read; subsequent Bill rulings are recorded at `7646941`.
No original answer, grader implementation or changing student source was read
for this preparation. No code/test/live execution or formal review round.

Bill's two decisions resolve the consequential source contradictions: final
mechanical quotas apply when the judge is absent or fails, without overriding
successful judge choices; search all memory files while skipping content already
in actual context. The historical newest-two/system-prompt premise is superseded.
The future review must distinguish a valid empty judge selection from failure,
and recorded recall bytes from a new search during replay. These are existing
chapter obligations, not additional features. Actual public turn entrypoints,
checkpoint lifetime, changed-file indexing and source attribution remain binding.

The chapter describes algorithms and defaults but supplies no concrete executable
grading protocol. Coordinator should recover only missing invocation/fixture
facts if needed: launch target/flags and input envelope; exact configuration
keys for source paths and judge settings; corpus record envelope, especially
the `long` field in learnings; any public inspection/output envelope required
by the named checks. Array judge replies, defensive parsing, score thresholds,
chunk boundaries and byte cap are already public and need no answer disclosure.
Fixed harness paths or source-shape checks cannot silently supersede Agent-owned
files and the optional GUI boundary. No requested contract recovery authorizes
private CodeRhapsody corpus access, a new diagnostic framework or model features.

The student is implementing from accepted original16 tree
`4a253c57dd62b4a6ad086618574ef46af7c2d724`. Formal comparison waits for frozen
submission, completed original grading, meaningful mutations and bounded live
evidence. Original22/24 stops and the stopped author phase remain. Required
handoff: `.agents/skills/ensemble-coding/SKILL.md`.

## October11 grader preflight — first green source, not acceptance

Read the current student checkpoint note and narrowly inspected frozen source
`fcd091f619f879341d4bee318e5f1bf8fb4bb21b`, main tree
`a2fd5a4f71d8fde7f58cfdefd18a0cfc31206d0f`, against the unchanged original
Chapter17 launch/check code. The coordinator's full original run is still in
progress at this writing; the observed log repeatedly reports
`skill "base" has invalid type ""`. No final score or successful runtime
behavior is inferred here. No grader was modified or repeated by the reviewer.

**The missing skill type is a justified fixture exception.**
`ch17_run.go:325–334` writes a base skill with name and description but no type.
The launch selects it as primary (`ch17_harness.go:71–72`). The inherited student
parser requires one of primary/loadable/dependency (`internal/skills/parse.go:96–99`);
original Chapter10 lists these choices and supplies typed examples without
defining an absent-type default. GUI construction consequently fails at
`root.NewAgent`, before HTTP startup or any judge/main request. This repeats the
independently accepted original13 fixture mismatch; no permissive parser change
or dummy primary is warranted merely to satisfy this fixture.

**Two separate assumptions would remain after fixing that launch fixture.**
The GUI selects Agent DataDir `cr` (`gui/cmd/main.go:77`), while grader fixtures
are written below cwd/docs, cwd/handoffs, cwd/learnings and cwd/memory/memory
(`ch17_run.go:82–206`). Recall's source scan correctly follows its Agent DataDir
and configured skills directory (`internal/recall/recall.go:79–96`). No launch
override selects the fixture's cwd as Agent identity. Skills remain separately
reachable through EN_SKILLS_DIR; that does not make the other corpus roots
reachable. This is a latent owner-path mismatch, not an observed BM25 defect.

The structural check (`ch17_run.go:831–915`) examines `internal/recall` beneath
the submitted GUI module. The required separate optional GUI consumes the core
module, where the actual recall spoke lives. Absence below GUI is not absence
of the implementation. Its import detector also exempts paths not containing
`/ensemble/`, so it cannot establish the student's `ensemble/internal/...`
star graph merely by returning success on a different target. Actual core
imports, common behavior and parent relationships still require review.

Nine behavioral checks totaling90 points are masked by launch failure:

| Check | Valid obligation still requiring final supplemental evidence |
| --- | --- |
| bm25-indexes (10) | Real owner corpora, including docs/handoffs/learnings/skills and all eligible memories, are searchable and refresh after edits. |
| bm25-scores (10) | Smoothed IDF, TF saturation and length normalization produce meaningful ranking on a representative corpus. |
| bm25-stop-words (5) | Query stop words are excluded while document length statistics retain them. |
| chunking-splits (10) | Relevant sections, oversized subsections and paragraph chunks retain attribution without unrelated whole-document injection. |
| per-source-quota (10) | Final mechanical selection reserves source shares and redistributes unused slots, while successful judge choices stand. |
| judge-filters (15) | An actual stateless model call changes selection; valid empty choices remain empty, and defensive reply parsing handles documented forms. |
| judge-fallback (10) | Malformed/error/timed-out judge calls degrade to mechanical selection; empty corpus and short/disabled paths remain usable. |
| recall-is-own-kind (10) | Actual entry/event carries immutable selected bytes separately from human input, persists until checkpoint, deduplicates, and replays without retrieval. |
| injection-capped (10) | The whole formatted attachment meets the chapter's6144-byte default cap, including attribution, rather than only the grader's looser8192 bound. |

The remaining10-point structural obligation is valid but its fixed submitted-path
expectation is not. None of these exceptions establishes the masked behavior
passed. Final review must inspect actual runtime tests, mutations, captured
requests/replay and live use; no weakened grader or fabricated total is needed.

One additional latent own-kind predicate deserves final assessment:
`ch17CheckOwnKind` rejects recalled text sharing an Anthropic wire message with
the prompt and categorically rejects role:user (`ch17_run.go:425–451`).
Chapter17§17.5 instead separates stored conversation identity from renderer
placement and explicitly permits an extra part on an existing turn. The existing
Anthropic renderer coalesces adjacent user data. A separate recorded entry is
therefore not equivalent to a separate vendor message or a universal midstream
system role. The current failed launch has not reached this predicate; preserve
that distinction and independently inspect valid wire placement at final review.
Do not add an unconditionally unsupported system message for this source check.

Source provenance: harness/run revisions and hashes match the student-safe
erratum. `ch17_checks.go` revision
`a91362b61fcdd3316bd15069f28add0151892f37`, SHA-256
`476677c924fd26c24923e82f72c3484c4836659b30e3cde6c82468c708a690a0`.
Unlike earlier chapters, the Chapter17 command directly evaluates its ten checks;
this runner adds no recursive Chapter16 parity stage. Formal comparative review
and final student evidence remain pending; no correction round begins here.

Completion update: read the coordinator's committed `c7a84f9`
[provenance](../../solutions/edition-2/evidence/ch17/original-grade-provenance.json)
and [full result](../../solutions/edition-2/evidence/ch17/original-grade.txt).
Actual original score is **0/100 FAIL** on the source/tree above, grader exit1 /
make exit2, completed12:57:45 UTC. All nine runtime checks report startup
failure (fallback stops on its first case); the structural check reports absent
`gui/internal/recall`. Output SHA-256:
`23da91957bf4259cc5376b69932e962d8418a1060397f5424d332c2395dab4f5`.
Thus the missing-type and wrong-module diagnoses are actual failed outcomes;
corpus-path and vendor-message predicates remain unreached limitations. The
0 score is not a measurement of failed BM25/judge behavior. Final acceptance
remains pending complete student submission and independent review.
