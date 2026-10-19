#!/usr/bin/env python3
import datetime, difflib, hashlib, json, os, pathlib, shutil, stat, subprocess, sys, time
E = pathlib.Path(__file__).resolve().parent
P = json.loads((E/'plan-v2.json').read_text())
C = pathlib.Path(P['source'])
M = pathlib.Path('/Users/bill/projects/ensemble')
INPUT = json.loads((M/'docs/edition-3-notes/ch10-retention-review-input.json').read_text())
name = sys.argv[1]
assert name in P['commands']
assert not (E/(name+'.json')).exists(), 'No reruns'
def verify():
    for v in INPUT['candidate_files']:
        f=C/v['path']; b=f.read_bytes()
        assert len(b)==v['bytes'] and hashlib.sha256(b).hexdigest()==v['sha256'], str(f)
        assert oct(stat.S_IMODE(f.stat().st_mode))==v['filesystem_mode'],str(f)
    actual={str(f.relative_to(C)) for f in C.rglob('*') if f.is_file()}
    assert actual=={v['path'] for v in INPUT['candidate_files']}
    return {'members':len(actual),'bytes_modes':'all match pinned candidate','extra_files':[]}
def green(n):
    assert json.loads((E/(n+'.json')).read_text())['exit_code']==0,n
pre=verify()
if name=='gui': green('root')
if name=='negative': green('root'); green('gui')
if name=='restored':
    prior=json.loads((E/'negative.json').read_text())
    assert prior['intended_negative'] and prior['restoration']['members']==142
free=shutil.disk_usage(C).free
assert free >= 1024**3, '1GiB disk floor'
original=None
if name=='negative':
    n=P['negative']; f=C/n['file']; original=f.read_bytes()
    old=n['old'].encode(); new=n['new'].encode(); assert original.count(old)==1
    mutant=original.replace(old,new)
    (E/'negative.diff').write_text(''.join(difflib.unified_diff(original.decode().splitlines(True),mutant.decode().splitlines(True),fromfile='candidate/'+n['file'],tofile='mutant/'+n['file'])))
    f.write_bytes(mutant)
cwd=C/'gui' if name=='gui' else C
meta={'name':name,'cwd':str(cwd),'argv':P['commands'][name],'environment':P['environment'],'source_precheck':pre,'disk_free_before':free,'started_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'merged_output':name+'.txt'}
if original is not None:
    meta['mutant_sha256']=hashlib.sha256(mutant).hexdigest()
    meta['original_sha256']=hashlib.sha256(original).hexdigest()
t=time.monotonic()
try:
    with (E/(name+'.txt')).open('wb') as out:
        proc=subprocess.run(P['commands'][name],cwd=cwd,env=P['environment'],stdout=out,stderr=subprocess.STDOUT)
    meta['exit_code']=proc.returncode
finally:
    meta['elapsed_seconds']=round(time.monotonic()-t,3)
    meta['finished_utc']=datetime.datetime.now(datetime.timezone.utc).isoformat()
    if original is not None: f.write_bytes(original)
    meta['restoration' if original is not None else 'source_postcheck']=verify()
    meta['disk_free_after']=shutil.disk_usage(C).free
    raw=(E/(name+'.txt')).read_text()
    meta['run_cases']=[line.removeprefix('=== RUN   ') for line in raw.splitlines() if line.startswith('=== RUN   ')]
    meta['pass_cases']=[line.split(' (')[0].removeprefix('--- PASS: ') for line in raw.splitlines() if line.startswith('--- PASS: ')]
    meta['fail_cases']=[line.split(' (')[0].removeprefix('--- FAIL: ') for line in raw.splitlines() if line.startswith('--- FAIL: ')]
    if name=='negative':
        meta['intended_negative']=meta.get('exit_code')==1 and P['negative']['expected_assertion'] in raw and meta['fail_cases']==['TestRetentionDialogueArchiveAndModeSwitch'] and '[build failed]' not in raw and 'panic:' not in raw and 'DATA RACE' not in raw
    (E/(name+'.json')).write_text(json.dumps(meta,indent=2)+'\n')
print(raw,end='')
print('RECEIPT:',json.dumps({k:meta.get(k) for k in ['name','exit_code','elapsed_seconds','run_cases','pass_cases','fail_cases','intended_negative']}))
if name=='negative': assert meta['intended_negative'],'Unexpected negative result; STOP'
else:
    expected=P['root_cases'] if name=='root' else P['gui_cases'] if name=='gui' else ['TestRetentionDialogueArchiveAndModeSwitch']
    assert meta['exit_code']==0 and sorted(meta['pass_cases'])==sorted(expected),'Unexpected assertion/setup/selection failure; STOP'
