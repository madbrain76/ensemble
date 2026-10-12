package agent

import (
	"fmt"
	"time"
)

// MaxIdleWakes is how many times the framework wakes a sub-agent that stopped
// without submitting before it gives up and reports the agent to its parent.
//
// Three is a judgement, not a measurement. One wake is too few, because the
// most common cause is a model that simply forgot the last step and fixes it
// the moment it is told. Waking forever is worse than failing, because a
// parent waiting on a child that is being woken in a loop never learns that
// anything is wrong.
const MaxIdleWakes = 3

// A sub-agent is given a task and is expected to call submit when the task is
// finished. Stopping without submitting is not a result. It is the model
// losing the thread, and the fix is to say so and let it continue.
//
// This is deliberately reactive. Nothing polls and no timer runs, because
// going idle is an event the framework already observes: the turn ended.
// Correcting it at that moment costs nothing and needs no scheduler.
func (m *SubAgentManager) correctIdle(s *SubAgent) {
	s.mu.Lock()
	if s.wakes >= MaxIdleWakes {
		s.mu.Unlock()
		m.escalate(s)
		return
	}
	s.wakes++
	n := s.wakes
	s.mu.Unlock()

	// Starting a turn from here recurses only through the goroutine Start
	// launches, so the count is what bounds this, not the stack.
	s.Start(wakeText(n))
}

// wakeText tells the agent what happened and what to do about it. It names
// the count because an agent on its third wake should treat the task
// differently from one on its first: if the work cannot be done, saying so
// through submit is a result, and stopping silently is not.
func wakeText(wake int) string {
	return fmt.Sprintf(
		"You stopped without calling submit, so your task is not finished. "+
			"This is wake %d of %d. Continue the task, and call submit when "+
			"it is done. If the task cannot be done, call submit and say so.",
		wake, MaxIdleWakes)
}

// escalate reports a stuck agent to its parent and stops waking it.
//
// The report goes to the parent rather than to a log because somebody has to
// decide what happens next, and the framework is not entitled to that
// decision. The parent can retry the work, route around it, or give up. A log
// line would leave the parent waiting on a child that is never going to
// answer.
func (m *SubAgentManager) escalate(s *SubAgent) {
	msg := &Message{
		ID:   fmt.Sprintf("escalation-%s-%d", s.ID, time.Now().UnixNano()),
		From: s.ID,
		To:   m.parent.Spec().ID,
		Content: fmt.Sprintf(
			"Sub-agent %q stopped %d times without calling submit and will "+
				"not be woken again. It has produced no result. Decide "+
				"whether to retry it, do the work another way, or continue "+
				"without it.", s.ID, MaxIdleWakes),
		Sent: time.Now(),
	}
	// A failure here means the child is already gone, which is not a problem
	// worth reporting: an agent that no longer exists cannot be stuck.
	_ = m.deliver(msg)
}

// resetWakes clears the count. Submitting is the only thing that clears it,
// because submitting is the only evidence the framework has that the agent
// did what it was asked rather than stopping again in a new place.
func (s *SubAgent) resetWakes() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wakes = 0
}

// Wakes reports how many consecutive times the framework has woken this agent
// without it submitting.
func (s *SubAgent) Wakes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wakes
}
