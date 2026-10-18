#!/usr/bin/env python3
"""Exactly two independently approved Chapter 9 exercises; no provider matrix, retry, or fallback."""
import argparse,json,os,pathlib,shutil,signal,subprocess,sys,time,websocket
from launcher import ROOT,PREP,POLICY,Failure,require,save,sha,verify_frozen,LocalPeer,Proxy,Artifact
PAID_PROMPT='Load code-tools and follow its instructions: write note.txt with exact content SKILLS_OK, then read note.txt to verify it. Request unload_skill for code-tools. Answer with the verified file content and explain that unload is pending. Explain each tool batch, perform each action once, and do not perform unrelated actions.'
LOCAL_PROMPT=PAID_PROMPT+' After requesting unload and before answering, use write_file to write after-unload.txt with exact content SKILLS_OK, demonstrating that the dynamic write capability remains available while unload is pending.'
def alarm(*args):raise Failure('finite exercise deadline reached')
def classify(messages,require_done=True):
 errors=[];starts={};ends={};turns=[];state=None;paused=False;jobs=[]
 for m in messages:
  kind=m.get('type')
  if kind in ('error','settings_error'):errors.append({'type':kind,'text':m.get('text','')})
  elif kind=='turn_ended':
   turns.append(m)
   if m.get('error'):errors.append({'type':kind,'error':m['error']})
  elif kind=='tool_dispatched':
   key=m.get('call_id');require(key,'dispatch missing call_id');require(key not in starts,'duplicate dispatch call_id');starts[key]=m
  elif kind=='tool_finished':
   key=m.get('call_id');require(key,'tool finish missing call_id');require(key not in ends,'duplicate finish call_id');ends[key]=m
   if m.get('is_error') is True:errors.append({'type':kind,'call_id':key,'result':m.get('result')})
  elif kind=='state_changed':state=m.get('to');paused=m.get('paused',paused)
  elif kind=='pause_changed':paused=m.get('paused',False)
  elif kind=='job_report':
   j=m.get('job') or {};jobs.append(j)
   if j.get('exit_code') not in (None,0) or j.get('status') in ('failed','killed','canceled','cancelled'):errors.append({'type':'job_report','job':j})
 complete=bool(turns) and set(starts)==set(ends) and state=='idle' and not paused
 if require_done and not complete:errors.append({'type':'incomplete typed turn','state':state,'paused':paused,'started':list(starts),'finished':list(ends),'turn_count':len(turns)})
 return dict(ok=not errors,errors=errors,complete=complete,state=state,paused=paused,dispatch_by_call_id=starts,finished_by_call_id=ends,turns=turns,jobs=jobs,notice_is_not_error_test=True)
def dom(b,s):
 return b.eval(s,"""(()=>{const rect=n=>{let r=n.getBoundingClientRect();return {x:r.x,y:r.y,width:r.width,height:r.height}};return {state:document.querySelector('#state').textContent,paused:!document.querySelector('#paused').hidden,notice:document.querySelector('#notice').textContent,read_aloud:document.querySelector('#read-aloud').checked,connected:document.querySelector('#connection').textContent,transcript:document.querySelector('#transcript').innerText,actions:document.querySelector('#actions').innerText,cards:[...document.querySelectorAll('#actions article')].map(n=>({id:n.dataset.partId,title:n.querySelector('.artifact-title')?.textContent,text:n.innerText,results:n.querySelectorAll('.tool-result').length})),pane_bounds:Object.fromEntries(['.sidebar','.chat-pane','.actions-pane'].map(x=>[x,rect(document.querySelector(x))])),body_overflow:document.documentElement.scrollWidth>innerWidth}}})()""")
def facts(directory,proxy,dynamic,peer=None):
 requests=proxy.records;require(0<len(requests)<=8,'request count missing/excess');require(all(r.get('complete') and r.get('upstream_status')==200 for r in requests),'failed/incomplete exchange')
 require(all(r['system_sha256']==requests[0]['system_sha256'] and r['primary_markers_present'] for r in requests),'primary rendering changed or absent')
 require(requests[0]['tools']==['load_skill','read_file','unload_skill'],'initial tools not gated')
 require(all('write_file' in r['tools'] and 'list_directory' in r['tools'] for r in requests[1:]),'dynamic or dependency tools missing')
 events=[json.loads(line) for line in (directory/'event.log').read_text().splitlines() if line.strip()];calls=[e for e in events if e.get('type')=='tool_called'];returns=[e for e in events if e.get('type')=='tool_returned'];expected=['load_skill','write_file','read_file','unload_skill']+(['write_file'] if dynamic else [])
 require([e['tool']['name'] for e in calls]==expected,'unexpected actual tool sequence')
 require([e['tool']['call_id'] for e in calls]==[e['tool']['call_id'] for e in returns] and len({e['tool']['call_id'] for e in calls})==len(expected),'missing/mismatched durable tool result')
 require(all(not e['tool'].get('is_error') for e in returns),'actual tool failure')
 payload=json.dumps(returns);require(all(x in payload for x in ('CH09_HELPER','CH09_CODE_TOOLS','CH09_CALLBACK_VALUE','$UNRESOLVED_CH09','SKILLS_OK')),'helper/rendering/read evidence missing');require('pending' in payload.lower(),'pending unload result missing')
 file_names=['note.txt']+(['after-unload.txt'] if dynamic else [])
 for name in file_names:require((directory/'workspace'/name).read_text()=='SKILLS_OK','actual file content mismatch')
 body=dict(request_count=len(requests),native_cap=requests[0]['max_output_value'],unchanged_system_hash=requests[0]['system_sha256'],initial_tools=requests[0]['tools'],next_tools=requests[-1]['tools'],calls=[dict(seq=e['seq'],call_id=e['tool']['call_id'],name=e['tool']['name']) for e in calls],results=[dict(seq=e['seq'],call_id=e['tool']['call_id'],is_error=e['tool'].get('is_error',False),parts=e['tool'].get('parts',[])) for e in returns],files={n:dict(content='SKILLS_OK',sha256=sha(directory/'workspace'/n)) for n in file_names})
 if peer:
  require(len(requests)==len(peer.requests)==len(expected)+1,'scripted request count mismatch');require(all(a['body_sha256']==b['body_sha256'] for a,b in zip(requests,peer.requests)),'forwarded body was rewritten');require(peer.requests[0]['header_delay_seconds']>3,'delayed-header rehearsal missing')
  body.update(unchanged_forwarded_body_hashes=True,delayed_headers_seconds=peer.requests[0]['header_delay_seconds'])
 save(directory/'sanitized-effects.json',body);return body
def socket_turn(app,directory):
 ws=websocket.create_connection(app.url.replace('http','ws')+'ws',timeout=2,suppress_origin=True);messages=[]
 try:
  ws.send(json.dumps({'type':'subscribe'}));ws.send(json.dumps({'type':'prompt','text':LOCAL_PROMPT}));end=time.monotonic()+90
  with (directory/'wire.jsonl').open('w') as out:
   while time.monotonic()<end:
    try:m=json.loads(ws.recv())
    except websocket.WebSocketTimeoutException:continue
    messages.append(m);out.write(json.dumps(m)+'\n');out.flush();state=classify(messages,False);require(state['ok'],'actual typed error')
    if state['complete']:break
  result=classify(messages);save(directory/'typed-facts.json',result);require(result['ok'],'typed turn incomplete')
 finally:ws.close()
def browser_turn(app,directory):
 from browser import Browser
 b=Browser.__new__(Browser)
 try:
  b.__init__(directory/'browser-profile',True,directory,'owned');s=b.page(app.url);before=dom(b,s);require(not before['read_aloud'],'readAloud must be false');initial=b.shot(s,'initial')
  b.draft(s,LOCAL_PROMPT);b.click(s,'#submit');end=time.monotonic()+540
  while time.monotonic()<end:
   view=dom(b,s);messages=[r['message'] for r in b.wire if r['session']==s];state=classify(messages,False);require(state['ok'],'actual typed browser error')
   if state['complete']:break
   time.sleep(.1)
  state=classify(messages);require(state['ok'],'browser turn incomplete')
  # Observe settled frames on the one owned page; no restart, second page or settings campaign.
  b.eval(s,'new Promise(r=>requestAnimationFrame(()=>requestAnimationFrame(r)))');after=dom(b,s)
  require(after['state']=='idle' and not after['paused'] and not after['read_aloud'],'final browser state unsafe')
  require(len(after['cards'])==5 and all(c['results']==1 for c in after['cards']),'missing/duplicate paired tool cards')
  require(all(x in after['actions'] for x in ('CH09_HELPER','CH09_CALLBACK_VALUE','SKILLS_OK')) and 'pending' in after['actions'].lower(),'rendered skill/result evidence missing')
  require('SKILLS_OK' in after['transcript'],'verified content absent from answer')
  save(directory/'browser-wire.json',messages);save(directory/'typed-facts.json',state);save(directory/'browser-facts.json',dict(before=before,after=after,screenshots=[initial,b.shot(s,'skills-result')],visual_inspection='Captured only; root must inspect screenshots.',native_speech=False,microphone=False,restart=False,extra_page=False))
 finally:b.close()
def verify_harness(pin_file,pin_sha):
 require(pin_file and pin_sha,'preauthorization/pins required');p=pathlib.Path(pin_file);require(p.parent==PREP and p.name in ('rehearsal-pins.json','seal.json'),'unapproved pin file');require(sha(p)==pin_sha,'harness pin mismatch')
 pins=json.loads(p.read_text())
 for row in pins['preparation_files']:require(sha(PREP/row['path'])==row['sha256'],'sealed preparation changed: '+row['path'])
 verify_frozen();return pins
def approve(path,exercise,policy_sha,seal_sha):
 require(path and policy_sha and seal_sha,'root scope and pin authorization required');require(exercise in POLICY['modes'],'exercise not permitted');require(sha(PREP/'policy.json')==policy_sha,'policy pin mismatch');verify_harness(str(PREP/'seal.json'),seal_sha)
 a=json.loads(pathlib.Path(path).read_text());require(a.get('policy_sha256')==policy_sha and a.get('seal_sha256')==seal_sha,'root approval pins differ');require(exercise in a.get('modes',[]) and set(a['modes'])<=set(POLICY['modes']),'root scope does not authorize mode');require(0<a.get('expires_unix',0)-time.time()<=3600,'approval expired/unbounded')
 run=pathlib.Path(a['run_directory']);require(run.is_absolute() and run.parent==pathlib.Path(POLICY['run_root']) and run.name not in ('','..','.'),'run location not owned');require(not (run/exercise).exists(),'one-shot exercise already reserved');return a,run

def session(directory,surface,local,browser=False):
 # Caller has verified pins and (if executing) exact independent root mode before this allocation.
 directory.mkdir(exist_ok=False);start=time.monotonic();peer=proxy=app=None;cleanup=[];result=dict(ok=False,surface=surface,local_only=local,browser=browser,started_unix=time.time());dynamic=surface=='gui'
 absolute=600 if dynamic else 240;signal.setitimer(signal.ITIMER_REAL,600 if dynamic else 60)
 try:
  verify_frozen();require(shutil.disk_usage(ROOT).free>=1024**3,'disk floor reached')
  if local:peer=LocalPeer(directory,dynamic);origin='http://127.0.0.1:'+str(peer.port)
  else:require(surface=='cli' and not browser,'paid GUI forbidden');origin=POLICY['origin']
  proxy=Proxy(directory,surface,origin,authorized_live=not local)
  app=Artifact.__new__(Artifact);app.__init__(directory,surface,proxy,stream=dynamic)
  setup=time.monotonic()-start;require(dynamic or setup<=60,'setup ceiling exceeded');active=min((600 if dynamic else 180),absolute-setup);require(active>0,'absolute deadline reached');proxy.activate(active);signal.setitimer(signal.ITIMER_REAL,active);result.update(setup_seconds=round(setup,3),active_budget_seconds=active,absolute_budget_seconds=absolute)
  if browser:browser_turn(app,directory)
  elif surface=='cli':
   app.p.communicate(json.dumps({'user':PAID_PROMPT})+'\n',timeout=active);require(app.p.returncode==0,'CLI nonzero exit');require('SKILLS_OK' in (directory/'artifact.stdout.txt').read_text(),'verified content absent from CLI answer')
  else:socket_turn(app,directory)
  end=time.monotonic()+3
  while any(not r.get('complete') for r in proxy.records) and time.monotonic()<end:time.sleep(.01)
  result['effects']=facts(directory,proxy,dynamic,peer);result['ok']=True
 except Exception as e:
  result['failure_type']=type(e).__name__;result['failure']=str(e) if isinstance(e,Failure) else 'harness exception; preserve owned local logs';raise
 finally:
  signal.setitimer(signal.ITIMER_REAL,0)
  for name,obj in [('artifact',app),('proxy',proxy),('upstream',peer)]:
   if obj:
    try:obj.close();cleanup.append(dict(owner=name,closed=True))
    except Exception as e:cleanup.append(dict(owner=name,closed=False,error_type=type(e).__name__));result['ok']=False
  result.update(cleanup=cleanup,elapsed_seconds=round(time.monotonic()-start,3))
  # Key never enters child env/argv; redact and fail if an upstream unexpectedly echoed it into local logs.
  if proxy and proxy.credential:
   needle=proxy.credential.encode()
   for p in directory.rglob('*'):
    if p.is_file() and 'browser-profile' not in p.parts and needle in p.read_bytes():p.write_bytes(p.read_bytes().replace(needle,b'[REDACTED_CREDENTIAL]'));result['ok']=False;result['secret_scan']='unexpected echo redacted; stop'
   proxy.credential=None
  save(directory/'result.json',result);require(all(c['closed'] for c in cleanup),'owned cleanup incomplete');require(result['ok'],'exercise failed; no next case')
 return result

RECEIPT_PATH=None
def main():
 global RECEIPT_PATH
 ap=argparse.ArgumentParser();ap.add_argument('mode',choices=['rehearse','execute']);ap.add_argument('--surface',choices=['cli','gui']);ap.add_argument('--run-dir');ap.add_argument('--pins');ap.add_argument('--pins-sha256');ap.add_argument('--exercise',choices=['browser-local','paid-cli']);ap.add_argument('--approval');ap.add_argument('--policy-sha256');ap.add_argument('--seal-sha256');a=ap.parse_args();signal.signal(signal.SIGALRM,alarm)
 if a.mode=='rehearse':
  verify_harness(a.pins,a.pins_sha256);require(a.surface and a.run_dir and not a.approval,'rehearsal surface/new local directory required');run=pathlib.Path(a.run_dir).resolve();require(run.parent==PREP/'rehearsals','rehearsal not owned');run.parent.mkdir(exist_ok=True);session(run,a.surface,True);print('LOCAL REAL-FORWARDING REHEARSAL PASSED '+str(run));return
 approval,run=approve(a.approval,a.exercise,a.policy_sha256,a.seal_sha256);run.mkdir(parents=True,exist_ok=True);RECEIPT_PATH=run/('receipt-'+a.exercise+'.json');save(run/('approval-'+a.exercise+'.json'),approval);save(run/('command-'+a.exercise+'.json'),dict(argv=sys.argv,cwd=os.getcwd(),environment_names=['PATH','HOME','TMPDIR','LANG','PYTHONDONTWRITEBYTECODE'],started_unix=time.time()))
 session(run/a.exercise,'gui' if a.exercise=='browser-local' else 'cli',a.exercise=='browser-local',browser=a.exercise=='browser-local');print('APPROVED SINGLE EXERCISE PASSED '+a.exercise)
if __name__=='__main__':
 started=time.time();tick=time.monotonic();code=0
 try:main()
 except Exception as e:
  code=1;print(json.dumps(dict(ok=False,type=type(e).__name__,message=str(e) if isinstance(e,Failure) else 'harness exception; stop and preserve local evidence')),file=sys.stderr)
 finally:
  if RECEIPT_PATH:save(RECEIPT_PATH,dict(argv=sys.argv,cwd=os.getcwd(),environment={k:os.environ[k] for k in ('PATH','HOME','TMPDIR','LANG','PYTHONDONTWRITEBYTECODE') if k in os.environ},started_unix=started,elapsed_seconds=time.monotonic()-tick,exit_code=code))
 sys.exit(code)
