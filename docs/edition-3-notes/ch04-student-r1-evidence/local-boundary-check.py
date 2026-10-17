"""C1 local boundary verification; no paid provider or credential-store access.
Supplies deterministic Anthropic-shaped HTTP responses on loopback, but uses
actual compiled application processes, Agent dispatch, jobs, files and shell.
"""
import datetime
import hashlib
import http.server
import json
import os
import pathlib
import subprocess
import threading
import tempfile
import time

BASE = pathlib.Path(__file__).resolve().parent
REPO = pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
BLOB = 'abcde' * 8192


def text(block):
    content = block.get('content', '')
    if isinstance(content, str):
        return content
    return ''.join(part.get('text', '') for part in content)


def check_one(mode):
    binary = pathlib.Path('/tmp/ch04-r1-' + mode)
    work = BASE / ('local-' + mode)
    work.mkdir(exist_ok=False)
    temporary = pathlib.Path(tempfile.mkdtemp(prefix='ch04-r1-local-' + mode + '-'))
    workspace = temporary / 'workspace'
    workspace.mkdir()
    (workspace / 'large.txt').write_text(BLOB)
    data_dir = temporary / 'data'
    data_dir.mkdir()
    log = work / 'events.jsonl'
    operation = ('read_file', {'path': 'large.txt'}) if mode == 'core' else ('shout', {'text': BLOB})
    batches = [
        [('tool_limits', {'max_output_bytes': 1024}), ('tool_limits', {'ai_callback_delay': 0.3}), operation],
        [('tool_limits', {'max_output_bytes': 1024}), ('tool_limits', {'ai_callback_delay': -1}), operation],
        [('tool_limits', {'max_output_bytes': 1024}), ('tool_limits', {'max_output_bytes': 2048}), ('run_command', {'command': 'cat large.txt', 'max_output_bytes': 8192})],
    ]
    observations = []
    failures = []
    requests_seen = []

    class Handler(http.server.BaseHTTPRequestHandler):
        def log_message(self, *_):
            pass

        def do_POST(self):
            try:
                request = json.loads(self.rfile.read(int(self.headers['Content-Length'])))
                n = len(requests_seen)
                requests_seen.append({'path': self.path, 'model': request.get('model'), 'tool_names': [x['name'] for x in request.get('tools', [])]})
                assert self.path == '/v1/messages'
                assert request['model'] == 'local-c1-boundary-fixture'
                assert 'tool_limits' in requests_seen[-1]['tool_names']
                if mode == 'consumer':
                    assert 'shout' in requests_seen[-1]['tool_names']
                if n:
                    results = [p for p in request['messages'][-1]['content'] if p.get('type') == 'tool_result']
                    assert len(results) == 3, ('result_count', len(results))
                    output = text(results[2])
                    entry = {'batch': n - 1, 'result_bytes': len(output.encode()), 'returns': results}
                    observations.append(entry)
                    assert 'tool_limits consumed by this tool_limits call' in text(results[1])
                    if n == 1:
                        assert 'max_output_bytes 1024' in text(results[1])
                        assert 'max_output_bytes 16384' in output
                        assert 'ai_callback_delay 300ms' in output
                        assert 'tool_limits consumed by this ' + operation[0] + ' call' in output
                        assert 8192 < len(output) < 17000, len(output)
                        assert 'bytes omitted' in output
                    elif n == 2:
                        assert results[1].get('is_error') is True
                        assert 'error:' in text(results[1])
                        assert 'tool_limits consumed' not in output
                        assert 8192 < len(output) < 17000, len(output)
                    elif n == 3:
                        assert 'max_output_bytes 2048' in output
                        assert 4096 < len(output) < 8500, len(output)
                        assert 'exit_code: 0' in output
                    else:
                        raise AssertionError('unexpected extra request')
                if n < len(batches):
                    content = [{'type': 'text', 'text': 'I will exercise the one-shot settings and inspect the actual result to check this local correction.'}]
                    for index, (name, args) in enumerate(batches[n]):
                        content.append({'type': 'tool_use', 'id': 'local-%d-%d' % (n, index), 'name': name, 'input': args})
                    stop = 'tool_use'
                else:
                    content = [{'type': 'text', 'text': 'LOCAL C1 CHECKS PASSED'}]
                    stop = 'end_turn'
                response = {'id': 'local-c1-fixture', 'type': 'message', 'role': 'assistant', 'model': 'local-c1-boundary-fixture', 'content': content, 'stop_reason': stop, 'usage': {'input_tokens': 1, 'output_tokens': 1}}
                payload = json.dumps(response).encode()
                self.send_response(200)
            except Exception as exc:
                failures.append(repr(exc))
                payload = json.dumps({'error': {'message': 'local fixture assertion failed: ' + repr(exc)}}).encode()
                self.send_response(500)
            self.send_header('Content-Type', 'application/json')
            self.send_header('Content-Length', str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)

    server = http.server.ThreadingHTTPServer(('127.0.0.1', 0), Handler)
    worker = threading.Thread(target=server.serve_forever)
    worker.start()
    env = {key: os.environ[key] for key in ['PATH', 'HOME', 'TMPDIR', 'SHELL'] if key in os.environ}
    env.update(LLM_VENDOR='anthropic', LLM_BASE_URL='http://127.0.0.1:' + str(server.server_port), LLM_MODEL='local-c1-boundary-fixture', LLM_API_KEY='LOCAL-PLACEHOLDER-NOT-A-CREDENTIAL', LLM_REQUIRE_TOOL_EXPLANATION='true', ENSEMBLE_WORKSPACE=str(workspace), ENSEMBLE_DATA_DIR=str(data_dir), CH02_LOG=str(log))
    command = [str(binary)] + (['chat'] if mode == 'core' else [])
    started = datetime.datetime.now(datetime.timezone.utc).isoformat()
    begin = time.monotonic()
    try:
        result = subprocess.run(command, input='Exercise the local one-shot settings fixture.\n' if mode == 'core' else '', env=env, cwd=workspace, capture_output=True, text=True)
    finally:
        server.shutdown()
        server.server_close()
        worker.join()
    elapsed = time.monotonic() - begin
    (work / 'stdout.txt').write_text(result.stdout)
    (work / 'stderr.txt').write_text(result.stderr)
    records = [json.loads(line) for line in log.read_text().splitlines()]
    assert records[0] == {'log_version': 1}, records[0]
    events = records[1:]
    returned = [e['tool'] for e in events if e['type'] == 'tool_returned' and e['tool']['name'] in [operation[0], 'run_command']]
    try:
        assert result.returncode == 0, result.returncode
        assert not failures, failures
        assert len(requests_seen) == 4, len(requests_seen)
        assert 'LOCAL C1 CHECKS PASSED' in result.stdout
        assert len(returned) == 3, len(returned)
        for item in returned:
            assert item['job']['status'] == 'done', item['job']
            locator = item['job']['output']['locator']
            raw = (data_dir / locator).read_text()
            assert len(raw) >= 40960, (locator, len(raw))
        if mode == 'consumer':
            assert result.stderr.count('shout called') == 2
    except Exception as exc:
        failures.append(repr(exc))
    receipt = {'mode': mode, 'classification': 'local HTTP provider-boundary fixture, NOT paid/live vendor evidence', 'command': command, 'started_utc': started, 'elapsed_seconds': elapsed, 'exit_code': result.returncode, 'binary_sha256': hashlib.sha256(binary.read_bytes()).hexdigest(), 'commit': subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=REPO, text=True).strip(), 'tree': subprocess.check_output(['git', 'rev-parse', 'HEAD^{tree}'], cwd=REPO, text=True).strip(), 'requests': requests_seen, 'observations': observations, 'failures': failures, 'passed': not failures, 'actual_ordinary_jobs_done': [item['job'] for item in returned], 'workspace': str(workspace), 'data_dir': str(data_dir), 'fresh_vendor_calls': False, 'credential_store_accessed': False}
    (work / 'receipt.json').write_text(json.dumps(receipt, indent=2) + '\n')
    print(json.dumps({key: receipt[key] for key in ['mode', 'passed', 'exit_code', 'elapsed_seconds', 'failures']}))
    assert not failures


for mode in ['core', 'consumer']:
    check_one(mode)
