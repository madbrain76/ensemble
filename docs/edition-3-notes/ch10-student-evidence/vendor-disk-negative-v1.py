import os,json,threading,subprocess,time,sys
from pathlib import Path
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
binary=Path(sys.argv[1]);work=Path(sys.argv[2]);expect=sys.argv[3];work.mkdir(exist_ok=False);requests=[]
class Handler(BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_POST(self):
  requests.append(json.loads(self.rfile.read(int(self.headers['Content-Length']))));b=json.dumps({'model':'claude-sonnet-4-6','content':[{'type':'text','text':'wire observation reply'}],'usage':{'input_tokens':2,'output_tokens':1}}).encode();self.send_response(200);self.send_header('Content-Length',str(len(b)));self.send_header('Content-Type','application/json');self.end_headers();self.wfile.write(b)
s=ThreadingHTTPServer(('127.0.0.1',0),Handler);t=threading.Thread(target=s.serve_forever);t.start()
e=dict(os.environ,ANTHROPIC_BASE_URL='http://127.0.0.1:'+str(s.server_port),LLM_API_KEY='offline-fixture-not-a-secret',LLM_VENDOR='anthropic',LLM_MODEL='claude-sonnet-4-6',LLM_SYSTEM='current runtime instruction',EN_DISABLE_STREAMING='1',CH02_LOG=str(work/'history.jsonl'))
def run(label):
 p=subprocess.run([str(binary)],input='{"user":"'+label+'"}\n',cwd=work,env=e,text=True,capture_output=True,timeout=15);(work/(label+'.stdout')).write_text(p.stdout);(work/(label+'.stderr')).write_text(p.stderr);return p
try:
 p=run('first');assert p.returncode==0,p.stderr;path=work/'save.json';a=json.loads(path.read_text());recorded=a['config']['vendor'];assert requests[0]['model']=='claude-sonnet-4-6' and 'current runtime instruction' in json.dumps(requests[0])
 a['config'].update(vendor='openai',model='historical-model',system_prompt='historical instruction');path.write_text(json.dumps(a));before=path.read_bytes();p=run('second')
 if expect=='negative':
  assert p.returncode!=0 and 'cannot unmarshal string' in p.stderr and len(requests)==1 and path.read_bytes()==before
 else:
  assert p.returncode==0,p.stderr;assert len(requests)==2 and 'first' in json.dumps(requests[1]) and requests[1]['model']=='claude-sonnet-4-6' and 'current runtime instruction' in json.dumps(requests[1]) and 'historical instruction' not in json.dumps(requests[1]);a=json.loads(path.read_text());assert a['config']['vendor']=='anthropic' and a['config']['model']=='claude-sonnet-4-6' and a['config']['system_prompt']=='current runtime instruction'
 print(json.dumps({'observed':expect,'original_saved_vendor':recorded,'archived_vendor_on_restart':'openai','restart_exit':p.returncode,'restart_stderr':p.stderr if expect=='negative' else '', 'fake_request_count':len(requests),'current_runtime_request_observed':True,'no_real_provider':True},indent=2))
finally:s.shutdown();s.server_close();t.join()
