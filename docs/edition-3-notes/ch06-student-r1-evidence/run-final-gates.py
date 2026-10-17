#!/usr/bin/env python3
"""Sequential R1 frozen gates; no provider traffic or credential/config reads."""
import hashlib,json,os,pathlib,subprocess,sys,time
E=pathlib.Path(__file__).resolve().parent
S=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
B=pathlib.Path('/Users/bill/projects/ensemble')
T=pathlib.Path('/tmp/e3-ch06-correction-r1')
M=json.loads((E.parent/'ch06-student-r1-source-manifest.json').read_text())
env={'PATH':os.environ['PATH'],'HOME':str(pathlib.Path.home()),'TMPDIR':str(T),'GOFLAGS':'-mod=readonly','GOPROXY':'off','GOSUMDB':'off','GOTOOLCHAIN':'local'}
def hashfile(p):return {'path':str(p),'bytes':p.stat().st_size,'sha256':hashlib.sha256(p.read_bytes()).hexdigest()}
def pins():
 assert subprocess.check_output(['git','rev-parse','HEAD'],cwd=S,text=True).strip()==M['commit']
 assert subprocess.check_output(['git','rev-parse','HEAD^{tree}'],cwd=S,text=True).strip()==M['tree']
 assert subprocess.check_output(['git','status','--porcelain'],cwd=S)==b''
 for p in M['artifacts']:
  got=hashfile(pathlib.Path(p['path']));assert got['sha256']==p['sha256'] and got['bytes']==p['bytes']
commands=[('root-test',S,['go','test','-count=1','-timeout=240s','./...'],300),('root-race',S,['go','test','-race','-count=1','-timeout=300s','./...'],360),('root-vet',S,['go','vet','./...'],90)]
for name,path in [('consumer','examples/ch04-consumer'),('workflow','examples/ch06'),('shape','examples/ch07')]:
 commands.extend([(name+'-test',S/path,['go','test','-count=1','-timeout=60s','./...'],90),(name+'-vet',S/path,['go','vet','./...'],90)])
commands.extend([('focused-final',S,['go','test','-count=1','-timeout=90s','-run','^TestChapter6|^TestLiveProtocol|^TestDisplay','./...'],120),('shape-fixture',S,[sys.executable,str(E/'shape_fixture.py'),str(T/'shape')],30)])
commands.extend([('grade-ch06',B,['go','run','./cmd/grade-edition3','-ch','6','-json','-core',str(S),'-consumer',str(S/'examples/ch04-consumer'),'-workflow',str(S/'examples/ch06'),'-stream-exercise',str(S/'examples/ch07'),'-stream-model','claude-sonnet-4-6','-media-models',str(E/'media-models.json'),str(T/'core')],300),('grade-ch01',B,['go','run','./cmd/grade-edition3','-ch','1','-json',str(T/'core')],300)])
receipts=[]
(E/'final-planned-commands.json').write_text(json.dumps([{'name':name,'cwd':str(cwd),'argv':argv,'timeout_seconds':timeout} for name,cwd,argv,timeout in commands],indent=2)+'\n')
for name,cwd,argv,timeout in commands:
 pins();out=E/(name+'.stdout.txt');err=E/(name+'.stderr.txt');start=time.monotonic();timed_out=False
 with out.open('xb') as fo,err.open('xb') as fe:
  try:code=subprocess.run(argv,cwd=cwd,env=env,stdout=fo,stderr=fe,timeout=timeout).returncode
  except subprocess.TimeoutExpired:code=124;timed_out=True
 exitpath=E/(name+'.exit.txt');exitpath.write_text(str(code)+'\n')
 entry={'name':name,'cwd':str(cwd),'argv':argv,'environment':env,'timeout_seconds':timeout,'timed_out':timed_out,'exit_code':code,'elapsed_seconds':time.monotonic()-start,'stdout':hashfile(out),'stderr':hashfile(err),'exit_receipt':hashfile(exitpath),'source_commit':M['commit'],'source_tree':M['tree'],'core':M['artifacts'][1]}
 if name.startswith('grade-') and code==0:
  grade=json.loads(out.read_text());entry['passed']=grade.get('passed');entry['score']=grade.get('score')
 receipts.append(entry);(E/'final-command-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
 print(name,'exit',code,'seconds',round(entry['elapsed_seconds'],3),'passed',entry.get('passed'),flush=True)
 if code!=0 or entry.get('passed') is False:raise SystemExit(code or 1)
pins()
(E/'final-reconciliation.json').write_text(json.dumps({'commit':M['commit'],'tree':M['tree'],'clean':True,'all_artifact_pins_match':True,'completed_commands':len(receipts),'all_exits_zero':all(x['exit_code']==0 for x in receipts),'grades_passed':[x['passed'] for x in receipts if x['name'].startswith('grade-')],'provider_calls':0,'credential_reads':0},indent=2)+'\n')
print('ALL FROZEN GATES COMPLETED',flush=True)
