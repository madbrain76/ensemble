# Chapter 11: KISS / load-bearing prose pilot

## Scope and result

Bill authorized one read-only pre-execution review, excluding his opening. Deletion or simplification must not introduce ambiguity. No invented requirements, new architecture, speculative student difficulties, word-reduction quota or wholesale rewrite.

Reviewer: `ch11-kiss-pilot`. Target: `book/edition-3/chapter-11.md`, 463 lines, 3,386 whitespace-delimited words. Before/after SHA256: `f1a313cde64be5d631abaaf96007121d0505e617169fe34528d59072e6c5d161`.

**Reviewer recommendation: zero worthwhile cuts or consolidations; keep the chapter as written.** Cadence agrees with the reasons below and made no manuscript edits. This is a pruning judgment, not implementation acceptance or a finding about the rest of the book.

## Tempting cuts rejected

- **Lines 102–107, MCP job lifetimes:** reducing this to “MCP calls use ordinary jobs” would erase the distinct behavior of kill, callback, interrupt and shutdown. Other mentions describe different mechanisms and do not replace the comparison.
- **Lines 126–145, slow automatic collection:** retaining only the completion rule would lose the actual 3.4-second failure, why a running notice or longer timeout is insufficient, and the separate interruption/stale-result guarantees. These explain the design rather than repeat it.
- **Lines 435–444, browser evidence limits:** deleting qualifications would conflate applied font state with disk persistence, and a real browser plus provider fake with a live-provider browser run. The qualifications prevent stronger claims than the evidence supports.

## Separate coordinator observation

The reviewer established no factual inconsistency and did not audit linked receipts or implementation. Cadence separately noticed line 124: “Ephemeral errors are logged but not fatal. A GUI snapshot failure should not abort a turn.” Its GUI example is potentially misleading under the Ephemeral Tools heading now that lines 7–9 and 212–221 make GUI observations ordinary on-demand tools. A future factual cleanup could use an automatic environment observation as the example instead. No change was made in this pruning pilot; do not count this as padding found by the reviewer.

## Coverage and preservation

Reviewer reports complete reads of AGENTS, the governing workflow, carryover, voice guide, chapter procedure, Edition 3 execution instructions and the target chapter, recovering truncated output. The reviewed execution-note hash was `da0d7386ad664abfedd110bdd241efa3436f2486cb9ddb1cc7d5526dd2d47b4c`; Cadence updated that note during this run to record Bill's new pilot authorization and no-new-ambiguity rule. The chapter itself stayed fixed throughout.

Cadence independently checked the unchanged chapter hash and protected files:

- `preface.md`: `915cf8c531169c5aacc4cb81bd743eeba8a181d20bc779d936ac5f7ab9998d96`
- `note-from-bill.md`: `669ccd62c4df8462aff74152328a2f10dc67d89c303636e96782c31148897948`

Reviewer submitted and was joined; its result was retrieved and idle state checked before closure. No student execution, graders, implementation review, publication build, commit, push or live application change. The pilot is not authorization for automatic whole-book rollout.
