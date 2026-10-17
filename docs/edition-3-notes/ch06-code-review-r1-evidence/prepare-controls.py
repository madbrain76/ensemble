from pathlib import Path
import json,shutil,hashlib
E=Path(__file__).resolve().parent;R=E.parents[2];S=Path((E/'scratch-path.txt').read_text().strip());old=R/'docs/edition-3-notes/ch06-code-review-evidence';expected=json.loads((E/'source-hashes.json').read_text())
D=S/'control-source';D.mkdir(exist_ok=True)
for rel,h in expected.items():
 p=D/rel;p.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(S/'source'/rel,p);assert hashlib.sha256(p.read_bytes()).hexdigest()==h
M=S/'public-checks';M.mkdir(exist_ok=True)
for name in ['public_checks_test.go','batch_interrupt_test.go']:
 data=(old/name).read_bytes();(M/name).write_bytes(data);(E/name).write_bytes(data)
g=(S/'source/go.mod').read_text().replace('module ensemble','module reviewchecks',1)+'\nrequire ensemble v0.0.0\nreplace ensemble => ../control-source\n';(M/'go.mod').write_text(g);(M/'go.sum').write_bytes((S/'source/go.sum').read_bytes());(E/'public-module.go.mod.txt').write_text(g)
print('prepared',M,'production copy files',len(expected),'old probe bytes unchanged')
