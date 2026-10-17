#!/usr/bin/env python3
"""Local HTTP exercise of the actual separate shape binary; no credentials."""
import http.server,json,os,pathlib,subprocess,sys,tempfile,threading,time
out=pathlib.Path(__file__).resolve().parent
binary=sys.argv[1]
results=[]
for plain in (False,True):
 state={'stream':None,'writes':[]}
 class Handler(http.server.BaseHTTPRequestHandler):
  def log_message(self,*args):pass
  def do_POST(self):
   req=json.loads(self.rfile.read(int(self.headers['Content-Length'])));state['stream']=bool(req.get('stream'))
   self.send_response(200);self.send_header('Content-Type','application/json' if plain else 'text/event-stream');self.end_headers()
   if plain:
    self.wfile.write(json.dumps({'model':'claude-sonnet-4-6','usage':{'input_tokens':0,'output_tokens':0},'content':[{'type':'text','text':'shape one two'}]}).encode());return
   frames=[{'type':'message_start','message':{'model':'claude-sonnet-4-6','usage':{'input_tokens':0,'output_tokens':0}}},{'type':'content_block_start','index':0,'content_block':{'type':'text','text':''}},{'type':'content_block_delta','index':0,'delta':{'type':'text_delta','text':'shape one'}},{'type':'content_block_delta','index':0,'delta':{'type':'text_delta','text':' two'}},{'type':'content_block_stop','index':0},{'type':'message_stop'}]
   for i,frame in enumerate(frames):
    self.wfile.write(('data: '+json.dumps(frame)+'\n\n').encode());self.wfile.flush();state['writes'].append(time.monotonic())
    if i==2:time.sleep(.15)
 server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);thread=threading.Thread(target=server.serve_forever);thread.start()
 work=tempfile.mkdtemp(prefix='e3-ch06-shape-fixture-');env=dict(os.environ,LLM_VENDOR='anthropic',LLM_MODEL='claude-sonnet-4-6',LLM_API_KEY='fixture-key',LLM_BASE_URL=f'http://127.0.0.1:{server.server_port}',CH07_NO_STREAM='1' if plain else '0');env.pop('EN_DISABLE_STREAMING',None)
 try:
  start=time.monotonic();p=subprocess.run([binary],input=b'fixture shape\n',stdout=subprocess.PIPE,stderr=subprocess.PIPE,env=env,cwd=work,timeout=10);elapsed=(time.monotonic()-start)*1000
 finally:server.shutdown();server.server_close();thread.join()
 label='plain' if plain else 'stream';(out/('shape-fixture-'+label+'.stdout.json')).write_bytes(p.stdout);(out/('shape-fixture-'+label+'.stderr.txt')).write_bytes(p.stderr)
 summary=json.loads(p.stdout);assert p.returncode==0 and state['stream']!=plain and summary['assistant'].strip()=='shape one two' and summary['parts']==1
 assert sum(summary['deltas'].values())>=1 and summary['first_delta_ms']>=0
 if not plain:assert sum(summary['deltas'].values())>=2 and summary['first_delta_ms']+100<elapsed
 results.append({'mode':label,'exit':p.returncode,'elapsed_ms':elapsed,'summary':summary,'private_fixture_dir':work,'binary':binary})
(out/'shape-fixture-summary.json').write_text(json.dumps(results,indent=2)+'\n');print(json.dumps(results,indent=2))
