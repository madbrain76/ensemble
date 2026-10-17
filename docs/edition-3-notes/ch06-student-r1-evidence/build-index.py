#!/usr/bin/env python3
"""Index only the explicit R1 outputs/new temp root; archives are always local-only."""
import hashlib,json,pathlib,subprocess
B=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes')
E=B/'ch06-student-r1-evidence';T=pathlib.Path('/tmp/e3-ch06-correction-r1')
S=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
manifest=B/'ch06-student-r1-source-manifest.json';M=json.loads(manifest.read_text())
assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=S,text=True).strip()==M['commit']
assert subprocess.check_output(['git','status','--porcelain'],cwd=S)==b''
receipts=json.loads((E/'final-command-receipts.json').read_text())
assert len(receipts)==13 and all(r['exit_code']==0 and not r['timed_out'] for r in receipts)
assert [r['passed'] for r in receipts if r['name'].startswith('grade-')]==[True,True]
artifacts=[]
def entry(p,portable,kind):
 assert not p.is_symlink()
 data=p.read_bytes()
 if portable:data.decode('utf-8')
 return {'path':str(p),'bytes':len(data),'sha256':hashlib.sha256(data).hexdigest(),'portable':portable,'kind':kind}
for p in sorted(E.rglob('*')):
 if p.is_file():
  portable=p.suffix.lower() in {'.txt','.json','.md','.py','.diff','.patch'} and '__pycache__' not in p.parts
  artifacts.append(entry(p,portable,'safe R1 student receipt/source/fixture' if portable else 'local-only artifact'))
for p in [B/'ch06-student-review-r1.md',manifest]:artifacts.append(entry(p,True,'R1 report/frozen source manifest'))
for p in sorted(T.rglob('*')):
 if p.is_file():artifacts.append(entry(p,False,'local-only frozen archive/executable or private local runtime; do not redistribute'))
paths=[x['path'] for x in artifacts];assert len(paths)==len(set(paths))
by_path={x['path']:x for x in artifacts}
for pin in M['artifacts']:
 x=by_path[pin['path']];assert not x['portable'] and x['bytes']==pin['bytes'] and x['sha256']==pin['sha256']
for r in receipts:
 for name in ['stdout','stderr','exit_receipt']:
  pin=r[name];x=by_path[pin['path']];assert x['bytes']==pin['bytes'] and x['sha256']==pin['sha256']
result={'schema':1,'correction_round':1,'status':'Completed R1 candidate for independent review; no acceptance or tag','source_commit':M['commit'],'source_tree':M['tree'],'source_clean':True,'original_candidate':'800da21d215b9ca48de41516be4773962fe38ccd','original_evidence':'Separate, immutable; original live failure and all paid receipts retained, not fresh evidence on R1','live_calls_r1':0,'credential_reads_r1':0,'owned_jobs_pending':0,'completed_final_commands':len(receipts),'source_manifest_sha256':hashlib.sha256(manifest.read_bytes()).hexdigest(),'counts':{'total':len(artifacts),'portable':sum(x['portable'] for x in artifacts),'local_only':sum(not x['portable'] for x in artifacts)},'self_excluded':True,'artifacts':artifacts}
path=B/'ch06-student-r1-evidence-manifest.json';path.write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps({'manifest':str(path),'sha256':hashlib.sha256(path.read_bytes()).hexdigest(),'bytes':path.stat().st_size,'counts':result['counts']},indent=2))
