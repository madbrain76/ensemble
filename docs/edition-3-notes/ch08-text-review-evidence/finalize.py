from pathlib import Path
import json,hashlib,subprocess,re,datetime
E=Path(__file__).parent
R=E.parent/'ch08-text-review.md'
I=E/'inventory.json'
assert not I.exists(), 'Seal already exists; do not rewrite sealed evidence'
def ident(p):
 b=p.read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def save(name,obj):
 (E/name).write_text(json.dumps(obj,indent=2)+'\n')
def command(name,argv):
 p=subprocess.run(argv,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
 (E/name).write_bytes(p.stdout)
 row={'argv':argv,'exit':p.returncode,'output_file':name,**ident(E/name)}
 with (E/'commands.jsonl').open('a') as f:f.write(json.dumps(row)+'\n')
 return row,p.stdout
pins=json.loads((E/'input-pins.json').read_text())['inputs']
pincheck=[{'path':p,**ident(Path(p)),'matches':ident(Path(p))==v} for p,v in pins.items()]
assert all(x['matches'] for x in pincheck)
reads=[json.loads(x) for x in (E/'reads.jsonl').read_text().splitlines()]
extra=[]
for p in ['cmd/lintprose/main.go','internal/gradeedition2/ch05_media.go']:
 r=next(x for x in reversed(reads) if x['path']==p)
 now=ident(Path(p));extra.append({'path':p,**now,'matches_read':now=={'bytes':r['file_bytes'],'sha256':r['file_sha256']}})
assert all(x['matches_read'] for x in extra)
ws,ws_output=command('report-diff-check.log',['git','--no-pager','diff','--no-index','--check','--','/dev/null',str(R)])
# --no-index reports differing files as 1; --check must emit no diagnostics.
assert ws['exit'] in (0,1) and not ws_output
assert all(line==line.rstrip(' \t') for line in R.read_text().splitlines())
ws['interpretation']='No whitespace diagnostics; exit 1 is expected for a new file compared with /dev/null. Independent trailing-space/tab check also passes.'
scope=subprocess.run(['git','--no-pager','status','--porcelain=v1','--untracked-files=all'],stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
assert scope.returncode==0
before=set((E/'scope-initial.txt').read_bytes().splitlines());after=set(scope.stdout.splitlines())
added=sorted(after-before);removed=sorted(before-after)
allowed=['docs/edition-3-notes/ch08-text-review.md','docs/edition-3-notes/ch08-text-review-evidence/']
assert not removed
assert all(any(x.decode()[3:]==a or (a.endswith('/') and x.decode()[3:].startswith(a)) for a in allowed) for x in added),added
with (E/'commands.jsonl').open('a') as f:
 f.write(json.dumps({'argv':['git','--no-pager','status','--porcelain=v1','--untracked-files=all'],'exit':scope.returncode,'output_bytes':len(scope.stdout),'output_sha256':hashlib.sha256(scope.stdout).hexdigest(),'retained':'path delta in final-checks.json; existing initial and intermediate full status receipts retained'})+'\n')
links=[]
for target in re.findall(r'\]\(([^)]+)\)',R.read_text()):
 if re.match(r'[a-z]+:|#',target):continue
 p=R.parent/target.split('#')[0]
 links.append({'target':target,'exists_before_seal':p.exists(),'planned_finalization_output':p in [I,E/'final-checks.json',E/'final-submission.json']})
assert all(x['exists_before_seal'] or x['planned_finalization_output'] for x in links)
coverage=json.loads((E/'read-coverage.json').read_text());assert all(x['full_read'] for x in coverage['complete_files'])
for p in E.rglob('*'):
 if p.is_file():
  b=p.read_bytes();b.decode('utf-8');assert b'\0' not in b,p
save('final-submission.json',{'verdict':'PASS','required_revisions':[],'report':{'path':str(R),**ident(R)},'chapter':{'path':'book/edition-3/chapter-08.md',**ident(Path('book/edition-3/chapter-08.md'))},'evidence_inventory':str(I),'inventory_self_hash':'excluded to avoid circularity; inventory generated last and validated read-only','single_formal_review':True,'children_spawned':0,'runtime_exercises':0,'chapter_edits':0,'limitations':['Existing runtime/browser/provider results reviewed as records, not rerun.','No human listening or direct reviewer pixel observation.','Path-status preservation is not publication-content preservation.','Setup/read recovery failures retained in failures.md.'],'completion':'Saved final submission; no further chapter edit/review requested.'})
save('final-checks.json',{'stage':'All checks and submission saved before inventory generation; post-seal read-only validation is printed by finalize.py, not written into sealed members.','input_pins':pincheck,'extra_source_pins':extra,'all_input_pins_match':True,'complete_read_files':len(coverage['complete_files']),'report':{'path':str(R),**ident(R)},'report_diff_check':ws,'report_links':links,'scope':{'added_status_lines':[x.decode() for x in added],'removed_status_lines':[],'only_allowed_output_additions':True,'qualification':'Name/status comparison only; no assertion of excluded publication content preservation. Final checks, saved submission and inventory are themselves new permitted outputs.'},'utf8_no_nul_preflight':True,'children_spawned':0,'runtime_services_started':0,'known_setup_failures':'failures.md'})
files=[]
for p in sorted(E.rglob('*')):
 if p.is_file() and p!=I:
  b=p.read_bytes();b.decode('utf-8');assert b'\0' not in b,p
  files.append({'path':str(p.relative_to(E)),**ident(p)})
external=[{'path':str(R),**ident(R)}]
manifest={'schema':'ch08-text-review-evidence-v1','verdict':'PASS','generated_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'inventory_self_exclusion':{'path':'inventory.json','reason':'Self-excluded to avoid circular hashing; no evidence/report member is changed after this inventory is written.'},'files':files,'external_outputs':external}
# LAST WRITE. Everything below is a read-only check.
I.write_text(json.dumps(manifest,indent=2)+'\n')
actual={str(p.relative_to(E)) for p in E.rglob('*') if p.is_file() and p!=I}
assert actual=={x['path'] for x in manifest['files']}
for x in manifest['files']:
 assert ident(E/x['path'])=={'bytes':x['bytes'],'sha256':x['sha256']}
for x in external:assert ident(Path(x['path']))=={'bytes':x['bytes'],'sha256':x['sha256']}
for x in links:assert (R.parent/x['target'].split('#')[0]).exists()
for p,want in pins.items():assert ident(Path(p))==want
for p in [R,I]+[E/x['path'] for x in manifest['files']]:
 b=p.read_bytes();b.decode('utf-8');assert b'\0' not in b
print(json.dumps({'post_seal_read_only_validation':'PASS','evidence_members':len(files),'external_outputs':len(external),'report':{'path':str(R),**ident(R)},'inventory':{'path':str(I),**ident(I)},'saved_submission':{'path':str(E/'final-submission.json'),**ident(E/'final-submission.json')},'input_pins_unchanged':len(pins),'all_report_links_exist':True,'no_writes_after_inventory':True},indent=2))
