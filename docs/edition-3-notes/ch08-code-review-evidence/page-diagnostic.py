import sys,json,difflib
from pathlib import Path
sys.dont_write_bytecode=True
E=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch08-code-review-evidence');sys.path.insert(0,str(E))
from runner import run,T,sha
S=Path('/Users/bill/projects/ensemble-edition-3-student');C=T/'candidate-copy';original=S/'gui/test/correction-r1-page.mjs';text=original.read_text()
def one(old,new):
 global text
 assert text.count(old)==1,old
 text=text.replace(old,new)
one('class Node {',"const css=new Map([['--sidebar-width','240px'],['--actions-width','320px']]);\nclass Node {")
one('this.events=new Map();',"this.events=new Map();this.attrs={};this.style={setProperty:(k,v)=>css.set(k,v)};this.scrollHeight=100;this.clientHeight=100;this.scrollTop=0;")
one(' append(...nodes){'," setAttribute(k,v){this.attrs[k]=v;}\n focus(){this.focused=true;}\n getBoundingClientRect(){return {width:1100};}\n setPointerCapture(id){this.pointer=id;} hasPointerCapture(id){return this.pointer===id;} releasePointerCapture(){this.pointer=null;}\n append(...nodes){")
one("const ids=new Map();for(const id of ['transcript','state','paused','submit','draft','composer','notice','connection','cancel','interrupt','unpause','read-aloud','microphone'])ids.set(id,new Node(id==='draft'?'textarea':'div'));", "const ids=new Map();const html=await readFile(resolve(root,'index.html'),'utf8');for(const [,id] of html.matchAll(/id=\"([^\"]+)\"/g))ids.set(id,new Node(id==='draft'?'textarea':'div'));")
one('const document={hidden:true,',"const document={documentElement:new Node('html'),hidden:true,")
one('const window={speechSynthesis:',"const window={matchMedia:()=>({matches:false,addEventListener(){}}),speechSynthesis:")
one('vm.createContext({window,document,',"vm.createContext({window,document,getComputedStyle:()=>({getPropertyValue:k=>css.get(k)||'0px'}),")
one("const card=id=>ids.get('transcript').children.find(n=>String(n.dataset.partId)===String(id));", "const card=id=>[...ids.get('transcript').children,...ids.get('actions').children].find(n=>String(n.dataset.partId)===String(id));")
one("cancel:()=>{ids.get('read-aloud').checked=false;ids.get('read-aloud').onchange()}", "cancel:()=>deliver({type:'settings_changed',settings:{theme:'dark',tts_enabled:false,tts_speed:1,font_size:16,max_tool_rounds:0}})")
p=C/'gui/test/reviewer-page-diagnostic.mjs';p.write_text(text);(E/'reviewer-page-diagnostic.mjs').write_text(text)
(E/'page-fixture-diagnostic.diff').write_text(''.join(difflib.unified_diff(original.read_text().splitlines(True),text.splitlines(True),fromfile='submitted/correction-r1-page.mjs',tofile='owned/reviewer-page-diagnostic.mjs')))
(E/'page-diagnostic-identity.json').write_text(json.dumps({'original_fixture_path':str(original),'original_fixture_sha256':sha(original),'owned_fixture_path':str(p),'owned_fixture_sha256':sha(p),'production_app_sha256':sha(C/'gui/web/gui/app.js'),'candidate_production_app_sha256':sha(S/'gui/web/gui/app.js'),'scope':'Controlled test-boundary/pane-location/accepted-settings cancellation adaptation only; original behavioral assertions retained. Not a candidate fix or browser/native proof.'},indent=2)+'\n')
node=json.loads((E/'runtime-tools.json').read_text())['node'];run('page-adapted-diagnostic',[node,'--experimental-vm-modules',str(p)],C,30)
