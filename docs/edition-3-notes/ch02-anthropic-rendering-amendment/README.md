# Chapter 2: Anthropic missing-result rendering amendment

Bill requested a clarification after reporting that an interrupted tool call with
no recorded result repeatedly produced an invalid Anthropic request. This is a
later, separately reviewed teaching amendment, not a replacement for the original
Chapter 2 acceptance record or proof of a deployed runtime repair.

The change adds three paragraphs inside existing §2.6 Exhibit B. No new section
or TL;DR was added. It teaches legal disposable wire projections, honest unknown
outcomes, real-result precedence, local rejection of unsafe projections, and
specific offline renderer regressions. It does not authorize tool re-execution
or synthetic facts in history. The exact baseline, final hashes, and lint outputs
are in [validation.json](validation.json).

- [Author report](author.md), preserved verbatim.
- [Independent text review](review.md): PASS, no author revision.
- Both prose-lint runs exit 0; existing soft warning categories remain.
- No runtime, grader, historical/reference implementation, or student source changed.
- Existing tests and live receipts do not establish coverage of these new cases.
- The prior accepted chapter remains recoverable from Git; its old whole-file hash
  does not identify the amended text.

Anthropic's [official tool-call reference](https://platform.claude.com/docs/en/agents-and-tools/tool-use/handle-tool-calls)
confirms the wire rules. The missing-record placeholder is explicitly this
bridge's recovery policy, not a vendor-mandated representation of tool failure.
