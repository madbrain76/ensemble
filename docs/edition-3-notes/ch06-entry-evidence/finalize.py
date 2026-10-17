#!/usr/bin/env python3
"""Read-only provenance verification, then write evidence summaries (no tests)."""
import hashlib,json,pathlib,subprocess
E=pathlib.Path(__file__).resolve().parent;R=E.parents[2];S=pathlib.Path((E/'scratch-path.txt').read_text().strip());T=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
def cmd(args,cwd=R):return subprocess.check_output(args,cwd=cwd,text=True).strip()
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def manifest(p):return {str(f.relative_to(p)):sha(f) for f in sorted(p.rglob('*')) if f.is_file() and '.git' not in f.parts}
main=cmd(['git','rev-parse','HEAD']);assert main=='b67386ab80cf2afe7985db5a32ab99ec2311f090'
assert cmd(['git','diff','HEAD','--'])==''
assert cmd(['git','status','--short'],T)==''
assert cmd(['git','rev-parse','HEAD'],T)=='f60ebf3aad2ece4251face9d95c106dd667abcd6'
assert manifest(T)==json.loads((E/'student-source-sha256.json').read_text())
base=manifest(S/'solutions/ch06');assert base==json.loads((E/'historical-source-sha256.json').read_text())
for rel,h in base.items():
 data=subprocess.check_output(['git','show',main+':solutions/ch06/'+rel],cwd=R)
 assert hashlib.sha256(data).hexdigest()==h,rel
for n in ['disable','correlation']:assert manifest(S/('mutant-'+n))==base
inputs=json.loads((E/'inputs.json').read_text())
for p,m in inputs.items():assert sha(R/p)==m['sha256'],p
bins={str(p):sha(p) for p in sorted(S.iterdir()) if p.is_file()}
bins['/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble']=sha(pathlib.Path('/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble'))
(E/'final-binary-sha256.json').write_text(json.dumps(bins,indent=2)+'\n')
reports=[]
for f in sorted(E.glob('*.meta.json')):
 m=json.loads(f.read_text());label=f.name.removesuffix('.meta.json');row={'label':label,'exit':m['exit'],'duration_seconds':m['duration_seconds']}
 try:
  report=json.loads((E/(label+'.stdout')).read_text())
  if isinstance(report,dict) and 'score' in report:
   row.update({k:report[k] for k in ['score','max_score','passed']});row['failed_checks']=[c['id'] for c in report['checks'] if not c['passed']]
 except (ValueError,FileNotFoundError):pass
 reports.append(row)
(E/'results.json').write_text(json.dumps(reports,indent=2)+'\n')
(E/'commands.txt').write_text('\n\n'.join(json.dumps({'label':f.name.removesuffix('.meta.json'),**json.loads(f.read_text())},indent=2) for f in sorted(E.glob('*.meta.json')))+'\n')
s=R.stat();disk=__import__('os').statvfs(R)
state={'main_head':main,'main_tracked_diff_empty':True,'protected_final_tree':cmd(['git','rev-parse','HEAD:solutions/edition-3']),'student_head':cmd(['git','rev-parse','HEAD'],T),'student_tree':cmd(['git','rev-parse','HEAD^{tree}'],T),'student_tag_commit':cmd(['git','rev-parse','refs/tags/edition-3-ch05-r1^{}'],T),'student_status_clean':True,'student_source_manifest_unchanged':True,'historical_source_equals_pinned_git_bytes':True,'both_restored_sources_equal_historical':True,'prerequisite_and_grader_input_hashes_unchanged':True,'absent_ch06_exercise_exists':(S/'absent-ch06-stream-exercise').exists(),'free_bytes':disk.f_bavail*disk.f_frsize,'scratch':str(S),'main_untracked':cmd(['git','ls-files','--others','--exclude-standard'])}
assert not state['absent_ch06_exercise_exists']
(E/'final-verification.json').write_text(json.dumps(state,indent=2)+'\n')
(E/'buildinfo.txt').write_text('\n'.join(cmd(['go','version','-m',p]) for p in bins)+'\n')
print(json.dumps(state,indent=2));print(json.dumps(bins,indent=2))
