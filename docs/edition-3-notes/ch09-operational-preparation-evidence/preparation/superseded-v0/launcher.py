#!/usr/bin/env python3
"""Chapter 9 owned frozen-artifact launcher and real counting proxy; no standalone paid entrypoint."""
import hashlib,http.client,http.server,json,os,pathlib,re,shutil,signal,socket,subprocess,threading,time,urllib.parse,urllib.request
ROOT=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
PREP=pathlib.Path(__file__).resolve().parent
POLICY=json.loads((PREP/'policy.json').read_text())
class Failure(Exception):pass
def require(ok,message):
 if not ok:raise Failure(message)
def sha(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def save(p,value):pathlib.Path(p).write_text(json.dumps(value,indent=2)+'\n')
def env(home):return {'PATH':os.environ.get('PATH','/usr/bin:/bin'),'HOME':str(home),'TMPDIR':os.environ.get('TMPDIR','/tmp'),'LANG':'en_US.UTF-8','PYTHONDONTWRITEBYTECODE':'1'}
def verify_frozen():
 require(subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=ROOT,text=True).strip()==POLICY['source_commit'],'source commit changed')
 require(subprocess.check_output(['git','--no-pager','rev-parse','HEAD^{tree}'],cwd=ROOT,text=True).strip()==POLICY['source_tree'],'source tree changed')
 require(not subprocess.check_output(['git','--no-pager','status','--porcelain'],cwd=ROOT,text=True).strip(),'source worktree changed')
 for r in POLICY['source_files']:require(sha(ROOT/r['path'])==r['sha256'],'tracked source hash changed: '+r['path'])
 for r in POLICY['binaries'].values():require(sha(r['path'])==r['sha256'],'frozen executable changed')
 for path,want in POLICY['fixtures'].items():require(sha(PREP/path)==want,'fixture changed')
def port_closed(port):
 with socket.socket() as s:s.settimeout(.3);return s.connect_ex(('127.0.0.1',port))!=0
def stop_process(p):
 if p is None:return None
 if p.poll() is None:
  os.killpg(p.pid,signal.SIGTERM)
  try:p.wait(timeout=15)
  except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);p.wait(timeout=10)
 return p.returncode
class Server(http.server.ThreadingHTTPServer):
 daemon_threads=True
 allow_reuse_address=False
class LocalPeer:
 """Scripted local upstream. Receives actual forwarded Agent requests, never synthetic proxy calls."""
 def __init__(self,directory):
  self.d=directory;self.requests=[];owner=self
  class Handler(http.server.BaseHTTPRequestHandler):
   def log_message(self,*args):pass
   def do_POST(self):
    request=json.loads(self.rfile.read(int(self.headers.get('Content-Length',0))))
    owner.requests.append({'path':self.path,'body':request});save(owner.d/'upstream-requests.json',owner.requests)
    n=len(owner.requests);calls=[('load_skill',{'name':'writer'}),('write_file',{'path':'first.txt','content':'CH09_FIRST'}),('unload_skill',{'name':'writer'}),('write_file',{'path':'after-unload.txt','content':'CH09_PENDING_STILL_ENABLED'})]
    blocks=[{'type':'text','text':'I will perform the next bounded skill action.'}]
    if n<=len(calls):
     name,args=calls[n-1];blocks.append({'type':'tool_use','id':'ch09-local-'+str(n),'name':name,'input':args})
    else:blocks=[{'type':'text','text':'CH09_DONE'}]
    if request.get('stream'):
     self.send_response(200);self.send_header('Content-Type','text/event-stream');self.end_headers()
     def event(v):self.wfile.write(('data: '+json.dumps(v)+'\n\n').encode());self.wfile.flush()
     event({'type':'message_start','message':{'id':'local','type':'message','role':'assistant','content':[],'model':request['model'],'usage':{'input_tokens':0,'output_tokens':0}}})
     for i,b in enumerate(blocks):
      start=dict(b);start['input' if b['type']=='tool_use' else 'text']={} if b['type']=='tool_use' else ''
      event({'type':'content_block_start','index':i,'content_block':start})
      delta={'type':'input_json_delta','partial_json':json.dumps(b['input'])} if b['type']=='tool_use' else {'type':'text_delta','text':b['text']}
      event({'type':'content_block_delta','index':i,'delta':delta});event({'type':'content_block_stop','index':i})
     event({'type':'message_delta','delta':{'stop_reason':'tool_use' if n<=4 else 'end_turn'},'usage':{'output_tokens':0}});event({'type':'message_stop'})
    else:
     body=json.dumps({'id':'local','type':'message','role':'assistant','model':request['model'],'content':blocks,'stop_reason':'tool_use' if n<=4 else 'end_turn','usage':{'input_tokens':0,'output_tokens':0}}).encode();self.send_response(200);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(body)));self.end_headers();self.wfile.write(body)
  self.server=Server(('127.0.0.1',0),Handler);self.port=self.server.server_port;require(self.port!=8088,'forbidden port')
  self.thread=threading.Thread(target=self.server.serve_forever,daemon=True);self.thread.start()
  save(self.d/'upstream-ownership.json',dict(pid=os.getpid(),listener=['127.0.0.1',self.port],mode='local scripted peer'))
 def close(self):self.server.shutdown();self.server.server_close();self.thread.join(3);require(port_closed(self.port),'upstream listener survived')
class Proxy:
 def __init__(self,directory,vendor,surface,upstream,credential=None):
  self.d=directory;self.vendor=vendor;self.surface=surface;self.upstream=urllib.parse.urlsplit(upstream);self.credential=credential;self.started=time.monotonic();self.records=[];self.lock=threading.Lock();self.connections=[];self.stopping=False;owner=self
  require((credential is None and self.upstream.scheme=='http' and self.upstream.hostname=='127.0.0.1') or (credential and upstream==POLICY['origins'][vendor]),'unauthorized upstream')
  class Handler(http.server.BaseHTTPRequestHandler):
   def log_message(self,*args):pass
   def do_POST(self):owner.forward(self)
  self.server=Server(('127.0.0.1',0),Handler);self.port=self.server.server_port;require(self.port!=8088,'forbidden port')
  self.thread=threading.Thread(target=self.server.serve_forever,daemon=True);self.thread.start()
  save(self.d/'proxy-ownership.json',dict(pid=os.getpid(),listener=['127.0.0.1',self.port],upstream=upstream,credential_in_artifacts=False))
 def route(self,path,body):
  v=self.vendor;model=POLICY['models'][v]
  if v=='anthropic':require(path=='/v1/messages' and body.get('model')==model,'Anthropic route/model');field='max_tokens';value=body.get(field)
  elif v=='openai':require(path=='/v1/responses' and body.get('model')==model,'OpenAI route/model');field='max_output_tokens';value=body.get(field)
  else:
   require(path in ('/v1beta/models/'+model+':generateContent','/v1beta/models/'+model+':streamGenerateContent?alt=sse'),'Gemini route/model');field='generationConfig.maxOutputTokens';value=body.get('generationConfig',{}).get('maxOutputTokens')
  require(type(value) is int and 0<value<=POLICY['output_caps'][self.surface],'output budget exceeded or absent')
  return field,value
 def forward(self,h):
  row=None;conn=None
  try:
   size=int(h.headers.get('Content-Length','0'));require(0<size<8*1024*1024,'request size')
   raw=h.rfile.read(size);body=json.loads(raw);field,value=self.route(h.path,body)
   with self.lock:
    require(not self.stopping and time.monotonic()-self.started<POLICY['session_seconds'],'session closed/deadline')
    require(len(self.records)<POLICY['max_forwarded_requests_per_session'],'request budget exhausted')
    row=dict(number=len(self.records)+1,path=h.path,method='POST',body=body,body_sha256=hashlib.sha256(raw).hexdigest(),max_output_field=field,max_output_value=value,elapsed_seconds=round(time.monotonic()-self.started,3),forward_attempted=True,complete=False);self.records.append(row);save(self.d/'proxy-requests.json',self.records)
   headers={'Content-Type':'application/json','Accept':h.headers.get('Accept','application/json')}
   if self.vendor=='anthropic':headers.update({'anthropic-version':'2023-06-01','x-api-key':self.credential or 'local-placeholder'})
   elif self.vendor=='openai':headers['Authorization']='Bearer '+(self.credential or 'local-placeholder')
   else:headers['x-goog-api-key']=self.credential or 'local-placeholder'
   cls=http.client.HTTPSConnection if self.upstream.scheme=='https' else http.client.HTTPConnection
   conn=cls(self.upstream.hostname,self.upstream.port,timeout=min(65,max(1,POLICY['session_seconds']-(time.monotonic()-self.started))))
   with self.lock:self.connections.append(conn)
   conn.request('POST',h.path,raw,headers);response=conn.getresponse();row['upstream_status']=response.status
   h.send_response(response.status);ctype=response.getheader('Content-Type','application/json');h.send_header('Content-Type',ctype);h.send_header('Connection','close');h.end_headers();digest=hashlib.sha256();total=0
   while True:
    chunk=response.readline() if 'event-stream' in ctype else response.read(65536)
    if not chunk:break
    digest.update(chunk);total+=len(chunk);h.wfile.write(chunk);h.wfile.flush()
   row.update(complete=True,response_bytes=total,response_sha256=digest.hexdigest());h.close_connection=True
  except Exception as e:
   if row is not None:row['failure_type']=type(e).__name__
   try:h.send_error(502,'bounded forwarding failed')
   except (OSError,ValueError):pass
   with (self.d/'proxy-errors.jsonl').open('a') as f:f.write(json.dumps({'type':type(e).__name__,'message':str(e) if isinstance(e,Failure) else 'upstream/transport failure; no exception secrets recorded'})+'\n')
  finally:
   if conn:
    conn.close()
    with self.lock:
     if conn in self.connections:self.connections.remove(conn)
   with self.lock:save(self.d/'proxy-requests.json',self.records)
 def close(self):
  with self.lock:
   self.stopping=True
   for c in self.connections:
    if c.sock:
     try:c.sock.shutdown(socket.SHUT_RDWR)
     except OSError:pass
    c.close()
  self.server.shutdown();self.server.server_close();self.thread.join(5)
  end=time.monotonic()+5
  while self.connections and time.monotonic()<end:time.sleep(.05)
  require(not self.connections and not self.thread.is_alive() and port_closed(self.port),'proxy cleanup incomplete')
  save(self.d/'proxy-close.json',dict(requests=len(self.records),listener_closed=True,active_upstreams=0,elapsed_seconds=round(time.monotonic()-self.started,3),max_requests=16,retries=0))
class Artifact:
 def __init__(self,directory,surface,vendor,proxy,stream=False):
  self.d=directory;self.surface=surface;self.p=None;self.port=None;self.stdout=None;self.stderr=None
  for name in ('home','workspace','cr'): (directory/name).mkdir()
  shutil.copytree(PREP/'skills',directory/'skills')
  save(directory/'cr/settings.json',dict(theme='dark',tts_enabled=False,tts_speed=1,font_size=16,max_tool_rounds=8))
  binary=POLICY['binaries'][surface];e=env(directory/'home');e.update(LLM_VENDOR=vendor,LLM_MODEL=POLICY['models'][vendor],LLM_BASE_URL='http://127.0.0.1:'+str(proxy.port),LLM_API_KEY='local-proxy-placeholder',EN_SKILLS_DIR=str(directory/'skills'),EN_PRIMARY_SKILL='primary',EN_CUSTOM_VAR='CH09_CALLBACK_VALUE')
  args=[binary['path']]
  if surface=='gui':
   args+=['-workspace',str(directory/'workspace'),'-event-log',str(directory/'event.log'),'-port','0']
   if not stream:args+=['-no-stream']
  else:
   e.update(ENSEMBLE_WORKSPACE=str(directory/'workspace'),ENSEMBLE_DATA_DIR=str(directory/'data'),CH02_LOG=str(directory/'event.log'),LLM_TOOL_BATCH_LIMIT='8')
   if not stream:e['EN_DISABLE_STREAMING']='1'
  self.stdout=(directory/'artifact.stdout.txt').open('w');self.stderr=(directory/'artifact.stderr.txt').open('w')
  self.p=subprocess.Popen(args,cwd=directory,env=e,stdin=subprocess.PIPE,stdout=self.stdout,stderr=self.stderr,text=True,start_new_session=True)
  save(directory/'artifact-ownership.json',dict(pid=self.p.pid,argv=args,cwd=str(directory),environment=e,binary_sha256=sha(binary['path']),source_commit=POLICY['source_commit'],surface=surface,started_unix=time.time()))
  if surface=='gui':
   end=time.monotonic()+20
   while time.monotonic()<end:
    text=(directory/'artifact.stderr.txt').read_text();m=re.search(r'GUI http://127\.0\.0\.1:(\d+)',text)
    if m:self.port=int(m[1]);break
    require(self.p.poll() is None,'GUI exited during startup');time.sleep(.05)
   require(self.port is not None and self.port!=8088,'GUI listener missing/forbidden')
   proof=subprocess.check_output(['/usr/sbin/lsof','-nP','-a','-p',str(self.p.pid),'-iTCP:'+str(self.port),'-sTCP:LISTEN'],text=True)
   require('127.0.0.1:'+str(self.port) in proof and 'LISTEN' in proof,'GUI listener not owned')
   save(directory/'gui-listener.json',dict(pid=self.p.pid,listener=['127.0.0.1',self.port],proof=proof,verified_before_client_access=True))
   assets=[]
   for name,want in POLICY['assets'].items():
    with urllib.request.urlopen(self.url+name,timeout=5) as r:raw=r.read()
    got=hashlib.sha256(raw).hexdigest();require(got==want,'served asset differs: '+name);assets.append(dict(path=name,sha256=got))
   save(directory/'served-assets.json',assets)
 @property
 def url(self):return 'http://127.0.0.1:'+str(self.port)+'/'
 def close(self):
  code=stop_process(self.p)
  if self.p and self.p.stdin and not self.p.stdin.closed:self.p.stdin.close()
  for f in (self.stdout,self.stderr):
   if f:f.close()
  closed=self.port is None or port_closed(self.port)
  save(self.d/'artifact-close.json',dict(pid=self.p.pid if self.p else None,exit_code=code,process_exited=self.p is None or self.p.poll() is not None,listener_closed=closed))
  require(closed,'GUI listener survived')
def credential_for(vendor):
 # Only controller.execute calls this AFTER approval, freeze, disk, deadline and run-allocation checks.
 name=POLICY['credential_names'][vendor];path=pathlib.Path.home()/'all_api_keys.txt';raw=path.read_text();matches=[]
 for line in raw.splitlines():
  m=re.fullmatch(r'(?:export\s+)?'+re.escape(name)+r'\s*=\s*(.*)',line.strip())
  if m:matches.append(m.group(1).strip().strip('\"\''))
 require(len(matches)==1 and matches[0] and '\n' not in matches[0],'required credential missing/ambiguous')
 return matches[0]
