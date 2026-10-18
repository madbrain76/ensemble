#!/usr/bin/env python3
"""Sequential command receipt runner. No deadline and no provider environment."""
import datetime,json,os,pathlib,subprocess,sys,time
root=pathlib.Path(__file__).resolve().parent
name,cwd,*argv=sys.argv[1:]
env={k:os.environ[k] for k in ('PATH','HOME','TMPDIR','LANG','LC_ALL','TERM') if k in os.environ}
env.update(GIT_PAGER='cat',GOPROXY='off',GOTOOLCHAIN='local',PYTHONDONTWRITEBYTECODE='1')
start=time.monotonic(); began=datetime.datetime.now(datetime.timezone.utc).isoformat()
with open(root/(name+'.log'),'wb') as out:
 result=subprocess.run(argv,cwd=cwd,env=env,stdout=out,stderr=subprocess.STDOUT)
receipt={'name':name,'cwd':cwd,'argv':argv,'environment_names':sorted(env),'started_utc':began,'elapsed_seconds':round(time.monotonic()-start,3),'exit_code':result.returncode,'output':name+'.log'}
(root/(name+'.receipt.json')).write_text(json.dumps(receipt,indent=2)+'\n')
print(json.dumps(receipt),flush=True);print((root/(name+'.log')).read_text(errors='replace'),flush=True)
sys.exit(result.returncode)
