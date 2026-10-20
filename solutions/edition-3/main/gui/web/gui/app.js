'use strict';
const $ = id => document.getElementById(id);
let socket, clientID = '', selected = '', retry = 500, session = '';
const agents = new Map(), cards = new Map(), order = [];
let audio = [], speakingNow = null;
const typing = new Set(), delayed = new Map(), requests = new Map(), voiceAgents = new Set();
let inputOccupied=false;
const element = (tag, cls, text) => { const e = document.createElement(tag); if(cls) e.className=cls; if(text !== undefined) e.textContent=String(text); return e; };
function send(type, data={}) { if(socket?.readyState !== WebSocket.OPEN) { $('errors').textContent='Disconnected; wait for session replay before sending.'; return false; } socket.send(JSON.stringify({type,agent_id:selected,...data})); return true; }
function nearBottom(){ return document.documentElement.scrollHeight-innerHeight-scrollY < 180; }
function follow(wasNear){ if(wasNear) requestAnimationFrame(()=>scrollTo(0,document.documentElement.scrollHeight)); else $('follow').hidden=false; }
$('follow').onclick=()=>{scrollTo({top:document.documentElement.scrollHeight,behavior:'smooth'});$('follow').hidden=true;};
function fullDetails(title,child,open){const details=element('details');details.open=!!open;details.append(element('summary','',title),child);return details;}
function toolContent(c){
 const body=element('div');
 if(c.name==='edit_file'&&c.args&&typeof c.args.old_text==='string'&&typeof c.args.new_text==='string'){
  body.append(element('div','muted',c.args.path||''));
  const diff=element('pre','diff');for(const line of c.args.old_text.split('\n'))diff.append(element('div','deletion','− '+line));for(const line of c.args.new_text.split('\n'))diff.append(element('div','insertion','+ '+line));body.append(diff);
 }
 body.append(element('pre','code',c.partial?c.text:JSON.stringify(c.args??{},null,2)));return body;
}
function jobInfo(job){const box=element('div','job-state');let text=`Job ${job.handle} · ${job.status} · ${job.bytes} bytes`;if(job.exit_code!==undefined)text+=` · exit ${job.exit_code}`;if(job.reason)text+=` · ${job.reason}`;box.append(element('div','',text));if(job.output?.locator)box.append(element('div','muted',job.output.locator));return box;}
function renderCard(record) {
 const {frame, node} = record, c=frame.part, wasOpen=node.querySelector('details')?.open;
 node.replaceChildren();node.className=`card ${c.type} ${c.role==='human'?'human':''} ${c.partial?'partial':''} ${c.is_error?'error':''}`;
 const heading=element('div','card-title');
 const names={text:c.role==='human'?'You':c.role==='system'?'System':'Assistant',thinking:'Thinking',tool_call:`Tool request · ${c.name||'pending'}`,tool_result:`Tool result · ${c.name||c.call_id||''}`,image:'Image',audio:'Audio',video:'Video',file:'File'};
 heading.append(element('span','',names[c.type]||c.type));
 if(c.partial) heading.append(element('span','','Streaming…'));
 if(c.type==='tool_call'&&frame.type!=='job_update') heading.append(element('span',c.dispatched?'badge':'',c.dispatched?'Dispatched':'Requested · not dispatched'));
 if(c.is_error) heading.append(element('span','','Error'));
 if(c.from) heading.append(element('span','muted',`${c.from.vendor} / ${c.from.model}`));
 if(!c.partial && c.text && ['text','thinking','tool_call','tool_result'].includes(c.type) && !c.text.startsWith('[opaque')) {
  const speak=element('button','','Speak');speak.title='Speak the complete readable part';speak.onclick=()=>send('speak',{agent_id:frame.agent_id,part_id:frame.part_id});heading.append(speak);
 }
 node.append(heading);
 if(c.type==='tool_call') {if(c.partial)node.append(toolContent(c));else {node.append(element('div','preview',`${c.name} · ${JSON.stringify(c.args??{}).slice(0,160)}`));node.append(fullDetails('Full request / edit diff',toolContent(c),wasOpen));}if(c.call_id)node.append(element('div','muted',`Provider call: ${c.call_id}`));}
 else if(c.type==='thinking'&&!c.partial){node.append(element('div','preview',(c.text||'').slice(0,180)));node.append(fullDetails('Full thinking',element('div','body',c.text||''),wasOpen));}
 else if(c.type==='tool_result') {
  const text=c.text||'',preview=element('div','body preview',text.length>400?text.slice(0,400)+'…':text);node.append(preview);
  const details=element('details');details.open=!!wasOpen;details.append(element('summary','',`Full result · ${text.length.toLocaleString()} characters`),element('pre','body',text));node.append(details);
 } else if(['image','audio','video','file'].includes(c.type)) {
  node.append(element('div','media',c.text||`${c.type} attachment · ${c.mime||'unknown type'}`));
  if(c.locator) {let safe=false;try { const u=new URL(c.locator);safe=['http:','https:'].includes(u.protocol); }catch{}
   if(safe) {const a=element('a','',c.locator);a.href=c.locator;a.target='_blank';a.rel='noopener noreferrer';node.append(a);}else node.append(element('div','muted',c.locator)); }
 } else node.append(element('div','body',c.text||''));
 if(c.job)node.append(jobInfo(c.job));
 node.hidden=frame.agent_id!==selected;
}
function apply(frame) {
 if(['part_final','part_delta','tool_dispatched','tool_finished'].includes(frame.type)) {
  const key=`${frame.agent_id}/${frame.part_id}`;let record=cards.get(key);
  if(frame.type==='part_delta') {
   if(record&&!record.frame.part.partial)return;
   frame={...frame,part:{type:frame.kind,name:frame.name,role:'agent',partial:true,text:(record?.frame.part.text||'')+(frame.chunk||'')}};
  }
  const dispatched=record?.frame.part.dispatched||frame.part?.dispatched;
  if(!record){record={frame,node:element('article','card')};cards.set(key,record);order.push(key);$('conversation').append(record.node);}
  record.frame=frame;record.frame.part.dispatched=!!dispatched;renderCard(record);
 } else if(frame.type==='job_update') {
  const key=`${frame.agent_id}/${frame.part_id}`;let record=cards.get(key);if(!record){record={frame,node:element('article','card')};cards.set(key,record);$('conversation').append(record.node);}else record.frame.part.job=frame.part.job;renderCard(record);
 } else if(frame.type==='part_removed'){const key=`${frame.agent_id}/${frame.part_id}`,r=cards.get(key);if(r?.frame.part.partial){r.node.remove();cards.delete(key);}
 } else if(frame.type==='state'||frame.type==='state_changed') {
  const a=agents.get(frame.agent_id), next=frame.type==='state'?frame.data:frame.state, previous=a?.state;if(a) a.state=next;
  if(frame.agent_id===selected) controls();
  if(!frame.replay&&next.continue_epoch!==previous?.continue_epoch)typing.delete(frame.agent_id);
  if(!frame.replay&&next.epoch!==previous?.epoch)cancelAudio(frame.agent_id,false);
 } else if(frame.type==='turn_ended') {for(const [key,r] of cards) {if(r.frame.agent_id===frame.agent_id&&r.frame.part.partial){r.node.remove();cards.delete(key);}}}
 else if(frame.type==='error') $('errors').textContent=String(frame.data);
 else if(frame.type==='speech_requested'&&!frame.replay){
  if(frame.data.epoch!==agents.get(frame.agent_id)?.state.epoch)return;
  requests.set(frame.data.id,{agent:frame.agent_id,epoch:frame.data.epoch});
  $('audio-status').textContent='Browser speech requested';send('claim_speech',{agent_id:frame.agent_id,speech_id:frame.data.id});
 }else if(frame.type==='speech_claimed'&&!frame.replay){
  const request=requests.get(frame.data.id);requests.delete(frame.data.id);
  if(frame.data.owner!==clientID||!request||request.epoch!==frame.data.epoch||frame.data.epoch!==agents.get(frame.agent_id)?.state.epoch)return;
  const task={agent:frame.agent_id,...frame.data};const occupied=voiceAgents.has(task.agent);voiceAgents.add(task.agent);if(!occupied)send('pause',{agent_id:task.agent});
  if(task.reason==='prompt'&&$('mic').checked){const timer=setTimeout(()=>{delayed.delete(task.id);enqueueAudio(task);},2000);delayed.set(task.id,{task,timer});}else enqueueAudio(task);
 }else if(frame.type==='speech_cancelled'){
  cancelAudio(frame.agent_id,false);
  if(frame.data.release&&frame.data.owner===clientID)maybeResume(frame.agent_id);
 }else if(frame.type==='speech_idle'){
  if(frame.data.owner===clientID&&frame.data.epoch===agents.get(frame.agent_id)?.state.epoch){voiceAgents.delete(frame.agent_id);maybeResume(frame.agent_id);}
 }else if(frame.type==='speech_done'){requests.delete(frame.data.id);$('audio-status').textContent=`Speech: ${frame.data.status}`;}

}
function controls(){const a=agents.get(selected);if(!a)return;const s=a.state;$('vendor').textContent=a.config.vendor;$('model').textContent=a.config.model;$('surface').textContent=a.config.surface;$('state').textContent=s.paused?'Paused':s.turn.replaceAll('_',' ');$('reason').textContent=s.reason||'';$('speaking').checked=s.speaking;$('pause').disabled=s.paused;$('continue').disabled=!s.paused;$('audio-status').textContent=s.speaking?'Speech enabled':'Speech disabled';}
function replay(data){
 const previousAgents=new Map(agents);
 cancelAudio(undefined,false);clientID=data.client;session=data.session;cards.clear();order.length=0;$('conversation').replaceChildren();agents.clear();$('agent').replaceChildren();
 for(const a of data.agents){if(!a.state.paused||a.state.continue_epoch!==previousAgents.get(a.id)?.state.continue_epoch)typing.delete(a.id);agents.set(a.id,a);const option=element('option','',a.id);option.value=a.id;$('agent').append(option);}
 if(!agents.has(selected))selected=data.agents[0]?.id||'';$('agent').value=selected;
 for(const frame of data.events) apply(frame);for(const frame of data.partials) apply(frame);controls();
}
function connect(){socket=new WebSocket(`${location.protocol==='https:'?'wss':'ws'}://${location.host}/ws`);
 socket.onopen=()=>{$('connection').textContent='Connected · synchronizing';retry=500;};
 socket.onmessage=event=>{let frame;try{frame=JSON.parse(event.data);}catch{return;}const near=nearBottom();if(frame.type==='replay'){replay(frame.data);$('connection').textContent='Connected';}else apply(frame);follow(near);};
 socket.onclose=()=>{$('connection').textContent='Disconnected · reconnecting';cancelAudio(undefined,false);setTimeout(connect,retry);retry=Math.min(retry*2,8000);};
 socket.onerror=()=>socket.close();
}
function busy(agent){return voiceAgents.has(agent)||[...requests.values()].some(r=>r.agent===agent)||speakingNow?.agent===agent||audio.some(t=>t.agent===agent)||[...delayed.values()].some(x=>x.task.agent===agent);}
function maybeResume(agent){if(!busy(agent)&&!typing.has(agent))send('unpause',{agent_id:agent});}
function enqueueAudio(task){audio.push(task);speakNext();}
function cancelAudio(agent,resume=false){
 for(const [id,r] of requests)if(!agent||r.agent===agent)requests.delete(id);
 if(agent)voiceAgents.delete(agent);else voiceAgents.clear();
 const removed=audio.filter(t=>!agent||t.agent===agent);audio=audio.filter(t=>agent&&t.agent!==agent);
 for(const [id,x] of delayed){if(!agent||x.task.agent===agent){clearTimeout(x.timer);delayed.delete(id);removed.push(x.task);}}
 if(speakingNow&&(!agent||speakingNow.agent===agent)){removed.push(speakingNow);speakingNow=null;window.speechSynthesis?.cancel();}
 for(const t of removed)send('speech_done',{agent_id:t.agent,speech_id:t.id,status:'cancelled'});
 if(resume)for(const id of new Set(removed.map(t=>t.agent)))maybeResume(id);
 speakNext();
}
function speakNext(){
 if(speakingNow||!audio.length)return;
 const task=audio.shift();speakingNow=task;
 const complete=status=>{if(speakingNow!==task)return;speakingNow=null;send('speech_done',{agent_id:task.agent,speech_id:task.id,status});speakNext();};
 if(!window.speechSynthesis||!window.SpeechSynthesisUtterance){complete('unsupported');return;}
 // Read every codepoint in order. No preview truncation or chunk dropping.
 const chunks=Array.from(task.text);let offset=0;
 const next=()=>{if(speakingNow!==task)return;if(offset>=chunks.length){complete('done');return;}const utterance=new SpeechSynthesisUtterance(chunks.slice(offset,offset+1200).join(''));offset+=1200;utterance.onend=next;utterance.onerror=()=>complete('failed');speechSynthesis.speak(utterance);};next();
}
$('agent').onchange=()=>{selected=$('agent').value;for(const r of cards.values())r.node.hidden=r.frame.agent_id!==selected;controls();};
$('menu-toggle').onclick=()=>{const open=$('menu').hidden;$('menu').hidden=!open;$('menu-toggle').setAttribute('aria-expanded',String(open));};
$('pause').onclick=()=>send('pause');$('continue').onclick=()=>{typing.delete(selected);cancelAudio(selected,false);send('unpause');};$('interrupt').onclick=()=>{cancelAudio(selected,false);send('interrupt');};
$('speaking').onchange=()=>{if(!$('speaking').checked)cancelAudio(selected);send('speaking',{enabled:$('speaking').checked});};
$('cancel-audio').onclick=()=>{cancelAudio(selected);send('speaking',{enabled:false});};
function submit(type){const text=$('input').value;if(text.trim()&&send(type,{text})){$('input').value='';inputOccupied=false;typing.delete(selected);maybeResume(selected);$('errors').textContent='';}}
$('input').oninput=()=>{const nonempty=!!$('input').value.trim();if(nonempty&&!inputOccupied){typing.add(selected);send('pause');}else if(!nonempty&&typing.delete(selected))maybeResume(selected);inputOccupied=nonempty;};
$('composer').onsubmit=e=>{e.preventDefault();submit('input');};$('hint').onclick=()=>submit('hint');$('input').onkeydown=e=>{if(e.key==='Enter'&&(e.ctrlKey||e.metaKey)){e.preventDefault();submit('input');}};
addEventListener('pagehide',()=>cancelAudio(undefined,false));connect();
