from pathlib import Path
import subprocess,json,sys
E=Path(__file__).resolve().parent;R=E.parents[2];S=Path((E/'scratch-path.txt').read_text().strip());ST=Path('/Users/bill/projects/ensemble-edition-3-student');core='/tmp/e3-ch06-correction-r1/core';grader='/tmp/cadence-ch06-grader'
seq=[('root-race',S/'source',['go','test','-race','-count=1','./...']),('root-vet',S/'source',['go','vet','./...'])]
for mod in ['ch04-consumer','ch06','ch07']:
 seq += [(mod+'-race',S/'source/examples'/mod,['go','test','-race','-count=1','./...']),(mod+'-vet',S/'source/examples'/mod,['go','vet','./...'])]
seq += [('grade-ch06',R,[grader,'-ch','6','-json','-core',str(ST),'-consumer',str(ST/'examples/ch04-consumer'),'-workflow',str(ST/'examples/ch06'),'-stream-exercise',str(ST/'examples/ch07'),'-stream-model','claude-sonnet-4-6','-media-models',str(R/'docs/edition-3-notes/ch06-student-r1-evidence/media-models.json'),core]),('grade-ch01',R,[grader,'-ch','1','-json',core])]
(E/'base-planned.json').write_text(json.dumps([{'label':l,'cwd':str(c),'argv':a} for l,c,a in seq],indent=2)+'\n');results=[]
for l,c,a in seq:
 n=subprocess.call([sys.executable,str(E/'run.py'),l,str(c),*a]);results.append({'label':l,'exit':n});(E/'base-results.json').write_text(json.dumps(results,indent=2)+'\n')
 if n:raise SystemExit(n)
 if l.startswith('grade-'):
  d=json.loads((E/(l+'.stdout')).read_text());assert d['passed'] and d['score']==100,d
