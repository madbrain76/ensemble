from pathlib import Path
import os,sys,json,hashlib,subprocess,time,datetime,shutil,errno,re
ROOT=Path('/Users/bill/projects/ensemble-edition-3-student')
NOTES=Path('/Users/bill/projects/ensemble/docs/edition-3-notes')
OLD=NOTES/'ch10-student-evidence'
NEW=NOTES/'ch10-student-final-reconciliation-evidence'
WORK=ROOT/'.ch10-local/final-reconciliation'
CONFIG=json.loads((WORK/'config.json').read_text())
ENV=CONFIG['environment']
MANIFEST=json.loads((NOTES/'ch10-student-source-manifest.json').read_text())
OLDSEAL=json.loads((OLD/'seal.json').read_text())
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def pin(p):return {'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)}
def git(*args):return subprocess.check_output(['git','--no-pager',*args],cwd=ROOT,env=ENV,text=True).strip()
def identity():
 head=git('rev-parse','HEAD');tree=git('rev-parse','HEAD^{tree}');status=git('status','--porcelain','--untracked-files=normal')
 assert head==CONFIG['source_head'] and tree==CONFIG['source_tree'] and not status
 for f in MANIFEST['tracked_files']:
  p=ROOT/f['path'];assert p.stat().st_size==f['bytes'] and sha(p)==f['sha256'],str(p)
 for f in OLDSEAL['reports']+OLDSEAL['evidence']+CONFIG['prior_reports']+CONFIG['artifacts']:
  p=Path(f['path']);assert p.stat().st_size==f['bytes'] and sha(p)==f['sha256'],str(p)
 return {'head':head,'tree':tree,'parent':git('rev-parse','HEAD^'),'status':status,'tracked_files_verified':len(MANIFEST['tracked_files']),'prior_reports':CONFIG['prior_reports'],'artifacts':CONFIG['artifacts']}
def group_exists(pid):
 try:os.killpg(pid,0);return True
 except ProcessLookupError:return False

def run(name,argv,cwd,chapter):
 before=identity();free=shutil.disk_usage(ROOT).free
 if free<1024**3:
  (NEW/(name+'.guard.json')).write_text(json.dumps({'launched':False,'free_bytes':free,'floor_bytes':1024**3,'source':before},indent=2)+'\n');return False
 out=NEW/(name+'.stdout.txt');err=NEW/(name+'.stderr.txt');receipt=NEW/(name+'.receipt.json')
 assert not out.exists() and not err.exists() and not receipt.exists()
 started=datetime.datetime.now(datetime.timezone.utc).isoformat();tick=time.monotonic()
 with out.open('wb') as stdout,err.open('wb') as stderr:
  process=subprocess.Popen(argv,cwd=cwd,env=ENV,stdout=stdout,stderr=stderr,start_new_session=True)
  code=process.wait() # No short total deadline: both actual 62-second lifetime phases may finish.
 elapsed=time.monotonic()-tick
 for _ in range(20):
  if not group_exists(process.pid):break
  time.sleep(0.1)
 clean=not group_exists(process.pid)
 after=identity();assert before==after
 parse_error=None
 try:result=json.loads(out.read_text())
 except Exception as e:result={};parse_error=str(e)
 match=re.search(r'\bChapter\s+(\d+)\b',result.get('title',''))
 reported_chapter=int(match[1]) if match else result.get('chapter')
 ok=code==0 and result.get('passed') is True and reported_chapter==chapter and clean
 r={'name':name,'argv':argv,'cwd':cwd,'env':ENV,'started_utc':started,'ended_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'duration_seconds':elapsed,'exit_code':code,'requested_chapter':chapter,'reported_chapter':reported_chapter,'reported_title':result.get('title'),'reported_passed':result.get('passed'),'success':ok,'json_parse_error':parse_error,'source_and_artifacts_before':before,'source_and_artifacts_after':after,'outputs':[pin(out),pin(err)],'free_bytes_before':free,'free_bytes_after':shutil.disk_usage(ROOT).free,'grader_pid_and_owned_process_group':process.pid,'owned_process_group_gone_after_wait':clean,'outer_timeout_seconds':None,'lifetime_phases_shortened':False,'browser_or_paid_provider_authority':False}
 receipt.write_text(json.dumps(r,indent=2)+'\n')
 print(json.dumps({'name':name,'exit_code':code,'chapter':reported_chapter,'passed':result.get('passed'),'seconds':elapsed,'owned_group_gone':clean,'receipt':str(receipt)},indent=2),flush=True)
 print(out.read_text(),flush=True)
 return ok
r=json.loads((NEW/'chapter10-identity-reconciliation.json').read_text());assert r['grader_exit_code']==0 and r['grader_passed'] and r['identity_from_title']==10
argv=['/tmp/cadence-ch10-grader-prompt-fixed','-ch','1','-json',str(ROOT/'.ch10-local/bin/frozen-cli')]
if not run('chapter1-first-same-cli',argv,'/Users/bill/projects/ensemble',1):sys.exit(1)
print('First same-CLI Chapter1 completed; no Chapter10 rerun.',flush=True)
