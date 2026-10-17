// Package history owns append-only facts, serialization and their projection.
package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"ensemble/internal/common"
)

// History has one owning Agent and no vendor dependencies.
type History struct {
	mu     sync.Mutex
	parent common.Agent
	events []common.Event
}

// New constructs a History for an existing owning Agent.
func New(parent common.Agent) (*History, error) {
	if parent == nil {
		return nil, fmt.Errorf("history requires an agent parent")
	}
	return &History{parent: parent}, nil
}

// Parent returns the owner interface.
func (h *History) Parent() common.Agent { return h.parent }

// Append assigns sequence and time when absent, validates, then stores a deep
// copy. External snapshots and caller slices can never edit recorded facts.
func (h *History) Append(event common.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	var previous uint64
	if len(h.events) > 0 {
		previous = h.events[len(h.events)-1].Seq
	}
	if event.Seq == 0 {
		event.Seq = previous + 1
	}
	if event.Seq <= previous {
		return fmt.Errorf("sequence must increase")
	}
	if event.Time.IsZero() {
		event.Time = time.Now().UTC()
	}
	if err := validate(event); err != nil {
		return err
	}
	h.events = append(h.events, copyEvent(event))
	return nil
}

// Events returns a detached snapshot of original facts.
func (h *History) Events() []common.Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]common.Event, len(h.events))
	for i, e := range h.events {
		out[i] = copyEvent(e)
	}
	return out
}

// Dump emits a version header followed by ascending original events.
func (h *History) Dump(w io.Writer) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(struct {
		// Version identifies the supported JSONL serialization schema.
		Version int `json:"log_version"`
	}{1}); err != nil {
		return err
	}
	for _, event := range h.events {
		if err := encoder.Encode(event); err != nil {
			return err
		}
	}
	return nil
}

// Load validates a complete log before replacing this History. Headerless logs
// are current; unknown versions, variants and nonascending sequences fail loudly.
func (h *History) Load(r io.Reader) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	reader := bufio.NewReader(r)
	var events []common.Event
	var previous uint64
	lineNo := 0
	header := false
	for {
		line, readErr := reader.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		if len(strings.TrimSpace(string(line))) > 0 {
			lineNo++
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(line, &fields); err != nil {
				return fmt.Errorf("log line %d: %w", lineNo, err)
			}
			if raw, ok := fields["log_version"]; ok {
				var version int
				if err := json.Unmarshal(raw, &version); err != nil || version != 1 || lineNo != 1 || header || len(fields) != 1 {
					return fmt.Errorf("unsupported log version/header")
				}
				header = true
			} else {
				var event common.Event
				if err := json.Unmarshal(line, &event); err != nil {
					return fmt.Errorf("log line %d: %w", lineNo, err)
				}
				if event.Seq == 0 || event.Seq <= previous {
					return fmt.Errorf("log line %d: sequence must increase", lineNo)
				}
				if err := validate(event); err != nil {
					return fmt.Errorf("log line %d: %w", lineNo, err)
				}
				events = append(events, event)
				previous = event.Seq
			}
		}
		if readErr == io.EOF {
			break
		}
	}
	h.events = events
	return nil
}

// Copies are structural, so snapshots preserve even empty strings and never
// depend on a second serialization path or ignored marshal errors.
func copyParts(parts []common.Part) []common.Part {
	if parts == nil {
		return nil
	}
	out := make([]common.Part, len(parts))
	for i, p := range parts {
		if p.Ref != nil {
			ref := *p.Ref
			p.Ref = &ref
		}
		p.Args = append(json.RawMessage(nil), p.Args...)
		p.Opaque = append(json.RawMessage(nil), p.Opaque...)
		p.Data = append(json.RawMessage(nil), p.Data...)
		p.Parts = copyParts(p.Parts)
		out[i] = p
	}
	return out
}
func copyJob(j *common.JobRecord) *common.JobRecord {
	if j == nil {
		return nil
	}
	v := *j
	if j.ExitCode != nil {
		n := *j.ExitCode
		v.ExitCode = &n
	}
	return &v
}
func copyEvent(e common.Event) common.Event {
	e.Job = copyJob(e.Job)
	if e.Message != nil {
		v := *e.Message
		v.Parts = copyParts(v.Parts)
		e.Message = &v
	}
	if e.Request != nil {
		v := *e.Request
		e.Request = &v
	}
	if e.Response != nil {
		v := *e.Response
		v.Parts = copyParts(v.Parts)
		e.Response = &v
	}
	if e.Tool != nil {
		v := *e.Tool
		v.Job = copyJob(e.Tool.Job)
		v.Parts = copyParts(v.Parts)
		v.Args = append(json.RawMessage(nil), v.Args...)
		e.Tool = &v
	}
	if e.Error != nil {
		v := *e.Error
		e.Error = &v
	}
	if e.Redact != nil {
		v := *e.Redact
		v.Replacement = copyParts(v.Replacement)
		e.Redact = &v
	}
	return e
}

func provenance(p common.Provenance) error {
	expected := common.Messages
	switch p.Vendor {
	case common.Anthropic:
	case common.OpenAI:
		expected = common.ChatCompletions
	case common.Gemini:
		expected = common.GenerateContent
	default:
		return fmt.Errorf("unknown vendor")
	}
	if p.Surface != expected || p.Model == "" {
		return fmt.Errorf("invalid model or surface")
	}
	return nil
}
func validate(event common.Event) error {
	count := 0
	for _, present := range []bool{event.Message != nil, event.Request != nil, event.Response != nil, event.Tool != nil, event.Error != nil, event.Redact != nil, event.Job != nil} {
		if present {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("event requires exactly one payload")
	}
	switch event.Type {
	case "message_received", "hint_received":
		if event.Message == nil {
			return fmt.Errorf("message_received requires message")
		}
		switch event.Message.Actor {
		case "human", "agent", "system", "tool":
		default:
			return fmt.Errorf("unknown actor %q", event.Message.Actor)
		}
		return partsValid(event.Message.Parts)
	case "request_sent":
		if event.Request == nil {
			return fmt.Errorf("request_sent requires request")
		}
		return provenance(event.Request.To)
	case "response_started", "response_ended":
		if event.Response == nil {
			return fmt.Errorf("response requires response")
		}
		if err := provenance(event.Response.From); err != nil {
			return err
		}
		u := event.Response.Usage
		if u.Input < 0 || u.Output < 0 || u.CacheWrite < 0 || u.CacheRead < 0 {
			return fmt.Errorf("negative usage")
		}
		return partsValid(event.Response.Parts)
	case "tool_called", "tool_returned":
		if event.Tool == nil {
			return fmt.Errorf("tool event requires tool")
		}
		if len(event.Tool.Args) > 0 && !json.Valid(event.Tool.Args) {
			return fmt.Errorf("invalid tool args")
		}
		return partsValid(event.Tool.Parts)
	case "job_report":
		if event.Job == nil {
			return fmt.Errorf("job_report requires job")
		}
	case "turn_ended":
		if event.Error == nil {
			return fmt.Errorf("turn_ended requires outcome")
		}
	case "job_killed":
		if event.Job == nil || event.Job.Status != "killed" || (event.Job.Reason != "kill_job" && event.Job.Reason != "shutdown" && event.Job.Reason != "cancellation") {
			return fmt.Errorf("invalid killed job")
		}
	case "error_occurred":
		if event.Error == nil {
			return fmt.Errorf("error_occurred requires error")
		}
	case "redacted":
		if event.Redact == nil {
			return fmt.Errorf("redacted requires redact")
		}
		r := event.Redact
		if r.From == 0 || r.To < r.From || r.To >= event.Seq {
			return fmt.Errorf("invalid redaction span")
		}
		switch r.Level {
		case "redact_result", "redact_tool", "redact_dialogue":
			if len(r.Replacement) > 0 {
				return fmt.Errorf("replacement requires redact_summary")
			}
		case "redact_summary":
			return partsValid(r.Replacement)
		default:
			return fmt.Errorf("unknown redaction level %q", r.Level)
		}
	default:
		return fmt.Errorf("unknown event type %q", event.Type)
	}
	return nil
}
func partsValid(parts []common.Part) error {
	for _, part := range parts {
		if part.From != (common.Provenance{}) {
			if err := provenance(part.From); err != nil {
				return err
			}
		}
		for _, raw := range []json.RawMessage{part.Args, part.Opaque, part.Data} {
			if len(raw) > 0 && !json.Valid(raw) {
				return fmt.Errorf("invalid part JSON")
			}
		}
		switch part.Type {
		case "text", "redacted":
		case "blob":
			if part.Ref == nil || part.Ref.Kind < 1 || part.Ref.Kind > 3 || part.Ref.Locator == "" {
				return fmt.Errorf("blob requires ref kind 1, 2 or 3 and locator")
			}
		case "opaque":
			if err := provenance(part.From); err != nil {
				return err
			}
			if !json.Valid(part.Data) {
				return fmt.Errorf("opaque data is not JSON")
			}
		case "tool_call":
			if err := provenance(part.From); err != nil {
				return err
			}
			if part.Name == "" || !json.Valid(part.Args) {
				return fmt.Errorf("tool call requires name and JSON args")
			}
			if len(part.Opaque) > 0 && !json.Valid(part.Opaque) {
				return fmt.Errorf("tool opaque is not JSON")
			}
		case "tool_result":
			if err := partsValid(part.Parts); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown part type %q", part.Type)
		}
	}
	return nil
}

// Context always replays detached facts; no cached mutable projection can become
// a second source of truth. Turn-state defaults are identity for known events.
func (h *History) Context() common.Context {
	c := common.Context{Turn: common.Idle, Dialogue: []common.Turn{}, Ephemera: []common.Part{}}
	for _, event := range h.Events() {
		apply(&c, event)
	}
	return c
}

// messageParts projects supplied addressing as explicit context without changing
// the sender's parts or confusing application names with provider roles.
func messageParts(m *common.MessageEvent) []common.Part {
	if m.Sender == "" && m.Recipient == "" && m.RequestID == "" && m.TurnID == "" {
		return m.Parts
	}
	data, _ := json.Marshal(struct {
		Sender    string `json:"sender,omitempty"`
		Recipient string `json:"recipient,omitempty"`
		RequestID string `json:"request_id,omitempty"`
		TurnID    string `json:"turn_id,omitempty"`
	}{m.Sender, m.Recipient, m.RequestID, m.TurnID})
	return append([]common.Part{{Type: "text", Text: "Message addressing: " + string(data)}}, m.Parts...)
}
func apply(c *common.Context, event common.Event) {
	switch event.Type {
	case "message_received":
		if event.Message.Actor == "system" {
			c.Ephemera = append(c.Ephemera, messageParts(event.Message)...)
		} else {
			c.Dialogue = append(c.Dialogue, common.Turn{Seq: event.Seq, Actor: event.Message.Actor, Parts: messageParts(event.Message)})
		}
		if c.Turn == common.Idle || c.Turn == common.Interrupted {
			c.Turn = common.InputPending
		}
	case "hint_received":
		c.Ephemera = append(c.Ephemera, messageParts(event.Message)...)
	case "turn_ended":
		c.Turn = common.Idle
		if event.Error.Message == common.ErrInterrupted.Error() {
			c.Turn = common.Interrupted
		}
	case "request_sent":
		c.Ephemera = []common.Part{}
		if c.Turn == common.InputPending {
			c.Turn = common.InFlight
		}
	case "response_ended":
		c.Dialogue = append(c.Dialogue, common.Turn{Seq: event.Seq, Actor: "agent", Parts: event.Response.Parts})
		u := event.Response.Usage
		c.Usage.Input += u.Input
		c.Usage.CacheWrite += u.CacheWrite
		c.Usage.CacheRead += u.CacheRead
		c.Usage.Output += u.Output
		if c.Turn == common.InFlight {
			c.Turn = common.Idle
			for _, p := range event.Response.Parts {
				if p.Type == "tool_call" {
					c.Turn = common.ToolsPending
					break
				}
			}
		}
	case "tool_returned":
		t := event.Tool
		c.Dialogue = append(c.Dialogue, common.Turn{Seq: event.Seq, Actor: "tool", Parts: []common.Part{{Type: "tool_result", CallID: t.CallID, Parts: t.Parts, IsError: t.IsError}}})
		if c.Turn == common.ToolsPending && !pendingTools(c.Dialogue) {
			c.Turn = common.InputPending
		}
	case "error_occurred":
		if c.Turn == common.InFlight {
			c.Turn = common.Idle
		}
	case "redacted":
		redact(c, event.Redact)
	}
}

// pendingTools is derived from the most recent agent contribution, not an
// ever-growing side map in Context. Only this synchronous turn is relevant.
func pendingTools(dialogue []common.Turn) bool {
	returned := map[string]bool{}
	for i := len(dialogue) - 1; i >= 0; i-- {
		for _, p := range dialogue[i].Parts {
			if p.Type == "tool_result" {
				returned[p.CallID] = true
			}
		}
		if dialogue[i].Actor == "agent" {
			for _, p := range dialogue[i].Parts {
				if p.Type == "tool_call" && !returned[p.CallID] {
					return true
				}
			}
			break
		}
	}
	return false
}
func stub(p common.Part) common.Part {
	size := len(p.Text)
	if p.Type != "text" {
		data, _ := json.Marshal(p)
		size = len(data)
	}
	return common.Part{Type: "redacted", Stub: fmt.Sprintf("[redacted %s: %d bytes]", p.Type, size), Ref: p.Ref}
}
func redact(c *common.Context, r *common.Redaction) {
	out := make([]common.Turn, 0, len(c.Dialogue))
	inserted := false
	for _, turn := range c.Dialogue {
		if turn.Seq < r.From || turn.Seq > r.To {
			out = append(out, turn)
			continue
		}
		if r.Level == "redact_summary" {
			if !inserted {
				out = append(out, common.Turn{Seq: r.From, Actor: "system", Parts: copyParts(r.Replacement)})
				inserted = true
			}
			continue
		}
		parts := []common.Part{}
		for _, p := range turn.Parts {
			switch r.Level {
			case "redact_result":
				if p.Type == "tool_result" {
					replaced := []common.Part{}
					for _, inner := range p.Parts {
						replaced = append(replaced, stub(inner))
					}
					p.Parts = replaced
				}
				parts = append(parts, p)
			case "redact_tool":
				if p.Type != "tool_call" && p.Type != "tool_result" {
					parts = append(parts, p)
				}
			case "redact_dialogue":
				if p.Type == "tool_call" || p.Type == "tool_result" {
					parts = append(parts, p)
				} else {
					parts = append(parts, stub(p))
				}
			}
		}
		turn.Parts = parts
		out = append(out, turn)
	}
	c.Dialogue = out
}
