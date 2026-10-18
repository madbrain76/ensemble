import hashlib,json,os,signal,subprocess,time,shutil
from pathlib import Path
E=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch08-code-review-evidence');T=Path((E/'scratch-path.txt').read_text().strip());ENV=json.loads((E/'clean-environment.json').read_text())
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def run(name,argv,cwd=None,timeout=240,expected=0,env=None):
 free=shutil.disk_usage(T).free
 if free<512*1024*1024:raise RuntimeError('actual disk pressure: '+str(free))
 before=time.time();out=E/(name+'.stdout.log');err=E/(name+'.stderr.log');timed=False
 with out.open('wb') as o,err.open('wb') as e:
  p=subprocess.Popen([str(x) for x in argv],cwd=str(cwd or T/'work'),env=env or ENV,stdin=subprocess.DEVNULL,stdout=o,stderr=e,start_new_session=True)
  try:code=p.wait(timeout=timeout)
  except subprocess.TimeoutExpired:
   timed=True;os.killpg(p.pid,signal.SIGKILL);code=p.wait()
 try:os.killpg(p.pid,0);remaining=True
 except ProcessLookupError:remaining=False
 if remaining:
  os.killpg(p.pid,signal.SIGKILL);time.sleep(.1)
  try:os.killpg(p.pid,0);remaining=True
  except ProcessLookupError:remaining=False
 result={'name':name,'argv':[str(x) for x in argv],'cwd':str(cwd or T/'work'),'environment_file':'clean-environment.json','pid_pgid':p.pid,'exit_code':code,'expected_exit':expected,'elapsed_seconds':round(time.time()-before,3),'timeout_seconds':timeout,'timed_out':timed,'remaining_process_group':remaining,'free_before':free,'free_after':shutil.disk_usage(T).free,'stdout_sha256':sha(out),'stderr_sha256':sha(err)}
 (E/(name+'.receipt.json')).write_text(json.dumps(result,indent=2)+'\n');print(name,code,result['elapsed_seconds'],'seconds',flush=True)
 assert not timed and not remaining,result
 assert code==expected,result
 return result
