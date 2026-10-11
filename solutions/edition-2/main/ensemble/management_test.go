package ensemble_test

import (
	"context"
	"encoding/json"
	"ensemble/ensemble"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type activityLog struct {
	mu     sync.Mutex
	phases []string
}

// Observe records actual actor phases independently of status tool text.
func (l *activityLog) Observe(o ensemble.Observation) {
	if a, ok := o.(ensemble.ActivityChanged); ok {
		l.mu.Lock()
		l.phases = append(l.phases, a.Activity.Phase)
		l.mu.Unlock()
	}
}

// Only the provider is replaced: real actors, tools, jobs and durable submission
// establish the phase, output cursor and disk fallback contracts together.
func TestManagementActivityOutputAndSubmission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(fleetResultText(body), "tool text") {
			fleetCall(w, "run_command", map[string]string{"command": "printf 'tool text'"})
		} else {
			fleetCall(w, "submit", map[string]string{"result": "DONE_RESULT"})
		}
	}))
	defer server.Close()
	a, _, _ := fleetAgent(t, server)
	phases := &activityLog{}
	root := a.Ensemble().(*ensemble.Ensemble)
	root.SetAgentFactory(func(parent ensemble.Agent, spec ensemble.AgentSpec) (ensemble.Agent, error) {
		c, err := root.NewChildAgent(parent, parent.Config(), spec)
		if err == nil {
			c.Observe(phases)
		}
		return c, err
	})
	if _, err := a.Fleet().Spawn("worker", "TASK", ""); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := a.Fleet().Join(ctx, []string{"worker"})
	if err != nil || !strings.Contains(result, "DONE_RESULT") {
		t.Fatal(result, err)
	}
	child, _ := a.Fleet().Find("worker")
	if err = child.Flush(ctx); err != nil {
		t.Fatal(err)
	}
	phases.mu.Lock()
	got := strings.Join(phases.phases, ",")
	phases.mu.Unlock()
	if !strings.Contains(got, "thinking,tool_call") || !strings.HasSuffix(got, "idle") {
		t.Fatal("missing ordered phases", got)
	}
	output, err := a.Fleet().Output(ctx, "worker", 0, false)
	if err != nil || !strings.Contains(output.Content, "tool text") || output.Status.Usage.Input == 0 {
		t.Fatal(output, err)
	}
	next, err := a.Fleet().Output(ctx, "worker", output.Cursor, false)
	if err != nil || next.Content != "" {
		t.Fatal("cursor repeats output", next, err)
	}
	child.Fleet().Begin()
	submission, err := a.Fleet().Submission("worker")
	if err != nil || submission != "DONE_RESULT" {
		t.Fatal("durable fallback", submission, err)
	}
}

// Guessing another tree's ID must reveal no more than a genuinely absent name.
func TestManagementOwnership(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fleetCall(w, "submit", map[string]string{"result": "owned"})
	}))
	defer server.Close()
	a, _, _ := fleetAgent(t, server)
	b, _, _ := fleetAgent(t, server)
	if _, err := b.Fleet().Spawn("foreign", "TASK", ""); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	for _, name := range []string{"shutdown_agent", "kill_agent", "interrupt_agent", "send_message", "ask_agent", "wait_for_agent", "get_submission", "get_sub_agent_status", "respond_to_parent_message", "join_agents", "wait_for_agent_change"} {
		var messages []string
		for _, id := range []string{"foreign", "absent"} {
			args, _ := json.Marshal(map[string]any{"agent_id": id, "agent_ids": []string{id}, "content": "hello", "response": "reply", "timeout_seconds": 0.1})
			result, err := a.Invoke(ctx, name, args)
			message := result
			if err != nil {
				message = err.Error()
			}
			if !strings.Contains(message, "unknown agent") {
				t.Fatalf("%s wrong rejection: %s", name, message)
			}
			messages = append(messages, strings.ReplaceAll(message, id, "TARGET"))
		}
		if messages[0] != messages[1] {
			t.Fatalf("%s distinguishes foreign and absent: %q", name, messages)
		}
	}
}

// A parked real HTTP request lets us distinguish waiting for completion from
// canceling it. The stuck detector is a timeout property with that working baseline.
func TestManagementWaitAndShutdownLifetimes(t *testing.T) {
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entered <- struct{}{}
		select {
		case <-release:
			fleetCall(w, "submit", map[string]string{"result": "finished"})
		case <-r.Context().Done():
		}
	}))
	defer server.Close()
	a, _, _ := fleetAgent(t, server)
	if _, err := a.Fleet().Spawn("worker", "TASK", ""); err != nil {
		t.Fatal(err)
	}
	<-entered
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stuck, err := a.Fleet().WaitChange(ctx, []string{"worker"}, "", 100*time.Millisecond)
	if err != nil || !strings.Contains(stuck, `"reason":"stuck"`) {
		t.Fatal("stuck wait", stuck, err)
	}
	graceful := make(chan error, 1)
	go func() { graceful <- a.Fleet().GracefulShutdown(ctx, "worker") }()
	select {
	case err := <-graceful:
		t.Fatal("graceful canceled active turn", err)
	case <-time.After(30 * time.Millisecond):
	}
	close(release)
	if err := <-graceful; err != nil {
		t.Fatal(err)
	}
	if _, err := a.Fleet().Find("worker"); err == nil {
		t.Fatal("graceful retained child")
	}
}
