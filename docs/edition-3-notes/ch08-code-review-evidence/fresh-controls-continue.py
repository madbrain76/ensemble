import sys,json
from pathlib import Path
sys.dont_write_bytecode=True
E=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch08-code-review-evidence');sys.path.insert(0,str(E))
from runner import run,T,ENV
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');G=Path('/tmp/cadence-ch08-grader');GO='/usr/local/go/bin/go';node=json.loads((E/'runtime-tools.json').read_text())['node']
plan=[]
for name,cwd in [('root',S),('gui',S/'gui'),('consumer',S/'examples/ch04-consumer'),('workflow',S/'examples/ch06'),('stream',S/'examples/ch07')]:
 plan.extend([(name+'-race',[GO,'test','-race','-count=1','./...'],cwd,300),(name+'-vet',[GO,'vet','./...'],cwd,120)])
plan.extend([('page-workbench',[node,str(S/'gui/test/workbench.mjs')],S,30),('page-queues',[node,str(S/'gui/test/queues.mjs')],S,30),('page-inherited-integration',[node,'--experimental-vm-modules',str(S/'gui/test/correction-r1-page.mjs')],S,30)])
mapped=json.loads((R/'docs/edition-3-notes/ch08-student-evidence/grade-ch08-mapped.receipt.json').read_text());a=mapped['argv'];assert a[0]==str(G);assert a[a.index('-ch')+1]=='8';a=[str(x) for x in a]
# Keep exact pinned executable/models/options, with a new owned report path.
plan.append(('chapter8',a,T/'work',240))
lifetime=json.loads((R/'docs/edition-3-notes/ch08-student-evidence/grade-ch01-final.receipt.json').read_text());b=[str(x) for x in lifetime['argv']];assert b[0]==str(G);assert b[b.index('-ch')+1]=='1';assert a[-1]==b[-1];plan.append(('chapter1-lifetime-once',b,T/'work',600))
(E/'fresh-plan.json').write_text(json.dumps([{'name':n,'argv':av,'cwd':str(cwd),'timeout':limit} for n,av,cwd,limit in plan],indent=2)+'\n')
for n,av,cwd,limit in plan:
 run(n,av,cwd,limit)
 if n.startswith('chapter'):
  report=json.loads((E/(n+'.stdout.log')).read_text());(E/('chapter1-report.json' if n=='chapter1-lifetime-once' else 'chapter8-report.json')).write_text(json.dumps(report,indent=2)+'\n')
for n in ['chapter8','chapter1']:
 report=json.loads((E/(n+'-report.json')).read_text());assert report['passed'],report
(E/'fresh-complete.json').write_text(json.dumps({'commands':len(plan),'all_expected_exits':True,'chapter1_full_lifetime_runs':1,'same_pinned_cli':a[-1],'remaining_jobs':0},indent=2)+'\n')
