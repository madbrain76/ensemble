#!/usr/bin/env python3
import ast,hashlib,json,os,shutil,stat,subprocess,time
from pathlib import Path
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';E=N/'ch10-code-review-remaining-evidence';REPORT=N/'ch10-code-review-remaining.md';M=json.loads((N/'ch10-code-review-remaining-input.json').read_text());C=json.loads((N/'ch10-student-source-manifest.json').read_text())
# Only this task's own three read-only functions; never execute a parent verifier.
tree=ast.parse((E/'run-remaining.py').read_text());defs=[n for n in tree.body if isinstance(n,ast.FunctionDef) and n.name in ['ident','git','verify']];exec(compile(ast.Module(body=defs,type_ignores=[]),'remaining-read-only-functions','exec'))
assert ident('/tmp/cadence-ch10-code-review-remaining-controls-task.txt')=={'bytes':6471,'sha256':'87b4595d2d5a7d9793d05c296375879751875e030a3fd9e7dac123cd202d025a'}
final=verify();assert final==json.loads((E/'verification-before.json').read_text())==json.loads((E/'verification-after.json').read_text());(E/'final-preservation.json').write_text(json.dumps(final,indent=2)+'\n')
rows=json.loads((E/'command-receipts.json').read_text());assert len(rows)==5 and [r['exit_code'] for r in rows]==[0,0,0,1,0]
selectors=['^TestReviewReplayModesAndRichRender$/anthropic$/zero-cached-plus-tail$','^TestReviewReplayModesAndRichRender$/(openai|gemini)$','^TestReviewLegacyCollisionOverWebSocket$','^TestReviewAtomicSavePreservesFile$','^TestReviewAtomicSavePreservesFile$'];logs=[];assertions=[]
for i,r in enumerate(rows):
 assert r['argv'][r['argv'].index('-run')+1]==selectors[i] and not r['timed_out'];assert r['disk_free_before']>=1024**3 and r['disk_free_after']>=1024**3
 for k in ['stdout','stderr']:assert ident(E/r[k]['path'])=={z:r[k][z] for z in ['bytes','sha256']}
 assert r['stderr']['bytes']==0
 for k,v in {'GOMAXPROCS':'2','GOFLAGS':'-mod=readonly -p=1','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off'}.items():assert r['environment'][k]==v
 assert 'LLM_BASE_URL' not in r['environment'] and not any('KEY' in k or 'TOKEN' in k for k in r['environment']);text=(E/r['stdout']['path']).read_text();logs.append(text);assert '[build failed]' not in text and '[no tests to run]' not in text and 'panic:' not in text
 assertions.append({'label':r['label'],'exit':r['exit_code'],'duration_seconds':r['duration_seconds'],'lines':[l for l in text.splitlines() if any(s in l for s in ['=== RUN','--- PASS:','--- FAIL:','ATOMIC_SENTINEL_CORRUPTED','sentinel unchanged','real WebSocket'])]})
assert '--- PASS: TestReviewReplayModesAndRichRender/anthropic/zero-cached-plus-tail' in logs[0]
for v in ['openai','gemini']:
 for mode in ['checkpoint-only','null-context-full-replay','zero-null-full-replay','zero-cached-plus-tail','ignored-prefix','strict-tail']:assert f'--- PASS: TestReviewReplayModesAndRichRender/{v}/{mode}' in logs[1]
assert '=== RUN   TestReviewReplayModesAndRichRender/anthropic' not in logs[1]
assert 'fake requests=1' in logs[2] and 'archived prefix unchanged' in logs[2] and '--- PASS: TestReviewLegacyCollisionOverWebSocket' in logs[2]
assert 'ATOMIC_SENTINEL_CORRUPTED: got="" want="REVIEW_ORIGINAL_REGULAR_FILE\\n"' in logs[3] and '--- FAIL: TestReviewAtomicSavePreservesFile' in logs[3]
assert '--- PASS: TestReviewAtomicSavePreservesFile' in logs[4] and 'sentinel unchanged' in logs[4]
for text in logs[:2]:assert '=== RUN   TestReviewMalformed' not in text and '=== RUN   TestReviewAtomic' not in text
original_probe=(N/'ch10-code-review-completion-evidence/probes/independent_ch10_review_test.go.txt').read_text();new_probe=(E/'probes/independent_ch10_review_test.go.txt').read_text();assert new_probe==original_probe.replace('"Idle"','"idle"').replace('gpt-5-course','gpt-4.1-2025-04-14').replace('gemini-2.5-pro','gemini-3.8-flash');assert original_probe[original_probe.index('func TestReviewAtomicSavePreservesFile('):]==new_probe[new_probe.index('func TestReviewAtomicSavePreservesFile('):];assert (E/'probes/reviewer_legacy_ws_test.go.txt').read_bytes()==(N/'ch10-code-review-completion-evidence/probes/reviewer_legacy_ws_test.go.txt').read_bytes()
carry=json.loads((E/'carry.json').read_text());base=carry['atomic_baseline_receipt'];O=N/'ch10-code-review-completion-evidence'
for k in ['stdout','stderr']:assert ident(O/base[k]['path'])=={z:base[k][z] for z in ['bytes','sha256']}
assert '--- PASS: TestReviewAtomicSavePreservesFile' in (O/base['stdout']['path']).read_text()
pins=json.loads((E/'source-copy-pins.json').read_text());CP=Path(pins['source_copy']);assert len(pins['files'])==137 and sorted(str(p.relative_to(CP)) for p in CP.rglob('*') if p.is_file())==sorted(x['path'] for x in pins['files'])
for x in pins['files']:
 assert ident(CP/x['path'])=={k:x[k] for k in ['bytes','sha256']} and stat.S_IMODE((CP/x['path']).stat().st_mode)==x['archive_mode'];assert bool(x['archive_mode']&0o111)==bool(x['git_mode']&0o111)
assert json.loads((E/'restoration.json').read_text())=={'files':137,'all_bytes_and_archive_modes_match':True};assert json.loads((E/'source-final-restore.json').read_text())=={'files':137,'all_bytes_and_archive_modes_match':True}
mp=json.loads((E/'mutation-pins.json').read_text());assert mp['before']==ident(CP/'internal/agent/save.go') and mp['mutant']['sha256']!=mp['before']['sha256']
w=json.loads((E/'workspace.json').read_text());W=Path(w['runtime']);assert not list((W/'tmp').iterdir());groups={r['process_group'] for r in rows};live=[]
for line in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,pgid=,command='],text=True).splitlines():
 q=line.strip().split(None,3)
 if len(q)==4 and (str(W) in q[3] or str(CP) in q[3] or int(q[2]) in groups):live.append(line.strip())
assert not live,live
completion=json.loads((E/'completion.json').read_text());assert completion['failure'] is None and completion['all_top_level_jobs_waited'] and completion['mutation_applied'] and completion['mutation_restored_before_focused_positive'] and completion['commands']==5
clean=json.loads((E/'cleanup-after.json').read_text());assert clean['durable_hashes_before_delete'] and clean['bytes_and_archive_modes_match'] and clean['source_copy_exact_original_files']==137 and clean['owned_tmp_empty'] and not clean['shared_cleanup']
assert json.loads((E/'pre-control-completion.json').read_text())['commands']==0
(E/'assertion-ledger.json').write_text(json.dumps(assertions,indent=2)+'\n')
receipt={'time':time.time(),'recommendation':'PASS for combined first whole-candidate code review; root decision and join required','candidate_defects_established':0,'five_expected_exits':[r['exit_code'] for r in rows],'public_anthropic_zero_passed':True,'openai_and_gemini_twelve_cases_passed':True,'gui_real_websocket_replay_reconnect_restart_continuation_passed':True,'gui_fake_provider_requests':1,'same_atomic_function_and_sentinel':True,'carried_valid_atomic_baseline':True,'compiled_mutant_actual_sentinel_failure':True,'exact_137_byte_mode_restore_before_focused_positive':True,'focused_restored_positive':True,'all_jobs_finished':True,'remaining_processes_or_groups':live,'owned_tmp_empty':True,'cleanup_hashes_durable_before_delete':True,'original_seals_reports_inputs_protected_candidate_binaries_archives_unchanged':True,'carried_grades_and_ten_suites_not_replayed':True,'no_parent_verifier_writer_executed':True,'pre_control_verifier_error_preserved':'wrong older inventory pathname; zero commands; only pathname corrected','no_behavioural_retries_or_extra_controls':True,'no_student_grader_guidance_edits_children_browser_paid_provider_real_keys_shared_cleanup':True,'operational_exercise_completion_not_inferred':True,'no_automatic_continuation':True,'disk_free':shutil.disk_usage(R).free};assert receipt['disk_free']>=1024**3;(E/'final-receipt.json').write_text(json.dumps(receipt,indent=2)+'\n')
I=E/'inventory.json';files=[REPORT]+sorted(p for p in E.rglob('*') if p.is_file() and p!=I);assert not any(p.is_symlink() for p in E.rglob('*'));members=[]
for p in files:
 assert p.name!='settings.json' and p.suffix.lower() not in ['.tar','.zip','.gz','.exe','.so','.dylib','.dll'];b=p.read_bytes();b.decode('utf-8');assert b'\x00' not in b;members.append({'path':str(p.relative_to(R)),'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()})
inv={'format':'exact-portable-inventory-v1','base':str(R),'self_excluded':str(I.relative_to(R)),'files':members,'file_count_excluding_inventory':len(members),'total_payload_bytes':sum(x['bytes'] for x in members),'utf8_no_nul':True,'no_missing_or_extra_members':True,'no_archive_binary_settings_private_capture_payload':True};I.write_text(json.dumps(inv,indent=2)+'\n')
assert {x['path'] for x in members}=={str(p.relative_to(R)) for p in [REPORT]+[p for p in E.rglob('*') if p.is_file() and p!=I]}
for x in members:assert ident(R/x['path'])=={k:x[k] for k in ['bytes','sha256']}
print(json.dumps({'report':{'path':str(REPORT),**ident(REPORT)},'inventory':{'path':str(I),**ident(I)},'members_excluding_inventory':len(members),'payload_bytes':inv['total_payload_bytes'],'verdict':'PASS recommendation; root decision required','all_five_controls_expected':True,'all_jobs_finished':True,'old_seals_preserved':True,'source_copy_exact137':True,'no_owned_processes':True,'disk_free':receipt['disk_free']},indent=2))
