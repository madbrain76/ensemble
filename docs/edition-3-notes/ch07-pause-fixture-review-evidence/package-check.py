from pathlib import Path
import json,hashlib,subprocess,datetime
r=Path('/Users/bill/projects/ensemble');w=Path('/tmp/cadence-ch07-pause-fixture-review-independent');student=Path('/Users/bill/projects/ensemble-edition-3-student');sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
facts={}
for name in ['positive-retry','pause-bypass','restored']:
 phases={}
 for stage in ['paused-check','return']:
  d=json.loads((w/f'capture-{name}.{stage}.json').read_text());x={k:d[k] for k in ['stage','work','error','second_exists','second_value','provider_requests']};x['tool_messages']=[]
  for m in d['messages']:
   if m.get('type') not in ['tool_dispatched','tool_finished']:continue
   q={k:m.get(k) for k in ['type','call_id','part_id','is_error']};q['replay']=m.get('replay',False)
   if m['type']=='tool_finished':q['known_result_markers']=[a for a in ['GUI_FIRST','GUI_SECOND','GUI_THIRD'] if a in str(m.get('result',''))];q['result_sha256']=hashlib.sha256(str(m.get('result','')).encode()).hexdigest()
   x['tool_messages'].append(q)
  x['dispatch_count']=sum(m['type']=='tool_dispatched' for m in x['tool_messages']);x['nonerror_finish_count']=sum(m['type']=='tool_finished' and m['is_error'] is False for m in x['tool_messages']);phases[stage]=x
 rows=[]
 for line in (Path(phases['return']['work'])/'gui.log').read_text().splitlines():
  pos=line.index('{');prefix=line[:pos].split();m=json.loads(line[pos:]);kind=m.get('type',m.get('command'))
  if kind in ['pause','unpause','tool_dispatched','tool_finished']:rows.append({'timestamp':prefix[0],'direction':prefix[1],'type':kind,'call_id':m.get('call_id'),'is_error':m.get('is_error')})
 phases['selected_server_rows']=rows;facts[name]=phases
for name in ['positive-retry','restored']:assert facts[name]['paused-check']['dispatch_count']==1 and not facts[name]['paused-check']['second_exists']
assert facts['pause-bypass']['paused-check']['dispatch_count']==3 and facts['pause-bypass']['paused-check']['second_value']=='GUI_SECOND'
(w/'observed-facts.json').write_text(json.dumps(facts,indent=2)+'\n')
expected=json.loads((w/'source-input-hashes.json').read_text());actual={str(p.relative_to(w/'source')):sha(p) for p in (w/'source').rglob('*') if p.is_file()};assert actual==expected
for n,h in expected.items():assert sha(student/n)==h,n
assert subprocess.check_output(['git','status','--porcelain'],cwd=student).strip()==b'';head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=student).decode().strip();assert head=='48e1d12d9c396b5988c0ce6eccdef67e581872ff'
inputs=json.loads((w/'harness-input-hashes.json').read_text());assert len(inputs)==166
for n,h in inputs.items():assert sha(r/n)==h,n
root_e=r/'docs/edition-3-notes/ch07-pause-fixture-fix-evidence';candidate=json.loads((root_e/'candidate.json').read_text())
for n,v in candidate['paths'].items():assert sha(r/n)==v['sha256'],n
assert sha(Path(candidate['binary']['path']))==candidate['binary']['sha256']
old=json.loads((r/'docs/edition-3-notes/ch07-pause-diagnostic-evidence/preservation.json').read_text());baseline=Path('/tmp/cadence-ch07-publication-preservation-after-body11.json');pub=json.loads(baseline.read_text());pubhash={}
def collect(o):
 if isinstance(o,dict):
  for k,v in o.items():
   if k in old['publication12']:pubhash[k]=v if isinstance(v,str) else v['sha256']
   else:collect(v)
 elif isinstance(o,list):
  for v in o:collect(v)
collect(pub);assert pubhash==old['publication12']
for n,h in pubhash.items():assert sha(r/n)==h,n
bundles={}
for name in ['ch07-entry-evidence','ch07-grader-review-evidence','ch07-pause-diagnostic-evidence']:
 d=r/'docs/edition-3-notes'/name;manifest=d/'SHA256SUMS';names={'SHA256SUMS'}
 for line in manifest.read_text().splitlines():
  h,n=line.split('  ',1);assert sha(r/n)==h,n
  if Path(n).parent==d.relative_to(r):names.add(Path(n).name)
 assert names=={p.name for p in d.iterdir()},name;bundles[name]={'manifest_sha256':sha(manifest),'files':len(names)}
initial=json.loads((r/'docs/edition-3-notes/ch07-pause-diagnostic-evidence/initial-verification.json').read_text());package=r/'docs/edition-3-notes/ch07-student-evidence';indices=[p for p in [package/'r4-discrepancy-package-index.json'] if json.loads(p.read_text())==initial['package_index']];assert len(indices)==1
for n,v in initial['package_index']['files'].items():assert sha(package/n)==v['sha256'],n
for n,h in initial['artifacts'].items():assert sha(student/n)==h,n
pres={'verified_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'student_head':head,'student_clean':True,'student_source117_match':True,'scratch_exact_source117_restored':True,'production_internal166_unchanged_since_repaired_pin':True,'candidate':candidate,'publication12':pubhash,'publication_baseline_path':str(baseline),'publication_baseline_sha256':sha(baseline),'prior_bundles':bundles,'package31_unchanged':True,'package_index_path':str(indices[0]),'package_index_sha256':sha(indices[0]),'frozen_artifacts':initial['artifacts'],'root_receipt_hashes':{p.name:sha(p) for p in root_e.iterdir() if p.is_file()},'main_head':subprocess.check_output(['git','rev-parse','HEAD'],cwd=r).decode().strip()}
(w/'preservation.json').write_text(json.dumps(pres,indent=2)+'\n')
print('PASS frozen117/exactrestoration/internal166/adapter/repairedgrader/binaries/31package/12publication/previous47+71+76bundles/studentclean')
for name,phases in facts.items():print(name,{s:{k:phases[s][k] for k in ['dispatch_count','nonerror_finish_count','second_exists','provider_requests','error']} for s in ['paused-check','return']})
