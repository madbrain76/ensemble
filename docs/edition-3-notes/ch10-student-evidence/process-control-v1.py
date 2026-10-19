import os,sys,json,time,threading,subprocess,socket,struct,base64,hashlib,re
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
from pathlib import Path
root=Path('/Users/bill/projects/ensemble-edition-3-student');work=root/'.ch10-local/process-control-v1';work.mkdir(exist_ok=False)
requests=[];children=[];sockets=[];results=[];cli_runs=0
class Handler(BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_POST(self):
  body=json.loads(self.rfile.read(int(self.headers['Content-Length'])));requests.append(body)
  n=len(requests);time.sleep(.03)
  b=json.dumps({'id':str(n),'type':'message','role':'assistant','model':'claude-sonnet-4-6','stop_reason':'end_turn','content':[{'type':'text','text':'fixture reply '+str(n)}],'usage':{'input_tokens':9,'output_tokens':3}}).encode()
  self.send_response(200);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(b)));self.end_headers();self.wfile.write(b)
server=ThreadingHTTPServer(('127.0.0.1',0),Handler);thread=threading.Thread(target=server.serve_forever);thread.start();origin='http://127.0.0.1:'+str(server.server_port)
env=dict(os.environ,LLM_VENDOR='anthropic',LLM_MODEL='claude-sonnet-4-6',LLM_SYSTEM='current process instruction',ANTHROPIC_BASE_URL=origin,LLM_API_KEY='offline-fixture-not-a-secret',EN_DISABLE_STREAMING='1')
assert 'LLM_BASE_URL' not in env

def record(name,**fields):results.append(dict(name=name,passed=True,**fields))
def digest(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def archive(p):
 s=json.loads(p.read_text());assert set(s)=={'config','context','asof','log'};assert 'offline-fixture-not-a-secret' not in p.read_text();return s

def cli(directory,text,args=(),success=True):
 global cli_runs
 cli_runs+=1
 directory.mkdir(exist_ok=True);e=dict(env,CH02_LOG=str(directory/'history.jsonl'))
 p=subprocess.run([str(root/'.ch10-local/bin/frozen-cli'),*args],input=text.encode(),stdout=subprocess.PIPE,stderr=subprocess.PIPE,cwd=directory,env=e,timeout=20)
 assert (p.returncode==0)==success,(p.returncode,p.stderr.decode())
 (directory/('cli-'+str(cli_runs)+'.stdout')).write_bytes(p.stdout);(directory/('cli-'+str(cli_runs)+'.stderr')).write_bytes(p.stderr);return p

def start_gui(directory,path):
 directory.mkdir(exist_ok=True);out=(directory/'stdout').open('wb');err=(directory/'stderr').open('wb')
 p=subprocess.Popen([str(root/'.ch10-local/bin/frozen-gui'),'--port','0','--no-stream','--save',str(path),'--base-url',origin,'--model','claude-sonnet-4-6'],cwd=directory,env=env,stdin=subprocess.PIPE,stdout=out,stderr=err);children.append(p)
 deadline=time.monotonic()+10
 while time.monotonic()<deadline:
  text=(directory/'stderr').read_text();m=re.search(r'GUI http://127.0.0.1:(\d+)',text)
  if m:return p,int(m.group(1))
  assert p.poll() is None,text
  time.sleep(.02)
 raise AssertionError('GUI startup timeout')

def send(sock,obj):
 b=json.dumps(obj).encode();mask=b'abcd';head=bytes([0x81,0x80|len(b)]) if len(b)<126 else bytes([0x81,0xfe])+struct.pack('!H',len(b));sock.sendall(head+mask+bytes(x^mask[i%4] for i,x in enumerate(b)))
def exact(sock,n):
 b=b''
 while len(b)<n:
  chunk=sock.recv(n-len(b));assert chunk,'unexpected socket EOF';b+=chunk
 return b

def recv(sock):
 a,n=exact(sock,2);n=n&127
 if n==126:n=struct.unpack('!H',exact(sock,2))[0]
 elif n==127:n=struct.unpack('!Q',exact(sock,8))[0]
 b=exact(sock,n);assert a&15==1,(a,b)
 return json.loads(b)
def ws(port):
 sock=socket.create_connection(('127.0.0.1',port),timeout=10);sockets.append(sock)
 key=base64.b64encode(b'0123456789abcdef').decode();sock.sendall((f'GET /ws HTTP/1.1\r\nHost: 127.0.0.1:{port}\r\nOrigin: http://127.0.0.1:{port}\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: {key}\r\nSec-WebSocket-Version: 13\r\n\r\n').encode())
 header=b''
 while not header.endswith(b'\r\n\r\n'):header+=exact(sock,1)
 assert header.startswith(b'HTTP/1.1 101'),header
 send(sock,{'type':'subscribe'});messages=[]
 while True:
  m=recv(sock);messages.append(m)
  if m.get('type')=='replay_end':return sock,messages

def final(sock,expected):
 messages=[]
 while True:
  m=recv(sock);messages.append(m)
  if m.get('type')=='part_final' and m.get('part',{}).get('text')==expected:return messages

def finish(p):
 p.stdin.close();assert p.wait(timeout=20)==0
try:
 # Actual CLI startup / EOF / continuation and current config.
 cli_dir=work/'cli';cli(cli_dir,'{"kind":"prompt","text":"first CLI"}\n');path=cli_dir/'save.json';a=archive(path);before=a['asof']
 env['LLM_SYSTEM']='changed runtime instruction';cli(cli_dir,'{"kind":"prompt","text":"second CLI"}\n');a=archive(path)
 assert a['asof']>before and 'first CLI' in json.dumps(requests[-1]) and 'fixture reply 1' in json.dumps(requests[-1]);assert 'changed runtime instruction' in json.dumps(requests[-1]);assert a['config']['system_prompt']=='changed runtime instruction'
 record('CLI default restart/current runtime',asof=a['asof'],archive_sha256=digest(path))
 alt=cli_dir/'selected.json';cli(cli_dir,'{"kind":"ephemeral","text":"pending CLI"}\n',['--save',str(alt)]);assert len(archive(alt)['context']['ephemera'])==1
 bad=cli_dir/'bad.json';bad.write_text('{"config":{}}');prior=bad.read_bytes();cli(cli_dir,'',['--save',str(bad)],False);assert bad.read_bytes()==prior
 record('CLI explicit path/pending EOF/malformed no overwrite')
 # Real GUI stdin + WebSocket share exactly one retained lineage.
 save=work/'gui-save.json';p,port=start_gui(work/'gui-first',save);sock,replay=ws(port)
 assert not [m for m in replay if m.get('type')=='part_final']
 count=len(requests);p.stdin.write(b'{"kind":"ephemeral","text":"one-use GUI context"}\n{"kind":"prompt","text":"stdin GUI first"}\n');p.stdin.flush();first=final(sock,'fixture reply '+str(count+1));assert 'one-use GUI context' in json.dumps(requests[-1])
 count=len(requests);send(sock,{'type':'prompt','text':'websocket GUI second'});second=final(sock,'fixture reply '+str(count+1));assert 'stdin GUI first' in json.dumps(requests[-1]);assert 'one-use GUI context' not in json.dumps(requests[-1])
 # EOF immediately after another admitted prompt must wait for that response.
 count=len(requests);p.stdin.write(b'{"kind":"prompt","text":"EOF drains this GUI prompt"}\n');finish(p);a=archive(save);assert len(requests)==count+1 and 'EOF drains this GUI prompt' in json.dumps(a['context']);sock.close()
 ids=[m['id'] for m in first+second if m.get('type')=='part_final'];assert len(ids)==len(set(map(str,ids)))
 record('GUI stdin+WebSocket shared Agent/one-use ephemera/EOF drain',asof=a['asof'],archive_sha256=digest(save))
 # Restart replay and reconnect identity: presentation stays stable, future IDs disjoint.
 p,port=start_gui(work/'gui-second',save);sock,replay=ws(port);old=[m for m in replay if m.get('type')=='part_final'];assert old
 sock2,replay2=ws(port);same=[m for m in replay2 if m.get('type')=='part_final'];assert [(m['id'],m.get('part')) for m in old]==[(m['id'],m.get('part')) for m in same]
 count=len(requests);send(sock,{'type':'prompt','text':'after GUI restart'});live=final(sock,'fixture reply '+str(count+1));oldids={str(m['id']) for m in old};newids={str(m['id']) for m in live if m.get('type')=='part_final'};assert not oldids&newids
 p.stdin.write(b'{"kind":"ephemeral","text":"pending across GUI EOF"}\n');finish(p);a=archive(save);assert len(a['context']['ephemera'])==1;assert 'websocket GUI second' in json.dumps(requests[-1]);sock.close();sock2.close()
 record('GUI reconnect/restart identity/continued request/pending EOF',archived_part_ids=[m['id'] for m in old],new_part_ids=sorted(newids),archive_sha256=digest(save))
 print(json.dumps({'passed':True,'controls':results,'fake_http_requests':len(requests),'browser_observed':False,'live_provider_used':False,'all_children_exited':all(p.poll() is not None for p in children)},indent=2))
finally:
 for sock in sockets:
  try:sock.close()
  except OSError:pass
 for p in children:
  if p.poll() is None:p.terminate();p.wait(timeout=10)
 server.shutdown();server.server_close();thread.join(timeout=10)
