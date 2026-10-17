#!/usr/bin/env python3
"""Exact five-case Chapter 6 launcher. Live mode requires separate parent approval.
Rehearsal uses actual frozen binaries and this same forwarding boundary, no keys.
"""
import concurrent.futures,hashlib,http.client,http.server,json,os,pathlib,selectors,subprocess,tempfile,threading,time,sys
OUT=pathlib.Path(__file__).resolve().parent
SOURCE=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
MANIFEST=OUT.parent/'ch06-student-source-manifest.json'
MANIFEST_HASH='e1a5fca437e3002f8afc0c42fed6c8aac3cae22251a11b3a78cf8e8364188028'
COMMIT='aa704ce00abe4203845e6dff86e4c96105fc7cef';TREE='c3541c685ce35b1491e311a98e719ddf2de92001'
CORE='/tmp/e3-ch06-frozen-core';SHAPE='/tmp/e3-ch06-frozen-shape'
PINS={CORE:'f86f4c415148ed6980e764f1ff8bac60d79b962c14ee3cddeda2012d58fff3e7',SHAPE:'6d010fde5927e7164fa60c3370f50cce9e5f9013f4fc608a0c5de1f21fb30605','/tmp/e3-ch06-candidate.tar':'518c09084aaba668cf7de3d29dd55b1fd976f6b1060e67a47d5498880dd0da7a'}
ROWS={'anthropic':('claude-sonnet-4-6','api.anthropic.com','directClaudeAPIKey'),'openai':('gpt-4.1-2025-04-14','api.openai.com','directOpenAIAPIKey'),'gemini':('gemini-3.8-flash','generativelanguage.googleapis.com','directGeminiAPIKey')}
CASES=[('anthropic','human-stream'),('openai','human-stream'),('gemini','human-stream'),('openai','shape-stream'),('openai','shape-plain')]
SAFE_ENV={'PATH':'/usr/bin:/bin:/usr/sbin:/sbin','LANG':'C.UTF-8'}
def digest(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def verify(pins=PINS):
 assert digest(MANIFEST)==MANIFEST_HASH,'manifest pin mismatch'
 for p,want in pins.items():assert digest(p)==want,'artifact pin mismatch'
 def git(*a):return subprocess.check_output(['git','-C',str(SOURCE),*a],env=SAFE_ENV,text=True).strip()
 assert git('rev-parse','HEAD')==COMMIT and git('rev-parse','HEAD^{tree}')==TREE,'source pin mismatch'
 assert git('status','--porcelain')=='','dirty candidate'
def write(name,value): (OUT/name).write_text(json.dumps(value,indent=2)+'\n')
class Server:
 def __init__(self,handler):
  self.server=http.server.ThreadingHTTPServer(('127.0.0.1',0),handler);self.server.daemon_threads=False
  self.port=self.server.server_port;self.thread=threading.Thread(target=self.server.serve_forever);self.thread.start()
 def close(self):self.server.shutdown();self.server.server_close();self.thread.join()
class Forward:
 def __init__(self,vendor,plain,key,mark,upstream=None):
  self.vendor,self.plain,self.key,self.mark=vendor,plain,key,mark;self.model,self.host,_=ROWS[vendor]
  self.upstream=upstream;self.lock=threading.Lock();self.attempts=[];self.refusals=[];self.start=time.monotonic();self.active=False
  owner=self
  class Handler(http.server.BaseHTTPRequestHandler):
   def log_message(self,*args):pass
   def do_POST(self):owner.handle(self)
  self.server=Server(Handler)
 def ms(self):return round((time.monotonic()-self.start)*1000,3)
 def path(self):
  if self.vendor=='anthropic':return '/v1/messages'
  if self.vendor=='openai':return '/v1/chat/completions'
  return '/v1beta/models/'+self.model+(':generateContent' if self.plain else ':streamGenerateContent?alt=sse')
 def handle(self,h):
  try:
   data=h.rfile.read(int(h.headers.get('Content-Length','0')));body=json.loads(data)
   verify() # Frozen source/manifest/binary pins checked BEFORE reservation/upstream.
   tokens=body.get('generationConfig',{}).get('maxOutputTokens') if self.vendor=='gemini' else body.get('max_tokens',body.get('max_completion_tokens'))
   valid=h.path==self.path() and isinstance(tokens,int) and not isinstance(tokens,bool) and 1<=tokens<=1024
   if self.vendor!='gemini':valid=valid and body.get('model')==self.model and bool(body.get('stream',False))!=self.plain
   with self.lock:
    if not valid or len(self.attempts)>=3:
     self.refusals.append({'at_ms':self.ms(),'reason':'invalid-selection' if not valid else 'request-cap'});h.send_response(429);h.end_headers();return
    attempt={'number':len(self.attempts)+1,'start_ms':self.ms(),'marker_in_request':self.mark in data.decode(),'tokens':tokens,'stream':not self.plain,'status':None};self.attempts.append(attempt);self.active=True
   headers={'Content-Type':'application/json'}
   if self.vendor=='anthropic':headers.update({'x-api-key':self.key,'anthropic-version':'2023-06-01'})
   elif self.vendor=='openai':headers['Authorization']='Bearer '+self.key
   else:headers['x-goog-api-key']=self.key
   conn=http.client.HTTPConnection('127.0.0.1',self.upstream,timeout=30) if self.upstream else http.client.HTTPSConnection(self.host,timeout=60)
   conn.request('POST',h.path,data,headers);response=conn.getresponse() # No redirects/retries.
   attempt['status']=response.status;attempt['headers_ms']=self.ms();h.send_response(response.status);h.send_header('Content-Type',response.getheader('Content-Type','application/json'));h.end_headers()
   while True:
    line=response.readline()
    if not line:break
    h.wfile.write(line.replace(self.key.encode(),b'[redacted]'));h.wfile.flush()
   attempt['end_ms']=self.ms();conn.close()
  except Exception as error:
   with self.lock:self.refusals.append({'at_ms':self.ms(),'reason':type(error).__name__})
   try:h.send_response(502);h.end_headers()
   except Exception:pass
  finally:
   with self.lock:self.active=False
 def close(self):self.server.close()

def frames(vendor,model,parts,plain):
 usage={'input_tokens':1,'output_tokens':1}
 if vendor=='anthropic':
  content=[{'type':'text','text':p['text']} if 'text' in p else {'type':'tool_use','id':p['id'],'name':p['name'],'input':p['args']} for p in parts]
  result={'model':model,'usage':usage,'content':content}
  events=[{'type':'message_start','message':{'model':model,'usage':usage}}]
  for i,p in enumerate(content):events += [{'type':'content_block_start','index':i,'content_block':p},{'type':'content_block_stop','index':i}]
  events += [{'type':'message_stop'}]
 elif vendor=='openai':
  text=''.join(p.get('text','') for p in parts);tools=[{'id':p['id'],'type':'function','function':{'name':p['name'],'arguments':json.dumps(p['args'])}} for p in parts if 'name' in p]
  msg={'role':'assistant','content':text,'tool_calls':tools};u={'prompt_tokens':1,'completion_tokens':1};result={'model':model,'choices':[{'message':msg}],'usage':u}
  events=[{'model':model,'choices':[{'delta':{'content':p['text']}}]} if 'text' in p else {'choices':[{'delta':{'tool_calls':[dict(index=sum('name' in x for x in parts[:i]),**tools[sum('name' in x for x in parts[:i])])]}}]} for i,p in enumerate(parts)]
  events += [{'choices':[{'delta':{},'finish_reason':'tool_calls' if tools else 'stop'}]},{'choices':[],'usage':u}]
 else:
  content=[{'text':p['text']} if 'text' in p else {'functionCall':{'name':p['name'],'args':p['args']}} for p in parts];u={'promptTokenCount':1,'candidatesTokenCount':1}
  result={'modelVersion':model,'candidates':[{'content':{'parts':content},'finishReason':'STOP'}],'usageMetadata':u}
  events=[{'modelVersion':model,'candidates':[{'content':{'parts':[p]}}]} for p in content]+[{'candidates':[{'finishReason':'STOP'}],'usageMetadata':u}]
 if plain:return [json.dumps(result).encode()]
 data=[('data: '+json.dumps(e)+'\n\n').encode() for e in events]
 if vendor=='openai':data.append(b'data: [DONE]\n\n')
 return data
class Fixture:
 def __init__(self,vendor,mode,mark,fail=False):
  self.count=0;self.release=threading.Event();owner=self;model=ROWS[vendor][0]
  class Handler(http.server.BaseHTTPRequestHandler):
   def log_message(self,*args):pass
   def do_POST(self):
    body=self.rfile.read(int(self.headers.get('Content-Length','0')));owner.count+=1;n=owner.count
    if fail:self.send_response(503);self.end_headers();self.wfile.write(b'{"error":"fixture unavailable"}');return
    human=mode=='human-stream';plain=mode=='shape-plain'
    if n==1:parts=[{'text':'I will inspect config.json now. '},{'id':'read1','name':'read_file','args':{'path':'config.json'}}]
    elif human and n==2:parts=[{'text':'I will write the requested file. '},{'id':'write1','name':'write_file','args':{'path':'output.txt','content':'STREAM_OK'}}]
    else:parts=[{'text':(('Done '+mark if mark in body.decode() else 'Done without hint') if human else 'stream-study setting 7') }]
    self.send_response(200);self.send_header('Content-Type','application/json' if plain else 'text/event-stream');self.end_headers()
    for index,data in enumerate(frames(vendor,model,parts,plain)):
     self.wfile.write(data);self.wfile.flush()
     if n==1 and human and index==(1 if vendor=='anthropic' else 0):owner.release.wait(5)
     if not plain:time.sleep(.015)
  self.server=Server(Handler)
 def close(self):self.release.set();self.server.close()

def case(vendor,mode,key,root,prefix,local=False,fail=False):
 verify();label=vendor+'-'+mode;mark='HINT_APPLIED_'+label.replace('-','_');hint='Include marker '+mark+' in the final answer.'
 work=root/label;work.mkdir();(work/'config.json').write_text('{"project":"stream-study","setting":7}\n')
 fixture=Fixture(vendor,mode,mark,fail) if local else None
 forward=Forward(vendor,mode=='shape-plain',key,mark,fixture.server.port if fixture else None)
 env=dict(SAFE_ENV,HOME=str(work),TMPDIR=str(work),LLM_VENDOR=vendor,LLM_MODEL=ROWS[vendor][0],LLM_BASE_URL=f'http://127.0.0.1:{forward.server.port}',LLM_API_KEY='local-harness-placeholder',LLM_TOOL_BATCH_LIMIT='2',LLM_REQUIRE_TOOL_EXPLANATION='true',LLM_MAX_TOKENS='1024',CH07_NO_STREAM='1' if mode=='shape-plain' else '0')
 binary=CORE if mode=='human-stream' else SHAPE;verify()
 proc=subprocess.Popen([binary]+(['chat'] if mode=='human-stream' else []),cwd=work,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
 prompt='Read config.json using read_file, then write output.txt containing exactly STREAM_OK using write_file. Explain briefly before each tool batch. Never run shell commands. End with a short sentence.' if mode=='human-stream' else 'Use read_file to inspect config.json, then report its project name and numeric setting in one short sentence. Explain before tools. Never write files or run shell commands.'
 proc.stdin.write((prompt+'\n').encode());proc.stdin.flush();closed=mode!='human-stream'
 if closed:proc.stdin.close()
 sel=selectors.DefaultSelector();sel.register(proc.stdout,selectors.EVENT_READ,'stdout');sel.register(proc.stderr,selectors.EVENT_READ,'stderr');raw={'stdout':bytearray(),'stderr':bytearray()}
 result={'case':label,'pin_check':True,'first_content_ms':None,'hint_sent_ms':None,'hint_ack_ms':None,'hint_sent_while_forward_active':False,'chunks':[]}
 try:
  while sel.get_map():
   if forward.ms()>90000:result['timed_out']=True;proc.kill()
   if not closed and forward.ms()>12000:proc.stdin.close();closed=True
   for item,_ in sel.select(.1):
    data=os.read(item.fileobj.fileno(),4096)
    if not data:sel.unregister(item.fileobj);continue
    raw[item.data].extend(data)
    if item.data=='stdout':
     content=raw['stdout'].decode(errors='replace').removeprefix('> ')
     result['chunks'].append({'ms':forward.ms(),'bytes':len(data)})
     if content.strip() and result['first_content_ms'] is None:result['first_content_ms']=forward.ms()
     if mode=='human-stream' and content.strip() and not closed:
      result['hint_sent_while_forward_active']=forward.active;result['hint_sent_ms']=forward.ms()
      proc.stdin.write(('/hint '+hint+'\n').encode());proc.stdin.flush();proc.stdin.close();closed=True
    if hint in raw['stderr'].decode(errors='replace') and result['hint_ack_ms'] is None:
     result['hint_ack_ms']=forward.ms()
     if fixture:fixture.release.set()
  result['exit']=proc.wait(timeout=5)
 except Exception as error:result['launcher_error']=type(error).__name__
 finally:
  if proc.poll() is None:proc.kill();proc.wait()
  if not closed:proc.stdin.close()
  if fixture:fixture.release.set()
  sel.close();forward.close()
  if fixture:fixture.close()
 for stream,data in raw.items():(work/(stream+'.raw')).write_bytes(data)
 result.update(requests=forward.attempts,refusals=forward.refusals,private_directory=str(work),elapsed_ms=forward.ms(),binary_sha256=digest(binary))
 stdout=raw['stdout'].decode(errors='replace').replace(key,'[redacted]');stderr=raw['stderr'].decode(errors='replace').replace(key,'[redacted]')
 (OUT/(prefix+'-'+label+'.stdout.txt')).write_text(stdout)
 (OUT/(prefix+'-'+label+'.stderr-safe.txt')).write_text('\n'.join(s for s in stderr.splitlines() if s.startswith('Tokens used:'))+'\n[Other stderr retained privately.]\n')
 checks={'exit_zero':result.get('exit')==0,'provider_success':bool(forward.attempts) and all(a['status']==200 for a in forward.attempts),'within_caps':len(forward.attempts)<=3 and not forward.refusals,'genuine_content':bool(stdout.removeprefix('> ').strip())}
 if mode=='human-stream':
  effect=work/'output.txt';checks['file_effect']=effect.exists() and effect.read_text()=='STREAM_OK';checks['hint_active']=bool(result['hint_sent_while_forward_active'] and result['hint_ack_ms'] and result['hint_ack_ms']<forward.attempts[0].get('end_ms',0));checks['hint_in_next_request']=any(a['marker_in_request'] for a in forward.attempts[1:]);checks['hint_in_answer']=mark in stdout;checks['early_stdout']=bool(result['first_content_ms'] and result['first_content_ms']<forward.attempts[0].get('end_ms',0))
 else:
  try:
   shape=json.loads(stdout);result['shape']=shape;checks['shape']=shape['observation']=='stream_stats' and shape['parts']>0 and sum(shape['deltas'].values())>0 and 0<=shape['first_delta_ms']<shape['elapsed_ms'] and shape['streaming']==(mode=='shape-stream') and 'stream-study' in shape['assistant'] and '7' in shape['assistant']
   if mode=='shape-stream':checks['multiple_deltas']=sum(shape['deltas'].values())>=2
  except (KeyError,ValueError,TypeError):checks['shape']=False
 result['checks']=checks;result['passed']=all(checks.values());write(prefix+'-'+label+'.json',result);return result

def series(keys,prefix,local=False,fail=False):
 root=pathlib.Path(tempfile.mkdtemp(prefix='e3-ch06-'+prefix+'-private-'));root.chmod(0o700);results=[]
 for vendor,mode in CASES:
  r=case(vendor,mode,keys[vendor],root,prefix,local,fail);results.append(r);write(prefix+'-summary.json',results)
  print(prefix,r['case'],r['passed'],r['checks'],flush=True)
  if not r['passed']:break
 return results

def negatives():
 bad=dict(PINS);bad[CORE]='0'*64
 try:verify(bad);raise RuntimeError('wrong pin accepted')
 except AssertionError:pass
 fixture=Fixture('openai','shape-stream','marker');f=Forward('openai',False,'fixture-only-key','marker',fixture.server.port)
 body={'model':ROWS['openai'][0],'stream':True,'max_tokens':1024}
 def post(payload,path=None):
  c=http.client.HTTPConnection('127.0.0.1',f.server.port,timeout=10);c.request('POST',path or f.path(),json.dumps(payload),{'Content-Type':'application/json'});r=c.getresponse();code=r.status;r.read();c.close();return code
 try:
  wrong_model=post(dict(body,model='not-approved'));wrong_tokens=post(dict(body,max_tokens=1025));wrong_path=post(body,'/wrong');assert wrong_model==wrong_tokens==wrong_path==429 and fixture.count==0
  with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:statuses=list(pool.map(lambda _:post(body),range(8)))
  assert statuses.count(200)==3 and statuses.count(429)==5 and fixture.count==3 and len(f.attempts)==3
 finally:f.close();fixture.close()
 r=series({v:'fixture-only-key' for v in ROWS},'rehearsal2-first-failure',True,True);assert len(r)==1 and not r[0]['passed']
 write('rehearsal2-negatives.json',{'wrong_pin_blocks_before_loader':True,'wrong_model':wrong_model,'wrong_tokens':wrong_tokens,'wrong_path':wrong_path,'concurrent_forward_statuses':statuses,'upstream_count':fixture.count,'stopped_after_first_failure':True,'settings_reads':0})
def main():
 if sys.argv[1:] not in (['--rehearse'],['--approved-live']):raise SystemExit('Exactly --rehearse or separately approved --approved-live is required.')
 verify()
 if sys.argv[1]=='--rehearse':
  negatives();r=series({v:'fixture-only-key' for v in ROWS},'rehearsal2',True);assert len(r)==5 and all(x['passed'] for x in r)
 else:
  # Reached only after frozen pins and clean-source preflight, and parent approval.
  settings=json.loads((pathlib.Path.home()/'.cr/settings.json').read_text());keys={v:settings.get(ROWS[v][2],'') for v in ROWS};del settings
  if not all(keys.values()):raise SystemExit('Approved credential field absent; stopped.')
  r=series(keys,'live-v2');raise SystemExit(0 if len(r)==5 and all(x['passed'] for x in r) else 1)
if __name__=='__main__':main()
