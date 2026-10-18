#!/usr/bin/env python3
"""Owned real GUI/stdin/WebSocket/restart test; external provider is local and fake.
Usage: settings_process.py /absolute/ensemble-gui /absolute/owned-run-directory
No browser, native API, credentials, or external provider activity.
"""
import base64, hashlib, http.server, json, os, pathlib, re, socket, struct
import subprocess, sys, threading, time

binary = pathlib.Path(sys.argv[1]).resolve()
root = pathlib.Path(sys.argv[2]).resolve()
root.mkdir(parents=True, exist_ok=False)
(root / 'cr').mkdir()
(root / 'workspace').mkdir()
settings = root / 'cr/settings.json'
initial = dict(theme='dark', tts_enabled=True, tts_speed=1.0, font_size=16, max_tool_rounds=2)
settings.write_text(json.dumps(initial))
entered, release = threading.Event(), threading.Event()
rounds, requests = {}, []
system_hashes = set()

class Provider(http.server.BaseHTTPRequestHandler):
    def log_message(self, *_): pass
    def do_POST(self):
        req = json.loads(self.rfile.read(int(self.headers['Content-Length'])))
        system_hashes.add(hashlib.sha256(json.dumps(req.get('system'),sort_keys=True).encode()).hexdigest())
        assert len(system_hashes)==1, 'settings changed base/system prompt'
        for msg in req['messages']:
            if msg['role'] != 'user': continue
            content = msg['content'] if isinstance(msg['content'],list) else [{'type':'text','text':msg['content']}]
            for part in content:
                if part.get('type') == 'text':
                    assert part.get('text') in ('FIRST_BUDGET','LATER_18','ZERO_20'), 'non-turn/settings context leaked into user text'
        label = next((part['text'] for msg in reversed(req['messages']) if msg['role'] == 'user'
                      for part in (msg['content'] if isinstance(msg['content'], list) else [{'type': 'text', 'text': msg['content']}])
                      if part.get('type') == 'text' and part.get('text') in ('FIRST_BUDGET', 'LATER_18', 'ZERO_20')), None)
        if label is None: raise AssertionError('missing turn label')
        n = rounds[label] = rounds.get(label, 0) + 1
        requests.append(dict(label=label, round=n, stream=req.get('stream'), path=self.path))
        if label == 'FIRST_BUDGET' and n == 1:
            entered.set()
            if not release.wait(30): raise AssertionError('response boundary not released')
        target = {'FIRST_BUDGET': 20, 'LATER_18': 18, 'ZERO_20': 20}[label]
        blocks = [{'type': 'text', 'text': 'Done.'}]
        if n <= target:
            blocks = [{'type': 'text', 'text': 'I will record this real tool batch in the owned workspace.'},
                      {'type': 'tool_use', 'id': f'{label}-{n}', 'name': 'write_file',
                       'input': {'path': label + '.txt', 'content': str(n) + '\n', 'append': True}}]
        self.send_response(200)
        self.send_header('Content-Type', 'text/event-stream')
        self.end_headers()
        def frame(obj):
            self.wfile.write(('data: ' + json.dumps(obj) + '\n\n').encode()); self.wfile.flush()
        frame({'type': 'message_start', 'message': {'usage': {'input_tokens': 1, 'output_tokens': 0}}})
        for i, b in enumerate(blocks):
            start = dict(b)
            if b['type'] == 'text': start['text'] = ''
            else: start['input'] = {}
            frame({'type': 'content_block_start', 'index': i, 'content_block': start})
            delta = ({'type': 'text_delta', 'text': b['text']} if b['type'] == 'text' else
                     {'type': 'input_json_delta', 'partial_json': json.dumps(b['input'])})
            frame({'type': 'content_block_delta', 'index': i, 'delta': delta})
            frame({'type': 'content_block_stop', 'index': i})
        frame({'type': 'message_delta', 'usage': {'output_tokens': 1}})
        frame({'type': 'message_stop'})

class WS:
    def __init__(self, addr):
        self.sock = socket.create_connection(addr, timeout=20)
        self.file = self.sock.makefile('rb')
        key = base64.b64encode(b'0123456789abcdef').decode()
        self.sock.sendall((f'GET /ws HTTP/1.1\r\nHost: {addr[0]}:{addr[1]}\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Version: 13\r\nSec-WebSocket-Key: {key}\r\n\r\n').encode())
        status = self.file.readline()
        assert b'101' in status, status
        while self.file.readline() != b'\r\n': pass
        self.events = []
        self.send({'type': 'subscribe'})
    def send(self, obj):
        raw = json.dumps(obj).encode(); mask = b'abcd'
        size = bytes([0x80 | len(raw)]) if len(raw) < 126 else bytes([0x80 | 126]) + struct.pack('!H', len(raw))
        self.sock.sendall(b'\x81' + size + mask + bytes(c ^ mask[i % 4] for i, c in enumerate(raw)))
    def read(self):
        header = self.file.read(2); assert len(header) == 2, 'closed WebSocket'
        n = header[1] & 127
        if n == 126: n = struct.unpack('!H', self.file.read(2))[0]
        if n == 127: n = struct.unpack('!Q', self.file.read(8))[0]
        msg = json.loads(self.file.read(n)); self.events.append(msg); return msg
    def until(self, kind):
        while True:
            msg = self.read()
            if msg['type'] == kind: return msg
    def snapshot(self, kind):
        got = self.until(kind)['settings']; assert set(got) == set(initial), got
        return got
    def turn(self):
        events = []
        while True:
            msg = self.read(); events.append(msg)
            if msg['type'] == 'state_changed' and msg['to'] == 'idle': return events
    def close(self): self.file.close(); self.sock.close()

provider = http.server.ThreadingHTTPServer(('127.0.0.1', 0), Provider)
thread = threading.Thread(target=provider.serve_forever); thread.start()
# Deliberate allowlist: no inherited credentials/provider/LLM override variables.
env = {k: os.environ[k] for k in ('PATH', 'TMPDIR', 'LANG') if k in os.environ}
env.update(ANTHROPIC_API_KEY='local-fixture-only', HOME=str(root))
commands = []
def start(number):
    args = [str(binary), '-port', '0', '-workspace', str(root / 'workspace'),
            '-vendor', 'anthropic', '-model', 'claude-sonnet-4-6', '-base-url', f'http://127.0.0.1:{provider.server_port}']
    commands.append({'argv': args, 'cwd': str(root), 'environment_names': sorted(env)})
    output = open(root / f'process-{number}.log', 'w')
    process = subprocess.Popen(args, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                               text=True, cwd=root, env=env, bufsize=1)
    address, ready = [], threading.Event()
    def drain():
        for line in process.stdout:
            output.write(line); output.flush()
            match = re.search(r'GUI http://127\.0\.0\.1:(\d+)', line)
            if match: address.append(('127.0.0.1', int(match.group(1)))); ready.set()
    reader = threading.Thread(target=drain); reader.start()
    owned.append((process, reader, output))
    assert ready.wait(20), 'GUI did not announce its actual listener (see owned process log)'
    return process, address[0]
def stop(process):
    process.stdin.close()
    process.terminate() # graceful SIGTERM; terminal EOF only detaches the real app
    assert process.wait(timeout=20) == 0, 'application failed shutdown'

owned, sockets = [], []
try:
    process, addr = start(1)
    a, b = WS(addr), WS(addr); sockets += [a, b]
    assert a.snapshot('current_settings') == initial
    assert b.snapshot('current_settings') == initial
    a.until('replay_end'); b.until('replay_end')
    durable_before = settings.read_bytes()
    # A valid sibling must not partially apply when another field is invalid.
    a.send({'type':'update_settings','settings':{'theme':'light','max_tool_rounds':'invalid'}})
    assert a.until('settings_error').get('text'), 'missing invalid-patch error'
    assert settings.read_bytes() == durable_before
    probe=WS(addr); sockets.append(probe); assert probe.snapshot('current_settings')==initial; probe.close(); sockets.remove(probe)
    # Force the real final rename to fail, without mocking saveSettings.
    backup=settings.with_name('review-original-settings.json'); settings.rename(backup); settings.mkdir()
    try:
        a.send({'type':'update_settings','settings':{'theme':'system','max_tool_rounds':7}})
        assert a.until('settings_error').get('text'), 'missing real-save-failure error'
        assert backup.read_bytes()==durable_before
        probe=WS(addr); sockets.append(probe); assert probe.snapshot('current_settings')==initial; probe.close(); sockets.remove(probe)
    finally:
        settings.rmdir(); backup.rename(settings)
    process.stdin.write('FIRST_BUDGET\n'); process.stdin.flush()
    assert entered.wait(20), 'terminal did not reach held provider boundary'
    a.send({'type': 'update_settings', 'settings': {'theme': 'light', 'tts_enabled': False, 'tts_speed': 1.2345, 'font_size': 18, 'max_tool_rounds': 18}})
    expected = dict(theme='light', tts_enabled=False, tts_speed=1.2345, font_size=18, max_tool_rounds=18)
    assert a.snapshot('settings_changed') == expected
    assert b.snapshot('settings_changed') == expected
    assert json.loads(settings.read_text()) == expected, 'published before disk commit'
    release.set()
    first = a.turn()
    assert sum(e['type'] == 'tool_dispatched' for e in first) == 2, 'live update rewrote active turn'
    assert (root / 'workspace/FIRST_BUDGET.txt').read_text().splitlines() == ['1', '2']
    assert rounds['FIRST_BUDGET'] == 3
    for ws in sockets: ws.close()
    sockets.clear(); stop(process)
    # Fresh real process restores accepted settings before either view may begin work.
    process, addr = start(2)
    a, b = WS(addr), WS(addr); sockets += [a, b]
    assert a.snapshot('current_settings') == expected
    assert b.snapshot('current_settings') == expected
    a.until('replay_end'); b.until('replay_end')
    a.send({'type': 'prompt', 'text': 'LATER_18'})
    later = a.turn()
    assert sum(e['type'] == 'tool_dispatched' for e in later) == 18
    assert len((root / 'workspace/LATER_18.txt').read_text().splitlines()) == 18
    a.send({'type': 'update_settings', 'settings': {'max_tool_rounds': 0}})
    expected['max_tool_rounds'] = 0
    assert a.snapshot('settings_changed') == expected
    assert b.snapshot('settings_changed') == expected
    assert json.loads(settings.read_text()) == expected
    process.stdin.write('ZERO_20\n'); process.stdin.flush()
    zero = a.turn()
    assert sum(e['type'] == 'tool_dispatched' for e in zero) == 20
    assert len((root / 'workspace/ZERO_20.txt').read_text().splitlines()) == 20
    # Reconnect gets the full accepted zero/false values, without resetting state.
    c = WS(addr); sockets.append(c)
    assert c.snapshot('current_settings') == expected
    for ws in sockets: ws.close()
    sockets.clear(); stop(process)
    process, addr = start(3)
    a = WS(addr); sockets.append(a)
    assert a.snapshot('current_settings') == expected
    a.close(); sockets.clear(); stop(process)
    assert all(request['stream'] for request in requests), 'fixture failed to exercise streaming'
    print('PASS real built GUI: terminal turn holds budget2 across live WS update; persisted18 survives process restart and GUI executes18; zero allows20 terminal batches; full false/zero snapshots, two peers, reconnect and second restart. Actual file tools, fake loopback streaming provider. NOT browser/native/live-provider proof.')
finally:
    release.set()
    for ws in sockets: ws.close()
    for process, reader, output in owned:
        if process.poll() is None:
            if not process.stdin.closed: process.stdin.close()
            process.terminate()
            try: process.wait(timeout=20)
            except subprocess.TimeoutExpired: process.terminate(); process.wait(timeout=20)
        reader.join(timeout=5); assert not reader.is_alive(), 'owned log drainer not finished'
        output.close()
    provider.shutdown(); provider.server_close(); thread.join()
    (root / 'review-privacy-summary.json').write_text(json.dumps({'distinct_system_hashes':sorted(system_hashes),'user_text_allowlist_asserted_per_request':True,'requests':len(requests)},indent=2)+'\n')
    (root / 'requests.json').write_text(json.dumps(requests, indent=2) + '\n')
    (root / 'commands.json').write_text(json.dumps(commands, indent=2) + '\n')
