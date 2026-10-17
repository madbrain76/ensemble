#!/usr/bin/env python3
"""Finite offline entry controls; writes only scratch copies and this evidence directory."""
import os, pathlib, subprocess, json, time, hashlib, shutil, signal
E = pathlib.Path(__file__).resolve().parent
S = pathlib.Path((E/'scratch.txt').read_text().strip())
M = S/'main'; R = M/'solutions/ch04'; ST = S/'student'
env = dict(os.environ)
for k in list(env):
    if any(x in k.upper() for x in ('API_KEY','TOKEN','SECRET','PASSWORD','CREDENTIAL')) or k.startswith(('ANTHROPIC_','OPENAI_','GOOGLE_','GEMINI_','XAI_','GROQ_','MISTRAL_','DEEPSEEK_','ENSEMBLE_')):
        env.pop(k, None)
env.update(GOPROXY='off', GOTOOLCHAIN='local', GIT_PAGER='cat', PAGER='cat', TERM='dumb')
records=json.loads((E/'runs.json').read_text()) if (E/'runs.json').exists() else []
def run(name, argv, cwd=M, timeout=120):
    if any(r['name']==name for r in records):
        print(name, 'already recorded; not repeated', flush=True); return 0
    start=time.monotonic()
    p=subprocess.Popen([str(x) for x in argv],cwd=cwd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True)
    timed=False
    try: out,_=p.communicate(timeout=timeout)
    except subprocess.TimeoutExpired:
        timed=True; os.killpg(p.pid,signal.SIGKILL); out,_=p.communicate()
    (E/(name+'.log')).write_bytes(out)
    rec=dict(name=name,argv=[str(x) for x in argv],cwd=str(cwd),exit_code=p.returncode,timeout=timed,seconds=round(time.monotonic()-start,3),sha256=hashlib.sha256(out).hexdigest())
    records.append(rec); (E/'runs.json').write_text(json.dumps(records,indent=2)+'\n')
    print(name, p.returncode, rec['seconds'], flush=True)
    return p.returncode
D=M/'cmd/ch04-entry-driver'; D.mkdir(exist_ok=True)
source='''package main
import("fmt";"os";"github.com/waywardgeek/ensemble/internal/gradeedition2")
func main(){ example,err:=gradeedition2.ReadConsumerCase(""); if err!=nil { panic(err) }; report,err:=gradeedition2.Chapter4(os.Args[1],os.Args[2],os.Args[3],example); if err!=nil { fmt.Fprintln(os.Stderr,err); os.Exit(2) }; if err=report.WriteJSON(os.Stdout); err!=nil { panic(err) }; if !report.Passed { os.Exit(1) } }
'''
(D/'main.go').write_text(source); (E/'driver.go.txt').write_text(source)
assert run('build-grade2',['go','build','-o',S/'grade2','./cmd/grade-edition2'])==0
assert run('build-grade3',['go','build','-o',S/'grade3','./cmd/grade-edition3'])==0
assert run('build-driver-corrected',['go','build','-o',S/'driver','./cmd/ch04-entry-driver'])==0
assert run('build-reference',['go','build','-o',S/'reference','.'],R)==0
assert run('build-student',['go','build','-o',S/'student-bin','./cmd/ch01'],ST)==0
run('student-root-build',['go','build','-o',S/'student-root','.'],ST)
run('student-root-file',['file',S/'student-root'])
run('grade3-ch4-rejected',[S/'grade3','-ch','4',S/'reference'])
run('cli-student-root',[S/'grade2','-ch','4','-json',ST])
run('cli-reference-explicit',[S/'grade2','-ch','4','-json','-consumer',R/'ch05',S/'reference'],timeout=180)
run('reference-separated',[S/'driver',S/'reference',R,R/'ch05'],timeout=180)
run('reference-lifetime-corrected',[S/'grade3','-ch','1','-json',S/'reference'],timeout=180)
run('reference-cap-test',['go','test','./ensemble','-run','^TestJobFileCapCursorAndRedaction$','-count=1'],R)
mut=S/'cap-mutant'; shutil.copytree(R,mut,dirs_exist_ok=True)
f=mut/'internal/jobs/jobs.go'; old=f.read_text(); anchor='\tvalue, err := limits(raw, base)'
assert old.count(anchor)==1
new=old.replace(anchor,'\t// ENTRY AUDIT MUTANT: silently ignore pending max_output_bytes.\n\tbase.MaxBytes = defaults().MaxBytes\n'+anchor)
f.write_text(new)
(E/'mutation.json').write_text(json.dumps(dict(path='internal/jobs/jobs.go',before_sha256=hashlib.sha256(old.encode()).hexdigest(),after_sha256=hashlib.sha256(new.encode()).hexdigest(),insertion='base.MaxBytes = defaults().MaxBytes before explicit arguments',effect='ignore only pending output cap; explicit override remains'),indent=2)+'\n')
# Resolve the consumer to its own mutant core; reference consumer has a relative replacement.
assert run('build-cap-mutant',['go','build','-o',S/'cap-mutant-bin','.'],mut)==0
run('cap-mutant-separated',[S/'driver',S/'cap-mutant-bin',mut,mut/'ch05'],timeout=180)
run('cap-mutant-test',['go','test','./ensemble','-run','^TestJobFileCapCursorAndRedaction$','-count=1'],mut)
identities={}
for p in [S/'grade2',S/'grade3',S/'driver',S/'reference',S/'student-bin',S/'student-root',S/'cap-mutant-bin']:
    identities[str(p)]=hashlib.sha256(p.read_bytes()).hexdigest()
(E/'binary-identities.json').write_text(json.dumps(identities,indent=2)+'\n')
