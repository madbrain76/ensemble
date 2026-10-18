#!/bin/sh
set -eu
root=/Users/bill/projects/ensemble-edition-3-student
cd "$root"
run() { python3 "$root/.ch08-local/run_gate.py" "$@"; }
for spec in root:. gui:gui consumer:examples/ch04-consumer workflow:examples/ch06 stream:examples/ch07; do
 name=${spec%%:*}; dir=${spec#*:}; cwd="$root/$dir"
 run "$name-test" "$cwd" go test -count=1 ./...
 run "$name-race" "$cwd" go test -race -count=1 ./...
 run "$name-vet" "$cwd" go vet ./...
 run "$name-build" "$cwd" go build ./...
 case "$name" in root) target=./cmd/ch01;; gui) target=./cmd;; *) target=.;; esac
 run "$name-artifact" "$cwd" go build -o "$root/.ch08-local/bin/$name-ch08" "$target"
done
run gui-process-frozen "$root" python3 gui/test/settings_process.py "$root/.ch08-local/bin/gui-ch08" "$root/.ch08-local/process-frozen"
run grade-ch08 /Users/bill/projects/ensemble /tmp/cadence-ch08-grader -ch 8 -json -core "$root" -consumer "$root/examples/ch04-consumer" -workflow "$root/examples/ch06" -stream-exercise "$root/examples/ch07" -stream-model claude-sonnet-4-6 -gui "$root/gui" -gui-model claude-sonnet-4-6 "$root/.ch08-local/bin/root-ch08"
run grade-ch01 /Users/bill/projects/ensemble /tmp/cadence-ch08-grader -ch 1 -json "$root/.ch08-local/bin/root-ch08"
