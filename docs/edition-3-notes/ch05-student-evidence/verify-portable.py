#!/usr/bin/env python3
"""Verify portable text evidence and report hashes; no execution/network/credentials."""
import hashlib,json,pathlib,sys
root=pathlib.Path(__file__).resolve().parent
index=pathlib.Path(sys.argv[1]) if len(sys.argv)>1 else root.parent/'ch05-student-evidence-manifest.json'
data=json.loads(index.read_text())
for name,item in data['text_files'].items():
    path=(root/name).resolve()
    if not path.is_relative_to(root):raise SystemExit('unsafe evidence path')
    raw=path.read_bytes()
    if hashlib.sha256(raw).hexdigest()!=item['sha256'] or len(raw)!=item['bytes']:raise SystemExit('mismatch: '+name)
report=index.parent/data['report']['relative_path']
if hashlib.sha256(report.read_bytes()).hexdigest()!=data['report']['sha256']:raise SystemExit('report mismatch')
print('verified',len(data['text_files']),'portable text files and report; no providers or credentials accessed')
