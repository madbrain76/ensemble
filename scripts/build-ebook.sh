#!/bin/bash
# Build the KDP ebook (EPUB) from the assembled manuscript.
# Usage: scripts/build-ebook.sh [cover-image]
#
# WHY THIS DIFFERS FROM THE PRINT BUILD
#
# The print interior and the ebook are built from the SAME manuscript but
# want opposite things from a table of contents, so they are NOT one recipe
# with a flag. Print is a fixed page count: a reader cannot tap a heading,
# so the manuscript's own "# Contents" section with its page-ordered
# chapter and section list is the only navigation that exists, and it must
# stay. The ebook is reflowable and Kindle can build its own navigation
# from the EPUB nav document, but that nav is a MENU, not a page in the
# reading flow — and it is not guaranteed to be surfaced. So the ebook
# keeps a visible contents page too.
#
# Shipping both produced a real defect: the published ebook showed the
# table of contents TWICE. The two copies came from different places, and
# an earlier version of this script removed BOTH of them. That overshot:
# the published ebook then had no visible contents page at all. Amazon did
# not reliably surface one of its own from the nav document, so readers got
# nothing. Exactly one of the two copies has to come back.
#
#   1. The manuscript's own "# Contents" section is KEPT. It is generated
#      by scripts/book-toc.py, it lists chapters and sections, and it is
#      the copy we control. This is the visible contents page.
#
#   2. Pandoc's --toc flag is NOT passed. For EPUB, --toc does not merely
#      "add a table of contents": it inserts the nav document INTO THE
#      SPINE, which is what makes it render as a visible page. Omitting it
#      leaves nav.xhtml and toc.ncx fully populated for Kindle while
#      keeping them out of the reading order.
#
# --toc-depth=1 IS still passed, and dropping it is the trap here. It
# works independently of --toc, and without it pandoc defaults to depth 3,
# so every "##" section lands in the nav document. Measured on this book:
# 30 nav entries with the flag, 426 without it. Amazon would build its
# table of contents from all 426.
#
# Verify a build with:
#   unzip -p "<epub>" EPUB/content.opf | grep -o 'idref="[^"]*"' | head
# "nav" must NOT appear in that spine listing.
set -euo pipefail

BOOK_DIR="$(cd "$(dirname "$0")/../book" && pwd)"
MANUSCRIPT="$BOOK_DIR/the-art-of-building-ai-coding-agents.md"
FRONT="$BOOK_DIR/front-matter.md"
OUT="$BOOK_DIR/The Art of Building AI Coding Agents.epub"

TITLE="The Art of Building AI Coding Agents"
AUTHOR="Bill Cox"

COVER="${1:-$HOME/Downloads/The Art of Building AI Coding Agents (1)_upscayl_4x_digital-art-4x.png}"

for f in "$MANUSCRIPT" "$FRONT"; do
    if [ ! -f "$f" ]; then
        echo "build-ebook: missing required input: $f" >&2
        exit 1
    fi
done

if [ ! -f "$COVER" ]; then
    echo "build-ebook: cover image not found: $COVER" >&2
    echo "build-ebook: pass one explicitly: scripts/build-ebook.sh /path/to/cover.png" >&2
    exit 1
fi

# The ebook is built from the manuscript as-is, "# Contents" section and
# all. Stripping that section is exactly what left the published ebook
# with no visible contents page, so assert the opposite now: refuse to
# build if it has gone missing upstream. scripts/book-toc.py regenerates
# it on every assemble, so absence means the assemble step was skipped.
if ! grep -q '^# Contents' "$MANUSCRIPT"; then
    echo "build-ebook: no '# Contents' section in the manuscript. The ebook" >&2
    echo "build-ebook: needs it as its visible contents page. Did the" >&2
    echo "build-ebook: assemble step skip scripts/book-toc.py?" >&2
    exit 1
fi

pandoc "$FRONT" "$MANUSCRIPT" \
    -o "$OUT" \
    --from=markdown --to=epub \
    --toc-depth=1 \
    --resource-path="$BOOK_DIR" \
    -M title="$TITLE" \
    -M author="$AUTHOR" \
    --epub-cover-image="$COVER"

# The defect this script exists to prevent is invisible in the output file
# size, so assert the structure instead of trusting the build.
spine=$(unzip -p "$OUT" EPUB/content.opf | grep -o 'idref="[^"]*"')
if echo "$spine" | grep -q '"nav"'; then
    echo "build-ebook: nav document is in the spine, so the ebook has a" >&2
    echo "build-ebook: redundant visible contents page. Was --toc added?" >&2
    exit 1
fi

navlinks=$(unzip -p "$OUT" EPUB/nav.xhtml | grep -o '<a href="[^"]*"' | wc -l | tr -d ' ')
if [ "$navlinks" -gt 60 ]; then
    echo "build-ebook: nav document has $navlinks entries, which means it" >&2
    echo "build-ebook: descended past chapter level. Was --toc-depth lost?" >&2
    exit 1
fi

# The contents page must be a real page in the READING FLOW, not just the
# nav menu. The nav document is a menu Amazon may or may not surface; a
# missing visible page is what readers actually complained about, and it
# is invisible to every other check in this script.
#
# Count rather than grep -q: under `set -o pipefail`, grep -q closes the
# pipe on its first match, unzip dies of SIGPIPE, and the pipeline reports
# 141 even though the page IS there. That false alarm cost a build.
contents_pages=$(unzip -p "$OUT" 'EPUB/text/*.xhtml' \
    | grep -cE '<h1[^>]*>[[:space:]]*Contents' || true)
if [ "$contents_pages" -lt 1 ]; then
    echo "build-ebook: no visible Contents page in the ebook text. The" >&2
    echo "build-ebook: manuscript's '# Contents' section did not survive" >&2
    echo "build-ebook: into the spine. Readers get no table of contents." >&2
    exit 1
fi

# A dropped illustration is invisible in the checks above and nearly
# invisible in the file size next to a 16MB cover, but pandoc only WARNS
# about one. Relative image paths resolve against --resource-path, so
# losing that flag silently ships a book with no figures.
refs=$(cat "$FRONT" "$MANUSCRIPT" | grep -o '!\[[^]]*\]([^)]*)' | wc -l | tr -d ' ')
embedded=$(unzip -l "$OUT" | grep -ciE '\.(png|jpg|jpeg|gif)$' || true)
if [ "$embedded" -ne "$((refs + 1))" ]; then
    echo "build-ebook: $refs images referenced plus 1 cover, but $embedded" >&2
    echo "build-ebook: embedded. Figures were dropped. Check --resource-path." >&2
    exit 1
fi

echo "build-ebook: wrote $OUT"
echo "build-ebook: $navlinks chapter entries in the nav document, nav not in spine"
echo "build-ebook: $embedded images embedded ($refs figures plus cover)"
