from pathlib import Path
import datetime,difflib,hashlib,json,os,shutil,stat,subprocess
ROOT=Path('/Users/bill/projects/ensemble-edition-3-student');NOTES=Path('/Users/bill/projects/ensemble/docs/edition-3-notes');BASE=ROOT/'.ch08-local/operational-recovery';RUNS=ROOT/'.ch08-local/operational-runs';META=RUNS/'controller-local-recovered-01';GUI=RUNS/'local-recovered-01';OUT=NOTES/'ch08-operational-recovery-evidence';REPORT=NOTES/'ch08-operational-recovery.md'
sha=lambda p:hashlib.sha256(Path(p).read_bytes()).hexdigest()
def save(p,v):Path(p).write_text(json.dumps(v,indent=2)+'\n')
def info(p):p=Path(p);return {'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)}
assert OUT.exists() and not (OUT/'MANIFEST.json').exists() and not REPORT.exists();task=Path('/tmp/cadence-ch08-operational-recovery-task.txt');assert task.stat().st_size==4287 and sha(task)=='13762ae07b849bec074bf058be254b54dddbe19ecaf4e3dc55fb425704aa3d1c'
source=json.loads((NOTES/'ch08-student-r1-source-manifest.json').read_text());commit=subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=ROOT,text=True).strip();tree=subprocess.check_output(['git','--no-pager','rev-parse','HEAD^{tree}'],cwd=ROOT,text=True).strip();assert commit=='b8f22ba5b3c9d1fc52d1a114810b90b77928684d' and subprocess.check_output(['git','--no-pager','status','--porcelain=v1'],cwd=ROOT,text=True)==''
for r in source['tracked']:assert sha(ROOT/r['path'])==r['sha256']
for r in source['unchanged_initial_artifacts']:assert sha(Path(r['path']))==r['sha256']
old=[]
for name in ['ch08-student-evidence-manifest.json','ch08-student-r1-evidence-manifest.json','ch08-operational-plan-evidence/MANIFEST.json','ch08-operational-exercise-evidence/MANIFEST.json']:
 p=NOTES/name;m=json.loads(p.read_text())
 for r in m['entries']:
  if r.get('sha256'):assert sha(NOTES/r['path'])==r['sha256']
 old.append(info(p))
public=json.loads((ROOT/'.ch08-local/public-input-identity.json').read_text())
for r in public['files']:assert sha(Path(r['path']))==r['sha256']
OUT.mkdir(exist_ok=True)
save(OUT/'candidate-and-prior-preservation.json',{'candidate_commit':commit,'candidate_tree':tree,'tracked_unchanged':source['tracked'],'retained_artifacts_unchanged':source['unchanged_initial_artifacts'],'old_manifests_verified':old,'git_clean':True,'source_changed':False,'old_original_failure_remains_failed':True,'all_other_production_and_test_work_unperformed':True})
scopes=[]
for name in ['local','paid']:
 p=Path('/tmp/ch08-operational-'+name+'-approval.json');a=json.loads(p.read_text());scopes.append({**info(p),'scopes':a['scopes'],'expires_utc':a['expires_utc'],'approved':a['approved'],'policy_sha256':a['policy_sha256'],'launcher_sha256':a['launcher_sha256']})
save(OUT/'authority-and-read-qualification.json',{'task':info(task),'full_task_read':True,'approval_records':scopes,'public_reads':'Retained full prerequisite reads and relevant predecessor sections; all nine existing file pins verified unchanged. GUI-debug loaded instruction guidance retained; no nonexistent filesystem skill hash claimed.','public_input_pins':public,'assistance':'Root supplied only the stated prior --use-mock-keychain observation and scoped task; no private logs/reviewer files/other workers or child agents consulted.'})
for name in ['smoke-summary.json','smoke-cleanup.json','final-cleanup.json','preservation-and-window-check.json']:shutil.copy2(BASE/name,OUT/name)
for p in (BASE/'preparation').iterdir():
 if p.is_file():shutil.copy2(p,OUT/p.name)
shutil.copy2(BASE/'seal-recovery.py',OUT/'seal-recovery.py')
shutil.copy2(BASE/'seal-recovery-v2.py',OUT/'seal-recovery-v2.py')
shutil.copy2(BASE/'seal-recovery-v3.py',OUT/'seal-recovery-v3.py')
(OUT/'second-packaging-failure-note.txt').write_text('Second assembly stopped on an existing unsealed draft subdirectory (FileExistsError). The following verifier correctly could not find a nonexistent manifest; a trailing status/disk command made the outer shell exit0, which is not a packaging pass. Draft2 and packager-v2 are preserved. Version3 changes only directory reuse for this new unsealed package. Final execution uses set -e. No operational rerun.\n')
(OUT/'packaging-failure-note.txt').write_text('First assembly stopped before report/manifest sealing on guessed exit-1.json. Actual exit receipt is PID-named; controller.json is the actual serving metadata. Original packager and unsealed draft are preserved. Version2 corrects receipt selection/provenance only. No browser, controller, native, provider, source or test rerun.\n')
for label in ['original-options','plus-mock-keychain']:
 d=BASE/('smoke-'+label);dest=OUT/('smoke-'+label);dest.mkdir(exist_ok=True)
 for name in ['result.json',label+'-browser-ownership.json',label+'-browser-attached.json',label+'-browser-exit.json']:shutil.copy2(d/name,dest/name)
for name in ['result.json','launcher-command.json','launcher-exit.json','visual-browser-ownership.json','visual-browser-attached.json','visual-browser-exit.json']:shutil.copy2(META/name,OUT/('recovered-'+name))
for name in ['ownership-1.json','assets-1.json','controller.json','completion.json','exit-'+str(json.loads((GUI/'ownership-1.json').read_text())['pid'])+'.json']:shutil.copy2(GUI/name,OUT/('serving-'+name))
shutil.copy2(GUI/'cr/settings.json',OUT/'actual-owned-gui-settings.json')
for name in ['operational-recovery-smoke','operational-recovered-local-native']:
 for suffix in ['.receipt.json','.log']:shutil.copy2(ROOT/'.ch08-local'/(name+suffix),OUT/(name+suffix))
# Export only preference/state metadata from the raw actual app frames; keep full CDP traffic local.
settings=[]
for line in (META/'visual-raw-cdp.jsonl').read_text().splitlines():
 row=json.loads(line);candidates=[row]+[v for v in row.values() if isinstance(v,dict)]
 for m in candidates:
  if m.get('method')!='Network.webSocketFrameReceived':continue
  try:v=json.loads(m['params']['response']['payloadData'])
  except (ValueError,KeyError):continue
  if v.get('type') in ['current_settings','settings_changed','state_changed','replay_end']:
   settings.append({'session_id':m.get('sessionId'),'type':v['type'],**{k:v[k] for k in ['settings','to','paused','reason'] if k in v}})
save(OUT/'safe-settings-state-wire-summary.json',{'records':settings,'qualification':'Selected actual received preference/state frames only; raw complete CDP trace retained locally.'})
result=json.loads((META/'result.json').read_text());n=next(x['data'] for x in result['observations'] if x['name']=='740px narrow actual layout');screens=[]
for p in META.glob('*.png'):screens.append({**info(p),'portable':False,'captured':True,'visually_inspected_by_student':False})
script=BASE/'preparation/controller-v2.py';failed_lines=[{'line':i,'source':s} for i,s in enumerate(script.read_text().splitlines(),1) if "narrow['panes']" in s or "narrow=b.dom" in s or "width':740" in s]
save(OUT/'failure-and-image-evidence.json',{'status':'failed exercise assertion; product defect not adjudicated','exact_failed_check':"assert narrow['panes']['.actions-pane']['x']+narrow['panes']['.actions-pane']['w']<=741",'script':info(script),'source_lines':failed_lines,'viewport':n['dom']['viewport'],'actions_rectangle':n['dom']['panes']['.actions-pane'],'right_edge':872,'screenshot_refs':screens,'captured_not_viewed':True,'ordering':'Emulation.setDeviceMetricsOverride(740x820), immediate DOM evaluation, screenshot capture, then assertion on previously retained DOM; no frame-stability wait.','interpretation':'Shows an Actions rectangle outside viewport at that measurement. Does not establish permanent clipping or usable/unusable page horizontal scrolling. Parent overflowX/scrollWidth/scrollLeft and settled post-resize geometry were not captured.','frozen_source_read_only_context':[{'path':'gui/web/gui/style.css','sha256':sha(ROOT/'gui/web/gui/style.css'),'facts':'#workbench min-width:652px; grid includes minmax(260px,1fr); no explicit body/workbench horizontal clipping rule; <=900px stylesheet defaults may be overridden by inline divider widths.'},{'path':'gui/web/gui/workbench.js','sha256':sha(ROOT/'gui/web/gui/workbench.js'),'facts':'Divider sets inline widths and registers window resize clamp handlers. Asynchronous reconciliation is plausible but unmeasured; not asserted as cause.'}]})
save(OUT/'scope-accounting.json',{'diagnostic_browser_launches':2,'application_browser_launches':1,'application_pages':2,'application_browser_mode':'headed, muted, fresh owned profile, plus --use-mock-keychain','gui_restarts':0,'local_provider_requests':0,'native_window_started':False,'native_nonempty_utterances_requested':0,'native_rate_start_end_checks':0,'human_listening_report':None,'observed_queue_cancel_records_in_last_DOM':len(n['dom']['speechEvents']),'qualification_cancel':'Existing app cancels empty speech queue on accepted false settings snapshots. Those cancel records are not a native utterance/rate/audibility exercise.','paid_attempted':False,'paid_model_requests':0,'real_proxy_key_reads':0,'all_upstream_events':json.loads((GUI/'completion.json').read_text())['upstream_events'],'completion_scope':'Serving completion browser_launched:false describes only sealed serving launcher; separately owned browser launches are proved by controller ownership/attachment records.','no_app_native_api_override':True,'no_production_or_sealed_launcher_change':True,'no_tests_builds_grades_or_lifetime_reruns':True,'no_children':True,'no_retries_after_exercise_assertion':True})
(OUT/'diagnostic-inspection-note.txt').write_text('One post-failure read-only inspection guessed visual-cdp-raw.jsonl and returned FileNotFoundError. Actual owned inventory located visual-raw-cdp.jsonl. No operational rerun or source change followed; this inspection failure is retained as a limitation, not hidden as a successful command.\n')
# Local-only index never dereferences profile symlinks.
entries=[]
for prefix in [BASE,META,GUI]:
 for current,dirs,files in os.walk(prefix,followlinks=False):
  for name in dirs+files:
   p=Path(current)/name;st=p.lstat();r={'path':str(p),'bytes':st.st_size}
   if stat.S_ISLNK(st.st_mode):r.update(kind='symlink-not-followed',target_sha256=hashlib.sha256(os.readlink(p).encode()).hexdigest())
   elif stat.S_ISREG(st.st_mode):r.update(kind='file',sha256=sha(p))
   elif stat.S_ISDIR(st.st_mode):continue
   else:r.update(kind='special-not-read')
   entries.append(r)
save(OUT/'local-only-index.json',{'as_of_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'roots':[str(BASE),str(META),str(GUI)],'raw_browser_profiles_CDP_logs_PNG_API_debug_material_local_only':True,'entries':entries})
s=json.loads((BASE/'smoke-summary.json').read_text());a,b=s['cases'];report=f'''# Chapter 8 bounded operational recovery — partial / stopped

## Status and immutable identity

The authorized comparison found an observed-working Chrome configuration, but the **single recovered application exercise FAILED an assertion** and stopped. Original operational report remains **FAILED and immutable**. This is not a production-defect adjudication, full GUI/native/model pass, independent closure, overall acceptance or tag.

Student commit **{commit}**, tree **{tree}**, is clean. All **126** tracked files and **nine** retained executables match their pins. All prior package seals and all **317** originally indexed local entries are unchanged. No production/test/sealed-launcher/policy edits, rebuilds, grades, long-test reruns, source commits/tags, children or author/Chapter 9 work occurred.

Authority: fully read recovery assignment, 4287 bytes, SHA-256 `13762ae07b849bec074bf058be254b54dddbe19ecaf4e3dc55fb425704aa3d1c`. Both root approvals were valid through 20:40:25 UTC. Existing full prerequisite reads qualified by unchanged exact public pins; GUI-debug's loaded instructions were retained, without inventing a filesystem skill pin. Root's single historical flag observation is acknowledged; no private diagnostic files or other workers were read.

## Two muted smoke launches — exact comparison, limited causality

One owned trivial HTTP peer served the same **136-byte** page (SHA-256 `{s['body_sha256']}`) to two fresh, sequential, muted Chrome profiles, using pinned controller-v1 and smoke script. Exact executables, versions/hashes, argv, profile, PID/listener ownership **before** attachment, environment-name sets, timings, HTTP observations and cleanup are retained.

| Configuration | Navigation / actual evidence | Outcome |
|---|---|---|
| Original options, without added keychain flag | `Page.navigate` timed out at its 30-second bound; **zero HTTP requests** reached peer; no DOM/response proof. Case elapsed **{a['elapsed_seconds']}s**. | Failed; Chrome exited 2. |
| Same options plus `--use-mock-keychain` | Navigation returned in **{b['navigation_seconds']}s**; actual DOM marker/title/body matched; browser retrieved identical 136-byte HTTP 200 response. Separate favicon request returned intentional 404. Case elapsed **{b['elapsed_seconds']}s**. | Passed; Chrome exited 0. |

The paired runner exited 0 because a working case existed; this does **not** turn the original case into a pass. It used a fresh profile/HOME for each case, in this order, at different moments. The observed difference is consistent with a configuration/environment sensitivity, **not proof that keychain handling caused the original failure or that the flag generally fixes it**. Allocator warnings appeared in both successful and failed Chrome logs, so that warning alone is not a failure classifier.

Original-options Chrome again automatically spawned its bundled updater install attempt (reported return 3; wake-all return 0). No installer/updater command was issued by these controllers and no shared updater/OS configuration was cleaned or changed by them. No install return was observed in the added-flag/app logs; absence of a log is not a host-wide no-write/no-network audit. All logs stay local.

## Single recovered exercise and exact stop

Controller-v2 SHA-256 **{sha(script)}** was pinned before its only run. New mutable controller plumbing selected the observed-working flag; before execution it also prepared a real-overflow scroll probe and a fixed-phrase card-identity locator. Full diffs are included. **Those later scroll/native refinements were not reached or credited.** No assertion was weakened or retried after failure.

The frozen GUI and its served asset hashes were checked by the unchanged serving launcher. New isolated GUI HOME/data/workspace, dummy GUI credential, owned loopback endpoints, owned headed **muted** browser/profile and two app pages were used. No shared browser/live8088 or focus automation was used.

Actual DOM/control evidence reached:
- Connected idle root tree and the three panes; no fabricated Agent fleet.
- Keyboard tab selection focused Settings and exposed its panel.
- Actual divider keyboard input changed sidebar **240 → 250**, and pointer drag changed Actions **320 → 350**.
- Light/dark theme states and computed colours; browser-emulated system light/dark preference changes. OS appearance was not changed.
- Accepted fractional speed **1.2345**, font control **20**, `tts_enabled:false`, and `max_tool_rounds:0` acknowledgement/save. Owned settings file is included. No text artifacts existed yet, so visible article font scaling was **not** tested.
- Light, dark and narrow screenshots were captured. **All three are local and hash-indexed; none was visually inspected by this student.** DOM/computed-style readings are not a human visual or contrast verdict.

The run exited **1** after **2.485s**, at:

```python
assert narrow['panes']['.actions-pane']['x'] + narrow['panes']['.actions-pane']['w'] <= 741
```

Retained viewport: **740 × 820**. Sidebar: **x=0,w=250**; Chat: **x=256,w=260**; Actions: **x=522,w=350**, right edge **872** (132px beyond 740). `failure-and-image-evidence.json` includes exact controller lines, frozen source references and the complete measured rectangle.

**Clipping versus intended scrolling is unresolved.** The controller evaluated DOM immediately after `Emulation.setDeviceMetricsOverride`, then captured `narrow.png`, then asserted against the earlier measurement. It did not wait for frame/resize stability. Frozen CSS has a 652px minimum workbench width and no explicit body/workbench horizontal clipping rule; frozen JS registers resize handlers that clamp divider widths. These rules are consistent with either document overflow at that instant or pending width reconciliation, but are not evidence of the actual later outcome. Parent overflow-X/scroll width/scroll position and settled geometry were not captured; page horizontal-scroll usability was not exercised. The rectangle proves out-of-viewport bounds at measurement time, **not permanent clipping or an established product bug**. Root should reconcile the local PNG and timing/public contract, not silently reinterpret this failed assertion as a pass.

No new browser/controller/exercise run occurred after this assertion. Read-only evidence/source inspection followed. One guessed trace filename lookup failed before the actual inventory located it; this diagnostic inspection failure is recorded without inventing an execution receipt or rerunning anything.

## Budgets and unattempted gates

- Exactly **two** diagnostic browser launches and **one** recovered muted application browser session; no extra mode/flag sweep.
- **0** application GUI restarts; **0** local provider requests / admitted tools / terminal prompts. Streaming, busy-tree, tool routing, pause/hint overlap, real overflowing-pane scrolling and restart checks were not reached.
- Native window **not started**; **0** authored nonempty utterances and **0** native requested-rate/getter/start/end observations. The empty queue's ordinary cancel records appeared during settings snapshots; they are not a native-rate or audibility exercise. No voice/locality claim or human listening report exists. No OS volume/default-voice/MIC change was performed.
- Paid **not attempted**; **0** paid MODEL requests and **0** real-key reads. No credential handling, provider discovery/fallback, paid retry, or standalone one-file model call. Original 4-request/90-second ceiling was not consumed or exceeded. Chrome host background activity is separately qualified above; proxy/model counters are not a blanket network audit.

A first packaging assembly failed before sealing on a guessed GUI exit filename. Its unsealed draft and original packager are retained; a versioned packager used the actual PID-named exit receipt and `controller.json`. A second assembly stopped on an already-created unsealed subdirectory; its subsequent verifier found no manifest. That attempt and packager are also retained. The final version made new-draft directory handling idempotent and used fail-fast execution. These changes affected packaging only, not the exercise or assertion.

The failure stopped all later gates as instructed. Prior code acceptance remains a distinct fact; native and model evidence remain open.

## Cleanup, preservation and sealed packet

Original-options Chrome exited 2; added-flag smoke and recovered Chrome exited 0. Recovered GUI and serving launcher exited 0. The trivial peer shut down and joined. **13** recorded GUI/launcher/browser/log-descendant PIDs are absent; **six** owned HTTP/GUI/proxy/CDP listeners are closed. No owned-profile survivor or Chrome/updater process from either launch window remains. No shared process was terminated or shared cache/profile/artifact deleted. Profiles, dumps, logs and screenshots were preserved, not relocated or discarded.

Only new main outputs are this report and `ch08-operational-recovery-evidence/`. The latter contains safe text: pinned executed controllers and diffs, both runner receipts and complete stdout, per-case HTTP/DOM evidence, application measurements, selected preference/state frames, actual owned preferences, approvals/pins, cleanup and full local-only artifact index. Raw browser/CDP/API/debug/profile/crash material and all PNG/archive/binary content remain local. `MANIFEST.json` inventories every portable file including this external report; its self-hash is excluded only inside itself and supplied in submission.

Disposition: sealed **partial recovery with failed application assertion**, not a successful full recovery. No further exercise is planned within this bounded run. Root can independently reconcile the exact screenshots/timing/contract while preserving both original and recovery failures and the frozen candidate.
'''
REPORT.write_text(report)
files=sorted(p for p in OUT.rglob('*') if p.is_file())+[REPORT];entries=[]
for p in files:
 raw=p.read_bytes();raw.decode('utf-8');assert b'\0' not in raw;assert p.suffix.lower() not in {'.png','.tar','.zip','.gz','.exe','.dmp'};entries.append({'path':str(p.relative_to(NOTES)),'bytes':len(raw),'sha256':sha(p)})
manifest={'schema':'ch08-operational-recovery-evidence-v1','sealed_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'base':str(NOTES),'candidate_commit':commit,'candidate_tree':tree,'status':'diagnostic selected observed-working configuration; single recovered application assertion FAILED; native/paid not attempted','entries':entries,'all_portable_files_including_manifest':sorted([r['path'] for r in entries]+['ch08-operational-recovery-evidence/MANIFEST.json']),'self_hash_exclusion':'MANIFEST.json is named in the exact inventory; its SHA-256 is provided out of band, not recursively embedded.'};save(OUT/'MANIFEST.json',manifest)
assert {str(p.relative_to(NOTES)) for p in OUT.rglob('*') if p.is_file()}|{'ch08-operational-recovery.md'}==set(manifest['all_portable_files_including_manifest'])
assert subprocess.check_output(['git','--no-pager','status','--porcelain=v1'],cwd=ROOT,text=True)==''
print(json.dumps({'report_sha256':sha(REPORT),'manifest_sha256':sha(OUT/'MANIFEST.json'),'portable_files':len(manifest['all_portable_files_including_manifest']),'local_index_entries':len(json.loads((OUT/'local-only-index.json').read_text())['entries']),'selected_state_frames':len(settings),'source_clean':True},indent=2))
