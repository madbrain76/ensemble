#!/bin/sh
set -eu
ROOT=/Users/bill/projects/ensemble-edition-3-student
EV=/Users/bill/projects/ensemble/docs/edition-3-notes/ch10-student-evidence
cd "$ROOT"
set +e
python3 .ch10-local/receipt.py sensitivity-mutant-v1 "$ROOT/.ch10-local/mutation-729bc91" go test -count=1 -timeout 120s -run '^TestChapter10SnapshotReplayTailAndTrim$' .
negative=$?
set -e
[ "$negative" -eq 1 ]
python3 - <<'PY'
from pathlib import Path
import json,hashlib
root=Path.cwd();ev=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch10-student-evidence')
r=json.loads((ev/'sensitivity-mutant-v1.json').read_text());text=(ev/'sensitivity-mutant-v1.stdout.txt').read_text()+(ev/'sensitivity-mutant-v1.stderr.txt').read_text()
assert r['exit_code']==1 and '--- FAIL: TestChapter10SnapshotReplayTailAndTrim' in text and '[build failed]' not in text
m=json.loads((ev.parent/'ch10-student-source-manifest.json').read_text());copy=root/'.ch10-local/mutation-729bc91';p=copy/'internal/history/history.go';p.write_bytes((root/'internal/history/history.go').read_bytes())
for entry in m['tracked_files']:
 b=(copy/entry['path']).read_bytes();assert len(b)==entry['bytes'] and hashlib.sha256(b).hexdigest()==entry['sha256'],entry['path']
(ev/'sensitivity-restoration-inventory.json').write_text(json.dumps({'source_commit':m['head'],'source_tree':m['tree'],'restored_files_verified':len(m['tracked_files']),'restored_history_sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'compile_valid_negative_assertion_observed':True,'source_worktree_was_not_mutated':True},indent=2)+'\n')
PY
python3 .ch10-local/receipt.py sensitivity-restored-v1 "$ROOT/.ch10-local/mutation-729bc91" go test -count=1 -timeout 120s -run '^TestChapter10SnapshotReplayTailAndTrim$' .
python3 .ch10-local/receipt.py final-real-process-v1 "$ROOT" python3 "$ROOT/.ch10-local/process_control.py"
python3 .ch10-local/receipt.py grader-ch10-final-v1 /Users/bill/projects/ensemble /tmp/cadence-ch10-grader -ch 10 -json -core "$ROOT" -consumer "$ROOT/examples/ch04-consumer" -consumer-case "$ROOT/examples/ch04-consumer/consumer-case.json" -workflow "$ROOT/examples/ch06" -media-models "$ROOT/.ch10-local/media-models.json" -stream-exercise "$ROOT/examples/ch07" -stream-model claude-sonnet-4-6 -gui "$ROOT/gui" -gui-model claude-sonnet-4-6 "$ROOT/.ch10-local/bin/frozen-cli"
python3 .ch10-local/receipt.py grader-ch1-final-v1 /Users/bill/projects/ensemble /tmp/cadence-ch10-grader -ch 1 -json "$ROOT/.ch10-local/bin/frozen-cli"
