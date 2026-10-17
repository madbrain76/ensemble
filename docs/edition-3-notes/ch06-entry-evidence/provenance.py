#!/usr/bin/env python3
import pathlib,hashlib,json,subprocess
E=pathlib.Path(__file__).resolve().parent; S=pathlib.Path((E/'scratch-path.txt').read_text().strip()); R=pathlib.Path('/Users/bill/projects/ensemble')
def manifest(p):
 return {str(f.relative_to(p)):hashlib.sha256(f.read_bytes()).hexdigest() for f in sorted(p.rglob('*')) if f.is_file() and '.git' not in f.parts}
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
for label,p in [('historical',S/'solutions/ch06'),('student',pathlib.Path('/Users/bill/projects/ensemble-edition-3-student'))]:
 (E/(label+'-source-sha256.json')).write_text(json.dumps(manifest(p),indent=2)+'\n')
(E/'binary-sha256.json').write_text(json.dumps({str(p):sha(p) for p in [S/'api',S/'grade-edition2',S/'positive',pathlib.Path('/tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble')]},indent=2)+'\n')
(E/'git-identities.txt').write_text(subprocess.check_output(['git','rev-parse','HEAD','HEAD:solutions/ch06','HEAD:solutions/edition-3'],cwd=R,text=True)+subprocess.check_output(['git','rev-parse','HEAD','HEAD^{tree}','refs/tags/edition-3-ch05-r1^{}'],cwd='/Users/bill/projects/ensemble-edition-3-student',text=True))
