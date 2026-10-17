package ensemble_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLI builds the shipped frontend and observes actual stdout/stderr streams.
func TestCLI(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "ch01")
	if output, err := exec.Command("go", "build", "-o", binary, "./cmd/ch01").CombinedOutput(); err != nil {
		t.Fatalf("build: %v %s", err, output)
	}
	for _, mode := range []string{"protocol", "chat", "error", "empty"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				var request struct {
					Messages []ensembleMessage `json:"messages"`
				}
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Error(err)
				}
				if mode == "error" {
					w.WriteHeader(429)
					fmt.Fprint(w, "private-provider-error-marker")
					return
				}
				if len(request.Messages) != calls*2-1 {
					t.Error("CLI history lost")
				}
				if mode == "protocol" && calls == 1 && len(request.Messages[0].Content) != 70000 {
					t.Error("long input truncated")
				}
				fmt.Fprint(w, `{"content":[{"type":"text","text":"one "},{"type":"text","text":"two"}],"usage":{"input_tokens":4,"output_tokens":3}}`)
			}))
			defer server.Close()
			cmd := exec.Command(binary)
			input := "{\"user\":\"" + strings.Repeat("x", 70000) + "\"}\n{\"user\":\"second\"}" // no final newline
			if mode == "chat" {
				cmd = exec.Command(binary, "chat")
				input = "\nhello\nsecond\n"
			}
			if mode == "empty" {
				input = ""
			}
			// Replace provider variables instead of inheriting an accidental live endpoint.
			for _, value := range os.Environ() {
				if !strings.HasPrefix(value, "ANTHROPIC_") {
					cmd.Env = append(cmd.Env, value)
				}
			}
			cmd.Env = append(cmd.Env, "ANTHROPIC_BASE_URL="+server.URL, "ANTHROPIC_API_KEY=local-only", "ANTHROPIC_MODEL=test-model")
			cmd.Env = append(cmd.Env, "CH02_LOG="+filepath.Join(t.TempDir(), "session.jsonl"))
			cmd.Stdin = strings.NewReader(input)
			var out, diagnostics bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &diagnostics
			err := cmd.Run()
			if mode == "error" {
				if err == nil || calls != 1 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "HTTP 429") || strings.Contains(diagnostics.String(), "private-provider-error-marker") || strings.Contains(diagnostics.String(), "local-only") {
					t.Fatalf("CLI failure boundary: err=%v calls=%d stdout=%q stderr=%q", err, calls, out.String(), diagnostics.String())
				}
				return
			}
			if err != nil {
				t.Fatalf("run: %v %s", err, diagnostics.String())
			}
			switch mode {
			case "protocol":
				if out.String() != "{\"assistant\":\"one two\"}\n{\"assistant\":\"one two\"}\n{\"usage\":{\"input\":8,\"cache_write\":0,\"cache_read\":0,\"output\":6}}\n" || calls != 2 {
					t.Fatalf("protocol output/calls: %q %d", out.String(), calls)
				}
			case "empty":
				if out.String() != "{\"usage\":{\"input\":0,\"cache_write\":0,\"cache_read\":0,\"output\":0}}\n" || calls != 0 {
					t.Fatalf("empty protocol: %q", out.String())
				}
			case "chat":
				if calls != 2 || strings.Count(out.String(), "one two") != 2 || !strings.Contains(out.String(), "> ") || !strings.Contains(diagnostics.String(), "cumulative usage: input=8 output=6") {
					t.Fatalf("chat output: %q %q", out.String(), diagnostics.String())
				}
			}
		})
	}
}

// ensembleMessage decodes the CLI's real wire without coupling the fixture to
// the production Message encoder.
type ensembleMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// UnmarshalJSON accepts both Chapter 1 strings and Chapter 2 text-block arrays
// at the real external boundary, while retaining the same exact text assertions.
func (m *ensembleMessage) UnmarshalJSON(data []byte) error {
	var wire struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	m.Role = wire.Role
	if err := json.Unmarshal(wire.Content, &m.Content); err == nil {
		return nil
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(wire.Content, &blocks); err != nil {
		return err
	}
	for _, block := range blocks {
		if block.Type == "text" {
			m.Content += block.Text
		}
	}
	return nil
}
