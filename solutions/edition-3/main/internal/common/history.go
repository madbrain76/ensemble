package common

import (
	"encoding/json"
	"fmt"
	"time"
)

// Closed vendor and surface sets have invalid zero values; model names stay open.
type Vendor uint8

const (
	Anthropic Vendor = iota + 1
	OpenAI
	Gemini
)

type Surface uint8

const (
	Messages Surface = iota + 1
	ChatCompletions
	GenerateContent
)

func (v Vendor) String() string { return enumName(int(v), []string{"anthropic", "openai", "gemini"}) }
func (v Surface) String() string {
	return enumName(int(v), []string{"messages", "chat_completions", "generate_content"})
}
func enumName(n int, names []string) string {
	if n < 1 || n > len(names) {
		return ""
	}
	return names[n-1]
}
func (v Vendor) MarshalJSON() ([]byte, error) {
	if v.String() == "" {
		return nil, fmt.Errorf("invalid vendor")
	}
	return json.Marshal(v.String())
}
func (v Surface) MarshalJSON() ([]byte, error) {
	if v.String() == "" {
		return nil, fmt.Errorf("invalid surface")
	}
	return json.Marshal(v.String())
}
func (v *Vendor) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	for _, x := range []Vendor{Anthropic, OpenAI, Gemini} {
		if s == x.String() {
			*v = x
			return nil
		}
	}
	return fmt.Errorf("unknown vendor %q", s)
}
func (v *Surface) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "generatecontent" {
		s = "generate_content"
	}
	for _, x := range []Surface{Messages, ChatCompletions, GenerateContent} {
		if s == x.String() {
			*v = x
			return nil
		}
	}
	return fmt.Errorf("unknown surface %q", s)
}

// Target records the request destination, including an unresolved routing alias.
// A recorded destination is not proof of an opaque part's producer identity.
func (c Config) Target() Provenance {
	model := c.ResolvedModel
	if model == "" {
		model = c.Model
	}
	return Provenance{c.Vendor, model, c.Surface}
}

type Provenance struct {
	Vendor  Vendor  `json:"vendor"`
	Model   string  `json:"model"`
	Surface Surface `json:"surface"`
}
type Seq uint64
type Actor string

const (
	Human     Actor = "human"
	Assistant Actor = "agent"
	System    Actor = "system"
	Tool      Actor = "tool"
)

type TurnState string

const (
	Idle         TurnState = "idle"
	Interrupted  TurnState = "interrupted"
	InputPending TurnState = "input_pending"
	InFlight     TurnState = "in_flight"
	ToolsPending TurnState = "tools_pending"
)

type RefKind = int

const (
	RefPath RefKind = 1
	RefURI  RefKind = 2
)

type Ref struct {
	Kind    RefKind `json:"kind"`
	Locator string  `json:"locator"`
}

// Part is a tagged union. Each variant keeps its own wire-neutral distinctions.
// Raw JSON preserves arbitrary argument numbers and uninterpreted replay bytes.
type Part struct {
	Type    string          `json:"type"`
	Text    string          `json:"text,omitempty"`
	MIME    string          `json:"mime,omitempty"`
	Ref     *Ref            `json:"ref,omitempty"`
	From    *Provenance     `json:"from,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	CallID  string          `json:"call_id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Args    json.RawMessage `json:"args,omitempty"`
	Opaque  json.RawMessage `json:"opaque,omitempty"`
	Parts   PartList        `json:"parts,omitempty"`
	IsError bool            `json:"is_error,omitempty"`
	Stub    string          `json:"stub,omitempty"`
}

// Text must be present even when empty; omission is not an empty reply.
func (p Part) MarshalJSON() ([]byte, error) {
	type plain Part
	if p.Type == "text" {
		return json.Marshal(struct {
			Type   string          `json:"type"`
			Text   string          `json:"text"`
			From   *Provenance     `json:"from,omitempty"`
			Opaque json.RawMessage `json:"opaque,omitempty"`
		}{p.Type, p.Text, p.From, p.Opaque})
	}
	return json.Marshal(plain(p))
}

type PartList []Part

func Text(s string) Part { return Part{Type: "text", Text: s} }

type Entry struct {
	Seq   Seq      `json:"seq"`
	Actor Actor    `json:"actor"`
	Parts PartList `json:"parts"`
}
type Context struct {
	Turn     TurnState `json:"turn"`
	Dialogue []Entry   `json:"dialogue"`
	Ephemera PartList  `json:"ephemera"`
	Usage    Usage     `json:"usage"`
	// Pending is bounded by the current response, not lifetime dispatch history.
	Pending []string `json:"pending,omitempty"`
}
type EventType string

const (
	MessageReceived EventType = "message_received"
	HintReceived    EventType = "hint_received"
	TurnEnded       EventType = "turn_ended"
	TurnInterrupted EventType = "turn_interrupted"
	JobReported     EventType = "job_report"
	RequestSent     EventType = "request_sent"
	ResponseStarted EventType = "response_started"
	ResponseEnded   EventType = "response_ended"
	ToolCalled      EventType = "tool_called"
	ToolReturned    EventType = "tool_returned"
	Redacted        EventType = "redacted"
	ErrorOccurred   EventType = "error_occurred"
	JobKilled       EventType = "job_killed"
)

type MessageData struct {
	Sender    string   `json:"sender,omitempty"`
	Recipient string   `json:"recipient,omitempty"`
	RequestID string   `json:"request_id,omitempty"`
	TurnID    string   `json:"turn_id,omitempty"`
	Actor     Actor    `json:"actor"`
	Parts     PartList `json:"parts"`
}
type RequestData struct {
	To Provenance `json:"to"`
}
type ResponseData struct {
	Parts    PartList   `json:"parts"`
	Usage    Usage      `json:"usage"`
	From     Provenance `json:"from"`
	Thinking string     `json:"thinking,omitempty"`
	// Missing identifies unmeasured or partially measured buckets. Usage adds
	// only known charges; an absent count must not masquerade as measured zero.
	UsageMissing []string `json:"usage_missing,omitempty"`
}
type ToolData struct {
	Job     *JobData        `json:"job,omitempty"`
	CallID  string          `json:"call_id"`
	Name    string          `json:"name,omitempty"`
	Args    json.RawMessage `json:"args,omitempty"`
	Parts   PartList        `json:"parts,omitempty"`
	IsError bool            `json:"is_error,omitempty"`
}
type Redaction string

const (
	RedactResult   Redaction = "redact_result"
	RedactTool     Redaction = "redact_tool"
	RedactDialogue Redaction = "redact_dialogue"
	RedactSummary  Redaction = "redact_summary"
)

type RedactData struct {
	From        Seq       `json:"from"`
	To          Seq       `json:"to"`
	Level       Redaction `json:"level"`
	Replacement PartList  `json:"replacement,omitempty"`
	Reason      string    `json:"reason,omitempty"`
}
type ErrorData struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Event is the interoperable eight-case envelope, never a vendor request dump.
type Event struct {
	Stale    bool          `json:"stale,omitempty"`
	TurnID   uint64        `json:"turn_id,omitempty"`
	Job      *JobData      `json:"job,omitempty"`
	Seq      Seq           `json:"seq"`
	Type     EventType     `json:"type"`
	Time     time.Time     `json:"time"`
	Message  *MessageData  `json:"message,omitempty"`
	Request  *RequestData  `json:"request,omitempty"`
	Response *ResponseData `json:"response,omitempty"`
	Tool     *ToolData     `json:"tool,omitempty"`
	Redact   *RedactData   `json:"redact,omitempty"`
	Error    *ErrorData    `json:"error,omitempty"`
}
