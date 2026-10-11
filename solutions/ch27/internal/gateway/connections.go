package gateway

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

const writeTimeout = 10 * time.Second

// conn is one authenticated agent connection.
//
// Teardown uses a dedicated done channel that the connection manager closes
// exactly once. Nothing ever closes a channel that senders write to, because
// a send racing a close is a panic that takes the whole gateway down, and
// the window is small enough to survive testing and show up in production.
type conn struct {
	agentID string
	ws      *websocket.Conn
	noise   *pub.NoiseKK

	sendMu sync.Mutex

	doneOnce sync.Once
	done     chan struct{}
}

func newConn(agentID string, ws *websocket.Conn, session *pub.NoiseKK) *conn {
	return &conn{
		agentID: agentID,
		ws:      ws,
		noise:   session,
		done:    make(chan struct{}),
	}
}

// close releases the connection. Safe to call more than once.
func (c *conn) close() {
	c.doneOnce.Do(func() {
		close(c.done)
		_ = c.ws.Close()
	})
}

// sendFrame encrypts and writes one frame.
//
// The lock covers both the Noise encryption and the socket write. Noise
// cipher states are stateful counters: two goroutines encrypting
// concurrently would produce frames the peer decrypts out of order and
// rejects, so this is not merely about serializing writes.
func (c *conn) sendFrame(frameType pub.FrameType, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal %s payload: %w", frameType, err)
	}
	frameJSON, err := json.Marshal(pub.Frame{Type: frameType, Payload: payloadJSON})
	if err != nil {
		return fmt.Errorf("marshal %s frame: %w", frameType, err)
	}

	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	select {
	case <-c.done:
		return fmt.Errorf("connection to %s is closed", c.agentID)
	default:
	}

	ciphertext, err := c.noise.Encrypt(frameJSON, nil)
	if err != nil {
		return fmt.Errorf("encrypt %s: %w", frameType, err)
	}
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}
	return c.ws.WriteMessage(websocket.BinaryMessage, ciphertext)
}

// readFrame reads and decrypts one frame.
func (c *conn) readFrame() (*pub.Frame, error) {
	_, ciphertext, err := c.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	plaintext, err := c.noise.Decrypt(ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	var frame pub.Frame
	if err := json.Unmarshal(plaintext, &frame); err != nil {
		return nil, fmt.Errorf("unmarshal frame: %w", err)
	}
	return &frame, nil
}

// connections tracks the currently connected agents.
type connections struct {
	mu sync.RWMutex
	m  map[string]*conn
}

func newConnections() *connections {
	return &connections{m: map[string]*conn{}}
}

// add registers a connection, displacing any previous one for the same
// agent. A reconnect is far more often a network blip than a genuine second
// instance, so the newest connection wins and the stale one is torn down.
func (cs *connections) add(c *conn) {
	cs.mu.Lock()
	previous := cs.m[c.agentID]
	cs.m[c.agentID] = c
	cs.mu.Unlock()

	if previous != nil && previous != c {
		previous.close()
	}
}

// remove drops the connection if it is still the current one for its agent.
// The guard matters: a slow teardown must not evict a newer connection that
// has already replaced it.
func (cs *connections) remove(c *conn) {
	cs.mu.Lock()
	if cs.m[c.agentID] == c {
		delete(cs.m, c.agentID)
	}
	cs.mu.Unlock()
	c.close()
}

func (cs *connections) get(agentID string) (*conn, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	c, ok := cs.m[agentID]
	return c, ok
}

func (cs *connections) online() []string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	out := make([]string, 0, len(cs.m))
	for id := range cs.m {
		out = append(out, id)
	}
	return out
}

func (cs *connections) closeAll() {
	cs.mu.Lock()
	all := make([]*conn, 0, len(cs.m))
	for _, c := range cs.m {
		all = append(all, c)
	}
	cs.m = map[string]*conn{}
	cs.mu.Unlock()

	for _, c := range all {
		c.close()
	}
}
