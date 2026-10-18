from pathlib import Path
import hashlib,json,shutil,subprocess,datetime
root=Path('/Users/bill/projects/ensemble-edition-3-student');local=root/'.ch08-local/operational-prep';notes=Path('/Users/bill/projects/ensemble/docs/edition-3-notes');out=notes/'ch08-operational-plan-evidence';plan=notes/'ch08-operational-plan.md'
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest();load=lambda p:json.loads(p.read_text())
def info(p):return {'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)}
def write(p,x):p.write_text(json.dumps(x,indent=2)+'\n')
assert not out.exists() and not plan.exists()
r1src=notes/'ch08-student-r1-source-manifest.json';r1inv=notes/'ch08-student-r1-evidence-manifest.json';s=load(r1src)
assert subprocess.check_output(['git','--no-pager','status','--porcelain=v1'],cwd=root,text=True)==''
assert subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=root,text=True).strip()==s['candidate']['commit']=='b8f22ba5b3c9d1fc52d1a114810b90b77928684d'
assert subprocess.check_output(['git','--no-pager','rev-parse','HEAD^{tree}'],cwd=root,text=True).strip()==s['candidate']['tree']
for r in s['tracked']:assert sha(root/r['path'])==r['sha256']
for r in s['unchanged_initial_artifacts']:assert sha(Path(r['path']))==r['sha256']
old=[]
for f in [notes/'ch08-student-evidence-manifest.json',r1inv]:
 x=load(f)
 for row in x['entries']:assert sha(notes/row['path'])==row['sha256']
 old.append(info(f))
for f in [notes/'ch08-student-source-manifest.json',r1src]:old.append(info(f))
public=load(root/'.ch08-local/public-input-identity.json')
for row in public['files']:assert sha(Path(row['path']))==row['sha256']
assignment=Path('/tmp/cadence-ch08-operational-plan-task.txt');assert assignment.stat().st_size==6418 and sha(assignment)=='3bba619428affbee88c709122c903284f0a2e47835d6ce1608342960a0f06721'
p=load(local/'policy.json');assert sha(local/'launcher.py')==p['launcher_sha256'];assert sha(local/'policy.json')==(local/'policy.sha256').read_text().strip()
for path,digest in p['scripts'].items():assert sha(Path(path))==digest
last=load(local/'attempt-final/receipt.json');assert last['passed'] and last['all_owned_processes_stopped'];assert last['policy_sha256']==sha(local/'policy.json')
for d in ['attempt-1','attempt-2','attempt-final']:
 assert load(local/d/'receipt.json')['all_owned_processes_stopped']
 exits=list((local/d).glob('exit-*.json'));assert len(exits)==1
 for f in exits:assert load(f)['exit_code']==0 and load(f).get('reader_finished',True)
out.mkdir();shutil.copy2(local/'plan.md',plan)
for name in ['launcher.py','rehearse.py','policy.json','policy.sha256','launcher-draft-before-clock-review.py','seal.py']:shutil.copy2(local/name,out/name)
for name in ['attempt-1-scripts','attempt-2-scripts']:shutil.copytree(local/name,out/name)
for name in ['attempt-1','attempt-2','attempt-final']:
 dst=out/name;dst.mkdir()
 for f in (local/name).iterdir():
  if f.is_file() and (f.name in {'receipt.json','wrongpin.stdout','wrongpin.stderr','assets-1.json','ownership-1.json'} or f.name.startswith('exit-')):shutil.copy2(f,dst/f.name)
for f in sorted((root/'.ch08-local').glob('operational-prep-rehearsal*')):
 if f.suffix in {'.json','.log'}:shutil.copy2(f,out/f.name)
shutil.copy2(root/'.ch08-local/run_gate.py',out/'receipt-wrapper.py')
write(out/'candidate-carry.json',{'candidate':s['candidate'],'clean':True,'all_tracked_files_verified_unchanged':s['tracked'],'all_retained_executables_verified_unchanged':s['unchanged_initial_artifacts'],'previous_immutable_manifests':old,'previous_package_entries_verified_unchanged':True,'no_rebuild_or_new_source_archive':True,'root_status_as_reported':'R1 independently PASS and root code-review accepted; no overall acceptance/tag','current_source_archive':'unchanged previous local-only archive; not copied'})
write(out/'input-pins.json',{'assignment':info(assignment),'retained_public_reads':'Earlier full prerequisites qualified by unchanged exact pins; predecessor selected sections remain section reads.','public':public,'no_private_review_inputs':True})
write(out/'local-only-index.json',{'scope':'owned preparation files; raw fake-provider/API/history/GUI logs remain local and are NOT portable evidence','files':[info(f) for f in sorted(local.rglob('*')) if f.is_file()]})
write(out/'environment.json',{'python':subprocess.check_output(['python3','--version'],text=True).strip(),'platform':subprocess.check_output(['uname','-sm'],text=True).strip(),'receipt_environments':'exact argv/cwd/environment names in per-attempt receipts; GUI got isolated HOME and dummy Anthropic key only','paid_and_browser_native_paths_executed':False,'children_created':False})
(out/'preparation-notes.txt').write_text('All execution attempts retained. Draft before first run corrected local-vs-paid clocks and added body hashes. Attempt 1 failed on /index.html redirect locator; first diagnostic lacked asset name. Attempt 2 passed after / route repair and per-asset receipts. Final revision added socket abort before exceptional server teardown; same bounded controls passed. No production bug or native/provider result claimed. One guessed inspection filename internal/engine/anthropic.go was absent; actual internal/engine/engine.go was then inspected. No private review or historical solution links followed. No wrapper/native interception was needed: actual production speechEvents includes rate getter/requested and callbacks. No script or approval file was used to execute external/browser/native paths. Root controls later permission and attachment.\n')
entries=[]
for f in sorted([plan]+[f for f in out.rglob('*') if f.is_file()]):
 b=f.read_bytes();assert b'\0' not in b;b.decode('utf-8');assert f.suffix not in {'.tar','.zip','.gz','.pyc'};entries.append({'path':str(f.relative_to(notes)),'bytes':len(b),'sha256':sha(f)})
manifest=out/'MANIFEST.json'
write(manifest,{'schema':'chapter8-operational-preparation-v1','status':'SEALED PREPARATION ONLY; permission and browser/native/provider evidence still required','created_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'candidate_commit':p['commit'],'candidate_tree':p['tree'],'gui_sha256':p['gui_sha256'],'current_policy_sha256':sha(local/'policy.json'),'current_launcher_sha256':sha(local/'launcher.py'),'current_rehearsal_sha256':sha(local/'rehearse.py'),'portable_utf8_non_nul':True,'entries':entries,'all_portable_files_including_manifest':sorted([r['path'] for r in entries]+[str(manifest.relative_to(notes))]),'self_hash':'Excluded recursively; final exact digest reported out-of-band.','no_binary_archive_pycache_raw_api_private_material_credentials':True,'sealed_no_append':True})
print('SEALED',len(entries)+1,'portable files')
for f in [plan,manifest]:print(f,sha(f))
print('CURRENT POLICY',sha(local/'policy.json'),'CURRENT LAUNCHER',sha(local/'launcher.py'))
