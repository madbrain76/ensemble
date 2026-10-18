from pathlib import Path
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
 (E/(name+'.receipt.json')).write_text(json.dumps(receipt,indent=2)+'\n');print(json.dumps({k:receipt[k] for k in ['name','exit','seconds','timed_out','process_group_remaining','free_bytes_after']}),flush=True)
 assert not timed and not alive
 return rc
