import argparse,datetime,hashlib,importlib.util,json,os,socket,subprocess,sys,threading,time,traceback
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
from pathlib import Path
ROOT=Path('/Users/bill/projects/ensemble-edition-3-student');BASE=ROOT/'.ch08-local/operational-recovery';PREP=BASE/'preparation'
sha=lambda p:hashlib.sha256(Path(p).read_bytes()).hexdigest()
def save(p,v):Path(p).write_text(json.dumps(v,indent=2)+'\n')
a=argparse.ArgumentParser();a.add_argument('--sha256',required=True);a=a.parse_args();assert sha(__file__)==a.sha256
modpath=PREP/'controller-v1.py';assert sha(modpath)==(PREP/'controller-v1.sha256').read_text().strip();spec=importlib.util.spec_from_file_location('owned_controller',modpath);c=importlib.util.module_from_spec(spec);spec.loader.exec_module(c)
body=b'<!doctype html><meta charset="utf-8"><title>Chapter 8 owned smoke</title><main id="owned-marker">chapter8-owned-loopback-smoke-v1</main>'
http=[];lock=threading.Lock()
class Peer(BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_GET(self):
  row={'method':self.command,'path':self.path,'received_monotonic':time.monotonic(),'status':200 if self.path=='/' else 404,'response_sha256':hashlib.sha256(body if self.path=='/' else b'').hexdigest(),'bytes':len(body) if self.path=='/' else 0,'write_completed':False}
  with lock:http.append(row)
  self.send_response(row['status']);self.send_header('Content-Type','text/html; charset=utf-8');self.send_header('Content-Length',str(row['bytes']));self.end_headers()
  try:self.wfile.write(body if self.path=='/' else b'');self.wfile.flush();row['write_completed']=True
  except (BrokenPipeError,ConnectionResetError):row['write_error']='client disconnected'
server=ThreadingHTTPServer(('127.0.0.1',0),Peer);server.daemon_threads=True;t=threading.Thread(target=server.serve_forever,kwargs={'poll_interval':.05});t.start();url='http://127.0.0.1:'+str(server.server_port)+'/'
summary={'started_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'argv':sys.argv,'cwd':str(ROOT),'script_sha256':sha(__file__),'controller_sha256':sha(modpath),'http_listener':['127.0.0.1',server.server_port],'body_sha256':hashlib.sha256(body).hexdigest(),'model_requests':0,'key_reads':0,'audio_exercise':False,'cases':[]}
try:
 for label,flag in [('original-options',False),('plus-mock-keychain',True)]:
  meta=BASE/('smoke-'+label);meta.mkdir(exist_ok=False);b=None;start=time.monotonic();http_start=len(http);row={'label':label,'mock_keychain':flag,'status':'running','script_sha256':sha(__file__)};save(meta/'result.json',row)
  try:
   b=c.Browser(meta/'profile',True,meta,label,mock_keychain=flag)
   targets=b.call('Target.getTargets')['targetInfos'];target=next(x['targetId'] for x in targets if x['type']=='page');s=b.call('Target.attachToTarget',{'targetId':target,'flatten':True})['sessionId']
   for domain in ['Page','Runtime','Network']:b.call(domain+'.enable',session=s)
   n=time.monotonic();nav=b.call('Page.navigate',{'url':url},s,timeout=30);row['navigation_seconds']=round(time.monotonic()-n,3);row['navigation_result']=nav
   dom=b.wait(s,"(()=>{let n=document.querySelector('#owned-marker');return n?{marker:n.textContent,title:document.title,url:location.href,ready:document.readyState,bodyText:document.body.innerText}:false})()",seconds=5)
   assert dom['marker']=='chapter8-owned-loopback-smoke-v1';row['actual_DOM']=dom
   responses=[e for e in b.events if e.get('method')=='Network.responseReceived' and e['params']['response']['url']==url];assert responses
   response=b.call('Network.getResponseBody',{'requestId':responses[-1]['params']['requestId']},s);assert not response['base64Encoded'];got=response['body'].encode();assert got==body
   row['browser_response']={'bytes':len(got),'sha256':hashlib.sha256(got).hexdigest(),'status':responses[-1]['params']['response']['status']};row['status']='pass'
  except Exception as e:row.update(status='failed',failure={'type':type(e).__name__,'message':str(e),'traceback':traceback.format_exc()})
  finally:
   if b:b.close()
   else:
    # Constructor failures still have durable ownership; terminate only that verified process if needed.
    p=meta/(label+'-browser-ownership.json')
    if p.exists():
     own=json.loads(p.read_text());pid=own['pid']
     try:
      cmd=subprocess.check_output(['ps','-p',str(pid),'-o','command='],text=True);assert own['profile'] in cmd;os.kill(pid,15)
     except subprocess.CalledProcessError:pass
   row['http_observations']=http[http_start:];row['elapsed_seconds']=round(time.monotonic()-start,3);save(meta/'result.json',row);summary['cases'].append(row);save(BASE/'smoke-summary.json',summary);print('SMOKE',label,row['status'],row.get('failure',{}).get('message'),'HTTP',len(row['http_observations']),flush=True)
finally:
 server.shutdown();server.server_close();t.join(timeout=5);assert not t.is_alive();summary['http_server_stopped']=True;summary['finished_utc']=datetime.datetime.now(datetime.timezone.utc).isoformat();summary['http_all']=http;save(BASE/'smoke-summary.json',summary)
if not any(r['status']=='pass' for r in summary['cases']):sys.exit(1)
