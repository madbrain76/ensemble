#!/usr/bin/env python3
"""Verify copied correction text/receipt/report identities; execute nothing else."""
import hashlib,json,pathlib
root=pathlib.Path(__file__).resolve().parent
index=root.parent/'ch05-student-r1-evidence-manifest.json'
data=json.loads(index.read_text())
for name,meta in data['text_files'].items():
 path=(root/name).resolve();assert path.is_relative_to(root)
 raw=path.read_bytes();assert len(raw)==meta['bytes'] and hashlib.sha256(raw).hexdigest()==meta['sha256'],name
for key in ['report','source_index']:
 meta=data[key];path=root.parent/meta['relative_path'];assert hashlib.sha256(path.read_bytes()).hexdigest()==meta['sha256'],key
print(json.dumps({'verified_text_files':len(data['text_files']),'report_verified':True,'source_index_verified':True,'candidate_commit':data['candidate']['source_commit'],'executed_candidate':False,'credential_reads':0},indent=2))
