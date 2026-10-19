#!/usr/bin/env python3
"""Custody-only seal after aborted review; not a test runner or fixture repair."""
import datetime, hashlib, json, pathlib, shutil, stat, subprocess
M = pathlib.Path('/Users/bill/projects/ensemble')
N = M / 'docs/edition-3-notes'
E = N / 'ch10-retention-code-review-evidence'
C = pathlib.Path('/Users/bill/projects/ensemble-edition-3-ch10-retention')
W = pathlib.Path('/tmp/ch10-retention-independent-review-20261019')
IO = pathlib.Path('/Users/bill/projects/cadence-memories/children/e3-ch10-retention-code-review/io')
P = json.loads((N / 'ch10-retention-review-input.json').read_text())
started = datetime.datetime.now(datetime.timezone.utc).isoformat()
def digest(f):
    b = f.read_bytes()
    return {'bytes': len(b), 'sha256': hashlib.sha256(b).hexdigest(), 'mode': oct(stat.S_IMODE(f.stat().st_mode))}
def put(name, value):
    (E/name).write_text(json.dumps(value, indent=2) + '\n')
def git(*args):
    return subprocess.check_output(['git', '--no-pager', '-C', str(C), *args], text=True).strip()
checks=[]
def check(f, expected, kind):
    actual = digest(f)
    ok = all(actual[k] == expected[k] for k in ('bytes','sha256'))
    mode = expected.get('filesystem_mode', expected.get('mode'))
    if mode is not None:
        ok = ok and actual['mode'] == mode
    checks.append({'path':str(f),'kind':kind,'verified':ok, **actual})
    if not ok:
        raise RuntimeError('Preservation mismatch: ' + str(f))
check(pathlib.Path(P['task']['path']), P['task'], 'task')
for section in ('inputs','protected_tracked'):
    for name, expected in P[section].items():
        check(M/name, expected, section)
for expected in P['candidate_files']:
    check(C/expected['path'], expected, 'frozen candidate')
    check(W/'source'/expected['path'], expected, 'owned copy; no tests or mutation')
student=json.loads((N/'ch10-retention-student-evidence/inventory.json').read_text())
for expected in student['members']:
    check(N/expected['path'], expected, 'student evidence custody only')
identity={'head':git('rev-parse','HEAD'),'tree':git('rev-parse','HEAD^{tree}'),'parent':git('rev-parse','HEAD^'),'status':git('status','--porcelain')}
assert identity == {'head':P['candidate'],'tree':'b2f86ba9fd48ff9a923a8a924a9612549aa762c0','parent':P['base'],'status':''}
tracked=git('ls-tree','-r','--name-only','HEAD').splitlines()
assert len(tracked)==142
assert sorted(tracked)==sorted(v['path'] for v in P['candidate_files'])
assert sorted(str(f.relative_to(W/'source')) for f in (W/'source').rglob('*') if f.is_file())==sorted(tracked)
for job,name in [('220','setup-failure.txt'),('221','custody-script-failure.txt')]:
    data=(IO/job).read_bytes(); data.decode('utf-8'); assert b'\0' not in data
    (E/name).write_bytes(data)
put('preservation.json',{'started_utc':started,'finished_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'status':'PASS: custody only, not behavioral acceptance','candidate':identity,'candidate_members':len(tracked),'student_members':len(student['members']),'checks':checks,'owned_root':str(W),'owned_source_changed':False,'mutation':None,'restoration':'not needed; no mutation or test injection occurred','cleanup':'none; owned source/home/tmp retained; no shared cleanup','free_bytes':shutil.disk_usage(W).free})
put('read-ledger.json',{'governing_input_manifest':{'path':str(N/'ch10-retention-review-input.json'),**digest(N/'ch10-retention-review-input.json')},'input_hashes':P['inputs'],'task':P['task'],'read_scope':['mandatory instructions and workflow/carryover','current Chapter 10 draft and student report','previous Chapter 10 review as historical context only','exact 18-file candidate diff','changed persistence, retention, lifecycle/settings source and test fixtures; adjacent constructors'],'qualification':'Large outputs were clipped and several were recovered with bounded shell reads. Source inspection was interrupted by reviewer setup failure. This is NOT a certification that the entire requested source review completed. No future/FINAL/paused Chapter11 implementation was read.','student_receipts':'Frozen root race excluding lifetime, root vet, GUI race/vet and student negative/restored receipts carried by verified identity only; none rerun. Old live/numerical acceptance is not new-binary proof.'})
put('execution-ledger.json',{'result':'INCOMPLETE; zero Go commands executed','setup':{'cwd':'/Users/bill/projects/ensemble (tool default workspace; not explicitly overridden)','command_kind':'python3 stdin script: git archive candidate into owned source; create home/tmp; hash copy; attempt write copy manifest using p["candidate_tree"]','source_command':['git','-C',str(C),'archive',P['candidate']],'exit':1,'stdout_stderr':'setup-failure.txt','time_basis':'raw tool-output file modification time, not an independently captured command start','output_mtime_utc':datetime.datetime.fromtimestamp((IO/'220').stat().st_mtime,datetime.timezone.utc).isoformat(),'failure':'KeyError: candidate_tree','effects':['created owned source/home/tmp directories','extracted exact candidate','no manifest written before failure','no tests injected or executed']},'subsequent_custody_failure':{'command_kind':'python3 stdin read-only verification script, successfully compared file pins then incorrectly indexed p["tree"]','exit':1,'stdout_stderr':'custody-script-failure.txt','output_mtime_utc':datetime.datetime.fromtimestamp((IO/'221').stat().st_mtime,datetime.timezone.utc).isoformat(),'failure':'KeyError: tree','effect':'no report/evidence written before failure; no fixture repair or execution'},'inspection_errors':['Read-only cat of nonexistent ch10-retention-assignment.md and ch10-retention-student-assignment.md. Governing task requires neither path.','Read-only git diff used mistyped revisions 3f5d60cfbab17394bf260ccb0eea2de105f8425fd4 and 8919f9f874931c6f460dcebd70bc07e323a800a3; fatal bad object; exact diff later read from manifest identities.','Some git inspection calls omitted the required --no-pager switch. Final custody calls include it.'], 'environment':{'actual_setup':'inherited tool shell environment; no provider call, Go command, or secret capture; sensitive inherited variables not enumerated','required_test_environment_not_used':{'GOMAXPROCS':'2','GOFLAGS':'-mod=readonly -p=1','GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOWORK':'off','GOENV':'off'},'owned_HOME':str(W/'home'),'owned_TMPDIR':str(W/'tmp'),'qualification':'These directories were created but no test process used them. No valid test environment or successful control is claimed.'},'required_controls':[{'control':'one focused root retention selection under -race','status':'UNRUN'},{'control':'one focused GUI retention/settings selection under -race','status':'UNRUN'},{'control':'independent small valid-input probe if warranted','status':'UNRUN'},{'control':'one compile-valid meaningful negative after green baseline','status':'UNRUN; zero mutations'},{'control':'exact restoration and focused restored positive','status':'UNRUN; no mutation to restore'}],'plan_deviation':'Pre-setup commentary proposed non-race positives plus a race slice and a diagnostic negative; that was not an executed or accepted substitute for the task-required two race selections and preferred complementary negative. Stop occurred before selection was sealed or commands launched.','jobs':'No children created; all owned command jobs returned; no Go/test/background work remains.','final_seal_command':{'argv':['python3',str(E/'seal.py')],'cwd':str(M),'environment':'inherited shell; custody only; no Go/provider calls','output':'seal-output.txt; inventory itself self-excluded'}})
report='''# Chapter 10 retention amendment — incomplete independent delta review

## Verdict

**INCOMPLETE — not approved; no student-product defect established.** The reviewer hit a setup failure before launching any behavioral control. The governing finite task requires stopping and sealing on setup failure. No correction round is charged to the student (zero formal rounds used).

Frozen candidate: `b03c28fdb1586f16dfb53849f8c763576ebb39ae`  
Tree: `b2f86ba9fd48ff9a923a8a924a9612549aa762c0`  
Accepted base: `729bc919e734405307b4ac0060e65c40e33652e5`

## What happened

The reviewer verified governing/input/protection pins, read the task and delta, and extracted the candidate with `git archive` into `/tmp/ch10-retention-independent-review-20261019/source`. During setup, a Python script incorrectly assumed that the private input manifest had a `candidate_tree` key. It raised `KeyError` before writing the intended copy manifest or test plan. This is a **reviewer harness error**, not a failed assertion against the implementation.

A subsequent custody-only script also incorrectly assumed a `tree` key and failed. Both raw failure outputs are preserved. Inspection then confirmed that the manifest pins the commit and per-file hashes; the tree pin lives in the task. Final custody verification uses the actual fields and task tree pin. This verification is not a fixture repair or a behavioral rerun.

Earlier read-only inspection mistakes (two nonexistent assignment filenames and mistyped revision strings), clipped/recovered reads, and missing `--no-pager` on some Git inspection calls are disclosed in the evidence ledger. No failing output was replaced by a claimed passing test. The required exact test command/environment/timing records do not exist because no tests launched; this report does not claim otherwise.

## Scope and unrun obligations

The exact 18-file delta was inspected, with persistence/retention/settings/lifecycle ownership paths and fixture declarations. Inspection did not mature into a complete independent contract verdict. In particular, **all new behavioral evidence remains unrun**:

- Focused root retention controls under `-race`.
- Focused GUI retention/settings controls under `-race`.
- Any small independent valid-input probe.
- A green baseline followed by one meaningful compile-valid negative.
- Exact mutation restoration and the focused restored positive.

No Go test, vet, benchmark, original Chapter 10 lifetime/grade/whole suite, browser, live provider, external HTTP, or child agent ran. No source mutation or test injection occurred, so there is no negative/restored result to claim. Preliminary commentary about a different selection was not a compliant sealed run plan and was never executed.

## Carried evidence and boundaries

The 34 portable student-evidence members and their input-pinned inventory/source manifest/parent verification were hash/mode checked. Their root race (excluding `TestRequestLifetime`), root vet, GUI race/vet, and student diagnostic-bypass negative/restored receipts remain **student evidence**, not independently rerun acceptance. Old Chapter 10 numerical/live acceptance is historical and not new-binary proof.

The student report's earlier unintended lifetime execution, incorrectly directed GUI invocation, missing wrapper completion footer, and setup/compile mistakes remain disclosed in their unchanged evidence. They were not used as new product findings or excuses to rerun suites.

Final custody checks establish all 142 candidate members and all 142 owned-copy bytes/modes match the pins; candidate HEAD/tree/parent and clean status match the task. Pinned governing/current/protected inputs and the student seal remain unchanged. The shared main worktree had unrelated pre-existing changes; no blanket clean-main claim is made. No frozen source, guide, chapter, or paused Chapter 11 content was edited. No commit/tag/acceptance marker was created.

## Seal and handoff

Evidence: `ch10-retention-code-review-evidence/inventory.json` (self-excluded, report included), `preservation.json`, `read-ledger.json`, `execution-ledger.json`, and the two raw failure outputs. The small sealing script is included for custody reproducibility; it is not a test harness. UTF-8/no-NUL and closed-world membership were checked when sealing.

No owned jobs remain active. No cleanup was performed; the exact owned source and unused private HOME/TMP directories remain for custody. A parent-authorized follow-up review is needed to complete the unrun obligations; this sealed incomplete result must not be treated as approval or a request for student production changes. Root retains acceptance, prose review, correction-round ownership, and later Chapter 11 integration.
'''
(N/'ch10-retention-code-review.md').write_text(report)
(E/'seal-output.txt').write_text('Custody verification PASS; behavioral review INCOMPLETE.\n142 frozen and 142 owned source members verified; 34 student evidence members verified.\nCandidate clean; protected inputs preserved; zero tests, zero mutants, zero children.\nExit status: 0 after inventory validation (see actual tool receipt).\n')
members=[]
for f in [N/'ch10-retention-code-review.md', *sorted(E.iterdir())]:
    if f.name=='inventory.json': continue
    b=f.read_bytes(); b.decode('utf-8'); assert b'\0' not in b
    members.append({'path':str(f.relative_to(N)),**digest(f)})
put('inventory.json',{'schema':1,'created_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'result':'INCOMPLETE','self_excluded':True,'encoding':'UTF-8 without NUL','members':members})
actual={str(f.relative_to(N)) for f in E.iterdir() if f.name!='inventory.json'} | {'ch10-retention-code-review.md'}
assert actual=={v['path'] for v in members}
for v in members:
    assert digest(N/v['path'])=={k:v[k] for k in ('bytes','sha256','mode')}
print('SEALED INCOMPLETE:',len(members),'members; all custody checks passed; no behavioral checks executed.')
print('inventory sha256:',digest(E/'inventory.json')['sha256'])
