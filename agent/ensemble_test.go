package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestEnsembleCreateAndReload verifies that creating an agent persists it
// to ensemble.json and that a fresh Ensemble loaded from the same directory
// finds the agent.
func TestEnsembleCreateAndReload(t *testing.T) {
	dir := t.TempDir()
	cfg := bareConfig(t)

	// Phase 1: create an agent and shut down.
	en, err := NewEnsemble(dir, cfg, AgentSpec{DataDir: dir})
	if err != nil {
		t.Fatalf("NewEnsemble: %v", err)
	}
	agent, err := en.CreateAgent("analyst", "", "")
	if err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}
	if agent == nil {
		t.Fatal("CreateAgent returned nil agent")
	}
	if en.Agent("analyst") == nil {
		t.Fatal("Agent lookup returned nil immediately after create")
	}
	_ = en.Shutdown()

	// Phase 2: reload from the same directory.
	en2, err := NewEnsemble(dir, cfg, AgentSpec{DataDir: dir})
	if err != nil {
		t.Fatalf("NewEnsemble (reload): %v", err)
	}
	if err := en2.Start(); err != nil {
		t.Fatalf("Start (reload): %v", err)
	}
	defer en2.Shutdown()

	if en2.Agent("analyst") == nil {
		t.Fatal("agent 'analyst' not found after reload")
	}

	// Verify the registry file on disk.
	data, err := os.ReadFile(filepath.Join(dir, "ensemble.json"))
	if err != nil {
		t.Fatalf("read ensemble.json: %v", err)
	}
	var entries []AgentEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("parse ensemble.json: %v", err)
	}
	if len(entries) != 1 || entries[0].Name != "analyst" {
		t.Fatalf("unexpected registry: %+v", entries)
	}
}

// TestEnsembleDelete verifies that deleting an agent removes it from the
// live map, the persistent registry, and cleans up its data directory.
func TestEnsembleDelete(t *testing.T) {
	dir := t.TempDir()
	cfg := bareConfig(t)

	en, err := NewEnsemble(dir, cfg, AgentSpec{DataDir: dir})
	if err != nil {
		t.Fatalf("NewEnsemble: %v", err)
	}
	defer en.Shutdown()

	if _, err := en.CreateAgent("doomed", "", ""); err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}
	if en.Agent("doomed") == nil {
		t.Fatal("agent not created")
	}

	if err := en.DeleteAgent("doomed"); err != nil {
		t.Fatalf("DeleteAgent: %v", err)
	}

	if en.Agent("doomed") != nil {
		t.Fatal("agent still in live map after delete")
	}

	// Verify the data directory was removed.
	if _, err := os.Stat(filepath.Join(dir, "doomed")); !os.IsNotExist(err) {
		t.Fatalf("data directory still exists after delete: %v", err)
	}

	// Verify the registry is empty.
	data, err := os.ReadFile(filepath.Join(dir, "ensemble.json"))
	if err != nil {
		t.Fatalf("read ensemble.json: %v", err)
	}
	var entries []AgentEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("parse ensemble.json: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty registry, got %+v", entries)
	}
}

// TestEnsembleDeleteRejectsWithActiveChildren verifies that DeleteAgent
// refuses when the agent has active sub-agents.
func TestEnsembleDeleteRejectsWithActiveChildren(t *testing.T) {
	dir := t.TempDir()

	// The parent's vendor returns a tool call so the turn stays running.
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		return toolUse("t1", "submit", map[string]any{"result": "busy"})
	})

	en, err := NewEnsemble(dir, cfg, AgentSpec{DataDir: dir})
	if err != nil {
		t.Fatalf("NewEnsemble: %v", err)
	}
	defer en.Shutdown()

	parent, err := en.CreateAgent("parent", "", "")
	if err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}

	// Spawn a child directly through the SubAgentManager.
	child, err := parent.subs.Spawn("child-1", "work", nil)
	if err != nil {
		t.Fatalf("Spawn: %v", err)
	}
	if err := child.Start("work"); err != nil {
		t.Fatalf("child.Start: %v", err)
	}

	// DeleteAgent should reject while the child is active.
	err = en.DeleteAgent("parent")
	if err == nil {
		t.Fatal("expected DeleteAgent to reject with active children")
	}
}

// TestEnsembleDuplicateRejected verifies that creating two agents with the
// same name fails.
func TestEnsembleDuplicateRejected(t *testing.T) {
	dir := t.TempDir()
	cfg := bareConfig(t)

	en, err := NewEnsemble(dir, cfg, AgentSpec{DataDir: dir})
	if err != nil {
		t.Fatalf("NewEnsemble: %v", err)
	}
	defer en.Shutdown()

	if _, err := en.CreateAgent("one", "", ""); err != nil {
		t.Fatalf("first CreateAgent: %v", err)
	}
	if _, err := en.CreateAgent("one", "", ""); err == nil {
		t.Fatal("expected duplicate to be rejected")
	}
}

// bareConfig returns a minimal Config that does not need a real vendor.
// It uses the same fake vendor as fleetConfig but with a simple response.
func bareConfig(t *testing.T) Config {
	t.Helper()
	return fleetConfig(t, func(req map[string]any) []map[string]any {
		return []map[string]any{{"type": "text", "text": "ok"}}
	})
}
