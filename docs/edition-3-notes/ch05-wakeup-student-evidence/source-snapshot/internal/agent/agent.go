// Package agent owns an Agent's configuration, History and Engine.
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ensemble/internal/common"
)

// Agent is constructed through Ensemble, never as an alternate application root.
type Agent struct {
	parent  common.Ensemble
	config  common.Config
	engine  common.Engine
	history common.History
	// tools are the operations this Agent owns and dispatches.
	tools        []common.Tool
	pending      *common.WaitOptions
	mu           sync.Mutex
	once         sync.Once
	mailMu       sync.Mutex
	mail         []any
	wake         chan struct{}
	done         chan struct{}
	stopAck      chan struct{}
	accepting    bool
	finishOnce   sync.Once
	workers      sync.WaitGroup
	observations chan common.HubEvent
	observers    []common.Observer
	stateMu      sync.RWMutex
	state        common.TurnState
	active       *turn
	queue        []*request
	generation   uint64
	stopping     bool
	shutdownErr  error
}

// New establishes every parent relationship before exposing the Agent.
// Credentials are needed at Send, not for offline log inspection and rendering.
func New(parent common.Ensemble, config common.Config) (*Agent, error) {
	if parent == nil {
		return nil, fmt.Errorf("agent requires an ensemble parent")
	}
	if strings.TrimSpace(config.System) == "" {
		return nil, fmt.Errorf("system prompt is required")
	}
	if config.MaxTokens <= 0 {
		return nil, fmt.Errorf("max tokens must be positive")
	}
	if config.Vendor == 0 {
		config.Vendor = common.Anthropic
	}
	if config.Workspace == "" {
		var err error
		config.Workspace, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	var pathErr error
	config.Workspace, pathErr = filepath.Abs(config.Workspace)
	if pathErr != nil {
		return nil, pathErr
	}
	if config.DataDir == "" {
		config.DataDir = config.Workspace
	}
	config.DataDir, pathErr = filepath.Abs(config.DataDir)
	if pathErr != nil {
		return nil, pathErr
	}
	if config.ToolBatchLimit < 0 {
		return nil, fmt.Errorf("tool batch limit cannot be negative")
	}
	a := &Agent{parent: parent, config: config, accepting: true, wake: make(chan struct{}, 1), done: make(chan struct{}), stopAck: make(chan struct{}), observations: make(chan common.HubEvent, 64), state: common.Idle}

	return a, nil
}

// AttachTools is composition-root assembly, before publication.
func AttachTools(a *Agent, tools []common.Tool) { a.tools = append([]common.Tool(nil), tools...) }

// Tools returns the Agent-owned operation registry as a snapshot.
func (a *Agent) Tools() []common.Tool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]common.Tool(nil), a.tools...)
}

// Parent exposes the required root relationship.
func (a *Agent) Parent() common.Ensemble { return a.parent }

// Config returns the Agent-owned fixed configuration.
func (a *Agent) Config() common.Config { return a.config }

// Engine exposes the concrete owned Engine through its interface.
func (a *Agent) Engine() common.Engine { return a.engine }

// Log exposes the concrete owned History through its interface.
func (a *Agent) Log() common.History { return a.history }

// AttachChildren is composition-root assembly, before the Agent is published.
// Both children already refer back to this Agent; callers must attach only once.
func AttachChildren(a *Agent, engine common.Engine, history common.History) {
	a.engine, a.history = engine, history
}

// History retains the Chapter 1 text view; events are authoritative storage.
func (a *Agent) History() []common.Message {
	var messages []common.Message
	var pending string
	for _, event := range a.history.Events() {
		if event.Type == "message_received" && event.Message.Actor == "human" {
			pending = ""
			for _, part := range event.Message.Parts {
				if part.Type == "text" {
					pending += part.Text
				}
			}
		}
		if event.Type == "response_ended" {
			hasCalls := false
			for _, part := range event.Response.Parts {
				if part.Type == "tool_call" {
					hasCalls = true
				}
			}
			if hasCalls {
				continue
			}
			var text strings.Builder
			for _, part := range event.Response.Parts {
				if part.Type == "text" {
					text.WriteString(part.Text)
				}
			}
			messages = append(messages, common.Message{Role: "user", Content: pending}, common.Message{Role: "assistant", Content: text.String()})
		}
	}
	return messages
}

// RegisterTool adds an operation without replacing the builtin registrations.
func (a *Agent) RegisterTool(def common.ToolDefinition, handler common.ToolHandler) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	var schema struct {
		Type string `json:"type"`
	}
	if strings.TrimSpace(def.Name) == "" || strings.TrimSpace(def.Description) == "" || handler == nil || json.Unmarshal(def.Parameters, &schema) != nil || schema.Type != "object" {
		return fmt.Errorf("tool requires name, parameters and handler")
	}
	for _, tool := range a.tools {
		if tool.Definition().Name == def.Name {
			return fmt.Errorf("tool %q already registered", def.Name)
		}
	}
	def.Parameters = append(json.RawMessage(nil), def.Parameters...)
	a.tools = append(a.tools, &registered{parent: a, def: def, handler: handler})
	return nil
}

type registered struct {
	parent  common.Agent
	def     common.ToolDefinition
	handler common.ToolHandler
}

func (t *registered) Parent() common.Agent { return t.parent }
func (t *registered) Definition() common.ToolDefinition {
	d := t.def
	d.Parameters = append(json.RawMessage(nil), d.Parameters...)
	return d
}
func (t *registered) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	j, err := t.parent.Parent().Jobs().Create(t.parent)
	if err != nil {
		return "", err
	}
	j.Start(ctx, t, args)
	return j.Wait(ctx, defaults())
}
func (t *registered) Start(ctx context.Context, j common.Job, args json.RawMessage) (string, bool, error) {
	out, err := t.handler(common.ToolContext{Context: ctx, Agent: t.parent, Job: j}, args)
	return out, false, err
}
