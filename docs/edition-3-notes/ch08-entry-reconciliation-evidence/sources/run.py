from pathlib import Path
import os,sys,subprocess,json,time,signal,shutil
R=Path('/Users/bill/projects/ensemble');E=R/'docs/edition-3-notes/ch08-entry-reconciliation-evidence';T=Path((E/'scratch-path.txt').read_text().strip());env=json.loads((E/'clean-env.json').read_text())
def run(name,argv,cwd,timeout):
 if shutil.disk_usage(R).free<300*1024*1024:raise RuntimeError('resource floor reached')
 en=dict(env,AUDIT_STAGE=name);start=time.time()
 with (E/(name+'.stdout')).open('wb') as out,(E/(name+'.stderr')).open('wb') as err:
  p=subprocess.Popen(argv,cwd=cwd,env=en,stdin=subprocess.DEVNULL,stdout=out,stderr=err,start_new_session=True)
  timed=False
  try:rc=p.wait(timeout=timeout)
  except subprocess.TimeoutExpired:
   timed=True;os.killpg(p.pid,signal.SIGKILL);rc=p.wait()
  try:os.killpg(p.pid,0);alive=True
  except ProcessLookupError:alive=False
  if alive:
   os.killpg(p.pid,signal.SIGKILL);time.sleep(.1)
   try:os.killpg(p.pid,0);alive=True
   except ProcessLookupError:alive=False
 rec={'name':name,'argv':argv,'cwd':str(cwd),'exit':rc,'timed_out':timed,'seconds':round(time.time()-start,3),'pid_pgid':p.pid,'process_group_remaining':alive,'free_bytes_after':shutil.disk_usage(R).free}
 (E/(name+'.receipt.json')).write_text(json.dumps(rec,indent=2)+'\n');print(json.dumps(rec),flush=True)
 return rc
if __name__=='__main__':
 if run('build-driver',['go','build','-o',str(T/'bin/driver'),'.'],T/'driver',120):sys.exit(2)
 if run('build-cli',['go','build','-o',str(T/'bin/ch08-cli'),'.'],R/'solutions/ch08',120):sys.exit(2)
 argv=[str(T/'bin/driver'),'-mode','full','-bin',str(T/'bin/ch08-cli'),'-core',str(R/'solutions/ch08'),'-gui',str(R/'solutions/ch08/gui'),'-consumer',str(R/'solutions/ch08/ch05'),'-workflow',str(R/'solutions/ch08/ch06'),'-stream',str(R/'solutions/ch08/ch07'),'-stream-model','claude-sonnet-4-6','-gui-model','claude-sonnet-4-6','-consumer-case',str(T/'consumer-case.json'),'-media-models',str(T/'media-models.json')]
 sys.exit(run('historical-full',argv,T/'work',420))
