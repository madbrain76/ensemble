#!/usr/bin/env python3
"""Two-pass print build: render, learn real page numbers, inject them into the TOC.

Chrome renders @page margin boxes (so footer page numbers work) but does NOT
support target-counter(), so a TOC cannot learn its own page references in one
pass. This renders once, reads where each chapter actually landed, injects those
numbers, renders again, then re-reads to prove pagination did not shift.
"""
import re, subprocess, sys
from collections import Counter
from pathlib import Path

import pymupdf

CHROME = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
BOOK = Path(__file__).resolve().parent


def render(html: Path, pdf: Path) -> None:
    pdf.unlink(missing_ok=True)
    subprocess.run(
        [CHROME, "--headless=new", "--no-sandbox", f"--print-to-pdf={pdf}",
         "--no-pdf-header-footer", f"file://{html}"],
        capture_output=True, timeout=600, check=True,
    )
    if not pdf.exists():
        sys.exit(f"FATAL: Chrome produced no PDF at {pdf}")


def page_texts(pdf: Path) -> list[str]:
    out = subprocess.run(["pdftotext", str(pdf), "-"],
                         capture_output=True, text=True, timeout=300, check=True).stdout
    return out.split("\f")


def check_type_size(pdf: Path) -> None:
    """Check printed type, since correct CSS does not prevent Chrome shrinking it.

    Count characters rather than spans so short headings and page numbers
    cannot outweigh body prose. Code must retain its intended size everywhere.
    Small rounding differences in Chrome's PDF coordinates are harmless.
    """
    body = Counter()
    code_sizes = []
    with pymupdf.open(pdf) as document:
        for page in document:
            for block in page.get_text("dict")["blocks"]:
                for line in block.get("lines", []):
                    for span in line["spans"]:
                        if not span["text"].strip():
                            continue
                        size = round(span["size"], 2)
                        if span["font"] == "Georgia":
                            body[size] += len(span["text"])
                        if span["font"].startswith("Menlo"):
                            code_sizes.append(size)
    if not body or not code_sizes:
        sys.exit("FATAL: could not measure the expected body and code fonts")
    body_size = body.most_common(1)[0][0]
    code_size = min(code_sizes)
    if abs(body_size - 12) > 0.1 or code_size < 9.9:
        sys.exit(f"FATAL: printed type is {body_size}pt body / {code_size}pt code; "
                 "expected 12pt / 10pt. Check overflowing code and tables.")
    print(f"  printed type: {body_size}pt body, minimum {code_size}pt code")


def first_line(page: str) -> str:
    """Opening text of a page, normalized.

    A heading wraps across several printed lines, and pandoc wraps the TOC
    title at a different point, so compare normalized prefixes rather than
    single lines.
    """
    lines = [l.strip() for l in page.splitlines() if l.strip()]
    return norm(" ".join(lines[:4]))


def norm(s: str) -> str:
    return re.sub(r"\s+", " ", s).strip()


def toc_entries(html: str) -> list[tuple[str, str]]:
    """Ordered (href_id, title) pairs from pandoc's #TOC nav.

    The anchor pattern must tolerate '<a\\nhref=' because pandoc wraps long
    tags. Assuming a single space here silently dropped the three longest
    chapters, and because the completeness check counted only what this
    function returned, it agreed that nothing was missing.
    """
    nav = re.search(r'<nav id="TOC".*?</nav>', html, re.S)
    if not nav:
        sys.exit("FATAL: no <nav id=\"TOC\"> found. Did pandoc run with --toc?")
    block = nav.group(0)
    found = re.findall(r'<a\s[^>]*href="#([^"]+)"[^>]*>(.*?)</a>', block, re.S)
    want = len(re.findall(r"<li>", block))
    if len(found) != want:
        sys.exit(f"FATAL: {want} TOC items but only {len(found)} anchors parsed")
    return found


def strip_tags(s: str) -> str:
    return norm(re.sub(r"<[^>]+>", "", s))


def build_map(pdf: Path, titles: list[str]) -> dict[str, int]:
    """Walk pages in order, matching each expected title to the page it opens.

    Chapters are ordered and each starts a page (h1 page-break-before: always),
    so an ordered walk avoids matching a title where it appears inside the TOC.
    """
    pages = page_texts(pdf)
    found, idx = {}, 0
    for pageno, text in enumerate(pages, start=1):
        if idx >= len(titles):
            break
        want = titles[idx]
        if first_line(text).startswith(want):
            found[want] = pageno
            idx += 1
    return found


def inject(html: str, entries: list[tuple[str, str]], mapping: dict[str, int]) -> str:
    nav = re.search(r'<nav id="TOC".*?</nav>', html, re.S)
    block = nav.group(0)
    new = block
    done = 0
    for href, raw in entries:
        title = strip_tags(raw)
        if title not in mapping:
            continue
        # Pandoc wraps long tags, so the anchor can read '<a\nhref="#id"'.
        # Matching a literal '<a href=' silently skipped every wrapped entry
        # and still produced a plausible looking TOC.
        pat = re.compile(r'<a\s[^>]*href="#' + re.escape(href) + r'"[^>]*>.*?</a>', re.S)
        m = pat.search(new)
        if not m:
            sys.exit(f"FATAL: no anchor found for #{href}")
        end = m.end()
        tail = new[end:end + 80]
        prev = re.match(r'\s*<span class="tocpg">\d+</span>', tail)
        if prev:  # drop the number a previous pass wrote
            end += prev.end()
        new = new[:end] + f'<span class="tocpg">{mapping[title]}</span>' + new[end:]
        done += 1
    if done != len(entries):
        sys.exit(f"FATAL: injected {done} numbers for {len(entries)} TOC entries")
    return html.replace(block, new)


def clean_html(html: str) -> str:
    """Fix pandoc output for print before the two-pass TOC build.

    1. Force a page break after the pandoc TOC so it doesn't flow into the
       front matter.
    2. Remove the self-referential "Contents" entry from the pandoc TOC
       (the markdown TOC heading appears as a TOC entry pointing at itself).
    3. Strip the markdown-generated "# Contents" section entirely — the
       pandoc --toc generates the real TOC with page numbers; the markdown
       one is for GitHub and has only anchor links (useless in print).
    """
    # 1. Page break after TOC nav
    html = html.replace(
        '</nav>',
        '</nav>\n<div style="page-break-after: always;"></div>',
        1,  # only the first </nav> (the TOC)
    )

    # 2. Remove "Contents" entry from pandoc TOC (self-referential)
    html = re.sub(
        r'<li><a\s[^>]*href="#contents"[^>]*>Contents</a>.*?</li>\s*',
        '',
        html,
        count=1,
        flags=re.S,
    )

    # 3. Strip the markdown "# Contents" section from the body.
    #    It starts with <h1 id="contents"> and ends just before the next <h1>.
    html = re.sub(
        r'<h1 id="contents">Contents</h1>.*?(?=<h1[ >])',
        '',
        html,
        count=1,
        flags=re.S,
    )

    return html


def run_pandoc(src: Path) -> None:
    """Regenerate interior-print.html from the assembled manuscript.

    This used to be a manual step documented in book/README.md, and skipping
    it failed SILENTLY: build-print.py would re-render a stale HTML left over
    from an earlier chapter, producing a freshly-timestamped PDF with an old
    TOC and the newest chapter missing entirely. Nothing in the two-pass build
    noticed, because every heading the stale TOC listed did match a page.
    Running pandoc here makes that class of staleness impossible.
    """
    manuscript = BOOK / "the-art-of-building-ai-coding-agents.md"
    front = BOOK / "front-matter-print.md"
    header = BOOK / "print-header.html"
    for f in (front, manuscript, header):
        if not f.exists():
            sys.exit(f"FATAL: missing build input {f}")
    cmd = [
        "pandoc", front.name, manuscript.name,
        "-o", src.name,
        "--from=markdown", "--to=html", "--standalone",
        "--toc", "--toc-depth=1",
        "-M", "title=The Art of Building AI Coding Agents",
        "-M", "author=Bill Cox",
        "-H", header.name,
        "--lua-filter", str(BOOK.parent / "scripts" / "publication-links.lua"),
    ]
    print("pass 0: pandoc regenerate interior-print.html")
    r = subprocess.run(cmd, cwd=BOOK, capture_output=True, text=True)
    if r.returncode != 0:
        sys.exit(f"FATAL: pandoc failed:\n{r.stderr}")


def build_reading_pdf(interior_path: Path) -> None:
    """Prepend the cover without losing Chrome's named contents destinations.

    Concatenation can retain link annotations while discarding their name tree.
    Resolve those names in the original PDF, then offset explicit page targets.
    resolve_link also converts PDF coordinates to PyMuPDF page coordinates.
    """
    output = BOOK / "The Art of Building AI Coding Agents.pdf"
    with pymupdf.open(interior_path) as interior, pymupdf.open() as reading:
        rect = interior[0].rect
        cover = reading.new_page(width=rect.width, height=rect.height)
        cover.insert_image(cover.rect, filename=str(BOOK / "cover.png"),
                           keep_proportion=True)
        reading.insert_pdf(interior, links=False)
        count = 0
        for index, page in enumerate(interior):
            for original in page.get_links():
                link = {k: v for k, v in original.items()
                        if k not in ("xref", "id", "nameddest")}
                if "nameddest" in original:
                    target, x, y = interior.resolve_link("#" + original["nameddest"])
                    if target < 0:
                        sys.exit(f"FATAL: unresolved contents link {original['nameddest']}")
                    link.update(kind=pymupdf.LINK_GOTO, page=target + 1,
                                to=pymupdf.Point(x, y))
                elif link["kind"] == pymupdf.LINK_GOTO:
                    link["page"] += 1
                reading[index + 1].insert_link(link)
                count += 1
        reading.save(output, deflate=True)
    print(f"  reading PDF: cover plus interior, {count} links preserved")


def main() -> None:
    src = BOOK / "interior-print.html"
    pdf = BOOK / "The Art of Building AI Coding Agents - Print Interior.pdf"

    run_pandoc(src)
    html = src.read_text()

    print("cleaning HTML: TOC page break, strip duplicate Contents")
    html = clean_html(html)
    src.write_text(html)

    entries = toc_entries(html)
    titles = [strip_tags(t) for _, t in entries]
    print(f"TOC entries: {len(titles)}")

    print("pass 1: render")
    render(src, pdf)
    check_type_size(pdf)
    m1 = build_map(pdf, titles)
    print(f"  matched {len(m1)}/{len(titles)} headings")
    missing = [t for t in titles if t not in m1]
    if missing:
        sys.exit(f"FATAL: {len(missing)} headings never matched a page: {missing[:5]}")

    print("pass 2: inject + render")
    src.write_text(inject(html, entries, m1))
    render(src, pdf)

    m2 = build_map(pdf, titles)
    shifted = {t: (m1.get(t), m2.get(t)) for t in titles if m1.get(t) != m2.get(t)}
    if shifted:
        print(f"  pagination SHIFTED for {len(shifted)}; re-injecting with pass-2 map")
        src.write_text(inject(src.read_text(), entries, m2))
        render(src, pdf)
        m3 = build_map(pdf, titles)
        still = {t for t in titles if m2.get(t) != m3.get(t)}
        if still:
            sys.exit(f"FATAL: pagination still shifting for {len(still)} headings")
        print("  CONVERGED")
    else:
        print("  stable: pagination identical across passes")

    check_type_size(pdf)
    n = subprocess.run(["pdfinfo", str(pdf)], capture_output=True, text=True).stdout
    print("\n" + "\n".join(l for l in n.splitlines() if l.startswith(("Pages", "Page size"))))
    build_reading_pdf(pdf)


if __name__ == "__main__":
    main()
