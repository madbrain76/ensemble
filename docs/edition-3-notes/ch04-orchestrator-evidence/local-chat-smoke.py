#!/usr/bin/env python3
"""Local synthetic human-CLI smoke; no provider access or real credentials."""
import hashlib,json,os,pathlib,re,subprocess,tempfile,threading,time
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
OUT=pathlib.Path(__file__).resolve().parent
ROOT=pathlib.Path(tempfile.mkdtemp(prefix='cadence-ch04-chat-'))
requests=[]; errors=[]; observations=[]
class Handler(BaseHTTPRequestHandler):
 def log_message(self,*args): pass
 def do_POST(self):
  try:
   body=json.loads(self.rfile.read(int(self.headers['Content-Length'])))
   idx=len(requests); requests.append(body); messages=body['messages']
   msg={'role':'assistant','content':''}
   if idx==0:
    msg['content']='I will run a short local command so I can observe its output before and after completion.'
    msg['tool_calls']=[{'id':'run-one','type':'function','function':{'name':'run_command','arguments':json.dumps({'command':"printf FIRST; sleep 0.3; printf SECOND",'ai_callback_delay':0.03})}}]
   elif idx==1:
    text=next(x['content'] for x in reversed(messages) if x['role']=='tool')
    assert 'running' in text and 'FIRST' in text and 'SECOND' not in text,text
    handle=int(re.search(r'job (\d+) running',text).group(1));observations.append({'stage':'running','return':text,'handle':handle})
    msg['content']='I have the running handle. I will wait for its completion and inspect only the new output.'
    msg['tool_calls']=[{'id':'wait-one','type':'function','function':{'name':'wait_for_job','arguments':json.dumps({'handle':handle,'ai_callback_delay':1})}}]
   elif idx==2:
    text=next(x['content'] for x in reversed(messages) if x['role']=='tool')
    assert 'SECOND' in text and 'FIRST' not in text and 'exit_code: 0' in text,text
    observations.append({'stage':'completed','return':text})
    msg['content']='JOB-PASS: observed the command running, then its unseen completion output.'
   elif idx==3:
    assert any(x.get('role')=='assistant' and 'JOB-PASS' in str(x.get('content','')) for x in messages)
    assert messages[-1]['role']=='user' and 'follow-up' in messages[-1]['content']
    msg['content']='CONTEXT-PASS: the earlier job observation remains in this conversation.'
   else: raise AssertionError('unexpected extra request')
   result={'id':'local-smoke','model':'fixture-ch04','object':'chat.completion','choices':[{'index':0,'message':msg,'finish_reason':'tool_calls' if 'tool_calls' in msg else 'stop'}],'usage':{'prompt_tokens':10,'completion_tokens':5,'total_tokens':15}}
   data=json.dumps(result).encode();self.send_response(200);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(data)));self.end_headers();self.wfile.write(data)
  except Exception as e:
   errors.append(repr(e));self.send_error(500)
server=ThreadingHTTPServer(('127.0.0.1',0),Handler);thread=threading.Thread(target=server.serve_forever,daemon=True);thread.start()
env={k:os.environ[k] for k in ('PATH','HOME','TMPDIR') if k in os.environ}
env.update(LLM_VENDOR='openai',LLM_BASE_URL=f'http://127.0.0.1:{server.server_port}',LLM_MODEL='fixture-ch04',LLM_API_KEY='local-fake-not-a-credential',LLM_REQUIRE_TOOL_EXPLANATION='true',ENSEMBLE_WORKSPACE=str(ROOT),ENSEMBLE_DATA_DIR=str(ROOT/'data'),CH02_LOG=str(ROOT/'history.jsonl'))
start=time.monotonic()
try:
 run=subprocess.run(['/tmp/ch04-core','chat'],input='Observe a short local job.\nThis is the follow-up; retain prior observations.\n',text=True,capture_output=True,cwd=ROOT,env=env,timeout=30)
 (OUT/'local-chat.stdout.txt').write_text(run.stdout);(OUT/'local-chat.stderr.txt').write_text(run.stderr)
 record={'kind':'synthetic local provider; actual human chat CLI','binary_sha256':hashlib.sha256(pathlib.Path('/tmp/ch04-core').read_bytes()).hexdigest(),'workspace':str(ROOT),'exit_code':run.returncode,'elapsed_seconds':time.monotonic()-start,'request_count':len(requests),'errors':errors,'observations':observations,'live_provider_calls':0}
 record['passed']=run.returncode==0 and not errors and len(requests)==4 and 'JOB-PASS' in run.stdout and 'CONTEXT-PASS' in run.stdout
 (OUT/'local-chat-receipt.json').write_text(json.dumps(record,indent=2)+'\n');print(json.dumps(record,indent=2));assert record['passed']
finally:
 server.shutdown();server.server_close();thread.join()
