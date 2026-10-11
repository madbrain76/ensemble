package agent

import (
	"testing"
	"time"

	"github.com/waywardgeek/ensemble/agent/internal/common"
	"github.com/waywardgeek/ensemble/agent/internal/llm"
)

// TestActivityTransitions verifies that a scripted turn emits at least
// thinking, tool_call, and idle activity observations in order.
//
// The turn: the fake vendor returns a tool call, the agent dispatches it,
// then on the next round the vendor returns plain text and the turn ends.
// The observer should see: thinking → tool_call → thinking → idle.
func TestActivityTransitions(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "all done"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "ok"})
	})

	agent, err := NewAgent(cfg, AgentSpec{DataDir: t.TempDir(), ID: "activity-test"})
	if err != nil {
		t.Fatalf("NewAgent: %v", err)
	}
	t.Cleanup(func() { _ = agent.Shutdown() })

	// Collect activity observations.
	var phases []common.ActivityPhase
	done := make(chan struct{})
	detach := agent.Actor().Attach(llm.ObserverFunc(func(obs common.Observation) {
		if ac, ok := obs.(common.ActivityChanged); ok {
			phases = append(phases, ac.Phase)
			if ac.Phase == common.ActivityIdle {
				select {
				case <-done:
				default:
					close(done)
				}
			}
		}
	}))
	defer detach()

	_, err = agent.Ask("do something")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}

	// Wait for the idle observation.
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for idle activity")
	}

	// Verify we saw the required transitions.
	if len(phases) < 3 {
		t.Fatalf("expected at least 3 activity transitions, got %d: %v", len(phases), phases)
	}

	// Must see thinking before tool_call, and tool_call before the final idle.
	sawThinking := false
	sawToolCall := false
	sawToolCallAfterThinking := false
	sawIdleAfterToolCall := false
	for _, p := range phases {
		switch p {
		case common.ActivityThinking:
			sawThinking = true
		case common.ActivityToolCall:
			sawToolCall = true
			if sawThinking {
				sawToolCallAfterThinking = true
			}
		case common.ActivityIdle:
			if sawToolCall {
				sawIdleAfterToolCall = true
			}
		}
	}
	if !sawThinking {
		t.Error("never saw ActivityThinking")
	}
	if !sawToolCallAfterThinking {
		t.Error("never saw ActivityToolCall after ActivityThinking")
	}
	if !sawIdleAfterToolCall {
		t.Error("never saw ActivityIdle after ActivityToolCall")
	}
}

// TestActivityToolCallIncludesName verifies that the tool_call activity
// carries the name of the tool being dispatched.
func TestActivityToolCallIncludesName(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "done"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "ok"})
	})

	agent, err := NewAgent(cfg, AgentSpec{DataDir: t.TempDir(), ID: "tool-name-test"})
	if err != nil {
		t.Fatalf("NewAgent: %v", err)
	}
	t.Cleanup(func() { _ = agent.Shutdown() })

	var toolName string
	detach := agent.Actor().Attach(llm.ObserverFunc(func(obs common.Observation) {
		if ac, ok := obs.(common.ActivityChanged); ok && ac.Phase == common.ActivityToolCall {
			toolName = ac.ToolName
		}
	}))
	defer detach()

	_, err = agent.Ask("run submit")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}

	if toolName != "submit" {
		t.Errorf("expected tool name 'submit', got %q", toolName)
	}
}
