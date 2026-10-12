# Edition 2 publication and layout migration

Bill authorized this sequence on October 12, 2026: finish the Gateway chapter,
archive Edition 1, promote Edition 2, and prepare its first publication. The
crossover has begun; Chapter 25 deliberately remains a placeholder. Producing
Edition 3 with Ensemble remains future work.

## Completed physical migration

Promotion commit `564e316f3f8aeac2ec26fe8edfbd8cdf42aa350b` follows accepted
Gateway review `63709086` and checkpoint `c35aab73`. It maps 8,386 tracked paths
and preserves their committed modes and blobs. The resulting tree has 9,020
tracked files, including separate original/shared grading support. All 8,386
moved working-file hashes and 443 protected local-file hashes match the pre-move
capture. Independent reviewer verification confirms those counts and identities.

The migration's per-file shared-index reconciliation collided with a concurrent
status refresh after the relocation commit was written. The coordinator rebuilt
only the changed index entries in one `update-index --index-info` operation,
verified the committed tree against the prepared entries and confirmed an empty
staged diff. No working bytes, commits or tags were reverted. The same 11 prior
modified files remain uncommitted at their new paths: ten active crossover files
and the historical first-edition chapter index. Their acceptance is not implied.

Current entry-point, link and publication checks follow below; physical
preservation alone does not establish that those commands work.

| Former path | Current location |
| --- | --- |
| `agent/` | `edition-1/agent/` |
| First-edition content in `book/` | `edition-1/book/` |
| First-edition `solutions/chNN/` | `edition-1/solutions/chNN/` |
| `solutions/edition-2/main/` | `agent/` |
| `book/edition-2/` | `book/` |
| Other contents of `solutions/edition-2/` | `solutions/` |
| Original `cmd/grade/` | `edition-1/cmd/grade/` |

Preserve the original grading module and its support under `edition-1/`, with
the original module path, so it can run from that directory. The current revised
graders still use `internal/grade` and the fake providers. Retain that shared
support at the top level as well; moving the old command does not justify a
rewrite of the grading library. Keep `cmd/grade-edition2` as the explicit current
command and make the top-level Makefile select it.

Preserve original build scripts, Makefile and README in the archive. Shared craft
guides, build assets and book tooling also remain available for the current book.
Archive old chapter-specific maintenance scripts; distinguish retired Edition 2
attempt scripts from the current exporter rather than silently making obsolete
checks look current.

The move preserves tracked blob identities for the original edition and frozen
solutions. Bill's dirty application files and local untracked material move with
their directories, without becoming accepted source or being staged. Historical
manifests, scores and receipts retain their source paths and revisions. This map
explains their new locations. There are no symlinks or new chapter numbers.

## Assignments and final checks

- Coordinator: physical moves, preservation check, current README and agent
  instructions, chapter paths, external launcher, final epilogue and assembly.
- Author: preface, front matter, author's note and narrow Chapter 0 corrections;
  preserve the documented voice and distinguish completed work from crossover.
- Coder: current grading defaults, affected test paths and publication tooling.
- Reviewer: independent Gateway acceptance, then relocated entry points,
  preservation and publication prose/artifact checks.

Verify relocated default grading and original archive grading, affected module
checks, runnable example paths and local links. Do not repeat paid demonstrations
for a path-only move. Build the assembled Markdown, EPUB, print interior and
cover-prefixed PDF. Preserve the established EPUB navigation behavior. Inspect
chapter coverage, edition metadata, page references, code layout and images.

The new epilogue is the coordinator's own first-person account from Codex, written
after the manuscript is complete. The archived epilogue's narrator running inside
Ensemble is a different contribution and stays intact. Local publication files
and commits are preparation; no external upload or push has occurred.

## Current preparation status

The external launcher now selects `agent/skills`; its existing binary and live
session were preserved. A normal session restart is needed to pick up that path.
Current README, chapter index and governing-document headers identify the promoted
layout. Historical literal paths and source hashes remain unchanged; relative
Markdown links in the notes now resolve to relocated current or archived artifacts.

Independent review accepts the new preface, front matter, author note and narrow
Chapter 0 corrections. The detailed [publication review](edition-2-notes/publication-review.md)
tracks the remaining tooling, chapter-path, epilogue and artifact gates.

## Uncommitted application work: provenance correction

The coordinator initially attributed the 11 pre-existing modified files to Bill's
crossover work without checking who wrote them. Bill reported that CodeRhapsody
had committed its own work and disclaimed these changes. Read-only inspection of
the separate live Ensemble session's event log established the source of the ten
application edits: Bill requested approximately 5× speech with Eloquence
calibration and ordinary chat completion for the human-facing root, reserving
`submit` for children. The session records the requests and corresponding edit
operations on October 12 UTC (October 11 locally). The existing untracked
`ch25-crossover.md` records the speed work and its validation limits.

These edits remain uncommitted, unchanged and outside the accepted snapshots;
publication preparation neither accepts nor discards them. The two added roadmap
items in `edition-1/book/chapters.md` also predate the move, but their writer has
not been established. No private conversation log was copied into the repository.
