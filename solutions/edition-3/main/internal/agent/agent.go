// Package agent owns configuration and the completed conversation, not HTTP.
package agent

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"ensemble/internal/common"
)

type agent struct {
	parent  common.Ensemble
	config  common.Config
	engine  common.Engine
	history []common.Message
}

// New validates configuration before any request and retains the actual owner.
// The composition root attaches Engine before publishing this Agent.
func New(parent common.Ensemble, config common.Config) (*agent, error) {
	if parent == nil {
		return nil, errors.New("Agent requires an Ensemble parent")
	}
	if strings.TrimSpace(config.Model) == "" {
		return nil, errors.New("ANTHROPIC_MODEL is not set; ask GET /v1/models which models exist")
	}
	if strings.TrimSpace(config.APIKey) == "" {
		return nil, errors.New("ANTHROPIC_API_KEY is not set")
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://api.anthropic.com"
	}
	config.BaseURL = strings.TrimRight(config.BaseURL, "/")
	u, err := url.Parse(config.BaseURL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		// Never echo a malformed URL: it might contain embedded credentials.
		return nil, errors.New("invalid ANTHROPIC_BASE_URL")
	}
	return &agent{parent: parent, config: config}, nil
}

// SetEngine completes the root's one-time wiring; it cannot replace an owner.
func (a *agent) SetEngine(engine common.Engine) {
	if engine == nil || a.engine != nil || engine.Parent() != a {
		panic("invalid Agent/Engine wiring")
	}
	a.engine = engine
}

// Parent exposes the actual Ensemble rather than a copied service bundle.
func (a *agent) Parent() common.Ensemble { return a.parent }

// Config returns a value copy of the Agent-owned configuration.
func (a *agent) Config() common.Config { return a.config }

// Engine exposes the owned transport and its parent relationship.
func (a *agent) Engine() common.Engine { return a.engine }

// History returns a defensive copy of completed turns.
func (a *agent) History() []common.Message { return append([]common.Message(nil), a.history...) }

// Usage delegates accounting to its single authoritative owner.
func (a *agent) Usage() common.Usage { return a.engine.Usage() }

// Ask stages the user turn so a failed exchange cannot poison alternation.
// No retry is hidden here: one question makes at most one paid request.
func (a *agent) Ask(ctx context.Context, question string) (string, error) {
	if strings.TrimSpace(question) == "" {
		return "", errors.New("user message is empty")
	}
	pending := append(a.history, common.Message{Role: "user", Content: question})
	reply, err := a.engine.Send(ctx, pending)
	if err != nil {
		return "", err
	}
	a.history = append(pending, common.Message{Role: "assistant", Content: reply})
	return reply, nil
}
