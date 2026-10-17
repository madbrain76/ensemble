package ensemble_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"ensemble"
)

func actorFixture(t *testing.T, h http.HandlerFunc, modify func(*ensemble.Config)) (ensemble.Agent, *ensemble.Ensemble) {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	root := ensemble.New(log.New(io.Discard, "", 0))
	cfg := ensemble.Config{Vendor: ensemble.Anthropic, Model: "claude-sonnet-4-6", BaseURL: server.URL, APIKey: "fixture", System: "Test actor behavior.", MaxTokens: 128, DataDir: t.TempDir(), HistoryPath: filepath.Join(t.TempDir(), "history.jsonl")}
	if modify != nil {
		modify(&cfg)
	}
	a, err := root.NewAgent(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := root.Shutdown(); err != nil {
			t.Error(err)
		}
	})
	return a, root
}
func wait5(t *testing.T, ch <-chan ensemble.Result) ensemble.Result {
	t.Helper()
	select {
	case r := <-ch:
		return r
	case <-time.After(3 * time.Second):
		t.Fatal("request completion timed out")
		return ensemble.Result{}
	}
}
func TestChapter5HintsPendingContextAndReliableQueue(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	var count atomic.Int32
	payloads := make(chan string, 3)
	a, _ := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		payloads <- string(b)
		n := count.Add(1)
		if n == 1 {
			close(entered)
			<-release
			json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, []ensemble.Part{prose("I will count."), callPart("c", "count", map[string]any{})}))
			return
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, []ensemble.Part{prose("done")}))
	}, nil)
	a.RegisterTool(ensemble.ToolDefinition{Name: "count", Description: "count", Parameters: json.RawMessage(`{"type":"object"}`)}, func(ensemble.ToolContext, json.RawMessage) (string, error) { return "one", nil })
	first := a.Submit(context.Background(), "first")
	<-entered
	second := a.Submit(context.Background(), "second")
	if err := a.Post(ensemble.Hint{Text: "be concise"}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := a.Wait(ctx, func(o ensemble.Observation) bool {
		v, ok := o.(ensemble.PartDelta)
		return ok && v.Chunk == "hint: be concise"
	})
	if err != nil {
		t.Fatal(err)
	}
	close(release)
	if r := wait5(t, first); r.Err != nil || r.Text != "done" {
		t.Fatal(r)
	}
	if r := wait5(t, second); r.Err != nil {
		t.Fatal(r)
	}
	initial, next, last := <-payloads, <-payloads, <-payloads
	if strings.Contains(initial, "be concise") || !strings.Contains(next, "be concise") || strings.Contains(last, "be concise") {
		t.Fatal("hint did not feed exactly the next request")
	}
	for _, turn := range a.Log().Context().Dialogue {
		if turn.Actor == "human" {
			for _, p := range turn.Parts {
				if p.Text == "be concise" {
					t.Fatal("hint became durable user dialogue")
				}
			}
		}
	}
	saved, err := os.ReadFile(a.Config().HistoryPath)
	if err != nil || !bytes.Contains(saved, []byte(`"type":"turn_ended"`)) {
		t.Fatalf("reply preceded saved end: %s %v", saved, err)
	}
}
func TestChapter5InterruptRetainsRealJobAndShutdownJoins(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	var requests atomic.Int32
	var skipped atomic.Int32
	a, root := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		n := requests.Add(1)
		if n > 1 {
			json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, []ensemble.Part{prose("fresh answer")}))
			return
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, []ensemble.Part{prose("I will wait then act."), callPart("slow", "slow", map[string]any{}), callPart("skip", "skip", map[string]any{})}))
	}, nil)
	a.RegisterTool(ensemble.ToolDefinition{Name: "slow", Description: "slow", Parameters: json.RawMessage(`{"type":"object"}`)}, func(ensemble.ToolContext, json.RawMessage) (string, error) {
		close(entered)
		<-release
		return "real late output", nil
	})
	a.RegisterTool(ensemble.ToolDefinition{Name: "skip", Description: "skip", Parameters: json.RawMessage(`{"type":"object"}`)}, func(ensemble.ToolContext, json.RawMessage) (string, error) { skipped.Add(1); return "bad", nil })
	reply := a.Submit(context.Background(), "begin")
	<-entered
	if err := a.Post(ensemble.Interrupt{}); err != nil {
		t.Fatal(err)
	}
	if r := wait5(t, reply); !errors.Is(r.Err, ensemble.ErrInterrupted) {
		t.Fatal(r)
	}
	if a.State() != ensemble.Interrupted {
		t.Fatal(a.State())
	}
	// The old handler is still blocked: interruption must free the Agent now,
	// not merely let a follow-up succeed after the old work has completed.
	fresh := wait5(t, a.Submit(context.Background(), "new turn now"))
	if fresh.Err != nil || fresh.Text != "fresh answer" {
		t.Fatal(fresh)
	}
	old, _ := root.Jobs().Get(1)
	if old.Snapshot().Status != "running" {
		t.Fatal("old job ended before new answer")
	}
	stop := make(chan error, 1)
	go func() { stop <- root.Shutdown() }()
	select {
	case err := <-stop:
		t.Fatalf("shutdown did not join Go handler: %v", err)
	case <-time.After(30 * time.Millisecond):
	}
	close(release)
	select {
	case err := <-stop:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown failed to finish")
	}
	if requests.Load() != 2 || skipped.Load() != 0 {
		t.Fatal("old turn resumed")
	}
	job, err := root.Jobs().Get(1)
	if err != nil {
		t.Fatal(err)
	}
	record := job.Snapshot()
	if record.Status != "done" {
		t.Fatal(record)
	}
	output, err := job.Wait(context.Background(), ensemble.WaitOptions{Delay: time.Millisecond, MaxBytes: 8192})
	if err != nil || !strings.Contains(output, "real late output") {
		t.Fatalf("late real output lost: %q %v", output, err)
	}
	var returned, ends, reports int
	for _, e := range a.Log().Events() {
		switch e.Type {
		case "tool_returned":
			returned++
		case "turn_ended":
			ends++
		case "job_report":
			reports++
		}
	}
	if returned != 2 || ends != 2 || reports != 1 {
		t.Fatalf("facts returned=%d ends=%d reports=%d", returned, ends, reports)
	}
	if err := a.Post(ensemble.Hint{Text: "late"}); !errors.Is(err, ensemble.ErrStopped) {
		t.Fatal(err)
	}
}
func TestChapter5LossyObservationCannotLoseReliableReply(t *testing.T) {
	a, root := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, []ensemble.Part{prose("ok")}))
	}, nil)
	// Do not drain either bounded observation queue.
	for i := 0; i < 90; i++ {
		if r := wait5(t, a.Submit(context.Background(), "hello")); r.Err != nil || r.Text != "ok" {
			t.Fatal(r)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	seen, err := root.WaitAny(ctx)
	if err != nil || seen.Agent != a {
		t.Fatal(seen, err)
	}
}
func TestChapter5ConfiguredBatchBudgetAndFinalAnswer(t *testing.T) {
	if ensemble.DefaultToolBatchLimit != 200 {
		t.Fatal("chapter default batch budget changed")
	}
	for _, final := range []bool{false, true} {
		t.Run(map[bool]string{true: "final", false: "bounded"}[final], func(t *testing.T) {
			var requests, calls atomic.Int32
			a, _ := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
				n := requests.Add(1)
				parts := []ensemble.Part{prose("I will count."), callPart("c", "count", map[string]any{})}
				if n == 3 && final {
					parts = []ensemble.Part{prose("final answer")}
				}
				json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, parts))
			}, func(c *ensemble.Config) { c.ToolBatchLimit = 2 })
			a.RegisterTool(ensemble.ToolDefinition{Name: "count", Description: "count", Parameters: json.RawMessage(`{"type":"object"}`)}, func(ensemble.ToolContext, json.RawMessage) (string, error) { calls.Add(1); return "counted", nil })
			r := wait5(t, a.Submit(context.Background(), "go"))
			if final && r.Err != nil {
				t.Fatal(r)
			}
			if !final && !errors.Is(r.Err, ensemble.ErrToolBatchLimit) {
				t.Fatal(r)
			}
			if calls.Load() != 2 || requests.Load() != 3 {
				t.Fatal(calls.Load(), requests.Load())
			}
		})
	}
}
func TestChapter5MediaBytesCapabilityAndRedaction(t *testing.T) {
	png, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/lVIAAAAASUVORK5CYII=")
	path := filepath.Join(t.TempDir(), "tiny.png")
	os.WriteFile(path, png, 0600)
	for _, v := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI, ensemble.Gemini} {
		t.Run(v.String(), func(t *testing.T) {
			root := ensemble.New(log.New(io.Discard, "", 0))
			defer root.Shutdown()
			a, err := root.NewAgent(ensemble.Config{Vendor: v, Model: map[ensemble.Vendor]string{ensemble.Anthropic: "claude-sonnet-4-6", ensemble.OpenAI: "gpt-4.1-2025-04-14", ensemble.Gemini: "gemini-3.8-flash"}[v], APIKey: "fixture", System: "media", MaxTokens: 128, BaseURL: "http://fixture.invalid"})
			if err != nil {
				t.Fatal(err)
			}
			part := ensemble.Part{Type: "blob", MIME: "image/png", Ref: &ensemble.Ref{Kind: ensemble.RefPath, Locator: path}}
			c := ensemble.Context{Dialogue: []ensemble.Turn{{Actor: "human", Parts: []ensemble.Part{part}}}}
			wire, err := a.Engine().Render(c, a.Engine().Target())
			if err != nil || !bytes.Contains(wire, []byte(base64.StdEncoding.EncodeToString(png))) {
				t.Fatalf("bytes not rendered: %s %v", wire, err)
			}
			if len(c.Dialogue[0].Parts[0].Data) != 0 {
				t.Fatal("render mutated history data")
			}
			a.Log().Append(ensemble.Event{Type: "message_received", Message: &ensemble.MessageEvent{Actor: "human", Parts: []ensemble.Part{part}}})
			var history bytes.Buffer
			a.Log().Dump(&history)
			if bytes.Contains(history.Bytes(), []byte(base64.StdEncoding.EncodeToString(png))) {
				t.Fatal("history contains bytes")
			}
			unknown := a.Engine().Target()
			unknown.Model = "unknown-model"
			if _, err := a.Engine().Render(c, unknown); err == nil || !strings.Contains(err.Error(), "unknown-model") || !strings.Contains(err.Error(), "image") {
				t.Fatal(err)
			}
			c.Dialogue[0].Parts[0] = ensemble.Part{Type: "redacted", Parts: []ensemble.Part{part}}
			os.Remove(path)
			if _, err := a.Engine().Render(c, unknown); err != nil {
				t.Fatal("redaction opened path", err)
			}
			os.WriteFile(path, png, 0600)
		})
	}
	root := ensemble.New(log.New(io.Discard, "", 0))
	defer root.Shutdown()
	a, _ := root.NewAgent(ensemble.Config{Vendor: ensemble.Gemini, Model: "gemini-3.8-flash", APIKey: "fixture", System: "media", MaxTokens: 128, BaseURL: "http://fixture.invalid"})
	for _, mime := range []string{"audio/wav", "video/mp4", "application/pdf"} {
		p := filepath.Join(t.TempDir(), "fixture")
		data := []byte("fixture " + mime)
		os.WriteFile(p, data, 0600)
		wire, err := a.Engine().Render(ensemble.Context{Dialogue: []ensemble.Turn{{Actor: "human", Parts: []ensemble.Part{{Type: "blob", MIME: mime, Ref: &ensemble.Ref{Kind: ensemble.RefPath, Locator: p}}}}}}, a.Engine().Target())
		if err != nil || !bytes.Contains(wire, []byte(base64.StdEncoding.EncodeToString(data))) {
			t.Fatal(mime, string(wire), err)
		}
	}
}
