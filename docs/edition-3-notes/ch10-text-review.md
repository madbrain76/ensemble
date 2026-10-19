# Edition 3 Chapter 10 — independent text review

**Verdict: REQUIRED_REVISION.** One formal whole-chapter review; two blocking findings below. No manuscript edits were made. Root may assign the sole author revision; this reviewer is not opening a second review loop.

## Reviewed candidate and scope

- Chapter: `book/edition-3/chapter-10.md`, **27,916 bytes**, SHA256 **`05e433823f0757ecb65db9b69d59d9814b1c734e056c3419a8bf156e30a562f9`**.
- Immutable input: `docs/edition-3-notes/ch10-text-review-input-chapter.md`, byte-identical to the chapter; input manifest `ch10-text-review-input.json` verified.
- Baseline: `ch10-author-evidence/chapter-before.md`, 21,655 bytes, SHA256 `abda8ed3a74c2ee2c64bc2208ea104e0547810d6e06e1d5e772c0bb4ee84b3ed`. Full baseline, current chapter and delivered diff reviewed.
- Accepted implementation: student commit `729bc919e734405307b4ac0060e65c40e33652e5`, tree `025ef3744b42c8f9daec34e76661edc40583ca6a`, clean at inspection; acceptance records identify tag `edition-3-ch10-r1`. Code and required operations remain accepted, with zero formal production corrections.

Root disclosed two inspection concerns—grade chronology and body first-person narration—as assistance, not a supplied verdict. Both findings below were independently established from the permitted sources. The whole chapter and changed commands were assessed, not merely those passages.

## Blocking findings

### R1 — Correct the cause of the original frozen 80/100 becoming 100/100

**Location:** chapter lines **508–514**, especially “The default-load probe needed a valid minimal running configuration, and inherited skill checks needed to honor the implementation's declared models” followed by “With those corrections, Chapter 10 passed 100/100.”

That connective account misidentifies the correction responsible for the retained first frozen failure. `ch10-student-final-reconciliation.md:24–28` explicitly records that the default-load fixture supplied a prompt with a **terminal newline**, while its unchanged exact expectation did not; the fixture input was corrected to match the expectation. It also explicitly says no student production defect was established, no trimming was requested, and no source correction was authorized or made.

The original receipt corroborates the distinction: `ch10-student-evidence/grader-ch10-final-v1.json` records exit 1 on the frozen candidate, and its `.stdout.txt` reports **80/100, passed:false**, with only the 20-point `default-load` check failing for core and GUI. The inherited `ch9-parity` check already earned **15/15 and passed**. Later independent fixture problems must not be presented as the cause of this original failure. The final reconciliation separately records Chapter 10 **100/100, passed:true, exit 0** and same-binary Chapter 1 **100/100, passed:true, exit 0**, including the real lifetime phases.

**Required correction:** replace the causal explanation with the supported newline/exact-expectation fixture mismatch and its fixture-only correction. Preserve the original failure, unchanged production source and subsequent successful gates. The minimal-configuration/model-fixture issues can be omitted here; if retained, distinguish their separate chronology rather than making them explain this 80-to-100 transition. Preserve the runtime-configuration ownership lesson, but do not imply that a student runtime-configuration defect or trimming fix was proved. A short accurate explanation suffices; no report diary or new investigation is requested.

### R2 — Restore third-person narration in the chapter body

**Locations:** chapter **189–190**, “I put it this early in the book because…”, and **360–361**, “I want the old card to keep its address…”. Both are introduced by the author diff and occur in §§10.1 and 10.5, after the opener and TL;DR.

`book/voice.md` §§1.1–1.3 requires third-person chapter bodies, permits Bill's first person only in the short motivational opener **before the TL;DR**, and treats Bill as a third-person character outside it. These are unquoted authorial statements, not reader instructions or sourced historical dialogue. The first also invents an unattributed authorial reason for the book's placement of persistence.

**Required correction:** recast both sentences in third person without deleting their engineering reasons. Preserve why persistence makes the replay claim testable and why old card addresses and newly allocated addresses must remain distinct without prescribing a new public identity API. The baseline's placement rationale already avoids first person. Do not extend the opener's exception into these sections or use first-person insertions to improve a lint metric.

## Whole-text and command assessment

No further blocking finding was established in this bounded review.

- The Edition 3 grading recipe agrees with `docs/edition-3-student.md:158–175` and public `cmd/grade-edition3/main.go`: five separate module inputs, both declared models, optional consumer-case/media-map inputs, all flags before the explicit executable. It immediately pairs Chapter 10 with Chapter 1 on the same already-built binary, requires parsed pass plus exit zero, and preserves the real 62-second header/body waits and invalidity of a 90-second outer deadline.
- The build example correctly distinguishes the accepted library root from nested `cmd/ch01` and retains alternate root-main layouts. The accepted student's complete public `README.md`, `CHAPTER10.md`, `cmd/ch01/main.go` and `main_test.go` support `LLM_VENDOR`/`LLM_MODEL`, `--save` before `chat`, and the separate live JSONL prompt envelope when `chat` is omitted. No invented protocol subcommand or new API was found. The examples were not executed.
- The opener, TL;DR and core explanations remain: Agent ownership; recorded versus effective/running configuration; intentionally empty primary instructions; rich provider-neutral replay; full-log, compacted-tail and context-only cases; sequence anchors rather than lengths; atomic refusal and preserving malformed bytes; EOF admission/drain/save ordering; and the legacy GUI-ID collision rationale. The new address discussion preserves provider payload/correlation identities and explicitly avoids a new identity API. No new capability/schema or coverage-framework requirement was found.
- Historical three-provider receipts remain explicitly historical, not a fresh Chapter 10 provider campaign. The accepted operations report and parent receipt support the bounded Claude tool/EOF/relaunch result (three HTTP 200 responses) after the preserved transport failure and explicit TLS recovery. The chapter leaves the original failure's cause unproved and distinguishes the fake-provider GUI DOM checks from human visual/audio observation. Screenshots were hashed, not viewed. This review adds no provider/browser/native-speech observation.
- The author report's delivered basename and JSON companion are preserved as delivered; neither is a manuscript defect.

## Checks, reads and preservation

All current prerequisites were read in full, including the complete voice guide, chapter-writing procedure, workflow and Edition 2 carryover. Long/clipped displays were recovered with bounded full or consecutive reads; no compaction was used. `ch10-text-review-evidence/required-reads.json` records actual full-read paths, byte sizes and current SHA256s. Public entrypoint inspection did not open private grader/reviewer probes, raw captures, credentials, historical final solutions or future implementations.

The author's actual same-basename final lint is **carried, not rerun**: its sealed receipt records exit 0/no hard failures and the disclosed 1,297-word person-gap warning. The pinned static receipt records clean diff checking, four successful shell-syntax checks, balanced fences and five existing relative links. Chapter and receipt pins match. Lint does not validate narrative truth or authorize the two first-person statements; the soft warning alone is not a revision demand.

All **37 supplied input-pin records** and **24 author-seal members** verified. The preservation evidence covers the chapter/input copy, governing reads, delivered author files, permitted acceptance reports and public interfaces; it is not a private-evidence or whole-filesystem audit. Only this report and `ch10-text-review-evidence/` were created. No code tests/builds, fresh lint, graders, providers, browser, paid calls, children, edits to prior evidence, commits, tags, cleanup or publication were performed. A failed unavailable `rg` search read no source and was replaced with public-only standard-tool reads. All read/analysis/write jobs finished before submission; disk remained above 1 GiB.

`ch10-text-review-evidence/seal.json` is the exact UTF-8/no-NUL inventory of this report and its evidence, excluding the seal itself. The sealed verdict requests only the two bounded prose corrections above.
