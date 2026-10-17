#!/usr/bin/env python3
"""Finite Chapter 6 observations. --controls never reads credentials/network.
Live mode is forbidden without parent approval and exact source/binary freeze.
"""
import argparse,hashlib,http.client,http.server,json,os,pathlib,selectors,subprocess,tempfile,threading,time
OUT=pathlib.Path(__file__).resolve().parent
VENDORS=[('anthropic','claude-sonnet-4-6','api.anthropic.com','directClaudeAPIKey'),('openai','gpt-4.1-2025-04-14','api.openai.com','directOpenAIAPIKey'),('gemini','gemini-3.8-flash','generativelanguage.googleapis.com','directGeminiAPIKey')]
CORE='/tmp/e3-ch06-frozen-core';SHAPE='/tmp/e3-ch06-frozen-shape'
LIMIT=3;TOKENS=1024

def allowed(body,vendor,model,count):
 if count>=LIMIT:return False
 if vendor=='gemini':return body.get('generationConfig',{}).get('maxOutputTokens',0) in range(1,TOKENS+1)
 return body.get('model')==model and body.get('max_tokens',body.get('max_completion_tokens',0)) in range(1,TOKENS+1)

def controls():
 result={}
 for vendor,model,_,_ in VENDORS:
  body={'model':model,'max_tokens':TOKENS,'generationConfig':{'maxOutputTokens':TOKENS}}
  result[vendor]={'within_cap':allowed(body,vendor,model,0),'blocks_fourth':not allowed(body,vendor,model,3)}
  if vendor=='gemini':body['generationConfig']['maxOutputTokens']=TOKENS+1
  else:body['max_tokens']=TOKENS+1
  result[vendor]['blocks_output_over_cap']=not allowed(body,vendor,model,0)
 assert all(all(v.values()) for v in result.values())
 (OUT/'live-harness-controls.json').write_text(json.dumps(result,indent=2)+'\n')
 print('Local safety controls passed; no credential reads or provider requests.')

def run_case(vendor,model,host,key,mode,root):
 case=f'{vendor}-{mode}';work=root/case;work.mkdir();(work/'config.json').write_text('{"project":"stream-study","setting":7}\n')
 stats={'case':case,'request_count':0,'statuses':[],'request_ended_ms':[],'max_output_tokens':TOKENS,'request_limit':LIMIT,'output_chunks':[],'hints_sent':0}
 start=time.monotonic();elapsed=lambda:round((time.monotonic()-start)*1000,3)
 class Handler(http.server.BaseHTTPRequestHandler):
  def log_message(self,*args):pass
  def do_POST(self):
   try:
    data=self.rfile.read(int(self.headers.get('Content-Length','0')));body=json.loads(data)
    if not allowed(body,vendor,model,stats['request_count']) or (vendor=='gemini' and not self.path.startswith('/v1beta/models/'+model+':')):
     stats['cap_refusal']=True;self.send_response(429);self.end_headers();return
    stats['request_count']+=1
    headers={'Content-Type':'application/json'}
    if vendor=='anthropic':headers.update({'x-api-key':key,'anthropic-version':'2023-06-01'})
    elif vendor=='openai':headers['Authorization']='Bearer '+key
    else:headers['x-goog-api-key']=key
    connection=http.client.HTTPSConnection(host,timeout=60)
    # No automatic retry, redirect, model substitution, or provider discovery.
    connection.request('POST',self.path,data,headers);response=connection.getresponse();stats['statuses'].append(response.status)
    self.send_response(response.status);self.send_header('Content-Type',response.getheader('Content-Type','application/json'));self.end_headers()
    # Newline framing retains prompt progress and permits complete key redaction.
    while True:
     chunk=response.readline()
     if not chunk:break
     self.wfile.write(chunk.replace(key.encode(),b'[redacted]'));self.wfile.flush()
    stats['request_ended_ms'].append(elapsed());connection.close()
   except Exception as error:
    stats['transport_failure']=type(error).__name__
    try:self.send_response(502);self.end_headers()
    except Exception:pass
 server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);server.daemon_threads=False
 thread=threading.Thread(target=server.serve_forever);thread.start()
 binary=CORE if mode=='human-stream' else SHAPE
 env={k:v for k,v in os.environ.items() if not k.startswith(('LLM_','ANTHROPIC_','OPENAI_','GEMINI_','EN_','CH07_'))}
 env.update(LLM_VENDOR=vendor,LLM_MODEL=model,LLM_BASE_URL=f'http://127.0.0.1:{server.server_port}',LLM_API_KEY='local-harness-placeholder',LLM_TOOL_BATCH_LIMIT='2',LLM_MAX_TOKENS=str(TOKENS),CH07_NO_STREAM='1' if mode=='shape-plain' else '0')
 prompt='Read config.json with read_file, then write output.txt containing exactly STREAM_OK using write_file. Explain briefly before tools. Do not run shell commands. End with one short sentence.' if mode=='human-stream' else 'Use read_file to inspect config.json, then report the project name and setting in one short sentence. Explain before tools; do not write files or run shell commands.'
 command=[binary]+(['chat'] if mode=='human-stream' else [])
 proc=subprocess.Popen(command,cwd=work,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
 proc.stdin.write((prompt+'\n').encode());proc.stdin.flush()
 if mode!='human-stream':proc.stdin.close()
 selector=selectors.DefaultSelector();selector.register(proc.stdout,selectors.EVENT_READ,'stdout');selector.register(proc.stderr,selectors.EVENT_READ,'stderr');raw={'stdout':bytearray(),'stderr':bytearray()}
 stdin_closed=mode!='human-stream';first=None
 try:
  while selector.get_map():
   if elapsed()>90000:
    stats['timed_out']=True;proc.terminate()
    try:proc.wait(5)
    except subprocess.TimeoutExpired:proc.kill();proc.wait()
   if not stdin_closed and elapsed()>10000:proc.stdin.close();stdin_closed=True
   for selected,_ in selector.select(.1):
    data=os.read(selected.fileobj.fileno(),4096)
    if not data:selector.unregister(selected.fileobj);continue
    raw[selected.data].extend(data)
    if selected.data=='stdout':
     text=raw['stdout'].decode(errors='replace').removeprefix('> ')
     if text.strip() and first is None:first=elapsed()
     stats['output_chunks'].append({'ms':elapsed(),'bytes':len(data)})
     if mode=='human-stream' and text.strip() and not stdin_closed:
      try:proc.stdin.write(b'Keep the final answer concise.\n');proc.stdin.flush();stats['hints_sent']=1
      except BrokenPipeError:pass
      proc.stdin.close();stdin_closed=True
  stats['exit']=proc.wait(timeout=5)
 finally:
  if proc.poll() is None:proc.kill();proc.wait()
  if not stdin_closed:proc.stdin.close()
  selector.close();server.shutdown();server.server_close();thread.join()
 stats['elapsed_ms']=elapsed();stats['first_reply_stdout_ms']=first
 for stream,data in raw.items():(work/(stream+'.raw')).write_bytes(data)
 # Human stderr contains private reasoning/tool chunks. Retain it locally, export
 # only original known-safe accounting/error lines, never opaque frame/history logs.
 stdout=raw['stdout'].decode(errors='replace').replace(key,'[redacted]')
 stderr=raw['stderr'].decode(errors='replace').replace(key,'[redacted]')
 safe_stderr='\n'.join(line for line in stderr.splitlines() if line.startswith(('Tokens used:','error: anthropic HTTP ','error: openai HTTP ','error: gemini HTTP ','error: model turn canceled','error: tool batch limit')))
 (OUT/(case+'.stdout.txt')).write_text(stdout);(OUT/(case+'.stderr-safe.txt')).write_text(safe_stderr+'\n[Other stderr retained privately; may include reasoning/tool progress.]\n')
 effect=work/'output.txt';stats['expected_file_effect']=effect.exists() and effect.read_text()=='STREAM_OK'
 stats['private_directory']=str(work);stats['binary_sha256']=hashlib.sha256(pathlib.Path(binary).read_bytes()).hexdigest()
 if mode!='human-stream':
  try:stats['shape']=json.loads(stdout)
  except ValueError:stats['shape']='no JSON summary'
 (OUT/(case+'.json')).write_text(json.dumps(stats,indent=2)+'\n')
 return stats

def live():
 # Only approved settings fields, process memory only; never print values.
 settings=json.loads((pathlib.Path.home()/'.cr/settings.json').read_text())
 keys={v:settings.get(field,'') for v,_,_,field in VENDORS};del settings
 if not all(keys.values()):raise SystemExit('Required approved credential field absent; stopped without substitutions.')
 root=pathlib.Path(tempfile.mkdtemp(prefix='e3-ch06-live-private-'));root.chmod(0o700)
 results=[]
 for vendor,model,host,_ in VENDORS:
  for mode in ['human-stream','shape-stream','shape-plain']:
   result=run_case(vendor,model,host,keys[vendor],mode,root);results.append(result)
   print(result['case'], 'exit',result['exit'],'requests',result['request_count'],'statuses',result['statuses'],flush=True)
 (OUT/'live-summary.json').write_text(json.dumps(results,indent=2)+'\n')
if __name__=='__main__':
 parser=argparse.ArgumentParser();parser.add_argument('--controls',action='store_true');parser.add_argument('--approved-live',action='store_true');args=parser.parse_args()
 if args.controls:controls()
 elif args.approved_live:live()
 else:parser.error('explicit mode required')
