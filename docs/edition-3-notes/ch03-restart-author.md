# Chapter 3 restart — author

Bill explicitly requested more useful `edit_file` refusals in the manuscript and current student implementation before Chapter 4. This is an authorized timing exception before final code acceptance, not an author-invented feature.

Added one paragraph in Chapter 3 §3.5 after the overlapping-anchor example: the student's strict exact-unique policy reports zero matches and suggests rereading; multiple matches report the actual count, including overlapping sites, and a bounded list of one-based starting lines, marking truncation without truncating the count. Refusals leave bytes unchanged and neither choose a match nor relax exact comparison. Existing valid empty `new_text` and other semantics are untouched. No excerpts, fixed output format or new read mechanism are required. This describes the student implementation, not a change to the running tool/runtime.

Checked the chapter diff: one added paragraph only. Prose lint was not run because the repository invokes it through `go run` and this assignment prohibits builds.

Only the chapter and this note are author-owned changes. No code, graders, skills, builds, browser/provider activity or commits. No formal text review now: retain this author context for any minimal post-student clarification grounded in actual experience, then one fresh text review at chapter end with at most one author revision.

## Bounded post-student pass

After acceptance of the 27-file `solutions/edition-3/ch03` snapshot, clarified three demonstrated lessons: Engine projects common schemas into provider dialects (omit only Gemini's unsupported top-level `additionalProperties` for current built-ins, preserving strict local validation/shared schemas/other providers); exact signature replay needs `LLM_RESOLVED_MODEL` from an observed returned identity on the same vendor/surface/route (both launch authorization and harness omitted it; no source repair); exercise fixtures/logs/HOME/temp stay outside the tool workspace, and recall uses the actual third answer plus inspected effects. Preserved the strict-edit diagnostics paragraph and existing registry/tool contracts.

One same-basename baseline/final prose-lint invocation: both exit 0; existing warning categories remain. Negation count 18→21; person gap 2342→2449 words. No runtime/provider campaign, code edits, children or commits. Held for the one fresh chapter-end text review; at most one author revision remains.

## One author response to text review

Addressed the sole required revision in §3.5: “Both refusals identify the target file and leave its bytes unchanged,” with no fixed message format or other prose changes. File identity was the round-2 code repair; round 3 addressed Gemini schema projection. Final prose lint ran once and exited 0 with existing soft warnings (negation forms 21; person gap 2449 words), not zero warnings. Final held for root condition verification; no formal re-review or further author revision.
