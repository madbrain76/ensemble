// Package history owns the append-only log and every derived conversation view.
package history

import (
	"bufio"
	"encoding/json"
	"ensemble/internal/common"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type history struct {
	mu      sync.Mutex
	parent  common.Agent
	events  []common.Event
	context common.Context
}

func New(parent common.Agent) common.History {
	return &history{parent: parent, context: empty()}
}
func empty() common.Context {
	return common.Context{Turn: common.Idle, Dialogue: []common.Entry{}, Ephemera: common.PartList{}}
}
func (h *history) Parent() common.Agent    { return h.parent }
func clone[T any](v T) T                   { b, _ := json.Marshal(v); var c T; _ = json.Unmarshal(b, &c); return c }
func (h *history) Events() []common.Event  { h.mu.Lock(); defer h.mu.Unlock(); return clone(h.events) }
func (h *history) Context() common.Context { h.mu.Lock(); defer h.mu.Unlock(); return clone(h.context) }
func (h *history) Append(e common.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	next := common.Seq(1)
	if len(h.events) > 0 {
		next = h.events[len(h.events)-1].Seq + 1
	}
	if e.Seq == 0 {
		e.Seq = next
	}
	if e.Seq < next {
		return fmt.Errorf("sequence must increase")
	}
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	if err := validate(e); err != nil {
		return err
	}
	e = clone(e)
	h.events = append(h.events, e)
	Apply(&h.context, e)
	return nil
}
func (h *history) Load(r io.Reader) error {
	// Unbounded lines, optional version header, atomic publication on success.
	reader := bufio.NewReader(r)
	next := &history{parent: h.parent, context: empty()}
	n := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) > 0 {
			n++
			if strings.TrimSpace(line) == "" {
				return fmt.Errorf("empty history line %d", n)
			}
			var header struct {
				Version *int `json:"log_version"`
			}
			if x := json.Unmarshal([]byte(line), &header); x != nil {
				return x
			}
			if header.Version != nil {
				if n != 1 || *header.Version != 1 {
					return fmt.Errorf("unsupported history version")
				}
			} else {
				var e common.Event
				if x := json.Unmarshal([]byte(line), &e); x != nil {
					return fmt.Errorf("history line %d: %w", n, x)
				}
				if e.Seq == 0 {
					return fmt.Errorf("missing sequence")
				}
				if x := next.Append(e); x != nil {
					return fmt.Errorf("history line %d: %w", n, x)
				}
			}
		}
		if err == io.EOF {
			break
		}
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events, h.context = next.events, next.context
	return nil
}
func (h *history) Dump(w io.Writer) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(struct {
		Version int `json:"log_version"`
	}{1}); err != nil {
		return err
	}
	for _, e := range h.Events() {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}
func provenance(p common.Provenance) error {
	if p.Vendor.String() == "" || p.Surface.String() == "" {
		return fmt.Errorf("incomplete provenance")
	}
	return nil
}
func parts(ps common.PartList) error {
	for _, p := range ps {
		switch p.Type {
		case "text":
		case "blob", "image", "audio", "video":
			if p.Ref == nil {
				return fmt.Errorf("blob missing ref")
			}
		case "opaque":
			if p.From == nil || len(p.Data) == 0 {
				return fmt.Errorf("opaque missing provenance or data")
			}
		case "tool_call":
			if p.From == nil || p.Name == "" || len(p.Args) == 0 {
				return fmt.Errorf("incomplete tool call")
			}
		case "tool_result":
			if err := parts(p.Parts); err != nil {
				return err
			}
		case "redacted":
		default:
			return fmt.Errorf("unknown part %q", p.Type)
		}
		if p.Ref != nil {
			switch p.Ref.Kind {
			case 1, 2, 3:
			default:
				return fmt.Errorf("unknown ref kind")
			}
		}
		if p.From != nil {
			if err := provenance(*p.From); err != nil {
				return err
			}
		}
	}
	return nil
}
func validate(e common.Event) error {
	n := 0
	for _, yes := range []bool{e.Message != nil, e.Request != nil, e.Response != nil, e.Tool != nil, e.Redact != nil, e.Error != nil, e.Job != nil} {
		if yes {
			n++
		}
	}
	if n != 1 {
		return fmt.Errorf("event requires exactly one typed payload")
	}
	switch e.Type {
	case common.MessageReceived, common.HintReceived, common.TurnEnded:
		if e.Message != nil {
			switch e.Message.Actor {
			case common.Human, common.Assistant, common.System, common.Tool:
			default:
				return fmt.Errorf("unknown actor")
			}
			return parts(e.Message.Parts)
		}
	case common.RequestSent:
		if e.Request != nil {
			return provenance(e.Request.To)
		}
	case common.ResponseStarted, common.ResponseEnded:
		if e.Response != nil {
			if err := provenance(e.Response.From); err != nil {
				return err
			}
			u := e.Response.Usage
			if u.Input < 0 || u.Output < 0 || u.CacheRead < 0 || u.CacheWrite < 0 {
				return fmt.Errorf("negative usage")
			}
			return parts(e.Response.Parts)
		}
	case common.JobReported:
		if e.Job != nil && e.Job.Handle > 0 {
			return nil
		}
		if e.Tool != nil {
			return parts(e.Tool.Parts)
		}
	case common.JobKilled:
		if e.Job != nil && e.Job.Handle > 0 && e.Job.Status == common.JobKilledStatus && (e.Job.Reason == "kill_job" || e.Job.Reason == "shutdown") {
			return nil
		}
	case common.ToolCalled:
		if e.Tool != nil {
			return nil
		}
	case common.ToolReturned:
		if e.Tool != nil {
			return parts(e.Tool.Parts)
		}
	case common.Redacted:
		if e.Redact != nil {
			r := e.Redact
			if r.From == 0 || r.To < r.From {
				return fmt.Errorf("invalid redaction span")
			}
			switch r.Level {
			case common.RedactResult, common.RedactTool, common.RedactDialogue, common.RedactSummary:
				return parts(r.Replacement)
			}
		}
	case common.ErrorOccurred, common.TurnInterrupted:
		if e.Error != nil {
			return nil
		}
	default:
		return fmt.Errorf("unknown event type %q", e.Type)
	}
	return fmt.Errorf("invalid %s payload", e.Type)
}
func add(a *common.Usage, b common.Usage) {
	a.Input += b.Input
	a.Output += b.Output
	a.CacheRead += b.CacheRead
	a.CacheWrite += b.CacheWrite
}

// Apply is total for known events. Dialogue changes and turn transitions are
// separate: a known event can add a fact without changing this turn's state.
func Apply(c *common.Context, e common.Event) {
	if e.Stale {
		if e.Type == common.ResponseEnded && e.Response != nil {
			add(&c.Usage, e.Response.Usage)
		}
		return
	}
	switch e.Type {
	case common.MessageReceived, common.HintReceived:
		if e.Message.Actor == common.System {
			// The event log is append-only; this pending projection is a replacement.
			c.Ephemera = append(common.PartList(nil), e.Message.Parts...)
			if len(c.Ephemera) == 1 && c.Ephemera[0].Type == "text" && c.Ephemera[0].Text == "" {
				c.Ephemera = nil
			}
		} else {
			c.Dialogue = append(c.Dialogue, common.Entry{Seq: e.Seq, Actor: e.Message.Actor, Parts: e.Message.Parts})
		}
		if c.Turn == common.Idle {
			c.Turn = common.InputPending
		}
	case common.RequestSent:
		c.Ephemera = common.PartList{}
		c.Turn = common.InFlight
	case common.ResponseEnded:
		c.Dialogue = append(c.Dialogue, common.Entry{Seq: e.Seq, Actor: common.Assistant, Parts: e.Response.Parts})
		add(&c.Usage, e.Response.Usage)
		if c.Turn == common.InFlight {
			c.Pending = nil
			for _, p := range e.Response.Parts {
				if p.Type == "tool_call" {
					c.Pending = append(c.Pending, p.CallID)
				}
			}
			if len(c.Pending) > 0 {
				c.Turn = common.ToolsPending
			} else {
				c.Turn = common.Idle
			}
		}
	case common.ToolReturned:
		c.Dialogue = append(c.Dialogue, common.Entry{Seq: e.Seq, Actor: common.Tool, Parts: common.PartList{{Type: "tool_result", CallID: e.Tool.CallID, Parts: e.Tool.Parts, IsError: e.Tool.IsError}}})
		if c.Turn == common.ToolsPending {
			for i, id := range c.Pending {
				if id == e.Tool.CallID {
					c.Pending = append(c.Pending[:i], c.Pending[i+1:]...)
					break
				}
			}
			if len(c.Pending) == 0 {
				c.Turn = common.InputPending
			}
		}
	case common.TurnEnded:
		c.Turn = common.Idle
		c.Pending = nil
	case common.TurnInterrupted:
		c.Turn = common.Interrupted
	case common.ErrorOccurred:
		// An error changes turn state, not projected causal facts. Only an
		// explicit redaction may remove or replace retained dialogue.
		if c.Turn == common.InFlight {
			c.Turn = common.Idle
		}
	case common.Redacted:
		redact(c, *e.Redact)
	}
}
func stub(p common.Part) common.Part {
	size := len([]byte(p.Text))
	if p.Type != "text" {
		b, _ := json.Marshal(p)
		size = len(b)
	}
	return common.Part{Type: "redacted", Stub: fmt.Sprintf("[redacted %d bytes]", size), Ref: p.Ref}
}
func redact(c *common.Context, r common.RedactData) {
	if r.Level == common.RedactSummary {
		out := []common.Entry{}
		inserted := false
		for _, e := range c.Dialogue {
			if e.Seq >= r.From && e.Seq <= r.To {
				if !inserted {
					out = append(out, common.Entry{Seq: r.From, Actor: common.System, Parts: r.Replacement})
					inserted = true
				}
			} else {
				out = append(out, e)
			}
		}
		c.Dialogue = out
		return
	}
	for i := range c.Dialogue {
		e := &c.Dialogue[i]
		if e.Seq < r.From || e.Seq > r.To {
			continue
		}
		out := common.PartList{}
		for _, p := range e.Parts {
			switch r.Level {
			case common.RedactResult:
				if p.Type == "tool_result" {
					replacement := common.PartList{}
					for _, inner := range p.Parts {
						replacement = append(replacement, stub(inner))
					}
					p.Parts = replacement
				}
				out = append(out, p)
			case common.RedactTool:
				if p.Type != "tool_call" && p.Type != "tool_result" {
					out = append(out, p)
				}
			case common.RedactDialogue:
				if p.Type == "text" || p.Type == "opaque" || p.Type == "blob" {
					out = append(out, stub(p))
				} else {
					out = append(out, p)
				}
			}
		}
		e.Parts = out
	}
}
