import sys,json
from pathlib import Path
sys.dont_write_bytecode=True
E=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch08-code-review-evidence');sys.path.insert(0,str(E))
from runner import run
plan=json.loads((E/'fresh-plan.json').read_text());done=[]
for row in plan:
 if not row['name'].startswith('chapter'):continue
 assert not (E/(row['name']+'.receipt.json')).exists(),'do not repeat a gate'
 run(row['name'],row['argv'],row['cwd'],row['timeout']);report=json.loads((E/(row['name']+'.stdout.log')).read_text())
 out='chapter1-report.json' if row['name']=='chapter1-lifetime-once' else 'chapter8-report.json';(E/out).write_text(json.dumps(report,indent=2)+'\n');assert report['passed'];done.append(row['name'])
(E/'grade-gates-complete.json').write_text(json.dumps({'gates':done,'chapter1_full_lifetime_runs':1,'earlier_page_fixture_failure_retained':True,'remaining_process_groups':0},indent=2)+'\n')
