from pathlib import Path
import json,hashlib,importlib.util,copy,sys
p=Path(__file__).with_name('controller-v4.py');expected=Path(__file__).with_name('controller-v4.sha256').read_text().strip();assert hashlib.sha256(p.read_bytes()).hexdigest()==expected
spec=importlib.util.spec_from_file_location('owned_notice_controller',p);m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m)
e=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch08-operational-resize-continuation-evidence');raw=json.loads((e/'selected-actual-app-protocol.json').read_text())['records'];messages=[x['message'] for x in raw];dom=json.loads((e/'last-actual-dom-before-failure.json').read_text());base=m.classify_actual(messages,dom['notice'],dom['state'],dom['paused']);assert base['ok'] and base['complete'] and len(base['dispatch_by_call_id'])==2
rows=[{'case':'actual successful retained-job notice and completed pairs','passed':True,'facts':base}]
fixtures=[('typed error',{'type':'error','text':'controlled real-shaped error'}),('typed settings error',{'type':'settings_error','text':'controlled settings error'}),('nonempty turn error',{'type':'turn_ended','error':'controlled turn error','text':''})]
failed=copy.deepcopy(next(x for x in messages if x['type']=='tool_finished'));failed['is_error']=True;fixtures.append(('paired tool result is_error',failed))
for name,event in fixtures:
 facts=m.classify_actual(messages+[event],dom['notice'],'idle',False);assert not facts['ok'];rows.append({'case':name,'passed':True,'facts':facts})
facts=m.classify_actual(messages,'Server message rejected: controlled malformed JSON','idle',False);assert not facts['ok'];rows.append({'case':'actual malformed/server-rejected notice prefix','passed':True,'facts':facts})
missing=[x for x in messages if x['type']!='tool_finished'];facts=m.classify_actual(missing,dom['notice'],'idle',False);assert not facts['ok'];rows.append({'case':'missing actual tool-result pairing','passed':True,'facts':facts})
report={'controller_sha256':expected,'control_script_sha256':hashlib.sha256(Path(__file__).read_bytes()).hexdigest(),'source_reference':str(e),'scope':'ONE compact no-browser control. Captured actual positive; explicit real-protocol-shaped negatives. No production mocking, browser/server/native/key/provider activity.','cases':rows,'status':'pass'};Path(__file__).with_name('classifier-control-result.json').write_text(json.dumps(report,indent=2)+'\n');print('PASS compact classifier control:',len(rows),'discriminating cases; no browser/native/provider/key execution')
