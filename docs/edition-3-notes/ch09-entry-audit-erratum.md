# Chapter 9 entry audit: factual erratum

This is a narrow reconciliation of existing source and receipts, **not a second audit or review result**. The original `ch09-entry-audit.md`, its sealed inventory, and all original outputs remain immutable. This erratum supersedes the specific statements below, including their repetition in the original submission. No controls, tests, builds, applications, providers or browser sessions were rerun; no children were launched.

## 1. Prerequisite counts and chapter citations

The original report's line 20 says Chapter 9 has 448 lines and describes carryover reading through 1216. The existing `ch09-entry-audit-evidence/prerequisite-reads.json` records the correct pinned file lengths:

| File | Correct lines | Bytes |
|---|---:|---:|
| `book/edition-3/chapter-09.md` | 454 | 19058 |
| `docs/edition-2-carryover.md` | 1110 | 81871 |

A requested read range ending at 1216 can exceed the carryover file's EOF; it is not a count of actual lines. The ledger describes recovery through 1110 and an overlapping request through 1216. This corrects the report's counts, not the ledger's existing record that prerequisites were fully read before controls. No new prerequisite-read certification or audit is asserted here.

For the public-contract paragraph at original report line 108, use these current pinned Chapter 9 ranges: TL;DR/steps 25-177; format 179-210; registry 212-265; progressive disclosure 267-281; variables 283-320; provenance 322-346; constitution 348-359; exercise and historical evidence discussion 361-454. Thus the combined final range is **348-454**, not 348-448. The focused pending-unload discussion is at 437-442. The earlier section ranges remain valid.

## 2. Actual sensitivity observations

The original report's line 100 incorrectly attributes a failure to `skills_test.go:122` with the text `pending unload revoked dynamic capability`. That location and message are **not in this control's saved output** and are withdrawn as an observation of this run.

The existing `sensitivity.stdout` instead reports exactly:

```text
    skills_test.go:133: tool not enabled: write_file
    skills_test.go:159: dynamic provenance: {Source: EventSeq:0}, expected 7
```

The saved command receipt records `go test -work -v -count=1 -run ^TestSkillsThroughLoop$ ./ensemble` in `/tmp/ch09-entry-audit/mutant`, exit 1. Its output shows `TestSkillsThroughLoop` ran and failed; it is not a compile failure or a JSON result.

**Source-grounded interpretation, separate from observed output:** `toolNames` reads `a.Declarations()` (test lines 48-54). The test requests `write_file` before `unload_skill` (86-88), then checks enabled declarations (131-134) and first-enabling provenance (152-160). The saved mutation deletes dynamic tool entries in `MarkUnload`. The two reported failures are consistent with premature loss of `write_file`'s declaration and provenance after unload. They do **not** report a failed post-unload write dispatch, nor the nonexistent line-122 assertion. The current chapter itself distinguishes the focused declaration check from the revised grader's post-unload write (437-442).

The previously saved baseline and restored outputs each report the focused test passing. Those are existing results only: no second negative, new restoration, or new PASS was produced for this erratum.

## 3. Entry recommendation

**Unchanged, with the narrower factual description above.** Existing receipts still contain Edition 3's Chapter 9 rejection, the complete historical API's 100/100 result, and the one compiling negative with baseline/restored focused results. The corrected negative remains evidence of sensitivity to premature pending-unload removal, specifically through declarations and provenance. It does not add dispatch coverage to that focused control.

The measured recommendation remains a thin Edition 3 Chapter 9 dispatch/validation adapter using existing checks and separate executable/core identities. It is not an implementation, chapter acceptance, architecture-completeness finding, expanded score contract or authorization for a new campaign. All original limitations remain in force.

## Evidence and boundary

`ch09-entry-audit-erratum-evidence/excerpts.md` contains small verbatim excerpts with locations. `pins.json` records relevant input hashes, distinguishes prior sealed pins from current hash reconciliation, and identifies the original report/inventory. `inventory.json` seals only this erratum and its small evidence files, excluding itself. One parent-authorized read/hash-only packaging script calculates and verifies these receipts; it runs no behavior controls. Original audit completion records, including final191/join2231, remain prior completed evidence rather than new results. No source, grader, guidance, student or manuscript files were edited.
