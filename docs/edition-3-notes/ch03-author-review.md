# Edition 3 Chapter 3 — post-student author draft

Status: finite draft for the root-owned independent text review. No formal text review has been performed by this author, no reviewer or other child was spawned, and no text-review verdict is recorded here. Root should pin the chapter and report hashes returned with this draft before handing it to its one independent reviewer. At most one author revision remains available if that review requires it.

## Scope and source identity

Main repository: `/Users/bill/projects/ensemble`, HEAD `f7530a6643ca00bd2f59a8da60021c567cd7d174`. The chapter started byte-identical to HEAD, SHA-256 `2366f74bdb873af57fe95e0c7c0b040187365a8db7efc09f5e63ce1f3c59c801`.

Accepted student: `/Users/bill/projects/ensemble-edition-3-student`, clean HEAD `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`; `edition-3-ch03-r1` peels to that same commit. The accepted Chapter 2 predecessor is `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`; the initial Chapter 3 candidate was `0eb6fdffe62e662b0ab96982b01287afba19ff42`. The acceptance receipt records one code-correction round, R1 closed, code PASS, and saved student/reviewer final turns 549/451. Those workers were already closed by the parent. There is no pending code exception.

The exact accepted executable is `/tmp/ensemble-ch03-r1/student-cli`, SHA-256 `db9a31df81f88cace5edcdde6af50dfa35ffe25e96fbc7d7585e51dd6e264dde`; source archive `/tmp/ensemble-ch03-r1/source.tar.gz`, SHA-256 `8efe096945bfb35cda09316c30a2956302f0275406a795de4e6d283a366558d3`. This author checked those hashes and the clean student HEAD/tag read-only, without building or executing the binary.

Read in full: the current author role, AGENTS, mandatory coding skill, workflow, Edition 2 carryover, voice, chapter-writing procedure, shared Edition 3 execution procedure, current Chapter 3, both student reports, entry audit, both code reviews, acceptance receipt and Chapter 3 execution log. Bounded reads recovered truncated output; summaries did not substitute for the required documents. Also read the R1 identity and complete sanitized live-chat receipt. The thin role and shared Edition 3 procedure supersede historical phase paths. The execution guide's older “Chapter 3 has not started” position is not the current acceptance state; the current Chapter 3 records govern that fact. No shared guide was edited.

Source SHA-256s (paths relative to main):

```text
fa363e7ba234019267a6e23fff2d48fb354592db7a56a6f83711bcaa3bd47d9d  .agents/skills/edition-3-author/SKILL.md
5f33f75766707d9437f57632e00e14ba790df738102ff8b556ac62bb3b3f599e  AGENTS.md
bd406de600809afc3cced088830f6c3122653b89e5555aa2551b2e9a29ef603f  .agents/skills/ensemble-coding/SKILL.md
40adbae14d2b85b537a20c707b3cbd243c570c3d9246be7b22584f0624790883  docs/agentic-codebook-workflow.md
d100b86984079428852ea8c2599005401b58522f846c2e97b087319eee88284a  docs/edition-2-carryover.md
35139b7897b13fdf9af038a2bc5bc341d6f1a3638e427578a70d94ced6bd9909  book/voice.md
81381b8892f22b9becfde4887b789f5ad0a1a7ac8f239d503e10b38811715d99  book/chapter-writing-procedure.md
0942416638179f6d251dd3c3f53da79001bfcdfac6de5302b61b2d301b5a3793  docs/edition-3-execution.md
e2501bb153ddb005b2afb361634a7d68868c86a48b2027b492f9a9e627f19487  docs/edition-3-notes/ch03-student-review.md
892bd579a35af4000038a89dc7a6e8c50d3a2b68a8f5066aceee41ec2a6787f7  docs/edition-3-notes/ch03-student-review-r1.md
1bb3b7a012115fe907cecb96c34170225339110bbdb6733f42cf21b4be23ebec  docs/edition-3-notes/ch03-entry-audit.md
a3e3264b8e2862c05b5fbe66289e6f9b80d43435e29886ea9615cf6b491bdf00  docs/edition-3-notes/ch03-code-review.md
31ee5b5515a98f39fdd4433e31f2ecdba828188b89afd7a336fd0e00cfb30fe7  docs/edition-3-notes/ch03-code-review-r1.md
c227109552ce2f90c6232ed44e21024939967fbba0e4472b9891a52ea051a8d9  docs/edition-3-notes/ch03-accepted.json
98174fad51ca4a48af0569bade769d1270cfc11b8737f54d586224bab7e16046  docs/edition-3-notes/ch03-execution.md
3ed047e5b7ce6b11bf1b1d9967961d64eec78c29076ed1669f8d924a6c2f1167  docs/edition-3-notes/ch03-student-r1-evidence/live-gemini-chat.json
```

The initial code-review identity above is its corrected narrative, not the earlier transcription that mixed some reviewer controls into the student's mutation inventory. That report-only correction did not open another student code-correction round.

## What the student demonstrated

These are attributed findings from the student, independent code reviewer and acceptance records, not new code tests by this author.

- The inherited module root is a library; the actual CLI is `./cmd/ch01`. The accepted build command names that package explicitly. Both the Edition 2 Chapter 3 grader and Edition 3 Chapter 1 regression/lifetime grader passed against the same pinned executable with `passed:true`, score 100 and exit 0. The Chapter 3 grader includes Chapter 2 parity. The student's early 41/100 result and implementation self-corrections remain historical failures, not additional formal code-correction rounds.
- Passing those graders did not establish the selected-model explanation gate. The entry audit measured a passing late-narration control and identified the need for executable/source-API policy checks. The student and independent reviewer exercised rejection before any batch effects, preservation and matching results, explained retry within the existing round budget, and an unaffected nonselected path. The accepted student's public selection seam is one implementation choice, not the book's required constructor, flag or configuration shape.
- The initial code review independently reproduced an inherited Gemini parser defect: ordinary text bearing a signature was classified as hidden. That caused blank final API/CLI answers, false explanation refusals and loss of ordinary text across models. R1 retained ordinary text as text, using the existing same-part `Opaque` and `From` for the signature and exact vendor/surface/model binding. Genuine hidden thought remained hidden, and same-source signed replay remained supported. This was a local semantic correction, not an established provider outage or token-cap problem.
- The accepted R1 review records fresh race/vet and both graders, including the real 62-second lifetime phases, plus eleven independently repeated mutation/restoration controls. This author neither repeated them nor converted their recorded results into a new review verdict.
- One authorized Gemini `gemini-3.8-flash` actual `chat` run used the exact corrected binary. The retained receipt records two nonblank terminal answers, eleven paired calls/results, all six tools succeeding at least once, one recoverable `search_files` error, real on-disk repair, independent test exit 1 before and 0 after, and follow-up recall. “All six succeeded at least once” does not mean all eleven results succeeded.
- Raw stdout was not retained. The launcher checked the two complete visible answers against actual stdout in process and retained its hash and sanitized assertions. Root and reviewer separately read ordinary terminal text from retained History and inspected the repair; they did not reread a raw stdout transcript. The retained receipt and code-review qualifications are complementary evidence, not permission to claim stronger capture. Prior paid Anthropic/OpenAI live exercises concern the prior revision; they were not rerun on the accepted R1 binary.
- The required Gemini minimum was 3.0 Flash or newer, and Chapter 3 actually used 3.8 Flash. Chapter 2's accepted 2.5 Flash Lite history remains attributed as recorded despite the missed restriction. It neither authorizes a downgrade here nor justifies rewriting the accepted earlier chapter.

## Exact chapter changes and reasons

There are three hunks only. The existing voice, examples, historical results and rationales remain; there is no broad cleanup.

1. **§3.2, immediately after the presence/order explanation check:** add a short account of the demonstrated signed-ordinary-text error. It makes “ordinary” operational without changing Chapter 2: text remains visible and portable, its signature retains same-part exact-source binding, and genuine hidden thinking stays hidden. It adds no new history representation or policy framework.
2. **§3.7 grading instructions:** name the actual nested CLI build and run the accepted paired graders on one explicit executable. State both `passed: true` and exit 0, and explain the measured policy-coverage limit. This replaces an ambiguous solution-directory invocation; it does not add an Edition 3 `-ch 3` mode, a duplicate Chapter 2 run, a hidden activation flag or a universal policy.
3. **§3.8 hands-on build:** correct the root-library build to the same conditional `cmd/ch01` example, with a cross-reference to §3.7. This removes a second instance of the demonstrated build trap while retaining the runnable example and allowing other CLI layouts.

Exact diff against main HEAD:

````diff
@@ -218,6 +218,13 @@
 after the first call is too late. The runtime checks presence and order, not
 the quality of the explanation.
 
+Do not mistake signed ordinary text for hidden thinking. The Edition 3
+student's inherited Gemini parser did, producing blank answers and false
+explanation refusals. Keep that text visible, with its signature in the same
+part's `Opaque` and its exact source in `From`. The text can cross models; the
+signature cannot cross vendor, API surface or model. Genuine hidden thinking
+remains hidden.
+
 Reject the whole unexplained batch before invoking any handler. Preserve the
 assistant response and return a matching error result for every call, saying
 that nothing ran and asking for explanation before retrying. The retry uses
@@ -675,12 +682,28 @@
 2's work, and folding it would hide the one failure a student is most likely
 to cause and least likely to notice.
 
-From the repository root, run:
+Build the CLI package, not the library at your module root. With the inherited
+`cmd/ch01` layout, run this from the student module (substitute your actual CLI
+package if it lives elsewhere):
+
+```bash
+go build -o /tmp/ch3agent ./cmd/ch01
+```
+
+Then, from the course repository root, run both graders against that same
+executable:
 
 ```bash
-go run ./cmd/grade-edition2 -ch 3 path/to/your/solution
+go run ./cmd/grade-edition2 -ch 3 -json /tmp/ch3agent
+go run ./cmd/grade-edition3 -ch 1 -json /tmp/ch3agent
 ```
 
+Both must report `passed: true` and exit 0. The Chapter 3 grader includes
+Chapter 2's checks; the Edition 3 Chapter 1 run preserves the transport and
+lifetime checks. A score alone does not establish the selected-model
+explanation policy: its rejection-before-effects, explained retry and
+unaffected nonselected path still need public-API tests and source review.
+
 The revised grader keeps all ten checks and their weights. It uses the
 corrected Chapter 2 checks and adds observations for search context and append,
 both already required by the prose. Empty-registry omission and isolation stay
@@ -735,11 +758,11 @@
 vaguely, the abstraction collapses into something physical.
 
 Start the serve-only fake from §3.7 and use its printed environment with your
-own binary. This costs nothing and needs no provider key. From your module
-directory, build once:
+own binary. This costs nothing and needs no provider key. If you kept the
+`cmd/ch01` layout, build the CLI from your module directory as in §3.7:
 
 ```
-go build -o /tmp/ch3agent .
+go build -o /tmp/ch3agent ./cmd/ch01
 ```
 
 Then launch `/tmp/ch3agent chat` from the disposable project, with `CH02_LOG`
````

## Changes deliberately not made

- No new public policy constructor, `Config.RequireToolExplanation` requirement or `LLM_REQUIRE_TOOL_EXPLANATION` flag in the chapter. The student's approved seam is not the only permitted implementation.
- No redundant Gemini typed-schema/extras inventory. The live schema issue is documented in student evidence, but the chapter already requires vendor-specific declarations and fake/live verification. It does not need a new schema abstraction or a catalog of provider fields to teach this case.
- No model-name exception, new defaults, token-cap advice, general empty-answer framework, streaming, asynchronous work or policy engine. The synchronous dispatch boundary and existing round budget stay intact.
- No rewrite of the inherited Chapter 2 facade/History contract or its accepted model history. No change to earlier chapters, protected preface, shared guides, app, grader, student code or reference solutions.
- No replacement of the chapter's historical repair account with the newer R1 run. Its useful narrative and dated results remain; the newer evidence and its qualifications belong in this report.
- No optional prose polish, link churn, broad inventory, word-count target or whole-book loop. The unchanged examples, rationale and repetition between instruction and exercise were not treated as padding.

## Mechanical diagnostics and limits

Before is main HEAD's chapter; after is this frozen draft. Words are whitespace-delimited (`len(text.split())`), not a prose-only estimate. Fence checks match delimiter type/length and closure. Link checks resolve the two existing inline local targets relative to the chapter; no network requests or publication rendering were performed.

| Diagnostic | Before | After |
|---|---:|---:|
| UTF-8 bytes | 47,686 | 48,852 |
| Lines | 895 | 918 |
| Whitespace-delimited words | 8,109 | 8,287 |
| Fence delimiters | 22 | 24 |
| Closed fenced blocks | 11 | 12 |
| Unclosed/mismatched fences detected | 0 | 0 |
| Inline Markdown links | 2 | 2 |
| Existing local link targets | 2 | 2 |
| Missing local link targets | 0 | 0 |
| Reference-link definitions | 0 | 0 |

Net: +178 words, +23 lines, +1,166 bytes. The added fenced block is the explicit CLI build. The unchanged targets are `../../edition-1/book/exhibit-ch03-tools.md` and `../../solutions/evidence/ch03/live.txt`. `git diff --check -- book/edition-3/chapter-03.md` produced no whitespace errors. These are mechanical diagnostics, not the independent text review.

Only `book/edition-3/chapter-03.md` and this report were written. Existing untracked files/evidence were left untouched. No code tests, graders, builds, live calls, installation, commit, push or restart were run by this author. A read-only Git diff opened a pager; only that owned pager was terminated, and the identity/hash command then completed. No material contract question arose that required inventing a design. Both text files are finished before the final hash/submission step; any subsequent change requires a new identity supplied to root.


## Sole formal review response — one authorized author revision

The preceding report is the unchanged historical record of the first frozen draft, not a claim that its draft status, three-hunk diff or after-counts describe the final revision. Its frozen identities were chapter `796ea212fb7acead95083e7942561eadfd3e1f566acdbf6b43e00b74b22177e6` and report `12f31a4bc8c57b8813daf8ce051c03f400ed4d52d0c55bca3b72bc77198375ac`.

Root delivered the sole independent text review, `docs/edition-3-notes/ch03-text-review.md`, SHA-256 `ab761baae75e1ec633c298f6b2503234cc199429d2091e5b3ae80c963c7e50ef`, verdict **REVISE**, with exactly two required corrections. This author read all 65 lines and checked that hash. Root reports verified reviewer final turn 387 and owns bounded finding closure; no second formal review is requested or performed here.

### Response to the two required findings

1. **Signed ordinary text:** removed the prescription of the accepted student's `Opaque`/`From` fields. The paragraph still tells the demonstrated blank-answer/false-refusal story, keeps ordinary text visible and portable, binds each signature to its exact text part and source vendor/API surface/model, limits signature replay to that same source, and keeps genuine hidden thinking hidden. These are behavioral obligations, independent of either permitted Chapter 2 representation. No replacement struct, field inventory, public seam or earlier-chapter edit was introduced.
2. **Final real-work build:** changed only the final command's package operand from `.` to `./cmd/ch01`. The actual-package qualification already supplied in §3.7 applies. The destination, directory setup, shell command, path and following chat invocation are unchanged.

Exactly these two replacements separate the final chapter from the reviewed draft; the overall chapter diff against main HEAD now has four hunks. Exact revision diff:

````diff
--- frozen-author-draft
+++ one-author-revision
@@ -220,10 +220,9 @@
 
 Do not mistake signed ordinary text for hidden thinking. The Edition 3
 student's inherited Gemini parser did, producing blank answers and false
-explanation refusals. Keep that text visible, with its signature in the same
-part's `Opaque` and its exact source in `From`. The text can cross models; the
-signature cannot cross vendor, API surface or model. Genuine hidden thinking
-remains hidden.
+explanation refusals. Keep that text visible and portable. Bind each signature
+to its exact text part and source vendor, API surface and model, and replay it
+only to that same source. Genuine hidden thinking remains hidden.
 
 Reject the whole unexplained batch before invoking any handler. Preserve the
 assistant response and return a matching error result for every call, saying
@@ -882,7 +881,7 @@
 
 ```
 mkdir -p "$HOME/bin"
-go build -o "$HOME/bin/ch3agent" .
+go build -o "$HOME/bin/ch3agent" ./cmd/ch01
 ```
 
 Then go to a project you care about and run it there:
````

### Final identities, diagnostics and limits

Final chapter SHA-256: `a9cc08a3cc86d34f612226c582bfb130a9589299d721d5c9b154c6f6a32442c1`. The final report SHA-256 is returned in the submission after both files are finished; it is not embedded as a self-hash.

| Diagnostic | Main HEAD | Reviewed draft | Final revision |
|---|---:|---:|---:|
| UTF-8 bytes | 47,686 | 48,852 | 48,836 |
| Lines | 895 | 918 | 917 |
| Whitespace-delimited words | 8,109 | 8,287 | 8,284 |
| Fence delimiters | 22 | 24 | 24 |
| Closed fenced blocks | 11 | 12 | 12 |
| Inline links / existing local targets | 2 | 2 | 2 |
| Unclosed/mismatched fences detected | 0 | 0 | 0 |
| Missing local link targets | 0 | 0 | 0 |

Same whitespace-word, delimiter-closure and local-link methods as the initial report. Relative to the reviewed draft: -3 words, -1 lines and -16 bytes. Both original link targets remain unchanged and resolve. The replacement construction and line diff were checked to contain only the two requested changes; no optional prose edits were made. Whitespace checks found no trailing whitespace.

Read-only protection checks: all tracked changes against main HEAD are confined to Chapter 3. Both protected prefaces still match HEAD and each other at SHA-256 `915cf8c531169c5aacc4cb81bd743eeba8a181d20bc779d936ac5f7ab9998d96`; Chapter 2 still matches HEAD at `04b112390fac40d3dc12c2b58533573faffc717b997dffc9ce94494c97a487da`. Accepted student HEAD and peeled tag remain `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, clean. Fifteen of the sixteen source hashes in the initial report remain current. The execution record has since gained root's text-review/revision record and is now SHA-256 `c2c94c8aa72bb32d791ecc4a92ac6cda4c70c9b0d140ded3394a9a9a4b5cedcf`; its earlier identity above remains historical. This author did not edit that log or the independent review.

Only Chapter 3 and this append-only author-report response were written. No children, new review, code/grader/live tests, build, installation, commit, push or restart were used. Accepted code, tag, graders and evidence remain frozen; all original live-evidence qualifications remain in force. This is the one authorized author revision, supplied for root's bounded finding-closure verification and final prose lint, not a self-issued review PASS or a request for another formal pass.
