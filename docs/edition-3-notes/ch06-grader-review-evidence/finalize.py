from pathlib import Path
import json,hashlib,subprocess
E=Path(__file__).resolve().parent;R=E.parents[2];O=R/'docs/edition-3-notes/ch06-entry-evidence';S=Path((E/'scratch-path.txt').read_text().strip())
def h(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def git(*args,cwd=R):return subprocess.check_output(['git','--no-pager',*args],cwd=cwd).decode().strip()
def manifest(root):return {str(p.relative_to(root)):h(p) for p in sorted(root.rglob('*')) if p.is_file() and '.git' not in p.relative_to(root).parts}
pins=json.loads((E/'candidate-pins.json').read_text())
for name,m in pins['files'].items():assert h(name)==m['sha256'],name
for name,m in json.loads((O/'inputs.json').read_text()).items():
 if name!='cmd/grade-edition3/main.go':assert h(R/name)==m['sha256'],name
assert subprocess.check_output(['git','--no-pager','diff','HEAD','--'],cwd=R)==(E/'initial-tracked.diff').read_bytes()
entry_names={'SHA256SUMS'}
for line in (O/'SHA256SUMS').read_text().splitlines():
 a,b=line.split('  ',1);assert h(O/b)==a,b;entry_names.add(b)
assert {str(p.relative_to(O)) for p in O.rglob('*') if p.is_file()}==entry_names
assert h(O.parent/'ch06-entry-audit.md')=='a47915de7ce6f9d317f70efe84208d32cbc5a465ceee80026e55a1bb7577f7d6'
H=Path('/tmp/ensemble-ch06-entry-7_7zbtqs/solutions/ch06');assert manifest(H)==json.loads((O/'historical-source-sha256.json').read_text())
ST=Path('/Users/bill/projects/ensemble-edition-3-student');assert manifest(ST)==json.loads((O/'student-source-sha256.json').read_text())
assert git('status','--porcelain',cwd=ST)==''
assert git('rev-parse','HEAD',cwd=ST)=='f60ebf3aad2ece4251face9d95c106dd667abcd6'
assert git('rev-parse','HEAD^{tree}',cwd=ST)=='3de5d0d24dd9f593ec017fceccfbff1b209c5d57'
assert git('rev-parse','edition-3-ch05-r1^{commit}',cwd=ST)=='f60ebf3aad2ece4251face9d95c106dd667abcd6'
assert git('rev-parse','HEAD')=='b67386ab80cf2afe7985db5a32ab99ec2311f090'
assert git('rev-parse','HEAD:solutions/edition-3')=='4ff207e66df6d49cbdea462df2f6753d421779d8'
B=Path('/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble');assert h(B)=='e51138885a359eaf801258a52bc21f28ca41691683256647f8c9911be410b233'
commands=[];results={}
for p in sorted(E.glob('*.meta.json')):
 d=json.loads(p.read_text());label=p.name.removesuffix('.meta.json')
 for suffix in ['stdout','stderr']:assert h(E/(label+'.'+suffix))==d[suffix+'_sha256'],(label,suffix)
 commands.append(d)
 try:r=json.loads((E/(label+'.stdout')).read_text())
 except ValueError:continue
 if isinstance(r,dict) and 'checks' in r:results[label]={'score':r['score'],'max_score':r['max_score'],'passed':r['passed'],'exit':d['exit'],'seconds':d['duration_seconds'],'failed_ids':[c['id'] for c in r['checks'] if not c['passed']]}
assert len(results)==7,results
assert json.loads((E/'measurement-sequence.meta.json').read_text())['exit']==0
assert not (S/'launched').exists()
(E/'commands.json').write_text(json.dumps(commands,indent=2)+'\n');(E/'results.json').write_text(json.dumps(results,indent=2)+'\n')
(E/'source-binary-pins.json').write_text(json.dumps({'candidate':pins,'historical_binary':{'path':'/tmp/ensemble-ch06-entry-7_7zbtqs/positive','sha256':h('/tmp/ensemble-ch06-entry-7_7zbtqs/positive')},'accepted_f60_binary':{'path':str(B),'sha256':h(B)},'sentinel':{'path':str(S/'sentinel'),'sha256':h(S/'sentinel'),'source_sha256':h(S/'sentinel.go')},'historical_source_manifest_sha256':h(O/'historical-source-sha256.json'),'student_source_manifest_sha256':h(O/'student-source-sha256.json')},indent=2)+'\n')
(E/'final-verification.json').write_text(json.dumps({'main_head':git('rev-parse','HEAD'),'main_tracked_diff_unchanged_from_root_candidate':True,'protected_FINAL_tree':git('rev-parse','HEAD:solutions/edition-3'),'student_head':git('rev-parse','HEAD',cwd=ST),'student_tree':git('rev-parse','HEAD^{tree}',cwd=ST),'student_tag_commit':git('rev-parse','edition-3-ch05-r1^{commit}',cwd=ST),'student_clean_and_all_source_bytes_unchanged':True,'historical_source_all_bytes_unchanged':True,'candidate_source_test_binary_pins_unchanged':True,'all_prerequisite_and_original_grader_inputs_unchanged':True,'entry_report_and_exact_98_file_evidence_set_unchanged':True,'all_command_output_hashes_verified':True,'minimum_recorded_free_bytes':min(d[k] for d in commands for k in ['free_bytes_before','free_bytes_after']),'sentinel_marker_absent_after_all_invalid_help_checks':True,'measurement_sequence_exit':0},indent=2)+'\n')
print(json.dumps(results,indent=2));print((E/'source-comparison.json').read_text());print((E/'final-verification.json').read_text())
