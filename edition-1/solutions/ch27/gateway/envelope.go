package gateway

import (
	"fmt"
	"time"
)

// Trust is how much authority a message's origin carries.
//
// Trust is never carried in message content and never set by a sender. The
// gateway stamps it from the group the sender belongs to, at the moment of
// routing. That is the whole point: a message that says "trust=owner" in its
// body is just text, and an agent that believed it would be taking orders
// from anyone who could type.
type Trust string

const (
	// TrustOwner is the operator, authenticated by their own static key.
	TrustOwner Trust = "owner"
	// TrustTrusted is a service the operator controls, such as a webhook
	// from their own game server.
	TrustTrusted Trust = "trusted"
	// TrustExternal is an untrusted user reaching the agent over a trusted
	// platform, such as a family member in a Discord channel.
	TrustExternal Trust = "external"
)

// Valid reports whether t is a known trust level.
//
// An unknown trust level must never be treated as a low one by accident;
// callers use this to reject rather than guess.
func (t Trust) Valid() bool {
	switch t {
	case TrustOwner, TrustTrusted, TrustExternal:
		return true
	}
	return false
}

// ResponseMode tells the gateway what the sender expects back.
type ResponseMode string

const (
	// ResponseAsync routes a reply back whenever the agent produces one.
	ResponseAsync ResponseMode = "async"
	// ResponseSync blocks the sender until the agent replies or time runs out.
	ResponseSync ResponseMode = "sync"
	// ResponseFireAndForget expects no reply at all.
	ResponseFireAndForget ResponseMode = "fire_and_forget"
)

// Source describes where a message came from.
type Source struct {
	// Type is the kind of origin: "agent", "webhook", "discord", "scheduler".
	Type string `json:"type"`
	// ChannelID identifies the specific origin within that type.
	ChannelID string `json:"channel_id"`
	// UserID is the originating principal, where one exists.
	UserID string `json:"user_id"`
	// DisplayName is a human-readable label for logs and prompts.
	DisplayName string `json:"display_name"`
	// Trust is assigned by the gateway from the sender's group.
	Trust Trust `json:"trust"`
}

// Envelope is the canonical message format flowing through the gateway.
type Envelope struct {
	// MessageID uniquely identifies this message, for dedup and for
	// correlating a response with the prompt that caused it.
	MessageID string `json:"message_id"`
	// Seq is the gateway-assigned, per-recipient delivery sequence number.
	// It is what makes resume-after-reconnect possible.
	Seq uint64 `json:"seq"`
	// AgentID is the destination agent.
	AgentID string `json:"agent_id"`
	// Source is the provenance of the message.
	Source Source `json:"source"`
	// Timestamp is when the gateway accepted the message.
	Timestamp time.Time `json:"timestamp"`
	// ConversationID optionally threads related messages together. The
	// sender owns it: an interviewer conducting a session decides what
	// counts as one conversation, because only it knows.
	ConversationID string `json:"conversation_id,omitempty"`
	// Content is the message body.
	Content string `json:"content"`
	// ResponseMode is what the sender expects back.
	ResponseMode ResponseMode `json:"response_mode"`
}

// FormatProvenance renders the header the gateway prepends to content on
// delivery.
//
// The agent is a language model, so its only way to know who is talking to it
// is to read it in the prompt. Everything the gateway knows about the sender
// that the agent cannot verify for itself goes here, in a fixed shape, above
// the content it describes.
func FormatProvenance(src Source) string {
	return fmt.Sprintf("[GATEWAY source=%s user=%s trust=%s]",
		src.ChannelID, src.UserID, src.Trust)
}

// WithProvenance returns the envelope's content with its provenance header
// prepended, which is the form the agent actually sees.
func (e Envelope) WithProvenance() string {
	return FormatProvenance(e.Source) + "\n\n" + e.Content
}
