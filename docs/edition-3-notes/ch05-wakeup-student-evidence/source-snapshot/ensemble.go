// Package ensemble is the public entry point for the conversation core.
// Construct a root, add Agents, then call Ask; clients need no internal imports.
package ensemble

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"ensemble/internal/agent"
	"ensemble/internal/common"
	"ensemble/internal/engine"
	"ensemble/internal/history"
	"ensemble/internal/jobs"
	"ensemble/internal/tools"
)

// Config is an Agent's immutable provider and system configuration.
type Config = common.Config

// Tool is an Agent-owned synchronous operation.
type Tool = common.Tool

// ToolDefinition is the public model-facing operation schema.
type ToolDefinition = common.ToolDefinition

// Message is a committed user or assistant entry.
type Message = common.Message

// Usage reports disjoint measured input, cache-write, cache-read and output tokens.
type Usage = common.Usage

// Agent is a sequential conversation owned by an Ensemble.
type Agent = common.Agent

// Engine is the Agent-owned provider connection and usage counter.
type Engine = common.Engine

// Ensemble assembles Agents and supplies their shared diagnostic logger.
// Root child admission and shutdown are synchronized; each Agent owns its loop.
type Ensemble struct {
	jobs         common.JobTable
	logger       *log.Logger
	agents       []Agent
	mu           sync.Mutex
	stopped      bool
	stopOnce     sync.Once
	stopErr      error
	observations chan common.AgentObservation
}

// New constructs a root; a nil logger selects stderr, never protocol stdout.
func New(logger *log.Logger) *Ensemble {
	if logger == nil {
		logger = log.New(os.Stderr, "", 0)
	}
	e := &Ensemble{logger: logger, observations: make(chan common.AgentObservation, 128)}
	e.jobs = jobs.New(e)
	return e
}

// Logger exposes the root-owned diagnostic service to descendants.
func (e *Ensemble) Logger() *log.Logger { return e.logger }

// Agents returns a snapshot of the children registered with this root.
func (e *Ensemble) Agents() []Agent {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]Agent(nil), e.agents...)
}

// NewAgent validates configuration and constructs both sides of the owner chain.
func (e *Ensemble) NewAgent(config Config) (Agent, error) {
	// Reject nil/zero roots before an interface could hide a typed-nil parent.
	if e == nil || e.logger == nil {
		return nil, fmt.Errorf("construct Ensemble with New before adding Agents")
	}
	// Rendering and loading need no credentials. Live-send validation belongs
	// to Engine.Send; malformed explicit endpoints still fail at construction.
	if config.BaseURL != "" {
		endpoint, err := url.Parse(config.BaseURL)
		if err != nil || endpoint.Host == "" || (endpoint.Scheme != "http" && endpoint.Scheme != "https") {
			return nil, fmt.Errorf("base URL must be an HTTP(S) origin")
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if e.stopped {
		return nil, common.ErrStopped
	}
	a, err := agent.New(e, config)
	if err != nil {
		return nil, err
	}
	eng, err := engine.New(a)
	if err != nil {
		return nil, err
	}
	hist, err := history.New(a)
	if err != nil {
		return nil, err
	}
	agent.AttachChildren(a, eng, hist)
	if config.EnableTools {
		agent.AttachTools(a, tools.NewBuiltins(a))
	}
	a.Observe(func(event common.HubEvent) {
		select {
		case e.observations <- common.AgentObservation{Agent: a, Event: event}:
		default:
		}
	})
	e.agents = append(e.agents, a)
	return a, nil
}

// Part is one neutral content block; blob references are never loaded here.
type Part = common.Part

// Ref is a path, URI or opaque-handle blob reference.
type Ref = common.Ref

// Provenance identifies an exact vendor, model and API surface.
type Provenance = common.Provenance

// Event is a versioned log fact with one type-specific payload.
type Event = common.Event

// MessageEvent records an incoming actor and its parts.
type MessageEvent = common.MessageEvent

// RequestEvent records the actual target of a request.
type RequestEvent = common.RequestEvent

// ResponseEvent preserves parsed content and measured usage.
type ResponseEvent = common.ResponseEvent

// ToolEvent records a dispatch or return, without executing a tool.
type ToolEvent = common.ToolEvent

// ErrorEvent records a sanitized operation failure.
type ErrorEvent = common.ErrorEvent

// Redaction describes a projection-only sequence-span transform.
type Redaction = common.Redaction

// Turn is one projected contribution tied to its originating sequence.
type Turn = common.Turn

// Context is the rebuildable turn, dialogue, ephemera and usage projection.
type Context = common.Context

// History is the sole owner API for append, dump, load and replay.
type History = common.History

// Vendor and Surface are closed provider vocabulary, with invalid zero values.
type Vendor = common.Vendor
type Surface = common.Surface

// Supported vendors and their chapter-specific API surfaces.
const (
	Anthropic       = common.Anthropic
	OpenAI          = common.OpenAI
	Gemini          = common.Gemini
	Messages        = common.Messages
	ChatCompletions = common.ChatCompletions
	GenerateContent = common.GenerateContent
)

// TurnState records the synchronous lifecycle derived by History.
type TurnState = common.TurnState

const (
	Idle         = common.Idle
	InputPending = common.InputPending
	InFlight     = common.InFlight
	ToolsPending = common.ToolsPending
	Interrupted  = common.Interrupted
)

// RefKind distinguishes paths, URIs and opaque handles without loading them.
type RefKind = common.RefKind

const (
	RefPath   = common.RefPath
	RefURI    = common.RefURI
	RefHandle = common.RefHandle
)

// Jobs returns this root's shared invocation registry.
func (e *Ensemble) Jobs() common.JobTable { return e.jobs }

// Shutdown stops live work before the application persists its final log.
func (e *Ensemble) Shutdown() error {
	e.stopOnce.Do(func() {
		e.mu.Lock()
		e.stopped = true
		agents := append([]Agent(nil), e.agents...)
		e.mu.Unlock()
		for _, a := range agents {
			a.StopAdmission()
		}
		e.stopErr = e.jobs.Shutdown()
		for _, a := range agents {
			if err := a.Shutdown(); err != nil && e.stopErr == nil {
				e.stopErr = err
			}
		}
		close(e.observations)
	})
	return e.stopErr
}

// WaitAny is bounded lossy observation; request replies are the completion authority.
func (e *Ensemble) WaitAny(ctx context.Context, pred ...func(AgentObservation) bool) (AgentObservation, error) {
	for {
		select {
		case event, ok := <-e.observations:
			if !ok {
				return AgentObservation{}, ErrStopped
			}
			if len(pred) > 0 && pred[0] != nil && !pred[0](event) {
				continue
			}
			return event, nil
		case <-ctx.Done():
			return AgentObservation{}, ctx.Err()
		}
	}
}

// ToolContext carries handler cancellation and owner access.
type ToolContext = common.ToolContext

// ToolHandler implements a registered external operation.
type ToolHandler = common.ToolHandler

// JobRecord is the persistent invocation receipt.
type JobRecord = common.JobRecord

// Job is a retained invocation and output cursor.
type Job = common.Job

// WaitOptions bounds observation rather than execution.
type WaitOptions = common.WaitOptions

// Input is the closed mailbox vocabulary; admission is not completion.
type Input = common.Input

// UserMessage carries copied rich input, state-classified by Post.
type UserMessage = common.UserMessage

// Hint adds pending context while busy and starts a turn while idle.
type Hint = common.Hint

// Interrupt ends one turn without killing retained jobs.
type Interrupt = common.Interrupt

// Result is the single reliable reply for an admitted request.
type Result = common.Result

// HubEvent is best-effort display progress, separate from durable history.
type HubEvent = common.HubEvent

// PartDelta carries whole-part or incremental visible progress.
type PartDelta = common.PartDelta

// PartFinal correlates a completed part with its durable event.
type PartFinal = common.PartFinal

// StateChanged describes a reduced state transition.
type StateChanged = common.StateChanged

// Observer must return promptly; it runs on the Agent owner loop.
type Observer = common.Observer

// AgentObservation tags root-level progress with its emitting Agent.
type AgentObservation = common.AgentObservation

// MediaCapabilities describes an explicit provider/model input-class match.
type MediaCapabilities = common.MediaCapabilities

// DefaultToolBatchLimit allows 200 dispatched batches before refusing another.
const DefaultToolBatchLimit = common.DefaultToolBatchLimit

// ErrInterrupted ends one request while the Agent remains usable.
var ErrInterrupted = common.ErrInterrupted

// ErrStopped marks the closed admission boundary.
var ErrStopped = common.ErrStopped

// ErrToolBatchLimit reports refusal to dispatch beyond the configured allowance.
var ErrToolBatchLimit = common.ErrToolBatchLimit

// ModelCapabilities requires an exact known provider/model match.
func ModelCapabilities(v Vendor, model string) MediaCapabilities {
	return engine.ModelCapabilities(v, model)
}

// Observation is the Hub event vocabulary used by Wait.
type Observation = common.HubEvent

// Media is a bitset of supported input classes.
type Media = common.Media

// ModelFeatures is an explicitly known model capability row.
type ModelFeatures = common.ModelFeatures

const (
	// MediaImage permits that input class when present in the model row.
	MediaImage = common.MediaImage
	// MediaAudio permits that input class when present in the model row.
	MediaAudio = common.MediaAudio
	// MediaVideo permits that input class when present in the model row.
	MediaVideo = common.MediaVideo
	// MediaDocument permits that input class when present in the model row.
	MediaDocument = common.MediaDocument
)

// LookupModel returns conservative metadata without inferring capabilities from names.
func LookupModel(model string) (ModelFeatures, bool) { return engine.LookupModel(model) }
