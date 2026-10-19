#!/usr/bin/env python3
"""One root-authorized environment-corrected full Chapter10 API invocation; no retries."""
from pathlib import Path
import hashlib,json,os,signal,subprocess,threading,time,shutil
R=Path('/Users/bill/projects/ensemble')
S=Path('/Users/bill/projects/ensemble-edition-3-student')
O=R/'docs/edition-3-notes/ch10-entry-audit-evidence'
E=R/'docs/edition-3-notes/ch10-entry-followup-evidence'
M=json.loads((R/'docs/edition-3-notes/ch10-entry-manifest.json').read_text())
W=Path(json.loads((O/'workspace.json').read_text())['owned_scratch'])
receipts=[];observed={};processes=set();stop=threading.Event()
def ident(p):
 b=Path(p).read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def write(n,v): (E/n).write_text(json.dumps(v,indent=2)+'\n')
def git(root,*args): return subprocess.check_output(['git','--no-pager',*args],cwd=root)
def verify():
 inv=O/'inventory.json';expected={'bytes':16859,'sha256':'1000f80deb8d88f6efd9c00a2c1c4d4b242b46df36b00c8f2409d33acc3dd26e'};assert ident(inv)==expected
 seal=json.loads(inv.read_text());actual={p.relative_to(R).as_posix() for p in O.rglob('*') if p.is_file() and p!=inv}|{'docs/edition-3-notes/ch10-entry-audit.md'};assert actual=={x['path'] for x in seal['files']}
 for row in seal['files']: assert ident(R/row['path'])=={k:row[k] for k in ('bytes','sha256')},row['path']
 out={'original_inventory':expected,'original_sealed_files':seal['files'],'groups':{},'repositories':{}}
 for group,root,items in [('inputs',R,M['inputs']),('protected_tracked',R,M['protected_tracked']),('accepted_predecessor',S,M['accepted_predecessor']['files']),('historical_reference',R/'solutions/ch10',json.loads((O/'reference-source-pins.json').read_text()))]:
  for rel,want in items.items():assert ident(root/rel)==want,rel
  out['groups'][group]={'root':str(root),'count':len(items),'verified_files':items}
 for name,root in [('main',R),('student',S)]:
  d=git(root,'diff','--binary');c=git(root,'diff','--cached','--binary')
  out['repositories'][name]={'head':git(root,'rev-parse','HEAD').decode().strip(),'tree':git(root,'rev-parse','HEAD^{tree}').decode().strip(),'status':git(root,'status','--porcelain=v1','--untracked-files=normal').decode(),'unstaged_diff':{'bytes':len(d),'sha256':hashlib.sha256(d).hexdigest()},'staged_diff':{'bytes':len(c),'sha256':hashlib.sha256(c).hexdigest()}}
 assert out['repositories']['student']['status']==''
 assert out['repositories']['student']['head']==M['accepted_predecessor']['commit']
 assert out['repositories']['student']['tree']==M['accepted_predecessor']['tree']
 out['disk_free_bytes']=shutil.disk_usage(R).free;assert out['disk_free_bytes']>=1024**3
 return out
old=json.loads((O/'launch-environment.json').read_text());env=dict(old);env.update(PATH='/Users/bill/go/bin:'+old['PATH'],GOMAXPROCS='2',GOFLAGS='-mod=readonly -p=1')
assert {k for k in set(old)|set(env) if old.get(k)!=env.get(k)}=={'PATH','GOMAXPROCS','GOFLAGS'}
write('environment-v1-original.json',old);write('environment-v2-corrected.json',env)
write('environment-delta.json',{k:{'before':old.get(k),'after':env.get(k)} for k in ('PATH','GOMAXPROCS','GOFLAGS')})
original=next(x for x in json.loads((O/'command-receipts.json').read_text()) if x['name']=='07-full-chapter10')
argv=original['argv'];assert ident(argv[0])==original['executable_identity'];assert ident(argv[0])==json.loads((O/'explicit-binaries.json').read_text())[argv[0]]
assert ident(argv[2])==ident(O/'consumer-case.json')
dlv='/Users/bill/go/bin/dlv';assert ident(dlv)=={'bytes':20038402,'sha256':'10dbb1ce8e6d2cb77ae4d30c7ca1db0302b83529cd6cd39453f009a7653c17fe'}
write('retained-identities.json',{'full_api':{'path':argv[0],**ident(argv[0])},'consumer_case':{'path':argv[2],**ident(argv[2])},'helper_source':{'path':str(O/'full-api.go.txt'),**ident(O/'full-api.go.txt')},'dlv':{'path':dlv,**ident(dlv)},'followup_task':{'path':'/tmp/cadence-ch10-entry-followup-task.txt',**ident('/tmp/cadence-ch10-entry-followup-task.txt')}})
before=verify();write('verification-before.json',before)

def sample():
 while not stop.wait(.10):
  try:
   for p in (W/'tmp').glob('course-grade-*/submission'):
    if str(p) not in observed and p.is_file():observed[str(p)]={'observed_unix':time.time(),**ident(p)}
   output=subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True)
   processes.update(x.strip() for x in output.splitlines() if str(W) in x)
  except (OSError,subprocess.SubprocessError):pass

def run(name,args,timeout):
 assert shutil.disk_usage(R).free>=1024**3
 start=time.time();p=subprocess.Popen(args,cwd=original['cwd'],env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE,start_new_session=True);timed=False
 try:out,err=p.communicate(timeout=timeout)
 except subprocess.TimeoutExpired:
  timed=True;os.killpg(p.pid,signal.SIGTERM)
  try:out,err=p.communicate(timeout=5)
  except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);out,err=p.communicate()
 for suffix,data in [('stdout',out),('stderr',err)]:
  assert b'\0' not in data;data.decode('utf-8');(E/f'{name}.{suffix}.txt').write_bytes(data)
 row={'name':name,'argv':args,'cwd':original['cwd'],'environment':'environment-v2-corrected.json','pid':p.pid,'start_unix':start,'elapsed_seconds':round(time.time()-start,3),'timeout_seconds':timeout,'exit_code':p.returncode,'timed_out':timed,'stdout':f'{name}.stdout.txt','stderr':f'{name}.stderr.txt','executable_identity':ident(args[0]),'disk_free_bytes_after':shutil.disk_usage(R).free}
 receipts.append(row);write('command-receipts.json',receipts);print(name,'exit',p.returncode,'seconds',row['elapsed_seconds'],flush=True)
 return row,out,err
thread=threading.Thread(target=sample,daemon=True)
status={'full_api_attempts':0,'accepted_positive_control':False}
try:
 row,out,err=run('01-command-resolution',['/bin/sh','-c','command -v dlv; command -v go'],30)
 assert row['exit_code']==0 and not row['timed_out'] and out.decode().splitlines()==['/Users/bill/go/bin/dlv','/usr/local/go/bin/go'], 'resolution preflight failed'
 row,out,err=run('02-dlv-version',[dlv,'version'],30)
 assert row['exit_code']==0 and not row['timed_out'] and b'Delve Debugger' in out, 'Delve version preflight failed'
 thread.start();status['full_api_attempts']=1
 row,out,err=run('03-full-chapter10-environment-corrected',argv,600)
 parsed=json.loads(out.decode());write('parsed-full-report.json',parsed)
 ids=['save-shape','default-load','replay-equals-snapshot','tail-applied-once','log-not-needed','bad-save-refused','ch9-parity']
 status.update(helper_exit_zero=row['exit_code']==0,not_timed_out=not row['timed_out'],passed_is_true=parsed.get('passed') is True,score_is_100=parsed.get('score')==100,max_score_is_100=parsed.get('max_score')==100,title_is_chapter10=parsed.get('title')=='Edition 2, Chapter 10 — Persistence',check_ids_match=[c.get('id') for c in parsed.get('checks',[])]==ids,all_checks_pass=all(c.get('passed') is True for c in parsed.get('checks',[])))
 status['accepted_positive_control']=all(status[k] for k in ['helper_exit_zero','not_timed_out','passed_is_true','score_is_100','max_score_is_100','title_is_chapter10','check_ids_match','all_checks_pass'])
 if not status['accepted_positive_control']:status['failure']='Full report did not satisfy all positive-control requirements; no retry authorized.'
except Exception as ex:status['failure']=type(ex).__name__+': '+str(ex)
finally:
 stop.set()
 if thread.is_alive():thread.join()
 write('run-result.json',status)
 write('observed-transient-binaries.json',{'coverage':'Sampled at 100ms, non-exhaustive. Hash receipts only; no binaries exported.','binaries':observed})
 (E/'observed-processes.txt').write_text('\n'.join(sorted(processes))+'\n')
 after=verify();write('verification-after.json',after)
 for name in ['main','student']:
  for key in ['head','tree','unstaged_diff','staged_diff']:assert before['repositories'][name][key]==after['repositories'][name][key],(name,key)
 live=[x.strip() for x in subprocess.check_output(['/bin/ps','-axo','pid=,ppid=,command='],text=True).splitlines() if str(W) in x]
 write('completion.json',{'owned_processes_remaining':live,'all_commands_waited':True,'original_seal_unchanged':True,'protected_sources_unchanged':True,'predecessor_clean':True,'shared_cleanup':'none','owned_cleanup':'No new helper or binary built. Existing retained scratch preserved; grade.Build managed its own per-run temporary children.','sampled_transient_binaries':len(observed),'no_retry':True,'disk_floor_bytes':1024**3})
 assert not live,live
print(json.dumps(status,indent=2),flush=True)
if not status['accepted_positive_control']:raise SystemExit(1)
