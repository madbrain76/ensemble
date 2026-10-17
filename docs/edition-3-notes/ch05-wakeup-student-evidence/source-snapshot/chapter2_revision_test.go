package ensemble_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"ensemble"
)

// TestChapter2GeminiChronologicalNames uses ordinary sequential ID-less calls.
// The second call must not rename the first result when the full log is replayed.
func TestChapter2GeminiChronologicalNames(t *testing.T) {
	a := offlineAgent(t)
	from := issuer(ensemble.Gemini)
	for _, name := range []string{"lookup_weather", "lookup_stock"} {
		body := fmt.Sprintf(`{"candidates":[{"content":{"parts":[{"functionCall":{"name":%q,"args":{}}}]}}],"usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1}}`, name)
		parsed, err := a.Engine().Parse([]byte(body), from)
		if err != nil {
			t.Fatal(err)
		}
		if len(parsed.Parts) != 1 || parsed.Parts[0].CallID != "" {
			t.Fatalf("native absent ID changed: %+v", parsed.Parts)
		}
		appendFact(t, a, ensemble.Event{Type: "response_ended", Response: &parsed})
		appendFact(t, a, ensemble.Event{Type: "tool_returned", Tool: &ensemble.ToolEvent{Parts: []ensemble.Part{{Type: "text", Text: name + " result"}}}})
	}
	var saved bytes.Buffer
	if err := a.Log().Dump(&saved); err != nil {
		t.Fatal(err)
	}
	b := offlineAgent(t)
	if err := b.Log().Load(bytes.NewReader(saved.Bytes())); err != nil {
		t.Fatal(err)
	}
	before := b.Log().Events()
	wire, err := b.Engine().Render(b.Log().Context(), from)
	if err != nil {
		t.Fatal(err)
	}
	var request struct {
		Contents []struct {
			Parts []struct {
				Response *struct {
					Name string `json:"name"`
				} `json:"functionResponse"`
			} `json:"parts"`
		} `json:"contents"`
	}
	if err := json.Unmarshal(wire, &request); err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, turn := range request.Contents {
		for _, part := range turn.Parts {
			if part.Response != nil {
				names = append(names, part.Response.Name)
			}
		}
	}
	if !reflect.DeepEqual(names, []string{"lookup_weather", "lookup_stock"}) {
		t.Fatalf("result names must match earlier calls, got %v: %s", names, wire)
	}
	again, err := b.Engine().Render(b.Log().Context(), from)
	if err != nil || !bytes.Equal(wire, again) || !reflect.DeepEqual(before, b.Log().Events()) {
		t.Fatal("render changed facts or chronological resolution")
	}
}

// TestChapter2MeasuredEmptyReplies exercises actual root-owned HTTP transport.
// Attempt accounting includes measured failures; reduced usage and committed
// responses do not. Each Ask still performs only one request and returns an error.
func TestChapter2MeasuredEmptyReplies(t *testing.T) {
	for _, tc := range []struct {
		vendor            ensemble.Vendor
		path, good, empty string
	}{
		{ensemble.Anthropic, "/v1/messages", `{"content":[{"type":"text","text":"ok"}],"usage":{"input_tokens":3,"output_tokens":1}}`, `{"content":[],"usage":{"input_tokens":11,"output_tokens":2}}`},
		{ensemble.OpenAI, "/v1/chat/completions", `{"choices":[{"message":{"content":"ok"}}],"usage":{"prompt_tokens":3,"completion_tokens":1}}`, `{"choices":[],"usage":{"prompt_tokens":11,"completion_tokens":2}}`},
		{ensemble.Gemini, "/v1beta/models/exact-model:generateContent", `{"candidates":[{"content":{"parts":[{"text":"ok"}]}}],"usageMetadata":{"promptTokenCount":3,"candidatesTokenCount":1}}`, `{"candidates":[],"usageMetadata":{"promptTokenCount":11,"candidatesTokenCount":2}}`},
	} {
		t.Run(tc.vendor.String(), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != http.MethodPost || r.URL.Path != tc.path {
					t.Errorf("unexpected wire request: %s %s", r.Method, r.URL.Path)
				}
				if calls == 1 {
					fmt.Fprint(w, tc.good)
				} else {
					fmt.Fprint(w, tc.empty)
				}
			}))
			defer server.Close()
			root := ensemble.New(log.New(io.Discard, "", 0))
			a, err := root.NewAgent(ensemble.Config{Vendor: tc.vendor, Model: "exact-model", BaseURL: server.URL, APIKey: "fake-key", System: "test", MaxTokens: 32})
			if err != nil {
				t.Fatal(err)
			}
			if answer, err := a.Ask(context.Background(), "first"); err != nil || answer != "ok" {
				t.Fatalf("successful baseline: %q %v", answer, err)
			}
			committed := a.History()
			for failures := 1; failures <= 2; failures++ {
				if answer, err := a.Ask(context.Background(), "empty reply"); err == nil || answer != "" {
					t.Fatalf("empty reply accepted: %q %v", answer, err)
				}
				want := ensemble.Usage{Input: 3 + failures*11, Output: 1 + failures*2}
				if got := a.Engine().Usage(); got != want {
					t.Errorf("measured attempt usage: got %+v, want %+v", got, want)
				}
				if calls != failures+1 || !reflect.DeepEqual(committed, a.History()) {
					t.Fatal("failed request retried or committed")
				}
				if got := a.Log().Context().Usage; got != (ensemble.Usage{Input: 3, Output: 1}) {
					t.Fatalf("attempt cost fabricated a completed response: %+v", got)
				}
				completed, errors := 0, 0
				for _, event := range a.Log().Events() {
					if event.Type == "response_ended" {
						completed++
					}
					if event.Type == "error_occurred" {
						errors++
					}
				}
				if completed != 1 || errors != failures {
					t.Fatalf("capture: %d completed responses, %d errors", completed, errors)
				}
			}
		})
	}
}

// TestChapter2OpenAIMultipleOpaqueParts checks both already-supported fields
// through parsing, durable facts, and reduction, rather than direct wire assembly.
func TestChapter2OpenAIMultipleOpaqueParts(t *testing.T) {
	a := offlineAgent(t)
	from := issuer(ensemble.OpenAI)
	parsed, err := a.Engine().Parse([]byte(`{"choices":[{"message":{"content":"answer","reasoning_content":"opaque-first","reasoning":{"text":"opaque-second"}}}],"usage":{"prompt_tokens":4,"completion_tokens":2}}`), from)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Parts) != 3 || parsed.Parts[1].Type != "opaque" || parsed.Parts[2].Type != "opaque" {
		t.Fatalf("expected two independent opaque parts: %+v", parsed.Parts)
	}
	appendFact(t, a, ensemble.Event{Type: "response_ended", Response: &parsed})
	var saved bytes.Buffer
	if err := a.Log().Dump(&saved); err != nil {
		t.Fatal(err)
	}
	b := offlineAgent(t)
	if err := b.Log().Load(bytes.NewReader(saved.Bytes())); err != nil {
		t.Fatal(err)
	}
	before := b.Log().Events()
	for _, reverse := range []bool{false, true} {
		c := b.Log().Context()
		if reverse {
			c.Dialogue[0].Parts[1], c.Dialogue[0].Parts[2] = c.Dialogue[0].Parts[2], c.Dialogue[0].Parts[1]
		}
		wire, err := b.Engine().Render(c, from)
		if err != nil {
			t.Fatal(err)
		}
		var request struct {
			Messages []struct {
				Role             string          `json:"role"`
				Content          string          `json:"content"`
				ReasoningContent json.RawMessage `json:"reasoning_content"`
				Reasoning        json.RawMessage `json:"reasoning"`
			} `json:"messages"`
		}
		if err := json.Unmarshal(wire, &request); err != nil {
			t.Fatal(err)
		}
		if len(request.Messages) != 2 {
			t.Fatalf("messages: %s", wire)
		}
		message := request.Messages[1]
		if message.Role != "assistant" || message.Content != "answer" || string(message.ReasoningContent) != `"opaque-first"` || string(message.Reasoning) != `{"text":"opaque-second"}` {
			t.Fatalf("lost opaque field (reversed=%v): %s", reverse, wire)
		}
	}
	for _, mismatch := range []string{"vendor", "model", "surface"} {
		c := b.Log().Context()
		for i := range c.Dialogue[0].Parts {
			p := &c.Dialogue[0].Parts[i]
			if p.Type != "opaque" {
				continue
			}
			switch mismatch {
			case "vendor":
				p.From.Vendor = ensemble.Anthropic
			case "model":
				p.From.Model = "other-model"
			case "surface":
				p.From.Surface = ensemble.Messages
			}
		}
		wire, err := b.Engine().Render(c, from)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(wire, []byte("opaque-first")) || bytes.Contains(wire, []byte("opaque-second")) {
			t.Fatalf("opaque replay crossed %s: %s", mismatch, wire)
		}
	}
	if !reflect.DeepEqual(before, b.Log().Events()) {
		t.Fatal("render mutated stored provenance")
	}
}
