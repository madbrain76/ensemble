# Chapter 19 restart author: HOLD

Date: 2026-10-22 (host local date).

Stable **HOLD** after the sole bounded response to the formal text review, for root's direct verification. No second formal review is requested. This note does not close Chapter 19 or authorize a commit.

## Scope and edits

Read the author skill and complete mandatory guidance: `AGENTS.md`, workflow, carryover, voice, chapter-writing procedure, current execution guide, full current Chapter 19, and current restart student/code-review records. Recovered clipped portions. The current lean restart controls over historical campaign instructions. Read the Chapter 18 author note for the established offline lint invocation.

Changed only `book/edition-3/chapter-19.md` and this new note; temporary work is under `/private/tmp/ch19-author.X0Ldgq`. Existing unrelated work was preserved.

- Section 19.2 distinguishes trusted application/corpus access, rooted in-process file tools, and OS-confined commands/descendants with deliberate system-runtime exceptions. `/` remains the host root; this is not chroot.
- Section 19.4 adds the concrete Apple `/var` developer-selection metadata lesson. The three added literals were a sufficient tested combination, not a necessity proof. No broad subtree/global metadata allowance or warning cleanup is proposed. Existing root-literal and ancestor rules remain unchanged.
- Section 19.11 adds a short October 22 Sonnet 4.6 API-key CLI observation: useful work beside actual refusals, synthetic-token and controlled-network checks, and narrower safe-mode behavior. The old Astra account remains historical; its no-further-paid-session statement now explicitly refers to those second-edition corrections.

First HOLD baseline/final diff inspected: 236 net whitespace-delimited words, three added paragraphs and the historical scoping phrase. Byte comparisons confirm Bill's preface, the TL;DR public contract, and the complete Checks section are unchanged. The saved first-held final copy matched the working chapter at that HOLD.

## Evidence limits retained

The new observation rests on the current final code-review note and execution guide, which record root's paired raw-result/artifact checks and the reviewer's bounded extractions. This author read the small controller receipts, not raw live bodies, and ran no application verification.

Both controller receipts remain FAILED: the first reached its 12-request cap before safe mode received an upstream request; the separately bounded five-request follow-up completed descendant and safe-mode sessions but asserted a newline where the model supplied literal backslash-n. Actual write/read/artifact bytes matched. No prose or product change adopts that arbitrary assertion. The completed confined, descendant and safe sessions exited 0 at EOF. Shell-selection warnings did not establish command failure. No claim covers universal security, real private files, optional GUI/browser/OAuth, installation or restart. The deferred subscription exercise is distinct from the API-key observations.

## Prose checks

Saved same-basename copies before and after editing:

- `/private/tmp/ch19-author.X0Ldgq/baseline/chapter-19.md`
- `/private/tmp/ch19-author.X0Ldgq/final/chapter-19.md`

Actual invocations, sequentially from the repository root:

```sh
/bin/sh /private/tmp/ch19-author.X0Ldgq/run-lint.sh baseline
/bin/sh /private/tmp/ch19-author.X0Ldgq/run-lint.sh final
```

The saved wrapper checks for at least 1,048,576 KiB free, then runs this command with `phase` set to `baseline` or `final`:

```sh
work=/private/tmp/ch19-author.X0Ldgq
/usr/bin/env -i \
  HOME="$work/home" TMPDIR="$work/tmp" \
  PATH=/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin \
  GOCACHE=/Users/bill/Library/Caches/go-build \
  GOMODCACHE=/Users/bill/go/pkg/mod \
  GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local \
  GOPROXY=off GOSUMDB=off GOWORK=off CGO_ENABLED=0 \
  /usr/local/go/bin/go run ./cmd/lintprose -v "$work/$phase/chapter-19.md"
```

Both lint commands exited **0**, with all hard checks passing. Soft warnings remain: negation forms **22 → 23** against budget 11; person gap after L22 **3593 → 3718** against 1200. Short closers remain within budget, **28/112 → 29/115** paragraphs. These limited additions do not justify rewriting the existing teaching voice to chase density warnings. Full output and exits are in `baseline-lint.log`, `final-lint.log`, `baseline-lint.exit`, and `final-lint.exit` under the temporary directory. Disk checks recorded **17,588,152 / 17,584,704 KiB** available, both above the floor.

No lint setup failure occurred. Two preparatory inspection issues were resolved: the older Chapter 18 temporary wrapper had expired (the note still supplied its invocation), and the initial Git diff opened an interactive pager; that owned inspection was deliberately stopped (exit 143) and subsequent diffs used `--no-pager` or `diff -u`. No valid lint check was repeated for paperwork. No children, network/provider requests, credential/settings access, application build/test/native/browser/install/restart, snapshot edits, guide edits or commit were performed.


## Sole bounded formal-review response (R1)

Read the complete `ch19-restart-text-review.md` REQUIRED_REVISION and only the permitted source excerpts: `internal/common/auth.go`, `internal/auth/auth.go` and the specified HTTP-boundary portions of `internal/engine/engine.go` under `solutions/edition-3/main`.

Corrected only the two existing §19.5 passages and the minor §19.4 cross-reference. Plain auth-owner/Engine descriptions replace nonexistent `CredentialProvider`, `Engine.sendRequest` and `requestFailure` references. The text now identifies OpenAI credential resolution, the Anthropic/OpenAI/Gemini header distinction, and generic safe HTTP failure or context cancellation instead of credential-bearing transport detail. Section 21.3 now correctly points to 19.3. No alias or requirement was added. Exact response diff inspected: three hunks, 24 net words. The three previously accepted paragraphs and historical-scoping phrase remain unchanged, as do Bill's preface, public contract and Checks section.

New stable response copy: `/private/tmp/ch19-author.X0Ldgq/r1-final/chapter-19.md`; byte-identical to the working chapter. The original baseline, first-held final and their receipts remain intact. Reused the baseline rather than rerunning it.

One new offline lint invocation:

```sh
/bin/sh /private/tmp/ch19-author.X0Ldgq/run-lint-r1.sh r1-final
```

The wrapper is the existing approved wrapper copied with only its accepted phase name changed to `r1-final`; isolation, existing caches, offline settings and the 1 GiB floor are unchanged. Exit **0**, all hard checks pass, no setup failure. Soft warnings: negation forms **23** (budget 11; unchanged from first HOLD), person gap **3742** (budget 1200; previously 3718). Short closers remain **29/115**, within budget. New receipts are `r1-final-lint.log`, `r1-final-lint.exit` and `r1-final-disk-kib.txt` under the same temporary directory; available disk was **17592464 KiB**, above the floor. No earlier receipt was overwritten and no additional operation or validation lane was opened.

The sole author response is complete. Stable HOLD for root's direct verification and closure decision; no second formal review, delegation, source/test changes, provider/network operation, settings access or commit.
