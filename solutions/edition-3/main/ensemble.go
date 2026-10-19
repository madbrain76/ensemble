// Package ensemble is the public composition root for Chapter 1's conversation.
package ensemble

import (
	"errors"
	"io"
	"log"
	"sync"

	"ensemble/internal/agent"
	"ensemble/internal/common"
	"ensemble/internal/engine"
	"ensemble/internal/history"
	"ensemble/internal/jobs"
	"ensemble/internal/tool"
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
type Ensemble struct{ privateState }

type privateState struct {
	mu         sync.Mutex
	closed     bool
	nextHandle int
	jobs       map[int]common.Job
	logger     *log.Logger
	agents     []common.Agent
}

// New creates an independent application. A nil logger discards diagnostics.
func New(logger *log.Logger) *Ensemble {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &Ensemble{privateState: privateState{logger: logger, jobs: make(map[int]common.Job)}}
}

// Logger supplies the application-owned logger to descendants through parents.
func (e *Ensemble) Logger() *log.Logger { return e.logger }

// NewAgent wires both ownership directions before exposing the conversation.
// Neither implementation imports the other; only this root connects them.
func (e *Ensemble) NewAgent(config Config) (Agent, error) {
	if e == nil {
		return nil, errors.New("Agent requires an Ensemble parent")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return nil, errors.New("Ensemble is shut down")
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
	a.SetHistory(history.New(a))
	a.SetJobs(jobs.New(a))
	registry, err := tool.New(a)
	if err != nil {
		return nil, err
	}
	a.SetTools(registry)
	e.agents = append(e.agents, a)
	return a, nil
}

// Public aliases expose the common contracts without exporting implementation owners.
type Event = common.Event
type Context = common.Context
type History = common.History
type Part = common.Part
type PartList = common.PartList
type MessageData = common.MessageData
type Provenance = common.Provenance
type Vendor = common.Vendor
type Surface = common.Surface

const (
	Anthropic       = common.Anthropic
	OpenAI          = common.OpenAI
	Gemini          = common.Gemini
	Messages        = common.Messages
	ChatCompletions = common.ChatCompletions
	GenerateContent = common.GenerateContent
	MessageReceived = common.MessageReceived
	System          = common.System
)

// Tool contracts expose declarations and dispatch without implementation state.
type ToolRegistry = common.ToolRegistry
type ToolDefinition = common.ToolDefinition

// RegisterJob atomically admits a job and allocates a process-local Ensemble handle.
func (e *Ensemble) RegisterJob(j common.Job) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return 0, errors.New("Ensemble is shut down")
	}
	e.nextHandle++
	e.jobs[e.nextHandle] = j
	return e.nextHandle, nil
}
func (e *Ensemble) LookupJob(handle int) (common.Job, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	j, ok := e.jobs[handle]
	return j, ok
}
func (e *Ensemble) JobTable() []common.Job {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]common.Job, 0, len(e.jobs))
	for n := 1; n <= e.nextHandle; n++ {
		out = append(out, e.jobs[n])
	}
	return out
}

// Shutdown stops admission before killing jobs, and is safe to repeat.
func (e *Ensemble) Shutdown() error {
	e.mu.Lock()
	e.closed = true
	agents := append([]common.Agent(nil), e.agents...)
	e.mu.Unlock()
	var errs []error
	for _, a := range agents {
		errs = append(errs, a.Shutdown())
	}
	return errors.Join(errs...)
}

type ToolHandler = common.ToolHandler
type ToolCall = common.ToolCall
type Jobs = common.Jobs
type Job = common.Job
type JobData = common.JobData
type JobReport = common.JobReport
type JobStatus = common.JobStatus
type ToolLimits = common.ToolLimits
type Ref = common.Ref

const (
	RefHandle       = common.RefHandle
	JobRunning      = common.JobRunning
	JobDone         = common.JobDone
	JobKilledStatus = common.JobKilledStatus
	JobKilled       = common.JobKilled
)

func DefaultToolLimits() ToolLimits { return common.DefaultToolLimits() }
