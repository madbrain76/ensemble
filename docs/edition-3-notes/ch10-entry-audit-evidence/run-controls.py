#!/usr/bin/env python3
"""Finite, sequential Chapter10 entry controls. Run from repository root.
Fresh output directories required; original grader/source files remain read-only.
"""
from pathlib import Path
import json,hashlib,os,subprocess,shutil,time,signal,threading,difflib
R=Path.cwd(); E=R/'docs/edition-3-notes/ch10-entry-audit-evidence'
W=Path(json.loads((E/'workspace.json').read_text())['owned_scratch'])
GO='/usr/local/go/bin/go'
env={'PATH':'/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(W/'home'),'TMPDIR':str(W/'tmp'),'GOCACHE':'/Users/bill/Library/Caches/go-build','GOMODCACHE':'/Users/bill/go/pkg/mod','GOPATH':str(W/'gopath'),'GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOFLAGS':'-mod=readonly','GOWORK':'off','LC_ALL':'C','TZ':'UTC','NO_PROXY':'127.0.0.1,localhost'}
(E/'launch-environment.json').write_text(json.dumps(env,indent=2)+'\n')
receipts=[]; bins={}; saves={}; commands=set(); stop=False
(E/'observed-saves').mkdir();
def digest(p):
 b=p.read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def sample():
 while not stop:
  for d in list((W/'tmp').glob('course-grade-*')):
   p=d/'submission'
   if str(p) not in bins and p.is_file():
    try: bins[str(p)]=digest(p)
    except OSError: pass
  for d in list((W/'tmp').glob('edition2-persistence-*')):
   for p in d.rglob('*.json'):
    if p.name not in ('save.json','chosen.json'):continue
    try:
     b=p.read_bytes();h=hashlib.sha256(b).hexdigest()
     if h not in saves:
      name=f'{len(saves):03d}-{p.parent.name}-{p.name}'
      (E/'observed-saves'/name).write_bytes(b)
      saves[h]={'source_path':str(p),'evidence_path':'observed-saves/'+name,'sha256':h,'bytes':len(b),'observed_at':time.time()}
    except OSError:pass
  try:
   ps=subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True)
   commands.update(x.strip() for x in ps.splitlines() if str(W) in x)
  except Exception:pass
  time.sleep(.075)
t=threading.Thread(target=sample,daemon=True);t.start()
def run(name,argv,cwd=R,seconds=180,expected=None):
 start=time.time();out=E/(name+'.stdout.txt');err=E/(name+'.stderr.txt')
 exe=Path(argv[0]); item={'name':name,'argv':list(map(str,argv)),'cwd':str(cwd),'timeout_seconds':seconds,'environment':'launch-environment.json','start_unix':start,'executable_identity':digest(exe) if exe.is_file() else None}
 with out.open('wb') as o,err.open('wb') as e:
  p=subprocess.Popen(list(map(str,argv)),cwd=cwd,env=env,stdout=o,stderr=e,start_new_session=True)
  item['pid']=p.pid
  try: item['exit_code']=p.wait(timeout=seconds);item['timed_out']=False
  except subprocess.TimeoutExpired:
   item['timed_out']=True;os.killpg(p.pid,signal.SIGTERM)
   try:p.wait(timeout=5)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);p.wait()
   item['exit_code']=p.returncode
 item['elapsed_seconds']=round(time.time()-start,3);item['stdout']=out.name;item['stderr']=err.name
 receipts.append(item);(E/'command-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
 print(name,item['exit_code'],'timeout',item['timed_out'],'seconds',item['elapsed_seconds'],flush=True)
 if expected is not None and item['exit_code'] not in expected:raise RuntimeError('unexpected command result '+name)
 return item
try:
 run('00-go-version',[GO,'version'],expected=[0])
 # Build the real adapter once. The separate marker proves the fixture can execute.
 run('01-build-edition3',[GO,'build','-o',str(W/'bin/grade-edition3'),'./cmd/grade-edition3'],expected=[0])
 marker=W/'marker';marker.mkdir()
 (marker/'go.mod').write_text('module marker\n\ngo 1.25\n')
 marker_source='package main\nimport "os"\nfunc main(){if err:=os.WriteFile(os.Args[1],[]byte("DIRECT_POSITIVE_MARKER\\n"),0600);err!=nil{panic(err)}}\n'
 (marker/'main.go').write_text(marker_source);(E/'marker-main.go.txt').write_text(marker_source)
 run('02-build-marker',[GO,'build','-o',str(W/'bin/marker'),'.'],marker,expected=[0])
 mark=W/'marker-positive.txt'
 run('03-marker-positive',[str(W/'bin/marker'),str(mark)],marker,expected=[0])
 before=digest(mark);(E/'marker-positive.txt').write_bytes(mark.read_bytes())
 run('04-unsupported-ch10-no-target',[str(W/'bin/grade-edition3'),'-ch','10','-json'],marker,expected=[2])
 run('05-unsupported-ch10-marker-target',[str(W/'bin/grade-edition3'),'-ch','10','-json',str(marker)],marker,expected=[2])
 after=digest(mark);assert before==after
 (E/'routing-marker-result.json').write_text(json.dumps({'positive_marker_before':before,'positive_marker_after':after,'unchanged':before==after,'no_target_cwd':str(marker),'scope':'unsupported chapter rejected before positional/default target build; positive marker proves executable fixture only'},indent=2)+'\n')
 # Full exported API once, unchanged reference and complete explicit auxiliaries.
 driver=R/'.ch10-entry-audit-local/full';driver.mkdir()
 shutil.copyfile(E/'full-api.go.txt',driver/'main.go')
 run('06-build-full-api',[GO,'build','-o',str(W/'bin/full-api'),'./.ch10-entry-audit-local/full'],expected=[0])
 case=Path('/Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json')
 shutil.copyfile(case,E/'consumer-case.json')
 run('07-full-chapter10',[str(W/'bin/full-api'),str(R/'solutions/ch10'),str(case)],seconds=600,expected=[0])
 # Exact grader/package test copied to owned scratch; only reference history.Save mutates.
 copy=W/'focused';copy.mkdir()
 for p in ('go.mod','go.sum'):shutil.copyfile(R/p,copy/p)
 shutil.copytree(R/'internal',copy/'internal')
 (copy/'solutions').mkdir();shutil.copytree(R/'solutions/ch10',copy/'solutions/ch10')
 run('08-build-focused-test',[GO,'test','-c','-o',str(W/'bin/focused.test'),'./internal/gradeedition2'],copy,expected=[0])
 testargs=[str(W/'bin/focused.test'),'-test.run','^TestPersistenceLifecycleThroughCLIAndGUI$','-test.v','-test.timeout','150s','-test.count','1']
 run('09-focused-baseline',testargs,copy/'internal/gradeedition2',seconds=170,expected=[0])
 victim=copy/'solutions/ch10/internal/history/save.go'; original=victim.read_text()
 old='SystemPrompt: h.parent.Engine().SystemPrompt()';new='SystemPrompt: "AUDIT_BROKEN_EFFECTIVE_PROMPT"'
 assert original.count(old)==1
 changed=original.replace(old,new); victim.write_text(changed)
 (E/'mutation.patch').write_text(''.join(difflib.unified_diff(original.splitlines(True),changed.splitlines(True),fromfile='a/solutions/ch10/internal/history/save.go',tofile='b/solutions/ch10/internal/history/save.go')))
 (E/'mutation-source-identity.json').write_text(json.dumps({'path':str(victim),'baseline_sha256':hashlib.sha256(original.encode()).hexdigest(),'mutant_sha256':hashlib.sha256(changed.encode()).hexdigest(),'property':'effective system_prompt metadata saved must equal actual wire instruction','expected':'save-shape and current-config default-load detect mismatched metadata; compilation must remain valid'},indent=2)+'\n')
 run('10-build-mutant-cli',[GO,'build','-o',str(W/'bin/mutant-cli'),'.'],copy/'solutions/ch10',expected=[0])
 run('11-build-mutant-gui',[GO,'build','-o',str(W/'bin/mutant-gui'),'./cmd'],copy/'solutions/ch10/gui',expected=[0])
 run('12-focused-mutation',testargs,copy/'internal/gradeedition2',seconds=170,expected=[1])
 victim.write_text(original);assert digest(victim)==digest(R/'solutions/ch10/internal/history/save.go')
 run('13-focused-restoration',testargs,copy/'internal/gradeedition2',seconds=170,expected=[0])
 (E/'restoration-result.json').write_text(json.dumps({'restored_file':str(victim),**digest(victim),'matches_reference':True},indent=2)+'\n')
 # Record explicit binaries; no binary is exported in portable evidence.
 binaries={str(p):digest(p) for p in sorted((W/'bin').iterdir()) if p.is_file()}
 (E/'explicit-binaries.json').write_text(json.dumps(binaries,indent=2)+'\n')
finally:
 stop=True;t.join()
 (E/'observed-transient-binaries.json').write_text(json.dumps({'scope':'sampled course-grade submission binaries only; not exhaustive child process identity','binaries':bins},indent=2)+'\n')
 (E/'observed-save-index.json').write_text(json.dumps({'scope':'non-exhaustive snapshots sampled while existing grader owns lifecycle; no claim of every phase','saves':list(saves.values())},indent=2)+'\n')
 (E/'observed-processes.txt').write_text('\n'.join(sorted(commands))+'\n')
 print('runner finished',flush=True)
