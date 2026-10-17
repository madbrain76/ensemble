#!/bin/sh
set -u
OUT=/tmp/cadence-ch05-evidence-r1/resume-r2
ROOT=/Users/bill/projects/ensemble-edition-3-student
: > "$OUT/gate-exits.txt"
step() {
 name=$1; shift
 df -k /tmp >> "$OUT/space-gates.txt"
 free=$(df -k /tmp | awk 'NR==2 {print $4}')
 if [ "$free" -lt 1048576 ]; then echo "$name BLOCKED free_kib=$free" | tee -a "$OUT/gate-exits.txt"; exit 90; fi
 printf '%s command:' "$name" >> "$OUT/gate-exits.txt"; printf ' %s' "$@" >> "$OUT/gate-exits.txt"; printf '\n' >> "$OUT/gate-exits.txt"
 "$@" > "$OUT/$name.txt" 2>&1; rc=$?
 echo "$name exit=$rc" | tee -a "$OUT/gate-exits.txt"
 if [ "$rc" -ne 0 ]; then cat "$OUT/$name.txt"; exit "$rc"; fi
}
cd "$ROOT" || exit 1
step root-vet go vet ./...
step root-race go test -race -count=1 -timeout 4m -coverpkg=./... -coverprofile="$OUT/coverage.out" ./...
step root-build go build -race ./...
for module in gui examples/ch04-consumer examples/ch06; do
 cd "$ROOT/$module" || exit 1
 name=$(basename "$module")
 step "$name-vet" env GOWORK=off go vet ./...
 step "$name-race" env GOWORK=off go test -race -count=1 -timeout 90s ./...
 step "$name-build" env GOWORK=off go build -race -o "$OUT/$name" .
done
cd "$ROOT" || exit 1
step cli-build go build -race -o "$OUT/ensemble" ./cmd/ch01
step coverage-summary go tool cover -func="$OUT/coverage.out"
