package gui_test

import (
	"context"
	gui "ensemble-gui-stub"
	"ensemble/ensemble"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// Reviewer found repeated selection could fill the reliable replay queue while
// holding the observation mutex. A slow receiver must disconnect, not deadlock
// every Agent and Server.Close. Keep the real socket writer and its backpressure.
func TestSelectionSlowReceiverClose(t *testing.T) {
	root := ensemble.New(io.Discard)
	a, err := root.NewAgent(ensemble.Config{Model: "fake", DataDir: t.TempDir(), Recall: ensemble.RecallSettings{Disabled: true}})
	if err != nil {
		t.Fatal(err)
	}
	defer a.Shutdown()
	if err = a.Record(ensemble.Event{Type: "message_received", Message: &ensemble.MessageData{Actor: ensemble.Human, Parts: []ensemble.Part{{Type: "text", Text: strings.Repeat("x", 64*1024)}}}}); err != nil {
		t.Fatal(err)
	}
	hub, err := gui.New(root, a, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(hub.Handler(""))
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	for i := 0; i < 320; i++ {
		if err = conn.WriteJSON(map[string]string{"type": "select_agent", "key": "main"}); err != nil {
			break
		}
	}
	// Allow the admitted requests to fill the socket and application queues. The
	// peer intentionally reads no replay; no fixture replaces the actual pressure.
	time.Sleep(3 * time.Second)
	conn.Close()
	done := make(chan struct{})
	go func() { hub.Close(); close(done) }()
	select {
	case <-done:
		server.Close()
	case <-time.After(time.Second):
		t.Fatal("GUI close blocked after slow selection client disconnected")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = a.Flush(ctx); err != nil {
		t.Fatal("slow receiver stranded actor", err)
	}
}
