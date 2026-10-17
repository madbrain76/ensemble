#!/usr/bin/env python3
"""Freeze revision-1 text receipts and a complete, self-excluding root index."""
from pathlib import Path
import hashlib
import json

V = Path(__file__).resolve().parent
E = V.parent
R = E.parents[2]
I = E / 'initial-submission'
h = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
chapter = R / 'book/edition-3/chapter-06.md'
report = R / 'docs/edition-3-notes/ch06-author.md'
manifest = E / 'SHA256SUMS'
assert h(chapter) == '12aced5a996a34c4b384b596b6bd63f6137c89c9f4ad848bf272f6c13689004f'
assert h(manifest) == 'bdae2167a55ff3382b3dedcd2ac6e1070e7c518f84c87fd106e04fe1c61da71b', 'Do not rerun after freeze'
assert manifest.read_bytes() == (I / 'docs/edition-3-notes/ch06-author-evidence/SHA256SUMS').read_bytes()
assert report.read_bytes().startswith((I / 'docs/edition-3-notes/ch06-author.md').read_bytes())
validation = json.loads((V / 'validation.json').read_text())
assert validation['result'] == 'PASS'
assert validation['chapter_sha256'] == h(chapter)
assert h(V / 'final/chapter-06.md') == h(chapter)
for name in ['baseline-prose-lint', 'final-prose-lint', 'chapter-diff-check', 'chapter-git-diff', 'validation']:
    assert (V / (name + '.exit')).read_text().strip() == '0'
assert (V / 'validation-attempt1.exit').read_text().strip() == '1'
for line in (I / 'docs/edition-3-notes/ch06-author-evidence/SHA256SUMS').read_text().splitlines():
    digest, name = line.split('  ', 1)
    p = I / name if name in ['book/edition-3/chapter-06.md', 'docs/edition-3-notes/ch06-author.md'] else R / name
    assert h(p) == digest, name
for p in list(E.rglob('*')) + [chapter, report]:
    if not p.is_file():
        continue
    b = p.read_bytes()
    b.decode('utf-8')
    assert b'\x00' not in b, str(p)
assert all(line.rstrip(' \t') == line for line in report.read_text().splitlines())
submission = f'''# Chapter 6 author revision 1 — saved submission

One authorized revision of at most one, limited to the Taking it for a spin recipe.
The actual CLI package is selected from the core root; ./cmd/ch01 is an example,
while . applies only to a root CLI main package. Both chat invocations survive.
No chapter bytes outside the exact required recipe replacement changed.

- Chapter: {chapter.stat().st_size} bytes; SHA-256 `{h(chapter)}`.
- Appended author report SHA-256: `{h(report)}`.
- Revision-only portable diff SHA-256: `{h(V / 'chapter.diff')}`.
- Formal review SHA-256: `dfbe0eedfd38d2ead897bb7db5ae408aab80cfcd7b73a58ac4c9c542c0f36e0a`.
- Original chapter/report/index preserved under initial-submission with mapping.
- Original 31 historical entries verify with mapped chapter/report snapshots;
  old receipts remain unchanged. The old index is not current-root verification.
- Same-basename sequential lint exits 0/0; 5515 to 5537 words; hard checks pass.
  Existing soft warnings remain: negation 19/19, person gap 5087/5109.
- Final revision validator exit 0; exact single replacement, seven links,
  whitespace, append-only report, seven pin/status checks, 58 protected files.
  Only the parent-owned ch06-execution.md differs in the protected comparison.
- Initial revision validator exit 1 from an incorrect one-occurrence assumption
  for the normal and opt-out chat examples; failed script/output/exit retained.
  No manuscript change was made to satisfy that validator correction.
- Diff whitespace and diff capture exits 0. Safe UTF-8 evidence only.

Root SHA256SUMS is regenerated after preserving the original index. Its exact
set is every author-evidence file except that current root index itself, plus
the current chapter/report. It includes the saved historical index. The current
index digest is returned in the tool/submission receipt, not embedded in this
indexed file, avoiding a hash cycle. The final checksum command verifies it.

No student build/source edit, new code review, child, paid/key access, GUI,
publication, staging, commit, tag, push, cleanup, or out-of-scope edit.
All jobs finish before final submission. Root verifies the sole review condition
and closes; no second formal review or further revision is requested.
'''
(V / 'submission.md').write_text(submission)
freeze = {'status': 'PASS', 'chapter_sha256': h(chapter), 'chapter_bytes': chapter.stat().st_size, 'report_sha256': h(report), 'portable_revision_diff_sha256': h(V / 'chapter.diff'), 'saved_submission_sha256': h(V / 'submission.md'), 'initial_history_preserved': True, 'report_initial_prefix_preserved': True, 'safe_utf8_no_nul': True, 'current_inventory_rule': 'Every author-evidence file except the current root SHA256SUMS itself, plus external current chapter and author report. Includes the initial historical SHA256SUMS. No self hash.'}
(V / 'freeze.json').write_text(json.dumps(freeze, indent=2) + '\n')
(V / 'freeze.exit').write_text('0\n')
files = sorted([p for p in E.rglob('*') if p.is_file() and p != manifest] + [chapter, report], key=lambda p: str(p.relative_to(R)))
manifest.write_text(''.join(f'{h(p)}  {p.relative_to(R)}\n' for p in files))
rows = [line.split('  ', 1) for line in manifest.read_text().splitlines()]
assert len(rows) == len(files)
assert {name for _, name in rows} == {str(p.relative_to(R)) for p in files}
assert len({name for _, name in rows}) == len(rows)
for digest, name in rows:
    p = R / name
    assert h(p) == digest
    assert b'\x00' not in p.read_bytes()
    p.read_bytes().decode('utf-8')
print(json.dumps(dict(freeze, inventory_entries=len(rows), inventory_sha256=h(manifest), inventory_exact_set_and_hashes='PASS'), indent=2))
