import subprocess,sys,time,json,os
from pathlib import Path
w=Path(__file__).resolve().parent
name,cwd,*cmd=sys.argv[1:]
env=os.environ.copy();env.update(TMPDIR=str(w/'tmp'),GOPROXY='off',GOTOOLCHAIN='local')
# Do not copy ambient credentials into an offline harness. Fixtures set their own local fake key.
for k in list(env):
 if any(x in k for x in ['API_KEY','AUTH_TOKEN','ACCESS_TOKEN']): env.pop(k,None)
t=time.time()
with (w/'logs'/f'{name}.stdout.txt').open('wb') as out,(w/'logs'/f'{name}.stderr.txt').open('wb') as err:
 p=subprocess.run(cmd,cwd=cwd,env=env,stdout=out,stderr=err)
r=dict(name=name,cwd=cwd,argv=cmd,start_unix=t,end_unix=time.time(),elapsed_seconds=time.time()-t,exit=p.returncode,environment_overrides={'TMPDIR':str(w/'tmp'),'GOPROXY':'off','GOTOOLCHAIN':'local'},ambient_credentials='API_KEY/AUTH_TOKEN/ACCESS_TOKEN names removed; no values recorded')
(w/'logs'/f'{name}.run.json').write_text(json.dumps(r,indent=2)+'\n')
print(json.dumps(r,indent=2));print((w/'logs'/f'{name}.stdout.txt').read_text());print((w/'logs'/f'{name}.stderr.txt').read_text())
