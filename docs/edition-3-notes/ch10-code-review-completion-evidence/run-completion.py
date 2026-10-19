#!/usr/bin/env python3
import ast,difflib,hashlib,json,os,shutil,signal,stat,subprocess,tarfile,tempfile,threading,time,traceback
from pathlib import Path,PurePosixPath
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';E=N/'ch10-code-review-completion-evidence';OLD=N/'ch10-code-review-evidence';M=json.loads((N/'ch10-code-review-completion-input.json').read_text());C=json.loads((N/'ch10-student-source-manifest.json').read_text())
# Import only the previous verified read-only pin functions, never a control body.
tree=ast.parse((OLD/'v4/continue-probes.py').read_text());defs=[n for n in tree.body if isinstance(n,ast.FunctionDef) and n.name in ['digest','ident','git','verify']];exec(compile(ast.Module(body=defs,type_ignores=[]),'carried-read-only-verifier','exec'))
W=None;COPY=None;env=None;before=None;original={};modes={};rows=[];observed={};processes=set();done=threading.Event();sampler=None;failure=None;stage='preflight';mutation_applied=False;mutation_restored=False

def write(name,obj):
 p=E/name;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(obj if isinstance(obj,str) else json.dumps(obj,indent=2)+'\n')
def floor():
 free=shutil.disk_usage(R).free
 if free<1024**3:raise RuntimeError('below 1 GiB disk floor')
 return free

def verify_original_seal():
 p=OLD/'inventory.json';assert ident(p)=={'bytes':24503,'sha256':'d781ef13cd6cdfafe8dcc32dcc50ffde64dda611d0faa7472bda0a53acc6e932'};j=json.loads(p.read_text());assert len(j['files'])==116 and sum(x['bytes'] for x in j['files'])==1138356
 for x in j['files']:assert ident(R/x['path'])=={k:x[k] for k in ['bytes','sha256']}
 assert {x['path'] for x in j['files']}=={str(p.relative_to(R)) for p in [N/'ch10-code-review.md']+list(OLD.rglob('*')) if p.is_file() and p!=OLD/'inventory.json'}
 return {'inventory':ident(OLD/'inventory.json'),'members':116,'payload_bytes':1138356}

def restored():
 bad=[]
 for rel,b in original.items():
  p=COPY/rel
  if not p.is_file() or p.read_bytes()!=b or stat.S_IMODE(p.stat().st_mode)!=modes[rel]:bad.append(rel)
 return {'original_files':len(original),'mismatches':bad,'all_original_bytes_and_archive_modes_match':not bad}

def restore():
 for rel,b in original.items():p=COPY/rel;p.write_bytes(b);p.chmod(modes[rel])
 assert restored()['all_original_bytes_and_archive_modes_match']

def sample():
 while not done.wait(.15):
  try:
   for line in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,pgid=,command='],text=True).splitlines():
    if str(W) in line:processes.add(line.strip())
   for p in (W/'tmp').rglob('*'):
    if p.is_file() and p.name in ['ensemble.test','gui.test','ensemble-gui.test','submission']:
     key=str(p)
     if key not in observed:
      try:observed[key]={'path':key,**ident(p),'phase':stage,'non_exhaustive':True}
      except FileNotFoundError:pass
  except Exception:pass

def run(label,args,cwd,deadline,expected):
 global stage
 stage=label;free=floor();stdout=E/(label+'.stdout.txt');stderr=E/(label+'.stderr.txt');start=time.time();print('START',label,args,flush=True)
 with stdout.open('wb') as o,stderr.open('wb') as e:
  p=subprocess.Popen(args,cwd=cwd,env=env,stdin=subprocess.DEVNULL,stdout=o,stderr=e,start_new_session=True);timed=False
  try:code=p.wait(timeout=deadline)
  except subprocess.TimeoutExpired:
   timed=True;os.killpg(p.pid,signal.SIGTERM)
   try:code=p.wait(timeout=5)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);code=p.wait()
 row={'label':label,'argv':args,'cwd':str(cwd),'environment':env,'pid':p.pid,'process_group':p.pid,'start_unix':start,'duration_seconds':round(time.time()-start,6),'deadline_seconds':deadline,'exit_code':code,'timed_out':timed,'stdout':{'path':stdout.name,**ident(stdout)},'stderr':{'path':stderr.name,**ident(stderr)},'disk_free_before':free,'disk_free_after':floor()};rows.append(row);write('command-receipts.json',rows);print('END',label,code,row['duration_seconds'],flush=True)
 if timed or code!=expected:raise RuntimeError('unexpected control result; stop without repair: '+label)
 return stdout.read_text(),stderr.read_text()

try:
 floor();before=verify();write('verification-before.json',before);write('original-seal-before.json',verify_original_seal())
 assert ident('/tmp/cadence-ch10-code-review-completion-task.txt')=={'bytes':6971,'sha256':'fb8cfe55bd957f304d998179d35ed38645eb901dbb6fd3ccc434e825f4661e0e'}
 carry=[]
 for version in ['v2','v3','v4']:
  for row in json.loads((OLD/version/'command-receipts.json').read_text()):
   for k in ['stdout','stderr']:assert ident(OLD/version/row[k]['path'])=={a:row[k][a] for a in ['bytes','sha256']}
   carry.append({'version':version,'receipt':row})
 assert len(carry)==13;fresh=carry[0]['receipt'];assert fresh['label']=='01-fresh-ch1' and fresh['exit_code']==0 and fresh['duration_seconds']>=124
 r=json.loads((OLD/'v2'/fresh['stdout']['path']).read_text());assert r['passed'] and r['score']==r['max_score']==100 and len(r['checks'])==9 and all(c['passed'] for c in r['checks']);assert all(x['receipt']['exit_code']==0 for x in carry[1:11]);assert carry[11]['receipt']['exit_code']==carry[12]['receipt']['exit_code']==1
 ch10rows=json.loads((N/'ch10-prompt-fixture-review-evidence/v2/command-receipts.json').read_text());cr=json.loads((N/'ch10-prompt-fixture-review-evidence/v2'/ch10rows[0]['stdout']['path']).read_text());assert cr['passed'] and cr['score']==100 and all(c['passed'] for c in cr['checks']);assert ch10rows[0]['argv'][-1]==fresh['argv'][-1];assert ident(fresh['argv'][-1])=={k:C['binaries'][0][k] for k in ['bytes','sha256']}
 write('carried-controls-and-failures.json',{'original_commands':carry,'chapter10_receipt':ch10rows[0],'same_cli_identity':ident(fresh['argv'][-1]),'no_gate_or_suite_reexecution':True})
 W=Path(tempfile.mkdtemp(prefix='cadence-ch10-code-review-completion-'));COPY=W/'source';COPY.mkdir();
 for d in ['home','tmp','bin']:(W/d).mkdir()
 write('workspace.json',{'workspace':str(W),'source_copy':str(COPY),'archive_path':C['archive']['local_path'],'archive':{k:C['archive'][k] for k in ['bytes','sha256']},'extraction':'manual regular-file extraction from verified TAR; no checkout copying; no links/special files admitted'})
 assert ident(C['archive']['local_path'])=={k:C['archive'][k] for k in ['bytes','sha256']}
 with tarfile.open(C['archive']['local_path']) as t:
  members=t.getmembers();regular=[m for m in members if m.isfile()];assert len(regular)==137;assert {m.name for m in regular}=={x['path'] for x in C['tracked_files']}
  for m in members:
   path=PurePosixPath(m.name);assert not path.is_absolute() and '..' not in path.parts and (m.isdir() or m.isfile())
  expected={x['path']:x for x in C['tracked_files']}
  for m in regular:
   b=t.extractfile(m).read();x=expected[m.name];assert {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}=={k:x[k] for k in ['bytes','sha256']};p=COPY/m.name;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes(b);p.chmod(m.mode&0o777);original[m.name]=b;modes[m.name]=m.mode&0o777
 write('archive-extraction-pins.json',{'archive':C['archive'],'original_files':[{'path':r,'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest(),'archive_mode':modes[r]} for r,b in sorted(original.items())]});assert restored()['all_original_bytes_and_archive_modes_match']
 env={'PATH':'/Users/bill/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(W/'home'),'TMPDIR':str(W/'tmp'),'GOPATH':'/Users/bill/go','GOMODCACHE':'/Users/bill/go/pkg/mod','GOCACHE':'/Users/bill/Library/Caches/go-build','GOMAXPROCS':'2','GOFLAGS':'-mod=readonly -p=1','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','CGO_ENABLED':'1','LANG':'en_US.UTF-8','LC_ALL':'en_US.UTF-8'}
 assert 'LLM_BASE_URL' not in env and not any('KEY' in k or 'TOKEN' in k for k in env);write('environment.json',env);go=shutil.which('go',path=env['PATH']);assert go;write('tool-identities.json',{'go':{'path':go,**ident(go)},'uid':os.geteuid()});assert os.geteuid()!=0
 probes=[('independent_ch10_review_test.go.txt','reviewer_ch10_test.go'),('reviewer_legacy_ws_test.go.txt','gui/reviewer_legacy_ws_test.go')]
 for source,dest in probes:assert dest not in original;(COPY/dest).write_bytes((E/'probes'/source).read_bytes())
 write('control-plan.json',{'required_public_source':'book/edition-3/chapter-10.md','probes':[{'source':'probes/'+a,'copy_destination':b,**ident(E/'probes'/a)} for a,b in probes],'commands':['one public persistence/renderer/primary/atomicity suite','one real GUI collision/reconnect/new-generation/restart','one early-truncation mutant test','one exact-restoration sentinel test'],'mutation_target':'internal/agent/save.go','mutation':'one early os.WriteFile(path,nil,0600) before os.CreateTemp; compile-valid; only owned copy','stop_rule':'any setup/compile/constructor/real behavioral error stops; no repair/retry','all_prior_grades_and_ten_suites_carried':True,'disk_free':floor()})
 sampler=threading.Thread(target=sample,daemon=True);sampler.start()
 out,err=run('01-public-controls',[go,'test','-count=1','-run','^TestReview','-v','-timeout=90s','.'],COPY,150,0)
 for name in ['TestReviewReplayModesAndRichRender','TestReviewMalformedLoadAtomicityAndFrozenPrompt','TestReviewAtomicSavePreservesFile']:assert '=== RUN   '+name in out and '--- PASS: '+name in out
 assert '[no tests to run]' not in out and 'max tokens must be positive' not in out and 'valid fixture constructor:' not in out
 out,err=run('02-gui-controls',[go,'test','-count=1','-run','^TestReviewLegacyCollisionOverWebSocket$','-v','-timeout=60s','.'],COPY/'gui',100,0);assert '--- PASS: TestReviewLegacyCollisionOverWebSocket' in out and 'fake requests=1' in out
 assert restored()['all_original_bytes_and_archive_modes_match'];target=COPY/'internal/agent/save.go';old=target.read_text();anchor='\tf, err := os.CreateTemp(filepath.Dir(path), ".save-*")';assert old.count(anchor)==1;new=old.replace(anchor,'\t_ = os.WriteFile(path, nil, 0600) // REVIEW MUTANT: forbidden early truncation\n'+anchor);target.write_text(new);mutation_applied=True
 write('mutation.diff.txt',''.join(difflib.unified_diff(old.splitlines(True),new.splitlines(True),fromfile='archive-original/internal/agent/save.go',tofile='owned-mutant/internal/agent/save.go')));write('mutation-pins.json',{'original':{'bytes':len(old.encode()),'sha256':hashlib.sha256(old.encode()).hexdigest()},'mutant':ident(target),'only_changed_original_file':restored()['mismatches']});assert restored()['mismatches']==['internal/agent/save.go']
 out,err=run('03-atomic-mutant',[go,'test','-count=1','-run','^TestReviewAtomicSavePreservesFile$','-v','-timeout=45s','.'],COPY,100,1);assert '=== RUN   TestReviewAtomicSavePreservesFile' in out and 'ATOMIC_SENTINEL_CORRUPTED' in out and '--- FAIL: TestReviewAtomicSavePreservesFile' in out and '[build failed]' not in out+err
 restore();mutation_restored=True;write('restoration-before-control.json',restored());out,err=run('04-atomic-restored',[go,'test','-count=1','-run','^TestReviewAtomicSavePreservesFile$','-v','-timeout=45s','.'],COPY,100,0);assert '--- PASS: TestReviewAtomicSavePreservesFile' in out and 'sentinel unchanged' in out
 write('outcome.json',{'passed':True,'missing_control_groups_completed':True,'runtime_commands':len(rows),'prior_gate_replays':0,'prior_suite_replays':0,'sensitivity':'baseline green -> early-truncation mutant specific sentinel failure -> exact restored green'})
except Exception as ex:
 failure={'stage':stage,'type':type(ex).__name__,'error':str(ex),'traceback':traceback.format_exc(),'stop_without_repair':True};write('outcome.json',{'passed':False,'failure':failure,'runtime_commands':len(rows),'remaining_controls_not_run':True});print('STOP',failure['traceback'],flush=True)
finally:
 done.set()
 if sampler is not None:sampler.join(15)
 if COPY is not None and original:
  restore();mutation_restored=mutation_applied or mutation_restored;write('copied-source-final.json',restored())
 write('observed-executables.json',list(observed.values()));write('observed-processes.txt','\n'.join(sorted(processes))+'\n')
 after=None
 try:
  after=verify();write('verification-after.json',after);assert before is None or before==after;write('original-seal-after.json',verify_original_seal())
 except Exception as ex:write('final-verification-error.json',{'type':type(ex).__name__,'error':str(ex)});failure=failure or {'stage':'final verification','error':str(ex)}
 live=[]
 if W is not None:
  groups={x['process_group'] for x in rows}
  for line in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,pgid=,command='],text=True).splitlines():
   q=line.strip().split(None,3)
   if len(q)==4 and (str(W) in q[3] or int(q[2]) in groups):live.append(line.strip())
 write('completion.json',{'controls_completed':len(rows),'all_top_level_jobs_waited':True,'sampler_joined':sampler is None or not sampler.is_alive(),'owned_processes_remaining':live,'before_after_equal':before is not None and before==after,'mutation_applied':mutation_applied,'mutation_restored':mutation_restored,'source_restored':restored() if original else None,'failure':failure,'owned_tmp_entries':sorted(p.name for p in (W/'tmp').iterdir()) if W else [],'disk_free':shutil.disk_usage(R).free});print('COMPLETION',json.dumps({'commands':len(rows),'failure':failure is not None,'owned_processes':len(live),'mutation_applied':mutation_applied,'mutation_restored':mutation_restored}),flush=True)
