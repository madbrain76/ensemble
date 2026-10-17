#!/usr/bin/env python3
"""Bounded Chapter4 launcher: secrets never leave process environment or memory."""
import hashlib, json, os, pathlib, subprocess, sys, time, socket
provider, mode = sys.argv[1:3]
models = {'anthropic': ('claude-sonnet-4-6','directClaudeAPIKey','https://api.anthropic.com'), 'openai': ('gpt-4.1-2025-04-14','directOpenAIAPIKey','https://api.openai.com'), 'gemini': ('gemini-3.8-flash','directGeminiAPIKey','https://generativelanguage.googleapis.com')}
model, field, endpoint = models[provider]
base = pathlib.Path('/tmp/ch04-live') / (provider+'-'+mode)
workspace, data = base/'workspace', base/'data'
workspace.mkdir(parents=True,exist_ok=True); data.mkdir(exist_ok=True)
evidence = pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes/ch04-student-evidence')
settings = json.loads((pathlib.Path.home()/'.cr/settings.json').read_text())
key = settings.get(field)
if not isinstance(key,str) or not key: raise SystemExit('authorized credential field unavailable: '+field)
secret_values = [v for k,v in settings.items() if 'key' in k.lower() and isinstance(v,str) and v]
del settings
binary = pathlib.Path('/tmp/ch04-core' if mode.startswith('core') else '/tmp/ch04-consumer')
env = os.environ.copy()
for name in ('LLM_API_KEY','ANTHROPIC_API_KEY','OPENAI_API_KEY','GEMINI_API_KEY'): env.pop(name,None)
env.update(LLM_VENDOR=provider,LLM_MODEL=model,LLM_BASE_URL=endpoint,LLM_API_KEY=key,LLM_REQUIRE_TOOL_EXPLANATION='true',ENSEMBLE_WORKSPACE=str(workspace),ENSEMBLE_DATA_DIR=str(data),CH02_LOG=str(base/'events.jsonl'))
if mode.startswith('core'):
    (workspace/'greeting.txt').write_text('helo\n')
    (workspace/'main.go').write_text('package main\nimport "fmt"\nfunc main() {\n value := 41\n fmt.Println(value + 1)\n}\n')
    if provider=='anthropic':
        subprocess.run(['go','build','-gcflags=all=-N -l','-o',str(workspace/'debug-target'),str(workspace/'main.go')],check=True)
        prompts = [
          'This disposable Chapter4 exercise has session token MAPLE. Explain visibly before EVERY tool batch. Repair greeting.txt from helo to hello using read/edit tools, then verify. Exercise one-shot tool_limits on a large file read: create large.txt with HEAD-MARK + 24000 x bytes + TAIL-MARK using python3 through run_command, set max_output_bytes 1024, read it and inspect the persisted output if needed. Start a PTY command `printf "READY>"; read line; printf "answer=%s\\n" "$line"` with ai_callback_pattern READY and ai_callback_delay 1; submit hello using send_input, then wait once on the finished handle and note no output replay. Finally run `/Users/bill/go/bin/dlv exec ./debug-target` with ai_callback_pattern "\\(dlv\\)" and delay 1; use send_input to break main.main, continue, next, print value, then quit (send y if asked). Do not claim success without tool output. Finish with a concise evidence summary.',
          'Continue the same session: name the earlier session token and verify greeting.txt is still hello. Last, start a retained job with command `echo EOF-SURVIVOR; sleep 60`, ai_callback_pattern EOF-SURVIVOR and ai_callback_delay 0.1. Once the tool reports it running, do NOT wait for or kill that job. State its handle and finish your answer so CLI EOF triggers root shutdown. No more tools after observing the running job.'
        ]
    else:
        prompts = [
          'Disposable Chapter4 exercise, session token MAPLE. Explain visibly before EVERY tool batch. Read greeting.txt, repair helo to hello using edit_file, then verify. Start `echo KEEP-RUNNING; sleep 60` using run_command with ai_callback_pattern KEEP-RUNNING and ai_callback_delay 0.1; record its running handle, kill_job that handle, then wait_for_job to confirm killed without replay. Finish concisely with observed evidence.',
          'Continue the same session: name the earlier session token and read greeting.txt to verify the repair persisted. Set tool_limits max_output_bytes 512 then run `printf "continuity-ok\\n"`; identify which call consumed the limits. Finish with concise observed facts.'
        ]
    if mode=='core-followup':
        sock=socket.socket();sock.bind(('127.0.0.1',0));port=sock.getsockname()[1];sock.close()
        (workspace/'go.mod').write_text('module ch04live\n\ngo 1.22\n')
        (workspace/'main_test.go').write_text('package main\nimport("os";"strings";"testing";"time")\nfunc TestGreeting(t *testing.T){time.Sleep(8*time.Second);b,e:=os.ReadFile("greeting.txt");if e!=nil{t.Fatal(e)};if strings.TrimSpace(string(b))!="hello"{t.Fatalf("want hello, got %q",b)}}\n')
        prompts=[
            f'Session token MAPLE. This is a small disposable Chapter4 test. Before EVERY tool batch emit visible explanation prose. Start the actual test command `go test -v .` with ai_callback_delay 0.1; its TestGreeting deliberately takes 8 seconds and may fail because greeting.txt contains helo. Retain the running handle; do not wait yet. Next start a real server using `python3 -m http.server {port} --bind 127.0.0.1` with ai_callback_pattern "Serving HTTP" and ai_callback_delay 1. Keep the server running across all turns. Use a separate run_command to curl http://127.0.0.1:{port}/greeting.txt. Conclude this turn with test/server handles and actual curl evidence; no extra tasks.',
            f'Continue session. Name our earlier token. Use wait_for_job on the earlier test handle and inspect its actual failure. Read greeting.txt, edit helo to hello, then start `go test -v .` again with ai_callback_delay 0.1 so it returns a running handle. Curl http://127.0.0.1:{port}/greeting.txt to prove the STILL-RUNNING server serves the repaired file. Do not wait for the new test yet. Conclude with observed facts and new test handle. Visible explanation before EACH batch, including follow-ups.',
            'Continue. First wait_for_job on the repaired test and verify actual PASS (wait longer if needed). Next start the real debugger `/Users/bill/go/bin/dlv exec ./debug-target` with ai_callback_pattern "\\(dlv\\)" and ai_callback_delay 1. Then send_input with these newline-separated commands in ONE input: break main.main; continue; next; print value; quit (use NEWLINES not semicolons). This is the provided tiny Go target, value should be 41 after next. Use wait_for_job or send y only if needed. Report actual breakpoint, stop, printed value, and debugger exit; never model assertions alone. Keep HTTP server running. Visible explanation before EACH tool batch.',
            f'Final turn: name our session token and curl http://127.0.0.1:{port}/greeting.txt once more. Then use wait_for_job on the ORIGINAL HTTP server handle with ai_callback_delay 0.01 to confirm it is still running. Do NOT kill that server, and do NOT make any more tool calls after that running observation. Finish your answer naming its running handle, so CLI input EOF triggers root shutdown. Visible explanation before each batch.'
        ]
    args = [str(binary),'chat']; stdin='\n'.join(prompts)+'\n' 
else: args=[str(binary)]; stdin=None
start=time.time(); observed={}
proc=subprocess.Popen(args,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True,env=env,cwd=workspace)
proc.stdin.write(stdin);proc.stdin.close();proc.stdin=None
while proc.poll() is None:
    for line in subprocess.check_output(['ps','-axo','pid,ppid,command'],text=True).splitlines()[1:]:
        parts=line.strip().split(None,2)
        if len(parts)==3 and f'-m http.server {port} ' in parts[2]:
            pid=int(parts[0]);item=observed.setdefault(pid,{'pid':pid,'ppid':int(parts[1]),'command':parts[2],'first_alive':time.time()});item['last_alive']=time.time()
    time.sleep(.1)
stdout,stderr=proc.communicate();result=subprocess.CompletedProcess(args,proc.returncode,stdout,stderr)
for pid,item in observed.items():
    try:os.kill(pid,0);item['absent_after_cli_exit']=False
    except ProcessLookupError:item['absent_after_cli_exit']=True
    except PermissionError:item['absent_after_cli_exit']='permission denied'
(evidence/'anthropic-core-followup.processes.json').write_text(json.dumps({'server_port':port,'cli_exit_at':time.time(),'processes':list(observed.values())},indent=2)+'\n')
def safe(text):
    for secret in secret_values: text=text.replace(secret,'[REDACTED CREDENTIAL]')
    return text
prefix=provider+'-'+mode
(evidence/(prefix+'.stdout.txt')).write_text(safe(result.stdout))
(evidence/(prefix+'.stderr.txt')).write_text(safe(result.stderr))
def sanitize(value):
    if isinstance(value,dict): return {k:('[opaque omitted]' if k=='opaque' else sanitize(v)) for k,v in value.items()}
    if isinstance(value,list): return [sanitize(v) for v in value]
    if isinstance(value,str): return safe(value)
    return value
log=base/'events.jsonl'; events=[]
if log.exists():
    for line in log.read_text().splitlines():
        if line.strip(): events.append(sanitize(json.loads(line)))
    # Do not retain a copied log containing signed opaque material.
    log.write_text(''.join(json.dumps(e)+'\n' for e in events))
(evidence/(prefix+'.events.jsonl')).write_text(''.join(json.dumps(e)+'\n' for e in events))
metadata={'provider':provider,'requested_model':model,'endpoint':endpoint,'surface':{'anthropic':'Messages','openai':'ChatCompletions','gemini':'GenerateContent'}[provider],'mode':mode,'binary':str(binary),'binary_sha256':hashlib.sha256(binary.read_bytes()).hexdigest(),'exit_code':result.returncode,'elapsed_seconds':round(time.time()-start,3),'workspace':str(workspace),'data_dir':str(data),'prompts':stdin,'events':len(events),'source_tree':pathlib.Path('/tmp/ch04-final-source-tree.txt').read_text().strip(),'git_head':subprocess.check_output(['git','rev-parse','HEAD'],cwd='/Users/bill/projects/ensemble-edition-3-student',text=True).strip()}
(evidence/(prefix+'.metadata.json')).write_text(json.dumps(metadata,indent=2)+'\n')
print(json.dumps({k:metadata[k] for k in ('provider','requested_model','mode','binary_sha256','exit_code','elapsed_seconds','events')}))
print(safe(result.stdout));print(safe(result.stderr),file=sys.stderr)
