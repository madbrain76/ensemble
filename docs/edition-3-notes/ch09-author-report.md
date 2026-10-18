# Edition 3 Chapter 9 — bounded post-student author report

## Status and authority

The bounded author pass is complete and ready for root's **one formal independent text review**. This report is not text acceptance, publication, or permission for another implementation/review cycle. No children were created. The author has not consumed the formal text-review pass or its maximum one permitted revision.

Authority: `/tmp/cadence-ch09-author-task.txt`, SHA-256 `6eceb70e1746991310d6d763868d647c9aa21eb2ece1da87f0dcb44c0031ec37`. The same assignment resumed after root reported transport error65, “Responses inference did not complete,” during prerequisite reading. Root observed no final/report/evidence then. The interruption is process history, not an implementation failure or a second author assignment. No timestamp or duration is invented for it.

The settled code and operational decisions were inputs, not reopened gates. This author ran no tests, graders, student builds, browser sessions, provider calls, or implementation investigation beyond the permitted CLI/media-map contract inspection. No credentials, branch/tag operations, commits, publication artifacts, other chapter files, or Bill's preface were opened or changed for this pass.

## Read coverage and preservation

All 20 mandatory full-file prerequisites, including the three role/coding skills, were read through EOF **before the first author write**. `ch09-author-evidence/read-coverage.json` lists each path, line count, byte count and SHA-256. The long workflow was covered in contiguous ranges through its actual EOF at line 2434; Edition 2 carryover lines 1–350 preceded the interruption and 351–1110 completed it afterward. Initial truncated AGENTS output was recovered with bounded reads. Later targeted rereads recovered exact task and paragraph/CLI details without replacing any required full read.

The source coverage includes AGENTS, the complete workflow and Edition 2 carryover, current canonical voice and chapter-writing procedure, Edition 3 execution/student/carryover guidance, Chapter 9, the student and code reviews, both root decisions and parent verifications, the operational exercise, and the accepted record. Additional permitted verification read the current `cmd/lintprose/main.go`, the `cmd/grade-edition3/main.go` CLI contract, and only lines 1–43 of `internal/gradeedition2/ch05_media.go` for `ReadMediaModels`. A mistaken record-read command also read `ch09-execution.md`; it did not modify that record or use it to expand scope.

All 24 baseline inputs matched before editing. Final verification found the expected Chapter 9 change and **all other 23 pins unchanged**, including the current concurrent voice, procedure, skills, execution/readiness/carryover documents, graders and settled records. `cmp -- book/voice.md book/edition-3/voice.md` exited zero. Both voice files retain SHA-256 `5a54dabd30b7846f4cca2502110481397343ad2a232b1fcaa0bd19b1ff3a4d75`. No repository-wide tracked-file scan or other-chapter/publication read was used. Preservation of unread files is a scoped-write assertion, not a claim that they were independently hashed.

## Changes and evidence

There are three bounded replacement sites, all in `book/edition-3/chapter-09.md`:

1. **Build the CLI, not the library root.** The core-module command now targets `./cmd/ensemble`. This follows the inherited importable-core/separate-command layout and the public student recipe's distinction between executable and source identities. The example was checked against the records; it was not executed by this author.
2. **Replace the historical Edition 2 grading recipe.** The example now uses `grade-edition3 -ch 9` with explicit core, consumer, workflow, stream-exercise and GUI modules; explicit streaming/GUI model identifiers; an explicit media-map path; and the built CLI as the final argument. It includes Chapter 1 on that same frozen CLI. The explanation distinguishes module roots from executable targets, keeps flags before the positional target, identifies the media-map keys checked by `ReadMediaModels`, and avoids historical model defaults or live-discovery claims. `-consumer-case` remains optional. Supplying a media map in the example is not a new grader requirement: it avoids implicit historical fixtures. The GUI module is required by the current CLI gate, even though a reader may use the headless interface day to day.
3. **Explain the verified registry boundary failure.** One paragraph in “The Registry” explains why a per-Agent lifetime does not permit importing a sibling registry implementation into Agent. Agent uses the common lifecycle interface; root constructs and injects the concrete registry. The student review records the actual first-candidate sibling-import failure and its interface-based repair. This makes an existing architecture obligation explicit, without prescribing a new exported API or new behavior.

`exact-diff.stdout.txt` is the complete unified diff against the pinned `/tmp` chapter copy. The no-index diff command's exit 1 means differences were found, not a validation failure. No unrelated prose was rewritten. The frontmatter contract, parsing and graph rules, exercise obligations and score table, failed-load atomicity, fixed system text, pending unload, declaration/dispatch boundary, rendering lock rationale, provenance, API example, memory limitation, scratch-file exercise and live/browser qualifications remain. No demonstrated acceptance limit was removed and no later feature was pulled forward.

## Validation and editorial judgment

The only Go execution in this author pass was the current prose linter. Exact argv, cwd, allowlisted credential-free environment, elapsed time, disk headroom, target/linter hashes and exit codes are in `lint-before.json` and `lint-after.json`; stdout and stderr are separate unedited files. Both used `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOWORK=off`, and `GOENV=off`.

| Check | Before | After |
| --- | --- | --- |
| Chapter bytes | 19,058 | 20,767 |
| Words (whitespace split, including examples) | 2,763 | 2,987 |
| Lines | 454 | 484 |
| Markdown headings outside fenced examples | 9 | 9, identical |
| Lint exit | 0 | 0 |
| Person gap | 1,133 / 1,200 | **WARN: 1,324 / 1,200** |
| Negation forms | 2 / 4 | 4 / 4 |
| Short closers | 13 / 16 allowed, 50 paragraphs | 12 / 17 allowed, 52 paragraphs |

All other current lint metrics remained within their printed thresholds. The after-run warning is retained, not called a clean warning-free pass. The new command block and explanatory contract lengthen the opening instructional span. I kept the needed distinctions and existing prose rather than inserting a personal aside or reshaping unrelated paragraphs to satisfy a heuristic. Root's text reviewer can evaluate that judgment in the same bounded review. There were no failed lint runs and no broad style cleanup.

The scoped `git diff --check -- book/edition-3/chapter-09.md` exited zero. A separate recorded whitespace inspection covers the final report and author text evidence, including untracked files. `chapter-statistics.json` distinguishes actual Markdown headings from hash-prefixed lines inside examples. All commands completed; no job, server, child, cache deletion, or temporary worktree remains owned by this author. Free space remained above the 1 GiB guard; the lint runs recorded roughly 3.20 GB available.

## Failures retained

- Root's reported transport error65 interrupted prerequisite reading. The same task resumed; remaining reads finished before edits.
- One read-only `cat` command guessed four nonexistent record filenames and exited 1. `failed-record-read.errors.txt` retains the exact four errors; `process-history.json` retains the command, original tool-log path/hash, observed exit, and recovery. The authoritative task was reread and all actual decision/verification/acceptance files were then read in full. No content was fabricated to replace the missing names.
- Output truncations are recorded in `process-history.json` and coverage notes. They were recovered where necessary for prerequisite coverage. A later truncated duplicate/targeted command-guidance read did not stand in for a required full read.

## Final identity and handoff

Baseline chapter: 19,058 bytes, SHA-256 `1fc63f4a8b2680dbfd89658028c0a53d6a7c8dcd762c881e8f7b8ee0ba26e4ca`.

Final chapter: 20,767 bytes, SHA-256 `137f7fe93628d6434e5be6266933637b4bb0362b1926a6fe1e6da061db16447d`.

`ch09-author-evidence/seal.json` is written last and pins the final chapter, this report, and every evidence file except itself. It also records the final prescribed-input comparison and the non-self-hashing rule. All author outputs are immutable after that seal. Root owns verification, the single formal text review and any acceptance decision; this author has made no publication or completion claim for those later gates.
