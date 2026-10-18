import sys,json,hashlib,subprocess,os,signal,time,shutil,difflib
from pathlib import Path
sys.dont_write_bytecode=True
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');E=R/'docs/edition-3-notes/ch08-code-review-r1-evidence';T=Path((E/'scratch-path.txt').read_text().strip());cfg=json.loads((E/'environment.json').read_text());ENV=cfg['env'];node=cfg['node_path'];sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
wire=Path(json.loads((R/'docs/edition-3-notes/ch08-student-r1-evidence/wire-provenance.json').read_text())['raw_wire_local_only']['path'])
C=T/'copy';(C/'gui/test').mkdir(parents=True);shutil.copytree(S/'gui/web',C/'gui/web');shutil.copy2(S/'gui/test/correction-r1-page.mjs',C/'gui/test/correction-r1-page.mjs')
if (S/'gui/package.json').exists():shutil.copy2(S/'gui/package.json',C/'gui/package.json')
copy_before={str(p.relative_to(C)):sha(p) for p in sorted(C.rglob('*')) if p.is_file()};(E/'owned-copy-before.json').write_text(json.dumps(copy_before,indent=2)+'\n')
receipts=[]
def run(name,args,expected=0):
 free=shutil.disk_usage(T).free;assert free>512*1024*1024,('resource floor',free)
 wd=T/'work'/name;wd.mkdir();home=wd/'home';home.mkdir();env={**ENV,'HOME':str(home)};out=E/(name+'.stdout.log');err=E/(name+'.stderr.log');start=time.monotonic();timeout=False
 with out.open('wb') as o,err.open('wb') as e:
  p=subprocess.Popen(args,cwd=wd,env=env,stdin=subprocess.DEVNULL,stdout=o,stderr=e,start_new_session=True)
  try:rc=p.wait(timeout=30)
  except subprocess.TimeoutExpired:
   timeout=True;os.killpg(p.pid,signal.SIGTERM)
   try:rc=p.wait(timeout=3)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);rc=p.wait()
 try:os.killpg(p.pid,0);alive=True
 except ProcessLookupError:alive=False
 if alive:
  os.killpg(p.pid,signal.SIGTERM)
  raise AssertionError(('unexpected remaining owned process group',p.pid))
 r={'name':name,'argv':list(map(str,args)),'cwd':str(wd),'home':str(home),'environment_ref':'environment.json','pid_pgid':p.pid,'exit_code':rc,'expected_exit':expected,'elapsed_seconds':round(time.monotonic()-start,3),'deadline_seconds':30,'timed_out':timeout,'remaining_process_group':alive,'stdout_sha256':sha(out),'stderr_sha256':sha(err),'free_bytes_after':shutil.disk_usage(T).free};(E/(name+'.receipt.json')).write_text(json.dumps(r,indent=2)+'\n');receipts.append(r);print(name,'exit',rc,'seconds',r['elapsed_seconds'],flush=True);assert rc==expected and not timeout,r;return r
fixture=S/'gui/test/correction-r1-page.mjs'
plan=[('current-page',[node,'--experimental-vm-modules',str(fixture)]),('current-page-real-wire',[node,'--experimental-vm-modules',str(fixture),str(wire)]),('current-workbench',[node,str(S/'gui/test/workbench.mjs')]),('current-queues',[node,str(S/'gui/test/queues.mjs')])]
(E/'fresh-plan.json').write_text(json.dumps({'controls':[{'name':n,'argv':a} for n,a in plan],'mutation':'Remove stale-utterance done guard in owned speech.js; parse/behavior old-negative-restored. Expected negative exit 1; no other controls.','original_source_write_allowed':False},indent=2)+'\n')
for name,args in plan:run(name,args)
p=C/'gui/web/gui/speech.js';original=p.read_bytes();text=original.decode();anchor='      if (this.active !== utterance || generation !== this.generation) return;';assert text.count(anchor)==1;mutant=text.replace(anchor,'      // REVIEW MUTANT: allow a canceled utterance to complete the new generation.').encode();(E/'speech-guard-mutation.diff').write_text(''.join(difflib.unified_diff(text.splitlines(True),mutant.decode().splitlines(True),fromfile='owned/speech.js (old)',tofile='owned/speech.js (negative)')))
F=C/'gui/test/correction-r1-page.mjs';results=[]
try:
 for label,source,expected in [('old',original,0),('negative',mutant,1),('restored',original,0)]:
  p.write_bytes(source);parse=run('guard-'+label+'-parse',[node,'--check',str(p)]);behavior=run('guard-'+label+'-behavior',[node,'--experimental-vm-modules',str(F)],expected)
  if label=='negative':assert 'AssertionError' in (E/'guard-negative-behavior.stderr.log').read_text()
  results.append({'stage':label,'speech_sha256':sha(p),'fixture_sha256':sha(F),'parse_exit':parse['exit_code'],'behavior_exit':behavior['exit_code']})
finally:p.write_bytes(original)
assert copy_before=={str(p.relative_to(C)):sha(p) for p in sorted(C.rglob('*')) if p.is_file()};assert results[0]['speech_sha256']==results[2]['speech_sha256']!=results[1]['speech_sha256']
(E/'mutation-results.json').write_text(json.dumps({'production_checkout_unchanged':True,'copy_restored':True,'stages':results,'scope':'Real production SpeechQueue and submitted fixture in inert Node page; no browser/native speech or synthesis proof.'},indent=2)+'\n');(E/'controls-summary.json').write_text(json.dumps({'commands':receipts,'command_count':len(receipts),'unexpected_failures':[],'expected_negative':'guard-negative-behavior','owned_copy_restored':True,'go_or_grader_runs':0,'all_owned_process_groups_ended':True},indent=2)+'\n');print('FINITE CONTROLS COMPLETE',len(receipts),flush=True)
