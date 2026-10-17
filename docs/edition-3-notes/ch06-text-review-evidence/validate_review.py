#!/usr/bin/env python3
"""Bounded text review checks only; no behavioral or provider validation."""
from pathlib import Path
from urllib.parse import urlsplit, unquote
import datetime, difflib, hashlib, json, re, subprocess, sys
R=Path('/Users/bill/projects/ensemble')
S=Path('/Users/bill/projects/ensemble-edition-3-student')
E=R/'docs/edition-3-notes/ch06-text-review-evidence'
A=R/'docs/edition-3-notes/ch06-author-evidence'
errors=[]
def sha(p): return hashlib.sha256(p.read_bytes()).hexdigest()
def run(args,cwd=R):
 p=subprocess.run(args,cwd=cwd,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
 return {'command':args,'cwd':str(cwd),'exit':p.returncode,'output':p.stdout}
def check(ok,msg):
 if not ok: errors.append(msg)
checks={}
for label,args,cwd,expected in [
 ('main_head',['git','--no-pager','rev-parse','HEAD'],R,'b67386ab80cf2afe7985db5a32ab99ec2311f090'),
 ('student_head',['git','--no-pager','rev-parse','HEAD'],S,'f595c787675e2a6343053ec42ea6c7e1e02a930e'),
 ('student_tree',['git','--no-pager','rev-parse','HEAD^{tree}'],S,'6fb36168020a22b8993412812987338aa4b6e010'),
 ('student_tag_object',['git','--no-pager','rev-parse','edition-3-ch06-r1'],S,'84d0bce5f2525b406affa64fb1a8eb70711b1e06'),
 ('student_tag_commit',['git','--no-pager','rev-parse','edition-3-ch06-r1^{commit}'],S,'f595c787675e2a6343053ec42ea6c7e1e02a930e'),
 ('student_tag_type',['git','--no-pager','cat-file','-t','edition-3-ch06-r1'],S,'tag')]:
 receipt=run(args,cwd);receipt['expected']=expected
 check(receipt['exit']==0 and receipt['output'].strip()==expected,label)
 checks[label]=receipt
checks['student_status']=run(['git','--no-pager','status','--short'],S)
check(checks['student_status']['exit']==0 and not checks['student_status']['output'].strip(),'student worktree no longer clean')
(E/'pins.json').write_text(json.dumps(checks,indent=2)+'\n')
pins={'book/edition-3/chapter-06.md':'864d806c63c019ede50810ae5618b6fa4545f7028ae3abcb060722458babb4e9','docs/edition-3-notes/ch06-author.md':'94af4bedf32dba09c01bea6a01117a90b9f1ffea2aa41171e3d06ded0763f2d7','docs/edition-3-notes/ch06-author-evidence/chapter.diff':'9eebf86b7aad8d0be503969803d82f455117715591677239702be09672324973','docs/edition-3-notes/ch06-author-evidence/SHA256SUMS':'bdae2167a55ff3382b3dedcd2ac6e1070e7c518f84c87fd106e04fe1c61da71b'}
for p,want in pins.items():check(sha(R/p)==want,'frozen hash '+p)
base=E/'baseline/chapter-06.md'; final=E/'final/chapter-06.md'; chapter=R/'book/edition-3/chapter-06.md'
check(sha(base)=='ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55','baseline sha')
check(base.read_bytes()==(A/'baseline/chapter-06.md').read_bytes()==Path('/tmp/cadence-ch06-poststudent-baseline/chapter-06.md').read_bytes(),'baseline copies')
check(final.read_bytes()==chapter.read_bytes()==(A/'final/chapter-06.md').read_bytes(),'final copies')
manifest=[]
for line in (A/'SHA256SUMS').read_text().splitlines():
 digest,path=line.split('  ',1); p=R/path; b=p.read_bytes(); b.decode('utf-8')
 match=sha(p)==digest;check(match,'author manifest hash '+path);check(b'\x00' not in b,'author NUL '+path)
 manifest.append({'path':path,'bytes':len(b),'expected_sha256':digest,'actual_sha256':sha(p),'match':match})
expected={str(p.relative_to(R)) for p in A.rglob('*') if p.is_file() and p.name!='SHA256SUMS'}|{'docs/edition-3-notes/ch06-author.md','book/edition-3/chapter-06.md'}
actual=[i['path'] for i in manifest]
check(set(actual)==expected and len(actual)==len(expected)==31,'author exact manifest set')
(E/'author-manifest-verification.json').write_text(json.dumps({'exact_set':set(actual)==expected,'entries':len(actual),'files':manifest},indent=2)+'\n')
old=base.read_bytes().decode('utf-8');new=final.read_bytes().decode('utf-8')
changes=[]
for tag,i,j,k,l in difflib.SequenceMatcher(None,old.splitlines(),new.splitlines()).get_opcodes():
 if tag!='equal':changes.append({'type':tag,'baseline_lines':[i+1,j],'final_lines':[k+1,l]})
check(changes==[{'type':'replace','baseline_lines':[183,195],'final_lines':[183,200]},{'type':'insert','baseline_lines':[439,438],'final_lines':[444,453]},{'type':'replace','baseline_lines':[652,652],'final_lines':[667,670]}],'bounded change locations')
diff=''.join(difflib.unified_diff(old.splitlines(True),new.splitlines(True),fromfile='baseline/chapter-06.md',tofile='final/chapter-06.md'))
(E/'chapter.diff').write_text(diff)
check(diff.encode()==(A/'chapter.diff').read_bytes(),'portable diff equality')
check((E/'chapter.git.diff').read_bytes()==(A/'chapter.git.diff').read_bytes(),'git diff equality')
links=[]
for m in re.finditer(r'!?\[[^\]\n]+\]\(([^)\s]+)\)',new):
 target=m.group(1);u=urlsplit(target)
 if u.scheme or u.netloc:continue
 p=(chapter.parent/unquote(u.path)).resolve() if u.path else chapter
 item={'target':target,'resolved':str(p),'exists':p.exists(),'fragment':u.fragment}
 check(item['exists'],'local link '+target);check(not u.fragment,'unexpected fragment needing validation')
 links.append(item)
check(len(links)==7,'seven local links')
(E/'local-links.json').write_text(json.dumps(links,indent=2)+'\n')
check(new.endswith('\n') and '\x00' not in new,'chapter text format')
check(not any(s.rstrip(' \t')!=s for s in new.splitlines()),'chapter trailing whitespace')
for name in ['baseline-prose-lint','final-prose-lint','chapter-diff-check','chapter-git-diff']:
 check((E/(name+'.exit')).read_text().strip()=='0',name+' exit')
for name in ['baseline-prose-lint','final-prose-lint']:
 reviewer=(E/(name+'.txt')).read_text().splitlines()[1:]
 author=(A/(name+'.txt')).read_text().splitlines()[1:]
 check(reviewer==author,name+' diagnostics match author apart from input path')
concurrent=json.loads(Path('/tmp/cadence-ch06-concurrent-publication-state.json').read_text())['files']
preservation=[]
for p,meta in concurrent.items():
 same=sha(R/p)==meta['sha256'];check(same,'concurrent publication hash '+p)
 preservation.append({'path':p,'sha256':sha(R/p),'matches_root_observation':same})
for meta in json.loads((A/'protected-before.json').read_text()):
 if meta['path']=='book/preface.md':
  same=sha(R/meta['path'])==meta['sha256'];check(same,'Bill preface preserved')
  preservation.append({'path':meta['path'],'sha256':sha(R/meta['path']),'matches_author_initial':same})
(E/'preservation.json').write_text(json.dumps(preservation,indent=2)+'\n')
reads=['AGENTS.md','.agents/skills/ensemble-coding/SKILL.md','.agents/skills/edition-3-text-reviewer/SKILL.md','.agents/skills/edition-3-author/SKILL.md','docs/edition-3-execution.md','docs/edition-3-readiness.md','docs/edition-3-student.md','docs/edition-3-carryover.md','docs/agentic-codebook-workflow.md','docs/edition-2-carryover.md','book/voice.md','book/chapter-writing-procedure.md','book/edition-3/chapter-06.md','docs/edition-3-notes/ch06-author.md','docs/edition-3-notes/ch06-student-review.md','docs/edition-3-notes/ch06-student-review-r1.md','docs/edition-3-notes/ch06-code-review.md','docs/edition-3-notes/ch06-code-review-r1.md','docs/edition-3-notes/ch06-accepted.json','docs/edition-3-notes/ch06-entry-audit.md','docs/edition-3-notes/ch06-grader-review.md','docs/edition-3-notes/ch06-artifact-portability.md','cmd/lintprose/main.go']
(E/'read-inputs.json').write_text(json.dumps([{'path':p,'bytes':(R/p).stat().st_size,'sha256':sha(R/p),'coverage':'full read; current bytes hashed at validation'} for p in reads],indent=2)+'\n')
(E/'git-status-final.txt').write_text(run(['git','--no-pager','status','--short'])['output'])
result={'checked_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'status':'PASS' if not errors else 'FAIL','baseline_bytes':len(base.read_bytes()),'final_bytes':len(final.read_bytes()),'net_bytes':len(final.read_bytes())-len(base.read_bytes()),'added_lines':sum(1 for s in diff.splitlines() if s.startswith('+') and not s.startswith('+++')),'removed_lines':sum(1 for s in diff.splitlines() if s.startswith('-') and not s.startswith('---')),'changed_locations':changes,'author_manifest_entries':len(manifest),'local_links':len(links),'concurrent_publication_files':len(concurrent),'limits':'Text-only checks; no new code acceptance, behavioral execution or live proof. Root-owned records may be appended concurrently.','errors':errors}
print(json.dumps(result,indent=2));sys.exit(bool(errors))
