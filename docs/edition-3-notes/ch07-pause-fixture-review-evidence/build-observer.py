#!/opt/homebrew/bin/python3
import os,sys,subprocess,pathlib,json,hashlib,shutil,time
w=pathlib.Path('/tmp/cadence-ch07-pause-fixture-review-independent');args=sys.argv[1:];cwd=pathlib.Path.cwd();p=subprocess.run(['/usr/local/go/bin/go']+args);rec={'argv':['/usr/local/go/bin/go']+args,'cwd':str(cwd),'exit':p.returncode,'phase':os.environ.get('REVIEW_PHASE'),'time':time.time()}
if p.returncode==0 and args and args[0]=='build' and '-o' in args and cwd.resolve()==(w/'source/gui').resolve():
 out=pathlib.Path(args[args.index('-o')+1]);dest=w/'bin'/('gui-'+os.environ['REVIEW_PHASE']);shutil.copy2(out,dest);rec.update({'captured':str(dest),'sha256':hashlib.sha256(dest.read_bytes()).hexdigest(),'bytes':dest.stat().st_size})
with (w/'go-invocations.jsonl').open('a') as f:f.write(json.dumps(rec)+'\n')
sys.exit(p.returncode)
