package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These tests call the same run function as main; it constructs the real public
// Ensemble and Agent. Only the external HTTP server is a protocol fixture.
func cliSettings(t *testing.T) string {
	t.Helper()
	for _, name := range []string{"LLM_REQUIRE_TOOL_EXPLANATION", "LLM_VENDOR", "LLM_MODEL", "LLM_RESOLVED_MODEL", "LLM_BASE_URL", "LLM_API_KEY", "LLM_SYSTEM", "LLM_MAX_OUTPUT_TOKENS", "ANTHROPIC_MODEL", "ANTHROPIC_BASE_URL", "ANTHROPIC_API_KEY"} {
		t.Setenv(name, "")
	}
	path := filepath.Join(t.TempDir(), "log.jsonl")
	t.Setenv("CH02_LOG", path)
	return path
}
func TestRealCommandSession(t *testing.T) {
	path := cliSettings(t)
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/v1/messages" || r.Header.Get("x-api-key") != "fixture-key" {
			t.Error("wrong live boundary")
		}
		fmt.Fprint(w, `{"model":"resolved-model","content":[{"type":"text","text":"answer"}],"usage":{"input_tokens":2,"cache_creation_input_tokens":3,"cache_read_input_tokens":4,"output_tokens":5}}`)
	}))
	defer server.Close()
	t.Setenv("LLM_MODEL", "model")
	t.Setenv("LLM_API_KEY", "fixture-key")
	t.Setenv("LLM_BASE_URL", server.URL)
	var stdout bytes.Buffer
	if err := execute(t, nil, strings.NewReader("{\"ephemeral\":\"one request\"}\n{\"user\":\"hello\"}\n"), &stdout); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || !strings.Contains(stdout.String(), `"cache_write":3`) {
		t.Fatalf("calls %d stdout %s", calls, &stdout)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("fixture-key")) || !bytes.Contains(data, []byte("response_ended")) {
		t.Fatal("bad saved capture")
	}
	stdout.Reset()
	if err := execute(t, []string{"dump"}, nil, &stdout); err != nil || !bytes.Equal(data, stdout.Bytes()) {
		t.Fatalf("dump: %v", err)
	}
	for _, vendor := range []string{"anthropic", "openai", "gemini"} {
		t.Setenv("LLM_VENDOR", vendor)
		stdout.Reset()
		if err := execute(t, []string{"render", path}, nil, &stdout); err != nil || !json.Valid(stdout.Bytes()) {
			t.Fatalf("render %s: %v %s", vendor, err, &stdout)
		}
	}
	if calls != 1 {
		t.Fatal("offline render made a live call")
	}
}
func TestCommandRefusalsAndFailureCapture(t *testing.T) {
	path := cliSettings(t)
	for _, args := range [][]string{{"render", "extra"}, {"--protocol"}, {"dump-log"}} {
		if err := execute(t, args, strings.NewReader(""), &bytes.Buffer{}); err == nil {
			t.Fatalf("accepted %v", args)
		}
	}
	for _, input := range []string{"not json\n", "{}\n", "{\"user\":\"a\",\"ephemeral\":\"b\"}\n", "{\"other\":1}\n", "{\"user\":\"a\"} {}\n", "{\"user\":\"needs key\"}\n"} {
		if err := execute(t, nil, strings.NewReader(input), &bytes.Buffer{}); err == nil {
			t.Fatalf("accepted %s", input)
		}
		data, err := os.ReadFile(path)
		if err != nil || !bytes.HasPrefix(data, []byte(`{"log_version":1}`)) {
			t.Fatalf("failure not saved %v", err)
		}
	}
	data, _ := os.ReadFile(path)
	if !bytes.Contains(data, []byte("error_occurred")) {
		t.Fatal("send failure not captured")
	}
	t.Setenv("LLM_VENDOR", "future-provider")
	if err := execute(t, nil, strings.NewReader(""), &bytes.Buffer{}); err == nil {
		t.Fatal("unknown vendor")
	}
	t.Setenv("LLM_VENDOR", "")
	t.Setenv("CH02_LOG", t.TempDir())
	if err := execute(t, nil, strings.NewReader(""), &bytes.Buffer{}); err == nil {
		t.Fatal("save to directory")
	}
	if err := os.WriteFile(path, []byte(`{"log_version":99}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := execute(t, []string{"render", path}, nil, &bytes.Buffer{}); err == nil {
		t.Fatal("unknown log version")
	}
}
func TestHumanCommand(t *testing.T) {
	cliSettings(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"content":[{"type":"text","text":"human answer"}],"usage":{"input_tokens":1,"output_tokens":2}}`)
	}))
	defer server.Close()
	t.Setenv("ANTHROPIC_MODEL", "model")
	t.Setenv("ANTHROPIC_API_KEY", "fixture-key")
	t.Setenv("ANTHROPIC_BASE_URL", server.URL)
	var stdout bytes.Buffer
	if err := execute(t, []string{"chat"}, strings.NewReader("\nhello\n"), &stdout); err != nil || !strings.Contains(stdout.String(), "human answer") {
		t.Fatalf("human %v %s", err, &stdout)
	}
	if err := execute(t, []string{"chat"}, strings.NewReader(""), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANTHROPIC_API_KEY", "")
	if err := execute(t, []string{"chat"}, strings.NewReader("hello\n"), &bytes.Buffer{}); err == nil {
		t.Fatal("human missing credentials")
	}
}

// execute supplies actual process files to the unchanged command path. It does
// not replace application construction, History, Engine, or any owner service.
func execute(t *testing.T, args []string, in io.Reader, out io.Writer) error {
	t.Helper()
	dir := t.TempDir()
	input, err := os.Create(filepath.Join(dir, "stdin"))
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	if in != nil {
		if _, err := io.Copy(input, in); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := input.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	output, err := os.Create(filepath.Join(dir, "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()
	oldIn, oldOut := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = input, output
	defer func() { os.Stdin, os.Stdout = oldIn, oldOut }()
	result := run(args)
	if _, err := output.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(out, output); err != nil {
		t.Fatal(err)
	}
	return result
}

// Public CLI selection is explicit and reversible; it does not infer model names.
func TestChapter3CLIPolicySelection(t *testing.T) {
	for _, value := range []string{"true", "false", ""} {
		t.Run("policy="+value, func(t *testing.T) {
			cliSettings(t)
			dir := t.TempDir()
			t.Chdir(dir)
			marker := filepath.Join(dir, "marker")
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if calls == 1 {
					call := map[string]any{"type": "tool_use", "id": "cli-call", "name": "write_file", "input": map[string]string{"path": marker, "content": "cli"}}
					_ = json.NewEncoder(w).Encode(map[string]any{"content": []any{call}, "usage": map[string]int{"input_tokens": 1, "output_tokens": 1}})
				} else {
					data, _ := io.ReadAll(r.Body)
					refused := bytes.Contains(data, []byte("tool batch refused"))
					if refused != (value == "true") {
						t.Errorf("policy=%q refusal=%t", value, refused)
					}
					fmt.Fprint(w, `{"content":[{"type":"text","text":"done"}],"usage":{"input_tokens":1,"output_tokens":1}}`)
				}
			}))
			defer server.Close()
			t.Setenv("LLM_REQUIRE_TOOL_EXPLANATION", value)
			t.Setenv("LLM_MODEL", "ordinary-public-name")
			t.Setenv("LLM_API_KEY", "fixture")
			t.Setenv("LLM_BASE_URL", server.URL)
			var output bytes.Buffer
			if err := execute(t, nil, strings.NewReader("{\"user\":\"act\"}\n"), &output); err != nil {
				t.Fatal(err)
			}
			_, err := os.Stat(marker)
			if (err == nil) != (value != "true") {
				t.Fatalf("policy %q stat=%v", value, err)
			}
			if calls != 2 || !strings.Contains(output.String(), `"assistant":"done"`) {
				t.Fatalf("calls %d output %s", calls, &output)
			}
		})
	}
	cliSettings(t)
	t.Setenv("LLM_REQUIRE_TOOL_EXPLANATION", "not-a-bool")
	if _, err := settings(); err == nil {
		t.Fatal("invalid policy accepted")
	}
}
