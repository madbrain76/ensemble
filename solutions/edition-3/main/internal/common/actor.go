package common

import "errors"

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
type PartDelta struct {
	PartID uint64 `json:"part_id"`
	Chunk  string `json:"chunk"`
}
type PartFinal struct {
	Seq    uint64 `json:"seq"`
	PartID uint64 `json:"part_id"`
	Part   Part   `json:"part"`
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

type ModelFeatures struct{ Media Media }

// Exact rows only. An absent row permits text, never guessed attachments.
func LookupModel(model string) (ModelFeatures, bool) {
	switch model {
	case "claude-3-haiku-20240307", "claude-haiku-5-5", "claude-sonnet-4-6", "claude-opus-4-6", "gpt-5.4-mini", "gpt-5.4":
		return ModelFeatures{Media: MediaImage}, true
	case "gemini-3.8-flash":
		return ModelFeatures{Media: MediaImage | MediaAudio | MediaVideo}, true
	}
	return ModelFeatures{}, false
}
