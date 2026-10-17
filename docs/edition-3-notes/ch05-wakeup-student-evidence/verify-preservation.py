import hashlib,json,pathlib
N=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes');P=pathlib.Path('/tmp/cadence-ch05-wakeup-evidence');sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
pins={'ch05-student-review.md':'2c47fb072bb6a93d17fa6738c21cdbd0f306c88e32947b0539cb961c596c9020','ch05-student-evidence-manifest.json':'7180c6e8ec7d1cc63d53d33e5e70a5c292e8aa4a3ca05c8121597c8f6d47bff0','ch05-student-review-r1.md':'9b6dca6bee16bf58eb739578f2974e2cecfe44ebc8cfe167ad85e95692662966','ch05-student-r1-source-manifest.json':'412830507167f81571642c1557ba3fcf16153baa682ecd969f54ec02ce366796','ch05-student-r1-evidence-manifest.json':'c89fc3177764450345e656d6e19350d3d9830af486c79c007c0eb17348dbb0fd'}
for name,digest in pins.items():assert sha(N/name)==digest,name
rows=[]
for index,local in [('ch05-student-evidence-manifest.json','/tmp/cadence-ch05-evidence-r1'),('ch05-student-r1-evidence-manifest.json','/tmp/cadence-ch05-correction-r1-evidence')]:
 data=json.loads((N/index).read_text());local=pathlib.Path(local);dest=N/data['portable_directory']
 for name,meta in data['text_files'].items():
  assert sha(dest/name)==meta['sha256'],name;assert sha(local/name)==meta['sha256'],name
 for name,meta in data['excluded_artifacts'].items():assert sha(pathlib.Path(meta['local_path']))==meta['sha256'],name
 rows.append({'index':index,'portable_texts':len(data['text_files']),'local_texts':len(data['text_files']),'excluded_local_artifacts':len(data['excluded_artifacts']),'passed':True})
# The original source snapshot remains intentionally at the initial candidate.
identity=json.loads(pathlib.Path('/tmp/cadence-ch05-evidence-r1/frozen/identity.json').read_text())
files=identity['source_files'];snapshot=pathlib.Path('/Users/bill/projects/ensemble/.agents/chapters/code/ch05-student')
for n,h in files.items():assert sha(snapshot/n)==h,n
(P/'prior-preservation.json').write_text(json.dumps({'passed':True,'pinned_report_and_indexes':pins,'evidence_sets':rows,'original_snapshot_files':len(files),'originals_rewritten':False},indent=2)+'\n');print(rows,len(files))
