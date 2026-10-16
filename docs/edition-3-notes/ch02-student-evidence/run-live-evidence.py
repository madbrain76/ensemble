#!/usr/bin/env python3
"""Bounded, authorized Chapter 2 CLI evidence. No keys are printed or saved.
Only the three explicitly authorized settings fields are selected into launch
memory. Outputs are normalized student logs, not raw vendor payloads/requests.
"""
import collections
import hashlib
import http.server
import json
import os
import pathlib
import re
import subprocess
import threading

OUT = pathlib.Path(__file__).resolve().parent
CWD = pathlib.Path('/Users/bill/projects/ensemble-edition-3-student')
BINARY = OUT / 'student-cli'
EXPECTED = '1f6494edd754535ff5d1f2ca0f2663ccfa4ea236a5655dd8894b1940dc70f7e4'
assert hashlib.sha256(BINARY.read_bytes()).hexdigest() == EXPECTED
COMMIT = subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=CWD, text=True).strip()
assert COMMIT == '74eb0265bd8bdbc7857d579704300b317b2365da'
assert not subprocess.check_output(['git', 'status', '--porcelain'], cwd=CWD)
providers = [
    ('anthropic', 'claude-haiku-4-5-20251001', 'https://api.anthropic.com', 'directClaudeAPIKey'),
    ('openai', 'gpt-4.1-nano', 'https://api.openai.com', 'directOpenAIAPIKey'),
    ('gemini', 'gemini-2.5-flash-lite', 'https://generativelanguage.googleapis.com', 'directGeminiAPIKey'),
]
with pathlib.Path.home().joinpath('.cr/settings.json').open() as handle:
    store = json.load(handle)
keys = {field: store.get(field, '') for _, _, _, field in providers}
del store
if not all(isinstance(value, str) and value for value in keys.values()):
    raise RuntimeError('One or more authorized key fields are unavailable; no call made')
base_env = {key: os.environ[key] for key in ['PATH', 'HOME', 'TMPDIR', 'LANG'] if key in os.environ}
prompts = (
    "I'm preparing a short teaching note for a chat application. Explain in two sentences why an append-only event log differs from a derived context view. Use a concrete chat-history example.\n"
    "Using your example, give one sentence explaining why rendering that saved log offline should not contact the model again.\n"
)
(OUT / 'live-prompts.txt').write_text(prompts)
results = []

def save_safe(path, data):
    if any(key.encode() in data for key in keys.values()):
        raise RuntimeError('Credential-like output detected; refusing to archive')
    path.write_bytes(data)

for vendor, model, endpoint, field in providers:
    log = OUT / f'live-{vendor}.jsonl'
    env = dict(base_env, LLM_VENDOR=vendor, LLM_MODEL=model, LLM_BASE_URL=endpoint,
               LLM_API_KEY=keys[field], CH02_LOG=str(log))
    command = [str(BINARY), 'chat']
    # Exactly one CLI process, two turns, no retries and no elapsed-time cap.
    process = subprocess.run(command, input=prompts.encode(), cwd=CWD, env=env, capture_output=True)
    del env
    save_safe(OUT / f'live-{vendor}.stdout.txt', process.stdout)
    save_safe(OUT / f'live-{vendor}.stderr.txt', process.stderr)
    result = dict(vendor=vendor, routing_model=model, command=command, cwd=str(CWD),
                  exit=process.returncode, commit=COMMIT, binary_sha256=EXPECTED,
                  live=True, automatic_retries=0, planned_turns=2, output_token_limit_per_call=1024)
    if log.exists():
        raw = log.read_bytes()
        save_safe(log, raw)
        facts = [json.loads(line) for line in raw.splitlines()][1:]
        result['event_counts'] = dict(collections.Counter(fact['type'] for fact in facts))
        usage = {name: 0 for name in ['input', 'cache_write', 'cache_read', 'output']}
        for fact in facts:
            if fact['type'] == 'response_ended':
                for name in usage:
                    usage[name] += fact['response']['usage'][name]
        result['captured_usage'] = usage
        result['returned_models'] = sorted({fact['response']['from']['model'] for fact in facts if fact['type'] == 'response_ended'})
        result['log_sha256'] = hashlib.sha256(raw).hexdigest()
        match = re.search(rb'Usage: input=(\d+) cache_write=(\d+) cache_read=(\d+) output=(\d+)', process.stderr)
        result['transcript_usage_matches_log'] = bool(match and list(map(int, match.groups())) == list(usage.values()))
    results.append(result)
    (OUT / 'live-summary.json').write_text(json.dumps(results, indent=2) + '\n')
    print(vendor, 'exit', process.returncode, 'captured usage', result.get('captured_usage'))
    if process.returncode != 0:
        print('Provider blocker recorded; not retrying.')

# Each offline command is a fresh process with no API key. Point all target
# endpoints at a local counter so unintended transport is additionally visible.
requests = []
class Sentinel(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        requests.append(self.path)
        self.send_response(500)
        self.end_headers()
    def log_message(self, *args):
        pass
server = http.server.ThreadingHTTPServer(('127.0.0.1', 0), Sentinel)
thread = threading.Thread(target=server.serve_forever, daemon=True)
thread.start()
offline = []
try:
    for result in results:
        vendor = result['vendor']
        log = OUT / f'live-{vendor}.jsonl'
        if not log.exists():
            continue
        env = dict(base_env, CH02_LOG=str(log))
        command = [str(BINARY), 'dump']
        process = subprocess.run(command, cwd=CWD, env=env, capture_output=True)
        save_safe(OUT / f'offline-{vendor}-dump.jsonl', process.stdout)
        save_safe(OUT / f'offline-{vendor}-dump.stderr.txt', process.stderr)
        offline.append(dict(command=command, cwd=str(CWD), source_vendor=vendor,
                            exit=process.returncode, exact_saved_log=process.stdout == log.read_bytes()))
        for target_vendor, target_model, _, _ in providers:
            env = dict(base_env, LLM_VENDOR=target_vendor, LLM_MODEL=target_model,
                       LLM_BASE_URL=f'http://127.0.0.1:{server.server_port}')
            command = [str(BINARY), 'render', str(log)]
            runs = [subprocess.run(command, cwd=CWD, env=env, capture_output=True) for _ in range(2)]
            for i, process in enumerate(runs):
                save_safe(OUT / f'offline-{vendor}-to-{target_vendor}-{i+1}.json', process.stdout)
                save_safe(OUT / f'offline-{vendor}-to-{target_vendor}-{i+1}.stderr.txt', process.stderr)
            offline.append(dict(command=command, cwd=str(CWD), source_vendor=vendor,
                                target_vendor=target_vendor, target_model=target_model,
                                exits=[p.returncode for p in runs], byte_identical=runs[0].stdout == runs[1].stdout,
                                request_json_valid=all(bool(json.loads(p.stdout)) for p in runs if p.returncode == 0),
                                output_sha256=hashlib.sha256(runs[0].stdout).hexdigest()))
finally:
    server.shutdown()
    server.server_close()
(OUT / 'offline-summary.json').write_text(json.dumps(dict(commit=COMMIT, binary_sha256=EXPECTED,
    api_keys_supplied=False, sentinel_post_count=len(requests), commands=offline), indent=2) + '\n')
print('offline sentinel POST count', len(requests))
assert not requests
assert all(result['exit'] == 0 and result.get('transcript_usage_matches_log') and result.get('event_counts',{}).get('response_ended') == 2 for result in results), 'Live blocker recorded'
assert all(item.get('exit', 0) == 0 and all(code == 0 for code in item.get('exits', [])) and item.get('byte_identical', True) and item.get('exact_saved_log', True) for item in offline)
