from pathlib import Path
import subprocess,sys,json
E=Path(__file__).resolve().parent;R=E.parents[2];run=str(E/'run.py');C='/tmp/cadence-ch06-grader';H='/tmp/ensemble-ch06-entry-7_7zbtqs';core=H+'/solutions/ch06';B=H+'/positive';model='claude-sonnet-4-6'
base=['-ch','6','-core',core,'-consumer',core+'/ch05','-workflow',core+'/ch06','-stream-exercise',core+'/ch07','-stream-model',model]
seq=[('positive',base+[B],0),('wrong-stream-model',base+['-stream-model','adapter-review-unlisted-stream-model',B],1),('wrong-stream-exercise',base+['-stream-exercise',str(Path((E/'scratch-path.txt').read_text().strip())/'module'),B],1),('custom-case-media',base+['-consumer-case',str(E/'custom-consumer-case.json'),'-media-models',str(E/'custom-media-models.json'),B],1),('legacy-ch4',['-ch','4','-core',core,'-consumer',core+'/ch05',B],0),('legacy-ch5',['-ch','5','-core',core,'-consumer',core+'/ch05','-workflow',core+'/ch06',B],0),('legacy-ch1',['-ch','1','/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble'],0)]
for label,args,want in seq:
 p=subprocess.run([sys.executable,run,label,str(R),C]+args)
 if p.returncode!=want:raise SystemExit(f'{label}: unexpected exit {p.returncode}; stop without changing fixture')
 r=json.loads((E/(label+'.stdout')).read_text());print(label,'score',r['score'],'passed',r['passed'],flush=True)
