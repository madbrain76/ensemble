from pathlib import Path
import os,json,hashlib,subprocess,time,signal,shutil,tarfile,difflib
R=Path('/Users/bill/projects/ensemble');T=Path('/Users/bill/projects/ensemble-edition-3-student');S=Path('/tmp/ch09-code-review');C=S/'candidate';B=S/'bin';E=R/'docs/edition-3-notes/ch09-code-review-evidence';GO='/usr/local/go/bin/go'
sha=lambda p:hashlib.sha256(Path(p).read_bytes()).hexdigest()
def save(n,x):
 p=E/n;assert not p.exists();p.write_text(json.dumps(x,indent=2)+'\n')
# This continuation has no completed suite or application build command.
old=json.loads((E/'commands-corrected.json').read_text());assert len(old)==21 and old[-1]['name']=='grade-ch09' and old[-1]['exit_code']==2
assert not (E/'commands-continuation.json').exists()
source=json.loads((R/'docs/edition-3-notes/ch09-student-source-manifest.json').read_text());archive=Path(source['source_archive']['path']);assert sha(archive)==source['source_archive']['sha256']
models=R/'docs/edition-3-notes/ch09-student-evidence/media-models.json';assert models.stat().st_size==92 and sha(models)=='cc75c8f40899ce7dd6046f26c413c7afff518425bb86c72b853306966de9796b';assert json.loads(models.read_text())==source['models']
grader=Path('/tmp/cadence-ch09-grader');assert sha(grader)==source['grader']['sha256']
records=json.loads((S/'source-map.json').read_text());expected={x['path']:(x['bytes'],x['sha256']) for x in records}
for n,(size,digest) in expected.items():assert (C/n).stat().st_size==size and sha(C/n)==digest,n
H=S/'home-continuation';TMP=S/'tmp-continuation';H.mkdir();TMP.mkdir()
env=json.loads((E/'environment-corrected.json').read_text())['environment'];env.update(HOME=str(H),TMPDIR=str(TMP)+'/',GOFLAGS='-mod=readonly -p=1');assert env['GOMAXPROCS']=='2' and all(not k.startswith(('LLM_','ANTHROPIC_','OPENAI_','GEMINI_','GOOGLE_')) for k in env)
save('environment-continuation.json',{'environment':env,'go_executable':GO,'go_executable_sha256':sha(GO),'isolation':'new empty owned HOME/TMPDIR; no provider variables, user config or key reads; offline local toolchain; helper fake servers provide dummy credentials'})
artifacts=[{'name':p.name,'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)} for p in sorted(B.iterdir())];assert len(artifacts)==5
args=[str(grader),'-ch','9','-json','-core',str(C),'-consumer',str(C/'examples/ch04-consumer'),'-workflow',str(C/'examples/ch06'),'-stream-exercise',str(C/'examples/ch07'),'-stream-model','claude-sonnet-4-6','-gui',str(C/'gui'),'-gui-model','claude-sonnet-4-6','-media-models',str(models),str(B/'cli')]
chapter1=[str(grader),'-ch','1','-json',str(B/'cli')]
for p in [C,C/'gui',C/'examples/ch04-consumer',C/'examples/ch06',C/'examples/ch07']:assert (p/'go.mod').is_file()
save('continuation-preflight.json',{'authorization':'main-16, resumed after provider transport failure seq758; prior receipts reconciled before execution','completed_before_continuation':'15 corrected normal/race/vet gates plus5 application builds; none rerun','prior_failures':'initial contaminated root-test exit1 and media-map validation exit2 preserved','remaining_runtime_controls_not_started_before_resume':True,'grader':{'path':str(grader),'bytes':grader.stat().st_size,'sha256':sha(grader)},'exact_existing_media_map':{'path':str(models),'bytes':models.stat().st_size,'sha256':sha(models),'content':json.loads(models.read_text())},'argv_chapter9':args,'argv_chapter1':chapter1,'module_and_flag_review':'cmd/grade-edition3/main.go validation/dispatch and internal/gradeedition2/ch05_media.go ReadMediaModels; all flags precede same existing CLI; core is separate source path; consumer-case uses accepted shout default','reused_artifacts':artifacts,'source_files_verified':len(expected),'archive':{'path':str(archive),'bytes':archive.stat().st_size,'sha256':sha(archive)},'chapter1_allowance_seconds':210,'phase_requirement':'unchanged genuine62-second headers and62-second body; no artificial request cap or shortened timing','free_bytes':shutil.disk_usage(S).free})
commands=[]
def run(name,argv,cwd=C,expected_exit=0,deadline=180):
 assert shutil.disk_usage(S).free>=1024**3
 out=E/(name+'.stdout');err=E/(name+'.stderr');assert not out.exists() and not err.exists();start=time.monotonic();free=shutil.disk_usage(S).free
 with out.open('wb') as o,err.open('wb') as e:
  p=subprocess.Popen(list(map(str,argv)),cwd=cwd,env=env,stdout=o,stderr=e,start_new_session=True);timed=False
  try:code=p.wait(timeout=deadline)
  except subprocess.TimeoutExpired:timed=True;os.killpg(p.pid,signal.SIGKILL);code=p.wait()
 record={'name':name,'argv':list(map(str,argv)),'cwd':str(cwd),'environment':'environment-continuation.json','pid':p.pid,'pgid':p.pid,'deadline_seconds':deadline,'duration_seconds':time.monotonic()-start,'exit_code':code,'expected_exit':expected_exit,'timeout':timed,'free_bytes_before':free,'free_bytes_after':shutil.disk_usage(S).free,'stdout':{'path':out.name,'bytes':out.stat().st_size,'sha256':sha(out)},'stderr':{'path':err.name,'bytes':err.stat().st_size,'sha256':sha(err)}};commands.append(record);(E/'commands-continuation.json').write_text(json.dumps(commands,indent=2)+'\n');print(name,'exit',code,'seconds',round(record['duration_seconds'],2),flush=True)
 assert not timed and code==expected_exit,'Gate outcome requires review: '+name
run('runtime-ch09',args)
r=json.loads((E/'runtime-ch09.stdout').read_text());assert r['passed'] and r['score']==r['max_score']==100 and len(r['checks'])==9 and all(x['passed'] for x in r['checks'])
run('runtime-ch01',chapter1,deadline=210)
r=json.loads((E/'runtime-ch01.stdout').read_text());assert r['passed'] and r['score']==r['max_score']==100 and all(x['passed'] for x in r['checks']);assert commands[-1]['duration_seconds']>=124
M=S/'sensitivity-archive';M.mkdir()
with tarfile.open(archive) as tar:
 members={m.name.removeprefix('./'):m for m in tar.getmembers() if m.isfile()};assert set(members)==set(expected)
 for n,(size,digest) in expected.items():
  assert not Path(n).is_absolute() and '..' not in Path(n).parts
  data=tar.extractfile(members[n]).read();assert len(data)==size and hashlib.sha256(data).hexdigest()==digest
  dest=M/n;dest.parent.mkdir(parents=True,exist_ok=True);dest.write_bytes(data)
run('sensitivity-baseline',[GO,'test','-count=1','-v','-run','^TestChapter9CapabilityLifecycle$','.'],M)
f=M/'internal/skills/registry.go';before=f.read_bytes();old='\tr.pending[name] = true\n';new='\tr.pending[name] = true\n\tfor _, tool := range r.defs[name].tools {\n\t\tdelete(r.enabled, tool)\n\t}\n';assert before.decode().count(old)==1;mutated=before.decode().replace(old,new).encode();f.write_bytes(mutated)
(E/'sensitivity.diff').write_text(''.join(difflib.unified_diff(before.decode().splitlines(True),mutated.decode().splitlines(True),fromfile='original/internal/skills/registry.go',tofile='mutant/internal/skills/registry.go')))
save('sensitivity-mutation.json',{'source':'pinned archive extracted into owned sensitivity-archive; all132 relative paths/bytes verified','file':'internal/skills/registry.go','original_sha256':hashlib.sha256(before).hexdigest(),'mutant_sha256':hashlib.sha256(mutated).hexdigest(),'change':'premature capability removal immediately on pending unload; tests unchanged'})
run('sensitivity-compile',[GO,'test','-c','-o',B/'pending-mutant.test','.'],M)
run('sensitivity-negative',[B/'pending-mutant.test','-test.v','-test.run=^TestChapter9CapabilityLifecycle$'],M,expected_exit=1)
negative=(E/'sensitivity-negative.stdout').read_text();assert 'pending unload removed capability' in negative and 'panic:' not in negative
f.write_bytes(before)
for n,(size,digest) in expected.items():assert (M/n).stat().st_size==size and sha(M/n)==digest,n
run('sensitivity-restored',[GO,'test','-count=1','-v','-run','^TestChapter9CapabilityLifecycle$','.'],M)
save('sensitivity.json',{'mutation':'premature removal of pending-unload capability','file':'internal/skills/registry.go','original_sha256':hashlib.sha256(before).hexdigest(),'mutant_sha256':hashlib.sha256(mutated).hexdigest(),'restored_sha256':sha(f),'all132_source_files_exactly_restored':True,'test_source_unchanged':True,'baseline_exit':0,'compile_exit':0,'negative_exit':1,'observed_diagnostic':'pending unload removed capability','restored_exit':0,'compiled_mutant_test_binary':{'path':str(B/'pending-mutant.test'),'bytes':(B/'pending-mutant.test').stat().st_size,'sha256':sha(B/'pending-mutant.test')}})
for a in artifacts:assert sha(a['path'])==a['sha256']
save('reused-artifacts-after.json',artifacts);print('AUTHORIZED REMAINING CONTROLS COMPLETE',flush=True)
