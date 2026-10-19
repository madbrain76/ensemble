# Chapter 2 post-student author note

Added one seven-line paragraph in §2.5 Usage, after the normalization rule. Its
worked OpenAI example has prompt total 100, invalid cached count `"bad"`, cache
write 0, and output 7: `100 − unknown − 0` leaves fresh input unknown, while the
independent output charge remains measured. It distinguishes an omitted optional
bucket's zero default from an invalid supplied value, without prescribing an
unknown-count representation or adding an API/schema. This addresses the repaired
usage bug recorded in `ch02-restart-student.md` and
`ch02-restart-code-review-r1.md`; the accepted snapshot's
`internal/engine/parse.go` (`parseUsage`) confirms the dependency rule.

Kept the existing nine-line TLDR clarification unchanged. The earlier rejected
attempt ignored an explicit public schema; missing teaching did not cause every
implementation bug. Existing §2.5 already distinguishes ordinary signed Gemini
text from thinking and keeps visible content portable, so it needed no addition.
Exact offline replay targets and legitimate Anthropic string content were repaired
grader assumptions, not student product failures. All pre-existing chapter text,
including the historical opening, remains unchanged.

One baseline/final prose-lint comparison used identically named `chapter-02.md`
files under `/tmp/ch02-author-prose.QpfHRC/{baseline,final}`. Both runs exited 0.
Soft warnings: negation forms 23 → 25; the existing person gap after line 634 grows
from 5,540 to 5,605 words. Other checks pass. Chapter `git diff --check` passed.
Only the chapter and this note were edited in this phase. No code edits, tests,
providers, grader changes, children, or commits; formal text review remains pending.
