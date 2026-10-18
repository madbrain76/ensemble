import sys;sys.dont_write_bytecode=True
from common import *
import datetime
report=R/'docs/edition-3-notes/ch08-grader-review.md';I=E/'inventory.json';F=E/'final-seal.json'
B=json.loads((E/'protected-baseline.json').read_text());assert all(Path(p).is_file() and sha(p)==x['sha256'] for p,x in B.items())
C=json.loads((E/'pins-before.json').read_text())['candidate'];assert all(sha(R/p)==x['sha256'] for p,x in C['sources'].items());assert sha(C['binary']['path'])==C['binary']['sha256']
receipts=[json.loads(p.read_text()) for p in E.glob('*.receipt.json')];live=[]
for x in receipts:
 assert not x['timed_out'] and not x['process_group_remaining'];assert sha(E/(x['name']+'.stdout'))==x['stdout_sha256'] and sha(E/(x['name']+'.stderr'))==x['stderr_sha256']
 try:os.killpg(x['pid_pgid'],0);live.append(x['pid_pgid'])
 except ProcessLookupError:pass
assert not live
A=json.loads((E/'local-artifacts.json').read_text());assert all(Path(x['path']).is_file() and sha(x['path'])==x['sha256'] for x in A['files'])
I.write_text('{}\n');F.write_text('{}\n')
def status(b):return {s[3:].decode():s[:2].decode() for s in b.split(b'\0') if s}
before=status((T/'status-before.z').read_bytes());after=status(subprocess.check_output(['git','status','--porcelain=v1','-z','--untracked-files=all'],cwd=R));changed=[p for p in sorted(set(before)|set(after)) if before.get(p)!=after.get(p)];bad=[p for p in changed if not (p=='docs/edition-3-notes/ch08-grader-review.md' or p.startswith('docs/edition-3-notes/ch08-grader-review-evidence/'))];assert not bad,bad
paths=sorted([report]+[p for p in E.rglob('*') if p.is_file() and p!=I]);unsafe=[]
for p in paths:
 if p.is_symlink():unsafe.append(str(p));continue
 b=p.read_bytes()
 try:b.decode('utf-8')
 except UnicodeDecodeError:unsafe.append(str(p))
 if b'\0' in b or p.suffix in ['.pyc','.zip','.tar','.gz'] or '__pycache__' in p.parts:unsafe.append(str(p))
assert not unsafe,unsafe
seal={'sealed_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'verdict':'PASS for exact pinned candidate only','protected_baseline_files':len(B),'protected_differences':[],'candidate_source_and_binary_pins_match':True,'portable_output_scope':'only new ch08-grader-review.md and ch08-grader-review-evidence/','changed_status_paths':changed,'out_of_scope_status_changes':bad,'owned_command_receipts':len(receipts),'all_owned_control_jobs_ended':True,'remaining_owned_process_groups':live,'scratch_retained':str(T),'local_artifact_hashes_reverified':len(A['files']),'portable_UTF8_no_NUL_no_binary_archive_pycache':True,'portable_path_count_excluding_inventory':len(paths),'free_bytes':shutil.disk_usage(R).free,'minimum_free_bytes_after_control':min(x['free_bytes_after'] for x in receipts),'inventory_self_excluded':True}
F.write_text(json.dumps(seal,indent=2)+'\n')
files=[{'path':str(p.relative_to(R)),'bytes':p.stat().st_size,'sha256':sha(p),'classification':'review report' if p==report else 'portable UTF-8 evidence'} for p in paths];data={'format':'sha256-size-portable-review-inventory-v1','self_excluded':str(I.relative_to(R)),'files':files,'file_count':len(files),'total_bytes':sum(x['bytes'] for x in files)};I.write_text(json.dumps(data,indent=2)+'\n')
assert all((R/x['path']).stat().st_size==x['bytes'] and sha(R/x['path'])==x['sha256'] for x in files)
print(json.dumps({'verdict':seal['verdict'],'files':len(files),'total_bytes':data['total_bytes'],'inventory_sha256':sha(I),'remaining_owned_groups':live,'out_of_scope_status_changes':bad,'protected_files_unchanged':len(B),'free_bytes':seal['free_bytes'],'local_retained_artifacts':len(A['files'])},indent=2))
