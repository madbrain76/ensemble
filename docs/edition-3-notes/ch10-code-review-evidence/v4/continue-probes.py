#!/usr/bin/env python3
import hashlib,json,os,shutil,signal,subprocess,tempfile,threading,time,traceback
from pathlib import Path
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';BASE=N/'ch10-code-review-evidence';E=BASE/'v4';PREV=BASE/'v3';M=json.loads((N/'ch10-code-review-input.json').read_text());C=json.loads((N/'ch10-student-source-manifest.json').read_text());W=None;COPY=None;env=None;before=None;rows=[];failure=None;phase='mechanically-corrected-probe-setup';stop=threading.Event();thread=None;seen={};processes=set();original={}
def digest(b):return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}

def ident(p):return digest(Path(p).read_bytes())

def write(n,x):p=E/n;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(json.dumps(x,indent=2,ensure_ascii=False)+'\n')

def git(root,*args):return subprocess.check_output(['git','--no-pager','-C',str(root),*args])

def verify():
 result={'inputs':[],'protected':[],'student':[],'old_seals':[]}
 for key,entries in [('inputs',M['inputs']),('protected',list(M['protected_tracked'].values())),('student',M['candidate']['files'])]:
  for x in entries:
   got=ident(x['path']);assert got=={k:x[k] for k in ['bytes','sha256']},x['path'];result[key].append({'path':x['path'],**got})
 assert M['candidate']['files']==[{'path':str(S/x['path']),**{k:x[k] for k in ['bytes','sha256']}} for x in C['tracked_files']]
 result['binaries']=[]
 for x in C['binaries']:
  assert ident(x['path'])=={k:x[k] for k in ['bytes','sha256']};result['binaries'].append({'path':x['path'],**ident(x['path'])})
 for k in ['local_path','evidence_path']:assert ident(C['archive'][k])=={a:C['archive'][a] for a in ['bytes','sha256']}
 for name in ['ch10-entry-audit-evidence','ch10-entry-followup-evidence','ch10-grader-review-evidence','ch10-grader-review-followup-evidence','ch10-default-load-diagnosis-evidence','ch10-prompt-fixture-review-evidence','ch10-student-final-reconciliation-evidence']:
  p=N/name/'inventory.json';j=json.loads(p.read_text());members=j['entries'] if name=='ch10-student-final-reconciliation-evidence' else j['files']
  for x in members:
   q=Path(x['path']);q=q if q.is_absolute() else R/q;assert ident(q)=={k:x[k] for k in ['bytes','sha256']},str(q)
  result['old_seals'].append({'path':str(p),**ident(p),'members':len(members)})
 result['repository']={'main_head':git(R,'rev-parse','HEAD').decode().strip(),'main_diff':digest(git(R,'diff','--binary')),'main_index':digest(git(R,'diff','--cached','--binary')),'student_head':git(S,'rev-parse','HEAD').decode().strip(),'student_tree':git(S,'rev-parse','HEAD^{tree}').decode().strip(),'student_status':git(S,'status','--porcelain=v1').decode()}
 assert result['repository']['student_head']==C['head'] and result['repository']['student_tree']==C['tree'] and result['repository']['student_status']==''
 return result

def sample():
 while not stop.wait(.15):
  if not W:continue
  try:
   for line in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True).splitlines():
    if str(W) in line:processes.add(line.strip())
   for root,ds,fs in os.walk(W/'tmp'):
    for f in fs:
     if f not in ['submission','root.test','main.test','gui.test','chapter1-cli','ch01']:continue
     p=Path(root)/f
     if p.is_file() and not p.is_symlink() and os.access(p,os.X_OK):
      st=p.stat();key=(str(p),st.st_size,st.st_mtime_ns)
      if key not in seen:seen[key]={'phase':phase,'path':str(p),**ident(p)}
  except (OSError,subprocess.SubprocessError):pass

def run(label,argv,cwd,deadline,allow=(0,1)):
 global phase
 assert shutil.disk_usage(R).free>=1024**3,'disk below 1GiB floor';phase=label
 out=E/(label+'.stdout.txt');err=E/(label+'.stderr.txt');start=time.monotonic();timed=False
 exe=shutil.which(argv[0],path=env['PATH']) if not Path(argv[0]).is_absolute() else argv[0]
 print('START',label,flush=True)
 with out.open('wb') as o,err.open('wb') as e:
  p=subprocess.Popen(argv,cwd=str(cwd),env=env,stdin=subprocess.DEVNULL,stdout=o,stderr=e,start_new_session=True)
  try:rc=p.wait(timeout=deadline)
  except subprocess.TimeoutExpired:timed=True;os.killpg(p.pid,signal.SIGTERM)
  if timed:
   try:rc=p.wait(timeout=5)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);rc=p.wait()
 row={'label':label,'argv':argv,'cwd':str(cwd),'environment':env,'executable':{'path':exe,**ident(exe)},'pid':p.pid,'process_group':p.pid,'deadline_seconds':deadline,'duration_seconds':round(time.monotonic()-start,6),'exit_code':rc,'timed_out':timed,'stdout':{'path':out.relative_to(E).as_posix(),**ident(out)},'stderr':{'path':err.relative_to(E).as_posix(),**ident(err)},'disk_free_bytes_after':shutil.disk_usage(R).free};rows.append(row);write('command-receipts.json',rows);print(label,'exit',rc,'seconds',row['duration_seconds'],flush=True)
 assert not timed and rc in allow,(label,rc,timed);assert shutil.disk_usage(R).free>=1024**3
 return row,out.read_text(),err.read_text()

def originals_ok():return all(ident(COPY/rel)==digest(b) for rel,b in original.items())
try:
 before=verify();write('verification-before.json',before);assert before==json.loads((PREV/'verification-before.json').read_text())==json.loads((PREV/'verification-after.json').read_text())
 gateRows=json.loads((PREV/'command-receipts.json').read_text());assert len(gateRows)==11 and all(x['exit_code']==0 and not x['timed_out'] for x in gateRows[:10]);assert gateRows[-1]['label']=='12-independent-public' and gateRows[-1]['exit_code']==1
 for x in gateRows:
  for k in ['stdout','stderr']:assert ident(PREV/x[k]['path'])=={a:x[k][a] for a in ['bytes','sha256']}
 oldout=(PREV/gateRows[-1]['stdout']['path']).read_text();assert '[build failed]' in oldout and '=== RUN' not in oldout
 ch1rows=json.loads((BASE/'v2/command-receipts.json').read_text());assert len(ch1rows)==1;c1=ch1rows[0];rep=json.loads((BASE/'v2/01-parsed-report.json').read_text());assert rep['passed'] is True and rep['score']==rep['max_score']==100 and c1['exit_code']==0 and c1['duration_seconds']==124.19327;assert [x['id'] for x in rep['checks']]==['protocol','wire','calls','replies','memory','growth','usage','generation_headers','generation_body'] and all(x['passed'] for x in rep['checks']);assert '1m2s headers delay' in ' '.join(rep['checks'][-2]['details']) and '1m2s body delay' in ' '.join(rep['checks'][-1]['details']);assert ident(BASE/'v2'/c1['stdout']['path'])=={k:c1['stdout'][k] for k in ['bytes','sha256']};assert ident(c1['argv'][0])['sha256']=='394a51635fe8538653fe9e49025c9c3817612c05b67106effd2ecbd445ab9fb2';assert any(x['path']==c1['argv'][-1] and ident(x['path'])=={k:x[k] for k in ['bytes','sha256']} for x in C['binaries'])
 write('carried-controls.json',{'root_authorization':'main-33','chapter1_receipt_identity':ident(BASE/'v2/command-receipts.json'),'chapter1':{'executions':1,'passed':True,'score':100,'exit_code':0,'duration_seconds':c1['duration_seconds'],'actual_ids':[x['id'] for x in rep['checks']]},'ten_completed_gate_receipt_identity':ident(PREV/'command-receipts.json'),'ten_completed_gates':gateRows[:10],'invalid_probe_compile':gateRows[-1],'behavioral_assertions_in_invalid_attempt':0,'all_old_output_hashes_reverified':True,'no_completed_gate_rerun':True})
 w=json.loads((PREV/'workspace.json').read_text());W=Path(w['owned_workspace']);COPY=Path(w['source_copy']);env=json.loads((PREV/'environment.json').read_text());write('workspace.json',w);write('environment.json',env);original={x['path']:(S/x['path']).read_bytes() for x in C['tracked_files']};assert originals_ok();assert shutil.disk_usage(R).free>=1024**3
 for name,target in [('independent_ch10_review_test.go.txt','reviewer_ch10_test.go'),('reviewer_legacy_ws_test.go.txt','gui/reviewer_legacy_ws_test.go')]:assert ident(COPY/target)==ident(PREV/'probes'/name)
 write('corrected-probe-plan.json',{'root_authorization':'main-33','new_commands_only':['12-independent-public','13-independent-gui','14-atomic-mutant','15-atomic-restored'],'stop_on_genuine_failure':True,'assertion_semantics_preserved':True,'mechanical_diffs':['public-probe-correction.diff.txt','gui-probe-correction.diff.txt'],'new_probe_pins':{p.name:ident(p) for p in (E/'probes').glob('*.txt')},'owned_tmp_entries_retained_from_completed_gates':len(list((W/'tmp').iterdir())),'disk_free_bytes_before':shutil.disk_usage(R).free})
 thread=threading.Thread(target=sample,daemon=True);thread.start()
 for source,target in [('independent_ch10_review_test.go.txt','reviewer_ch10_test.go'),('reviewer_legacy_ws_test.go.txt','gui/reviewer_legacy_ws_test.go')]:shutil.copyfile(E/'probes'/source,COPY/target)
 row,stdout,stderr=run('12-independent-public',['go','test','-count=1','-run','^TestReview','-v','-timeout=90s','.'],COPY,150)
 assert '=== RUN   TestReview' in stdout and '[build failed]' not in stdout and '[build failed]' not in stderr,'reviewer probe compile/setup failure';assert '--- PASS: TestReviewAtomicSavePreservesFile' in stdout,'atomic oracle baseline did not pass';write('12-independent-summary.json',{'exit_code':row['exit_code'],'assertion_failure_present':'--- FAIL:' in stdout,'actual_tests_ran':True,'atomic_oracle_baseline_passed':True,'not_retried':True})
 if row['exit_code']!=0:raise RuntimeError('Independent public behavioral assertions failed: stop under root main-33; no probe retry')
 row,stdout,stderr=run('13-independent-gui',['go','test','-count=1','-run','^TestReviewLegacyCollisionOverWebSocket$','-v','-timeout=45s','.'],COPY/'gui',90)
 assert '=== RUN   TestReviewLegacyCollisionOverWebSocket' in stdout and '[build failed]' not in stdout and '[build failed]' not in stderr,'GUI probe compile/setup failure'
 if row['exit_code']!=0:raise RuntimeError('Independent GUI assertions failed: stop under root main-33; no probe retry')
 assert originals_ok()
 target=COPY/'internal/agent/save.go';beforeb=target.read_bytes();anchor=b'func atomicWrite(path string, write func(*os.File) error) error {\n';assert beforeb.count(anchor)==1;mut=beforeb.replace(anchor,anchor+b'\t_ = os.WriteFile(path, nil, 0600) // reviewer sensitivity: wrong early truncation\n',1);target.write_bytes(mut);write('mutation-pins.json',{'path':'internal/agent/save.go','before':digest(beforeb),'mutated':digest(mut),'other_original_files_unchanged':all(ident(COPY/r)==digest(b) for r,b in original.items() if r!='internal/agent/save.go'),'expected_assertion':'ATOMIC_SENTINEL_CORRUPTED'})
 import difflib
 (E/'mutation.diff.txt').write_text(''.join(difflib.unified_diff(beforeb.decode().splitlines(True),mut.decode().splitlines(True),fromfile='original/internal/agent/save.go',tofile='owned-mutant/internal/agent/save.go')))
 try:
  row,stdout,stderr=run('14-atomic-mutant',['go','test','-count=1','-run','^TestReviewAtomicSavePreservesFile$','-v','-timeout=45s','.'],COPY,90)
  assert row['exit_code']==1 and '=== RUN   TestReviewAtomicSavePreservesFile' in stdout and 'ATOMIC_SENTINEL_CORRUPTED' in stdout and '[build failed]' not in stdout+stderr,'not a compile-valid behavioral negative';write('14-sensitivity.json',{'compiled_and_executed':True,'exit_code':1,'exact_assertion_reached':True,'mutation_is_owned_copy_only':True})
 finally:
  for rel,b in original.items():(COPY/rel).write_bytes(b)
  assert originals_ok();write('restored-copy-pins.json',{'all_original_files_byte_exact':True,'count':len(original),'files':{rel:ident(COPY/rel) for rel in original}})
 row,stdout,stderr=run('15-atomic-restored',['go','test','-count=1','-run','^TestReviewAtomicSavePreservesFile$','-v','-timeout=45s','.'],COPY,90)
 assert row['exit_code']==0 and '--- PASS: TestReviewAtomicSavePreservesFile' in stdout;assert originals_ok();write('outcome.json',{'completed_control_commands':len(rows),'chapter1':{'passed':rep['passed'],'score':rep['score'],'exit_code':0},'all_non_mutant_commands_green':all(x['exit_code']==0 for x in rows if x['label']!='14-atomic-mutant'),'atomic_sensitivity':'baseline green -> behavioral mutant red -> restored green','findings_must_be_reported_from_actual_logs':True,'no_runtime_retries':True})
except Exception as ex:
 failure={'exception':str(ex),'traceback':traceback.format_exc()};write('failure.json',failure);print('STOP',failure['traceback'],flush=True)
finally:
 if COPY and original:
  for rel,b in original.items():
   if (COPY/rel).read_bytes()!=b:(COPY/rel).write_bytes(b)
  write('copy-final-restoration.json',{'all_original_files_byte_exact':originals_ok(),'files':{rel:ident(COPY/rel) for rel in original}})
 stop.set()
 if thread:thread.join(timeout=10)
 write('observed-transient-executables.json',list(seen.values()));(E/'observed-processes.txt').write_text('\n'.join(sorted(processes))+'\n');tail={'failure':failure,'all_top_level_jobs_waited':True,'sampler_joined':not thread or not thread.is_alive(),'completed_commands':len(rows),'owned_workspace':str(W) if W else None}
 try:after=verify();write('verification-after.json',after);tail['before_after_equal']=after==before
 except Exception:tail['verification_error']=traceback.format_exc()
 if W:
  ps=subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True);tail['owned_processes_remaining']=[x.strip() for x in ps.splitlines() if str(W) in x];tail['owned_tmp_entries']=[p.name for p in (W/'tmp').iterdir()];tail['disk_free_bytes_after']=shutil.disk_usage(R).free
 write('completion.json',tail);print('COMPLETION',json.dumps({k:v for k,v in tail.items() if k!='owned_tmp_entries'}),flush=True)
