# Chapter 2 restart — formal text review

## Verdict: PASS

No required revision. This is the single formal text review, not a reopening of accepted code or operations.

Read the current AGENTS, coding/shared-execution instructions, text-reviewer role and required author/voice/chapter-writing/workflow/carryover prerequisites in full; the complete current Chapter 2; its diff against HEAD; and `ch02-restart-{student,code-review-r1,author,tldr-note}.md`.

## Findings

- **TL;DR:** The nine added lines restate the already-explicit frozen event/content/JSONL and CLI contract (§2.4 and Exercise). Causal turn state, four-bucket cumulative usage and ephemeral arrival/delivery agree with the chapter's reducer and context teaching. The pseudocode exception preserves implementation freedom without making public wire names negotiable. The initial attempt's reduced transcript/protocol was a failure to implement the supplied contract, not an omission in the book.
- **Usage example:** The seven added lines correctly evaluate `100 − unknown − 0` as unknown fresh input, retain unknown cache-read usage, and preserve independently known output of 7. They distinguish an omitted optional bucket's zero default from an invalid supplied value and require incomplete accounting to be reported. They introduce no null, sentinel, flag or other unknown-count schema into the existing four-field public shape.
- **Whole-chapter consistency and voice:** The additions fit the existing explanation, preserve the History/Context/request boundary and shape-before-capability scope, and leave the hard-won prose intact. The diff is exactly 16 inserted lines with no deletions. Existing text-versus-thinking/opaque replay teaching is unchanged; ordinary signed text is not reclassified as thinking or grounds for discarding visible text. Exact producer/target identity remains distinct from routing aliases. Historical measurements, original grading exceptions and later student results remain separately attributed; this edit does not upgrade them into fresh evidence.
- **Handoff accuracy:** Code/operations acceptance followed the fresh retry and two narrow fixes. The exact-target replay configuration and Anthropic string-content corrections were grader-assumption fixes, not student defects. The author-recorded paired lint-0 result is carried evidence, not a check performed by this reviewer.

## Boundary

Only this review note was written. No code, runtime, build, test, lint, provider or browser checks were run; no children, private raw traces, other student contexts, source changes, packages/manifests/seals or commits were used or changed. No optional polish campaign or formal rereview is requested.
