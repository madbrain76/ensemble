#!/usr/bin/env python3
"""Local text-only checks for the bounded Chapter 6 author pass."""
from pathlib import Path
import difflib
import hashlib
import json
import re
import sys
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[3]
EVIDENCE = Path(__file__).resolve().parent
CHAPTER = ROOT / 'book/edition-3/chapter-06.md'
BASELINE = EVIDENCE / 'baseline/chapter-06.md'
EXPECTED_BASELINE = 'ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55'

def digest(data):
    return hashlib.sha256(data).hexdigest()

old_bytes, new_bytes = BASELINE.read_bytes(), CHAPTER.read_bytes()
old, new = old_bytes.decode('utf-8'), new_bytes.decode('utf-8')
errors = []
if digest(old_bytes) != EXPECTED_BASELINE:
    errors.append('Baseline SHA mismatch')
if b'\x00' in new_bytes:
    errors.append('NUL byte in chapter')
for number, line in enumerate(new.splitlines(), 1):
    if line.rstrip(' \t') != line:
        errors.append(f'Trailing whitespace at chapter line {number}')
if not new.endswith('\n'):
    errors.append('Chapter lacks final newline')

# Resolve baseline links in their original chapter context, not in the evidence directory.
def links(text):
    results = []
    for match in re.finditer(r'!?\[[^\]\n]+\]\(([^)\s]+)\)', text):
        target = match.group(1)
        parsed = urlsplit(target)
        if parsed.scheme or parsed.netloc:
            results.append({'target': target, 'check': 'external-not-requested'})
            continue
        path = (CHAPTER.parent / unquote(parsed.path)).resolve() if parsed.path else CHAPTER
        exists = path.exists()
        item = {'target': target, 'resolved': str(path.relative_to(ROOT)), 'exists': exists}
        if parsed.fragment and exists:
            headings = re.findall(r'^#{1,6}\s+(.+?)\s*#*$', path.read_text(), re.M)
            slugs = [re.sub(r'[^\w\- ]', '', h.lower()).replace(' ', '-') for h in headings]
            item['fragment_exists'] = unquote(parsed.fragment) in slugs
        results.append(item)
    return results

baseline_links, final_links = links(old), links(new)
for item in final_links:
    if item.get('exists') is False or item.get('fragment_exists') is False:
        errors.append('Unresolved final local link: ' + item['target'])
(EVIDENCE / 'local-links.json').write_text(json.dumps({'method': 'Inline Markdown link targets resolved from original chapter directory; no network access or source implementation reads.', 'baseline': baseline_links, 'final': final_links}, indent=2) + '\n')

changes = []
for tag, i, j, k, l in difflib.SequenceMatcher(None, old.splitlines(), new.splitlines()).get_opcodes():
    if tag == 'equal':
        continue
    changes.append({'type': tag, 'baseline_lines': [i+1, j], 'final_lines': [k+1, l]})
    # Three teaching locations only: grading identity, Gemini accumulation, JSON completion.
    if not any(lo <= i+1 <= hi and j <= hi for lo, hi in [(181, 195), (440, 442), (652, 653)]):
        errors.append('Diff outside the three intended locations: ' + str(changes[-1]))

(EVIDENCE / 'final').mkdir(exist_ok=True)
(EVIDENCE / 'final/chapter-06.md').write_bytes(new_bytes)
(EVIDENCE / 'chapter.diff').write_text(''.join(difflib.unified_diff(old.splitlines(True), new.splitlines(True), fromfile='baseline/chapter-06.md', tofile='final/chapter-06.md')))

protected = []
for item in json.loads((EVIDENCE / 'protected-before.json').read_text()):
    data = (ROOT / item['path']).read_bytes()
    protected.append({'path': item['path'], 'before_sha256': item['sha256'], 'after_sha256': digest(data), 'unchanged': item['sha256'] == digest(data)})
(EVIDENCE / 'protected-comparison.json').write_text(json.dumps(protected, indent=2) + '\n')

summary = {'baseline': {'bytes': len(old_bytes), 'sha256': digest(old_bytes)}, 'final': {'bytes': len(new_bytes), 'sha256': digest(new_bytes)}, 'changed_locations': changes, 'local_links': len(final_links), 'protected_changes': [p['path'] for p in protected if not p['unchanged']], 'errors': errors}
print(json.dumps(summary, indent=2))
sys.exit(bool(errors))
