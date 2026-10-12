package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/flynn/noise"
	"github.com/gorilla/websocket"
	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// ServerVersion is reported to agents in the welcome frame.
const ServerVersion = "ensemble-gateway/1"

// Server is the gateway: it accepts authenticated agent connections and
// routes messages between them according to group membership.
//
// The gateway owns no agents. It is a router that happens to hold the
// connections, and every policy question it answers comes from config plus
// runtime declarations, never from the content of a message.
type Server struct {
	cfg      *Config
	registry *Registry
	queue    *Queue
	key      noise.DHKey
	conns    *connections
	logf     func(format string, args ...any)

	// origins remembers who sent a message, so a response can be routed
	// back to them. Replying to a message you received is always allowed:
	// the sender chose to start the exchange, so the reply needs no
	// separate contactability grant.
	originsMu sync.Mutex
	origins   map[string]string

	upgrader websocket.Upgrader
}

// NewServer builds a gateway over a validated config, restoring any runtime
// state previously persisted to the config's data directory.
func NewServer(cfg *Config, key noise.DHKey) (*Server, error) {
	registry := NewRegistry(cfg)
	if err := registry.Load(filepath.Join(cfg.Gateway.DataDir, "registry.json")); err != nil {
		return nil, err
	}
	queue := NewQueue()
	if err := queue.Load(filepath.Join(cfg.Gateway.DataDir, "queue.json")); err != nil {
		return nil, err
	}

	return &Server{
		cfg:      cfg,
		registry: registry,
		queue:    queue,
		key:      key,
		conns:    newConnections(),
		logf:     log.Printf,
		origins:  map[string]string{},
		upgrader: websocket.Upgrader{
			// Origin checking belongs to the HTTP layer and is meaningless
			// here: authentication is the Noise-KK handshake, which a
			// browser cannot perform and an attacker cannot forge without
			// the agent's static private key.
			CheckOrigin: func(*http.Request) bool { return true },
		},
	}, nil
}

// SetLogger replaces the server's log function, which tests use to keep
// output quiet and deterministic.
func (s *Server) SetLogger(logf func(format string, args ...any)) { s.logf = logf }

// Registry exposes the authorization state for the composition root.
func (s *Server) Registry() *Registry { return s.registry }

// Queue exposes the delivery queue for the composition root.
func (s *Server) Queue() *Queue { return s.queue }

// PublicKey returns the gateway's static public key, which every agent must
// be configured with before it can connect.
func (s *Server) PublicKey() []byte { return s.key.Public }

// Handler returns the gateway's HTTP routes.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/ws", s.handleWS)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
}

// Persist writes runtime state to the data directory.
func (s *Server) Persist() error {
	dir := s.cfg.Gateway.DataDir
	if err := s.registry.Save(filepath.Join(dir, "registry.json")); err != nil {
		return err
	}
	return s.queue.Save(filepath.Join(dir, "queue.json"))
}

// Close tears down every connection.
func (s *Server) Close() { s.conns.closeAll() }

// Online reports which agents currently hold a connection.
func (s *Server) Online() []string { return s.conns.online() }

// --- connection lifecycle ---

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logf("gateway: websocket upgrade failed: %v", err)
		return
	}

	agentID, session, err := s.doHandshake(ws)
	if err != nil {
		// The peer is unauthenticated at this point, so there is no
		// encrypted channel to report the error over. Log and hang up.
		s.logf("gateway: handshake failed: %v", err)
		_ = ws.Close()
		return
	}

	c := newConn(agentID, ws, session)
	s.conns.add(c)
	defer s.conns.remove(c)

	s.logf("gateway: %s connected", agentID)
	s.readLoop(c)
	s.logf("gateway: %s disconnected", agentID)
}

// doHandshake performs the Noise-KK handshake over raw WebSocket frames.
//
// The first message is not encrypted and carries the agent's id in the
// clear, prefixed by its length. The responder has no other way to know
// which static public key to expect: KK requires both peers' static keys
// before the session can be constructed, so the initiator must say who it
// claims to be first. Claiming an identity proves nothing. The handshake
// then either completes, proving possession of that identity's private key,
// or fails.
func (s *Server) doHandshake(ws *websocket.Conn) (string, *pub.NoiseKK, error) {
	_, msg1Raw, err := ws.ReadMessage()
	if err != nil {
		return "", nil, fmt.Errorf("read handshake msg1: %w", err)
	}
	if len(msg1Raw) < 2 {
		return "", nil, fmt.Errorf("handshake msg1 too short")
	}
	idLen := int(msg1Raw[0])
	if len(msg1Raw) < 1+idLen {
		return "", nil, fmt.Errorf("handshake msg1 agent id truncated")
	}
	agentID := string(msg1Raw[1 : 1+idLen])
	noiseMsg1 := msg1Raw[1+idLen:]

	agentKey, err := s.registry.PublicKey(agentID)
	if err != nil {
		return "", nil, fmt.Errorf("handshake from %q: %w", agentID, err)
	}

	session, err := pub.NewNoiseKK("responder", s.key, agentKey, []byte(pub.ProtocolVersion))
	if err != nil {
		return "", nil, fmt.Errorf("create noise session: %w", err)
	}
	if _, err := session.ReadHandshakeMessage(noiseMsg1); err != nil {
		return "", nil, fmt.Errorf("noise read msg1 from %q: %w", agentID, err)
	}
	msg2, err := session.WriteHandshakeMessage(nil)
	if err != nil {
		return "", nil, fmt.Errorf("noise write msg2: %w", err)
	}
	if err := ws.WriteMessage(websocket.BinaryMessage, msg2); err != nil {
		return "", nil, fmt.Errorf("send handshake msg2: %w", err)
	}
	if !session.IsHandshakeComplete() {
		return "", nil, fmt.Errorf("handshake incomplete after two messages")
	}
	return agentID, session, nil
}

func (s *Server) readLoop(c *conn) {
	for {
		frame, err := c.readFrame()
		if err != nil {
			return
		}
		if err := s.handleFrame(c, frame); err != nil {
			s.logf("gateway: %s frame %s: %v", c.agentID, frame.Type, err)
		}
	}
}

func (s *Server) handleFrame(c *conn, frame *pub.Frame) error {
	switch frame.Type {
	case pub.FrameHello:
		return s.handleHello(c, frame)
	case pub.FrameAck:
		var p pub.AckPayload
		if err := json.Unmarshal(frame.Payload, &p); err != nil {
			return err
		}
		s.queue.Ack(c.agentID, p.Seq)
		return nil
	case pub.FrameResponse:
		return s.handleResponse(c, frame)
	case pub.FramePing:
		return c.sendFrame(pub.FramePong, struct{}{})
	case pub.FrameMCPRequest:
		return s.handleMCPRequest(c, frame)
	default:
		return c.sendFrame(pub.FrameError, pub.ErrorPayload{
			Code:    pub.ErrCodeBadRequest,
			Message: fmt.Sprintf("unknown frame type %q", frame.Type),
		})
	}
}

// handleHello completes the session: it replays whatever the agent missed.
func (s *Server) handleHello(c *conn, frame *pub.Frame) error {
	var p pub.HelloPayload
	if len(frame.Payload) > 0 {
		if err := json.Unmarshal(frame.Payload, &p); err != nil {
			return err
		}
	}

	backlog := s.queue.Resume(c.agentID, p.LastAckedSeq)
	if err := c.sendFrame(pub.FrameWelcome, pub.WelcomePayload{
		AgentID:       c.agentID,
		QueuedCount:   len(backlog),
		ServerVersion: ServerVersion,
	}); err != nil {
		return err
	}

	for _, env := range backlog {
		if err := s.deliver(c, env); err != nil {
			return err
		}
	}
	return nil
}

// handleResponse routes an agent's reply back to whoever sent the prompt.
func (s *Server) handleResponse(c *conn, frame *pub.Frame) error {
	var p pub.ResponsePayload
	if err := json.Unmarshal(frame.Payload, &p); err != nil {
		return err
	}

	s.originsMu.Lock()
	origin, ok := s.origins[p.MessageID]
	delete(s.origins, p.MessageID)
	s.originsMu.Unlock()

	if !ok {
		return c.sendFrame(pub.FrameError, pub.ErrorPayload{
			Code:    pub.ErrCodeBadRequest,
			Message: fmt.Sprintf("no outstanding message %s to respond to", p.MessageID),
		})
	}

	reply := pub.Envelope{
		MessageID:      pub.NewMessageID(),
		AgentID:        origin,
		Timestamp:      time.Now().UTC(),
		ConversationID: p.Metadata["conversation_id"],
		Content:        p.Content,
		ResponseMode:   pub.ResponseFireAndForget,
		Source: pub.Source{
			Type:        "agent",
			ChannelID:   c.agentID,
			UserID:      c.agentID,
			DisplayName: s.displayName(c.agentID),
			Trust:       s.trustForReply(c.agentID, origin),
		},
	}
	s.dispatch(reply)
	return nil
}

// trustForReply labels a reply with the trust the responder would carry if
// it had initiated contact, and with the least privilege otherwise. A reply
// is always delivered, but it is never labeled with authority the
// relationship does not actually carry.
func (s *Server) trustForReply(responder, recipient string) pub.Trust {
	if _, trust, err := s.registry.Authorize(responder, recipient); err == nil {
		return trust
	}
	return pub.TrustExternal
}

func (s *Server) displayName(agentID string) string {
	if a, ok := s.cfg.Agents[agentID]; ok && a.DisplayName != "" {
		return a.DisplayName
	}
	return agentID
}

// --- routing ---

// SendMessage routes a message from one agent to another, enforcing
// group-based authorization.
//
// This is the single choke point for agent-to-agent traffic. Authorization
// happens here, before the message is queued, so an unauthorized send never
// becomes a durable obligation to deliver.
func (s *Server) SendMessage(sender string, msg pub.OutboundMessage) (pub.Envelope, error) {
	if !s.registry.Known(msg.ToAgent) {
		return pub.Envelope{}, &RouteError{Code: pub.ErrCodeUnknownAgent,
			Message: fmt.Sprintf("unknown agent %q", msg.ToAgent)}
	}
	if !s.registry.HasCapability(sender, pub.CapSendMessage) {
		return pub.Envelope{}, &RouteError{Code: pub.ErrCodeForbidden,
			Message: fmt.Sprintf("%s lacks the %s capability", sender, pub.CapSendMessage)}
	}

	_, trust, err := s.registry.Authorize(sender, msg.ToAgent)
	if err != nil {
		return pub.Envelope{}, &RouteError{Code: pub.ErrCodeUnauthorized, Message: err.Error()}
	}

	mode := msg.ResponseMode
	if mode == "" {
		mode = pub.ResponseAsync
	}
	env := pub.Envelope{
		MessageID:      pub.NewMessageID(),
		AgentID:        msg.ToAgent,
		Timestamp:      time.Now().UTC(),
		ConversationID: msg.ConversationID,
		Content:        msg.Content,
		ResponseMode:   mode,
		Source: pub.Source{
			Type:        "agent",
			ChannelID:   sender,
			UserID:      sender,
			DisplayName: s.displayName(sender),
			// Trust comes from the group that authorized this contact,
			// never from anything the sender supplied.
			Trust: trust,
		},
	}

	if mode != pub.ResponseFireAndForget {
		s.originsMu.Lock()
		s.origins[env.MessageID] = sender
		s.originsMu.Unlock()
	}

	return s.dispatch(env), nil
}

// dispatch queues an envelope and delivers it if the recipient is online.
//
// Queue first, then deliver. The reverse order would lose a message that
// was written to a socket the recipient never read, because nothing would
// remain to redeliver.
func (s *Server) dispatch(env pub.Envelope) pub.Envelope {
	env = s.queue.Enqueue(env)

	if c, ok := s.conns.get(env.AgentID); ok {
		if err := s.deliver(c, env); err != nil {
			s.logf("gateway: deliver to %s failed, leaving queued: %v", env.AgentID, err)
		}
	}
	return env
}

func (s *Server) deliver(c *conn, env pub.Envelope) error {
	// The provenance header is added at delivery, not at acceptance, so
	// what is stored in the queue stays the original message and cannot
	// accumulate a second header if it is redelivered.
	delivered := env
	delivered.Content = env.WithProvenance()
	return c.sendFrame(pub.FrameDeliver, pub.DeliverPayload{Seq: env.Seq, Envelope: delivered})
}

// RouteError is a routing refusal carrying a stable, matchable code.
type RouteError struct {
	Code    string
	Message string
}

func (e *RouteError) Error() string { return e.Message }

// --- lifecycle helpers ---

// Serve runs the gateway until ctx is cancelled.
func (s *Server) Serve(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		s.Close()
		return s.Persist()
	}
}
