package ensemble_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ensemble/ensemble"
	"ensemble/internal/common"
)

// The fake replaces only HTTP; actual tools, jobs, skill registry, actor and
// reducer run. Requests are receipts, not a scripted mock of internal calls.
func contextSession(t *testing.T, model string, script func(int) []map[string]any, spec ensemble.AgentSpec) (ensemble.Agent, [][]byte) {
	t.Helper()
	var requests [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		requests = append(requests, raw)
		if model == "claude-opus-5-course" && r.Header.Get("anthropic-beta") != "inline-tools-2026-09-15" {
			t.Error("inline beta missing")
		}
		content := script(len(requests))
		if content == nil {
			content = []map[string]any{{"type": "text", "text": "done"}}
		}
		json.NewEncoder(w).Encode(map[string]any{"content": content, "usage": map[string]int{"input_tokens": 1, "output_tokens": 1}})
	}))
	t.Cleanup(server.Close)
	a, err := ensemble.New(io.Discard).NewAgent(ensemble.Config{Model: model, BaseURL: server.URL, BuiltinTools: true, DataDir: t.TempDir(), ContextTarget: 20000}, spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { a.Shutdown() })
	if _, err := a.Ask(context.Background(), "Read, checkpoint and continue."); err != nil {
		t.Fatal(err)
	}
	return a, requests
}
func contextCall(id, name string, args map[string]any) map[string]any {
	return map[string]any{"type": "tool_use", "id": id, "name": name, "input": args}
}
func toolResult(t *testing.T, raw []byte, id string) string {
	t.Helper()
	var wire struct {
		// Messages contains the captured provider conversation.
		Messages []struct {
			// Content holds result blocks, separately from earlier call arguments.
			Content []struct {
				// ID identifies the actual result under inspection.
				ID string `json:"tool_use_id"`
				// Content is its nested text, never the model's call input.
				Content []struct {
					// Text is the actual nested result text.
					Text string `json:"text"`
				} `json:"content"`
			} `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatal(err)
	}
	for _, m := range wire.Messages {
		for _, p := range m.Content {
			if p.ID == id {
				var text string
				for _, b := range p.Content {
					text += b.Text
				}
				return text
			}
		}
	}
	return ""
}

// Keep governs the preceding batch, even alongside the next file read. A
// non-stubbing model is the positive control for the feature-table distinction.
func TestContextKeepPreviousBatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "result.txt")
	os.WriteFile(path, []byte(strings.Repeat("EVIDENCE ", 100)), 0600)
	for _, model := range []string{"claude-opus-5-course", "claude-sonnet-5-course"} {
		t.Run(model, func(t *testing.T) {
			a, requests := contextSession(t, model, func(n int) []map[string]any {
				switch n {
				case 1:
					return []map[string]any{contextCall("a", "read_file", map[string]any{"path": path})}
				case 2:
					return []map[string]any{contextCall("keep", "keep_tool_results", map[string]any{}), contextCall("b", "read_file", map[string]any{"path": path})}
				case 3:
					return []map[string]any{contextCall("c", "read_file", map[string]any{"path": path})}
				}
				return nil
			}, ensemble.AgentSpec{})
			if len(requests) != 4 {
				t.Fatalf("requests=%d", len(requests))
			}
			if !strings.Contains(toolResult(t, requests[1], "a"), "EVIDENCE") || !strings.Contains(toolResult(t, requests[2], "b"), "EVIDENCE") {
				t.Fatal("unseen result cut")
			}
			if !strings.Contains(toolResult(t, requests[3], "a"), "EVIDENCE") {
				t.Fatal("keep did not preserve preceding batch")
			}
			stubbed := strings.Contains(toolResult(t, requests[3], "b"), "redacted")
			if stubbed != (model == "claude-opus-5-course") {
				t.Fatal("stubbing ignored model or kept wrong batch")
			}
			if strings.Contains(toolResult(t, requests[3], "keep"), "redacted") {
				t.Fatal("keep result stubbed")
			}
			if len(a.History().EventsAfter(0)) == 0 {
				t.Fatal("real history was not exercised")
			}
		})
	}
}

// A handoff shares a reply with another tool, which must finish before clearing.
// Skill text and declaration changes survive; the prefix is byte-identical.
func TestContextSurvivorsAndFrozenPrefix(t *testing.T) {
	dir := skillFiles(t, map[string]string{
		"base":  "---\nname: base\ntype: primary\ntools: read_file\nloadable-skills: extra\n---\nBASE",
		"extra": "---\nname: extra\ntype: loadable\ntools: list_directory\n---\nSURVIVOR MANUAL",
	})
	path := filepath.Join(t.TempDir(), "file")
	os.WriteFile(path, []byte("TOOL BYTES"), 0600)
	a, requests := contextSession(t, "claude-opus-5-course", func(n int) []map[string]any {
		switch n {
		case 1:
			return []map[string]any{contextCall("load", "load_skill", map[string]any{"name": "extra"})}
		case 2:
			return []map[string]any{contextCall("note", "micro_handoff", map[string]any{"text": "HANDOFF NOTE"}), contextCall("read", "read_file", map[string]any{"path": path})}
		case 3:
			return []map[string]any{contextCall("unload", "unload_skill", map[string]any{"name": "extra"})}
		}
		return nil
	}, ensemble.AgentSpec{SkillsDir: dir, PrimarySkill: "base"})
	if len(requests) != 4 {
		t.Fatalf("requests=%d", len(requests))
	}
	var initial map[string]json.RawMessage
	json.Unmarshal(requests[0], &initial)
	for _, request := range requests[1:] {
		var got map[string]json.RawMessage
		json.Unmarshal(request, &got)
		if !bytes.Equal(got["tools"], initial["tools"]) || !bytes.Equal(got["system"], initial["system"]) {
			t.Fatal("prefix changed")
		}
	}
	if strings.Contains(toolResult(t, requests[1], "load"), "SURVIVOR") {
		t.Fatal("manual still inside tool result")
	}
	if strings.Count(string(requests[2]), "HANDOFF NOTE") != 1 || !bytes.Contains(requests[2], []byte("SURVIVOR MANUAL")) {
		t.Fatal("survivor missing or handoff duplicated")
	}
	if bytes.Contains(requests[2], []byte(`"tool_use"`)) || bytes.Contains(requests[2], []byte(`"tool_result"`)) {
		t.Fatal("handoff left batch tool parts")
	}
	if !bytes.Contains(requests[1], []byte(`"tool_addition"`)) || !bytes.Contains(requests[3], []byte(`"tool_removal"`)) {
		t.Fatal("declaration delta missing")
	}
	for _, d := range a.Declarations() {
		if d.Name == "list_directory" {
			t.Fatal("unloaded tool still executable")
		}
	}
	// Raw chronology retains the ordinary result before the handoff event, and
	// the sibling read still really ran despite its projection being cleared.
	result, handoff, read := uint64(0), uint64(0), uint64(0)
	for _, event := range a.History().EventsAfter(0) {
		if event.Type == "micro_handoff" {
			handoff = event.Seq
		}
		if event.Type == "tool_returned" {
			if event.Tool.CallID == "note" {
				result = event.Seq
			}
			if event.Tool.CallID == "read" && !event.Tool.IsError {
				read = event.Seq
			}
		}
	}
	if !(result > 0 && result < handoff && handoff < read) {
		t.Fatal("handoff chronology or sibling execution changed")
	}
}

// Long real paths fill the calls band as well as the result band. Changing the
// target after capture cannot change recorded cuts during snapshot/log replay.
func TestContextLadderRecordedAndReplay(t *testing.T) {
	dir := filepath.Join(t.TempDir(), strings.Repeat("p", 200), strings.Repeat("q", 200))
	os.MkdirAll(dir, 0700)
	path := filepath.Join(dir, "data")
	os.WriteFile(path, []byte(strings.Repeat("LADDER DATA ", 800)), 0600)
	a, requests := contextSession(t, "claude-sonnet-5-course", func(n int) []map[string]any {
		if n <= 12 {
			return []map[string]any{contextCall(fmt.Sprint(n), "read_file", map[string]any{"path": path})}
		}
		return nil
	}, ensemble.AgentSpec{})
	cuts := map[string]int{}
	for _, e := range a.History().EventsAfter(0) {
		if e.Redact != nil {
			cuts[e.Redact.Level]++
			if e.Redact.To == 0 {
				t.Fatal("cut has no numeric boundary")
			}
		}
	}
	if cuts["redact_result"] == 0 || cuts["redact_tool"] == 0 {
		t.Fatalf("ladder did not exercise both bands: %v", cuts)
	}
	for i := 1; i < len(requests); i++ {
		if !strings.Contains(toolResult(t, requests[i], fmt.Sprint(i)), "LADDER DATA") {
			t.Fatalf("result %d cut before first presentation", i)
		}
	}
	a.SetContextSettings(900000, 0)
	var save, log bytes.Buffer
	a.History().Save(&save)
	a.History().Dump(&log)
	want, _ := a.Engine().Render(a.History().Context())
	b, _ := ensemble.New(io.Discard).NewAgent(ensemble.Config{Model: "claude-sonnet-5-course", BuiltinTools: true, DataDir: t.TempDir(), ContextTarget: 900000})
	for _, restore := range []func() error{func() error { return b.History().Restore(bytes.NewReader(save.Bytes())) }, func() error { return b.History().Load(bytes.NewReader(log.Bytes())) }} {
		if err := restore(); err != nil {
			t.Fatal(err)
		}
		got, err := b.Engine().Render(b.History().Context())
		if err != nil || !bytes.Equal(got, want) {
			t.Fatal("recorded cuts changed on replay", err)
		}
	}
}

// Known malformed facts are atomic skips. Later valid facts still apply while
// unknown event types refuse the staged load and preserve the original state.
func TestContextTotalReducerAndJournal(t *testing.T) {
	dir := t.TempDir()
	var diagnostics bytes.Buffer
	a, _ := ensemble.New(&diagnostics).NewAgent(ensemble.Config{Model: "model", DataDir: dir})
	if err := a.Load(); err != nil {
		t.Fatal(err)
	}
	record := func(text string) {
		t.Helper()
		if err := a.Record(common.Event{Type: "message_received", Message: &common.MessageData{Actor: common.Human, Parts: []common.Part{{Type: "text", Text: text}}}}); err != nil {
			t.Fatal(err)
		}
	}
	record("BEFORE")
	if err := a.Save(); err != nil {
		t.Fatal(err)
	}
	first, _ := os.ReadFile(filepath.Join(dir, "save.json"))
	record("JOURNAL TAIL")
	b, _ := ensemble.New(&diagnostics).NewAgent(ensemble.Config{Model: "model", DataDir: dir})
	if err := b.Load(); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(b.History().Context())
	if !bytes.Contains(raw, []byte("JOURNAL TAIL")) {
		t.Fatal("journal tail lost")
	}
	b.SetContextSettings(0, 1)
	if err := b.Save(); err != nil {
		t.Fatal(err)
	}
	backup, _ := os.ReadFile(filepath.Join(dir, "save.json.bak"))
	if !bytes.Equal(first, backup) {
		t.Fatal("previous snapshot backup lost")
	}
	journal, _ := os.ReadFile(filepath.Join(dir, "save.json.journal"))
	if len(journal) != 0 {
		t.Fatal("checkpoint did not reset journal")
	}
	saved, _ := os.ReadFile(filepath.Join(dir, "save.json"))
	var sf common.SaveFile
	json.Unmarshal(saved, &sf)
	if len(sf.Log) != 1 || sf.Log[0].Seq > sf.AsOf {
		t.Fatal("retention crossed anchor")
	}
	good := `{"seq":1,"type":"message_received","message":{"actor":"human","parts":[{"type":"text","text":"FIRST"}]}}` + "\n"
	bad := `{"seq":2,"type":"message_received","message":{"actor":"human","parts":[{"type":"text","text":"MUST NOT LAND"},{"type":"mystery"}]}}` + "\n" + `{"seq":3,"type":"redacted","redact":{"from":999,"to":999,"level":"redact_tool"}}` + "\n"
	last := strings.Replace(strings.Replace(good, `"seq":1`, `"seq":4`, 1), "FIRST", "LAST", 1)
	if err := b.History().Load(strings.NewReader(good + bad + last)); err != nil {
		t.Fatal(err)
	}
	raw, _ = json.Marshal(b.History().Context())
	if bytes.Contains(raw, []byte("MUST NOT LAND")) || !bytes.Contains(raw, []byte("FIRST")) || !bytes.Contains(raw, []byte("LAST")) || strings.Count(diagnostics.String(), "skipped event") != 2 {
		t.Fatal("bad event not skipped atomically with diagnostic")
	}
	if err := b.History().Load(strings.NewReader(strings.Replace(good, "message_received", "future_event", 1))); err == nil {
		t.Fatal("unknown event accepted")
	}
	a.Shutdown()
	b.Shutdown()
}
