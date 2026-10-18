import subprocess,os,time,json,hashlib,datetime,sys
from pathlib import Path
w=Path('/tmp/cadence-ch07-r4-pause-independent');label=sys.argv[1];cmd=json.loads(sys.argv[2]);cwd=sys.argv[3]
env={'PATH':str(w/'bin')+':/usr/local/go/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin','HOME':'/Users/bill','LANG':'C.UTF-8','TMPDIR':str(w/'tmp'),'GOTOOLCHAIN':'local','GOPATH':'/Users/bill/go','GOPROXY':'off','GOSUMDB':'off'}
t=time.monotonic();start=datetime.datetime.now(datetime.timezone.utc).isoformat();p=subprocess.run(cmd,cwd=cwd,env=env,capture_output=True)
(w/'logs'/f'{label}.stdout.txt').write_bytes(p.stdout);(w/'logs'/f'{label}.stderr.txt').write_bytes(p.stderr)
meta={'argv':cmd,'cwd':cwd,'env':env,'started':start,'seconds':time.monotonic()-t,'exit':p.returncode,'stdout_sha256':hashlib.sha256(p.stdout).hexdigest(),'stderr_sha256':hashlib.sha256(p.stderr).hexdigest(),'score':None,'passed':None}
try:r=json.loads(p.stdout);meta.update({k:r.get(k) for k in ['score','max_score','passed']})
except Exception:pass
(w/'logs'/f'{label}.run.json').write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta,indent=2));print(p.stdout.decode());print(p.stderr.decode())
