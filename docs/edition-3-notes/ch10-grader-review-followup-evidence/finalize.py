#!/usr/bin/env python3
"""Read-only final verification and portable evidence sealing; no runtime controls."""
from pathlib import Path
import hashlib,json,shutil,subprocess
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');E=R/'docs/edition-3-notes/ch10-grader-review-followup-evidence';P=R/'docs/edition-3-notes/ch10-grader-review-followup.md';M=json.loads((R/'docs/edition-3-notes/ch10-grader-review-input.json').read_text());W=Path(json.loads((E/'workspace.json').read_text())['owned_scratch'])
def ident(p):
 b=Path(p).read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def write(n,x):
 p=E/n;assert not p.exists(),p;p.write_text(json.dumps(x,indent=2)+'\n')
def git(root,*args):return subprocess.check_output(['git','--no-pager',*args],cwd=root)
def pins(root):return {p.relative_to(root).as_posix():ident(p) for p in sorted(root.rglob('*')) if p.is_file()}
before=json.loads((E/'verification-before.json').read_text());after=json.loads((E/'verification-after.json').read_text());assert before['repositories']==after['repositories'];assert before['prior_seals']==after['prior_seals'];outcome=json.loads((E/'run-outcome.json').read_text());assert outcome['failure'] is None and outcome['verdict'].startswith('PASS:');rc=json.loads((E/'command-receipts.json').read_text());assert len(rc)==22 and all(not c['timed_out'] and c['process_exit'] in c['allowed_process_exits'] for c in rc)
for c in rc:
 for key in ['stdout','stderr']:assert ident(E/c[key]['path'])=={k:c[key][k] for k in ['bytes','sha256']}
 assert ident(c['argv'][0])==c['executable_identity'],c['name']
for label in ['01-ch9-real-reference','07-mutant-ch10-wrong-route','09-restored-ch10']:
 raw=json.loads((E/(label+'.stdout.txt')).read_text());parsed=json.loads((E/(label+'.report.json')).read_text());assert raw==parsed
for label in ['01-ch9-real-reference','09-restored-ch10']:
 d=json.loads((E/(label+'.detector.json')).read_text());assert d['actual_process_exit']==0 and d['report_passed'] is True and d['report_score']==d['report_max_score']==100 and d['detector']['calculated_detector_status']==0
m=json.loads((E/'07-mutant-ch10-wrong-route.detector.json').read_text());assert m['actual_process_exit']==1 and m['detector']['calculated_detector_status']==1 and m['detector']['actual_title']=='Edition 2, Chapter 9 — Skills'
guards=json.loads((E/'build-guard-results.json').read_text());assert len(guards)==14 and all(c['build_marker_absent'] and c['adapter_process_exit']==c['expected_adapter_exit'] for c in guards)
positive=json.loads((E/'04-build-entry-positive.json').read_text());adapter_pid=next(c['pid'] for c in rc if c['name']=='04-valid-directory-reaches-build');assert positive['marker']['ppid']==adapter_pid and positive['actual_build_marker'] and positive['adapter_process_exit']==2
write('control-ledger.json',{'verdict':'PASS for bounded review scope; root inspection/join/acceptance pending','actual_commands':len(rc),'unexpected_control_failures':[],'retries':0,'expected_nonzero_process_exits':[{'name':c['name'],'actual_process_exit':c['process_exit']} for c in rc if c['process_exit']!=0],'mutant_detector_calculated_status':m['detector']['calculated_detector_status'],'mutant_process_exit_separate':m['actual_process_exit'],'new_ch10_unmutated_baselines_before_sensitivity':0,'new_full_compositions':3,'original_ch9_85_failed_unchanged':True,'no_causal_diagnosis_of_old_parity_failure':True,'original_build_marker_omission_now_supplied':True,'original_wrong_historical_mutation_plan_never_executed':True,'actual_mutation':'only copied cmd/grade-edition3/main.go dispatch Chapter10 -> Chapter9','copied_files_all_restored':134,'other_old_runtime_routes_carried_not_grading_rerun':True,'four_older_empty_flag_validation_guards_only':True,'prior_80_sentinel_cases_carried_not_repeated':True,'no_package_gates_or_campaigns_repeated':True,'setup_note':'Read-only dependency-inspection helper emitted a Python positional maxsplit deprecation warning; runner uses keyword form. No control failed or retried.'})
final={'pin_groups':{},'prior_seals':{},'repositories':{}}
for label,root,group in [('inputs',R,M['inputs']),('protected_tracked',R,M['protected_tracked']),('predecessor',S,M['predecessor']['files'])]:
 for rel,want in group.items():assert ident(root/rel)==want,rel
 final['pin_groups'][label]={'count':len(group),'all_match':True}
for rel,want in M['candidate']['source'].items():assert ident(R/rel)==want,rel
assert ident(M['candidate']['binary']['path'])=={k:M['candidate']['binary'][k] for k in ['bytes','sha256']}
assert ident(R/'internal/gradeedition2/ch10.go')==M['candidate']['historical_grader_unchanged'];assert git(R,'diff','--','cmd/grade-edition3/main.go','cmd/grade-edition3/main_test.go')==(R/'docs/edition-3-notes/ch10-grader-evidence/candidate.diff').read_bytes()
refpins=json.loads((E/'reference-source-pins.json').read_text())
for ch,expected in refpins.items():assert pins(R/'solutions'/ch)==expected;final['pin_groups'][ch]={'count':len(expected),'all_match':True}
cp=json.loads((E/'copy-original-pins.json').read_text())['files'];assert len(cp)==134;assert pins(W/'adapter-copy')==cp
for rel,want in cp.items():assert ident(R/rel)==want,rel
for n in ['copy-before.json','copy-restored-before-build.json','copy-final-restoration.json']:assert json.loads((E/n).read_text())['files']==cp
mut=json.loads((E/'copy-mutated.json').read_text())['files'];assert set(mut)==set(cp);assert [x for x in cp if cp[x]!=mut[x]]==['cmd/grade-edition3/main.go'];final['copied_file_restoration']={'count':134,'exact_file_set_and_bytes':True,'all_original_sources_still_match':True,'root':str(W/'adapter-copy')}
for name,want in before['prior_seals'].items():
 d=R/'docs/edition-3-notes'/name;I=d/'inventory.json';j=json.loads(I.read_text());assert ident(I)==want['inventory']
 for row in j['files']:assert ident(R/row['path'])=={k:row[k] for k in ['bytes','sha256']},row['path']
 outside={row['path'] for row in j['files'] if not row['path'].startswith(d.relative_to(R).as_posix()+'/')};actual={p.relative_to(R).as_posix() for p in d.rglob('*') if p.is_file() and p!=I}|outside;assert actual=={row['path'] for row in j['files']};final['prior_seals'][name]={'inventory':ident(I),'file_count':len(j['files']),'exact_closed_set':True}
for name,root in [('main',R),('student',S)]:
 d=git(root,'diff','--binary');c=git(root,'diff','--cached','--binary');s=git(root,'status','--porcelain=v1','--untracked-files=no');v={'head':git(root,'rev-parse','HEAD').decode().strip(),'tree':git(root,'rev-parse','HEAD^{tree}').decode().strip(),'unstaged_diff':{'bytes':len(d),'sha256':hashlib.sha256(d).hexdigest()},'staged_diff':{'bytes':len(c),'sha256':hashlib.sha256(c).hexdigest()},'tracked_status':{'bytes':len(s),'sha256':hashlib.sha256(s).hexdigest()}};assert v==before['repositories'][name];final['repositories'][name]=v
assert git(S,'status','--porcelain=v1','--untracked-files=normal')==b''
live=[x.strip() for x in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True).splitlines() if str(W) in x];assert not live,live;assert not list((W/'tmp').iterdir());assert not (W/'build-entry-marker').exists();disk=shutil.disk_usage(R).free;assert disk>=1024**3
final.update({'all_jobs_waited':True,'owned_processes_remaining':live,'owned_runtime_tmp_empty':True,'build_marker_removed':True,'owned_scratch_retained_hash_only':str(W),'all_prior_seals_and_candidate_preserved':True,'predecessor_clean':True,'no_shared_cleanup':True,'disk_free_bytes':disk,'no_new_control_during_finalization':True});write('final-preservation-and-completion.json',final)
I=E/'inventory.json';assert not I.exists();assert not any(p.is_symlink() for p in E.rglob('*'));members=[P]+sorted(p for p in E.rglob('*') if p.is_file());rows=[]
for p in members:
 assert p.name!='settings.json';b=p.read_bytes();assert b'\0' not in b,p;b.decode('utf-8');rows.append({'path':p.relative_to(R).as_posix(),**ident(p)})
rows.sort(key=lambda r:r['path']);write('inventory.json',{'schema':'ch10-grader-review-followup-portable-inventory-v1','path_base':'repository root','self_exclusion':'Only this inventory.json is excluded; external inventory hash supplied in final submission','scope':'Exact standalone report plus every regular evidence file except inventory itself','encoding':'UTF-8 without NUL','file_count':len(rows),'total_bytes':sum(r['bytes'] for r in rows),'files':rows})
j=json.loads(I.read_text());actual={P.relative_to(R).as_posix()}|{p.relative_to(R).as_posix() for p in E.rglob('*') if p.is_file() and p!=I};assert actual=={r['path'] for r in j['files']}
for row in j['files']:assert ident(R/row['path'])=={k:row[k] for k in ['bytes','sha256']}
b=I.read_bytes();assert b'\0' not in b;b.decode('utf-8')
print(json.dumps({'verdict':'PASS bounded scope; root inspection/join pending','report':ident(P),'inventory':ident(I),'file_count_excluding_inventory':len(rows),'bytes_excluding_inventory':j['total_bytes'],'exact_closed_set_reverified':True,'all_utf8_nul_free':True,'final':final},indent=2))
