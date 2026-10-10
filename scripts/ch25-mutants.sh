#!/usr/bin/env bash
#
# Mutation audit for the chapter 25 grader.
#
# A grader that passes is evidence of nothing until each check has been shown
# to fail for its own reason. So delete one protected behavior at a time from
# the reference implementation and assert the exact set of checks that dies.
# A survivor means the check is decorative, or that the mutant was invalid and
# removed nothing. Those two look identical in a summary, which is why the
# script insists the mutation landed and that the mutant still compiles before
# it believes any result at all.
#
# Usage: ./scripts/ch25-mutants.sh
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

# run_mutant NUMBER NAME EXPECTED FILES APPLY VERIFY
#
# EXPECTED is the space separated set of check ids that must fail, and nothing
# else may fail. VERIFY is a shell test proving the edit landed, because a sed
# that silently matches nothing produces a green run that means nothing.
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
	"$GRADER" -ch 25 -json ./agent >"$report" 2>"$OUT/stderr-$n.log"

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

# 1. Submitting records nothing. The child still finishes and still reports a
#    clean outcome, so this separates "the child stopped" from "the child said
#    what it found", which is the whole point of a submitted payload.
run_mutant 1 "submit records no payload" \
	"spawn-reports-back join-fans-in" \
	"agent/fleet.go" \
	"perl -i -pe 's/a\.msg\.submission = &Submission\{.*\}/_ = payload/' agent/fleet.go" \
	"! grep -q 'a.msg.submission = &Submission' agent/fleet.go"

# 2. The child is handed its parent's directory instead of its own.
#
#    This one kills six checks, and the breadth is the finding rather than a
#    defect in the probe. A child sharing the parent's directory does not
#    merely lose its private copy of the conversation, it takes zero turns:
#    two agents writing one save file is not a degraded fleet, it is no fleet.
#    Isolation reads like housekeeping and is load bearing.
#
#    The sixth check tells us the cascade is honest. wait-interrupt survives,
#    and it is the one scenario that spawns no child at all, so the damage is
#    confined to children rather than being a grader that fell over. A sharper
#    mutant was attempted, creating the directory unwritable instead, and it
#    killed the same six: there is no way to take away a child's own storage
#    and leave it able to work.
run_mutant 2 "child shares the parent's data directory" \
	"spawn-reports-back child-isolated fleet-visibility join-fans-in addressed-messaging idle-corrected" \
	"agent/agent.go" \
	"perl -i -pe 's/c\.DataDir = dataDir/c.DataDir = s.DataDir/' agent/agent.go" \
	"grep -q 'c.DataDir = s.DataDir' agent/agent.go"

# 3. Shutting an agent down stops it but leaves it listed. A roster that still
#    names a dead agent is worse than no roster, because the next instruction
#    addressed to it goes nowhere and reports success.
run_mutant 3 "shutdown leaves the child on the roster" \
	"fleet-visibility" \
	"agent/subagent.go" \
	"perl -i -pe 's/delete\(m\.children, id\)/_ = id/' agent/subagent.go" \
	"! grep -q 'delete(m.children, id)' agent/subagent.go"

# 4. The join collects only the first agent it was given. One child reporting
#    looks exactly like a working fan in until a second child exists, which is
#    why the probe always joins two.
run_mutant 4 "join collects only the first agent" \
	"join-fans-in" \
	"agent/fleet.go" \
	"perl -i -pe 's/for _, id := range ids \{/for _, id := range ids[:1] {/' agent/fleet.go" \
	"grep -q 'range ids\[:1\]' agent/fleet.go"

# 5. The reply arrives with no sender. The answer is still delivered and the
#    child still resumes, so the only thing lost is the address, and the only
#    thing lost with it is the ability to ask a second question.
run_mutant 5 "reply carries no sender address" \
	"addressed-messaging" \
	"agent/messaging.go" \
	"perl -i -pe 's/m\.parent\.Spec\(\)\.ID,/\"\",/' agent/messaging.go" \
	"! grep -q 'From:.*m.parent.Spec().ID' agent/messaging.go"

# 6. A message arriving during a tool wait is attached to the conversation but
#    does not end the wait. This is the state the chapter started from: the
#    text is read, and nothing can act on it until the tool it was meant to
#    interrupt has already finished.
run_mutant 6 "an arriving message does not end the wait" \
	"wait-interrupt" \
	"agent/internal/llm/actor.go" \
	"perl -i -pe 's/a\.eng\.Jobs\.Interrupt\(\)/_ = a.eng.Jobs/' agent/internal/llm/actor.go" \
	"! grep -q 'a.eng.Jobs.Interrupt()' agent/internal/llm/actor.go"

# 7. An agent that stops without submitting is left alone. Nothing hangs and
#    nothing errors, which is precisely why this one is worth a check: the
#    fleet simply reports a child as unfinished forever.
run_mutant 7 "an idle child is never woken" \
	"idle-corrected" \
	"agent/idle.go" \
	"perl -i -pe 's/s\.Start\(wakeText\(n\)\)/_ = n/' agent/idle.go" \
	"! grep -q 's.Start(wakeText(n))' agent/idle.go"

echo
echo "================================"
echo "killed cleanly: $PASSED"
echo "problems:      $FAILED"
echo "================================"
[ "$FAILED" -eq 0 ]
