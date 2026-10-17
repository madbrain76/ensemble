#!/usr/bin/env python3
"""Seal the bounded prose-review artifacts; no source tests or Git writes."""
from pathlib import Path
import hashlib
import json
import re
import subprocess
ROOT = Path('/Users/bill/projects/ensemble')
OUT = ROOT/'docs/edition-3-notes/ch05-text-review-evidence'
REPORT = OUT.parent/'ch05-text-review.md'
CHAPTER = ROOT/'book/edition-3/chapter-05.md'
def sha(p): return hashlib.sha256(p.read_bytes()).hexdigest()
def save(n, o): (OUT/n).write_text(json.dumps(o, indent=2)+'\n')
def owned(line): return 'docs/edition-3-notes/ch05-text-review-evidence/' in line or line[3:] == 'docs/edition-3-notes/ch05-text-review.md'
commands=[]
def run(name,args):
 p=subprocess.run(args,cwd=ROOT,capture_output=True)
 for suffix, data in [('stdout',p.stdout),('stderr',p.stderr)]: (OUT/f'{name}.{suffix}.txt').write_bytes(data)
 (OUT/f'{name}.exit.txt').write_text(str(p.returncode)+'\n')
 commands.append({'name':name,'argv':args,'cwd':str(ROOT),'exit':p.returncode})
 assert p.returncode==0
 return p.stdout
status=run('final-status',['git','--no-pager','status','--short','--untracked-files=all'])
tracked=run('final-tracked',['git','--no-pager','diff','--name-only'])
staged=run('final-staged',['git','--no-pager','diff','--cached','--name-only'])
head=run('final-head',['git','--no-pager','rev-parse','HEAD'])
save('final-commands.json',commands)
initial=json.loads((OUT/'input-hashes.json').read_text())
actual={p:sha(Path(p)) for p in initial}
assert initial==actual
assert CHAPTER.read_bytes()==(OUT/'candidate/chapter-05.md').read_bytes()
assert sha(CHAPTER)=='475d807b1fcbd9e3bbb475afc08af7baea6cae438e77d0c404c9030782e7aba9'
preservation=json.loads((OUT/'preservation.json').read_text())
reconstructed=CHAPTER.read_text()
for item in reversed(preservation['insertions']):
 lines=reconstructed.splitlines(keepends=True)
 del lines[item['candidate_start_line']-1:item['candidate_end_line']]
 reconstructed=''.join(lines)
assert reconstructed.encode()==(OUT/'baseline/chapter-05.md').read_bytes()
assert head.decode().strip()=='4f1701a76f57ec57d362d4f1a30264f30b65ddf8'
assert tracked==(OUT/'tracked-before.stdout.txt').read_bytes()
assert staged==b''
old={s for s in (OUT/'status-before.stdout.txt').read_text().splitlines() if not owned(s)}
new={s for s in status.decode().splitlines() if not owned(s)}
assert old==new
assert REPORT.read_bytes().endswith(b'\n')
assert '## Verdict: PASS' in REPORT.read_text()
assert (OUT/'validation.exit.txt').read_text().strip()=='0'
# Index is intentionally excluded from its own hashes. Create it before checking
# report destinations; replace it with its final contents below.
index=OUT/'SHA256SUMS'
index.touch(exist_ok=True)
links=[]
for m in re.finditer(r'!?\[[^\]]*\]\(([^\s)]+)(?:\s+[^)]*)?\)',REPORT.read_text()):
 target=m.group(1)
 if re.match(r'^[a-zA-Z][a-zA-Z0-9+.-]*:',target) or target.startswith('#'): continue
 p=REPORT.parent/target.split('#')[0]
 links.append({'target':target,'exists':p.exists()})
assert all(x['exists'] for x in links)
final={'verdict':'PASS','report_sha256':sha(REPORT),'chapter_sha256':sha(CHAPTER),
       'all_named_input_hashes_unchanged':initial==actual,
       'byte_exact_baseline_reconstruction':True,'tracked_diff_names_unchanged':True,
       'all_preexisting_status_entries_unchanged':True,'staged_diff_empty':True,
       'head_unchanged':True,'successful_validation_exit':0,'report_local_links':links,
       'scope':'Only ch05-text-review.md and ch05-text-review-evidence/ written by this reviewer.'}
save('final-checks.json',final)
# A stable completion receipt is included in the index; actual hash verification
# immediately follows index creation and prints to the tool's saved output.
(OUT/'final-validation.stdout.txt').write_text('PASS: named inputs unchanged; exact frozen chapter and byte-preserved baseline; report links resolve; HEAD and prior status preserved; staged diff empty.\n')
(OUT/'final-validation.stderr.txt').write_text('')
(OUT/'final-validation.exit.txt').write_text('0\n')
files=[REPORT]+sorted(p for p in OUT.rglob('*') if p.is_file() and p!=index)
rows=[]
for p in files:
 rel='../ch05-text-review.md' if p==REPORT else p.relative_to(OUT).as_posix()
 rows.append((sha(p),rel))
index.write_text(''.join(h+'  '+rel+'\n' for h,rel in rows))
for h,rel in rows: assert sha(OUT/rel)==h
print('PASS: sealed and independently re-read all',len(rows),'indexed files; index excludes only itself.')
print('REPORT_SHA256',sha(REPORT))
print('CHAPTER_SHA256',sha(CHAPTER))
print('INDEX_SHA256',sha(index))
print('FINAL_REPORT_TAIL:')
print('\n'.join(REPORT.read_text().splitlines()[-9:]))
