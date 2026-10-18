from pathlib import Path
import subprocess, json, hashlib, datetime, sys
R=Path('/Users/bill/projects/ensemble-edition-3-student')
O=Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-r1-evidence')
B=R/'.ch07-local/correction-r1/frozen'; RAW=R/'.ch07-local/correction-r1/raw'
F=json.loads((O/'frozen-source.json').read_text())
OLD=json.loads(Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-harness-reconciliation-evidence/repaired-chapter7-invocation.json').read_text())
ENV=OLD['environment']; rows=[]
def sha(p): return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def now(): return datetime.datetime.now(datetime.timezone.utc).isoformat()
def clean():
 assert subprocess.check_output(['git','--no-pager','rev-parse','HEAD'],cwd=R,text=True).strip()==F['commit']
 assert not subprocess.check_output(['git','--no-pager','status','--porcelain=v1'],cwd=R,text=True).strip()
 for f,h in F['source_files'].items(): assert sha(R/f)==h,f

def run(name,argv,cwd,extra=None,want_failure=None,json_report=False):
 clean(); stdout=O/(name+('.json' if json_report else '.stdout.log'));stderr=O/(name+'.stderr.log')
 assert not stdout.exists() and not stderr.exists(),name
 env=dict(ENV);env.update(extra or {});row={'name':name,'command':list(map(str,argv)),'cwd':str(cwd),'environment':env,'start':now(),'candidate':F['commit'],'source_tree':F['tree'],'stdout':str(stdout),'stderr':str(stderr),'external_timeout':None}
 rows.append(row);(O/'sequential-in-progress.json').write_text(json.dumps(rows,indent=2)+'\n');print('START',name,flush=True)
 # No subprocess timeout, including genuine delayed lifetime phases.
 with stdout.open('wb') as out,stderr.open('wb') as err:p=subprocess.run(row['command'],cwd=cwd,env=env,stdout=out,stderr=err)
 row.update(exit_code=p.returncode,end=now(),stdout_sha256=sha(stdout),stderr_sha256=sha(stderr));(O/(name+'-command.json')).write_text(json.dumps(row,indent=2)+'\n');(O/'sequential-in-progress.json').write_text(json.dumps(rows,indent=2)+'\n');print('END',name,p.returncode,flush=True)
 clean()
 if want_failure:
  data=stdout.read_bytes()+stderr.read_bytes();assert p.returncode!=0 and want_failure.encode() in data and b'[build failed]' not in data and b'SyntaxError' not in data and b'test timed out' not in data,name
 else:assert p.returncode==0,name
 if json_report:assert json.loads(stdout.read_text())['passed'] is True,name
 return row
try:
 run('frozen-root-race',['go','test','-race','-count=1','./...'],R)
 run('frozen-root-vet',['go','vet','./...'],R)
 run('frozen-root-build',['go','build','-trimpath','-o',B/'ensemble-cli','./cmd/ch01'],R)
 run('frozen-gui-race',['go','test','-race','-count=1','./...'],R/'gui',{'CH07_CORRECTION_WIRE':str(RAW/'frozen-full-gui-wire.json')})
 run('frozen-gui-vet',['go','vet','./...'],R/'gui')
 run('frozen-gui-build',['go','build','-trimpath','-o',B/'ensemble-gui','./cmd'],R/'gui')
 for module,label in [('examples/ch04-consumer','consumer'),('examples/ch06','workflow'),('examples/ch07','streaming')]:
  run('frozen-'+label+'-race',['go','test','-race','-count=1','./...'],R/module)
  run('frozen-'+label+'-vet',['go','vet','./...'],R/module)
  run('frozen-'+label+'-build',['go','build','-trimpath','-o',B/label,'.'],R/module)
 focused=['go','test','-race','-count=1','-timeout=30s','-run','^TestCorrectionR1RefusedProposalIsNotDispatch$','-v','.']
 run('frozen-focused',['go',*focused[1:]],R/'gui',{'CH07_CORRECTION_WIRE':str(RAW/'frozen-focused-wire.json')})
 page=['node','--experimental-vm-modules','test/correction-r1-page.mjs',str(RAW/'frozen-focused-wire.json')]
 run('frozen-page-controlled',page,R/'gui')
 run('frozen-existing-queues',['node','test/queues.mjs'],R/'gui')
 copy=R/'.ch07-local/correction-r1/sensitivity-final-copy'; mutations=[]
 for label,file,old,new,cmd,needle in [('r1','gui/hub.go',b'if t.Dispatched || t.Job != nil || h.toolStarts[t.DisplayID] {',b'if true { // deliberate isolated regression',focused,'falsely claims 3 dispatches'),('r2','gui/web/gui/app.js',b'a.toolResult = {text: m.result, error: Boolean(m.is_error)}; a.readable = readableTool(a.toolCall, a.toolResult);',b'a.toolResult = null; a.readable = readableTool(a.toolCall, a.toolResult);',page,'actual tool speaker omitted visible result')]:
  p=copy/file;before=p.read_bytes();assert sha(p)==F['source_files'][file] and before.count(old)==1;bad=before.replace(old,new);p.write_bytes(bad)
  try:run('frozen-control-'+label+'-broken',cmd,copy/'gui',want_failure=needle)
  finally:p.write_bytes(before)
  run('frozen-control-'+label+'-restored',cmd,copy/'gui')
  mutations.append({'case':label,'path':str(p),'original_sha256':hashlib.sha256(before).hexdigest(),'broken_sha256':hashlib.sha256(bad).hexdigest(),'restored_sha256':sha(p),'expected_semantic_failure':True})
 for f,h in F['source_files'].items():assert sha(copy/f)==h,f
 (O/'frozen-control-results.json').write_text(json.dumps(mutations,indent=2)+'\n')
 artifacts={p.name:{'path':str(p),'bytes':p.stat().st_size,'sha256':sha(p)} for p in sorted(B.iterdir()) if p.is_file()}
 assert sha(OLD['grader']['path'])==OLD['grader']['sha256']==F['grader']['sha256']
 (O/'frozen-artifacts-before-graders.json').write_text(json.dumps(artifacts,indent=2)+'\n')
 run('frozen-cli-build-info',['go','version','-m',B/'ensemble-cli'],R)
 run('frozen-gui-build-info',['go','version','-m',B/'ensemble-gui'],R)
 cmd=list(OLD['command']);cmd[-1]=str(B/'ensemble-cli')
 run('chapter7-final',cmd,Path(OLD['cwd']),json_report=True)
 assert sha(B/'ensemble-cli')==artifacts['ensemble-cli']['sha256']
 run('chapter1-final',[OLD['grader']['path'],'-ch','1','-json',B/'ensemble-cli'],Path(OLD['cwd']),json_report=True)
 for name,a in artifacts.items():assert sha(a['path'])==a['sha256'],name
 clean();(O/'sequential-completed.json').write_text(json.dumps({'completed':now(),'candidate':F['commit'],'tree':F['tree'],'all_expected_outcomes_met':True,'commands':rows,'artifacts':artifacts,'same_cli_across_both_graders':True,'grader_sha256':sha(OLD['grader']['path']),'source_clean':True,'no_browser_native_audio_or_paid_run':True},indent=2)+'\n');print('ALL SEQUENTIAL GATES COMPLETE',flush=True)
except BaseException as e:
 (O/'sequential-failure.json').write_text(json.dumps({'at':now(),'error':repr(e),'commands':rows},indent=2)+'\n');raise
