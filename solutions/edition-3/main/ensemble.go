// Package ensemble is the public composition root for Chapter 1's conversation.
package ensemble

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
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

// Agent is a conversation actor with synchronous and asynchronous adapters.
type Agent = common.Agent

// Message is one retained user or assistant turn.
type Message = common.Message

// Usage contains cumulative input and output token counts.
type Usage = common.Usage

// Ensemble owns its Agents and logger; it has no invented parent of its own.
// Construction wires children before the single conversation actor starts.
type Ensemble struct {
	mu           sync.Mutex
	closed       bool
	nextHandle   int
	watchers     map[uint64]chan AgentObservation
	nextWatcher  uint64
	shutdownDone chan struct{}
	shutdownErr  error
	jobs         map[int]common.Job
	trace        io.Writer
	logger       *log.Logger
	agents       []common.Agent
}

// New creates an independent application. A nil logger discards diagnostics.
func New(logger *log.Logger) *Ensemble {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &Ensemble{logger: logger, jobs: make(map[int]common.Job), watchers: make(map[uint64]chan AgentObservation), shutdownDone: make(chan struct{})}
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
	if config.StreamCapability == StreamUnknown {
		model := config.ResolvedModel
		if model == "" {
			model = config.Model
		}
		features, _ := LookupModel(model)
		config.StreamCapability = features.Streaming
	}
	if config.StreamCapability != StreamUnknown && config.StreamCapability != StreamPlain && config.StreamCapability != StreamSSE {
		return nil, errors.New("invalid stream capability")
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
	a.Start()
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
	if e.closed {
		e.mu.Unlock()
		<-e.shutdownDone
		return e.shutdownErr
	}
	e.closed = true
	agents := append([]common.Agent(nil), e.agents...)
	e.mu.Unlock()
	// Stop every child before waiting for any potentially uncooperative work.
	errs := make([]error, len(agents))
	var joins sync.WaitGroup
	for i, a := range agents {
		joins.Add(1)
		go func() { defer joins.Done(); errs[i] = a.Shutdown() }()
	}
	joins.Wait()
	e.shutdownErr = errors.Join(errs...)
	close(e.shutdownDone)
	return e.shutdownErr
}

type ToolHandler = common.ToolHandler
type ToolCall = common.ToolCall
type Jobs = common.Jobs
type Job = common.Job
type JobData = common.JobData
type ToolData = common.ToolData
type JobReport = common.JobReport
type JobStatus = common.JobStatus
type ToolLimits = common.ToolLimits
type Ref = common.Ref

const (
	ToolCalled      = common.ToolCalled
	ToolReturned    = common.ToolReturned
	RefHandle       = common.RefHandle
	JobRunning      = common.JobRunning
	JobDone         = common.JobDone
	JobKilledStatus = common.JobKilledStatus
	JobKilled       = common.JobKilled
)

func DefaultToolLimits() ToolLimits { return common.DefaultToolLimits() }

// Mailbox and observation contracts are shared by importers and adapters.
type Inbound = common.Inbound
type UserMessage = common.UserMessage
type Hint = common.Hint
type Ephemeral = common.Ephemeral
type Question = common.Question
type Interrupt = common.Interrupt
type Result = common.Result
type Observation = common.Observation
type Observer = common.Observer
type PartDelta = common.PartDelta
type PartFinal = common.PartFinal
type HintObserved = common.HintObserved
type ObservationQueue = common.ObservationQueue
type ParseCallbacks = common.ParseCallbacks
type DeltaKind = common.DeltaKind
type Frame = common.Frame

const (
	DeltaText     = common.DeltaText
	DeltaThinking = common.DeltaThinking
	DeltaToolCall = common.DeltaToolCall
)

type TraceRecord = common.TraceRecord
type StreamCapability = common.StreamCapability

const (
	StreamUnknown = common.StreamUnknown
	StreamPlain   = common.StreamPlain
	StreamSSE     = common.StreamSSE
)

type StateChanged = common.StateChanged
type ModelFeatures = common.ModelFeatures
type Media = common.Media
type RefKind = common.RefKind

const (
	MediaImage    = common.MediaImage
	MediaAudio    = common.MediaAudio
	MediaVideo    = common.MediaVideo
	MediaDocument = common.MediaDocument
	RefPath       = common.RefPath
	RefURI        = common.RefURI
)
const DefaultMaxToolRounds = common.DefaultMaxToolRounds
const HintReceived = common.HintReceived
const JobReported = common.JobReported
const TurnEnded = common.TurnEnded
const TurnInterrupted = common.TurnInterrupted

var ErrInterrupted = common.ErrInterrupted
var ErrStopped = common.ErrStopped
var ErrToolRoundLimit = common.ErrToolRoundLimit

func LookupModel(model string) (ModelFeatures, bool) { return common.LookupModel(model) }

// Agent identity belongs to the parent envelope, not conversation events.
type AgentObservation struct {
	Agent       Agent
	Observation Observation
}

func (e *Ensemble) Publish(a Agent, o Observation) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, ch := range e.watchers {
		select {
		case ch <- AgentObservation{a, o}:
		default:
		}
	}
}
func (e *Ensemble) WaitAny(ctx context.Context, predicate func(AgentObservation) bool) (AgentObservation, error) {
	e.mu.Lock()
	e.nextWatcher++
	id := e.nextWatcher
	ch := make(chan AgentObservation, 64)
	e.watchers[id] = ch
	e.mu.Unlock()
	defer func() { e.mu.Lock(); delete(e.watchers, id); e.mu.Unlock() }()
	for {
		select {
		case o := <-ch:
			if predicate == nil || predicate(o) {
				return o, nil
			}
		case <-ctx.Done():
			return AgentObservation{}, ctx.Err()
		case <-e.shutdownDone:
			return AgentObservation{}, ErrStopped
		}
	}
}

// SetTrace enables a sanitized process-local diagnostic stream, separate from
// semantic observations and History. The caller owns the writer's lifetime.
func (e *Ensemble) SetTrace(w io.Writer) { e.mu.Lock(); defer e.mu.Unlock(); e.trace = w }
func (e *Ensemble) Trace(r TraceRecord) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.trace != nil {
		b, err := json.Marshal(r)
		if err != nil {
			return
		}
		for _, a := range e.agents {
			key := a.Config().APIKey
			if key == "" {
				continue
			}
			escaped, _ := json.Marshal(key)
			b = []byte(strings.ReplaceAll(string(b), string(escaped[1:len(escaped)-1]), "[redacted]"))
		}
		_, _ = e.trace.Write(append(b, '\n'))
	}
}
