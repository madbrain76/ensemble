package gui

import (
	"encoding/json"
	"ensemble"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const clientQueue = 64

type controlState struct {
	Paused        bool               `json:"paused"`
	Speaking      bool               `json:"speaking"`
	Reason        string             `json:"reason,omitempty"`
	Turn          ensemble.TurnState `json:"turn"`
	Epoch         uint64             `json:"epoch"`
	ContinueEpoch uint64             `json:"continue_epoch"`
}

func controlsOf(s ensemble.ViewSnapshot) controlState {
	return controlState{s.Paused, s.Speaking, s.Reason, s.State, s.ControlEpoch, s.ContinueEpoch}
}

type agentView struct {
	agent              ensemble.Agent
	cursor             ensemble.Seq
	revision, speechID uint64
	state              controlState
	partial            map[uint64]Card
	spoken             map[uint64]int
	turn               uint64
}
type client struct {
	id    string
	queue chan []byte
	done  chan struct{}
	once  sync.Once
}

func (c *client) stop() { c.once.Do(func() { close(c.done) }) }

type Hub struct {
	mu          sync.Mutex
	en          *ensemble.Ensemble
	agents      map[string]*agentView
	clients     map[*client]bool
	nextClient  uint64
	session     string
	unsubscribe func()
	closed      bool
	speech      *speaker
}

func NewHub(en *ensemble.Ensemble) *Hub {
	h := &Hub{en: en, agents: map[string]*agentView{}, clients: map[*client]bool{}, session: fmt.Sprint(time.Now().UnixNano())}
	h.speech = newSpeaker(h)
	for _, a := range en.Agents() {
		v := h.view(a)
		s := a.View()
		v.cursor = s.Through
		v.revision = s.Revision
		v.speechID = s.SpeechID
		v.state = controlsOf(s)
		v.partial = partialCards(s)
		v.turn = s.TurnID
	}
	h.unsubscribe = en.SubscribeIdentified(h.observe)
	return h
}
func (h *Hub) Close() {
	h.unsubscribe()
	h.mu.Lock()
	defer h.mu.Unlock()
	h.closed = true
	h.speech.close()
	for c := range h.clients {
		c.stop()
	}
	h.clients = map[*client]bool{}
}
func (h *Hub) trace(kind, agent, detail string) {
	h.en.Trace(ensemble.TraceRecord{Stage: "gui_" + kind, AgentID: agent, Detail: detail})
}
func (h *Hub) view(a ensemble.Agent) *agentView {
	v := h.agents[a.ID()]
	if v == nil {
		cfg := a.Config()
		v = &agentView{agent: a, state: controlState{Paused: cfg.PauseOnStart, Speaking: cfg.Speaking, Turn: ensemble.Idle}, partial: map[uint64]Card{}, spoken: map[uint64]int{}}
		h.agents[a.ID()] = v
	}
	return v
}
func partialCards(s ensemble.ViewSnapshot) map[uint64]Card {
	parts := make(map[uint64]Card, len(s.Partials))
	for _, p := range s.Partials {
		parts[p.PartID] = Card{Kind: string(p.Kind), Role: "agent", Name: p.Name, Text: p.Chunk, Partial: true, TurnID: p.TurnID}
	}
	return parts
}

// sync consumes exactly one immutable saved checkpoint. No view callback runs
// on the actor, and a dropped/coalesced wakeup cannot lose content or a final.
func (h *Hub) sync(v *agentView, s ensemble.ViewSnapshot) {
	if s.Revision == v.revision {
		return
	}
	id := v.agent.ID()
	old := v.state
	v.state = controlsOf(s)
	v.turn = s.TurnID
	if old != v.state {
		h.broadcast(Wire{Type: "state_changed", AgentID: id, From: old.Turn, To: v.state.Turn, State: &v.state})
	}
	if old.Epoch != v.state.Epoch {
		h.speech.cancelAgent(id, old.Speaking && !v.state.Speaking)
	}
	for _, e := range v.agent.Log().EventsRange(v.cursor, s.Through) {
		for _, f := range projectEvent(id, e) {
			if f.Type != "job_update" {
				delete(v.partial, f.PartID)
			}
			h.broadcast(f)
			if v.state.Speaking && f.Type != "job_update" {
				if f.Part != nil {
					h.narrate(v, f.PartID, *f.Part, true)
				}
			}
		}
		if e.Type == ensemble.TurnEnded || e.Type == ensemble.TurnInterrupted {
			v.spoken = map[uint64]int{}
		}
	}
	v.cursor = s.Through
	next := partialCards(s)
	for id := range v.partial {
		if _, ok := next[id]; !ok {
			h.broadcast(Wire{Type: "part_removed", AgentID: v.agent.ID(), PartID: id})
		}
	}
	ids := make([]uint64, 0, len(next))
	for id := range next {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for _, part := range ids {
		c := next[part]
		old, ok := v.partial[part]
		if !ok || old.Text != c.Text || old.Name != c.Name {
			h.broadcast(Wire{Type: "part_delta", AgentID: id, PartID: part, TurnID: c.TurnID, Kind: c.Kind, Name: c.Name, Chunk: c.Text[len(old.Text):]})
			if v.state.Speaking {
				h.narrate(v, part, c, false)
			}
		}
	}
	v.partial = next
	v.revision = s.Revision
	if s.SpeechID > v.speechID {
		v.speechID = s.SpeechID
		if s.SpeechEpoch == s.ControlEpoch {
			h.speech.request(id, s.Speech)
		}
	}
}
func (h *Hub) observe(p ensemble.AgentObservation) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.closed {
		h.sync(h.view(p.Agent), p.Agent.View())
	}
}
func (h *Hub) narrate(v *agentView, id uint64, c Card, final bool) {
	if c.Kind == "tool_call" {
		if final && v.spoken[id] == 0 {
			v.spoken[id] = 1
			h.speech.request(v.agent.ID(), ensemble.SpeechRequested{Parts: speechParts(c.Name + ": " + abbreviate(string(c.Args), 160)), Reason: "automatic"})
		}
		return
	}
	if c.Kind != "text" && c.Kind != "thinking" {
		return
	}
	if c.Text == "" || strings.HasPrefix(c.Text, "[opaque") {
		return
	}
	reason := "automatic"
	if c.Role == "human" {
		reason = "prompt"
	}
	n := v.spoken[id]
	if n > len(c.Text) {
		n = 0
	}
	end := len(c.Text)
	if !final {
		i := strings.LastIndexAny(c.Text[n:], ".!?\n")
		if i < 0 {
			return
		}
		end = n + i + 1
	}
	if end > n {
		h.speech.request(v.agent.ID(), ensemble.SpeechRequested{Parts: speechParts(c.Text[n:end]), Reason: reason})
		v.spoken[id] = end
	}
	if final {
		delete(v.spoken, id)
	}
}
func abbreviate(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}
func (h *Hub) broadcast(w Wire) {
	raw, err := json.Marshal(w)
	if err != nil {
		h.trace("encode_error", w.AgentID, err.Error())
		return
	}
	for c := range h.clients {
		h.enqueue(c, raw, w.Type, w.AgentID)
	}
}
func (h *Hub) enqueue(c *client, raw []byte, kind, agent string) {
	select {
	case c.queue <- raw:
	default:
		h.trace("slow_client_disconnect", agent, c.id)
		delete(h.clients, c)
		c.stop()
	}
}

type replayAgent struct {
	ID     string `json:"id"`
	Config struct {
		Vendor  ensemble.Vendor  `json:"vendor"`
		Model   string           `json:"model"`
		Surface ensemble.Surface `json:"surface"`
	} `json:"config"`
	State controlState `json:"state"`
}
type replayEnvelope struct {
	Session  string        `json:"session"`
	Client   string        `json:"client"`
	Agents   []replayAgent `json:"agents"`
	Events   []Wire        `json:"events"`
	Partials []Wire        `json:"partials"`
}

func (h *Hub) connect() *client {
	h.mu.Lock()
	defer h.mu.Unlock()
	snapshots := map[string]ensemble.ViewSnapshot{}
	for _, a := range h.en.Agents() {
		snapshots[a.ID()] = a.View()
	}
	return h.connectAt(snapshots)
}

// connectAt consumes the already captured checkpoint under the hub lock.
func (h *Hub) connectAt(snapshots map[string]ensemble.ViewSnapshot) *client {
	h.nextClient++
	c := &client{id: fmt.Sprint(h.nextClient), queue: make(chan []byte, clientQueue), done: make(chan struct{})}
	if h.closed {
		c.stop()
		return c
	}
	for id, s := range snapshots {
		if a, ok := h.en.Agent(id); ok {
			h.sync(h.view(a), s)
		}
	}
	data := replayEnvelope{Session: h.session, Client: c.id, Agents: []replayAgent{}, Events: []Wire{}, Partials: []Wire{}}
	ids := make([]string, 0, len(snapshots))
	for id := range snapshots {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		v := h.agents[id]
		s := snapshots[id]
		for _, e := range v.agent.Log().RecentThrough(s.Through, time.Now().Add(-time.Hour), 100, eventVisible) {
			frames := projectEvent(id, e)
			for i := range frames {
				frames[i].Replay = true
			}
			data.Events = append(data.Events, frames...)
		}
		cfg := v.agent.Config()
		a := replayAgent{ID: id, State: controlsOf(s)}
		a.Config.Vendor = cfg.Vendor
		a.Config.Model = cfg.Model
		a.Config.Surface = cfg.Surface
		data.Agents = append(data.Agents, a)
		data.Events = append(data.Events, Wire{Type: "state_changed", AgentID: id, Replay: true, From: s.State, To: s.State, State: &a.State})
		for _, p := range s.Partials {
			data.Partials = append(data.Partials, Wire{Type: "part_delta", AgentID: id, PartID: p.PartID, TurnID: p.TurnID, Kind: string(p.Kind), Name: p.Name, Chunk: p.Chunk, Replay: true})
		}
	}
	raw, err := json.Marshal(Wire{Type: "replay", Data: data, Replay: true})
	if err != nil {
		c.stop()
		return c
	}
	c.queue <- raw
	h.clients[c] = true
	h.trace("client_connect", "", c.id)
	return c
}
func (h *Hub) disconnect(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
	c.stop()
	h.speech.cancelClient(c.id)
}
func (h *Hub) sendError(c *client, agent string, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	b, _ := json.Marshal(Wire{Type: "error", AgentID: agent, Data: err.Error()})
	h.enqueue(c, b, "error", agent)
}
func (h *Hub) emit(w Wire) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.closed {
		h.broadcast(w)
	}
}
