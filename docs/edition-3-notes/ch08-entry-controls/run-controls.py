import subprocess,os,signal,json,time,pathlib
T=pathlib.Path('/tmp/ch08-entry-audit.aTPJZB'); records=[]
def run(name,args,timeout=30):
 start=time.monotonic()
 with (T/(name+'.log')).open('wb') as log:
  p=subprocess.Popen([str(T/'clean-env'),*map(str,args)],cwd=T,stdin=subprocess.DEVNULL,stdout=log,stderr=subprocess.STDOUT,start_new_session=True)
  timed=False
  try: rc=p.wait(timeout=timeout)
  except subprocess.TimeoutExpired:
   timed=True; os.killpg(p.pid,signal.SIGKILL);rc=p.wait()
  alive=True
  try: os.killpg(p.pid,0)
  except ProcessLookupError:alive=False
  if alive:
   os.killpg(p.pid,signal.SIGKILL)
  record={'name':name,'argv':list(map(str,args)),'pid_pgid':p.pid,'exit':rc,'timeout':timed,'seconds':round(time.monotonic()-start,3),'remaining_process_group_after_exit':alive}
  records.append(record); print(json.dumps(record),flush=True)
run('publication-help',[T/'grade-edition3','-h',T/'sentinel'])
run('publication-ch08',[T/'grade-edition3','-ch','8',T/'sentinel'])
for mode in ['initial-budget-omitted','no-persistence']:
 run('mutant-'+mode,[T/'settings-driver',T/mode])
(T/'control-results.json').write_text(json.dumps({'controls':records,'sentinel_marker_exists':(T/'SENTINEL-LAUNCHED').exists()},indent=2))
