package gateway

import (
	"fmt"

	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// Group administration is itself a capability.
//
// Editing groups is how every other permission in the system is defined, so
// the right to do it has to be protected by the same mechanism it protects.
// It is conferred by the group-admin capability, which a config-defined root
// group holds. The system bootstraps from the operator's file and is
// self-protecting from that point on: there is no runtime path by which an
// agent without group-admin acquires it.
//
// Runtime groups live beside config groups, never on top of them. A group
// defined in the operator's file cannot be edited or deleted over the wire,
// for two reasons. The config is re-read at every start, so a runtime edit
// to a config group would silently revert and leave an operator debugging a
// permission that works until the next restart. And the operator's file is
// the one place authority is supposed to be auditable by reading it.

// requireGroupAdmin returns an error unless the caller may administer groups.
func (r *Registry) requireGroupAdmin(callerID string) error {
	if !r.HasCapability(callerID, pub.CapGroupAdmin) {
		return fmt.Errorf("%s lacks the %s capability", callerID, pub.CapGroupAdmin)
	}
	return nil
}

// CreateGroup defines a new runtime group.
func (r *Registry) CreateGroup(callerID, name string, group GroupConfig) error {
	if err := r.requireGroupAdmin(callerID); err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("group name is required")
	}
	if group.Trust != "" && !group.Trust.Valid() {
		return fmt.Errorf("unknown trust %q", group.Trust)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.groupLocked(name); exists {
		return fmt.Errorf("group %q already exists", name)
	}
	if r.state.Groups == nil {
		r.state.Groups = map[string]GroupConfig{}
	}
	r.state.Groups[name] = group
	return nil
}

// EditGroup replaces a runtime group's description, capabilities and trust.
func (r *Registry) EditGroup(callerID, name string, group GroupConfig) error {
	if err := r.requireGroupAdmin(callerID); err != nil {
		return err
	}
	if group.Trust != "" && !group.Trust.Valid() {
		return fmt.Errorf("unknown trust %q", group.Trust)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, isConfig := r.cfg.Groups[name]; isConfig {
		return fmt.Errorf("group %q is defined in config and cannot be edited at runtime", name)
	}
	if _, exists := r.state.Groups[name]; !exists {
		return fmt.Errorf("unknown group %q", name)
	}
	r.state.Groups[name] = group
	return nil
}

// DeleteGroup removes a runtime group and every membership and
// contactability declaration that referred to it.
//
// Leaving dangling references would be worse than refusing: a membership in
// a deleted group confers nothing but still reads as access when an operator
// lists an agent.
func (r *Registry) DeleteGroup(callerID, name string) error {
	if err := r.requireGroupAdmin(callerID); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, isConfig := r.cfg.Groups[name]; isConfig {
		return fmt.Errorf("group %q is defined in config and cannot be deleted at runtime", name)
	}
	if _, exists := r.state.Groups[name]; !exists {
		return fmt.Errorf("unknown group %q", name)
	}

	delete(r.state.Groups, name)
	for id, groups := range r.state.DynamicGroups {
		r.state.DynamicGroups[id] = without(groups, name)
	}
	for id, groups := range r.state.Contactable {
		r.state.Contactable[id] = without(groups, name)
	}
	return nil
}

// Groups returns every defined group, config and runtime.
func (r *Registry) Groups() map[string]GroupConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]GroupConfig, len(r.cfg.Groups)+len(r.state.Groups))
	for name, g := range r.cfg.Groups {
		out[name] = g
	}
	for name, g := range r.state.Groups {
		out[name] = g
	}
	return out
}

// groupLocked resolves a group by name from either namespace.
func (r *Registry) groupLocked(name string) (GroupConfig, bool) {
	if g, ok := r.cfg.Groups[name]; ok {
		return g, true
	}
	g, ok := r.state.Groups[name]
	return g, ok
}

func without(list []string, drop string) []string {
	kept := make([]string, 0, len(list))
	for _, v := range list {
		if v != drop {
			kept = append(kept, v)
		}
	}
	return kept
}
