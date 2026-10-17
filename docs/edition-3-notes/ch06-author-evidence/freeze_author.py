#!/usr/bin/env python3
"""Freeze text-only author evidence; SHA256SUMS necessarily excludes itself."""
from pathlib import Path
import hashlib
import json
import re
import subprocess

E = Path(__file__).resolve().parent
R = E.parents[2]
chapter = R / 'book/edition-3/chapter-06.md'
report = R / 'docs/edition-3-notes/ch06-author.md'
expected = '864d806c63c019ede50810ae5618b6fa4545f7028ae3abcb060722458babb4e9'
hash_of = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
assert hash_of(chapter) == expected
assert chapter.read_bytes() == (E / 'final/chapter-06.md').read_bytes()
assert subprocess.check_output(['git', '--no-pager', 'rev-parse', 'HEAD'], cwd=R).decode().strip() == 'b67386ab80cf2afe7985db5a32ab99ec2311f090'
checked = []
for p in sorted(E.rglob('*')) + [report, chapter]:
    if not p.is_file() or p.name == 'SHA256SUMS':
        continue
    data = p.read_bytes()
    data.decode('utf-8')
    assert b'\x00' not in data, str(p)
    checked.append(str(p.relative_to(R)))
for p in [report, chapter]:
    s = p.read_text()
    assert s.endswith('\n')
    assert all(line.rstrip(' \t') == line for line in s.splitlines()), str(p)
for target in re.findall(r'\[[^\]\n]+\]\(([^)\s]+)\)', report.read_text()):
    assert (report.parent / target).resolve().exists(), target
for name, value in [('baseline-prose-lint.exit', '0'), ('final-prose-lint.exit', '0'), ('text-validation.exit', '0'), ('chapter-diff-check.exit', '0'), ('text-validation-attempt1.exit', '1')]:
    assert (E / name).read_text().strip() == value, name
result = {'status': 'PASS', 'chapter_sha256': hash_of(chapter), 'report_sha256': hash_of(report), 'utf8_payloads_checked_before_receipt_generation': len(checked), 'report_and_chapter_whitespace': 'PASS', 'report_local_links': 'PASS', 'head_unchanged': True, 'manifest_rule': 'All files in this author evidence directory, plus report and chapter; SHA256SUMS alone is excluded to avoid self-reference. No source/build/provider validation implied.'}
(E / 'freeze-check.txt').write_text(json.dumps(result, indent=2) + '\n')
(E / 'freeze-check.exit').write_text('0\n')
files = sorted([p for p in E.rglob('*') if p.is_file() and p.name != 'SHA256SUMS'] + [report, chapter], key=lambda p: str(p.relative_to(R)))
(E / 'SHA256SUMS').write_text(''.join(f'{hash_of(p)}  {p.relative_to(R)}\n' for p in files))
print(json.dumps(dict(result, manifest_entries=len(files), inventory_sha256=hash_of(E / 'SHA256SUMS'), portable_diff_sha256=hash_of(E / 'chapter.diff')), indent=2))
