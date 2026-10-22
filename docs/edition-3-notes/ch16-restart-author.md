# Chapter 16 restart author — 2026-10-22

**Status: stable HOLD for the single formal text review.** Initial author pass complete; no formal text-review revision used. At most one bounded revision remains. Only this note and `book/edition-3/chapter-16.md` were edited. Required role/guidance/chapter/evidence reads completed, including recovery of clipped output; current lean rules governed.

## Changes

- Explicitly attributed the historical subscription failure, summary observations, and paid-run narrative to the earlier student run. Preserved its failures, measurements, rationale, and shutdown limitation; preserved Bill's opening.
- Added one paragraph explaining the demonstrated final-summary suffix gap: final cards survived, while the live-delta path missed text. Existing public-summary catch-up makes the unseen suffix available without repeating its prefix. No claim that canonical summary content was lost or that speech was heard.
- Added a dated restart status separating local passing grades from the historical paid runs and Bill's explicit live-OAuth deferral.
- Left the existing auth-command, completed-item fallback, model-capability, tool-pairing, SSE, developer/cache, and other requirements intact: they already teach the accepted contract. No new requirement or design/style campaign.

## Evidence and limits

Evidence reused from `ch16-restart-student.md`, `ch16-restart-code-review.md`, and the root handoff; no production behavior was rerun by this author. Code was accepted after three production correction rounds plus separately classified test-only repairs; the 216-file snapshot was already saved. Chapter 16's third grade was 100/100, `passed:true`, exit 0. The earlier 15/100 and 85/100 grades remain failures. Chapter 1's first grade was 100/100, `passed:true`, exit 0, with actual 62+62-second checks and 511/65 usage. Root/GUI race+vet passed; other-module/JS evidence retains its stated no-test/legacy exclusions.

Bill said: "Ok. Please defer the live oauth test and continue". Live OAuth/account/provider compatibility remains undemonstrated, not waived into success. No fresh native GUI, browser, or listening observation is claimed. No OAuth login, private-credential access, network, install/restart, student adapter, production test, child, or commit was performed by this author. This runtime's Codex credential route is separate.

## Checks

Inspected `cmd/lintprose/main.go`; used its existing same-basename procedure. The supplied baseline initially matched the chapter byte-for-byte (`cmp`, exit 0).

- `GOPROXY=off GOSUMDB=off go run ./cmd/lintprose /private/tmp/ch16-author-baseline/chapter-16.md` — exit 0.
- `GOPROXY=off GOSUMDB=off go run ./cmd/lintprose book/edition-3/chapter-16.md` — final bytes, exit 0.
- Outputs: `/private/tmp/ch16-author.Ph5buv/baseline-lint.txt` and `final-lint.txt`.
- Hard failures: **0 → 0**. Soft warnings: negation forms **13 → 14** (allowance **11 → 12**); person gap **3,937 → 3,940** words (limit 1,200). These are reported separately, not treated as hard failures or a reason for unrelated rewriting.
- Baseline/final `diff -u` inspected (exit 1 denotes the intended changes). `git --no-pager diff --check -- book/edition-3/chapter-16.md` — exit 0. An initial exact-anchor edit failed before writing; the subsequent edit used the inspected actual wording.
- Disk remained above the 1 GiB floor (about 22 GiB available at final lint).

No acceptance, publication, or broader restart advance is asserted. Await the root's one formal text reviewer.
