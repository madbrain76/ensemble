package agent

import (
	"strings"
	"testing"
)

// TestOwnershipEnforcementCrossTreeRejected verifies that agent A cannot
// shut down agent B's child. The SubAgentManager is per-agent, so cross-
// tree access is rejected structurally: B's child simply does not exist in
// A's manager.
func TestOwnershipEnforcementCrossTreeRejected(t *testing.T) {
	dir := t.TempDir()
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "done"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "ok"})
	})

	en, err := NewEnsemble(dir, cfg)
	if err != nil {
		t.Fatalf("NewEnsemble: %v", err)
	}
	defer en.Shutdown()

	agentA, err := en.CreateAgent("agent-a", "", "")
	if err != nil {
		t.Fatalf("CreateAgent A: %v", err)
	}
	agentB, err := en.CreateAgent("agent-b", "", "")
	if err != nil {
		t.Fatalf("CreateAgent B: %v", err)
	}

	// B spawns a child. The child will submit and finish quickly.
	_, err = agentB.subs.Spawn("b-child", "work", nil)
	if err != nil {
		t.Fatalf("B.Spawn: %v", err)
	}

	// A tries to shut down B's child. Must fail.
	err = agentA.subs.Shutdown("b-child")
	if err == nil {
		t.Fatal("agent A should not be able to shut down agent B's child")
	}
	if !strings.Contains(err.Error(), "no such sub-agent") {
		t.Fatalf("unexpected error: %v", err)
	}

	// A tries to send a message to B's child. Must fail.
	_, ok := agentA.subs.Get("b-child")
	if ok {
		t.Fatal("agent A should not be able to find agent B's child via Get")
	}

	// B can still reach its own child.
	_, ok = agentB.subs.Get("b-child")
	if !ok {
		t.Fatal("agent B should be able to find its own child")
	}
}

// TestOwnershipEnforcementJoinRejectsUnknown verifies that join_agents
// reports "unknown_agent" for IDs the caller did not spawn.
func TestOwnershipEnforcementJoinRejectsUnknown(t *testing.T) {
	dir := t.TempDir()
	cfg := bareConfig(t)

	en, err := NewEnsemble(dir, cfg)
	if err != nil {
		t.Fatalf("NewEnsemble: %v", err)
	}
	defer en.Shutdown()

	agentA, err := en.CreateAgent("agent-a", "", "")
	if err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}

	// Join with an ID that A never spawned.
	result := agentA.subs.Join([]string{"nonexistent"}, 1)
	if len(result.Outcomes) != 1 {
		t.Fatalf("expected 1 outcome, got %d", len(result.Outcomes))
	}
	if result.Outcomes[0].Outcome != "unknown_agent" {
		t.Fatalf("expected outcome 'unknown_agent', got %q", result.Outcomes[0].Outcome)
	}
}
