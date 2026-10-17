#!/usr/bin/env python3
"""Finite Chapter 2 mutation checks; edits are always restored, never committed.
Run from the student repository: python3 scripts/check-mutations.py ABS_EVIDENCE_DIR
No live calls, grader imports, credentials or production testing hooks are used.
"""
import hashlib
import pathlib
import subprocess
import sys

root = pathlib.Path(__file__).resolve().parent.parent
out = pathlib.Path(sys.argv[1]).resolve()
out.mkdir(parents=True, exist_ok=True)
cases = [
    ("ephemera", "internal/history/history.go", "c.Ephemera = []common.Part{}", "c.Ephemera = c.Ephemera", "TestChapter2History"),
    ("cache-write", "internal/engine/parse.go", "wire.Usage.Input-wire.Usage.Details.Cached-wire.Usage.Details.Written", "wire.Usage.Input-wire.Usage.Details.Cached", "TestChapter2Parse"),
    ("opaque-isolation", "internal/engine/render.go", 'func anthropicParts(parts []common.Part, to common.Provenance) []json.RawMessage {', 'func anthropicParts(parts []common.Part, to common.Provenance) []json.RawMessage {\n for _,p:=range parts {if p.Type=="opaque" {to=p.From;break}}', "TestChapter2ParseAndProvenance"),
    ("anthropic-order", "internal/engine/render.go", "result.Messages[i].Content = append(results, rest...)", "result.Messages[i].Content = append(rest, results...)", "TestChapter2AnthropicOrderingAndIDs"),
    ("ref-stub", "internal/history/history.go", "size), Ref: p.Ref}", "size)}", "TestChapter2Redaction"),
]
passed = True
for name, filename, old, new, test in cases:
    path = root / filename
    original = path.read_bytes()
    source = original.decode()
    if source.count(old) != 1:
        raise RuntimeError(f"{name}: mutation anchor is not unique")
    record = f"cwd: {root}\ncommand: go test . -count=1 -run {test}\nfile: {filename}\noriginal sha256: {hashlib.sha256(original).hexdigest()}\nmutation: {old!r} -> {new!r}\n"
    try:
        path.write_text(source.replace(old, new))
        failed = subprocess.run(["go", "test", ".", "-count=1", "-run", test], cwd=root, capture_output=True, text=True)
        record += f"\nMUTATED exit={failed.returncode}\n{failed.stdout}{failed.stderr}"
    finally:
        path.write_bytes(original)
    restored = subprocess.run(["go", "test", ".", "-count=1", "-run", test], cwd=root, capture_output=True, text=True)
    record += f"\nRESTORED exit={restored.returncode}\n{restored.stdout}{restored.stderr}\nrestored sha256: {hashlib.sha256(path.read_bytes()).hexdigest()}\n"
    ok = failed.returncode != 0 and "--- FAIL:" in failed.stdout and restored.returncode == 0 and "no tests to run" not in restored.stdout and original == path.read_bytes()
    record += f"behavioral failure and restored pass: {ok}\n"
    (out / f"mutation-{name}.txt").write_text(record)
    print(name, "PASS" if ok else "FAIL")
    passed = passed and ok
sys.exit(0 if passed else 1)
