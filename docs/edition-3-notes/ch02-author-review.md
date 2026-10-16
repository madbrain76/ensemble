# Edition 3 Chapter 2 — post-student author review

Status: independent text review PASS; no author revision required or used.
Ready for parent acceptance; shared status remains the parent's responsibility.

## Boundary and evidence

This is a bounded improvement of the existing chapter after accepted student
work, not a new specification. Repository HEAD is
`5351583208988f8370937cd0494727645efa8b46`. Pre-edit chapter SHA-256:
`1a67ee9e3b6ae5f0f776a30f58908c6d53b53be2d8136ef44661728a677dc2d2`.
Frozen chapter SHA-256:
`04b112390fac40d3dc12c2b58533573faffc717b997dffc9ce94494c97a487da`.

The accepted student commit is `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`,
tag `edition-3-ch02-r1`, in `/Users/bill/projects/ensemble-edition-3-student`.
One code correction round was used. The parent verified the final source,
reviewer tail and all 133 evidence checksums; this author pass relies on that
acceptance and does not claim to repeat the implementation audit or live run.
The acceptance record retains the initial live candidate separately and says
that the final candidate's live exercises were not repeated.

Read the full current author skill, AGENTS, workflow, Edition 2 and Edition 3
carryovers, voice guide, chapter procedure, Edition 3 execution rules, chapter,
public student handoff, and these chapter records:

- [Student review](ch02-student-review.md) and [correction feedback](ch02-student-review-r1.md).
- [Code review](ch02-code-review.md) and [round-one acceptance review](ch02-code-review-r1.md).
- [Execution record](ch02-execution.md) and [accepted source/evidence identity](ch02-accepted.json).
- The narrow [official cache-field documentation check](ch02-entry-evidence/openai-cache-field-doc-check.json), used as evidence, not a new contract.

Oversized tool reads were recovered in bounded ranges rather than treating
truncated output as a completed read.

## Changes and reasons

1. **Build the actual CLI and grade the same executable twice.** The student
   inherited a root library with `./cmd/ch01`. A directory target invokes the
   historical root build and is not an executable for that layout. The exercise
   now shows an explicit build-package and absolute executable path, followed by
   the existing Edition 2 Chapter 2 and Edition 3 Chapter 1 commands against that
   same file. A root `main` still uses `.`; no flattening or package rename is
   required. The live example likewise builds the nested CLI into the existing
   `ch02` binary name. Successful exits and `passed: true` matter, not merely
   100 points: inherited request-lifetime checks remain in force. No timing
   inventory or grader rewrite was added.
2. **Name the optional OpenAI Chat Completions cache-field path.** Usage already
   named the counters and subtraction formula but did not locate them. A short
   paragraph supplies `usage.prompt_tokens_details.cached_tokens` and
   `usage.prompt_tokens_details.cache_write_tokens`, links the current official
   Chat Completions reference, and keeps absent counters at zero. The parent's
   recorded HTTP-200 documentation check supports the optional field; the
   prompt-caching guide's Responses example is not substituted for this
   chapter's interface. The historical live numbers and normalization formula
   are unchanged. “Current” qualifies the new documentation link; no claim is
   made that this optional field was documented in December 2025.

Only three chapter locations change: Usage, the grading commands, and the live
CLI build paragraph. All other chapter bytes are preserved by exact-anchor
replacement. No schemas, protocol behavior, renderer order, defaults, point
weights, historical attributions, rationale, or evidence qualifications change.

## Deliberately not changed

- R1 chronological ID-less Gemini pairing, R2 measured failed-reply usage, and
  R3 retaining multiple OpenAI opaque fields were already explicit requirements.
  The correction feedback calls for better test coverage, not more prose.
  They are not converted into a new regression inventory or architecture.
- No Responses API requirement, broader provider rewrite, tools, GUI, streaming,
  retries, media exercise, or new feature/default.
- No score repair: the chapter's original 92/100 FAIL and accepted media
  exception remain history, distinct from the accepted Edition 3 code result.
- No speculative polish, arbitrary word target, audit timing dump, or new
  claims about live signature/media coverage.
- Bill's preface is untouched. No application, grader, student, reference,
  original-book, other-chapter, procedure, or shared-status file is edited.
  Unrelated untracked items are preserved; no commit, push, install, restart,
  build, grader execution, or live API request was performed in this pass.

## Validation at freeze

- Chapter words (`len(text.split())`): **11,257 → 11,411** (+154).
- Fenced code blocks: **22 → 23** (44 → 46 fence markers), balanced; the new
  block isolates the student-directory build from book-root grader commands.
- Both chapter Markdown link targets were enumerated. The relative target
  `../../solutions/evidence/ch02/live.txt` resolves on disk. The only new
  external link is the official Chat Completions reference; its HTTP 200 and
  optional nested field are supported by the recorded parent check, not by a
  new network request from this author.
- The named grader commands' `-json` flags were checked in their entrypoints.
  Accepted evidence supplies the executable-target behavior; no grader was run.
- Before writing, an in-memory unified diff was inspected. Failed exact-anchor
  and incorrectly transcribed checksum assertions wrote nothing. The final
  write required the verified complete baseline SHA-256 and unique anchors.
- `git diff --check` passes. The chapter diff is three hunks, 28 inserted and
  seven deleted lines; no other tracked file is modified.
- Bill's preface is byte-identical to HEAD; SHA-256:
  `915cf8c531169c5aacc4cb81bd743eeba8a181d20bc779d936ac5f7ab9998d96`.

## Independent review

Reviewer: independent read-only child `ch02-text-review`.
Verdict: **PASS — no blocking findings; no author revision required.**
One review completed and joined; zero author revisions used. The reviewed and
final chapter SHA-256 remains
`04b112390fac40d3dc12c2b58533573faffc717b997dffc9ce94494c97a487da`.

The reviewer read the full current governing documents, chapter, author report,
student/code reviews, execution/acceptance records, and narrow official-field
evidence. Its findings by location:

- Lines 729–733: exact Chat Completions nesting, current optional-field
  qualification, omitted-zero semantics, and unchanged normalization are correct.
- Lines 1260–1278: nested CLI and root-main alternative, same explicit executable
  for both existing graders, and boolean/exit acceptance are correct; no new
  contract, layout mandate, or timing inventory is introduced.
- Lines 1308–1319: the live build remains consistent with the inherited layout,
  root-main alternative, existing binary name, and three-provider exercise.

It independently confirmed the three-hunk diff, word/fence counts, clean
whitespace check, local chapter link, unchanged preface, faithful acceptance/live
provenance, and absence of scope creep. It made no edits or prohibited runs.
All seven report-relative evidence links also resolve on disk. No second review
or speculative revision was requested.
