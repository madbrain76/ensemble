#!/usr/bin/env python3
import hashlib,json,pathlib,re,subprocess,tarfile
P=pathlib.Path('/tmp/cadence-ch05-correction-r1-evidence');N=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes');D=N/'ch05-student-r1-evidence';R=N/'ch05-student-review-r1.md';I=N/'ch05-student-r1-evidence-manifest.json';S=N/'ch05-student-r1-source-manifest.json'
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
identity=json.loads((P/'frozen/identity.json').read_text());root=pathlib.Path(identity['source_root'])
assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip()==identity['source_commit'];assert not subprocess.check_output(['git','status','--porcelain'],cwd=root)
for name,digest in identity['source_files'].items():assert sha(root/name)==digest,name
assert sha(pathlib.Path(identity['archive']['path']))==identity['archive']['sha256']
with tarfile.open(identity['archive']['path']) as archive:
 members={m.name:m for m in archive.getmembers() if m.isfile()};assert set(members)==set(identity['source_files'])
 for name,digest in identity['source_files'].items():assert hashlib.sha256(archive.extractfile(members[name]).read()).hexdigest()==digest
for item in identity['binaries'].values():assert sha(pathlib.Path(item['path']))==item['sha256']
state={'source_commit':identity['source_commit'],'source_tree':identity['source_tree'],'worktree_clean':True,'archive_verified_files':len(identity['source_files']),'binaries_verified':identity['binaries'],'paid_calls_this_round':0,'children_spawned':0,'acceptance_tag_created':False}
(P/'final-state.json').write_text(json.dumps(state,indent=2)+'\n')
files={};excluded={};suspected=[];pattern=re.compile(rb'(?:sk-ant-[A-Za-z0-9_-]{20,}|sk-proj-[A-Za-z0-9_-]{20,}|AIza[0-9A-Za-z_-]{30,})')
for path in sorted(P.rglob('*')):
 if not path.is_file():continue
 rel=str(path.relative_to(P));raw=path.read_bytes();meta={'sha256':hashlib.sha256(raw).hexdigest(),'bytes':len(raw)}
 binary=path.suffix.lower() in ['.tar','.gz','.png','.pyc','.pdf','.mp4','.jpg','.wav','.mp3'] or b'\0' in raw
 if not binary:
  try:raw.decode('utf-8')
  except UnicodeError:binary=True
 if binary:excluded[rel]={**meta,'local_path':str(path),'reason':'binary/archive/media/cache stays outside Git'};continue
 if pattern.search(raw):suspected.append(rel)
 files[rel]=meta
assert not suspected,('possible token-shaped data; copying stopped',suspected)
if D.exists():
 existing={str(p.relative_to(D)) for p in D.rglob('*') if p.is_file()};assert existing<=set(files),'unexpected destination files; not overwriting'
 for rel in existing:assert sha(D/rel)==files[rel]['sha256'],rel
D.mkdir(exist_ok=True)
for rel,meta in files.items():
 target=D/rel;target.parent.mkdir(parents=True,exist_ok=True);target.write_bytes((P/rel).read_bytes());assert sha(target)==meta['sha256']
source_bytes=(P/'frozen/identity.json').read_bytes()
if S.exists():assert S.read_bytes()==source_bytes
S.write_bytes(source_bytes)
data={'schema':1,'portable_directory':D.name,'text_files':files,'excluded_artifacts':excluded,'report':{'relative_path':R.name,'sha256':sha(R),'bytes':R.stat().st_size},'source_index':{'relative_path':S.name,'sha256':sha(S)},'candidate':identity,'privacy_checks':{'possible_token_shaped_matches':0,'provider_calls_this_round':0,'credential_loader_used':False,'new_runtime_histories_are_local_fake_provider_only':True},'original_preservation_receipt':{'relative_path':'original-preservation.json','sha256':files['original-preservation.json']['sha256']},'evidence_credit':{'original_compile_exit':0,'original_behavioral_controls':{'R1':1,'R2':2,'R3_failing_layouts':7,'R3_legitimate_passing_layouts':1},'corrected_focused_exit':0,'all_frozen_local_gates_exit':0,'coverage_percent':86.2,'actual_executable_local_regressions_exit':0,'public_chapter5':{'exit':0,'score':100,'elapsed_seconds':159.073},'public_chapter1':{'exit':0,'score':100,'elapsed_seconds':127.260,'headers_delay_seconds':62,'body_delay_seconds':62},'setup_failures_count_as_detections':False,'initial_live_evidence_is_not_corrected_live_evidence':True},'relocation':'verify-portable.py checks relative copied texts plus adjacent report/source index; binaries and archives remain outside Git at indexed paths'}
I.write_text(json.dumps(data,indent=2)+'\n')
print(json.dumps({'copied_text_files':len(files),'excluded_artifacts':len(excluded),'source_commit':identity['source_commit'],'report_sha256':sha(R),'source_index_sha256':sha(S),'evidence_index_sha256':sha(I),'original_preservation_passed':json.loads((P/'original-preservation.json').read_text())['passed']},indent=2))
