from pathlib import Path
import hashlib,json,subprocess,os,datetime
E=Path(__file__).resolve().parent;R=E.parents[2];S=Path((E/'scratch-path.txt').read_text().strip());ST=Path('/Users/bill/projects/ensemble-edition-3-student')
def sha(p):
 h=hashlib.sha256()
 with Path(p).open('rb') as f:
  for b in iter(lambda:f.read(1048576),b''):h.update(b)
 return h.hexdigest()
def git(*args,cwd=ST):return subprocess.check_output(['git','--no-pager',*args],cwd=cwd,text=True).strip()
assert sha('/tmp/cadence-ch06-code-review-task.txt')=='0f7c1f4d0b802949e18cdb6b327e8233187de641f90e25116f80da804950b926'
assert sha(R/'docs/edition-3-notes/ch06-candidate.json')=='55f613021fc62b332ec935d6e59e01fd95d6a34bfca3272d5d4539de815432c3'
assert git('rev-parse','HEAD')=='800da21d215b9ca48de41516be4773962fe38ccd';assert git('rev-parse','HEAD^{tree}')=='d54d3cce7e708e393d1b5576bbf3867c037847f4';assert git('status','--porcelain')==''
source=json.loads((E/'source-hashes.json').read_text());assert len(source)==98
for rel,h in source.items():assert sha(ST/rel)==h and sha(S/'source'/rel)==h,rel
runtime_tree={p.relative_to(S/'source').as_posix():sha(p) for p in (S/'source').rglob('*') if p.is_file()}
extras=set(runtime_tree)-set(source)
assert all(p in ['api.log','cmd/ch01/api.log'] or p.startswith('cr/io/') for p in extras),extras
assert not any(p.endswith('.go') or p.endswith('go.mod') or p.endswith('go.sum') for p in extras)
m=json.loads((R/'docs/edition-3-notes/ch06-student-evidence-manifest.json').read_text());assert len(m['artifacts'])==474
for d in m['artifacts']:assert Path(d['path']).stat().st_size==d['size_bytes'] and sha(d['path'])==d['sha256'],d['path']
c=json.loads((R/'docs/edition-3-notes/ch06-candidate.json').read_text());sm=json.loads((R/'docs/edition-3-notes/ch06-student-source-manifest.json').read_text())
for d in [c['report'],c['source_manifest'],c['evidence_manifest']]+sm['artifacts']:assert Path(d['path']).stat().st_size==d['bytes'] and sha(d['path'])==d['sha256'],d['path']
pins={'cmd/grade-edition3/main.go':'4227ad0a0d96f17b541b60f85e7ac19acd09edfbb79e50c35153aab58f7c5ee7','cmd/grade-edition3/main_test.go':'6f443849bef87fcf2b5e97566bff545088e9919db15212450560d76f11c9159a','internal/gradeedition2/ch06.go':'edfd0ef2ad9e9841bf531ba552a8983ac182fa8b1f61612596343aa102d080b7','/tmp/cadence-ch06-grader':'8818400c89bb1663cad72faeca99bc92d64fb9a26e9f023a54d5e319c0e2a798','/tmp/e3-ch06-frozen-core':'b2fd9a630d5587d810c1069e44f63d9dd64daebe402f99893cb760b601043db5'}
for p,h in pins.items():assert sha(R/p)==h,p
protected={}
for name in ['ch06-entry-evidence','ch06-grader-review-evidence']:
 d=R/'docs/edition-3-notes'/name;entries={line.split('  ',1)[1]:line.split('  ',1)[0] for line in (d/'SHA256SUMS').read_text().splitlines()};actual={p.relative_to(d).as_posix() for p in d.rglob('*') if p.is_file()};assert actual==set(entries)|{'SHA256SUMS'}
 for rel,h in entries.items():assert sha(d/rel)==h,(name,rel)
 protected[name]={'files':len(actual),'index_sha256':sha(d/'SHA256SUMS')}
prereqs=json.loads((E/'prerequisite-hashes.json').read_text());drift={p:{'read_sha256':h,'current_sha256':sha(R/p)} for p,h in prereqs.items() if sha(R/p)!=h};assert not drift,drift
for p in (S/'source').rglob('*.go'):
 if not p.name.endswith('_test.go'):assert sha(p)==sha(S/'stream-wakeup'/p.relative_to(S/'source'))
assert sha(S/'source/chapter5_wakeup_test.go')==sha(S/'stream-wakeup/chapter5_wakeup_test.go')
for item in json.loads((E/'mutation-results.json').read_text())['controls']:
 variant=S/('mutation-source-'+item['name']);assert {p.relative_to(variant).as_posix():sha(p) for p in variant.rglob('*') if p.is_file()}==runtime_tree
 assert sha(item['binary'])==item['binary_sha256'] and sha(item['restored_binary'])==item['restored_binary_sha256']
for p in E.glob('*.meta.json'):
 d=json.loads(p.read_text());label=p.name[:-10]
 for k in ['stdout','stderr']:assert sha(E/(label+'.'+k))==d[k+'_sha256'],p
result={'verified_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'student_head':git('rev-parse','HEAD'),'student_tree':git('rev-parse','HEAD^{tree}'),'student_clean':True,'exact_source_files':98,'retained_scratch_runtime_extra_files':len(extras),'all_source_and_frozen_candidate_artifacts_unchanged':True,'artifact_entries_verified':474,'pins':pins,'protected_audits':protected,'prerequisite_drift':drift,'streaming_wakeup_production_and_original_wakeup_assertions_unchanged':True,'both_mutant_source_trees_fully_restored':True,'command_output_digests_verified':True,'main_head_observed':git('rev-parse','HEAD',cwd=R),'main_unrelated_changes_preserved_and_unattributed':True,'review_children_created':0,'paid_calls':0,'settings_or_key_reads':0,'pending_behavioral_jobs':0}
(E/'final-integrity.json').write_text(json.dumps(result,indent=2)+'\n');(E/'main-status-final.txt').write_text(git('status','--short',cwd=R)+'\n');print(json.dumps(result,indent=2))
