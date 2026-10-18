import subprocess,sys,time,json,os,hashlib,shutil,datetime
from pathlib import Path
W=Path(__file__).resolve().parent
ENV=os.environ.copy();ENV.update(TMPDIR=str(W/'tmp'),GOPROXY='off',GOTOOLCHAIN='local',GIT_PAGER='cat')
for k in list(ENV):
 if any(x in k for x in ['API_KEY','AUTH_TOKEN','ACCESS_TOKEN']): ENV.pop(k,None)
def run(name,cmd,cwd='/Users/bill/projects/ensemble',extra=None):
 env=ENV.copy();env.update(extra or {});t=time.time();free=shutil.disk_usage(W).free
 p=subprocess.run([str(x) for x in cmd],cwd=str(cwd),env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
 out=W/'logs'/name;out.with_suffix('.stdout.txt').write_bytes(p.stdout);out.with_suffix('.stderr.txt').write_bytes(p.stderr)
 report=None
 try:report=json.loads(p.stdout)
 except (ValueError,UnicodeError):pass
 m=dict(name=name,argv=[str(x) for x in cmd],cwd=str(cwd),start_utc=datetime.datetime.fromtimestamp(t,datetime.timezone.utc).isoformat(),seconds=time.time()-t,exit=p.returncode,free_before=free,free_after=shutil.disk_usage(W).free,environment_overrides={'TMPDIR':str(W/'tmp'),'GOPROXY':'off','GOTOOLCHAIN':'local','GIT_PAGER':'cat',**(extra or {})},credentials='ambient API_KEY/AUTH_TOKEN/ACCESS_TOKEN names removed; no values recorded',stdout_sha256=hashlib.sha256(p.stdout).hexdigest(),stderr_sha256=hashlib.sha256(p.stderr).hexdigest(),score=report.get('score') if isinstance(report,dict) else None,max_score=report.get('max_score') if isinstance(report,dict) else None,passed=report.get('passed') if isinstance(report,dict) else None)
 out.with_suffix('.run.json').write_text(json.dumps(m,indent=2)+'\n');print(json.dumps(m),flush=True);return m,report,p.stdout.decode(errors='replace'),p.stderr.decode(errors='replace')
if __name__=='__main__':
 name,cwd,*cmd=sys.argv[1:];m,r,o,e=run(name,cmd,cwd);print(o);print(e);sys.exit(m['exit'])
