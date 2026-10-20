// Package common holds the vocabulary shared by the owners, not runtime state.
package common

import (
	"context"
	"io"
	"log"
)

// Config belongs to one Agent. Model names must come from provider discovery.
type Config struct {
	// BaseDir resolves relative tool paths; DataDir retains job output.
	// Empty BaseDir uses the launch directory; empty DataDir creates an Agent-local directory.
	BaseDir string
	DataDir string
	// BaseURL selects the provider or a local HTTP peer; empty means Anthropic.
	BaseURL       string
	Vendor        Vendor
	Surface       Surface
	ResolvedModel string
	// APIKey authenticates requests and must never appear in diagnostics.
	APIKey string
	// Model is required; there is deliberately no remembered model default.
	Model string
	// Tools selects enabled built-ins: nil enables all ten; non-nil empty disables all.
	Tools []string
	// RequireToolExplanation selects this model's before-effects explanation policy.
	RequireToolExplanation bool
	LogPath                string
	MaxToolRounds          int
	SystemPrompt           string
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
	Output     int64 `json:"output"`
	CacheWrite int64 `json:"cache_write"`
	CacheRead  int64 `json:"cache_read"`
}

// Ensemble is the Agent's actual owner and the source of application services.
type Ensemble interface {
	// Logger returns this application's diagnostic destination.
	Logger() *log.Logger
	Publish(Agent, Observation)
	RegisterJob(Job) (int, error)
	LookupJob(int) (Job, bool)
	JobTable() []Job
	Shutdown() error
}

// Agent owns one mailbox and one conversation actor.
type Agent interface {
	// Parent returns the Ensemble which owns this conversation.
	Parent() Ensemble
	// Config returns a value copy, leaving configuration with its owner.
	Config() Config
	// Engine returns this Agent's transport and accounting owner.
	Engine() Engine
	// Tools is this Agent's independently owned synchronous registry.
	Tools() ToolRegistry
	Jobs() Jobs
	Shutdown() error
	Post(Inbound) error
	Submit(context.Context, string) <-chan Result
	// SubmitLive starts a turn if idle, or adds a hint to the current turn.
	// Joined callers receive the same reliable result and TurnID at completion.
	SubmitLive(context.Context, string) <-chan Result
	Subscribe(Observer)
	Wait(context.Context, func(Observation) bool) (Observation, error)
	ReportJob(JobReport)
	InvokeTool(context.Context, Part, string) error
	// Ask records every attempt; the compatibility History view contains completed text pairs.
	Ask(context.Context, string) (string, error)
	// History returns a copy so callers cannot corrupt retained turns.
	History() []Message
	Log() History
	Render() ([]byte, error)
	ParseResponse(int, []byte) ([]Event, error)
	// Usage reads measured totals from Engine rather than keeping a second copy.
	Usage() Usage
}

// Engine exchanges messages; it never retains the conversation itself.
type Engine interface {
	// Parent returns the Agent which owns transport configuration.
	Parent() Agent
	// Send makes one HTTP request, with cancellation controlled by the caller.
	Send(context.Context, []Message) (string, error)
	Exchange(context.Context, Context) ([]Event, error)
	ExchangeRequest(context.Context, []byte) ([]Event, error)
	Render(Context) ([]byte, error)
	Parse(int, []byte) ([]Event, error)
	// Usage returns cumulative measured totals, not estimates.
	Usage() Usage
}

// History owns causal facts and their projection; its parent is the real Agent.
type History interface {
	Parent() Agent
	Append(Event) error
	Events() []Event
	Context() Context
	Load(io.Reader) error
	Dump(io.Writer) error
}
