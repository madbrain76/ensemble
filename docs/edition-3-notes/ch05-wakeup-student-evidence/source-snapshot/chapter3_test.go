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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ensemble"
)

func toolAgent(t *testing.T, vendor ensemble.Vendor, enabled, policy bool, url string) ensemble.Agent {
	t.Helper()
	a, err := ensemble.New(log.New(io.Discard, "", 0)).NewAgent(ensemble.Config{Vendor: vendor, Model: "test-model", APIKey: "benign-fixture", BaseURL: url, System: "Use tools with an explanation.", MaxTokens: 1024, EnableTools: enabled, RequireToolExplanation: policy})
	if err != nil {
		t.Fatal(err)
	}
	return a
}
func callPart(id, name string, args any) ensemble.Part {
	data, _ := json.Marshal(args)
	return ensemble.Part{Type: "tool_call", CallID: id, Name: name, Args: data}
}
func prose(text string) ensemble.Part { return ensemble.Part{Type: "text", Text: text} }

// providerReply exercises actual wire parsing, HTTP and dispatch, not a substitute loop.
func providerReply(v ensemble.Vendor, parts []ensemble.Part) any {
	switch v {
	case ensemble.Anthropic:
		blocks := []any{}
		for _, p := range parts {
			switch p.Type {
			case "text":
				blocks = append(blocks, map[string]any{"type": "text", "text": p.Text})
			case "opaque":
				blocks = append(blocks, map[string]any{"type": "thinking", "thinking": "not visible", "signature": "fixture-signature"})
			case "tool_call":
				blocks = append(blocks, map[string]any{"type": "tool_use", "id": p.CallID, "name": p.Name, "input": p.Args})
			}
		}
		return map[string]any{"model": "test-model", "content": blocks, "usage": map[string]int{"input_tokens": 3, "output_tokens": 2}}
	case ensemble.OpenAI:
		msg := map[string]any{"content": nil}
		calls := []any{}
		for _, p := range parts {
			switch p.Type {
			case "text":
				msg["content"] = p.Text
			case "opaque":
				msg["reasoning_content"] = "not visible"
			case "tool_call":
				calls = append(calls, map[string]any{"id": p.CallID, "type": "function", "function": map[string]any{"name": p.Name, "arguments": string(p.Args)}})
			}
		}
		if len(calls) > 0 {
			msg["tool_calls"] = calls
		}
		return map[string]any{"model": "test-model", "choices": []any{map[string]any{"message": msg}}, "usage": map[string]int{"prompt_tokens": 3, "completion_tokens": 2}}
	case ensemble.Gemini:
		blocks := []any{}
		for _, p := range parts {
			switch p.Type {
			case "text":
				blocks = append(blocks, map[string]any{"text": p.Text})
			case "opaque":
				blocks = append(blocks, map[string]any{"text": "not visible", "thought": true})
			case "tool_call":
				blocks = append(blocks, map[string]any{"functionCall": map[string]any{"id": p.CallID, "name": p.Name, "args": p.Args}})
			}
		}
		return map[string]any{"modelVersion": "test-model", "candidates": []any{map[string]any{"content": map[string]any{"parts": blocks}}}, "usageMetadata": map[string]int{"promptTokenCount": 3, "candidatesTokenCount": 2}}
	}
	panic("unsupported fixture vendor")
}

func TestChapter3SelectedPolicyWholeBatchAndRetry(t *testing.T) {
	for _, vendor := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI, ensemble.Gemini} {
		t.Run(vendor.String(), func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			write := filepath.Join(dir, "written")
			command := filepath.Join(dir, "command")
			edited := filepath.Join(dir, "existing")
			if err := os.WriteFile(edited, []byte("before"), 0600); err != nil {
				t.Fatal(err)
			}
			batch := []ensemble.Part{callPart("write", "write_file", map[string]string{"path": write, "content": "actual"}), callPart("command", "run_command", map[string]string{"command": "printf shell > '" + command + "'"}), callPart("edit", "edit_file", map[string]string{"path": edited, "old_text": "before", "new_text": "after"})}
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, _ := io.ReadAll(r.Body)
				requests++
				switch requests {
				case 1:
					_ = json.NewEncoder(w).Encode(providerReply(vendor, append([]ensemble.Part{{Type: "opaque"}}, batch...)))
				case 2:
					for _, path := range []string{write, command} {
						if _, err := os.Stat(path); !os.IsNotExist(err) {
							t.Errorf("unexplained batch changed %s: %v", path, err)
						}
					}
					text, _ := os.ReadFile(edited)
					if string(text) != "before" {
						t.Error("unexplained batch edited disk")
					}
					if bytes.Count(data, []byte("tool batch refused")) != 3 {
						t.Errorf("missing paired refusal feedback: %s", data)
					}
					_ = json.NewEncoder(w).Encode(providerReply(vendor, append([]ensemble.Part{prose("I will write, run the shell marker and edit the file to perform the requested repair.")}, batch...)))
				case 3:
					for _, path := range []string{write, command} {
						if _, err := os.Stat(path); err != nil {
							t.Errorf("explained retry did not execute %s: %v", path, err)
						}
					}
					text, _ := os.ReadFile(edited)
					if string(text) != "after" {
						t.Error("explained retry did not edit")
					}
					if !bytes.Contains(data, []byte("file written")) || !bytes.Contains(data, []byte("exit_code")) {
						t.Errorf("missing real retry feedback: %s", data)
					}
					_ = json.NewEncoder(w).Encode(providerReply(vendor, []ensemble.Part{prose("repaired")}))
				default:
					t.Error("extra request")
					http.Error(w, "unexpected", 500)
				}
			}))
			defer server.Close()
			a := toolAgent(t, vendor, true, true, server.URL)
			answer, err := a.Ask(context.Background(), "repair")
			if err != nil || answer != "repaired" {
				t.Fatalf("%q %v", answer, err)
			}
			if requests != 3 {
				t.Fatal(requests)
			}
			events := a.Log().Events()
			called, returned, responses := 0, 0, 0
			for i, e := range events {
				switch e.Type {
				case "response_ended":
					responses++
					if responses == 1 && len(e.Response.Parts) != 4 {
						t.Fatal("assistant response not preserved")
					}
					if responses == 1 {
						opaque := false
						for _, part := range e.Response.Parts {
							if part.Type == "opaque" {
								opaque = true
							}
						}
						if !opaque {
							t.Fatal("opaque response fact lost")
						}
					}
				case "tool_called":
					called++
					// Actor job completion reports may arrive between call and observation.
					next := i + 1
					for next < len(events) && events[next].Type == "job_report" {
						next++
					}
					if next >= len(events) || events[next].Type != "tool_returned" || events[next].Tool.CallID != e.Tool.CallID {
						t.Fatal("unpaired event ordering")
					}
				case "tool_returned":
					returned++
					if e.Tool.IsError != (returned <= 3) {
						t.Fatal("wrong refusal/success flags")
					}
				}
			}
			if called != 6 || returned != 6 || responses != 3 {
				t.Fatalf("events calls=%d results=%d responses=%d", called, returned, responses)
			}
			if a.Log().Context().Turn != ensemble.Idle || len(a.History()) != 2 {
				t.Fatal("followup state/compatibility history")
			}
			if u := a.Engine().Usage(); u.Input != 9 || u.Output != 6 {
				t.Fatal(u)
			}
		})
	}
}

func TestChapter3PolicySelectionAndPrecedingText(t *testing.T) {
	cases := []struct {
		name     string
		policy   bool
		position string
		executes bool
	}{
		{"disabled-call-only", false, "none", true}, {"selected-call-only", true, "none", false},
		{"whitespace", true, "blank", false}, {"opaque-only", true, "opaque", false},
		{"late-text", true, "late", false}, {"preceding-text", true, "before", true},
		{"prior-response-not-current", true, "prior", false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "marker")
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				parts := []ensemble.Part{prose("done")}
				if test.position == "prior" && calls == 1 {
					parts = []ensemble.Part{prose("I will act on the next turn.")}
				} else if calls == 1 || (test.position == "prior" && calls == 2) {
					p := callPart("mark", "write_file", map[string]string{"path": path, "content": "ran"})
					parts = []ensemble.Part{p}
					switch test.position {
					case "before":
						parts = []ensemble.Part{prose("I will write a marker to record completion."), p}
					case "blank":
						parts = []ensemble.Part{prose(" \n\t"), p}
					case "opaque":
						parts = []ensemble.Part{{Type: "opaque"}, p}
					case "late":
						parts = []ensemble.Part{p, prose("Too late.")}
					}
				}
				_ = json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, parts))
			}))
			defer server.Close()
			a := toolAgent(t, ensemble.Anthropic, true, test.policy, server.URL)
			if test.position == "prior" {
				if _, err := a.Ask(context.Background(), "first"); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := a.Ask(context.Background(), "act"); err != nil {
				t.Fatal(err)
			}
			_, err := os.Stat(path)
			if (err == nil) != test.executes {
				t.Fatalf("executes=%t stat=%v", test.executes, err)
			}
		})
	}
}

func TestChapter3RoundBudgetAndErrors(t *testing.T) {
	for _, policy := range []bool{true, false} {
		t.Run(fmt.Sprint(policy), func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				_ = json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, []ensemble.Part{callPart(fmt.Sprint(requests), "missing", map[string]string{})}))
			}))
			defer server.Close()
			a := toolAgent(t, ensemble.Anthropic, true, policy, server.URL)
			_, err := a.Ask(context.Background(), "keep going")
			if err == nil || !strings.Contains(err.Error(), "tool batch limit") || requests != ensemble.DefaultToolBatchLimit+1 {
				t.Fatalf("requests=%d err=%v", requests, err)
			}
			returned := 0
			for _, e := range a.Log().Events() {
				if e.Type == "tool_returned" {
					returned++
					if !e.Tool.IsError {
						t.Fatal("not error")
					}
				}
			}
			if returned != ensemble.DefaultToolBatchLimit+1 {
				t.Fatal(returned)
			}
		})
	}
}

func TestChapter3RegistryAndSchemas(t *testing.T) {
	for _, vendor := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI, ensemble.Gemini} {
		for _, enabled := range []bool{false, true} {
			a := toolAgent(t, vendor, enabled, false, "")
			wire, err := a.Engine().Render(a.Log().Context(), a.Engine().Target())
			if err != nil {
				t.Fatal(err)
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(wire, &body); err != nil {
				t.Fatal(err)
			}
			if _, has := body["tools"]; has != enabled {
				t.Fatalf("tools enabled=%v %s", enabled, wire)
			}
			if vendor == ensemble.Gemini && bytes.Contains(wire, []byte("additionalProperties")) {
				t.Fatal("Gemini Schema rejects additionalProperties; local registry must still be strict")
			}
			if !enabled {
				if len(a.Tools()) != 0 {
					t.Fatal("ordinary agent has tools")
				}
				continue
			}
			if len(a.Tools()) != 10 {
				t.Fatal(len(a.Tools()))
			}
			for _, tool := range a.Tools() {
				if tool.Parent() != a {
					t.Fatal("broken tool owner")
				}
				d := tool.Definition()
				if !bytes.Contains(wire, []byte(d.Name)) || !json.Valid(d.Parameters) {
					t.Fatal(d)
				}
				var schema map[string]any
				_ = json.Unmarshal(d.Parameters, &schema)
				if schema["type"] != "object" || schema["additionalProperties"] != false {
					t.Fatal(schema)
				}
				d.Parameters[0] = '!'
				if !json.Valid(tool.Definition().Parameters) {
					t.Fatal("mutable definition")
				}
			}
			snapshot := a.Tools()
			snapshot[0] = nil
			if a.Tools()[0] == nil {
				t.Fatal("mutable registry")
			}
		}
	}
}

func invoke(t *testing.T, a ensemble.Agent, name string, args any) (string, error) {
	t.Helper()
	raw, _ := json.Marshal(args)
	for _, tool := range a.Tools() {
		if tool.Definition().Name == name {
			return tool.Execute(context.Background(), raw)
		}
	}
	t.Fatalf("missing %s", name)
	return "", nil
}

func TestChapter3FileTools(t *testing.T) {
	a := toolAgent(t, ensemble.Anthropic, true, false, "")
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "f.txt")
	if _, err := invoke(t, a, "write_file", map[string]string{"path": path, "content": "alpha\nbeta\ngamma"}); err != nil {
		t.Fatal(err)
	}
	if text, err := invoke(t, a, "read_file", map[string]any{"path": path, "start_line": 2, "end_line": 2}); err != nil || text != "beta\n" {
		t.Fatalf("%q %v", text, err)
	}
	if text, err := invoke(t, a, "list_directory", map[string]string{"path": dir}); err != nil || text != "nested/\n" {
		t.Fatalf("%q %v", text, err)
	}
	if text, err := invoke(t, a, "search_files", map[string]string{"path": dir, "pattern": "b.ta"}); err != nil || text != path+":2:beta\n" {
		t.Fatalf("%q %v", text, err)
	}
	if _, err := invoke(t, a, "edit_file", map[string]string{"path": path, "old_text": "beta", "new_text": "fixed"}); err != nil {
		t.Fatal(err)
	}
	if text, err := invoke(t, a, "read_file", map[string]string{"path": path}); err != nil || text != "alpha\nfixed\ngamma" {
		t.Fatalf("%q %v", text, err)
	}
	for _, old := range []string{"", "missing", "a"} {
		if _, err := invoke(t, a, "edit_file", map[string]string{"path": path, "old_text": old, "new_text": "BAD"}); err == nil {
			t.Fatalf("accepted anchor %q", old)
		}
	}
	data, _ := os.ReadFile(path)
	if string(data) != "alpha\nfixed\ngamma" {
		t.Fatal("failed edit changed disk")
	}
	if _, err := invoke(t, a, "write_file", map[string]any{"path": path, "content": "aaa", "overwrite": true}); err != nil {
		t.Fatal(err)
	}
	if _, err := invoke(t, a, "edit_file", map[string]string{"path": path, "old_text": "aa", "new_text": "BAD"}); err == nil {
		t.Fatal("overlapping ambiguous anchor accepted")
	}
	if _, err := invoke(t, a, "write_file", map[string]any{"path": path, "content": "", "overwrite": true}); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); len(data) != 0 {
		t.Fatal("empty write failed")
	}
	for _, test := range []struct {
		name string
		args any
	}{
		{"read_file", map[string]any{"path": path, "start_line": -1}}, {"read_file", map[string]any{"path": path, "end_line": 99}},
		{"read_file", map[string]string{"path": dir}}, {"read_file", map[string]string{"path": filepath.Join(dir, "absent")}},
		{"write_file", map[string]string{"path": dir, "content": "bad"}}, {"write_file", map[string]string{"path": "", "content": "bad"}},
		{"list_directory", map[string]string{"path": path}}, {"search_files", map[string]string{"path": dir, "pattern": "["}},
		{"search_files", map[string]string{"path": filepath.Join(dir, "absent"), "pattern": "x"}},
		{"edit_file", map[string]string{"path": filepath.Join(dir, "absent"), "old_text": "a", "new_text": "b"}},
	} {
		if _, err := invoke(t, a, test.name, test.args); err == nil {
			t.Fatalf("accepted %+v", test)
		}
	}
}

func TestChapter3StrictToolArguments(t *testing.T) {
	a := toolAgent(t, ensemble.Anthropic, true, false, "")
	for _, tool := range a.Tools() {
		switch tool.Definition().Name {
		case "wait_for_job", "send_input", "kill_job", "tool_limits":
			continue
		}
		for _, raw := range []string{`{`, `null`, `[]`, `{}`, `{"unknown":1}`, `{"path":null}`, `{"command":42}`, `{"path":"x","command":"x"}`, `{"command":" "}`} {
			if _, err := tool.Execute(context.Background(), json.RawMessage(raw)); err == nil {
				t.Fatalf("%s accepted %s", tool.Definition().Name, raw)
			}
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, tool := range a.Tools() {
		if _, err := tool.Execute(ctx, json.RawMessage(`{}`)); err != context.Canceled {
			t.Fatal(err)
		}
	}
}

func TestChapter3CommandContract(t *testing.T) {
	t.Chdir(t.TempDir())
	a := toolAgent(t, ensemble.Anthropic, true, false, "")
	for _, key := range []string{"LLM_API_KEY", "ANTHROPIC_API_KEY", "OPENAI_API_KEY", "GEMINI_API_KEY"} {
		t.Setenv(key, "benign-secret-canary")
	}
	t.Setenv("CH03_BUILD_SETTING", "preserved-canary")
	output, err := invoke(t, a, "run_command", map[string]string{"command": `printf '%s' "$CH03_BUILD_SETTING"; printf problem >&2; printf '%s%s%s%s' "$LLM_API_KEY" "$ANTHROPIC_API_KEY" "$OPENAI_API_KEY" "$GEMINI_API_KEY"; printf marker > created`})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "preserved-canaryproblem") || !strings.Contains(output, "exit_code: 0") || strings.Contains(output, "benign-secret-canary") {
		t.Fatal(output)
	}

	if data, err := os.ReadFile("created"); err != nil || string(data) != "marker" {
		t.Fatalf("%s %v", data, err)
	}
	if output, err := invoke(t, a, "run_command", map[string]string{"command": "printf out; printf err >&2; exit 7"}); err != nil || !strings.Contains(output, "exit_code: 7") || !strings.Contains(output, "outerr") {
		t.Fatalf("%s %v", output, err)
	}
	if _, err := invoke(t, a, "run_command", map[string]string{"command": "cd /; export CH03_SHELL_LOCAL=value"}); err != nil {
		t.Fatal(err)
	}
	if output, err := invoke(t, a, "run_command", map[string]string{"command": "test -z \"$CH03_SHELL_LOCAL\" && test -f created"}); err != nil {
		t.Fatalf("shell leaked state: %s %v", output, err)
	}
	if output, err := invoke(t, a, "run_command", map[string]string{"command": "head -c 70000 /dev/zero | tr '\\0' x"}); err != nil || !strings.Contains(output, "bytes omitted") {
		t.Fatalf("bounded output: %d %v", len(output), err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	for _, tool := range a.Tools() {
		if tool.Definition().Name == "run_command" {
			_, err = tool.Execute(ctx, json.RawMessage(`{"command":"sleep 30 & wait"}`))
		}
	}
	if err == nil || !strings.Contains(err.Error(), "deadline exceeded") || time.Since(start) > 3*time.Second {
		t.Fatalf("command cancellation %v %v", time.Since(start), err)
	}
}

func TestChapter3MissingIDsAndFollowup(t *testing.T) {
	dir := t.TempDir()
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		requests++
		parts := []ensemble.Part{prose("done")}
		if requests == 1 {
			parts = []ensemble.Part{prose("I will write two markers."), callPart("", "write_file", map[string]string{"path": filepath.Join(dir, "a"), "content": "A"}), callPart("", "write_file", map[string]string{"path": filepath.Join(dir, "b"), "content": "B"})}
		}
		if requests == 3 && (!bytes.Contains(data, []byte("file written")) || !bytes.Contains(data, []byte("done"))) {
			t.Errorf("followup lost results: %s", data)
		}
		_ = json.NewEncoder(w).Encode(providerReply(ensemble.Gemini, parts))
	}))
	defer server.Close()
	a := toolAgent(t, ensemble.Gemini, true, true, server.URL)
	for _, question := range []string{"first", "followup"} {
		if _, err := a.Ask(context.Background(), question); err != nil {
			t.Fatal(err)
		}
	}
	ids := map[string]bool{}
	for _, e := range a.Log().Events() {
		if e.Type == "tool_returned" {
			if e.Tool.CallID == "" || ids[e.Tool.CallID] {
				t.Fatal("missing/duplicate correlation")
			}
			ids[e.Tool.CallID] = true
		}
	}
	if len(ids) != 2 || len(a.History()) != 4 || a.Log().Context().Turn != ensemble.Idle {
		t.Fatalf("ids=%v history=%v state=%s", ids, a.History(), a.Log().Context().Turn)
	}
	var saved bytes.Buffer
	if err := a.Log().Dump(&saved); err != nil {
		t.Fatal(err)
	}
	loaded := toolAgent(t, ensemble.Gemini, true, true, "")
	if err := loaded.Log().Load(&saved); err != nil {
		t.Fatal(err)
	}
	if len(loaded.Log().Events()) != len(a.Log().Events()) {
		t.Fatal("log lost tool events")
	}
}

func TestChapter3ReadCapsWriteGuardAndSearchWindows(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	a := toolAgent(t, ensemble.Anthropic, true, false, "")
	path := "f.txt"
	planted := "one\nhit\nthree\nhit\nfive\nsix\nseven\neight\nhit\nten\n"
	if _, err := invoke(t, a, "write_file", map[string]string{"path": path, "content": planted}); err != nil {
		t.Fatal(err)
	}
	if _, err := invoke(t, a, "write_file", map[string]string{"path": path, "content": "BAD"}); err == nil || !strings.Contains(err.Error(), "bytes") || !strings.Contains(err.Error(), "overwrite:true") {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); string(data) != planted {
		t.Fatal("unguarded overwrite")
	}
	if output, err := invoke(t, a, "read_file", map[string]any{"path": path, "start_line": 2, "end_line": 4, "max_bytes": 5}); err != nil || output != "hit\nt\n[truncated at 5 bytes]" {
		t.Fatalf("%q %v", output, err)
	}
	if output, err := invoke(t, a, "read_file", map[string]any{"path": path, "start_line": 0, "end_line": 0, "max_bytes": 0}); err != nil || output != planted {
		t.Fatalf("%q %v", output, err)
	}
	expected := "f.txt-1-one\nf.txt:2:hit\nf.txt-3-three\nf.txt:4:hit\nf.txt-5-five\n--\nf.txt-8-eight\nf.txt:9:hit\nf.txt-10-ten\n"
	if output, err := invoke(t, a, "search_files", map[string]any{"path": path, "pattern": "hit", "context_lines": 1}); err != nil || output != expected {
		t.Fatalf("%q %v", output, err)
	}
	if output, err := invoke(t, a, "search_files", map[string]any{"path": path, "pattern": "one|three", "context_lines": 0}); err != nil || output != "f.txt:1:one\nf.txt:3:three\n" {
		t.Fatalf("%q %v", output, err)
	}
	if output, err := invoke(t, a, "search_files", map[string]any{"path": path, "pattern": "one|four", "context_lines": 2}); err != nil || output != "f.txt:1:one\nf.txt-2-hit\nf.txt-3-three\n" {
		t.Fatalf("%q %v", output, err)
	}
	for _, flags := range []map[string]any{{"append": true}, {"append": true, "overwrite": true}} {
		flags["path"] = path
		flags["content"] = "tail"
		if _, err := invoke(t, a, "write_file", flags); err != nil {
			t.Fatal(err)
		}
	}
	if data, _ := os.ReadFile(path); string(data) != planted+"tailtail" {
		t.Fatal("append lost bytes")
	}
	if _, err := invoke(t, a, "write_file", map[string]any{"path": path, "content": "replacement", "overwrite": true}); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); string(data) != "replacement" {
		t.Fatal("explicit overwrite failed")
	}
	if _, err := invoke(t, a, "write_file", map[string]any{"path": "new-append", "content": "new", "append": true}); err != nil {
		t.Fatal(err)
	}
	if _, err := invoke(t, a, "write_file", map[string]any{"path": "large", "content": strings.Repeat("a", 40000)}); err != nil {
		t.Fatal(err)
	}
	if output, err := invoke(t, a, "read_file", map[string]any{"path": "large"}); err != nil || !strings.Contains(output, "bytes omitted") || len(output) > 16384 {
		t.Fatalf("length %d %v", len(output), err)
	}
	for _, bad := range []struct{ name, raw string }{
		{"write_file", `{"path":"f.txt","content":"x","overwrite":"true"}`},
		{"write_file", `{"path":"f.txt","content":null}`},
		{"edit_file", `{"path":"f.txt","old_text":"x"}`},
		{"edit_file", `{"path":"f.txt","old_text":"x","new_text":null}`},
		{"search_files", `{"path":"f.txt","pattern":"x","context_lines":-1}`},
		{"search_files", `{"path":"f.txt","pattern":"x","context_lines":1.5}`},
		{"read_file", `{"path":"f.txt","max_bytes":-1}`},
		{"read_file", `{"path":"f.txt","start_line":1.5}`},
	} {
		for _, tool := range a.Tools() {
			if tool.Definition().Name == bad.name {
				if _, err := tool.Execute(context.Background(), json.RawMessage(bad.raw)); err == nil {
					t.Fatalf("accepted %s", bad.raw)
				}
			}
		}
	}
	if _, err := invoke(t, a, "edit_file", map[string]string{"path": path, "old_text": "replacement", "new_text": ""}); err != nil {
		t.Fatal(err)
	}
	if output, err := invoke(t, a, "read_file", map[string]string{"path": path}); err != nil || output != "" {
		t.Fatalf("empty %q %v", output, err)
	}
}

func TestChapter3NativeResultIDsAndEmptyRegistryFeedback(t *testing.T) {
	for _, vendor := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI, ensemble.Gemini} {
		t.Run(vendor.String(), func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, _ := io.ReadAll(r.Body)
				requests++
				if requests == 1 {
					_ = json.NewEncoder(w).Encode(providerReply(vendor, []ensemble.Part{callPart("native-A", "unknown", map[string]any{}), callPart("native-B", "unknown", map[string]any{})}))
					return
				}
				if bytes.Contains(data, []byte(`"tools"`)) {
					t.Error("empty registry declared tools")
				}
				// Inspect actual result objects, not a substring also present in the call.
				var body any
				_ = json.Unmarshal(data, &body)
				seen := map[string]bool{}
				var visit func(any)
				visit = func(node any) {
					switch n := node.(type) {
					case []any:
						for _, v := range n {
							visit(v)
						}
					case map[string]any:
						if n["type"] == "tool_result" {
							id, _ := n["tool_use_id"].(string)
							seen[id] = true
						}
						if n["role"] == "tool" {
							id, _ := n["tool_call_id"].(string)
							seen[id] = true
						}
						if v, ok := n["functionResponse"].(map[string]any); ok {
							id, _ := v["id"].(string)
							seen[id] = true
						}
						for _, v := range n {
							visit(v)
						}
					}
				}
				visit(body)
				if len(seen) != 2 || !seen["native-A"] || !seen["native-B"] {
					t.Errorf("result ids=%v", seen)
				}
				if bytes.Count(data, []byte("unknown tool")) != 2 {
					t.Error("unknown calls unanswered")
				}
				_ = json.NewEncoder(w).Encode(providerReply(vendor, []ensemble.Part{prose("finished")}))
			}))
			defer server.Close()
			a := toolAgent(t, vendor, false, false, server.URL)
			if output, err := a.Ask(context.Background(), "call"); err != nil || output != "finished" {
				t.Fatalf("%q %v", output, err)
			}
			if requests != 2 {
				t.Fatal(requests)
			}
		})
	}
}
