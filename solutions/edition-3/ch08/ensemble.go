// Package ensemble is the public composition root for Chapter 1's conversation.
package ensemble

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	identified   map[uint64]*identifiedSubscriber
	traceQueue   chan traceWrite
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
	e := &Ensemble{traceQueue: make(chan traceWrite, 128), logger: logger, jobs: make(map[int]common.Job), watchers: make(map[uint64]chan AgentObservation), identified: make(map[uint64]*identifiedSubscriber), shutdownDone: make(chan struct{})}
	go e.runTrace()
	return e
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
	a.SetID(fmt.Sprintf("agent-%d", len(e.agents)+1))
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
type Pause = common.Pause
type Unpause = common.Unpause
type SetSpeaking = common.SetSpeaking
type SpeakParts = common.SpeakParts
type Paused = common.Paused
type SpeakingChanged = common.SpeakingChanged
type SpeechRequested = common.SpeechRequested
type TurnEndedObservation = common.TurnEndedObservation
type ToolDispatched = common.ToolDispatched
type ToolFinished = common.ToolFinished
type Seq = common.Seq
type Result = common.Result
type Observation = common.Observation
type Observer = common.Observer
type ViewSnapshot = common.ViewSnapshot
type TurnState = common.TurnState
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

// Agents returns a registry snapshot; all mutation stays with the root.
func (e *Ensemble) Agents() []Agent {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]Agent(nil), e.agents...)
}
func (e *Ensemble) Agent(id string) (Agent, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, a := range e.agents {
		if a.ID() == id {
			return a, true
		}
	}
	return nil, false
}

// Agent identity belongs to the parent envelope, not conversation events.
type AgentObservation struct {
	Agent       Agent
	Observation Observation
}

type identifiedSubscriber struct {
	wake chan struct{}
	done chan struct{}
}

// A bounded wakeup tells the asynchronous subscriber to recover all Agents
// from their latest immutable View. No subscriber code executes in an actor.
func (e *Ensemble) SubscribeIdentified(fn func(AgentObservation)) func() {
	e.mu.Lock()
	e.nextWatcher++
	id := e.nextWatcher
	sub := &identifiedSubscriber{wake: make(chan struct{}, 1), done: make(chan struct{})}
	e.identified[id] = sub
	e.mu.Unlock()
	go func() {
		for {
			select {
			case <-sub.done:
				return
			case <-e.shutdownDone:
				return
			case <-sub.wake:
				for _, a := range e.Agents() {
					select {
					case <-sub.done:
						return
					default:
					}
					fn(AgentObservation{a, common.ViewChanged{}})
				}
			}
		}
	}()
	var once sync.Once
	return func() { once.Do(func() { e.mu.Lock(); delete(e.identified, id); close(sub.done); e.mu.Unlock() }) }
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
	for _, sub := range e.identified {
		select {
		case sub.wake <- struct{}{}:
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

type traceWrite struct {
	writer io.Writer
	data   []byte
}

func (e *Ensemble) runTrace() {
	for {
		select {
		case <-e.shutdownDone:
			return
		case item := <-e.traceQueue:
			_, _ = item.writer.Write(item.data)
		}
	}
}
func (e *Ensemble) Trace(r TraceRecord) {
	e.mu.Lock()
	writer := e.trace
	keys := make([]string, 0, len(e.agents))
	for _, a := range e.agents {
		keys = append(keys, a.Config().APIKey)
	}
	e.mu.Unlock()
	if writer == nil {
		return
	}
	body, err := json.Marshal(r)
	if err != nil {
		return
	}
	for _, key := range keys {
		if key != "" {
			escaped, _ := json.Marshal(key)
			body = []byte(strings.ReplaceAll(string(body), string(escaped[1:len(escaped)-1]), "[redacted]"))
		}
	}
	select {
	case e.traceQueue <- traceWrite{writer, append(body, '\n')}:
	default:
	} // diagnostics cannot block progress
}

const (
	Idle          = common.Idle
	ResponseEnded = common.ResponseEnded
	ToolReleased  = common.ToolReleased
	ErrorOccurred = common.ErrorOccurred
)
