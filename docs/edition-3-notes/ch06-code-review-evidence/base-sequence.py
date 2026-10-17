from pathlib import Path
import json,subprocess,sys
E=Path(__file__).resolve().parent;R=E.parents[2];S=Path((E/'scratch-path.txt').read_text().strip())/'source';runner=E/'run.py';ST=Path('/Users/bill/projects/ensemble-edition-3-student');G='/tmp/cadence-ch06-grader';B='/tmp/e3-ch06-frozen-core'
steps=[('root-race',S,['go','test','-race','-count=1','./...']),('root-vet',S,['go','vet','./...'])]
for d in ['ch04-consumer','ch06','ch07']:
 steps += [(d+'-race',S/'examples'/d,['go','test','-race','-count=1','./...']),(d+'-vet',S/'examples'/d,['go','vet','./...'])]
steps += [('grade-ch06',R,[G,'-ch','6','-json','-core',str(ST),'-consumer',str(ST/'examples/ch04-consumer'),'-workflow',str(ST/'examples/ch06'),'-stream-exercise',str(ST/'examples/ch07'),'-stream-model','claude-sonnet-4-6','-media-models',str(R/'docs/edition-3-notes/ch06-student-evidence/media-models.json'),B]),('grade-ch01',R,[G,'-ch','1','-json',B])]
results=[]
for label,cwd,argv in steps:
 p=subprocess.run([sys.executable,str(runner),label,str(cwd)]+argv);results.append({'label':label,'exit':p.returncode})
 if label.startswith('grade-') and p.returncode==0:
  report=json.loads((E/(label+'.stdout')).read_text());assert report['passed'] is True
(E/'base-sequence-results.json').write_text(json.dumps(results,indent=2)+'\n')
sys.exit(int(any(r['exit']!=0 for r in results)))
