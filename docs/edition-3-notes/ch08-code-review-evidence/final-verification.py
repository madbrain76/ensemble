import sys,json,hashlib,os,subprocess,shutil,datetime
from pathlib import Path
sys.dont_write_bytecode=True
R=Path('/Users/bill/projects/ensemble');N=R/'docs/edition-3-notes';E=N/'ch08-code-review-evidence';S=Path('/Users/bill/projects/ensemble-edition-3-student')
sha=lambda p:hashlib.sha256(Path(p).read_bytes()).hexdigest()
M=json.loads((N/'ch08-code-review-input.json').read_text());B=json.loads((E/'input-verification-before.json').read_text());assert sha(N/'ch08-code-review-input.json')==B['input_manifest_sha256']
def verify(p,x):
 assert p.is_file(),str(p)
 assert p.stat().st_size==x['bytes'] and sha(p)==x['sha256'],str(p)
for x in M['inputs']:verify(Path(x['path']),x)
for x in M['source_members']:verify(S/x['path'],x)
for x in M['protected_tracked'].values():verify(Path(x['path']),x)
def git(*av):return subprocess.check_output(['git',*av],cwd=S).decode().strip()
assert git('rev-parse','HEAD')==M['student']['commit'];assert git('rev-parse','HEAD^{tree}')==M['student']['tree'];assert not git('status','--porcelain=v1','--untracked-files=all')
assert git('rev-parse',M['predecessor']['tag'])==M['predecessor']['tag_object'];assert git('rev-parse',M['predecessor']['tag']+'^{}')==M['predecessor']['commit']
P=json.loads((N/'ch08-student-evidence-manifest.json').read_text())
for x in P['entries']:verify(N/x['path'],x)
ref=P['local_only_index'];lp=Path(ref if isinstance(ref,str) else ref['path']);lp=lp if lp.is_absolute() else N/lp;L=json.loads(lp.read_text())
localroots=[S,S/'.ch08-local',lp.parent,N]
for k,v in L.items():
 if isinstance(v,str) and ('root' in k or 'directory' in k) and Path(v).is_dir():localroots.insert(0,Path(v))
for x in L['files']:
 q=Path(x['path']);candidates={str((base/q).resolve()) for base in localroots if (base/q).is_file()} if not q.is_absolute() else {str(q.resolve())};assert len(candidates)==1,(x['path'],candidates);verify(Path(candidates.pop()),x)
A=json.loads((E/'artifact-verification.json').read_text())
for x in A['six_executables']:verify(Path(x['path']),x)
prior=[]
for rel in ['ch08-grader-review-evidence/inventory.json','ch08-entry-reconciliation-evidence/inventory.json']:
 d=json.loads((N/rel).read_text())
 for x in d['files']:verify(R/x['path'],x)
 prior.append({'index':str(N/rel),'members':len(d['files']),'all_hashes_match':True})
records=[]
for p in sorted(E.glob('*.receipt.json')):
 r=json.loads(p.read_text());assert not r['timed_out'] and not r['remaining_process_group'],p.name
 try:os.killpg(r['pid_pgid'],0);raise AssertionError(('remaining group',r['pid_pgid']))
 except ProcessLookupError:pass
 assert sha(E/(r['name']+'.stdout.log'))==r['stdout_sha256'];assert sha(E/(r['name']+'.stderr.log'))==r['stderr_sha256'];records.append(r)
assert len([r for r in records if r['name']=='chapter1-lifetime-once'])==1
actual_failures=[r['name'] for r in records if r['exit_code']!=r['expected_exit']]
assert actual_failures==['page-inherited-integration'],actual_failures
for n in ['chapter8','chapter1']:
 d=json.loads((E/(n+'-report.json')).read_text());assert d['passed'] and d['score']==d['max_score']==100
T=Path((E/'scratch-path.txt').read_text().strip());base=json.loads((E/'owned-copy-before.json').read_text())
for x in base:assert sha(T/'candidate-copy'/x['path'])==x['sha256']
local=[]
for p in sorted(T.rglob('*')):
 if p.is_file():local.append({'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p),'classification':'local-only owned scratch; binary/process/tool/source artifacts not portable evidence'})
(E/'local-artifacts.json').write_text(json.dumps({'root':str(T),'retained':True,'files':local,'count':len(local),'bytes':sum(x['bytes'] for x in local),'remaining_process_groups':0},indent=2)+'\n')
(E/'controls-summary.json').write_text(json.dumps({'commands':records,'actual_unexpected_failures':actual_failures,'expected_mutant_failures':['budget-negative-behavior','speech-negative-behavior'],'chapter1_full_runs':1,'all_process_groups_ended':True},indent=2)+'\n')
V={'checked_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'candidate':M['student'],'predecessor':M['predecessor'],'input_manifest_sha256':sha(N/'ch08-code-review-input.json'),'inputs_verified':len(M['inputs']),'source_files_verified':len(M['source_members']),'protected_files_verified':len(M['protected_tracked']),'student_clean':True,'indexed_student_members_verified':len(P['entries']),'local_student_members_verified':len(L['files']),'six_runtime_identities_unchanged':True,'prior_bundles':prior,'original_fixture_failure_retained':True,'owned_copy_restored':True,'owned_command_receipts':len(records),'remaining_process_groups':0,'free_bytes':shutil.disk_usage(R).free,'local_retained_files':len(local),'local_retained_bytes':sum(x['bytes'] for x in local),'archive_qualification':M['archive_classification'],'no_source_guide_old_report_edits':True}
(E/'final-verification.json').write_text(json.dumps(V,indent=2)+'\n');print(json.dumps(V,indent=2))
notes={'control_plan_preparation_failure':json.loads((E/'fresh-preparation-failure.json').read_text()),'read_only_lookup_errors':'Guessed source/evidence names were absent in several read-only commands; actual paths were resolved from tracked inventory and authoritative manifests. No guessed source or receipt was substituted. Examples include nonexistent internal/llm paths (implementation is internal/engine), gui/cmd/gui/main.go (actual gui/cmd/main.go), guessed local-artifact JSON paths and .tar.gz (actual source.tar).','artifact_summary_assertion':'A six-build assertion initially counted all nine retained student build records; corrected selection is five current executables plus grader. Archive membership was already correct at 126 files/26 directories; an early visible suspicion of archive-count mismatch was corrected.','summary_schema_error':'A read-only summary tried process_group_alive_after; actual receipt key is remaining_process_group. Original receipts remain intact.','fixture_failure':'Checked-in correction-r1-page.mjs fails at missing window.matchMedia before behavioral assertions; source and failing logs retained. An owned diagnostic adapts only inert boundaries/locators/settings cancellation and passes; it is not a candidate fix.','not_executed_plans':'An early spoken Chapter7/validation-sentinel proposal was discarded on rereading this whole-code task. No Chapter7/entry/sentinel control was executed. The failed plan construction started no command.','prep_anchor_checks':'Before execution, independent helper corrected error field to text and used wrong-typed budget rather than negative budget (negative is contractually clamped). No failed behavior run was hidden.'}
(E/'review-attempt-notes.json').write_text(json.dumps(notes,indent=2)+'\n')
prod=[]
for x in M['source_members']:
 p=Path(x['path'])
 if p.suffix in ['.go','.js','.css','.html'] and not p.name.endswith('_test.go') and not str(p).startswith(('evidence/','gui/test/')):prod.append(x)
(E/'source-review-coverage.json').write_text(json.dumps({'scope':'Whole production candidate, not only the 19-path diff; tracked source/archive identities independently verified. Reads split by subsystem, with recovery of capped output. Tests/fixtures and docs additionally inspected as described in report.','production_files':prod,'changed_paths':B['changed_paths'],'inherited_contracts':'Relevant Chapter1-6 publication/behavior sections read this turn; full Chapter7 and historical workflow/carryover/craft reads carried only at unchanged hashes. Current execution/student/carryover/readiness/Chapter8/report read in full; see prerequisite-carry-forward.json.','reviewed_boundaries':['public API and common vocabulary','one-request Engine renderer/stream provider boundary','Agent admission/dispatch/history/tool budget ownership','retained jobs/root/process lifetime','GUI startup/terminal/socket hub/full settings atomicity','production workbench/theme/font/TTS/pane and queue integration','external consumer/workflow/stream modules and CLI protocol','tests, fixtures, evidence packaging and preserved failure qualifications']},indent=2)+'\n')
