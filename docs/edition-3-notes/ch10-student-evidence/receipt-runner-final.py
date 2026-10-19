import subprocess,sys,time,datetime,json,hashlib,pathlib,shutil
name,cwd,*cmd=sys.argv[1:]
base=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch10-student-evidence')
if (base/(name+'.json')).exists():raise SystemExit('refuse receipt overwrite')
if shutil.disk_usage(cwd).free < 1024**3: raise SystemExit('disk floor blocker')
env={'PATH':'/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin','HOME':'/Users/bill','TMPDIR':'/tmp','GOMAXPROCS':'2','GOFLAGS':'-mod=readonly -p=1','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','LANG':'en_US.UTF-8'}
source=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
def identity():
 def git(*args):return subprocess.check_output(['git','--no-pager',*args],cwd=source,env=env,text=True).strip()
 return {'head':git('rev-parse','HEAD'),'tree':git('rev-parse','HEAD^{tree}'),'status':git('status','--porcelain','--untracked-files=normal')}
def binaries():
 return {p.name:{'path':str(p),'bytes':p.stat().st_size,'sha256':hashlib.sha256(p.read_bytes()).hexdigest()} for p in [source/'.ch10-local/bin/frozen-cli',source/'.ch10-local/bin/frozen-gui'] if p.exists()}
before=identity();expected=json.loads((source/'.ch10-local/freeze.json').read_text())
if before['status'] or any(before[k]!=v for k,v in expected.items()):raise SystemExit('source freeze mismatch')
r={'source_before':before,'binaries_before':binaries(),'argv':cmd,'cwd':cwd,'env':env,'start_utc':datetime.datetime.now(datetime.timezone.utc).isoformat()};start=time.monotonic()
with (base/(name+'.stdout.txt')).open('wb') as out,(base/(name+'.stderr.txt')).open('wb') as err:
 p=subprocess.run(cmd,cwd=cwd,env=env,stdout=out,stderr=err)
r.update(exit_code=p.returncode,duration_seconds=time.monotonic()-start,end_utc=datetime.datetime.now(datetime.timezone.utc).isoformat(),source_after=identity(),binaries_after=binaries())
r['outputs']=[]
for suffix in ['stdout.txt','stderr.txt']:
 path=base/(name+'.'+suffix);b=path.read_bytes();r['outputs'].append({'path':str(path),'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()})
(base/(name+'.json')).write_text(json.dumps(r,indent=2)+'\n')
print(json.dumps({'name':name,'exit':p.returncode,'duration':r['duration_seconds'],'free_bytes':shutil.disk_usage(cwd).free}))
for suffix in ['stdout.txt','stderr.txt']:
 b=(base/(name+'.'+suffix)).read_bytes()
 if b:print(b.decode(errors='replace')[-7000:])
if r['source_after']!=before:raise SystemExit('source changed during gate')
if not name.startswith('build-') and r['binaries_after']!=r['binaries_before']:raise SystemExit('binary changed during gate')
sys.exit(p.returncode)
