package common

import (
	"context"
	"errors"
)

// DefaultToolBatchLimit counts dispatched batches, not calls or HTTP requests.
const DefaultToolBatchLimit = 200

// ErrInterrupted ends one turn without stopping the Agent or retained jobs.
var ErrInterrupted = errors.New("turn interrupted")

// ErrStopped rejects work after the Agent's admission boundary closes.
var ErrStopped = errors.New("agent stopped")

// ErrToolBatchLimit means a provider requested another batch after the allowance.
var ErrToolBatchLimit = errors.New("tool batch limit reached")

// Input is the closed external mailbox vocabulary. Post is admission, not completion.
type Input interface{ agentInput() }

// UserMessage starts an idle Agent; while busy, Post treats it as pending context.
// SubmitMessage instead queues a distinct request with a reliable completion.
type UserMessage struct {
	// Sender and Recipient preserve application-supplied addressing, not routing.
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	// RequestID and TurnID preserve application-supplied correlation verbatim.
	RequestID string `json:"request_id,omitempty"`
	TurnID    string `json:"turn_id,omitempty"`

	// Text is used when Parts is empty.
	Text string
	// Parts carries a full neutral message, copied when admitted.
	Parts []Part
	// Ephemera supplies one-request system context without creating dialogue.
	Ephemera []Part
}

func (UserMessage) agentInput() {}

// Hint is an explicit text form of state-classified live input.
type Hint struct {
	// Text becomes pending context while busy, or an ordinary prompt while idle.
	Text string
}

func (Hint) agentInput() {}

// Interrupt ends the current turn, not the Agent or any independently owned job.
type Interrupt struct{}

func (Interrupt) agentInput() {}

// Result is delivered once on a request's private buffered reply channel.
type Result struct {
	// Text contains visible assistant prose (or the direct tool outcome).
	Text string
	// Err reports request, cancellation or persistence failure; tool errors are feedback.
	Err error
}

// HubEvent is deliberately smaller than History's durable event vocabulary.
type HubEvent interface{ hubEvent() }

// PartDelta is best-effort display progress; whole-part deltas are supported here.
type PartDelta struct {
	// Part is a copied neutral part associated with this update.
	Part Part
	// PartID identifies the part within a finalized response.
	PartID int
	// Chunk is displayable text, including immediate hint acknowledgement.
	Chunk string
}

func (PartDelta) hubEvent() {}

// PartFinal publishes a completed visible or structural part after its fact commits.
type PartFinal struct {
	// Part preserves tool call/result correlation and owns its copied payload.
	Part Part
	// Seq identifies the durable event containing the part.
	Seq uint64
	// PartID identifies the part within that event.
	PartID int
}

func (PartFinal) hubEvent() {}

// StateChanged describes a reduced state transition, never an event-log command.
type StateChanged struct {
	// State is the current state; it equals To.
	State TurnState
	// From is the prior reduced state.
	From TurnState
	// To is the new reduced state.
	To TurnState
}

func (StateChanged) hubEvent() {}

// Observer runs on the actor, after reduction. It must not block or call Ask.
type Observer func(HubEvent)

// AgentObservation identifies the child behind a root-level observation.
type AgentObservation struct {
	// Agent is the emitting child, not a copied name or configuration.
	Agent Agent
	// Event is best-effort progress, not reliable completion.
	Event HubEvent
}

// Actor is embedded in the public Agent contract.
type Actor interface {
	// Start idempotently starts the owner loop; public admission also starts it.
	Start()
	// Post copies and admits input; its return is not a turn-completion receipt.
	Post(Input) error
	// Submit queues a distinct prompt with one reliable buffered reply.
	Submit(context.Context, string) <-chan Result
	// SubmitMessage is Submit's rich-message form.
	SubmitMessage(context.Context, UserMessage) <-chan Result
	// Observe installs a nonblocking callback on the actor's publish path.
	Observe(Observer)
	// Wait consumes bounded, lossy progress; it must not gate a workflow handoff.
	Wait(context.Context, ...func(HubEvent) bool) (HubEvent, error)
	// Drain waits for already-admitted turns, for executable EOF coordination.
	Drain(context.Context) error
	// State returns the current reduced state, not a speculative worker state.
	State() TurnState
	// StopAdmission closes external admission and cancels observation of the active turn.
	StopAdmission()
	// Shutdown joins owned work, drains terminal reports and saves final history.
	Shutdown() error
}

// Media is a bitset of supported input classes, independent of wire formats.
type Media uint8

const (
	// MediaImage permits local raster-image attachments.
	MediaImage Media = 1 << iota
	// MediaAudio permits local audio attachments.
	MediaAudio
	// MediaVideo permits local video attachments.
	MediaVideo
	// MediaDocument permits local PDF document attachments.
	MediaDocument
)

// ModelFeatures is the data-only capability row for an explicitly known model.
type ModelFeatures struct {
	// Media lists supported input classes; zero is conservative text-only behavior.
	Media Media
}

// MediaCapabilities combines a provider/model match with its input classes.
type MediaCapabilities struct {
	// Images is the image capability, retained as a convenient readable check.
	Images bool
	// Media is the complete input-class bitset for this exact provider/model pair.
	Media Media
}
