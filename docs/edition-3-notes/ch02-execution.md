# Edition 3 Chapter 2 execution

## Entry — 2026-10-16

Bill authorized continuing after the completed Chapter 1 / GUI-fix handoff. Chapter 2 implementation has **not started**. This entry check is not an author-first rewrite.

- Repository entry: `5351583208988f8370937cd0494727645efa8b46` (tracked/staged clean at entry; unrelated untracked files preserved).
- Student predecessor: `1aad3975842d502b9144b7b9bddfd30ca7309db7`, annotated tag `edition-3-ch01-r1` peeled to that commit, clean.
- Public contract: `book/edition-3/chapter-02.md`, read completely. Current shared procedure: `docs/edition-3-execution.md`; the Edition 3 sequence supersedes historical author-first phases.
- Parent completed current repository/coding/role instructions, execution/readiness/carryover, workflow, Edition 2 carryover, voice and chapter-writing procedure. Combined reads that omitted content were recovered separately.
- `solutions/edition-3` remains the immutable final reference, never student starting code.

## Gates and ownership

| Gate | Owner / status | Evidence / next step |
| --- | --- | --- |
| Chapter 2 contract and harness entry | Parent + read-only `e3-ch02-entry-audit`; in progress | Audit actual grader, retained checks, positive/negative controls and supported executable invocation before launch |
| Three-provider access readiness | `e3-ch02-provider-readiness`; in progress | Current discovery/docs and one minimal inference per vendor; not student live proof |
| Fresh student | Not launched | Only current chapter, public brief/role, accepted Chapter 1 predecessor and public launch clarifications |
| Code acceptance | Not started | Independent exact-revision review; at most three correction rounds |
| Student-informed author/text review | Not started | After code acceptance; one review and at most one revision |
| Accepted Chapter 2 tag/commit | Not created | Requires all preceding gates |

## Bounded grader decision under verification

The Edition 3 CLI currently supports Chapter 1 only. Do not invoke an unimplemented `-ch 2` or omit any submission target. Chapter 2's historical root-package build example does not fit the accepted public-library plus `cmd/ch01` layout; this is a launch mismatch, not a reason to flatten the student's architecture.

Proposed entry gate: build the actual CLI explicitly and pass the same absolute executable path to both commands, from the repository root:

```sh
go run ./cmd/grade-edition2 -ch 2 -json /absolute/path/to/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/student-cli
```

The first retains the current Chapter 2 100-point checks; the second retains Edition 3 Chapter 1's two required zero-point 62-second non-streaming header/body lifetime checks. Read `passed` and exit status, not score alone. Source review still checks ownership, replay preservation and arbitrary blanket generation deadlines. No grader/manuscript change has been authorized merely to modernize command naming.

The entry auditor initially described the lifetime checks as SSE/error probes, then corrected that to delayed non-streaming header/body checks after reading the source. The incorrect description is not a new contract.

## Access and input separation

Approved credential store: `~/.cr/settings.json`; keys stay in process memory. Parent initially checked nonexistent snake-case field names and recovered the actual fields by inspecting structure only: `directClaudeAPIKey`, `directOpenAIAPIKey`, `directGeminiAPIKey`. All three are nonempty, which alone does not establish inference access. No credential values were printed or copied.

The future student must independently perform actual Chapter 2 CLI use, including all three vendors and offline saved-log rendering. Setup probes and historical controls are not student acceptance. Historical/private reference findings stay outside the fresh student's context; separation is instruction-based, not an OS sandbox.

Background compaction belongs to the runtime. No settings, restart, deployment, live-GUI tests or manual compactor lifecycle changes are part of this chapter entry work.

## Entry evidence corrections before acceptance

The entry auditor submitted before executing the parent's requested full-duration two-command controls. Its first audit correctly labels the existing timeout as source-confirmed and its shorter mechanism tests as such; those are not substitutes for the requested real-duration control. The parent reopened a bounded task for the unchanged historical Chapter 2 negative and timeout-removal-only scratch positive, on explicitly built artifacts. No student launched on the incomplete evidence.

The provider setup task exposed a parent error: the parent instructed native OpenAI Responses and later incorrectly asserted that the current Chapter 2 required it. The checker challenged that assertion. Direct absolute-file inspection, equality to the committed `HEAD:book/edition-3/chapter-02.md` blob, and SHA-256 `1a67ee9e3b6ae5f0f776a30f58908c6d53b53be2d8136ef44661728a677dc2d2` establish the actual contract: OpenAI **Chat Completions**, explicitly `chat_completions` at line 1129 and `choices[0].message` at line 872. No Responses requirement is added. The three original minimal setup calls remain evidence of their actual interfaces; the parent authorized exactly one additional OpenAI Chat Completions probe to verify the correct interface, not a retry/spend loop. Provider receipt correction also must resolve its initially wrong repository revision. These checks are setup, not student live acceptance.

The independent audit demonstrated that a missing-call-ID defect can survive the existing black-box Chapter 2 grader while failing a targeted source regression. The entry decision therefore retains explicit independent source/test duties for deterministic/distinct missing IDs, preservation of issued IDs, exact Ref carryforward and redaction, pairing/signature provenance, construction/isolation and deadline inspection. These are existing chapter/architecture obligations, not new features or private fixture answers for the student. No grader rewrite is required merely because automated checks are not exhaustive.

Runtime observation during entry: the handoff's Conversation→Session worker committed at parent event 6880; no further graduation launch was seen in that metadata check. This records durable outcome only, not a new semantic or outgoing-request audit, and did not change settings or worker lifecycle.

## Entry accepted; student launch authorized

Parent accepted the bounded entry after full receipt reads, unchanged input hashes, verified control evidence hashes, and saved final turns (auditor event 267; provider checker event 198). Both tasks joined. The provider checker is closed; the independent auditor remains submitted for the later exact-revision code-review turn, not running a premature acceptance review.

Actual composition controls: unchanged historical Chapter 2 executable fails Edition 3 Chapter 1 at both lifetime gates (100 points, passed:false, exit 1, 120.060s). A scratch copy changing only the blanket timeout and unused import passes Edition 2 Chapter 2 (100, true, exit 0) and Edition 3 Chapter 1 (100, true, exit 0, 124.057s) on the same unchanged executable hash. This is grader sensitivity evidence, not student acceptance. All three taught provider interfaces have actual minimal inference setup evidence. The extra OpenAI Responses call remains parent-requested setup outside the chapter.

Public launch clarifications: preserve the accepted library/CLI layout; build its actual command package explicitly rather than blindly building the library root; use both explicitly targeted grader commands on the same executable. The current chapter governs interfaces and behavior. The launch may name the credential field references and recently discovered economical model IDs, but supplies no private grader details, source mutations, historical implementation, reference solution or repair answers. Any further help must be recorded.

Fresh student ID: `e3-ch02-student-r1`. Assigned source predecessor: `1aad3975842d502b9144b7b9bddfd30ca7309db7` / `edition-3-ch01-r1`. No Chapter 2 code is accepted and no Chapter 2 tag exists. Code correction rounds used: zero. Author/text-review phase has not started.

## Student question — OpenAI cache-write field

During initial development (not a review correction round), the student reported passing several Chapter 2 checks but failing the OpenAI cache-write accounting probe. It requested public wire-field guidance rather than inspect the grader. Parent confirmed the chapter's Usage section names `cache_write_tokens` and the subset formula (lines 705, 718), but does not spell out the nesting. Current official Chat Completions create reference returned HTTP 200 and documents optional `usage.prompt_tokens_details.cache_write_tokens`: “The unadjusted number of prompt tokens written to cache.” Parent supplied that public source/path, not fixture implementation or historical code, and directed the student to verify/record the assistance. The prompt-caching guide also returned HTTP 200 but its retrieved usage example used Responses `input_tokens_details`; the parent explicitly warned not to copy that different wire path into Chat Completions.

Sources:
- https://developers.openai.com/api/reference/resources/chat/subresources/completions/methods/create
- https://developers.openai.com/api/docs/guides/prompt-caching

Bounded fetched-document evidence: `ch02-entry-evidence/openai-cache-field-doc-check.json` (full fetched HTML searched, retained excerpts only; no credentialed fetch). Parent's initial suspicion that this might be only a compatible-provider extension was superseded by the current official schema, not treated as fact. No chapter/grader changes; possible exact-path teaching clarification belongs to the later student-informed author phase. Student implementation/acceptance still pending.

Parent noticed `gui/` in student status and inspected its small `main.go`: it is the separate-module public-root construction/import stub required by the coding rules, explicitly no GUI/transport/live connection. No demonstrated scope expansion; no new GUI feature was authorized.

## Initial candidate and independent review

Student submitted clean `74eb0265bd8bdbc7857d579704300b317b2365da`, preserving the accepted Chapter 1 predecessor. Saved final turn verified at child event 460; student joined. Parent pinned the candidate in `ch02-candidate.json`, verified source archive and binary hashes, read full report/index, and verified every artifact in the student's evidence manifest. Binary `student-cli` and `gui-stub` are local test artifacts, **excluded from any Git documentation commit**; identity records remain. The initial manifest path guard rejected the intentional `../ch02-student-review.md` entry; rerun allowed only that exact report path, not arbitrary traversal. Basic credential-pattern check found no matches; this is not a complete security proof.

The separate entry-audit context now owns the finite independent code-review turn, reading the code-review role and exact candidate. It receives all assistance/known coverage limits and preserves the distinction between automated score and source acceptance. Candidate remains frozen; zero correction rounds used. The student report records six actual CLI live calls and an evidence-verifier regex failure, with offline-only repair of the verification over unchanged logs. No parent acceptance is implied by this report. Author/text review have not begun.

## Initial review: required revisions; correction round 1

Independent review completed with three reproduced blockers and no architecture/grader redesign request. Parent fully read the report (recovering omitted middle), checked cited source and public contracts, verified saved reviewer `turn_ended` at event 430 and unchanged candidate, then assigned exactly R1–R3 to the student as **correction round 1 of at most 3**:

- R1: Gemini names were collected from the whole dialogue before rendering, so a future ID-less call overwrote an earlier result's matching name. Preserve preceding-call identity, as Chapter 2 requires.
- R2: OpenAI/Gemini parsers returned on empty choices/candidates before retaining measured usage. Preserve measured attempt cost while still erroring and not committing a completed response, as Chapter 1 requires.
- R3: rendering multiple same-provenance OpenAI opaque parts assigned absent fields and cleared already-preserved reasoning data. Preserve both already-supported fields without relaxing provenance or adding Responses.

The reviewer independently reran both graders on the exact candidate, real-duration gates, source race/coverage, separate import smoke, five mutations, and keyless saved-log replay. It audited rather than repeated paid live evidence. Grader passes did not establish the missing invariants. Current governance has no numeric source-coverage threshold; the review explicitly corrects its earlier entry-audit claim.

Parent caught a report-summary error before using its evidence: the actual repeated mutation names are ephemera, cache-write, opaque-isolation, anthropic-order, ref-stub. The report had described a different set. Requested a finite **report-only correction**, preserving findings/verdict/raw outputs; this is not another student correction round or a rerun. Initial submitted review SHA-256 before wording correction: `34c5bdcf79500fb9c44e6008351e5ed96bf7ec7642681cf2c6202fc21ed692bd`.

Student received only scoped failures/public contract bases, not reviewer test code, private reference comparison, or repair answers. Revised evidence goes to new `ch02-student-review-r1.md` and `ch02-student-r1-evidence/`, leaving original source/evidence identifiable. Require narrow real-root regressions/sensitivity, final race/coverage, same-binary paired grader checks including actual 62-second phases, and affected offline replay. No new live text calls required merely to repeat unchanged functionality; assess exact delta and label retained live evidence honestly. No accepted Chapter 2 tag, author pass or Chapter 3 work yet.

## Code accepted; author/text phase begins separately

Round-one candidate `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323` independently **ACCEPTED**, R1–R3 closed, no required findings or exceptions remaining. Parent read the complete report, verified all 133 review/evidence checksums, source/archive/binary pins, clean tree and saved review final turn 512 (student final 636), then created immutable accepted tag **`edition-3-ch02-r1`**. One of three allowed correction rounds used; this is the first accepted chapter revision. Code acceptance receipt: `ch02-accepted.json`. Text phase remains pending at this point.

Reviewer independently rebuilt byte-identical CLI, ran source race/coverage (749/838 statements, 89.4% displayed), both same-binary graders including actual 62-second phases, previous public regression, all three scratch fix-reverts, and affected keyless/local CLI evidence. Parent also directly read the six original live stdout/stderr files after correcting a filename-filter mistake; they support the original candidate's two-turn text conversations and usage. No new-binary live call is claimed. All earlier source/evidence/failures remain attributed and preserved. No additional correction round or architecture/grader change needed.

Next bounded author phase may clarify demonstrated nested CLI invocation/paired grading and exact OpenAI cache-field path, preserving the existing contract. Student's review fixes are testing gaps in already explicit requirements, not authority to add a new test-inventory lecture or future features. One text-review pass and at most one author revision; no Chapter 3 yet.

## Text acceptance and review ownership

The author independently delegated the single text review to its fresh read-only child `ch02-text-review`. Parent had separately prepared `e3-ch02-text-review-r1` for prerequisite reading only; parent did **not** send that worker a formal review task. The author task did not clearly reserve reviewer launch to the parent. Once this was discovered, parent explicitly prohibited a duplicate review and requested the actual nested review provenance/lifecycle. This is a coordination correction, not a second text pass. Exactly one formal independent text review was performed; zero author revisions were required.

Parent read the complete author report/diff and nested child's saved submission. The child returned **PASS**, no blocking findings, on exact chapter SHA-256 `04b112390fac40d3dc12c2b58533573faffc717b997dffc9ce94494c97a487da`. Parent directly verified nested reviewer saved final turn 97 and unused prerequisite-only reviewer final turn 213. The nested review's submitted result is preserved verbatim in `ch02-text-review.md`; original submission envelope SHA-256 `29083b1a27cbe1381e548fd1d7ca83de56a0935f9dbf1d312789865212924011`. Author report appends the received verdict after the review; no post-review chapter change. Parent accepts the bounded text, pending final cleanup/checkpoint verification.

Changes are only the documented cache counter path and actual nested-CLI build/paired-grader instructions, plus consistent live build. Chapter +154 split words, 22→23 balanced fenced blocks, three hunks. No new exercise, architecture, runtime feature or regression inventory. Student source and graders stay frozen. No second text review or optional revision will be opened.

## Final closure

Author lifecycle follow-up independently checked its child's fresh context, full read coverage, frozen inputs, saved submission/final turn, join, and shutdown; no second review or repository edit occurred in that follow-up. Review-time author report SHA-256 was `c938e6cecce2c4599bfe630fbd128ff5625fbb0ad45d76b3198cabde1d9f30da`; the final report only appends the received verdict/status. Parent recovered and read the complete lifecycle receipt, verified author final turn 260 and unused readiness reviewer final turn 213, then joined/closed both. All workers, including the nested reviewer, are closed. No author revision was required.

Parent aligned the shared procedure/current-position and public brief with the accepted chapter's same-executable paired invocation; these are status/reference updates, not new exercise requirements or a second manuscript review. `ch02-completion.md` and final `ch02-accepted.json` identify accepted code, tested/live provenance, reviewed chapter, and next gate. Five local executable artifacts are explicitly excluded from Git in `ch02-artifact-exclusions.json`; retained text receipts preserve their hashes and actual outcomes. Nothing is installed, restarted, pushed, or published by this chapter closure. Chapter 3 has not started. Exact staging scope, integrity and whitespace/security checks precede the chapter commit.

### Final staging qualifications

Parent's first explicit staging-list check stopped before staging because it guessed `ch02-launch-manifest.json`; the actual file is `ch02-entry-manifest.json`. The corrected 500-file list was checked byte-for-byte against staging, with five local executables excluded. All 231 immutable reference files and 47 original manuscript files match the copy baseline; 34 local closure links resolve. No tracked application/grader/original-book changes.

The staged whitespace check reported one extra EOF blank line in the submitted author report plus 13 retained-evidence warnings. Parent removed only that one report blank line (no prose/chapter change): original submitted SHA-256 `63d2de51315425d9506d7e654faca5dae944d22d4a82eec772f7868fbd23d54a`, normalized committed SHA-256 `d4cf331827a744988ab92f711c103d32a11f354751ddffd57825d16460d599c4`. Both identities are recorded in `ch02-accepted.json`; this is not an author revision or another text-review pass. The 13 raw warnings remain deliberately: four trailing tabs emitted by `go version -m`, three literal CLI `> ` prompts, and six whitespace-bearing lines from the fetched public documentation excerpt. Those artifacts stay byte-identical to their evidence manifests. The chapter, current guidance, and closure prose pass whitespace checks separately; the credential hook remains enabled.

Final mandatory prose check: parent found the required lint step not claimed in the author receipt, ran `go run ./cmd/lintprose -v` separately on a baseline copy named `chapter-02.md` and the accepted chapter, and saved `ch02-prose-lint.json`. Both exit 0; all hard ceilings pass. Existing soft word/person-gap warnings persist; prose word count 9793→9932, negation forms 19→21 create an additional soft density warning. Parent accepts the useful boundary explanations rather than reopen optional polish. No manuscript edit or additional formal text-review pass. The one new lint receipt makes the final explicit commit scope **501 files**.
