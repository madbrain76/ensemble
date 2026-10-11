#!/usr/bin/env bash
#
# Mutation audit for the chapter 26 grader.
#
# Each mutant deletes one behavior and asserts the exact set of checks that
# dies. A survivor means the check is decorative; a compile failure means the
# mutant is invalid (a script bug, not a finding).
#
# Usage: ./scripts/ch26-mutants.sh
set -uo pipefail
cd "$(dirname "$0")/.."

OUT=$(mktemp -d)
GRADER="$OUT/grade"
PASSED=0
FAILED=0

cleanup() { git checkout -- agent/ >/dev/null 2>&1 || true; }
trap cleanup EXIT

echo "building the grader once, so a mutant cannot change it"
go build -o "$GRADER" ./cmd/grade || exit 1

run_mutant() {
	local n="$1" name="$2" expect="$3" files="$4" apply="$5" verify="$6"

	echo
	echo "=== mutant $n: $name"
	echo "  expecting to kill: $expect"

	git checkout -- $files >/dev/null 2>&1
	eval "$apply"

	if ! eval "$verify"; then
		echo "  MUTATION DID NOT LAND - invalid mutant (this is a script bug)"
		git checkout -- $files
		FAILED=$((FAILED + 1))
		return
	fi
	echo "  mutation landed."

	if ! (cd agent && go build ./... >/dev/null 2>&1); then
		echo "  MUTANT DOES NOT COMPILE - invalid mutant (this is a script bug)"
		git checkout -- $files
		FAILED=$((FAILED + 1))
		return
	fi
	echo "  mutant compiles."

	local report="$OUT/report-$n.json"
	"$GRADER" -ch 26 -json ./agent >"$report" 2>"$OUT/stderr-$n.log"

	local killed
	killed=$(
		python3 - "$report" <<'PY'
import json, sys
def walk(x, out):
    if isinstance(x, dict):
        if "id" in x and "passed" in x:
            if not x["passed"]:
                out.append(x["id"])
        for v in x.values():
            walk(v, out)
    elif isinstance(x, list):
        for v in x:
            walk(v, out)
out = []
try:
    walk(json.load(open(sys.argv[1])), out)
except Exception as e:
    print("PARSE-ERROR", e)
    sys.exit(0)
print(" ".join(sorted(set(out))))
PY
	)

	local want
	want=$(echo "$expect" | tr ' ' '\n' | sort -u | tr '\n' ' ' | sed 's/ $//')
	local got
	got=$(echo "$killed" | tr ' ' '\n' | sort -u | tr '\n' ' ' | sed 's/ $//')

	if [ "$got" = "$want" ]; then
		echo "  KILLED exactly: $got"
		PASSED=$((PASSED + 1))
	else
		echo "  MISMATCH"
		echo "    wanted: $want"
		echo "    got:    $got"
		FAILED=$((FAILED + 1))
	fi

	git checkout -- $files >/dev/null 2>&1
}

# 1. Activity tracking is removed from the status tool. The observer still
#    emits, but get_sub_agent_status no longer reports it, so the grader's
#    probe finds no "activity" in the result.
run_mutant 1 "get_sub_agent_status omits activity" \
	"activity-observed" \
	"agent/management_tools.go" \
	"perl -i -pe 's/Activity: %s/_hidden_: %s/' agent/management_tools.go" \
	"grep -q '_hidden_' agent/management_tools.go"

# 2. Named agents from the persistent registry are never instantiated on
#    restart. The registry loads, but Start() is removed, so a create-restart
#    cycle finds no agent directory.
run_mutant 2 "en.Start removed from composition root" \
	"ensemble-persists" \
	"agent/cmd/main.go" \
	"perl -i -pe 's/en\.Start\(\)/\/\/ en.Start()/' agent/cmd/main.go" \
	"grep -q '// en.Start()' agent/cmd/main.go"

# 3. SubAgentManager.Shutdown accepts unknown agents instead of rejecting.
#    The grader calls shutdown_agent with a non-child and expects an error;
#    this mutant makes it silently succeed.
run_mutant 3 "shutdown accepts unknown agents" \
	"ownership-enforced" \
	"agent/subagent.go" \
	"perl -i -pe 's/return fmt\.Errorf\(\"no such sub-agent: %s\", id\)/return nil \/\/ mutant/' agent/subagent.go" \
	"grep -q 'return nil // mutant' agent/subagent.go"

# 4. Status-change events are reported with a wrong kind. The poll loop still
#    detects transitions, but labels them as "heartbeat" instead of
#    "status_change", so the grader's pattern match misses them.
run_mutant 4 "wait_for_agent_change mislabels status changes" \
	"wait-change-wakes" \
	"agent/management_tools.go" \
	"perl -i -pe 's/\"status_change\"/\"heartbeat\"/' agent/management_tools.go" \
	"grep -q 'heartbeat' agent/management_tools.go"

# 5. Both tree-broadcasting paths are disabled in main.go: the subscribe-time
#    tree callback (h.tree) is set to nil, and the BroadcastTree call is
#    removed. Neither path sends agent_tree to the WebSocket client.
run_mutant 5 "no tree data reaches the GUI" \
	"gui-tree-message" \
	"agent/cmd/main.go" \
	"perl -i -pe 's/Tree:\s+func\(\) any \{ return en\.TreeSnapshot\(\) \},/\/\/ Tree removed by mutant/' agent/cmd/main.go; perl -i -pe 's/guiServer\.BroadcastTree\(en\.TreeSnapshot\(\)\)/\/\/ mutant: tree disabled/' agent/cmd/main.go" \
	"grep -q 'Tree removed by mutant' agent/cmd/main.go && grep -q 'mutant: tree disabled' agent/cmd/main.go"

# 6. Submit rejection with active children is disabled. The parent submits
#    while the child is still active, and the grader expects a rejection but
#    gets a silent success.
run_mutant 6 "submit does not reject with active children" \
	"submit-idle-regression" \
	"agent/workflow_tools.go" \
	"perl -i -pe 's/if ids := a\.ActiveChildIDs\(\); len\(ids\) > 0/if ids := a.ActiveChildIDs(); false \&\& len(ids) > 0/' agent/workflow_tools.go" \
	"grep -q 'false && len(ids)' agent/workflow_tools.go"

# 7. An idle child is never woken. The framework notices it stopped without
#    submitting and does nothing, so the child takes exactly 1 turn (its own)
#    instead of the expected 4.
run_mutant 7 "idle correction disabled" \
	"submit-idle-regression" \
	"agent/idle.go" \
	"perl -i -pe 's/s\.Start\(wakeText\(n\)\)/_ = n/' agent/idle.go" \
	"! grep -q 's.Start(wakeText(n))' agent/idle.go"

echo
echo "================================"
echo "killed cleanly: $PASSED"
echo "problems:      $FAILED"
echo "================================"
[ "$FAILED" -eq 0 ]
