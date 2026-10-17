import datetime, hashlib, json, os, pathlib, signal, subprocess, tarfile, tempfile, time
R=pathlib.Path('/Users/bill/projects/ensemble'); S=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student'); E=R/'docs/edition-3-notes/ch04-code-review-r1-evidence'; SE=R/'docs/edition-3-notes/ch04-student-r1-evidence'; OLD=R/'docs/edition-3-notes/ch04-code-review-evidence'
def sha(p): return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def git(p,*args): return subprocess.check_output(['git','--no-pager','-C',str(p),*args],text=True).strip()
def env():
 e={k:v for k,v in os.environ.items() if not any(s in k.upper() for s in ['KEY','TOKEN','SECRET','CREDENTIAL','ANTHROPIC','OPENAI','GEMINI','GOOGLE','ENSEMBLE','LLM_'])};e.update(GOPROXY='off',GOTOOLCHAIN='local',GOFLAGS='-mod=readonly',GIT_PAGER='cat');return e
def run(name,argv,cwd,timeout,records,receipt):
 start=time.monotonic();timed=False
 p=subprocess.Popen([str(x) for x in argv],cwd=cwd,env=env(),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True)
 try:out,_=p.communicate(timeout=timeout)
 except subprocess.TimeoutExpired:
  timed=True;os.killpg(p.pid,signal.SIGKILL);out,_=p.communicate()
 (E/(name+'.log')).write_bytes(out);row={'name':name,'argv':[str(x) for x in argv],'cwd':str(cwd),'exit_code':p.returncode,'timeout':timed,'seconds':round(time.monotonic()-start,3),'output_sha256':sha(E/(name+'.log'))};records.append(row);(E/receipt).write_text(json.dumps(records,indent=2)+'\n');print(name,p.returncode,row['seconds'],flush=True);return row,out
if __name__=='__main__':
 pin=json.loads((R/'docs/edition-3-notes/ch04-candidate-r1.json').read_text());assert git(S,'rev-parse','HEAD')==pin['commit'];assert git(S,'rev-parse','HEAD^{tree}')==pin['tree'];assert not git(S,'status','--porcelain')
 for p,h in pin['pins'].items():assert sha(p)==h,p
 original=json.loads((OLD/'input-identities.json').read_text())
 for p,h in {**original['docs'],**original['adapter']}.items():assert sha(R/p)==h,p
 sources=json.loads((SE/'source-manifest.json').read_text())['files'];current=json.loads((SE/'evidence-manifest.json').read_text())['files'];local=json.loads((SE/'local-artifacts.json').read_text())['files'];initial=json.loads((R/'docs/edition-3-notes/ch04-student-evidence/evidence-manifest.json').read_text())['files']
 for group,base in [(sources,S),(current,SE),(local,pathlib.Path('/')),(initial,R/'docs/edition-3-notes/ch04-student-evidence')]:
  for x in group:assert sha(base/x['path'])==x['sha256'],str(base/x['path'])
 T=pathlib.Path(tempfile.mkdtemp(prefix='e3-ch04-code-review-r1-'));C=T/'candidate';C.mkdir();M=T/'sensitivity';M.mkdir();(E/'scratch.txt').write_text(str(T)+'\n')
 with tarfile.open('/tmp/ch04-r1-candidate-source.tar.gz','r:gz') as t:
  members=t.getmembers();assert all(not pathlib.PurePosixPath(x.name).is_absolute() and '..' not in pathlib.PurePosixPath(x.name).parts and not x.issym() and not x.islnk() for x in members)
  t.extractall(C,filter='data');t.extractall(M,filter='data')
 for x in sources:
  data=(C/x['path']).read_bytes();assert sha(C/x['path'])==x['sha256'];assert data==subprocess.check_output(['git','--no-pager','-C',str(S),'show',pin['commit']+':'+x['path']])
 diff=subprocess.check_output(['git','--no-pager','-C',str(S),'diff',pin['commit']+'^',pin['commit'],'--','internal/agent/dispatch.go','chapter4_limits_test.go']);(E/'candidate.patch').write_bytes(diff)
 metadata=''
 for p in ['/tmp/ch04-r1-core','/tmp/ch04-r1-consumer']:
  text=subprocess.check_output(['go','version','-m',p],text=True,env=env());assert 'vcs.revision='+pin['commit'] in text and 'vcs.modified=false' in text;metadata+=text
 (E/'binary-build-metadata.txt').write_text(metadata)
 inputs={'at_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'pin':pin,'main_head':git(R,'rev-parse','HEAD'),'student_parent':git(S,'rev-parse','HEAD^'),'prerequisites':original['docs'],'adapter':original['adapter'],'execution_chronology_sha256':sha(R/'docs/edition-3-notes/ch04-execution.md'),'source_records':len(sources),'r1_evidence_records':len(current),'local_records':len(local),'initial_unchanged_records':len(initial),'archive_exact_committed_bytes':True,'source_files':sources,'initial_review_sha256':sha(R/'docs/edition-3-notes/ch04-code-review.md')};(E/'input-identities.json').write_text(json.dumps(inputs,indent=2)+'\n')
 records=[]
 for name,args,cwd,bound in [('full',['go','test','-count=1','./...'],C,300),('race',['go','test','-race','-count=1','./...'],C,300),('vet',['go','vet','./...'],C,90),('grade4',['go','run','./cmd/grade-edition3','-ch','4','-json','-core',S,'-consumer',S/'examples/ch04-consumer','-consumer-case',S/'examples/ch04-consumer/consumer-case.json','/tmp/ch04-r1-core'],R,150),('grade1',['go','run','./cmd/grade-edition3','-ch','1','-json','/tmp/ch04-r1-core'],R,230)]:
  row,out=run(name,args,cwd,bound,records,'baseline-runs.json');assert row['exit_code']==0 and not row['timeout'],name
  if name.startswith('grade'):
   j=json.loads(out);assert j['passed'] is True;row['report']={k:j[k] for k in ['passed','score','max_score']};(E/'baseline-runs.json').write_text(json.dumps(records,indent=2)+'\n')
 for p,h in pin['pins'].items():assert sha(p)==h,p
 for x in sources:assert sha(S/x['path'])==x['sha256'] and sha(C/x['path'])==x['sha256']
 assert not git(S,'status','--porcelain');(E/'baseline-complete.json').write_text(json.dumps({'completed_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'all_gates_completed':True,'all_pins_and_source_unchanged':True},indent=2)+'\n')
