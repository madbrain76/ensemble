#!/usr/bin/env python3
"""Receipt/preflight wrapper only. Does not replace any sealed function."""
import hashlib,json,os,pathlib,shutil,subprocess,sys,time
ROOT=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');HERE=pathlib.Path(__file__).resolve().parent;PREP=ROOT/'.ch09-local/operational-preparation';APPROVAL=pathlib.Path('/tmp/cadence-ch09-operational-execution-approval-v2.json')
POLICY='b8d8525398503564367003eaa3c34dc73dab68be971238c86377752e213a0e79';SEAL='070e03c3703a629fe3b29cd1adadc9292563a6bcf07ce1839aaefe910fb43ab1'
def sha(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def save(p,data):p.write_text(json.dumps(data,indent=2)+'\n')
def require(x,n):
 if not x:raise RuntimeError(n)
def cleanup(mode,run):
 d=run/mode;result=json.loads((d/'result.json').read_text());require(result['ok'],'sealed result failed');require(all(row['closed'] for row in result['cleanup']),'sealed cleanup failed')
 app=json.loads((d/'artifact-close.json').read_text());proxy=json.loads((d/'proxy-close.json').read_text());require(app['process_exited'] and app['listener_closed'],'artifact cleanup failed');require(proxy['listener_closed'] and proxy['active_upstreams']==0 and proxy['requests']<=8,'proxy cleanup failed')
 if mode=='browser-local':
  exits=list((d/'browser').glob('*browser-exit.json'));require(len(exits)==1,'browser exit receipt missing');br=json.loads(exits[0].read_text());require(br['listener_closed'] and br['exit_code'] is not None,'browser cleanup failed')
  if br['pid']:
   try:os.kill(br['pid'],0)
   except ProcessLookupError:pass
   else:raise RuntimeError('browser PID remains')
 return dict(result_ok=True,cleanup=result['cleanup'],artifact=app,proxy=proxy)
mode=sys.argv[1];require(mode in ('browser-local','paid-cli'),'scope');require(not (HERE/(mode+'.receipt.json')).exists(),'wrapper attempt already consumed');require(shutil.disk_usage(ROOT).free>=1024**3,'disk floor')
require(sha(APPROVAL)=='4207ec9d40b046a7ffea927b0c6e685b33094c31b48fb441bff9aeef9939e31b','approval pin');a=json.loads(APPROVAL.read_text());require(0<a['expires_unix']-time.time()<=3600 and mode in a['modes'],'approval expiry/scope');run=pathlib.Path(a['run_directory']);require(not (run/mode).exists(),'sealed one-shot directory already exists')
require(sha(PREP/'seal.json')==SEAL and sha(PREP/'policy.json')==POLICY,'seal/policy pins');seal=json.loads((PREP/'seal.json').read_text());policy=json.loads((PREP/'policy.json').read_text())
for row in seal['preparation_files']:require(sha(PREP/row['path'])==row['sha256'],'script/fixture pin')
require(sha(policy['source_manifest']['path'])==policy['source_manifest']['sha256'],'manifest pin');manifest=json.loads(pathlib.Path(policy['source_manifest']['path']).read_text())
for row in policy['source_files']:require(sha(ROOT/row['path'])==row['sha256'],'source file pin')
for row in manifest['binaries']:require(sha(row['path'])==row['sha256'] and pathlib.Path(row['path']).stat().st_size==row['bytes'],'five binary pins')
for name,digest in policy['assets'].items():require(sha(ROOT/'gui/web/gui'/name)==digest,'asset pin')
c=a['chrome_pin'];require(pathlib.Path(c['path']).stat().st_size==c['bytes'] and sha(c['path'])==c['sha256'],'Chrome pin')
require(subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=ROOT,text=True).strip()==policy['source_commit'],'HEAD');require(not subprocess.check_output(['git','--no-pager','status','--porcelain'],cwd=ROOT,text=True).strip(),'worktree')
if mode=='paid-cli':
 prior=json.loads((HERE/'browser-local.receipt.json').read_text());require(prior['exit_code']==0 and prior['gate_passed'],'browser prerequisite');cleanup('browser-local',run)
env={k:os.environ[k] for k in ('PATH','HOME','TMPDIR') if k in os.environ};env['PYTHONDONTWRITEBYTECODE']='1'
controller=['.ch09-local/operational-preparation/controller.py','execute','--exercise',mode,'--approval',str(APPROVAL),'--policy-sha256',POLICY,'--seal-sha256',SEAL]
command=['python3','-B',str(HERE/'audit-entry.py'),str(HERE/(mode+'.audit.json'))]+controller
receipt=dict(mode=mode,actual_argv=command,effective_controller_argv=['python3','-B']+controller,cwd=str(ROOT),environment=env,started_unix=time.time(),wrapper_sha256=sha(__file__),audit_wrapper_sha256=sha(HERE/'audit-entry.py'),approval_sha256=sha(APPROVAL),policy_sha256=POLICY,seal_sha256=SEAL,preflight_passed=True,preflight_source_files=len(policy['source_files']),preflight_binaries=len(manifest['binaries']),chrome_pin=c)
save(HERE/(mode+'.started.json'),receipt);tick=time.monotonic()
with (HERE/(mode+'.stdout.txt')).open('w') as out,(HERE/(mode+'.stderr.txt')).open('w') as err:process=subprocess.run(command,cwd=ROOT,env=env,stdout=out,stderr=err)
receipt.update(exit_code=process.returncode,elapsed_seconds=round(time.monotonic()-tick,3),finished_unix=time.time(),stdout_sha256=sha(HERE/(mode+'.stdout.txt')),stderr_sha256=sha(HERE/(mode+'.stderr.txt')),gate_passed=False)
if process.returncode==0:
 try:receipt['cleanup_verification']=cleanup(mode,run);receipt['gate_passed']=True
 except Exception as e:receipt['gate_error']=str(e)
save(HERE/(mode+'.receipt.json'),receipt);print(json.dumps(receipt,indent=2));print((HERE/(mode+'.stdout.txt')).read_text());print((HERE/(mode+'.stderr.txt')).read_text(),file=sys.stderr);sys.exit(0 if receipt['gate_passed'] else 1)
