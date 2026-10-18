import hashlib,json,pathlib,subprocess,sys,tarfile
base=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student/.ch09-local')
checks=[('hidden-dispatch-r2','internal/agent/dispatch.go','if a.skills != nil && !a.skills.Enabled(call.Name) {','if false && a.skills != nil && !a.skills.Enabled(call.Name) {'),('first-source','internal/skills/registry.go','if _, ok := r.enabled[tool]; !ok {','if true {')]
records=[]
for name,file,old,new in checks:
 target=base/('control-'+name);target.mkdir()
 with tarfile.open(base/'source.tar') as archive: archive.extractall(target,filter='data')
 path=target/file;original=path.read_bytes();text=original.decode()
 if old not in text: raise SystemExit('missing mutation anchor')
 changed=text.replace(old,new,1)
 if name=='hidden-dispatch-r2':
  allocation='(a.skills == nil || a.skills.Enabled(call.Name))'
  lookup='for _, candidate := range a.Tools() {'
  if allocation not in changed or lookup not in changed: raise SystemExit('missing complete-dispatch anchors')
  changed=changed.replace(allocation,'true',1).replace(lookup,'a.mu.RLock(); installed := append([]common.Tool(nil), a.tools...); a.mu.RUnlock(); for _, candidate := range installed {',1)
 path.write_text(changed)
 mutated=path.read_bytes()
 cmd=['go','test','-count=1','-v','-run','^TestChapter9CapabilityLifecycle$','.']
 bad=subprocess.run([sys.executable,str(base/'gate.py'),'control-'+name+'-negative',str(target),*cmd]).returncode
 stdout=(base/('control-'+name+'-negative.stdout.txt')).read_text();stderr=(base/('control-'+name+'-negative.stderr.txt')).read_text()
 diagnostic='write_file result=' if name=='hidden-dispatch-r2' else 'first provenance overwritten'
 credited=bad!=0 and diagnostic in stdout and 'build failed' not in stdout.lower() and 'build failed' not in stderr.lower()
 path.write_bytes(original)
 good=subprocess.run([sys.executable,str(base/'gate.py'),'control-'+name+'-restored',str(target),*cmd]).returncode
 records.append(dict(name=name,file=file,original_sha256=hashlib.sha256(original).hexdigest(),mutated_sha256=hashlib.sha256(mutated).hexdigest(),restored_sha256=hashlib.sha256(path.read_bytes()).hexdigest(),mutation_from=old,mutation_to=new,negative_exit=bad,restored_exit=good,diagnostic=diagnostic,credited=credited and good==0))
 (base/'sensitivity-summary-r2.json').write_text(json.dumps(records,indent=2)+'\n')
 if not credited or good: raise SystemExit('uncredited control; inspect receipts')
