import vm from 'node:vm';
import fs from 'node:fs/promises';
import path from 'node:path';
import assert from 'node:assert/strict';
const root=path.resolve(process.argv[2]);
class Element {
  constructor(tag='div'){this.tagName=tag.toUpperCase();this.children=[];this.dataset={};this.listeners=new Map();this.value='';this.checked=false;this.hidden=false;this.textContent='';this.className='';this.parent=null;this.classList={remove:()=>{},add:()=>{}};}
  addEventListener(t,f){this.listeners.set(t,f)}
  append(...xs){for(const x of xs){this.children.push(x);if(x&&typeof x==='object')x.parent=this;}}
  replaceChildren(...xs){this.children=[];this.append(...xs)}
  querySelector(q){const test=e=>q.startsWith('.')?e.className?.split(' ').includes(q.slice(1)):e.tagName?.toLowerCase()===q;for(const c of this.children){if(test(c))return c;const x=c.querySelector?.(q);if(x)return x;}return null;}
  remove(){if(this.parent)this.parent.children=this.parent.children.filter(x=>x!==this)}
  setAttribute(k,v){this[k]=v}
  requestSubmit(){this.listeners.get('submit')?.({preventDefault(){}})}
}
async function page(){
 const elements=new Map(),events=new Map(),frames=new Map(),sockets=[],spoken=[],errors=[];let seq=0;
 const get=id=>{if(!elements.has(id))elements.set(id,new Element());return elements.get(id)};get('read-aloud').checked=true;
 const document={getElementById:get,createElement:t=>new Element(t),createDocumentFragment:()=>new Element('#fragment'),createTextNode:text=>{const n=new Element('#text');n.textContent=String(text);return n},addEventListener:(t,f)=>events.set('document:'+t,f),scrollingElement:{scrollHeight:100,scrollTop:0,clientHeight:100}};
 class Socket {static OPEN=1;constructor(url){this.url=url;this.readyState=1;this.sent=[];sockets.push(this)}send(x){this.sent.push(JSON.parse(x))}close(){this.readyState=3;this.onclose?.()}}
 class Utterance{constructor(text){this.text=text}}
 const synth={speak(u){spoken.push(u)},cancel(){}};
 const sandbox={console:{log(){},error(e){errors.push(String(e))}},document,location:{protocol:'http:',host:'inert.invalid'},WebSocket:Socket,speechSynthesis:synth,SpeechSynthesisUtterance:Utterance,addEventListener:(t,f)=>events.set(t,f),requestAnimationFrame:f=>{frames.set(++seq,f);return seq},cancelAnimationFrame:id=>frames.delete(id),setTimeout:()=>++seq,clearTimeout(){},setInterval:()=>++seq,clearInterval(){}};sandbox.window=sandbox;
 const ctx=vm.createContext(sandbox),cache=new Map();
 async function load(filename){filename=path.resolve(filename);if(cache.has(filename))return cache.get(filename);const m=new vm.SourceTextModule(await fs.readFile(filename,'utf8'),{context:ctx,identifier:filename});cache.set(filename,m);await m.link((name,ref)=>load(path.resolve(path.dirname(ref.identifier),name)));return m;}
 const module=await load(path.join(root,'app.js'));await module.evaluate();const socket=sockets[0];assert(socket);socket.onopen();
 const flush=()=>{let n=0;while(frames.size){if(++n>100)throw Error('bounded RAF drain exceeded');const batch=[...frames.values()];frames.clear();for(const f of batch)f();}};
 const receive=m=>{socket.onmessage({data:JSON.stringify(m)});flush();assert.equal(errors.length,0,errors.join('\n'));};
 return {get,events,spoken,socket,receive,flush,errors,input(text){get('draft').value=text;get('draft').listeners.get('input')({type:'input'});},commands(){return socket.sent.map(m=>m.type)}};
}
const results=[];
{
 const p=await page();p.input('keep this gate paused');assert.equal(p.commands().at(-1),'pause');const before=p.socket.sent.length;p.events.get('pagehide')();const outgoing=p.socket.sent.slice(before);
 results.push({name:'disconnect-must-not-release-agent-gate',expected:'close without unpause',outgoing,closed:p.socket.readyState===3,passed:!outgoing.some(m=>m.type==='unpause')});
}
{
 const p=await page();p.input('existing draft');p.receive({type:'part_final',part_id:101,part:{type:'text',text:'A readable sentence.'}});assert.equal(p.spoken.length,1);p.receive({type:'pause_changed',paused:false});const before=p.socket.sent.length;p.spoken[0].onend({type:'end'});const outgoing=p.socket.sent.slice(before);
 results.push({name:'stale-draft-does-not-repause-when-speech-ends-after-remote-unpause',outgoing,draft:p.get('draft').value,passed:!outgoing.some(m=>m.type==='pause')});
}
{
 const p=await page();p.get('read-aloud').checked=false;p.receive({type:'tool_dispatched',part_id:201,call_id:'review-call',name:'read_file',input:{path:'note.txt',start_line:4,end_line:8}});p.receive({type:'tool_finished',part_id:201,call_id:'review-call',name:'read_file',result:'DISTINCT_TOOL_RESULT_41',is_error:false});const card=p.get('transcript').children[0],button=card.querySelector('.speaker');assert(button);assert(card.querySelector('.tool-result'));button.onclick();const spoken=p.spoken.map(u=>u.text);
 results.push({name:'tool-result-readable-content-is-available-to-card-speech',resultMarker:'DISTINCT_TOOL_RESULT_41',spoken,passed:spoken.some(s=>s.includes('DISTINCT_TOOL_RESULT_41'))});
}
const out={scope:'Non-browser execution of unchanged production JS; inert interfaces and controlled speech callbacks, not native/browser proof',sourceRoot:root,results,passed:results.every(r=>r.passed)};console.log(JSON.stringify(out,null,2));process.exitCode=out.passed?0:1;
