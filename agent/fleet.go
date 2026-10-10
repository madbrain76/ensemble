package agent

import (
	"fmt"
	"strings"
	"time"
)

// Submission is what a child reports when it finishes.
//
// It is stored on the child itself rather than pushed to the parent, which
// means a child needs no back-pointer to its parent in order to finish. The
// parent reads it through the SubAgent handle it already holds.
type Submission struct {
	// Payload is the child's structured result, as JSON text.
	Payload string

	// At is when the child submitted.
	At time.Time
}

// Submit records this agent's result and marks it finished.
//
// Submit is the ONLY signal that means finished. An agent that merely stops
// has not finished, it has stopped, and stage 4 treats the difference as a
// mistake to correct. Keeping a separate explicit verb for "done" is what
// makes that distinction checkable: without it, the framework would have to
// infer completion from silence, and silence is exactly what a confused
// model produces.
func (a *Agent) Submit(payload string) {
	a.msg.mu.Lock()
	a.msg.submission = &Submission{Payload: payload, At: time.Now()}
	a.msg.mu.Unlock()
}

// Submission returns this agent's result, or nil if it has not submitted.
func (a *Agent) Submission() *Submission {
	a.msg.mu.Lock()
	defer a.msg.mu.Unlock()
	return a.msg.submission
}

// Start runs one turn for this child, in the background.
//
// The turn goes through the same path a top-level agent's turn does. There
// is no separate "sub-agent turn" machinery, because a child IS an agent and
// a second path would be a second thing to keep correct.
//
// State is set from what the framework watches happen, never from what the
// child says about itself. Running is set before the turn starts and the
// terminal state after it returns, so a child cannot report itself finished
// by asserting it.
func (s *SubAgent) Start(task string) error {
	s.mu.Lock()
	if s.state == SubAgentRunning {
		s.mu.Unlock()
		return fmt.Errorf("sub-agent %s is already running", s.ID)
	}
	s.state = SubAgentRunning
	s.turns++
	s.mu.Unlock()
	s.notifyFleet()

	go func() {
		_, err := s.Agent.Ask(task)

		// A turn that ended in an error is FAILED, which is a different
		// thing from idle. A child killed on its first turn by a bad model
		// name used to be indistinguishable from one that ran and finished,
		// because both simply sat there.
		next := SubAgentIdle
		if err != nil {
			next = SubAgentFailed
		} else if s.Agent.Submission() != nil {
			next = SubAgentDone
		}
		s.setState(next, err)
		s.notifyFleet()
	}()
	return nil
}

// notifyFleet wakes anyone blocked on this child's fleet.
func (s *SubAgent) notifyFleet() {
	if s.mgr != nil {
		s.mgr.wake()
	}
}

// Finished reports whether this child has stopped for good.
//
// Submitted and failed both count. Idle deliberately does not: an idle child
// that never submitted is the silent failure this chapter exists to catch,
// and counting it as finished here would hide it from join and hand the
// parent an empty result as though it were an answer.
func (s *SubAgent) Finished() bool {
	switch s.State() {
	case SubAgentDone, SubAgentFailed:
		return true
	default:
		return false
	}
}

// JoinOutcome is how one child ended, as observed by the framework.
type JoinOutcome struct {
	// ID is the child's address.
	ID string

	// Outcome is one of "submitted", "failed", "idle_without_submit",
	// "timed_out" or "unknown_agent". It describes HOW the child ended as
	// the framework saw it, never as the child claimed.
	//
	// Note what is deliberately absent: a success flag. There is no
	// framework-level "it worked", because a self-reported one invites an
	// agent to tick its own box. "submitted" means a payload arrived in the
	// right shape, which is a claim to check and not a fact.
	Outcome string

	// Payload is the child's submitted result, empty unless it submitted.
	Payload string

	// Err is the error that ended the child's last turn, if any.
	Err error
}

// JoinResult is what a fan-in returns.
type JoinResult struct {
	Outcomes []JoinOutcome

	// Messages are questions from children that arrived while waiting. A
	// non-empty list is why Join returned early.
	Messages []*Message

	// AllFinished is true only if every named child ended. Check it first:
	// a caller that reads Outcomes without checking this can mistake an
	// early return for a completed fleet.
	AllFinished bool

	// Interrupted says the wait ended because a child asked a question
	// rather than because the fleet finished.
	Interrupted bool
}

// Join blocks until every named child finishes, or until a child asks a
// question, or until the timeout expires.
//
// Returning early on a question is what stops a parent deadlocking against
// its own child. The child blocks waiting for an answer, the parent blocks
// waiting for the child to finish, and without this neither ever moves
// again. The parent gets control back, answers, and calls Join again.
func (m *SubAgentManager) Join(ids []string, timeout time.Duration) *JoinResult {
	deadline := time.After(timeout)
	for {
		res := m.collect(ids)
		if res.AllFinished || res.Interrupted {
			return res
		}
		select {
		case <-m.notify:
			// Something changed: a child finished, or asked a question.
			// Loop and look rather than deciding here, so one wakeup that
			// covers several changes is handled in one pass.
		case <-deadline:
			res = m.collect(ids)
			for i := range res.Outcomes {
				if res.Outcomes[i].Outcome == "running" {
					res.Outcomes[i].Outcome = "timed_out"
				}
			}
			// A timed-out child is NOT stopped. Join reports the deadline;
			// killing work because a caller grew impatient is a different
			// decision, and it belongs to whoever asked.
			return res
		}
	}
}

// collect takes one reading of the fleet.
func (m *SubAgentManager) collect(ids []string) *JoinResult {
	res := &JoinResult{AllFinished: true}
	for _, id := range ids {
		sub, ok := m.Get(id)
		if !ok {
			// A typo'd id must not look like a child that finished, or a
			// fleet of five with one misspelled name reports four results
			// and success.
			res.Outcomes = append(res.Outcomes, JoinOutcome{ID: id, Outcome: "unknown_agent"})
			res.AllFinished = false
			continue
		}
		out := JoinOutcome{ID: id, Err: sub.Err()}
		switch sub.State() {
		case SubAgentDone:
			out.Outcome = "submitted"
			if s := sub.Agent.Submission(); s != nil {
				out.Payload = s.Payload
			}
		case SubAgentFailed:
			out.Outcome = "failed"
		case SubAgentIdle:
			if s := sub.Agent.Submission(); s != nil {
				out.Outcome = "submitted"
				out.Payload = s.Payload
			} else {
				out.Outcome = "idle_without_submit"
				res.AllFinished = false
			}
		default:
			out.Outcome = "running"
			res.AllFinished = false
		}
		res.Outcomes = append(res.Outcomes, out)
	}

	if msgs := m.PendingMessages(); len(msgs) > 0 {
		res.Messages = msgs
		res.Interrupted = true
	}
	return res
}

// Render writes a join result as text for the model.
func (r *JoinResult) Render() string {
	var b strings.Builder
	if r.Interrupted {
		b.WriteString("Returned early: a sub-agent is waiting for your answer.\n\n")
	}
	for _, o := range r.Outcomes {
		fmt.Fprintf(&b, "%s: %s", o.ID, o.Outcome)
		if o.Err != nil {
			fmt.Fprintf(&b, " (%v)", o.Err)
		}
		b.WriteString("\n")
		if o.Payload != "" {
			fmt.Fprintf(&b, "  result: %s\n", o.Payload)
		}
	}
	for _, msg := range r.Messages {
		fmt.Fprintf(&b, "\nmessage %s from %s: %s\n", msg.ID, msg.From, msg.Content)
		b.WriteString("  reply with respond_to_parent_message\n")
	}
	return b.String()
}

// Status renders the fleet for the model.
//
// Ordered by id so two readings a turn apart can be compared by eye, and so
// an unchanged fleet renders identically and does not cost a cache hit.
func (m *SubAgentManager) Status() string {
	subs := m.List()
	if len(subs) == 0 {
		return "No sub-agents."
	}
	var b strings.Builder
	for _, sub := range subs {
		submitted := "no"
		if sub.Agent.Submission() != nil {
			submitted = "yes"
		}
		fmt.Fprintf(&b, "%s: state=%s turns=%d submitted=%s", sub.ID, sub.State(), sub.Turns(), submitted)
		if err := sub.Err(); err != nil {
			fmt.Fprintf(&b, " error=%v", err)
		}
		// Turns is reported because a child with zero turns never ran at
		// all, which sends you to different logs than one that ran and
		// stopped.
		b.WriteString("\n")
	}
	if msgs := m.PendingMessages(); len(msgs) > 0 {
		b.WriteString("\nWaiting for your answer:\n")
		for _, msg := range msgs {
			fmt.Fprintf(&b, "  %s from %s: %s\n", msg.ID, msg.From, msg.Content)
		}
	}
	return b.String()
}
