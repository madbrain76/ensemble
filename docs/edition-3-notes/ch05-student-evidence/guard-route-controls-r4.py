#!/usr/bin/env python3
import contextlib,hashlib,importlib.util,io,json,pathlib,sys,tempfile,urllib.request,urllib.error
from unittest import mock
p=pathlib.Path('/tmp/cadence-ch05-evidence-r1');out=p/'guard-route-r4';spec=importlib.util.spec_from_file_location('launch',p/'launch.py');m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m)
results=[]
with tempfile.TemporaryDirectory(prefix='cadence-ch05-guard-control-') as temp:
    root=pathlib.Path(temp);(root/'cr/io').mkdir(parents=True)
    (root/'cr/io/1').write_text('CH05_JOB_STARTED\nCH05_JOB_FINISHED\nexit_code: 0\n')
    (root/'stdout.txt').write_text('No `KEY_ENV_LEAK` was printed, confirming the guard did not fire.')
    m.validate_process_guard(root);results.append({'case':'benign-assistant-prose','passed':True})
    (root/'cr/io/1').write_text('KEY_ENV_LEAK\nexit_code: 3\n')
    try:m.validate_process_guard(root)
    except RuntimeError:results.append({'case':'malicious-process-output','rejected':True})
    else:raise AssertionError('real guard output accepted')
# Reconcile the actual saved human output; do not rerun any inference.
saved=p/'live-r1/human';m.validate_process_guard(saved)
facts=json.loads((p/'live-r1-human-observed-facts.json').read_text())
for name,digest in facts['artifact_hashes'].items():assert m.sha(saved/name)==digest
assert facts['original_launcher_passed'] is False and facts['new_turn_precedes_job_seconds']>0
results.append({'case':'saved-human-guard-reconciliation','passed':True,'original_launcher_passed':False,'new_inference_requests':0,'facts_sha256':m.sha(p/'live-r1-human-observed-facts.json')})
# Exercise actual argument parsing/main dispatch, with any human call forbidden.
with tempfile.TemporaryDirectory(prefix='cadence-ch05-route-control-') as temp:
    called=[]
    def step(name):
        def f(*args,**kwargs):called.append(name)
        return f
    argv=['launch.py','--manifest',str(p/'frozen/identity.json'),'--scenario','remaining','--out',str(pathlib.Path(temp)/'route')]
    with mock.patch.object(sys,'argv',argv),mock.patch.object(m,'human',side_effect=AssertionError('human route forbidden')),mock.patch.object(m,'key_for',side_effect=AssertionError('real credentials forbidden')),mock.patch.object(m,'workflow',side_effect=step('workflow')),mock.patch.object(m,'image',side_effect=step('image')),mock.patch.object(m,'unsupported',side_effect=step('unsupported')),contextlib.redirect_stdout(io.StringIO()):assert m.main() is None
    assert called==['workflow','image','unsupported']
    summary=json.loads((pathlib.Path(temp)/'route/summary.json').read_text());assert summary['human_run'] is False
    results.append({'case':'no-human-route','passed':True,'called':called})
fake=m.Fake('counter');forward=m.Forward(fake.url,1)
try:
    def request():return urllib.request.urlopen(urllib.request.Request(forward.url+'/v1/test',data=b'{}',headers={'Content-Type':'application/json'},method='POST'),timeout=5)
    with request() as response:assert response.status==200;response.read()
    try:
        with request() as response:response.read()
        raise AssertionError('cap allowed extra upstream request')
    except urllib.error.HTTPError as err:assert err.code==429
    assert forward.forwarded==1 and forward.attempts==2 and len(fake.requests)==1
    results.append({'case':'hard-cap','passed':True,'attempts':2,'upstream':1,'second_status':429})
finally:forward.close(out/'cap');fake.close()
old=(out/'launch.py.before').read_text();new=(p/'launch.py').read_text()
for start,end in [('def select_key','def environment'),('class NoRedirect','class Capture'),('def validate_workflow','def main')]:assert old[old.index(start):old.index(end)]==new[new.index(start):new.index(end)]
m.save(out/'results.json',{'passed':True,'results':results,'real_settings_reads':0,'additional_paid_calls':0,'launcher_sha256':m.sha(p/'launch.py'),'controls_sha256':m.sha(__file__),'production_manifest_sha256':m.sha(p/'frozen/identity.json'),'scope':'guard and main route only; selected key/cap/model/remaining scenario bodies unchanged'})
print(json.dumps({'passed':True,'controls':len(results),'additional_paid_calls':0}))
