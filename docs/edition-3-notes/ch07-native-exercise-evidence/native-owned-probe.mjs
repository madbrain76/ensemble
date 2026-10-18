// Bill authorized one bounded audible exercise using the unchanged default voice.
// Voice service locality is unknown; no local-service proof is a chapter requirement.
// Original unexecuted proposal remains preserved in its sealed evidence directory.
import {spawn, execFile as execFileCallback} from 'node:child_process';
import {promisify} from 'node:util';
import {createServer} from 'node:http';
import {createServer as reserveServer} from 'node:net';
import {readFile, writeFile, mkdtemp, mkdir} from 'node:fs/promises';
import {fileURLToPath} from 'node:url';
import {dirname, join} from 'node:path';
import {createHash} from 'node:crypto';
import {isDeepStrictEqual} from 'node:util';
const execFile = promisify(execFileCallback), script = fileURLToPath(import.meta.url), directory = dirname(script);
const argv = process.argv.slice(2), has = key => argv.includes(key);
const option = key => argv.find(arg => arg.startsWith(key + '='))?.slice(key.length + 1);
const allowed = ['--explicit-later-root-authorization', '--headed-separately-authorized', '--audible-separately-authorized'];
if (argv.some(arg => !allowed.includes(arg) && !arg.startsWith('--expected-script-sha256=') && !arg.startsWith('--expected-policy-sha256='))) throw Error('Unknown option; refused');
if (!has('--explicit-later-root-authorization')) throw Error('REFUSED: preparation grants no execution authority');
const headed = has('--headed-separately-authorized'), audible = has('--audible-separately-authorized');
if (audible && !headed) throw Error('Audible output additionally requires separate headed consent');
const sha = bytes => createHash('sha256').update(bytes).digest('hex');
if (sha(await readFile(script)) !== option('--expected-script-sha256')) throw Error('Script identity mismatch');
const policyBytes = await readFile(join(directory, 'native-policy.json'));
if (sha(policyBytes) !== option('--expected-policy-sha256')) throw Error('Policy identity mismatch');
const policy = JSON.parse(policyBytes), delay = ms => new Promise(resolve => setTimeout(resolve, ms));
// Observation only: do not select a voice or replace native speech/callbacks.
async function recordDefaultVoicePolicy({stage,receipt}) { receipt.native.voicePolicy = {stage, selection:'unchanged browser default', serviceLocality:'unknown', userConsent:true}; }
const receipt = {status:'started-under-explicit-later-authorization', script, script_sha256:sha(await readFile(script)), policy_sha256:sha(policyBytes),
  command:[process.execPath, ...process.argv.slice(1)], started:new Date().toISOString(), headed, audible, ownership:[], commands:[], requests:[], assets:[], observations:[],
  native:{attempted:false, listening:'not observed; separately consenting human required', microphone:'never requested'}, cleanup:[]};
const children = []; let run, fixture, browser, session, contextId, controller, fatalReject, deadline, fixtureCount=0;
let fatalError = null, pageReady = false;
const fatalPromise = new Promise((_, reject) => { fatalReject=reject; });
fatalPromise.catch(() => {});
function fatal(error) { if (fatalError) return; fatalError=error instanceof Error ? error : Error(String(error)); controller?.abort(); fatalReject(fatalError); }
function check() { if (fatalError) throw fatalError; }
async function bounded(test, ms, description) { const end=Date.now()+ms; while(Date.now()<end) { check(); const value=await test(); if(value) return value; await delay(25); } throw Error(description+' unavailable within '+ms+'ms; no retry/fallback'); }
async function command(file, args) { check(); const result=await execFile(file,args,{env:{PATH:'/usr/bin:/bin:/usr/sbin:/sbin',HOME:run||'/var/empty',GIT_CONFIG_NOSYSTEM:'1',GIT_CONFIG_GLOBAL:'/dev/null'},timeout:2000,maxBuffer:1024*1024}); return result.stdout; }
async function processFact(pid, executable, requiredArgs=[]) {
  const facts={pid, ps:await command('/bin/ps',['-p',String(pid),'-o','pid=,ppid=,command=']), expectedExecutable:executable, requiredArgs};
  if(!facts.ps.trim().startsWith(String(pid)+' ') || !facts.ps.includes(executable) || requiredArgs.some(a=>!facts.ps.includes(a))) throw Error('Owned PID/command mismatch: '+pid);
  receipt.ownership.push(facts); return facts;
}
function listeningRows(text) {
  let pid; const rows=[]; for(const line of text.trim().split('\n')) { if(line[0]==='p') pid=Number(line.slice(1)); if(line[0]==='n') rows.push({pid,address:line.slice(1)}); } return rows;
}
async function listenerFact(pid, port) {
  if(!Number.isInteger(port)||port<=0||port>65535||port===8088) throw Error('Forbidden/invalid port');
  // Query the actual port, not just a claimed endpoint. No environment or other process command lines are read.
  const raw=await command('/usr/sbin/lsof',['-nP','-iTCP:'+port,'-sTCP:LISTEN','-Fpcn']);
  const rows=listeningRows(raw); if(!rows.length || rows.some(r=>r.pid!==pid || !['127.0.0.1:'+port,'[::1]:'+port].includes(r.address))) throw Error('Listener is not exclusively the exact owned PID on loopback');
  receipt.ownership.push({pid,port,lsof:raw}); return true;
}
async function awaitListener(child, executable, args, port) {
  await processFact(child.pid,executable,args);
  await bounded(async()=>{ if(child.exitCode!==null||child.signalCode!==null) throw Error('Owned child exited before listening');
    try {return await listenerFact(child.pid,port);} catch(e) { if(e.code===1&&!e.stdout) return false; throw e; }
  },5000,'Owned listener');
  await processFact(child.pid,executable,args); // Recheck PID/argv after listener allocation.
}
function launch(file,args) {
  check(); const env={PATH:'/usr/bin:/bin:/usr/sbin:/sbin',HOME:run,TMPDIR:join(run,'tmp'),LLM_API_KEY:'LOCAL_FIXTURE_NOT_A_CREDENTIAL'};
  const child=spawn(file,args,{cwd:run,env,stdio:['pipe','pipe','pipe']});
  const row={file,args,cwd:run,env,pid:child.pid,stdout:'',stderr:''}; receipt.commands.push(row); child.row=row;
  child.on('error',fatal); child.stdout.on('data',b=>row.stdout+=b); child.stderr.on('data',b=>row.stderr+=b);
  child.on('exit',(code,signal)=>{row.exit=code;row.signal=signal;row.exitedAt=new Date().toISOString();});
  children.push(child); return child;
}
class CDP {
  constructor(ws) {this.ws=ws;this.next=0;this.pending=new Map();this.events=[];
    ws.onmessage=event=>{const m=JSON.parse(event.data);if(m.id){const p=this.pending.get(m.id);if(!p)return;clearTimeout(p.timer);this.pending.delete(m.id);m.error?p.reject(Error(JSON.stringify(m.error))):p.resolve(m.result);}
      else {this.events.push(m);this.event?.(m);}};
    ws.onclose=()=>{for(const p of this.pending.values()){clearTimeout(p.timer);p.reject(Error('Owned CDP closed'));}this.pending.clear();};
  }
  static async open(url){const ws=new WebSocket(url);await new Promise((resolve,reject)=>{const t=setTimeout(()=>{ws.close();reject(Error('CDP open deadline'));},2000);ws.onopen=()=>{clearTimeout(t);resolve();};ws.onerror=e=>{clearTimeout(t);reject(e);};});return new CDP(ws);}
  call(method,params={},sessionId,cleanup=false){if(!cleanup)check();return new Promise((resolve,reject)=>{const id=++this.next,timer=setTimeout(()=>{this.pending.delete(id);reject(Error('CDP timeout '+method));},2000);this.pending.set(id,{resolve,reject,timer});this.ws.send(JSON.stringify({id,method,params,...(sessionId?{sessionId}:{})}));});}
  async eval(expression,cleanup=false){const r=await this.call('Runtime.evaluate',{expression,returnByValue:true,awaitPromise:true,userGesture:true},session,cleanup);if(r.exceptionDetails)throw Error(JSON.stringify(r.exceptionDetails));return r.result.value;}
}
async function snapshot(label) {
  const value=await browser.eval(`({events:window.speechEvents.slice(),native:{speaking:speechSynthesis.speaking,pending:speechSynthesis.pending,paused:speechSynthesis.paused},pauseVisible:!document.querySelector('#paused').hidden,draft:document.querySelector('#draft').value})`);
  receipt.observations.push({label,time:Date.now(),...value});return value;
}
function currentError(snapshot,offset) { const error=snapshot.events.slice(offset).find(e=>e.event==='error');if(error)throw Error('Native callback error: '+JSON.stringify(error)); }
async function exercise() {
  const root=policy.source_root, git=async(...args)=>(await command('/usr/bin/git',['--no-pager','-C',root,...args])).trim();
  if(await git('rev-parse','HEAD')!==policy.commit || await git('rev-parse','HEAD^{tree}')!==policy.tree || await git('status','--porcelain=v1'))throw Error('Source identity/cleanliness mismatch');
  for(const [name,expected] of Object.entries(policy.source_files))if(sha(await readFile(join(root,name)))!==expected)throw Error('Source hash mismatch '+name);
  if(sha(await readFile(policy.binary))!==policy.binary_sha256)throw Error('Frozen GUI binary mismatch');
  for(const row of Object.values(policy.assets))if(sha(await readFile(join(root,row.source)))!==row.sha256)throw Error('Expected asset source changed');
  const requestBytes=await readFile(join(directory,policy.expected_request_file));if(sha(requestBytes)!==policy.expected_request_sha256)throw Error('Fixed request policy changed');
  receipt.preflight={commit:policy.commit,tree:policy.tree,sourceCount:Object.keys(policy.source_files).length,binary:policy.binary,binary_sha256:policy.binary_sha256};
  // Explicit user consent covers a default voice of unknown service locality.
  await recordDefaultVoicePolicy({stage:'before-process-launch',policy,receipt});
  const expectedRequest=JSON.parse(requestBytes), parts=policy.request_policy.response_parts;
  fixture=createServer((req,res)=>{
    const row={number:++fixtureCount,method:req.method,path:req.url,bytes:0};receipt.requests.push(row);
    function reject(why){if(!res.headersSent)res.writeHead(400);res.end('fixture refused');fatal(Error(why));}
    if(fixtureCount!==1||req.method!=='POST'||req.url!=='/v1/messages'){reject('Fixed request method/path/count violated');return;}
    const chunks=[];req.setTimeout(2000,()=>reject('Fixture request time bound'));req.on('error',fatal);
    req.on('data',b=>{row.bytes+=b.length;if(row.bytes>65536){reject('Fixture body too large');req.destroy();}else chunks.push(b);});
    req.on('end',()=>{try{const bytes=Buffer.concat(chunks);row.body_sha256=sha(bytes);row.body=JSON.parse(bytes);if(!isDeepStrictEqual(row.body,expectedRequest))throw Error('Exact fixed JSON body mismatch');
      res.writeHead(200,{'content-type':'application/json'});res.end(JSON.stringify({id:'owned-native-fixed',type:'message',role:'assistant',model:'claude-sonnet-4-6',content:parts.map(text=>({type:'text',text})),stop_reason:'end_turn',stop_sequence:null,usage:{input_tokens:1,output_tokens:12}}));
    }catch(e){reject(String(e));}});
  });
  await new Promise((resolve,reject)=>{fixture.once('error',reject);fixture.listen(0,'127.0.0.1',resolve);});
  const fixturePort=fixture.address().port;await processFact(process.pid,process.execPath,[script]);await listenerFact(process.pid,fixturePort);
  const reservation=reserveServer();await new Promise(r=>reservation.listen(0,'127.0.0.1',r));const port=reservation.address().port;await new Promise(r=>reservation.close(r));if(port===8088)throw Error('Forbidden GUI port');
  const origin=`http://127.0.0.1:${port}`;const guiArgs=['--port',String(port),'--vendor','anthropic','--model','claude-sonnet-4-6','--base-url',`http://127.0.0.1:${fixturePort}`,'--workspace',run,'--no-stream'];
  const gui=launch(policy.binary,guiArgs);await awaitListener(gui,policy.binary,guiArgs,port);
  // Source/binary/expected assets were verified before GUI birth. Actual served bytes are verified
  // after its owned listener exists, but BEFORE Chrome launch or page/native activity.
  for(const [path,row] of Object.entries(policy.assets)){await listenerFact(gui.pid,port);const response=await fetch(origin+path,{signal:controller.signal,redirect:'error'});const bytes=Buffer.from(await response.arrayBuffer());if(response.status!==200||sha(bytes)!==row.sha256||bytes.length!==row.bytes)throw Error('Served asset mismatch '+path);receipt.assets.push({path,sha256:sha(bytes),bytes:bytes.length});}
  const profile=await mkdtemp(join(run,'chrome-profile-')),chrome='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';
  const chromeArgs=['--user-data-dir='+profile,'--remote-debugging-port=0','--remote-debugging-address=127.0.0.1','--no-startup-window','--no-first-run','--no-default-browser-check','--disable-background-networking','--disable-sync','--disable-extensions','--password-store=basic','--use-mock-keychain','--host-resolver-rules=MAP * ~NOTFOUND','--disable-component-update'];
  if(!headed)chromeArgs.push('--headless=new');if(!audible)chromeArgs.push('--mute-audio');const chromeChild=launch(chrome,chromeArgs);
  await processFact(chromeChild.pid,chrome,['--user-data-dir='+profile,'--remote-debugging-port=0']);
  let endpoint;await bounded(async()=>{try{endpoint=(await readFile(join(profile,'DevToolsActivePort'),'utf8')).trim().split('\n');return endpoint.length===2;}catch(e){if(e.code==='ENOENT')return false;throw e;}},5000,'Owned CDP endpoint file');
  const cdpPort=Number(endpoint[0]);if(!/^\/devtools\/browser\/[a-zA-Z0-9-]+$/.test(endpoint[1]))throw Error('Invalid newly owned CDP endpoint');
  await awaitListener(chromeChild,chrome,['--user-data-dir='+profile,'--remote-debugging-port=0'],cdpPort);
  browser=await CDP.open(`ws://127.0.0.1:${cdpPort}${endpoint[1]}`);receipt.browser=await browser.call('Browser.getVersion');
  contextId=(await browser.call('Target.createBrowserContext',{disposeOnDetach:true})).browserContextId;
  const {targetId}=await browser.call('Target.createTarget',{url:'about:blank',background:true,browserContextId:contextId});receipt.target={targetId,contextId,profile,cdpPort,pid:chromeChild.pid};
  session=(await browser.call('Target.attachToTarget',{targetId,flatten:true})).sessionId;
  await browser.call('Page.enable',{},session);await browser.call('Runtime.enable',{},session);await browser.call('Network.enable',{},session);await browser.call('Fetch.enable',{patterns:[{urlPattern:'*'}]},session);
  const pageRequests=new Map();let pageSockets=0;
  browser.event=m=>{if(m.sessionId!==session)return;
    if(m.method==='Fetch.requestPaused'){const r=m.params.request,u=new URL(r.url),count=(pageRequests.get(u.pathname)||0)+1;pageRequests.set(u.pathname,count);const valid=count===1&&r.method==='GET'&&u.origin===origin&&!u.search&&(Object.hasOwn(policy.assets,u.pathname)||u.pathname==='/favicon.ico');
      browser.call(valid?'Fetch.continueRequest':'Fetch.failRequest',{requestId:m.params.requestId,...(valid?{}:{errorReason:'BlockedByClient'})},session).catch(fatal);if(!valid)fatal(Error('Page network policy violated: '+r.method+' '+r.url));}
    if(m.method==='Network.webSocketCreated'&&(++pageSockets!==1||m.params.url!==`ws://127.0.0.1:${port}/ws`))fatal(Error('Repeated or non-owned page WebSocket attempted'));
    if(m.method==='Runtime.exceptionThrown')fatal(Error('Actual page exception '+JSON.stringify(m.params)));
  };
  // No Page.bringToFront, Target.activateTarget, focus, microphone or recognition calls.
  await browser.call('Page.navigate',{url:origin},session);
  await bounded(()=>browser.eval("document.querySelector('#connection')?.textContent==='Connected' && !!window.speechEvents"),3000,'Production page ready');pageReady=true;
  await browser.eval("document.querySelector('#read-aloud').checked=false;document.querySelector('#read-aloud').dispatchEvent(new Event('change'));true");
  await recordDefaultVoicePolicy({stage:'actual-unmodified-page-before-prompt-or-speech',browser,session,policy,receipt});
  await browser.eval(`document.querySelector('#draft').value=${JSON.stringify(expectedRequest.messages[0].content[0].text)};document.querySelector('#draft').dispatchEvent(new Event('input'));document.querySelector('#composer').requestSubmit();true`);
  await bounded(()=>browser.eval(`document.querySelector('#state').textContent==='idle' && ${JSON.stringify(parts)}.every(t=>Array.from(document.querySelectorAll('article[data-role="agent"]')).some(a=>a.innerText.includes(t)&&!a.classList.contains('streaming')))`),3000,'Two exact fixed assistant artifacts');
  if(fixtureCount!==1)throw Error('Fixture count not exactly one');
  const selector=text=>`Array.from(document.querySelectorAll('article[data-role="agent"]')).find(a=>a.innerText.includes(${JSON.stringify(text)})).querySelector('.speaker').click()`;
  const initial=await snapshot('before-queue'),offset=initial.events.length;receipt.native.attempted=true;
  await browser.eval(`${selector(parts[0])};${selector(parts[1])};true`);
  await bounded(async()=>{const s=await snapshot('waiting-first-start');currentError(s,offset);if(s.events.slice(offset).some(e=>e.event==='start'&&e.text===parts[1]))throw Error('Queued second began before intended cancellation');return s.events.slice(offset).some(e=>e.event==='start'&&e.text===parts[0])&&s.native.speaking&&s.pauseVisible;},3000,'First actual native start + global pause');
  const atStart=await snapshot('first-start');if(!parts.every(t=>atStart.events.slice(offset).some(e=>e.event==='queued'&&e.text===t)))throw Error('Both samples were not actually queued');
  // Real existing page control: production SpeechQueue.cancel calls native synth.cancel.
  // Canceled callbacks are intentionally generation-guarded, so no invented oncancel/error assertion.
  await browser.eval("document.querySelector('#read-aloud').checked=false;document.querySelector('#read-aloud').dispatchEvent(new Event('change'));true");
  await bounded(async()=>{const s=await snapshot('cancellation-release');currentError(s,offset);return s.events.slice(offset).some(e=>e.event==='cancel')&&!s.native.speaking&&!s.native.pending&&!s.pauseVisible;},2000,'Actual cancellation and pause release');
  await delay(300);const quiet=await snapshot('post-cancel-quiet');currentError(quiet,offset);if(quiet.native.speaking||quiet.native.pending||quiet.events.slice(offset).some(e=>e.event==='start'&&e.text===parts[1]))throw Error('Canceled queued utterance leaked');
  const recoveryOffset=quiet.events.length;await browser.eval(`${selector(parts[1])};true`);
  await bounded(async()=>{const s=await snapshot('recovery-start');currentError(s,recoveryOffset);return s.events.slice(recoveryOffset).some(e=>e.event==='start'&&e.text===parts[1])&&s.pauseVisible;},3000,'Subsequent queued actual native start');
  await bounded(async()=>{const s=await snapshot('recovery-end');currentError(s,recoveryOffset);return s.events.slice(recoveryOffset).some(e=>e.event==='end'&&e.text===parts[1])&&!s.native.speaking&&!s.native.pending&&!s.pauseVisible;},6000,'Subsequent actual native END and pause release');
  receipt.native.completedActualSequence=true;receipt.status='bounded-native-observations-complete-not-human-listening';
}
try {
  run=await mkdtemp('/tmp/cadence-ch07-owned-native-');await mkdir(join(run,'tmp'));receipt.run=run;controller=new AbortController();
  deadline=setTimeout(()=>fatal(Error('55-second exercise deadline; no repeat')),55000);
  await Promise.race([exercise(),fatalPromise]);
}catch(e){receipt.status='blocked-or-failed';receipt.error=String(e);process.exitCode=1;}
finally {
  clearTimeout(deadline);controller?.abort();const cleanupEnd=Date.now()+10000;
  if(browser&&pageReady){try{await browser.eval("document.querySelector('#read-aloud').checked=false;document.querySelector('#read-aloud').dispatchEvent(new Event('change'));true",true);}catch(e){receipt.cleanup.push({pageCancelError:String(e)});}}
  if(browser){receipt.cdpEvents=browser.events;try{if(contextId)await browser.call('Target.disposeBrowserContext',{browserContextId:contextId},undefined,true);}catch(e){receipt.cleanup.push({contextCloseError:String(e)});}try{await browser.call('Browser.close',{},undefined,true);}catch(e){receipt.cleanup.push({browserCloseError:String(e)});}browser.ws.close();}
  for(const child of children.reverse()){
    if(child.exitCode===null&&child.signalCode===null)child.kill('SIGTERM');
    const stop=Math.min(Date.now()+1000,cleanupEnd);while(child.exitCode===null&&child.signalCode===null&&Date.now()<stop)await delay(20);
    if(child.exitCode===null&&child.signalCode===null)child.kill('SIGKILL');
    while(child.exitCode===null&&child.signalCode===null&&Date.now()<cleanupEnd)await delay(20);
    const exited=child.exitCode!==null||child.signalCode!==null;receipt.cleanup.push({pid:child.pid,actualExitObserved:exited,exit:child.exitCode,signal:child.signalCode});
    if(!exited){receipt.status='cleanup-incomplete-stop-for-root';process.exitCode=1;child.unref();child.stdin.destroy();child.stdout.destroy();child.stderr.destroy();}
  }
  if(fixture){fixture.closeAllConnections();const closed=await Promise.race([new Promise(resolve=>fixture.close(()=>resolve(true))),delay(Math.max(0,cleanupEnd-Date.now())).then(()=>false)]);receipt.cleanup.push({fixtureClosed:closed&&!fixture.listening});if(!closed){receipt.status='cleanup-incomplete-stop-for-root';process.exitCode=1;fixture.unref();}}
  receipt.finished=new Date().toISOString();receipt.profileDisposition='Owned scratch/profile retained, never shared/reused; no historical artifacts/cache deleted';
  if(run)await writeFile(join(run,'partial-or-final-receipt.json'),JSON.stringify(receipt,null,2)+'\n');
  console.log(run?join(run,'partial-or-final-receipt.json'):'No run directory created');
}
