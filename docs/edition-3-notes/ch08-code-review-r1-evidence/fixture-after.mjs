// Controlled page regression: actual production app modules/handlers/buttons and
// SpeechQueue; minimal deterministic DOM/socket/synthesis boundaries. NOT a
// browser, native speech, audio, permission, microphone, or listening exercise.
import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import {fileURLToPath,pathToFileURL} from 'node:url';
import {resolve,dirname} from 'node:path';
import vm from 'node:vm';
const root=resolve(dirname(fileURLToPath(import.meta.url)),'../web/gui');
class Node {
 constructor(tag=''){this.tagName=tag.toUpperCase();this.children=[];this.parentNode=null;this._text='';this.className='';this.dataset={};this.hidden=false;this.value='';this.checked=false;this.events=new Map();this.attributes=new Map();this.style={values:new Map(),setProperty(k,v){this.values.set(k,v)},getPropertyValue(k){return this.values.get(k)||''}};this.scrollHeight=10;this.scrollTop=0;this.clientHeight=100;this.classList={add:(...v)=>{this.className=[...new Set([...this.className.split(' ').filter(Boolean),...v])].join(' ')},remove:(...v)=>{this.className=this.className.split(' ').filter(s=>!v.includes(s)).join(' ')},contains:v=>this.className.split(' ').includes(v)};}
 append(...nodes){for(let node of nodes){if(typeof node==='string'){const n=new Node('#text');n.textContent=node;node=n}if(node.tagName==='#FRAGMENT'){this.append(...node.children);continue}node.parentNode=this;this.children.push(node)}}
 replaceChildren(...nodes){for(const n of this.children)n.parentNode=null;this.children=[];this._text='';this.append(...nodes)}
 set textContent(v){this.replaceChildren();this._text=String(v??'')}get textContent(){return this._text+this.children.map(n=>n.textContent).join('')}
 addEventListener(k,f){if(!this.events.has(k))this.events.set(k,[]);this.events.get(k).push(f)}
 dispatchEvent(event){event.target=this;this['on'+event.type]?.(event);for(const f of this.events.get(event.type)||[])f(event)}
 click(){this.dispatchEvent({type:'click',preventDefault(){}})}
 setAttribute(k,v){this.attributes.set(k,String(v))}
 getBoundingClientRect(){return {width:1100}}
 remove(){if(this.parentNode)this.parentNode.children=this.parentNode.children.filter(n=>n!==this);this.parentNode=null}
 querySelector(selector){return this.querySelectorAll(selector)[0]||null}
 querySelectorAll(selector){let out=[];const match=n=>selector.startsWith('.')?n.className.split(' ').includes(selector.slice(1)):n.tagName===selector.toUpperCase();for(const n of this.children){if(match(n))out.push(n);out.push(...n.querySelectorAll(selector))}return out}
}
async function page(){
 // Inert current-page boundaries, not an HTML/layout engine. Obtain IDs/tags
 // from the real shell and connect only the two artifact/scroll containers.
 const html=await readFile(resolve(root,'index.html'),'utf8'),ids=new Map();
 for(const [,tag,attrs,id] of html.matchAll(/<([a-z][\w-]*)\b([^>]*\bid="([^"]+)"[^>]*)>/g)){
  assert(!ids.has(id),`duplicate actual page id: ${id}`);const node=new Node(tag);node.checked=/\bchecked\b/.test(attrs);ids.set(id,node);
 }
 ids.get('chat-scroll').append(ids.get('transcript'));
 ids.get('actions-scroll').append(ids.get('actions'));ids.get('actions').append(ids.get('actions-empty'));
 const doc=new Node('html');doc.style.setProperty('--sidebar-width','240px');doc.style.setProperty('--actions-width','320px');
 const document={hidden:true,documentElement:doc,getElementById:id=>{assert(ids.has(id),`missing actual page id: ${id}`);return ids.get(id)},createElement:tag=>new Node(tag),createTextNode:text=>{const n=new Node('#text');n.textContent=text;return n},createDocumentFragment:()=>new Node('#fragment'),addEventListener(){}};
 let frames=new Map(),next=0;const spoken=[],sent=[],errors=[];let canceled=0;
 class Utterance{constructor(text){this.text=text}}
 const synth={speak(u){spoken.push(u)},cancel(){canceled++}};
 const sockets=[];class Socket{static OPEN=1;constructor(url){this.url=url;this.readyState=1;sockets.push(this)}send(data){sent.push(JSON.parse(data))}close(){this.readyState=3}}
 const media={matches:false,addEventListener(){}};
 const window={speechSynthesis:synth,SpeechSynthesisUtterance:Utterance,matchMedia:()=>media,addEventListener(){}};
 const context=vm.createContext({window,document,getComputedStyle:node=>node.style,WebSocket:Socket,SpeechSynthesisUtterance:Utterance,speechSynthesis:synth,location:{protocol:'http:',host:'controlled.invalid'},URL,console:{log(){},error:e=>errors.push(String(e))},requestAnimationFrame:f=>{frames.set(++next,f);return next},cancelAnimationFrame:id=>frames.delete(id),setTimeout:()=>{throw Error('unexpected reconnect timer')},clearTimeout(){}});
 const modules=new Map();async function load(path){if(modules.has(path))return modules.get(path);assert(path.startsWith(root+'/'));const module=new vm.SourceTextModule(await readFile(path,'utf8'),{context,identifier:pathToFileURL(path).href});modules.set(path,module);await module.link((specifier,owner)=>load(resolve(dirname(fileURLToPath(owner.identifier)),specifier)));return module}
 await (await load(resolve(root,'app.js'))).evaluate();const socket=sockets[0];assert.equal(sockets.length,1);socket.onopen();
 const deliver=m=>{socket.onmessage({data:JSON.stringify(m)});assert.deepEqual(errors,[]);assert(!ids.get('notice').textContent.startsWith('Server message rejected:'))};
 const settings={theme:'dark',tts_enabled:true,tts_speed:1,font_size:16,max_tool_rounds:0};
 deliver({type:'current_settings',settings});
 const card=id=>{const found=[...ids.get('transcript').children,...ids.get('actions').children].filter(n=>String(n.dataset.partId)===String(id));assert.equal(found.length,1,`one artifact per display identity ${id}`);return found[0]};
 const cancel=()=>{
  const control=ids.get('read-aloud'),before=canceled;assert.equal(control.disabled,false);assert.equal(control.checked,true);
  control.checked=false;control.dispatchEvent({type:'change'});
  assert.deepEqual(sent.at(-1),{type:'update_settings',settings:{tts_enabled:false}},'real control must send sparse settings request');
  assert.equal(control.checked,true,'draft must await accepted snapshot');assert.equal(canceled,before,'draft must not cancel speech');
  deliver({type:'settings_changed',settings:{...settings,tts_enabled:false}});
  assert.equal(control.checked,false);assert.equal(canceled,before+1,'accepted disabled snapshot must cancel speech');
 };
 return {ids,spoken,sent,deliver,card,window,cancel,canceled:()=>canceled,finish:()=>{const u=spoken.at(-1);u?.onstart?.({type:'start'});u?.onend?.({type:'end'})}};
}
function call(id,name,args){return {type:'part_final',part_id:id,role:'agent',part:{type:'tool_call',name,args,call_id:'same-model-id',display_id:id}}}
function dispatch(id,name,input){return {type:'tool_dispatched',part_id:id,name,input,call_id:'same-model-id'}}
const checks=[];
for(const scenario of [
 {name:'read_file',args:{path:'safe.txt',start_line:2,end_line:3},result:'LINE_TWO visible result\n<script>literal only</script>\nLINE_THREE',error:false},
 {name:'run_command',args:{command:'printf safe'},result:JSON.stringify({stdout:'RESULT_STDOUT',stderr:'RESULT_STDERR',exit_code:7}),error:true},
]){
 const p=await page(),id=41;p.deliver(call(id,scenario.name,scenario.args));p.deliver(dispatch(id,scenario.name,scenario.args));
 assert.equal(p.spoken.length,1,'one short automatic dispatch summary');assert(!p.spoken[0].text.includes(scenario.result));p.finish();
 p.deliver({type:'tool_finished',part_id:id,call_id:'same-model-id',result:scenario.result,is_error:scenario.error,signature:'OPAQUE_PROVIDER_SENTINEL'});
 assert.equal(p.spoken.length,1,'results must not auto-read');p.deliver(call(id,scenario.name,scenario.args));p.deliver({type:'turn_ended'});const c=p.card(id);assert.equal(c.parentNode,p.ids.get('actions'),'tool artifact belongs in Actions');assert(c.querySelector('.tool-result').textContent.includes(scenario.result));
 c.querySelector('.speaker').click();const utterance=p.spoken.at(-1);console.log(JSON.stringify({scenario:scenario.name,phase:'live speaker',queued:utterance.text}));
 assert(utterance.text.includes(scenario.result),'actual tool speaker omitted visible result');assert(utterance.text.includes(scenario.name.replaceAll('_',' ')),'call context lost');assert(utterance.text.includes(scenario.args.path||scenario.args.command),'arguments lost');assert(utterance.text.includes(scenario.error?'Tool error':'Tool result'),'result/error distinction missing');assert(!utterance.text.includes('OPAQUE_PROVIDER_SENTINEL'));
 assert.equal(p.sent.at(-1).type,'pause');p.finish();assert.equal(p.sent.at(-1).type,'unpause');
 // Replay is silent but the real recreated speaker retains complete content.
 const count=p.spoken.length;p.deliver({type:'replay_begin'});p.deliver({...call(id,scenario.name,scenario.args),replay:true});p.deliver({...dispatch(id,scenario.name,scenario.args),replay:true});p.deliver({type:'tool_finished',part_id:id,result:scenario.result,is_error:scenario.error,replay:true});p.deliver({type:'replay_end'});assert.equal(p.spoken.length,count,'replay spoke automatically');p.card(id).querySelector('.speaker').click();assert(p.spoken.at(-1).text.includes(scenario.result),'replayed speaker omitted result');
 // Actual production callback-generation guard; controlled callbacks only.
 const stale=p.spoken.at(-1);p.cancel();assert.equal(p.sent.at(-1).type,'unpause');p.card(id).querySelector('.speaker').click();const current=p.spoken.at(-1);stale.onend({type:'end'});stale.onerror({type:'error',error:'controlled stale'});assert.equal(p.sent.at(-1).type,'pause');current.onerror({type:'error',error:'controlled current'});assert.equal(p.sent.at(-1).type,'unpause');
 checks.push(scenario.name+' actual handler/speaker: result, error, replay silence, cancellation/stale callback/current error release');
}
// Empty/non-readable provider blocks never become speaker text or automatic audio.
{
 const p=await page();for(const [id,part] of [[51,{type:'thinking',text:'',signature:'SECRET'}],[52,{type:'opaque',data:{signature:'SECRET',redacted:'SECRET'}}]]){p.deliver({type:'replay_begin'});p.deliver({type:'part_final',part_id:id,part});p.deliver({type:'replay_end'});p.card(id).querySelector('.speaker').click()};assert.equal(p.spoken.length,0);checks.push('empty/opaque actual artifact speakers silent');
}
// Optional real Agent/hub wire from the Go regression, not constructed dispatch input.
if(process.argv[2]){
 const wire=JSON.parse(await readFile(process.argv[2],'utf8'));const p=await page();for(const m of wire.refused_live){p.deliver(m);p.finish()}assert.equal(p.spoken.filter(u=>/^(write file|run command|tool limits):/.test(u.text)).length,0,'refused batch generated execution summary');assert.equal(p.ids.get('actions').querySelectorAll('.tool-result').length,3);for(const c of p.ids.get('actions').children.filter(n=>n.querySelector('.tool-result'))){c.querySelector('.speaker').click();assert(p.spoken.at(-1).text.includes('no tool in this batch ran'));p.finish()}
 const before=p.spoken.length;for(const m of wire.refused_replay)p.deliver(m);assert.equal(p.spoken.length,before,'refusal replay spoke');
 const fresh=await page();for(const m of wire.durable_replay)fresh.deliver(m);assert.equal(fresh.spoken.length,0,'durable reconnect auto-spoke');assert.equal(fresh.ids.get('actions').querySelectorAll('.tool-result').length,6,'refusal/admitted artifacts lost or duplicated');const socketPage=await page();assert(wire.websocket_refused?.length,'real WebSocket receipt required');for(const m of wire.websocket_refused){socketPage.deliver(m);socketPage.finish()}assert.equal(socketPage.spoken.filter(u=>/^(write file|run command|tool limits):/.test(u.text)).length,0,'real WebSocket refusal spoke execution');const legacy=await page();for(const m of wire.legacy_replay)legacy.deliver(m);assert.equal(legacy.spoken.length,0);assert.equal(legacy.ids.get('actions').querySelectorAll('.tool-result').length,6);checks.push('actual core/hub refusal live + reconnect + serialized replay feed actual app without execution summaries; on-demand errors retained');
}
console.log(JSON.stringify({passed:true,scope:'Controlled actual production page handlers and SpeechQueue, NOT native/browser/permission/listening proof',checks},null,2));
