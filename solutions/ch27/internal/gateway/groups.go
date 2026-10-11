package gateway

import (
	"fmt"
	"sort"
	"sync"

	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// Registry is the gateway's authorization state: who exists, what groups
// they belong to, and who may contact whom.
//
// Authority has two layers. Base memberships come from config and are the
// seed; dynamic memberships are granted at runtime by a parent to a child it
// registered. The rule joining them is that a parent may grant only groups
// it already holds. Authority can be narrowed on the way down and never
// widened, so no sequence of runtime calls can produce a privilege that was
// not in the config to begin with.
type Registry struct {
	mu  sync.RWMutex
	cfg *Config

	// state is everything that changes at runtime and must survive restart.
	state registryState
}

// registryState is the persisted half of the registry.
type registryState struct {
	// DynamicGroups maps an agent to memberships granted by its parent.
	DynamicGroups map[string][]string `json:"dynamic_groups,omitempty"`
	// Parents maps a child agent to the agent that registered it. Only a
	// parent may grant or revoke its child's memberships.
	Parents map[string]string `json:"parents,omitempty"`
	// ChildKeys holds the static public keys of agents registered at
	// runtime, which are not present in the config file.
	ChildKeys map[string]string `json:"child_keys,omitempty"`
	// Contactable maps an agent to the groups it accepts contact from. An
	// agent with no declaration accepts nothing: inbound access is opt-in,
	// so an agent that has not spoken cannot be reached by surprise.
	Contactable map[string][]string `json:"contactable,omitempty"`
	// Groups holds groups created at runtime by a group-admin. Config
	// groups are not copied here; the two namespaces stay separate so a
	// restart cannot silently revert a runtime edit.
	Groups map[string]GroupConfig `json:"groups,omitempty"`
}

// NewRegistry builds a registry over a validated config.
func NewRegistry(cfg *Config) *Registry {
	return &Registry{
		cfg: cfg,
		state: registryState{
			DynamicGroups: map[string][]string{},
			Parents:       map[string]string{},
			ChildKeys:     map[string]string{},
			Contactable:   map[string][]string{},
			Groups:        map[string]GroupConfig{},
		},
	}
}

// Known reports whether the gateway recognizes this agent at all.
func (r *Registry) Known(agentID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.knownLocked(agentID)
}

func (r *Registry) knownLocked(agentID string) bool {
	if _, ok := r.cfg.Agents[agentID]; ok {
		return true
	}
	_, ok := r.state.ChildKeys[agentID]
	return ok
}

// PublicKey returns an agent's static Noise public key, from config for
// configured agents or from runtime registration for children.
func (r *Registry) PublicKey(agentID string) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	encoded := ""
	if a, ok := r.cfg.Agents[agentID]; ok {
		encoded = a.PublicKey
	} else if k, ok := r.state.ChildKeys[agentID]; ok {
		encoded = k
	} else {
		return nil, fmt.Errorf("unknown agent %q", agentID)
	}
	return pub.ParsePublicKey(encoded)
}

// GroupsOf returns the agent's full membership set, base plus dynamic.
func (r *Registry) GroupsOf(agentID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.groupsOfLocked(agentID)
}

func (r *Registry) groupsOfLocked(agentID string) []string {
	set := map[string]bool{}
	if a, ok := r.cfg.Agents[agentID]; ok {
		for _, g := range a.Groups {
			set[g] = true
		}
	}
	for _, g := range r.state.DynamicGroups[agentID] {
		set[g] = true
	}
	out := make([]string, 0, len(set))
	for g := range set {
		out = append(out, g)
	}
	sort.Strings(out)
	return out
}

// HasGroup reports membership in one group.
func (r *Registry) HasGroup(agentID, group string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.hasGroupLocked(agentID, group)
}

func (r *Registry) hasGroupLocked(agentID, group string) bool {
	for _, g := range r.groupsOfLocked(agentID) {
		if g == group {
			return true
		}
	}
	return false
}

// HasCapability reports whether any of the agent's groups confers cap.
func (r *Registry) HasCapability(agentID, capability string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, g := range r.groupsOfLocked(agentID) {
		def, ok := r.groupLocked(g)
		if !ok {
			continue
		}
		for _, c := range def.Capabilities {
			if c == capability {
				return true
			}
		}
	}
	return false
}

// DeclareContactable records which groups may send messages to agentID.
//
// This is the agent controlling its own inbound access rather than the
// config hardcoding who may talk to whom. Declaring an undefined group is
// refused, because it would read as a working restriction while admitting
// nobody.
func (r *Registry) DeclareContactable(agentID string, groups []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.knownLocked(agentID) {
		return fmt.Errorf("unknown agent %q", agentID)
	}
	for _, g := range groups {
		if _, ok := r.groupLocked(g); !ok {
			return fmt.Errorf("undefined group %q", g)
		}
	}
	clean := append([]string(nil), groups...)
	sort.Strings(clean)
	r.state.Contactable[agentID] = clean
	return nil
}

// ContactableBy returns the groups the agent accepts contact from.
func (r *Registry) ContactableBy(agentID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]string(nil), r.state.Contactable[agentID]...)
}

// Authorize decides whether sender may deliver to recipient, and at what
// trust level.
//
// Delivery is permitted when the sender belongs to at least one group the
// recipient has declared it accepts. The returned trust is the highest trust
// among the groups that authorized this particular contact: a sender reaches
// the recipient with the authority of the relationship they actually share,
// not with every privilege the sender happens to hold elsewhere.
func (r *Registry) Authorize(sender, recipient string) (via []string, trust pub.Trust, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.knownLocked(sender) {
		return nil, "", fmt.Errorf("unknown agent %q", sender)
	}
	if !r.knownLocked(recipient) {
		return nil, "", fmt.Errorf("unknown agent %q", recipient)
	}

	accepted := map[string]bool{}
	for _, g := range r.state.Contactable[recipient] {
		accepted[g] = true
	}
	for _, g := range r.groupsOfLocked(sender) {
		if accepted[g] {
			via = append(via, g)
		}
	}
	if len(via) == 0 {
		return nil, "", fmt.Errorf("%s is not in any group %s accepts contact from", sender, recipient)
	}

	trust = pub.TrustExternal
	for _, g := range via {
		def, ok := r.groupLocked(g)
		if !ok {
			continue
		}
		if trustRank(def.TrustLevel()) > trustRank(trust) {
			trust = def.TrustLevel()
		}
	}
	return via, trust, nil
}

// trustRank orders trust levels so the highest can be selected.
func trustRank(t pub.Trust) int {
	switch t {
	case pub.TrustOwner:
		return 3
	case pub.TrustTrusted:
		return 2
	case pub.TrustExternal:
		return 1
	}
	return 0
}

// RegisterChild records a child agent's public key and its initial
// memberships, granted by the parent.
//
// Two rules apply, and both are refusals rather than silent trims. The
// parent may grant only groups it holds, so authority narrows. And the
// parent becomes the child's only administrator, so no third agent can later
// alter what this child may do.
func (r *Registry) RegisterChild(parentID, childID, publicKey string, groups []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.knownLocked(parentID) {
		return fmt.Errorf("unknown agent %q", parentID)
	}
	if childID == "" {
		return fmt.Errorf("child id is required")
	}
	if childID == parentID {
		return fmt.Errorf("agent %q cannot register itself as its own child", childID)
	}
	if _, ok := r.cfg.Agents[childID]; ok {
		return fmt.Errorf("agent %q is defined in config and cannot be re-registered", childID)
	}
	if existing, ok := r.state.Parents[childID]; ok && existing != parentID {
		return fmt.Errorf("agent %q is already registered by %q", childID, existing)
	}
	if _, err := pub.ParsePublicKey(publicKey); err != nil {
		return fmt.Errorf("child %q: %w", childID, err)
	}
	if err := r.checkGrantableLocked(parentID, groups); err != nil {
		return err
	}

	r.state.ChildKeys[childID] = publicKey
	r.state.Parents[childID] = parentID
	clean := append([]string(nil), groups...)
	sort.Strings(clean)
	r.state.DynamicGroups[childID] = clean
	return nil
}

// checkGrantableLocked enforces narrow-only inheritance.
func (r *Registry) checkGrantableLocked(parentID string, groups []string) error {
	for _, g := range groups {
		if _, ok := r.groupLocked(g); !ok {
			return fmt.Errorf("undefined group %q", g)
		}
		if !r.hasGroupLocked(parentID, g) {
			return fmt.Errorf("%s cannot grant group %q it does not hold", parentID, g)
		}
	}
	return nil
}

// GrantMembership adds a group to a child the caller registered.
func (r *Registry) GrantMembership(parentID, childID, group string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state.Parents[childID] != parentID {
		return fmt.Errorf("%s is not the parent of %s", parentID, childID)
	}
	if err := r.checkGrantableLocked(parentID, []string{group}); err != nil {
		return err
	}
	for _, g := range r.state.DynamicGroups[childID] {
		if g == group {
			return nil
		}
	}
	r.state.DynamicGroups[childID] = append(r.state.DynamicGroups[childID], group)
	sort.Strings(r.state.DynamicGroups[childID])
	return nil
}

// RevokeMembership removes a dynamic group from a child the caller
// registered.
func (r *Registry) RevokeMembership(parentID, childID, group string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state.Parents[childID] != parentID {
		return fmt.Errorf("%s is not the parent of %s", parentID, childID)
	}
	kept := make([]string, 0, len(r.state.DynamicGroups[childID]))
	for _, g := range r.state.DynamicGroups[childID] {
		if g != group {
			kept = append(kept, g)
		}
	}
	r.state.DynamicGroups[childID] = kept
	return nil
}

// ParentOf returns the agent that registered childID, if any.
func (r *Registry) ParentOf(childID string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.state.Parents[childID]
	return p, ok
}

// AgentSummary describes an agent for listing.
type AgentSummary struct {
	AgentID     string   `json:"agent_id"`
	DisplayName string   `json:"display_name"`
	Groups      []string `json:"groups"`
	Contactable []string `json:"contactable"`
	Parent      string   `json:"parent,omitempty"`
}

// ListAgents returns every known agent, configured and runtime-registered.
func (r *Registry) ListAgents() []AgentSummary {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := map[string]bool{}
	for id := range r.cfg.Agents {
		ids[id] = true
	}
	for id := range r.state.ChildKeys {
		ids[id] = true
	}
	ordered := make([]string, 0, len(ids))
	for id := range ids {
		ordered = append(ordered, id)
	}
	sort.Strings(ordered)

	out := make([]AgentSummary, 0, len(ordered))
	for _, id := range ordered {
		out = append(out, AgentSummary{
			AgentID:     id,
			DisplayName: r.cfg.Agents[id].DisplayName,
			Groups:      r.groupsOfLocked(id),
			Contactable: append([]string(nil), r.state.Contactable[id]...),
			Parent:      r.state.Parents[id],
		})
	}
	return out
}
