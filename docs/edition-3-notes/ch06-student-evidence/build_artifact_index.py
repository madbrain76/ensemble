#!/usr/bin/env python3
"""Index only this student's allowed evidence and explicitly referenced own runtimes."""
import hashlib,json,pathlib,subprocess,sys
BASE=pathlib.Path(__file__).resolve().parent.parent
E=BASE/'ch06-student-evidence';S=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
DEST=BASE/'ch06-student-evidence-manifest.json'
def info(p,portable,kind):
 return {'path':str(p),'relative_path':str(p.relative_to(BASE)) if p.is_relative_to(BASE) else None,'size_bytes':p.stat().st_size,'sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'portable':portable,'kind':kind}
files=[BASE/'ch06-student-review.md',BASE/'ch06-student-source-manifest.json']+[p for p in E.rglob('*') if p.is_file() and not p.is_symlink()]
entries=[];private=set()
def refs(v):
 if isinstance(v,dict):
  for k,x in v.items():
   if k=='private_directory' and isinstance(x,str):
    p=pathlib.Path(x)
    # Only explicit directories created/reported by these student launchers.
    if any(part.startswith('e3-ch06-') and 'private-' in part for part in p.parts):private.add(p)
   else:refs(x)
 elif isinstance(v,list):
  for x in v:refs(x)
for p in sorted(set(files)):
 archive=p in {E/'candidate-source.tar',E/'candidate1/source.tar'}
 portable='__pycache__' not in p.parts and not archive
 kind='local-only source archive; do not redistribute' if archive else ('portable student evidence/source' if portable else 'generated local-only Python cache')
 entries.append(info(p,portable,kind))
 if p.suffix=='.json':
  try:refs(json.loads(p.read_text()))
  except (ValueError,UnicodeError):pass
for directory in sorted(private):
 if directory.is_dir():
  for p in sorted(directory.rglob('*')):
   if p.is_file() and not p.is_symlink():entries.append(info(p,False,'own private runtime: raw trace/history/output; do not redistribute'))
for name in ['candidate.tar','frozen-core','frozen-consumer','frozen-workflow','frozen-shape','candidate1-candidate.tar','root-api-trace.log']:
 p=pathlib.Path('/tmp/e3-ch06-'+name)
 if p.is_file():entries.append(info(p,False,'local-only source archive/executable or owned fixture trace; archive copies also local-only'))
# Retained candidate-1 executables are enumerated by the preserved artifact map,
# not guessed filenames; assert their existing pinned bytes before indexing.
expected={f'e3-ch06-candidate1-frozen-{name}' for name in ['core','consumer','workflow','shape']}
found=set()
for artifact in json.loads((E/'candidate1/preserved-artifacts.json').read_text()):
 p=pathlib.Path(artifact['path'])
 if p.name not in expected:continue
 assert p.is_file() and p.stat().st_size==artifact['bytes']
 assert hashlib.sha256(p.read_bytes()).hexdigest()==artifact['sha256']
 assert not any(entry['path']==str(p) for entry in entries)
 entries.append(info(p,False,'retained candidate-1 executable; local-only'))
 found.add(p.name)
assert found==expected
# Include committed safe test/control receipts without duplicating the entire source archive.
for p in sorted((S/'evidence/ch06').rglob('*')):
 if p.is_file() and not p.is_symlink():entries.append(info(p,True,'committed safe local test/control receipt; also in source archive'))
def git(*a):return subprocess.check_output(['git','-C',str(S),*a],text=True).strip()
assert git('rev-parse','HEAD')=='800da21d215b9ca48de41516be4773962fe38ccd'
assert git('status','--porcelain')==''
result={'schema':1,'role':'student evidence, not independent acceptance','self_excluded':str(DEST),'source_commit':git('rev-parse','HEAD'),'source_tree':git('rev-parse','HEAD^{tree}'),'source_clean':True,'local_and_public_gates':'all required final-byte gates passed; see exact paired-grade command receipt','original_live':'failed first Anthropic assertion; three HTTP200 calls, remaining original cases not run','remaining_live':'see remaining-live-summary and command exit if present; absence means unrun, not success','completion_note':sys.argv[1] if len(sys.argv)>1 else 'index draft; not a completion assertion','private_policy':'Only hashes, sizes and paths exported for own raw private artifacts. Never redistribute those bytes. Settings credentials are not indexed or reread.','artifacts':entries}
result['inventory_correction']={'scope':'Metadata packaging only; no source/artifact byte changes, tests, builds, provider calls or credential reads. Two source archive copies reclassified local-only and four existing candidate-1 executables added local-only from their preserved map.','retained_prior_boundary':'/tmp/cadence-ch06-pre-inventory-correction (parent-held; untouched)','prior_report_sha256':'7022fcc7331d04a7a6cd21f770aebe86e81e35234612384b82602b8261252429','prior_index_sha256':'a45b8c21ab677dc892622ab522dcf452e95c9d0b0ac9b1051bf622491238f0c4','prior_counts':{'total':470,'portable':285,'local_only':185}}
result['artifact_counts']={'total':len(entries),'portable':sum(x['portable'] for x in entries),'local_only':sum(not x['portable'] for x in entries)}
DEST.write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps({'manifest':str(DEST),'sha256':hashlib.sha256(DEST.read_bytes()).hexdigest(),'size_bytes':DEST.stat().st_size,'portable_artifacts':sum(x['portable'] for x in entries),'local_only_artifacts':sum(not x['portable'] for x in entries)},indent=2))
