from pathlib import Path
import json,hashlib,re,os,socket,subprocess,datetime
R=Path('/Users/bill/projects/ensemble-edition-3-student');B=R/'.ch08-local/operational-scroll-continuation';L=B/'runs/controller-scroll-local-native-01';LG=B/'runs/scroll-local-native-01';P=B/'runs/controller-scroll-paid-01';PG=B/'runs/scroll-paid-01'
h=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
def info(p):return {'path':str(p),'bytes':p.stat().st_size,'sha256':h(p)}
def save(p,x):p.write_text(json.dumps(x,indent=2,ensure_ascii=False)+'\n')
l=json.loads((L/'result.json').read_text());paid=json.loads((P/'result.json').read_text());lc=json.loads((LG/'completion.json').read_text());pc=json.loads((PG/'completion.json').read_text());assert l['status']==paid['status']=='pass';assert len(lc['upstream_events'])==3 and lc['key_reads']==0 and lc['external_requests']==0;assert pc['key_reads']==1 and pc['external_requests']==2 and len(pc['upstream_events'])==2
obs={x['name']:x['data'] for x in l['observations']};rates=[]
for i in range(1,4):
 a=obs['native phrase '+str(i)];q=next(e for e in a['events'] if e['event']=='rate');st=next(e for e in a['events'] if e['event']=='start');en=next(e for e in a['events'] if e['event']=='end');assert q['requested']==a['requested'] and abs(q['rate']-a['expectedNativeRate'])<1e-5 and not any(e['event']=='error' for e in a['events']);rates.append({'requested':a['requested'],'expected_native':a['expectedNativeRate'],'native_getter':q['rate'],'events':a['events'],'start_to_end_ms':en['at']-st['at']})
fr=obs['fresh disable-cancel correlation'];release=obs['native fresh cancellation / remaining draft pause / Interrupt release'];assert fr['freshCancelIndex']>fr['fourthStartIndex'] and fr['freshCancelEvent']['at']>=fr['before']['at'] and fr['before']['speaking'] and not fr['after']['speaking'];assert fr['after']['paused'] and fr['after']['draft']=='held draft' and not release['after']['paused'] and release['typedFacts']['ok'];assert l['native_nonempty_utterances']==4
save(B/'local-native-safe-facts.json',{'source_result':info(L/'result.json'),'status':'pass','local_request_receipt':lc,'actual_scroll':obs['actual tool completion / pane scrolling']['scroll'],'overflow_observation':obs['constrained-scroll viewport observation'],'restored_normal_observation':obs['restored viewport observation'],'two_page_restart_persistence':obs['one owned restart and two-page persistence'],'two_GUI_ownership_records':[info(p) for p in LG.glob('ownership-*.json')],'native_requested_nonempty':4,'native_rates':rates,'fourth_boundary':obs['fourth utterance and fresh-cancel boundary'],'fresh_cancel':fr,'release':{'overlap_paused':release['overlap']['paused'],'after_paused':release['after']['paused'],'after_draft':release['after'].get('draft'),'typed_facts':release['typedFacts'],'elapsed_seconds':release['elapsed_seconds']},'default_voice_metadata':obs['native platform/default qualification'],'qualifications':['Actual production native events/getters; no native API/callback replacement.','Default voice left unmodified; reported default metadata does not independently prove actual backend/locality. Locality not an acceptance condition.','Three completed phrases plus fourth actively canceled. Stale-generation callbacks may be intentionally ignored by production; do not claim observation of every native callback.','No human listening report or screenshot viewing claim.']})
# Read only this owned API response trace. Never read any real-user settings/key file.
api=PG/'api.log';events=[];decoder=json.JSONDecoder();allowed={'message_start','message_delta','message_stop','content_block_start','content_block_delta','content_block_stop'}
def walk(x):
 if isinstance(x,dict):
  if x.get('type') in allowed:events.append(x);return
  for v in x.values():walk(v)
 elif isinstance(x,list):
  for v in x:walk(v)
for line in api.read_text().splitlines():
 for match in re.finditer(r'\{',line):
  try:x,end=decoder.raw_decode(line[match.start():])
  except ValueError:continue
  walk(x);break
responses=[];current=None;blocks={};usage_keys=['input_tokens','output_tokens','cache_creation_input_tokens','cache_read_input_tokens']
def usage(x):return {k:x[k] for k in usage_keys if isinstance(x.get(k),int)}
for e in events:
 t=e['type']
 if t=='message_start':
  msg=e['message'];current={'message_id':msg['id'],'model':msg['model'],'start_usage':usage(msg.get('usage',{})),'tools':[]};responses.append(current);blocks={}
 elif current is not None and t=='content_block_start' and e.get('content_block',{}).get('type')=='tool_use':
  q=e['content_block'];blocks[e['index']]={'id':q['id'],'name':q['name'],'input':q.get('input',{}),'partial':''}
 elif current is not None and t=='content_block_delta' and e.get('index') in blocks and e.get('delta',{}).get('type')=='input_json_delta':blocks[e['index']]['partial']+=e['delta']['partial_json']
 elif current is not None and t=='content_block_stop' and e.get('index') in blocks:
  q=blocks.pop(e['index']);part=q.pop('partial');q['input']=json.loads(part) if part else q['input'];current['tools'].append(q)
 elif current is not None and t=='message_delta':current['final_usage']=usage(e.get('usage',{}));current['stop_reason']=e.get('delta',{}).get('stop_reason')
 elif current is not None and t=='message_stop':current['message_stop_seen']=True
assert len(responses)==2 and all(x.get('message_stop_seen') for x in responses)
d=paid['paid']['dispatch'];f=paid['paid']['finished'];assert len(d)==len(f)==1 and d[0]['call_id']==f[0]['call_id'] and d[0]['name']=='read_file' and d[0]['input']['path']=='exercise.txt' and not f[0]['is_error'];assert responses[0]['tools']==[{'id':d[0]['call_id'],'name':'read_file','input':d[0]['input']}] and not responses[1]['tools'];assert responses[0]['stop_reason']=='tool_use' and responses[1]['stop_reason']=='end_turn'
for x in responses:x['reported_usage']={k:x.get('final_usage',{}).get(k,x['start_usage'].get(k)) for k in usage_keys};assert all(isinstance(v,int) for v in x['reported_usage'].values())
totals={k:sum(x['reported_usage'][k] for x in responses) for k in usage_keys};turns=[{'text':x.get('text',''),'error':x.get('error','')} for x in paid['paid']['turns']];assert turns and not any(t['error'] for t in turns);seed=PG/'workspace/exercise.txt';assert seed.read_text()=='Chapter 8 boundary check.\n';assert 'Chapter 8 boundary check.' in f[0]['result']
safe={'source_result':info(P/'result.json'),'source_API_trace':info(api),'status':'pass','paid_invocations':1,'upstream_receipt':pc,'responses':responses,'usage_totals':totals,'usage_method':'Final message_delta usage per response; fall back to same response message_start only for fields not present at final. Never add provisional output1 to final64/24. No USD price/cost claim.','request_context_qualification':'Proxy captured exact outbound body hashes/model/max_tokens; API trace captured response IDs/usage/tool-use. Full next-request content/headers are not claimed as captured here.','dispatch':[{'call_id':x['call_id'],'name':x['name'],'input':x['input']} for x in d],'finished':[{'call_id':x['call_id'],'is_error':x['is_error'],'result':x['result']} for x in f],'turns':turns,'completion_facts':paid['paid_completion_facts'],'owned_seed_file':{**info(seed),'text':seed.read_text()},'actual_final_DOM_state':{'state':paid['paid']['final']['state'],'paused':paid['paid']['final']['paused'],'notice':paid['paid']['final']['notice']},'no_retry_fallback_discovery_or_body_rewrite':True,'key_privacy':'One sealed-proxy-only key read; HTTPS authorization only. GUI dummy environment. No real key/headers/settings exported. No standalone key read by these tools.','no_raw_thinking_or_API_body_export':True};save(B/'paid-safe-facts.json',safe)
# Confirm all recorded owned lifecycle identities without any new browser/provider call.
pids=set();ports=set();profiles=[];exit_receipts=[]
def collect(o):
 if isinstance(o,dict):
  for k,v in o.items():
   if k in ('pid','launcher_pid') and isinstance(v,int) and v>1:pids.add(v)
   if 'listener' in k and isinstance(v,list) and len(v)==2 and isinstance(v[1],int):ports.add(v[1])
   collect(v)
 elif isinstance(o,list):
  for v in o:collect(v)
for meta,gui in [(L,LG),(P,PG)]:
 for p in list(meta.glob('*ownership.json'))+list(meta.glob('*attached.json'))+[meta/'launcher-command.json']+list(gui.glob('*.json')):collect(json.loads(p.read_text()))
 for p in meta.glob('*-browser-ownership.json'):profiles.append(json.loads(p.read_text())['profile'])
 for p in list(meta.glob('*-browser-exit.json'))+[meta/'launcher-exit.json']+list(gui.glob('exit-*.json')):
  q=json.loads(p.read_text());assert q['exit_code']==0;exit_receipts.append(info(p));collect(q)
 for p in meta.glob('*chrome.log'):pids.update(int(x) for x in re.findall(r'^\[(\d+):',p.read_text(errors='replace'),re.M) if int(x)>1)
for pid in pids:
 try:os.kill(pid,0);raise AssertionError('owned PID survives')
 except ProcessLookupError:pass
port_receipts=[]
for port in sorted(ports):
 assert port!=8088;s=socket.socket();s.settimeout(.15);code=s.connect_ex(('127.0.0.1',port));s.close();assert code!=0;port_receipts.append({'port':port,'connect_ex':code,'closed':True})
ps=subprocess.check_output(['ps','-axo','pid=,ppid=,command='],text=True);survivors=[x for x in ps.splitlines() if any(p in x for p in profiles)];assert not survivors
save(B/'owned-cleanup.json',{'utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'recorded_pids':sorted(pids),'alive':[],'listeners':port_receipts,'owned_profile_survivors':survivors,'exit_receipt_references':exit_receipts,'all_recorded_exits_zero':True,'serving_processes_stopped':lc['processes_stopped'] and pc['processes_stopped'],'no_shared_cleanup':True,'component_scope':'Serving completion browser_launched:false describes launcher only; three actual owned Chrome sessions have independent ownership/attachment/exit receipts.'})
shots=[]
for meta in [L,P]:
 for p in meta.glob('*.png'):shots.append({**info(p),'captured':True,'visually_inspected':False})
save(B/'screenshot-index.json',{'qualification':'PNG capture and DOM/geometry/event reads only. Student did not visually view any captured screenshot via an image/browser observation capability; no new image/provider call. Raw PNGs stay local.','entries':shots})
print('SAFE FACTS:',len(responses),'paid responses; usage',totals,'; one read_file ID matches native API tool-use / dispatch / finished; normal final answer preserved.')
print('LOCAL/NATIVE:3 fake requests,1 GUI restart,4 native nonempty phrases; three completed native rates and fresh fourth cancellation. Human listening not reported.')
print('CLEANUP:',len(pids),'recorded PIDs absent,',len(ports),'listeners closed; all exits0;',len(shots),'PNGs captured, none visually viewed.')
