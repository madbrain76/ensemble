from pathlib import Path
import hashlib,json,subprocess,datetime
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');W=Path('/tmp/cadence-ch07-code-r1-rereview-independent');N=R/'docs/edition-3-notes'
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def expected(v):return v if isinstance(v,str) else v['sha256']
def git(*args,cwd=S):return subprocess.check_output(['git',*args],cwd=cwd,text=True).strip()
d=json.loads((W/'input.json').read_text());out={'checked_at':datetime.datetime.now(datetime.timezone.utc).isoformat()}
for n,v in d['pins'].items():
 p=Path(n) if Path(n).is_absolute() else R/n;assert sha(p)==expected(v),n
out['review_input_pins']=len(d['pins']);assert git('rev-parse','HEAD')==d['commit'];assert git('rev-parse','HEAD^{tree}')==d['tree'];assert git('status','--porcelain')==''
out['student']={'commit':d['commit'],'tree':d['tree'],'clean':True,'parent':git('rev-parse','HEAD^')};assert out['student']['parent']==d['parent']
sm=json.loads((W/'source-manifest.json').read_text());assert git('rev-parse','edition-3-ch06-r1^{commit}')==sm['ch06_tag_commit'];out['predecessor_tag_commit']=sm['ch06_tag_commit']
for n,v in sm['source_files'].items():
 assert sha(S/n)==expected(v),n;assert sha(W/'source'/n)==expected(v),n
out['source_original_members']=len(sm['source_files']);out['all_source_hashes_unchanged']=True
em=json.loads((W/'evidence-manifest.json').read_text());assert len(em['files'])==138
for n,v in em['files'].items():assert sha(N/n)==expected(v),n
out['r1_portable_files_verified']=len(em['files'])
runtime=json.loads((W/'runtime-before.json').read_text());assert len(runtime)==13
for n,v in runtime.items():assert sha(R/n)==expected(v),n
out['runtime_13']=runtime
pub=json.loads(Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch07-code-review-evidence/final-preservation.json').read_text())['publication_files'];assert len(pub)==12
for n,v in pub.items():assert sha(R/n)==expected(v),n
out['publication_12']={n:expected(v) for n,v in pub.items()}
old=[]
for name in ['ch07-entry-evidence','ch07-grader-review-evidence','ch07-pause-diagnostic-evidence','ch07-pause-fixture-review-evidence','ch07-code-review-evidence']:
 base=N/name;mf=base/'SHA256SUMS';count=0;declared={str(mf.relative_to(R))}
 for line in mf.read_text().splitlines():
  h,n=line.split('  ',1);assert sha(R/n)==h,n;count+=1
  if n.startswith(str(base.relative_to(R))+'/'):declared.add(n)
 actual={str(p.relative_to(R)) for p in base.rglob('*') if p.is_file()};assert actual==declared,name
 old.append({'directory':name,'manifest_sha256':sha(mf),'verified_manifest_entries':count,'exact_set':len(actual)})
out['prior_audit_bundles']=old
# Preserve all indexed previous student/native records, without copying their payloads.
for filename in ['ch07-evidence-manifest.json','ch07-native-exercise-evidence/full-review-index.json']:
 m=json.loads((N/filename).read_text());fs=m['files'];
 for n,v in fs.items():assert sha(N/n)==expected(v),n
 out[filename]={'sha256':sha(N/filename),'indexed_members':len(fs)}
orig=W/'original-r4/gui/web/gui';cur=W/'source/gui/web/gui';assets={}
for n in ['speech.js','queue.js','artifact-scroll.js','format.js','index.html','style.css']:
 assert (orig/n).read_bytes()==(cur/n).read_bytes();assets[n]=sha(cur/n)
anchor="$('draft').addEventListener('input', localPause);";a=(orig/'app.js').read_text().split(anchor,1)[1];b=(cur/'app.js').read_text().split(anchor,1)[1];assert a==b
out['unchanged_native_mechanism_assets']=assets;out['unchanged_C1_control_tail_sha256']=hashlib.sha256((anchor+b).encode()).hexdigest();out['main_head']=git('rev-parse','HEAD',cwd=R)
extras={}
for base,label,known in [(W/'source','corrected',set(sm['source_files'])),(W/'original-r4','original',set(json.loads((W/'original-r4-source-pins.json').read_text())))]:
 extra={str(p.relative_to(base)):{'bytes':p.stat().st_size,'sha256':sha(p)} for p in base.rglob('*') if p.is_file() and str(p.relative_to(base)) not in known};extras[label]=extra
out['private_extra_file_counts']={n:len(x) for n,x in extras.items()};(W/'private-extra-file-inventory.json').write_text(json.dumps(extras,indent=2)+'\n');(W/'preservation.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps({k:v for k,v in out.items() if k not in ['runtime_13','publication_12','unchanged_native_mechanism_assets']},indent=2))
