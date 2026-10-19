// Package ensemble is the public composition root for Chapter 1's conversation.
package ensemble

import (
	"errors"
	"io"
	"log"

	"ensemble/internal/agent"
	"ensemble/internal/common"
	"ensemble/internal/engine"
)

// Config configures one conversation; APIKey must be treated as a secret.
type Config = common.Config

// Agent is the public synchronous conversation interface.
type Agent = common.Agent

// Message is one retained user or assistant turn.
type Message = common.Message

// Usage contains cumulative input and output token counts.
type Usage = common.Usage

// Ensemble owns its Agents and logger; it has no invented parent of its own.
// Agent creation, like conversation calls, is synchronous in this chapter.
type Ensemble struct {
	logger *log.Logger
	agents []common.Agent
}

// New creates an independent application. A nil logger discards diagnostics.
func New(logger *log.Logger) *Ensemble {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &Ensemble{logger: logger}
}

// Logger supplies the application-owned logger to descendants through parents.
func (e *Ensemble) Logger() *log.Logger { return e.logger }

// NewAgent wires both ownership directions before exposing the conversation.
// Neither implementation imports the other; only this root connects them.
func (e *Ensemble) NewAgent(config Config) (Agent, error) {
	if e == nil {
		return nil, errors.New("Agent requires an Ensemble parent")
	}
	a, err := agent.New(e, config)
	if err != nil {
		return nil, err
	}
	transport, err := engine.New(a)
	if err != nil {
		return nil, err
	}
	a.SetEngine(transport)
	e.agents = append(e.agents, a)
	return a, nil
}
