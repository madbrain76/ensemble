from pathlib import Path
import datetime,hashlib,json,os,shutil,stat,subprocess
root=Path('/Users/bill/projects/ensemble-edition-3-student');base=root/'.ch08-local/operational-runs';meta=base/'controller-local-01';run=base/'local-01';notes=Path('/Users/bill/projects/ensemble/docs/edition-3-notes');out=notes/'ch08-operational-exercise-evidence';report=notes/'ch08-operational-exercise.md'
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest();load=lambda p:json.loads(p.read_text())
def write(p,v):p.write_text(json.dumps(v,indent=2)+'\n')
def info(p):return {'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)}
assert out.exists() and not (out/'MANIFEST.json').exists() and not report.exists();task=Path('/tmp/cadence-ch08-operational-execution-task.txt');assert task.stat().st_size==7696 and sha(task)=='0709ca122d7dc35a5ac4ad6817288d3c0759ace23d25a08bc968cd36f81eca77'
src=load(notes/'ch08-student-r1-source-manifest.json');assert subprocess.check_output(['git','--no-pager','status','--porcelain=v1'],cwd=root,text=True)==''
assert subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=root,text=True).strip()==src['candidate']['commit']=='b8f22ba5b3c9d1fc52d1a114810b90b77928684d'
for r in src['tracked']:assert sha(root/r['path'])==r['sha256']
for r in src['unchanged_initial_artifacts']:assert sha(Path(r['path']))==r['sha256']
old=[]
for f in [notes/'ch08-student-evidence-manifest.json',notes/'ch08-student-r1-evidence-manifest.json',notes/'ch08-operational-plan-evidence/MANIFEST.json']:
 m=load(f)
 for r in m['entries']:assert sha(notes/r['path'])==r['sha256']
 old.append(info(f))
public=load(root/'.ch08-local/public-input-identity.json')
for r in public['files']:assert sha(Path(r['path']))==r['sha256']
policy=load(root/'.ch08-local/operational-prep/policy.json');assert sha(root/'.ch08-local/operational-prep/policy.json')=='eef19bf051df623c71d6227f1e375b48098177009e1d801f50f713e3533e6b27';assert sha(root/'.ch08-local/operational-prep/launcher.py')==policy['launcher_sha256']
result=load(meta/'result.json');receipt=load(root/'.ch08-local/operational-local-native.receipt.json');done=load(run/'completion.json');browser=load(meta/'visual-browser-attached.json');cleanup=load(meta/'final-cleanup.json');window=load(meta/'launch-window-process-check.json')
assert result['status']=='failed' and len(result['observations'])==0 and result['failure']['message']=='CDP timeout Page.navigate';assert receipt['exit_code']==1 and receipt['elapsed_seconds']==41.056
assert done['key_reads']==0 and done['external_requests']==0 and done['upstream_events']==[]
assert not cleanup['surviving_recorded_pids'] and not cleanup['profile_matching_processes'] and all(x['absent'] for x in cleanup['owned_ports']);assert not window['remaining_chrome_updater_processes_started_in_owned_window']
assert sha(Path(browser['argv'][0]))==browser['executable_sha256'];assert sha(base/'preparation/exercise.py')==result['script_sha256']
approvals=[]
for mode in ['local','paid']:
 p=Path('/tmp/ch08-operational-'+mode+'-approval.json');a=load(p);approvals.append({**info(p),**{k:a[k] for k in ['approved','scopes','expires_utc','policy_sha256','launcher_sha256']}})
rows=[json.loads(l) for l in (meta/'visual-raw-cdp.jsonl').read_text().splitlines()];navigation=[]
for r in rows:
 if 'out' in r:
  x=r['out'];navigation.append({'direction':'command','id':x['id'],'method':x['method'],'session':x.get('sessionId'),'navigation_url':x.get('params',{}).get('url')})
 else:
  x=r['in'];navigation.append({'direction':'response' if 'id' in x else 'event','id':x.get('id'),'method':x.get('method'),'session':x.get('sessionId'),'has_error':'error' in x})
out.mkdir(exist_ok=True)
for name in ['exercise.py','exercise.sha256','exercise-before-protocol-review.py']:shutil.copy2(base/'preparation'/name,out/name)
shutil.copy2(base/'seal-execution.py',out/'seal-execution.py')
shutil.copy2(base/'seal-execution-v2.py',out/'seal-execution-v2.py')
(out/'packaging-correction.txt').write_text('First assembly failed before report/manifest sealing on a nonexistent guessed filesystem GUI-debug skill path. Skill had been successfully loaded via functions.load_skill and complete appended instructions read before browser execution. Unsealed draft and original packager retained locally. Version2 changes only packaging/provenance and retains both scripts. No browser/native/provider/GUI exercise rerun. An earlier pre-browser cat of that same guessed path also failed; no physical skill-file hash/read is claimed.\n')
for name in ['operational-local-native.receipt.json','operational-local-native.log']:shutil.copy2(root/'.ch08-local'/name,out/name)
for name in ['result.json','launcher-command.json','launcher-exit.json','visual-browser-ownership.json','visual-browser-attached.json','visual-browser-exit.json','final-cleanup.json','launch-window-process-check.json']:shutil.copy2(meta/name,out/name)
for name in ['assets-1.json','ownership-1.json','exit-62716.json','completion.json']:shutil.copy2(run/name,out/('gui-'+name))
write(out/'scope-and-navigation-summary.json',{'status':'FAILED OWNED BROWSER SETUP; no candidate defect established','application_DOM_observations':0,'captured_or_visually_inspected_exercise_screenshots':0,'native_exercise_attempted':False,'nonempty_native_utterances_requested':0,'human_listening_report':None,'paid_exercise_attempted':False,'local_upstream_requests':0,'external_model_requests':0,'key_reads':0,'gui_restarts':0,'intended_page_target_count':1,'cdp_summary_not_raw_trace':navigation,'browser_failure_diagnostics_summary':['Page.navigate navigation-start events observed; no completion response before 30-second controller timeout','Chrome emitted allocator multiple-load warning','Chrome emitted teardown-watchdog-expired diagnostic and exited 2','Chrome automatically spawned bundled GoogleUpdater --install child PID62764, parent62717; reported return3','automatic wake-all reported return0; no update/installer command was issued by exercise controller'],'host_side_effect_qualification':'Not a host-wide no-write/no-network or no-installation audit. Known automatic activity is disclosed; installed state not inspected or undone. Main Chrome executable hash unchanged.','serving_completion_field_qualification':'gui-completion.json browser_launched=false describes sealed serving launcher only. The separately owned controller DID launch muted Chrome; its attached/ownership/exit receipts govern browser activity.','diagnostic_lookup_failure':'One diagnostic read guessed nonexistent controller-final.json; actual completion.json was subsequently located. No exercise retry occurred.'})
write(out/'candidate-and-prior-preservation.json',{'candidate':src['candidate'],'clean':True,'tracked_unchanged':src['tracked'],'retained_artifacts_unchanged':src['unchanged_initial_artifacts'],'old_manifests_verified':old,'old_manifest_entries_all_unchanged':True,'launcher_policy_unchanged':True,'source_build_test_grade_commit_tag_activity':'none'})
write(out/'authority-and-public-pins.json',{'assignment':info(task),'approval_metadata':approvals,'retained_public_read_pins':public,'gui_debug_skill_read_this_turn':{'method':'functions.load_skill(name=gui-debug) succeeded; complete appended instructions read before browser work','physical_file_hash_claimed':False,'qualification':'Guessed on-disk path was absent; loaded instruction provenance is the tool result, not that path'},'no_private_coordinator_review_inputs':True})
local=[]
for p in sorted(base.rglob('*')):
 st=p.lstat();rel=str(p.relative_to(base))
 if stat.S_ISLNK(st.st_mode):local.append({'relative_path':rel,'kind':'symlink-not-followed','target_sha256':hashlib.sha256(os.readlink(p).encode()).hexdigest()})
 elif stat.S_ISREG(st.st_mode):local.append({'relative_path':rel,'kind':'file','bytes':st.st_size,'sha256':sha(p)})
 elif not stat.S_ISDIR(st.st_mode):local.append({'relative_path':rel,'kind':'special-not-read'})
write(out/'local-only-index.json',{'base':str(base),'snapshot_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'qualification':'Raw browser/CDP/API/event/debug logs, browser profiles/caches/crash material stay local. Hashing does not export contents or dereference symlinks. No exercise screenshot was captured.','entries':local})
text=f'''# Chapter 8 — bounded operational attempt (sealed partial evidence)

## Disposition

**Stopped at the first unresolved owned-browser setup failure. Not an application defect finding, operational pass, or overall acceptance.** R1 code remains independently PASS/root code-review accepted as reported by root. Candidate `{src['candidate']['commit']}` / tree `{src['candidate']['tree']}` remains clean. All 126 tracked members and nine retained artifacts are byte-identical. Existing reports, source manifests, and sealed plan/evidence were verified unchanged. No production/test/source edits, rebuilds, grades/lifetime reruns, commits/tags or child agents occurred.

## Authority, scope and actual chronology

Read the complete 7696-byte execution task (SHA-256 `0709ca122d7dc35a5ac4ad6817288d3c0759ace23d25a08bc968cd36f81eca77`), loaded/read GUI-debug guidance, and verified both scoped approvals against the immutable launcher/policy before launch. Earlier full prerequisites were retained only after exact public-file pin verification. No private reviewer/coordinator material was used. Root's execution authority was followed without another routine permission request.

Created a separately pinned, owned CDP controller; the sealed launcher and policy were **not modified**. Before execution, source inspection corrected the new controller's subscriber ordering and added a single default-voice qualification read to its then-unexecuted native branch. The pre-protocol version is retained and was not exercised. Executed controller SHA-256: `{result['script_sha256']}`. Full script, command/receipt, version and pin records are in this package; this run does not validate the controller's unexecuted UI/native/paid branches.

The only execution command was the retained receipt's exact `python3 -B .../exercise.py --phase local-native --script-sha256 ... --run-id local-01`, through the existing clean-environment receipt wrapper. It invoked the exact sealed local serving command and supplied approval. Frozen GUI hash: `{policy['gui_sha256']}`. The launcher verified nine served asset hashes before browser attachment. GUI used its owned directory, owned HOME, dummy Anthropic key, and read-aloud off.

One new **headed, muted** Chrome/profile was started, with loopback dynamic CDP. Version `{browser['version']}`; executable SHA-256 `{browser['executable_sha256']}`. Browser PID62717, profile under the owned `local-01/browser-profile-visual`, CDP127.0.0.1:52107. Exact argv/profile/executable/`ps` and PID-specific `lsof` ownership were durably recorded **before** WebSocket attachment. No8088, shared/personal browser attachment, focus automation, OS-volume/voice/MIC changes, or native callback replacement was performed.

CDP successfully enumerated the owned target, attached a flattened session, and acknowledged Page/Runtime/Network enablement. `Page.navigate` to the owned GUI127.0.0.1:52097 emitted navigation-start/loading events but did not return within 30 seconds. There was **no completed application navigation or application DOM/geometry/control assertion**. No screenshot was captured or visually inspected. The run stopped, issued Browser.close and shut down the owned serving process. Total measured command time: **41.056 seconds**, exit1. No second browser style/profile or paid retry followed.

## Failure, host-side observations and limits

The actual failure is `RuntimeError: CDP timeout Page.navigate`. Chrome later reported a teardown-watchdog expiry and exited **2**. Its retained log also contains an allocator multiple-load warning. These are observed symptoms; the root cause is **unproven**. Nothing here establishes that the frozen application caused the failure.

**Unexpected browser-internal side effect disclosed:** despite the controller's background/component-update disabling flags, Chrome automatically spawned its bundled `GoogleUpdater --install` child (PID62764, parent62717), which reported return3; a later wake-all reported return0. The controller did not invoke an installer/updater or change updater settings. This is not proof that no host-side writes/installation/network activity occurred. No host installed-state audit or undo was attempted, and no shared updater/cache/profile was cleaned. The main Chrome executable's hash remained unchanged. Raw diagnostics remain local and hash-indexed; only this safe observation is exported.

A packaging-only failure attempted to hash a nonexistent guessed GUI-debug filesystem path. The skill was actually loaded via `load_skill`; no filesystem skill hash is claimed. The unsealed attempt and original packager were preserved; a versioned packager corrected provenance without any operational rerun. An earlier pre-browser filesystem lookup of that path also failed.

A subsequent diagnostic guessed the wrong completion filename, then located the actual `completion.json`; that lookup failure is retained as a qualification, not a new exercise attempt. The sealed launcher's `browser_launched=false` field refers only to that serving process: the separate owned controller really did launch Chrome. Browser-specific ownership/attachment/exit receipts govern that fact.

## Evidence gates and counters

| Gate | Actual result |
| --- | --- |
| Frozen pins, scoped approvals, owned serving setup, nine asset hashes | Verified |
| Owned Chrome PID/profile/listener before CDP attachment | Verified |
| Initial application navigation | Failed to complete; unresolved setup/Chrome failure |
| Rendered panes, controls, themes, font/contrast, resizing/narrow layout, two-page synchronization, streaming/tool/pause/scroll/reconnect | **Not exercised**; no DOM/screenshot pass claimed |
| One allowed GUI restart | **0 used** |
| Local fake upstream requests / completed GUI model turns | **0 / 0** |
| Native rate/start/end/cancellation and audible window | **Not attempted**; 0 nonempty utterances requested |
| Human listening | **No report**; no hearing/silence claim |
| Paid model turn / external model requests / real proxy key reads | **Not attempted / 0 / 0** |

Native and paid were correctly gated off by the unresolved earlier failure. The proxy's zero external/model count is not a claim about every possible Chrome background network operation. No provider discovery, fallback, standalone key reading, provider settings export, or credential handling occurred. The owned GUI's dummy credential is not a real provider credential. No paid request budget was consumed.

## Cleanup and preservation

Actual exits: GUI PID62716 **0**, serving launcher PID62709 **0**, Chrome PID62717 **2**, outer controller command **1**. All recorded PIDs are absent, all eight Chrome-log helper PIDs are absent, no process matches the owned profile, and an additional executable/start-window check found no surviving Chrome/updater process from this launch window. Proxy52096, GUI52097 and CDP52107 all reject connection with code61. No broad process kill was used. Profiles, raw traces, logs and crash material remain on disk; no shared cache/profile deletion occurred. The task has **zero surviving owned jobs/processes/listeners**.

## Seal and next boundary

`ch08-operational-exercise-evidence/MANIFEST.json` specifies the exact safe-text portable set, notes-directory-relative base, every included file's size/SHA-256, and explicit self-hash exclusion. Raw CDP/browser/API/history/debug logs, profiles, binary/cache/crash material and any images are local-only and separately hash-indexed; symlinks are not dereferenced. No TAR/PNG/binary/raw trace enters the portable set. The report and package are sealed: no later append/reseal or silent reinterpretation of the failed receipt.

Root can inspect this partial evidence and reconcile the owned Chrome startup/navigation failure before deciding any bounded continuation. Actual GUI, native and model evidence remains open. This submission does not reopen the closed code round, claim operational success, accept/tag the chapter, or start author work.
'''
report.write_text(text)
entries=[]
for p in sorted([report]+[p for p in out.rglob('*') if p.is_file()]):
 b=p.read_bytes();assert b'\0' not in b;b.decode('utf-8');assert p.suffix not in ['.png','.tar','.zip','.dmp','.pyc'];entries.append({'path':str(p.relative_to(notes)),'bytes':len(b),'sha256':sha(p)})
manifest=out/'MANIFEST.json';write(manifest,{'schema':'chapter8-operational-exercise-v1','status':'SEALED FAILED BROWSER SETUP / NATIVE AND PAID UNATTEMPTED','base':str(notes),'candidate_commit':src['candidate']['commit'],'candidate_tree':src['candidate']['tree'],'entries':entries,'all_portable_files_including_manifest':sorted([r['path'] for r in entries]+[str(manifest.relative_to(notes))]),'self_hash_exclusion':'This manifest excludes its own size/SHA to avoid recursion; exact digest supplied out-of-band.','utf8_non_nul_safe_text_only':True,'sealed_no_append':True})
print('SEALED',len(entries)+1,'portable files')
print('REPORT',sha(report));print('MANIFEST',sha(manifest));print('LOCAL INDEX MEMBERS',len(local))
