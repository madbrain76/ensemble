#!/usr/bin/env python3
"""Actual executable R1 regression against local-only HTTP, with real role tools."""
import codecs,hashlib,http.server,json,os,pathlib,shutil,subprocess,tempfile,threading,time
P=pathlib.Path('/tmp/cadence-ch05-correction-r1-evidence');manifest=json.loads((P/'frozen/identity.json').read_text())
OUT=P/'local-interface-r2';OUT.mkdir(exist_ok=False)
HINT='R1_WIRE_HINT_71d';AUTHOR='The quiet library welcomed tired readers.';EDITED='The quiet library welcomed tired readers and offered warm shelter.';REVIEW={'decision':'accept','notes':'The edited manuscript is clearer.'}
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=manifest['source_root'],text=True).strip()==manifest['source_commit']
assert not subprocess.check_output(['git','status','--porcelain'],cwd=manifest['source_root'])
for name,digest in manifest['source_files'].items():assert sha(pathlib.Path(manifest['source_root'])/name)==digest
for item in manifest['binaries'].values():assert sha(pathlib.Path(item['path']))==item['sha256']
class Boundary:
 def __init__(self,mode):
  self.mode=mode;self.requests=[];self.roles={};self.entered=threading.Event();self.release=threading.Event();owner=self
  class Handler(http.server.BaseHTTPRequestHandler):
   def log_message(self,*args):pass
   def do_POST(self):
    body=json.loads(self.rfile.read(int(self.headers['Content-Length'])));owner.requests.append(body)
    if len(owner.requests)==1:
     owner.entered.set();owner.release.wait();reply=owner.text('initial ordinary text without tools')
    elif owner.mode=='cli':reply=owner.text('CLI_HINT_HANDLED' if len(owner.requests)==2 else 'CLI_UNRELATED_REPLY')
    else:reply=owner.workflow(body)
    raw=json.dumps(reply).encode();self.send_response(200);self.send_header('Content-Type','application/json');self.end_headers();self.wfile.write(raw)
  self.server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);self.thread=threading.Thread(target=self.server.serve_forever);self.thread.start();self.url='http://127.0.0.1:'+str(self.server.server_port)
 def text(self,text):
  if self.mode=='cli':return {'content':[{'type':'text','text':text}],'usage':{'input_tokens':10,'output_tokens':10}}
  return {'choices':[{'message':{'role':'assistant','content':text}}],'usage':{'prompt_tokens':10,'completion_tokens':10}}
 def workflow(self,body):
  names=[x.get('function',{}).get('name') for x in body['tools']];role='author' if 'write_draft' in names else 'editor' if 'edit_draft' in names else 'reviewer';n=self.roles.get(role,0)+1;self.roles[role]=n
  message={'role':'assistant','content':'I will use the manuscript tools.' if n==1 else role+' final.'}
  if n==1:
   calls={'author':[('write_draft',{'text':AUTHOR}),('word_count',{})],'editor':[('read_draft',{}),('edit_draft',{'text':EDITED})],'reviewer':[('review',REVIEW)]}[role]
   message['tool_calls']=[{'id':role+'-'+name,'type':'function','function':{'name':name,'arguments':json.dumps(args)}} for name,args in calls]
  return {'choices':[{'message':message}],'usage':{'prompt_tokens':10,'completion_tokens':10}}
 def close(self):self.release.set();self.server.shutdown();self.thread.join();self.server.server_close()
class Process:
 def __init__(self,args,env,cwd,out):
  self.out=out;out.mkdir();self.started=time.monotonic();self.buffers={'stdout':'','stderr':''};self.raw={'stdout':bytearray(),'stderr':bytearray()};self.inputs=[];self.p=subprocess.Popen(args,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,cwd=cwd,env=env);self.readers=[];self.args=args
  for name,pipe in [('stdout',self.p.stdout),('stderr',self.p.stderr)]:
   def reader(name=name,pipe=pipe):
    decoder=codecs.getincrementaldecoder('utf-8')('replace')
    while True:
     part=pipe.read(1)
     if not part:
      self.buffers[name]+=decoder.decode(b'',final=True);break
     self.raw[name].extend(part);self.buffers[name]+=decoder.decode(part)
   t=threading.Thread(target=reader);t.start();self.readers.append(t)
 def send(self,text):self.inputs.append({'elapsed_seconds':time.monotonic()-self.started,'text':text});self.p.stdin.write((text+'\n').encode());self.p.stdin.flush()
 def until(self,fn,label):
  deadline=time.monotonic()+8
  while not fn():
   if self.p.poll() is not None:raise AssertionError('executable exited while waiting for '+label)
   if time.monotonic()>deadline:raise AssertionError('timeout waiting for '+label)
   time.sleep(.005)
 def finish(self):
  self.p.stdin.close();code=self.p.wait(timeout=20)
  for t in self.readers:t.join()
  self.save(code);assert code==0,('executable exit',code)
 def save(self,code):
  for k,v in self.raw.items():(self.out/(k+'.txt')).write_bytes(v)
  (self.out/'inputs.json').write_text(json.dumps(self.inputs,indent=2)+'\n');(self.out/'process.json').write_text(json.dumps({'argv':self.args,'exit':code,'elapsed_seconds':time.monotonic()-self.started,'fake_http_only':True},indent=2)+'\n')
 def abort(self):
  if self.p.poll() is None:self.p.terminate()
  try:code=self.p.wait(timeout=5)
  except subprocess.TimeoutExpired:self.p.kill();code=self.p.wait()
  for t in self.readers:t.join()
  self.save(code)
def rows(process):
 out=[]
 for line in process.buffers['stdout'].splitlines():
  try:out.append(json.loads(line))
  except json.JSONDecodeError:pass
 return out
results=[]
for mode,vendor,model in [('cli','anthropic','claude-sonnet-4-6'),('workflow','openai','gpt-4.1-2025-04-14')]:
 boundary=Boundary(mode);process=None;folder=OUT/mode
 try:
  with tempfile.TemporaryDirectory(prefix='cadence-ch05-r1-interface-') as temp:
   runtime=pathlib.Path(temp)
   env={k:v for k,v in os.environ.items() if k in ['PATH','HOME','TMPDIR','LANG','LC_ALL']}
   env.update({'LLM_VENDOR':vendor,'LLM_MODEL':model,'LLM_API_KEY':'local-correction-fixture','LLM_BASE_URL':boundary.url,'LLM_TOOL_BATCH_LIMIT':'1','CH06_LOG':str(runtime/'workflow.jsonl')})
   assert env['LLM_BASE_URL'].startswith('http://127.0.0.1:')
   binary=manifest['binaries'][mode if mode=='cli' else 'workflow']['path'];args=[binary]+(['chat'] if mode=='cli' else [])
   process=Process(args,env,runtime,folder)
   process.send('initial prompt' if mode=='cli' else json.dumps({'kind':'prompt','text':'Write and improve a short library manuscript.'}))
   assert boundary.entered.wait(5),'initial request not held at local HTTP'
   assert HINT not in json.dumps(boundary.requests[0])
   process.send('/hint '+HINT if mode=='cli' else json.dumps({'kind':'hint','agent':'author','text':HINT}))
   if mode=='cli':process.until(lambda:'Hint:' in process.buffers['stderr'] and HINT in process.buffers['stderr'],'actual CLI hint acknowledgement')
   else:process.until(lambda:any(r.get('observation')=='part_delta' and HINT in json.dumps(r) for r in rows(process)),'actual workflow hint acknowledgement')
   assert len(boundary.requests)==1,'hint mutated/duplicated in-flight request'
   acknowledged=time.monotonic()-process.started;boundary.release.set()
   if mode=='cli':
    process.until(lambda:'CLI_HINT_HANDLED' in process.buffers['stdout'],'same-turn CLI hint continuation')
    process.send('unrelated next prompt');process.until(lambda:'CLI_UNRELATED_REPLY' in process.buffers['stdout'],'independent later CLI reply')
   else:process.until(lambda:any(r.get('pipeline')=='reviewer_done' for r in rows(process)),'reliable workflow handoffs')
   process.finish()
   assert len(boundary.requests)==(3 if mode=='cli' else 7),('request count',len(boundary.requests))
   assert json.dumps(boundary.requests[1]).count(HINT)==1
   assert all(HINT not in json.dumps(body) for body in boundary.requests[2:]),'hint leaked after immediate continuation'
   proof={'mode':mode,'exit':0,'requests':len(boundary.requests),'hint_ack_before_first_response_release_seconds':acknowledged,'immutable_first_request':True,'same_turn_continuation':True,'hint_exactly_once':True,'no_later_leak':True,'configured_dispatched_batch_limit':1,'fake_http_only':True}
   if mode=='workflow':
    stream=rows(process);assert [r['pipeline'] for r in stream if 'pipeline' in r]==['author_done','editor_done','reviewer_done']
    tools={r['part']['name'] for r in stream if r.get('part',{}).get('type')=='tool_result' and not r['part'].get('is_error')}
    assert {'write_draft','word_count','read_draft','edit_draft','review'}<=tools
    assert (runtime/'manuscript.txt').read_text()==EDITED
    assert json.loads((runtime/'review.json').read_text())==REVIEW
    proof.update({'successful_role_tools':sorted(tools),'genuine_edited_manuscript':True,'genuine_nonempty_review':True,'single_dispatched_batch_per_role_preserved':True})
   for path in runtime.rglob('*'):
    if path.is_file() and path.suffix in ['.jsonl','.json','.txt']:
     dest=folder/'artifacts'/path.relative_to(runtime);dest.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(path,dest)
   (folder/'requests.json').write_text(json.dumps(boundary.requests,indent=2)+'\n');(folder/'result.json').write_text(json.dumps(proof,indent=2)+'\n');results.append(proof)
 except BaseException:
  boundary.release.set()
  if process:process.abort()
  raise
 finally:boundary.close()
(OUT/'summary.json').write_text(json.dumps({'passed':True,'source_commit':manifest['source_commit'],'binaries':manifest['binaries'],'actual_provider_calls':0,'results':results},indent=2)+'\n');print(json.dumps(results,indent=2))
