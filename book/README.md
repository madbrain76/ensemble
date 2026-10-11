# Book Build Procedures

*The Art of Building AI Coding Agents* by Bill Cox.

## Prerequisites

- **pandoc** — markdown to HTML/EPUB conversion
- **Google Chrome** — headless PDF rendering (page numbers, CSS @page)
- **pdftotext** (poppler) — page number extraction for two-pass TOC
- **pdfinfo** (poppler) — PDF verification
- **pdfunite** (poppler) — cover prepending for website PDF
- **Python 3** — `build-print.py` and `book-toc.py`

Install poppler on macOS: `brew install poppler`

## File Layout

```
book/
├── README.md                          ← this file
├── preface.md                         ← preface
├── chapter-01.md ... chapter-26.md    ← chapter source (single source of truth)
├── epilogue.md                        ← epilogue
├── note-from-bill.md                  ← author's note (update edition number here)
├── front-matter.md                    ← EPUB front matter (update edition number here)
├── front-matter-print.md             ← print PDF front matter (update edition number here)
├── cover.png                          ← book cover image
├── print-header.html                  ← CSS for print PDF (margins, fonts, page numbers)
├── build-print.py                     ← two-pass print PDF builder
│
├── the-art-of-building-ai-coding-agents.md   ← assembled manuscript (COMMITTED)
├── interior.html                              ← generated (pandoc, no page numbers)
├── interior-print.html                        ← generated (pandoc + TOC page refs)
├── The Art of Building AI Coding Agents.epub  ← KDP ebook
├── The Art of Building AI Coding Agents - Print Interior.pdf  ← KDP paperback
└── The Art of Building AI Coding Agents.pdf   ← website (cover + interior)

scripts/
├── assemble-book.sh                   ← assembles chapters into single markdown
├── build-ebook.sh                     ← builds EPUB with structural checks
└── book-toc.py                        ← generates linked TOC for assembled markdown
```

## Publishing Checklist

Amazon will reject an update that changes more than 10% of the book,
insisting it must be a new edition. Publish after every chapter to stay
under the threshold.

### Step 0: Update edition numbers

Three files carry the edition number — update all three:

```bash
# book/front-matter.md        →  "Seventh Update, October 2026"
# book/front-matter-print.md  →  "<p>Seventh Update, October 2026</p>"
# book/note-from-bill.md      →  "This is the seventh update of ..."
```

### Step 1: Assemble the manuscript

```bash
scripts/assemble-book.sh > book/the-art-of-building-ai-coding-agents.md
```

This concatenates preface + all chapters + epilogue + note-from-bill,
and runs `book-toc.py` to generate a linked table of contents. The output
is committed on purpose — the README links to it so the book is readable
on GitHub.

### Step 2: Build the EPUB (ebook)

```bash
scripts/build-ebook.sh [cover-image]
```

If no cover image is passed, it defaults to `~/Downloads/The Art of
Building AI Coding Agents (1)_upscayl_4x_digital-art-4x.png`. The
`book/cover.png` also works:

```bash
scripts/build-ebook.sh book/cover.png
```

The script:
- KEEPS the manuscript's "# Contents" section as the ebook's visible
  contents page. An earlier version stripped it on the theory that Kindle
  builds its own nav; Amazon did not reliably surface one, and the
  published ebook ended up with no table of contents at all.
- Does NOT pass `--toc` (would ALSO put nav in the spine — that is what
  produced two contents pages originally). One visible page, ours.
- Passes `--toc-depth=1` (keeps nav entries at chapter level, not 426 sections)
- Verifies: nav not in spine, nav entries < 60, all images embedded, and
  that a visible Contents page survived into the reading flow

Output: `book/The Art of Building AI Coding Agents.epub`

### Step 3: Build the print interior PDF (paperback)

```bash
python3 book/build-print.py
```

This is a two-pass build with automatic HTML cleanup:
0. **Pass 0 — regenerates `interior-print.html` from the assembled manuscript
   via pandoc.** You no longer run pandoc by hand; see the warning below.
1. Cleans `interior-print.html`: forces page break after the pandoc TOC,
   strips the self-referential "Contents" entry from the pandoc TOC, and
   removes the markdown-generated "# Contents" section (duplicate of
   pandoc's `--toc` — it's for GitHub, not print)
2. Pass 1: renders PDF with Chrome headless, reads where each chapter landed
3. Pass 2: injects page numbers into the TOC, re-renders
4. Verifies pagination didn't shift between passes (re-runs if it did)

The script handles everything — no manual steps needed.

**If you need to regenerate `interior-print.html` first** (e.g. after CSS
changes), run pandoc manually before the two-pass build:

```bash
cd book
pandoc front-matter-print.md the-art-of-building-ai-coding-agents.md \
  -o interior-print.html \
  --from=markdown --to=html --standalone --toc --toc-depth=1 \
  -M title="The Art of Building AI Coding Agents" \
  -M author="Bill Cox" \
  -H print-header.html
```

The CSS header file (`book/print-header.html`) includes:
- `@page` with mirror margins: outside=0.25in, gutter=0.625in,
  top/bottom=0.75in. Even/odd pages swap left/right via `@page :left`
  and `@page :right` rules. Amazon's preview shows the mirrored layout.
- `@page { @bottom-center { content: counter(page); ... } }` — page numbers
- `.tocpg { float: right; }` — TOC page number positioning
- `h1 { page-break-before: always }` — chapter breaks (h1 only; h2 subsections flow continuously)
- `pre { white-space: pre-wrap }` — prevent code overflow on 6-inch page

Output: `book/The Art of Building AI Coding Agents - Print Interior.pdf`

**Important:** The Chrome command inside `build-print.py` does NOT use
`--no-pdf-header-footer`. That flag suppresses CSS `@page` margin boxes,
which is where page numbers live. If page numbers disappear, check that
flag first.

### Step 4: Build the website PDF (cover + interior)

```bash
cd book

# Convert cover to 6×9 PDF page at 300 DPI
python3 -c "
from PIL import Image
img = Image.open('cover.png')
img_resized = img.resize((1800, 2700), Image.LANCZOS)
img_resized.save('/tmp/cover-page.pdf', 'PDF', resolution=300.0)
"

# Prepend cover to interior
pdfunite /tmp/cover-page.pdf \
  "The Art of Building AI Coding Agents - Print Interior.pdf" \
  "The Art of Building AI Coding Agents.pdf"
```

Output: `book/The Art of Building AI Coding Agents.pdf` (one more page
than the interior — the cover)

### Step 5: Copy artifacts

```bash
# To ~/Downloads/ for KDP upload
cp "book/The Art of Building AI Coding Agents.epub" ~/Downloads/
cp "book/The Art of Building AI Coding Agents - Print Interior.pdf" ~/Downloads/
cp "book/The Art of Building AI Coding Agents.pdf" ~/Downloads/

# To website repo for coderhapsody.ai
cp "book/The Art of Building AI Coding Agents.pdf" \
   ~/projects/coderhapsody/web/website/public/the-art-of-building-ai-coding-agents.pdf
```

### Step 6: Commit

```bash
# Ensemble repo (commit as CodeRhapsody)
git add book/the-art-of-building-ai-coding-agents.md \
       book/front-matter.md book/front-matter-print.md book/note-from-bill.md \
       "book/The Art of Building AI Coding Agents.epub" \
       "book/The Art of Building AI Coding Agents - Print Interior.pdf" \
       "book/The Art of Building AI Coding Agents.pdf" \
       book/interior-print.html
# NEVER git add -A — Astra's edition-2 files will contaminate the commit

# Website repo (commit as Bill Cox for Vercel)
cd ~/projects/coderhapsody
git add web/website/public/the-art-of-building-ai-coding-agents.pdf
git -c user.name='Bill Cox' -c user.email='waywardgeest@gmail.com' commit -m "..."
```

Bill pushes both repos.

## Verification

```bash
# Print PDF: page size, page count, renderer
pdfinfo "book/The Art of Building AI Coding Agents - Print Interior.pdf" \
  | grep -E "Pages|Page size|Producer"
# Expect: 6×9 inches (432 × 648 pts), HeadlessChrome, ~500+ pages

# EPUB: nav not in spine, correct entry count
unzip -p "book/The Art of Building AI Coding Agents.epub" EPUB/content.opf \
  | grep -o 'idref="[^"]*"' | grep nav
# Expect: no output (nav NOT in spine)

# Website PDF: cover page + interior
pdfinfo "book/The Art of Building AI Coding Agents.pdf" \
  | grep -E "Pages|Page size"
# Expect: 1 more page than the interior

# Compare with previous version for Amazon's 10% rule
# Old versions can be extracted from the website repo:
#   git show <hash>:web/website/public/the-art-of-building-ai-coding-agents.pdf > old.pdf
```

## What Does NOT Work

- **DOCX for ebooks** — KDP's conversion breaks TOC and navigation
- **DOCX for print** — KDP reads pandoc DOCX as 8.5×11 regardless of settings
- **`--css=` with pandoc** — CSS not embedded; use `-H file.html` instead
- **`--no-pdf-header-footer` with page numbers** — suppresses CSS `@page` margin boxes
- **`--toc` for EPUB** — inserts nav into spine, creating a duplicate contents page
