# Edition 3 readiness closure

Bill authorized closing the known preparation issues before execution on 2026-10-16. This is a bounded preparation gate, not a new author-led design pass. Baseline preparation commit: `2be995ebc12c6a7efb21cee3b9740a758489adda`. No student has started.

This checklist records the facts submitted for independent acceptance. The frozen review manifest and final verdict are `edition-3-notes/readiness-review-manifest.json` and `edition-3-notes/readiness-review.md`; the checklist alone is not a review verdict.

## Checklist

- **Integrated — crossover ownership:** the [audit](edition-3-notes/crossover-closure-audit.md) and [integration record](edition-3-notes/crossover-integration.md) distinguish verified requirements from withdrawn preliminary suggestions. Dispatch belongs in 3/6, restart/legacy identity in 10, provider policy in 16, child composition in 20, actual application file mode in 21, and memory/recall extensions in 22/23. Chapter 25 is historical-only; Chapter 26 explicitly extends accepted Chapter 24 code. No new early configuration loader or identity API was prescribed.
- **Validated — Chapter 1 harness:** `cmd/grade-edition3` requires an explicit target and reuses the inherited checks. Two real 62-second header/body gates reject the historical 60-second client and accept a temporary corrected control. Focused race tests also passed in the orchestrator's independent run. These are grader controls, not student acceptance. See [the report](edition-3-notes/ch01-readiness.md).
- **Verified — live Anthropic access:** authorized model discovery and one minimal generation both returned HTTP 200 using `claude-sonnet-5`; the reply matched `OK`. The key was resolved in memory from the previously authorized store, never printed or copied into repository evidence. Codex access was not used as evidence.
- **Packaged — five thin roles:** `.agents/skills/edition-3-{coder,code-reviewer,author,text-reviewer,orchestrator}/SKILL.md`. Explicit full-file reads require no runtime registration/restart. They share the current execution procedure, including three code-correction rounds after initial review and the separate single text-review/single author-revision loop.
- **Created — empty student repository:** `/Users/bill/projects/ensemble-edition-3-student`, separate from the copied final reference. Empty-tree baseline `69859da3328c3f710f8143203b99f934fed778a1`; no source files and no accepted student chapters. Exact accepted commits and immutable `edition-3-chNN-rN` tags will identify later predecessors.
- **Recorded — mechanical checks:** `edition-3-notes/readiness-checks.json` captures hashes, link checks, protected originals/reference integrity, the empty student tree and the review snapshot. Independent review and final child-save verification remain separately recorded acceptance steps, not consequences of a successful file write.

## Execution boundaries

The student reads [the public brief](edition-3-student.md), not historical answer-bearing repair records. Mandatory instructions may mention historical paths without authorizing following them into answers. Authors/reviewers load their full governing documents. The copied `solutions/edition-3` remains an immutable final reference, not a predecessor; instruction-based input separation is not an OS sandbox.

Align affected later graders just before their chapters, with independent review and sensitivity evidence, rather than rewriting every grader before Chapter 1. This scheduled gate is not a claim that alignment already passed. In particular, the known artifact-identity defect must earn acceptance at Chapter 10; its prose obligation is not a runtime fix. Chapter 1 still requires the actual student's live multi-turn exercise and independent source review beyond the finite automated checks.

Publications are separate: copied PDFs/EPUBs/assembled Markdown remain old artifacts, and eventual builds need isolated Edition 3 outputs and independent format verification. No publication build, installation, restart or push is part of readiness. Preserve unrelated tracked, untracked and staged work; any preparation commit must name only this accepted scope.
