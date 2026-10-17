# Chapter 5 completion — Edition 3

## Accepted boundary

Chapter 5's code and text are accepted. The accepted student commit is
`f60ebf3aad2ece4251face9d95c106dd667abcd6`, tree
`3de5d0d24dd9f593ec017fceccfbff1b209c5d57`. The immutable annotated
`edition-3-ch05-r1` tag object is `477fb958d4d65b7629501cffcfd6247cda4a5d36`.
It extends the accepted Chapter 4 student, not the frozen final reference.
The final chapter SHA256 is
`475d807b1fcbd9e3bbb475afc08af7baea6cae438e77d0c404c9030782e7aba9`.
Exact artifacts, review pins and lifecycle receipts are recorded in
[ch05-accepted.json](ch05-accepted.json); the detailed chronology is
[ch05-execution.md](ch05-execution.md).

One formal code-correction round was used. Bill's subsequent new wake-up
requirement was an explicit contract addition, not a retroactive student
failure or another correction round. One normal post-student text review
accepted the final two-paragraph clarification unchanged; zero author
revisions were requested. All code and text workers were joined, their
saved final turns verified, and then closed. No acceptance exception remains.

## Entry and code acceptance

The [entry audit](ch05-entry-audit.md) demonstrated that the historical CLI
conflated executable and source topology. A thin, independently
[reviewed adapter](ch05-grader-review.md) exposes the unchanged Chapter 5
API with explicit executable, core, consumer, workflow and optional model-map
identities. The Edition 3 command supports only Chapters 1, 4 and 5;
a numeric `-ch 6` must not be assumed to work. A public root library with a
nested CLI remains valid. Both the chapter grade and same-binary Chapter 1
lifetime gate require `passed: true` and exit zero, not merely a score.

The initial `00489b05` implementation scored 100 in both gates, but
[independent whole-code review](ch05-code-review.md) found three real
lifecycle defects: a pending hint stranded after a text-only HTTP reply,
a queued `Ask` delaying caller cancellation, and shutdown waiting for an
ordinary Go handler before stopping another managed process. The student
corrected these at `6c4561cb`; the [independent re-review](ch05-code-review-r1.md)
closed all three under the then-current contract. Original-source behavioral
negatives and restored passes were retained, including the favorable map
ordering that did not reproduce the shutdown failure. Setup and compilation
failures were not counted as sensitivity evidence.

Bill then authorized the [job-wake-up amendment](job-wakeup-amendment.md),
committed separately in `4f1701a76f57ec57d362d4f1a30264f30b65ddf8`.
It makes foreground observation interruptible by hints and addressed
questions without cancelling the underlying work, preserves common job
handles and capped cursor results, and gives each undispatched batch-tail
call an honest, correlated not-executed error. Completed and started work
retain their true results/state. The actor can answer the question and
resume the same job; no future fleet framework was imported into this chapter.
That amendment had its own explicit user-authorized text review before
implementation. It is distinct from the normal post-student pass below.

The [new implementation report](ch05-wakeup-student.md) preserves a genuine
intermediate PTY failure: treating every arriving byte as observation
completion made `send_input` return too early. The accepted distinction is
immediate already-available output for resumed `wait_for_job`, versus
notification/rechecking of pattern, terminal, delay and wake conditions
for a blocked observation. Parent assistance and the superseded freeze are
recorded, not erased.

The [final independent code review](ch05-wakeup-code-review.md) accepted
`f60ebf3` without findings. On the exact pinned final CLI, Chapter 5 scored
100/100 with exit 0 in 160.028 seconds, and Chapter 1 scored 100/100 with
exit 0 in 127.273 seconds, including actual 62-second header/body phases.
The reviewer also ran pristine-archive race/vet and module gates, and
independent real-public-interface controls for wake-up, true batch results,
metadata delivery, caps/cursors, resumed waits, process input, shutdown,
and inherited behavior. Original production compiled and then failed the
intended behavioral controls; restored code passed. A tentative one-shot
limit concern did not reproduce and is not a finding. Finite schedules and
source review are evidence, not a proof of every possible interleaving.

## Live evidence and its limits

All twelve paid requests belong to the initial `00489b05` production bytes:
four Anthropic human-interaction calls, seven OpenAI Chat Completions workflow
calls, and one Gemini image call. They are not fresh paid demonstrations of
`6c4561cb` or `f60ebf3`. Final wake-up and adjacent behavior have independent
local public-interface and actual CLI/workflow evidence using controlled
HTTP boundaries; no extra paid calls were made for the amendment.

The original human launcher failed its blanket sentinel assertion because
it matched the model's sentence denying a leak. The actual process output
contained only the intended start/finish/exit report; its CLI exited zero,
and a new turn answered before the old job ended. The failed launcher
summary remains failed. A local verifier correction and retained evidence
reconciliation established the exercise without repeating paid requests.
The real workflow changed the manuscript and produced a substantive review;
the image answer correctly identified left red/right blue. Unsupported
video was refused locally with zero HTTP requests. An attempted mocked
route control remains explicitly uncredited; actual executable/HTTP-boundary
controls provide separate evidence. No native GUI, runtime installation,
provider-wire-wide correctness, or exhaustive media support is claimed.

## Post-student teaching and text review

The [author](ch05-author.md) added exactly two paragraphs: pending hints at
a text-only response boundary, and notification versus observation completion.
They add 15 lines and 134 lint words while preserving every baseline byte.
The [single text review](ch05-text-review.md) accepted them unchanged after
reading the complete chapter and underlying reports. Same-basename baseline
and candidate prose lint both exited zero; existing soft warnings were not
used to delete useful failure stories. The actual additions are in §5.5;
the author's §5.6 locator is a disclosed report-label error, not a reason
to rewrite the frozen evidence. Bill's preface and the already accepted
wake-up amendment remain unchanged.

Author/reviewer validation setup errors and corrected transcriptions are
preserved. The parent also initially omitted the externally indexed review
report from an exact-set assertion; every hash matched, and the corrected
159-file comparison passed without changing artifacts. Author final turn
187 and text-reviewer final turn 224 were saved and joined before closure.
Code student/reviewer final turns were 1825 and 903 respectively.

## Scope, preservation and next step

The closure scope is the five explicit chapter/entry-guide/adapter files and
named `ch05-*` reports/evidence. It excludes the old `.agents/chapters/code/`
convenience source copy, unrelated untracked work, private journals, settings,
keys, local archives and executables. The accepted named source snapshot in
student evidence is an indexed student artifact, not the final reference.
Raw failed logs and patch-format whitespace remain byte-preserved. Ignored
text evidence is included only as an inspected named subset; ignore rules
are unchanged. Exact staging and committed-blob checks are retained in the
parent's closure receipt, with the normal staged credential hook enabled.

The 231 immutable final-reference files and 47 original manuscript files
match their copy-baseline hashes. All 166 entry inputs were checked; the
known changed input is the deliberately amended Chapter 5 manuscript
(before the final progress-marker update). Grading code remains at its
independently accepted pins. Local links resolve. Credential-pattern scans
are limited checks, not comprehensive security proof.

No push, installation, deployment, live restart or publication build occurred.
Chapter 6 entry is next, from accepted `f60ebf3`/`edition-3-ch05-r1`, with a
fresh contract and supported-grader check before student implementation.
