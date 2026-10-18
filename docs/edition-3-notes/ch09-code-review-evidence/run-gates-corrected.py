from pathlib import Path
import os,json,hashlib,subprocess,time,signal,shutil,datetime,difflib
R=Path('/Users/bill/projects/ensemble');S=Path('/tmp/ch09-code-review');E=R/'docs/edition-3-notes/ch09-code-review-evidence';C=S/'candidate';B=S/'bin';B.mkdir(exist_ok=True);H=S/'home-corrected';H.mkdir();TMP=S/'tmp-corrected';TMP.mkdir();sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
GO=shutil.which('go');env={'PATH':str(Path(GO).parent)+':/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(H),'TMPDIR':str(TMP)+'/', 'GOPATH':'/Users/bill/go','GOMODCACHE':'/Users/bill/go/pkg/mod','GOCACHE':'/Users/bill/Library/Caches/go-build','GOPROXY':'off','GOSUMDB':'off','GOTOOLCHAIN':'local','GOWORK':'off','GOFLAGS':'-mod=readonly -p=2','GOMAXPROCS':'2','LANG':'C','LC_ALL':'C'}
(E/'environment-corrected.json').write_text(json.dumps({'environment':env,'go_executable':GO,'go_executable_sha256':sha(Path(GO)),'config_isolation':'owned initially empty HOME; minimal explicit environment; no user keys or configuration inherited; cached dependencies only'},indent=2)+'\n');records=[]
def run(name,argv,cwd=C,deadline=120,expect=0):
 free=shutil.disk_usage(S).free
 if free<1024**3:raise RuntimeError('disk floor before '+name)
 start=time.monotonic();out=E/(name+'.stdout');err=E/(name+'.stderr')
 with out.open('wb') as a,err.open('wb') as b:
  p=subprocess.Popen([str(x) for x in argv],cwd=cwd,env=env,stdout=a,stderr=b,start_new_session=True);timed=False
  try:code=p.wait(timeout=deadline)
  except subprocess.TimeoutExpired:
   timed=True;os.killpg(p.pid,signal.SIGTERM)
   try:p.wait(timeout=5)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);p.wait()
   code=p.returncode
 row={'name':name,'argv':[str(x) for x in argv],'cwd':str(cwd),'environment':'environment-corrected.json','pid':p.pid,'pgid':p.pid,'deadline_seconds':deadline,'duration_seconds':time.monotonic()-start,'exit_code':code,'expected_exit':expect,'timeout':timed,'free_bytes_before':free,'free_bytes_after':shutil.disk_usage(S).free,'stdout':{'path':out.name,'bytes':out.stat().st_size,'sha256':sha(out)},'stderr':{'path':err.name,'bytes':err.stat().st_size,'sha256':sha(err)}};records.append(row);(E/'commands-corrected.json').write_text(json.dumps(records,indent=2)+'\n');print(name,'exit',code,'seconds',round(row['duration_seconds'],2),flush=True)
 if timed or code!=expect:raise RuntimeError('Gate outcome requires review: '+name)
 return row
run('root-test-corrected',[GO,'test','-count=1','./...'],deadline=480)
run('root-race',[GO,'test','-race','-count=1','./...'],deadline=480)
run('root-vet',[GO,'vet','./...'])
modules=[('gui',C/'gui','./cmd'),('consumer',C/'examples/ch04-consumer','.'),('workflow',C/'examples/ch06','.'),('stream',C/'examples/ch07','.')]
for name,work,target in modules:
 for mode,flags in [('test',['test','-count=1']),('race',['test','-race','-count=1']),('vet',['vet'])]:run(name+'-'+mode,[GO,*flags,'./...'],work)
run('build-cli',[GO,'build','-o',B/'cli','./cmd/ch01'])
for name,work,target in modules:run('build-'+name,[GO,'build','-o',B/name,target],work)
models=E/'media-models.json';models.write_text('{"anthropic":"claude-sonnet-4-6"}\n')
args=['/tmp/cadence-ch09-grader','-ch','9','-json','-core',C,'-consumer',C/'examples/ch04-consumer','-workflow',C/'examples/ch06','-stream-exercise',C/'examples/ch07','-stream-model','claude-sonnet-4-6','-gui',C/'gui','-gui-model','claude-sonnet-4-6','-media-models',models,B/'cli']
run('grade-ch09',args,deadline=180)
a=json.loads((E/'grade-ch09.stdout').read_text());assert a['passed'] and a['score']==a['max_score']==100
run('grade-ch01',['/tmp/cadence-ch09-grader','-ch','1','-json',B/'cli'],deadline=210)
a=json.loads((E/'grade-ch01.stdout').read_text());assert a['passed'] and a['score']==100 and all(x['passed'] for x in a['checks'] if x['id'] in ['generation_headers','generation_body'])
M=S/'pending-unload-mutant';shutil.copytree(C,M)
run('sensitivity-baseline',[GO,'test','-count=1','-v','-run','^TestChapter9CapabilityLifecycle$','.'],M)
f=M/'internal/skills/registry.go';before=f.read_bytes();old='\tr.pending[name] = true\n';new='\tr.pending[name] = true\n\tfor _, tool := range r.defs[name].tools {\n\t\tdelete(r.enabled, tool)\n\t}\n';assert before.decode().count(old)==1;f.write_text(before.decode().replace(old,new));mutated=f.read_bytes();(E/'sensitivity.diff').write_text(''.join(difflib.unified_diff(before.decode().splitlines(True),mutated.decode().splitlines(True),fromfile='original/internal/skills/registry.go',tofile='mutant/internal/skills/registry.go')))
run('sensitivity-compile',[GO,'test','-c','-o',B/'pending-mutant.test','.'],M)
run('sensitivity-negative',[B/'pending-mutant.test','-test.v','-test.run=^TestChapter9CapabilityLifecycle$'],M,expect=1)
negative=(E/'sensitivity-negative.stdout').read_text();assert 'pending unload removed capability' in negative and 'panic:' not in negative
f.write_bytes(before);maps={n:hashlib.sha256((M/n).read_bytes()).hexdigest() for n in [a['path'] for a in json.loads((S/'source-map.json').read_text())]};expected={a['path']:a['sha256'] for a in json.loads((S/'source-map.json').read_text())};assert maps==expected
run('sensitivity-restored',[GO,'test','-count=1','-v','-run','^TestChapter9CapabilityLifecycle$','.'],M)
(E/'sensitivity.json').write_text(json.dumps({'mutation':'prematurely remove pending-unload tool capability','file':'internal/skills/registry.go','original_sha256':hashlib.sha256(before).hexdigest(),'mutated_sha256':hashlib.sha256(mutated).hexdigest(),'restored_sha256':sha(f),'all132_source_files_exactly_restored':True,'negative_exit':1,'observed_diagnostic':'pending unload removed capability','baseline_exit':0,'restored_exit':0,'compiled_test_binary':{'path':str(B/'pending-mutant.test'),'bytes':(B/'pending-mutant.test').stat().st_size,'sha256':sha(B/'pending-mutant.test')}},indent=2)+'\n')
artifacts=[]
for name in ['cli','gui','consumer','workflow','stream','pending-mutant.test']:
 p=B/name;artifacts.append({'name':name,'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)})
(E/'artifacts.json').write_text(json.dumps(artifacts,indent=2)+'\n');print('BOUNDED GATES COMPLETE',flush=True)
