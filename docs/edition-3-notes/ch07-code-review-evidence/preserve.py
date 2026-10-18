from pathlib import Path
import json,hashlib,subprocess,datetime
w=Path('/tmp/cadence-ch07-whole-code-review-independent');r=Path('/Users/bill/projects/ensemble');s=Path('/Users/bill/projects/ensemble-edition-3-student')
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def check(p,h):assert sha(p)==h,(str(p),sha(p),h)
def git(c,*args):return subprocess.check_output(['git','--no-pager',*args],cwd=c,text=True).strip()
d=json.loads((w/'review-input.json').read_text());pins={}
for n,v in d['pins'].items():
 p=Path(n) if n.startswith('/') else r/n;check(p,v['sha256']);pins[n]=v['sha256']
assert git(s,'rev-parse','HEAD')=='48e1d12d9c396b5988c0ce6eccdef67e581872ff';assert git(s,'rev-parse','HEAD^{tree}')==d['tree'];assert git(s,'status','--porcelain')==''
assert git(s,'rev-parse','edition-3-ch06-r1^{commit}')=='f595c787675e2a6343053ec42ea6c7e1e02a930e'
src=json.loads((w/'private-source-verification.json').read_text())['files'];assert len(src)==117
for n,v in src.items():
 h=v if isinstance(v,str) else v['sha256'];check(s/n,h);check(w/'source'/n,h)
actual={str(p.relative_to(w/'source')) for p in (w/'source').rglob('*') if p.is_file()};generated=json.loads((w/'private-generated-files.json').read_text());assert actual==set(src)|set(generated);assert {n for n in generated if n.endswith('.go')}=={'gui/reviewer_contract_test.go'}
runtime=json.loads((w/'protected-runtime-source-hashes.json').read_text());assert len(runtime)==13
for n,h in runtime.items():check(r/n,h)
pubfile=Path('/tmp/cadence-ch07-publication-preservation-after-body11.json');pub=json.loads(pubfile.read_text())['files'];assert len(pub)==12
for n,v in pub.items():check(r/n,v['sha256'])
prior={}
for dirname in ['ch07-entry-evidence','ch07-grader-review-evidence','ch07-pause-diagnostic-evidence','ch07-pause-fixture-review-evidence']:
 folder=r/'docs/edition-3-notes'/dirname;manifest=folder/'SHA256SUMS';count=0;own={manifest.name}
 for line in manifest.read_text().splitlines():
  h,n=line.split('  ',1);check(r/n,h);count+=1
  prefix=str(folder.relative_to(r))+'/'
  if n.startswith(prefix):own.add(n[len(prefix):])
 assert {str(p.relative_to(folder)) for p in folder.rglob('*') if p.is_file()}==own
 prior[dirname]={'sha256sums_sha256':sha(manifest),'verified_entries':count,'exact_file_count':len(own)}
indexes={}
for rel in ['docs/edition-3-notes/ch07-native-exercise-evidence/full-review-index.json','docs/edition-3-notes/ch07-student-evidence/r4-discrepancy-package-index.json']:
 p=r/rel;j=json.loads(p.read_text());files=j['files']
 for n,v in files.items():
  candidates=[p.parent/n,p.parent.parent/n,r/n];matches=[x for x in candidates if x.is_file()];assert matches,(rel,n);q=matches[0];check(q,v if isinstance(v,str) else v['sha256'])
 indexes[rel]={'index_sha256':sha(p),'verified_files':len(files)}
out={'checked_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'main_head':git(r,'rev-parse','HEAD'),'student_head':git(s,'rev-parse','HEAD'),'student_tree':git(s,'rev-parse','HEAD^{tree}'),'student_clean':True,'source_files':117,'private_extra_authored_source':'gui/reviewer_contract_test.go','private_generated_artifact_count':len(generated)-1,'input_pins':pins,'protected_runtime_sources':runtime,'publication_baseline':str(pubfile),'publication_baseline_sha256':sha(pubfile),'publication_files':pub,'prior_audit_bundles':prior,'evidence_indexes':indexes}
(w/'final-preservation.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps({'source_files':117,'runtime':len(runtime),'publication':len(pub),'student_clean':True,'prior_bundles':prior,'indexes':indexes},indent=2))
