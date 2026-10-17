from pathlib import Path
import subprocess,json,time,hashlib
E=Path(__file__).resolve().parent;S=Path((E/'scratch-path.txt').read_text().strip());C='/tmp/cadence-ch06-grader';T=str(S/'sentinel');M=str(S/'module');marker=S/'launched'
assert not marker.exists()
p=subprocess.run([T],capture_output=True);assert p.returncode==99 and marker.exists();marker.rename(S/'sentinel-selfcheck-marker')
results=[{'name':'sentinel self-check','argv':[T],'exit':p.returncode,'stdout':p.stdout.decode(),'stderr':p.stderr.decode(),'marker_written':True}]
base=['-ch','6','-core',M,'-consumer',M,'-workflow',M,'-stream-exercise',M,'-stream-model','claude-sonnet-4-6']
cases=[]
for flag in ['-core','-consumer','-workflow','-stream-exercise','-stream-model']:
 pairs=[base[i:i+2] for i in range(0,len(base),2)];without=sum((x for x in pairs if x[0]!=flag),[])
 cases += [(f'missing {flag}',without+[T],2),(f'empty {flag}',base+[flag+'=',T],2)]
for flag in ['-core','-consumer','-workflow','-stream-exercise']:
 for kind,root in [('missing',S/'missing'),('no-go.mod',S),('nonregular-go.mod',S/'nonregular'),('file-not-root',S/'sentinel')]:
  cases.append((kind+' '+flag,base+[flag,str(root),T],2))
cases.append(('whitespace-only model',base+['-stream-model',' \t\n',T],2))
for flag in ['-consumer-case','-media-models']:
 for kind,path in [('malformed',S/'bad-json.json'),('missing',S/'missing.json'),('directory',S/'module'),('schema',S/('bad-case.json' if flag=='-consumer-case' else 'bad-models.json'))]:
  cases.append((kind+' '+flag,base+[flag,str(path),T],2))
for ch in ['1','4','5']:
 old=['-ch',ch]
 if ch!='1':old+=['-core',M,'-consumer',M]
 if ch=='5':old+=['-workflow',M]
 for flag in ['-stream-exercise','-stream-model']:
  for v in ['', 'otherwise-valid-value']:cases.append((f'ch{ch} {flag} {v!r}',old+[flag+'='+v,T],2))
 cases.append((f'ch{ch} help',old+['-h',T],0))
cases += [('ch6 help',base+['-h',T],0),('default help',['--help',T],0),('target absent',base,2),('two targets',base+[T,T],2),('post-positional flag',base+[T,'-stream-model','claude-sonnet-4-6'],2),('flags entirely post-positional',[T]+base,2),('unknown flag',base+['-no-such-flag',T],2),('unsupported chapter',['-ch','7',T],2),('bad chapter',['-ch','NaN',T],2)]
for flag in ['-consumer-case','-media-models']:
 cases.append(('JSON validation before invalid build target '+flag,base+[flag,str(S/'bad-json.json'),str(S/'nonexistent-target')],2))
for name,args,want in cases:
 started=time.monotonic();p=subprocess.run([C]+args,capture_output=True)
 item={'name':name,'argv':[C]+args,'exit':p.returncode,'expected_exit':want,'seconds':time.monotonic()-started,'stdout':p.stdout.decode(),'stderr':p.stderr.decode(),'marker_absent':not marker.exists()}
 item['stdout_sha256']=hashlib.sha256(p.stdout).hexdigest();item['stderr_sha256']=hashlib.sha256(p.stderr).hexdigest();results.append(item)
 (E/'runtime-gates.json').write_text(json.dumps(results,indent=2)+'\n')
 assert p.returncode==want,(name,p.returncode,p.stderr)
 assert not marker.exists(),name
 assert not p.stdout,name
 if name.startswith('JSON validation before'):assert 'nonexistent-target' not in p.stderr.decode() and ('consumer case:' in p.stderr.decode() or 'media models:' in p.stderr.decode())
print(json.dumps({'candidate_sha256':hashlib.sha256(Path(C).read_bytes()).hexdigest(),'sentinel_sha256':hashlib.sha256(Path(T).read_bytes()).hexdigest(),'gates':len(cases),'all_expected_exits':True,'marker_absent_after_all_gates':not marker.exists()},indent=2))
