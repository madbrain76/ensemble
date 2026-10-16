# Edition 3 Chapter 3 — independent text review

## Verdict: REVISE

One formal, bounded text review of the frozen inputs below. Two actionable text corrections are required. Neither reopens accepted code, asks for a new feature or test suite, nor creates an optional polish cycle. The author has at most one revision; root owns closure and any commit/tag/status work.

## Exact inputs and coverage

- Main repository HEAD: `f7530a6643ca00bd2f59a8da60021c567cd7d174`.
- Frozen working `book/edition-3/chapter-03.md`: SHA-256 `796ea212fb7acead95083e7942561eadfd3e1f566acdbf6b43e00b74b22177e6`; read all 918 lines.
- Frozen `docs/edition-3-notes/ch03-author-review.md`: SHA-256 `12f31a4bc8c57b8813daf8ce051c03f400ed4d52d0c55bca3b72bc77198375ac`; read all 151 lines.
- Committed Chapter 3 baseline read during readiness through `git show HEAD:book/edition-3/chapter-03.md`: SHA-256 `2366f74bdb873af57fe95e0c7c0b040187365a8db7efc09f5e63ce1f3c59c801`, 895 lines.
- Read the entire exact three-hunk chapter diff: 28 additions, 5 deletions. Only Chapter 3 was a tracked working-tree change before this report.
- Prerequisite readiness included complete reviewer/author overlays, shared E3 execution procedure, AGENTS, coding skill, workflow, E2 carryover, voice, chapter-writing procedure, current E3 carryover/readiness, both Chapter 3 student reports, both code reviews, acceptance receipt, and execution record. Truncated prerequisite reads were recovered. Current reviewer/execution instructions were reread for this pass; the execution output truncated in a combined command was recovered separately.
- Read Chapter 2's public representation, replay/provenance, and acceptance sections in context (`chapter-02.md:255–435,478–575,920–975,1125–1180`). Inspected accepted student source only for the specific representation and build-layout questions, not a code-acceptance rerun.
- Student accepted commit: `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`. Independently peeled annotated tag `edition-3-ch03-r1^{commit}` to that exact commit. The unpeeled annotated-tag object is not the source commit.

## Required corrections

### 1. State the signed-text invariant without requiring the student's particular Part fields

**Location:** `book/edition-3/chapter-03.md:221–226`, the dispatch-policy paragraph beginning “Do not mistake signed ordinary text for hidden thinking.”

The failure story and invariant are useful and accurate: a Gemini signature does not make ordinary text hidden, hidden reasoning does not explain a tool call, and replay-only metadata must retain its exact source/part association. However, the instruction “Keep that text visible, with its signature in the same part's `Opaque` and its exact source in `From`” prescribes the accepted student's tagged-union field layout as if Chapter 2 universally supplied those fields on text parts.

Chapter 2 explicitly allows either a small interface with private concrete variants or a tagged union. Its displayed `TextPart` has only `Text`; `OpaquePart` has `From` and `Data`; the displayed `Opaque` field belongs to `ToolCallPart`. Its accompanying per-part binding explanation expressly describes tool-call metadata. Consequently a reader who followed the permitted concrete-variant example has no text-part `Opaque` or `From` fields to use. The student's accepted `internal/common/common.go` does have both fields on its tagged `Part`, but that implementation choice is not a new public prerequisite for every reader.

**Required bounded revision:** retain the visible-text/hidden-content distinction and exact text-part/source signature binding, but express the obligation independently of struct field names. If the field names remain as an example, identify them as the accepted student's tagged representation rather than a required inherited text-part API. Do not introduce a new representation, History layer, policy framework, or earlier-chapter change. This is a clarification of the same behavioral obligation, not a request for another seam.

### 2. Apply the nested-CLI correction to the final real-work build instruction

**Location:** `book/edition-3/chapter-03.md:881–886`, “Then use it for real work,” especially line 885:

```sh
go build -o "$HOME/bin/ch3agent" .
```

The new setup and grader instructions correctly build `./cmd/ch01` and tell readers to substitute their actual CLI package. This remaining root-package command contradicts that guidance at the point where the reader prepares the real chat executable. The accepted source's root `ensemble.go` is `package ensemble`, while `cmd/ch01/main.go` is `package main`; building the library root is not building the executable that the following `ch3agent chat` command needs. This is a concrete continuation of the actual student's nested-CLI workflow problem, not cosmetic consistency.

**Required bounded revision:** use the actual CLI package in this final build command as well (the same `./cmd/ch01` example, with the existing actual-package qualification applying). No new flag, build system, or grader mode is needed.

## Findings that do not require further changes

- The two added grader commands correctly name Edition 2 Chapter 3 and Edition 3 Chapter 1, use the same explicit executable, and require both `passed:true` and process exit 0. Narrow source inspection of `cmd/grade-edition2/main.go` and `cmd/grade-edition3/main.go` confirms the dispatch and failure-exit contract. No E3 Chapter 3 mode is invented.
- The grader-limit paragraph correctly leaves selected-policy refusals/retries to the student's local behavioral proof. It does not turn an old harness score into policy evidence.
- The chapter's typed provider-wire, fake-server, and real-use obligations already cover tool declarations sufficiently for this bounded author pass. No additional `Schema` field inventory or new test/coverage requirement is requested.
- The actual failure story earns its place: it explains why preserving opaque metadata alone is insufficient if ordinary visible text is misclassified. The initial Gemini failure is historical evidence, not an open exception for the accepted R1 code.
- The chapter retains the useful one-loop ownership explanation, synchronous-tool scope, honest shell authority, and visible-explanation rationale. No optional style rewrite, chapter-length target, streaming/tool-async work, model-default/cap change, or new policy selection mechanism is requested.

## Mechanical checks and evidence limits

Independent static checks performed in this pass:

- SHA-256 checks reproduced both frozen input hashes and the baseline chapter hash.
- `git diff --check -- book/edition-3/chapter-03.md` exited 0.
- A line-based fence scan found 24 fence delimiters (12 closed blocks), no unmatched fence, and no mismatched delimiter.
- Both chapter Markdown link targets exist: `../../edition-1/book/exhibit-ch03-tools.md` and `../../solutions/evidence/ch03/live.txt`. This is a local target-existence check, not a new validation of the historical live log's contents or current provider behavior.
- `book/preface.md` and `book/edition-3/preface.md` are byte-identical to their respective HEAD versions and to each other; SHA-256 for each is `915cf8c531169c5aacc4cb81bd743eeba8a181d20bc779d936ac5f7ab9998d96`.
- Read-only inspection confirmed the accepted student's public library root and `./cmd/ch01` executable layout. No build, grader, test, or provider invocation was performed for this text review. No prose-lint threshold was imposed.

The author report and existing code/student records are evidence with their recorded limits, not new executions by this reviewer. Corrected Gemini 3.8 live evidence concerns two nonblank human answers from the same pinned binary; raw stdout was not retained, so the observation is qualified by launcher checks plus retained History/project audit rather than a replayable raw terminal transcript. Earlier Anthropic/OpenAI live evidence belongs to the prior revision. The recoverable search error is not an all-tool-results-success claim. Nothing in this verdict inflates those observations or reopens accepted code.

A first source lookup used an incorrect `common/types.go` path; the correct accepted `internal/common/common.go` was then read. An attempted author-check-file lookup found no file at the guessed path and supplies no evidence; the mechanical results above come from independent commands in this pass. A Git pager opened during the first diff command; it was terminated, and the complete diff was reread with paging disabled.

Only this review report was written. No chapter, author report, guide, preface, application, grader, student source, or reference solution was modified. No children, paid provider calls, commits, installation, restart, or push were used. This completes the sole formal pass; the two corrections above are the complete requested revision set.
