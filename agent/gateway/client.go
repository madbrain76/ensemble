package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/flynn/noise"
	"github.com/gorilla/websocket"
)

// ClientConfig configures a connection to a gateway.
type ClientConfig struct {
	// GatewayURL is the WebSocket endpoint, such as
	// wss://gateway.example.com/v1/ws.
	GatewayURL string
	// AgentID is this agent's identity, which must be known to the gateway.
	AgentID string
	// StaticKey is this agent's long-term Noise keypair.
	StaticKey noise.DHKey
	// GatewayPubKey is the gateway's known static public key. KK requires
	// it up front: an agent that would accept any server key has no way to
	// tell the real gateway from one that merely answers on that address.
	GatewayPubKey []byte
	// LastAckedSeq resumes delivery after this sequence number. An agent
	// with durable state of its own sets this on reconnect; one without
	// leaves it zero and accepts redelivery of everything outstanding.
	LastAckedSeq uint64
	// HandshakeTimeout bounds the dial and handshake. Zero means ten
	// seconds.
	HandshakeTimeout time.Duration
}

// Client is a connection to a gateway.
//
// It owns the Noise session and a single receive goroutine. All outbound
// frames are serialized, because a Noise cipher state is a counter and two
// concurrent encryptions would produce frames the peer rejects.
type Client struct {
	cfg     ClientConfig
	handler Handler

	ws    *websocket.Conn
	noise *NoiseKK

	sendMu sync.Mutex

	doneOnce sync.Once
	done     chan struct{}

	// pending correlates MCP tool calls with their responses.
	pendingMu sync.Mutex
	pending   map[string]chan MCPResponsePayload

	// seen is the dedup window. Delivery is at-least-once, so the same
	// message can legitimately arrive twice after a reconnect; the handler
	// should see it once.
	seenMu sync.Mutex
	seen   map[string]bool

	// lastAcked is the highest sequence this client has acknowledged.
	lastAcked uint64
}

// Connect dials the gateway, completes the Noise-KK handshake, sends hello,
// and starts the receive goroutine.
func Connect(ctx context.Context, cfg ClientConfig, handler Handler) (*Client, error) {
	if handler == nil {
		return nil, fmt.Errorf("gateway: a Handler is required")
	}
	if cfg.AgentID == "" {
		return nil, fmt.Errorf("gateway: AgentID is required")
	}
	if len(cfg.GatewayPubKey) != PublicKeySize {
		return nil, fmt.Errorf("gateway: GatewayPubKey must be %d bytes, got %d",
			PublicKeySize, len(cfg.GatewayPubKey))
	}
	timeout := cfg.HandshakeTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	dialer := websocket.Dialer{HandshakeTimeout: timeout}
	ws, resp, err := dialer.DialContext(ctx, cfg.GatewayURL, http.Header{})
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("dial gateway %s: %w (http %s)", cfg.GatewayURL, err, resp.Status)
		}
		return nil, fmt.Errorf("dial gateway %s: %w", cfg.GatewayURL, err)
	}

	session, err := newClientSession(ws, cfg, timeout)
	if err != nil {
		_ = ws.Close()
		return nil, err
	}

	c := &Client{
		cfg:       cfg,
		handler:   handler,
		ws:        ws,
		noise:     session,
		done:      make(chan struct{}),
		pending:   map[string]chan MCPResponsePayload{},
		seen:      map[string]bool{},
		lastAcked: cfg.LastAckedSeq,
	}

	if err := c.sendFrame(FrameHello, HelloPayload{
		AgentID:      cfg.AgentID,
		Version:      ProtocolVersion,
		LastAckedSeq: cfg.LastAckedSeq,
	}); err != nil {
		_ = ws.Close()
		return nil, fmt.Errorf("send hello: %w", err)
	}

	go c.receiveLoop()
	return c, nil
}

// newClientSession performs the initiator half of the Noise-KK handshake.
//
// The first message carries the agent id in the clear, length-prefixed,
// because the gateway cannot build a KK session without knowing which
// static key to expect. The id is a claim, not a credential: the handshake
// completes only if this client holds the matching private key.
func newClientSession(ws *websocket.Conn, cfg ClientConfig, timeout time.Duration) (*NoiseKK, error) {
	session, err := NewNoiseKK("initiator", cfg.StaticKey, cfg.GatewayPubKey, []byte(ProtocolVersion))
	if err != nil {
		return nil, fmt.Errorf("create noise session: %w", err)
	}

	msg1, err := session.WriteHandshakeMessage(nil)
	if err != nil {
		return nil, fmt.Errorf("noise write msg1: %w", err)
	}
	if len(cfg.AgentID) > 255 {
		return nil, fmt.Errorf("agent id %q is too long for the handshake frame", cfg.AgentID)
	}
	framed := make([]byte, 0, 1+len(cfg.AgentID)+len(msg1))
	framed = append(framed, byte(len(cfg.AgentID)))
	framed = append(framed, cfg.AgentID...)
	framed = append(framed, msg1...)

	if err := ws.WriteMessage(websocket.BinaryMessage, framed); err != nil {
		return nil, fmt.Errorf("send handshake msg1: %w", err)
	}

	if err := ws.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	_, msg2, err := ws.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read handshake msg2: %w", err)
	}
	if err := ws.SetReadDeadline(time.Time{}); err != nil {
		return nil, err
	}
	if _, err := session.ReadHandshakeMessage(msg2); err != nil {
		return nil, fmt.Errorf("noise read msg2: %w", err)
	}
	if !session.IsHandshakeComplete() {
		return nil, fmt.Errorf("handshake incomplete after two messages")
	}
	return session, nil
}

// Done is closed when the connection ends.
func (c *Client) Done() <-chan struct{} { return c.done }

// Close disconnects.
func (c *Client) Close() error {
	c.doneOnce.Do(func() {
		close(c.done)
		_ = c.ws.Close()
	})
	return nil
}

func (c *Client) sendFrame(frameType FrameType, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	frameJSON, err := json.Marshal(Frame{Type: frameType, Payload: payloadJSON})
	if err != nil {
		return err
	}

	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	select {
	case <-c.done:
		return fmt.Errorf("gateway: connection closed")
	default:
	}

	ciphertext, err := c.noise.Encrypt(frameJSON, nil)
	if err != nil {
		return fmt.Errorf("encrypt %s: %w", frameType, err)
	}
	return c.ws.WriteMessage(websocket.BinaryMessage, ciphertext)
}

func (c *Client) receiveLoop() {
	reason := "connection closed"
	defer func() {
		c.Close()
		c.handler.OnDisconnected(reason)
	}()

	for {
		_, ciphertext, err := c.ws.ReadMessage()
		if err != nil {
			reason = err.Error()
			return
		}
		plaintext, err := c.noise.Decrypt(ciphertext, nil)
		if err != nil {
			reason = fmt.Sprintf("decrypt: %v", err)
			return
		}
		var frame Frame
		if err := json.Unmarshal(plaintext, &frame); err != nil {
			c.handler.OnError(fmt.Errorf("unmarshal frame: %w", err))
			continue
		}
		c.handleFrame(frame)
	}
}

func (c *Client) handleFrame(frame Frame) {
	switch frame.Type {
	case FrameWelcome:
		var p WelcomePayload
		if err := json.Unmarshal(frame.Payload, &p); err != nil {
			c.handler.OnError(err)
			return
		}
		c.handler.OnConnected(p)

	case FrameDeliver:
		var p DeliverPayload
		if err := json.Unmarshal(frame.Payload, &p); err != nil {
			c.handler.OnError(err)
			return
		}
		c.deliver(p)

	case FrameError:
		var p ErrorPayload
		if err := json.Unmarshal(frame.Payload, &p); err != nil {
			c.handler.OnError(err)
			return
		}
		c.handler.OnError(fmt.Errorf("gateway error %s: %s", p.Code, p.Message))

	case FrameMCPResponse:
		var p MCPResponsePayload
		if err := json.Unmarshal(frame.Payload, &p); err != nil {
			c.handler.OnError(err)
			return
		}
		c.pendingMu.Lock()
		ch, ok := c.pending[p.ID]
		delete(c.pending, p.ID)
		c.pendingMu.Unlock()
		if ok {
			ch <- p
		}

	case FramePong:
		// Keepalive answered; nothing to do.

	default:
		c.handler.OnError(fmt.Errorf("unknown frame type %q", frame.Type))
	}
}

// deliver hands one envelope to the handler, then acknowledges it.
//
// The ack follows the handler rather than preceding it, so the gateway only
// forgets a message once the application has accepted it. A duplicate that
// arrives after a reconnect is suppressed here rather than troubling the
// handler, which is the whole purpose of an at-least-once channel.
func (c *Client) deliver(p DeliverPayload) {
	c.seenMu.Lock()
	duplicate := c.seen[p.Envelope.MessageID]
	if !duplicate {
		c.seen[p.Envelope.MessageID] = true
	}
	c.seenMu.Unlock()

	if !duplicate {
		c.handler.OnDeliver(p.Envelope)
	}

	if err := c.sendFrame(FrameAck, AckPayload{Seq: p.Seq}); err != nil {
		c.handler.OnError(fmt.Errorf("ack seq %d: %w", p.Seq, err))
		return
	}
	c.seenMu.Lock()
	if p.Seq > c.lastAcked {
		c.lastAcked = p.Seq
	}
	c.seenMu.Unlock()
}

// LastAckedSeq reports the highest sequence this client has acknowledged,
// which is what to pass as ClientConfig.LastAckedSeq when reconnecting.
func (c *Client) LastAckedSeq() uint64 {
	c.seenMu.Lock()
	defer c.seenMu.Unlock()
	return c.lastAcked
}

// --- administration tool calls ---

// OutboundMessage is a message to another agent.
type OutboundMessage struct {
	ToAgent        string
	Content        string
	ConversationID string
	ResponseMode   ResponseMode
}

// Call invokes a gateway administration tool and waits for its result.
func (c *Client) Call(ctx context.Context, tool string, args any) (json.RawMessage, error) {
	var raw json.RawMessage
	if args != nil {
		encoded, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}
		raw = encoded
	}

	id := NewMessageID()
	replyCh := make(chan MCPResponsePayload, 1)
	c.pendingMu.Lock()
	c.pending[id] = replyCh
	c.pendingMu.Unlock()

	cleanup := func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}

	if err := c.sendFrame(FrameMCPRequest, MCPRequestPayload{ID: id, Tool: tool, Arguments: raw}); err != nil {
		cleanup()
		return nil, err
	}

	select {
	case reply := <-replyCh:
		if reply.IsError {
			return nil, fmt.Errorf("%s: %s", tool, reply.Error)
		}
		return reply.Result, nil
	case <-ctx.Done():
		cleanup()
		return nil, ctx.Err()
	case <-c.done:
		cleanup()
		return nil, fmt.Errorf("gateway: connection closed while calling %s", tool)
	}
}

// Send delivers a message to another agent through the gateway.
func (c *Client) Send(ctx context.Context, msg OutboundMessage) error {
	_, err := c.Call(ctx, ToolSendMessage, sendMessageArgs{
		ToAgent:        msg.ToAgent,
		Content:        msg.Content,
		ConversationID: msg.ConversationID,
		ResponseMode:   string(msg.ResponseMode),
	})
	return err
}

// Respond replies to a previously delivered envelope. The gateway routes
// the reply back to whoever sent the original message.
func (c *Client) Respond(ctx context.Context, messageID, content string) error {
	return c.sendFrame(FrameResponse, ResponsePayload{MessageID: messageID, Content: content})
}

// ChildRegistration describes a child agent to register with the gateway.
type ChildRegistration struct {
	AgentID   string
	PublicKey []byte
	// Groups must be a subset of the registering agent's own memberships.
	// The gateway refuses the whole registration otherwise rather than
	// silently granting the part that was allowed.
	Groups []string
}

// RegisterChild registers a child agent's public key and grants it a subset
// of this agent's groups.
func (c *Client) RegisterChild(ctx context.Context, reg ChildRegistration) error {
	_, err := c.Call(ctx, ToolRegisterChild, registerChildArgs{
		AgentID:   reg.AgentID,
		PublicKey: EncodePublicKey(reg.PublicKey),
		Groups:    reg.Groups,
	})
	return err
}

// DeclareContactable tells the gateway which groups may send messages to
// this agent. Inbound contact is opt-in, so an agent that never calls this
// receives nothing.
func (c *Client) DeclareContactable(ctx context.Context, groups []string) error {
	_, err := c.Call(ctx, ToolDeclareContactable, declareContactableArgs{Groups: groups})
	return err
}

// Ping sends a keepalive.
func (c *Client) Ping() error { return c.sendFrame(FramePing, struct{}{}) }

// argument shapes shared with the server side of the protocol.

type sendMessageArgs struct {
	ToAgent        string `json:"to_agent"`
	Content        string `json:"content"`
	ConversationID string `json:"conversation_id,omitempty"`
	ResponseMode   string `json:"response_mode,omitempty"`
}

type registerChildArgs struct {
	AgentID   string   `json:"agent_id"`
	PublicKey string   `json:"public_key"`
	Groups    []string `json:"groups"`
}

type declareContactableArgs struct {
	Groups []string `json:"groups"`
}
