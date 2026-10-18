# Edition 3 Chapter 9 completion

Chapter 9 code, operational exercises, and text are accepted. Repository closure follows this record; the later commit receipt records the actual commit/hook outcome, not this precommit document.

## Accepted construction

- Student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Commit: `251ce88fb1677897dbcc6b33661096c5cb39cd8e`; tree `a2add41c306d567c05d98ae760911bc13c18b75a`.
- Immutable annotated tag: `edition-3-ch09-r1`, object `eaa70a1a1adb0396586d735211620e5d839cd06f`, peeling to that commit.
- 132 source/support files; accepted predecessor is Chapter 8's `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`, not the protected final reference.
- Zero formal code-correction rounds. Initial implementation genuinely failed ownership parity at 85/100 when Agent imported its sibling registry implementation. The student corrected composition through the shared interface before independent whole-candidate review; root did not patch student code or waive the requirement.

The accepted implementation preserves per-Agent skill ownership, complete graph validation before activation, creation-time frozen primary instructions, dynamic tool-result instructions and provenance, actual hidden-tool admission checks, and pending unload without premature revocation. Finite checks are not universal proof. See [student report](ch09-student-review.md), [independent code review](ch09-code-review.md), and [code acceptance](ch09-code-root-decision.json).

## Evidence and limitations

The student and independent reviewer each passed Chapter 9 and Chapter 1 against their own same frozen CLI, including genuine separate 62-second delayed-header and delayed-body phases. The review's actual Chapter 9 result was 100/100, passed true, exit 0; Chapter 1 completed in 124.087434 seconds. Root race, GUI race, vet and module checks passed; modules with no tests are compile-only evidence. Independent pending-unload mutation compiled and failed behaviorally, then exact restoration passed. Student hidden-dispatch and first-provenance controls likewise exposed real behavior; earlier panic/compile setup failures remain uncredited.

Review setup failures remain visible: provider-environment pollution, malformed media map, resource/copy-origin deviations, and transport errors. The valid saved PASS at reviewer final850 preceded a later queued-turn failure864; root verified both rather than inventing a new finale or repeating gates. All code workers' saved tails were joined and retained before closure.

[Operational evidence](ch09-operational-exercise.md) closes the remaining chapter exercise:

- One muted owned Chrome/local-provider GUI run, six fake requests, actual paired tool cards/results, `note.txt` and post-unload `after-unload.txt` containing `SKILLS_OK`, frozen primary instructions and cleanup. The sealed controller exited 0; the original outer wrapper exited 1 because it looked for cleanup at the wrong path. Root verified the actual receipt without replaying the browser.
- One Claude Sonnet 4.6 CLI attempt, five HTTP 200 requests, unchanged native 1024-token maximum per request, 8,358 input/729 output tokens, zero cache, one scoped credential-file read, no retries. Actual load/write/read/pending-unload/final response were correlated. No paid GUI or three-provider campaign occurred.
- The unsuccessful optional decoded-log byte comparison remains uncredited. Usage and the sealed forwarding checks do not magically turn that check into proof of upstream byte equality.
- Screenshots were captured, not directly viewed. GUI assets are byte-identical to accepted Chapter 8, so its qualified visual assessment carries by identity; Chapter 9's fresh DOM/typed-message evidence establishes its new content and interactions. No native-speech or human-listening claim is made.

Root's first execution task used the wrong sealed-controller interface; the student stopped prelaunch and the corrected version is retained. Transport recovery resumed the same context and unused paid attempt, not a duplicate browser exercise. The private-data/key audit is scoped, not a universal host/network privacy proof.

## Teaching acceptance

The author updated the separate-artifact Edition 3 grading recipe and explained the demonstrated distinction between per-Agent state ownership and sibling-package imports. The ONE formal [text review](ch09-text-review.md) returned REQUIRED_REVISION, not PASS. The same author's ONE [revision](ch09-author-revision.md) corrected the real `./cmd/ch01` target, supported `chat` invocation, and accepted `Config.Skills`/`SkillSpec` construction example. It explicitly retracts the original false claim that the nonexistent `./cmd/ensemble` path had been checked.

Root directly verified R1–R3; no second formal review or later PASS was invented. Current paired same-basename prose lint exited 0 with the existing soft person-gap warning retained. Original drafts/reports/seals remain historical evidence. Final chapter: 20,888 bytes, SHA-256 `e2c82e17dc2cdeb4c580afe182ecef60bf6393fe10e234f56a7a255bbc9a5fc0`. See [root text decision](ch09-text-root-decision.json) and [revision verification](ch09-author-revision-parent-verification.json). Both text workers are joined and closed.

## Repository boundary and next chapter

Only the Chapter 9 manuscript, accepted thin adapter/main tests, public invocation guides, and inspected Chapter 9 text evidence belong in this closure. Concurrent runtime, publication, voice, linter and main README work are excluded. No push, installation, live restart, publication rebuild, or new operational exercise is part of closure. The credential hook and exact committed-blob checks must run during the actual scoped commit. [Portability](ch09-artifact-portability.md) explains local artifacts and historical pin mappings.

Next: inspect Chapter 10's current contract and existing grader route afresh from accepted `251ce88`, before launching a student. The current Edition 3 adapter supports only 1/4/5/6/7/8/9; no Chapter 10 route is presumed. Chapter 25 remains historical-only; Chapter 26 follows accepted Chapter 24. No memory migration or crossover chapter has been implemented by this work.
