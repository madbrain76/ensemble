#!/usr/bin/env python3
"""Local-only evidence assertions; never loads real settings or calls providers."""
import copy, importlib.util, json, pathlib, subprocess, tempfile, urllib.request, urllib.error
from unittest import mock
P=pathlib.Path('/tmp/cadence-ch05-evidence-r1')
spec=importlib.util.spec_from_file_location('launch',P/'launch.py');m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m)
out=P/'launcher-controls-r2'
if out.exists():raise SystemExit('refusing overwrite')
out.mkdir();results=[]
def reject(name,fn):
    try:fn()
    except (RuntimeError,ValueError,KeyError) as err:results.append({'case':name,'rejected':True,'reason':str(err)});return
    raise AssertionError('negative control incorrectly accepted: '+name)
# The actual loader is exercised only with read_text replaced in memory.
dummy={'directClaudeAPIKey':'dummy-anthropic','directOpenAIAPIKey':'dummy-openai','directGeminiAPIKey':'dummy-gemini'};reads=[]
def fixture(path,*args,**kwargs):
    assert path==pathlib.Path.home()/'.cr/settings.json','unexpected credential source';reads.append(str(path));return json.dumps(dummy)
with mock.patch.object(pathlib.Path,'read_text',fixture):
    for vendor in ['anthropic','openai','gemini']:assert m.key_for(vendor)=='dummy-'+vendor
assert len(reads)==3
results.append({'case':'approved-settings-key-selection','passed':True,'source':'~/.cr/settings.json','real_credentials_read':False})
reject('no-environment-key-alias-fallback',lambda:m.select_key({'OPENAI_API_KEY':'dummy'},'openai'))
reject('empty-selected-key',lambda:m.select_key({'directOpenAIAPIKey':'  '},'openai'))
assert m.validate_image(json.dumps({'assistant':'{"left":"red","right":"blue"}'}))=={'left':'red','right':'blue'}
reject('reversed-image-orientation',lambda:m.validate_image(json.dumps({'assistant':'{"left":"blue","right":"red"}'})))
reject('image-color-words-without-orientation',lambda:m.validate_image(json.dumps({'assistant':'{"colors":["red","blue"],"sides":["left","right"]}'})))
fixture_dir=P/'local-launcher-r5/workflow';rows=[json.loads(x) for x in (fixture_dir/'stdout.txt').read_text().splitlines()]
with tempfile.TemporaryDirectory(prefix='cadence-ch05-control-') as temp:
    root=pathlib.Path(temp)
    def reset():
        for name in ['manuscript.txt','review.json']:(root/name).write_bytes((fixture_dir/name).read_bytes())
    reset();proof=m.validate_workflow(root,rows);results.append({'case':'nonempty-correlated-workflow','passed':True,'proof':proof})
    reject('fabricated-done-markers',lambda:m.validate_workflow(root,[{'pipeline':n} for n in ['author_done','editor_done','reviewer_done']]))
    (root/'manuscript.txt').write_text('');reject('empty-manuscript-file',lambda:m.validate_workflow(root,rows));reset()
    changed=copy.deepcopy(rows)
    for r in changed:
        if r.get('part',{}).get('type')=='tool_call' and r['part'].get('name')=='edit_draft':r['part']['args']['text']=''
    (root/'manuscript.txt').write_text('');reject('empty-edit-with-matching-file',lambda:m.validate_workflow(root,changed));reset()
    authored=next(r['part']['args']['text'] for r in rows if r.get('part',{}).get('type')=='tool_call' and r['part'].get('name')=='write_draft')
    changed=copy.deepcopy(rows)
    for r in changed:
        if r.get('part',{}).get('type')=='tool_call' and r['part'].get('name')=='edit_draft':r['part']['args']['text']=authored
    (root/'manuscript.txt').write_text(authored);reject('unchanged-editor-manuscript',lambda:m.validate_workflow(root,changed));reset()
    changed=copy.deepcopy(rows)
    for r in changed:
        if r.get('part',{}).get('type')=='tool_call' and r['part'].get('name')=='review':r['part']['args']['notes']=''
    reject('empty-review-notes',lambda:m.validate_workflow(root,changed))
    changed=[r for r in rows if not (r.get('part',{}).get('type')=='tool_call' and r['part'].get('name')=='edit_draft')]
    reject('uncorrelated-tool-result',lambda:m.validate_workflow(root,changed))
# Two local attempts, but only one reaches the upstream fake.
fake=m.Fake('counter');forward=m.Forward(fake.url,1)
try:
    def request():return urllib.request.urlopen(urllib.request.Request(forward.url+'/v1/test',data=b'{}',headers={'Content-Type':'application/json'},method='POST'),timeout=5)
    with request() as response:assert response.status==200;response.read()
    try:
        with request() as response:response.read()
        raise AssertionError('budget allowed extra upstream request')
    except urllib.error.HTTPError as err:assert err.code==429
    assert forward.forwarded==1 and forward.attempts==2 and len(fake.requests)==1
    results.append({'case':'hard-request-cap','passed':True,'attempts':2,'upstream_requests':1,'second_status':429})
finally:forward.close(out/'cap');fake.close()
identity=json.loads((P/'frozen/identity.json').read_text())
for name in ['wrong-source','wrong-binary']:
    changed=copy.deepcopy(identity)
    if name=='wrong-source':changed['source_files']['ensemble.go']='0'*64
    else:changed['binaries']['cli']['sha256']='0'*64
    path=out/(name+'.json');m.save(path,changed)
    result=subprocess.run(['python3',str(P/'launch.py'),'--manifest',str(path),'--verify-only'],capture_output=True,text=True)
    (out/(name+'.txt')).write_text(result.stdout+result.stderr)
    assert result.returncode!=0 and 'identity refused' in result.stdout
    results.append({'case':name,'rejected':True,'exit_code':result.returncode})
m.save(out/'results.json',results)
m.save(out/'identity.json',{'launcher_sha256':m.sha(P/'launch.py'),'controls_sha256':m.sha(__file__),'frozen_manifest_sha256':m.sha(P/'frozen/identity.json'),'fixture_stdout_sha256':m.sha(fixture_dir/'stdout.txt'),'real_credentials_read':False,'paid_calls':0})
print(json.dumps({'passed':True,'controls':len(results),'paid_calls':0}))
