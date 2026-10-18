#!/usr/bin/env python3
"""Finite Chapter 9 preparation rehearsal and root-approved execution. No retries or implicit matrix expansion."""
import argparse,collections,json,os,pathlib,shutil,signal,subprocess,sys,time,traceback,websocket
from launcher import ROOT,PREP,POLICY,Failure,require,save,sha,verify_frozen,LocalPeer,Proxy,Artifact,credential_for
PROMPT='Load the writer skill. Write first.txt with exact content CH09_FIRST. Request unload_skill for writer. Since unload is pending, then write after-unload.txt with exact content CH09_PENDING_STILL_ENABLED using the still-enabled tool. Perform each action once, explain each tool batch, and finish by replying CH09_DONE. Do not read other files or perform other actions.'
def alarm(*args):raise Failure('finite session or overall execution deadline reached')
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
def event_facts(directory):
 events=[json.loads(line) for line in (directory/'event.log').read_text().splitlines() if line.strip()]
 calls=[e for e in events if e.get('type')=='tool_called'];results=[e for e in events if e.get('type')=='tool_returned']
 def tool(e):return e.get('tool',{})
 names=[tool(e).get('name') for e in calls];call_ids=[tool(e).get('call_id') for e in calls];result_ids=[tool(e).get('call_id') for e in results]
 require(names==['load_skill','write_file','unload_skill','write_file'],'unexpected durable tool sequence: '+repr(names))
 require(call_ids==result_ids and len(set(call_ids))==4,'durable call pairing failed')
 require(all(not tool(e).get('is_error') for e in results),'durable tool error')
 text=json.dumps(results);require('CH09_DYNAMIC' in text and 'CH09_CALLBACK_VALUE' in text and '$UNRESOLVED_CH09' in text,'dynamic body variable evidence absent')
 require('pending' in text.lower(),'pending unload result absent')
 return dict(calls=[dict(seq=e.get('seq'),name=tool(e).get('name'),call_id=tool(e).get('call_id')) for e in calls],returns=[dict(seq=e.get('seq'),call_id=tool(e).get('call_id'),is_error=tool(e).get('is_error',False)) for e in results],registry_provenance_claim='Durable enabling event sequence recorded; internal registry source query was covered by frozen local tests, not exposed by CLI/GUI.')
def tool_names(body,vendor):
 tools=body.get('tools',[])
 if vendor=='gemini':return sorted(x.get('name') for t in tools for x in t.get('functionDeclarations',[]))
 return sorted(t.get('name') for t in tools)
def system(body,vendor):return body.get({'anthropic':'system','openai':'instructions','gemini':'systemInstruction'}[vendor])
def assert_effects(directory,proxy,local):
 requests=proxy.records;require(bool(requests) and len(requests)<=16,'missing/excess requests')
 require(all(r.get('complete') and r.get('upstream_status')==200 for r in requests),'incomplete or failed provider exchange')
 vendor=proxy.vendor;initial=system(requests[0]['body'],vendor);serialized=json.dumps(initial)
 require('CH09_PRIMARY' in serialized and 'CH09_CALLBACK_VALUE' in serialized and '$UNRESOLVED_CH09' in serialized,'primary rendering absent')
 require(all(system(r['body'],vendor)==initial for r in requests),'system text mutated by load')
 require(tool_names(requests[0]['body'],vendor)==['load_skill','read_file','unload_skill'],'initial declarations not gated')
 require(all('write_file' in tool_names(r['body'],vendor) and 'list_directory' in tool_names(r['body'],vendor) for r in requests[1:]),'dynamic/dependency declarations absent')
 for name,expected in [('first.txt','CH09_FIRST'),('after-unload.txt','CH09_PENDING_STILL_ENABLED')]:require((directory/'workspace'/name).read_text()==expected,'actual file effect mismatch: '+name)
 facts=event_facts(directory)
 if local:require(len(requests)==5,'local scripted lifecycle should forward exactly five requests')
 return dict(request_count=len(requests),output_fields=[dict(field=r['max_output_field'],value=r['max_output_value']) for r in requests],system_unchanged=True,initial_tools=tool_names(requests[0]['body'],vendor),dynamic_tools=tool_names(requests[-1]['body'],vendor),files={p.name:dict(bytes=p.stat().st_size,sha256=sha(p)) for p in sorted((directory/'workspace').iterdir()) if p.is_file()},durable=facts)
def websocket_turn(app,directory):
 ws=websocket.create_connection(app.url.replace('http','ws')+'ws',timeout=2,suppress_origin=True);messages=[]
 try:
  ws.send(json.dumps({'type':'subscribe'}));ws.send(json.dumps({'type':'prompt','text':PROMPT}));end=time.monotonic()+100
  with (directory/'wire.jsonl').open('w') as f:
   while time.monotonic()<end:
    try:m=json.loads(ws.recv())
    except websocket.WebSocketTimeoutException:continue
    messages.append(m);f.write(json.dumps(m)+'\n');f.flush();facts=classify(messages,False)
    require(facts['ok'],'actual typed GUI error')
    if facts['complete']:break
  facts=classify(messages);save(directory/'typed-facts.json',facts);require(facts['ok'],'incomplete actual GUI turn');return facts
 finally:ws.close()
def dom(b,s):
 return b.eval(s,"""(()=>{const rect=n=>{let r=n.getBoundingClientRect();return {x:r.x,y:r.y,width:r.width,height:r.height}};return {state:document.querySelector('#state').textContent,paused:!document.querySelector('#paused').hidden,notice:document.querySelector('#notice').textContent,read_aloud:document.querySelector('#read-aloud').checked,connected:document.querySelector('#connection').textContent,transcript:document.querySelector('#transcript').innerText,actions:document.querySelector('#actions').innerText,cards:[...document.querySelectorAll('#actions article')].map(n=>({id:n.dataset.partId,title:n.querySelector('.artifact-title')?.textContent,text:n.innerText,results:n.querySelectorAll('.tool-result').length})),pane_bounds:Object.fromEntries(['.sidebar','.chat-pane','.actions-pane'].map(x=>[x,rect(document.querySelector(x))])),body_overflow:document.documentElement.scrollWidth>innerWidth}}})()""")
def browser_turn(app,directory):
 from browser import Browser
 b=Browser.__new__(Browser);result={}
 try:
  b.__init__(directory/'browser-profile',True,directory,'owned');s=b.page(app.url);initial=dom(b,s)
  require(initial['connected']=='Connected' and not initial['read_aloud'],'read-aloud must be off');require(not initial['body_overflow'],'initial layout overflow')
  result['initial']=initial;result['screenshots']=[b.shot(s,'initial')]
  b.draft(s,PROMPT);b.click(s,'#submit');end=time.monotonic()+600
  while time.monotonic()<end:
   view=dom(b,s);messages=[r['message'] for r in b.wire if r['session']==s];facts=classify(messages,False);require(facts['ok'],'browser typed GUI error')
   if facts['complete'] and 'CH09_DONE' in view['transcript']:break
   time.sleep(.1)
  messages=[r['message'] for r in b.wire if r['session']==s];facts=classify(messages);save(directory/'browser-wire.json',messages);save(directory/'typed-facts.json',facts);require(facts['ok'],'browser turn incomplete')
  final=dom(b,s);require(not final['read_aloud'] and not final['paused'] and final['state']=='idle','browser completion state')
  require(len(final['cards'])==4 and all(x['results']==1 for x in final['cards']),'skill cards must have one paired result each')
  require('CH09_DYNAMIC' in final['actions'] and 'pending' in final['actions'].lower() and 'CH09_DONE' in final['transcript'],'actual rendered lifecycle content missing')
  result.update(final=final,typed=facts);result['screenshots'].append(b.shot(s,'lifecycle'))
  b.call('Page.reload',{},s);b.wait(s,"document.querySelector('#connection')?.textContent==='Connected' && document.querySelectorAll('#actions .tool-result').length===4")
  replay=dom(b,s);require(len(replay['cards'])==4 and all(x['results']==1 for x in replay['cards']) and not replay['read_aloud'],'replay duplication/settings mismatch')
  result['replay']=replay;result['screenshots'].append(b.shot(s,'replay'));result['visual_inspection']='Screenshots captured, not human-credited; root must inspect them.'
  save(directory/'browser-facts.json',result)
 finally:b.close()
def scan_secrets(directory,credential):
 if not credential:return
 needle=credential.encode()
 # Profiles contain non-text browser state and never receive credentials; raw user/API/controller text is checked.
 for p in directory.rglob('*'):
  if p.is_file() and not any(x in p.parts for x in ('browser-profile','owned-home')):
   raw=p.read_bytes()
   if needle in raw:
    p.write_bytes(raw.replace(needle,b'[REDACTED_CREDENTIAL]'));raise Failure('credential leak redacted; operational run failed')
def session(directory,surface,vendor,local,browser=False,deadline=None):
 directory.mkdir(exist_ok=False);began=time.monotonic();peer=proxy=app=None;result=dict(ok=False,surface=surface,vendor=vendor,local_only=local,browser=browser,started_unix=time.time());credential=None;cleanup=[]
 signal.setitimer(signal.ITIMER_REAL,min(POLICY['session_seconds'],deadline-time.monotonic()) if deadline else POLICY['session_seconds'])
 try:
  verify_frozen();require(shutil.disk_usage(ROOT).free>=1024**3,'disk guardrail')
  if local:peer=LocalPeer(directory);origin='http://127.0.0.1:'+str(peer.port)
  else:credential=credential_for(vendor);origin=POLICY['origins'][vendor]
  proxy=Proxy(directory,vendor,surface,origin,credential)
  app=Artifact.__new__(Artifact);app.__init__(directory,surface,vendor,proxy,stream=browser)
  if browser:browser_turn(app,directory)
  elif surface=='cli':
   app.p.communicate(json.dumps({'user':PROMPT})+'\n',timeout=100 if local else max(1,(deadline or (began+720))-time.monotonic()));require(app.p.returncode==0,'CLI nonzero exit')
   require('CH09_DONE' in (directory/'artifact.stdout.txt').read_text(),'CLI final answer absent')
  else:websocket_turn(app,directory)
  # Every HTTP response has drained before the typed turn completes; wait only for its bookkeeping.
  end=time.monotonic()+3
  while any(not r.get('complete') for r in proxy.records) and time.monotonic()<end:time.sleep(.01)
  result['effects']=assert_effects(directory,proxy,local);result['ok']=True
 except Exception as e:
  result['failure_type']=type(e).__name__;result['failure']=str(e) if isinstance(e,Failure) else 'execution exception; see owned nonsecret process logs';raise
 finally:
  signal.setitimer(signal.ITIMER_REAL,0)
  for name,obj in [('artifact',app),('proxy',proxy),('upstream',peer)]:
   if obj:
    try:obj.close();cleanup.append({'owner':name,'closed':True})
    except Exception as e:cleanup.append({'owner':name,'closed':False,'error_type':type(e).__name__});result['ok']=False
  result['cleanup']=cleanup;result['elapsed_seconds']=round(time.monotonic()-began,3)
  try:scan_secrets(directory,credential)
  except Exception:result['ok']=False;result['secret_scan']='redaction required';raise
  finally:credential=None;save(directory/'result.json',result)
  require(all(x['closed'] for x in cleanup),'cleanup incomplete')
 return result
def approval(path,expected_policy,expected_seal):
 require(sha(PREP/'policy.json')==expected_policy and sha(PREP/'seal.json')==expected_seal,'approval CLI pin mismatch')
 seal=json.loads((PREP/'seal.json').read_text())
 for r in seal['preparation_files']:require(sha(PREP/r['path'])==r['sha256'],'sealed harness changed')
 a=json.loads(pathlib.Path(path).read_text());require(a.get('policy_sha256')==expected_policy and a.get('seal_sha256')==expected_seal,'approval pin mismatch')
 require(a.get('modes')==POLICY['execution_order'],'approval must name exact finite baseline order')
 require(0<a.get('expires_unix',0)-time.time()<=7200,'approval expired or unbounded')
 run=pathlib.Path(a['run_directory']);allowed=pathlib.Path(POLICY['run_root']);require(run.is_absolute() and run.parent==allowed and run.name not in ('','..','.'),'unapproved run path')
 require(not run.exists(),'one-shot run directory already consumed');verify_frozen();return a,run
def main():
 p=argparse.ArgumentParser();p.add_argument('mode',choices=['rehearse','execute']);p.add_argument('--run-dir');p.add_argument('--approval');p.add_argument('--policy-sha256');p.add_argument('--seal-sha256');args=p.parse_args();signal.signal(signal.SIGALRM,alarm)
 if args.mode=='rehearse':
  require(args.run_dir and not args.approval,'rehearsal takes only a new local run directory');run=pathlib.Path(args.run_dir).resolve();require(run.parent==PREP/'rehearsals','rehearsal outside owned preparation');run.mkdir(parents=True,exist_ok=False)
  save(run/'command.json',dict(argv=sys.argv,cwd=os.getcwd(),mode='local real-forwarding rehearsal; no browser/provider/key access',started_unix=time.time()))
  for surface in ('cli','gui'):session(run/surface,surface,'anthropic',True)
  save(run/'summary.json',dict(ok=True,CLI_measured_cap=1024,GUI_measured_cap=4096,real_forwarding=True,browser=False,provider=False,credentials=False));print('REHEARSAL PASSED '+str(run));return
 require(args.approval and args.policy_sha256 and args.seal_sha256,'execute requires explicit root approval and exact pins')
 a,run=approval(args.approval,args.policy_sha256,args.seal_sha256);run.parent.mkdir(exist_ok=True);run.mkdir(exist_ok=False);save(run/'approval-copy.json',a);save(run/'command.json',dict(argv=sys.argv,cwd=os.getcwd(),policy_sha256=args.policy_sha256,seal_sha256=args.seal_sha256,source_commit=POLICY['source_commit']))
 paid_deadline=None;results=[]
 try:
  for mode in POLICY['execution_order']:
   local=mode=='browser-local';is_gui=local or mode=='paid-gui';vendor='anthropic' if is_gui else mode.removeprefix('paid-')
   if not local and paid_deadline is None:paid_deadline=time.monotonic()+POLICY['paid_seconds']
   if paid_deadline:require(time.monotonic()<paid_deadline,'48-minute paid ceiling reached')
   results.append(session(run/mode,'gui' if is_gui else 'cli',vendor,local,browser=is_gui,deadline=paid_deadline))
  print('FINITE EXECUTION PASSED '+str(run))
 finally:save(run/'execution-summary.json',dict(ok=len(results)==5 and all(r['ok'] for r in results),completed=len(results),planned=5,results=results,no_retry=True))
if __name__=='__main__':
 try:main()
 except Exception as e:
  print(json.dumps(dict(ok=False,type=type(e).__name__,message=str(e) if isinstance(e,Failure) else 'harness exception; preserve owned artifacts and stop')),file=sys.stderr);sys.exit(1)
