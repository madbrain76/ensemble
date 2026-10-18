import sys;sys.dont_write_bytecode=True
from common import *
results=json.loads((E/'mutation-partial-results.json').read_text());assert len(results)==2 and results[0]['passed'] and results[1]['failed_checks']==['settings-on-connect']
G=T/'real-gui';p=G/'gui.go';assert p.read_bytes()==(E/'sources/gui-original.go.txt').read_bytes();driver=T/'bin/candidate-settings-driver';name='candidate-check-restored';rc=run(name,[str(driver),'-mode','settings','-gui',str(G),'-gui-model','claude-sonnet-4-6'],T/'work',90);r=json.loads((E/(name+'.stdout')).read_text());results.append({'name':name,'source_sha256':sha(p),'driver_sha256':sha(driver),'exit':rc,'score':r['score'],'max_score':r['max_score'],'passed':r['passed'],'failed_checks':[c['id'] for c in r['checks'] if not c['passed']]});assert rc==0 and r['passed'] and r['score']==80
(E/'mutation-results.json').write_text(json.dumps(results,indent=2)+'\n');print(json.dumps(results,indent=2))
