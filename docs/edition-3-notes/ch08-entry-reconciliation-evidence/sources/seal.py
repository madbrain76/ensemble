from pathlib import Path
import json,hashlib,os,subprocess,shutil,datetime
R=Path('/Users/bill/projects/ensemble');E=R/'docs/edition-3-notes/ch08-entry-reconciliation-evidence';T=Path((E/'scratch-path.txt').read_text().strip())
receipts=[json.loads(p.read_text()) for p in E.glob('*.receipt.json')];live=[]
for x in receipts:
 try:os.killpg(x['pid_pgid'],0);live.append(x['pid_pgid'])
 except ProcessLookupError:pass
ps=subprocess.check_output(['ps','-axo','pid=,command=']).decode().splitlines();residual=[x.strip() for x in ps if str(T) in x or str(T.resolve()) in x]
assert not live and not residual,(live,residual)
art=json.loads((E/'local-artifacts.json').read_text());bad=[]
for x in art['files']:
 b=Path(x['path']).read_bytes()
 if len(b)!=x['bytes'] or hashlib.sha256(b).hexdigest()!=x['sha256']:bad.append(x['path'])
assert not bad,bad
current=set(subprocess.check_output(['git','status','--porcelain=v1','-z','--untracked-files=all'],cwd=R).split(b'\0'));before=set((E/'status-before.z').read_bytes().split(b'\0'));removed=sorted(x.decode() for x in before-current);added=sorted(x.decode() for x in current-before)
assert not removed,removed
assert all(x[3:].startswith('docs/edition-3-notes/ch08-entry-reconciliation-evidence/') or x[3:]=='docs/edition-3-notes/ch08-entry-reconciliation.md' for x in added),added
v=json.loads((E/'verification.json').read_text());assert not any(v[k]['differences'] for k in ['inputs','student_files','protected_tracked_changes','followup_source_baseline'])
full=json.loads((E/'historical-full-corrected.stdout').read_text());assert full['passed'] and full['score']==100 and full['max_score']==100
mut=json.loads((E/'mutation-results.json').read_text());assert [(x['score'],x['passed'],x['driver_exit']) for x in mut]==[(80,True,0),(65,False,1),(80,True,0)]
state={'sealed_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'all_owned_control_jobs_ended':True,'receipt_count':len(receipts),'owned_groups_remaining':live,'owned_path_processes_remaining':residual,'local_artifact_files_verified':len(art['files']),'local_artifact_bytes':sum(x['bytes'] for x in art['files']),'local_artifact_hash_differences':bad,'scratch_retained_for_root_verification':str(T),'free_bytes':shutil.disk_usage(R).free,'resource_floor_bytes':300*1024*1024,'status_removed':removed,'status_added_at_seal':added,'status_note':'Only the new report/evidence paths differ. final-state.json, tail.txt and inventory.json are created after this comparison under the same authorized evidence directory.','full_route_score':full['score'],'full_route_max_score':full['max_score'],'full_route_passed':full['passed'],'no_additional_cases':True}
assert state['free_bytes']>state['resource_floor_bytes']
(E/'final-state.json').write_text(json.dumps(state,indent=2)+'\n')
(E/'tail.txt').write_text('SEALED: finite Chapter8 reconciliation complete.\nValid full historical route: 100/100 Passed=true exit=0.\nReal GUI persistence: old 80/80 exit=0; negative 65/80 exit=1 (settings-persist only); restored 80/80 exit=0. All compiled.\nInvalid first full invocation and pre-compile mutation-anchor failure retained and qualified.\nAll owned jobs ended; no process groups/path processes remain. Local artifacts retained for root verification.\nOriginal bundle preserved; root-restored carryover remains pinned and was not reappended. No new permission/API prerequisite imposed.\nNo additional cases; only new reconciliation report/evidence outputs.\n')
paths=[R/'docs/edition-3-notes/ch08-entry-reconciliation.md']+sorted(p for p in E.rglob('*') if p.is_file() and p!=E/'inventory.json');files=[]
for p in paths:
 b=p.read_bytes();files.append({'path':str(p.relative_to(R)),'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()})
inv={'self_excluded':'docs/edition-3-notes/ch08-entry-reconciliation-evidence/inventory.json','file_count':len(files),'total_bytes':sum(x['bytes'] for x in files),'files':files}
(E/'inventory.json').write_text(json.dumps(inv,indent=2)+'\n')
for x in files:
 b=(R/x['path']).read_bytes();assert len(b)==x['bytes'] and hashlib.sha256(b).hexdigest()==x['sha256']
print(json.dumps({'sealed':True,'inventory_files':len(files),'inventory_bytes':inv['total_bytes'],'hash_differences':0,'all_jobs_ended':True,'groups_remaining':live,'path_processes_remaining':residual,'free_bytes':state['free_bytes'],'local_artifacts_retained':len(art['files'])},indent=2))
