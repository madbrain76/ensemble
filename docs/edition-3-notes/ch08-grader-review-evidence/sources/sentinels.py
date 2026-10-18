import sys;sys.dont_write_bytecode=True
from common import *
D=T/'sentinel';D.mkdir();(D/'shim').mkdir();(D/'marks').mkdir();(D/'target').mkdir();(D/'target/go.mod').write_text('module sentineltarget\ngo 1.25\n');(D/'target/main.go').write_text('package main\nfunc main() {}\n');(D/'nonregular').mkdir();(D/'nonregular/go.mod').mkdir();(D/'bad.json').write_text('{');src='''package main
import("os";"path/filepath";"encoding/json")
func main(){kind:="target";code:=98;if filepath.Base(os.Args[0])=="go"{kind="go";code=97};data,_:=json.Marshal(map[string]any{"argv":os.Args,"kind":kind});if err:=os.WriteFile(os.Getenv("SENTINEL_MARKER")+"."+kind,data,0600);err!=nil{panic(err)};os.Exit(code)}
''';(D/'marker.go').write_text(src);(E/'sources/sentinel-marker.go.txt').write_text(src)
assert run('build-compiled-sentinel',['go','build','-o',str(D/'marker'),str(D/'marker.go')],T/'work')==0
(D/'shim/go').symlink_to(D/'marker');baseenv={'PATH':str(D/'shim')+':/usr/local/go/bin:/usr/bin:/bin'}
for name,argv,kind,exitcode in [('sentinel-target-selftest',[str(D/'marker')],'target',98),('sentinel-build-selftest',[str(D/'shim/go'),'build','-o',str(D/'unused'),'.'],'go',97)]:
 prefix=D/'marks'/name;rc=run(name,argv,T/'work',20,dict(baseenv,SENTINEL_MARKER=str(prefix)));assert rc==exitcode and Path(str(prefix)+'.'+kind).is_file();shutil.copyfile(str(prefix)+'.'+kind,E/(name+'.marker.json'))
valid=['-core',str(R/'solutions/ch08'),'-consumer',str(R/'solutions/ch08/ch05'),'-workflow',str(R/'solutions/ch08/ch06'),'-stream-exercise',str(R/'solutions/ch08/ch07'),'-stream-model','claude-sonnet-4-6','-gui',str(R/'solutions/ch08/gui'),'-gui-model','claude-sonnet-4-6','-consumer-case',str(T/'consumer-case.json'),'-media-models',str(T/'media-models.json')]
def alter(flag,value=None,omit=False):
 a=list(valid);i=a.index(flag)
 if omit:del a[i:i+2]
 else:a[i+1]=value
 return a
cases=[]
def case(name,args,expected=2,target='dir'):
 end=[] if target=='none' else [str(D/'marker' if target=='bin' else D/'target')];cases.append((name,args+end,expected))
for ch in [8,7,1]:case('help-ch'+str(ch),['-ch',str(ch),'-h'],0,'bin')
for f in ['core','consumer','workflow','stream-exercise','gui']:case('ch8-missing-'+f,['-ch','8']+alter('-'+f,omit=True),target='bin' if f=='core' else 'dir')
for f in ['stream-model','gui-model']:
 for label,value in [('empty',''),('blank',' \t ')]:case('ch8-'+label+'-'+f,['-ch','8']+alter('-'+f,value),target='bin' if label=='empty' else 'dir')
case('ch8-nonregular-go-mod',['-ch','8']+alter('-gui',str(D/'nonregular')))
for f in ['consumer-case','media-models']:case('ch8-invalid-'+f,['-ch','8']+alter('-'+f,str(D/'bad.json')),target='bin' if f=='media-models' else 'dir')
for ch in [2,9]:case('unsupported-ch'+str(ch),['-ch',str(ch)])
case('ch8-missing-target',['-ch','8']+valid,target='none');case('ch8-extra-target',['-ch','8']+valid+[str(D/'marker')])
case('ch7-missing-gui',['-ch','7']+alter('-gui',omit=True),target='bin');case('ch7-empty-gui-model',['-ch','7']+alter('-gui-model',''),target='bin');case('ch7-blank-stream-model',['-ch','7']+alter('-stream-model',' \t '),target='dir')
for ch,flags in [(1,['core','workflow','stream-model','gui']),(4,['workflow','stream-model','gui']),(5,['stream-model','gui']),(6,['gui'])]:
 for f in flags:case('ch'+str(ch)+'-reject-empty-'+f,['-ch',str(ch),'-'+f,''],target='bin')
(E/'sentinel-plan.json').write_text(json.dumps([{'name':n,'argv':a,'expected_exit':v} for n,a,v in cases],indent=2)+'\n')
binary=json.loads((E/'pins-before.json').read_text())['candidate']['binary']['path'];results=[]
for name,args,expected in cases:
 prefix=D/'marks'/name;rc=run(name,[binary]+args,T/'work',20,dict(baseenv,SENTINEL_MARKER=str(prefix)));markers=[str(p) for p in prefix.parent.glob(prefix.name+'.*')];result={'name':name,'exit':rc,'expected_exit':expected,'markers':markers,'passed':rc==expected and not markers};results.append(result);(E/'sentinel-results.json').write_text(json.dumps(results,indent=2)+'\n');assert result['passed'],result
print(json.dumps({'cases':len(results),'all_expected_exits':all(x['passed'] for x in results),'target_and_build_selftests':True,'compiled_marker_sha256':sha(D/'marker'),'free_bytes':shutil.disk_usage(R).free},indent=2))
