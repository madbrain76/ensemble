#!/usr/bin/env python3
"""Read-only Python audit observer; runs the unchanged controller as __main__."""
import json,os,pathlib,runpy,sys,time
output=pathlib.Path(sys.argv[1]);controller=pathlib.Path(sys.argv[2]);args=sys.argv[2:]
secret_path=str(pathlib.Path.home()/'.cr/settings.json');events=[];start=time.time();credential_opens=0

def observe(event,args):
 global credential_opens
 if event=='open' and isinstance(args[0],(str,bytes,os.PathLike)) and os.fsdecode(args[0])==secret_path:
  credential_opens+=1;events.append(dict(event='authorized_credential_file_open',at_unix=time.time()))
 elif event=='socket.connect':events.append(dict(event=event,address=args[1],at_unix=time.time()))
 elif event=='socket.getaddrinfo':events.append(dict(event=event,host=args[0],port=args[1],at_unix=time.time()))
 elif event=='subprocess.Popen':events.append(dict(event=event,executable=os.fsdecode(args[0]),at_unix=time.time()))

sys.addaudithook(observe);sys.argv=args;sys.path.insert(0,str(controller.resolve().parent))
try:runpy.run_path(str(controller),run_name='__main__')
finally:
 output.write_text(json.dumps(dict(started_unix=start,finished_unix=time.time(),credential_file_open_count=credential_opens,credential_values_observed=False,scope='Python controller process only; does not observe Chrome background effects or native child internals',events=events),indent=2)+'\n')
