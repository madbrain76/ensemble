#!/usr/bin/env python3
import pathlib,hashlib,json,shutil
E=pathlib.Path(__file__).resolve().parent;S=pathlib.Path((E/'scratch-path.txt').read_text().strip());base=S/'solutions/ch06'
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def manifest(p):return {str(f.relative_to(p)):sha(f) for f in sorted(p.rglob('*')) if f.is_file()}
b=manifest(base);out=[]
for m in json.loads((E/'mutations.json').read_text()):
 p=pathlib.Path(m['source']);before=manifest(p)
 differences=[f for f in set(b)|set(before) if b.get(f)!=before.get(f)]
 assert differences==[m['file']],differences
 (E/('mutant-'+m['name']+'-source-sha256.json')).write_text(json.dumps(before,indent=2)+'\n')
 h=sha(S/m['name']);shutil.copyfile(base/m['file'],p/m['file']);after=manifest(p)
 assert after==b
 out.append({'name':m['name'],'mutant_binary':str(S/m['name']),'mutant_binary_sha256':h,'changed_files':differences,'restored_source_matches_pinned':True,'restored_file_sha256':sha(p/m['file'])})
(E/'restoration.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
