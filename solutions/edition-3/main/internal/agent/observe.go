package agent

import (
	"ensemble/internal/common"
	"sync"
)

// Queues are owned by the subscriber, not an actor callback executing terminal
// I/O. There is no fixed buffer that could silently drop text or final notices.
type observationQueue struct {
	mu      sync.Mutex
	pending []common.Observation
	wake    chan struct{}
	owner   *agent
	closed  bool
}

func (a *agent) Observe() common.ObservationQueue {
	q := &observationQueue{owner: a, wake: make(chan struct{}, 1)}
	a.observersMu.Lock()
	if a.queues == nil {
		a.queues = make(map[*observationQueue]struct{})
	}
	a.queues[q] = struct{}{}
	a.observersMu.Unlock()
	return q
}
func (q *observationQueue) Wake() <-chan struct{} { return q.wake }
func (q *observationQueue) Drain() []common.Observation {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := q.pending
	q.pending = nil
	return out
}
func (q *observationQueue) push(o common.Observation) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	q.pending = append(q.pending, o)
	select {
	case q.wake <- struct{}{}:
	default:
	}
}
func (q *observationQueue) Close() {
	q.owner.observersMu.Lock()
	delete(q.owner.queues, q)
	q.owner.observersMu.Unlock()
	q.mu.Lock()
	q.closed = true
	q.pending = nil
	q.mu.Unlock()
}
