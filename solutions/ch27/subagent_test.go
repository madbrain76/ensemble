package agent

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testFactory returns a factory that builds real child agents and records
// how many times it was called.
//
// A real *Agent rather than a stand-in: the thing under test is that a child
// is a COMPLETE agent with its own files, and a fake agent would have
// whatever files the fake chose to have, which would prove nothing.
func testFactory(t *testing.T, calls *int) AgentFactory {
	t.Helper()
	return func(spec AgentSpec, task string) (*Agent, error) {
		*calls++
		return NewAgent(testSpecConfig(), spec)
	}
}

func newTestParent(t *testing.T) *Agent {
	t.Helper()
	p, err := NewAgent(testSpecConfig(), AgentSpec{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = p.Shutdown() })
	return p
}

// TestSpawnedChildOwnsItsOwnFiles is the stage 1 contract: a child is a full
// agent with its own data directory and its own save file.
//
// The assertion is on the PATHS, not on the manager's bookkeeping. A manager
// that recorded a child correctly while handing it the parent's save file
// would pass any check that asked the manager what it did, and would corrupt
// the parent's history the first time either one saved.
func TestSpawnedChildOwnsItsOwnFiles(t *testing.T) {
	parent := newTestParent(t)
	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))

	sub, err := m.Spawn("worker", "do the thing", nil)
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	if calls != 1 {
		t.Errorf("factory called %d times, want 1: the library must not construct children itself", calls)
	}

	childDir := sub.Agent.Spec().DataDir
	parentDir := parent.Spec().DataDir

	if childDir == parentDir {
		t.Fatalf("child shares the parent's data directory %q", childDir)
	}
	if !strings.HasPrefix(childDir, parentDir) {
		t.Errorf("child directory %q is not under the parent's %q", childDir, parentDir)
	}
	if info, err := os.Stat(childDir); err != nil || !info.IsDir() {
		t.Errorf("child data directory %q was not created on disk: %v", childDir, err)
	}

	// The save file is the specific shared-state hazard. Two agents writing
	// one save file interleave their event logs and neither can be replayed.
	childSave := sub.Agent.SavePath()
	if childSave == parent.SavePath() {
		t.Fatalf("child and parent share save file %q", childSave)
	}
	if filepath.Dir(childSave) != childDir {
		t.Errorf("child save file %q is not inside the child's directory %q", childSave, childDir)
	}
}

// TestSpawnRejectsDuplicateID checks that a reused id is refused rather than
// silently replacing a running agent.
func TestSpawnRejectsDuplicateID(t *testing.T) {
	parent := newTestParent(t)
	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))

	if _, err := m.Spawn("dup", "first", nil); err != nil {
		t.Fatalf("first spawn: %v", err)
	}
	if _, err := m.Spawn("dup", "second", nil); err == nil {
		t.Fatal("second spawn with the same id succeeded, want an error")
	}
	if got := len(m.List()); got != 1 {
		t.Errorf("fleet has %d children, want 1", got)
	}
}

// TestSpawnWithoutFactoryRefuses checks that a manager with no factory fails
// loudly instead of inventing a default agent.
//
// A default would have to pick a model, which is exactly the decision the
// factory indirection exists to keep out of the library.
func TestSpawnWithoutFactoryRefuses(t *testing.T) {
	parent := newTestParent(t)
	m := NewSubAgentManager(parent, nil)

	if _, err := m.Spawn("orphan", "task", nil); err == nil {
		t.Fatal("spawn with no factory succeeded, want an error")
	}
}

// TestSpawnClampsChildPermissions checks that the library narrows a child's
// permissions before the factory ever sees the spec.
//
// This is the positive control's twin: it is not enough that a child CAN be
// restricted, the restriction has to survive a caller who explicitly asks
// for more. The adjust callback here asks for exactly that.
func TestSpawnClampsChildPermissions(t *testing.T) {
	parent, err := NewAgent(testSpecConfig(), AgentSpec{
		DataDir:         t.TempDir(),
		SafeMode:        true,
		EnableWebSearch: false,
	})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	var seen AgentSpec
	m := NewSubAgentManager(parent, func(spec AgentSpec, task string) (*Agent, error) {
		seen = spec
		return NewAgent(testSpecConfig(), spec)
	})

	// A child that asks to escape both restrictions.
	_, err = m.Spawn("greedy", "task", func(s *AgentSpec) {
		s.SafeMode = false
		s.EnableWebSearch = true
	})
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}

	if !seen.SafeMode {
		t.Error("child escaped the parent's safe mode")
	}
	if seen.EnableWebSearch {
		t.Error("child escaped the parent's web search restriction")
	}
}

// TestSpawnInheritsRestrictionsByDefault is the positive control for the
// clamp test above. A test that only ever checks a restriction survives an
// explicit attempt to widen it would also pass if the manager restricted
// every child unconditionally, which would be a different bug.
func TestSpawnInheritsRestrictionsByDefault(t *testing.T) {
	parent, err := NewAgent(testSpecConfig(), AgentSpec{
		DataDir:         t.TempDir(),
		EnableWebSearch: true,
	})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	var seen AgentSpec
	m := NewSubAgentManager(parent, func(spec AgentSpec, task string) (*Agent, error) {
		seen = spec
		return NewAgent(testSpecConfig(), spec)
	})
	if _, err := m.Spawn("heir", "task", nil); err != nil {
		t.Fatalf("spawn: %v", err)
	}

	if !seen.EnableWebSearch {
		t.Error("child did not inherit the parent's web search capability")
	}
	if seen.SafeMode {
		t.Error("child was put in safe mode its parent was not in")
	}
}

// TestShutdownRemovesChild checks teardown, including that a child is
// removed from the fleet so its id can be reused.
func TestShutdownRemovesChild(t *testing.T) {
	parent := newTestParent(t)
	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))

	if _, err := m.Spawn("tmp", "task", nil); err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if _, ok := m.Get("tmp"); !ok {
		t.Fatal("spawned child is not in the fleet")
	}
	if err := m.Shutdown("tmp"); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if _, ok := m.Get("tmp"); ok {
		t.Error("child still in the fleet after shutdown")
	}
	if err := m.Shutdown("tmp"); err == nil {
		t.Error("shutting down an unknown child succeeded, want an error")
	}
	// The id is free again only if removal actually happened.
	if _, err := m.Spawn("tmp", "again", nil); err != nil {
		t.Errorf("id not reusable after shutdown: %v", err)
	}
}

// TestListIsSortedByID pins the ordering that keeps status output stable
// between turns.
func TestListIsSortedByID(t *testing.T) {
	parent := newTestParent(t)
	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))

	for _, id := range []string{"charlie", "alpha", "bravo"} {
		if _, err := m.Spawn(id, "task", nil); err != nil {
			t.Fatalf("spawn %s: %v", id, err)
		}
	}
	var got []string
	for _, sub := range m.List() {
		got = append(got, sub.ID)
	}
	want := []string{"alpha", "bravo", "charlie"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("List() = %v, want %v", got, want)
		}
	}
}

// TestShutdownAllLeavesNoOrphans checks that tearing down a parent tears
// down its whole fleet.
func TestShutdownAllLeavesNoOrphans(t *testing.T) {
	parent := newTestParent(t)
	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))

	for _, id := range []string{"a", "b", "c"} {
		if _, err := m.Spawn(id, "task", nil); err != nil {
			t.Fatalf("spawn %s: %v", id, err)
		}
	}
	m.ShutdownAll()
	if got := len(m.List()); got != 0 {
		t.Errorf("%d children left after ShutdownAll, want 0", got)
	}
}

// TestStateZeroValueIsNotAValidState guards the iota+1 choice. If someone
// "tidies" the constant block back to starting at zero, an unregistered
// SubAgent starts reading as running.
func TestStateZeroValueIsNotAValidState(t *testing.T) {
	var zero SubAgentState
	if zero.String() != "unknown" {
		t.Errorf("zero SubAgentState renders as %q, want \"unknown\"", zero.String())
	}
	for _, s := range []SubAgentState{SubAgentRunning, SubAgentIdle, SubAgentDone, SubAgentFailed} {
		if s == 0 {
			t.Errorf("state %v has the zero value", s)
		}
		if s.String() == "unknown" {
			t.Errorf("valid state %d renders as \"unknown\"", int(s))
		}
	}
}

// TestSpawnErrorsAreNotWidenedPermissions is a guard against a future edit
// turning the clamp diagnostic into a spawn failure. The clamp must narrow
// and continue, not refuse: the narrowing is the wall, the error is only
// the alarm.
func TestSpawnErrorsAreNotWidenedPermissions(t *testing.T) {
	parent, err := NewAgent(testSpecConfig(), AgentSpec{
		DataDir:  t.TempDir(),
		SafeMode: true,
	})
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	t.Cleanup(func() { _ = parent.Shutdown() })

	calls := 0
	m := NewSubAgentManager(parent, testFactory(t, &calls))
	sub, err := m.Spawn("greedy", "task", func(s *AgentSpec) { s.SafeMode = false })
	if errors.Is(err, WidenedPermissionsError{}) {
		t.Fatal("spawn refused an over-broad request; it must clamp and continue")
	}
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	if !sub.Agent.Spec().SafeMode {
		t.Error("child is not in safe mode despite a safe-mode parent")
	}
}
