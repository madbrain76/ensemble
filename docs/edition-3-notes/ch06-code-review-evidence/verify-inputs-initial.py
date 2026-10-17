from pathlib import Path
import hashlib,json,subprocess,tarfile
E=Path(__file__).resolve().parent;ROOT=E.parents[2];S=Path((E/'scratch-path.txt').read_text().strip());ST=Path('/Users/bill/projects/ensemble-edition-3-student')
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
assert sha('/tmp/cadence-ch06-code-review-task.txt')=='0f7c1f4d0b802949e18cdb6b327e8233187de641f90e25116f80da804950b926'
assert sha(ROOT/'docs/edition-3-notes/ch06-candidate.json')=='55f613021fc62b332ec935d6e59e01fd95d6a34bfca3272d5d4539de815432c3'
c=json.loads((ROOT/'docs/edition-3-notes/ch06-candidate.json').read_text()); records=[]
for d in c['artifacts']+c['grader_inputs']:
 p=Path(d['path']);assert p.stat().st_size==d['bytes'],p;assert sha(p)==d['sha256'],p;records.append(d)
m=json.loads((ROOT/'docs/edition-3-notes/ch06-student-evidence-manifest.json').read_text());assert len(m['artifacts'])==474
for d in m['artifacts']:
 p=Path(d['path']);assert p.stat().st_size==d['size_bytes'],p;assert sha(p)==d['sha256'],p
assert subprocess.check_output(['git','-C',str(ST),'rev-parse','HEAD'],text=True).strip()==c['student_commit']
assert subprocess.check_output(['git','-C',str(ST),'status','--porcelain'],text=True)==''
sm=json.loads((ROOT/'docs/edition-3-notes/ch06-student-source-manifest.json').read_text())
source=S/'source';source.mkdir()
archive=ROOT/'docs/edition-3-notes/ch06-student-evidence/candidate-source.tar'
with tarfile.open(archive) as t:t.extractall(source,filter='data')
files=subprocess.check_output(['git','-C',str(ST),'ls-tree','-r','--name-only',c['student_commit']],text=True).splitlines();assert len(files)==98
source_hashes={p:sha(source/p) for p in files}
for p,h in source_hashes.items():assert sha(ST/p)==h,p
assert len([p for p in source.rglob('*') if p.is_file()])==98
prereqs=['AGENTS.md','.agents/skills/ensemble-coding/SKILL.md','.agents/skills/edition-3-code-reviewer/SKILL.md','docs/edition-3-execution.md','docs/edition-3-readiness.md','docs/edition-3-carryover.md','docs/edition-3-student.md','docs/agentic-codebook-workflow.md','docs/edition-2-carryover.md','book/voice.md','book/chapter-writing-procedure.md','book/edition-3/chapter-04.md','book/edition-3/chapter-05.md','book/edition-3/chapter-06.md','docs/edition-3-notes/ch06-entry-audit.md','docs/edition-3-notes/ch06-grader-review.md','docs/edition-3-notes/ch06-student-review.md','docs/edition-3-notes/ch06-student-evidence/COMMANDS.md','docs/edition-3-notes/ch06-student-source-manifest.json','docs/edition-3-notes/ch06-student-evidence-manifest.json']
(E/'source-hashes.json').write_text(json.dumps(source_hashes,indent=2)+'\n')
(E/'prerequisite-hashes.json').write_text(json.dumps({p:sha(ROOT/p) for p in prereqs},indent=2)+'\n')
(E/'input-verification.json').write_text(json.dumps({'candidate':c,'verified_artifact_count':474,'verified_source_files':98,'artifact_inventory_metadata':{k:v for k,v in m.items() if k!='artifacts'},'source_scratch':str(source),'raw_private_contents_disclosed':False},indent=2)+'\n')
print('Verified all 474 artifacts; extracted 98 exact source files; candidate and grader pins match. No private payload contents printed.')
print('SOURCE',source)
