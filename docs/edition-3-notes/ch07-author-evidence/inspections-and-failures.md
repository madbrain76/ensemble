# Author inspection corrections (not runtime results)

This record preserves unsuccessful local author steps rather than silently presenting a straight-line process. No such failure is a student or grader outcome.

- Some combined prerequisite reads exceeded inline output caps. The omitted portions were recovered by separate complete reads or bounded overlapping ranges before writing. The workflow and carryover range byte sizes were subsequently checked against their explicit caps.
- Several guessed supplemental paths did not exist: `docs/edition-3-problem-sets/ch07-artifact-stream.md`, `docs/edition-3-notes/ch07-live-exercise-evidence/parent-toolcall-link-crosscheck.json`, and `docs/edition-3-notes/ch07-native-exercise-evidence/attempt-4-child-parent-link-check.json`. No content from those paths was assumed. Actual required records and public CLI source support the patch; no chapter link to a guessed file was added.
- `rg` was unavailable (`/bin/sh: rg: command not found`). Inspection continued with standard `grep`, `find`, and bounded `sed`/`cat` reads.
- First editing attempt stopped with `KeyError: 'protected'`: the baseline uses `tracked_protection`. No chapter write occurred. Its trailing Git check unexpectedly opened a `less` pager. The process tree was inspected and only the confirmed pager child (PID 49087) was terminated. Later Git commands used `git --no-pager` explicitly.
- Second editing attempt verified the real protection field but refused an inexact prose anchor (`That tuple is the GUI’s public contract`). No chapter write occurred. The actual full chapter was recovered as overlapping line ranges, and exact exercise/observation/speech anchors were used next.
- The successful edit ran fail-fast, verified the supplied chapter baseline and all 30 protected hashes, asserted exactly one match per replacement, and wrote only the chapter. Its diff and check output are retained as `chapter.patch` and `checks.json`.
- A later exact edit to label a model default historical refused with zero matches. Reading the actual nearby passage showed that it already says “historical default”; no change was needed or made. The sealed chapter pair therefore still reflects only the three intended hunks.

No commit, grader invocation, build, paid call, browser, native application, or publication was used to recover these author inspection errors. The example shell commands were syntax-checked only, without running their actions.
