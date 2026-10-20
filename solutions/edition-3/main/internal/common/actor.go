package common

import (
	"context"
	"errors"
)

const DefaultMaxToolRounds = 200

var ErrInterrupted = errors.New("turn interrupted")
var ErrStopped = errors.New("Agent stopped")
var ErrToolRoundLimit = errors.New("tool round limit reached")

type Inbound interface{ inbound() }
type UserMessage struct{ Text string }
type Hint struct{ Text string }
type Ephemeral struct{ Text string }

// Question preserves the application's addressing and correlation metadata.
type Question struct{ Text, Sender, Recipient, RequestID, TurnID string }
type Interrupt struct{}

func (UserMessage) inbound() {}
func (Hint) inbound()        {}
func (Ephemeral) inbound()   {}
func (Question) inbound()    {}
func (Interrupt) inbound()   {}

type Result struct {
	TurnID uint64
	Text   string
	Err    error
}
type Observation interface{ observation() }
type DeltaKind string

const (
	DeltaText     DeltaKind = "text"
	DeltaThinking DeltaKind = "thinking"
	DeltaToolCall DeltaKind = "tool_call"
)

type PartDelta struct {
	TurnID     uint64    `json:"turn_id"`
	ResponseID uint64    `json:"response_id"`
	Kind       DeltaKind `json:"kind"`
	// Name is supplemental metadata; tool-call Chunk content carries the
	// arriving name fragments followed by argument fragments for appenders.
	Name   string `json:"name,omitempty"`
	PartID uint64 `json:"part_id"`
	Chunk  string `json:"chunk"`
}
type PartFinal struct {
	TurnID     uint64 `json:"turn_id"`
	ResponseID uint64 `json:"response_id"`
	Seq        uint64 `json:"seq"`
	PartID     uint64 `json:"part_id"`
	Part       Part   `json:"part"`
}
type StateChanged struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (PartDelta) observation()    {}
func (PartFinal) observation()    {}
func (StateChanged) observation() {}

type Observer interface{ Observe(Observation) }

type Media uint8

const (
	MediaImage Media = 1 << iota
	MediaAudio
	MediaVideo
	MediaDocument
)

type ModelFeatures struct {
	Media     Media
	Streaming StreamCapability
}

// Exact rows only. An absent row permits text, never guessed attachments.
func LookupModel(model string) (ModelFeatures, bool) {
	switch model {
	case "claude-3-haiku-20240307", "claude-haiku-5-5", "claude-sonnet-4-6", "claude-opus-4-6", "gpt-5.4-mini", "gpt-5.4":
		return ModelFeatures{Media: MediaImage, Streaming: StreamSSE}, true
	case "gemini-3.8-flash":
		return ModelFeatures{Media: MediaImage | MediaAudio | MediaVideo, Streaming: StreamSSE}, true
	}
	return ModelFeatures{}, false
}

// HintObserved is emitted only after the actor classifies live input as a hint.
// Sender distinguishes local human acknowledgements from child/future clients.
type HintObserved struct {
	TurnID uint64
	Sender string
	Parts  PartList
}

func (HintObserved) observation() {}

// ObservationQueue is an opt-in lossless pull subscription. Actor publishing
// only appends to memory; slow terminal I/O cannot block the actor or HTTP.
// Close detaches it and releases retained observations.
type ObservationQueue interface {
	Wake() <-chan struct{}
	Drain() []Observation
	Close()
}

// ToolAdmission lets a low-level Prepare/Await caller use the same Agent-owned
// call-record barrier as Invoke and model dispatch. It is internal graph traffic.
type ToolAdmission struct {
	Context   context.Context
	Operation ToolOperation
	Reply     chan ToolData
}

func (ToolAdmission) inbound() {}
