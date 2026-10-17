#!/usr/bin/env python3
"""Actual frozen shape against local Gemini boundary fixture; no credentials/network provider."""
import http.server,json,os,pathlib,subprocess,sys,tempfile,threading,time
out=pathlib.Path(__file__).resolve().parent
root=pathlib.Path((out/'scratch-path.txt').read_text().strip())/'shape-fixture';root.mkdir(exist_ok=True);raw=root/'raw';raw.mkdir(exist_ok=True)
binary=sys.argv[1];results=[]
parts=[{'text':'shape one','thoughtSignature':'r1-synthetic-shape-text'},{'thoughtSignature':'r1-synthetic-shape-opaque'},{'text':' two'}]
def wire(items,done):
 c={'content':{'parts':items}}
 if done:c['finishReason']='STOP'
 return {'modelVersion':'gemini-3.8-flash','candidates':[c],'usageMetadata':{'promptTokenCount':0,'candidatesTokenCount':0}}
for plain in (False,True):
 state={'stream':None}
 class Handler(http.server.BaseHTTPRequestHandler):
  def log_message(self,*args):pass
  def do_POST(self):
   self.rfile.read(int(self.headers['Content-Length']));state['stream']='streamGenerateContent' in self.path
   self.send_response(200);self.send_header('Content-Type','application/json' if plain else 'text/event-stream');self.end_headers()
   if plain:self.wfile.write(json.dumps(wire(parts,True)).encode());return
   for i,p in enumerate(parts):
    self.wfile.write(('data: '+json.dumps(wire([p],i==len(parts)-1))+'\n\n').encode());self.wfile.flush()
    if i==0:time.sleep(.15)
 server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);thread=threading.Thread(target=server.serve_forever);thread.start()
 work=tempfile.mkdtemp(prefix='shape-',dir=raw)
 env={'PATH':os.defpath,'HOME':work,'TMPDIR':str(raw),'LLM_VENDOR':'gemini','LLM_MODEL':'gemini-3.8-flash','LLM_API_KEY':'r1-local-only','LLM_BASE_URL':f'http://127.0.0.1:{server.server_port}','CH07_NO_STREAM':'1' if plain else '0','EN_DISABLE_STREAMING':'0'}
 try:
  start=time.monotonic();p=subprocess.run([binary],input=b'fixture shape\n',stdout=subprocess.PIPE,stderr=subprocess.PIPE,env=env,cwd=work,timeout=10);elapsed=(time.monotonic()-start)*1000
 finally:server.shutdown();server.server_close();thread.join()
 label='plain' if plain else 'stream';(out/('shape-fixture-'+label+'.stdout.json')).write_bytes(p.stdout);(out/('shape-fixture-'+label+'.stderr.txt')).write_bytes(p.stderr)
 summary=json.loads(p.stdout)
 assert p.returncode==0 and state['stream']!=plain and summary['streaming']!=plain and summary['assistant']=='shape one two' and summary['parts']==3
 assert summary['deltas']=={'text':2,'thinking':0,'tool_call':0} and summary['first_delta_ms']>=0
 if not plain:assert summary['first_delta_ms']+100<elapsed
 results.append({'mode':label,'exit':p.returncode,'elapsed_ms':elapsed,'summary':summary,'private_fixture_dir':work,'binary':binary})
(out/'shape-fixture-summary.json').write_text(json.dumps(results,indent=2)+'\n');print(json.dumps(results,indent=2))
