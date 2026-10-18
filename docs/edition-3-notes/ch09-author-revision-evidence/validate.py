#!/usr/bin/env python3
"""Scoped revision validation; never execute student code or graders."""
import datetime
import hashlib
import json
import os
import pathlib
import shutil
import subprocess
import time

root = pathlib.Path('/Users/bill/projects/ensemble')
evidence = root / 'docs/edition-3-notes/ch09-author-revision-evidence'
review = root / 'docs/edition-3-notes/ch09-text-review-evidence'

def pin(path):
    data = path.read_bytes()
    return {'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest()}

def run(name, argv, cwd, expected):
    env = {key: os.environ[key] for key in ('PATH', 'HOME', 'TMPDIR') if key in os.environ}
    env.update({'GIT_PAGER': 'cat', 'PAGER': 'cat', 'GIT_TERMINAL_PROMPT': '0'})
    assert shutil.disk_usage(root).free >= 1024**3
    started = datetime.datetime.now(datetime.timezone.utc).isoformat()
    t0 = time.monotonic()
    result = subprocess.run(argv, cwd=cwd, env=env, capture_output=True)
    for suffix, data in [('stdout.txt', result.stdout), ('stderr.txt', result.stderr)]:
        with (evidence / f'{name}.{suffix}').open('xb') as stream:
            stream.write(data)
    receipt = {'argv': argv, 'cwd': str(cwd), 'environment': env,
               'started_utc': started, 'duration_seconds': time.monotonic() - t0,
               'exit_code': result.returncode, 'expected_exit_codes': expected,
               'stdout': f'{name}.stdout.txt', 'stderr': f'{name}.stderr.txt'}
    with (evidence / f'{name}.json').open('x') as stream:
        json.dump(receipt, stream, indent=2)
        stream.write('\n')
    print(name, 'exit', result.returncode)
    assert result.returncode in expected, name
    return result.stdout

run('revision-diff', ['git', '--no-pager', 'diff', '--no-index', '--',
    'docs/edition-3-notes/ch09-author-revision-evidence/reviewed-draft/chapter-09.md',
    'book/edition-3/chapter-09.md'], root, [1])
run('chapter-whitespace', ['git', '--no-pager', 'diff', '--check', '--',
    'book/edition-3/chapter-09.md'], root, [0])
run('voice-comparison', ['cmp', '--', 'book/voice.md', 'book/edition-3/voice.md'], root, [0])
source_receipt = json.loads((review / 'student-tag.json').read_text())
source_cwd = source_receipt['cwd']
resolved = run('accepted-tag', ['git', '--no-pager', 'rev-parse',
    'edition-3-ch09-r1^{commit}'], source_cwd, [0])
assert resolved == b'251ce88fb1677897dbcc6b33661096c5cb39cd8e\n'

initial = json.loads((evidence / 'initial-pins.json').read_text())
protected = {}
chapter_name = 'book/edition-3/chapter-09.md'
for name, expected in initial['inputs'].items():
    if name == chapter_name:
        continue
    actual = pin(root / name)
    assert actual == expected, name
    protected[name] = actual

before = (evidence / 'reviewed-draft/chapter-09.md').read_text()
after = (root / chapter_name).read_text()
grade_start = 'From the book repository root, grade that executable'
grade_end = '## The Format'
registry_start = '## The Registry'
registry_end = '## Progressive Disclosure'
for start, end in [(grade_start, grade_end), (registry_start, registry_end)]:
    assert before.split(start, 1)[1].split(end, 1)[0] == after.split(start, 1)[1].split(end, 1)[0]
assert './cmd/ensemble' not in after
assert 'prompt TEXT' not in after
assert '/tmp/ensemble-ch09 prompt ' not in after
assert 'ensemble.AgentSpec' not in after
assert 'go build -o /tmp/ensemble-ch09 ./cmd/ch01' in after
assert '/tmp/ensemble-ch09 chat' in after
assert 'cfg.Skills = &ensemble.SkillSpec{' in after
assert 'agent, err := app.NewAgent(cfg)' in after
assert 'if err != nil {\n    return err\n}' in after

def headings(text):
    result = []
    fenced = False
    for line in text.splitlines():
        if line.startswith('```'):
            fenced = not fenced
        elif not fenced and line.startswith('#'):
            result.append(line)
    return result

assert headings(before) == headings(after)
result = {'checked_utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
          'all_protected_inputs_unchanged': True, 'protected_inputs': protected,
          'historical_chapter_pin': pin(evidence / 'reviewed-draft/chapter-09.md'),
          'final_chapter_pin': pin(root / chapter_name),
          'historical_pin_rule': 'Old author/review seals remain byte-identical; their chapter entry now names the preserved reviewed snapshot, not current manuscript bytes.',
          'grade_recipe_and_registry_section_byte_identical': True,
          'required_corrections_present': ['R1', 'R2', 'R3'],
          'before': {'words_whitespace_split': len(before.split()), 'lines': len(before.splitlines()), 'headings': headings(before)},
          'after': {'words_whitespace_split': len(after.split()), 'lines': len(after.splitlines()), 'headings': headings(after)},
          'free_bytes': shutil.disk_usage(root).free}
with (evidence / 'validation.json').open('x') as stream:
    json.dump(result, stream, indent=2)
    stream.write('\n')
print('protected inputs unchanged:', len(protected))
print('final chapter:', result['final_chapter_pin'])
print('before/after words:', len(before.split()), len(after.split()))
print('before/after lines:', len(before.splitlines()), len(after.splitlines()))
