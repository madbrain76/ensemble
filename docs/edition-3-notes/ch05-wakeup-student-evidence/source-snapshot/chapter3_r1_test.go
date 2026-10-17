package ensemble_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"ensemble"
)

func r1Signed(text string, flag *bool) map[string]any {
	p := map[string]any{"text": text, "thoughtSignature": "r1-public-signature"}
	if flag != nil {
		p["thought"] = *flag
	}
	return p
}
func r1Wire(parts ...any) any {
	return map[string]any{"modelVersion": "test-model", "candidates": []any{map[string]any{"content": map[string]any{"parts": parts}}}, "usageMetadata": map[string]int{"promptTokenCount": 3, "candidatesTokenCount": 2, "thoughtsTokenCount": 1}}
}
func r1DecodeParts(t *testing.T, data []byte) []map[string]any {
	t.Helper()
	var wire struct {
		Contents []struct {
			Role  string           `json:"role"`
			Parts []map[string]any `json:"parts"`
		} `json:"contents"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	parts := []map[string]any{}
	for _, c := range wire.Contents {
		if c.Role == "model" {
			parts = append(parts, c.Parts...)
		}
	}
	return parts
}
func r1CheckReplay(t *testing.T, data []byte, signed bool) {
	t.Helper()
	parts := r1DecodeParts(t, data)
	counts := map[string]int{}
	for _, p := range parts {
		if text, ok := p["text"].(string); ok {
			counts[text]++
			if text == "r1 visible explanation" || text == "final answer" {
				if signed && p["thoughtSignature"] != "r1-public-signature" {
					t.Error("same-source signature not attached to its ordinary text part")
				}
				if !signed && p["thoughtSignature"] != nil {
					t.Error("cross-model signature leaked")
				}
				if p["thought"] == true || p["functionCall"] != nil {
					t.Error("malformed ordinary text replay part")
				}
			}
		}
	}
	for _, text := range []string{"r1 visible explanation", "final answer"} {
		if counts[text] != 1 {
			t.Errorf("ordinary text %q count=%d, want exactly one", text, counts[text])
		}
	}
	if signed && counts["r1-hidden-only"] != 1 {
		t.Error("same-source hidden part not preserved")
	}
	if !signed && (bytes.Contains(data, []byte("r1-hidden-only")) || bytes.Contains(data, []byte("r1-public-signature")) || bytes.Contains(data, []byte("r1-hidden-signature"))) {
		t.Error("model-bound material leaked")
	}
}

func TestR1SignedTextPublicAPI(t *testing.T) {
	no := false
	for _, fixture := range []struct {
		name string
		flag *bool
	}{{"absent", nil}, {"false", &no}} {
		t.Run(fixture.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "marker")
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, _ := io.ReadAll(r.Body)
				requests++
				switch requests {
				case 1:
					_ = json.NewEncoder(w).Encode(r1Wire(map[string]any{"text": "r1-hidden-only", "thought": true, "thoughtSignature": "r1-hidden-signature"}, r1Signed("r1 visible explanation", fixture.flag), map[string]any{"functionCall": map[string]any{"name": "write_file", "args": map[string]string{"path": path, "content": "executed"}}}))
				case 2:
					content, err := os.ReadFile(path)
					if err != nil || string(content) != "executed" {
						t.Error("signed visible explanation was not executed")
					}
					if !bytes.Contains(data, []byte("file written")) || bytes.Contains(data, []byte("tool batch refused")) {
						t.Error("missing successful real tool feedback")
					}
					parts := r1DecodeParts(t, data)
					if len(parts) != 3 || parts[1]["text"] != "r1 visible explanation" || parts[1]["thoughtSignature"] != "r1-public-signature" {
						t.Error("explaining signed text lost, duplicated, or reordered in retry")
					}
					_ = json.NewEncoder(w).Encode(r1Wire(r1Signed("final answer", fixture.flag)))
				default:
					t.Error("unexpected retry")
					http.Error(w, "unexpected", 500)
				}
			}))
			defer server.Close()
			a := toolAgent(t, ensemble.Gemini, true, true, server.URL)
			answer, err := a.Ask(context.Background(), "fixture request")
			if err != nil || answer != "final answer" {
				t.Errorf("ordinary signed final = %q, err=%v", answer, err)
			}
			if requests != 2 {
				t.Errorf("requests=%d", requests)
			}
			if a.Engine().Usage() != (ensemble.Usage{Input: 6, Output: 6}) {
				t.Error("usage changed")
			}
			calls, returns := []string{}, []string{}
			for _, e := range a.Log().Events() {
				if e.Type == "response_ended" {
					for _, p := range e.Response.Parts {
						if p.Text == "r1 visible explanation" || p.Text == "final answer" {
							if p.Type != "text" || p.From != (ensemble.Provenance{Vendor: ensemble.Gemini, Model: "test-model", Surface: ensemble.GenerateContent}) || string(p.Opaque) != `"r1-public-signature"` {
								t.Error("signed text normalization lost public text or exact provenance")
							}
						}
					}
				}
				if e.Type == "tool_called" {
					calls = append(calls, e.Tool.CallID)
				}
				if e.Type == "tool_returned" {
					returns = append(returns, e.Tool.CallID)
					if e.Tool.IsError {
						t.Error("signed ordinary text falsely refused")
					}
				}
			}
			if len(calls) != 1 || calls[0] == "" || !reflect.DeepEqual(calls, returns) {
				t.Error("tool result pairing changed")
			}
			var dump bytes.Buffer
			if err := a.Log().Dump(&dump); err != nil {
				t.Fatal(err)
			}
			loaded := offlineAgent(t)
			if err := loaded.Log().Load(bytes.NewReader(dump.Bytes())); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(a.Log().Context(), loaded.Log().Context()) {
				t.Error("dump/load projection changed signed text")
			}
			before := append([]byte(nil), dump.Bytes()...)
			source := ensemble.Provenance{Vendor: ensemble.Gemini, Model: "test-model", Surface: ensemble.GenerateContent}
			for _, owner := range []ensemble.Agent{a, loaded} {
				data, err := owner.Engine().Render(owner.Log().Context(), source)
				if err != nil {
					t.Fatal(err)
				}
				r1CheckReplay(t, data, true)
				target := source
				target.Model = "other-model"
				data, err = owner.Engine().Render(owner.Log().Context(), target)
				if err != nil {
					t.Fatal(err)
				}
				r1CheckReplay(t, data, false)
				for _, vendor := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI} {
					target := issuer(vendor)
					target.Model = source.Model
					data, err = owner.Engine().Render(owner.Log().Context(), target)
					if err != nil {
						t.Fatal(err)
					}
					for _, word := range []string{"r1 visible explanation", "final answer"} {
						if bytes.Count(data, []byte(word)) != 1 {
							t.Errorf("cross-surface ordinary text %q lost or duplicated", word)
						}
					}
					for _, secret := range []string{"r1-hidden-only", "r1-public-signature", "r1-hidden-signature"} {
						if bytes.Contains(data, []byte(secret)) {
							t.Error("cross-surface private material leaked")
						}
					}
				}
			}
			// Render is public and must check all provenance dimensions even for a
			// detached, manually prepared projection, not rely solely on parser input.
			for _, alter := range []string{"vendor", "surface"} {
				c := loaded.Log().Context()
				for i := range c.Dialogue {
					for j := range c.Dialogue[i].Parts {
						p := &c.Dialogue[i].Parts[j]
						if p.Type == "text" && len(p.Opaque) > 0 {
							if alter == "vendor" {
								p.From.Vendor = ensemble.Anthropic
							} else {
								p.From.Surface = ensemble.Messages
							}
						}
					}
				}
				data, err := loaded.Engine().Render(c, source)
				if err != nil {
					t.Fatal(err)
				}
				if bytes.Contains(data, []byte("r1-public-signature")) {
					t.Errorf("%s equality weakened", alter)
				}
				if bytes.Count(data, []byte("final answer")) != 1 {
					t.Error("foreign provenance dropped ordinary text")
				}
			}
			dump.Reset()
			_ = a.Log().Dump(&dump)
			if !bytes.Equal(before, dump.Bytes()) {
				t.Error("render mutated History")
			}
		})
	}
}

func TestR1TrueThoughtCannotExplain(t *testing.T) {
	path := filepath.Join(t.TempDir(), "must-not-exist")
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			_ = json.NewEncoder(w).Encode(r1Wire(map[string]any{"text": "r1-hidden-only", "thought": true, "thoughtSignature": "r1-hidden-signature"}, map[string]any{"functionCall": map[string]any{"name": "write_file", "args": map[string]string{"path": path, "content": "bad"}}}))
		} else {
			data, _ := io.ReadAll(r.Body)
			if !bytes.Contains(data, []byte("tool batch refused")) {
				t.Error("hidden thought qualified as explanation")
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Error("hidden-only batch executed")
			}
			_ = json.NewEncoder(w).Encode(r1Wire(map[string]any{"text": "refusal acknowledged"}))
		}
	}))
	defer server.Close()
	a := toolAgent(t, ensemble.Gemini, true, true, server.URL)
	answer, err := a.Ask(context.Background(), "fixture")
	if err != nil || answer != "refusal acknowledged" || requests != 2 {
		t.Fatalf("answer=%q err=%v requests=%d", answer, err, requests)
	}
}

func TestR1FallbackIDsUniqueAcrossBatches(t *testing.T) {
	dir := t.TempDir()
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests <= 2 {
			blocks := []any{map[string]any{"text": "I will inspect the disposable directory."}}
			for i := 0; i < 2; i++ {
				blocks = append(blocks, map[string]any{"functionCall": map[string]any{"name": "list_directory", "args": map[string]string{"path": dir}}})
			}
			_ = json.NewEncoder(w).Encode(r1Wire(blocks...))
		} else {
			_ = json.NewEncoder(w).Encode(r1Wire(map[string]any{"text": "done"}))
		}
	}))
	defer server.Close()
	a := toolAgent(t, ensemble.Gemini, true, true, server.URL)
	if answer, err := a.Ask(context.Background(), "inspect twice"); err != nil || answer != "done" {
		t.Fatal(answer, err)
	}
	calls, returned := []string{}, []string{}
	seen := map[string]bool{}
	for _, e := range a.Log().Events() {
		if e.Type == "tool_called" {
			id := e.Tool.CallID
			if id == "" || seen[id] {
				t.Error("empty or colliding fallback execution ID")
			}
			seen[id] = true
			calls = append(calls, id)
		}
		if e.Type == "tool_returned" {
			returned = append(returned, e.Tool.CallID)
			if e.Tool.IsError {
				t.Error("unexpected refusal")
			}
		}
	}
	if len(calls) != 4 || !reflect.DeepEqual(calls, returned) {
		t.Error("fallback calls/results not paired in order")
	}
	var dump bytes.Buffer
	_ = a.Log().Dump(&dump)
	b := offlineAgent(t)
	if err := b.Log().Load(&dump); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a.Log().Context(), b.Log().Context()) {
		t.Error("fallback pairing lost on replay")
	}
}

func TestR1ActualCLISignedText(t *testing.T) {
	exe := filepath.Join(t.TempDir(), "student-cli")
	if out, err := exec.Command("go", "build", "-o", exe, "./cmd/ch01").CombinedOutput(); err != nil {
		t.Fatalf("build: %v %s", err, out)
	}
	binary, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("fixture CLI SHA256=%x", sha256.Sum256(binary))
	no := false
	for _, fixture := range []struct {
		name string
		flag *bool
	}{{"absent", nil}, {"false", &no}} {
		t.Run(fixture.name, func(t *testing.T) {
			dir := t.TempDir()
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, _ := io.ReadAll(r.Body)
				requests++
				switch requests {
				case 1:
					_ = json.NewEncoder(w).Encode(r1Wire(r1Signed("r1 visible explanation", fixture.flag), map[string]any{"functionCall": map[string]any{"name": "write_file", "args": map[string]string{"path": "marker", "content": "actual CLI effect"}}}))
				case 2:
					content, err := os.ReadFile(filepath.Join(dir, "marker"))
					if err != nil || string(content) != "actual CLI effect" {
						t.Error("CLI falsely refused signed explanation")
					}
					if !bytes.Contains(data, []byte("file written")) {
						t.Error("CLI lost real tool feedback")
					}
					_ = json.NewEncoder(w).Encode(r1Wire(r1Signed("first final answer", fixture.flag)))
				case 3:
					parts := r1DecodeParts(t, data)
					found := 0
					for _, p := range parts {
						if p["text"] == "first final answer" {
							found++
							if p["thoughtSignature"] != "r1-public-signature" {
								t.Error("CLI followup lost exact-source signature")
							}
						}
					}
					if found != 1 {
						t.Error("CLI followup lost or duplicated earlier ordinary text")
					}
					_ = json.NewEncoder(w).Encode(r1Wire(r1Signed("second final answer", fixture.flag)))
				default:
					t.Error("extra CLI request")
					http.Error(w, "unexpected", 500)
				}
			}))
			defer server.Close()
			env := []string{}
			for _, v := range os.Environ() {
				name := strings.SplitN(v, "=", 2)[0]
				if strings.HasPrefix(name, "LLM_") || strings.HasPrefix(name, "ANTHROPIC_") || strings.HasPrefix(name, "OPENAI_") || strings.HasPrefix(name, "GEMINI_") || name == "CH02_LOG" {
					continue
				}
				env = append(env, v)
			}
			env = append(env, "LLM_VENDOR=gemini", "LLM_MODEL=test-model", "LLM_API_KEY=benign-local-fixture", "LLM_BASE_URL="+server.URL, "LLM_REQUIRE_TOOL_EXPLANATION=true", "CH02_LOG="+filepath.Join(dir, "history.jsonl"))
			cmd := exec.Command(exe, "chat")
			cmd.Dir = dir
			cmd.Env = env
			cmd.Stdin = strings.NewReader("repair fixture\nrecall fixture\n")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("actual CLI: %v %s", err, &stderr)
			}
			for _, answer := range []string{"first final answer", "second final answer"} {
				if strings.Count(stdout.String(), answer) != 1 {
					t.Errorf("actual chat answer %q missing or duplicated", answer)
				}
			}
			if strings.Contains(stdout.String(), "r1-public-signature") {
				t.Error("CLI displayed signature")
			}
			if requests != 3 {
				t.Errorf("CLI requests=%d", requests)
			}
			logdata, err := os.ReadFile(filepath.Join(dir, "history.jsonl"))
			if err != nil {
				t.Fatal(err)
			}
			a := offlineAgent(t)
			if err := a.Log().Load(bytes.NewReader(logdata)); err != nil {
				t.Fatal(err)
			}
			data, err := a.Engine().Render(a.Log().Context(), ensemble.Provenance{Vendor: ensemble.Gemini, Model: "test-model", Surface: ensemble.GenerateContent})
			if err != nil {
				t.Fatal(err)
			}
			for _, answer := range []string{"first final answer", "second final answer"} {
				if bytes.Count(data, []byte(answer)) != 1 {
					t.Error("CLI saved/reloaded ordinary text lost or duplicated")
				}
			}
		})
	}
}
