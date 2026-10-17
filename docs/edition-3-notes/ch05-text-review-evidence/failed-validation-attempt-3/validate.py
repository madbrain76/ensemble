#!/usr/bin/env python3
"""Read-only prose validation; writes only this review's evidence directory."""
from pathlib import Path
import difflib
import hashlib
import json
import os
import re
import shutil
import subprocess
import time

ROOT = Path('/Users/bill/projects/ensemble')
NOTES = ROOT / 'docs/edition-3-notes'
OUT = NOTES / 'ch05-text-review-evidence'
AUTHOR = NOTES / 'ch05-author-evidence'
CHAPTER = ROOT / 'book/edition-3/chapter-05.md'
BASE = Path('/tmp/cadence-ch05-poststudent-baseline/chapter-05.md')
REPORT = NOTES / 'ch05-text-review.md'
OUT.mkdir(exist_ok=True)

def sha(p):
    return hashlib.sha256(p.read_bytes()).hexdigest()

def save(name, data):
    p = OUT / name
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(json.dumps(data, indent=2) + '\n')

commands = []
def run(name, argv, env=None):
    start = time.monotonic()
    proc = subprocess.run(argv, cwd=ROOT, env={**os.environ, **(env or {})}, capture_output=True)
    (OUT / (name + '.stdout.txt')).write_bytes(proc.stdout)
    (OUT / (name + '.stderr.txt')).write_bytes(proc.stderr)
    (OUT / (name + '.exit.txt')).write_text(str(proc.returncode) + '\n')
    rec = {'name': name, 'argv': argv, 'cwd': str(ROOT), 'environment_overrides': env or {},
           'exit': proc.returncode, 'elapsed_seconds': round(time.monotonic() - start, 3),
           'stdout': name + '.stdout.txt', 'stderr': name + '.stderr.txt'}
    commands.append(rec)
    save('commands.json', commands)
    print(name, 'exit', proc.returncode, flush=True)
    return proc

def owned_status(line):
    return 'docs/edition-3-notes/ch05-text-review-evidence/' in line or line[3:] == 'docs/edition-3-notes/ch05-text-review.md'

before = run('status-before', ['git', '--no-pager', 'status', '--short', '--untracked-files=all'])
tracked = run('tracked-before', ['git', '--no-pager', 'diff', '--name-only'])
staged = run('staged-before', ['git', '--no-pager', 'diff', '--cached', '--name-only'])
head = run('main-head', ['git', '--no-pager', 'rev-parse', 'HEAD'])
tag = run('accepted-tag', ['git', '--no-pager', '-C', '/Users/bill/projects/ensemble-edition-3-student', 'rev-parse', 'edition-3-ch05-r1^{commit}'])

inputs = [ROOT / 'AGENTS.md', ROOT / '.agents/skills/edition-3-text-reviewer/SKILL.md',
          ROOT / '.agents/skills/edition-3-author/SKILL.md',
          *[ROOT / ('docs/' + s + '.md') for s in ['edition-3-execution','edition-3-readiness','edition-3-student','edition-3-carryover','agentic-codebook-workflow','edition-2-carryover']],
          ROOT / 'book/voice.md', ROOT / 'book/chapter-writing-procedure.md', ROOT / 'book/edition-3/chapter-04.md', CHAPTER, BASE,
          *[NOTES / (s + '.md') for s in ['ch05-author','ch05-student-review','ch05-student-review-r1','ch05-code-review','ch05-code-review-r1','ch05-wakeup-student','ch05-wakeup-code-review','job-wakeup-amendment','job-wakeup-amendment-author','job-wakeup-amendment-review','ch05-execution']],
          NOTES / 'ch05-accepted.json', NOTES / 'ch05-wakeup-student-evidence-manifest.json', AUTHOR / 'SHA256SUMS',
          *[ROOT / p for p in tracked.stdout.decode().splitlines()]]
inputs = sorted(set(inputs))
protected_before = {str(p): sha(p) for p in inputs}
save('input-hashes.json', protected_before)
expected = {str(CHAPTER): '475d807b1fcbd9e3bbb475afc08af7baea6cae438e77d0c404c9030782e7aba9',
            str(BASE): '4f4c2595ef65a08511d5cc68553f853feaf854e4cb0fe7d4d9094d4ae7f5a344',
            str(NOTES/'ch05-author.md'): '534b9ecdef4f691f6715eca9bb5f50572c6f9aada9673deeb39ec5b8ff98fa5d',
            str(AUTHOR/'SHA256SUMS'): '468f351f542e0be8abc70385a07d7c192679dacb5978238acb2ef4db121c3253',
            str(AUTHOR/'chapter.diff'): 'f6797409e78be3eabda0bb53f2415c5e98bec61f65a8558bc9a890ff4de4c3b7'}
pins = {p: {'expected': h, 'actual': sha(Path(p)), 'matches': sha(Path(p)) == h} for p,h in expected.items()}
save('frozen-pins.json', pins)
assert all(v['matches'] for v in pins.values())

manifest = []
for line in (AUTHOR/'SHA256SUMS').read_text().splitlines():
    h, rel = line.split(None, 1)
    p = AUTHOR / rel.strip()
    manifest.append({'path': rel.strip(), 'expected': h, 'actual': sha(p), 'matches': sha(p) == h})
save('author-manifest-verification.json', manifest)
assert len(manifest) == 23 and all(r['matches'] for r in manifest)

for name, p in [('baseline', BASE), ('candidate', CHAPTER)]:
    dest = OUT/name/'chapter-05.md'
    dest.parent.mkdir(exist_ok=True)
    shutil.copyfile(p, dest)
assert (AUTHOR/'baseline/chapter-05.md').read_bytes() == BASE.read_bytes()
assert (AUTHOR/'final/chapter-05.md').read_bytes() == CHAPTER.read_bytes()
shutil.copyfile(AUTHOR/'chapter.diff', OUT/'chapter.diff')
receipt_names = ['lint-baseline.command.txt','lint-baseline.exit.txt','lint-baseline.stdout.txt','lint-baseline.stderr.txt',
                 'lint-final.command.txt','lint-final.exit.txt','lint-final.stdout.txt','lint-final.stderr.txt',
                 'preservation-checks.json','preservation-initial-failed.json','final-initial-failed.json',
                 'final-checks.txt','validation-notes.txt','parent-confirmation.txt']
(OUT/'author-receipts').mkdir(exist_ok=True)
for n in receipt_names:
    shutil.copyfile(AUTHOR/n, OUT/'author-receipts'/n)
(OUT/'historical-pty').mkdir(exist_ok=True)
for rel, dest in [('gates/root-race.txt','intermediate-root-race.txt'), ('gates/commands.json','intermediate-commands.json'),
                  ('candidate-r2/frozen/grade-ch05.json','accepted-grade-ch05.json'),
                  ('candidate-r2/frozen/grade-commands-and-exits.json','accepted-grade-commands-and-exits.json')]:
    source = NOTES/'ch05-wakeup-student-evidence'/rel
    meta = json.loads((NOTES/'ch05-wakeup-student-evidence-manifest.json').read_text())['text_files'][rel]
    assert sha(source) == meta['sha256']
    shutil.copyfile(source, OUT/'historical-pty'/dest)

b, c = BASE.read_text(), CHAPTER.read_text()
ops = difflib.SequenceMatcher(a=b.splitlines(keepends=True), b=c.splitlines(keepends=True), autojunk=False).get_opcodes()
changed = [o for o in ops if o[0] != 'equal']
assert len(changed) == 2 and all(o[0] == 'insert' for o in changed)
assert sum(o[4]-o[3] for o in changed) == 15
rebuilt = ''.join(''.join(c.splitlines(keepends=True)[j1:j2]) for tag, i1, i2, j1, j2 in ops if tag == 'equal')
assert rebuilt == b
fences = lambda s: re.findall(r'^```[^\n]*\n.*?^```[^\n]*$', s, flags=re.M|re.S)
heading = lambda s: [x for x in s.splitlines() if x.startswith('#')]
tables = lambda s: [x for x in s.splitlines() if x.startswith('|')]
checks = {'two_insertions_only': True, 'added_lines': 15, 'deleted_or_replaced_lines': 0,
          'removing_insertions_reconstructs_exact_baseline': rebuilt == b,
          'code_fences_identical': fences(b) == fences(c), 'headings_identical': heading(b) == heading(c),
          'table_lines_identical': tables(b) == tables(c),
          'all_preexisting_text_including_Bill_preface_and_accepted_amendment_preserved': rebuilt == b,
          'insertions': [{'baseline_after_line': i1, 'candidate_start_line': j1+1, 'candidate_end_line': j2,
                          'text': ''.join(c.splitlines(keepends=True)[j1:j2])} for _,i1,i2,j1,j2 in changed]}
save('preservation.json', checks)
assert all(checks[x] for x in ['code_fences_identical','headings_identical','table_lines_identical'])

links = []
for source in [CHAPTER, NOTES/'ch05-author.md']:
    for match in re.finditer(r'!?\[[^\]]*\]\(([^\s)]+)(?:\s+[^)]*)?\)', source.read_text()):
        target = match.group(1)
        if re.match(r'^[a-zA-Z][a-zA-Z0-9+.-]*:', target) or target.startswith('#'):
            continue
        path = source.parent / target.split('#')[0]
        links.append({'source': str(source), 'target': target, 'resolved': str(path.resolve()), 'exists': path.exists()})
save('local-links.json', {'scope': 'Explicit inline Markdown file destinations in chapter and author report; no network requests, no anchor/content verification.', 'links': links})
assert all(r['exists'] for r in links)

for label in ['baseline','candidate']:
    result = run('lint-'+label, ['go','run','./cmd/prose-lint',str(OUT/label/'chapter-05.md')],
                 {'GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off'})
    assert result.returncode == 0
    original = (AUTHOR/('lint-baseline.stdout.txt' if label == 'baseline' else 'lint-final.stdout.txt')).read_text()
    assert result.stdout.decode().splitlines()[1:] == original.splitlines()[1:]
    assert result.stderr == b''
ws = run('chapter-whitespace', ['git','--no-pager','diff','--check','--','book/edition-3/chapter-05.md'])
assert ws.returncode == 0

after = run('status-after-validation', ['git','--no-pager','status','--short','--untracked-files=all'])
tracked_after = run('tracked-after-validation', ['git','--no-pager','diff','--name-only'])
staged_after = run('staged-after-validation', ['git','--no-pager','diff','--cached','--name-only'])
protected_after = {str(p): sha(p) for p in inputs}
a = set(x for x in before.stdout.decode().splitlines() if not owned_status(x))
z = set(x for x in after.stdout.decode().splitlines() if not owned_status(x))
scope = {'protected_input_hashes_unchanged': protected_before == protected_after,
         'preexisting_status_entries_unchanged': a == z,
         'status_added_outside_review': sorted(z-a), 'status_removed_outside_review': sorted(a-z),
         'tracked_diff_names_unchanged': tracked.stdout == tracked_after.stdout,
         'staged_diff_unchanged_and_empty': staged.stdout == staged_after.stdout == b'',
         'existing_untracked_file_entries': len([x for x in a if x.startswith('?? ')]),
         'qualification': 'Hashes cover named governing documents, reports, chapter, baseline, and all five already-modified tracked files, not every unrelated untracked file. Full before/after untracked path lists retained.'}
save('scope-validation.json', scope)
assert all(scope[k] for k in ['protected_input_hashes_unchanged','preexisting_status_entries_unchanged','tracked_diff_names_unchanged','staged_diff_unchanged_and_empty'])
assert head.stdout.decode().strip() == '4f1701a76f57ec57d362d4f1a30264f30b65ddf8'
assert tag.stdout.decode().strip() == 'f60ebf3aad2ece4251face9d95c106dd667abcd6'
print('PASS: pins, 23 author files, two insertion-only paragraphs, exact preservation, saved lint parity, links, whitespace and bounded checkout scope.')
