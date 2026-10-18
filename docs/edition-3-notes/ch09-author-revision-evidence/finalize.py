#!/usr/bin/env python3
"""Seal the single bounded R1-R3 revision, preserving historical artifacts."""
import datetime
import hashlib
import json
import pathlib
import shutil

root = pathlib.Path('/Users/bill/projects/ensemble')
evidence = root / 'docs/edition-3-notes/ch09-author-revision-evidence'
report = root / 'docs/edition-3-notes/ch09-author-revision.md'
chapter = root / 'book/edition-3/chapter-09.md'
assert not (evidence / 'seal.json').exists()

def pin(path):
    data = path.read_bytes()
    return {'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest()}

initial = json.loads((evidence / 'initial-pins.json').read_text())
protected = {}
for name, expected in initial['inputs'].items():
    if name == 'book/edition-3/chapter-09.md':
        continue
    actual = pin(root / name)
    assert actual == expected, name
    protected[name] = actual
assert len(protected) == 92
assert pin(chapter) == {
    'bytes': 20888,
    'sha256': 'e2c82e17dc2cdeb4c580afe182ecef60bf6393fe10e234f56a7a255bbc9a5fc0',
}
snapshot = evidence / 'reviewed-draft/chapter-09.md'
assert pin(snapshot) == initial['inputs']['book/edition-3/chapter-09.md']
assert snapshot.read_bytes() == pathlib.Path('/tmp/cadence-ch09-text-review-initial/chapter-09.md').read_bytes()
assert (root / 'book/voice.md').read_bytes() == (root / 'book/edition-3/voice.md').read_bytes()
assert shutil.disk_usage(root).free >= 1024**3

# Preserve raw captures and copied snapshots verbatim; check authored text only.
whitespace_files = [chapter, report] + sorted(evidence.glob('*.py')) + sorted(evidence.glob('*.json'))
for path in whitespace_files:
    text = path.read_text()
    assert all(line.rstrip() == line for line in text.splitlines()), str(path)
    assert not text or text.endswith('\n'), str(path)

check = {
    'argv': ['python3', 'docs/edition-3-notes/ch09-author-revision-evidence/finalize.py'],
    'cwd': str(root),
    'checked_utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
    'checks_passed': True,
    'protected_non_manuscript_inputs_unchanged': len(protected),
    'reviewed_draft_matches_immutable_snapshot': True,
    'final_chapter_matches_after_lint_target': True,
    'voice_byte_identical': True,
    'authored_text_whitespace_and_newline_checked': [str(p.relative_to(root)) for p in whitespace_files],
    'raw_outputs_and_source_snapshots_not_normalized': True,
    'free_bytes': shutil.disk_usage(root).free,
}
with (evidence / 'final-check.json').open('x') as stream:
    json.dump(check, stream, indent=2)
    stream.write('\n')

artifacts = {}
for path in [chapter, report] + sorted(evidence.rglob('*')):
    assert not path.is_symlink(), str(path)
    if path.is_dir():
        continue
    assert path.is_file(), str(path)
    data = path.read_bytes()
    data.decode('utf-8')
    assert b'\x00' not in data, str(path)
    artifacts[str(path.relative_to(root))] = {**pin(path), 'utf8': True, 'nul_free': True}
seal = {
    'sealed_utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
    'task': pin(evidence / 'task.txt'),
    'hash_algorithm': 'SHA-256',
    'inventory_rule': 'Final chapter, revision report, and every new revision evidence file; exact recursive inventory; seal.json excluded to avoid self-hash. No subsequent writes.',
    'artifacts': artifacts,
    'protected_inputs': protected,
    'historical_chapter_rule': 'Original author/review chapter pins describe the unchanged portable reviewed-draft snapshot, not the revised working manuscript. Original seals and all their non-manuscript artifacts remain unchanged.',
    'historical_reviewed_chapter': pin(snapshot),
    'all_artifacts_utf8_and_nul_free': True,
    'sole_revision_requests': ['R1', 'R2', 'R3'],
    'second_formal_review': False,
    'acceptance': 'not claimed; root directly verifies this sole revision',
    'unfinished_jobs': [],
    'children_spawned': 0,
}
with (evidence / 'seal.json').open('x') as stream:
    json.dump(seal, stream, indent=2)
    stream.write('\n')
print(json.dumps({'checks_passed': True, 'sealed_artifacts_excluding_seal': len(artifacts),
                  'report': pin(report), 'chapter': pin(chapter),
                  'seal_identity_for_handoff_only': pin(evidence / 'seal.json')}, indent=2))
