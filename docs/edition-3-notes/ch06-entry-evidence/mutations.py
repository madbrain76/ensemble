#!/usr/bin/env python3
import pathlib,hashlib,json,shutil,difflib
E=pathlib.Path(__file__).resolve().parent; S=pathlib.Path((E/'scratch-path.txt').read_text().strip()); base=S/'solutions/ch06'
cases=[('disable','internal/engine/media.go','if cfg.DisableStreaming {','if cfg.DisableStreaming && false {'),('correlation','internal/engine/stream.go','cb.OnPartFinal(ids[i], p)','cb.OnPartFinal(ids[i]+100, p)')]
records=[]
for label,rel,old,new in cases:
 target=S/('mutant-'+label);shutil.copytree(base,target)
 p=target/rel;before=p.read_text();assert before.count(old)==1
 after=before.replace(old,new);p.write_text(after)
 patch=''.join(difflib.unified_diff(before.splitlines(True),after.splitlines(True),fromfile=rel,tofile=rel))
 (E/(label+'.diff')).write_text(patch)
 records.append({'name':label,'source':str(target),'file':rel,'before_sha256':hashlib.sha256(before.encode()).hexdigest(),'after_sha256':hashlib.sha256(after.encode()).hexdigest(),'old':old,'new':new})
(E/'mutations.json').write_text(json.dumps(records,indent=2)+'\n')
print(json.dumps(records,indent=2))
