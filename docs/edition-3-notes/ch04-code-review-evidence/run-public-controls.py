#!/usr/bin/env python3
import hashlib,json,os,pathlib,re,shutil,signal,subprocess,time
E=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch04-code-review-evidence');S=pathlib.Path((E/'scratch.txt').read_text().strip());C=S/'candidate';P=S/'public-consumer-r2';P.mkdir();shutil.copyfile(E/'public_api_test.go.txt',P/'public_api_test.go');shutil.copyfile(C/'go.sum',P/'go.sum')
(P/'go.mod').write_text('module ch04-review-consumer\n\ngo 1.25.0\n\nrequire ensemble v0.0.0\nrequire github.com/creack/pty v1.1.24 // indirect\n\nreplace ensemble => '+str(C)+'\n')
env={k:v for k,v in os.environ.items() if not re.search(r'API_KEY|TOKEN|SECRET|PASSWORD|CREDENTIAL|^(ANTHROPIC|OPENAI|GEMINI|GOOGLE|LLM|ENSEMBLE|AWS|AZURE)',k,re.I)};env.update(GOPROXY='off',GOTOOLCHAIN='local',GOFLAGS='-mod=readonly')
records=[]
for name,pat,race in [('positive-public-controls','TestPublic(PendingResetOverrideAndIsolation|PendingRearmExplicitDefaultControl|PluginLifecycleCursorAndTerminalHistory|PTYNormalInputKillControl|CwdEnvironmentAndUnseenPattern)$',True),('pending-rearm-reproduction','^TestPublicPendingRearm$',False),('pending-explicit-reset-control','^TestPublicPendingRearmExplicitDefaultControl$',False),('pty-blocked-input-reproduction','^TestPublicPTYBlockedInputKill$',True),('pty-small-input-restored-control','^TestPublicPTYNormalInputKillControl$',True)]:
 args=['go','test','-count=1','-v','-timeout','20s','-run',pat];args+=['-race'] if race else [];args+=['.'];start=time.monotonic();p=subprocess.Popen(args,cwd=P,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True);expired=False
 try:out=p.communicate(timeout=30)[0]
 except subprocess.TimeoutExpired:expired=True;os.killpg(p.pid,signal.SIGKILL);out=p.communicate()[0]
 (E/(name+'.log')).write_bytes(out);records.append({'name':name,'argv':args,'cwd':str(P),'exit_code':p.returncode,'timeout':expired,'seconds':round(time.monotonic()-start,3),'output_sha256':hashlib.sha256(out).hexdigest()});(E/'public-api-runs.json').write_text(json.dumps(records,indent=2)+'\n');print(name,p.returncode,records[-1]['seconds'],flush=True)
print('COMPLETE',flush=True)
