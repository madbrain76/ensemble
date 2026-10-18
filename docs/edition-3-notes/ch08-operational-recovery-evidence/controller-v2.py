#!/usr/bin/env python3
"""Bounded actual-owned-GUI/CDP orchestration. No key access, native replacement, or provider client."""
import argparse,base64,hashlib,http.client,json,os,signal,socket,subprocess,sys,time,traceback,urllib.request
from pathlib import Path
import websocket
ROOT=Path('/Users/bill/projects/ensemble-edition-3-student');PREP=ROOT/'.ch08-local/operational-prep'
POLICY='eef19bf051df623c71d6227f1e375b48098177009e1d801f50f713e3533e6b27'
CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def save(p,v):Path(p).write_text(json.dumps(v,indent=2)+'\n')
def env(home):
 d={k:os.environ[k] for k in ['PATH','TMPDIR','LANG'] if k in os.environ};d.update(HOME=str(home),PYTHONDONTWRITEBYTECODE='1');return d
def waitfile(p,seconds=25):
 end=time.monotonic()+seconds
 while time.monotonic()<end:
  if p.exists():return json.loads(p.read_text())
  time.sleep(.05)
 raise RuntimeError('setup receipt not ready: '+str(p))
class Failure(Exception):pass
class GUI:
 def __init__(self,mode,directory,meta):
  self.d=directory;self.meta=meta;self.out=(meta/'launcher.stdout').open('w');self.began=time.monotonic()
  args=[sys.executable,'-B',str(PREP/'launcher.py'),'serve-'+mode,'--policy',str(PREP/'policy.json'),'--policy-sha256',POLICY,'--approval','/tmp/ch08-operational-'+('local' if mode=='local' else 'paid')+'-approval.json','--run-dir',str(directory)]
  self.p=subprocess.Popen(args,cwd=ROOT,env=env('/Users/bill'),stdout=self.out,stderr=subprocess.STDOUT)
  save(meta/'launcher-command.json',{'pid':self.p.pid,'argv':args,'cwd':str(ROOT),'env_names':sorted(env('/Users/bill')),'launcher_sha256':sha(PREP/'launcher.py'),'policy_sha256':POLICY})
  self.number=1;self.own=self.ready()
 def ready(self):
  p=self.d/f'ownership-{self.number}.json';end=time.monotonic()+25
  while time.monotonic()<end:
   if self.p.poll() is not None:raise RuntimeError('owned launcher exited during setup')
   if p.exists():
    r=json.loads(p.read_text())
    if r.get('served_assets_verified'):return r
   time.sleep(.05)
  raise RuntimeError('GUI setup timeout')
 @property
 def url(self):return 'http://127.0.0.1:'+str(self.own['gui_listener'][1])+'/'
 def restart(self):
  assert self.number==1;self.p.send_signal(signal.SIGUSR1);self.number=2;self.own=self.ready()
 def stop(self):
  if self.p.poll() is None:
   self.p.terminate()
   try:self.p.wait(timeout=25)
   except subprocess.TimeoutExpired:self.p.kill();self.p.wait(timeout=10)
  self.out.close();save(self.meta/'launcher-exit.json',{'pid':self.p.pid,'exit_code':self.p.returncode,'elapsed_seconds':round(time.monotonic()-self.began,3)})
class Browser:
 def __init__(self,directory,muted,meta,label,mock_keychain=True):
  self.meta=meta;self.label=label;self.trace=(meta/(label+'-raw-cdp.jsonl')).open('w');self.events=[];self.wire=[];self.ids=0;self.ws=None
  directory.mkdir(exist_ok=False);home=meta/(label+'-home');home.mkdir();log=(meta/(label+'-chrome.log')).open('w');self.log=log
  args=[CHROME,'--user-data-dir='+str(directory),'--remote-debugging-address=127.0.0.1','--remote-debugging-port=0','--no-first-run','--no-default-browser-check','--disable-background-networking','--disable-component-update','--disable-sync','--disable-extensions','--password-store=basic','--window-size=1280,900']
  if mock_keychain:args+=['--use-mock-keychain']
  if muted:args+=['--mute-audio']
  args+=['about:blank'];self.p=subprocess.Popen(args,env=env(home),stdout=log,stderr=subprocess.STDOUT)
  self.own={'pid':self.p.pid,'argv':args,'executable_sha256':sha(CHROME),'version':subprocess.check_output([CHROME,'--version'],text=True).strip(),'profile':str(directory),'environment_names':sorted(env(home)),'muted':muted,'controller_sha256':sha(__file__),'attached':False}
  save(meta/(label+'-browser-ownership.json'),self.own)
  portfile=directory/'DevToolsActivePort';end=time.monotonic()+25
  while not portfile.exists() and time.monotonic()<end:
   if self.p.poll() is not None:raise RuntimeError('owned Chrome exited during setup')
   time.sleep(.05)
  lines=portfile.read_text().splitlines();port=int(lines[0]);assert port!=8088
  proof=subprocess.check_output(['/usr/sbin/lsof','-nP','-a','-p',str(self.p.pid),'-iTCP:'+str(port),'-sTCP:LISTEN'],text=True)
  assert '127.0.0.1:'+str(port) in proof and 'LISTEN' in proof
  command=subprocess.check_output(['ps','-p',str(self.p.pid),'-o','command='],text=True);assert str(directory) in command
  self.own.update(cdp_listener=['127.0.0.1',port],listener_proof=proof,process_command=command)
  save(meta/(label+'-browser-ownership.json'),self.own) # durable verified ownership BEFORE websocket attachment
  self.ws=websocket.create_connection('ws://127.0.0.1:'+str(port)+lines[1],timeout=2,suppress_origin=True)
  self.own['attached']=True;save(meta/(label+'-browser-attached.json'),self.own)
 def call(self,method,params=None,session=None,timeout=30):
  self.ids+=1;n=self.ids;m={'id':n,'method':method,'params':params or {}}
  if session:m['sessionId']=session
  self.trace.write(json.dumps({'out':m})+'\n');self.trace.flush();self.ws.send(json.dumps(m));end=time.monotonic()+timeout
  while time.monotonic()<end:
   self.ws.settimeout(min(2,max(.01,end-time.monotonic())))
   try:r=json.loads(self.ws.recv())
   except websocket.WebSocketTimeoutException:continue
   self.trace.write(json.dumps({'in':r})+'\n');self.trace.flush()
   if r.get('id')==n:
    if 'error'in r:raise RuntimeError('CDP error '+str(r['error']))
    return r.get('result',{})
   self.events.append(r)
   if r.get('method')=='Network.webSocketFrameReceived':
    try:self.wire.append({'session':r.get('sessionId'),'message':json.loads(r['params']['response']['payloadData'])})
    except (KeyError,ValueError):pass
  raise RuntimeError('CDP timeout '+method)
 def page(self,url,second=False):
  if second:t=self.call('Target.createTarget',{'url':'about:blank','background':True})['targetId']
  else:t=next(x['targetId'] for x in self.call('Target.getTargets')['targetInfos'] if x['type']=='page')
  s=self.call('Target.attachToTarget',{'targetId':t,'flatten':True})['sessionId']
  for domain in ['Page','Runtime','Network']:self.call(domain+'.enable',session=s)
  self.call('Page.navigate',{'url':url},s);self.wait(s,"document.querySelector('#connection')?.textContent==='Connected'")
  return s
 def eval(self,s,expression):
  r=self.call('Runtime.evaluate',{'expression':expression,'returnByValue':True,'awaitPromise':True},s,timeout=95)
  if 'exceptionDetails'in r:raise RuntimeError('page evaluation: '+str(r['exceptionDetails']))
  return r.get('result',{}).get('value')
 def wait(self,s,expression,seconds=20):
  end=time.monotonic()+seconds;last=None
  while time.monotonic()<end:
   last=self.eval(s,expression)
   if last:return last
   time.sleep(.05)
  raise Failure('observation timeout: '+expression+' last='+repr(last))
 def click(self,s,selector):
  box=self.eval(s,"(()=>{let n=document.querySelector("+json.dumps(selector)+");if(!n)throw Error('missing visible selector');let r=n.getBoundingClientRect();return {x:r.x+r.width/2,y:r.y+r.height/2,w:r.width,h:r.height}})()")
  if not box['w'] or not box['h']:raise RuntimeError('not a visible selector: '+selector)
  for kind in ['mousePressed','mouseReleased']:self.call('Input.dispatchMouseEvent',{'type':kind,'x':box['x'],'y':box['y'],'button':'left','clickCount':1},s)
 def key(self,s,key):
  for t in ['keyDown','keyUp']:self.call('Input.dispatchKeyEvent',{'type':t,'key':key,'code':key},s)
 def widget(self,s,selector,value):
  return self.eval(s,"(()=>{const n=document.querySelector("+json.dumps(selector)+");n.value="+json.dumps(str(value))+";n.dispatchEvent(new Event('change',{bubbles:true}));return true})()")
 def draft(self,s,text):
  self.eval(s,"(()=>{const n=document.querySelector('#draft');n.value="+json.dumps(text)+";n.dispatchEvent(new Event('input',{bubbles:true}));return n.value})()")
 def shot(self,s,name):
  r=self.call('Page.captureScreenshot',{'format':'png','captureBeyondViewport':False},s);p=self.meta/(name+'.png');p.write_bytes(base64.b64decode(r['data']));return {'path':str(p),'sha256':sha(p),'visually_inspected':False}
 def close(self):
  if self.p.poll() is None:
   if self.ws:
    try:self.call('Browser.close',timeout=3)
    except Exception:pass
   try:self.p.wait(timeout=12)
   except subprocess.TimeoutExpired:self.p.terminate();self.p.wait(timeout=10)
  if self.ws:self.ws.close()
  self.log.close();self.trace.close();save(self.meta/(self.label+'-browser-exit.json'),{'pid':self.p.pid,'exit_code':self.p.returncode,'cdp_listener':self.own.get('cdp_listener'),'profile_preserved':True})
 def messages(self,s,kind):return [r['message'] for r in self.wire if r['session']==s and r['message'].get('type')==kind]

def dom(b,s):return b.eval(s,"""(()=>{let rect=n=>{let r=n.getBoundingClientRect();let c=getComputedStyle(n);return {x:r.x,y:r.y,w:r.width,h:r.height,scrollTop:n.scrollTop,scrollHeight:n.scrollHeight,clientHeight:n.clientHeight,overflowY:c.overflowY,color:c.color,background:c.backgroundColor,font:c.fontSize}};return {title:document.title,viewport:[innerWidth,innerHeight],bodyText:document.body.innerText,theme:document.documentElement.dataset.theme,settings:{theme:document.querySelector('#theme').value,tts_enabled:document.querySelector('#read-aloud').checked,tts_speed:document.querySelector('#tts-speed').value,speedLabel:document.querySelector('#speed-value').textContent,font_size:document.querySelector('#font-size').value,fontLabel:document.querySelector('#font-value').textContent},state:document.querySelector('#state').textContent,paused:!document.querySelector('#paused').hidden,notice:document.querySelector('#notice').textContent,tree:document.querySelector('#agent-tree').innerText,panes:Object.fromEntries(['.sidebar','.chat-pane','.actions-pane','#chat-scroll','#actions-scroll','#sidebar-divider','#actions-divider','#draft'].map(sel=>[sel,rect(document.querySelector(sel))])),cards:[...document.querySelectorAll('article')].map(n=>({id:n.dataset.partId,role:n.dataset.role,class:n.className,text:n.innerText})),speechEvents:window.speechEvents}})()""")
def patch(b,s,settings):
 return b.eval(s,"""new Promise((resolve,reject)=>{const w=new WebSocket(location.origin.replace('http','ws')+'/ws');let timer=setTimeout(()=>{w.close();reject(Error('settings ack timeout'))},15000);w.onopen=()=>{w.send(JSON.stringify({type:'subscribe'}));w.send(JSON.stringify({type:'update_settings',settings:"""+json.dumps(settings)+"""}));};w.onmessage=e=>{const m=JSON.parse(e.data);if(m.type==='settings_changed'){clearTimeout(timer);w.close();resolve(m)}else if(m.type==='settings_error'){clearTimeout(timer);w.close();reject(Error(m.text))}}})""")

def local_native(gui,meta,result):
 b=None;native=None
 def rec(name,data):result['observations'].append({'name':name,'data':data});save(meta/'result.json',result);print('OBSERVED '+name,flush=True)
 try:
  # The launcher reserved browser-profile; Chrome receives a fresh sibling, never a shared profile.
  b=Browser(gui.d/'browser-profile-visual',True,meta,'visual');s=b.page(gui.url);s2=b.page(gui.url,True)
  initial=dom(b,s);rec('initial actual rendered DOM',initial);assert not initial['settings']['tts_enabled']
  b.click(s,'#agents-tab');b.key(s,'ArrowRight');b.wait(s,"document.querySelector('#settings-tab').getAttribute('aria-selected')==='true'")
  rec('keyboard tab order',b.eval(s,"({active:document.activeElement.id,settingsVisible:!document.querySelector('#settings-panel').hidden})"))
  old=b.eval(s,"document.querySelector('#sidebar-divider').getAttribute('aria-valuenow')");b.click(s,'#sidebar-divider');b.key(s,'ArrowRight');new=b.eval(s,"document.querySelector('#sidebar-divider').getAttribute('aria-valuenow')");assert float(new)>float(old)
  r=b.eval(s,"(()=>{let r=document.querySelector('#actions-divider').getBoundingClientRect();return{x:r.x+r.width/2,y:r.y+80}})()")
  before=b.eval(s,"document.querySelector('#actions-divider').getAttribute('aria-valuenow')")
  b.call('Input.dispatchMouseEvent',{'type':'mousePressed','button':'left','buttons':1,'clickCount':1,**r},s);b.call('Input.dispatchMouseEvent',{'type':'mouseMoved','button':'left','buttons':1,'x':r['x']-30,'y':r['y']},s);b.call('Input.dispatchMouseEvent',{'type':'mouseReleased','button':'left','buttons':0,'clickCount':1,'x':r['x']-30,'y':r['y']},s)
  after=b.eval(s,"document.querySelector('#actions-divider').getAttribute('aria-valuenow')");assert before!=after;rec('actual divider keyboard and pointer',{'sidebar':[old,new],'actions':[before,after]})
  for theme in ['light','dark','system']:
   b.widget(s,'#theme',theme);b.wait(s2,"document.querySelector('#theme').value==="+json.dumps(theme))
   if theme=='system':
    for media in ['light','dark']:
     b.call('Emulation.setEmulatedMedia',{'features':[{'name':'prefers-color-scheme','value':media}]},s);b.wait(s,"document.documentElement.dataset.theme==="+json.dumps(media));rec('system '+media+' emulated',dom(b,s))
   else:rec('theme '+theme,{'dom':dom(b,s),'screenshot':b.shot(s,'theme-'+theme)})
  b.widget(s,'#theme','light');b.widget(s,'#font-size','20');b.widget(s,'#tts-speed','1.2345');b.wait(s2,"document.querySelector('#tts-speed').value==='1.2345' && document.querySelector('#font-size').value==='20'")
  ack=patch(b,s,{'tts_enabled':False,'max_tool_rounds':0});assert ack['settings']['tts_enabled'] is False and ack['settings']['max_tool_rounds']==0
  b.call('Emulation.setDeviceMetricsOverride',{'width':740,'height':820,'deviceScaleFactor':1,'mobile':False},s);narrow=dom(b,s);rec('740px narrow actual layout',{'dom':narrow,'screenshot':b.shot(s,'narrow')})
  assert narrow['panes']['.actions-pane']['x']+narrow['panes']['.actions-pane']['w']<=741
  assert narrow['panes']['#draft']['w']>100
  b.call('Emulation.clearDeviceMetricsOverride',session=s)
  b.draft(s,'tools slow');b.click(s,'#submit');b.wait(s,"document.querySelector('#state').textContent==='in_flight'")
  b.draft(s,'pause overlap');b.wait(s,"!document.querySelector('#paused').hidden");rec('real in-flight draft pause',dom(b,s));b.click(s,'#cancel');b.wait(s,"document.querySelector('#paused').hidden")
  b.wait(s,"document.querySelector('#state').textContent==='idle' && document.querySelectorAll('#actions .tool-result').length===2",90)
  final=dom(b,s);assert not final['notice'];assert len([x for x in final['cards'] if 'tool_call' in x['class']])==2
  normal_shot=b.shot(s,'tool-completion')
  b.call('Emulation.setDeviceMetricsOverride',{'width':1280,'height':320,'deviceScaleFactor':1,'mobile':False},s)
  scroll=b.eval(s,"(()=>{let a=document.querySelector('#actions-scroll'),c=document.querySelector('#chat-scroll');a.scrollTop=a.scrollHeight;c.scrollTop=c.scrollHeight;let before=[a.scrollTop,c.scrollTop];a.scrollTop=0;let one=[a.scrollTop,c.scrollTop];c.scrollTop=0;return {before,afterActions:one,afterChat:[a.scrollTop,c.scrollTop],dimensions:[a.scrollHeight,a.clientHeight,c.scrollHeight,c.clientHeight]}})()")
  assert scroll['before'][0]>0 and scroll['before'][1]>0
  assert scroll['afterActions']==[0,scroll['before'][1]] and scroll['afterChat']==[0,0]
  rec('actual tool completion / pane scrolling',{'dom':final,'scroll':scroll,'scrollProbeViewport':[1280,320],'wire':b.messages(s,'tool_dispatched')+b.messages(s,'tool_finished'),'screenshot':normal_shot,'scrollProbeScreenshot':b.shot(s,'constrained-scroll')})
  b.call('Emulation.clearDeviceMetricsOverride',session=s)
  gui.restart();b.call('Page.navigate',{'url':gui.url},s);b.call('Page.navigate',{'url':gui.url},s2)
  b.wait(s,"document.querySelector('#connection')?.textContent==='Connected'");b.wait(s2,"document.querySelector('#connection')?.textContent==='Connected'")
  b.wait(s,"document.querySelector('#tts-speed').value==='1.2345' && document.querySelector('#font-size').value==='20' && !document.querySelector('#read-aloud').checked")
  snapshots=b.messages(s,'current_settings');assert snapshots[-1]['settings']['max_tool_rounds']==0 and snapshots[-1]['settings']['tts_enabled'] is False
  rec('one owned restart and two-page persistence',{'page1':dom(b,s),'page2':dom(b,s2),'snapshot':snapshots[-1]})
  b.draft(s,'Hello');b.click(s,'#submit');b.wait(s,"document.querySelector('#state').textContent==='idle' && [...document.querySelectorAll('#transcript article.text[data-role=agent]')].some(n=>n.textContent.includes('First phrase. Final phrase.'))",90)
  rec('native fixed phrase prepared muted',dom(b,s));b.close();b=None
  print('MILESTONE LOCAL PASS; starting authorized brief native window',flush=True)
  native=Browser(gui.d/'browser-profile-native',False,meta,'native');ns=native.page(gui.url);native_start=time.monotonic();deadline=native_start+120;count=0
  native.click(ns,'#settings-tab');native.wait(ns,"!!document.querySelector('#transcript article.text[data-role=agent] .speaker')")
  rec('native platform/default qualification',native.eval(ns,"({userAgent:navigator.userAgent,visibility:document.visibilityState,utteranceVoiceSelection:'production does not assign a voice',reportedDefaults:speechSynthesis.getVoices().filter(v=>v.default).map(v=>({name:v.name,lang:v.lang,localService:v.localService,default:v.default})),humanListening:'not reported'})"))
  speaker=native.eval(ns,'''(()=>{let n=[...document.querySelectorAll('#transcript article.text[data-role=agent]')].find(n=>n.innerText.includes('First phrase. Final phrase.'));if(!n)throw Error('fixed phrase artifact absent');return '#transcript article[data-part-id=\"'+CSS.escape(n.dataset.partId)+'\"] .speaker'})()''')
  rec('native fixed phrase actual locator',{'selector':speaker})
  def remaining():return max(.01,deadline-time.monotonic())
  for requested,expected in [(1.75,1.75),(1.2345,1.2345),(5,3)]:
   native.widget(ns,'#tts-speed',requested);native.wait(ns,"document.querySelector('#tts-speed').value==="+json.dumps(str(requested)))
   base=native.eval(ns,'window.speechEvents.length');native.click(ns,speaker);count+=1
   native.wait(ns,f"window.speechEvents.slice({base}).some(e=>e.event==='end'||e.event==='error')",remaining())
   events=native.eval(ns,f'window.speechEvents.slice({base})');rec('native phrase '+str(count),{'requested':requested,'expectedNativeRate':expected,'events':events,'dom':dom(native,ns)})
   rate=next(e for e in events if e['event']=='rate');assert abs(rate['rate']-expected)<.00001 and rate['requested']==requested
   if not any(e['event']=='start' for e in events) or not any(e['event']=='end' for e in events) or any(e['event']=='error' for e in events):raise Failure('genuine native start/end/error failure')
   native.wait(ns,"document.querySelector('#paused').hidden",remaining())
  native.widget(ns,'#tts-speed',1.75);native.wait(ns,"document.querySelector('#tts-speed').value==='1.75'");native.click(ns,'#read-aloud');native.wait(ns,"document.querySelector('#read-aloud').checked")
  base=native.eval(ns,'window.speechEvents.length');native.click(ns,speaker);count+=1
  native.wait(ns,f"window.speechEvents.slice({base}).some(e=>e.event==='start')",remaining())
  native.draft(ns,'held draft');native.click(ns,'#read-aloud');native.wait(ns,"!document.querySelector('#read-aloud').checked && !speechSynthesis.speaking",remaining())
  overlap=dom(native,ns);assert overlap['paused'];native.click(ns,'#cancel');native.wait(ns,"document.querySelector('#paused').hidden",remaining())
  rec('native cancellation / remaining draft pause',{'overlap':overlap,'after':dom(native,ns),'utterances':count,'elapsed_seconds':time.monotonic()-native_start,'human_listening':'not reported; no hearing claim'})
  native.close();native=None;result['local_native']='pass';result['native_nonempty_utterances']=count;print('MILESTONE NATIVE PASS; no human listening claim',flush=True)
 finally:
  if b:b.close()
  if native:native.close()

def paid(gui,meta,result):
 b=None
 try:
  b=Browser(gui.d/'browser-profile-paid',True,meta,'paid');s=b.page(gui.url)
  if time.monotonic()-gui.began>=60:raise Failure('paid setup exceeded 60 seconds; no prompt sent')
  initial=dom(b,s);assert initial['settings']['tts_enabled'] is False
  prompt='Read exercise.txt using read_file once; quote its one line in one short sentence. Do not run commands or write files.'
  b.draft(s,prompt);started=time.monotonic();b.click(s,'#submit')
  while time.monotonic()-started<90:
   state=dom(b,s);dispatched=b.messages(s,'tool_dispatched')
   if any(m.get('name')!='read_file' or m.get('input',{}).get('path')!='exercise.txt' for m in dispatched):raise Failure('model attempted outside requested one-file task')
   if state['notice']:raise Failure('paid notice: '+state['notice'])
   if b.messages(s,'turn_ended') and state['state']=='idle':break
   if gui.p.poll() is not None:raise Failure('paid controller stopped before completed turn')
   time.sleep(.1)
  else:raise Failure('paid observation ceiling reached')
  finished=b.messages(s,'tool_finished');dispatch=b.messages(s,'tool_dispatched');assert len(dispatch)==1 and len(finished)==1
  assert dispatch[0]['call_id']==finished[0]['call_id'] # tool_finished has no required name
  if finished[0].get('is_error'):raise Failure('paid read_file failed')
  assert 'Chapter 8 boundary check.' in finished[0]['result']
  final=dom(b,s);assert any('Chapter 8 boundary check.' in c['text'] for c in final['cards'] if c['role']=='agent' and 'text' in c['class'])
  result['paid']={'status':'pass','elapsed_seconds':time.monotonic()-started,'prompt':prompt,'initial':initial,'dispatch':dispatch,'finished':finished,'turns':b.messages(s,'turn_ended'),'final':final,'screenshot':b.shot(s,'paid-final')}
  save(meta/'result.json',result);print('MILESTONE PAID completed; no retry',flush=True)
 finally:
  if b:b.close()

def main():
 ap=argparse.ArgumentParser();ap.add_argument('--phase',choices=['local-native','paid'],required=True);ap.add_argument('--script-sha256',required=True);ap.add_argument('--run-id',required=True);a=ap.parse_args();assert sha(__file__)==a.script_sha256
 meta=ROOT/'.ch08-local/operational-runs'/('controller-'+a.run_id);meta.mkdir(exist_ok=False);directory=ROOT/'.ch08-local/operational-runs'/a.run_id
 result={'phase':a.phase,'script_sha256':sha(__file__),'argv':sys.argv,'cwd':str(ROOT),'observations':[],'status':'running','screenshots':'captured, not visually inspected by this agent','human_listening':'not reported'};save(meta/'result.json',result);gui=None;start=time.monotonic()
 try:
  gui=GUI('local' if a.phase=='local-native' else 'paid',directory,meta)
  if a.phase=='local-native':local_native(gui,meta,result)
  else:paid(gui,meta,result)
  result['status']='pass'
 except Exception as e:
  result['status']='failed';result['failure']={'type':type(e).__name__,'message':str(e),'traceback':traceback.format_exc()};print('MILESTONE STOP '+type(e).__name__+': '+str(e),flush=True);raise
 finally:
  if gui:gui.stop()
  result['elapsed_seconds']=round(time.monotonic()-start,3);result['owned_cleanup_completed']=True;save(meta/'result.json',result)
if __name__=='__main__':main()
