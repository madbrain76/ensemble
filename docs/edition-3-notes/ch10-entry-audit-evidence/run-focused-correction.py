#!/usr/bin/env python3
"""Correct the empty test selector; do not repeat full composition or routing."""
from pathlib import Path
import os,json,hashlib,subprocess,time,signal
R=Path.cwd();E=R/'docs/edition-3-notes/ch10-entry-audit-evidence';W=Path(json.loads((E/'workspace.json').read_text())['owned_scratch']);env=json.loads((E/'launch-environment.json').read_text());receipts=json.loads((E/'command-receipts.json').read_text());copy=W/'focused'
def digest(p):
 b=p.read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def run(name,argv,cwd,expected):
 out=E/(name+'.stdout.txt');err=E/(name+'.stderr.txt');start=time.time();item={'name':name,'argv':list(map(str,argv)),'cwd':str(cwd),'timeout_seconds':170,'environment':'launch-environment.json','start_unix':start,'executable_identity':digest(Path(argv[0]))}
 with out.open('wb') as o,err.open('wb') as e:
  p=subprocess.Popen(argv,cwd=cwd,env=env,stdout=o,stderr=e,start_new_session=True);item['pid']=p.pid
  try:item['exit_code']=p.wait(timeout=170);item['timed_out']=False
  except subprocess.TimeoutExpired:
   item['timed_out']=True;os.killpg(p.pid,signal.SIGTERM)
   try:p.wait(timeout=5)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);p.wait()
   item['exit_code']=p.returncode
 item.update(elapsed_seconds=round(time.time()-start,3),stdout=out.name,stderr=err.name);receipts.append(item);(E/'command-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n');print(name,item['exit_code'],item['elapsed_seconds'],flush=True)
 assert item['exit_code']==expected,item
 assert '=== RUN   TestChapter10Persistence' in out.read_text(), 'empty test selection'
 return item
v=copy/'solutions/ch10/internal/history/save.go';original=(R/'solutions/ch10/internal/history/save.go').read_text();changed=original.replace('SystemPrompt: h.parent.Engine().SystemPrompt()','SystemPrompt: "AUDIT_BROKEN_EFFECTIVE_PROMPT"');assert original!=changed
args=[str(W/'bin/focused.test'),'-test.run','^TestChapter10Persistence$','-test.v','-test.timeout','150s','-test.count','1']
(E/'selector-correction.json').write_text(json.dumps({'invalid_attempts':['09-focused-baseline','12-focused-mutation'],'invalid_selector':'^TestPersistenceLifecycleThroughCLIAndGUI$','actual_selector':'^TestChapter10Persistence$','reason':'warning no tests to run; these passes carry zero persistence evidence','full_composition_repeated':False},indent=2)+'\n')
v.write_text(original)
try:
 run('14-focused-baseline-correct-selector',args,copy/'internal/gradeedition2',0)
 v.write_text(changed)
 # Same mutant source already built successfully as both real clients in 10/11.
 run('15-focused-mutation-correct-selector',args,copy/'internal/gradeedition2',1)
finally:
 v.write_text(original)
 (E/'restoration-result.json').write_text(json.dumps({'restored_file':str(v),**digest(v),'matches_reference':digest(v)==digest(R/'solutions/ch10/internal/history/save.go')},indent=2)+'\n')
run('16-focused-restoration-correct-selector',args,copy/'internal/gradeedition2',0)
(E/'explicit-binaries.json').write_text(json.dumps({str(p):digest(p) for p in sorted((W/'bin').iterdir()) if p.is_file()},indent=2)+'\n')
checks={}
for rel in ('ch10.go','ch10_test.go'):
 checks[rel]={'original':digest(R/'internal/gradeedition2'/rel),'copy':digest(copy/'internal/gradeedition2'/rel)};assert checks[rel]['original']==checks[rel]['copy']
(E/'focused-grader-copy-identity.json').write_text(json.dumps(checks,indent=2)+'\n')
