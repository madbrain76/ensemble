package ensemble_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"ensemble"
)

func newAgent(t *testing.T, handler http.HandlerFunc) (*ensemble.Ensemble, ensemble.Agent, *bytes.Buffer) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	logs := new(bytes.Buffer)
	root := ensemble.New(log.New(logs, "", 0))
	a, err := root.NewAgent(ensemble.Config{BaseURL: server.URL, APIKey: "test-key", Model: "test-model", System: "Be precise.", MaxTokens: 123})
	if err != nil {
		t.Fatal(err)
	}
	return root, a, logs
}

// TestConversation checks real wire data, owners, snapshots and exact replay.
func TestConversation(t *testing.T) {
	var requests [][]ensemble.Message
	root, a, logs := newAgent(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/messages" || r.Header.Get("x-api-key") != "test-key" || r.Header.Get("anthropic-version") != "2023-06-01" || r.Header.Get("content-type") != "application/json" {
			t.Error("wire method/path/headers")
		}
		var body struct {
			Model     string            `json:"model"`
			MaxTokens int               `json:"max_tokens"`
			System    string            `json:"system"`
			Stream    bool              `json:"stream"`
			Messages  []ensembleMessage `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.Model != "test-model" || body.MaxTokens != 123 || body.System != "Be precise." || body.Stream {
			t.Error("wire configuration")
		}
		var messages []ensemble.Message
		for _, m := range body.Messages {
			messages = append(messages, ensemble.Message{Role: m.Role, Content: m.Content})
		}
		requests = append(requests, messages)
		// Two text blocks must be concatenated exactly, without added separators.
		fmt.Fprint(w, `{"content":[{"type":"text","text":"hello "},{"type":"text","text":"世界"}],"usage":{"input_tokens":17,"output_tokens":9}}`)
	})
	if a.Parent() != root || a.Engine().Parent() != a || len(root.Agents()) != 1 || root.Agents()[0] != a {
		t.Fatal("owner chain")
	}
	config := a.Config()
	config.Model = "changed"
	agents := root.Agents()
	agents[0] = nil
	if a.Config().Model != "test-model" || root.Agents()[0] != a {
		t.Fatal("configuration/root snapshot aliases")
	}
	var want []ensemble.Message
	for _, user := range []string{"My boat is Juniper Arc.", "How are you?", "Remind me.", strings.Repeat("x", 70000)} {
		answer, err := a.Ask(context.Background(), user)
		if err != nil || answer != "hello 世界" {
			t.Fatalf("all text blocks: %q %v", answer, err)
		}
		want = append(want, ensemble.Message{Role: "user", Content: user})
		if !reflect.DeepEqual(requests[len(requests)-1], want) {
			t.Fatal("full exact history not resent")
		}
		want = append(want, ensemble.Message{Role: "assistant", Content: answer})
		if !reflect.DeepEqual(a.History(), want) {
			t.Fatal("committed history")
		}
	}
	if len(requests) != 4 {
		t.Fatalf("one request per turn: %d", len(requests))
	}
	if a.Engine().Usage() != (ensemble.Usage{Input: 68, Output: 36}) {
		t.Fatal("usage must sum")
	}
	history := a.History()
	history[0].Content = "changed"
	if a.History()[0].Content != "My boat is Juniper Arc." {
		t.Fatal("history snapshot aliases")
	}
	if !strings.Contains(logs.String(), "input=17 output=9") {
		t.Fatal("parent logger not reached")
	}
	// Multiple Agents share services but not state.
	other, err := root.NewAgent(a.Config())
	if err != nil || other.Parent() != root || other.Engine().Parent() != other || len(root.Agents()) != 2 || len(other.History()) != 0 || other.Engine().Usage() != (ensemble.Usage{}) {
		t.Fatal("Agent isolation")
	}
}

// TestFailuresDoNotCommit separates measured cost from committed history.
func TestFailuresDoNotCommit(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		measured   ensemble.Usage
	}{
		{"status", `{"error":"no"}`, 429, ensemble.Usage{}},
		{"missing usage", `{"content":[{"type":"text","text":"unaccounted"}]}`, 200, ensemble.Usage{}},
		{"missing output", `{"content":[{"type":"text","text":"unaccounted"}],"usage":{"input_tokens":0}}`, 200, ensemble.Usage{}},
		{"malformed", `{"content":`, 200, ensemble.Usage{}},
		{"empty", `{"content":[],"usage":{"input_tokens":11,"output_tokens":2}}`, 200, ensemble.Usage{Input: 11, Output: 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			_, a, _ := newAgent(t, func(w http.ResponseWriter, r *http.Request) {
				calls++
				if calls == 2 {
					w.WriteHeader(tc.status)
					fmt.Fprint(w, tc.body)
					return
				}
				fmt.Fprint(w, `{"content":[{"type":"text","text":"ok"}],"usage":{"input_tokens":3,"output_tokens":1}}`)
			})
			if _, err := a.Ask(context.Background(), "first"); err != nil {
				t.Fatal(err)
			}
			before := a.History()
			if _, err := a.Ask(context.Background(), "failed"); err == nil {
				t.Fatal("failure must return error")
			}
			if calls != 2 {
				t.Fatal("failure retried")
			}
			if !reflect.DeepEqual(a.History(), before) {
				t.Fatal("failed turn committed")
			}
			if a.Engine().Usage() != (ensemble.Usage{Input: 3 + tc.measured.Input, Output: 1 + tc.measured.Output}) {
				t.Fatal("measured failure usage")
			}
			if _, err := a.Ask(context.Background(), "   "); err == nil || calls != 2 {
				t.Fatal("empty user must fail before request")
			}
		})
	}
}

// TestCancellation checks cancellation before headers and during body reads.
func TestCancellation(t *testing.T) {
	for _, phase := range []string{"headers", "body"} {
		t.Run(phase, func(t *testing.T) {
			started := make(chan struct{})
			release := make(chan struct{})
			_, a, _ := newAgent(t, func(w http.ResponseWriter, r *http.Request) {
				if phase == "body" {
					fmt.Fprint(w, `{"content":`)
					w.(http.Flusher).Flush()
				}
				close(started)
				select {
				case <-r.Context().Done():
				case <-release:
				}
			})
			defer close(release)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { _, err := a.Ask(ctx, "wait"); done <- err }()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("request not started")
			}
			cancel()
			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("cancellation: %v", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("cancellation did not unblock request")
			}
			if len(a.History()) != 0 || a.Engine().Usage() != (ensemble.Usage{}) {
				t.Fatal("canceled turn changed state")
			}
		})
	}
}

// TestConfigValidation rejects incomplete configuration before registration.
func TestConfigValidation(t *testing.T) {
	root := ensemble.New(log.New(new(bytes.Buffer), "", 0))
	valid := ensemble.Config{BaseURL: "http://localhost", APIKey: "key", Model: "model", System: "system", MaxTokens: 1}
	for _, field := range []string{"url", "system", "tokens"} {
		bad := valid
		switch field {
		case "url":
			bad.BaseURL = "not-url"
		case "key":
			bad.APIKey = ""
		case "model":
			bad.Model = ""
		case "system":
			bad.System = ""
		case "tokens":
			bad.MaxTokens = 0
		}
		if _, err := root.NewAgent(bad); err == nil || len(root.Agents()) != 0 {
			t.Fatalf("invalid %s registered", field)
		}
	}
}

// TestMeasuredZero distinguishes known zero usage from absent counters.
func TestMeasuredZero(t *testing.T) {
	_, a, _ := newAgent(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"content":[{"type":"text","text":"zero"}],"usage":{"input_tokens":0,"output_tokens":0}}`)
	})
	if answer, err := a.Ask(context.Background(), "zero"); err != nil || answer != "zero" || len(a.History()) != 2 {
		t.Fatalf("measured zero rejected: %q %v", answer, err)
	}
}

// TestRequestLifetime crosses the real one-minute boundary before headers and
// after headers. These are intentionally slow, not a simulated HTTP client.
func TestRequestLifetime(t *testing.T) {
	for _, phase := range []string{"headers", "body"} {
		t.Run(phase, func(t *testing.T) {
			_, a, _ := newAgent(t, func(w http.ResponseWriter, r *http.Request) {
				if phase == "body" {
					w.WriteHeader(http.StatusOK)
					w.(http.Flusher).Flush()
				}
				select {
				case <-time.After(62 * time.Second):
				case <-r.Context().Done():
					return
				}
				fmt.Fprint(w, `{"content":[{"type":"text","text":"patient"}],"usage":{"input_tokens":1,"output_tokens":1}}`)
			})
			// A test watchdog is not a production response timeout.
			ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
			defer cancel()
			if answer, err := a.Ask(ctx, "wait"); err != nil || answer != "patient" {
				t.Fatalf("valid delayed %s response lost: %q %v", phase, answer, err)
			}
		})
	}
}

// TestUninitializedRoot protects the non-nil, service-bearing parent invariant.
func TestUninitializedRoot(t *testing.T) {
	for _, root := range []*ensemble.Ensemble{nil, {}} {
		child, err := root.NewAgent(ensemble.Config{BaseURL: "http://localhost", APIKey: "key", Model: "model", System: "system", MaxTokens: 1})
		if err == nil || child != nil {
			t.Fatal("uninitialized root returned an invalid owner chain")
		}
	}
}
