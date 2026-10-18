from pathlib import Path
import json,hashlib
E=Path(__file__).parent
rows=[json.loads(x) for x in (E/'reads.jsonl').read_text().splitlines()]
req={
'AGENTS.md':None,'.agents/skills/ensemble-coding/SKILL.md':None,'.agents/skills/edition-3-text-reviewer/SKILL.md':None,
'docs/edition-3-execution.md':None,'docs/edition-3-readiness.md':None,
'docs/agentic-codebook-workflow.md':[(1,300),(301,750),(751,1200),(1201,1650),(1651,2050),(2051,2434)],
'docs/edition-2-carryover.md':[(1,230),(231,630),(631,790),(791,949),(950,1110)],
'book/voice.md':[(1,378),(379,756)],'book/edition-3/voice.md':[(1,378),(379,756)],
'book/chapter-writing-procedure.md':[(1,311),(312,622)],'book/edition-3/chapter-08.md':None}
pins=json.loads((E/'input-pins.json').read_text())['inputs']
for p in pins:
 if p not in req: req[p]=None
for suffix in ['baseline/chapter-08.md','read-coverage.json','failures.md','chapter-08.diff','initial-chapter-08.diff','baseline-checks.json','initial-final-checks.json','final-checks.json','baseline-lint.log','initial-final-lint.log','final-lint.log','baseline-links.json','report-links.json','seal-check.json','verify.py']:
 req['docs/edition-3-notes/ch08-author-evidence/'+suffix]=None
out=[]
for p,segments in req.items():
 matching=[(i+1,x) for i,x in enumerate(rows) if x['path']==p]
 assert matching,p
 n=matching[-1][1]['total_lines']; segments=segments or [(1,n)]
 chosen=[]
 for a,z in segments:
  choices=[(i,x) for i,x in matching if x['start_line']==a and x['end_line']==z]
  assert choices,(p,a,z)
  i,x=choices[-1];chosen.append({'read_record':i,'lines':[a,z],'output_bytes':x['output_bytes'],'output_sha256':x['output_sha256']})
 covered=set()
 for a,z in segments:covered.update(range(a,z+1))
 assert len(covered)==n,p
 last=matching[-1][1];out.append({'path':p,'lines':n,'bytes':last['file_bytes'],'sha256':last['file_sha256'],'full_read':True,'selected_complete_reads':chosen})
(E/'read-coverage.json').write_text(json.dumps({'timing':'All prerequisite and required evidence reads completed before the first verdict/report. No chapter edits.','complete_files':out,'excluded_attempts':'Oversized combined prerequisite reads and guessed nonexistent receipt paths are not relied on. reads.jsonl records emitted bytes, not proof that an oversized tool output was fully displayed. The selected ranges above use complete/recovered reads.','narrow_source':'internal/gradeedition2/ch05_media.go lines 1-60 read after parent main-12 permission; relevant validation lines 18-41; next function opening 44-60 incidental, not assessed.','additional_empty_outputs_read':['docs/edition-3-notes/ch08-author-evidence/baseline-diff-check.log','docs/edition-3-notes/ch08-author-evidence/final-diff-check.log']},indent=2)+'\n')
print('Complete-read coverage:',len(out),'files; no uncovered prerequisite lines.')
