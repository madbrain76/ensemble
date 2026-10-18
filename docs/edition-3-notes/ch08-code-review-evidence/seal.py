import sys,json,hashlib,os,subprocess,shutil,datetime
from pathlib import Path
sys.dont_write_bytecode=True
R=Path('/Users/bill/projects/ensemble');N=R/'docs/edition-3-notes';E=N/'ch08-code-review-evidence';S=Path('/Users/bill/projects/ensemble-edition-3-student');report=N/'ch08-code-review.md'
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
M=json.loads((N/'ch08-code-review-input.json').read_text());V=json.loads((E/'final-verification.json').read_text());assert sha(N/'ch08-code-review-input.json')==V['input_manifest_sha256']
for x in M['inputs']+list(M['protected_tracked'].values()):
 p=Path(x['path']);assert p.stat().st_size==x['bytes'] and sha(p)==x['sha256'],p
for x in M['source_members']:
 p=S/x['path'];assert p.stat().st_size==x['bytes'] and sha(p)==x['sha256'],p
assert subprocess.check_output(['git','status','--porcelain=v1','--untracked-files=all'],cwd=S)==b''
C=json.loads((E/'controls-summary.json').read_text());assert len(C['commands'])==28
for r in C['commands']:
 assert not r['remaining_process_group'] and not r['timed_out']
 try:os.killpg(r['pid_pgid'],0);raise AssertionError(('live owned process group',r['pid_pgid']))
 except ProcessLookupError:pass
L=json.loads((E/'local-artifacts.json').read_text())
for x in L['files']:
 p=Path(x['path']);assert p.stat().st_size==x['bytes'] and sha(p)==x['sha256'],p
seal={'checked_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'review_verdict':'NEEDS CORRECTION — R1 checked-in inherited page fixture compatibility','candidate':M['student']['commit'],'student_clean':True,'pinned_inputs':len(M['inputs']),'source_files':len(M['source_members']),'protected_files':len(M['protected_tracked']),'all_pinned_hashes_match':True,'owned_terminal_command_receipts':len(C['commands']),'chapter1_full_lifetime_runs':1,'remaining_owned_process_groups':0,'local_retained_artifact_hashes_match':True,'source_guidance_old_reports_modified':False,'browser_native_paid_or_secret_work':False,'new_outputs_only':['docs/edition-3-notes/ch08-code-review.md','docs/edition-3-notes/ch08-code-review-evidence/'],'scratch_root':L['root'],'free_bytes':shutil.disk_usage(R).free,'original_page_failure_preserved':True,'expected_mutant_negatives_preserved':True}
(E/'final-seal.json').write_text(json.dumps(seal,indent=2)+'\n')
paths=[report]+sorted(p for p in E.rglob('*') if p.is_file() and p!=E/'inventory.json');entries=[]
for p in paths:
 assert not p.is_symlink(),p
 data=p.read_bytes();assert b'\0' not in data,('NUL',p);data.decode('utf-8')
 assert data[:4] not in (b'\x7fELF',b'\xcf\xfa\xed\xfe',b'\xfe\xed\xfa\xcf',b'\xca\xfe\xba\xbe',b'PK\x03\x04'),p
 rel=str(p.relative_to(R));assert rel=='docs/edition-3-notes/ch08-code-review.md' or rel.startswith('docs/edition-3-notes/ch08-code-review-evidence/')
 entries.append({'path':rel,'bytes':len(data),'sha256':hashlib.sha256(data).hexdigest(),'classification':'portable UTF-8 review report' if p==report else 'portable UTF-8 evidence/source/metadata; not a binary artifact'})
assert len({x['path'] for x in entries})==len(entries)
inv={'format':'ch08-initial-whole-code-review-v1','self_excluded':True,'files':entries,'file_count':len(entries),'total_bytes':sum(x['bytes'] for x in entries),'local_binary_and_process_artifacts':'separately indexed in local-artifacts.json; no payloads exported','existing_student_tar':'preserved/excluded/local-only exception; not copied here'}
(E/'inventory.json').write_text(json.dumps(inv,indent=2)+'\n')
for x in json.loads((E/'inventory.json').read_text())['files']:assert sha(R/x['path'])==x['sha256'] and (R/x['path']).stat().st_size==x['bytes']
print(json.dumps({'seal':seal,'portable_file_count':inv['file_count'],'portable_total_bytes':inv['total_bytes'],'report_sha256':sha(report),'inventory_sha256':sha(E/'inventory.json'),'verified_inventory_differences':0},indent=2))
