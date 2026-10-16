#!/usr/bin/env python3
"""No credential reads or paid calls: offline CLI checks and local fake responses."""
import hashlib
import http.server
import json
import os
from pathlib import Path
import subprocess
import threading

ROOT = Path('/Users/bill/projects/ensemble-edition-3-student')
OUT = Path(__file__).resolve().parent
OLD = OUT.parent/'ch02-student-evidence'
CLI = OUT/'student-cli'
COMMIT = '89cc874b3c32e65ee43f44a72ad9ca48c9dd3323'
SHA = 'b9b29b73e0da79393dc8c499aaee1437e73cfcb260da57b1ee1fef7f4325624a'
assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip() == COMMIT
assert subprocess.check_output(['git','status','--porcelain'],cwd=ROOT) == b''
assert hashlib.sha256(CLI.read_bytes()).hexdigest() == SHA
original = {str(p): hashlib.sha256(p.read_bytes()).hexdigest() for p in OLD.rglob('*') if p.is_file()}
review = OUT.parent/'ch02-student-review.md'
original[str(review)] = hashlib.sha256(review.read_bytes()).hexdigest()
requests = []
bodies = {
 '/v1/messages': {'content':[], 'usage':{'input_tokens':11,'output_tokens':2}},
 '/v1/chat/completions': {'choices':[], 'usage':{'prompt_tokens':11,'completion_tokens':2}},
 '/v1beta/models/exact-model:generateContent': {'candidates':[], 'usageMetadata':{'promptTokenCount':11,'candidatesTokenCount':2}},
}
class Handler(http.server.BaseHTTPRequestHandler):
    def log_message(self, *args): pass
    def do_POST(self):
        self.rfile.read(int(self.headers.get('Content-Length','0')))
        requests.append(self.path)
        data=json.dumps(bodies[self.path]).encode()
        self.send_response(200);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(data)));self.end_headers();self.wfile.write(data)
    def do_GET(self):
        requests.append(self.path);self.send_response(500);self.end_headers()
server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler)
thread=threading.Thread(target=server.serve_forever,daemon=True);thread.start()
base={k:os.environ[k] for k in ('PATH','HOME','TMPDIR','LANG') if k in os.environ}
base['LLM_BASE_URL']='http://127.0.0.1:'+str(server.server_address[1])
records=[]
def run(label,args,vendor,model='exact-model',binary=CLI,stdin='',fake=False,log=None):
    env=base | {'LLM_VENDOR':vendor,'LLM_MODEL':model}
    if fake: env['LLM_API_KEY']='fake-local-only'
    if log: env['CH02_LOG']=str(log)
    result=subprocess.run([str(binary)]+args,cwd=ROOT,env=env,input=stdin,text=True,capture_output=True)
    (OUT/(label+'.stdout')).write_text(result.stdout)
    (OUT/(label+'.stderr')).write_text(result.stderr)
    records.append({'label':label,'argv':[str(binary)]+args,'cwd':str(ROOT),'vendor':vendor,'model':model,'exit':result.returncode,'real_credentials':False,'fake_local_transport':fake})
    return result

def render(label,log,vendor,model='exact-model',binary=CLI):
    result=run(label,['render',str(log)],vendor,model,binary)
    assert result.returncode==0,(label,result.stderr)
    return result.stdout,json.loads(result.stdout)

def save_log(name,events):
    path=OUT/name
    path.write_text(json.dumps({'log_version':1})+'\n'+''.join(json.dumps(dict(seq=i+1,time='2026-01-01T00:00:00Z',**event))+'\n' for i,event in enumerate(events)))
    return path

models={'anthropic':'claude-haiku-4-5-20251001','openai':'gpt-4.1-nano','gemini':'gemini-2.5-flash-lite'}
try:
    for origin in models:
        path=OLD/('live-'+origin+'.jsonl')
        dumped=run('dump-original-'+origin,['dump'],origin,models[origin],log=path)
        assert dumped.returncode==0 and dumped.stdout==path.read_text()
        for target,model in models.items():
            stem='replay-'+origin+'-to-'+target
            first,_=render(stem+'-new1',path,target,model)
            second,_=render(stem+'-new2',path,target,model)
            baseline,_=render(stem+'-old-binary-now',path,target,model,OLD/'student-cli')
            assert first==second==baseline,stem
    from_g={'vendor':'gemini','model':'exact-model','surface':'generate_content'}
    events=[]
    for name in ('lookup_weather','lookup_stock'):
        events += [
            {'type':'response_ended','response':{'from':from_g,'parts':[{'type':'tool_call','name':name,'args':{},'from':from_g}],'usage':{'input':1,'output':1}}},
            {'type':'tool_returned','tool':{'parts':[{'type':'text','text':name+' result'}]}},
        ]
    fixture=save_log('fake-R1-sequential-idless.jsonl',events)
    _,body=render('R1-sequential-idless',fixture,'gemini')
    names=[part['functionResponse']['name'] for turn in body['contents'] for part in turn['parts'] if 'functionResponse' in part]
    assert names==['lookup_weather','lookup_stock'],names
    from_o={'vendor':'openai','model':'exact-model','surface':'chat_completions'}
    parts=[{'type':'text','text':'answer'}, {'type':'opaque','from':from_o,'data':{'reasoning_content':'opaque-first'}}, {'type':'opaque','from':from_o,'data':{'reasoning':{'text':'opaque-second'}}}]
    fixture=save_log('fake-R3-multiple-opaque.jsonl',[{'type':'response_ended','response':{'from':from_o,'parts':parts,'usage':{'input':1,'output':1}}}])
    _,body=render('R3-exact',fixture,'openai')
    message=body['messages'][1]
    assert message['reasoning_content']=='opaque-first' and message['reasoning']=={'text':'opaque-second'}
    foreign,_=render('R3-other-model',fixture,'openai','other-model')
    assert 'opaque-first' not in foreign and 'opaque-second' not in foreign
    for vendor in models:
        text,_=render('redaction-'+vendor,OLD/'fake-redaction-after.jsonl',vendor)
        assert 'PRIVATE-FAKE-RESULT' not in text
    assert requests==[],requests
    offline_requests=len(requests)
    for vendor in models:
        path=OUT/('fake-R2-empty-'+vendor+'.jsonl')
        result=run('R2-empty-'+vendor,[],vendor,stdin='{"user":"Give a short answer"}\n',fake=True,log=path)
        assert result.returncode!=0,(vendor,result.stdout,result.stderr)
        assert 'input=11 output=2' in result.stderr,(vendor,result.stderr)
        events=[json.loads(line) for line in path.read_text().splitlines()][1:]
        assert [event['type'] for event in events]==['message_received','request_sent','error_occurred'],events
    assert sorted(requests)==sorted(bodies),requests
finally:
    server.shutdown();server.server_close();thread.join()
assert all(hashlib.sha256(Path(name).read_bytes()).hexdigest()==digest for name,digest in original.items())
assert subprocess.check_output(['git','status','--porcelain'],cwd=ROOT)==b''
(OUT/'offline-checks-summary.json').write_text(json.dumps({'candidate_commit':COMMIT,'binary_sha256':SHA,'new_paid_calls':0,'original_live_commit':'74eb0265bd8bdbc7857d579704300b317b2365da','original_artifacts_unchanged':True,'original_artifact_hashes':original,'offline_http_requests':offline_requests,'fake_http_paths':requests,'records':records},indent=2)+'\n')
print('PASS: affected CLI rendering, measured empty-response captures, original-log replay comparison, no offline HTTP requests, three local fake calls, no real credentials or paid calls, original evidence unchanged.')
