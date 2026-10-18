from pathlib import Path
import json,hashlib,subprocess,shutil,os
R=Path('/Users/bill/projects/ensemble');E=R/'docs/edition-3-notes/ch08-grader-review-evidence';T=Path((E/'scratch-path.txt').read_text().strip())
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
C=json.loads((R/'docs/edition-3-notes/ch08-grader-evidence/candidate.json').read_text());M=json.loads((R/'docs/edition-3-notes/ch08-entry-manifest.json').read_text());V={}
for p,x in C['sources'].items():assert sha(R/p)==x['sha256'],p
assert sha(C['binary']['path'])==C['binary']['sha256'];assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=R).decode().strip()==M['main_head'];S=Path('/Users/bill/projects/ensemble-edition-3-student');assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=S).decode().strip()==M['student_head'];assert not subprocess.check_output(['git','status','--porcelain'],cwd=S)
for key in ['inputs','student_files','protected_tracked_changes']:
 diffs=[];carried=[];overridden=[]
 for p,x in M[key].items():
  path=Path(p) if Path(p).is_absolute() else (S/p if key=='student_files' else R/p)
  expected=C['sources'][p]['sha256'] if p in C['sources'] else x['sha256'];actual=sha(path)
  if actual!=expected:diffs.append(str(path))
  (overridden if p in C['sources'] else carried).append(p)
 assert not diffs,(key,diffs);V[key]={'verified_count':len(M[key]),'candidate_overrides':overridden,'unchanged_carried':carried,'differences':diffs}
prior=R/'docs/edition-3-notes/ch08-entry-reconciliation-evidence';P=Path((prior/'scratch-path.txt').read_text().strip());oldart=json.loads((prior/'local-artifacts.json').read_text());indexed={x['path']:x for x in oldart['files']};reuse=[]
for name in ['bin/ch08-cli','consumer-case.json','media-models.json']:
 p=P/name;x=indexed[str(p)];assert sha(p)==x['sha256'];reuse.append(x)
V.update({'main_head':M['main_head'],'student_head':M['student_head'],'student_clean':True,'candidate':C,'reused_artifacts':reuse,'free_bytes':shutil.disk_usage(R).free})
(E/'pins-before.json').write_text(json.dumps(V,indent=2)+'\n');(T/'status-before.z').write_bytes(subprocess.check_output(['git','status','--porcelain=v1','-z','--untracked-files=all'],cwd=R))
base={};oldbase=json.loads((prior/'source-baseline.json').read_text())
for p in oldbase:
 path=Path(p) if Path(p).is_absolute() else R/p;base[str(path)]={'sha256':sha(path),'bytes':path.stat().st_size}
for p in R.joinpath('docs/edition-3-notes').rglob('*'):
 if p.is_file() and E not in p.parents:base[str(p)]={'sha256':sha(p),'bytes':p.stat().st_size}
for p in M['inputs']:
 path=Path(p) if Path(p).is_absolute() else R/p;base[str(path)]={'sha256':sha(path),'bytes':path.stat().st_size}
(E/'protected-baseline.json').write_text(json.dumps(base,indent=2)+'\n')
for d in ['home','tmp','work','shim','bin','driver']:(T/d).mkdir()
for name in ['consumer-case.json','media-models.json']:shutil.copyfile(P/name,T/name);shutil.copyfile(P/name,E/name)
for p in C['sources']:shutil.copyfile(R/p,E/'sources'/((p.replace('/','_'))+'.txt'))
shutil.copyfile('/tmp/cadence-ch08-adapter-before/internal/gradeedition2/ch08.go',T/'old-ch08.go');shutil.copyfile(T/'old-ch08.go',E/'sources/old-ch08.go.txt')
assert sha(T/'old-ch08.go')==M['inputs']['internal/gradeedition2/ch08.go']['sha256']
(E/'source-provenance.json').write_text(json.dumps({'old_settings_source':{'path':'/tmp/cadence-ch08-adapter-before/internal/gradeedition2/ch08.go','sha256':sha(T/'old-ch08.go'),'matches_entry_manifest':True},'reconciliation_reused_artifacts':reuse},indent=2)+'\n')
shim='''#!/usr/bin/python3
import os,sys,subprocess,json,hashlib,time
from pathlib import Path
a=['/usr/local/go/bin/go']+sys.argv[1:];start=time.time();rc=subprocess.call(a)
if len(a)>2 and a[1]=='build':
 r={'stage':os.environ.get('AUDIT_STAGE'),'argv':a,'cwd':os.getcwd(),'exit':rc,'seconds':round(time.time()-start,3)}
 if '-o' in a:
  p=Path(a[a.index('-o')+1]);r['output']=str(p)
  if rc==0 and p.is_file():r.update(bytes=p.stat().st_size,sha256=hashlib.sha256(p.read_bytes()).hexdigest())
 with open(os.environ['AUDIT_BUILD_LOG'],'a') as f:f.write(json.dumps(r)+'\\n')
sys.exit(rc)
'''
(T/'shim/go').write_text(shim);(T/'shim/go').chmod(0o755);(E/'sources/go-shim.py').write_text(shim)
env={'PATH':str(T/'shim')+':/usr/local/go/bin:/usr/bin:/bin','HOME':str(T/'home'),'TMPDIR':str(T/'tmp'),'GOCACHE':'/Users/bill/Library/Caches/go-build','GOMODCACHE':'/Users/bill/go/pkg/mod','GOPATH':'/Users/bill/go','GOPROXY':'off','GOSUMDB':'off','GOTOOLCHAIN':'local','GOFLAGS':'-p=1 -mod=readonly','GOMAXPROCS':'2','AUDIT_BUILD_LOG':str(E/'builds.jsonl'),'PYTHONDONTWRITEBYTECODE':'1'}
(E/'clean-env.json').write_text(json.dumps(env,indent=2)+'\n')
common='''from pathlib import Path
import json,os,subprocess,time,signal,shutil,hashlib
R=Path('/Users/bill/projects/ensemble');E=R/'docs/edition-3-notes/ch08-grader-review-evidence';T=Path((E/'scratch-path.txt').read_text().strip());ENV=json.loads((E/'clean-env.json').read_text())
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def run(name,argv,cwd,timeout=180,override=None):
 assert shutil.disk_usage(R).free>300*1024*1024,'resource floor'
 env=dict(ENV,AUDIT_STAGE=name);env.update(override or {});start=time.time()
 with (E/(name+'.stdout')).open('wb') as out,(E/(name+'.stderr')).open('wb') as err:
  p=subprocess.Popen(argv,cwd=cwd,env=env,stdin=subprocess.DEVNULL,stdout=out,stderr=err,start_new_session=True);timed=False
  try:rc=p.wait(timeout=timeout)
  except subprocess.TimeoutExpired:timed=True;os.killpg(p.pid,signal.SIGKILL);rc=p.wait()
 alive=False
 try:os.killpg(p.pid,0);alive=True
 except ProcessLookupError:pass
 if alive:
  os.killpg(p.pid,signal.SIGKILL);time.sleep(.1)
  try:os.killpg(p.pid,0)
  except ProcessLookupError:alive=False
 receipt={'name':name,'argv':argv,'cwd':str(cwd),'env':env,'exit':rc,'seconds':round(time.time()-start,3),'timed_out':timed,'pid_pgid':p.pid,'process_group_remaining':alive,'stdout_sha256':sha(E/(name+'.stdout')),'stderr_sha256':sha(E/(name+'.stderr')),'free_bytes_after':shutil.disk_usage(R).free}
 (E/(name+'.receipt.json')).write_text(json.dumps(receipt,indent=2)+'\\n');print(json.dumps({k:receipt[k] for k in ['name','exit','seconds','timed_out','process_group_remaining','free_bytes_after']}),flush=True)
 assert not timed and not alive
 return rc
'''
(E/'sources/common.py').write_text(common)
route='''import sys;sys.dont_write_bytecode=True
from common import *
pins=json.loads((E/'pins-before.json').read_text());cli=pins['reused_artifacts'][0]['path'];binary=pins['candidate']['binary']['path']
args=['-json','-core',str(R/'solutions/ch08'),'-gui',str(R/'solutions/ch08/gui'),'-consumer',str(R/'solutions/ch08/ch05'),'-workflow',str(R/'solutions/ch08/ch06'),'-stream-exercise',str(R/'solutions/ch08/ch07'),'-stream-model','claude-sonnet-4-6','-gui-model','claude-sonnet-4-6','-consumer-case',str(T/'consumer-case.json'),'-media-models',str(T/'media-models.json'),cli]
for ch in [8,7]:
 name='historical-ch'+str(ch);assert run(name,[binary,'-ch',str(ch)]+args,T/'work',420)==0
 report=json.loads((E/(name+'.stdout')).read_text());assert report['passed'] and report['score']==100 and report['max_score']==100,report
'''
(E/'sources/routes.py').write_text(route)
print(json.dumps({'scratch':str(T),'candidate_match':True,'baseline_files':len(base),'free_bytes':V['free_bytes'],'reused_cli':reuse[0]},indent=2))
