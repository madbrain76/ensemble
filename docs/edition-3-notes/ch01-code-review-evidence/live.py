"""Bounded reviewer live chat. Never persist credentials or raw provider errors."""
from pathlib import Path
import errno, json, os, pty, select, subprocess, time, urllib.request
E = Path(__file__).resolve().parent
key = json.loads((Path.home() / '.cr/settings.json').read_text()).get('directClaudeAPIKey')
if not isinstance(key, str) or not key:
 raise SystemExit('Authorized directClaudeAPIKey missing; no live run.')
log = (E / 'live.txt').open('w')
def record(s):
 log.write(s.replace(key, '[REDACTED]'))
 log.flush()
record('Independent reviewer live PTY run; exact d2baff64055ac33702a91c6d6ee5a4246e84c0c9 binary.\n')
try:
 request = urllib.request.Request('https://api.anthropic.com/v1/models', headers={'x-api-key':key, 'anthropic-version':'2023-06-01'})
 with urllib.request.urlopen(request, timeout=30) as response:
  models = [item['id'] for item in json.load(response)['data']]
 model = 'claude-sonnet-5' if 'claude-sonnet-5' in models else next(m for m in models if 'sonnet' in m)
 record(f'GET /v1/models succeeded; discovered {len(models)} models; selected {model}.\n')
except Exception as error:
 record(f'Model discovery failed: {type(error).__name__}; no chat attempted.\n')
 raise SystemExit(1)
env = os.environ.copy()
env.update(ANTHROPIC_API_KEY=key, ANTHROPIC_BASE_URL='https://api.anthropic.com', ANTHROPIC_MODEL=model)
master, slave = pty.openpty()
binary = str(Path((E/'scratch-path.txt').read_text().strip())/'ch01')
child = subprocess.Popen([binary, 'chat'], stdin=slave, stdout=slave, stderr=slave, env=env, close_fds=True)
os.close(slave)
def until_prompt():
 data = b''
 # Reviewer safety watchdog, not an application generation deadline.
 end = time.monotonic() + 300
 while not data.endswith(b'> '):
  if time.monotonic() >= end:
   raise TimeoutError('reviewer watchdog')
  if select.select([master], [], [], 1)[0]:
   try:
    chunk = os.read(master, 65536)
   except OSError as e:
    if e.errno == errno.EIO: chunk = b''
    else: raise
   if not chunk:
    record(data.decode(errors='replace'))
    raise RuntimeError('chat ended before prompt')
   data += chunk
 record(data.decode(errors='replace'))
try:
 until_prompt()
 for prompt in ['My canoe is named Cobalt Wren. Remember that during this conversation.', 'What is my canoe named?', 'In one sentence, is that name stored outside this conversation or supplied in its messages?']:
  os.write(master, (prompt+'\n').encode())
  until_prompt()
 os.write(master, b'\x04')
 result = child.wait(timeout=10)
 record(f'\nCtrl-D sent; chat exit={result}.\n')
 if result != 0: raise RuntimeError('nonzero chat exit')
 print('Independent live PTY chat completed; sanitized evidence saved.')
except Exception as error:
 record(f'\nReviewer run failed: {type(error).__name__}.\n')
 child.kill(); child.wait()
 raise SystemExit(1)
finally:
 os.close(master)
 log.close()
