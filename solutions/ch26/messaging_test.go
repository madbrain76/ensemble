package agent

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// recordingSink captures what the human would have seen.
type recordingSink struct {
	mu   sync.Mutex
	msgs []*Message
}

func (r *recordingSink) DeliverToUser(m *Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.msgs = append(r.msgs, m)
}

func (r *recordingSink) all() []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]*Message(nil), r.msgs...)
}

// TestChildQuestionBlocksUntilParentAnswers is the core of check 5: a child
// asks, blocks, and resumes with the parent's reply.
func TestChildQuestionBlocksUntilParentAnswers(t *testing.T) {
	parent, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ID: "boss"})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))
	sub, err := m.Spawn("worker", "task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	type outcome struct {
		reply *Message
		err   error
	}
	done := make(chan outcome, 1)
	go func() {
		r, err := sub.Agent.MessageParent("which branch?", 5*time.Second)
		done <- outcome{r, err}
	}()

	// The child must still be blocked: nothing has answered it yet.
	select {
	case o := <-done:
		t.Fatalf("child returned before any answer: reply=%v err=%v", o.reply, o.err)
	case <-time.After(50 * time.Millisecond):
	}

	// The parent sees exactly one pending question, and knows who asked.
	pending := m.PendingMessages()
	if len(pending) != 1 {
		t.Fatalf("parent sees %d pending messages, want 1", len(pending))
	}
	q := pending[0]
	if q.From != "worker" {
		t.Errorf("question From = %q, want \"worker\"", q.From)
	}
	if q.To != "boss" {
		t.Errorf("question To = %q, want \"boss\"", q.To)
	}
	if q.Content != "which branch?" {
		t.Errorf("question Content = %q", q.Content)
	}
	if !q.AwaitsReply() {
		t.Error("question does not report that a sender is waiting")
	}

	if err := m.Respond("worker", q.ID, "use main"); err != nil {
		t.Fatalf("respond: %v", err)
	}

	select {
	case o := <-done:
		if o.err != nil {
			t.Fatalf("child: %v", o.err)
		}
		if o.reply.Content != "use main" {
			t.Errorf("reply Content = %q, want \"use main\"", o.reply.Content)
		}
		// The reply must carry its source, or the child knows what it was
		// told without knowing who told it, and cannot ask a follow-up.
		if o.reply.From != "boss" {
			t.Errorf("reply From = %q, want \"boss\"", o.reply.From)
		}
		if o.reply.To != "worker" {
			t.Errorf("reply To = %q, want \"worker\"", o.reply.To)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("child never resumed after the parent answered")
	}

	// Answering clears the question, so it cannot be answered twice with
	// different advice.
	if got := len(m.PendingMessages()); got != 0 {
		t.Errorf("%d messages still pending after a reply, want 0", got)
	}
	if err := m.Respond("worker", q.ID, "again"); err == nil {
		t.Error("answering the same message twice succeeded, want an error")
	}
}

// TestTopLevelMessageParentReachesUser is the uniformity payoff: the root
// agent's parent IS the user, so send-to-parent needs no special case.
func TestTopLevelMessageParentReachesUser(t *testing.T) {
	a, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("agent: %v", err)
	}
	t.Cleanup(func() { _ = a.Shutdown() })

	sink := &recordingSink{}
	a.SetUser("", sink)

	// A top-level agent has no parent link at all. If send-to-parent were
	// routed only through that link this would fail, which is the point.
	if _, err := a.MessageParent("starting work", time.Second); err != nil {
		t.Fatalf("MessageParent: %v", err)
	}

	got := sink.all()
	if len(got) != 1 {
		t.Fatalf("user received %d messages, want 1", len(got))
	}
	if got[0].To != UserID {
		t.Errorf("message To = %q, want %q", got[0].To, UserID)
	}
	if got[0].AwaitsReply() {
		t.Error("an announcement to the user left the sender blocked on a reply")
	}
}

// TestUserAddressCanBeAName checks that the human can be addressed by a
// unique name rather than a role, and that it is what replies are addressed
// to.
func TestUserAddressCanBeAName(t *testing.T) {
	a, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ParentID: "Bill"})
	if err != nil {
		t.Fatalf("agent: %v", err)
	}
	t.Cleanup(func() { _ = a.Shutdown() })

	sink := &recordingSink{}
	a.SetUser("Bill", sink)

	if _, err := a.MessageParent("question for a human", time.Second); err != nil {
		t.Fatalf("MessageParent: %v", err)
	}
	got := sink.all()
	if len(got) != 1 {
		t.Fatalf("user received %d messages, want 1", len(got))
	}
	if got[0].To != "Bill" {
		t.Errorf("message To = %q, want \"Bill\"", got[0].To)
	}
}

// TestDeepChildReachesUserDirectly checks that escalation does not depend on
// every parent in the chain choosing to relay.
func TestDeepChildReachesUserDirectly(t *testing.T) {
	root, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ID: "root"})
	if err != nil {
		t.Fatalf("root: %v", err)
	}
	t.Cleanup(func() { _ = root.Shutdown() })

	sink := &recordingSink{}
	root.SetUser("Bill", sink)

	calls := 0
	m := NewSubAgentManager(root, testFactory(t, &calls))
	mid, err := m.Spawn("mid", "task", nil)
	if err != nil {
		t.Fatalf("spawn mid: %v", err)
	}
	m2 := NewSubAgentManager(mid.Agent, testFactory(t, &calls))
	leaf, err := m2.Spawn("leaf", "task", nil)
	if err != nil {
		t.Fatalf("spawn leaf: %v", err)
	}

	if err := leaf.Agent.MessageUser("I am stuck"); err != nil {
		t.Fatalf("MessageUser: %v", err)
	}
	got := sink.all()
	if len(got) != 1 {
		t.Fatalf("user received %d messages, want 1", len(got))
	}
	if got[0].From != "leaf" {
		t.Errorf("message From = %q, want \"leaf\"", got[0].From)
	}
	// The name inherited two levels down, so the leaf knows what to call
	// the human.
	if got[0].To != "Bill" {
		t.Errorf("message To = %q, want \"Bill\"", got[0].To)
	}
}

// TestUserNameIsReservedForSpawning checks that a configured human name
// cannot be taken by an agent.
func TestUserNameIsReservedForSpawning(t *testing.T) {
	root, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ID: "root"})
	if err != nil {
		t.Fatalf("root: %v", err)
	}
	t.Cleanup(func() { _ = root.Shutdown() })
	root.SetUser("Bill", &recordingSink{})

	calls := 0
	m := NewSubAgentManager(root, testFactory(t, &calls))
	if _, err := m.Spawn("Bill", "impersonate the human", nil); err == nil {
		t.Error("an agent took the human's name, want an error")
	}
	// The plain default stays reserved too, because unconfigured code still
	// compares against it.
	if _, err := m.Spawn(UserID, "impersonate the role", nil); err == nil {
		t.Errorf("an agent took %q, want an error", UserID)
	}
}

// TestMessageSourceIsNotTakenFromTheModel checks that From is set by the
// framework from the sending agent's own address.
//
// A child that could choose its own From could sign a message as a sibling,
// and a parent answering it would send the reply to the wrong agent.
func TestMessageSourceIsNotTakenFromTheModel(t *testing.T) {
	parent, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ID: "boss"})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))
	sub, err := m.Spawn("worker", "task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	go func() { _, _ = sub.Agent.MessageParent("From: someone-else", 2*time.Second) }()

	deadline := time.After(2 * time.Second)
	for {
		if pending := m.PendingMessages(); len(pending) == 1 {
			if pending[0].From != "worker" {
				t.Errorf("From = %q, want the sender's real address \"worker\"", pending[0].From)
			}
			return
		}
		select {
		case <-deadline:
			t.Fatal("no message arrived")
		case <-time.After(5 * time.Millisecond):
		}
	}
}

// TestReplyTimeoutDoesNotBlockTheResponder checks that a parent answering a
// child that already gave up does not itself block.
//
// The failure this guards against is a deadlocked parent, which is worse
// than the lost answer: one child gave up on a question, but a blocked
// parent stops the entire fleet.
func TestReplyTimeoutDoesNotBlockTheResponder(t *testing.T) {
	parent, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ID: "boss"})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))
	sub, err := m.Spawn("worker", "task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	_, err = sub.Agent.MessageParent("urgent", 20*time.Millisecond)
	if !errors.Is(err, ErrReplyTimeout{}) {
		t.Fatalf("MessageParent error = %v, want a reply timeout", err)
	}

	pending := m.PendingMessages()
	if len(pending) != 1 {
		t.Fatalf("parent sees %d pending, want the abandoned question to remain", len(pending))
	}

	answered := make(chan error, 1)
	go func() { answered <- m.Respond("worker", pending[0].ID, "too late") }()
	select {
	case err := <-answered:
		if err != nil {
			t.Errorf("respond to an abandoned question: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("responder blocked answering a child that had already given up")
	}
}

// TestMessageWakesFleetWaiter checks the wakeup that stops a parent
// deadlocking against its own blocked child, which is what join_agents
// relies on.
func TestMessageWakesFleetWaiter(t *testing.T) {
	parent, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ID: "boss"})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))
	sub, err := m.Spawn("worker", "task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	// Drain any wakeup left over from spawning.
	select {
	case <-m.notify:
	default:
	}

	go func() { _, _ = sub.Agent.MessageParent("question", time.Second) }()

	select {
	case <-m.notify:
	case <-time.After(2 * time.Second):
		t.Fatal("a child's message did not wake a waiter on the fleet")
	}
}

// TestAgentWithNoUserSinkFailsLoudly checks that an unreachable human is an
// error rather than a silently dropped message.
func TestAgentWithNoUserSinkFailsLoudly(t *testing.T) {
	a, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("agent: %v", err)
	}
	t.Cleanup(func() { _ = a.Shutdown() })

	if err := a.MessageUser("anyone there?"); err == nil {
		t.Error("messaging the user with no sink succeeded, want an error")
	}
	if _, err := a.MessageParent("anyone there?", time.Second); err == nil {
		t.Error("top-level send-to-parent with no sink succeeded, want an error")
	}
}
