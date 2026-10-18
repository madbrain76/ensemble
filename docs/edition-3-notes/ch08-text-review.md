# Chapter 8 formal text review

**Verdict: PASS.** This is the one formal text review of the final Chapter 8 edit. **Required revisions: none.** No chapter edit, child/review commission, grader run, application run, browser exercise, provider call, publication operation, or student-code review was performed.

## Reviewed identity and scope

- Current chapter: `book/edition-3/chapter-08.md`, 27,335 bytes, SHA-256 `fb16fce247ad7ae955113ee9e1793681e17beee14b754b0b5505fa4f06c3c662`.
- Author baseline: 26,124 bytes, SHA-256 `8eb2d52071ae59ca095bdad3431ba72903419fbec1eae4bf0779ec99600ba14d`.
- Commission: 6,460 bytes, SHA-256 `404b09fbf595f9129557356f3531dc1ff4f3b848a0acce3fab6651a75452af78`. The assignment and all 27 current input pins are saved in the [review evidence](ch08-text-review-evidence/).
- All required current prerequisites and bounded evidence were read in full before this verdict, including the complete 2,434-line historical workflow and 1,110-line carryover. The [coverage receipt](ch08-text-review-evidence/read-coverage.json) selects complete reads, not truncated attempts. Current assignment, repository instructions, and Edition 3 roles govern; historical phases do not authorize another exercise.
- The author inventory verified: 22 internal files and two external outputs match their sizes/hashes, with no missing or unlisted author-evidence members. The saved final chapter is byte-identical to the current chapter. The independently generated diff exactly matches the author's diff. See [seal verification](ch08-text-review-evidence/author-seal-verification.json) and [chapter checks](ch08-text-review-evidence/chapter-checks.json).

The change is confined to the grading block: baseline line 50 and lines 53–61 become final lines 50–51 and 54–84. The opening, settings/GUI teaching, examples, historical evidence, and closing are unchanged. A bounded correction is appropriate: the implementation review found no further chapter ambiguity requiring an enlarged rewrite.

## Assessment

### Teaching, voice, and self-wielding: pass

The opening at lines 3–32 remains about starting to live inside one's own agent, with the reader responsible for deciding whether the GUI is useful. Its personal TTS motivation is specific, not an invented anecdote. The TL;DR leads to a small settings store and the Agent-owned execution policy, then supplies the corrected grading route.

The chapter teaches the distinctions a reader needs: preferences versus effective Agent state; received, persisted, broadcast, and consumed values; browser layout versus server checks; next-turn policy versus a live setting; native speech callbacks versus listening. The fractional-speed example at lines 490–540 connects the saved number, slider sanitization, native consumer, mutation, and tolerance. This is useful causal explanation rather than a generic demand for more testing. The single-node tree at lines 408–458 explicitly does not invent child-agent ownership ahead of the curriculum.

The added procedural material stays inside the grading block. It does not turn the opener or narrative into an acceptance ledger. The settings tutorial and concrete browser actions retain the chapter's readability and the personal reason for self-wielding.

### Edition 3 recipe and exact media-map claim: pass

Final lines 50–84 now provide the supported `grade-edition3 -ch 8 -json` route, an explicit core executable and all five module directories (`core`, `consumer`, `workflow`, `stream-exercise`, `gui`), and explicit streaming and GUI model selections. Optional flags precede the positional executable. The separate Chapter 1 invocation uses the same frozen core CLI, requires both JSON `passed=true` and exit zero, and preserves the two genuine 62-second delayed-header/body phases. Chapter 7 parity is not misrepresented as replacing that lifetime gate.

This agrees with the public recipe in `docs/edition-3-student.md`, the accepted grader/gate records, and the inspected CLI: `cmd/grade-edition3/main.go` lines 24–34 define the flags; 58–87 validate the route, model selections and module identities; 97–103 load the media map; 121–122 dispatch Chapter 8; 130–141 write JSON and enforce the success exit. `internal/gradeedition2/ch08.go` supports the retained settings/parity teaching and restart distinction. No grader was executed for this review.

The stronger wording at lines 69–71—**exactly** `anthropic`, `openai`, and `gemini`, each with a nonempty model ID—is correct, not merely plausible. The CLI delegates this validation outside the two originally allowed source files. Parent permission `main-12` authorized a narrow additional read, recorded in [the permission receipt](ch08-text-review-evidence/narrow-source-permission.json). `internal/gradeedition2/ch05_media.go` lines 21–41 decodes into a fresh string map when a file is supplied, rejects any key outside those three and any whitespace-only value (33–35), and requires map length three (38–39). Together those conditions enforce the stated key set. The function's default values do not establish live capability. The chapter correctly says the map selects local fixtures rather than discovers live-provider capabilities. Source identity and the exact read range are retained in `reads.jsonl` and the coverage receipt.

### Evidence attribution and acceptance: pass

The complete initial and r1 student reports, initial and follow-up independent code reviews, grader acceptance, final code acceptance, scroll continuation, delegated visual report, and root visual decision were read. They are not interchangeable kinds of proof:

- The initial independent review's blocker was one-tab TTS cancellation after another tab changed the global setting. The r1 review closes that specific defect and related coverage; the earlier failure remains part of the record.
- The operational scroll continuation is separately identified and bounded. It is not silently presented as original student success or as a new text requirement.
- Final implementation acceptance includes root-recovered artifacts after the student's failed final packaging. This review does not turn recovery into a successful student final submission.
- Visual acceptance relies on the named delegated visual inspection and the root's recorded decision, not direct pixel observation by this text reviewer or the root.
- Native speech API evidence is not human listening evidence. Chapter lines 538–540 state the limitation explicitly.

The manuscript does not import current student's implementation details or author-workflow failures into its historical story. Its original fractional-speed check is expressly limited to the earlier 2× range (507–517), and the headed-browser account is labeled historical (557–563). The current acceptance records therefore do not silently upgrade old receipts into evidence for the later extension. Historical solution links were checked for existence only, not opened; no solution-tree source was read. The complete bounded records supplied by the assignment are the basis for judging current acceptance. Browser/provider/runtime results cited here are those records' results, not fresh observations.

### GUI/settings boundaries and completeness: pass

The prose keeps browser preferences in the optional GUI module and restores the Agent-owned `max_tool_rounds` through the public interface (36–43). It preserves the server check table and the distinction between fixture inputs and mandated defaults/UI endpoints (86–105). Runtime enforcement, invalid loads, failed writes, browser controls, and visible behavior are explicitly outside what those settings echoes alone prove. The new recipe retains Chapter 7's launch/dependency/provider-environment boundary and explains why restarting in the same disposable working directory matters (80–84; 544–548).

Settings propagation at 482–488 includes the sender, other connected clients, fresh subscriptions, and restart while retaining the existing slow-client limitation. The speech-rate explanation follows the value to its consumer instead of trusting a restored label. No new model row, module boundary, speech design, runtime policy, or product feature is required by this text review.

## Mechanical checks and judgment

The same current linter was read and run sequentially against paths both named `chapter-08.md`: the author's baseline and the live final chapter. Both exited zero. Both report 11 negation forms and a person-gap warning; the gap changes from 3,044 to 3,156 with the recipe expansion. Other reported checks pass, including zero em-dashes and zero LLM crutches. See [baseline](ch08-text-review-evidence/paired-baseline-lint.log) and [final](ch08-text-review-evidence/paired-final-lint.log). These soft diagnostics do not justify padding the chapter with personal flourishes or removing useful technical contrasts. There is no lint-based required revision.

The chapter is UTF-8, has no NUL or CR bytes, and all five relative link targets exist. Their contents were not read. The corrected `git --no-pager diff --check -- book/edition-3/chapter-08.md` exits zero. All 27 pinned inputs were unchanged on verification. Final report/link/scope checks and identities are in [final checks](ch08-text-review-evidence/final-checks.json). Status comparison is path-level evidence only, not a content-preservation assertion about excluded publication assets.

The author's disclosed workflow-tail read-order deviation and earlier out-of-scope publication hashing are retained in its report/failure log; the present verdict does not certify those actions as compliant. They do not falsify the now-verified chapter diff, source-backed recipe, or acceptance provenance.

This review also had setup failures: oversized combined reads, guessed nonexistent receipt names, an incorrect wait handle and repeated permission question, unavailable `rg`, one whitespace-check invocation missing mandated `--no-pager`, and a pre-seal assertion that initially misinterpreted Git's no-index files-differ exit status. All are disclosed in [review failures](ch08-text-review-evidence/failures.md); required text was recovered and the whitespace command corrected before the verdict. The authorized validator window incidentally included the next function's opening; no broader fixture assessment was made. These errors do not support a claim of flawless process, but leave no unresolved text finding or evidence gap.

## Disposition

**PASS; required revisions: none; no further chapter edit or review is requested.** Only this report and `ch08-text-review-evidence/` were created by this assignment. No children or runtime services were started. The [saved final submission](ch08-text-review-evidence/final-submission.json) records the verdict and report identity. The [inventory](ch08-text-review-evidence/inventory.json) lists every evidence member plus the external report by exact size and SHA-256, excludes itself explicitly, and is generated last and then validated read-only.
