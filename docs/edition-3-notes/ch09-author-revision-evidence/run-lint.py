#!/usr/bin/env python3
"""Bounded revision lint receipts; no tests, graders, or provider calls."""
import datetime
import hashlib
import json
import os
import pathlib
import shutil
import subprocess
import sys
import time

root = pathlib.Path('/Users/bill/projects/ensemble')
evidence = root / 'docs/edition-3-notes/ch09-author-revision-evidence'
phase = sys.argv[1]
assert phase in ('before', 'after')
target = 'docs/edition-3-notes/ch09-author-revision-evidence/reviewed-draft/chapter-09.md' if phase == 'before' else 'book/edition-3/chapter-09.md'
command = ['go', 'run', './cmd/lintprose', '-v', target]
env = {key: os.environ[key] for key in ('PATH', 'HOME', 'TMPDIR') if key in os.environ}
env.update({'GOPROXY': 'off', 'GOSUMDB': 'off', 'GOTOOLCHAIN': 'local', 'GOWORK': 'off', 'GOENV': 'off'})
free_before = shutil.disk_usage(root).free
assert free_before >= 1024**3, 'less than 1 GiB free'
started = datetime.datetime.now(datetime.timezone.utc).isoformat()
t0 = time.monotonic()
result = subprocess.run(command, cwd=root, env=env, capture_output=True)
duration = time.monotonic() - t0
for suffix, data in [('stdout.txt', result.stdout), ('stderr.txt', result.stderr)]:
    with (evidence / f'lint-{phase}.{suffix}').open('xb') as stream:
        stream.write(data)
def pin(path):
    data = path.read_bytes()
    return {'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest()}
receipt = {
    'argv': command, 'cwd': str(root), 'environment': env,
    'started_utc': started, 'duration_seconds': duration,
    'exit_code': result.returncode, 'free_bytes_before': free_before,
    'free_bytes_after': shutil.disk_usage(root).free,
    'target': {'path': target, **pin(root / target)},
    'linter': {'path': 'cmd/lintprose/main.go', **pin(root / 'cmd/lintprose/main.go')},
    'stdout': f'lint-{phase}.stdout.txt', 'stderr': f'lint-{phase}.stderr.txt',
}
with (evidence / f'lint-{phase}.json').open('x') as stream:
    json.dump(receipt, stream, indent=2)
    stream.write('\n')
print(json.dumps(receipt, indent=2))
print(result.stdout.decode(), end='')
print(result.stderr.decode(), end='', file=sys.stderr)
sys.exit(result.returncode)
