# Edition 2 publication and layout migration

## October 14: reordered manuscript inputs

The current chapter files and index now follow the approved dependency order:
old17–23 become15–21, old15/16 become22/23;0–14 and24–26 stay fixed.
`assemble-book.sh` consumes the sorted numbered files, and `book-toc.py`
derives titles and anchors from those inputs. No publication artifact is rebuilt
by this change. The reading-copy counts and hashes below retain their earlier
source identity, and Bill's preface remains unchanged. Revised cumulative
solution exports are a separate reviewed step; see the
[chapter map](edition-2-chapter-map.md) and
[renumbering record](edition-2-notes/codex-compaction-renumber.md).

## Print contrast update

At Bill's request, print code now uses black text on white backgrounds, including
comments and strings. Bold and italic syntax cues, spacing and wrapping remain.
The rebuilt PDFs retain 371 interior pages / 372 with the cover and the same 30
contents destinations. Bill's preface, all manuscript text and the EPUB are unchanged.

Implementation `4e21dfad` and independent review `8217c6d3` establish the bounded
print correction. See the [build and browser evidence](../solutions/evidence/publication-print-light/README.md)
and [independent rendered-page review](edition-2-notes/publication-review.md).
This is a local reading copy; no publication or physical proof is claimed.

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

## October 13: correct print scaling after publication

Bill reports the book is published and requests a second-edition readability
fix. Comparing the assembled manuscripts showed 158,623 whitespace-separated
words in Edition 2 versus 129,538 in Edition 1. The 372-page second-edition
interior nevertheless printed its declared 11pt body at 7.33pt and its 9pt
code at 6pt. Long code tokens and table cells overflowed; Chrome scaled the
whole document down. The archived first edition also shrinks, but less.
Earlier checks of page dimensions and representative screenshots missed this.

The coordinator corrected print CSS to wrap long tokens and override pandoc's
inline-block code spans and hanging indent. At the printable width, the full
HTML now has no horizontal overflow. The PDF measures 10.99pt body (rounding
from 11pt) and minimum 9pt code. The print build now checks those actual sizes
with PyMuPDF; running that check against the previous PDF correctly rejects
its 7.33pt/6pt type. No manuscript or first-edition content was changed.

The rebuilt interior has 697 pages; the cover-prefixed PDF has 698. All pages
remain 6 x 9 inches, all 30 contents references match actual chapter pages,
and no extracted word extends outside its page. Visual inspection covered the
previously overflowing JSON example, code, tables and a chapter opening.
The preface, chapter sources, assembled Markdown, cover image and EPUB retain
their pre-fix hashes. These are local replacement PDFs; no upload or change to
the published listing was made. A physical print proof remains unperformed,
and a paperback cover's spine must match the new page count before reprinting.

## October 13: larger trim and type

Bill selected the recommended 7 × 10 inch format with 12pt body text and 10pt
code. The print stylesheet now uses those sizes, 11pt tables and larger section
headings. Mirrored inner margins are 0.875 inches, outer margins 0.5 inches,
and top/bottom margins 0.75 inches. Removing pandoc's browser body padding lets
the page margins determine the text width. The overflow fix remains in place;
the build's measured-type check now requires 12pt/10pt.

The rebuilt interior has 562 pages; the reading PDF has 563, including its
cover. Every page measures 504 × 720 points. The existing cover artwork is
centered proportionally on the larger reading-PDF page rather than stretched;
the source image is unchanged. All 30 printed contents references match actual
chapter pages. Actual type measures 12pt body and minimum 10pt code, the HTML
has no horizontal overflow at the printable width, and no extracted word falls
outside a PDF page. Representative table, code and chapter-opening pages were
visually inspected. Manuscript, preface, chapter sources and EPUB are unchanged.
No external publication action was taken; the physical paperback cover still
needs a separate trim/spine layout and a print proof.

## October 13: full artifact rebuild

At Bill's request, rebuilt assembled Markdown, EPUB, print HTML, print interior
and cover-prefixed PDF from the current sources, including Chapter 25's link
to the live crossover record. The PDFs remain 562/563 pages at 7 × 10 inches,
with measured 12pt body and minimum 10pt code; all 30 print contents references
match. EPUB contains Chapters 0–26 and all closing sections, its cover and code
stylesheet, with no navigation document in the spine or duplicate Contents
section. All 1,895 internal links/resource references resolve. Its navigation
has 31 TOC entries (30 book sections plus front matter); the build script's
33-link count also includes two landmark links. Rebuild only: no source prose
edits, external upload or paperback wraparound-cover generation.
