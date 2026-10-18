#!/usr/bin/env python3
"""Owned operational harness. Rehearsal is loopback-only; serving requires later approval.
No browser launcher, CDP client, voice selector, microphone or OS-settings API.
"""
import argparse, contextlib, datetime, hashlib, http.client, http.server, io, json, os
from pathlib import Path
import re, signal, socket, subprocess, sys, threading, time, types

def sha(p): return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def save(p, value): Path(p).write_text(json.dumps(value, indent=2)+'\n')
def clean_env(home):
    env={k:os.environ[k] for k in ('PATH','TMPDIR','LANG') if k in os.environ}
    env.update(HOME=str(home), ANTHROPIC_API_KEY='owned-loopback-dummy');return env

def verify(policy_path, expected):
    if sha(policy_path)!=expected: raise ValueError('policy pin mismatch')
    p=json.loads(Path(policy_path).read_text());root=Path(p['root'])
    if sha(__file__)!=p['launcher_sha256']: raise ValueError('launcher pin mismatch')
    for script,digest in p['scripts'].items():
        if sha(script)!=digest: raise ValueError('harness script pin mismatch')
    for key,ref in [('HEAD',p['commit']),('HEAD^{tree}',p['tree'])]:
        if subprocess.check_output(['git','--no-pager','rev-parse',key],cwd=root,text=True).strip()!=ref: raise ValueError('candidate pin mismatch')
    if subprocess.check_output(['git','--no-pager','status','--porcelain=v1'],cwd=root,text=True): raise ValueError('dirty source')
    for row in p['source_files']:
        if sha(root/row['path'])!=row['sha256']: raise ValueError('source pin mismatch: '+row['path'])
    if sha(p['gui'])!=p['gui_sha256']: raise ValueError('GUI pin mismatch')
    return p

def approval(path, p, policy_hash, mode):
    if not path: raise ValueError('separate later authorization required')
    a=json.loads(Path(path).read_text())
    required={'serve-local':{'local-visual','native-audible'},'serve-paid':{'paid-muted'}}[mode]
    if a.get('approved') is not True or a.get('policy_sha256')!=policy_hash or a.get('launcher_sha256')!=p['launcher_sha256'] or not required.intersection(a.get('scopes',[])) or not a.get('consent_reference'): raise ValueError('approval mismatch')
    if datetime.datetime.fromisoformat(a['expires_utc'])<=datetime.datetime.now(datetime.timezone.utc): raise ValueError('approval expired')
    return a

@contextlib.contextmanager
def server(handler):
    s=http.server.ThreadingHTTPServer(('127.0.0.1',0),handler);s.daemon_threads=False;s.timeout=1
    t=threading.Thread(target=lambda:s.serve_forever(poll_interval=.05));t.start()
    try: yield s
    finally: s.shutdown();s.server_close();t.join()

def fake_peer(p, directory):
    # Reuse the exact checked-in Peer; exclude only its top-level CLI server footer.
    text=(Path(p['root'])/'gui/test/fixture.py').read_text();assert text.count('\nserver=')==1
    ns={};exec(compile(text.split('\nserver=')[0],'pinned-existing-fixture.py','exec'),ns)
    ns['sys']=types.SimpleNamespace(argv=['fixture.py','0',str(directory/'raw-fake-upstream.jsonl')])
    class Peer(ns['Peer']):
        delayed=False;redirect=False;lock=threading.Lock();bodies=[]
        def do_POST(self):
            with self.lock:
                first=not self.__class__.delayed;self.__class__.delayed=True
            raw=self.rfile.read(int(self.headers.get('Content-Length','0')));self.rfile=io.BytesIO(raw)
            self.__class__.bodies.append(hashlib.sha256(raw).hexdigest())
            if first: time.sleep(3.25)
            if self.__class__.redirect:
                self.send_response(302);self.send_header('Location','https://example.invalid/no-follow');self.end_headers();return
            super().do_POST()
    return Peer

class Forward:
    def __init__(self,p,policy_path,policy_hash,upstream=None,paid_approval=None):
        self.p=p;self.policy_path=policy_path;self.policy_hash=policy_hash;self.upstream=upstream;self.paid_approval=paid_approval
        self.duration=p['observation_seconds'] if upstream is None else p['local_session_seconds']
        self.lock=threading.Lock();self.slots=0;self.started=None;self.failed=False;self.active=set();self.events=[];self.key=None;self.key_reads=0;self.external_requests=0
    def abort(self):
        with self.lock: self.failed=True;active=list(self.active)
        for conn in active:
            try:
                if conn.sock: conn.sock.shutdown(socket.SHUT_RDWR)
                conn.close()
            except OSError: pass
    def handler(self):
        owner=self
        class Proxy(http.server.BaseHTTPRequestHandler):
            def log_message(self,*args): pass
            def do_POST(self):
                conn=None;status=0;slot=None;start=time.monotonic();sent=False
                try:
                    verify(owner.policy_path,owner.policy_hash)
                    length=int(self.headers.get('Content-Length','0'))
                    if not 0<length<=4*1024*1024: raise ValueError('invalid body length')
                    raw=self.rfile.read(length);body=json.loads(raw)
                    if self.path!='/v1/messages': raise ValueError('wrong upstream path')
                    if body.get('model')!='claude-sonnet-4-6': raise ValueError('wrong model')
                    if type(body.get('max_tokens')) is not int or body['max_tokens']!=4096: raise ValueError('wrong native token maximum')
                    if body.get('stream') is not True: raise ValueError('streaming required')
                    with owner.lock:
                        now=time.monotonic()
                        if owner.failed: raise ValueError('terminal failure already recorded')
                        if owner.slots>=owner.p['upstream_request_cap']: raise ValueError('request cap reached')
                        if owner.started is None: owner.started=now
                        remaining=owner.duration-(now-owner.started)
                        if remaining<=0: raise ValueError('observation cap reached')
                        owner.slots+=1;slot=owner.slots
                    # Slots are irrevocably reserved before key access or outbound connection.
                    headers={'Content-Type':'application/json','anthropic-version':'2023-06-01','x-api-key':'owned-loopback-dummy'}
                    if owner.upstream is None:
                        approval(owner.paid_approval,owner.p,owner.policy_hash,'serve-paid')
                        with owner.lock:
                            if owner.key is None:
                                data=json.loads((Path.home()/'.cr/settings.json').read_text())
                                owner.key=data.get('directClaudeAPIKey');owner.key_reads+=1
                                if not isinstance(owner.key,str) or not owner.key: raise ValueError('approved credential field unavailable')
                        headers['x-api-key']=owner.key
                        conn=http.client.HTTPSConnection('api.anthropic.com',timeout=remaining);owner.external_requests+=1
                    else:
                        host,port=owner.upstream
                        if host!='127.0.0.1' or port==8088: raise ValueError('non-owned fake destination')
                        conn=http.client.HTTPConnection(host,port,timeout=remaining)
                    with owner.lock: owner.active.add(conn)
                    conn.request('POST','/v1/messages',body=raw,headers=headers);response=conn.getresponse();status=response.status
                    if not 200<=status<300: raise ValueError('upstream failed or redirected; no retry')
                    self.send_response(status);self.send_header('Content-Type',response.getheader('Content-Type','application/octet-stream'));self.send_header('Connection','close');self.end_headers();sent=True
                    while True:
                        if time.monotonic()-owner.started>=owner.duration: raise ValueError('observation cap reached')
                        chunk=response.read1(65536)
                        if not chunk: break
                        self.wfile.write(chunk);self.wfile.flush()
                    owner.events.append({'slot':slot,'status':status,'body_sha256':hashlib.sha256(raw).hexdigest(),'model':body['model'],'max_tokens':body['max_tokens'],'elapsed_seconds':round(time.monotonic()-start,3),'outcome':'forwarded'})
                except Exception as exc:
                    with owner.lock: owner.failed=True
                    # Never serialize exception text or response payloads: they may contain upstream/secret data.
                    owner.events.append({'slot':slot,'status':status,'elapsed_seconds':round(time.monotonic()-start,3),'outcome':'rejected','error_class':type(exc).__name__})
                    if not sent:
                        try:self.send_response(502);self.end_headers();self.wfile.write(b'Guard rejected; see safe receipt. No retry.')
                        except (BrokenPipeError,ConnectionResetError):pass
                finally:
                    if conn:
                        with owner.lock: owner.active.discard(conn)
                        conn.close()
        return Proxy

class GUI:
    def __init__(self,p,proxy_port,directory,number):
        verify(p['_policy_path'],p['_policy_hash'])
        self.directory=directory;self.path=directory/f'gui-{number}.log';self.output=self.path.open('w');self.ready=threading.Event();self.port=None
        args=[p['gui'],'-port','0','-workspace',str(directory/'workspace'),'-vendor','anthropic','-model','claude-sonnet-4-6','-base-url',f'http://127.0.0.1:{proxy_port}']
        env=clean_env(directory)
        self.process=subprocess.Popen(args,cwd=directory,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True,bufsize=1)
        self.ownership={'pid':self.process.pid,'argv':args,'cwd':str(directory),'environment_names':sorted(env),'dummy_credential_only':True,'profile':str(directory/'browser-profile'),'proxy_listener':['127.0.0.1',proxy_port]}
        save(directory/f'ownership-{number}.json',self.ownership) # before readiness/any future attach
        def drain():
            for line in self.process.stdout:
                self.output.write(line);self.output.flush();m=re.search(r'GUI http://127\.0\.0\.1:(\d+)',line)
                if m:self.port=int(m.group(1));self.ready.set()
        self.reader=threading.Thread(target=drain);self.reader.start()
        if not self.ready.wait(20):self.stop();raise RuntimeError('owned GUI startup failed')
        assert self.port!=8088
        asset_receipts=[]
        for row in p['assets']:
            path='/' if row['name']=='index.html' else '/'+row['name']
            c=http.client.HTTPConnection('127.0.0.1',self.port,timeout=20);c.request('GET',path);r=c.getresponse();data=r.read();c.close()
            asset_receipts.append({'asset':row['name'],'path':path,'status':r.status,'actual_sha256':hashlib.sha256(data).hexdigest(),'expected_sha256':row['sha256']});save(directory/f'assets-{number}.json',asset_receipts)
            if r.status!=200 or hashlib.sha256(data).hexdigest()!=row['sha256']:self.stop();raise ValueError('served asset mismatch: '+row['name'])
        self.ownership['gui_listener']=['127.0.0.1',self.port];self.ownership['served_assets_verified']=True;save(directory/f'ownership-{number}.json',self.ownership)
    def stop(self):
        if self.process.poll() is None:
            self.process.terminate()
            try:self.process.wait(timeout=20)
            except subprocess.TimeoutExpired:self.process.kill();self.process.wait(timeout=20)
        self.reader.join(timeout=5);assert not self.reader.is_alive();self.output.close()
        if not self.process.stdin.closed:self.process.stdin.close()
        self.ownership['reader_finished']=not self.reader.is_alive();self.ownership['exit_code']=self.process.returncode;save(self.directory/f'exit-{self.process.pid}.json',self.ownership)

def setup(path):
    path.mkdir(parents=True,exist_ok=False);(path/'cr').mkdir();(path/'workspace').mkdir();(path/'browser-profile').mkdir()
    (path/'workspace/exercise.txt').write_text('Chapter 8 boundary check.\n')
    save(path/'cr/settings.json',{'theme':'dark','tts_enabled':False,'tts_speed':1.75,'font_size':16,'max_tool_rounds':1})

def main():
    ap=argparse.ArgumentParser();ap.add_argument('mode',choices=['serve-local','serve-paid']);ap.add_argument('--policy',required=True);ap.add_argument('--policy-sha256',required=True);ap.add_argument('--approval',required=True);ap.add_argument('--run-dir',required=True);a=ap.parse_args()
    p=verify(a.policy,a.policy_sha256);grant=approval(a.approval,p,a.policy_sha256,a.mode)
    p.update(_policy_path=a.policy,_policy_hash=a.policy_sha256);directory=Path(a.run_dir).resolve();setup(directory)
    stop=threading.Event();restart=threading.Event();signal.signal(signal.SIGTERM,lambda *_:stop.set());signal.signal(signal.SIGINT,lambda *_:stop.set());signal.signal(signal.SIGUSR1,lambda *_:restart.set())
    save(directory/'controller.json',{'pid':os.getpid(),'argv':sys.argv,'mode':a.mode,'approval_sha256':sha(a.approval),'policy_sha256':a.policy_sha256,'launcher_sha256':sha(__file__),'browser_launched':False,'profile_reserved':str(directory/'browser-profile'),'audio_mode':'muted unless separately native-audible approved','approved_scopes':grant['scopes']})
    gui=None;forward=None;began=time.monotonic();number=1
    try:
        with contextlib.ExitStack() as stack:
            peer=stack.enter_context(server(fake_peer(p,directory))) if a.mode=='serve-local' else None
            forward=Forward(p,a.policy,a.policy_sha256,('127.0.0.1',peer.server_port) if peer else None,a.approval)
            proxy=stack.enter_context(server(forward.handler()))
            stack.callback(forward.abort) # abort sockets BEFORE server_close waits for handler threads, including exceptional exits
            gui=GUI(p,proxy.server_port,directory,number)
            print(json.dumps({'ready':True,'ownership':gui.ownership,'controller_pid':os.getpid()}),flush=True)
            while not stop.wait(.1):
                elapsed=time.monotonic()-began
                if forward.failed or gui.process.poll() is not None:break
                if forward.started and time.monotonic()-forward.started>=forward.duration:break
                if elapsed>=(p['local_session_seconds'] if a.mode=='serve-local' else 150):break
                if restart.is_set():
                    if a.mode!='serve-local' or number>=2:raise ValueError('restart not permitted')
                    restart.clear();gui.stop();number+=1;gui=GUI(p,proxy.server_port,directory,number);print(json.dumps({'restarted':True,'ownership':gui.ownership}),flush=True)
            forward.abort();gui.stop();gui=None
    finally:
        if forward:forward.abort()
        if gui:gui.stop()
        save(directory/'completion.json',{'elapsed_seconds':round(time.monotonic()-began,3),'browser_launched':False,'upstream_events':forward.events if forward else [],'key_reads':forward.key_reads if forward else 0,'external_requests':forward.external_requests if forward else 0,'processes_stopped':True})
if __name__=='__main__':main()
