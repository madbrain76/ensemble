#!/usr/bin/env python3
"""Only CDP ownership/mechanics reused from the explicitly permitted controller-v5 scaffold.
No Chapter 8 scenario or authority is imported. Constructed only after root approval.
"""
import base64,json,os,subprocess,time,websocket
from pathlib import Path
from launcher import Failure,require,sha,save,env,stop_process,port_closed
CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'
class Browser:
 def __init__(self,directory,muted,meta,label,mock_keychain=True):
  self.meta=meta;self.label=label;self.trace=(meta/(label+'-raw-cdp.jsonl')).open('w');self.events=[];self.wire=[];self.ids=0;self.ws=None
  directory.mkdir(exist_ok=False);home=meta/(label+'-home');home.mkdir();log=(meta/(label+'-chrome.log')).open('w');self.log=log
  args=[CHROME,'--user-data-dir='+str(directory),'--remote-debugging-address=127.0.0.1','--remote-debugging-port=0','--no-first-run','--no-default-browser-check','--disable-background-networking','--disable-component-update','--disable-sync','--disable-extensions','--password-store=basic','--window-size=1280,900']
  if mock_keychain:args+=['--use-mock-keychain']
  if muted:args+=['--mute-audio']
  args+=['about:blank'];self.p=subprocess.Popen(args,env=env(home),stdout=log,stderr=subprocess.STDOUT,start_new_session=True)
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
  p=getattr(self,'p',None);ws=getattr(self,'ws',None)
  if p and p.poll() is None and ws:
   try:self.call('Browser.close',timeout=3)
   except Exception:pass
  code=stop_process(p)
  if ws:ws.close()
  for name in ('log','trace'):
   f=getattr(self,name,None)
   if f:f.close()
  own=getattr(self,'own',{});listener=own.get('cdp_listener');closed=not listener or port_closed(listener[1])
  save(self.meta/(self.label+'-browser-exit.json'),{'pid':p.pid if p else None,'exit_code':code,'cdp_listener':listener,'listener_closed':closed,'profile_preserved':True})
  require(closed,'owned CDP listener survived')
 def messages(self,s,kind):return [r['message'] for r in self.wire if r['session']==s and r['message'].get('type')==kind]
