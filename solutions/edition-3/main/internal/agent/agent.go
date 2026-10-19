// Package agent owns turn orchestration and configuration, not HTTP or replay.
package agent

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"ensemble/internal/common"
)

type agent struct{ privateState }

type privateState struct {
	parent  common.Ensemble
	config  common.Config
	engine  common.Engine
	history common.History
	tools   common.ToolRegistry
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
	if config.Tools != nil {
		config.Tools = append([]string{}, config.Tools...)
	}
	return &agent{privateState: privateState{parent: parent, config: config}}, nil
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
func (a *agent) Config() common.Config {
	c := a.config
	if c.Tools != nil {
		c.Tools = append([]string{}, c.Tools...)
	}
	return c
}

func (a *agent) SetTools(t common.ToolRegistry) {
	if t == nil || a.tools != nil || t.Parent() != a {
		panic("invalid Agent/Tools wiring")
	}
	a.tools = t
}
func (a *agent) Tools() common.ToolRegistry { return a.tools }

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
			hasCalls := false
			for _, p := range e.Response.Parts {
				if p.Type == "tool_call" {
					hasCalls = true
				}
			}
			if hasCalls {
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

// Ask adds each fully returned response before dispatching its ordered calls.
// A budget counts batches, not elapsed time, stop reasons, or usage estimates.
func (a *agent) Ask(ctx context.Context, question string) (string, error) {
	if strings.TrimSpace(question) == "" {
		return "", errors.New("user message is empty")
	}
	if err := a.history.Append(common.Event{Type: common.MessageReceived, Message: &common.MessageData{Actor: common.Human, Parts: common.PartList{common.Text(question)}}}); err != nil {
		return "", err
	}
	const maxBatches = 16
	for batches := 0; ; batches++ {
		projection := a.history.Context()
		// Validate before declaring a request sent. This snapshot delivers ephemera
		// once even though RequestSent consumes it from the retained projection.
		if _, err := a.engine.Render(projection); err != nil {
			return "", err
		}
		if err := a.history.Append(common.Event{Type: common.RequestSent, Request: &common.RequestData{To: a.config.Target()}}); err != nil {
			return "", err
		}
		events, exchangeErr := a.engine.Exchange(ctx, projection)
		// Missing live IDs need durable identities, unlike offline render-only IDs.
		idErr := a.identifyCalls(events)
		var parts common.PartList
		for _, event := range events {
			if err := a.history.Append(event); err != nil {
				return "", err
			}
			if event.Type == common.ResponseEnded {
				parts = event.Response.Parts
			}
		}
		if exchangeErr != nil {
			return "", exchangeErr
		}
		if idErr != nil {
			return "", a.recordError(idErr.Error())
		}
		var calls common.PartList
		var reply strings.Builder
		explained := false
		for _, p := range parts {
			if p.Type == "text" {
				reply.WriteString(p.Text)
				if len(calls) == 0 && strings.TrimSpace(p.Text) != "" {
					explained = true
				}
			}
			if p.Type == "tool_call" {
				calls = append(calls, p)
			}
		}
		if len(calls) == 0 {
			return reply.String(), nil
		}
		refusal := ""
		if a.config.RequireToolExplanation && !explained {
			refusal = "No tools ran: provide nonblank ordinary assistant text explaining what will run and why before the first tool call, then retry."
		}
		if batches == maxBatches {
			refusal = "No tools ran: the limit of 16 tool batches has been reached."
		}
		// Decide whole-batch refusal before invoking even the first handler.
		for _, call := range calls {
			if err := a.history.Append(common.Event{Type: common.ToolCalled, Tool: &common.ToolData{CallID: call.CallID, Name: call.Name, Args: call.Args}}); err != nil {
				return "", err
			}
			result, failed := common.PartList{common.Text(refusal)}, true
			if refusal == "" {
				result, failed = a.tools.Dispatch(ctx, call.Name, call.Args)
			}
			if err := a.history.Append(common.Event{Type: common.ToolReturned, Tool: &common.ToolData{CallID: call.CallID, Name: call.Name, Parts: result, IsError: failed}}); err != nil {
				return "", err
			}
		}
		// The request after batch sixteen consumed its actual results. A seventeenth
		// request for tools is paired with refusals, logged, and terminates locally.
		if batches == maxBatches {
			return "", a.recordError("tool batch limit reached (16)")
		}
	}
}

func (a *agent) recordError(message string) error {
	err := a.history.Append(common.Event{Type: common.ErrorOccurred, Error: &common.ErrorData{Message: message}})
	if err != nil {
		return err
	}
	return errors.New(message)
}

func (a *agent) identifyCalls(events []common.Event) error {
	used := map[string]bool{}
	prior := a.history.Events()
	for _, e := range prior {
		if e.Response != nil {
			for _, p := range e.Response.Parts {
				if p.Type == "tool_call" && p.CallID != "" {
					used[p.CallID] = true
				}
			}
		}
	}
	// Reserve supplied IDs before synthesizing missing ones, preserving provenance
	// and all provider IDs. Duplicate supplied IDs cannot safely be dispatched.
	for _, e := range events {
		if e.Type != common.ResponseEnded {
			continue
		}
		for _, p := range e.Response.Parts {
			if p.Type != "tool_call" || p.CallID == "" {
				continue
			}
			if used[p.CallID] {
				return fmt.Errorf("duplicate tool call ID")
			}
			used[p.CallID] = true
		}
	}
	for i, e := range events {
		if e.Type != common.ResponseEnded {
			continue
		}
		for j := range e.Response.Parts {
			p := &e.Response.Parts[j]
			if p.Type != "tool_call" || p.CallID != "" {
				continue
			}
			base := fmt.Sprintf("live_call_%d_%d_%d", len(prior), i, j)
			id := base
			for suffix := 1; used[id]; suffix++ {
				id = fmt.Sprintf("%s_%d", base, suffix)
			}
			p.CallID = id
			used[id] = true
		}
	}
	return nil
}
