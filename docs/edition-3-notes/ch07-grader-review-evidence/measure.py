from run import W,run
import json
C='/tmp/cadence-ch07-grader';P='/tmp/cadence-ch07-entry-auditor/positive';B='/tmp/cadence-ch07-entry-auditor/bin/positive-core'
base={'-ch':'7','-core':P,'-consumer':P+'/ch05','-workflow':P+'/ch06','-stream-exercise':P+'/ch07','-stream-model':'claude-sonnet-4-6','-gui':P+'/gui','-gui-model':'claude-sonnet-4-6','-consumer-case':'/tmp/cadence-ch07-entry-auditor/consumer-case.json','-media-models':'/tmp/cadence-ch07-entry-auditor/media-models.json'}
def invoke(name,opts,expected=None):
 argv=[C,'-json',*[x for pair in opts.items() for x in pair],B]
 m,r,o,e=run(name,argv);results.append({'name':name,'exit':m['exit'],'seconds':m['seconds'],'report':r});(W/'results.json').write_text(json.dumps(results,indent=2)+'\n');print(o,flush=True);print(e,flush=True)
 assert isinstance(r,dict),name
 if expected is not None:assert m['exit']==expected,(name,m)
results=[]
invoke('positive',base,0)
x=base.copy();x['-gui']=str(W/'alternate-gui');invoke('wrong-gui-root',x,1)
x=base.copy();x['-gui-model']='adapter-review-unlisted-gui-model';invoke('nondefault-gui-model',x)
x=base.copy();x['-stream-model']='adapter-review-unlisted-stream-model';invoke('wrong-stream-model',x,1)
x=base.copy();x['-consumer-case']=str(W/'wrong-case.json');x['-media-models']=str(W/'wrong-media.json');invoke('custom-case-media',x,1)
for ch in ['4','5','6']:
 x={k:v for k,v in base.items() if k not in ['-gui','-gui-model']};x['-ch']=ch
 if ch in ['4','5']:x.pop('-stream-exercise');x.pop('-stream-model')
 if ch=='4':x.pop('-workflow')
 invoke('legacy-ch'+ch,x,0)
print('Finite real-grader sequence complete.',flush=True)
