#!/usr/bin/env python3
"""Compact real fail-closed guard rehearsal. Loopback only; no browser or credential access."""
import http.client,json,os,pathlib,subprocess,sys,time
from launcher import PREP,POLICY,LocalPeer,Proxy,require,save,sha
D=PREP/'guard-rehearsal';D.mkdir(exist_ok=False)
env={k:os.environ[k] for k in ('PATH','HOME','TMPDIR') if k in os.environ};env['PYTHONDONTWRITEBYTECODE']='1'
audit="""import sys,runpy,json,pathlib
out=pathlib.Path(sys.argv[1]);sys.argv=sys.argv[2:];events=[]
def hook(event,args):
 if event=='subprocess.Popen':events.append({'type':event,'executable':str(args[0])})
 if event=='socket.connect':events.append({'type':event,'address':str(args[1])})
 if event=='open' and str(args[0]).endswith('/.cr/settings.json'):events.append({'type':'credential_open'})
sys.addaudithook(hook)
try:runpy.run_path(sys.argv[0],run_name='__main__')
finally:out.write_text(json.dumps(events,indent=2)+'\\n')
"""
def rejected(name,args):
 cmd=[sys.executable,'-B','-c',audit,str(D/(name+'-audit.json')),str(PREP/'controller.py')]+args;t=time.monotonic();start=time.time()
 with (D/(name+'.stdout.txt')).open('w') as out,(D/(name+'.stderr.txt')).open('w') as err:r=subprocess.run(cmd,cwd=PREP.parent.parent,env=env,stdout=out,stderr=err)
 save(D/(name+'-receipt.json'),dict(argv=cmd,cwd=str(PREP.parent.parent),environment=env,started_unix=start,elapsed_seconds=time.monotonic()-t,exit_code=r.returncode));require(r.returncode!=0,'entrypoint incorrectly authorized');require(json.loads((D/(name+'-audit.json')).read_text())==[],'rejected entrypoint spawned/connected/read credential');return cmd
rejected('wrong-pin',['rehearse','--surface','cli','--run-dir',str(PREP/'rehearsals/must-not-exist'),'--pins',str(PREP/'rehearsal-pins.json'),'--pins-sha256','0'*64])
rejected('preauthorization',['execute','--exercise','paid-cli'])
# Execute a genuine fail-fast command sequence; the second process would create a witness file.
commands=[[sys.executable,'-B',str(PREP/'controller.py'),'execute','--exercise','paid-cli'],[sys.executable,'-c','import pathlib;pathlib.Path('+repr(str(D/'next-case-ran'))+').write_text("UNEXPECTED")']];attempts=[]
for cmd in commands:
 r=subprocess.run(cmd,env=env,cwd=PREP.parent.parent,capture_output=True,text=True);attempts.append(dict(argv=cmd,exit_code=r.returncode,stdout=r.stdout,stderr=r.stderr))
 if r.returncode:break
require(len(attempts)==1 and not (D/'next-case-ran').exists(),'first failure entered next case');save(D/'first-failure.json',dict(attempts=attempts,next_case_entered=False))
local=D/'http';local.mkdir();peer=proxy=None;checks=[]
try:
 peer=LocalPeer(local,False);proxy=Proxy(local,'cli','http://127.0.0.1:'+str(peer.port));proxy.activate(60)
 def request(body):
  c=http.client.HTTPConnection('127.0.0.1',proxy.port,timeout=15)
  try:c.request('POST','/v1/messages',json.dumps(body).encode(),{'Content-Type':'application/json'});r=c.getresponse();r.read();return r.status
  finally:c.close()
 good=dict(model=POLICY['model'],max_tokens=1024,messages=[dict(role='user',content='guard probe')],system='nonsecret guard probe',tools=[],stream=False)
 for name,changed in [('wrong-model',dict(good,model='unapproved-model')),('wrong-token',dict(good,max_tokens=1025))]:
  status=request(changed);require(status==502 and len(peer.requests)==0,'invalid request reached upstream');checks.append(dict(name=name,status=status,upstream_requests=0))
 for n in range(8):require(request(good)==200,'allowed local reservation failed')
 status=request(good);require(status==502 and len(peer.requests)==8 and len(proxy.records)==8,'ninth request forwarded');checks.append(dict(name='over-budget',status=status,reservations=8,upstream_requests=8))
 save(D/'http-guards.json',dict(ok=True,checks=checks,real_http=True,external_http=False,credentials=False))
finally:
 if proxy:proxy.close()
 if peer:peer.close()
save(D/'summary.json',dict(ok=True,wrong_pin_no_entry=True,preauthorization_no_entry=True,wrong_model_no_forward=True,wrong_tokens_no_forward=True,ninth_request_no_forward=True,first_failure_no_next_case=True));print('COMPACT GUARDS PASSED')
