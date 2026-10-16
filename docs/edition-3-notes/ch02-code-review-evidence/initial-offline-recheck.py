#!/usr/bin/env python3
"""Reviewer-only evidence verification. No live-provider path or credentials."""
import hashlib,http.server,json,os,pathlib,subprocess,sys,threading,time
bundle=pathlib.Path(sys.argv[1]); binary=bundle/'student-cli'; out=pathlib.Path(sys.argv[2])
def sha(b):return hashlib.sha256(b).hexdigest()
m=json.loads((bundle/'evidence-manifest.json').read_text());checked=[]
for name,want in m['files'].items():
 b=(bundle/name).read_bytes();assert sha(b)==want['sha256'] and len(b)==want['bytes'],name;checked.append(name)
assert sha(binary.read_bytes())==m['binary_sha256']
(out/'retained-manifest-verification.json').write_text(json.dumps({'checked_file_count':len(checked),'all_hashes_and_sizes_match':True,'manifest_sha256':sha((bundle/'evidence-manifest.json').read_bytes()),'binary_sha256':sha(binary.read_bytes()),'checked_files':checked},indent=2)+'\n')
class Sentinel(http.server.BaseHTTPRequestHandler):
 posts=0
 def log_message(self,*args):pass
 def do_POST(self):
  type(self).posts+=1;self.send_response(500);self.end_headers()
server=http.server.HTTPServer(('127.0.0.1',0),Sentinel);thread=threading.Thread(target=server.serve_forever,daemon=True);thread.start()
env={k:os.environ[k] for k in ('PATH','HOME','TMPDIR') if k in os.environ};env['LLM_BASE_URL']='http://127.0.0.1:'+str(server.server_port)
summary=json.loads((bundle/'live-summary.json').read_text());rows=[]
def run(args):
 start=time.monotonic();p=subprocess.run([str(binary),*map(str,args)],env=env,capture_output=True,timeout=15);assert p.returncode==0 and not p.stderr,(args,p.returncode,len(p.stderr));return p.stdout,round(time.monotonic()-start,4)
try:
 for source in summary['vendors']:
  log=bundle/(source['vendor']+'-live.jsonl');raw=log.read_bytes();dump,elapsed=run(['dump','--log',log]);assert dump==raw==(bundle/(source['vendor']+'-dump.jsonl')).read_bytes()
  rows.append({'kind':'dump','source':source['vendor'],'exit_code':0,'byte_identical':True,'sha256':sha(dump),'elapsed_seconds':elapsed})
  for dest in summary['vendors']:
   target=dest['vendor'];model=dest['responses'][0]['from']['model'];args=['render','--log',log,'--vendor',target,'--model',model]
   first,elapsed1=run(args);second,elapsed2=run(args);assert first==second
   saved1=(bundle/(source['vendor']+'-to-'+target+'-1.json')).read_bytes();saved2=(bundle/(source['vendor']+'-to-'+target+'-2.json')).read_bytes();assert first==saved1==saved2
   rows.append({'kind':'render-pair','source':source['vendor'],'target':target,'model':model,'exit_codes':[0,0],'repeat_equal':True,'retained_pair_equal':True,'sha256':sha(first),'elapsed_seconds':[elapsed1,elapsed2]})
 assert Sentinel.posts==0
finally:server.shutdown();server.server_close();thread.join()
result={'additional_live_calls':0,'key_environment':'allowlist PATH/HOME/TMPDIR only; LLM_BASE_URL=loopback sentinel; no API keys','cli_sha256':sha(binary.read_bytes()),'sentinel_posts':Sentinel.posts,'records':rows,'original_live_launcher_exit':(bundle/'live-launcher.exit').read_text().strip(),'retained_offline_verifier_exit':(bundle/'capture-verification.exit').read_text().strip()}
(out/'independent-offline-recheck.json').write_text(json.dumps(result,indent=2)+'\n')
print('Verified',len(checked),'retained artifact hashes/sizes; 3 new dumps and 18 new keyless render processes match retained bytes; sentinel POSTs=0; additional live calls=0.')
