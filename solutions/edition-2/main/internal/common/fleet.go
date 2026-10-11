package common

import "context"

// UserSink is the application's delivery boundary for questions and escalation.
// Deliver must return promptly; a question's reply comes through Fleet.Reply.
type UserSink interface {
	// Deliver presents the sender and message without impersonating a fleet member.
	Deliver(FleetMessage) error
}

// FleetMessage retains sender identity so questions and replies can be continued.
type FleetMessage struct {
	// ID identifies a pending question; notices need no reply.
	ID string `json:"message_id,omitempty"`
	// From identifies the sender in the recipient's immediate ownership tree.
	From string `json:"from"`
	// To is the receiving Agent or configured human address.
	To string `json:"to"`
	// Content carries the question, answer or escalation.
	Content string `json:"content"`
}

// ChildStatus is a snapshot of one owned child's completion, not a global registry.
type ChildStatus struct {
	// ID names this parent's child.
	ID string `json:"agent_id"`
	// State is running, submitted, or failed.
	State string `json:"state"`
	// Result is the explicit submitted payload, never inferred from an idle reply.
	Result string `json:"result,omitempty"`
	// Error records a failed turn without retrying it.
	Error string `json:"error,omitempty"`
}

// SubAgentManager owns child lifetimes and pending questions for one Agent.
type SubAgentManager interface {
	// Agent returns the required owner.
	Agent() Agent
	// Spawn constructs, restores and starts a child through the application factory.
	Spawn(id, task, model string) (ChildStatus, error)
	// List returns only this Agent's children in stable ID order.
	List() []ChildStatus
	// Join waits for all requested submissions, yielding early to pending questions.
	Join(context.Context, []string) (string, error)
	// ShutdownChild joins and saves one child, then removes it, preserving its files.
	ShutdownChild(string) error
	// Shutdown closes admission and recursively joins all owned children.
	Shutdown() error
	// AskParent blocks until the addressed parent or human replies.
	AskParent(context.Context, string) (FleetMessage, error)
	// Receive delivers a child question or notice into this parent's mailbox.
	Receive(FleetMessage, chan FleetMessage) error
	// Reply answers the specified pending question; omitted ID selects the sole match.
	Reply(agentID, messageID, content string) error
	// Submit records completion only when no unfinished children remain.
	Submit(string) error
	// Status reports this Agent's explicit completion state.
	Status() ChildStatus
	// Changed wakes fan-in observers after a child transition.
	Changed()
	// Begin starts a new task without forgetting consecutive idle corrections.
	Begin()
	// Idle returns a correction or escalation after a successful unsubmitted turn.
	Idle() (correction, escalation string)
	// Escalate reports failure to the parent without opening a blocking question.
	Escalate(string) error
}
