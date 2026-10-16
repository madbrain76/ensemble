"""Local real-CLI observation, not live-provider or exhaustive acceptance proof."""
import hashlib, http.server, json, os, pathlib, subprocess, tempfile, threading
binary=pathlib.Path('/tmp/ensemble-ch03-evidence/student-cli')
assert hashlib.sha256(binary.read_bytes()).hexdigest()=='e2ce182d3220ffa05c3584b9c92ef84536925ffbbe0f88f5bccba9ac29b2aa2c'
with tempfile.TemporaryDirectory(prefix='cadence-ch03-chat-') as td:
    work=pathlib.Path(td); requests=[]; observations=[]; errors=[]
    class Handler(http.server.BaseHTTPRequestHandler):
        def log_message(self,*args): pass
        def do_POST(self):
            try:
                body=json.loads(self.rfile.read(int(self.headers['Content-Length'])))
                requests.append(body); n=len(requests)
                assert {t['name'] for t in body['tools']}=={'read_file','write_file','edit_file','list_directory','search_files','run_command'}
                def call(i,path):return {'type':'tool_use','id':i,'name':'write_file','input':{'path':path,'content':'Juniper local fixture\n'}}
                if n==1: content=[call('rejected-a','rejected-a.txt'),call('rejected-b','rejected-b.txt')]
                elif n==2:
                    results=[b for m in body['messages'] for b in m.get('content',[]) if isinstance(b,dict) and b.get('type')=='tool_result']
                    assert len(results)==2 and all(b.get('is_error') for b in results)
                    assert not (work/'rejected-a.txt').exists() and not (work/'rejected-b.txt').exists()
                    observations.append('unexplained whole batch refused with two paired errors and zero writes')
                    content=[{'type':'text','text':'I will write the requested local marker now.'},call('accepted','marker.txt')]
                elif n==3:
                    assert (work/'marker.txt').read_text()=='Juniper local fixture\n'
                    assert 'accepted' in json.dumps(body['messages'])
                    observations.append('explained retry executed and actual file verified')
                    content=[{'type':'text','text':'Wrote Juniper local fixture.'}]
                elif n==4:
                    assert 'Please write the local marker.' in json.dumps(body['messages']) and 'Wrote Juniper local fixture.' in json.dumps(body['messages'])
                    observations.append('second human request includes original user and assistant text')
                    content=[{'type':'text','text':'The marker is Juniper local fixture.'}]
                else: raise AssertionError('unexpected extra provider request')
                data=json.dumps({'type':'message','role':'assistant','model':'fixture-only','content':content,'usage':{'input_tokens':10,'output_tokens':3}}).encode()
                self.send_response(200);self.send_header('Content-Type','application/json');self.end_headers();self.wfile.write(data)
            except Exception as e:
                errors.append(type(e).__name__+': '+str(e));self.send_response(500);self.end_headers()
    server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);thread=threading.Thread(target=server.serve_forever,daemon=True);thread.start()
    try:
        env={'PATH':os.environ['PATH'],'HOME':td,'LLM_VENDOR':'anthropic','LLM_MODEL':'fixture-only','LLM_API_KEY':'local-fixture-not-a-secret','LLM_BASE_URL':f'http://127.0.0.1:{server.server_port}','LLM_REQUIRE_TOOL_EXPLANATION':'true'}
        result=subprocess.run([str(binary),'chat'],input='Please write the local marker.\nWhat marker did you write?\n',text=True,capture_output=True,cwd=td,env=env,timeout=30)
    finally: server.shutdown();server.server_close();thread.join()
    receipt={'candidate':'0eb6fdffe62e662b0ab96982b01287afba19ff42','binary_sha256':hashlib.sha256(binary.read_bytes()).hexdigest(),'kind':'local fake-provider real chat CLI; not paid/live-provider proof','external_watchdog_seconds':30,'requests':len(requests),'exit_code':result.returncode,'stdout':result.stdout,'stderr':result.stderr,'observations':observations,'fixture_errors':errors}
    out=pathlib.Path(__file__).with_name('chat-smoke.json');out.write_text(json.dumps(receipt,indent=2)+'\n')
    assert not errors and result.returncode==0 and len(requests)==4 and len(observations)==3,receipt
    assert 'Wrote Juniper local fixture.' in result.stdout and 'The marker is Juniper local fixture.' in result.stdout
    print(json.dumps(receipt,indent=2))
