package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/waywardgeek/ensemble/agent/internal/common"
)

// textReply makes the fake vendor answer with prose and never a tool call, so
// the child finishes its turn without submitting. That is the mistake this
// stage exists to correct.
func idleTextReply(turns *atomic.Int64) func(map[string]any) []map[string]any {
	return func(map[string]any) []map[string]any {
		turns.Add(1)
		return []map[string]any{{"type": "text", "text": "Thinking about it."}}
	}
}

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// The threshold is part of the contract, so it is pinned here rather than
// only in assertions written in terms of the constant. Those move with it: a
// test expecting MaxIdleWakes+1 turns keeps passing when the value changes,
// which is exactly how a published number drifts without anyone noticing.
func TestWakeBudgetIsThree(t *testing.T) {
	if MaxIdleWakes != 3 {
		t.Fatalf("MaxIdleWakes = %d, want 3", MaxIdleWakes)
	}
}

// An agent that stops without submitting has not finished, and the framework
// says so rather than accepting silence as a result. The wake budget is what
// stops that correction becoming a loop.
func TestIdleChildIsWokenThenEscalatedToItsParent(t *testing.T) {
	var turns atomic.Int64
	parent := fleetParent(t, fleetConfig(t, idleTextReply(&turns)))

	sub, err := parent.SubAgents().Spawn("stuck", "do the thing", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("do the thing"); err != nil {
		t.Fatalf("start: %v", err)
	}

	waitFor(t, "the escalation to reach the parent", func() bool {
		return len(parent.SubAgents().PendingMessages()) > 0
	})

	// One turn for the task, then exactly MaxIdleWakes wakes. A fifth turn
	// would mean the budget is not holding.
	if got := turns.Load(); got != int64(MaxIdleWakes+1) {
		t.Fatalf("ran %d turns, want %d (one task plus %d wakes)", got, MaxIdleWakes+1, MaxIdleWakes)
	}

	msgs := parent.SubAgents().PendingMessages()
	if len(msgs) != 1 {
		t.Fatalf("parent has %d messages, want exactly 1 escalation", len(msgs))
	}
	// The parent has to act on this, which means knowing which agent is
	// stuck and that no result is coming.
	if msgs[0].From != "stuck" {
		t.Errorf("escalation From = %q, want the stuck agent's address", msgs[0].From)
	}
	if !strings.Contains(msgs[0].Content, "stuck") {
		t.Errorf("escalation does not name the agent: %q", msgs[0].Content)
	}

	// Stopped means stopped. Waking again after escalating is the loop the
	// budget exists to prevent.
	before := turns.Load()
	time.Sleep(200 * time.Millisecond)
	if after := turns.Load(); after != before {
		t.Errorf("ran %d more turns after escalating, want none", after-before)
	}
}

// The count measures stopping in a row. An agent that submits has done what
// it was asked, so earlier lapses must not be held against it.
//
// The child here stops once, is woken, and then submits, so the reset runs
// through the path a real agent takes rather than being called directly.
func TestSubmittingResetsTheWakeCount(t *testing.T) {
	var turns atomic.Int64
	parent := fleetParent(t, fleetConfig(t, func(map[string]any) []map[string]any {
		if turns.Add(1) == 1 {
			return []map[string]any{{"type": "text", "text": "Thinking about it."}}
		}
		return toolUse("t1", "submit", map[string]any{"result": "done"})
	}))

	sub, err := parent.SubAgents().Spawn("resets", "do the thing", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("do the thing"); err != nil {
		t.Fatalf("start: %v", err)
	}

	waitFor(t, "the child to submit", func() bool { return sub.State() == SubAgentDone })

	// Two turns proves a wake happened: the first stopped without
	// submitting, the second was the framework waking it. Counting the
	// agent's turns rather than the vendor's requests matters, because one
	// turn makes several requests while it works through a tool call.
	if got := sub.Turns(); got != 2 {
		t.Fatalf("ran %d turns, want 2 (one stop, one wake that submitted)", got)
	}
	if got := sub.Wakes(); got != 0 {
		t.Fatalf("wakes = %d after submitting, want 0", got)
	}
}

// A failed turn is a real error, not a model losing the thread. Waking it
// spends another turn arriving at the same failure.
func TestFailedChildIsNotWoken(t *testing.T) {
	var turns atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		turns.Add(1)
		http.Error(w, "upstream is down", http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)
	parent := fleetParent(t, Config{
		Model: "claude-sonnet-5", Vendor: common.VendorAnthropic,
		BaseURL: srv.URL, APIKey: "test", DisableStreaming: true,
	})

	sub, err := parent.SubAgents().Spawn("broken", "do the thing", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if err := sub.Start("do the thing"); err != nil {
		t.Fatalf("start: %v", err)
	}

	waitFor(t, "the child to fail", func() bool { return sub.State() == SubAgentFailed })

	before := turns.Load()
	time.Sleep(200 * time.Millisecond)
	if after := turns.Load(); after != before {
		t.Errorf("made %d more requests after failing, want none", after-before)
	}
	if got := sub.Wakes(); got != 0 {
		t.Errorf("wakes = %d for a failed turn, want 0", got)
	}
	if n := len(parent.SubAgents().PendingMessages()); n != 0 {
		t.Errorf("a failed turn escalated %d messages; escalation is for an agent that stops, not one that errors", n)
	}
}
