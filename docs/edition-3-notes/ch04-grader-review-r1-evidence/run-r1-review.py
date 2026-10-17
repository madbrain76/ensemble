#!/usr/bin/env python3
import datetime,difflib,hashlib,json,os,pathlib,re,shutil,signal,subprocess,tempfile,time
ROOT=pathlib.Path('/Users/bill/projects/ensemble');E=ROOT/'docs/edition-3-notes/ch04-grader-review-r1-evidence';PRE=ROOT/'docs/edition-3-notes/ch04-grader-review-evidence'
expected={'cmd/grade-edition3/main.go':'539ee3e6cc7313171ab57d268803ddd4c2d49be33fb028c38fa43b3b31c12834','cmd/grade-edition3/main_test.go':'c94e75cc895bfcf5e2625b523ff5571955c565ad8ef9f261b1f7ae3f7637fa2f'}
def sha(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
prior=json.loads((PRE/'input-identities.json').read_text());old=pathlib.Path((PRE/'scratch.txt').read_text().strip())
assert sha(old/'source.tar')==prior['archive_sha256']
branch=b'\t\tif err == flag.ErrHelp {\n\t\t\treturn 0\n\t\t}\n'
for name,h in expected.items():
 assert sha(ROOT/name)==h
 before=(old/'main'/name).read_bytes();after=(ROOT/name).read_bytes();assert hashlib.sha256(before).hexdigest()==prior['candidate'][name]
 if name.endswith('/main.go'):
  assert after.count(branch)==1;assert after.replace(branch,b'',1)==before
 else:
  anchor=b'\nfunc TestHelpSucceedsWithoutLaunchingTarget('
  assert after.count(anchor)==1;assert after[:after.index(anchor)]==before
 (E/(pathlib.Path(name).name+'.patch')).write_text(''.join(difflib.unified_diff(before.decode().splitlines(True),after.decode().splitlines(True),fromfile='previous/'+name,tofile='r1/'+name)))
for name,h in prior['historical_sources'].items():assert sha(ROOT/name)==h,name
S=pathlib.Path(tempfile.mkdtemp(prefix='e3-ch04-r1-review-'));M=S/'main';M.mkdir();(E/'scratch.txt').write_text(str(S)+'\n')
subprocess.run(['tar','-xf',str(old/'source.tar'),'-C',str(M)],check=True)
for name in expected:shutil.copyfile(ROOT/name,M/name)
(E/'input-identities.json').write_text(json.dumps({'at_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'base':prior['base'],'candidate':expected,'prior_candidate':prior['candidate'],'reverse_delta_exact':True,'historical_source_hashes_verified':len(prior['historical_sources']),'prior_manifest_sha256':sha(PRE/'input-identities.json'),'go_version':subprocess.check_output(['go','version'],text=True).strip()},indent=2)+'\n')
env={k:v for k,v in os.environ.items() if not re.search(r'API_KEY|TOKEN|SECRET|PASSWORD|CREDENTIAL|^(ANTHROPIC|OPENAI|GEMINI|GOOGLE|LLM|ENSEMBLE|AWS|AZURE)',k,re.I)};env.update(GOPROXY='off',GOTOOLCHAIN='local')
records=[]
def run(name,args,expected_exit=0):
 args=list(map(str,args));start=time.monotonic();p=subprocess.Popen(args,cwd=M,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True);timeout=False
 try:out=p.communicate(timeout=90)[0]
 except subprocess.TimeoutExpired:
  timeout=True;os.killpg(p.pid,signal.SIGKILL);out=p.communicate()[0]
 (E/(name+'.log')).write_bytes(out)
 r={'name':name,'argv':args,'cwd':str(M),'exit_code':p.returncode,'timeout':timeout,'seconds':round(time.monotonic()-start,3),'output_sha256':hashlib.sha256(out).hexdigest()};records.append(r);(E/'runs.json').write_text(json.dumps(records,indent=2)+'\n');print(name,p.returncode,r['seconds'],flush=True)
 assert not timeout and p.returncode==expected_exit,r
 return out
run('race',['go','test','-race','-count=1','-v','./cmd/grade-edition3'])
run('vet',['go','vet','./cmd/grade-edition3'])
B=S/'grade-edition3';run('build',['go','build','-o',B,'./cmd/grade-edition3'])
marker=S/'invoked';target=S/'sentinel';target.write_text('#!/bin/sh\nprintf invoked >> '+str(marker)+'\nexit 1\n');target.chmod(0o700)
cases=[('help-short',['-h'],0),('help-long',['-help'],0),('help-ch1-target',['-ch','1','-h',target],0),('help-ch4-target',['-ch','4','-h',target],0),('unknown-flag',['-unknown',target],2),('missing-target',['-ch','4'],2),('ch1-empty-core',['-ch','1','-core=',target],2),('ch1-empty-consumer',['-ch','1','-consumer=',target],2),('ch1-empty-case',['-ch','1','-consumer-case=',target],2)]
for name,args,code in cases:
 out=run(name,[B,*args],code);assert not marker.exists(),name
 if name.startswith('help-'):assert b'Usage of grade-edition3:' in out,name
# Reverse only production fix in scratch; leave new test intact.
p=M/'cmd/grade-edition3/main.go';fixed=p.read_bytes()
try:
 p.write_bytes(fixed.replace(branch,b'',1))
 out=run('reverted-help-mutation',['go','test','-count=1','-run','^TestHelpSucceedsWithoutLaunchingTarget$','-v','./cmd/grade-edition3'],1)
 assert out.count(b' = 2, want 0')==4,out
finally:p.write_bytes(fixed)
assert sha(p)==expected['cmd/grade-edition3/main.go']
run('restored-help-test',['go','test','-count=1','-run','^TestHelpSucceedsWithoutLaunchingTarget$','-v','./cmd/grade-edition3'])
# Carried evidence remains evidence for the PRIOR adapter, not this new binary.
carry=[]
for r in json.loads((PRE/'runs.json').read_text()):
 if r['name'] in ['adapter-ch4-executable-default','adapter-ch4-directory-explicit-case','adapter-ch4-wrong-core','adapter-ch1-same-binary-lifetime','adapter-case-forwarding-negative','misrouted-adapter-control']:
  f=PRE/(r['name']+'.log');assert sha(f)==r['output_sha256'];carry.append({'source':str(f),'receipt':r,'not_rerun_on_r1':True})
(E/'carried-route-evidence.json').write_text(json.dumps({'justification':'exact byte-level delta only changes flag.ErrHelp return; all valid-route code and historical graders unchanged','prior_adapter_sha256':json.loads((PRE/'binary-identities.json').read_text())[str(old/'grade-edition3')],'same_reference_sha256':json.loads((PRE/'completed.json').read_text())['paired_reference_sha256'],'prior_report_sha256':sha(ROOT/'docs/edition-3-notes/ch04-grader-review.md'),'controls':carry},indent=2)+'\n')
for name,h in expected.items():assert sha(ROOT/name)==h
(E/'completed.json').write_text(json.dumps({'candidate_unchanged':True,'new_binary_sha256':sha(B),'all_expected_outcomes':True,'sentinel_launched':marker.exists(),'at_utc':datetime.datetime.now(datetime.timezone.utc).isoformat()},indent=2)+'\n');print('COMPLETE',flush=True)
