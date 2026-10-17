"""Local evidence verification only: no provider request and no credential output."""
import hashlib,json,pathlib,subprocess
base=pathlib.Path(__file__).resolve().parent
repo=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
def sha(data):return hashlib.sha256(data).hexdigest()
source=[]
for entry in subprocess.check_output(['git','ls-files','-s','-z'],cwd=repo).decode().split('\0'):
 if not entry:continue
 meta,path=entry.split('\t',1);mode,blob,stage=meta.split();assert stage=='0'
 data=subprocess.check_output(['git','show',':'+path],cwd=repo)
 source.append({'path':path,'mode':mode,'git_blob':blob,'sha256':sha(data),'bytes':len(data)})
source_info={'tree':subprocess.check_output(['git','write-tree'],cwd=repo,text=True).strip(),'files':source}
(base/'source-manifest.json').write_text(json.dumps(source_info,indent=2)+'\n')
local=[]
for path in [pathlib.Path('/tmp/ch04-core'),pathlib.Path('/tmp/ch04-consumer')]+sorted(pathlib.Path('/tmp/ch04-live').rglob('*')):
 if path.is_file():
  content=path.read_bytes();local.append({'path':str(path),'bytes':len(content),'sha256':sha(content),'classification':'temporary local artifact; not committed or included as a portable evidence blob'})
(base/'local-artifacts-manifest.json').write_text(json.dumps(local,indent=2)+'\n')
# Read only the authorized store here, never write or print the field values.
settings=json.loads((pathlib.Path.home()/'.cr/settings.json').read_text())
keys=[settings.get(n) for n in ['directClaudeAPIKey','directOpenAIAPIKey','directGeminiAPIKey']]
assert all(isinstance(k,str) and len(k)>12 for k in keys),'authorized credential unavailable for scan'
keys=[k.encode() for k in keys];del settings
findings=[];count=0
paths=list(p for p in base.rglob('*') if p.is_file())+[base.parent/'ch04-student-review.md']
paths += [repo/s['path'] for s in source]
paths += [pathlib.Path(v['path']) for v in local]
for path in paths:
 data=path.read_bytes();count+=1
 if any(key in data for key in keys):findings.append({'path':str(path),'problem':'configured credential value found'})
opaque_findings=[]
def check_opaque(obj,where):
 if isinstance(obj,dict):
  for key,value in obj.items():
   if key.lower()=='opaque' and value not in (None,'[opaque omitted]'):opaque_findings.append(where)
   else:check_opaque(value,where)
 elif isinstance(obj,list):
  for value in obj:check_opaque(value,where)
for path in base.glob('*.events.jsonl'):
 for n,line in enumerate(path.read_text().splitlines(),1):check_opaque(json.loads(line),str(path.name)+':'+str(n))
result={'files_scanned':count,'credential_fields_checked':3,'credential_findings':findings,'unredacted_opaque_findings':opaque_findings,'scope':'own evidence/report, all staged source files, and own temporary binaries/live artifacts only','no_provider_calls':True}
(base/'redaction-verification.json').write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps(result));assert not findings and not opaque_findings
