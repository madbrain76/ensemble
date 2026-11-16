// Run from an isolated fixture workspace after starting --crossover-host. This is a
// real-browser consumer of the mounted, embedded GUI, not a fake WebSocket UI.
const { launchBrowser, launchPersistentProfile } = require('./browser-launch.cjs');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const crypto = require('node:crypto');
const path = require('node:path');
const root = process.cwd();
const evidence = path.join(root, '.crossover/evidence');
const base = 'http://127.0.0.1:18202/assistant/';
const check = (condition, message) => { assert.ok(condition, message); console.log('PASS '+message); };
// renderControls writes "<state> · <budget>" into #state, so a state comparison
// has to take the state word on its own, as job-status.cjs does.
// ea446b77 lets the server's own live Job count outrank the turn word, so an Agent
// with a background Job never reads plain "idle", and an explicitly paused Agent
// never reads plain "paused" either. The fixture keeps a 300-second Job running for
// the whole pass, so the turn word alone can neither settle nor prove a hold here.
// The Interrupt control is the discriminator: the actor enables it only while the
// turn is busy or paused, and a background Job does not enable it.
const settled = (pg) => pg.waitForFunction(() => {
  const word = document.querySelector('#state').textContent.split(' \u00b7 ')[0];
  return ['idle', 'inputpending'].includes(word) || /^[0-9]+ jobs? running$/.test(word);
});
const held = (pg) => pg.evaluate(() => {
  const word = document.querySelector('#state').textContent.split(' \u00b7 ')[0];
  return !['inflight', 'toolspending'].includes(word)
    && (word === 'paused' || !document.querySelector('#interrupt').disabled);
});
const waitHeld = (pg) => pg.waitForFunction(() => {
  const word = document.querySelector('#state').textContent.split(' \u00b7 ')[0];
  return !['inflight', 'toolspending'].includes(word)
    && (word === 'paused' || !document.querySelector('#interrupt').disabled);
});
const waitRunning = (pg) => pg.waitForFunction(() => {
  const word = document.querySelector('#state').textContent.split(' \u00b7 ')[0];
  return (['idle', 'inputpending'].includes(word) || /^[0-9]+ jobs? running$/.test(word))
    && document.querySelector('#interrupt').disabled && !document.querySelector('#model-choice').disabled;
});
(async () => {
 const context = await launchPersistentProfile(path.join(root,'.crossover/profile-browser'), {
  executablePath:'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome', headless:true,
  viewport:{width:1600,height:1000},
  // A fresh Agent now defaults to the "system" theme, so the dark evidence this
  // suite inspects has to come from an explicitly dark operating system rather
  // than from a stored preference. This also proves the default follows the OS.
  colorScheme:'dark',
  args:['--no-first-run','--disable-background-networking']
 });
 const page = await context.newPage();
 const fetchOutput = (url, options={}) => context.request.get(url, {
  ...options, headers: { 'X-Ensemble-Output':'1', ...options.headers }
 });
 const errors=[]; const external=[]; const requests=[]; const frames=[];
 const controls=[];
 // Check actual rendered colors for representative text and input boundaries.
 // This is scoped contrast evidence, not an automated accessibility verdict.
 const checkContrast=async(theme)=>{
  const report=await page.evaluate(async()=>{
   const rgb=s=>s.match(/[\d.]+/g).slice(0,3).map(Number);
   const luminance=c=>c.map(v=>{v/=255;return v<=.04045?v/12.92:((v+.055)/1.055)**2.4}).reduce((sum,v,i)=>sum+v*[.2126,.7152,.0722][i],0);
   const ratio=(a,b)=>{const x=[luminance(rgb(a)),luminance(rgb(b))].sort((a,b)=>a-b);return(x[1]+.05)/(x[0]+.05)};
   const background=e=>{while(e){const c=getComputedStyle(e).backgroundColor;if(c!=='rgba(0, 0, 0, 0)'&&c!=='transparent')return c;e=e.parentElement;}throw Error('No opaque background');};
   // #send is transparent artwork and paints no text, so its declared colour is
   // never rendered. The affordance is measured below, where it actually appears.
   const text=['#connection','#state','#agents-tab','#settings-open','#agent-tree .tree-state','#prompt','.card-heading strong','.card-heading .btn'].map(selector=>{
    const e=document.querySelector(selector);return{selector,ratio:ratio(getComputedStyle(e).color,background(e))};
   });
   const prompt=document.querySelector('#prompt'),style=getComputedStyle(prompt);
   const toggle=document.querySelector('.tree-expand').getBoundingClientRect();
   const artwork=await (async()=>{
    const icon=document.querySelector('#send .send-icon');
    const base=luminance(rgb(background(icon)));
    const bitmap=await createImageBitmap(icon);
    const width=bitmap.width,height=bitmap.height;
    const canvas=new OffscreenCanvas(width,height);
    const context=canvas.getContext('2d',{willReadFrequently:true});
    context.drawImage(bitmap,0,0);
    bitmap.close();
    const {data}=context.getImageData(0,0,width,height);
    // The theme darkens the pale artwork with a CSS filter. Canvas filtering is
    // unevenly supported, so the declared factor is applied to the pixels here.
    const shade=Number((getComputedStyle(icon).filter.match(/brightness\(([\d.]+)\)/)||[null,'1'])[1]);
    let best=0;
    for(let index=0;index<data.length;index+=4){
     if(data[index+3]<200)continue;
     const pixel=[data[index],data[index+1],data[index+2]].map(value=>value*shade);
     const pair=[luminance(pixel),base].sort((first,second)=>first-second);
     best=Math.max(best,(pair[1]+.05)/(pair[0]+.05));
    }
    return {ratio:best};
   })();
   const bubbles=['human','assistant'].map(role=>{
    const e=document.querySelector(`#transcript article[data-role="${role}"][data-kind="text"] > .message`),s=getComputedStyle(e);
    const link=document.createElement('a');link.href='#bubble-contrast-test';link.textContent='Link';e.append(link);
    const linkContrast=ratio(getComputedStyle(link).color,s.backgroundColor);link.remove();
    return {role,border:s.borderTopWidth,leftBorder:s.borderInlineStartWidth,borderStyle:s.borderTopStyle,borderColor:s.borderTopColor,background:s.backgroundColor,radius:s.borderTopLeftRadius,padding:s.padding,textContrast:ratio(s.color,s.backgroundColor),linkContrast,borderContrast:ratio(s.borderTopColor,s.backgroundColor)};
   });
   return{text,artwork,bubbles,inputBorder:ratio(style.borderTopColor,background(prompt)),target:{width:toggle.width,height:toggle.height}};
  });
  for(const item of report.text)check(item.ratio>=4.5,`${theme} text contrast ${item.selector}: ${item.ratio.toFixed(2)}:1`);
  check(report.artwork.ratio>=4.5,`${theme} send artwork is legible against the composer: ${report.artwork.ratio.toFixed(2)}:1`);
  check(report.inputBorder>=3,`${theme} input boundary contrast: ${report.inputBorder.toFixed(2)}:1`);
  check(report.target.width>=44&&report.target.height>=44,`${theme} tree expand target at least 44 pixels`);
  const [human,agent]=report.bubbles;
  // 87f5fcb3 removed the outline from every bubble and left the rounded fill with
  // the speaker name above it as the identity cue, but never updated these
  // expectations. What still has to hold is what that commit left behind: no
  // outline, a real radius, and fills that tell the two speakers apart.
  check(report.bubbles.every(b=>parseFloat(b.border)===0&&b.borderStyle==='none'&&parseFloat(b.radius)>0),`${theme} bubbles are unoutlined and rounded`);
  check(['radius','padding'].every(key=>human[key]===agent[key]),`${theme} both bubbles retain compact spacing`);
  // The heavier leading edge is gone; the non-color cue is now the speaker name
  // above the bubble, which this file checks with the accessible-label pass.
  check(agent.background!==human.background,`${theme} Agent has a fill distinct from the human bubble`);
  check(agent.textContrast>=4.5&&agent.linkContrast>=4.5,`${theme} Agent bubble text and links retain contrast`);
 };

 const modelRequests=[];
 // Observe real outgoing controls, without replacing the socket or the actor.
 page.on('websocket', ws=>ws.on('framesent', frame=>{
  try { const value=JSON.parse(frame.payload); if(value.type==='set_model')modelRequests.push(value); if(['pause','unpause','interrupt'].includes(value.type)) controls.push(value.type); } catch {}
 }));
 page.on('pageerror', e=>errors.push(e.stack));
 page.on('request', r=>{ requests.push(r.url()); if (!/^(http:\/\/127\.0\.0\.1:18202\/|data:)/.test(r.url())) external.push(r.url()); });
 page.on('websocket', ws=>ws.on('framereceived', f=>{try { frames.push(JSON.parse(f.payload)); } catch {} }));
 try {
  await page.goto(base);
  await page.waitForFunction(()=>document.querySelector('#connection').textContent==='Connected' && !document.querySelector('#send').disabled);
  check(await page.locator('html').getAttribute('data-bs-theme')==='dark'
   &&await page.locator('html').getAttribute('data-theme')==='system','fresh settings follow the system theme');
  await page.locator('#settings-open').click(); await page.locator('#settings-display-tab').click();
  await page.locator('#human-name').fill('Taylor <friend>');
  await page.locator('#human-description').fill('Writes $TOOLS <examples>');
  await page.locator('#user-settings button').click();
  await page.waitForFunction(()=>[...document.querySelectorAll('article[data-role="human"] .card-heading strong')].some(el=>el.textContent==='Taylor <friend>'));
  await page.reload();
  await page.waitForFunction(()=>document.querySelector('#connection').textContent==='Connected' && !document.querySelector('#send').disabled);
  await page.waitForFunction(()=>document.querySelector('#human-name').value==='Taylor <friend>');
  check(await page.locator('#human-description').inputValue()==='Writes $TOOLS <examples>','User description survives reload without interpretation');
  await page.waitForFunction(()=>document.querySelector('article[data-role="assistant"] .card-heading strong')?.textContent==='captain');
  const names=await page.evaluate(()=>[...document.querySelectorAll('#transcript article[data-kind="text"]')].map(el=>({
   name:el.querySelector('.card-heading strong').textContent, title:el.querySelector('.card-heading strong').title,
   speak:el.querySelector('[data-action="speak"]').getAttribute('aria-label'),
   above:el.querySelector('.card-heading').getBoundingClientRect().bottom<=el.querySelector('.message').getBoundingClientRect().top+1
  })));
  check(names.length>0&&names.every(x=>x.above&&!/\b(Assistant|You)\b/.test(x.name+x.title+x.speak)),'names appear above messages, including accessible labels');
  check(await page.locator('.card-heading friend').count()===0,'user name is plain text, not HTML');
  if (await page.locator('#settings-dialog').isVisible()) await page.locator('#settings-close').click();

  check(!(await page.locator('#narrate').isChecked()),'speech starts opt-in');
  await page.locator('.brand-mark').evaluate(image=>image.decode());
  const artwork=await page.locator('.brand-mark').evaluate(image=>({width:image.naturalWidth,height:image.naturalHeight,ratio:image.getBoundingClientRect().width/image.getBoundingClientRect().height}));
  check(artwork.width===128&&artwork.height===120&&Math.abs(artwork.ratio-128/120)<.01,'original proportional artwork loaded without crop');
  const favicon=await context.request.get(new URL('artwork/ensemble-32.png',base).href);
  check(favicon.ok()&&favicon.headers()['content-type']==='image/png','embedded PNG favicon has correct media type');
  check(await page.locator('#theme-toggle use').getAttribute('href')==='#icon-moon','dark theme uses moon with explicit switch action');
  check(await page.locator('#theme-toggle').getAttribute('aria-label')==='Switch to light theme','theme icon has accessible action');
  check(await page.locator('[data-action="speak"] svg').count()>0,'per-message Speak actions use locally defined icons');
  check(await page.locator('#stop-speech').getAttribute('aria-label')==='Stop speech','Stop speech remains explicitly named');
  await page.locator('#narrate').focus(); await page.keyboard.press('Space');
  await page.waitForFunction(()=>document.querySelector('#narrate').checked&&document.querySelector('.speech-toggle').title.startsWith('Speech on'));
  await page.keyboard.press('Space');
  await page.waitForFunction(()=>!document.querySelector('#narrate').checked&&document.querySelector('.speech-toggle').title.startsWith('Speech off'));
  check(true,'native speech icon toggles by keyboard and acknowledges saved on/off state');
  await page.locator('#send img').evaluate(image=>image.decode());
  const sendImage=await context.request.get(new URL('artwork/paper-airplane-send.svg',base).href);
  check(sendImage.ok()&&sendImage.headers()['content-type']==='image/svg+xml','supplied Send SVG loads with image media type');
  // 63fa7279 grouped the Send control inside #compose-actions and left this
  // expectation naming it as a direct child of the form, so the assertion could
  // not hold against upstream's own markup. The guarantee is unchanged: the form
  // shows exactly the input and one actions group, and Send is that group's
  // submit control on the same row as the input.
  const composer=await page.locator('#compose').evaluate(el=>{
   const visible=(node)=>[...node.children].filter(x=>x.getBoundingClientRect().height>0);
   const input=el.querySelector('textarea'),button=el.querySelector('button[type="submit"]');
   const a=input.getBoundingClientRect(),b=button.getBoundingClientRect();
   return {children:visible(el).map(x=>x.id),submit:button.id,
     belowInput:b.top>=a.bottom,placeholder:input.placeholder};
  });
  assert.deepEqual(composer.children,['prompt','compose-actions','voice-status']);
  assert.equal(composer.submit,'send');
  // d117cc03 moved the controls to their own row under a full-width input, which
  // composer-layout.cjs pins; the one-row wording here could not survive it.
  check(composer.belowInput,'composer is exactly input above one actions group with Send');
  check(composer.placeholder===`Ask ${await page.locator('#agent-choice').inputValue()} anything…`,'placeholder names the selected Agent');
  const measureComposer=()=>page.locator('#prompt').evaluate(el=>{const c=getComputedStyle(el);return{height:el.getBoundingClientRect().height,line:parseFloat(c.lineHeight),padding:parseFloat(c.paddingTop)+parseFloat(c.paddingBottom)+parseFloat(c.borderTopWidth)+parseFloat(c.borderBottomWidth),scroll:el.scrollHeight,client:el.clientHeight,overflow:c.overflowY};});
  const one=await measureComposer();
  check(Math.abs(one.height-one.line-one.padding)<2,'composer starts at one line');
  // resizeComposer fits rendered lines within the CSS one-to-five-line bounds, so
  // a multiline draft grows and a cleared one shrinks back. The earlier one-line
  // assertions predated that growth and could not both hold with it.
  await page.locator('#prompt').fill('first\nsecond');const two=await measureComposer();
  check(two.height>one.height&&Math.abs(two.height-2*two.line-two.padding)<2,'multiline input grows to two lines');
  await page.locator('#prompt').fill('a wrapped message '.repeat(100));const many=await measureComposer();
  check(many.height>two.height&&many.scroll>many.client&&many.overflow==='auto','long wrapped input caps and scrolls inside its area');
  await page.locator('#prompt').fill('');const empty=await measureComposer();
  check(Math.abs(empty.height-one.height)<2,'clearing input shrinks back to one line');
  await settled(page);

  // Structural presentation checks do not certify pixel-level beauty. These
  // protect the selected design's distinction between navigation and actions.
  await page.waitForFunction(()=>document.querySelector('#agent-choice').value==='captain');
  const design=await page.evaluate(()=>{
   const tab=document.querySelector('#agents-tab'), tree=document.querySelector('.tree-select');
   const t=getComputedStyle(tab), row=getComputedStyle(tree);
   return {
    tabButtonClass:tab.classList.contains('btn'),tabBottom:t.borderBottomColor,
    tabBackground:t.backgroundColor,tabBottomRadius:t.borderBottomLeftRadius,
    treeBorder:row.borderTopWidth,treeBackground:row.backgroundColor,
    sendPrimary:document.querySelector('#send').classList.contains('btn-primary'),
    // 7e46d046 removed the manual pause/play pair, which this file itself
    // asserts below, so only the surviving compact control can be named.
    controlNames:['interrupt'].map(id=>document.getElementById(id).getAttribute('aria-label'))};
  });
  check(await page.locator('#chat-pane > .pane-heading').count()===0,'redundant conversation title row removed');
  // 89c644e3 gave the Agents heading its own bottom rule in the same commit that
  // stopped it being a tab, so a border color can no longer equal a transparent
  // background. What the check still has to protect is the distinction it named:
  // this is navigation, not one of the action buttons beside it.
  check(!design.tabButtonClass && design.tabBottomRadius==='0px' && design.tabBackground==='rgba(0, 0, 0, 0)','tabs connect to their pane instead of appearing as action buttons');
  check(design.treeBorder==='0px' && design.treeBackground==='rgba(0, 0, 0, 0)','tree rows use quiet navigation styling');
  check(design.sendPrimary && design.controlNames.join(',')==='Interrupt','primary Send and named compact Agent controls');
  await checkContrast('dark');
  await page.screenshot({path:path.join(evidence,'design-dark-initial.png')});

  // Markdown is the shipped prose renderer module now, not a page global: this is
  // the same import renderers.js performs for the "prose" selection.
  const markdownShape=await page.evaluate(async ()=>{
   const target=document.createElement('div');
   target.className='part';document.body.append(target);
   const {default: markdown}=await import(new URL('./skills/shared/markdown.js',location.href).href);
   markdown(target,'- first\n- second\n\n1. one\n2. two\n\nA wrapped\nparagraph with **bold** and *emphasis*.');
   const shape={whiteSpace:getComputedStyle(target.querySelector('p')).whiteSpace,ul:target.querySelectorAll('ul').length,ol:target.querySelectorAll('ol').length,
    li:target.querySelectorAll('li').length,p:target.querySelectorAll('p').length,
    text:target.querySelector('p')?.textContent.replace(/\s+/g,' '),strong:target.querySelectorAll('strong').length,
    em:target.querySelectorAll('em').length};
   target.remove();return shape;
  });
  assert.deepEqual(markdownShape,{whiteSpace:'normal',ul:1,ol:1,li:4,p:1,text:'A wrapped paragraph with bold and emphasis.',strong:1,em:1});
  console.log('PASS semantic Markdown lists, wrapped paragraphs and inline emphasis');
  // Geometry-only stress: this does not test Engine accounting. Preserve and
  // restore the real values; later assertions test actual acknowledged usage.
  const geometry=await page.evaluate(()=>{
   const ids=['model-choice','narrate','interrupt','theme-toggle'];
   const measure=()=>Object.fromEntries(ids.map(id=>{
    const r=document.getElementById(id).getBoundingClientRect();return[id,[r.x,r.y,r.width,r.height]];
   }));
   const targets={state:'a-long-changing-state-value',connection:'Disconnected — reconnecting',
    'usage-cache':'123,456,789','usage-hit':'99.99%',
    'usage-all':'999,999,999','usage-cost':'$1234.56789'};
   const saved=Object.fromEntries(Object.keys(targets).map(id=>[id,document.getElementById(id).textContent]));
   const before=measure();
   for(const [id,text] of Object.entries(targets))document.getElementById(id).textContent=text;
   const after=measure();
   for(const [id,text] of Object.entries(saved))document.getElementById(id).textContent=text;
   return {before,after};
  });
  for(const id of Object.keys(geometry.before)) for(let i=0;i<4;i++)
   assert.ok(Math.abs(geometry.before[id][i]-geometry.after[id][i])<1,`${id} moved during status-only changes`);
  console.log('PASS status changes leave control rectangles stable');
  check(await page.locator('#set-model, #funding-source, #selected-name').count()===0,'header removes Apply, repeated Agent name and funding status');
  check(await page.locator('.model-caption').count()===0,'caption removed from header');
  const header=await page.evaluate(()=>({
   height:document.querySelector('#status-bar').getBoundingClientRect().height,
   groups:['.run-state','#model-controls','#usage-meter','.status-actions'].map(selector=>{const r=document.querySelector(selector).getBoundingClientRect();return {left:r.left,right:r.right,centerX:(r.left+r.right)/2,centerY:(r.top+r.bottom)/2};}),
   slots:[...document.querySelectorAll('#usage-meter output')].map(e=>({width:e.getBoundingClientRect().width,align:getComputedStyle(e).textAlign,gap:parseFloat(getComputedStyle(e.parentElement).columnGap)}))
  }));
  check(header.height<=55,'desktop header is at most 55px at the unchanged 18px font');
  check(header.groups.every(g=>Math.abs(g.centerY-header.groups[0].centerY)<1),'all header groups share one centered row');
  check(header.groups.every((g,i)=>!i||g.left>=header.groups[i-1].right),'status, model, usage and actions appear in order without overlap');
  check(Math.abs(header.groups[1].centerX-(page.viewportSize().width/2))<1,'model is centered in the entire header');
  check(Math.abs(header.groups[2].centerX-(header.groups[1].right+header.groups[3].left)/2)<1,'usage is centered between model and actions');
  check(header.slots.every(s=>s.width>0&&s.align==='left'&&s.gap<5),'usage figures occupy fixed left-aligned slots beside labels');
  const usageFits=await page.evaluate(()=>{
   const values={'usage-cost':'$1234.56','usage-all':'1,234,567','usage-cache':'1,234,567','usage-hit':'100.0%'};
   return Object.entries(values).every(([id,value])=>{const e=document.getElementById(id),saved=e.textContent;e.textContent=value;const fits=e.scrollWidth<=e.clientWidth+1;e.textContent=saved;return fits;});
  });
  check(usageFits,'tight numeric slots retain seven-digit counts, cost and full hit rate');
  const originalModel=await page.locator('#model-choice').inputValue();
  // showUsage lists every configured offering and disables the ones the current
  // provider will not serve, so "another model" has to mean another *enabled*
  // one, which is what the check below already claims to require.
  const otherModel=await page.locator('#model-choice option').evaluateAll((options,current)=>options.find(o=>o.value!==current&&!o.disabled)?.value,originalModel);
  check(!!otherModel,'fixture supplies another same-provider model');
  const startRequests=modelRequests.length;
  const startFrames=frames.length;
  await page.locator('#model-choice').selectOption(otherModel);
  await page.waitForFunction(model=>!document.querySelector('#model-choice').disabled&&document.querySelector('#model-choice').value===model,otherModel);
  check(modelRequests.length===startRequests+1&&modelRequests.at(-1).text===otherModel,'one selection immediately sends one model change');
  check(frames.slice(startFrames).some(f=>f.type==='usage'&&f.selected_model===otherModel),'server acknowledges changed model');
  // showUsage() rebuilds the picker on every usage frame, so an option added in one
  // round trip can be gone before a second one selects it. Adding it, selecting it and
  // firing the page's own change handler in one evaluate keeps the refusal the check
  // is about without depending on that race going the right way.
  await page.locator('#model-choice').evaluate(el=>{
      el.add(new Option('Invalid empty model',' '));
      el.value=' ';
      el.dispatchEvent(new Event('change',{bubbles:true}));
  });
  await page.waitForFunction(model=>!document.querySelector('#model-choice').disabled&&document.querySelector('#model-choice').value===model,otherModel);
  check((await page.locator('#notice').textContent()).includes('Model selection refused'),'refused model restores acknowledged choice');
  console.log('PASS immediate model selection, acknowledgment and refusal recovery');

  // The public agents_status regression adds one real result to the six seeded artifacts.
  // The fixture gained a real work-product offer and further tool observations, so
  // a literal card count no longer describes it. What this check exists to protect
  // is the split itself: every Tools card is an observation or its media, and the
  // conversation pane below carries none of them.
  const toolKinds=await page.evaluate(()=>[...document.querySelectorAll('#tool-transcript .part')]
   .map(node=>node.dataset.kind||node.querySelector('.card-heading')?.dataset.kind||''));
  check(toolKinds.includes('tool_call')&&toolKinds.filter(kind=>kind==='tool_result').length>=4&&toolKinds.includes('blob'),
   'tool calls, results and media appear in Tools');
  check(toolKinds.every(kind=>['tool_call','tool_result','blob',''].includes(kind)),'Tools holds no conversation prose');
  check(await page.locator('#transcript [data-kind="tool_call"],#transcript [data-kind="tool_result"]').count()===0,'Chat excludes tools');
  check(await page.locator('#tool-transcript img').evaluate(e=>e.complete && e.naturalWidth>0),'embedded image renders');
  const originalIDs=await page.locator('[data-part-id]').evaluateAll(nodes=>nodes.map(n=>n.dataset.partId));
  check(new Set(originalIDs).size===originalIDs.length,'replay IDs are unique across both panes');
  // Bootstrap owns settings tabs; native dialog owns modal focus/Escape.
  await page.locator('#settings-open').click();
  await page.locator('#settings-ai-tab').focus(); await page.keyboard.press('ArrowRight');
  check(await page.locator('#settings-accounts-tab').getAttribute('aria-selected')==='true','tab ArrowRight selects Accounts');
  check(await page.locator('#settings-ai-panel').isHidden(),'inactive settings panel is hidden');
  await page.keyboard.press('ArrowLeft');
  check(await page.locator('#settings-ai-tab').getAttribute('aria-selected')==='true','tab ArrowLeft returns to AI');
  await page.keyboard.press('Escape');
  check(await page.evaluate(()=>document.activeElement.id)==='settings-open','closing settings restores gear focus');
  await page.locator('[data-agent="captain"]').focus(); await page.keyboard.press('ArrowRight');
  check(await page.evaluate(()=>document.activeElement.dataset.agent)==='captain/worker','tree ArrowRight moves to child');
  await page.keyboard.press('ArrowRight');
  check(await page.evaluate(()=>document.activeElement.dataset.agent)==='captain/worker/nested','tree reaches grandchild');
  await page.keyboard.press('ArrowLeft');
  check(await page.evaluate(()=>document.activeElement.dataset.agent)==='captain/worker','tree ArrowLeft returns to parent');
  const treeGeometry=await page.evaluate(()=>{
    const row=key=>{const e=document.querySelector(`[data-agent="${key}"]`);return{left:e.querySelector('.tree-name').getBoundingClientRect().left,level:Number(e.closest('[role="treeitem"]').getAttribute('aria-level'))};};
    return{root:row('captain'),peer:row('peer'),child:row('captain/worker'),grandchild:row('captain/worker/nested'),retained:row('captain/retained')};
  });
  check(Math.abs(treeGeometry.root.left-treeGeometry.peer.left)<1,'leaf and branch roots align');
  check(treeGeometry.child.left-treeGeometry.root.left>=18&&treeGeometry.grandchild.left-treeGeometry.child.left>=18,'child and grandchild have clear indentation');
  check(treeGeometry.child.level===2&&treeGeometry.grandchild.level===3,'tree exposes nesting levels');
  check(await page.locator('#agent-tree [data-delete-agent]').count()===0,'one action set replaces the per-row trash buttons');
  check(await page.locator('#agent-actions button[data-delete-agent]').count()===0,'mounted root has no trash action');
  check(await page.locator('#agent-actions button[aria-label="Delete selected Agent"]').isDisabled(),'the shared delete control is inert for the mounted root');
  check(await page.locator('#peer-delete,#peer-delete-confirm').count()===0,'standalone delete controls removed');
  check(await page.locator('#agent-choice').isHidden(),'tree is the only visible conversation selector');
  check((await page.locator('[data-agent="captain/retained"] .tree-state').innerText()).includes('submitted'),'submitted child remains visible with result status');
  await page.locator('[data-toggle-agent="captain/worker"]').click();
  check(await page.locator('[data-agent="captain/worker/nested"]').isHidden(),'branch collapse hides descendants');
  const protectedDelete=await page.evaluate(async()=>{const r=await fetch('./peers/delete',{method:'POST',headers:{'Content-Type':'application/json','X-Ensemble-Peers':'1'},body:JSON.stringify({name:'captain'})});return r.status;});
  check(protectedDelete===409,'server independently refuses mounted-root deletion');
  await page.locator('#peer-open').click();
  check(await page.locator('#peer-dialog').evaluate(e=>e.open),'Create Agent opens dialog');
  await page.locator('#peer-name').fill('tree-created');
  await page.locator('#peer-create').click();
  await page.waitForFunction(()=>document.querySelector('[data-agent="tree-created"]'));
  check(await page.locator('#peer-dialog').isHidden(),'creation closes dialog while tree remains available');
  check(await page.locator('[data-agent="captain/worker/nested"]').isHidden(),'collapse survives a new server tree snapshot');
  await page.locator('[data-toggle-agent="captain/worker"]').click();
  check(await page.locator('[data-agent="captain/worker/nested"]').isVisible(),'branch re-expands');
  await page.locator('[data-agent="captain/retained"]').click();
  await page.waitForFunction(()=>document.querySelector('#agent-choice').value==='captain/retained'&&!document.querySelector('#send').disabled);
  // The shared control targets whichever Agent is selected, and still confirms.
  const trash=()=>page.locator('#agent-actions button[data-delete-agent]');
  await page.waitForFunction(()=>document.querySelector('#agent-actions button[data-delete-agent="captain/retained"]'));
  page.once('dialog',dialog=>dialog.dismiss());await trash().click();
  check(await page.locator('[data-agent="captain/retained"]').count()===1,'canceling deletion preserves child');
  page.once('dialog',dialog=>dialog.accept());await trash().click();
  await page.waitForFunction(()=>!document.querySelector('[data-agent="captain/retained"]')&&document.querySelector('#agent-choice').value==='captain');
  check((await page.locator('#peer-notice').innerText()).includes('saved files retained'),'child deletion acknowledges retained files');
  check(fs.existsSync(JSON.parse(fs.readFileSync(path.join(root,'.crossover/runtime/fixture.json'),'utf8')).retainedArtifact),'child deletion actually retains submit.md');
  await page.locator('[data-agent="tree-created"]').click();
  await page.waitForFunction(()=>document.querySelector('#agent-actions button[data-delete-agent="tree-created"]'));
  page.once('dialog',dialog=>dialog.accept());await page.locator('#agent-actions button[data-delete-agent="tree-created"]').click();
  await page.waitForFunction(()=>!document.querySelector('[data-agent="tree-created"]'));
  console.log('PASS real tree hierarchy, collapse persistence, retained submission, creation dialog, and guarded row deletion');
  await page.locator('#left-splitter').focus(); const before=await page.locator('#left-splitter').getAttribute('aria-valuenow');
  await page.keyboard.press('ArrowRight');
  check(await page.locator('#left-splitter').getAttribute('aria-valuenow')!==before,'separator supports keyboard resizing');
  const savedWidths=await page.evaluate(()=>['left','right'].map(side=>document.documentElement.style.getPropertyValue(`--${side}-width`)));
  await page.locator('#left-splitter').press('Home');await page.locator('#right-splitter').press('End');
  const paneGeometry=()=>page.evaluate(()=>Object.fromEntries(['workspace','left-pane','chat-pane','tools-pane'].map(id=>[id,document.getElementById(id).getBoundingClientRect().width])));
  let panes=await paneGeometry();
  check(panes['tools-pane']>panes.workspace*.7&&panes['chat-pane']>=159&&panes['chat-pane']<170,'Actions can occupy most of the workspace while chat is narrow');
  await page.setViewportSize({width:1200,height:1000});
  await page.waitForFunction(()=>document.querySelector('#chat-pane').getBoundingClientRect().width>=159);
  panes=await paneGeometry();check(['left-pane','chat-pane','tools-pane'].every(id=>panes[id]>=159),'window resizing preserves reachable minimum panes');
  await page.setViewportSize({width:1600,height:1000});
  const handle=await page.locator('#right-splitter').boundingBox();
  await page.mouse.move(handle.x+handle.width/2,handle.y+20);await page.mouse.down();await page.mouse.move(1595,handle.y+20,{steps:8});await page.mouse.up();
  panes=await paneGeometry();check(panes['chat-pane']>panes.workspace*.7,'mouse dragging can give the conversation most of the workspace');
  check(await page.locator('#transcript article[data-kind="text"]').first().evaluate(el=>el.getBoundingClientRect().width)>900,'messages use the available wide conversation instead of a fixed content cap');
  await page.evaluate(values=>['left','right'].forEach((side,i)=>{const style=document.documentElement.style;values[i]?style.setProperty(`--${side}-width`,values[i]):style.removeProperty(`--${side}-width`);}),savedWidths);

  // Full output is owned by an Agent and Job, not by the browser's displayed path.
  const fixture=JSON.parse(fs.readFileSync('.crossover/runtime/fixture.json'));
  const output=await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.handle}`);
  const complete=await output.text();
  check(output.status()===200 && complete.includes('fixture-line-1999') && complete.length>60000,'complete output extends far beyond bounded preview');
  check(!output.headers()['content-type'].includes('html') && output.headers()['x-content-type-options']==='nosniff','output is passive text');
  const downloaded=await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.handle}&download=1`);
  check(downloaded.headers()['content-disposition'].includes('attachment') && await downloaded.text()===complete,'download returns complete bytes');
  for (const [query,status] of [
   [`agent=peer&handle=${fixture.handle}`,404], [`agent=captain/worker&handle=${fixture.handle}`,404],
   [`agent=removed&handle=${fixture.removedHandle}`,404], [`agent=captain&handle=999999`,404],
   [`agent=captain&handle=${fixture.handle}&path=/etc/passwd`,400], [`agent=captain&handle=${fixture.handle}&agent=peer`,400]
  ]) check((await fetchOutput(base+'jobs/output?'+query)).status()===status,'output refusal '+query);
  check((await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.handle}`,{headers:{Origin:'https://foreign.invalid'}})).status()===403,'foreign Origin refused');
  check((await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.handle}`,{headers:{'Sec-Fetch-Site':'cross-site'}})).status()===403,'cross-site navigation refused');
  const running=await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.runningHandle}`);
  check(running.status()===409,'running output explicitly refuses as incomplete');
  const nested=await fetchOutput(base+`jobs/output?agent=captain/worker&handle=${fixture.nestedHandle}`);
  check(nested.status()===200 && await nested.text()==='nested-owned-output\nexit_code: 0\n','nested owner address resolves its own output');
  check((await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.nestedHandle}`)).status()===404,'parent cannot claim child Job by handle');
  check((await context.request.get(base+`jobs/output?agent=captain&handle=${fixture.handle}`)).status()===403,'missing custom header refuses');
  for(const [agent,handle,label] of [['files',1,'unreadable'],['files',2,'FIFO'],['files',3,'missing file'],['confined',1,'rooted symlink escape']]) {
   const refusal=await fetchOutput(base+`jobs/output?agent=${agent}&handle=${handle}`, {timeout:3000});
   check(refusal.status()===404 && !(await refusal.text()).includes('must-not-cross-root'),label+' backing file refuses without blocking or revealing content');
  }

  check((await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.oversizedHandle}`)).status()===413,'oversized output explicitly refuses');
  check(crypto.createHash('sha256').update(await output.body()).digest('hex').toUpperCase()===fixture.outputHash,'complete response exact bytes match owned output');
  // The display window retires older cards, so the download button left on screen
  // belongs to whichever Job the window kept rather than to the first one. The
  // check owns a lossless guarded download, so it compares the saved file with
  // what the endpoint serves for the very Job the button asked about, and still
  // insists that Job is one this fixture owns.
  await page.evaluate(() => {
      window.downloadedOutputURL = null;
      const real = window.fetch;
      window.fetch = async (input, init) => {
          const response = await real(input, init);
          if (String(input).includes('jobs/output')) window.downloadedOutputURL = String(input);
          return response;
      };
  });
  const downloadEvent=page.waitForEvent('download');
  await page.getByRole('button',{name:'Download complete output',exact:true}).first().click();
  const browserDownload=await downloadEvent;
  await browserDownload.saveAs(path.join(evidence,'downloaded-output.txt'));
  const downloadedURL = await page.evaluate(() => window.downloadedOutputURL);
  const downloadedHandle = Number(new URL(downloadedURL).searchParams.get('handle'));
  const servedBytes = await (await fetchOutput(downloadedURL)).body();
  const servedHash = crypto.createHash('sha256').update(servedBytes).digest('hex').toUpperCase();
  const savedHash = crypto.createHash('sha256')
      .update(fs.readFileSync(path.join(evidence,'downloaded-output.txt'))).digest('hex').toUpperCase();
  // The download must be a guarded read of this Agent's own output, and it must
  // arrive without a byte changing. Which Job that is follows the card the window
  // kept, so it is taken from the request the page made rather than assumed.
  check(savedHash === servedHash && new URL(downloadedURL).pathname.endsWith('/jobs/output')
      && new URL(downloadedURL).searchParams.get('agent') === 'captain',
      'guarded browser download exact bytes');
  const cursorOutput=await fetchOutput(base+`jobs/output?agent=captain&handle=${fixture.cursorHandle}`);
  check(cursorOutput.status()===200 && (await cursorOutput.text()).includes('unobserved-suffix'),'download completed two-stage output before ordinary observation');
  const observed=await context.request.get('http://127.0.0.1:18202/fixture/observe');
  check((await observed.text()).includes('unobserved-suffix'),'later ordinary Job observer still receives unconsumed tail');
  // Copy controls are a shipped feature now: the prose renderer puts "Copy code"
  // on a closed fence, which markdown-copy.cjs pins, and a snippet box carries its
  // own. What still has to hold is the boundary that check was drawn around — a
  // control belongs to code, never to spoken message prose.
  const copyScopes=await page.evaluate(()=>[...document.querySelectorAll('#transcript button,#tool-transcript button')]
   .filter(node=>/^Copy(?: code| text)?$/.test(node.textContent.trim()))
   .map(node=>Boolean(node.closest('.code-snippet,.snippet-controls'))));
  check(copyScopes.every(Boolean),'Copy controls belong to code, never to message prose');
  check((await page.locator('details pre code').allTextContents()).some(text=>text.includes('value65')),'complete code remains available for expansion and selection');
  await page.locator('#transcript details').first().locator('summary').click();
  check(await page.locator('#transcript .code-keyword').count()>0,'expanded code has safe token highlighting');
  await page.screenshot({path:path.join(evidence,'desktop-dark.png')});
  await page.locator('[data-agent="peer"]').click();
  await page.waitForFunction(()=>document.querySelector('#transcript').textContent.includes('PEER-ONLY'));
  check(!(await page.locator('#transcript').innerText()).includes('Offline crossover'),'selection replaces Chat replay');
  check(await page.locator('#tool-transcript .part').count()===0,'selection clears old tool display');
  check(frames.filter(frame=>frame.type==='connection').at(-1)?.relay_agent_id==='captain','selection leaves the fixed relay authority unchanged');
  await page.locator('[data-agent="captain"]').click();
  await page.waitForFunction(()=>document.querySelector('#transcript').textContent.includes('Offline crossover'));
  await page.locator('#prompt').fill('Stream the fixture response'); await page.locator('#send').click();
  await page.waitForFunction(()=>document.querySelector('#transcript').textContent.includes('STREAM-IDENTITY first'));
  check(await page.locator('#model-choice').isDisabled(),'model mutation disabled while running');
  check(await page.locator('#send').isEnabled(),'ordinary input remains available during a turn');
  check(await page.locator('#interrupt').isEnabled(),'interrupt enabled while running');
  const streamed=await page.locator('#transcript .part').filter({hasText:'STREAM-IDENTITY first'}).getAttribute('data-part-id');
  await page.waitForFunction(()=>document.querySelector('#transcript').textContent.includes('and final.'));
  await page.waitForFunction(()=>!document.querySelector('#model-choice').disabled);
  check(await page.locator('#transcript .part').filter({hasText:'STREAM-IDENTITY first and final.'}).count()===1,'stream final updates one card');
  check(await page.locator('#transcript .part').filter({hasText:'STREAM-IDENTITY first and final.'}).getAttribute('data-part-id')===streamed,'stream/final part identity stable');
  check((await page.locator('#usage-all').textContent())==='31','meter uses acknowledged Engine input');
  const paused=()=>waitHeld(page);
  const resumed=()=>waitRunning(page);
  const barrier=async()=>{
   await page.evaluate(()=>{window.priorSettings=settings;send('get_settings');});
   await page.waitForFunction(()=>settings!==window.priorSettings);
  };
  // Explicit runtime pauses remain supported, but the GUI has no manual pair.
  check(await page.locator('#pause,#resume').count()===0,'manual pause/play pair removed');
  check(await page.locator('#speech-controls #stop-speech').count()===1,'Stop speech is reachable beside mic');
  // 3f7879a7 replaced connection-owned holds with one Agent-wide pause, where
  // Stop speech and Escape explicitly resume and a cleared draft resumes too.
  // Those edges are automatic-pause.cjs's contract; what belongs here is the
  // draft the operator keeps while the Agent stays held for a real reason.
  await page.locator('#prompt').fill('interrupt retains draft hold');await paused();
  await page.locator('#interrupt').click();await barrier();
  check(await held(page),'Interrupt does not erase typed draft hold');
  await page.locator('#prompt').fill('');await resumed();
  const beforeReconnect=await page.locator('#transcript .part').evaluateAll(nodes=>nodes.map(n=>n.dataset.partId));
  // Disconnect/reconnect must rebuild the same canonical cards, not append them.
  await context.setOffline(true); await page.waitForTimeout(300); await context.setOffline(false);
  await page.waitForFunction(()=>document.querySelector('#connection').textContent==='Connected' && !document.querySelector('#send').disabled);
  await page.reload(); await page.waitForFunction(()=>!document.querySelector('#send').disabled);
  assert.deepEqual(await page.locator('#transcript .part').evaluateAll(nodes=>nodes.map(n=>n.dataset.partId)),beforeReconnect);
  console.log('PASS reconnect/reload preserves canonical card identities');
  // Acknowledged theme/font are durable, then deliberately tested together.
  await page.locator('#theme-toggle').click(); await page.waitForFunction(()=>document.documentElement.dataset.bsTheme==='light');
  check(await page.locator('#theme-toggle use').getAttribute('href')==='#icon-sun','light theme uses sun');
  await checkContrast('light');
  await page.screenshot({path:path.join(evidence,'desktop-light.png')});
  await page.locator('#settings-open').click(); await page.locator('#settings-display-tab').click(); await page.locator('#font-size').fill('32'); await page.locator('#font-size').dispatchEvent('change');
  await page.waitForFunction(()=>getComputedStyle(document.documentElement).fontSize==='32px');
  // 784e841d moved font size into the browser profile, so the durable record this
  // choice has to survive in is localStorage, not the server's appearance file.
  check(JSON.parse(await page.evaluate(()=>localStorage.getItem('ensemble.device.v1'))).font===32,
    'font size durable in this browser profile');
  await page.locator('#settings-close').click();
  await page.setViewportSize({width:390,height:844});
  await page.waitForFunction(()=>document.body.classList.contains('compact'));
  check(await page.locator('#pane-navigation').isVisible(),'pane navigation remains reachable at narrow largest text');
  for (const pane of ['agents','chat','tools']) {
   await page.locator(`button[data-pane="${pane}"]`).click();
   check(await page.locator('#workspace').getAttribute('data-pane')===pane,'narrow pane switch '+pane);
   check(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth+1),'no page horizontal overflow in '+pane);
   if (pane==='chat') {
    const bounds=await page.locator('#compose').boundingBox();
    check(bounds.height>=44 && bounds.y+bounds.height<=844,'composer remains on-screen at largest text/narrow width');
    await page.locator('#prompt').fill('narrow-window text');
    check(await page.locator('#send').isVisible(),'Send remains reachable');
   }
   await page.screenshot({path:path.join(evidence,`narrow-32-${pane}.png`)});
  }
  await page.reload(); await page.waitForFunction(()=>!document.querySelector('#theme-toggle').disabled);
  check(await page.locator('html').getAttribute('data-bs-theme')==='light','light preference survives reload');
  check(await page.evaluate(()=>getComputedStyle(document.documentElement).fontSize)==='32px','largest font survives reload');
  await page.locator('button[data-pane="agents"]').click(); await page.locator('#settings-open').click(); await page.locator('#settings-display-tab').click();
  await page.locator('#font-size').fill('18'); await page.locator('#font-size').dispatchEvent('change');
  await page.waitForFunction(()=>getComputedStyle(document.documentElement).fontSize==='18px');
  await page.locator('#settings-close').click();
  await page.locator('#theme-toggle').click(); await page.waitForFunction(()=>document.documentElement.dataset.bsTheme==='dark');
  // Cross enlarged text with narrow/intermediate windows. Manual textarea
  // sizing must not push Send outside the viewport. Scrollable header groups
  // must permit reaching each control without page-wide horizontal overflow.
  for(const size of [18,32]) {
   await page.locator('button[data-pane="agents"]').click();await page.locator('#settings-open').click(); await page.locator('#settings-display-tab').click();
   await page.locator('#font-size').fill(String(size));await page.locator('#font-size').dispatchEvent('change');
   await page.waitForFunction(value=>getComputedStyle(document.documentElement).fontSize===`${value}px`,size);
  await page.locator('#settings-close').click();
   for(const [width,height] of [[390,844],[390,1000],[768,1000],[1536,1000]]) {
    await page.setViewportSize({width,height});
    await page.waitForFunction(()=>document.body.classList.contains('compact')===(innerWidth<62*parseFloat(getComputedStyle(document.documentElement).fontSize)));
    if(await page.locator('#pane-navigation').isVisible())await page.locator('button[data-pane="chat"]').click();
    await page.locator('#prompt').evaluate(el=>el.style.height='1000px');
    const layout=await page.evaluate(()=>{
     const ids=['model-choice','narrate','interrupt','theme-toggle'];
     const measure=()=>ids.map(id=>{
      const r=document.getElementById(id).getBoundingClientRect();return[r.x,r.y,r.width,r.height];
     });
     const saved=['state','usage-all','usage-cost'].map(id=>[id,document.getElementById(id).textContent]);
     const before=measure();
     for(const [id] of saved)document.getElementById(id).textContent='Long changing status: 123,456,789.12345';
     const after=measure();
     for(const [id,text] of saved)document.getElementById(id).textContent=text;
     const send=document.querySelector('#send').getBoundingClientRect();
     return {before,after,overflow:document.documentElement.scrollWidth>innerWidth+1,send:{top:send.top,bottom:send.bottom,right:send.right}};
    });
    assert.equal(layout.overflow,false,`page overflow ${width}/${size}`);
    assert.ok(layout.send.top>=0 && layout.send.bottom<=height && layout.send.right<=width,`Send unreachable ${width}x${height}/${size}: ${JSON.stringify(layout.send)}`);
    for(let i=0;i<layout.before.length;i++)for(let j=0;j<4;j++)assert.ok(Math.abs(layout.before[i][j]-layout.after[i][j])<1,`status moved control ${i} at ${width}/${size}`);
    // 7e46d046 removed the manual pause/play pair from the header.
    for(const id of ['model-choice','narrate','interrupt','theme-toggle']) {
     await page.locator('#'+id).scrollIntoViewIfNeeded();
     assert.ok(await page.locator('#'+id).evaluate(el=>{
      const r=el.getBoundingClientRect(),bar=document.querySelector('#status-bar').getBoundingClientRect();
      return r.left>=0 && r.right<=innerWidth+1 && r.top>=bar.top-1 && r.bottom<=bar.bottom+1;
     }),`status control ${id} clipped at ${width}/${size}`);
    }
    console.log(`PASS stable header, reachable controls/composer and no page overflow at ${width}x${height}px / font ${size}`);
    await page.locator('#prompt').evaluate(el=>el.style.height='');
   }
   await page.setViewportSize({width:390,height:844});
   await page.waitForFunction(()=>document.body.classList.contains('compact'));
  }
  check(external.length===0,'no external page requests'); check(errors.length===0,'no browser script errors');
  fs.writeFileSync(path.join(evidence,'browser-network.json'),JSON.stringify({requests,external,errors,frameTypes:[...new Set(frames.map(f=>f.type))]},null,2));
  console.log('PASS crossover browser checks complete');
 } catch(error) {
  await page.screenshot({path:path.join(evidence,'browser-failure.png')}).catch(()=>{});
  fs.writeFileSync(path.join(evidence,'browser-failure.json'),JSON.stringify({error:String(error),errors,external,controls,
   ui:await page.evaluate(()=>({state:document.querySelector('#state').textContent,
    connection:document.querySelector('#connection').textContent,
    stopSpeech:!!document.querySelector('#stop-speech'),draft:document.querySelector('#prompt').value.length})).catch(()=>null),
   frames:frames.slice(-30)},null,2));
  throw error;
 } finally { await context.close(); }
})().catch(e=>{console.error(e);process.exitCode=1;});
