#!/usr/bin/env python3
"""Bounded independent compiled-route review; all writes are owned evidence/scratch."""
import pathlib,json,hashlib,subprocess,shutil,os,time,signal,threading,difflib,stat
R=pathlib.Path('/Users/bill/projects/ensemble'); S=pathlib.Path('/tmp/ch09-grader-review'); E=R/'docs/edition-3-notes/ch09-grader-review-evidence'
for n in ['home','tmp','bin','work','local-artifacts','guard']:(S/n).mkdir(exist_ok=True)
go='/usr/local/go/bin/go'; candidate=pathlib.Path('/tmp/cadence-ch09-grader')
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
env={'PATH':'/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(S/'home'),'TMPDIR':str(S/'tmp'),'GOCACHE':'/Users/bill/Library/Caches/go-build','GOMODCACHE':'/Users/bill/go/pkg/mod','GOPATH':'/Users/bill/go','GOWORK':'off','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOFLAGS':'-p=1 -mod=readonly','GOMAXPROCS':'2','LANG':'C','TZ':'UTC'}
(E/'environment.json').write_text(json.dumps({'env':env,'inheritance':'allowlist only; no external provider credentials or endpoint env copied','provider_network':'unchanged helpers create local fake endpoints and inject fake-key; no external providers','go':go},indent=2)+'\n')
records=[]
def write(name,obj):(E/name).write_text(json.dumps(obj,indent=2)+'\n')
def tree(p):return {str(f.relative_to(p)):{'bytes':f.stat().st_size,'sha256':sha(f)} for f in sorted(p.rglob('*')) if f.is_file()}
def digest(m):return hashlib.sha256(json.dumps(m,sort_keys=True,separators=(',',':')).encode()).hexdigest()
def run(name,args,cwd,limit=180,override=None):
 if shutil.disk_usage(S).free<1024**3:raise RuntimeError('below 1 GiB disk floor')
 args=list(map(str,args)); active=dict(env); active.update(override or {}); begin=time.time()
 with (E/(name+'.stdout')).open('wb') as out,(E/(name+'.stderr')).open('wb') as err:
  p=subprocess.Popen(args,cwd=cwd,env=active,stdout=out,stderr=err,start_new_session=True); timed=False
  try: code=p.wait(timeout=limit)
  except subprocess.TimeoutExpired:
   timed=True;os.killpg(p.pid,signal.SIGTERM)
   try:code=p.wait(timeout=5)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);code=p.wait()
 rec={'name':name,'argv':args,'cwd':str(cwd),'environment':'environment.json','overrides':override or {},'pid':p.pid,'pgid':p.pid,'exit_code':code,'outer_deadline_seconds':limit,'timeout':timed,'started_unix':begin,'duration_seconds':round(time.time()-begin,3)}
 for field in ['stdout','stderr']:
  f=E/(name+'.'+field);rec[field]={'path':f.name,'bytes':f.stat().st_size,'sha256':sha(f)}
 records.append(rec);write('commands.json',records);print(name,'exit',code,'seconds',rec['duration_seconds'],flush=True);return code
# Real source is copied unchanged; adapter dependencies are copied only for one route mutation.
C=S/'source';shutil.copytree(R/'solutions/ch09',C); source=tree(C);assert source==tree(R/'solutions/ch09')
(S/'local-artifacts/source-map.json').write_text(json.dumps(source,sort_keys=True,indent=2)+'\n')
A=S/'adapter';(A/'cmd/grade-edition3').mkdir(parents=True)
for f in ['go.mod','go.sum','cmd/grade-edition3/main.go','cmd/grade-edition3/main_test.go']:shutil.copy2(R/f,A/f)
for d in ['grade','gradeedition2','gradeedition3','fakevendor','fakeanthropic']:shutil.copytree(R/'internal'/d,A/'internal'/d,ignore=shutil.ignore_patterns('*.log','__pycache__'))
original_adapter=tree(A);(S/'local-artifacts/adapter-map.json').write_text(json.dumps(original_adapter,sort_keys=True,indent=2)+'\n')
write('source-identities.json',{'real_source_original':str(R/'solutions/ch09'),'real_copy':str(C),'real_files':len(source),'real_tree_sha256':digest(source),'adapter_copy':str(A),'adapter_files':len(original_adapter),'adapter_tree_sha256':digest(original_adapter),'module_paths':{n:str(C/d) for n,d in [('core',''),('consumer','ch05'),('workflow','ch06'),('stream','ch07'),('gui','gui')]},'explicit_executable':str(S/'bin/cli'),'source_identity_not_inferred_from_executable':True,'source_maps_local_only':['/tmp/ch09-grader-review/local-artifacts/source-map.json','/tmp/ch09-grader-review/local-artifacts/adapter-map.json']})
# Fail-fast setup. No retry path exists.
if run('build-real-cli',[go,'build','-o',S/'bin/cli','.'],C):raise SystemExit('CLI build failed')
sentinel=S/'guard/sentinel.go';sentinel.write_text('package main\nimport "os"\nfunc main(){if err:=os.WriteFile(os.Getenv("REVIEW_TARGET_MARKER"),[]byte("target launched\\n"),0600);err!=nil{panic(err)};os.Exit(73)}\n')
if run('build-sentinel',[go,'build','-o',S/'bin/sentinel',sentinel],S/'work'):raise SystemExit('sentinel build failed')
buildguard=S/'guard/go';buildguard.write_text('#!/bin/sh\nprintf "build invoked\\n" > "$REVIEW_BUILD_MARKER"\nexit 97\n');buildguard.chmod(0o700)
targetmark=S/'guard/target.marker';buildmark=S/'guard/build.marker';guardenv={'PATH':str(S/'guard')+':'+env['PATH'],'REVIEW_TARGET_MARKER':str(targetmark),'REVIEW_BUILD_MARKER':str(buildmark)}
assert run('target-marker-selfcheck',[S/'bin/sentinel'],S/'work',15,guardenv)==73 and targetmark.read_text()=='target launched\n';targetmark.rename(S/'local-artifacts/target-selfcheck.marker')
assert run('build-marker-selfcheck',[buildguard,'version'],S/'work',15,guardenv)==97 and buildmark.read_text()=='build invoked\n';buildmark.rename(S/'local-artifacts/build-selfcheck.marker')
module=S/'guard/module';module.mkdir();(module/'go.mod').write_text('module example.test/reviewguard\n\ngo 1.25\n');(module/'main.go').write_text('package main\nfunc main() {}\n')
nonregular=S/'guard/nonregular';(nonregular/'go.mod').mkdir(parents=True)
for n,text in [('malformed.json','{'),('bad-case.json','{"name":"x"}'),('bad-models.json','{"anthropic":""}')]: (S/'guard'/n).write_text(text)
base=['-ch','9','-core',str(module),'-consumer',str(module),'-workflow',str(module),'-stream-exercise',str(module),'-stream-model','claude-sonnet-5','-gui',str(module),'-gui-model','claude-sonnet-5']
cases=[('help',['-ch','9','-help'],0),('unsupported',['-ch','10'],2),('unknown-flag',['-ch','9','-unknown'],2),('older-ch6-gui',['-ch','6','-gui='],2)]
for flag in ['-core','-consumer','-workflow','-stream-exercise','-stream-model','-gui','-gui-model']:
 cases.append(('missing-'+flag[1:],[x for i in range(0,len(base),2) if base[i]!=flag for x in base[i:i+2]],2))
for flag in ['-stream-model','-gui-model']:cases.append(('blank-'+flag[1:],base+[flag,' \t\n'],2))
for flag,value,label in [('-core',S/'guard/absent','absent-core'),('-consumer',S/'guard','nonmodule-consumer'),('-gui',nonregular,'nonregular-gui')]:cases.append((label,base+[flag,str(value)],2))
for flag in ['-consumer-case','-media-models']:
 for suffix,path in [('malformed',S/'guard/malformed.json'),('missing',S/'guard/absent.json'),('schema',S/('guard/bad-case.json' if flag=='-consumer-case' else 'guard/bad-models.json'))]:cases.append((suffix+'-'+flag[1:],base+[flag,str(path)],2))
guards=[]
for name,flags,want in cases:
 for kind,target in [('executable',S/'bin/sentinel'),('build',module)]:
  assert not targetmark.exists() and not buildmark.exists(); code=run('guard-'+name+'-'+kind,[candidate]+flags+[str(target)],S/'work',15,guardenv)
  result={'case':name,'target_kind':kind,'exit_code':code,'expected_exit':want,'target_marker_absent':not targetmark.exists(),'build_marker_absent':not buildmark.exists()};guards.append(result);write('guard-results.json',guards)
  assert code==want and result['target_marker_absent'] and result['build_marker_absent'],result
for name,flags,want in [('no-target',base,2),('help-no-target',['-ch','9','-help'],0)]:
 code=run('guard-'+name,[candidate]+flags,S/'work',15,guardenv); result={'case':name,'target_kind':'none','exit_code':code,'expected_exit':want,'target_marker_absent':not targetmark.exists(),'build_marker_absent':not buildmark.exists()};guards.append(result);write('guard-results.json',guards);assert code==want and result['target_marker_absent'] and result['build_marker_absent']
# Observe actual generated executable files, without wrapping Go or changing launch behavior.
observed={};errors=[];stop=threading.Event();stage=['none']
def observe():
 while not stop.is_set():
  try:
   for directory in (S/'tmp').iterdir():
    if not directory.is_dir() or directory.name.startswith('go-build'):continue
    for f in directory.rglob('*'):
     try:
      if not f.is_file() or not f.stat().st_mode&stat.S_IXUSR:continue
      key=(str(f),f.stat().st_mtime_ns,f.stat().st_size)
      if key in observed:continue
      with f.open('rb') as h:magic=h.read(4)
      if magic not in [b'\xcf\xfa\xed\xfe',b'\xfe\xed\xfa\xcf',b'\xca\xfe\xba\xbe',b'\x7fELF']:continue
      data=f.read_bytes();h=hashlib.sha256(data).hexdigest(); dest=S/'local-artifacts'/('generated-'+h);dest.write_bytes(data)
      observed[key]={'path_at_observation':str(f),'retained':str(dest),'bytes':len(data),'sha256':h,'route_stage':stage[0]}
     except FileNotFoundError:pass
  except Exception as ex:errors.append(repr(ex))
  stop.wait(.01)
thread=threading.Thread(target=observe);thread.start()
ids9=['initial-tools','system-prompt','ensemble-tools','load-skill','progressive-disclosure','depends-autoload','var-substitution','blocked-skill','ch8-parity']
def args(ch):return ['-ch',str(ch),'-core',str(C),'-consumer',str(C/'ch05'),'-workflow',str(C/'ch06'),'-stream-exercise',str(C/'ch07'),'-stream-model','claude-sonnet-5','-gui',str(C/'gui'),'-gui-model','claude-sonnet-5','-json',str(S/'bin/cli')]
def inspect_report(name,want_ch):
 report=json.loads((E/(name+'.stdout')).read_text());checks=report.get('checks',[]);ids=[c['id'] for c in checks];title=report.get('title','');missing=[x for x in ids9 if x not in ids]
 matched=(title=='Edition 2, Chapter 9 — Skills' and ids==ids9) if want_ch==9 else (title=='Edition 2, Chapter 8 — Build Your Dream GUI' and 'ch7-parity' in ids and 'initial-tools' not in ids)
 result={'name':name,'expected_chapter':want_ch,'actual_title':title,'actual_ids':ids,'missing_ch9_ids':missing,'route_matches':matched,'report_passed':report.get('passed'),'all_checks_pass':all(c['passed'] for c in checks),'score':report.get('score'),'total':report.get('total'),'route_detector_exit':0 if matched else 1};write(name+'-route-detection.json',result);return result
try:
 stage[0]='frozen-ch09';assert run('frozen-ch09',[candidate]+args(9),S/'work',300)==0;result=inspect_report('frozen-ch09',9);assert result['route_matches'] and result['report_passed'] and result['all_checks_pass']
 stage[0]='frozen-ch08';assert run('frozen-ch08',[candidate]+args(8),S/'work',300)==0;result=inspect_report('frozen-ch08',8);assert result['route_matches'] and result['report_passed'] and result['all_checks_pass']
 main=A/'cmd/grade-edition3/main.go';original=main.read_bytes();old=b'report, err = gradeedition2.Chapter9(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)';new=old.replace(b'.Chapter9(',b'.Chapter8(');assert original.count(old)==1;mutated=original.replace(old,new);main.write_bytes(mutated)
 changed=[k for k,v in tree(A).items() if original_adapter.get(k)!=v];assert changed==['cmd/grade-edition3/main.go'];(E/'route-mutation.diff').write_text(''.join(difflib.unified_diff(original.decode().splitlines(True),mutated.decode().splitlines(True),fromfile='baseline/cmd/grade-edition3/main.go',tofile='mutant/cmd/grade-edition3/main.go')))
 write('route-mutation.json',{'only_changed_file':changed[0],'baseline_sha256':hashlib.sha256(original).hexdigest(),'mutant_sha256':sha(main),'baseline_adapter_tree':digest(original_adapter),'mutant_adapter_tree':digest(tree(A)),'behavioral_oracle':'Actual -ch 9 JSON must contain Chapter9 title and exact existing nine IDs, including skills and ch8-parity; a passing Chapter8 report fails this oracle.'})
 if run('build-route-mutant',[go,'build','-o',S/'bin/route-mutant','./cmd/grade-edition3'],A):raise SystemExit('mutant did not compile')
 stage[0]='route-mutant';assert run('route-mutant',[S/'bin/route-mutant']+args(9),S/'work',300)==0;result=inspect_report('route-mutant',9);assert not result['route_matches'] and result['report_passed'] and result['actual_title']=='Edition 2, Chapter 8 — Build Your Dream GUI' and len(result['missing_ch9_ids'])==9
 main.write_bytes(original);assert tree(A)==original_adapter;write('route-restoration.json',{'exact_all_files_restored':True,'source_tree_sha256':digest(tree(A)),'restored_main_sha256':sha(main),'baseline_main_sha256':hashlib.sha256(original).hexdigest()})
 if run('build-route-restored',[go,'build','-o',S/'bin/route-restored','./cmd/grade-edition3'],A):raise SystemExit('restored route did not compile')
 stage[0]='route-restored';assert run('route-restored',[S/'bin/route-restored']+args(9),S/'work',300)==0;result=inspect_report('route-restored',9);assert result['route_matches'] and result['report_passed'] and result['all_checks_pass']
finally:
 stop.set();thread.join();write('generated-artifacts.json',{'artifacts':list(observed.values()),'observer_errors':errors,'observer_joined':True,'limitations':'Sampling observer is not an exhaustive process tracer; unobserved transient artifacts are not invented.'})
assert tree(C)==source and tree(A)==original_adapter
write('artifact-identities.json',[{'role':name,'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)} for name,p in [('frozen-candidate',candidate),('real-cli',S/'bin/cli'),('target-sentinel',S/'bin/sentinel'),('build-marker-wrapper',buildguard),('route-mutant',S/'bin/route-mutant'),('route-restored',S/'bin/route-restored')]])
write('completion.json',{'guards':len(guards),'real_source_unchanged':True,'adapter_exactly_restored':True,'all_controls_waited':True,'observer_joined':True,'no_retry':True,'no_deadline_fired':all(not r['timeout'] for r in records),'behavioral_controls':['frozen-ch09','frozen-ch08','route-mutant','route-restored']})
print('BOUNDED CONTROLS COMPLETE',flush=True)
