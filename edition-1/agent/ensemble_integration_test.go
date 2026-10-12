package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/waywardgeek/ensemble/agent/internal/common"
)

// TestEnsembleLifecycle exercises the full create → observe → persist → reload
// → delete lifecycle in a single test. This is the integration test for ch26:
// it proves the composition root, activity tracking, persistence, and tree
// snapshot work end-to-end.
func TestEnsembleLifecycle(t *testing.T) {
	dir := t.TempDir()
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		// Just end the turn with text.
		return []map[string]any{{"type": "text", "text": "ok"}}
	})

	// ── Phase 1: Create an agent and observe the tree ──
	en, err := NewEnsemble(dir, cfg, AgentSpec{DataDir: dir})
	if err != nil {
		t.Fatal(err)
	}

	// Tree starts with just the root.
	tree := en.TreeSnapshot()
	if tree.ID != "main" {
		t.Fatalf("root id = %q, want main", tree.ID)
	}
	if len(tree.Children) != 0 {
		t.Fatalf("root children = %d, want 0", len(tree.Children))
	}

	// Create a named agent.
	a, err := en.CreateAgent("worker-1", "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Tree now has the child.
	tree = en.TreeSnapshot()
	if len(tree.Children) != 1 {
		t.Fatalf("children = %d after create, want 1", len(tree.Children))
	}
	if tree.Children[0].ID != "worker-1" {
		t.Errorf("child id = %q, want worker-1", tree.Children[0].ID)
	}

	// ── Phase 2: Activity tracking ──
	var activities []common.ActivityPhase
	en.Root().Actor().Notify(common.ActivityChanged{Phase: common.ActivityThinking})
	act := en.Root().Activity()
	if act.Phase != common.ActivityThinking {
		t.Errorf("root activity = %v, want thinking", act.Phase)
	}

	// Reset to idle.
	en.Root().Actor().Notify(common.ActivityChanged{Phase: common.ActivityIdle})
	act = en.Root().Activity()
	if act.Phase != common.ActivityIdle {
		t.Errorf("root activity = %v, want idle", act.Phase)
	}

	// Collect tree change events.
	var treeChanged int
	en.SetOnTreeChange(func() { treeChanged++ })
	en.Root().Actor().Notify(common.ActivityChanged{Phase: common.ActivityToolCall, Detail: "read_file"})
	if treeChanged < 1 {
		t.Error("expected tree change callback on activity transition")
	}

	// ── Phase 3: Persistence ──
	// Verify ensemble.json was written.
	data, err := os.ReadFile(filepath.Join(dir, "ensemble.json"))
	if err != nil {
		t.Fatal(err)
	}
	var reg []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &reg); err != nil {
		t.Fatal(err)
	}
	if len(reg) != 1 || reg[0].Name != "worker-1" {
		t.Errorf("registry = %s, want [{name: worker-1}]", data)
	}

	// Shut down agent and reopen.
	a.Shutdown()
	en.Shutdown()

	en2, err := NewEnsemble(dir, cfg, AgentSpec{DataDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	if err := en2.Start(); err != nil {
		t.Fatal(err)
	}
	defer en2.Shutdown()

	// Agent survives restart.
	tree2 := en2.TreeSnapshot()
	if len(tree2.Children) != 1 {
		t.Fatalf("children after reload = %d, want 1", len(tree2.Children))
	}
	if tree2.Children[0].ID != "worker-1" {
		t.Errorf("child id after reload = %q, want worker-1", tree2.Children[0].ID)
	}

	// ── Phase 4: Delete ──
	if err := en2.DeleteAgent("worker-1"); err != nil {
		t.Fatal(err)
	}
	tree3 := en2.TreeSnapshot()
	if len(tree3.Children) != 0 {
		t.Fatalf("children after delete = %d, want 0", len(tree3.Children))
	}

	// Registry is updated.
	data, _ = os.ReadFile(filepath.Join(dir, "ensemble.json"))
	var reg2 []struct {
		Name string `json:"name"`
	}
	json.Unmarshal(data, &reg2)
	if len(reg2) != 0 {
		t.Errorf("registry after delete = %s, want empty", data)
	}

	_ = activities // suppress unused warning
}
