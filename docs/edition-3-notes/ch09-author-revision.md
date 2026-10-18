# Chapter 9 — sole bounded author revision

## Status and correction of the record

This is the same author's **one permitted revision** after the sole formal text review returned `REQUIRED_REVISION` (root-verified final154/join3004). Only R1–R3 were addressed. No second formal review was requested or performed. Root owns direct verification and the final decision; this report does not declare text acceptance or publication.

Revision authority: `/tmp/cadence-ch09-author-revision-task.txt`, 4,346 bytes, SHA-256 `a52a445d743b51449cc070de177e973d56f92981ec2555e016c1ced1dd89a945`.

**I retract the original author report's statement that `./cmd/ensemble` followed the checked accepted/public recipe. That claim was false.** I inferred the command-package spelling rather than verifying it against the accepted tree. The accepted executable package is `./cmd/ch01`; its root is the importable `ensemble` module. The prior prose-lint result provided no evidence that `./cmd/ensemble` existed. The original report, receipts and seal remain unchanged as historical records; this new record corrects, rather than silently replaces, that unsupported claim.

## Read coverage and source basis

Before revision writes, I read the revision task, full formal review and parent verification, current reviewed chapter and original author diff. The original read ledger's complete mandatory reads were carried forward only after confirming their unchanged bytes; the current chapter was fully reread. All 24 prior prescribed-input pins and all 28 original author artifacts matched at entry. `read-coverage.json` distinguishes original-ledger identities from the chapter bytes read for this revision.

The source basis is the formal review's sealed, safe receipts for `edition-3-ch09-r1`: command tree, `go.mod`, README, the CLI's configuration/command parsing and `skillSpec` function, and exported construction aliases/signature. Their pins were verified against the review seal before use. A fresh read-only `git --no-pager rev-parse edition-3-ch09-r1^{commit}` confirmed `251ce88fb1677897dbcc6b33661096c5cb39cd8e`. No common-package/private/deeper implementation inspection was needed. No student executable was built or run, and the Go example was not compiled; its shape is checked against the accepted public source, not claimed as new runtime evidence.

Specific evidence:

- `student-cmd-tree.stdout.txt` lists the accepted `cmd/ch01` files. README lines 3–5 and 12 preserve that executable path; `go.mod` identifies the root library module.
- `student-cli.stdout.txt` lines 74–141 show the supported command selection (`chat`, `render LOG`, `dump`) and the no-argument JSON-lines path. An unknown command is rejected. R2 therefore uses existing human chat rather than inventing `prompt TEXT`.
- `student-public-construction.stdout.txt` contains the public `Config`/`SkillSpec` aliases and `NewAgent(config Config) (Agent, error)` signature. README lines 186–190 describe `Config.Skills` and nil-spec behavior. CLI lines 483–490 confirm the actual `SkillSpec{Directory, Primary, Variables}` shape. No new API is required.

These source receipts remain in `ch09-text-review-evidence/`, unchanged. The new evidence contains the full command outputs used to read them and the direct tag-verification receipt.

## Finite R1–R3 changes

The exact diff from the reviewed draft is `ch09-author-revision-evidence/revision-diff.stdout.txt`.

| Request | Correction | Preserved contract |
| --- | --- | --- |
| R1 | Replace both `./cmd/ensemble` references in the build recipe with the actual accepted `./cmd/ch01` path. | Importable library root remains distinct from its executable. No student rename or implementation change. |
| R2 | Run `/tmp/ensemble-ch09 chat`; tell the human to type `load the code-tools skill` at the chat prompt. Remove the unsupported `prompt TEXT` command and alternative-mode claim. | CLI and GUI still translate the same environment values into per-Agent construction inputs. The separate GUI module/launch remains. |
| R3 | Set `cfg.Skills = &ensemble.SkillSpec{Directory: ..., Primary: ..., Variables: ...}`, then call `app.NewAgent(cfg)`. | Existing error handling and future `RegisterVar` behavior remain. The prose explicitly retains nil-spec compatibility and permits different identifiers while preserving variables-before-primary-rendering. |

The grading recipe and complete registry section are byte-identical to the reviewed draft. In particular, the useful explanation separating per-Agent state ownership from sibling-package imports remains. All nine Markdown headings remain identical. The scratch demonstration was already using `chat` and was left untouched. There is no optional prose sweep, invented transcript, new feature, new gate, or manuscript process diary.

## Validation and limits

The current prose linter ran sequentially on the portable reviewed draft and the revised manuscript, both named `chapter-09.md`. `lint-before.json` and `lint-after.json` retain exact argv, cwd, allowlisted environment, target/linter pins, times, disk readings and exits; raw stdout/stderr are separate. Both used offline Go settings and no credential environment. The only Go commands in this revision were these prose-lint invocations.

| Check | Reviewed draft | Revised chapter |
| --- | --- | --- |
| Bytes | 20,767 | 20,888 |
| Whitespace-split words, including examples | 2,987 | 3,003 |
| Lines | 484 | 486 |
| Actual Markdown headings, excluding fenced examples | 9 | 9, identical |
| Lint exit | 0 | 0 |
| Person-gap metric | WARN 1324 / 1200 | WARN 1325 / 1200 |
| Negation forms | 4 / 4 | 4 / 4 |
| Short closers | 12 / 17 allowed, 52 paragraphs | 12 / 17 allowed, 52 paragraphs |

The existing soft lint warning is retained and disclosed. It is not a failed command, a warning-free pass, or a reason to add filler. Other printed metrics remain within thresholds.

The scoped chapter whitespace check and voice byte comparison exited zero. The no-index diff exited 1 because the files differ; that is its expected result. All 92 protected non-manuscript inputs/artifacts in `initial-pins.json` remained byte-identical in validation, including original reports, original author/reviewer seals and evidence, current guidance and voice files, and both immutable `/tmp` chapter copies. The accepted tag remains unchanged. This is a scoped verification, not a claim to have inventoried the entire student tree or repository; root owns broader protection checks. No publication inputs were read or hashed.

Free space stayed above the 1 GiB floor; the lint receipts show about 3.05 GB available. No tests, builds of student code, graders, browser/provider calls, credential reads, children, commits/tags/pushes, server restarts, shared-cache cleanup, or other chapter/source/guidance edits occurred. All invoked jobs finished.

## Process history, not erased

- The original root-reported transport error65 interrupted prerequisite reading. The original assignment resumed and its remaining full reads completed before edits. Its unchanged process ledger preserves this history.
- The original author read-only command guessed four nonexistent record paths and exited 1. The original exact errors and recovery remain preserved; the actual authoritative files were subsequently read in full.
- The formal-review parent verification discloses root's read-only job3041 pin-list/schema assumption failure and recovery. The reviewer separately disclosed its media-helper overread. This revision neither repeats those actions nor presents the earlier work as failure-free.
- The revision task refers to prior pager failures without incident details. The supplied original author/review process ledgers and parent verification do not identify such an incident; this report does not invent an exit code, timing or recovery. No pager was invoked in this revision; every git command used `--no-pager`.
- Two targeted/duplicate revision reads were clipped in inline output. Their complete tool outputs are retained in the new evidence; they did not substitute for mandatory full reads. No new command failure or lint failure occurred.

`process-history.json` records these distinctions and the retained command-output locations.

## Immutable history and final handoff

Before edits, the reviewed chapter was copied byte-for-byte to `ch09-author-revision-evidence/reviewed-draft/chapter-09.md`. It remains 20,767 bytes, SHA-256 `137f7fe93628d6434e5be6266933637b4bb0362b1926a6fe1e6da061db16447d`, matching `/tmp/cadence-ch09-text-review-initial/chapter-09.md`.

The old author/review seals are **unchanged**. Their `book/edition-3/chapter-09.md` entries now describe that preserved historical reviewed snapshot, not the revised working manuscript. No old pin is silently updated or claimed to match current chapter bytes.

The revised chapter is 20,888 bytes, SHA-256 `e2c82e17dc2cdeb4c580afe182ecef60bf6393fe10e234f56a7a255bbc9a5fc0`.

The new `ch09-author-revision-evidence/seal.json` is written last. Its exact, self-excluded inventory pins this report, the final chapter and every new evidence file; each listed artifact is verified UTF-8 with no NUL bytes. After sealing, only read-only verification occurs. This finishes the one authorized R1–R3 revision; root directly verifies the diff, report, pins and handoff without another formal review.
