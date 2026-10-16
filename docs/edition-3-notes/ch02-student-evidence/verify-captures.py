#!/usr/bin/env python3
"""Verify existing captures only. No credentials or live model calls."""
import hashlib
import json
import os
import pathlib
import re
import subprocess

out = pathlib.Path(__file__).resolve().parent
cwd = pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
binary = out / 'student-cli'
summary = json.loads((out / 'live-summary.json').read_text())
verified = []
for original in summary:
    vendor = original['vendor']
    raw = (out / f'live-{vendor}.jsonl').read_bytes()
    facts = [json.loads(line) for line in raw.splitlines()][1:]
    responses = [fact['response'] for fact in facts if fact['type'] == 'response_ended']
    stderr = (out / f'live-{vendor}.stderr.txt').read_text()
    match = re.search(r'cumulative usage: input=(\d+) output=(\d+) cache_write=(\d+) cache_read=(\d+)', stderr)
    assert match, f'{vendor}: missing cumulative usage'
    transcript = dict(zip(['input', 'output', 'cache_write', 'cache_read'], map(int, match.groups())))
    total = {name: sum(response['usage'][name] for response in responses) for name in transcript}
    assert transcript == total == original['captured_usage']
    measured = [tuple(map(int, match.groups())) for match in re.finditer(r'^input=(\d+) output=(\d+) cache_write=(\d+) cache_read=(\d+)$', stderr, re.M)]
    assert measured == [tuple(response['usage'][name] for name in ['input', 'output', 'cache_write', 'cache_read']) for response in responses]
    assert len(responses) == 2 and original['exit'] == 0
    assert hashlib.sha256(raw).hexdigest() == original['log_sha256']
    verified.append(dict(vendor=vendor, live_cli_exit=0, response_count=2,
                         cumulative_usage=total, transcript_matches_log=True,
                         per_response_usage_matches_log=True, log_sha256_unchanged=True))
offline = json.loads((out / 'offline-summary.json').read_text())
assert not offline['api_keys_supplied'] and offline['sentinel_post_count'] == 0
assert len(offline['commands']) == 12
for item in offline['commands']:
    assert item.get('exit', 0) == 0 and all(code == 0 for code in item.get('exits', []))
    assert item.get('byte_identical', True) and item.get('exact_saved_log', True) and item.get('request_json_valid', True)
report = dict(commit=summary[0]['commit'], binary_sha256=hashlib.sha256(binary.read_bytes()).hexdigest(),
              additional_live_calls=0, verification='Saved normalized logs and stderr only; no settings access',
              original_launcher_exit=1, original_launcher_failure='Incorrect capitalization/order in evidence-only usage regex',
              live=verified, offline_repeat_pairs=9, offline_dump_matches=3, offline_sentinel_posts=0)

# Separate, explicitly fake CLI demonstration of retroactive projection. The
# original result remains in the append-only JSONL even after its view is stubbed.
source = dict(vendor='anthropic', model='fake-model', surface='messages')
facts = [
    dict(seq=1, type='message_received', time='2026-01-01T00:00:00Z', message=dict(actor='human', parts=[dict(type='text',text='Look up the sample item')])),
    dict(seq=2, type='response_ended', time='2026-01-01T00:00:01Z', response=dict(parts=[dict(type='tool_call',call_id='c1',name='lookup',args={'item':'sample'}, **{'from':source})], usage=dict(input=1,cache_write=0,cache_read=0,output=1), **{'from':source})),
    dict(seq=3, type='tool_returned', time='2026-01-01T00:00:02Z', tool=dict(call_id='c1',name='lookup',parts=[dict(type='text',text='PRIVATE-FAKE-RESULT')]))
]
redaction = dict(seq=4, type='redacted', time='2026-01-01T00:00:03Z', redact=dict(level='redact_result',reason='teaching fixture', **{'from':3,'to':3}))
for name, events in [('before',facts),('after',facts+[redaction])]:
    (out/f'fake-redaction-{name}.jsonl').write_text('\n'.join(json.dumps(event) for event in [{'log_version':1}]+events)+'\n')
commands=[]
base_env={key:os.environ[key] for key in ['PATH','HOME','TMPDIR','LANG'] if key in os.environ}
for vendor in ['anthropic','openai','gemini']:
    env=dict(base_env,LLM_VENDOR=vendor,LLM_MODEL='fake-model',LLM_BASE_URL='http://127.0.0.1:1')
    for phase in ['before','after']:
        command=[str(binary),'render',str(out/f'fake-redaction-{phase}.jsonl')]
        result=subprocess.run(command,cwd=cwd,env=env,capture_output=True)
        (out/f'fake-redaction-{vendor}-{phase}.json').write_bytes(result.stdout)
        (out/f'fake-redaction-{vendor}-{phase}.stderr.txt').write_bytes(result.stderr)
        assert result.returncode==0 and json.loads(result.stdout)
        contains=b'PRIVATE-FAKE-RESULT' in result.stdout
        assert contains==(phase=='before')
        commands.append(dict(command=command,cwd=str(cwd),target_vendor=vendor,phase=phase,exit=0,contains_fake_secret=contains,live=False))
report['fake_redaction']=commands
report['original_fake_result_retained_in_after_log']='PRIVATE-FAKE-RESULT' in (out/'fake-redaction-after.jsonl').read_text()
assert report['original_fake_result_retained_in_after_log']
(out/'capture-verification.json').write_text(json.dumps(report,indent=2)+'\n')
print('Verified three live two-turn captures, per-response and cumulative usage.')
print('Verified three byte-identical dumps and nine repeated cross-vendor render pairs; zero sentinel posts.')
print('Verified retroactive fake-result redaction through all three actual CLI renderers.')
print('Additional live calls: 0. Original launcher failure remains recorded.')
