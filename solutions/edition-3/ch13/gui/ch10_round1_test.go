package gui

import (
	"context"
	"encoding/json"
	"ensemble"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func ch10Event(t *testing.T, raw string) ensemble.Event {
	t.Helper()
	var e ensemble.Event
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatal(err)
	}
	return e
}
func TestCh10R1SnapshotPrefixTailWS(t *testing.T) {
	for _, mode := range []string{"snapshot_tail", "context_only", "full_carry"} {
		t.Run(mode, func(t *testing.T) { ch10PrefixTail(t, mode) })
	}
}
func ch10PrefixTail(t *testing.T, mode string) {
	dir := t.TempDir()
	path := filepath.Join(dir, "save.json")
	seed := ensemble.New(nil)
	a, err := seed.NewAgent(ensemble.Config{Model: "fixture", Tools: []string{}, SavePath: path})
	if err != nil {
		t.Fatal(err)
	}
	if err = a.Log().Append(ch10Event(t, `{"type":"response_ended","turn_id":9,"response":{"from":{"vendor":"anthropic","surface":"messages","model":"fixture"},"parts":[{"type":"tool_call","from":{"vendor":"anthropic","surface":"messages","model":"fixture"},"display_id":7,"call_id":"native-call","name":"echo","args":{}}]}}`)); err != nil {
		t.Fatal(err)
	}
	prefix := a.Log().Context()
	for _, raw := range []string{
		`{"type":"tool_called","turn_id":9,"tool":{"call_id":"native-call","name":"echo","display_id":7,"call_part_id":7,"args":{}}}`,
		`{"type":"tool_dispatched","turn_id":9,"tool":{"call_id":"native-call","name":"echo","display_id":7,"call_part_id":7,"args":{},"dispatched":true}}`,
		`{"type":"tool_returned","turn_id":9,"tool":{"call_id":"native-call","name":"echo","display_id":8,"call_part_id":7,"parts":[{"type":"text","text":"tail-result"}]}}`,
	} {
		if err = a.Log().Append(ch10Event(t, raw)); err != nil {
			t.Fatal(err)
		}
	}
	if err = seed.Shutdown(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var s ensemble.SaveFile
	if err = json.Unmarshal(raw, &s); err != nil {
		t.Fatal(err)
	}
	tail := append([]ensemble.Event(nil), s.Log[1:]...)
	if mode != "full_carry" {
		s.AsOf = 1
		s.Context = &prefix
		s.Log = tail
	}
	if mode == "context_only" {
		s.Log = nil
		s.LogMode = ensemble.LogDialogueOnly
	}
	raw, err = json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	peer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"model":"fixture","content":[{"type":"text","text":"new-display"}],"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`)
	}))
	defer peer.Close()
	root := ensemble.New(nil)
	defer root.Shutdown()
	cfg := ensemble.Config{Vendor: ensemble.Anthropic, Model: "fixture", APIKey: "fake-local-only", BaseURL: peer.URL, DisableStreaming: true, Tools: []string{}, SavePath: path}
	loaded, err := root.NewAgent(cfg)
	if err != nil {
		t.Fatal(err)
	}
	server, err := NewServerWithSettings(root, "web/gui", filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	web := httptest.NewServer(server.Handler())
	defer web.Close()

	connect := func(url string) (*websocket.Conn, []Wire) {
		t.Helper()
		ws, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(url, "http")+"/ws", nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = ws.SetReadDeadline(time.Now().Add(3 * time.Second))
		var frame struct {
			Type string
			Data replayEnvelope
		}
		if err = ws.ReadJSON(&frame); err != nil {
			t.Fatal(err)
		}
		if frame.Type != "replay" {
			t.Fatal(frame.Type)
		}
		return ws, frame.Data.Events
	}
	ws, initial := connect(web.URL)
	defer ws.Close()
	var firstID uint64
	for _, f := range initial {
		if f.Part != nil && f.Part.Kind == "tool_call" && f.Part.CallID == "native-call" {
			firstID = f.PartID
			break
		}
	}
	if firstID == 0 {
		t.Fatal("missing prefix card")
	}
	if mode == "context_only" {
		for _, e := range tail {
			e.Seq = 0
			if err = loaded.Log().Append(e); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err = loaded.Log().Append(ch10Event(t, `{"type":"job_report","turn_id":9,"job":{"handle":1,"status":"done","turn_id":9,"call_id":"native-call","call_part_id":7,"name":"echo","bytes":5,"exit_code":0}}`)); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err = loaded.Ask(ctx, "new prompt"); err != nil {
		t.Fatal(err)
	}
	// Ask supplies reliable completion. Read the actual live WS projection through
	// its final new text, including retained tail and late job frames before it.
	frames := append([]Wire(nil), initial...)
	for {
		var f Wire
		if err = ws.ReadJSON(&f); err != nil {
			t.Fatal(err)
		}
		frames = append(frames, f)
		if f.Part != nil && f.Part.Text == "new-display" {
			break
		}
	}
	var newID, resultID uint64
	check := func(frames []Wire) {
		t.Helper()
		var calls, jobs, results int
		var gotNew uint64
		for _, f := range frames {
			if f.Part == nil {
				continue
			}
			if f.Part.Text == "new-display" {
				gotNew = f.PartID
			}
			if f.Part.CallID != "native-call" {
				continue
			}
			if f.Type == "job_update" {
				jobs++
				if f.PartID != firstID || f.Part.CallPartID != firstID || f.Part.Job == nil || f.Part.Job.CallPartID != firstID {
					t.Errorf("late job association: %+v", f)
				}
			}
			switch f.Part.Kind {
			case "tool_call":
				calls++
				if f.PartID != firstID {
					t.Errorf("split prefix card: %d != %d", f.PartID, firstID)
				}
			case "tool_result":
				results++
				if f.Part.CallPartID != firstID || f.PartID == firstID {
					t.Errorf("result association: %+v", f)
				}
				if resultID == 0 {
					resultID = f.PartID
				} else if resultID != f.PartID {
					t.Error("result ID changed")
				}

			}
		}
		if calls < 2 || results == 0 || jobs == 0 || gotNew == 0 {
			t.Fatalf("missing projected frames: calls=%d results=%d jobs=%d new=%d", calls, results, jobs, gotNew)
		}
		if gotNew == firstID || gotNew == resultID {
			t.Error("new display reused loaded ID")
		}
		if newID == 0 {
			newID = gotNew
		} else if newID != gotNew {
			t.Error("new display ID changed on replay/restore")
		}
	}
	check(frames)
	again, replay := connect(web.URL)
	_ = again.Close()
	check(replay)
	if loaded.Log().Context().Dialogue[0].Parts[0].DisplayID != 7 {
		t.Fatal("canonical prefix ID rewritten")
	}
	if err = loaded.Save(); err != nil {
		t.Fatal(err)
	}
	if err = root.Shutdown(); err != nil {
		t.Fatal(err)
	}
	// Reopen exactly the produced archive, keeping current config independently.
	reopened := ensemble.New(nil)
	defer reopened.Shutdown()
	if _, err = reopened.NewAgent(cfg); err != nil {
		t.Fatal(err)
	}
	restored, err := NewServerWithSettings(reopened, "web/gui", filepath.Join(dir, "restored-settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()
	httpRestored := httptest.NewServer(restored.Handler())
	defer httpRestored.Close()
	restoredWS, replayed := connect(httpRestored.URL)
	_ = restoredWS.Close()
	check(replayed)
}

func TestCh10R1RetentionSettings(t *testing.T) {
	root := ensemble.New(nil)
	defer root.Shutdown()
	a, err := root.NewAgent(ensemble.Config{Model: "fixture", Tools: []string{}, DataDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewServerWithSettings(root, "web/gui", filepath.Join(t.TempDir(), "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if s.Hub.currentSettings().DiagnosticMaxRecords != 200 || *a.Config().DiagnosticMaxRecords != 200 {
		t.Fatal("wrong default")
	}
	if err = s.Hub.updateSettings(json.RawMessage(`{"diagnostic_max_records":0}`)); err != nil {
		t.Fatal(err)
	}
	saved, err := s.Hub.store.load()
	if err != nil || saved.DiagnosticMaxRecords != 0 || *a.Config().DiagnosticMaxRecords != 0 {
		t.Fatalf("zero not retained: %+v %v", saved, err)
	}
	if err = s.Hub.updateSettings(json.RawMessage(`{"diagnostic_max_records":-1}`)); err == nil {
		t.Fatal("negative limit accepted")
	}
	after, err := s.Hub.store.load()
	if err != nil || after != saved || *a.Config().DiagnosticMaxRecords != 0 {
		t.Fatal("negative patch changed accepted state")
	}
}
