// Package engine owns HTTP and vendor translation, never retained dialogue.
package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"ensemble/internal/common"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const systemPrompt = "You are a helpful coding assistant. Use the available tools when needed. Base claims about files and commands on returned results, and answer clearly."

type engine struct {
	mu        sync.Mutex
	parent    common.Agent
	client    *http.Client
	liveUsage common.Usage // Only Exchange charges this process; Parse and replay never do.
}

func New(parent common.Agent) (common.Engine, error) {
	if parent == nil {
		return nil, errors.New("Engine requires an Agent parent")
	}
	return &engine{parent: parent, client: &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}
func (e *engine) Parent() common.Agent { return e.parent }
func (e *engine) Usage() common.Usage  { e.mu.Lock(); defer e.mu.Unlock(); return e.liveUsage }
func (e *engine) failure(message string) error {
	e.parent.Parent().Logger().Print(message)
	return errors.New(message)
}
func (e *engine) Render(c common.Context) ([]byte, error) {
	cfg := e.parent.Config()
	snapshot := e.parent.Skills().Snapshot()
	if len(c.ActiveSkills) == 0 {
		cfg.SystemPrompt += snapshot.Instructions
	}
	cfg.SystemPrompt += offers(e.parent.Skills().LoadableSkills())
	disclosure := anthropicDisclosure{}
	if cfg.Vendor == common.Anthropic && cfg.Skills != nil && supportsToolReferences(cfg.Model) {
		disclosure.catalog = make([]common.ToolSpec, 0, len(snapshot.Tools))
		disclosure.references = map[string][]string{}
		for _, t := range snapshot.Tools {
			// Tools retains capability truth for admission and automatic gathering.
			// Automatic tools never enter declarations or enabling-call references.
			// Dormant compiled tools remain deferred; external declarations require
			// an enabled grant, even while their skill's instructions remain pending.
			if t.Ephemeral != "" || (!t.Enabled && (strings.HasPrefix(t.Source, "mcp:") || t.Skill != "")) {
				continue
			}
			disclosure.catalog = append(disclosure.catalog, t)
			if !t.Enabled || t.Origin != "tool" || t.EventSeq == 0 {
				continue
			}
			id := ""
			// Only live projected provenance can supply a native reference.
			// Never hydrate the lifetime archive during ordinary rendering.
			for _, entry := range c.Dialogue {
				if entry.Actor != common.Assistant {
					continue
				}
				for _, part := range entry.Parts {
					if part.Type == "tool_call" && part.Name == "load_skill" && part.AdmissionSeq == t.EventSeq && part.From != nil && part.From.Vendor == cfg.Vendor && part.From.Surface == cfg.Surface {
						id = part.CallID
					}
				}
			}
			if id == "" {
				continue
			}
			present := false
			for _, entry := range c.Dialogue {
				for _, part := range entry.Parts {
					if entry.Seq > t.EventSeq && part.Type == "tool_result" && part.CallID == id && !part.IsError {
						present = true
					}
				}
			}
			if !present {
				continue
			}

			disclosure.references[id] = append(disclosure.references[id], t.Name)
		}
	}
	return render(c, cfg, snapshot.Definitions, disclosure)
}

// EffectiveSystemPrompt uses the same system composition as all renderers.
func (e *engine) EffectiveSystemPrompt(c common.Context) (string, error) {
	cfg := e.parent.Config()
	snapshot := e.parent.Skills().Snapshot()
	if len(c.ActiveSkills) == 0 {
		cfg.SystemPrompt += snapshot.Instructions
	}
	cfg.SystemPrompt += offers(e.parent.Skills().LoadableSkills())
	c, cfg, _ = projectPrompt(c, cfg, snapshot.Definitions)
	return system(c, cfg)
}

// Parse owns transport framing only. The caller owns resp.Body and callbacks;
// no retained History, observation subscription, or reply channel lives here.
func (e *engine) Parse(resp *http.Response, callbacks common.ParseCallbacks) (events []common.Event, err error) {
	defer func() {
		if callbacks.OnEvent != nil {
			for _, event := range events {
				callbacks.OnEvent(event)
			}
		}
	}()
	if resp == nil || resp.Body == nil {
		return parseFailure("missing response body")
	}
	if resp.StatusCode != http.StatusOK {
		return parseFailure(fmt.Sprintf("HTTP status %d", resp.StatusCode))
	}
	if strings.EqualFold(strings.TrimSpace(strings.Split(resp.Header.Get("Content-Type"), ";")[0]), "text/event-stream") {
		events, err = parseSSE(resp.Body, e.parent.Config(), callbacks)
	} else {
		var raw []byte
		raw, err = io.ReadAll(resp.Body)
		if err != nil {
			return parseFailure("reading response failed")
		}
		if callbacks.OnFrame != nil {
			callbacks.OnFrame(common.Frame{Data: string(raw)})
		}
		events, err = parse(resp.StatusCode, raw, e.parent.Config())
		for _, event := range events {
			if event.Type == common.ResponseEnded {
				for _, p := range event.Response.Parts {
					plainDelta(p, callbacks)
				}
			}
		}
	}
	return events, err
}
func errorEvent(message string) common.Event {
	return common.Event{Type: common.ErrorOccurred, Error: &common.ErrorData{Message: message}}
}

// Exchange preserves caller-owned header/body cancellation and refuses redirects.
// Only locally chosen diagnostics escape: bodies and transport errors may hold keys.
func (e *engine) Exchange(ctx context.Context, c common.Context) ([]common.Event, error) {
	b, err := e.Render(c)
	if err != nil {
		return []common.Event{errorEvent(err.Error())}, err
	}
	return e.ExchangeRequest(ctx, b)
}

// ExchangeRequest sends the immutable bytes frozen by the actor before its
// request_sent marker. The Engine performs exactly one exchange.
func (e *engine) ExchangeRequest(ctx context.Context, b []byte, callbacks ...common.ParseCallbacks) ([]common.Event, error) {
	var err error
	config := e.parent.Config()
	if strings.TrimSpace(config.APIKey) == "" {
		err = e.failure("API key is not set")
		return []common.Event{errorEvent(err.Error())}, err
	}
	endpoint := config.BaseURL
	switch config.Vendor {
	case common.Anthropic:
		endpoint += "/v1/messages"
	case common.OpenAI:
		if strings.HasSuffix(endpoint, "/v1") {
			endpoint += "/chat/completions"
		} else {
			endpoint += "/v1/chat/completions"
		}
	case common.Gemini:
		if !strings.HasSuffix(endpoint, "/v1beta") {
			endpoint += "/v1beta"
		}
		endpoint += "/models/" + url.PathEscape(config.Model)
		if config.Streaming() {
			endpoint += ":streamGenerateContent?alt=sse"
		} else {
			endpoint += ":generateContent"
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return []common.Event{errorEvent("creating request failed")}, e.failure("creating request failed")
	}
	req.Header.Set("content-type", "application/json")
	if config.Streaming() {
		req.Header.Set("Accept", "text/event-stream")
	}
	switch config.Vendor {
	case common.Anthropic:
		req.Header.Set("x-api-key", config.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")
		if inlineFeatures(config) {
			req.Header.Set("anthropic-beta", "inline-tools-2026-09-15")
		}
	case common.OpenAI:
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
	case common.Gemini:
		req.Header.Set("x-goog-api-key", config.APIKey)
	}
	e.parent.Parent().Trace(common.TraceRecord{Stage: "request_body", Payload: string(b)})
	resp, err := e.client.Do(req)
	if err != nil {
		err = e.failure("HTTP request failed")
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return []common.Event{errorEvent(err.Error())}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = e.failure(fmt.Sprintf("HTTP status %d", resp.StatusCode))
		return []common.Event{errorEvent(err.Error())}, err
	}
	cb := common.ParseCallbacks{}
	if len(callbacks) > 0 {
		cb = callbacks[0]
	}
	originalDelta := cb.OnDelta
	cb.OnDelta = func(d common.PartDelta) {
		e.parent.Parent().Trace(common.TraceRecord{Stage: "parse_delta", PartID: d.PartID, Kind: string(d.Kind), Bytes: len(d.Chunk)})
		if originalDelta != nil {
			originalDelta(d)
		}
	}
	originalFrame := cb.OnFrame
	cb.OnFrame = func(f common.Frame) {
		e.parent.Parent().Trace(common.TraceRecord{Stage: "response_frame", Event: f.Event, Payload: f.Data})
		if originalFrame != nil {
			originalFrame(f)
		}
	}
	e.parent.Parent().Trace(common.TraceRecord{Stage: "response_parse"})
	events, err := e.Parse(resp, cb)
	if ctx.Err() != nil {
		err = ctx.Err()
	}

	// The same measured response facts feed two distinct owners: this process's
	// HTTP bill here, and replay totals when Agent appends them to History.
	e.mu.Lock()
	for _, event := range events {
		if event.Type == common.ResponseEnded {
			u := event.Response.Usage
			e.liveUsage.Input += u.Input
			e.liveUsage.Output += u.Output
			e.liveUsage.CacheRead += u.CacheRead
			e.liveUsage.CacheWrite += u.CacheWrite
		}
	}
	e.mu.Unlock()
	if err != nil {
		e.parent.Parent().Logger().Print(err.Error())
	}
	return events, err
}

// Send preserves the accepted text API as a disposable view, not stored state.
func (e *engine) Send(ctx context.Context, messages []common.Message) (string, error) {
	c := common.Context{Turn: common.InputPending}
	for i, m := range messages {
		actor := common.Human
		if m.Role == "assistant" {
			actor = common.Assistant
		}
		c.Dialogue = append(c.Dialogue, common.Entry{Seq: common.Seq(i + 1), Actor: actor, Parts: common.PartList{common.Text(m.Content)}})
	}
	events, err := e.Exchange(ctx, c)
	return responseText(events), err
}
func responseText(events []common.Event) string {
	var b strings.Builder
	for _, e := range events {
		if e.Type == common.ResponseEnded {
			for _, p := range e.Response.Parts {
				if p.Type == "text" {
					b.WriteString(p.Text)
				}
			}
		}
	}
	return b.String()
}

// Canonical JSON preserves integers beyond float64's exact range.
func canonical(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return json.RawMessage(`{}`), nil
	}
	var v any
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return nil, errors.New("invalid JSON value")
	}
	return json.Marshal(v)
}

// Plain delivery uses exactly the same notification vocabulary, in one piece.
func plainDelta(p common.Part, cb common.ParseCallbacks) {
	if cb.OnDelta == nil {
		return
	}
	d := common.PartDelta{PartID: p.ID}
	switch p.Type {
	case "text":
		d.Kind = common.DeltaText
		d.Chunk = p.Text
	case "tool_call":
		d.Kind = common.DeltaToolCall
		d.Name = p.Name
		d.Chunk = p.Name + string(p.Args)
	case "opaque":
		var fields map[string]json.RawMessage
		if json.Unmarshal(p.Data, &fields) != nil {
			return
		}
		switch {
		case string(fields["type"]) == `"thinking"`:
			d.Chunk = str(fields["thinking"])
		case string(fields["thought"]) == "true":
			d.Chunk = str(fields["text"])
		case fields["reasoning_content"] != nil:
			d.Chunk = str(fields["reasoning_content"])
		case fields["reasoning"] != nil:
			d.Chunk = str(fields["reasoning"])
		default:
			return
		}
		d.Kind = common.DeltaThinking
	default:
		return
	}
	cb.OnDelta(d)
}
