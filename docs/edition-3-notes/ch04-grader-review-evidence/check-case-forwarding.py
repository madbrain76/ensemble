#!/usr/bin/env python3
"""One additional negative proving caller-supplied case reaches existing grader."""
import hashlib,json,os,pathlib,re,signal,subprocess,time
ROOT=pathlib.Path('/Users/bill/projects/ensemble');E=ROOT/'docs/edition-3-notes/ch04-grader-review-evidence'
assert (E/'completed.json').exists(),'run only after primary finite matrix completes'
S=pathlib.Path((E/'scratch.txt').read_text().strip());R=S/'reference-lifetime-corrected'
case=S/'wrong-expectation.json';d=json.loads((E/'explicit-case.json').read_text());d['result']='DELIBERATELY WRONG EXPECTATION';case.write_text(json.dumps(d)+'\n');(E/'wrong-expectation.json').write_bytes(case.read_bytes())
env={k:v for k,v in os.environ.items() if not re.search(r'API_KEY|TOKEN|SECRET|PASSWORD|CREDENTIAL|^(ANTHROPIC|OPENAI|GEMINI|GOOGLE|LLM|ENSEMBLE|AWS|AZURE)',k,re.I)};env.update(GOPROXY='off',GOTOOLCHAIN='local')
args=list(map(str,[S/'grade-edition3','-ch','4','-json','-core',R,'-consumer',R/'ch05','-consumer-case',case,S/'reference-cli']))
start=time.monotonic();proc=subprocess.Popen(args,cwd=S/'main',env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True);timeout=False
try:out=proc.communicate(timeout=90)[0]
except subprocess.TimeoutExpired:
 timeout=True;os.killpg(proc.pid,signal.SIGKILL);out=proc.communicate()[0]
p=E/'adapter-case-forwarding-negative.log';p.write_bytes(out)
rec={'name':'adapter-case-forwarding-negative','argv':args,'cwd':str(S/'main'),'exit_code':proc.returncode,'timeout':timeout,'seconds':round(time.monotonic()-start,3),'output_sha256':hashlib.sha256(out).hexdigest()}
records=json.loads((E/'runs.json').read_text());records.append(rec);(E/'runs.json').write_text(json.dumps(records,indent=2)+'\n')
assert not timeout and proc.returncode==1,rec
d=json.loads(out);assert (d['score'],d['max_score'],d['passed'])==(135,160,False),d
assert [c['id'] for c in d['checks'] if not c['passed']]==['external-consumer'],d
(E/'case-forwarding-completed.json').write_text(json.dumps({'expected_failure_detected':True,'case_kind':'valid caller fixture with deliberately wrong expected result'},indent=2)+'\n')
print(rec)
