# Building the second edition

The chapter files are the manuscript source. The assembled Markdown is committed
for reading on GitHub; EPUB and PDF are generated publication artifacts. The
first-edition files and build instructions remain in `edition-1/book/`.

## Prerequisites

Use Python 3 and pandoc. Print production also needs Google Chrome and the
Poppler tools `pdftotext`, `pdfinfo` and `pdfunite`. The cover-page conversion
below uses Pillow. The print build uses PyMuPDF to check actual PDF font sizes
(`python3 -m pip install pymupdf`).

Run from the repository root:

```sh
make book
make ebook
python3 book/build-print.py
```

Assembly includes the preface, Chapters 0–26, epilogue and author's note. Chapter
25 is the deliberate crossover placeholder. Both front-matter files identify
this as the second edition, October 2026. Keep them and the author's note
consistent when preparing a later release.

| Artifact in `book/` | Purpose |
| --- | --- |
| `the-art-of-building-ai-coding-agents.md` | Assembled manuscript with linked contents |
| `The Art of Building AI Coding Agents.epub` | Ebook |
| `The Art of Building AI Coding Agents - Print Interior.pdf` | Print interior |
| `The Art of Building AI Coding Agents.pdf` | Cover followed by the print interior |

## EPUB navigation and code

The EPUB build removes the manuscript's printed Contents section and omits
pandoc's `--toc`, keeping navigation out of the reading spine. It retains
`--toc-depth=1` and populated `nav.xhtml` and `toc.ncx` navigation. This is Bill's
explicit choice after a previous build displayed two contents pages. Do not
reverse it as a routine cleanup.

The build combines pandoc's default stylesheet with `book/epub.css`; supplying
only the override would discard the base styling. The override keeps highlighted
code spans inline and suppresses their empty anchors to avoid excess line spacing
in ebook readers. Browser rendering is useful evidence, but the final EPUB still
needs inspection in the intended reader.

`book/cover.png` is the default cover. An explicit cover can be supplied to
`scripts/build-ebook.sh`. Check that the EPUB includes it, has chapter-level
navigation, and has no duplicate printed contents or navigation page in its spine.

## Print pagination and cover

`book/build-print.py` generates print HTML from the assembled manuscript and print
front matter, removes duplicate contents, renders with Chrome, then adds measured
page numbers to the contents and renders again. It checks that those references
agree with the resulting chapter pages. The print CSS supplies mirrored margins,
page numbers, chapter breaks and wrapping code blocks on a 6 × 9 inch page.
Print code uses black text on white, retaining syntax bold and italic emphasis;
the EPUB keeps its separate color styling.
Long code tokens and table cells wrap rather than making Chrome shrink the
entire document. The build measures the PDF itself and rejects reduced body
or code type: body text must remain 11pt and code at least 9pt, allowing small
PDF-coordinate rounding differences. Checking CSS declarations alone misses
Chrome's automatic scaling.
Inspect the resulting PDF as well as the script's result, particularly long code
lines, tables, chapter starts and page footers.

To create the cover-prefixed PDF after a successful interior build:

```sh
python3 - <<'PY'
from pathlib import Path
from PIL import Image
cover = Image.open('book/cover.png').convert('RGB')
cover.resize((1800, 2700), Image.Resampling.LANCZOS).save(
    '/tmp/ensemble-cover-page.pdf', 'PDF', resolution=300.0)
PY
pdfunite /tmp/ensemble-cover-page.pdf \
  'book/The Art of Building AI Coding Agents - Print Interior.pdf' \
  'book/The Art of Building AI Coding Agents.pdf'
```

The cover-prefixed PDF should contain exactly one more page than the interior.
Verify both page counts and the 432 × 648 point page size with `pdfinfo`.

## Release check

Verify chapter coverage, edition metadata, contents links/page references, embedded
images and representative code/table layout. Check links in the assembled
Markdown and the generated formats; a repository-relative link needs a usable
published destination in an ebook. Keep the current publication record in
[docs/edition-2-publication.md](../docs/edition-2-publication.md), including any
remaining reader or print-preview limitations.

Commit only the intended source and publication files, preserving concurrent
application work. A local build or commit does not publish the book. Bill chooses
the destination and pushes or uploads the reviewed artifacts.
