#!/usr/bin/env python3
"""Sequential command recorder; no timeout and no ambient credential capture."""
import datetime,hashlib,json,os,pathlib,subprocess,sys,time
E=pathlib.Path(__file__).resolve().parent
label,cwd,*argv=sys.argv[1:]
env={k:os.environ[k] for k in ['PATH','HOME','USER','TMPDIR','GOCACHE','GOMODCACHE','GOPATH'] if k in os.environ}
env.update({'GOMAXPROCS':'2','GOFLAGS':'-p=1 -mod=readonly','GOPROXY':'off','GOSUMDB':'off','GOTOOLCHAIN':'local','TMPDIR':(E/'scratch-path.txt').read_text().strip()})
def disk():
 s=os.statvfs(E);return s.f_bavail*s.f_frsize
m={'argv':argv,'cwd':cwd,'started_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'env_policy':'allowlist build/runtime environment; no credentials inherited','env_overrides':{k:env[k] for k in ['GOMAXPROCS','GOFLAGS','GOPROXY','GOSUMDB','GOTOOLCHAIN','TMPDIR']},'free_bytes_before':disk()}
t=time.monotonic()
with open(E/(label+'.stdout'),'wb') as out,open(E/(label+'.stderr'),'wb') as err:
 p=subprocess.Popen(argv,cwd=cwd,env=env,stdout=out,stderr=err)
 m['exit']=p.wait()
m['duration_seconds']=time.monotonic()-t;m['free_bytes_after']=disk()
for suffix in ['stdout','stderr']:
 m[suffix+'_sha256']=hashlib.sha256((E/(label+'.'+suffix)).read_bytes()).hexdigest()
(E/(label+'.meta.json')).write_text(json.dumps(m,indent=2)+'\n')
print(json.dumps(m))
print((E/(label+'.stdout')).read_text(errors='replace')[-9000:])
print((E/(label+'.stderr')).read_text(errors='replace')[-3000:])
sys.exit(m['exit'])
