package gateway

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Save writes the registry's runtime state to path atomically.
//
// Only the runtime half is persisted. Config is never written back, so a
// restart always re-reads the operator's file as the source of base
// authority; a bug in the gateway cannot rewrite its own permissions.
func (r *Registry) Save(path string) error {
	r.mu.RLock()
	data, err := json.MarshalIndent(r.state, "", "  ")
	r.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("marshal registry state: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	// Write to a temporary file and rename, so a crash midway leaves the
	// previous state intact rather than a truncated file that will not
	// parse on the next start.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("write registry state: %w", err)
	}
	return os.Rename(tmp, path)
}

// Load restores runtime state written by Save. A missing file is not an
// error: a gateway starting for the first time has no dynamic state.
func (r *Registry) Load(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read registry state %s: %w", path, err)
	}

	var state registryState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("parse registry state %s: %w", path, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if state.DynamicGroups == nil {
		state.DynamicGroups = map[string][]string{}
	}
	if state.Parents == nil {
		state.Parents = map[string]string{}
	}
	if state.ChildKeys == nil {
		state.ChildKeys = map[string]string{}
	}
	if state.Contactable == nil {
		state.Contactable = map[string][]string{}
	}
	if state.Groups == nil {
		state.Groups = map[string]GroupConfig{}
	}
	r.state = state
	return nil
}
