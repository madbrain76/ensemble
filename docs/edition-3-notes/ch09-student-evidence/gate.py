import datetime, json, os, pathlib, shutil, subprocess, sys, time
base=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student/.ch09-local').resolve()
name,cwd,*command=sys.argv[1:]
env={k:os.environ[k] for k in ('PATH','HOME','TMPDIR') if k in os.environ}
env.update(GOPROXY='off',GOSUMDB='off',GOTOOLCHAIN='local')
free=shutil.disk_usage(cwd).free
if free<1024**3: raise SystemExit('resource guardrail: less than 1GiB free')
start=time.monotonic();stamp=datetime.datetime.now(datetime.timezone.utc).isoformat()
with open(base/(name+'.stdout.txt'),'wb') as out,open(base/(name+'.stderr.txt'),'wb') as err:
 p=subprocess.run(command,cwd=cwd,env=env,stdout=out,stderr=err)
receipt={'name':name,'cwd':cwd,'command':command,'environment':env,'started_utc':stamp,'duration_seconds':time.monotonic()-start,'exit_code':p.returncode,'free_bytes_before':free,'free_bytes_after':shutil.disk_usage(cwd).free,'stdout':name+'.stdout.txt','stderr':name+'.stderr.txt'}
(base/(name+'.receipt.json')).write_text(json.dumps(receipt,indent=2)+'\n')
print(json.dumps(receipt,indent=2))
print((base/(name+'.stdout.txt')).read_text()[-6000:]);print((base/(name+'.stderr.txt')).read_text()[-3000:])
raise SystemExit(p.returncode)
