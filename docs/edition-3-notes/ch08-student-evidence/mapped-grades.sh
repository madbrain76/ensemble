#!/bin/sh
set -eu
root=/Users/bill/projects/ensemble-edition-3-student
run() { python3 "$root/.ch08-local/run_gate.py" "$@"; }
run grade-ch08-mapped /Users/bill/projects/ensemble /tmp/cadence-ch08-grader -ch 8 -json -core "$root" -consumer "$root/examples/ch04-consumer" -workflow "$root/examples/ch06" -stream-exercise "$root/examples/ch07" -stream-model claude-sonnet-4-6 -gui "$root/gui" -gui-model claude-sonnet-4-6 -media-models "$root/.ch08-local/media-models.json" "$root/.ch08-local/bin/root-ch08"
run grade-ch01-final /Users/bill/projects/ensemble /tmp/cadence-ch08-grader -ch 1 -json "$root/.ch08-local/bin/root-ch08"
