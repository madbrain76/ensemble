#!/usr/bin/env python3
import hashlib,json,os,pathlib,shutil,subprocess,time
root=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');out=pathlib.Path('/tmp/cadence-ch05-wakeup-evidence/candidate-r2');frozen=out/'frozen';logs=out/'gates';logs.mkdir(exist_ok=False)
identity=json.loads((frozen/'source-identity.json').read_text());env={**os.environ,'GOWORK':'off','GIT_PAGER':'cat'};rows=[]
steps=[('root-vet',root,['go','vet','./...']),('root-race',root,['go','test','-race','-count=1','-timeout','5m','-coverpkg=./...','-coverprofile='+str(logs/'coverage.out'),'./...']),('root-build',root,['go','build','-race','./...']),('cli-build',root,['go','build','-race','-o',str(frozen/'ensemble'),'./cmd/ch01'])]
for label,location,binary in [('gui','gui','gui'),('consumer','examples/ch04-consumer','consumer'),('workflow','examples/ch06','workflow')]:
 cwd=root/location;steps.extend([(label+'-vet',cwd,['go','vet','./...']),(label+'-race',cwd,['go','test','-race','-count=1','-timeout','3m','./...']),(label+'-build',cwd,['go','build','-race','-o',str(frozen/binary),'.'])])
for name,cwd,args in steps:
 free=shutil.disk_usage('/tmp').free
 if free<1024**3:raise SystemExit('disk below1GiB; no cleanup authorized')
 assert not subprocess.check_output(['git','status','--porcelain'],cwd=root)
 start=time.monotonic()
 with (logs/(name+'.txt')).open('w') as log:p=subprocess.run(args,cwd=cwd,env=env,stdout=log,stderr=subprocess.STDOUT)
 row={'name':name,'cwd':str(cwd),'command':args,'environment_overrides':{'GOWORK':'off'},'source_commit':identity['source_commit'],'exit':p.returncode,'elapsed_seconds':round(time.monotonic()-start,3),'free_bytes_before':free};rows.append(row);(logs/'commands.json').write_text(json.dumps(rows,indent=2)+'\n')
 print(name,'exit',p.returncode,flush=True)
 if p.returncode:raise SystemExit(p.returncode)
with (logs/'coverage-summary.txt').open('w') as log:subprocess.run(['go','tool','cover','-func',str(logs/'coverage.out')],cwd=root,env=env,stdout=log,stderr=subprocess.STDOUT,check=True)
for label in ['cli','workflow']:
 path=frozen/('ensemble' if label=='cli' else label);info=subprocess.check_output(['go','version','-m',str(path)],cwd=root,text=True)
 assert 'vcs.revision='+identity['source_commit'] in info and 'vcs.modified=false' in info
 (frozen/(label+'-build-info.txt')).write_text(info);identity['binaries'][label]={'path':str(path),'sha256':hashlib.sha256(path.read_bytes()).hexdigest()}
for name,digest in identity['source_files'].items():assert hashlib.sha256((root/name).read_bytes()).hexdigest()==digest
assert not subprocess.check_output(['git','status','--porcelain'],cwd=root)
(frozen/'identity.json').write_text(json.dumps(identity,indent=2)+'\n');print('ALL GATES PASS; frozen identity written')
