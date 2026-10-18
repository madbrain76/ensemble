import datetime, hashlib, json, os, pathlib, shutil, subprocess, tarfile
root=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');local=root/'.ch09-local'
notes=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes')
evidence=notes/'ch09-student-evidence'
evidence_manifest=notes/'ch09-student-evidence-manifest.json'
manifest_path=notes/'ch09-student-source-manifest.json';review_path=notes/'ch09-student-review.md'
def git(*args):return subprocess.check_output(['git','--no-pager',*args],cwd=root).decode().strip()
def digest(path):return hashlib.sha256(path.read_bytes()).hexdigest()
def read(name):return json.loads((local/name).read_text())
commit=git('rev-parse','HEAD');tree=git('rev-parse','HEAD^{tree}')
assert commit=='251ce88fb1677897dbcc6b33661096c5cb39cd8e'
assert git('status','--porcelain')==''
assert not manifest_path.exists() and not review_path.exists() and not evidence.exists() and not evidence_manifest.exists(), 'refusing to overwrite prior evidence'
for chapter in ('09','01'):
 assert read(f'grade2-ch{chapter}.receipt.json')['exit_code']==0
 assert read(f'grade2-ch{chapter}.stdout.txt')['passed'] is True
required=['root','root-race','root-vet','gui','gui-race','gui-vet','consumer','consumer-race','consumer-vet','workflow','workflow-race','workflow-vet','stream','stream-race','stream-vet']
for name in required:assert read('final2-'+name+'.receipt.json')['exit_code']==0
assert all(row['credited'] for row in read('sensitivity-summary-final2.json'))
assert (local/'binary2-hashes-before-adapters.txt').read_bytes()==(local/'binary2-hashes-after-adapters.txt').read_bytes()
files=git('ls-files','-z').split('\x00')
tracked=[dict(path=name,bytes=(root/name).stat().st_size,sha256=digest(root/name)) for name in files if name]
# Prove restored copies are the same tracked source, not merely that one file matches.
for name in ('final2-hidden-dispatch','final2-first-source'):
 copy=local/('control-'+name)
 for row in tracked:assert digest(copy/row['path'])==row['sha256'], (name,row['path'])
with tarfile.open(local/'source-r2.tar') as archive:
 archive_members=[dict(path=m.name,bytes=m.size,sha256=hashlib.sha256(archive.extractfile(m).read()).hexdigest()) for m in archive.getmembers() if m.isfile()]
assert sorted(archive_members,key=lambda x:x['path'])==sorted(tracked,key=lambda x:x['path'])
binaries=[]
for name,module,target in [('cli','.', './cmd/ch01'),('gui','gui','./cmd'),('consumer','examples/ch04-consumer','.'),('workflow','examples/ch06','.'),('stream','examples/ch07','.')]:
 path=local/('ch09-'+name+'-r2');assert read('build2-'+name+'.receipt.json')['exit_code']==0
 binaries.append(dict(role=name,path=str(path),bytes=path.stat().st_size,sha256=digest(path),module=module,target=target,build_receipt='build2-'+name+'.receipt.json'))
manifest=dict(format_version=1,status='frozen-local-candidate; root browser/provider continuation pending',created_utc=datetime.datetime.now(datetime.timezone.utc).isoformat(),repository=str(root),predecessor=dict(tag='edition-3-ch08-r1',tag_object='bbae919d6c595f879eb27f602155c975e1cb8489',tree='fa896caa404b87b08de61eb38ce079eb38138fbe',commit='b8f22ba5b3c9d1fc52d1a114810b90b77928684d'),commit=commit,tree=tree,parents=git('show','-s','--format=%P','HEAD').split(),worktree_clean=True,tracked_files=tracked,binaries=binaries,grader=dict(path='/tmp/cadence-ch09-grader',sha256=digest(pathlib.Path('/tmp/cadence-ch09-grader'))),source_archive=dict(path=str(local/'source-r2.tar'),bytes=(local/'source-r2.tar').stat().st_size,sha256=digest(local/'source-r2.tar'),members=archive_members),source_diff=dict(path='source-diff-r2.txt',sha256=digest(local/'source-diff-r2.txt')),superseded_candidate=dict(commit='f66cfd1a48e85e91ce6ff122a134a673d1ef0ea0',archive_sha256=digest(local/'source.tar'),reason='direct sibling implementation import rejected by inherited ownership parity; replaced with injected interface'),evidence_index=str(evidence_manifest),models=read('media-models.json'),toolchain=(local/'toolchain.txt').read_text(),public_inputs='public-input-hashes.txt',manual_browser='NOT RUN: requires finite root authorization',paid_provider='NOT RUN: requires finite root authorization')
manifest_path.write_text(json.dumps(manifest,indent=2)+'\n')
runtime=read('final2-root.receipt.json')['duration_seconds'];race=read('final2-root-race.receipt.json')['duration_seconds']
g9=read('grade2-ch09.stdout.txt');g1=read('grade2-ch01.stdout.txt')
review=f'''# Edition 3 Chapter 9 — frozen local candidate

## Identity and scope

- Predecessor: `edition-3-ch08-r1`, `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`.
- Final candidate: `{commit}`; tree `{tree}`. Worktree clean.
- Exact tracked source, archives, executable SHA-256s, construction targets and public input pins: `ch09-student-source-manifest.json`.
- Portable receipts: `ch09-student-evidence-manifest.json`. `.ch09-local` retains local executables and both source archives; binaries/archives are deliberately not copied into portable UTF-8 evidence.
- Local work only. No browser, paid provider, credentials, children, branch/tag promotion, push, or acceptance claim. No private grader/audit/history/reference inputs were read. The public Chapter 9 document and its current Edition 3 recipe were the contract.

## Ownership and behavior

The root constructs a per-Agent registry, installs lifecycle declarations, renders the primary, then publishes the Agent. The Agent sees an injected lifecycle interface, not the registry implementation. The registry has immutable discovered definitions and separately locked loaded/enabled/offered/provenance/pending/variable state; it has a parent Agent. No mutable application globals were added. CLI and GUI translate the three chapter environment variables into explicit `SkillSpec` input. No-spec creation preserves inherited registration.

Frontmatter names, not directory names, identify skills. Parsing supports the specified scalar and indented dash-list fields, detects malformed input with paths, and leaves later Markdown rules in the body. Activation validates the complete dependency closure first; dependencies precede dependents. Declaration filtering and pre-control dispatch gating share enabled state. Hidden calls do not allocate jobs or execute effects. The lifecycle uses ordinary paired tool history, with actual `tool_called` sequence numbers and first-enabling sources. Repeated loads preserve those sources. Unload is pending only; dynamic writes remain executable. The primary text stays fixed. Rendering snapshots callbacks/state under lock and invokes callbacks after unlock, preserving unknown variables. Existing engine request rendering and GUI history/wire paths are reused.

## Final-source validation

All commands, exact cwd/environment, start time, duration, exit, stdout/stderr and disk headroom are in individual JSON receipts. Receipt-wrapper processes receive PATH/HOME/TMPDIR plus offline Go settings only, not provider environment. No suite was interrupted or shortened. Guards required at least 1 GiB free; receipts show actual headroom.

- Root `go test -count=1 ./...`: pass, {runtime:.2f}s. Root race suite: pass, {race:.2f}s. Root vet: pass.
- GUI, consumer, workflow, and stream modules: tests/race/vet all exit zero. Consumer and stream report no test files; their build and adapter checks are the applicable coverage, not invented unit-test counts.
- Five final executable builds passed, with explicit module/target identity and unchanged pre/post-adapter hashes.
- Chapter 9 adapter: {g9['score']}/{g9['max_score']}, `passed: true`, exit zero. Chapter 1 adapter: {g1['score']}/{g1['max_score']}, `passed: true`, exit zero, same final CLI.
- Focused tests cover parsing, duplicate identities, initial/dynamic dependency behavior, missing/cyclic/unknown/reachable-vs-hidden loads and atomicity, real write/process effects, hidden controls, repeat/provenance, pending dynamic write, Agent isolation, reentrant callbacks, concurrent snapshots, legacy tools, and fixed system/dynamic declarations over actual local HTTP requests.
- The GUI test connects a real local WebSocket through the production hub, validates launch-variable translation, hidden write refusal, dynamic rendered instructions in the normal wire transcript, pending unload, and a real post-unload write. This is local wire integration, **not** browser evidence.
- Two final-source disposable controls are compile-valid behavior negatives: removing dispatch admission/lookup guards makes the hidden writer actually run; overwriting first-source metadata fails the provenance assertion. Each exact-restored copy passes, and every tracked file of both restored copies matches the final source manifest. See `sensitivity-summary-final2.json` and scripts for complete mutations.

## Preserved failures and repairs

1. Initial pre-freeze root compile failed because a scripted field insertion missed the predecessor's spacing. The insertion and slice-based wiring were repaired; the full inherited suite then passed (129.048s reported by Go). Early merged logs remain; their wrapper start/duration was not measured and is not retroactively invented.
2. A metadata `printf` command failed before hashing; corrected immediately. One guessed public brief path was absent; the prescribed path was recovered from the public assignment. A later metadata check guessed a nonexistent predecessor tag; it was corrected to the assignment's annotated edition-3-ch08-r1 and commit ancestry verified. None changed source.
3. First GUI build used nonexistent `./cmd/ensemble-gui`. That exit-one receipt is retained; corrected `./cmd` and all final builds passed.
4. The first dispatch sensitivity mutation caused a nil-job panic, and its second revision did not compile (`RLock` on a Mutex). Both are explicitly uncredited; both were restored and passed. The corrected third control produced the intended real-write failure and restored pass. Final-source controls were then rerun from the second frozen archive.
5. First frozen adapter attempt scored 85/100, `passed: false`, exit one: all skill checks passed, but inherited star-topology parity correctly rejected Agent's concrete sibling registry import. This was an ownership defect, not an accepted exception. Agent now depends on an injected interface; the root still owns concrete construction. The new commit was rebuilt, both adapters rerun, and all local final gates repeated. The first archive, logs and report are preserved and labelled superseded. Its fail-fast command did not run Chapter 1.

## Assistance, interpretation and limits

Assistance was the public operational assignments, root-approved plan, repository AGENTS/coding-skill guidance, current Edition 3 guides, public chapter, and pinned black-box diagnostics. No child/reviewer or private reference input was used. Implementation iterations above are not formal independent-review rounds; that review remains root-owned. The earned design lesson is the concrete sibling-import defect: injection through an interface is necessary even when the registry is per-Agent. The chapter leaves exact exported API spelling and callback error signatures to implementation; the documented `SkillSpec`/read-only `SkillRegistry` API supplies that contract without claiming a new chapter ambiguity. Frontmatter remains deliberately narrower than YAML. A successful finite sensitivity check is not universal proof, and no reproducible-build claim is made; embedded build metadata and exact current binary hashes are retained.

## Browser/provider continuation — NOT EXECUTED

Root must supply its later finite authorization. Use final CLI/GUI hashes from the manifest; do not rebuild silently. Prepare a fresh scratch skill directory with `base` (primary/read_file/offers code-tools), `code-tools` (loadable/write_file/depends helper), and `helper` (dependency/list_directory), matching the public chapter. Launch the GUI from scratch with `EN_SKILLS_DIR`, `EN_PRIMARY_SKILL=base`, and a distinctive `EN_CUSTOM_VAR`; use the frozen GUI executable and root-approved local fixture/model endpoint.

The local-only fixture should attempt a hidden write, load code-tools, write/read a scratch note, unload, then **write again using the dynamic write_file capability**. Capture clean request bodies showing initial versus expanded declarations and identical system text. The real browser should connect through the production GUI, visibly show the skill instructions and pending result in ordinary tool panels, inspect the real note content, and exercise pause/Stop, a later ask, and shutdown with a bounded delayed fixture where needed. Check responsive layout/console through the approved browser tools. Keep screenshots and sanitized console/network evidence distinct from local Go/WebSocket tests. No dummy browser or provider success is claimed here.

Paid CLI provider work, if authorized, is limited to the public chapter graph and a finite root-selected matrix. Current declared identities used locally are claude-sonnet-4-6, gpt-4.1-2025-04-14 and gemini-3.8-flash; root must confirm account availability and budget. Inspect the real note and post-unload dynamic write, retain credential-free request structure and exact exits, and terminate owned processes. No credentials have been opened or calls made by this student.

## Status

Frozen local implementation/evidence complete. No genuine implementation blocker remains. Manual browser/provider validation and acceptance/promotion are intentionally pending root action; there is no new ownership plan required.
'''
review_path.write_text(review)
local_rows=[dict(path=str(p.relative_to(local)),bytes=p.stat().st_size,sha256=digest(p)) for p in sorted(local.rglob('*')) if p.is_file() and p.name!='local-artifact-index.json']
(local/'local-artifact-index.json').write_text(json.dumps(dict(format_version=1,root=str(local),excludes_self=True,qualifiers='Local-only binaries, source archives and restored disposable copies are not portable evidence. Earlier candidates and uncredited controls are retained explicitly; final2 receipt prefixes and source-r2.tar bind the final candidate.',files=local_rows),indent=2)+'\n')
evidence.mkdir(parents=True)
for path in sorted(local.iterdir()):
 if path.is_file() and path.suffix in ('.txt','.json','.py'):
  content=path.read_bytes();content.decode('utf-8');assert b'\0' not in content
  shutil.copyfile(path,evidence/path.name)
shutil.copyfile(manifest_path,evidence/'source-manifest.json');shutil.copyfile(review_path,evidence/'review.md')
rows=[dict(path=str(p.relative_to(notes)),bytes=p.stat().st_size,sha256=digest(p)) for p in sorted(evidence.iterdir()) if p.is_file()]
rows.extend(dict(path=p.name,bytes=p.stat().st_size,sha256=digest(p)) for p in (manifest_path,review_path))
(evidence_manifest).write_text(json.dumps(dict(format_version=1,source_commit=commit,source_tree=tree,review_sha256=digest(review_path),manifest_sha256=digest(manifest_path),final_receipt_prefixes=['final2-','build2-','grade2-','control-final2-'],superseded_or_prefreeze='Other receipts are retained failures or earlier-candidate evidence, not final-source credit.',binary_policy='UTF-8 evidence only; executable and archive hashes retained in manifest',files=rows),indent=2)+'\n')
print(json.dumps(dict(manifest=str(manifest_path),manifest_sha256=digest(manifest_path),review=str(review_path),review_sha256=digest(review_path),index=str(evidence_manifest),index_sha256=digest(evidence_manifest),artifacts=len(rows),tracked_files=len(tracked)),indent=2))
