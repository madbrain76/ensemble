#!/usr/bin/env python3
"""Root-authorized finite missing-control follow-up. No automatic retries."""
from pathlib import Path
import difflib,hashlib,json,os,re,shutil,signal,subprocess,tempfile,threading,time
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');E=R/'docs/edition-3-notes/ch10-grader-review-followup-evidence';O=R/'docs/edition-3-notes/ch10-grader-review-evidence';M=json.loads((R/'docs/edition-3-notes/ch10-grader-review-input.json').read_text())
W=Path(tempfile.mkdtemp(prefix='cadence-ch10-review-followup-',dir='/tmp'));receipts=[];phase='preparation';stop=threading.Event();observed={};processes=set();run_outcome={'verdict':'not completed','controls':{},'failure':None};copy_created=False;copied_restored=False
for n in ['home','tmp','gopath','bin','fake-bin','fixtures']:(W/n).mkdir()
def ident(p):
 b=Path(p).read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def write(n,x):(E/n).write_text(json.dumps(x,indent=2)+'\n')
def git(root,*args):return subprocess.check_output(['git','--no-pager',*args],cwd=root)
def treepins(root):
 assert not any(p.is_symlink() for p in root.rglob('*'))
 return {p.relative_to(root).as_posix():ident(p) for p in sorted(root.rglob('*')) if p.is_file()}
def checkpins(root,pins):
 for rel,want in pins.items():assert ident(root/rel)==want,(str(root),rel)
def verify():
 v={'pinned_groups':{},'prior_seals':{},'repositories':{},'candidate_binary':ident(M['candidate']['binary']['path'])}
 assert v['candidate_binary']=={k:M['candidate']['binary'][k] for k in ['bytes','sha256']}
 for label,root,pins in [('inputs',R,M['inputs']),('protected_tracked',R,M['protected_tracked']),('predecessor',S,M['predecessor']['files'])]:
  checkpins(root,pins);v['pinned_groups'][label]={'count':len(pins),'all_match':True}
 checkpins(R,M['candidate']['source']);assert ident(R/'internal/gradeedition2/ch10.go')==M['candidate']['historical_grader_unchanged']
 for dirname in ['ch10-entry-audit-evidence','ch10-entry-followup-evidence','ch10-grader-review-evidence']:
  d=R/'docs/edition-3-notes'/dirname;inv=d/'inventory.json';j=json.loads(inv.read_text())
  for row in j['files']:assert ident(R/row['path'])=={k:row[k] for k in ['bytes','sha256']},row['path']
  outside={row['path'] for row in j['files'] if not row['path'].startswith(d.relative_to(R).as_posix()+'/')};actual={p.relative_to(R).as_posix() for p in d.rglob('*') if p.is_file() and p!=inv}|outside;assert actual=={row['path'] for row in j['files']}
  v['prior_seals'][dirname]={'inventory':ident(inv),'file_count':len(j['files']),'exact_closed_set':True}
 for name,root in [('main',R),('student',S)]:
  d=git(root,'diff','--binary');c=git(root,'diff','--cached','--binary');s=git(root,'status','--porcelain=v1','--untracked-files=no')
  v['repositories'][name]={'head':git(root,'rev-parse','HEAD').decode().strip(),'tree':git(root,'rev-parse','HEAD^{tree}').decode().strip(),'unstaged_diff':{'bytes':len(d),'sha256':hashlib.sha256(d).hexdigest()},'staged_diff':{'bytes':len(c),'sha256':hashlib.sha256(c).hexdigest()},'tracked_status':{'bytes':len(s),'sha256':hashlib.sha256(s).hexdigest()}}
 assert git(S,'status','--porcelain=v1','--untracked-files=normal')==b''
 assert v['repositories']['student']['head']==M['predecessor']['commit'];assert v['repositories']['student']['tree']==M['predecessor']['tree']
 assert git(R,'diff','--','cmd/grade-edition3/main.go','cmd/grade-edition3/main_test.go')==(R/'docs/edition-3-notes/ch10-grader-evidence/candidate.diff').read_bytes()
 v['disk_free_bytes']=shutil.disk_usage(R).free;assert v['disk_free_bytes']>=1024**3
 return v
ENV=json.loads((O/'environment.json').read_text());ENV.update(HOME=str(W/'home'),TMPDIR=str(W/'tmp'),GOPATH=str(W/'gopath'))
assert ENV['GOMAXPROCS']=='2' and ENV['GOFLAGS']=='-mod=readonly -p=1' and ENV['GOTOOLCHAIN']=='local' and ENV['GOPROXY']=='off' and ENV['GOSUMDB']=='off';assert ENV['PATH'].startswith('/Users/bill/go/bin:');assert 'LLM_BASE_URL' not in ENV
GUARD_ENV=dict(ENV);GUARD_ENV['PATH']=str(W/'fake-bin')+':'+ENV['PATH']
write('environment.json',ENV);write('environment-build-guard.json',GUARD_ENV);write('workspace.json',{'owned_scratch':str(W),'existing_file_edits_forbidden':True})
write('task-and-input.json',{'task':ident('/tmp/cadence-ch10-grader-review-followup-task.txt'),'review_input':ident(R/'docs/edition-3-notes/ch10-grader-review-input.json'),'carried_read_rule':'Previous full reads carried only after all exact original input pins and sealed evidence reverified; new task fully read. Actual ch09 module paths and CLI entry inspected.'})
assert ident('/tmp/cadence-ch10-grader-review-followup-task.txt')=={'bytes':5498,'sha256':'4145951c34a75dcbf3657718687907b190e3d8a11b35c3b8fc623b555c13006d'}
before=verify();write('verification-before.json',before)
refpins={ch:treepins(R/'solutions'/ch) for ch in ['ch09','ch10']};write('reference-source-pins.json',refpins)
for ch in ['ch09','ch10']:
 for rel in ['go.mod','gui/go.mod','ch05/go.mod','ch06/go.mod','ch07/go.mod']:assert (R/'solutions'/ch/rel).is_file()
go=Path('/usr/local/go/bin/go');dlv=Path('/Users/bill/go/bin/dlv');assert ident(dlv)=={'bytes':20038402,'sha256':'10dbb1ce8e6d2cb77ae4d30c7ca1db0302b83529cd6cd39453f009a7653c17fe'}
write('toolchain-pins.json',{str(go):ident(go),str(dlv):ident(dlv)})
# Compute local build import closure without running an extra Go campaign.
prefix='github.com/waywardgeek/ensemble/';todo=['cmd/grade-edition3'];seen=set();copyfiles={'go.mod','go.sum','cmd/grade-edition3/main_test.go'}
while todo:
 rel=todo.pop()
 if rel in seen:continue
 seen.add(rel)
 for p in sorted((R/rel).glob('*.go')):
  if p.name.endswith('_test.go'):continue
  copyfiles.add(p.relative_to(R).as_posix());text=p.read_text();head=re.split(r'(?m)^(?:func|type|var|const)\b',text,maxsplit=1)[0]
  for block in re.findall(r'(?ms)^import\s+(?:\([^)]*\)|"[^"\n]+")',head):
   for path in re.findall(r'"([^"]+)"',block):
    if path.startswith(prefix):todo.append(path[len(prefix):])
  for line in re.findall(r'(?m)^//go:embed\s+([^\n]+)$',text):
   for pat in line.split():
    for f in p.parent.glob(pat):assert f.is_file();copyfiles.add(f.relative_to(R).as_posix())
copy_pins={rel:ident(R/rel) for rel in sorted(copyfiles)};write('copy-original-pins.json',{'root':str(R),'local_packages':sorted(seen),'file_count':len(copy_pins),'files':copy_pins});assert len(copy_pins)==134
COPY=W/'adapter-copy';mainrel='cmd/grade-edition3/main.go';needle='report, err = gradeedition2.Chapter10(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)';replacement=needle.replace('Chapter10(', 'Chapter9(');assert (R/mainrel).read_text().count(needle)==1
for n in ['media-models.json','consumer-case.json']:(E/n).write_bytes((O/n).read_bytes())
case=E/'consumer-case.json';media=E/'media-models.json';models=json.loads(media.read_text());assert set(models)=={'anthropic','openai','gemini'};assert all(v.strip() for v in models.values())
assert ident(case)==ident(S/'examples/ch04-consumer/consumer-case.json')
ids10=['save-shape','default-load','replay-equals-snapshot','tail-applied-once','log-not-needed','bad-save-refused','ch9-parity'];ids9=['initial-tools','system-prompt','ensemble-tools','load-skill','progressive-disclosure','depends-autoload','var-substitution','blocked-skill','ch8-parity'];title10='Edition 2, Chapter 10 — Persistence';title9='Edition 2, Chapter 9 — Skills'
def identity(r,expected_title=title10,expected_ids=ids10):
 actual_ids=[c['id'] for c in r['checks']];ok=r.get('title')==expected_title and actual_ids==expected_ids
 return {'expected_title':expected_title,'expected_ids':expected_ids,'actual_title':r.get('title'),'actual_ids':actual_ids,'identity_matches':ok,'calculated_detector_status':0 if ok else 1,'status_is_not_process_exit':True}
saved=json.loads((O/'07-ch10-positive.stdout.txt').read_text());saved_receipt=next(c for c in json.loads((O/'command-receipts.json').read_text()) if c['name']=='07-ch10-positive');baseline_detector=identity(saved);assert baseline_detector['identity_matches'] and saved['passed'] is True and saved['score']==saved['max_score']==100 and saved_receipt['exit']==0
write('carried-ch10-baseline.json',{'source':'immutable prior successful 07-ch10-positive.stdout.txt','source_report':ident(O/'07-ch10-positive.stdout.txt'),'original_receipt':saved_receipt,'candidate_source_and_binary_reverified':True,'detector':baseline_detector,'no_fresh_baseline_run':True})
# Compact representative guard fixtures; no real target runs.
f=W/'fixtures';valid=f/'valid-module';valid.mkdir();(valid/'go.mod').write_text('module example.test/buildguard\n\ngo 1.25\n');(valid/'main.go').write_text('package main\nfunc main(){panic("real guard target must never execute")}\n')
nonregular=f/'nonregular';(nonregular/'go.mod').mkdir(parents=True);bad=f/'malformed.json';bad.write_text('{')
marker=W/'build-entry-marker';fakego=W/'fake-bin/go';fake_source=f/'fake-go.go'
fake_text='package main\nimport ("encoding/json"; "fmt"; "os")\nfunc main(){wd,_:=os.Getwd(); event:=map[string]any{"argv":os.Args,"cwd":wd,"pid":os.Getpid(),"ppid":os.Getppid(),"declared_exit":73}; b,err:=json.Marshal(event);if err!=nil{panic(err)};if err=os.WriteFile('+json.dumps(str(marker))+',b,0600);err!=nil{panic(err)};fmt.Println("owned fake go reached: exiting 73");os.Exit(73)}\n'
fake_source.write_text(fake_text);(E/'fake-go-source.txt').write_text(fake_text)
base=['-ch','10','-core',str(valid),'-consumer',str(valid),'-workflow',str(valid),'-stream-exercise',str(valid),'-stream-model','declared-stream','-gui',str(valid),'-gui-model','declared-gui'];matrix=[]
def add(name,args,exit=2):matrix.append({'name':name,'args':args,'expected_adapter_exit':exit})
add('missing-core',[v for i in range(0,len(base),2) if base[i]!='-core' for v in base[i:i+2]]+[str(valid)])
for name,path in [('missing-module',f/'absent'),('nonmodule',f),('nonregular-go-mod',nonregular)]:add(name,base+['-core',str(path),str(valid)])
for flag in ['-stream-model','-gui-model']:add('blank'+flag,base+[flag,' \t\n',str(valid)])
for flag in ['-consumer-case','-media-models']:add('malformed'+flag,base+[flag,str(bad),str(valid)])
add('unsupported-11',['-ch','11',str(valid)]);add('help-10',['-ch','10','-h',str(valid)],0)
for ch,flag in [('1','-core'),('4','-workflow'),('5','-stream-model'),('6','-gui-model')]:add('older'+ch+flag+'-explicit-empty',['-ch',ch,flag,'',str(valid)])
assert len(matrix)==14;write('build-guard-plan.json',matrix)
write('control-plan.json',{'order':['ONE actual ch09 complete regression with frozen binary','compile fake go; direct positive marker check; valid directory Build-entry positive check; 14 invalid/help no-Build checks','compile adapter-copy Chapter10-to-Chapter9 mutation; ONE actual report and independent wrong-identity rejection','restore EVERY copied file byte-for-byte, rebuild, ONE restored full Chapter10 pass'],'no_fresh_ch10_baseline':True,'no_package_race_vet_or_root_build':True,'no_historical_app_mutation':True,'full_deadline_seconds':600,'build_deadline_seconds':180,'guard_deadline_seconds':30,'no_retries':True,'model_map':models,'candidate_binary':M['candidate']['binary'],'copy_file_count':len(copy_pins),'mutation':{'file':mainrel,'old':needle,'new':replacement},'original_review_failure_preserved':True})

def sample():
 while not stop.wait(.12):
  try:
   # Sampling is non-exhaustive. Local binary/config/save bodies never exported.
   for root,dirs,files in os.walk(W/'tmp'):
    dirs[:]=[d for d in dirs if not d.startswith('go-build')]
    for name in files:
     if name not in ['submission','settings.json','save.json']:continue
     p=Path(root)/name
     try:
      h=ident(p);key=str(p)+'#'+h['sha256']
      if key not in observed:observed[key]={'local_path':str(p),'kind':'executable' if name=='submission' else 'local-config-or-save','phase':phase,'observed_unix':time.time(),**h}
     except OSError:pass
   processes.update(x.strip() for x in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True).splitlines() if str(W) in x)
  except (OSError,subprocess.SubprocessError):pass

def run(name,argv,allowed,timeout,cwd=R,env=ENV):
 global phase
 phase=name;assert shutil.disk_usage(R).free>=1024**3;allowed=[allowed] if isinstance(allowed,int) else allowed;start=time.time();p=subprocess.Popen([str(x) for x in argv],cwd=cwd,env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE,start_new_session=True);timed=False
 try:out,err=p.communicate(timeout=timeout)
 except subprocess.TimeoutExpired:
  timed=True;os.killpg(p.pid,signal.SIGTERM)
  try:out,err=p.communicate(timeout=5)
  except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);out,err=p.communicate()
 for suffix,data in [('stdout',out),('stderr',err)]:assert b'\0' not in data;data.decode('utf-8');(E/f'{name}.{suffix}.txt').write_bytes(data)
 row={'name':name,'argv':[str(x) for x in argv],'cwd':str(cwd),'environment':'environment-build-guard.json' if env is GUARD_ENV else 'environment.json','pid':p.pid,'start_unix':start,'elapsed_seconds':round(time.time()-start,3),'deadline_seconds':timeout,'process_exit':p.returncode,'allowed_process_exits':allowed,'timed_out':timed,'stdout':{'path':f'{name}.stdout.txt',**ident(E/f'{name}.stdout.txt')},'stderr':{'path':f'{name}.stderr.txt',**ident(E/f'{name}.stderr.txt')},'executable_identity':ident(argv[0]),'disk_free_bytes_after':shutil.disk_usage(R).free};receipts.append(row);write('command-receipts.json',receipts);print(name,'process exit',p.returncode,'seconds',row['elapsed_seconds'],flush=True)
 assert not timed and p.returncode in allowed,{'unexpected_command_result':row}
 return out,err,p.returncode

def report_run(name,bin,ch,ref,mode):
 args=[bin,'-ch',str(ch),'-json','-core',ref,'-consumer',ref/'ch05','-consumer-case',case,'-workflow',ref/'ch06','-media-models',media,'-stream-exercise',ref/'ch07','-stream-model',models['anthropic'],'-gui',ref/'gui','-gui-model',models['anthropic'],ref]
 out,err,exit=run(name,args,[0,1] if mode=='mutant' else 0,600);r=json.loads(out);write(name+'.report.json',r)
 detector=identity(r,title9,ids9) if mode=='ch9' else identity(r);receipt={'actual_process_exit':exit,'report_passed':r.get('passed'),'report_score':r.get('score'),'report_max_score':r.get('max_score'),'detector':detector};write(name+'.detector.json',receipt)
 if mode=='mutant':
  assert r.get('title')==title9 and [c['id'] for c in r['checks']]==ids9
  assert detector['identity_matches'] is False and detector['calculated_detector_status']==1
 else:
  assert detector['identity_matches'] and r['passed'] is True and r['score']==r['max_score']==100 and all(c['passed'] is True for c in r['checks']) and exit==0
 run_outcome['controls'][name]=receipt;write('run-outcome.json',run_outcome)

def copycheck(expected):
 current=treepins(COPY);assert current==expected,'copied file set/content differs from pins';return current

def restore_copy():
 global copied_restored
 checkpins(R,copy_pins)
 for rel in copy_pins:shutil.copy2(R/rel,COPY/rel)
 copycheck(copy_pins);copied_restored=True

sampler=threading.Thread(target=sample,daemon=True);sampler.start()
try:
 report_run('01-ch9-real-reference',Path(M['candidate']['binary']['path']),9,R/'solutions/ch09','ch9')
 run('02-build-fake-go',[go,'build','-o',fakego,fake_source],0,180)
 write('fake-go-binary.json',{'path':str(fakego),**ident(fakego),'source':ident(fake_source)})
 assert not marker.exists();out,err,exit=run('03-fake-go-direct-positive',[fakego,'direct-positive-selfcheck'],73,30,env=GUARD_ENV);event=json.loads(marker.read_text());write('03-marker-positive.json',{'marker_existed':True,'marker':event,'marker_identity':ident(marker),'observed_direct_fake_go_process_exit':exit});marker.unlink()
 out,err,exit=run('04-valid-directory-reaches-build',[M['candidate']['binary']['path'],'-ch','1',valid],2,30,env=GUARD_ENV);assert marker.exists();event=json.loads(marker.read_text());assert event['argv'][1]=='build' and event['cwd']==str(valid)
 write('04-build-entry-positive.json',{'adapter_process_exit':exit,'actual_build_marker':True,'marker':event,'marker_identity':ident(marker),'fake_go_declared_exit':73,'indirect_child_exit_not_separately_waited':'same compiled fixture unconditionally exits73; direct exit73 observed separately; adapter exit2 is NOT fake-go exit73','real_target_executed':False});marker.unlink()
 results=[]
 for i,c in enumerate(matrix,1):
  assert not marker.exists();out,err,exit=run(f'05-guard-{i:02}',[M['candidate']['binary']['path']]+c['args'],c['expected_adapter_exit'],30,env=GUARD_ENV);assert not marker.exists();results.append({**c,'adapter_process_exit':exit,'build_marker_absent':True,'actual_target':'owned valid directory; fake Go intercepts any build'});write('build-guard-results.json',results)
 run_outcome['controls']['build-entry-guard']={'direct_fake_go_process_exit':73,'valid_directory_adapter_process_exit':2,'valid_directory_build_marker':True,'invalid_help_cases':len(results),'all_invalid_help_build_markers_absent':True};write('run-outcome.json',run_outcome)
 # The sole mutation: one typed adapter dispatch call. Graders/tests/apps unchanged.
 checkpins(R,copy_pins)
 for rel in copy_pins:
  target=COPY/rel;target.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(R/rel,target)
 copy_created=True;copycheck(copy_pins);write('copy-before.json',{'root':str(COPY),'files':treepins(COPY)})
 original=(COPY/mainrel).read_text();mutated=original.replace(needle,replacement);(COPY/mainrel).write_text(mutated);mutated_pins=dict(copy_pins);mutated_pins[mainrel]=ident(COPY/mainrel);copycheck(mutated_pins)
 (E/'adapter-only-mutation.diff.txt').write_text(''.join(difflib.unified_diff(original.splitlines(True),mutated.splitlines(True),fromfile='original/'+mainrel,tofile='mutant/'+mainrel)))
 write('copy-mutated.json',{'root':str(COPY),'changed_files':[mainrel],'files':mutated_pins,'historical_app_mutated':False})
 mutant_bin=W/'bin/adapter-route-mutant';run('06-compile-adapter-mutant',[go,'build','-o',mutant_bin,'./cmd/grade-edition3'],0,180,cwd=COPY);copycheck(mutated_pins);write('mutant-binary.json',{'path':str(mutant_bin),**ident(mutant_bin),'source_manifest':'copy-mutated.json'})
 report_run('07-mutant-ch10-wrong-route',mutant_bin,10,R/'solutions/ch10','mutant')
 restore_copy();write('copy-restored-before-build.json',{'root':str(COPY),'all_134_files_exact_original':True,'files':treepins(COPY)})
 restored_bin=W/'bin/adapter-route-restored';run('08-compile-restored-adapter',[go,'build','-o',restored_bin,'./cmd/grade-edition3'],0,180,cwd=COPY);copycheck(copy_pins);write('restored-binary.json',{'path':str(restored_bin),**ident(restored_bin),'source_manifest':'copy-restored-before-build.json'})
 report_run('09-restored-ch10',restored_bin,10,R/'solutions/ch10','restored');copycheck(copy_pins)
 run_outcome['verdict']='PASS: missing assigned controls completed, subject to root inspection/join'
except Exception as ex:
 run_outcome['verdict']='NOT ACCEPTED: stopped on unexpected result, no retry';run_outcome['failure']=type(ex).__name__+': '+str(ex);print(run_outcome['failure'],flush=True)
finally:
 if copy_created and not copied_restored:restore_copy()
 stop.set();sampler.join();write('run-outcome.json',run_outcome)
 write('observed-local-artifacts.json',{'coverage':'120ms sampling, NON-EXHAUSTIVE; local binaries/settings/save bodies NOT exported','artifacts':list(observed.values())});(E/'observed-processes.txt').write_text('\n'.join(sorted(processes))+'\n')
 for ch,pins in refpins.items():assert treepins(R/'solutions'/ch)==pins,ch
 checkpins(R,copy_pins)
 if copy_created:copycheck(copy_pins);write('copy-final-restoration.json',{'root':str(COPY),'all_original_files_restored':True,'files':treepins(COPY),'no_extra_files':True})
 after=verify();write('verification-after.json',after);assert before['repositories']==after['repositories'];assert before['prior_seals']==after['prior_seals']
 live=[x.strip() for x in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True).splitlines() if str(W) in x];assert not live,live
 write('completion.json',{'all_jobs_waited':True,'owned_processes_remaining':live,'candidate_and_historical_sources_unchanged':True,'copy_created':copy_created,'copy_restored':copied_restored,'owned_scratch_retained':str(W),'prior_seals_unchanged':True,'predecessor_clean':True,'no_shared_cleanup':True,'disk_floor_bytes':1024**3,'disk_free_bytes':shutil.disk_usage(R).free,'no_automatic_continuation':True})
print(json.dumps(run_outcome,indent=2),flush=True)
if run_outcome['failure']:raise SystemExit(1)
