#!/usr/bin/env python3
"""Finite independent frozen-adapter controls; writes only own evidence/scratch."""
import datetime,difflib,hashlib,json,os,pathlib,re,shutil,signal,subprocess,tempfile,time
ROOT=pathlib.Path('/Users/bill/projects/ensemble')
E=ROOT/'docs/edition-3-notes/ch04-grader-review-evidence'
BASE='bac07e644facb25ba3e5bffb2961674763fed4c4'
EXPECTED={'cmd/grade-edition3/main.go':'d08616eee707b4b91363f528abafa85f9b2268af964a0c6d903efbae95c291a4','cmd/grade-edition3/main_test.go':'f505aebe9a8ef9bea5495c6ee723eea0880b1ae3735e4b36adadff2cceb71a5f'}
def sha(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
for p,h in EXPECTED.items():assert sha(ROOT/p)==h,p
S=pathlib.Path(tempfile.mkdtemp(prefix='e3-ch04-review-'));M=S/'main';M.mkdir()
(E/'scratch.txt').write_text(str(S)+'\n')
archive=S/'source.tar'
with archive.open('wb') as f:subprocess.run(['git','archive',BASE],cwd=ROOT,stdout=f,check=True)
subprocess.run(['tar','-xf',str(archive),'-C',str(M)],check=True)
for p in EXPECTED:shutil.copyfile(ROOT/p,M/p)
# Confirm all historical grader and provider sources equal frozen HEAD export.
tracked=subprocess.check_output(['git','ls-tree','-r','--name-only',BASE,'internal/grade','internal/gradeedition2','internal/gradeedition3','internal/fakevendor','internal/fakeanthropic'],cwd=ROOT,text=True).splitlines()
for p in tracked:assert (ROOT/p).read_bytes()==(M/p).read_bytes(),p
(E/'input-identities.json').write_text(json.dumps({'at_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'base':BASE,'candidate':EXPECTED,'historical_sources':{p:sha(M/p) for p in tracked},'archive_sha256':sha(archive),'go_version':subprocess.check_output(['go','version'],text=True).strip()},indent=2)+'\n')
ENV={k:v for k,v in os.environ.items() if not re.search(r'API_KEY|TOKEN|SECRET|PASSWORD|CREDENTIAL|^(ANTHROPIC|OPENAI|GEMINI|GOOGLE|LLM|ENSEMBLE|AWS|AZURE)',k,re.I)}
ENV.update(GOPROXY='off',GOTOOLCHAIN='local',GIT_PAGER='cat',PAGER='cat')
records=[]
def run(name,args,cwd=M,timeout=90,expected=0,check_report=None):
 args=list(map(str,args));start=time.monotonic();timed_out=False
 proc=subprocess.Popen(args,cwd=cwd,env=ENV,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True)
 try:out=proc.communicate(timeout=timeout)[0]
 except subprocess.TimeoutExpired:
  timed_out=True;os.killpg(proc.pid,signal.SIGKILL);out=proc.communicate()[0]
 p=E/(name+'.log');p.write_bytes(out)
 rec={'name':name,'argv':args,'cwd':str(cwd),'exit_code':proc.returncode,'timeout':timed_out,'seconds':round(time.monotonic()-start,3),'output_sha256':sha(p)}
 records.append(rec);(E/'runs.json').write_text(json.dumps(records,indent=2)+'\n')
 print(name,proc.returncode,rec['seconds'],flush=True)
 assert not timed_out and proc.returncode==expected,(name,rec)
 if check_report is not None:
  d=json.loads(out);assert (d['score'],d['max_score'],d['passed'])==check_report,(name,d)
 return out
run('independent-race',['go','test','-race','-count=1','-v','./cmd/grade-edition3'])
run('independent-vet',['go','vet','./cmd/grade-edition3'])
B=S/'grade-edition3';run('build-adapter',['go','build','-o',B,'./cmd/grade-edition3'])
R=S/'reference-lifetime-corrected';shutil.copytree(M/'solutions/ch04',R)
f=R/'internal/engine/engine.go';before=f.read_bytes();after=before
for old,new in [(b'\n\t"time"',b''),(b'&http.Client{Timeout: 60 * time.Second}',b'&http.Client{}')]:
 assert after.count(old)==1,old;after=after.replace(old,new,1)
f.write_bytes(after)
(E/'reference-timeout-only.patch').write_text(''.join(difflib.unified_diff(before.decode().splitlines(True),after.decode().splitlines(True),fromfile='original/internal/engine/engine.go',tofile='lifetime-corrected/internal/engine/engine.go')))
changed=[str(p.relative_to(R)) for p in R.rglob('*') if p.is_file() and p.read_bytes()!=(M/'solutions/ch04'/p.relative_to(R)).read_bytes()]
assert changed==['internal/engine/engine.go'],changed
(E/'reference-change.json').write_text(json.dumps({'changed_files':changed,'original_sha256':hashlib.sha256(before).hexdigest(),'corrected_sha256':sha(f),'purpose':'remove only previously measured whole-request timeout; no grader/probe delay change'},indent=2)+'\n')
T=S/'reference-cli';run('build-positive-reference',['go','build','-o',T,'.'],cwd=R)
case=S/'explicit-case.json';case.write_text(json.dumps({'name':'shout','arguments':{'text':'independent route'},'result':'INDEPENDENT ROUTE','diagnostic':'shout called'})+'\n');shutil.copyfile(case,E/'explicit-case.json')
common=['-ch','4','-json','-core',R,'-consumer',R/'ch05']
run('adapter-ch4-executable-default',[B,*common,T],check_report=(160,160,True))
run('adapter-ch4-directory-explicit-case',[B,*common,'-consumer-case',case,R],check_report=(160,160,True))
wrong=run('adapter-ch4-wrong-core',[B,'-ch','4','-json','-core',M/'solutions/ch04','-consumer',R/'ch05',T],expected=1,check_report=(135,160,False))
assert [c['id'] for c in json.loads(wrong)['checks'] if not c['passed']]==['star-topology']
# Real compiled pre-dispatch validation with an observable launch sentinel.
marker=S/'sentinel-invoked';sentinel=S/'sentinel';sentinel.write_text('#!/bin/sh\nprintf invoked >> '+str(marker)+'\nexit 1\n');sentinel.chmod(0o700)
bad=S/'malformed-case.json';bad.write_text('not JSON\n')
incomplete=S/'incomplete-case.json';incomplete.write_text('{}\n')
for name,extra in [('malformed',[*common,'-consumer-case',bad]),('missing',[*common,'-consumer-case',S/'missing.json']),('incomplete',[*common,'-consumer-case',incomplete]),('ch1-empty-core',['-ch','1','-core=']),('ch1-empty-consumer',['-ch','1','-consumer=']),('ch1-empty-case',['-ch','1','-consumer-case='])]:
 run('compiled-invalid-'+name,[B,*extra,sentinel],expected=2)
 assert not marker.exists(),name
# Private full-route sensitivity: a wrong forwarding argument must fail.
D=M/'cmd/ch04-review-mutant';D.mkdir();src=(M/'cmd/grade-edition3/main.go').read_bytes();old=b'gradeedition2.Chapter4(bin, *core, *consumer, example)';new=b'gradeedition2.Chapter4(bin, *consumer, *consumer, example)';assert src.count(old)==1
(D/'main.go').write_bytes(src.replace(old,new,1))
(E/'adapter-mutation.patch').write_text(''.join(difflib.unified_diff(src.decode().splitlines(True),(D/'main.go').read_text().splitlines(True),fromfile='candidate/main.go',tofile='mutant/main.go')))
MB=S/'misrouted-grader';run('build-misrouted-adapter',['go','build','-o',MB,'./cmd/ch04-review-mutant'])
mut=run('misrouted-adapter-control',[MB,*common,T],expected=1,check_report=(135,160,False))
assert [c['id'] for c in json.loads(mut)['checks'] if not c['passed']]==['star-topology']
# Same immutable binary for both acceptance reports. Real 62s header AND body phases.
paired_hash=sha(T)
run('adapter-ch1-same-binary-lifetime',[B,'-ch','1','-json',T],timeout=210,check_report=(100,100,True))
assert sha(T)==paired_hash
life=json.loads((E/'adapter-ch1-same-binary-lifetime.log').read_text())
assert all(c['passed'] for c in life['checks'] if c['id'] in ['generation_headers','generation_body'])
assert len([c for c in life['checks'] if c['id'] in ['generation_headers','generation_body']])==2
for p,h in EXPECTED.items():assert sha(ROOT/p)==h,p
(E/'binary-identities.json').write_text(json.dumps({str(p):sha(p) for p in [B,T,MB]},indent=2)+'\n')
(E/'completed.json').write_text(json.dumps({'at_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'candidate_unchanged':True,'paired_reference_sha256':paired_hash,'all_expected_outcomes':True,'credential_or_live_calls':False},indent=2)+'\n')
print('COMPLETE',flush=True)
