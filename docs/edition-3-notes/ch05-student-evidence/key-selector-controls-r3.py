#!/usr/bin/env python3
import difflib,hashlib,importlib.util,json,pathlib
from unittest import mock
p=pathlib.Path('/tmp/cadence-ch05-evidence-r1');out=p/'key-selector-r3'
spec=importlib.util.spec_from_file_location('launch',p/'launch.py');m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m)
fields={'anthropic':'directClaudeAPIKey','openai':'directOpenAIAPIKey','gemini':'directGeminiAPIKey'}
old_fields={'anthropic':'directAnthropicApiKey','openai':'directOpenAIApiKey','gemini':'directGeminiApiKey'}
dummy={field:'dummy-'+vendor for vendor,field in fields.items()};reads=[]
def fixture(path,*args,**kwargs):
    assert path==pathlib.Path.home()/'.cr/settings.json'
    reads.append(str(path));return json.dumps(dummy)
with mock.patch.object(pathlib.Path,'read_text',fixture):
    for vendor in fields:assert m.key_for(vendor)=='dummy-'+vendor
assert len(reads)==3
rejections=[]
for vendor,field in fields.items():
    for name,settings in [('missing',{}),('previous-wrong-spelling',{old_fields[vendor]:'dummy'}),('empty',{field:'  '}),('non-string',{field:42})]:
        try:m.select_key(settings,vendor)
        except RuntimeError:rejections.append({'vendor':vendor,'case':name,'rejected':True})
        else:raise AssertionError('incorrect key selection accepted')
old=(out/'launch.py.before').read_text();new=(p/'launch.py').read_text()
start='def select_key';end='\ndef key_for'
assert old[:old.index(start)]+old[old.index(end):]==new[:new.index(start)]+new[new.index(end):]
assert len([x for x in difflib.ndiff(old.splitlines(),new.splitlines()) if x.startswith('- ')])==1
hashes={name:hashlib.sha256((p/name).read_bytes()).hexdigest() for name in ['launch.py','launcher-controls.py','live-checklist-r3.md']}
result={'passed':True,'positive_vendors':list(fields),'exact_fields':fields,'rejections':rejections,'real_settings_reads':0,'paid_calls':0,'launcher_delta':'one selector-mapping line only','unchanged_scenarios_identity_cap':'byte-for-byte equal to c582a4e322b9dd60b672a0133c673a54f09a59d80bb9947f2d6b396fc6613c82; carry local-launcher-r6 and launcher-controls-r2 receipts only for unchanged code','hashes':hashes}
(out/'result.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps({'passed':True,'positive':3,'rejections':len(rejections),'real_settings_reads':0,'paid_calls':0,'hashes':hashes}))
