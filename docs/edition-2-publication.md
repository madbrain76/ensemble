# Edition 2 publication and layout migration

## Latest reading copy: Bill's preface

At Bill's request, regenerate the book from his rewritten `book/preface.md`,
preserving it verbatim. Bill retains editorial ownership of the preface; future
author passes must leave it alone unless he explicitly requests an edit.
This supersedes the previous preface draft and its review, without relabeling
those historical records as review of Bill's new text.

The rebuilt interior has 371 pages; the cover-prefixed PDF has 372. All 30
contents page references agree with the new heading pages. EPUB chapter order,
navigation, cover and links pass the existing structural checks. The preface
is unchanged on disk, appears verbatim in the assembled Markdown, and its opening
and closing text are present in the EPUB and rendered PDF. Chapter text and the
11 pre-existing application/archive edits remain unchanged.

Build inputs, commands, artifact hashes and verification are retained in
[the rebuild record](../solutions/evidence/publication-preface-r2/build.json) and
[the checks](../solutions/evidence/publication-preface-r2/full-verification.json).
This is a mechanical regeneration for Bill's next reading, not a new editorial
review or a publication. No push or upload has occurred. The earlier independent
artifact review below describes its original version.

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

## Local reading copy, October 12

Bill will read the book before publication. No push, upload or copy to an external
publishing destination is authorized by this preparation. Crossover Chapter 25
remains reserved.

Artifact checkpoint `ca3f4764` contains assembled Markdown, EPUB, print HTML, the
370-page print interior and the 371-page cover-prefixed PDF. Every PDF page is
6 × 9 inches. All 30 printed contents references match the actual heading pages.
The EPUB includes all chapters in order, its cover and populated chapter-level
navigation, with no navigation page in the reading spine or duplicate printed
Contents. Repository links distinguish files from directories. The final artifact
hashes and measured checks are in [the verification record](../solutions/evidence/publication-tooling/full-verification.json).
The coordinator independently matched all five artifact hashes.

Current default grades for Chapters 1 and 26 and archived original Chapter 1 each
score 100/100 in new relocation runs. The real nested-module fakevendor loop,
maintained tooling build/vet, selected regressions and exporter tests pass. Current
CLI and GUI build to temporary binaries; this establishes buildability, not
acceptance of the separate uncommitted crossover changes. Full preface/chapter
prose lint exits zero with soft warnings retained. Reader-entry local links pass.
See [the tooling report](../solutions/evidence/publication-tooling/README.md) for
exact scope and the distinction between compilation and executed tests.

The epilogue is Astra's own first-person account from Codex. Independent review
resolved one attribution correction: the student raised the memory-trigger
contradiction, and Bill settled the intended behavior. The final text preserves
failed attempts, actual results and the unfinished crossover. Front matter and
all mechanically relocated chapter paths are independently accepted.

Independent rendered-page review accepts the local reading copy: cover, contents,
front matter, representative chapter openings, multi-page tables, long command/JSON
blocks, placeholder25, epilogue and author note are readable. All 370 interior
footers and all 30 contents references were checked. An extracted EPUB Gateway
page also renders readably in Chrome. The [independent review](edition-2-notes/publication-review.md)
records the precise scope. Kindle Previewer, KDP print preview, physical proof
and Bill's editorial read-through have not been performed.
These local artifacts are reading copies, not a published edition.
