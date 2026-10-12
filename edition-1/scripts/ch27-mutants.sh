#!/usr/bin/env bash
#
# Chapter 27 mutation audit.
#
# A passing grader is evidence of nothing until each check has been shown to
# fail when the behavior it claims to protect is removed. Each mutant below
# deletes exactly one behavior, must still compile, and must kill a named
# set of checks. A survivor means the check is decorative, the mutant is
# invalid, or the harness is broken, and all three are worth knowing.
set -u

cd "$(dirname "$0")/.."
ROOT="$(pwd)"
OUT="$(mktemp -d)"
trap 'rm -rf "$OUT"' EXIT

PASSED=0
FAILED=0

echo "building the grader..."
GRADER="$OUT/grade"
if ! go build -o "$GRADER" ./cmd/grade; then
	echo "FATAL: the grader does not build"
	exit 1
fi

# Baseline: the unmutated tree must score 100. Proving the harness can
# report a kill is worthless if it cannot first report a pass.
echo "baseline..."
"$GRADER" -ch 27 -json ./agent >"$OUT/baseline.json" 2>"$OUT/baseline.log"
BASE_FAILS=$(python3 - "$OUT/baseline.json" <<'PY'
import json, sys
def walk(x, out):
    if isinstance(x, dict):
        if "id" in x and "passed" in x and not x["passed"]:
            out.append(x["id"])
        for v in x.values(): walk(v, out)
    elif isinstance(x, list):
        for v in x: walk(v, out)
out = []
try: walk(json.load(open(sys.argv[1])), out)
except Exception as e: print("PARSE-ERROR", e); sys.exit(0)
print(" ".join(sorted(set(out))))
PY
)
if [ -n "$BASE_FAILS" ]; then
	echo "FATAL: baseline is not clean, these already fail: $BASE_FAILS"
	echo "       (fix the tree before auditing; a dirty baseline makes every"
	echo "        mutant look like it killed something)"
	exit 1
fi
echo "  baseline clean."

# run_mutant <n> <description> <expected-failing-ids> <files> <mutate-cmd> <verify-cmd>
run_mutant() {
	local n="$1" desc="$2" expect="$3" files="$4" mutate="$5" verify="$6"

	# Optional first argument selects a single mutant, so a fix can be
	# retested in seconds instead of re-running the whole audit.
	if [ -n "${ONLY:-}" ] && [ "${ONLY}" != "$n" ]; then
		return
	fi

	echo ""
	echo "=== mutant $n: $desc"
	git checkout -- $files >/dev/null 2>&1

	if ! eval "$mutate"; then
		echo "  MUTATION COMMAND FAILED - script bug"
		git checkout -- $files
		FAILED=$((FAILED + 1))
		return
	fi
	# A mutation that silently did not apply produces a false survivor,
	# which reads exactly like an insensitive check.
	if ! eval "$verify" >/dev/null 2>&1; then
		echo "  MUTATION DID NOT APPLY - script bug"
		git checkout -- $files
		FAILED=$((FAILED + 1))
		return
	fi
	echo "  mutation landed."

	if ! (cd agent && go build ./... >/dev/null 2>&1); then
		echo "  MUTANT DOES NOT COMPILE - invalid mutant (script bug)"
		git checkout -- $files
		FAILED=$((FAILED + 1))
		return
	fi
	echo "  mutant compiles."

	local report="$OUT/report-$n.json"
	"$GRADER" -ch 27 -json ./agent >"$report" 2>"$OUT/stderr-$n.log"

	local killed
	killed=$(python3 - "$report" <<'PY'
import json, sys
def walk(x, out):
    if isinstance(x, dict):
        if "id" in x and "passed" in x and not x["passed"]:
            out.append(x["id"])
        for v in x.values(): walk(v, out)
    elif isinstance(x, list):
        for v in x: walk(v, out)
out = []
try: walk(json.load(open(sys.argv[1])), out)
except Exception as e: print("PARSE-ERROR", e); sys.exit(0)
print(" ".join(sorted(set(out))))
PY
	)

	local want got
	want=$(echo "$expect" | tr ' ' '\n' | sort -u | tr '\n' ' ' | sed 's/ $//')
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

# 1. Both handshake error paths are swallowed, so a peer that fails to
#    authenticate is reported as one that succeeded. Only the failure paths
#    change, so a legitimate connection is untouched and the damage lands
#    exactly on the check that presents a wrong static key.
#
#    Swallowing only the read error is NOT enough: the crypto fails again at
#    the next step, so the rejection still happens and the mutant survives.
#    That survivor was an invalid mutant rather than a blind check.
run_mutant 1 "handshake errors are swallowed, so bad keys are accepted" \
	"C1" \
	"agent/gateway/noise.go" \
	"perl -i -pe 's/return nil, fmt\.Errorf\(\"failed to (read|write) handshake message: %v\", err\)/return nil, nil/' agent/gateway/noise.go" \
	"! grep -q 'handshake message: %v' agent/gateway/noise.go"

# 2. The provenance header is no longer prepended at delivery, so the agent
#    receives content with no indication of who sent it or at what trust.
run_mutant 2 "provenance header is not prepended on delivery" \
	"C2" \
	"agent/internal/gateway/server.go" \
	"perl -i -pe 's/delivered\.Content = env\.WithProvenance\(\)/delivered.Content = env.Content/' agent/internal/gateway/server.go" \
	"grep -q 'delivered.Content = env.Content' agent/internal/gateway/server.go"

# 3. The welcome frame no longer identifies the connection it belongs to, so
#    an agent cannot confirm which identity the gateway authenticated it as.
run_mutant 3 "welcome frame reports the wrong agent id" \
	"C3" \
	"agent/internal/gateway/server.go" \
	"perl -i -pe 's/AgentID:       c\.agentID,/AgentID:       \"unknown\",/' agent/internal/gateway/server.go" \
	"grep -q 'AgentID:       \"unknown\"' agent/internal/gateway/server.go"

# 4. Authorization always succeeds, so contactability stops being a
#    restriction and any agent can reach any other.
run_mutant 4 "Authorize never refuses" \
	"C4" \
	"agent/internal/gateway/groups.go" \
	"perl -i -pe 's/return nil, \"\", fmt\.Errorf\(\"%s is not in any group %s accepts contact from\", sender, recipient\)/return []string{\"agents\"}, pub.TrustTrusted, nil/' agent/internal/gateway/groups.go" \
	"grep -q 'return \[\]string{\"agents\"}, pub.TrustTrusted, nil' agent/internal/gateway/groups.go"

# 5. Resume ignores the sequence the agent says it already handled, so a
#    reconnect replays messages the agent has durably processed.
run_mutant 5 "resume ignores the client's last-acked sequence" \
	"C5" \
	"agent/internal/gateway/queue.go" \
	"perl -i -pe 's/q\.ackLocked\(agentID, lastAcked\)/q.ackLocked(agentID, 0)/' agent/internal/gateway/queue.go" \
	"grep -q 'q.ackLocked(agentID, 0)' agent/internal/gateway/queue.go"

# 6. The public client library reaches into an internal package, which makes
#    it uncompilable from outside the module no matter what it exports.
run_mutant 6 "public client library imports an internal package" \
	"C6" \
	"agent/gateway/protocol.go" \
	"perl -0777 -i -pe 's|import \"encoding/json\"|import (\n\t\"encoding/json\"\n\n\t_ \"github.com/waywardgeek/ensemble/agent/internal/common\"\n)|' agent/gateway/protocol.go" \
	"grep -q 'agent/internal/common' agent/gateway/protocol.go"

# 7. The narrow-only rule is removed from child registration, so a parent can
#    grant a child a group the parent does not itself hold.
run_mutant 7 "child registration skips the narrow-only check" \
	"C7" \
	"agent/internal/gateway/groups.go" \
	"perl -i -pe 's/if err := r\.checkGrantableLocked\(parentID, groups\); err != nil \{/if err := error(nil); err != nil {/' agent/internal/gateway/groups.go" \
	"grep -q 'if err := error(nil); err != nil {' agent/internal/gateway/groups.go"

echo ""
echo "=================================="
echo "mutants killed as expected: $PASSED"
echo "problems:                   $FAILED"
echo "=================================="

git status --porcelain agent/ | head

[ "$FAILED" -eq 0 ]
