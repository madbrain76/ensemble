#!/usr/bin/env python3
"""Seal only completed receipts; never launches an exercise or reads credentials."""
import pathlib,json,hashlib,os,subprocess,shutil,time,stat
ROOT=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');D=ROOT/'.ch09-local/operational-execution';P=ROOT/'.ch09-local/operational-preparation';RUN=P/'execution-runs/ch09-root-execution-01';NOTES=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes');OUT=NOTES/'ch09-operational-exercise-evidence';REPORT=NOTES/'ch09-operational-exercise.md'
def sha(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def load(p):return json.loads(pathlib.Path(p).read_text())
def save(p,x):pathlib.Path(p).write_text(json.dumps(x,indent=2)+'\n')
assert not OUT.exists() and not REPORT.exists();assert shutil.disk_usage(ROOT).free>=1024**3
policy=load(P/'policy.json');seal=load(P/'seal.json');manifest=load(policy['source_manifest']['path']);assert sha(P/'seal.json')=='070e03c3703a629fe3b29cd1adadc9292563a6bcf07ce1839aaefe910fb43ab1';assert sha(P/'policy.json')=='b8d8525398503564367003eaa3c34dc73dab68be971238c86377752e213a0e79'
for row in seal['preparation_files']:assert sha(P/row['path'])==row['sha256']
for row in policy['source_files']:assert sha(ROOT/row['path'])==row['sha256']
for row in manifest['binaries']:assert sha(row['path'])==row['sha256'] and pathlib.Path(row['path']).stat().st_size==row['bytes']
for name,digest in policy['assets'].items():assert sha(ROOT/'gui/web/gui'/name)==digest
assert sha(policy['source_manifest']['path'])==policy['source_manifest']['sha256']
assert subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=ROOT,text=True).strip()==policy['source_commit'];assert subprocess.check_output(['git','--no-pager','rev-parse','HEAD^{tree}'],cwd=ROOT,text=True).strip()==policy['source_tree'];assert not subprocess.check_output(['git','--no-pager','status','--porcelain'],cwd=ROOT,text=True).strip()
prior=load(P/'preservation.json');previous=[]
for row in prior['prior_seals']+prior['public_read_carry']:
 assert sha(row['path'])==row['sha256'];previous.append(dict(path=row['path'],sha256=row['sha256']))
for row in load(D/'preflight.json')['prior_inputs']:assert sha(row['path'])==row['sha256']
# Preserve the prior portable inventory exactly, including superseded report snapshots.
prior_dir=NOTES/'ch09-operational-preparation-evidence';prior_index=load(prior_dir/'index.json')
for row in prior_index['files']:
 f=prior_dir/row['path'];assert f.stat().st_size==row['bytes'] and sha(f)==row['sha256']
previous.extend([dict(path=str(prior_dir/'index.json'),sha256=sha(prior_dir/'index.json')),dict(path=str(NOTES/'ch09-operational-preparation.md'),sha256=sha(NOTES/'ch09-operational-preparation.md'))])
original=load(D/'browser-local.receipt.json');assert original['exit_code']==0 and original['gate_passed'] is False and original['gate_error']=='browser exit receipt missing';assert sha(D/'run-once.py')==original['wrapper_sha256'];assert sha(D/'audit-entry.py')==original['audit_wrapper_sha256']
paid_receipt=load(D/'paid-cli.receipt.json');assert paid_receipt['exit_code']==0 and paid_receipt['gate_passed'];assert sha(D/'run-paid-once.py')==paid_receipt['wrapper_sha256']
checked=[];summaries={};pids=set()
for mode,count,cap in (('browser-local',6,4096),('paid-cli',5,1024)):
 p=RUN/mode;result=load(p/'result.json');req=load(p/'proxy-requests.json');app=load(p/'artifact-close.json');proxy=load(p/'proxy-close.json');audit=load(D/(mode+'.audit.json'))
 assert result['ok'] and all(x['closed'] for x in result['cleanup']);assert len(req)==count and all(x['complete'] and x['forward_attempted'] and x['upstream_status']==200 and x['model']=='claude-sonnet-4-6' and x['path']=='/v1/messages' and x['method']=='POST' and x['max_output_field']=='max_tokens' and x['max_output_value']==cap for x in req)
 assert app['process_exited'] and app['listener_closed'] and app['exit_code']==0;assert proxy['listener_closed'] and proxy['active_upstreams']==0 and proxy['requests']==count and proxy['retries']==0
 assert audit['credential_file_open_count']==(0 if mode=='browser-local' else 1)
 assert (p/'workspace/note.txt').read_text()=='SKILLS_OK';assert result['effects']['request_count']==count
 for f in p.glob('*ownership.json'):
  v=load(f)
  if isinstance(v.get('pid'),int) and v['pid']>0:pids.add(v['pid'])
 checked.extend([dict(path=str(p/n),sha256=sha(p/n)) for n in ('artifact-close.json','proxy-close.json','result.json')])
 summaries[mode]=dict(result_ok=True,controller_elapsed_seconds=result['elapsed_seconds'],setup_seconds=result['setup_seconds'],requests=count,http_statuses=[x['upstream_status'] for x in req],max_tokens=cap,credential_file_open_count=audit['credential_file_open_count'],call_names=[c['name'] for c in result['effects']['calls']],retry_count=0,cleanup=result['cleanup'])
br=load(RUN/'browser-local/owned-browser-exit.json');assert br['exit_code']==0 and br['listener_closed'];pids.add(br['pid']);assert (RUN/'browser-local/workspace/after-unload.txt').read_text()=='SKILLS_OK'
for pid in pids:
 try:os.kill(pid,0)
 except ProcessLookupError:pass
 else:raise AssertionError('owned PID remains: '+str(pid))
assert set(x.name for x in RUN.iterdir() if x.is_dir())=={'browser-local','paid-cli'}
att=load(D/'root-protection-attestation.json');assert sha(D/'root-protection-attestation.json')=='632ec51d6c1e8de4bc7c2d8de22c62fc23c37394837780bd89a322fd036f74b6'
usage=load(D/'paid-usage.json');assert usage['numeric_totals_match']
preservation=dict(checked_unix=time.time(),source_commit=policy['source_commit'],source_tree=policy['source_tree'],tracked_source_clean=True,source_file_count=len(policy['source_files']),five_binary_pins=manifest['binaries'],served_assets=policy['assets'],sealed_preparation=dict(seal_sha256=sha(P/'seal.json'),policy_sha256=sha(P/'policy.json'),preparation_files=seal['preparation_files']),prior_seals_and_public_read_pins=previous,prior_portable_inventory_verified=len(prior_index['files']),root_protection_attestation=att,root_protection_note='Root-only time-qualified comparison; student did not inspect private baseline.',owned_pids_absent=sorted(pids),cleanup_receipts=checked,disk_free_bytes=shutil.disk_usage(ROOT).free)
save(D/'final-preservation.json',preservation)
qualification=dict(transport_error707='Parent reported inference/provider transport interruption after browser completion; recovered existing job2843, not another exercise. No browser replay.',browser_actual='Sealed controller exit0/result.ok:true; one completed browser-local attempt.',outer_wrapper='Original run-once.py outer exit1 and gate_passed:false preserved. It looked under browser-local/browser/*browser-exit.json, but actual receipt is browser-local/owned-browser-exit.json. No global green relabel.',reconciliation='Root main-18 authorized receipt-only correction after independently verifying actual browser cleanup; new paid-only wrapper changed only scope/receipt lookup/prerequisite reconciliation. Diff preserved.',original_root_setup_mismatch='Original task command syntax, approval fields and run-root did not match sealed controller. Stopped before artifact/key/upstream; root supplied task/approval v2 and confirmed zero consumed attempts.',local_preflight_lookup='Read-only wrapper initially used root/app.js instead of gui/web/gui/app.js; failure preserved, corrected before launch.',preparation_scope='Root removed unapproved original three-provider/paid-GUI matrix before any such execution. Superseded preparation remains untouched.',fixture_qualification='Earlier sealed preparation report says invalid field dependencies; root verified actual saved initial stderr says depends-on. Corrected only here, not in old report.',marker_qualification='CH09_CALLBACK_VALUE is the CLI/GUI custom-variable marker. File content is SKILLS_OK. This exercises interpolation; marker name is not additional proof of RegisterVar API callbacks or a new product requirement.',uncredited_log_byte_extraction='Optional post-execution decoded API-log byte comparison did not establish wire-byte equality. Negative extraction result retained, not a provider failure. Numeric usage was extracted from all five actual responses and reconciled with CLI totals. No additional call/diagnosis/retry was made for this.',visual_qualification='Two settled PNGs captured locally; neither personally visually inspected by student. DOM/typed facts passed. No speech/MIC/restart/extra page/paid GUI.',host_qualification='Audit observes Python controller opens/connections only. It does not establish absence of all Chrome/host background I/O; no universal host-no-write/no-network claim.')
save(D/'qualifications.json',qualification);save(D/'exercise-summary.json',dict(scope='one browser-local, then one paid-cli; no reruns',exercises=summaries,usage=usage['cli_usage'],cost_usd=None,owned_pids_absent=sorted(pids)))
# Local artifact inventory contains hashes, never copies of binaries/profiles/PNG/API/opaque traces.
local=[]
for base,label in ((D,'execution-metadata'),(RUN,'exercise-runs')):
 for f in sorted(base.rglob('*')):
  if f.name=='local-artifact-index.json':continue
  info=f.lstat();row=dict(root=label,path=str(f.relative_to(base)))
  if stat.S_ISREG(info.st_mode):row.update(kind='file',bytes=info.st_size,sha256=sha(f))
  elif stat.S_ISLNK(info.st_mode):row.update(kind='symlink',target_sha256=hashlib.sha256(os.readlink(f).encode()).hexdigest())
  elif stat.S_ISDIR(info.st_mode):continue
  else:row.update(kind='nonregular',mode=stat.S_IFMT(info.st_mode))
  local.append(row)
save(D/'local-artifact-index.json',dict(self_excluded='local-artifact-index.json',roots=dict(execution_metadata=str(D),exercise_runs=str(RUN)),files=local,frozen_binaries_hash_only=manifest['binaries'],local_only_roles=['binaries','source archives','browser profile','PNG screenshots','raw API/event/browser/CDP traces'],no_private_contents_exported=True))
OUT.mkdir();copied=[]
def copy(source,relative):
 target=OUT/relative;target.parent.mkdir(parents=True,exist_ok=True);raw=pathlib.Path(source).read_bytes();raw.decode('utf-8');assert b'\0' not in raw;target.write_bytes(raw);copied.append(relative)
# Exact new wrapper/receipt metadata. No old bundle contents are rewritten.
for f in sorted(D.iterdir()):
 if f.is_file() and f.suffix in ('.py','.json','.txt','.diff'):copy(f,'execution/'+f.name)
authority=[('task-original.txt','/tmp/cadence-ch09-operational-execution-task.txt'),('approval-original.json','/tmp/cadence-ch09-operational-execution-approval.json'),('task-v2.txt','/tmp/cadence-ch09-operational-execution-task-v2.txt'),('approval-v2.json','/tmp/cadence-ch09-operational-execution-approval-v2.json'),('root-protection-attestation.json','/tmp/cadence-ch09-operational-root-protection-attestation.json')]
for name,path in authority:copy(path,'authority/'+name)
common=('result.json','sanitized-effects.json','proxy-requests.json','proxy-close.json','proxy-ownership.json','artifact-close.json','artifact-ownership.json')
for mode in ('browser-local','paid-cli'):
 for name in common:copy(RUN/mode/name,mode+'/'+name)
 if mode=='browser-local':
  for name in ('browser-facts.json','typed-facts.json','served-assets.json','gui-listener.json','owned-browser-attached.json','owned-browser-exit.json','owned-browser-ownership.json','upstream-requests.json','upstream-ownership.json'):copy(RUN/mode/name,mode+'/'+name)
for f in sorted(RUN.glob('*-receipt.json')):copy(f,'controller/'+f.name)
bf=load(RUN/'browser-local/browser-facts.json');tf=load(RUN/'browser-local/typed-facts.json');assert tf['ok'] and not tf['errors'] and tf['complete'];assert not bf['after']['paused'] and not bf['after']['read_aloud'] and not bf['after']['body_overflow'];assert len(bf['after']['cards'])==5
REPORT.write_text(f'''# Chapter 9 — bounded operational exercise

**Both actual exercises passed.** Exactly one browser-local attempt then one paid Claude CLI attempt ran. No rerun/retry, new browser, paid GUI, provider matrix, production change, rebuild/test/grader, commit/tag or child agent. No owned process/job remains. **An outer browser receipt-wrapper failure is preserved, not relabelled green.**

## Actual results and limits

| Exercise | Actual result | Requests/status | Native cap | Sealed elapsed / wrapper |
|---|---|---|---|---|
| Muted owned Chrome + frozen GUI, local fake upstream | Controller exit 0, `ok:true` | 6 / six HTTP 200 | 4096 | 6.053s / 6.192s; outer wrapper exit **1** (lookup defect below) |
| Frozen CLI, real `claude-sonnet-4-6` | Controller and wrapper exit 0, `ok:true` | 5 / five HTTP 200 | **1024**, unchanged | 15.183s / 15.305s |

Browser authority was one normal muted profile, ≤600 seconds/8 fake requests. Paid authority was one session, ≤8 total requests, ≤60 seconds setup/180 active/240 absolute; measured setup was 0.054s. No automatic next attempt, redirect, fallback, token/body change or retry occurred. HTTP/request counters are actual proxy reservations, including all turns. Only Anthropic POST `/v1/messages` with the pinned model was forwarded in the paid case.

The graph `base → code-tools → helper` activated through real calls. Both branches loaded `code-tools`, wrote **`note.txt = SKILLS_OK`**, read that exact value, marked unload pending and answered verified content. Primary system hash remained identical across requests; next declarations retained activated tools. Browser-local additionally wrote `after-unload.txt = SKILLS_OK` through the still-enabled dynamic tool. Paid made no extra post-unload tool call. Durable call/result IDs, file hashes, pending result, declarations and system hashes are in each exercise’s `sanitized-effects.json` / `proxy-requests.json`.

## GUI, privacy and usage

The real page, production WebSocket and frozen assets were used. Typed messages correlate five dispatch/result pairs, completed turn and idle/unpaused state with five DOM action cards. Helper instructions and `CH09_CALLBACK_VALUE` interpolation are visible; transcript answers the verified file value. DOM facts show three pane bounds, no body overflow, connected/idle, read-aloud false. Informational completed-job notices were not treated as errors. Two settled screenshots were captured, but **not personally visually inspected**; they remain local/hash-only at `{RUN}/browser-local/initial.png` and `skills-result.png`. No speech/MIC/audio/focus, restart, second-page or paid GUI credit is claimed.

The read-only Python audit counted **0 credential-file opens in browser-local and 1 in paid-cli**. The latter is the sealed proxy’s approved `~/.cr/settings.json` access after pin/reservation; values were not logged. Python-observed browser destinations were loopback only; paid had five TLS connections to the resolved Anthropic endpoint. This does **not** assert all Chrome/host background networking or writes were absent. Child HOME/environment were isolated, with dummy keys, and Chrome was muted. Raw API/auth/event/CDP traces and the profile are local/hash-only, not portable evidence.

Actual CLI usage was **8,358 input / 729 output tokens**, cache write/read both 0. Five retained provider response usage objects sum to those exact values; per-response outputs were 127, 195, 76, 131 and 200. No price lookup or dollar estimate was made. An optional post-hoc decoder did **not** establish byte equality between reconstructed API-log JSON and raw proxy fingerprints; that negative extraction is preserved and uncredited, not an HTTP/provider failure. Local upstream body-hash equality is established; paid body preservation is supported by the unchanged sealed forwarder and actual request fingerprints, not a claimed independent upstream byte receipt.

## Failures, corrections and transport recovery

1. Root’s original command/approval/run-root did not fit the frozen controller interface. Prelaunch stopped before artifact/key/upstream. Root supplied exact v2 task/approval and confirmed no attempt consumed. Original inputs remain in evidence. A separate read-only asset-path lookup failure was corrected before launch, with its failure retained.
2. The actual browser exercise passed, but new outer `run-once.py` searched `browser-local/browser/*browser-exit.json`; the real receipt is `browser-local/owned-browser-exit.json`. Original wrapper, outer exit 1 and `gate_passed:false` remain unchanged. Root independently confirmed successful actual cleanup and explicitly authorized receipt-only reconciliation (main-18). New paid-only wrapper/diff checks the direct receipt and actual prerequisite without changing sealed functions/assertions/timeouts. There is no global all-green relabel.
3. **Error 707 was an inference/transport interruption after browser completion**, not a browser or paid result. Existing job2843 was collected; durable results were reconciled; no browser replay occurred. The unused paid attempt then ran once under still-valid v2 authority.
4. The earlier unapproved three-provider/paid-GUI preparation was removed after root’s scope correction before execution. Prior bundles remain untouched. Earlier sealed report’s `dependencies` wording is inaccurate: root verified the original failing fixture field was **`depends-on`**. This qualification does not rewrite that report.
5. `CH09_CALLBACK_VALUE` is the accepted interpolation marker; **`SKILLS_OK` is file content**. Its name is not extra proof of public `RegisterVar` callback execution and is not a new product requirement.

## Cleanup, identity and exact evidence

Chrome exit0/listener closed; GUI and paid CLI exit0/process exited; both proxies closed with zero active upstreams and zero retries. All recorded owned PIDs are absent. Source remains `{policy['source_commit']}`, tree `{policy['source_tree']}`; 132 tracked source files, all five binaries, served assets, sealed scripts/policy/fixtures, prior public-read hashes and prior bundles verified unchanged. Root’s separately attributed attestation reports **33/33 protected tracked inputs unchanged, 0 drift at {att['checked_utc']}**; student did not read the private baseline. This is time-qualified, not a universal preservation claim.

The exact commands, cwd, sanitized environment, preflight pins, timing, exit status, read-only audit, original failures, v2 authority and root reconciliation are retained. Preparation seal remains `070e03c3703a629fe3b29cd1adadc9292563a6bcf07ce1839aaefe910fb43ab1`; policy remains `b8d8525398503564367003eaa3c34dc73dab68be971238c86377752e213a0e79`. Wrapper/collector scripts are new receipt-only artifacts, not production patches. No private audit/grader input was read; assistance was limited to the exact public tasks, approved scaffolding/preparation and root’s explicit operational corrections/attestation.

`ch09-operational-exercise-evidence/index.json` is the exact self-excluded UTF-8/no-NUL portable inventory and includes this report. `execution/local-artifact-index.json` hashes the owned raw traces/profile/PNGs without copying them. `execution/final-preservation.json` binds cleanup and unchanged identities; `execution/qualifications.json` preserves each limitation. No genuine blocker remains; root’s final inspection/acceptance is still its own decision.
''')
rows=[]
for f in sorted(OUT.rglob('*')):
 if f.is_file():
  raw=f.read_bytes();raw.decode('utf-8');assert b'\0' not in raw;rows.append(dict(path=str(f.relative_to(OUT)),bytes=len(raw),sha256=sha(f)))
raw=REPORT.read_bytes();raw.decode();assert b'\0' not in raw;rows.append(dict(path='../'+REPORT.name,bytes=len(raw),sha256=sha(REPORT)))
index=dict(version=1,self_excluded='index.json',external_report='../'+REPORT.name,source_commit=policy['source_commit'],preparation_seal_sha256=sha(P/'seal.json'),policy_sha256=sha(P/'policy.json'),all_utf8_no_nul=True,raw_private_traces_profiles_png_binaries_copied=False,files=rows)
save(OUT/'index.json',index)
print(json.dumps(dict(report=str(REPORT),report_sha256=sha(REPORT),evidence=str(OUT),index_sha256=sha(OUT/'index.json'),inventory_files=len(rows),owned_pids_absent=sorted(pids),actual_exercises=summaries,remaining_live_jobs=0),indent=2))
