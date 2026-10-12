// Package gateway implements the gateway server: the edge service that sits
// at the trust boundary and routes messages between agents.
//
// It is private to the gateway binary. The client half that other programs
// import lives in the public agent/gateway package, which this package does
// not import in the other direction.
package gateway

import (
	"encoding/json"
	"fmt"
	"os"

	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// Config is the gateway's static configuration: the groups that exist, the
// agents that may connect, and the memberships they start with.
//
// The invariant this file protects is that the config declares the base of
// all authority. Nothing an agent says at runtime can introduce a group or
// grant itself a membership it was not given.
type Config struct {
	Gateway ServerConfig           `json:"gateway"`
	Agents  map[string]AgentConfig `json:"agents"`
	Groups  map[string]GroupConfig `json:"groups"`
}

// ServerConfig holds process-level settings.
type ServerConfig struct {
	// Listen is the bind address, such as ":8093".
	Listen string `json:"listen"`
	// DataDir holds runtime state: dynamic memberships and queued messages.
	DataDir string `json:"data_dir"`
}

// AgentConfig identifies one agent permitted to connect.
type AgentConfig struct {
	DisplayName string `json:"display_name"`
	// PublicKey is the agent's static Noise key, "x25519:<base64>".
	PublicKey string `json:"public_key"`
	// Groups are the agent's base memberships, the seed of its authority.
	Groups []string `json:"groups,omitempty"`
}

// GroupConfig describes one group.
type GroupConfig struct {
	Description string `json:"description"`
	// Capabilities are the privileged operations members may perform.
	Capabilities []string `json:"capabilities,omitempty"`
	// Trust is the trust level stamped on messages authorized by this
	// group. An empty value means external, the least privileged reading,
	// so that forgetting to declare trust cannot accidentally grant it.
	Trust pub.Trust `json:"trust,omitempty"`
}

// TrustLevel returns the group's trust, defaulting to external.
func (g GroupConfig) TrustLevel() pub.Trust {
	if g.Trust.Valid() {
		return g.Trust
	}
	return pub.TrustExternal
}

// LoadConfig reads and validates a gateway config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read gateway config %s: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse gateway config %s: %w", path, err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("gateway config %s: %w", path, err)
	}
	return &cfg, nil
}

// Validate rejects a config that would behave surprisingly at runtime.
//
// Every failure here is one that would otherwise surface much later as a
// handshake that will not complete or a message that silently will not
// route. Refusing to start with a clear reason is the cheaper outcome.
func (c *Config) Validate() error {
	if c.Gateway.Listen == "" {
		return fmt.Errorf("gateway.listen is required")
	}
	if c.Gateway.DataDir == "" {
		return fmt.Errorf("gateway.data_dir is required")
	}
	for name, group := range c.Groups {
		if name == "" {
			return fmt.Errorf("a group has an empty name")
		}
		if group.Trust != "" && !group.Trust.Valid() {
			return fmt.Errorf("group %q has unknown trust %q", name, group.Trust)
		}
	}
	for id, agent := range c.Agents {
		if id == "" {
			return fmt.Errorf("an agent has an empty id")
		}
		if agent.PublicKey == "" {
			return fmt.Errorf("agent %q has no public_key", id)
		}
		if _, err := pub.ParsePublicKey(agent.PublicKey); err != nil {
			return fmt.Errorf("agent %q: %w", id, err)
		}
		// A membership in a group that does not exist is always a typo, and
		// one that would quietly confer nothing instead of what was meant.
		for _, g := range agent.Groups {
			if _, ok := c.Groups[g]; !ok {
				return fmt.Errorf("agent %q is a member of undefined group %q", id, g)
			}
		}
	}
	return nil
}
