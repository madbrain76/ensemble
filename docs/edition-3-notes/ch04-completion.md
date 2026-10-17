# Edition 3 Chapter 4 completion

Chapter 4 is accepted in student code and manuscript. Chapter 5 contract/grader entry is next under Bill’s authorization to proceed autonomously, stopping for genuine blockers or consequential scope/architecture decisions. No Chapter 5 implementation is included here.

## Accepted identities

- Student: `ec4a350f564e820d398f9125fd73a9b3b27129cd`, tree `932a5fb7305526276a1c98b53fac9b5a92e55d16`, immutable annotated tag `edition-3-ch04-r1` (peeled commit verified).
- Predecessor: `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, `edition-3-ch03-r1`. The immutable final reference was not student starting code.
- Initial candidate: `062c146cafa9b105745f6b44921160198bbaee32`, tree `840b5b3bab4daf7490b5555451bb6506b273d7b9`.
- Accepted core executable SHA256: `7342ee1d5e91791c23bc3ce32a73db51a6ed157a8d02ce07c2d3296fdf028cd3`.
- Accepted consumer executable SHA256: `7e032ea5e2692838e2a2a612acc929611e644105809d6b2bc6c6a7856fac95f7`.
- Accepted source archive SHA256: `e66358c5ecfc1a63701bd387da1cae4147b8ca95cd4886c21f34e0c7452add8c`.
- Final chapter SHA256: `ae7f6fbe1048400f857e8fb3eeefbe513d0a9bab492e1eff1d13a10ff266f37f`.
- Final author report SHA256: `f30a71586a009123578e81aa9c1e3e3cd1c76545db58ab9f499518348355e27e`.
- Independent text review SHA256: `51635df3bcb3ee7845c2e490ef2ebfa940c4dccae523f9c26ccc0f82cd10c424`.

[Acceptance receipt](ch04-accepted.json), [initial candidate pin](ch04-candidate.json), and [corrected candidate pin](ch04-candidate-r1.json) preserve the exact artifacts and qualifications. Binaries and source archives remain local artifacts outside this commit; their retained manifests are not a promise that a fresh clone already contains them.

## Entry gate and implementation

The existing Chapter 4 grading function already separates executable, core module and consumer module. Its historical CLI conflated executable and source paths, rejecting a valid public-root-library/nested-CLI layout in topology checks. A thin Edition 3 command adapter now exposes the existing split without modifying historical grader checks. An independently reviewed help-exit regression was corrected before student launch; that adapter correction was not a student correction round. See [entry audit](ch04-entry-audit.md), [initial adapter review](ch04-grader-review.md), and [accepted adapter review](ch04-grader-review-r1.md).

Both required commands run from the main repository root against the same explicit student executable:

```sh
go run ./cmd/grade-edition3 -ch 4 -json \
  -core /absolute/path/to/student \
  -consumer /absolute/path/to/student/examples/ch04-consumer \
  -consumer-case /absolute/path/to/consumer-case.json \
  /absolute/path/to/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/student-cli
```

The optional consumer case is supplied when needed and flags precede the target. Both reports must have `passed: true` and exit successfully. The adapter supports Chapters 1 and 4, not arbitrary new chapter modes. The retained 160-point checks do not exhaust the chapter’s ownership, cancellation, cursor or one-shot-setting obligations.

The fresh student extended its own accepted predecessor with a per-Ensemble job table, Agent-owned settings/tool registration, a shared ordinary-tool job dispatcher, PTY/process handling, observation waits and output cursors, explicit shutdown, and a separate public-API consumer. The ownership plan was approved before implementation. Public invocation, credential/model setup and live-checklist clarifications are recorded in the [execution record](ch04-execution.md); private grader controls, historical solutions and final-reference code were not student inputs. Context separation was not an OS sandbox.

## Recovery and code acceptance

A transport interruption and restart removed the active workers. Bill subsequently identified his phone hotspot loss; this is not evidence of an Ensemble defect. The parent recovered durable source, reports and receipts rather than rerunning completed paid work. Completed sequential tests survived, while the student had not submitted or committed. Its reopened context reconciled the records and committed the already-tested tree without implementation changes or new paid calls.

The initial debugger verifier incorrectly required two messages in one output read. Retained events showed breakpoint setup, hit, stepping, value and exit across successive reads of the same job. The parent matched those excerpts to the saved sanitized events and read the actual human CLI output. The failed verifier remains retained. A separate report-only UTC-table transcription correction was verified against unchanged raw receipts. Neither recovery nor report correction consumed a code-correction round. See the [initial student report](ch04-student-review.md).

The [full independent review](ch04-code-review.md) reproduced both grader passes but found one required C1: a partial `tool_limits` call could re-arm omitted settings already consumed by that call. Correction round 1 changed eight production lines and added a focused regression file, without redesigning the job architecture. The [student correction report](ch04-student-review-r1.md) and [independent re-review](ch04-code-review-r1.md) close C1. Reverting the fix independently caused behavioral failures for the leaked fields, then restoring it passed.

The corrected source passed full/race suites, vet, the 160/160 Chapter 4 gate and 100/100 inherited Chapter 1 gate, including actual 62-second pre-header and body waits. The reviewer tested the exact pinned executables and separately built test/mutation copies; no byte-identical rebuild claim is made. Focused local CLI and registered-consumer evidence exercises corrected settings behavior. There were no paid calls in the correction or independent re-review. One student code-correction round was used, with no acceptance exception.

## Live evidence and limits

Live receipts belong to the initial candidate’s production bytes, not fresh corrected-commit calls. The Anthropic follow-up and Gemini exercises, registered consumer, actual job/debugger use and shutdown observations are retained. OpenAI’s refusal-only behavior is not represented as successful coverage of all features. Earlier failed, refused, incomplete and over-budget attempts remain visible. Approved models included Gemini 3.8 Flash; no fallback to 2.5 was authorized.

The long-test observation was an actual eight-second test, not a measured ninety-second example. Server cleanup evidence observes the process alive before EOF and absent afterward. The debugger is Delve; no C/GDB requirement was invented. The parent’s separate keyless local chat smoke verifies running output, unseen-tail waiting and follow-up context, not live-provider correctness. See [the retained smoke receipt](ch04-orchestrator-evidence/local-chat-receipt.json) and the execution record for the actual saved artifacts.

The acceptance is bounded by these runs and source review, not exhaustive scheduler, filesystem, platform or process-failure proof. Credentials remained in approved process-local loading paths; no secret store or private runtime journal is part of this checkpoint.

## Bounded manuscript phase and closure

After code acceptance and closure of the code workers, the author changed three hunks: the real executable/core/consumer grading route, the demonstrated chained-settings confusion, and the corresponding live-build command. The paragraph teaches observable reset behavior without prescribing implementation helpers. The [author report](ch04-author-review.md) preserves rationale and exact diff. A separate root-owned reviewer first submitted readiness, then made [one formal text review](ch04-text-review.md): PASS, zero author revisions. No duplicate or nested reviewer was launched.

[Prose lint](ch04-prose-lint.json) passed all hard checks on both baseline and final chapter with the same basename. Existing length, negation-density and person-gap warning categories remain; no arbitrary word target or cosmetic review loop displaced the teaching. No PDF/EPUB build was performed.

Saved final turns were checked before joining and closing workers: student 673, code reviewer 704, author 154, text reviewer 367. The original manuscript and immutable final reference remain protected. The main-repository scope is the accepted command adapter/tests, current guidance, Chapter 4 manuscript and its text evidence. Student implementation remains in its separate repository. No push, deployment, installation or live restart is part of this completion.
