#!/usr/bin/env python3
"""Only one corrected full Chapter10 and one exact historical persistence test."""
from pathlib import Path
import hashlib,json,os,shutil,signal,subprocess,tempfile,threading,time,traceback
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';E=N/'ch10-prompt-fixture-review-evidence';M=json.loads((N/'ch10-prompt-fixture-review-input.json').read_text());C=M['candidate'];F=M['fixture_candidate'];W=None;receipts=[];stop=threading.Event();observed={};phase='setup';sampler=None;before=None;failure=None
IDS=['save-shape','default-load','replay-equals-snapshot','tail-applied-once','log-not-needed','bad-save-refused','ch9-parity'];POINTS=[15,20,20,15,10,5,15];TITLE='Edition 2, Chapter 10 — Persistence'
def ident(p):
 b=Path(p).read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def digest(b):return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def write(name,obj):
 p=E/name;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(json.dumps(obj,indent=2,ensure_ascii=False)+'\n')
def git(p,*args):return subprocess.check_output(['git','--no-pager','-C',str(p),*args])
def check_rows(rows,base):
 out=[]
 for row in rows:
  p=Path(row['path']);p=p if p.is_absolute() else base/p;want={k:row[k] for k in ['bytes','sha256']};assert ident(p)==want,p;out.append({'path':str(p),**want})
 return out
def verify():
 inputs=check_rows(M['inputs'],R);protected={}
 for rel,want in M['protected_tracked'].items():assert ident(R/rel)==want,rel;protected[rel]=want
 files={}
 for rel,want in C['tracked_files'].items():assert ident(S/rel)==want,rel;files[rel]=want
 assert git(S,'rev-parse','HEAD').decode().strip()==C['head'];assert git(S,'rev-parse','HEAD^{tree}').decode().strip()==C['tree'];assert git(S,'rev-parse','HEAD^').decode().strip()==C['parent'];assert not git(S,'status','--porcelain=v1').strip()
 bins=check_rows(C['binaries'],S)
 for p in [C['archive']['path'],C['archive']['local_evidence_copy']]:assert ident(p)=={k:C['archive'][k] for k in ['bytes','sha256']},p
 for key in ['before','after','binary','commands']:
  row=F[key];assert ident(row['path'])=={k:row[k] for k in ['bytes','sha256']}
 sealpath=N/'ch10-student-evidence/seal.json';seal=json.loads(sealpath.read_text());members=check_rows(seal['entries'],R)
 inventories=[]
 for rel in ['ch10-entry-audit-evidence/inventory.json','ch10-entry-followup-evidence/inventory.json','ch10-grader-review-evidence/inventory.json','ch10-grader-review-followup-evidence/inventory.json','ch10-default-load-diagnosis-evidence/inventory.json']:
  p=N/rel;j=json.loads(p.read_text());rows=check_rows(j['files'],R);inventories.append({'path':str(p),**ident(p),'verified_members':len(rows),'verified_payload_bytes':sum(x['bytes'] for x in rows)})
 localref=json.loads((N/'ch10-entry-audit-evidence/reference-source-pins.json').read_text());reference={}
 for rel,h in localref.items():
  p=R/'solutions/ch10'/rel;want=h if isinstance(h,dict) else {'sha256':h};got=ident(p);assert all(got[k]==v for k,v in want.items()),p;reference[rel]=got
 rootchecks=json.loads(Path(F['commands']['path']).read_text());rootout=[]
 for row in rootchecks:
  assert row['exit_code']==0
  for key in ['stdout','stderr']:
   x=row[key];assert ident(x['path'])=={k:x[k] for k in ['bytes','sha256']};rootout.append({'command':row['label'],'stream':key,**x})
 old=json.loads((N/'ch10-student-evidence/grader-ch10-final-v1.json').read_text());declarations={flag:{'path':old['argv'][old['argv'].index(flag)+1],**ident(old['argv'][old['argv'].index(flag)+1])} for flag in ['-consumer-case','-media-models']}
 return {'inputs':inputs,'protected':protected,'candidate_files':files,'candidate_binaries':bins,'candidate_head':C['head'],'candidate_tree':C['tree'],'candidate_parent':C['parent'],'candidate_status':'','blocked_seal':{'path':str(sealpath),**ident(sealpath),'verified_members':len(members)},'prior_seals':inventories,'historical_reference_source_files':reference,'root_command_outputs':rootout,'runtime_declarations':declarations,'main_head':git(R,'rev-parse','HEAD').decode().strip(),'main_tree':git(R,'rev-parse','HEAD^{tree}').decode().strip(),'main_unstaged_diff_identity':digest(git(R,'diff','--binary')),'main_staged_diff_identity':digest(git(R,'diff','--cached','--binary'))}
def source_and_carry():
 old=Path(F['before']['path']).read_bytes();new=Path(F['after']['path']).read_bytes();needle=b'CURRENT_RUNTIME_PROMPT\\n"), 0600)';replacement=b'CURRENT_RUNTIME_PROMPT"), 0600)';assert old.count(needle)==1;assert new==old.replace(needle,replacement,1)
 exact=git(R,'diff','--','internal/gradeedition2/ch10.go');assert exact==(N/'ch10-prompt-fixture-evidence/candidate.diff').read_bytes();(E/'reviewed-exact.diff.txt').write_bytes(exact)
 active=new.decode();assert 'request.JSON()["model"] != "gpt-5-course" || persistenceSystem(request, "openai") != "CURRENT_RUNTIME_PROMPT"' in active
 write('exact-diff-proof.json',{'before':F['before'],'after':F['after'],'replacement_count':1,'old_source_literal_suffix':needle.decode(),'new_source_literal_suffix':replacement.decode(),'source_bytes_removed':2,'runtime_fixture_bytes_removed':1,'runtime_byte_removed':'0a','everything_else_in_checker_byte_identical':True,'exact_assertion_unchanged':True,'exact_diff_identity':digest(exact),'protected_adapter_public_student_pins_unchanged':True})
 d=N/'ch10-default-load-diagnosis-evidence';o=json.loads((d/'corrected-v2/observations.json').read_text());assert len(o)==3 and [x['wire']['historical_conjunction_predicate'] for x in o]==[False,False,True];assert [x['wire']['system_utf8_bytes'] for x in o]==[23,23,22];assert all(x['wire']['historical_model_predicate'] and x['rewritten_save']['config_matches_actual_request'] for x in o)
 write('carried-negative-control.json',{'new_runtime_replay':False,'diagnosis_inventory':{'path':str(d/'inventory.json'),**ident(d/'inventory.json')},'observations_identity':{'path':str(d/'corrected-v2/observations.json'),**ident(d/'corrected-v2/observations.json')},'observed_old_predicates':[{'label':x['label'],'model':x['wire']['model'],'system_text':x['wire']['system_text'],'utf8_bytes':x['wire']['system_utf8_bytes'],'model_predicate':x['wire']['historical_model_predicate'],'exact_system_predicate':x['wire']['historical_exact_system_predicate'],'conjunction':x['wire']['historical_conjunction_predicate'],'metadata_matches_wire':x['rewritten_save']['config_matches_actual_request']} for x in o],'old_source_fixture_matches_lf_negative':True,'new_source_fixture_matches_no_lf_contrast':True,'limitations':['Prior source seed was handcrafted snapshot-only, not original archived grader fixture.','Original diagnosis invalid comma-list attempt and omitted inert Go env fields remain retained; not repeated.','Old 80/100 exit1 report remains unchanged.']})
 oldresult=json.loads((N/'ch10-student-evidence/grader-ch10-final-v1.stdout.txt').read_text());assert oldresult['score']==80 and oldresult['passed'] is False;assert oldresult['checks'][1]['id']=='default-load' and oldresult['checks'][1]['passed'] is False
 return oldresult

def sample():
 while not stop.wait(.075):
  try:
   for p in (W/'tmp').glob('course-grade-*/submission'):
    try:
     st=p.stat()
     if st.st_size and os.access(p,os.X_OK):
      k=(str(p),st.st_size,st.st_mtime_ns)
      if k not in observed:observed[k]={'phase':phase,'path':str(p),**ident(p),'observed_epoch':time.time(),'scope':'non-exhaustive transient executable sampling; local/hash-only'}
    except FileNotFoundError:pass
  except Exception as e:write('sampler-error.json',{'error':repr(e)})
def run(label,argv,deadline):
 global phase
 phase=label;assert shutil.disk_usage(R).free>=1024**3;start=time.monotonic();p=subprocess.Popen(argv,cwd=R,env=env,stdin=subprocess.DEVNULL,stdout=subprocess.PIPE,stderr=subprocess.PIPE,start_new_session=True);timeout=False
 try:out,err=p.communicate(timeout=deadline)
 except subprocess.TimeoutExpired:
  timeout=True;os.killpg(p.pid,signal.SIGTERM)
  try:out,err=p.communicate(timeout=5)
  except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);out,err=p.communicate()
 (E/(label+'.stdout.txt')).write_bytes(out);(E/(label+'.stderr.txt')).write_bytes(err)
 executable=Path(shutil.which(argv[0],path=env['PATH']) or argv[0]);row={'label':label,'argv':argv,'cwd':str(R),'environment':env,'pid':p.pid,'process_group':p.pid,'executable':{'path':str(executable),**ident(executable)},'deadline_seconds':deadline,'duration_seconds':time.monotonic()-start,'exit_code':p.returncode,'timed_out':timeout,'stdout':{'path':label+'.stdout.txt',**digest(out)},'stderr':{'path':label+'.stderr.txt',**digest(err)}};receipts.append(row);write('command-receipts.json',receipts);print(label,'EXIT',p.returncode,'SECONDS',round(row['duration_seconds'],3),flush=True)
 return row,out,err

try:
 E.mkdir(parents=True,exist_ok=True);before=verify();write('verification-before.json',before);old=source_and_carry();assert shutil.disk_usage(R).free>=1024**3
 W=Path(tempfile.mkdtemp(prefix='cadence-ch10-prompt-fixture-review-'));(W/'home').mkdir();(W/'tmp').mkdir()
 env={'PATH':'/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(W/'home'),'TMPDIR':str(W/'tmp'),'LANG':'en_US.UTF-8','GOMAXPROCS':'2','GOFLAGS':'-mod=readonly -p=1','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOCACHE':'/Users/bill/Library/Caches/go-build','GOMODCACHE':'/Users/bill/go/pkg/mod','GOPATH':'/Users/bill/go','CGO_ENABLED':'1'}
 assert all(k not in env for k in ['LLM_API_KEY','OPENAI_API_KEY','ANTHROPIC_API_KEY','GOOGLE_API_KEY','LLM_BASE_URL','EN_PRIMARY_SKILL','EN_SKILLS_DIR'])
 go=shutil.which('go',path=env['PATH']);dlv=shutil.which('dlv',path=env['PATH']);assert go and dlv;write('environment.json',env);write('workspace.json',{'owned_scratch':str(W),'go':{'path':go,**ident(go)},'dlv':{'path':dlv,**ident(dlv)},'disk_free_bytes_before':shutil.disk_usage(R).free,'no_inherited_provider_credentials_or_settings':True,'real_launch_fake_provider_environment':'unchanged grader supplies loopback URLs and dummy keys to its child clients; no paid provider'})
 oldcmd=json.loads((N/'ch10-student-evidence/grader-ch10-final-v1.json').read_text());full=oldcmd['argv'].copy();assert full[0]=='/tmp/cadence-ch10-grader';full[0]=F['binary']['path'];assert full[-1]==str(S/'.ch10-local/bin/frozen-cli')
 test=['go','test','-count=1','-run','^TestChapter10Persistence$','-v','-timeout=180s','./internal/gradeedition2']
 write('control-plan.json',{'order':[{'label':'01-corrected-full-chapter10','argv':full,'deadline_seconds':600,'expected_title':TITLE,'expected_ids':IDS,'expected_points':POINTS,'required_exit':0,'required_score':100,'required_passed':True},{'label':'02-historical-focused-persistence','argv':test,'deadline_seconds':240,'required_exit':0,'required_run_and_pass':'TestChapter10Persistence','authorized_root_message':'main-31','only_if_first_control_passes':True}],'no_other_runtime_controls':True,'chapter1':'NOT RUN','old_negative_carried_not_replayed':True})
 write('prerequisite-and-selector-ledger.json',{'all_current_inputs_and_previous_seals_verified':True,'full_historical_governing_reads_carried_by_unchanged_hash':True,'changed_checker_and_new_parent_decision_verification_fully_read':True,'root_task_selector_error':'TestChapter10CLISavesAndLoads absent from pinned test file','actual_only_function':'TestChapter10Persistence','root_authorization':'main-31: exact -run ^TestChapter10Persistence$, only after student gate passes; require actual RUN/PASS','no_empty_selector_run':True,'root_compile_race_vet_build_whitespace_receipts_verified_not_replayed':True})
 sampler=threading.Thread(target=sample,daemon=True);sampler.start()
 row,out,err=run('01-corrected-full-chapter10',full,600)
 try:parsed=json.loads(out)
 except Exception as e:write('01-parse-failure.json',{'exception':repr(e)});raise
 write('01-parsed-report.json',parsed)
 expected=parsed.get('title')==TITLE and [x['id'] for x in parsed.get('checks',[])]==IDS and [x['points'] for x in parsed.get('checks',[])]==POINTS and parsed.get('score')==100 and parsed.get('max_score')==100 and parsed.get('passed') is True and all(x['passed'] is True and x['earned']==x['points'] for x in parsed['checks'])
 write('01-report-verdict.json',{'process_exit':row['exit_code'],'timed_out':row['timed_out'],'parsed_report_complete_green':expected,'accepted_runtime_result':row['exit_code']==0 and not row['timed_out'] and expected,'old_80_report_replaced':False,'same_full_argv_except_grader_executable':full[1:]==oldcmd['argv'][1:]})
 assert row['exit_code']==0 and not row['timed_out'] and expected,'corrected full Chapter10 did not fully pass; no further controls'
 row,out,err=run('02-historical-focused-persistence',test,240);text=out.decode();real='=== RUN   TestChapter10Persistence' in text and '--- PASS: TestChapter10Persistence' in text and '\nPASS\n' in text and 'no tests to run' not in text;write('02-test-verdict.json',{'exit_code':row['exit_code'],'timed_out':row['timed_out'],'actual_named_test_ran_and_passed':real,'only_exact_selector':True,'both_historical_clients_evidenced_by_unchanged_test_source':True});assert row['exit_code']==0 and not row['timed_out'] and real,'focused historical compatibility failed or did not execute'
except BaseException as e:
 failure={'exception':repr(e),'traceback':traceback.format_exc(),'no_retry_or_tuning':True,'remaining_controls_not_run':len(receipts)<2};write('failure.json',failure);print('STOP',repr(e),flush=True)
finally:
 stop.set()
 if sampler:sampler.join()
 write('observed-transient-executable-identities.json',list(observed.values()))
 tail={'failure':failure,'completed_runtime_commands':len(receipts),'all_top_level_jobs_waited':True,'sampler_joined':sampler is None or not sampler.is_alive(),'chapter1':'NOT RUN','owned_scratch':str(W) if W else None,'disk_free_bytes_after':shutil.disk_usage(R).free}
 try:
  after=verify();write('verification-after.json',after);tail['before_after_exactly_equal']=after==before;assert after==before
 except BaseException as e:tail['preservation_verification_error']=repr(e)
 if W:
  ps=subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True);tail['owned_processes_remaining']=[s.strip() for s in ps.splitlines() if str(W) in s];tail['owned_tmp_entries']=[p.name for p in (W/'tmp').iterdir()]
 write('completion.json',tail);print(json.dumps({'completed_runtime_commands':len(receipts),'failure':failure,'before_after_exactly_equal':tail.get('before_after_exactly_equal'),'owned_processes_remaining':tail.get('owned_processes_remaining'),'owned_tmp_entries':tail.get('owned_tmp_entries')}),flush=True)
