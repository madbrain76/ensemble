#!/usr/bin/env python3
import ast,hashlib,json,os,shutil,stat,subprocess,time
from pathlib import Path
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';E=N/'ch10-code-review-completion-evidence';OLD=N/'ch10-code-review-evidence';REPORT=N/'ch10-code-review-completion.md';M=json.loads((N/'ch10-code-review-completion-input.json').read_text());C=json.loads((N/'ch10-student-source-manifest.json').read_text())
tree=ast.parse((OLD/'v4/continue-probes.py').read_text());defs=[n for n in tree.body if isinstance(n,ast.FunctionDef) and n.name in ['digest','ident','git','verify']];exec(compile(ast.Module(body=defs,type_ignores=[]),'read-only-preservation-functions','exec'))
final=verify();assert final==json.loads((E/'verification-before.json').read_text())==json.loads((E/'verification-after.json').read_text());(E/'final-preservation.json').write_text(json.dumps(final,indent=2)+'\n')
assert ident('/tmp/cadence-ch10-code-review-completion-task.txt')=={'bytes':6971,'sha256':'fb8cfe55bd957f304d998179d35ed38645eb901dbb6fd3ccc434e825f4661e0e'}
assert ident(OLD/'inventory.json')=={'bytes':24503,'sha256':'d781ef13cd6cdfafe8dcc32dcc50ffde64dda611d0faa7472bda0a53acc6e932'};old=json.loads((OLD/'inventory.json').read_text());assert len(old['files'])==116 and sum(x['bytes'] for x in old['files'])==1138356
for x in old['files']:assert ident(R/x['path'])=={k:x[k] for k in ['bytes','sha256']}
assert {x['path'] for x in old['files']}=={str(p.relative_to(R)) for p in [N/'ch10-code-review.md']+list(OLD.rglob('*')) if p.is_file() and p!=OLD/'inventory.json'}
rows=json.loads((E/'command-receipts.json').read_text());assert len(rows)==1;row=rows[0];assert row['label']=='01-public-controls' and row['exit_code']==1 and not row['timed_out']
for k in ['stdout','stderr']:assert ident(E/row[k]['path'])=={z:row[k][z] for z in ['bytes','sha256']}
for k,v in {'GOMAXPROCS':'2','GOFLAGS':'-mod=readonly -p=1','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off'}.items():assert row['environment'][k]==v
assert 'LLM_BASE_URL' not in row['environment'];assert not any('KEY' in k or 'TOKEN' in k for k in row['environment']);assert row['disk_free_before']>=1024**3 and row['disk_free_after']>=1024**3
log=(E/row['stdout']['path']).read_text();assert '[build failed]' not in log and '[no tests to run]' not in log and '=== RUN   TestReviewReplayModesAndRichRender' in log
for s in ['invalid context turn','model "gpt-5-course" does not support image media "image/png"','model "gemini-2.5-pro" does not support image media "image/png"']:assert s in log
for s in ['--- PASS: TestReviewMalformedLoadAtomicityAndFrozenPrompt','--- PASS: TestReviewAtomicSavePreservesFile','sentinel unchanged','valid selected and persisted primary bytes=0','valid selected and persisted primary bytes=3','valid selected and persisted primary bytes=18']:assert s in log
assert not any(p.name.startswith(('02-gui-controls','03-atomic-mutant','04-atomic-restored')) for p in E.iterdir())
completion=json.loads((E/'completion.json').read_text());assert completion['controls_completed']==1 and completion['all_top_level_jobs_waited'] and completion['sampler_joined'] and completion['before_after_equal'];assert not completion['owned_processes_remaining'] and not completion['mutation_applied'] and not completion['mutation_restored'];assert completion['source_restored']['all_original_bytes_and_archive_modes_match']
cleanup=json.loads((E/'cleanup-after.json').read_text());assert cleanup['hashes_written_before_cleanup'] and cleanup['cleanup_complete'] and not cleanup['owned_processes_or_group_members_remaining'];w=json.loads((E/'workspace.json').read_text());W=Path(w['workspace']);CP=Path(w['source_copy']);assert not list((W/'tmp').iterdir());pins=json.loads((E/'archive-extraction-pins.json').read_text());assert len(pins['original_files'])==137
assert sorted(str(p.relative_to(CP)) for p in CP.rglob('*') if p.is_file())==sorted(x['path'] for x in pins['original_files'])
mode_rows={}
for b in subprocess.check_output(['git','--no-pager','ls-tree','-rz','--full-tree',C['head']],cwd=S).split(b'\0'):
 if b:
  info,path=b.split(b'\t',1);mode_rows[path.decode()]=int(info.split()[0],8)
for x in pins['original_files']:
 p=CP/x['path'];assert ident(p)=={k:x[k] for k in ['bytes','sha256']};assert stat.S_IMODE(p.stat().st_mode)==x['archive_mode'];assert bool(x['archive_mode']&0o111)==bool(mode_rows[x['path']]&0o111)
live=[]
for line in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,pgid=,command='],text=True).splitlines():
 q=line.strip().split(None,3)
 if len(q)==4 and (str(W) in q[3] or int(q[2])==row['process_group']):live.append(line.strip())
assert not live,live
classified=json.loads((E/'failure-and-assertion-classification.json').read_text());assert classified['candidate_defects_established']==0 and len(classified['fixture_failures'])==3
receipt={'time':time.time(),'verdict':'QUALIFIED / PARTIAL — NOT ACCEPTED','candidate_defects_established':0,'root_decision_required':True,'executed_controls':1,'executed_command_exit':1,'duration_seconds':row['duration_seconds'],'compiled':True,'fixture_failures':classified['fixture_failures'],'exact_pass_lines':classified['exact_pass_lines'],'unrun':['GUI collision/reconnect/restart/new-generation','compile-valid atomicity mutation','focused restored positive'],'no_probe_repair_or_retry_after_execution':True,'no_gate_or_standard_suite_replays':True,'all_jobs_finished':True,'sampler_joined':True,'owned_processes_or_group_members_remaining':live,'archive_extracted_exact_137':True,'archive_modes_preserved':True,'git_executable_class_checked_at_cleanup_and_finalization':True,'owned_tmp_empty':True,'two_reviewer_files_hashed_before_cleanup':True,'no_shared_cleanup':True,'all_35_inputs_and_33_protected_files_and_candidate_bytes_archives_binaries_repository_state_preserved':True,'all_seven_prior_inventories_and_original_116_member_seal_preserved':True,'original_inventory_identity':ident(OLD/'inventory.json'),'original_report_identity':ident(N/'ch10-code-review.md'),'source_copy_method':'pinned TAR extraction, not checkout copy','no_student_source_or_public_guidance_or_grader_changes':True,'no_delegated_children_browser_paid_provider_real_keys':True,'no_automatic_continuation':True,'disk_free':shutil.disk_usage(R).free}
assert receipt['disk_free']>=1024**3;(E/'final-receipt.json').write_text(json.dumps(receipt,indent=2)+'\n');(E/'command-output-hash-ledger.json').write_text(json.dumps(rows,indent=2)+'\n')
I=E/'inventory.json';files=[REPORT]+sorted(p for p in E.rglob('*') if p.is_file() and p!=I);assert not any(p.is_symlink() for p in E.rglob('*'));members=[]
for p in files:
 assert p.name!='settings.json' and p.suffix.lower() not in ['.tar','.zip','.gz','.exe','.so','.dll','.dylib'];b=p.read_bytes();b.decode('utf-8');assert b'\x00' not in b;members.append({'path':str(p.relative_to(R)),'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()})
inv={'format':'exact-portable-inventory-v1','base':str(R),'self_excluded':str(I.relative_to(R)),'files':members,'file_count_excluding_inventory':len(members),'total_payload_bytes':sum(x['bytes'] for x in members),'utf8_no_nul_verified':True,'no_missing_or_extra_payload_files':True,'no_archive_binary_settings_or_private_raw_capture_payload':True};I.write_text(json.dumps(inv,indent=2)+'\n')
assert {x['path'] for x in members}=={str(p.relative_to(R)) for p in [REPORT]+[p for p in E.rglob('*') if p.is_file() and p!=I]}
for x in members:assert ident(R/x['path'])=={k:x[k] for k in ['bytes','sha256']}
print(json.dumps({'report':{'path':str(REPORT),**ident(REPORT)},'inventory':{'path':str(I),**ident(I)},'members_excluding_inventory':len(members),'payload_bytes':inv['total_payload_bytes'],'verdict':receipt['verdict'],'candidate_defects':0,'all_jobs_finished':True,'no_owned_processes':True,'old_seals_unchanged':True,'archive_copy_restored_exact137':True,'disk_free':receipt['disk_free']},indent=2))
