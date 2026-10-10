package agent

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

// AgentFactory builds a sub-agent from a specification and a task.
//
// The library spawns children. It does not decide what they are made of.
// Which model a child runs, which credentials it presents and which endpoint
// it talks to are policy, and policy belongs to the program that assembled
// the agent in the first place. This function is how the composition root
// keeps that decision: the manager calls it and uses whatever comes back.
//
// It is a function rather than an interface because there is exactly one
// method's worth of behaviour here, and a closure over the root's own Config
// is the natural implementation. Declaring a type for it would add a name
// without adding a choice.
//
// Note what is absent. There is no Config parameter, no model name and no
// API key. If any of those appeared here the library would be back to
// deciding them, just with an extra hop.
type AgentFactory func(spec AgentSpec, task string) (*Agent, error)

// SubAgentState is what a child is doing, as the FRAMEWORK observes it
// rather than as the child reports it.
//
// The distinction is the whole reason this type exists. A child that has
// talked itself into believing it finished is exactly the failure this
// chapter is about, so a state it could set by saying so would be useless
// for detecting it. Every value here is set by the manager from something it
// watched happen: a turn started, a turn ended, submit was called.
type SubAgentState int

// Constants start at iota+1 so the zero value is never a valid state. A
// SubAgent that somehow skipped registration reads as 0, which matches no
// case and shows up as "unknown" rather than silently masquerading as the
// first real state in the list.
const (
	// SubAgentRunning means a turn is in flight right now.
	SubAgentRunning SubAgentState = iota + 1

	// SubAgentIdle means the child finished a turn and is waiting. This is
	// the suspicious state: for a child that was given a job and never
	// submitted, idle means stopped early, not finished.
	SubAgentIdle

	// SubAgentDone means the child called submit. This is the only state
	// that means finished, and the only one a child can cause deliberately.
	SubAgentDone

	// SubAgentFailed means the child's last turn ended in an error. A child
	// that died on its first turn, say because its model name does not
	// exist, used to be indistinguishable from one that ran and finished
	// cleanly, because both sat at idle.
	SubAgentFailed
)

// String renders a state for status output.
func (s SubAgentState) String() string {
	switch s {
	case SubAgentRunning:
		return "running"
	case SubAgentIdle:
		return "idle"
	case SubAgentDone:
		return "done"
	case SubAgentFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// SubAgent is one child, plus everything its parent knows about it.
//
// The child itself is a complete *Agent: its own actor, tool registry,
// settings, save file and data directory. Nothing here is a lightweight
// handle onto a shared engine, and that is deliberate. A child that shared
// its parent's save file would corrupt it, and a child that shared its tool
// registry could not be given a narrower one.
type SubAgent struct {
	// ID is the caller's chosen name, unique within one parent.
	ID string

	// Task is the instruction the child was spawned with, kept so status
	// output can say what a child is for without asking it.
	Task string

	// Agent is the child itself.
	Agent *Agent

	// mu guards every field below it. The parent's actor reads this
	// concurrently with the child's own turn goroutine writing it.
	mu sync.Mutex

	state SubAgentState
	turns int
	err   error
}

// State reports what the child is doing.
func (s *SubAgent) State() SubAgentState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

// Turns reports how many turns the child has started. A child with zero
// turns never ran at all, which is a different failure from one that ran
// and stopped, and status output that conflated them sent people reading
// the wrong logs.
func (s *SubAgent) Turns() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.turns
}

// Err reports the error that ended the child's last turn, if one did.
func (s *SubAgent) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

// setState records an observed transition.
func (s *SubAgent) setState(st SubAgentState, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state, s.err = st, err
}

// SubAgentManager spawns, tracks and shuts down one agent's children.
//
// Every agent has one, including children, which is what makes a fleet a
// tree rather than a flat list: a child that spawns its own children manages
// them itself, and its parent neither knows nor needs to.
type SubAgentManager struct {
	mu       sync.Mutex
	parent   *Agent
	factory  AgentFactory
	children map[string]*SubAgent
}

// NewSubAgentManager builds a manager for one parent agent.
//
// The factory may be nil, and a manager with no factory refuses to spawn
// rather than inventing a default. A default would have to choose a model,
// which is the decision this whole indirection exists to avoid making.
func NewSubAgentManager(parent *Agent, factory AgentFactory) *SubAgentManager {
	return &SubAgentManager{
		parent:   parent,
		factory:  factory,
		children: make(map[string]*SubAgent),
	}
}

// Spawn creates a child agent and registers it.
//
// The child's specification is DERIVED from the parent's with Child, never
// built fresh, so every restriction the parent carries is inherited by
// construction rather than by anyone remembering to copy it. Clamp then
// narrows whatever the caller changed. Both steps run here, in the library,
// so a factory cannot skip them: the factory receives a spec that has
// already been clamped and gets no say in the matter.
func (m *SubAgentManager) Spawn(id, task string, adjust func(*AgentSpec)) (*SubAgent, error) {
	if id == "" {
		return nil, fmt.Errorf("sub-agent id must not be empty")
	}

	m.mu.Lock()
	if _, exists := m.children[id]; exists {
		m.mu.Unlock()
		return nil, fmt.Errorf("sub-agent %q already exists", id)
	}
	factory := m.factory
	m.mu.Unlock()

	if factory == nil {
		return nil, fmt.Errorf("sub-agent %q: no agent factory configured", id)
	}

	// The child's directory lives under the parent's, named for the child.
	// Deriving the path rather than accepting one means two children cannot
	// be pointed at the same directory by a caller who reused an id.
	spec := m.parent.Spec().Child(m.parent.Spec().Path(agentsDir(id)))
	if adjust != nil {
		adjust(&spec)
	}

	// Clamp unconditionally, then report. The narrowing already happened by
	// the time the error is read, so a caller that drops the error still
	// gets a child no wider than its parent.
	spec, widenErr := spec.Clamp(m.parent.Spec())
	if widenErr != nil {
		m.parent.Logf("subagent %s: %v", id, widenErr)
	}

	if err := os.MkdirAll(spec.DataDir, 0o755); err != nil {
		return nil, fmt.Errorf("sub-agent %q: data directory: %w", id, err)
	}

	child, err := factory(spec, task)
	if err != nil {
		return nil, fmt.Errorf("sub-agent %q: %w", id, err)
	}

	sub := &SubAgent{ID: id, Task: task, Agent: child, state: SubAgentIdle}

	m.mu.Lock()
	defer m.mu.Unlock()
	// Re-check under the lock. Two concurrent Spawn calls with the same id
	// both passed the first check, and without this the second would
	// silently replace the first, orphaning a running agent.
	if _, exists := m.children[id]; exists {
		_ = child.Shutdown()
		return nil, fmt.Errorf("sub-agent %q already exists", id)
	}
	m.children[id] = sub
	return sub, nil
}

// Get returns a child by id.
func (m *SubAgentManager) Get(id string) (*SubAgent, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sub, ok := m.children[id]
	return sub, ok
}

// List returns every child, ordered by id.
//
// Sorted rather than map order because this backs status output that a
// model reads. Map iteration would reshuffle the list between two otherwise
// identical turns, which costs a prompt cache hit for no reason and makes
// two status readings impossible to diff by eye.
func (m *SubAgentManager) List() []*SubAgent {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*SubAgent, 0, len(m.children))
	for _, sub := range m.children {
		out = append(out, sub)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Shutdown stops one child and removes it from the fleet.
//
// Removal happens even when the child's own shutdown reports an error. The
// alternative leaves an entry that can never be spawned again under the same
// id and can never be cleaned up, which is a leak that compounds.
func (m *SubAgentManager) Shutdown(id string) error {
	m.mu.Lock()
	sub, ok := m.children[id]
	if ok {
		delete(m.children, id)
	}
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("no such sub-agent: %s", id)
	}
	return sub.Agent.Shutdown()
}

// ShutdownAll stops every child. It is what a parent's own shutdown calls,
// so that killing a parent does not leave orphans running.
func (m *SubAgentManager) ShutdownAll() {
	m.mu.Lock()
	subs := make([]*SubAgent, 0, len(m.children))
	for _, sub := range m.children {
		subs = append(subs, sub)
	}
	m.children = make(map[string]*SubAgent)
	m.mu.Unlock()

	for _, sub := range subs {
		if err := sub.Agent.Shutdown(); err != nil {
			m.parent.Logf("subagent %s: shutdown: %v", sub.ID, err)
		}
	}
}

// agentsDir names a child's directory within its parent's.
func agentsDir(id string) string { return "agents/" + id }
