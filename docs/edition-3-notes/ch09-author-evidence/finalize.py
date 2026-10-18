#!/usr/bin/env python3
"""Seal only the bounded Chapter 9 author outputs after final checks."""
import datetime
import hashlib
import json
import pathlib
import shutil

root = pathlib.Path('/Users/bill/projects/ensemble')
evidence = root / 'docs/edition-3-notes/ch09-author-evidence'
report = root / 'docs/edition-3-notes/ch09-author-report.md'
chapter = root / 'book/edition-3/chapter-09.md'
assert not (evidence / 'seal.json').exists(), 'seal already exists'

def pin(path):
    data = path.read_bytes()
    return {'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest()}

baseline = json.loads((evidence / 'baseline-inputs.json').read_text())
inputs = {}
for name, before in baseline['inputs'].items():
    final = pin(root / name)
    inputs[name] = {'baseline': before, 'final': final, 'matches': final == before}
    if name == 'book/edition-3/chapter-09.md':
        assert final == {
            'bytes': 20767,
            'sha256': '137f7fe93628d6434e5be6266933637b4bb0362b1926a6fe1e6da061db16447d',
        }
    else:
        assert final == before, name
assert (root / 'book/voice.md').read_bytes() == (root / 'book/edition-3/voice.md').read_bytes()
assert shutil.disk_usage(root).free >= 1024**3

# Check authored files, not verbatim command output or copied source snapshots.
files = [chapter, report] + sorted(evidence.glob('*.py')) + sorted(evidence.glob('*.json'))
checked = []
for path in files:
    text = path.read_text()
    bad = [i for i, line in enumerate(text.splitlines(), 1) if line.rstrip() != line]
    assert not bad, (str(path), bad)
    assert not text or text.endswith('\n'), str(path)
    checked.append(str(path.relative_to(root)))
result = {
    'argv': ['python3', 'docs/edition-3-notes/ch09-author-evidence/finalize.py'],
    'cwd': str(root),
    'checked_utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
    'checks_passed': True,
    'non_chapter_pins_matching': 23,
    'chapter_expected_final_hash': True,
    'voice_byte_comparison': True,
    'authored_text_whitespace_and_final_newline_checked': checked,
    'raw_outputs_and_source_snapshots_not_normalized': True,
    'free_bytes': shutil.disk_usage(root).free,
}
(evidence / 'final-check.json').write_text(json.dumps(result, indent=2) + '\n')

artifacts = {}
for path in [chapter, report] + sorted(evidence.iterdir()):
    assert path.is_file(), str(path)
    artifacts[str(path.relative_to(root))] = pin(path)
seal = {
    'sealed_utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
    'task_sha256': '6eceb70e1746991310d6d763868d647c9aa21eb2ece1da87f0dcb44c0031ec37',
    'hash_algorithm': 'SHA-256',
    'rule': 'Final chapter, report and all evidence files are sealed. seal.json excludes itself; no self-hash and no modifications after sealing.',
    'artifacts': artifacts,
    'prescribed_inputs': inputs,
    'all_non_chapter_pins_unchanged': True,
    'formal_text_review': 'not performed by author; root owns one pass and maximum one revision',
    'unfinished_jobs': [],
    'children_spawned': 0,
}
with (evidence / 'seal.json').open('x') as stream:
    json.dump(seal, stream, indent=2)
    stream.write('\n')
print(json.dumps({'checks_passed': True, 'artifact_count_excluding_seal': len(artifacts),
                  'seal': str((evidence / 'seal.json').relative_to(root)),
                  'seal_identity_for_handoff_only': pin(evidence / 'seal.json')}, indent=2))
