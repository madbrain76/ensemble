# Edition 3 Chapter 4 — independent text review

## Verdict: PASS

No required author revision. This is the single formal independent text-review pass on the frozen inputs identified below, not code re-acceptance or final chapter acceptance. No optional polishing round is proposed. The parent retains prose lint, final acceptance and lifecycle closure.

## Frozen identities and scope

Hashes were verified before reading/review and again after the complete read, before writing this report:

| Input | SHA-256 | Complete coverage |
| --- | --- | --- |
| `book/edition-3/chapter-04.md` | `ae7f6fbe1048400f857e8fb3eeefbe513d0a9bab492e1eff1d13a10ff266f37f` | Lines 1–1029, full frozen chapter |
| `docs/edition-3-notes/ch04-author-review.md` | `f30a71586a009123578e81aa9c1e3e3cd1c76545db58ab9f499518348355e27e` | Lines 1–94, full frozen report |
| Baseline `book/edition-3/chapter-04.md` from commit `bac07e644facb25ba3e5bffb2961674763fed4c4` | `016c4e88af4bcb657aad6733e5f65aa09509b0397e84a37059ca245c588eaa55` | Lines 1–1000, full immutable chapter read during readiness; identity rechecked this turn |

The exact `git --no-pager diff bac07e644facb25ba3e5bffb2961674763fed4c4 -- book/edition-3/chapter-04.md` was read in full. Its ordinary-context output has three hunks; `--numstat` reports 40 insertions and 11 deletions. A zero-context diff was also inspected for precise changed-line locations. It shows only the chained-settings paragraph, the grading/build route, and the matching live-exercise build instruction. The review did not expand this into another chapter-design or implementation-review campaign.

Accepted code context remains `ec4a350f564e820d398f9125fd73a9b3b27129cd`, tag `edition-3-ch04-r1`, with the one code correction C1 closed. That acceptance is evidence for this text review, not a claim that this reviewer reran the code gates.

## Findings and rationale

### 1. The grading route now describes the actual supported entry point

Frozen chapter lines 794–825 distinguish the built CLI executable, the core source module and the separate consumer module. The reader builds the actual CLI package from the core module, then runs the grader from the book repository root. The placeholder `./cmd/your-cli` is explicitly replaceable with the real package, including `.` for a root command. It does not require a nested-command layout or conversion of a public library into `main`.

I read all 91 lines of current `cmd/grade-edition3/main.go` and compared the instructions with it and with the complete current public student brief. The adapter parses `-ch`, `-json`, `-core`, `-consumer` and `-consumer-case`, requires one explicit positional target, supports Chapters 1 and 4, passes the built executable and separate source identities to the historical Chapter 4 checks, and dispatches Chapter 1 to its lifetime-aware gate. Its report/exit handling supports the chapter's instruction to require both `passed: true` and successful exit status. The example places every flag before the positional executable, as required by its Go flag parser.

The paired commands use the same `/tmp/ch4agent` path without an intervening rebuild. The text correctly distinguishes the retained 160-point Chapter 4 report from the inherited Chapter 1 lifetime obligation. It does not claim that the 160-point score alone proves all inherited or source-ownership obligations. The surrounding source-review limitations remain intact.

The default `shout` case remains usable without a case flag; a different case is supplied explicitly before the executable. The consumer still builds separately and does not replace the core CLI or acquire its JSON-lines protocol. The actual-module warning avoids substituting a similarly named source copy for the library resolved by the consumer.

These are justified corrections to the documented binary/source entry mismatch, consistent with the accepted adapter and already-public brief. They are not a new grader architecture or a new student feature.

### 2. The C1 paragraph clarifies the existing contract without prescribing representation

Frozen chapter lines 408–414, in §4.6, follow the existing explicit rule that the very next call consumes pending settings, including another `tool_limits`. The paragraph separates the old setting's consumption from installation of the next setting. Its byte-cap/delay example makes the demonstrated omitted-field leakage understandable: the new pending setting uses defaults plus its own explicit fields, not values already spent by the previous setter.

This matches the accepted C1 finding and focused re-review. It does not recast C1 as a missing original requirement: the preceding baseline text already covered consumption by another control call. Nor does it import `limits`, `defaults`, or other implementation helper names as public types or APIs. A reader may represent pending state as a complete value or as explicit overrides; the prose constrains the observable result, not the storage mechanism. Existing precedence, ordinary failure behavior, loud consumption reporting and per-call waiting remain consistent with the clarification.

The small addition earns its place despite the existing general rule because it identifies the specific confusion that survived implementation and initial checks. No architecture redesign, target argument, persistent setting, extra scheduler or new control tool is implied.

### 3. The live build instruction is consistent and does not invent an exercise

Frozen chapter lines 914–925 apply the same actual-CLI-package distinction to the live build example. Provider configuration, key-handling cautions, disposable-project setup and `CH02_LOG` invocation remain in place. This prevents the same library-root build mistake in a second executable instruction rather than leaving the exercise internally inconsistent.

The debugger remains Delve (`dlv`). Nothing in the diff adds a C or GDB requirement. The existing live exercise remains an exercise; no new provider campaign is requested by this review.

### 4. Preservation, voice and evidence attribution are sound

The full chapter still develops the freeze, ownership and asynchronous job argument before the exercise. The diff leaves the star/parent-chain rationale, one-shot tradeoff, PTY and cursor explanations, external-consumer proof, source-review limits and closing interruption motivation intact. Cuts replace obsolete command/build wording rather than removing useful reasoning. No Bill preface text is altered by the reviewed diff.

The additions use concrete calls, paths and consequences instead of reporting the editorial workflow to the reader. They do not turn the chapter into an acceptance report, optimize toward a word target or introduce gratuitous new abstractions.

The author report's bounded account is supported by the accepted records: entry correction plus one demonstrated state bug, not a fresh architecture failure. Historical live receipts remain historical. The chapter diff adds no claim that initial live observations were rerun on `ec4a350`, that an eight-second observation measured ninety seconds, or that rebuilt artifacts were byte-identical. Preserved second-edition demonstrations and scores are not relabeled as third-edition runs.

## Prerequisite and evidence coverage

The readiness turn fully read the following current governing documents. This formal turn recomputed their SHA-256 hashes and found all twelve unchanged from the readiness coverage record:

| Document | Complete lines read |
| --- | ---: |
| `AGENTS.md` | 1–382 |
| `.agents/skills/ensemble-coding/SKILL.md` | 1–397 |
| `docs/agentic-codebook-workflow.md` | 1–2434 |
| `docs/edition-2-carryover.md` | 1–1110 |
| `book/voice.md` | 1–632 |
| `book/chapter-writing-procedure.md` | 1–622 |
| `docs/edition-3-execution.md` | 1–53 |
| `docs/edition-3-student.md` | 1–71 |
| `docs/edition-3-readiness.md` | 1–22 |
| `docs/edition-3-carryover.md` | 1–71 |
| `.agents/skills/edition-3-author/SKILL.md` | 1–12 |
| `.agents/skills/edition-3-text-reviewer/SKILL.md` | 1–12 |

The public student brief and text-reviewer overlay were additionally reread in full during this turn. The completed Chapter 4 records were fully read during readiness and their identities checked again this turn:

| Record under `docs/edition-3-notes/` | Complete lines read | Identity check |
| --- | ---: | --- |
| `ch04-accepted.json` | 1–28 | Unchanged |
| `ch04-entry-audit.md` | 1–118 | Unchanged |
| `ch04-grader-review-r1.md` | 1–58 | Unchanged |
| `ch04-student-review.md` | 1–110 | Unchanged |
| `ch04-student-review-r1.md` | 1–96 | Unchanged |
| `ch04-code-review.md` | 1–158 | Unchanged |
| `ch04-code-review-r1.md` | 1–132 | Unchanged; lines 1–85 additionally reread for C1 and gate attribution |
| `ch04-execution.md` | 1–107 at readiness; appended 108–111 read this turn | Original prefix unchanged; handoff metadata appended |

The execution record's first 107 lines still hash to `a07fb431822c8b6e150b51892ef1840d8bae8e59b6e836cb349c867df6996593`. Its complete current 111-line hash is `5a7cd60006de5ec8db56131364626ad9c6993421963f20c16b6d761755263d89`. No changed governing material was silently carried forward.

Additional direct command-route check: `cmd/grade-edition3/main.go`, full lines 1–91, SHA-256 `539ee3e6cc7313171ab57d268803ddd4c2d49be33fb028c38fa43b3b31c12834`.

Frozen chapter reads covered 1–170, 171–340, 341–510, 511–680, 681–850 and 851–1029; report reads covered 1–48 and 49–94. These bounded reads were complete without truncation. Truncated initial readiness reads were recovered in smaller ranges before readiness was submitted. No moving author draft or private author conversation was used.

## Checks and limitations

- Performed read-only identity, coverage and exact-diff checks; full frozen manuscript/report review; comparison with completed student/acceptance records; and direct inspection of the current adapter and public brief.
- Did not run tests, builds, graders, live models, provider calls, debugger sessions, GUI work or a new sensitivity campaign. Accepted code-review and student results are attributed to their records, not asserted as newly reproduced observations.
- Did not independently rehash every nested source, binary or raw-evidence artifact referenced by those records. No claim of byte-identical rebuild reproducibility or fresh corrected-commit live evidence is made.
- Did not inspect private credentials, author history/private conversation or historical solution repair code. No child workers were used.
- This review writes only `docs/edition-3-notes/ch04-text-review.md`. No chapter, author report, source, model, guide or acceptance record was edited; no staging, commit, tag, push or restart was performed.

**Disposition:** PASS on the two exact frozen inputs above. No material teaching, technical, preservation or scope issue requires the available author revision. Parent-owned prose lint and final acceptance remain outstanding, not implied by this text verdict.
