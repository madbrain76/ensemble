from pathlib import Path
import json,hashlib,subprocess,difflib,re,datetime
E=Path(__file__).parent; A=Path('docs/edition-3-notes/ch08-author-evidence'); C=Path('book/edition-3/chapter-08.md')
def ident(p):
 b=p.read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def save(n,v): (E/n).write_text(json.dumps(v,indent=2)+'\n')
def run(n,argv):
 r=subprocess.run(argv,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);(E/n).write_bytes(r.stdout)
 rec={'argv':argv,'exit':r.returncode,'output_file':n,'output_identity':ident(E/n)}
 with (E/'commands.jsonl').open('a') as f:f.write(json.dumps(rec)+'\n')
 print(json.dumps(rec)); print(r.stdout.decode());return r.returncode
inv=json.loads((A/'inventory.json').read_text());rows=[]
for group in ['files','external_outputs']:
 for row in inv[group]:
  p=A/row['path'] if group=='files' else Path(row['path']); b=p.read_bytes();b.decode('utf-8')
  rows.append({'path':str(p),**ident(p),'matches':all(ident(p)[k]==row[k] for k in ('bytes','sha256')),'utf8':True,'nul_bytes':b.count(b'\0')})
actual={str(p.relative_to(A)) for p in A.rglob('*') if p.is_file() and p.name!='inventory.json'}
expected={x['path'] for x in inv['files']}
save('author-seal-verification.json',{'inventory_identity':ident(A/'inventory.json'),'members':rows,'unlisted':sorted(actual-expected),'missing':sorted(expected-actual)})
assert actual==expected and all(x['matches'] and not x['nul_bytes'] for x in rows)
pins=json.loads((E/'input-pins.json').read_text());checks=[]
for p,w in pins['inputs'].items():checks.append({'path':p,**ident(Path(p)),'matches':ident(Path(p))==w})
save('pins-check.json',checks);assert all(x['matches'] for x in checks)
b=(A/'baseline/chapter-08.md').read_bytes(); f=C.read_bytes();assert f==(A/'final/chapter-08.md').read_bytes()
bt=b.decode();ft=f.decode();ops=difflib.SequenceMatcher(None,bt.splitlines(keepends=True),ft.splitlines(keepends=True)).get_opcodes()
changes=[{'kind':tag,'baseline_lines':[i+1,j],'final_lines':[k+1,l]} for tag,i,j,k,l in ops if tag!='equal']
diff=''.join(difflib.unified_diff(bt.splitlines(keepends=True),ft.splitlines(keepends=True),fromfile='baseline/book/edition-3/chapter-08.md',tofile='book/edition-3/chapter-08.md'))
(E/'verified-chapter.diff').write_text(diff,newline='')
links=[]
for target in re.findall(r'\]\(([^)]+)\)',ft):
 if not re.match(r'[a-z]+:|#',target):links.append({'target':target,'exists':(C.parent/target.split('#')[0]).exists()})
save('chapter-checks.json',{'baseline':ident(A/'baseline/chapter-08.md'),'final':ident(C),'saved_final_identical':True,'changes':changes,'diff_matches_author':diff.encode()==(A/'chapter-08.diff').read_bytes(),'relative_links':links,'link_targets_content_read':False,'utf8':True,'nul_bytes':f.count(b'\0'),'cr_bytes':f.count(b'\r')})
run('chapter-lint.log',['go','run','./cmd/lintprose',str(C)])
run('chapter-diff-check.log',['git','diff','--check','--',str(C)])
run('scope-check.txt',['git','--no-pager','status','--porcelain=v1','--untracked-files=all'])
print('Author inventory:',len(rows),'members verified; input pins:',len(checks),'verified; changed ranges:',changes)
