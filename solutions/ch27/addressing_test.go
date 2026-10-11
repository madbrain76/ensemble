package agent

import "testing"

// TestTopLevelAgentReportsToUser pins the terminus of the parent chain.
//
// This is what makes escalation total. Stage 4 escalates a stuck agent to
// its parent, and without this an agent at the top of the tree would have
// nowhere to send the one report a human most needs to see.
func TestTopLevelAgentReportsToUser(t *testing.T) {
	a, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("agent: %v", err)
	}
	t.Cleanup(func() { _ = a.Shutdown() })

	if got := a.Spec().ParentID; got != UserID {
		t.Errorf("top-level agent ParentID = %q, want %q", got, UserID)
	}
}

// TestChildAddressesItsParent checks that a spawned child knows both its own
// address and its parent's.
func TestChildAddressesItsParent(t *testing.T) {
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

	spec := sub.Agent.Spec()
	if spec.ID != "worker" {
		t.Errorf("child ID = %q, want \"worker\"", spec.ID)
	}
	if spec.ParentID != "boss" {
		t.Errorf("child ParentID = %q, want \"boss\"", spec.ParentID)
	}
	// The child must not inherit its parent's own address, which is the
	// specific way two agents end up sharing one id.
	if spec.ID == parent.Spec().ID {
		t.Errorf("child inherited its parent's id %q", spec.ID)
	}
}

// TestGrandchildAddressesItsOwnParent checks the chain is relative at every
// level rather than always pointing at the root.
func TestGrandchildAddressesItsOwnParent(t *testing.T) {
	root, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir(), ID: "root"})
	if err != nil {
		t.Fatalf("root: %v", err)
	}
	t.Cleanup(func() { _ = root.Shutdown() })

	calls := 0
	m := NewSubAgentManager(root, testFactory(t, &calls))
	child, err := m.Spawn("mid", "task", nil)
	if err != nil {
		t.Fatalf("spawn child: %v", err)
	}

	m2 := NewSubAgentManager(child.Agent, testFactory(t, &calls))
	grand, err := m2.Spawn("leaf", "task", nil)
	if err != nil {
		t.Fatalf("spawn grandchild: %v", err)
	}

	if got := grand.Agent.Spec().ParentID; got != "mid" {
		t.Errorf("grandchild ParentID = %q, want \"mid\"", got)
	}
	if got := grand.Agent.Spec().ParentID; got == UserID {
		t.Error("grandchild reports to the user, skipping its real parent")
	}
}

// TestUserIDIsReserved checks that no agent may occupy the user's address.
//
// An agent registered there could receive messages a child addressed to the
// human, so this is refused at spawn rather than detected afterwards.
func TestUserIDIsReserved(t *testing.T) {
	parent := newTestParent(t)
	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))

	if _, err := m.Spawn(UserID, "impersonate the human", nil); err == nil {
		t.Fatalf("spawning a child with id %q succeeded, want an error", UserID)
	}
	if calls != 0 {
		t.Error("factory was called for a reserved id; the refusal must come first")
	}
}

// TestAdjustCannotRenameAgent checks that the manager owns identity.
//
// The map key and the agent's own address have to agree, or a reply routes
// to the wrong child. A caller hook that could set ID would break that
// agreement silently.
func TestAdjustCannotRenameAgent(t *testing.T) {
	parent := newTestParent(t)
	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))

	sub, err := m.Spawn("real", "task", func(s *AgentSpec) {
		s.ID = "impostor"
		s.ParentID = "elsewhere"
	})
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if got := sub.Agent.Spec().ID; got != "real" {
		t.Errorf("adjust renamed the agent to %q, want \"real\"", got)
	}
	// Reparenting is refused for the same reason renaming is. A child
	// pointed at a stranger would send its escalations there.
	if got := sub.Agent.Spec().ParentID; got != parent.Spec().ID {
		t.Errorf("adjust reparented the agent to %q, want %q", got, parent.Spec().ID)
	}
	// Look the agent up the way a message would, and check it is the same one.
	found, ok := m.Get("real")
	if !ok || found != sub {
		t.Error("agent is not reachable under the id it was registered with")
	}
}

// TestUnnamedParentDoesNotRerouteChildrenToUser is a regression test.
//
// The bug: a parent constructed without an explicit id had an empty address,
// so Spawn copied "" into each child's ParentID, and NewAgent then read that
// empty value as "top-level" and defaulted it to the user. The children of an
// unnamed parent therefore escalated PAST their real parent, directly to the
// human, and nothing reported an error.
//
// The failure is silent, which is why it is pinned here rather than left to
// the general addressing tests: every one of those named its agents
// explicitly, and naming them is exactly what hid the bug.
func TestUnnamedParentDoesNotRerouteChildrenToUser(t *testing.T) {
	// Deliberately no ID. This is the shape that broke.
	parent, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	if parent.Spec().ID == "" {
		t.Fatal("agent was constructed with an empty address")
	}

	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))
	sub, err := m.Spawn("worker", "task", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	if got := sub.Agent.Spec().ParentID; got == UserID {
		t.Error("child of an unnamed parent reports to the user, skipping its real parent")
	}
	if got, want := sub.Agent.Spec().ParentID, parent.Spec().ID; got != want {
		t.Errorf("child ParentID = %q, want its real parent %q", got, want)
	}
}
