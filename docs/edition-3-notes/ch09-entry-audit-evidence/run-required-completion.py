#!/usr/bin/env python3
import os,pathlib,subprocess,json,hashlib,shutil,time,signal,tarfile,difflib
R=pathlib.Path('/Users/bill/projects/ensemble'); E=R/'docs/edition-3-notes/ch09-entry-audit-evidence'; S=pathlib.Path('/tmp/ch09-entry-audit')
for d in ['home','tmp','bin','work']: (S/d).mkdir(exist_ok=True)
go=shutil.which('go'); env={'PATH':str(pathlib.Path(go).parent)+':/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(S/'home'),'TMPDIR':str(S/'tmp'),'GOCACHE':'/Users/bill/Library/Caches/go-build','GOMODCACHE':'/Users/bill/go/pkg/mod','GOPATH':'/Users/bill/go','GOWORK':'off','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOFLAGS':'-p=1 -mod=readonly','GOMAXPROCS':'2','LANG':'C','TZ':'UTC'}
(E/'environment-completion.json').write_text(json.dumps({'go':go,'env':env,'provider_credentials':'none; unchanged graders inject fake values and local endpoints','network':'Go proxy/sumdb disabled; test fake providers only'},indent=2)+'\n')
def tree(p):
 return {str(f.relative_to(p)):{'sha256':hashlib.sha256(f.read_bytes()).hexdigest(),'bytes':f.stat().st_size} for f in sorted(p.rglob('*')) if f.is_file()}
def run(name,args,cwd,limit):
 start=time.time(); meta={'argv':list(map(str,args)),'cwd':str(cwd),'timeout_seconds':limit,'environment_file':'environment-completion.json','start_unix':start}
 with (E/(name+'.stdout')).open('wb') as out,(E/(name+'.stderr')).open('wb') as err:
  p=subprocess.Popen(list(map(str,args)),cwd=cwd,env=env,stdout=out,stderr=err,start_new_session=True); meta['pid']=p.pid
  try: code=p.wait(timeout=limit)
  except subprocess.TimeoutExpired:
   meta['timeout']=True; os.killpg(p.pid,signal.SIGTERM)
   try: code=p.wait(timeout=3)
   except subprocess.TimeoutExpired: os.killpg(p.pid,signal.SIGKILL); code=p.wait()
  meta.update(exit_code=code,elapsed_seconds=round(time.time()-start,3))
 (E/(name+'.json')).write_text(json.dumps(meta,indent=2)+'\n'); print(name,code,flush=True); return code

import threading
C=S/'positive'; original=tree(R/'solutions/ch09'); assert tree(C)==original
sentinel=S/'sentinel.go'; sentinel.write_text('package main\nimport "os"\nfunc main(){if err:=os.WriteFile("/tmp/ch09-entry-audit/sentinel.marker",[]byte("launched\\n"),0600);err!=nil{panic(err)};os.Exit(73)}\n'.replace('\\n','\n'))
# Write unambiguous literal Go source instead.
sentinel.write_text('package main\nimport "os"\nfunc main(){if err:=os.WriteFile("/tmp/ch09-entry-audit/sentinel.marker",[]byte("launched"),0600);err!=nil{panic(err)};os.Exit(73)}\n'.replace('\\n','\n').replace('\\"','"'))
if run('build-sentinel',[go,'build','-o',S/'bin/sentinel',sentinel],S,120): raise SystemExit('sentinel build failed')
marker=S/'sentinel.marker'; assert not marker.exists()
assert run('sentinel-selfcheck',[S/'bin/sentinel'],S/'work',15)==73
assert marker.read_text()=='launched'; (S/'local-artifacts/sentinel-selfcheck.marker').write_bytes(marker.read_bytes()); marker.unlink()
run('entry-no-target',[S/'bin/grade-edition3','-ch','9','-json'],R,15); absent_no_target=not marker.exists()
run('entry-sentinel',[S/'bin/grade-edition3','-ch','9','-core',C,'-consumer',C/'ch05','-workflow',C/'ch06','-stream-exercise',C/'ch07','-stream-model','claude-sonnet-5','-gui',C/'gui','-gui-model','claude-sonnet-5','-json',S/'bin/sentinel'],R,15)
(E/'sentinel-result.json').write_text(json.dumps({'selfcheck_exit':73,'selfcheck_marker':'launched','marker_removed_before_entry':True,'no_target_marker_absent':absent_no_target,'sentinel_marker_absent':not marker.exists(),'sentinel_sha256':hashlib.sha256((S/'bin/sentinel').read_bytes()).hexdigest()},indent=2)+'\n'); assert not marker.exists()
if run('build-driver-corrected',[go,'build','-o',S/'bin/driver-corrected','./docs/edition-3-notes/ch09-entry-audit-evidence/driver-corrected'],R,120): raise SystemExit('corrected driver build failed')
stop=threading.Event(); captured={}; errors=[]
def observe():
 while not stop.is_set():
  for f in (S/'tmp').glob('**/submission'):
   try:
    if not f.is_file() or not os.access(f,os.X_OK) or str(f) in captured: continue
    b=f.read_bytes(); h=hashlib.sha256(b).hexdigest(); dest=S/'local-artifacts'/('generated-'+h); dest.write_bytes(b); dest.chmod(0o755)
    captured[str(f)]={'original_path':str(f),'retained_path':str(dest),'sha256':h,'bytes':len(b)}
   except FileNotFoundError: pass
   except Exception as e: errors.append(str(e))
  stop.wait(.01)
t=threading.Thread(target=observe); t.start()
try: code=run('historical-positive',[S/'bin/driver-corrected',C],S/'work',300)
finally: stop.set(); t.join()
for item in captured.values():
 p=subprocess.run([go,'version','-m',item['retained_path']],env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True); item['go_version_m']=p.stdout; item['inspection_exit']=p.returncode
(E/'generated-artifacts.json').write_text(json.dumps({'artifacts':list(captured.values()),'observer_errors':errors},indent=2)+'\n')
assert tree(C)==original
if code: raise SystemExit('real historical positive failed; no repairs or retries')
# Restore precisely the sole mutated file, then run the exact failed focused command.
M=S/'mutant'; f=M/'internal/skills/skills.go'; f.write_bytes((C/'internal/skills/skills.go').read_bytes()); assert tree(M)==original
run('restored',[go,'test','-work','-v','-count=1','-run','^TestSkillsThroughLoop$','./ensemble'],M,90)
(E/'restored-source.json').write_text(json.dumps({'path':str(M),'exact_tree_matches_baseline':tree(M)==original,'restored_file_sha256':hashlib.sha256(f.read_bytes()).hexdigest()},indent=2)+'\n')
