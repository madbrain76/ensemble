import importlib.util,json,pathlib,shutil,subprocess
HERE=pathlib.Path(__file__).resolve().parent
spec=importlib.util.spec_from_file_location('baseline',HERE/'run-baseline.py');b=importlib.util.module_from_spec(spec);spec.loader.exec_module(b)
T=pathlib.Path((b.E/'scratch.txt').read_text().strip());M=T/'sensitivity';P=T/'public-consumer';P.mkdir(exist_ok=True);records=json.loads((b.E/'sensitivity-runs.json').read_text()) if (b.E/'sensitivity-runs.json').exists() else []
oldmanifest=json.loads((b.OLD/'manifest.json').read_text());oldentry=next(x for x in oldmanifest['files'] if x['path']=='public_api_test.go.txt');probe=b.OLD/'public_api_test.go.txt';assert b.sha(probe)==oldentry['sha256'];shutil.copyfile(probe,b.E/'reused-public-probe.go.txt');shutil.copyfile(probe,P/'public_api_test.go');shutil.copyfile(M/'go.sum',P/'go.sum');(P/'go.mod').write_text(f'module independent-review\n\ngo 1.25\n\nrequire ensemble v0.0.0\nrequire github.com/creack/pty v1.1.24 // indirect\nreplace ensemble => {M}\n');(b.E/'public-consumer-module.txt').write_text((P/'go.mod').read_text())
file=M/'internal/agent/dispatch.go';fixed=file.read_bytes();old=subprocess.check_output(['git','--no-pager','-C',str(b.S),'show','062c146cafa9b105745f6b44921160198bbaee32:internal/agent/dispatch.go']);assert fixed!=old
meta={'corrected_dispatch_sha256':b.sha(file),'reverted_dispatch_sha256':__import__('hashlib').sha256(old).hexdigest(),'regression_test_sha256':b.sha(M/'chapter4_limits_test.go'),'public_probe_sha256':b.sha(probe),'reversion_is_exact_initial_candidate_file':True};(b.E/'sensitivity-inputs.json').write_text(json.dumps(meta,indent=2)+'\n')
def check(name,where,selector,expected,race=False):
 previous=next((x for x in records if x['name']==name),None)
 if previous:
  assert previous['exit_code']==expected and not previous['timeout'];return (b.E/(name+'.log')).read_bytes()
 args=['go','test']+(['-race'] if race else [])+['-count=1','-v','-run',selector,'.'];row,out=b.run(name,args,where,60,records,'sensitivity-runs.json');assert row['exit_code']==expected and not row['timeout'];assert b'[build failed]' not in out and b'[no tests to run]' not in out
 return out
check('public-corrected',P,'^TestPublic(PendingRearm|PendingResetOverrideAndIsolation)$',0,True)
check('regressions-corrected',M,'^TestChapter4ChainedLimits',0,True)
try:
 file.write_bytes(old)
 out=check('regressions-reverted',M,'^TestChapter4ChainedLimits',1)
 for marker in [b'omitted cap inherited',b'--- FAIL: TestChapter4ChainedLimitsUseFreshDefaults/omitted_delay',b'--- FAIL: TestChapter4ChainedLimitsUseFreshDefaults/omitted_pattern']:assert marker in out,marker
 out=check('public-reverted',P,'^TestPublicPendingRearm$',1);assert b'old 1024-byte pending cap survived consumption' in out
finally:file.write_bytes(fixed)
assert b.sha(file)==meta['corrected_dispatch_sha256']
check('regressions-restored',M,'^TestChapter4ChainedLimits',0,True)
check('public-restored',P,'^TestPublic(PendingRearm|PendingResetOverrideAndIsolation)$',0,True)
for x in json.loads((b.SE/'source-manifest.json').read_text())['files']:assert b.sha(M/x['path'])==x['sha256']
(b.E/'sensitivity-complete.json').write_text(json.dumps({'behavioral_reversion_failures_verified':True,'restored_source_exact':True,'candidate_source_never_edited':True},indent=2)+'\n')
