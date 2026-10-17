from pathlib import Path
import json,sys,subprocess,shutil,hashlib,difflib
E=Path(__file__).resolve().parent;S=Path((E/'scratch-path.txt').read_text().strip());R=Path('/Users/bill/projects/ensemble');M=S/'mutations';M.mkdir();results=[]
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
def run(label,cwd,argv,want=0):
 code=subprocess.call([sys.executable,str(E/'run.py'),label,str(cwd),*map(str,argv)],cwd=R)
 assert code==want,(label,code,want)
 return json.loads((E/(label+'.meta.json')).read_text())
pattern='^(TestRealRequestModesAndUnknownFallback|TestOpenAISparseInterleavedToolsAndUsage)$'
base=M/'baseline.test';run('sensitivity-baseline-compile',S/'public-checks',['go','test','-race','-c','-o',base,'.']);run('sensitivity-baseline',S/'public-checks',[base,'-test.v','-test.run',pattern]);basehash=sha(base)
for name,rel,old,new,selected in [
 ('ignore-disable','internal/engine/models.go','if cfg.DisableStreaming {','if cfg.DisableStreaming && false {','^TestRealRequestModesAndUnknownFallback$'),
 ('wrong-final-id','internal/engine/stream.go','cb.OnPartFinal(ids[i], notify)','cb.OnPartFinal(ids[i]+100, notify)','^TestOpenAISparseInterleavedToolsAndUsage$')]:
 src=S/('mutation-source-'+name);mod=S/('mutation-tests-'+name);shutil.copytree(S/'source',src);shutil.copytree(S/'public-checks',mod)
 p=mod/'go.mod';text=p.read_text();assert 'replace ensemble => ../source' in text;p.write_text(text.replace('replace ensemble => ../source','replace ensemble => '+str(src)))
 target=src/rel;original=target.read_text();assert original.count(old)==1;changed=original.replace(old,new);target.write_text(changed);assert sha(target)!=sha(S/'source'/rel)
 (E/(name+'.diff')).write_text(''.join(difflib.unified_diff(original.splitlines(True),changed.splitlines(True),fromfile='pinned/'+rel,tofile='mutant/'+rel)))
 before=sha(S/'source'/rel);mutated=sha(target);binary=M/(name+'.test');restored=M/(name+'-restored.test')
 run(name+'-compile',mod,['go','test','-race','-c','-o',binary,'.']);broken=run(name+'-behavior',mod,[binary,'-test.v','-test.run',selected],1)
 target.write_text(original);assert sha(target)==before
 expected={p.relative_to(S/'source').as_posix():sha(p) for p in (S/'source').rglob('*') if p.is_file()};actual={p.relative_to(src).as_posix():sha(p) for p in src.rglob('*') if p.is_file()};assert actual==expected
 run(name+'-restored-compile',mod,['go','test','-race','-c','-o',restored,'.']);good=run(name+'-restored-behavior',mod,[restored,'-test.v','-test.run',selected])
 results.append({'name':name,'source_file':rel,'original_sha256':before,'mutant_sha256':mutated,'restored_sha256':sha(target),'restored_full_source_equal':True,'changed_bytes_compile_exit':0,'targeted_failure_exit':broken['exit'],'restored_behavior_exit':good['exit'],'binary':str(binary),'binary_sha256':sha(binary),'restored_binary':str(restored),'restored_binary_sha256':sha(restored)})
 (E/'mutation-results.json').write_text(json.dumps({'baseline_binary':str(base),'baseline_binary_sha256':basehash,'controls':results},indent=2)+'\n')
print(json.dumps(results,indent=2))
