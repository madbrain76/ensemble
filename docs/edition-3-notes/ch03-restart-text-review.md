# Chapter 3 restart text review

**Verdict: REQUIRED_REVISION.** One bounded author response; no formal re-review.

Read the complete current governing prerequisites, Chapter 3 and its diff against HEAD, the restart student/author notes, and the code-review findings and final dispositions. Checked accepted `main` source only for the relevant factual claims. This reviews prose, not the accepted code or public contract.

## Required revision

In `book/edition-3/chapter-03.md:435–441`, explicitly say that both strict-edit refusals identify the target file. A short clause such as “Identify the target file in either refusal” is sufficient; no fixed message string is required.

The new paragraph specifies zero-match advice, actual overlapping-match counts, bounded one-based locations, truncation disclosure, and unchanged bytes, but leaves file identity implicit. The first graded candidate demonstrated this exact gap: otherwise useful diagnostics omitted the path and the edit-contract grade scored 95/100. Correction round 2 added the supplied path; accepted `internal/tool/tool.go` now identifies the file in both refusals. Round 3 addressed the separate Gemini schema issue. Carry that concrete recoverability lesson into this paragraph without excerpts, fallback matching, or a universal strict-edit mandate. This does not make the student's earlier misses of explicit contracts evidence of book ambiguity.

## Otherwise accepted

- The Gemini addition correctly locates schema projection in Engine, limits the current adaptation to top-level `additionalProperties` on the six built-ins, and retains shared schemas and strict dispatch validation. It promises no general schema converter.
- Private-extras replay requires observed exact resolved identity, not matching aliases. Unknown identity correctly strips extras. The recorded live failure came from omitted root/harness identity, not a parser defect or a fourth formal code round; the prose does not misrepresent it.
- Exercise controls sit outside tool cwd; the relative log example is consistent with launching from the scratch project. Checking the actual third answer, tool log, and file effects addresses the demonstrated evidence mistakes without adding a feature.
- Historical narrative, alternative edit policies, failed attempts, and qualifications on earlier live receipts remain intact. The additions fit the existing voice; no other edits are requested.

No code/prose edits, runtime checks, grades, or provider calls were performed. The author reports baseline and final prose-lint exit status 0 with existing soft warnings; lint was not repeated.

## Root condition check

The sole author revision explicitly states that both refusals identify the target file and leave its bytes unchanged. Root verified that condition and the otherwise unchanged accepted prose; final lint exited 0 with existing soft warnings. Text closure is complete after one REQUIRED_REVISION and one author response—not a second formal PASS.
