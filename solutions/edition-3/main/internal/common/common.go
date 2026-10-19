// Package common holds the vocabulary shared by the owners, not runtime state.
package common

import (
	"context"
	"log"
)

// Config belongs to one Agent. Model names must come from provider discovery.
type Config struct {
	// BaseURL selects the provider or a local HTTP peer; empty means Anthropic.
	BaseURL string
	// APIKey authenticates requests and must never appear in diagnostics.
	APIKey string
	// Model is required; there is deliberately no remembered model default.
	Model string
}

// Message is a retained user or assistant turn, sent verbatim on later requests.
type Message struct {
	// Role alternates between user and assistant, beginning with user.
	Role string `json:"role"`
	// Content is nonempty text; request content is not a response block array.
	Content string `json:"content"`
}

// Usage is measured cumulative token consumption, including unusable answers.
type Usage struct {
	// Input counts input tokens reported by the provider.
	Input int64 `json:"input"`
	// Output counts output tokens reported by the provider.
	Output int64 `json:"output"`
}

// Ensemble is the Agent's actual owner and the source of application services.
type Ensemble interface {
	// Logger returns this application's diagnostic destination.
	Logger() *log.Logger
}

// Agent owns one synchronous conversation. Calls on one Agent must be serialized.
type Agent interface {
	// Parent returns the Ensemble which owns this conversation.
	Parent() Ensemble
	// Config returns a value copy, leaving configuration with its owner.
	Config() Config
	// Engine returns this Agent's transport and accounting owner.
	Engine() Engine
	// Ask completes one exchange, retaining history only on success.
	Ask(context.Context, string) (string, error)
	// History returns a copy so callers cannot corrupt retained turns.
	History() []Message
	// Usage reads measured totals from Engine rather than keeping a second copy.
	Usage() Usage
}

// Engine exchanges messages; it never retains the conversation itself.
type Engine interface {
	// Parent returns the Agent which owns transport configuration.
	Parent() Agent
	// Send makes one HTTP request, with cancellation controlled by the caller.
	Send(context.Context, []Message) (string, error)
	// Usage returns cumulative measured totals, not estimates.
	Usage() Usage
}
