# Edition 3 Chapter 10 completion

Chapter 10 is accepted in code, operational exercises, and teaching. The accepted student is `729bc919e734405307b4ac0060e65c40e33652e5`, tree `025ef3744b42c8f9daec34e76661edc40583ca6a`, immutable annotated tag `edition-3-ch10-r1`. There were **zero formal student code-correction rounds**, one formal text review requiring revision, and one author revision verified directly by root. All code and text workers were joined and closed.

## What established acceptance

- [Code decision](ch10-code-root-decision.json): independent source assessment, real persistence/replay and legacy GUI identity controls, and a compile-valid atomic-write mutation that failed behaviorally before exact restoration passed. Invalid reviewer fixtures remain recorded; they were not student defects.
- [Final student gates](ch10-student-final-reconciliation.md): Chapter 10 and Chapter 1 passed on the same frozen CLI, including actual separate 62-second delayed-header and delayed-body phases. Independent review carried its own Chapter 10 result by exact identity and ran its own same-CLI lifetime gate. Compile-only module checks remain qualified.
- [Operational verification](ch10-final-operations-parent-verification.json) and [final exercise report](ch10-operational-harness-v2-evidence/tls-recovery-v1/final-report.md): an owned muted browser demonstrated GUI save/restart, stable old cards and immediate new output; a real Claude CLI exercise read a file, saved through EOF, restarted and used its retained result. No application source changed after the freeze.
- [Text review](ch10-text-review.md), [sole revision](ch10-author-revision.md), and [root text decision](ch10-text-root-decision.json): the final manuscript teaches the actual nested CLI and separate-artifact grading commands, preserves persistence/replay reasoning, and correctly attributes the earlier default-load mismatch to the fixture's terminal newline. The two unsupported body first-person statements were removed. Final chapter SHA256: `c156708aae6fde042cb34340da239279702041b245e096efb03a5d19a7178aca`.

## Failures and limits retained

The original 80/100 grade remains failed. Its synthetic prompt contained a terminal newline that the unchanged exact expectation omitted; the one-line fixture correction did not require trimming student prompts or changing production. Earlier numerical, setup, harness, wrapper, and reviewer failures remain in their original seals and reconciliations.

The initial operational preparation overran scope and was rejected. Subsequent browser setup and harness assertion failures were preserved. The final GUI result uses its own seal and is not replayed under the paid recovery's configuration. PNGs were captured and hashed, not human-viewed; no listening claim is made.

Paid exercise accounting is three successful HTTP 200 responses, 5,616 input and 139 output tokens, zero cache usage. Across the failed transport attempt and explicit TLS recovery there were four reserved slots and two scoped credential opens. The first exception detail was lost: an empty Python trust store was subsequently observed and a pinned public CA bundle used with hostname/certificate verification enabled, but that does not prove the lost exception's sole cause. No further calls were made for cleaner evidence.

The sole author revision passed prose lint with a disclosed soft person-gap warning. Text acceptance is root verification of the single REQUIRED_REVISION verdict, not a second formal PASS. Finite checks and limited credential scanning are not universal correctness or security proofs.

## Delivery and next boundary

[Acceptance](ch10-accepted.json), [execution chronology](ch10-execution.md), [artifact portability](ch10-artifact-portability.md), and [closure validation](ch10-closure-validation.json) preserve exact identities and qualifications. Student source remains in its separate repository/tag. The accidentally placed source TAR, binaries, private captures, browser profiles, PNGs and forbidden credential filenames remain local; no private settings or credentials are publication material.

The closure commit is restricted to this chapter, its adapter/fixture/public guidance and inspected text evidence. Concurrent runtime and publication work is excluded. No push, publication rebuild, installation or live restart was performed. Chapter 11 begins from this accepted student tag after its own contract/grader entry check; the current adapter must not be invoked blindly for Chapter 11.
