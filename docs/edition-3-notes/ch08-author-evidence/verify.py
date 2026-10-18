#!/usr/bin/env python3
"""Read-only author checks; receipts are written only into this evidence directory."""
import datetime
import difflib
import hashlib
import json
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
EVIDENCE = Path(__file__).resolve().parent
BASELINE = Path('/tmp/cadence-ch08-text-phase-baseline/baseline.json')
CHAPTER = 'book/edition-3/chapter-08.md'
# Publication ownership is root's: do not open/hash/assert preservation here.
EXCLUDED_PUBLICATION = {
    'README.md', 'book/README.md',
    'book/The Art of Building AI Coding Agents - Print Interior.pdf',
    'book/The Art of Building AI Coding Agents.epub',
    'book/The Art of Building AI Coding Agents.pdf',
    'book/build-print.py', 'book/front-matter-print.md', 'book/front-matter.md',
    'book/interior-print.html', 'book/note-from-bill.md', 'book/print-header.html',
    'book/the-art-of-building-ai-coding-agents.md', 'docs/edition-2-publication.md',
}


def identity(path):
    data = path.read_bytes()
    return {'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest()}


baseline = json.loads(BASELINE.read_text())
checks = {'checked_utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
          'root_baseline': str(BASELINE), 'root_baseline_identity': identity(BASELINE)}
for group in ('inputs', 'protected_tracked'):
    mismatches = []
    checked_count = 0
    for name, expected in baseline[group].items():
        if group == 'inputs' and name == CHAPTER:
            continue
        if name in EXCLUDED_PUBLICATION:
            continue
        checked_count += 1
        actual = identity(ROOT / name)
        if actual != expected:
            mismatches.append({'path': name, 'expected': expected, 'actual': actual})
    checks[group] = {'checked_count': checked_count,
                     'mismatches': mismatches}

checks['excluded_publication_paths'] = sorted(EXCLUDED_PUBLICATION)
checks['publication_qualification'] = 'No publication preservation assertion; earlier broad hash-check receipts were out of scope and are retained as process failures.'
status = subprocess.run(['git', '--no-pager', 'status', '--porcelain=v1', '-z', '--untracked-files=no'], cwd=ROOT, stdout=subprocess.PIPE, check=True).stdout
before_status = set(bytes.fromhex(baseline['status_hex']).split(b'\0')) - {b''}
after_status = set(status.split(b'\0')) - {b''}
checks['tracked_status_delta'] = {
    'added': sorted(s.decode() for s in after_status - before_status),
    'removed': sorted(s.decode() for s in before_status - after_status),
    'only_chapter_added': after_status - before_status == {b' M book/edition-3/chapter-08.md'} and not (before_status - after_status),
    'qualification': 'Status-name comparison only; not a publication content/preservation assertion.',
}
before = (EVIDENCE / 'baseline/chapter-08.md').read_text()
after = (ROOT / CHAPTER).read_text()
old_start = 'From the repository root, grade the core and its optional GUI separately:'
new_start = 'Build your actual core CLI target.'
end = '| check | pts | what it tests |'
checks['bounded_edit'] = {
    'prefix_unchanged': before.split(old_start)[0] == after.split(new_start)[0],
    'suffix_unchanged': before.split(end, 1)[1] == after.split(end, 1)[1],
    'before': identity(EVIDENCE / 'baseline/chapter-08.md'),
    'after': identity(ROOT / CHAPTER),
    'saved_final_matches_chapter': identity(EVIDENCE / 'final/chapter-08.md') == identity(ROOT / CHAPTER),
    'saved_baseline_matches_root': identity(EVIDENCE / 'baseline/chapter-08.md') == baseline['inputs'][CHAPTER],
}
(EVIDENCE / 'chapter-08.diff').write_text(''.join(difflib.unified_diff(
    before.splitlines(True), after.splitlines(True),
    fromfile='baseline/book/edition-3/chapter-08.md', tofile=CHAPTER)))

links = []
for target in re.findall(r'\[[^\]]*\]\(([^)]+)\)', after):
    if '://' in target or target.startswith('#'):
        continue
    links.append({'target': target, 'exists': (ROOT / CHAPTER).parent.joinpath(target.split('#')[0]).exists()})
checks['relative_links'] = {'check': 'target existence only; no linked content opened',
                            'links': links,
                            'unchanged': links == json.loads((EVIDENCE / 'baseline-links.json').read_text())['links']}
checks['commands'] = []
commands = [
    ('final-lint', ['go', 'run', './cmd/lintprose', str(EVIDENCE / 'final/chapter-08.md')]),
    ('final-diff-check', ['git', 'diff', '--check']),
    ('tracked-status', ['git', '--no-pager', 'status', '--short', '--untracked-files=no']),
    ('scoped-status', ['git', '--no-pager', 'status', '--short', '--untracked-files=all', '--', CHAPTER,
                       'docs/edition-3-notes/ch08-author.md', 'docs/edition-3-notes/ch08-author-evidence']),
]
for name, argv in commands:
    proc = subprocess.run(argv, cwd=ROOT, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
    (EVIDENCE / (name + '.log')).write_text(proc.stdout)
    checks['commands'].append({'name': name, 'argv': argv, 'cwd': str(ROOT),
                                'exit_code': proc.returncode, 'output': name + '.log'})
    print(name, 'exit', proc.returncode)
    print(proc.stdout)
checks['passed'] = (
    all(not checks[g]['mismatches'] for g in ('inputs', 'protected_tracked'))
    and checks['tracked_status_delta']['only_chapter_added']
    and checks['bounded_edit']['saved_final_matches_chapter']
    and checks['bounded_edit']['prefix_unchanged']
    and checks['bounded_edit']['suffix_unchanged']
    and checks['bounded_edit']['saved_baseline_matches_root']
    and checks['relative_links']['unchanged']
    and all(link['exists'] for link in links)
    and all(cmd['exit_code'] == 0 for cmd in checks['commands'])
)
(EVIDENCE / 'final-checks.json').write_text(json.dumps(checks, indent=2) + '\n')
print('Mechanical checks passed:', checks['passed'])
raise SystemExit(0 if checks['passed'] else 1)
