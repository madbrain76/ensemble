# Edition 3 Chapter 4 — frozen post-student author draft

## Status and ownership

Draft for the parent's one formal independent text review; **not text acceptance**. No reviewer or other child worker was launched. No author revision round has been used. Await the parent's exact findings, if any; at most one findings-driven author revision is authorized. No optional polishing or new requirements are proposed.

Only `book/edition-3/chapter-04.md` and this new report were written. Application, grader, public brief, student repository, earlier chapters, historical book, references and existing evidence were left alone. Bill's preface was preserved verbatim. No staging, commit, tag, push, reset, clean, installation, restart, publication build, provider call or grader run was performed by this author.

## Exact baseline and frozen identities

- Main HEAD: `bac07e644facb25ba3e5bffb2961674763fed4c4`.
- Chapter baseline: HEAD's unchanged `book/edition-3/chapter-04.md`, 1,000 lines; raw SHA256 `016c4e88af4bcb657aad6733e5f65aa09509b0397e84a37059ca245c588eaa55`.
- Frozen chapter: same path, 1,029 lines; raw SHA256 `ae7f6fbe1048400f857e8fb3eeefbe513d0a9bab492e1eff1d13a10ff266f37f`.
- Report baseline: `docs/edition-3-notes/ch04-author-review.md` did not exist. There is no baseline report hash.
- Frozen report: this file. Its exact **raw file SHA256 is supplied in the accompanying frozen submission to the parent**, computed after the final write. A raw full-file digest cannot be embedded in the same file without changing that digest; no self-referential or pre-final hash is represented as the final report hash.
- Accepted student tag, verified in the separate student repository: `edition-3-ch04-r1`, peeled commit `ec4a350f564e820d398f9125fd73a9b3b27129cd`; student HEAD matches, and student status was clean. The tag is not a tag in the main book repository.
- Both `book/preface.md` and `book/edition-3/preface.md`: SHA256 `915cf8c531169c5aacc4cb81bd743eeba8a181d20bc779d936ac5f7ab9998d96`.
- Protected original `book/chapter-04.md`: SHA256 `b5b89281bdb8aaeee58fc2d5777c2fec8d70f6d048fb83f675487219188a308d`.

The pre-edit preservation snapshot covers all 633 pre-existing modified/untracked files reported by `git ls-files -m -o --exclude-standard -z`, excluding the two owned paths. Its aggregate SHA256 is `ba166fda9861e11348d0061fef9ed0521f1f4865501ca1365d806d7c91c0e998`. The aggregate is the SHA256 of compact JSON containing sorted unique `[relative path, raw file SHA256]` pairs for existing regular files. This records dirty work without incorporating or publishing its contents. Final preservation comparison is part of the frozen submission checks, not a claim that those files were clean.

## Full read coverage before formal editing

`CURRENTAGENTS.md` is not a literal root file; the current root instructions are `AGENTS.md`. The actual current files below were read completely. Long reads were split; truncated combined outputs were recovered with standalone reads. The workflow's final range 2,284–2,434 was explicitly recovered. The current assigned chapter was also reread completely in ranges 1–350, 351–700 and 701–1,000 before any edit. Historical author-first orders were treated as superseded by the current Edition 3 student-first procedure. No compaction/handoff was used in place of these readings.

| Required file | Complete coverage |
| --- | --- |
| `AGENTS.md` | 1–382 |
| `.agents/skills/ensemble-coding/SKILL.md` | 1–397 |
| `docs/agentic-codebook-workflow.md` | 1–2,434 |
| `docs/edition-2-carryover.md` | 1–1,110 |
| `book/voice.md` | 1–632 |
| `book/chapter-writing-procedure.md` | 1–622 |
| `docs/edition-3-execution.md` | 1–53 |
| `docs/edition-3-student.md` | 1–71 |
| `docs/edition-3-readiness.md` | 1–22 |
| `docs/edition-3-carryover.md` | 1–71 |
| `.agents/skills/edition-3-author/SKILL.md` | 1–12 |
| `book/edition-3/chapter-04.md` | baseline 1–1,000 |

Authoritative chapter records, all fully read (871 lines total):

| File under `docs/edition-3-notes/` | Complete coverage |
| --- | --- |
| `ch04-accepted.json` | 1–28 |
| `ch04-execution.md` | 1–107 |
| `ch04-entry-audit.md` | 1–118 |
| `ch04-grader-review-r1.md` | 1–58 |
| `ch04-student-review.md` | 1–110 |
| `ch04-student-review-r1.md` | 1–96 |
| `ch04-code-review.md` | 1–158 |
| `ch04-code-review-r1.md` | 1–132 |
| `ch04-candidate.json` | 1–33 |
| `ch04-candidate-r1.json` | 1–31 |

Direct verification beyond the records: read the complete current `cmd/grade-edition3/main.go` and the relevant accepted `internal/agent/dispatch.go` excerpt (lines 45–185) in the student repository; verified the accepted tag/HEAD and clean student status. The adapter source explicitly restricts its route to Chapters 1 and 4, takes an executable positional target, accepts distinct `-core` and `-consumer` module paths, and documents flag-before-positional ordering. No old answer/reference implementation was needed for this author pass.

## Three bounded hunks and demonstrated basis

The chapter diff is **40 insertions, 11 deletions, three hunks**. All headings are unchanged. The purpose is clarity after the observed student run, not author-led redesign.

1. **§4.5, after the existing one-shot explanation (baseline after line 407; frozen 408–414).** Add one paragraph separating the pending options consumed to observe the current setter call from the newly pending defaults plus that setter's explicit fields. Use one omitted-cap example, without code, helper names or a regression inventory. The baseline already explicitly included `tool_limits` in the next-call consumption rule; this does not newly require that behavior. Nevertheless the initial candidate merged spent settings back into new pending state. C1 and its accepted narrow fix demonstrate why the distinction earns a brief explanation. `ch04-student-review-r1.md` describes the public counterexample and three omitted-field tests; `ch04-code-review-r1.md` closes C1 with public-API and mutation evidence. C1 is closed, not an open code question.

2. **§4.9, the build/grading instructions (baseline 787–798; frozen 794–825).** Replace the historical directory-target command with a build of the actual CLI package, then the supported Edition 3 Chapter 4 invocation with executable positional target and explicit actual core/separate consumer module paths. Show Chapter 1 lifetime grading on the **same executable** from the main book repository root. Explain optional `-consumer-case` placement before the positional target, success flags/exits rather than score alone, and the adapter's present Chapters 1-and-4-only scope. Keep the existing case format and default `shout` contract. The entry audit demonstrates why the historical target could not serve both source inspection and execution for a legitimate public-root-library/nested-CLI layout. The accepted thin adapter and public student brief already supply this route; this hunk makes the chapter agree with them. `.` remains explicitly valid for a module-root command; the nested command spelling is a placeholder, not a mandated package name or layout. Historical graders were not edited.

3. **§4.10, live build instructions (baseline 887–894; frozen 914–923).** Apply the same actual-command-package distinction to the live core build, retaining a root-command alternative. This prevents `go build ... .` from targeting the public root library when the runnable command is nested. The provider setup, scratch-directory workflow, JSON-lines/chat distinction, Delve prompt, useful server task and consumer exercise remain unchanged. The consumer example already builds its separate runnable module and did not need a new layout rule.

## Why the rest stays unchanged

The chapter already teaches observation rather than cancellation, full disk output, ordered dispatch, per-Agent/per-Ensemble ownership, source-level boundary review, root logging, a real consuming module, and process-level verification beyond a model's claim. The evidence did not justify rewriting those explanations or adding requirements. The deadlock diagnosis, debugger-paused-by-design lesson, disk-spill story, callback-versus-deadline reasoning, and the helpful/harmful power of a sharp tool remain load bearing. So do the long-form explanation and the existing 160-point check table. No arbitrary length target or wholesale pruning was applied.

The existing chapter correctly asks for **Delve**, not GDB; no C/GDB exercise was added. No revision was earned for provider budgets, model selection, logging policy, coverage thresholds, job lifetime rules or process architecture. The old/historical live anecdotes and evidence links retain their attribution. The optional historical reference remains an optional reference, not a new student prerequisite. All these areas were left alone deliberately, not inferred correct from the score alone.

## Author-side checks actually performed

- Manually read the whole current chapter and the final three-hunk diff against HEAD, checking the new sentences against the accepted adapter/public brief and the closed C1 evidence. This is an author self-check, **not independent text acceptance**.
- `git --no-pager diff --check -- book/edition-3/chapter-04.md`: exit 0, no whitespace errors. For the new untracked report, `git --no-pager diff --no-index --check -- /dev/null docs/edition-3-notes/ch04-author-review.md` also exited 0; a separate line scan found no trailing whitespace.
- A Python fence scan found 12 complete pairs in the baseline and 13 in the draft, with no unclosed fence. The additional pair is the explicit build command.
- The same scan confirmed all chapter headings are byte-for-byte unchanged.
- Checked all three local Markdown link paths from the chapter's actual directory: `../../docs/edition-2-notes/ch04.md`, `../../solutions/ch04/ch05/`, and `../../solutions/evidence/ch05/live.txt` all exist. This is path existence, not publication rendering or a new evidence audit.
- Ran `/bin/bash -n` on both labeled shell blocks (frozen opening/closing fence lines 798/800 and 804/807): both exit 0. This checks syntax, not executable builds or placeholder-path existence.
- Compared preserved preface/original chapter and dirty-work hashes at freeze; exact final results accompany the submission.

The parent will run prose lint and owns the one independent text review after freezing the inputs. I did not run that lint or start a competing reviewer. No publication rendering, broken-anchor crawl, code compilation, application tests, graders, provider sessions or new paid/live exercise are claimed by this author.

## Code and live evidence limits

Code acceptance is established by `ch04-accepted.json`, not by this prose pass. The R1 code reviewer independently records full/race/vet success, public-API C1 closure and sensitivity controls, Chapter 4 `passed:true` 160/160, and Chapter 1 `passed:true` 100/100 on the same corrected executable, with real 62-second header/body phases. These are recorded earlier runs, **not author reruns**. The grader adapter itself was accepted after its one help/diagnostic correction; there is no reopened grader issue here.

The real-provider/debugger/server evidence remains attached to initial candidate `062c146cafa9b105745f6b44921160198bbaee32` and its original artifacts. R1 `ec4a350…` has source, regression, mutation, public-grader and audited local synthetic-provider CLI/consumer evidence, not a fresh vendor campaign. Preserve the initial Anthropic failure and bounded recovery, the OpenAI refusal-only result, and the eight-second workload as qualified observations. Eight seconds is not a 90-second demonstration, and a model's report is not process-side-effect proof. None becomes a new chapter specification or a fabricated fresh success. No additional provider or paid test is requested.

## Handoff

Freeze the chapter/report pair using the exact raw hashes in the accompanying submission. The parent owns formal prose lint, the one independent text review and final text acceptance. If that review supplies concrete findings, return them verbatim for the sole allowed author revision. Otherwise no further author change is requested.
