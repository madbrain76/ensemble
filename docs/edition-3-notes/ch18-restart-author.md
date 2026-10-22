# Chapter 18 restart — minimal post-student author pass

2026-10-22. Read the author skill and mandatory current guidance, full chapter, and restart student/code-review notes through the final R4 disposition; recovered clipped reads. No files exist under `.agents/chapters`. Current lean restart governs historical campaign instructions.

Only book changes: a short deferred-definition/cache-marker explanation beside progressive disclosure, and a dated CLI observation. Deferred definitions remain in the sent catalogue; paired inline `tool_reference` content is not a cache marker. Anthropic forbids `cache_control` on a deferred definition: choose the last non-deferred definition or omit the tool breakpoint, preserving catalogue/disclosure and system/history markers. Root had already verified the linked public documentation; no fresh fetch. The HTTP-400 teaching mention is deliberately brief.

Preserved Bill's preface, historical observations, useful rationale, and existing primary/three prompts. The earlier reconstructed-primary failure was root setup, not a book defect. No new requirements or private-fixture answers.

Evidence read: existing restart records and only the bounded `verification.txt` and `receipt.json` in `/private/tmp/ch18-live-r4.FWfiti`. API-key Sonnet 4.6 → Opus 4.6; keyless Firecrawl one search/one scrape; eight model POSTs (6+2), all 200; normal EOF/exit 0. Five CLI snapshots and two model status results matched actual response prefixes. Usage: 14 fresh / 39,375 cache-write / 36,687 cache-read / 1,895 output tokens. The $0.23440875 model-tariff estimate excludes Firecrawl and earlier attempts; it is not invoice proof or universal exact billing. The release filename/citation is supported; disclosed truncation does not establish all omitted regions. No Codex/OAuth or optional GUI/browser/native/install/authenticated-MCP success claim. Earlier failures remain historical.

## Prose lint only

Canonical area: `/private/tmp/ch18-restart-author.tHedLq`; same-basename copies `baseline/chapter-18.md` and `final/chapter-18.md`. From `/Users/bill/projects/ensemble`, once per stage, sequentially:

```sh
work=/private/tmp/ch18-restart-author.tHedLq
# stage=baseline, then stage=final, in separate invocations
env -i PATH=/usr/local/go/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin \
  HOME="$work/home" TMPDIR="$work/tmp" GOMAXPROCS=2 \
  GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/Users/bill/go/pkg/mod GOCACHE=/Users/bill/Library/Caches/go-build \
  go run ./cmd/lintprose "$work/$stage/chapter-18.md" \
  > "$work/$stage-lint.txt" 2>&1
```

Both exits 0, saved in `baseline-lint.exit` and `final-lint.exit`; no lint setup/execution failure. Both >=1 GiB disk-floor checks passed (21,229,656 / 21,213,264 KiB available). Same two WARN categories: negation forms 29→32; person gap 1,764→1,873 words (anchor L906→L922). Added qualifications account for the increases; all other checks remain OK. Inspected the exact addition-only `diff -u` against baseline (expected difference exit 1). Current book matches the linted final copy. A pre-dispatch interruption prevented the first note write; state was checked and completed lint was not rerun.

No children, network, credentials/settings/raw live bodies, other gates, application builds, installs, or commits. Production, snapshot, graders, guides, other books/papers and unrelated work untouched. Root owns one independent text review and at most one bounded author response. Steady HOLD.
