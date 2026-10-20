package ensemble_test

import (
	"context"
	"ensemble"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStreamingArrivalAndSavedFinal(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"model\":\"actual\",\"choices\":[{\"delta\":{\"content\":\"early\"}}]}\n\n")
		w.(http.Flusher).Flush()
		select {
		case <-release:
		case <-r.Context().Done():
			return
		}
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\" final\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":2,\"completion_tokens\":3}}\n\ndata: [DONE]\n\n")
	}))
	defer server.Close()
	app := ensemble.New(nil)
	defer app.Shutdown()
	logPath := filepath.Join(t.TempDir(), "agent.jsonl")
	a, err := app.NewAgent(ensemble.Config{Vendor: ensemble.OpenAI, Model: "test", APIKey: "fake", BaseURL: server.URL, StreamCapability: ensemble.StreamSSE, LogPath: logPath})
	if err != nil {
		t.Fatal(err)
	}
	q := a.Observe()
	defer q.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	reply := a.Submit(ctx, "hello")
	observed := false
	for !observed {
		select {
		case <-q.Wake():
			for _, o := range q.Drain() {
				if d, ok := o.(ensemble.PartDelta); ok {
					if d.PartID != 0 || d.Chunk != "early" {
						t.Fatal(d)
					}
					observed = true
				}
			}
		case result := <-reply:
			t.Fatalf("completed before releasing final frame: %+v", result)
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
	}
	close(release)
	select {
	case result := <-reply:
		if result.Err != nil || result.Text != "early final" {
			t.Fatal(result)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	finalCount := 0
	for _, o := range q.Drain() {
		if f, ok := o.(ensemble.PartFinal); ok {
			finalCount++
			if f.PartID != 0 || f.Part.ID != 0 || f.Part.Text != "early final" {
				t.Fatal(f)
			}
			saved, err := os.ReadFile(logPath)
			if err != nil || !strings.Contains(string(saved), `"response_ended"`) || !strings.Contains(string(saved), "early final") {
				t.Fatalf("final not saved: %v %s", err, saved)
			}
		}
	}
	if finalCount != 1 {
		t.Fatalf("finals=%d", finalCount)
	}
}
