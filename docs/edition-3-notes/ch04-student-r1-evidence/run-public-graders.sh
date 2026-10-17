#!/bin/sh
set -u
evidence=/Users/bill/projects/ensemble/docs/edition-3-notes/ch04-student-r1-evidence
cd /Users/bill/projects/ensemble || exit 1
run_grade() {
 name=$1; shift
 printf '\n%s start %s\n' "$name" "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
 printf 'command:'; printf ' %s' "$@"; printf '\n'
 "$@" > "$evidence/raw/$name.json" 2> "$evidence/raw/$name.stderr"
 code=$?
 printf '%s exit=%s end %s\n' "$name" "$code" "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
 cat "$evidence/raw/$name.json" "$evidence/raw/$name.stderr"
 [ "$code" -eq 0 ] || exit "$code"
 python3 -c 'import json,sys; r=json.load(open(sys.argv[1]));assert r["passed"] is True and all(c["passed"] for c in r["checks"])' "$evidence/raw/$name.json" || exit 1
}
run_grade public-chapter4 go run ./cmd/grade-edition3 -ch 4 -core /Users/bill/projects/ensemble-edition-3-student -consumer /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer -consumer-case /Users/bill/projects/ensemble-edition-3-student/examples/ch04-consumer/consumer-case.json -json /tmp/ch04-r1-core
run_grade public-chapter1 go run ./cmd/grade-edition3 -ch 1 -json /tmp/ch04-r1-core
printf 'both public graders complete %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
