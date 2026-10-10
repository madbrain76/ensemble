package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/waywardgeek/ensemble/agent/internal/common"
)

// fleetConfig returns a Config pointed at a fake vendor.
//
// A fake server rather than a mocked engine: the child has to run a REAL
// turn through the real actor for the result to mean anything. A test that
// stubbed the turn would prove the manager's bookkeeping and nothing about
// whether a child can actually work.
func fleetConfig(t *testing.T, reply func(map[string]any) []map[string]any) Config {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "test", "type": "message", "role": "assistant",
			"model": "claude-sonnet-5", "content": reply(req),
			"stop_reason": "end_turn",
			"usage":       map[string]int{"input_tokens": 1, "output_tokens": 1},
		})
	}))
	t.Cleanup(srv.Close)
	return Config{
		Model: "claude-sonnet-5", Vendor: common.VendorAnthropic,
		BaseURL: srv.URL, APIKey: "test", DisableStreaming: true,
	}
}

// sawToolResult reports whether the conversation already contains a tool
// result, which is how these fakes tell a first turn from a later one
// without counting calls. Counting breaks the moment a retry happens.
func sawToolResult(req map[string]any) bool {
	raw, _ := json.Marshal(req)
	return strings.Contains(string(raw), "tool_result")
}

func toolUse(id, name string, input map[string]any) []map[string]any {
	return []map[string]any{{"type": "tool_use", "id": id, "name": name, "input": input}}
}

// fleetParent builds a parent wired to a factory that builds real children
// against the same fake vendor.
func fleetParent(t *testing.T, cfg Config) *Agent {
	t.Helper()
	parent, err := NewAgent(cfg, AgentSpec{DataDir: t.TempDir(), ID: "boss"})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })
	parent.SetAgentFactory(func(spec AgentSpec, task string) (*Agent, error) {
		return NewAgent(cfg, spec)
	})
	return parent
}

// TestChildRunsAndReportsBack is check 1: a spawned child runs a real turn,
// calls submit, and its payload reaches the parent.
func TestChildRunsAndReportsBack(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "finished"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "the answer is 42"})
	})
	parent := fleetParent(t, cfg)

	sub, err := parent.SubAgents().Spawn("worker", "find the answer", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("find the answer"); err != nil {
		t.Fatalf("start: %v", err)
	}

	res := parent.SubAgents().Join([]string{"worker"}, 10*time.Second)
	if !res.AllFinished {
		t.Fatalf("join did not finish: %+v", res.Outcomes)
	}
	if len(res.Outcomes) != 1 {
		t.Fatalf("got %d outcomes, want 1", len(res.Outcomes))
	}
	if res.Outcomes[0].Outcome != "submitted" {
		t.Errorf("outcome = %q, want \"submitted\"", res.Outcomes[0].Outcome)
	}
	if !strings.Contains(res.Outcomes[0].Payload, "42") {
		t.Errorf("payload = %q, want the child's result", res.Outcomes[0].Payload)
	}
	// The child did real work, so it must have run at least one turn. A
	// child that never ran and one that ran and finished both end up idle,
	// and that is the confusion this guards.
	if sub.Turns() == 0 {
		t.Error("child reports zero turns yet produced a result")
	}
}

// TestJoinReportsIdleWithoutSubmitAsUnfinished is the heart of the chapter:
// an agent that stops without submitting has NOT finished.
//
// If join treated idle as finished, the parent would receive an empty result
// and carry on as though the work were done, which is the silent failure the
// whole chapter is about.
func TestJoinReportsIdleWithoutSubmitAsUnfinished(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		// The child just talks and stops. It never submits.
		return []map[string]any{{"type": "text", "text": "I think I am done"}}
	})
	parent := fleetParent(t, cfg)

	sub, err := parent.SubAgents().Spawn("slacker", "do the job", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("do the job"); err != nil {
		t.Fatalf("start: %v", err)
	}

	res := parent.SubAgents().Join([]string{"slacker"}, 2*time.Second)
	if res.AllFinished {
		t.Error("join reported a fleet finished when a child never submitted")
	}
	if got := res.Outcomes[0].Outcome; got != "idle_without_submit" {
		t.Errorf("outcome = %q, want \"idle_without_submit\"", got)
	}
	if res.Outcomes[0].Payload != "" {
		t.Error("a child that never submitted produced a payload")
	}
}

// TestJoinReturnsEarlyWhenChildAsks is check 5's deadlock guard: the parent
// gets control back so it can answer.
func TestJoinReturnsEarlyWhenChildAsks(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return toolUse("t2", "submit", map[string]any{"result": "done after asking"})
		}
		return toolUse("t1", "send_message_to_parent", map[string]any{
			"content": "which branch?", "timeout_seconds": 30,
		})
	})
	parent := fleetParent(t, cfg)

	sub, err := parent.SubAgents().Spawn("asker", "do the job", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("do the job"); err != nil {
		t.Fatalf("start: %v", err)
	}

	// Without the early return this blocks until the timeout: the child
	// waits for an answer, the parent waits for the child.
	res := parent.SubAgents().Join([]string{"asker"}, 10*time.Second)
	if !res.Interrupted {
		t.Fatalf("join did not return early for a question: %+v", res.Outcomes)
	}
	if len(res.Messages) != 1 {
		t.Fatalf("got %d messages, want 1", len(res.Messages))
	}
	q := res.Messages[0]
	if q.From != "asker" {
		t.Errorf("question From = %q, want \"asker\"", q.From)
	}
	// The rendered form must name the sender and the message, or the model
	// cannot call respond_to_parent_message with the right arguments.
	rendered := res.Render()
	if !strings.Contains(rendered, q.ID) || !strings.Contains(rendered, "asker") {
		t.Errorf("rendered join result lacks the ids needed to reply:\n%s", rendered)
	}

	if err := parent.SubAgents().Respond("asker", q.ID, "use main"); err != nil {
		t.Fatalf("respond: %v", err)
	}

	res = parent.SubAgents().Join([]string{"asker"}, 10*time.Second)
	if !res.AllFinished {
		t.Fatalf("child did not finish after being answered: %+v", res.Outcomes)
	}
	if !strings.Contains(res.Outcomes[0].Payload, "done after asking") {
		t.Errorf("payload = %q", res.Outcomes[0].Payload)
	}
}

// TestJoinUnknownAgentIsNotSuccess guards a quiet way to report success for
// work nobody did: a misspelled id must not look like a finished child.
func TestJoinUnknownAgentIsNotSuccess(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		return []map[string]any{{"type": "text", "text": "hi"}}
	})
	parent := fleetParent(t, cfg)

	res := parent.SubAgents().Join([]string{"never-spawned"}, 500*time.Millisecond)
	if res.AllFinished {
		t.Error("join reported success for an agent that was never spawned")
	}
	if got := res.Outcomes[0].Outcome; got != "unknown_agent" {
		t.Errorf("outcome = %q, want \"unknown_agent\"", got)
	}
}

// TestAgentsStatusListsChildren is check 3.
func TestAgentsStatusListsChildren(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any {
		if sawToolResult(req) {
			return []map[string]any{{"type": "text", "text": "ok"}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "r"})
	})
	parent := fleetParent(t, cfg)

	if _, err := parent.SubAgents().Spawn("alpha", "task a", nil); err != nil {
		t.Fatalf("spawn: %v", err)
	}
	status := parent.SubAgents().Status()
	if !strings.Contains(status, "alpha") {
		t.Errorf("status does not list the child:\n%s", status)
	}

	if err := parent.SubAgents().Shutdown("alpha"); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if status := parent.SubAgents().Status(); strings.Contains(status, "alpha") {
		t.Errorf("status still lists a shut-down child:\n%s", status)
	}
}

// TestWorkflowToolsAreRegistered pins the seven names the grader enforces.
//
// Spelling is a contract here: a model calls these by name, and a rename
// that compiles cleanly breaks every conversation that learned the old one.
func TestWorkflowToolsAreRegistered(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any { return nil })
	parent := fleetParent(t, cfg)

	want := []string{
		"spawn_sub_agent", "join_agents", "agents_status", "shutdown_agent",
		"send_message_to_parent", "respond_to_parent_message", "submit",
	}
	for _, name := range want {
		if _, err := parent.reg.Lookup(name); err != nil {
			t.Errorf("tool %q is not registered: %v", name, err)
		}
	}
}

// TestBlockingWorkflowToolsRunAsJobs checks that the two tools which wait do
// NOT run inline.
//
// An inline blocking tool holds the actor loop, which would make stage 3's
// interruption impossible: there would be no loop left to deliver the
// message that is supposed to interrupt it. This is a positive control as
// well, since it asserts the non-blocking tools DO run inline.
func TestBlockingWorkflowToolsRunAsJobs(t *testing.T) {
	cfg := fleetConfig(t, func(req map[string]any) []map[string]any { return nil })
	parent := fleetParent(t, cfg)

	blocking := map[string]bool{
		"join_agents": true, "send_message_to_parent": true,
		"spawn_sub_agent": false, "agents_status": false,
		"shutdown_agent": false, "respond_to_parent_message": false, "submit": false,
	}
	for name, wantBlocking := range blocking {
		tool, err := parent.reg.Lookup(name)
		if err != nil {
			t.Fatalf("lookup %q: %v", name, err)
		}
		if wantBlocking && tool.NoJob {
			t.Errorf("%q blocks but runs inline; it would hold the actor loop", name)
		}
		if !wantBlocking && !tool.NoJob {
			t.Errorf("%q returns immediately but runs as a job", name)
		}
	}
}

// TestSubmitIsRecordedOnTheChildItself checks that finishing needs no
// back-pointer to a parent.
func TestSubmitIsRecordedOnTheChildItself(t *testing.T) {
	a, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("agent: %v", err)
	}
	t.Cleanup(func() { _ = a.Shutdown() })

	if a.Submission() != nil {
		t.Fatal("a fresh agent already has a submission")
	}
	a.Submit("my result")
	got := a.Submission()
	if got == nil || got.Payload != "my result" {
		t.Fatalf("Submission() = %+v, want the submitted payload", got)
	}
}

// TestFailedChildTurnIsNotIdle checks that a child whose turn errored is
// reported as failed rather than idle.
//
// The two look identical from outside: both are children that are not
// running and have not submitted. Conflating them sends you to the wrong
// logs, because a child that died on its first turn, say because its model
// name does not exist, has a completely different cause from one that ran
// and stopped early.
//
// This test exists because a mutation audit found the distinction was
// claimed in a comment and checked nowhere.
func TestFailedChildTurnIsNotIdle(t *testing.T) {
	// A vendor that refuses every request, so the child's turn errors.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "upstream is down", http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)
	cfg := Config{
		Model: "claude-sonnet-5", Vendor: common.VendorAnthropic,
		BaseURL: srv.URL, APIKey: "test", DisableStreaming: true,
	}
	parent := fleetParent(t, cfg)

	sub, err := parent.SubAgents().Spawn("doomed", "task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("task"); err != nil {
		t.Fatalf("start: %v", err)
	}

	res := parent.SubAgents().Join([]string{"doomed"}, 30*time.Second)
	if got := res.Outcomes[0].Outcome; got != "failed" {
		t.Errorf("outcome = %q, want \"failed\"", got)
	}
	if sub.State() != SubAgentFailed {
		t.Errorf("state = %v, want failed", sub.State())
	}
	if sub.Err() == nil {
		t.Error("a failed child reports no error")
	}
	// A failed child HAS finished: it will not do anything else, so join
	// must not block on it forever waiting for a result that cannot come.
	if !sub.Finished() {
		t.Error("a failed child does not count as finished, so join would hang on it")
	}
}
