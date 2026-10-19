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
)

const systemPrompt = "You are a helpful assistant built from raw HTTP calls in Chapter 1 of The Art of Building AI Coding Agents. Answer briefly."

type engine struct {
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
func (e *engine) Usage() common.Usage  { return e.liveUsage }
func (e *engine) failure(message string) error {
	e.parent.Parent().Logger().Print(message)
	return errors.New(message)
}
func (e *engine) Render(c common.Context) ([]byte, error) { return render(c, e.parent.Config()) }
func (e *engine) Parse(status int, b []byte) ([]common.Event, error) {
	return parse(status, b, e.parent.Config())
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
		endpoint += "/models/" + url.PathEscape(config.Model) + ":generateContent"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return []common.Event{errorEvent("creating request failed")}, e.failure("creating request failed")
	}
	req.Header.Set("content-type", "application/json")
	switch config.Vendor {
	case common.Anthropic:
		req.Header.Set("x-api-key", config.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	case common.OpenAI:
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
	case common.Gemini:
		req.Header.Set("x-goog-api-key", config.APIKey)
	}
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
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		err = e.failure("reading response failed")
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return []common.Event{errorEvent(err.Error())}, err
	}
	events, err := e.Parse(resp.StatusCode, raw)
	// The same measured response facts feed two distinct owners: this process's
	// HTTP bill here, and replay totals when Agent appends them to History.
	for _, event := range events {
		if event.Type == common.ResponseEnded {
			u := event.Response.Usage
			e.liveUsage.Input += u.Input
			e.liveUsage.Output += u.Output
			e.liveUsage.CacheRead += u.CacheRead
			e.liveUsage.CacheWrite += u.CacheWrite
		}
	}
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
