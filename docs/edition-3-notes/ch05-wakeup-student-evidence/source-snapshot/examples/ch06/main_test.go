package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestApplicationOwnedWorkflow(t *testing.T) {
	var mu sync.Mutex
	counts := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			System   string `json:"system"`
			Messages []any  `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		system := body.System
		role := "reviewer"
		if strings.Contains(system, "You are the author") {
			role = "author"
		} else if strings.Contains(system, "You are the editor") {
			role = "editor"
		}
		mu.Lock()
		n := counts[role]
		counts[role]++
		mu.Unlock()
		parts := []any{map[string]any{"type": "text", "text": "done"}}
		if n == 0 {
			parts = []any{map[string]any{"type": "text", "text": "I will use the role tool."}}
			add := func(id, name string, args any) {
				parts = append(parts, map[string]any{"type": "tool_use", "id": id, "name": name, "input": args})
			}
			switch role {
			case "author":
				add("w", "write_draft", map[string]any{"text": "first draft"})
				add("c", "word_count", map[string]any{})
			case "editor":
				add("r", "read_draft", map[string]any{})
				add("e", "edit_draft", map[string]any{"text": "edited manuscript"})
			case "reviewer":
				wire, _ := json.Marshal(body.Messages)
				if !strings.Contains(string(wire), "edited manuscript") {
					t.Error("reviewer did not receive edited manuscript")
				}
				add("v", "review", map[string]any{"decision": "accept", "notes": "clear"})
			}
		}
		json.NewEncoder(w).Encode(map[string]any{"content": parts, "usage": map[string]int{"input_tokens": 10, "output_tokens": 4}})
	}))
	defer server.Close()
	t.Setenv("LLM_VENDOR", "anthropic")
	t.Setenv("LLM_MODEL", "claude-sonnet-4-6")
	t.Setenv("LLM_API_KEY", "fixture")
	t.Setenv("LLM_BASE_URL", server.URL)
	t.Setenv("LLM_REQUIRE_TOOL_EXPLANATION", "true")
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("CH06_LOG", filepath.Join(dir, "log.jsonl"))
	in, err := os.CreateTemp(dir, "input")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	in.WriteString("{\"kind\":\"prompt\",\"text\":\"write briefly\"}\n")
	in.Seek(0, 0)
	out, _ := os.CreateTemp(dir, "out")
	defer out.Close()
	oldIn, oldOut := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = in, out
	defer func() { os.Stdin, os.Stdout = oldIn, oldOut }()
	if err := run(); err != nil {
		t.Fatal(err)
	}
	text, _ := os.ReadFile("manuscript.txt")
	if string(text) != "edited manuscript" {
		t.Fatal(string(text))
	}
	review, _ := os.ReadFile("review.json")
	if !strings.Contains(string(review), `"decision":"accept"`) {
		t.Fatal(string(review))
	}
	output, _ := os.ReadFile(out.Name())
	for _, marker := range []string{"author_done", "editor_done", "reviewer_done", "part_final", "turn_ended", `"type":"tool_result"`, `"call_id":"w"`} {
		if !strings.Contains(string(output), marker) {
			t.Fatal("missing", marker)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	for _, role := range []string{"author", "editor", "reviewer"} {
		if counts[role] != 2 {
			t.Fatal(counts)
		}
	}
}
