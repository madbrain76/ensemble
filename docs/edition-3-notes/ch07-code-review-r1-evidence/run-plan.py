from pathlib import Path
import json,subprocess,sys
w=Path('/tmp/cadence-ch07-code-r1-rereview-independent')
for label,cwd,cmd in json.loads((w/'gate-plan.json').read_text()):
 subprocess.run([sys.executable,str(w/'run.py'),label,json.dumps(cmd),cwd],check=True)
 receipt=json.loads((w/'logs'/f'{label}.run.json').read_text());assert receipt['exit']==0,(label,receipt['exit'])
print('SEQUENTIAL FINAL GATES COMPLETE')
