package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Ensemble owns the agent map and replaces the manual wiring that used to
// live in cmd/main.go. It IS the agent factory: when any agent spawns a
// child, the Ensemble builds it, wires the factory and user sink, and
// assigns the data directory. The five composition-root bugs from chapter 25
// all traced to the host forgetting to wire one of these, and they are all
// prevented structurally here because they are done in one place.
type Ensemble struct {
	dir      string                // root data directory
	cfg      Config                // shared base config (model, API key, etc.)
	registry []AgentEntry          // persisted to ensemble.json
	agents   map[string]*Agent     // live top-level agent map
	mu       sync.RWMutex
	userAddr string                // what to call the human ("user" by default)
	userSink func(*Message)        // delivers messages to the human
}

// AgentEntry is one row in the persistent registry.
type AgentEntry struct {
	Name    string    `json:"name"`
	Model   string    `json:"model,omitempty"`
	Skill   string    `json:"skill,omitempty"`
	Created time.Time `json:"created"`
}

// registryPath is where the persistent registry lives on disk.
func (e *Ensemble) registryPath() string {
	return filepath.Join(e.dir, "ensemble.json")
}

// NewEnsemble creates an Ensemble rooted at dir with the given base config.
//
// It does NOT start any agents. Call Start after wiring the user sink
// and any other configuration. Two-phase construction (New + Start) lets
// the caller install the user sink before any agent can message it.
func NewEnsemble(dir string, cfg Config) (*Ensemble, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("ensemble: data directory: %w", err)
	}

	e := &Ensemble{
		dir:      dir,
		cfg:      cfg,
		agents:   make(map[string]*Agent),
		userAddr: UserID,
	}

	// Load the existing registry, or start empty.
	data, err := os.ReadFile(e.registryPath())
	if err == nil {
		if err := json.Unmarshal(data, &e.registry); err != nil {
			return nil, fmt.Errorf("ensemble: parse %s: %w", e.registryPath(), err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("ensemble: read %s: %w", e.registryPath(), err)
	}

	return e, nil
}

// SetUser configures the human's address and message sink. Call this before
// Start so that agents created during Start can reach the human.
func (e *Ensemble) SetUser(addr string, sink func(*Message)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if addr != "" {
		e.userAddr = addr
	}
	e.userSink = sink
}

// Start creates an Agent for every entry in the persistent registry.
//
// Each agent is fully wired: factory, user sink, data directory. Errors
// from individual agents are logged but do not prevent the rest from
// starting, because a broken agent should not take down the fleet.
func (e *Ensemble) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, entry := range e.registry {
		if _, err := e.createAgentLocked(entry); err != nil {
			// Log but continue. The alternative is refusing to start the
			// whole ensemble because one agent's settings are broken.
			fmt.Fprintf(os.Stderr, "ensemble: start %s: %v\n", entry.Name, err)
		}
	}
	return nil
}

// createAgentLocked builds, wires, and registers a live agent.
// The caller must hold e.mu.
func (e *Ensemble) createAgentLocked(entry AgentEntry) (*Agent, error) {
	spec := AgentSpec{
		DataDir: filepath.Join(e.dir, entry.Name),
		ID:      entry.Name,
	}
	if entry.Skill != "" {
		spec.Skills = []string{entry.Skill}
	}

	// Ensure the data directory exists before NewAgent tries to open files
	// inside it.
	if err := os.MkdirAll(spec.DataDir, 0o755); err != nil {
		return nil, fmt.Errorf("data directory %s: %w", spec.DataDir, err)
	}

	cfg := e.cfg
	if entry.Model != "" {
		cfg.Model = entry.Model
	}

	agent, err := NewAgent(cfg, spec)
	if err != nil {
		return nil, err
	}

	// Wire the factory so children can spawn.
	agent.SetAgentFactory(e.agentFactory())

	// Wire the user sink so agents can reach the human.
	if e.userSink != nil {
		agent.SetUser(e.userAddr, UserSinkFunc(e.userSink))
	}

	e.agents[entry.Name] = agent
	return agent, nil
}

// agentFactory returns the factory closure that builds sub-agents.
//
// Every child built by this factory inherits the same factory, so the fleet
// can grow to any depth without anyone remembering to propagate it.
func (e *Ensemble) agentFactory() AgentFactory {
	return func(spec AgentSpec, task string) (*Agent, error) {
		child, err := NewAgent(e.cfg, spec)
		if err != nil {
			return nil, err
		}
		// Propagate the factory. This is what chapter 25 bug #3 was: a child
		// whose manager held nil because the factory was never propagated.
		child.SetAgentFactory(e.agentFactory())
		return child, nil
	}
}

// saveRegistryLocked writes ensemble.json. The caller must hold e.mu.
func (e *Ensemble) saveRegistryLocked() error {
	data, err := json.MarshalIndent(e.registry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(e.registryPath(), data, 0o644)
}

// CreateAgent adds a top-level agent to the registry, creates it, wires it,
// and persists the registry. Returns the live Agent.
func (e *Ensemble) CreateAgent(name, model, skill string) (*Agent, error) {
	if name == "" {
		return nil, fmt.Errorf("ensemble: agent name must not be empty")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Reject duplicates.
	if _, exists := e.agents[name]; exists {
		return nil, fmt.Errorf("ensemble: agent %q already exists", name)
	}

	entry := AgentEntry{
		Name:    name,
		Model:   model,
		Skill:   skill,
		Created: time.Now(),
	}

	agent, err := e.createAgentLocked(entry)
	if err != nil {
		return nil, err
	}

	e.registry = append(e.registry, entry)
	if err := e.saveRegistryLocked(); err != nil {
		// Roll back the live state. The alternative is an agent that exists
		// in memory but not on disk, which would vanish on restart.
		delete(e.agents, name)
		_ = agent.Shutdown()
		return nil, fmt.Errorf("ensemble: save registry: %w", err)
	}

	return agent, nil
}

// DeleteAgent removes a top-level agent. It rejects if the agent has
// active sub-agents, then removes it from the registry, shuts it down,
// cleans up its data directory, and saves.
func (e *Ensemble) DeleteAgent(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	agent, exists := e.agents[name]
	if !exists {
		return fmt.Errorf("ensemble: no such agent: %s", name)
	}

	// Reject if the agent has active children.
	if ids := agent.ActiveChildIDs(); len(ids) > 0 {
		return fmt.Errorf("ensemble: cannot delete %s: %d active sub-agent(s) (%s)",
			name, len(ids), joinIDs(ids))
	}

	// Remove from the live map first. An agent that is removed from the map
	// but whose shutdown fails is preferable to one that stays in the map
	// after its directory has been cleaned up.
	delete(e.agents, name)

	// Remove from the persistent registry.
	for i, entry := range e.registry {
		if entry.Name == name {
			e.registry = append(e.registry[:i], e.registry[i+1:]...)
			break
		}
	}

	if err := e.saveRegistryLocked(); err != nil {
		return fmt.Errorf("ensemble: save registry after delete: %w", err)
	}

	// Shut down the agent (stops children, saves state).
	_ = agent.Shutdown()

	// Remove the data directory.
	dataDir := filepath.Join(e.dir, name)
	_ = os.RemoveAll(dataDir)

	return nil
}

// Agent returns a live agent by name, or nil if it does not exist.
func (e *Ensemble) Agent(name string) *Agent {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.agents[name]
}

// Agents returns all live top-level agents, ordered by name.
func (e *Ensemble) Agents() []*Agent {
	e.mu.RLock()
	defer e.mu.RUnlock()
	names := make([]string, 0, len(e.agents))
	for name := range e.agents {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]*Agent, len(names))
	for i, name := range names {
		out[i] = e.agents[name]
	}
	return out
}

// Shutdown stops all agents and saves.
func (e *Ensemble) Shutdown() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	var firstErr error
	for name, agent := range e.agents {
		if err := agent.Shutdown(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("ensemble: shutdown %s: %w", name, err)
		}
	}
	e.agents = make(map[string]*Agent)
	return firstErr
}

// joinIDs formats a list of IDs for error messages.
func joinIDs(ids []string) string {
	if len(ids) <= 3 {
		return fmt.Sprintf("%v", ids)
	}
	return fmt.Sprintf("%v and %d more", ids[:3], len(ids)-3)
}
