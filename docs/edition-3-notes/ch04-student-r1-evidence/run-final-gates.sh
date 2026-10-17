#!/bin/sh
set -u
repo=/Users/bill/projects/ensemble-edition-3-student
evidence=/Users/bill/projects/ensemble/docs/edition-3-notes/ch04-student-r1-evidence
cd "$repo" || exit 1
run_gate() {
 name=$1; shift
 printf '\n%s start %s\n' "$name" "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
 printf 'command:'; printf ' %s' "$@"; printf '\n'
 "$@" > "$evidence/raw/$name.log" 2>&1
 code=$?
 printf '%s exit=%s end %s\n' "$name" "$code" "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
 cat "$evidence/raw/$name.log"
 [ "$code" -eq 0 ] || exit "$code"
}
run_gate full go test -count=1 ./...
run_gate race go test -race -count=1 ./...
run_gate vet go vet ./...
run_gate chapter4 /Users/bill/projects/ensemble/grader/cli/grade-ensemble -edition 3 -chapter 4 -core "$repo" -consumer "$repo/examples/ch04-consumer" -consumer-case "$repo/examples/ch04-consumer/consumer-case.json" -run /tmp/ch04-r1-core
python3 -c 'import json,sys; r=json.load(open(sys.argv[1]));assert r["passed"] is True and all(c["passed"] for c in r["checks"])' "$evidence/raw/chapter4.log" || exit 1
run_gate chapter1 /Users/bill/projects/ensemble/grader/cli/grade-ensemble -edition 3 -chapter 1 -run /tmp/ch04-r1-core
python3 -c 'import json,sys; r=json.load(open(sys.argv[1]));assert r["passed"] is True and all(c["passed"] for c in r["checks"])' "$evidence/raw/chapter1.log" || exit 1
printf 'all gates completed %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
