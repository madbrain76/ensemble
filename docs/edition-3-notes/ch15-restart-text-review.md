# Chapter 15 restart — final text review

**Verdict: REQUIRED_REVISION.** One bounded full-text pass; one concrete finding. The current student-first, ten-round/one-prose-pass instructions govern, not the historical campaign rules.

## Required correction

**§15.7, `book/edition-3/chapter-15.md:831–833`: distinguish calculated cost from displayed precision.** The sentence says provider usage and the DOM meter agree “on $0.0755091.” In the retained live evidence, `usage-4.json` records calculated `cost` and wire `cost_usd` of **0.0755091**, but `dom-turn-4.json` records visible `meter` **$0.0755** and tooltip subtotal **$0.075509**. These agree after formatting, not at the precision claimed for the visible meter. Provider usage supplies token counts; the dollar figure is calculated from them.

A sufficient replacement, keeping the existing restart and qualification sentences, is:

> The retained October 22 Sonnet 4.6 run returned four correct answers. Provider usage yielded a calculated cost of $0.0755091, matching the DOM meter's rounded $0.0755 display.

Keep this attributed, dated observation distinct from the following October 11 historical sessions. Do not turn the retained DOM evidence into a claim of independently viewed screenshots or heard speech.

No other concrete revision is required from this pass. The marker-after-adjacency explanation, optional-zero versus invalid/required-unknown distinction, Sonnet 4.6 manual route, and inspect-before-restart lens instructions fit the chapter's existing ownership and teaching sequence. The section-body and nested-total contracts are already explicit. The opening, historical Sonnet 5 stories/prices, and provider/model qualifications need no change for this finding.

## Sources and limits

Read the complete current text-reviewer and author skills; `AGENTS.md`; `docs/agentic-codebook-workflow.md`; `docs/edition-2-carryover.md`; `docs/edition-3-execution.md`; `book/voice.md`; `book/chapter-writing-procedure.md`; the full current chapter; and `docs/edition-3-notes/ch15-restart-{student,code-review,author}.md`, including the final code-review/root acceptance. Recovered clipped document reads. Read the exact author diff against `/tmp/ch15-author.9Faef0/chapter-15.md` (13 additions, 3 deletions).

Read retained local `receipt.log`, usage counts/totals/cost and wire/tooltip comparison fields in `usage-1.json` through `usage-4.json`, and `dom-turn-1.json` through `dom-turn-4.json` plus `dom-restart.json`, under `/private/var/folders/p_/11vk6y114xq7k9g5r_pc3yv80000gn/T/ensemble-ch15-live-loopback-uvfoe9fz/`. These support the rounded display, four recorded answers, normal EOF, replay and zero process meter. Screenshots were not viewed. Current lens files reset at restart; the retained per-turn wire/tooltip fields supply the comparison evidence.

Command/ownership consistency was checked against the accepted notes, not by executing the manual recipe or reopening source. No new external research, product tests/builds, operations, provider/browser/key calls, source inspection, or manuscript/code/guide/grader edits occurred. Only this note was written; available disk space exceeded 1 GiB. Existing root lint and accepted-code results were not rerun or represented as this reviewer's executions.

**Final handoff:** at most one author response to this finding, with root verification directly afterward. No second formal text review is requested.

## Root condition verification

The same author made the sole permitted response: dated October 22 and distinguished calculated $0.0755091, header $0.0755 and tooltip $0.075509. Root checked the final sentence against retained usage/DOM evidence; the restart, zero-process-meter and no-guarantee limits remain. The required correction is satisfied. This is direct root verification, not a second formal text PASS. No product checks or operations were repeated. An initial read-only verifier incorrectly treated the text lifecycle log as a JSON object; it failed before writing this note and was corrected without any rerun.
