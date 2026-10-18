import socket,subprocess,pathlib,time,json,base64,hashlib,struct,sys
base=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');ev=pathlib.Path(__file__).parent;label=sys.argv[1];binary=base/'.ch07-local'/sys.argv[2];run=base/'.ch07-local'/label;run.mkdir(parents=True,exist_ok=True)
s=socket.socket();s.bind(('127.0.0.1',0));port=s.getsockname()[1];s.close();env={'PATH':'/usr/bin:/bin','HOME':str(run),'LLM_API_KEY':'fixture'}
p=subprocess.Popen([str(binary),'--port',str(port),'--vendor','anthropic','--model','claude-sonnet-4-6','--base-url','http://127.0.0.1:1','--workspace',str(run)],cwd=run,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
def connect():
 for _ in range(100):
  try:c=socket.create_connection(('127.0.0.1',port),.1);break
  except OSError:time.sleep(.02)
 c.settimeout(2);c.sendall(f'GET /ws HTTP/1.1\r\nHost: 127.0.0.1:{port}\r\nConnection: Upgrade\r\nUpgrade: websocket\r\nSec-WebSocket-Version: 13\r\nSec-WebSocket-Key: MDEyMzQ1Njc4OWFiY2RlZg==\r\n\r\n'.encode());header=b''
 while not header.endswith(b'\r\n\r\n'):header+=c.recv(1)
 assert b'101' in header;return c
def send(c,op,data,fin=True):
 mask=b'1234';n=len(data);h=bytes([(128 if fin else 0)|op]);h+=bytes([128|n]) if n<126 else bytes([128|126])+struct.pack('!H',n);c.sendall(h+mask+bytes(v^mask[i%4] for i,v in enumerate(data)))
def recv(c):
 h=c.recv(2)
 if not h:return {'eof':True}
 n=h[1]&127
 if n==126:n=struct.unpack('!H',c.recv(2))[0]
 if n==127:n=struct.unpack('!Q',c.recv(8))[0]
 data=b''
 while len(data)<n:data+=c.recv(n-len(data))
 return {'opcode':h[0]&15,'data':data.decode()}
r={'binary':str(binary),'sha256':hashlib.sha256(binary.read_bytes()).hexdigest(),'port':port,'checks':[]}
try:
 c=connect();send(c,9,b'ping');r['checks'].append({'name':'ping','response':recv(c)});c.close()
 c=connect();send(c,1,b'{"type":',False);send(c,0,b'"subscribe"}');r['checks'].append({'name':'fragmented-subscribe','response':recv(c)});c.close()
finally:
 p.terminate();out,err=p.communicate(timeout=5);r['exit']=p.returncode;(ev/(label+'-stdout.log')).write_bytes(out);(ev/(label+'-stderr.log')).write_bytes(err);(ev/(label+'.json')).write_text(json.dumps(r,indent=2)+'\n');print(json.dumps(r,indent=2))
