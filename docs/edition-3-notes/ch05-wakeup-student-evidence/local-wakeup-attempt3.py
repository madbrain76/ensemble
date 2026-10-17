#!/usr/bin/env python3
"""Loopback-only executable wakeup proof. No credentials or provider route."""
import hashlib,http.server,json,os,pathlib,re,subprocess,sys,threading,time,traceback
E=pathlib.Path('/tmp/cadence-ch05-wakeup-evidence');OUT=pathlib.Path(sys.argv[1]);OUT.mkdir(parents=True,exist_ok=False)
identity=json.loads((E/'candidate-r2/frozen/identity.json').read_text());sha=lambda p:hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
for n,h in identity['source_files'].items():assert sha(pathlib.Path(identity['source_root'])/n)==h
for b in identity['binaries'].values():assert sha(b['path'])==b['sha256']
assert not subprocess.check_output(['git','status','--porcelain'],cwd=identity['source_root'])
def rows(lines):
 out=[]
 for line in lines:
  try:out.append(json.loads(line))
  except ValueError:pass
 return out
def wait(pred,label,seconds=5):
 end=time.monotonic()+seconds
 while time.monotonic()<end:
  if pred():return
  time.sleep(.005)
 raise AssertionError(label)
class Process:
 def __init__(self,binary,name,url):
  self.root=OUT/name;self.root.mkdir();self.stdout=[];self.stderr=[]
  env={k:v for k,v in os.environ.items() if not(k.startswith(('LLM_','ANTHROPIC_','OPENAI_','GEMINI_','CH02_','CH06_','ENSEMBLE_')) or any(w in k.upper() for w in ('KEY','TOKEN','SECRET','PASSWORD')))}
  env.update(LLM_VENDOR='anthropic',LLM_API_KEY='local-fake-only',LLM_MODEL='claude-sonnet-4-6',LLM_BASE_URL=url,LLM_REQUIRE_TOOL_EXPLANATION='true',ENSEMBLE_WORKSPACE=str(self.root),ENSEMBLE_DATA_DIR=str(self.root/'runtime'),CH02_LOG=str(self.root/'history.jsonl'),CH06_LOG=str(self.root/'workflow.jsonl'))
  self.start=time.monotonic();self.p=subprocess.Popen([binary],cwd=self.root,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True,bufsize=1)
  def drain(stream,target,path):
   with path.open('w') as f:
    for line in stream:target.append(line);f.write(line);f.flush()
  self.threads=[threading.Thread(target=drain,args=(self.p.stdout,self.stdout,self.root/'stdout.txt')),threading.Thread(target=drain,args=(self.p.stderr,self.stderr,self.root/'stderr.txt'))]
  for t in self.threads:t.start()
 def send(self,data):self.p.stdin.write(json.dumps(data)+'\n');self.p.stdin.flush()
 def finish(self):
  self.p.stdin.close();code=self.p.wait(timeout=25)
  for t in self.threads:t.join()
  receipt={'binary':self.p.args[0],'binary_sha256':sha(self.p.args[0]),'exit':code,'elapsed_seconds':time.monotonic()-self.start,'http':'loopback fake only','paid_requests':0};(self.root/'exit.json').write_text(json.dumps(receipt,indent=2)+'\n');assert code==0,(code,self.stderr)
 def cleanup(self):
  if self.p.poll() is None:
   try:self.p.stdin.close();self.p.wait(timeout=8)
   except (subprocess.TimeoutExpired,BrokenPipeError):self.p.terminate();self.p.wait(timeout=5)
  for t in self.threads:t.join()
def text(s):return {'type':'text','text':s}
def call(id,name,args):return {'type':'tool_use','id':id,'name':name,'input':args}
def paired(body):
 out={}
 for m in body['messages']:
  if isinstance(m.get('content'),list):
   for p in m['content']:
    if p['type']=='tool_result':assert p['tool_use_id'] not in out;out[p['tool_use_id']]=p
 return out
def content(p):return json.dumps(p.get('content'))
class Boundary(http.server.BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_POST(self):
  try:
   body=json.loads(self.rfile.read(int(self.headers['Content-Length'])));self.server.bodies.append(body);n=len(self.server.bodies)
   (self.server.dest/f'request-{n}.json').write_text(json.dumps(body,indent=2)+'\n')
   parts=self.server.respond(body);data=json.dumps({'model':'claude-sonnet-4-6','content':parts,'usage':{'input_tokens':1,'output_tokens':1},'stop_reason':'tool_use' if any(p['type']=='tool_use' for p in parts) else 'end_turn'}).encode();self.send_response(200)
  except Exception as e:
   self.server.errors.append(traceback.format_exc());data=json.dumps({'error':str(e)}).encode();self.send_response(500)
  self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(data)));self.end_headers();self.wfile.write(data)
def serve(dest,respond):
 dest.mkdir();s=http.server.ThreadingHTTPServer(('127.0.0.1',0),Boundary);s.bodies=[];s.errors=[];s.respond=respond;s.dest=dest;threading.Thread(target=s.serve_forever,daemon=True).start();return s
results=[]
def cli():
 state={'id':None,'times':{}};p=None
 def reply(b):
  n=len(s.bodies);state['times'][str(n)]=time.monotonic();r=paired(b)
  if n==1:return [text('I will start a held process.'),call('launch','run_command',{'command':"printf ready > ready; while [ ! -f release ]; do sleep .02; done; printf 'ACTUAL_PROCESS_FINAL\\n'",'ai_callback_delay':3600,'max_output_bytes':512}),call('skipped','write_file',{'path':'must-not-exist','content':'wrong'})]
  if n==2:
   assert 'CLI_INITIAL_WAKE' in json.dumps(b);assert not r['launch'].get('is_error');m=re.search(r'job (\d+) running',content(r['launch']));assert m;state['id']=int(m[1]);assert r['skipped']['is_error'] and 'not executed' in content(r['skipped']);assert not (p.root/'release').exists();return [text('I will answer while the process continues.'),call('answer','write_file',{'path':'answer.txt','content':'ANSWER_BEFORE_JOB_END'})]
  if n==3:
   assert not r['answer'].get('is_error');assert not (p.root/'release').exists();return [text('I will observe the same job.'),call('later','wait_for_job',{'handle':state['id'],'ai_callback_delay':3600})]
  if n==4:
   assert 'CLI_LATER_WAKE' in json.dumps(b);assert not r['later'].get('is_error');assert f"job {state['id']} running" in content(r['later']);assert 'CLI_INITIAL_WAKE' not in json.dumps(b);assert not (p.root/'release').exists();(p.root/'release').touch();return [text('I will collect the same job.'),call('final','wait_for_job',{'handle':state['id'],'ai_callback_delay':3600})]
  assert n==5 and 'ACTUAL_PROCESS_FINAL' in content(r['final']);assert f"job {state['id']} done" in content(r['final']);return [text('CLI_WAKEUP_FINISHED')]
 s=serve(OUT/'cli-http',reply)
 try:
  p=Process(identity['binaries']['cli']['path'],'cli',f'http://127.0.0.1:{s.server_port}');p.send({'kind':'prompt','text':'start held process'})
  wait(lambda:(p.root/'ready').exists(),'real process ready');t1=time.monotonic();p.send({'kind':'hint','text':'CLI_INITIAL_WAKE'})
  wait(lambda:'2' in state['times'] and state['id'] is not None,'initial observation next HTTP');assert state['times']['2']-t1<2
  wait(lambda:sum(x.get('observation')=='state_changed' and x.get('to')=='tools_pending' for x in rows(p.stdout))>=3,'later wait admitted');t2=time.monotonic();p.send({'kind':'hint','text':'CLI_LATER_WAKE'})
  wait(lambda:len(s.bodies)>=5,'later wait next HTTP and real terminal result');assert state['times']['4']-t2<2
  p.finish();assert not s.errors,s.errors;assert not (p.root/'must-not-exist').exists();assert (p.root/'answer.txt').read_text()=='ANSWER_BEFORE_JOB_END'
  logs=[json.loads(x) for x in (p.root/'history.jsonl').read_text().splitlines()];returns=[x for x in logs if x.get('type')=='tool_returned'];assert len(returns)==5 and len({x['tool']['call_id'] for x in returns})==5
  files=list((p.root/'runtime/cr/io').iterdir());assert len(files)==2;assert 'ACTUAL_PROCESS_FINAL' in (p.root/f"runtime/cr/io/{state['id']}").read_text()
  results.append({'scenario':'CLI initial and later 3600s observations','passed':True,'initial_wake_seconds':state['times']['2']-t1,'later_wake_seconds':state['times']['4']-t2,'same_job':state['id'],'retained_jobs':2,'requests':len(s.bodies),'EOF_saved':True,'no_tail_effect':True})
 finally:
  if p:(p.root/'release').touch();p.cleanup()
  s.shutdown()
def workflow():
 counts={};state={};p=None
 def reply(b):
  sysmsg=json.dumps(b.get('system',''));role='author' if 'You are the author' in sysmsg else 'editor' if 'You are the editor' in sysmsg else 'reviewer';counts[role]=counts.get(role,0)+1;n=counts[role];r=paired(b)
  if role=='author':
   if n==1:return [text('I will think before drafting.'),call('limit','tool_limits',{'ai_callback_delay':3600}),call('held','think',{'seconds':4,'thought':'ACTUAL_ORDINARY_THOUGHT'}),call('skipped','write_draft',{'text':'WRONG_UNDISPATCHED_DRAFT'})]
   if n==2:
    state['wake_http']=time.monotonic();assert 'WORKFLOW_WAKE' in json.dumps(b);assert not r['held'].get('is_error');m=re.search(r'job (\d+) running',content(r['held']));assert m;state['id']=int(m[1]);assert r['skipped']['is_error'] and 'not executed' in content(r['skipped']);assert not (p.root/'manuscript.txt').exists();return [text('I will write now while thinking continues.'),call('write','write_draft',{'text':'New draft after waking the foreground observation.'}),call('count','word_count',{})]
   assert n==3 and not r['write'].get('is_error');return [text('AUTHOR_WAKE_FINISHED')]
  if role=='editor':return [text('I will read and edit.'),call('read','read_draft',{}),call('edit','edit_draft',{'text':'Edited wakeup draft with real application-owned handoffs.'})] if n==1 else [text('EDITOR_FINISHED')]
  return [text('I will record acceptance.'),call('review','review',{'decision':'accept','notes':'Verified actual revised wakeup draft.'})] if n==1 else [text('REVIEWER_FINISHED')]
 s=serve(OUT/'workflow-http',reply)
 try:
  p=Process(identity['binaries']['workflow']['path'],'workflow',f'http://127.0.0.1:{s.server_port}');p.send({'kind':'start','text':'write wakeup example'})
  wait(lambda:(p.root/'runtime/cr/io/1').exists(),'ordinary handler observation admitted');sent=time.monotonic();p.send({'kind':'hint','agent':'author','text':'WORKFLOW_WAKE'})
  wait(lambda:'wake_http' in state,'ordinary handler next HTTP');assert state['wake_http']-sent<2
  wait(lambda:any(x.get('pipeline')=='done' for x in rows(p.stdout)),'workflow pipeline done');eof=time.monotonic();assert (p.root/'runtime/cr/io/1').read_text()=='';p.finish();assert time.monotonic()-eof>.5
  assert not s.errors,s.errors;assert 'ACTUAL_ORDINARY_THOUGHT' in (p.root/'runtime/cr/io/1').read_text();assert (p.root/'manuscript.txt').read_text()=='Edited wakeup draft with real application-owned handoffs.';assert json.loads((p.root/'review.json').read_text())['decision']=='accept'
  events=[json.loads(x) for x in (p.root/'workflow.jsonl.author.jsonl').read_text().splitlines()];assert any(x.get('type')=='job_reported' and x['job']['handle']==state['id'] and x['job']['status']=='done' for x in events)
  results.append({'scenario':'ordinary workflow handler wake and EOF join','passed':True,'wake_seconds':state['wake_http']-sent,'same_job':state['id'],'requests':len(s.bodies),'pipeline_before_old_handler_terminal':True,'EOF_joined_real_output':True,'no_undispatched_draft':True})
 finally:
  if p:p.cleanup()
  s.shutdown()
try:
 cli();workflow();(OUT/'summary.json').write_text(json.dumps({'source_commit':identity['source_commit'],'binaries':identity['binaries'],'passed':True,'scenarios':results,'paid_requests':0},indent=2)+'\n')
except Exception as e:
 (OUT/'failure.json').write_text(json.dumps({'error':str(e),'traceback':traceback.format_exc(),'completed':results},indent=2)+'\n');raise
