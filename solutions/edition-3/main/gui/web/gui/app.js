'use strict';
const $ = id => document.getElementById(id);
let socket, clientID = '', selected = '', retry = 500, session = '';
const skillsByAgent=new Map();
const agents = new Map(), cards = new Map(), order = [];
let audio = [], speakingNow = null;
const typing = new Set(), delayed = new Map(), requests = new Map(), voiceAgents = new Set();
let inputOccupied=false;
let settings={theme:'dark',tts_enabled:true,tts_speed:1,font_size:16,max_tool_rounds:0,log_mode:'full',diagnostic_max_records:200},settingsReady=false;
const preferenceFields={'theme':'theme','tts-enabled':'tts_enabled','tts-speed':'tts_speed','font-size':'font_size','max-tool-rounds':'max_tool_rounds','log-mode':'log_mode','diagnostic-max-records':'diagnostic_max_records'};
const systemTheme=window.matchMedia?.('(prefers-color-scheme: dark)');
function applyTheme(){document.documentElement.dataset.theme=settings.theme==='system'?(systemTheme?.matches?'dark':'light'):settings.theme;}
if(systemTheme?.addEventListener)systemTheme.addEventListener('change',()=>{if(settings.theme==='system')applyTheme();});
function preferenceControls(){
 for(const [id,key] of Object.entries(preferenceFields)){const control=$(id);if(control.type==='checkbox')control.checked=settings[key];else control.value=String(settings[key]);}
 $('font-size-value').textContent=String(settings.font_size);$('tts-speed-value').textContent=`${settings.tts_speed}×`;$('preferences').disabled=!settingsReady;
}
function automaticSpeech(task){return !['on demand','user'].includes(task.reason);}
function cancelAutomaticAudio(){
 const removed=audio.filter(automaticSpeech);audio=audio.filter(task=>!automaticSpeech(task));
 for(const [id,item] of delayed){if(automaticSpeech(item.task)){clearTimeout(item.timer);delayed.delete(id);removed.push(item.task);}}
 if(speakingNow&&automaticSpeech(speakingNow)){removed.push(speakingNow);speakingNow=null;window.speechSynthesis?.cancel();}
 for(const task of removed)send('speech_done',{agent_id:task.agent,speech_id:task.id,status:'disabled'});
 speakNext();
}
function applySettings(value){
 if(!value)return;
 const readAloudChanged=settingsReady&&settings.tts_enabled!==value.tts_enabled;
 settings=value;settingsReady=true;applyTheme();document.documentElement.style.setProperty('--font-size',`${settings.font_size}px`);preferenceControls();
 // A Read Aloud transition relinquishes the old draft pause, not the draft.
 // Fresh input may claim a new pause; owned speech still releases through idle
 // and the existing epoch checks, so stale callbacks cannot resume new work.
 if(readAloudChanged){typing.delete(selected);inputOccupied=false;}
 if(!settings.tts_enabled)cancelAutomaticAudio();
 if(readAloudChanged&&selected)maybeResume(selected);
 $('settings-status').textContent='Saved · shared across your open tabs';
}
for(const [id,key] of Object.entries(preferenceFields))$(id).onchange=()=>{
 const control=$(id);let value=control.type==='checkbox'?control.checked:control.value;
 if(['font_size','tts_speed','max_tool_rounds','diagnostic_max_records'].includes(key))value=Number(value);
 if((typeof value==='number'&&!Number.isFinite(value))||(['max_tool_rounds','diagnostic_max_records'].includes(key)&&!Number.isSafeInteger(value))){$('settings-status').textContent='Enter a finite whole-number budget.';preferenceControls();return;}
 // A server-accepted snapshot owns rendering, never a speculative local edit.
 preferenceControls();if(send('update_settings',{settings:{[key]:value}}))$('settings-status').textContent='Saving…';
};
$('font-size').oninput=()=>{$('font-size-value').textContent=$('font-size').value;};
$('tts-speed').oninput=()=>{$('tts-speed-value').textContent=`${$('tts-speed').value}×`;};
const element = (tag, cls, text) => { const e = document.createElement(tag); if(cls) e.className=cls; if(text !== undefined) e.textContent=String(text); return e; };
function send(type, data={}) { if(socket?.readyState !== WebSocket.OPEN) { $('errors').textContent='Disconnected; wait for session replay before sending.'; return false; } socket.send(JSON.stringify({type,agent_id:selected,...data})); return true; }
function nearBottom(pane=$('conversation')){return pane.scrollHeight-pane.clientHeight-pane.scrollTop<100;}
function follow(wasNear){if(wasNear)requestAnimationFrame(()=>{$('conversation').scrollTop=$('conversation').scrollHeight;});else $('follow').hidden=false;}
$('follow').onclick=()=>{$('conversation').scrollTo({top:$('conversation').scrollHeight,behavior:'smooth'});$('follow').hidden=true;};
$('conversation').addEventListener('scroll',()=>{$('follow').hidden=nearBottom();},{passive:true});
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
 const destination=$(['tool_call','tool_result'].includes(c.type)?'tool-cards':'conversation');
 if(node.parentElement!==destination){destination.querySelector('.empty-state')?.remove();destination.append(node);}
}
function skillView(){
 const box=$('skills');if(!box)return;box.replaceChildren();const state=skillsByAgent.get(selected);if(!state){box.textContent='No skills configured.';return;}
 const loadable=new Set((state.loadable_skills||[]).map(s=>s.name));
 for(const s of state.skills||[]){
  const row=document.createElement('div'),label=document.createElement('span'),button=document.createElement('button');
  label.textContent=`${s.name} (${s.type}) — ${s.enabled?(s.pending_unload?'pending unload; still active':'enabled'):(loadable.has(s.name)?'loadable':'discovered')} ${s.description||''} `;
  row.append(label);if(s.enabled&&!s.pending_unload){button.textContent='Mark pending unload';button.onclick=()=>send('unload_skill',{name:s.name});row.append(button);}else if(loadable.has(s.name)||s.pending_unload&&s.type==='loadable'){button.textContent=s.pending_unload?'Keep loaded':'Load';button.onclick=()=>send('load_skill',{name:s.name});row.append(button);}box.append(row);
 }
 if(!(state.skills||[]).length)box.textContent='No skills configured.';
}
function apply(frame) {
 if(['part_final','part_delta','tool_dispatched','tool_finished'].includes(frame.type)) {
  const key=`${frame.agent_id}/${frame.part_id}`;let record=cards.get(key);
  if(frame.type==='part_delta') {
   if(record&&!record.frame.part.partial)return;
   frame={...frame,part:{type:frame.kind,name:frame.name,role:'agent',partial:true,text:(record?.frame.part.text||'')+(frame.chunk||'')}};
  }
  const dispatched=record?.frame.part.dispatched||frame.part?.dispatched;
  if(!record){record={frame,node:element('article','card')};cards.set(key,record);order.push(key);}
  record.frame=frame;record.frame.part.dispatched=!!dispatched;renderCard(record);
 } else if(frame.type==='job_update') {
  const key=`${frame.agent_id}/${frame.part_id}`;let record=cards.get(key);if(!record){record={frame,node:element('article','card')};cards.set(key,record);}else record.frame.part.job=frame.part.job;renderCard(record);
 } else if(frame.type==='part_removed'){const key=`${frame.agent_id}/${frame.part_id}`,r=cards.get(key);if(r?.frame.part.partial){r.node.remove();cards.delete(key);}
 } else if(frame.type==='state'||frame.type==='state_changed') {
  const a=agents.get(frame.agent_id), next=frame.type==='state'?frame.data:frame.state, previous=a?.state;if(a) a.state=next;
  if(frame.agent_id===selected) controls();
  // A state acknowledgement may follow newer input. Only explicit Continue
  // authority (below), not an automatic continuation epoch, clears typing.
  if(!frame.replay&&next.epoch!==previous?.epoch)cancelAudio(frame.agent_id,false);
 } else if(frame.type==='turn_ended') {for(const [key,r] of cards) {if(r.frame.agent_id===frame.agent_id&&r.frame.part.partial){r.node.remove();cards.delete(key);}}}
 else if(frame.type==='continued'&&!frame.replay){
  if(frame.data.client!==clientID&&frame.data.epoch===agents.get(frame.agent_id)?.state.epoch)explicitContinue(frame.agent_id);
 }else if(frame.type==='skill_state'){skillsByAgent.set(frame.agent_id,frame);skillView();
 }else if(frame.type==='current_settings'||frame.type==='settings_changed') applySettings(frame.settings);
 else if(frame.type==='settings_error') {$('settings-status').textContent=String(frame.data);preferenceControls();}
 else if(frame.type==='error') $('errors').textContent=String(frame.data);
 else if(frame.type==='speech_requested'&&!frame.replay){
  if(frame.data.epoch!==agents.get(frame.agent_id)?.state.epoch)return;
  requests.set(frame.data.id,{agent:frame.agent_id,epoch:frame.data.epoch});
  $('audio-status').textContent='Browser speech requested';send('claim_speech',{agent_id:frame.agent_id,speech_id:frame.data.id});
 }else if(frame.type==='speech_claimed'&&!frame.replay){
  const request=requests.get(frame.data.id);requests.delete(frame.data.id);
  if(frame.data.owner!==clientID||!request||request.epoch!==frame.data.epoch||frame.data.epoch!==agents.get(frame.agent_id)?.state.epoch)return;
  const task={agent:frame.agent_id,...frame.data};const occupied=voiceAgents.has(task.agent);voiceAgents.add(task.agent);if(!occupied)send('pause',{agent_id:task.agent});
  if(task.reason==='prompt'&&$('mic').checked&&settings.tts_enabled){const timer=setTimeout(()=>{delayed.delete(task.id);enqueueAudio(task);},2000);delayed.set(task.id,{task,timer});}else enqueueAudio(task);
 }else if(frame.type==='speech_cancelled'){
  cancelAudio(frame.agent_id,false);
  if(frame.data.release&&frame.data.owner===clientID)maybeResume(frame.agent_id);
 }else if(frame.type==='speech_idle'){
  if(frame.data.owner===clientID&&frame.data.epoch===agents.get(frame.agent_id)?.state.epoch){voiceAgents.delete(frame.agent_id);maybeResume(frame.agent_id);}
 }else if(frame.type==='speech_done'){requests.delete(frame.data.id);$('audio-status').textContent=`Speech: ${frame.data.status}`;}

}
function renderAgentTree(a){
 const node=$('agent-root');node.disabled=!a;if(!a)return;
 const state=a.state.paused?'paused':a.state.turn,label=state.replaceAll('_',' ');
 $('agent-name').textContent=`Agent ${a.id}`;$('agent-activity').textContent=label.charAt(0).toUpperCase()+label.slice(1);
 $('agent-cue').dataset.busy=String(state!=='idle');node.dataset.agentId=a.id;node.title=a.config.model;node.setAttribute('aria-selected',String(a.id===selected));
}
function controls(){skillView();const a=agents.get(selected);renderAgentTree(a);if(!a)return;const s=a.state;$('vendor').textContent=a.config.vendor;$('model').textContent=a.config.model;$('surface').textContent=a.config.surface;$('state').textContent=s.paused?'Paused':s.turn.replaceAll('_',' ');$('reason').textContent=s.reason||'';$('pause').disabled=s.paused;$('continue').disabled=!s.paused;$('audio-status').textContent=s.speaking?'Speech enabled':'Speech disabled';}
function replay(data){
 const previousAgents=new Map(agents);
 cancelAudio(undefined,false);clientID=data.client;session=data.session;cards.clear();order.length=0;$('conversation').replaceChildren();$('tool-cards').replaceChildren();agents.clear();skillsByAgent.clear();
 for(const a of data.agents){if(!a.state.paused||a.state.continue_epoch!==previousAgents.get(a.id)?.state.continue_epoch)typing.delete(a.id);agents.set(a.id,a);}
 if(!agents.has(selected))selected=data.agents[0]?.id||'';
 for(const frame of data.events) apply(frame);for(const frame of data.partials) apply(frame);controls();
}
function connect(){socket=new WebSocket(`${location.protocol==='https:'?'wss':'ws'}://${location.host}/ws`);
 socket.onopen=()=>{$('connection').textContent='Connected · synchronizing';retry=500;};
 socket.onmessage=event=>{let frame;try{frame=JSON.parse(event.data);}catch{return;}const near=nearBottom(),toolNear=nearBottom($('tool-cards'));if(frame.type==='replay'){replay(frame.data);$('connection').textContent='Connected';}else apply(frame);follow(near);if(toolNear)requestAnimationFrame(()=>{$('tool-cards').scrollTop=$('tool-cards').scrollHeight;});};
 socket.onclose=()=>{$('connection').textContent='Disconnected · reconnecting';settingsReady=false;preferenceControls();cancelAudio(undefined,false);setTimeout(connect,retry);retry=Math.min(retry*2,8000);};
 socket.onerror=()=>socket.close();
}
function busy(agent){return voiceAgents.has(agent)||[...requests.values()].some(r=>r.agent===agent)||speakingNow?.agent===agent||audio.some(t=>t.agent===agent)||[...delayed.values()].some(x=>x.task.agent===agent);}
function maybeResume(agent){if(!busy(agent)&&!typing.has(agent))send('unpause',{agent_id:agent,automatic:true});}
function explicitContinue(agent){typing.delete(agent);if(agent===selected)inputOccupied=false;cancelAudio(agent,false);}
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
 if(automaticSpeech(task)&&!settings.tts_enabled){complete('disabled');return;}
 if(!window.speechSynthesis||!window.SpeechSynthesisUtterance){complete('unsupported');return;}
 // Read every codepoint in order. No preview truncation or chunk dropping.
 const chunks=Array.from(task.text);let offset=0;
 const next=()=>{if(speakingNow!==task)return;if(offset>=chunks.length){complete('done');return;}const utterance=new SpeechSynthesisUtterance(chunks.slice(offset,offset+1200).join(''));offset+=1200;utterance.rate=Math.max(0.5,Math.min(3,settings.tts_speed));utterance.onend=next;utterance.onerror=()=>complete('failed');speechSynthesis.speak(utterance);};next();
}
$('agent-root').onclick=()=>{selected=$('agent-root').dataset.agentId;for(const r of cards.values())r.node.hidden=r.frame.agent_id!==selected;controls();};
function showSidebar(open){$('sidebar').hidden=!open;$('sidebar-resizer').hidden=!open;$('menu-toggle').setAttribute('aria-expanded',String(open));}
showSidebar(!window.matchMedia?.('(max-width: 1050px)')?.matches);
$('menu-toggle').onclick=()=>showSidebar($('sidebar').hidden);
function selectTab(name){for(const tab of ['settings','activity']){const active=tab===name;$(`${tab}-tab`).setAttribute('aria-selected',String(active));$(`${tab}-tab`).tabIndex=active?0:-1;$(`${tab}-panel`).hidden=!active;}}
for(const name of ['settings','activity']){
 $(`${name}-tab`).onclick=()=>selectTab(name);
 $(`${name}-tab`).onkeydown=event=>{if(['ArrowLeft','ArrowRight','Home','End'].includes(event.key)){event.preventDefault();const next=event.key==='Home'?'settings':event.key==='End'?'activity':name==='settings'?'activity':'settings';selectTab(next);$(`${next}-tab`).focus();}};
}
// Pane widths are transient UI state. Bounds leave useful room for conversation.
for(const [id,pane,property,minimum,maximum,direction] of [['sidebar-resizer','sidebar','--sidebar-width',220,420,1],['tools-resizer','tools-pane','--tools-width',260,700,-1]]){
 const handle=$(id);
 const resize=width=>{const other=$(pane==='sidebar'?'tools-pane':'sidebar');const available=$('workspace').clientWidth-(other.hidden?0:other.getBoundingClientRect().width)-334;const value=Math.max(minimum,Math.min(maximum,available,width));document.documentElement.style.setProperty(property,`${value}px`);handle.setAttribute('aria-valuenow',String(Math.round($(pane).getBoundingClientRect().width)));};
 handle.onkeydown=event=>{if(!['ArrowLeft','ArrowRight','Home','End'].includes(event.key))return;event.preventDefault();resize(event.key==='Home'?minimum:event.key==='End'?maximum:$(pane).getBoundingClientRect().width+(event.key==='ArrowRight'?20:-20)*direction);};
 handle.onpointerdown=event=>{if(event.button!==0)return;event.preventDefault();const start=event.clientX,width=$(pane).getBoundingClientRect().width;handle.setPointerCapture(event.pointerId);document.body.classList.add('resizing');handle.onpointermove=move=>resize(width+(move.clientX-start)*direction);};
 const stop=()=>{handle.onpointermove=null;document.body.classList.remove('resizing');};handle.onpointerup=stop;handle.onpointercancel=stop;handle.onlostpointercapture=stop;
}
$('pause').onclick=()=>send('pause');$('continue').onclick=()=>{explicitContinue(selected);send('unpause');};$('interrupt').onclick=()=>{cancelAudio(selected,false);send('interrupt');};
$('cancel-audio').onclick=()=>{cancelAudio(selected);send('update_settings',{settings:{tts_enabled:false}});};
function submit(type){const text=$('input').value;if(text.trim()&&send(type,{text})){$('input').value='';inputOccupied=false;typing.delete(selected);maybeResume(selected);$('errors').textContent='';}}
function draftChanged(){const nonempty=!!$('input').value.trim();if(nonempty&&!inputOccupied){typing.add(selected);send('pause');}else if(!nonempty&&typing.delete(selected))maybeResume(selected);inputOccupied=nonempty;}
$('input').oninput=draftChanged;
$('composer').onsubmit=e=>{e.preventDefault();submit('input');};$('hint').onclick=()=>submit('hint');$('input').onkeydown=e=>{if(e.key==='Enter'&&(e.ctrlKey||e.metaKey)){e.preventDefault();submit('input');}};
addEventListener('pagehide',()=>cancelAudio(undefined,false));connect();

$("save").onclick=()=>send("save");
