package agent

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Message is one message from an agent to whoever it reports to.
//
// A message is addressed, not merely posted, because a parent with several
// children has to know which one is asking before it can answer. The From
// field is set by the framework from the sending agent's own address and is
// never taken from the model, so a child cannot sign a message as a sibling.
type Message struct {
	// ID identifies this message so a reply can name which one it answers.
	// A parent handling three questions at once needs that, and so does a
	// parent answering a question asked two turns ago.
	ID string

	// From and To are agent addresses. To is either the sender's parent or
	// UserID.
	From string
	To   string

	// Content is the message text.
	Content string

	// Sent is when the message was posted, so a parent can see which of
	// several blocked children has been waiting longest.
	Sent time.Time

	// reply carries the answer back to the blocked sender. It is unexported
	// and buffered with room for exactly one answer, so a responder never
	// blocks even if the sender has already given up and walked away.
	//
	// It carries a whole Message rather than a bare string so the answer
	// arrives with its sender's address attached. An agent that receives
	// only text knows what it was told and not who told it, which is enough
	// to act on once and not enough to ask a follow-up question.
	//
	// A nil reply channel means nobody is waiting. That is how a message to
	// the user differs from a question to a parent: the first is an
	// announcement, the second is a question, and the channel is what
	// distinguishes them.
	reply chan *Message
}

// AwaitsReply reports whether a sender is blocked on this message.
func (m *Message) AwaitsReply() bool { return m.reply != nil }

// UserSink is how an agent reaches the human.
//
// It is an interface with one method rather than a channel because the human
// is reached differently depending on what is hosting the agent: a terminal
// prints, a GUI pushes over a websocket, a test records. None of those is
// the library's business, which is the same reason AgentFactory exists.
//
// It deliberately does not return a reply. A message to the user is an
// announcement; blocking an agent until a human answers is a different
// feature with different failure modes, and conflating the two would make
// every status update a potential deadlock.
type UserSink interface {
	DeliverToUser(m *Message)
}

// UserSinkFunc adapts a function to UserSink.
type UserSinkFunc func(m *Message)

// DeliverToUser calls f.
func (f UserSinkFunc) DeliverToUser(m *Message) { f(m) }

// parentLink is how a child reaches its parent. It is set by the manager
// that spawned the child, so a child cannot forge one or point it elsewhere.
type parentLink interface {
	deliver(m *Message) error
}

// messageSeq numbers messages. A counter rather than a timestamp because two
// messages posted in the same nanosecond must still have different ids, and
// on some platforms the clock is coarse enough for that to happen.
var messageSeq atomic.Uint64

func nextMessageID() string {
	return fmt.Sprintf("msg-%d", messageSeq.Add(1))
}

// ErrNoParentLink reports that an agent tried to reach a parent it has no
// link to. An empty comparable struct, so errors.Is still matches it through
// the wrapping that names the agent.
type ErrNoParentLink struct{}

func (ErrNoParentLink) Error() string { return "agent has no link to its parent" }

// ErrReplyTimeout reports that a sender gave up waiting for an answer.
type ErrReplyTimeout struct{}

func (ErrReplyTimeout) Error() string { return "timed out waiting for a reply" }

// MessageParent sends a message to whoever this agent reports to and blocks
// until the answer arrives or the timeout expires.
//
// For a top-level agent, "whoever it reports to" IS the user, and this
// routes there with no special case. That uniformity is the whole reason the
// user has an address: the alternative is every caller asking "am I at the
// top of the tree?" before deciding where to send, and getting it wrong once
// is a message that silently goes nowhere.
//
// Messaging is available at any point in an agent's work, not only when it
// finishes. Finishing is what submit means, and keeping the two verbs apart
// is what lets a child ask a question without claiming to be done.
func (a *Agent) MessageParent(content string, timeout time.Duration) (*Message, error) {
	m := &Message{
		ID:      nextMessageID(),
		From:    a.Spec().ID,
		To:      a.Spec().ParentID,
		Content: content,
		Sent:    time.Now(),
		reply:   make(chan *Message, 1),
	}

	// A top-level agent's parent is the user. The user does not answer, so
	// this is an announcement and the sender is not left blocked on a reply
	// that cannot come.
	if m.To == a.userAddr() {
		sink := a.userSink()
		if sink == nil {
			return nil, fmt.Errorf("agent %s: no user sink configured", m.From)
		}
		m.reply = nil
		sink.DeliverToUser(m)
		return nil, nil
	}

	link := a.parentLink()
	if link == nil {
		return nil, fmt.Errorf("%w: agent %s", ErrNoParentLink{}, m.From)
	}
	if err := link.deliver(m); err != nil {
		return nil, err
	}

	select {
	case answer := <-m.reply:
		return answer, nil
	case <-time.After(timeout):
		// The message stays in the parent's pending list. A late answer is
		// still delivered into the buffered channel and simply discarded,
		// which is why that channel has room for one value: a responder must
		// never block on a sender that has stopped listening.
		return nil, fmt.Errorf("%w after %s", ErrReplyTimeout{}, timeout)
	}
}

// MessageUser sends an announcement to the human, whatever depth this agent
// sits at.
//
// A deep child reaches the human directly rather than hoping each parent in
// the chain relays, because a relay only happens if every agent above
// chooses to perform it, and an agent that is stuck is exactly the one that
// will not.
func (a *Agent) MessageUser(content string) error {
	sink := a.userSink()
	if sink == nil {
		return fmt.Errorf("agent %s: no user sink configured", a.Spec().ID)
	}
	sink.DeliverToUser(&Message{
		ID:      nextMessageID(),
		From:    a.Spec().ID,
		To:      a.userAddr(),
		Content: content,
		Sent:    time.Now(),
	})
	return nil
}

// deliver places a child's message in that child's record and wakes anyone
// waiting on the fleet.
//
// The message is filed under the SENDER, not in one shared queue. A parent
// answering a question needs to know which child asked, and a flat queue
// makes that a search rather than a lookup.
func (m *SubAgentManager) deliver(msg *Message) error {
	m.mu.Lock()
	sub, ok := m.children[msg.From]
	if ok {
		// Under the child's own lock as well as the manager's, because
		// PendingMessages reads this map under the child's lock. The order
		// is always manager first, then child.
		sub.mu.Lock()
		if sub.pending == nil {
			sub.pending = make(map[string]*Message)
		}
		sub.pending[msg.ID] = msg
		sub.mu.Unlock()
	}
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("no such sub-agent: %s", msg.From)
	}
	m.wake()
	return nil
}

// wake releases anyone blocked on the fleet, such as a parent inside
// join_agents.
//
// This is what stops a parent deadlocking against its own child. Without it
// a child blocks waiting for an answer, the parent blocks waiting for the
// child to finish, and neither will move again.
//
// The send is non-blocking onto a channel with room for one. A pending
// wakeup that nobody has collected yet already says everything a second one
// would, so dropping it loses nothing.
func (m *SubAgentManager) wake() {
	select {
	case m.notify <- struct{}{}:
	default:
	}
}

// PendingMessages returns every message from a child that is waiting for an
// answer, ordered by child id and then by message id.
func (m *SubAgentManager) PendingMessages() []*Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	var out []*Message
	for _, sub := range m.sortedLocked() {
		sub.mu.Lock()
		for _, msg := range sub.pending {
			out = append(out, msg)
		}
		sub.mu.Unlock()
	}
	sortMessages(out)
	return out
}

// Respond answers a child's message and unblocks it.
//
// Answering removes the message from the pending list, so the same question
// cannot be answered twice with different advice.
func (m *SubAgentManager) Respond(agentID, messageID, response string) error {
	m.mu.Lock()
	sub, ok := m.children[agentID]
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("no such sub-agent: %s", agentID)
	}

	sub.mu.Lock()
	msg, found := sub.pending[messageID]
	if found {
		delete(sub.pending, messageID)
	}
	sub.mu.Unlock()

	if !found {
		return fmt.Errorf("sub-agent %s has no pending message %s", agentID, messageID)
	}
	if msg.reply == nil {
		return fmt.Errorf("message %s is an announcement and expects no reply", messageID)
	}
	// The answer carries the responder's address, so the child learns who
	// replied and can address a follow-up rather than only knowing what it
	// was told. Buffered with room for one, so this never blocks even when
	// the sender has already timed out and gone.
	msg.reply <- &Message{
		ID:      nextMessageID(),
		From:    m.parent.Spec().ID,
		To:      agentID,
		Content: response,
		Sent:    time.Now(),
	}
	return nil
}

// childLink is the parentLink handed to one child. It carries the manager
// and nothing else: the child's identity rides on the message itself, set
// from the child's own address by MessageParent.
type childLink struct {
	mgr *SubAgentManager
}

func (l childLink) deliver(msg *Message) error { return l.mgr.deliver(msg) }

// sortMessages orders messages by sender, then by send time, so that output
// a model reads is stable between turns.
func sortMessages(msgs []*Message) {
	for i := 1; i < len(msgs); i++ {
		for j := i; j > 0; j-- {
			a, b := msgs[j-1], msgs[j]
			if a.From < b.From || (a.From == b.From && !a.Sent.After(b.Sent)) {
				break
			}
			msgs[j-1], msgs[j] = b, a
		}
	}
}

// messagingState is embedded in Agent. It is a separate struct so the
// mutex that guards the links is obviously about the links and nothing else.
type messagingState struct {
	mu     sync.Mutex
	parent parentLink
	user   UserSink
	// addr is the human's address. Empty means the UserID default.
	addr string

	// submission is this agent's result, nil until it submits. It lives
	// here rather than on the parent so that finishing needs no
	// back-pointer: the parent reads it through the handle it already has.
	submission *Submission
}

func (a *Agent) parentLink() parentLink {
	a.msg.mu.Lock()
	defer a.msg.mu.Unlock()
	return a.msg.parent
}

func (a *Agent) setParentLink(l parentLink) {
	a.msg.mu.Lock()
	defer a.msg.mu.Unlock()
	a.msg.parent = l
}

func (a *Agent) userSink() UserSink {
	a.msg.mu.Lock()
	defer a.msg.mu.Unlock()
	return a.msg.user
}

// userAddr returns the human's address, defaulting to UserID.
func (a *Agent) userAddr() string {
	a.msg.mu.Lock()
	defer a.msg.mu.Unlock()
	if a.msg.addr == "" {
		return UserID
	}
	return a.msg.addr
}

// SetUser names the human and says how to reach them.
//
// Identity and delivery are set together deliberately. They are two halves
// of one fact, and setting only one of them produces a human who can be
// addressed but not reached, or reached but not named, with nothing to
// report either mistake.
//
// An address may be a name rather than a role. A message that arrives from
// "Bill" is more use to an agent than one from "user", and because the
// address is what a reply is sent to, the name has to be unique in the same
// way an agent id is. Pass an empty address to accept the UserID default.
//
// Every agent gets this, not just the top-level agent, because a deep child
// escalates to the human directly rather than relaying through parents that
// may themselves be stuck.
func (a *Agent) SetUser(addr string, s UserSink) {
	a.msg.mu.Lock()
	defer a.msg.mu.Unlock()
	a.msg.user = s
	a.msg.addr = addr
}
