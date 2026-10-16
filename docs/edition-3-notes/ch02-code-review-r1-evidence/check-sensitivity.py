#!/usr/bin/env python3
"""Revert each reviewed fix only in an isolated, disposable student source copy."""
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path('/Users/bill/projects/ensemble-edition-3-student')
OUT = Path(__file__).resolve().parent

def replace_once(text, old, new):
    if text.count(old) != 1:
        raise AssertionError('mutation anchor is not unique')
    return text.replace(old, new)

def r1(text):
    text = replace_once(text, '\t\t\t// Resolve only calls already visited, including sequential absent IDs.\n\t\t\tnames[p.CallID] = p.Name\n', '')
    return replace_once(text, '\tnames := map[string]string{}\n', '''\tnames := map[string]string{}
\tfor _, turn := range c.Dialogue {
\t\tfor _, p := range turn.Parts {
\t\t\tif p.Type == "tool_call" {
\t\t\t\tnames[p.CallID] = p.Name
\t\t\t}
\t\t}
\t}
''')

def r2(text):
    for start, end in [('func parseOpenAI', 'func parseGemini'), ('func parseGemini', '// measuredUsage')]:
        a, b = text.index(start), text.index(end)
        function = text[a:b]
        usage = next(line for line in function.splitlines(keepends=True) if 'result.Usage = common.Usage' in line)
        function = replace_once(function, usage, '')
        function = replace_once(function, '\tif len(result.Parts) == 0 {', usage + '\tif len(result.Parts) == 0 {')
        text = text[:a] + function + text[b:]
    return text

def r3(text):
    return replace_once(text, '''\t\t\t\t\tif opaque.ReasoningContent != nil {
\t\t\t\t\t\tmessage.ReasoningContent = opaque.ReasoningContent
\t\t\t\t\t}
\t\t\t\t\tif opaque.Reasoning != nil {
\t\t\t\t\t\tmessage.Reasoning = opaque.Reasoning
\t\t\t\t\t}''', '''\t\t\t\t\tmessage.ReasoningContent = opaque.ReasoningContent
\t\t\t\t\tmessage.Reasoning = opaque.Reasoning''')

tracked = subprocess.check_output(['git', 'ls-files', '-z'], cwd=ROOT).decode().split('\0')
files = set(filter(None, tracked)) | {'chapter2_revision_test.go'}
source_hashes = {name: hashlib.sha256((ROOT/name).read_bytes()).hexdigest() for name in sorted(files)}
results = []
with tempfile.TemporaryDirectory(prefix='scratch-sensitivity-', dir=OUT) as directory:
    scratch = Path(directory)
    for name in files:
        target = scratch/name
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(ROOT/name, target)
    for name, path, selector, mutate in [
        ('R1', 'internal/engine/render.go', 'TestChapter2GeminiChronologicalNames', r1),
        ('R2', 'internal/engine/parse.go', 'TestChapter2MeasuredEmptyReplies', r2),
        ('R3', 'internal/engine/render.go', 'TestChapter2OpenAIMultipleOpaqueParts', r3),
    ]:
        target = scratch/path
        original = target.read_bytes()
        command = ['go', 'test', '-count=1', '-run', '^'+selector+'$', '.']
        target.write_text(mutate(original.decode()))
        try:
            failed = subprocess.run(command, cwd=scratch, capture_output=True, text=True)
        finally:
            target.write_bytes(original)
        restored = subprocess.run(command, cwd=scratch, capture_output=True, text=True)
        record = {'finding': name, 'command': command, 'scratch_only': True,
                  'reverted_exit': failed.returncode, 'restored_exit': restored.returncode,
                  'restored_source_sha256': hashlib.sha256(target.read_bytes()).hexdigest()}
        (OUT/(name+'-sensitivity.txt')).write_text(json.dumps(record, indent=2)+'\n\nREVERTED FIX\n'+failed.stdout+failed.stderr+'\nRESTORED FIX\n'+restored.stdout+restored.stderr)
        assert failed.returncode != 0 and '--- FAIL: '+selector in failed.stdout, record
        assert 'build failed' not in failed.stdout+failed.stderr and 'no tests to run' not in failed.stdout+failed.stderr, record
        assert restored.returncode == 0, record
        assert target.read_bytes() == original
        print(name, 'behaviorally failed when reverted; restored pass')
        results.append(record)
assert source_hashes == {name: hashlib.sha256((ROOT/name).read_bytes()).hexdigest() for name in source_hashes}
(OUT/'sensitivity-summary.json').write_text(json.dumps({'production_unchanged': True, 'source_files_sha256': source_hashes, 'checks': results}, indent=2)+'\n')
