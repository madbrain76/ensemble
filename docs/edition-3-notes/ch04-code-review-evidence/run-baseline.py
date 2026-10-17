#!/usr/bin/env python3
import datetime,hashlib,json,os,pathlib,re,signal,subprocess,tempfile,time
ROOT=pathlib.Path('/Users/bill/projects/ensemble');ST=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');E=ROOT/'docs/edition-3-notes/ch04-code-review-evidence';SE=ROOT/'docs/edition-3-notes/ch04-student-evidence'
def sha(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def git(*a):return subprocess.check_output(['git','--no-pager','-C',str(ST),*a],text=True).strip()
pin=json.loads((ROOT/'docs/edition-3-notes/ch04-candidate.json').read_text())
assert git('rev-parse','HEAD')==pin['candidate']['commit'];assert git('rev-parse','HEAD^{tree}')==pin['candidate']['tree'];assert not git('status','--porcelain')
for p,h in pin['pins'].items():assert sha(p)==h,p
for p,h in pin['candidate']['binaries'].items():assert sha(p)==h,p
src=json.loads((SE/'source-manifest.json').read_text())
for x in src['files']:assert sha(ST/x['path'])==x['sha256'],x
em=json.loads((SE/'evidence-manifest.json').read_text())
for x in em['files']:assert sha(SE/x['path'])==x['sha256'],x
S=pathlib.Path(tempfile.mkdtemp(prefix='e3-ch04-code-review-'));C=S/'candidate';C.mkdir();archive=S/'candidate.tar'
with archive.open('wb') as f:subprocess.run(['git','--no-pager','-C',str(ST),'archive',pin['candidate']['commit']],stdout=f,check=True)
subprocess.run(['tar','-xf',archive,'-C',C],check=True);(E/'scratch.txt').write_text(str(S)+'\n')
docs=['AGENTS.md','.agents/skills/ensemble-coding/SKILL.md','.agents/skills/edition-3-code-reviewer/SKILL.md','docs/agentic-codebook-workflow.md','docs/edition-2-carryover.md','book/voice.md','book/chapter-writing-procedure.md','docs/edition-3-execution.md','docs/edition-3-student.md','docs/edition-3-readiness.md','docs/edition-3-carryover.md','book/edition-3/chapter-04.md']
(E/'input-identities.json').write_text(json.dumps({'at_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'candidate':pin['candidate'],'pins':pin['pins'],'source_files_verified':len(src['files']),'evidence_files_verified':len(em['files']),'docs':{p:sha(ROOT/p) for p in docs},'adapter':{p:sha(ROOT/p) for p in ['cmd/grade-edition3/main.go','cmd/grade-edition3/main_test.go']},'go_version':subprocess.check_output(['go','version'],text=True).strip(),'archive_sha256':sha(archive)},indent=2)+'\n')
env={k:v for k,v in os.environ.items() if not re.search(r'API_KEY|TOKEN|SECRET|PASSWORD|CREDENTIAL|^(ANTHROPIC|OPENAI|GEMINI|GOOGLE|LLM|ENSEMBLE|AWS|AZURE)',k,re.I)};env.update(GOPROXY='off',GOTOOLCHAIN='local',GOFLAGS='-mod=readonly',GIT_PAGER='cat')
records=[]
def run(name,args,cwd,timeout=240):
 args=list(map(str,args));start=time.monotonic();p=subprocess.Popen(args,cwd=cwd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True);expired=False
 try:out=p.communicate(timeout=timeout)[0]
 except subprocess.TimeoutExpired:expired=True;os.killpg(p.pid,signal.SIGKILL);out=p.communicate()[0]
 (E/(name+'.log')).write_bytes(out);r={'name':name,'argv':args,'cwd':str(cwd),'exit_code':p.returncode,'timeout':expired,'seconds':round(time.monotonic()-start,3),'output_sha256':hashlib.sha256(out).hexdigest()};records.append(r);(E/'baseline-runs.json').write_text(json.dumps(records,indent=2)+'\n');print(name,p.returncode,r['seconds'],flush=True)
 if name.startswith('grade') and p.returncode in [0,1]:
  report=json.loads(out);r['passed']=report['passed'];r['score']=report['score'];r['max_score']=report['max_score'];(E/'baseline-runs.json').write_text(json.dumps(records,indent=2)+'\n')
 return r
run('full',['go','test','-count=1','./...'],C)
run('race',['go','test','-race','-count=1','./...'],C)
run('vet',['go','vet','./...'],C)
run('grade4',['go','run','./cmd/grade-edition3','-ch','4','-json','-core',ST,'-consumer',ST/'examples/ch04-consumer','-consumer-case',ST/'examples/ch04-consumer/consumer-case.json','/tmp/ch04-core'],ROOT)
run('grade1',['go','run','./cmd/grade-edition3','-ch','1','-json','/tmp/ch04-core'],ROOT)
for p,h in pin['candidate']['binaries'].items():assert sha(p)==h,p
assert not git('status','--porcelain');assert git('rev-parse','HEAD')==pin['candidate']['commit']
(E/'baseline-completed.json').write_text(json.dumps({'all_commands_exit0':all(r['exit_code']==0 for r in records),'binaries_unchanged':True,'candidate_clean_unchanged':True,'completed_utc':datetime.datetime.now(datetime.timezone.utc).isoformat()},indent=2)+'\n');print('COMPLETE',flush=True)
