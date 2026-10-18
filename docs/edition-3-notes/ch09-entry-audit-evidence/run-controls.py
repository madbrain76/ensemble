#!/usr/bin/env python3
import os,pathlib,subprocess,json,hashlib,shutil,time,signal,tarfile,difflib
R=pathlib.Path('/Users/bill/projects/ensemble'); E=R/'docs/edition-3-notes/ch09-entry-audit-evidence'; S=pathlib.Path('/tmp/ch09-entry-audit')
for d in ['home','tmp','bin','work']: (S/d).mkdir(exist_ok=True)
go=shutil.which('go'); env={'PATH':str(pathlib.Path(go).parent)+':/usr/bin:/bin:/usr/sbin:/sbin','HOME':str(S/'home'),'TMPDIR':str(S/'tmp'),'GOCACHE':'/Users/bill/Library/Caches/go-build','GOMODCACHE':'/Users/bill/go/pkg/mod','GOPATH':'/Users/bill/go','GOWORK':'off','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOFLAGS':'-mod=readonly','LANG':'C','TZ':'UTC'}
(E/'environment.json').write_text(json.dumps({'go':go,'env':env,'provider_credentials':'none; unchanged graders inject fake values and local endpoints','network':'Go proxy/sumdb disabled; test fake providers only'},indent=2)+'\n')
def tree(p):
 return {str(f.relative_to(p)):{'sha256':hashlib.sha256(f.read_bytes()).hexdigest(),'bytes':f.stat().st_size} for f in sorted(p.rglob('*')) if f.is_file()}
def run(name,args,cwd,limit):
 start=time.time(); meta={'argv':list(map(str,args)),'cwd':str(cwd),'timeout_seconds':limit,'environment_file':'environment.json','start_unix':start}
 with (E/(name+'.stdout')).open('wb') as out,(E/(name+'.stderr')).open('wb') as err:
  p=subprocess.Popen(list(map(str,args)),cwd=cwd,env=env,stdout=out,stderr=err,start_new_session=True); meta['pid']=p.pid
  try: code=p.wait(timeout=limit)
  except subprocess.TimeoutExpired:
   meta['timeout']=True; os.killpg(p.pid,signal.SIGTERM)
   try: code=p.wait(timeout=3)
   except subprocess.TimeoutExpired: os.killpg(p.pid,signal.SIGKILL); code=p.wait()
  meta.update(exit_code=code,elapsed_seconds=round(time.time()-start,3))
 (E/(name+'.json')).write_text(json.dumps(meta,indent=2)+'\n'); print(name,code,flush=True); return code
original=tree(R/'solutions/ch09'); C=S/'positive'; shutil.copytree(R/'solutions/ch09',C); assert tree(C)==original
(E/'positive-source.json').write_text(json.dumps({'source':str(R/'solutions/ch09'),'copy':str(C),'files':original},indent=2)+'\n')
with tarfile.open(E/'positive-source.tar','w') as tar: tar.add(C,arcname='ch09')
# Build-only preparation; stop if any build fails. No repairs/retries.
for name,args,cwd in [
 ('build-entry',[go,'build','-o',S/'bin/grade-edition3','./cmd/grade-edition3'],R),
 ('build-driver',[go,'build','-o',S/'bin/driver','./docs/edition-3-notes/ch09-entry-audit-evidence/driver'],R),
 ('build-cli',[go,'build','-o',S/'bin/cli','.'],C),
 ('build-gui',[go,'build','-o',S/'bin/gui','./cmd'],C/'gui')]:
 if run(name,args,cwd,120): raise SystemExit('preparation failed: '+name)
(E/'binaries.json').write_text(json.dumps(tree(S/'bin'),indent=2)+'\n')
run('entry-command',[S/'bin/grade-edition3','-ch','9','-json',C],R,15)
run('positive',[S/'bin/driver',C,S/'bin/cli',S/'bin/gui',R/'web/gui'],S/'work',240)
# Mutation changes actual source only. Test, graders and deadlines are unchanged.
M=S/'mutant'; shutil.copytree(C,M); f=M/'internal/skills/skills.go'; old=f.read_text(); anchor='\tr.pending[name] = true\n\treturn nil\n'; assert old.count(anchor)==1
new=old.replace(anchor,'\tr.pending[name] = true\n\tfor _, tool := range r.all[name].tools {\n\t\tif r.tools[tool].Source == common.SourceDynamic {\n\t\t\tdelete(r.tools, tool)\n\t\t}\n\t}\n\treturn nil\n'); f.write_text(new)
(E/'mutation.diff').write_text(''.join(difflib.unified_diff(old.splitlines(True),new.splitlines(True),fromfile='positive/internal/skills/skills.go',tofile='mutant/internal/skills/skills.go')))
mutant=tree(M); changes=[k for k in original if mutant.get(k)!=original[k]]; assert changes==['internal/skills/skills.go']; assert set(mutant)==set(original)
(E/'mutant-source.json').write_text(json.dumps({'source':str(M),'changed_files':changes,'files':mutant},indent=2)+'\n')
run('sensitivity',[go,'test','-count=1','-run','^TestSkillsThroughLoop$','./ensemble'],M,90)
assert tree(C)==original
(E/'source-postcheck.json').write_text(json.dumps({'positive_unchanged':True,'mutant_only_expected_change':tree(M)==mutant},indent=2)+'\n')
