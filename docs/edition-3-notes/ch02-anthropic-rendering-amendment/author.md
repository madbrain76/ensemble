# Chapter 2 Anthropic amendment — frozen author draft

## Authorization and scope

Bill explicitly authorized this narrow amendment, not a restarted chapter execution or normal post-student author loop. Only `/Users/bill/projects/ensemble/book/edition-3/chapter-02.md` was edited in the repository. Added three paragraphs (213 words, 23 lines including paragraph breaks) at lines 871–893, after the last paragraph of §2.6 Exhibit B and before Exhibit C. No new heading, TL;DR, existing-prose rewrite, source/grader edit, shared-instruction edit, staging, or commit. No children were started. The draft is frozen pending the parent's one independent review; no review acceptance is claimed.

## Prerequisite reads

Completed active full reads of `AGENTS.md` (382 lines), `.agents/skills/edition-3-author/SKILL.md` (12), `docs/edition-3-execution.md` (53), `docs/agentic-codebook-workflow.md` (2434), `docs/edition-2-carryover.md` (1110), `book/voice.md` (632), `book/chapter-writing-procedure.md` (622), the entire baseline chapter (1375), and `docs/edition-3-notes/ch02-completion.md` (66). Bounded reads were used; capped workflow excerpts were reread in 175-line ranges, execution guidance was recovered, and length checks identified and recovered remaining prerequisite tails. Historical phases were not restarted.

## Amendment and self-review

The addition distinguishes missing recorded evidence from actual execution outcome. It requires real correlated results, including later arrivals; otherwise it prescribes this bridge's deterministic wire-only `is_error: true` result with the exact call ID and an honest unknown-outcome statement. It forbids rerunning the tool, fabricated success, or synthetic history/Context facts. It requires valid pairing, adjacency, and order before sending, with a local error when safe rendering is impossible. Renderer-test obligations cover one missing result, a partially completed multi-call batch with intervening user text, later-result replacement, deterministic replay, and an unchanged dumped log. Actor/tool-execution machinery remains out of scope. Self-review checked these obligations and the three-noun distinction against the complete surrounding chapter; existing prose was left intact.

## Exact artifacts and checks

- Baseline SHA-256: `04b112390fac40d3dc12c2b58533573faffc717b997dffc9ce94494c97a487da`
- Frozen final SHA-256: `7bae4bccc5df6da02a049de7fa1a6d3696ca729f0dff676816fcc52f23223349`
- Exact baseline copy: `/tmp/cadence-ch02-anthropic-amendment-baseline/chapter-02.md`
- Exact frozen copy: `/tmp/cadence-ch02-anthropic-amendment-final/chapter-02.md`
- Diff: `/tmp/cadence-ch02-anthropic-amendment.diff` (one insertion hunk; no deletions).
- Byte assertion: final equals baseline plus exactly the saved insertion before Exhibit C; PASS.
- `git diff --check -- book/edition-3/chapter-02.md`: exit 0; empty output retained in `/tmp/cadence-ch02-anthropic-amendment-diff-check.txt`, with `.exit` sidecar.

Ran `go run ./cmd/lintprose -v PATH` from the repository root against the baseline and final copies above, both named exactly `chapter-02.md`. Both exited 0. Exact combined stdout/stderr is retained in `/tmp/cadence-ch02-anthropic-amendment-baseline-lint.txt` and `/tmp/cadence-ch02-anthropic-amendment-final-lint.txt`; corresponding `-lint.exit` files record exits. Neither run is warning-free:

| Soft warning | Baseline | Final |
| --- | --- | --- |
| Words (4000–7500) | 9932 | 10145 |
| Negation forms | 21, limit 19 | 22, limit 20 |
| Person gap after L625 (limit 1200) | 5327 | 5540 |

All other lint categories remained OK. Warnings were retained, not used as permission for unrelated rewrites.

## Preservation and limitations

The before/after Git-status delta is only the chapter. A SHA-256 manifest verified all 1,194 pre-existing dirty files and protected Chapter 2 record/evidence files unchanged, including Chapter 5 work. Evidence: `/tmp/cadence-ch02-anthropic-amendment-protected-before.json`, `-status-before.txt`, `-status-after.txt`, and `-verification.txt` under the same prefix. Accepted evidence and completion records were not amended.

Approximately 3.9 GiB was free before the lint invocation; no book/publication builds were launched. No API calls, credential access, runtime fixes, live regressions, or new grader tests were performed. The request's source-inspection finding motivated the amendment; this author pass did not independently diagnose an incident or inspect engine implementation. The prose states requirements and tests to perform, not implemented behavior or existing grader coverage. Parent review is the next step.
