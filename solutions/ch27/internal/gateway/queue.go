package gateway

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// Queue holds messages per recipient until the recipient acknowledges them.
//
// Delivery is at-least-once, which is a deliberate choice over at-most-once.
// A duplicated prompt is a visible annoyance the recipient can suppress by
// message ID; a dropped prompt is silent and unrecoverable. So a message
// stays queued until the recipient says it has it, and anything unacked is
// redelivered on reconnect.
//
// Sequence numbers are assigned per recipient, not globally. A recipient
// resuming after a disconnect reports the last sequence it handled, and the
// gateway replays from the next one. A global counter would force every
// recipient to reason about gaps caused by traffic that was never theirs.
type Queue struct {
	mu    sync.Mutex
	state queueState
}

type queueState struct {
	// NextSeq is the next sequence number to assign per recipient.
	NextSeq map[string]uint64 `json:"next_seq,omitempty"`
	// Pending holds unacknowledged envelopes per recipient, in sequence
	// order.
	Pending map[string][]pub.Envelope `json:"pending,omitempty"`
}

// NewQueue returns an empty queue.
func NewQueue() *Queue {
	return &Queue{state: queueState{
		NextSeq: map[string]uint64{},
		Pending: map[string][]pub.Envelope{},
	}}
}

// Enqueue assigns the next sequence number for the recipient, stores the
// envelope as unacknowledged, and returns it with Seq filled in.
func (q *Queue) Enqueue(env pub.Envelope) pub.Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.state.NextSeq[env.AgentID]++
	env.Seq = q.state.NextSeq[env.AgentID]
	q.state.Pending[env.AgentID] = append(q.state.Pending[env.AgentID], env)
	return env
}

// Ack discards every message for agentID at or below seq.
//
// Acking a sequence implicitly acks everything before it, which means a
// recipient that processes a batch needs one ack rather than one per
// message, and a lost ack is repaired by the next one.
func (q *Queue) Ack(agentID string, seq uint64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.ackLocked(agentID, seq)
}

func (q *Queue) ackLocked(agentID string, seq uint64) {
	pending := q.state.Pending[agentID]
	kept := pending[:0]
	for _, env := range pending {
		if env.Seq > seq {
			kept = append(kept, env)
		}
	}
	if len(kept) == 0 {
		delete(q.state.Pending, agentID)
		return
	}
	q.state.Pending[agentID] = kept
}

// Resume acks everything through lastAcked and returns what remains, which
// is exactly the backlog a reconnecting agent must be sent.
func (q *Queue) Resume(agentID string, lastAcked uint64) []pub.Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.ackLocked(agentID, lastAcked)
	return append([]pub.Envelope(nil), q.state.Pending[agentID]...)
}

// Pending returns the unacknowledged messages for agentID.
func (q *Queue) Pending(agentID string) []pub.Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]pub.Envelope(nil), q.state.Pending[agentID]...)
}

// Depth reports how many messages are waiting for agentID.
func (q *Queue) Depth(agentID string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.state.Pending[agentID])
}

// Save persists the queue atomically.
//
// A queue that did not survive a gateway restart would turn a routine
// deploy into silent message loss for every agent that happened to be
// offline.
func (q *Queue) Save(path string) error {
	q.mu.Lock()
	data, err := json.MarshalIndent(q.state, "", "  ")
	q.mu.Unlock()
	if err != nil {
		return fmt.Errorf("marshal queue: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("write queue: %w", err)
	}
	return os.Rename(tmp, path)
}

// Load restores a persisted queue. A missing file means an empty queue.
func (q *Queue) Load(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read queue %s: %w", path, err)
	}
	var state queueState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("parse queue %s: %w", path, err)
	}
	if state.NextSeq == nil {
		state.NextSeq = map[string]uint64{}
	}
	if state.Pending == nil {
		state.Pending = map[string][]pub.Envelope{}
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	q.state = state
	return nil
}
