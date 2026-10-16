"""Scoped independent R1 re-review; immutable archives and disposable mutation only."""
from pathlib import Path
import difflib, hashlib, json, shutil, subprocess, tempfile
M = Path('/Users/bill/projects/ensemble')
S = Path('/Users/bill/projects/ensemble-edition-3-student')
E = Path(__file__).resolve().parent
OLD = 'd2baff64055ac33702a91c6d6ee5a4246e84c0c9'
NEW = '1aad3975842d502b9144b7b9bddfd30ca7309db7'
MAIN = '7f3d069420112c2136bc95d1d551e12fad6964c1'
def git(root, *args):
 return subprocess.check_output(['git', '--no-pager', *args], cwd=root)
def sha(data): return hashlib.sha256(data).hexdigest()
assert git(S, 'rev-parse', 'HEAD').decode().strip() == NEW
assert not git(S, 'status', '--porcelain')
assert git(M, 'rev-parse', 'HEAD').decode().strip() == MAIN
paths = git(S, 'ls-tree', '-r', '--name-only', NEW).decode().splitlines()
assert paths == git(S, 'ls-tree', '-r', '--name-only', OLD).decode().splitlines()
old = {p: git(S, 'show', OLD+':'+p) for p in paths}
new = {p: git(S, 'show', NEW+':'+p) for p in paths}
assert [p for p in paths if old[p] != new[p]] == ['ensemble_test.go']
a, b = old['ensemble_test.go'].splitlines(), new['ensemble_test.go'].splitlines()
assert len(a) == len(b)
assert [i+1 for i, (x,y) in enumerate(zip(a,b)) if x != y] == [54,55]
fixture = json.loads(b[54].decode().split('`')[1])
assert fixture == {'content':[{'type':'text','text':'hello '},{'type':'text','text':'世界'}], 'usage':{'input_tokens':17,'output_tokens':9}}
assert b[53].decode().strip() == '// Two text blocks must be concatenated exactly, without added separators.'
(E/'correction.patch').write_bytes(git(S, 'diff', OLD, NEW, '--'))
protected = [M/'docs/edition-3-notes/ch01-code-review-r1.md', M/'docs/edition-3-notes/ch01-student-review.md', M/'docs/edition-3-notes/ch01-candidate.json', M/'docs/edition-3-notes/ch01-student-correction-1.md']
for folder in ['ch01-code-review-evidence', 'ch01-student-evidence', 'ch01-student-correction-1-evidence']:
 protected += [p for p in (M/'docs/edition-3-notes'/folder).rglob('*') if p.is_file()]
preserved = {str(p.relative_to(M)): sha(p.read_bytes()) for p in protected}
prior = json.loads((M/'docs/edition-3-notes/ch01-code-review-evidence/provenance.json').read_text())
for p, digest in prior['sha256'].items(): assert sha((M/p).read_bytes()) == digest, p
T = Path(tempfile.mkdtemp(prefix='e3-ch01-rereview1.'))
subprocess.run(['tar','-x','-C',str(T)], input=git(S,'archive',NEW), check=True)
(E/'scratch-path.txt').write_text(str(T)+'\n')
proof = {'candidate':NEW, 'previous_candidate':OLD, 'book_grader':MAIN, 'changed_paths':['ensemble_test.go'], 'only_changed_lines':[54,55], 'same_path_set':True, 'all_other_test_lines_identical':True, 'all_production_and_module_bytes_identical':True, 'files':{p:{'previous_sha256':sha(old[p]),'candidate_sha256':sha(new[p]),'identical':old[p]==new[p]} for p in paths}, 'preserved_artifact_sha256':preserved, 'pinned_input_sha256':prior['sha256'], 'scratch_path':str(T)}
(E/'provenance.json').write_text(json.dumps(proof,indent=2)+'\n')
def run(name, args, cwd=T):
 with (E/(name+'.txt')).open('w') as log:
  log.write('COMMAND '+repr(args)+'\nCWD '+str(cwd)+'\n'); log.flush()
  result = subprocess.run(args, cwd=cwd, stdout=log, stderr=subprocess.STDOUT)
 (E/(name+'.exit')).write_text(str(result.returncode)+'\n')
 return result.returncode
assert run('go-version',['go','version']) == 0
assert run('build',['go','build','-o',str(T/'ch01'),'./cmd/ch01']) == 0
assert run('vet',['go','vet','./...']) == 0
focus = ['go','test','-count=1','-v','-run','^TestConversation$','.']
assert run('focused-control',focus) == 0
with tempfile.TemporaryDirectory(prefix='e3-ch01-r1-firstblock.') as mutant:
 D = Path(mutant)
 shutil.copytree(T,D,dirs_exist_ok=True,ignore=lambda p,n: ['ch01'] if Path(p)==T else [])
 engine = D/'internal/engine/engine.go'
 original = engine.read_bytes()
 anchor = b'text.WriteString(block.Text)'
 assert original.count(anchor)==1
 engine.write_bytes(original.replace(anchor,anchor+b'\n break'))
 status = run('first-block-mutant',focus,D)
 log = (E/'first-block-mutant.txt').read_text()
 assert status == 1 and 'all text blocks: "hello " <nil>' in log and '--- FAIL: TestConversation' in log and '[build failed]' not in log
 engine.write_bytes(original)
 assert run('focused-restored',focus,D) == 0
 for p in paths: assert (D/p).read_bytes() == new[p], p
assert run('tests',['go','test','-race','-count=1','-v','./...']) == 0
assert run('format',['gofmt','-l',*[str(T/p) for p in paths if p.endswith('.go')]]) == 0
format_lines = (E/'format.txt').read_text().splitlines()
assert len(format_lines) == 2, format_lines
for p in paths:
 assert (T/p).read_bytes() == new[p], p
 assert (S/p).read_bytes() == new[p], p
for p,digest in preserved.items(): assert sha((M/p).read_bytes()) == digest, p
for p,digest in prior['sha256'].items(): assert sha((M/p).read_bytes()) == digest, p
assert git(S,'rev-parse','HEAD').decode().strip() == NEW
assert not git(S,'status','--porcelain')
assert git(M,'rev-parse','HEAD').decode().strip() == MAIN
proof['final_checks_passed'] = True
proof['built_binary_sha256'] = sha((T/'ch01').read_bytes())
(E/'provenance.json').write_text(json.dumps(proof,indent=2)+'\n')
(E/'hygiene.txt').write_text('Exact candidate remains clean.\nAll archived tracked files match the pinned candidate.\nAll production/module bytes match the initially reviewed revision.\nOnly ensemble_test.go lines 54-55 differ.\nFirst-block mutation failed behaviorally; restored control passed.\nBuild/vet/full race suite/format checks passed.\nAll protected original and student correction artifacts remain unchanged.\nPinned chapter/grader hashes and HEAD unchanged.\nNo grader rerun, provider call or credential access performed in this round.\n')
print('Scoped R1 verification completed successfully; evidence saved.')
