import sys;sys.dont_write_bytecode=True
from common import *
import datetime
C=json.loads((R/'docs/edition-3-notes/ch08-grader-evidence/candidate.json').read_text());M=json.loads((R/'docs/edition-3-notes/ch08-entry-manifest.json').read_text());S=Path('/Users/bill/projects/ensemble-edition-3-student');B=json.loads((E/'protected-baseline.json').read_text());differences=[]
for p,x in B.items():
 if not Path(p).is_file() or sha(p)!=x['sha256']:differences.append(p)
assert not differences,differences
verification={'protected_baseline_files':len(B),'protected_baseline_differences':differences,'candidate_sources':{},'candidate_binary':{'path':C['binary']['path'],'sha256':sha(C['binary']['path']),'matches':sha(C['binary']['path'])==C['binary']['sha256']}}
for p,x in C['sources'].items():verification['candidate_sources'][p]={'sha256':sha(R/p),'matches':sha(R/p)==x['sha256']};assert verification['candidate_sources'][p]['matches']
assert verification['candidate_binary']['matches']
for key in ['inputs','student_files','protected_tracked_changes']:
 bad=[]
 for p,x in M[key].items():
  path=Path(p) if Path(p).is_absolute() else (S/p if key=='student_files' else R/p);expected=C['sources'][p]['sha256'] if p in C['sources'] else x['sha256']
  if sha(path)!=expected:bad.append(str(path))
 verification[key]={'count':len(M[key]),'differences':bad,'candidate_override_paths':[p for p in M[key] if p in C['sources']]};assert not bad
verification['main_head']=subprocess.check_output(['git','rev-parse','HEAD'],cwd=R).decode().strip();verification['student_head']=subprocess.check_output(['git','rev-parse','HEAD'],cwd=S).decode().strip();verification['student_clean']=not subprocess.check_output(['git','status','--porcelain'],cwd=S);assert verification['main_head']==M['main_head'] and verification['student_head']==M['student_head'] and verification['student_clean']
receipts=[json.loads(p.read_text()) for p in E.glob('*.receipt.json')];live=[]
for x in receipts:
 assert not x['timed_out'] and not x['process_group_remaining']
 assert sha(E/(x['name']+'.stdout'))==x['stdout_sha256'] and sha(E/(x['name']+'.stderr'))==x['stderr_sha256']
 try:os.killpg(x['pid_pgid'],0);live.append(x['pid_pgid'])
 except ProcessLookupError:pass
assert not live,live
verification.update({'receipt_count':len(receipts),'remaining_owned_groups':live,'all_control_jobs_ended':True,'free_bytes':shutil.disk_usage(R).free,'minimum_free_bytes_after_control':min(x['free_bytes_after'] for x in receipts),'verified_at':datetime.datetime.now(datetime.timezone.utc).isoformat()})
(E/'verification.json').write_text(json.dumps(verification,indent=2)+'\n')
root=R/'docs/edition-3-notes/ch08-grader-evidence';commands=json.loads((root/'commands.json').read_text());carry=[]
for i in [2,3,4]:
 c=commands[i];logs={}
 for kind in ['stdout','stderr']:
  p=root/('check-'+str(i+1)+'.'+kind+'.log');h=sha(p);assert h==c[kind+'_sha256'];logs[kind]={'path':str(p.relative_to(R)),'sha256':h,'bytes':p.stat().st_size}
 assert c['exit']==0;carry.append({'provenance':'root carried, not independently rerun','command':c,'logs':logs})
(E/'carried-root-controls.json').write_text(json.dumps({'candidate_manifest_path':str((root/'candidate.json').relative_to(R)),'candidate_manifest_sha256':sha(root/'candidate.json'),'commands_path':str((root/'commands.json').relative_to(R)),'commands_sha256':sha(root/'commands.json'),'controls':carry,'root_preparation_failures_path':str((root/'preparation-failures.md').relative_to(R)),'root_preparation_failures_sha256':sha(root/'preparation-failures.md'),'limits':'Root race/focused settings, vet, build carried only. No independent race/vet rerun; standalone legacy1/4/5/6 full gates not rerun or claimed carried.'},indent=2)+'\n')
art=[]
for p in sorted(T.rglob('*')):
 if not p.is_file():continue
 b=p.read_bytes();kind='local-only-source-or-fixture'
 if p.name=='status-before.z':kind='local-only-NUL-git-status'
 elif p in [T/'bin/candidate-settings-driver',T/'sentinel/marker',T/'sentinel/shim/go']:kind='local-only-compiled-executable'
 art.append({'path':str(p),'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest(),'classification':kind,'symlink_target':str(p.readlink()) if p.is_symlink() else None})
(E/'local-artifacts.json').write_text(json.dumps({'scratch':str(T),'retained':True,'files':art,'file_count':len(art),'total_bytes_including_symlink_target_accounting':sum(x['bytes'] for x in art),'portable_payloads_excluded':['compiled executables','NUL git status'],'note':'Hash metadata only. Internally built grader targets were cleaned by existing grade.Build; exact compile hashes remain in builds.jsonl. Prior retained/reused artifacts are indexed separately in pins-before/settings-driver-provenance.'},indent=2)+'\n')
(E/'read-lookup-notes.json').write_text(json.dumps({'missing_read_candidates':['ch08-grader-evidence/receipt.json','ch08-grader-evidence/controls/full.stdout','ch08-grader-evidence/controls/focused.stdout','ch08-entry-reconciliation-evidence/sources/main.go','ch08-grader-evidence/check-0.stdout.log'],'qualification':'Read-only guessed paths were absent, not builds/tests/control failures. Actual receipt/log/source paths were then resolved. The fixtures were copied from hash-verified prior scratch. See preparation-failures.json for setup, result-wrapper, and sealing errors.'},indent=2)+'\n')
print(json.dumps({'verified':True,'baseline':len(B),'receipts':len(receipts),'remaining_groups':live,'local_artifacts':len(art),'local_bytes':sum(x['bytes'] for x in art),'candidate':C,'free_bytes':verification['free_bytes']},indent=2))
