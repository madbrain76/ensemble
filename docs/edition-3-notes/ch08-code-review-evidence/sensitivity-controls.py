import sys,json,hashlib,subprocess,shutil,difflib
from pathlib import Path
sys.dont_write_bytecode=True
E=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch08-code-review-evidence');sys.path.insert(0,str(E))
from runner import run,T,ENV,sha
S=Path('/Users/bill/projects/ensemble-edition-3-student');C=T/'candidate-copy';assert not C.exists();C.mkdir()
files=subprocess.check_output(['git','ls-files'],cwd=S).decode().splitlines();copy=[]
for q in files:
 p=C/q;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes((S/q).read_bytes());assert sha(p)==sha(S/q);copy.append({'path':q,'sha256':sha(p)})
(E/'owned-copy-before.json').write_text(json.dumps(copy,indent=2)+'\n')
# Reuse only reviewed transport framing and process cleanup; add independent
# assertions for rejected partial patches, a real filesystem save failure, and
# immutable prompt/absence of settings as user context.
text=(S/'gui/test/settings_process.py').read_text()
text=text.replace('rounds, requests = {}, []','rounds, requests = {}, []\nsystem_hashes = set()')
anchor="        label = next((part['text']"
assert text.count(anchor)==1
extra="""        system_hashes.add(hashlib.sha256(json.dumps(req.get('system'),sort_keys=True).encode()).hexdigest())
        assert len(system_hashes)==1, 'settings changed base/system prompt'
        for msg in req['messages']:
            if msg['role'] != 'user': continue
            content = msg['content'] if isinstance(msg['content'],list) else [{'type':'text','text':msg['content']}]
            for part in content:
                if part.get('type') == 'text':
                    assert part.get('text') in ('FIRST_BUDGET','LATER_18','ZERO_20'), 'non-turn/settings context leaked into user text'
"""
text=text.replace(anchor,extra+anchor)
anchor="    a.until('replay_end'); b.until('replay_end')\n    process.stdin.write('FIRST_BUDGET\\n'); process.stdin.flush()";assert text.count(anchor)==1
extra="""    a.until('replay_end'); b.until('replay_end')
    durable_before = settings.read_bytes()
    # A valid sibling must not partially apply when another field is invalid.
    a.send({'type':'update_settings','settings':{'theme':'light','max_tool_rounds':'invalid'}})
    assert a.until('settings_error').get('text'), 'missing invalid-patch error'
    assert settings.read_bytes() == durable_before
    probe=WS(addr); sockets.append(probe); assert probe.snapshot('current_settings')==initial; probe.close(); sockets.remove(probe)
    # Force the real final rename to fail, without mocking saveSettings.
    backup=settings.with_name('review-original-settings.json'); settings.rename(backup); settings.mkdir()
    try:
        a.send({'type':'update_settings','settings':{'theme':'system','max_tool_rounds':7}})
        assert a.until('settings_error').get('text'), 'missing real-save-failure error'
        assert backup.read_bytes()==durable_before
        probe=WS(addr); sockets.append(probe); assert probe.snapshot('current_settings')==initial; probe.close(); sockets.remove(probe)
    finally:
        settings.rmdir(); backup.rename(settings)
    process.stdin.write('FIRST_BUDGET\\n'); process.stdin.flush()"""
text=text.replace(anchor,extra)
text=text.replace("    (root / 'requests.json').write_text", "    (root / 'review-privacy-summary.json').write_text(json.dumps({'distinct_system_hashes':sorted(system_hashes),'user_text_allowlist_asserted_per_request':True,'requests':len(requests)},indent=2)+'\\n')\n    (root / 'requests.json').write_text")
helper=E/'real-process-review.py';helper.write_text(text)
(E/'real-process-helper.diff').write_text(''.join(difflib.unified_diff((S/'gui/test/settings_process.py').read_text().splitlines(True),text.splitlines(True),fromfile='submitted/gui/test/settings_process.py',tofile='owned/reviewer-process.py')))
GO='/usr/local/go/bin/go';PY=json.loads((E/'runtime-tools.json').read_text())['python'];NODE=json.loads((E/'runtime-tools.json').read_text())['node']
actor=C/'internal/agent/actor.go';old=actor.read_bytes();anchor=b'if t.batches >= t.batchLimit {';assert old.count(anchor)==1
negative=old.replace(anchor,b'reviewLiveLimit := a.ToolBatchLimit()\n\tif reviewLiveLimit == 0 { reviewLiveLimit = common.DefaultToolBatchLimit }\n\tif t.batches >= reviewLiveLimit {')
(E/'budget-mutation.diff').write_text(''.join(difflib.unified_diff(old.decode().splitlines(True),negative.decode().splitlines(True),fromfile='candidate/internal/agent/actor.go',tofile='owned-active-turn-reread-mutant/actor.go')))
records=[]
def gui(stage,expected):
 b=T/'bin'/('gui-'+stage);run('budget-'+stage+'-compile',[GO,'build','-o',str(b),'./cmd'],C/'gui',120)
 r={'stage':stage,'actor_sha256':sha(actor),'binary_path':str(b),'binary_sha256':sha(b),'compile_exit':0,'expected_behavior_exit':expected};records.append(r)
 run('budget-'+stage+'-behavior',[PY,str(helper),str(b),str(T/('real-'+stage))],timeout=100,expected=expected)
 if expected:assert 'live update rewrote active turn' in (E/('budget-'+stage+'-behavior.stderr.log')).read_text()
 r['behavior_exit']=expected
 (E/'budget-mutation-results.json').write_text(json.dumps(records,indent=2)+'\n')
gui('old',0)
try:
 actor.write_bytes(negative);gui('negative',1)
finally:actor.write_bytes(old)
assert actor.read_bytes()==old;gui('restored',0)
assert records[0]['binary_sha256']==records[2]['binary_sha256']
# Independent actual SpeechQueue assertions, with deterministic synthetic native
# boundaries. No browser/native/audio claims and no fixed-preset substitute.
speech=E/'speech-review.mjs';speech.write_text("""import assert from 'node:assert/strict';
import {pathToFileURL} from 'node:url';
const {SpeechQueue,nativeSpeechRate}=await import(pathToFileURL(process.argv[2]).href);
let rate=1.2345;const spoken=[];let cancels=0;class U{constructor(text){this.text=text}}
const q=new SpeechQueue(()=>{}, {speak:u=>spoken.push(u),cancel:()=>cancels++},U,()=>rate);
q.enqueue('fraction one');q.enqueue('fraction two');assert.equal(spoken[0].rate,1.2345);
rate=1.8765;assert.equal(spoken[0].rate,1.2345,'active utterance must not be re-rated');spoken[0].onend({type:'end'});assert.equal(spoken[1].rate,1.8765);
rate=4.25;q.enqueue('extended approximate');spoken[1].onend({type:'end'});assert.equal(spoken[2].rate,2.75);
const stale=spoken[2];q.cancel();assert(!q.busy);assert.equal(cancels,1);rate=.7654;q.enqueue('fraction below one');stale.onend({type:'end'});assert.equal(spoken[3].rate,.7654);assert(q.busy);spoken[3].onerror({type:'error',error:'controlled'});assert(!q.busy);
assert.equal(nativeSpeechRate(2),2);assert(nativeSpeechRate(3.125)>nativeSpeechRate(3));
console.log(JSON.stringify({passed:true,rates:spoken.map(u=>u.rate),scope:'Actual production SpeechQueue with synthetic synthesis; NOT browser/native speech proof'}));
""")
sp=C/'gui/web/gui/speech.js';spold=sp.read_bytes();anchor=b'utterance.rate = nativeSpeechRate(this.speed());';assert spold.count(anchor)==1;spbad=spold.replace(anchor,b'utterance.rate = Math.round(nativeSpeechRate(this.speed()));');sr=[]
(E/'speech-mutation.diff').write_text(''.join(difflib.unified_diff(spold.decode().splitlines(True),spbad.decode().splitlines(True),fromfile='candidate/speech.js',tofile='owned-rounding-mutant/speech.js')))
for stage,data,expected in [('old',spold,0),('negative',spbad,1),('restored',spold,0)]:
 try:
  sp.write_bytes(data);run('speech-'+stage+'-parse',[NODE,'--check',str(sp)],timeout=30)
  run('speech-'+stage+'-behavior',[NODE,str(speech),str(sp)],timeout=30,expected=expected)
  sr.append({'stage':stage,'source_sha256':sha(sp),'parse_exit':0,'behavior_exit':expected})
 finally:sp.write_bytes(spold)
 (E/'speech-mutation-results.json').write_text(json.dumps(sr,indent=2)+'\n')
# Retain only small text summaries from real process controls in the portable
# bundle; raw histories/process logs/tools and binaries remain local-only.
for stage in ['old','negative','restored']:
 root=T/('real-'+stage)
 for n in ['commands.json','requests.json','review-privacy-summary.json']:
  (E/(stage+'-'+n)).write_bytes((root/n).read_bytes())
assert all(sha(C/x['path'])==x['sha256'] for x in copy)
(E/'sensitivity-complete.json').write_text(json.dumps({'copied_source_files':len(copy),'copy_restored':True,'mutation_count':2,'production_original_edited':False,'browser_native_paid':False,'all_owned_groups_ended':True},indent=2)+'\n')
