import json,hashlib,subprocess,os,shutil,datetime
from pathlib import Path
R=Path('/Users/bill/projects/ensemble');S=Path('/Users/bill/projects/ensemble-edition-3-student');N=R/'docs/edition-3-notes';E=N/'ch08-code-review-r1-evidence';T=Path((E/'scratch-path.txt').read_text().strip());B=json.loads((E/'preservation-baseline.json').read_text());M=json.loads((N/'ch08-code-review-r1-input.json').read_text());sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
def verify(x,base=None):
 p=Path(x['path']);p=p if p.is_absolute() else base/p;assert p.stat().st_size==x['bytes'] and sha(p)==x['sha256'],p
for x in B['inputs']:verify(x)
for x in B['source_members']:verify(x,S)
for x in B['protected_tracked'].values():verify(x)
for x in B['student_portable']['entries']:verify(x,N)
for x in B['student_portable']['local_only_exclusions']:verify(x)
for x in B['source_manifest']['unchanged_initial_artifacts']:verify(x)
for x in B['source_manifest']['old_immutable_package'].values():verify(x)
for x in json.loads((E/'prerequisite-qualification.json').read_text()):assert sha(R/x['path'])==x['sha256']
I=json.loads((N/'ch08-code-review-evidence/inventory.json').read_text());[verify(x,R) for x in I['files']]
old_artifacts=json.loads((N/'ch08-code-review-evidence/artifact-verification.json').read_text())['six_executables'];assert len(old_artifacts)==6
for x in old_artifacts:verify(x)
wire=json.loads((N/'ch08-student-r1-evidence/wire-provenance.json').read_text());verify(wire['raw_wire_local_only']);verify(wire['producer_source'])
G=lambda *a:subprocess.check_output(['git','--no-pager',*a],cwd=S,text=True).strip()
assert G('rev-parse','HEAD')==M['student']['commit'];assert G('rev-parse','HEAD^{tree}')==M['student']['tree'];assert G('rev-parse','HEAD^')==M['student']['parent'];assert G('status','--porcelain=v1','--untracked-files=all')=='';assert G('show-ref','--tags')+'\n'==(T/'student-tags-before.txt').read_text()
assert G('diff','--name-only',M['student']['parent'],M['student']['commit']).splitlines()==['gui/test/correction-r1-page.mjs']
C=json.loads((E/'controls-summary.json').read_text());assert len(C['commands'])==10
for r in C['commands']:
 assert not r['timed_out'] and not r['remaining_process_group'];assert r['exit_code']==r['expected_exit']
 assert sha(E/(r['name']+'.stdout.log'))==r['stdout_sha256'] and sha(E/(r['name']+'.stderr.log'))==r['stderr_sha256']
 try:os.killpg(r['pid_pgid'],0);raise AssertionError(('owned group remains',r['pid_pgid']))
 except ProcessLookupError:pass
try:os.killpg(61610,0);raise AssertionError('prior pager/read process group remains')
except ProcessLookupError:pass
copy_before=json.loads((E/'owned-copy-before.json').read_text());assert copy_before=={str(p.relative_to(T/'copy')):sha(p) for p in (T/'copy').rglob('*') if p.is_file()}
local=[{'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p),'classification':'local-only status snapshot' if p.suffix=='.z' else 'owned local source/status; not portable deliverable'} for p in sorted(T.rglob('*')) if p.is_file()];(E/'local-artifacts.json').write_text(json.dumps({'root':str(T),'files':local,'file_count':len(local),'total_bytes':sum(x['bytes'] for x in local)},indent=2)+'\n')
carried=[]
old_controls=json.loads((N/'ch08-code-review-evidence/controls-summary.json').read_text())['commands']
selected=[r for r in old_controls if r['name'].endswith(('-race','-vet')) or r['name'].startswith('chapter')]
assert len(selected)==12
for name in [r['name'] for r in selected]:
 p=N/'ch08-code-review-evidence'/(name+'.receipt.json');d=json.loads(p.read_text());assert d['exit_code']==0 and not d['timed_out'];carried.append({'path':str(p.relative_to(R)),'bytes':p.stat().st_size,'sha256':sha(p),'exit_code':d['exit_code'],'elapsed_seconds':d['elapsed_seconds'],'qualification':'CARRIED, not rerun; production/source/module/test dependencies except the corrected JS fixture unchanged. Original candidate and artifact metadata remain original.'})
(E/'carried-controls.json').write_text(json.dumps({'unchanged_source_files':125,'only_changed_path':'gui/test/correction-r1-page.mjs','initial_review_inventory_verified_members':len(I['files']),'runtime_identities':old_artifacts,'controls':carried,'fresh_go_suite_runs':0,'fresh_grader_runs':0,'fresh_lifetime_runs':0},indent=2)+'\n')
V={'checked_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'verdict':'PASS — fixture-only R1 resolved','student':M['student'],'input_manifest_sha256':sha(N/'ch08-code-review-r1-input.json'),'inputs_verified':len(B['inputs']),'source_files_verified':len(B['source_members']),'protected_files_verified':len(B['protected_tracked']),'student_portable_entries_verified':len(B['student_portable']['entries']),'prior_review_members_verified':len(I['files']),'six_runtime_identities_verified':len(old_artifacts),'prior_retained_binaries_verified':len(B['source_manifest']['unchanged_initial_artifacts']),'old_student_manifests_unchanged':True,'wire_and_producer_unchanged':True,'candidate_clean':True,'student_tags_unchanged':True,'owned_copy_restored':True,'terminal_command_receipts':10,'remaining_owned_process_groups':0,'pager_read_group_ended':True,'local_file_count':len(local),'local_bytes':sum(x['bytes'] for x in local),'free_bytes':shutil.disk_usage(R).free}
(E/'final-verification.json').write_text(json.dumps(V,indent=2)+'\n')
rows='\n'.join(f"| `{r['name']}` | {r['exit_code']} | {r['elapsed_seconds']:.3f}s | {'Expected behavioral failure' if r['expected_exit'] else 'PASS'} |" for r in C['commands'])
text=f'''# Chapter 8 — fixture-only formal R1 re-review

**Verdict: PASS — R1 resolved; the initial whole-code review’s sole blocker is cleared.**

This closes the narrow fixture-compatibility finding against the exact frozen correction. It does not claim browser/native/live acceptance, create an acceptance tag, alter the original review, or turn carried grades into fresh runs. No additional mandatory correction was found in this one-file review.

## Identity and reviewed scope

- Task SHA-256: `c0f94a707152430d1a6d450df4c86f2d5e589bfba5940c70e4c8029f60f95584`.
- Authoritative input `ch08-code-review-r1-input.json`: `{V['input_manifest_sha256']}`.
- Candidate: `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`; tree `fa896caa404b87b08de61eb38ce079eb38138fbe`; parent `f97a6850468a9cdc6875507f2821ccf6c397654d`.
- The only changed tracked path is `gui/test/correction-r1-page.mjs`, SHA-256 `8adfb4dbd88e883d50be7556ce4b72f137771c1ef140a36c4332ef0004762ae0`. All other 125 source members match the initial reviewed parent. Production, Go tests, module files, grader inputs and APIs are unchanged.
- Required `AGENTS.md`, coding/reviewer skills, and execution instructions were reloaded fully. Other prerequisite reads carry by their exact initial-review hashes, retaining prior full-versus-section qualifications; see `prerequisite-qualification.json`.
- Read the complete student R1 report, changed fixture/diff, source/evidence manifests, wire provenance, relevant production app/workbench/SpeechQueue/tool projection and the unchanged Go producer. This is a fixture-only re-review, not another entry audit or whole-production replay.

## R1 closure rationale

The initial failure and its original seal remain unchanged. The submitted correction now repairs the actual boundary rather than skipping its assertions:

| Location in `gui/test/correction-r1-page.mjs` | Review evidence |
|---|---|
| 10–22, 24–44 | Minimal event/style/media/scroll boundaries cover current production imports. IDs/tags come from actual `index.html`; missing/duplicate IDs, module escape, server-message rejection and unexpected reconnect timers fail explicitly. The actual `app.js` and its modules execute; no alternate production implementation is substituted. |
| 45–55 | A full accepted initial settings snapshot enables the real control. Cancellation dispatches `change`, checks the sparse `update_settings` message, verifies draft rollback/no premature cancellation, and then verifies cancellation on the accepted snapshot. This replaces the obsolete direct `.onchange` assumption. |
| 47, 65–75 | Artifact lookup spans Transcript and Actions but asserts exactly one matching display identity. Tool location is explicitly asserted to be Actions. Visible call/result/error, arguments, opaque exclusion, replay silence/manual readout, pause/unpause and stale-callback protection remain active. |
| 78–86 | Empty/opaque artifacts remain silent. The optional real-wire branch checks current Actions placement, refused batches, on-demand errors, replay/reconnect silence and six retained result artifacts; it is not silently bypassed. |

Fresh standalone and real-wire runs both pass. The workbench and queue fixtures also pass. The correction maintains the original behavioral assertions and strengthens current entrypoint/identity checks; it does not weaken a failing production behavior into a green fixture.

## Independent fresh controls

Node `{json.loads((E/'environment.json').read_text())['node_version']}`; exact executable hash, scrubbed environment, cwd/HOME, argv, output hashes, deadlines, timing and exit statuses are in the indexed receipts. Commands are sequential, with 30-second deadlines and owned process groups. No credentials are inherited.

| Control | Exit | Time | Actual outcome |
|---|---:|---:|---|
{rows}

### Independent executable negative and restoration

The owned-copy mutation removes only the production `SpeechQueue` completion guard that rejects a canceled/stale utterance. The submitted fixture is unchanged in all stages. This is independent of the student’s carried tool-result-projection mutant.

- Old: parse **0**, behavior **0**.
- Negative: parse **0**, behavior **1**, specifically `AssertionError: 'unpause' !== 'pause'` at the stale-callback assertion (fixture line 75). This is executable wrong behavior, not a syntax/import failure.
- Restored: parse **0**, behavior **0**; original/restored `speech.js` SHA-256 `818693251c3bea40bea32e0c10cc8c8dc31a908055c2320dc564546d14cc5ad6`; negative `935b2e5adf6eb8abec8d5b93e08817ac1d18640c98ecb23564cd812b08797344`.

All copied production files restore to their original hashes. `mutation-results.json`, `speech-guard-mutation.diff` and old/negative/restored receipts preserve the distinction.

## Real-wire provenance and limits

The fresh optional branch consumes the pinned local-only `.ch08-local/r1-wire.json` ({wire['raw_wire_local_only']['bytes']:,} bytes; SHA-256 `99a7078ba9a64843b42d40823d1ae0df86fb5b9b15855ed1f5937380d66f635c`). Its producer is unchanged `gui/correction_r1_test.go` (SHA-256 `{wire['producer_source']['sha256']}`), whose carried receipt records the real local Agent/hub/WebSocket path. The producer was not rerun here.

The fresh consumer exercises the named refusal/replay/WebSocket streams through the actual production page handlers. Current settings in the retained replay/WebSocket streams have TTS enabled, so silence assertions are not explained away by a disabled preference. This remains a controlled DOM/socket/synthesis boundary with captured real local wire: **not browser rendering, native speech, permission, microphone, audio/listening or paid-provider evidence**. Raw wire is not copied into the portable bundle.

## Carried evidence; not rerun

`carried-controls.json` records exact hashes for the initial five-module race/vet receipts and both full grade receipts. All six runtime identities and nine retained historical binaries verify unchanged; their original commit/VCS metadata remains original, not relabeled as newly built R1 artifacts. The prior whole-code review’s 132 sealed members are intact.

- Initial root/GUI/consumer/workflow/stream race/vet remain **carried passes**.
- Initial Chapter 8 **100/100** and the single Chapter 1 lifetime **100/100** remain **carried passes** on the originally pinned same CLI, supported by the unchanged production inputs. Fresh R1 grade/lifetime runs: **zero**.
- Initial real-process budget/settings and fractional-rate evidence remains qualified as before. The original failed page fixture remains a historical failure; the new repaired fixture’s pass does not rewrite it.
- Student producer and result-projection sensitivity receipts are carried student evidence, not fresh reviewer runs.
- No browser/native, paid/live provider, secret/key/config, student/child agent, broad Go suite, entry matrix, commit/tag, source/grader/guide/old-report edit, or shared cleanup was performed. No new API or framework recommendation is warranted.

## Failed preparations and preserved qualifications

Read/setup mistakes are preserved, not presented as candidate failures: the initial diff entered a pager because `--no-pager` was omitted; only its uniquely verified owned pager was terminated, the read job ended, and the full diff was recovered with paging disabled. A guessed adjacent filename and a manifest-key assumption caused read-only lookup errors. Two original preparation scripts stopped before controls on a guessed raw-wire path and an incorrect envelope-shape assumption. The corrected preparation derives the path/hash from provenance and validates the actual named streams. `attempt-notes.json`, `preparation-failure.json`, `preparation-envelope-failure.json` and all three scripts retain the sequence. There was no unexpected runtime-control failure or repeat-to-green run.

Both source TARs remain local-only/excluded and unchanged; no TAR, binary payload, raw private wire, bytecode or symlink is exported in this review bundle.

## Preservation, inventory and next boundary

Before/after verification: **20 inputs, 126 source members, 31 protected files, 44 student R1 indexed entries, 132 prior review members, six runtime identities and nine retained binaries** match. Candidate and tags remain unchanged/clean, the one-file diff is unchanged, and the owned copy is restored. All **10** command receipts are terminal, and the earlier pager/read group is also ended.

`ch08-code-review-r1-evidence/inventory.json` is the self-excluded portable hash/size index. `local-artifacts.json` separately records retained owned scratch metadata; no local payload is implicitly called portable. `final-verification.json` and `final-seal.json` are the preservation/process/scope receipts. Only the new R1 report/evidence paths and owned scratch were written.

**Next boundary:** coordinator may record R1/code-review closure using this exact candidate and seal. Subsequent rendered/native/live acceptance remains separate; this review neither performs it nor adds routine permission or architecture prerequisites.
'''
cfg=json.loads((E/'environment.json').read_text());text=text.replace('{cfg_version_placeholder}',cfg['node_version']);(N/'ch08-code-review-r1.md').write_text(text)
print(json.dumps(V,indent=2));print('REPORT WRITTEN',N/'ch08-code-review-r1.md')
