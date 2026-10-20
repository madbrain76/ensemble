package gui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"ensemble"
	"github.com/gorilla/websocket"
)

type Server struct {
	Root         *ensemble.Ensemble
	Hub          *Hub
	Assets       string
	DefaultAgent string
	Trace        io.Writer
	traceMu      sync.Mutex
}

func NewServer(root *ensemble.Ensemble, assets string) (*Server, error) {
	return NewServerWithSettings(root, assets, filepath.Join("cr", "settings.json"))
}

// NewServerWithSettings accepts the settings path selected by the launcher.
func NewServerWithSettings(root *ensemble.Ensemble, assets, path string) (*Server, error) {
	if path == "" {
		path = filepath.Join("cr", "settings.json")
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(assets)
	if err != nil {
		return nil, err
	}
	for _, name := range []string{"index.html", "app.js", "style.css"} {
		if info, err := os.Stat(filepath.Join(abs, name)); err != nil || info.IsDir() {
			return nil, fmt.Errorf("GUI asset missing: %s", filepath.Join(abs, name))
		}
	}
	store := settingsStore{path: path}
	settings, loadErr := store.load()
	if loadErr != nil {
		return nil, fmt.Errorf("load GUI settings: %w", loadErr)
	}
	return &Server{Root: root, Hub: newHub(root, settings, store), Assets: abs}, nil
}
func (s *Server) wireTrace(direction string, raw []byte) {
	if s.Trace == nil {
		return
	}
	var payload bytes.Buffer
	if err := json.Compact(&payload, raw); err != nil {
		// Invalid inbound frames still leave one readable diagnostic record.
		// Quoting keeps malformed input readable without extra record lines.
		payload.Reset()
		quoted, _ := json.Marshal(string(raw))
		payload.Write(quoted)
	}
	s.traceMu.Lock()
	defer s.traceMu.Unlock()
	if _, err := fmt.Fprintf(s.Trace, "%s %s %s\n", time.Now().Format(time.RFC3339Nano), direction, payload.Bytes()); err != nil {
		s.Hub.en.Logger().Printf("GUI diagnostic write: %v", err)
	}
}

func (s *Server) Close() { s.Hub.Close() }
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.socket)
	mux.Handle("/", http.FileServer(http.Dir(s.Assets)))
	return mux
}
func (s *Server) socket(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		u, err := url.Parse(origin)
		return err == nil && u.Host == r.Host
	}}
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := s.Hub.connect()
	defer s.Hub.disconnect(c)
	defer conn.Close()
	writerDone := make(chan struct{})
	// Ping payloads are at most 125 bytes. The reader admits a bounded reply
	// without writing the socket or waiting behind a slow socket writer.
	pongs := make(chan string, 16)
	// Exactly one goroutine writes every data/control frame for this socket.
	go func() {
		defer close(writerDone)
		defer conn.Close()
		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-c.done:
				return
			case raw := <-c.queue:
				_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				start := time.Now()
				s.wireTrace("<", raw)
				if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
					s.Hub.trace("socket_write_error", "", err.Error())
					return
				}
				s.Hub.trace("socket_write", "", fmt.Sprintf("client=%s bytes=%d elapsed=%s", c.id, len(raw), time.Since(start)))
			case payload := <-pongs:
				_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PongMessage, []byte(payload)); err != nil {
					s.Hub.trace("socket_write_error", "", err.Error())
					return
				}
			case <-ticker.C:
				_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
	conn.SetReadLimit(1 << 20)
	// Avoid the library's default reader-side pong/close writers.
	conn.SetPingHandler(func(payload string) error {
		select {
		case <-c.done:
			return io.EOF
		default:
		}
		select {
		case pongs <- payload:
			return nil
		default:
			// Do not silently discard a required Pong or block command reads.
			c.stop()
			s.Hub.trace("socket_control_overflow", "", c.id)
			return errors.New("websocket pong queue full")
		}
	})
	conn.SetCloseHandler(func(int, string) error { return nil })
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}
		s.wireTrace(">", raw)
		var cmd command
		if err := json.Unmarshal(raw, &cmd); err != nil {
			s.Hub.sendError(c, "", err)
			continue
		}
		s.Hub.trace("command", cmd.AgentID, cmd.Type)
		if err := s.command(c, cmd); err != nil {
			s.Hub.sendError(c, cmd.AgentID, err)
		}
	}
	c.stop()
	<-writerDone
}

type command struct {
	Name      string            `json:"name,omitempty"`
	Type      string            `json:"type"`
	Settings  json.RawMessage   `json:"settings"`
	AgentID   string            `json:"agent_id"`
	Text      string            `json:"text"`
	Parts     ensemble.PartList `json:"parts"`
	Enabled   bool              `json:"enabled"`
	Automatic bool              `json:"automatic,omitempty"`
	PartID    uint64            `json:"part_id"`
	SpeechID  string            `json:"speech_id"`
	Status    string            `json:"status"`
}

func (s *Server) command(c *client, m command) error {
	if m.Type == "update_settings" {
		if err := s.Hub.updateSettings(m.Settings); err != nil {
			s.Hub.sendSettingsError(c, err)
		}
		return nil
	}
	id := m.AgentID
	if id == "" {
		id = s.DefaultAgent
	}
	a, ok := s.Root.Agent(id)
	if !ok {
		return errors.New("unknown agent")
	}
	switch m.Type {
	case "save":
		return a.Save()
	case "subscribe":
		return nil
	case "input", "prompt":
		if strings.TrimSpace(m.Text) == "" {
			return errors.New("empty input")
		}
		return a.Post(ensemble.UserMessage{Text: m.Text})
	case "hint":
		return a.Post(ensemble.Hint{Text: m.Text})
	case "pause":
		return a.Pause()
	case "unpause":
		if m.Automatic {
			return a.Unpause()
		}
		// Explicit Continue is user authority in every tab. Publish it before
		// the asynchronous state acknowledgement so fresh typing after this
		// notification cannot be cleared again by an older continuation epoch.
		s.Hub.mu.Lock()
		defer s.Hub.mu.Unlock()
		s.Hub.sync(s.Hub.view(a), a.View())
		if err := a.Unpause(); err != nil {
			return err
		}
		s.Hub.speech.cancelAgent(id, false)
		s.Hub.broadcast(Wire{Type: "continued", AgentID: id, Data: map[string]any{"client": c.id, "epoch": s.Hub.agents[id].state.Epoch}})
		return nil
	case "interrupt":
		return a.Post(ensemble.Interrupt{})
	case "load_skill":
		return a.Skills().LoadSkill(m.Name)
	case "unload_skill":
		return a.Skills().UnloadSkill(m.Name)
	case "skill_state":
		s.Hub.mu.Lock()
		s.Hub.broadcast(skillFrame(a))
		s.Hub.mu.Unlock()
		return nil
	case "speaking":
		return a.Post(ensemble.SetSpeaking{Enabled: m.Enabled})
	case "speak":
		// The server resolves a durable card; a browser cannot inject private replay
		// payloads or make a shortened preview stand in for the complete result.
		for _, e := range a.Log().EventsRange(0, a.View().Through) {
			for _, w := range projectEvent(id, e) {
				if w.PartID == m.PartID {
					if card := w.Part; card != nil && card.Text != "" && !strings.HasPrefix(card.Text, "[opaque") {
						s.Hub.mu.Lock()
						s.Hub.sync(s.Hub.view(a), a.View())
						s.Hub.speech.request(id, ensemble.SpeechRequested{Parts: speechParts(card.Text), Reason: "on demand"})
						s.Hub.mu.Unlock()
						return nil
					}
				}
			}
		}
		return errors.New("no readable final part")
	case "claim_speech":
		return s.Hub.speech.claim(c, id, m.SpeechID)
	case "speech_done":
		return s.Hub.speech.complete(c, id, m.SpeechID, m.Status)
	default:
		return errors.New("unknown command")
	}
}
