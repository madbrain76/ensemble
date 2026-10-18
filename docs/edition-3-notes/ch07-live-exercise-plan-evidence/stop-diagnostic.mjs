// One Chapter 7 exercise, not a product module. Default: no execution.
// Local controls use only loopback fake upstreams; live/key/browser require later explicit approval.
import http from 'node:http';
import https from 'node:https';
import {spawn,execFile} from 'node:child_process';
import {promisify} from 'node:util';
import {readFile,writeFile,mkdir,mkdtemp,realpath,lstat} from 'node:fs/promises';
import {createWriteStream} from 'node:fs';
import {createHash} from 'node:crypto';
import {resolve,dirname,basename,join} from 'node:path';
import {fileURLToPath} from 'node:url';
import assert from 'node:assert/strict';
const exec=promisify(execFile),here=dirname(fileURLToPath(import.meta.url));
const sha=b=>createHash('sha256').update(b).digest('hex');
const canonical=x=>JSON.stringify(x&&typeof x==='object'?(Array.isArray(x)?x.map(y=>JSON.parse(canonical(y))):Object.fromEntries(Object.keys(x).sort().map(k=>[k,JSON.parse(canonical(x[k]))]))):x);
const sleep=ms=>new Promise(r=>setTimeout(r,ms));
const closedEnv=home=>({PATH:'/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin',HOME:home,TMPDIR:home});
const fail=code=>Object.assign(new Error(code),{safeCode:code});
const demand=(ok,code)=>{if(!ok)throw fail(code)};
const stats={credentialReads:0,externalAttempts:0,browserLaunches:0,guiLaunches:0};
let policy,run,activeChildren=[],activeServers=[],activeSockets=[];
async function pins(p=policy){
 demand((await exec('/usr/bin/git',['--no-pager','rev-parse','HEAD'],{cwd:p.source_root,env:{...closedEnv(run),GIT_CONFIG_NOSYSTEM:'1',GIT_CONFIG_GLOBAL:'/dev/null'}})).stdout.trim()===p.commit,'PIN_COMMIT');
 demand(!(await exec('/usr/bin/git',['--no-pager','status','--porcelain=v1'],{cwd:p.source_root,env:{...closedEnv(run),GIT_CONFIG_NOSYSTEM:'1',GIT_CONFIG_GLOBAL:'/dev/null'}})).stdout.trim(),'PIN_DIRTY');
 for(const [path,hash] of Object.entries(p.source_files))demand(sha(await readFile(join(p.source_root,path)))===hash,'PIN_SOURCE');
 demand(sha(await readFile(p.gui.path))===p.gui.sha256,'PIN_GUI');
 for(const row of Object.values(p.assets))demand(sha(await readFile(join(p.source_root,row.source)))===row.sha256,'PIN_ASSET');
}
async function listen(server){await new Promise((r,j)=>{server.once('error',j);server.listen(0,'127.0.0.1',r)});activeServers.push(server);const port=server.address().port;demand(port!==8088,'PORT_8088');return `http://127.0.0.1:${port}`}
async function ownership(child,binary,args,port){
 demand(child.exitCode===null&&child.signalCode===null&&child.pid>0,'OWNED_CHILD_EXITED');
 const env=closedEnv(run),command=(await exec('/bin/ps',['-ww','-p',String(child.pid),'-o','pid=,ppid=,command='],{env})).stdout.trim();
 const match=command.match(/^(\d+)\s+(\d+)\s+([\s\S]*)$/);demand(match&&+match[1]===child.pid&&+match[2]===process.pid,'PID_PARENT');
 // Resolve executable symlinks before spawn. Match whole known executable prefix, not whitespace tokenization of a macOS app path.
 demand(match[3].startsWith(binary+' ')||match[3]===binary,'PID_COMMAND');
 for(const arg of args.filter(x=>x.startsWith('--user-data-dir=')||x.startsWith('--remote-debugging-port=')))demand(match[3].includes(arg),'PID_PROFILE');
 const listeners=(await exec('/usr/sbin/lsof',['-nP',`-iTCP:${port}`,'-sTCP:LISTEN','-Fpn'],{env})).stdout;
 const ids=[...new Set([...listeners.matchAll(/^p(\d+)$/gm)].map(x=>+x[1]))];demand(ids.length===1&&ids[0]===child.pid,'LISTENER_OWNER');
 demand(listeners.split('\n').filter(x=>x.startsWith('n')).every(x=>x===`n127.0.0.1:${port}`||x===`n[::1]:${port}`),'LISTENER_LOOPBACK');
 return {pid:child.pid,ppid:process.pid,binary,command,port,listeners};
}
async function waitUntil(fn,ms,code){const end=performance.now()+ms;while(performance.now()<end){const v=await fn();if(v)return v;await sleep(25)}throw fail(code)}
async function startGUI(proxy,dir,p=policy){
 await pins(p);await mkdir(dir,{recursive:true,mode:0o700});await writeFile(join(dir,p.fixture.name),p.fixture.text,{mode:0o600});
 const binary=await realpath(p.gui.path),args=['--port','0','--vendor',p.vendor,'--model',p.model,'--base-url',proxy.origin,'--workspace',dir,'--event-log',join(dir,'event.log')];
 const child=spawn(binary,args,{cwd:dir,env:{...closedEnv(dir),LLM_API_KEY:'LOCAL_GUI_DUMMY_ONLY'},stdio:['pipe','pipe','pipe']});stats.guiLaunches++;activeChildren.push(child);
 child.exited=new Promise(r=>child.once('exit',(code,signal)=>r({code,signal})));let stderr='';child.stderr.on('data',b=>stderr+=b.toString());
 child.stdout.pipe(createWriteStream(join(dir,'stdout.log'),{mode:0o600}));child.stderr.pipe(createWriteStream(join(dir,'stderr.log'),{mode:0o600}));
 child.once('error',()=>{child.launchError=true});
 const port=await waitUntil(()=>{demand(!child.launchError&&child.exitCode===null,'GUI_START');return Number(stderr.match(/GUI http:\/\/127\.0\.0\.1:(\d+)/)?.[1])},6000,'GUI_LISTEN_TIMEOUT');demand(port!==8088,'PORT_8088');
 const own=await ownership(child,binary,args,port),origin=`http://127.0.0.1:${port}`;
 const assets={};for(const [path,row] of Object.entries(p.assets)){const res=await fetch(origin+path,{redirect:'error',signal:AbortSignal.timeout(3000)});demand(res.ok,'SERVED_ASSET_STATUS');const b=Buffer.from(await res.arrayBuffer());demand(sha(b)===row.sha256,'SERVED_ASSET_PIN');assets[path]={sha256:sha(b),bytes:b.length}}
 return {child,origin,ownership:own,assets,args,dir};
}
class StreamGate{
 constructor(state,res){this.s=state;this.res=res;this.buffer=Buffer.alloc(0);this.bytes=0;this.terminal=false;this.pending=null;this.usage={};}
 push(chunk){this.bytes+=chunk.length;demand(this.bytes<=policy.request.maximum_response_bytes,'RESPONSE_BYTES');this.buffer=Buffer.concat([this.buffer,chunk]);for(;;){let n=this.buffer.indexOf('\n\n'),end=2;const c=this.buffer.indexOf('\r\n\r\n');if(c>=0&&(n<0||c<n)){n=c;end=4}if(n<0){demand(this.buffer.length<=65536,'SSE_FRAME_BYTES');break}const raw=this.buffer.subarray(0,n+end);this.buffer=this.buffer.subarray(n+end);this.frame(raw)}}
 frame(raw){const text=new TextDecoder('utf-8',{fatal:true}).decode(raw);const data=text.split(/\r?\n/).filter(l=>l.startsWith('data:')).map(l=>l.slice(5).replace(/^ /,'')).join('\n');if(!data){if(this.pending)this.pending.frames.push(raw);else this.res.write(raw);return}let m;try{m=JSON.parse(data)}catch{throw fail('SSE_JSON')}
 demand(!this.terminal,'SSE_AFTER_TERMINAL');demand(m.type!=='error'&&!m.error,'PROVIDER_STREAM_ERROR');
 if(m.type==='message_start'){if(m.message?.model)demand(m.message.model===policy.model,'RESPONSE_MODEL');this.usage.input=m.message?.usage?.input_tokens;}
 if(m.type==='message_delta'){this.usage.output=m.usage?.output_tokens;this.stopReason=m.delta?.stop_reason;}
 if(m.type==='content_block_start'&&m.content_block?.type==='tool_use'){
  demand(!this.pending&&m.content_block.name==='read_file'&&this.s.readCalls<policy.fixture.maximum_read_calls,'TOOL_SCOPE');
  this.pending={index:m.index,initial:m.content_block.input??{},args:'',frames:[raw]};return;
 }
 if(this.pending){const t=this.pending;demand(m.index===t.index,'TOOL_BLOCK_ORDER');t.frames.push(raw);
  if(m.type==='content_block_delta'){demand(m.delta?.type==='input_json_delta','TOOL_DELTA');t.args+=m.delta.partial_json;demand(t.args.length<=4096,'TOOL_ARGUMENT_BYTES');return}
  demand(m.type==='content_block_stop','TOOL_BLOCK_END');let args;try{args=t.args?JSON.parse(t.args):t.initial}catch{throw fail('TOOL_ARGUMENT_JSON')}
  demand(args&&typeof args==='object'&&!Array.isArray(args)&&typeof args.path==='string','TOOL_ARGUMENTS');
  demand(resolve(this.s.workspace,args.path)===join(this.s.workspace,policy.fixture.name),'TOOL_PATH');
  demand(Object.keys(args).every(k=>['path','start_line','end_line','max_bytes'].includes(k)),'TOOL_ARGUMENT_KEYS');
  for(const k of ['start_line','end_line','max_bytes'])if(k in args)demand(Number.isSafeInteger(args[k])&&args[k]>=0,'TOOL_ARGUMENT_LIMIT');
  this.s.readCalls++;for(const frame of t.frames)this.res.write(frame);this.pending=null;return;
 }
 if(m.type==='content_block_start')demand(['text','thinking','redacted_thinking'].includes(m.content_block?.type),'RESPONSE_BLOCK_SCOPE');
 if(m.type==='message_stop'){demand(!this.pending,'TOOL_INCOMPLETE');this.terminal=true;}
 this.res.write(raw);
 }
 finish(){demand(!this.pending&&this.terminal&&!this.buffer.toString().trim(),'INCOMPLETE_STREAM');}
}
async function makeProxy(target,workspace,key='LOCAL_UPSTREAM_DUMMY_ONLY'){
 const u=new URL(target.url);demand(target.mode==='local'?u.protocol==='http:'&&u.hostname==='127.0.0.1':target.mode==='live'&&target.approved===true&&u.href===policy.endpoint,'TARGET_DISABLED');
 const s={mode:target.mode,workspace,attempts:0,reserved:0,completed:0,rejected:[],readCalls:0,seen:new Set(),records:[],requests:new Set(),halted:null,deadline:null};
 function reject(res,code,status=400){s.halted??=code;s.rejected.push({at:Date.now(),code});res.writeHead(status,{'content-type':'application/json'});res.end(JSON.stringify({error:code}));}
 const server=http.createServer((req,res)=>{
  s.attempts++;if(s.halted)return reject(res,'HALTED');if(req.method!==policy.request.method||req.url!==policy.request.route)return reject(res,'ROUTE');
  if(req.headers['x-api-key']!=='LOCAL_GUI_DUMMY_ONLY'||req.headers['anthropic-version']!==policy.request.anthropic_version)return reject(res,'LOCAL_AUTH_OR_VERSION');
  let size=0,chunks=[];req.setTimeout(3000,()=>{req.destroy();s.halted='REQUEST_BODY_TIMEOUT'});
  req.on('data',b=>{size+=b.length;if(size>policy.request.maximum_request_bytes){req.destroy();s.halted='REQUEST_BYTES'}else chunks.push(b)});
  req.on('error',()=>{});
  req.on('end',()=>{
   if(s.halted)return reject(res,'HALTED');
   const raw=Buffer.concat(chunks);let b;try{b=JSON.parse(raw)}catch{return reject(res,'REQUEST_JSON')}
   if(b.model!==policy.model)return reject(res,'MODEL');if(b.max_tokens!==4096)return reject(res,'TOKENS');if(b.stream!==true)return reject(res,'STREAM');
   if(canonical(Object.keys(b).sort())!==canonical(policy.request.top_keys)||b.system!==policy.request.system||sha(canonical(b.tools))!==policy.request.tools_sha256||!Array.isArray(b.messages))return reject(res,'REQUEST_BOUNDARY');
   const digest=sha(canonical(b));if(s.seen.has(digest))return reject(res,'DUPLICATE_NO_RETRY');
   if(s.deadline&&performance.now()>=s.deadline)return reject(res,'DEADLINE');
   // Atomic reservation in this one Node event loop. No await between check/increment/network construction.
   if(s.reserved>=policy.request.maximum_requests)return reject(res,'REQUEST_BUDGET',429);
   s.reserved++;s.seen.add(digest);s.deadline??=performance.now()+policy.request.observation_ms;
   const row={number:s.reserved,at:Date.now(),request_sha256:sha(raw),request_bytes:raw.length,model:b.model,max_tokens:b.max_tokens,stream:b.stream};s.records.push(row);
   if(target.mode==='live')stats.externalAttempts++;
   const transport=target.mode==='local'?http:https;
   const upstream=transport.request(u,{method:'POST',headers:{'content-type':'application/json','content-length':raw.length,'x-api-key':key,'anthropic-version':policy.request.anthropic_version,'accept-encoding':'identity'},rejectUnauthorized:true},reply=>{
    row.status=reply.statusCode;
    if(reply.statusCode!==200||!String(reply.headers['content-type']).startsWith('text/event-stream')){s.halted=`UPSTREAM_HTTP_${reply.statusCode}`;reply.destroy();return reject(res,s.halted,502)}
    res.writeHead(200,{'content-type':'text/event-stream','cache-control':'no-store'});res.flushHeaders();const gate=new StreamGate(s,res);
    reply.on('data',chunk=>{try{gate.push(chunk)}catch(err){(row.transitions??=[]).push({stage:'data',code:err.safeCode||'STREAM_INVALID'});s.halted=err.safeCode||'STREAM_INVALID';row.failure=s.halted;reply.destroy();res.destroy()}});
    reply.on('end',()=>{try{gate.finish();row.usage=gate.usage;row.stop_reason=gate.stopReason;row.response_bytes=gate.bytes;s.completed++;res.end()}catch(err){(row.transitions??=[]).push({stage:'end',code:err.safeCode||'STREAM_INVALID'});s.halted=err.safeCode||'STREAM_INVALID';row.failure=s.halted;res.destroy()}});
    reply.on('error',()=>{s.halted??='UPSTREAM_READ';row.failure??=s.halted;res.destroy()});
   });s.requests.add(upstream);
   const timer=setTimeout(()=>{s.halted='DEADLINE';upstream.destroy();res.destroy()},Math.max(1,s.deadline-performance.now()));
   upstream.once('close',()=>{clearTimeout(timer);s.requests.delete(upstream)});
   upstream.on('error',()=>{s.halted??='UPSTREAM_CONNECTION';row.failure??=s.halted;if(!res.headersSent)reject(res,s.halted,502);else res.destroy()});
   upstream.end(raw); // Accepted body bytes unchanged. No redirect handling, retry, model substitution, or request transformation.
  });
 });const origin=await listen(server);
 return {origin,state:s,stop(){s.halted??='STOPPED';for(const q of s.requests)q.destroy();server.closeAllConnections()},summary(){return {...s,seen:undefined,requests:undefined,workspace:undefined}}};
}
async function fakeProvider(mode){
 const records=[];const server=http.createServer(async(req,res)=>{let chunks=[];for await(const c of req)chunks.push(c);const raw=Buffer.concat(chunks),b=JSON.parse(raw);records.push({body:b,sha256:sha(raw),path:req.url,max_tokens:b.max_tokens,stream:b.stream});
  res.on('error',()=>{});res.writeHead(200,{'content-type':'text/event-stream'});
  const emit=m=>res.write('data: '+JSON.stringify(m)+'\n\n');
  emit({type:'message_start',message:{model:policy.model,usage:{input_tokens:10,output_tokens:0}}});
  const result=b.messages.some(m=>Array.isArray(m.content)&&m.content.some(p=>p.type==='tool_result'));
  const tool=!result&&mode!=='budget';const text=tool?'I will read the provided note.':'Lantern Orchard needs three reviewers; use the blue folder.';
  emit({type:'content_block_start',index:0,content_block:{type:'text',text:''}});
  for(const t of [text.slice(0,12),text.slice(12)]){emit({type:'content_block_delta',index:0,delta:{type:'text_delta',text:t}});await sleep(50)}
  emit({type:'content_block_stop',index:0});
  if(tool){emit({type:'content_block_start',index:1,content_block:{type:'tool_use',id:'owned-read-1',name:mode==='unsafe'?'run_command':'read_file',input:{}}});emit({type:'content_block_delta',index:1,delta:{type:'input_json_delta',partial_json:JSON.stringify(mode==='unsafe'?{command:'printf unsafe > SHOULD_NOT_EXIST'}:{path:policy.fixture.name})}});emit({type:'content_block_stop',index:1});}
  emit({type:'message_delta',delta:{stop_reason:tool?'tool_use':'end_turn'},usage:{output_tokens:20}});emit({type:'message_stop'});res.end();
 });return {origin:await listen(server),records,server};
}
async function wireTurn(gui,proxy){const ws=new WebSocket(gui.origin.replace('http:','ws:')+'/ws');activeSockets.push(ws);const frames=[];let ended=false,connected=false;ws.onmessage=e=>{const m=JSON.parse(e.data);frames.push({at:Date.now(),...m});if(m.type==='turn_ended')ended=true};ws.onerror=()=>{};await new Promise((r,j)=>{ws.onopen=()=>{connected=true;ws.send(JSON.stringify({type:'subscribe'}));ws.send(JSON.stringify({type:'prompt',text:policy.fixture.prompt}));r()};ws.onerror=j});await waitUntil(()=>ended||proxy.state.halted,10000,'LOCAL_TURN_TIMEOUT');await sleep(100);ws.close();return {frames,connected,ended};}
async function post(origin,body){return new Promise((resolve,reject)=>{const b=Buffer.from(JSON.stringify(body)),q=http.request(origin+'/v1/messages',{method:'POST',headers:{'content-type':'application/json','content-length':b.length,'x-api-key':'LOCAL_GUI_DUMMY_ONLY','anthropic-version':policy.request.anthropic_version}},res=>{res.resume();res.on('end',()=>resolve(res.statusCode))});q.on('error',reject);q.end(b)})}
async function stopChild(child){if(child.exitCode!==null||child.signalCode!==null)return {code:child.exitCode,signal:child.signalCode};child.kill('SIGTERM');const stopped=await Promise.race([child.exited,sleep(2500).then(()=>null)]);if(stopped)return stopped;child.kill('SIGKILL');const killed=await Promise.race([child.exited,sleep(2500).then(()=>null)]);demand(killed,'OWNED_CHILD_DID_NOT_EXIT');return killed;}
async function localControls(){
 const result={mode:'local controls; NOT browser/native/live-provider evidence',checks:[],failures:[],stats,source:policy.commit};let base;
 try{
  const uf=await fakeProvider('unsafe'),ud=join(run,'unsafe'),up=await makeProxy({mode:'local',url:uf.origin+'/v1/messages'},ud),ug=await startGUI(up,ud),uw=await wireTurn(ug,up);
  let exists=true;try{await lstat(join(ud,'SHOULD_NOT_EXIST'))}catch(e){if(e.code==='ENOENT')exists=false;else throw e}
  result.stopDiagnostic={proxy:up.summary(),wireTypes:uw.frames.map(m=>m.type),forbiddenEffectExists:exists};
  demand(up.state.halted==='TOOL_SCOPE'&&!exists&&!uw.frames.some(m=>m.type==='tool_dispatched'),'UNSAFE_TOOL_STOP');await stopChild(ug.child);up.stop();result.unsafe={reason:up.state.halted,upstream_effects:uf.records.length,tool_starts:0,forbidden_effect_exists:false};result.checks.push('out-of-scope tool block stopped before actionable completion; real GUI/core creates no effect');
  demand(stats.credentialReads===0&&stats.externalAttempts===0&&stats.browserLaunches===0,'LOCAL_AUTHORITY');await pins();result.passed=true;
 }catch(e){result.passed=false;result.failures.push(e.safeCode||e.code||'LOCAL_CONTROL_ERROR');process.exitCode=1}
 return result;
}
// Minimal RPC shape adapted from the pinned own gui/test/browser.mjs helper. No framework/install.
class CDP{
 constructor(ws){this.ws=ws;this.n=0;this.pending=new Map();this.events=[];ws.onmessage=e=>{const m=JSON.parse(e.data);if(m.id){const p=this.pending.get(m.id);if(p){this.pending.delete(m.id);clearTimeout(p.timer);m.error?p.reject(fail('CDP_ERROR')):p.resolve(m.result)}}else this.events.push({at:Date.now(),...m})};ws.onclose=()=>{for(const p of this.pending.values()){clearTimeout(p.timer);p.reject(fail('CDP_CLOSED'))}this.pending.clear()}}
 static async open(url){const ws=new WebSocket(url);activeSockets.push(ws);await new Promise((r,j)=>{ws.onopen=r;ws.onerror=()=>j(fail('CDP_CONNECT'))});return new CDP(ws)}
 call(method,params={}){const id=++this.n;return new Promise((resolve,reject)=>{const timer=setTimeout(()=>{this.pending.delete(id);reject(fail('CDP_TIMEOUT'))},4000);this.pending.set(id,{resolve,reject,timer});this.ws.send(JSON.stringify({id,method,params}))})}
 async eval(expression){const r=await this.call('Runtime.evaluate',{expression,returnByValue:true,awaitPromise:true});demand(!r.exceptionDetails,'PAGE_EXCEPTION');return r.result.value}
 close(){this.ws.close()}
}
async function approvedLive(){
 // This function is never called by controls. Caller must supply BOTH explicit later authorization and exact hashes.
 let key;try{stats.credentialReads++;const settings=JSON.parse(await readFile(policy.credential.path,'utf8'));key=settings[policy.credential.field];demand(typeof key==='string'&&key.trim(),'KEY_UNAVAILABLE')}catch{throw fail('KEY_UNAVAILABLE')}
 const dir=join(run,'live'),proxy=await makeProxy({mode:'live',approved:true,url:policy.endpoint},dir,key),gui=await startGUI(proxy,dir),result={mode:'live, later-authorized only',stats,gui:gui.ownership,assets:gui.assets,observations:[],provider:null};
 const chrome=await realpath(policy.browser.binary),profile=await mkdtemp(join(run,'chrome-'));
 const args=['--headless=new','--mute-audio','--disable-gpu','--disable-background-networking','--disable-sync','--disable-background-timer-throttling','--disable-renderer-backgrounding','--no-first-run','--no-default-browser-check','--no-proxy-server','--password-store=basic','--use-mock-keychain',`--host-resolver-rules=${policy.browser.host_resolver_rules}`,`--user-data-dir=${profile}`,'--remote-debugging-port=0',gui.origin];
 const child=spawn(chrome,args,{cwd:run,env:closedEnv(run),stdio:['ignore','pipe','pipe']});stats.browserLaunches++;activeChildren.push(child);child.exited=new Promise(r=>child.once('exit',(code,signal)=>r({code,signal})));child.stdout.pipe(createWriteStream(join(run,'chrome.stdout.log'),{mode:0o600}));child.stderr.pipe(createWriteStream(join(run,'chrome.stderr.log'),{mode:0o600}));
 const active=await waitUntil(async()=>{try{return await readFile(join(profile,'DevToolsActivePort'),'utf8')}catch{return null}},6000,'CDP_LISTEN_TIMEOUT');const [portText]=active.split('\n'),port=Number(portText);demand(Number.isInteger(port)&&port>0&&port!==8088,'CDP_PORT');result.chrome=await ownership(child,chrome,args,port);
 const targets=await fetch(`http://127.0.0.1:${port}/json/list`,{redirect:'error',signal:AbortSignal.timeout(3000)}).then(r=>r.json());const pages=targets.filter(x=>x.type==='page'&&[gui.origin,gui.origin+'/'].includes(x.url));demand(pages.length===1,'OWNED_PAGE');
 const endpoint=new URL(pages[0].webSocketDebuggerUrl);demand(endpoint.hostname==='127.0.0.1'&&Number(endpoint.port)===port,'OWNED_CDP_URL');const page=await CDP.open(endpoint.href);await page.call('Runtime.enable');await page.call('Page.enable');await page.call('Network.enable');
 await waitUntil(()=>page.eval(`document.querySelector('#connection')?.textContent==='Connected'`),5000,'PAGE_READY');
 // No speaker, mic, voice enumeration, synthesis callback probe, focus/frontmost, hint or reconnect exercise.
 await page.eval(`document.querySelector('#read-aloud').checked=false`);
 const snapshot=()=>page.eval(`({at:Date.now(),state:document.querySelector('#state').textContent,readAloud:document.querySelector('#read-aloud').checked,connected:document.querySelector('#connection').textContent,artifacts:Array.from(document.querySelectorAll('.artifact'),n=>({id:n.dataset.partId,role:n.dataset.role,streaming:n.classList.contains('streaming'),text:n.innerText}))})`);
 result.before=await snapshot();await page.eval(`document.querySelector('#draft').value=${JSON.stringify(policy.fixture.prompt)};document.querySelector('#draft').dispatchEvent(new Event('input',{bubbles:true}));document.querySelector('#composer').requestSubmit()`);
 const deadline=performance.now()+policy.request.observation_ms;let partial=false,tool=false,ended=false,finalText='',started=false;
 try{
  while(performance.now()<deadline){demand(!proxy.state.halted,'LIVE_PROXY_STOP');const s=await snapshot();demand(!s.readAloud&&s.connected==='Connected','LIVE_PAGE_CHANGED');
   const events=page.events.splice(0);for(const e of events){if(e.method==='Runtime.exceptionThrown')throw fail('PAGE_EXCEPTION');if(e.method==='Network.webSocketClosed')throw fail('PAGE_DISCONNECTED');if(e.method==='Network.webSocketFrameReceived'){let m;try{m=JSON.parse(e.params.response.payloadData)}catch{continue}if(m.type==='tool_finished'&&m.name==='read_file'&&!m.is_error)tool=['Lantern Orchard','Reviewers needed: 3','Folder: blue'].every(x=>String(m.result).includes(x));if(m.type==='turn_ended'){ended=true;finalText=m.text||'';result.turnError=!!m.error}}}
   started ||= s.state!=='idle';if(!partial&&s.artifacts.some(a=>a.role==='agent'&&a.streaming&&a.text.trim())){partial=true;result.observations.push({kind:'actual_partial_dom',...s});const png=await page.call('Page.captureScreenshot',{format:'png'});await writeFile(join(run,'partial.png'),Buffer.from(png.data,'base64'));}
   if(ended){result.observations.push({kind:'actual_final_dom',...s});break}await sleep(25);
  }
  result.provider=proxy.summary();result.actualReadResult=tool;result.actualPartialDOM=partial;result.ended=ended;result.finalText=finalText;result.started=started;
  result.passed=partial&&tool&&ended&&!result.turnError&&finalText.trim().length>0;result.semanticReview='Root must compare actual final meaning with the three fixture facts; no arbitrary exact answer phrase required.';
  const png=await page.call('Page.captureScreenshot',{format:'png'});await writeFile(join(run,'final.png'),Buffer.from(png.data,'base64'));
 }catch(e){result.passed=false;result.failure=e.safeCode||'LIVE_ERROR';result.provider=proxy.summary();try{result.partialReceipt=await snapshot()}catch{}}
 finally{proxy.stop();page.close();key=undefined}
 return result;
}
async function cleanup(){const receipts=[];for(const ws of activeSockets)try{ws.close()}catch{}for(const server of activeServers){server.closeIdleConnections();server.closeAllConnections();await new Promise(r=>server.close(r))}for(const child of activeChildren.toReversed())receipts.push({pid:child.pid,...await stopChild(child)});return receipts}
const mode=process.argv[2]||'plan';
if(!['controls','live'].includes(mode)){console.log('Preparation only. controls runs loopback tests; live requires later explicit root authorization and exact hashes.');process.exit(0)}
const opts=Object.fromEntries(process.argv.slice(3).map(x=>{const [k,...v]=x.replace(/^--/,'').split('=');return[k,v.join('=')||true]}));
demand(sha(await readFile(fileURLToPath(import.meta.url)))===opts['launcher-sha256'],'LAUNCHER_PIN');demand(sha(await readFile(join(here,'policy.json')))===opts['policy-sha256'],'POLICY_PIN');policy=JSON.parse(await readFile(join(here,'policy.json'),'utf8'));
if(mode==='live')demand(opts['explicit-later-root-paid-browser-authorization']===true,'LIVE_NOT_AUTHORIZED');
const label=opts.label;demand(typeof label==='string'&&/^[a-z0-9-]+$/.test(label),'LABEL');run=join(policy.source_root,'.ch07-local/live-exercise-prep',label);await mkdir(run,{mode:0o700});let result;
try{await pins();result=mode==='controls'?await localControls():await approvedLive()}catch(e){result={passed:false,mode,failure:e.safeCode||'HARNESS_ERROR',stats};process.exitCode=1}
finally{try{result.cleanup=await cleanup()}catch(e){result.passed=false;result.cleanupFailure=e.safeCode||'CLEANUP_ERROR';process.exitCode=1}result.stats=stats;result.source=policy.commit;result.run=run;await writeFile(join(run,'receipt.json'),JSON.stringify(result,null,2)+'\n',{mode:0o600});if(mode==='controls')await writeFile(join(here,label+'-result.json'),JSON.stringify(result,null,2)+'\n');console.log(JSON.stringify({mode,passed:result.passed,failure:result.failure||result.failures,checks:result.checks,stats,cleanup:result.cleanup,receipt:join(run,'receipt.json')}));if(!result.passed)process.exitCode=1}
