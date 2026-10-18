#!/usr/bin/python3
import subprocess,sys,os,json,hashlib,shutil
from pathlib import Path
w=Path('/tmp/cadence-ch07-r4-pause-independent');a=sys.argv[1:];real=(w/'real-go.txt').read_text().strip();p=subprocess.run([real]+a)
row={'argv':[real]+a,'cwd':os.getcwd(),'exit':p.returncode}
if p.returncode==0 and 'build' in a and '-o' in a and Path.cwd().resolve()==(w/'source/gui').resolve():
 b=Path(a[a.index('-o')+1]);shutil.copy2(b,w/'built-gui');row['gui_sha256']=hashlib.sha256(b.read_bytes()).hexdigest()
with (w/'go-invocations.jsonl').open('a') as f:f.write(json.dumps(row)+'\n')
sys.exit(p.returncode)
