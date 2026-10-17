#!/usr/bin/env python3
import hashlib,json,os,pathlib,shutil,subprocess,time
p=pathlib.Path('/tmp/cadence-ch05-wakeup-evidence/candidate-r2');f=p/'frozen';root=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');main=pathlib.Path('/Users/bill/projects/ensemble')
identity=json.loads((f/'identity.json').read_text());mapping=p/'media-models.json';mapping.write_bytes(pathlib.Path('/tmp/cadence-ch05-evidence-r1/media-models.json').read_bytes());rows=[]
cli=identity['binaries']['cli'];assert hashlib.sha256(pathlib.Path(cli['path']).read_bytes()).hexdigest()==cli['sha256']
for chapter in [5,1]:
 if shutil.disk_usage('/tmp').free<1024**3:raise SystemExit('disk below1GiB; stop without cleanup')
 assert not subprocess.check_output(['git','status','--porcelain'],cwd=root)
 args=['go','run','./cmd/grade-edition3','-ch',str(chapter),'-json']
 if chapter==5:args+=['-core',str(root),'-consumer',str(root/'examples/ch04-consumer'),'-workflow',str(root/'examples/ch06'),'-media-models',str(mapping)]
 args.append(cli['path']);start=time.monotonic();print('starting chapter',chapter,flush=True)
 with (f/('grade-ch%02d.json'%chapter)).open('w') as out,(f/('grade-ch%02d.stderr'%chapter)).open('w') as err:
  result=subprocess.run(args,cwd=main,env={**os.environ,'GIT_PAGER':'cat'},stdout=out,stderr=err)
 row={'chapter':chapter,'command':args,'cwd':str(main),'grading_main_head':subprocess.check_output(['git','rev-parse','HEAD'],cwd=main,text=True).strip(),'candidate_commit':identity['source_commit'],'cli_sha256':cli['sha256'],'workflow_sha256':identity['binaries']['workflow']['sha256'],'media_map_sha256':hashlib.sha256(mapping.read_bytes()).hexdigest(),'exit':result.returncode,'elapsed_seconds':round(time.monotonic()-start,3)}
 rows.append(row);(f/'grade-commands-and-exits.json').write_text(json.dumps(rows,indent=2)+'\n');print('chapter',chapter,'exit',result.returncode,flush=True)
 if result.returncode:raise SystemExit(result.returncode)
 assert not subprocess.check_output(['git','status','--porcelain'],cwd=root)
 assert hashlib.sha256(pathlib.Path(cli['path']).read_bytes()).hexdigest()==cli['sha256']
print('both exact-binary public routes completed')
