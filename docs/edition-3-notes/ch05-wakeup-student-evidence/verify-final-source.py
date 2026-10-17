import pathlib,hashlib,json,subprocess,tarfile,io
P=pathlib.Path('/tmp/cadence-ch05-wakeup-evidence');d=json.loads((P/'candidate-r2/frozen/identity.json').read_text());r=pathlib.Path(d['source_root']);sha=lambda b:hashlib.sha256(b).hexdigest()
assert not subprocess.check_output(['git','status','--porcelain'],cwd=r)
assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=r,text=True).strip()==d['source_commit']
assert subprocess.check_output(['git','rev-parse','HEAD^{tree}'],cwd=r,text=True).strip()==d['source_tree']
archive=pathlib.Path(d['source_archive']['path']).read_bytes();assert sha(archive)==d['source_archive']['sha256'];assert archive==subprocess.check_output(['git','archive',d['source_commit']],cwd=r)
snapshot=P/'source-snapshot';snapshot.mkdir(exist_ok=True);seen={}
with tarfile.open(fileobj=io.BytesIO(archive)) as tar:
 for member in tar.getmembers():
  if not member.isfile():continue
  n=member.name;assert not n.startswith('/') and '..' not in pathlib.Path(n).parts
  b=tar.extractfile(member).read();assert sha(b)==d['source_files'][n];assert b==(r/n).read_bytes();assert b==subprocess.check_output(['git','show',d['source_commit']+':'+n],cwd=r)
  b.decode('utf-8');f=snapshot/n;f.parent.mkdir(parents=True,exist_ok=True);f.write_bytes(b);seen[n]=sha(b)
assert seen==d['source_files']
for name,m in d['binaries'].items():
 assert sha(pathlib.Path(m['path']).read_bytes())==m['sha256'];info=subprocess.check_output(['go','version','-m',m['path']],text=True);assert 'vcs.revision='+d['source_commit'] in info;assert 'vcs.modified=false' in info;assert '-race=true' in info
 (P/f'final-{name}-buildinfo.txt').write_text(info)
(P/'final-source-verification.json').write_text(json.dumps({'passed':True,'source_commit':d['source_commit'],'tree':d['source_tree'],'clean':True,'git_worktree_archive_snapshot_files':len(seen),'archive':d['source_archive'],'race_binaries':d['binaries']},indent=2)+'\n');print('verified',len(seen),'source files, exact archive, clean VCS, and race binaries')
