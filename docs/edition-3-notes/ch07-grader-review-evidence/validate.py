from run import W,ENV,run
import json,subprocess,time,os
C='/tmp/cadence-ch07-grader';M=str(W/'module');T=str(W/'bin/sentinel')
(W/'observe-bin').mkdir();go=W/'observe-bin/go';go.write_text('#!/bin/sh\nprintf BUILT >> "$REVIEW_BUILD_MARKER"\nexit 97\n');go.chmod(0o700)
realmarker=W/'sentinel-positive.marker';m,*_=run('sentinel-positive',[T],extra={'REVIEW_SENTINEL_MARKER':str(realmarker)});assert m['exit']==99 and realmarker.read_text()=='EXECUTED\n'
base={'-ch':'7','-core':M,'-consumer':M,'-workflow':M,'-stream-exercise':M,'-stream-model':'declared-stream','-gui':M,'-gui-model':'declared-gui'}
def args(d): return [x for pair in d.items() for x in pair]
cases=[]
def add(n,d=None,a=None,code=2):cases.append((n,args(d) if d is not None else a,code))
for flag in list(base)[1:]:
 d=base.copy();del d[flag];add('missing-'+flag[1:],d)
 d=base.copy();d[flag]='';add('empty-'+flag[1:],d)
for flag in ['-core','-consumer','-workflow','-stream-exercise','-gui']:
 for label,path in [('absent',W/'absent'),('nonmodule',W/'nonmodule'),('nonregular',W/'go-mod-directory')]:
  d=base.copy();d[flag]=str(path);add(label+'-'+flag[1:],d)
for flag in ['-stream-model','-gui-model']:
 d=base.copy();d[flag]=' \t\n';add('blank-'+flag[1:],d)
for flag in ['-consumer-case','-media-models']:
 paths=[('missing',W/'absent.json'),('malformed',W/'bad.json'),('schema',W/('bad-case.json' if flag=='-consumer-case' else 'bad-media.json')),('directory',W/'module')]
 for label,path in paths:
  d=base.copy();d[flag]=str(path);add(label+'-'+flag[1:],d)
for ch in ['1','4','5','6']:
 d={'-ch':ch}
 if ch!='1': d.update({'-core':M,'-consumer':M})
 if ch in ['5','6']:d['-workflow']=M
 if ch=='6':d.update({'-stream-exercise':M,'-stream-model':'declared-stream'})
 for flag in ['-gui','-gui-model']:
  for label,value in [('empty',''),('value','not-legal-here')]:
   n=d.copy();n[flag]=value;add('old'+ch+'-'+label+'-'+flag[1:],n)
for ch in ['0','2','3','8','-1']:add('unsupported-'+ch,{'-ch':ch})
add('unknown-flag',a=['-unknown-adapter-flag']);add('help',a=['-h'],code=0);add('long-help',a=['-help'],code=0);add('ch7-help',a=['-ch','7','-h'],code=0)
rows=[]
for i,(name,flags,code) in enumerate(cases):
 for kind,target in [('compiled',T),('build-observer',M)]:
  marker=W/f'gate-{i}-{kind}.marker';buildmarker=W/f'build-{i}-{kind}.marker'
  env=ENV.copy();env.update(REVIEW_SENTINEL_MARKER=str(marker),REVIEW_BUILD_MARKER=str(buildmarker),PATH=str(W/'observe-bin')+os.pathsep+ENV['PATH'])
  argv=[C,*flags,target];start=time.time();p=subprocess.run(argv,cwd=str(W),env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
  row=dict(name=name,kind=kind,argv=argv,cwd=str(W),exit=p.returncode,expected_exit=code,seconds=time.time()-start,stdout=p.stdout.decode(),stderr=p.stderr.decode(),target_launched=marker.exists(),build_invoked=buildmarker.exists(),score=None,passed=None)
  rows.append(row)
  (W/'validation-matrix.json').write_text(json.dumps(rows,indent=2)+'\n')
  assert p.returncode==code and not marker.exists() and not buildmarker.exists(),row
print(json.dumps({'cases':len(cases),'invocations':len(rows),'all_expected_exits':True,'target_launches':0,'build_invocations':0}))
