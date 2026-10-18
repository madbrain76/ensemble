# Edition 3 Chapter 9 — one-pass independent text review

## Verdict and reviewed identity

**Required revision.** The original author draft is not ready for text acceptance: its advertised CLI build target does not exist in the accepted student tree, and its headless example invokes an unsupported subcommand. A third bounded accuracy correction must distinguish an illustrative Go constructor from the accepted public API. These are manuscript/example defects, not grounds to reopen accepted code or operational gates.

This is the **one formal text-review pass**, conducted by `e3-ch09-text-review`. Root owns direct verification of at most one author revision; no second review round, implementation campaign, or operational replay is requested. No manuscript or original author output was edited.

Authority: `/tmp/cadence-ch09-text-review-task.txt`, 5,960 bytes, SHA-256 `de6cc3dbfcf419da265660bbccf8aa23b2753a7aea6d721f3d1af394ebd3df50`. The coordinator's private input is `docs/edition-3-notes/ch09-text-review-input.json`. Its 28 `inputs` pins matched before review and are checked again before sealing. Its `protected_tracked` inventory was not used to open or hash publication/manuscript artifacts.

- Baseline: `/tmp/cadence-ch09-text-baseline/chapter-09.md`, 19,058 bytes, SHA-256 `1fc63f4a8b2680dbfd89658028c0a53d6a7c8dcd762c881e8f7b8ee0ba26e4ca`.
- Reviewed draft: `book/edition-3/chapter-09.md`, 20,767 bytes, SHA-256 `137f7fe93628d6434e5be6266933637b4bb0362b1926a6fe1e6da061db16447d`; byte-identical to `/tmp/cadence-ch09-text-review-initial/chapter-09.md`.
- Accepted student tag: `edition-3-ch09-r1`, independently resolved to `251ce88fb1677897dbcc6b33661096c5cb39cd8e`.
- Original author report, diff, lint receipts, read/process history and seal remain immutable. All 28 entries in the author seal's artifact map matched on initial verification.

Line references below refer to that unchanged reviewed draft and to line-numbered accepted-source receipts in `ch09-text-review-evidence/`.

## Finite required-revision list

### R1 — high: use a real executable package in the build recipe

**Location:** chapter lines 137–143, particularly line 142; original `ch09-author-report.md` line 23.

The new command is `go build -o /tmp/ensemble-ch09 ./cmd/ensemble`. At the accepted tag, the `cmd` tree contains only the `cmd/ch01` package (main/display source and their tests); there is no `cmd/ensemble` package. The root module is `module ensemble`, and the accepted README's CLI recipe builds `./cmd/ch01`. This independently confirms the coordinator's concern without executing a build. The output executable may still be named `/tmp/ensemble-ch09`; its name does not determine the source-package path.

**Required correction:** use `./cmd/ch01` for a concrete recipe against the accepted student layout, or explicitly introduce a reader-supplied executable-package path and give the accepted path as the example. Preserve the important distinction between the importable core module and its command package. Do not add or rename student code to make the erroneous prose true.

**Record reconciliation:** the original author report says the nonexistent target follows the public student recipe and was checked against records. The revision handoff must explicitly correct that claim in a new revision record, while preserving the original report and seal. Do not silently rewrite the original evidence.

**Evidence:** `student-tag.*`, `student-cmd-tree.*`, `student-module.*`, `student-readme.*`. The tree listing and README are from the immutable accepted tag, not the current working tree. No failed-build transcript is claimed or needed.

### R2 — high: remove the unsupported `prompt TEXT` invocation

**Location:** chapter lines 143–146, especially the shell command at line 143.

`EN_SKILLS_DIR=path/to/skills EN_PRIMARY_SKILL=base /tmp/ensemble-ch09 prompt "load the code-tools skill"` is presented as the public headless launch, with `chat` offered as an alternative. The accepted CLI accepts `chat`, `render LOG`, or `dump`; an unrecognized first argument follows the `unknown command` branch. No-argument execution selects the JSON-lines protocol. There is no `prompt` subcommand. This is an inherited example defect, not an author-introduced code regression, but it is part of the runnable instructions in the reviewed chapter.

**Required correction:** either use the already-supported human `chat` mode and instruct the reader to enter the quoted text there, or show a correctly labeled, valid JSON-lines request on standard input to the no-argument executable. If retaining a headless example, explain that its wire input is not human chat. Do not invent a new subcommand or require an implementation change. Preserve the goal: both clients reach the same Agent/registry/tool path.

**Evidence:** `student-cli.stdout.txt` lines 74–77, 97–141 show parsing and mode selection; its `protocol` function shows the `prompt` JSON field. The accepted README's Chapter 1 and Chapter 2 descriptions agree. These are read-only source observations, not a new provider or CLI run.

### R3 — medium: distinguish illustrative API spelling from accepted runnable Go

**Location:** chapter lines 315–334, especially the `app.NewAgent(cfg, ensemble.AgentSpec{...})` fragment at lines 322–328.

The surrounding prose introduces “the existing Ensemble root `app` and its provider configuration `cfg`,” followed by an unqualified Go example. The accepted root API instead exposes `NewAgent(cfg Config) (Agent, error)`. Its README sets `Config.Skills` to `&ensemble.SkillSpec{Directory: ..., Primary: ..., Variables: ...}`. The public root exports `SkillSpec` and no `AgentSpec`; the accepted CLI also calls `root.NewAgent(config)` with one argument. A reader cannot copy the displayed constructor into this accepted API.

**Required correction:** either label this fragment explicitly as illustrative API spelling, not a compilable call against the accepted student snapshot, and identify the accepted `Config.Skills` equivalent; or update only the fragment to that existing public construction shape, preserving its error handling. Keep skills opt-in and preserve nil-spec compatibility. This is not a request to standardize student identifier choices, add a new API, or reimplement the registry. The conceptual names are allowable; their status must be clear.

**Evidence:** `student-readme.stdout.txt` lines 186–190; `student-public-construction.stdout.txt` records the public aliases and `NewAgent` declaration at source line 106; `student-cli.stdout.txt` lines 71 and 93 show actual configuration and construction.

These three items are the complete requested revision scope. The lint warning and the observations below do not authorize a polish sweep or additional gates.

## What the draft gets right and should preserve

### Student experience and architecture

The added registry paragraph at lines 264–269 is a useful, evidence-backed clarification. It separates runtime ownership from Go package imports, states that Agent depends on the common lifecycle interface, and leaves construction/injection of the concrete registry at the root. The student review records the actual first-candidate sibling import failure and its repair. The paragraph does not pretend the initial attempt passed, invent an anecdote, or prescribe a new exported API. Preserve this explanation.

The chapter retains the successful teaching sequence: concrete frontmatter and discovery contract; primary/loadable/dependency distinctions; per-Agent state; progressive disclosure; declaration/dispatch enforcement; fixed initial system text; pending unload; validation-before-mutation and provenance. The human-facing observation that a skill cannot be unwritten motivates the pending-unload rule. The rendering section explains the callback/lock deadlock risk rather than merely ordering the reader to avoid a lock. These reasons remain load-bearing.

The bounded author diff changes three sites: executable-package selection, Edition 3 grading instructions, and the registry import-boundary explanation. It does not remove skill contracts, rewrite the scratch-file exercise, weaken provenance or unknown-variable behavior, or claim future compaction is implemented. It preserves the opener's concrete capability-discovery problem and the useful contrast between tool installation and capability publication.

### Grade command and media map

Apart from R1's invalid build target, the updated grading command agrees with the current `cmd/grade-edition3/main.go` CLI contract: Chapter 9; explicit core/consumer/workflow/stream-exercise/GUI module roots; nonblank stream and GUI model identifiers; flags before the single positional target; and a separate Chapter 1 gate on the same executable. `-consumer-case` is optional. The GUI module remains required by this grading gate even when the reader normally uses the headless client.

The text correctly distinguishes module directories from the executable and does not pretend placeholder model values are discovered live IDs. `ReadMediaModels` at lines 21–42 confirms a JSON string map with exactly the `anthropic`, `openai` and `gemini` keys and nonblank values. This validator does not establish whether each identifier is declared by the student; the chapter separately instructs the reader to select declared identifiers. The chapter's explicit map avoids a historical fallback; it does not falsely make `-media-models` universally mandatory. No grader was run and no deeper grader behavior was inferred from this CLI inspection.

### Demonstration and evidence limits

The retained exercise tells the reader to load a skill, observe the actual read/write actions and file changes, inspect the resulting render, mark unload pending, and try a hidden dispatch. It does not fabricate a transcript. Its final section labels linked evidence as historical, distinguishes deterministic rejection controls from real human-mode browser/model runs, and leaves later compaction application outside the demonstrated claim. Current accepted records preserve qualifications such as deterministic-only variable/constructor cases, no live schema/transport/rate-limit campaign, and the distinction between pending unload and later removal. Nothing in this review upgrades those qualifications to new demonstrations.

All five local Markdown link targets exist. I checked existence only; I did not open or hash their historical target contents or independently certify their older receipts. The pinned current student/code/operational records remain the evidence for current accepted experience. No new live-success claim is made.

## Fresh lint and human judgment

The current `cmd/lintprose/main.go` was read and pinned. Fresh baseline/current lint ran sequentially from `/Users/bill/projects/ensemble` under an explicit credential-free environment with `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOWORK=off`, `GOENV=off`, and scratch-only Go caches. Exact argv, environment, times, exits, free-space readings and output hashes are in `lint-baseline.json` and `lint-current.json`; stdout and stderr are separate unedited files.

| Metric | Baseline | Reviewed draft |
| --- | --- | --- |
| Exit status | 0 | 0 |
| Person gap | 1133 / 1200 | **WARN 1324 / 1200** |
| Negation forms | 2 / 4 | 4 / 4 |
| Short closers | 13 / 16 of 50 paragraphs | 12 / 17 of 52 paragraphs |
| Other printed metrics | within thresholds | within thresholds |
| Stderr | empty | empty |

Exit zero is not a warning-free or editorial-acceptance claim. The linter's person-gap warning is soft: the expanded command contract delays the next first-person anchor. The explanation is useful, the later first-person discussion still earns its place, and inserting a gratuitous personal aside would be worse. No separate revision is required for that heuristic. The instruction-heavy TL;DR earns its length through concrete contracts, while the body supplies motivation and causal explanation. No length target, generic shortening, or optional stylistic clean-up is requested.

## Read coverage, process qualifications and preservation

The complete current governing documents, role/coding skills, chapter, student/review/decision records, coordinator input, and required author diff/seal/read/process/lint evidence were read. Long documents were read in bounded ranges; truncated combined outputs were recovered with bounded or standalone reads. `read-coverage.json` records paths, identities and qualifications. The byte-identical Edition 3 voice copy was also included in the reads and compared to the fully read canonical voice. Baseline/current texts and exact author diff were read, rather than inferring the chapter solely from a patch.

**Scope deviation retained:** one later permitted-source reread mistakenly requested `sed -n '1,180p' internal/gradeedition2/ch05_media.go` rather than stopping at line 43. It displayed lines 44–180 of adjacent media-helper implementation in addition to the allowed validator. I stopped there, did not inspect any further grader implementation, did not execute it, and did not use that overread to add findings or claims about gate behavior. This report therefore does not claim perfect read-scope compliance. Root must retain this qualification when deciding how to use the review. Hashing the whole pinned file for preservation did not itself justify that extra content read.

The accepted student source inspection was limited to the README, module declaration, tree listings, CLI main file and exported construction API. Source was fetched at the accepted tag, not built or run. The public construction receipt explicitly distinguishes its full-source hash from its displayed declaration snippets. Initial/final pins cover the actual named inputs; no recursive protected-publication audit was used.

No credentials were opened. No provider/network/browser calls, tests, race/vet campaigns, graders, student builds, servers, children, git mutations, broad status/diff scan, publication outputs, or chapter edits were performed. Command batches were sequential and all completed. The only Go execution was the two prose-linter runs. Free space stayed above 1 GiB. Original author outputs and both immutable manuscript copies remain unchanged. Unopened protected material, including Bill's preface, is covered by a scoped-write assertion only; root owns the independent protected-file audit.

## Handoff and seal

Root may request at most one bounded author revision resolving R1–R3 and the R1 record correction, then verify it directly. Preserve this reviewed draft, this report, all original author receipts and their seal. Do not relabel them as evidence of the revision. Material unresolved issues go to root, not another reviewer loop.

`ch09-text-review-evidence/seal.json` is written last and pins this exact report plus every evidence file in the review directory except the seal itself. It records reviewed identities, input preservation results and the non-self-hashing rule. The submitted seal byte count and SHA-256 identify the exact final seal. This report makes no text-acceptance or publication claim.
