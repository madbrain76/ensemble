#!/usr/bin/env python3
import hashlib,json,os,pathlib,shutil,subprocess,tarfile,tempfile,time
root=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
out=pathlib.Path('/tmp/cadence-ch05-correction-r1-evidence/sensitivity');out.mkdir(exist_ok=False)
original=pathlib.Path('/tmp/cadence-ch05-evidence-r1/frozen')
identity=json.loads((original/'identity.json').read_text())
newtest=root/'chapter5_correction_test.go'
def digest(p):return hashlib.sha256(p.read_bytes()).hexdigest()
assert identity['source_commit']=='00489b05b65cfe5209f26cb655667df44b3847e2'
assert digest(pathlib.Path(identity['archive']['path']))==identity['archive']['sha256']
rows=[]
def run(name,args,cwd):
 if shutil.disk_usage('/tmp').free<1024**3:raise RuntimeError('disk below1GiB; no cleanup authorized')
 start=time.monotonic()
 with (out/(name+'.txt')).open('w') as log:
  p=subprocess.run(args,cwd=cwd,env={**os.environ,'GOWORK':'off','GIT_PAGER':'cat'},stdout=log,stderr=subprocess.STDOUT)
 row={'name':name,'command':args,'cwd':str(cwd),'exit':p.returncode,'elapsed_seconds':round(time.monotonic()-start,3)};rows.append(row)
 (out/'commands.json').write_text(json.dumps(rows,indent=2)+'\n');return p.returncode,(out/(name+'.txt')).read_text()
with tempfile.TemporaryDirectory(prefix='cadence-ch05-r1-original-') as temp:
 scratch=pathlib.Path(temp)
 with tarfile.open(identity['archive']['path']) as archive:archive.extractall(scratch,filter='data')
 for name,expected in identity['source_files'].items():assert digest(scratch/name)==expected,name
 shutil.copyfile(newtest,scratch/newtest.name)
 (out/'source-proof.json').write_text(json.dumps({'original_commit':identity['source_commit'],'original_archive_sha256':identity['archive']['sha256'],'all_original_files_verified':len(identity['source_files']),'only_added_file':newtest.name,'added_test_sha256':digest(newtest),'original_production_replaced_in_scratch_only':True},indent=2)+'\n')
 rc,text=run('original-compile',['go','test','-race','-run','^$','.'],scratch)
 if rc:raise SystemExit('original compile failed: INCONCLUSIVE')
 expected={'R1':'R1 successful text-only turn stranded its acknowledged hint','R2':'R2 cancelled Ask blocked on unrelated held request','R3':'R3 managed process cleanup blocked behind ordinary Go handler; actual process alive=true'}
 results=[]
 for finding,needle in expected.items():
  rc,text=run('original-'+finding,['go','test','-race','-count=1','-timeout','90s','-run','^TestChapter5Correction'+finding,'.'],scratch)
  environment=any(s in text for s in ['no space left','[build failed]','panic:','test timed out','DATA RACE','fatal error:'])
  detected=rc!=0 and needle in text and not environment
  results.append({'finding':finding,'behavioral_detection':detected,'exit':rc,'expected_assertion':needle,'environment_or_unrelated_failure':environment,'assertion_occurrences':text.count(needle)})
  (out/'results.json').write_text(json.dumps(results,indent=2)+'\n')
  if not detected:raise SystemExit(finding+' sensitivity inconclusive/failed; inspect receipt')
 rc,text=run('corrected-focused',['go','test','-race','-count=1','-timeout','90s','-run','^TestChapter5Correction','.'],root)
 if rc:raise SystemExit('corrected focused suite failed')
print(json.dumps({'original_behavioral_detections':results,'corrected_exit':rc},indent=2))
