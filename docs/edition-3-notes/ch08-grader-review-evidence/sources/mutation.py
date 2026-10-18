import sys;sys.dont_write_bytecode=True
from common import *
prior=R/'docs/edition-3-notes/ch08-entry-reconciliation-evidence';P=Path((prior/'scratch-path.txt').read_text().strip());index={x['path']:x for x in json.loads((prior/'local-artifacts.json').read_text())['files']};reused=[]
for n in ['bin/driver','driver/main.go','driver/go.mod','driver/go.sum']:
 p=P/n;assert sha(p)==index[str(p)]['sha256'];reused.append(index[str(p)])
for n in ['main.go','go.mod','go.sum']:shutil.copyfile(P/'driver'/n,T/'driver'/n);shutil.copyfile(P/'driver'/n,E/'sources'/('settings-driver-'+n+'.txt'))
(E/'settings-driver-provenance.json').write_text(json.dumps({'old_driver':reused[0],'source_files':reused[1:],'old_check_source':json.loads((E/'source-provenance.json').read_text())['old_settings_source'],'candidate_check_source':json.loads((E/'pins-before.json').read_text())['candidate']['sources']['internal/gradeedition2/ch08.go'],'note':'Same retained caller source. Old retained binary links pre-candidate checks; new build links pinned candidate. Only settings mode used.'},indent=2)+'\n')
assert run('build-candidate-settings-driver',['go','build','-o',str(T/'bin/candidate-settings-driver'),'.'],T/'driver')==0
G=T/'real-gui';shutil.copytree(R/'solutions/ch08/gui',G)
mod=(G/'go.mod').read_text();assert mod.count('replace ensemble => ..')==1;(G/'go.mod').write_text(mod.replace('replace ensemble => ..','replace ensemble => '+str(R/'solutions/ch08')))
p=G/'gui.go';before=p.read_bytes();text=before.decode();anchor='\tadd(map[string]any{"type": "current_settings", "settings": s.snapshotSettings()})';assert text.count(anchor)==1
replacement='''\t// OWNED AUDIT MUTANT: omit only pre-update current_settings budget.
\tvar settings any = s.snapshotSettings()
\tif _, err := os.Stat(s.settingsPath()); os.IsNotExist(err) {
\t\traw, _ := json.Marshal(settings)
\t\tvar initial map[string]any
\t\tjson.Unmarshal(raw, &initial)
\t\tdelete(initial, "max_tool_rounds")
\t\tsettings = initial
\t}
\tadd(map[string]any{"type": "current_settings", "settings": settings})'''
assert text.count('\t"net/http"\n')==1;mutant=text.replace('\t"net/http"\n','\t"net/http"\n\t"os"\n').replace(anchor,replacement);p.write_text(mutant)
(E/'sources/gui-original.go.txt').write_bytes(before);(E/'sources/gui-mutant.go.txt').write_bytes(p.read_bytes());(E/'sources/scratch-gui-go.mod.txt').write_bytes((G/'go.mod').read_bytes())
base={str(x.relative_to(G)):{'bytes':x.stat().st_size,'sha256':sha(x)} for x in G.rglob('*') if x.is_file()};(E/'scratch-gui-mutant-hashes.json').write_text(json.dumps(base,indent=2)+'\n')
results=[]
try:
 for name,driver,expected in [('old-check-omitted-budget',P/'bin/driver',True),('candidate-check-omitted-budget',T/'bin/candidate-settings-driver',False)]:
  rc=run(name,[str(driver),'-mode','settings','-gui',str(G),'-gui-model','claude-sonnet-4-6'],T/'work',90);r=json.loads((E/(name+'.stdout')).read_text());failed=[c['id'] for c in r['checks'] if not c['passed']];results.append({'name':name,'source_sha256':sha(p),'driver_sha256':sha(driver),'exit':rc,'score':r['score'],'max_score':r['max_score'],'passed':r['passed'],'failed_checks':failed});assert r['passed']==expected and rc==(0 if expected else 1),(name,r)
  if not expected:assert failed==['settings-onconnect'],r
finally:
 p.write_bytes(before);(E/'sources/gui-restored.go.txt').write_bytes(p.read_bytes());(E/'mutation-partial-results.json').write_text(json.dumps(results,indent=2)+'\n')
name='candidate-check-restored';driver=T/'bin/candidate-settings-driver';rc=run(name,[str(driver),'-mode','settings','-gui',str(G),'-gui-model','claude-sonnet-4-6'],T/'work',90);r=json.loads((E/(name+'.stdout')).read_text());results.append({'name':name,'source_sha256':sha(p),'driver_sha256':sha(driver),'exit':rc,'score':r['score'],'max_score':r['max_score'],'passed':r['passed'],'failed_checks':[c['id'] for c in r['checks'] if not c['passed']]});assert rc==0 and r['passed'] and r['score']==80 and p.read_bytes()==before
(E/'mutation-results.json').write_text(json.dumps(results,indent=2)+'\n');print(json.dumps(results,indent=2))
