import pathlib,json,hashlib,subprocess,datetime
R=pathlib.Path('/Users/bill/projects/ensemble'); S=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student'); T=pathlib.Path('/tmp/ch08-entry-audit.aTPJZB')
m=json.loads((R/'docs/edition-3-notes/ch08-entry-manifest.json').read_text())
def git(root,*args): return subprocess.check_output(['git','--no-pager','-C',str(root),*args]).decode().strip()
def digest(p):
 b=p.read_bytes(); return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
report={'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'identities':{},'hash_groups':{}}
for key,root,args in [('main_head',R,['rev-parse','HEAD']),('student_head',S,['rev-parse','HEAD']),('student_tree',S,['rev-parse','HEAD^{tree}']),('student_tag_commit',S,['rev-parse','edition-3-ch07-r1^{}'])]:
 actual=git(root,*args); report['identities'][key]={'actual':actual,'matches':actual==m['student_head' if key=='student_tag_commit' else key]}
for key,root in [('inputs',R),('student_files',S),('protected_tracked_changes',R)]:
 diffs=[]
 for f,want in m[key].items():
  actual=digest(root/f) if (root/f).is_file() else None
  if actual!=want: diffs.append(f)
 report['hash_groups'][key]={'checked':len(m[key]),'differences':diffs}
report['student_status']=git(S,'status','--porcelain=v1','--untracked-files=all')
protected={}
for prefix in ['internal/grade','internal/gradeedition2','internal/gradeedition3','cmd/grade-edition2','cmd/grade-edition3','solutions/ch08','solutions/edition-3']:
 for p in sorted((R/prefix).rglob('*')):
  if p.is_file(): protected[str(p.relative_to(R))]=digest(p)
baseline=T/'protected-baseline.json'
if baseline.exists():
 old=json.loads(baseline.read_text()); report['owned_baseline']={'checked':len(old),'differences':[k for k in old.keys()|protected.keys() if old.get(k)!=protected.get(k)]}
else: baseline.write_text(json.dumps(protected,indent=2)); report['owned_baseline']={'recorded':len(protected)}
print(json.dumps(report,indent=2))
