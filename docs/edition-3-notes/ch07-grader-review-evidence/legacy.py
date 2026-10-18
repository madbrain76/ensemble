from run import W,run
import json
C='/tmp/cadence-ch07-grader';P='/tmp/cadence-ch07-entry-auditor/positive';B='/tmp/cadence-ch07-entry-auditor/bin/positive-core';rows=[]
for ch in ['4','5','6']:
 opts={'-ch':ch,'-core':P,'-consumer':P+'/ch05','-consumer-case':'/tmp/cadence-ch07-entry-auditor/consumer-case.json'}
 if ch in ['5','6']:opts.update({'-workflow':P+'/ch06','-media-models':'/tmp/cadence-ch07-entry-auditor/media-models.json'})
 if ch=='6':opts.update({'-stream-exercise':P+'/ch07','-stream-model':'claude-sonnet-4-6'})
 m,r,o,e=run('legacy-ch'+ch+'-corrected',[C,'-json',*[x for p in opts.items() for x in p],B]);print(o,flush=True);print(e,flush=True);rows.append({'name':m['name'],'exit':m['exit'],'seconds':m['seconds'],'report':r});(W/'legacy-results.json').write_text(json.dumps(rows,indent=2)+'\n');assert m['exit']==0 and r['passed'] is True
print('Corrected legacy sequence complete.')
