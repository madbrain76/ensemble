package agent

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/waywardgeek/ensemble/agent/internal/common"
)

// TestAskAgent verifies that a parent can synchronously ask a child a
// question and receive the response text.
func TestAskAgent(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		// The child just responds with text. When it sees a submit tool
		// result, it responds with its final text.
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "done"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "42"})
	})
	parent := fleetParent(t, cfg)

	sub, err := parent.subs.Spawn("worker", "compute", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	// Start the first turn manually.
	if err := sub.Start("compute"); err != nil {
		t.Fatalf("start: %v", err)
	}
	// Wait for the first turn to complete.
	parent.subs.Join([]string{"worker"}, 5*time.Second)

	// Now ask_agent as a tool call. This calls Agent.Ask internally.
	result, err := parent.toolAskAgent(nil, json.RawMessage(`{"agent_id":"worker","question":"what was the answer?"}`))
	if err != nil {
		t.Fatalf("ask_agent: %v", err)
	}
	// The fake always responds with text. The exact text depends on
	// whether it's the first or subsequent turn.
	if result == "" {
		t.Error("ask_agent returned empty result")
	}
}

// TestGetSubAgentStatus verifies the detailed status tool returns
// state, task, and turns.
func TestGetSubAgentStatus(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "done"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "ok"})
	})
	parent := fleetParent(t, cfg)

	sub, err := parent.subs.Spawn("scout", "investigate threats", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("investigate threats"); err != nil {
		t.Fatalf("start: %v", err)
	}
	parent.subs.Join([]string{"scout"}, 5*time.Second)

	result, err := parent.toolGetSubAgentStatus(nil, json.RawMessage(`{"agent_id":"scout"}`))
	if err != nil {
		t.Fatalf("get_sub_agent_status: %v", err)
	}
	if !strings.Contains(result, "Agent: scout") {
		t.Errorf("status missing agent name, got: %s", result)
	}
	if !strings.Contains(result, "Task: investigate threats") {
		t.Errorf("status missing task, got: %s", result)
	}
	if !strings.Contains(result, "State: done") {
		t.Errorf("status should show done, got: %s", result)
	}
}

// TestGetSubmission verifies the submission payload is retrievable.
func TestGetSubmission(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "done"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "found 3 vulnerabilities"})
	})
	parent := fleetParent(t, cfg)

	sub, err := parent.subs.Spawn("scanner", "scan network", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("scan network"); err != nil {
		t.Fatalf("start: %v", err)
	}
	parent.subs.Join([]string{"scanner"}, 5*time.Second)

	result, err := parent.toolGetSubmission(nil, json.RawMessage(`{"agent_id":"scanner"}`))
	if err != nil {
		t.Fatalf("get_submission: %v", err)
	}
	if !strings.Contains(result, "vulnerabilities") {
		t.Errorf("submission should contain the child's payload, got: %s", result)
	}
}

// TestGetSubmission_NotSubmitted verifies the error when a child hasn't submitted.
func TestGetSubmission_NotSubmitted(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		// Child just talks and stops, never submits.
		return []map[string]any{{"type": "text", "text": "I'm thinking"}}
	})
	parent := fleetParent(t, cfg)

	_, err := parent.subs.Spawn("lazy", "do work", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	_, err = parent.toolGetSubmission(nil, json.RawMessage(`{"agent_id":"lazy"}`))
	if err == nil {
		t.Error("get_submission should fail when agent has not submitted")
	}
	if err != nil && !strings.Contains(err.Error(), "has not submitted") {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestSendMessage verifies that a message can be delivered to a child.
func TestSendMessage(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		return []map[string]any{{"type": "text", "text": "acknowledged"}}
	})
	parent := fleetParent(t, cfg)

	_, err := parent.subs.Spawn("worker", "wait for orders", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	result, err := parent.toolSendMessage(nil, json.RawMessage(`{"agent_id":"worker","prompt":"go"}`))
	if err != nil {
		t.Fatalf("send_message: %v", err)
	}
	if !strings.Contains(result, "delivered") {
		t.Errorf("expected delivery confirmation, got: %s", result)
	}
}

// TestInterruptAgent verifies that interrupting a child returns a
// confirmation and the child is still alive.
func TestInterruptAgent(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		return []map[string]any{{"type": "text", "text": "working"}}
	})
	parent := fleetParent(t, cfg)

	_, err := parent.subs.Spawn("worker", "long task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	result, err := parent.toolInterruptAgent(nil, json.RawMessage(`{"agent_id":"worker"}`))
	if err != nil {
		t.Fatalf("interrupt_agent: %v", err)
	}
	if !strings.Contains(result, "interrupted") {
		t.Errorf("expected interrupt confirmation, got: %s", result)
	}
	// Agent should still be in the fleet.
	if _, ok := parent.subs.Get("worker"); !ok {
		t.Error("interrupted agent was removed from fleet")
	}
}

// TestKillAgent verifies that killing a child removes it from the fleet.
func TestKillAgent(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		return []map[string]any{{"type": "text", "text": "working"}}
	})
	parent := fleetParent(t, cfg)

	_, err := parent.subs.Spawn("doomed", "dangerous task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	result, err := parent.toolKillAgent(nil, json.RawMessage(`{"agent_id":"doomed"}`))
	if err != nil {
		t.Fatalf("kill_agent: %v", err)
	}
	if !strings.Contains(result, "killed") {
		t.Errorf("expected kill confirmation, got: %s", result)
	}
	// Agent should be gone from the fleet.
	if _, ok := parent.subs.Get("doomed"); ok {
		t.Error("killed agent still present in fleet")
	}
}

// TestUnknownAgentRejected verifies that management tools return errors
// for agents the caller didn't spawn.
func TestUnknownAgentRejected(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		return []map[string]any{{"type": "text", "text": "ok"}}
	})
	parent := fleetParent(t, cfg)

	tools := []struct {
		name string
		call func() error
	}{
		{"send_message", func() error {
			_, err := parent.toolSendMessage(nil, json.RawMessage(`{"agent_id":"ghost","prompt":"hi"}`))
			return err
		}},
		{"ask_agent", func() error {
			_, err := parent.toolAskAgent(nil, json.RawMessage(`{"agent_id":"ghost","question":"hi"}`))
			return err
		}},
		{"get_submission", func() error {
			_, err := parent.toolGetSubmission(nil, json.RawMessage(`{"agent_id":"ghost"}`))
			return err
		}},
		{"get_sub_agent_status", func() error {
			_, err := parent.toolGetSubAgentStatus(nil, json.RawMessage(`{"agent_id":"ghost"}`))
			return err
		}},
		{"kill_agent", func() error {
			_, err := parent.toolKillAgent(nil, json.RawMessage(`{"agent_id":"ghost"}`))
			return err
		}},
		{"interrupt_agent", func() error {
			_, err := parent.toolInterruptAgent(nil, json.RawMessage(`{"agent_id":"ghost"}`))
			return err
		}},
	}

	for _, tc := range tools {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Errorf("%s should reject unknown agent", tc.name)
			}
			if err != nil && !strings.Contains(err.Error(), "unknown agent") {
				t.Errorf("%s: expected 'unknown agent' error, got: %v", tc.name, err)
			}
		})
	}
}

// TestWaitForAgentChange_StatusChange verifies the epoll/select mechanism
// wakes when a child's state changes.
func TestWaitForAgentChange_StatusChange(t *testing.T) {
	// The child will respond with text on the second request (after tool result).
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "done"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "ok"})
	})
	parent := fleetParent(t, cfg)

	// Spawn but don't start yet. Start will happen while wait is blocking.
	sub, err := parent.subs.Spawn("target", "work", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	// Start the child in a goroutine so the main thread can call wait.
	started := make(chan struct{})
	go func() {
		close(started)
		_ = sub.Start("work")
	}()
	<-started

	result, err := parent.toolWaitForAgentChange(nil, json.RawMessage(
		`{"agent_ids":["target"],"timeout_seconds":5}`))
	if err != nil {
		t.Fatalf("wait_for_agent_change: %v", err)
	}

	var resp struct {
		Events []agentEvent `json:"events"`
	}
	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		t.Fatalf("parse events: %v", err)
	}
	if len(resp.Events) == 0 {
		t.Fatal("wait_for_agent_change returned no events")
	}
	found := false
	for _, ev := range resp.Events {
		if ev.Kind == "status_change" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected status_change event, got: %+v", resp.Events)
	}
}

// TestActivityTracking verifies that the Agent.Activity() method tracks
// activity changes emitted by the actor.
func TestActivityTracking(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		return []map[string]any{{"type": "text", "text": "done"}}
	})
	a, err := NewAgent(cfg, AgentSpec{DataDir: t.TempDir(), ID: "test"})
	if err != nil {
		t.Fatalf("NewAgent: %v", err)
	}
	t.Cleanup(func() { _ = a.Shutdown() })

	// Initially idle.
	act := a.Activity()
	if act.Phase != common.ActivityIdle && act.Phase != 0 {
		t.Errorf("initial activity = %v, want idle or zero", act.Phase)
	}

	// Run a turn and check that activity changed.
	reply, err := a.Ask("hello")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}
	_ = reply

	// After the turn, activity should be back to idle.
	act = a.Activity()
	if act.Phase != common.ActivityIdle {
		t.Errorf("post-turn activity = %v, want idle", act.Phase)
	}
}
