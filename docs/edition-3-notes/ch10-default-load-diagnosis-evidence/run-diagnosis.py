#!/usr/bin/env python3
"""Three-launch maximum, read-only-source default-load diagnosis. No build or grade."""
from pathlib import Path
import collections,hashlib,http.server,json,os,shutil,signal,subprocess,tempfile,threading,time,traceback
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');E=R/'docs/edition-3-notes/ch10-default-load-diagnosis-evidence';M=json.loads((R/'docs/edition-3-notes/ch10-default-load-diagnosis-input.json').read_text());C=M['candidate']
def ident(p):
 b=Path(p).read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def blob(b):return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def write(n,x):
 p=E/n;assert not p.exists(),p;p.write_text(json.dumps(x,indent=2,ensure_ascii=False)+'\n')
def git(root,*args):return subprocess.check_output(['git','--no-pager',*args],cwd=root)
def verify():
 out={'inputs':[],'protected':[],'candidate_files':[],'candidate_binaries':[],'prior_inventory_seals':{},'repositories':{}}
 for x in M['inputs']:
  h=ident(x['path']);assert h=={k:x[k] for k in h};out['inputs'].append(x)
 for rel,want in M['protected_tracked'].items():assert ident(R/rel)==want;out['protected'].append({'path':rel,**want})
 for x in C['tracked_files']:assert ident(S/x['path'])=={k:x[k] for k in ['bytes','sha256']};out['candidate_files'].append(x)
 for x in C['binaries']:assert ident(x['path'])=={k:x[k] for k in ['bytes','sha256']};out['candidate_binaries'].append(x)
 for p in [C['archive']['local_path'],C['archive']['evidence_path']]:assert ident(p)=={k:C['archive'][k] for k in ['bytes','sha256']}
 for key,rev in [('head','HEAD'),('tree','HEAD^{tree}'),('parent','HEAD^')]:assert git(S,'rev-parse',rev).decode().strip()==C[key]
 assert git(S,'status','--porcelain=v1')==b''
 sealpath=R/'docs/edition-3-notes/ch10-student-evidence/seal.json';seal=json.loads(sealpath.read_text());rows=seal['reports']+seal['evidence']
 for x in rows:assert ident(x['path'])=={k:x[k] for k in ['bytes','sha256']},x['path']
 assert {p.resolve().as_posix() for p in sealpath.parent.rglob('*') if p.is_file() and p!=sealpath}=={x['path'] for x in seal['evidence']}
 out['blocked_seal']={'identity':ident(sealpath),'member_count':len(rows),'closed_evidence_set':True,'all_match':True,'chapter10':seal['chapter10'],'chapter1':seal['chapter1']}
 for name in ['ch10-entry-audit-evidence','ch10-entry-followup-evidence','ch10-grader-review-evidence','ch10-grader-review-followup-evidence']:
  d=R/'docs/edition-3-notes'/name;I=d/'inventory.json';j=json.loads(I.read_text())
  for x in j['files']:assert ident(R/x['path'])=={k:x[k] for k in ['bytes','sha256']},x['path']
  outside={x['path'] for x in j['files'] if not x['path'].startswith(d.relative_to(R).as_posix()+'/')};actual={p.relative_to(R).as_posix() for p in d.rglob('*') if p.is_file() and p!=I}|outside;assert actual=={x['path'] for x in j['files']}
  out['prior_inventory_seals'][name]={'identity':ident(I),'member_count':len(j['files']),'exact_closed_set':True}
 for name,root in [('main',R),('student',S)]:out['repositories'][name]={'head':git(root,'rev-parse','HEAD').decode().strip(),'tree':git(root,'rev-parse','HEAD^{tree}').decode().strip(),'unstaged_diff':blob(git(root,'diff','--binary')),'staged_diff':blob(git(root,'diff','--cached','--binary')),'tracked_status':blob(git(root,'status','--porcelain=v1','--untracked-files=no'))}
 assert shutil.disk_usage(R).free>=1024**3
 return out
before=verify();write('verification-before.json',before)
W=Path(tempfile.mkdtemp(prefix='cadence-ch10-default-load-diagnosis-'));(W/'home').mkdir();(W/'tmp').mkdir();(W/'skills/current').mkdir(parents=True);(W/'cli').mkdir();(W/'gui').mkdir();write('workspace.json',{'owned_scratch':str(W),'no_existing_sources_or_fixtures_written':True})
B={}
for b in C['binaries']:
 n=Path(b['path']).name.lower()
 if 'gui' in n:B['gui']=b
 elif 'cli' in n or 'ch01' in n:B['cli']=b
assert set(B)=={'cli','gui'},C['binaries']
head='---\nname: current\ntype: primary\ndescription: current runtime primary\ntools: read_file, load_skill, unload_skill\n---\n';body='CURRENT_RUNTIME_PROMPT\n';skill=W/'skills/current/SKILL.md';skill.write_text(head+body)
fromv={'vendor':'anthropic','model':'fake-model','surface':'messages'};dialogue=[]
for i,word in enumerate(['FIRST_WORD','SECOND_WORD','THIRD_WORD']):
 dialogue.append({'seq':2*i+1,'actor':'human','parts':[{'type':'text','text':word,'from':fromv}]});dialogue.append({'seq':2*i+2,'actor':'agent','parts':[{'type':'text','text':'PERSISTED_REPLY_'+word,'from':fromv}]})
seed={'config':{'model':'fake-model','vendor':'anthropic','system_prompt':'RECORDED_PROMPT_MUST_NOT_RESTORE','tools':[{'name':'write_file','description':'Recorded-only declaration; no handler.','schema':{'type':'object','properties':{}}}]},'as_of':6,'context':{'turn':'idle','dialogue':dialogue,'ephemera':[],'usage':{'input':0,'cache_write':0,'cache_read':0,'output':0}},'log':[]};seedbytes=(json.dumps(seed,separators=(',',':'))+'\n').encode()
(E/'snapshot-fixture.txt').write_bytes(seedbytes);(E/'primary-fixture-with-lf.txt').write_text(head+body);write('control-plan.json',{'scope':'default-load configuration mismatch only','maximum_real_client_launches':3,'builds_or_grades_or_package_suites':0,'first':'frozen CLI default-load, literal-LF primary','second':'frozen GUI default-load, literal-LF primary, -port 0','conditional_third':'same CLI binary/cwd/env/save bytes/prompt/reply, remove only terminal LF from owned primary fixture','request_model':'gpt-5-course-fast','provider':'loopback OpenAI chat completions','historical_expected_system_exact':'CURRENT_RUNTIME_PROMPT','fixture_origin':'Handauthored minimal public four-field snapshot-only archive, not a captured historical graded save. as_of=6, 6 turns, log=[]; deliberately opposing recorded config.','baseline_seed_identity':blob(seedbytes),'source_or_grader_mutations':0,'explanation':'This isolates the current-config phase without replaying seed/restart campaigns. Existing score remains 80/100 exit1; Chapter1 remains NOT RUN.'})
write('prerequisite-and-setup-ledger.json',{'carried_full_reads':'Unchanged manifest-pinned governing/craft/historical documents compared with prior review input. Current changed execution/student guides and blocked records read. Chapter10 and relevant exact source reread.','source_manifest_equals_input_candidate':json.loads((R/'docs/edition-3-notes/ch10-student-source-manifest.json').read_text())==C,'read_only_setup_errors':['Incorrect assumed renderer filename; corrected using frozen tracked-file list.','Incorrect assumed common history/vocabulary filenames; actual common.go inspected.','Assumed seal.files key; actual reports/evidence schema inspected.','Oversized combined source/seal outputs truncated; needed code ranges reread in bounded sections.'],'runtime_retries':0})
requests=[];receipts=[];obs=[];errors=[];lock=threading.Lock();active=None;server=None;thread=None;failure=None
class Handler(http.server.BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_POST(self):
  try:
   raw=self.rfile.read(int(self.headers.get('Content-Length','0')));p=json.loads(raw);assert self.path=='/v1/chat/completions',self.path;assert p.get('stream',False) is False
   with lock:
    n=len(requests)+1;local=W/('raw-request-%02d.json'%n);local.write_bytes(raw);requests.append({'method':'POST','path':self.path,'auth_is_expected_dummy':self.headers.get('Authorization')=='Bearer diagnosis-dummy-key','raw_local_path':str(local),'raw_identity':blob(raw),'body':p})
   reply={'id':'diagnosis-reply','object':'chat.completion','created':0,'model':p['model'],'choices':[{'index':0,'message':{'role':'assistant','content':'DIAGNOSIS_OK'},'finish_reason':'stop'}],'usage':{'prompt_tokens':40,'completion_tokens':4,'total_tokens':44}};data=json.dumps(reply,separators=(',',':')).encode();self.send_response(200);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(data)));self.end_headers();self.wfile.write(data)
  except Exception as ex:
   errors.append(repr(ex));self.send_error(400)
def text(content):
 if isinstance(content,str):return content
 if isinstance(content,list):return ''.join(x.get('text','') for x in content if isinstance(x,dict) and x.get('type')=='text')
 return ''
def env_for(kind,base):
 d=W/kind
 return {'PATH':'/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(W/'home'),'TMPDIR':str(W/'tmp'),'LANG':'C','LC_ALL':'C','TZ':'UTC','GOMAXPROCS':'2','GOFLAGS':'-mod=readonly -p=1','LLM_VENDOR':'openai','LLM_MODEL':'gpt-5-course-fast','LLM_BASE_URL':base,'LLM_API_KEY':'diagnosis-dummy-key','OPENAI_BASE_URL':base,'OPENAI_API_KEY':'diagnosis-dummy-key','EN_DISABLE_STREAMING':'1','EN_SKILLS_DIR':str(W/'skills'),'EN_PRIMARY_SKILL':'current','ENSEMBLE_DATA_DIR':str(d),'CH02_LOG':str(d/'events.jsonl')}
def launch(label,kind,env):
 global active
 assert len(receipts)<3;assert not errors
 d=W/kind;(d/'save.json').write_bytes(seedbytes);inputsha=ident(d/'save.json');skillsha=ident(skill);argv=[B[kind]['path']]+(['-port','0'] if kind=='gui' else []);before_n=len(requests);start=time.monotonic();active=subprocess.Popen(argv,cwd=d,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,start_new_session=True);pid=active.pid;timed=False
 try:out,err=active.communicate(b'{"user":"SWITCH_WORD"}\n',timeout=60)
 except subprocess.TimeoutExpired:
  timed=True;os.killpg(pid,signal.SIGTERM)
  try:out,err=active.communicate(timeout=5)
  except subprocess.TimeoutExpired:os.killpg(pid,signal.SIGKILL);out,err=active.communicate()
 rc=active.returncode;active=None;elapsed=time.monotonic()-start;(E/(label+'.stdout.txt')).write_bytes(out);(E/(label+'.stderr.txt')).write_bytes(err)
 receipt={'label':label,'argv':argv,'cwd':str(d),'environment':env,'pid':pid,'process_group':pid,'executable_identity':ident(argv[0]),'stdin_utf8':'{"user":"SWITCH_WORD"}\n','input_save_identity':inputsha,'primary_fixture_identity':skillsha,'duration_seconds':elapsed,'deadline_seconds':60,'timed_out':timed,'process_exit':rc,'stdout':{'path':label+'.stdout.txt',**blob(out)},'stderr':{'path':label+'.stderr.txt',**blob(err)}};receipts.append(receipt);write(label+'.receipt.json',receipt);assert rc==0 and not timed,(label,rc);assert not errors;assert len(requests)==before_n+1,(label,len(requests)-before_n)
 q=requests[before_n];p=q['body'];systems=[text(x.get('content')) for x in p['messages'] if x.get('role')=='system'];system='\n'.join(systems);users=[text(x.get('content')) for x in p['messages'] if x.get('role')=='user'];tools=sorted(x['function']['name'] for x in p.get('tools',[]));saved=(d/'save.json').read_bytes();(W/(label+'-saved.json')).write_bytes(saved);s=json.loads(saved);c=s['config'];config_tools=sorted(x['name'] for x in c['tools']);messages=[{'role':x.get('role'),'text':text(x.get('content'))} for x in p['messages']];want='CURRENT_RUNTIME_PROMPT';wire={'method':q['method'],'path':q['path'],'provider_evidenced_by':'OpenAI chat-completions endpoint and message/tool wire shape; no vendor JSON field is asserted','auth_is_expected_dummy':q['auth_is_expected_dummy'],'raw_request_local_only':q['raw_local_path'],'raw_request_identity':q['raw_identity'],'model':p.get('model'),'stream':p.get('stream',False),'system_messages':systems,'system_text':system,'system_utf8_hex':system.encode().hex(),'system_utf8_bytes':len(system.encode()),'system_has_one_terminal_lf':system==want+'\n','tool_names':tools,'sanitized_messages':messages,'user_texts':users,'historical_model_predicate':p.get('model')=='gpt-5-course-fast','historical_exact_system_predicate':system==want,'historical_conjunction_predicate':p.get('model')=='gpt-5-course-fast' and system==want,'desired_history_counts':{word:users.count(word) for word in ['FIRST_WORD','SECOND_WORD','THIRD_WORD','SWITCH_WORD']}}
 save_summary={'local_only_capture':str(W/(label+'-saved.json')),'identity':blob(saved),'top_level_fields':sorted(s),'config':{'model':c['model'],'vendor':c['vendor'],'system_prompt':c['system_prompt'],'system_utf8_hex':c['system_prompt'].encode().hex(),'system_utf8_bytes':len(c['system_prompt'].encode()),'tool_names':config_tools},'as_of':s['as_of'],'log_entry_count':len(s['log']),'context_turn_count':len(s['context']['dialogue']),'config_matches_actual_request':c['model']==p['model'] and c['vendor']=='openai' and c['system_prompt']==system and config_tools==tools,'recorded_seed_config_not_restored':c['vendor']!='anthropic' and c['model']!='fake-model' and c['system_prompt']!='RECORDED_PROMPT_MUST_NOT_RESTORE'}
 observation={'label':label,'client':kind,'process_exit':rc,'wire':wire,'rewritten_save':save_summary};obs.append(observation);write(label+'.sanitized-observation.json',observation);print(label,'exit',rc,'seconds',round(elapsed,3),'model',p.get('model'),'system',repr(system),'tools',tools,'history',wire['desired_history_counts'],'saved_current',save_summary['config_matches_actual_request'],flush=True)
 assert all(n==1 for n in wire['desired_history_counts'].values());assert save_summary['config_matches_actual_request'];assert save_summary['recorded_seed_config_not_restored'];assert tools==['load_skill','read_file','unload_skill'];assert q['auth_is_expected_dummy'];assert wire['historical_model_predicate']
 return observation
try:
 server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);server.daemon_threads=False;server.block_on_close=True;thread=threading.Thread(target=server.serve_forever,name='owned-loopback-vendor');thread.start();base='http://127.0.0.1:'+str(server.server_port);write('loopback-server.json',{'pid':os.getpid(),'listen_address':server.server_address,'base_url':base,'allowlisted_route':'/v1/chat/completions','reply_text':'DIAGNOSIS_OK','no_request_headers_exported':True,'authorization_checked_only_against_fixed_dummy':True})
 envs={k:env_for(k,base) for k in ['cli','gui']};write('launch-environments.json',envs)
 first=launch('01-cli-default-load','cli',envs['cli']);second=launch('02-gui-default-load','gui',envs['gui'])
 if all(x['wire']['system_has_one_terminal_lf'] and not x['wire']['historical_conjunction_predicate'] for x in [first,second]):
  before_skill=skill.read_bytes();assert before_skill.endswith(b'\n');skill.write_bytes(before_skill[:-1]);(E/'primary-fixture-without-lf.txt').write_bytes(before_skill[:-1]);write('contrast-definition.json',{'one_input_change':'Remove final byte 0a from owned primary SKILL.md; all preceding bytes identical.','before':blob(before_skill),'after':ident(skill),'same_binary_cwd_environment_stdin_seed_save_and_fake_reply_as':'01-cli-default-load','not_changed':'No source, grader, model, vendor, tool list, source archive, or previous evidence changed.'})
  contrast=launch('03-cli-without-terminal-lf','cli',envs['cli']);assert contrast['wire']['historical_conjunction_predicate'];assert receipts[0]['environment']==receipts[2]['environment'] and receipts[0]['cwd']==receipts[2]['cwd'] and receipts[0]['argv']==receipts[2]['argv'] and receipts[0]['input_save_identity']==receipts[2]['input_save_identity'];skill.write_bytes(before_skill)
 else:write('contrast-not-run.json',{'reason':'Observed systems did not match the sole source-predicted LF mismatch; no opportunistic control substituted.'})
except Exception:
 failure=traceback.format_exc();(E/'runtime-failure.txt').write_text(failure);print(failure,flush=True)
finally:
 if active is not None and active.poll() is None:
  os.killpg(active.pid,signal.SIGKILL);active.communicate()
 if server is not None:server.shutdown();server.server_close()
 if thread is not None:thread.join(timeout=10);assert not thread.is_alive()
 write('command-receipts.json',receipts);write('observations.json',obs);write('server-errors.json',errors)
 try:
  after=verify();write('verification-after.json',after);assert before==after
 except Exception:
  failure=(failure or '')+'\nFINAL VERIFY '+traceback.format_exc();(E/'final-verification-failure.txt').write_text(failure)
 live=[x.strip() for x in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True).splitlines() if str(W) in x];localpins=[]
 for p in sorted(W.rglob('*')):
  if p.is_file():localpins.append({'path':str(p),**ident(p)})
 write('local-only-artifact-identities.json',localpins);write('run-outcome.json',{'failure':failure,'all_jobs_waited':True,'real_client_launches':len(receipts),'server_request_count':len(requests),'server_closed':server is not None,'server_thread_joined':thread is not None and not thread.is_alive(),'owned_processes_remaining':live,'owned_tmp_entries':[p.name for p in (W/'tmp').iterdir()],'source_and_previous_evidence_unchanged':before==after if 'after' in locals() else False,'owned_skill_restored_to_initial_lf':skill.read_text()==head+body,'disk_free_bytes':shutil.disk_usage(R).free,'full_grade_runs':0,'chapter1':'NOT RUN','candidate_remains_blocked':True});assert not live
 if failure:raise SystemExit(1)
