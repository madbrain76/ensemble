import json,hashlib,subprocess,os,shutil,datetime
from pathlib import Path
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';E=N/'ch08-code-review-r1-evidence';T=Path((E/'scratch-path.txt').read_text().strip());report=N/'ch08-code-review-r1.md';sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
B=json.loads((E/'preservation-baseline.json').read_text());V=json.loads((E/'final-verification.json').read_text());assert sha(N/'ch08-code-review-r1-input.json')==V['input_manifest_sha256']
def check(x,base=None):
 p=Path(x['path']);p=p if p.is_absolute() else base/p;assert p.stat().st_size==x['bytes'] and sha(p)==x['sha256'],p
for x in B['inputs']+list(B['protected_tracked'].values()):check(x)
for x in B['source_members']:check(x,S)
for x in B['student_portable']['entries']:check(x,N)
for x in B['student_portable']['local_only_exclusions']+B['source_manifest']['unchanged_initial_artifacts']:check(x)
for x in json.loads((N/'ch08-code-review-evidence/inventory.json').read_text())['files']:check(x,R)
for x in json.loads((E/'carried-controls.json').read_text())['runtime_identities']:check(x)
for x in json.loads((E/'local-artifacts.json').read_text())['files']:check(x)
assert subprocess.check_output(['git','--no-pager','status','--porcelain=v1','--untracked-files=all'],cwd=S)==b''
assert subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=S,text=True).strip()==V['student']['commit']
assert subprocess.check_output(['git','--no-pager','show-ref','--tags'],cwd=S,text=True)==(T/'student-tags-before.txt').read_text()
C=json.loads((E/'controls-summary.json').read_text());assert len(C['commands'])==10
for r in C['commands']:
 assert r['exit_code']==r['expected_exit'] and not r['remaining_process_group'] and not r['timed_out']
 for stream in ['stdout','stderr']:assert sha(E/(r['name']+'.'+stream+'.log'))==r[stream+'_sha256']
 try:os.killpg(r['pid_pgid'],0);raise AssertionError(('owned group remains',r['pid_pgid']))
 except ProcessLookupError:pass
try:os.killpg(61610,0);raise AssertionError('pager/read group remains')
except ProcessLookupError:pass
for filename in ['inventory.json','final-seal.json']:
 p=E/filename
 if not p.exists():p.write_text('{}\n')
def status(data):
 parts=data.split(b'\0');out={};i=0
 while i<len(parts):
  x=parts[i];i+=1
  if not x:continue
  code=x[:2].decode();path=x[3:].decode();out[path]=code
  if code[0] in 'RC' or code[1] in 'RC':
   assert i<len(parts);out['rename-source:'+parts[i].decode()]=code;i+=1
 return out
before=status((T/'main-status-before.z').read_bytes());after=status(subprocess.check_output(['git','--no-pager','status','--porcelain=v1','-z','--untracked-files=all'],cwd=R));changes=[{'path':k,'before':before.get(k),'after':after.get(k)} for k in sorted(before.keys()|after.keys()) if before.get(k)!=after.get(k)];allowed=lambda p:p=='docs/edition-3-notes/ch08-code-review-r1.md' or p.startswith('docs/edition-3-notes/ch08-code-review-r1-evidence/');outside=[x for x in changes if not allowed(x['path'])];assert not outside,outside
(E/'status-scope.json').write_text(json.dumps({'changes':changes,'outside_scope_changes':outside,'qualification':'Status before/after plus exact protected/source/input/prior-seal hashes; only new allowed paths written.'},indent=2)+'\n')
seal={'checked_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'verdict':'PASS — fixture-only R1 resolved','student_commit':V['student']['commit'],'student_clean_and_tags_unchanged':True,'all_input_source_protected_and_carried_runtime_hashes_match':True,'owned_command_receipts':10,'unexpected_runtime_failures':0,'expected_behavioral_negatives':1,'remaining_owned_process_groups':0,'pager_read_group_ended':True,'owned_copy_restored':True,'fresh_go_suite_runs':0,'fresh_grader_or_lifetime_runs':0,'outside_scope_status_changes':0,'old_reports_and_manifests_unchanged':True,'archive_and_raw_wire_payloads_not_exported':True,'free_bytes':shutil.disk_usage(R).free}
(E/'final-seal.json').write_text(json.dumps(seal,indent=2)+'\n')
entries=[]
for p in [report]+sorted(p for p in E.rglob('*') if p.is_file() and p!=E/'inventory.json'):
 assert not p.is_symlink(),p;data=p.read_bytes();assert b'\0' not in data,p;data.decode('utf-8');assert data[:4] not in [b'\x7fELF',b'\xcf\xfa\xed\xfe',b'\xfe\xed\xfa\xcf',b'PK\x03\x04'],p;assert data[257:262]!=b'ustar',p
 rel=str(p.relative_to(R));assert allowed(rel),rel;entries.append({'path':rel,'bytes':len(data),'sha256':hashlib.sha256(data).hexdigest(),'classification':'portable UTF-8 review/source/evidence metadata; no binary or raw-wire payload'})
assert len(entries)==len({x['path'] for x in entries})
inv={'format':'ch08-code-r1-independent-review-v1','self_excluded':True,'files':entries,'file_count':len(entries),'total_bytes':sum(x['bytes'] for x in entries),'local_payloads':'separately indexed in local-artifacts.json; raw wire and TARs remain at pinned local-only locations, not copied here'}
(E/'inventory.json').write_text(json.dumps(inv,indent=2)+'\n')
for x in json.loads((E/'inventory.json').read_text())['files']:check(x,R)
print(json.dumps({'seal':seal,'portable_file_count':inv['file_count'],'portable_bytes':inv['total_bytes'],'report_sha256':sha(report),'inventory_sha256':sha(E/'inventory.json'),'verified_inventory_differences':0},indent=2))
