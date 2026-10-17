# Chapter 2 Anthropic missing-result amendment — independent text review

## Verdict: PASS

No substantive findings and no author revision requested. This is one formal independent review of Bill's authorized, narrow amendment to an already accepted chapter, not a reopened student cycle or acceptance of a runtime fix.

## Frozen subject and scope

Reviewed `/Users/bill/projects/ensemble/book/edition-3/chapter-02.md` at SHA-256 `7bae4bccc5df6da02a049de7fa1a6d3696ca729f0dff676816fcc52f23223349` against `/tmp/cadence-ch02-anthropic-amendment-baseline/chapter-02.md` at `04b112390fac40d3dc12c2b58533573faffc717b997dffc9ce94494c97a487da`.

The author report `/tmp/cadence-ch02-anthropic-amendment-author.md` matched supplied SHA-256 `059fa0f3e25e16f07b2ed5d423d37cc5dbbec29e4ffd485a7a011780a9ac4079`.

Independent byte comparisons verified exactly one insertion: 1,422 bytes, 23 lines, 213 whitespace-delimited words, three paragraphs at current lines 871–892, within §2.6 Exhibit B immediately before Exhibit C. Removing those inserted bytes reproduces the baseline exactly. The supplied `/tmp/cadence-ch02-anthropic-amendment.diff` contains exactly those added bytes and no deletion. There is no new heading or section. The frozen final lint copy also hashes identically to the reviewed chapter.

## Required context

Fully read the current reviewer and author skills, repository `AGENTS.md`, complete `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, shared `docs/edition-3-execution.md`, and the entire current Chapter 2. Capped reads were recovered in bounded ranges. The parent confirmed that “sharedguide” denotes the execution guide, not another file.

Also read Chapter 2's completion record, initial and revision student reports, and revision code review, together with the amendment author report, exact diff, saved primary-source excerpt, and existing paired lint artifacts. Historical instructions and accepted receipts were treated as historical evidence, not authority to restart a prior phase. No child agents were created.

## Semantic and instructional assessment

- **History remains evidence.** The insertion makes an absent result an absence of a record, not evidence that execution failed, succeeded, or never occurred. It confines recovery to a disposable wire projection and explicitly forbids mutating history or Context or appending synthetic facts to either. This preserves the chapter's existing distinction between facts, reduced context, and provider rendering.
- **The placeholder is local policy, not a vendor mandate.** “Use this bridge's recovery policy” clearly owns the deterministic placeholder convention. The exact call ID becomes `tool_use_id`; the fixed text states “No result was recorded; execution outcome unknown.” The next sentence explicitly explains that `is_error: true` describes the missing record, not a diagnosed tool failure. Together these sentences adequately prevent the bridge's recovery convention from being silently attributed to Anthropic or confused with a historical tool-execution error.
- **Real evidence takes precedence.** The first paragraph requires searching for a real correlated result, including one recorded later, and using it whenever present. The third requires that a later real result supersede the placeholder on a fresh render. “Do not rerun the tool to repair a transcript” removes a dangerous implication that protocol repair authorizes another execution.
- **Wire validity is actionable.** All results must follow their assistant `tool_use` message immediately and precede intervening user text. The pairing, adjacency, and block-order checks apply before sending; unsafe histories must fail locally rather than produce malformed HTTP. This gives the reader both a recovery rule and a refusal boundary rather than suggesting that every malformed history can safely be rearranged.
- **The offline obligations are concrete.** The prose calls for a single missing result, a partially completed multi-call batch, intervening user text, later-real-result precedence, deterministic replay, and an unchanged dumped log. These are renderer tests over supplied history, not a demand for actor orchestration or tool execution. They cover the requested cases without importing later-chapter machinery.
- **Placement and usefulness fit the amendment.** The existing Exhibit B discusses Anthropic result ordering. The added three paragraphs extend that same wire-boundary lesson to incomplete history, then give bounded tests. They teach the missing behavior without adding a section, changing the accepted narrative, or launching an optional polishing campaign.

## Primary-source check and evidence limits

Examined `/tmp/cadence-ch02-anthropic-doc-check.json`, which records HTTP 200 for:

`https://platform.claude.com/docs/en/agents-and-tools/tool-use/handle-tool-calls`

Its excerpt defines `tool_use_id` as the ID being answered, requires immediate result adjacency and results before text, and defines `is_error` as whether tool execution resulted in an error. That last definition does **not** itself prescribe a synthetic missing-record placeholder. The insertion's explicit bridge-policy ownership, unknown-outcome text, and missing-record explanation provide the necessary distinction. PASS does not assert vendor endorsement of this recovery convention.

A separate reviewer attempt to refresh that public URL with Python urllib returned HTTP 403. Therefore the successful HTTP-200 evidence cited here is the supplied saved check, not a claimed successful fresh reviewer fetch. No paid/model calls were made.

The motivation is limited to Bill's reported invalid-RPC problem and the author report's source-identified gap. This review does not independently reproduce or comprehensively audit the incident. The insertion makes no claim that the runtime has been fixed, that these new regression cases already exist, or that the current grader covers them. Existing live demonstrations and accepted implementation receipts remain evidence for their original, narrower scopes; their unchanged presence is not new missing-result coverage.

## Existing lint evidence

Inspected, without rerunning prose lint:

- `/tmp/cadence-ch02-anthropic-amendment-baseline-lint.txt` and `.exit`
- `/tmp/cadence-ch02-anthropic-amendment-final-lint.txt` and `.exit`

Both runs use `chapter-02.md` as the tested basename and both exit files contain `0`. Both report the same three soft-warning categories: words, negation forms, and person gap. Counts move from 9,932 to 10,145 words, 21 to 22 negation forms, and 5,327 to 5,540 words for the flagged person gap. No new warning category appears. These unchanged categories do not warrant an unrelated prose revision in this authorized narrow review.

## Preservation and disposition

Only this requested report was written by the reviewer. No manuscript, runtime/source, grader, guide, accepted record, Chapter 5 file, or other repository artifact was edited. Pre-existing dirty/untracked work was left alone. No graders, application tests, or paid live probes were run; existing lint results were reused as requested.

Post-reading byte verification still matched all three supplied hashes and the final lint copy. The exact-diff verification passed. This PASS applies only to the specified frozen text amendment; runtime remediation and implementation/regression coverage remain separate work. No revision is required.
