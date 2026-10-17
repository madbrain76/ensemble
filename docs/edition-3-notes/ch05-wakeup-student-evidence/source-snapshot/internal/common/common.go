// Package common holds shared vocabulary, not implementation behavior.
package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

// ToolDefinition is the model-facing name, explanation and JSON object schema.
type ToolDefinition struct {
	// Name identifies a registered operation.
	Name string `json:"name"`
	// Description explains when the model should use this operation.
	Description string `json:"description"`
	// Parameters describes the strict JSON arguments accepted by Execute.
	Parameters json.RawMessage `json:"parameters"`
}

// Tool is an Agent-owned synchronous operation. Errors become paired results.
type Tool interface {
	// Parent returns the owner responsible for dispatch and policy.
	Parent() Agent
	// Definition returns independent model-facing metadata.
	Definition() ToolDefinition
	// Execute performs one synchronous operation; errors become result feedback.
	Execute(context.Context, json.RawMessage) (string, error)
}

// Config belongs to an Agent and stays fixed for its conversation.
type Config struct {
	// ToolBatchLimit: zero selects DefaultToolBatchLimit; positive values bound dispatched batches.
	ToolBatchLimit int
	// HistoryPath optionally persists History before each reliable request reply.
	HistoryPath string
	// Workspace is the launch-time command directory, separate from persisted data.
	Workspace string
	// DataDir contains job output under cr/io.
	DataDir string
	// EnableTools installs the six Chapter 3 built-ins; false keeps a text-only Agent.
	EnableTools bool
	// RequireToolExplanation selects this configured model's visible-prose policy.
	// It is opt-in, not inferred from model names, and applies to every call batch.
	RequireToolExplanation bool
	// BaseURL is the provider origin, without the API path.
	BaseURL string
	// APIKey is required only when sending a live request.
	APIKey string
	// Model is the provider's model identifier, required for live requests.
	Model string
	// ResolvedModel is an optional exact replay identity, not the routing name.
	ResolvedModel string
	// Vendor selects anthropic (the default), openai or gemini.
	Vendor Vendor
	// System is stable instruction text, outside dialogue.
	System string
	// MaxTokens bounds output quantity, never model-work time.
	MaxTokens int
}

// Message is the Chapter 1 text-only compatibility view of committed dialogue.
type Message struct {
	// Role is user or assistant.
	Role string `json:"role"`
	// Content concatenates text without inserting separators.
	Content string `json:"content"`
}

// Usage is disjoint provider-measured accounting, not an estimate.
type Usage struct {
	// Input excludes tokens counted in either cache category.
	Input int `json:"input"`
	// CacheWrite counts newly cached input tokens.
	CacheWrite int `json:"cache_write"`
	// CacheRead counts cached input tokens reused.
	CacheRead int `json:"cache_read"`
	// Output includes generated thinking tokens when the provider separates them.
	Output int `json:"output"`
}

// Provenance binds replay-only material to its exact issuing model and surface.
type Provenance struct {
	// Vendor is anthropic, openai or gemini.
	Vendor Vendor `json:"vendor"`
	// Model is the issuing model identifier.
	Model string `json:"model"`
	// Surface is messages, chat_completions or generate_content.
	Surface Surface `json:"surface"`
}

// Ref is data only: this chapter never opens or fetches its locator.
type Ref struct {
	// Kind is 1 for path, 2 for URI or 3 for an opaque handle; zero is invalid.
	Kind RefKind `json:"kind"`
	// Locator preserves the original reference bytes.
	Locator string `json:"locator"`
}

// Part is a closed tagged union. Type selects the fields that carry meaning.
// Raw JSON fields carry JSON values, never double-encoded JSON strings.
type Part struct {
	// Type is text, blob, opaque, tool_call, tool_result or redacted.
	Type string `json:"type"`
	// Text is ordinary prose, including a deliberately empty string.
	Text string `json:"text,omitempty"`
	// MIME describes a blob or a removed blob.
	MIME string `json:"mime,omitempty"`
	// Ref locates a blob without granting this layer permission to load it.
	Ref *Ref `json:"ref,omitempty"`
	// From binds opaque, signed-text, and tool-call replay material to its issuer.
	From Provenance `json:"from"`
	// Data preserves a standalone opaque vendor block.
	Data json.RawMessage `json:"data,omitempty"`
	// CallID pairs a tool call and result; a renderer synthesizes one from sequence if absent.
	CallID string `json:"call_id,omitempty"`
	// Name identifies a requested tool, not an executable registered here.
	Name string `json:"name,omitempty"`
	// Args contains tool arguments as a JSON value.
	Args json.RawMessage `json:"args,omitempty"`
	// Opaque preserves model-bound replay material attached to this part.
	// For signed Gemini text it is the original signature JSON value, not prose.
	Opaque json.RawMessage `json:"opaque,omitempty"`
	// Parts contains a tool result's parts, in their original order.
	Parts []Part `json:"parts,omitempty"`
	// IsError records whether a tool result failed.
	IsError bool `json:"is_error,omitempty"`
	// Stub describes content removed from the projection.
	Stub string `json:"stub,omitempty"`
}

// UnmarshalJSON rejects obsolete blob fields rather than silently dropping data.
func (p *Part) UnmarshalJSON(data []byte) error {
	type plain Part
	var decoded plain
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, ok := fields["path"]; ok {
		return fmt.Errorf("obsolete path; use ref")
	}
	*p = Part(decoded)
	return nil
}

// MessageEvent records who delivered parts to the Agent.
type MessageEvent struct {
	// Sender and Recipient preserve application-supplied addressing, not routing.
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	// RequestID and TurnID preserve application-supplied correlation verbatim.
	RequestID string `json:"request_id,omitempty"`
	TurnID    string `json:"turn_id,omitempty"`

	// Actor is human, agent, tool or system.
	Actor string `json:"actor"`
	// Parts is the delivered content.
	Parts []Part `json:"parts"`
}

// RequestEvent records the actual target of an outbound request.
type RequestEvent struct {
	// To identifies the target vendor, model and surface.
	To Provenance `json:"to"`
}

// ResponseEvent records parsed model content and measured usage.
type ResponseEvent struct {
	// Parts contains everything the model said, including opaque blocks.
	Parts []Part `json:"parts"`
	// Usage is normalized disjoint accounting.
	Usage Usage `json:"usage"`
	// From identifies the response's exact source.
	From Provenance `json:"from"`
}

// ToolEvent records dispatch or return, without implementing execution.
type ToolEvent struct {
	Job *JobRecord `json:"job,omitempty"`
	// CallID pairs dispatch and return with the model's call.
	CallID string `json:"call_id"`
	// Name and Args describe a dispatch.
	Name string `json:"name,omitempty"`
	// Args is the dispatched JSON argument value.
	Args json.RawMessage `json:"args,omitempty"`
	// Parts describes a return.
	Parts []Part `json:"parts,omitempty"`
	// IsError distinguishes failed tool returns.
	IsError bool `json:"is_error,omitempty"`
}

// ErrorEvent is a failed operation's sanitized diagnostic.
type ErrorEvent struct {
	// Message explains the failure without credentials or raw provider bodies.
	Message string `json:"message"`
}

// Redaction removes projected content, never recorded facts.
type Redaction struct {
	// From and To are an inclusive sequence span.
	From uint64 `json:"from"`
	// To is the inclusive last affected sequence.
	To uint64 `json:"to"`
	// Level is redact_result, redact_tool, redact_dialogue or redact_summary.
	Level string `json:"level"`
	// Reason explains why content was removed.
	Reason string `json:"reason"`
	// Replacement supplies a summary for redact_summary.
	Replacement []Part `json:"replacement,omitempty"`
}

// Event is one append-only fact. Exactly one payload matches Type.
type Event struct {
	Job *JobRecord `json:"job,omitempty"`
	// Seq orders events independently of wall-clock time.
	Seq uint64 `json:"seq"`
	// Type selects the payload.
	Type string `json:"type"`
	// Time is metadata and does not determine replay order.
	Time time.Time `json:"time"`
	// Message carries message_received.
	Message *MessageEvent `json:"message,omitempty"`
	// Request carries request_sent.
	Request *RequestEvent `json:"request,omitempty"`
	// Response carries response_started or response_ended.
	Response *ResponseEvent `json:"response,omitempty"`
	// Tool carries tool_called or tool_returned.
	Tool *ToolEvent `json:"tool,omitempty"`
	// Error carries error_occurred.
	Error *ErrorEvent `json:"error,omitempty"`
	// Redact carries redacted.
	Redact *Redaction `json:"redact,omitempty"`
}

// Turn is a projected contribution with its originating sequence.
type Turn struct {
	// Seq connects this projection to the original fact.
	Seq uint64 `json:"seq"`
	// Actor is human, agent or tool; redaction summaries may use system.
	Actor string `json:"actor"`
	// Parts is the projected (possibly redacted) content.
	Parts []Part `json:"parts"`
}

// Context is disposable replay output, independent of any vendor wire shape.
type Context struct {
	// Turn is the synchronous lifecycle state derived from known events.
	Turn TurnState `json:"turn"`
	// Dialogue contains retained conversational turns and replacement summaries.
	Dialogue []Turn `json:"dialogue"`
	// Ephemera is pending one-request material.
	Ephemera []Part `json:"ephemera"`
	// Usage accumulates measured response-ended facts.
	Usage Usage `json:"usage"`
}

// Ensemble is the parent service visible to Agents.
type Ensemble interface {
	Jobs() JobTable
	// Logger supplies diagnostics, never protocol stdout.
	Logger() *log.Logger
}

// Agent owns configuration, History, Engine and conversation serialization.
type Agent interface {
	Actor
	RegisterTool(ToolDefinition, ToolHandler) error
	CallTool(context.Context, string, json.RawMessage) (string, error)
	// Parent returns the owning root.
	Parent() Ensemble
	// Config returns immutable configuration by value.
	Config() Config
	// Engine returns the Agent-owned provider seam.
	Engine() Engine
	// Tools returns a snapshot of the registered operations.
	Tools() []Tool
	// History returns the Chapter 1 text-only view of completed turns.
	History() []Message
	// Log returns the History owner, the sole route into stored facts.
	Log() History
	// Ask captures a user turn, sends it and records the resulting facts.
	Ask(context.Context, string) (string, error)
}

// History owns append, load, dump and reduction, not vendor translation.
type History interface {
	// Parent returns the owning Agent.
	Parent() Agent
	// Append assigns sequence and time, validates and copies a fact.
	Append(Event) error
	// Events returns an independent snapshot of recorded facts.
	Events() []Event
	// Context rebuilds the projection from recorded facts.
	Context() Context
	// Dump emits a version header and one event per line.
	Dump(io.Writer) error
	// Load atomically replaces the log with validated JSONL facts.
	Load(io.Reader) error
}

// Engine owns the vendor seam, HTTP connection and live usage counter.
type Engine interface {
	// SendPrepared sends immutable bytes rendered before request_sent.
	SendPrepared(context.Context, json.RawMessage) (ResponseEvent, error)
	// Parent returns the owning Agent.
	Parent() Agent
	// Target returns this Engine's configured provider provenance.
	Target() Provenance
	// Send renders one context and returns parsed facts, without editing History.
	Send(context.Context, Context) (ResponseEvent, error)
	// Render is pure vendor projection; it never makes an HTTP request.
	Render(Context, Provenance) (json.RawMessage, error)
	// Parse converts one full vendor response to neutral facts.
	Parse([]byte, Provenance) (ResponseEvent, error)
	// Usage returns cumulative measured live-response accounting.
	Usage() Usage
}

// MarshalJSON omits irrelevant provenance rather than serializing an invalid
// zero vendor on text. Tagged values remain ordinary inspectable JSON.
func (p Part) MarshalJSON() ([]byte, error) {
	type plain Part
	var from *Provenance
	if p.From != (Provenance{}) {
		value := p.From
		from = &value
	}
	return json.Marshal(struct {
		plain
		// From binds optional replay-only material to its issuer.
		From *Provenance `json:"from,omitempty"`
	}{plain(p), from})
}

// Vendor is a closed set with an invalid zero value.
type Vendor uint8

const (
	// Anthropic uses the Messages API.
	Anthropic Vendor = iota + 1
	// OpenAI uses the ChatCompletions API.
	OpenAI
	// Gemini uses generateContent.
	Gemini
)

// String returns the canonical wire name, or invalid for the zero value.
func (v Vendor) String() string {
	switch v {
	case Anthropic:
		return "anthropic"
	case OpenAI:
		return "openai"
	case Gemini:
		return "gemini"
	}
	return "invalid"
}

// MarshalJSON keeps logs human-readable and refuses an unset vendor.
func (v Vendor) MarshalJSON() ([]byte, error) {
	if v < Anthropic || v > Gemini {
		return nil, fmt.Errorf("invalid vendor")
	}
	return json.Marshal(v.String())
}

// UnmarshalJSON rejects unknown vendors instead of guessing.
func (v *Vendor) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	switch name {
	case "anthropic":
		*v = Anthropic
	case "openai":
		*v = OpenAI
	case "gemini":
		*v = Gemini
	default:
		return fmt.Errorf("unknown vendor %q", name)
	}
	return nil
}

// Surface identifies a closed API shape independently of its vendor.
type Surface uint8

const (
	// Messages is Anthropic's message-block surface.
	Messages Surface = iota + 1
	// ChatCompletions is the OpenAI chat-message surface.
	ChatCompletions
	// GenerateContent is Gemini's contents/parts surface.
	GenerateContent
)

// String returns the canonical surface name, or invalid for the zero value.
func (s Surface) String() string {
	switch s {
	case Messages:
		return "messages"
	case ChatCompletions:
		return "chat_completions"
	case GenerateContent:
		return "generate_content"
	}
	return "invalid"
}

// MarshalJSON always emits the canonical spelling.
func (s Surface) MarshalJSON() ([]byte, error) {
	if s < Messages || s > GenerateContent {
		return nil, fmt.Errorf("invalid surface")
	}
	return json.Marshal(s.String())
}

// UnmarshalJSON also accepts the public historical Gemini alias.
func (s *Surface) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	switch name {
	case "messages":
		*s = Messages
	case "chat_completions":
		*s = ChatCompletions
	case "generate_content", "generatecontent":
		*s = GenerateContent
	default:
		return fmt.Errorf("unknown surface %q", name)
	}
	return nil
}

// TurnState is the reduced conversation lifecycle.
type TurnState string

const (
	// Idle has no pending user input or tool return.
	Idle TurnState = "idle"
	// InputPending awaits the next synchronous model request.
	InputPending TurnState = "input_pending"
	// InFlight has a sent request without a completed response.
	InFlight TurnState = "in_flight"
	// ToolsPending awaits one or more recorded tool results.
	ToolsPending TurnState = "tools_pending"
	Interrupted  TurnState = "interrupted"
)

// RefKind is a closed locator vocabulary; zero is never a reference kind.
type RefKind uint8

const (
	// RefPath preserves a path without opening it.
	RefPath RefKind = iota + 1
	// RefURI preserves a URI without fetching it.
	RefURI
	// RefHandle preserves an opaque loader-owned handle.
	RefHandle
)
