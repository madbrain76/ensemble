#!/usr/bin/env python3
"""Bounded text-revision checks; never build or execute student code."""
from pathlib import Path
import difflib
import hashlib
import json
import re
import subprocess

V = Path(__file__).resolve().parent
E = V.parent
R = E.parents[2]
I = E / 'initial-submission'
h = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
def save(name, value):
    (V / name).write_text(json.dumps(value, indent=2) + '\n')
chapter = R / 'book/edition-3/chapter-06.md'
report = R / 'docs/edition-3-notes/ch06-author.md'
before = (V / 'baseline/chapter-06.md').read_text()
after = chapter.read_text()
assert h(V / 'baseline/chapter-06.md') == '864d806c63c019ede50810ae5618b6fa4545f7028ae3abcb060722458babb4e9'
assert h(chapter) == '12aced5a996a34c4b384b596b6bd63f6137c89c9f4ad848bf272f6c13689004f'
assert chapter.read_bytes() == (V / 'final/chapter-06.md').read_bytes()
replacement = json.loads((V / 'required-replacement.json').read_text())
assert before.count(replacement['old']) == 1
assert after == before.replace(replacement['old'], replacement['new'])
assert before.count('/tmp/ch6agent chat') == after.count('/tmp/ch6agent chat') == 1
assert "use `.` only if the CLI's main\npackage lives at the module root" in after
assert 'go build -o /tmp/ch6agent ./cmd/ch01' in after
for p in [chapter, report]:
    text = p.read_text()
    assert text.endswith('\n')
    assert all(line.rstrip(' \t') == line for line in text.splitlines()), str(p)
assert report.read_bytes().startswith((I / 'docs/edition-3-notes/ch06-author.md').read_bytes())
assert len(report.read_bytes()) > len((I / 'docs/edition-3-notes/ch06-author.md').read_bytes())
(V / 'chapter.diff').write_text(''.join(difflib.unified_diff(before.splitlines(True), after.splitlines(True), fromfile='baseline/chapter-06.md', tofile='final/chapter-06.md')))
links = []
for target in re.findall(r'\[[^\]\n]+\]\(([^)\s]+)\)', after):
    if '://' in target or target.startswith('mailto:'):
        continue
    path, _, fragment = target.partition('#')
    assert not fragment, target
    resolved = (chapter.parent / path).resolve()
    assert resolved.exists(), target
    links.append({'target': target, 'exists': True})
assert len(links) == 7
save('local-links.json', links)
initial_manifest = I / 'docs/edition-3-notes/ch06-author-evidence/SHA256SUMS'
assert h(initial_manifest) == 'bdae2167a55ff3382b3dedcd2ac6e1070e7c518f84c87fd106e04fe1c61da71b'
external = {'book/edition-3/chapter-06.md', 'docs/edition-3-notes/ch06-author.md'}
historical = []
for line in initial_manifest.read_text().splitlines():
    expected, original = line.split('  ', 1)
    mapped = (I / original) if original in external else (R / original)
    assert h(mapped) == expected, original
    historical.append({'original_path': original, 'verified_path': str(mapped.relative_to(R)), 'sha256': expected})
assert len(historical) == 31
save('initial-history-verification.json', {'entries': historical, 'all_match': True, 'root_external_paths_not_used_for_historical_chapter_or_report': True})
protected = []
allowed_parent = {'docs/edition-3-notes/ch06-execution.md', 'docs/edition-3-notes/ch06-accepted.json'}
for row in json.loads((V / 'protected-before.json').read_text()):
    current = h(R / row['path'])
    same = current == row['sha256']
    assert same or row['path'] in allowed_parent, row['path']
    protected.append(dict(row, current_sha256=current, unchanged=same, owner_if_changed='parent' if not same else None))
save('protected-comparison.json', protected)
pins = {}
for name, spec in json.loads((R / 'docs/edition-3-notes/ch06-text-review-evidence/pins.json').read_text()).items():
    proc = subprocess.run(spec['command'], cwd=spec['cwd'], text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    expected = spec.get('expected', '')
    assert proc.returncode == 0 and proc.stdout.strip() == expected, name
    pins[name] = {'command': spec['command'], 'cwd': spec['cwd'], 'exit': proc.returncode, 'stdout': proc.stdout, 'stderr': proc.stderr, 'expected': expected}
save('pins.json', pins)
for name in ['baseline-prose-lint.exit', 'final-prose-lint.exit', 'chapter-diff-check.exit', 'chapter-git-diff.exit']:
    assert (V / name).read_text().strip() == '0', name
(V / 'git-status-final.txt').write_bytes(subprocess.check_output(['git', '--no-pager', 'status', '--short'], cwd=R))
result = {'result': 'PASS', 'exact_single_required_recipe_replacement': True, 'all_other_chapter_bytes_unchanged': True, 'chapter_bytes': chapter.stat().st_size, 'chapter_sha256': h(chapter), 'report_append_only': True, 'local_links': len(links), 'historical_entries_verified_with_mapping': len(historical), 'protected_files': len(protected), 'parent_owned_changes': [x['path'] for x in protected if not x['unchanged']], 'read_only_pin_checks': len(pins), 'student_source_builds_or_behavioral_tests': 0}
save('validation.json', result)
print(json.dumps(result, indent=2))
