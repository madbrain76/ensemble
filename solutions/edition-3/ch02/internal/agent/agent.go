// Package agent owns turn orchestration and configuration, not HTTP or replay.
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
	history common.History
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

	if config.Vendor == 0 {
		config.Vendor = common.Anthropic
	}
	switch config.Vendor {
	case common.Anthropic:
		if config.Surface == 0 {
			config.Surface = common.Messages
		}
		if config.BaseURL == "" {
			config.BaseURL = "https://api.anthropic.com"
		}
	case common.OpenAI:
		if config.Surface == 0 {
			config.Surface = common.ChatCompletions
		}
		if config.BaseURL == "" {
			config.BaseURL = "https://api.openai.com"
		}
	case common.Gemini:
		if config.Surface == 0 {
			config.Surface = common.GenerateContent
		}
		if config.BaseURL == "" {
			config.BaseURL = "https://generativelanguage.googleapis.com"
		}
	default:
		return nil, errors.New("unknown vendor")
	}
	if (config.Vendor == common.Anthropic && config.Surface != common.Messages) || (config.Vendor == common.OpenAI && config.Surface != common.ChatCompletions) || (config.Vendor == common.Gemini && config.Surface != common.GenerateContent) {
		return nil, errors.New("unsupported vendor surface")
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

// SetHistory completes the other child edge at the composition root.
func (a *agent) SetHistory(h common.History) {
	if h == nil || a.history != nil || h.Parent() != a {
		panic("invalid Agent/History wiring")
	}
	a.history = h
}
func (a *agent) Log() common.History     { return a.history }
func (a *agent) Render() ([]byte, error) { return a.engine.Render(a.history.Context()) }
func (a *agent) ParseResponse(status int, b []byte) ([]common.Event, error) {
	return a.engine.Parse(status, b)
}

// History preserves Chapter 1's completed-text view, derived from the log.
func (a *agent) History() []common.Message {
	var messages, pending []common.Message
	events := a.history.Events()
	for i, e := range events {
		switch e.Type {
		case common.MessageReceived:
			if e.Message.Actor == common.Human {
				var b strings.Builder
				for _, p := range e.Message.Parts {
					if p.Type == "text" {
						b.WriteString(p.Text)
					}
				}
				pending = append(pending, common.Message{Role: "user", Content: b.String()})
			}
		case common.ResponseEnded:
			if i+1 < len(events) && events[i+1].Type == common.ErrorOccurred {
				pending = nil
				continue
			}
			var b strings.Builder
			for _, p := range e.Response.Parts {
				if p.Type == "text" {
					b.WriteString(p.Text)
				}
			}
			if strings.TrimSpace(b.String()) != "" {
				messages = append(messages, pending...)
				messages = append(messages, common.Message{Role: "assistant", Content: b.String()})
			}
			pending = nil
		case common.ErrorOccurred:
			pending = nil
		}
	}
	return messages
}
func (a *agent) Usage() common.Usage { return a.engine.Usage() }

// Ask adds causal facts through History only. No hidden retry or tool execution.
func (a *agent) Ask(ctx context.Context, question string) (string, error) {
	if strings.TrimSpace(question) == "" {
		return "", errors.New("user message is empty")
	}
	if err := a.history.Append(common.Event{Type: common.MessageReceived, Message: &common.MessageData{Actor: common.Human, Parts: common.PartList{common.Text(question)}}}); err != nil {
		return "", err
	}
	projection := a.history.Context()
	// Validate the disposable request before declaring it sent; ephemera is
	// consumed by RequestSent, but this snapshot still delivers it exactly once.
	if _, err := a.engine.Render(projection); err != nil {
		return "", err
	}
	if err := a.history.Append(common.Event{Type: common.RequestSent, Request: &common.RequestData{To: a.config.Target()}}); err != nil {
		return "", err
	}
	events, err := a.engine.Exchange(ctx, projection)
	var reply strings.Builder
	for _, event := range events {
		if appendErr := a.history.Append(event); appendErr != nil {
			return "", appendErr
		}
		if event.Type == common.ResponseEnded {
			for _, p := range event.Response.Parts {
				if p.Type == "text" {
					reply.WriteString(p.Text)
				}
			}
		}
	}
	return reply.String(), err
}
