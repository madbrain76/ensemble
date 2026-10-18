#!/usr/bin/python3
import os,sys,subprocess,json,hashlib,time
from pathlib import Path
a=['/usr/local/go/bin/go']+sys.argv[1:];start=time.time();rc=subprocess.call(a)
if len(a)>2 and a[1]=='build':
 r={'stage':os.environ.get('AUDIT_STAGE'),'argv':a,'cwd':os.getcwd(),'exit':rc,'seconds':round(time.time()-start,3)}
 if '-o' in a:
  p=Path(a[a.index('-o')+1]);r['output']=str(p)
  if rc==0 and p.is_file():r.update(bytes=p.stat().st_size,sha256=hashlib.sha256(p.read_bytes()).hexdigest())
 with open(os.environ['AUDIT_BUILD_LOG'],'a') as f:f.write(json.dumps(r)+'\n')
sys.exit(rc)
