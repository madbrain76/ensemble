import hashlib,json,pathlib,re
R=pathlib.Path('/Users/bill/projects/ensemble');SE=R/'docs/edition-3-notes/ch04-student-r1-evidence';E=R/'docs/edition-3-notes/ch04-code-review-r1-evidence';sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest();summary=[];blob='abcde'*8192
def text(parts):return parts if isinstance(parts,str) else ''.join(x.get('text','') for x in parts)
for mode in ['core','consumer']:
 D=SE/('local-'+mode);r=json.loads((D/'receipt.json').read_text());records=[json.loads(s) for s in (D/'events.jsonl').read_text().splitlines()];assert records[0]=={'log_version':1};events=records[1:];called={x['tool']['call_id']:x['tool'] for x in events if x.get('type')=='tool_called'};returned={x['tool']['call_id']:x['tool'] for x in events if x.get('type')=='tool_returned'};assert len(called)==len(returned)==9;assert r['exit_code']==0 and r['passed'] and not r['failures'];assert r['binary_sha256']==sha(pathlib.Path(r['command'][0]));assert len(r['requests'])==4;assert all(q['path']=='/v1/messages' and q['model']=='local-c1-boundary-fixture' for q in r['requests']);assert 'LOCAL C1 CHECKS PASSED' in (D/'stdout.txt').read_text();sizes=[]
 for batch in r['observations']:
  n=batch['batch'];rs=batch['returns'];assert len(rs)==3
  for i,v in enumerate(rs):
   ident=f'local-{n}-{i}';assert v['tool_use_id']==ident;assert called[ident]['name']==returned[ident]['name'];assert text(v['content'])==text(returned[ident]['parts']);assert bool(v.get('is_error'))==bool(returned[ident].get('is_error'))
  control=text(rs[1]['content']);out=text(rs[2]['content']);assert 'tool_limits consumed by this tool_limits call' in control;assert len(out.encode())==batch['result_bytes'];sizes.append(batch['result_bytes'])
  if n==0:
   assert 'max_output_bytes 1024' in control and 'max_output_bytes 16384' in out and 'ai_callback_delay 300ms' in out;assert 8192<len(out)<17000 and 'bytes omitted' in out
  elif n==1:
   assert rs[1]['is_error'] and 'error:' in control and 'tool_limits consumed' not in out;assert 8192<len(out)<17000
  elif n==2:
   assert called['local-2-2']['args']['max_output_bytes']==8192;assert 'max_output_bytes 2048' in out and 4096<len(out)<8500 and 'exit_code: 0' in out
  else:raise AssertionError(n)
 files=[]
 for i,j in enumerate(r['actual_ordinary_jobs_done']):
  assert j['status']=='done';p=pathlib.Path(r['data_dir'])/j['output']['locator'];data=p.read_text();assert len(data)>=40960;expected=blob.upper() if mode=='consumer' and i<2 else blob;assert data.startswith(expected);assert len(data.encode())==j['bytes'];files.append({'handle':j['handle'],'bytes':len(data.encode()),'sha256':sha(p)})
 assert len(files)==3
 diag=[s for s in (D/'stderr.txt').read_text().splitlines() if 'shout called' in s]
 if mode=='consumer':assert len(diag)==2 and all(re.search(r'\d{4}[-/]\d{2}[-/]\d{2}',s) for s in diag)
 summary.append({'mode':mode,'receipt_sha256':sha(D/'receipt.json'),'events_sha256':sha(D/'events.jsonl'),'binary_sha256':r['binary_sha256'],'commit':r['commit'],'tree':r['tree'],'requests':4,'correlated_tool_calls':9,'observed_result_bytes':sizes,'persisted_outputs':files,'registered_diagnostics':len(diag),'audit_passed':True,'classification':'independent audit of retained loopback-HTTP run; not a newly executed run or paid provider evidence'})
(E/'retained-local-audit.json').write_text(json.dumps(summary,indent=2)+'\n');print(json.dumps(summary,indent=2))
