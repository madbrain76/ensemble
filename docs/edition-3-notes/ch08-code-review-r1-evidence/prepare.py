import json,hashlib,subprocess,tempfile,shutil,tarfile,os,datetime
from pathlib import Path
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';E=N/'ch08-code-review-r1-evidence';E.mkdir(exist_ok=True)
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
M=json.loads((N/'ch08-code-review-r1-input.json').read_text());old=json.loads((N/'ch08-code-review-input.json').read_text());SM=json.loads((N/'ch08-student-r1-source-manifest.json').read_text());PM=json.loads((N/'ch08-student-r1-evidence-manifest.json').read_text())
def dump(n,x):(E/n).write_text(json.dumps(x,indent=2)+'\n')
def verify(x,base=None):
 p=Path(x['path']);p=p if p.is_absolute() else base/p;assert p.is_file(),p;assert p.stat().st_size==x['bytes'] and sha(p)==x['sha256'],p
for x in M['inputs']:verify(x)
for x in M['protected_tracked'].values():verify(x)
for x in M['source_members']:verify(x,S)
assert SM['tracked']==M['source_members'];assert SM['candidate']==M['student']
for x in PM['entries']:verify(x,N);b=(N/x['path']).read_bytes();assert b'\0' not in b;b.decode('utf-8')
for x in PM['local_only_exclusions']:verify(x)
for x in SM['unchanged_initial_artifacts']:verify(x)
G=lambda *a:subprocess.check_output(['git','--no-pager',*a],cwd=S,text=True).strip()
assert G('rev-parse','HEAD')==M['student']['commit'];assert G('rev-parse','HEAD^{tree}')==M['student']['tree'];assert G('rev-parse','HEAD^')==M['student']['parent'];assert G('status','--porcelain=v1','--untracked-files=all')==''
changed=G('diff','--name-only',M['student']['parent'],M['student']['commit']).splitlines();assert changed==['gui/test/correction-r1-page.mjs']
oldby={x['path']:x for x in old['source_members']};newby={x['path']:x for x in M['source_members']};assert oldby.keys()==newby.keys();unchanged=[p for p in oldby if oldby[p]['sha256']==newby[p]['sha256']];assert len(unchanged)==125
assert len(oldby)==126 and sorted(set(oldby)-set(unchanged))==changed
old_fixture=subprocess.check_output(['git','--no-pager','show',M['student']['parent']+':'+changed[0]],cwd=S);assert hashlib.sha256(old_fixture).hexdigest()==oldby[changed[0]]['sha256']
(E/'fixture-before.mjs').write_bytes(old_fixture);(E/'fixture-after.mjs').write_bytes((S/changed[0]).read_bytes());(E/'candidate.diff').write_bytes(subprocess.check_output(['git','--no-pager','diff',M['student']['parent'],M['student']['commit'],'--',changed[0]],cwd=S))
I=json.loads((N/'ch08-code-review-evidence/inventory.json').read_text());[verify(x,R) for x in I['files']]
prior_reads=json.loads((N/'ch08-code-review-evidence/prerequisite-carry-forward.json').read_text());read_records=[]
for x in prior_reads:
 p=R/x['path'];assert sha(p)==x['current_sha256'],p;read_records.append({'path':x['path'],'sha256':sha(p),'qualification':'reloaded fully in this turn' if x['path'] in ['AGENTS.md','.agents/skills/ensemble-coding/SKILL.md','.agents/skills/edition-3-code-reviewer/SKILL.md','docs/edition-3-execution.md'] else 'unchanged from complete initial review reading; prior section-only reads remain section-only'})
archive=SM['archive'];verify(archive);afile=Path(archive['path']);members=[]
with tarfile.open(afile) as a:
 for m in a.getmembers():
  assert not m.name.startswith('/') and '..' not in Path(m.name).parts,m.name
  if m.isfile():
   data=a.extractfile(m).read();assert m.name in newby;assert hashlib.sha256(data).hexdigest()==newby[m.name]['sha256'];members.append(m.name)
  else:assert m.isdir(),m.name
assert set(members)==set(newby)
wire_record=json.loads((N/'ch08-student-r1-evidence/wire-provenance.json').read_text());wire=Path('/Users/bill/projects/ensemble-edition-3-student/.ch08-local-r1/artifacts/r1-wire.json');assert sha(wire)=='5be64a1bda890831bd5dbda9e5d5d6354521c6ef9e991ee3eb14f5b7e0ba97ae';wiredata=json.loads(wire.read_text());assert len(wiredata)==2
for messages in wiredata:
 finals=[m for m in messages if m.get('type')=='part_final' and m.get('part',{}).get('type')=='tool_result'];replays=[m for m in messages if m.get('type')=='part_snapshot' and m.get('part',{}).get('type')=='tool_result'];assert finals and replays;assert finals[-1]['part_id']==replays[-1]['part_id'];assert finals[-1]['role']=='tool' and replays[-1]['role']=='tool';assert finals[-1]['part']==replays[-1]['part']
T=Path(tempfile.mkdtemp(prefix='cadence-ch08-code-r1-review-',dir='/tmp'));(E/'scratch-path.txt').write_text(str(T)+'\n');(T/'main-status-before.z').write_bytes(subprocess.check_output(['git','--no-pager','status','--porcelain=v1','-z','--untracked-files=all'],cwd=R));(T/'student-tags-before.txt').write_text(G('show-ref','--tags')+'\n')
for d in ['home','tmp','work']: (T/d).mkdir()
node=Path(shutil.which('node')).resolve();env={'PATH':str(node.parent)+':/usr/local/bin:/usr/bin:/bin','HOME':str(T/'home'),'TMPDIR':str(T/'tmp'),'LANG':'C.UTF-8','LC_ALL':'C.UTF-8','PAGER':'cat','GIT_PAGER':'cat','NODE_NO_WARNINGS':'0'}
dump('environment.json',{'env':env,'node_path':str(node),'node_sha256':sha(node),'node_version':subprocess.check_output([str(node),'--version'],env=env,text=True).strip(),'credentials_inherited':False,'network_or_browser_invoked':False})
dump('verification-before.json',{'checked_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'input_manifest_sha256':sha(N/'ch08-code-review-r1-input.json'),'student':M['student'],'source_count':126,'unchanged_source_count':125,'changed_paths':changed,'input_count':len(M['inputs']),'protected_count':len(M['protected_tracked']),'student_portable_entries':len(PM['entries']),'prior_review_verified_members':len(I['files']),'unchanged_prior_binaries':len(SM['unchanged_initial_artifacts']),'archive_source_members':len(members),'archive_local_only_preserved':True,'wire_sha256':sha(wire),'wire_streams':len(wiredata),'wire_tool_result_live_replay_identity_verified':True,'free_bytes':shutil.disk_usage(R).free})
dump('prerequisite-qualification.json',read_records);dump('preservation-baseline.json',{'inputs':M['inputs'],'source_members':M['source_members'],'protected_tracked':M['protected_tracked'],'student_portable':PM,'source_manifest':SM})
dump('attempt-notes.json',{'read_only_errors':[{'kind':'pager','detail':'First diff read omitted --no-pager. Owned shell/group 61610, git 61612, less 61613 identified; only verified pager 61613 terminated. Read job 1544 then ended exit 0. Full diff recovered using --no-pager. Not a runtime control.'},{'kind':'absent guessed adjacent filenames','detail':'speech-queue.js and tool-readable.js guesses absent; actual speech.js and tool-status.js read instead.'},{'kind':'manifest schema','detail':'Read summary assumed files instead of R1 entries; actual schemas inspected before verification. No control started.'}],'controls_rerun':False})
print(json.dumps(json.loads((E/'verification-before.json').read_text()),indent=2));print('SCRATCH',T)
