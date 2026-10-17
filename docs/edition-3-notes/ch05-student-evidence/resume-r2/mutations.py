import pathlib,shutil,subprocess,json,tempfile
source=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
mutations=[
 ('drop_hint','internal/history/history.go','case "hint_received":\n\t\tc.Ephemera = append(c.Ephemera, event.Message.Parts...)','case "hint_received":','^TestChapter5Hints'),
 ('ignore_batch_cap','internal/agent/actor.go','if t.batches >= limit {','if t.batches >= limit && false {','^TestChapter5Configured'),
 ('drop_late_job_report','internal/agent/actor.go','_ = a.append(common.Event{Type: "job_report", Job: &r})','_ = r','^TestChapter5Interrupt'),
 ('kill_ordinary_on_shutdown','internal/jobs/jobs.go','process := j.processTool || j.process','process := true','^TestChapter5Interrupt'),
 ('skip_request_save','internal/agent/actor.go','if saveErr := a.save(); saveErr != nil {','if saveErr := error(nil); saveErr != nil {','^TestChapter5Hints'),
 ('unknown_media_allowed','internal/engine/models.go','return common.MediaCapabilities{}','return common.MediaCapabilities{Images:true,Media:common.MediaImage|common.MediaAudio|common.MediaVideo|common.MediaDocument}','^TestChapter5Media'),
]
results=[]
output_root=pathlib.Path('/tmp/cadence-ch05-evidence-r1/resume-r2')
for name,file,before,after,test in mutations:
 if shutil.disk_usage(source).free < 1024**3: raise RuntimeError('disk headroom below 1GiB; no mutation launched')
 with tempfile.TemporaryDirectory(prefix='ch05-mutation-') as temp:
  root=pathlib.Path(temp)
  for p in source.iterdir():
   if p.suffix in ['.go','.mod','.sum'] and p.is_file():shutil.copy2(p,root/p.name)
  for folder in ['internal','cmd']:shutil.copytree(source/folder,root/folder)
  p=root/file;data=p.read_text();assert data.count(before)==1,(name,data.count(before));p.write_text(data.replace(before,after))
  try:
   proc=subprocess.run(['go','test','-timeout','12s','-run',test,'./...'],cwd=root,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=35);code=proc.returncode;output=proc.stdout
  except subprocess.TimeoutExpired as err:code=124;output=str(err.stdout)
  (output_root/('mutation-'+name+'.txt')).write_text(output)
  environmental=any(marker in output.lower() for marker in ['no space left','build failed','cannot allocate memory','resource temporarily unavailable'])
  assertion='--- FAIL: TestChapter5' in output and 'chapter5_test.go:' in output
  detected=code!=0 and code!=124 and not environmental and assertion
  results.append({'name':name,'file':file,'before':before,'after':after,'test':test,'exit':code,'detected':detected,'classification':'test assertion rejected mutation' if detected else 'INCONCLUSIVE; inspect receipt'});print(results[-1],flush=True)
(output_root/'mutations.json').write_text(json.dumps(results,indent=2))
assert all(r['detected'] for r in results)
