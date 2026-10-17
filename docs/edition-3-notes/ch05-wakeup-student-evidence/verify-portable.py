#!/usr/bin/env python3
"""Hash-only verification of this portable amendment; no execution or credentials."""
import hashlib,json,pathlib
root=pathlib.Path(__file__).resolve().parent;data=json.loads((root.parent/'ch05-wakeup-student-evidence-manifest.json').read_text())
for name,m in data['text_files'].items():
 p=(root/name).resolve();assert p.is_relative_to(root);b=p.read_bytes();assert len(b)==m['bytes'] and hashlib.sha256(b).hexdigest()==m['sha256'],name
for key in ['report','source_index']:
 m=data[key];p=(root.parent/m['relative_path']).resolve();assert p.parent==root.parent;assert hashlib.sha256(p.read_bytes()).hexdigest()==m['sha256'],key
source=json.loads((root.parent/data['source_index']['relative_path']).read_text())
for name,digest in source['source_files'].items():assert hashlib.sha256((root/'source-snapshot'/name).read_bytes()).hexdigest()==digest,name
print(json.dumps({'verified_text_files':len(data['text_files']),'source_snapshot_files':len(source['source_files']),'report_verified':True,'source_index_verified':True,'candidate_commit':source['source_commit'],'executed_candidate':False,'credential_reads':0},indent=2))
