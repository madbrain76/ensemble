import hashlib,json,os,pathlib,shutil,subprocess,time
root=pathlib.Path('/Users/bill/projects/ensemble-edition-3-student');out=pathlib.Path('/tmp/cadence-ch05-wakeup-evidence/sensitivity-r2');out.mkdir();scratch=out/'source';scratch.mkdir();base='6c4561cbd6cbae080d952da34a5fcd588a76920f'
archive=out/'original-source.tar';archive.write_bytes(subprocess.check_output(['git','archive',base],cwd=root));subprocess.run(['tar','xf',str(archive),'-C',str(scratch)],check=True)
files=subprocess.check_output(['git','ls-tree','-r','--name-only',base],cwd=root,text=True).splitlines();sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
expected={n:hashlib.sha256(subprocess.check_output(['git','show',base+':'+n],cwd=root)).hexdigest() for n in files}
assert all(sha(scratch/n)==h for n,h in expected.items())
shutil.copyfile(root/'chapter5_wakeup_test.go',scratch/'chapter5_wakeup_test.go');env={**os.environ,'GOWORK':'off','GIT_PAGER':'cat'};rows=[]
def run(label,pattern,location,want):
 if shutil.disk_usage('/tmp').free<1024**3:raise SystemExit('disk below1GiB')
 args=['go','test','-race','-count=1','-timeout','90s','-run',pattern,'-v','.'];start=time.monotonic();p=subprocess.run(args,cwd=location,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);log=p.stdout.decode();(out/(label+'.log')).write_text(log)
 valid=(p.returncode==0) if want=='pass' else (p.returncode==1 and '--- FAIL: TestChapter5Wakeup' in log and 'WAKEUP' in log and not any(x in log for x in ['[build failed]','no space left','panic: test timed out','DATA RACE']))
 row={'label':label,'command':args,'cwd':str(location),'exit':p.returncode,'elapsed_seconds':round(time.monotonic()-start,3),'expected':want,'valid':valid};rows.append(row);(out/'commands.json').write_text(json.dumps(rows,indent=2)+'\n');print(label,p.returncode,valid,flush=True);assert valid,log
for name,pattern in [('initial','InitialMixedBatch'),('setter','SetterDispatchBoundary'),('later','LaterProcessWait'),('available','TimerAndAvailableOutput'),('timer','ShortTimerCannotReleaseLaterWait'),('failed-report','FailedRetainedReport')]:run('original-'+name,'TestChapter5Wakeup'+pattern,scratch,'behavioral-fail')
changed=subprocess.check_output(['git','diff','--name-only'],cwd=root,text=True).splitlines()
for n in changed:shutil.copyfile(root/n,scratch/n)
run('restored-corrected','TestChapter5Wakeup|TestChapter5Correction',scratch,'pass')
(out/'source-controls.json').write_text(json.dumps({'original_commit':base,'original_sha256':expected,'test_sha256':sha(root/'chapter5_wakeup_test.go'),'restored_files':{n:sha(root/n) for n in changed}},indent=2)+'\n')
